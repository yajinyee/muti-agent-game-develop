// torpedo_handler.go — 魚雷武器系統 handler（DAY-155）
// 業界依據：jiligames.com 2026 Mega Fishing「With torpedoes and railgun, you can easily catch sea monsters.」
//   + megafishing.click「Special Weapons Railgun (15x stake), Torpedo (6x stake)」
// 設計：費用 6x betLevel，範圍 250px，擊破機率 85%（普通）/65%（特殊）/40%（BOSS）
//   比炸彈（150px/費用固定 500）更大範圍，但費用更高（動態費用）
//   是業界標準的「高費用高回報」武器，讓高 betLevel 玩家有更強的武器選擇
package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"digital-twin/server/internal/game/specialweapon"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// handleTorpedoFire 使用魚雷（由 handleUseSpecialWeapon 呼叫）
// 費用 = 6x betLevel，範圍 250px，擊破機率 85%
func (g *Game) handleTorpedoFire(p *player.Player, clickX, clickY float64) {
	// 計算費用（6x betLevel）
	cost := p.BetLevel * specialweapon.TorpedoCostMultiplier
	if cost <= 0 {
		cost = p.BetLevel * 6
	}

	// 扣除費用
	g.mu.Lock()
	if p.Coins < cost {
		g.mu.Unlock()
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "insufficient_coins", Message: fmt.Sprintf("魚雷費用 %d 金幣不足", cost)},
		})
		return
	}
	p.Coins -= cost
	g.mu.Unlock()

	log.Printf("[Torpedo] player=%s fired torpedo at (%.0f, %.0f), cost=%d", p.ID, clickX, clickY, cost)

	// 廣播魚雷發射（讓 Client 播放魚雷飛行動畫）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTorpedoResult,
		Payload: ws.TorpedoResultPayload{
			ShooterID:   p.ID,
			ShooterName: p.DisplayName,
			Phase:       "torpedo_launch",
			TargetX:     clickX,
			TargetY:     clickY,
			Cost:        cost,
		},
	})

	// 等待魚雷飛行動畫（0.6 秒）
	time.Sleep(600 * time.Millisecond)

	// 收集命中目標
	g.mu.RLock()
	targets := make([]specialweapon.TargetPos, 0, len(g.Targets))
	for _, t := range g.Targets {
		if t.HP > 0 {
			targets = append(targets, specialweapon.TargetPos{
				InstanceID: t.InstanceID,
				X:          t.X,
				Y:          t.Y,
				Multiplier: t.Multiplier,
			})
		}
	}
	g.mu.RUnlock()

	hitIDs := specialweapon.CalcTorpedoTargets(clickX, clickY, targets)

	// 廣播爆炸（讓 Client 播放爆炸動畫）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTorpedoResult,
		Payload: ws.TorpedoResultPayload{
			ShooterID:   p.ID,
			ShooterName: p.DisplayName,
			Phase:       "explosion",
			TargetX:     clickX,
			TargetY:     clickY,
		},
	})

	// 等待爆炸動畫（0.3 秒）
	time.Sleep(300 * time.Millisecond)

	// 處理命中
	hitEntries := make([]ws.TorpedoKillEntry, 0, len(hitIDs))
	totalReward := 0

	g.mu.Lock()
	for _, id := range hitIDs {
		t, ok := g.Targets[id]
		if !ok || t.HP <= 0 {
			continue
		}

		// 判斷擊破機率（依目標類型）
		killChance := specialweapon.TorpedoNormalKillChance
		isBoss := t.DefID == "B001"
		if isBoss {
			killChance = specialweapon.TorpedoBossKillChance
		} else if t.Multiplier >= 20 {
			killChance = specialweapon.TorpedoSpecialKillChance
		}

		entry := ws.TorpedoKillEntry{
			InstanceID: t.InstanceID,
			DefID:      t.DefID,
			Multiplier: t.Multiplier,
		}

		if rand.Float64() < killChance {
			// 擊破！獎勵 = 倍率 × betLevel × 0.75（魚雷是範圍武器，比直接擊破略低）
			reward := int(float64(p.BetLevel) * t.Multiplier * 0.75)
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
			delete(g.Targets, id)

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
		hitEntries = append(hitEntries, entry)
	}
	g.mu.Unlock()

	// 廣播最終結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTorpedoResult,
		Payload: ws.TorpedoResultPayload{
			ShooterID:   p.ID,
			ShooterName: p.DisplayName,
			Phase:       "result",
			TargetX:     clickX,
			TargetY:     clickY,
			HitTargets:  hitEntries,
			TotalReward: totalReward,
			NewBalance:  p.Coins,
			Cost:        cost,
		},
	})

	// 全服公告：魚雷大豐收（≥4 個擊破）
	killedCount := 0
	for _, e := range hitEntries {
		if e.Killed {
			killedCount++
		}
	}
	if killedCount >= 4 {
		g.announceTorpedo(p.DisplayName, killedCount, totalReward)
	}

	// 更新玩家狀態
	g.sendPlayerUpdate(p)

	log.Printf("[Torpedo] player=%s hit=%d killed=%d total_reward=%d cost=%d",
		p.ID, len(hitEntries), killedCount, totalReward, cost)
}

// announceTorpedo 全服公告魚雷大豐收
func (g *Game) announceTorpedo(playerName string, killedCount int, totalReward int) {
	msg := fmt.Sprintf("🚀 %s 的魚雷擊破 %d 個目標，獲得 %d 金幣！", playerName, killedCount, totalReward)

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "torpedo",
			"message":    msg,
			"color":      "#FFD700",
			"duration":   4.0,
			"priority":   2,
		},
	})
}
