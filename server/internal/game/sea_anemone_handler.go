// sea_anemone_handler.go — 海葵觸手攻擊系統 handler（DAY-174）
// 業界依據：JILI Jackpot Fishing「Sea Anemone introduces unique effects —
// tentacle attacks that spread to nearby fish, chain lightning or explosive torpedoes」
// jackpotfishing-game.com 2026「Sea Anemone introduce unique effects, such as chain lightning
// or explosive torpedoes, adding layers of strategy and excitement」
// 擊破 T132 後觸手向 8 個方向延伸，每個方向命中最近的目標（300px 範圍），
// 命中目標有 70% 機率擊破，獲得 0.6x 倍率獎勵
// 設計差異：與閃電鰻（連鎖跳躍，隨機目標）不同，海葵是**方向性觸手**（8方向固定延伸），
// 讓玩家感受到「觸手從中心向四周蔓延」的視覺爽感；
// 與炸彈（圓形範圍爆炸）不同，海葵是「線性觸手」（每個方向只命中最近的一個目標），
// 更有策略性（玩家可以預判觸手方向）
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

const (
	// SeaAnemoneRadius 觸手攻擊半徑（px）
	SeaAnemoneRadius = 300.0
	// SeaAnemoneKillChance 觸手命中目標的擊破機率
	SeaAnemoneKillChance = 0.70
	// SeaAnemoneRewardMult 觸手擊破獎勵倍率（比直接擊破低）
	SeaAnemoneRewardMult = 0.60
	// SeaTentacleIntervalMs 觸手延伸間隔（ms）
	SeaTentacleIntervalMs = 100
	// SeaAnemoneAnnounceThreshold 全服公告門檻（擊破數）
	SeaAnemoneAnnounceThreshold = 4
	// SeaAnemoneDirections 觸手方向數（8方向）
	SeaAnemoneDirections = 8
)

// isSeaAnemone 判斷是否為海葵（T132）
func isSeaAnemone(defID string) bool {
	return defID == "T132"
}

// seaAnemoneHitEntry 觸手命中記錄
type seaAnemoneHitEntry struct {
	instanceID string
	defID      string
	x, y       float64
	multiplier float64
	direction  int // 0-7（8方向）
	isKill     bool
	reward     int
}

