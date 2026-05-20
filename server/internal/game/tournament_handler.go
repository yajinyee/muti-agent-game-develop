// tournament_handler.go — 週賽 + 每日賽 + 多格式賽 handler（DAY-093 / DAY-111）
package game

import (
	"log"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
	"digital-twin/server/internal/game/tournament"
)

// handleGetTournament 處理玩家主動查詢週賽/每日賽/多格式賽狀態（DAY-093）
func (g *Game) handleGetTournament(p *player.Player) {
	// 發送週賽狀態
	g.sendTournamentUpdate(p.ID)
	// 發送每日賽狀態
	g.sendDailyTournamentUpdate(p.ID)
	// 發送多格式賽狀態（DAY-111）
	g.sendMultiFormatUpdate(p.ID)
}

// sendTournamentUpdate 發送週賽排名給特定玩家
func (g *Game) sendTournamentUpdate(playerID string) {
	snap := g.tournamentMgr.GetSnapshot()
	rank, points := g.tournamentMgr.GetPlayerRank(playerID)

	rankings := make([]ws.TournamentRankEntry, len(snap.Rankings))
	for i, r := range snap.Rankings {
		rankings[i] = ws.TournamentRankEntry{
			Rank:        r.Rank,
			PlayerID:    r.PlayerID,
			DisplayName: r.DisplayName,
			Points:      r.Points,
			Prize:       r.Prize,
			PrizeLabel:  r.PrizeLabel,
			IsSelf:      r.PlayerID == playerID,
		}
	}

	if err := g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgTournamentUpdate,
		Payload: ws.TournamentUpdatePayload{
			WeekStart:    snap.WeekStart,
			WeekEnd:      snap.WeekEnd,
			SecondsLeft:  snap.SecondsLeft,
			Rankings:     rankings,
			TotalPlayers: snap.TotalPlayers,
			PlayerRank:   rank,
			PlayerPoints: points,
		},
	}); err != nil {
		log.Printf("[Tournament] send weekly update error: %v", err)
	}
}

// sendDailyTournamentUpdate 發送每日賽排名給特定玩家（DAY-093）
func (g *Game) sendDailyTournamentUpdate(playerID string) {
	snap := g.dailyTournamentMgr.GetDailySnapshot()
	rank, points := g.dailyTournamentMgr.GetPlayerRank(playerID)

	rankings := make([]ws.TournamentRankEntry, len(snap.Rankings))
	for i, r := range snap.Rankings {
		rankings[i] = ws.TournamentRankEntry{
			Rank:        r.Rank,
			PlayerID:    r.PlayerID,
			DisplayName: r.DisplayName,
			Points:      r.Points,
			Prize:       r.Prize,
			PrizeLabel:  r.PrizeLabel,
			IsSelf:      r.PlayerID == playerID,
		}
	}

	if err := g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgDailyTournamentUpdate,
		Payload: ws.DailyTournamentUpdatePayload{
			DayStart:     snap.DayStart,
			DayEnd:       snap.DayEnd,
			SecondsLeft:  snap.SecondsLeft,
			Rankings:     rankings,
			TotalPlayers: snap.TotalPlayers,
			PlayerRank:   rank,
			PlayerPoints: points,
		},
	}); err != nil {
		log.Printf("[DailyTournament] send daily update error: %v", err)
	}
}

// broadcastDailyTournament 廣播每日賽排名給所有玩家（每 30 秒，DAY-093）
func (g *Game) broadcastDailyTournament() {
	snap := g.dailyTournamentMgr.GetDailySnapshot()

	g.mu.RLock()
	playerIDs := make([]string, 0, len(g.Players))
	for id := range g.Players {
		playerIDs = append(playerIDs, id)
	}
	g.mu.RUnlock()

	rankings := make([]ws.TournamentRankEntry, len(snap.Rankings))
	for i, r := range snap.Rankings {
		rankings[i] = ws.TournamentRankEntry{
			Rank:        r.Rank,
			PlayerID:    r.PlayerID,
			DisplayName: r.DisplayName,
			Points:      r.Points,
			Prize:       r.Prize,
			PrizeLabel:  r.PrizeLabel,
		}
	}

	for _, pid := range playerIDs {
		rank, points := g.dailyTournamentMgr.GetPlayerRank(pid)

		personalRankings := make([]ws.TournamentRankEntry, len(rankings))
		copy(personalRankings, rankings)
		for i := range personalRankings {
			personalRankings[i].IsSelf = (personalRankings[i].PlayerID == pid)
		}

		g.Hub.Send(pid, &ws.Message{
			Type: ws.MsgDailyTournamentUpdate,
			Payload: ws.DailyTournamentUpdatePayload{
				DayStart:     snap.DayStart,
				DayEnd:       snap.DayEnd,
				SecondsLeft:  snap.SecondsLeft,
				Rankings:     personalRankings,
				TotalPlayers: snap.TotalPlayers,
				PlayerRank:   rank,
				PlayerPoints: points,
			},
		})
	}
}

// notifyDailyTournamentKill 每次擊破目標後更新每日賽積分（DAY-093）
func (g *Game) notifyDailyTournamentKill(p *player.Player, multiplier float64) {
	g.dailyTournamentMgr.AddPoints(p.ID, p.DisplayName, tournament.PointKill, multiplier)
}

