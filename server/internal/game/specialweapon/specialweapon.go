// Package specialweapon 特殊武器系統（DAY-089）
// 業界依據：Fish Road 2026 有 8 tier 武器系統，炸彈/雷射是標配特殊武器
package specialweapon

import (
	"math"
	"sync"
)

// WeaponType 特殊武器類型
type WeaponType string

const (
	WeaponBomb   WeaponType = "bomb"   // 炸彈砲：範圍爆炸，命中半徑 200px 內所有目標
	WeaponLaser  WeaponType = "laser"  // 雷射砲：穿透，命中一條線上所有目標（Y 軸 ±60px）
	WeaponFreeze WeaponType = "freeze" // 冰凍砲：全畫面減速所有目標 5 秒
)

// WeaponDef 特殊武器定義
type WeaponDef struct {
	Type        WeaponType `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Cost        int        `json:"cost"`        // 購買一發的金幣費用
	MaxCharges  int        `json:"max_charges"` // 最多持有發數
	Icon        string     `json:"icon"`        // 圖示（emoji）
	Color       string     `json:"color"`       // 顏色（hex）
}

// AvailableWeapons 所有特殊武器定義
var AvailableWeapons = []WeaponDef{
	{
		Type:        WeaponBomb,
		Name:        "炸彈砲",
		Description: "範圍爆炸，命中半徑200px內所有目標",
		Cost:        500,
		MaxCharges:  3,
		Icon:        "💣",
		Color:       "#FF6B35",
	},
	{
		Type:        WeaponLaser,
		Name:        "雷射砲",
		Description: "穿透射擊，命中同一水平線上所有目標",
		Cost:        800,
		MaxCharges:  3,
		Icon:        "⚡",
		Color:       "#00FFFF",
	},
	{
		Type:        WeaponFreeze,
		Name:        "冰凍砲",
		Description: "全畫面冰凍，所有目標減速5秒",
		Cost:        300,
		MaxCharges:  3,
		Icon:        "❄️",
		Color:       "#87CEEB",
	},
}

// PlayerWeaponState 玩家特殊武器狀態
type PlayerWeaponState struct {
	BombCharges   int `json:"bomb_charges"`
	LaserCharges  int `json:"laser_charges"`
	FreezeCharges int `json:"freeze_charges"`
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
	m.states[playerID] = &PlayerWeaponState{
		BombCharges:   clampCharges(bomb, 3),
		LaserCharges:  clampCharges(laser, 3),
		FreezeCharges: clampCharges(freeze, 3),
	}
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
