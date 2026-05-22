// anglerfish_electric_handler.go — 巨型鮟鱇魚電擊寶箱系統（DAY-196）
// 業界依據：JILI Mega Fishing「Giant Anglerfish can shoot electricity to open treasure chests,
// giant crocodiles awaken to hunt fish on the fish farm to accumulate big prizes!」
// 設計：T154 巨型鮟鱇魚出現後，每 3 秒用「電擊」攻擊場上一個隨機目標：
//   1. 若命中 T102 寶箱怪 → 強制開箱：給觸發玩家 3-5x 倍率加成（寶箱獎勵）
//   2. 若命中普通目標 → 70% 擊破機率（0.60x 倍率），全服共享獎勵
//   3. 每次電擊有 5% 機率「超級電擊」→ 全場所有目標同時受到電擊
//   4. 玩家擊破鮟鱇魚本身 → 獲得鮟鱇魚倍率 + 累積電擊獎池 40%
//   5. 全服廣播每次電擊，讓所有玩家看到「鮟鱇魚在電擊全場」
// 設計差異：
//   - 與巨型鱷魚獵食（每 2 秒獵食，累積獎池）不同，鮟鱇魚是「電擊型」，
//     有「超級電擊」全場清場的爆發感，且能強制開寶箱製造驚喜
//   - 與閃電魚自動連鎖（時間驅動，8 秒）不同，鮟鱇魚是「目標驅動」（在場上持續存在），
//     玩家需要決策：讓它繼續電擊累積獎池，還是立刻擊破？
//   - 「強制開寶箱」機制是業界首創設計感：鮟鱇魚的電擊能打開寶箱怪，
//     讓玩家看到「鮟鱇魚幫我開寶箱」的驚喜感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// 巨型鮟鱇魚常數
const (
	AnglerfishCooldownSec      = 50    // 全服冷卻 50 秒
	AnglerfishZapIntervalMs    = 3000  // 每 3 秒電擊一次
	AnglerfishMaxZaps          = 8     // 最多電擊 8 次
	AnglerfishZapChance        = 0.70  // 普通目標擊破機率 70%
	AnglerfishZapMult          = 0.60  // 普通目標獎勵倍率 0.60x
	AnglerfishSuperZapChance   = 0.05  // 超級電擊機率 5%
	AnglerfishTreasureMultMin  = 3.0   // 寶箱開箱最小倍率 3x
	AnglerfishTreasureMultMax  = 5.0   // 寶箱開箱最大倍率 5x
	AnglerfishPoolSharePct     = 0.40  // 玩家擊破時獲得獎池 40%
)

// anglerfishManager 巨型鮟鱇魚管理器（全服共享）
type anglerfishManager struct {
	mu          sync.Mutex
	isActive    bool
	instanceID  string    // 當前鮟鱇魚實例 ID
	zapCount    int       // 已電擊次數
	totalPool   int       // 累積電擊獎池
	cooldownEnd time.Time
}

func newAnglerfishManager() *anglerfishManager {
	return &anglerfishManager{}
}

// isAnglerfishElectric 判斷是否為巨型鮟鱇魚（T154）
func isAnglerfishElectric(defID string) bool {
	return defID == "T154"
}

// notifyAnglerfishSpawn T154 生成時觸發電擊模式
func (g *Game) notifyAnglerfishSpawn(instanceID string) {
	mgr := g.AnglerfishElectric
	mgr.mu.Lock()

	// 全服冷卻檢查
	if mgr.isActive || time.Now().Before(mgr.cooldownEnd) {
		mgr.mu.Unlock()
		return
	}
	mgr.isActive = true
	mgr.instanceID = instanceID
	mgr.zapCount = 0
	mgr.totalPool = 0
	mgr.mu.Unlock()

	// 廣播鮟鱇魚出現（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnglerfishElectric,
		Payload: ws.AnglerfishElectricPayload{
			Phase:      "anglerfish_appear",
			InstanceID: instanceID,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventMegaWin, "巨型鮟鱇魚", 0, map[string]string{
		"message": "⚡ 巨型鮟鱇魚出現！牠的電擊能強制開啟寶箱！快去打！",
	})
	g.broadcastAnnouncement(ann)

	log.Printf("[AnglerfishElectric] instance=%s spawned, starting zap loop", instanceID)

	// 啟動電擊循環
	go g.runAnglerfishZapLoop(instanceID)
}

