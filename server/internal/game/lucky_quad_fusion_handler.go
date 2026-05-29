// lucky_quad_fusion_handler.go — T248 幸運四重終極融合魚
// 設計：Quad Fusion Ultimate 機制（里程碑：全服 ×56.0）
//       四重機制融合：Wild Collector + Lightning Eel + Domino Chain + Immortal Boss
//       Phase 1: Wild 收集（×2→×3→×10）
//       Phase 2: 閃電鰻連鎖（8條 ×90.0）
//       Phase 3: 骨牌連鎖（20個 ×50.0）
//       Phase 4: 不死 BOSS（5次復活 ×100-×300）
//       全服 ×56.0 加成 112 秒（新史上最高，超越 T246 的 ×55.0）
//       業界依據：四機制終極融合（2026 里程碑）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type quadFusionBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyQuadFusionManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *quadFusionBoost
}

func newLuckyQuadFusionManager() *luckyQuadFusionManager {
	return &luckyQuadFusionManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyQuadFusionFish(defID string) bool {
	return defID == "T248"
}

func (m *luckyQuadFusionManager) getQuadFusionMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyQuadFusionManager) tryLuckyQuadFusionFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(620 * time.Second)
	m.personalCD[p.ID] = now.Add(560 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyQuadFusion,
		Payload: map[string]interface{}{
			"event":        "quad_fusion_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"phases":       4,
			"milestone":    "GLOBAL_56X",
		},
	})
	g.sendAnnounce(fmt.Sprintf("QUAD FUSION ULTIMATE! %s activated 4-Phase Fusion! Wild+Eel+Domino+Boss! MILESTONE: GLOBAL x56.0!", p.GetDisplayName()), "critical", "#FF00FF")
	log.Printf("[LuckyQuadFusion] %s triggered Quad Fusion Ultimate fish - MILESTONE x56.0", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		totalReward := 0
		totalMult := 0.0

		// Phase 1: Wild Collector（×2→×3→×10，10 次旋轉）
		time.Sleep(500 * time.Millisecond)
		phase1Mults := []float64{2.0, 2.0, 2.0, 3.0, 3.0, 3.0, 10.0, 10.0, 10.0, 10.0}
		for i, mult := range phase1Mults {
			time.Sleep(200 * time.Millisecond)
			r := int(mult * betCost * 5.0)
			totalReward += r
			totalMult += mult
			g.mu.Lock()
			p.Coins += r
			g.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyQuadFusion,
				Payload: map[string]interface{}{"event": "phase1_spin", "spin_no": i + 1, "mult": mult},
			})
		}
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyQuadFusion,
			Payload: map[string]interface{}{"event": "phase1_complete", "phase_mult": totalMult},
		})

		// Phase 2: Lightning Eel（8條 ×90.0）
		time.Sleep(500 * time.Millisecond)
		for eel := 1; eel <= 8; eel++ {
			time.Sleep(200 * time.Millisecond)
			eelMult := 90.0
			r := int(eelMult * betCost)
			totalReward += r
			totalMult += eelMult
			g.mu.Lock()
			p.Coins += r
			g.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyQuadFusion,
				Payload: map[string]interface{}{"event": "phase2_eel", "eel_no": eel, "mult": eelMult},
			})
		}
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyQuadFusion,
			Payload: map[string]interface{}{"event": "phase2_complete", "phase_mult": totalMult},
		})

		// Phase 3: Domino Chain（20個 ×50.0）
		time.Sleep(500 * time.Millisecond)
		for domino := 1; domino <= 20; domino++ {
			time.Sleep(150 * time.Millisecond)
			dominoMult := 50.0 * (1.0 + float64(domino-1)*0.05)
			r := int(dominoMult * betCost)
			totalReward += r
			totalMult += dominoMult
			g.mu.Lock()
			p.Coins += r
			g.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyQuadFusion,
				Payload: map[string]interface{}{"event": "phase3_domino", "domino_no": domino, "mult": dominoMult},
			})
		}
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyQuadFusion,
			Payload: map[string]interface{}{"event": "phase3_complete", "phase_mult": totalMult},
		})

		// Phase 4: Immortal Boss（5次復活 ×100→×300）
		time.Sleep(500 * time.Millisecond)
		bossMults := []float64{100.0, 150.0, 200.0, 250.0, 300.0}
		for i, mult := range bossMults {
			time.Sleep(400 * time.Millisecond)
			r := int(mult * betCost)
			totalReward += r
			totalMult += mult
			g.mu.Lock()
			p.Coins += r
			g.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyQuadFusion,
				Payload: map[string]interface{}{"event": "phase4_boss", "revive_no": i + 1, "mult": mult},
			})
		}

		// 里程碑：全服 ×56.0（新史上最高）
		globalBonus := 56.0
		globalDuration := 112

		m.mu.Lock()
		m.boost = &quadFusionBoost{
			mult:      globalBonus,
			expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyQuadFusion,
			Payload: map[string]interface{}{
				"event":          "quad_fusion_milestone",
				"total_mult":     totalMult,
				"total_reward":   totalReward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
				"milestone":      "GLOBAL_56X",
			},
		})
		g.sendAnnounce(fmt.Sprintf("QUAD FUSION MILESTONE! %s: 4 phases complete, total x%.1f! GLOBAL x%.1f for %ds! NEW RECORD!", p.GetDisplayName(), totalMult, globalBonus, globalDuration), "critical", "#FF00FF")
		log.Printf("[LuckyQuadFusion] MILESTONE! %s: total_mult=%.1f, global=x%.1f (NEW RECORD x56.0)", p.GetDisplayName(), totalMult, globalBonus)
	}()

	return true
}
