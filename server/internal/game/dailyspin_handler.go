// Package game — 每日簽到轉盤 handler（DAY-092）
package game

import (
	"log"

	"digital-twin/server/internal/game/dailyspin"
	"digital-twin/server/internal/game/mysterybox"
	"digital-twin/server/internal/game/specialweapon"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// handleGetDailySpin 處理查詢每日轉盤狀態（Client → Server）
func (g *Game) handleGetDailySpin(playerID string) {
	g.mu.RLock()
	_, ok := g.Players[playerID]
	g.mu.RUnlock()
	if !ok {
		return
	}

	snap := g.DailySpin.GetSnapshot(playerID)
	payload := buildDailySpinStatePayload(snap)
	g.Hub.Send(playerID, &ws.Message{
		Type:    ws.MsgDailySpinState,
		Payload: payload,
	})
}

// handleDailySpin 處理執行每日轉盤（Client → Server）
func (g *Game) handleDailySpin(playerID string) {
	g.mu.RLock()
	p, ok := g.Players[playerID]
	g.mu.RUnlock()
	if !ok {
		return
	}

	// 取得玩家登入連續天數
	loginStreak, _ := p.GetLoginInfo()

	// 執行轉盤
	result := g.DailySpin.Spin(playerID, loginStreak)
	if result == nil {
		// 今天已轉過
		g.Hub.Send(playerID, &ws.Message{
			Type: ws.MsgError,
			Payload: ws.ErrorPayload{
				Code:    "already_spun",
				Message: "今天已經轉過了，明天再來！",
			},
		})
		return
	}

	// 發放獎勵
	resultPayload := g.applyDailySpinReward(p, result)

	log.Printf("[DailySpin] Player %s spun (streak=%d, super=%v, slot=%s, amount=%d)",
		playerID, result.LoginStreak, result.IsSuper, result.Slot.Type, result.Slot.Amount)

	// 發送結果
	g.Hub.Send(playerID, &ws.Message{
		Type:    ws.MsgDailySpinResult,
		Payload: resultPayload,
	})

	// 更新玩家狀態
	g.sendPlayerUpdate(p)
}

// applyDailySpinReward 發放每日轉盤獎勵，回傳結果 Payload
func (g *Game) applyDailySpinReward(p *player.Player, result *dailyspin.SpinResult) ws.DailySpinResultPayload {
	slot := result.Slot
	newBalance := p.GetCoins()
	seasonPoints := 0
	multBonus := 0.0
	mysteryBoxRarity := ""

	switch slot.Type {
	case dailyspin.RewardCoins:
		p.AddCoins(slot.Amount)
		newBalance = p.GetCoins()

	case dailyspin.RewardBombCharge:
		if g.SpecialWeapon != nil {
			for i := 0; i < slot.Amount; i++ {
				g.SpecialWeapon.AddCharge(p.ID, specialweapon.WeaponBomb)
			}
		}

	case dailyspin.RewardLaserCharge:
		if g.SpecialWeapon != nil {
			for i := 0; i < slot.Amount; i++ {
				g.SpecialWeapon.AddCharge(p.ID, specialweapon.WeaponLaser)
			}
		}

	case dailyspin.RewardFreezeCharge:
		if g.SpecialWeapon != nil {
			for i := 0; i < slot.Amount; i++ {
				g.SpecialWeapon.AddCharge(p.ID, specialweapon.WeaponFreeze)
			}
		}

	case dailyspin.RewardMysteryBox:
		// 普通轉盤給普通寶箱，超級轉盤給稀有寶箱
		rarity := "common"
		if result.IsSuper {
			rarity = "rare"
		}
		if g.MysteryBox != nil {
			g.MysteryBox.AddBox(p.ID, mysterybox.BoxRarity(rarity))
		}
		mysteryBoxRarity = rarity

	case dailyspin.RewardSeasonPoints:
		if g.Season != nil {
			g.addSeasonPoints(p.ID, slot.Amount)
		}
		seasonPoints = slot.Amount

	case dailyspin.RewardJackpotTicket:
		// Jackpot 券：直接給金幣（簡化實作，每張 = 1000 金幣）
		jackpotValue := slot.Amount * 1000
		p.AddCoins(jackpotValue)
		newBalance = p.GetCoins()

	case dailyspin.RewardMultBonus:
		// 下次攻擊倍率加成（存到 player 的 pending mult）
		multBonus = float64(slot.Amount)
		// 用 MysteryBox 的 pending mult 機制
		if g.MysteryBox != nil {
			g.MysteryBox.SetPendingMultiplier(p.ID, multBonus)
		}
	}

	return ws.DailySpinResultPayload{
		SlotIndex: result.SlotIndex,
		Slot: ws.DailySpinSlotPayload{
			ID:      slot.ID,
			Type:    string(slot.Type),
			Amount:  slot.Amount,
			Label:   slot.Label,
			Icon:    slot.Icon,
			Color:   slot.Color,
			IsSuper: slot.IsSuper,
		},
		IsSuper:          result.IsSuper,
		LoginStreak:      result.LoginStreak,
		NextSpinAt:       result.NextSpinAt,
		NewBalance:       newBalance,
		SeasonPoints:     seasonPoints,
		MultBonus:        multBonus,
		MysteryBoxRarity: mysteryBoxRarity,
	}
}

// buildDailySpinStatePayload 從快照建立狀態 Payload
func buildDailySpinStatePayload(snap map[string]interface{}) ws.DailySpinStatePayload {
	canSpin, _ := snap["can_spin"].(bool)
	isSuper, _ := snap["is_super"].(bool)
	loginStreak, _ := snap["login_streak"].(int)
	totalSpins, _ := snap["total_spins"].(int)
	nextSpinAt, _ := snap["next_spin_at"].(int64)

	normalSlots := convertSlots(snap["normal_slots"])
	superSlots := convertSlots(snap["super_slots"])

	return ws.DailySpinStatePayload{
		CanSpin:     canSpin,
		IsSuper:     isSuper,
		LoginStreak: loginStreak,
		TotalSpins:  totalSpins,
		NextSpinAt:  nextSpinAt,
		NormalSlots: normalSlots,
		SuperSlots:  superSlots,
	}
}

// convertSlots 轉換 dailyspin.Slot 到 ws.DailySpinSlotPayload
func convertSlots(raw interface{}) []ws.DailySpinSlotPayload {
	slots, ok := raw.([]dailyspin.Slot)
	if !ok {
		return nil
	}
	result := make([]ws.DailySpinSlotPayload, len(slots))
	for i, s := range slots {
		result[i] = ws.DailySpinSlotPayload{
			ID:      s.ID,
			Type:    string(s.Type),
			Amount:  s.Amount,
			Label:   s.Label,
			Icon:    s.Icon,
			Color:   s.Color,
			IsSuper: s.IsSuper,
		}
	}
	return result
}
