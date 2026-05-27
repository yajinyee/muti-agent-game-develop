// lucky_time_loop_handler.go — T177 幸運時間迴圈魚
// 業界依據：「time loop mechanic」
// 設計：擊破後 15 秒時間迴圈，每次迴圈重置目標 HP 並提高獎勵 ×1.5
//       最多 3 次迴圈，全部完成 → 時間完美：全服 ×10.0 加成 22 秒
//       個人冷卻 72 秒；全服冷卻 115 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyTimeLoopManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	active       *timeLoopSession
	perfectBoost *timeLoopPerfectBoost
}

type timeLoopSession struct {
	triggerID   string
	triggerName string
	loopCount   int
	maxLoops    int
	loopMult    float64
	expiresAt   time.Time
}

type timeLoopPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyTimeLoopManager() *luckyTimeLoopManager {
	return &luckyTimeLoopManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyTimeLoopFish(defID string) bool {
	return defID == "T177"
}

func (m *luckyTimeLoopManager) getTimeLoopPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyTimeLoopManager) getTimeLoopMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.active != nil && time.Now().Before(m.active.expiresAt) {
		return m.active.loopMult
	}
	return 1.0
}

func (m *luckyTimeLoopManager) tryLuckyTimeLoopFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(72 * time.Second)
	m.globalCD = now.Add(115 * time.Second)
	m.active = &timeLoopSession{
		triggerID:   p.ID,
		triggerName: p.GetDisplayName(),
		loopCount:   0,
		maxLoops:    3,
		loopMult:    1.5,
		expiresAt:   now.Add(15 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_time_loop",
		Payload: map[string]interface{}{
			"event":        "time_loop_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"max_loops":    3,
			"loop_mult":    1.5,
			"duration":     15,
		},
	})
	g.sendAnnounce(fmt.Sprintf("⏰ 時間迴圈！%s 開啟時間迴圈！每次迴圈獎勵 ×1.5！最多 3 次！", p.GetDisplayName()), "critical", "#1565C0")
	log.Printf("[LuckyTimeLoop] %s 觸發時間迴圈魚", p.GetDisplayName())

	go m.runLoops(g, p)
	return true
}

func (m *luckyTimeLoopManager) runLoops(g *Game, p *Player) {
	for loop := 1; loop <= 3; loop++ {
		time.Sleep(15 * time.Second)
		m.mu.Lock()
		if m.active == nil {
			m.mu.Unlock()
			return
		}
		m.active.loopCount = loop
		m.active.loopMult = 1.5 * float64(loop+1)
		m.active.expiresAt = time.Now().Add(15 * time.Second)
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_time_loop",
			Payload: map[string]interface{}{
				"event":        "loop_reset",
				"loop_no":      loop,
				"loop_mult":    m.active.loopMult,
				"trigger_name": p.GetDisplayName(),
			},
		})
		g.sendAnnounce(fmt.Sprintf("⏰ 第 %d 次迴圈！獎勵倍率 ×%.1f！", loop, m.active.loopMult), "high", "#1976D2")

		if loop == 3 {
			time.Sleep(15 * time.Second)
			m.mu.Lock()
			m.active = nil
			m.perfectBoost = &timeLoopPerfectBoost{
				mult:      10.0,
				expiresAt: time.Now().Add(22 * time.Second),
			}
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_time_loop",
				Payload: map[string]interface{}{
					"event":        "time_loop_perfect",
					"trigger_name": p.GetDisplayName(),
					"boost_mult":   10.0,
					"boost_secs":   22,
				},
			})
			g.sendAnnounce(fmt.Sprintf("⏰✨ 時間迴圈完美！%s 全服 ×10.0 加成 22 秒！", p.GetDisplayName()), "critical", "#42A5F5")
			return
		}
	}
}
