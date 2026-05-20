// halloffame_handler.go — 名人堂系統 handler（DAY-110）
package game

import (
	"fmt"
	"log"
	"time"

	"digital-twin/server/internal/game/halloffame"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// betLevelToCharID 根據投注等級推算角色 ID（1=吉伊卡哇, 2=小八, 3=烏薩奇）
func betLevelToCharID(betLevel int) int {
	if betLevel <= 3 {
		return 1
	} else if betLevel <= 7 {
		return 2
	}
	return 3
}

// handleGetHallOfFame 處理查詢名人堂請求
func (g *Game) handleGetHallOfFame(p *player.Player) {
	snap := g.HallOfFame.GetAll()
	payload := buildHallOfFamePayload(snap)
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type:    ws.MsgHallOfFameUpdate,
		Payload: payload,
	}); err != nil {
		log.Printf("[HallOfFame] send error: %v", err)
	}
}

// notifyHallOfFameKill 在擊破目標後嘗試更新名人堂記錄
func (g *Game) notifyHallOfFameKill(p *player.Player, multiplier float64, reward int) {
	if p.Stats == nil {
		return
	}
	snap := p.Stats.Snapshot()

	// 最高倍率
	if multiplier > 0 {
		desc := fmt.Sprintf("%.0fx 倍率（LV%d）", multiplier, p.BetLevel)
		isNew, old := g.HallOfFame.TryUpdate(
			p.ID, p.DisplayName,
			halloffame.RecordBestMultiplier,
			multiplier, desc,
			p.BetLevel, betLevelToCharID(p.BetLevel),
		)
		if isNew {
			g.broadcastHallOfFameNewRecord(halloffame.RecordBestMultiplier, p, multiplier, desc, old)
		}
	}

	// 最高金幣
	if p.Coins > 0 {
		desc := fmt.Sprintf("%d 金幣", p.Coins)
		isNew, old := g.HallOfFame.TryUpdate(
			p.ID, p.DisplayName,
			halloffame.RecordMaxCoins,
			float64(p.Coins), desc,
			p.BetLevel, betLevelToCharID(p.BetLevel),
		)
		if isNew {
			g.broadcastHallOfFameNewRecord(halloffame.RecordMaxCoins, p, float64(p.Coins), desc, old)
		}
	}

	// 最高 RTP（需 >= 100 次攻擊）
	if snap.TotalShots >= 100 && snap.RTP > 0 {
		desc := fmt.Sprintf("RTP %.0f%%（%d 次攻擊）", snap.RTP*100, snap.TotalShots)
		isNew, old := g.HallOfFame.TryUpdate(
			p.ID, p.DisplayName,
			halloffame.RecordBestRTP,
			snap.RTP, desc,
			p.BetLevel, betLevelToCharID(p.BetLevel),
		)
		if isNew {
			g.broadcastHallOfFameNewRecord(halloffame.RecordBestRTP, p, snap.RTP, desc, old)
		}
	}
}

// notifyHallOfFameStreak 在連擊更新後嘗試更新名人堂
func (g *Game) notifyHallOfFameStreak(p *player.Player, streak int) {
	if streak < 5 {
		return
	}
	desc := fmt.Sprintf("%d 連擊（LV%d）", streak, p.BetLevel)
	isNew, old := g.HallOfFame.TryUpdate(
		p.ID, p.DisplayName,
		halloffame.RecordBestStreak,
		float64(streak), desc,
		p.BetLevel, betLevelToCharID(p.BetLevel),
	)
	if isNew {
		g.broadcastHallOfFameNewRecord(halloffame.RecordBestStreak, p, float64(streak), desc, old)
	}
}

// notifyHallOfFameBonus 在 Bonus 結算後嘗試更新名人堂
func (g *Game) notifyHallOfFameBonus(p *player.Player, reward int) {
	if reward <= 0 {
		return
	}
	desc := fmt.Sprintf("Bonus 獎勵 %d 金幣（LV%d）", reward, p.BetLevel)
	isNew, old := g.HallOfFame.TryUpdate(
		p.ID, p.DisplayName,
		halloffame.RecordBestBonusReward,
		float64(reward), desc,
		p.BetLevel, betLevelToCharID(p.BetLevel),
	)
	if isNew {
		g.broadcastHallOfFameNewRecord(halloffame.RecordBestBonusReward, p, float64(reward), desc, old)
	}
}

