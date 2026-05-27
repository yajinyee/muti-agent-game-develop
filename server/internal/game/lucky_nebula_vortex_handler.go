// lucky_nebula_vortex_handler.go — T189 幸運星雲漩渦魚
// 業界依據：Fishing Carnival「vortex anemone」+ 星雲能量吸收概念
// 設計：擊破後星雲漩渦 20 秒（每秒全場 HP -8%，每個獎勵 ×1.5）
//       20 秒內累積傷害 ≥160% → 星雲完美：全服 ×18.5 加成 39 秒
//       個人冷卻 112 秒；全服冷卻 172 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyNebulaVortexManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *nebulaVortexPerfectBoost
}

type nebulaVortexPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyNebulaVortexManager() *luckyNebulaVortexManager {
	return &luckyNebulaVortexManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyNebulaVortexFish(defID string) bool {
	return defID == "T189"
}

func (m *luckyNebulaVortexManager) getNebulaVortexMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyNebulaVortexManager) tryLuckyNebulaVortexFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) {
		m.mu.Unlock()
		return false
	}
	if cd, ok := m.personalCD[p.ID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return false
	}
	m.personalCD[p.ID] = now.Add(112 * time.Second)
	m.globalCD = now.Add(172 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_nebula_vortex",
		Payload: map[string]interface{}{
			"event":        "nebula_vortex_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     20,
			"damage_per_sec": 0.08,
			"reward_mult":  1.5,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🌌🌀 星雲漩渦！%s 召喚星雲漩渦！每秒全場 HP -8%%！持續 20 秒！", p.GetDisplayName()), "critical", "#4A148C")
	log.Printf("[LuckyNebulaVortex] %s 觸發星雲漩渦魚", p.GetDisplayName())

	go func() {
		totalDamagePct := 0.0
		totalHits := 0

		// 每秒施加 8% 傷害，持續 20 秒
		for wave := 1; wave <= 20; wave++ {
			time.Sleep(1 * time.Second)

			hitCount := g.applyAOEDamage(0, 0, 99999, 0.08)
			totalDamagePct += 8.0
			totalHits += hitCount

			g.broadcast(protocol.Envelope{
				Type: "lucky_nebula_vortex",
				Payload: map[string]interface{}{
					"event":            "nebula_wave",
					"wave":             wave,
					"hit_count":        hitCount,
					"total_damage_pct": totalDamagePct,
					"trigger_id":       p.ID,
					"trigger_name":     p.GetDisplayName(),
				},
			})
		}

		// 判定完美：累積傷害 ≥160%（即 20 波全部命中 ≥1 個目標）
		if totalDamagePct >= 160.0 && totalHits >= 20 {
			boostMult := 18.5
			boostSecs := 39
			m.mu.Lock()
			m.perfectBoost = &nebulaVortexPerfectBoost{
				mult:      boostMult,
				expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
			}
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_nebula_vortex",
				Payload: map[string]interface{}{
					"event":        "nebula_perfect",
					"total_hits":   totalHits,
					"boost_mult":   boostMult,
					"boost_secs":   boostSecs,
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
				},
			})
			g.sendAnnounce(fmt.Sprintf("🌌🏆 星雲完美！%s 累積命中 %d 次！全服 ×%.1f 加成 %d 秒！", p.GetDisplayName(), totalHits, boostMult, boostSecs), "critical", "#6A1B9A")
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_nebula_vortex",
				Payload: map[string]interface{}{
					"event":        "nebula_vortex_end",
					"total_hits":   totalHits,
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
				},
			})
			g.sendAnnounce(fmt.Sprintf("🌌 星雲漩渦結束！%s 累積命中 %d 次", p.GetDisplayName(), totalHits), "normal", "#7B1FA2")
		}
	}()
	return true
}
