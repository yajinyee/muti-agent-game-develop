// raid_handler.go — Co-op Boss Raid 系統 handler（DAY-115）
// 全服玩家合作討伐超強 BOSS，依貢獻度分配獎勵池
package game

import (
	"log"
	"time"

	raidboss "digital-twin/server/internal/game/raidBoss"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// Raid BOSS 設定
const (
	RaidBossName   = "吉伊卡哇大魔王"
	RaidBossHP     = 50000 // 超強 BOSS，需要全服合作
	RaidRewardPool = 200000 // 獎勵池總金幣
	RaidWarningDur = 30 * time.Second
)

// triggerRaid 觸發 Co-op Boss Raid（每日一次，或手動觸發）
func (g *Game) triggerRaid() {
	todayDate := time.Now().Format("2006-01-02")
	if !g.RaidBoss.CanTrigger(todayDate) {
		log.Printf("[Raid] already triggered today or not idle")
		return
	}

	raidID := g.RaidBoss.StartWarning()
	log.Printf("[Raid] warning started: raidID=%s", raidID)

	// 廣播警告
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRaidWarning,
		Payload: ws.RaidWarningPayload{
			RaidID:     raidID,
			BossName:   RaidBossName,
			MaxHP:      RaidBossHP,
			RewardPool: RaidRewardPool,
			StartsIn:   int(RaidWarningDur.Seconds()),
		},
	})

	// 全服公告
	g.announceEventStart("⚔️ Co-op Boss Raid 即將開始！全服合作討伐「" + RaidBossName + "」！獎勵池 20萬 金幣！")

	// 30 秒後開始討伐
	g.safeAfterFunc(RaidWarningDur, func() {
		g.startRaid()
	})
}

// startRaid 開始討伐
func (g *Game) startRaid() {
	todayDate := time.Now().Format("2006-01-02")
	g.RaidBoss.StartRaid(RaidBossName, RaidBossHP, RaidRewardPool, todayDate)
	log.Printf("[Raid] started: boss=%s HP=%d pool=%d", RaidBossName, RaidBossHP, RaidRewardPool)

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRaidStart,
		Payload: ws.RaidStartPayload{
			RaidID:     g.RaidBoss.GetSnapshot().RaidID,
			BossName:   RaidBossName,
			HP:         RaidBossHP,
			MaxHP:      RaidBossHP,
			RewardPool: RaidRewardPool,
			Duration:   int(raidboss.RaidDuration.Seconds()),
		},
	})
}

// notifyRaidKill 玩家擊破目標時，對 Raid BOSS 造成傷害（由 handleKill 呼叫）
// damage = result.Reward（用獎勵值作為傷害，讓高投注玩家貢獻更多）
func (g *Game) notifyRaidKill(p *player.Player, damage int) {
	if !g.RaidBoss.IsActive() {
		return
	}

	newHP, killed := g.RaidBoss.RecordDamage(p.ID, p.DisplayName, damage)
	if killed {
		g.handleRaidKill()
	} else {
		// 每次傷害後廣播更新（節流：由 tickRaidUpdate 每 3 秒廣播）
		_ = newHP
	}
}

