// lucky_spacetime_rift_handler.go — T168 幸運時空裂縫魚
// 業界依據：Fishing Fortune「time warp + multiplier cascade」組合升級版
// 設計：擊破後時空裂縫 20 秒，每 4 秒隨機選 3 個目標瞬間擊破（獎勵 ×4.0）
//       裂縫期間玩家擊破 ≥12 個 → 時空完美：全服 ×7.5 加成 16 秒
//       個人冷卻 58 秒；全服冷卻 90 秒
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckySpacetimeRiftManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	sessions     map[string]*spacetimeRiftSession
	perfectBoost *spacetimeRiftPerfectBoost
}

type spacetimeRiftSession struct {
	playerID  string
	killCount int
	expiresAt time.Time
}

type spacetimeRiftPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckySpacetimeRiftManager() *luckySpacetimeRiftManager {
	return &luckySpacetimeRiftManager{
		personalCD: make(map[string]time.Time),
		sessions:   make(map[string]*spacetimeRiftSession),
	}
}

func isLuckySpacetimeRiftFish(defID string) bool {
	return defID == "T168"
}

func (m *luckySpacetimeRiftManager) getSpacetimeRiftPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckySpacetimeRiftManager) onKillDuringRift(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok || time.Now().After(sess.expiresAt) {
		return
	}
	sess.killCount++
}

func (m *luckySpacetimeRiftManager) tryLuckySpacetimeRiftFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(58 * time.Second)
	m.globalCD = now.Add(90 * time.Second)

	sess := &spacetimeRiftSession{
		playerID:  p.ID,
		killCount: 0,
		expiresAt: now.Add(20 * time.Second),
	}
	m.sessions[p.ID] = sess
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_spacetime_rift",
		Payload: map[string]interface{}{
			"event":        "rift_open",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     20,
			"wave_count":   5,
		},
	})
	g.sendAnnounce(fmt.Sprintf("⏳🌀 時空裂縫！%s 撕裂時空！每 4 秒瞬間擊破 3 個目標！", p.GetDisplayName()), "high", "#1565C0")
	log.Printf("[LuckySpacetimeRift] %s 觸發時空裂縫魚", p.GetDisplayName())

	// 5 波，每 4 秒一波
	go func() {
		for wave := 1; wave <= 5; wave++ {
			time.Sleep(4 * time.Second)

			m.mu.Lock()
			sess, ok := m.sessions[p.ID]
			if !ok || time.Now().After(sess.expiresAt) {
				m.mu.Unlock()
				break
			}
			m.mu.Unlock()

			// 瞬間擊破 3 個目標，獎勵 ×4.0
			killCount := g.applyRiftInstantKill(p, 3, 4.0)

			g.broadcast(protocol.Envelope{
				Type: "lucky_spacetime_rift",
				Payload: map[string]interface{}{
					"event":      "rift_wave",
					"wave":       wave,
					"kill_count": killCount,
					"trigger_id": p.ID,
				},
			})
		}

		// 結算
		m.mu.Lock()
		finalKills := 0
		if s, ok := m.sessions[p.ID]; ok {
			finalKills = s.killCount
			delete(m.sessions, p.ID)
		}
		m.mu.Unlock()

		if finalKills >= 12 {
			m.mu.Lock()
			m.perfectBoost = &spacetimeRiftPerfectBoost{
				mult:      7.5,
				expiresAt: time.Now().Add(16 * time.Second),
			}
			m.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: "lucky_spacetime_rift",
				Payload: map[string]interface{}{
					"event":        "rift_perfect",
					"kill_count":   finalKills,
					"boost_mult":   7.5,
					"boost_secs":   16,
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
				},
			})
			g.sendAnnounce(fmt.Sprintf("⏳💫 時空完美！%s 裂縫期間擊破 %d 個！全服 ×7.5 加成 16 秒！", p.GetDisplayName(), finalKills), "critical", "#0D47A1")
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_spacetime_rift",
				Payload: map[string]interface{}{
					"event":      "rift_end",
					"kill_count": finalKills,
					"trigger_id": p.ID,
				},
			})
		}
	}()

	return true
}

// applyRiftInstantKill — 瞬間擊破最多 n 個目標，獎勵 rewardMult 倍
func (g *Game) applyRiftInstantKill(p *Player, n int, rewardMult float64) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	var candidates []*Target
	for _, t := range g.targets {
		if t.Def.Type != "boss" {
			candidates = append(candidates, t)
		}
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	if len(candidates) < n {
		n = len(candidates)
	}

	bet := p.GetBetDef()
	killCount := 0
	var toDelete []string
	for i := 0; i < n; i++ {
		t := candidates[i]
		reward := int(float64(bet.BetCost) * rewardMult)
		p.AddCoins(reward)
		toDelete = append(toDelete, t.InstanceID)
		g.hub.Broadcast(protocol.MsgTargetKill, protocol.TargetKillPayload{
			InstanceID: t.InstanceID,
			DefID:      t.Def.ID,
			Multiplier: t.Multiplier * rewardMult,
			Reward:     reward,
			LaborGain:  0,
			KillerID:   p.ID,
		})
		killCount++
	}
	for _, id := range toDelete {
		delete(g.targets, id)
	}
	return killCount
}
