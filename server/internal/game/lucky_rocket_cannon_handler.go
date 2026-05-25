// lucky_rocket_cannon_handler.go — T118 幸運火箭砲魚
// server-event-agent 負責維護
// 業界依據：Jili Games 2026「Rocket Cannon — 3 rockets, each AOE r=200px HP -50%」
package game

import (
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyRocketCannonManager struct {
	mu                sync.Mutex
	personalCooldowns map[string]time.Time
	globalCooldown    time.Time
}

func newLuckyRocketCannonManager() *luckyRocketCannonManager {
	return &luckyRocketCannonManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

func isLuckyRocketCannonFish(defID string) bool {
	return defID == "T118"
}

func (m *luckyRocketCannonManager) tryLuckyRocketCannon(g *Game, playerID, playerName string) {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCooldown) {
		m.mu.Unlock()
		return
	}
	if cd, ok := m.personalCooldowns[playerID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	// 個人冷卻 22 秒，全服冷卻 36 秒
	m.personalCooldowns[playerID] = now.Add(22 * time.Second)
	m.globalCooldown = now.Add(36 * time.Second)
	m.mu.Unlock()

	// 廣播觸發
	g.hub.Broadcast(protocol.MsgLuckyRocketCannon, protocol.LuckyRocketCannonPayload{
		Event:       "trigger",
		TriggerID:   playerID,
		TriggerName: playerName,
	})

	// 3 枚火箭砲依序發射
	go func() {
		totalReward := 0

		for rocketNo := 1; rocketNo <= 3; rocketNo++ {
			time.Sleep(700 * time.Millisecond)

			// 隨機選擇爆炸位置（場地中心附近）
			explodeX := 200.0 + rand.Float64()*880.0
			explodeY := 100.0 + rand.Float64()*520.0

			// 廣播火箭發射
			g.hub.Broadcast(protocol.MsgLuckyRocketCannon, protocol.LuckyRocketCannonPayload{
				Event:       "rocket_launch",
				TriggerID:   playerID,
				TriggerName: playerName,
				RocketNo:    rocketNo,
				ExplodeX:    explodeX,
				ExplodeY:    explodeY,
			})

			time.Sleep(400 * time.Millisecond)

			// AOE 傷害 r=200px HP -50%
			g.mu.Lock()
			hitIDs := make([]string, 0)
			p, pOk := g.players[playerID]

			for _, t := range g.targets {
				dx := t.X - explodeX
				dy := t.Y - explodeY
				dist := dx*dx + dy*dy
				if dist <= 200*200 {
					dmg := t.MaxHP * 50 / 100
					t.HP -= dmg
					if t.HP < 1 {
						t.HP = 1
					}
					hitIDs = append(hitIDs, t.InstanceID)

					g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
						InstanceID: t.InstanceID,
						HP:         t.HP,
						MaxHP:      t.MaxHP,
						X:          t.X,
						Y:          t.Y,
					})

					if t.HP <= 1 && pOk {
						bet := p.GetBetDef()
						reward := int(float64(bet.BetCost) * t.Multiplier * 0.6)
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
			}
			g.mu.Unlock()

			// 廣播爆炸結果
			g.hub.Broadcast(protocol.MsgLuckyRocketCannon, protocol.LuckyRocketCannonPayload{
				Event:       "rocket_explode",
				TriggerID:   playerID,
				TriggerName: playerName,
				RocketNo:    rocketNo,
				ExplodeX:    explodeX,
				ExplodeY:    explodeY,
				HitTargets:  hitIDs,
			})
		}

		// 結算
		g.hub.Broadcast(protocol.MsgLuckyRocketCannon, protocol.LuckyRocketCannonPayload{
			Event:       "settle",
			TriggerID:   playerID,
			TriggerName: playerName,
			TotalReward: totalReward,
		})

		g.sendPlayerUpdate(playerID)
	}()
}
