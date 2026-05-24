// lucky_dragon_wrath_handler.go — 幸運龍怒隕石魚系統（DAY-284）
// 業界依據：Royal Fishing Jili「Dragon Wrath meteors」機制（2026 最熱門）
// 龍怒召喚隕石雨，全場 AOE 傷害 + 累積倍率 + 完美爆發
//
// 設計：
//   - 擊破 T242 後，召喚 4-7 顆「龍怒隕石」
//   - 每顆隕石在 4 秒內依序墜落（每 600ms 一顆）
//   - 每顆隕石墜落時，命中範圍內（r=120px）所有目標 HP -45%
//   - 每顆隕石命中至少 1 個目標 → 觸發玩家獲得 ×1.4 累積倍率（最高 ×8.0）
//   - 若所有隕石都命中目標（無空砸）→「龍怒完美」：全服 ×2.5 加成 6 秒
//   - 全服廣播隕石墜落位置和結果
//   - 個人冷卻 26 秒；全服冷卻 42 秒
//
// 設計差異：
//   - 與黃金颶風（T234，螺旋掃場 HP-30%）不同，龍怒隕石是「多點 AOE 墜落」
//   - 「每顆隕石命中才累積倍率」讓玩家有「要是每顆都命中就賺大了」的期待感
//   - 「龍怒完美（全部命中）→ 全服 ×2.5」是最高全服加成，製造「全服一起爽」的社交感
//   - 「HP -45%」比黃金颶風（-30%）更強，讓玩家感受到「龍怒的威力」
//   - 「全服廣播隕石位置」讓所有玩家看到「隕石在哪裡砸」，製造全服緊張感
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
	LuckyDragonWrathPersonalCD = 26 * time.Second // 個人冷卻
	LuckyDragonWrathGlobalCD   = 42 * time.Second // 全服冷卻

	// 隕石設計
	DragonWrathMeteorInterval = 600 * time.Millisecond // 每顆隕石間隔
	DragonWrathMeteorHPDmg    = 0.45                   // 每顆隕石 HP -45%
	DragonWrathMeteorRadius   = 120.0                  // 隕石命中半徑（px）
	DragonWrathAccumMult      = 1.4                    // 每顆命中累積倍率
	DragonWrathMaxMult        = 8.0                    // 最高累積倍率

	// 龍怒完美：全服加成
	DragonWrathPerfectMult     = 2.5                   // 全服 ×2.5
	DragonWrathPerfectDuration = 6 * time.Second       // 持續 6 秒
)

// dragonWrathPerfectBoost 龍怒完美全服加成
type dragonWrathPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

// luckyDragonWrathManager 幸運龍怒隕石魚管理器
type luckyDragonWrathManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 龍怒完美全服加成
	perfectBoost *dragonWrathPerfectBoost
}

func newLuckyDragonWrathManager() *luckyDragonWrathManager {
	return &luckyDragonWrathManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyDragonWrathFish 判斷是否為幸運龍怒隕石魚
func isLuckyDragonWrathFish(defID string) bool {
	return defID == "T242"
}

// getDragonWrathPerfectMult 取得龍怒完美全服加成倍率（供 handleKill 使用）
func (m *luckyDragonWrathManager) getDragonWrathPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

// tryLuckyDragonWrathFish 擊破 T242 後觸發龍怒隕石（供 handleKill 使用）
func (g *Game) tryLuckyDragonWrathFish(p *player.Player) {
	mgr := g.LuckyDragonWrath
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
	mgr.personalCooldowns[p.ID] = now.Add(LuckyDragonWrathPersonalCD)
	mgr.globalCooldownUntil = now.Add(LuckyDragonWrathGlobalCD)
	mgr.mu.Unlock()

	// 決定隕石數量（4-7 顆）
	meteorCount := 4 + rand.Intn(4)

	log.Printf("[DragonWrath] player=%s meteors=%d", p.ID, meteorCount)

	// 全服廣播：龍怒隕石開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyDragonWrath,
		Payload: ws.LuckyDragonWrathPayload{
			Event:       "wrath_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			MeteorCount: meteorCount,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyDragonWrath, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🐉🔥 %s 召喚龍怒隕石！%d 顆隕石即將墜落！",
			p.DisplayName, meteorCount),
		"color": "#FF4500",
	})
	g.broadcastAnnouncement(ann)

	// 啟動隕石序列 goroutine
	go g.runLuckyDragonWrathMeteors(p, meteorCount)
}

