// lucky_kraken_handler.go — 幸運深海克拉肯魚系統（DAY-286）
// 業界依據：Kraken Unleashed「Kraken Reel + 多段觸手攻擊 + Grand Jackpot」機制
//          Fishing Fortune 2026「multiplier chains + jackpot triggers」機制
// 業界原創「克拉肯觸手連擊+累積倍率+狂怒爆發」機制
//
// 設計：
//   - 擊破 T244 後，召喚「深海克拉肯」：
//     1. 8 條觸手依序攻擊（每 500ms 一條）
//     2. 每條觸手命中 1-3 個目標（HP -35%）
//     3. 每條觸手命中至少 1 個目標 → 觸發玩家獲得 ×1.3 累積倍率（最高 ×10.0）
//     4. 若 8 條觸手全部命中（無空揮）→「克拉肯狂怒」：全服 ×2.8 加成 7 秒
//   - 全服廣播觸手攻擊位置和結果
//   - 個人冷卻 30 秒；全服冷卻 48 秒
//
// 設計差異：
//   - 與龍怒隕石（T242，4-7 顆隕石 AOE）不同，克拉肯是「8 條觸手精準攻擊」
//   - 「每條觸手命中 1-3 個目標」讓每次攻擊都有「這條觸手打到幾條魚？」的期待感
//   - 「8 條觸手全部命中觸發狂怒」比龍怒完美（所有隕石命中）更難，但獎勵更高
//   - 「最高 ×10.0 累積倍率」是所有幸運魚中最高的個人倍率
//   - 「全服 ×2.8 加成 7 秒」讓所有玩家都受益，製造「全服一起爽」的社交感
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
	LuckyKrakenPersonalCD = 30 * time.Second // 個人冷卻
	LuckyKrakenGlobalCD   = 48 * time.Second // 全服冷卻

	// 克拉肯觸手設計
	KrakenTentacleCount    = 8                      // 觸手數量
	KrakenTentacleInterval = 500 * time.Millisecond // 每條觸手間隔
	KrakenTentacleHPDmg    = 0.35                   // 每條觸手 HP -35%
	KrakenTentacleRadius   = 100.0                  // 觸手命中半徑（px）
	KrakenAccumMult        = 1.3                    // 每條命中累積倍率
	KrakenMaxMult          = 10.0                   // 最高累積倍率

	// 克拉肯狂怒：全服加成
	KrakenFuryMult     = 2.8                   // 全服 ×2.8
	KrakenFuryDuration = 7 * time.Second       // 持續 7 秒
)

// krakenFuryBoost 克拉肯狂怒全服加成
type krakenFuryBoost struct {
	mult      float64
	expiresAt time.Time
}

// luckyKrakenManager 幸運深海克拉肯魚管理器
type luckyKrakenManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 克拉肯狂怒全服加成
	furyBoost *krakenFuryBoost
}

func newLuckyKrakenManager() *luckyKrakenManager {
	return &luckyKrakenManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyKrakenFish 判斷是否為幸運深海克拉肯魚
func isLuckyKrakenFish(defID string) bool {
	return defID == "T244"
}

// getKrakenFuryMult 取得克拉肯狂怒全服加成倍率（供 handleKill 使用）
func (m *luckyKrakenManager) getKrakenFuryMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.furyBoost != nil && time.Now().Before(m.furyBoost.expiresAt) {
		return m.furyBoost.mult
	}
	return 1.0
}

// tryLuckyKrakenFish 擊破 T244 後觸發克拉肯（供 handleKill 使用）
func (g *Game) tryLuckyKrakenFish(p *player.Player) {
	mgr := g.LuckyKraken
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
	mgr.personalCooldowns[p.ID] = now.Add(LuckyKrakenPersonalCD)
	mgr.globalCooldownUntil = now.Add(LuckyKrakenGlobalCD)
	mgr.mu.Unlock()

	log.Printf("[Kraken] player=%s tentacles=%d", p.ID, KrakenTentacleCount)

	// 全服廣播：克拉肯召喚開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyKraken,
		Payload: ws.LuckyKrakenPayload{
			Event:          "kraken_start",
			PlayerID:       p.ID,
			PlayerName:     p.DisplayName,
			TentacleCount:  KrakenTentacleCount,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyKraken, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🦑🌊 %s 召喚深海克拉肯！%d 條觸手即將攻擊！",
			p.DisplayName, KrakenTentacleCount),
		"color": "#1E3A5F",
	})
	g.broadcastAnnouncement(ann)

	// 啟動觸手攻擊序列 goroutine
	go g.runKrakenTentacles(p)
}

