// Package announce 全服公告系統（DAY-097）
// 當重大事件發生時，廣播全服通知，增加社交感和緊張感
package announce

import (
	"fmt"
	"sync"
	"time"
)

// EventType 公告事件類型
type EventType string

const (
	// 玩家事件
	EventJackpotWin    EventType = "jackpot_win"    // Jackpot 中獎
	EventBigWin        EventType = "big_win"         // 大獎（≥50x）
	EventMegaWin       EventType = "mega_win"        // 超大獎（≥100x）
	EventBossKill      EventType = "boss_kill"       // BOSS 擊殺
	EventStreakRecord  EventType = "streak_record"   // 連擊新記錄（≥20）
	EventPlayerJoin    EventType = "player_join"     // 玩家加入
	EventPlayerLeave   EventType = "player_leave"    // 玩家離開

	// 系統事件
	EventWeatherChange EventType = "weather_change"  // 天氣變化
	EventEventStart    EventType = "event_start"     // 限時活動開始
	EventDailyReset    EventType = "daily_reset"     // 每日重置
	EventBossWarning   EventType = "boss_warning"    // BOSS 即將出現
	EventGrandJackpot  EventType = "grand_jackpot"   // Grand Jackpot 中獎（最高優先）
)

// Priority 公告優先級
type Priority int

const (
	PriorityLow    Priority = 1 // 低優先（玩家加入/離開）
	PriorityNormal Priority = 2 // 普通（大獎/天氣）
	PriorityHigh   Priority = 3 // 高優先（BOSS/Jackpot）
	PriorityCritical Priority = 4 // 最高（Grand Jackpot）
)

// Announcement 單筆公告
type Announcement struct {
	ID         string    `json:"id"`
	EventType  EventType `json:"event_type"`
	Priority   Priority  `json:"priority"`
	Title      string    `json:"title"`       // 公告標題
	Message    string    `json:"message"`     // 公告內容
	PlayerName string    `json:"player_name"` // 相關玩家名稱（可空）
	Amount     int       `json:"amount"`      // 相關金額（可空）
	Icon       string    `json:"icon"`        // 顯示圖示
	Color      string    `json:"color"`       // 顯示顏色
	Duration   int       `json:"duration"`    // 顯示時長（毫秒）
	CreatedAt  time.Time `json:"created_at"`
	CreatedAtMs int64    `json:"created_at_ms"`
}

// Manager 公告管理器
type Manager struct {
	mu      sync.RWMutex
	history []Announcement
	maxSize int
	counter int
}

// NewManager 建立公告管理器
func NewManager() *Manager {
	return &Manager{
		history: make([]Announcement, 0, 50),
		maxSize: 50,
	}
}

// Create 建立一筆公告
func (m *Manager) Create(eventType EventType, playerName string, amount int, extra map[string]string) Announcement {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter++
	now := time.Now()

	title, message, icon, color, priority, duration := m.buildContent(eventType, playerName, amount, extra)

	ann := Announcement{
		ID:          fmt.Sprintf("ann_%d_%d", now.UnixMilli(), m.counter),
		EventType:   eventType,
		Priority:    priority,
		Title:       title,
		Message:     message,
		PlayerName:  playerName,
		Amount:      amount,
		Icon:        icon,
		Color:       color,
		Duration:    duration,
		CreatedAt:   now,
		CreatedAtMs: now.UnixMilli(),
	}

	// 加入歷史
	m.history = append([]Announcement{ann}, m.history...)
	if len(m.history) > m.maxSize {
		m.history = m.history[:m.maxSize]
	}

	return ann
}

