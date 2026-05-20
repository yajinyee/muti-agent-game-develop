// roulette_handler.go — 雙層倍率輪盤系統 handler（DAY-113）
// 參考 JILI Royal Fishing 的 ChainLong King 雙層輪盤機制
// 觸發條件：擊破 B001 BOSS / T103 流星(5%) / T104 金草(8%) / T105 金幣魚(10%)
// 機制：內圈(8格,1-10x) × 外圈(12格,1-100x) = 最終倍率（最高 1000x）
package game

import (
	"log"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/roulette"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// notifyRouletteKill 在擊破目標後判斷是否觸發雙層輪盤（由 handleKill 呼叫）
// 若觸發，廣播輪盤開始，執行旋轉，廣播結果，回傳額外獎勵金額
func (g *Game) notifyRouletteKill(p *player.Player, defID string, baseReward int) int {
	if g.Roulette == nil {
		return 0
	}

	betLevel := p.BetLevel
	if !g.Roulette.ShouldTrigger(defID, betLevel) {
		return 0
	}

	// 開始 session
	session := g.Roulette.StartSession(p.ID, defID, baseReward)

	// 取得目標名稱
	targetName := defID
	if def, ok := data.Targets[defID]; ok {
		targetName = def.Name
	}

	// 建立格子列表（供 Client 顯示）
	innerSegs := make([]ws.RouletteSegmentPayload, len(roulette.InnerSegments))
	for i, s := range roulette.InnerSegments {
		innerSegs[i] = ws.RouletteSegmentPayload{
			Multiplier: s.Multiplier,
			Label:      s.Label,
			Color:      s.Color,
		}
	}
	outerSegs := make([]ws.RouletteSegmentPayload, len(roulette.OuterSegments))
	for i, s := range roulette.OuterSegments {
		outerSegs[i] = ws.RouletteSegmentPayload{
			Multiplier: s.Multiplier,
			Label:      s.Label,
			Color:      s.Color,
		}
	}

	// 廣播輪盤開始（所有玩家都能看到）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRouletteStart,
		Payload: ws.RouletteStartPayload{
			SessionID:      session.ID,
			PlayerID:       p.ID,
			PlayerName:     p.DisplayName,
			TargetDefID:    defID,
			TargetName:     targetName,
			BaseReward:     baseReward,
			InnerSegments:  innerSegs,
			OuterSegments:  outerSegs,
			SpinDurationMs: 3000, // 3 秒旋轉動畫
		},
	})

	log.Printf("[Roulette] player=%s triggered roulette on %s, base=%d",
		p.ID, defID, baseReward)

	// 執行旋轉（立即計算結果，Client 端做動畫）
	result, ok := g.Roulette.Resolve(p.ID)
	if !ok {
		log.Printf("[Roulette] resolve failed for player=%s", p.ID)
		return 0
	}

	// 計算額外獎勵（最終獎勵 - 基礎獎勵）
	extraReward := result.FinalReward - baseReward
	if extraReward < 0 {
		extraReward = 0
	}

	// 發放額外獎勵
	if extraReward > 0 {
		p.AddCoins(extraReward)
	}

	// 廣播結果（所有玩家都能看到，但 new_balance 只對觸發玩家有意義）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRouletteResult,
		Payload: ws.RouletteResultPayload{
			SessionID:  session.ID,
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Inner: ws.RouletteSpinPayload{
				SegmentIndex: result.Inner.SegmentIndex,
				Multiplier:   result.Inner.Multiplier,
				Label:        result.Inner.Label,
				Color:        result.Inner.Color,
			},
			Outer: ws.RouletteSpinPayload{
				SegmentIndex: result.Outer.SegmentIndex,
				Multiplier:   result.Outer.Multiplier,
				Label:        result.Outer.Label,
				Color:        result.Outer.Color,
			},
			FinalMult:   result.FinalMult,
			BaseReward:  baseReward,
			FinalReward: result.FinalReward,
			NewBalance:  p.GetCoins(),
			IsJackpot:   result.IsJackpot,
			IsMegaWin:   result.IsMegaWin,
		},
	})

	log.Printf("[Roulette] player=%s result: inner=%.0fx outer=%.0fx final=%.0fx reward=%d (extra=%d)",
		p.ID, result.Inner.Multiplier, result.Outer.Multiplier, result.FinalMult, result.FinalReward, extraReward)

	// 動態牆：超大獎通知（≥100x）
	if result.IsMegaWin {
		go g.notifyFeedMegaWin(p, result.FinalMult, result.FinalReward)
	}

	// 名人堂：最高倍率記錄
	go g.notifyHallOfFameKill(p, result.FinalMult, result.FinalReward)

	// 玩家統計：記錄擊破
	g.notifyStatsKill(p, result.FinalMult, result.FinalReward)

	return extraReward
}
