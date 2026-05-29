// lucky_genesis_ultimate_handler.go — T238 幸運創世終極魚
// 設計：Genesis Ultimate 機制（里程碑：全服 ×50.0）
//       全場清空（每個目標 ×150.0），觸發創世終極光柱（12 道）
//       全服 ×50.0 加成 100 秒（里程碑：史上第一個全服 ×50.0）
//       業界依據：終極清場機制 + 創世機制融合終極版（2026）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type genesisUltimateBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyGenesisUltimateManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *genesisUltimateBoost
}

func newLuckyGenesisUltimateManager() *luckyGenesisUltimateManager {
	return &luckyGenesisUltimateManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyGenesisUltimateFish(defID string) bool {
	return defID == "T238"
}

func (m *luckyGenesisUltimateManager) getGenesisUltimateMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyGenesisUltimateManager) tryLuckyGenesisUltimateFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(520 * time.Second)
	m.personalCD[p.ID] = now.Add(460 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyGenesisUltimate,
		Payload: map[string]interface{}{
			"event":         "genesis_start",
			"trigger_id":    p.ID,
			"trigger_name":  p.GetDisplayName(),
			"per_mult":      150.0,
			"pillar_count":  12,
			"global_target": 50.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("GENESIS ULTIMATE! %s triggered the ULTIMATE CREATION! 12 pillars! Every target x150.0! MILESTONE: Global x50.0!", p.GetDisplayName()), "critical", "#FFD700")
	log.Printf("[LuckyGenesisUltimate] %s triggered Genesis Ultimate fish - MILESTONE x50.0", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		time.Sleep(2 * time.Second)

		// 全場清空
		g.mu.Lock()
		targetCount := len(g.targets)
		for _, t := range g.targets {
			t.HP = 0
		}
		g.mu.Unlock()

		if targetCount == 0 {
			targetCount = 15
		}

		perTargetMult := 150.0
		totalMult := float64(targetCount) * perTargetMult
		reward := int(totalMult * betCost)
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		// 12 道創世光柱動畫
		for pillar := 1; pillar <= 12; pillar++ {
			time.Sleep(250 * time.Millisecond)
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyGenesisUltimate,
				Payload: map[string]interface{}{
					"event":     "genesis_pillar",
					"pillar_no": pillar,
				},
			})
		}

		// 里程碑：全服 ×50.0（無條件觸發）
		globalBonus := 50.0
		globalDuration := 100

		m.mu.Lock()
		m.boost = &genesisUltimateBoost{
			mult:      globalBonus,
			expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyGenesisUltimate,
			Payload: map[string]interface{}{
				"event":          "genesis_milestone",
				"target_count":   targetCount,
				"per_mult":       perTargetMult,
				"total_mult":     totalMult,
				"reward":         reward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
				"milestone":      "GLOBAL_50X",
			},
		})
		g.sendAnnounce(fmt.Sprintf("GENESIS ULTIMATE MILESTONE! %s cleared %d targets! Total x%.1f! GLOBAL x%.1f for %ds! HISTORY MADE!", p.GetDisplayName(), targetCount, totalMult, globalBonus, globalDuration), "critical", "#FFD700")
		log.Printf("[LuckyGenesisUltimate] MILESTONE! %s: targets=%d, total_mult=%.1f, global=x%.1f (FIRST EVER x50.0)", p.GetDisplayName(), targetCount, totalMult, globalBonus)
	}()

	return true
}
