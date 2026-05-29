// lucky_shark_spark_handler.go — T239 幸運鯊魚閃電魚
// 設計：Shark & Spark 機制（業界依據：BGaming Shark & Spark Hold & Win 2026-05-30 最新）
//       鯊魚閃電 + 珍珠倍率組合：場上每個目標獲得珍珠倍率（×1-×200）
//       閃電連鎖 6 條（每條 ×80.0），完美觸發（連鎖 ≥6）→ 全服 ×51.0 加成 102 秒（新里程碑）
//       業界依據：BGaming Shark & Spark Hold & Win（2026-05-30 最新發布）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type sharkSparkBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckySharkSparkManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *sharkSparkBoost
}

func newLuckySharkSparkManager() *luckySharkSparkManager {
	return &luckySharkSparkManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckySharkSparkFish(defID string) bool {
	return defID == "T239"
}

func (m *luckySharkSparkManager) getSharkSparkMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckySharkSparkManager) tryLuckySharkSparkFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(530 * time.Second)
	m.personalCD[p.ID] = now.Add(470 * time.Second)
	m.mu.Unlock()

	// 珍珠倍率：場上每個目標獲得隨機珍珠倍率
	pearlMults := []float64{1, 2, 5, 10, 20, 50, 100, 200}
	pearlWeights := []int{30, 25, 20, 12, 7, 4, 1, 1}

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckySharkSpark,
		Payload: map[string]interface{}{
			"event":        "shark_spark_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"chain_count":  6,
			"per_chain":    80.0,
			"global_target": 51.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("SHARK & SPARK! %s triggered the Shark Lightning! 6 chains x80.0! Pearl Multipliers activated!", p.GetDisplayName()), "critical", "#00BFFF")
	log.Printf("[LuckySharkSpark] %s triggered Shark & Spark fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)

		// 珍珠倍率分配
		g.mu.RLock()
		targetIDs := make([]string, 0, len(g.targets))
		for id := range g.targets {
			targetIDs = append(targetIDs, id)
		}
		g.mu.RUnlock()

		pearlAssignments := make(map[string]float64)
		for _, id := range targetIDs {
			// 加權隨機選珍珠倍率
			totalW := 0
			for _, w := range pearlWeights {
				totalW += w
			}
			r := rand.Intn(totalW)
			cum := 0
			chosen := pearlMults[0]
			for i, w := range pearlWeights {
				cum += w
				if r < cum {
					chosen = pearlMults[i]
					break
				}
			}
			pearlAssignments[id] = chosen
		}

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckySharkSpark,
			Payload: map[string]interface{}{
				"event":       "pearl_assigned",
				"pearl_count": len(pearlAssignments),
			},
		})

		time.Sleep(1 * time.Second)

		// 閃電連鎖 6 條
		chainCount := 0
		totalChainMult := 0.0
		for chain := 1; chain <= 6; chain++ {
			time.Sleep(400 * time.Millisecond)
			chainMult := 80.0
			totalChainMult += chainMult
			chainCount++

			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckySharkSpark,
				Payload: map[string]interface{}{
					"event":      "chain_strike",
					"chain_no":   chain,
					"chain_mult": chainMult,
				},
			})
		}

		// 計算珍珠獎勵
		totalPearlMult := 0.0
		for _, pm := range pearlAssignments {
			totalPearlMult += pm
		}

		totalMult := totalChainMult + totalPearlMult
		reward := int(totalMult * betCost)
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		// 完美觸發（連鎖 ≥6）→ 全服 ×51.0（新里程碑）
		globalBonus := 51.0
		globalDuration := 102

		m.mu.Lock()
		m.boost = &sharkSparkBoost{
			mult:      globalBonus,
			expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckySharkSpark,
			Payload: map[string]interface{}{
				"event":          "shark_spark_complete",
				"chain_count":    chainCount,
				"total_chain":    totalChainMult,
				"total_pearl":    totalPearlMult,
				"total_mult":     totalMult,
				"reward":         reward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
				"milestone":      "GLOBAL_51X",
			},
		})
		g.sendAnnounce(fmt.Sprintf("SHARK & SPARK COMPLETE! %s: %d chains + Pearl x%.1f! Total x%.1f! MILESTONE: Global x%.1f for %ds!", p.GetDisplayName(), chainCount, totalPearlMult, totalMult, globalBonus, globalDuration), "critical", "#00BFFF")
		log.Printf("[LuckySharkSpark] MILESTONE! %s: chains=%d, pearl=%.1f, total=%.1f, global=x%.1f", p.GetDisplayName(), chainCount, totalPearlMult, totalMult, globalBonus)
	}()

	return true
}
