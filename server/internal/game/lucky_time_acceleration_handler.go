// lucky_time_acceleration_handler.go — T188 幸運時間加速魚
// 業界依據：Fishing Fortune「time warp」升級版 + 時間操控概念
// 設計：擊破後時間加速 30 秒（目標速度 ×0.15，射擊速度 ×3.0，獎勵 ×2.5）
//       30 秒內擊破 ≥20 個 → 時間完美：全服 ×18.0 加成 38 秒（新最高）
//       個人冷卻 110 秒；全服冷卻 170 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyTimeAccelerationManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *timeAccelerationPerfectBoost
	isActive     bool
	killCount    int
	triggerID    string
}

type timeAccelerationPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyTimeAccelerationManager() *luckyTimeAccelerationManager {
	return &luckyTimeAccelerationManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyTimeAccelerationFish(defID string) bool {
	return defID == "T188"
}

func (m *luckyTimeAccelerationManager) getTimeAccelerationMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyTimeAccelerationManager) isTimeAccelerationActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isActive
}

func (m *luckyTimeAccelerationManager) onKill(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isActive && playerID == m.triggerID {
		m.killCount++
	}
}

func (m *luckyTimeAccelerationManager) tryLuckyTimeAccelerationFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(110 * time.Second)
	m.globalCD = now.Add(170 * time.Second)
	m.isActive = true
	m.killCount = 0
	m.triggerID = p.ID
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_time_acceleration",
		Payload: map[string]interface{}{
			"event":          "time_acceleration_start",
			"trigger_id":     p.ID,
			"trigger_name":   p.GetDisplayName(),
			"duration":       30,
			"speed_factor":   0.15,
			"fire_rate_mult": 3.0,
			"reward_mult":    2.5,
		},
	})
	g.sendAnnounce(fmt.Sprintf("⚡🕐 時間加速！%s 啟動時間加速！目標速度 ×0.15！射擊速度 ×3.0！30 秒！", p.GetDisplayName()), "critical", "#E65100")
	log.Printf("[LuckyTimeAcceleration] %s 觸發時間加速魚", p.GetDisplayName())

	go func() {
		time.Sleep(30 * time.Second)

		m.mu.Lock()
		m.isActive = false
		kills := m.killCount
		m.mu.Unlock()

		if kills >= 20 {
			boostMult := 18.0
			boostSecs := 38
			m.mu.Lock()
			m.perfectBoost = &timeAccelerationPerfectBoost{
				mult:      boostMult,
				expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
			}
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_time_acceleration",
				Payload: map[string]interface{}{
					"event":        "time_acceleration_perfect",
					"kill_count":   kills,
					"boost_mult":   boostMult,
					"boost_secs":   boostSecs,
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
				},
			})
			g.sendAnnounce(fmt.Sprintf("⚡🏆 時間完美！%s 擊破 %d 個！全服 ×%.1f 加成 %d 秒！（新最高）", p.GetDisplayName(), kills, boostMult, boostSecs), "critical", "#BF360C")
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_time_acceleration",
				Payload: map[string]interface{}{
					"event":        "time_acceleration_end",
					"kill_count":   kills,
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
				},
			})
			g.sendAnnounce(fmt.Sprintf("⚡ 時間加速結束！%s 擊破 %d 個（需 20 個觸發完美）", p.GetDisplayName(), kills), "normal", "#FF6D00")
		}
	}()
	return true
}
