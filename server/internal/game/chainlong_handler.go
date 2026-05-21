// chainlong_handler.go — 千龍王強化輪盤系統 handler（DAY-148）
// 業界依據：Royal Fishing JILI 2026「ChainLong King — capture this golden dragon to trigger
// the dual-ring roulette. The ChainLong King itself can award up to 1000X mega wins.」
// 設計：T112 千龍王（150-1000x）擊破後觸發強化版雙環輪盤
// 內環（5x/10x/20x/30x/50x）× 外環（2x/3x/5x/7x/10x/20x）= 最高 1000x
// 比普通雙環輪盤（最高 150x）強 6.7 倍，是全遊戲最高倍率的個人機制
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/chainlongwheel"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// isChainLongKing 判斷是否為千龍王目標
func isChainLongKing(defID string) bool {
	return defID == "T112"
}

// tryChainLongWheel 擊破千龍王後觸發強化輪盤（由 handleKill 呼叫）
func (g *Game) tryChainLongWheel(p *player.Player, instanceID string, baseReward int, targetMult float64) {
	if g.ChainLongWheel == nil {
		return
	}

	// 千龍王必定觸發輪盤（不需要機率判斷，擊破即觸發）
	if !g.ChainLongWheel.CanTrigger(p.ID) {
		log.Printf("[ChainLongWheel] player=%s cannot trigger (cooldown or active session)", p.ID)
		return
	}

	// 開始 session（結果預先決定）
	s := g.ChainLongWheel.StartSession(p.ID, targetMult, baseReward)
	if s == nil {
		return
	}

	log.Printf("[ChainLongWheel] player=%s triggered! targetMult=%.0fx, baseReward=%d, inner=%.0fx, outer=%.0fx, combined=%.0fx",
		p.ID, targetMult, baseReward, s.InnerResult, s.OuterResult, s.InnerResult*s.OuterResult)

	// 廣播千龍王輪盤開始（全服可見）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgChainLongWheelStart,
		Payload: ws.ChainLongWheelStartPayload{
			InstanceID:  instanceID,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			TargetMult:  targetMult,
			BaseReward:  baseReward,
			InnerSlots:  chainlongwheel.InnerRing,
			OuterSlots:  chainlongwheel.OuterRing,
			SpinSecs:    chainlongwheel.SpinDuration,
			Message:     fmt.Sprintf("🐉 %s 擊破了千龍王！觸發強化輪盤！最高 1000x！", p.DisplayName),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, baseReward, map[string]string{
		"message": fmt.Sprintf("🐉 千龍王出現！%s 觸發強化輪盤！最高 1000x！", p.DisplayName),
	})
	g.broadcastAnnouncement(ann)
}

// handleChainLongWheelStop 處理玩家停止千龍王輪盤（Client → Server）
func (g *Game) handleChainLongWheelStop(playerID string) {
	if g.ChainLongWheel == nil {
		return
	}

	g.mu.RLock()
	p, ok := g.Players[playerID]
	g.mu.RUnlock()
	if !ok {
		return
	}

	s := g.ChainLongWheel.StopSession(playerID)
	if s == nil {
		return
	}

	g.processChainLongWheelResult(p, s)
}

// processChainLongWheelResult 處理千龍王輪盤結果，發放獎勵並廣播
func (g *Game) processChainLongWheelResult(p *player.Player, s *chainlongwheel.Session) {
	bonusReward := s.BonusReward()
	combined := s.FinalMultiplier()
	isMegaWin := combined >= 200.0

	// 發放額外獎勵
	p.AddCoins(bonusReward)

	log.Printf("[ChainLongWheel] result: player=%s, inner=%.0fx, outer=%.0fx, combined=%.0fx, bonus=%d",
		p.ID, s.InnerResult, s.OuterResult, combined, bonusReward)

	// 個人通知（完整結果）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgChainLongWheelResult,
		Payload: ws.ChainLongWheelResultPayload{
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			TargetMult:  s.TargetMult,
			InnerResult: s.InnerResult,
			OuterResult: s.OuterResult,
			Combined:    combined,
			BaseReward:  s.BaseReward,
			BonusReward: bonusReward,
			NewBalance:  p.GetCoins(),
			IsMegaWin:   isMegaWin,
			IsPersonal:  true,
			Message:     fmt.Sprintf("🐉 千龍王輪盤結果：%.0fx × %.0fx = %.0fx！額外獎勵 %d 金幣！", s.InnerResult, s.OuterResult, combined, bonusReward),
		},
	})

	// 全服廣播（讓其他玩家看到）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgChainLongWheelResult,
		Payload: ws.ChainLongWheelResultPayload{
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			TargetMult:  s.TargetMult,
			InnerResult: s.InnerResult,
			OuterResult: s.OuterResult,
			Combined:    combined,
			BaseReward:  s.BaseReward,
			BonusReward: bonusReward,
			IsMegaWin:   isMegaWin,
			IsPersonal:  false,
			Message:     fmt.Sprintf("🐉 %s 千龍王輪盤：%.0fx × %.0fx = %.0fx！", p.DisplayName, s.InnerResult, s.OuterResult, combined),
		},
	})

	// 大獎（≥200x）全服公告
	if isMegaWin {
		ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, bonusReward, map[string]string{
			"message": fmt.Sprintf("🐉🌟 %s 千龍王輪盤大獎！%.0fx × %.0fx = %.0fx！獎勵 %d 金幣！", p.DisplayName, s.InnerResult, s.OuterResult, combined, bonusReward),
		})
		g.broadcastAnnouncement(ann)
		// 動態牆：千龍王大獎
		go g.notifyFeedMegaWin(p, combined, bonusReward)
	}

	// 更新玩家狀態
	g.sendPlayerUpdate(p)
}

// tickChainLongWheel 每秒檢查千龍王輪盤超時（由 gameLoop 呼叫）
func (g *Game) tickChainLongWheel() {
	if g.ChainLongWheel == nil {
		return
	}

	expired := g.ChainLongWheel.TickAutoStop()
	for _, s := range expired {
		g.mu.RLock()
		p, ok := g.Players[s.PlayerID]
		g.mu.RUnlock()
		if !ok {
			continue
		}
		log.Printf("[ChainLongWheel] auto-stop: player=%s", s.PlayerID)
		g.processChainLongWheelResult(p, s)
	}
}

// sendChainLongWheelStatus 發送千龍王輪盤冷卻狀態給玩家（登入時呼叫）
func (g *Game) sendChainLongWheelStatus(p *player.Player) {
	if g.ChainLongWheel == nil {
		return
	}

	cooldown := g.ChainLongWheel.GetCooldownLeft(p.ID)
	hasActive := g.ChainLongWheel.HasActiveSession(p.ID)

	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgChainLongWheelStatus,
		Payload: ws.ChainLongWheelStatusPayload{
			CooldownLeft: cooldown,
			HasActive:    hasActive,
		},
	})
}
