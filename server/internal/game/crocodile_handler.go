// crocodile_handler.go — 巨型鹹水鱷魚獵魚累積 handler（DAY-146）
// 業界依據：jiligames.com 2026「giant crocodiles awaken to hunt fish on the fish farm to accumulate big prizes!」
// + megafishinggame.top「Giant Saltwater Crocodile」
// 擊破 T110 後觸發「鱷魚獵魚」模式：鱷魚在 8 秒內自動獵殺場上的普通目標，累積獎勵給觸發玩家
package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// CrocodileHuntDuration 鱷魚獵魚持續時間（秒）
	CrocodileHuntDuration = 8
	// CrocodileHuntInterval 每次獵殺間隔（ms）
	CrocodileHuntIntervalMs = 800
	// CrocodileMaxHunts 最多獵殺目標數
	CrocodileMaxHunts = 6
	// CrocodileRewardMult 獵殺獎勵倍率（觸發玩家獲得目標獎勵的 30%）
	CrocodileRewardMult = 0.30
	// CrocodileAnnounceThreshold 全服公告門檻（獵殺數）
	CrocodileAnnounceThreshold = 4
)

// isCrocodile 判斷是否為巨型鹹水鱷魚（T110）
func isCrocodile(defID string) bool {
	return defID == "T110"
}

// isBasicTarget 判斷是否為普通目標（T001-T006，鱷魚獵殺對象）
func isBasicTarget(defID string) bool {
	switch defID {
	case "T001", "T002", "T003", "T004", "T005", "T006":
		return true
	}
	return false
}

// crocodileHuntEntry 鱷魚獵殺記錄
type crocodileHuntEntry struct {
	instanceID string
	defID      string
	multiplier float64
	reward     int
}

// tryCrocodileHunt 擊破 T110 後觸發鱷魚獵魚模式（DAY-146）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryCrocodileHunt(p *player.Player, triggerID string, triggerX, triggerY float64) {
	// 廣播鱷魚覺醒（讓 Client 播放覺醒動畫）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCrocodileHunt,
		Payload: ws.CrocodileHuntPayload{
			TriggerID:    triggerID,
			TriggerX:     triggerX,
			TriggerY:     triggerY,
			Phase:        "awaken",
			HuntDuration: CrocodileHuntDuration,
			MaxHunts:     CrocodileMaxHunts,
			KillerID:     p.ID,
			KillerName:   p.DisplayName,
		},
	})

	log.Printf("[Crocodile] player=%s crocodile awakened, hunting for %ds", p.ID, CrocodileHuntDuration)

	// 鱷魚獵魚循環
	var huntedEntries []ws.CrocodileHuntEntry
	totalReward := 0
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for hunt := 0; hunt < CrocodileMaxHunts; hunt++ {
		time.Sleep(time.Duration(CrocodileHuntIntervalMs) * time.Millisecond)

		// 隨機選擇一個普通目標
		g.mu.RLock()
		var candidates []string
		for id, t := range g.Targets {
			if t.HP > 0 && isBasicTarget(t.DefID) {
				candidates = append(candidates, id)
			}
		}
		g.mu.RUnlock()

		if len(candidates) == 0 {
			break // 沒有普通目標了
		}

		// 隨機選一個
		targetID := candidates[rng.Intn(len(candidates))]

		g.mu.Lock()
		t, ok := g.Targets[targetID]
		if !ok || t.HP <= 0 || !isBasicTarget(t.DefID) {
			g.mu.Unlock()
			continue
		}

		reward := int(float64(p.BetLevel) * t.Multiplier * CrocodileRewardMult)
		if reward < 1 {
			reward = 1
		}
		huntedEntry := ws.CrocodileHuntEntry{
			InstanceID: t.InstanceID,
			DefID:      t.DefID,
			Multiplier: t.Multiplier,
			Reward:     reward,
			HuntIndex:  hunt,
		}
		t.HP = 0
		delete(g.Targets, targetID)
		g.mu.Unlock()

		totalReward += reward
		huntedEntries = append(huntedEntries, huntedEntry)

		// 廣播鱷魚獵殺（讓 Client 播放獵殺動畫）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgCrocodileHunt,
			Payload: ws.CrocodileHuntPayload{
				TriggerID:  triggerID,
				Phase:      "hunt",
				HuntIndex:  hunt,
				HuntedID:   targetID,
				HuntReward: reward,
				KillerID:   p.ID,
			},
		})

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetKill,
			Payload: ws.TargetKillPayload{
				InstanceID: targetID,
				KillerID:   p.ID,
				Reward:     reward,
				Multiplier: t.Multiplier,
			},
		})

		log.Printf("[Crocodile] hunt[%d] target=%s mult=%.0f reward=%d",
			hunt, targetID, t.Multiplier, reward)
	}

	if totalReward <= 0 {
		// 廣播結束（即使沒有獵殺）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgCrocodileHunt,
			Payload: ws.CrocodileHuntPayload{
				TriggerID:     triggerID,
				Phase:         "result",
				HuntedTargets: huntedEntries,
				TotalReward:   0,
				KillerID:      p.ID,
				KillerName:    p.DisplayName,
			},
		})
		return
	}

	// 發放總獎勵
	p.AddReward(totalReward)

	// 廣播鱷魚獵魚結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCrocodileHunt,
		Payload: ws.CrocodileHuntPayload{
			TriggerID:     triggerID,
			TriggerX:      triggerX,
			TriggerY:      triggerY,
			Phase:         "result",
			HuntedTargets: huntedEntries,
			TotalReward:   totalReward,
			KillerID:      p.ID,
			KillerName:    p.DisplayName,
		},
	})

	// 個人結果通知
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgReward,
		Payload: ws.RewardPayload{
			Source:     "crocodile_hunt",
			Amount:     totalReward,
			Multiplier: float64(len(huntedEntries)),
			NewBalance: p.Coins,
		},
	})

	// 全服公告：獵殺 ≥4 個目標
	if len(huntedEntries) >= CrocodileAnnounceThreshold {
		g.announceCrocodileHunt(p.DisplayName, len(huntedEntries), totalReward)
	}

	log.Printf("[Crocodile] player=%s hunted=%d total_reward=%d",
		p.ID, len(huntedEntries), totalReward)
}

// announceCrocodileHunt 全服公告鱷魚獵魚（DAY-146）
func (g *Game) announceCrocodileHunt(playerName string, huntCount int, reward int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "crocodile_hunt",
			"message":    fmt.Sprintf("🐊 %s 的巨型鱷魚獵殺 %d 個目標！累積獲得 %d 金幣！", playerName, huntCount, reward),
			"color":      "#228B22",
			"duration":   4.5,
			"priority":   2,
		},
	})
}
