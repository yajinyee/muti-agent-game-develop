// lucky_vortex_handler.go — T108 幸運渦旋海葵系統
// server-event-agent 負責維護
// 業界依據：Jackpot Fishing Jili「Sea Anemone — whirlpool pulls fish to center」
// 設計：擊破 T108 後，全場渦旋 5 秒（所有目標 HP -30%），每秒廣播一次
// 渦旋結束時，場上所有目標 HP 再 -20%（渦旋爆炸）
// 個人冷卻 20 秒；全服冷卻 35 秒
package game

import (
	"log"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyVortexManager struct {
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
	isActive        bool
}

func newLuckyVortexManager() *luckyVortexManager {
	return &luckyVortexManager{
		playerCooldowns: make(map[string]time.Time),
	}
}

func isLuckyVortexFish(defID string) bool {
	return defID == "T108"
}

func (m *luckyVortexManager) canTrigger(playerID string) bool {
	now := time.Now()
	if m.isActive {
		return false
	}
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

func (g *Game) tryLuckyVortex(playerID string, killerName string) {
	m := g.luckyVortex
	if !m.canTrigger(playerID) {
		return
	}

	now := time.Now()
	m.playerCooldowns[playerID] = now.Add(20 * time.Second)
	m.globalCooldown = now.Add(35 * time.Second)
	m.isActive = true

	g.hub.Broadcast(protocol.MsgLuckyVortex, protocol.LuckyVortexPayload{
		Event:       "trigger",
		TriggerID:   playerID,
		TriggerName: killerName,
		TimeLeft:    5.0,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "🌀 " + killerName + " 召喚渦旋海葵！全場 HP -30%！",
		Priority: "high",
		Color:    "#7B2FBE",
	})

	go g.runVortex(playerID, killerName)
}

func (g *Game) runVortex(playerID string, killerName string) {
	g.mu.RLock()
	p, ok := g.players[playerID]
	betCost := 1
	if ok {
		betCost = p.GetBetDef().BetCost
	}
	g.mu.RUnlock()

	totalReward := 0

	// 5 秒渦旋，每秒傷害一次
	for tick := 0; tick < 5; tick++ {
		time.Sleep(1 * time.Second)

		g.mu.Lock()
		hitCount := 0
		for _, t := range g.targets {
			if t.Def.Type == "boss" {
				continue
			}
			damage := int(float64(t.MaxHP) * 0.06) // 每秒 6%，5 秒共 30%
			if damage < 1 {
				damage = 1
			}
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
			hitCount++
		}
		reward := hitCount * betCost / 5
		totalReward += reward

		g.hub.Broadcast(protocol.MsgLuckyVortex, protocol.LuckyVortexPayload{
			Event:       "pull",
			TriggerID:   playerID,
			TriggerName: killerName,
			TimeLeft:    float64(4 - tick),
			HitCount:    hitCount,
			TotalReward: totalReward,
		})
		g.mu.Unlock()
	}

	// 渦旋爆炸：所有目標 HP -20%
	time.Sleep(200 * time.Millisecond)
	g.mu.Lock()
	explosionHits := 0
	for _, t := range g.targets {
		if t.Def.Type == "boss" {
			continue
		}
		damage := int(float64(t.MaxHP) * 0.20)
		if damage < 1 {
			damage = 1
		}
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
		explosionHits++
	}
	explosionReward := explosionHits * betCost / 3
	totalReward += explosionReward

	if p2, ok2 := g.players[playerID]; ok2 {
		p2.AddCoins(totalReward)
		g.sendPlayerUpdate(playerID)
	}

	g.luckyVortex.isActive = false
	g.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyVortex, protocol.LuckyVortexPayload{
		Event:       "end",
		TriggerID:   playerID,
		TriggerName: killerName,
		TimeLeft:    0,
		HitCount:    explosionHits,
		TotalReward: totalReward,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "🌀 渦旋爆炸！全場 HP -20%！",
		Priority: "normal",
		Color:    "#C77DFF",
	})

	log.Printf("[Vortex] Player %s: reward=%d", playerID, totalReward)
}
