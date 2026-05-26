// Package game — T147 幸運超新星魚 handler
// server-event-agent 負責維護
// 業界依據：Supernova explosion mechanic — massive area damage + temporary multiplier boost
// 設計：擊破後超新星爆炸，全場 HP -70%；
//       爆炸後 5 秒內所有目標倍率 ×3.0；
//       命中 ≥ 8 個 → 「超新星完美」：全服 ×5.5 加成 12 秒；
//       個人冷卻 44 秒；全服冷卻 72 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckySupernovaManager struct {
	mu              sync.Mutex
	personalCD      map[string]time.Time
	globalCD        time.Time
	perfectBoost    *supernovaPerfectBoost
	multBoostActive bool
	multBoostExpiry time.Time
}

type supernovaPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckySupernovaManager() *luckySupernovaManager {
	return &luckySupernovaManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckySupernovaFish(defID string) bool {
	return defID == "T147"
}

func (m *luckySupernovaManager) getSupernovaPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckySupernovaManager) getSupernovaMultBoost() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.multBoostActive && time.Now().Before(m.multBoostExpiry) {
		return 3.0
	}
	return 1.0
}

func (m *luckySupernovaManager) tryLuckySupernovaFish(g *Game, playerID, playerName string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	if cd, ok := m.personalCD[playerID]; ok && now.Before(cd) {
		return false
	}
	if now.Before(m.globalCD) {
		return false
	}

	m.personalCD[playerID] = now.Add(44 * time.Second)
	m.globalCD = now.Add(72 * time.Second)

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckySupernova,
		Payload: protocol.LuckySupernovaPayload{
			Event:      "supernova_explode",
			PlayerID:   playerID,
			PlayerName: playerName,
		},
	})

	log.Printf("[LuckySupernova] %s 觸發超新星爆炸", playerName)

	go m.runSupernovaExplosion(g, playerID, playerName)
	return true
}

func (m *luckySupernovaManager) runSupernovaExplosion(g *Game, playerID, playerName string) {
	time.Sleep(300 * time.Millisecond)

	// 超新星爆炸：全場 HP -70%
	hitCount := g.applyAOEDamage(0, 0, 99999, 0.70)

	// 啟動 5 秒倍率加成
	m.mu.Lock()
	m.multBoostActive = true
	m.multBoostExpiry = time.Now().Add(5 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckySupernova,
		Payload: protocol.LuckySupernovaPayload{
			Event:      "supernova_boost",
			PlayerID:   playerID,
			PlayerName: playerName,
			HitCount:   hitCount,
			MultBoost:  3.0,
			BoostSec:   5,
		},
	})

	time.Sleep(5 * time.Second)
	m.mu.Lock()
	m.multBoostActive = false
	m.mu.Unlock()

	if hitCount >= 8 {
		m.doSupernovaPerfect(g, playerID, playerName, hitCount)
	} else {
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckySupernova,
			Payload: protocol.LuckySupernovaPayload{
				Event:      "supernova_end",
				PlayerID:   playerID,
				PlayerName: playerName,
				HitCount:   hitCount,
			},
		})
	}
}

func (m *luckySupernovaManager) doSupernovaPerfect(g *Game, playerID, playerName string, hitCount int) {
	m.mu.Lock()
	m.perfectBoost = &supernovaPerfectBoost{
		mult:      5.5,
		expiresAt: time.Now().Add(12 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckySupernova,
		Payload: protocol.LuckySupernovaPayload{
			Event:      "supernova_perfect",
			PlayerID:   playerID,
			PlayerName: playerName,
			HitCount:   hitCount,
			BoostMult:  5.5,
			BoostSec:   12,
		},
	})

	g.sendAnnounce("💥 超新星完美！"+playerName+" 命中 "+fmt.Sprintf("%d", hitCount)+" 個！全服 ×5.5 加成 12 秒！", "critical", "#FF6B35")

	time.Sleep(12 * time.Second)
	m.mu.Lock()
	m.perfectBoost = nil
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckySupernova,
		Payload: protocol.LuckySupernovaPayload{
			Event:      "supernova_perfect_end",
			PlayerID:   playerID,
			PlayerName: playerName,
		},
	})
}
