// Package multstorm 全服倍率風暴系統（DAY-138）
// 業界依據：findingdulcinea.com 2026「admin events multiply luck by thousands」
// Fisch 2026 server-wide bonus events + Royal Fishing dual-ring roulette
// 隨機觸發全服倍率提升，讓所有目標獎勵暫時翻倍，製造「全場瘋狂」高峰體驗
package multstorm

import (
	"math/rand"
	"sync"
	"time"
)

// StormTier 風暴等級
type StormTier struct {
	Name        string  // 等級名稱
	Icon        string  // 圖示
	Color       string  // 顏色（hex）
	MultBoost   float64 // 倍率加成（疊加到所有獎勵）
	Duration    float64 // 持續秒數
	TriggerProb float64 // 觸發機率（每次 tick）
}

// 三個風暴等級
var StormTiers = []StormTier{
	{
		Name:        "⚡ 閃電風暴",
		Icon:        "⚡",
		Color:       "#FFE066",
		MultBoost:   2.0,
		Duration:    20.0,
		TriggerProb: 0.008, // 0.8%
	},
	{
		Name:        "🌊 海嘯風暴",
		Icon:        "🌊",
		Color:       "#4A90D9",
		MultBoost:   3.0,
		Duration:    15.0,
		TriggerProb: 0.003, // 0.3%
	},
	{
		Name:        "🌈 彩虹風暴",
		Icon:        "🌈",
		Color:       "#FF69B4",
		MultBoost:   5.0,
		Duration:    10.0,
		TriggerProb: 0.001, // 0.1%
	},
}

// StormConfig 風暴設定
type StormConfig struct {
	CooldownSecs float64 // 冷卻時間（預設 180 秒）
	TickInterval float64 // 觸發檢查間隔（預設 1 秒）
}

// DefaultConfig 預設設定
func DefaultConfig() StormConfig {
	return StormConfig{
		CooldownSecs: 180.0,
		TickInterval: 1.0,
	}
}

// StormSession 風暴 session
type StormSession struct {
	Tier      StormTier
	StartAt   time.Time
	EndAt     time.Time
	IsActive  bool
}

// Manager 倍率風暴管理器
type Manager struct {
	mu        sync.RWMutex
	config    StormConfig
	session   *StormSession
	lastEndAt time.Time
	rng       *rand.Rand
}

// New 建立倍率風暴管理器
func New(cfg StormConfig) *Manager {
	return &Manager{
		config: cfg,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewDefault 使用預設設定建立管理器
func NewDefault() *Manager {
	return New(DefaultConfig())
}

// TryTrigger 嘗試觸發風暴（每秒呼叫一次）
// 回傳觸發的 session（nil = 未觸發）
func (m *Manager) TryTrigger() *StormSession {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 已有活躍風暴
	if m.session != nil && m.session.IsActive {
		return nil
	}

	// 冷卻檢查
	if !m.lastEndAt.IsZero() && time.Since(m.lastEndAt).Seconds() < m.config.CooldownSecs {
		return nil
	}

	// 依機率嘗試觸發各等級（從高到低，高等級優先）
	for i := len(StormTiers) - 1; i >= 0; i-- {
		tier := StormTiers[i]
		if m.rng.Float64() < tier.TriggerProb {
			now := time.Now()
			m.session = &StormSession{
				Tier:     tier,
				StartAt:  now,
				EndAt:    now.Add(time.Duration(tier.Duration * float64(time.Second))),
				IsActive: true,
			}
			return m.session
		}
	}
	return nil
}

// ForceStart 強制觸發指定等級的風暴（Prototype 展示用）
func (m *Manager) ForceStart(tierIndex int) *StormSession {
	m.mu.Lock()
	defer m.mu.Unlock()

	if tierIndex < 0 || tierIndex >= len(StormTiers) {
		tierIndex = 0
	}
	tier := StormTiers[tierIndex]
	now := time.Now()
	m.session = &StormSession{
		Tier:     tier,
		StartAt:  now,
		EndAt:    now.Add(time.Duration(tier.Duration * float64(time.Second))),
		IsActive: true,
	}
	return m.session
}

// CheckExpiry 檢查風暴是否結束
// 回傳 true 表示剛剛結束（需要廣播結束）
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

// GetMultBoost 取得當前倍率加成（1.0 = 無加成）
func (m *Manager) GetMultBoost() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.session == nil || !m.session.IsActive {
		return 1.0
	}
	return m.session.Tier.MultBoost
}

// IsActive 是否有活躍風暴
func (m *Manager) IsActive() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.session != nil && m.session.IsActive
}

// GetSnapshot 取得當前風暴快照
func (m *Manager) GetSnapshot() StormSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.session == nil || !m.session.IsActive {
		return StormSnapshot{IsActive: false}
	}

	secsLeft := time.Until(m.session.EndAt).Seconds()
	if secsLeft < 0 {
		secsLeft = 0
	}

	return StormSnapshot{
		IsActive:    true,
		TierName:    m.session.Tier.Name,
		TierIcon:    m.session.Tier.Icon,
		TierColor:   m.session.Tier.Color,
		MultBoost:   m.session.Tier.MultBoost,
		SecondsLeft: secsLeft,
	}
}

// StormSnapshot 風暴快照
type StormSnapshot struct {
	IsActive    bool
	TierName    string
	TierIcon    string
	TierColor   string
	MultBoost   float64
	SecondsLeft float64
}
