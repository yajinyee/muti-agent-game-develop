// Package game — T145 幸運神龍魚 handler
// server-event-agent 負責維護
// 業界依據：Royal Fishing Jili「Divine Dragon — legendary creature descends from the heavens,
//           unleashing devastating claw strikes that devastate the entire battlefield」
// 設計：擊破後「神龍降臨」20 秒；
//       每 4 秒神龍爪擊（HP -50%，AOE 全場）；
//       5 次爪擊全部命中 ≥ 5 個目標 → 「神龍完美」：全服 ×5.0 加成 12 秒；
//       個人冷卻 40 秒；全服冷卻 65 秒
package game

import (
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyDivineDragonManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	activeSession *divineDragonSession
	perfectBoost  *divineDragonPerfectBoost
}

type divineDragonPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type divineDragonSession struct {
	playerID       string
	playerName     string
	perfectClaws   int // 命中 ≥5 個目標的爪擊數
	expiresAt      time.Time
	settled        bool
}

func newLuckyDivineDragonManager() *luckyDivineDragonManager {
	return &luckyDivineDragonManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyDivineDragonFish(defID string) bool {
	return defID == "T145"
}

func (m *luckyDivineDragonManager) getDivineDragonPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyDivineDragonManager) tryLuckyDivineDragonFish(g *Game, playerID, playerName string) bool {
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

	m.personalCD[playerID] = now.Add(40 * time.Second)
	m.globalCD = now.Add(65 * time.Second)

	session := &divineDragonSession{
		playerID:   playerID,
		playerName: playerName,
		expiresAt:  now.Add(20 * time.Second),
		settled:    false,
	}
	m.activeSession = session

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyDivineDragon,
		Payload: protocol.LuckyDivineDragonPayload{
			Event:      "dragon_descend",
			PlayerID:   playerID,
			PlayerName: playerName,
			Duration:   20.0,
			ClawCount:  5,
		},
	})

	log.Printf("[LuckyDivineDragon] %s 觸發神龍降臨", playerName)

	go m.runDivineDragonClaws(g, session)
	return true
}

func (m *luckyDivineDragonManager) runDivineDragonClaws(g *Game, session *divineDragonSession) {
	perfectClaws := 0

	for claw := 1; claw <= 5; claw++ {
		time.Sleep(4 * time.Second)

		m.mu.Lock()
		if session.settled {
			m.mu.Unlock()
			return
		}
		m.mu.Unlock()

		// 神龍爪擊：全場 HP -50%
		hitCount := g.applyAOEDamage(0, 0, 99999, 0.50)
		if hitCount >= 5 {
			perfectClaws++
		}

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyDivineDragon,
			Payload: protocol.LuckyDivineDragonPayload{
				Event:        "dragon_claw",
				PlayerID:     session.playerID,
				PlayerName:   session.playerName,
				ClawNum:      claw,
				HitCount:     hitCount,
				PerfectClaws: perfectClaws,
			},
		})
	}

	m.mu.Lock()
	session.perfectClaws = perfectClaws
	session.settled = true
	m.mu.Unlock()

	if perfectClaws >= 5 {
		m.doDivineDragonPerfect(g, session.playerID, session.playerName)
	} else {
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyDivineDragon,
			Payload: protocol.LuckyDivineDragonPayload{
				Event:        "dragon_leave",
				PlayerID:     session.playerID,
				PlayerName:   session.playerName,
				PerfectClaws: perfectClaws,
			},
		})
	}
}

func (m *luckyDivineDragonManager) doDivineDragonPerfect(g *Game, playerID, playerName string) {
	m.mu.Lock()
	m.perfectBoost = &divineDragonPerfectBoost{
		mult:      5.0,
		expiresAt: time.Now().Add(12 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyDivineDragon,
		Payload: protocol.LuckyDivineDragonPayload{
			Event:      "dragon_perfect",
			PlayerID:   playerID,
			PlayerName: playerName,
			BoostMult:  5.0,
			BoostSec:   12,
		},
	})

	g.sendAnnounce("🐉 神龍完美！"+playerName+" 觸發全服 ×5.0 加成 12 秒！", "critical", "#FFD700")

	time.Sleep(12 * time.Second)
	m.mu.Lock()
	m.perfectBoost = nil
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyDivineDragon,
		Payload: protocol.LuckyDivineDragonPayload{
			Event:      "dragon_perfect_end",
			PlayerID:   playerID,
			PlayerName: playerName,
		},
	})
}
