// chain_bomb_handler.go — 連鎖爆炸魚系統 handler（DAY-187）
// 業界依據：Royal Fishing「chain reaction mechanic — players can trigger multiple explosions
// to capture additional fish within a blast radius」
// T145 連鎖爆炸魚機制：
//   1. 擊破 T145 後，在原位觸發爆炸（200px 半徑）
//   2. 爆炸範圍內的目標有 75% 機率被擊破（獎勵 0.65x 倍率）
//   3. 若爆炸命中的目標也是 T145，則繼續引爆（連鎖反應，最多 5 層）
//   4. 全服廣播每層爆炸，讓所有玩家看到「連鎖爆炸在全場蔓延」
//   5. 連鎖結束後廣播總結（層數/擊破數/總獎勵）
// 設計差異：
//   - 與炸彈武器（玩家主動放置，即時爆炸）不同，連鎖爆炸魚是「被動觸發的連鎖反應」
//   - 與漩渦魚（吸引同類）不同，連鎖爆炸魚是「爆炸傳播」（位置驅動，不是類型驅動）
//   - 最多 5 層連鎖，讓玩家有「一顆引爆全場」的爽感，但不會無限連鎖（平衡 RTP）
//   - 每層爆炸間隔 300ms，讓玩家看到「爆炸在場上蔓延」的視覺過程
package game

import (
	"log"
	"math"
	"math/rand"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// ChainBombRadius 爆炸半徑（px）
	ChainBombRadius = 200.0
	// ChainBombKillChance 爆炸擊破機率（75%）
	ChainBombKillChance = 0.75
	// ChainBombRewardMult 爆炸擊破獎勵倍率（比直接擊破低，平衡 RTP）
	ChainBombRewardMult = 0.65
	// ChainBombMaxChain 最大連鎖層數
	ChainBombMaxChain = 5
	// ChainBombIntervalMs 每層爆炸間隔（毫秒）
	ChainBombIntervalMs = 300
	// ChainBombAnnounceMinKills 全服公告最低擊破數
	ChainBombAnnounceMinKills = 4
)

// isChainBomb 判斷是否為連鎖爆炸魚
func isChainBomb(defID string) bool {
	return defID == "T145"
}

// chainBombEntry 連鎖爆炸單次爆炸記錄
type chainBombEntry struct {
	InstanceID string  `json:"instance_id"` // 爆炸中心目標 ID
	X          float64 `json:"x"`           // 爆炸位置 X
	Y          float64 `json:"y"`           // 爆炸位置 Y
	KillCount  int     `json:"kill_count"`  // 本次爆炸擊破數
	Reward     int     `json:"reward"`      // 本次爆炸獎勵
	ChainDepth int     `json:"chain_depth"` // 連鎖層數（1=初始）
}

