// lucky_wild_collector_handler.go — T244 幸運野生收集魚
// 設計：Wild Collector 機制（BGaming Big Boat Big Catch 升級版）
//       Wild 符號收集：每 4 個 Wild → 額外旋轉（×2→×3→×10 倍率升級）
//       最高 10 次額外旋轉，每次旋轉倍率遞增，全服 ×54.0 加成 108 秒
//       業界依據：BGaming Big Boat Big Catch（2026-03）Wild Collector 機制
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type wildCollectorBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyWildCollectorManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *wildCollectorBoost
}

func newLuckyWildCollectorManager() *luckyWildCollectorManager {
	return &luckyWildCollectorManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyWildCollectorFish(defID string) bool {
	return defID == "T244"
}

func (m *luckyWildCollectorManager) getWildCollectorMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyWildCollectorManager) tryLuckyWildCollectorFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(580 * time.Second)
	m.personalCD[p.ID] = now.Add(520 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyWildCollector,
		Payload: map[string]interface{}{
			"event":        "wild_collector_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"max_spins":    10,
			"mult_stages":  []float64{2.0, 3.0, 10.0},
		},
	})
	g.sendAnnounce(fmt.Sprintf("WILD COLLECTOR! %s activated Wild Collector! Collecting Wilds for bonus spins x2→x3→x10!", p.GetDisplayName()), "critical", "#FFD700")
	log.Printf("[LuckyWildCollector] %s triggered Wild Collector fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)

		// Wild 收集階段：3 組（每組 4 個 Wild）
		wildStages := []struct {
			wilds int
			mult  float64
			spins int
		}{
			{4, 2.0, 3},
			{4, 3.0, 3},
			{4, 10.0, 4},
		}

		totalReward := 0
		totalMult := 0.0

		for stageIdx, stage := range wildStages {
			time.Sleep(800 * time.Millisecond)

			// 廣播 Wild 收集進度
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyWildCollector,
				Payload: map[string]interface{}{
					"event":       "wild_collected",
					"stage":       stageIdx + 1,
					"wild_count":  stage.wilds,
					"spin_mult":   stage.mult,
					"bonus_spins": stage.spins,
				},
			})

			time.Sleep(400 * time.Millisecond)

			// 執行額外旋轉
			for spin := 1; spin <= stage.spins; spin++ {
				time.Sleep(300 * time.Millisecond)
				spinMult := stage.mult
				spinReward := int(spinMult * betCost * 5.0)
				totalReward += spinReward
				totalMult += spinMult

				g.mu.Lock()
				p.Coins += spinReward
				g.mu.Unlock()

				g.broadcast(protocol.Envelope{
					Type: protocol.MsgLuckyWildCollector,
					Payload: map[string]interface{}{
						"event":      "bonus_spin",
						"stage":      stageIdx + 1,
						"spin_no":    spin,
						"spin_mult":  spinMult,
						"spin_reward": spinReward,
					},
				})
			}
		}

		// 全服 ×54.0 加成 108 秒
		globalBonus := 54.0
		globalDuration := 108

		m.mu.Lock()
		m.boost = &wildCollectorBoost{
			mult:      globalBonus,
			expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyWildCollector,
			Payload: map[string]interface{}{
				"event":          "wild_collector_complete",
				"total_spins":    10,
				"total_mult":     totalMult,
				"total_reward":   totalReward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
			},
		})
		g.sendAnnounce(fmt.Sprintf("WILD COLLECTOR COMPLETE! %s collected all Wilds! Total x%.1f! GLOBAL x%.1f for %ds!", p.GetDisplayName(), totalMult, globalBonus, globalDuration), "critical", "#FFD700")
		log.Printf("[LuckyWildCollector] %s: total_mult=%.1f, reward=%d, global=x%.1f", p.GetDisplayName(), totalMult, totalReward, globalBonus)
	}()

	return true
}
