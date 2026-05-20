// boss_handler.go — BOSS 相關邏輯（DAY-058）
// 從 game.go 拆分：triggerBoss, spawnBoss, updateBossBattle, handleBossKill
package game

import (
	"log"
	"math/rand"
	"time"

	"digital-twin/server/internal/analytics"
	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/achievement"
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

// triggerBoss 觸發 BOSS（規格書 28.1：每 3-5 分鐘自動觸發）
func (g *Game) triggerBoss() {
	g.mu.RLock()
	currentState := g.State
	g.mu.RUnlock()

	if currentState != state.StateNormalPlay {
		return
	}

	// BOSS 警告
	g.transitionState(state.StateBossWarning)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBossEvent,
		Payload: ws.BossEventPayload{Event: "warning"},
	})

	g.safeAfterFunc(3*time.Second, func() {
		g.spawnBoss()
	})
}

// spawnBoss 生成 BOSS（HP 依玩家平均 bet 等級動態縮放）
func (g *Game) spawnBoss() {
	def := data.Targets["B001"]
	instanceID := uuid.New().String()

	// BOSS HP 依玩家平均 bet 等級縮放（DAY-044b）
	// 設計原則：玩家在 60 秒內有 ~50% 機率打死 BOSS
	g.mu.RLock()
	avgBetLevel := 5
	playerCount := len(g.Players)
	if playerCount > 0 {
		total := 0
		for _, p := range g.Players {
			total += p.BetLevel
		}
		avgBetLevel = total / playerCount
		if avgBetLevel < 1 {
			avgBetLevel = 1
		}
	}
	g.mu.RUnlock()

	betDef := data.GetBetDef(avgBetLevel)
	effectivePlayers := playerCount
	if effectivePlayers < 1 {
		effectivePlayers = 1
	}
	if effectivePlayers > 4 {
		effectivePlayers = 4 // 最多 4 人效果，避免 HP 過高
	}
	bossHP := int(betDef.FireRate * 60 * float64(betDef.BetCost) * 0.5 * float64(effectivePlayers))
	if bossHP < 100 {
		bossHP = 100
	}
	if bossHP > 10000 {
		bossHP = 10000
	}

	// 建立動態 HP 的 BOSS def（不修改原始 def）
	bossDef := *def
	bossDef.HP = bossHP

	t := target.NewTarget(instanceID, &bossDef, 1100, 360)

	g.mu.Lock()
	g.Targets[instanceID] = t
	g.bossInstanceID = instanceID
	g.bossSpawnedAt = time.Now()
	g.mu.Unlock()

	log.Printf("[Game] BOSS spawned: HP=%d (avgBetLV=%d, players=%d)", bossHP, avgBetLevel, playerCount)

	g.transitionState(state.StateBossBattle)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBossEvent,
		Payload: ws.BossEventPayload{
			Event:      "spawn",
			InstanceID: instanceID,
			HP:         bossHP,
			MaxHP:      bossHP,
		},
	})

	// 埋點：BOSS 生成
	if tracker := analytics.Get(); tracker != nil {
		tracker.Track(analytics.EventBossSpawn, "system", map[string]interface{}{
			"instance_id": instanceID,
			"boss_def":    "B001",
			"hp":          bossHP,
			"avg_bet_lv":  avgBetLevel,
		})
	}
}

