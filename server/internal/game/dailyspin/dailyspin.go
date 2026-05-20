// Package dailyspin 每日簽到轉盤系統（DAY-092）
// 每天登入可以免費轉一次，連續 7 天可轉超級轉盤
// 業界依據：iGaming 2026 最熱門留存機制，每日驚喜感提升留存率 35%+
package dailyspin

import (
	"math/rand"
	"sync"
	"time"
)

// RewardType 獎勵類型
type RewardType string

const (
	RewardCoins         RewardType = "coins"          // 金幣
	RewardBombCharge    RewardType = "bomb_charge"     // 炸彈充能
	RewardLaserCharge   RewardType = "laser_charge"    // 雷射充能
	RewardFreezeCharge  RewardType = "freeze_charge"   // 冰凍充能
	RewardMysteryBox    RewardType = "mystery_box"     // 神秘寶箱（普通）
	RewardSeasonPoints  RewardType = "season_points"   // 賽季積分
	RewardJackpotTicket RewardType = "jackpot_ticket"  // Jackpot 券
	RewardMultBonus     RewardType = "mult_bonus"      // 下次攻擊倍率加成
)

// Slot 轉盤格子
type Slot struct {
	ID          int        `json:"id"`
	Type        RewardType `json:"type"`
	Amount      int        `json:"amount"`
	Label       string     `json:"label"`
	Icon        string     `json:"icon"`
	Color       string     `json:"color"`
	Weight      int        `json:"-"` // 加權機率（不回傳給 Client）
	IsSuper     bool       `json:"is_super"` // 是否是超級轉盤專屬格子
}

// SpinResult 轉盤結果
type SpinResult struct {
	SlotIndex   int        `json:"slot_index"`
	Slot        Slot       `json:"slot"`
	IsSuper     bool       `json:"is_super"`
	LoginStreak int        `json:"login_streak"`
	NextSpinAt  int64      `json:"next_spin_at"` // Unix ms，下次可轉時間
}

// PlayerState 玩家轉盤狀態
type PlayerState struct {
	LastSpinDate  string // UTC+8 日期（"2006-01-02"）
	LoginStreak   int    // 連續登入天數（用於超級轉盤判斷）
	TotalSpins    int    // 累計轉盤次數
}

// Manager 每日轉盤管理器
type Manager struct {
	mu      sync.RWMutex
	players map[string]*PlayerState
	rng     *rand.Rand

	// 普通轉盤格子（8格）
	normalSlots []Slot
	// 超級轉盤格子（8格，連續 7 天解鎖，獎勵翻倍）
	superSlots []Slot
}

// NewManager 建立每日轉盤管理器
func NewManager() *Manager {
	m := &Manager{
		players: make(map[string]*PlayerState),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	m.initSlots()
	return m
}

func (m *Manager) initSlots() {
	// 普通轉盤（8格）
	m.normalSlots = []Slot{
		{ID: 0, Type: RewardCoins, Amount: 500, Label: "500 金幣", Icon: "🪙", Color: "#FFD700", Weight: 30},
		{ID: 1, Type: RewardCoins, Amount: 1000, Label: "1000 金幣", Icon: "🪙", Color: "#FFD700", Weight: 25},
		{ID: 2, Type: RewardBombCharge, Amount: 1, Label: "炸彈 ×1", Icon: "💣", Color: "#FF6B35", Weight: 15},
		{ID: 3, Type: RewardMysteryBox, Amount: 1, Label: "普通寶箱 ×1", Icon: "📦", Color: "#4CAF50", Weight: 12},
		{ID: 4, Type: RewardSeasonPoints, Amount: 50, Label: "賽季積分 +50", Icon: "⭐", Color: "#2196F3", Weight: 10},
		{ID: 5, Type: RewardLaserCharge, Amount: 1, Label: "雷射 ×1", Icon: "⚡", Color: "#9C27B0", Weight: 4},
		{ID: 6, Type: RewardJackpotTicket, Amount: 1, Label: "Jackpot 券 ×1", Icon: "🎰", Color: "#FF9800", Weight: 3},
		{ID: 7, Type: RewardCoins, Amount: 5000, Label: "5000 金幣", Icon: "💰", Color: "#FF5722", Weight: 1},
	}

	// 超級轉盤（連續 7 天，獎勵翻倍）
	m.superSlots = []Slot{
		{ID: 0, Type: RewardCoins, Amount: 2000, Label: "2000 金幣", Icon: "🪙", Color: "#FFD700", Weight: 25, IsSuper: true},
		{ID: 1, Type: RewardCoins, Amount: 5000, Label: "5000 金幣", Icon: "🪙", Color: "#FFD700", Weight: 20, IsSuper: true},
		{ID: 2, Type: RewardBombCharge, Amount: 2, Label: "炸彈 ×2", Icon: "💣", Color: "#FF6B35", Weight: 15, IsSuper: true},
		{ID: 3, Type: RewardMysteryBox, Amount: 1, Label: "稀有寶箱 ×1", Icon: "📦", Color: "#2196F3", Weight: 12, IsSuper: true},
		{ID: 4, Type: RewardSeasonPoints, Amount: 200, Label: "賽季積分 +200", Icon: "⭐", Color: "#2196F3", Weight: 10, IsSuper: true},
		{ID: 5, Type: RewardMultBonus, Amount: 3, Label: "下次攻擊 ×3.0", Icon: "🔥", Color: "#FF9800", Weight: 8, IsSuper: true},
		{ID: 6, Type: RewardJackpotTicket, Amount: 5, Label: "Jackpot 券 ×5", Icon: "🎰", Color: "#FF9800", Weight: 5, IsSuper: true},
		{ID: 7, Type: RewardCoins, Amount: 20000, Label: "20000 金幣", Icon: "💰", Color: "#FF5722", Weight: 5, IsSuper: true},
	}
}

// todayUTC8 取得 UTC+8 今日日期字串
func todayUTC8() string {
	loc := time.FixedZone("UTC+8", 8*60*60)
	return time.Now().In(loc).Format("2006-01-02")
}

// tomorrowUTC8Midnight 取得 UTC+8 明日 00:00 的 Unix ms
func tomorrowUTC8Midnight() int64 {
	loc := time.FixedZone("UTC+8", 8*60*60)
	now := time.Now().In(loc)
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)
	return tomorrow.UnixMilli()
}

