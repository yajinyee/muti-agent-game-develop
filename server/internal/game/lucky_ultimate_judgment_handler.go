// lucky_ultimate_judgment_handler.go — T160 幸運終極審判魚
// 業界依據：終極機制 — 全場目標 HP 歸零（每個獎勵 ×6.0），觸發全服 ×10.0 加成 20 秒
// 設計：擊破後全場所有目標 HP 歸零（每個獎勵 ×6.0），觸發全服 ×10.0 加成 20 秒
//       這是遊戲中倍率最高的機制，個人冷卻 60 秒，全服冷卻 90 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyUltimateJudgmentManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *ultimateJudgmentPerfectBoost
}

type ultimateJudgmentPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyUltimateJudgmentManager() *luckyUltimateJudgmentManager {
	return &luckyUltimateJudgmentManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyUltimateJudgmentFish(defID string) bool {
	return defID == "T160"
}

func (m *luckyUltimateJudgmentManager) getUltimateJudgmentPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyUltimateJudgmentManager) tryLuckyUltimateJudgmentFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(60 * time.Second)
	m.globalCD = now.Add(90 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_ultimate_judgment",
		Payload: map[string]interface{}{
			"event":        "judgment_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("⚖️💥 終極審判！%s 降下審判！全場 HP 歸零！", p.GetDisplayName()), "critical", "#B71C1C")
	log.Printf("[LuckyUltimateJudgment] %s 觸發終極審判魚", p.GetDisplayName())

	go func() {
		time.Sleep(300 * time.Millisecond)

		// 全場 HP 歸零，每個獎勵 ×6.0
		hitCount := g.applyUltimateJudgment(p, 6.0)

		g.broadcast(protocol.Envelope{
			Type: "lucky_ultimate_judgment",
			Payload: map[string]interface{}{
				"event":        "judgment_execute",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"hit_count":    hitCount,
				"reward_mult":  6.0,
			},
		})

		// 無論命中多少，都觸發全服 ×10.0 加成 20 秒
		m.mu.Lock()
		m.perfectBoost = &ultimateJudgmentPerfectBoost{
			mult:      10.0,
			expiresAt: time.Now().Add(20 * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_ultimate_judgment",
			Payload: map[string]interface{}{
				"event":        "judgment_boost",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"hit_count":    hitCount,
				"boost_mult":   10.0,
				"boost_secs":   20,
			},
		})
		g.sendAnnounce(fmt.Sprintf("⚖️✨ 終極審判完成！%s 清場 %d 個！全服 ×10.0 加成 20 秒！", p.GetDisplayName(), hitCount), "critical", "#FFD700")

		time.AfterFunc(20*time.Second, func() {
			m.mu.Lock()
			m.perfectBoost = nil
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_ultimate_judgment",
				Payload: map[string]interface{}{
					"event":      "judgment_boost_end",
					"trigger_id": p.ID,
				},
			})
		})
	}()
	return true
}
