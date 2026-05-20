// friendchallenge_handler.go — 好友挑戰系統 handler（DAY-102）
// 好友間 1v1 挑戰：3 分鐘內比較分數，勝者獲得全部賭注
package game

import (
	"log"
	"time"

	"digital-twin/server/internal/game/friendchallenge"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// handleSendChallengeRequest 處理發起挑戰請求（DAY-102）
func (g *Game) handleSendChallengeRequest(p *player.Player, msg *ws.Message) {
	var payload ws.SendChallengeRequestPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	if payload.FriendID == "" || payload.FriendID == p.ID {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgChallengeError,
			Payload: ws.ChallengeErrorPayload{ErrorCode: "invalid_target", Message: "無效的挑戰對象"},
		})
		return
	}

	// 必須是好友
	if !g.Friends.IsFriend(p.ID, payload.FriendID) {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgChallengeError,
			Payload: ws.ChallengeErrorPayload{ErrorCode: "not_friend", Message: "只能挑戰好友"},
		})
		return
	}

	// 檢查金幣是否足夠
	if p.Coins < friendchallenge.ChallengeStake {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgChallengeError,
			Payload: ws.ChallengeErrorPayload{ErrorCode: "insufficient_coins", Message: "金幣不足（需要 1000🪙）"},
		})
		return
	}

	// 取得好友顯示名稱
	friendName := payload.FriendID
	g.mu.RLock()
	friendPlayer, friendOnline := g.Players[payload.FriendID]
	g.mu.RUnlock()
	if friendOnline && friendPlayer != nil {
		friendName = friendPlayer.DisplayName
	}

	c, errCode := g.FriendChallenge.CreateChallenge(p.ID, p.DisplayName, payload.FriendID, friendName)
	if errCode != "" {
		var msg string
		switch errCode {
		case "already_in_challenge":
			msg = "你已在挑戰中"
		case "opponent_in_challenge":
			msg = "對方已在挑戰中"
		default:
			msg = "無法發起挑戰"
		}
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgChallengeError,
			Payload: ws.ChallengeErrorPayload{ErrorCode: errCode, Message: msg},
		})
		return
	}

	// 通知挑戰者
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgChallengeUpdate,
		Payload: ws.ChallengeUpdatePayload{
			ChallengeID:    c.ID,
			Status:         string(c.Status),
			OpponentID:     payload.FriendID,
			OpponentName:   friendName,
			Stake:          c.Stake,
			MyScore:        0,
			OpponentScore:  0,
			TimeRemaining:  0,
		},
	})

	// 通知被挑戰者（如果在線）
	if friendOnline && friendPlayer != nil {
		g.Hub.Send(payload.FriendID, &ws.Message{
			Type: ws.MsgChallengeRequest,
			Payload: ws.ChallengeRequestPayload{
				ChallengeID:     c.ID,
				ChallengerID:    p.ID,
				ChallengerName:  p.DisplayName,
				Stake:           c.Stake,
				ExpiresInSec:    int(friendchallenge.PendingTimeout.Seconds()),
			},
		})
	}

	log.Printf("[Challenge] 玩家 %s 向 %s 發起挑戰（ID: %s）", p.ID, payload.FriendID, c.ID)
}

