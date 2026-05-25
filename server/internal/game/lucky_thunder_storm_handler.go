// lucky_thunder_storm_handler.go — T124 幸運雷暴魚系統
// server-event-agent 負責維護
// 業界依據：Royal Fishing Jili「Thunder Storm — random lightning strikes across the field for 10 seconds」
// 設計：擊破 T124 後，觸發「雷暴」：10 秒內每 1.5 秒在隨機位置落下一道閃電
// 每道閃電命中半徑 80px 內所有目標 HP -30%
// 共 6-7 道閃電；每道命中至少 1 個目標 → 觸發玩家 ×1.2 累積倍率（最高 ×7.0）
// 若 6 道以上全部命中 → 「雷暴完美」：全服 ×2.3 加成 6 秒
// 個人冷卻 24 秒；全服冷卻 40 秒
package game

import (
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyThunderStormManager struct {
	mu              sync.Mutex
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
	stormBoost      *thunderStormBoost
}

type thunderStormBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyThunderStormManager() *luckyThunderStormManager {
	return &luckyThunderStormManager{
		playerCooldowns: make(map[string]time.Time),
	}
}

func isLuckyThunderStormFish(defID string) bool {
	return defID == "T124"
}

func (m *luckyThunderStormManager) getThunderStormMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.stormBoost != nil && time.Now().Before(m.stormBoost.expiresAt) {
		return m.stormBoost.mult
	}
	return 1.0
}

func (g *Game) tryLuckyThunderStorm(playerID string, killerName string) {
	m := g.luckyThunderStorm
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCooldown) {
		m.mu.Unlock()
		return
	}
	if cd, ok := m.playerCooldowns[playerID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	m.playerCooldowns[playerID] = now.Add(24 * time.Second)
	m.globalCooldown = now.Add(40 * time.Second)
	m.mu.Unlock()

	// 決定閃電數量（6-7 道）
	lightningCount := 6 + rand.Intn(2)

	g.hub.Broadcast(protocol.MsgLuckyThunderStorm, protocol.LuckyThunderStormPayload{
		Event:          "storm_start",
		TriggerID:      playerID,
		TriggerName:    killerName,
		LightningCount: lightningCount,
		Duration:       10.0,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "⛈️ " + killerName + " 召喚雷暴！" + string(rune('0'+lightningCount)) + " 道閃電！",
		Priority: "high",
		Color:    "#FFD700",
	})

	go g.runThunderStorm(playerID, killerName, lightningCount)
}

func (g *Game) runThunderStorm(playerID string, killerName string, lightningCount int) {
	accumMult := 1.0
	totalReward := 0
	hitStrikes := 0
	lightningRadius := 80.0

	for i := 0; i < lightningCount; i++ {
		time.Sleep(1500 * time.Millisecond)

		// 隨機閃電落點
		strikeX := float64(80 + rand.Intn(1120))
		strikeY := float64(80 + rand.Intn(560))

		g.mu.Lock()
		p, ok := g.players[playerID]
		if !ok {
			g.mu.Unlock()
			break
		}
		betCost := p.GetBetDef().BetCost

		hitTargets := []string{}
		for _, t := range g.targets {
			if t.Def.Type == "boss" {
				continue
			}
			dx := t.X - strikeX
			dy := t.Y - strikeY
			if math.Sqrt(dx*dx+dy*dy) <= lightningRadius {
				damage := t.MaxHP * 30 / 100
				t.HP -= damage
				if t.HP < 1 {
					t.HP = 1
				}
				hitTargets = append(hitTargets, t.InstanceID)
				g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
					InstanceID: t.InstanceID,
					HP:         t.HP,
					MaxHP:      t.MaxHP,
					X:          t.X,
					Y:          t.Y,
				})
			}
		}

		strikeReward := 0
		if len(hitTargets) > 0 {
			hitStrikes++
			accumMult += 1.2
			if accumMult > 7.0 {
				accumMult = 7.0
			}
			strikeReward = betCost * len(hitTargets) / 2
			totalReward += strikeReward
			p.AddCoins(strikeReward)
			g.sendPlayerUpdate(playerID)
		}
		g.mu.Unlock()

		g.hub.Broadcast(protocol.MsgLuckyThunderStorm, protocol.LuckyThunderStormPayload{
			Event:          "lightning_strike",
			TriggerID:      playerID,
			TriggerName:    killerName,
			StrikeX:        strikeX,
			StrikeY:        strikeY,
			HitTargets:     hitTargets,
			StrikeNo:       i + 1,
			AccumMult:      accumMult,
			TotalReward:    totalReward,
		})
	}

	// 完美雷暴判定（6 道以上全部命中）
	isPerfect := hitStrikes >= 6
	if isPerfect {
		m := g.luckyThunderStorm
		m.mu.Lock()
		m.stormBoost = &thunderStormBoost{
			mult:      2.3,
			expiresAt: time.Now().Add(6 * time.Second),
		}
		m.mu.Unlock()

		g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
			Message:  "⛈️✨ 雷暴完美！" + killerName + " 全服 ×2.3 加成 6 秒！",
			Priority: "high",
			Color:    "#FFD700",
		})
		g.hub.Broadcast(protocol.MsgLuckyThunderStorm, protocol.LuckyThunderStormPayload{
			Event:       "perfect_storm",
			TriggerID:   playerID,
			TriggerName: killerName,
			HitStrikes:  hitStrikes,
			TotalReward: totalReward,
		})

		go func() {
			time.Sleep(6 * time.Second)
			g.hub.Broadcast(protocol.MsgLuckyThunderStorm, protocol.LuckyThunderStormPayload{
				Event:       "perfect_end",
				TriggerID:   playerID,
				TriggerName: killerName,
			})
		}()
	}

	g.hub.Broadcast(protocol.MsgLuckyThunderStorm, protocol.LuckyThunderStormPayload{
		Event:       "storm_end",
		TriggerID:   playerID,
		TriggerName: killerName,
		HitStrikes:  hitStrikes,
		AccumMult:   accumMult,
		TotalReward: totalReward,
	})

	log.Printf("[ThunderStorm] Player %s: strikes=%d/%d, mult=%.1f, reward=%d, perfect=%v",
		playerID, hitStrikes, lightningCount, accumMult, totalReward, isPerfect)
}
