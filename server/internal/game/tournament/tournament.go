// Package tournament 週賽 + 每日賽系統（DAY-093 升級）
// 週賽：每週重置排行榜，前三名獲得大獎
// 每日賽：每日 UTC+8 00:00 重置，前三名獲得每日獎勵
// 積分來源：擊破目標（依倍率）+ BOSS 擊殺 + Bonus 完成
package tournament

import (
	"sort"
	"sync"
	"time"
)

// PointSource 積分來源
type PointSource string

const (
	PointKill  PointSource = "kill"  // 擊破目標
	PointBoss  PointSource = "boss"  // 擊殺 BOSS
	PointBonus PointSource = "bonus" // 完成 Bonus
)

// PrizeConfig 獎勵設定
type PrizeConfig struct {
	Rank   int
	Coins  int
	Label  string
}

// DefaultPrizes 預設獎勵（前三名）
var DefaultPrizes = []PrizeConfig{
	{Rank: 1, Coins: 50000, Label: "🥇 週賽冠軍"},
	{Rank: 2, Coins: 25000, Label: "🥈 週賽亞軍"},
	{Rank: 3, Coins: 10000, Label: "🥉 週賽季軍"},
}

// DefaultDailyPrizes 每日賽預設獎勵（前三名）
var DefaultDailyPrizes = []PrizeConfig{
	{Rank: 1, Coins: 5000, Label: "🥇 日賽冠軍"},
	{Rank: 2, Coins: 2000, Label: "🥈 日賽亞軍"},
	{Rank: 3, Coins: 1000, Label: "🥉 日賽季軍"},
}

// Entry 週賽參賽者記錄
type Entry struct {
	PlayerID    string
	DisplayName string
	Points      int
	KillCount   int
	BossKills   int
	BonusCount  int
	LastUpdated time.Time
}

// WeeklyResult 週賽結算結果
type WeeklyResult struct {
	WeekStart time.Time
	WeekEnd   time.Time
	Rankings  []RankEntry
	SettledAt time.Time
}

// RankEntry 排名記錄
type RankEntry struct {
	Rank        int
	PlayerID    string
	DisplayName string
	Points      int
	Prize       int
	PrizeLabel  string
}

// Tournament 週賽管理器
type Tournament struct {
	mu        sync.RWMutex
	entries   map[string]*Entry // playerID → Entry
	weekStart time.Time
	weekEnd   time.Time
	history   []WeeklyResult // 最近 4 週歷史
}

// New 建立新的週賽管理器
func New() *Tournament {
	now := time.Now()
	start, end := currentWeekRange(now)
	return &Tournament{
		entries:   make(map[string]*Entry),
		weekStart: start,
		weekEnd:   end,
		history:   make([]WeeklyResult, 0, 4),
	}
}

// currentWeekRange 計算當前週的開始（週一 00:00 UTC+8）和結束（週日 23:59:59 UTC+8）
func currentWeekRange(t time.Time) (start, end time.Time) {
	loc := time.FixedZone("UTC+8", 8*3600)
	local := t.In(loc)

	// 找到本週一
	weekday := int(local.Weekday())
	if weekday == 0 {
		weekday = 7 // 週日 = 7
	}
	daysToMonday := weekday - 1

	monday := time.Date(local.Year(), local.Month(), local.Day()-daysToMonday, 0, 0, 0, 0, loc)
	sunday := monday.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return monday, sunday
}

// AddPoints 增加玩家積分
// multiplier 用於計算擊破積分（倍率越高積分越多）
func (t *Tournament) AddPoints(playerID, displayName string, source PointSource, multiplier float64) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 檢查是否需要重置（新的一週）
	t.checkAndReset()

	entry, ok := t.entries[playerID]
	if !ok {
		entry = &Entry{
			PlayerID:    playerID,
			DisplayName: displayName,
		}
		t.entries[playerID] = entry
	}

	// 更新顯示名稱（可能已更改）
	if displayName != "" {
		entry.DisplayName = displayName
	}

	// 計算積分
	var pts int
	switch source {
	case PointKill:
		// 擊破積分 = max(1, floor(multiplier))
		pts = int(multiplier)
		if pts < 1 {
			pts = 1
		}
		entry.KillCount++
	case PointBoss:
		pts = 50 // BOSS 擊殺固定 50 分
		entry.BossKills++
	case PointBonus:
		pts = 20 // Bonus 完成固定 20 分
		entry.BonusCount++
	}

	entry.Points += pts
	entry.LastUpdated = time.Now()

	return entry.Points
}

