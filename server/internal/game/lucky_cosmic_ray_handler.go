// Package game — T144 幸運星際魚 handler
// server-event-agent 負責維護
// 業界依據：Fishing Fortune 2026「Cosmic ray — 8-directional energy beams sweep the screen,
//           hitting all fish in their path for massive damage」
// 設計：擊破後觸發「星際射線」8 方向光束掃射；
//       每方向光束 HP -30%，命中沿線所有目標；
//       8 方向全部命中 ≥ 16 個目標 → 「完美星際」：全服 ×4.5 加成 10 秒；
//       個人冷卻 36 秒；全服冷卻 58 秒
package game

import (
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyCosmicRayManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	activeSession *cosmicRaySession
	perfectBoost  *cosmicRayPerfectBoost
}

type cosmicRayPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type cosmicRaySession struct {
	playerID      string
	playerName    string
	totalHitCount int
	expiresAt     time.Time
	settled       bool
}

func newLuckyCosmicRayManager() *luckyCosmicRayManager {
	return &luckyCosmicRayManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyCosmicRayFish(defID string) bool {
	return defID == "T144"
}

func (m *luckyCosmicRayManager) getCosmicRayPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyCosmicRayManager) tryLuckyCosmicRayFish(g *Game, playerID, playerName string) bool {
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

	m.personalCD[playerID] = now.Add(36 * time.Second)
	m.globalCD = now.Add(58 * time.Second)

	session := &cosmicRaySession{
		playerID:   playerID,
		playerName: playerName,
		expiresAt:  now.Add(10 * time.Second),
		settled:    false,
	}
	m.activeSession = session

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyCosmicRay,
		Payload: protocol.LuckyCosmicRayPayload{
			Event:      "cosmic_start",
			PlayerID:   playerID,
			PlayerName: playerName,
			RayCount:   8,
		},
	})

	log.Printf("[LuckyCosmicRay] %s 觸發星際射線", playerName)

	go m.runCosmicRays(g, session)
	return true
}

func (m *luckyCosmicRayManager) runCosmicRays(g *Game, session *cosmicRaySession) {
	totalHit := 0

	// 8 方向依序發射，每 0.5 秒一道
	for dir := 0; dir < 8; dir++ {
		time.Sleep(500 * time.Millisecond)

		m.mu.Lock()
		if session.settled {
			m.mu.Unlock()
			return
		}
		m.mu.Unlock()

		// 每道光束對全場造成 HP -30%（簡化：全場 AOE）
		hitCount := g.applyAOEDamage(0, 0, 99999, 0.30)
		totalHit += hitCount

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyCosmicRay,
			Payload: protocol.LuckyCosmicRayPayload{
				Event:         "cosmic_ray",
				PlayerID:      session.playerID,
				PlayerName:    session.playerName,
				Direction:     dir,
				HitCount:      hitCount,
				TotalHitCount: totalHit,
			},
		})
	}

	m.mu.Lock()
	session.totalHitCount = totalHit
	session.settled = true
	m.mu.Unlock()

	if totalHit >= 16 {
		m.doCosmicRayPerfect(g, session.playerID, session.playerName)
	} else {
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyCosmicRay,
			Payload: protocol.LuckyCosmicRayPayload{
				Event:         "cosmic_end",
				PlayerID:      session.playerID,
				PlayerName:    session.playerName,
				TotalHitCount: totalHit,
			},
		})
	}
}

func (m *luckyCosmicRayManager) doCosmicRayPerfect(g *Game, playerID, playerName string) {
	m.mu.Lock()
	m.perfectBoost = &cosmicRayPerfectBoost{
		mult:      4.5,
		expiresAt: time.Now().Add(10 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyCosmicRay,
		Payload: protocol.LuckyCosmicRayPayload{
			Event:      "cosmic_perfect",
			PlayerID:   playerID,
			PlayerName: playerName,
			BoostMult:  4.5,
			BoostSec:   10,
		},
	})

	g.sendAnnounce("✨ 完美星際！"+playerName+" 觸發全服 ×4.5 加成 10 秒！", "high", "#9C27B0")

	time.Sleep(10 * time.Second)
	m.mu.Lock()
	m.perfectBoost = nil
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyCosmicRay,
		Payload: protocol.LuckyCosmicRayPayload{
			Event:      "cosmic_perfect_end",
			PlayerID:   playerID,
			PlayerName: playerName,
		},
	})
}
