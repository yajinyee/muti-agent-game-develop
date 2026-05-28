// lucky_cosmic_singularity_handler.go — T205 幸運宇宙奇點魚
// 設計：宇宙奇點，全場 HP 歸零（每個獎勵 ×30.0）
//       觸發後全服 ×30.0 加成 60 秒（史上最高，超越 T200 的 ×25.0）
//       觸發率：0.003%（最稀有）；個人冷卻 240 秒；全服冷卻 300 秒
//       業界依據：終極宇宙奇點機制 + 2026 最高倍率設計
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyCosmicSingularityManager struct {
	mu               sync.Mutex
	personalCD       map[string]time.Time
	globalCD         time.Time
	singularityBoost *cosmicSingularityBoost
}

type cosmicSingularityBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyCosmicSingularityManager() *luckyCosmicSingularityManager {
	return &luckyCosmicSingularityManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyCosmicSingularityFish(defID string) bool {
	return defID == "T205"
}

func (m *luckyCosmicSingularityManager) getCosmicSingularityMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.singularityBoost != nil && time.Now().Before(m.singularityBoost.expiresAt) {
		return m.singularityBoost.mult
	}
	return 1.0
}

func (m *luckyCosmicSingularityManager) tryLuckyCosmicSingularityFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(240 * time.Second)
	m.globalCD = now.Add(300 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_cosmic_singularity",
		Payload: map[string]interface{}{
			"event":        "singularity_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("🌌💥 宇宙奇點！%s 開啟宇宙奇點！全場 HP 歸零！每個獎勵 ×30.0！全服 ×30.0！", p.GetDisplayName()), "critical", "#FF00FF")
	log.Printf("[LuckyCosmicSingularity] %s 觸發宇宙奇點魚（史上最高全服倍率 ×30.0）", p.GetDisplayName())

	go func() {
		time.Sleep(1200 * time.Millisecond)

		// 宇宙奇點：全場 HP 歸零，每個獎勵 ×30.0（史上最高單次清場倍率）
		hitCount := g.applyUltimateJudgment(p, 30.0)

		// 觸發全服 ×30.0 加成 60 秒（史上最高全服倍率，超越 T200 的 ×25.0）
		boostMult := 30.0
		boostSecs := 60
		m.mu.Lock()
		m.singularityBoost = &cosmicSingularityBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_cosmic_singularity",
			Payload: map[string]interface{}{
				"event":        "singularity_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"hit_count":    hitCount,
				"reward_mult":  30.0,
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
			},
		})
		log.Printf("[LuckyCosmicSingularity] 宇宙奇點完成！命中 %d 個目標，全服 ×%.1f 加成 %d 秒（史上最高）", hitCount, boostMult, boostSecs)
	}()
	return true
}
