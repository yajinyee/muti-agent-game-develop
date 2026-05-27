// lucky_fisher_wild_handler.go — T183 幸運漁夫野生魚
// 業界依據：Big Bass Splash 1000（Fisherman Wild + Fish Cash mechanic）
// 設計：擊破後標記 3 個 Wild 目標（HP -50%，擊破獎勵 ×5.0）
//       30 秒內全部擊破 → 全服 ×17.0 加成 35 秒
//       個人冷卻 98 秒；全服冷卻 150 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/data"
	"chiikawa-game/internal/protocol"
)

type luckyFisherWildManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	activeSession *fisherWildSession
	perfectBoost  *fisherWildPerfectBoost
}

type fisherWildSession struct {
	triggerID    string
	triggerName  string
	wildTargets  []string // instance IDs of marked wild targets
	killedWilds  int
	expiresAt    time.Time
}

type fisherWildPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyFisherWildManager() *luckyFisherWildManager {
	return &luckyFisherWildManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyFisherWildFish(defID string) bool {
	return defID == "T183"
}

func (m *luckyFisherWildManager) getFisherWildMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyFisherWildManager) isWildTarget(instanceID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil || time.Now().After(m.activeSession.expiresAt) {
		return false
	}
	for _, id := range m.activeSession.wildTargets {
		if id == instanceID {
			return true
		}
	}
	return false
}

func (m *luckyFisherWildManager) onWildKilled(g *Game, instanceID string) {
	m.mu.Lock()
	if m.activeSession == nil || time.Now().After(m.activeSession.expiresAt) {
		m.mu.Unlock()
		return
	}
	// 確認是 Wild 目標
	isWild := false
	for _, id := range m.activeSession.wildTargets {
		if id == instanceID {
			isWild = true
			break
		}
	}
	if !isWild {
		m.mu.Unlock()
		return
	}
	m.activeSession.killedWilds++
	killed := m.activeSession.killedWilds
	total := len(m.activeSession.wildTargets)
	triggerID := m.activeSession.triggerID
	triggerName := m.activeSession.triggerName
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_fisher_wild",
		Payload: map[string]interface{}{
			"event":        "wild_killed",
			"killed":       killed,
			"total":        total,
			"instance_id":  instanceID,
			"trigger_id":   triggerID,
			"trigger_name": triggerName,
		},
	})

	// 全部 Wild 目標擊破 → 完美漁夫
	if killed >= total {
		boostMult := 17.0
		boostSecs := 35
		m.mu.Lock()
		m.perfectBoost = &fisherWildPerfectBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.activeSession = nil
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_fisher_wild",
			Payload: map[string]interface{}{
				"event":        "fisher_wild_perfect",
				"killed":       killed,
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
				"trigger_id":   triggerID,
				"trigger_name": triggerName,
			},
		})
		g.sendAnnounce(fmt.Sprintf("🎣🏆 完美漁夫！%s 擊破全部 Wild 目標！全服 ×%.1f 加成 %d 秒！", triggerName, boostMult, boostSecs), "critical", "#1565C0")
	}
}

func (m *luckyFisherWildManager) tryLuckyFisherWildFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(98 * time.Second)
	m.globalCD = now.Add(150 * time.Second)
	m.mu.Unlock()

	// 標記場上 3 個隨機目標為 Wild
	wildIDs := []string{}
	g.mu.Lock()
	for id, t := range g.targets {
		if t.HP > 0 && t.Def.Type != data.TypeBoss && len(wildIDs) < 3 {
			// Wild 目標 HP -50%
			damage := t.HP / 2
			t.HP -= damage
			if t.HP < 1 {
				t.HP = 1
			}
			wildIDs = append(wildIDs, id)
		}
	}
	g.mu.Unlock()

	if len(wildIDs) == 0 {
		return false
	}

	m.mu.Lock()
	m.activeSession = &fisherWildSession{
		triggerID:   p.ID,
		triggerName: p.GetDisplayName(),
		wildTargets: wildIDs,
		killedWilds: 0,
		expiresAt:   time.Now().Add(30 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_fisher_wild",
		Payload: map[string]interface{}{
			"event":        "fisher_wild_start",
			"wild_targets": wildIDs,
			"wild_count":   len(wildIDs),
			"wild_mult":    5.0,
			"duration":     30,
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("🎣✨ 漁夫野生！%s 標記 %d 個 Wild 目標！擊破獎勵 ×5.0！", p.GetDisplayName(), len(wildIDs)), "special", "#1976D2")
	log.Printf("[LuckyFisherWild] %s 觸發漁夫野生魚，標記 %d 個目標", p.GetDisplayName(), len(wildIDs))

	// 30 秒後超時結算
	go func() {
		time.Sleep(30 * time.Second)
		m.mu.Lock()
		if m.activeSession != nil && m.activeSession.triggerID == p.ID {
			killed := m.activeSession.killedWilds
			total := len(m.activeSession.wildTargets)
			m.activeSession = nil
			m.mu.Unlock()
			if killed < total {
				g.broadcast(protocol.Envelope{
					Type: "lucky_fisher_wild",
					Payload: map[string]interface{}{
						"event":        "fisher_wild_timeout",
						"killed":       killed,
						"total":        total,
						"trigger_id":   p.ID,
						"trigger_name": p.GetDisplayName(),
					},
				})
			}
		} else {
			m.mu.Unlock()
		}
	}()
	return true
}
