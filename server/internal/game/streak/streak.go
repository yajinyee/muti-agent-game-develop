// Package streak 連擊系統（DAY-083）
// 玩家連續擊破目標時，連擊數遞增，獎勵倍率提升
// 超過 3 秒未擊破則連擊重置
package streak

import (
	"sync"
	"time"
)

// 連擊等級定義
type Level struct {
	MinStreak  int     // 最低連擊數
	MultBonus  float64 // 獎勵倍率加成（乘法，1.0 = 無加成）
	Name       string  // 等級名稱
	Color      string  // 顯示顏色（hex）
}

// Levels 連擊等級表
var Levels = []Level{
	{MinStreak: 1,  MultBonus: 1.0,  Name: "開始",   Color: "#FFFFFF"},
	{MinStreak: 3,  MultBonus: 1.1,  Name: "連擊！",  Color: "#FFFF00"},
	{MinStreak: 5,  MultBonus: 1.2,  Name: "熱身中",  Color: "#FFA500"},
	{MinStreak: 8,  MultBonus: 1.35, Name: "火力全開", Color: "#FF6600"},
	{MinStreak: 12, MultBonus: 1.5,  Name: "無法阻擋", Color: "#FF0000"},
	{MinStreak: 20, MultBonus: 2.0,  Name: "傳說連擊", Color: "#FF00FF"},
}

// ResetTimeout 連擊重置超時（秒）
const ResetTimeout = 3 * time.Second

// Manager 連擊管理器（per-player）
type Manager struct {
	mu          sync.Mutex
	current     int       // 當前連擊數
	maxStreak   int       // 本局最高連擊數
	lastKillAt  time.Time // 最後一次擊破時間
	totalKills  int       // 本局總擊破數
}

// NewManager 建立新連擊管理器
func NewManager() *Manager {
	return &Manager{}
}

// RecordKill 記錄一次擊破，回傳 (currentStreak, multBonus, isNewLevel)
// isNewLevel = 這次擊破剛好升到新的連擊等級
func (m *Manager) RecordKill() (currentStreak int, multBonus float64, isNewLevel bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	// 檢查是否超時重置
	if m.current > 0 && now.Sub(m.lastKillAt) > ResetTimeout {
		m.current = 0
	}

	prevLevel := m.getLevel(m.current)
	m.current++
	m.totalKills++
	m.lastKillAt = now

	if m.current > m.maxStreak {
		m.maxStreak = m.current
	}

	newLevel := m.getLevel(m.current)
	isNewLevel = newLevel.MinStreak > prevLevel.MinStreak && m.current == newLevel.MinStreak

	return m.current, newLevel.MultBonus, isNewLevel
}

// CheckTimeout 檢查是否超時，若超時則重置連擊（由外部定期呼叫）
// 回傳是否發生了重置
func (m *Manager) CheckTimeout() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.current > 0 && time.Since(m.lastKillAt) > ResetTimeout {
		m.current = 0
		return true
	}
	return false
}

// Reset 強制重置連擊（如玩家離線）
func (m *Manager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.current = 0
}

// GetSnapshot 取得當前快照
func (m *Manager) GetSnapshot() Snapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	level := m.getLevel(m.current)
	return Snapshot{
		Current:   m.current,
		MaxStreak: m.maxStreak,
		MultBonus: level.MultBonus,
		LevelName: level.Name,
		LevelColor: level.Color,
		TotalKills: m.totalKills,
	}
}

// getLevel 取得指定連擊數對應的等級（需在鎖內呼叫）
func (m *Manager) getLevel(streak int) Level {
	result := Levels[0]
	for _, lv := range Levels {
		if streak >= lv.MinStreak {
			result = lv
		}
	}
	return result
}

// Snapshot 連擊快照
type Snapshot struct {
	Current    int     `json:"current"`
	MaxStreak  int     `json:"max_streak"`
	MultBonus  float64 `json:"mult_bonus"`
	LevelName  string  `json:"level_name"`
	LevelColor string  `json:"level_color"`
	TotalKills int     `json:"total_kills"`
}
