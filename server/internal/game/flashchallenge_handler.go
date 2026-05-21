// flashchallenge_handler.go — 閃電挑戰系統 handler（DAY-123）
// 業界依據：Infingame（2026-05-19）確認 Challenges 工具是 2026 年最熱門留存機制
// 限時 90 秒的特殊目標挑戰，完成獎勵豐厚，全服可見增加社交競爭感
package game

import (
	"log"

	"digital-twin/server/internal/game/activityfeed"
	"digital-twin/server/internal/game/flashchallenge"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// tryStartFlashChallenge 嘗試觸發閃電挑戰（由 game loop 或 BOSS 擊殺後呼叫）
func (g *Game) tryStartFlashChallenge(triggerType string) {
	if !g.FlashChallenge.ShouldTrigger(triggerType) {
		return
	}

	snap := g.FlashChallenge.StartChallenge()
	if snap == nil {
		return
	}

	log.Printf("[FlashChallenge] started: type=%s title=%s target=%d duration=%ds",
		snap.Type, snap.Title, snap.Target, snap.Duration)

	// 廣播給所有玩家
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgFlashChallengeStart,
		Payload: ws.FlashChallengeStartPayload{
			Type:        string(snap.Type),
			Title:       snap.Title,
			Description: snap.Description,
			Icon:        snap.Icon,
			Color:       snap.Color,
			Target:      snap.Target,
			TargetDefID: snap.TargetDefID,
			Duration:    snap.Duration,
			TimeLeft:    snap.TimeLeft,
			BaseReward:  snap.BaseReward,
			BonusReward: snap.BonusReward,
			TopPlayers:  buildFlashPlayerSnaps(snap.TopPlayers),
		},
	})
}

// notifyFlashChallengeKill 在擊破目標後更新閃電挑戰進度（由 handleKill 呼叫）
func (g *Game) notifyFlashChallengeKill(p *player.Player, targetDefID string, mult float64, streak int) {
	if !g.FlashChallenge.IsActive() {
		return
	}

	progress, completed, firstComplete := g.FlashChallenge.RecordKill(p.ID, p.DisplayName, targetDefID, mult, streak)
	if progress == 0 && !completed {
		return // 這次擊破不計入挑戰
	}

	snap := g.FlashChallenge.GetSnapshot()
	if snap == nil {
		return
	}

	// 廣播進度更新（讓所有玩家看到排行榜變化）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgFlashChallengeUpdate,
		Payload: ws.FlashChallengeUpdatePayload{
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Progress:   progress,
			Target:     snap.Target,
			Completed:  completed,
			TimeLeft:   snap.TimeLeft,
			TopPlayers: buildFlashPlayerSnaps(snap.TopPlayers),
		},
	})

	// 首次完成：個人獎勵通知
	if firstComplete {
		reward := g.FlashChallenge.CalcReward(progress, true)
		p.Coins += reward
		if p.Coins > p.MaxCoins {
			p.MaxCoins = p.Coins
		}

		log.Printf("[FlashChallenge] player=%s completed! reward=%d", p.ID, reward)

		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgFlashChallengeReward,
			Payload: ws.FlashChallengeRewardPayload{
				PlayerID:   p.ID,
				Progress:   progress,
				Target:     snap.Target,
				Completed:  true,
				Reward:     reward,
				NewBalance: p.Coins,
				Message:    "🎉 挑戰完成！",
			},
		})

		// 動態牆廣播
		go g.notifyFeedFlashChallenge(p, snap.Title, reward)
	}
}

// tickFlashChallenge 定期檢查閃電挑戰狀態（由 game loop 呼叫）
func (g *Game) tickFlashChallenge() {
	// 檢查是否超時
	if expired := g.FlashChallenge.CheckExpiry(); expired {
		g.handleFlashChallengeEnd()
		return
	}

	// 嘗試隨機觸發新挑戰
	g.tryStartFlashChallenge("random")
}

