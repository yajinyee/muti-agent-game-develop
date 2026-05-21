// megaoctopus_handler.go — 巨型章魚轉盤 handler（DAY-144）
// 業界依據：JILI Mega Fishing「Mega Octopus Wheel – Defeat that giant octopus and enter
// the bonus wheel round where you have a chance to win massive guaranteed prizes up to 950x.」
package game

import (
	"fmt"
	"log"
	"time"

	"digital-twin/server/internal/game/megaoctopus"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// isMegaOctopus 判斷是否為巨型章魚（T108）
func isMegaOctopus(defID string) bool {
	return defID == "T108"
}

// tryMegaOctopusWheel 擊破 T108 後觸發個人轉盤（DAY-144）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryMegaOctopusWheel(p *player.Player, triggerID string) {
	// 如果玩家已有活躍 session，不重複觸發
	if g.MegaOctopus.HasActiveSession(p.ID) {
		return
	}

	// 開始轉盤 session（預先決定結果）
	session := g.MegaOctopus.StartSession(p.ID)
	if session == nil {
		return
	}

	// 廣播轉盤開始（只發給觸發玩家）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgMegaOctopusWheelStart,
		Payload: ws.MegaOctopusWheelStartPayload{
			TriggerID:    triggerID,
			SpinDuration: megaoctopus.SpinDuration,
			Slots:        buildWheelSlotsPayload(),
		},
	})

	log.Printf("[MegaOctopus] player=%s wheel started, result_index=%d (mult=%dx)",
		p.ID, session.ResultIndex, megaoctopus.WheelSlots[session.ResultIndex].Multiplier)

	// 等待玩家停止或超時自動停止
	// 超時由 tickMegaOctopus 處理
}

// handleMegaOctopusStop 玩家點擊停止轉盤（DAY-144）
func (g *Game) handleMegaOctopusStop(p *player.Player) {
	resultIndex, multiplier, ok := g.MegaOctopus.StopSession(p.ID)
	if !ok {
		return
	}

	g.processMegaOctopusResult(p, resultIndex, multiplier)
}

// processMegaOctopusResult 處理轉盤結果（DAY-144）
func (g *Game) processMegaOctopusResult(p *player.Player, resultIndex int, multiplier int) {
	// 計算獎勵
	reward := multiplier * p.BetLevel
	if reward < 1 {
		reward = 1
	}

	// 發放獎勵
	p.AddReward(reward)

	slot := megaoctopus.GetResultSlot(resultIndex)

	// 發送結果給玩家
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgMegaOctopusWheelResult,
		Payload: ws.MegaOctopusWheelResultPayload{
			ResultIndex: resultIndex,
			Multiplier:  multiplier,
			Reward:      reward,
			NewBalance:  p.Coins,
			SlotLabel:   slot.Label,
			SlotColor:   slot.Color,
		},
	})

	// 更新玩家狀態
	g.sendPlayerUpdate(p)

	// 全服公告：≥300x
	if multiplier >= megaoctopus.AnnounceThreshold {
		g.announceMegaOctopusWin(p.DisplayName, multiplier, reward)
	}

	log.Printf("[MegaOctopus] player=%s result=%dx reward=%d",
		p.ID, multiplier, reward)
}

// tickMegaOctopus 超時自動停止（由 game loop 呼叫）
func (g *Game) tickMegaOctopus() {
	expired := g.MegaOctopus.AutoStop()
	for _, playerID := range expired {
		g.mu.RLock()
		p, ok := g.Players[playerID]
		g.mu.RUnlock()
		if !ok {
			continue
		}

		// 自動停止：重新執行 StopSession
		resultIndex, multiplier, ok := g.MegaOctopus.StopSession(playerID)
		if !ok {
			continue
		}
		go g.processMegaOctopusResult(p, resultIndex, multiplier)
	}
}

// announceMegaOctopusWin 全服公告巨型章魚轉盤大獎（DAY-144）
func (g *Game) announceMegaOctopusWin(playerName string, multiplier int, reward int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "mega_octopus_win",
			"message":    fmt.Sprintf("🐙 %s 的巨型章魚轉盤獲得 %dx！獲得 %d 金幣！", playerName, multiplier, reward),
			"color":      "#9400D3",
			"duration":   5.0,
			"priority":   3,
		},
	})
}

// buildWheelSlotsPayload 建立轉盤格子 payload
func buildWheelSlotsPayload() []ws.OctopusWheelSlotPayload {
	slots := make([]ws.OctopusWheelSlotPayload, len(megaoctopus.WheelSlots))
	for i, slot := range megaoctopus.WheelSlots {
		slots[i] = ws.OctopusWheelSlotPayload{
			Index:      i,
			Multiplier: slot.Multiplier,
			Color:      slot.Color,
			Label:      slot.Label,
		}
	}
	return slots
}

// sendMegaOctopusStatus 登入時發送狀態（如果有活躍 session）
func (g *Game) sendMegaOctopusStatus(p *player.Player) {
	if g.MegaOctopus.HasActiveSession(p.ID) {
		// 有活躍 session，重新發送轉盤開始（斷線重連）
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgMegaOctopusWheelStart,
			Payload: ws.MegaOctopusWheelStartPayload{
				TriggerID:    "",
				SpinDuration: megaoctopus.SpinDuration,
				Slots:        buildWheelSlotsPayload(),
			},
		})
	}
}

// 每 2 秒執行一次超時檢查
func (g *Game) startMegaOctopusTicker() {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				g.tickMegaOctopus()
			case <-g.stopCh:
				return
			}
		}
	}()
}
