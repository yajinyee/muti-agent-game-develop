// friend_handler.go — 好友系統 handler（DAY-073）
// DAY-101：新增禮物贈送系統 + 好友持久化
package game

import (
	"log"

	"digital-twin/server/internal/game/friend"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/store"
	"digital-twin/server/internal/ws"
)

// sendFriendList 發送好友列表給玩家
func (g *Game) sendFriendList(p *player.Player) {
	friendIDs := g.Friends.GetFriendIDs(p.ID)
	pendingReqs := g.Friends.GetPendingRequests(p.ID)

	friends := make([]ws.FriendInfoPayload, 0, len(friendIDs))
	for _, fid := range friendIDs {
		info := g.buildFriendInfo(fid)
		friends = append(friends, info)
	}

	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgFriendList,
		Payload: ws.FriendListPayload{
			Friends:      friends,
			PendingCount: len(pendingReqs),
		},
	})
}

// buildFriendInfo 建立好友資訊（從在線玩家或快取）
func (g *Game) buildFriendInfo(playerID string) ws.FriendInfoPayload {
	g.mu.RLock()
	p, online := g.Players[playerID]
	g.mu.RUnlock()

	info := ws.FriendInfoPayload{
		PlayerID: playerID,
		IsOnline: online,
	}

	if online && p != nil {
		snap := p.Snapshot()
		info.DisplayName = snap.DisplayName
		info.Coins = snap.Coins
		info.KillCount = snap.KillCount
		info.TitleName = snap.TitleName
		info.TitleIcon = snap.TitleIcon
		// 賽季資料
		seasonSnap := g.Season.GetSnapshot(playerID)
		info.SeasonLevel = seasonSnap.CurrentLevel
		info.SeasonPoints = seasonSnap.SeasonPoints
	} else {
		// 離線玩家：從 Store 取得基本資料
		if g.store != nil {
			if state, err := g.store.LoadPlayer(playerID); err == nil && state != nil {
				info.DisplayName = state.DisplayName
				info.Coins = int(state.Coins)
			}
		}
		if info.DisplayName == "" {
			// 取 ID 前 8 碼作為顯示名稱
			if len(playerID) > 8 {
				info.DisplayName = playerID[:8]
			} else {
				info.DisplayName = playerID
			}
		}
	}

	return info
}

// handleSendFriendRequest 處理發送好友請求（DAY-073）
func (g *Game) handleSendFriendRequest(p *player.Player, msg *ws.Message) {
	var payload ws.SendFriendRequestPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	if payload.TargetID == "" || payload.TargetID == p.ID {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: map[string]string{"message": "無效的目標玩家"},
		})
		return
	}

	ok := g.Friends.SendRequest(p.ID, payload.TargetID)
	if !ok {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: map[string]string{"message": "無法發送好友請求（已是好友或已有待處理請求）"},
		})
		return
	}

	// 如果目標玩家在線，通知他有好友請求
	g.mu.RLock()
	targetPlayer, targetOnline := g.Players[payload.TargetID]
	g.mu.RUnlock()

	if targetOnline && targetPlayer != nil {
		g.Hub.Send(payload.TargetID, &ws.Message{
			Type: ws.MsgFriendRequest,
			Payload: ws.FriendRequestPayload{
				FromID:      p.ID,
				DisplayName: p.DisplayName,
			},
		})
	}

	// 如果已成為好友（互相發請求），更新雙方好友列表
	if g.Friends.IsFriend(p.ID, payload.TargetID) {
		g.sendFriendList(p)
		if targetOnline && targetPlayer != nil {
			g.sendFriendList(targetPlayer)
		}
		log.Printf("[Friend] 玩家 %s 和 %s 互相發請求，已成為好友", p.ID, payload.TargetID)
	} else {
		log.Printf("[Friend] 玩家 %s 向 %s 發送好友請求", p.ID, payload.TargetID)
	}
}

// handleAcceptFriendRequest 處理接受好友請求（DAY-073）
func (g *Game) handleAcceptFriendRequest(p *player.Player, msg *ws.Message) {
	var payload ws.AcceptFriendRequestPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	ok := g.Friends.AcceptRequest(payload.FromID, p.ID)
	if !ok {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: map[string]string{"message": "找不到好友請求"},
		})
		return
	}

	// 更新雙方好友列表
	g.sendFriendList(p)

	g.mu.RLock()
	fromPlayer, fromOnline := g.Players[payload.FromID]
	g.mu.RUnlock()

	if fromOnline && fromPlayer != nil {
		g.sendFriendList(fromPlayer)
		// 通知對方請求已被接受
		g.Hub.Send(payload.FromID, &ws.Message{
			Type: ws.MsgFriendUpdate,
			Payload: ws.FriendUpdatePayload{
				FriendID:    p.ID,
				DisplayName: p.DisplayName,
				IsOnline:    true,
				Event:       "accepted",
			},
		})
	}

	// 持久化雙方好友關係（DAY-101）
	go g.saveFriendState(p.ID)
	go g.saveFriendState(payload.FromID)

	log.Printf("[Friend] 玩家 %s 接受了 %s 的好友請求", p.ID, payload.FromID)
}

