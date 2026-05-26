// Package game — T138 幸運傳說龍魚 handler
// server-event-agent 負責維護
// 業界依據：Royal Fishing Jili「Legend Dragon 120-200x from 20x base multiplier」
// 設計：擊破後觸發「傳說龍降臨」，龍在場上隨機位置出現（持續 15 秒）；
//       龍每 3 秒噴火一次，噴火命中範圍內所有目標 HP -35%；
//       噴火 4 次全部命中 ≥ 3 個目標 → 「傳說龍怒」：全服 ×4.0 加成 10 秒；
//       個人冷卻 35 秒；全服冷卻 55 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyLegendDragonManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	dragonRage   *legendDragonRage
}

type legendDragonRage struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyLegendDragonManager() *luckyLegendDragonManager {
	return &luckyLegendDragonManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyLegendDragonFish(defID string) bool {
	return defID == "T138"
}

func (m *luckyLegendDragonManager) getLegendDragonRageMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.dragonRage != nil && time.Now().Before(m.dragonRage.expiresAt) {
		return m.dragonRage.mult
	}
	return 1.0
}

func (g *Game) tryLuckyLegendDragonFish(playerID, playerName string) {
	m := g.luckyLegendDragon
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
	m.personalCD[playerID] = now.Add(35 * time.Second)
	m.globalCD = now.Add(55 * time.Second)
	m.mu.Unlock()

	log.Printf("[LuckyLegendDragon] Triggered by %s", playerName)

	g.hub.Broadcast(protocol.MsgLuckyLegendDragon, protocol.LuckyLegendDragonPayload{
		Event:      "dragon_appear",
		PlayerID:   playerID,
		PlayerName: playerName,
		Duration:   15.0,
	})

	go g.runLegendDragonBreath(playerID, playerName)
}

func (g *Game) runLegendDragonBreath(playerID, playerName string) {
	perfectBreaths := 0

	for breath := 1; breath <= 4; breath++ {
		time.Sleep(3 * time.Second)

		g.mu.Lock()
		hitCount := 0
		for _, t := range g.targets {
			if t.HP <= 0 {
				continue
			}
			dmg := int(float64(t.MaxHP) * 0.35)
			t.HP -= dmg
			if t.HP < 0 {
				t.HP = 0
			}
			hitCount++
			g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
				InstanceID: t.InstanceID,
				HP:         t.HP,
				MaxHP:      t.MaxHP,
				X:          float64(t.X),
				Y:          float64(t.Y),
			})
		}
		g.mu.Unlock()

		if hitCount >= 3 {
			perfectBreaths++
		}

		log.Printf("[LuckyLegendDragon] Breath %d: hit=%d", breath, hitCount)

		g.hub.Broadcast(protocol.MsgLuckyLegendDragon, protocol.LuckyLegendDragonPayload{
			Event:          "dragon_breath",
			PlayerID:       playerID,
			PlayerName:     playerName,
			BreathNum:      breath,
			HitCount:       hitCount,
			PerfectBreaths: perfectBreaths,
		})
	}

	// 傳說龍怒：4 次噴火全部命中 ≥ 3 個
	if perfectBreaths >= 4 {
		g.doLegendDragonRage(playerID, playerName, perfectBreaths)
	} else {
		g.hub.Broadcast(protocol.MsgLuckyLegendDragon, protocol.LuckyLegendDragonPayload{
			Event:          "dragon_leave",
			PlayerID:       playerID,
			PlayerName:     playerName,
			PerfectBreaths: perfectBreaths,
		})
	}
}

func (g *Game) doLegendDragonRage(playerID, playerName string, perfectBreaths int) {
	m := g.luckyLegendDragon
	m.mu.Lock()
	m.dragonRage = &legendDragonRage{
		mult:      4.0,
		expiresAt: time.Now().Add(10 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyLegendDragon] RAGE! %s perfect_breaths=%d → global ×4.0 for 10s", playerName, perfectBreaths)

	g.hub.Broadcast(protocol.MsgLuckyLegendDragon, protocol.LuckyLegendDragonPayload{
		Event:          "dragon_rage",
		PlayerID:       playerID,
		PlayerName:     playerName,
		PerfectBreaths: perfectBreaths,
		BoostMult:      4.0,
		BoostSec:       10,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  fmt.Sprintf("🐲 傳說龍怒！%s 完美噴火 %d 次！全服 ×4.0 加成 10 秒！", playerName, perfectBreaths),
		Priority: "critical",
		Color:    "#FF6B00",
	})

	go func() {
		time.Sleep(10 * time.Second)
		m.mu.Lock()
		m.dragonRage = nil
		m.mu.Unlock()
		g.hub.Broadcast(protocol.MsgLuckyLegendDragon, protocol.LuckyLegendDragonPayload{
			Event: "dragon_rage_end",
		})
	}()
}
