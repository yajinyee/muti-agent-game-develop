// lucky_immortal_boss_handler.go — T155 幸運不死 BOSS 魚
// 業界依據：Royal Fishing「Immortal Boss mechanic — consecutive wins 50X-150X until they leave the screen」
// 設計：擊破後召喚不死 BOSS（5 條命，每次擊破倍率 +0.5x），18 秒內耗盡 5 條命 → 完美不死全服 ×5.0 加成 12 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyImmortalBossManager struct {
	mu             sync.Mutex
	personalCD     map[string]time.Time
	globalCD       time.Time
	activeSessions map[string]*immortalBossSession
	perfectBoost   *immortalBossPerfectBoost
}

type immortalBossPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type immortalBossSession struct {
	playerID   string
	playerName string
	livesLeft  int
	currentMult float64
	totalReward int
	expiresAt  time.Time
	settled    bool
}

func newLuckyImmortalBossManager() *luckyImmortalBossManager {
	return &luckyImmortalBossManager{
		personalCD:     make(map[string]time.Time),
		activeSessions: make(map[string]*immortalBossSession),
	}
}

func isLuckyImmortalBossFish(defID string) bool {
	return defID == "T155"
}

func (m *luckyImmortalBossManager) getImmortalBossPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyImmortalBossManager) notifyImmortalBossKill(g *Game, p *Player) {
	m.mu.Lock()
	sess, ok := m.activeSessions[p.ID]
	if !ok || sess.settled || time.Now().After(sess.expiresAt) {
		m.mu.Unlock()
		return
	}
	sess.livesLeft--
	sess.currentMult += 0.5
	livesLeft := sess.livesLeft
	currentMult := sess.currentMult
	reward := int(float64(p.GetBetDef().BetCost) * currentMult)
	p.AddCoins(reward)
	sess.totalReward += reward
	totalReward := sess.totalReward
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_immortal_boss",
		Payload: map[string]interface{}{
			"event":        "immortal_kill",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"lives_left":   livesLeft,
			"current_mult": currentMult,
			"reward":       reward,
		},
	})

	if livesLeft <= 0 {
		m.mu.Lock()
		sess.settled = true
		delete(m.activeSessions, p.ID)
		m.perfectBoost = &immortalBossPerfectBoost{
			mult:      5.0,
			expiresAt: time.Now().Add(12 * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_immortal_boss",
			Payload: map[string]interface{}{
				"event":        "immortal_perfect",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"total_reward": totalReward,
				"boost_mult":   5.0,
				"boost_secs":   12,
			},
		})
		g.sendAnnounce(fmt.Sprintf("💀✨ 完美不死！%s 耗盡 5 條命！全服 ×5.0 加成 12 秒！", p.GetDisplayName()), "critical", "#FFD700")
		time.AfterFunc(12*time.Second, func() {
			m.mu.Lock()
			m.perfectBoost = nil
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_immortal_boss",
				Payload: map[string]interface{}{
					"event":      "immortal_perfect_end",
					"trigger_id": p.ID,
				},
			})
		})
	}
}

func (m *luckyImmortalBossManager) tryLuckyImmortalBossFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(32 * time.Second)
	m.globalCD = now.Add(50 * time.Second)
	sess := &immortalBossSession{
		playerID:    p.ID,
		playerName:  p.GetDisplayName(),
		livesLeft:   5,
		currentMult: 2.0,
		expiresAt:   now.Add(18 * time.Second),
	}
	m.activeSessions[p.ID] = sess
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_immortal_boss",
		Payload: map[string]interface{}{
			"event":        "immortal_spawn",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"lives":        5,
			"init_mult":    2.0,
			"duration":     18,
		},
	})
	g.sendAnnounce(fmt.Sprintf("💀 %s 召喚不死 BOSS！5 條命，倍率遞增！18 秒！", p.GetDisplayName()), "high", "#B71C1C")
	log.Printf("[LuckyImmortalBoss] %s 觸發不死 BOSS", p.GetDisplayName())

	go func() {
		time.Sleep(18 * time.Second)
		m.mu.Lock()
		sess, ok := m.activeSessions[p.ID]
		if !ok || sess.settled {
			m.mu.Unlock()
			return
		}
		sess.settled = true
		livesLeft := sess.livesLeft
		totalReward := sess.totalReward
		delete(m.activeSessions, p.ID)
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_immortal_boss",
			Payload: map[string]interface{}{
				"event":        "immortal_timeout",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"lives_left":   livesLeft,
				"total_reward": totalReward,
			},
		})
	}()
	return true
}
