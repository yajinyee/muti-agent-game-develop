// Package game — T134 幸運賞金獵人魚 handler
// server-event-agent 負責維護
// 業界依據：Fishing Frenzy Chapter 3「Guild Wars + Boss Fish + quality values」
//           「賞金系統 — 隨機標記 3 條目標為賞金目標，擊破賞金目標獲得額外獎勵」
// 設計：擊破後隨機標記場上 3 條目標為「賞金目標」（HP -20% 弱化 + 金色標記）；
//       每擊破一個賞金目標 → 觸發玩家獲得 ×3.0 獎勵；
//       30 秒內擊破全部 3 個賞金目標 → 「完美賞金」：全服 ×3.5 加成 8 秒；
//       個人冷卻 26 秒；全服冷卻 42 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyBountyHunterManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	activeSession *bountyHunterSession
	perfectBoost  *bountyPerfectBoost
}

type bountyPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type bountyHunterSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	bountyTargets     map[string]bool // instanceID -> killed
	killCount         int
	expiresAt         time.Time
	settled           bool
}

func newLuckyBountyHunterManager() *luckyBountyHunterManager {
	return &luckyBountyHunterManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyBountyHunterFish(defID string) bool {
	return defID == "T134"
}

func (m *luckyBountyHunterManager) getBountyPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyBountyHunterManager) isBountyTarget(instanceID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil || m.activeSession.settled {
		return false
	}
	_, ok := m.activeSession.bountyTargets[instanceID]
	return ok
}

func (g *Game) tryLuckyBountyHunterFish(playerID, playerName string) {
	m := g.luckyBountyHunter
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) {
		m.mu.Unlock()
		return
	}
	if cd, ok := m.personalCD[playerID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	if m.activeSession != nil && !m.activeSession.settled {
		m.mu.Unlock()
		return
	}
	m.personalCD[playerID] = now.Add(26 * time.Second)
	m.globalCD = now.Add(42 * time.Second)
	m.mu.Unlock()

	// 選取 3 個賞金目標
	g.mu.Lock()
	var candidates []string
	for id, t := range g.targets {
		if t.HP > 0 {
			candidates = append(candidates, id)
		}
	}
	// 隨機選 3 個
	bountyIDs := make(map[string]bool)
	maxBounty := 3
	if len(candidates) < maxBounty {
		maxBounty = len(candidates)
	}
	// Fisher-Yates shuffle 取前 maxBounty 個
	for i := len(candidates) - 1; i > 0; i-- {
		j := int(time.Now().UnixNano()) % (i + 1)
		if j < 0 {
			j = -j
		}
		candidates[i], candidates[j] = candidates[j], candidates[i]
	}
	for i := 0; i < maxBounty; i++ {
		bountyIDs[candidates[i]] = false
		// 弱化賞金目標 HP -20%
		if t, ok := g.targets[candidates[i]]; ok {
			dmg := int(float64(t.MaxHP) * 0.20)
			t.HP -= dmg
			if t.HP < 0 {
				t.HP = 1
			}
		}
	}
	g.mu.Unlock()

	if len(bountyIDs) == 0 {
		return
	}

	m.mu.Lock()
	session := &bountyHunterSession{
		triggerPlayerID:   playerID,
		triggerPlayerName: playerName,
		bountyTargets:     bountyIDs,
		killCount:         0,
		expiresAt:         time.Now().Add(30 * time.Second),
		settled:           false,
	}
	m.activeSession = session
	m.mu.Unlock()

	log.Printf("[LuckyBountyHunter] Triggered by %s, %d bounty targets", playerName, len(bountyIDs))

	// 廣播賞金目標列表
	var bountyList []string
	for id := range bountyIDs {
		bountyList = append(bountyList, id)
	}
	g.hub.Broadcast(protocol.MsgLuckyBountyHunter, protocol.LuckyBountyHunterPayload{
		Event:         "bounty_start",
		PlayerID:      playerID,
		PlayerName:    playerName,
		BountyTargets: bountyList,
		TotalBounty:   len(bountyIDs),
		Duration:      30.0,
	})

	// 30 秒後超時結算
	go func() {
		time.Sleep(30 * time.Second)
		m.mu.Lock()
		if session.settled {
			m.mu.Unlock()
			return
		}
		session.settled = true
		killCount := session.killCount
		m.mu.Unlock()

		g.hub.Broadcast(protocol.MsgLuckyBountyHunter, protocol.LuckyBountyHunterPayload{
			Event:      "bounty_timeout",
			PlayerID:   playerID,
			PlayerName: playerName,
			KillCount:  killCount,
		})
	}()
}

func (g *Game) notifyBountyKill(instanceID, killerID, killerName string) {
	m := g.luckyBountyHunter
	m.mu.Lock()
	if m.activeSession == nil || m.activeSession.settled {
		m.mu.Unlock()
		return
	}
	session := m.activeSession
	if _, ok := session.bountyTargets[instanceID]; !ok {
		m.mu.Unlock()
		return
	}
	session.bountyTargets[instanceID] = true
	session.killCount++
	killCount := session.killCount
	totalBounty := len(session.bountyTargets)
	playerID := session.triggerPlayerID
	playerName := session.triggerPlayerName
	m.mu.Unlock()

	log.Printf("[LuckyBountyHunter] Bounty kill! %s killed %s (%d/%d)", killerName, instanceID, killCount, totalBounty)

	g.hub.Broadcast(protocol.MsgLuckyBountyHunter, protocol.LuckyBountyHunterPayload{
		Event:       "bounty_kill",
		PlayerID:    playerID,
		PlayerName:  playerName,
		KillerID:    killerID,
		KillerName:  killerName,
		KillCount:   killCount,
		TotalBounty: totalBounty,
		BoostMult:   3.0,
	})

	// 全部擊破 → 完美賞金
	if killCount >= totalBounty {
		m.mu.Lock()
		session.settled = true
		m.mu.Unlock()
		g.doBountyPerfect(playerID, playerName, killCount)
	}
}

func (g *Game) doBountyPerfect(playerID, playerName string, killCount int) {
	m := g.luckyBountyHunter
	m.mu.Lock()
	m.perfectBoost = &bountyPerfectBoost{
		mult:      3.5,
		expiresAt: time.Now().Add(8 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyBountyHunter] Perfect! %s killed all %d bounties → global ×3.5 for 8s", playerName, killCount)

	g.hub.Broadcast(protocol.MsgLuckyBountyHunter, protocol.LuckyBountyHunterPayload{
		Event:      "bounty_perfect",
		PlayerID:   playerID,
		PlayerName: playerName,
		KillCount:  killCount,
		BoostMult:  3.5,
		BoostSec:   8,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  fmt.Sprintf("🎯 完美賞金！%s 獵殺全部 %d 個賞金目標！全服 ×3.5 加成 8 秒！", playerName, killCount),
		Priority: "high",
		Color:    "#FF6B35",
	})

	go func() {
		time.Sleep(8 * time.Second)
		m.mu.Lock()
		m.perfectBoost = nil
		m.mu.Unlock()
		g.hub.Broadcast(protocol.MsgLuckyBountyHunter, protocol.LuckyBountyHunterPayload{
			Event: "bounty_perfect_end",
		})
	}()
}
