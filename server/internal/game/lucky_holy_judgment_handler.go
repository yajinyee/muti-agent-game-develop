// lucky_holy_judgment_handler.go — T169 幸運神聖審判魚
// 業界依據：Jili「Super Awakening 3000x」升級版
// 設計：擊破後神聖審判 25 秒，每 5 秒一波神聖光柱（全場 HP -30%）
//       5 波全部命中 ≥5 個目標 → 神聖完美：全服 ×8.5 加成 18 秒
//       個人冷卻 62 秒；全服冷卻 95 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyHolyJudgmentManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	sessions     map[string]*holyJudgmentSession
	perfectBoost *holyJudgmentPerfectBoost
}

type holyJudgmentSession struct {
	playerID      string
	waveHits      [5]int // 每波命中數
	currentWave   int
	expiresAt     time.Time
}

type holyJudgmentPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyHolyJudgmentManager() *luckyHolyJudgmentManager {
	return &luckyHolyJudgmentManager{
		personalCD: make(map[string]time.Time),
		sessions:   make(map[string]*holyJudgmentSession),
	}
}

func isLuckyHolyJudgmentFish(defID string) bool {
	return defID == "T169"
}

func (m *luckyHolyJudgmentManager) getHolyJudgmentPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyHolyJudgmentManager) tryLuckyHolyJudgmentFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(62 * time.Second)
	m.globalCD = now.Add(95 * time.Second)

	sess := &holyJudgmentSession{
		playerID:  p.ID,
		expiresAt: now.Add(25 * time.Second),
	}
	m.sessions[p.ID] = sess
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_holy_judgment",
		Payload: map[string]interface{}{
			"event":        "judgment_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     25,
			"wave_count":   5,
		},
	})
	g.sendAnnounce(fmt.Sprintf("✨⚖️ 神聖審判！%s 召喚神聖光柱！5 波審判降臨！", p.GetDisplayName()), "high", "#F57F17")
	log.Printf("[LuckyHolyJudgment] %s 觸發神聖審判魚", p.GetDisplayName())

	// 5 波，每 5 秒一波
	go func() {
		allPerfect := true
		for wave := 1; wave <= 5; wave++ {
			time.Sleep(5 * time.Second)

			// 神聖光柱：全場 HP -30%
			hitCount := g.applyAOEDamage(640, 360, 99999, 0.3)

			m.mu.Lock()
			if sess, ok := m.sessions[p.ID]; ok && wave <= 5 {
				sess.waveHits[wave-1] = hitCount
				if hitCount < 5 {
					allPerfect = false
				}
			}
			m.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: "lucky_holy_judgment",
				Payload: map[string]interface{}{
					"event":      "judgment_wave",
					"wave":       wave,
					"hit_count":  hitCount,
					"trigger_id": p.ID,
				},
			})
		}

		// 結算
		m.mu.Lock()
		delete(m.sessions, p.ID)
		m.mu.Unlock()

		if allPerfect {
			m.mu.Lock()
			m.perfectBoost = &holyJudgmentPerfectBoost{
				mult:      8.5,
				expiresAt: time.Now().Add(18 * time.Second),
			}
			m.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: "lucky_holy_judgment",
				Payload: map[string]interface{}{
					"event":        "judgment_perfect",
					"boost_mult":   8.5,
					"boost_secs":   18,
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
				},
			})
			g.sendAnnounce(fmt.Sprintf("✨🏆 神聖完美！%s 5 波全部命中！全服 ×8.5 加成 18 秒！", p.GetDisplayName()), "critical", "#FF8F00")
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_holy_judgment",
				Payload: map[string]interface{}{
					"event":      "judgment_end",
					"trigger_id": p.ID,
				},
			})
		}
	}()

	return true
}