// handleRaidKill 討伐成功（BOSS 被擊殺）
func (g *Game) handleRaidKill() {
	snap := g.RaidBoss.GetSnapshot()
	log.Printf("[Raid] BOSS defeated! distributing rewards to %d contributors", len(snap.Contributors))

	// 發放獎勵給所有貢獻者
	g.mu.RLock()
	players := make(map[string]*player.Player, len(g.Players))
	for id, p := range g.Players {
		players[id] = p
	}
	g.mu.RUnlock()

	contributors := buildRaidContributorPayloads(snap.Contributors)

	// 廣播結算結果（含每個玩家的獎勵）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRaidResult,
		Payload: ws.RaidResultPayload{
			RaidID:       snap.RaidID,
			BossName:     snap.BossName,
			Defeated:     true,
			RewardPool:   snap.RewardPool,
			Contributors: contributors,
		},
	})

	// 發放金幣給在線玩家
	for _, entry := range snap.Contributors {
		if p, ok := players[entry.PlayerID]; ok && entry.Reward > 0 {
			p.AddCoins(entry.Reward)
			g.Hub.Send(p.ID, &ws.Message{
				Type: ws.MsgReward,
				Payload: ws.RewardPayload{
					Source:     "raid_boss",
					Amount:     entry.Reward,
					Multiplier: 1.0,
					NewBalance: p.GetCoins(),
				},
			})
			log.Printf("[Raid] reward player=%s rank=%d reward=%d", p.ID, entry.Rank, entry.Reward)
		}
	}

	// 動態牆：Raid 勝利
	if len(snap.Contributors) > 0 {
		topPlayer := snap.Contributors[0]
		go g.notifyFeedBossKill(&player.Player{
			ID:          topPlayer.PlayerID,
			DisplayName: topPlayer.DisplayName,
		}, snap.RewardPool)
	}

	// 全服公告
	g.announceEventStart("🏆 Raid 勝利！全服合作擊敗「" + snap.BossName + "」！獎勵池已分配！")

	// 5 秒後重置
	g.safeAfterFunc(5*time.Second, func() {
		g.RaidBoss.Reset()
	})
}

// handleRaidTimeout 討伐超時（BOSS 未被擊殺）
func (g *Game) handleRaidTimeout() {
	snap := g.RaidBoss.GetSnapshot()
	log.Printf("[Raid] timeout: boss=%s HP=%d/%d", snap.BossName, snap.HP, snap.MaxHP)

	contributors := buildRaidContributorPayloads(snap.Contributors)

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRaidResult,
		Payload: ws.RaidResultPayload{
			RaidID:       snap.RaidID,
			BossName:     snap.BossName,
			Defeated:     false,
			RewardPool:   snap.RewardPool,
			Contributors: contributors,
		},
	})

	// 全服公告
	g.announceEventStart("💀 Raid 失敗！「" + snap.BossName + "」逃脫了！明天再來！")

	g.safeAfterFunc(5*time.Second, func() {
		g.RaidBoss.Reset()
	})
}

// tickRaidUpdate 每 3 秒廣播討伐狀態（由 gameLoop 呼叫）
func (g *Game) tickRaidUpdate() {
	if !g.RaidBoss.IsActive() {
		// 檢查超時
		if g.RaidBoss.CheckTimeout() {
			g.handleRaidTimeout()
		}
		return
	}

	// 檢查超時
	if g.RaidBoss.CheckTimeout() {
		g.handleRaidTimeout()
		return
	}

	snap := g.RaidBoss.GetSnapshot()
	contributors := buildRaidContributorPayloads(snap.Contributors)

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRaidUpdate,
		Payload: ws.RaidUpdatePayload{
			RaidID:       snap.RaidID,
			HP:           snap.HP,
			MaxHP:        snap.MaxHP,
			TimeLeft:     snap.TimeLeft,
			Contributors: contributors,
		},
	})
}

// handleGetRaidStatus 處理玩家查詢討伐狀態
func (g *Game) handleGetRaidStatus(p *player.Player) {
	snap := g.RaidBoss.GetSnapshot()
	todayDate := time.Now().Format("2006-01-02")
	contributors := buildRaidContributorPayloads(snap.Contributors)

	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgRaidStatus,
		Payload: ws.RaidStatusPayload{
			State:        string(snap.State),
			RaidID:       snap.RaidID,
			BossName:     snap.BossName,
			HP:           snap.HP,
			MaxHP:        snap.MaxHP,
			RewardPool:   snap.RewardPool,
			TimeLeft:     snap.TimeLeft,
			Contributors: contributors,
			CanTrigger:   g.RaidBoss.CanTrigger(todayDate),
		},
	})
}

// buildRaidContributorPayloads 將 ContributorEntry 轉換為 Payload
func buildRaidContributorPayloads(entries []*raidboss.ContributorEntry) []*ws.RaidContributorPayload {
	result := make([]*ws.RaidContributorPayload, 0, len(entries))
	for _, e := range entries {
		result = append(result, &ws.RaidContributorPayload{
			PlayerID:    e.PlayerID,
			DisplayName: e.DisplayName,
			Damage:      e.Damage,
			Reward:      e.Reward,
			Rank:        e.Rank,
		})
	}
	return result
}