// lucky_elemental_fusion_handler.go — T163 幸運元素融合魚
// 業界依據：Elemental system — 三元素（火/冰/雷）融合，全部觸發 → 元素爆發
// 設計：擊破後觸發元素融合模式，25 秒內依序觸發火/冰/雷三元素
//       每個元素觸發：全場 HP -25%；三元素全部觸發 → 元素爆發全服 ×6.5 加成 14 秒
//       個人冷卻 55 秒；全服冷卻 80 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyElementalFusionManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	sessions   map[string]*elementalFusionSession
	perfectBoost *elementalFusionPerfectBoost
}

type elementalFusionSession struct {
	playerID    string
	fireTriggered  bool
	iceTriggered   bool
	thunderTriggered bool
	expiresAt   time.Time
}

type elementalFusionPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyElementalFusionManager() *luckyElementalFusionManager {
	return &luckyElementalFusionManager{
		personalCD: make(map[string]time.Time),
		sessions:   make(map[string]*elementalFusionSession),
	}
}

func isLuckyElementalFusionFish(defID string) bool {
	return defID == "T163"
}

func (m *luckyElementalFusionManager) getElementalFusionPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyElementalFusionManager) tryLuckyElementalFusionFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(55 * time.Second)
	m.globalCD = now.Add(80 * time.Second)
	m.sessions[p.ID] = &elementalFusionSession{
		playerID:  p.ID,
		expiresAt: now.Add(25 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyElementalFusion] Player %s triggered elemental fusion (25s)", p.ID)

	g.hub.Broadcast(protocol.MsgLuckyElementalFusion, map[string]interface{}{
		"event":     "start",
		"player_id": p.ID,
		"duration":  25,
		"elements":  []string{"fire", "ice", "thunder"},
	})

	// 依序觸發三元素
	go func() {
		elements := []struct {
			name     string
			delay    time.Duration
			field    *bool
		}{
			{"fire", 5 * time.Second, nil},
			{"ice", 12 * time.Second, nil},
			{"thunder", 20 * time.Second, nil},
		}

		for i, elem := range elements {
			time.Sleep(elem.delay)
			m.mu.Lock()
			sess, ok := m.sessions[p.ID]
			if !ok || time.Now().After(sess.expiresAt) {
				m.mu.Unlock()
				break
			}
			switch i {
			case 0:
				sess.fireTriggered = true
			case 1:
				sess.iceTriggered = true
			case 2:
				sess.thunderTriggered = true
			}
			m.mu.Unlock()

			// 元素攻擊：全場 HP -25%
			g.mu.Lock()
			hitCount := 0
			for _, t := range g.targets {
				if t.HP > 0 {
					dmg := int(float64(t.MaxHP) * 0.25)
					t.HP -= dmg
					if t.HP < 0 {
						t.HP = 0
					}
					hitCount++
				}
			}
			g.mu.Unlock()

			g.hub.Broadcast(protocol.MsgLuckyElementalFusion, map[string]interface{}{
				"event":     "element_trigger",
				"player_id": p.ID,
				"element":   elem.name,
				"hit_count": hitCount,
				"dmg_pct":   0.25,
			})
		}

		// 25 秒後結算
		time.Sleep(5 * time.Second)
		m.mu.Lock()
		sess, ok := m.sessions[p.ID]
		allTriggered := false
		if ok {
			allTriggered = sess.fireTriggered && sess.iceTriggered && sess.thunderTriggered
			delete(m.sessions, p.ID)
		}
		m.mu.Unlock()

		if allTriggered {
			m.mu.Lock()
			m.perfectBoost = &elementalFusionPerfectBoost{
				mult:      6.5,
				expiresAt: time.Now().Add(14 * time.Second),
			}
			m.mu.Unlock()
			g.hub.Broadcast(protocol.MsgAnnounce, map[string]interface{}{
				"key":     "elemental_fusion_perfect",
				"message": fmt.Sprintf("⚡🔥❄️ 元素爆發！%s 三元素融合完成！全服 ×6.5 加成 14 秒！", p.ID),
				"mult":    6.5,
				"duration": 14,
			})
		}

		g.hub.Broadcast(protocol.MsgLuckyElementalFusion, map[string]interface{}{
			"event":        "end",
			"player_id":    p.ID,
			"all_triggered": allTriggered,
			"perfect":      allTriggered,
		})
	}()

	return true
}
