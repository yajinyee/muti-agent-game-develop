// Package game — T146 幸運量子魚 handler
// server-event-agent 負責維護
// 業界依據：量子力學 Crash mechanic 升級版 — 量子疊加態觀測機制
// 設計：擊破後觸發「量子觀測」；
//       場上所有目標同時有 50% 機率被觀測到（HP -60%）；
//       觀測到 ≥ 10 個 → 「量子坍縮」：全服 ×5.5 加成 12 秒；
//       個人冷卻 42 秒；全服冷卻 70 秒
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyQuantumManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *quantumPerfectBoost
}

type quantumPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyQuantumManager() *luckyQuantumManager {
	return &luckyQuantumManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyQuantumFish(defID string) bool {
	return defID == "T146"
}

func (m *luckyQuantumManager) getQuantumPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyQuantumManager) tryLuckyQuantumFish(g *Game, playerID, playerName string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	if cd, ok := m.personalCD[playerID]; ok && now.Before(cd) {
		return false
	}
	if now.Before(m.globalCD) {
		return false
	}

	m.personalCD[playerID] = now.Add(42 * time.Second)
	m.globalCD = now.Add(70 * time.Second)

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyQuantum,
		Payload: protocol.LuckyQuantumPayload{
			Event:      "quantum_observe",
			PlayerID:   playerID,
			PlayerName: playerName,
		},
	})

	log.Printf("[LuckyQuantum] %s 觸發量子觀測", playerName)

	go m.runQuantumObservation(g, playerID, playerName)
	return true
}

func (m *luckyQuantumManager) runQuantumObservation(g *Game, playerID, playerName string) {
	time.Sleep(500 * time.Millisecond)

	// 量子觀測：每個目標 50% 機率被觀測（HP -60%）
	g.mu.Lock()
	observedCount := 0
	for _, t := range g.targets {
		if t.Def.ID == "B001" {
			continue
		}
		if rand.Float64() < 0.50 {
			dmg := int(float64(t.MaxHP) * 0.60)
			t.HP -= dmg
			if t.HP < 0 {
				t.HP = 0
			}
			observedCount++
		}
	}
	g.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyQuantum,
		Payload: protocol.LuckyQuantumPayload{
			Event:         "quantum_result",
			PlayerID:      playerID,
			PlayerName:    playerName,
			ObservedCount: observedCount,
		},
	})

	if observedCount >= 10 {
		m.doQuantumCollapse(g, playerID, playerName, observedCount)
	}
}

func (m *luckyQuantumManager) doQuantumCollapse(g *Game, playerID, playerName string, observed int) {
	m.mu.Lock()
	m.perfectBoost = &quantumPerfectBoost{
		mult:      5.5,
		expiresAt: time.Now().Add(12 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyQuantum,
		Payload: protocol.LuckyQuantumPayload{
			Event:         "quantum_collapse",
			PlayerID:      playerID,
			PlayerName:    playerName,
			ObservedCount: observed,
			BoostMult:     5.5,
			BoostSec:      12,
		},
	})

	g.sendAnnounce("⚛️ 量子坍縮！"+playerName+" 觀測 "+fmt.Sprintf("%d", observed)+" 個！全服 ×5.5 加成 12 秒！", "critical", "#00E5FF")

	time.Sleep(12 * time.Second)
	m.mu.Lock()
	m.perfectBoost = nil
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyQuantum,
		Payload: protocol.LuckyQuantumPayload{
			Event:      "quantum_collapse_end",
			PlayerID:   playerID,
			PlayerName: playerName,
		},
	})
}
