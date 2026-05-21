// specialweapon_handler.go — 特殊武器系統 handler（DAY-089，升級 DAY-134，DAY-141）
// 業界依據：
//   - Fish Road 2026 有 8 tier 武器系統，炸彈/雷射是標配特殊武器
//   - Royal Fishing 2026 Tornado Cannon — 龍捲風掃場，旋轉吸入所有目標
//   - JILI 2026 Auto-Charge — 每次擊破目標自動累積充能，不需要花金幣
//   - thechipotlemenu.com 2026 Automatic Target Locking Weapon — AI 自動追蹤最高倍率目標，100% 命中
package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"

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
		// 金幣不足、已達上限、或不可購買（龍捲風砲）
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "buy_weapon_failed", Message: "購買失敗（金幣不足、已達上限或此武器不可購買）"},
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
	case specialweapon.WeaponTornado:
		// 龍捲風：全螢幕，50% 機率擊破每個目標（DAY-134）
		hitIDs = specialweapon.CalcTornadoTargets(targets)
	case specialweapon.WeaponHoming:
		// 追蹤飛彈：自動追蹤倍率最高的目標（DAY-141）
		bestID := specialweapon.CalcHomingTarget(targets)
		if bestID != "" {
			hitIDs = []string{bestID}
		}
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

		switch wtype {
		case specialweapon.WeaponFreeze:
			// 冰凍：不擊破，只減速（Client 端處理視覺）
			entry.Killed = false
			entry.Reward = 0

		case specialweapon.WeaponTornado:
			// 龍捲風：50% 機率擊破（DAY-134）
			if rand.Float64() < specialweapon.TornadoKillChance {
				entry.Killed = true
				// 龍捲風獎勵 = 基礎獎勵 × 0.6（比炸彈稍高，因為是全螢幕）
				baseReward := int(float64(p.BetLevel) * t.Multiplier * 0.6)
				if baseReward < 1 {
					baseReward = 1
				}
				entry.Reward = baseReward
				totalReward += baseReward
				t.HP = 0
				delete(g.Targets, id)
			}

		case specialweapon.WeaponHoming:
			// 追蹤飛彈：100% 命中，獎勵 ×1.5（DAY-141）
			entry.Killed = true
			baseReward := int(float64(p.BetLevel) * t.Multiplier)
			if baseReward < 1 {
				baseReward = 1
			}
			finalReward := int(float64(baseReward) * specialweapon.HomingRewardMult)
			entry.Reward = finalReward
			totalReward += finalReward
			t.HP = 0
			delete(g.Targets, id)

		default:
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

	// 龍捲風：延遲廣播製造「旋轉掃場」的連續感（DAY-134）
	if wtype == specialweapon.WeaponTornado {
		go g.broadcastTornadoEffect(p, hitEntries, totalReward)
		g.sendSpecialWeaponUpdate(p, false)
		return
	}

	// 追蹤飛彈：廣播鎖定追蹤效果（DAY-141）
	if wtype == specialweapon.WeaponHoming {
		go g.broadcastHomingMissileEffect(p, hitEntries, totalReward)
		g.sendSpecialWeaponUpdate(p, false)
		return
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
			NewBalance:  p.Coins,
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

// broadcastTornadoEffect 龍捲風效果：分批廣播擊破，製造旋轉掃場的連續感（DAY-134）
func (g *Game) broadcastTornadoEffect(p *player.Player, hitEntries []ws.SpecialWeaponHitEntry, totalReward int) {
	killedCount := countKilled(hitEntries)

	// 先廣播龍捲風開始（全螢幕旋轉動畫）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgSpecialWeaponFired,
		Payload: ws.SpecialWeaponFiredPayload{
			PlayerID:    p.ID,
			WeaponType:  string(specialweapon.WeaponTornado),
			ClickX:      0,
			ClickY:      0,
			HitTargets:  nil, // 先不帶目標，讓 Client 播放旋轉動畫
			TotalReward: 0,
			NewBalance:  0,
			FreezeTime:  0,
		},
	})

	// 等待旋轉動畫開始（0.5 秒）
	time.Sleep(500 * time.Millisecond)

	// 分批廣播擊破（每 80ms 一批，製造連續掃場感）
	batchSize := 3
	for i := 0; i < len(hitEntries); i += batchSize {
		end := i + batchSize
		if end > len(hitEntries) {
			end = len(hitEntries)
		}
		batch := hitEntries[i:end]

		for _, entry := range batch {
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
		time.Sleep(80 * time.Millisecond)
	}

	// 最後廣播龍捲風結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgSpecialWeaponFired,
		Payload: ws.SpecialWeaponFiredPayload{
			PlayerID:    p.ID,
			WeaponType:  string(specialweapon.WeaponTornado),
			ClickX:      -1, // -1 表示結果廣播
			ClickY:      -1,
			HitTargets:  hitEntries,
			TotalReward: totalReward,
			NewBalance:  p.Coins,
			FreezeTime:  0,
		},
	})

	log.Printf("[SpecialWeapon] player=%s tornado hit=%d killed=%d reward=%d",
		p.ID, len(hitEntries), killedCount, totalReward)
}

