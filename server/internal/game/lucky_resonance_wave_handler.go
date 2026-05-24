// lucky_resonance_wave_handler.go — 幸運共鳴波魚系統（DAY-273）
// 業界依據：Royal Fishing / Jili 2026「連鎖閃電+群體攻擊」趨勢的進化版
//           業界原創「共鳴波擴散+全場同步爆發」機制，讓玩家有「一波帶一片，波波相連」的視覺爽感
//
// 設計：擊破 T231 後，觸發「共鳴波」：
//   - Server 以觸發點為中心，發出 3 層同心圓共鳴波（每層間隔 400ms）
//   - 第 1 層（r=150px）：波及目標 HP -20%，35% 機率「共鳴引爆」（×2.0 倍率，全服共享）
//   - 第 2 層（r=250px）：波及目標 HP -15%，25% 機率「共鳴引爆」（×1.8 倍率，全服共享）
//   - 第 3 層（r=350px）：波及目標 HP -10%，15% 機率「共鳴引爆」（×1.5 倍率，全服共享）
//   - 3 層波完成後，若引爆數 ≥ 5，觸發「共鳴爆發」：全服 ×1.5 加成 8 秒
//   - 個人冷卻 25 秒；全服冷卻 40 秒
//
// 設計差異：
//   - 與連鎖爆炸（T224，從一點向外擴散）不同，共鳴波是「同心圓擴散」，讓玩家看到「波紋從中心向外擴散」的視覺爽感
//   - 「3 層同心圓」讓玩家有「波波相連，越來越大」的期待感
//   - 「35%/25%/15% 機率引爆」讓每層波都有不確定性，製造「這層波會引爆幾個？」的期待感
//   - 「引爆數 ≥ 5 觸發全服爆發」讓玩家有「要趁波及範圍大時多引爆幾個」的策略感
//   - 「全服 ×1.5 加成 8 秒」讓所有玩家都受益，製造「全服一起爽」的社交感
//   - 「全服廣播每層波的引爆數和位置」讓所有玩家看到「波紋在哪裡、引爆了幾個」，製造視覺衝擊感
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

const (
	LuckyResonanceWavePersonalCD = 25 * time.Second // 個人冷卻
	LuckyResonanceWaveGlobalCD   = 40 * time.Second // 全服冷卻
	LuckyResonanceWaveInterval   = 400 * time.Millisecond // 每層波間隔

	// 三層波參數
	LuckyResonanceWaveLayer1Range  = 150.0 // 第 1 層範圍（px）
	LuckyResonanceWaveLayer2Range  = 250.0 // 第 2 層範圍（px）
	LuckyResonanceWaveLayer3Range  = 350.0 // 第 3 層範圍（px）
	LuckyResonanceWaveLayer1HP     = 0.20  // 第 1 層 HP 傷害比例
	LuckyResonanceWaveLayer2HP     = 0.15  // 第 2 層 HP 傷害比例
	LuckyResonanceWaveLayer3HP     = 0.10  // 第 3 層 HP 傷害比例
	LuckyResonanceWaveLayer1Prob   = 0.35  // 第 1 層引爆機率
	LuckyResonanceWaveLayer2Prob   = 0.25  // 第 2 層引爆機率
	LuckyResonanceWaveLayer3Prob   = 0.15  // 第 3 層引爆機率
	LuckyResonanceWaveLayer1Mult   = 2.0   // 第 1 層引爆倍率
	LuckyResonanceWaveLayer2Mult   = 1.8   // 第 2 層引爆倍率
	LuckyResonanceWaveLayer3Mult   = 1.5   // 第 3 層引爆倍率
	LuckyResonanceWaveBurstThresh  = 5     // 觸發全服爆發所需引爆數
	LuckyResonanceWaveBurstMult    = 1.5   // 全服爆發倍率
	LuckyResonanceWaveBurstDur     = 8 * time.Second // 全服爆發持續時間
)

// resonanceWaveBurst 全服共鳴爆發加成
type resonanceWaveBurst struct {
	mult      float64
	expiresAt time.Time
}

// luckyResonanceWaveManager 幸運共鳴波魚管理器
type luckyResonanceWaveManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 全服共鳴爆發加成（nil = 無加成）
	activeBurst *resonanceWaveBurst
}

func newLuckyResonanceWaveManager() *luckyResonanceWaveManager {
	return &luckyResonanceWaveManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyResonanceWaveFish 判斷是否為幸運共鳴波魚
func isLuckyResonanceWaveFish(defID string) bool {
	return defID == "T231"
}

// getResonanceWaveBurstMult 取得全服共鳴爆發倍率（供 handleKill 使用）
func (m *luckyResonanceWaveManager) getResonanceWaveBurstMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeBurst == nil || time.Now().After(m.activeBurst.expiresAt) {
		m.activeBurst = nil
		return 1.0
	}
	return m.activeBurst.mult
}

// tryLuckyResonanceWaveFish 擊破 T231 後觸發共鳴波
func (g *Game) tryLuckyResonanceWaveFish(p *player.Player, triggerX, triggerY float64) {
	m := g.LuckyResonanceWave

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
	m.personalCooldowns[p.ID] = now.Add(LuckyResonanceWavePersonalCD)
	m.globalCooldownUntil = now.Add(LuckyResonanceWaveGlobalCD)
	m.mu.Unlock()

	log.Printf("[ResonanceWave] player=%s 觸發共鳴波！位置=(%.0f,%.0f)", p.ID, triggerX, triggerY)

	// 個人通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyResonanceWave,
		Payload: ws.LuckyResonanceWavePayload{
			Event:      "wave_start",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			X:          triggerX,
			Y:          triggerY,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyResonanceWave,
		Payload: ws.LuckyResonanceWavePayload{
			Event:      "wave_broadcast",
			PlayerName: p.DisplayName,
			X:          triggerX,
			Y:          triggerY,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyResonanceWave, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌊 %s 觸發共鳴波！3 層同心圓擴散！",
			p.DisplayName),
		"color": "#00BFFF",
	})

	// 執行共鳴波（goroutine）
	go g.runResonanceWave(p, triggerX, triggerY)
}

