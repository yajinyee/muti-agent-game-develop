// megacatch.go — 全服 Mega Catch 事件系統（DAY-140）
// 業界依據：Ocean King 系列「Mega Catch」— 全場所有目標同時出現高倍率
// 持續 10-15 秒，所有目標擊破獎勵翻倍，製造「全場瘋狂搶魚」高峰體驗
// 與倍率風暴（MultStorm）的差異：
//   - MultStorm：疊加倍率（×2-5），持續 10-20 秒，機率觸發
//   - MegaCatch：直接提升目標生成倍率（讓高倍率目標大量湧現），持續 10-15 秒，BOSS 擊殺後觸發
package megacatch

import (
	"math/rand"
	"sync"
	"time"
)

// EventTier Mega Catch 等級
type EventTier struct {
	Name          string  // 等級名稱
	Icon          string  // 圖示
	Color         string  // 顏色（hex）
	SpawnMultBoost float64 // 高倍率目標生成加成（加到 rareBonus）
	RewardBoost   float64 // 獎勵倍率加成（疊加到所有獎勵）
	Duration      float64 // 持續秒數
}

// 三個等級
var EventTiers = []EventTier{
	{
		Name:          "🎣 大豐收",
		Icon:          "🎣",
		Color:         "#66CCFF",
		SpawnMultBoost: 0.20, // 稀有目標生成 +20%
		RewardBoost:   1.5,   // 獎勵 ×1.5
		Duration:      12.0,
	},
	{
		Name:          "🌟 超級豐收",
		Icon:          "🌟",
		Color:         "#FFD700",
		SpawnMultBoost: 0.35, // 稀有目標生成 +35%
		RewardBoost:   2.0,   // 獎勵 ×2.0
		Duration:      10.0,
	},
	{
		Name:          "💎 傳說豐收",
		Icon:          "💎",
		Color:         "#FF69B4",
		SpawnMultBoost: 0.50, // 稀有目標生成 +50%
		RewardBoost:   3.0,   // 獎勵 ×3.0
		Duration:      8.0,
	},
}

// Session 一次 Mega Catch 事件
type Session struct {
	Tier      EventTier
	StartedAt time.Time
	EndsAt    time.Time
}

// IsActive 是否仍在進行中
func (s *Session) IsActive() bool {
	return time.Now().Before(s.EndsAt)
}

// SecondsLeft 剩餘秒數
func (s *Session) SecondsLeft() float64 {
	left := time.Until(s.EndsAt).Seconds()
	if left < 0 {
		return 0
	}
	return left
}

// Snapshot 事件快照（用於廣播）
type Snapshot struct {
	IsActive      bool    `json:"is_active"`
	TierName      string  `json:"tier_name"`
	TierIcon      string  `json:"tier_icon"`
	TierColor     string  `json:"tier_color"`
	SpawnBoost    float64 `json:"spawn_boost"`
	RewardBoost   float64 `json:"reward_boost"`
	SecondsLeft   float64 `json:"seconds_left"`
	TotalDuration float64 `json:"total_duration"`
}

// Config 設定
type Config struct {
	CooldownSecs    int     // 冷卻秒數（預設 120）
	BossKillChance  float64 // BOSS 擊殺後觸發機率（預設 0.6）
	RandomChance    float64 // 每分鐘隨機觸發機率（預設 0.05）
}

// DefaultConfig 預設設定
func DefaultConfig() Config {
	return Config{
		CooldownSecs:   120,
		BossKillChance: 0.60,
		RandomChance:   0.05,
	}
}

// Manager Mega Catch 管理器
type Manager struct {
	mu       sync.Mutex
	config   Config
	session  *Session
	cooldown time.Time
	rng      *rand.Rand
}

// New 建立管理器
func New(cfg Config) *Manager {
	return &Manager{
		config: cfg,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewDefault 使用預設設定建立管理器
func NewDefault() *Manager {
	return New(DefaultConfig())
}

// IsActive 是否有活躍事件
func (m *Manager) IsActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.session != nil && m.session.IsActive()
}

// GetRewardBoost 取得當前獎勵倍率加成（無事件時回傳 1.0）
func (m *Manager) GetRewardBoost() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.session == nil || !m.session.IsActive() {
		return 1.0
	}
	return m.session.Tier.RewardBoost
}

