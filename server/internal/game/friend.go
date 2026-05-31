// Package game — 好友系統
// DAY-349：好友排行榜、好友在線狀態
// social-ui-agent 負責維護
package game

import (
	"sort"
	"sync"
	"time"
)

// FriendEntry 好友記錄
type FriendEntry struct {
	PlayerID    string
	DisplayName string
	LastSeen    time.Time
	IsOnline    bool
	SeasonXP    int
	TotalKills  int
	BestMult    float64 // 歷史最高倍率
}

// FriendSystem 好友系統（純記憶體，以 DisplayName 為識別）
// 簡化版：同一房間的玩家自動互為「同場玩家」，可查看彼此排名
type FriendSystem struct {
	mu      sync.RWMutex
	// playerID -> FriendEntry（記錄所有曾在線的玩家）
	players map[string]*FriendEntry
}

func newFriendSystem() *FriendSystem {
	return &FriendSystem{
		players: make(map[string]*FriendEntry),
	}
}

// UpdatePlayer 更新玩家資訊（每次 player_update 時呼叫）
func (f *FriendSystem) UpdatePlayer(playerID, displayName string, seasonXP, totalKills int, bestMult float64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	entry, ok := f.players[playerID]
	if !ok {
		entry = &FriendEntry{PlayerID: playerID}
		f.players[playerID] = entry
	}
	entry.DisplayName = displayName
	entry.LastSeen = time.Now()
	entry.IsOnline = true
	entry.SeasonXP = seasonXP
	entry.TotalKills = totalKills
	if bestMult > entry.BestMult {
		entry.BestMult = bestMult
	}
}

// SetOffline 玩家離線
func (f *FriendSystem) SetOffline(playerID string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if entry, ok := f.players[playerID]; ok {
		entry.IsOnline = false
		entry.LastSeen = time.Now()
	}
}

// GetRoomLeaderboard 取得同場玩家排行榜（依賽季 XP 降序）
// 只回傳在線玩家
func (f *FriendSystem) GetRoomLeaderboard() []*FriendEntry {
	f.mu.RLock()
	defer f.mu.RUnlock()
	var entries []*FriendEntry
	for _, e := range f.players {
		if e.IsOnline {
			// 複製一份避免外部修改
			cp := *e
			entries = append(entries, &cp)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].SeasonXP != entries[j].SeasonXP {
			return entries[i].SeasonXP > entries[j].SeasonXP
		}
		return entries[i].TotalKills > entries[j].TotalKills
	})
	return entries
}

// GetPlayerRank 取得玩家在同場中的排名（1-based）
func (f *FriendSystem) GetPlayerRank(playerID string) int {
	entries := f.GetRoomLeaderboard()
	for i, e := range entries {
		if e.PlayerID == playerID {
			return i + 1
		}
	}
	return -1
}

// GetOnlineCount 取得在線玩家數
func (f *FriendSystem) GetOnlineCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	count := 0
	for _, e := range f.players {
		if e.IsOnline {
			count++
		}
	}
	return count
}
