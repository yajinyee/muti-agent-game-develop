// season_handler.go — 賽季通行證 handler（DAY-072）
// 整合到 game.go 的 HandleMessage 和積分增加流程
package game

import (
	"log"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// sendSeasonUpdate 發送賽季通行證更新給玩家
func (g *Game) sendSeasonUpdate(p *player.Player) {
	snap := g.Season.GetSnapshot(p.ID)
	levels := make([]ws.SeasonLevelStatus, len(snap.Levels))
	for i, l := range snap.Levels {
		levels[i] = ws.SeasonLevelStatus{
			Level:        l.Level,
			PointsNeeded: l.PointsNeeded,
			CoinReward:   l.CoinReward,
			SpecialType:  l.SpecialType,
			SpecialID:    l.SpecialID,
			SpecialName:  l.SpecialName,
			Icon:         l.Icon,
			Claimed:      l.Claimed,
			Unlocked:     l.Unlocked,
		}
	}
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgSeasonUpdate,
		Payload: ws.SeasonUpdatePayload{
			PlayerID:     p.ID,
			SeasonPoints: snap.SeasonPoints,
			CurrentLevel: snap.CurrentLevel,
			NextLevel:    snap.NextLevel,
			PointsToNext: snap.PointsToNext,
			Progress:     snap.Progress,
			Levels:       levels,
		},
	})
}

// addSeasonPoints 增加賽季積分（在週賽積分增加時同步呼叫）
// 回傳可領取的新等級列表
func (g *Game) addSeasonPoints(playerID string, points int) []int {
	_, newLevels := g.Season.AddPoints(playerID, points)
	return newLevels
}

// handleClaimSeasonLevel 處理領取賽季等級獎勵（DAY-072）
func (g *Game) handleClaimSeasonLevel(p *player.Player, msg *ws.Message) {
	var payload ws.ClaimSeasonLevelPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	result, ok := g.Season.ClaimLevel(p.ID, payload.Level)
	if !ok {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: map[string]string{"message": "無法領取此等級獎勵（積分不足或已領取）"},
		})
		return
	}

	// 發放金幣獎勵
	p.AddCoins(result.CoinReward)

	// 處理特殊獎勵
	switch result.SpecialType {
	case "skin":
		// 解鎖賽季限定皮膚（免費加入擁有列表）
		if !p.BuySkin(result.SpecialID, 0) {
			// 已擁有，直接裝備
			p.EquipSkin(result.SpecialID)
		} else {
			p.EquipSkin(result.SpecialID)
		}
		// 通知皮膚更新
		equippedSkin, ownedSkins := p.GetSkinInfo()
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgSkinUpdate,
			Payload: ws.SkinUpdatePayload{
				PlayerID:     p.ID,
				EquippedSkin: equippedSkin,
				OwnedSkins:   ownedSkins,
				NewBalance:   p.Coins,
			},
		})
		log.Printf("[Season] 玩家 %s 解鎖賽季皮膚 %s", p.ID, result.SpecialID)

	case "title":
		// 解鎖賽季稱號（加入成就系統）
		// 直接廣播稱號解鎖通知
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgTitleUnlocked,
			Payload: ws.TitleUnlockedPayload{
				TitleID:     result.SpecialID,
				TitleName:   result.SpecialName,
				TitleIcon:   "👑",
				TitleColor:  "#FFD700",
				Description: "賽季通行證等級 10 獎勵",
			},
		})
		log.Printf("[Season] 玩家 %s 解鎖賽季稱號 %s", p.ID, result.SpecialID)
	}

	// 通知等級升級
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgSeasonLevelUp,
		Payload: ws.SeasonLevelUpPayload{
			PlayerID:    p.ID,
			Level:       result.NewLevel,
			CoinReward:  result.CoinReward,
			NewBalance:  p.Coins,
			SpecialType: result.SpecialType,
			SpecialID:   result.SpecialID,
			SpecialName: result.SpecialName,
		},
	})

	// 更新玩家狀態（讓 Client 看到新金幣餘額）
	g.sendPlayerUpdate(p)

	// 更新賽季快照
	g.sendSeasonUpdate(p)

	log.Printf("[Season] 玩家 %s 領取賽季等級 %d 獎勵（%d 金幣）",
		p.ID, result.NewLevel, result.CoinReward)
}

// seasonPointsFromTournamentPoints 從週賽積分計算賽季積分
// 賽季積分 = 週賽積分（1:1 比例）
func seasonPointsFromTournamentPoints(tournamentPoints int) int {
	return tournamentPoints
}

// checkSeasonLevelNotify 檢查是否有新等級可領取，發送通知
func (g *Game) checkSeasonLevelNotify(p *player.Player, newLevels []int) {
	if len(newLevels) == 0 {
		return
	}
	// 發送賽季更新（讓 Client 知道有新等級可領取）
	g.sendSeasonUpdate(p)
	log.Printf("[Season] 玩家 %s 有 %d 個新等級可領取: %v", p.ID, len(newLevels), newLevels)
	// 動態牆：賽季升級（DAY-112，只廣播最高等級）
	if len(newLevels) > 0 {
		maxLevel := newLevels[len(newLevels)-1]
		go g.notifyFeedSeasonLevel(p, maxLevel)
	}
}
