// lucky_chain_lightning_handler.go — T106 幸運連鎖閃電魚系統
// server-event-agent 負責維護
// 業界依據：Royal Fishing Jili「Thunder Shark — chain lightning that shocks nearby fish consecutively」
// 設計：擊破 T106 後，連鎖閃電攻擊附近 3 條魚（HP -50%），每條命中 ×1.5 累積倍率（最高 ×4.5）
// 若 3 條全部命中 → 「完美連鎖」：全服廣播 + 觸發玩家額外 ×2.0 加成
package game

import (
	"log"
	"math"
	"math/rand"
	"time"

	"chiikawa-game/internal/protocol"
)

// luckyChainLightningManager 管理連鎖閃電系統
type luckyChainLightningManager struct {
	// 個人冷卻：同一玩家 15 秒內不能再次觸發
	playerCooldowns map[string]time.Time
	// 全服冷卻：25 秒內全服只能觸發一次
	globalCooldown time.Time
}

func newLuckyChainLightningManager() *luckyChainLightningManager {
	return &luckyChainLightningManager{
		playerCooldowns: make(map[string]time.Time),
	}
}

// isLuckyChainLightningFish 判斷是否為連鎖閃電魚
func isLuckyChainLightningFish(defID string) bool {
	return defID == "T106"
}

// canTrigger 判斷是否可以觸發
func (m *luckyChainLightningManager) canTrigger(playerID string) bool {
	now := time.Now()
	if now.Before(m.globalCooldown) {
		return false
	}
	if cd, ok := m.playerCooldowns[playerID]; ok {
		if now.Before(cd) {
			return false
		}
	}
	return true
}

// tryLuckyChainLightning 嘗試觸發連鎖閃電
// 在 handleKill 中，擊破 T106 後呼叫此函數
func (g *Game) tryLuckyChainLightning(playerID string, killerName string) {
	m := g.luckyChainLightning
	if !m.canTrigger(playerID) {
		return
	}

	// 設定冷卻
	now := time.Now()
	m.playerCooldowns[playerID] = now.Add(15 * time.Second)
	m.globalCooldown = now.Add(25 * time.Second)

	// 廣播觸發事件
	g.hub.Broadcast(protocol.MsgLuckyChainLightning, protocol.LuckyChainLightningPayload{
		Event:       "trigger",
		TriggerID:   playerID,
		TriggerName: killerName,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "⚡ " + killerName + " 觸發連鎖閃電！",
		Priority: "high",
		Color:    "#FFD700",
	})

	// 在 goroutine 中執行連鎖閃電邏輯
	go g.runChainLightning(playerID, killerName)
}

// runChainLightning 執行連鎖閃電邏輯
func (g *Game) runChainLightning(playerID string, killerName string) {
	// 等待 300ms 讓 Client 顯示觸發動畫
	time.Sleep(300 * time.Millisecond)

	g.mu.Lock()
	p, ok := g.players[playerID]
	if !ok {
		g.mu.Unlock()
		return
	}
	betCost := p.GetBetDef().BetCost

	// 找附近最多 3 個目標（排除 BOSS）
	type targetInfo struct {
		id   string
		dist float64
	}
	var candidates []targetInfo
	// 以場地中心為基準，找最近的目標
	centerX := GameWidth / 2
	centerY := GameHeight / 2
	for id, t := range g.targets {
		if t.Def.Type == "boss" {
			continue
		}
		dx := t.X - centerX
		dy := t.Y - centerY
		dist := math.Sqrt(dx*dx + dy*dy)
		candidates = append(candidates, targetInfo{id: id, dist: dist})
	}
	g.mu.Unlock()

	// 隨機選最多 3 個（不按距離，模擬閃電隨機跳躍）
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	maxChain := 3
	if len(candidates) < maxChain {
		maxChain = len(candidates)
	}
	selected := candidates[:maxChain]

	// 逐一執行閃電攻擊（每 400ms 一次）
	hitTargets := []string{}
	accumMult := 1.0
	totalReward := 0

	for i, c := range selected {
		time.Sleep(400 * time.Millisecond)

		g.mu.Lock()
		t, exists := g.targets[c.id]
		if !exists {
			g.mu.Unlock()
			continue
		}

		// HP -50%
		damage := t.MaxHP / 2
		t.HP -= damage
		if t.HP < 1 {
			t.HP = 1
		}

		// 廣播 HP 更新
		g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
			InstanceID: t.InstanceID,
			HP:         t.HP,
			MaxHP:      t.MaxHP,
			X:          t.X,
			Y:          t.Y,
		})

		hitTargets = append(hitTargets, c.id)
		accumMult += 1.5 // 每次命中 +1.5x
		chainReward := int(float64(betCost) * 1.5)
		totalReward += chainReward

		// 廣播連鎖命中
		g.hub.Broadcast(protocol.MsgLuckyChainLightning, protocol.LuckyChainLightningPayload{
			Event:       "chain_hit",
			TriggerID:   playerID,
			TriggerName: killerName,
			HitTargets:  hitTargets,
			ChainCount:  i + 1,
			TotalReward: totalReward,
			Multiplier:  accumMult,
		})
		g.mu.Unlock()
	}

	// 結算
	g.mu.Lock()
	p2, ok2 := g.players[playerID]
	if ok2 {
		p2.AddCoins(totalReward)
		g.sendPlayerUpdate(playerID)
	}

	// 完美連鎖（3 條全部命中）
	isPerfect := len(hitTargets) == 3
	if isPerfect {
		// 額外 ×2.0 加成
		bonusReward := int(float64(betCost) * 2.0)
		if ok2 {
			p2.AddCoins(bonusReward)
			totalReward += bonusReward
		}
		g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
			Message:  "⚡✨ 完美連鎖！" + killerName + " 獲得額外 ×2.0 加成！",
			Priority: "high",
			Color:    "#00E5FF",
		})
	}
	g.mu.Unlock()

	// 廣播結算
	g.hub.Broadcast(protocol.MsgLuckyChainLightning, protocol.LuckyChainLightningPayload{
		Event:       "settle",
		TriggerID:   playerID,
		TriggerName: killerName,
		HitTargets:  hitTargets,
		ChainCount:  len(hitTargets),
		TotalReward: totalReward,
		Multiplier:  accumMult,
	})

	log.Printf("[ChainLightning] Player %s: %d chains, mult=%.1f, reward=%d, perfect=%v",
		playerID, len(hitTargets), accumMult, totalReward, isPerfect)
}
