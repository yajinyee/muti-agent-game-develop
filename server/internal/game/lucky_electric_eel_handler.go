// Package game — T131 幸運電鰻魚 handler
// server-event-agent 負責維護
// 業界依據：Royal Fishing Jili「60x lightning eel creates chain reactions that jump between nearby fish
//           consecutively until targeting disengages, creating cascading capture sequences」
// 設計：擊破後持續放電 12 秒，每 1.5 秒電擊最近 3 條魚（HP -25%），
//       每次電擊命中 ≥2 條 → 連鎖加速（間隔縮短 0.1s，最短 0.5s）
//       12 秒內累積電擊 ≥ 8 次 → 「超級放電」：全服 ×2.5 加成 7 秒
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyElectricEelManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	superBoost   *eelSuperBoost
}

type eelSuperBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyElectricEelManager() *luckyElectricEelManager {
	return &luckyElectricEelManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyElectricEelFish(defID string) bool {
	return defID == "T131"
}

func (m *luckyElectricEelManager) getEelSuperMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.superBoost != nil && time.Now().Before(m.superBoost.expiresAt) {
		return m.superBoost.mult
	}
	return 1.0
}

func (g *Game) tryLuckyElectricEelFish(playerID, playerName string) {
	m := g.luckyElectricEel
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
	m.personalCD[playerID] = now.Add(22 * time.Second)
	m.globalCD = now.Add(38 * time.Second)
	m.mu.Unlock()

	log.Printf("[LuckyElectricEel] Triggered by %s", playerName)

	g.hub.Broadcast(protocol.MsgLuckyElectricEel, protocol.LuckyElectricEelPayload{
		Event:      "eel_start",
		PlayerID:   playerID,
		PlayerName: playerName,
		Duration:   12.0,
	})

	go g.runElectricEelShock(playerID, playerName)
}

func (g *Game) runElectricEelShock(playerID, playerName string) {
	interval := 1500 * time.Millisecond
	totalDuration := 12 * time.Second
	deadline := time.Now().Add(totalDuration)
	shockCount := 0

	for time.Now().Before(deadline) {
		time.Sleep(interval)
		if time.Now().After(deadline) {
			break
		}

		// 找最近 3 條目標
		g.mu.Lock()
		type targetDist struct {
			id   string
			dist float64
		}
		var candidates []targetDist
		for id, t := range g.targets {
			if t.HP <= 0 {
				continue
			}
			dist := math.Sqrt(float64(t.X*t.X + t.Y*t.Y))
			candidates = append(candidates, targetDist{id, dist})
		}
		// 隨機選 3 個（模擬電鰻攻擊最近的）
		rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
		maxHit := 3
		if len(candidates) < maxHit {
			maxHit = len(candidates)
		}
		hitCount := 0
		for i := 0; i < maxHit; i++ {
			t, ok := g.targets[candidates[i].id]
			if !ok || t.HP <= 0 {
				continue
			}
			dmg := int(float64(t.MaxHP) * 0.25)
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

		shockCount++

		// 連鎖加速：命中 ≥2 條，間隔縮短 0.1s（最短 0.5s）
		if hitCount >= 2 {
			interval -= 100 * time.Millisecond
			if interval < 500*time.Millisecond {
				interval = 500 * time.Millisecond
			}
		}

		g.hub.Broadcast(protocol.MsgLuckyElectricEel, protocol.LuckyElectricEelPayload{
			Event:      "eel_shock",
			PlayerID:   playerID,
			PlayerName: playerName,
			HitCount:   hitCount,
			ShockCount: shockCount,
			TimeLeft:   time.Until(deadline).Seconds(),
		})
	}

	// 結算：累積電擊 ≥ 8 次 → 超級放電
	if shockCount >= 8 {
		g.doEelSuperCharge(playerID, playerName, shockCount)
	} else {
		g.hub.Broadcast(protocol.MsgLuckyElectricEel, protocol.LuckyElectricEelPayload{
			Event:      "eel_end",
			PlayerID:   playerID,
			PlayerName: playerName,
			ShockCount: shockCount,
		})
	}
}

func (g *Game) doEelSuperCharge(playerID, playerName string, shockCount int) {
	m := g.luckyElectricEel
	m.mu.Lock()
	m.superBoost = &eelSuperBoost{
		mult:      2.5,
		expiresAt: time.Now().Add(7 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyElectricEel] SUPER CHARGE! %s shocks=%d", playerName, shockCount)

	g.hub.Broadcast(protocol.MsgLuckyElectricEel, protocol.LuckyElectricEelPayload{
		Event:      "eel_super",
		PlayerID:   playerID,
		PlayerName: playerName,
		ShockCount: shockCount,
		BoostMult:  2.5,
		BoostSec:   7,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  fmt.Sprintf("⚡ 超級放電！%s 電擊 %d 次！全服 ×2.5 加成 7 秒！", playerName, shockCount),
		Priority: "high",
		Color:    "#FFD700",
	})

	go func() {
		time.Sleep(7 * time.Second)
		g.hub.Broadcast(protocol.MsgLuckyElectricEel, protocol.LuckyElectricEelPayload{
			Event: "eel_super_end",
		})
	}()
}
