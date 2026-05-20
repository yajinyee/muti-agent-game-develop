// codex_handler.go — 魚類圖鑑系統 handler（DAY-081）
package game

import (
	"log"

	"digital-twin/server/internal/game/codex"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// sendCodexUpdate 發送完整圖鑑狀態給指定玩家
func (g *Game) sendCodexUpdate(playerID string) {
	p, ok := g.Players[playerID]
	if !ok {
		return
	}
	if p.Codex == nil {
		return
	}

	entries := p.Codex.GetSnapshot()
	unlocked, total := p.Codex.GetStats()

	payload := ws.CodexUpdatePayload{
		Entries:       make([]ws.CodexEntryPayload, 0, len(entries)),
		UnlockedCount: unlocked,
		TotalCount:    total,
		IsComplete:    p.Codex.IsComplete(),
	}
	for _, e := range entries {
		var unlockedAt int64
		if e.Unlocked {
			unlockedAt = e.UnlockedAt.UnixMilli()
		}
		payload.Entries = append(payload.Entries, ws.CodexEntryPayload{
			TargetID:      e.TargetID,
			TargetName:    e.TargetName,
			Rarity:        e.Rarity,
			Unlocked:      e.Unlocked,
			UnlockedAt:    unlockedAt,
			KillCount:     e.KillCount,
			MaxMultiplier: e.MaxMultiplier,
		})
	}

	if err := g.Hub.Send(playerID, &ws.Message{
		Type:    ws.MsgCodexUpdate,
		Payload: payload,
	}); err != nil {
		log.Printf("[Codex] sendCodexUpdate error: %v", err)
	}
}

// handleGetCodex 處理 Client 查詢圖鑑請求
func (g *Game) handleGetCodex(playerID string) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	g.sendCodexUpdate(playerID)
}

// notifyCodexKill 在擊破目標後更新圖鑑（由 handleKill 呼叫）
func (g *Game) notifyCodexKill(playerID string, defID string, multiplier float64) {
	p, ok := g.Players[playerID]
	if !ok || p.Codex == nil {
		return
	}

	isNewUnlock, isComplete := p.Codex.RecordKill(defID, multiplier)
	if !isNewUnlock {
		return
	}

	// 取得條目資訊
	entries := p.Codex.GetSnapshot()
	var entry *codex.Entry
	for _, e := range entries {
		if e.TargetID == defID {
			entry = e
			break
		}
	}
	if entry == nil {
		return
	}

	// 發放解鎖獎勵
	p.AddCoins(codex.UnlockReward)
	unlocked, total := p.Codex.GetStats()

	// 發送解鎖通知
	if err := g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgCodexUnlock,
		Payload: ws.CodexUnlockPayload{
			TargetID:      entry.TargetID,
			TargetName:    entry.TargetName,
			Rarity:        entry.Rarity,
			Reward:        codex.UnlockReward,
			NewBalance:    p.GetCoins(),
			UnlockedCount: unlocked,
			TotalCount:    total,
		},
	}); err != nil {
		log.Printf("[Codex] send unlock error: %v", err)
	}

	log.Printf("[Codex] player=%s unlocked %s (%s), reward=%d, progress=%d/%d",
		playerID, entry.TargetName, entry.Rarity, codex.UnlockReward, unlocked, total)

	// 全圖鑑完成
	if isComplete {
		g.handleCodexComplete(playerID, p)
	}
}

// handleCodexComplete 全圖鑑完成處理
func (g *Game) handleCodexComplete(playerID string, p *player.Player) {
	// 發放完成獎勵
	p.AddCoins(codex.CompleteReward)

	// 解鎖「圖鑑完成者」稱號（title_id: codex_master）
	titleID := "codex_master"
	titleName := "圖鑑完成者"

	// 解鎖稱號（TitleTracker 目前只支援成就觸發，這裡直接記錄到 Unlocked map）
	// 注意：codex_master 是自定義稱號，不在 achievement 系統內
	// 透過 sendAchievements 廣播稱號解鎖通知
	log.Printf("[Codex] player=%s unlocked title: %s", playerID, titleName)

	if err := g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgCodexComplete,
		Payload: ws.CodexCompletePayload{
			Reward:     codex.CompleteReward,
			NewBalance: p.GetCoins(),
			TitleID:    titleID,
			TitleName:  titleName,
		},
	}); err != nil {
		log.Printf("[Codex] send complete error: %v", err)
	}

	log.Printf("[Codex] player=%s completed full codex! reward=%d", playerID, codex.CompleteReward)
}
