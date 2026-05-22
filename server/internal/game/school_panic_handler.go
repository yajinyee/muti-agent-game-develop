// school_panic_handler.go — 魚群驚嚇連帶系統 handler（DAY-191）
// 業界靈感：Ocean King 3 Plus「School of Fish — when one fish in a school is caught,
// the others scatter in panic, but a lucky shot can trigger a chain reaction that catches the entire school」
// 設計：擊破 T149 魚群領袖後觸發「魚群驚嚇」：
//   - 場上所有基礎目標（T001-T006）HP 降低 50%（更容易擊破）
//   - 持續 8 秒，讓玩家在「魚群驚嚇」中快速收割
//   - 全服廣播：讓所有玩家看到「魚群驚嚇中，快打！」
// 設計差異：與漩渦魚（直接擊破所有基礎目標）不同，魚群驚嚇是「降低 HP」而非直接擊破，
// 讓玩家仍需要主動射擊，製造「緊張但有利」的感覺；
// 與冰凍炸彈魚（凍結特殊目標）不同，魚群驚嚇影響基礎目標，讓普通遊戲也有爽感
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	SchoolPanicDefID       = "T149"
	SchoolPanicDuration    = 8.0  // 驚嚇持續秒數
	SchoolPanicHPReduction = 0.50 // HP 降低 50%
	SchoolPanicCooldown    = 25   // 全服冷卻秒數
)

// schoolPanicManager 魚群驚嚇管理器
type schoolPanicManager struct {
	mu          sync.Mutex
	isActive    bool
	panicEnd    time.Time
	cooldownEnd time.Time
	instanceID  string // 觸發的魚群領袖 InstanceID
}

func newSchoolPanicManager() *schoolPanicManager {
	return &schoolPanicManager{}
}

// isSchoolLeader 判斷是否為魚群領袖
func isSchoolLeader(defID string) bool {
	return defID == SchoolPanicDefID
}

// isSchoolPanicActive 查詢魚群驚嚇是否活躍（供 combat 使用）
func (g *Game) isSchoolPanicActive() bool {
	if g.SchoolPanic == nil {
		return false
	}
	g.SchoolPanic.mu.Lock()
	defer g.SchoolPanic.mu.Unlock()
	if !g.SchoolPanic.isActive {
		return false
	}
	if time.Now().After(g.SchoolPanic.panicEnd) {
		g.SchoolPanic.isActive = false
		return false
	}
	return true
}

// isSchoolBasicTarget 判斷是否為基礎目標（T001-T006），供魚群驚嚇使用
func isSchoolBasicTarget(defID string) bool {
	switch defID {
	case "T001", "T002", "T003", "T004", "T005", "T006":
		return true
	}
	return false
}

// trySchoolPanic 擊破魚群領袖後觸發魚群驚嚇（由 handleKill 呼叫）
func (g *Game) trySchoolPanic(p *player.Player, killedInstanceID string) {
	if g.SchoolPanic == nil {
		return
	}

	g.SchoolPanic.mu.Lock()
	// 檢查冷卻
	if time.Now().Before(g.SchoolPanic.cooldownEnd) {
		g.SchoolPanic.mu.Unlock()
		log.Printf("[SchoolPanic] on cooldown, skip")
		return
	}
	// 防止重複觸發
	if g.SchoolPanic.isActive {
		g.SchoolPanic.mu.Unlock()
		return
	}

	now := time.Now()
	g.SchoolPanic.isActive = true
	g.SchoolPanic.panicEnd = now.Add(time.Duration(SchoolPanicDuration * float64(time.Second)))
	g.SchoolPanic.cooldownEnd = now.Add(time.Duration(SchoolPanicCooldown) * time.Second)
	g.SchoolPanic.instanceID = killedInstanceID
	g.SchoolPanic.mu.Unlock()

	// 收集場上所有基礎目標，降低 HP
	g.mu.Lock()
	panicTargets := make([]struct {
		instanceID string
		defID      string
		newHP      int
		x, y       float64
	}, 0, 16)

	for _, t := range g.Targets {
		if !isSchoolBasicTarget(t.DefID) || t.HP <= 0 {
			continue
		}
		// HP 降低 50%（最少 1）
		newHP := t.HP / 2
		if newHP < 1 {
			newHP = 1
		}
		t.HP = newHP
		panicTargets = append(panicTargets, struct {
			instanceID string
			defID      string
			newHP      int
			x, y       float64
		}{t.InstanceID, t.DefID, newHP, t.X, t.Y})
	}
	g.mu.Unlock()

	log.Printf("[SchoolPanic] player=%s triggered panic: %d basic targets HP halved",
		p.ID, len(panicTargets))

	// 廣播驚嚇開始（全服）
	panicIDs := make([]string, 0, len(panicTargets))
	for _, pt := range panicTargets {
		panicIDs = append(panicIDs, pt.instanceID)
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgSchoolPanic,
		Payload: ws.SchoolPanicPayload{
			Phase:        "panic_start",
			TriggerID:    killedInstanceID,
			KillerID:     p.ID,
			KillerName:   p.DisplayName,
			PanicTargets: panicIDs,
			TargetCount:  len(panicTargets),
			Duration:     SchoolPanicDuration,
			Message:      fmt.Sprintf("🐟 %s 觸發魚群驚嚇！%d 條基礎魚 HP 減半！快打！", p.DisplayName, len(panicTargets)),
		},
	})

	// 全服公告
	if len(panicTargets) >= 3 {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgAnnouncement,
			Payload: map[string]interface{}{
				"event_type": "school_panic",
				"message":    fmt.Sprintf("🐟 %s 觸發魚群驚嚇！%d 條魚 HP 減半！快搶！", p.DisplayName, len(panicTargets)),
				"color":      "#FF8C00",
				"duration":   4.0,
				"priority":   2,
			},
		})
	}

	// 8 秒後廣播驚嚇結束
	go func() {
		time.Sleep(time.Duration(SchoolPanicDuration) * time.Second)

		g.SchoolPanic.mu.Lock()
		g.SchoolPanic.isActive = false
		g.SchoolPanic.mu.Unlock()

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgSchoolPanic,
			Payload: ws.SchoolPanicPayload{
				Phase:   "panic_end",
				Message: "🐟 魚群驚嚇結束",
			},
		})
		log.Printf("[SchoolPanic] panic ended")
	}()
}
