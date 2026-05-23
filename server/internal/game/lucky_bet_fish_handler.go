// lucky_bet_fish_handler.go — 幸運賭注魚系統（DAY-240）
// 業界原創「玩家主動風險決策+賭注翻倍」機制
//
// 設計：擊破 T198 後，玩家面臨「賭注選擇」（10 秒決策時間）：
//   - 選擇 A（保守）：下一次擊破 ×2.0 倍率，100% 觸發
//   - 選擇 B（激進）：下一次擊破 ×5.0 倍率，50% 觸發；失敗則 ×0.5 倍率
//   - 選擇 C（瘋狂）：下一次擊破 ×10.0 倍率，25% 觸發；失敗則 ×0.3 倍率
//   - 10 秒內未選擇 → 自動選擇 A（保守）
//   - 個人冷卻 30 秒
//
// 設計差異：
//   - 與幸運量子魚（DAY-228，50% 機率坍縮，玩家無法選擇）不同，賭注魚是「玩家主動選擇風險等級」，
//     讓玩家有「我要不要賭一把」的真實賭注感
//   - 「保守/激進/瘋狂」三個選項讓不同風格的玩家都有對應策略
//   - 「失敗懲罰」讓選擇 B/C 有真實風險，不是無腦選最高倍率
//   - 「10 秒決策時間」製造緊迫感，讓玩家在壓力下做決定
//   - 全服廣播玩家的選擇和結果，製造「看他敢不敢賭」的社交觀看感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyBetFishPersonalCD   = 30 * time.Second // 個人冷卻
	LuckyBetFishDecisionTime = 10 * time.Second // 決策時間

	// 選擇 A（保守）— 期望值 ×2.0，零風險
	LuckyBetFishAMult    = 2.0  // 倍率
	LuckyBetFishAChance  = 1.0  // 成功機率（100%）
	LuckyBetFishAFailMult = 1.0 // 失敗倍率（不會失敗）

	// 選擇 B（激進）— 期望值 ×2.0，高方差
	// 0.5 × 4.0 + 0.5 × 0.0 = 2.0x（與 A 相同期望值）
	LuckyBetFishBMult    = 4.0  // 倍率
	LuckyBetFishBChance  = 0.50 // 成功機率（50%）
	LuckyBetFishBFailMult = 0.0 // 失敗倍率（歸零，真實風險）

	// 選擇 C（瘋狂）— 期望值 ×2.0，極高方差
	// 0.25 × 8.0 + 0.75 × 0.0 = 2.0x（與 A 相同期望值）
	LuckyBetFishCMult    = 8.0  // 倍率
	LuckyBetFishCChance  = 0.25 // 成功機率（25%）
	LuckyBetFishCFailMult = 0.0 // 失敗倍率（歸零，真實風險）
)

// betChoice 賭注選擇
type betChoice string

const (
	BetChoiceA betChoice = "A" // 保守
	BetChoiceB betChoice = "B" // 激進
	BetChoiceC betChoice = "C" // 瘋狂
)

// betSession 個人賭注 session
type betSession struct {
	choice      betChoice // 玩家選擇（空字串=尚未選擇）
	decideUntil time.Time // 決策截止時間
	pendingMult float64   // 待套用的倍率（決策後設定）
	decided     bool      // 是否已決策
	instanceID  string    // 本次 session ID（防止重複觸發）
}

// luckyBetFishManager 幸運賭注魚管理器
type luckyBetFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 賭注 session（playerID → betSession）
	sessions map[string]*betSession
}

func newLuckyBetFishManager() *luckyBetFishManager {
	return &luckyBetFishManager{
		personalCooldowns: make(map[string]time.Time),
		sessions:          make(map[string]*betSession),
	}
}

// isLuckyBetFish 判斷是否為幸運賭注魚
func isLuckyBetFish(defID string) bool {
	return defID == "T198"
}

