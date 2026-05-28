// lucky_crystal_resonance_handler.go — T202 幸運水晶共鳴魚
// 設計：水晶共鳴，場上所有目標同時爆炸（每個獎勵 ×30.0）
//       觸發後全服 ×27.0 加成 54 秒（超越 T201 的 ×26.0）
//       觸發率：0.007%；個人冷卻 180 秒；全服冷卻 240 秒
//       業界依據：Fishing Legend 2025「Crystal Resonance」全場共鳴爆炸機制
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyCrystalResonanceManager struct {
	mu              sync.Mutex
	personalCD      map[string]time.Time
	globalCD        time.Time
	resonanceBoost  *crystalResonanceBoost
}

type crystalResonanceBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyCrystalResonanceManager() *luckyCrystalResonanceManager {
	return &luckyCrystalResonanceManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyCrystalResonanceFish(defID string) bool {
	return defID == "T202"
}

func (m *luckyCrystalResonanceManager) getCrystalResonanceMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.resonanceBoost != nil && time.Now().Before(m.resonanceBoost.expiresAt) {
		return m.resonanceBoost.mult
	}
	return 1.0
}

func (m *luckyCrystalResonanceManager) tryLuckyCrystalResonanceFish(g *Game, p *Player) bool {
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
		Type: "lucky_crystal_resonance",
		Payload: map[string]interface{}{
			"event":        "crystal_resonance_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("💎✨ 水晶共鳴！%s 觸發全場共鳴爆炸！每個獎勵 ×30.0！全服 ×27.0！", p.GetDisplayName()), "critical", "#E0E0FF")
	log.Printf("[LuckyCrystalResonance] %s 觸發水晶共鳴魚（全場共鳴爆炸 ×30.0）", p.GetDisplayName())

	go func() {
		time.Sleep(1000 * time.Millisecond)

		// 水晶共鳴：全場 HP 歸零，每個獎勵 ×30.0
		hitCount := g.applyUltimateJudgment(p, 30.0)

		// 觸發全服 ×27.0 加成 54 秒
		boostMult := 27.0
		boostSecs := 54
		m.mu.Lock()
		m.resonanceBoost = &crystalResonanceBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_crystal_resonance",
			Payload: map[string]interface{}{
				"event":        "crystal_resonance_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"hit_count":    hitCount,
				"reward_mult":  30.0,
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
			},
		})
		log.Printf("[LuckyCrystalResonance] 水晶共鳴完成！命中 %d 個目標，全服 ×%.1f 加成 %d 秒", hitCount, boostMult, boostSecs)
	}()
	return true
}
