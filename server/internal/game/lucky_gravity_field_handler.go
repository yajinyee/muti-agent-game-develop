// lucky_gravity_field_handler.go — T187 幸運引力場魚
// 業界依據：Black Hole Fishing（SDG Games, 2026）引力吸引機制
// 設計：擊破後引力場 15 秒（所有目標速度 ×0.1，向中心聚集）
//       15 秒後引力爆炸（全場 HP -55%，每個獎勵 ×9.0）
//       爆炸命中 ≥12 → 引力完美：全服 ×17.5 加成 37 秒
//       個人冷卻 108 秒；全服冷卻 168 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyGravityFieldManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *gravityFieldPerfectBoost
	isActive     bool
}

type gravityFieldPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyGravityFieldManager() *luckyGravityFieldManager {
	return &luckyGravityFieldManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyGravityFieldFish(defID string) bool {
	return defID == "T187"
}

func (m *luckyGravityFieldManager) getGravityFieldMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyGravityFieldManager) isGravityActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isActive
}

func (m *luckyGravityFieldManager) tryLuckyGravityFieldFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(108 * time.Second)
	m.globalCD = now.Add(168 * time.Second)
	m.isActive = true
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_gravity_field",
		Payload: map[string]interface{}{
			"event":        "gravity_field_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     15,
			"slow_factor":  0.1,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🌀⚡ 引力場！%s 啟動引力場！所有目標速度 ×0.1！15 秒後引力爆炸！", p.GetDisplayName()), "critical", "#311B92")
	log.Printf("[LuckyGravityField] %s 觸發引力場魚", p.GetDisplayName())

	go func() {
		// 引力場持續 15 秒
		time.Sleep(15 * time.Second)

		m.mu.Lock()
		m.isActive = false
		m.mu.Unlock()

		// 引力爆炸：全場 HP -55%，每個獎勵 ×9.0
		hitCount := g.applyAOEDamage(0, 0, 99999, 0.55)

		g.broadcast(protocol.Envelope{
			Type: "lucky_gravity_field",
			Payload: map[string]interface{}{
				"event":        "gravity_explosion",
				"hit_count":    hitCount,
				"damage_pct":   0.55,
				"reward_mult":  9.0,
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
			},
		})

		if hitCount >= 12 {
			boostMult := 17.5
			boostSecs := 37
			m.mu.Lock()
			m.perfectBoost = &gravityFieldPerfectBoost{
				mult:      boostMult,
				expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
			}
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_gravity_field",
				Payload: map[string]interface{}{
					"event":        "gravity_perfect",
					"boost_mult":   boostMult,
					"boost_secs":   boostSecs,
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
				},
			})
			g.sendAnnounce(fmt.Sprintf("🌀🏆 引力完美！%s 命中 %d 個！全服 ×%.1f 加成 %d 秒！", p.GetDisplayName(), hitCount, boostMult, boostSecs), "critical", "#1A237E")
		} else {
			g.sendAnnounce(fmt.Sprintf("🌀💥 引力爆炸！%s 命中 %d 個！獎勵 ×9.0！", p.GetDisplayName(), hitCount), "special", "#4527A0")
		}
	}()
	return true
}
