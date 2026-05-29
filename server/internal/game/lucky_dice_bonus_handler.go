// lucky_dice_bonus_handler.go — T221 幸運骰子獎勵魚
// 設計：Dice Bonus 機制（BGaming「Shark & Spark Hold & Win」，2026-05-25 最新）
//       擲骰決定獎勵：1-3 點 ×50.0，4-5 點 ×150.0，6 點 ×300.0
//       連續擲骰 3 次，每次結果累加，全服 ×41.5 加成 83 秒
//       觸發率：0.001%（最稀有）；個人冷卻 315 秒；全服冷卻 375 秒
//       業界依據：BGaming「Shark & Spark Hold & Win」Dice Bonus（2026-05-25）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyDiceBonusManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
}

func newLuckyDiceBonusManager() *luckyDiceBonusManager {
	return &luckyDiceBonusManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyDiceBonusFish(defID string) bool {
	return defID == "T221"
}

func (m *luckyDiceBonusManager) tryLuckyDiceBonusFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(375 * time.Second)
	m.personalCD[p.ID] = now.Add(315 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_dice_bonus",
		Payload: map[string]interface{}{
			"event":        "dice_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"rolls":        3,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🎲✨ 骰子獎勵！%s 觸發骰子獎勵魚！連續擲骰 3 次，最高 ×300.0！", p.GetDisplayName()), "critical", "#FFD700")
	log.Printf("[LuckyDiceBonus] %s 觸發骰子獎勵魚（3 次擲骰）", p.GetDisplayName())

	go func() {
		totalMult := 0.0
		rollResults := []int{}

		for i := 0; i < 3; i++ {
			time.Sleep(1500 * time.Millisecond)
			dice := rand.Intn(6) + 1 // 1-6
			rollResults = append(rollResults, dice)

			var rollMult float64
			switch {
			case dice <= 3:
				rollMult = 50.0
			case dice <= 5:
				rollMult = 150.0
			default:
				rollMult = 300.0
			}
			totalMult += rollMult

			g.broadcast(protocol.Envelope{
				Type: "lucky_dice_bonus",
				Payload: map[string]interface{}{
					"event":      "dice_roll",
					"roll_index": i + 1,
					"dice_value": dice,
					"roll_mult":  rollMult,
					"total_mult": totalMult,
				},
			})
		}

		// 計算最終獎勵
		betCost := float64(p.GetBetDef().BetCost)
		reward := int(totalMult * betCost)

		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_dice_bonus",
			Payload: map[string]interface{}{
				"event":        "dice_settle",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"roll_results": rollResults,
				"total_mult":   totalMult,
				"reward":       reward,
			},
		})

		// 全服加成 ×41.5，持續 83 秒
		g.broadcast(protocol.Envelope{
			Type: "lucky_dice_bonus",
			Payload: map[string]interface{}{
				"event":         "global_boost",
				"global_mult":   41.5,
				"duration":      83,
				"trigger_name":  p.GetDisplayName(),
				"total_mult":    totalMult,
			},
		})
		g.sendAnnounce(fmt.Sprintf("🎲🌟 骰子結算！%s 擲出 %.0f 倍！全服 ×41.5 加成 83 秒！", p.GetDisplayName(), totalMult), "critical", "#FFD700")
		log.Printf("[LuckyDiceBonus] %s 骰子結算 %.0fx，全服 ×41.5 加成 83 秒", p.GetDisplayName(), totalMult)
	}()

	return true
}
