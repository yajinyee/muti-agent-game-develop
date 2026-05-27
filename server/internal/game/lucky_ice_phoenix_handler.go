// lucky_ice_phoenix_handler.go — T156 幸運冰鳳凰魚
// 業界依據：Royal Fishing「Ice Phoenix 180-300x — freezes all fish, then phoenix rebirth explosion」
// 設計：擊破後冰凍全場 10 秒（傷害 ×1.5），冰凍結束鳳凰重生爆炸（HP -60%），
//       爆炸命中 ≥8 → 完美鳳凰：全服 ×5.5 加成 12 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyIcePhoenixManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	activeSession *icePhoenixSession
	perfectBoost  *icePhoenixPerfectBoost
}

type icePhoenixPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type icePhoenixSession struct {
	playerID   string
	playerName string
	killCount  int
	expiresAt  time.Time
	settled    bool
	isFrozen   bool
}

func newLuckyIcePhoenixManager() *luckyIcePhoenixManager {
	return &luckyIcePhoenixManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyIcePhoenixFish(defID string) bool {
	return defID == "T156"
}

func (m *luckyIcePhoenixManager) getIcePhoenixPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyIcePhoenixManager) isIcePhoenixFrozen() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession != nil && m.activeSession.isFrozen && time.Now().Before(m.activeSession.expiresAt) {
		return true
	}
	return false
}

func (m *luckyIcePhoenixManager) getIcePhoenixDamageMult() float64 {
	if m.isIcePhoenixFrozen() {
		return 1.5
	}
	return 1.0
}

func (m *luckyIcePhoenixManager) notifyIcePhoenixKill(g *Game, p *Player) {
	m.mu.Lock()
	sess := m.activeSession
	if sess == nil || sess.settled || time.Now().After(sess.expiresAt) {
		m.mu.Unlock()
		return
	}
	if !sess.isFrozen {
		m.mu.Unlock()
		return
	}
	sess.killCount++
	killCount := sess.killCount
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_ice_phoenix",
		Payload: map[string]interface{}{
			"event":        "phoenix_kill",
			"trigger_id":   sess.playerID,
			"trigger_name": sess.playerName,
			"kill_count":   killCount,
		},
	})
}

func (m *luckyIcePhoenixManager) tryLuckyIcePhoenixFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(35 * time.Second)
	m.globalCD = now.Add(55 * time.Second)
	sess := &icePhoenixSession{
		playerID:   p.ID,
		playerName: p.GetDisplayName(),
		expiresAt:  now.Add(10 * time.Second),
		isFrozen:   true,
	}
	m.activeSession = sess
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_ice_phoenix",
		Payload: map[string]interface{}{
			"event":        "freeze_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     10,
			"damage_mult":  1.5,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🧊🔥 %s 召喚冰鳳凰！全場凍結 10 秒！傷害 ×1.5！", p.GetDisplayName()), "high", "#00BCD4")
	log.Printf("[LuckyIcePhoenix] %s 觸發冰鳳凰", p.GetDisplayName())

	go func() {
		time.Sleep(10 * time.Second)
		m.mu.Lock()
		if m.activeSession == nil || m.activeSession.settled {
			m.mu.Unlock()
			return
		}
		m.activeSession.isFrozen = false
		killCount := m.activeSession.killCount
		m.mu.Unlock()

		// 鳳凰重生爆炸：全場 HP -60%
		hitCount := g.applyAOEDamage(GameWidth/2, GameHeight/2, 9999, 0.60)
		g.broadcast(protocol.Envelope{
			Type: "lucky_ice_phoenix",
			Payload: map[string]interface{}{
				"event":        "phoenix_rebirth",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"hit_count":    hitCount,
				"kill_count":   killCount,
			},
		})

		if hitCount >= 8 {
			m.mu.Lock()
			m.perfectBoost = &icePhoenixPerfectBoost{
				mult:      5.5,
				expiresAt: time.Now().Add(12 * time.Second),
			}
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_ice_phoenix",
				Payload: map[string]interface{}{
					"event":        "phoenix_perfect",
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
					"boost_mult":   5.5,
					"boost_secs":   12,
				},
			})
			g.sendAnnounce(fmt.Sprintf("🔥✨ 完美鳳凰！%s 命中 %d 個！全服 ×5.5 加成 12 秒！", p.GetDisplayName(), hitCount), "critical", "#FF6F00")
			time.AfterFunc(12*time.Second, func() {
				m.mu.Lock()
				m.perfectBoost = nil
				m.mu.Unlock()
				g.broadcast(protocol.Envelope{
					Type: "lucky_ice_phoenix",
					Payload: map[string]interface{}{
						"event":      "phoenix_perfect_end",
						"trigger_id": p.ID,
					},
				})
			})
		}

		m.mu.Lock()
		if m.activeSession != nil {
			m.activeSession.settled = true
			m.activeSession = nil
		}
		m.mu.Unlock()
	}()
	return true
}
