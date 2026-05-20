// vip_handler.go VIP 等級系統 handler（DAY-078）
package game

import (
	"log"

	"digital-twin/server/internal/game/vip"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// sendVIPUpdate 發送 VIP 狀態給指定玩家
func (g *Game) sendVIPUpdate(p *player.Player) {
	snap := g.VIP.GetSnapshot(p.ID)
	g.Hub.Send(p.ID, &ws.Message{
		Type:    ws.MsgVIPUpdate,
		Payload: vipSnapshotToPayload(snap),
	})
}

// handleGetVIPStatus 處理查詢 VIP 狀態請求
func (g *Game) handleGetVIPStatus(p *player.Player) {
	g.sendVIPUpdate(p)
}

// handleClaimVIPWeekly 處理領取 VIP 週獎勵請求
func (g *Game) handleClaimVIPWeekly(p *player.Player) {
	result := g.VIP.ClaimWeeklyBonus(p.ID)
	if result == nil {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "vip_weekly_unavailable", Message: "週獎勵尚未可領取或您不是 VIP 會員"},
		})
		return
	}

	// 發放獎勵
	p.AddCoins(result.Coins)

	log.Printf("[VIP] Player %s claimed weekly bonus: level=%d, coins=%d", p.ID, result.VIPLevel, result.Coins)

	// 通知玩家
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgVIPWeeklyClaimed,
		Payload: ws.VIPWeeklyClaimedPayload{
			PlayerID:   p.ID,
			VIPLevel:   result.VIPLevel,
			TierName:   result.TierName,
			Coins:      result.Coins,
			NewBalance: p.Coins,
		},
	})

	// 更新 VIP 狀態
	g.sendVIPUpdate(p)
}

// notifyVIPSpend 記錄 VIP 消費並處理升級/返還（每次攻擊後呼叫）
func (g *Game) notifyVIPSpend(p *player.Player, spendAmount int) {
	// 計算返還（在升級前計算，用舊等級）
	cashback := g.VIP.GetCashback(p.ID, spendAmount)

	// 記錄消費，檢查是否升級
	newLevel, levelUp := g.VIP.AddSpend(p.ID, spendAmount)
	_ = newLevel

	// 發放返還金幣（若有）
	if cashback > 0 {
		p.AddCoins(cashback)
	}

	// 若升級，發送升級通知
	if levelUp != nil {
		log.Printf("[VIP] Player %s leveled up to VIP %d (%s)", p.ID, levelUp.NewLevel, levelUp.TierName)
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgVIPLevelUp,
			Payload: ws.VIPLevelUpPayload{
				PlayerID:    p.ID,
				NewLevel:    levelUp.NewLevel,
				TierName:    levelUp.TierName,
				TierIcon:    levelUp.TierIcon,
				TierColor:   levelUp.TierColor,
				TitleID:     levelUp.TitleID,
				TitleName:   levelUp.TitleName,
				WeeklyBonus: levelUp.WeeklyBonus,
			},
		})
		// 更新 VIP 狀態
		g.sendVIPUpdate(p)
	}
}

// applyVIPDailyBonusMult 取得 VIP 每日登入獎勵倍率（供 dailybonus_handler 使用）
func (g *Game) applyVIPDailyBonusMult(playerID string) float64 {
	return g.VIP.GetDailyBonusMult(playerID)
}

// vipSnapshotToPayload 將 VIP 快照轉換為 WebSocket Payload
func vipSnapshotToPayload(snap vip.VIPSnapshot) ws.VIPUpdatePayload {
	return ws.VIPUpdatePayload{
		PlayerID:       snap.PlayerID,
		TotalSpend:     snap.TotalSpend,
		VIPLevel:       snap.VIPLevel,
		TierName:       snap.TierName,
		TierIcon:       snap.TierIcon,
		TierColor:      snap.TierColor,
		CashbackRate:   snap.CashbackRate,
		DailyBonusMult: snap.DailyBonusMult,
		WeeklyBonus:    snap.WeeklyBonus,
		NextLevel:      snap.NextLevel,
		SpendToNext:    snap.SpendToNext,
		Progress:       snap.Progress,
		CanClaimWeekly: snap.CanClaimWeekly,
	}
}
