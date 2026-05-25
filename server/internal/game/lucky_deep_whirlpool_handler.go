// lucky_deep_whirlpool_handler.go — T119 幸運深海漩渦魚
// server-event-agent 負責維護
// 業界依據：Jili Games 2026「Free Deep Sea Whirlpool — all fish HP -50%, 6 seconds continuous damage」
package game

import (
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyDeepWhirlpoolManager struct {
	mu                sync.Mutex
	personalCooldowns map[string]time.Time
	globalCooldown    time.Time
	isActive          bool
}

func newLuckyDeepWhirlpoolManager() *luckyDeepWhirlpoolManager {
	return &luckyDeepWhirlpoolManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

func isLuckyDeepWhirlpoolFish(defID string) bool {
	return defID == "T119"
}

func (m *luckyDeepWhirlpoolManager) tryLuckyDeepWhirlpool(g *Game, playerID, playerName string) {
	m.mu.Lock()
	now := time.Now()
	if m.isActive {
		m.mu.Unlock()
		return
	}
	if now.Before(m.globalCooldown) {
		m.mu.Unlock()
		return
	}
	if cd, ok := m.personalCooldowns[playerID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	// 個人冷卻 25 秒，全服冷卻 40 秒
	m.personalCooldowns[playerID] = now.Add(25 * time.Second)
	m.globalCooldown = now.Add(40 * time.Second)
	m.isActive = true
	m.mu.Unlock()

	// 廣播觸發
	g.hub.Broadcast(protocol.MsgLuckyDeepWhirlpool, protocol.LuckyDeepWhirlpoolPayload{
		Event:       "trigger",
		TriggerID:   playerID,
		TriggerName: playerName,
	})

	// 6 秒漩渦傷害（每秒 HP -8%，共 6 次）
	go func() {
		totalReward := 0

		for tick := 0; tick < 6; tick++ {
			time.Sleep(1 * time.Second)

			g.mu.Lock()
			hitCount := 0
			p, pOk := g.players[playerID]

			for _, t := range g.targets {
				dmg := t.MaxHP * 8 / 100
				if dmg < 1 {
					dmg = 1
				}
				t.HP -= dmg
				if t.HP < 1 {
					t.HP = 1
				}
				hitCount++

				g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
					InstanceID: t.InstanceID,
					HP:         t.HP,
					MaxHP:      t.MaxHP,
					X:          t.X,
					Y:          t.Y,
				})

				if t.HP <= 1 && pOk {
					bet := p.GetBetDef()
					reward := int(float64(bet.BetCost) * t.Multiplier * 0.4)
					p.AddCoins(reward)
					totalReward += reward
					delete(g.targets, t.InstanceID)
					g.hub.Broadcast(protocol.MsgTargetKill, protocol.TargetKillPayload{
						InstanceID: t.InstanceID,
						DefID:      t.Def.ID,
						Multiplier: t.Multiplier,
						Reward:     reward,
						KillerID:   playerID,
					})
				}
			}
			g.mu.Unlock()

			g.hub.Broadcast(protocol.MsgLuckyDeepWhirlpool, protocol.LuckyDeepWhirlpoolPayload{
				Event:       "whirlpool_damage",
				TriggerID:   playerID,
				TriggerName: playerName,
				HitCount:    hitCount,
			})
		}

		// 結算
		m.mu.Lock()
		m.isActive = false
		m.mu.Unlock()

		g.hub.Broadcast(protocol.MsgLuckyDeepWhirlpool, protocol.LuckyDeepWhirlpoolPayload{
			Event:       "settle",
			TriggerID:   playerID,
			TriggerName: playerName,
			TotalReward: totalReward,
		})

		g.sendPlayerUpdate(playerID)
	}()
}
