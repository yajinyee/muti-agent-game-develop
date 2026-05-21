// giant_prize_fish_handler.go — 夢幻巨型獎勵魚系統 handler（DAY-147）
// 業界依據：jiligames.com 2026「The dreamy Giant Prize Fish lets you easily win great prizes,
// with the chance for 5x multipliers」
// 擊破 T111 後觸發「夢幻獎勵模式」：觸發玩家在 10 秒內所有擊破獎勵 ×5
// 設計理念：低 HP（容易擊破）+ 中等倍率（40-60x）+ 觸發後 10 秒 5x 加成
// 這是「容易觸發的短期爆發」機制，讓玩家感受到「夢幻大獎」的爽感
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
	// GiantPrizeFishMultBonus 夢幻獎勵模式倍率加成
	GiantPrizeFishMultBonus = 5.0
	// GiantPrizeFishDuration 夢幻獎勵模式持續時間（秒）
	GiantPrizeFishDuration = 10
	// GiantPrizeFishCooldown 每個玩家的冷卻時間（秒）
	GiantPrizeFishCooldown = 60
	// GiantPrizeFishAnnounceThreshold 全服公告門檻（夢幻模式期間擊破數）
	GiantPrizeFishAnnounceThreshold = 5
)

// giantPrizeFishSession 夢幻獎勵模式 session（每個玩家獨立）
type giantPrizeFishSession struct {
	playerID    string
	startAt     time.Time
	endAt       time.Time
	totalReward int
	killCount   int
}

// giantPrizeFishManager 夢幻獎勵魚管理器（輕量，不需要獨立套件）
type giantPrizeFishManager struct {
	mu       sync.RWMutex
	sessions map[string]*giantPrizeFishSession // playerID → session
	cooldown map[string]time.Time              // playerID → 冷卻結束時間
}

// newGiantPrizeFishManager 建立管理器
func newGiantPrizeFishManager() *giantPrizeFishManager {
	return &giantPrizeFishManager{
		sessions: make(map[string]*giantPrizeFishSession),
		cooldown: make(map[string]time.Time),
	}
}

// isGiantPrizeFish 判斷是否為夢幻巨型獎勵魚（T111）
func isGiantPrizeFish(defID string) bool {
	return defID == "T111"
}

// tryGiantPrizeFish 擊破 T111 後觸發夢幻獎勵模式（DAY-147）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryGiantPrizeFish(p *player.Player, triggerID string, triggerX, triggerY float64) {
	if g.GiantPrizeFish == nil {
		return
	}

	mgr := g.GiantPrizeFish
	mgr.mu.Lock()

	// 冷卻檢查
	if cd, ok := mgr.cooldown[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		log.Printf("[GiantPrizeFish] player=%s still in cooldown", p.ID)
		return
	}

	// 建立 session
	now := time.Now()
	session := &giantPrizeFishSession{
		playerID: p.ID,
		startAt:  now,
		endAt:    now.Add(GiantPrizeFishDuration * time.Second),
	}
	mgr.sessions[p.ID] = session
	mgr.cooldown[p.ID] = now.Add(GiantPrizeFishCooldown * time.Second)
	mgr.mu.Unlock()

	log.Printf("[GiantPrizeFish] player=%s activated dreamy mode for %ds (×%.0f)",
		p.ID, GiantPrizeFishDuration, GiantPrizeFishMultBonus)

	// 全服廣播夢幻模式啟動
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGiantPrizeFish,
		Payload: ws.GiantPrizeFishPayload{
			TriggerID:  triggerID,
			TriggerX:   triggerX,
			TriggerY:   triggerY,
			Phase:      "activate",
			MultBonus:  GiantPrizeFishMultBonus,
			Duration:   GiantPrizeFishDuration,
			KillerID:   p.ID,
			KillerName: p.DisplayName,
		},
	})

	// 等待夢幻模式結束
	time.Sleep(GiantPrizeFishDuration * time.Second)

	// 讀取結果
	mgr.mu.Lock()
	finalSession := mgr.sessions[p.ID]
	delete(mgr.sessions, p.ID)
	mgr.mu.Unlock()

	totalReward := 0
	killCount := 0
	if finalSession != nil {
		totalReward = finalSession.totalReward
		killCount = finalSession.killCount
	}

	// 全服廣播夢幻模式結束
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGiantPrizeFish,
		Payload: ws.GiantPrizeFishPayload{
			TriggerID:   triggerID,
			Phase:       "end",
			MultBonus:   GiantPrizeFishMultBonus,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			TotalReward: totalReward,
			KillCount:   killCount,
		},
	})

	// 全服公告：夢幻模式期間擊破 ≥5 個目標
	if killCount >= GiantPrizeFishAnnounceThreshold {
		g.announceGiantPrizeFish(p.DisplayName, killCount, totalReward)
	}

	log.Printf("[GiantPrizeFish] player=%s dreamy mode ended: kills=%d total_reward=%d",
		p.ID, killCount, totalReward)
}

// getGiantPrizeFishMult 取得夢幻獎勵模式倍率（供 handleKill 使用）
// 如果玩家正在夢幻模式中，回傳 5.0；否則回傳 1.0
func (g *Game) getGiantPrizeFishMult(playerID string) float64 {
	if g.GiantPrizeFish == nil {
		return 1.0
	}
	mgr := g.GiantPrizeFish
	mgr.mu.RLock()
	defer mgr.mu.RUnlock()

	session, ok := mgr.sessions[playerID]
	if !ok {
		return 1.0
	}
	if time.Now().After(session.endAt) {
		return 1.0
	}
	return GiantPrizeFishMultBonus
}

// recordGiantPrizeFishKill 記錄夢幻模式期間的擊破（供 handleKill 使用）
func (g *Game) recordGiantPrizeFishKill(playerID string, reward int) {
	if g.GiantPrizeFish == nil {
		return
	}
	mgr := g.GiantPrizeFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	session, ok := mgr.sessions[playerID]
	if !ok {
		return
	}
	if time.Now().After(session.endAt) {
		return
	}
	session.totalReward += reward
	session.killCount++
}

// announceGiantPrizeFish 全服公告夢幻獎勵模式（DAY-147）
func (g *Game) announceGiantPrizeFish(playerName string, killCount int, reward int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "giant_prize_fish",
			"message":    fmt.Sprintf("✨ %s 的夢幻獎勵魚在 10 秒內擊破 %d 個目標！獲得 %d 金幣！", playerName, killCount, reward),
			"color":      "#FF69B4",
			"duration":   4.5,
			"priority":   2,
		},
	})
}
