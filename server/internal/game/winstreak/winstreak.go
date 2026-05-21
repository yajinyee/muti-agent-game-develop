// Package winstreak — 連勝獎勵系統（DAY-131）
// 業界依據：BGaming Fishing Club 2026 Best Win/Best Catch 里程碑
// 追蹤玩家在 session 中的連勝次數，達到里程碑時給予額外獎勵。
// 與連擊系統（2秒內連續）不同，這是更長時間的累積（30秒超時重置）。
package winstreak

import (
	"sync"
	"time"
)

// MilestoneLevel 里程碑等級
type MilestoneLevel int

const (
	MilestoneNone    MilestoneLevel = 0
	MilestoneBronze  MilestoneLevel = 10  // 10 連勝：銅牌
	MilestoneSilver  MilestoneLevel = 25  // 25 連勝：銀牌
	MilestoneGold    MilestoneLevel = 50  // 50 連勝：金牌
	MilestoneLegend  MilestoneLevel = 100 // 100 連勝：傳說
)

// MilestoneDef 里程碑定義
type MilestoneDef struct {
	Level      MilestoneLevel
	Name       string
	Icon       string
	Color      string
	BonusMult  float64 // 獎勵倍率（相對於 betCost）
	Broadcast  bool    // 是否全服廣播
}

// Milestones 里程碑定義表（依等級排序）
var Milestones = []MilestoneDef{
	{MilestoneBronze, "銅牌連勝", "🥉", "#CD7F32", 20.0, false},
	{MilestoneSilver, "銀牌連勝", "🥈", "#C0C0C0", 50.0, false},
	{MilestoneGold, "金牌連勝", "🥇", "#FFD700", 100.0, true},
	{MilestoneLegend, "傳說連勝", "🏆", "#FF69B4", 300.0, true},
}

// GetMilestoneDef 取得里程碑定義
func GetMilestoneDef(level MilestoneLevel) *MilestoneDef {
	for i := range Milestones {
		if Milestones[i].Level == level {
			return &Milestones[i]
		}
	}
	return nil
}

// PlayerStreak 玩家連勝狀態
type PlayerStreak struct {
	PlayerID       string
	Current        int            // 當前連勝次數
	MaxStreak      int            // 本 session 最高連勝
	LastKillAt     time.Time      // 上次擊破時間
	NextMilestone  MilestoneLevel // 下一個里程碑
	ClaimedLevels  map[MilestoneLevel]bool // 已達成的里程碑
}

// IsExpired 是否超時（30 秒未擊破則重置）
func (s *PlayerStreak) IsExpired() bool {
	if s.LastKillAt.IsZero() {
		return false
	}
	return time.Since(s.LastKillAt) > 30*time.Second
}

// RecordKill 記錄擊破，回傳（新連勝數, 達成的里程碑, 是否重置）
func (s *PlayerStreak) RecordKill() (int, *MilestoneDef, bool) {
	wasReset := false
	if s.IsExpired() {
		s.Current = 0
		wasReset = true
	}

	s.Current++
	s.LastKillAt = time.Now()
	if s.Current > s.MaxStreak {
		s.MaxStreak = s.Current
	}

	// 檢查是否達成里程碑
	var reached *MilestoneDef
	for i := range Milestones {
		m := &Milestones[i]
		if s.Current == int(m.Level) && !s.ClaimedLevels[m.Level] {
			s.ClaimedLevels[m.Level] = true
			reached = m
			break
		}
	}

	// 更新下一個里程碑
	s.NextMilestone = MilestoneNone
	for i := range Milestones {
		if !s.ClaimedLevels[Milestones[i].Level] {
			s.NextMilestone = Milestones[i].Level
			break
		}
	}

	return s.Current, reached, wasReset
}

// GetProgressToNext 取得到下一個里程碑的進度（0.0-1.0）
func (s *PlayerStreak) GetProgressToNext() float64 {
	if s.NextMilestone == MilestoneNone {
		return 1.0
	}
	// 找上一個里程碑
	prevLevel := 0
	for i := range Milestones {
		if Milestones[i].Level == s.NextMilestone {
			break
		}
		prevLevel = int(Milestones[i].Level)
	}
	span := int(s.NextMilestone) - prevLevel
	progress := s.Current - prevLevel
	if span <= 0 {
		return 1.0
	}
	p := float64(progress) / float64(span)
	if p > 1.0 {
		return 1.0
	}
	if p < 0 {
		return 0
	}
	return p
}

// Snapshot 快照
type Snapshot struct {
	Current         int
	MaxStreak       int
	NextMilestone   MilestoneLevel
	NextMilestoneName string
	ProgressToNext  float64
	SecondsToExpiry float64
}

// Manager 連勝管理器
type Manager struct {
	mu      sync.RWMutex
	players map[string]*PlayerStreak
}

// New 建立新管理器
func New() *Manager {
	return &Manager{
		players: make(map[string]*PlayerStreak),
	}
}

// EnsurePlayer 確保玩家記錄存在
func (m *Manager) EnsurePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.players[playerID]; !ok {
		m.players[playerID] = &PlayerStreak{
			PlayerID:      playerID,
			ClaimedLevels: make(map[MilestoneLevel]bool),
			NextMilestone: MilestoneBronze,
		}
	}
}

// RecordKill 記錄擊破
func (m *Manager) RecordKill(playerID string) (int, *MilestoneDef, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.players[playerID]
	if !ok {
		s = &PlayerStreak{
			PlayerID:      playerID,
			ClaimedLevels: make(map[MilestoneLevel]bool),
			NextMilestone: MilestoneBronze,
		}
		m.players[playerID] = s
	}
	return s.RecordKill()
}

// GetSnapshot 取得快照
func (m *Manager) GetSnapshot(playerID string) Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.players[playerID]
	if !ok {
		return Snapshot{NextMilestone: MilestoneBronze, NextMilestoneName: "銅牌連勝"}
	}

	nextName := ""
	if def := GetMilestoneDef(s.NextMilestone); def != nil {
		nextName = def.Name
	}

	secsToExpiry := 0.0
	if !s.LastKillAt.IsZero() {
		remaining := 30.0 - time.Since(s.LastKillAt).Seconds()
		if remaining > 0 {
			secsToExpiry = remaining
		}
	}

	return Snapshot{
		Current:           s.Current,
		MaxStreak:         s.MaxStreak,
		NextMilestone:     s.NextMilestone,
		NextMilestoneName: nextName,
		ProgressToNext:    s.GetProgressToNext(),
		SecondsToExpiry:   secsToExpiry,
	}
}

// CheckExpiry 檢查並重置過期的連勝（由 gameLoop 每秒呼叫）
// 回傳過期的玩家 ID 列表
func (m *Manager) CheckExpiry() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	expired := make([]string, 0)
	for id, s := range m.players {
		if s.Current > 0 && s.IsExpired() {
			s.Current = 0
			// 重置里程碑（讓玩家可以再次達成）
			s.ClaimedLevels = make(map[MilestoneLevel]bool)
			s.NextMilestone = MilestoneBronze
			expired = append(expired, id)
		}
	}
	return expired
}

// RemovePlayer 移除玩家
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.players, playerID)
}
