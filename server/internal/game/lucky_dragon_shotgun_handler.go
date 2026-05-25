// lucky_dragon_shotgun_handler.go — T117 幸運龍力散彈魚
// server-event-agent 負責維護
// 業界依據：Jili Games 2026「Dragon Power Shotgun — 8-direction spread attack, each direction HP -40%」
package game

import (
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyDragonShotgunManager struct {
	mu                sync.Mutex
	personalCooldowns map[string]time.Time
	globalCooldown    time.Time
}

func newLuckyDragonShotgunManager() *luckyDragonShotgunManager {
	return &luckyDragonShotgunManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

func isLuckyDragonShotgunFish(defID string) bool {
	return defID == "T117"
}

func (m *luckyDragonShotgunManager) tryLuckyDragonShotgun(g *Game, playerID, playerName string) {
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
	// 個人冷卻 20 秒，全服冷卻 32 秒
	m.personalCooldowns[playerID] = now.Add(20 * time.Second)
	m.globalCooldown = now.Add(32 * time.Second)
	m.mu.Unlock()

	// 廣播觸發
	g.hub.Broadcast(protocol.MsgLuckyDragonShotgun, protocol.LuckyDragonShotgunPayload{
		Event:       "trigger",
		TriggerID:   playerID,
		TriggerName: playerName,
	})

	// 8 方向散彈攻擊
	go func() {
		totalHits := 0
		totalReward := 0

		for dir := 0; dir < 8; dir++ {
			time.Sleep(150 * time.Millisecond)

			g.mu.Lock()
			// 找方向扇形範圍內的目標（每方向最多 2 個）
			hitTargets := g.findTargetsInDirection(dir, 2)
			hitIDs := make([]string, 0, len(hitTargets))
			dirReward := 0

			p, pOk := g.players[playerID]
			for _, t := range hitTargets {
				// HP -40%
				dmg := t.MaxHP * 40 / 100
				t.HP -= dmg
				if t.HP < 1 {
					t.HP = 1
				}
				hitIDs = append(hitIDs, t.InstanceID)
				totalHits++

				// 廣播 HP 更新
				g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
					InstanceID: t.InstanceID,
					HP:         t.HP,
					MaxHP:      t.MaxHP,
					X:          t.X,
					Y:          t.Y,
				})

				// 若 HP <= 0 則擊破
				if t.HP <= 1 && pOk {
					bet := p.GetBetDef()
					reward := int(float64(bet.BetCost) * t.Multiplier * 0.5)
					p.AddCoins(reward)
					dirReward += reward
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

			if len(hitIDs) > 0 {
				g.hub.Broadcast(protocol.MsgLuckyDragonShotgun, protocol.LuckyDragonShotgunPayload{
					Event:       "shotgun_fire",
					TriggerID:   playerID,
					TriggerName: playerName,
					Direction:   dir,
					HitTargets:  hitIDs,
					TotalHits:   totalHits,
				})
			}
		}

		// 結算
		g.hub.Broadcast(protocol.MsgLuckyDragonShotgun, protocol.LuckyDragonShotgunPayload{
			Event:       "settle",
			TriggerID:   playerID,
			TriggerName: playerName,
			TotalHits:   totalHits,
			TotalReward: totalReward,
		})

		g.sendPlayerUpdate(playerID)
	}()
}

// findTargetsInDirection 找指定方向扇形範圍內的目標（最多 maxHit 個）
func (g *Game) findTargetsInDirection(direction int, maxHit int) []*Target {
	// 8 方向角度（度）
	angles := []float64{0, 45, 90, 135, 180, 225, 270, 315}
	_ = angles[direction]

	// 簡化：隨機選取場上目標（實際應依方向角度篩選）
	result := make([]*Target, 0, maxHit)
	for _, t := range g.targets {
		if len(result) >= maxHit {
			break
		}
		// 50% 機率命中（模擬散彈命中率）
		if rand.Float64() < 0.5 {
			result = append(result, t)
		}
	}
	return result
}