// notifyDailyTournamentBoss 擊殺 BOSS 後更新每日賽積分（DAY-093）
func (g *Game) notifyDailyTournamentBoss(p *player.Player) {
	g.dailyTournamentMgr.AddPoints(p.ID, p.DisplayName, tournament.PointBoss, 0)
}

// notifyDailyTournamentBonus 完成 Bonus 後更新每日賽積分（DAY-093）
func (g *Game) notifyDailyTournamentBonus(p *player.Player) {
	g.dailyTournamentMgr.AddPoints(p.ID, p.DisplayName, tournament.PointBonus, 0)
}

// GetDailyTournamentSnapshot 取得每日賽快照（供 HTTP 端點使用，DAY-093）
func (g *Game) GetDailyTournamentSnapshot() tournament.DailySnapshot {
	return g.dailyTournamentMgr.GetDailySnapshot()
}

// ---- 多格式每日賽（DAY-111）----

// sendMultiFormatUpdate 發送多格式賽排名給特定玩家（DAY-111）
func (g *Game) sendMultiFormatUpdate(playerID string) {
	snap := g.multiFormatMgr.GetSnapshot()

	rankings := make([]ws.MultiFormatRankEntry, len(snap.Rankings))
	for i, r := range snap.Rankings {
		rankings[i] = ws.MultiFormatRankEntry{
			Rank:        r.Rank,
			PlayerID:    r.PlayerID,
			DisplayName: r.DisplayName,
			Score:       r.Score,
			ScoreLabel:  r.ScoreLabel,
			Prize:       r.Prize,
			PrizeLabel:  r.PrizeLabel,
			IsSelf:      r.PlayerID == playerID,
		}
	}

	// 取得玩家自己的排名和分數
	playerRank, playerScore := g.multiFormatMgr.GetPlayerRank(playerID)

	if err := g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgMultiFormatUpdate,
		Payload: ws.MultiFormatUpdatePayload{
			DayStart:       snap.DayStart,
			DayEnd:         snap.DayEnd,
			SecondsLeft:    snap.SecondsLeft,
			TodayFormat:    string(snap.TodayFormat),
			FormatName:     snap.FormatDef.Name,
			FormatIcon:     snap.FormatDef.Icon,
			FormatUnit:     snap.FormatDef.Unit,
			FormatDesc:     snap.FormatDef.Description,
			Rankings:       rankings,
			TotalPlayers:   snap.TotalPlayers,
			PlayerRank:     playerRank,
			PlayerScore:    playerScore,
			NextFormat:     string(snap.NextFormat),
			NextFormatName: snap.NextFormatDef.Name,
			NextFormatIcon: snap.NextFormatDef.Icon,
		},
	}); err != nil {
		log.Printf("[MultiFormat] send update error: %v", err)
	}
}

// broadcastMultiFormat 廣播多格式賽排名給所有玩家（每 30 秒，DAY-111）
func (g *Game) broadcastMultiFormat() {
	snap := g.multiFormatMgr.GetSnapshot()

	g.mu.RLock()
	playerIDs := make([]string, 0, len(g.Players))
	for id := range g.Players {
		playerIDs = append(playerIDs, id)
	}
	g.mu.RUnlock()

	baseRankings := make([]ws.MultiFormatRankEntry, len(snap.Rankings))
	for i, r := range snap.Rankings {
		baseRankings[i] = ws.MultiFormatRankEntry{
			Rank:        r.Rank,
			PlayerID:    r.PlayerID,
			DisplayName: r.DisplayName,
			Score:       r.Score,
			ScoreLabel:  r.ScoreLabel,
			Prize:       r.Prize,
			PrizeLabel:  r.PrizeLabel,
		}
	}

	for _, pid := range playerIDs {
		playerRank, playerScore := g.multiFormatMgr.GetPlayerRank(pid)

		personalRankings := make([]ws.MultiFormatRankEntry, len(baseRankings))
		copy(personalRankings, baseRankings)
		for i := range personalRankings {
			personalRankings[i].IsSelf = (personalRankings[i].PlayerID == pid)
		}

		g.Hub.Send(pid, &ws.Message{
			Type: ws.MsgMultiFormatUpdate,
			Payload: ws.MultiFormatUpdatePayload{
				DayStart:       snap.DayStart,
				DayEnd:         snap.DayEnd,
				SecondsLeft:    snap.SecondsLeft,
				TodayFormat:    string(snap.TodayFormat),
				FormatName:     snap.FormatDef.Name,
				FormatIcon:     snap.FormatDef.Icon,
				FormatUnit:     snap.FormatDef.Unit,
				FormatDesc:     snap.FormatDef.Description,
				Rankings:       personalRankings,
				TotalPlayers:   snap.TotalPlayers,
				PlayerRank:     playerRank,
				PlayerScore:    playerScore,
				NextFormat:     string(snap.NextFormat),
				NextFormatName: snap.NextFormatDef.Name,
				NextFormatIcon: snap.NextFormatDef.Icon,
			},
		})
	}
}

// GetMultiFormatSnapshot 取得多格式賽快照（供 HTTP 端點使用，DAY-111）
func (g *Game) GetMultiFormatSnapshot() tournament.MultiFormatSnapshot {
	return g.multiFormatMgr.GetSnapshot()
}
