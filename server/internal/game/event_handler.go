// event_handler.go 限時活動系統 handler（DAY-079）
package game

import (
	"log"

	"digital-twin/server/internal/game/event"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// sendEventUpdate 發送限時活動狀態給指定玩家
func (g *Game) sendEventUpdate(p *player.Player) {
	snap := g.Event.GetSnapshot()
	g.Hub.Send(p.ID, &ws.Message{
		Type:    ws.MsgEventUpdate,
		Payload: eventSnapshotToPayload(snap),
	})
}

// handleGetEventStatus 處理查詢限時活動狀態請求
func (g *Game) handleGetEventStatus(p *player.Player) {
	g.sendEventUpdate(p)
}

// tickAndBroadcastEvent 定期 Tick 活動並廣播（在 gameLoop 中呼叫）
func (g *Game) tickAndBroadcastEvent() {
	changed := g.Event.Tick()
	snap := g.Event.GetSnapshot()
	payload := eventSnapshotToPayload(snap)

	if changed {
		if snap.IsActive {
			log.Printf("[Event] New event started: %s (%s)", snap.Name, snap.Type)
		} else {
			log.Printf("[Event] Event ended, no active event")
		}
	}

	// 廣播給所有玩家
	g.Hub.Broadcast(&ws.Message{
		Type:    ws.MsgEventUpdate,
		Payload: payload,
	})
}

// getEventKillChanceAdd 取得當前活動的擊破率加成（供 combat 使用）
func (g *Game) getEventKillChanceAdd() float64 {
	return g.Event.GetKillChanceAdd()
}

// getEventSpawnMult 取得當前活動的目標生成倍率（供 spawnTarget 使用）
func (g *Game) getEventSpawnMult() float64 {
	return g.Event.GetSpawnMult()
}

// eventSnapshotToPayload 將活動快照轉換為 WebSocket Payload
func eventSnapshotToPayload(snap event.EventSnapshot) ws.EventUpdatePayload {
	return ws.EventUpdatePayload{
		Type:          snap.Type,
		Name:          snap.Name,
		Description:   snap.Description,
		Icon:          snap.Icon,
		Color:         snap.Color,
		IsActive:      snap.IsActive,
		EndAt:         snap.EndAt,
		TimeLeft:      snap.TimeLeft,
		RewardMult:    snap.RewardMult,
		SpawnMult:     snap.SpawnMult,
		KillChanceAdd: snap.KillChanceAdd,
	}
}
