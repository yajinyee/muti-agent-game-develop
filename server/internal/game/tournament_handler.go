// tournament_handler.go — 週賽 + 每日賽 handler（DAY-093）
package game

import (
	"log"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
	"digital-twin/server/internal/game/tournament"
)

// handleGetTournament 處理玩家主動查詢週賽/每日賽狀態（DAY-093）
func (g *Game) handleGetTournament(p *player.Player) {
	// 發送週賽狀態
	g.sendTournamentUpdate(p.ID)
	// 發送每日賽狀態
	g.sendDailyTournamentUpdate(p.ID)
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
