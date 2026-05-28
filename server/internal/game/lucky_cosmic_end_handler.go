// lucky_cosmic_end_handler.go — T195 幸運宇宙終焉魚
// 設計：宇宙終焉，全場 HP 歸零，每個獎勵 ×20.0
//       觸發後全服 ×22.0 加成 45 秒（新最高，超越 T190 的 ×19.0）
//       觸發率：0.02%（最稀有）；個人冷卻 140 秒；全服冷卻 200 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyCosmicEndManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	endBoost   *cosmicEndPerfectBoost
}

type cosmicEndPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyCosmicEndManager() *luckyCosmicEndManager {
	return &luckyCosmicEndManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyCosmicEndFish(defID string) bool {
	return defID == "T195"
}

func (m *luckyCosmicEndManager) getCosmicEndMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.endBoost != nil && time.Now().Before(m.endBoost.expiresAt) {
		return m.endBoost.mult
	}
	return 1.0
}

func (m *luckyCosmicEndManager) tryLuckyCosmicEndFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(140 * time.Second)
	m.globalCD = now.Add(200 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_cosmic_end",
		Payload: map[string]interface{}{
			"event":        "cosmic_end_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("☄️💀 宇宙終焉！%s 召喚宇宙終焉！全場 HP 歸零！每個獎勵 ×20.0！全服 ×22.0！", p.GetDisplayName()), "critical", "#000000")
	log.Printf("[LuckyCosmicEnd] %s 觸發宇宙終焉魚（最稀有）", p.GetDisplayName())

	go func() {
		time.Sleep(700 * time.Millisecond)

		// 宇宙終焉：全場 HP 歸零，每個獎勵 ×20.0（超越 T190 的 ×14.0）
		hitCount := g.applyUltimateJudgment(p, 20.0)

		// 觸發全服 ×22.0 加成 45 秒（新最高，超越 T190 的 ×19.0）
		boostMult := 22.0
		boostSecs := 45
		m.mu.Lock()
		m.endBoost = &cosmicEndPerfectBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_cosmic_end",
			Payload: map[string]interface{}{
				"event":        "cosmic_end_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"hit_count":    hitCount,
				"reward_mult":  20.0,
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
			},
		})
		g.sendAnnounce(fmt.Sprintf("☄️👑 宇宙終焉完成！%s 清場 %d 個！全服 ×%.1f 加成 %d 秒！（史上最高）",
			p.GetDisplayName(), hitCount, boostMult, boostSecs), "critical", "#1A0000")
	}()
	return true
}