// notifySpecialWeaponCharge 擊破目標後累積充能，若充滿則通知玩家（DAY-134）
// 由 handleKill 呼叫
func (g *Game) notifySpecialWeaponCharge(p *player.Player, multiplier float64) {
	if g.SpecialWeapon == nil {
		return
	}

	results := g.SpecialWeapon.RecordKill(p.ID, multiplier)

	// 檢查是否有武器充滿
	hasNewCharge := false
	for _, r := range results {
		if r.ChargeUnlocked {
			hasNewCharge = true
			// 找到武器定義取得名稱和圖示
			weaponName, weaponIcon := getWeaponNameIcon(r.WeaponType)

			// 通知玩家充能完成
			if err := g.Hub.Send(p.ID, &ws.Message{
				Type: ws.MsgSpecialWeaponCharged,
				Payload: ws.SpecialWeaponChargedPayload{
					PlayerID:   p.ID,
					WeaponType: string(r.WeaponType),
					WeaponName: weaponName,
					WeaponIcon: weaponIcon,
					NewCharges: r.NewCharges,
					Message:    fmt.Sprintf("%s %s 充能完成！可以使用了！", weaponIcon, weaponName),
				},
			}); err != nil {
				log.Printf("[SpecialWeapon] send charge notify error: %v", err)
			}

			log.Printf("[SpecialWeapon] player=%s %s charged! charges=%d", p.ID, r.WeaponType, r.NewCharges)
		}
	}

	// 如果有新充能，更新武器狀態
	if hasNewCharge {
		g.sendSpecialWeaponUpdate(p, false)
	}
}

// handleGetSpecialWeapons 查詢特殊武器狀態（Client → Server）
func (g *Game) handleGetSpecialWeapons(p *player.Player) {
	g.sendSpecialWeaponUpdate(p, false)
}