// runKrakenTentacles 執行克拉肯觸手攻擊序列
func (g *Game) runKrakenTentacles(p *player.Player) {
	accumMult := 1.0
	hitCount := 0   // 命中目標的觸手數
	totalHit := 0   // 命中的目標總數

	for i := 0; i < KrakenTentacleCount; i++ {
		time.Sleep(KrakenTentacleInterval)

		// 觸手攻擊位置（隨機，偏向場地中央）
		tentacleX := 120.0 + rand.Float64()*760.0  // 120-880px
		tentacleY := 80.0 + rand.Float64()*440.0   // 80-520px

		// 對範圍內目標造成 HP -35%（最多命中 3 個）
		hitTargets := g.applyKrakenTentacleDamage(tentacleX, tentacleY, KrakenTentacleRadius, KrakenTentacleHPDmg, 3)
		totalHit += hitTargets

		// 命中至少 1 個目標 → 累積倍率
		if hitTargets > 0 {
			hitCount++
			newMult := accumMult * KrakenAccumMult
			if newMult > KrakenMaxMult {
				newMult = KrakenMaxMult
			}
			accumMult = newMult
		}

		// 廣播觸手攻擊結果
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyKraken,
			Payload: ws.LuckyKrakenPayload{
				Event:         "kraken_tentacle",
				PlayerID:      p.ID,
				TentacleX:     tentacleX,
				TentacleY:     tentacleY,
				HitTargets:    hitTargets,
				AccumMult:     accumMult,
				TentacleIdx:   i + 1,
				TotalCount:    KrakenTentacleCount,
			},
		})

		log.Printf("[Kraken] tentacle %d/%d pos=(%.0f,%.0f) hit=%d accumMult=%.2f",
			i+1, KrakenTentacleCount, tentacleX, tentacleY, hitTargets, accumMult)
	}

	// 結算
	g.settleKraken(p, hitCount, totalHit, accumMult)
}

// applyKrakenTentacleDamage 對範圍內目標造成 HP 傷害（最多 maxHit 個），回傳命中目標數
func (g *Game) applyKrakenTentacleDamage(cx, cy, radius, dmgPct float64, maxHit int) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	hitCount := 0
	for _, t := range g.Targets {
		if t == nil || t.HP <= 0 {
			continue
		}
		if hitCount >= maxHit {
			break
		}
		// 計算距離
		dx := t.X - cx
		dy := t.Y - cy
		dist := dx*dx + dy*dy
		if dist <= radius*radius {
			// HP -35%
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

// settleKraken 克拉肯結算
func (g *Game) settleKraken(p *player.Player, hitCount, totalHit int, accumMult float64) {
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

	// 判斷是否克拉肯狂怒（所有觸手都命中）
	isFury := hitCount == KrakenTentacleCount

	log.Printf("[Kraken] settle player=%s hitCount=%d/%d totalHit=%d accumMult=%.2f reward=%d fury=%v",
		p.ID, hitCount, KrakenTentacleCount, totalHit, accumMult, reward, isFury)

	// 廣播結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyKraken,
		Payload: ws.LuckyKrakenPayload{
			Event:         "kraken_end",
			PlayerID:      p.ID,
			PlayerName:    p.DisplayName,
			TentacleCount: KrakenTentacleCount,
			HitCount:      hitCount,
			TotalHit:      totalHit,
			AccumMult:     accumMult,
			Reward:        reward,
			IsFury:        isFury,
		},
	})

	// 克拉肯狂怒：全服 ×2.8 加成 7 秒
	if isFury {
		g.doKrakenFury(p)
	} else {
		// 普通結算公告（累積倍率 ≥ 4.0 才公告）
		if accumMult >= 4.0 {
			ann := g.Announce.Create(announce.EventLuckyKraken, p.DisplayName, 0, map[string]string{
				"message": fmt.Sprintf("🦑 %s 克拉肯結算！命中 %d 個目標，累積 ×%.1f！獲得 %d 金幣！",
					p.DisplayName, totalHit, accumMult, reward),
				"color": "#4A90D9",
			})
			g.broadcastAnnouncement(ann)
		}
	}
}

// doKrakenFury 克拉肯狂怒：全服 ×2.8 加成 7 秒
func (g *Game) doKrakenFury(p *player.Player) {
	mgr := g.LuckyKraken
	mgr.mu.Lock()
	mgr.furyBoost = &krakenFuryBoost{
		mult:      KrakenFuryMult,
		expiresAt: time.Now().Add(KrakenFuryDuration),
	}
	mgr.mu.Unlock()

	log.Printf("[Kraken] FURY! player=%s global x%.1f for %v",
		p.ID, KrakenFuryMult, KrakenFuryDuration)

	// 全服廣播克拉肯狂怒
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyKraken,
		Payload: ws.LuckyKrakenPayload{
			Event:      "kraken_fury",
			PlayerName: p.DisplayName,
			FuryMult:   KrakenFuryMult,
			Duration:   int(KrakenFuryDuration.Seconds()),
		},
	})

	// 全服最高優先公告
	ann := g.Announce.Create(announce.EventLuckyKraken, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🦑🌊🦑 %s 克拉肯狂怒！全服 ×%.1f 加成 %d 秒！",
			p.DisplayName, KrakenFuryMult, int(KrakenFuryDuration.Seconds())),
		"color": "#0A1628",
	})
	g.broadcastAnnouncement(ann)

	// 7 秒後廣播狂怒結束
	go func() {
		time.Sleep(KrakenFuryDuration)
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyKraken,
			Payload: ws.LuckyKrakenPayload{
				Event: "kraken_fury_end",
			},
		})
	}()
}

// end of lucky_kraken_handler.go
