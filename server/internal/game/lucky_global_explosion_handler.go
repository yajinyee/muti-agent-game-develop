// lucky_global_explosion_handler.go — T193 幸運全服大爆炸魚
// 設計：全服大爆炸，全場 HP 歸零，每個獎勵 ×15.0
//       觸發後全服 ×20.5 加成 40 秒（新最高）
//       觸發率：0.03%；個人冷卻 130 秒；全服冷卻 190 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyGlobalExplosionManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	explBoost  *globalExplosionPerfectBoost
}

type globalExplosionPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyGlobalExplosionManager() *luckyGlobalExplosionManager {
	return &luckyGlobalExplosionManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyGlobalExplosionFish(defID string) bool {
	return defID == "T193"
}

func (m *luckyGlobalExplosionManager) getGlobalExplosionMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.explBoost != nil && time.Now().Before(m.explBoost.expiresAt) {
		return m.explBoost.mult
	}
	return 1.0
}

func (m *luckyGlobalExplosionManager) tryLuckyGlobalExplosionFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(130 * time.Second)
	m.globalCD = now.Add(190 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_global_explosion",
		Payload: map[string]interface{}{
			"event":        "global_explosion_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("💥🌍 全服大爆炸！%s 引爆全服大爆炸！全場 HP 歸零！每個獎勵 ×15.0！", p.GetDisplayName()), "critical", "#B71C1C")
	log.Printf("[LuckyGlobalExplosion] %s 觸發全服大爆炸魚", p.GetDisplayName())

	go func() {
		time.Sleep(600 * time.Millisecond)

		// 全服大爆炸：全場 HP 歸零，每個獎勵 ×15.0
		hitCount := g.applyUltimateJudgment(p, 15.0)

		// 觸發全服 ×20.5 加成 40 秒（新最高）
		boostMult := 20.5
		boostSecs := 40
		m.mu.Lock()
		m.explBoost = &globalExplosionPerfectBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_global_explosion",
			Payload: map[string]interface{}{
				"event":        "global_explosion_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"hit_count":    hitCount,
				"reward_mult":  15.0,
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
			},
		})
		g.sendAnnounce(fmt.Sprintf("💥🏆 全服大爆炸完成！%s 清場 %d 個！全服 ×%.1f 加成 %d 秒！（新最高）",
			p.GetDisplayName(), hitCount, boostMult, boostSecs), "critical", "#7F0000")
	}()
	return true
}
