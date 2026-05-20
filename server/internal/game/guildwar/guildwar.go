// Package guildwar 公會戰系統（DAY-076）
// 每週一次公會間積分競爭，週一 UTC+8 00:00 開始，週日 23:59 結算
// 積分來源：公會成員的擊殺/BOSS/Bonus 貢獻
// 前三名公會獲得獎勵（會長代領，分配給所有成員）
package guildwar

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// WarStatus 公會戰狀態
type WarStatus string

const (
	WarStatusActive   WarStatus = "active"   // 進行中
	WarStatusSettling WarStatus = "settling" // 結算中
	WarStatusIdle     WarStatus = "idle"     // 等待下一場
)

// WarRank 公會戰排名獎勵
type WarRank struct {
	Rank        int    `json:"rank"`
	GuildID     string `json:"guild_id"`
	GuildName   string `json:"guild_name"`
	GuildIcon   string `json:"guild_icon"`
	Score       int    `json:"score"`
	MemberCount int    `json:"member_count"`
	Reward      int    `json:"reward"` // 每人獎勵金幣
}

// WarResult 公會戰結算結果
type WarResult struct {
	WeekID    string     `json:"week_id"`    // 格式：2026-W21
	StartAt   time.Time  `json:"start_at"`
	EndAt     time.Time  `json:"end_at"`
	Rankings  []*WarRank `json:"rankings"`
	SettledAt time.Time  `json:"settled_at"`
}

// GuildWarScore 公會在本週的積分記錄
type GuildWarScore struct {
	GuildID     string `json:"guild_id"`
	GuildName   string `json:"guild_name"`
	GuildIcon   string `json:"guild_icon"`
	MemberCount int    `json:"member_count"`
	Score       int    `json:"score"`
	KillScore   int    `json:"kill_score"`   // 擊殺積分
	BossScore   int    `json:"boss_score"`   // BOSS 積分
	BonusScore  int    `json:"bonus_score"`  // Bonus 積分
}

// Manager 公會戰管理器
type Manager struct {
	mu          sync.RWMutex
	scores      map[string]*GuildWarScore // guildID → score
	status      WarStatus
	weekID      string
	startAt     time.Time
	endAt       time.Time
	lastResult  *WarResult
	history     []*WarResult // 最近 4 週歷史
}

// New 建立新的公會戰管理器
func New() *Manager {
	m := &Manager{
		scores:  make(map[string]*GuildWarScore),
		status:  WarStatusIdle,
		history: make([]*WarResult, 0, 4),
	}
	m.startNewWar()
	return m
}

// startNewWar 開始新一週的公會戰（非 thread-safe）
func (m *Manager) startNewWar() {
	now := time.Now()
	weekID := getWeekID(now)
	start, end := getWeekRange(now)

	m.weekID = weekID
	m.startAt = start
	m.endAt = end
	m.status = WarStatusActive
	m.scores = make(map[string]*GuildWarScore)
}

// EnsureGuildRegistered 確保公會已在本週公會戰中登記
func (m *Manager) EnsureGuildRegistered(guildID, guildName, guildIcon string, memberCount int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.scores[guildID]; !ok {
		m.scores[guildID] = &GuildWarScore{
			GuildID:     guildID,
			GuildName:   guildName,
			GuildIcon:   guildIcon,
			MemberCount: memberCount,
		}
	} else {
		// 更新成員數和名稱（可能有變化）
		m.scores[guildID].GuildName = guildName
		m.scores[guildID].GuildIcon = guildIcon
		m.scores[guildID].MemberCount = memberCount
	}
}

// AddKillScore 增加擊殺積分（每個目標 = 1 分，高倍率目標額外加分）
func (m *Manager) AddKillScore(guildID string, multiplier int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	score, ok := m.scores[guildID]
	if !ok {
		return
	}

	// 基礎 1 分，高倍率額外加分
	pts := 1
	if multiplier >= 50 {
		pts = 5
	} else if multiplier >= 20 {
		pts = 3
	} else if multiplier >= 10 {
		pts = 2
	}

	score.KillScore += pts
	score.Score += pts
}

// AddBossScore 增加 BOSS 積分（每次 BOSS 擊殺 = 50 分）
func (m *Manager) AddBossScore(guildID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	score, ok := m.scores[guildID]
	if !ok {
		return
	}

	score.BossScore += 50
	score.Score += 50
}

