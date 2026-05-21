// Package specialweapon 特殊武器系統（DAY-089，升級 DAY-134）
// 業界依據：
//   - Fish Road 2026 有 8 tier 武器系統，炸彈/雷射是標配特殊武器
//   - Royal Fishing 2026 Tornado Cannon — 龍捲風掃場，旋轉吸入所有目標
//   - JILI 2026 Auto-Charge — 每次擊破目標自動累積充能，不需要花金幣
package specialweapon

import (
	"math"
	"sync"
)

// WeaponType 特殊武器類型
type WeaponType string

const (
	WeaponBomb    WeaponType = "bomb"    // 炸彈砲：範圍爆炸，命中半徑 200px 內所有目標
	WeaponLaser   WeaponType = "laser"   // 雷射砲：穿透，命中一條線上所有目標（Y 軸 ±60px）
	WeaponFreeze  WeaponType = "freeze"  // 冰凍砲：全畫面減速所有目標 5 秒
	WeaponTornado WeaponType = "tornado" // 龍捲風砲：全螢幕旋轉，50% 機率擊破所有目標（DAY-134）
)

// WeaponDef 特殊武器定義
type WeaponDef struct {
	Type           WeaponType `json:"type"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Cost           int        `json:"cost"`            // 購買一發的金幣費用（0 = 只能充能獲得）
	MaxCharges     int        `json:"max_charges"`     // 最多持有發數
	Icon           string     `json:"icon"`            // 圖示（emoji）
	Color          string     `json:"color"`           // 顏色（hex）
	ChargeRequired int        `json:"charge_required"` // 自動充能所需點數（DAY-134）
	ChargePerKill  int        `json:"charge_per_kill"` // 每次擊破獲得的充能點數（DAY-134）
}

// AvailableWeapons 所有特殊武器定義
var AvailableWeapons = []WeaponDef{
	{
		Type:           WeaponBomb,
		Name:           "炸彈砲",
		Description:    "範圍爆炸，命中半徑200px內所有目標",
		Cost:           500,
		MaxCharges:     3,
		Icon:           "💣",
		Color:          "#FF6B35",
		ChargeRequired: 20, // 擊破 20 個目標自動充能一發
		ChargePerKill:  1,
	},
	{
		Type:           WeaponLaser,
		Name:           "雷射砲",
		Description:    "穿透射擊，命中同一水平線上所有目標",
		Cost:           800,
		MaxCharges:     3,
		Icon:           "⚡",
		Color:          "#00FFFF",
		ChargeRequired: 30, // 擊破 30 個目標自動充能一發
		ChargePerKill:  1,
	},
	{
		Type:           WeaponFreeze,
		Name:           "冰凍砲",
		Description:    "全畫面冰凍，所有目標減速5秒",
		Cost:           300,
		MaxCharges:     3,
		Icon:           "❄️",
		Color:          "#87CEEB",
		ChargeRequired: 15, // 擊破 15 個目標自動充能一發（最容易充能）
		ChargePerKill:  1,
	},
	{
		Type:           WeaponTornado,
		Name:           "龍捲風砲",
		Description:    "全螢幕龍捲風，50%機率擊破所有目標，製造掃場爽感",
		Cost:           0,    // 只能透過充能獲得，不能購買
		MaxCharges:     2,    // 最多持有 2 發（稀有武器）
		Icon:           "🌪️",
		Color:          "#9B59B6",
		ChargeRequired: 50, // 擊破 50 個目標自動充能一發（最難充能，最強效果）
		ChargePerKill:  1,
	},
}

// PlayerWeaponState 玩家特殊武器狀態
type PlayerWeaponState struct {
	BombCharges    int `json:"bomb_charges"`
	LaserCharges   int `json:"laser_charges"`
	FreezeCharges  int `json:"freeze_charges"`
	TornadoCharges int `json:"tornado_charges"` // DAY-134

	// 自動充能進度（DAY-134）
	BombChargeProgress    int `json:"bomb_charge_progress"`
	LaserChargeProgress   int `json:"laser_charge_progress"`
	FreezeChargeProgress  int `json:"freeze_charge_progress"`
	TornadoChargeProgress int `json:"tornado_charge_progress"`
}

// ChargeResult 充能結果（DAY-134）
type ChargeResult struct {
	WeaponType    WeaponType
	NewProgress   int
	Required      int
	ChargeUnlocked bool // 是否剛充滿一發
	NewCharges    int
}

// Manager 特殊武器管理器
type Manager struct {
	mu     sync.RWMutex
	states map[string]*PlayerWeaponState // playerID -> state
}

// New 建立特殊武器管理器
func New() *Manager {
	return &Manager{
		states: make(map[string]*PlayerWeaponState),
	}
}

// GetOrCreate 取得或建立玩家武器狀態
func (m *Manager) GetOrCreate(playerID string) *PlayerWeaponState {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.states[playerID]; ok {
		return s
	}
	s := &PlayerWeaponState{}
	m.states[playerID] = s
	return s
}

// GetSnapshot 取得玩家武器狀態快照（thread-safe）
func (m *Manager) GetSnapshot(playerID string) PlayerWeaponState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if s, ok := m.states[playerID]; ok {
		return *s
	}
	return PlayerWeaponState{}
}

// BuyWeapon 購買特殊武器（回傳是否成功，以及扣除的金幣）
func (m *Manager) BuyWeapon(playerID string, wtype WeaponType, currentCoins int) (success bool, cost int) {
	def := getWeaponDef(wtype)
	if def == nil {
		return false, 0
	}
	// 龍捲風砲不能購買，只能充能
	if def.Cost == 0 {
		return false, 0
	}
	if currentCoins < def.Cost {
		return false, 0
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	s := m.getOrCreateLocked(playerID)
	charges := m.getChargesLocked(s, wtype)
	if charges >= def.MaxCharges {
		return false, 0 // 已達上限
	}
	m.setChargesLocked(s, wtype, charges+1)
	return true, def.Cost
}

// UseWeapon 使用特殊武器（回傳是否成功）
func (m *Manager) UseWeapon(playerID string, wtype WeaponType) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := m.getOrCreateLocked(playerID)
	charges := m.getChargesLocked(s, wtype)
	if charges <= 0 {
		return false
	}
	m.setChargesLocked(s, wtype, charges-1)
	return true
}

// RecordKill 記錄擊破，累積所有武器的充能進度（DAY-134）
// 回傳所有剛充滿的武器列表（可能同時充滿多個）
func (m *Manager) RecordKill(playerID string, multiplier float64) []ChargeResult {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := m.getOrCreateLocked(playerID)
	var results []ChargeResult

	// 高倍率目標給更多充能（≥10x 給 2 點，≥30x 給 3 點）
	chargePoints := 1
	if multiplier >= 30.0 {
		chargePoints = 3
	} else if multiplier >= 10.0 {
		chargePoints = 2
	}

	for _, def := range AvailableWeapons {
		if def.ChargeRequired <= 0 {
			continue
		}

		progress := m.getProgressLocked(s, def.Type)
		charges := m.getChargesLocked(s, def.Type)

		// 已達上限，不再累積
		if charges >= def.MaxCharges {
			continue
		}

		newProgress := progress + chargePoints*def.ChargePerKill
		chargeUnlocked := false

		if newProgress >= def.ChargeRequired {
			// 充滿一發
			newProgress = newProgress - def.ChargeRequired
			if newProgress < 0 {
				newProgress = 0
			}
			charges++
			if charges > def.MaxCharges {
				charges = def.MaxCharges
			}
			m.setChargesLocked(s, def.Type, charges)
			chargeUnlocked = true
		}

		m.setProgressLocked(s, def.Type, newProgress)

		results = append(results, ChargeResult{
			WeaponType:     def.Type,
			NewProgress:    newProgress,
			Required:       def.ChargeRequired,
			ChargeUnlocked: chargeUnlocked,
			NewCharges:     charges,
		})
	}

	return results
}

// RemovePlayer 移除玩家狀態
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.states, playerID)
}

// ---- 特殊武器效果計算 ----

// TargetPos 目標物位置（用於範圍計算）
type TargetPos struct {
	InstanceID string
	X, Y       float64
	Multiplier float64
}

// BombRadius 炸彈爆炸半徑（px）
const BombRadius = 200.0

// LaserYRange 雷射 Y 軸容差（px）
const LaserYRange = 60.0

// TornadoKillChance 龍捲風擊破機率（DAY-134）
const TornadoKillChance = 0.50

// CalcBombTargets 計算炸彈命中的目標（圓形範圍）
func CalcBombTargets(clickX, clickY float64, targets []TargetPos) []string {
	var hit []string
	for _, t := range targets {
		dx := t.X - clickX
		dy := t.Y - clickY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= BombRadius {
			hit = append(hit, t.InstanceID)
		}
	}
	return hit
}

// CalcLaserTargets 計算雷射命中的目標（水平穿透）
func CalcLaserTargets(clickY float64, targets []TargetPos) []string {
	var hit []string
	for _, t := range targets {
		if math.Abs(t.Y-clickY) <= LaserYRange {
			hit = append(hit, t.InstanceID)
		}
	}
	return hit
}

// CalcFreezeTargets 冰凍命中所有目標
func CalcFreezeTargets(targets []TargetPos) []string {
	ids := make([]string, 0, len(targets))
	for _, t := range targets {
		ids = append(ids, t.InstanceID)
	}
	return ids
}

// CalcTornadoTargets 龍捲風命中所有目標（全螢幕，DAY-134）
// 回傳所有目標 ID（由 handler 決定哪些被擊破，使用 TornadoKillChance）
func CalcTornadoTargets(targets []TargetPos) []string {
	ids := make([]string, 0, len(targets))
	for _, t := range targets {
		ids = append(ids, t.InstanceID)
	}
	return ids
}

// ---- 內部輔助函數 ----

func getWeaponDef(wtype WeaponType) *WeaponDef {
	for i := range AvailableWeapons {
		if AvailableWeapons[i].Type == wtype {
			return &AvailableWeapons[i]
		}
	}
	return nil
}

func (m *Manager) getOrCreateLocked(playerID string) *PlayerWeaponState {
	if s, ok := m.states[playerID]; ok {
		return s
	}
	s := &PlayerWeaponState{}
	m.states[playerID] = s
	return s
}

func (m *Manager) getChargesLocked(s *PlayerWeaponState, wtype WeaponType) int {
	switch wtype {
	case WeaponBomb:
		return s.BombCharges
	case WeaponLaser:
		return s.LaserCharges
	case WeaponFreeze:
		return s.FreezeCharges
	case WeaponTornado:
		return s.TornadoCharges
	}
	return 0
}

func (m *Manager) setChargesLocked(s *PlayerWeaponState, wtype WeaponType, v int) {
	switch wtype {
	case WeaponBomb:
		s.BombCharges = v
	case WeaponLaser:
		s.LaserCharges = v
	case WeaponFreeze:
		s.FreezeCharges = v
	case WeaponTornado:
		s.TornadoCharges = v
	}
}

func (m *Manager) getProgressLocked(s *PlayerWeaponState, wtype WeaponType) int {
	switch wtype {
	case WeaponBomb:
		return s.BombChargeProgress
	case WeaponLaser:
		return s.LaserChargeProgress
	case WeaponFreeze:
		return s.FreezeChargeProgress
	case WeaponTornado:
		return s.TornadoChargeProgress
	}
	return 0
}

func (m *Manager) setProgressLocked(s *PlayerWeaponState, wtype WeaponType, v int) {
	switch wtype {
	case WeaponBomb:
		s.BombChargeProgress = v
	case WeaponLaser:
		s.LaserChargeProgress = v
	case WeaponFreeze:
		s.FreezeChargeProgress = v
	case WeaponTornado:
		s.TornadoChargeProgress = v
	}
}

// AddCharge 直接增加充能（開箱獎勵用，不扣金幣）
func (m *Manager) AddCharge(playerID string, wtype WeaponType) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := m.getOrCreateLocked(playerID)
	def := getWeaponDef(wtype)
	if def == nil {
		return
	}
	charges := m.getChargesLocked(s, wtype)
	if charges < def.MaxCharges {
		m.setChargesLocked(s, wtype, charges+1)
	}
}

// LoadState 從持久化資料恢復特殊武器充能數（DAY-100）
func (m *Manager) LoadState(playerID string, bomb int, laser int, freeze int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.states[playerID]
	if !ok {
		s = &PlayerWeaponState{}
		m.states[playerID] = s
	}
	s.BombCharges = clampCharges(bomb, 3)
	s.LaserCharges = clampCharges(laser, 3)
	s.FreezeCharges = clampCharges(freeze, 3)
	// TornadoCharges 不持久化（每次登入重新充能，保持稀有感）
}

// clampCharges 限制充能數在 0-max 範圍內
func clampCharges(v, max int) int {
	if v < 0 {
		return 0
	}
	if v > max {
		return max
	}
	return v
}