// sendSpecialWeaponUpdate 發送特殊武器狀態更新給玩家
func (g *Game) sendSpecialWeaponUpdate(p *player.Player, withDefs bool) {
	snap := g.SpecialWeapon.GetSnapshot(p.ID)

	payload := ws.SpecialWeaponUpdatePayload{
		PlayerID:              p.ID,
		BombCharges:           snap.BombCharges,
		LaserCharges:          snap.LaserCharges,
		FreezeCharges:         snap.FreezeCharges,
		TornadoCharges:        snap.TornadoCharges,
		HomingCharges:         snap.HomingCharges,
		NewBalance:            p.Coins,
		BombChargeProgress:    snap.BombChargeProgress,
		LaserChargeProgress:   snap.LaserChargeProgress,
		FreezeChargeProgress:  snap.FreezeChargeProgress,
		TornadoChargeProgress: snap.TornadoChargeProgress,
		HomingChargeProgress:  snap.HomingChargeProgress,
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

// getWeaponNameIcon 取得武器名稱和圖示
func getWeaponNameIcon(wtype specialweapon.WeaponType) (name, icon string) {
	for _, d := range specialweapon.AvailableWeapons {
		if d.Type == wtype {
			return d.Name, d.Icon
		}
	}
	return string(wtype), "🔫"
}

// broadcastHomingMissileEffect 追蹤飛彈效果廣播（DAY-141）
// 追蹤飛彈自動鎖定倍率最高的目標，100% 命中，獎勵 ×1.5
// 廣播鎖定動畫（0.8s 追蹤飛行）→ 命中爆炸 → 結果通知
func (g *Game) broadcastHomingMissileEffect(p *player.Player, hitEntries []ws.SpecialWeaponHitEntry, totalReward int) {
	if len(hitEntries) == 0 {
		// 無目標，廣播空結果
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgHomingMissileResult,
			Payload: ws.HomingMissileResultPayload{
				PlayerID:    p.ID,
				TargetID:    "",
				Killed:      false,
				FinalReward: 0,
				NewBalance:  p.Coins,
				Message:     "🎯 沒有可追蹤的目標",
			},
		})
		return
	}

	entry := hitEntries[0] // 追蹤飛彈只命中一個目標

	// 先廣播追蹤飛彈發射（讓所有玩家看到飛彈追蹤動畫）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgSpecialWeaponFired,
		Payload: ws.SpecialWeaponFiredPayload{
			PlayerID:    p.ID,
			WeaponType:  string(specialweapon.WeaponHoming),
			ClickX:      0,
			ClickY:      0,
			HitTargets:  nil, // 先不帶目標，讓 Client 播放追蹤動畫
			TotalReward: 0,
			NewBalance:  0,
			FreezeTime:  0,
		},
	})

	// 等待追蹤飛行動畫（0.8 秒）
	time.Sleep(800 * time.Millisecond)

	// 廣播命中爆炸
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

	// 發送個人結果通知
	baseReward := int(float64(entry.Reward) / specialweapon.HomingRewardMult)
	msg := fmt.Sprintf("🎯 追蹤飛彈命中 ×%.0f 目標！獎勵 ×1.5 = %d 金幣", entry.Multiplier, entry.Reward)
	if !entry.Killed {
		msg = "🎯 追蹤飛彈命中但未擊破目標"
	}

	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgHomingMissileResult,
		Payload: ws.HomingMissileResultPayload{
			PlayerID:    p.ID,
			TargetID:    entry.InstanceID,
			DefID:       entry.DefID,
			Multiplier:  entry.Multiplier,
			BaseReward:  baseReward,
			FinalReward: entry.Reward,
			NewBalance:  p.Coins,
			Killed:      entry.Killed,
			Message:     msg,
		},
	})

	// 全服公告：高倍率目標被追蹤飛彈命中（≥20x）
	if entry.Multiplier >= 20 && entry.Killed {
		g.announceHomingMissileHit(p.DisplayName, entry.Multiplier, entry.Reward)
	}

	log.Printf("[HomingMissile] player=%s hit target=%s mult=%.0f reward=%d",
		p.ID, entry.InstanceID, entry.Multiplier, entry.Reward)
}

// announceHomingMissileHit 全服公告追蹤飛彈命中高倍率目標（DAY-141）
func (g *Game) announceHomingMissileHit(playerName string, multiplier float64, reward int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "homing_missile_hit",
			"message":    fmt.Sprintf("🎯 %s 的追蹤飛彈精準命中 ×%.0f 目標，獲得 %d 金幣！", playerName, multiplier, reward),
			"color":      "#FF0080",
			"duration":   3.5,
			"priority":   2,
		},
	})
}
