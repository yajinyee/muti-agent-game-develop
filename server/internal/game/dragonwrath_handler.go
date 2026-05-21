// dragonwrath_handler.go — 龍怒蓄力大招系統 handler（DAY-128）
// 業界依據：JILI Royal Fishing 2026 Dragon Wrath — 累積怒氣值釋放全螢幕大招
// 設計：每次射擊 +1 怒氣，擊破目標 +multiplier 怒氣，滿 100 可釋放「吉伊卡哇大討伐」
// 大招效果：全螢幕攻擊所有目標，每個目標 30% 擊破機率，持續 2 秒
package game

import (
	"log"
	"math/rand"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	WrathPerShot       = 1    // 每次射擊 +1
	WrathPerKillBase   = 2    // 每次擊破基礎 +2
	WrathKillMultScale = 0.5  // 擊破倍率加成係數（multiplier × 0.5）
	WrathKillCap       = 10   // 單次擊破最多 +10 怒氣
	WrathHitChance     = 0.30 // 大招命中每個目標的機率 30%
)

// notifyWrathShot 每次射擊時累積怒氣（由 handleAttack 呼叫）
func (g *Game) notifyWrathShot(p *player.Player) {
	newCharge, justFull := p.AddWrathCharge(WrathPerShot)

	// 每 10 點發送一次更新（減少訊息量）
	if newCharge%10 == 0 || justFull {
		g.sendWrathUpdate(p, newCharge, justFull)
	}
}

// notifyWrathKill 擊破目標時累積怒氣（由 handleKill 呼叫）
func (g *Game) notifyWrathKill(p *player.Player, multiplier float64) {
	// 怒氣加成：基礎 +2，倍率加成（最多 +10）
	gain := WrathPerKillBase + int(multiplier*WrathKillMultScale)
	if gain > WrathKillCap {
		gain = WrathKillCap
	}
	newCharge, justFull := p.AddWrathCharge(gain)

	// 發送更新
	g.sendWrathUpdate(p, newCharge, justFull)

	if justFull {
		log.Printf("[DragonWrath] Player %s wrath FULL (100/100)", p.ID)
	}
}

// handleUseWrath 處理玩家釋放大招（Client → Server）
func (g *Game) handleUseWrath(playerID string) {
	g.mu.RLock()
	p, ok := g.Players[playerID]
	g.mu.RUnlock()
	if !ok {
		return
	}

	// 嘗試消耗怒氣（ConsumeWrath 內部做所有檢查）
	if !p.ConsumeWrath() {
		charge := p.GetWrathCharge()
		cooldown := p.GetWrathCooldownSecs()
		if charge < player.WrathMaxCharge {
			g.Hub.Send(playerID, &ws.Message{
				Type: ws.MsgError,
				Payload: ws.ErrorPayload{
					Code:    "wrath_not_ready",
					Message: "怒氣值尚未充滿！",
				},
			})
		} else if cooldown > 0 {
			g.Hub.Send(playerID, &ws.Message{
				Type: ws.MsgError,
				Payload: ws.ErrorPayload{
					Code:    "wrath_cooldown",
					Message: "大招冷卻中！",
				},
			})
		}
		return
	}

	log.Printf("[DragonWrath] Player %s unleashed Dragon Wrath!", playerID)

	// 廣播大招開始（全服可見）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgWrathStart,
		Payload: ws.WrathStartPayload{
			PlayerID:   playerID,
			PlayerName: p.DisplayName,
			Icon:       "🐉",
			Message:    p.DisplayName + " 釋放了吉伊卡哇大討伐！",
		},
	})

	// 更新怒氣值顯示（已清零）
	g.sendWrathUpdate(p, 0, false)

	// 執行大招效果（goroutine，避免阻塞）
	go g.executeWrathAttack(p)
}

// executeWrathAttack 執行大招攻擊效果
func (g *Game) executeWrathAttack(p *player.Player) {
	// 收集當前所有目標
	g.mu.RLock()
	targetIDs := make([]string, 0, len(g.Targets))
	for id := range g.Targets {
		targetIDs = append(targetIDs, id)
	}
	g.mu.RUnlock()

	totalReward := 0
	killedCount := 0
	killedTargets := make([]ws.WrathKillEntry, 0)

	// 取得投注定義
	betDef := p.GetBetDef()
	if betDef == nil {
		return
	}

	// 對每個目標進行攻擊（30% 機率擊破）
	for _, instanceID := range targetIDs {
		if rand.Float64() > WrathHitChance {
			continue // 未命中
		}

		g.mu.Lock()
		t, exists := g.Targets[instanceID]
		if !exists || !t.IsAlive {
			g.mu.Unlock()
			continue
		}

		// 計算獎勵
		reward := int(float64(betDef.BetCost) * t.Multiplier)

		// 擊破目標
		t.IsAlive = false
		t.HP = 0
		delete(g.Targets, instanceID)
		g.mu.Unlock()

		// 發放獎勵
		p.AddCoins(reward)
		totalReward += reward
		killedCount++

		killedTargets = append(killedTargets, ws.WrathKillEntry{
			InstanceID: instanceID,
			DefID:      t.Def.ID,
			Reward:     reward,
			Multiplier: t.Multiplier,
		})

		// 廣播目標消失
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetKill,
			Payload: ws.TargetKillPayload{
				InstanceID: instanceID,
				KillerID:   p.ID,
				Reward:     reward,
				Multiplier: t.Multiplier,
			},
		})

		// 短暫延遲，製造連續擊破的視覺效果
		time.Sleep(50 * time.Millisecond)
	}

	// 廣播大招結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgWrathResult,
		Payload: ws.WrathResultPayload{
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			KilledCount: killedCount,
			TotalReward: totalReward,
			NewBalance:  p.GetCoins(),
			Targets:     killedTargets,
		},
	})

	log.Printf("[DragonWrath] Player %s wrath result: killed=%d, reward=%d",
		p.ID, killedCount, totalReward)

	// 更新玩家狀態
	g.sendPlayerUpdate(p)

	// 動態牆：大招擊破 ≥ 5 個目標時廣播
	if killedCount >= 5 {
		go g.notifyFeedBossKill(p, totalReward)
	}
}

// sendWrathUpdate 發送怒氣值更新給玩家
func (g *Game) sendWrathUpdate(p *player.Player, charge int, isReady bool) {
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgWrathUpdate,
		Payload: ws.WrathUpdatePayload{
			Charge:    charge,
			MaxCharge: player.WrathMaxCharge,
			IsReady:   isReady,
			Cooldown:  0,
		},
	})
}

// sendWrathStatus 登入時發送怒氣值狀態
func (g *Game) sendWrathStatus(p *player.Player) {
	charge := p.GetWrathCharge()
	cooldown := p.GetWrathCooldownSecs()
	isReady := charge >= player.WrathMaxCharge && cooldown == 0

	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgWrathUpdate,
		Payload: ws.WrathUpdatePayload{
			Charge:    charge,
			MaxCharge: player.WrathMaxCharge,
			IsReady:   isReady,
			Cooldown:  cooldown,
		},
	})
}
