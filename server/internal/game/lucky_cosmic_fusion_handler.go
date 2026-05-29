// lucky_cosmic_fusion_handler.go — T228 幸運宇宙大融合魚
// 設計：終極融合機制 — 融合所有已知 Lucky 機制的精華
//       Phase 1：Coin Respin（3格，Bronze/Silver/Gold/Diamond）
//       Phase 2：Cascade Lock（4 波，每波 ×5-×20）
//       Phase 3：Legend Awaken（3 次連續獎勵）
//       Phase 4：全場 HP 歸零（每個獎勵 ×20.0）
//       全部完成 → 全服 ×45.0 加成 90 秒（新史上最高，超越 T227 的 ×44.5）
//       觸發率：0.0002%（最稀有）；個人冷卻 350 秒；全服冷卻 410 秒
//       業界依據：終極融合設計，整合 2026 年最新業界機制
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyCosmicFusionManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
}

func newLuckyCosmicFusionManager() *luckyCosmicFusionManager {
	return &luckyCosmicFusionManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyCosmicFusionFish(defID string) bool {
	return defID == "T228"
}

func (m *luckyCosmicFusionManager) tryLuckyCosmicFusionFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(410 * time.Second)
	m.personalCD[p.ID] = now.Add(350 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_cosmic_fusion",
		Payload: map[string]interface{}{
			"event":        "fusion_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"phases":       4,
		},
	})
	g.sendAnnounce(fmt.Sprintf("COSMIC FUSION! %s triggered the ultimate fusion! 4 phases of destruction!", p.GetDisplayName()), "critical", "#FF00FF")
	log.Printf("[LuckyCosmicFusion] %s triggered Cosmic Fusion fish", p.GetDisplayName())

	go func() {
		totalMult := 0.0
		betCost := float64(p.GetBetDef().BetCost)

		// ── Phase 1: Mini Coin Respin (3 slots) ──────────────────
		g.broadcast(protocol.Envelope{
			Type: "lucky_cosmic_fusion",
			Payload: map[string]interface{}{
				"event": "phase_start",
				"phase": 1,
				"name":  "Coin Respin",
			},
		})

		phase1Mult := 0.0
		grid := make([]float64, 3)
		spinsLeft := 3
		for spinsLeft > 0 {
			if rand.Float64() < 0.70 {
				emptySlots := []int{}
				for i, v := range grid {
					if v == 0 {
						emptySlots = append(emptySlots, i)
					}
				}
				if len(emptySlots) > 0 {
					slot := emptySlots[rand.Intn(len(emptySlots))]
					coin := rollGoldenPotCoin()
					grid[slot] = coin.Mult
					phase1Mult += coin.Mult
					spinsLeft = 3
					g.broadcast(protocol.Envelope{
						Type: "lucky_cosmic_fusion",
						Payload: map[string]interface{}{
							"event":      "phase1_coin",
							"slot":       slot,
							"coin_name":  coin.Name,
							"coin_mult":  coin.Mult,
							"phase_mult": phase1Mult,
						},
					})
				}
			}
			spinsLeft--
			time.Sleep(500 * time.Millisecond)
		}
		totalMult += phase1Mult
		g.broadcast(protocol.Envelope{
			Type: "lucky_cosmic_fusion",
			Payload: map[string]interface{}{
				"event":      "phase_complete",
				"phase":      1,
				"phase_mult": phase1Mult,
				"total_mult": totalMult,
			},
		})
		time.Sleep(800 * time.Millisecond)

		// ── Phase 2: Cascade Lock (4 waves) ──────────────────────
		g.broadcast(protocol.Envelope{
			Type: "lucky_cosmic_fusion",
			Payload: map[string]interface{}{
				"event": "phase_start",
				"phase": 2,
				"name":  "Cascade Lock",
			},
		})

		phase2Mult := 0.0
		for wave := 1; wave <= 4; wave++ {
			waveMult := 5.0 + rand.Float64()*15.0 // 5-20x per wave
			if rand.Float64() < 0.25 {             // 25% Pearl bonus
				waveMult *= float64(2 + rand.Intn(4)) // ×2-×5
			}
			phase2Mult += waveMult
			g.broadcast(protocol.Envelope{
				Type: "lucky_cosmic_fusion",
				Payload: map[string]interface{}{
					"event":      "phase2_wave",
					"wave":       wave,
					"wave_mult":  waveMult,
					"phase_mult": phase2Mult,
				},
			})
			time.Sleep(600 * time.Millisecond)
		}
		totalMult += phase2Mult
		g.broadcast(protocol.Envelope{
			Type: "lucky_cosmic_fusion",
			Payload: map[string]interface{}{
				"event":      "phase_complete",
				"phase":      2,
				"phase_mult": phase2Mult,
				"total_mult": totalMult,
			},
		})
		time.Sleep(800 * time.Millisecond)

		// ── Phase 3: Legend Awaken (3 rounds) ────────────────────
		g.broadcast(protocol.Envelope{
			Type: "lucky_cosmic_fusion",
			Payload: map[string]interface{}{
				"event": "phase_start",
				"phase": 3,
				"name":  "Legend Awaken",
			},
		})

		phase3Mult := 0.0
		for round := 1; round <= 3; round++ {
			roundMult := 10.0 + float64(round)*5.0 + rand.Float64()*20.0 // 15-45x per round
			phase3Mult += roundMult
			roundReward := int(roundMult * betCost)
			g.mu.Lock()
			p.Coins += roundReward
			g.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_cosmic_fusion",
				Payload: map[string]interface{}{
					"event":        "phase3_awaken",
					"round":        round,
					"round_mult":   roundMult,
					"round_reward": roundReward,
					"phase_mult":   phase3Mult,
				},
			})
			time.Sleep(700 * time.Millisecond)
		}
		totalMult += phase3Mult
		g.broadcast(protocol.Envelope{
			Type: "lucky_cosmic_fusion",
			Payload: map[string]interface{}{
				"event":      "phase_complete",
				"phase":      3,
				"phase_mult": phase3Mult,
				"total_mult": totalMult,
			},
		})
		time.Sleep(800 * time.Millisecond)

		// ── Phase 4: Full Field Clear ─────────────────────────────
		g.broadcast(protocol.Envelope{
			Type: "lucky_cosmic_fusion",
			Payload: map[string]interface{}{
				"event": "phase_start",
				"phase": 4,
				"name":  "Cosmic Clear",
			},
		})

		g.mu.Lock()
		clearedCount := 0
		phase4Mult := 0.0
		for _, t := range g.targets {
			if t.HP > 0 {
				t.HP = 0
				clearedCount++
				phase4Mult += 20.0
			}
		}
		g.mu.Unlock()

		if clearedCount == 0 {
			phase4Mult = 20.0 // minimum
		}
		totalMult += phase4Mult

		g.broadcast(protocol.Envelope{
			Type: "lucky_cosmic_fusion",
			Payload: map[string]interface{}{
				"event":         "phase4_clear",
				"cleared_count": clearedCount,
				"phase_mult":    phase4Mult,
				"total_mult":    totalMult,
			},
		})
		g.sendAnnounce("Cosmic Clear! All targets destroyed!", "critical", "#FF00FF")
		time.Sleep(1000 * time.Millisecond)

		// Final reward
		reward := int(totalMult * betCost)
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_cosmic_fusion",
			Payload: map[string]interface{}{
				"event":        "fusion_settle",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"total_mult":   totalMult,
				"reward":       reward,
				"phase1_mult":  phase1Mult,
				"phase2_mult":  phase2Mult,
				"phase3_mult":  phase3Mult,
				"phase4_mult":  phase4Mult,
			},
		})

		// Global boost x45.0 for 90 seconds (new all-time high)
		g.broadcast(protocol.Envelope{
			Type: "lucky_cosmic_fusion",
			Payload: map[string]interface{}{
				"event":        "global_boost",
				"global_mult":  45.0,
				"duration":     90,
				"trigger_name": p.GetDisplayName(),
				"total_mult":   totalMult,
			},
		})
		g.sendAnnounce(fmt.Sprintf("COSMIC FUSION COMPLETE! %s fused x%.0f! Global x45.0 for 90s! NEW ALL-TIME HIGH!", p.GetDisplayName(), totalMult), "critical", "#FF00FF")
		log.Printf("[LuckyCosmicFusion] %s Cosmic Fusion settled %.0fx, global x45.0 for 90s", p.GetDisplayName(), totalMult)
	}()

	return true
}
