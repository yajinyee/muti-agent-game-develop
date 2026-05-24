// lucky_thunder_lobster_handler.go — T110 幸運雷霆龍蝦系統
// server-event-agent 負責維護
// 業界依據：Royal Fishing Jili「Thunderbolt Lobster — 15 seconds free play + automatic shooting」
// 設計：擊破 T110 後，觸發玩家進入「雷霆模式」（15 秒）
// 雷霆模式期間：每 0.5 秒自動對場上最高倍率目標發射一次（免費，不扣 bet_cost）
// 每次命中 → 觸發玩家獲得 ×1.0 倍率獎勵（等同 bet_cost）
// 雷霆模式結束時廣播結算
// 個人冷卻 22 秒；全服冷卻 38 秒
package game

import (
	"log"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyThunderLobsterManager struct {
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
	activeSessions  map[string]bool // playerID -> isActive
}

func newLuckyThunderLobsterManager() *luckyThunderLobsterManager {
	return &luckyThunderLobsterManager{
		playerCooldowns: make(map[string]time.Time),
		activeSessions:  make(map[string]bool),
	}
}

func isLuckyThunderLobsterFish(defID string) bool {
	return defID == "T110"
}

func (m *luckyThunderLobsterManager) canTrigger(playerID string) bool {
	now := time.Now()
	if now.Before(m.globalCooldown) {
		return false
	}
	if m.activeSessions[playerID] {
		return false
	}
	if cd, ok := m.playerCooldowns[playerID]; ok {
		if now.Before(cd) {
			return false
		}
	}
	return true
}

func (g *Game) tryLuckyThunderLobster(playerID string, killerName string) {
	m := g.luckyThunderLobster
	if !m.canTrigger(playerID) {
		return
	}

	now := time.Now()
	m.playerCooldowns[playerID] = now.Add(22 * time.Second)
	m.globalCooldown = now.Add(38 * time.Second)
	m.activeSessions[playerID] = true

	g.hub.Broadcast(protocol.MsgLuckyThunderLobster, protocol.LuckyThunderLobsterPayload{
		Event:       "trigger",
		TriggerID:   playerID,
		TriggerName: killerName,
		TimeLeft:    15.0,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "🦞⚡ " + killerName + " 觸發雷霆龍蝦！15 秒免費自動射擊！",
		Priority: "high",
		Color:    "#FF4500",
	})

	go g.runThunderLobster(playerID, killerName)
}

func (g *Game) runThunderLobster(playerID string, killerName string) {
	const duration = 15 * time.Second
	const fireInterval = 500 * time.Millisecond
	const totalShots = int(duration / fireInterval) // 30 shots

	g.mu.RLock()
	p, ok := g.players[playerID]
	betCost := 1
	if ok {
		betCost = p.GetBetDef().BetCost
	}
	g.mu.RUnlock()

	killCount := 0
	totalReward := 0
	startTime := time.Now()

	for i := 0; i < totalShots; i++ {
		time.Sleep(fireInterval)

		timeLeft := duration - time.Since(startTime)
		if timeLeft <= 0 {
			break
		}

		g.mu.Lock()
		// 找最高倍率目標
		var bestTarget *Target
		bestMult := 0.0
		for _, t := range g.targets {
			if t.Def.Type == "boss" {
				continue
			}
			if t.Multiplier > bestMult {
				bestMult = t.Multiplier
				bestTarget = t
			}
		}

		if bestTarget != nil {
			// 免費射擊：直接嘗試擊破
			isKill := bestTarget.TryKill(betCost)
			if isKill {
				reward := int(float64(betCost) * bestTarget.Multiplier)
				totalReward += reward
				killCount++

				if p2, ok2 := g.players[playerID]; ok2 {
					p2.AddCoins(reward)
				}

				// 廣播擊破
				g.hub.Broadcast(protocol.MsgTargetKill, protocol.TargetKillPayload{
					InstanceID: bestTarget.InstanceID,
					DefID:      bestTarget.Def.ID,
					Multiplier: bestTarget.Multiplier,
					Reward:     reward,
					LaborGain:  bestTarget.Def.LaborGain,
					KillerID:   playerID,
				})
				delete(g.targets, bestTarget.InstanceID)
			} else {
				// 未擊破，更新 HP
				bestTarget.HP = max(1, bestTarget.HP-betCost/5)
				g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
					InstanceID: bestTarget.InstanceID,
					HP:         bestTarget.HP,
					MaxHP:      bestTarget.MaxHP,
					X:          bestTarget.X,
					Y:          bestTarget.Y,
				})
			}
		}

		// 每 5 次廣播一次進度
		if i%5 == 0 {
			g.hub.Broadcast(protocol.MsgLuckyThunderLobster, protocol.LuckyThunderLobsterPayload{
				Event:       "auto_fire",
				TriggerID:   playerID,
				TriggerName: killerName,
				TimeLeft:    timeLeft.Seconds(),
				KillCount:   killCount,
				TotalReward: totalReward,
			})
		}
		g.mu.Unlock()
	}

	// 結算
	g.mu.Lock()
	if p2, ok2 := g.players[playerID]; ok2 {
		g.sendPlayerUpdate(playerID)
		_ = p2
	}
	g.luckyThunderLobster.activeSessions[playerID] = false
	g.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyThunderLobster, protocol.LuckyThunderLobsterPayload{
		Event:       "end",
		TriggerID:   playerID,
		TriggerName: killerName,
		TimeLeft:    0,
		KillCount:   killCount,
		TotalReward: totalReward,
	})

	log.Printf("[ThunderLobster] Player %s: kills=%d, reward=%d", playerID, killCount, totalReward)
}
