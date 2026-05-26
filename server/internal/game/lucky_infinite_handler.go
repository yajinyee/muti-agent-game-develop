// Package game — T148 幸運無限魚 handler
// server-event-agent 負責維護
// 業界依據：Infinite multiplier accumulation — 每次擊破倍率無限累積
// 設計：擊破後觸發「無限模式」20 秒；
//       每次擊破倍率 +1.0x（無上限）；
//       20 秒後結算，最終倍率 ≥ 20x → 「無限完美」：全服 ×6.0 加成 15 秒；
//       個人冷卻 46 秒；全服冷卻 75 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyInfiniteManager struct {
	mu            sync.Mutex
	personalCD    map[string]time.Time
	globalCD      time.Time
	activeSession *infiniteSession
	perfectBoost  *infinitePerfectBoost
}

type infinitePerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type infiniteSession struct {
	playerID    string
	playerName  string
	accumMult   float64
	killCount   int
	expiresAt   time.Time
	settled     bool
}

func newLuckyInfiniteManager() *luckyInfiniteManager {
	return &luckyInfiniteManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyInfiniteFish(defID string) bool {
	return defID == "T148"
}

func (m *luckyInfiniteManager) getInfinitePerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyInfiniteManager) isInfiniteActive(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.activeSession != nil &&
		m.activeSession.playerID == playerID &&
		!m.activeSession.settled &&
		time.Now().Before(m.activeSession.expiresAt)
}

func (m *luckyInfiniteManager) notifyInfiniteKill(g *Game, playerID string) {
	m.mu.Lock()
	if m.activeSession == nil || m.activeSession.playerID != playerID || m.activeSession.settled {
		m.mu.Unlock()
		return
	}
	m.activeSession.accumMult += 1.0
	m.activeSession.killCount++
	mult := m.activeSession.accumMult
	kills := m.activeSession.killCount
	name := m.activeSession.playerName
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyInfinite,
		Payload: protocol.LuckyInfinitePayload{
			Event:      "infinite_kill",
			PlayerID:   playerID,
			PlayerName: name,
			AccumMult:  mult,
			KillCount:  kills,
		},
	})
}

func (m *luckyInfiniteManager) tryLuckyInfiniteFish(g *Game, playerID, playerName string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	if cd, ok := m.personalCD[playerID]; ok && now.Before(cd) {
		return false
	}
	if now.Before(m.globalCD) {
		return false
	}
	if m.activeSession != nil && !m.activeSession.settled && now.Before(m.activeSession.expiresAt) {
		return false
	}

	m.personalCD[playerID] = now.Add(46 * time.Second)
	m.globalCD = now.Add(75 * time.Second)

	session := &infiniteSession{
		playerID:   playerID,
		playerName: playerName,
		accumMult:  1.0,
		killCount:  0,
		expiresAt:  now.Add(20 * time.Second),
		settled:    false,
	}
	m.activeSession = session

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyInfinite,
		Payload: protocol.LuckyInfinitePayload{
			Event:      "infinite_start",
			PlayerID:   playerID,
			PlayerName: playerName,
			Duration:   20.0,
		},
	})

	log.Printf("[LuckyInfinite] %s 觸發無限模式", playerName)

	go m.runInfiniteTimeout(g, session)
	return true
}

func (m *luckyInfiniteManager) runInfiniteTimeout(g *Game, session *infiniteSession) {
	time.Sleep(20 * time.Second)

	m.mu.Lock()
	if session.settled {
		m.mu.Unlock()
		return
	}
	session.settled = true
	finalMult := session.accumMult
	kills := session.killCount
	m.mu.Unlock()

	if finalMult >= 20.0 {
		m.doInfinitePerfect(g, session.playerID, session.playerName, finalMult, kills)
	} else {
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyInfinite,
			Payload: protocol.LuckyInfinitePayload{
				Event:      "infinite_end",
				PlayerID:   session.playerID,
				PlayerName: session.playerName,
				AccumMult:  finalMult,
				KillCount:  kills,
			},
		})
	}
}

func (m *luckyInfiniteManager) doInfinitePerfect(g *Game, playerID, playerName string, finalMult float64, kills int) {
	m.mu.Lock()
	m.perfectBoost = &infinitePerfectBoost{
		mult:      6.0,
		expiresAt: time.Now().Add(15 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyInfinite,
		Payload: protocol.LuckyInfinitePayload{
			Event:      "infinite_perfect",
			PlayerID:   playerID,
			PlayerName: playerName,
			AccumMult:  finalMult,
			KillCount:  kills,
			BoostMult:  6.0,
			BoostSec:   15,
		},
	})

	g.sendAnnounce("♾️ 無限完美！"+playerName+" 累積 ×"+fmt.Sprintf("%.1f", finalMult)+"！全服 ×6.0 加成 15 秒！", "critical", "#7B2FBE")

	time.Sleep(15 * time.Second)
	m.mu.Lock()
	m.perfectBoost = nil
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyInfinite,
		Payload: protocol.LuckyInfinitePayload{
			Event:      "infinite_perfect_end",
			PlayerID:   playerID,
			PlayerName: playerName,
		},
	})
}
