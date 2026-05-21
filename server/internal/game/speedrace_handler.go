// speedrace_handler.go — 全服競速獵殺系統 handler（DAY-136）
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/speedrace"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// tryStartSpeedRace 嘗試對新生成的目標啟動競速獵殺
// 由 spawnTarget 呼叫（高倍率目標生成時）
func (g *Game) tryStartSpeedRace(instanceID, defID, name string, mult float64) {
	if g.SpeedRace == nil {
		return
	}
	sess := g.SpeedRace.StartRace(instanceID, defID, name, mult)
	if sess == nil {
		return
	}

	log.Printf("[SpeedRace] Race started: target=%s (%s) mult=%.0fx, duration=%.0fs",
		name, instanceID, mult, sess.EndAt.Sub(sess.StartAt).Seconds())

	// 全服廣播競速開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgSpeedRaceStart,
		Payload: ws.SpeedRaceStartPayload{
			TargetInstanceID: instanceID,
			TargetDefID:      defID,
			TargetName:       name,
			TargetMult:       mult,
			SecondsLeft:      sess.EndAt.Sub(sess.StartAt).Seconds(),
			BonusMult:        3.0,
			Message:          fmt.Sprintf("🏆 競速獵殺！搶先擊破【%s】(×%.0f) 獲得 3x 獎勵！", name, mult),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventSpeedRace, "", 0, map[string]string{
		"target_name": name,
		"mult":        fmt.Sprintf("%.0f", mult),
	})
	g.broadcastAnnouncement(ann)
}

// notifySpeedRaceKill 處理競速目標被擊破
// 由 handleKill 呼叫，回傳競速獎勵倍率（1.0 = 無加成）
func (g *Game) notifySpeedRaceKill(p *player.Player, instanceID string) float64 {
	if g.SpeedRace == nil {
		return 1.0
	}

	rank, bonusMult, isRaceTarget := g.SpeedRace.RecordKill(instanceID, p.ID, p.DisplayName)
	if !isRaceTarget {
		return 1.0
	}
	if rank == 0 {
		return 1.0 // 競速已結束（第一名已出現）
	}

	snap := g.SpeedRace.GetSnapshot()

	// 發送個人競速結果
	rankIcon := ""
	switch rank {
	case 1:
		rankIcon = "🥇"
	case 2:
		rankIcon = "🥈"
	case 3:
		rankIcon = "🥉"
	}

	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgSpeedRaceResult,
		Payload: ws.SpeedRaceResultPayload{
			PlayerID:    p.ID,
			DisplayName: p.DisplayName,
			Rank:        rank,
			BonusMult:   bonusMult,
			RankIcon:    rankIcon,
			Message:     fmt.Sprintf("%s 競速第 %d 名！獎勵 ×%.1f", rankIcon, rank, bonusMult),
		},
	})

	// 第一名：全服廣播 + 公告
	if rank == 1 {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgSpeedRaceEnd,
			Payload: ws.SpeedRaceEndPayload{
				WinnerID:   p.ID,
				WinnerName: p.DisplayName,
				TargetName: snap.TargetName,
				TargetMult: snap.TargetMult,
				BonusMult:  bonusMult,
				Message:    fmt.Sprintf("🥇 %s 搶先擊破【%s】！獲得 %.1fx 獎勵！", p.DisplayName, snap.TargetName, bonusMult),
			},
		})

		ann := g.Announce.Create(announce.EventSpeedRaceWin, p.DisplayName, 0, map[string]string{
			"target_name": snap.TargetName,
			"bonus_mult":  fmt.Sprintf("%.1f", bonusMult),
		})
		g.broadcastAnnouncement(ann)

		log.Printf("[SpeedRace] Winner: player=%s target=%s mult=%.0fx bonus=%.1fx",
			p.DisplayName, snap.TargetName, snap.TargetMult, bonusMult)
	}

	return bonusMult
}

// tickSpeedRace 競速超時檢查（由 gameLoop 每次 update 呼叫）
func (g *Game) tickSpeedRace() {
	if g.SpeedRace == nil {
		return
	}

	snap := g.SpeedRace.GetSnapshot()
	if !snap.IsActive {
		return
	}

	if g.SpeedRace.CheckExpiry() {
		// 競速超時，廣播取消
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgSpeedRaceCancel,
			Payload: ws.SpeedRaceCancelPayload{
				TargetInstanceID: snap.TargetInstanceID,
				TargetName:       snap.TargetName,
				Message:          fmt.Sprintf("⏰ 競速獵殺超時！【%s】無人搶先擊破。", snap.TargetName),
			},
		})
		log.Printf("[SpeedRace] Race expired: target=%s", snap.TargetName)
	}
}

// cancelSpeedRaceIfTarget 目標消失時取消競速（由 updateNormalPlay 清除過期目標後呼叫）
func (g *Game) cancelSpeedRaceIfTarget(instanceID string) {
	if g.SpeedRace == nil {
		return
	}
	if g.SpeedRace.CancelRace(instanceID) {
		snap := g.SpeedRace.GetSnapshot()
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgSpeedRaceCancel,
			Payload: ws.SpeedRaceCancelPayload{
				TargetInstanceID: instanceID,
				TargetName:       snap.TargetName,
				Message:          "競速目標已消失，競速取消。",
			},
		})
	}
}

// sendSpeedRaceStatus 登入時發送競速狀態
func (g *Game) sendSpeedRaceStatus(p *player.Player) {
	if g.SpeedRace == nil {
		return
	}
	snap := g.SpeedRace.GetSnapshot()
	if !snap.IsActive {
		return
	}
	// 告知玩家當前有進行中的競速
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgSpeedRaceStart,
		Payload: ws.SpeedRaceStartPayload{
			TargetInstanceID: snap.TargetInstanceID,
			TargetDefID:      snap.TargetDefID,
			TargetName:       snap.TargetName,
			TargetMult:       snap.TargetMult,
			SecondsLeft:      snap.SecondsLeft,
			BonusMult:        snap.BonusMult,
			Message:          fmt.Sprintf("🏆 競速進行中！搶先擊破【%s】(×%.0f) 獲得 3x 獎勵！", snap.TargetName, snap.TargetMult),
		},
	})
}

// getSpeedRaceSnapshot 取得競速快照（供 HTTP 端點使用）
func (g *Game) getSpeedRaceSnapshot() *speedrace.RaceSnapshot {
	if g.SpeedRace == nil {
		return &speedrace.RaceSnapshot{IsActive: false}
	}
	return g.SpeedRace.GetSnapshot()
}
