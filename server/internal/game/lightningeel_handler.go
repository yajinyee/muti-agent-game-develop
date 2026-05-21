// lightningeel_handler.go — 閃電鰻連鎖攻擊系統 handler（DAY-132）
// 業界依據：JILI Royal Fishing 2026 — 「The 60x lightning eel creates chain reactions
// that jump between nearby fish. Once activated, electric shocks continue spreading
// until targeting disengages, creating cascading capture sequences.」
// 閃電鰻（T103）擊破後釋放閃電連鎖，在附近目標之間跳躍傳導，
// 每次跳躍 45% 機率直接擊破目標，最多跳 5 次，製造「一箭多雕」的爽感。
package game

import (
	"fmt"
	"log"
	"math"

	"digital-twin/server/internal/game/activityfeed"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/lightningeel"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// LightningEelDefID 閃電鰻的目標定義 ID（T103 流星/閃電鰻）
// 規格書中 T103 是「流星」，這裡重新定義為閃電鰻（視覺上是流星形狀的電鰻）
const LightningEelDefID = "T103"

// isLightningEel 判斷目標是否為閃電鰻
func isLightningEel(defID string) bool {
	return defID == LightningEelDefID
}

// tryLightningEelChain 嘗試觸發閃電鰻連鎖攻擊（由 handleKill 呼叫）
// 只有擊破閃電鰻（T103）且冷卻結束時才觸發
func (g *Game) tryLightningEelChain(p *player.Player, killedInstanceID string, killedX, killedY float64) {
	if g.LightningEel == nil {
		return
	}
	if !g.LightningEel.CanTrigger(p.ID) {
		return
	}

	// 收集附近目標（排除已擊破的閃電鰻本身）
	g.mu.RLock()
	nearby := make([]lightningeel.NearbyTarget, 0, 10)
	cfg := g.LightningEel.GetConfig()
	for _, t := range g.Targets {
		if t.InstanceID == killedInstanceID {
			continue
		}
		// 計算距離
		dx := t.X - killedX
		dy := t.Y - killedY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= cfg.JumpRangeUnits {
			nearby = append(nearby, lightningeel.NearbyTarget{
				InstanceID: t.InstanceID,
				DefID:      t.DefID,
				Name:       t.Def.Name,
				Multiplier: float64(t.Multiplier),
				X:          t.X,
				Y:          t.Y,
			})
		}
	}
	g.mu.RUnlock()

	if len(nearby) == 0 {
		// 沒有附近目標，不觸發（不消耗冷卻，讓玩家可以立即再試）
		return
	}

	// 執行連鎖攻擊
	betDef := p.GetBetDef()
	session := g.LightningEel.ExecuteChain(p.ID, killedInstanceID, nearby, int64(betDef.BetCost))

	// 處理連鎖結果：擊破目標、發放獎勵
	jumps := make([]ws.LightningEelJumpEntry, 0, len(session.Jumps))
	for _, jump := range session.Jumps {
		if jump.Killed {
			// 從 Targets 中移除被連鎖擊破的目標
			g.mu.Lock()
			t, exists := g.Targets[jump.TargetInstanceID]
			if exists {
				delete(g.Targets, jump.TargetInstanceID)
			}
			g.mu.Unlock()

			if exists {
				// 發放獎勵給觸發玩家
				p.AddCoins(int(jump.Reward))
				// 廣播目標擊破（讓 Client 播放死亡動畫）
				g.Hub.Broadcast(&ws.Message{
					Type: ws.MsgTargetKill,
					Payload: ws.TargetKillPayload{
						InstanceID: jump.TargetInstanceID,
						DefID:      t.DefID,
						Multiplier: jump.Multiplier,
						Reward:     int(jump.Reward),
						LaborGain:  0, // 連鎖擊破不累積勞動值
						KillerID:   p.ID,
						Quality:    string(t.Quality),
					},
				})
				log.Printf("[LightningEel] chain kill: player=%s target=%s mult=%.1f reward=%d",
					p.ID, jump.TargetInstanceID, jump.Multiplier, jump.Reward)
			}
		}

		jumps = append(jumps, ws.LightningEelJumpEntry{
			TargetInstanceID: jump.TargetInstanceID,
			TargetDefID:      jump.TargetDefID,
			TargetName:       jump.TargetName,
			Killed:           jump.Killed,
			Multiplier:       jump.Multiplier,
			Reward:           jump.Reward,
			JumpIndex:        jump.JumpIndex,
		})
	}

	// 廣播連鎖結果（全服可見，讓所有玩家看到閃電效果）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLightningEelChain,
		Payload: ws.LightningEelChainPayload{
			PlayerID:        p.ID,
			PlayerName:      p.DisplayName,
			TriggerTargetID: killedInstanceID,
			Jumps:           jumps,
			TotalKills:      session.TotalKills,
			TotalReward:     session.TotalReward,
			NewBalance:      int64(p.Coins),
		},
	})

	if session.TotalKills > 0 {
		log.Printf("[LightningEel] player=%s triggered chain: %d jumps, %d kills, reward=%d",
			p.ID, len(session.Jumps), session.TotalKills, session.TotalReward)
	}

	// 全服公告：連鎖擊破 ≥3 個時廣播
	if session.TotalKills >= 3 {
		g.announceLightningChain(p.DisplayName, session.TotalKills, int(session.TotalReward))
	}

	// 動態牆：連鎖擊破 ≥4 個時廣播
	if session.TotalKills >= 4 {
		go g.notifyFeedLightningChain(p, session.TotalKills, int(session.TotalReward))
	}
}

// sendLightningEelStatus 發送閃電鰻冷卻狀態給玩家（登入時呼叫）
func (g *Game) sendLightningEelStatus(p *player.Player) {
	if g.LightningEel == nil {
		return
	}
	snap := g.LightningEel.GetSnapshot(p.ID)
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLightningEelStatus,
		Payload: ws.LightningEelStatusPayload{
			PlayerID:     p.ID,
			CooldownLeft: snap.CooldownLeft,
			MaxJumps:     snap.Config.MaxJumps,
			JumpRange:    snap.Config.JumpRangeUnits,
		},
	})
}

// announceLightningChain 全服公告：閃電鰻連鎖擊破
func (g *Game) announceLightningChain(playerName string, kills int, reward int) {
	icon := "⚡"
	if kills >= 5 {
		icon = "🌩️"
	}
	extra := map[string]string{
		"kills":  fmt.Sprintf("%d", kills),
		"reward": fmt.Sprintf("%d", reward),
	}
	ann := g.Announce.Create(announce.EventLightningChain, playerName, reward, extra)
	_ = icon
	g.broadcastAnnouncement(ann)
}

// notifyFeedLightningChain 動態牆：閃電鰻連鎖記錄
func (g *Game) notifyFeedLightningChain(p *player.Player, kills int, reward int) {
	if g.ActivityFeed == nil {
		return
	}
	icon := "⚡"
	if kills >= 5 {
		icon = "🌩️"
	}
	rarity := activityfeed.RarityRare
	if kills >= 5 {
		rarity = activityfeed.RarityEpic
	}
	evt := g.ActivityFeed.Push(&activityfeed.FeedEvent{
		EventType:   activityfeed.EventLightningChain,
		PlayerID:    p.ID,
		DisplayName: p.DisplayName,
		Icon:        icon,
		Title:       "閃電連鎖",
		Detail:      fmt.Sprintf("連鎖擊破 %d 個目標，獲得 %d 金幣", kills, reward),
		Rarity:      rarity,
	})
	go g.broadcastFeedEvent(evt)
}
