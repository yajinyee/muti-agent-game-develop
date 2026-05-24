// lucky_crab_torpedo_handler.go — T107 幸運螃蟹魚雷系統
// server-event-agent 負責維護
// 業界依據：Jackpot Fishing Jili「Crab Torpedoes — explosion AOE damage」
// 設計：擊破 T107 後，在場上隨機 3 個位置依序觸發 AOE 爆炸（r=150px，HP -40%）
// 每次爆炸命中至少 1 個目標 → 觸發玩家獲得 ×1.8 累積倍率（最高 ×5.4）
// 個人冷卻 18 秒；全服冷卻 30 秒
package game

import (
	"log"
	"math"
	"math/rand"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyCrabTorpedoManager struct {
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
}

func newLuckyCrabTorpedoManager() *luckyCrabTorpedoManager {
	return &luckyCrabTorpedoManager{
		playerCooldowns: make(map[string]time.Time),
	}
}

func isLuckyCrabTorpedoFish(defID string) bool {
	return defID == "T107"
}

func (m *luckyCrabTorpedoManager) canTrigger(playerID string) bool {
	now := time.Now()
	if now.Before(m.globalCooldown) {
		return false
	}
	if cd, ok := m.playerCooldowns[playerID]; ok {
		if now.Before(cd) {
			return false
		}
	}
	return true
}

func (g *Game) tryLuckyCrabTorpedo(playerID string, killerName string) {
	m := g.luckyCrabTorpedo
	if !m.canTrigger(playerID) {
		return
	}

	now := time.Now()
	m.playerCooldowns[playerID] = now.Add(18 * time.Second)
	m.globalCooldown = now.Add(30 * time.Second)

	g.hub.Broadcast(protocol.MsgLuckyCrabTorpedo, protocol.LuckyCrabTorpedoPayload{
		Event:       "trigger",
		TriggerID:   playerID,
		TriggerName: killerName,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "🦀 " + killerName + " 發射螃蟹魚雷！",
		Priority: "high",
		Color:    "#FF6B35",
	})

	go g.runCrabTorpedoes(playerID, killerName)
}

func (g *Game) runCrabTorpedoes(playerID string, killerName string) {
	time.Sleep(300 * time.Millisecond)

	g.mu.RLock()
	p, ok := g.players[playerID]
	betCost := 1
	if ok {
		betCost = p.GetBetDef().BetCost
	}
	g.mu.RUnlock()

	// 3 次爆炸，每次隨機位置
	const explosionRadius = 150.0
	const hpDamageRatio = 0.40 // HP -40%
	const rewardPerHit = 1.8   // 每次命中 ×1.8

	totalReward := 0
	accumMult := 1.0

	for i := 0; i < 3; i++ {
		time.Sleep(600 * time.Millisecond)

		// 隨機爆炸位置（避開邊緣）
		ex := 200.0 + rand.Float64()*(GameWidth-400)
		ey := 100.0 + rand.Float64()*(GameHeight-200)

		g.mu.Lock()
		hitTargets := []string{}
		for _, t := range g.targets {
			if t.Def.Type == "boss" {
				continue
			}
			dx := t.X - ex
			dy := t.Y - ey
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= explosionRadius {
				damage := int(float64(t.MaxHP) * hpDamageRatio)
				t.HP -= damage
				if t.HP < 1 {
					t.HP = 1
				}
				g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
					InstanceID: t.InstanceID,
					HP:         t.HP,
					MaxHP:      t.MaxHP,
					X:          t.X,
					Y:          t.Y,
				})
				hitTargets = append(hitTargets, t.InstanceID)
			}
		}

		if len(hitTargets) > 0 {
			accumMult += rewardPerHit
			reward := int(float64(betCost) * rewardPerHit)
			totalReward += reward
		}

		g.hub.Broadcast(protocol.MsgLuckyCrabTorpedo, protocol.LuckyCrabTorpedoPayload{
			Event:       "explosion",
			TriggerID:   playerID,
			TriggerName: killerName,
			ExplosionX:  ex,
			ExplosionY:  ey,
			HitTargets:  hitTargets,
			ExplosionNo: i + 1,
			TotalReward: totalReward,
		})
		g.mu.Unlock()
	}

	// 結算
	g.mu.Lock()
	if p2, ok2 := g.players[playerID]; ok2 {
		p2.AddCoins(totalReward)
		g.sendPlayerUpdate(playerID)
	}
	g.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyCrabTorpedo, protocol.LuckyCrabTorpedoPayload{
		Event:       "settle",
		TriggerID:   playerID,
		TriggerName: killerName,
		TotalReward: totalReward,
	})

	log.Printf("[CrabTorpedo] Player %s: mult=%.1f, reward=%d", playerID, accumMult, totalReward)
}
