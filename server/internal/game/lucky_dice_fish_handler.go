// lucky_dice_fish_handler.go — 幸運骰子魚系統 handler（DAY-175）
// 業界依據：Ocean King 3 Plus「Fast Bomb — randomly triggered bonus that instantly destroys
// multiple fish」+ 捕魚機業界「Dice Roll bonus — roll dice to determine reward multiplier」
// + Fishing Carnival 2026「Dice Fish — catching triggers a dice roll, sum determines payout」
// 擊破 T133 後觸發「幸運骰子」：擲 2 顆骰子（1-6），點數之和決定獎勵倍率
//   - 點數 2（蛇眼）：20x betLevel（特殊彩蛋）
//   - 點數 7（最常見）：7x betLevel
//   - 點數 12（大六）：50x betLevel（大獎）
//   - 其他點數：點數 × betLevel
// 設計差異：與輪盤（多格選擇）不同，骰子是「兩顆骰子點數之和」，
// 機率分布符合真實骰子（7最常見，2和12最稀有），讓玩家有「期待骰子停下來」的緊張感；
// 與彩蛋（隨機獎勵類型）不同，骰子是「純數值獎勵」，更直觀
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
	// LuckyDiceCooldownSec 個人冷卻時間（秒）
	LuckyDiceCooldownSec = 25
	// LuckyDiceRollDurationMs 骰子滾動時間（ms，Client 動畫用）
	LuckyDiceRollDurationMs = 2000
	// LuckyDiceAnnounceThreshold 全服公告門檻（點數之和）
	LuckyDiceAnnounceThreshold = 10
)

// luckyDiceManager 幸運骰子魚管理器
type luckyDiceManager struct {
	mu       sync.Mutex
	cooldown map[string]time.Time // playerID → 冷卻結束時間
}

// newLuckyDiceManager 建立幸運骰子魚管理器
func newLuckyDiceManager() *luckyDiceManager {
	return &luckyDiceManager{
		cooldown: make(map[string]time.Time),
	}
}

// isLuckyDiceFish 判斷是否為幸運骰子魚（T133）
func isLuckyDiceFish(defID string) bool {
	return defID == "T133"
}

// isOnCooldown 檢查玩家是否在冷卻中
func (m *luckyDiceManager) isOnCooldown(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	cd, ok := m.cooldown[playerID]
	if !ok {
		return false
	}
	return time.Now().Before(cd)
}

// setCooldown 設定玩家冷卻
func (m *luckyDiceManager) setCooldown(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cooldown[playerID] = time.Now().Add(time.Duration(LuckyDiceCooldownSec) * time.Second)
}

// calcDiceReward 根據骰子點數計算獎勵
// 點數 2（蛇眼）：20x betLevel
// 點數 7（最常見）：7x betLevel
// 點數 12（大六）：50x betLevel
// 其他：點數 × betLevel
func calcDiceReward(die1, die2, betLevel int) (int, string) {
	sum := die1 + die2
	var mult int
	var label string

	switch sum {
	case 2:
		mult = 20
		label = "🎲🎲 蛇眼！×20"
	case 7:
		mult = 7
		label = "🎲 幸運7！×7"
	case 11:
		mult = 11
		label = "🎲 幸運11！×11"
	case 12:
		mult = 50
		label = "🎲🎲 大六！×50"
	default:
		mult = sum
		label = fmt.Sprintf("🎲 點數%d！×%d", sum, sum)
	}

	reward := mult * betLevel
	if reward < 1 {
		reward = 1
	}
	return reward, label
}

// tryLuckyDiceFish 擊破 T133 後觸發幸運骰子（DAY-175）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryLuckyDiceFish(p *player.Player, triggerX, triggerY float64) {
	// 個人冷卻檢查
	if g.LuckyDice.isOnCooldown(p.ID) {
		return
	}
	g.LuckyDice.setCooldown(p.ID)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 預先決定骰子結果（公平性保證）
	die1 := rng.Intn(6) + 1 // 1-6
	die2 := rng.Intn(6) + 1 // 1-6
	sum := die1 + die2
	reward, label := calcDiceReward(die1, die2, p.BetLevel)

	log.Printf("[LuckyDice] player=%s die1=%d die2=%d sum=%d reward=%d",
		p.ID, die1, die2, sum, reward)

	// 廣播骰子開始滾動（個人）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyDiceFish,
		Payload: ws.LuckyDiceFishPayload{
			Phase:       "dice_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			RollMs:      LuckyDiceRollDurationMs,
			TriggerX:    triggerX,
			TriggerY:    triggerY,
		},
	})

	// 全服廣播：有人觸發骰子
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyDiceFish,
		Payload: ws.LuckyDiceFishPayload{
			Phase:      "dice_broadcast",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
		},
	})

	// 等待骰子滾動時間
	time.Sleep(time.Duration(LuckyDiceRollDurationMs) * time.Millisecond)

	// 發放獎勵
	p.AddReward(reward)

	// 廣播骰子結果（個人）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyDiceFish,
		Payload: ws.LuckyDiceFishPayload{
			Phase:      "dice_result",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Die1:       die1,
			Die2:       die2,
			Sum:        sum,
			Reward:     reward,
			Label:      label,
			NewBalance: p.Coins,
		},
	})

	// 全服公告：點數 ≥10 時廣播
	if sum >= LuckyDiceAnnounceThreshold {
		g.announceLuckyDiceFish(p.DisplayName, die1, die2, sum, reward, label)
	}

	// 特殊彩蛋：點數 12（大六）全服廣播
	if sum == 12 {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyDiceFish,
			Payload: ws.LuckyDiceFishPayload{
				Phase:      "dice_jackpot",
				PlayerID:   p.ID,
				PlayerName: p.DisplayName,
				Die1:       die1,
				Die2:       die2,
				Sum:        sum,
				Reward:     reward,
				Label:      label,
			},
		})
	}

	log.Printf("[LuckyDice] player=%s result: %d+%d=%d reward=%d label=%s",
		p.ID, die1, die2, sum, reward, label)
}

// announceLuckyDiceFish 全服公告幸運骰子魚（DAY-175）
func (g *Game) announceLuckyDiceFish(playerName string, die1, die2, sum, reward int, label string) {
	color := "#FFD700"
	if sum == 12 {
		color = "#FF4500" // 大六用橙紅色
	} else if sum == 2 {
		color = "#9400D3" // 蛇眼用紫色
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "lucky_dice_fish",
			"message":    fmt.Sprintf("🎲 %s 擲出 %d+%d=%d！%s 獲得 %d 金幣！", playerName, die1, die2, sum, label, reward),
			"color":      color,
			"duration":   4.0,
			"priority":   3,
		},
	})
}
