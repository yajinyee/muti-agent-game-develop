// dualroulette.go — 雙環輪盤系統（DAY-139）
// 業界依據：Royal Fishing JILI 2026 ChainLong King Dual-Ring Roulette
// 擊破高倍率目標後觸發，內外圈相乘最高 150x，製造「技巧感」
package dualroulette

import (
	"math/rand"
	"sync"
	"time"
)

// InnerRing 內環倍率選項
var InnerRing = []float64{2.0, 3.0, 5.0, 8.0, 10.0}

// OuterRing 外環倍率選項
var OuterRing = []float64{2.0, 3.0, 5.0, 7.0, 10.0, 15.0}

// TriggerConfig 觸發設定
type TriggerConfig struct {
	MinMultiplier float64 // 觸發所需最低目標倍率（預設 30x）
	TriggerChance float64 // 觸發機率（0.0-1.0，預設 0.15）
	SpinDuration  float64 // 旋轉持續秒數（預設 3.0）
	CooldownSecs  int     // 冷卻秒數（預設 60）
}

// DefaultConfig 預設設定
func DefaultConfig() TriggerConfig {
	return TriggerConfig{
		MinMultiplier: 30.0,
		TriggerChance: 0.15,
		SpinDuration:  3.0,
		CooldownSecs:  60,
	}
}

// Session 一次輪盤 session
type Session struct {
	PlayerID    string
	TargetMult  float64   // 觸發目標的倍率
	BaseReward  int       // 觸發時的基礎獎勵
	InnerResult float64   // 內環結果（停止後確定）
	OuterResult float64   // 外環結果（停止後確定）
	Combined    float64   // 最終倍率 = InnerResult × OuterResult
	StartedAt   time.Time
	StoppedAt   time.Time
	IsStopped   bool
}

// FinalMultiplier 計算最終倍率
func (s *Session) FinalMultiplier() float64 {
	if !s.IsStopped {
		return 0
	}
	return s.InnerResult * s.OuterResult
}

// BonusReward 計算額外獎勵（基礎獎勵 × 最終倍率）
func (s *Session) BonusReward() int {
	if !s.IsStopped {
		return 0
	}
	return int(float64(s.BaseReward) * s.FinalMultiplier())
}

// Snapshot 輪盤快照（用於廣播）
type Snapshot struct {
	PlayerID    string  `json:"player_id"`
	TargetMult  float64 `json:"target_mult"`
	BaseReward  int     `json:"base_reward"`
	InnerResult float64 `json:"inner_result"`
	OuterResult float64 `json:"outer_result"`
	Combined    float64 `json:"combined"`
	BonusReward int     `json:"bonus_reward"`
	IsStopped   bool    `json:"is_stopped"`
}

// Manager 雙環輪盤管理器
type Manager struct {
	mu       sync.Mutex
	config   TriggerConfig
	sessions map[string]*Session // playerID -> active session
	cooldown map[string]time.Time // playerID -> cooldown end
	rng      *rand.Rand
}

