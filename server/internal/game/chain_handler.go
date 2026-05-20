// chain_handler.go — 連鎖爆炸系統 handler（DAY-088）
package game

import (
	"log"

	"digital-twin/server/internal/game/chain"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// notifyChainKill 嘗試觸發連鎖爆炸（由 handleKill 呼叫）
// triggerID: 觸發連鎖的目標 ID
// triggerX, triggerY: 觸發目標的位置
// triggerMult: 觸發目標的倍率
// triggerDefID: 觸發目標的定義 ID
// 回傳連鎖總獎勵（已加入玩家金幣）
func (g *Game) notifyChainKill(p *player.Player, triggerID string, triggerX, triggerY, triggerMult float64, triggerDefID string) int {
	if g.Chain == nil {
		return 0
	}

	// 收集場上所有目標資訊
	g.mu.RLock()
	allTargets := make([]chain.TargetInfo, 0, len(g.Targets))
	for id, t := range g.Targets {
		allTargets = append(allTargets, chain.TargetInfo{
			ID:         id,
			X:          t.X,
			Y:          t.Y,
			Multiplier: t.Multiplier,
			DefID:      t.DefID,
		})
	}
	g.mu.RUnlock()

	trigger := chain.TargetInfo{
		ID:         triggerID,
		X:          triggerX,
		Y:          triggerY,
		Multiplier: triggerMult,
		DefID:      triggerDefID,
	}

	result := g.Chain.TryChain(trigger, allTargets, 0)
	if result.Level == chain.ChainNone || len(result.TargetIDs) == 0 {
		return 0
	}

	// 執行連鎖擊破
	totalReward := 0
	chainEntries := make([]ws.ChainKillEntry, 0, len(result.TargetIDs))

	g.mu.Lock()
	for _, chainTargetID := range result.TargetIDs {
		ct, exists := g.Targets[chainTargetID]
		if !exists {
			continue
		}
		// 計算連鎖獎勵（基礎獎勵 × 連鎖倍率加成）
		baseReward := int(float64(p.BetLevel*10) * ct.Multiplier)
		chainReward := int(float64(baseReward) * result.BonusMult)
		totalReward += chainReward

		chainEntries = append(chainEntries, ws.ChainKillEntry{
			InstanceID: chainTargetID,
			DefID:      ct.DefID,
			Multiplier: ct.Multiplier,
			Reward:     chainReward,
		})

		// 從場上移除連鎖目標
		delete(g.Targets, chainTargetID)
	}
	g.mu.Unlock()

	if len(chainEntries) == 0 {
		return 0
	}

	// 發放連鎖獎勵
	p.AddCoins(totalReward)

	// 廣播連鎖爆炸事件（所有玩家都能看到）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgChainExplosion,
		Payload: ws.ChainExplosionPayload{
			TriggerID:   triggerID,
			Level:       int(result.Level),
			LevelName:   result.LevelName,
			LevelColor:  result.LevelColor,
			Chains:      chainEntries,
			TotalReward: totalReward,
			BonusMult:   result.BonusMult,
			PlayerID:    p.ID,
		},
	})

	log.Printf("[Chain] player=%s triggered %s (%d targets, reward=%d)",
		p.ID, result.LevelName, len(chainEntries), totalReward)

	// 節日任務：記錄連鎖爆炸（DAY-109）
	go g.notifyFestivalChain(p)

	return totalReward
}
