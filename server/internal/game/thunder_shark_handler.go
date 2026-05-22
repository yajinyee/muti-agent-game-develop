// thunder_shark_handler.go — 雷霆鯊魚全場連鎖閃電系統 handler（DAY-181）
// 業界依據：JILI Jackpot Fishing「Thunder Shark brings unique abilities —
// chain lightning that jumps between nearby fish, with no distance limit」
// 擊破 T139 後觸發「雷霆連鎖閃電」：
//   1. 閃電從觸發位置開始，全場隨機跳躍（不限距離），最多 20 跳
//   2. 每跳有 75% 機率擊破目標，獎勵 0.65x 倍率
//   3. 每跳間隔 150ms，製造「閃電連鎖」的視覺爽感
//   4. 全服廣播每一跳，讓所有玩家看到閃電在全場跳躍
// 設計差異：
//   - 與 T103 閃電鰻（5跳/200px範圍/50%機率）不同，雷霆鯊魚是「全場無限距離」（20跳/75%機率）
//   - 與 T118 皇家閃電鰻（15跳/300px範圍/60%機率）不同，雷霆鯊魚是「全場隨機跳躍」（不限距離），
//     讓玩家看到閃電在全場「隨機亂跳」的混亂爽感
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
	// ThunderSharkMaxJumps 最大跳躍次數
	ThunderSharkMaxJumps = 20
	// ThunderSharkKillChance 每跳擊破機率（75%）
	ThunderSharkKillChance = 0.75
	// ThunderSharkRewardMult 擊破獎勵倍率（比直接擊破低，平衡 RTP）
	ThunderSharkRewardMult = 0.65
	// ThunderSharkJumpDelayMs 每跳間隔（ms）
	ThunderSharkJumpDelayMs = 150
	// ThunderSharkAnnounceMinJumps 全服公告最低跳數門檻
	ThunderSharkAnnounceMinJumps = 10
)

// isThunderShark 判斷是否為雷霆鯊魚（T139）
func isThunderShark(defID string) bool {
	return defID == "T139"
}

// thunderSharkTarget 雷霆閃電目標（內部結構）
type thunderSharkTarget struct {
	instanceID string
	defID      string
	x, y       float64
	multiplier float64
}

// tryThunderSharkChain 擊破 T139 後觸發全場連鎖閃電（DAY-181）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryThunderSharkChain(p *player.Player, triggerID string, triggerX, triggerY float64) {
	// 廣播閃電開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgThunderSharkChain,
		Payload: ws.ThunderSharkChainPayload{
			Phase:      "chain_start",
			TriggerID:  triggerID,
			TriggerX:   triggerX,
			TriggerY:   triggerY,
			KillerID:   p.ID,
			KillerName: p.DisplayName,
		},
	})

	totalReward := 0
	totalKills := 0
	totalJumps := 0
	usedTargets := map[string]bool{triggerID: true}

	for jump := 0; jump < ThunderSharkMaxJumps; jump++ {
		time.Sleep(time.Duration(ThunderSharkJumpDelayMs) * time.Millisecond)

		// 從全場隨機選一個未被閃電過的存活目標
		g.mu.RLock()
		var candidates []thunderSharkTarget
		for id, t := range g.Targets {
			if usedTargets[id] || t.HP <= 0 || t.DefID == "B001" {
				continue
			}
			candidates = append(candidates, thunderSharkTarget{
				instanceID: t.InstanceID,
				defID:      t.DefID,
				x:          t.X,
				y:          t.Y,
				multiplier: t.Multiplier,
			})
		}
		g.mu.RUnlock()

		if len(candidates) == 0 {
			break // 沒有可跳躍的目標，結束
		}

		// 隨機選一個目標
		dt := candidates[rand.Intn(len(candidates))]
		usedTargets[dt.instanceID] = true
		totalJumps++

		// 廣播閃電跳躍（全服）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgThunderSharkChain,
			Payload: ws.ThunderSharkChainPayload{
				Phase:      fmt.Sprintf("jump_%d", jump+1),
				TriggerID:  triggerID,
				JumpTarget: dt.instanceID,
				JumpX:      dt.x,
				JumpY:      dt.y,
				JumpNum:    jump + 1,
				KillerID:   p.ID,
			},
		})

		// 75% 機率擊破
		if rand.Float64() >= ThunderSharkKillChance {
			continue // 未擊破，繼續下一跳
		}

		// 擊破目標
		g.mu.Lock()
		t, ok := g.Targets[dt.instanceID]
		if !ok || t.HP <= 0 {
			g.mu.Unlock()
			continue
		}
		reward := int(float64(p.BetLevel) * dt.multiplier * ThunderSharkRewardMult)
		if reward < 1 {
			reward = 1
		}
		t.HP = 0
		delete(g.Targets, dt.instanceID)
		g.mu.Unlock()

		totalReward += reward
		totalKills++

		// 廣播目標擊破
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetKill,
			Payload: ws.TargetKillPayload{
				InstanceID: dt.instanceID,
				KillerID:   p.ID,
				Reward:     reward,
				Multiplier: dt.multiplier,
			},
		})

		log.Printf("[ThunderShark] jump[%d] target=%s mult=%.0f reward=%d",
			jump+1, dt.instanceID, dt.multiplier, reward)
	}

	if totalReward > 0 {
		p.AddReward(totalReward)
	}

	// 廣播連鎖結果（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgThunderSharkChain,
		Payload: ws.ThunderSharkChainPayload{
			Phase:       "result",
			TriggerID:   triggerID,
			TotalJumps:  totalJumps,
			TotalKills:  totalKills,
			TotalReward: totalReward,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
		},
	})

	// 個人結果通知
	if totalReward > 0 {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgReward,
			Payload: ws.RewardPayload{
				Source:     "thunder_shark",
				Amount:     totalReward,
				Multiplier: float64(totalKills),
				NewBalance: p.Coins,
			},
		})
	}

	// 全服公告：跳數 ≥ 10
	if totalJumps >= ThunderSharkAnnounceMinJumps {
		g.announceThunderSharkChain(p.DisplayName, totalJumps, totalKills, totalReward)
	}

	log.Printf("[ThunderShark] player=%s jumps=%d kills=%d total_reward=%d",
		p.ID, totalJumps, totalKills, totalReward)
}

// announceThunderSharkChain 全服公告雷霆鯊魚連鎖閃電（DAY-181）
func (g *Game) announceThunderSharkChain(playerName string, jumps, kills, reward int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "thunder_shark_chain",
			"message":    fmt.Sprintf("⚡ %s 的雷霆鯊魚連鎖 %d 跳！擊破 %d 個目標！獲得 %d 金幣！", playerName, jumps, kills, reward),
			"color":      "#FFD700", // 金黃色（閃電感）
			"duration":   5.0,
			"priority":   3,
		},
	})
}
