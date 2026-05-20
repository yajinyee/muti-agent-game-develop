// Package event 限時活動系統（DAY-079）
// 每 30 分鐘自動輪換活動，提供特殊加成效果
// 3 種活動類型：黃金時段/魚群爆發/幸運時刻
package event

import (
	"sync"
	"time"
)

// EventType 活動類型
type EventType string

const (
	EventGoldenHour   EventType = "golden_hour"   // 黃金時段：獎勵倍率 ×1.5
	EventFishFrenzy   EventType = "fish_frenzy"   // 魚群爆發：目標數量 ×2
	EventLuckyMoment  EventType = "lucky_moment"  // 幸運時刻：擊破率 +20%
	EventNone         EventType = "none"           // 無活動
)

// EventDef 活動定義
type EventDef struct {
	Type        EventType `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Color       string    `json:"color"`
	Duration    time.Duration `json:"-"` // 活動持續時間
	// 效果參數
	RewardMult    float64 `json:"reward_mult"`    // 獎勵倍率加成（1.0=無加成）
	SpawnMult     float64 `json:"spawn_mult"`     // 目標生成倍率（1.0=無加成）
	KillChanceAdd float64 `json:"kill_chance_add"` // 擊破率加成（0.0=無加成）
}

// EventDefs 活動定義列表
var EventDefs = map[EventType]*EventDef{
	EventGoldenHour: {
		Type:          EventGoldenHour,
		Name:          "黃金時段",
		Description:   "所有獎勵倍率提升 50%！",
		Icon:          "✨",
		Color:         "#FFD700",
		Duration:      30 * time.Minute,
		RewardMult:    1.5,
		SpawnMult:     1.0,
		KillChanceAdd: 0.0,
	},
	EventFishFrenzy: {
		Type:          EventFishFrenzy,
		Name:          "魚群爆發",
		Description:   "目標數量大幅增加，機會更多！",
		Icon:          "🐟",
		Color:         "#00BFFF",
		Duration:      30 * time.Minute,
		RewardMult:    1.0,
		SpawnMult:     2.0,
		KillChanceAdd: 0.0,
	},
	EventLuckyMoment: {
		Type:          EventLuckyMoment,
		Name:          "幸運時刻",
		Description:   "擊破率提升 20%，輕鬆大豐收！",
		Icon:          "🍀",
		Color:         "#00FF7F",
		Duration:      30 * time.Minute,
		RewardMult:    1.0,
		SpawnMult:     1.0,
		KillChanceAdd: 0.20,
	},
}

// EventRotation 活動輪換順序
var EventRotation = []EventType{
	EventGoldenHour,
	EventNone,
	EventFishFrenzy,
	EventNone,
	EventLuckyMoment,
	EventNone,
}

// ActiveEvent 當前活動狀態
type ActiveEvent struct {
	Type        EventType `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Color       string    `json:"color"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	// 效果參數
	RewardMult    float64 `json:"reward_mult"`
	SpawnMult     float64 `json:"spawn_mult"`
	KillChanceAdd float64 `json:"kill_chance_add"`
}

// IsActive 是否有活動進行中
func (e *ActiveEvent) IsActive() bool {
	if e == nil || e.Type == EventNone {
		return false
	}
	now := time.Now()
	return !now.Before(e.StartAt) && now.Before(e.EndAt)
}

// TimeLeft 剩餘時間（秒）
func (e *ActiveEvent) TimeLeft() float64 {
	if e == nil {
		return 0
	}
	remaining := time.Until(e.EndAt).Seconds()
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Manager 限時活動管理器
type Manager struct {
	mu           sync.RWMutex
	current      *ActiveEvent
	rotationIdx  int
	slotDuration time.Duration // 每個輪換槽的持續時間（包含無活動期）
}

// New 建立新的活動管理器
// slotDuration：每個輪換槽的持續時間（預設 30 分鐘）
func New(slotDuration time.Duration) *Manager {
	if slotDuration <= 0 {
		slotDuration = 30 * time.Minute
	}
	m := &Manager{
		slotDuration: slotDuration,
	}
	m.advance() // 初始化第一個活動
	return m
}

// advance 推進到下一個活動
func (m *Manager) advance() {
	eventType := EventRotation[m.rotationIdx%len(EventRotation)]
	m.rotationIdx++

	now := time.Now()
	endAt := now.Add(m.slotDuration)

	if eventType == EventNone {
		m.current = &ActiveEvent{
			Type:    EventNone,
			StartAt: now,
			EndAt:   endAt,
		}
		return
	}

	def, ok := EventDefs[eventType]
	if !ok {
		m.current = &ActiveEvent{
			Type:    EventNone,
			StartAt: now,
			EndAt:   endAt,
		}
		return
	}

	m.current = &ActiveEvent{
		Type:          def.Type,
		Name:          def.Name,
		Description:   def.Description,
		Icon:          def.Icon,
		Color:         def.Color,
		StartAt:       now,
		EndAt:         endAt,
		RewardMult:    def.RewardMult,
		SpawnMult:     def.SpawnMult,
		KillChanceAdd: def.KillChanceAdd,
	}
}

// Tick 定期呼叫，檢查是否需要切換活動
// 回傳 true 表示活動已切換（需要廣播）
func (m *Manager) Tick() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.current == nil || time.Now().After(m.current.EndAt) {
		m.advance()
		return true
	}
	return false
}

// GetCurrent 取得當前活動（thread-safe）
func (m *Manager) GetCurrent() *ActiveEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current == nil {
		return &ActiveEvent{Type: EventNone}
	}
	// 回傳副本
	cp := *m.current
	return &cp
}

// GetRewardMult 取得當前獎勵倍率加成（thread-safe）
func (m *Manager) GetRewardMult() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current == nil || !m.current.IsActive() {
		return 1.0
	}
	return m.current.RewardMult
}

// GetSpawnMult 取得當前目標生成倍率（thread-safe）
func (m *Manager) GetSpawnMult() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current == nil || !m.current.IsActive() {
		return 1.0
	}
	return m.current.SpawnMult
}

// GetKillChanceAdd 取得當前擊破率加成（thread-safe）
func (m *Manager) GetKillChanceAdd() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current == nil || !m.current.IsActive() {
		return 0.0
	}
	return m.current.KillChanceAdd
}

// GetSnapshot 取得活動快照（用於 WebSocket 廣播）
func (m *Manager) GetSnapshot() EventSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.current == nil {
		return EventSnapshot{Type: string(EventNone)}
	}

	return EventSnapshot{
		Type:          string(m.current.Type),
		Name:          m.current.Name,
		Description:   m.current.Description,
		Icon:          m.current.Icon,
		Color:         m.current.Color,
		IsActive:      m.current.IsActive(),
		EndAt:         m.current.EndAt.UnixMilli(),
		TimeLeft:      m.current.TimeLeft(),
		RewardMult:    m.current.RewardMult,
		SpawnMult:     m.current.SpawnMult,
		KillChanceAdd: m.current.KillChanceAdd,
	}
}

// EventSnapshot 活動快照（用於 WebSocket 廣播）
type EventSnapshot struct {
	Type          string  `json:"type"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Icon          string  `json:"icon"`
	Color         string  `json:"color"`
	IsActive      bool    `json:"is_active"`
	EndAt         int64   `json:"end_at"`    // Unix ms
	TimeLeft      float64 `json:"time_left"` // 秒
	RewardMult    float64 `json:"reward_mult"`
	SpawnMult     float64 `json:"spawn_mult"`
	KillChanceAdd float64 `json:"kill_chance_add"`
}
