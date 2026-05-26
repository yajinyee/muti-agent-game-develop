// lucky_vampire_v2_handler.go — T152 幸運吸血鬼升級魚
// 業界依據：Jili「Vampire multiplier increases the more you fight, chance to enter multiplier mode up to X5」升級版
// 設計：擊破後 25 秒吸血模式，每次擊破 +1.5x（最高 ×10.0），吸收 ≥10 次 → 完美吸血全服 ×4.0 加成 10 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyVampireV2Manager struct {
	mu             sync.Mutex
	personalCD     map[string]time.Time
	globalCD       time.Time
	activeSessions map[string]*vampireV2Session
	perfectBoost   *vampireV2PerfectBoost
}

type vampireV2PerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type vampireV2Session struct {
	playerID    string
	playerName  string
	absorbCount int
	currentMult float64
	inMultMode  bool
	multExpires time.Time
	expiresAt   time.Time
	settled     bool
}

func newLuckyVampireV2Manager() *luckyVampireV2Manager {
	return &luckyVampireV2Manager{
		personalCD:     make(map[string]time.Time),
		activeSessions: make(map[string]*vampireV2Session),
	}
}

func isLuckyVampireV2Fish(defID string) bool {
	return defID == "T152"
}

func (m *luckyVampireV2Manager) getVampireV2PerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyVampireV2Manager) getVampireV2KillMult(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	if !ok || sess.settled {
		return 1.0
	}
	if sess.inMultMode && time.Now().Before(sess.multExpires) {
		return sess.currentMult
	}
	return 1.0
}

func (m *luckyVampireV2Manager) notifyVampireV2Kill(g *Game, p *Player) {
	m.mu.Lock()
	sess, ok := m.activeSessions[p.ID]
	if !ok || sess.settled || time.Now().After(sess.expiresAt) {
		m.mu.Unlock()
		return
	}
	sess.absorbCount++
	newMult := 1.0 + float64(sess.absorbCount)*1.5
	if newMult > 10.0 {
		newMult = 10.0
	}
	sess.currentMult = newMult
	absorbCount := sess.absorbCount
	currentMult := sess.currentMult
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_vampire_v2",
		Payload: map[string]interface{}{
			"event":        "absorb_v2",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"absorb_count": absorbCount,
			"current_mult": currentMult,
		},
	})

	if absorbCount == 10 {
		m.mu.Lock()
		sess.inMultMode = true
		sess.multExpires = time.Now().Add(12 * time.Second)
		m.mu.Unlock()
		g.broadcast(protocol.Envelope{
			Type: "lucky_vampire_v2",
			Payload: map[string]interface{}{
				"event":        "mult_mode_v2",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"current_mult": currentMult,
				"time_left":    12.0,
			},
		})
		g.sendAnnounce(fmt.Sprintf("🧛 %s 進入吸血鬼倍率模式！×%.1f！12 秒！", p.GetDisplayName(), currentMult), "high", "#9C27B0")
	}
}

func (m *luckyVampireV2Manager) tryLuckyVampireV2Fish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(22 * time.Second)
	m.globalCD = now.Add(36 * time.Second)
	sess := &vampireV2Session{
		playerID:    p.ID,
		playerName:  p.GetDisplayName(),
		currentMult: 1.0,
		expiresAt:   now.Add(25 * time.Second),
	}
	m.activeSessions[p.ID] = sess
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_vampire_v2",
		Payload: map[string]interface{}{
			"event":        "vampire_v2_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     25,
			"max_mult":     10.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🧛 %s 觸發吸血鬼升級！25 秒吸血模式，最高 ×10.0！", p.GetDisplayName()), "high", "#9C27B0")
	log.Printf("[LuckyVampireV2] %s 觸發吸血鬼升級", p.GetDisplayName())

	go func() {
		time.Sleep(25 * time.Second)
		m.mu.Lock()
		sess.settled = true
		absorbCount := sess.absorbCount
		currentMult := sess.currentMult
		delete(m.activeSessions, p.ID)
		m.mu.Unlock()

		if absorbCount >= 10 {
			m.mu.Lock()
			m.perfectBoost = &vampireV2PerfectBoost{
				mult:      4.0,
				expiresAt: time.Now().Add(10 * time.Second),
			}
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_vampire_v2",
				Payload: map[string]interface{}{
					"event":        "vampire_v2_perfect",
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
					"absorb_count": absorbCount,
					"final_mult":   currentMult,
					"boost_mult":   4.0,
					"boost_secs":   10,
				},
			})
			g.sendAnnounce(fmt.Sprintf("🧛✨ 完美吸血！%s 吸收 %d 次！全服 ×4.0 加成 10 秒！", p.GetDisplayName(), absorbCount), "high", "#FFD700")
			time.AfterFunc(10*time.Second, func() {
				m.mu.Lock()
				m.perfectBoost = nil
				m.mu.Unlock()
				g.broadcast(protocol.Envelope{
					Type: "lucky_vampire_v2",
					Payload: map[string]interface{}{
						"event":      "vampire_v2_perfect_end",
						"trigger_id": p.ID,
					},
				})
			})
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_vampire_v2",
				Payload: map[string]interface{}{
					"event":        "vampire_v2_end",
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
					"absorb_count": absorbCount,
					"final_mult":   currentMult,
				},
			})
		}
	}()
	return true
}
