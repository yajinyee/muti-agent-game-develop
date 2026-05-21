// winstreak_handler.go — 連勝獎勵系統 handler（DAY-131）
// 業界依據：BGaming Fishing Club 2026 Best Win/Best Catch 里程碑
// 追蹤玩家連勝次數，達到里程碑（10/25/50/100次）時給予額外獎勵。
// 與連擊系統（2秒內連續）不同，這是更長時間的累積（30秒超時重置）。
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/winstreak"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// notifyWinStreakKill 在擊破目標後更新連勝（由 handleKill 呼叫）
func (g *Game) notifyWinStreakKill(p *player.Player) {
	if g.WinStreak == nil {
		return
	}

	current, milestone, wasReset := g.WinStreak.RecordKill(p.ID)

	if wasReset && current == 1 {
		// 連勝被重置後的第一次擊破，不需要特別通知
	}

	// 發送連勝更新（個人）
	snap := g.WinStreak.GetSnapshot(p.ID)
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgWinStreakUpdate,
		Payload: ws.WinStreakUpdatePayload{
			Current:           current,
			MaxStreak:         snap.MaxStreak,
			NextMilestone:     int(snap.NextMilestone),
			NextMilestoneName: snap.NextMilestoneName,
			ProgressToNext:    snap.ProgressToNext,
			SecondsToExpiry:   snap.SecondsToExpiry,
		},
	}); err != nil {
		log.Printf("[WinStreak] send update error: %v", err)
	}

	// 里程碑達成
	if milestone != nil {
		betDef := p.GetBetDef()
		bonusReward := int(float64(betDef.BetCost) * milestone.BonusMult)
		p.AddCoins(bonusReward)

		log.Printf("[WinStreak] player=%s reached %s (streak=%d), bonus=%d",
			p.ID, milestone.Name, current, bonusReward)

		// 發送里程碑通知（個人）
		if err := g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgWinStreakMilestone,
			Payload: ws.WinStreakMilestonePayload{
				PlayerID:    p.ID,
				PlayerName:  p.DisplayName,
				Streak:      current,
				Level:       int(milestone.Level),
				LevelName:   milestone.Name,
				Icon:        milestone.Icon,
				Color:       milestone.Color,
				BonusReward: bonusReward,
				NewBalance:  p.GetCoins(),
				Broadcast:   milestone.Broadcast,
			},
		}); err != nil {
			log.Printf("[WinStreak] send milestone error: %v", err)
		}

		// 高等級里程碑全服廣播（金牌/傳說）
		if milestone.Broadcast {
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgWinStreakMilestone,
				Payload: ws.WinStreakMilestonePayload{
					PlayerID:    p.ID,
					PlayerName:  p.DisplayName,
					Streak:      current,
					Level:       int(milestone.Level),
					LevelName:   milestone.Name,
					Icon:        milestone.Icon,
					Color:       milestone.Color,
					BonusReward: bonusReward,
					Broadcast:   true,
				},
			})

			// 全服公告
			ann := g.Announce.Create(announce.EventBigWin, p.DisplayName, bonusReward, map[string]string{
				"message": fmt.Sprintf("%s %s 達成 %s！連勝 %d 次！獎勵 %d 金幣！", milestone.Icon, p.DisplayName, milestone.Name, current, bonusReward),
			})
			g.broadcastAnnouncement(ann)

			// 動態牆：傳說連勝
			if milestone.Level == winstreak.MilestoneLegend {
				go g.notifyFeedMegaWin(p, float64(milestone.BonusMult), bonusReward)
			}
		}
	}
}

// tickWinStreakExpiry 定期檢查連勝超時（由 gameLoop 每秒呼叫）
func (g *Game) tickWinStreakExpiry() {
	if g.WinStreak == nil {
		return
	}

	expiredIDs := g.WinStreak.CheckExpiry()
	for _, playerID := range expiredIDs {
		g.mu.RLock()
		p := g.Players[playerID]
		g.mu.RUnlock()

		if p == nil {
			continue
		}

		snap := g.WinStreak.GetSnapshot(playerID)
		if err := g.Hub.Send(playerID, &ws.Message{
			Type: ws.MsgWinStreakReset,
			Payload: ws.WinStreakResetPayload{
				FinalStreak: 0,
				MaxStreak:   snap.MaxStreak,
			},
		}); err != nil {
			log.Printf("[WinStreak] send reset error: %v", err)
		}
	}
}
