// Package game — T150 幸運重生魚 handler
// server-event-agent 負責維護
// 業界依據：Phoenix rebirth mechanic — 死亡目標復活再擊破，雙重獎勵
// 設計：擊破後觸發「重生之力」15 秒；
//       15 秒內死亡的所有目標立即以 HP 50% 復活一次；
//       復活目標被擊破獎勵 ×3.0；
//       復活擊破 ≥ 8 個 → 「完美重生」：全服 ×6.5 加成 15 秒；
//       個人冷卻 48 秒；全服冷卻 78 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyRebirthManager struct {
	mu              sync.Mutex
	personalCD      map[string]time.Time
	globalCD        time.Time
	activeSession   *rebirthSession
	perfectBoost    *rebirthPerfectBoost
	rebirthTargets  map[string]bool // instance_id -> is_rebirth
}

type rebirthPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type rebirthSession struct {
	playerID     string
	playerName   string
	rebirthKills int
	expiresAt    time.Time
	settled      bool
}

func newLuckyRebirthManager() *luckyRebirthManager {
	return &luckyRebirthManager{
		personalCD:     make(map[string]time.Time),
		rebirthTargets: make(map[string]bool),
	}
}

func isLuckyRebirthFish(defID string) bool {
	return defID == "T150"
}

func (m *luckyRebirthManager) getRebirthPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyRebirthManager) getRebirthKillMult(instanceID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.rebirthTargets[instanceID] {
		return 3.0
	}
	return 1.0
}

func (m *luckyRebirthManager) isRebirthActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.activeSession != nil && !m.activeSession.settled && time.Now().Before(m.activeSession.expiresAt)
}

func (m *luckyRebirthManager) notifyRebirthKill(g *Game, instanceID string) {
	m.mu.Lock()
	if !m.rebirthTargets[instanceID] {
		m.mu.Unlock()
		return
	}
	delete(m.rebirthTargets, instanceID)
	if m.activeSession == nil || m.activeSession.settled {
		m.mu.Unlock()
		return
	}
	m.activeSession.rebirthKills++
	kills := m.activeSession.rebirthKills
	name := m.activeSession.playerName
	pid := m.activeSession.playerID
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyRebirth,
		Payload: protocol.LuckyRebirthPayload{
			Event:        "rebirth_kill",
			PlayerID:     pid,
			PlayerName:   name,
			RebirthKills: kills,
		},
	})
}

func (m *luckyRebirthManager) tryLuckyRebirthFish(g *Game, playerID, playerName string) bool {
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

	m.personalCD[playerID] = now.Add(48 * time.Second)
	m.globalCD = now.Add(78 * time.Second)

	session := &rebirthSession{
		playerID:     playerID,
		playerName:   playerName,
		rebirthKills: 0,
		expiresAt:    now.Add(15 * time.Second),
		settled:      false,
	}
	m.activeSession = session

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyRebirth,
		Payload: protocol.LuckyRebirthPayload{
			Event:      "rebirth_start",
			PlayerID:   playerID,
			PlayerName: playerName,
			Duration:   15.0,
		},
	})

	log.Printf("[LuckyRebirth] %s 觸發重生之力", playerName)

	go m.runRebirthTimeout(g, session)
	return true
}

func (m *luckyRebirthManager) runRebirthTimeout(g *Game, session *rebirthSession) {
	time.Sleep(15 * time.Second)

	m.mu.Lock()
	if session.settled {
		m.mu.Unlock()
		return
	}
	session.settled = true
	kills := session.rebirthKills
	// 清除所有重生標記
	m.rebirthTargets = make(map[string]bool)
	m.mu.Unlock()

	if kills >= 8 {
		m.doRebirthPerfect(g, session.playerID, session.playerName, kills)
	} else {
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyRebirth,
			Payload: protocol.LuckyRebirthPayload{
				Event:        "rebirth_end",
				PlayerID:     session.playerID,
				PlayerName:   session.playerName,
				RebirthKills: kills,
			},
		})
	}
}

func (m *luckyRebirthManager) doRebirthPerfect(g *Game, playerID, playerName string, kills int) {
	m.mu.Lock()
	m.perfectBoost = &rebirthPerfectBoost{
		mult:      6.5,
		expiresAt: time.Now().Add(15 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyRebirth,
		Payload: protocol.LuckyRebirthPayload{
			Event:        "rebirth_perfect",
			PlayerID:     playerID,
			PlayerName:   playerName,
			RebirthKills: kills,
			BoostMult:    6.5,
			BoostSec:     15,
		},
	})

	g.sendAnnounce("🔥 完美重生！"+playerName+" 重生擊破 "+fmt.Sprintf("%d", kills)+" 個！全服 ×6.5 加成 15 秒！", "critical", "#FF4500")

	time.Sleep(15 * time.Second)
	m.mu.Lock()
	m.perfectBoost = nil
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyRebirth,
		Payload: protocol.LuckyRebirthPayload{
			Event:      "rebirth_perfect_end",
			PlayerID:   playerID,
			PlayerName: playerName,
		},
	})
}
