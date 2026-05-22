// electric_jellyfish_handler.go — 電流水母電流網路系統 handler（DAY-193）
// 業界依據：King of Ocean 2026「electric jellyfish chains current between adjacent targets,
// paying multipliers from every link in the chain」
// T151 電流水母機制：
//   1. 擊破 T151 後，掃描場上所有目標，建立「電流網路」
//   2. 相鄰目標（200px 內）之間建立電流連接（每對目標只連接一次）
//   3. 每條電流連接有 65% 機率擊破兩端目標之一（較低 HP 的那個），獎勵 0.55x 倍率
//   4. 全服廣播每條電流連接，讓所有玩家看到「電流在場上形成網路」
//   5. 結果廣播：連接數/擊破數/總獎勵
// 設計差異：
//   - 與閃電鰻（隨機跳躍，單一路徑，5跳）不同，電流水母是「網路拓撲」（所有相鄰目標同時建立連接），
//     讓玩家看到「電流在整個場上形成網路」的壯觀視覺
//   - 與雷霆鯊魚（全場隨機跳躍，20跳）不同，電流水母是「距離驅動」（只連接相鄰目標），
//     密集的目標群會形成更多連接，製造「越多魚越爽」的策略感
//   - 電流連接是「雙向的」（A→B 和 B→A 是同一條連接），不會重複計算
//   - 每條連接只擊破一個目標（較低 HP 的那個），讓玩家感受到「電流選擇弱者」的自然感
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// ElectricJellyfishRadius 電流連接半徑（px）— 相鄰目標的最大距離
	ElectricJellyfishRadius = 200.0
	// ElectricJellyfishKillChance 每條電流連接的擊破機率（65%）
	ElectricJellyfishKillChance = 0.65
	// ElectricJellyfishRewardMult 電流擊破獎勵倍率（比直接擊破低，平衡 RTP）
	ElectricJellyfishRewardMult = 0.55
	// ElectricJellyfishIntervalMs 每條電流連接的廣播間隔（ms）
	ElectricJellyfishIntervalMs = 80
	// ElectricJellyfishCooldownSec 全服冷卻時間（秒）
	ElectricJellyfishCooldownSec = 30
	// ElectricJellyfishAnnounceMinKills 全服公告最低擊破數
	ElectricJellyfishAnnounceMinKills = 5
	// ElectricJellyfishAnnounceMinLinks 全服公告最低連接數
	ElectricJellyfishAnnounceMinLinks = 8
)

// isElectricJellyfish 判斷是否為電流水母（T151）
func isElectricJellyfish(defID string) bool {
	return defID == "T151"
}

// electricLink 電流連接（兩個目標之間）
type electricLink struct {
	idA, idB       string
	xA, yA         float64
	xB, yB         float64
	multA, multB   float64
	hpA, hpB       int
	isKill         bool
	killedID       string
	killedMult     float64
	reward         int
}

