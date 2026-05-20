// Package dm 玩家私訊系統（DAY-103）
// 好友間可以互相發送私訊，離線訊息暫存，上線後自動發送
// 業界依據：optikpi.com 2026 確認 in-app messaging 是留存核心機制
package dm

import (
	"sync"
	"time"
)

const (
	MaxMessageLength  = 200  // 單則訊息最大字元數
	MaxStoredMessages = 50   // 每個玩家最多暫存離線訊息數
	MaxDailyMessages  = 100  // 每日最多發送訊息數（防刷）
)

// Message 私訊記錄
type Message struct {
	ID          string    `json:"id"`
	FromID      string    `json:"from_id"`
	FromName    string    `json:"from_name"`
	ToID        string    `json:"to_id"`
	Content     string    `json:"content"`
	SentAt      time.Time `json:"sent_at"`
	IsRead      bool      `json:"is_read"`
}

// SendResult 發送結果
type SendResult struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

// DailyCount 每日發送計數
type DailyCount struct {
	Count    int    `json:"count"`
	LastDate string `json:"last_date"` // UTC+8 格式 2006-01-02
}

// Manager 私訊管理器
type Manager struct {
	mu sync.RWMutex
	// 離線訊息暫存：toID → []Message
	pending map[string][]*Message
	// 每日發送計數：fromID → DailyCount
	dailyCounts map[string]*DailyCount
	// 訊息 ID 計數器
	msgCounter int64
}

// New 建立新的私訊管理器
func New() *Manager {
	return &Manager{
		pending:     make(map[string][]*Message),
		dailyCounts: make(map[string]*DailyCount),
	}
}

// todayDate 取得今日日期字串（UTC+8）
func todayDate() string {
	loc := time.FixedZone("UTC+8", 8*60*60)
	return time.Now().In(loc).Format("2006-01-02")
}

// Send 發送私訊
// 如果接收者在線（由 deliverFn 判斷），直接發送；否則暫存
// deliverFn 回傳 true 表示已成功發送給在線玩家
func (m *Manager) Send(fromID, fromName, toID, content string, deliverFn func(*Message) bool) SendResult {
	if fromID == toID {
		return SendResult{ErrorCode: "self_message", ErrorMsg: "不能傳訊息給自己"}
	}
	if len(content) == 0 {
		return SendResult{ErrorCode: "empty_content", ErrorMsg: "訊息不能為空"}
	}
	if len([]rune(content)) > MaxMessageLength {
		return SendResult{ErrorCode: "too_long", ErrorMsg: "訊息太長（最多 200 字）"}
	}

	// 每日發送限制
	m.mu.Lock()
	today := todayDate()
	dc, ok := m.dailyCounts[fromID]
	if !ok {
		dc = &DailyCount{}
		m.dailyCounts[fromID] = dc
	}
	if dc.LastDate != today {
		dc.Count = 0
		dc.LastDate = today
	}
	if dc.Count >= MaxDailyMessages {
		m.mu.Unlock()
		return SendResult{ErrorCode: "daily_limit", ErrorMsg: "今日訊息已達上限（100 則）"}
	}
	dc.Count++

	m.msgCounter++
	counter := m.msgCounter
	msgID := fromID[:min(8, len(fromID))] + "-" + time.Now().Format("150405.000") + "-" + itoa(counter)
	msg := &Message{
		ID:       msgID,
		FromID:   fromID,
		FromName: fromName,
		ToID:     toID,
		Content:  content,
		SentAt:   time.Now(),
		IsRead:   false,
	}
	m.mu.Unlock()

	// 嘗試即時發送
	if deliverFn != nil && deliverFn(msg) {
		return SendResult{Success: true, MessageID: msgID}
	}

	// 接收者離線，暫存
	m.mu.Lock()
	pending := m.pending[toID]
	if len(pending) >= MaxStoredMessages {
		// 超過上限，移除最舊的
		pending = pending[1:]
	}
	m.pending[toID] = append(pending, msg)
	m.mu.Unlock()

	return SendResult{Success: true, MessageID: msgID}
}

// GetPending 取得並清除玩家的離線訊息
func (m *Manager) GetPending(playerID string) []*Message {
	m.mu.Lock()
	defer m.mu.Unlock()

	msgs := m.pending[playerID]
	if len(msgs) == 0 {
		return nil
	}
	delete(m.pending, playerID)
	return msgs
}

// GetDailyCount 取得今日發送計數
func (m *Manager) GetDailyCount(playerID string) (sent int, remaining int) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	today := todayDate()
	dc, ok := m.dailyCounts[playerID]
	if !ok || dc.LastDate != today {
		return 0, MaxDailyMessages
	}
	s := dc.Count
	if s > MaxDailyMessages {
		s = MaxDailyMessages
	}
	return s, MaxDailyMessages - s
}

// PendingCount 取得玩家待接收的離線訊息數量
func (m *Manager) PendingCount(playerID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.pending[playerID])
}

// min 輔助函數
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// itoa 整數轉字串
func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
