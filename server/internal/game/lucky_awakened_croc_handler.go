// lucky_awakened_croc_handler.go — T151 幸運覺醒鱷魚
// 業界依據：Jili「Giant Crocodile awakens to hunt fish on the fish farm to accumulate big prizes」
// 設計：擊破後覺醒鱷魚在場上自動獵魚 20 秒，每次獵魚 ×3.0 獎勵，獵魚 ≥8 → 完美覺醒全服 ×3.5 加成 9 秒
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyAwakenedCrocManager struct {
	mu             sync.Mutex
	personalCD     map[string]time.Time
	globalCD       time.Time
	activeSessions map[string]*crocSession
	perfectBoost   *crocPerfectBoost
}

type crocPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type crocSession struct {
	playerID   string
	playerName string
	huntCount  int
	expiresAt  time.Time
	settled    bool
}

func newLuckyAwakenedCrocManager() *luckyAwakenedCrocManager {
	return &luckyAwakenedCrocManager{
		personalCD:     make(map[string]time.Time),
		activeSessions: make(map[string]*crocSession),
	}
}

func isLuckyAwakenedCrocFish(defID string) bool {
	return defID == "T151"
}

func (m *luckyAwakenedCrocManager) getCrocPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyAwakenedCrocManager) tryLuckyAwakenedCrocFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(28 * time.Second)
	m.globalCD = now.Add(45 * time.Second)
	sess := &crocSession{
		playerID:   p.ID,
		playerName: p.GetDisplayName(),
		expiresAt:  now.Add(20 * time.Second),
	}
	m.activeSessions[p.ID] = sess
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_awakened_croc",
		Payload: map[string]interface{}{
			"event":        "croc_awaken",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     20,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🐊 %s 觸發覺醒鱷魚！鱷魚開始獵魚 20 秒！", p.GetDisplayName()), "high", "#00C853")
	log.Printf("[LuckyAwakenedCroc] %s 觸發覺醒鱷魚", p.GetDisplayName())

	go m.runCrocHunt(g, p, sess)
	return true
}

func (m *luckyAwakenedCrocManager) runCrocHunt(g *Game, p *Player, sess *crocSession) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	deadline := time.NewTimer(20 * time.Second)
	defer deadline.Stop()

	for {
		select {
		case <-ticker.C:
			g.mu.Lock()
			// 隨機選一個非 BOSS 目標獵殺
			var targetIDs []string
			for id, t := range g.targets {
				if t.HP > 0 {
					targetIDs = append(targetIDs, id)
				}
			}
			if len(targetIDs) > 0 {
				idx := rand.Intn(len(targetIDs))
				tid := targetIDs[idx]
				t := g.targets[tid]
				// 鱷魚獵殺：HP -60%
				dmg := int(float64(t.MaxHP) * 0.60)
				t.HP -= dmg
				if t.HP <= 0 {
					t.HP = 0
					delete(g.targets, tid)
					reward := int(float64(p.GetBetDef().BetCost) * t.Def.Multiplier * 3.0)
					p.AddCoins(reward)
					m.mu.Lock()
					sess.huntCount++
					m.mu.Unlock()
					g.broadcast(protocol.Envelope{
						Type: "lucky_awakened_croc",
						Payload: map[string]interface{}{
							"event":      "croc_hunt",
							"trigger_id": p.ID,
							"hunt_count": sess.huntCount,
							"reward":     reward,
							"target_id":  tid,
						},
					})
					g.broadcast(protocol.Envelope{
						Type: protocol.MsgTargetKill,
						Payload: protocol.TargetKillPayload{
							InstanceID: t.InstanceID,
							DefID:      t.Def.ID,
							Multiplier: t.Def.Multiplier * 3.0,
							Reward:     reward,
							LaborGain:  0,
							KillerID:   p.ID,
						},
					})
				} else {
					g.broadcast(protocol.Envelope{
						Type: protocol.MsgTargetUpdate,
						Payload: protocol.TargetUpdatePayload{
							InstanceID: t.InstanceID,
							HP:         t.HP,
							MaxHP:      t.MaxHP,
						},
					})
				}
			}
			g.mu.Unlock()

		case <-deadline.C:
			m.mu.Lock()
			sess.settled = true
			huntCount := sess.huntCount
			delete(m.activeSessions, p.ID)
			m.mu.Unlock()

			if huntCount >= 8 {
				m.mu.Lock()
				m.perfectBoost = &crocPerfectBoost{
					mult:      3.5,
					expiresAt: time.Now().Add(9 * time.Second),
				}
				m.mu.Unlock()
				g.broadcast(protocol.Envelope{
					Type: "lucky_awakened_croc",
					Payload: map[string]interface{}{
						"event":        "croc_perfect",
						"trigger_id":   p.ID,
						"trigger_name": p.GetDisplayName(),
						"hunt_count":   huntCount,
						"boost_mult":   3.5,
						"boost_secs":   9,
					},
				})
				g.sendAnnounce(fmt.Sprintf("🐊✨ 完美覺醒！%s 獵魚 %d 條！全服 ×3.5 加成 9 秒！", p.GetDisplayName(), huntCount), "high", "#FFD700")
				time.AfterFunc(9*time.Second, func() {
					m.mu.Lock()
					m.perfectBoost = nil
					m.mu.Unlock()
					g.broadcast(protocol.Envelope{
						Type: "lucky_awakened_croc",
						Payload: map[string]interface{}{
							"event":      "croc_perfect_end",
							"trigger_id": p.ID,
						},
					})
				})
			} else {
				g.broadcast(protocol.Envelope{
					Type: "lucky_awakened_croc",
					Payload: map[string]interface{}{
						"event":        "croc_end",
						"trigger_id":   p.ID,
						"trigger_name": p.GetDisplayName(),
						"hunt_count":   huntCount,
					},
				})
			}
			return
		}
	}
}