// runLuckyDragonWrathMeteors 執行龍怒隕石序列
func (g *Game) runLuckyDragonWrathMeteors(p *player.Player, meteorCount int) {
	accumMult := 1.0
	hitCount := 0   // 命中目標的隕石數
	totalHit := 0   // 命中的目標總數

	for i := 0; i < meteorCount; i++ {
		time.Sleep(DragonWrathMeteorInterval)

		// 隕石墜落位置（隨機）
		meteorX := 80.0 + rand.Float64()*840.0  // 80-920px
		meteorY := 100.0 + rand.Float64()*400.0 // 100-500px

		// 對範圍內目標造成 HP -45%
		hitTargets := g.applyDragonWrathMeteorDamage(meteorX, meteorY, DragonWrathMeteorRadius, DragonWrathMeteorHPDmg)
		hitCount++
		totalHit += hitTargets

		// 命中至少 1 個目標 → 累積倍率
		if hitTargets > 0 {
			newMult := accumMult * DragonWrathAccumMult
			if newMult > DragonWrathMaxMult {
				newMult = DragonWrathMaxMult
			}
			accumMult = newMult
		}

		// 廣播隕石墜落結果
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyDragonWrath,
			Payload: ws.LuckyDragonWrathPayload{
				Event:      "wrath_meteor",
				PlayerID:   p.ID,
				MeteorX:    meteorX,
				MeteorY:    meteorY,
				HitTargets: hitTargets,
				AccumMult:  accumMult,
				MeteorIdx:  i + 1,
				TotalCount: meteorCount,
			},
		})

		log.Printf("[DragonWrath] meteor %d/%d pos=(%.0f,%.0f) hit=%d accumMult=%.2f",
			i+1, meteorCount, meteorX, meteorY, hitTargets, accumMult)
	}

	// 結算
	g.settleDragonWrath(p, meteorCount, hitCount, totalHit, accumMult)
}

// applyDragonWrathMeteorDamage 對範圍內目標造成 HP 傷害，回傳命中目標數
func (g *Game) applyDragonWrathMeteorDamage(cx, cy, radius, dmgPct float64) int {
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
		dist := dx*dx + dy*dy
		if dist <= radius*radius {
			// HP -45%
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

// settleDragonWrath 龍怒隕石結算
func (g *Game) settleDragonWrath(p *player.Player, meteorCount, hitCount, totalHit int, accumMult float64) {
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

	// 判斷是否龍怒完美（所有隕石都命中）
	isPerfect := hitCount == meteorCount && totalHit >= meteorCount

	log.Printf("[DragonWrath] settle player=%s meteors=%d hitCount=%d totalHit=%d accumMult=%.2f reward=%d perfect=%v",
		p.ID, meteorCount, hitCount, totalHit, accumMult, reward, isPerfect)

	// 廣播結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyDragonWrath,
		Payload: ws.LuckyDragonWrathPayload{
			Event:      "wrath_end",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			MeteorCount: meteorCount,
			TotalHit:   totalHit,
			AccumMult:  accumMult,
			Reward:     reward,
			IsPerfect:  isPerfect,
		},
	})

	// 龍怒完美：全服 ×2.5 加成 6 秒
	if isPerfect {
		g.doDragonWrathPerfect(p)
	} else {
		// 普通結算公告
		if accumMult >= 3.0 {
			ann := g.Announce.Create(announce.EventLuckyDragonWrath, p.DisplayName, 0, map[string]string{
				"message": fmt.Sprintf("🔥 %s 龍怒隕石結算！命中 %d 個目標，累積 ×%.1f！獲得 %d 金幣！",
					p.DisplayName, totalHit, accumMult, reward),
				"color": "#FF6B35",
			})
			g.broadcastAnnouncement(ann)
		}
	}
}

// doDragonWrathPerfect 龍怒完美：全服 ×2.5 加成 6 秒
func (g *Game) doDragonWrathPerfect(p *player.Player) {
	mgr := g.LuckyDragonWrath
	mgr.mu.Lock()
	mgr.perfectBoost = &dragonWrathPerfectBoost{
		mult:      DragonWrathPerfectMult,
		expiresAt: time.Now().Add(DragonWrathPerfectDuration),
	}
	mgr.mu.Unlock()

	log.Printf("[DragonWrath] PERFECT! player=%s global x%.1f for %v",
		p.ID, DragonWrathPerfectMult, DragonWrathPerfectDuration)

	// 全服廣播龍怒完美
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyDragonWrath,
		Payload: ws.LuckyDragonWrathPayload{
			Event:      "wrath_perfect",
			PlayerName: p.DisplayName,
			PerfectMult: DragonWrathPerfectMult,
			Duration:   int(DragonWrathPerfectDuration.Seconds()),
		},
	})

	// 全服最高優先公告
	ann := g.Announce.Create(announce.EventLuckyDragonWrath, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🐉🔥🐉 %s 龍怒完美！全服 ×%.1f 加成 %d 秒！",
			p.DisplayName, DragonWrathPerfectMult, int(DragonWrathPerfectDuration.Seconds())),
		"color": "#FF0000",
	})
	g.broadcastAnnouncement(ann)

	// 6 秒後廣播完美結束
	go func() {
		time.Sleep(DragonWrathPerfectDuration)
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyDragonWrath,
			Payload: ws.LuckyDragonWrathPayload{
				Event: "wrath_perfect_end",
			},
		})
	}()
}

// end of lucky_dragon_wrath_handler.go
