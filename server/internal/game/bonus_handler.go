// bonus_handler.go — Bonus Game 相關邏輯（DAY-058）
// 從 game.go 拆分：triggerBonusReady, startBonusGame, spawnBonusTargets,
//                  pickBonusTarget, endBonusGame, updateBonusGame, handleBonusClick
package game

import (
	"log"
	"math/rand"
	"time"

	"digital-twin/server/internal/analytics"
	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/combat"
	"digital-twin/server/internal/game/guild"
	"digital-twin/server/internal/game/mission"
	"digital-twin/server/internal/game/state"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/game/tournament"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"

	"github.com/google/uuid"
)

// triggerBonusReady 觸發 Bonus Ready（防止 90 秒內重複觸發）
func (g *Game) triggerBonusReady() {
	g.mu.RLock()
	currentState := g.State
	g.mu.RUnlock()

	if currentState != state.StateNormalPlay {
		return
	}

	// 防止 Bonus 觸發過於頻繁（至少間隔 90 秒）
	g.mu.Lock()
	if !g.lastBonusAt.IsZero() && time.Since(g.lastBonusAt).Seconds() < 90 {
		g.mu.Unlock()
		return
	}
	g.lastBonusAt = time.Now()
	g.mu.Unlock()

	g.transitionState(state.StateBonusReady)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBonusEvent,
		Payload: ws.BonusEventPayload{Event: "ready"},
	})

	// 3秒後自動進入 Bonus Game
	g.safeAfterFunc(3*time.Second, func() {
		g.startBonusGame()
	})
}

// startBonusGame 開始 Bonus Game
func (g *Game) startBonusGame() {
	g.mu.Lock()
	g.bonusStartedAt = time.Now()
	// 記錄所有玩家的進場 Bet
	for id, p := range g.Players {
		g.bonusEntryBet[id] = data.GetBetDef(p.BetLevel).BetCost
		g.bonusScores[id] = 0
	}
	// 清除一般目標，生成 Bonus 目標
	g.Targets = make(map[string]*target.Target)
	g.mu.Unlock()

	g.transitionState(state.StateBonusGame)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBonusEvent,
		Payload: ws.BonusEventPayload{
			Event:    "start",
			TimeLeft: data.BonusDuration,
		},
	})

	// 埋點：Bonus 開始
	if tracker := analytics.Get(); tracker != nil {
		tracker.Track(analytics.EventBonusStart, "system", map[string]interface{}{
			"duration": data.BonusDuration,
		})
	}

	// 生成 Bonus 目標
	g.spawnBonusTargets()
}

// spawnBonusTargets 生成 20 個 Bonus 目標
func (g *Game) spawnBonusTargets() {
	for i := 0; i < 20; i++ {
		def := g.pickBonusTarget()
		instanceID := uuid.New().String()
		x := 100.0 + rand.Float64()*1000
		y := 100.0 + rand.Float64()*500

		// 用 HP 欄位存 ClickScore（Bonus 目標特殊處理）
		bonusDef := &data.TargetDef{
			ID:       def.ID,
			Name:     def.Name,
			Type:     data.TargetTypeBonus,
			HP:       def.ClickScore,
			Lifetime: data.BonusDuration,
		}
		t := target.NewTarget(instanceID, bonusDef, x, y)

		g.mu.Lock()
		g.Targets[instanceID] = t
		g.mu.Unlock()

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetSpawn,
			Payload: ws.TargetSpawnPayload{
				InstanceID: instanceID,
				DefID:      def.ID,
				Name:       def.Name,
				Type:       "bonus",
				X:          x,
				Y:          y,
				HP:         def.ClickScore,
				MaxHP:      def.ClickScore,
				Behavior:   def.SpecialEffect,
			},
		})
	}
}

// pickBonusTarget 依權重隨機選取 Bonus 目標定義
func (g *Game) pickBonusTarget() *data.BonusTargetDef {
	total := 0
	for _, d := range data.BonusTargets {
		total += d.SpawnWeight
	}
	r := rand.Intn(total)
	cumulative := 0
	for _, d := range data.BonusTargets {
		cumulative += d.SpawnWeight
		if r < cumulative {
			return d
		}
	}
	return data.BonusTargets[0]
}

