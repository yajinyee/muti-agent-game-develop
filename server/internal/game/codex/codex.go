// Package codex 魚類圖鑑收集系統（DAY-081）
// 玩家首次擊破每種目標物時解鎖圖鑑條目，記錄擊破次數和最高倍率
package codex

import (
	"sync"
	"time"
)

// Entry 圖鑑條目
type Entry struct {
	TargetID    string    `json:"target_id"`
	TargetName  string    `json:"target_name"`
	Unlocked    bool      `json:"unlocked"`
	UnlockedAt  time.Time `json:"unlocked_at,omitempty"`
	KillCount   int       `json:"kill_count"`
	MaxMultiplier float64 `json:"max_multiplier"`
	Rarity      string    `json:"rarity"` // common/rare/epic/legendary
}

// Manager 圖鑑管理器
type Manager struct {
	mu      sync.RWMutex
	entries map[string]*Entry // targetID -> Entry
}

// AllTargetIDs 所有可收集的目標物 ID（依稀有度排序）
var AllTargetIDs = []string{
	"T001", "T002", "T003", "T004", "T005", "T006", // 普通
	"T101", "T102", "T103", "T104", "T105",          // 特殊
	"B001",                                           // BOSS
}

// targetMeta 目標物元資料（名稱 + 稀有度）
var targetMeta = map[string]struct {
	Name   string
	Rarity string
}{
	"T001": {"像素雜草", "common"},
	"T002": {"綠色小蟲", "common"},
	"T003": {"紅色小蟲", "common"},
	"T004": {"藍色小蟲", "common"},
	"T005": {"會走路的布丁", "common"},
	"T006": {"巨大蘑菇", "common"},
	"T101": {"擬態型怪物", "rare"},
	"T102": {"寶箱怪", "rare"},
	"T103": {"流星", "epic"},
	"T104": {"金色雜草", "epic"},
	"T105": {"巨大金幣魚", "epic"},
	"B001": {"那個孩子", "legendary"},
}

// UnlockReward 解鎖單一條目的金幣獎勵
const UnlockReward = 200

// CompleteReward 全圖鑑完成的金幣獎勵
const CompleteReward = 5000

// NewManager 建立新圖鑑管理器
func NewManager() *Manager {
	m := &Manager{
		entries: make(map[string]*Entry),
	}
	// 初始化所有條目（未解鎖狀態）
	for _, id := range AllTargetIDs {
		meta := targetMeta[id]
		m.entries[id] = &Entry{
			TargetID:   id,
			TargetName: meta.Name,
			Unlocked:   false,
			Rarity:     meta.Rarity,
		}
	}
	return m
}

// RecordKill 記錄擊破事件，回傳 (isNewUnlock, isComplete)
// isNewUnlock = 首次解鎖此條目
// isComplete = 全圖鑑完成（剛好這次觸發）
func (m *Manager) RecordKill(targetID string, multiplier float64) (isNewUnlock bool, isComplete bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry, ok := m.entries[targetID]
	if !ok {
		return false, false
	}

	entry.KillCount++
	if multiplier > entry.MaxMultiplier {
		entry.MaxMultiplier = multiplier
	}

	if !entry.Unlocked {
		entry.Unlocked = true
		entry.UnlockedAt = time.Now()
		isNewUnlock = true

		// 檢查是否全圖鑑完成
		isComplete = m.isAllUnlocked()
	}

	return isNewUnlock, isComplete
}

// isAllUnlocked 檢查是否所有條目都已解鎖（需在鎖內呼叫）
func (m *Manager) isAllUnlocked() bool {
	for _, e := range m.entries {
		if !e.Unlocked {
			return false
		}
	}
	return true
}

// GetSnapshot 取得圖鑑快照（用於 WebSocket 廣播）
func (m *Manager) GetSnapshot() []*Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*Entry, 0, len(AllTargetIDs))
	for _, id := range AllTargetIDs {
		if e, ok := m.entries[id]; ok {
			// 複製一份避免外部修改
			copy := *e
			result = append(result, &copy)
		}
	}
	return result
}

// GetStats 取得圖鑑統計
func (m *Manager) GetStats() (unlocked int, total int) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total = len(AllTargetIDs)
	for _, e := range m.entries {
		if e.Unlocked {
			unlocked++
		}
	}
	return unlocked, total
}

// IsComplete 是否已完成全圖鑑
func (m *Manager) IsComplete() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isAllUnlocked()
}

// LoadState 從持久化資料恢復圖鑑狀態
func (m *Manager) LoadState(entries []*Entry) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, saved := range entries {
		if e, ok := m.entries[saved.TargetID]; ok {
			e.Unlocked = saved.Unlocked
			e.UnlockedAt = saved.UnlockedAt
			e.KillCount = saved.KillCount
			e.MaxMultiplier = saved.MaxMultiplier
		}
	}
}
