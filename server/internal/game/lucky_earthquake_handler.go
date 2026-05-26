// Package game — T142 幸運地震魚 handler
// server-event-agent 負責維護
// 業界依據：Fishing Fortune 2026「Earthquake shockwave — seismic waves radiate outward
//           in concentric rings, dealing escalating damage with each ring」
// 設計：擊破後觸發「地震波」3 波同心圓衝擊；
//       第 1 波 HP -25%，第 2 波 HP -35%，第 3 波 HP -45%；
//       三波命中總數 ≥ 12 → 「完美地震」：全服 ×4.0 加成 9 秒；
//       個人冷卻 32 秒；全服冷卻 52 秒
package game

import (
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyEarthquakeManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	activeSession *earthquakeSession
	perfectBoost  *earthquakePerfectBoost
}

type earthquakePerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type earthquakeSession struct {
	playerID      string
	playerName    string
	totalHitCount int
	expiresAt     time.Time
	settled       bool
}

func newLuckyEarthquakeManager() *luckyEarthquakeManager {
	return &luckyEarthquakeManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyEarthquakeFish(defID string) bool {
	return defID == "T142"
}

func (m *luckyEarthquakeManager) getEarthquakePerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyEarthquakeManager) tryLuckyEarthquakeFish(g *Game, playerID, playerName string) bool {
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

	m.personalCD[playerID] = now.Add(32 * time.Second)
	m.globalCD = now.Add(52 * time.Second)

	session := &earthquakeSession{
		playerID:   playerID,
		playerName: playerName,
		expiresAt:  now.Add(15 * time.Second),
		settled:    false,
	}
	m.activeSession = session

	// 廣播地震警告
	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyEarthquake,
		Payload: protocol.LuckyEarthquakePayload{
			Event:      "quake_warning",
			PlayerID:   playerID,
			PlayerName: playerName,
			WaveCount:  3,
		},
	})

	log.Printf("[LuckyEarthquake] %s 觸發地震波", playerName)

	go m.runEarthquakeWaves(g, session)
	return true
}

func (m *luckyEarthquakeManager) runEarthquakeWaves(g *Game, session *earthquakeSession) {
	damagePcts := []float64{0.25, 0.35, 0.45}
	totalHit := 0

	for i, pct := range damagePcts {
		time.Sleep(3 * time.Second)

		m.mu.Lock()
		if session.settled {
			m.mu.Unlock()
			return
		}
		m.mu.Unlock()

		hitCount := g.applyAOEDamage(0, 0, 99999, pct)
		totalHit += hitCount

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyEarthquake,
			Payload: protocol.LuckyEarthquakePayload{
				Event:         "quake_wave",
				PlayerID:      session.playerID,
				PlayerName:    session.playerName,
				WaveNum:       i + 1,
				DamagePct:     pct,
				HitCount:      hitCount,
				TotalHitCount: totalHit,
			},
		})
	}

	m.mu.Lock()
	session.totalHitCount = totalHit
	session.settled = true
	m.mu.Unlock()

	if totalHit >= 12 {
		m.doEarthquakePerfect(g, session.playerID, session.playerName)
	} else {
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyEarthquake,
			Payload: protocol.LuckyEarthquakePayload{
				Event:         "quake_end",
				PlayerID:      session.playerID,
				PlayerName:    session.playerName,
				TotalHitCount: totalHit,
			},
		})
	}
}

func (m *luckyEarthquakeManager) doEarthquakePerfect(g *Game, playerID, playerName string) {
	m.mu.Lock()
	m.perfectBoost = &earthquakePerfectBoost{
		mult:      4.0,
		expiresAt: time.Now().Add(9 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyEarthquake,
		Payload: protocol.LuckyEarthquakePayload{
			Event:      "quake_perfect",
			PlayerID:   playerID,
			PlayerName: playerName,
			BoostMult:  4.0,
			BoostSec:   9,
		},
	})

	g.sendAnnounce("🌍 完美地震！"+playerName+" 觸發全服 ×4.0 加成 9 秒！", "high", "#FF6B35")

	time.Sleep(9 * time.Second)
	m.mu.Lock()
	m.perfectBoost = nil
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyEarthquake,
		Payload: protocol.LuckyEarthquakePayload{
			Event:      "quake_perfect_end",
			PlayerID:   playerID,
			PlayerName: playerName,
		},
	})
}