// updateBossBattle 每 tick 更新 BOSS 戰狀態
func (g *Game) updateBossBattle() {
	g.mu.RLock()
	bossID := g.bossInstanceID
	bossSpawnedAt := g.bossSpawnedAt
	g.mu.RUnlock()

	if bossID == "" {
		return
	}

	// BOSS 超時
	elapsed := time.Since(bossSpawnedAt).Seconds()
	if elapsed >= data.BossDuration {
		g.mu.Lock()
		delete(g.Targets, bossID)
		g.bossInstanceID = ""
		g.mu.Unlock()

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgBossEvent,
			Payload: ws.BossEventPayload{
				Event: "timeout",
			},
		})
		g.transitionState(state.StateNormalPlay)
		return
	}

	// 規格書 9章：BOSS 期間 Max Targets = 8（不含 BOSS 本身）
	const MaxTargetsDuringBoss = 8
	g.mu.Lock()
	type targetWithTime struct {
		id        string
		spawnedAt time.Time
	}
	nonBossTargets := make([]targetWithTime, 0)
	for id, t := range g.Targets {
		if id != bossID && t.Def.Type != data.TargetTypeBoss {
			nonBossTargets = append(nonBossTargets, targetWithTime{id: id, spawnedAt: t.SpawnedAt})
		}
	}
	if len(nonBossTargets) > MaxTargetsDuringBoss {
		// 依生成時間排序（最舊的在前）
		for i := 0; i < len(nonBossTargets); i++ {
			for j := i + 1; j < len(nonBossTargets); j++ {
				if nonBossTargets[j].spawnedAt.Before(nonBossTargets[i].spawnedAt) {
					nonBossTargets[i], nonBossTargets[j] = nonBossTargets[j], nonBossTargets[i]
				}
			}
		}
		// 移除最舊的目標
		for i := 0; i < len(nonBossTargets)-MaxTargetsDuringBoss; i++ {
			delete(g.Targets, nonBossTargets[i].id)
		}
	}
	g.mu.Unlock()
}

// handleBossKill 處理 BOSS 被擊殺
func (g *Game) handleBossKill(p *player.Player, t *target.Target, result *combat.AttackResult) {
	g.mu.Lock()
	g.bossInstanceID = ""
	g.mu.Unlock()

	// 成就：首次擊敗 BOSS
	if u := p.TryUnlockAchievement(achievement.AchKillBoss); u != nil {
		g.sendAchievements(p.ID, []*achievement.AchievementUnlock{u})
	}

	// 任務進度：擊敗 BOSS（DAY-037）
	go g.updateMissionProgress(p.ID, mission.MissionKillBoss, 1)

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBossEvent,
		Payload: ws.BossEventPayload{
			Event:      "kill",
			InstanceID: t.InstanceID,
			Reward:     result.Reward,
			Multiplier: result.Multiplier,
		},
	})

	// 埋點：BOSS 擊敗
	if tracker := analytics.Get(); tracker != nil {
		tracker.Track(analytics.EventBossKill, p.ID, map[string]interface{}{
			"instance_id": t.InstanceID,
			"reward":      result.Reward,
			"multiplier":  result.Multiplier,
		})
	}

	// 週賽積分：擊殺 BOSS（DAY-066）
	g.tournamentMgr.AddPoints(p.ID, p.DisplayName, tournament.PointBoss, 0)
	// 每日賽積分：擊殺 BOSS（DAY-093）
	g.notifyDailyTournamentBoss(p)
	// 賽季積分同步（DAY-072）：BOSS 擊殺 = 50 分
	newLevels := g.addSeasonPoints(p.ID, 50)
	g.checkSeasonLevelNotify(p, newLevels)
	// 公會任務進度：擊殺 BOSS（DAY-074）
	guildID := g.Guild.GetPlayerGuildID(p.ID)
	if guildID != "" {
		completedTasks := g.Guild.UpdateTaskProgress(p.ID, guild.TaskKillBoss, 1)
		g.notifyGuildTaskComplete(guildID, completedTasks)
	}
	// 公會戰積分：BOSS 擊殺（DAY-076）
	go g.notifyGuildWarBoss(p.ID)
	// 隱藏挑戰：首次擊敗 BOSS（DAY-085）
	g.notifyChallengeBoss(p)
	// 玩家統計：記錄 BOSS 擊殺（DAY-096）
	g.notifyStatsBossKill(p)

	g.transitionState(state.StateBossResult)
	g.safeAfterFunc(3*time.Second, func() {
		g.transitionState(state.StateNormalPlay)
	})
}

// nextBossSchedule 計算下次 BOSS 觸發時間（3-5 分鐘）
func nextBossSchedule() time.Duration {
	return time.Duration(180+rand.Intn(120)) * time.Second
}
