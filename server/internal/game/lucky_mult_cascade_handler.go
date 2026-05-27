// lucky_mult_cascade_handler.go — T158 幸運倍率瀑布魚
// 業界依據：Fishing Fortune「Multiplier Cascade — consecutive rare catches build 2x→500x」
// 設計：擊破後 30 秒倍率瀑布（每次擊破 +0.5x，最高 ×20.0），
//       30 秒內達到 ×15.0 → 完美瀑布：全服 ×6.5 加成 14 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyMultCascadeManager struct {
	mu            sync.Mutex
	personalCD    map[string]time.Time
	globalCD      time.Time
	activeSessions map[string]*multCascadeSession
	perfectBoost  *multCascadePerfectBoost
}

type multCascadePerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type multCascadeSession struct {
	playerID    string
	playerName  string
	currentMult float64
	killCount   int
	expiresAt   time.Time
	settled     bool
	perfect     bool
}

func newLuckyMultCascadeManager() *luckyMultCascadeManager {
	return &luckyMultCascadeManager{
		personalCD:     make(map[string]time.Time),
		activeSessions: make(map[string]*multCascadeSession),
	}
}

func isLuckyMultCascadeFish(defID string) bool {
	return defID == "T158"
}

func (m *luckyMultCascadeManager) getMultCascadePerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyMultCascadeManager) getMultCascadeKillBonus(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	if !ok || sess.settled || time.Now().After(sess.expiresAt) {
		return 1.0
	}
	return sess.currentMult
}

func (m *luckyMultCascadeManager) notifyMultCascadeKill(g *Game, p *Player) {
	m.mu.Lock()
	sess, ok := m.activeSessions[p.ID]
	if !ok || sess.settled || time.Now().After(sess.expiresAt) {
		m.mu.Unlock()
		return
	}
	sess.currentMult += 0.5
	if sess.currentMult > 20.0 {
		sess.currentMult = 20.0
	}
	sess.killCount++
	currentMult := sess.currentMult
	killCount := sess.killCount
	perfect := sess.perfect
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_mult_cascade",
		Payload: map[string]interface{}{
			"event":        "cascade_rise",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"current_mult": currentMult,
			"kill_count":   killCount,
		},
	})

	if currentMult >= 15.0 && !perfect {
		m.mu.Lock()
		if sess, ok := m.activeSessions[p.ID]; ok && !sess.perfect {
			sess.perfect = true
			m.perfectBoost = &multCascadePerfectBoost{
				mult:      6.5,
				expiresAt: time.Now().Add(14 * time.Second),
			}
		}
		m.mu.Unlock()
		g.broadcast(protocol.Envelope{
			Type: "lucky_mult_cascade",
			Payload: map[string]interface{}{
				"event":        "cascade_perfect",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"peak_mult":    currentMult,
				"boost_mult":   6.5,
				"boost_secs":   14,
			},
		})
		g.sendAnnounce(fmt.Sprintf("🌊✨ 完美瀑布！%s 倍率達 ×%.1f！全服 ×6.5 加成 14 秒！", p.GetDisplayName(), currentMult), "critical", "#1565C0")
		time.AfterFunc(14*time.Second, func() {
			m.mu.Lock()
			m.perfectBoost = nil
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_mult_cascade",
				Payload: map[string]interface{}{
					"event":      "cascade_perfect_end",
					"trigger_id": p.ID,
				},
			})
		})
	}
}

func (m *luckyMultCascadeManager) tryLuckyMultCascadeFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(40 * time.Second)
	m.globalCD = now.Add(65 * time.Second)
	sess := &multCascadeSession{
		playerID:    p.ID,
		playerName:  p.GetDisplayName(),
		currentMult: 1.0,
		expiresAt:   now.Add(30 * time.Second),
	}
	m.activeSessions[p.ID] = sess
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_mult_cascade",
		Payload: map[string]interface{}{
			"event":        "cascade_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     30,
			"max_mult":     20.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🌊 %s 啟動倍率瀑布！30 秒內每次擊破 +0.5x！最高 ×20.0！", p.GetDisplayName()), "high", "#1976D2")
	log.Printf("[LuckyMultCascade] %s 觸發倍率瀑布魚", p.GetDisplayName())

	go func() {
		time.Sleep(30 * time.Second)
		m.mu.Lock()
		sess, ok := m.activeSessions[p.ID]
		if !ok || sess.settled {
			m.mu.Unlock()
			return
		}
		finalMult := sess.currentMult
		killCount := sess.killCount
		sess.settled = true
		delete(m.activeSessions, p.ID)
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_mult_cascade",
			Payload: map[string]interface{}{
				"event":        "cascade_end",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"final_mult":   finalMult,
				"kill_count":   killCount,
			},
		})
	}()
	return true
}
