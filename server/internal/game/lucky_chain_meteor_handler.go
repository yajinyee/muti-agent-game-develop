// Package game — T129 幸運連鎖隕石魚 handler
// server-event-agent 負責維護
package game

import (
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyChainMeteorManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time // 個人冷卻
	globalCD     time.Time            // 全服冷卻
	perfectBoost *chainMeteorPerfectBoost
}

type chainMeteorPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyChainMeteorManager() *luckyChainMeteorManager {
	return &luckyChainMeteorManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyChainMeteorFish(defID string) bool {
	return defID == "T129"
}

func (m *luckyChainMeteorManager) getChainMeteorPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (g *Game) tryLuckyChainMeteor(playerID, playerName string) {
	m := g.luckyChainMeteor
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
	m.personalCD[playerID] = now.Add(26 * time.Second)
	m.globalCD = now.Add(42 * time.Second)
	m.mu.Unlock()

	log.Printf("[LuckyChainMeteor] Triggered by %s", playerName)
	g.hub.Broadcast(protocol.MsgLuckyChainMeteor, protocol.LuckyChainMeteorPayload{
		Event:      "meteor_start",
		PlayerID:   playerID,
		PlayerName: playerName,
		AOERadius:  150,
	})

	go g.runChainMeteors(playerID, playerName)
}

func (g *Game) runChainMeteors(playerID, playerName string) {
	aoeRadius := 150.0
	allHit := true

	for i := 1; i <= 5; i++ {
		time.Sleep(600 * time.Millisecond)

		g.mu.Lock()
		hitCount := g.applyChainMeteorDamage(aoeRadius)
		g.mu.Unlock()

		if hitCount == 0 {
			allHit = false
			g.hub.Broadcast(protocol.MsgLuckyChainMeteor, protocol.LuckyChainMeteorPayload{
				Event:       "meteor_miss",
				PlayerID:    playerID,
				PlayerName:  playerName,
				MeteorIndex: i,
				AOERadius:   aoeRadius,
				HitCount:    0,
			})
		} else {
			// 連鎖：AOE 半徑 +30px（最大 300px）
			if aoeRadius < 300 {
				aoeRadius += 30
			}
			g.hub.Broadcast(protocol.MsgLuckyChainMeteor, protocol.LuckyChainMeteorPayload{
				Event:       "meteor_hit",
				PlayerID:    playerID,
				PlayerName:  playerName,
				MeteorIndex: i,
				AOERadius:   aoeRadius,
				HitCount:    hitCount,
			})
		}
	}

	if allHit {
		g.doChainMeteorPerfect(playerID, playerName)
	}
}

func (g *Game) applyChainMeteorDamage(radius float64) int {
	hitCount := 0
	for _, t := range g.targets {
		// 隨機位置落下，對所有目標造成傷害（簡化：對所有目標造成 HP -40%）
		dmg := int(float64(t.HP) * 0.40)
		if dmg < 1 {
			dmg = 1
		}
		t.HP -= dmg
		if t.HP < 0 {
			t.HP = 0
		}
		hitCount++
		g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
			InstanceID: t.InstanceID,
			HP:         t.HP,
			MaxHP:      t.MaxHP,
			X:          t.X,
			Y:          t.Y,
		})
	}
	return hitCount
}

func (g *Game) doChainMeteorPerfect(playerID, playerName string) {
	m := g.luckyChainMeteor
	m.mu.Lock()
	expiresAt := time.Now().Add(7 * time.Second)
	m.perfectBoost = &chainMeteorPerfectBoost{
		mult:      2.5,
		expiresAt: expiresAt,
	}
	m.mu.Unlock()

	log.Printf("[LuckyChainMeteor] Perfect! ×2.5 boost for 7s")
	g.hub.Broadcast(protocol.MsgLuckyChainMeteor, protocol.LuckyChainMeteorPayload{
		Event:       "meteor_perfect",
		PlayerID:    playerID,
		PlayerName:  playerName,
		PerfectMult: 2.5,
		ExpiresAt:   expiresAt.UnixMilli(),
	})

	time.AfterFunc(7*time.Second, func() {
		m.mu.Lock()
		m.perfectBoost = nil
		m.mu.Unlock()
		g.hub.Broadcast(protocol.MsgLuckyChainMeteor, protocol.LuckyChainMeteorPayload{
			Event:      "meteor_perfect_end",
			PlayerID:   playerID,
			PlayerName: playerName,
		})
	})
}