// notifyHallOfFameJackpot 在 Jackpot 中獎後嘗試更新名人堂
func (g *Game) notifyHallOfFameJackpot(p *player.Player, level string, amount int) {
	if p.Stats == nil {
		return
	}
	snap := p.Stats.Snapshot()

	// Grand Jackpot 記錄
	if level == "grand" {
		desc := fmt.Sprintf("Grand Jackpot %d 金幣", amount)
		isNew, old := g.HallOfFame.TryUpdate(
			p.ID, p.DisplayName,
			halloffame.RecordGrandJackpot,
			float64(amount), desc,
			p.BetLevel, betLevelToCharID(p.BetLevel),
		)
		if isNew {
			g.broadcastHallOfFameNewRecord(halloffame.RecordGrandJackpot, p, float64(amount), desc, old)
		}
	}

	// 最多 Jackpot 次數
	totalJackpots := snap.JackpotWins
	if totalJackpots > 0 {
		desc := fmt.Sprintf("累計 %d 次 Jackpot", totalJackpots)
		isNew, old := g.HallOfFame.TryUpdate(
			p.ID, p.DisplayName,
			halloffame.RecordMostJackpots,
			float64(totalJackpots), desc,
			p.BetLevel, betLevelToCharID(p.BetLevel),
		)
		if isNew {
			g.broadcastHallOfFameNewRecord(halloffame.RecordMostJackpots, p, float64(totalJackpots), desc, old)
		}
	}
}

// notifyHallOfFameBossKill 在 BOSS 擊殺後嘗試更新名人堂
func (g *Game) notifyHallOfFameBossKill(p *player.Player) {
	if p.Stats == nil {
		return
	}
	snap := p.Stats.Snapshot()
	bossKills := snap.TotalBossKills
	if bossKills <= 0 {
		return
	}
	desc := fmt.Sprintf("擊殺 %d 隻 BOSS", bossKills)
	isNew, old := g.HallOfFame.TryUpdate(
		p.ID, p.DisplayName,
		halloffame.RecordBossKills,
		float64(bossKills), desc,
		p.BetLevel, betLevelToCharID(p.BetLevel),
	)
	if isNew {
		g.broadcastHallOfFameNewRecord(halloffame.RecordBossKills, p, float64(bossKills), desc, old)
	}
}

// broadcastHallOfFameNewRecord 廣播新記錄誕生
func (g *Game) broadcastHallOfFameNewRecord(
	rt halloffame.RecordType,
	p *player.Player,
	value float64,
	desc string,
	old *halloffame.HallEntry,
) {
	entry := ws.HallEntryPayload{
		PlayerID:    p.ID,
		DisplayName: p.DisplayName,
		RecordType:  string(rt),
		RecordLabel: halloffame.RecordTypeLabel(rt),
		RecordIcon:  halloffame.RecordTypeIcon(rt),
		Value:       value,
		Description: desc,
		AchievedAt:  time.Now().UnixMilli(),
		BetLevel:    p.BetLevel,
		CharacterID: betLevelToCharID(p.BetLevel),
	}

	oldHolder := ""
	oldValue := 0.0
	isFirst := true
	if old != nil {
		oldHolder = old.DisplayName
		oldValue = old.Value
		isFirst = false
	}

	payload := ws.HallOfFameNewRecordPayload{
		Entry:       entry,
		OldHolder:   oldHolder,
		OldValue:    oldValue,
		IsFirstTime: isFirst,
	}

	g.Hub.Broadcast(&ws.Message{
		Type:    ws.MsgHallOfFameNewRecord,
		Payload: payload,
	})

	log.Printf("[HallOfFame] NEW RECORD %s %s: %.2f by %s",
		halloffame.RecordTypeIcon(rt), halloffame.RecordTypeLabel(rt), value, p.DisplayName)

	// 動態牆：名人堂新記錄（DAY-112）
	recordLabel := halloffame.RecordTypeIcon(rt) + " " + halloffame.RecordTypeLabel(rt) + ": " + desc
	go g.notifyFeedHallOfFame(p, string(rt), recordLabel)
}

// buildHallOfFamePayload 將名人堂快照轉換為 Payload
func buildHallOfFamePayload(snap halloffame.HallSnapshot) ws.HallOfFameUpdatePayload {
	var entries []ws.HallEntryPayload
	for rt, e := range snap.Records {
		entries = append(entries, ws.HallEntryPayload{
			PlayerID:    e.PlayerID,
			DisplayName: e.DisplayName,
			RecordType:  string(rt),
			RecordLabel: halloffame.RecordTypeLabel(rt),
			RecordIcon:  halloffame.RecordTypeIcon(rt),
			Value:       e.Value,
			Description: e.Description,
			AchievedAt:  e.AchievedAt.UnixMilli(),
			BetLevel:    e.BetLevel,
			CharacterID: e.CharacterID,
		})
	}
	return ws.HallOfFameUpdatePayload{
		Records:   entries,
		UpdatedAt: snap.UpdatedAt,
	}
}

// GetHallOfFameSnapshot 取得名人堂快照（供 HTTP 端點使用）
func (g *Game) GetHallOfFameSnapshot() interface{} {
	snap := g.HallOfFame.GetAll()
	return buildHallOfFamePayload(snap)
}
