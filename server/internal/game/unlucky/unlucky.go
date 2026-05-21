// Package unlucky 失敗補償系統（DAY-135）
// 業界依據：Funrize 2026 的「Unlucky Bonus」
// 連續花費超過一定金額但獲得低回報時，自動給予補償獎勵
// 防止玩家因為「運氣太差」而離開，是 2026 年業界最新的留存機制
package unlucky

import (
	"sync"
	"time"
)

// UnluckyConfig 失敗補償設定
type UnluckyConfig struct {
	// 觸發條件：連續 N 次射擊的花費超過回報的比例
	TrackingShots    int     // 追蹤最近 N 次射擊（預設 30）
	SpendThreshold   float64 // 花費/回報比例門檻（預設 3.0 = 花了 3 倍才回收 1 倍）
	MinSpend         int     // 最低花費門檻（避免低投注觸發，預設 betCost × 20）
	CooldownSecs     int     // 補償後冷卻秒數（預設 120）
	// 補償獎勵
	BaseRewardMult   float64 // 補償獎勵 = 淨虧損 × BaseRewardMult（預設 0.3，補償 30% 虧損）
	MaxRewardMult    float64 // 最高補償倍率（預設 0.5）
	MinReward        int     // 最低補償金額（預設 100）
}

// DefaultConfig 預設設定
var DefaultConfig = UnluckyConfig{
	TrackingShots:  30,
	SpendThreshold: 3.0,
	MinSpend:       200,
	CooldownSecs:   120,
	BaseRewardMult: 0.3,
	MaxRewardMult:  0.5,
	MinReward:      100,
}

// ShotRecord 單次射擊記錄
type ShotRecord struct {
	Spend  int // 花費（betCost）
	Reward int // 獲得（0 = 未擊破）
}

// PlayerState 玩家失敗補償狀態
type PlayerState struct {
	Shots        []ShotRecord // 最近 N 次射擊記錄（環形緩衝）
	ShotIdx      int          // 環形緩衝索引
	TotalSpend   int          // 追蹤期間總花費
	TotalReward  int          // 追蹤期間總回報
	LastBonusAt  time.Time    // 上次補償時間
	BonusCount   int          // 累計補償次數（本次登入）
}

// Manager 失敗補償管理器
type Manager struct {
	mu     sync.RWMutex
	states map[string]*PlayerState
	cfg    UnluckyConfig
}

// New 建立失敗補償管理器
func New(cfg UnluckyConfig) *Manager {
	return &Manager{
		states: make(map[string]*PlayerState),
		cfg:    cfg,
	}
}

// NewDefault 使用預設設定建立管理器
func NewDefault() *Manager {
	return New(DefaultConfig)
}

// RecordShot 記錄一次射擊（花費 + 回報）
// 回傳是否觸發補償，以及補償金額
func (m *Manager) RecordShot(playerID string, spend int, reward int) (triggered bool, bonusAmount int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := m.getOrCreateLocked(playerID)

	// 更新環形緩衝
	if len(s.Shots) < m.cfg.TrackingShots {
		s.Shots = append(s.Shots, ShotRecord{Spend: spend, Reward: reward})
	} else {
		// 移除最舊的記錄
		old := s.Shots[s.ShotIdx]
		s.TotalSpend -= old.Spend
		s.TotalReward -= old.Reward
		s.Shots[s.ShotIdx] = ShotRecord{Spend: spend, Reward: reward}
		s.ShotIdx = (s.ShotIdx + 1) % m.cfg.TrackingShots
	}
	s.TotalSpend += spend
	s.TotalReward += reward

	// 只有追蹤滿 N 次才判斷
	if len(s.Shots) < m.cfg.TrackingShots {
		return false, 0
	}

	// 冷卻檢查
	if !s.LastBonusAt.IsZero() {
		elapsed := time.Since(s.LastBonusAt)
		if elapsed < time.Duration(m.cfg.CooldownSecs)*time.Second {
			return false, 0
		}
	}

	// 觸發條件：花費 >= MinSpend 且 花費/回報 >= SpendThreshold
	if s.TotalSpend < m.cfg.MinSpend {
		return false, 0
	}
	if s.TotalReward <= 0 {
		// 完全沒有回報，直接觸發
	} else {
		ratio := float64(s.TotalSpend) / float64(s.TotalReward)
		if ratio < m.cfg.SpendThreshold {
			return false, 0
		}
	}

	// 計算補償金額
	netLoss := s.TotalSpend - s.TotalReward
	if netLoss <= 0 {
		return false, 0
	}

	bonus := int(float64(netLoss) * m.cfg.BaseRewardMult)
	maxBonus := int(float64(netLoss) * m.cfg.MaxRewardMult)
	if bonus > maxBonus {
		bonus = maxBonus
	}
	if bonus < m.cfg.MinReward {
		bonus = m.cfg.MinReward
	}

	// 記錄補償時間，重置追蹤
	s.LastBonusAt = time.Now()
	s.BonusCount++
	s.TotalSpend = 0
	s.TotalReward = 0
	s.Shots = s.Shots[:0]
	s.ShotIdx = 0

	return true, bonus
}

// GetSnapshot 取得玩家狀態快照（thread-safe）
type Snapshot struct {
	TotalSpend   int
	TotalReward  int
	ShotCount    int
	TrackingMax  int
	NetLoss      int
	RatioPercent int // 花費/回報 百分比（100 = 1:1，300 = 3:1）
	CooldownLeft int // 冷卻剩餘秒數
	BonusCount   int
}

func (m *Manager) GetSnapshot(playerID string) Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.states[playerID]
	if !ok {
		return Snapshot{TrackingMax: m.cfg.TrackingShots}
	}

	netLoss := s.TotalSpend - s.TotalReward
	if netLoss < 0 {
		netLoss = 0
	}

	ratioPercent := 0
	if s.TotalReward > 0 {
		ratioPercent = int(float64(s.TotalSpend) / float64(s.TotalReward) * 100)
	} else if s.TotalSpend > 0 {
		ratioPercent = 9999 // 完全沒有回報
	}

	cooldownLeft := 0
	if !s.LastBonusAt.IsZero() {
		elapsed := time.Since(s.LastBonusAt)
		remaining := time.Duration(m.cfg.CooldownSecs)*time.Second - elapsed
		if remaining > 0 {
			cooldownLeft = int(remaining.Seconds())
		}
	}

	return Snapshot{
		TotalSpend:   s.TotalSpend,
		TotalReward:  s.TotalReward,
		ShotCount:    len(s.Shots),
		TrackingMax:  m.cfg.TrackingShots,
		NetLoss:      netLoss,
		RatioPercent: ratioPercent,
		CooldownLeft: cooldownLeft,
		BonusCount:   s.BonusCount,
	}
}

// RemovePlayer 移除玩家狀態
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.states, playerID)
}

// ---- 內部輔助函數 ----

func (m *Manager) getOrCreateLocked(playerID string) *PlayerState {
	if s, ok := m.states[playerID]; ok {
		return s
	}
	s := &PlayerState{
		Shots: make([]ShotRecord, 0, m.cfg.TrackingShots),
	}
	m.states[playerID] = s
	return s
}
