// rainbow_lucky_fish_handler.go — 彩虹幸運魚系統 handler（DAY-173）
// 業界依據：Fisch Roblox 2026「Rainbow Leviathan — rare rainbow fish that triggers a luck boost event」
// + Fish It 2026「Rainbow Throw — triggered by Prismatic enchant, increases luck for rare fish」
// + Ocean King 2026「Rainbow Fish — when caught, all players receive a luck boost for 10 seconds」
// 擊破 T131 後觸發「彩虹幸運時間」（Rainbow Lucky Time）：
//   - 持續 10 秒，全服所有玩家的擊破機率提升 20%（BASE_RTP × 1.2）
//   - 全服廣播彩虹光效，讓所有玩家感受到「幸運時間到了！」
// 設計差異：與幸運星魚（個人 ×2 倍率，10秒）不同，彩虹幸運魚是**全服共享的擊破機率加成**，
// 讓所有玩家在 10 秒內更容易擊破目標，製造「全服一起爽」的社交感；
// 與黃金鯊魚（全服 ×1.5 倍率）不同，彩虹幸運魚是**機率加成**（更容易擊破），
// 讓玩家感受到「這 10 秒打什麼都中」的爽感
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/ws"
)

const (
	// RainbowLuckyDurationSec 彩虹幸運時間持續時間（秒）
	RainbowLuckyDurationSec = 10
	// RainbowLuckyKillChanceBoost 擊破機率加成（+20%）
	RainbowLuckyKillChanceBoost = 0.20
	// RainbowLuckyCooldownSec 全服冷卻時間（秒）
	RainbowLuckyCooldownSec = 60
)

// rainbowLuckyManager 彩虹幸運魚管理器（全服共享）
type rainbowLuckyManager struct {
	mu        sync.Mutex
	isActive  bool
	expiresAt time.Time
	cooldown  time.Time
}

// newRainbowLuckyManager 建立彩虹幸運魚管理器
func newRainbowLuckyManager() *rainbowLuckyManager {
	return &rainbowLuckyManager{}
}

// isRainbowLuckyFish 判斷是否為彩虹幸運魚（T131）
func isRainbowLuckyFish(defID string) bool {
	return defID == "T131"
}

// IsRainbowLuckyActive 查詢彩虹幸運時間是否激活（供 combat 系統使用）
func (g *Game) IsRainbowLuckyActive() bool {
	if g.RainbowLucky == nil {
		return false
	}
	g.RainbowLucky.mu.Lock()
	defer g.RainbowLucky.mu.Unlock()
	if !g.RainbowLucky.isActive {
		return false
	}
	if time.Now().After(g.RainbowLucky.expiresAt) {
		g.RainbowLucky.isActive = false
		return false
	}
	return true
}

// GetRainbowLuckyBoost 取得彩虹幸運時間的擊破機率加成（供 combat 系統使用）
// 回傳 0.20（有加成）或 0.0（無加成）
func (g *Game) GetRainbowLuckyBoost() float64 {
	if g.IsRainbowLuckyActive() {
		return RainbowLuckyKillChanceBoost
	}
	return 0.0
}

// tryRainbowLuckyFish 擊破 T131 後觸發彩虹幸運時間（DAY-173）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryRainbowLuckyFish(playerName string, triggerX, triggerY float64) {
	g.RainbowLucky.mu.Lock()

	// 全服冷卻檢查
	if time.Now().Before(g.RainbowLucky.cooldown) {
		g.RainbowLucky.mu.Unlock()
		return
	}

	// 激活彩虹幸運時間
	g.RainbowLucky.isActive = true
	g.RainbowLucky.expiresAt = time.Now().Add(time.Duration(RainbowLuckyDurationSec) * time.Second)
	g.RainbowLucky.cooldown = time.Now().Add(time.Duration(RainbowLuckyCooldownSec) * time.Second)
	g.RainbowLucky.mu.Unlock()

	log.Printf("[RainbowLucky] activated by %s, duration=%ds, boost=+%.0f%%",
		playerName, RainbowLuckyDurationSec, RainbowLuckyKillChanceBoost*100)

	// 全服廣播：彩虹幸運時間開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRainbowLuckyFish,
		Payload: ws.RainbowLuckyFishPayload{
			Phase:       "lucky_start",
			PlayerName:  playerName,
			DurationSec: RainbowLuckyDurationSec,
			KillBoost:   RainbowLuckyKillChanceBoost,
			TriggerX:    triggerX,
			TriggerY:    triggerY,
		},
	})

	// 全服公告
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "rainbow_lucky_fish",
			"message":    fmt.Sprintf("🌈 %s 擊破彩虹幸運魚！全服幸運時間 %d 秒！擊破機率 +%.0f%%！", playerName, RainbowLuckyDurationSec, RainbowLuckyKillChanceBoost*100),
			"color":      "#FF69B4",
			"duration":   4.0,
			"priority":   3,
		},
	})

	// 10 秒後廣播結束
	time.Sleep(time.Duration(RainbowLuckyDurationSec) * time.Second)

	g.RainbowLucky.mu.Lock()
	g.RainbowLucky.isActive = false
	g.RainbowLucky.mu.Unlock()

	// 全服廣播：彩虹幸運時間結束
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRainbowLuckyFish,
		Payload: ws.RainbowLuckyFishPayload{
			Phase: "lucky_end",
		},
	})

	log.Printf("[RainbowLucky] ended")
}
