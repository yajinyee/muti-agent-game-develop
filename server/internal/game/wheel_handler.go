// wheel_handler.go - DAY-084 幸運轉盤 handler
package game

import (
	"log"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/wheel"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// notifyWheelKill 在擊破目標後判斷是否觸發轉盤（由 handleKill 呼叫）
// 若觸發，執行轉盤並發放額外獎勵，回傳額外獎勵金額
func (g *Game) notifyWheelKill(p *player.Player, defID string, baseReward int) int {
	if g.Wheel == nil {
		return 0
	}
	if !g.Wheel.ShouldTrigger(defID) {
		return 0
	}

	// 執行轉盤
	slotIndex, slot := g.Wheel.Spin()
	finalReward := int(float64(baseReward) * slot.Multiplier)
	extraReward := finalReward - baseReward

	// 發放額外獎勵
	p.AddCoins(extraReward)

	// 建立格子列表
	slots := make([]ws.WheelSlotPayload, len(wheel.Slots))
	for i, s := range wheel.Slots {
		slots[i] = ws.WheelSlotPayload{
			Multiplier: s.Multiplier,
			Label:      s.Label,
			Color:      s.Color,
		}
	}

	// 取得目標名稱
	targetName := defID
	if def, ok := data.Targets[defID]; ok {
		targetName = def.Name
	}

	// 廣播轉盤結果（只給觸發的玩家）
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgWheelTrigger,
		Payload: ws.WheelTriggerPayload{
			PlayerID:    p.ID,
			TargetID:    defID,
			TargetName:  targetName,
			Slots:       slots,
			WinIndex:    slotIndex,
			Multiplier:  slot.Multiplier,
			BaseReward:  baseReward,
			FinalReward: finalReward,
			NewBalance:  p.GetCoins(),
		},
	}); err != nil {
		log.Printf("[Wheel] send trigger error: %v", err)
	}

	log.Printf("[Wheel] player=%s triggered wheel on %s: %.0fx, base=%d, final=%d",
		p.ID, defID, slot.Multiplier, baseReward, finalReward)

	// 隱藏挑戰：轉盤 100x（DAY-085）
	g.notifyChallengeWheel(p, slot.Multiplier)

	return extraReward
}
