// Package goldentime — 黃金時間系統（DAY-125）
// 業界依據：Fire Kirin / Ocean King 系列的 Golden Time 機制
// 全場目標物倍率暫時提升，製造「全場瘋狂」的高峰體驗
// 業界研究顯示 Golden Time 讓短期參與度提升 40%+
package goldentime

import (
	"sync"
	"time"
)

// TierDef 黃金時間等級定義
type TierDef struct {
	Name       string  // 等級名稱
	MultBoost  float64 // 倍率加成（1.5 = 全場 ×1.5）
	Duration   int     // 持續秒數
	Icon       string  // 圖示
	Color      string  // 顏色（hex）
	BgColor    string  // 背景顏色
}

// Tier 黃金時間等級
type Tier int

const (
	TierSilver Tier = iota // 銀色時間（×1.5，30秒）
	TierGold               // 黃金時間（×2.0，45秒）
	TierRainbow            // 彩虹時間（×3.0，60秒）
)

// TierDefs 等級定義表
var TierDefs = map[Tier]TierDef{
	TierSilver: {
		Name:      "⚡ 銀色時間",
		MultBoost: 1.5,
		Duration:  30,
		Icon:      "⚡",
		Color:     "#C0C0C0",
		BgColor:   "#2A2A3A",
	},
	TierGold: {
		Name:      "✨ 黃金時間",
		MultBoost: 2.0,
		Duration:  45,
		Icon:      "✨",
		Color:     "#FFD700",
		BgColor:   "#3A2A00",
	},
	TierRainbow: {
		Name:      "🌈 彩虹時間",
		MultBoost: 3.0,
		Duration:  60,
		Icon:      "🌈",
		Color:     "#FF69B4",
		BgColor:   "#2A003A",
	},
}

// TriggerType 觸發類型
type TriggerType string

const (
	TriggerBossKill    TriggerType = "boss_kill"    // BOSS 擊殺後觸發
	TriggerRandom      TriggerType = "random"       // 隨機觸發
	TriggerFlashCombo  TriggerType = "flash_combo"  // 閃電挑戰完成後觸發
	TriggerRaidVictory TriggerType = "raid_victory" // Raid 勝利後觸發
)

// Session 黃金時間 session
type Session struct {
	Tier        Tier
	TriggerType TriggerType
	StartedAt   time.Time
	EndsAt      time.Time
	MultBoost   float64
}

// IsActive 是否仍在進行中
func (s *Session) IsActive() bool {
	return time.Now().Before(s.EndsAt)
}

// SecondsLeft 剩餘秒數
func (s *Session) SecondsLeft() int {
	d := time.Until(s.EndsAt)
	if d < 0 {
		return 0
	}
	return int(d.Seconds())
}

// Snapshot 快照
type Snapshot struct {
	IsActive    bool
	Tier        int
	TierName    string
	MultBoost   float64
	SecondsLeft int
	Icon        string
	Color       string
	BgColor     string
	TriggerType string
}

// Manager 黃金時間管理器
type Manager struct {
	mu      sync.RWMutex
	session *Session
	cooldown time.Time // 冷卻結束時間（防止連續觸發）
}

// New 建立新管理器
func New() *Manager {
	return &Manager{}
}

// CanTrigger 是否可以觸發（冷卻檢查）
func (m *Manager) CanTrigger() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.session != nil && m.session.IsActive() {
		return false // 已有進行中的 session
	}
	return time.Now().After(m.cooldown)
}

// Start 開始黃金時間
// 回傳 session 快照，nil 表示無法觸發
func (m *Manager) Start(tier Tier, trigger TriggerType) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 已有進行中的 session，不重複觸發
	if m.session != nil && m.session.IsActive() {
		return nil
	}
	// 冷卻中
	if time.Now().Before(m.cooldown) {
		return nil
	}

	def := TierDefs[tier]
	now := time.Now()
	s := &Session{
		Tier:        tier,
		TriggerType: trigger,
		StartedAt:   now,
		EndsAt:      now.Add(time.Duration(def.Duration) * time.Second),
		MultBoost:   def.MultBoost,
	}
	m.session = s

	// 設定冷卻：黃金時間結束後 3 分鐘才能再次觸發
	m.cooldown = s.EndsAt.Add(3 * time.Minute)

	return s
}

// GetMultBoost 取得當前倍率加成（1.0 = 無加成）
func (m *Manager) GetMultBoost() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.session == nil || !m.session.IsActive() {
		return 1.0
	}
	return m.session.MultBoost
}

// IsActive 是否有進行中的黃金時間
func (m *Manager) IsActive() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.session != nil && m.session.IsActive()
}

// CheckExpiry 檢查是否剛剛結束（由 game loop 呼叫）
// 回傳 true 表示剛剛結束，需要廣播結束事件
func (m *Manager) CheckExpiry() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.session == nil {
		return false
	}
	if !m.session.IsActive() {
		m.session = nil
		return true
	}
	return false
}

// GetSnapshot 取得當前快照
func (m *Manager) GetSnapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.session == nil || !m.session.IsActive() {
		return Snapshot{IsActive: false}
	}

	def := TierDefs[m.session.Tier]
	return Snapshot{
		IsActive:    true,
		Tier:        int(m.session.Tier),
		TierName:    def.Name,
		MultBoost:   m.session.MultBoost,
		SecondsLeft: m.session.SecondsLeft(),
		Icon:        def.Icon,
		Color:       def.Color,
		BgColor:     def.BgColor,
		TriggerType: string(m.session.TriggerType),
	}
}

// ShouldTriggerRandom 隨機觸發判斷（由 game loop 呼叫）
// 每次呼叫有 0.5% 機率觸發（約每 200 次 tick = 每 200 秒觸發一次）
func (m *Manager) ShouldTriggerRandom() bool {
	if !m.CanTrigger() {
		return false
	}
	// 0.5% 機率
	return randFloat() < 0.005
}

// SelectTier 根據觸發類型選擇等級
// boss_kill → 70% Gold + 30% Rainbow
// random → 60% Silver + 30% Gold + 10% Rainbow
// flash_combo → 50% Gold + 50% Rainbow
// raid_victory → 100% Rainbow
func SelectTier(trigger TriggerType) Tier {
	r := randFloat()
	switch trigger {
	case TriggerBossKill:
		if r < 0.70 {
			return TierGold
		}
		return TierRainbow
	case TriggerFlashCombo:
		if r < 0.50 {
			return TierGold
		}
		return TierRainbow
	case TriggerRaidVictory:
		return TierRainbow
	default: // TriggerRandom
		if r < 0.60 {
			return TierSilver
		} else if r < 0.90 {
			return TierGold
		}
		return TierRainbow
	}
}
