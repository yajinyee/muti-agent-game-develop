// golden_jellyfish_handler.go — 黃金水母全場電擊系統 handler（DAY-149）
// 業界依據：Ocean King 3 2026「Electric Jellyfish chain shocks across multiple targets.
// Devastating against clustered schools.」
// 設計：T113 黃金水母擊破後觸發「全場電擊」，對畫面上所有目標發動電擊
// 比閃電鰻（T103，200px 範圍跳躍 5 次）更強：全場範圍，最多 8 個目標，40% 擊破機率
// 電擊間隔 150ms，製造「電流掃場」的連續感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	GoldenJellyfishDefID    = "T113"
	GoldenJellyfishKillChance = 0.40  // 每個目標 40% 擊破機率
	GoldenJellyfishMaxTargets = 8     // 最多電擊 8 個目標
	GoldenJellyfishShockInterval = 150 * time.Millisecond // 電擊間隔
	GoldenJellyfishRewardMult = 0.65  // 連帶擊破獎勵係數（比直接擊破低，平衡 RTP）
	GoldenJellyfishCooldownSecs = 10  // 冷卻秒數（比閃電鰻的 8 秒略長）
)

// isGoldenJellyfish 判斷是否為黃金水母目標
func isGoldenJellyfish(defID string) bool {
	return defID == GoldenJellyfishDefID
}

// tryGoldenJellyfishShock 擊破黃金水母後觸發全場電擊（由 handleKill 呼叫）
func (g *Game) tryGoldenJellyfishShock(p *player.Player, killedInstanceID string, killedX, killedY float64) {
	// 收集所有存活目標（排除已擊破的黃金水母本身）
	g.mu.RLock()
	type candidateTarget struct {
		instanceID string
		defID      string
		name       string
		multiplier float64
		x, y       float64
	}
	candidates := make([]candidateTarget, 0, 16)
	for _, t := range g.Targets {
		if t.InstanceID == killedInstanceID {
			continue
		}
		candidates = append(candidates, candidateTarget{
			instanceID: t.InstanceID,
			defID:      t.DefID,
			name:       t.Def.Name,
			multiplier: float64(t.Multiplier),
			x:          t.X,
			y:          t.Y,
		})
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		return
	}

	// 隨機打亂順序（讓電擊看起來更自然）
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	// 限制最多 8 個目標
	if len(candidates) > GoldenJellyfishMaxTargets {
		candidates = candidates[:GoldenJellyfishMaxTargets]
	}

	betDef := p.GetBetDef()
	triggerID := killedInstanceID

	log.Printf("[GoldenJellyfish] player=%s triggered global shock: %d candidates", p.ID, len(candidates))

	// 廣播電擊開始（全服可見）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGoldenJellyfishShock,
		Payload: ws.GoldenJellyfishShockPayload{
			TriggerID:  triggerID,
			TriggerX:   killedX,
			TriggerY:   killedY,
			Phase:      "shock_start",
			KillerID:   p.ID,
			KillerName: p.DisplayName,
			Message:    fmt.Sprintf("⚡ %s 的黃金水母觸發全場電擊！", p.DisplayName),
		},
	})

	// 逐一電擊目標（分批廣播，製造連續感）
	totalKills := 0
	totalReward := 0
	shockEntries := make([]ws.GoldenJellyfishShockEntry, 0, len(candidates))

	for i, c := range candidates {
		time.Sleep(GoldenJellyfishShockInterval)

		// 判斷是否擊破（40% 機率）
		killed := rand.Float64() < GoldenJellyfishKillChance
		reward := 0

		if killed {
			g.mu.Lock()
			t, exists := g.Targets[c.instanceID]
			if exists && t.IsAlive {
				// 計算獎勵（比直接擊破低）
				reward = int(float64(betDef.BetCost) * c.multiplier * GoldenJellyfishRewardMult)
				t.IsAlive = false
				t.HP = 0
				delete(g.Targets, c.instanceID)
				g.mu.Unlock()

				// 發放獎勵
				p.AddCoins(reward)
				totalKills++
				totalReward += reward

				// 廣播目標擊破
				g.Hub.Broadcast(&ws.Message{
					Type: ws.MsgTargetKill,
					Payload: ws.TargetKillPayload{
						InstanceID: c.instanceID,
						DefID:      c.defID,
						Multiplier: c.multiplier,
						Reward:     reward,
						LaborGain:  0,
						KillerID:   p.ID,
						Quality:    "normal",
					},
				})
			} else {
				g.mu.Unlock()
				killed = false // 目標已不存在
			}
		}

		entry := ws.GoldenJellyfishShockEntry{
			TargetInstanceID: c.instanceID,
			TargetDefID:      c.defID,
			TargetName:       c.name,
			Killed:           killed,
			Multiplier:       c.multiplier,
			Reward:           reward,
			ShockIndex:       i,
		}
		shockEntries = append(shockEntries, entry)

		// 廣播單次電擊（讓 Client 播放電擊動畫）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgGoldenJellyfishShock,
			Payload: ws.GoldenJellyfishShockPayload{
				TriggerID:  triggerID,
				TriggerX:   killedX,
				TriggerY:   killedY,
				Phase:      "shock",
				KillerID:   p.ID,
				KillerName: p.DisplayName,
				Targets:    []ws.GoldenJellyfishShockEntry{entry},
				Message:    "",
			},
		})
	}

	// 廣播電擊結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGoldenJellyfishShock,
		Payload: ws.GoldenJellyfishShockPayload{
			TriggerID:   triggerID,
			TriggerX:    killedX,
			TriggerY:    killedY,
			Phase:       "result",
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			Targets:     shockEntries,
			TotalKills:  totalKills,
			TotalReward: totalReward,
			NewBalance:  p.GetCoins(),
			Message:     fmt.Sprintf("⚡ %s 黃金水母電擊結果：擊破 %d 個目標，獲得 %d 金幣！", p.DisplayName, totalKills, totalReward),
		},
	})

	log.Printf("[GoldenJellyfish] player=%s shock result: kills=%d, reward=%d", p.ID, totalKills, totalReward)

	// 全服公告：擊破 ≥3 個目標時廣播
	if totalKills >= 3 {
		g.announceGoldenJellyfishShock(p.DisplayName, totalKills, totalReward)
	}

	// 動態牆：擊破 ≥4 個目標時廣播
	if totalKills >= 4 {
		go g.notifyFeedMegaWin(p, float64(totalKills)*10.0, totalReward)
	}

	// 更新玩家狀態
	g.sendPlayerUpdate(p)
}

// announceGoldenJellyfishShock 全服公告：黃金水母電擊
func (g *Game) announceGoldenJellyfishShock(playerName string, kills int, reward int) {
	extra := map[string]string{
		"kills":  fmt.Sprintf("%d", kills),
		"reward": fmt.Sprintf("%d", reward),
	}
	ann := g.Announce.Create(announce.EventLightningChain, playerName, reward, extra)
	g.broadcastAnnouncement(ann)
}
