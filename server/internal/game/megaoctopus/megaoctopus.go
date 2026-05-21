// Package megaoctopus 實作巨型章魚轉盤系統（DAY-144）
// 業界依據：JILI Mega Fishing「Mega Octopus Wheel – Defeat that giant octopus and enter
// the bonus wheel round where you have a chance to win massive guaranteed prizes up to 950x.」
// 擊破 T108 後觸發個人轉盤，最高 950x 獎勵
package megaoctopus

import (
	"math/rand"
	"sync"
	"time"
)

// WheelSlot 轉盤格子定義
type WheelSlot struct {
	Multiplier int    // 倍率（乘以 betLevel）
	Weight     int    // 權重（越高越容易出現）
	Color      string // 顯示顏色
	Label      string // 顯示文字
}

// WheelSlots 轉盤格子（8格，對應業界 Mega Octopus Wheel）
var WheelSlots = []WheelSlot{
	{Multiplier: 50, Weight: 30, Color: "#C0C0C0", Label: "50x"},
	{Multiplier: 100, Weight: 25, Color: "#FFD700", Label: "100x"},
	{Multiplier: 150, Weight: 18, Color: "#FFD700", Label: "150x"},
	{Multiplier: 200, Weight: 12, Color: "#FF8C00", Label: "200x"},
	{Multiplier: 300, Weight: 8, Color: "#FF4500", Label: "300x"},
	{Multiplier: 500, Weight: 4, Color: "#FF0080", Label: "500x"},
	{Multiplier: 750, Weight: 2, Color: "#9400D3", Label: "750x"},
	{Multiplier: 950, Weight: 1, Color: "#FF0000", Label: "950x 👑"},
}

// SpinDuration 轉盤旋轉時間（秒）
const SpinDuration = 3

// AnnounceThreshold 全服公告門檻（倍率）
const AnnounceThreshold = 300

// Session 玩家轉盤 session
type Session struct {
	PlayerID    string
	ResultIndex int // 預先決定的結果格子索引
	StartedAt   time.Time
	Stopped     bool
}

// Manager 巨型章魚轉盤管理器
type Manager struct {
	mu       sync.RWMutex
	sessions map[string]*Session // playerID -> session
	rng      *rand.Rand
}

// NewManager 建立管理器
func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// StartSession 開始轉盤 session（擊破 T108 後呼叫）
// 預先決定結果（公平性保證），回傳 session
func (m *Manager) StartSession(playerID string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 計算總權重
	totalWeight := 0
	for _, slot := range WheelSlots {
		totalWeight += slot.Weight
	}

	// 加權隨機選擇結果
	roll := m.rng.Intn(totalWeight)
	resultIndex := 0
	cumulative := 0
	for i, slot := range WheelSlots {
		cumulative += slot.Weight
		if roll < cumulative {
			resultIndex = i
			break
		}
	}

	session := &Session{
		PlayerID:    playerID,
		ResultIndex: resultIndex,
		StartedAt:   time.Now(),
		Stopped:     false,
	}
	m.sessions[playerID] = session
	return session
}

// StopSession 玩家停止轉盤（點擊停止按鈕）
// 回傳結果格子索引和倍率，如果 session 不存在或已停止則回傳 -1
func (m *Manager) StopSession(playerID string) (resultIndex int, multiplier int, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[playerID]
	if !exists || session.Stopped {
		return -1, 0, false
	}

	session.Stopped = true
	resultIndex = session.ResultIndex
	multiplier = WheelSlots[resultIndex].Multiplier
	delete(m.sessions, playerID)
	return resultIndex, multiplier, true
}

// AutoStop 超時自動停止（SpinDuration 秒後）
// 回傳需要自動停止的 session 列表
func (m *Manager) AutoStop() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var expired []string
	now := time.Now()
	for playerID, session := range m.sessions {
		if !session.Stopped && now.Sub(session.StartedAt) >= time.Duration(SpinDuration+2)*time.Second {
			expired = append(expired, playerID)
		}
	}
	return expired
}

// HasActiveSession 是否有活躍 session
func (m *Manager) HasActiveSession(playerID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.sessions[playerID]
	return exists
}

// RemovePlayer 移除玩家（斷線時清理）
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, playerID)
}

// GetResultSlot 取得結果格子資訊
func GetResultSlot(index int) WheelSlot {
	if index < 0 || index >= len(WheelSlots) {
		return WheelSlots[0]
	}
	return WheelSlots[index]
}
