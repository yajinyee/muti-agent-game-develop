// multstorm_handler.go — 全服倍率風暴系統 handler（DAY-138）
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/ws"
)

// tickMultStorm 倍率風暴 tick（每秒呼叫一次）
// 嘗試觸發風暴 + 檢查過期
func (g *Game) tickMultStorm() {
	if g.MultStorm == nil {
		return
	}

	// 嘗試觸發新風暴
	if sess := g.MultStorm.TryTrigger(); sess != nil {
		log.Printf("[MultStorm] Storm triggered: %s (×%.1f, %.0fs)",
			sess.Tier.Name, sess.Tier.MultBoost, sess.Tier.Duration)

		// 全服廣播風暴開始
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgMultStormStart,
			Payload: ws.MultStormStartPayload{
				TierName:    sess.Tier.Name,
				TierIcon:    sess.Tier.Icon,
				TierColor:   sess.Tier.Color,
				MultBoost:   sess.Tier.MultBoost,
				SecondsLeft: sess.Tier.Duration,
				Message:     fmt.Sprintf("%s 全場倍率 ×%.0f！持續 %.0f 秒！", sess.Tier.Name, sess.Tier.MultBoost, sess.Tier.Duration),
			},
		})

		// 全服公告
		ann := g.Announce.Create(announce.EventMultStorm, "", 0, map[string]string{
			"tier_name":  sess.Tier.Name,
			"tier_icon":  sess.Tier.Icon,
			"mult_boost": fmt.Sprintf("%.0f", sess.Tier.MultBoost),
			"duration":   fmt.Sprintf("%.0f", sess.Tier.Duration),
		})
		g.broadcastAnnouncement(ann)
	}

	// 檢查風暴是否結束
	if g.MultStorm.CheckExpiry() {
		log.Printf("[MultStorm] Storm ended")
		g.Hub.Broadcast(&ws.Message{
			Type:    ws.MsgMultStormEnd,
			Payload: ws.MultStormEndPayload{Message: "倍率風暴結束，回歸正常！"},
		})
	}
}

// sendMultStormStatus 登入時發送風暴狀態
func (g *Game) sendMultStormStatus(playerID string) {
	if g.MultStorm == nil {
		return
	}
	snap := g.MultStorm.GetSnapshot()
	if !snap.IsActive {
		return
	}
	g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgMultStormStart,
		Payload: ws.MultStormStartPayload{
			TierName:    snap.TierName,
			TierIcon:    snap.TierIcon,
			TierColor:   snap.TierColor,
			MultBoost:   snap.MultBoost,
			SecondsLeft: snap.SecondsLeft,
			Message:     fmt.Sprintf("%s 進行中！全場倍率 ×%.0f！", snap.TierName, snap.MultBoost),
		},
	})
}

// getMultStormBoost 取得當前風暴倍率加成（供 handleKill 使用）
func (g *Game) getMultStormBoost() float64 {
	if g.MultStorm == nil {
		return 1.0
	}
	return g.MultStorm.GetMultBoost()
}
