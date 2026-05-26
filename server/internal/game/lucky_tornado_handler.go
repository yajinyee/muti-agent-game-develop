// Package game — T141 幸運龍捲風魚 handler
// server-event-agent 負責維護
// 業界依據：Fishing Fortune 2026「Tornado sweep — spinning vortex sweeps across screen,
//           pulling fish into the funnel and dealing massive area damage」
// 設計：擊破後觸發「龍捲風橫掃」10 秒；
//       龍捲風從左向右橫掃，每 2 秒造成全場 HP -40%；
//       龍捲風期間擊破 ≥ 8 個目標 → 「完美龍捲風」：全服 ×3.8 加成 9 秒；
//       個人冷卻 30 秒；全服冷卻 50 秒
package game

import (
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyTornadoManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	activeSession *tornadoSession
	perfectBoost  *tornadoPerfectBoost
}

type tornadoPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type tornadoSession struct {
	playerID   string
	playerName string
	killCount  int
	expiresAt  time.Time
	settled    bool
}

func newLuckyTornadoManager() *luckyTornadoManager {
	return &luckyTornadoManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyTornadoFish(defID string) bool {
	return defID == "T141"
}

func (m *luckyTornadoManager) getTornadoPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyTornadoManager) isTornadoActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.activeSession != nil && !m.activeSession.settled && time.Now().Before(m.activeSession.expiresAt)
}

func (m *luckyTornadoManager) notifyTornadoKill(g *Game, playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil || m.activeSession.settled {
		return
	}
	m.activeSession.killCount++
}

func (m *luckyTornadoManager) tryLuckyTornadoFish(g *Game, playerID, playerName string) bool {
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

	m.personalCD[playerID] = now.Add(30 * time.Second)
	m.globalCD = now.Add(50 * time.Second)

	session := &tornadoSession{
		playerID:   playerID,
		playerName: playerName,
		killCount:  0,
		expiresAt:  now.Add(10 * time.Second),
		settled:    false,
	}
	m.activeSession = session

	// 廣播龍捲風開始
	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyTornado,
		Payload: protocol.LuckyTornadoPayload{
			Event:      "tornado_start",
			PlayerID:   playerID,
			PlayerName: playerName,
			Duration:   10.0,
		},
	})

	log.Printf("[LuckyTornado] %s 觸發龍捲風橫掃", playerName)

	go m.runTornadoSweep(g, session)
	return true
}

func (m *luckyTornadoManager) runTornadoSweep(g *Game, session *tornadoSession) {
	// 每 2 秒一波，共 5 波
	for wave := 1; wave <= 5; wave++ {
		time.Sleep(2 * time.Second)

		m.mu.Lock()
		if session.settled {
			m.mu.Unlock()
			return
		}
		m.mu.Unlock()

		// 對全場目標造成 HP -40%
		hitCount := g.applyAOEDamage(0, 0, 99999, 0.40)

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyTornado,
			Payload: protocol.LuckyTornadoPayload{
				Event:      "tornado_sweep",
				PlayerID:   session.playerID,
				PlayerName: session.playerName,
				WaveNum:    wave,
				HitCount:   hitCount,
			},
		})
	}

	// 龍捲風結束
	m.mu.Lock()
	if session.settled {
		m.mu.Unlock()
		return
	}
	killCount := session.killCount
	session.settled = true
	m.mu.Unlock()

	// 判定完美龍捲風
	if killCount >= 8 {
		m.doTornadoPerfect(g, session.playerID, session.playerName)
	} else {
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyTornado,
			Payload: protocol.LuckyTornadoPayload{
				Event:      "tornado_end",
				PlayerID:   session.playerID,
				PlayerName: session.playerName,
				KillCount:  killCount,
			},
		})
	}
}

func (m *luckyTornadoManager) doTornadoPerfect(g *Game, playerID, playerName string) {
	m.mu.Lock()
	m.perfectBoost = &tornadoPerfectBoost{
		mult:      3.8,
		expiresAt: time.Now().Add(9 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyTornado,
		Payload: protocol.LuckyTornadoPayload{
			Event:      "tornado_perfect",
			PlayerID:   playerID,
			PlayerName: playerName,
			BoostMult:  3.8,
			BoostSec:   9,
		},
	})

	g.sendAnnounce("🌪️ 完美龍捲風！"+playerName+" 觸發全服 ×3.8 加成 9 秒！", "high", "#00E5FF")

	time.Sleep(9 * time.Second)
	m.mu.Lock()
	m.perfectBoost = nil
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyTornado,
		Payload: protocol.LuckyTornadoPayload{
			Event:      "tornado_perfect_end",
			PlayerID:   playerID,
			PlayerName: playerName,
		},
	})
}