// runAnglerfishZapLoop 電擊循環（goroutine）
func (g *Game) runAnglerfishZapLoop(instanceID string) {
	for zapIdx := 1; zapIdx <= AnglerfishMaxZaps; zapIdx++ {
		time.Sleep(AnglerfishZapIntervalMs * time.Millisecond)

		// 確認鮟鱇魚還在場上
		mgr := g.AnglerfishElectric
		mgr.mu.Lock()
		if !mgr.isActive || mgr.instanceID != instanceID {
			mgr.mu.Unlock()
			return
		}
		mgr.mu.Unlock()

		g.mu.RLock()
		_, stillAlive := g.Targets[instanceID]
		g.mu.RUnlock()
		if !stillAlive {
			// 鮟鱇魚已被擊破，停止電擊
			return
		}

		// 判斷是否超級電擊
		isSuperZap := rand.Float64() < AnglerfishSuperZapChance

		if isSuperZap {
			g.doAnglerfishSuperZap(instanceID, zapIdx)
		} else {
			g.doAnglerfishSingleZap(instanceID, zapIdx)
		}
	}

	// 達到最大電擊次數，鮟鱇魚離開
	g.onAnglerfishLeave(instanceID, "max_zaps")
}

// doAnglerfishSingleZap 單次電擊（選一個目標）
func (g *Game) doAnglerfishSingleZap(instanceID string, zapIdx int) {
	// 選擇目標：優先選 T102 寶箱怪，否則隨機選普通目標
	type zapTarget struct {
		instanceID string
		defID      string
		multiplier float64
		x, y       float64
	}

	g.mu.RLock()
	var treasureTargets []zapTarget
	var normalTargets []zapTarget
	for _, t := range g.Targets {
		if t.InstanceID == instanceID {
			continue // 跳過鮟鱇魚自身
		}
		info := zapTarget{
			instanceID: t.InstanceID,
			defID:      t.DefID,
			multiplier: t.Multiplier,
			x:          t.X,
			y:          t.Y,
		}
		if t.DefID == "T102" {
			treasureTargets = append(treasureTargets, info)
		} else {
			normalTargets = append(normalTargets, info)
		}
	}
	g.mu.RUnlock()

	// 優先選寶箱怪
	var chosen *zapTarget
	isTreasure := false
	if len(treasureTargets) > 0 {
		t := treasureTargets[rand.Intn(len(treasureTargets))]
		chosen = &t
		isTreasure = true
	} else if len(normalTargets) > 0 {
		t := normalTargets[rand.Intn(len(normalTargets))]
		chosen = &t
	}

	if chosen == nil {
		// 場上沒有目標，廣播空電擊
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgAnglerfishElectric,
			Payload: ws.AnglerfishElectricPayload{
				Phase:      fmt.Sprintf("zap_%d", zapIdx),
				InstanceID: instanceID,
				ZapIndex:   zapIdx,
				IsEmpty:    true,
			},
		})
		return
	}

	// 執行電擊
	isKill := false
	reward := 0
	isTreasureOpen := false
	treasureMult := 0.0

	if isTreasure {
		// 強制開箱：給觸發玩家（取第一個玩家）寶箱獎勵
		isTreasureOpen = true
		treasureMult = AnglerfishTreasureMultMin + rand.Float64()*(AnglerfishTreasureMultMax-AnglerfishTreasureMultMin)

		// 找一個玩家給獎勵（取平均 betLevel）
		avgBet := g.getAverageBetLevelForAnglerfish()
		reward = int(treasureMult * float64(avgBet))

		// 移除寶箱怪
		g.mu.Lock()
		if _, ok := g.Targets[chosen.instanceID]; ok {
			delete(g.Targets, chosen.instanceID)
			isKill = true
			// 給所有玩家分享寶箱獎勵
			for _, p := range g.Players {
				p.Coins += reward / max(1, len(g.Players))
			}
		}
		g.mu.Unlock()

		// 更新獎池
		mgr := g.AnglerfishElectric
		mgr.mu.Lock()
		mgr.zapCount++
		mgr.totalPool += reward
		mgr.mu.Unlock()

		log.Printf("[AnglerfishElectric] zap_%d: TREASURE OPEN! mult=%.1fx reward=%d",
			zapIdx, treasureMult, reward)
	} else {
		// 普通電擊：70% 擊破機率
		if rand.Float64() < AnglerfishZapChance {
			isKill = true
			avgBet := g.getAverageBetLevelForAnglerfish()
			reward = int(chosen.multiplier * float64(avgBet) * AnglerfishZapMult)

			g.mu.Lock()
			if _, ok := g.Targets[chosen.instanceID]; ok {
				delete(g.Targets, chosen.instanceID)
				// 全服共享獎勵
				for _, p := range g.Players {
					p.Coins += reward / max(1, len(g.Players))
				}
			}
			g.mu.Unlock()

			// 更新獎池
			mgr := g.AnglerfishElectric
			mgr.mu.Lock()
			mgr.zapCount++
			mgr.totalPool += reward
			mgr.mu.Unlock()
		}
	}

	// 廣播電擊結果（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnglerfishElectric,
		Payload: ws.AnglerfishElectricPayload{
			Phase:          fmt.Sprintf("zap_%d", zapIdx),
			InstanceID:     instanceID,
			ZapIndex:       zapIdx,
			TargetID:       chosen.instanceID,
			TargetDefID:    chosen.defID,
			TargetX:        chosen.x,
			TargetY:        chosen.y,
			IsKill:         isKill,
			IsTreasure:     isTreasureOpen,
			TreasureMult:   treasureMult,
			ZapReward:      reward,
			IsSuperZap:     false,
		},
	})

	// 寶箱開箱全服公告
	if isTreasureOpen {
		ann := g.Announce.Create(announce.EventMegaWin, "鮟鱇魚電擊", reward, map[string]string{
			"message": fmt.Sprintf("⚡💰 鮟鱇魚電擊開箱！寶箱爆出 %.1fx 大獎！全服共享 %d 金幣！",
				treasureMult, reward),
		})
		g.broadcastAnnouncement(ann)
	}
}

