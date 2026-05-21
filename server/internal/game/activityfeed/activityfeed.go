// Package activityfeed 成就動態牆系統（DAY-112）
// 記錄全服玩家的重要事件，製造社交證明和 FOMO 效應
package activityfeed

import (
	"fmt"
	"sync"
	"time"
)

// EventType 動態事件類型
type EventType string

const (
	EventAchievement  EventType = "achievement"   // 成就解鎖
	EventTitle        EventType = "title"          // 稱號獲得
	EventJackpot      EventType = "jackpot"        // Jackpot 中獎
	EventBossKill     EventType = "boss_kill"      // BOSS 擊殺
	EventMegaWin      EventType = "mega_win"       // 超大獎（≥50x）
	EventStreakRecord EventType = "streak_record"  // 連擊記錄（≥20）
	EventHallOfFame   EventType = "hall_of_fame"   // 名人堂新記錄
	EventSeasonLevel  EventType = "season_level"   // 賽季升級
	EventMilestone    EventType = "milestone"      // 登入里程碑
	EventLightningChain EventType = "lightning_chain" // 閃電鰻連鎖擊破（DAY-132）
)

// Rarity 事件稀有度（影響 UI 顯示顏色）
type Rarity string

const (
	RarityCommon    Rarity = "common"    // 灰色
	RarityUncommon  Rarity = "uncommon"  // 綠色
	RarityRare      Rarity = "rare"      // 藍色
	RarityEpic      Rarity = "epic"      // 紫色
	RarityLegendary Rarity = "legendary" // 金色
)

// FeedEvent 動態牆事件
type FeedEvent struct {
	ID          string    `json:"id"`           // 唯一 ID（時間戳+序號）
	EventType   EventType `json:"event_type"`   // 事件類型
	PlayerID    string    `json:"player_id"`    // 玩家 ID
	DisplayName string    `json:"display_name"` // 玩家顯示名稱
	Icon        string    `json:"icon"`         // 事件圖示（emoji）
	Title       string    `json:"title"`        // 事件標題（如「解鎖成就」）
	Detail      string    `json:"detail"`       // 事件詳情（如「討伐傳說」）
	Rarity      Rarity    `json:"rarity"`       // 稀有度
	Timestamp   int64     `json:"timestamp"`    // Unix ms
}

// Manager 動態牆管理器
type Manager struct {
	mu      sync.RWMutex
	events  []*FeedEvent
	maxSize int
	counter int64 // 事件序號（用於生成唯一 ID）
}

// New 建立新的動態牆管理器
func New() *Manager {
	return &Manager{
		events:  make([]*FeedEvent, 0, 50),
		maxSize: 50,
	}
}

// Push 推入新事件，回傳事件（供廣播用）
func (m *Manager) Push(evt *FeedEvent) *FeedEvent {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter++
	evt.ID = generateID(evt.Timestamp, m.counter)
	evt.Timestamp = time.Now().UnixMilli()

	m.events = append(m.events, evt)
	if len(m.events) > m.maxSize {
		m.events = m.events[len(m.events)-m.maxSize:]
	}
	return evt
}

// GetRecent 取得最近 N 條事件（由新到舊）
func (m *Manager) GetRecent(n int) []*FeedEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := len(m.events)
	if n > total {
		n = total
	}
	result := make([]*FeedEvent, n)
	// 從最新的開始取
	for i := 0; i < n; i++ {
		result[i] = m.events[total-1-i]
	}
	return result
}

// generateID 生成唯一事件 ID
func generateID(ts int64, counter int64) string {
	return fmt.Sprintf("feed_%d_%d", ts, counter)
}

// ---- 事件建構輔助函數 ----

// NewAchievementEvent 建立成就解鎖事件
func NewAchievementEvent(playerID, displayName, achName, achIcon, achType string) *FeedEvent {
	rarity := rarityFromAchType(achType)
	return &FeedEvent{
		EventType:   EventAchievement,
		PlayerID:    playerID,
		DisplayName: displayName,
		Icon:        achIcon,
		Title:       displayName + " 解鎖成就",
		Detail:      achName,
		Rarity:      rarity,
	}
}

// NewTitleEvent 建立稱號獲得事件
func NewTitleEvent(playerID, displayName, titleName, titleIcon string, priority int) *FeedEvent {
	rarity := rarityFromTitlePriority(priority)
	return &FeedEvent{
		EventType:   EventTitle,
		PlayerID:    playerID,
		DisplayName: displayName,
		Icon:        titleIcon,
		Title:       displayName + " 獲得稱號",
		Detail:      titleName,
		Rarity:      rarity,
	}
}

