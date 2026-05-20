// player_card_handler.go — 玩家名片系統（DAY-106）
// 玩家可以查看其他在線玩家的名片（稱號/VIP/公會/統計亮點）
package game

import (
	"log"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// handleGetPlayerCard 處理查詢玩家名片請求（DAY-106）
func (g *Game) handleGetPlayerCard(p *player.Player, msg *ws.Message) {
	var payload ws.GetPlayerCardPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	targetID := payload.TargetPlayerID
	if targetID == "" {
		targetID = p.ID // 查詢自己
	}

	card := g.buildPlayerCard(targetID)
	if card == nil {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgError,
			Payload: ws.ErrorPayload{
				Code:    "player_not_found",
				Message: "玩家不存在或已離線",
			},
		})
		return
	}

	g.Hub.Send(p.ID, &ws.Message{
		Type:    ws.MsgPlayerCard,
		Payload: card,
	})

	log.Printf("[PlayerCard] %s viewed card of %s", p.ID, targetID)
}

// buildPlayerCard 建立玩家名片（DAY-106）
func (g *Game) buildPlayerCard(playerID string) *ws.PlayerCardPayload {
	g.mu.RLock()
	target, online := g.Players[playerID]
	g.mu.RUnlock()

	if !online || target == nil {
		return nil
	}

	snap := target.Snapshot()
	leaderSnap := target.LeaderboardSnapshot()
	loginStreak, _ := target.GetLoginInfo()

	// VIP 資訊
	vipSnap := g.VIP.GetSnapshot(playerID)

	// 公會資訊
	guildName := ""
	guildRole := ""
	if guildID := g.Guild.GetPlayerGuildID(playerID); guildID != "" {
		if gd := g.Guild.GetGuild(guildID); gd != nil {
			guildName = gd.Name
			if member, ok := gd.Members[playerID]; ok {
				guildRole = string(member.Role)
			}
		}
	}

	// 統計亮點
	bestStreak := 0
	bestMult := 0.0
	jackpotWins := 0
	rtp := 0.0
	if target.Stats != nil {
		statsSnap := target.Stats.Snapshot()
		bestStreak = statsSnap.BestStreak
		bestMult = statsSnap.BestMultiplier
		jackpotWins = statsSnap.JackpotWins
		rtp = statsSnap.RTP
	}

	// 成就數量
	achUnlocked := target.GetAchievements()

	return &ws.PlayerCardPayload{
		PlayerID:         playerID,
		DisplayName:      snap.DisplayName,
		TitleName:        snap.TitleName,
		TitleIcon:        snap.TitleIcon,
		TitleColor:       snap.TitleColor,
		VIPLevel:         vipSnap.VIPLevel,
		VIPName:          vipSnap.TierName,
		GuildName:        guildName,
		GuildRole:        guildRole,
		KillCount:        snap.KillCount,
		MaxCoins:         leaderSnap.MaxCoins,
		BestStreak:       bestStreak,
		BestMult:         bestMult,
		JackpotWins:      jackpotWins,
		AchievementCount: len(achUnlocked),
		LoginStreak:      loginStreak,
		RTP:              rtp,
		IsOnline:         true,
	}
}