// trySeaAnemone 擊破 T132 後觸發觸手攻擊（DAY-174）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) trySeaAnemone(p *player.Player, triggerID string, triggerX, triggerY float64) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 廣播觸手攻擊開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgSeaAnemone,
		Payload: ws.SeaAnemonePayload{
			Phase:      "tentacle_start",
			TriggerID:  triggerID,
			TriggerX:   triggerX,
			TriggerY:   triggerY,
			KillerID:   p.ID,
			KillerName: p.DisplayName,
			Directions: SeaAnemoneDirections,
		},
	})

	log.Printf("[SeaAnemone] player=%s triggered at (%.0f,%.0f), extending %d tentacles",
		p.ID, triggerX, triggerY, SeaAnemoneDirections)

	// 8方向觸手逐一延伸
	var hitEntries []ws.SeaAnemoneHitEntry
	totalReward := 0
	killCount := 0

	for dir := 0; dir < SeaAnemoneDirections; dir++ {
		if dir > 0 {
			time.Sleep(time.Duration(SeaTentacleIntervalMs) * time.Millisecond)
		}

		// 計算方向角度（0=右, 45=右下, 90=下, ...）
		angle := float64(dir) * (360.0 / float64(SeaAnemoneDirections))
		angleRad := angle * math.Pi / 180.0
		dirX := math.Cos(angleRad)
		dirY := math.Sin(angleRad)

		// 找這個方向上最近的目標
		g.mu.RLock()
		var bestTarget *seaAnemoneHitEntry
		bestScore := math.MaxFloat64

		for id, t := range g.Targets {
			if id == triggerID || t.HP <= 0 {
				continue
			}
			dx := t.X - triggerX
			dy := t.Y - triggerY
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > SeaAnemoneRadius {
				continue
			}

			// 計算目標與觸手方向的對齊程度（點積）
			if dist < 1.0 {
				continue
			}
			normDx := dx / dist
			normDy := dy / dist
			dot := normDx*dirX + normDy*dirY

			// 只考慮方向對齊的目標（點積 > 0.5，即 ±60° 範圍內）
			if dot < 0.5 {
				continue
			}

			// 評分：距離越近越好，方向越對齊越好
			score := dist / dot
			if score < bestScore {
				bestScore = score
				bestTarget = &seaAnemoneHitEntry{
					instanceID: id,
					defID:      t.DefID,
					x:          t.X,
					y:          t.Y,
					multiplier: t.Multiplier,
					direction:  dir,
				}
			}
		}
		g.mu.RUnlock()

		if bestTarget == nil {
			// 這個方向沒有目標，廣播空觸手
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgSeaAnemone,
				Payload: ws.SeaAnemonePayload{
					Phase:     "tentacle_miss",
					TriggerX:  triggerX,
					TriggerY:  triggerY,
					Direction: dir,
					Angle:     angle,
					KillerID:  p.ID,
				},
			})
			continue
		}

		// 嘗試擊破目標
		g.mu.Lock()
		t, ok := g.Targets[bestTarget.instanceID]
		if !ok || t.HP <= 0 {
			g.mu.Unlock()
			continue
		}

		isKill := rng.Float64() < SeaAnemoneKillChance
		reward := 0
		if isKill {
			reward = int(float64(p.BetLevel) * bestTarget.multiplier * SeaAnemoneRewardMult)
			if reward < 1 {
				reward = 1
			}
			t.HP = 0
			delete(g.Targets, bestTarget.instanceID)
			killCount++
			totalReward += reward
		} else {
			// 未擊破，造成傷害
			t.HP -= 1
			if t.HP < 0 {
				t.HP = 0
			}
		}
		g.mu.Unlock()

		bestTarget.isKill = isKill
		bestTarget.reward = reward
		hitEntries = append(hitEntries, ws.SeaAnemoneHitEntry{
			InstanceID: bestTarget.instanceID,
			DefID:      bestTarget.defID,
			X:          bestTarget.x,
			Y:          bestTarget.y,
			Multiplier: bestTarget.multiplier,
			Direction:  dir,
			Angle:      angle,
			IsKill:     isKill,
			Reward:     reward,
		})

		// 廣播觸手命中（全服）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgSeaAnemone,
			Payload: ws.SeaAnemonePayload{
				Phase:      "tentacle_hit",
				TriggerX:   triggerX,
				TriggerY:   triggerY,
				Direction:  dir,
				Angle:      angle,
				HitID:      bestTarget.instanceID,
				HitX:       bestTarget.x,
				HitY:       bestTarget.y,
				IsKill:     isKill,
				Reward:     reward,
				Multiplier: bestTarget.multiplier,
				KillerID:   p.ID,
			},
		})

		if isKill {
			// 廣播目標擊破
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgTargetKill,
				Payload: ws.TargetKillPayload{
					InstanceID: bestTarget.instanceID,
					KillerID:   p.ID,
					Reward:     reward,
					Multiplier: bestTarget.multiplier,
				},
			})
		}

		log.Printf("[SeaAnemone] dir=%d angle=%.0f target=%s kill=%v reward=%d",
			dir, angle, bestTarget.instanceID, isKill, reward)
	}

	if totalReward > 0 {
		p.AddReward(totalReward)
	}

	// 廣播觸手攻擊結果（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgSeaAnemone,
		Payload: ws.SeaAnemonePayload{
			Phase:       "tentacle_result",
			TriggerX:    triggerX,
			TriggerY:    triggerY,
			HitEntries:  hitEntries,
			KillCount:   killCount,
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
				Source:     "sea_anemone",
				Amount:     totalReward,
				Multiplier: float64(killCount),
				NewBalance: p.Coins,
			},
		})
	}

	// 全服公告：擊破 ≥4 個目標
	if killCount >= SeaAnemoneAnnounceThreshold {
		g.announceSeaAnemone(p.DisplayName, killCount, totalReward)
	}

	log.Printf("[SeaAnemone] player=%s kills=%d total_reward=%d",
		p.ID, killCount, totalReward)
}

// announceSeaAnemone 全服公告海葵觸手攻擊（DAY-174）
func (g *Game) announceSeaAnemone(playerName string, killCount int, reward int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "sea_anemone",
			"message":    fmt.Sprintf("🪸 %s 的海葵觸手攻擊擊破 %d 個目標！獲得 %d 金幣！", playerName, killCount, reward),
			"color":      "#FF69B4",
			"duration":   4.0,
			"priority":   2,
		},
	})
}
