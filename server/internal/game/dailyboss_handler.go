// dailyboss_handler.go — 每日特殊 BOSS 挑戰 handler（DAY-077）
package game

import (
	"log"

	"digital-twin/server/internal/game/dailyboss"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// broadcastDailyBoss 廣播每日 BOSS 狀態給所有玩家（每 30 秒，DAY-077）
func (g *Game) broadcastDailyBoss() {
	// 檢查是否需要重置
	if reset, old := g.DailyBoss.CheckAndReset(); reset && old != nil {
		log.Printf("[DailyBoss] Day %s expired (status: %s), new boss spawned", old.DateID, old.Status)
	}

	g.mu.RLock()
	players := make([]*player.Player, 0, len(g.Players))
	for _, p := range g.Players {
		players = append(players, p)
	}
	g.mu.RUnlock()

	for _, p := range players {
		g.sendDailyBossUpdate(p)
	}
}

// sendDailyBossUpdate 發送每日 BOSS 狀態給指定玩家
func (g *Game) sendDailyBossUpdate(p *player.Player) {
	snap := g.DailyBoss.GetSnapshot()
	if snap == nil {
		return
	}

	// 取得前 5 名貢獻者
	topContribs := g.DailyBoss.GetTopContributors(5)
	contribEntries := make([]ws.DailyBossContributorEntry, 0, len(topContribs))
	for i, c := range topContribs {
		contribEntries = append(contribEntries, ws.DailyBossContributorEntry{
			Rank:        i + 1,
			PlayerID:    c.PlayerID,
			DisplayName: c.DisplayName,
			Damage:      c.Damage,
			Reward:      c.Reward,
			IsMe:        c.PlayerID == p.ID,
		})
	}

	// 取得玩家自己的貢獻
	myDamage := 0
	myReward := 0
	if myContrib := g.DailyBoss.GetPlayerContribution(p.ID); myContrib != nil {
		myDamage = myContrib.Damage
		myReward = myContrib.Reward
	}

	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgDailyBossUpdate,
		Payload: ws.DailyBossUpdatePayload{
			DateID:        snap.DateID,
			BossID:        snap.BossType.ID,
			BossName:      snap.BossType.Name,
			BossIcon:      snap.BossType.Icon,
			BossColor:     snap.BossType.Color,
			Description:   snap.BossType.Description,
			MaxHP:         snap.MaxHP,
			CurrentHP:     snap.CurrentHP,
			HPPercent:     g.DailyBoss.GetHPPercent(),
			Status:        string(snap.Status),
			EndAt:         snap.EndAt.UnixMilli(),
			RewardPool:    snap.RewardPool,
			TopContribs:   contribEntries,
			MyDamage:      myDamage,
			MyReward:      myReward,
			DifficultyMod: snap.DifficultyMod,
		},
	})
}

// handleGetDailyBoss 處理查詢每日 BOSS 狀態請求
func (g *Game) handleGetDailyBoss(p *player.Player) {
	g.sendDailyBossUpdate(p)
}

// handleDailyBossAttack 處理對每日 BOSS 的攻擊
// 每次普通攻擊命中時，自動貢獻 1 點傷害（由 handleKill 呼叫）
// 也可以由 Client 主動發送 daily_boss_attack 訊息
func (g *Game) handleDailyBossAttack(p *player.Player, msg *ws.Message) {
	var payload ws.DailyBossAttackPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	// 驗證傷害值（防止作弊，最大 100 點/次）
	damage := payload.Damage
	if damage <= 0 {
		damage = 1
	}
	if damage > 100 {
		damage = 100
	}

	g.applyDailyBossDamage(p, damage)
}

// applyDailyBossDamage 對每日 BOSS 造成傷害（由 handleKill 呼叫）
func (g *Game) applyDailyBossDamage(p *player.Player, damage int) {
	defeated, reward := g.DailyBoss.AddDamage(p.ID, p.DisplayName, damage)

	if defeated {
		log.Printf("[DailyBoss] Defeated by %s! Reward: %d", p.ID, reward)

		// 發放獎勵給擊殺者
		if reward > 0 {
			g.mu.Lock()
			if player, ok := g.Players[p.ID]; ok {
				player.AddCoins(reward)
			}
			g.mu.Unlock()
		}

		// 廣播擊殺通知給所有玩家
		g.broadcastDailyBossDefeated(p.ID, p.DisplayName)
	}
}

// broadcastDailyBossDefeated 廣播每日 BOSS 擊殺通知
func (g *Game) broadcastDailyBossDefeated(killerID, killerName string) {
	snap := g.DailyBoss.GetSnapshot()
	if snap == nil {
		return
	}

	// 取得所有貢獻者排名
	allContribs := g.DailyBoss.GetTopContributors(20)
	rankings := make([]ws.DailyBossContributorEntry, 0, len(allContribs))
	for i, c := range allContribs {
		rankings = append(rankings, ws.DailyBossContributorEntry{
			Rank:        i + 1,
			PlayerID:    c.PlayerID,
			DisplayName: c.DisplayName,
			Damage:      c.Damage,
			Reward:      c.Reward,
		})
	}

	// 發放所有貢獻者的獎勵
	g.mu.Lock()
	for _, c := range allContribs {
		if c.Reward > 0 {
			if player, ok := g.Players[c.PlayerID]; ok {
				player.AddCoins(c.Reward)
				log.Printf("[DailyBoss] Player %s received %d coins (damage: %d)", c.PlayerID, c.Reward, c.Damage)
			}
		}
	}
	players := make([]*player.Player, 0, len(g.Players))
	for _, p := range g.Players {
		players = append(players, p)
	}
	g.mu.Unlock()

	// 廣播給所有玩家（每人看到自己的獎勵）
	for _, p := range players {
		myReward := 0
		for _, c := range allContribs {
			if c.PlayerID == p.ID {
				myReward = c.Reward
				break
			}
		}

		// 標記自己的排名
		personalRankings := make([]ws.DailyBossContributorEntry, len(rankings))
		copy(personalRankings, rankings)
		for i := range personalRankings {
			personalRankings[i].IsMe = personalRankings[i].PlayerID == p.ID
		}

		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgDailyBossDefeated,
			Payload: ws.DailyBossDefeatedPayload{
				DateID:      snap.DateID,
				BossName:    snap.BossType.Name,
				BossIcon:    snap.BossType.Icon,
				KillerID:    killerID,
				KillerName:  killerName,
				Rankings:    personalRankings,
				MyReward:    myReward,
				TotalDamage: snap.TotalDamage,
			},
		})
	}
}

// notifyDailyBossKill 每次擊破目標時，自動貢獻每日 BOSS 傷害（由 handleKill 呼叫）
// 傷害 = floor(multiplier)，最少 1 點
func (g *Game) notifyDailyBossKill(p *player.Player, multiplier float64) {
	if g.DailyBoss.GetStatus() != dailyboss.BossStatusActive {
		return
	}
	damage := int(multiplier)
	if damage < 1 {
		damage = 1
	}
	if damage > 50 {
		damage = 50 // 單次最多 50 點
	}
	g.applyDailyBossDamage(p, damage)
}
