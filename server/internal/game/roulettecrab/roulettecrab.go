// Package roulettecrab — 黃金輪盤螃蟹系統（DAY-167）
// 業界依據：King of Treasures Plus 2026「Roulette Crab — triggers Golden Roulette bonus game,
// player hits SHOOT to stop wheel, wins the amount listed where it stops.」
// 設計：單環輪盤（8格：10x-200x），比千龍王雙環更簡單直接，適合中等 betLevel 玩家
// 結果預先決定（公平性保證），玩家「停止」只是視覺互動（業界標準做法）
package roulettecrab

import (
	"math/rand"
	"sync"
	"time"
)

// WheelSlots 輪盤格子（8格，加權隨機）
// 業界依據：King of Treasures Plus 輪盤設計，低倍率高機率，高倍率低機率
var WheelSlots = []float64{10.0, 20.0, 30.0, 50.0, 80.0, 100.0, 150.0, 200.0}

// WheelWeights 輪盤加權（對應 WheelSlots）
// 10x×30, 20x×25, 30x×18, 50x×12, 80x×7, 100x×4, 150x×3, 200x×1
var WheelWeights = []int{30, 25, 18, 12, 7, 4, 3, 1}

// SpinDuration 旋轉持續秒數（玩家有 4 秒可以停止）
const SpinDuration = 4.0

// CooldownSecs 冷卻秒數（螃蟹比千龍王更常見，冷卻較短）
const CooldownSecs = 20

// Session 一次輪盤螃蟹 session
type Session struct {
	PlayerID   string
	TargetMult float64   // 螃蟹本身的倍率（20-40x）
	BaseReward int       // 螃蟹擊破的基礎獎勵
	WheelResult float64  // 輪盤結果（預先決定）
	SlotIndex  int       // 輪盤格子索引（0-7，用於 Client 動畫定位）
	StartedAt  time.Time
	StoppedAt  time.Time
	IsStopped  bool
}

// BonusReward 計算額外獎勵（基礎獎勵 × 輪盤倍率）
func (s *Session) BonusReward() int {
	if !s.IsStopped {
		return 0
	}
	return int(float64(s.BaseReward) * s.WheelResult)
}

// Snapshot 輪盤快照（用於廣播）
type Snapshot struct {
	PlayerID    string  `json:"player_id"`
	TargetMult  float64 `json:"target_mult"`
	BaseReward  int     `json:"base_reward"`
	WheelResult float64 `json:"wheel_result"`
	SlotIndex   int     `json:"slot_index"`
	BonusReward int     `json:"bonus_reward"`
	IsStopped   bool    `json:"is_stopped"`
}

// Manager 輪盤螃蟹管理器
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
	if _, exists := m.sessions[playerID]; exists {
		return false
	}
	if cd, ok := m.cooldown[playerID]; ok {
		if time.Now().Before(cd) {
			return false
		}
	}
	return true
}

// StartSession 開始一次輪盤螃蟹 session
// 結果預先決定（公平性保證），玩家「停止」只是視覺互動
func (m *Manager) StartSession(playerID string, targetMult float64, baseReward int) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 加權隨機選擇輪盤格子
	slotIdx := m.weightedRandIndex(WheelWeights)

	s := &Session{
		PlayerID:    playerID,
		TargetMult:  targetMult,
		BaseReward:  baseReward,
		WheelResult: WheelSlots[slotIdx],
		SlotIndex:   slotIdx,
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
		WheelResult: s.WheelResult,
		SlotIndex:   s.SlotIndex,
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