// getLuckyBetFishMult 取得賭注倍率（供 handleKill 使用）
// 回傳倍率並消耗 session（一次性）
// 注意：失敗時回傳 0.0（歸零），成功時回傳 2.0/4.0/8.0，未觸發時回傳 1.0
func (g *Game) getLuckyBetFishMult(playerID string) float64 {
	mgr := g.LuckyBetFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	sess, ok := mgr.sessions[playerID]
	if !ok || !sess.decided {
		return 1.0
	}

	mult := sess.pendingMult
	// 消耗 session（一次性）
	delete(mgr.sessions, playerID)

	log.Printf("[LuckyBetFish] player=%s consumed bet mult=%.1f", playerID, mult)
	return mult
}

// notifyBetChoice 玩家選擇賭注（由 handleBetChoice 呼叫）
func (g *Game) notifyBetChoice(p *player.Player, choice betChoice) {
	mgr := g.LuckyBetFish
	mgr.mu.Lock()

	sess, ok := mgr.sessions[p.ID]
	if !ok || sess.decided || time.Now().After(sess.decideUntil) {
		mgr.mu.Unlock()
		return
	}

	// 標記已決策
	sess.choice = choice
	sess.decided = true

	// 計算結果
	var mult float64
	var success bool
	var successChance float64

	switch choice {
	case BetChoiceA:
		mult = LuckyBetFishAMult
		success = true
		successChance = LuckyBetFishAChance
	case BetChoiceB:
		successChance = LuckyBetFishBChance
		if rand.Float64() < LuckyBetFishBChance {
			mult = LuckyBetFishBMult
			success = true
		} else {
			mult = LuckyBetFishBFailMult // 0.0 = 歸零（失去這次擊破獎勵）
			success = false
		}
	case BetChoiceC:
		successChance = LuckyBetFishCChance
		if rand.Float64() < LuckyBetFishCChance {
			mult = LuckyBetFishCMult
			success = true
		} else {
			mult = LuckyBetFishCFailMult // 0.0 = 歸零（失去這次擊破獎勵）
			success = false
		}
	default:
		// 未知選擇，預設 A
		mult = LuckyBetFishAMult
		success = true
		successChance = LuckyBetFishAChance
	}

	sess.pendingMult = mult
	instanceID := sess.instanceID
	mgr.mu.Unlock()

	log.Printf("[LuckyBetFish] player=%s chose %s, success=%v, mult=%.1f (chance=%.0f%%)",
		p.ID, choice, success, mult, successChance*100)

	// 廣播決策結果
	choiceLabel := map[betChoice]string{
		BetChoiceA: "保守 ×2.0",
		BetChoiceB: "激進 ×5.0",
		BetChoiceC: "瘋狂 ×10.0",
	}[choice]

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyBetFish,
		Payload: ws.LuckyBetFishPayload{
			Event:         "bet_decided",
			PlayerID:      p.ID,
			PlayerName:    p.DisplayName,
			Choice:        string(choice),
			ChoiceLabel:   choiceLabel,
			Success:       success,
			ResultMult:    mult,
			SuccessChance: successChance,
			InstanceID:    instanceID,
		},
	})

	// 全服公告
	var annMsg, annColor string
	if success {
		annMsg = fmt.Sprintf("🎲 %s 選擇「%s」— 成功！下次擊破 ×%.1f！",
			p.DisplayName, choiceLabel, mult)
		annColor = map[betChoice]string{
			BetChoiceA: "#27AE60",
			BetChoiceB: "#F39C12",
			BetChoiceC: "#E74C3C",
		}[choice]
	} else {
		annMsg = fmt.Sprintf("🎲 %s 選擇「%s」— 失敗！下次擊破 ×%.1f...",
			p.DisplayName, choiceLabel, mult)
		annColor = "#7F8C8D"
	}

	ann := g.Announce.Create(announce.EventLuckyBetFish, p.DisplayName, 0, map[string]string{
		"message": annMsg,
		"color":   annColor,
	})
	g.broadcastAnnouncement(ann)
}

