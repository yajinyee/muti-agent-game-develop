// friend_handler.go — 好友系統 handler（DAY-073）
package game

import (
	"log"

	"digital-twin/server/internal/player"
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
