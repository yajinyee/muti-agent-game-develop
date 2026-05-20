// shop_handler.go — 商店系統 handler（DAY-094）
package game

import (
	"log"

	"digital-twin/server/internal/game/specialweapon"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// handleGetShop 處理玩家查詢商店（DAY-094）
func (g *Game) handleGetShop(p *player.Player) {
	g.sendShopUpdate(p)
}

// sendShopUpdate 發送商店狀態給特定玩家（DAY-094）
func (g *Game) sendShopUpdate(p *player.Player) {
	snap := g.Shop.GetSnapshot()
	purchases := g.Shop.GetPlayerDailyPurchases(p.ID)

	// 轉換商品列表
	items := make([]ws.ShopItem, len(snap.Items))
	for i, item := range snap.Items {
		items[i] = ws.ShopItem{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Type:        string(item.Type),
			Price:       item.Price,
			OrigPrice:   item.OrigPrice,
			Reward: ws.ShopItemReward{
				Coins:        item.Reward.Coins,
				BombCharge:   item.Reward.BombCharge,
				LaserCharge:  item.Reward.LaserCharge,
				FreezeCharge: item.Reward.FreezeCharge,
				AttackMult:   item.Reward.AttackMult,
				SeasonPoints: item.Reward.SeasonPoints,
			},
			Stock:          item.Stock,
			LimitPerDay:    item.LimitPerDay,
			IsFlashSale:    item.IsFlashSale,
			FlashEndAt:     item.FlashEndAt,
			PurchasedToday: purchases[item.ID],
		}
	}

	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgShopUpdate,
		Payload: ws.ShopUpdatePayload{
			Items:          items,
			FlashSaleEndAt: snap.FlashSaleEndAt,
			SecondsLeft:    snap.SecondsLeft,
		},
	}); err != nil {
		log.Printf("[Shop] send update error: %v", err)
	}
}

// handleBuyShopItem 處理玩家購買商品（DAY-094）
func (g *Game) handleBuyShopItem(p *player.Player, msg *ws.Message) {
	payload, ok := msg.Payload.(ws.BuyShopItemPayload)
	if !ok {
		log.Printf("[Shop] invalid payload type: %T", msg.Payload)
		return
	}

	result := g.Shop.BuyItem(p.ID, payload.ItemID, p.Coins)
	if !result.Success {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgShopError,
			Payload: ws.ShopErrorPayload{
				ItemID: payload.ItemID,
				Reason: result.Reason,
			},
		})
		return
	}

	// 扣除金幣
	if result.Item.Price > 0 {
		p.Coins -= result.Item.Price
	}

	// 發放獎勵
	reward := result.Reward
	if reward.Coins > 0 {
		p.AddCoins(reward.Coins)
	}
	if reward.BombCharge > 0 {
		for i := 0; i < reward.BombCharge; i++ {
			g.SpecialWeapon.AddCharge(p.ID, specialweapon.WeaponBomb)
		}
	}
	if reward.LaserCharge > 0 {
		for i := 0; i < reward.LaserCharge; i++ {
			g.SpecialWeapon.AddCharge(p.ID, specialweapon.WeaponLaser)
		}
	}
	if reward.FreezeCharge > 0 {
		for i := 0; i < reward.FreezeCharge; i++ {
			g.SpecialWeapon.AddCharge(p.ID, specialweapon.WeaponFreeze)
		}
	}
	if reward.AttackMult > 0 {
		// 攻擊倍率加成：複用 MysteryBox 的 PendingMultiplier 機制
		g.MysteryBox.SetPendingMultiplier(p.ID, reward.AttackMult)
		log.Printf("[Shop] player=%s got attack mult bonus x%.1f", p.ID, reward.AttackMult)
	}
	if reward.SeasonPoints > 0 {
		newLevels := g.addSeasonPoints(p.ID, reward.SeasonPoints)
		g.checkSeasonLevelNotify(p, newLevels)
	}

	log.Printf("[Shop] player=%s bought item=%s price=%d", p.ID, result.Item.ID, result.Item.Price)

	// 發送購買成功通知
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgShopPurchased,
		Payload: ws.ShopPurchasedPayload{
			ItemID:     result.Item.ID,
			ItemName:   result.Item.Name,
			Price:      result.Item.Price,
			NewBalance: p.Coins,
			Reward: ws.ShopItemReward{
				Coins:        reward.Coins,
				BombCharge:   reward.BombCharge,
				LaserCharge:  reward.LaserCharge,
				FreezeCharge: reward.FreezeCharge,
				AttackMult:   reward.AttackMult,
				SeasonPoints: reward.SeasonPoints,
			},
		},
	})

	// 更新商店狀態（讓玩家看到最新庫存和購買次數）
	g.sendShopUpdate(p)

	// 如果有特殊武器充能，更新特殊武器狀態
	if reward.BombCharge > 0 || reward.LaserCharge > 0 || reward.FreezeCharge > 0 {
		g.sendSpecialWeaponUpdate(p, false)
	}
}

// GetShopSnapshot 取得商店快照（供 HTTP 端點使用，DAY-094）
func (g *Game) GetShopSnapshot() interface{} {
	return g.Shop.GetSnapshot()
}