// runResonanceWave 執行共鳴波主邏輯（3 層同心圓）
func (g *Game) runResonanceWave(p *player.Player, cx, cy float64) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalExploded := 0
	totalReward := 0

	// 三層波參數
	layers := []struct {
		radius   float64
		hpRatio  float64
		prob     float64
		mult     float64
		layerNum int
	}{
		{LuckyResonanceWaveLayer1Range, LuckyResonanceWaveLayer1HP, LuckyResonanceWaveLayer1Prob, LuckyResonanceWaveLayer1Mult, 1},
		{LuckyResonanceWaveLayer2Range, LuckyResonanceWaveLayer2HP, LuckyResonanceWaveLayer2Prob, LuckyResonanceWaveLayer2Mult, 2},
		{LuckyResonanceWaveLayer3Range, LuckyResonanceWaveLayer3HP, LuckyResonanceWaveLayer3Prob, LuckyResonanceWaveLayer3Mult, 3},
	}

	// 已被引爆的目標（避免重複引爆）
	explodedIDs := make(map[string]bool)

	for _, layer := range layers {
		// 找範圍內目標
		targets := g.findTargetsInRange(cx, cy, layer.radius, "")
		layerExploded := 0
		layerReward := 0

		for _, t := range targets {
			if explodedIDs[t.instanceID] || isLuckyResonanceWaveFish(t.defID) {
				continue
			}

			// HP 傷害
			g.applyExplosionDamage(t.instanceID, layer.hpRatio)

			// 機率引爆
			if rng.Float64() < layer.prob {
				reward := g.doExplosionKill(p, t.instanceID, t.multiplier, layer.mult)
				if reward > 0 {
					layerExploded++
					layerReward += reward
					totalExploded++
					totalReward += reward
					explodedIDs[t.instanceID] = true
				}
			}
		}

		log.Printf("[ResonanceWave] 第 %d 層：波及 %d 個目標，引爆 %d 個，獎勵 %d",
			layer.layerNum, len(targets), layerExploded, layerReward)

		// 廣播每層波結果
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyResonanceWave,
			Payload: ws.LuckyResonanceWavePayload{
				Event:         "wave_layer",
				PlayerName:    p.DisplayName,
				Layer:         layer.layerNum,
				X:             cx,
				Y:             cy,
				Radius:        layer.radius,
				AffectedCount: len(targets),
				ExplodeCount:  layerExploded,
				Mult:          layer.mult,
				Reward:        layerReward,
			},
		})

		// 等待下一層
		time.Sleep(LuckyResonanceWaveInterval)
	}

	log.Printf("[ResonanceWave] 共鳴波完成！總引爆 %d 個，總獎勵 %d", totalExploded, totalReward)

	// 判斷是否觸發全服爆發
	burstTriggered := totalExploded >= LuckyResonanceWaveBurstThresh

	if burstTriggered {
		// 設定全服爆發加成
		g.LuckyResonanceWave.mu.Lock()
		g.LuckyResonanceWave.activeBurst = &resonanceWaveBurst{
			mult:      LuckyResonanceWaveBurstMult,
			expiresAt: time.Now().Add(LuckyResonanceWaveBurstDur),
		}
		g.LuckyResonanceWave.mu.Unlock()

		log.Printf("[ResonanceWave] 共鳴爆發觸發！全服 ×%.1f 加成 %v",
			LuckyResonanceWaveBurstMult, LuckyResonanceWaveBurstDur)

		// 全服廣播爆發
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyResonanceWave,
			Payload: ws.LuckyResonanceWavePayload{
				Event:        "wave_burst",
				PlayerName:   p.DisplayName,
				TotalExplode: totalExploded,
				TotalReward:  totalReward,
				BurstMult:    LuckyResonanceWaveBurstMult,
				BurstDurSec:  int(LuckyResonanceWaveBurstDur.Seconds()),
			},
		})

		// 全服公告
		g.Announce.Create(announce.EventLuckyResonanceWave, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🌊 %s 共鳴爆發！%d 個目標引爆！全服 ×%.1f 加成 %d 秒！",
				p.DisplayName, totalExploded, LuckyResonanceWaveBurstMult,
				int(LuckyResonanceWaveBurstDur.Seconds())),
			"color": "#FFD700",
		})

		// 爆發結束後清除
		go func() {
			time.Sleep(LuckyResonanceWaveBurstDur)
			g.LuckyResonanceWave.mu.Lock()
			if g.LuckyResonanceWave.activeBurst != nil &&
				time.Now().After(g.LuckyResonanceWave.activeBurst.expiresAt) {
				g.LuckyResonanceWave.activeBurst = nil
			}
			g.LuckyResonanceWave.mu.Unlock()

			// 廣播爆發結束
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgLuckyResonanceWave,
				Payload: ws.LuckyResonanceWavePayload{
					Event: "wave_burst_end",
				},
			})
		}()
	} else {
		// 未達到爆發門檻，廣播結算
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyResonanceWave,
			Payload: ws.LuckyResonanceWavePayload{
				Event:        "wave_result",
				PlayerName:   p.DisplayName,
				TotalExplode: totalExploded,
				TotalReward:  totalReward,
				BurstMult:    1.0,
			},
		})
	}
}