// New 建立管理器
func New(cfg TriggerConfig) *Manager {
	return &Manager{
		config:   cfg,
		sessions: make(map[string]*Session),
		cooldown: make(map[string]time.Time),
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewDefault 使用預設設定建立管理器
func NewDefault() *Manager {
	return New(DefaultConfig())
}

// CanTrigger 檢查是否可以觸發（冷卻 + 無活躍 session）
func (m *Manager) CanTrigger(playerID string, targetMult float64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 倍率門檻
	if targetMult < m.config.MinMultiplier {
		return false
	}
	// 已有活躍 session
	if _, exists := m.sessions[playerID]; exists {
		return false
	}
	// 冷卻中
	if cd, ok := m.cooldown[playerID]; ok {
		if time.Now().Before(cd) {
			return false
		}
	}
	// 機率觸發
	return m.rng.Float64() < m.config.TriggerChance
}

// StartSession 開始一次輪盤 session
// 回傳 session（含預先決定的內外環結果，但不告訴玩家）
func (m *Manager) StartSession(playerID string, targetMult float64, baseReward int) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 預先決定結果（公平性：結果已定，玩家只是「選擇停止時機」）
	innerIdx := m.rng.Intn(len(InnerRing))
	outerIdx := m.rng.Intn(len(OuterRing))

	s := &Session{
		PlayerID:    playerID,
		TargetMult:  targetMult,
		BaseReward:  baseReward,
		InnerResult: InnerRing[innerIdx],
		OuterResult: OuterRing[outerIdx],
		StartedAt:   time.Now(),
	}
	m.sessions[playerID] = s
	return s
}

// StopSession 玩家停止輪盤，回傳 session 結果
// 若無活躍 session 或已停止，回傳 nil
func (m *Manager) StopSession(playerID string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, exists := m.sessions[playerID]
	if !exists || s.IsStopped {
		return nil
	}

	s.IsStopped = true
	s.StoppedAt = time.Now()

	// 設定冷卻
	m.cooldown[playerID] = time.Now().Add(time.Duration(m.config.CooldownSecs) * time.Second)
	// 移除 session
	delete(m.sessions, playerID)

	return s
}

// AutoStop 自動停止（超時後由 server 自動停止）
// 回傳 session 結果，若無活躍 session 回傳 nil
func (m *Manager) AutoStop(playerID string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, exists := m.sessions[playerID]
	if !exists || s.IsStopped {
		return nil
	}

	// 超時自動停止
	if time.Since(s.StartedAt) < time.Duration(m.config.SpinDuration*float64(time.Second)) {
		return nil // 還沒到時間
	}

	s.IsStopped = true
	s.StoppedAt = time.Now()

	m.cooldown[playerID] = time.Now().Add(time.Duration(m.config.CooldownSecs) * time.Second)
	delete(m.sessions, playerID)

	return s
}

// GetActiveSession 取得玩家的活躍 session（唯讀）
func (m *Manager) GetActiveSession(playerID string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sessions[playerID]
}

// HasActiveSession 檢查玩家是否有活躍 session
func (m *Manager) HasActiveSession(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, exists := m.sessions[playerID]
	return exists
}

// GetCooldownLeft 取得冷卻剩餘秒數
func (m *Manager) GetCooldownLeft(playerID string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	if cd, ok := m.cooldown[playerID]; ok {
		left := time.Until(cd)
		if left > 0 {
			return int(left.Seconds())
		}
	}
	return 0
}

// RemovePlayer 玩家離線時清理
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, playerID)
	delete(m.cooldown, playerID)
}

// GetSnapshot 取得 session 快照（用於廣播）
func (m *Manager) GetSnapshot(s *Session) Snapshot {
	return Snapshot{
		PlayerID:    s.PlayerID,
		TargetMult:  s.TargetMult,
		BaseReward:  s.BaseReward,
		InnerResult: s.InnerResult,
		OuterResult: s.OuterResult,
		Combined:    s.FinalMultiplier(),
		BonusReward: s.BonusReward(),
		IsStopped:   s.IsStopped,
	}
}

// TickAutoStop 每秒檢查所有活躍 session 是否超時
// 回傳所有超時自動停止的 session 列表
func (m *Manager) TickAutoStop() []*Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	var expired []*Session
	now := time.Now()
	spinDur := time.Duration(m.config.SpinDuration * float64(time.Second))

	for playerID, s := range m.sessions {
		if s.IsStopped {
			delete(m.sessions, playerID)
			continue
		}
		if now.Sub(s.StartedAt) >= spinDur {
			s.IsStopped = true
			s.StoppedAt = now
			m.cooldown[playerID] = now.Add(time.Duration(m.config.CooldownSecs) * time.Second)
			delete(m.sessions, playerID)
			expired = append(expired, s)
		}
	}
	return expired
}
