// lucky_clover_fish_handler.go — 幸運草魚系統 handler（DAY-179）
// 業界依據：Ocean King 3 Plus「Lucky Shamrock Leprechaun Boss」
// + Fisch Roblox 2026「Lucky Gold Pool — rainbow event triggers lucky fish spawns」
// 擊破 T137 後觸發「幸運草爆發」：
//   1. 場上所有目標物獎勵 +50% 持續 10 秒（全服共享）
//   2. 隨機為 1-3 個在線玩家發放「幸運草金幣」（betLevel × 10-30x）
// 設計差異：與黃金鯊魚（全服 ×1.5 倍率乘法）不同，幸運草是「+50% 加成」（加法），
// 且有「隨機發放金幣給玩家」的社交機制，讓被選中的玩家感到「幸運」
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// LuckyCloverCooldownSec 全服冷卻時間（秒）
	LuckyCloverCooldownSec = 35
	// LuckyCloverBoostDurationSec 幸運加成持續時間（秒）
	LuckyCloverBoostDurationSec = 10
	// LuckyCloverBoostPercent 幸運加成百分比（+50%）
	LuckyCloverBoostPercent = 0.50
	// LuckyCloverMinGiftPlayers 最少發放金幣玩家數
	LuckyCloverMinGiftPlayers = 1
	// LuckyCloverMaxGiftPlayers 最多發放金幣玩家數
	LuckyCloverMaxGiftPlayers = 3
	// LuckyCloverGiftMinMult 幸運草金幣最小倍率
	LuckyCloverGiftMinMult = 10
	// LuckyCloverGiftMaxMult 幸運草金幣最大倍率
	LuckyCloverGiftMaxMult = 30
)

// luckyCloverManager 幸運草魚管理器（全服共享）
type luckyCloverManager struct {
	mu           sync.Mutex
	isActive     bool
	cooldownAt   time.Time
	boostExpiresAt time.Time
}

// newLuckyCloverManager 建立幸運草魚管理器
func newLuckyCloverManager() *luckyCloverManager {
	return &luckyCloverManager{}
}

// isLuckyCloverFish 判斷是否為幸運草魚（T137）
func isLuckyCloverFish(defID string) bool {
	return defID == "T137"
}

// canTrigger 檢查是否可以觸發
func (m *luckyCloverManager) canTrigger() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isActive {
		return false
	}
	return time.Now().After(m.cooldownAt)
}

// activate 激活幸運草爆發
func (m *luckyCloverManager) activate() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isActive = true
	m.boostExpiresAt = time.Now().Add(time.Duration(LuckyCloverBoostDurationSec) * time.Second)
}

// deactivate 停用幸運草爆發，設定冷卻
func (m *luckyCloverManager) deactivate() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isActive = false
	m.cooldownAt = time.Now().Add(time.Duration(LuckyCloverCooldownSec) * time.Second)
}

// getLuckyCloverBoost 取得幸運草加成（供 handleKill 使用）
// 回傳額外加成比例（0.0 = 無加成，0.50 = +50%）
func (m *luckyCloverManager) getLuckyCloverBoost() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.isActive {
		return 0.0
	}
	if time.Now().After(m.boostExpiresAt) {
		m.isActive = false
		return 0.0
	}
	return LuckyCloverBoostPercent
}

// tryLuckyCloverFish 擊破 T137 後觸發幸運草爆發（DAY-179）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryLuckyCloverFish(p *player.Player) {
	if !g.LuckyClover.canTrigger() {
		return
	}
	g.LuckyClover.activate()
	defer g.LuckyClover.deactivate()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	log.Printf("[LuckyClover] player=%s triggered, boost +%.0f%% for %ds",
		p.ID, LuckyCloverBoostPercent*100, LuckyCloverBoostDurationSec)

	// 廣播幸運草爆發開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyCloverFish,
		Payload: ws.LuckyCloverFishPayload{
			Phase:            "clover_start",
			PlayerID:         p.ID,
			PlayerName:       p.DisplayName,
			BoostPercent:     LuckyCloverBoostPercent,
			BoostDurationSec: LuckyCloverBoostDurationSec,
		},
	})

	// 隨機選取 1-3 個在線玩家發放幸運草金幣
	g.mu.RLock()
	var onlinePlayers []*player.Player
	for _, pp := range g.Players {
		onlinePlayers = append(onlinePlayers, pp)
	}
	g.mu.RUnlock()

	giftCount := LuckyCloverMinGiftPlayers + rng.Intn(LuckyCloverMaxGiftPlayers-LuckyCloverMinGiftPlayers+1)
	if giftCount > len(onlinePlayers) {
		giftCount = len(onlinePlayers)
	}

	// 隨機打亂玩家列表
	rng.Shuffle(len(onlinePlayers), func(i, j int) {
		onlinePlayers[i], onlinePlayers[j] = onlinePlayers[j], onlinePlayers[i]
	})

	// 發放幸運草金幣
	for i := 0; i < giftCount; i++ {
		giftPlayer := onlinePlayers[i]
		giftMult := LuckyCloverGiftMinMult + rng.Intn(LuckyCloverGiftMaxMult-LuckyCloverGiftMinMult+1)
		giftAmount := giftPlayer.BetLevel * giftMult

		giftPlayer.AddReward(giftAmount)

		// 廣播幸運草金幣（個人）
		g.Hub.Send(giftPlayer.ID, &ws.Message{
			Type: ws.MsgLuckyCloverFish,
			Payload: ws.LuckyCloverFishPayload{
				Phase:      "clover_gift",
				PlayerID:   giftPlayer.ID,
				PlayerName: giftPlayer.DisplayName,
				GiftAmount: giftAmount,
				GiftMult:   giftMult,
				NewBalance: giftPlayer.Coins,
			},
		})

		log.Printf("[LuckyClover] gift player=%s amount=%d (betLevel=%d × %dx)",
			giftPlayer.ID, giftAmount, giftPlayer.BetLevel, giftMult)
	}

	// 全服公告
	g.announceLuckyCloverFish(p.DisplayName, giftCount)

	// 等待幸運加成結束
	time.Sleep(time.Duration(LuckyCloverBoostDurationSec) * time.Second)

	// 廣播幸運草爆發結束（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyCloverFish,
		Payload: ws.LuckyCloverFishPayload{
			Phase: "clover_end",
		},
	})

	log.Printf("[LuckyClover] player=%s boost ended", p.ID)
}

// getLuckyCloverBoost 取得幸運草加成（供 handleKill 使用）
func (g *Game) getLuckyCloverBoost() float64 {
	return g.LuckyClover.getLuckyCloverBoost()
}

// announceLuckyCloverFish 全服公告幸運草魚（DAY-179）
func (g *Game) announceLuckyCloverFish(playerName string, giftCount int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "lucky_clover_fish",
			"message":    fmt.Sprintf("🍀 %s 觸發幸運草爆發！全服 +50%% 加成！%d 位玩家獲得幸運草金幣！", playerName, giftCount),
			"color":      "#00FF7F", // 春綠色（幸運草感）
			"duration":   5.0,
			"priority":   3,
		},
	})
}
