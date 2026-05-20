// milestone_handler.go — 登入里程碑獎勵 handler（DAY-107）
// 業界依據：ilogos.biz（2026）確認 gamified login streaks 讓留存率提升 75%
// 設計：連續登入達到 3/7/14/30/60/100 天時給予特殊獎勵（寶箱 + 金幣 + 稱號）
package game

import (
	"log"

	"digital-twin/server/internal/game/achievement"
	"digital-twin/server/internal/game/dailybonus"
	"digital-twin/server/internal/game/mysterybox"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// checkAndGrantLoginMilestone 檢查並發放登入里程碑獎勵
// 在 checkAndSendDailyBonus 確認是新的一天後呼叫
func (g *Game) checkAndGrantLoginMilestone(p *player.Player, newStreak int) {
	milestone := dailybonus.CheckMilestone(newStreak)
	if milestone == nil {
		return // 未達到里程碑
	}

	log.Printf("[Milestone] Player %s reached day %d milestone: %s %s",
		p.ID, milestone.Days, milestone.Icon, milestone.Name)

	// 計算並發放獎勵
	coinsGained := 0
	var rewardPayloads []ws.MilestoneRewardPayload

	for _, reward := range milestone.Rewards {
		rp := ws.MilestoneRewardPayload{
			Type:    string(reward.Type),
			Amount:  reward.Amount,
			Rarity:  reward.Rarity,
			TitleID: reward.TitleID,
		}

		switch reward.Type {
		case dailybonus.MilestoneRewardCoins:
			// 發放金幣
			p.AddCoins(reward.Amount)
			coinsGained += reward.Amount
			log.Printf("[Milestone] Player %s received %d coins", p.ID, reward.Amount)

		case dailybonus.MilestoneRewardMysteryBox:
			// 發放神秘寶箱
			for i := 0; i < reward.Amount; i++ {
				g.MysteryBox.AddBox(p.ID, mysterybox.BoxRarity(reward.Rarity))
			}
			log.Printf("[Milestone] Player %s received %d %s mystery box(es)",
				p.ID, reward.Amount, reward.Rarity)

		case dailybonus.MilestoneRewardTitle:
			// 解鎖特殊稱號
			if reward.TitleID != "" {
				titleID := achievement.TitleID(reward.TitleID)
				if titleDef := p.Titles.TryUnlockByID(titleID); titleDef != nil {
					log.Printf("[Milestone] Player %s unlocked title: %s", p.ID, titleDef.Name)
					// 通知稱號解鎖
					g.Hub.Send(p.ID, &ws.Message{
						Type: ws.MsgTitleUnlocked,
						Payload: ws.TitleUnlockedPayload{
							TitleID:     string(titleDef.ID),
							TitleName:   titleDef.Name,
							TitleIcon:   titleDef.Icon,
							TitleColor:  titleDef.Color,
							Description: titleDef.Description,
						},
					})
				}
			}
		}

		rewardPayloads = append(rewardPayloads, rp)
	}

	// 通知玩家里程碑達成
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLoginMilestone,
		Payload: ws.LoginMilestonePayload{
			Days:        milestone.Days,
			Name:        milestone.Name,
			Description: milestone.Description,
			Icon:        milestone.Icon,
			Color:       milestone.Color,
			Rewards:     rewardPayloads,
			CoinsGained: coinsGained,
			NewBalance:  p.GetCoins(),
		},
	})

	// 同步更新玩家狀態
	g.sendPlayerUpdate(p)

	// 更新神秘寶箱狀態（如果有發放寶箱）
	for _, reward := range milestone.Rewards {
		if reward.Type == dailybonus.MilestoneRewardMysteryBox {
			g.sendMysteryBoxUpdate(p)
			break
		}
	}
	// 動態牆：登入里程碑（DAY-112）
	go g.notifyFeedMilestone(p, milestone.Days, milestone.Name)
}

// handleGetLoginProgress 處理查詢登入進度請求（DAY-107）
func (g *Game) handleGetLoginProgress(p *player.Player) {
	g.mu.RLock()
	currentStreak := p.LoginStreak
	maxStreak := p.MaxLoginStreak
	g.mu.RUnlock()

	// 取得所有里程碑
	allMilestones := dailybonus.GetAllMilestones()
	milestonePayloads := make([]ws.MilestoneInfoPayload, 0, len(allMilestones))

	for _, m := range allMilestones {
		rewards := make([]ws.MilestoneRewardPayload, 0, len(m.Rewards))
		for _, r := range m.Rewards {
			rewards = append(rewards, ws.MilestoneRewardPayload{
				Type:    string(r.Type),
				Amount:  r.Amount,
				Rarity:  r.Rarity,
				TitleID: r.TitleID,
			})
		}
		milestonePayloads = append(milestonePayloads, ws.MilestoneInfoPayload{
			Days:      m.Days,
			Name:      m.Name,
			Icon:      m.Icon,
			Color:     m.Color,
			Rewards:   rewards,
			IsReached: currentStreak >= m.Days,
		})
	}

	// 計算下一個里程碑
	nextMilestoneDays := 0
	daysToNext := 0
	if next := dailybonus.GetNextMilestone(currentStreak); next != nil {
		nextMilestoneDays = next.Days
		daysToNext = next.Days - currentStreak
	}

	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLoginProgress,
		Payload: ws.LoginProgressPayload{
			CurrentStreak:     currentStreak,
			MaxStreak:         maxStreak,
			NextMilestoneDays: nextMilestoneDays,
			DaysToNext:        daysToNext,
			Milestones:        milestonePayloads,
		},
	})
}