// GetRankings 取得當前排名（前 N 名）
func (t *Tournament) GetRankings(topN int) []RankEntry {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// 收集所有參賽者
	entries := make([]*Entry, 0, len(t.entries))
	for _, e := range t.entries {
		entries = append(entries, e)
	}

	// 依積分排序（降序），積分相同依 KillCount 排序
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Points != entries[j].Points {
			return entries[i].Points > entries[j].Points
		}
		return entries[i].KillCount > entries[j].KillCount
	})

	// 取前 N 名
	if topN > 0 && len(entries) > topN {
		entries = entries[:topN]
	}

	result := make([]RankEntry, len(entries))
	for i, e := range entries {
		rank := i + 1
		re := RankEntry{
			Rank:        rank,
			PlayerID:    e.PlayerID,
			DisplayName: e.DisplayName,
			Points:      e.Points,
		}
		// 加入獎勵資訊
		for _, prize := range DefaultPrizes {
			if prize.Rank == rank {
				re.Prize = prize.Coins
				re.PrizeLabel = prize.Label
				break
			}
		}
		result[i] = re
	}

	return result
}

// GetPlayerRank 取得特定玩家的排名和積分
func (t *Tournament) GetPlayerRank(playerID string) (rank int, points int) {
	rankings := t.GetRankings(0) // 取全部
	for _, r := range rankings {
		if r.PlayerID == playerID {
			return r.Rank, r.Points
		}
	}
	return 0, 0
}

// GetWeekInfo 取得當前週的時間資訊
func (t *Tournament) GetWeekInfo() (start, end time.Time, secondsLeft int64) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	left := time.Until(t.weekEnd)
	if left < 0 {
		left = 0
	}
	return t.weekStart, t.weekEnd, int64(left.Seconds())
}

// GetHistory 取得歷史週賽結果
func (t *Tournament) GetHistory() []WeeklyResult {
	t.mu.RLock()
	defer t.mu.RUnlock()
	result := make([]WeeklyResult, len(t.history))
	copy(result, t.history)
	return result
}

// checkAndReset 檢查是否需要重置（新的一週），必須在持有鎖的情況下呼叫
func (t *Tournament) checkAndReset() {
	now := time.Now()
	if now.After(t.weekEnd) {
		// 結算本週
		t.settle()
		// 重置為新的一週
		start, end := currentWeekRange(now)
		t.weekStart = start
		t.weekEnd = end
		t.entries = make(map[string]*Entry)
	}
}

// settle 結算本週（必須在持有鎖的情況下呼叫）
func (t *Tournament) settle() {
	if len(t.entries) == 0 {
		return
	}

	// 建立結算結果
	entries := make([]*Entry, 0, len(t.entries))
	for _, e := range t.entries {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Points > entries[j].Points
	})

	rankings := make([]RankEntry, 0, len(entries))
	for i, e := range entries {
		rank := i + 1
		re := RankEntry{
			Rank:        rank,
			PlayerID:    e.PlayerID,
			DisplayName: e.DisplayName,
			Points:      e.Points,
		}
		for _, prize := range DefaultPrizes {
			if prize.Rank == rank {
				re.Prize = prize.Coins
				re.PrizeLabel = prize.Label
				break
			}
		}
		rankings = append(rankings, re)
	}

	result := WeeklyResult{
		WeekStart: t.weekStart,
		WeekEnd:   t.weekEnd,
		Rankings:  rankings,
		SettledAt: time.Now(),
	}

	// 保留最近 4 週歷史
	t.history = append(t.history, result)
	if len(t.history) > 4 {
		t.history = t.history[len(t.history)-4:]
	}
}

