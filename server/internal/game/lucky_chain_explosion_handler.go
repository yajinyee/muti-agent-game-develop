// lucky_chain_explosion_handler.go — 幸運連鎖爆炸魚系統（DAY-266）
// 業界原創「連鎖爆炸+空間擴散+三層引爆」機制
//
// 設計：擊破 T224 後，觸發「連鎖爆炸」：
//   - 第 1 層：隨機選場上 1 個目標引爆（×2.0 倍率，全服共享）
//   - 爆炸後，距離爆炸點 200px 內的所有目標 HP -50%，各自 40% 機率「二次引爆」（×1.5 倍率）
//   - 第 2 層：二次引爆目標對 150px 內目標 HP -30%，各自 25% 機率「三次引爆」（×1.2 倍率）
//   - 第 3 層：三次引爆目標直接擊破（×1.0 倍率，全服共享）
//   - 最多 3 層連鎖；個人冷卻 22 秒；全服冷卻 35 秒
//
// 設計差異：
//   - 與隕石雨（T211，隨機轟炸）不同，連鎖爆炸是「空間擴散」，讓玩家看到「爆炸從一點向外擴散」的視覺爽感
//   - 「三層連鎖」讓玩家有「一炸帶一片，一片再帶一片」的連鎖快感
//   - 「40%/25% 機率引爆」讓每次爆炸都有不確定性，製造「這次會不會連鎖」的期待感
//   - 「距離衰減（200→150px）」讓爆炸有「中心強、邊緣弱」的真實感
//   - 「全服廣播爆炸位置和連鎖數」讓所有玩家看到「爆炸在哪裡、連鎖了幾層」，製造社交感
//   - 業界依據：Ocean King 3 的 Chain Explosion 系統，2026 年最熱門 AOE 連鎖方向
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyChainExplosionPersonalCD  = 22 * time.Second // 個人冷卻
	LuckyChainExplosionGlobalCD    = 35 * time.Second // 全服冷卻
	LuckyChainExplosionLayer1Mult  = 2.0              // 第 1 層引爆倍率
	LuckyChainExplosionLayer2Mult  = 1.5              // 第 2 層引爆倍率
	LuckyChainExplosionLayer3Mult  = 1.2              // 第 3 層引爆倍率
	LuckyChainExplosionLayer1Range = 200.0            // 第 1 層爆炸範圍（px）
	LuckyChainExplosionLayer2Range = 150.0            // 第 2 層爆炸範圍（px）
	LuckyChainExplosionLayer1Prob  = 0.40             // 第 2 層引爆機率
	LuckyChainExplosionLayer2Prob  = 0.25             // 第 3 層引爆機率
	LuckyChainExplosionLayer1HP    = 0.50             // 第 1 層 HP 傷害比例
	LuckyChainExplosionLayer2HP    = 0.30             // 第 2 層 HP 傷害比例
	LuckyChainExplosionMaxLayers   = 3                // 最大連鎖層數
)

// chainExplosionEvent 連鎖爆炸事件記錄（用於廣播）
type chainExplosionEvent struct {
	layer      int
	targetName string
	x, y       float64
	mult       float64
	reward     int
}

// luckyChainExplosionManager 幸運連鎖爆炸魚管理器
type luckyChainExplosionManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time
}

