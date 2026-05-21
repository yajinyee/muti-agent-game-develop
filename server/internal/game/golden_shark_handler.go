// golden_shark_handler.go — 黃金鯊魚狂暴模式 handler（DAY-161）
// 業界依據：King of Ocean 2026「sharks climb into x50-x300 zone」
// + 捕魚機業界「rage/berserk mode」機制 — 擊破黃金鯊魚後觸發「全服狂暴模式」
// 全場所有目標物獎勵倍率 ×1.5，持續 12 秒，全服廣播
// 設計：全服共享（不是個人），任何玩家擊破都讓全服受益，製造「全場爆發」的社交爽感
// 與幸運星魚（個人 ×2）不同：黃金鯊魚是全服 ×1.5，社交性更強
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// goldenSharkManager 黃金鯊魚狂暴模式管理器（全服共享）
type goldenSharkManager struct {
	mu          sync.Mutex
	isActive    bool
	startAt     time.Time
	duration    float64
	killerID    string
	killerName  string
	cooldownEnd time.Time
}

func newGoldenSharkManager() *goldenSharkManager {
	return &goldenSharkManager{}
}

// isGoldenShark 判斷是否為黃金鯊魚
func isGoldenShark(defID string) bool {
	return defID == "T121"
}

// getGoldenSharkMult 取得黃金鯊魚狂暴倍率加成（供 handleKill 使用）
// 若全服狂暴模式活躍，回傳 1.5，否則回傳 1.0
func (g *Game) getGoldenSharkMult() float64 {
	if g.GoldenShark == nil {
		return 1.0
	}
	g.GoldenShark.mu.Lock()
	defer g.GoldenShark.mu.Unlock()

	if !g.GoldenShark.isActive {
		return 1.0
	}
	elapsed := time.Since(g.GoldenShark.startAt).Seconds()
	if elapsed >= g.GoldenShark.duration {
		g.GoldenShark.isActive = false
		return 1.0
	}
	return 1.5
}

// tryGoldenSharkBerserk 嘗試觸發黃金鯊魚狂暴模式（擊破 T121 後呼叫）
func (g *Game) tryGoldenSharkBerserk(p *player.Player, instanceID string, x, y float64) {
	if g.GoldenShark == nil {
		return
	}

	const duration = 12.0 // 12 秒狂暴模式
	const cooldown = 90.0 // 90 秒冷卻（全服共享，冷卻更長）

	g.GoldenShark.mu.Lock()
	if g.GoldenShark.isActive {
		g.GoldenShark.mu.Unlock()
		return // 已在狂暴模式中
	}
	if time.Now().Before(g.GoldenShark.cooldownEnd) {
		g.GoldenShark.mu.Unlock()
		return // 冷卻中
	}

	// 啟動狂暴模式
	g.GoldenShark.isActive = true
	g.GoldenShark.startAt = time.Now()
	g.GoldenShark.duration = duration
	g.GoldenShark.killerID = p.ID
	g.GoldenShark.killerName = p.DisplayName
	g.GoldenShark.mu.Unlock()

	log.Printf("[GoldenShark] player=%s triggered berserk mode x1.5 for %.0fs (global)", p.ID, duration)

	// 廣播狂暴模式開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGoldenSharkBerserk,
		Payload: ws.GoldenSharkBerserkPayload{
			TriggerID:    instanceID,
			TriggerX:     x,
			TriggerY:     y,
			KillerID:     p.ID,
			KillerName:   p.DisplayName,
			Phase:        "berserk_start",
			DurationSecs: duration,
			MultBonus:    1.5,
		},
	})

	// 全服公告
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "golden_shark_berserk",
			"message":    fmt.Sprintf("🦈 %s 擊破黃金鯊魚！全場狂暴模式！所有獎勵 ×1.5 持續 %.0f 秒！", p.DisplayName, duration),
			"color":      "#FF6600",
			"duration":   5.0,
			"priority":   4,
		},
	})

	// 等待狂暴模式結束
	go func() {
		time.Sleep(time.Duration(duration * float64(time.Second)))

		g.GoldenShark.mu.Lock()
		g.GoldenShark.isActive = false
		g.GoldenShark.cooldownEnd = time.Now().Add(time.Duration(cooldown * float64(time.Second)))
		g.GoldenShark.mu.Unlock()

		// 廣播狂暴模式結束
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgGoldenSharkBerserk,
			Payload: ws.GoldenSharkBerserkPayload{
				KillerID:   p.ID,
				KillerName: p.DisplayName,
				Phase:      "berserk_end",
			},
		})

		log.Printf("[GoldenShark] berserk mode ended, cooldown %.0fs", cooldown)
	}()
}
