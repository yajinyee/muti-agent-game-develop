// Package season 賽季通行證系統（DAY-072）
// 賽季積分 = 週賽積分累積（跨週不重置）
// 10 個等級，每級 100 積分，每級有金幣獎勵
// 等級 5：解鎖「賽季限定」皮膚（season_gold）
// 等級 10：解鎖「賽季傳說」稱號
package season

import (
	"sync"
	"time"
)

// SeasonLevel 賽季等級定義
type SeasonLevel struct {
	Level       int    `json:"level"`
	PointsNeeded int   `json:"points_needed"` // 達到此等級所需累積積分
	CoinReward  int    `json:"coin_reward"`   // 金幣獎勵
	SpecialType string `json:"special_type"`  // "" / "skin" / "title"
	SpecialID   string `json:"special_id"`    // 特殊獎勵 ID
	SpecialName string `json:"special_name"`  // 特殊獎勵名稱
	Icon        string `json:"icon"`
}

// SeasonLevels 賽季等級定義（10 個等級）
var SeasonLevels = []SeasonLevel{
	{Level: 1, PointsNeeded: 100, CoinReward: 500, Icon: "⭐"},
	{Level: 2, PointsNeeded: 200, CoinReward: 800, Icon: "⭐⭐"},
	{Level: 3, PointsNeeded: 350, CoinReward: 1200, Icon: "⭐⭐⭐"},
	{Level: 4, PointsNeeded: 550, CoinReward: 1800, Icon: "🌟"},
	{Level: 5, PointsNeeded: 800, CoinReward: 2500,
		SpecialType: "skin", SpecialID: "season_gold", SpecialName: "賽季黃金砲台",
		Icon: "🌟🌟"},
	{Level: 6, PointsNeeded: 1100, CoinReward: 3000, Icon: "💫"},
	{Level: 7, PointsNeeded: 1500, CoinReward: 3500, Icon: "💫💫"},
	{Level: 8, PointsNeeded: 2000, CoinReward: 4000, Icon: "✨"},
	{Level: 9, PointsNeeded: 2600, CoinReward: 4500, Icon: "✨✨"},
	{Level: 10, PointsNeeded: 3300, CoinReward: 5000,
		SpecialType: "title", SpecialID: "season_legend", SpecialName: "賽季傳說",
		Icon: "👑"},
}

// PlayerSeasonData 玩家賽季資料
type PlayerSeasonData struct {
	PlayerID      string    `json:"player_id"`
	SeasonPoints  int       `json:"season_points"`   // 累積賽季積分
	CurrentLevel  int       `json:"current_level"`   // 當前等級（0=未達等級1）
	ClaimedLevels []int     `json:"claimed_levels"`  // 已領取獎勵的等級列表
	LastUpdated   time.Time `json:"last_updated"`
}

// LevelUpResult 升級結果
type LevelUpResult struct {
	NewLevel    int
	CoinReward  int
	SpecialType string
	SpecialID   string
	SpecialName string
}

// Manager 賽季通行證管理器
type Manager struct {
	mu      sync.RWMutex
	players map[string]*PlayerSeasonData // playerID → data
}

// New 建立新的賽季管理器
func New() *Manager {
	return &Manager{
		players: make(map[string]*PlayerSeasonData),
	}
}

// GetOrCreate 取得或建立玩家賽季資料
func (m *Manager) GetOrCreate(playerID string) *PlayerSeasonData {
	m.mu.Lock()
	defer m.mu.Unlock()

	if data, ok := m.players[playerID]; ok {
		return data
	}
	data := &PlayerSeasonData{
		PlayerID:      playerID,
		SeasonPoints:  0,
		CurrentLevel:  0,
		ClaimedLevels: []int{},
		LastUpdated:   time.Now(),
	}
	m.players[playerID] = data
	return data
}

// AddPoints 增加賽季積分，回傳新積分和可領取的等級列表
func (m *Manager) AddPoints(playerID string, points int) (newTotal int, newLevels []int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, ok := m.players[playerID]
	if !ok {
		data = &PlayerSeasonData{
			PlayerID:      playerID,
			SeasonPoints:  0,
			CurrentLevel:  0,
			ClaimedLevels: []int{},
			LastUpdated:   time.Now(),
		}
		m.players[playerID] = data
	}

	data.SeasonPoints += points
	data.LastUpdated = time.Now()

	// 檢查是否有新等級可領取
	newLevels = []int{}
	for _, lvl := range SeasonLevels {
		if data.SeasonPoints >= lvl.PointsNeeded && !m.hasClaimed(data, lvl.Level) {
			newLevels = append(newLevels, lvl.Level)
		}
	}

	return data.SeasonPoints, newLevels
}

