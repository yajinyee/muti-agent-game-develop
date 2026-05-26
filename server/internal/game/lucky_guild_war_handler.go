// Package game — T139 幸運公會戰魚 handler
// server-event-agent 負責維護
// 業界依據：Fishing Frenzy Chapter 3「Guild Wars — guilds compete to capture tiles
//           through Boss Fish battles, territory control and leaderboard rankings」
// 設計：擊破後觸發「公會戰」，全服玩家 30 秒內共同擊破目標積分；
//       積分目標依在線玩家數動態調整（1人=15分, 2人=22分, 4人=35分）；
//       達成積分目標 → 「公會勝利」：全服 ×4.5 加成 10 秒；
//       依貢獻比例分配個人獎勵；個人冷卻 30 秒；全服冷卻 50 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyGuildWarManager struct {
	mu            sync.Mutex
	personalCD    map[string]time.Time
	globalCD      time.Time
	activeSession *guildWarSession
	victoryBoost  *guildWarVictoryBoost
}

type guildWarVictoryBoost struct {
	mult      float64
	expiresAt time.Time
}

type guildWarSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	targetPoints      int
	currentPoints     int
	contributions     map[string]int // playerID -> points
	expiresAt         time.Time
	settled           bool
}

func newLuckyGuildWarManager() *luckyGuildWarManager {
	return &luckyGuildWarManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyGuildWarFish(defID string) bool {
	return defID == "T139"
}

func (m *luckyGuildWarManager) getGuildWarVictoryMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.victoryBoost != nil && time.Now().Before(m.victoryBoost.expiresAt) {
		return m.victoryBoost.mult
	}
	return 1.0
}

func (m *luckyGuildWarManager) isGuildWarActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.activeSession != nil && !m.activeSession.settled
}

func (g *Game) tryLuckyGuildWarFish(playerID, playerName string) {
	m := g.luckyGuildWar
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
	m.personalCD[playerID] = now.Add(30 * time.Second)
	m.globalCD = now.Add(50 * time.Second)

	// 動態積分目標
	g.mu.RLock()
	playerCount := len(g.players)
	g.mu.RUnlock()
	targetPoints := 15
	if playerCount >= 4 {
		targetPoints = 35
	} else if playerCount >= 2 {
		targetPoints = 22
	}

	session := &guildWarSession{
		triggerPlayerID:   playerID,
		triggerPlayerName: playerName,
		targetPoints:      targetPoints,
		currentPoints:     0,
		contributions:     make(map[string]int),
		expiresAt:         now.Add(30 * time.Second),
		settled:           false,
	}
	m.activeSession = session
	m.mu.Unlock()

	log.Printf("[LuckyGuildWar] Triggered by %s, target=%d points", playerName, targetPoints)

	g.hub.Broadcast(protocol.MsgLuckyGuildWar, protocol.LuckyGuildWarPayload{
		Event:        "war_start",
		PlayerID:     playerID,
		PlayerName:   playerName,
		TargetPoints: targetPoints,
		Duration:     30.0,
	})

	go func() {
		time.Sleep(30 * time.Second)
		m.mu.Lock()
		if session.settled {
			m.mu.Unlock()
			return
		}
		session.settled = true
		currentPoints := session.currentPoints
		m.mu.Unlock()

		g.hub.Broadcast(protocol.MsgLuckyGuildWar, protocol.LuckyGuildWarPayload{
			Event:         "war_timeout",
			PlayerID:      playerID,
			PlayerName:    playerName,
			CurrentPoints: currentPoints,
			TargetPoints:  targetPoints,
		})
	}()
}

func (g *Game) notifyGuildWarKill(killerID, killerName string) {
	m := g.luckyGuildWar
	m.mu.Lock()
	if m.activeSession == nil || m.activeSession.settled {
		m.mu.Unlock()
		return
	}
	session := m.activeSession
	session.currentPoints++
	session.contributions[killerID]++
	currentPoints := session.currentPoints
	targetPoints := session.targetPoints
	triggerID := session.triggerPlayerID
	triggerName := session.triggerPlayerName
	m.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyGuildWar, protocol.LuckyGuildWarPayload{
		Event:         "war_progress",
		PlayerID:      triggerID,
		PlayerName:    triggerName,
		KillerID:      killerID,
		KillerName:    killerName,
		CurrentPoints: currentPoints,
		TargetPoints:  targetPoints,
	})

	if currentPoints >= targetPoints {
		m.mu.Lock()
		if session.settled {
			m.mu.Unlock()
			return
		}
		session.settled = true
		m.mu.Unlock()
		g.doGuildWarVictory(triggerID, triggerName, currentPoints)
	}
}

func (g *Game) doGuildWarVictory(playerID, playerName string, points int) {
	m := g.luckyGuildWar
	m.mu.Lock()
	m.victoryBoost = &guildWarVictoryBoost{
		mult:      4.5,
		expiresAt: time.Now().Add(10 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyGuildWar] Victory! %s points=%d → global ×4.5 for 10s", playerName, points)

	g.hub.Broadcast(protocol.MsgLuckyGuildWar, protocol.LuckyGuildWarPayload{
		Event:         "war_victory",
		PlayerID:      playerID,
		PlayerName:    playerName,
		CurrentPoints: points,
		BoostMult:     4.5,
		BoostSec:      10,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  fmt.Sprintf("⚔️ 公會勝利！全服達成 %d 積分！全服 ×4.5 加成 10 秒！", points),
		Priority: "critical",
		Color:    "#FFD700",
	})

	go func() {
		time.Sleep(10 * time.Second)
		m.mu.Lock()
		m.victoryBoost = nil
		m.mu.Unlock()
		g.hub.Broadcast(protocol.MsgLuckyGuildWar, protocol.LuckyGuildWarPayload{
			Event: "war_victory_end",
		})
	}()
}
