// fevermode_handler.go — 狂熱模式系統 handler（DAY-133）
// 業界依據：Fire Kirin / Ocean King 系列的 Fever Mode
// 玩家在 5 秒內擊破 5 個目標觸發狂熱模式，期間所有獎勵 ×1.5，
// 繼續快速擊破可延長時間（最多 30 秒），製造「停不下來」的爽感。
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/game/activityfeed"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// notifyFeverModeKill 在擊破目標後更新狂熱模式（由 handleKill 呼叫）
// 回傳狂熱模式倍率加成（用於最終獎勵計算）
func (g *Game) notifyFeverModeKill(p *player.Player) float64 {
	if g.FeverMode == nil {
		return 1.0
	}

	triggered, extended, multBoost := g.FeverMode.RecordKill(p.ID)

	if triggered {
		// 狂熱模式觸發！廣播給所有玩家
		snap := g.FeverMode.GetSnapshot(p.ID)
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgFeverModeStart,
			Payload: ws.FeverModeStartPayload{
				PlayerID:    p.ID,
				PlayerName:  p.DisplayName,
				SecondsLeft: snap.SecondsLeft,
				MultBoost:   multBoost,
			},
		})
		log.Printf("[FeverMode] player=%s triggered fever! mult=%.1f, secs=%d",
			p.ID, multBoost, snap.SecondsLeft)

		// 全服公告
		g.announceFeverMode(p.DisplayName)

		// 動態牆
		go g.notifyFeedFeverMode(p)
	} else if extended {
		// 狂熱延長：只發送給個人（避免頻繁廣播）
		snap := g.FeverMode.GetSnapshot(p.ID)
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgFeverModeStatus,
			Payload: ws.FeverModeStatusPayload{
				PlayerID:     p.ID,
				IsActive:     true,
				SecondsLeft:  snap.SecondsLeft,
				CooldownLeft: 0,
				MultBoost:    multBoost,
				KillProgress: 0,
				TriggerKills: g.FeverMode.GetConfig().TriggerKills,
				TotalFevered: snap.TotalFevered,
			},
		})
	} else if multBoost == 1.0 {
		// 未觸發：發送進度更新（讓 Client 顯示進度條）
		snap := g.FeverMode.GetSnapshot(p.ID)
		if snap.KillProgress > 0 {
			g.Hub.Send(p.ID, &ws.Message{
				Type: ws.MsgFeverModeStatus,
				Payload: ws.FeverModeStatusPayload{
					PlayerID:     p.ID,
					IsActive:     false,
					SecondsLeft:  0,
					CooldownLeft: snap.CooldownLeft,
					MultBoost:    1.0,
					KillProgress: snap.KillProgress,
					TriggerKills: g.FeverMode.GetConfig().TriggerKills,
					TotalFevered: snap.TotalFevered,
				},
			})
		}
	}

	return multBoost
}

// tickFeverModeExpiry 定期檢查狂熱模式過期（由 gameLoop 每秒呼叫）
func (g *Game) tickFeverModeExpiry() {
	if g.FeverMode == nil {
		return
	}

	expired := g.FeverMode.TickExpiry()
	for _, playerID := range expired {
		g.mu.RLock()
		p := g.Players[playerID]
		g.mu.RUnlock()

		snap := g.FeverMode.GetSnapshot(playerID)
		if err := g.Hub.Send(playerID, &ws.Message{
			Type: ws.MsgFeverModeEnd,
			Payload: ws.FeverModeEndPayload{
				PlayerID:     playerID,
				TotalFevered: snap.TotalFevered,
				CooldownLeft: snap.CooldownLeft,
			},
		}); err != nil {
			log.Printf("[FeverMode] send end error: %v", err)
		}

		if p != nil {
			log.Printf("[FeverMode] player=%s fever ended, total=%d", playerID, snap.TotalFevered)
		}
	}
}

// sendFeverModeStatus 發送狂熱模式狀態給玩家（登入時呼叫）
func (g *Game) sendFeverModeStatus(p *player.Player) {
	if g.FeverMode == nil {
		return
	}
	snap := g.FeverMode.GetSnapshot(p.ID)
	cfg := g.FeverMode.GetConfig()
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgFeverModeStatus,
		Payload: ws.FeverModeStatusPayload{
			PlayerID:     p.ID,
			IsActive:     snap.State == 1, // FeverStateActive
			SecondsLeft:  snap.SecondsLeft,
			CooldownLeft: snap.CooldownLeft,
			MultBoost:    snap.MultBoost,
			KillProgress: snap.KillProgress,
			TriggerKills: cfg.TriggerKills,
			TotalFevered: snap.TotalFevered,
		},
	})
}

// announceFeverMode 全服公告：狂熱模式觸發
func (g *Game) announceFeverMode(playerName string) {
	extra := map[string]string{
		"mult": fmt.Sprintf("%.1f", g.FeverMode.GetConfig().MultBoost),
	}
	ann := g.Announce.Create(announce.EventFeverMode, playerName, 0, extra)
	g.broadcastAnnouncement(ann)
}

// notifyFeedFeverMode 動態牆：狂熱模式觸發
func (g *Game) notifyFeedFeverMode(p *player.Player) {
	if g.ActivityFeed == nil {
		return
	}
	evt := g.ActivityFeed.Push(&activityfeed.FeedEvent{
		EventType:   activityfeed.EventFeverMode,
		PlayerID:    p.ID,
		DisplayName: p.DisplayName,
		Icon:        "🔥",
		Title:       "狂熱模式",
		Detail:      fmt.Sprintf("進入狂熱模式！獎勵 ×%.1f！", g.FeverMode.GetConfig().MultBoost),
		Rarity:      activityfeed.RarityRare,
	})
	go g.broadcastFeedEvent(evt)
}
