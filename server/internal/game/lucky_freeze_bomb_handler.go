// lucky_freeze_bomb_handler.go — T123 幸運冰凍炸彈魚系統
// server-event-agent 負責維護
// 業界依據：Fishing Fortune 2026「Ice Bomb — freezes all fish in a radius, then explodes for massive damage」
// 設計：擊破 T123 後，投擲「冰凍炸彈」到場地中心
// 第一階段（3秒）：半徑 300px 內所有目標凍結（停止移動，傷害 ×1.5）
// 第二階段（爆炸）：凍結結束時，半徑 300px 內所有目標 HP -60%
// 若爆炸命中 ≥ 5 個目標 → 「冰爆完美」：全服 ×2.2 加成 5 秒
// 個人冷卻 22 秒；全服冷卻 38 秒
package game

import (
	"log"
	"math"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyFreezeBombManager struct {
	mu              sync.Mutex
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
	activeSession   *freezeBombSession
	perfectBoost    *freezeBombPerfectBoost
}

type freezeBombSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	frozenTargets     []string // 被凍結的目標 instance_id
	expiresAt         time.Time
	settled           bool
}

type freezeBombPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyFreezeBombManager() *luckyFreezeBombManager {
	return &luckyFreezeBombManager{
		playerCooldowns: make(map[string]time.Time),
	}
}

func isLuckyFreezeBombFish(defID string) bool {
	return defID == "T123"
}

func (m *luckyFreezeBombManager) getFreezeBombPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (g *Game) tryLuckyFreezeBomb(playerID string, killerName string) {
	m := g.luckyFreezeBomb
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
	m.playerCooldowns[playerID] = now.Add(22 * time.Second)
	m.globalCooldown = now.Add(38 * time.Second)

	// 炸彈落點：場地中心
	bombX := float64(GameWidth / 2)
	bombY := float64(GameHeight / 2)
	freezeRadius := 300.0

	// 找半徑內的目標
	g.mu.RLock()
	var frozenTargets []string
	for id, t := range g.targets {
		if t.Def.Type == "boss" {
			continue
		}
		dx := t.X - bombX
		dy := t.Y - bombY
		if math.Sqrt(dx*dx+dy*dy) <= freezeRadius {
			frozenTargets = append(frozenTargets, id)
		}
	}
	g.mu.RUnlock()

	m.activeSession = &freezeBombSession{
		triggerPlayerID:   playerID,
		triggerPlayerName: killerName,
		frozenTargets:     frozenTargets,
		expiresAt:         now.Add(3 * time.Second),
	}
	m.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyFreezeBomb, protocol.LuckyFreezeBombPayload{
		Event:         "freeze_start",
		TriggerID:     playerID,
		TriggerName:   killerName,
		BombX:         bombX,
		BombY:         bombY,
		FreezeRadius:  freezeRadius,
		FrozenTargets: frozenTargets,
		Duration:      3.0,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "❄️💣 " + killerName + " 投擲冰凍炸彈！" + string(rune('0'+len(frozenTargets))) + " 個目標凍結！",
		Priority: "high",
		Color:    "#00E5FF",
	})

	// 3 秒後爆炸
	go func() {
		time.Sleep(3 * time.Second)
		g.doFreezeBombExplosion(playerID, killerName, bombX, bombY, freezeRadius)
	}()
}

func (g *Game) doFreezeBombExplosion(playerID string, killerName string, bombX float64, bombY float64, radius float64) {
	g.mu.Lock()
	p, ok := g.players[playerID]
	if !ok {
		g.mu.Unlock()
		return
	}
	betCost := p.GetBetDef().BetCost

	// 爆炸：半徑內所有目標 HP -60%
	hitCount := 0
	totalReward := 0
	for _, t := range g.targets {
		if t.Def.Type == "boss" {
			continue
		}
		dx := t.X - bombX
		dy := t.Y - bombY
		if math.Sqrt(dx*dx+dy*dy) <= radius {
			damage := t.MaxHP * 60 / 100
			t.HP -= damage
			if t.HP < 1 {
				t.HP = 1
			}
			hitCount++
			reward := betCost / 2
			totalReward += reward
			g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
				InstanceID: t.InstanceID,
				HP:         t.HP,
				MaxHP:      t.MaxHP,
				X:          t.X,
				Y:          t.Y,
			})
		}
	}

	if ok {
		p.AddCoins(totalReward)
		g.sendPlayerUpdate(playerID)
	}
	g.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyFreezeBomb, protocol.LuckyFreezeBombPayload{
		Event:       "bomb_explode",
		TriggerID:   playerID,
		TriggerName: killerName,
		HitCount:    hitCount,
		TotalReward: totalReward,
	})

	// 完美冰爆判定
	if hitCount >= 5 {
		m := g.luckyFreezeBomb
		m.mu.Lock()
		m.perfectBoost = &freezeBombPerfectBoost{
			mult:      2.2,
			expiresAt: time.Now().Add(5 * time.Second),
		}
		m.mu.Unlock()

		g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
			Message:  "❄️💥✨ 冰爆完美！" + killerName + " 命中 " + string(rune('0'+hitCount)) + " 個！全服 ×2.2 加成 5 秒！",
			Priority: "high",
			Color:    "#00E5FF",
		})
		g.hub.Broadcast(protocol.MsgLuckyFreezeBomb, protocol.LuckyFreezeBombPayload{
			Event:       "perfect_freeze",
			TriggerID:   playerID,
			TriggerName: killerName,
			HitCount:    hitCount,
			TotalReward: totalReward,
		})

		go func() {
			time.Sleep(5 * time.Second)
			g.hub.Broadcast(protocol.MsgLuckyFreezeBomb, protocol.LuckyFreezeBombPayload{
				Event:       "perfect_end",
				TriggerID:   playerID,
				TriggerName: killerName,
			})
		}()
	}

	log.Printf("[FreezeBomb] Player %s: hit=%d, reward=%d, perfect=%v",
		playerID, hitCount, totalReward, hitCount >= 5)
}
