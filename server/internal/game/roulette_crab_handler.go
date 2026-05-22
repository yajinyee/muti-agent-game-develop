// roulette_crab_handler.go — 黃金輪盤螃蟹系統 handler（DAY-167）
// 業界依據：King of Treasures Plus 2026「Roulette Crab — triggers Golden Roulette bonus game,
// player hits SHOOT to stop wheel, wins the amount listed where it stops.」
// 設計：T125 黃金輪盤螃蟹，擊破後觸發個人黃金輪盤（8格：10x-200x）
// 與千龍王輪盤（雙環，最高 1000x）不同：輪盤螃蟹是單環輪盤，更簡單直接
// 結果預先決定（公平性保證），玩家「停止」只是視覺互動（業界標準做法）
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/game/activityfeed"
	"digital-twin/server/internal/game/roulettecrab"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// isRouletteCrab 判斷是否為黃金輪盤螃蟹（T125）
func isRouletteCrab(defID string) bool {
	return defID == "T125"
}

// tryRouletteCrabWheel 擊破 T125 後觸發黃金輪盤（DAY-167）
// 由 handleKill 呼叫（goroutine）
func (g *Game) tryRouletteCrabWheel(p *player.Player, targetMult float64, baseReward int) {
	if g.RouletteCrab == nil {
		return
	}

	if !g.RouletteCrab.CanTrigger(p.ID) {
		// 冷卻中，不觸發
		log.Printf("[RouletteCrab] player=%s cooldown, skip", p.ID)
		return
	}

	// 開始 session（結果預先決定）
	session := g.RouletteCrab.StartSession(p.ID, targetMult, baseReward)

	log.Printf("[RouletteCrab] player=%s started wheel, target_mult=%.0f, base_reward=%d, wheel_result=%.0f",
		p.ID, targetMult, baseReward, session.WheelResult)

	// 廣播輪盤開始（全服看到）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRouletteCrabStart,
		Payload: ws.RouletteCrabStartPayload{
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			TargetMult: targetMult,
			BaseReward: baseReward,
			SpinSecs:   roulettecrab.SpinDuration,
			WheelSlots: roulettecrab.WheelSlots,
		},
	})

	// 全服公告
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "roulette_crab",
			"message":    fmt.Sprintf("🦀 %s 擊破了黃金輪盤螃蟹！黃金輪盤開始旋轉！", p.DisplayName),
			"color":      "#FFD700",
			"duration":   3.5,
			"priority":   2,
		},
	})
}

// handleRouletteCrabWheelStop 玩家停止輪盤（Client → Server）
// 由 HandleMessage 呼叫
func (g *Game) handleRouletteCrabWheelStop(p *player.Player) {
	if g.RouletteCrab == nil {
		return
	}

	session := g.RouletteCrab.StopSession(p.ID)
	if session == nil {
		// 沒有活躍 session 或已停止
		return
	}

	g.processRouletteCrabResult(p, session, false)
}

// processRouletteCrabResult 處理輪盤結果（停止或超時）
func (g *Game) processRouletteCrabResult(p *player.Player, session *roulettecrab.Session, isAutoStop bool) {
	bonusReward := session.BonusReward()

	// 發放獎勵
	if bonusReward > 0 {
		p.Coins += bonusReward
		if p.Coins > p.MaxCoins {
			p.MaxCoins = p.Coins
		}
	}

	log.Printf("[RouletteCrab] player=%s wheel_result=%.0f bonus_reward=%d auto_stop=%v",
		p.ID, session.WheelResult, bonusReward, isAutoStop)

	// 廣播結果（全服看到）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRouletteCrabResult,
		Payload: ws.RouletteCrabResultPayload{
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			WheelResult: session.WheelResult,
			SlotIndex:   session.SlotIndex,
			BaseReward:  session.BaseReward,
			BonusReward: bonusReward,
			NewBalance:  p.Coins,
			IsAutoStop:  isAutoStop,
		},
	})

	// 高倍率全服公告（≥100x）
	if session.WheelResult >= 100.0 {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgAnnouncement,
			Payload: map[string]interface{}{
				"event_type": "roulette_crab_big",
				"message":    fmt.Sprintf("🦀✨ %s 的黃金輪盤螃蟹轉出 ×%.0f！獲得 %d 金幣！", p.DisplayName, session.WheelResult, bonusReward),
				"color":      "#FFD700",
				"duration":   5.0,
				"priority":   4,
			},
		})
	}

	// 動態牆記錄（≥150x）
	if session.WheelResult >= 150.0 {
		go g.notifyFeedRouletteCrab(p, session.WheelResult, bonusReward)
	}
}

// tickRouletteCrabWheel 每秒檢查輪盤超時（由 gameLoop 呼叫）
func (g *Game) tickRouletteCrabWheel() {
	if g.RouletteCrab == nil {
		return
	}

	expired := g.RouletteCrab.TickAutoStop()
	for _, session := range expired {
		// 找到對應玩家
		g.mu.RLock()
		var p *player.Player
		for _, pl := range g.Players {
			if pl.ID == session.PlayerID {
				p = pl
				break
			}
		}
		g.mu.RUnlock()

		if p != nil {
			g.processRouletteCrabResult(p, session, true)
		}
	}
}

// sendRouletteCrabStatus 登入時發送輪盤螃蟹冷卻狀態（DAY-167）
func (g *Game) sendRouletteCrabStatus(p *player.Player) {
	if g.RouletteCrab == nil {
		return
	}
	cooldownLeft := g.RouletteCrab.GetCooldownLeft(p.ID)
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgRouletteCrabStatus,
		Payload: ws.RouletteCrabStatusPayload{
			PlayerID:     p.ID,
			CooldownLeft: cooldownLeft,
		},
	})
}

// notifyFeedRouletteCrab 動態牆記錄輪盤螃蟹大獎（DAY-167）
func (g *Game) notifyFeedRouletteCrab(p *player.Player, wheelResult float64, bonusReward int) {
	// 使用現有的 MegaWin 事件類型（輪盤螃蟹大獎 ≥150x 相當於 mega win）
	if g.ActivityFeed == nil {
		return
	}
	evt := g.ActivityFeed.Push(activityfeed.NewMegaWinEvent(
		p.ID, p.DisplayName, wheelResult, bonusReward,
	))
	if evt != nil {
		g.broadcastFeedEvent(evt)
	}
}
