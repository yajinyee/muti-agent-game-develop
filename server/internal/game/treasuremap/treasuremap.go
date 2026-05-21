// Package treasuremap 實作寶藏地圖系統（DAY-122）
// 業界依據：bsu.edu（2026）確認「Hidden Treasure Unlocks」是 2026 年捕魚機最新趨勢
// 玩家擊破特定目標物收集地圖格子，集滿一行/列/對角線觸發寶藏獎勵（類似賓果）
// 集滿整張地圖觸發傳說寶藏大獎，每日 UTC 重置
package treasuremap

import (
	"sync"
	"time"
)

// GridSize 地圖格子數（3×3）
const GridSize = 3

// 獎勵類型
const (
	RewardLine   = "line"   // 一行/列/對角線（3格）
	RewardFull   = "full"   // 整張地圖（9格）
)

// CellDef 地圖格子定義（對應目標物）
type CellDef struct {
	Row    int    // 0-2
	Col    int    // 0-2
	DefID  string // 對應目標物 ID
	Name   string // 顯示名稱
	Icon   string // 顯示圖示
}

// DefaultGrid 預設 3×3 地圖格子定義
// 每格對應一種目標物，玩家擊破對應目標物即可填滿該格
var DefaultGrid = [GridSize][GridSize]CellDef{
	{
		{Row: 0, Col: 0, DefID: "T001", Name: "像素雜草", Icon: "🌿"},
		{Row: 0, Col: 1, DefID: "T003", Name: "紅色小蟲", Icon: "🐛"},
		{Row: 0, Col: 2, DefID: "T005", Name: "布丁怪", Icon: "🍮"},
	},
	{
		{Row: 1, Col: 0, DefID: "T002", Name: "綠色小蟲", Icon: "🐝"},
		{Row: 1, Col: 1, DefID: "T101", Name: "擬態怪物", Icon: "👾"},
		{Row: 1, Col: 2, DefID: "T104", Name: "金色雜草", Icon: "✨"},
	},
	{
		{Row: 2, Col: 0, DefID: "T006", Name: "巨大蘑菇", Icon: "🍄"},
		{Row: 2, Col: 1, DefID: "T102", Name: "寶箱怪", Icon: "📦"},
		{Row: 2, Col: 2, DefID: "T105", Name: "金幣魚", Icon: "🪙"},
	},
}

// PlayerMap 玩家的寶藏地圖狀態
type PlayerMap struct {
	PlayerID    string
	Cells       [GridSize][GridSize]bool // 已填滿的格子
	FilledCount int                      // 已填滿格子數
	Lines       []LineResult             // 已完成的行/列/對角線
	FullDone    bool                     // 是否已完成整張地圖
	Date        string                   // UTC 日期（YYYY-MM-DD），用於每日重置
}

// LineResult 完成的行/列/對角線
type LineResult struct {
	Type  string // "row0"/"row1"/"row2"/"col0"/"col1"/"col2"/"diag0"/"diag1"
	Cells [3][2]int // 三個格子的 [row][col]
}

// Manager 寶藏地圖管理器
type Manager struct {
	mu      sync.Mutex
	players map[string]*PlayerMap
}

// New 建立新的寶藏地圖管理器
func New() *Manager {
	return &Manager{
		players: make(map[string]*PlayerMap),
	}
}

// getOrCreate 取得或建立玩家地圖（自動每日重置）
func (m *Manager) getOrCreate(playerID string) *PlayerMap {
	today := todayUTC()
	pm, ok := m.players[playerID]
	if !ok || pm.Date != today {
		// 新建或重置
		pm = &PlayerMap{
			PlayerID: playerID,
			Date:     today,
		}
		m.players[playerID] = pm
	}
	return pm
}

