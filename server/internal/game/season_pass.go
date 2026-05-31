package game

import (
	"sync"
	"time"
)

// SeasonPassTier 賽季通行證等級
type SeasonPassTier struct {
	Level       int    `json:"level"`
	Name        string `json:"name"`
	RequiredXP  int    `json:"required_xp"`
	FreeReward  int    `json:"free_reward"`  // 免費獎勵（金幣）
	PremiumReward int  `json:"premium_reward"` // 付費獎勵（任務幣）
	BadgeName   string `json:"badge_name"`
}

// SeasonPassState 玩家賽季通行證狀態
type SeasonPassState struct {
	PlayerID    string    `json:"player_id"`
	CurrentXP   int       `json:"current_xp"`
	CurrentLevel int      `json:"current_level"`
	IsPremium   bool      `json:"is_premium"`
	ClaimedFree []int     `json:"claimed_free"`    // 已領取的免費等級
	ClaimedPrem []int     `json:"claimed_premium"` // 已領取的付費等級
	LastUpdated time.Time `json:"last_updated"`
}

// SeasonPassManager 賽季通行證管理器
type SeasonPassManager struct {
	mu       sync.RWMutex
	states   map[string]*SeasonPassState
	tiers    []SeasonPassTier
	seasonID string
	seasonEnd time.Time
}

// 賽季通行證等級定義（10個等級，每季 30 天）
var defaultSeasonTiers = []SeasonPassTier{
	{Level: 1,  Name: "新手探索者",   RequiredXP: 0,    FreeReward: 100,   PremiumReward: 50,   BadgeName: "🌱"},
	{Level: 2,  Name: "初級獵人",     RequiredXP: 100,  FreeReward: 200,   PremiumReward: 100,  BadgeName: "🗡️"},
	{Level: 3,  Name: "中級戰士",     RequiredXP: 300,  FreeReward: 300,   PremiumReward: 150,  BadgeName: "⚔️"},
	{Level: 4,  Name: "高級勇者",     RequiredXP: 600,  FreeReward: 500,   PremiumReward: 200,  BadgeName: "🛡️"},
	{Level: 5,  Name: "精英鬥士",     RequiredXP: 1000, FreeReward: 800,   PremiumReward: 300,  BadgeName: "🏆"},
	{Level: 6,  Name: "傳說英雄",     RequiredXP: 1500, FreeReward: 1200,  PremiumReward: 500,  BadgeName: "⭐"},
	{Level: 7,  Name: "神話戰神",     RequiredXP: 2100, FreeReward: 1800,  PremiumReward: 700,  BadgeName: "🌟"},
	{Level: 8,  Name: "宇宙霸主",     RequiredXP: 2800, FreeReward: 2500,  PremiumReward: 1000, BadgeName: "💫"},
	{Level: 9,  Name: "時空征服者",   RequiredXP: 3600, FreeReward: 3500,  PremiumReward: 1500, BadgeName: "🌌"},
	{Level: 10, Name: "終極大師",     RequiredXP: 4500, FreeReward: 5000,  PremiumReward: 2000, BadgeName: "👑"},
}

// XP 獲取規則
const (
	XPPerKill        = 1   // 每次擊破 +1 XP
	XPPerBossKill    = 10  // 擊破 BOSS +10 XP
	XPPerBonus       = 5   // 完成 Bonus +5 XP
	XPPerCombo5      = 3   // 5連擊 +3 XP
	XPPerCombo10     = 8   // 10連擊 +8 XP
	XPPerDailyQuest  = 20  // 完成每日任務 +20 XP
	XPPerWeeklyChallenge = 50 // 完成每週挑戰 +50 XP
)

// NewSeasonPassManager 建立賽季通行證管理器
func NewSeasonPassManager() *SeasonPassManager {
	// 計算本季結束時間（每月1日 UTC+8 00:00 重置）
	now := time.Now().In(time.FixedZone("UTC+8", 8*60*60))
	// 下個月1日
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	
	return &SeasonPassManager{
		states:    make(map[string]*SeasonPassState),
		tiers:     defaultSeasonTiers,
		seasonID:  now.Format("2006-01"),
		seasonEnd: nextMonth,
	}
}

// GetOrCreateState 取得或建立玩家狀態
func (m *SeasonPassManager) GetOrCreateState(playerID string) *SeasonPassState {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if state, ok := m.states[playerID]; ok {
		return state
	}
	
	state := &SeasonPassState{
		PlayerID:     playerID,
		CurrentXP:    0,
		CurrentLevel: 1,
		IsPremium:    false,
		ClaimedFree:  []int{},
		ClaimedPrem:  []int{},
		LastUpdated:  time.Now(),
	}
	m.states[playerID] = state
	return state
}

// AddXP 增加 XP 並更新等級
func (m *SeasonPassManager) AddXP(playerID string, xp int) (levelUp bool, newLevel int, newXP int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	state, ok := m.states[playerID]
	if !ok {
		state = &SeasonPassState{
			PlayerID:     playerID,
			CurrentXP:    0,
			CurrentLevel: 1,
			IsPremium:    false,
			ClaimedFree:  []int{},
			ClaimedPrem:  []int{},
			LastUpdated:  time.Now(),
		}
		m.states[playerID] = state
	}
	
	oldLevel := state.CurrentLevel
	state.CurrentXP += xp
	state.LastUpdated = time.Now()
	
	// 計算新等級
	newLvl := 1
	for _, tier := range m.tiers {
		if state.CurrentXP >= tier.RequiredXP {
			newLvl = tier.Level
		}
	}
	state.CurrentLevel = newLvl
	
	return newLvl > oldLevel, newLvl, state.CurrentXP
}

// GetSnapshot 取得玩家賽季狀態快照
func (m *SeasonPassManager) GetSnapshot(playerID string) map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	state, ok := m.states[playerID]
	if !ok {
		return map[string]interface{}{
			"current_xp":    0,
			"current_level": 1,
			"is_premium":    false,
			"season_id":     m.seasonID,
			"season_end":    m.seasonEnd.Format(time.RFC3339),
			"days_left":     int(time.Until(m.seasonEnd).Hours() / 24),
			"tiers":         m.tiers,
		}
	}
	
	// 計算下一等級所需 XP
	nextLevelXP := -1
	if state.CurrentLevel < len(m.tiers) {
		nextLevelXP = m.tiers[state.CurrentLevel].RequiredXP
	}
	
	return map[string]interface{}{
		"current_xp":     state.CurrentXP,
		"current_level":  state.CurrentLevel,
		"next_level_xp":  nextLevelXP,
		"is_premium":     state.IsPremium,
		"claimed_free":   state.ClaimedFree,
		"claimed_premium": state.ClaimedPrem,
		"season_id":      m.seasonID,
		"season_end":     m.seasonEnd.Format(time.RFC3339),
		"days_left":      int(time.Until(m.seasonEnd).Hours() / 24),
		"tiers":          m.tiers,
	}
}

// GetSeasonID 取得當前賽季 ID
func (m *SeasonPassManager) GetSeasonID() string {
	return m.seasonID
}

// GetSeasonEnd 取得賽季結束時間
func (m *SeasonPassManager) GetSeasonEnd() time.Time {
	return m.seasonEnd
}

// GetPlayerCount 取得參與玩家數
func (m *SeasonPassManager) GetPlayerCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.states)
}
