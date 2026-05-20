// Package halloffame 全服名人堂系統（DAY-110）
// 追蹤並展示全服最佳記錄，激勵玩家挑戰極限
package halloffame

import (
	"sync"
	"time"
)

// RecordType 名人堂記錄類型
type RecordType string

const (
	RecordBestStreak      RecordType = "best_streak"       // 最高連擊
	RecordBestMultiplier  RecordType = "best_multiplier"   // 最高單次倍率
	RecordBestBonusReward RecordType = "best_bonus_reward" // 最高 Bonus 獎勵
	RecordMostJackpots    RecordType = "most_jackpots"     // 最多 Jackpot 次數
	RecordGrandJackpot    RecordType = "grand_jackpot"     // Grand Jackpot 中獎
	RecordBossKills       RecordType = "boss_kills"        // 最多 BOSS 擊殺
	RecordMaxCoins        RecordType = "max_coins"         // 歷史最高金幣
	RecordBestRTP         RecordType = "best_rtp"          // 最高 RTP（需 >= 100 次攻擊）
)

// HallEntry 名人堂條目
type HallEntry struct {
	PlayerID    string     `json:"player_id"`
	DisplayName string     `json:"display_name"`
	RecordType  RecordType `json:"record_type"`
	Value       float64    `json:"value"`       // 記錄數值（連擊數/倍率/金幣等）
	Description string     `json:"description"` // 人類可讀描述
	AchievedAt  time.Time  `json:"achieved_at"`
	BetLevel    int        `json:"bet_level"`   // 達成時的投注等級
	CharacterID int        `json:"character_id"` // 達成時使用的角色
}

// HallSnapshot 名人堂快照（用於傳送給 Client）
type HallSnapshot struct {
	Records   map[RecordType]*HallEntry `json:"records"`
	UpdatedAt int64                     `json:"updated_at_ms"`
}

// Manager 名人堂管理器
type Manager struct {
	mu      sync.RWMutex
	records map[RecordType]*HallEntry
}

// New 建立名人堂管理器
func New() *Manager {
	return &Manager{
		records: make(map[RecordType]*HallEntry),
	}
}

// TryUpdate 嘗試更新名人堂記錄
// 回傳 (isNewRecord bool, oldEntry *HallEntry)
func (m *Manager) TryUpdate(
	playerID, displayName string,
	recordType RecordType,
	value float64,
	description string,
	betLevel, characterID int,
) (bool, *HallEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()

	existing, ok := m.records[recordType]
	if ok && existing.Value >= value {
		return false, nil
	}

	oldEntry := existing
	newEntry := &HallEntry{
		PlayerID:    playerID,
		DisplayName: displayName,
		RecordType:  recordType,
		Value:       value,
		Description: description,
		AchievedAt:  time.Now(),
		BetLevel:    betLevel,
		CharacterID: characterID,
	}
	m.records[recordType] = newEntry
	return true, oldEntry
}

// GetAll 取得所有名人堂記錄
func (m *Manager) GetAll() HallSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	snap := HallSnapshot{
		Records:   make(map[RecordType]*HallEntry, len(m.records)),
		UpdatedAt: time.Now().UnixMilli(),
	}
	for k, v := range m.records {
		entry := *v // copy
		snap.Records[k] = &entry
	}
	return snap
}

// GetRecord 取得特定類型的記錄
func (m *Manager) GetRecord(rt RecordType) *HallEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if e, ok := m.records[rt]; ok {
		copy := *e
		return &copy
	}
	return nil
}

// IsRecordHolder 檢查玩家是否持有某項記錄
func (m *Manager) IsRecordHolder(playerID string, rt RecordType) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if e, ok := m.records[rt]; ok {
		return e.PlayerID == playerID
	}
	return false
}

// GetPlayerRecords 取得玩家持有的所有記錄
func (m *Manager) GetPlayerRecords(playerID string) []*HallEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*HallEntry
	for _, e := range m.records {
		if e.PlayerID == playerID {
			copy := *e
			result = append(result, &copy)
		}
	}
	return result
}

// RecordTypeLabel 取得記錄類型的中文標籤
func RecordTypeLabel(rt RecordType) string {
	switch rt {
	case RecordBestStreak:
		return "最高連擊王"
	case RecordBestMultiplier:
		return "最高倍率王"
	case RecordBestBonusReward:
		return "Bonus 大師"
	case RecordMostJackpots:
		return "Jackpot 收集者"
	case RecordGrandJackpot:
		return "Grand Jackpot 傳說"
	case RecordBossKills:
		return "BOSS 獵人"
	case RecordMaxCoins:
		return "金幣大亨"
	case RecordBestRTP:
		return "效率之王"
	default:
		return string(rt)
	}
}

// RecordTypeIcon 取得記錄類型的圖示
func RecordTypeIcon(rt RecordType) string {
	switch rt {
	case RecordBestStreak:
		return "🔥"
	case RecordBestMultiplier:
		return "⚡"
	case RecordBestBonusReward:
		return "🌾"
	case RecordMostJackpots:
		return "🎰"
	case RecordGrandJackpot:
		return "👑"
	case RecordBossKills:
		return "⚔️"
	case RecordMaxCoins:
		return "💰"
	case RecordBestRTP:
		return "📊"
	default:
		return "🏆"
	}
}
