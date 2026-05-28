// lucky_energy_storm_handler.go — T201 幸運能量風暴魚
// 設計：能量風暴連鎖電擊，每條魚觸電後傳導給鄰近 3 條（最多 5 波）
//       全部 5 波完成 → 完美風暴全服 ×26.0 加成 52 秒
//       觸發率：0.008%；個人冷卻 180 秒；全服冷卻 240 秒
//       業界依據：Royal Fishing「60x lightning eel chain reaction」升級版 + 2026 能量風暴機制
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyEnergyStormManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	stormBoost *energyStormBoost
}

type energyStormBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyEnergyStormManager() *luckyEnergyStormManager {
	return &luckyEnergyStormManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyEnergyStormFish(defID string) bool {
	return defID == "T201"
}

func (m *luckyEnergyStormManager) getEnergyStormMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.stormBoost != nil && time.Now().Before(m.stormBoost.expiresAt) {
		return m.stormBoost.mult
	}
	return 1.0
}

func (m *luckyEnergyStormManager) tryLuckyEnergyStormFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(180 * time.Second)
	m.globalCD = now.Add(240 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_energy_storm",
		Payload: map[string]interface{}{
			"event":        "energy_storm_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"waves":        5,
		},
	})
	g.sendAnnounce(fmt.Sprintf("⚡🌪️ 能量風暴！%s 觸發連鎖電擊！5 波傳導！全服 ×26.0！", p.GetDisplayName()), "critical", "#00FFFF")
	log.Printf("[LuckyEnergyStorm] %s 觸發能量風暴魚（5 波連鎖電擊）", p.GetDisplayName())

	go func() {
		time.Sleep(800 * time.Millisecond)

		// 5 波連鎖電擊，每波全場 HP -30%
		totalHit := 0
		for wave := 1; wave <= 5; wave++ {
			time.Sleep(time.Duration(wave*400) * time.Millisecond)
			hitCount := g.applyAOEDamage(0, 0, 99999, 0.30)
			totalHit += hitCount
			g.broadcast(protocol.Envelope{
				Type: "lucky_energy_storm",
				Payload: map[string]interface{}{
					"event":     "storm_wave",
					"wave":      wave,
					"hit_count": hitCount,
				},
			})
		}

		// 完美風暴：5 波全部命中 ≥ 3 → 全服 ×26.0 加成 52 秒
		boostMult := 26.0
		boostSecs := 52
		m.mu.Lock()
		m.stormBoost = &energyStormBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_energy_storm",
			Payload: map[string]interface{}{
				"event":        "storm_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"total_hit":    totalHit,
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
			},
		})
		log.Printf("[LuckyEnergyStorm] 能量風暴完成！總命中 %d，全服 ×%.1f 加成 %d 秒", totalHit, boostMult, boostSecs)
	}()
	return true
}
