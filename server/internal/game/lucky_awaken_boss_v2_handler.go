// lucky_awaken_boss_v2_handler.go — T159 幸運覺醒 BOSS 魚 v2
// 業界依據：Royal Fishing「Awaken Boss — Power Up attacks multiply by 6x-10x for devastating rewards」升級版
// 設計：覺醒後 8 次 Power Up（每次 8x-15x 隨機），全部命中 → 完美覺醒：全服 ×7.0 加成 15 秒
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyAwakenBossV2Manager struct {
	mu             sync.Mutex
	personalCD     map[string]time.Time
	globalCD       time.Time
	activeSessions map[string]*awakenBossV2Session
	perfectBoost   *awakenBossV2PerfectBoost
}

type awakenBossV2PerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type awakenBossV2Session struct {
	playerID    string
	playerName  string
	shotsLeft   int
	hitCount    int
	totalReward int
	expiresAt   time.Time
	settled     bool
}

func newLuckyAwakenBossV2Manager() *luckyAwakenBossV2Manager {
	return &luckyAwakenBossV2Manager{
		personalCD:     make(map[string]time.Time),
		activeSessions: make(map[string]*awakenBossV2Session),
	}
}

func isLuckyAwakenBossV2Fish(defID string) bool {
	return defID == "T159"
}

func (m *luckyAwakenBossV2Manager) getAwakenBossV2PerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyAwakenBossV2Manager) notifyAwakenBossV2Kill(g *Game, p *Player) {
	m.mu.Lock()
	sess, ok := m.activeSessions[p.ID]
	if !ok || sess.settled || time.Now().After(sess.expiresAt) || sess.shotsLeft <= 0 {
		m.mu.Unlock()
		return
	}
	// 每次命中觸發 Power Up（8x-15x 隨機）
	powerUpMult := 8.0 + rand.Float64()*7.0 // 8.0 ~ 15.0
	reward := int(float64(p.GetBetDef().BetCost) * powerUpMult)
	p.AddCoins(reward)
	sess.shotsLeft--
	sess.hitCount++
	sess.totalReward += reward
	shotsLeft := sess.shotsLeft
	hitCount := sess.hitCount
	totalReward := sess.totalReward
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_awaken_boss_v2",
		Payload: map[string]interface{}{
			"event":         "power_up",
			"trigger_id":    p.ID,
			"trigger_name":  p.GetDisplayName(),
			"power_up_mult": powerUpMult,
			"shots_left":    shotsLeft,
			"hit_count":     hitCount,
			"reward":        reward,
		},
	})

	if shotsLeft <= 0 {
		// 8 次全部命中 → 完美覺醒
		m.mu.Lock()
		sess.settled = true
		delete(m.activeSessions, p.ID)
		m.perfectBoost = &awakenBossV2PerfectBoost{
			mult:      7.0,
			expiresAt: time.Now().Add(15 * time.Second),
		}
		m.mu.Unlock()
		g.broadcast(protocol.Envelope{
			Type: "lucky_awaken_boss_v2",
			Payload: map[string]interface{}{
				"event":        "awaken_perfect",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"hit_count":    hitCount,
				"total_reward": totalReward,
				"boost_mult":   7.0,
				"boost_secs":   15,
			},
		})
		g.sendAnnounce(fmt.Sprintf("⚡✨ 完美覺醒！%s 8 次全命中！全服 ×7.0 加成 15 秒！", p.GetDisplayName()), "critical", "#FFD700")
		time.AfterFunc(15*time.Second, func() {
			m.mu.Lock()
			m.perfectBoost = nil
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_awaken_boss_v2",
				Payload: map[string]interface{}{
					"event":      "awaken_perfect_end",
					"trigger_id": p.ID,
				},
			})
		})
	}
}

func (m *luckyAwakenBossV2Manager) tryLuckyAwakenBossV2Fish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(42 * time.Second)
	m.globalCD = now.Add(68 * time.Second)
	sess := &awakenBossV2Session{
		playerID:   p.ID,
		playerName: p.GetDisplayName(),
		shotsLeft:  8,
		expiresAt:  now.Add(25 * time.Second),
	}
	m.activeSessions[p.ID] = sess
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_awaken_boss_v2",
		Payload: map[string]interface{}{
			"event":        "awaken_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"shots":        8,
			"duration":     25,
			"min_mult":     8.0,
			"max_mult":     15.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("⚡ %s 覺醒 BOSS！8 次 Power Up（8x-15x）！25 秒！", p.GetDisplayName()), "high", "#F57F17")
	log.Printf("[LuckyAwakenBossV2] %s 觸發覺醒 BOSS v2", p.GetDisplayName())

	go func() {
		time.Sleep(25 * time.Second)
		m.mu.Lock()
		sess, ok := m.activeSessions[p.ID]
		if !ok || sess.settled {
			m.mu.Unlock()
			return
		}
		shotsLeft := sess.shotsLeft
		hitCount := sess.hitCount
		totalReward := sess.totalReward
		sess.settled = true
		delete(m.activeSessions, p.ID)
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_awaken_boss_v2",
			Payload: map[string]interface{}{
				"event":        "awaken_timeout",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"shots_left":   shotsLeft,
				"hit_count":    hitCount,
				"total_reward": totalReward,
			},
		})
	}()
	return true
}