// ClaimLevel 領取等級獎勵，回傳獎勵資訊
func (m *Manager) ClaimLevel(playerID string, level int) (*LevelUpResult, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, ok := m.players[playerID]
	if !ok {
		return nil, false
	}

	// 找到等級定義
	var lvlDef *SeasonLevel
	for i := range SeasonLevels {
		if SeasonLevels[i].Level == level {
			lvlDef = &SeasonLevels[i]
			break
		}
	}
	if lvlDef == nil {
		return nil, false
	}

	// 確認積分足夠且未領取
	if data.SeasonPoints < lvlDef.PointsNeeded {
		return nil, false
	}
	if m.hasClaimed(data, level) {
		return nil, false
	}

	// 標記已領取
	data.ClaimedLevels = append(data.ClaimedLevels, level)
	if level > data.CurrentLevel {
		data.CurrentLevel = level
	}

	return &LevelUpResult{
		NewLevel:    level,
		CoinReward:  lvlDef.CoinReward,
		SpecialType: lvlDef.SpecialType,
		SpecialID:   lvlDef.SpecialID,
		SpecialName: lvlDef.SpecialName,
	}, true
}

// GetSnapshot 取得玩家賽季快照
func (m *Manager) GetSnapshot(playerID string) PlayerSeasonSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, ok := m.players[playerID]
	if !ok {
		return PlayerSeasonSnapshot{
			PlayerID:     playerID,
			SeasonPoints: 0,
			CurrentLevel: 0,
			NextLevel:    1,
			PointsToNext: SeasonLevels[0].PointsNeeded,
			Progress:     0.0,
			Levels:       buildLevelStatus(nil),
		}
	}

	// 計算下一個等級
	nextLevel := 0
	pointsToNext := 0
	progress := 1.0
	for _, lvl := range SeasonLevels {
		if !m.hasClaimed(data, lvl.Level) {
			nextLevel = lvl.Level
			pointsToNext = lvl.PointsNeeded - data.SeasonPoints
			if pointsToNext < 0 {
				pointsToNext = 0
			}
			// 計算進度（到下一個等級）
			prevPoints := 0
			if lvl.Level > 1 {
				prevPoints = SeasonLevels[lvl.Level-2].PointsNeeded
			}
			span := lvl.PointsNeeded - prevPoints
			earned := data.SeasonPoints - prevPoints
			if span > 0 {
				progress = float64(earned) / float64(span)
				if progress > 1.0 {
					progress = 1.0
				}
				if progress < 0 {
					progress = 0
				}
			}
			break
		}
	}

	return PlayerSeasonSnapshot{
		PlayerID:     playerID,
		SeasonPoints: data.SeasonPoints,
		CurrentLevel: data.CurrentLevel,
		NextLevel:    nextLevel,
		PointsToNext: pointsToNext,
		Progress:     progress,
		Levels:       buildLevelStatus(data),
	}
}

// PlayerSeasonSnapshot 玩家賽季快照（用於 WebSocket 廣播）
type PlayerSeasonSnapshot struct {
	PlayerID     string        `json:"player_id"`
	SeasonPoints int           `json:"season_points"`
	CurrentLevel int           `json:"current_level"`
	NextLevel    int           `json:"next_level"`    // 0 = 已滿級
	PointsToNext int           `json:"points_to_next"` // 距離下一等級所需積分
	Progress     float64       `json:"progress"`       // 0.0-1.0 當前等級進度
	Levels       []LevelStatus `json:"levels"`
}

// LevelStatus 等級狀態
type LevelStatus struct {
	Level       int    `json:"level"`
	PointsNeeded int   `json:"points_needed"`
	CoinReward  int    `json:"coin_reward"`
	SpecialType string `json:"special_type"`
	SpecialID   string `json:"special_id"`
	SpecialName string `json:"special_name"`
	Icon        string `json:"icon"`
	Claimed     bool   `json:"claimed"`
	Unlocked    bool   `json:"unlocked"` // 積分已達到但未領取
}

// buildLevelStatus 建立等級狀態列表
func buildLevelStatus(data *PlayerSeasonData) []LevelStatus {
	result := make([]LevelStatus, len(SeasonLevels))
	for i, lvl := range SeasonLevels {
		claimed := false
		unlocked := false
		if data != nil {
			claimed = hasClaimed(data, lvl.Level)
			unlocked = data.SeasonPoints >= lvl.PointsNeeded
		}
		result[i] = LevelStatus{
			Level:        lvl.Level,
			PointsNeeded: lvl.PointsNeeded,
			CoinReward:   lvl.CoinReward,
			SpecialType:  lvl.SpecialType,
			SpecialID:    lvl.SpecialID,
			SpecialName:  lvl.SpecialName,
			Icon:         lvl.Icon,
			Claimed:      claimed,
			Unlocked:     unlocked,
		}
	}
	return result
}

// hasClaimed 檢查是否已領取（非 thread-safe，呼叫前需持有鎖）
func (m *Manager) hasClaimed(data *PlayerSeasonData, level int) bool {
	return hasClaimed(data, level)
}

// hasClaimed 獨立函式（供 buildLevelStatus 使用）
func hasClaimed(data *PlayerSeasonData, level int) bool {
	for _, l := range data.ClaimedLevels {
		if l == level {
			return true
		}
	}
	return false
}
