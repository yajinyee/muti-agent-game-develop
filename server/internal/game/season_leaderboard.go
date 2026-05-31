// Package game — 賽季排行榜（DAY-348）
// 設計：依賽季 XP 排名，每月重置，顯示前 20 名
// 整合：SeasonPassManager 的 XP 資料
package game

import (
	"sort"
	"sync"
	"time"
)

// LeaderboardEntry 排行榜條目
type LeaderboardEntry struct {
	Rank        int    `json:"rank"`
	PlayerID    string `json:"player_id"`
	DisplayName string `json:"display_name"`
	SeasonXP    int    `json:"season_xp"`
	Level       int    `json:"level"`
	LevelName   string `json:"level_name"`
	Badge       string `json:"badge"`
}

// SeasonLeaderboard 賽季排行榜
type SeasonLeaderboard struct {
	mu          sync.RWMutex
	entries     map[string]*leaderboardPlayerData // playerID -> data
	seasonID    string
	lastUpdated time.Time
}

type leaderboardPlayerData struct {
	PlayerID    string
	DisplayName string
	SeasonXP    int
	Level       int
	LevelName   string
	Badge       string
}

// NewSeasonLeaderboard 建立賽季排行榜
func NewSeasonLeaderboard(seasonID string) *SeasonLeaderboard {
	return &SeasonLeaderboard{
		entries:     make(map[string]*leaderboardPlayerData),
		seasonID:    seasonID,
		lastUpdated: time.Now(),
	}
}

// UpdatePlayer 更新玩家排行榜資料
func (lb *SeasonLeaderboard) UpdatePlayer(playerID, displayName string, xp, level int, levelName, badge string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	
	lb.entries[playerID] = &leaderboardPlayerData{
		PlayerID:    playerID,
		DisplayName: displayName,
		SeasonXP:    xp,
		Level:       level,
		LevelName:   levelName,
		Badge:       badge,
	}
	lb.lastUpdated = time.Now()
}

// GetTop 取得前 N 名排行榜
func (lb *SeasonLeaderboard) GetTop(n int) []LeaderboardEntry {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	
	// 轉換為切片並排序
	players := make([]*leaderboardPlayerData, 0, len(lb.entries))
	for _, p := range lb.entries {
		players = append(players, p)
	}
	
	// 依 XP 降序排序，XP 相同時依 PlayerID 排序（確保穩定）
	sort.Slice(players, func(i, j int) bool {
		if players[i].SeasonXP != players[j].SeasonXP {
			return players[i].SeasonXP > players[j].SeasonXP
		}
		return players[i].PlayerID < players[j].PlayerID
	})
	
	// 取前 N 名
	if n > len(players) {
		n = len(players)
	}
	
	result := make([]LeaderboardEntry, n)
	for i := 0; i < n; i++ {
		p := players[i]
		displayName := p.DisplayName
		if displayName == "" {
			displayName = "玩家" + p.PlayerID[:4]
		}
		result[i] = LeaderboardEntry{
			Rank:        i + 1,
			PlayerID:    p.PlayerID,
			DisplayName: displayName,
			SeasonXP:    p.SeasonXP,
			Level:       p.Level,
			LevelName:   p.LevelName,
			Badge:       p.Badge,
		}
	}
	return result
}

// GetPlayerRank 取得特定玩家的排名
func (lb *SeasonLeaderboard) GetPlayerRank(playerID string) (rank int, entry *LeaderboardEntry) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	
	// 轉換為切片並排序
	players := make([]*leaderboardPlayerData, 0, len(lb.entries))
	for _, p := range lb.entries {
		players = append(players, p)
	}
	sort.Slice(players, func(i, j int) bool {
		if players[i].SeasonXP != players[j].SeasonXP {
			return players[i].SeasonXP > players[j].SeasonXP
		}
		return players[i].PlayerID < players[j].PlayerID
	})
	
	for i, p := range players {
		if p.PlayerID == playerID {
			displayName := p.DisplayName
			if displayName == "" {
				displayName = "玩家" + p.PlayerID[:4]
			}
			e := &LeaderboardEntry{
				Rank:        i + 1,
				PlayerID:    p.PlayerID,
				DisplayName: displayName,
				SeasonXP:    p.SeasonXP,
				Level:       p.Level,
				LevelName:   p.LevelName,
				Badge:       p.Badge,
			}
			return i + 1, e
		}
	}
	return -1, nil
}

// GetSnapshot 取得排行榜快照（用於發送給 Client）
func (lb *SeasonLeaderboard) GetSnapshot(playerID string) map[string]interface{} {
	top20 := lb.GetTop(20)
	rank, playerEntry := lb.GetPlayerRank(playerID)
	
	result := map[string]interface{}{
		"season_id":    lb.seasonID,
		"top20":        top20,
		"last_updated": lb.lastUpdated.UnixMilli(),
	}
	
	if rank > 0 && playerEntry != nil {
		result["my_rank"]  = rank
		result["my_entry"] = playerEntry
	} else {
		result["my_rank"] = -1
	}
	
	return result
}

// Reset 重置排行榜（新賽季開始時呼叫）
func (lb *SeasonLeaderboard) Reset(newSeasonID string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	
	lb.entries = make(map[string]*leaderboardPlayerData)
	lb.seasonID = newSeasonID
	lb.lastUpdated = time.Now()
}

// GetPlayerCount 取得排行榜玩家數
func (lb *SeasonLeaderboard) GetPlayerCount() int {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	return len(lb.entries)
}
