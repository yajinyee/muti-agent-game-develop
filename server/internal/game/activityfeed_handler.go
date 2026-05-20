// activityfeed_handler.go — 成就動態牆 handler（DAY-112）
package game

import (
	"log"

	"digital-twin/server/internal/game/activityfeed"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// handleGetActivityFeed 處理玩家查詢最近動態（DAY-112）
func (g *Game) handleGetActivityFeed(p *player.Player) {
	g.sendActivityFeedHistory(p.ID)
}

// sendActivityFeedHistory 發送最近 10 條動態給特定玩家（DAY-112）
func (g *Game) sendActivityFeedHistory(playerID string) {
	events := g.ActivityFeed.GetRecent(10)
	payloads := make([]ws.ActivityFeedEventPayload, len(events))
	for i, evt := range events {
		payloads[i] = feedEventToPayload(evt)
	}

	if err := g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgActivityFeedHistory,
		Payload: ws.ActivityFeedHistoryPayload{
			Events: payloads,
			Total:  len(g.ActivityFeed.GetRecent(50)),
		},
	}); err != nil {
		log.Printf("[ActivityFeed] send history error: %v", err)
	}
}

// broadcastFeedEvent 廣播新動態事件給所有玩家（DAY-112）
func (g *Game) broadcastFeedEvent(evt *activityfeed.FeedEvent) {
	if evt == nil {
		return
	}
	payload := feedEventToPayload(evt)

	g.mu.RLock()
	playerIDs := make([]string, 0, len(g.Players))
	for id := range g.Players {
		playerIDs = append(playerIDs, id)
	}
	g.mu.RUnlock()

	for _, pid := range playerIDs {
		g.Hub.Send(pid, &ws.Message{
			Type:    ws.MsgActivityFeedEvent,
			Payload: payload,
		})
	}
}

// ---- 各系統觸發點 ----

// notifyFeedAchievement 成就解鎖時推送動態（DAY-112）
func (g *Game) notifyFeedAchievement(p *player.Player, achName, achIcon, achType string) {
	evt := g.ActivityFeed.Push(activityfeed.NewAchievementEvent(
		p.ID, p.DisplayName, achName, achIcon, achType,
	))
	go g.broadcastFeedEvent(evt)
}

// notifyFeedTitle 稱號獲得時推送動態（DAY-112）
func (g *Game) notifyFeedTitle(p *player.Player, titleName, titleIcon string, priority int) {
	// 只廣播優先級 ≥ 20 的稱號（避免太多低優先級稱號刷屏）
	if priority < 20 {
		return
	}
	evt := g.ActivityFeed.Push(activityfeed.NewTitleEvent(
		p.ID, p.DisplayName, titleName, titleIcon, priority,
	))
	go g.broadcastFeedEvent(evt)
}

// notifyFeedJackpot Jackpot 中獎時推送動態（DAY-112）
func (g *Game) notifyFeedJackpot(p *player.Player, levelName, levelIcon string, amount int) {
	evt := g.ActivityFeed.Push(activityfeed.NewJackpotEvent(
		p.ID, p.DisplayName, levelName, levelIcon, amount,
	))
	go g.broadcastFeedEvent(evt)
}

// notifyFeedBossKill BOSS 擊殺時推送動態（DAY-112）
func (g *Game) notifyFeedBossKill(p *player.Player, reward int) {
	evt := g.ActivityFeed.Push(activityfeed.NewBossKillEvent(
		p.ID, p.DisplayName, reward,
	))
	go g.broadcastFeedEvent(evt)
}

// notifyFeedMegaWin 超大獎時推送動態（≥50x，DAY-112）
func (g *Game) notifyFeedMegaWin(p *player.Player, multiplier float64, reward int) {
	if multiplier < 50 {
		return
	}
	evt := g.ActivityFeed.Push(activityfeed.NewMegaWinEvent(
		p.ID, p.DisplayName, multiplier, reward,
	))
	go g.broadcastFeedEvent(evt)
}

// notifyFeedStreakRecord 連擊記錄時推送動態（≥20，DAY-112）
func (g *Game) notifyFeedStreakRecord(p *player.Player, streak int, levelName string) {
	if streak < 20 {
		return
	}
	evt := g.ActivityFeed.Push(activityfeed.NewStreakRecordEvent(
		p.ID, p.DisplayName, streak, levelName,
	))
	go g.broadcastFeedEvent(evt)
}

// notifyFeedHallOfFame 名人堂新記錄時推送動態（DAY-112）
func (g *Game) notifyFeedHallOfFame(p *player.Player, recordType, recordLabel string) {
	evt := g.ActivityFeed.Push(activityfeed.NewHallOfFameEvent(
		p.ID, p.DisplayName, recordType, recordLabel,
	))
	go g.broadcastFeedEvent(evt)
}

// notifyFeedSeasonLevel 賽季升級時推送動態（DAY-112）
func (g *Game) notifyFeedSeasonLevel(p *player.Player, level int) {
	// 只廣播重要等級（5的倍數）
	if level%5 != 0 {
		return
	}
	evt := g.ActivityFeed.Push(activityfeed.NewSeasonLevelEvent(
		p.ID, p.DisplayName, level,
	))
	go g.broadcastFeedEvent(evt)
}

// notifyFeedMilestone 登入里程碑時推送動態（DAY-112）
func (g *Game) notifyFeedMilestone(p *player.Player, days int, milestoneName string) {
	evt := g.ActivityFeed.Push(activityfeed.NewMilestoneEvent(
		p.ID, p.DisplayName, days, milestoneName,
	))
	go g.broadcastFeedEvent(evt)
}

// ---- 工具函數 ----

// feedEventToPayload 將 FeedEvent 轉換為 ws.ActivityFeedEventPayload
func feedEventToPayload(evt *activityfeed.FeedEvent) ws.ActivityFeedEventPayload {
	return ws.ActivityFeedEventPayload{
		ID:          evt.ID,
		EventType:   string(evt.EventType),
		PlayerID:    evt.PlayerID,
		DisplayName: evt.DisplayName,
		Icon:        evt.Icon,
		Title:       evt.Title,
		Detail:      evt.Detail,
		Rarity:      string(evt.Rarity),
		Timestamp:   evt.Timestamp,
	}
}
