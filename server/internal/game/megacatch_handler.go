// megacatch_handler.go — 全服 Mega Catch 事件系統 handler（DAY-140）
// 業界依據：Ocean King 系列「Mega Catch」— 全場高倍率目標湧現 + 獎勵翻倍
// BOSS 擊殺後 60% 機率觸發，或每分鐘 5% 機率隨機觸發
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/megacatch"
	"digital-twin/server/internal/ws"
)

// tryMegaCatchBossKill BOSS 擊殺後嘗試觸發 Mega Catch
// 由 handleBossKill 呼叫
func (g *Game) tryMegaCatchBossKill() {
	s := g.MegaCatch.TryTriggerBossKill()
	if s == nil {
		return
	}
	g.broadcastMegaCatchStart(s)
}

// tryMegaCatchRandom 每分鐘隨機嘗試觸發 Mega Catch
// 由 updateNormalPlay 每分鐘呼叫
func (g *Game) tryMegaCatchRandom() {
	s := g.MegaCatch.TryTriggerRandom()
	if s == nil {
		return
	}
	g.broadcastMegaCatchStart(s)
}

// broadcastMegaCatchStart 廣播 Mega Catch 開始
func (g *Game) broadcastMegaCatchStart(s *megacatch.Session) {
	log.Printf("[MegaCatch] Event started: %s (reward×%.1f, spawn+%.0f%%, %.0fs)",
		s.Tier.Name, s.Tier.RewardBoost, s.Tier.SpawnMultBoost*100, s.Tier.Duration)

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgMegaCatchStart,
		Payload: ws.MegaCatchStartPayload{
			TierName:      s.Tier.Name,
			TierIcon:      s.Tier.Icon,
			TierColor:     s.Tier.Color,
			RewardBoost:   s.Tier.RewardBoost,
			SpawnBoost:    s.Tier.SpawnMultBoost,
			Duration:      s.Tier.Duration,
			SecondsLeft:   s.SecondsLeft(),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventMegaCatch, "", 0, map[string]string{
		"tier_name":    s.Tier.Name,
		"tier_icon":    s.Tier.Icon,
		"reward_boost": fmt.Sprintf("%.0f", s.Tier.RewardBoost),
		"duration":     fmt.Sprintf("%.0f", s.Tier.Duration),
	})
	g.broadcastAnnouncement(ann)
}

// tickMegaCatch 每秒檢查 Mega Catch 過期（由 updateNormalPlay 呼叫）
func (g *Game) tickMegaCatch() {
	expired := g.MegaCatch.CheckExpiry()
	if expired == nil {
		return
	}

	log.Printf("[MegaCatch] Event ended: %s", expired.Tier.Name)

	// 全服廣播結束
	g.Hub.Broadcast(&ws.Message{
		Type:    ws.MsgMegaCatchEnd,
		Payload: ws.MegaCatchEndPayload{Message: "Mega Catch 結束，繼續加油！"},
	})
}

// getMegaCatchRewardBoost 取得當前 Mega Catch 獎勵倍率（供 handleKill 使用）
func (g *Game) getMegaCatchRewardBoost() float64 {
	return g.MegaCatch.GetRewardBoost()
}

// getMegaCatchSpawnBoost 取得當前 Mega Catch 稀有目標生成加成（供 spawnTarget 使用）
func (g *Game) getMegaCatchSpawnBoost() float64 {
	return g.MegaCatch.GetSpawnBoost()
}

// sendMegaCatchStatus 登入時發送 Mega Catch 狀態
func (g *Game) sendMegaCatchStatus(playerID string) {
	snap := g.MegaCatch.GetSnapshot()
	g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgMegaCatchStatus,
		Payload: ws.MegaCatchStatusPayload{
			IsActive:    snap.IsActive,
			TierName:    snap.TierName,
			TierIcon:    snap.TierIcon,
			TierColor:   snap.TierColor,
			RewardBoost: snap.RewardBoost,
			SpawnBoost:  snap.SpawnBoost,
			SecondsLeft: snap.SecondsLeft,
		},
	})
}