// GetSnapshot 取得當前週賽快照（用於 HTTP 端點）
type Snapshot struct {
	WeekStart   int64       `json:"week_start"`    // Unix ms
	WeekEnd     int64       `json:"week_end"`      // Unix ms
	SecondsLeft int64       `json:"seconds_left"`  // 距離結束秒數
	Rankings    []RankEntry `json:"rankings"`      // 前 10 名
	TotalPlayers int        `json:"total_players"` // 本週參賽人數
	Prizes      []PrizeConfig `json:"prizes"`      // 獎勵設定
}

// GetSnapshot 取得快照
func (t *Tournament) GetSnapshot() Snapshot {
	t.mu.RLock()
	totalPlayers := len(t.entries)
	t.mu.RUnlock()

	start, end, left := t.GetWeekInfo()
	rankings := t.GetRankings(10)

	return Snapshot{
		WeekStart:    start.UnixMilli(),
		WeekEnd:      end.UnixMilli(),
		SecondsLeft:  left,
		Rankings:     rankings,
		TotalPlayers: totalPlayers,
		Prizes:       DefaultPrizes,
	}
}

// ============================================================
// DailyTournament — 每日賽管理器（DAY-093）
// 每日 UTC+8 00:00 重置，前三名獲得每日獎勵
// ============================================================

// DailyResult 每日賽結算結果
type DailyResult struct {
	Date      string      // "2026-05-20"
	Rankings  []RankEntry
	SettledAt time.Time
}

// DailyTournament 每日賽管理器
type DailyTournament struct {
	mu          sync.RWMutex
	entries     map[string]*Entry // playerID → Entry
	dayStart    time.Time
	dayEnd      time.Time
	history     []DailyResult // 最近 7 天歷史
}

// NewDaily 建立新的每日賽管理器
func NewDaily() *DailyTournament {
	now := time.Now()
	start, end := currentDayRange(now)
	return &DailyTournament{
		entries:  make(map[string]*Entry),
		dayStart: start,
		dayEnd:   end,
		history:  make([]DailyResult, 0, 7),
	}
}

// currentDayRange 計算當前日的開始（UTC+8 00:00）和結束（UTC+8 23:59:59）
func currentDayRange(t time.Time) (start, end time.Time) {
	loc := time.FixedZone("UTC+8", 8*3600)
	local := t.In(loc)
	start = time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
	end = start.Add(24*time.Hour - time.Second)
	return start, end
}

// AddPoints 增加玩家每日積分
func (d *DailyTournament) AddPoints(playerID, displayName string, source PointSource, multiplier float64) int {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.checkAndReset()

	entry, ok := d.entries[playerID]
	if !ok {
		entry = &Entry{
			PlayerID:    playerID,
			DisplayName: displayName,
		}
		d.entries[playerID] = entry
	}
	if displayName != "" {
		entry.DisplayName = displayName
	}

	var pts int
	switch source {
	case PointKill:
		pts = int(multiplier)
		if pts < 1 {
			pts = 1
		}
		entry.KillCount++
	case PointBoss:
		pts = 50
		entry.BossKills++
	case PointBonus:
		pts = 20
		entry.BonusCount++
	}

	entry.Points += pts
	entry.LastUpdated = time.Now()
	return entry.Points
}

// GetRankings 取得每日賽當前排名（前 N 名）
func (d *DailyTournament) GetRankings(topN int) []RankEntry {
	d.mu.RLock()
	defer d.mu.RUnlock()

	entries := make([]*Entry, 0, len(d.entries))
	for _, e := range d.entries {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Points != entries[j].Points {
			return entries[i].Points > entries[j].Points
		}
		return entries[i].KillCount > entries[j].KillCount
	})
	if topN > 0 && len(entries) > topN {
		entries = entries[:topN]
	}

	result := make([]RankEntry, len(entries))
	for i, e := range entries {
		rank := i + 1
		re := RankEntry{
			Rank:        rank,
			PlayerID:    e.PlayerID,
			DisplayName: e.DisplayName,
			Points:      e.Points,
		}
		for _, prize := range DefaultDailyPrizes {
			if prize.Rank == rank {
				re.Prize = prize.Coins
				re.PrizeLabel = prize.Label
				break
			}
		}
		result[i] = re
	}
	return result
}

