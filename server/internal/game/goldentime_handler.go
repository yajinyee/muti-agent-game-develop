// goldentime_handler.go — 黃金時間系統 handler（DAY-125）
// 業界依據：Fire Kirin / Ocean King 系列的 Golden Time 機制
// 全場目標物倍率暫時提升，製造「全場瘋狂」的高峰體驗
// 業界研究顯示 Golden Time 讓短期參與度提升 40%+
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/goldentime"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// triggerGoldenTime 觸發黃金時間（由 boss_handler / raid_handler / flashchallenge_handler 呼叫）
func (g *Game) triggerGoldenTime(trigger goldentime.TriggerType) {
	if !g.GoldenTime.CanTrigger() {
		return
	}
	tier := goldentime.SelectTier(trigger)
	session := g.GoldenTime.Start(tier, trigger)
	if session == nil {
		return
	}
	def := goldentime.TierDefs[tier]
	log.Printf("[GoldenTime] triggered tier=%s mult=%.1fx duration=%ds trigger=%s",
		def.Name, def.MultBoost, def.Duration, trigger)

	// 廣播黃金時間開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGoldenTimeStart,
		Payload: ws.GoldenTimeStartPayload{
			Tier:        int(tier),
			TierName:    def.Name,
			MultBoost:   def.MultBoost,
			Duration:    def.Duration,
			SecondsLeft: def.Duration,
			Icon:        def.Icon,
			Color:       def.Color,
			BgColor:     def.BgColor,
			TriggerType: string(trigger),
		},
	})

	// 全服公告
	multStr := fmt.Sprintf("%.1f", def.MultBoost)
	ann := g.Announce.Create(announce.EventGoldenTime, "", 0, map[string]string{
		"tier_name": def.Name,
		"mult":      multStr,
	})
	g.broadcastAnnouncement(ann)
}

// tickGoldenTime 檢查黃金時間是否結束（由 gameLoop 每秒呼叫）
func (g *Game) tickGoldenTime() {
	// 檢查是否剛剛結束
	if g.GoldenTime.CheckExpiry() {
		log.Printf("[GoldenTime] session ended")
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgGoldenTimeEnd,
			Payload: ws.GoldenTimeEndPayload{
				Message: "⏰ 黃金時間結束！",
			},
		})
		return
	}

	// 隨機觸發（每次 tick 有 0.5% 機率）
	if g.GoldenTime.ShouldTriggerRandom() {
		go g.triggerGoldenTime(goldentime.TriggerRandom)
	}
}

// handleGetGoldenTime 處理玩家查詢黃金時間狀態
func (g *Game) handleGetGoldenTime(p *player.Player) {
	snap := g.GoldenTime.GetSnapshot()
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgGoldenTimeStatus,
		Payload: ws.GoldenTimeStatusPayload{
			IsActive:    snap.IsActive,
			Tier:        snap.Tier,
			TierName:    snap.TierName,
			MultBoost:   snap.MultBoost,
			SecondsLeft: snap.SecondsLeft,
			Icon:        snap.Icon,
			Color:       snap.Color,
			BgColor:     snap.BgColor,
			TriggerType: snap.TriggerType,
		},
	}); err != nil {
		log.Printf("[GoldenTime] send status to %s error: %v", p.ID, err)
	}
}