// NewJackpotEvent 建立 Jackpot 中獎事件
func NewJackpotEvent(playerID, displayName, levelName, levelIcon string, amount int) *FeedEvent {
	rarity := RarityLegendary
	if levelName == "Mini" || levelName == "Minor" {
		rarity = RarityRare
	} else if levelName == "Major" {
		rarity = RarityEpic
	}
	return &FeedEvent{
		EventType:   EventJackpot,
		PlayerID:    playerID,
		DisplayName: displayName,
		Icon:        levelIcon,
		Title:       displayName + " 中了 " + levelName + " Jackpot",
		Detail:      fmt.Sprintf("%d 金幣", amount),
		Rarity:      rarity,
	}
}

// NewBossKillEvent 建立 BOSS 擊殺事件
func NewBossKillEvent(playerID, displayName string, reward int) *FeedEvent {
	return &FeedEvent{
		EventType:   EventBossKill,
		PlayerID:    playerID,
		DisplayName: displayName,
		Icon:        "🔥",
		Title:       displayName + " 擊敗了 BOSS",
		Detail:      fmt.Sprintf("獲得 %d 金幣", reward),
		Rarity:      RarityEpic,
	}
}

// NewMegaWinEvent 建立超大獎事件（≥50x）
func NewMegaWinEvent(playerID, displayName string, multiplier float64, reward int) *FeedEvent {
	rarity := RarityRare
	if multiplier >= 100 {
		rarity = RarityEpic
	}
	if multiplier >= 200 {
		rarity = RarityLegendary
	}
	return &FeedEvent{
		EventType:   EventMegaWin,
		PlayerID:    playerID,
		DisplayName: displayName,
		Icon:        "💎",
		Title:       displayName + " 獲得超大獎",
		Detail:      fmt.Sprintf("%.0fx = %d 金幣", multiplier, reward),
		Rarity:      rarity,
	}
}

// NewStreakRecordEvent 建立連擊記錄事件
func NewStreakRecordEvent(playerID, displayName string, streak int, levelName string) *FeedEvent {
	return &FeedEvent{
		EventType:   EventStreakRecord,
		PlayerID:    playerID,
		DisplayName: displayName,
		Icon:        "⚡",
		Title:       displayName + " 達成連擊記錄",
		Detail:      fmt.Sprintf("%d 連擊（%s）", streak, levelName),
		Rarity:      RarityUncommon,
	}
}

// NewHallOfFameEvent 建立名人堂新記錄事件
func NewHallOfFameEvent(playerID, displayName, recordType, recordLabel string) *FeedEvent {
	return &FeedEvent{
		EventType:   EventHallOfFame,
		PlayerID:    playerID,
		DisplayName: displayName,
		Icon:        "🏆",
		Title:       displayName + " 創下名人堂記錄",
		Detail:      recordLabel,
		Rarity:      RarityLegendary,
	}
}

// NewSeasonLevelEvent 建立賽季升級事件
func NewSeasonLevelEvent(playerID, displayName string, level int) *FeedEvent {
	rarity := RarityCommon
	if level >= 20 {
		rarity = RarityUncommon
	}
	if level >= 50 {
		rarity = RarityRare
	}
	if level >= 80 {
		rarity = RarityEpic
	}
	return &FeedEvent{
		EventType:   EventSeasonLevel,
		PlayerID:    playerID,
		DisplayName: displayName,
		Icon:        "🌟",
		Title:       displayName + " 賽季升級",
		Detail:      fmt.Sprintf("達到第 %d 級", level),
		Rarity:      rarity,
	}
}

// NewMilestoneEvent 建立登入里程碑事件
func NewMilestoneEvent(playerID, displayName string, days int, milestoneName string) *FeedEvent {
	rarity := RarityUncommon
	if days >= 30 {
		rarity = RarityRare
	}
	if days >= 60 {
		rarity = RarityEpic
	}
	if days >= 100 {
		rarity = RarityLegendary
	}
	return &FeedEvent{
		EventType:   EventMilestone,
		PlayerID:    playerID,
		DisplayName: displayName,
		Icon:        "📅",
		Title:       displayName + " 達成登入里程碑",
		Detail:      fmt.Sprintf("%d 天 — %s", days, milestoneName),
		Rarity:      rarity,
	}
}

// ---- 稀有度輔助函數 ----

func rarityFromAchType(achType string) Rarity {
	switch achType {
	case "special":
		return RarityEpic
	case "boss":
		return RarityRare
	case "bonus":
		return RarityUncommon
	default:
		return RarityCommon
	}
}

func rarityFromTitlePriority(priority int) Rarity {
	switch {
	case priority >= 80:
		return RarityLegendary
	case priority >= 60:
		return RarityEpic
	case priority >= 40:
		return RarityRare
	case priority >= 20:
		return RarityUncommon
	default:
		return RarityCommon
	}
}
