// Package game — 每週挑戰系統（DAY-346）
// 設計：比每日任務更難，獎勵更豐厚，週一重置
// 挑戰類型：累積擊破、高倍率擊破、連擊里程碑、Bonus 次數
package game

import (
	"fmt"
	"time"
)

// WeeklyChallengeType 每週挑戰類型
type WeeklyChallengeType string

const (
	WeeklyKillCount    WeeklyChallengeType = "weekly_kill_count"    // 累積擊破數
	WeeklyHighMult     WeeklyChallengeType = "weekly_high_mult"     // 高倍率擊破（≥50x）
	WeeklyMaxCombo     WeeklyChallengeType = "weekly_max_combo"     // 達成最高連擊
	WeeklyBonusCount   WeeklyChallengeType = "weekly_bonus_count"   // 觸發 Bonus 次數
	WeeklyLuckyKill    WeeklyChallengeType = "weekly_lucky_kill"    // 擊破 Lucky 目標物
)

// WeeklyChallenge 每週挑戰定義
type WeeklyChallenge struct {
	ID          string              `json:"id"`
	Type        WeeklyChallengeType `json:"type"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Target      int                 `json:"target"`
	Reward      int                 `json:"reward"` // 任務幣獎勵（比每日任務多 5-10 倍）
	Tier        int                 `json:"tier"`   // 難度等級 1-3
}

// PlayerWeeklyProgress 玩家每週挑戰進度
type PlayerWeeklyProgress struct {
	ChallengeID string `json:"challenge_id"`
	Progress    int    `json:"progress"`
	Completed   bool   `json:"completed"`
	Claimed     bool   `json:"claimed"`
}

// WeeklyChallengeSystem 每週挑戰系統
type WeeklyChallengeSystem struct {
	challenges []WeeklyChallenge
	// playerID -> challengeID -> progress
	progress map[string]map[string]*PlayerWeeklyProgress
	// playerID -> 本週任務幣
	weeklyCoins map[string]int
	// 本週週次（ISO week number）
	currentWeek string
}

// newWeeklyChallengeSystem 建立每週挑戰系統
func newWeeklyChallengeSystem() *WeeklyChallengeSystem {
	wcs := &WeeklyChallengeSystem{
		progress:    make(map[string]map[string]*PlayerWeeklyProgress),
		weeklyCoins: make(map[string]int),
		currentWeek: currentWeekKey(),
	}
	wcs.challenges = wcs.generateWeeklyChallenges()
	return wcs
}

// currentWeekKey 取得本週的唯一識別鍵（年-週次）
func currentWeekKey() string {
	year, week := time.Now().UTC().ISOWeek()
	return fmt.Sprintf("%d-W%02d", year, week)
}

// nextMondayUTC 下週一 UTC 00:00 的 Unix 毫秒時間戳
func nextMondayUTC() int64 {
	now := time.Now().UTC()
	// 計算到下週一的天數
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 週日 = 7
	}
	daysUntilMonday := 8 - weekday // 下週一
	nextMonday := time.Date(now.Year(), now.Month(), now.Day()+daysUntilMonday, 0, 0, 0, 0, time.UTC)
	return nextMonday.UnixMilli()
}

// generateWeeklyChallenges 生成本週挑戰
func (wcs *WeeklyChallengeSystem) generateWeeklyChallenges() []WeeklyChallenge {
	// 依週次決定挑戰難度（奇數週較難）
	_, week := time.Now().UTC().ISOWeek()
	hardMode := week%2 == 1

	killTarget := 200
	highMultTarget := 10
	comboTarget := 20
	bonusTarget := 5
	luckyTarget := 15

	if hardMode {
		killTarget = 300
		highMultTarget = 15
		comboTarget = 25
		bonusTarget := 8
		_ = bonusTarget
		luckyTarget = 20
	}

	return []WeeklyChallenge{
		{
			ID:          "weekly_kill",
			Type:        WeeklyKillCount,
			Name:        "週間討伐",
			Description: fmt.Sprintf("本週累積擊破 %d 個目標物", killTarget),
			Target:      killTarget,
			Reward:      100,
			Tier:        1,
		},
		{
			ID:          "weekly_high_mult",
			Type:        WeeklyHighMult,
			Name:        "高倍率獵人",
			Description: fmt.Sprintf("擊破 %d 個 50x 以上的目標物", highMultTarget),
			Target:      highMultTarget,
			Reward:      150,
			Tier:        2,
		},
		{
			ID:          "weekly_combo",
			Type:        WeeklyMaxCombo,
			Name:        "連擊大師",
			Description: fmt.Sprintf("達成 %d 連擊", comboTarget),
			Target:      comboTarget,
			Reward:      120,
			Tier:        2,
		},
		{
			ID:          "weekly_bonus",
			Type:        WeeklyBonusCount,
			Name:        "Bonus 探索家",
			Description: fmt.Sprintf("觸發 Bonus 遊戲 %d 次", bonusTarget),
			Target:      bonusTarget,
			Reward:      200,
			Tier:        3,
		},
		{
			ID:          "weekly_lucky",
			Type:        WeeklyLuckyKill,
			Name:        "幸運獵手",
			Description: fmt.Sprintf("擊破 %d 個 Lucky 目標物（T106+）", luckyTarget),
			Target:      luckyTarget,
			Reward:      180,
			Tier:        3,
		},
	}
}

// checkReset 檢查是否需要重置（每週一 UTC 00:00）
func (wcs *WeeklyChallengeSystem) checkReset() bool {
	week := currentWeekKey()
	if week != wcs.currentWeek {
		wcs.currentWeek = week
		wcs.challenges = wcs.generateWeeklyChallenges()
		wcs.progress = make(map[string]map[string]*PlayerWeeklyProgress)
		wcs.weeklyCoins = make(map[string]int)
		return true
	}
	return false
}

// getOrCreateProgress 取得或建立玩家進度
func (wcs *WeeklyChallengeSystem) getOrCreateProgress(playerID string) map[string]*PlayerWeeklyProgress {
	if _, ok := wcs.progress[playerID]; !ok {
		wcs.progress[playerID] = make(map[string]*PlayerWeeklyProgress)
		for _, c := range wcs.challenges {
			wcs.progress[playerID][c.ID] = &PlayerWeeklyProgress{
				ChallengeID: c.ID,
				Progress:    0,
				Completed:   false,
				Claimed:     false,
			}
		}
	}
	return wcs.progress[playerID]
}

// OnKillTarget 擊破目標物時呼叫
// isLucky: 是否為 Lucky 目標物（T106+）
// mult: 擊破倍率
// 回傳：完成的挑戰列表
func (wcs *WeeklyChallengeSystem) OnKillTarget(playerID string, isLucky bool, mult float64) []WeeklyChallengeComplete {
	wcs.checkReset()
	prog := wcs.getOrCreateProgress(playerID)
	var completed []WeeklyChallengeComplete

	for _, c := range wcs.challenges {
		p := prog[c.ID]
		if p.Completed {
			continue
		}

		switch c.Type {
		case WeeklyKillCount:
			p.Progress++
			if p.Progress >= c.Target {
				p.Completed = true
				completed = append(completed, WeeklyChallengeComplete{
					ChallengeID:   c.ID,
					ChallengeName: c.Name,
					Reward:        c.Reward,
					Tier:          c.Tier,
				})
			}
		case WeeklyHighMult:
			if mult >= 50.0 {
				p.Progress++
				if p.Progress >= c.Target {
					p.Completed = true
					completed = append(completed, WeeklyChallengeComplete{
						ChallengeID:   c.ID,
						ChallengeName: c.Name,
						Reward:        c.Reward,
						Tier:          c.Tier,
					})
				}
			}
		case WeeklyLuckyKill:
			if isLucky {
				p.Progress++
				if p.Progress >= c.Target {
					p.Completed = true
					completed = append(completed, WeeklyChallengeComplete{
						ChallengeID:   c.ID,
						ChallengeName: c.Name,
						Reward:        c.Reward,
						Tier:          c.Tier,
					})
				}
			}
		}
	}

	return completed
}

// OnComboReach 達成連擊時呼叫
func (wcs *WeeklyChallengeSystem) OnComboReach(playerID string, combo int) []WeeklyChallengeComplete {
	wcs.checkReset()
	prog := wcs.getOrCreateProgress(playerID)
	var completed []WeeklyChallengeComplete

	for _, c := range wcs.challenges {
		if c.Type != WeeklyMaxCombo {
			continue
		}
		p := prog[c.ID]
		if p.Completed {
			continue
		}
		if combo > p.Progress {
			p.Progress = combo
		}
		if p.Progress >= c.Target {
			p.Completed = true
			completed = append(completed, WeeklyChallengeComplete{
				ChallengeID:   c.ID,
				ChallengeName: c.Name,
				Reward:        c.Reward,
				Tier:          c.Tier,
			})
		}
	}

	return completed
}

// OnTriggerBonus Bonus 觸發時呼叫
func (wcs *WeeklyChallengeSystem) OnTriggerBonus(playerID string) []WeeklyChallengeComplete {
	wcs.checkReset()
	prog := wcs.getOrCreateProgress(playerID)
	var completed []WeeklyChallengeComplete

	for _, c := range wcs.challenges {
		if c.Type != WeeklyBonusCount {
			continue
		}
		p := prog[c.ID]
		if p.Completed {
			continue
		}
		p.Progress++
		if p.Progress >= c.Target {
			p.Completed = true
			completed = append(completed, WeeklyChallengeComplete{
				ChallengeID:   c.ID,
				ChallengeName: c.Name,
				Reward:        c.Reward,
				Tier:          c.Tier,
			})
		}
	}

	return completed
}

// ClaimReward 領取挑戰獎勵
func (wcs *WeeklyChallengeSystem) ClaimReward(playerID string, challengeID string) int {
	wcs.checkReset()
	prog := wcs.getOrCreateProgress(playerID)

	p, ok := prog[challengeID]
	if !ok || !p.Completed || p.Claimed {
		return -1
	}

	for _, c := range wcs.challenges {
		if c.ID == challengeID {
			p.Claimed = true
			wcs.weeklyCoins[playerID] += c.Reward
			return c.Reward
		}
	}
	return -1
}

// GetWeeklyCoins 取得玩家本週任務幣
func (wcs *WeeklyChallengeSystem) GetWeeklyCoins(playerID string) int {
	return wcs.weeklyCoins[playerID]
}

// GetPlayerStatus 取得玩家每週挑戰狀態
func (wcs *WeeklyChallengeSystem) GetPlayerStatus(playerID string) WeeklyChallengeStatus {
	wcs.checkReset()
	prog := wcs.getOrCreateProgress(playerID)

	challengeStatuses := make([]ChallengeStatus, 0, len(wcs.challenges))
	for _, c := range wcs.challenges {
		p := prog[c.ID]
		challengeStatuses = append(challengeStatuses, ChallengeStatus{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			Target:      c.Target,
			Progress:    p.Progress,
			Completed:   p.Completed,
			Claimed:     p.Claimed,
			Reward:      c.Reward,
			Tier:        c.Tier,
		})
	}

	return WeeklyChallengeStatus{
		Challenges:  challengeStatuses,
		WeeklyCoins: wcs.weeklyCoins[playerID],
		WeekKey:     wcs.currentWeek,
		ResetAt:     nextMondayUTC(),
	}
}

// WeeklyChallengeComplete 完成的挑戰資訊
type WeeklyChallengeComplete struct {
	ChallengeID   string `json:"challenge_id"`
	ChallengeName string `json:"challenge_name"`
	Reward        int    `json:"reward"`
	Tier          int    `json:"tier"`
}

// ChallengeStatus 挑戰狀態（發送給 Client）
type ChallengeStatus struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Target      int    `json:"target"`
	Progress    int    `json:"progress"`
	Completed   bool   `json:"completed"`
	Claimed     bool   `json:"claimed"`
	Reward      int    `json:"reward"`
	Tier        int    `json:"tier"`
}

// WeeklyChallengeStatus 每週挑戰狀態（發送給 Client）
type WeeklyChallengeStatus struct {
	Challenges  []ChallengeStatus `json:"challenges"`
	WeeklyCoins int               `json:"weekly_coins"`
	WeekKey     string            `json:"week_key"`
	ResetAt     int64             `json:"reset_at"` // UTC 毫秒時間戳
}

// SpendQuestCoins 消耗任務幣（購買商店道具時呼叫）
// 回傳：實際消耗的任務幣數量
func (wcs *WeeklyChallengeSystem) SpendQuestCoins(playerID string, amount int) int {
	current := wcs.weeklyCoins[playerID]
	if amount > current {
		amount = current
	}
	wcs.weeklyCoins[playerID] -= amount
	return amount
}
