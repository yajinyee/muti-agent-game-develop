// lucky_cosmic_pulse_handler.go — 幸運宇宙脈衝魚系統（DAY-287）
// 業界依據：TaDa Gaming 2026「Cosmic」主題 + Fishing Fortune 2026「pulse wave mechanics」
//          業界原創「宇宙脈衝波+全場共振+脈衝連鎖爆發」機制
//
// 設計：
//   - 擊破 T245 後，發出「宇宙脈衝波」（從場地中心向外擴散）
//   - 脈衝波分 3 層（每 800ms 一層），每層命中範圍內所有目標 HP -20%
//   - 每層脈衝波命中的目標數 × 0.2 = 累積倍率加成（最高 ×5.0）
//   - 若 3 層脈衝波命中目標總數 ≥ 15 → 「宇宙共振」：全服 ×2.2 加成 6 秒
//   - 全服廣播脈衝波擴散和命中結果
//   - 個人冷卻 24 秒；全服冷卻 40 秒
//
// 設計差異：
//   - 與克拉肯（T244，8 條觸手精準攻擊）不同，宇宙脈衝是「同心圓擴散波」
//   - 「每層命中目標數 × 0.2 累積倍率」讓玩家有「場上魚越多，脈衝越值錢」的策略感
//   - 「3 層脈衝波命中總數 ≥ 15 觸發宇宙共振」讓玩家有「要趁魚多的時候觸發」的時機感
//   - 「全服 ×2.2 加成 6 秒」讓所有玩家都受益，製造「全服一起爽」的社交感
//   - 「HP -20% 弱化」比克拉肯（-35%）更溫和，但 3 層疊加讓魚更容易打
package game

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyCosmicPulsePersonalCD = 24 * time.Second // 個人冷卻
	LuckyCosmicPulseGlobalCD   = 40 * time.Second // 全服冷卻

	// 宇宙脈衝波設計
	CosmicPulseWaveCount    = 3                      // 脈衝波層數
	CosmicPulseInterval     = 800 * time.Millisecond // 每層間隔
	CosmicPulseHPDmg        = 0.20                   // 每層 HP -20%
	CosmicPulseMultPerHit   = 0.2                    // 每命中 1 個目標 +0.2 倍率
	CosmicPulseMaxMult      = 5.0                    // 最高累積倍率
	CosmicPulseResonanceMin = 15                     // 觸發宇宙共振的最低命中總數

	// 宇宙共振：全服加成
	CosmicResonanceMult     = 2.2                   // 全服 ×2.2
	CosmicResonanceDuration = 6 * time.Second       // 持續 6 秒

	// 脈衝波半徑（每層擴大）
	CosmicPulseRadiusBase = 200.0 // 第 1 層半徑（px）
	CosmicPulseRadiusStep = 150.0 // 每層增加半徑（px）
)

// cosmicResonanceBoost 宇宙共振全服加成
type cosmicResonanceBoost struct {
	mult      float64
	expiresAt time.Time
}

// luckyCosmicPulseManager 幸運宇宙脈衝魚管理器
type luckyCosmicPulseManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 宇宙共振全服加成
	resonanceBoost *cosmicResonanceBoost
}

func newLuckyCosmicPulseManager() *luckyCosmicPulseManager {
	return &luckyCosmicPulseManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyCosmicPulseFish 判斷是否為幸運宇宙脈衝魚
func isLuckyCosmicPulseFish(defID string) bool {
	return defID == "T245"
}

// getCosmicResonanceMult 取得宇宙共振全服加成倍率（供 handleKill 使用）
func (m *luckyCosmicPulseManager) getCosmicResonanceMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.resonanceBoost != nil && time.Now().Before(m.resonanceBoost.expiresAt) {
		return m.resonanceBoost.mult
	}
	return 1.0
}

// tryLuckyCosmicPulseFish 擊破 T245 後觸發宇宙脈衝（供 handleKill 使用）
func (g *Game) tryLuckyCosmicPulseFish(p *player.Player) {
	mgr := g.LuckyCosmicPulse
	mgr.mu.Lock()

	now := time.Now()

	// 全服冷卻檢查
	if now.Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && now.Before(cd) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = now.Add(LuckyCosmicPulsePersonalCD)
	mgr.globalCooldownUntil = now.Add(LuckyCosmicPulseGlobalCD)
	mgr.mu.Unlock()

	log.Printf("[CosmicPulse] player=%s waves=%d", p.ID, CosmicPulseWaveCount)

	// 全服廣播：宇宙脈衝開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyCosmicPulse,
		Payload: ws.LuckyCosmicPulsePayload{
			Event:      "pulse_start",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			WaveCount:  CosmicPulseWaveCount,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyCosmicPulse, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌌✨ %s 觸發宇宙脈衝！%d 層脈衝波即將擴散！",
			p.DisplayName, CosmicPulseWaveCount),
		"color": "#1A0A3E",
	})
	g.broadcastAnnouncement(ann)

	// 啟動脈衝波序列 goroutine
	go g.runCosmicPulseWaves(p)
}

