// guild_handler.go — 公會系統 handler（DAY-074）
package game

import (
	"log"

	"digital-twin/server/internal/game/guild"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// sendGuildUpdate 發送公會資訊給玩家
func (g *Game) sendGuildUpdate(p *player.Player) {
	guildData := g.Guild.GetPlayerGuild(p.ID)
	if guildData == nil {
		// 不在公會，發送空的更新
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgGuildUpdate,
			Payload: ws.GuildUpdatePayload{
				GuildID: "",
				MyRole:  "",
			},
		})
		return
	}

	payload := g.buildGuildUpdatePayload(guildData, p.ID)
	g.Hub.Send(p.ID, &ws.Message{
		Type:    ws.MsgGuildUpdate,
		Payload: payload,
	})
}

// buildGuildUpdatePayload 建立公會更新 Payload
func (g *Game) buildGuildUpdatePayload(guildData *guild.Guild, requesterID string) ws.GuildUpdatePayload {
	members := make([]ws.GuildMemberInfo, 0, len(guildData.Members))
	myRole := ""

	for _, m := range guildData.Members {
		if m.PlayerID == requesterID {
			myRole = string(m.Role)
		}
		members = append(members, ws.GuildMemberInfo{
			PlayerID:     m.PlayerID,
			DisplayName:  m.DisplayName,
			Role:         string(m.Role),
			IsOnline:     m.IsOnline,
			Contribution: m.Contribution,
		})
	}

	tasks := make([]ws.GuildTaskInfo, 0, len(guildData.Tasks))
	for _, t := range guildData.Tasks {
		tasks = append(tasks, ws.GuildTaskInfo{
			ID:          t.ID,
			Type:        string(t.Type),
			Name:        t.Name,
			Description: t.Description,
			Icon:        t.Icon,
			Target:      t.Target,
			Current:     t.Current,
			Reward:      t.Reward,
			Completed:   t.Completed,
			ResetAt:     t.ResetAt.UnixMilli(),
		})
	}

	return ws.GuildUpdatePayload{
		GuildID:     guildData.ID,
		Name:        guildData.Name,
		Description: guildData.Description,
		Icon:        guildData.Icon,
		Level:       guildData.Level,
		Exp:         guildData.Exp,
		Members:     members,
		Tasks:       tasks,
		TotalKills:  guildData.TotalKills,
		TotalCoins:  guildData.TotalCoins,
		MyRole:      myRole,
	}
}

// broadcastGuildUpdate 廣播公會更新給所有在線成員
func (g *Game) broadcastGuildUpdate(guildID string) {
	memberIDs := g.Guild.GetGuildMemberIDs(guildID)
	guildData := g.Guild.GetGuild(guildID)
	if guildData == nil {
		return
	}

	for _, memberID := range memberIDs {
		g.mu.RLock()
		p, online := g.Players[memberID]
		g.mu.RUnlock()

		if online && p != nil {
			payload := g.buildGuildUpdatePayload(guildData, memberID)
			g.Hub.Send(memberID, &ws.Message{
				Type:    ws.MsgGuildUpdate,
				Payload: payload,
			})
		}
	}
}

// handleCreateGuild 處理建立公會（DAY-074）
func (g *Game) handleCreateGuild(p *player.Player, msg *ws.Message) {
	var payload ws.CreateGuildPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	guildID, err := g.Guild.CreateGuild(p.ID, p.DisplayName, payload.Name, payload.Description)
	if err != nil {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgGuildError,
			Payload: ws.GuildErrorPayload{
				Operation: "create_guild",
				Message:   err.Error(),
			},
		})
		return
	}

	log.Printf("[Guild] 玩家 %s 建立公會 %s（%s）", p.ID, payload.Name, guildID)
	g.sendGuildUpdate(p)
}

// handleJoinGuild 處理加入公會（DAY-074）
func (g *Game) handleJoinGuild(p *player.Player, msg *ws.Message) {
	var payload ws.JoinGuildPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	err := g.Guild.JoinGuild(p.ID, p.DisplayName, payload.GuildID)
	if err != nil {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgGuildError,
			Payload: ws.GuildErrorPayload{
				Operation: "join_guild",
				Message:   err.Error(),
			},
		})
		return
	}

	log.Printf("[Guild] 玩家 %s 加入公會 %s", p.ID, payload.GuildID)
	// 廣播給所有公會成員
	g.broadcastGuildUpdate(payload.GuildID)
}