// handleRejectFriendRequest 處理拒絕好友請求（DAY-073）
func (g *Game) handleRejectFriendRequest(p *player.Player, msg *ws.Message) {
	var payload ws.RejectFriendRequestPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	g.Friends.RejectRequest(payload.FromID, p.ID)
	log.Printf("[Friend] 玩家 %s 拒絕了 %s 的好友請求", p.ID, payload.FromID)
}

// handleRemoveFriend 處理移除好友（DAY-073）
func (g *Game) handleRemoveFriend(p *player.Player, msg *ws.Message) {
	var payload ws.RemoveFriendPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	ok := g.Friends.RemoveFriend(p.ID, payload.FriendID)
	if !ok {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: map[string]string{"message": "找不到此好友"},
		})
		return
	}

	// 更新雙方好友列表
	g.sendFriendList(p)

	g.mu.RLock()
	friendPlayer, friendOnline := g.Players[payload.FriendID]
	g.mu.RUnlock()

	if friendOnline && friendPlayer != nil {
		g.sendFriendList(friendPlayer)
		g.Hub.Send(payload.FriendID, &ws.Message{
			Type: ws.MsgFriendUpdate,
			Payload: ws.FriendUpdatePayload{
				FriendID:    p.ID,
				DisplayName: p.DisplayName,
				IsOnline:    false,
				Event:       "removed",
			},
		})
	}

	// 持久化雙方好友關係（DAY-101）
	go g.saveFriendState(p.ID)
	go g.saveFriendState(payload.FriendID)

	log.Printf("[Friend] 玩家 %s 移除了好友 %s", p.ID, payload.FriendID)
}

// handleGetFriendList 處理查詢好友列表（DAY-073）
func (g *Game) handleGetFriendList(p *player.Player, msg *ws.Message) {
	g.sendFriendList(p)
}

// notifyFriendsOnline 通知好友玩家上線（DAY-073）
func (g *Game) notifyFriendsOnline(playerID string, displayName string) {
	friendIDs := g.Friends.GetFriendIDs(playerID)
	for _, fid := range friendIDs {
		g.mu.RLock()
		fp, online := g.Players[fid]
		g.mu.RUnlock()
		if online && fp != nil {
			g.Hub.Send(fid, &ws.Message{
				Type: ws.MsgFriendUpdate,
				Payload: ws.FriendUpdatePayload{
					FriendID:    playerID,
					DisplayName: displayName,
					IsOnline:    true,
					Event:       "online",
				},
			})
		}
	}
}

// notifyFriendsOffline 通知好友玩家下線（DAY-073）
func (g *Game) notifyFriendsOffline(playerID string, displayName string) {
	friendIDs := g.Friends.GetFriendIDs(playerID)
	for _, fid := range friendIDs {
		g.mu.RLock()
		fp, online := g.Players[fid]
		g.mu.RUnlock()
		if online && fp != nil {
			g.Hub.Send(fid, &ws.Message{
				Type: ws.MsgFriendUpdate,
				Payload: ws.FriendUpdatePayload{
					FriendID:    playerID,
					DisplayName: displayName,
					IsOnline:    false,
					Event:       "offline",
				},
			})
		}
	}
}

// ---- 好友禮物系統（DAY-101）----

