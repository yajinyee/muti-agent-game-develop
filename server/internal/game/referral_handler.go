// referral_handler.go — 推薦碼系統 handler（DAY-082）
package game

import (
	"log"

	"digital-twin/server/internal/game/referral"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// sendReferralInfo 發送推薦碼資訊給指定玩家
func (g *Game) sendReferralInfo(p *player.Player) {
	info := g.Referral.GetInfo(p.ID)
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgReferralInfo,
		Payload: ws.ReferralInfoPayload{
			MyCode:         info.MyCode,
			UsedCode:       info.UsedCode,
			ReferredBy:     info.ReferredBy,
			ReferralCount:  info.ReferralCount,
			TotalReward:    info.TotalReward,
			ReferrerReward: referral.ReferrerReward,
			RefereeReward:  referral.RefereeReward,
			MaxReferrals:   referral.MaxReferrals,
		},
	}); err != nil {
		log.Printf("[Referral] sendReferralInfo error: %v", err)
	}
}

// handleGetReferralInfo 處理查詢推薦碼請求
func (g *Game) handleGetReferralInfo(p *player.Player) {
	g.sendReferralInfo(p)
}

// handleUseReferralCode 處理使用推薦碼請求
func (g *Game) handleUseReferralCode(p *player.Player, msg *ws.Message) {
	var payload ws.UseReferralCodePayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	if payload.Code == "" {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgReferralError,
			Payload: ws.ReferralErrorPayload{Code: "", Reason: "推薦碼不能為空"},
		})
		return
	}

	referrerID, err := g.Referral.UseCode(p.ID, payload.Code)
	if err != nil {
		log.Printf("[Referral] player=%s use code=%s failed: %v", p.ID, payload.Code, err)
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgReferralError,
			Payload: ws.ReferralErrorPayload{Code: payload.Code, Reason: err.Error()},
		})
		return
	}

	// 發放被推薦人獎勵
	p.AddCoins(referral.RefereeReward)

	// 通知被推薦人
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgReferralSuccess,
		Payload: ws.ReferralSuccessPayload{
			Code:       payload.Code,
			ReferrerID: referrerID,
			Reward:     referral.RefereeReward,
			NewBalance: p.GetCoins(),
			Message:    "推薦碼使用成功！",
		},
	})

	// 發放推薦人獎勵（若推薦人在線）
	g.mu.RLock()
	referrer, referrerOnline := g.Players[referrerID]
	g.mu.RUnlock()

	if referrerOnline {
		referrer.AddCoins(referral.ReferrerReward)
		// 通知推薦人
		g.Hub.Send(referrerID, &ws.Message{
			Type: ws.MsgReferralSuccess,
			Payload: ws.ReferralSuccessPayload{
				Code:       payload.Code,
				ReferrerID: referrerID,
				Reward:     referral.ReferrerReward,
				NewBalance: referrer.GetCoins(),
				Message:    "你的推薦碼被使用了！",
			},
		})
		log.Printf("[Referral] referrer=%s online, rewarded %d coins", referrerID, referral.ReferrerReward)
	} else {
		// 推薦人不在線，記錄待發放（簡化版：下次登入時補發）
		log.Printf("[Referral] referrer=%s offline, reward pending", referrerID)
	}

	// 更新雙方推薦資訊
	g.sendReferralInfo(p)

	log.Printf("[Referral] player=%s used code=%s (referrer=%s), referee_reward=%d",
		p.ID, payload.Code, referrerID, referral.RefereeReward)
}