// doAnglerfishSuperZap 超級電擊（全場所有目標同時受到電擊）
func (g *Game) doAnglerfishSuperZap(instanceID string, zapIdx int) {
	type superZapTarget struct {
		instanceID string
		defID      string
		multiplier float64
		x, y       float64
	}

	g.mu.RLock()
	var targets []superZapTarget
	for _, t := range g.Targets {
		if t.InstanceID == instanceID {
			continue
		}
		targets = append(targets, superZapTarget{
			instanceID: t.InstanceID,
			defID:      t.DefID,
			multiplier: t.Multiplier,
			x:          t.X,
			y:          t.Y,
		})
	}
	g.mu.RUnlock()

	// 廣播超級電擊開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnglerfishElectric,
		Payload: ws.AnglerfishElectricPayload{
			Phase:      "super_zap_start",
			InstanceID: instanceID,
			ZapIndex:   zapIdx,
			IsSuperZap: true,
			TargetCount: len(targets),
		},
	})

	// 全服公告超級電擊
	ann := g.Announce.Create(announce.EventMegaWin, "超級電擊", 0, map[string]string{
		"message": fmt.Sprintf("⚡⚡⚡ 鮟鱇魚超級電擊！全場 %d 個目標同時受到電擊！", len(targets)),
	})
	g.broadcastAnnouncement(ann)

	superKills := 0
	superReward := 0
	avgBet := g.getAverageBetLevelForAnglerfish()

	for i, t := range targets {
		time.Sleep(80 * time.Millisecond) // 每個目標間隔 80ms，製造「電流蔓延」感

		isKill := false
		reward := 0
		isTreasure := t.defID == "T102"

		if isTreasure {
			// 寶箱強制開箱
			treasureMult := AnglerfishTreasureMultMin + rand.Float64()*(AnglerfishTreasureMultMax-AnglerfishTreasureMultMin)
			reward = int(treasureMult * float64(avgBet))
			g.mu.Lock()
			if _, ok := g.Targets[t.instanceID]; ok {
				delete(g.Targets, t.instanceID)
				isKill = true
				for _, p := range g.Players {
					p.Coins += reward / max(1, len(g.Players))
				}
			}
			g.mu.Unlock()
		} else if rand.Float64() < AnglerfishZapChance {
			// 普通目標 70% 擊破
			reward = int(t.multiplier * float64(avgBet) * AnglerfishZapMult)
			g.mu.Lock()
			if _, ok := g.Targets[t.instanceID]; ok {
				delete(g.Targets, t.instanceID)
				isKill = true
				for _, p := range g.Players {
					p.Coins += reward / max(1, len(g.Players))
				}
			}
			g.mu.Unlock()
		}

		if isKill {
			superKills++
			superReward += reward
		}

		// 廣播每個目標的電擊結果
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgAnglerfishElectric,
			Payload: ws.AnglerfishElectricPayload{
				Phase:       fmt.Sprintf("super_zap_%d", i+1),
				InstanceID:  instanceID,
				ZapIndex:    zapIdx,
				TargetID:    t.instanceID,
				TargetDefID: t.defID,
				TargetX:     t.x,
				TargetY:     t.y,
				IsKill:      isKill,
				IsTreasure:  isTreasure,
				ZapReward:   reward,
				IsSuperZap:  true,
			},
		})
	}

	// 更新獎池
	mgr := g.AnglerfishElectric
	mgr.mu.Lock()
	mgr.zapCount++
	mgr.totalPool += superReward
	mgr.mu.Unlock()

	// 廣播超級電擊結果（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnglerfishElectric,
		Payload: ws.AnglerfishElectricPayload{
			Phase:       "super_zap_result",
			InstanceID:  instanceID,
			ZapIndex:    zapIdx,
			IsSuperZap:  true,
			SuperKills:  superKills,
			SuperReward: superReward,
		},
	})

	// 全服公告超級電擊結果（≥3 個擊破）
	if superKills >= 3 {
		ann2 := g.Announce.Create(announce.EventMegaWin, "超級電擊", superReward, map[string]string{
			"message": fmt.Sprintf("⚡💥 超級電擊結算！擊破 %d 個目標！全服共享 %d 金幣！",
				superKills, superReward),
		})
		g.broadcastAnnouncement(ann2)
	}

	log.Printf("[AnglerfishElectric] super_zap_%d: kills=%d reward=%d", zapIdx, superKills, superReward)
}