// tryElectricJellyfishNetwork 擊破 T151 後觸發電流網路（DAY-193）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryElectricJellyfishNetwork(p *player.Player, triggerID string, triggerX, triggerY float64) {
	// 全服冷卻檢查（使用 GoldenJellyfish 的冷卻機制，獨立管理）
	// 注意：電流水母有自己的冷卻，不依賴其他系統
	log.Printf("[ElectricJellyfish] player=%s triggered at (%.0f,%.0f)", p.ID, triggerX, triggerY)

	// 收集場上所有存活目標
	g.mu.RLock()
	type targetInfo struct {
		instanceID string
		defID      string
		x, y       float64
		multiplier float64
		hp         int
	}
	var targets []targetInfo
	for id, t := range g.Targets {
		if id == triggerID || !t.IsAlive || t.DefID == "B001" {
			continue
		}
		targets = append(targets, targetInfo{
			instanceID: id,
			defID:      t.DefID,
			x:          t.X,
			y:          t.Y,
			multiplier: t.Multiplier,
			hp:         t.HP,
		})
	}
	g.mu.RUnlock()

	if len(targets) < 2 {
		log.Printf("[ElectricJellyfish] not enough targets (%d), skipping", len(targets))
		return
	}

	// 建立電流網路：找出所有相鄰目標對（距離 ≤ 200px）
	var links []electricLink
	usedPairs := make(map[string]bool)

	for i := 0; i < len(targets); i++ {
		for j := i + 1; j < len(targets); j++ {
			a := targets[i]
			b := targets[j]
			dx := a.x - b.x
			dy := a.y - b.y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > ElectricJellyfishRadius {
				continue
			}
			// 確保每對只計算一次
			pairKey := a.instanceID + ":" + b.instanceID
			if usedPairs[pairKey] {
				continue
			}
			usedPairs[pairKey] = true

			links = append(links, electricLink{
				idA:   a.instanceID,
				idB:   b.instanceID,
				xA:    a.x,
				yA:    a.y,
				xB:    b.x,
				yB:    b.y,
				multA: a.multiplier,
				multB: b.multiplier,
				hpA:   a.hp,
				hpB:   b.hp,
			})
		}
	}

	if len(links) == 0 {
		log.Printf("[ElectricJellyfish] no adjacent targets found, skipping")
		return
	}

	// 按距離排序（近的先連接，視覺上更自然）
	sort.Slice(links, func(i, j int) bool {
		dxi := links[i].xA - links[i].xB
		dyi := links[i].yA - links[i].yB
		dxj := links[j].xA - links[j].xB
		dyj := links[j].yA - links[j].yB
		return math.Sqrt(dxi*dxi+dyi*dyi) < math.Sqrt(dxj*dxj+dyj*dyj)
	})

	log.Printf("[ElectricJellyfish] found %d links among %d targets", len(links), len(targets))

	// 廣播電流網路開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgElectricJellyfish,
		Payload: ws.ElectricJellyfishPayload{
			Phase:      "network_start",
			TriggerID:  triggerID,
			TriggerX:   triggerX,
			TriggerY:   triggerY,
			KillerID:   p.ID,
			KillerName: p.DisplayName,
			LinkCount:  len(links),
		},
	})

	// 逐條電流連接執行
	totalReward := 0
	totalKills := 0
	killedIDs := make(map[string]bool) // 追蹤已擊破的目標（防止重複）

	var resultLinks []ws.ElectricLinkResult

	for i, link := range links {
		if i > 0 {
			time.Sleep(time.Duration(ElectricJellyfishIntervalMs) * time.Millisecond)
		}

		// 跳過已擊破的目標
		if killedIDs[link.idA] && killedIDs[link.idB] {
			continue
		}

		// 65% 機率觸發電流擊破
		isKill := rand.Float64() < ElectricJellyfishKillChance
		killedID := ""
		killedMult := 0.0
		reward := 0

		if isKill {
			// 選擇較低 HP 的目標擊破（電流選擇弱者）
			targetToKill := link.idA
			targetMult := link.multA
			if link.hpB < link.hpA && !killedIDs[link.idB] {
				targetToKill = link.idB
				targetMult = link.multB
			}
			if killedIDs[targetToKill] {
				// 已擊破，改選另一個
				if targetToKill == link.idA && !killedIDs[link.idB] {
					targetToKill = link.idB
					targetMult = link.multB
				} else if targetToKill == link.idB && !killedIDs[link.idA] {
					targetToKill = link.idA
					targetMult = link.multA
				} else {
					isKill = false
				}
			}

			if isKill {
				g.mu.Lock()
				t, ok := g.Targets[targetToKill]
				if !ok || !t.IsAlive {
					g.mu.Unlock()
					isKill = false
				} else {
					reward = int(float64(p.BetLevel) * targetMult * ElectricJellyfishRewardMult)
					if reward < 1 {
						reward = 1
					}
					t.IsAlive = false
					delete(g.Targets, targetToKill)
					g.mu.Unlock()

					killedID = targetToKill
					killedMult = targetMult
					killedIDs[targetToKill] = true
					totalReward += reward
					totalKills++

					// 廣播目標擊破
					g.Hub.Broadcast(&ws.Message{
						Type: ws.MsgTargetKill,
						Payload: ws.TargetKillPayload{
							InstanceID: targetToKill,
							KillerID:   p.ID,
							Reward:     reward,
							Multiplier: targetMult,
						},
					})
				}
			}
		}

		resultLinks = append(resultLinks, ws.ElectricLinkResult{
			IDA:        link.idA,
			IDB:        link.idB,
			XA:         link.xA,
			YA:         link.yA,
			XB:         link.xB,
			YB:         link.yB,
			IsKill:     isKill,
			KilledID:   killedID,
			KilledMult: killedMult,
			Reward:     reward,
		})

		// 廣播電流連接（全服）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgElectricJellyfish,
			Payload: ws.ElectricJellyfishPayload{
				Phase:      fmt.Sprintf("link_%d", i+1),
				LinkIndex:  i + 1,
				IDA:        link.idA,
				IDB:        link.idB,
				XA:         link.xA,
				YA:         link.yA,
				XB:         link.xB,
				YB:         link.yB,
				IsKill:     isKill,
				KilledID:   killedID,
				KilledMult: killedMult,
				Reward:     reward,
				KillerID:   p.ID,
			},
		})

		log.Printf("[ElectricJellyfish] link[%d] A=%s B=%s kill=%v reward=%d",
			i+1, link.idA, link.idB, isKill, reward)
	}

	// 發放總獎勵
	if totalReward > 0 {
		p.AddCoins(totalReward)
		g.sendPlayerUpdate(p)
	}

	// 廣播電流網路結果（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgElectricJellyfish,
		Payload: ws.ElectricJellyfishPayload{
			Phase:       "network_result",
			TriggerID:   triggerID,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			LinkCount:   len(links),
			TotalKills:  totalKills,
			TotalReward: totalReward,
			Links:       resultLinks,
		},
	})

	// 個人結果通知
	if totalReward > 0 {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgReward,
			Payload: ws.RewardPayload{
				Source:     "electric_jellyfish",
				Amount:     totalReward,
				Multiplier: float64(totalKills),
				NewBalance: p.Coins,
			},
		})
	}

	// 全服公告：擊破 ≥5 個 或 連接 ≥8 條
	if totalKills >= ElectricJellyfishAnnounceMinKills || len(links) >= ElectricJellyfishAnnounceMinLinks {
		g.announceElectricJellyfish(p.DisplayName, len(links), totalKills, totalReward)
	}

	log.Printf("[ElectricJellyfish] player=%s links=%d kills=%d total_reward=%d",
		p.ID, len(links), totalKills, totalReward)
}

// announceElectricJellyfish 全服公告電流水母電流網路（DAY-193）
func (g *Game) announceElectricJellyfish(playerName string, links, kills, reward int) {
	icon := "⚡🪼"
	if links >= 12 {
		icon = "🌐⚡"
	}
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "electric_jellyfish",
			"message":    fmt.Sprintf("%s %s 的電流水母建立 %d 條電流連接！擊破 %d 個目標！獲得 %d 金幣！", icon, playerName, links, kills, reward),
			"color":      "#00FFFF", // 青色（電流感）
			"duration":   5.0,
			"priority":   3,
		},
	})
}
