// stats_handler.go — 玩家統計系統 handler（DAY-096）
package game

import (
	"log"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// handleGetPlayerStats 處理玩家查詢個人統計請求
func (g *Game) handleGetPlayerStats(p *player.Player) {
	if p.Stats == nil {
		return
	}
	g.sendPlayerStats(p)
}

// sendPlayerStats 發送玩家統計給指定玩家
func (g *Game) sendPlayerStats(p *player.Player) {
	if p.Stats == nil {
		return
	}
	snap := p.Stats.Snapshot()
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgPlayerStatsUpdate,
		Payload: ws.PlayerStatsPayload{
			TotalSessions:      snap.TotalSessions,
			TotalPlayTimeSec:   snap.TotalPlayTimeSec,
			TotalShots:         snap.TotalShots,
			TotalKills:         snap.TotalKills,
			TotalBet:           snap.TotalBet,
			TotalReward:        snap.TotalReward,
			TotalBonuses:       snap.TotalBonuses,
			TotalBossKills:     snap.TotalBossKills,
			BestMultiplier:     snap.BestMultiplier,
			BestStreak:         snap.BestStreak,
			BestSessionScore:   snap.BestSessionScore,
			BestBonusReward:    snap.BestBonusReward,
			MaxCoins:           snap.MaxCoins,
			JackpotWins:        snap.JackpotWins,
			JackpotMiniWins:    snap.JackpotMiniWins,
			JackpotMinorWins:   snap.JackpotMinorWins,
			JackpotMajorWins:   snap.JackpotMajorWins,
			JackpotGrandWins:   snap.JackpotGrandWins,
			TotalJackpotPayout: snap.TotalJackpotPayout,
			HitRate:            snap.HitRate,
			RTP:                snap.RTP,
			FirstPlayAtMs:      snap.FirstPlayAt,
			LastPlayAtMs:       snap.LastPlayAt,
		},
	}); err != nil {
		log.Printf("[Stats] send error: %v", err)
	}
}

// notifyStatsKill 擊破目標時更新統計
func (g *Game) notifyStatsKill(p *player.Player, multiplier float64, reward int) {
	if p.Stats == nil {
		return
	}
	p.Stats.RecordKill(multiplier, reward)
	p.Stats.UpdateMaxCoins(p.Coins)
}

// notifyStatsShot 射擊時更新統計
func (g *Game) notifyStatsShot(p *player.Player, betCost int) {
	if p.Stats == nil {
		return
	}
	p.Stats.RecordShot(betCost)
}

// notifyStatsBonus 觸發 Bonus 時更新統計
func (g *Game) notifyStatsBonus(p *player.Player, reward int) {
	if p.Stats == nil {
		return
	}
	p.Stats.RecordBonus(reward)
}

// notifyStatsBossKill 擊殺 BOSS 時更新統計
func (g *Game) notifyStatsBossKill(p *player.Player) {
	if p.Stats == nil {
		return
	}
	p.Stats.RecordBossKill()
}

// notifyStatsJackpot Jackpot 中獎時更新統計
func (g *Game) notifyStatsJackpot(p *player.Player, level string, amount int) {
	if p.Stats == nil {
		return
	}
	p.Stats.RecordJackpot(level, amount)
}

// notifyStatsStreak 連擊更新時記錄最高連擊
func (g *Game) notifyStatsStreak(p *player.Player, streak int) {
	if p.Stats == nil {
		return
	}
	p.Stats.RecordStreak(streak)
}
