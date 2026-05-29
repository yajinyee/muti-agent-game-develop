// lucky_atlantis_frenzy_handler.go — T241 幸運大西洋狂潮魚
// 設計：Big Atlantis Frenzy 機制（業界依據：BGaming Big Atlantis Frenzy 2025-2026）
//       亞特蘭提斯爆炸 + 連鎖消除：Fish 符號隨機獎勵（×5-×500）+ Buy Chance 三倍機率
//       連鎖消除 ≥5 波 → 全服 ×52.0 加成 104 秒
//       業界依據：BGaming Big Atlantis Frenzy（2025-2026）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type atlantisFrenzyBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyAtlantisFrenzyManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *atlantisFrenzyBoost
}

func newLuckyAtlantisFrenzyManager() *luckyAtlantisFrenzyManager {
	return &luckyAtlantisFrenzyManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyAtlantisFrenzyFish(defID string) bool {
	return defID == "T241"
}

func (m *luckyAtlantisFrenzyManager) getAtlantisFrenzyMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyAtlantisFrenzyManager) tryLuckyAtlantisFrenzyFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(550 * time.Second)
	m.personalCD[p.ID] = now.Add(490 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyAtlantisFrenzy,
		Payload: map[string]interface{}{
			"event":        "atlantis_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"waves":        7,
			"global_target": 52.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("BIG ATLANTIS FRENZY! %s triggered the Atlantis Explosion! 7 cascade waves! Fish symbols x5-x500!", p.GetDisplayName()), "critical", "#1E90FF")
	log.Printf("[LuckyAtlantisFrenzy] %s triggered Big Atlantis Frenzy fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)

		// Fish 符號倍率表（加權）
		fishMults := []float64{5, 10, 20, 50, 100, 200, 500}
		fishWeights := []int{30, 25, 20, 12, 7, 4, 2}

		totalMult := 0.0
		waveCount := 0

		for wave := 1; wave <= 7; wave++ {
			time.Sleep(500 * time.Millisecond)

			// 每波 2-4 個 Fish 符號
			fishCount := 2 + rand.Intn(3)
			waveMult := 0.0
			for i := 0; i < fishCount; i++ {
				totalW := 0
				for _, w := range fishWeights {
					totalW += w
				}
				r := rand.Intn(totalW)
				cum := 0
				chosen := fishMults[0]
				for j, w := range fishWeights {
					cum += w
					if r < cum {
						chosen = fishMults[j]
						break
					}
				}
				waveMult += chosen
			}
			totalMult += waveMult
			waveCount++

			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyAtlantisFrenzy,
				Payload: map[string]interface{}{
					"event":      "cascade_wave",
					"wave_no":    wave,
					"fish_count": fishCount,
					"wave_mult":  waveMult,
				},
			})
		}

		reward := int(totalMult * betCost)
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		// 連鎖消除 ≥5 波 → 全服 ×52.0
		globalBonus := 52.0
		globalDuration := 104

		m.mu.Lock()
		m.boost = &atlantisFrenzyBoost{
			mult:      globalBonus,
			expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyAtlantisFrenzy,
			Payload: map[string]interface{}{
				"event":          "atlantis_complete",
				"wave_count":     waveCount,
				"total_mult":     totalMult,
				"reward":         reward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
			},
		})
		g.sendAnnounce(fmt.Sprintf("ATLANTIS FRENZY COMPLETE! %s: %d waves! Total x%.1f! Global x%.1f for %ds!", p.GetDisplayName(), waveCount, totalMult, globalBonus, globalDuration), "critical", "#1E90FF")
		log.Printf("[LuckyAtlantisFrenzy] %s: waves=%d, total=%.1f, global=x%.1f", p.GetDisplayName(), waveCount, totalMult, globalBonus)
	}()

	return true
}
