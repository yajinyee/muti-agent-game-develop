// announce_handler.go — 全服公告系統 handler（DAY-097）
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/ws"
)

// broadcastAnnouncement 廣播全服公告
func (g *Game) broadcastAnnouncement(ann announce.Announcement) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: ws.AnnouncementPayload{
			ID:          ann.ID,
			EventType:   string(ann.EventType),
			Priority:    int(ann.Priority),
			Title:       ann.Title,
			Message:     ann.Message,
			PlayerName:  ann.PlayerName,
			Amount:      ann.Amount,
			Icon:        ann.Icon,
			Color:       ann.Color,
			Duration:    ann.Duration,
			CreatedAtMs: ann.CreatedAtMs,
		},
	})
	log.Printf("[Announce] [%s] %s", ann.EventType, ann.Message)
}

// announceJackpotWin 公告 Jackpot 中獎
func (g *Game) announceJackpotWin(playerName string, level string, levelName string, amount int) {
	eventType := announce.EventJackpotWin
	if level == "grand" {
		eventType = announce.EventGrandJackpot
	}
	ann := g.Announce.Create(eventType, playerName, amount, map[string]string{
		"level":      level,
		"level_name": levelName,
	})
	g.broadcastAnnouncement(ann)
}

// announceBigWin 公告大獎（≥50x 為 Mega，≥20x 為 Big）
func (g *Game) announceBigWin(playerName string, multiplier float64, reward int) {
	eventType := announce.EventBigWin
	if multiplier >= 100 {
		eventType = announce.EventMegaWin
	}
	ann := g.Announce.Create(eventType, playerName, reward, map[string]string{
		"multiplier": formatMult(multiplier),
	})
	g.broadcastAnnouncement(ann)
}

// announceBossKill 公告 BOSS 擊殺
func (g *Game) announceBossKill(playerName string, bossName string, reward int) {
	ann := g.Announce.Create(announce.EventBossKill, playerName, reward, map[string]string{
		"boss_name": bossName,
	})
	g.broadcastAnnouncement(ann)
}

// announceStreakRecord 公告連擊記錄（≥20 才公告）
func (g *Game) announceStreakRecord(playerName string, streak int) {
	if streak < 20 {
		return
	}
	ann := g.Announce.Create(announce.EventStreakRecord, playerName, streak, nil)
	g.broadcastAnnouncement(ann)
}

// announcePlayerJoin 公告玩家加入
func (g *Game) announcePlayerJoin(playerName string) {
	ann := g.Announce.Create(announce.EventPlayerJoin, playerName, 0, nil)
	g.broadcastAnnouncement(ann)
}

// announcePlayerLeave 公告玩家離開
func (g *Game) announcePlayerLeave(playerName string) {
	ann := g.Announce.Create(announce.EventPlayerLeave, playerName, 0, nil)
	g.broadcastAnnouncement(ann)
}

// announceWeatherChange 公告天氣變化
func (g *Game) announceWeatherChange(weatherName string) {
	ann := g.Announce.Create(announce.EventWeatherChange, "", 0, map[string]string{
		"weather_name": weatherName,
	})
	g.broadcastAnnouncement(ann)
}

// announceEventStart 公告限時活動開始
func (g *Game) announceEventStart(eventName string) {
	ann := g.Announce.Create(announce.EventEventStart, "", 0, map[string]string{
		"event_name": eventName,
	})
	g.broadcastAnnouncement(ann)
}

// announceBossWarning 公告 BOSS 即將出現
func (g *Game) announceBossWarning() {
	ann := g.Announce.Create(announce.EventBossWarning, "", 0, nil)
	g.broadcastAnnouncement(ann)
}

// formatMult 格式化倍率顯示
func formatMult(mult float64) string {
	if mult == float64(int(mult)) {
		return fmt.Sprintf("%dx", int(mult))
	}
	return fmt.Sprintf("%.1fx", mult)
}
