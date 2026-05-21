// chainlongwheel.go — 千龍王強化輪盤系統（DAY-148）
// 業界依據：Royal Fishing JILI 2026「ChainLong King — capture this golden dragon to trigger
// the dual-ring roulette. The ChainLong King itself can award up to 1000X mega wins.」
// 設計：內環（5x/10x/20x/30x/50x）× 外環（2x/3x/5x/7x/10x/20x）= 最高 1000x
// 比普通雙環輪盤（最高 150x）強 6.7 倍，是全遊戲最高倍率的個人機制
package chainlongwheel

import (
	"math/rand"
	"sync"
	"time"
)

// InnerRing 千龍王內環倍率選項（高倍率版）
// 加權：5x×35, 10x×28, 20x×18, 30x×12, 50x×7
var InnerRing = []float64{5.0, 10.0, 20.0, 30.0, 50.0}

// InnerWeights 內環加權（對應 InnerRing）
var InnerWeights = []int{35, 28, 18, 12, 7}

// OuterRing 千龍王外環倍率選項（高倍率版）
// 加權：2x×30, 3x×25, 5x×20, 7x×13, 10x×8, 20x×4
var OuterRing = []float64{2.0, 3.0, 5.0, 7.0, 10.0, 20.0}

// OuterWeights 外環加權（對應 OuterRing）
var OuterWeights = []int{30, 25, 20, 13, 8, 4}

// SpinDuration 旋轉持續秒數
const SpinDuration = 4.0

// CooldownSecs 冷卻秒數（千龍王極稀有，冷卻較短）
const CooldownSecs = 30

// Session 一次千龍王輪盤 session
type Session struct {
	PlayerID    string
	TargetMult  float64   // 千龍王本身的倍率（150-1000x）
	BaseReward  int       // 千龍王擊破的基礎獎勵
	InnerResult float64   // 內環結果（預先決定）
	OuterResult float64   // 外環結果（預先決定）
	StartedAt   time.Time
	StoppedAt   time.Time
	IsStopped   bool
}

// FinalMultiplier 計算最終倍率（內環 × 外環）
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

// Manager 千龍王輪盤管理器
type Manager struct {
	mu       sync.Mutex
	sessions map[string]*Session  // playerID -> active session
	cooldown map[string]time.Time // playerID -> cooldown end
	rng      *rand.Rand
}

// New 建立管理器
func New() *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
		cooldown: make(map[string]time.Time),
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// weightedRandIndex 加權隨機選擇索引
func (m *Manager) weightedRandIndex(weights []int) int {
	total := 0
	for _, w := range weights {
		total += w
	}
	r := m.rng.Intn(total)
	cumulative := 0
	for i, w := range weights {
		cumulative += w
		if r < cumulative {
			return i
		}
	}
	return len(weights) - 1
}

// HasActiveSession 檢查玩家是否有活躍 session
func (m *Manager) HasActiveSession(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, exists := m.sessions[playerID]
	return exists
}

// CanTrigger 檢查是否可以觸發（無活躍 session + 不在冷卻中）
func (m *Manager) CanTrigger(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

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
	return true
}

// StartSession 開始一次千龍王輪盤 session
// 結果預先決定（公平性保證），玩家「停止」只是視覺互動
func (m *Manager) StartSession(playerID string, targetMult float64, baseReward int) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 加權隨機選擇內外環結果
	innerIdx := m.weightedRandIndex(InnerWeights)
	outerIdx := m.weightedRandIndex(OuterWeights)

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
	m.cooldown[playerID] = time.Now().Add(CooldownSecs * time.Second)
	delete(m.sessions, playerID)

	return s
}

// TickAutoStop 每秒檢查所有活躍 session 是否超時
// 回傳所有超時自動停止的 session 列表
func (m *Manager) TickAutoStop() []*Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	var expired []*Session
	now := time.Now()
	spinDur := time.Duration(SpinDuration * float64(time.Second))

	for playerID, s := range m.sessions {
		if s.IsStopped {
			delete(m.sessions, playerID)
			continue
		}
		if now.Sub(s.StartedAt) >= spinDur {
			s.IsStopped = true
			s.StoppedAt = now
			m.cooldown[playerID] = now.Add(CooldownSecs * time.Second)
			delete(m.sessions, playerID)
			expired = append(expired, s)
		}
	}
	return expired
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
