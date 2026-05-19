// Package mission 每日任務系統（DAY-037）
// 提供每日任務定義、進度追蹤、完成獎勵
package mission

import (
	"sync"
	"time"
)

// MissionType 任務類型
type MissionType string

const (
	MissionKillTargets  MissionType = "kill_targets"   // 擊破 N 個目標
	MissionKillBoss     MissionType = "kill_boss"       // 擊敗 BOSS N 次
	MissionPlayBonus    MissionType = "play_bonus"      // 完成 Bonus Game N 次
	MissionEarnCoins    MissionType = "earn_coins"      // 累積獲得 N 金幣
	MissionKillHighMult MissionType = "kill_high_mult"  // 擊破高倍率目標（30x+）N 個
	MissionCombo        MissionType = "combo"           // 達成 N 連擊
)

// Mission 任務定義
type Mission struct {
	ID          string      `json:"id"`
	Type        MissionType `json:"type"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Target      int         `json:"target"`      // 目標數量
	Reward      int         `json:"reward"`      // 完成獎勵（金幣）
	Icon        string      `json:"icon"`        // 顯示圖示
}

// PlayerProgress 玩家任務進度
type PlayerProgress struct {
	MissionID   string    `json:"mission_id"`
	Current     int       `json:"current"`     // 當前進度
	Target      int       `json:"target"`      // 目標數量
	Completed   bool      `json:"completed"`   // 是否已完成
	RewardClaimed bool    `json:"reward_claimed"` // 是否已領取獎勵
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// DailyMissions 每日任務集合（每天固定 6 個任務）
var DailyMissions = []Mission{
	{
		ID:          "daily_kill_10",
		Type:        MissionKillTargets,
		Name:        "討伐大作戰",
		Description: "擊破 10 個目標",
		Target:      10,
		Reward:      500,
		Icon:        "⚔️",
	},
	{
		ID:          "daily_kill_boss",
		Type:        MissionKillBoss,
		Name:        "那個孩子的剋星",
		Description: "擊敗 BOSS 1 次",
		Target:      1,
		Reward:      2000,
		Icon:        "👹",
	},
	{
		ID:          "daily_bonus",
		Type:        MissionPlayBonus,
		Name:        "瘋狂拔草達人",
		Description: "完成 Bonus Game 2 次",
		Target:      2,
		Reward:      1000,
		Icon:        "🌿",
	},
	{
		ID:          "daily_earn_5000",
		Type:        MissionEarnCoins,
		Name:        "金幣收集家",
		Description: "累積獲得 5000 金幣",
		Target:      5000,
		Reward:      800,
		Icon:        "🪙",
	},
	{
		ID:          "daily_high_mult",
		Type:        MissionKillHighMult,
		Name:        "高手中的高手",
		Description: "擊破 3 個 30x+ 高倍率目標",
		Target:      3,
		Reward:      1500,
		Icon:        "⭐",
	},
	{
		ID:          "daily_combo_5",
		Type:        MissionCombo,
		Name:        "連擊達人",
		Description: "達成 5 連擊",
		Target:      5,
		Reward:      1200,
		Icon:        "🔥",
	},
}

// Manager 任務管理器（per-player 進度追蹤）
type Manager struct {
	mu       sync.RWMutex
	progress map[string]map[string]*PlayerProgress // playerID → missionID → progress
	resetAt  time.Time                              // 下次重置時間（每日 00:00）
}

// NewManager 建立任務管理器
func NewManager() *Manager {
	m := &Manager{
		progress: make(map[string]map[string]*PlayerProgress),
		resetAt:  nextMidnight(),
	}
	return m
}

// nextMidnight 計算下一個午夜時間（UTC+8，台灣/亞洲標準時間）
// 業界標準：每日任務以 UTC+8 00:00 為重置基準，確保所有玩家同步重置
func nextMidnight() time.Time {
	loc := time.FixedZone("UTC+8", 8*60*60)
	now := time.Now().In(loc)
	return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)
}

// GetOrInitProgress 取得或初始化玩家進度
func (m *Manager) GetOrInitProgress(playerID string) map[string]*PlayerProgress {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 檢查是否需要重置（每日）
	if time.Now().After(m.resetAt) {
		m.progress = make(map[string]map[string]*PlayerProgress)
		m.resetAt = nextMidnight()
	}

	if _, ok := m.progress[playerID]; !ok {
		m.progress[playerID] = make(map[string]*PlayerProgress)
		for _, mission := range DailyMissions {
			m.progress[playerID][mission.ID] = &PlayerProgress{
				MissionID: mission.ID,
				Current:   0,
				Target:    mission.Target,
				Completed: false,
			}
		}
	}
	return m.progress[playerID]
}

// UpdateProgress 更新任務進度，回傳新完成的任務列表
func (m *Manager) UpdateProgress(playerID string, mType MissionType, amount int) []Mission {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 確保玩家進度已初始化
	if _, ok := m.progress[playerID]; !ok {
		m.mu.Unlock()
		m.GetOrInitProgress(playerID)
		m.mu.Lock()
	}

	playerProgress := m.progress[playerID]
	var newlyCompleted []Mission

	for _, mission := range DailyMissions {
		if mission.Type != mType {
			continue
		}
		prog, ok := playerProgress[mission.ID]
		if !ok || prog.Completed {
			continue
		}

		prog.Current += amount
		if prog.Current >= prog.Target {
			prog.Current = prog.Target
			prog.Completed = true
			prog.CompletedAt = time.Now()
			newlyCompleted = append(newlyCompleted, mission)
		}
	}

	return newlyCompleted
}

// ClaimReward 領取任務獎勵，回傳獎勵金幣數（0 = 無法領取）
func (m *Manager) ClaimReward(playerID, missionID string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	playerProgress, ok := m.progress[playerID]
	if !ok {
		return 0
	}
	prog, ok := playerProgress[missionID]
	if !ok || !prog.Completed || prog.RewardClaimed {
		return 0
	}

	// 找到對應任務的獎勵
	for _, mission := range DailyMissions {
		if mission.ID == missionID {
			prog.RewardClaimed = true
			return mission.Reward
		}
	}
	return 0
}

// GetPlayerMissions 取得玩家所有任務進度（含任務定義）
func (m *Manager) GetPlayerMissions(playerID string) []MissionStatus {
	progress := m.GetOrInitProgress(playerID)

	m.mu.RLock()
	defer m.mu.RUnlock()

	statuses := make([]MissionStatus, 0, len(DailyMissions))
	for _, mission := range DailyMissions {
		prog := progress[mission.ID]
		status := MissionStatus{
			Mission:  mission,
			Progress: *prog,
		}
		statuses = append(statuses, status)
	}
	return statuses
}

// MissionStatus 任務狀態（任務定義 + 玩家進度）
type MissionStatus struct {
	Mission  Mission        `json:"mission"`
	Progress PlayerProgress `json:"progress"`
}

// ResetAt 取得下次重置時間
func (m *Manager) ResetAt() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.resetAt
}