// handleFlashChallengeEnd 挑戰結束（超時或全員完成）
func (g *Game) handleFlashChallengeEnd() {
	snap := g.FlashChallenge.GetSnapshot()
	if snap == nil {
		return
	}

	// 發放安慰獎給未完成但有進度的玩家
	g.mu.RLock()
	players := make([]*player.Player, 0, len(g.Players))
	for _, p := range g.Players {
		players = append(players, p)
	}
	g.mu.RUnlock()

	completedCount := 0
	for _, p := range players {
		progress, completed := g.FlashChallenge.GetPlayerProgress(p.ID)
		if progress <= 0 {
			continue
		}
		reward := g.FlashChallenge.CalcReward(progress, completed)
		if reward <= 0 {
			continue
		}
		if completed {
			completedCount++
			continue // 已在 notifyFlashChallengeKill 發放
		}
		// 安慰獎
		p.Coins += reward
		if p.Coins > p.MaxCoins {
			p.MaxCoins = p.Coins
		}
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgFlashChallengeReward,
			Payload: ws.FlashChallengeRewardPayload{
				PlayerID:   p.ID,
				Progress:   progress,
				Target:     snap.Target,
				Completed:  false,
				Reward:     reward,
				NewBalance: p.Coins,
				Message:    "💪 努力了！獲得安慰獎",
			},
		})
	}

	// 廣播挑戰結束
	success := completedCount > 0
	message := "挑戰時間到！"
	if success {
		message = "有玩家完成了挑戰！"
	}

	log.Printf("[FlashChallenge] ended: success=%v completedCount=%d", success, completedCount)

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgFlashChallengeEnd,
		Payload: ws.FlashChallengeEndPayload{
			Success:    success,
			Title:      snap.Title,
			Icon:       snap.Icon,
			TopPlayers: buildFlashPlayerSnaps(snap.TopPlayers),
			Message:    message,
		},
	})
}

// handleGetFlashChallenge 查詢閃電挑戰狀態（Client → Server）
func (g *Game) handleGetFlashChallenge(p *player.Player) {
	snap := g.FlashChallenge.GetSnapshot()
	active := snap != nil && snap.State == flashchallenge.StateActive

	myProgress := 0
	myCompleted := false
	if active {
		myProgress, myCompleted = g.FlashChallenge.GetPlayerProgress(p.ID)
	}

	payload := ws.FlashChallengeStatusPayload{
		Active:      active,
		MyProgress:  myProgress,
		MyCompleted: myCompleted,
	}
	if snap != nil {
		payload.Type = string(snap.Type)
		payload.Title = snap.Title
		payload.Description = snap.Description
		payload.Icon = snap.Icon
		payload.Color = snap.Color
		payload.Target = snap.Target
		payload.TargetDefID = snap.TargetDefID
		payload.Duration = snap.Duration
		payload.TimeLeft = snap.TimeLeft
		payload.BaseReward = snap.BaseReward
		payload.BonusReward = snap.BonusReward
		payload.TopPlayers = buildFlashPlayerSnaps(snap.TopPlayers)
	}

	g.Hub.Send(p.ID, &ws.Message{
		Type:    ws.MsgFlashChallengeUpdate,
		Payload: payload,
	})
}

// notifyFeedFlashChallenge 動態牆：閃電挑戰完成（DAY-123）
func (g *Game) notifyFeedFlashChallenge(p *player.Player, title string, reward int) {
	if g.ActivityFeed == nil {
		return
	}
	event := activityfeed.NewMegaWinEvent(p.ID, p.DisplayName, 1.0, reward)
	event.Title = "⚡ 閃電挑戰完成"
	event.Detail = p.DisplayName + " 完成了「" + title + "」，獲得 " + formatCoins(reward) + " 金幣"
	g.broadcastFeedEvent(event)
}

// buildFlashPlayerSnaps 轉換玩家快照格式
func buildFlashPlayerSnaps(snaps []flashchallenge.PlayerSnap) []ws.FlashChallengePlayerSnap {
	result := make([]ws.FlashChallengePlayerSnap, 0, len(snaps))
	for _, s := range snaps {
		result = append(result, ws.FlashChallengePlayerSnap{
			PlayerID:   s.PlayerID,
			PlayerName: s.PlayerName,
			Progress:   s.Progress,
			Completed:  s.Completed,
		})
	}
	return result
}
