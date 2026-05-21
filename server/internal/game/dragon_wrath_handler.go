// dragon_wrath_handler.go — 龍怒流星雨武器系統 handler（DAY-154）
// 業界依據：royalfishing.co.uk 2026 Dragon Wrath
// 「Once the wrath meter fills, players unleash a massive meteorite attack across the centre screen,
//  simultaneously targeting multiple fish including Immortal Bosses and the ChainLong King.
//  The wrath value converts proportionally to your bet amount, meaning higher stakes generate faster charge rates.」
// 設計：每次射擊累積怒氣值（60次充滿），釋放 5 波流星雨打擊全場，可命中不死 BOSS 和千龍王
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

// notifyDragonWrathShot 每次射擊時累積龍怒怒氣值（由 handleAttack 呼叫）
// 回傳是否剛充滿一發
func (g *Game) notifyDragonWrathShot(p *player.Player) bool {
	if g.SpecialWeapon == nil {
		return false
	}

	chargeUnlocked, newCharges, newProgress := g.SpecialWeapon.RecordShot(p.ID)

	// 發送怒氣值更新（每次射擊都發，讓 Client 更新怒氣條）
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgDragonWrathCharge,
		Payload: ws.DragonWrathChargePayload{
			PlayerID:    p.ID,
			Progress:    newProgress,
			Required:    60,
			Charges:     newCharges,
			MaxCharges:  1,
			JustCharged: chargeUnlocked,
		},
	}); err != nil {
		log.Printf("[DragonWrath] send charge error: %v", err)
	}

	if chargeUnlocked {
		log.Printf("[DragonWrath] player=%s wrath charged! charges=%d", p.ID, newCharges)
		// 充滿時也更新武器面板
		g.sendSpecialWeaponUpdate(p, false)
	}

	return chargeUnlocked
}

// handleDragonWrathFire 使用龍怒流星雨（由 handleUseSpecialWeapon 呼叫）
// 5 波流星雨，每波間隔 300ms，打擊畫面中央區域
// 可以命中所有目標，包括不死 BOSS 和千龍王（T112）
func (g *Game) handleDragonWrathFire(p *player.Player) {
	log.Printf("[DragonWrath] player=%s firing dragon wrath!", p.ID)

	// 廣播龍怒開始（全螢幕紅色龍形動畫）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgDragonWrathResult,
		Payload: ws.DragonWrathResultPayload{
			KillerID:   p.ID,
			KillerName: p.DisplayName,
			Phase:      "wrath_start",
		},
	})

	go g.runDragonWrathMeteors(p)
}

