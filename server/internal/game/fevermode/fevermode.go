// Package fevermode — 狂熱模式系統（DAY-133）
// 業界依據：Fire Kirin / Ocean King 系列的 Fever Mode
// 玩家在 5 秒內擊破 5 個目標觸發狂熱模式，期間所有獎勵 ×1.5，
// 繼續快速擊破可延長時間，製造「停不下來」的爽感。
// 這是業界最經典的短期留存機制之一。
package fevermode

import (
	"sync"
	"time"
)

// FeverConfig 狂熱模式設定
type FeverConfig struct {
	TriggerKills    int           // 觸發所需擊破數（預設 5）
	TriggerWindow   time.Duration // 觸發時間窗口（預設 5 秒）
	BaseDuration    time.Duration // 基礎持續時間（預設 15 秒）
	MaxDuration     time.Duration // 最大持續時間（預設 30 秒）
	ExtendPerKill   time.Duration // 每次擊破延長時間（預設 1 秒）
	MultBoost       float64       // 獎勵倍率加成（預設 1.5）
	CooldownSecs    int           // 冷卻時間（預設 20 秒）
}

// DefaultConfig 預設狂熱模式設定
var DefaultConfig = FeverConfig{
	TriggerKills:  5,
	TriggerWindow: 5 * time.Second,
	BaseDuration:  15 * time.Second,
	MaxDuration:   30 * time.Second,
	ExtendPerKill: 1 * time.Second,
	MultBoost:     1.5,
	CooldownSecs:  20,
}

// FeverState 玩家狂熱模式狀態
type FeverState int

const (
	FeverStateIdle    FeverState = iota // 未觸發
	FeverStateActive                    // 狂熱中
	FeverStateCooldown                  // 冷卻中
)

// PlayerFever 玩家的狂熱模式狀態
type PlayerFever struct {
	PlayerID    string
	State       FeverState
	KillTimes   []time.Time // 最近的擊破時間（用於觸發判斷）
	FeverEndAt  time.Time   // 狂熱結束時間
	CooldownAt  time.Time   // 冷卻開始時間
	TotalFevered int        // 本 session 觸發次數
}

// IsActive 是否正在狂熱中
func (p *PlayerFever) IsActive() bool {
	return p.State == FeverStateActive && time.Now().Before(p.FeverEndAt)
}

// SecondsLeft 狂熱剩餘秒數
func (p *PlayerFever) SecondsLeft() int {
	if !p.IsActive() {
		return 0
	}
	remaining := time.Until(p.FeverEndAt)
	if remaining <= 0 {
		return 0
	}
	return int(remaining.Seconds()) + 1
}

// CooldownLeft 冷卻剩餘秒數
func (p *PlayerFever) CooldownLeft(cooldownSecs int) int {
	if p.State != FeverStateCooldown {
		return 0
	}
	elapsed := time.Since(p.CooldownAt)
	cd := time.Duration(cooldownSecs) * time.Second
	if elapsed >= cd {
		return 0
	}
	return int((cd - elapsed).Seconds()) + 1
}

// Manager 狂熱模式管理器
type Manager struct {
	mu      sync.RWMutex
	config  FeverConfig
	players map[string]*PlayerFever
}

// New 建立新的狂熱模式管理器
func New() *Manager {
	return &Manager{
		config:  DefaultConfig,
		players: make(map[string]*PlayerFever),
	}
}

// getOrCreate 取得或建立玩家狀態
func (m *Manager) getOrCreate(playerID string) *PlayerFever {
	if p, ok := m.players[playerID]; ok {
		return p
	}
	p := &PlayerFever{
		PlayerID:  playerID,
		State:     FeverStateIdle,
		KillTimes: make([]time.Time, 0, 10),
	}
	m.players[playerID] = p
	return p
}

