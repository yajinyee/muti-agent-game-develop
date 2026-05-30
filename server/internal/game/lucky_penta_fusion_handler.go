// lucky_penta_fusion_handler.go — T253 幸運五重終極魚
// 設計：Penta Fusion Ultimate 機制（里程碑：全服 ×58.5）
//       五重機制融合：電擊框架 + 磁力連鎖 + 漁夫路徑 + 黃金鰓 Jackpot + Quad Fusion
//       Phase 1: 電擊框架（10次命中，×1→×1024）
//       Phase 2: 磁力連鎖 Respin（8次，×75.0）
//       Phase 3: 漁夫路徑（10節點，最高 ×500）
//       Phase 4: 黃金鰓 Jackpot（Grand ×2000）
//       Phase 5: Quad Fusion 終極（全場清空 ×200.0）
//       全服 ×58.5 加成 117 秒（新史上最高，超越 T248 的 ×56.0）
//       業界依據：五機制終極融合（2026 里程碑）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type pentaFusionBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyPentaFusionManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *pentaFusionBoost
}

func newLuckyPentaFusionManager() *luckyPentaFusionManager {
	return &luckyPentaFusionManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyPentaFusionFish(defID string) bool {
	return defID == "T253"
}

func (m *luckyPentaFusionManager) getPentaFusionMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyPentaFusionManager) tryLuckyPentaFusionFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(680 * time.Second)
	m.personalCD[p.ID] = now.Add(620 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyPentaFusion,
		Payload: map[string]interface{}{
			"event":        "penta_fusion_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"phases":       5,
			"milestone":    "GLOBAL_58.5X",
		},
	})
	g.sendAnnounce(fmt.Sprintf("PENTA FUSION ULTIMATE! %s activated 5-Phase Fusion! Frame+Magnetic+Trail+Gills+Quad! MILESTONE: GLOBAL x58.5!", p.GetDisplayName()), "critical", "#FF69B4")
	log.Printf("[LuckyPentaFusion] %s triggered Penta Fusion Ultimate fish - MILESTONE x58.5", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		totalReward := 0
		totalMult := 0.0

		// Phase 1: 電擊框架（10次命中，×1→×1024）
		time.Sleep(500 * time.Millisecond)
		globalMult := 1.0
		for i := 0; i < 10; i++ {
			time.Sleep(200 * time.Millisecond)
			globalMult *= 2.0
			if globalMult > 1024.0 {
				globalMult = 1024.0
			}
			r := int(globalMult * betCost * 0.5)
			totalReward += r
			totalMult += globalMult
			g.mu.Lock()
			p.Coins += r
			g.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyPentaFusion,
				Payload: map[string]interface{}{"event": "phase1_frame", "hit_no": i + 1, "global_mult": globalMult},
			})
		}
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyPentaFusion,
			Payload: map[string]interface{}{"event": "phase1_complete", "final_mult": globalMult},
		})

		// Phase 2: 磁力連鎖 Respin（8次，×75.0）
		time.Sleep(500 * time.Millisecond)
		for respin := 1; respin <= 8; respin++ {
			time.Sleep(200 * time.Millisecond)
			r := int(75.0 * betCost)
			totalReward += r
			totalMult += 75.0
			g.mu.Lock()
			p.Coins += r
			g.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyPentaFusion,
				Payload: map[string]interface{}{"event": "phase2_respin", "respin_no": respin, "mult": 75.0},
			})
		}
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyPentaFusion,
			Payload: map[string]interface{}{"event": "phase2_complete"},
		})

		// Phase 3: 漁夫路徑（10節點，最高 ×500）
		time.Sleep(500 * time.Millisecond)
		nodeMults := []float64{10.0, 20.0, 30.0, 50.0, 75.0, 100.0, 150.0, 200.0, 350.0, 500.0}
		for i, mult := range nodeMults {
			time.Sleep(200 * time.Millisecond)
			r := int(mult * betCost)
			totalReward += r
			totalMult += mult
			g.mu.Lock()
			p.Coins += r
			g.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyPentaFusion,
				Payload: map[string]interface{}{"event": "phase3_node", "node_no": i + 1, "mult": mult},
			})
		}
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyPentaFusion,
			Payload: map[string]interface{}{"event": "phase3_complete"},
		})

		// Phase 4: 黃金鰓 Jackpot（Grand ×2000）
		time.Sleep(500 * time.Millisecond)
		grandReward := int(2000.0 * betCost)
		totalReward += grandReward
		totalMult += 2000.0
		g.mu.Lock()
		p.Coins += grandReward
		g.mu.Unlock()
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyPentaFusion,
			Payload: map[string]interface{}{"event": "phase4_grand_jackpot", "mult": 2000.0, "reward": grandReward},
		})
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyPentaFusion,
			Payload: map[string]interface{}{"event": "phase4_complete"},
		})

		// Phase 5: Quad Fusion 終極（全場清空 ×200.0）
		time.Sleep(500 * time.Millisecond)
		clearMults := []float64{200.0, 200.0, 200.0, 200.0, 200.0}
		for i, mult := range clearMults {
			time.Sleep(300 * time.Millisecond)
			r := int(mult * betCost)
			totalReward += r
			totalMult += mult
			g.mu.Lock()
			p.Coins += r
			g.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyPentaFusion,
				Payload: map[string]interface{}{"event": "phase5_clear", "clear_no": i + 1, "mult": mult},
			})
		}

		// 里程碑：全服 ×58.5（新史上最高）
		globalBonus := 58.5
		globalDuration := 117

		m.mu.Lock()
		m.boost = &pentaFusionBoost{
			mult:      globalBonus,
			expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyPentaFusion,
			Payload: map[string]interface{}{
				"event":          "penta_fusion_milestone",
				"total_mult":     totalMult,
				"total_reward":   totalReward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
				"milestone":      "GLOBAL_58.5X",
			},
		})
		g.sendAnnounce(fmt.Sprintf("PENTA FUSION MILESTONE! %s: 5 phases complete, total x%.1f! GLOBAL x%.1f for %ds! NEW RECORD!", p.GetDisplayName(), totalMult, globalBonus, globalDuration), "critical", "#FF69B4")
		log.Printf("[LuckyPentaFusion] MILESTONE! %s: total_mult=%.1f, global=x%.1f (NEW RECORD x58.5)", p.GetDisplayName(), totalMult, globalBonus)
	}()

	return true
}
