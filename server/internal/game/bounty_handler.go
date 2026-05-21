// bounty_handler.go — 全服目標懸賞系統 handler（DAY-137）
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// handlePostBounty 玩家下懸賞（Client → Server）
func (g *Game) handlePostBounty(p *player.Player, msg *ws.Message) {
	if g.Bounty == nil {
		return
	}

	var payload ws.PostBountyPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	// 驗證目標是否存在
	g.mu.RLock()
	t, exists := g.Targets[payload.TargetInstanceID]
	g.mu.RUnlock()
	if !exists {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgBountyError,
			Payload: ws.BountyErrorPayload{Code: "target_not_found", Message: "目標不存在"},
		})
		return
	}

	// 驗證金額（玩家是否有足夠金幣）
	// 注意：DeductCoins 會做最終的金幣檢查，這裡只做快速預檢
	if p.Coins < payload.Amount {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgBountyError,
			Payload: ws.BountyErrorPayload{Code: "insufficient_coins", Message: "金幣不足"},
		})
		return
	}

	// 嘗試下懸賞
	bountyID, errCode := g.Bounty.PostBounty(
		p.ID, p.DisplayName,
		t.InstanceID, t.DefID, t.Def.Name, t.Multiplier,
		payload.Amount,
	)

	if errCode != "" {
		var msg string
		switch errCode {
		case "cooldown":
			cooldown := g.Bounty.GetPlayerCooldown(p.ID)
			msg = fmt.Sprintf("下懸賞冷卻中，還需等待 %d 秒", cooldown)
		case "full":
			msg = "目前懸賞已滿（最多 3 個），請等待現有懸賞結束"
		case "invalid_amount":
			msg = fmt.Sprintf("懸賞金額需在 %d-%d 之間", 100, 5000)
		default:
			msg = "下懸賞失敗"
		}
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgBountyError,
			Payload: ws.BountyErrorPayload{Code: errCode, Message: msg},
		})
		return
	}

	// 扣除金幣
	if _, ok := p.DeductCoins(payload.Amount); !ok {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgBountyError,
			Payload: ws.BountyErrorPayload{Code: "insufficient_coins", Message: "金幣不足"},
		})
		return
	}
	g.sendPlayerUpdate(p)

	log.Printf("[Bounty] Player %s posted bounty %s on %s (×%.0f) for %d coins",
		p.DisplayName, bountyID, t.Def.Name, t.Multiplier, payload.Amount)

	// 全服廣播懸賞發布
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBountyPosted,
		Payload: ws.BountyPostedPayload{
			BountyID:         bountyID,
			TargetInstanceID: t.InstanceID,
			TargetDefID:      t.DefID,
			TargetName:       t.Def.Name,
			TargetMult:       t.Multiplier,
			PosterID:         p.ID,
			PosterName:       p.DisplayName,
			Amount:           payload.Amount,
			SecondsLeft:      60.0,
			Message:          fmt.Sprintf("💰 %s 對【%s】(×%.0f) 懸賞 %d 金幣！", p.DisplayName, t.Def.Name, t.Multiplier, payload.Amount),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventBountyPosted, p.DisplayName, payload.Amount, map[string]string{
		"target_name": t.Def.Name,
		"mult":        fmt.Sprintf("%.0f", t.Multiplier),
	})
	g.broadcastAnnouncement(ann)
}

// notifyBountyKill 擊破懸賞目標，發放懸賞（由 handleKill 呼叫）
// 回傳懸賞總金額（0 = 無懸賞）
func (g *Game) notifyBountyKill(p *player.Player, instanceID string) int {
	if g.Bounty == nil {
		return 0
	}

	totalAmount, claimed, isAny := g.Bounty.ClaimBounty(instanceID, p.ID, p.DisplayName)
	if !isAny || totalAmount == 0 {
		return 0
	}

	// 發放懸賞金幣
	p.AddCoins(totalAmount)
	g.sendPlayerUpdate(p)

	log.Printf("[Bounty] Player %s claimed %d bounty coins for killing %s",
		p.DisplayName, totalAmount, instanceID)

	// 發送個人懸賞領取通知
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgBountyClaimed,
		Payload: ws.BountyClaimedPayload{
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			TotalAmount: totalAmount,
			BountyCount: len(claimed),
			NewBalance:  p.Coins,
			Message:     fmt.Sprintf("💰 你領取了 %d 金幣懸賞！", totalAmount),
		},
	})

	// 全服廣播懸賞被領取
	targetName := ""
	if len(claimed) > 0 {
		targetName = claimed[0].TargetName
	}
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBountyKilled,
		Payload: ws.BountyKilledPayload{
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			TargetName:  targetName,
			TotalAmount: totalAmount,
			BountyCount: len(claimed),
			Message:     fmt.Sprintf("💰 %s 擊破懸賞目標【%s】！獲得 %d 金幣懸賞！", p.DisplayName, targetName, totalAmount),
		},
	})

	// 全服公告（懸賞金額 >= 500 才公告）
	if totalAmount >= 500 {
		ann := g.Announce.Create(announce.EventBountyClaimed, p.DisplayName, totalAmount, map[string]string{
			"target_name": targetName,
		})
		g.broadcastAnnouncement(ann)
	}

	return totalAmount
}

