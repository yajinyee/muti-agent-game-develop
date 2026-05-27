// lucky_arctic_storm_handler.go — T182 幸運北極風暴魚
// 業界依據：Arctic Mechanics（500x multiplier, fast pace rounds）
// 設計：擊破後快速節奏 8 波冰雪攻擊（每 0.3 秒一波，每波 HP -15%）
//       全部 8 波命中 ≥3 個目標 → 全服 ×16.5 加成 33 秒
//       個人冷卻 95 秒；全服冷卻 148 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/data"
	"chiikawa-game/internal/protocol"
)

type luckyArcticStormManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *arcticStormPerfectBoost
}

type arcticStormPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyArcticStormManager() *luckyArcticStormManager {
	return &luckyArcticStormManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyArcticStormFish(defID string) bool {
	return defID == "T182"
}

func (m *luckyArcticStormManager) getArcticStormMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyArcticStormManager) tryLuckyArcticStormFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(95 * time.Second)
	m.globalCD = now.Add(148 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_arctic_storm",
		Payload: map[string]interface{}{
			"event":        "arctic_storm_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"total_waves":  8,
		},
	})
	g.sendAnnounce(fmt.Sprintf("❄️⚡ 北極風暴！%s 引動 8 波冰雪攻擊！", p.GetDisplayName()), "special", "#0288D1")
	log.Printf("[LuckyArcticStorm] %s 觸發北極風暴魚", p.GetDisplayName())

	go func() {
		totalHit := 0
		perfectWaves := 0
		const totalWaves = 8
		const waveInterval = 300 * time.Millisecond
		const damagePercent = 0.15

		for wave := 1; wave <= totalWaves; wave++ {
			time.Sleep(waveInterval)

			waveHit := 0
			g.mu.Lock()
			for _, t := range g.targets {
				if t.HP > 0 && t.Def.Type != data.TypeBoss {
					damage := int(float64(t.HP) * damagePercent)
					if damage < 1 {
						damage = 1
					}
					t.HP -= damage
					if t.HP <= 0 {
						t.HP = 0
					}
					waveHit++
					totalHit++
				}
			}
			g.mu.Unlock()

			if waveHit >= 3 {
				perfectWaves++
			}

			g.broadcast(protocol.Envelope{
				Type: "lucky_arctic_storm",
				Payload: map[string]interface{}{
					"event":      "arctic_wave",
					"wave":       wave,
					"wave_hit":   waveHit,
					"total_hit":  totalHit,
					"trigger_id": p.ID,
				},
			})
		}

		// 全部 8 波都命中 ≥3 個目標 → 完美北極風暴
		if perfectWaves >= 8 {
			boostMult := 16.5
			boostSecs := 33
			m.mu.Lock()
			m.perfectBoost = &arcticStormPerfectBoost{
				mult:      boostMult,
				expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
			}
			m.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: "lucky_arctic_storm",
				Payload: map[string]interface{}{
					"event":         "arctic_storm_perfect",
					"total_hit":     totalHit,
					"perfect_waves": perfectWaves,
					"boost_mult":    boostMult,
					"boost_secs":    boostSecs,
					"trigger_id":    p.ID,
					"trigger_name":  p.GetDisplayName(),
				},
			})
			g.sendAnnounce(fmt.Sprintf("❄️🏆 完美北極風暴！%s 8 波全中！全服 ×%.1f 加成 %d 秒！", p.GetDisplayName(), boostMult, boostSecs), "critical", "#01579B")
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_arctic_storm",
				Payload: map[string]interface{}{
					"event":         "arctic_storm_complete",
					"total_hit":     totalHit,
					"perfect_waves": perfectWaves,
					"trigger_id":    p.ID,
					"trigger_name":  p.GetDisplayName(),
				},
			})
		}
	}()
	return true
}
