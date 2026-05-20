// Package friend 好友系統（DAY-073）
// 玩家可以加好友、查看好友列表、比較積分
// 好友關係是雙向的（A 加 B，B 也要接受）
package friend

import (
	"sync"
	"time"
)

// FriendStatus 好友狀態
type FriendStatus string

const (
	StatusPending  FriendStatus = "pending"  // 等待對方接受
	StatusAccepted FriendStatus = "accepted" // 已成為好友
	StatusBlocked  FriendStatus = "blocked"  // 已封鎖
)

// FriendRequest 好友請求記錄
type FriendRequest struct {
	FromID      string       `json:"from_id"`
	ToID        string       `json:"to_id"`
	Status      FriendStatus `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// FriendInfo 好友資訊（用於列表顯示）
type FriendInfo struct {
	PlayerID    string `json:"player_id"`
	DisplayName string `json:"display_name"`
	IsOnline    bool   `json:"is_online"`
	Coins       int    `json:"coins"`
	KillCount   int    `json:"kill_count"`
	TitleName   string `json:"title_name"`
	TitleIcon   string `json:"title_icon"`
	SeasonLevel int    `json:"season_level"`
	SeasonPoints int   `json:"season_points"`
}

// Manager 好友系統管理器
type Manager struct {
	mu       sync.RWMutex
	requests map[string]*FriendRequest // key: fromID+":"+toID
	friends  map[string][]string       // playerID → 好友 ID 列表
}

// New 建立新的好友管理器
func New() *Manager {
	return &Manager{
		requests: make(map[string]*FriendRequest),
		friends:  make(map[string][]string),
	}
}

// requestKey 生成請求 key（確保 A→B 和 B→A 是不同的 key）
func requestKey(fromID, toID string) string {
	return fromID + ":" + toID
}

// SendRequest 發送好友請求
// 回傳 true=成功，false=已是好友或已有待處理請求
func (m *Manager) SendRequest(fromID, toID string) bool {
	if fromID == toID {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	// 已是好友
	if m.areFriends(fromID, toID) {
		return false
	}

	// 已有待處理請求（任一方向）
	key1 := requestKey(fromID, toID)
	key2 := requestKey(toID, fromID)
	if req, ok := m.requests[key1]; ok && req.Status == StatusPending {
		return false
	}
	if req, ok := m.requests[key2]; ok && req.Status == StatusPending {
		// 對方已發請求，直接接受
		req.Status = StatusAccepted
		req.UpdatedAt = time.Now()
		m.addFriendPair(fromID, toID)
		return true
	}

	// 建立新請求
	m.requests[key1] = &FriendRequest{
		FromID:    fromID,
		ToID:      toID,
		Status:    StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return true
}

// AcceptRequest 接受好友請求
func (m *Manager) AcceptRequest(fromID, toID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := requestKey(fromID, toID)
	req, ok := m.requests[key]
	if !ok || req.Status != StatusPending {
		return false
	}

	req.Status = StatusAccepted
	req.UpdatedAt = time.Now()
	m.addFriendPair(fromID, toID)
	return true
}

// RejectRequest 拒絕好友請求
func (m *Manager) RejectRequest(fromID, toID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := requestKey(fromID, toID)
	req, ok := m.requests[key]
	if !ok || req.Status != StatusPending {
		return false
	}

	delete(m.requests, key)
	return true
}

// RemoveFriend 移除好友
func (m *Manager) RemoveFriend(playerID, friendID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.areFriends(playerID, friendID) {
		return false
	}

	// 移除雙向好友關係
	m.friends[playerID] = removeFromSlice(m.friends[playerID], friendID)
	m.friends[friendID] = removeFromSlice(m.friends[friendID], playerID)

	// 清除請求記錄
	delete(m.requests, requestKey(playerID, friendID))
	delete(m.requests, requestKey(friendID, playerID))
	return true
}

// GetFriendIDs 取得好友 ID 列表
func (m *Manager) GetFriendIDs(playerID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids, ok := m.friends[playerID]
	if !ok {
		return []string{}
	}
	result := make([]string, len(ids))
	copy(result, ids)
	return result
}

// GetPendingRequests 取得待處理的好友請求（別人發給我的）
func (m *Manager) GetPendingRequests(playerID string) []*FriendRequest {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*FriendRequest
	for _, req := range m.requests {
		if req.ToID == playerID && req.Status == StatusPending {
			result = append(result, req)
		}
	}
	return result
}

// IsFriend 檢查是否為好友
func (m *Manager) IsFriend(playerID, friendID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.areFriends(playerID, friendID)
}

// GetFriendCount 取得好友數量
func (m *Manager) GetFriendCount(playerID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.friends[playerID])
}

// areFriends 檢查是否為好友（非 thread-safe，呼叫前需持有鎖）
func (m *Manager) areFriends(playerID, friendID string) bool {
	for _, id := range m.friends[playerID] {
		if id == friendID {
			return true
		}
	}
	return false
}

// addFriendPair 建立雙向好友關係（非 thread-safe，呼叫前需持有鎖）
func (m *Manager) addFriendPair(playerA, playerB string) {
	if m.friends[playerA] == nil {
		m.friends[playerA] = []string{}
	}
	if m.friends[playerB] == nil {
		m.friends[playerB] = []string{}
	}
	if !m.areFriends(playerA, playerB) {
		m.friends[playerA] = append(m.friends[playerA], playerB)
		m.friends[playerB] = append(m.friends[playerB], playerA)
	}
}

// removeFromSlice 從 slice 中移除元素
func removeFromSlice(slice []string, item string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}
