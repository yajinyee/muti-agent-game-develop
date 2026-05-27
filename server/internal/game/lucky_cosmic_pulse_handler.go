// lucky_cosmic_pulse_handler.go — T185 幸運宇宙脈衝魚
// 業界依據：Fishing Fortune multiplier cascade（2x→500x）升級版
// 設計：擊破後宇宙脈衝波（全場 HP -45%，每個獎勵 ×12.0）
//       觸發全服 ×16.0 加成 35 秒（超越 T180 成為新最高倍率機制）
//       個人冷卻 100 秒；全服冷卻 160 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyCosmicPulseManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *cosmicPulsePerfectBoost
}

type cosmicPulsePerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyCosmicPulseManager() *luckyCosmicPulseManager {
	return &luckyCosmicPulseManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyCosmicPulseFish(defID string) bool {
	return defID == "T185"
}

func (m *luckyCosmicPulseManager) getCosmicPulseMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyCosmicPulseManager) tryLuckyCosmicPulseFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(100 * time.Second)
	m.globalCD = now.Add(160 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_cosmic_pulse",
		Payload: map[string]interface{}{
			"event":        "cosmic_pulse_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("🌌💥 宇宙脈衝！%s 引動宇宙脈衝波！全場 HP -45%%！全服 ×16.0！", p.GetDisplayName()), "critical", "#4A148C")
	log.Printf("[LuckyCosmicPulse] %s 觸發宇宙脈衝魚", p.GetDisplayName())

	go func() {
		time.Sleep(600 * time.Millisecond)

		// 宇宙脈衝波：全場 HP -45%，每個獎勵 ×12.0（超越 T180 的 ×10.0）
		hitCount := g.applyUltimateJudgment(p, 12.0)

		// 觸發全服 ×16.0 加成 35 秒（超越 T185 成為新最高倍率機制）
		boostMult := 16.0
		boostSecs := 35
		m.mu.Lock()
		m.perfectBoost = &cosmicPulsePerfectBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_cosmic_pulse",
			Payload: map[string]interface{}{
				"event":        "cosmic_pulse_complete",
				"hit_count":    hitCount,
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
			},
		})
		g.sendAnnounce(fmt.Sprintf("🌌🏆 宇宙脈衝完成！%s 清場 %d 個！全服 ×%.1f 加成 %d 秒！（新最高）", p.GetDisplayName(), hitCount, boostMult, boostSecs), "critical", "#311B92")
	}()
	return true
}