// CanSpin 檢查玩家今天是否可以轉盤
func (m *Manager) CanSpin(playerID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	state, ok := m.players[playerID]
	if !ok {
		return true
	}
	return state.LastSpinDate != todayUTC8()
}

// IsSuper 檢查玩家是否可以轉超級轉盤（連續 7 天）
func (m *Manager) IsSuper(playerID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	state, ok := m.players[playerID]
	if !ok {
		return false
	}
	return state.LoginStreak >= 7
}

// GetSnapshot 取得玩家轉盤狀態快照
func (m *Manager) GetSnapshot(playerID string) map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	canSpin := true
	loginStreak := 0
	totalSpins := 0
	isSuper := false

	if state, ok := m.players[playerID]; ok {
		canSpin = state.LastSpinDate != todayUTC8()
		loginStreak = state.LoginStreak
		totalSpins = state.TotalSpins
		isSuper = state.LoginStreak >= 7
	}

	nextSpinAt := int64(0)
	if !canSpin {
		nextSpinAt = tomorrowUTC8Midnight()
	}

	return map[string]interface{}{
		"can_spin":     canSpin,
		"is_super":     isSuper,
		"login_streak": loginStreak,
		"total_spins":  totalSpins,
		"next_spin_at": nextSpinAt,
		"normal_slots": m.normalSlots,
		"super_slots":  m.superSlots,
	}
}

// Spin 執行轉盤（回傳結果，若今天已轉過則回傳 nil）
func (m *Manager) Spin(playerID string, loginStreak int) *SpinResult {
	m.mu.Lock()
	defer m.mu.Unlock()

	today := todayUTC8()

	// 取得或建立玩家狀態
	state, ok := m.players[playerID]
	if !ok {
		state = &PlayerState{}
		m.players[playerID] = state
	}

	// 今天已轉過
	if state.LastSpinDate == today {
		return nil
	}

	// 更新狀態
	state.LastSpinDate = today
	state.LoginStreak = loginStreak
	state.TotalSpins++

	// 決定使用普通還是超級轉盤
	isSuper := loginStreak >= 7
	slots := m.normalSlots
	if isSuper {
		slots = m.superSlots
	}

	// 加權隨機選擇格子
	totalWeight := 0
	for _, s := range slots {
		totalWeight += s.Weight
	}
	r := m.rng.Intn(totalWeight)
	cumulative := 0
	selectedIdx := 0
	for i, s := range slots {
		cumulative += s.Weight
		if r < cumulative {
			selectedIdx = i
			break
		}
	}

	return &SpinResult{
		SlotIndex:   selectedIdx,
		Slot:        slots[selectedIdx],
		IsSuper:     isSuper,
		LoginStreak: loginStreak,
		NextSpinAt:  tomorrowUTC8Midnight(),
	}
}

// GetSlots 取得轉盤格子定義（用於 Client 顯示）
func (m *Manager) GetSlots(isSuper bool) []Slot {
	if isSuper {
		return m.superSlots
	}
	return m.normalSlots
}
