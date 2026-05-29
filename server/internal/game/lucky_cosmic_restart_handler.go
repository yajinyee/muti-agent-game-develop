// lucky_cosmic_restart_handler.go — T233 幸運宇宙重啟魚
// 設計：Cosmic Restart 機制
//       全場 HP 歸零（每個目標獎勵 ×100.0），觸發全服 ×47.5 加成 95 秒（新史上最高）
//       業界依據：終極清場機制升級版 + Cosmic Restart 概念（2026）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type cosmicRestartBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyCosmicRestartManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *cosmicRestartBoost
}

func newLuckyCosmicRestartManager() *luckyCosmicRestartManager {
	return &luckyCosmicRestartManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyCosmicRestartFish(defID string) bool {
	return defID == "T233"
}

func (m *luckyCosmicRestartManager) getCosmicRestartMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyCosmicRestartManager) tryLuckyCosmicRestartFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(460 * time.Second)
	m.personalCD[p.ID] = now.Add(400 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyCosmicRestart,
		Payload: map[string]interface{}{
			"event":        "restart_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"per_mult":     100.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("COSMIC RESTART! %s restarted the universe! Every target x100.0!", p.GetDisplayName()), "critical", "#FF00FF")
	log.Printf("[LuckyCosmicRestart] %s triggered Cosmic Restart fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		time.Sleep(2 * time.Second)

		g.mu.Lock()
		targetCount := len(g.targets)
		for _, t := range g.targets {
			t.HP = 0
		}
		g.mu.Unlock()

		if targetCount == 0 {
			targetCount = 10
		}

		perTargetMult := 100.0
		totalMult := float64(targetCount) * perTargetMult
		reward := int(totalMult * betCost)
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		globalBonus := 47.5
		globalDuration := 95

		m.mu.Lock()
		m.boost = &cosmicRestartBoost{
			mult:      globalBonus,
			expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyCosmicRestart,
			Payload: map[string]interface{}{
				"event":          "restart_result",
				"target_count":   targetCount,
				"per_mult":       perTargetMult,
				"total_mult":     totalMult,
				"reward":         reward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
			},
		})

		g.sendAnnounce(fmt.Sprintf("UNIVERSE RESTARTED! %s cleared %d targets! Total x%.1f! Global x%.1f for %ds! NEW ALL-TIME HIGH!", p.GetDisplayName(), targetCount, totalMult, globalBonus, globalDuration), "critical", "#FF00FF")
		log.Printf("[LuckyCosmicRestart] %s: targets=%d, total_mult=%.1f, global=x%.1f", p.GetDisplayName(), targetCount, totalMult, globalBonus)
	}()

	return true
}
