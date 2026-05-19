// Package game — Client 效能上報 handler（DAY-057 拆分自 game.go）
package game

import (
	"log"

	"digital-twin/server/internal/ws"
)

// handleClientPerf 處理 Client 端效能數據上報（DAY-045）
// Client 每 30 秒發送一次，Server 記錄並暴露到 /metrics
// 同時檢查高延遲玩家並輸出警告 log
func (g *Game) handleClientPerf(clientID string, msg *ws.Message) {
	var payload ws.ClientPerfPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	// 更新 Hub 中的 Client 效能快照
	g.Hub.UpdateClientPerf(clientID, payload.FPS, payload.MemoryMB, payload.DrawCalls, payload.Quality)

	// 高延遲警告（DAY-045）：Client 端 ping > 200ms 輸出警告 log
	// 這讓運維人員能識別網路品質差的玩家
	if payload.PingMs > 200 {
		log.Printf("[PerfAlert] High latency player %s: ping=%dms fps=%.1f quality=%s",
			clientID, payload.PingMs, payload.FPS, payload.Quality)
	}

	// 低 FPS 警告：Client 端 FPS < 20 輸出警告 log
	if payload.FPS > 0 && payload.FPS < 20 {
		log.Printf("[PerfAlert] Low FPS player %s: fps=%.1f memory=%.1fMB drawcalls=%d quality=%s",
			clientID, payload.FPS, payload.MemoryMB, payload.DrawCalls, payload.Quality)
	}
}
