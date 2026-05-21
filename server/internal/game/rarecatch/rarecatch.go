// Package rarecatch — 稀有連擊累積倍率系統（DAY-126）
// 業界依據：fishingfortune.app（2026-05-21）確認「multiplier cascade system」
// 連續在 90 秒內擊破稀有目標，倍率從 2x 累積到最高 15x
// 業界研究顯示稀有目標專屬倍率讓玩家主動追求高價值目標，提升策略深度
package rarecatch

import (
	"sync"
	"time"
)

// RareTargetIDs 稀有目標 DefID 集合（T101-T105）
var RareTargetIDs = map[string]bool{
	"T101": true, // 擬態怪物
	"T102": true, // 寶箱怪
	"T103": true, // 流星
	"T104": true, // 金色雜草
	"T105": true, // 金幣魚
}

// IsRareTarget 判斷是否為稀有目標
func IsRareTarget(defID string) bool {
	return RareTargetIDs[defID]
}

// CascadeLevel 累積等級定義
type CascadeLevel struct {
	Count    int     // 需要的連擊數
	MultBoost float64 // 倍率加成
	Name     string  // 等級名稱
	Icon     string  // 圖示
	Color    string  // 顏色（hex）
}

// CascadeLevels 等級定義表（依連擊數升序）
var CascadeLevels = []CascadeLevel{
	{Count: 1, MultBoost: 2.0, Name: "稀有連擊", Icon: "💎", Color: "#00BFFF"},
	{Count: 2, MultBoost: 3.0, Name: "稀有連擊 ×2", Icon: "💎💎", Color: "#00FF7F"},
	{Count: 3, MultBoost: 5.0, Name: "稀有連擊 ×3", Icon: "💎💎💎", Color: "#FFD700"},
	{Count: 4, MultBoost: 8.0, Name: "稀有連擊 ×4", Icon: "🌟", Color: "#FF8C00"},
	{Count: 5, MultBoost: 15.0, Name: "稀有連擊 MAX", Icon: "🌈", Color: "#FF1493"},
}

// MaxCascadeCount 最大連擊數（超過後維持最高倍率）
const MaxCascadeCount = 5

// CascadeTimeout 連擊超時（90 秒未擊破稀有目標則重置）
const CascadeTimeout = 90 * time.Second

// Session 玩家的稀有連擊 session
type Session struct {
	Count     int       // 當前連擊數
	LastHitAt time.Time // 最後一次擊破稀有目標的時間
}

// IsExpired 是否已超時
func (s *Session) IsExpired() bool {
	return time.Since(s.LastHitAt) > CascadeTimeout
}

// GetMultBoost 取得當前倍率加成
func (s *Session) GetMultBoost() float64 {
	if s.Count <= 0 || s.IsExpired() {
		return 1.0
	}
	idx := s.Count - 1
	if idx >= len(CascadeLevels) {
		idx = len(CascadeLevels) - 1
	}
	return CascadeLevels[idx].MultBoost
}

// GetLevel 取得當前等級定義
func (s *Session) GetLevel() *CascadeLevel {
	if s.Count <= 0 || s.IsExpired() {
		return nil
	}
	idx := s.Count - 1
	if idx >= len(CascadeLevels) {
		idx = len(CascadeLevels) - 1
	}
	l := CascadeLevels[idx]
	return &l
}

// Snapshot 快照
type Snapshot struct {
	Count       int     `json:"count"`        // 當前連擊數
	MultBoost   float64 `json:"mult_boost"`   // 當前倍率加成
	LevelName   string  `json:"level_name"`   // 等級名稱
	Icon        string  `json:"icon"`         // 圖示
	Color       string  `json:"color"`        // 顏色
	SecondsLeft int     `json:"seconds_left"` // 距離超時的剩餘秒數
	IsActive    bool    `json:"is_active"`    // 是否有效
}

// Manager 稀有連擊管理器（每個玩家一個 session）
type Manager struct {
	mu       sync.RWMutex
	sessions map[string]*Session // playerID → session
}

// New 建立管理器
func New() *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
	}
}

// RecordKill 記錄稀有目標擊破
// 回傳：(新連擊數, 倍率加成, 是否升級, 是否達到廣播門檻)
func (m *Manager) RecordKill(playerID string) (count int, multBoost float64, isLevelUp bool, shouldBroadcast bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessions[playerID]
	if !ok || s.IsExpired() {
		// 新 session 或已超時，從 1 開始
		s = &Session{Count: 1, LastHitAt: time.Now()}
		m.sessions[playerID] = s
		isLevelUp = true
	} else {
		// 累積連擊
		prevCount := s.Count
		s.Count++
		if s.Count > MaxCascadeCount {
			s.Count = MaxCascadeCount
		}
		s.LastHitAt = time.Now()
		isLevelUp = s.Count > prevCount
	}

	count = s.Count
	multBoost = s.GetMultBoost()
	// 達到 ×5.0（第3次）以上時廣播
	shouldBroadcast = isLevelUp && multBoost >= 5.0

	return
}

// GetMultBoost 取得玩家當前倍率加成（1.0 = 無加成）
func (m *Manager) GetMultBoost(playerID string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.sessions[playerID]
	if !ok {
		return 1.0
	}
	return s.GetMultBoost()
}

// GetSnapshot 取得玩家快照
func (m *Manager) GetSnapshot(playerID string) Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.sessions[playerID]
	if !ok || s.IsExpired() {
		return Snapshot{IsActive: false}
	}

	level := s.GetLevel()
	if level == nil {
		return Snapshot{IsActive: false}
	}

	remaining := CascadeTimeout - time.Since(s.LastHitAt)
	secondsLeft := int(remaining.Seconds())
	if secondsLeft < 0 {
		secondsLeft = 0
	}

	return Snapshot{
		Count:       s.Count,
		MultBoost:   level.MultBoost,
		LevelName:   level.Name,
		Icon:        level.Icon,
		Color:       level.Color,
		SecondsLeft: secondsLeft,
		IsActive:    true,
	}
}

// CheckExpiry 檢查並清理過期 session（由 game loop 定期呼叫）
// 回傳過期的 playerID 列表
func (m *Manager) CheckExpiry() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var expired []string
	for playerID, s := range m.sessions {
		if s.IsExpired() {
			expired = append(expired, playerID)
			delete(m.sessions, playerID)
		}
	}
	return expired
}

// RemovePlayer 移除玩家 session
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, playerID)
}