// RecordKill 記錄擊破，回傳（是否觸發狂熱, 是否延長狂熱, 當前倍率加成）
func (m *Manager) RecordKill(playerID string) (triggered bool, extended bool, multBoost float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p := m.getOrCreate(playerID)
	now := time.Now()

	// 更新冷卻狀態
	if p.State == FeverStateCooldown {
		if p.CooldownLeft(m.config.CooldownSecs) == 0 {
			p.State = FeverStateIdle
		}
	}

	// 更新狂熱狀態
	if p.State == FeverStateActive {
		if !now.Before(p.FeverEndAt) {
			// 狂熱結束，進入冷卻
			p.State = FeverStateCooldown
			p.CooldownAt = now
		}
	}

	// 狂熱中：延長時間
	if p.State == FeverStateActive {
		newEnd := p.FeverEndAt.Add(m.config.ExtendPerKill)
		if newEnd.Sub(now) <= m.config.MaxDuration {
			p.FeverEndAt = newEnd
		} else {
			p.FeverEndAt = now.Add(m.config.MaxDuration)
		}
		return false, true, m.config.MultBoost
	}

	// 冷卻中：不記錄擊破
	if p.State == FeverStateCooldown {
		return false, false, 1.0
	}

	// Idle：記錄擊破時間，判斷是否觸發
	p.KillTimes = append(p.KillTimes, now)

	// 清除超出時間窗口的舊記錄
	cutoff := now.Add(-m.config.TriggerWindow)
	valid := p.KillTimes[:0]
	for _, t := range p.KillTimes {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	p.KillTimes = valid

	// 判斷是否觸發
	if len(p.KillTimes) >= m.config.TriggerKills {
		p.State = FeverStateActive
		p.FeverEndAt = now.Add(m.config.BaseDuration)
		p.KillTimes = p.KillTimes[:0] // 清空觸發記錄
		p.TotalFevered++
		return true, false, m.config.MultBoost
	}

	return false, false, 1.0
}

// GetMultBoost 取得當前倍率加成（狂熱中返回 MultBoost，否則返回 1.0）
func (m *Manager) GetMultBoost(playerID string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.players[playerID]
	if !ok {
		return 1.0
	}
	if p.IsActive() {
		return m.config.MultBoost
	}
	return 1.0
}

// CheckExpiry 檢查並更新過期狀態，回傳是否剛結束
func (m *Manager) CheckExpiry(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, ok := m.players[playerID]
	if !ok {
		return false
	}
	if p.State == FeverStateActive && !time.Now().Before(p.FeverEndAt) {
		p.State = FeverStateCooldown
		p.CooldownAt = time.Now()
		return true
	}
	return false
}

// GetSnapshot 取得玩家快照
type Snapshot struct {
	PlayerID     string
	State        FeverState
	SecondsLeft  int
	CooldownLeft int
	MultBoost    float64
	KillProgress int // 觸發進度（0-TriggerKills）
	TotalFevered int
}

func (m *Manager) GetSnapshot(playerID string) Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.players[playerID]
	if !ok {
		return Snapshot{PlayerID: playerID, MultBoost: 1.0}
	}
	snap := Snapshot{
		PlayerID:     playerID,
		State:        p.State,
		SecondsLeft:  p.SecondsLeft(),
		CooldownLeft: p.CooldownLeft(m.config.CooldownSecs),
		MultBoost:    1.0,
		TotalFevered: p.TotalFevered,
	}
	if p.IsActive() {
		snap.MultBoost = m.config.MultBoost
	}
	// 計算觸發進度（只在 Idle 狀態有意義）
	if p.State == FeverStateIdle {
		now := time.Now()
		cutoff := now.Add(-m.config.TriggerWindow)
		count := 0
		for _, t := range p.KillTimes {
			if t.After(cutoff) {
				count++
			}
		}
		snap.KillProgress = count
	}
	return snap
}

// RemovePlayer 移除玩家狀態
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.players, playerID)
}

// GetConfig 取得設定
func (m *Manager) GetConfig() FeverConfig {
	return m.config
}

// TickExpiry 批次檢查所有玩家的過期狀態，回傳剛結束的玩家 ID 列表
func (m *Manager) TickExpiry() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	expired := make([]string, 0)
	for id, p := range m.players {
		if p.State == FeverStateActive && !now.Before(p.FeverEndAt) {
			p.State = FeverStateCooldown
			p.CooldownAt = now
			expired = append(expired, id)
		}
	}
	return expired
}
