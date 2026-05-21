// anglerfish_handler.go — 巨型鮟鱇魚電擊寶箱 handler（DAY-145）
// 業界依據：jiligames.com 2026「Giant Anglerfish can shoot electricity to open treasure chests,
// giant crocodiles awaken to hunt fish on the fish farm to accumulate big prizes!」
// 擊破 T109 後觸發電擊，電流傳導到附近的寶箱目標（T102），強制開啟寶箱獲得額外獎勵
package game

import (
	"fmt"
	"log"
	"math"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// AnglerfishShockRadius 電擊傳導半徑（px）
	AnglerfishShockRadius = 250.0
	// AnglerfishChestRewardMult 強制開啟寶箱的獎勵倍率（比正常擊破略低）
	AnglerfishChestRewardMult = 0.80
	// AnglerfishShockDelayMs 電擊傳導延遲（ms，製造電流跳躍感）
	AnglerfishShockDelayMs = 120
	// AnglerfishAnnounceThreshold 全服公告門檻（開啟寶箱數）
	AnglerfishAnnounceThreshold = 2
)

// isAnglerfish 判斷是否為巨型鮟鱇魚（T109）
func isAnglerfish(defID string) bool {
	return defID == "T109"
}

// isChestTarget 判斷是否為寶箱目標（T102）
func isChestTarget(defID string) bool {
	return defID == "T102"
}

// anglerfishChestEntry 電擊開啟的寶箱記錄
type anglerfishChestEntry struct {
	instanceID string
	defID      string
	x, y       float64
	multiplier float64
}

// tryAnglerfishShock 擊破 T109 後觸發電擊開寶箱（DAY-145）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryAnglerfishShock(p *player.Player, triggerID string, triggerX, triggerY float64) {
	// 收集電擊範圍內的寶箱目標（T102）
	g.mu.RLock()
	var chestTargets []anglerfishChestEntry
	for id, t := range g.Targets {
		if id == triggerID || t.HP <= 0 {
			continue
		}
		if !isChestTarget(t.DefID) {
			continue
		}
		dx := t.X - triggerX
		dy := t.Y - triggerY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= AnglerfishShockRadius {
			chestTargets = append(chestTargets, anglerfishChestEntry{
				instanceID: id,
				defID:      t.DefID,
				x:          t.X,
				y:          t.Y,
				multiplier: t.Multiplier,
			})
		}
	}
	g.mu.RUnlock()

	if len(chestTargets) == 0 {
		// 沒有寶箱，只廣播電擊特效（視覺效果）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgAnglerfishShock,
			Payload: ws.AnglerfishShockPayload{
				TriggerID:  triggerID,
				TriggerX:   triggerX,
				TriggerY:   triggerY,
				Phase:      "shock_start",
				ChestIDs:   []string{},
			},
		})
		return
	}

	// 廣播電擊開始（讓 Client 播放電流動畫）
	chestIDs := make([]string, 0, len(chestTargets))
	for _, ct := range chestTargets {
		chestIDs = append(chestIDs, ct.instanceID)
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnglerfishShock,
		Payload: ws.AnglerfishShockPayload{
			TriggerID:  triggerID,
			TriggerX:   triggerX,
			TriggerY:   triggerY,
			Phase:      "shock_start",
			ChestIDs:   chestIDs,
		},
	})

	// 逐一電擊開啟寶箱（每 120ms 一個，製造電流跳躍感）
	totalReward := 0
	var openedEntries []ws.AnglerfishChestEntry

	for i, ct := range chestTargets {
		if i > 0 {
			time.Sleep(time.Duration(AnglerfishShockDelayMs) * time.Millisecond)
		}

		g.mu.Lock()
		t, ok := g.Targets[ct.instanceID]
		if !ok || t.HP <= 0 {
			g.mu.Unlock()
			continue
		}
		reward := int(float64(p.BetLevel) * ct.multiplier * AnglerfishChestRewardMult)
		if reward < 1 {
			reward = 1
		}
		t.HP = 0
		delete(g.Targets, ct.instanceID)
		g.mu.Unlock()

		totalReward += reward
		openedEntries = append(openedEntries, ws.AnglerfishChestEntry{
			InstanceID: ct.instanceID,
			Multiplier: ct.multiplier,
			Reward:     reward,
			X:          ct.x,
			Y:          ct.y,
		})

		// 廣播寶箱被電擊開啟
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetKill,
			Payload: ws.TargetKillPayload{
				InstanceID: ct.instanceID,
				KillerID:   p.ID,
				Reward:     reward,
				Multiplier: ct.multiplier,
			},
		})

		log.Printf("[Anglerfish] shock chest[%d] id=%s mult=%.0f reward=%d",
			i, ct.instanceID, ct.multiplier, reward)
	}

	if totalReward <= 0 {
		return
	}

	// 發放總獎勵
	p.AddReward(totalReward)

	// 廣播電擊結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnglerfishShock,
		Payload: ws.AnglerfishShockPayload{
			TriggerID:     triggerID,
			TriggerX:      triggerX,
			TriggerY:      triggerY,
			Phase:         "result",
			ChestIDs:      chestIDs,
			OpenedChests:  openedEntries,
			TotalReward:   totalReward,
			KillerID:      p.ID,
			KillerName:    p.DisplayName,
		},
	})

	// 個人結果通知
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgReward,
		Payload: ws.RewardPayload{
			Source:     "anglerfish_shock",
			Amount:     totalReward,
			Multiplier: float64(len(openedEntries)),
			NewBalance: p.Coins,
		},
	})

	// 全服公告：開啟 ≥2 個寶箱
	if len(openedEntries) >= AnglerfishAnnounceThreshold {
		g.announceAnglerfishShock(p.DisplayName, len(openedEntries), totalReward)
	}

	log.Printf("[Anglerfish] player=%s opened=%d chests total_reward=%d",
		p.ID, len(openedEntries), totalReward)
}

// announceAnglerfishShock 全服公告鮟鱇魚電擊開寶箱（DAY-145）
func (g *Game) announceAnglerfishShock(playerName string, chestCount int, reward int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "anglerfish_shock",
			"message":    fmt.Sprintf("⚡ %s 的巨型鮟鱇魚電擊開啟 %d 個寶箱！獲得 %d 金幣！", playerName, chestCount, reward),
			"color":      "#00BFFF",
			"duration":   4.0,
			"priority":   2,
		},
	})
}
