// golden_turtle_handler.go — 黃金海龜時間停止系統 handler（DAY-159）
// 業界依據：Ocean King 系列「Time Stop」機制
// 擊破黃金海龜後觸發「全場時間停止」8 秒，所有目標物暫停移動
// 玩家可以在 8 秒內輕鬆瞄準並大量擊破其他目標，是「輔助型特殊目標」
// 設計：不直接給高獎勵，但讓玩家在 8 秒內大量擊破其他目標，間接提升收益
// 全服廣播讓所有玩家都能享受時間停止的爽感
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// goldenTurtleManager 黃金海龜時間停止管理器
type goldenTurtleManager struct {
	mu          sync.Mutex
	isActive    bool
	activatedBy string // 觸發玩家 ID
	activatedAt time.Time
	duration    float64 // 停止時間（秒）
	cooldownEnd time.Time
}

func newGoldenTurtleManager() *goldenTurtleManager {
	return &goldenTurtleManager{
		duration: 8.0, // 8 秒時間停止
	}
}

// isGoldenTurtle 判斷是否為黃金海龜
func isGoldenTurtle(defID string) bool {
	return defID == "T119"
}

// tryGoldenTurtleTimeStop 嘗試觸發時間停止（擊破 T119 後呼叫）
func (g *Game) tryGoldenTurtleTimeStop(p *player.Player, instanceID string, x, y float64) {
	if g.GoldenTurtle == nil {
		return
	}

	g.GoldenTurtle.mu.Lock()
	// 冷卻中不觸發
	if time.Now().Before(g.GoldenTurtle.cooldownEnd) {
		g.GoldenTurtle.mu.Unlock()
		return
	}
	// 已有活躍 session 不觸發
	if g.GoldenTurtle.isActive {
		g.GoldenTurtle.mu.Unlock()
		return
	}

	g.GoldenTurtle.isActive = true
	g.GoldenTurtle.activatedBy = p.ID
	g.GoldenTurtle.activatedAt = time.Now()
	g.GoldenTurtle.mu.Unlock()

	log.Printf("[GoldenTurtle] time stop activated by player=%s, duration=%.0fs", p.ID, g.GoldenTurtle.duration)

	// 廣播時間停止開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGoldenTurtleTimeStop,
		Payload: ws.GoldenTurtleTimeStopPayload{
			TriggerID:    instanceID,
			TriggerX:     x,
			TriggerY:     y,
			KillerID:     p.ID,
			KillerName:   p.DisplayName,
			Phase:        "time_stop_start",
			DurationSecs: g.GoldenTurtle.duration,
		},
	})

	// 全服公告
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "golden_turtle",
			"message":    fmt.Sprintf("🐢 %s 擊破黃金海龜！全場時間停止 %.0f 秒！", p.DisplayName, g.GoldenTurtle.duration),
			"color":      "#FFD700",
			"duration":   4.0,
			"priority":   3,
		},
	})

	// 等待時間停止結束
	go func() {
		time.Sleep(time.Duration(g.GoldenTurtle.duration * float64(time.Second)))

		g.GoldenTurtle.mu.Lock()
		g.GoldenTurtle.isActive = false
		g.GoldenTurtle.cooldownEnd = time.Now().Add(60 * time.Second) // 60 秒冷卻
		g.GoldenTurtle.mu.Unlock()

		// 廣播時間停止結束
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgGoldenTurtleTimeStop,
			Payload: ws.GoldenTurtleTimeStopPayload{
				KillerID:   p.ID,
				KillerName: p.DisplayName,
				Phase:      "time_stop_end",
			},
		})

		log.Printf("[GoldenTurtle] time stop ended, cooldown 60s")
	}()
}

// IsTimeStopActive 查詢時間停止是否活躍（供 spawnTarget 和 updateNormalPlay 使用）
func (g *Game) IsTimeStopActive() bool {
	if g.GoldenTurtle == nil {
		return false
	}
	g.GoldenTurtle.mu.Lock()
	defer g.GoldenTurtle.mu.Unlock()
	return g.GoldenTurtle.isActive
}
