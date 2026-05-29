// lucky_fishing_time_wheel_handler.go — T242 幸運釣魚時間魚
// 設計：Fishing Time Wheel 機制（業界依據：BGaming Fishing Time 2026-04）
//       命運輪盤 + 倍率疊加：Wheel of Fortune 機制，每次旋轉倍率疊加
//       5 次旋轉（最高 ×10000），完美收割（≥5000x）→ 全服 ×52.5 加成 105 秒
//       業界依據：BGaming Fishing Time（2026-04）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type fishingTimeWheelBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyFishingTimeWheelManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *fishingTimeWheelBoost
}

func newLuckyFishingTimeWheelManager() *luckyFishingTimeWheelManager {
	return &luckyFishingTimeWheelManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyFishingTimeWheelFish(defID string) bool {
	return defID == "T242"
}

func (m *luckyFishingTimeWheelManager) getFishingTimeWheelMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyFishingTimeWheelManager) tryLuckyFishingTimeWheelFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(560 * time.Second)
	m.personalCD[p.ID] = now.Add(500 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyFishingTimeWheel,
		Payload: map[string]interface{}{
			"event":        "fishing_time_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"spins":        5,
			"max_mult":     10000.0,
			"global_target": 52.5,
		},
	})
	g.sendAnnounce(fmt.Sprintf("FISHING TIME WHEEL! %s triggered the Fortune Wheel! 5 spins! Max x10000!", p.GetDisplayName()), "critical", "#FFD700")
	log.Printf("[LuckyFishingTimeWheel] %s triggered Fishing Time Wheel fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)

		// 輪盤倍率表（5次旋轉，倍率疊加）
		wheelMults := []float64{10, 20, 50, 100, 200, 500, 1000, 2000, 5000, 10000}
		wheelWeights := []int{25, 20, 18, 15, 10, 6, 3, 2, 1, 0} // 最後一格極低機率

		totalMult := 0.0
		bestSpin := 0.0

		for spin := 1; spin <= 5; spin++ {
			time.Sleep(700 * time.Millisecond)

			totalW := 0
			for _, w := range wheelWeights {
				totalW += w
			}
			r := rand.Intn(totalW)
			cum := 0
			chosen := wheelMults[0]
			for i, w := range wheelWeights {
				cum += w
				if r < cum {
					chosen = wheelMults[i]
					break
				}
			}

			totalMult += chosen
			if chosen > bestSpin {
				bestSpin = chosen
			}

			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyFishingTimeWheel,
				Payload: map[string]interface{}{
					"event":     "wheel_spin",
					"spin_no":   spin,
					"spin_mult": chosen,
					"total_so_far": totalMult,
				},
			})
		}

		reward := int(totalMult * betCost)
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		// 全服 ×52.5
		globalBonus := 52.5
		globalDuration := 105

		m.mu.Lock()
		m.boost = &fishingTimeWheelBoost{
			mult:      globalBonus,
			expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyFishingTimeWheel,
			Payload: map[string]interface{}{
				"event":          "fishing_time_complete",
				"spins":          5,
				"best_spin":      bestSpin,
				"total_mult":     totalMult,
				"reward":         reward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
			},
		})
		g.sendAnnounce(fmt.Sprintf("FISHING TIME COMPLETE! %s: best x%.1f! Total x%.1f! Global x%.1f for %ds!", p.GetDisplayName(), bestSpin, totalMult, globalBonus, globalDuration), "critical", "#FFD700")
		log.Printf("[LuckyFishingTimeWheel] %s: best=%.1f, total=%.1f, global=x%.1f", p.GetDisplayName(), bestSpin, totalMult, globalBonus)
	}()

	return true
}
