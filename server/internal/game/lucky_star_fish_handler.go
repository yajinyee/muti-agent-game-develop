// lucky_star_fish_handler.go — 幸運星魚全場倍率翻倍系統 handler（DAY-160）
// 業界依據：捕魚機業界標準「倍率爆發」機制
// 擊破幸運星魚後觸發「全場倍率翻倍」10 秒，所有目標物的獎勵倍率 ×2
// 是「爆發型特殊目標」，讓玩家在 10 秒內所有擊破獎勵翻倍，製造「大豐收」的爽感
// 設計：每個玩家獨立 session，觸發者享受 10 秒倍率翻倍，60 秒冷卻
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// luckyStarFishManager 幸運星魚倍率翻倍管理器
type luckyStarFishManager struct {
	mu       sync.Mutex
	sessions map[string]*luckyStarSession // playerID -> session
}

type luckyStarSession struct {
	PlayerID    string
	StartAt     time.Time
	Duration    float64
	CooldownEnd time.Time
}

func newLuckyStarFishManager() *luckyStarFishManager {
	return &luckyStarFishManager{
		sessions: make(map[string]*luckyStarSession),
	}
}

// isLuckyStarFish 判斷是否為幸運星魚
func isLuckyStarFish(defID string) bool {
	return defID == "T120"
}

// getLuckyStarMult 取得幸運星魚倍率加成（供 handleKill 使用）
// 若玩家有活躍 session，回傳 2.0（翻倍），否則回傳 1.0
func (g *Game) getLuckyStarMult(playerID string) float64 {
	if g.LuckyStarFish == nil {
		return 1.0
	}
	g.LuckyStarFish.mu.Lock()
	defer g.LuckyStarFish.mu.Unlock()

	sess, ok := g.LuckyStarFish.sessions[playerID]
	if !ok {
		return 1.0
	}
	elapsed := time.Since(sess.StartAt).Seconds()
	if elapsed >= sess.Duration {
		return 1.0
	}
	return 2.0
}

// tryLuckyStarFish 嘗試觸發幸運星魚倍率翻倍（擊破 T120 後呼叫）
func (g *Game) tryLuckyStarFish(p *player.Player, instanceID string, x, y float64) {
	if g.LuckyStarFish == nil {
		return
	}

	const duration = 10.0 // 10 秒倍率翻倍
	const cooldown = 60.0 // 60 秒冷卻

	g.LuckyStarFish.mu.Lock()
	sess, exists := g.LuckyStarFish.sessions[p.ID]
	if exists && time.Now().Before(sess.CooldownEnd) {
		g.LuckyStarFish.mu.Unlock()
		return // 冷卻中
	}

	// 建立新 session
	newSess := &luckyStarSession{
		PlayerID: p.ID,
		StartAt:  time.Now(),
		Duration: duration,
	}
	g.LuckyStarFish.sessions[p.ID] = newSess
	g.LuckyStarFish.mu.Unlock()

	log.Printf("[LuckyStarFish] player=%s activated x2 mult for %.0fs", p.ID, duration)

	// 廣播倍率翻倍開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyStarFish,
		Payload: ws.LuckyStarFishPayload{
			TriggerID:    instanceID,
			TriggerX:     x,
			TriggerY:     y,
			KillerID:     p.ID,
			KillerName:   p.DisplayName,
			Phase:        "lucky_start",
			DurationSecs: duration,
			MultBonus:    2.0,
		},
	})

	// 全服公告
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "lucky_star_fish",
			"message":    fmt.Sprintf("⭐ %s 擊破幸運星魚！全場獎勵翻倍 %.0f 秒！", p.DisplayName, duration),
			"color":      "#FFD700",
			"duration":   4.0,
			"priority":   3,
		},
	})

	// 等待倍率翻倍結束
	go func() {
		time.Sleep(time.Duration(duration * float64(time.Second)))

		g.LuckyStarFish.mu.Lock()
		if s, ok := g.LuckyStarFish.sessions[p.ID]; ok {
			s.CooldownEnd = time.Now().Add(time.Duration(cooldown * float64(time.Second)))
		}
		g.LuckyStarFish.mu.Unlock()

		// 廣播倍率翻倍結束
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyStarFish,
			Payload: ws.LuckyStarFishPayload{
				KillerID:   p.ID,
				KillerName: p.DisplayName,
				Phase:      "lucky_end",
			},
		})

		log.Printf("[LuckyStarFish] player=%s x2 mult ended, cooldown %.0fs", p.ID, cooldown)
	}()
}