// GetPlayerRank 取得特定玩家的每日排名和積分
func (d *DailyTournament) GetPlayerRank(playerID string) (rank int, points int) {
	rankings := d.GetRankings(0)
	for _, r := range rankings {
		if r.PlayerID == playerID {
			return r.Rank, r.Points
		}
	}
	return 0, 0
}

// GetDayInfo 取得當前日的時間資訊
func (d *DailyTournament) GetDayInfo() (start, end time.Time, secondsLeft int64) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	left := time.Until(d.dayEnd)
	if left < 0 {
		left = 0
	}
	return d.dayStart, d.dayEnd, int64(left.Seconds())
}

// GetHistory 取得歷史每日賽結果（最近 7 天）
func (d *DailyTournament) GetHistory() []DailyResult {
	d.mu.RLock()
	defer d.mu.RUnlock()
	result := make([]DailyResult, len(d.history))
	copy(result, d.history)
	return result
}

// checkAndReset 檢查是否需要重置（新的一天），必須在持有鎖的情況下呼叫
func (d *DailyTournament) checkAndReset() {
	now := time.Now()
	if now.After(d.dayEnd) {
		d.settleDaily()
		start, end := currentDayRange(now)
		d.dayStart = start
		d.dayEnd = end
		d.entries = make(map[string]*Entry)
	}
}

// settleDaily 結算當日（必須在持有鎖的情況下呼叫）
func (d *DailyTournament) settleDaily() {
	if len(d.entries) == 0 {
		return
	}
	entries := make([]*Entry, 0, len(d.entries))
	for _, e := range d.entries {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Points > entries[j].Points
	})

	rankings := make([]RankEntry, 0, len(entries))
	for i, e := range entries {
		rank := i + 1
		re := RankEntry{
			Rank:        rank,
			PlayerID:    e.PlayerID,
			DisplayName: e.DisplayName,
			Points:      e.Points,
		}
		for _, prize := range DefaultDailyPrizes {
			if prize.Rank == rank {
				re.Prize = prize.Coins
				re.PrizeLabel = prize.Label
				break
			}
		}
		rankings = append(rankings, re)
	}

	loc := time.FixedZone("UTC+8", 8*3600)
	dateStr := d.dayStart.In(loc).Format("2006-01-02")
	result := DailyResult{
		Date:      dateStr,
		Rankings:  rankings,
		SettledAt: time.Now(),
	}
	d.history = append(d.history, result)
	if len(d.history) > 7 {
		d.history = d.history[len(d.history)-7:]
	}
}

// DailySnapshot 每日賽快照
type DailySnapshot struct {
	DayStart     int64         `json:"day_start"`     // Unix ms
	DayEnd       int64         `json:"day_end"`       // Unix ms
	SecondsLeft  int64         `json:"seconds_left"`  // 距離結束秒數
	Rankings     []RankEntry   `json:"rankings"`      // 前 10 名
	TotalPlayers int           `json:"total_players"` // 今日參賽人數
	Prizes       []PrizeConfig `json:"prizes"`        // 獎勵設定
}

// GetDailySnapshot 取得每日賽快照
func (d *DailyTournament) GetDailySnapshot() DailySnapshot {
	d.mu.RLock()
	totalPlayers := len(d.entries)
	d.mu.RUnlock()

	start, end, left := d.GetDayInfo()
	rankings := d.GetRankings(10)

	return DailySnapshot{
		DayStart:     start.UnixMilli(),
		DayEnd:       end.UnixMilli(),
		SecondsLeft:  left,
		Rankings:     rankings,
		TotalPlayers: totalPlayers,
		Prizes:       DefaultDailyPrizes,
	}
}
