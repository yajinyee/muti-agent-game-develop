// dualroulette_handler.go — 雙環輪盤系統 handler（DAY-139）
// 業界依據：Royal Fishing JILI 2026 ChainLong King Dual-Ring Roulette
// 擊破高倍率目標後觸發，內外圈相乘最高 150x，製造「技巧感」
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/dualroulette"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// tryDualRoulette 擊破高倍率目標後嘗試觸發雙環輪盤
// 由 handleKill 呼叫（在倍率疊加之後）
func (g *Game) tryDualRoulette(p *player.Player, targetMult float64, baseReward int) {
	if !g.DualRoulette.CanTrigger(p.ID, targetMult) {
		return
	}

	s := g.DualRoulette.StartSession(p.ID, targetMult, baseReward)
	if s == nil {
		return
	}

	log.Printf("[DualRoulette] player=%s triggered, targetMult=%.1f, baseReward=%d, inner=%.1f, outer=%.1f",
		p.ID, targetMult, baseReward, s.InnerResult, s.OuterResult)

	// 廣播輪盤開始（只傳給觸發者）
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgDualRouletteStart,
		Payload: ws.DualRouletteStartPayload{
			PlayerID:     p.ID,
			TargetMult:   targetMult,
			BaseReward:   baseReward,
			SpinDuration: 3.0,
			InnerRing:    dualroulette.InnerRing,
			OuterRing:    dualroulette.OuterRing,
		},
	}); err != nil {
		log.Printf("[DualRoulette] send start error: %v", err)
	}
}

// handleDualRouletteStop 玩家點擊停止輪盤
func (g *Game) handleDualRouletteStop(playerID string) {
	g.mu.RLock()
	p, exists := g.Players[playerID]
	g.mu.RUnlock()
	if !exists {
		return
	}

	s := g.DualRoulette.StopSession(playerID)
	if s == nil {
		return
	}

	g.processDualRouletteResult(p, s)
}

// processDualRouletteResult 處理輪盤結果（停止或超時）
func (g *Game) processDualRouletteResult(p *player.Player, s *dualroulette.Session) {
	snap := g.DualRoulette.GetSnapshot(s)
	bonusReward := snap.BonusReward

	log.Printf("[DualRoulette] player=%s result: inner=%.1f × outer=%.1f = %.1f, bonus=%d",
		p.ID, snap.InnerResult, snap.OuterResult, snap.Combined, bonusReward)

	// 發放獎勵
	if bonusReward > 0 {
		p.AddReward(bonusReward)
	}

	// 發送結果給玩家
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgDualRouletteResult,
		Payload: ws.DualRouletteResultPayload{
			PlayerID:    p.ID,
			InnerResult: snap.InnerResult,
			OuterResult: snap.OuterResult,
			Combined:    snap.Combined,
			BonusReward: bonusReward,
			NewBalance:  p.Coins,
		},
	}); err != nil {
		log.Printf("[DualRoulette] send result error: %v", err)
	}

	// 全服公告：高倍率組合（≥50x）
	if snap.Combined >= 50.0 {
		ann := g.Announce.Create(announce.EventDualRoulette, p.DisplayName, snap.BonusReward, map[string]string{
			"combined":     fmt.Sprintf("%.0f", snap.Combined),
			"bonus_reward": fmt.Sprintf("%d", snap.BonusReward),
		})
		g.broadcastAnnouncement(ann)
	}
}

// tickDualRoulette 每秒檢查超時的輪盤 session（由 gameLoop 呼叫）
func (g *Game) tickDualRoulette() {
	expired := g.DualRoulette.TickAutoStop()
	for _, s := range expired {
		g.mu.RLock()
		p, exists := g.Players[s.PlayerID]
		g.mu.RUnlock()
		if !exists {
			continue
		}
		log.Printf("[DualRoulette] auto-stop player=%s", s.PlayerID)
		g.processDualRouletteResult(p, s)
	}
}

// sendDualRouletteStatus 登入時發送輪盤狀態（冷卻剩餘）
func (g *Game) sendDualRouletteStatus(playerID string) {
	cdLeft := g.DualRoulette.GetCooldownLeft(playerID)
	g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgDualRouletteStatus,
		Payload: ws.DualRouletteStatusPayload{
			CooldownLeft: cdLeft,
		},
	})
}
