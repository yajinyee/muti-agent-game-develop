// specialweapon_handler.go — 特殊武器系統 handler（DAY-089）
// 業界依據：Fish Road 2026 有 8 tier 武器系統，炸彈/雷射是標配特殊武器
package game

import (
	"log"

	"digital-twin/server/internal/game/combat"
	"digital-twin/server/internal/game/specialweapon"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// handleBuySpecialWeapon 購買特殊武器（Client → Server）
func (g *Game) handleBuySpecialWeapon(p *player.Player, msg *ws.Message) {
	var payload ws.BuySpecialWeaponPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	wtype := specialweapon.WeaponType(payload.WeaponType)
	ok, cost := g.SpecialWeapon.BuyWeapon(p.ID, wtype, p.Coins)
	if !ok {
		// 金幣不足或已達上限
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "buy_weapon_failed", Message: "購買失敗（金幣不足或已達上限）"},
		})
		return
	}

	// 扣除金幣
	p.Coins -= cost
	if p.Coins < 0 {
		p.Coins = 0
	}

	log.Printf("[SpecialWeapon] player=%s bought %s, cost=%d, balance=%d", p.ID, wtype, cost, p.Coins)

	// 發送武器狀態更新
	g.sendSpecialWeaponUpdate(p, true)
}

// handleUseSpecialWeapon 使用特殊武器（Client → Server）
func (g *Game) handleUseSpecialWeapon(p *player.Player, msg *ws.Message) {
	var payload ws.UseSpecialWeaponPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	wtype := specialweapon.WeaponType(payload.WeaponType)
	if !g.SpecialWeapon.UseWeapon(p.ID, wtype) {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "no_weapon_charges", Message: "沒有可用的特殊武器"},
		})
		return
	}

	log.Printf("[SpecialWeapon] player=%s fired %s at (%.0f, %.0f)", p.ID, wtype, payload.ClickX, payload.ClickY)

	// 計算命中目標
	g.mu.RLock()
	targets := make([]specialweapon.TargetPos, 0, len(g.Targets))
	for _, t := range g.Targets {
		if t.HP > 0 && t.DefID != "B001" { // BOSS 不受特殊武器影響（防止 RTP 爆炸）
			targets = append(targets, specialweapon.TargetPos{
				InstanceID: t.InstanceID,
				X:          t.X,
				Y:          t.Y,
				Multiplier: t.Multiplier,
			})
		}
	}
	g.mu.RUnlock()

	var hitIDs []string
	var freezeTime float64

	switch wtype {
	case specialweapon.WeaponBomb:
		hitIDs = specialweapon.CalcBombTargets(payload.ClickX, payload.ClickY, targets)
	case specialweapon.WeaponLaser:
		hitIDs = specialweapon.CalcLaserTargets(payload.ClickY, targets)
	case specialweapon.WeaponFreeze:
		hitIDs = specialweapon.CalcFreezeTargets(targets)
		freezeTime = 5.0 // 冰凍 5 秒
	}

	// 處理命中目標
	hitEntries := make([]ws.SpecialWeaponHitEntry, 0, len(hitIDs))
	totalReward := 0

	g.mu.Lock()
	for _, id := range hitIDs {
		t, ok := g.Targets[id]
		if !ok || t.HP <= 0 {
			continue
		}

		entry := ws.SpecialWeaponHitEntry{
			InstanceID: t.InstanceID,
			DefID:      t.DefID,
			Multiplier: t.Multiplier,
		}

		if wtype == specialweapon.WeaponFreeze {
			// 冰凍：不擊破，只減速（Client 端處理視覺）
			entry.Killed = false
			entry.Reward = 0
		} else {
			// 炸彈/雷射：嘗試擊破（使用基礎命中率）
			req := combat.AttackRequest{
				PlayerID: p.ID,
				TargetID: t.InstanceID,
				BetLevel: p.BetLevel,
				IsAuto:   false,
				IsLock:   false,
			}
			result := combat.ProcessAttack(req, t)
			if result.IsKill {
				entry.Killed = true
				// 計算獎勵（特殊武器獎勵 = 基礎獎勵 × 0.5，避免 RTP 爆炸）
				baseReward := int(float64(p.BetLevel) * t.Multiplier * 0.5)
				if baseReward < 1 {
					baseReward = 1
				}
				entry.Reward = baseReward
				totalReward += baseReward
				// 擊破目標
				t.HP = 0
				delete(g.Targets, id)
			}
		}
		hitEntries = append(hitEntries, entry)
	}
	g.mu.Unlock()

	// 發放獎勵
	if totalReward > 0 {
		p.Coins += totalReward
		if p.Coins > p.MaxCoins {
			p.MaxCoins = p.Coins
		}
	}

	// 廣播特殊武器發射效果給所有玩家
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgSpecialWeaponFired,
		Payload: ws.SpecialWeaponFiredPayload{
			PlayerID:    p.ID,
			WeaponType:  string(wtype),
			ClickX:      payload.ClickX,
			ClickY:      payload.ClickY,
			HitTargets:  hitEntries,
			TotalReward: totalReward,
			NewBalance:  p.Coins, // 只有發射者有值（Client 端只更新自己的餘額）
			FreezeTime:  freezeTime,
		},
	})

	// 更新武器狀態（扣除一發）
	g.sendSpecialWeaponUpdate(p, false)

	// 廣播被擊破的目標
	for _, entry := range hitEntries {
		if entry.Killed {
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgTargetKill,
				Payload: ws.TargetKillPayload{
					InstanceID: entry.InstanceID,
					KillerID:   p.ID,
					Reward:     entry.Reward,
					Multiplier: entry.Multiplier,
				},
			})
		}
	}

	log.Printf("[SpecialWeapon] player=%s %s hit=%d killed=%d reward=%d",
		p.ID, wtype, len(hitEntries), countKilled(hitEntries), totalReward)
}

// handleGetSpecialWeapons 查詢特殊武器狀態（Client → Server）
func (g *Game) handleGetSpecialWeapons(p *player.Player) {
	g.sendSpecialWeaponUpdate(p, false)
}

// sendSpecialWeaponUpdate 發送特殊武器狀態更新給玩家
func (g *Game) sendSpecialWeaponUpdate(p *player.Player, withDefs bool) {
	snap := g.SpecialWeapon.GetSnapshot(p.ID)

	payload := ws.SpecialWeaponUpdatePayload{
		PlayerID:      p.ID,
		BombCharges:   snap.BombCharges,
		LaserCharges:  snap.LaserCharges,
		FreezeCharges: snap.FreezeCharges,
		NewBalance:    p.Coins,
	}

	// 首次發送時附帶武器定義
	if withDefs {
		defs := make([]ws.SpecialWeaponDef, 0, len(specialweapon.AvailableWeapons))
		for _, d := range specialweapon.AvailableWeapons {
			defs = append(defs, ws.SpecialWeaponDef{
				Type:        string(d.Type),
				Name:        d.Name,
				Description: d.Description,
				Cost:        d.Cost,
				MaxCharges:  d.MaxCharges,
				Icon:        d.Icon,
				Color:       d.Color,
			})
		}
		payload.Definitions = defs
	}

	if err := g.Hub.Send(p.ID, &ws.Message{
		Type:    ws.MsgSpecialWeaponUpdate,
		Payload: payload,
	}); err != nil {
		log.Printf("[SpecialWeapon] send update error: %v", err)
	}
}

// countKilled 計算被擊破的目標數
func countKilled(entries []ws.SpecialWeaponHitEntry) int {
	n := 0
	for _, e := range entries {
		if e.Killed {
			n++
		}
	}
	return n
}
