// Package speedrace 全服競速獵殺系統（DAY-136）
// 業界依據：soup.io 2025「PvP modes where every shot counts」
// 隨機選定高價值目標，全服搶先擊破，第一名獲得 3x 獎勵加成
package speedrace

import (
	"sync"
	"time"
)

// RaceConfig 競速獵殺設定
type RaceConfig struct {
	Duration      float64 // 競速持續秒數（預設 30 秒）
	BonusMult     float64 // 第一名獎勵倍率（預設 3.0x）
	SecondMult    float64 // 第二名獎勵倍率（預設 1.5x）
	ThirdMult     float64 // 第三名獎勵倍率（預設 1.2x）
	CooldownSecs  float64 // 冷卻時間（預設 90 秒）
	MinMultiplier float64 // 觸發競速的最低目標倍率（預設 10x）
}

// DefaultConfig 預設設定
func DefaultConfig() RaceConfig {
	return RaceConfig{
		Duration:      30.0,
		BonusMult:     3.0,
		SecondMult:    1.5,
		ThirdMult:     1.2,
		CooldownSecs:  90.0,
		MinMultiplier: 10.0,
	}
}

// RaceResult 競速結果（單一玩家）
type RaceResult struct {
	PlayerID    string
	DisplayName string
	Rank        int     // 1/2/3
	BonusMult   float64 // 獎勵倍率
	KilledAt    time.Time
}

// RaceSession 競速 session
type RaceSession struct {
	TargetInstanceID string
	TargetDefID      string
	TargetName       string
	TargetMult       float64
	StartAt          time.Time
	EndAt            time.Time
	Results          []RaceResult // 依擊破順序排列
	IsActive         bool
}

// Manager 競速獵殺管理器
type Manager struct {
	mu       sync.RWMutex
	config   RaceConfig
	session  *RaceSession
	lastEndAt time.Time // 上次競速結束時間（冷卻用）
}

// New 建立競速獵殺管理器
func New(cfg RaceConfig) *Manager {
	return &Manager{
		config: cfg,
	}
}

// NewDefault 使用預設設定建立管理器
func NewDefault() *Manager {
	return New(DefaultConfig())
}

// CanStart 是否可以開始新競速（冷卻檢查）
func (m *Manager) CanStart() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.session != nil && m.session.IsActive {
		return false // 已有進行中的競速
	}
	if m.lastEndAt.IsZero() {
		return true
	}
	return time.Since(m.lastEndAt).Seconds() >= m.config.CooldownSecs
}

// StartRace 開始競速（由 spawnTarget 或 triggerSpecialEvent 呼叫）
// 回傳 session（nil 表示無法開始）
func (m *Manager) StartRace(instanceID, defID, name string, mult float64) *RaceSession {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session != nil && m.session.IsActive {
		return nil
	}
	if !m.lastEndAt.IsZero() && time.Since(m.lastEndAt).Seconds() < m.config.CooldownSecs {
		return nil
	}
	if mult < m.config.MinMultiplier {
		return nil
	}

	now := time.Now()
	m.session = &RaceSession{
		TargetInstanceID: instanceID,
		TargetDefID:      defID,
		TargetName:       name,
		TargetMult:       mult,
		StartAt:          now,
		EndAt:            now.Add(time.Duration(m.config.Duration * float64(time.Second))),
		Results:          make([]RaceResult, 0, 3),
		IsActive:         true,
	}
	return m.session
}

// RecordKill 記錄玩家擊破競速目標
// 回傳：(rank, bonusMult, isRaceTarget)
// isRaceTarget=false 表示此目標不是競速目標
func (m *Manager) RecordKill(instanceID, playerID, displayName string) (rank int, bonusMult float64, isRaceTarget bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session == nil || !m.session.IsActive {
		return 0, 1.0, false
	}
	if m.session.TargetInstanceID != instanceID {
		return 0, 1.0, false
	}

	// 已有 3 名以上，不再記錄
	if len(m.session.Results) >= 3 {
		return 0, 1.0, true
	}

	// 計算名次
	rank = len(m.session.Results) + 1
	var mult float64
	switch rank {
	case 1:
		mult = m.config.BonusMult
	case 2:
		mult = m.config.SecondMult
	case 3:
		mult = m.config.ThirdMult
	default:
		mult = 1.0
	}

	m.session.Results = append(m.session.Results, RaceResult{
		PlayerID:    playerID,
		DisplayName: displayName,
		Rank:        rank,
		BonusMult:   mult,
		KilledAt:    time.Now(),
	})

	// 第一名擊破後結束競速
	if rank == 1 {
		m.session.IsActive = false
		m.lastEndAt = time.Now()
	}

	return rank, mult, true
}

// CheckExpiry 檢查競速是否超時
// 回傳 true 表示剛剛超時（需要廣播取消）
func (m *Manager) CheckExpiry() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session == nil || !m.session.IsActive {
		return false
	}
	if time.Now().After(m.session.EndAt) {
		m.session.IsActive = false
		m.lastEndAt = time.Now()
		return true
	}
	return false
}

// CancelRace 取消競速（目標消失時呼叫）
func (m *Manager) CancelRace(instanceID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session == nil || !m.session.IsActive {
		return false
	}
	if m.session.TargetInstanceID != instanceID {
		return false
	}
	m.session.IsActive = false
	m.lastEndAt = time.Now()
	return true
}

// GetSnapshot 取得當前競速快照
func (m *Manager) GetSnapshot() *RaceSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.session == nil {
		return &RaceSnapshot{IsActive: false}
	}

	results := make([]RaceResultSnap, len(m.session.Results))
	for i, r := range m.session.Results {
		results[i] = RaceResultSnap{
			PlayerID:    r.PlayerID,
			DisplayName: r.DisplayName,
			Rank:        r.Rank,
			BonusMult:   r.BonusMult,
		}
	}

	secsLeft := 0.0
	if m.session.IsActive {
		secsLeft = time.Until(m.session.EndAt).Seconds()
		if secsLeft < 0 {
			secsLeft = 0
		}
	}

	return &RaceSnapshot{
		IsActive:         m.session.IsActive,
		TargetInstanceID: m.session.TargetInstanceID,
		TargetDefID:      m.session.TargetDefID,
		TargetName:       m.session.TargetName,
		TargetMult:       m.session.TargetMult,
		SecondsLeft:      secsLeft,
		Results:          results,
		BonusMult:        m.config.BonusMult,
	}
}

// IsRaceTarget 快速檢查某個目標是否是當前競速目標
func (m *Manager) IsRaceTarget(instanceID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.session != nil && m.session.IsActive && m.session.TargetInstanceID == instanceID
}

// RaceSnapshot 競速快照（用於廣播）
type RaceSnapshot struct {
	IsActive         bool
	TargetInstanceID string
	TargetDefID      string
	TargetName       string
	TargetMult       float64
	SecondsLeft      float64
	Results          []RaceResultSnap
	BonusMult        float64 // 第一名倍率（顯示用）
}

// RaceResultSnap 競速結果快照
type RaceResultSnap struct {
	PlayerID    string
	DisplayName string
	Rank        int
	BonusMult   float64
}
