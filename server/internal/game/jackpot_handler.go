// Package game — Jackpot 相關 handler（DAY-057 拆分，DAY-095 升級四層）
package game

import (
	"log"
	"time"

	"digital-twin/server/internal/game/jackpot"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// GetJackpotSnapshot 取得 Jackpot 池當前金額快照（thread-safe）
func (g *Game) GetJackpotSnapshot() map[string]int {
	snap := g.jackpotMgr.GetSnapshot()
	return map[string]int{
		"mini":  snap[jackpot.LevelMini],
		"minor": snap[jackpot.LevelMinor],
		"major": snap[jackpot.LevelMajor],
		"grand": snap[jackpot.LevelGrand],
	}
}

// GetJackpotHistory 取得最近 Jackpot 中獎記錄（DAY-048e）
func (g *Game) GetJackpotHistory(n int) []jackpot.WinRecord {
	return g.jackpotHist.GetRecent(n)
}

// GetJackpotDailyStats 取得今日 Jackpot 統計（DAY-049）
func (g *Game) GetJackpotDailyStats() jackpot.DailyStats {
	return g.jackpotHist.GetDailyStats()
}

// handleJackpotWin 處理 Jackpot 中獎（DAY-048，DAY-095 升級）
// 發放獎勵、記錄歷史、廣播中獎通知 + 動畫通知
func (g *Game) handleJackpotWin(p *player.Player, win *jackpot.JackpotWin) {
	// 發放獎勵給中獎玩家
	p.AddReward(win.Amount)

	// 取得顯示名稱
	displayName := p.DisplayName
	if displayName == "" {
		displayName = p.ID[:8]
	}

	// 取得等級顯示資訊
	levelName, levelColor, levelIcon := jackpot.GetLevelInfo(win.Level)
	isGrand := win.Level == jackpot.LevelGrand
	isMajor := win.Level == jackpot.LevelMajor

	// 記錄到歷史（DAY-048e）
	g.jackpotHist.Add(win, displayName)

	// 廣播動畫通知給所有玩家（DAY-095）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgJackpotAnimation,
		Payload: ws.JackpotAnimationPayload{
			Level:      string(win.Level),
			LevelName:  levelName,
			LevelColor: levelColor,
			LevelIcon:  levelIcon,
			Amount:     win.Amount,
			WinnerName: displayName,
			IsGrand:    isGrand,
			IsMajor:    isMajor,
		},
	})

	// 廣播中獎通知給所有玩家
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgJackpotWin,
		Payload: ws.JackpotWinPayload{
			Level:      string(win.Level),
			LevelName:  levelName,
			LevelColor: levelColor,
			LevelIcon:  levelIcon,
			Amount:     win.Amount,
			WinnerID:   p.ID,
			WinnerName: displayName,
			NewBalance: p.Coins,
			IsGrand:    isGrand,
		},
	})

	// 更新中獎玩家的狀態
	g.sendPlayerUpdate(p)

	// 隱藏挑戰：Jackpot 中獎（DAY-085）
	g.notifyChallengeJackpot(p.ID)
	// 玩家統計：記錄 Jackpot 中獎（DAY-096）
	g.notifyStatsJackpot(p, string(win.Level), win.Amount)
	// 全服公告：Jackpot 中獎（DAY-097）
	g.announceJackpotWin(displayName, string(win.Level), levelName, win.Amount)

	log.Printf("[Jackpot] %s won %s jackpot: %d coins (player: %s, grand=%v)",
		p.ID, win.Level, win.Amount, displayName, isGrand)
}

// broadcastJackpot 廣播 Jackpot 池當前金額（每 5 秒，DAY-048，DAY-095 升級四層）
func (g *Game) broadcastJackpot() {
	snap := g.jackpotMgr.GetSnapshot()
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgJackpotUpdate,
		Payload: ws.JackpotUpdatePayload{
			Mini:  snap[jackpot.LevelMini],
			Minor: snap[jackpot.LevelMinor],
			Major: snap[jackpot.LevelMajor],
			Grand: snap[jackpot.LevelGrand],
		},
	})
}

// saveJackpotState 儲存 Jackpot 池狀態到 Store（DAY-049d）
// 每 30 秒自動呼叫，確保 Server 重啟後能恢復 Jackpot 池
func (g *Game) saveJackpotState() {
	if g.store == nil {
		return
	}
	state := g.jackpotMgr.SaveState()
	key := "jackpot_state:" + g.ID
	if err := g.store.SetJSON(key, state, 7*24*time.Hour); err != nil {
		log.Printf("[Jackpot] Failed to save state: %v", err)
	}
}

// loadJackpotState 從 Store 恢復 Jackpot 池狀態（DAY-049d）
// 在 Game 啟動時呼叫
func (g *Game) loadJackpotState() {
	if g.store == nil {
		return
	}
	key := "jackpot_state:" + g.ID
	var state jackpot.PoolState
	if err := g.store.GetJSON(key, &state); err != nil {
		// 找不到或解析失敗，使用預設值（正常情況）
		return
	}
	g.jackpotMgr.LoadState(state)
	log.Printf("[Jackpot] Restored state: mini=%d minor=%d major=%d grand=%d",
		state.Mini, state.Minor, state.Major, state.Grand)
}
