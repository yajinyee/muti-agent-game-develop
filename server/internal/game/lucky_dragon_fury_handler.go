// lucky_dragon_fury_handler.go — T157 幸運龍怒能量魚
// 業界依據：Royal Fishing「Dragon Fury — energy accumulation → full-screen attack」
// 設計：擊破後能量累積 15 秒（每次擊破 +10 能量），滿 100 → 龍怒全場（HP -80%）
//       龍怒命中 ≥10 → 完美龍怒：全服 ×6.0 加成 13 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyDragonFuryManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	activeSession *dragonFurySession
	perfectBoost  *dragonFuryPerfectBoost
}

type dragonFuryPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type dragonFurySession struct {
	playerID   string
	playerName string
	energy     int
	expiresAt  time.Time
	settled    bool
	furyActive bool
}

func newLuckyDragonFuryManager() *luckyDragonFuryManager {
	return &luckyDragonFuryManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyDragonFuryFish(defID string) bool {
	return defID == "T157"
}

func (m *luckyDragonFuryManager) getDragonFuryPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyDragonFuryManager) notifyDragonFuryKill(g *Game, p *Player) {
	m.mu.Lock()
	sess := m.activeSession
	if sess == nil || sess.settled || time.Now().After(sess.expiresAt) {
		m.mu.Unlock()
		return
	}
	if sess.furyActive {
		m.mu.Unlock()
		return
	}
	sess.energy += 10
	if sess.energy > 100 {
		sess.energy = 100
	}
	energy := sess.energy
	playerID := sess.playerID
	playerName := sess.playerName
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_dragon_fury",
		Payload: map[string]interface{}{
			"event":        "energy_gain",
			"trigger_id":   playerID,
			"trigger_name": playerName,
			"energy":       energy,
		},
	})

	if energy >= 100 {
		m.mu.Lock()
		if m.activeSession != nil && !m.activeSession.furyActive {
			m.activeSession.furyActive = true
		}
		m.mu.Unlock()
		go m.triggerDragonFury(g, p, playerID, playerName)
	}
}

func (m *luckyDragonFuryManager) triggerDragonFury(g *Game, p *Player, playerID, playerName string) {
	g.broadcast(protocol.Envelope{
		Type: "lucky_dragon_fury",
		Payload: map[string]interface{}{
			"event":        "fury_unleash",
			"trigger_id":   playerID,
			"trigger_name": playerName,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🐉💥 龍怒爆發！%s 能量滿載！全場 HP -80%%！", playerName), "critical", "#FF3D00")

	time.Sleep(500 * time.Millisecond)
	hitCount := g.applyAOEDamage(GameWidth/2, GameHeight/2, 9999, 0.80)

	g.broadcast(protocol.Envelope{
		Type: "lucky_dragon_fury",
		Payload: map[string]interface{}{
			"event":        "fury_hit",
			"trigger_id":   playerID,
			"trigger_name": playerName,
			"hit_count":    hitCount,
		},
	})

	if hitCount >= 10 {
		m.mu.Lock()
		m.perfectBoost = &dragonFuryPerfectBoost{
			mult:      6.0,
			expiresAt: time.Now().Add(13 * time.Second),
		}
		m.mu.Unlock()
		g.broadcast(protocol.Envelope{
			Type: "lucky_dragon_fury",
			Payload: map[string]interface{}{
				"event":        "fury_perfect",
				"trigger_id":   playerID,
				"trigger_name": playerName,
				"boost_mult":   6.0,
				"boost_secs":   13,
			},
		})
		g.sendAnnounce(fmt.Sprintf("🐉✨ 完美龍怒！%s 命中 %d 個！全服 ×6.0 加成 13 秒！", playerName, hitCount), "critical", "#FFD700")
		time.AfterFunc(13*time.Second, func() {
			m.mu.Lock()
			m.perfectBoost = nil
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_dragon_fury",
				Payload: map[string]interface{}{
					"event":      "fury_perfect_end",
					"trigger_id": playerID,
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
}

func (m *luckyDragonFuryManager) tryLuckyDragonFuryFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(38 * time.Second)
	m.globalCD = now.Add(60 * time.Second)
	sess := &dragonFurySession{
		playerID:   p.ID,
		playerName: p.GetDisplayName(),
		energy:     0,
		expiresAt:  now.Add(15 * time.Second),
	}
	m.activeSession = sess
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_dragon_fury",
		Payload: map[string]interface{}{
			"event":        "energy_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     15,
			"max_energy":   100,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🐉 %s 啟動龍怒蓄能！15 秒內累積 100 能量觸發全場攻擊！", p.GetDisplayName()), "high", "#FF6D00")
	log.Printf("[LuckyDragonFury] %s 觸發龍怒能量魚", p.GetDisplayName())

	go func() {
		time.Sleep(15 * time.Second)
		m.mu.Lock()
		sess := m.activeSession
		if sess == nil || sess.settled {
			m.mu.Unlock()
			return
		}
		energy := sess.energy
		sess.settled = true
		m.activeSession = nil
		m.mu.Unlock()

		if energy < 100 {
			g.broadcast(protocol.Envelope{
				Type: "lucky_dragon_fury",
				Payload: map[string]interface{}{
					"event":        "energy_timeout",
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
					"energy":       energy,
				},
			})
		}
	}()
	return true
}