// AddBonusScore 增加 Bonus 積分（每次 Bonus 完成 = 20 分）
func (m *Manager) AddBonusScore(guildID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	score, ok := m.scores[guildID]
	if !ok {
		return
	}

	score.BonusScore += 20
	score.Score += 20
}

// GetRankings 取得當前排名（依積分降序）
func (m *Manager) GetRankings() []*GuildWarScore {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*GuildWarScore, 0, len(m.scores))
	for _, s := range m.scores {
		result = append(result, s)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Score != result[j].Score {
			return result[i].Score > result[j].Score
		}
		return result[i].GuildID < result[j].GuildID
	})

	return result
}

// GetGuildRank 取得特定公會的排名（1-based，0=未登記）
func (m *Manager) GetGuildRank(guildID string) int {
	rankings := m.GetRankings()
	for i, r := range rankings {
		if r.GuildID == guildID {
			return i + 1
		}
	}
	return 0
}

// GetStatus 取得公會戰狀態快照
func (m *Manager) GetStatus() (WarStatus, string, time.Time, time.Time) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.status, m.weekID, m.startAt, m.endAt
}

// GetLastResult 取得上週結算結果
func (m *Manager) GetLastResult() *WarResult {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastResult
}

// GetHistory 取得歷史結算記錄
func (m *Manager) GetHistory() []*WarResult {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*WarResult, len(m.history))
	copy(result, m.history)
	return result
}

// CheckAndSettle 檢查是否需要結算（每分鐘呼叫一次）
// 回傳結算結果（nil = 未結算）
func (m *Manager) CheckAndSettle() *WarResult {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	if m.status != WarStatusActive || now.Before(m.endAt) {
		// 檢查是否需要開始新一週
		if m.status == WarStatusIdle || (m.status == WarStatusSettling && now.After(m.endAt.Add(5*time.Minute))) {
			m.startNewWar()
		}
		return nil
	}

	// 結算
	m.status = WarStatusSettling
	result := m.settle()
	m.lastResult = result

	// 保留最近 4 週歷史
	m.history = append(m.history, result)
	if len(m.history) > 4 {
		m.history = m.history[1:]
	}

	// 5 分鐘後開始新一週（由下次 CheckAndSettle 觸發）
	return result
}

// settle 執行結算（非 thread-safe）
func (m *Manager) settle() *WarResult {
	rankings := make([]*GuildWarScore, 0, len(m.scores))
	for _, s := range m.scores {
		rankings = append(rankings, s)
	}

	sort.Slice(rankings, func(i, j int) bool {
		if rankings[i].Score != rankings[j].Score {
			return rankings[i].Score > rankings[j].Score
		}
		return rankings[i].GuildID < rankings[j].GuildID
	})

	// 計算獎勵
	rewardTable := []int{10000, 5000, 2000} // 前三名每人獎勵
	warRanks := make([]*WarRank, 0, len(rankings))
	for i, s := range rankings {
		reward := 0
		if i < len(rewardTable) {
			reward = rewardTable[i]
		}
		warRanks = append(warRanks, &WarRank{
			Rank:        i + 1,
			GuildID:     s.GuildID,
			GuildName:   s.GuildName,
			GuildIcon:   s.GuildIcon,
			Score:       s.Score,
			MemberCount: s.MemberCount,
			Reward:      reward,
		})
	}

	return &WarResult{
		WeekID:    m.weekID,
		StartAt:   m.startAt,
		EndAt:     m.endAt,
		Rankings:  warRanks,
		SettledAt: time.Now(),
	}
}

// GetGuildScore 取得特定公會的積分
func (m *Manager) GetGuildScore(guildID string) *GuildWarScore {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.scores[guildID]
}

// GetParticipatingGuildCount 取得參與公會數
func (m *Manager) GetParticipatingGuildCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.scores)
}

// getWeekID 取得週 ID（格式：2026-W21）
func getWeekID(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-W%02d", year, week)
}

// getWeekRange 取得本週的開始和結束時間（UTC+8，週一 00:00 ~ 週日 23:59:59）
func getWeekRange(t time.Time) (start, end time.Time) {
	loc := time.FixedZone("UTC+8", 8*60*60)
	now := t.In(loc)

	// 找到本週週一
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 週日 = 7
	}
	daysToMonday := weekday - 1

	monday := time.Date(now.Year(), now.Month(), now.Day()-daysToMonday, 0, 0, 0, 0, loc)
	sunday := monday.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return monday.UTC(), sunday.UTC()
}
