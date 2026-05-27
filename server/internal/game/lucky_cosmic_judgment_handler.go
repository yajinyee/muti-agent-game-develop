// lucky_cosmic_judgment_handler.go — T190 幸運宇宙審判魚
// 業界依據：Fishing Fortune「ultimate judgment」+ 宇宙終極機制
// 設計：擊破後宇宙審判（全場 HP 歸零，每個獎勵 ×14.0）
//       觸發全服 ×19.0 加成 40 秒（新最高全服倍率機制，超越 T185 的 ×16.0）
//       個人冷卻 115 秒；全服冷卻 175 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyCosmicJudgmentManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *cosmicJudgmentPerfectBoost
}

type cosmicJudgmentPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyCosmicJudgmentManager() *luckyCosmicJudgmentManager {
	return &luckyCosmicJudgmentManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyCosmicJudgmentFish(defID string) bool {
	return defID == "T190"
}

func (m *luckyCosmicJudgmentManager) getCosmicJudgmentMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyCosmicJudgmentManager) tryLuckyCosmicJudgmentFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(115 * time.Second)
	m.globalCD = now.Add(175 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_cosmic_judgment",
		Payload: map[string]interface{}{
			"event":        "cosmic_judgment_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("⚖️🌌 宇宙審判！%s 召喚宇宙審判！全場 HP 歸零！每個獎勵 ×14.0！全服 ×19.0！", p.GetDisplayName()), "critical", "#B71C1C")
	log.Printf("[LuckyCosmicJudgment] %s 觸發宇宙審判魚", p.GetDisplayName())

	go func() {
		time.Sleep(700 * time.Millisecond)

		// 宇宙審判：全場 HP 歸零，每個獎勵 ×14.0（超越 T185 的 ×12.0）
		hitCount := g.applyUltimateJudgment(p, 14.0)

		// 觸發全服 ×19.0 加成 40 秒（新最高全服倍率機制）
		boostMult := 19.0
		boostSecs := 40
		m.mu.Lock()
		m.perfectBoost = &cosmicJudgmentPerfectBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_cosmic_judgment",
			Payload: map[string]interface{}{
				"event":        "cosmic_judgment_complete",
				"hit_count":    hitCount,
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
			},
		})
		g.sendAnnounce(fmt.Sprintf("⚖️🏆 宇宙審判完成！%s 清場 %d 個！全服 ×%.1f 加成 %d 秒！（新最高）", p.GetDisplayName(), hitCount, boostMult, boostSecs), "critical", "#7F0000")
	}()
	return true
}
