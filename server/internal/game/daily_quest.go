// Package game — 每日任務系統（DAY-345）
// 靈感來源：BGaming Quests（2026-05-27 發布）
// 設計：3個每日任務，完成獲得任務幣，任務幣可換取 BET 加成
// 重置：每日 UTC 00:00 自動重置
package game

import (
	"fmt"
	"time"
)

// QuestType 任務類型
type QuestType string

const (
	QuestKillTargets  QuestType = "kill_targets"  // 擊破目標物
	QuestComboReach   QuestType = "combo_reach"   // 達成連擊
	QuestTriggerBonus QuestType = "trigger_bonus" // 觸發 Bonus
)

// DailyQuest 每日任務定義
type DailyQuest struct {
	ID          string    `json:"id"`
	Type        QuestType `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Target      int       `json:"target"`   // 目標數量
	Reward      int       `json:"reward"`   // 完成獎勵（任務幣）
}

// PlayerQuestProgress 玩家任務進度
type PlayerQuestProgress struct {
	QuestID   string `json:"quest_id"`
	Progress  int    `json:"progress"`
	Completed bool   `json:"completed"`
	Claimed   bool   `json:"claimed"`
}

// DailyQuestSystem 每日任務系統
type DailyQuestSystem struct {
	quests    []DailyQuest
	// playerID -> questID -> progress
	progress  map[string]map[string]*PlayerQuestProgress
	// playerID -> 任務幣數量
	questCoins map[string]int
	// 今日日期（UTC）
	today     string
}

// newDailyQuestSystem 建立每日任務系統
func newDailyQuestSystem() *DailyQuestSystem {
	dqs := &DailyQuestSystem{
		progress:   make(map[string]map[string]*PlayerQuestProgress),
		questCoins: make(map[string]int),
		today:      todayUTC(),
	}
	dqs.quests = dqs.generateDailyQuests()
	return dqs
}

// todayUTC 取得今日 UTC 日期字串
func todayUTC() string {
	return time.Now().UTC().Format("2006-01-02")
}

// generateDailyQuests 生成今日任務（依日期決定難度）
func (dqs *DailyQuestSystem) generateDailyQuests() []DailyQuest {
	// 依星期幾決定任務難度
	weekday := int(time.Now().UTC().Weekday())

	// 基礎難度係數（週一最輕鬆，週末最難）
	difficulty := 1 + weekday/3 // 1, 1, 1, 2, 2, 2, 2

	return []DailyQuest{
		{
			ID:          "daily_kill",
			Type:        QuestKillTargets,
			Name:        "討伐任務",
			Description: fmt.Sprintf("今日擊破 %d 個目標物", 20*difficulty),
			Target:      20 * difficulty,
			Reward:      10 * difficulty,
		},
		{
			ID:          "daily_combo",
			Type:        QuestComboReach,
			Name:        "連擊挑戰",
			Description: fmt.Sprintf("達成 %d 連擊", 5*difficulty),
			Target:      5 * difficulty,
			Reward:      15 * difficulty,
		},
		{
			ID:          "daily_bonus",
			Type:        QuestTriggerBonus,
			Name:        "Bonus 探索",
			Description: "觸發 Bonus 遊戲 1 次",
			Target:      1,
			Reward:      20 * difficulty,
		},
	}
}

// checkReset 檢查是否需要重置（每日 UTC 00:00）
func (dqs *DailyQuestSystem) checkReset() bool {
	today := todayUTC()
	if today != dqs.today {
		dqs.today = today
		dqs.quests = dqs.generateDailyQuests()
		dqs.progress = make(map[string]map[string]*PlayerQuestProgress)
		// 任務幣不重置（累積制）
		return true
	}
	return false
}

// getOrCreateProgress 取得或建立玩家進度
func (dqs *DailyQuestSystem) getOrCreateProgress(playerID string) map[string]*PlayerQuestProgress {
	if _, ok := dqs.progress[playerID]; !ok {
		dqs.progress[playerID] = make(map[string]*PlayerQuestProgress)
		for _, q := range dqs.quests {
			dqs.progress[playerID][q.ID] = &PlayerQuestProgress{
				QuestID:   q.ID,
				Progress:  0,
				Completed: false,
				Claimed:   false,
			}
		}
	}
	return dqs.progress[playerID]
}

// OnKillTarget 擊破目標物時呼叫
// 回傳：是否有任務完成
func (dqs *DailyQuestSystem) OnKillTarget(playerID string) (completed bool, questName string, reward int) {
	dqs.checkReset()
	prog := dqs.getOrCreateProgress(playerID)

	for _, q := range dqs.quests {
		if q.Type != QuestKillTargets {
			continue
		}
		p := prog[q.ID]
		if p.Completed {
			continue
		}
		p.Progress++
		if p.Progress >= q.Target {
			p.Completed = true
			return true, q.Name, q.Reward
		}
	}
	return false, "", 0
}

// OnComboReach 達成連擊時呼叫
func (dqs *DailyQuestSystem) OnComboReach(playerID string, combo int) (completed bool, questName string, reward int) {
	dqs.checkReset()
	prog := dqs.getOrCreateProgress(playerID)

	for _, q := range dqs.quests {
		if q.Type != QuestComboReach {
			continue
		}
		p := prog[q.ID]
		if p.Completed {
			continue
		}
		if combo >= q.Target {
			p.Progress = combo
			p.Completed = true
			return true, q.Name, q.Reward
		}
		if combo > p.Progress {
			p.Progress = combo
		}
	}
	return false, "", 0
}

// OnTriggerBonus Bonus 觸發時呼叫
func (dqs *DailyQuestSystem) OnTriggerBonus(playerID string) (completed bool, questName string, reward int) {
	dqs.checkReset()
	prog := dqs.getOrCreateProgress(playerID)

	for _, q := range dqs.quests {
		if q.Type != QuestTriggerBonus {
			continue
		}
		p := prog[q.ID]
		if p.Completed {
			continue
		}
		p.Progress++
		if p.Progress >= q.Target {
			p.Completed = true
			return true, q.Name, q.Reward
		}
	}
	return false, "", 0
}

// ClaimReward 領取任務獎勵
// 回傳：領取的任務幣數量，-1 表示失敗
func (dqs *DailyQuestSystem) ClaimReward(playerID string, questID string) int {
	dqs.checkReset()
	prog := dqs.getOrCreateProgress(playerID)

	p, ok := prog[questID]
	if !ok || !p.Completed || p.Claimed {
		return -1
	}

	// 找到對應任務的獎勵
	for _, q := range dqs.quests {
		if q.ID == questID {
			p.Claimed = true
			dqs.questCoins[playerID] += q.Reward
			return q.Reward
		}
	}
	return -1
}

// GetQuestCoins 取得玩家任務幣數量
func (dqs *DailyQuestSystem) GetQuestCoins(playerID string) int {
	return dqs.questCoins[playerID]
}

// GetPlayerStatus 取得玩家任務狀態（用於發送給 Client）
func (dqs *DailyQuestSystem) GetPlayerStatus(playerID string) DailyQuestStatus {
	dqs.checkReset()
	prog := dqs.getOrCreateProgress(playerID)

	questStatuses := make([]QuestStatus, 0, len(dqs.quests))
	for _, q := range dqs.quests {
		p := prog[q.ID]
		questStatuses = append(questStatuses, QuestStatus{
			ID:          q.ID,
			Name:        q.Name,
			Description: q.Description,
			Target:      q.Target,
			Progress:    p.Progress,
			Completed:   p.Completed,
			Claimed:     p.Claimed,
			Reward:      q.Reward,
		})
	}

	return DailyQuestStatus{
		Quests:     questStatuses,
		QuestCoins: dqs.questCoins[playerID],
		ResetAt:    nextUTCMidnight(),
	}
}

// nextUTCMidnight 下次 UTC 午夜的 Unix 毫秒時間戳
func nextUTCMidnight() int64 {
	now := time.Now().UTC()
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	return midnight.UnixMilli()
}

// QuestStatus 任務狀態（發送給 Client）
type QuestStatus struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Target      int    `json:"target"`
	Progress    int    `json:"progress"`
	Completed   bool   `json:"completed"`
	Claimed     bool   `json:"claimed"`
	Reward      int    `json:"reward"`
}

// DailyQuestStatus 每日任務狀態（發送給 Client）
type DailyQuestStatus struct {
	Quests     []QuestStatus `json:"quests"`
	QuestCoins int           `json:"quest_coins"`
	ResetAt    int64         `json:"reset_at"` // UTC 毫秒時間戳
}

// SpendQuestCoins 消耗任務幣（購買商店道具時呼叫）
// 回傳：實際消耗的任務幣數量
func (dqs *DailyQuestSystem) SpendQuestCoins(playerID string, amount int) int {
	current := dqs.questCoins[playerID]
	if amount > current {
		amount = current
	}
	dqs.questCoins[playerID] -= amount
	return amount
}