// runDragonWrathMeteors 執行 5 波流星雨（goroutine）
func (g *Game) runDragonWrathMeteors(p *player.Player) {
	const meteorCount = 5
	const meteorInterval = 300 * time.Millisecond

	// 流星落點（分散在畫面中央區域）
	meteorPositions := []struct{ x, y float64 }{
		{640, 360}, // 中央
		{400, 250}, // 左上
		{880, 250}, // 右上
		{400, 470}, // 左下
		{880, 470}, // 右下
	}

	allHitEntries := make([]ws.DragonWrathMeteorEntry, 0, 20)
	totalReward := 0
	immortalHits := 0
	immortalReward := 0

	// 等待開始動畫（0.8 秒）
	time.Sleep(800 * time.Millisecond)

	for i := 0; i < meteorCount; i++ {
		pos := meteorPositions[i]

		// 收集這波流星的命中目標
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
		// 取得不死 BOSS instanceID
		immortalBossID := ""
		if g.ImmortalBoss != nil {
			immortalBossID = g.ImmortalBoss.GetActiveInstanceID()
		}
		g.mu.RUnlock()

		// 計算這波流星的命中範圍（以落點為中心，半徑 200px）
		hitIDs := calcMeteorHitTargets(pos.x, pos.y, targets, immortalBossID)

		// 廣播這波流星（讓 Client 播放流星落下動畫）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgDragonWrathResult,
			Payload: ws.DragonWrathResultPayload{
				KillerID:    p.ID,
				KillerName:  p.DisplayName,
				Phase:       "meteor",
				MeteorIndex: i + 1,
				MeteorX:     pos.x,
				MeteorY:     pos.y,
			},
		})

		// 等待流星落下動畫（200ms）
		time.Sleep(200 * time.Millisecond)

		// 處理命中
		g.mu.Lock()
		for _, id := range hitIDs {
			// 先檢查是否是不死 BOSS
			if id == immortalBossID && immortalBossID != "" {
				// 不死 BOSS：每次命中給獎勵（不死，不擊破）
				mult, reward, ok := g.ImmortalBoss.RecordHit(immortalBossID, p.ID, p.DisplayName, p.BetLevel)
				if ok {
					entry := ws.DragonWrathMeteorEntry{
						InstanceID: immortalBossID,
						DefID:      "immortal_boss",
						Multiplier: mult,
						Reward:     reward,
						Killed:     false, // 不死 BOSS 不會被擊破
						IsImmortal: true,
					}
					allHitEntries = append(allHitEntries, entry)
					totalReward += reward
					immortalHits++
					immortalReward += reward
					p.Coins += reward
					if p.Coins > p.MaxCoins {
						p.MaxCoins = p.Coins
					}
				}
				continue
			}

			t, ok := g.Targets[id]
			if !ok || t.HP <= 0 {
				continue
			}

			// 判斷擊破機率（依目標類型）
			killChance := specialweapon.DragonWrathNormalKillChance
			isBoss := t.DefID == "B001"
			if isBoss {
				killChance = specialweapon.DragonWrathBossKillChance
			} else if t.Multiplier >= 30 {
				killChance = specialweapon.DragonWrathSpecialKillChance
			}

			entry := ws.DragonWrathMeteorEntry{
				InstanceID: t.InstanceID,
				DefID:      t.DefID,
				Multiplier: t.Multiplier,
				IsBoss:     isBoss,
			}

			if rand.Float64() < killChance {
				// 擊破！獎勵 = 倍率 × betLevel × 0.65（流星雨是範圍武器，比直接擊破低）
				reward := int(float64(p.BetLevel) * t.Multiplier * 0.65)
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
			allHitEntries = append(allHitEntries, entry)
		}
		g.mu.Unlock()

		// 等待下一波（300ms 間隔）
		if i < meteorCount-1 {
			time.Sleep(meteorInterval)
		}
	}

	// 廣播最終結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgDragonWrathResult,
		Payload: ws.DragonWrathResultPayload{
			KillerID:       p.ID,
			KillerName:     p.DisplayName,
			Phase:          "result",
			HitTargets:     allHitEntries,
			TotalReward:    totalReward,
			NewBalance:     p.Coins,
			ImmortalHits:   immortalHits,
			ImmortalReward: immortalReward,
		},
	})

	// 全服公告：龍怒大豐收（≥5個擊破或命中不死 BOSS）
	killedCount := 0
	for _, e := range allHitEntries {
		if e.Killed {
			killedCount++
		}
	}
	if killedCount >= 5 || immortalHits >= 2 {
		g.announceDragonWrath(p.DisplayName, killedCount, totalReward, immortalHits)
	}

	log.Printf("[DragonWrath] player=%s meteors=%d killed=%d immortal_hits=%d total_reward=%d",
		p.ID, meteorCount, killedCount, immortalHits, totalReward)
}

// calcMeteorHitTargets 計算流星落點命中的目標（半徑 200px）
func calcMeteorHitTargets(cx, cy float64, targets []specialweapon.TargetPos, immortalBossID string) []string {
	const meteorRadius = 200.0
	var hit []string

	for _, t := range targets {
		dx := t.X - cx
		dy := t.Y - cy
		dist := dx*dx + dy*dy
		if dist <= meteorRadius*meteorRadius {
			hit = append(hit, t.InstanceID)
		}
	}

	// 不死 BOSS 也加入（如果有活躍的，不管位置）
	if immortalBossID != "" {
		hit = append(hit, immortalBossID)
	}

	return hit
}

// announceDragonWrath 全服公告龍怒大豐收
func (g *Game) announceDragonWrath(playerName string, killedCount int, totalReward int, immortalHits int) {
	var msg string
	if immortalHits >= 2 {
		msg = fmt.Sprintf("🐉 %s 的龍怒流星雨命中不死BOSS %d 次，獲得 %d 金幣！", playerName, immortalHits, totalReward)
	} else {
		msg = fmt.Sprintf("🐉 %s 的龍怒流星雨擊破 %d 個目標，獲得 %d 金幣！", playerName, killedCount, totalReward)
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "dragon_wrath",
			"message":    msg,
			"color":      "#FF4500",
			"duration":   4.0,
			"priority":   3,
		},
	})
}