// notifyAnglerfishKill 玩家擊破鮟鱇魚時結算
func (g *Game) notifyAnglerfishKill(p *player.Player, instanceID string, baseMult float64) {
	mgr := g.AnglerfishElectric
	mgr.mu.Lock()
	if !mgr.isActive || mgr.instanceID != instanceID {
		mgr.mu.Unlock()
		return
	}
	zapCount := mgr.zapCount
	totalPool := mgr.totalPool
	mgr.isActive = false
	mgr.cooldownEnd = time.Now().Add(AnglerfishCooldownSec * time.Second)
	mgr.mu.Unlock()

	// 計算獎勵：基礎倍率 + 獎池 40%
	baseReward := int(baseMult * float64(p.BetLevel))
	poolBonus := int(float64(totalPool) * AnglerfishPoolSharePct)
	totalReward := baseReward + poolBonus

	g.mu.Lock()
	p.Coins += totalReward
	g.mu.Unlock()

	// 廣播擊破結算（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnglerfishElectric,
		Payload: ws.AnglerfishElectricPayload{
			Phase:       "anglerfish_killed",
			InstanceID:  instanceID,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			ZapCount:    zapCount,
			TotalPool:   totalPool,
			PoolBonus:   poolBonus,
			BaseReward:  baseReward,
			TotalReward: totalReward,
		},
	})

	// 全服公告（獎池 > 0 或電擊次數 ≥ 3）
	if zapCount >= 3 || poolBonus > 0 {
		ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("⚡🎉 %s 擊破巨型鮟鱇魚！電擊 %d 次累積獎池 %d！獲得 %d 金幣！",
				p.DisplayName, zapCount, totalPool, totalReward),
		})
		g.broadcastAnnouncement(ann)
	}

	log.Printf("[AnglerfishElectric] player=%s killed anglerfish: zaps=%d pool=%d bonus=%d total=%d",
		p.ID, zapCount, totalPool, poolBonus, totalReward)
}

// onAnglerfishLeave 鮟鱇魚達到最大電擊次數或超時離開
func (g *Game) onAnglerfishLeave(instanceID string, reason string) {
	mgr := g.AnglerfishElectric
	mgr.mu.Lock()
	if !mgr.isActive || mgr.instanceID != instanceID {
		mgr.mu.Unlock()
		return
	}
	zapCount := mgr.zapCount
	totalPool := mgr.totalPool
	mgr.isActive = false
	mgr.cooldownEnd = time.Now().Add(AnglerfishCooldownSec * time.Second)
	mgr.mu.Unlock()

	// 廣播離開（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnglerfishElectric,
		Payload: ws.AnglerfishElectricPayload{
			Phase:      "anglerfish_leave",
			InstanceID: instanceID,
			ZapCount:   zapCount,
			TotalPool:  totalPool,
		},
	})

	log.Printf("[AnglerfishElectric] instance=%s left (%s): zaps=%d pool=%d",
		instanceID, reason, zapCount, totalPool)
}

// getAverageBetLevelForAnglerfish 計算全服平均 betLevel（供電擊獎勵計算使用）
func (g *Game) getAverageBetLevelForAnglerfish() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if len(g.Players) == 0 {
		return 1
	}
	total := 0
	for _, p := range g.Players {
		total += p.BetLevel
	}
	return total / len(g.Players)
}

// max 輔助函數（Go 1.21 前需要手動定義）
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
