// dm_handler.go — 玩家私訊系統 handler（DAY-103）
// 好友間可以互相發送私訊，離線訊息暫存，上線後自動發送
package game

import (
	"log"

	"digital-twin/server/internal/game/dm"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// handleSendDM 處理發送私訊（DAY-103）
func (g *Game) handleSendDM(p *player.Player, msg *ws.Message) {
	var payload ws.SendDMPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	if payload.ToID == "" {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgDMError,
			Payload: ws.DMErrorPayload{ErrorCode: "invalid_target", Message: "無效的接收者"},
		})
		return
	}

	// 必須是好友（防止陌生人騷擾）
	if !g.Friends.IsFriend(p.ID, payload.ToID) {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgDMError,
			Payload: ws.DMErrorPayload{ErrorCode: "not_friend", Message: "只能傳訊息給好友"},
		})
		return
	}

	// 嘗試即時發送
	result := g.DM.Send(p.ID, p.DisplayName, payload.ToID, payload.Content,
		func(dmMsg *dm.Message) bool {
			g.mu.RLock()
			toPlayer, online := g.Players[payload.ToID]
			g.mu.RUnlock()
			if !online || toPlayer == nil {
				return false
			}
			// 即時發送給在線玩家
			g.Hub.Send(payload.ToID, &ws.Message{
				Type: ws.MsgDMReceived,
				Payload: ws.DMReceivedPayload{
					MessageID:   dmMsg.ID,
					FromID:      dmMsg.FromID,
					FromName:    dmMsg.FromName,
					Content:     dmMsg.Content,
					SentAt:      dmMsg.SentAt.UnixMilli(),
				},
			})
			return true
		},
	)

	if !result.Success {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgDMError,
			Payload: ws.DMErrorPayload{ErrorCode: result.ErrorCode, Message: result.ErrorMsg},
		})
		return
	}

	// 取得今日發送計數
	sent, remaining := g.DM.GetDailyCount(p.ID)

	// 通知發送者成功
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgDMSent,
		Payload: ws.DMSentPayload{
			MessageID: result.MessageID,
			ToID:      payload.ToID,
			SentToday: sent,
			Remaining: remaining,
		},
	})

	log.Printf("[DM] 玩家 %s → %s: %q（今日第 %d 則）",
		p.ID, payload.ToID, payload.Content, sent)
}

// deliverPendingDMs 玩家上線時發送離線期間收到的私訊（DAY-103）
func (g *Game) deliverPendingDMs(p *player.Player) {
	msgs := g.DM.GetPending(p.ID)
	if len(msgs) == 0 {
		return
	}

	for _, dmMsg := range msgs {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgDMReceived,
			Payload: ws.DMReceivedPayload{
				MessageID:   dmMsg.ID,
				FromID:      dmMsg.FromID,
				FromName:    dmMsg.FromName,
				Content:     dmMsg.Content,
				SentAt:      dmMsg.SentAt.UnixMilli(),
				IsOffline:   true,
			},
		})
	}

	log.Printf("[DM] 玩家 %s 收到 %d 則離線訊息", p.ID, len(msgs))
}
