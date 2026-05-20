// Package friend 好友系統（DAY-073）
// 玩家可以加好友、查看好友列表、比較積分
// 好友關係是雙向的（A 加 B，B 也要接受）
// DAY-101：新增禮物贈送系統 + 好友關係持久化
package friend

import (
	"sync"
	"time"
)

// GiftRecord 禮物贈送記錄
type GiftRecord struct {
	SentCount int       `json:"sent_count"`  // 今日已送出次數
	LastDate  string    `json:"last_date"`   // 最後送出日期（UTC+8，格式 2006-01-02）
}

// GiftResult 禮物贈送結果
type GiftResult struct {
	Success     bool   `json:"success"`
	Amount      int    `json:"amount"`
	ErrorCode   string `json:"error_code,omitempty"`
	ErrorMsg    string `json:"error_msg,omitempty"`
}

// FriendState 好友關係持久化狀態（DAY-101）
type FriendState struct {
	PlayerID  string   `json:"player_id"`
	FriendIDs []string `json:"friend_ids"`
}

const (
	MaxDailyGifts  = 3    // 每日最多送出禮物次數
	GiftAmount     = 500  // 每次禮物金幣數量
	MaxFriends     = 50   // 最多好友數量
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
	gifts    map[string]*GiftRecord    // playerID → 今日禮物記錄（DAY-101）
}

// New 建立新的好友管理器
func New() *Manager {
	return &Manager{
		requests: make(map[string]*FriendRequest),
		friends:  make(map[string][]string),
		gifts:    make(map[string]*GiftRecord),
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

// ---- 禮物贈送系統（DAY-101）----

// todayDate 取得今日日期字串（UTC+8）
func todayDate() string {
	loc := time.FixedZone("UTC+8", 8*60*60)
	return time.Now().In(loc).Format("2006-01-02")
}

// SendGift 向好友贈送禮物
// 回傳 GiftResult（含成功/失敗原因）
func (m *Manager) SendGift(fromID, toID string) GiftResult {
	if fromID == toID {
		return GiftResult{ErrorCode: "self_gift", ErrorMsg: "不能送禮物給自己"}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 必須是好友
	if !m.areFriends(fromID, toID) {
		return GiftResult{ErrorCode: "not_friend", ErrorMsg: "只能送禮物給好友"}
	}

	// 檢查今日送出次數
	today := todayDate()
	rec, ok := m.gifts[fromID]
	if !ok {
		rec = &GiftRecord{}
		m.gifts[fromID] = rec
	}
	// 跨日重置
	if rec.LastDate != today {
		rec.SentCount = 0
		rec.LastDate = today
	}
	if rec.SentCount >= MaxDailyGifts {
		return GiftResult{ErrorCode: "daily_limit", ErrorMsg: "今日禮物已送完（每日上限 3 次）"}
	}

	rec.SentCount++
	rec.LastDate = today
	return GiftResult{Success: true, Amount: GiftAmount}
}

// GetGiftStatus 取得今日禮物狀態（已送次數 / 剩餘次數）
func (m *Manager) GetGiftStatus(playerID string) (sentToday int, remaining int) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	today := todayDate()
	rec, ok := m.gifts[playerID]
	if !ok || rec.LastDate != today {
		return 0, MaxDailyGifts
	}
	sent := rec.SentCount
	if sent > MaxDailyGifts {
		sent = MaxDailyGifts
	}
	return sent, MaxDailyGifts - sent
}

// ---- 持久化支援（DAY-101）----

// GetFriendState 取得好友關係持久化狀態
func (m *Manager) GetFriendState(playerID string) *FriendState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := m.friends[playerID]
	result := make([]string, len(ids))
	copy(result, ids)
	return &FriendState{
		PlayerID:  playerID,
		FriendIDs: result,
	}
}

// LoadFriendState 從持久化狀態恢復好友關係
// 只恢復單向（playerID 的好友列表），不重複建立雙向關係
func (m *Manager) LoadFriendState(state *FriendState) {
	if state == nil || len(state.FriendIDs) == 0 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.friends[state.PlayerID] == nil {
		m.friends[state.PlayerID] = []string{}
	}
	for _, fid := range state.FriendIDs {
		if !m.areFriends(state.PlayerID, fid) {
			m.friends[state.PlayerID] = append(m.friends[state.PlayerID], fid)
		}
	}
}
