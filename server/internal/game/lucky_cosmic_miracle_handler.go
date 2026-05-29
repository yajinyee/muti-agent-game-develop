// lucky_cosmic_miracle_handler.go — T237 幸運宇宙奇蹟魚
// 設計：Cosmic Miracle 機制
//       全場 HP 歸零（每個目標獎勵 ×120.0），觸發宇宙奇蹟光柱（8 道）
//       命中 ≥12 個目標 → 完美奇蹟，全服 ×49.5 加成 99 秒（超越 T236 的 ×49.0）
//       業界依據：終極清場機制 + 神聖光柱機制融合升級版（2026）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type cosmicMiracleBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyCosmicMiracleManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *cosmicMiracleBoost
}

func newLuckyCosmicMiracleManager() *luckyCosmicMiracleManager {
	return &luckyCosmicMiracleManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyCosmicMiracleFish(defID string) bool {
	return defID == "T237"
}

func (m *luckyCosmicMiracleManager) getCosmicMiracleMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyCosmicMiracleManager) tryLuckyCosmicMiracleFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(510 * time.Second)
	m.personalCD[p.ID] = now.Add(450 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyCosmicMiracle,
		Payload: map[string]interface{}{
			"event":        "miracle_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"per_mult":     120.0,
			"pillar_count": 8,
		},
	})
	g.sendAnnounce(fmt.Sprintf("COSMIC MIRACLE! %s summoned 8 cosmic pillars! Every target x120.0!", p.GetDisplayName()), "critical", "#9400D3")
	log.Printf("[LuckyCosmicMiracle] %s triggered Cosmic Miracle fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		time.Sleep(2 * time.Second)

		// 全場 HP 歸零
		g.mu.Lock()
		targetCount := len(g.targets)
		for _, t := range g.targets {
			t.HP = 0
		}
		g.mu.Unlock()

		if targetCount == 0 {
			targetCount = 12
		}

		perTargetMult := 120.0
		totalMult := float64(targetCount) * perTargetMult
		reward := int(totalMult * betCost)
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		// 8 道宇宙光柱動畫
		for pillar := 1; pillar <= 8; pillar++ {
			time.Sleep(300 * time.Millisecond)
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyCosmicMiracle,
				Payload: map[string]interface{}{
					"event":      "miracle_pillar",
					"pillar_no":  pillar,
					"hit_count":  targetCount / 8,
				},
			})
		}

		isPerfect := targetCount >= 12
		globalBonus := 49.5
		globalDuration := 99

		if isPerfect {
			m.mu.Lock()
			m.boost = &cosmicMiracleBoost{
				mult:      globalBonus,
				expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
			}
			m.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyCosmicMiracle,
				Payload: map[string]interface{}{
					"event":          "miracle_perfect",
					"target_count":   targetCount,
					"per_mult":       perTargetMult,
					"total_mult":     totalMult,
					"reward":         reward,
					"global_bonus":   globalBonus,
					"global_seconds": globalDuration,
				},
			})
			g.sendAnnounce(fmt.Sprintf("PERFECT COSMIC MIRACLE! %s cleared %d targets! Total x%.1f! Global x%.1f for %ds!", p.GetDisplayName(), targetCount, totalMult, globalBonus, globalDuration), "critical", "#9400D3")
		} else {
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyCosmicMiracle,
				Payload: map[string]interface{}{
					"event":        "miracle_end",
					"target_count": targetCount,
					"per_mult":     perTargetMult,
					"total_mult":   totalMult,
					"reward":       reward,
				},
			})
		}

		log.Printf("[LuckyCosmicMiracle] %s: targets=%d, total_mult=%.1f, perfect=%v", p.GetDisplayName(), targetCount, totalMult, isPerfect)
	}()

	return true
}