func newLuckyChainExplosionManager() *luckyChainExplosionManager {
	return &luckyChainExplosionManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyChainExplosionFish 判斷是否為幸運連鎖爆炸魚
func isLuckyChainExplosionFish(defID string) bool {
	return defID == "T224"
}

// tryLuckyChainExplosionFish 擊破 T224 後觸發連鎖爆炸
func (g *Game) tryLuckyChainExplosionFish(p *player.Player) {
	m := g.LuckyChainExplosion

	m.mu.Lock()
	now := time.Now()

	// 全服冷卻檢查
	if now.Before(m.globalCooldownUntil) {
		m.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if cd, ok := m.personalCooldowns[p.ID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyChainExplosionPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyChainExplosionGlobalCD)
	m.mu.Unlock()

	log.Printf("[ChainExplosion] player=%s 觸發連鎖爆炸！", p.ID)

	// 個人訊息：觸發者
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyChainExplosion,
		Payload: ws.LuckyChainExplosionPayload{
			Event:      "explosion_start",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Layer1Mult: LuckyChainExplosionLayer1Mult,
			Layer2Mult: LuckyChainExplosionLayer2Mult,
			Layer3Mult: LuckyChainExplosionLayer3Mult,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyChainExplosion,
		Payload: ws.LuckyChainExplosionPayload{
			Event:      "explosion_broadcast",
			PlayerName: p.DisplayName,
			Layer1Mult: LuckyChainExplosionLayer1Mult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyChainExplosion, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("💥 %s 觸發連鎖爆炸！最多 3 層連鎖！第 1 層 ×%.1f！",
			p.DisplayName, LuckyChainExplosionLayer1Mult),
		"color": "#FF4500",
	})

	// 執行連鎖爆炸（在 goroutine 中，避免阻塞）
	go g.runChainExplosion(p)
}

// runChainExplosion 執行連鎖爆炸主邏輯
func (g *Game) runChainExplosion(p *player.Player) {
	// 第 1 層：隨機選一個目標引爆
	g.mu.RLock()
	targets := make([]targetSnapshot, 0, len(g.Targets))
	for _, t := range g.Targets {
		if t.HP > 0 && !isLuckyChainExplosionFish(t.DefID) {
			name := t.DefID
			if t.Def != nil {
				name = t.Def.Name
			}
			targets = append(targets, targetSnapshot{
				instanceID: t.InstanceID,
				defID:      t.DefID,
				name:       name,
				x:          t.X,
				y:          t.Y,
				hp:         t.HP,
				multiplier: t.Multiplier,
			})
		}
	}
	g.mu.RUnlock()

	if len(targets) == 0 {
		log.Printf("[ChainExplosion] 場上無目標，連鎖爆炸取消")
		return
	}

	// 隨機選第 1 個引爆目標
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := rng.Intn(len(targets))
	layer1Target := targets[idx]

	// 第 1 層引爆
	reward1 := g.doExplosionKill(p, layer1Target.instanceID, layer1Target.multiplier, LuckyChainExplosionLayer1Mult)

	events := []chainExplosionEvent{
		{layer: 1, targetName: layer1Target.name, x: layer1Target.x, y: layer1Target.y,
			mult: LuckyChainExplosionLayer1Mult, reward: reward1},
	}

	// 廣播第 1 層爆炸
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyChainExplosion,
		Payload: ws.LuckyChainExplosionPayload{
			Event:      "explosion_layer",
			PlayerName: p.DisplayName,
			Layer:      1,
			TargetName: layer1Target.name,
			X:          layer1Target.x,
			Y:          layer1Target.y,
			Mult:       LuckyChainExplosionLayer1Mult,
			Reward:     reward1,
			Range:      LuckyChainExplosionLayer1Range,
		},
	})

	// 短暫延遲，讓 Client 有時間顯示第 1 層爆炸動畫
	time.Sleep(300 * time.Millisecond)

	// 第 2 層：找 200px 內的目標，各自 40% 機率引爆
	layer2Targets := g.findTargetsInRange(layer1Target.x, layer1Target.y,
		LuckyChainExplosionLayer1Range, layer1Target.instanceID)

	var layer3Candidates []targetSnapshot
	totalLayer2Reward := 0

	for _, t2 := range layer2Targets {
		// HP -50%
		g.applyExplosionDamage(t2.instanceID, LuckyChainExplosionLayer1HP)

		// 40% 機率二次引爆
		if rng.Float64() < LuckyChainExplosionLayer1Prob {
			reward2 := g.doExplosionKill(p, t2.instanceID, t2.multiplier, LuckyChainExplosionLayer2Mult)
			totalLayer2Reward += reward2
			events = append(events, chainExplosionEvent{
				layer: 2, targetName: t2.name, x: t2.x, y: t2.y,
				mult: LuckyChainExplosionLayer2Mult, reward: reward2,
			})
			layer3Candidates = append(layer3Candidates, t2)
		}
	}

	if len(layer2Targets) > 0 {
		// 廣播第 2 層爆炸
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyChainExplosion,
			Payload: ws.LuckyChainExplosionPayload{
				Event:         "explosion_layer",
				PlayerName:    p.DisplayName,
				Layer:         2,
				AffectedCount: len(layer2Targets),
				ExplodeCount:  len(layer3Candidates),
				Mult:          LuckyChainExplosionLayer2Mult,
				Reward:        totalLayer2Reward,
				Range:         LuckyChainExplosionLayer2Range,
				X:             layer1Target.x,
				Y:             layer1Target.y,
			},
		})
		time.Sleep(300 * time.Millisecond)
	}

	// 第 3 層：從二次引爆的目標出發，找 150px 內目標，各自 25% 機率三次引爆
	totalLayer3Reward := 0
	layer3ExplodeCount := 0

	for _, t2 := range layer3Candidates {
		layer3Targets := g.findTargetsInRange(t2.x, t2.y,
			LuckyChainExplosionLayer2Range, t2.instanceID)

		for _, t3 := range layer3Targets {
			// HP -30%
			g.applyExplosionDamage(t3.instanceID, LuckyChainExplosionLayer2HP)

			// 25% 機率三次引爆
			if rng.Float64() < LuckyChainExplosionLayer2Prob {
				reward3 := g.doExplosionKill(p, t3.instanceID, t3.multiplier, LuckyChainExplosionLayer3Mult)
				totalLayer3Reward += reward3
				layer3ExplodeCount++
				events = append(events, chainExplosionEvent{
					layer: 3, targetName: t3.name, x: t3.x, y: t3.y,
					mult: LuckyChainExplosionLayer3Mult, reward: reward3,
				})
			}
		}
	}

	if layer3ExplodeCount > 0 {
		// 廣播第 3 層爆炸
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyChainExplosion,
			Payload: ws.LuckyChainExplosionPayload{
				Event:        "explosion_layer",
				PlayerName:   p.DisplayName,
				Layer:        3,
				ExplodeCount: layer3ExplodeCount,
				Mult:         LuckyChainExplosionLayer3Mult,
				Reward:       totalLayer3Reward,
			},
		})
	}

	// 計算總連鎖層數和總獎勵
	totalLayers := 1
	if len(layer3Candidates) > 0 {
		totalLayers = 2
	}
	if layer3ExplodeCount > 0 {
		totalLayers = 3
	}
	totalReward := reward1 + totalLayer2Reward + totalLayer3Reward
	totalExploded := 1 + len(layer3Candidates) + layer3ExplodeCount

	log.Printf("[ChainExplosion] 連鎖完成！%d 層，共 %d 個目標引爆，總獎勵 %d",
		totalLayers, totalExploded, totalReward)

	// 最終結算廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyChainExplosion,
		Payload: ws.LuckyChainExplosionPayload{
			Event:        "explosion_result",
			PlayerName:   p.DisplayName,
			TotalLayers:  totalLayers,
			TotalExplode: totalExploded,
			TotalReward:  totalReward,
		},
	})

	// 連鎖達到 3 層時全服公告
	if totalLayers >= 3 {
		g.Announce.Create(announce.EventLuckyChainExplosion, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("💥 %s 連鎖爆炸達到 3 層！共 %d 個目標引爆！總獎勵 +%d！",
				p.DisplayName, totalExploded, totalReward),
			"color": "#FFD700",
		})
	}
}