// runCosmicPulseWaves 執行宇宙脈衝波序列
func (g *Game) runCosmicPulseWaves(p *player.Player) {
	accumMult := 1.0
	totalHit := 0 // 3 層脈衝波命中目標總數

	// 場地中心
	centerX := 512.0
	centerY := 300.0

	for wave := 0; wave < CosmicPulseWaveCount; wave++ {
		time.Sleep(CosmicPulseInterval)

		// 本層脈衝波半徑（逐層擴大）
		radius := CosmicPulseRadiusBase + float64(wave)*CosmicPulseRadiusStep

		// 對範圍內目標造成 HP -20%
		hitTargets := g.applyCosmicPulseDamage(centerX, centerY, radius, CosmicPulseHPDmg)
		totalHit += hitTargets

		// 每命中 1 個目標 +0.2 倍率（最高 ×5.0）
		if hitTargets > 0 {
			newMult := accumMult + float64(hitTargets)*CosmicPulseMultPerHit
			if newMult > CosmicPulseMaxMult {
				newMult = CosmicPulseMaxMult
			}
			accumMult = newMult
		}

		// 廣播脈衝波結果
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyCosmicPulse,
			Payload: ws.LuckyCosmicPulsePayload{
				Event:      "pulse_wave",
				PlayerID:   p.ID,
				WaveIdx:    wave + 1,
				TotalWaves: CosmicPulseWaveCount,
				Radius:     radius,
				HitTargets: hitTargets,
				AccumMult:  accumMult,
			},
		})

		log.Printf("[CosmicPulse] wave %d/%d radius=%.0f hit=%d accumMult=%.2f",
			wave+1, CosmicPulseWaveCount, radius, hitTargets, accumMult)
	}

	// 結算
	g.settleCosmicPulse(p, totalHit, accumMult)
}

// applyCosmicPulseDamage 對範圍內所有目標造成 HP 傷害，回傳命中目標數
func (g *Game) applyCosmicPulseDamage(cx, cy, radius, dmgPct float64) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	hitCount := 0
	for _, t := range g.Targets {
		if t == nil || t.HP <= 0 {
			continue
		}
		// 計算距離
		dx := t.X - cx
		dy := t.Y - cy
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= radius {
			// HP -20%
			dmg := int(float64(t.HP) * dmgPct)
			if dmg < 1 {
				dmg = 1
			}
			t.HP -= dmg
			if t.HP < 0 {
				t.HP = 0
			}
			hitCount++

			// 廣播目標 HP 更新
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgTargetUpdate,
				Payload: ws.TargetUpdatePayload{
					InstanceID: t.InstanceID,
					HP:         t.HP,
					MaxHP:      t.MaxHP,
				},
			})
		}
	}
	return hitCount
}

// settleCosmicPulse 宇宙脈衝結算
func (g *Game) settleCosmicPulse(p *player.Player, totalHit int, accumMult float64) {
	// 計算個人獎勵
	reward := 0
	if accumMult > 1.0 {
		betDef := p.GetBetDef()
		bet := 0
		if betDef != nil {
			bet = betDef.BetCost
		}
		reward = int(float64(bet) * accumMult)
		p.AddCoins(reward)
	}

	// 判斷是否觸發宇宙共振（命中總數 ≥ 15）
	isResonance := totalHit >= CosmicPulseResonanceMin

	log.Printf("[CosmicPulse] settle player=%s totalHit=%d accumMult=%.2f reward=%d resonance=%v",
		p.ID, totalHit, accumMult, reward, isResonance)

	// 廣播結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyCosmicPulse,
		Payload: ws.LuckyCosmicPulsePayload{
			Event:       "pulse_end",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			TotalHit:    totalHit,
			AccumMult:   accumMult,
			Reward:      reward,
			IsResonance: isResonance,
		},
	})

	// 宇宙共振：全服 ×2.2 加成 6 秒
	if isResonance {
		g.doCosmicResonance(p, totalHit)
	} else {
		// 普通結算公告（累積倍率 ≥ 3.0 才公告）
		if accumMult >= 3.0 {
			ann := g.Announce.Create(announce.EventLuckyCosmicPulse, p.DisplayName, 0, map[string]string{
				"message": fmt.Sprintf("🌌 %s 宇宙脈衝結算！命中 %d 個目標，累積 ×%.1f！獲得 %d 金幣！",
					p.DisplayName, totalHit, accumMult, reward),
				"color": "#4B0082",
			})
			g.broadcastAnnouncement(ann)
		}
	}
}

// doCosmicResonance 宇宙共振：全服 ×2.2 加成 6 秒
func (g *Game) doCosmicResonance(p *player.Player, totalHit int) {
	mgr := g.LuckyCosmicPulse
	mgr.mu.Lock()
	mgr.resonanceBoost = &cosmicResonanceBoost{
		mult:      CosmicResonanceMult,
		expiresAt: time.Now().Add(CosmicResonanceDuration),
	}
	mgr.mu.Unlock()

	log.Printf("[CosmicPulse] RESONANCE! player=%s totalHit=%d global x%.1f for %v",
		p.ID, totalHit, CosmicResonanceMult, CosmicResonanceDuration)

	// 全服廣播宇宙共振
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyCosmicPulse,
		Payload: ws.LuckyCosmicPulsePayload{
			Event:      "pulse_resonance",
			PlayerName: p.DisplayName,
			TotalHit:   totalHit,
			ResMult:    CosmicResonanceMult,
			Duration:   int(CosmicResonanceDuration.Seconds()),
		},
	})

	// 全服最高優先公告
	ann := g.Announce.Create(announce.EventLuckyCosmicPulse, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌌🌟🌌 %s 宇宙共振！命中 %d 個目標！全服 ×%.1f 加成 %d 秒！",
			p.DisplayName, totalHit, CosmicResonanceMult, int(CosmicResonanceDuration.Seconds())),
		"color": "#0D0025",
	})
	g.broadcastAnnouncement(ann)

	// 6 秒後廣播共振結束
	go func() {
		time.Sleep(CosmicResonanceDuration)
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyCosmicPulse,
			Payload: ws.LuckyCosmicPulsePayload{
				Event: "pulse_resonance_end",
			},
		})
	}()
}

// end of lucky_cosmic_pulse_handler.go
