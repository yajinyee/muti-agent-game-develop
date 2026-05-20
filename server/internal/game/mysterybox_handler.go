// mysterybox_handler.go — 神秘寶箱系統 handler（DAY-090）
// 業界依據：nerdbot.com 2026-05-02 確認「mystery rewards」是 2026 年 iGaming 最熱門留存機制
package game

import (
	"log"

	"digital-twin/server/internal/game/mysterybox"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// notifyMysteryBoxKill 在擊破目標後嘗試掉落神秘寶箱（由 handleKill 呼叫）
func (g *Game) notifyMysteryBoxKill(p *player.Player, targetX, targetY float64, isBoss bool) {
	box := g.MysteryBox.TryDropBox(isBoss)
	if box == nil {
		return
	}

	// 加入玩家背包
	g.MysteryBox.AddBox(p.ID, box.Rarity)

	log.Printf("[MysteryBox] player=%s dropped %s box at (%.0f, %.0f)", p.ID, box.Rarity, targetX, targetY)

	// 廣播掉落通知（所有玩家可見，增加驚喜感）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgMysteryBoxDrop,
		Payload: ws.MysteryBoxDropPayload{
			PlayerID:  p.ID,
			Rarity:    string(box.Rarity),
			Name:      box.Name,
			Icon:      box.Icon,
			Color:     box.Color,
			GlowColor: box.GlowColor,
			DropX:     targetX,
			DropY:     targetY,
		},
	})

	// 更新玩家背包狀態
	g.sendMysteryBoxUpdate(p)
}

// handleOpenMysteryBox 開箱請求（Client → Server）
func (g *Game) handleOpenMysteryBox(p *player.Player, msg *ws.Message) {
	var payload ws.OpenMysteryBoxPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	rarity := mysterybox.BoxRarity(payload.Rarity)

	// 確認玩家有這個稀有度的寶箱
	if !g.MysteryBox.HasBox(p.ID, rarity) {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "no_box", Message: "沒有這個稀有度的寶箱"},
		})
		return
	}

	// 消耗寶箱
	g.MysteryBox.RemoveBox(p.ID, rarity)

	// 開箱獲得獎勵
	reward := g.MysteryBox.OpenBox(rarity)
	if reward == nil {
		return
	}

	// 發放獎勵
	pendingMult := 0.0
	switch reward.Type {
	case mysterybox.RewardCoins:
		p.Coins += reward.Amount
		if p.Coins > p.MaxCoins {
			p.MaxCoins = p.Coins
		}
	case mysterybox.RewardBombCharge:
		for i := 0; i < reward.Amount; i++ {
			g.SpecialWeapon.AddCharge(p.ID, "bomb")
		}
		g.sendSpecialWeaponUpdate(p, false)
	case mysterybox.RewardLaserCharge:
		for i := 0; i < reward.Amount; i++ {
			g.SpecialWeapon.AddCharge(p.ID, "laser")
		}
		g.sendSpecialWeaponUpdate(p, false)
	case mysterybox.RewardFreezeCharge:
		for i := 0; i < reward.Amount; i++ {
			g.SpecialWeapon.AddCharge(p.ID, "freeze")
		}
		g.sendSpecialWeaponUpdate(p, false)
	case mysterybox.RewardMultiplier:
		// amount = 倍率 × 10（例如 20 = 2.0x，50 = 5.0x）
		mult := float64(reward.Amount) / 10.0
		g.MysteryBox.SetPendingMultiplier(p.ID, mult)
		pendingMult = mult
	case mysterybox.RewardJackpotTicket:
		// Jackpot 券：增加 Jackpot 貢獻
		g.jackpotMgr.Contribute(reward.Amount*100, p.ID) // 每張券 = 100 金幣貢獻
	}

	// 取得寶箱定義
	boxDef := mysterybox.GetBoxDef(rarity)
	boxName := ""
	boxIcon := ""
	if boxDef != nil {
		boxName = boxDef.Name
		boxIcon = boxDef.Icon
	}

	// 取得剩餘數量
	remaining := g.MysteryBox.GetBoxCount(p.ID, rarity)

	log.Printf("[MysteryBox] player=%s opened %s box, reward=%s amount=%d", p.ID, rarity, reward.Type, reward.Amount)

	// 發送開箱結果
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgMysteryBoxOpened,
		Payload: ws.MysteryBoxOpenedPayload{
			PlayerID: p.ID,
			Rarity:   string(rarity),
			BoxName:  boxName,
			BoxIcon:  boxIcon,
			Reward: ws.MysteryBoxRewardPayload{
				Type:   string(reward.Type),
				Amount: reward.Amount,
				Label:  reward.Label,
				Icon:   reward.Icon,
				Color:  reward.Color,
			},
			NewBalance:      p.Coins,
			PendingMultMult: pendingMult,
			RemainingBoxes:  remaining,
		},
	})

	// 更新背包狀態
	g.sendMysteryBoxUpdate(p)
}

// handleGetMysteryBoxes 查詢持有寶箱（Client → Server）
func (g *Game) handleGetMysteryBoxes(p *player.Player) {
	g.sendMysteryBoxUpdate(p)
}

// sendMysteryBoxUpdate 發送持有寶箱狀態更新給玩家
func (g *Game) sendMysteryBoxUpdate(p *player.Player) {
	inventory := g.MysteryBox.GetInventory(p.ID)
	entries := make([]ws.MysteryBoxInventoryEntry, 0, len(inventory))
	total := 0

	for _, box := range mysterybox.AvailableBoxes {
		count := inventory[box.Rarity]
		if count > 0 {
			entries = append(entries, ws.MysteryBoxInventoryEntry{
				Rarity:    string(box.Rarity),
				Name:      box.Name,
				Icon:      box.Icon,
				Color:     box.Color,
				GlowColor: box.GlowColor,
				Count:     count,
			})
			total += count
		}
	}

	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgMysteryBoxUpdate,
		Payload: ws.MysteryBoxUpdatePayload{
			PlayerID:  p.ID,
			Inventory: entries,
			Total:     total,
		},
	}); err != nil {
		log.Printf("[MysteryBox] send update error: %v", err)
	}
}
