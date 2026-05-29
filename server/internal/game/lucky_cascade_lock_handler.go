// lucky_cascade_lock_handler.go — T225 幸運瀑布鎖定魚
// 設計：Cascading Wins + 鎖定倍率機制（BGaming「Shark & Spark Hold & Win」2026-05-28）
//       每次擊破觸發瀑布連鎖（最多 8 波），每波倍率鎖定並累積
//       Pearl 符號隨機出現（×2-×10 倍率加成），鎖定後不消失
//       8 波全部完成 → 全服 ×43.5 加成 87 秒（超越 T224 的 ×43.0）
//       觸發率：0.00035%；個人冷卻 335 秒；全服冷卻 395 秒
//       業界依據：BGaming「Shark & Spark Hold & Win」Cascading Wins + Pearl Multipliers（2026-05-28）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyCascadeLockManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
}

func newLuckyCascadeLockManager() *luckyCascadeLockManager {
	return &luckyCascadeLockManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyCascadeLockFish(defID string) bool {
	return defID == "T225"
}

func (m *luckyCascadeLockManager) tryLuckyCascadeLockFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(395 * time.Second)
	m.personalCD[p.ID] = now.Add(335 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_cascade_lock",
		Payload: map[string]interface{}{
			"event":        "cascade_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"max_waves":    8,
		},
	})
	g.sendAnnounce(fmt.Sprintf("Cascade Lock! %s triggered Cascading Wins! 8 waves incoming!", p.GetDisplayName()), "critical", "#00BFFF")
	log.Printf("[LuckyCascadeLock] %s triggered Cascade Lock fish", p.GetDisplayName())

	go func() {
		totalMult := 0.0
		lockedMults := []float64{}
		perfectCascade := true

		for wave := 1; wave <= 8; wave++ {
			// Base wave multiplier: 3.0 + wave * 1.5
			waveMult := 3.0 + float64(wave)*1.5

			// 30% chance Pearl symbol appears (×2-×10 bonus)
			pearlMult := 1.0
			hasPearl := rand.Float64() < 0.30
			if hasPearl {
				pearlMult = float64(2 + rand.Intn(9)) // 2-10
				waveMult *= pearlMult
			}

			// 85% chance wave hits (15% miss = not perfect)
			waveHit := rand.Float64() < 0.85
			if !waveHit {
				perfectCascade = false
				g.broadcast(protocol.Envelope{
					Type: "lucky_cascade_lock",
					Payload: map[string]interface{}{
						"event":      "wave_miss",
						"wave":       wave,
						"total_mult": totalMult,
					},
				})
				time.Sleep(400 * time.Millisecond)
				continue
			}

			totalMult += waveMult
			lockedMults = append(lockedMults, waveMult)

			g.broadcast(protocol.Envelope{
				Type: "lucky_cascade_lock",
				Payload: map[string]interface{}{
					"event":       "wave_hit",
					"wave":        wave,
					"wave_mult":   waveMult,
					"has_pearl":   hasPearl,
					"pearl_mult":  pearlMult,
					"total_mult":  totalMult,
					"locked_mults": lockedMults,
				},
			})
			time.Sleep(600 * time.Millisecond)
		}

		// Perfect cascade bonus
		if perfectCascade {
			bonusMult := 50.0
			totalMult += bonusMult
			g.broadcast(protocol.Envelope{
				Type: "lucky_cascade_lock",
				Payload: map[string]interface{}{
					"event":      "perfect_cascade",
					"bonus_mult": bonusMult,
					"total_mult": totalMult,
				},
			})
			g.sendAnnounce("Perfect Cascade! All 8 waves hit! Extra x50.0!", "critical", "#00FF7F")
		}

		betCost := float64(p.GetBetDef().BetCost)
		reward := int(totalMult * betCost)
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_cascade_lock",
			Payload: map[string]interface{}{
				"event":            "cascade_settle",
				"trigger_id":       p.ID,
				"trigger_name":     p.GetDisplayName(),
				"total_mult":       totalMult,
				"reward":           reward,
				"perfect_cascade":  perfectCascade,
				"locked_mults":     lockedMults,
			},
		})

		// Global boost x43.5 for 87 seconds
		g.broadcast(protocol.Envelope{
			Type: "lucky_cascade_lock",
			Payload: map[string]interface{}{
				"event":        "global_boost",
				"global_mult":  43.5,
				"duration":     87,
				"trigger_name": p.GetDisplayName(),
				"total_mult":   totalMult,
			},
		})
		g.sendAnnounce(fmt.Sprintf("Cascade Lock settled! %s locked x%.0f! Global x43.5 for 87s!", p.GetDisplayName(), totalMult), "critical", "#00BFFF")
		log.Printf("[LuckyCascadeLock] %s Cascade Lock settled %.0fx, global x43.5 for 87s", p.GetDisplayName(), totalMult)
	}()

	return true
}
