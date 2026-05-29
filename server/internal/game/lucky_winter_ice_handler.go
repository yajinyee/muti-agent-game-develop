// lucky_winter_ice_handler.go — T240 幸運冬季冰釣魚
// 設計：Winter Ice Fishing 機制（業界依據：BGaming Winter Fishing Club 2026-01）
//       冰下魚群 + 53格輪盤：Leaf×2（1:1-10:1）、Lil'Blues（3x-100x）、Big Oranges（4x-200x）、Huge Reds（10x-500x）
//       最高單次 ≥300x → 全服 ×51.5 加成 103 秒
//       業界依據：BGaming Winter Fishing Club（2026-01）+ Evolution Ice Fishing Live（2025-2026）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type winterIceBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyWinterIceManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *winterIceBoost
}

func newLuckyWinterIceManager() *luckyWinterIceManager {
	return &luckyWinterIceManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyWinterIceFish(defID string) bool {
	return defID == "T240"
}

func (m *luckyWinterIceManager) getWinterIceMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyWinterIceManager) tryLuckyWinterIceFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(540 * time.Second)
	m.personalCD[p.ID] = now.Add(480 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyWinterIce,
		Payload: map[string]interface{}{
			"event":        "winter_ice_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"spins":        3,
			"global_target": 51.5,
		},
	})
	g.sendAnnounce(fmt.Sprintf("WINTER ICE FISHING! %s triggered the 53-segment wheel! 3 spins! Max x500!", p.GetDisplayName()), "critical", "#87CEEB")
	log.Printf("[LuckyWinterIce] %s triggered Winter Ice Fishing fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)

		// 53格輪盤：Leaf1×23, Leaf2×23, Lil'Blues×4, BigOranges×2, HugeReds×1
		type segment struct {
			name    string
			count   int
			minMult float64
			maxMult float64
		}
		segments := []segment{
			{"Leaf1", 23, 1.0, 10.0},
			{"Leaf2", 23, 1.0, 10.0},
			{"LilBlues", 4, 3.0, 100.0},
			{"BigOranges", 2, 4.0, 200.0},
			{"HugeReds", 1, 10.0, 500.0},
		}

		totalMult := 0.0
		bestSpin := 0.0
		spinResults := make([]map[string]interface{}, 0, 3)

		for spin := 1; spin <= 3; spin++ {
			time.Sleep(600 * time.Millisecond)

			// 隨機選 segment（加權）
			totalSegs := 53
			r := rand.Intn(totalSegs)
			cum := 0
			chosen := segments[0]
			for _, seg := range segments {
				cum += seg.count
				if r < cum {
					chosen = seg
					break
				}
			}

			// 在 segment 範圍內隨機倍率
			spinMult := chosen.minMult + rand.Float64()*(chosen.maxMult-chosen.minMult)
			totalMult += spinMult
			if spinMult > bestSpin {
				bestSpin = spinMult
			}

			spinResults = append(spinResults, map[string]interface{}{
				"spin_no":  spin,
				"segment":  chosen.name,
				"spin_mult": spinMult,
			})

			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyWinterIce,
				Payload: map[string]interface{}{
					"event":     "wheel_spin",
					"spin_no":   spin,
					"segment":   chosen.name,
					"spin_mult": spinMult,
				},
			})
		}

		reward := int(totalMult * betCost)
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		// 最高單次 ≥300x → 全服 ×51.5
		globalBonus := 51.5
		globalDuration := 103

		m.mu.Lock()
		m.boost = &winterIceBoost{
			mult:      globalBonus,
			expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyWinterIce,
			Payload: map[string]interface{}{
				"event":          "winter_ice_complete",
				"spins":          3,
				"best_spin":      bestSpin,
				"total_mult":     totalMult,
				"reward":         reward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
			},
		})
		g.sendAnnounce(fmt.Sprintf("WINTER ICE COMPLETE! %s: best spin x%.1f! Total x%.1f! Global x%.1f for %ds!", p.GetDisplayName(), bestSpin, totalMult, globalBonus, globalDuration), "critical", "#87CEEB")
		log.Printf("[LuckyWinterIce] %s: best=%.1f, total=%.1f, global=x%.1f", p.GetDisplayName(), bestSpin, totalMult, globalBonus)
	}()

	return true
}
