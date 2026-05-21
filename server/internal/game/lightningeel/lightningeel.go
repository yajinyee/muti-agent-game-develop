// Package lightningeel — 閃電鰻連鎖攻擊系統（DAY-132）
// 業界依據：JILI Royal Fishing 2026 — 「The 60x lightning eel creates chain reactions
// that jump between nearby fish. Once activated, electric shocks continue spreading
// until targeting disengages, creating cascading capture sequences.」
// 閃電鰻是一種特殊目標物，擊破後釋放閃電連鎖，在附近目標之間跳躍傳導，
// 每次跳躍有機率直接擊破目標，製造「一箭多雕」的爽感。
package lightningeel

import (
	"math/rand"
	"sync"
	"time"
)

// ChainConfig 連鎖攻擊設定
type ChainConfig struct {
	MaxJumps       int     // 最大跳躍次數
	JumpKillChance float64 // 每次跳躍的擊破機率
	JumpMultMod    float64 // 跳躍獎勵倍率修正（相對於目標原始倍率）
	JumpRangeUnits float64 // 跳躍範圍（遊戲單位，Client 用於視覺）
	CooldownSecs   int     // 觸發冷卻（秒）
}

// DefaultConfig 預設連鎖設定
var DefaultConfig = ChainConfig{
	MaxJumps:       5,    // 最多跳 5 次
	JumpKillChance: 0.45, // 每次跳躍 45% 擊破機率
	JumpMultMod:    0.6,  // 跳躍獎勵 = 目標倍率 × 0.6（略低於直接擊破）
	JumpRangeUnits: 200,  // 200 單位範圍內的目標都可能被連鎖
	CooldownSecs:   8,    // 8 秒冷卻（每個玩家獨立）
}

// JumpResult 單次跳躍結果
type JumpResult struct {
	TargetInstanceID string  // 被跳躍的目標 instance ID
	TargetDefID      string  // 目標定義 ID
	TargetName       string  // 目標名稱
	Killed           bool    // 是否擊破
	Multiplier       float64 // 目標原始倍率
	Reward           int64   // 實際獎勵（Killed 才有）
	JumpIndex        int     // 第幾次跳躍（1-based）
}

// ChainSession 一次連鎖攻擊的 session
type ChainSession struct {
	PlayerID        string
	TriggerTargetID string    // 觸發連鎖的閃電鰻 instance ID
	StartAt         time.Time
	Jumps           []JumpResult
	TotalReward     int64
	TotalKills      int
}

// PlayerState 玩家的閃電鰻狀態
type PlayerState struct {
	PlayerID    string
	LastChainAt time.Time // 上次觸發連鎖的時間（冷卻計算）
}

// Manager 閃電鰻連鎖攻擊管理器
type Manager struct {
	mu      sync.RWMutex
	config  ChainConfig
	players map[string]*PlayerState
	rng     *rand.Rand
}

// New 建立新的閃電鰻管理器
func New() *Manager {
	return &Manager{
		config:  DefaultConfig,
		players: make(map[string]*PlayerState),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// getOrCreatePlayer 取得或建立玩家狀態
func (m *Manager) getOrCreatePlayer(playerID string) *PlayerState {
	if s, ok := m.players[playerID]; ok {
		return s
	}
	s := &PlayerState{PlayerID: playerID}
	m.players[playerID] = s
	return s
}

// CanTrigger 檢查玩家是否可以觸發連鎖（冷卻檢查）
func (m *Manager) CanTrigger(playerID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.players[playerID]
	if !ok {
		return true
	}
	if s.LastChainAt.IsZero() {
		return true
	}
	return time.Since(s.LastChainAt) >= time.Duration(m.config.CooldownSecs)*time.Second
}

// CooldownLeft 取得冷卻剩餘秒數
func (m *Manager) CooldownLeft(playerID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.players[playerID]
	if !ok {
		return 0
	}
	if s.LastChainAt.IsZero() {
		return 0
	}
	elapsed := time.Since(s.LastChainAt)
	cooldown := time.Duration(m.config.CooldownSecs) * time.Second
	if elapsed >= cooldown {
		return 0
	}
	return int((cooldown - elapsed).Seconds()) + 1
}

// NearbyTarget 附近目標的資訊（由 game 層傳入）
type NearbyTarget struct {
	InstanceID string
	DefID      string
	Name       string
	Multiplier float64
	X, Y       float64 // 目標位置（用於距離計算）
}

// ExecuteChain 執行連鎖攻擊
// nearbyTargets：閃電鰻附近的目標列表（已按距離排序，由 game 層提供）
// betCost：玩家當前投注金額（用於計算獎勵）
// 回傳：ChainSession（包含所有跳躍結果）
func (m *Manager) ExecuteChain(
	playerID string,
	triggerTargetID string,
	nearbyTargets []NearbyTarget,
	betCost int64,
) *ChainSession {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := m.getOrCreatePlayer(playerID)
	s.LastChainAt = time.Now()

	session := &ChainSession{
		PlayerID:        playerID,
		TriggerTargetID: triggerTargetID,
		StartAt:         time.Now(),
		Jumps:           make([]JumpResult, 0, m.config.MaxJumps),
	}

	// 已跳躍過的目標（避免重複跳躍）
	jumped := make(map[string]bool)
	jumped[triggerTargetID] = true

	maxJumps := m.config.MaxJumps
	if len(nearbyTargets) < maxJumps {
		maxJumps = len(nearbyTargets)
	}

	for i := 0; i < maxJumps; i++ {
		// 找下一個未跳躍的目標
		var target *NearbyTarget
		for j := range nearbyTargets {
			if !jumped[nearbyTargets[j].InstanceID] {
				target = &nearbyTargets[j]
				break
			}
		}
		if target == nil {
			break
		}
		jumped[target.InstanceID] = true

		// 判斷是否擊破
		killed := m.rng.Float64() < m.config.JumpKillChance
		var reward int64
		if killed {
			reward = int64(float64(betCost) * target.Multiplier * m.config.JumpMultMod)
			if reward < 1 {
				reward = 1
			}
			session.TotalKills++
			session.TotalReward += reward
		}

		session.Jumps = append(session.Jumps, JumpResult{
			TargetInstanceID: target.InstanceID,
			TargetDefID:      target.DefID,
			TargetName:       target.Name,
			Killed:           killed,
			Multiplier:       target.Multiplier,
			Reward:           reward,
			JumpIndex:        i + 1,
		})
	}

	return session
}

// RemovePlayer 移除玩家狀態（玩家離線時清理）
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.players, playerID)
}

// GetConfig 取得連鎖設定（供 Client 初始化用）
func (m *Manager) GetConfig() ChainConfig {
	return m.config
}

// Snapshot 快照（供 Client 查詢）
type Snapshot struct {
	PlayerID     string
	CooldownLeft int
	Config       ChainConfig
}

// GetSnapshot 取得玩家快照
func (m *Manager) GetSnapshot(playerID string) Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.players[playerID]
	cooldown := 0
	if ok && !s.LastChainAt.IsZero() {
		elapsed := time.Since(s.LastChainAt)
		cd := time.Duration(m.config.CooldownSecs) * time.Second
		if elapsed < cd {
			cooldown = int((cd-elapsed).Seconds()) + 1
		}
	}
	return Snapshot{
		PlayerID:     playerID,
		CooldownLeft: cooldown,
		Config:       m.config,
	}
}