// doExplosionKill 執行爆炸擊破（給予獎勵，廣播 kill 事件）
// 回傳實際獎勵
func (g *Game) doExplosionKill(p *player.Player, instanceID string, multiplier float64, multBonus float64) int {
	g.mu.Lock()
	t, ok := g.Targets[instanceID]
	if !ok || t.HP <= 0 {
		g.mu.Unlock()
		return 0
	}
	t.HP = 0
	g.mu.Unlock()

	// 計算獎勵（基礎獎勵 × 爆炸倍率加成）
	betDef := p.GetBetDef()
	betCost := 1
	if betDef != nil {
		betCost = betDef.BetCost
	}
	reward := int(float64(betCost) * multiplier * multBonus)
	if reward < 1 {
		reward = 1
	}

	// 給予玩家獎勵
	p.AddCoins(reward)

	// 廣播 kill 事件（讓 Client 顯示死亡動畫）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTargetKill,
		Payload: ws.TargetKillPayload{
			InstanceID: instanceID,
			KillerID:   p.ID,
			Reward:     reward,
			Multiplier: multiplier,
		},
	})

	// 移除目標
	g.mu.Lock()
	delete(g.Targets, instanceID)
	g.mu.Unlock()

	return reward
}

// applyExplosionDamage 對目標施加爆炸傷害（HP 百分比扣除）
func (g *Game) applyExplosionDamage(instanceID string, hpRatio float64) {
	g.mu.Lock()
	t, ok := g.Targets[instanceID]
	if !ok || t.HP <= 0 {
		g.mu.Unlock()
		return
	}
	damage := int(float64(t.HP) * hpRatio)
	if damage < 1 {
		damage = 1
	}
	t.HP -= damage
	if t.HP < 0 {
		t.HP = 0
	}
	hp := t.HP
	g.mu.Unlock()

	// 廣播 HP 更新
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTargetUpdate,
		Payload: ws.TargetUpdatePayload{
			InstanceID: instanceID,
			HP:         hp,
		},
	})
}

// targetSnapshot 目標快照（用於連鎖爆炸計算）
type targetSnapshot struct {
	instanceID string
	defID      string
	name       string
	x, y       float64
	hp         int
	multiplier float64
}

// findTargetsInRange 找指定範圍內的所有存活目標（排除指定 instanceID）
func (g *Game) findTargetsInRange(cx, cy, radius float64, excludeID string) []targetSnapshot {
	g.mu.RLock()
	defer g.mu.RUnlock()

	result := make([]targetSnapshot, 0)
	for _, t := range g.Targets {
		if t.InstanceID == excludeID || t.HP <= 0 {
			continue
		}
		dx := t.X - cx
		dy := t.Y - cy
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= radius {
			name := t.DefID
			if t.Def != nil {
				name = t.Def.Name
			}
			result = append(result, targetSnapshot{
				instanceID: t.InstanceID,
				defID:      t.DefID,
				name:       name,
				x:          t.X,
				y:          t.Y,
				hp:         t.HP,
				multiplier: t.Multiplier,
			})
		}
	}
	return result
}
