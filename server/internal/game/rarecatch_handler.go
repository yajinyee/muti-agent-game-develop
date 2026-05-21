// rarecatch_handler.go — 稀有連擊累積倍率系統 handler（DAY-126）
// 業界依據：fishingfortune.app（2026-05-21）確認「multiplier cascade system」
// 連續在 90 秒內擊破稀有目標（T101-T105），倍率從 2x 累積到最高 15x
// 業界研究顯示稀有目標專屬倍率讓玩家主動追求高價值目標，提升策略深度
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/game/rarecatch"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// notifyRareCatchKill 在擊破稀有目標後更新稀有連擊（由 handleKill 呼叫）
// 回傳稀有連擊倍率加成（用於最終獎勵計算）
func (g *Game) notifyRareCatchKill(p *player.Player, defID string) float64 {
	if !rarecatch.IsRareTarget(defID) {
		return 1.0
	}

	count, multBoost, isLevelUp, shouldBroadcast := g.RareCatch.RecordKill(p.ID)
	snap := g.RareCatch.GetSnapshot(p.ID)

	log.Printf("[RareCatch] player=%s defID=%s count=%d mult=%.1fx levelUp=%v",
		p.ID, defID, count, multBoost, isLevelUp)

	// 發送個人更新
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgRareCatchUpdate,
		Payload: ws.RareCatchUpdatePayload{
			Count:       count,
			MultBoost:   multBoost,
			LevelName:   snap.LevelName,
			Icon:        snap.Icon,
			Color:       snap.Color,
			SecondsLeft: snap.SecondsLeft,
			IsLevelUp:   isLevelUp,
		},
	}); err != nil {
		log.Printf("[RareCatch] send update error: %v", err)
	}

	// 達到 ×5.0 以上時全服廣播
	if shouldBroadcast {
		multStr := fmt.Sprintf("%.0f", multBoost)
		msg := fmt.Sprintf("%s 達成稀有連擊 ×%s！", p.DisplayName, multStr)
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgRareCatchBroadcast,
			Payload: ws.RareCatchBroadcastPayload{
				PlayerID:   p.ID,
				PlayerName: p.DisplayName,
				Count:      count,
				MultBoost:  multBoost,
				LevelName:  snap.LevelName,
				Icon:       snap.Icon,
				Color:      snap.Color,
				Message:    msg,
			},
		})
		log.Printf("[RareCatch] broadcast: %s", msg)
	}

	return multBoost
}

// tickRareCatchExpiry 定期檢查稀有連擊超時（由 game loop 每 5 秒呼叫）
func (g *Game) tickRareCatchExpiry() {
	expired := g.RareCatch.CheckExpiry()
	for _, playerID := range expired {
		// 通知玩家連擊已重置
		if err := g.Hub.Send(playerID, &ws.Message{
			Type: ws.MsgRareCatchReset,
			Payload: ws.RareCatchResetPayload{
				FinalCount: 0,
				Message:    "稀有連擊超時，倍率重置",
			},
		}); err != nil {
			// 玩家可能已離線，忽略錯誤
			log.Printf("[RareCatch] send reset to %s error: %v", playerID, err)
		}
	}
}
