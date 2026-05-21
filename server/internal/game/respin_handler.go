// respin_handler.go — Rapid Respin 系統 handler（DAY-121）
// 業界依據：Reflex Gaming Big Game Fishing Rapid Riches（2026-05-14）
// Rapid Respin 讓玩家在擊破目標後有機率觸發場上目標全部重新整理，
// 連鎖觸發最多 5 次，倍率遞增（1x→1.5x→2x→3x→5x）
package game

import (
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"digital-twin/server/internal/game/state"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// notifyRespinKill 在擊破目標後嘗試觸發 Rapid Respin（由 handleKill 呼叫）
// 回傳連鎖倍率加成（用於最終獎勵計算）
func (g *Game) notifyRespinKill(p *player.Player, finalReward int) float64 {
	if g.RespinMgr == nil {
		return 1.0
	}

	// Bonus Game 中不觸發
	g.mu.RLock()
	currentState := g.State
	g.mu.RUnlock()
	if currentState != state.StateNormalPlay {
		return 1.0
	}

	randFloat := rand.Float64()
	triggered, isChain, chainCount := g.RespinMgr.ShouldTrigger(p.ID, p.BetLevel, randFloat)
	if !triggered {
		return 1.0
	}

	// 取得當前連鎖倍率
	sess := g.RespinMgr.GetSession(p.ID)
	chainMult := 1.0
	if sess != nil {
		chainMult = sess.GetCurrentMult()
	}

	log.Printf("[Respin] player=%s triggered Rapid Respin chain=%d mult=%.1fx isChain=%v",
		p.ID, chainCount, chainMult, isChain)

	// 廣播 Rapid Respin 觸發（全服可見）
	icon := "⚡🔄"
	if isChain {
		icon = "🔥🔄"
	}
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRapidRespin,
		Payload: ws.RapidRespinPayload{
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			ChainCount: chainCount,
			ChainMult:  chainMult,
			IsChain:    isChain,
			MaxChain:   5,
			Icon:       icon,
		},
	})

	// 執行 Respin：清除場上所有非 BOSS 目標，重新生成
	go g.executeRespin(p, chainCount, chainMult)

	return chainMult
}

// executeRespin 執行 Rapid Respin 邏輯
// 清除場上所有非 BOSS 目標，重新生成一批新目標
func (g *Game) executeRespin(p *player.Player, chainCount int, chainMult float64) {
	// 短暫延遲讓 Client 先播放動畫
	time.Sleep(300 * time.Millisecond)

	g.mu.Lock()
	// 收集要清除的目標（非 BOSS）
	toRemove := make([]string, 0)
	for id, t := range g.Targets {
		if t.DefID != "B001" { // 保留 BOSS
			toRemove = append(toRemove, id)
		}
	}
	// 清除目標
	for _, id := range toRemove {
		delete(g.Targets, id)
	}
	g.mu.Unlock()

	// 廣播目標清除（讓 Client 移除這些目標）
	for _, id := range toRemove {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetKill,
			Payload: ws.TargetKillPayload{
				InstanceID: id,
				DefID:      "respin_clear", // 特殊標記，Client 不顯示獎勵
				Multiplier: 0,
				Reward:     0,
				LaborGain:  0,
				KillerID:   "respin",
			},
		})
	}

	// 短暫等待後生成新目標（讓 Client 有時間清除舊目標）
	time.Sleep(200 * time.Millisecond)

	// 生成新一批目標（數量依連鎖次數遞增）
	spawnCount := 6 + chainCount*2 // 第一次 6 個，每次連鎖多 2 個
	if spawnCount > 14 {
		spawnCount = 14
	}

	g.mu.Lock()
	// 取平均 bet level
	betLevel := 5
	if len(g.Players) > 0 {
		total := 0
		for _, pl := range g.Players {
			total += pl.BetLevel
		}
		betLevel = total / len(g.Players)
		if betLevel < 1 {
			betLevel = 1
		}
	}
	for i := 0; i < spawnCount; i++ {
		def := g.SpawnSys.PickTargetDef(betLevel, 0.1) // 稍微提高特殊目標機率
		if def == nil {
			continue
		}
		instanceID := uuid.New().String()
		x := 1280.0 + rand.Float64()*100
		y := 100.0 + rand.Float64()*500
		t := target.NewTarget(instanceID, def, x, y)
		// Respin 期間目標倍率加成
		if chainMult > 1.0 {
			t.Multiplier = t.Multiplier * chainMult
		}
		g.Targets[instanceID] = t
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetSpawn,
			Payload: ws.TargetSpawnPayload{
				InstanceID:   instanceID,
				DefID:        def.ID,
				Name:         def.Name,
				Type:         string(def.Type),
				X:            x,
				Y:            y,
				HP:           def.HP,
				MaxHP:        def.HP,
				Speed:        def.Speed,
				Lifetime:     def.Lifetime,
				Behavior:     def.SpecialBehavior,
				Multiplier:   t.Multiplier,
				Quality:      string(t.Quality),
				QualityColor: t.QualityColor,
			},
		})
	}
	g.mu.Unlock()

	log.Printf("[Respin] executed: cleared=%d spawned=%d chain=%d mult=%.1fx",
		len(toRemove), spawnCount, chainCount, chainMult)
}

// checkRespinSessionExpiry 定期檢查 Respin session 是否過期（由 gameLoop 呼叫）
func (g *Game) checkRespinSessionExpiry() {
	if g.RespinMgr == nil {
		return
	}
	g.mu.RLock()
	players := make([]*player.Player, 0, len(g.Players))
	for _, p := range g.Players {
		players = append(players, p)
	}
	g.mu.RUnlock()

	for _, p := range players {
		sess := g.RespinMgr.GetSession(p.ID)
		if sess == nil {
			continue
		}
		// session 已過期（連鎖視窗結束）
		if !sess.CanChain() && sess.Active {
			g.RespinMgr.EndSession(p.ID)
			// 通知玩家連鎖結束
			g.Hub.Send(p.ID, &ws.Message{
				Type: ws.MsgRapidRespinEnd,
				Payload: ws.RapidRespinEndPayload{
					PlayerID:   p.ID,
					PlayerName: p.DisplayName,
					TotalChain: sess.ChainCount + 1,
				},
			})
			log.Printf("[Respin] session ended for player=%s totalChain=%d", p.ID, sess.ChainCount+1)
		}
	}
}