// tryLuckyBetFish 擊破 T198 後觸發賭注選擇（供 handleKill 使用）
func (g *Game) tryLuckyBetFish(p *player.Player) {
	mgr := g.LuckyBetFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyBetFishPersonalCD)

	// 建立 session
	instanceID := fmt.Sprintf("bet_%s_%d", p.ID, time.Now().UnixNano())
	sess := &betSession{
		decideUntil: time.Now().Add(LuckyBetFishDecisionTime),
		instanceID:  instanceID,
	}
	mgr.sessions[p.ID] = sess
	mgr.mu.Unlock()

	log.Printf("[LuckyBetFish] player=%s triggered bet fish, decision window=%v",
		p.ID, LuckyBetFishDecisionTime)

	// 個人訊息：顯示選擇介面
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyBetFish,
		Payload: ws.LuckyBetFishPayload{
			Event:      "bet_start",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			DecisionSec: int(LuckyBetFishDecisionTime.Seconds()),
			InstanceID: instanceID,
			// 三個選項的說明
			OptionA: ws.BetOption{
				Choice:        "A",
				Label:         "保守",
				Mult:          LuckyBetFishAMult,
				SuccessChance: LuckyBetFishAChance,
				FailMult:      LuckyBetFishAFailMult,
			},
			OptionB: ws.BetOption{
				Choice:        "B",
				Label:         "激進",
				Mult:          LuckyBetFishBMult,
				SuccessChance: LuckyBetFishBChance,
				FailMult:      LuckyBetFishBFailMult,
			},
			OptionC: ws.BetOption{
				Choice:        "C",
				Label:         "瘋狂",
				Mult:          LuckyBetFishCMult,
				SuccessChance: LuckyBetFishCChance,
				FailMult:      LuckyBetFishCFailMult,
			},
		},
	})

	// 全服廣播：通知其他玩家有人觸發了賭注魚
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyBetFish,
		Payload: ws.LuckyBetFishPayload{
			Event:      "bet_broadcast",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			InstanceID: instanceID,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyBetFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🎲 %s 觸發幸運賭注魚！正在選擇賭注...", p.DisplayName),
		"color":   "#9B59B6",
	})
	g.broadcastAnnouncement(ann)

	// 10 秒後自動決策（若玩家未選擇）
	go func() {
		time.Sleep(LuckyBetFishDecisionTime)

		mgr.mu.Lock()
		sess, ok := mgr.sessions[p.ID]
		if !ok || sess.instanceID != instanceID || sess.decided {
			mgr.mu.Unlock()
			return
		}
		// 自動選擇 A（保守）
		sess.choice = BetChoiceA
		sess.decided = true
		sess.pendingMult = LuckyBetFishAMult
		mgr.mu.Unlock()

		log.Printf("[LuckyBetFish] player=%s auto-chose A (timeout)", p.ID)

		// 廣播超時自動選擇
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyBetFish,
			Payload: ws.LuckyBetFishPayload{
				Event:       "bet_timeout",
				PlayerID:    p.ID,
				PlayerName:  p.DisplayName,
				Choice:      "A",
				ChoiceLabel: "保守 ×2.0（自動）",
				Success:     true,
				ResultMult:  LuckyBetFishAMult,
				InstanceID:  instanceID,
			},
		})
	}()
}

// handleLuckyBetChoice 處理玩家賭注選擇（由 HandleMessage 呼叫）
func (g *Game) handleLuckyBetChoice(p *player.Player, msg *ws.Message) {
	var payload ws.LuckyBetChoicePayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		log.Printf("[LuckyBetFish] invalid bet choice payload: %v", err)
		return
	}

	// 驗證選擇
	choice := betChoice(payload.Choice)
	if choice != BetChoiceA && choice != BetChoiceB && choice != BetChoiceC {
		log.Printf("[LuckyBetFish] invalid choice=%s from player=%s", payload.Choice, p.ID)
		return
	}

	// 驗證 instanceID 匹配
	mgr := g.LuckyBetFish
	mgr.mu.Lock()
	sess, ok := mgr.sessions[p.ID]
	if !ok || sess.decided || sess.instanceID != payload.InstanceID {
		mgr.mu.Unlock()
		return
	}
	mgr.mu.Unlock()

	g.notifyBetChoice(p, choice)
}