// tickBountyExpiry 懸賞過期檢查（由 gameLoop 每次 update 呼叫）
func (g *Game) tickBountyExpiry() {
	if g.Bounty == nil {
		return
	}

	expired := g.Bounty.CheckExpiry()
	for _, b := range expired {
		// 退款給下懸賞的玩家
		g.mu.RLock()
		poster, exists := g.Players[b.PosterID]
		g.mu.RUnlock()
		if exists {
			poster.AddCoins(b.Amount)
			g.sendPlayerUpdate(poster)
			g.Hub.Send(b.PosterID, &ws.Message{
				Type: ws.MsgBountyExpired,
				Payload: ws.BountyExpiredPayload{
					BountyID:   b.ID,
					TargetName: b.TargetName,
					Amount:     b.Amount,
					Message:    fmt.Sprintf("⏰ 懸賞超時！【%s】的 %d 金幣懸賞已退還。", b.TargetName, b.Amount),
				},
			})
		}

		// 全服廣播懸賞過期
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgBountyExpired,
			Payload: ws.BountyExpiredPayload{
				BountyID:   b.ID,
				TargetName: b.TargetName,
				Amount:     b.Amount,
				Message:    fmt.Sprintf("⏰ 懸賞超時！【%s】的懸賞已取消。", b.TargetName),
			},
		})
		log.Printf("[Bounty] Bounty %s expired: target=%s amount=%d", b.ID, b.TargetName, b.Amount)
	}
}

// cancelBountyForTarget 目標消失時取消懸賞並退款
func (g *Game) cancelBountyForTarget(instanceID string) {
	if g.Bounty == nil {
		return
	}

	cancelled := g.Bounty.CancelBountyForTarget(instanceID)
	for _, b := range cancelled {
		// 退款
		g.mu.RLock()
		poster, exists := g.Players[b.PosterID]
		g.mu.RUnlock()
		if exists {
			poster.AddCoins(b.Amount)
			g.sendPlayerUpdate(poster)
			g.Hub.Send(b.PosterID, &ws.Message{
				Type: ws.MsgBountyExpired,
				Payload: ws.BountyExpiredPayload{
					BountyID:   b.ID,
					TargetName: b.TargetName,
					Amount:     b.Amount,
					Message:    fmt.Sprintf("目標消失！【%s】的 %d 金幣懸賞已退還。", b.TargetName, b.Amount),
				},
			})
		}
	}
}

// sendBountyStatus 登入時發送懸賞狀態
func (g *Game) sendBountyStatus(p *player.Player) {
	if g.Bounty == nil {
		return
	}
	bounties := g.Bounty.GetActiveBounties()
	if len(bounties) == 0 {
		return
	}

	snaps := make([]ws.BountySnap, len(bounties))
	for i, b := range bounties {
		snaps[i] = ws.BountySnap{
			BountyID:         b.ID,
			TargetInstanceID: b.TargetInstanceID,
			TargetDefID:      b.TargetDefID,
			TargetName:       b.TargetName,
			TargetMult:       b.TargetMult,
			PosterID:         b.PosterID,
			PosterName:       b.PosterName,
			Amount:           b.Amount,
			SecondsLeft:      b.SecondsLeft,
		}
	}

	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgBountyList,
		Payload: ws.BountyListPayload{
			Bounties:     snaps,
			CooldownLeft: g.Bounty.GetPlayerCooldown(p.ID),
		},
	})
}
