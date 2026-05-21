// bomb_crab_handler.go — 炸彈蟹連帶效果 handler（DAY-143）
// 業界依據：royal-fishing.uk 2026「Worth 70x, this explosive crustacean triggers multiple
// large-scale detonations. Each bomb creates expanding capture zones for massive multi-target eliminations.」
// 擊破 T107 後觸發 3 波爆炸，每波爆炸半徑 150px，每波間隔 400ms，連帶擊破爆炸範圍內所有目標
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
	// BombCrabExplosionRadius 每波爆炸半徑（px）
	BombCrabExplosionRadius = 150.0
	// BombCrabExplosionWaves 爆炸波數
	BombCrabExplosionWaves = 3
	// BombCrabWaveDelayMs 每波爆炸間隔（ms）
	BombCrabWaveDelayMs = 400
	// BombCrabRewardMult 連帶擊破獎勵倍率（比直接擊破低，平衡 RTP）
	BombCrabRewardMult = 0.50
	// BombCrabAnnounceThreshold 全服公告門檻（連帶擊破數）
	BombCrabAnnounceThreshold = 4
)

// isBombCrab 判斷是否為炸彈蟹（T107）
func isBombCrab(defID string) bool {
	return defID == "T107"
}

// bombCrabTarget 炸彈蟹連帶目標（內部結構）
type bombCrabTarget struct {
	instanceID string
	defID      string
	x, y       float64
	multiplier float64
}

// tryBombCrabChain 擊破 T107 後觸發多波爆炸連帶效果（DAY-143）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryBombCrabChain(p *player.Player, triggerID string, triggerX, triggerY float64) {
	totalReward := 0
	var allKilledEntries []ws.BombCrabKillEntry

	// 廣播炸彈蟹爆炸開始（讓 Client 播放爆炸動畫）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBombCrabChain,
		Payload: ws.BombCrabChainPayload{
			TriggerID:  triggerID,
			TriggerX:   triggerX,
			TriggerY:   triggerY,
			Phase:      "bomb_start",
			WaveIndex:  0,
			TotalWaves: BombCrabExplosionWaves,
		},
	})

	// 3 波爆炸
	for wave := 0; wave < BombCrabExplosionWaves; wave++ {
		if wave > 0 {
			time.Sleep(time.Duration(BombCrabWaveDelayMs) * time.Millisecond)
		}

		// 爆炸中心：第一波在觸發點，後續波在隨機偏移位置（製造擴散感）
		explodeX := triggerX
		explodeY := triggerY
		if wave == 1 {
			// 第二波：向右偏移
			explodeX = triggerX + 120.0
			if explodeX > 1200 {
				explodeX = triggerX - 120.0
			}
		} else if wave == 2 {
			// 第三波：向左偏移
			explodeX = triggerX - 120.0
			if explodeX < 80 {
				explodeX = triggerX + 120.0
			}
		}

		// 收集爆炸範圍內的目標
		g.mu.RLock()
		var waveTargets []bombCrabTarget
		for id, t := range g.Targets {
			if id == triggerID || t.HP <= 0 || t.DefID == "B001" {
				continue
			}
			dx := t.X - explodeX
			dy := t.Y - explodeY
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= BombCrabExplosionRadius {
				waveTargets = append(waveTargets, bombCrabTarget{
					instanceID: id,
					defID:      t.DefID,
					x:          t.X,
					y:          t.Y,
					multiplier: t.Multiplier,
				})
			}
		}
		g.mu.RUnlock()

		// 廣播本波爆炸（讓 Client 播放爆炸特效）
		waveIDs := make([]string, 0, len(waveTargets))
		for _, dt := range waveTargets {
			waveIDs = append(waveIDs, dt.instanceID)
		}

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgBombCrabChain,
			Payload: ws.BombCrabChainPayload{
				TriggerID:  triggerID,
				TriggerX:   explodeX,
				TriggerY:   explodeY,
				Phase:      "explosion",
				WaveIndex:  wave,
				TotalWaves: BombCrabExplosionWaves,
				ExplodeIDs: waveIDs,
			},
		})

		time.Sleep(80 * time.Millisecond)

		// 擊破爆炸範圍內的目標
		for _, dt := range waveTargets {
			g.mu.Lock()
			t, ok := g.Targets[dt.instanceID]
			if !ok || t.HP <= 0 {
				g.mu.Unlock()
				continue
			}
			reward := int(float64(p.BetLevel) * dt.multiplier * BombCrabRewardMult)
			if reward < 1 {
				reward = 1
			}
			t.HP = 0
			delete(g.Targets, dt.instanceID)
			g.mu.Unlock()

			totalReward += reward
			allKilledEntries = append(allKilledEntries, ws.BombCrabKillEntry{
				InstanceID: dt.instanceID,
				DefID:      dt.defID,
				Multiplier: dt.multiplier,
				Reward:     reward,
				WaveIndex:  wave,
			})

			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgTargetKill,
				Payload: ws.TargetKillPayload{
					InstanceID: dt.instanceID,
					KillerID:   p.ID,
					Reward:     reward,
					Multiplier: dt.multiplier,
				},
			})

			log.Printf("[BombCrab] wave=%d target=%s mult=%.0f reward=%d",
				wave, dt.instanceID, dt.multiplier, reward)
		}
	}

	if totalReward <= 0 {
		return
	}

	// 發放總獎勵
	p.AddReward(totalReward)

	// 廣播炸彈蟹連帶結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBombCrabChain,
		Payload: ws.BombCrabChainPayload{
			TriggerID:     triggerID,
			TriggerX:      triggerX,
			TriggerY:      triggerY,
			Phase:         "result",
			TotalWaves:    BombCrabExplosionWaves,
			KilledTargets: allKilledEntries,
			TotalReward:   totalReward,
			KillerID:      p.ID,
			KillerName:    p.DisplayName,
		},
	})

	// 個人結果通知
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgReward,
		Payload: ws.RewardPayload{
			Source:     "bomb_crab",
			Amount:     totalReward,
			Multiplier: float64(len(allKilledEntries)),
			NewBalance: p.Coins,
		},
	})

	// 全服公告：連帶擊破 ≥4 個目標
	if len(allKilledEntries) >= BombCrabAnnounceThreshold {
		g.announceBombCrabChain(p.DisplayName, len(allKilledEntries), totalReward)
	}

	log.Printf("[BombCrab] player=%s waves=%d total_killed=%d total_reward=%d",
		p.ID, BombCrabExplosionWaves, len(allKilledEntries), totalReward)
}

// announceBombCrabChain 全服公告炸彈蟹連帶效果（DAY-143）
func (g *Game) announceBombCrabChain(playerName string, killCount int, reward int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "bomb_crab_chain",
			"message":    fmt.Sprintf("💣 %s 的炸彈蟹連環爆炸擊破 %d 個目標！獲得 %d 金幣！", playerName, killCount, reward),
			"color":      "#FF4500",
			"duration":   4.5,
			"priority":   2,
		},
	})
}
