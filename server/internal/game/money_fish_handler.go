// money_fish_handler.go — 金幣魚王即時獎勵 handler（DAY-162）
// 業界依據：King of Ocean 2026（Galaxsys）「Money Fish trigger instant payouts」
// 擊破金幣魚王後立即給予玩家一筆即時獎勵（betLevel × 20-50 隨機），
// 不走正常的 kill 倍率計算，是「保底即時獎勵」型特殊目標
// 設計：讓玩家在任何 betLevel 都能獲得有感的即時金幣，製造「爆金幣」的爽感
// 全服廣播讓其他玩家看到「有人的金幣魚王爆出大量金幣」
package game

import (
	"fmt"
	"log"
	"math/rand"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// MoneyFishMinMult 金幣魚王即時獎勵最小倍率（betLevel 的倍數）
const MoneyFishMinMult = 20

// MoneyFishMaxMult 金幣魚王即時獎勵最大倍率（betLevel 的倍數）
const MoneyFishMaxMult = 50

// MoneyFishAnnounceThreshold 全服公告門檻（即時獎勵 >= betLevel × 40 才公告）
const MoneyFishAnnounceThreshold = 40

// isMoneyFish 判斷是否為金幣魚王
func isMoneyFish(defID string) bool {
	return defID == "T122"
}

// notifyMoneyFishKill 擊破金幣魚王後觸發即時獎勵（由 handleKill 呼叫）
// 注意：金幣魚王的即時獎勵是「額外獎勵」，在正常 kill 獎勵之外額外發放
func (g *Game) notifyMoneyFishKill(p *player.Player, instanceID string, x, y float64) {
	// 計算即時獎勵：betLevel × 隨機倍率（20-50）
	mult := MoneyFishMinMult + rand.Intn(MoneyFishMaxMult-MoneyFishMinMult+1)
	instantReward := p.BetLevel * mult

	// 直接發放即時獎勵（不走正常 kill 計算）
	p.AddReward(instantReward)

	log.Printf("[MoneyFish] player=%s instant reward=%d (betLevel=%d × %d)",
		p.ID, instantReward, p.BetLevel, mult)

	// 廣播即時獎勵（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgMoneyFishReward,
		Payload: ws.MoneyFishRewardPayload{
			TriggerID:     instanceID,
			TriggerX:      x,
			TriggerY:      y,
			KillerID:      p.ID,
			KillerName:    p.DisplayName,
			InstantReward: instantReward,
			MultUsed:      mult,
			BetLevel:      p.BetLevel,
		},
	})

	// 全服公告（高倍率才公告）
	if mult >= MoneyFishAnnounceThreshold {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgAnnouncement,
			Payload: map[string]interface{}{
				"event_type": "money_fish_king",
				"message":    fmt.Sprintf("💰 %s 擊破金幣魚王！即時獲得 %d 金幣（×%d）！", p.DisplayName, instantReward, mult),
				"color":      "#FFD700",
				"duration":   4.0,
				"priority":   3,
			},
		})
	}
}