// handleSendGift 處理送禮物請求（DAY-101）
func (g *Game) handleSendGift(p *player.Player, msg *ws.Message) {
	var payload ws.SendGiftPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	if payload.FriendID == "" {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgGiftError,
			Payload: ws.GiftErrorPayload{ErrorCode: "invalid_target", Message: "無效的好友 ID"},
		})
		return
	}

	result := g.Friends.SendGift(p.ID, payload.FriendID)
	if !result.Success {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgGiftError,
			Payload: ws.GiftErrorPayload{ErrorCode: result.ErrorCode, Message: result.ErrorMsg},
		})
		return
	}

	// 取得好友顯示名稱
	friendDisplayName := payload.FriendID
	g.mu.RLock()
	friendPlayer, friendOnline := g.Players[payload.FriendID]
	g.mu.RUnlock()
	if friendOnline && friendPlayer != nil {
		friendDisplayName = friendPlayer.DisplayName
		// 發放金幣給好友
		friendPlayer.Coins += result.Amount
		// 通知好友收到禮物
		g.Hub.Send(payload.FriendID, &ws.Message{
			Type: ws.MsgGiftReceived,
			Payload: ws.GiftReceivedPayload{
				FromID:      p.ID,
				DisplayName: p.DisplayName,
				Amount:      result.Amount,
				NewBalance:  friendPlayer.Coins,
			},
		})
		// 更新好友的玩家狀態顯示
		g.sendPlayerUpdate(friendPlayer)
	} else {
		// 好友離線：儲存待領取禮物（用 KV store 暫存）
		if g.store != nil {
			pendingKey := "pending_gift:" + payload.FriendID
			var pending []int
			_ = g.store.GetJSON(pendingKey, &pending)
			pending = append(pending, result.Amount)
			_ = g.store.SetJSON(pendingKey, pending, 0)
		}
		// 嘗試從 store 取得好友名稱
		if g.store != nil {
			if state, err := g.store.LoadPlayer(payload.FriendID); err == nil && state != nil {
				friendDisplayName = state.DisplayName
			}
		}
	}

	// 取得今日禮物狀態
	sentToday, remaining := g.Friends.GetGiftStatus(p.ID)

	// 通知送禮者成功
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgGiftSent,
		Payload: ws.GiftSentPayload{
			ToID:        payload.FriendID,
			DisplayName: friendDisplayName,
			Amount:      result.Amount,
			SentToday:   sentToday,
			Remaining:   remaining,
		},
	})

	log.Printf("[Gift] 玩家 %s 送禮物 %d 金幣給 %s（今日第 %d 次）",
		p.ID, result.Amount, payload.FriendID, sentToday)
}

// handleGetGiftStatus 處理查詢禮物狀態（DAY-101）
func (g *Game) handleGetGiftStatus(p *player.Player) {
	sentToday, remaining := g.Friends.GetGiftStatus(p.ID)
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgGiftStatus,
		Payload: ws.GiftStatusPayload{
			SentToday: sentToday,
			Remaining: remaining,
			MaxDaily:  3,
			Amount:    500,
		},
	})
}

// deliverPendingGifts 玩家上線時發放離線期間收到的禮物（DAY-101）
func (g *Game) deliverPendingGifts(p *player.Player) {
	if g.store == nil {
		return
	}
	pendingKey := "pending_gift:" + p.ID
	var pending []int
	if err := g.store.GetJSON(pendingKey, &pending); err != nil || len(pending) == 0 {
		return
	}

	total := 0
	for _, amount := range pending {
		total += amount
	}
	if total <= 0 {
		return
	}

	p.Coins += total
	// 清除待領取禮物
	_ = g.store.SetJSON(pendingKey, []int{}, 0)

	// 通知玩家收到離線禮物
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgGiftReceived,
		Payload: ws.GiftReceivedPayload{
			FromID:      "system",
			DisplayName: "離線禮物",
			Amount:      total,
			NewBalance:  p.Coins,
		},
	})
	log.Printf("[Gift] 玩家 %s 收到離線禮物共 %d 金幣（%d 份）", p.ID, total, len(pending))
}

// ---- 好友持久化（DAY-101）----

// saveFriendState 儲存玩家好友關係到 FileStore
func (g *Game) saveFriendState(playerID string) {
	fs, ok := g.store.(*store.FileStore)
	if !ok {
		return
	}
	friendIDs := g.Friends.GetFriendIDs(playerID)
	if err := fs.SaveFriends(playerID, friendIDs); err != nil {
		log.Printf("[Friend] Failed to save friends for %s: %v", playerID, err)
	}
}

// restoreFriendState 從 FileStore 恢復玩家好友關係
func (g *Game) restoreFriendState(playerID string) {
	fs, ok := g.store.(*store.FileStore)
	if !ok {
		return
	}
	friendIDs, err := fs.LoadFriends(playerID)
	if err != nil {
		log.Printf("[Friend] Failed to load friends for %s: %v", playerID, err)
		return
	}
	if len(friendIDs) == 0 {
		return
	}
	g.Friends.LoadFriendState(&friend.FriendState{
		PlayerID:  playerID,
		FriendIDs: friendIDs,
	})
	log.Printf("[Friend] Player %s restored %d friends", playerID, len(friendIDs))
}
