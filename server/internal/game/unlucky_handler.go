// unlucky_handler.go — 失敗補償系統 handler（DAY-135）
// 業界依據：Funrize 2026 的「Unlucky Bonus」
// 連續花費超過一定金額但獲得低回報時，自動給予補償獎勵
// 防止玩家因為「運氣太差」而離開，是 2026 年業界最新的留存機制
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// notifyUnluckyShot 記錄射擊，若觸發失敗補償則發放獎勵（由 handleAttack 呼叫）
// spend: 本次射擊花費（betCost）
// reward: 本次射擊獲得的獎勵（0 = 未擊破）
func (g *Game) notifyUnluckyShot(p *player.Player, spend int, reward int) {
	if g.UnluckyBonus == nil {
		return
	}

	triggered, bonusAmount := g.UnluckyBonus.RecordShot(p.ID, spend, reward)
	if !triggered {
		return
	}

	// 發放補償獎勵
	p.AddCoins(bonusAmount)

	log.Printf("[UnluckyBonus] player=%s triggered! bonus=%d, balance=%d",
		p.ID, bonusAmount, p.GetCoins())

	// 通知玩家
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgUnluckyBonus,
		Payload: ws.UnluckyBonusPayload{
			PlayerID:   p.ID,
			BonusAmount: bonusAmount,
			NewBalance: p.GetCoins(),
			Message:    fmt.Sprintf("🍀 運氣補償！獲得 %d 金幣！繼續加油！", bonusAmount),
		},
	}); err != nil {
		log.Printf("[UnluckyBonus] send notify error: %v", err)
	}

	// 全服公告（補償金額較大時）
	if bonusAmount >= 500 {
		ann := g.Announce.Create(announce.EventUnluckyBonus, p.DisplayName, bonusAmount, map[string]string{
			"message": fmt.Sprintf("🍀 %s 獲得運氣補償 %d 金幣！", p.DisplayName, bonusAmount),
		})
		g.broadcastAnnouncement(ann)
	}
}

// sendUnluckyBonusStatus 發送失敗補償狀態給玩家（登入時呼叫）
func (g *Game) sendUnluckyBonusStatus(p *player.Player) {
	if g.UnluckyBonus == nil {
		return
	}

	snap := g.UnluckyBonus.GetSnapshot(p.ID)
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgUnluckyBonusStatus,
		Payload: ws.UnluckyBonusStatusPayload{
			PlayerID:     p.ID,
			ShotCount:    snap.ShotCount,
			TrackingMax:  snap.TrackingMax,
			NetLoss:      snap.NetLoss,
			RatioPercent: snap.RatioPercent,
			CooldownLeft: snap.CooldownLeft,
			BonusCount:   snap.BonusCount,
		},
	}); err != nil {
		log.Printf("[UnluckyBonus] send status error: %v", err)
	}
}
