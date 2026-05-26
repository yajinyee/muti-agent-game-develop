// Package game — T143 幸運火山魚 handler
// server-event-agent 負責維護
// 業界依據：Jili Games 2026「Volcano eruption — lava bombs rain down randomly,
//           each creating an AOE explosion on impact」
// 設計：擊破後觸發「火山爆發」，10 顆熔岩彈隨機落下（每 0.8 秒一顆）；
//       每顆熔岩彈 AOE r=140px，HP -35%；
//       10 顆全部命中至少 1 個目標 → 「完美火山」：全服 ×4.2 加成 10 秒；
//       個人冷卻 34 秒；全服冷卻 55 秒
package game

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyVolcanoManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	activeSession *volcanoSession
	perfectBoost  *volcanoPerfectBoost
}

type volcanoPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type volcanoSession struct {
	playerID      string
	playerName    string
	hitBombs      int // 命中至少 1 個目標的熔岩彈數
	expiresAt     time.Time
	settled       bool
}

func newLuckyVolcanoManager() *luckyVolcanoManager {
	return &luckyVolcanoManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyVolcanoFish(defID string) bool {
	return defID == "T143"
}

func (m *luckyVolcanoManager) getVolcanoPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyVolcanoManager) tryLuckyVolcanoFish(g *Game, playerID, playerName string) bool {
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

	m.personalCD[playerID] = now.Add(34 * time.Second)
	m.globalCD = now.Add(55 * time.Second)

	session := &volcanoSession{
		playerID:   playerID,
		playerName: playerName,
		expiresAt:  now.Add(12 * time.Second),
		settled:    false,
	}
	m.activeSession = session

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyVolcano,
		Payload: protocol.LuckyVolcanoPayload{
			Event:      "volcano_erupt",
			PlayerID:   playerID,
			PlayerName: playerName,
			BombCount:  10,
		},
	})

	log.Printf("[LuckyVolcano] %s 觸發火山爆發", playerName)

	go m.runVolcanoBombs(g, session)
	return true
}

func (m *luckyVolcanoManager) runVolcanoBombs(g *Game, session *volcanoSession) {
	hitBombs := 0

	for bomb := 1; bomb <= 10; bomb++ {
		time.Sleep(800 * time.Millisecond)

		m.mu.Lock()
		if session.settled {
			m.mu.Unlock()
			return
		}
		m.mu.Unlock()

		// 隨機落點
		bx := 100.0 + rand.Float64()*900.0
		by := 100.0 + rand.Float64()*500.0

		hitCount := g.applyAOEDamage(bx, by, 140, 0.35)
		if hitCount > 0 {
			hitBombs++
		}

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyVolcano,
			Payload: protocol.LuckyVolcanoPayload{
				Event:      "lava_bomb",
				PlayerID:   session.playerID,
				PlayerName: session.playerName,
				BombNum:    bomb,
				BombX:      bx,
				BombY:      by,
				HitCount:   hitCount,
				HitBombs:   hitBombs,
			},
		})
	}

	m.mu.Lock()
	session.hitBombs = hitBombs
	session.settled = true
	m.mu.Unlock()

	if hitBombs >= 10 {
		m.doVolcanoPerfect(g, session.playerID, session.playerName)
	} else {
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyVolcano,
			Payload: protocol.LuckyVolcanoPayload{
				Event:      "volcano_end",
				PlayerID:   session.playerID,
				PlayerName: session.playerName,
				HitBombs:   hitBombs,
			},
		})
	}
}

func (m *luckyVolcanoManager) doVolcanoPerfect(g *Game, playerID, playerName string) {
	m.mu.Lock()
	m.perfectBoost = &volcanoPerfectBoost{
		mult:      4.2,
		expiresAt: time.Now().Add(10 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyVolcano,
		Payload: protocol.LuckyVolcanoPayload{
			Event:      "volcano_perfect",
			PlayerID:   playerID,
			PlayerName: playerName,
			BoostMult:  4.2,
			BoostSec:   10,
		},
	})

	g.sendAnnounce("🌋 完美火山！"+playerName+" 觸發全服 ×4.2 加成 10 秒！", "high", "#FF4500")

	time.Sleep(10 * time.Second)
	m.mu.Lock()
	m.perfectBoost = nil
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyVolcano,
		Payload: protocol.LuckyVolcanoPayload{
			Event:      "volcano_perfect_end",
			PlayerID:   playerID,
			PlayerName: playerName,
		},
	})
}