// endBonusGame 結束 Bonus Game，計算並發放獎勵
func (g *Game) endBonusGame() {
	g.mu.Lock()
	scores := make(map[string]int)
	entryBets := make(map[string]int)
	for id, score := range g.bonusScores {
		scores[id] = score
	}
	for id, bet := range g.bonusEntryBet {
		entryBets[id] = bet
	}
	g.Targets = make(map[string]*target.Target)
	g.mu.Unlock()

	// 計算每個玩家的獎勵
	for playerID, score := range scores {
		g.mu.RLock()
		p := g.Players[playerID]
		g.mu.RUnlock()

		if p == nil {
			continue
		}

		entryBet := entryBets[playerID]
		reward, multiplier := combat.CalcBonusReward(entryBet, score)
		p.AddReward(reward)
		p.ResetLaborValue()
		// 玩家統計：記錄 Bonus（DAY-096）
		g.notifyStatsBonus(p, reward)
		// 異常偵測：記錄 Bonus 觸發（DAY-105）
		if alert := g.AntiCheat.RecordBonus(playerID); alert != nil {
			log.Printf("[AntiCheat] Bonus Abuse Alert for player %s: %s", playerID, alert.Message)
		}

		g.Hub.Send(playerID, &ws.Message{
			Type: ws.MsgBonusEvent,
			Payload: ws.BonusEventPayload{
				Event:      "end",
				Score:      score,
				Multiplier: multiplier,
				Reward:     reward,
			},
		})

		g.Hub.Send(playerID, &ws.Message{
			Type: ws.MsgReward,
			Payload: ws.RewardPayload{
				Source:     "bonus",
				Amount:     reward,
				Multiplier: multiplier,
				NewBalance: p.Coins,
			},
		})

		// 埋點：Bonus 結束（每個玩家的獎勵）
		if tracker := analytics.Get(); tracker != nil {
			tracker.Track(analytics.EventBonusEnd, playerID, map[string]interface{}{
				"score":      score,
				"multiplier": multiplier,
				"reward":     reward,
				"entry_bet":  entryBet,
			})
			tracker.Track(analytics.EventReward, playerID, map[string]interface{}{
				"source":     "bonus",
				"amount":     reward,
				"multiplier": multiplier,
			})
		}

		// 任務進度：完成 Bonus Game（DAY-037）
		go func(pid string, r int) {
			g.updateMissionProgress(pid, mission.MissionPlayBonus, 1)
			g.updateMissionProgress(pid, mission.MissionEarnCoins, r)
		}(playerID, reward)

		// 週賽積分：完成 Bonus Game（DAY-066）
		g.tournamentMgr.AddPoints(playerID, p.DisplayName, tournament.PointBonus, 0)
		// 每日賽積分：完成 Bonus Game（DAY-093）
		g.notifyDailyTournamentBonus(p)
		// 賽季積分同步（DAY-072）：Bonus 完成 = 20 分
		newLevels := g.addSeasonPoints(playerID, 20)
		g.checkSeasonLevelNotify(p, newLevels)
		// 公會任務進度：賺取金幣（DAY-074）
		guildID := g.Guild.GetPlayerGuildID(playerID)
		if guildID != "" {
			completedTasks := g.Guild.UpdateTaskProgress(playerID, guild.TaskEarnCoins, reward)
			g.notifyGuildTaskComplete(guildID, completedTasks)
		}
		// 公會戰積分：Bonus 完成（DAY-076）
		go g.notifyGuildWarBonus(playerID)
	}

	g.transitionState(state.StateBonusResult)
	g.safeAfterFunc(3*time.Second, func() {
		g.transitionState(state.StateNormalPlay)
	})
}

// updateBonusGame 每 tick 更新 Bonus Game 狀態（廣播剩餘時間）
func (g *Game) updateBonusGame() {
	g.mu.RLock()
	bonusStart := g.bonusStartedAt
	g.mu.RUnlock()

	elapsed := time.Since(bonusStart).Seconds()
	timeLeft := data.BonusDuration - elapsed

	if timeLeft <= 0 {
		g.endBonusGame()
		return
	}

	// 廣播剩餘時間（每秒一次，避免過度廣播）
	now := time.Now()
	g.mu.Lock()
	shouldTick := now.Sub(g.lastBonusTickAt) >= time.Second
	if shouldTick {
		g.lastBonusTickAt = now
	}
	g.mu.Unlock()

	if shouldTick {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgBonusEvent,
			Payload: ws.BonusEventPayload{
				Event:    "tick",
				TimeLeft: timeLeft,
			},
		})
	}
}

// handleBonusClick 處理 Bonus Game 中的點擊（拔草）
func (g *Game) handleBonusClick(p *player.Player, msg *ws.Message) {
	g.mu.RLock()
	currentState := g.State
	g.mu.RUnlock()

	if currentState != state.StateBonusGame {
		return
	}

	var payload ws.BonusClickPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	g.mu.Lock()
	t := g.Targets[payload.TargetID]
	if t == nil || !t.IsAlive {
		g.mu.Unlock()
		return
	}

	defID := t.DefID
	score := t.Def.HP // HP 欄位存 ClickScore

	// BG002 硬雜草：需連點 2 次（規格書 29.3）
	// 用 HitCount 追蹤已點擊次數
	if defID == "BG002" {
		t.HitCount++
		if t.HitCount < 2 {
			// 第一次點擊：只廣播受擊效果，不消滅
			g.mu.Unlock()
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgTargetUpdate,
				Payload: ws.TargetUpdatePayload{
					InstanceID: payload.TargetID,
					HP:         2 - t.HitCount, // 剩餘點擊次數（視覺用）
					MaxHP:      2,
					X:          t.X,
					Y:          t.Y,
				},
			})
			return
		}
		// 第二次點擊：消滅
	}

	// 消滅目標
	t.IsAlive = false
	delete(g.Targets, payload.TargetID)
	g.bonusScores[p.ID] += score

	g.mu.Unlock()

	// 廣播 target_kill（讓 Client 播放拔草動畫）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTargetKill,
		Payload: ws.TargetKillPayload{
			InstanceID: payload.TargetID,
			DefID:      defID,
			Multiplier: 1,
			Reward:     score,
			LaborGain:  score,
			KillerID:   p.ID,
		},
	})

	// 特殊雜草效果
	switch defID {
	case "BG003": // 發光雜草：增加倍率（加分）
		g.mu.Lock()
		g.bonusScores[p.ID] += 5 // 額外加分
		g.mu.Unlock()
	case "BG004": // 金色雜草：觸發巨大金幣（大量加分，規格書 29.3）
		g.mu.Lock()
		g.bonusScores[p.ID] += 20 // 金色雜草本身 20 分 + 額外 10 分獎勵
		g.mu.Unlock()
		// 廣播金幣特效事件（Client 播放金幣雨動畫）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgBonusEvent,
			Payload: ws.BonusEventPayload{
				Event: "coin_shower",
			},
		})
	case "BG005": // 搗亂怪草：扣分
		g.mu.Lock()
		if g.bonusScores[p.ID] > 5 {
			g.bonusScores[p.ID] -= 5
		}
		g.mu.Unlock()
	}
}
