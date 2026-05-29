// lucky_holy_pillar_handler.go — T231 幸運神聖光柱魚
// 設計：Holy Pillar 機制
//       12 道神聖光柱同時降下（每道 HP -50%），命中 ≥8 道 → 完美神聖
//       完美神聖 → 全服 ×46.5 加成 93 秒（超越 T230 的 ×46.0）
//       業界依據：神聖審判機制升級版 + Holy Pillar 概念（2026）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type holyPillarBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyHolyPillarManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *holyPillarBoost
}

func newLuckyHolyPillarManager() *luckyHolyPillarManager {
	return &luckyHolyPillarManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyHolyPillarFish(defID string) bool {
	return defID == "T231"
}

func (m *luckyHolyPillarManager) getHolyPillarMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyHolyPillarManager) tryLuckyHolyPillarFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(440 * time.Second)
	m.personalCD[p.ID] = now.Add(380 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyHolyPillar,
		Payload: map[string]interface{}{
			"event":        "pillar_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"pillar_count": 12,
		},
	})
	g.sendAnnounce(fmt.Sprintf("HOLY PILLAR! %s summoned 12 holy pillars! Divine judgment descends!", p.GetDisplayName()), "critical", "#FFFF00")
	log.Printf("[LuckyHolyPillar] %s triggered Holy Pillar fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		totalMult := 0.0
		hitPillars := 0

		for i := 1; i <= 12; i++ {
			time.Sleep(600 * time.Millisecond)
			if rand.Float64() < 0.85 {
				pillarMult := 15.0 + rand.Float64()*10.0
				totalMult += pillarMult
				hitPillars++
				reward := int(pillarMult * betCost)
				g.mu.Lock()
				p.Coins += reward
				for _, t := range g.targets {
					damage := int(float64(t.MaxHP) * 0.50)
					t.HP -= damage
					if t.HP <= 0 {
						t.HP = 0
					}
				}
				g.mu.Unlock()
				g.broadcast(protocol.Envelope{
					Type: protocol.MsgLuckyHolyPillar,
					Payload: map[string]interface{}{
						"event":        "pillar_hit",
						"pillar_index": i,
						"pillar_mult":  pillarMult,
						"total_mult":   totalMult,
						"reward":       reward,
					},
				})
			} else {
				g.broadcast(protocol.Envelope{
					Type: protocol.MsgLuckyHolyPillar,
					Payload: map[string]interface{}{
						"event":        "pillar_miss",
						"pillar_index": i,
					},
				})
			}
		}

		time.Sleep(1 * time.Second)
		isPerfect := hitPillars >= 8
		globalBonus := 46.5
		globalDuration := 93

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyHolyPillar,
			Payload: map[string]interface{}{
				"event":          "pillar_result",
				"hit_pillars":    hitPillars,
				"total_mult":     totalMult,
				"is_perfect":     isPerfect,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
			},
		})

		if isPerfect {
			m.mu.Lock()
			m.boost = &holyPillarBoost{
				mult:      globalBonus,
				expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
			}
			m.mu.Unlock()
			g.sendAnnounce(fmt.Sprintf("PERFECT HOLY PILLAR! %s hit %d/12 pillars! Total x%.1f! Global x%.1f for %ds!", p.GetDisplayName(), hitPillars, totalMult, globalBonus, globalDuration), "critical", "#FFFF00")
		}

		log.Printf("[LuckyHolyPillar] %s: hit=%d/12, total_mult=%.1f, perfect=%v", p.GetDisplayName(), hitPillars, totalMult, isPerfect)
	}()

	return true
}
