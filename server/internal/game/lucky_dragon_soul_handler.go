// lucky_dragon_soul_handler.go — T167 幸運龍魂融合魚
// 業界依據：Royal Fishing「Dragon Wrath energy accumulation」升級版
// 設計：擊破後龍魂融合 30 秒，每次擊破吸收龍魂（最高 50 魂）
//       50 魂 → 龍魂爆發全場 HP -90%，全服 ×9.0 加成 18 秒
//       個人冷卻 60 秒；全服冷卻 95 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyDragonSoulManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	sessions     map[string]*dragonSoulSession
	perfectBoost *dragonSoulPerfectBoost
}

type dragonSoulSession struct {
	playerID  string
	soulCount int
	expiresAt time.Time
}

type dragonSoulPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyDragonSoulManager() *luckyDragonSoulManager {
	return &luckyDragonSoulManager{
		personalCD: make(map[string]time.Time),
		sessions:   make(map[string]*dragonSoulSession),
	}
}

func isLuckyDragonSoulFish(defID string) bool {
	return defID == "T167"
}

func (m *luckyDragonSoulManager) getDragonSoulPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyDragonSoulManager) onKillDuringDragonSoul(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok || time.Now().After(sess.expiresAt) {
		return
	}
	if sess.soulCount < 50 {
		sess.soulCount++
	}
}

func (m *luckyDragonSoulManager) tryLuckyDragonSoulFish(g *Game, p *Player) bool {
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
	m.globalCD = now.Add(95 * time.Second)

	sess := &dragonSoulSession{
		playerID:  p.ID,
		soulCount: 0,
		expiresAt: now.Add(30 * time.Second),
	}
	m.sessions[p.ID] = sess
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_dragon_soul",
		Payload: map[string]interface{}{
			"event":        "soul_fusion_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     30,
			"max_souls":    50,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🐉🔥 龍魂融合！%s 開始吸收龍魂！30 秒內集滿 50 魂！", p.GetDisplayName()), "high", "#D32F2F")
	log.Printf("[LuckyDragonSoul] %s 觸發龍魂融合魚", p.GetDisplayName())

	go func() {
		time.Sleep(30 * time.Second)

		m.mu.Lock()
		finalSouls := 0
		if s, ok := m.sessions[p.ID]; ok {
			finalSouls = s.soulCount
			delete(m.sessions, p.ID)
		}
		m.mu.Unlock()

		// 龍魂爆發：全場 HP -90%
		hitCount := g.applyAOEDamage(640, 360, 99999, 0.9)

		g.broadcast(protocol.Envelope{
			Type: "lucky_dragon_soul",
			Payload: map[string]interface{}{
				"event":       "soul_burst",
				"soul_count":  finalSouls,
				"hit_count":   hitCount,
				"trigger_id":  p.ID,
			},
		})

		// 判定是否完美（50 魂）
		if finalSouls >= 50 {
			m.mu.Lock()
			m.perfectBoost = &dragonSoulPerfectBoost{
				mult:      9.0,
				expiresAt: time.Now().Add(18 * time.Second),
			}
			m.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: "lucky_dragon_soul",
				Payload: map[string]interface{}{
					"event":        "soul_perfect",
					"boost_mult":   9.0,
					"boost_secs":   18,
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
				},
			})
			g.sendAnnounce(fmt.Sprintf("🐉💥 龍魂完美！%s 集滿 50 魂！全場 HP -90%%！全服 ×9.0 加成 18 秒！", p.GetDisplayName()), "critical", "#FF6F00")
		} else {
			g.sendAnnounce(fmt.Sprintf("🐉 龍魂爆發！%s 吸收 %d 魂！全場 HP -90%%！", p.GetDisplayName(), finalSouls), "high", "#D32F2F")
		}
	}()

	return true
}