// buildContent 根據事件類型建立公告內容
func (m *Manager) buildContent(eventType EventType, playerName string, amount int, extra map[string]string) (title, message, icon, color string, priority Priority, duration int) {
	name := playerName
	if name == "" {
		name = "玩家"
	}

	switch eventType {
	case EventGrandJackpot:
		title = "👑 GRAND JACKPOT！"
		message = fmt.Sprintf("%s 中了 Grand Jackpot！獲得 🪙%d！", name, amount)
		icon = "👑"
		color = "#FF0080"
		priority = PriorityCritical
		duration = 8000

	case EventJackpotWin:
		level := "JACKPOT"
		if extra != nil {
			if l, ok := extra["level_name"]; ok {
				level = l
			}
		}
		title = fmt.Sprintf("🎰 %s 中獎！", level)
		message = fmt.Sprintf("%s 中了 %s！獲得 🪙%d！", name, level, amount)
		icon = "🎰"
		color = "#FFD700"
		priority = PriorityHigh
		duration = 5000

	case EventMegaWin:
		mult := ""
		if extra != nil {
			if m, ok := extra["multiplier"]; ok {
				mult = m + "x "
			}
		}
		title = "🌟 MEGA WIN！"
		message = fmt.Sprintf("%s 獲得 %s超大獎！🪙%d！", name, mult, amount)
		icon = "🌟"
		color = "#FF6B35"
		priority = PriorityHigh
		duration = 5000

	case EventBigWin:
		title = "✨ BIG WIN！"
		message = fmt.Sprintf("%s 獲得大獎！🪙%d！", name, amount)
		icon = "✨"
		color = "#FFD700"
		priority = PriorityNormal
		duration = 3500

	case EventBossKill:
		bossName := "BOSS"
		if extra != nil {
			if b, ok := extra["boss_name"]; ok {
				bossName = b
			}
		}
		title = "⚔️ BOSS 擊殺！"
		message = fmt.Sprintf("%s 擊敗了 %s！獲得 🪙%d！", name, bossName, amount)
		icon = "⚔️"
		color = "#FF4444"
		priority = PriorityHigh
		duration = 5000

	case EventStreakRecord:
		title = "🔥 連擊記錄！"
		message = fmt.Sprintf("%s 達成 %d 連擊！", name, amount)
		icon = "🔥"
		color = "#FF8C00"
		priority = PriorityNormal
		duration = 3000

	case EventPlayerJoin:
		title = "👋 玩家加入"
		message = fmt.Sprintf("%s 加入了遊戲！", name)
		icon = "👋"
		color = "#4CAF50"
		priority = PriorityLow
		duration = 2500

	case EventPlayerLeave:
		title = "👋 玩家離開"
		message = fmt.Sprintf("%s 離開了遊戲。", name)
		icon = "🚪"
		color = "#9E9E9E"
		priority = PriorityLow
		duration = 2000

	case EventWeatherChange:
		weatherName := "天氣變化"
		if extra != nil {
			if w, ok := extra["weather_name"]; ok {
				weatherName = w
			}
		}
		title = "🌤️ 天氣變化"
		message = fmt.Sprintf("天氣變為「%s」！", weatherName)
		icon = "🌤️"
		color = "#64B5F6"
		priority = PriorityNormal
		duration = 3000

	case EventEventStart:
		eventName := "限時活動"
		if extra != nil {
			if e, ok := extra["event_name"]; ok {
				eventName = e
			}
		}
		title = "⚡ 限時活動開始！"
		message = fmt.Sprintf("「%s」開始了！", eventName)
		icon = "⚡"
		color = "#FFC107"
		priority = PriorityNormal
		duration = 4000

	case EventBossWarning:
		title = "⚠️ BOSS 即將出現！"
		message = "強大的 BOSS 即將現身！準備好了嗎？"
		icon = "⚠️"
		color = "#FF5722"
		priority = PriorityHigh
		duration = 4000

	case EventDailyReset:
		title = "🌅 每日重置"
		message = "新的一天開始了！任務和獎勵已重置。"
		icon = "🌅"
		color = "#81C784"
		priority = PriorityNormal
		duration = 3000

	default:
		title = "📢 公告"
		message = "系統公告"
		icon = "📢"
		color = "#FFFFFF"
		priority = PriorityLow
		duration = 2500
	}

	return
}

// GetRecent 取得最近 n 筆公告
func (m *Manager) GetRecent(n int) []Announcement {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if n > len(m.history) {
		n = len(m.history)
	}
	result := make([]Announcement, n)
	copy(result, m.history[:n])
	return result
}

// Count 取得公告總數
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.history)
}
