// Package game — Mission（每日任務）相關 handler（DAY-057 拆分自 game.go）
package game

import (
	"log"

	"digital-twin/server/internal/game/mission"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// sendMissionUpdate 傳送任務列表給指定玩家
func (g *Game) sendMissionUpdate(playerID string) {
	statuses := g.missionMgr.GetPlayerMissions(playerID)
	payloads := make([]ws.MissionPayload, 0, len(statuses))
	for _, s := range statuses {
		payloads = append(payloads, ws.MissionPayload{
			ID:            s.Mission.ID,
			Name:          s.Mission.Name,
			Description:   s.Mission.Description,
			Icon:          s.Mission.Icon,
			Target:        s.Mission.Target,
			Current:       s.Progress.Current,
			Completed:     s.Progress.Completed,
			RewardClaimed: s.Progress.RewardClaimed,
			Reward:        s.Mission.Reward,
		})
	}
	g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgMissionUpdate,
		Payload: ws.MissionUpdatePayload{
			PlayerID:      playerID,
			Missions:      payloads,
			ResetAt:       g.missionMgr.ResetAt().UnixMilli(),
			ResetTimezone: "UTC+8",
		},
	})
}

// updateMissionProgress 更新任務進度並通知玩家
// 由各遊戲事件（擊殺、BOSS、Bonus）呼叫
func (g *Game) updateMissionProgress(playerID string, mType mission.MissionType, amount int) {
	completed := g.missionMgr.UpdateProgress(playerID, mType, amount)

	// 通知任務完成
	for _, m := range completed {
		log.Printf("[Mission] Player %s completed: %s", playerID, m.Name)
		g.Hub.Send(playerID, &ws.Message{
			Type: ws.MsgMissionComplete,
			Payload: ws.MissionCompletePayload{
				MissionID: m.ID,
				Name:      m.Name,
				Icon:      m.Icon,
				Reward:    m.Reward,
			},
		})
	}

	// 更新任務進度（有變化才發送）
	if amount > 0 {
		g.sendMissionUpdate(playerID)
	}
}

// handleClaimMission 處理領取任務獎勵
func (g *Game) handleClaimMission(p *player.Player, msg *ws.Message) {
	var payload ws.ClaimMissionPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}
	if payload.MissionID == "" {
		return
	}

	reward := g.missionMgr.ClaimReward(p.ID, payload.MissionID)
	if reward <= 0 {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "mission_not_claimable", Message: "任務未完成或已領取"},
		})
		return
	}

	// 發放獎勵
	p.AddReward(reward)
	log.Printf("[Mission] Player %s claimed reward %d for mission %s", p.ID, reward, payload.MissionID)

	// 通知玩家
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgReward,
		Payload: ws.RewardPayload{
			Source:     "mission",
			Amount:     reward,
			Multiplier: 1.0,
			NewBalance: p.Coins,
		},
	})
	g.sendPlayerUpdate(p)
	g.sendMissionUpdate(p.ID)
}