// GetSpawnBoost 取得當前稀有目標生成加成（無事件時回傳 0.0）
func (m *Manager) GetSpawnBoost() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.session == nil || !m.session.IsActive() {
		return 0.0
	}
	return m.session.Tier.SpawnMultBoost
}

// CanTrigger 是否可以觸發（冷卻 + 無活躍事件）
func (m *Manager) CanTrigger() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.session != nil && m.session.IsActive() {
		return false
	}
	return time.Now().After(m.cooldown)
}

// TryTriggerBossKill BOSS 擊殺後嘗試觸發
// 回傳觸發的 Session（nil 表示未觸發）
func (m *Manager) TryTriggerBossKill() *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session != nil && m.session.IsActive() {
		return nil
	}
	if !time.Now().After(m.cooldown) {
		return nil
	}
	if m.rng.Float64() >= m.config.BossKillChance {
		return nil
	}

	return m.startSessionLocked()
}

// TryTriggerRandom 每分鐘隨機嘗試觸發
// 回傳觸發的 Session（nil 表示未觸發）
func (m *Manager) TryTriggerRandom() *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session != nil && m.session.IsActive() {
		return nil
	}
	if !time.Now().After(m.cooldown) {
		return nil
	}
	if m.rng.Float64() >= m.config.RandomChance {
		return nil
	}

	return m.startSessionLocked()
}

// ForceStart 強制觸發（Prototype 展示用）
func (m *Manager) ForceStart(tierIndex int) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	if tierIndex < 0 || tierIndex >= len(EventTiers) {
		tierIndex = 0
	}
	return m.startSessionWithTierLocked(EventTiers[tierIndex])
}

// startSessionLocked 內部觸發（已持鎖）
func (m *Manager) startSessionLocked() *Session {
	// 加權隨機選擇等級（傳說最稀有）
	roll := m.rng.Float64()
	var tier EventTier
	if roll < 0.10 {
		tier = EventTiers[2] // 傳說豐收 10%
	} else if roll < 0.35 {
		tier = EventTiers[1] // 超級豐收 25%
	} else {
		tier = EventTiers[0] // 大豐收 65%
	}
	return m.startSessionWithTierLocked(tier)
}

// startSessionWithTierLocked 用指定等級觸發（已持鎖）
func (m *Manager) startSessionWithTierLocked(tier EventTier) *Session {
	now := time.Now()
	s := &Session{
		Tier:      tier,
		StartedAt: now,
		EndsAt:    now.Add(time.Duration(tier.Duration * float64(time.Second))),
	}
	m.session = s
	m.cooldown = s.EndsAt.Add(time.Duration(m.config.CooldownSecs) * time.Second)
	return s
}

// CheckExpiry 檢查是否過期，回傳過期的 Session（nil 表示未過期或無事件）
func (m *Manager) CheckExpiry() *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session == nil {
		return nil
	}
	if m.session.IsActive() {
		return nil
	}
	// 已過期
	expired := m.session
	m.session = nil
	return expired
}

// GetSnapshot 取得當前快照
func (m *Manager) GetSnapshot() Snapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session == nil || !m.session.IsActive() {
		return Snapshot{IsActive: false}
	}
	return Snapshot{
		IsActive:      true,
		TierName:      m.session.Tier.Name,
		TierIcon:      m.session.Tier.Icon,
		TierColor:     m.session.Tier.Color,
		SpawnBoost:    m.session.Tier.SpawnMultBoost,
		RewardBoost:   m.session.Tier.RewardBoost,
		SecondsLeft:   m.session.SecondsLeft(),
		TotalDuration: m.session.Tier.Duration,
	}
}

// GetCooldownLeft 取得冷卻剩餘秒數
func (m *Manager) GetCooldownLeft() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	left := time.Until(m.cooldown)
	if left <= 0 {
		return 0
	}
	return int(left.Seconds())
}
