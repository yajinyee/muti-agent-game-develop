// lucky_big_bang_handler.go — T170 幸運宇宙大爆炸魚
// 業界依據：Fishing Carnival「Big Bang mechanic」
// 設計：擊破後宇宙大爆炸，全場所有目標 HP 歸零（每個獎勵 ×8.0）
//       觸發全服 ×12.0 加成 25 秒（超越 T160 的最高倍率機制）
//       個人冷卻 70 秒；全服冷卻 110 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyBigBangManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *bigBangPerfectBoost
}

type bigBangPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyBigBangManager() *luckyBigBangManager {
	return &luckyBigBangManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyBigBangFish(defID string) bool {
	return defID == "T170"
}

func (m *luckyBigBangManager) getBigBangPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyBigBangManager) tryLuckyBigBangFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(70 * time.Second)
	m.globalCD = now.Add(110 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_big_bang",
		Payload: map[string]interface{}{
			"event":        "big_bang_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("💥🌌 宇宙大爆炸！%s 引爆宇宙！全場 HP 歸零！全服 ×12.0！", p.GetDisplayName()), "critical", "#B71C1C")
	log.Printf("[LuckyBigBang] %s 觸發宇宙大爆炸魚", p.GetDisplayName())

	go func() {
		time.Sleep(500 * time.Millisecond)

		// 全場 HP 歸零，每個獎勵 ×8.0
		hitCount := g.applyUltimateJudgment(p, 8.0)

		// 觸發全服 ×12.0 加成 25 秒（遊戲最高倍率）
		m.mu.Lock()
		m.perfectBoost = &bigBangPerfectBoost{
			mult:      12.0,
			expiresAt: time.Now().Add(25 * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_big_bang",
			Payload: map[string]interface{}{
				"event":        "big_bang_complete",
				"hit_count":    hitCount,
				"boost_mult":   12.0,
				"boost_secs":   25,
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
			},
		})
		g.sendAnnounce(fmt.Sprintf("💥🏆 宇宙大爆炸完成！%s 清場 %d 個！全服 ×12.0 加成 25 秒！", p.GetDisplayName(), hitCount), "critical", "#FF1744")
	}()

	return true
}
