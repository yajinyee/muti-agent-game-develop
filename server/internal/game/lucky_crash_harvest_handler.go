// lucky_crash_harvest_handler.go — T227 幸運崩潰收割魚
// 設計：Crash Harvest 機制（Lucky Fish AbraCadabra 2026 + Crash mechanic 升級版）
//       倍率從 ×1.0 開始持續上升（每 0.5 秒 +0.5x），玩家可隨時收割
//       崩潰機率隨倍率增加（×10 時 5%，×20 時 15%，×50 時 40%）
//       完美收割（≥50x 且未崩潰）→ 全服 ×44.5 加成 89 秒（超越 T226 的 ×44.0）
//       最高倍率：×1000（極低機率）
//       觸發率：0.00025%；個人冷卻 345 秒；全服冷卻 405 秒
//       業界依據：Lucky Fish AbraCadabra「Crash mechanic」（2026-05）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyCrashHarvestManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
}

func newLuckyCrashHarvestManager() *luckyCrashHarvestManager {
	return &luckyCrashHarvestManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyCrashHarvestFish(defID string) bool {
	return defID == "T227"
}

func (m *luckyCrashHarvestManager) tryLuckyCrashHarvestFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(405 * time.Second)
	m.personalCD[p.ID] = now.Add(345 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_crash_harvest",
		Payload: map[string]interface{}{
			"event":        "crash_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"max_mult":     1000.0,
			"duration":     30,
		},
	})
	g.sendAnnounce(fmt.Sprintf("Crash Harvest! %s triggered Crash! Multiplier rising... cash out before it crashes!", p.GetDisplayName()), "critical", "#FF4500")
	log.Printf("[LuckyCrashHarvest] %s triggered Crash Harvest fish", p.GetDisplayName())

	go func() {
		currentMult := 1.0
		crashed := false
		perfectHarvest := false
		maxReached := 0.0

		// Simulate crash game: multiplier rises until crash
		for tick := 0; tick < 60; tick++ { // max 30 seconds (0.5s per tick)
			currentMult += 0.5 + rand.Float64()*0.5 // +0.5 to +1.0 per tick

			// Crash probability increases with multiplier
			crashProb := 0.0
			switch {
			case currentMult >= 100.0:
				crashProb = 0.60
			case currentMult >= 50.0:
				crashProb = 0.40
			case currentMult >= 20.0:
				crashProb = 0.15
			case currentMult >= 10.0:
				crashProb = 0.05
			default:
				crashProb = 0.01
			}

			if rand.Float64() < crashProb {
				crashed = true
				break
			}

			maxReached = currentMult

			g.broadcast(protocol.Envelope{
				Type: "lucky_crash_harvest",
				Payload: map[string]interface{}{
					"event":        "mult_tick",
					"current_mult": currentMult,
					"tick":         tick + 1,
				},
			})
			time.Sleep(500 * time.Millisecond)
		}

		// Auto-harvest at max if not crashed
		harvestMult := maxReached
		if crashed {
			// Crashed: get 20% of current mult as consolation
			harvestMult = currentMult * 0.20
			if harvestMult < 1.0 {
				harvestMult = 1.0
			}
			g.broadcast(protocol.Envelope{
				Type: "lucky_crash_harvest",
				Payload: map[string]interface{}{
					"event":         "crashed",
					"crash_mult":    currentMult,
					"harvest_mult":  harvestMult,
					"consolation":   true,
				},
			})
			g.sendAnnounce(fmt.Sprintf("CRASHED at x%.1f! Consolation x%.1f", currentMult, harvestMult), "warning", "#FF0000")
		} else {
			// Perfect harvest
			if harvestMult >= 50.0 {
				perfectHarvest = true
			}
			g.broadcast(protocol.Envelope{
				Type: "lucky_crash_harvest",
				Payload: map[string]interface{}{
					"event":          "harvested",
					"harvest_mult":   harvestMult,
					"perfect":        perfectHarvest,
				},
			})
		}

		betCost := float64(p.GetBetDef().BetCost)
		reward := int(harvestMult * betCost)
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_crash_harvest",
			Payload: map[string]interface{}{
				"event":           "crash_settle",
				"trigger_id":      p.ID,
				"trigger_name":    p.GetDisplayName(),
				"harvest_mult":    harvestMult,
				"reward":          reward,
				"crashed":         crashed,
				"perfect_harvest": perfectHarvest,
			},
		})

		if perfectHarvest {
			// Global boost x44.5 for 89 seconds
			g.broadcast(protocol.Envelope{
				Type: "lucky_crash_harvest",
				Payload: map[string]interface{}{
					"event":        "global_boost",
					"global_mult":  44.5,
					"duration":     89,
					"trigger_name": p.GetDisplayName(),
					"harvest_mult": harvestMult,
				},
			})
			g.sendAnnounce(fmt.Sprintf("Perfect Harvest! %s cashed out x%.1f! Global x44.5 for 89s!", p.GetDisplayName(), harvestMult), "critical", "#FF4500")
		} else {
			g.sendAnnounce(fmt.Sprintf("Crash Harvest settled! %s got x%.1f!", p.GetDisplayName(), harvestMult), "info", "#FFA500")
		}
		log.Printf("[LuckyCrashHarvest] %s Crash Harvest settled %.1fx (crashed=%v, perfect=%v)", p.GetDisplayName(), harvestMult, crashed, perfectHarvest)
	}()

	return true
}