// tryChainBombExplosion 連鎖爆炸魚擊破後觸發（由 handleKill 呼叫）
func (g *Game) tryChainBombExplosion(p *player.Player, triggerID string, triggerX, triggerY float64) {
	log.Printf("[ChainBomb] triggered by player=%s at (%.0f,%.0f)", p.ID, triggerX, triggerY)

	// 廣播連鎖開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgChainBomb,
		Payload: ws.ChainBombPayload{
			Phase:      "chain_start",
			TriggerID:  triggerID,
			TriggerX:   triggerX,
			TriggerY:   triggerY,
			KillerID:   p.ID,
			KillerName: p.DisplayName,
		},
	})

	// 執行連鎖爆炸（遞迴，最多 5 層）
	var entries []chainBombEntry
	totalKills := 0
	totalReward := 0

	// 追蹤已爆炸的目標（防止重複爆炸）
	exploded := map[string]bool{triggerID: true}

	// 第一層：以觸發目標為中心爆炸
	pendingExplosions := []struct {
		x, y  float64
		depth int
		srcID string
	}{{triggerX, triggerY, 1, triggerID}}

	for len(pendingExplosions) > 0 && len(entries) < ChainBombMaxChain {
		current := pendingExplosions[0]
		pendingExplosions = pendingExplosions[1:]

		// 間隔廣播（讓玩家看到爆炸蔓延）
		if current.depth > 1 {
			time.Sleep(ChainBombIntervalMs * time.Millisecond)
		}

		// 收集爆炸範圍內的目標
		g.mu.RLock()
		type nearbyTarget struct {
			instanceID string
			defID      string
			x, y       float64
			multiplier float64
		}
		var nearby []nearbyTarget
		for id, t := range g.Targets {
			if exploded[id] || !t.IsAlive {
				continue
			}
			dx := t.X - current.x
			dy := t.Y - current.y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= ChainBombRadius {
				nearby = append(nearby, nearbyTarget{
					instanceID: id,
					defID:      t.DefID,
					x:          t.X,
					y:          t.Y,
					multiplier: t.Multiplier,
				})
			}
		}
		g.mu.RUnlock()

		// 對範圍內目標執行爆炸
		entry := chainBombEntry{
			InstanceID: current.srcID,
			X:          current.x,
			Y:          current.y,
			ChainDepth: current.depth,
		}

		for _, nt := range nearby {
			if rand.Float64() >= ChainBombKillChance {
				continue // 未命中
			}

			// 擊破目標
			g.mu.Lock()
			t, exists := g.Targets[nt.instanceID]
			if !exists || !t.IsAlive {
				g.mu.Unlock()
				continue
			}
			t.IsAlive = false
			delete(g.Targets, nt.instanceID)
			g.mu.Unlock()

			exploded[nt.instanceID] = true

			// 計算獎勵
			reward := int(nt.multiplier * float64(p.BetLevel) * ChainBombRewardMult)
			p.AddCoins(reward)
			entry.KillCount++
			entry.Reward += reward
			totalKills++
			totalReward += reward

			// 廣播目標被爆炸擊破
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgTargetKill,
				Payload: ws.TargetKillPayload{
					InstanceID: nt.instanceID,
					KillerID:   p.ID,
					Reward:     reward,
					Multiplier: nt.multiplier,
				},
			})

			// 如果被爆炸的目標也是連鎖爆炸魚，加入下一層爆炸隊列
			if isChainBomb(nt.defID) && current.depth < ChainBombMaxChain {
				pendingExplosions = append(pendingExplosions, struct {
					x, y  float64
					depth int
					srcID string
				}{nt.x, nt.y, current.depth + 1, nt.instanceID})
			}
		}

		entries = append(entries, entry)

		// 廣播本層爆炸
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgChainBomb,
			Payload: ws.ChainBombPayload{
				Phase:      "chain_explode",
				TriggerID:  current.srcID,
				TriggerX:   current.x,
				TriggerY:   current.y,
				ChainDepth: current.depth,
				KillCount:  entry.KillCount,
				Reward:     entry.Reward,
			},
		})
	}

	// 更新玩家狀態
	g.sendPlayerUpdate(p)

	// 廣播連鎖結束
	chainDepth := len(entries)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgChainBomb,
		Payload: ws.ChainBombPayload{
			Phase:       "chain_result",
			TriggerID:   triggerID,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			ChainDepth:  chainDepth,
			TotalKills:  totalKills,
			TotalReward: totalReward,
		},
	})

	log.Printf("[ChainBomb] player=%s chain=%d kills=%d reward=%d",
		p.ID, chainDepth, totalKills, totalReward)

	// 全服公告（≥4 個擊破才公告）
	if totalKills >= ChainBombAnnounceMinKills {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgAnnouncement,
			Payload: map[string]interface{}{
				"event_type": "chain_bomb",
				"content":    chainBombAnnounceMsg(p.DisplayName, chainDepth, totalKills, totalReward),
				"color":      chainBombColor(chainDepth),
				"duration":   5,
				"priority":   4,
			},
		})
	}
}

// chainBombAnnounceMsg 格式化連鎖爆炸公告訊息
func chainBombAnnounceMsg(playerName string, chainDepth, totalKills, totalReward int) string {
	icon := "💥"
	if chainDepth >= 3 {
		icon = "🔥💥"
	}
	if chainDepth >= 5 {
		icon = "🌟💥💥"
	}
	return icon + " " + playerName + " 觸發 " + chainBombItoa(chainDepth) + " 層連鎖爆炸！擊破 " +
		chainBombItoa(totalKills) + " 個目標，獎勵 " + chainBombItoa(totalReward) + " 金幣！"
}

// chainBombColor 依連鎖層數決定公告顏色
func chainBombColor(chainDepth int) string {
	switch {
	case chainDepth >= 5:
		return "#FFD700" // 金色（5層連鎖）
	case chainDepth >= 3:
		return "#FF6600" // 橙色（3-4層）
	default:
		return "#FF4444" // 紅色（1-2層）
	}
}

// chainBombItoa 整數轉字串（獨立命名防止衝突）
func chainBombItoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	buf := make([]byte, 20)
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
