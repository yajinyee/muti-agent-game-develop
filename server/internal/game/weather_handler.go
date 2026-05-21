// weather_handler.go — 天氣系統 handler（DAY-087）
package game

import (
	"log"

	"digital-twin/server/internal/ws"
)

// sendWeatherUpdate 發送天氣狀態給單一玩家
func (g *Game) sendWeatherUpdate(playerID string, isNew bool) {
	snap := g.Weather.GetSnapshot(isNew)
	if err := g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgWeatherUpdate,
		Payload: ws.WeatherUpdatePayload{
			Type:             string(snap.Type),
			Name:             snap.Name,
			Icon:             snap.Icon,
			Description:      snap.Description,
			RemainingSeconds: snap.RemainingSeconds,
			SpawnRateMult:    snap.SpawnRateMult,
			RewardMult:       snap.RewardMult,
			SpeedMult:        snap.SpeedMult,
			RareChanceBonus:  snap.RareChanceBonus,
			GoldFishBonus:    snap.GoldFishBonus,
			BossChanceBonus:  snap.BossChanceBonus,
			FogEffect:        snap.FogEffect,
			IsNew:            isNew,
		},
	}); err != nil {
		log.Printf("[Weather] send update to %s error: %v", playerID, err)
	}
}

// broadcastWeatherUpdate 廣播天氣狀態給所有玩家
func (g *Game) broadcastWeatherUpdate(isNew bool) {
	snap := g.Weather.GetSnapshot(isNew)
	payload := ws.WeatherUpdatePayload{
		Type:             string(snap.Type),
		Name:             snap.Name,
		Icon:             snap.Icon,
		Description:      snap.Description,
		RemainingSeconds: snap.RemainingSeconds,
		SpawnRateMult:    snap.SpawnRateMult,
		RewardMult:       snap.RewardMult,
		SpeedMult:        snap.SpeedMult,
		RareChanceBonus:  snap.RareChanceBonus,
		GoldFishBonus:    snap.GoldFishBonus,
		BossChanceBonus:  snap.BossChanceBonus,
		FogEffect:        snap.FogEffect,
		IsNew:            isNew,
	}
	g.Hub.Broadcast(&ws.Message{
		Type:    ws.MsgWeatherUpdate,
		Payload: payload,
	})
}

// tickAndBroadcastWeather 檢查天氣是否需要切換，並廣播（由 gameLoop 每 30 秒呼叫）
func (g *Game) tickAndBroadcastWeather() {
	changed, snap := g.Weather.CheckAndRotate()
	if changed {
		log.Printf("[Weather] changed to %s (%s) — reward×%.1f, speed×%.1f",
			snap.Name, snap.Icon, snap.RewardMult, snap.SpeedMult)
		g.broadcastWeatherUpdate(true)
		// 全服公告：天氣變化（DAY-097）
		g.announceWeatherChange(snap.Name)
		// 天氣湧現事件：特定天氣觸發稀有目標群湧（DAY-127）
		go g.tryTriggerWeatherSurge(snap.Type)
	}
}
