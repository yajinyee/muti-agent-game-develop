// format.go — 錦標賽多格式系統（DAY-111）
// 支援 4 種競賽格式，每日輪換，讓玩家每天有不同的競爭目標
package tournament

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// FormatType 錦標賽格式類型
type FormatType string

const (
	FormatScore      FormatType = "score"       // 積分賽（現有）— 累積擊破積分
	FormatMultiplier FormatType = "multiplier"  // 最高倍率賽 — 比誰單次倍率最高
	FormatReward     FormatType = "reward"      // 最高獎勵賽 — 比誰單次獎勵最高
	FormatBet        FormatType = "bet"         // 投注競賽 — 比誰總投注最多
)

// FormatDef 格式定義
type FormatDef struct {
	Type        FormatType `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Icon        string     `json:"icon"`
	Unit        string     `json:"unit"` // 分數單位（"分"/"x"/"金幣"）
}

// AllFormats 所有格式定義
var AllFormats = []FormatDef{
	{
		Type:        FormatScore,
		Name:        "積分賽",
		Description: "擊破目標累積積分，倍率越高積分越多",
		Icon:        "⭐",
		Unit:        "分",
	},
	{
		Type:        FormatMultiplier,
		Name:        "最高倍率賽",
		Description: "比誰能打出最高單次倍率，流星和寶箱怪是關鍵",
		Icon:        "⚡",
		Unit:        "x",
	},
	{
		Type:        FormatReward,
		Name:        "最高獎勵賽",
		Description: "比誰能獲得最高單次獎勵金幣，高投注高回報",
		Icon:        "💰",
		Unit:        "金幣",
	},
	{
		Type:        FormatBet,
		Name:        "投注競賽",
		Description: "比誰今日總投注最多，高投注玩家的主場",
		Icon:        "🎯",
		Unit:        "金幣",
	},
}

// GetFormatDef 取得格式定義
func GetFormatDef(ft FormatType) FormatDef {
	for _, f := range AllFormats {
		if f.Type == ft {
			return f
		}
	}
	return AllFormats[0]
}

// FormatEntry 多格式錦標賽參賽者記錄
type FormatEntry struct {
	PlayerID    string
	DisplayName string
	// 積分賽
	TotalPoints int
	KillCount   int
	BossKills   int
	BonusCount  int
	// 最高倍率賽
	BestMultiplier float64
	// 最高獎勵賽
	BestReward int
	// 投注競賽
	TotalBet int
	// 通用
	LastUpdated time.Time
}

// GetScore 依格式取得分數
func (e *FormatEntry) GetScore(ft FormatType) float64 {
	switch ft {
	case FormatScore:
		return float64(e.TotalPoints)
	case FormatMultiplier:
		return e.BestMultiplier
	case FormatReward:
		return float64(e.BestReward)
	case FormatBet:
		return float64(e.TotalBet)
	}
	return float64(e.TotalPoints)
}

// FormatRankEntry 多格式排名記錄
type FormatRankEntry struct {
	Rank        int        `json:"rank"`
	PlayerID    string     `json:"player_id"`
	DisplayName string     `json:"display_name"`
	Score       float64    `json:"score"`
	ScoreLabel  string     `json:"score_label"` // 格式化後的分數文字
	Prize       int        `json:"prize"`
	PrizeLabel  string     `json:"prize_label"`
	IsSelf      bool       `json:"is_self"`
}

// MultiFormatTournament 多格式每日錦標賽
type MultiFormatTournament struct {
	mu          sync.RWMutex
	entries     map[string]*FormatEntry
	dayStart    time.Time
	dayEnd      time.Time
	todayFormat FormatType // 今日格式（依日期輪換）
	history     []FormatDailyResult
}

// FormatDailyResult 多格式每日賽結算結果
type FormatDailyResult struct {
	Date      string          `json:"date"`
	Format    FormatType      `json:"format"`
	Rankings  []FormatRankEntry `json:"rankings"`
	SettledAt time.Time       `json:"settled_at"`
}

// NewMultiFormat 建立多格式每日錦標賽
func NewMultiFormat() *MultiFormatTournament {
	now := time.Now()
	start, end := currentDayRange(now)
	return &MultiFormatTournament{
		entries:     make(map[string]*FormatEntry),
		dayStart:    start,
		dayEnd:      end,
		todayFormat: getTodayFormat(now),
		history:     make([]FormatDailyResult, 0, 7),
	}
}

// getTodayFormat 依日期決定今日格式（4天一輪）
func getTodayFormat(t time.Time) FormatType {
	loc := time.FixedZone("UTC+8", 8*3600)
	dayOfYear := t.In(loc).YearDay()
	formats := []FormatType{FormatScore, FormatMultiplier, FormatReward, FormatBet}
	return formats[dayOfYear%4]
}

// GetTodayFormat 取得今日格式
func (m *MultiFormatTournament) GetTodayFormat() FormatType {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.todayFormat
}

// RecordKill 記錄擊破事件（更新所有相關格式的分數）
func (m *MultiFormatTournament) RecordKill(playerID, displayName string, multiplier float64, reward int, betCost int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.checkAndReset()

	entry := m.getOrCreate(playerID, displayName)

	// 積分賽：依倍率計算積分
	pts := int(multiplier)
	if pts < 1 {
		pts = 1
	}
	entry.TotalPoints += pts
	entry.KillCount++

	// 最高倍率賽：更新最高倍率
	if multiplier > entry.BestMultiplier {
		entry.BestMultiplier = multiplier
	}

	// 最高獎勵賽：更新最高單次獎勵
	if reward > entry.BestReward {
		entry.BestReward = reward
	}

	// 投注競賽：累積投注
	entry.TotalBet += betCost

	entry.LastUpdated = time.Now()
}

// RecordBoss 記錄 BOSS 擊殺
func (m *MultiFormatTournament) RecordBoss(playerID, displayName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.checkAndReset()

	entry := m.getOrCreate(playerID, displayName)
	entry.TotalPoints += 50
	entry.BossKills++
	entry.LastUpdated = time.Now()
}

// RecordBonus 記錄 Bonus 完成
func (m *MultiFormatTournament) RecordBonus(playerID, displayName string, reward int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.checkAndReset()

	entry := m.getOrCreate(playerID, displayName)
	entry.TotalPoints += 20
	entry.BonusCount++
	// Bonus 獎勵也計入最高獎勵賽
	if reward > entry.BestReward {
		entry.BestReward = reward
	}
	entry.LastUpdated = time.Now()
}

// RecordShot 記錄射擊（投注競賽用）
func (m *MultiFormatTournament) RecordShot(playerID, displayName string, betCost int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.checkAndReset()

	entry := m.getOrCreate(playerID, displayName)
	entry.TotalBet += betCost
	entry.LastUpdated = time.Now()
}

// GetRankings 依今日格式取得排名
func (m *MultiFormatTournament) GetRankings(topN int) []FormatRankEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.getRankingsByFormat(m.todayFormat, topN)
}

// GetRankingsByFormat 依指定格式取得排名
func (m *MultiFormatTournament) GetRankingsByFormat(ft FormatType, topN int) []FormatRankEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.getRankingsByFormat(ft, topN)
}

func (m *MultiFormatTournament) getRankingsByFormat(ft FormatType, topN int) []FormatRankEntry {
	entries := make([]*FormatEntry, 0, len(m.entries))
	for _, e := range m.entries {
		entries = append(entries, e)
	}

	// 依格式排序
	sort.Slice(entries, func(i, j int) bool {
		si := entries[i].GetScore(ft)
		sj := entries[j].GetScore(ft)
		if si != sj {
			return si > sj
		}
		return entries[i].KillCount > entries[j].KillCount
	})

	if topN > 0 && len(entries) > topN {
		entries = entries[:topN]
	}

	result := make([]FormatRankEntry, len(entries))
	for i, e := range entries {
		rank := i + 1
		score := e.GetScore(ft)
		re := FormatRankEntry{
			Rank:        rank,
			PlayerID:    e.PlayerID,
			DisplayName: e.DisplayName,
			Score:       score,
			ScoreLabel:  formatScore(ft, score),
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

// GetPlayerRank 取得玩家在今日格式的排名
func (m *MultiFormatTournament) GetPlayerRank(playerID string) (rank int, score float64) {
	rankings := m.GetRankings(0)
	for _, r := range rankings {
		if r.PlayerID == playerID {
			return r.Rank, r.Score
		}
	}
	return 0, 0
}

// GetDayInfo 取得今日時間資訊
func (m *MultiFormatTournament) GetDayInfo() (start, end time.Time, secondsLeft int64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	left := time.Until(m.dayEnd)
	if left < 0 {
		left = 0
	}
	return m.dayStart, m.dayEnd, int64(left.Seconds())
}

// GetSnapshot 取得多格式錦標賽快照
type MultiFormatSnapshot struct {
	DayStart     int64              `json:"day_start"`
	DayEnd       int64              `json:"day_end"`
	SecondsLeft  int64              `json:"seconds_left"`
	TodayFormat  FormatType         `json:"today_format"`
	FormatDef    FormatDef          `json:"format_def"`
	Rankings     []FormatRankEntry  `json:"rankings"`
	TotalPlayers int                `json:"total_players"`
	Prizes       []PrizeConfig      `json:"prizes"`
	NextFormat   FormatType         `json:"next_format"`    // 明日格式
	NextFormatDef FormatDef         `json:"next_format_def"` // 明日格式定義
}

func (m *MultiFormatTournament) GetSnapshot() MultiFormatSnapshot {
	m.mu.RLock()
	totalPlayers := len(m.entries)
	todayFmt := m.todayFormat
	m.mu.RUnlock()

	start, end, left := m.GetDayInfo()
	rankings := m.GetRankings(10)

	// 計算明日格式
	tomorrow := time.Now().Add(24 * time.Hour)
	nextFmt := getTodayFormat(tomorrow)

	return MultiFormatSnapshot{
		DayStart:      start.UnixMilli(),
		DayEnd:        end.UnixMilli(),
		SecondsLeft:   left,
		TodayFormat:   todayFmt,
		FormatDef:     GetFormatDef(todayFmt),
		Rankings:      rankings,
		TotalPlayers:  totalPlayers,
		Prizes:        DefaultDailyPrizes,
		NextFormat:    nextFmt,
		NextFormatDef: GetFormatDef(nextFmt),
	}
}

// GetHistory 取得歷史結果
func (m *MultiFormatTournament) GetHistory() []FormatDailyResult {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]FormatDailyResult, len(m.history))
	copy(result, m.history)
	return result
}

// getOrCreate 取得或建立玩家記錄（必須在持有鎖的情況下呼叫）
func (m *MultiFormatTournament) getOrCreate(playerID, displayName string) *FormatEntry {
	entry, ok := m.entries[playerID]
	if !ok {
		entry = &FormatEntry{
			PlayerID:    playerID,
			DisplayName: displayName,
		}
		m.entries[playerID] = entry
	}
	if displayName != "" {
		entry.DisplayName = displayName
	}
	return entry
}

// checkAndReset 檢查是否需要重置（必須在持有鎖的情況下呼叫）
func (m *MultiFormatTournament) checkAndReset() {
	now := time.Now()
	if now.After(m.dayEnd) {
		m.settleFormat()
		start, end := currentDayRange(now)
		m.dayStart = start
		m.dayEnd = end
		m.todayFormat = getTodayFormat(now)
		m.entries = make(map[string]*FormatEntry)
	}
}

// settleFormat 結算今日格式賽（必須在持有鎖的情況下呼叫）
func (m *MultiFormatTournament) settleFormat() {
	if len(m.entries) == 0 {
		return
	}
	rankings := m.getRankingsByFormat(m.todayFormat, 0)
	loc := time.FixedZone("UTC+8", 8*3600)
	dateStr := m.dayStart.In(loc).Format("2006-01-02")
	result := FormatDailyResult{
		Date:      dateStr,
		Format:    m.todayFormat,
		Rankings:  rankings,
		SettledAt: time.Now(),
	}
	m.history = append(m.history, result)
	if len(m.history) > 7 {
		m.history = m.history[len(m.history)-7:]
	}
}

// formatScore 格式化分數顯示
func formatScore(ft FormatType, score float64) string {
	switch ft {
	case FormatMultiplier:
		return fmt.Sprintf("%.0fx", score)
	case FormatReward, FormatBet:
		return fmt.Sprintf("%d", int(score))
	default:
		return fmt.Sprintf("%d", int(score))
	}
}
