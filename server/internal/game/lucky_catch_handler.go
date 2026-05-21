// lucky_catch_handler.go — 幸運捕獲系統 handler（DAY-119）
// 業界依據：betway.com Lucky Catch Pick and Win（2026-04）確認「即時獎勵」機制
// 讓玩家留存率提升 22%；Ice Fishing Live（Evolution）的隨機 Bonus 觸發是 2026 年最熱門機制
package game

import (
	"log"
	"math/rand"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// luckyCatchCooldown 幸運捕獲冷卻時間（每個玩家 60 秒）
const luckyCatchCooldown = 60 * time.Second

// luckyCatchTargetNames 目標物名稱對照表
var luckyCatchTargetNames = map[string]string{
	"T001": "吉伊卡哇草",
	"T002": "綠色小蟲",
	"T003": "紅色小蟲",
	"T004": "藍色小蟲",
	"T005": "布丁怪",
	"T006": "金魚",
	"T101": "擬態怪物",
	"T102": "寶箱怪",
	"T103": "流星",
	"T104": "金草",
	"T105": "金幣魚",
	"B001": "吉伊卡哇大魔王",
}

// tryLuckyCatch 嘗試觸發幸運捕獲（由 notifyStreakKill 和 handleKill 呼叫）
// triggerType: "streak"（連擊觸發）/ "weather"（天氣觸發）/ "festival"（節日觸發）
// 回傳是否觸發
func (g *Game) tryLuckyCatch(p *player.Player, triggerType string) bool {
	if p == nil {
		return false
	}

	// 冷卻檢查
	if time.Since(p.LastLuckyCatchAt) < luckyCatchCooldown {
		return false
	}

	// 計算觸發機率
	chance := 0.0
	icon := "🍀"
	switch triggerType {
	case "streak":
		// 連擊≥10 時有 3% 機率
		if p.Streak == nil {
			return false
		}
		snap := p.Streak.GetSnapshot()
		if snap.Current < 10 {
			return false
		}
		chance = 0.03
		icon = "⚡🍀"
	case "weather":
		// 天氣加成期間有 5% 機率（每次擊破時）
		chance = 0.05
		icon = "🌟🍀"
	case "festival":
		// 節日期間有 8% 機率（每次擊破時）
		chance = 0.08
		icon = "🎊🍀"
	default:
		return false
	}

	// 機率判定
	if rand.Float64() >= chance {
		return false
	}

	// 收集場上可用目標（排除 BOSS）
	g.mu.RLock()
	type candidateTarget struct {
		id         string
		defID      string
		multiplier float64
	}
	candidates := make([]candidateTarget, 0, len(g.Targets))
	for id, t := range g.Targets {
		if t.DefID == "B001" {
			continue // 不捕獲 BOSS
		}
		if !t.IsAlive {
			continue
		}
		candidates = append(candidates, candidateTarget{
			id:         id,
			defID:      t.DefID,
			multiplier: t.Multiplier,
		})
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		return false
	}

	// 隨機選取一個目標
	chosen := candidates[rand.Intn(len(candidates))]

	// 幸運加成倍率：2.0-5.0x（隨機）
	bonusMult := 2.0 + rand.Float64()*3.0

	// 計算獎勵
	baseReward := int(float64(p.BetLevel*10) * chosen.multiplier)
	finalReward := int(float64(baseReward) * bonusMult)

	// 從場上移除目標
	g.mu.Lock()
	delete(g.Targets, chosen.id)
	g.mu.Unlock()

	// 發放獎勵
	p.AddCoins(finalReward)
	p.LastLuckyCatchAt = time.Now()

	// 取得目標名稱
	targetName, ok := luckyCatchTargetNames[chosen.defID]
	if !ok {
		targetName = chosen.defID
	}

	// 廣播幸運捕獲事件（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyCatch,
		Payload: ws.LuckyCatchPayload{
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			TargetDefID: chosen.defID,
			TargetName:  targetName,
			Multiplier:  chosen.multiplier,
			BonusMult:   bonusMult,
			Reward:      finalReward,
			TriggerType: triggerType,
			Icon:        icon,
		},
	})

	log.Printf("[LuckyCatch] player=%s triggered %s via %s — target=%s mult=%.1f bonus=%.1fx reward=%d",
		p.ID, p.DisplayName, triggerType, chosen.defID, chosen.multiplier, bonusMult, finalReward)

	// 動態牆：幸運捕獲（≥50x 才廣播）
	if chosen.multiplier*bonusMult >= 50 {
		go g.notifyFeedMegaWin(p, chosen.multiplier*bonusMult, finalReward)
	}

	return true
}