// handleAcceptChallenge 處理接受挑戰（DAY-102）
func (g *Game) handleAcceptChallenge(p *player.Player, msg *ws.Message) {
	var payload ws.AcceptChallengePayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	// 檢查金幣是否足夠
	if p.Coins < friendchallenge.ChallengeStake {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgChallengeError,
			Payload: ws.ChallengeErrorPayload{ErrorCode: "insufficient_coins", Message: "金幣不足（需要 1000🪙）"},
		})
		return
	}

	c, errCode := g.FriendChallenge.AcceptChallenge(payload.ChallengeID, p.ID)
	if errCode != "" {
		var errMsg string
		switch errCode {
		case "not_found":
			errMsg = "找不到挑戰"
		case "expired":
			errMsg = "挑戰已過期"
		default:
			errMsg = "無法接受挑戰"
		}
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgChallengeError,
			Payload: ws.ChallengeErrorPayload{ErrorCode: errCode, Message: errMsg},
		})
		return
	}

	// 扣除雙方賭注
	p.Coins -= c.Stake
	g.mu.RLock()
	challenger, challengerOnline := g.Players[c.ChallengerID]
	g.mu.RUnlock()
	if challengerOnline && challenger != nil {
		challenger.Coins -= c.Stake
		g.sendPlayerUpdate(challenger)
	}
	g.sendPlayerUpdate(p)

	// 通知雙方挑戰開始
	startPayload := ws.ChallengeUpdatePayload{
		ChallengeID:   c.ID,
		Status:        string(c.Status),
		OpponentID:    c.ChallengerID,
		OpponentName:  c.ChallengerName,
		Stake:         c.Stake,
		MyScore:       0,
		OpponentScore: 0,
		TimeRemaining: c.TimeRemaining(),
	}
	g.Hub.Send(p.ID, &ws.Message{
		Type:    ws.MsgChallengeUpdate,
		Payload: startPayload,
	})

	if challengerOnline && challenger != nil {
		g.Hub.Send(c.ChallengerID, &ws.Message{
			Type: ws.MsgChallengeUpdate,
			Payload: ws.ChallengeUpdatePayload{
				ChallengeID:   c.ID,
				Status:        string(c.Status),
				OpponentID:    c.ChallengedID,
				OpponentName:  c.ChallengedName,
				Stake:         c.Stake,
				MyScore:       0,
				OpponentScore: 0,
				TimeRemaining: c.TimeRemaining(),
			},
		})
	}

	log.Printf("[Challenge] 挑戰 %s 開始！%s vs %s（各賭 %d 金幣，3分鐘）",
		c.ID, c.ChallengerName, c.ChallengedName, c.Stake)
}

// handleDeclineChallenge 處理拒絕挑戰（DAY-102）
func (g *Game) handleDeclineChallenge(p *player.Player, msg *ws.Message) {
	var payload ws.DeclineChallengePayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	ok := g.FriendChallenge.DeclineChallenge(payload.ChallengeID, p.ID)
	if !ok {
		return
	}

	// 通知挑戰者被拒絕
	c := g.FriendChallenge.GetChallengeByID(payload.ChallengeID)
	if c == nil {
		return
	}
	g.mu.RLock()
	challenger, challengerOnline := g.Players[c.ChallengerID]
	g.mu.RUnlock()
	if challengerOnline && challenger != nil {
		g.Hub.Send(c.ChallengerID, &ws.Message{
			Type: ws.MsgChallengeError,
			Payload: ws.ChallengeErrorPayload{
				ErrorCode: "declined",
				Message:   p.DisplayName + " 拒絕了你的挑戰",
			},
		})
	}
	log.Printf("[Challenge] 玩家 %s 拒絕了 %s 的挑戰", p.ID, c.ChallengerID)
}

// notifyChallengeKillScore 擊破目標時更新挑戰分數（由 handleKill 呼叫）
func (g *Game) notifyChallengeKillScore(p *player.Player, reward int) {
	if !g.FriendChallenge.IsInChallenge(p.ID) {
		return
	}

	// 挑戰分數 = 獎勵金幣（反映玩家的實際表現）
	c := g.FriendChallenge.AddScore(p.ID, reward)
	if c == nil {
		return
	}

	// 取得我的分數和對手分數
	myScore := c.ChallengerScore
	opponentScore := c.ChallengedScore
	opponentID := c.ChallengedID
	if p.ID == c.ChallengedID {
		myScore = c.ChallengedScore
		opponentScore = c.ChallengerScore
		opponentID = c.ChallengerID
	}

	// 通知雙方分數更新
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgChallengeUpdate,
		Payload: ws.ChallengeUpdatePayload{
			ChallengeID:   c.ID,
			Status:        string(c.Status),
			OpponentID:    opponentID,
			MyScore:       myScore,
			OpponentScore: opponentScore,
			TimeRemaining: c.TimeRemaining(),
		},
	})

	// 通知對手分數更新
	g.mu.RLock()
	opponent, opponentOnline := g.Players[opponentID]
	g.mu.RUnlock()
	if opponentOnline && opponent != nil {
		g.Hub.Send(opponentID, &ws.Message{
			Type: ws.MsgChallengeUpdate,
			Payload: ws.ChallengeUpdatePayload{
				ChallengeID:   c.ID,
				Status:        string(c.Status),
				OpponentID:    p.ID,
				MyScore:       opponentScore,
				OpponentScore: myScore,
				TimeRemaining: c.TimeRemaining(),
			},
		})
	}
}