// handleLeaveGuild 處理退出公會（DAY-074）
func (g *Game) handleLeaveGuild(p *player.Player, msg *ws.Message) {
	guildID, err := g.Guild.LeaveGuild(p.ID)
	if err != nil {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgGuildError,
			Payload: ws.GuildErrorPayload{
				Operation: "leave_guild",
				Message:   err.Error(),
			},
		})
		return
	}

	log.Printf("[Guild] 玩家 %s 退出公會 %s", p.ID, guildID)

	// 通知玩家已退出（空公會資訊）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgGuildUpdate,
		Payload: ws.GuildUpdatePayload{
			GuildID: "",
			MyRole:  "",
		},
	})

	// 廣播給剩餘公會成員
	if g.Guild.GetGuild(guildID) != nil {
		g.broadcastGuildUpdate(guildID)
	}
}

// handleKickGuildMember 處理踢出成員（DAY-074）
func (g *Game) handleKickGuildMember(p *player.Player, msg *ws.Message) {
	var payload ws.KickGuildMemberPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	guildID := g.Guild.GetPlayerGuildID(p.ID)
	err := g.Guild.KickMember(p.ID, payload.TargetID)
	if err != nil {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgGuildError,
			Payload: ws.GuildErrorPayload{
				Operation: "kick_member",
				Message:   err.Error(),
			},
		})
		return
	}

	log.Printf("[Guild] 玩家 %s 踢出了 %s", p.ID, payload.TargetID)

	// 通知被踢出的玩家
	g.mu.RLock()
	targetPlayer, targetOnline := g.Players[payload.TargetID]
	g.mu.RUnlock()

	if targetOnline && targetPlayer != nil {
		g.Hub.Send(payload.TargetID, &ws.Message{
			Type: ws.MsgGuildUpdate,
			Payload: ws.GuildUpdatePayload{
				GuildID: "",
				MyRole:  "",
			},
		})
	}

	// 廣播給剩餘公會成員
	g.broadcastGuildUpdate(guildID)
}

// handlePromoteGuildMember 處理升職成員（DAY-074）
func (g *Game) handlePromoteGuildMember(p *player.Player, msg *ws.Message) {
	var payload ws.PromoteGuildMemberPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	guildID := g.Guild.GetPlayerGuildID(p.ID)
	err := g.Guild.PromoteMember(p.ID, payload.TargetID)
	if err != nil {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgGuildError,
			Payload: ws.GuildErrorPayload{
				Operation: "promote_member",
				Message:   err.Error(),
			},
		})
		return
	}

	log.Printf("[Guild] 玩家 %s 升職了 %s", p.ID, payload.TargetID)
	g.broadcastGuildUpdate(guildID)
}

// handleGetGuildInfo 處理查詢公會資訊（DAY-074）
func (g *Game) handleGetGuildInfo(p *player.Player, msg *ws.Message) {
	g.sendGuildUpdate(p)
}

// handleGetGuildList 處理查詢公會列表（DAY-074）
func (g *Game) handleGetGuildList(p *player.Player, msg *ws.Message) {
	allGuilds := g.Guild.GetAllGuilds()
	entries := make([]ws.GuildListEntry, 0, len(allGuilds))

	for _, gd := range allGuilds {
		onlineCount := 0
		for _, m := range gd.Members {
			if m.IsOnline {
				onlineCount++
			}
		}
		entries = append(entries, ws.GuildListEntry{
			GuildID:     gd.ID,
			Name:        gd.Name,
			Description: gd.Description,
			Icon:        gd.Icon,
			Level:       gd.Level,
			MemberCount: len(gd.Members),
			OnlineCount: onlineCount,
		})
	}

	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgGuildList,
		Payload: ws.GuildListPayload{
			Guilds: entries,
		},
	})
}

// notifyGuildTaskComplete 通知公會任務完成（DAY-074）
// 給所有在線公會成員發放獎勵並通知
func (g *Game) notifyGuildTaskComplete(guildID string, completedTasks []*guild.GuildTask) {
	if len(completedTasks) == 0 {
		return
	}

	memberIDs := g.Guild.GetGuildMemberIDs(guildID)
	guildData := g.Guild.GetGuild(guildID)
	if guildData == nil {
		return
	}

	for _, task := range completedTasks {
		for _, memberID := range memberIDs {
			g.mu.RLock()
			p, online := g.Players[memberID]
			g.mu.RUnlock()

			if !online || p == nil {
				continue
			}

			// 發放獎勵
			p.AddCoins(task.Reward)
			newBalance := p.Snapshot().Coins

			g.Hub.Send(memberID, &ws.Message{
				Type: ws.MsgGuildTaskComplete,
				Payload: ws.GuildTaskCompletePayload{
					GuildID:    guildID,
					GuildName:  guildData.Name,
					TaskID:     task.ID,
					TaskName:   task.Name,
					TaskIcon:   task.Icon,
					Reward:     task.Reward,
					NewBalance: newBalance,
				},
			})
		}
		log.Printf("[Guild] 公會 %s 完成任務 %s，每人獎勵 %d 金幣", guildID, task.Name, task.Reward)
	}
}