// RecordKill 記錄擊破目標物，回傳（是否有新格子填滿, 新完成的行/列/對角線, 是否完成整張地圖）
func (m *Manager) RecordKill(playerID, defID string) (bool, []LineResult, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pm := m.getOrCreate(playerID)

	// 找到對應格子
	row, col := findCell(defID)
	if row < 0 {
		return false, nil, false // 此目標物不在地圖上
	}

	// 已填滿，不重複計算
	if pm.Cells[row][col] {
		return false, nil, false
	}

	// 填滿格子
	pm.Cells[row][col] = true
	pm.FilledCount++

	// 檢查新完成的行/列/對角線
	newLines := checkNewLines(pm)
	for _, line := range newLines {
		pm.Lines = append(pm.Lines, line)
	}

	// 檢查是否完成整張地圖
	fullDone := false
	if pm.FilledCount == GridSize*GridSize && !pm.FullDone {
		pm.FullDone = true
		fullDone = true
	}

	return true, newLines, fullDone
}

// GetSnapshot 取得玩家地圖快照（唯讀）
func (m *Manager) GetSnapshot(playerID string) *PlayerMap {
	m.mu.Lock()
	defer m.mu.Unlock()
	pm := m.getOrCreate(playerID)
	// 回傳副本
	snap := *pm
	return &snap
}

// RemovePlayer 清理玩家資料（離線時呼叫）
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.players, playerID)
}

// findCell 找到目標物對應的格子位置，回傳 (row, col)，找不到回傳 (-1, -1)
func findCell(defID string) (int, int) {
	for r := 0; r < GridSize; r++ {
		for c := 0; c < GridSize; c++ {
			if DefaultGrid[r][c].DefID == defID {
				return r, c
			}
		}
	}
	return -1, -1
}

// checkNewLines 檢查是否有新完成的行/列/對角線（不包含已記錄的）
func checkNewLines(pm *PlayerMap) []LineResult {
	var newLines []LineResult

	// 已完成的 type 集合
	doneTypes := make(map[string]bool)
	for _, l := range pm.Lines {
		doneTypes[l.Type] = true
	}

	// 檢查三行
	for r := 0; r < GridSize; r++ {
		t := "row" + string(rune('0'+r))
		if !doneTypes[t] && pm.Cells[r][0] && pm.Cells[r][1] && pm.Cells[r][2] {
			newLines = append(newLines, LineResult{
				Type:  t,
				Cells: [3][2]int{{r, 0}, {r, 1}, {r, 2}},
			})
		}
	}

	// 檢查三列
	for c := 0; c < GridSize; c++ {
		t := "col" + string(rune('0'+c))
		if !doneTypes[t] && pm.Cells[0][c] && pm.Cells[1][c] && pm.Cells[2][c] {
			newLines = append(newLines, LineResult{
				Type:  t,
				Cells: [3][2]int{{0, c}, {1, c}, {2, c}},
			})
		}
	}

	// 檢查對角線（左上→右下）
	if !doneTypes["diag0"] && pm.Cells[0][0] && pm.Cells[1][1] && pm.Cells[2][2] {
		newLines = append(newLines, LineResult{
			Type:  "diag0",
			Cells: [3][2]int{{0, 0}, {1, 1}, {2, 2}},
		})
	}

	// 檢查對角線（右上→左下）
	if !doneTypes["diag1"] && pm.Cells[0][2] && pm.Cells[1][1] && pm.Cells[2][0] {
		newLines = append(newLines, LineResult{
			Type:  "diag1",
			Cells: [3][2]int{{0, 2}, {1, 1}, {2, 0}},
		})
	}

	return newLines
}

// todayUTC 取得今日 UTC 日期字串（YYYY-MM-DD）
func todayUTC() string {
	now := time.Now().UTC()
	return now.Format("2006-01-02")
}

// CalcLineReward 計算行/列/對角線獎勵（依投注等級）
// 業界設計：行獎勵 = betCost × 50（中等獎勵，讓玩家有成就感但不影響 RTP）
func CalcLineReward(betCost int) int {
	return betCost * 50
}

// CalcFullReward 計算整張地圖獎勵（依投注等級）
// 業界設計：全圖獎勵 = betCost × 500（大獎，讓玩家有強烈動機完成地圖）
func CalcFullReward(betCost int) int {
	return betCost * 500
}

// GetCellDef 取得格子定義
func GetCellDef(row, col int) *CellDef {
	if row < 0 || row >= GridSize || col < 0 || col >= GridSize {
		return nil
	}
	def := DefaultGrid[row][col]
	return &def
}