// tickAndFinishChallenges 定期檢查並結算到期的挑戰（由 gameLoop 每 5 秒呼叫）
func (g *Game) tickAndFinishChallenges() {
	finished := g.FriendChallenge.CheckAndFinish()
	for _, c := range finished {
		g.settleChallengeResult(c)
	}
}

// settleChallengeResult 結算挑戰結果並發放獎勵
func (g *Game) settleChallengeResult(c *friendchallenge.Challenge) {
	g.mu.RLock()
	challenger, challengerOnline := g.Players[c.ChallengerID]
	challenged, challengedOnline := g.Players[c.ChallengedID]
	g.mu.RUnlock()

	isDraw := c.WinnerID == ""

	// 發放獎勵
	if isDraw {
		// 平局：各退回賭注
		if challengerOnline && challenger != nil {
			challenger.Coins += c.Stake
			g.sendPlayerUpdate(challenger)
		}
		if challengedOnline && challenged != nil {
			challenged.Coins += c.Stake
			g.sendPlayerUpdate(challenged)
		}
	} else {
		// 勝者獲得全部賭注
		g.mu.RLock()
		winner, winnerOnline := g.Players[c.WinnerID]
		g.mu.RUnlock()
		if winnerOnline && winner != nil {
			winner.Coins += c.Prize
			g.sendPlayerUpdate(winner)
		}
	}

	// 建立結果 payload
	buildResultPayload := func(playerID string) ws.ChallengeResultPayload {
		myScore := c.ChallengerScore
		opponentScore := c.ChallengedScore
		opponentID := c.ChallengedID
		opponentName := c.ChallengedName
		if playerID == c.ChallengedID {
			myScore = c.ChallengedScore
			opponentScore = c.ChallengerScore
			opponentID = c.ChallengerID
			opponentName = c.ChallengerName
		}
		isWinner := c.WinnerID == playerID
		prize := 0
		if isWinner {
			prize = c.Prize
		} else if isDraw {
			prize = c.Stake // 退回賭注
		}
		return ws.ChallengeResultPayload{
			ChallengeID:   c.ID,
			IsWinner:      isWinner,
			IsDraw:        isDraw,
			WinnerName:    c.WinnerName,
			MyScore:       myScore,
			OpponentScore: opponentScore,
			OpponentID:    opponentID,
			OpponentName:  opponentName,
			Prize:         prize,
		}
	}

	// 通知雙方結果
	if challengerOnline && challenger != nil {
		g.Hub.Send(c.ChallengerID, &ws.Message{
			Type:    ws.MsgChallengeResult,
			Payload: buildResultPayload(c.ChallengerID),
		})
	}
	if challengedOnline && challenged != nil {
		g.Hub.Send(c.ChallengedID, &ws.Message{
			Type:    ws.MsgChallengeResult,
			Payload: buildResultPayload(c.ChallengedID),
		})
	}

	if isDraw {
		log.Printf("[Challenge] 挑戰 %s 平局！%s=%d vs %s=%d，各退回 %d 金幣",
			c.ID, c.ChallengerName, c.ChallengerScore, c.ChallengedName, c.ChallengedScore, c.Stake)
	} else {
		log.Printf("[Challenge] 挑戰 %s 結束！勝者：%s（%d 金幣）",
			c.ID, c.WinnerName, c.Prize)
	}
}

// startChallengeTicker 啟動挑戰結算計時器（每 5 秒檢查一次）
func (g *Game) startChallengeTicker() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				g.tickAndFinishChallenges()
			case <-g.stopCh:
				return
			}
		}
	}()
}
