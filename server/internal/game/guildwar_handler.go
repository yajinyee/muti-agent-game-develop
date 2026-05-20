// guildwar_handler.go — 公會戰系統 handler（DAY-076）
package game

import (
	"log"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// broadcastGuildWar 廣播公會戰排名給所有玩家（每 60 秒，DAY-076）
func (g *Game) broadcastGuildWar() {
	g.mu.RLock()
	players := make([]*player.Player, 0, len(g.Players))
	for _, p := range g.Players {
		players = append(players, p)
	}
	g.mu.RUnlock()

	// 檢查是否需要結算
	if result := g.GuildWar.CheckAndSettle(); result != nil {
		log.Printf("[GuildWar] Week %s settled, %d guilds participated", result.WeekID, len(result.Rankings))
		// 廣播結算結果給所有玩家
		for _, p := range players {
			g.sendGuildWarResult(p, result)
		}
		return
	}

	// 廣播當前排名
	for _, p := range players {
		g.sendGuildWarUpdate(p)
	}
}

// sendGuildWarUpdate 發送公會戰排名更新給指定玩家
func (g *Game) sendGuildWarUpdate(p *player.Player) {
	status, weekID, _, endAt := g.GuildWar.GetStatus()
	rankings := g.GuildWar.GetRankings()

	myGuildID := g.Guild.GetPlayerGuildID(p.ID)
	myRank := 0
	myScore := 0

	entries := make([]ws.GuildWarScoreEntry, 0, len(rankings))
	for i, r := range rankings {
		isMyGuild := r.GuildID == myGuildID
		if isMyGuild {
			myRank = i + 1
			myScore = r.Score
		}
		entries = append(entries, ws.GuildWarScoreEntry{
			Rank:        i + 1,
			GuildID:     r.GuildID,
			GuildName:   r.GuildName,
			GuildIcon:   r.GuildIcon,
			MemberCount: r.MemberCount,
			Score:       r.Score,
			KillScore:   r.KillScore,
			BossScore:   r.BossScore,
			BonusScore:  r.BonusScore,
			IsMyGuild:   isMyGuild,
		})
	}

	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgGuildWarUpdate,
		Payload: ws.GuildWarUpdatePayload{
			WeekID:       weekID,
			Status:       string(status),
			EndAt:        endAt.UnixMilli(),
			Rankings:     entries,
			MyGuildRank:  myRank,
			MyGuildScore: myScore,
			TotalGuilds:  g.GuildWar.GetParticipatingGuildCount(),
		},
	})
}

// sendGuildWarResult 發送公會戰結算結果給指定玩家
func (g *Game) sendGuildWarResult(p *player.Player, result interface{}) {
	// 使用 type assertion 取得 WarResult
	type warResult interface {
		GetWeekID() string
	}

	// 直接使用 lastResult
	lastResult := g.GuildWar.GetLastResult()
	if lastResult == nil {
		return
	}

	myGuildID := g.Guild.GetPlayerGuildID(p.ID)
	myRank := 0
	myReward := 0

	resultEntries := make([]ws.GuildWarResultEntry, 0, len(lastResult.Rankings))
	for _, r := range lastResult.Rankings {
		if r.GuildID == myGuildID {
			myRank = r.Rank
			myReward = r.Reward
		}
		resultEntries = append(resultEntries, ws.GuildWarResultEntry{
			Rank:      r.Rank,
			GuildID:   r.GuildID,
			GuildName: r.GuildName,
			GuildIcon: r.GuildIcon,
			Score:     r.Score,
			Reward:    r.Reward,
		})
	}

	// 發放獎勵
	if myReward > 0 {
		g.mu.Lock()
		if player, ok := g.Players[p.ID]; ok {
			player.AddCoins(myReward)
			log.Printf("[GuildWar] Player %s received %d coins (rank %d)", p.ID, myReward, myRank)
		}
		g.mu.Unlock()
	}

	// 計算下週開始時間
	_, _, nextStart, _ := g.GuildWar.GetStatus()

	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgGuildWarResult,
		Payload: ws.GuildWarResultPayload{
			WeekID:    lastResult.WeekID,
			Rankings:  resultEntries,
			MyRank:    myRank,
			MyReward:  myReward,
			NextWarAt: nextStart.UnixMilli(),
		},
	})
}

// handleGetGuildWarStatus 處理查詢公會戰狀態請求
func (g *Game) handleGetGuildWarStatus(p *player.Player) {
	// 確保玩家所在公會已登記
	guildID := g.Guild.GetPlayerGuildID(p.ID)
	if guildID != "" {
		guildData := g.Guild.GetGuild(guildID)
		if guildData != nil {
			g.GuildWar.EnsureGuildRegistered(
				guildID,
				guildData.Name,
				guildData.Icon,
				len(guildData.Members),
			)
		}
	}
	g.sendGuildWarUpdate(p)
}

// notifyGuildWarKill 通知公會戰擊殺積分（由 handleKill 呼叫）
func (g *Game) notifyGuildWarKill(playerID string, multiplier int) {
	guildID := g.Guild.GetPlayerGuildID(playerID)
	if guildID == "" {
		return
	}
	// 確保公會已登記
	guildData := g.Guild.GetGuild(guildID)
	if guildData != nil {
		g.GuildWar.EnsureGuildRegistered(guildID, guildData.Name, guildData.Icon, len(guildData.Members))
	}
	g.GuildWar.AddKillScore(guildID, multiplier)
}

// notifyGuildWarBoss 通知公會戰 BOSS 積分（由 boss_handler 呼叫）
func (g *Game) notifyGuildWarBoss(playerID string) {
	guildID := g.Guild.GetPlayerGuildID(playerID)
	if guildID == "" {
		return
	}
	guildData := g.Guild.GetGuild(guildID)
	if guildData != nil {
		g.GuildWar.EnsureGuildRegistered(guildID, guildData.Name, guildData.Icon, len(guildData.Members))
	}
	g.GuildWar.AddBossScore(guildID)
}

// notifyGuildWarBonus 通知公會戰 Bonus 積分（由 bonus_handler 呼叫）
func (g *Game) notifyGuildWarBonus(playerID string) {
	guildID := g.Guild.GetPlayerGuildID(playerID)
	if guildID == "" {
		return
	}
	guildData := g.Guild.GetGuild(guildID)
	if guildData != nil {
		g.GuildWar.EnsureGuildRegistered(guildID, guildData.Name, guildData.Icon, len(guildData.Members))
	}
	g.GuildWar.AddBonusScore(guildID)
}
