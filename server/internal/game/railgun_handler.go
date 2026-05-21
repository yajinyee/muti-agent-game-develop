// railgun_handler.go — 軌道炮武器系統 handler（DAY-157）
// 業界依據：megafishing.click 2026「Special Weapons Railgun (15x stake), Torpedo (6x stake)」
//   + jiligames.com 2026「With torpedoes and railgun, you can easily catch sea monsters.」
// 設計：費用 15x betLevel（比魚雷 6x 更貴），穿透全場高能光束（Y 軸 ±40px）
//   普通目標 100% 擊破，特殊目標 90%，BOSS 60%
//   是「終極清場武器」，費用極高但效果最強，讓高 betLevel 玩家有最強的武器選擇
//   充能：擊破 40 個目標自動充能一發（比魚雷 25 更難充能，保持稀有感）
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

// handleRailgunFire 使用軌道炮（由 handleUseSpecialWeapon 呼叫）
// 費用 = 15x betLevel，穿透全場，100% 擊破普通目標
func (g *Game) handleRailgunFire(p *player.Player, clickY float64) {
	// 計算費用（15x betLevel）
	cost := p.BetLevel * specialweapon.RailgunCostMultiplier

	// 扣除費用
	g.mu.Lock()
	if p.Coins < cost {
		g.mu.Unlock()
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "insufficient_coins", Message: fmt.Sprintf("軌道炮費用 %d 金幣不足", cost)},
		})
		return
	}
	p.Coins -= cost
	g.mu.Unlock()

	log.Printf("[Railgun] player=%s fired railgun at Y=%.0f, cost=%d", p.ID, clickY, cost)

	// 廣播軌道炮充能（讓 Client 播放充能動畫）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRailgunResult,
		Payload: ws.RailgunResultPayload{
			ShooterID:   p.ID,
			ShooterName: p.DisplayName,
			Phase:       "railgun_charge",
			TargetY:     clickY,
			Cost:        cost,
		},
	})

	// 等待充能動畫（0.8 秒，比魚雷更長，製造「蓄力」感）
	time.Sleep(800 * time.Millisecond)

	// 廣播軌道炮發射（讓 Client 播放光束穿透動畫）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRailgunResult,
		Payload: ws.RailgunResultPayload{
			ShooterID:   p.ID,
			ShooterName: p.DisplayName,
			Phase:       "railgun_fire",
			TargetY:     clickY,
		},
	})

	// 等待光束動畫（0.4 秒）
	time.Sleep(400 * time.Millisecond)

	// 收集命中目標（包含 BOSS，軌道炮可以命中 BOSS）
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

	hitIDs := specialweapon.CalcRailgunTargets(clickY, targets)

	// 處理命中（按 X 座標排序，製造「從左到右穿透」的視覺感）
	hitEntries := make([]ws.RailgunKillEntry, 0, len(hitIDs))
	totalReward := 0

	g.mu.Lock()
	for _, id := range hitIDs {
		t, ok := g.Targets[id]
		if !ok || t.HP <= 0 {
			continue
		}

		// 判斷擊破機率（依目標類型）
		killChance := specialweapon.RailgunNormalKillChance
		isBoss := t.DefID == "B001"
		if isBoss {
			killChance = specialweapon.RailgunBossKillChance
		} else if t.Multiplier >= 20 {
			killChance = specialweapon.RailgunSpecialKillChance
		}

		entry := ws.RailgunKillEntry{
			InstanceID: t.InstanceID,
			DefID:      t.DefID,
			X:          t.X,
			Multiplier: t.Multiplier,
		}

		if rand.Float64() < killChance {
			// 擊破！獎勵 = 倍率 × betLevel × 0.80（軌道炮是穿透武器，比直接擊破略低）
			reward := int(float64(p.BetLevel) * t.Multiplier * 0.80)
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
		Type: ws.MsgRailgunResult,
		Payload: ws.RailgunResultPayload{
			ShooterID:   p.ID,
			ShooterName: p.DisplayName,
			Phase:       "result",
			TargetY:     clickY,
			HitTargets:  hitEntries,
			TotalReward: totalReward,
			NewBalance:  p.Coins,
			Cost:        cost,
		},
	})

	// 全服公告：軌道炮大豐收（≥3 個擊破，比魚雷 4 個門檻低，因為軌道炮更難充能）
	killedCount := 0
	for _, e := range hitEntries {
		if e.Killed {
			killedCount++
		}
	}
	if killedCount >= 3 {
		g.announceRailgun(p.DisplayName, killedCount, totalReward)
	}

	// 更新玩家狀態
	g.sendPlayerUpdate(p)

	log.Printf("[Railgun] player=%s hit=%d killed=%d total_reward=%d cost=%d",
		p.ID, len(hitEntries), killedCount, totalReward, cost)
}

// announceRailgun 全服公告軌道炮大豐收
func (g *Game) announceRailgun(playerName string, killedCount int, totalReward int) {
	msg := fmt.Sprintf("🔫 %s 的軌道炮穿透擊破 %d 個目標，獲得 %d 金幣！", playerName, killedCount, totalReward)

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "railgun",
			"message":    msg,
			"color":      "#00FFFF",
			"duration":   4.0,
			"priority":   2,
		},
	})
}
