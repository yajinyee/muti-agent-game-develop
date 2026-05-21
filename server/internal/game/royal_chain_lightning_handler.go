// royal_chain_lightning_handler.go — 皇家閃電鰻持續連鎖電擊系統 handler（DAY-156）
// 業界依據：royal-fishing.co.uk 2026「Creates chain lightning that shocks nearby fish
//   consecutively until targeting turns off. Devastating against clustered schools.」
//   + royal-fishing.uk 2026「The 60x lightning eel creates chain reactions that jump between
//   nearby fish. Once activated, electric shocks continue spreading until targeting disengages,
//   creating cascading capture sequences across the underwater battlefield.」
// 設計：擊破 T118 後觸發持續連鎖電擊，每 200ms 跳一次（最多 15 跳），每跳 300px 範圍
//   比 T103 閃電鰻（一次性 5 跳/200px/50%）更強：15 跳/300px/60%
//   每跳廣播讓所有玩家看到「電擊在連鎖跳躍」，製造視覺爽感
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// isRoyalChainLightning 判斷是否為皇家閃電鰻（T118）
func isRoyalChainLightning(defID string) bool {
	return defID == "T118"
}

// notifyRoyalChainLightningKill 擊破皇家閃電鰻後觸發持續連鎖電擊（由 handleKill 呼叫）
func (g *Game) notifyRoyalChainLightningKill(p *player.Player, triggerID string, triggerX, triggerY float64) {
	log.Printf("[RoyalChainLightning] player=%s triggered chain lightning at (%.0f, %.0f)", p.ID, triggerX, triggerY)
	go g.runRoyalChainLightning(p, triggerID, triggerX, triggerY)
}

// runRoyalChainLightning 執行持續連鎖電擊（goroutine）
// 最多 15 跳，每跳 200ms，每跳 300px 範圍，60% 擊破機率
func (g *Game) runRoyalChainLightning(p *player.Player, triggerID string, triggerX, triggerY float64) {
	const maxJumps = 15
	const jumpInterval = 200 * time.Millisecond
	const jumpRadius = 300.0
	const killChance = 0.60

	// 廣播連鎖開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRoyalChainLightning,
		Payload: ws.RoyalChainLightningPayload{
			TriggerID:  triggerID,
			TriggerX:   triggerX,
			TriggerY:   triggerY,
			KillerID:   p.ID,
			KillerName: p.DisplayName,
			Phase:      "chain_start",
		},
	})

	allEntries := make([]ws.RoyalChainLightningEntry, 0, maxJumps)
	totalReward := 0
	currentX := triggerX
	currentY := triggerY
	usedIDs := map[string]bool{triggerID: true} // 已電擊過的目標，不重複

	for jumpIdx := 1; jumpIdx <= maxJumps; jumpIdx++ {
		time.Sleep(jumpInterval)

		// 找下一個跳躍目標（在當前位置 300px 範圍內，未被電擊過）
		g.mu.RLock()
		var bestID string
		var bestDist float64 = math.MaxFloat64
		var bestX, bestY float64
		for _, t := range g.Targets {
			if t.HP <= 0 || usedIDs[t.InstanceID] {
				continue
			}
			dx := t.X - currentX
			dy := t.Y - currentY
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= jumpRadius && dist < bestDist {
				bestDist = dist
				bestID = t.InstanceID
				bestX = t.X
				bestY = t.Y
			}
		}
		g.mu.RUnlock()

		if bestID == "" {
			// 沒有可跳躍的目標，連鎖結束
			log.Printf("[RoyalChainLightning] player=%s chain ended at jump=%d (no target)", p.ID, jumpIdx-1)
			break
		}

		usedIDs[bestID] = true

		// 處理電擊
		g.mu.Lock()
		t, ok := g.Targets[bestID]
		if !ok || t.HP <= 0 {
			g.mu.Unlock()
			continue
		}

		entry := ws.RoyalChainLightningEntry{
			InstanceID: t.InstanceID,
			DefID:      t.DefID,
			Multiplier: t.Multiplier,
			JumpIndex:  jumpIdx,
			FromX:      currentX,
			FromY:      currentY,
			ToX:        bestX,
			ToY:        bestY,
		}

		if rand.Float64() < killChance {
			// 擊破！獎勵 = 倍率 × betLevel × 0.60（連鎖電擊是連帶效果，比直接擊破低）
			reward := int(float64(p.BetLevel) * t.Multiplier * 0.60)
			if reward < 1 {
				reward = 1
			}
			entry.Killed = true
			entry.Reward = reward
			totalReward += reward
			p.Coins += reward
			if p.Coins > p.MaxCoins {
				p.MaxCoins = p.Coins
			}
			t.HP = 0
			delete(g.Targets, bestID)

			// 廣播目標被擊破
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgTargetKill,
				Payload: ws.TargetKillPayload{
					InstanceID: t.InstanceID,
					KillerID:   p.ID,
					Reward:     reward,
					Multiplier: t.Multiplier,
				},
			})
		}
		g.mu.Unlock()

		allEntries = append(allEntries, entry)

		// 廣播本跳（讓 Client 播放電擊跳躍動畫）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgRoyalChainLightning,
			Payload: ws.RoyalChainLightningPayload{
				TriggerID:  triggerID,
				TriggerX:   triggerX,
				TriggerY:   triggerY,
				KillerID:   p.ID,
				KillerName: p.DisplayName,
				Phase:      "jump",
				JumpIndex:  jumpIdx,
				JumpEntry:  &entry,
			},
		})

		// 更新當前位置（電擊從這個目標繼續跳）
		currentX = bestX
		currentY = bestY
	}

	// 廣播最終結果
	killedCount := 0
	for _, e := range allEntries {
		if e.Killed {
			killedCount++
		}
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRoyalChainLightning,
		Payload: ws.RoyalChainLightningPayload{
			TriggerID:   triggerID,
			TriggerX:    triggerX,
			TriggerY:    triggerY,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			Phase:       "result",
			AllEntries:  allEntries,
			TotalReward: totalReward,
			TotalJumps:  len(allEntries),
		},
	})

	// 全服公告：連鎖大豐收（≥6 跳或 ≥4 個擊破）
	if len(allEntries) >= 6 || killedCount >= 4 {
		g.announceRoyalChainLightning(p.DisplayName, len(allEntries), killedCount, totalReward)
	}

	log.Printf("[RoyalChainLightning] player=%s jumps=%d killed=%d total_reward=%d",
		p.ID, len(allEntries), killedCount, totalReward)
}

// announceRoyalChainLightning 全服公告皇家閃電鰻連鎖大豐收
func (g *Game) announceRoyalChainLightning(playerName string, jumps int, killedCount int, totalReward int) {
	var msg string
	if killedCount >= 4 {
		msg = fmt.Sprintf("⚡ %s 的皇家閃電鰻連鎖 %d 跳，擊破 %d 個目標，獲得 %d 金幣！", playerName, jumps, killedCount, totalReward)
	} else {
		msg = fmt.Sprintf("⚡ %s 的皇家閃電鰻連鎖電擊 %d 跳，獲得 %d 金幣！", playerName, jumps, totalReward)
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "royal_chain_lightning",
			"message":    msg,
			"color":      "#00BFFF",
			"duration":   4.0,
			"priority":   2,
		},
	})
}
