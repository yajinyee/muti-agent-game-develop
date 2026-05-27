// lucky_final_power_handler.go — T180 幸運終焉之力魚
// 業界依據：「ultimate power mechanic」
// 設計：擊破後全場 HP 歸零（每個獎勵 ×10.0）
//       觸發全服 ×15.0 加成 30 秒（超越 T170 成為新最高倍率機制）
//       個人冷卻 90 秒；全服冷卻 140 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyFinalPowerManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *finalPowerPerfectBoost
}

type finalPowerPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyFinalPowerManager() *luckyFinalPowerManager {
	return &luckyFinalPowerManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyFinalPowerFish(defID string) bool {
	return defID == "T180"
}

func (m *luckyFinalPowerManager) getFinalPowerMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyFinalPowerManager) tryLuckyFinalPowerFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(90 * time.Second)
	m.globalCD = now.Add(140 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_final_power",
		Payload: map[string]interface{}{
			"event":        "final_power_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("💀🌌 終焉之力！%s 引動終焉！全場 HP 歸零！全服 ×15.0！", p.GetDisplayName()), "critical", "#B71C1C")
	log.Printf("[LuckyFinalPower] %s 觸發終焉之力魚", p.GetDisplayName())

	go func() {
		time.Sleep(600 * time.Millisecond)

		// 全場 HP 歸零，每個獎勵 ×10.0（超越 T170 的 ×8.0）
		hitCount := g.applyUltimateJudgment(p, 10.0)

		// 觸發全服 ×15.0 加成 30 秒（超越 T170 的 ×12.0）
		m.mu.Lock()
		m.perfectBoost = &finalPowerPerfectBoost{
			mult:      15.0,
			expiresAt: time.Now().Add(30 * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_final_power",
			Payload: map[string]interface{}{
				"event":        "final_power_complete",
				"hit_count":    hitCount,
				"boost_mult":   15.0,
				"boost_secs":   30,
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
			},
		})
		g.sendAnnounce(fmt.Sprintf("💀🏆 終焉之力完成！%s 清場 %d 個！全服 ×15.0 加成 30 秒！", p.GetDisplayName(), hitCount), "critical", "#D50000")
	}()
	return true
}
