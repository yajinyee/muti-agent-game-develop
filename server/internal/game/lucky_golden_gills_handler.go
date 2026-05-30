// lucky_golden_gills_handler.go — T252 幸運黃金鰓魚
// 設計：Golden Gills Jackpot Respin 機制（Atomic Slot Lab 2026）
//       磁力連鎖 + 4 層 Jackpot（Mini/Minor/Major/Grand）
//       Respin 期間收集 Jackpot 符號，填滿對應層 → 觸發 Jackpot
//       Grand Jackpot 觸發 → 全服 ×58.0 加成 116 秒
//       業界依據：Atomic Slot Lab「Golden Gills」Jackpot Respin + 磁力連鎖（2026-02）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type goldenGillsBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyGoldenGillsManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *goldenGillsBoost
}

func newLuckyGoldenGillsManager() *luckyGoldenGillsManager {
	return &luckyGoldenGillsManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyGoldenGillsFish(defID string) bool {
	return defID == "T252"
}

func (m *luckyGoldenGillsManager) getGoldenGillsMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyGoldenGillsManager) tryLuckyGoldenGillsFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(660 * time.Second)
	m.personalCD[p.ID] = now.Add(600 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyGoldenGills,
		Payload: map[string]interface{}{
			"event":        "golden_gills_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"jackpot_tiers": []string{"Mini", "Minor", "Major", "Grand"},
		},
	})
	g.sendAnnounce(fmt.Sprintf("GOLDEN GILLS! %s activated Jackpot Respin system! 4-tier Jackpot: Mini/Minor/Major/Grand!", p.GetDisplayName()), "critical", "#FFD700")
	log.Printf("[LuckyGoldenGills] %s triggered Golden Gills fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		totalReward := 0

		// Jackpot 層級定義
		jackpotMults := map[string]float64{
			"Mini":  50.0,
			"Minor": 150.0,
			"Major": 500.0,
			"Grand": 2000.0,
		}
		jackpotOrder := []string{"Mini", "Minor", "Major", "Grand"}

		// 磁力連鎖 Respin：6 次，每次可能觸發 Jackpot
		triggeredJackpots := []string{}
		for respin := 1; respin <= 6; respin++ {
			time.Sleep(400 * time.Millisecond)
			// 隨機觸發 Jackpot（機率遞增）
			roll := rand.Float64()
			var triggered string
			if roll < 0.15 {
				triggered = "Grand"
			} else if roll < 0.30 {
				triggered = "Major"
			} else if roll < 0.55 {
				triggered = "Minor"
			} else if roll < 0.80 {
				triggered = "Mini"
			}

			spinReward := int(75.0 * betCost) // 基礎旋轉獎勵
			if triggered != "" {
				jackpotReward := int(jackpotMults[triggered] * betCost)
				spinReward += jackpotReward
				triggeredJackpots = append(triggeredJackpots, triggered)
				g.broadcast(protocol.Envelope{
					Type: protocol.MsgLuckyGoldenGills,
					Payload: map[string]interface{}{
						"event":          "jackpot_triggered",
						"respin_no":      respin,
						"jackpot_tier":   triggered,
						"jackpot_reward": jackpotReward,
					},
				})
			}
			totalReward += spinReward
			g.mu.Lock()
			p.Coins += spinReward
			g.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyGoldenGills,
				Payload: map[string]interface{}{
					"event":     "respin",
					"respin_no": respin,
					"reward":    spinReward,
				},
			})
		}

		// 確保至少觸發 Grand Jackpot 一次（保底）
		if len(triggeredJackpots) == 0 {
			grandReward := int(jackpotMults["Grand"] * betCost)
			totalReward += grandReward
			g.mu.Lock()
			p.Coins += grandReward
			g.mu.Unlock()
			triggeredJackpots = append(triggeredJackpots, "Grand")
		}

		// Grand Jackpot 觸發 → 全服 ×58.0 加成 116 秒
		globalBonus := 58.0
		globalDuration := 116
		m.mu.Lock()
		m.boost = &goldenGillsBoost{
			mult:      globalBonus,
			expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
		}
		m.mu.Unlock()

		_ = jackpotOrder // 保留供未來使用
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyGoldenGills,
			Payload: map[string]interface{}{
				"event":             "golden_gills_complete",
				"triggered_jackpots": triggeredJackpots,
				"total_reward":      totalReward,
				"global_bonus":      globalBonus,
				"global_seconds":    globalDuration,
			},
		})
		g.sendAnnounce(fmt.Sprintf("GOLDEN GILLS COMPLETE! %s: Jackpots=%v! GLOBAL x%.1f for %ds!", p.GetDisplayName(), triggeredJackpots, globalBonus, globalDuration), "critical", "#FFD700")
		log.Printf("[LuckyGoldenGills] %s: jackpots=%v, global=x%.1f", p.GetDisplayName(), triggeredJackpots, globalBonus)
	}()

	return true
}
