// lucky_lightning_hammer_handler.go — 幸運閃電錘魚系統（DAY-277）
// 業界依據：Battle of Luck「Lucky Slammer」機制（2026）進化版
//
// 設計：擊破 T235 後，觸發「閃電錘」：
//   - 瞬間選定場上 3-6 個目標（隨機）
//   - 對每個目標造成「閃電錘擊」（HP -60%，30% 機率直接擊破）
//   - 每個被錘擊的目標，觸發玩家獲得 ×1.2 倍率加成（累積）
//   - 被直接擊破的目標，額外給予 ×2.0 倍率獎勵（個人）
//   - 全服廣播錘擊結果（錘擊數/擊破數/總倍率）
//   - 個人冷卻 22 秒；全服冷卻 35 秒
//
// 設計差異：
//   - 與黃金颶風（T234，螺旋掃場 HP-30%）不同，閃電錘是「瞬間多目標錘擊」，
//     讓玩家有「一錘打多條魚，有些直接死掉」的爽感
//   - 「30% 機率直接擊破」讓每次錘擊都有「這條會不會直接死？」的期待感
//   - 「累積倍率 ×1.2/次」讓錘擊數越多倍率越高
//   - 「直接擊破額外 ×2.0」讓玩家有「要是全部都直接死就賺大了」的動力
//   - 「全服廣播錘擊結果」讓所有玩家看到「有幾條被直接錘死」，製造社交話題感
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
	LuckyLightningHammerPersonalCD  = 22 * time.Second // 個人冷卻
	LuckyLightningHammerGlobalCD    = 35 * time.Second // 全服冷卻
	LuckyLightningHammerHPDamage    = 0.60             // 錘擊 HP 傷害比例（-60%）
	LuckyLightningHammerKillChance  = 0.30             // 直接擊破機率（30%）
	LuckyLightningHammerMultPerHit  = 1.2              // 每次錘擊的倍率加成（累積乘法）
	LuckyLightningHammerKillBonus   = 2.0              // 直接擊破的額外倍率獎勵
	LuckyLightningHammerMinTargets  = 3                // 最少錘擊目標數
	LuckyLightningHammerMaxTargets  = 6                // 最多錘擊目標數
	LuckyLightningHammerHitDelay    = 200 * time.Millisecond // 每次錘擊間隔（視覺效果）
)

// luckyLightningHammerManager 幸運閃電錘魚管理器
type luckyLightningHammerManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time
}

func newLuckyLightningHammerManager() *luckyLightningHammerManager {
	return &luckyLightningHammerManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyLightningHammerFish 判斷是否為幸運閃電錘魚
func isLuckyLightningHammerFish(defID string) bool {
	return defID == "T235"
}

// tryLuckyLightningHammerFish 擊破 T235 後觸發閃電錘（供 handleKill 使用）
func (g *Game) tryLuckyLightningHammerFish(p *player.Player) {
	mgr := g.LuckyLightningHammer
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}
	// 全服冷卻檢查
	if time.Now().Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyLightningHammerPersonalCD)
	mgr.globalCooldownUntil = time.Now().Add(LuckyLightningHammerGlobalCD)
	mgr.mu.Unlock()

	log.Printf("[LuckyLightningHammer] player=%s triggered lightning hammer", p.ID)

	// 選定錘擊目標（3-6 個）
	g.mu.RLock()
	var candidates []string
	for iid, t := range g.Targets {
		if t.HP > 0 && !isLuckyLightningHammerFish(t.Def.ID) {
			candidates = append(candidates, iid)
		}
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		log.Printf("[LuckyLightningHammer] no targets available")
		return
	}

	// 隨機選 3-6 個目標（不重複）
	hammerCount := LuckyLightningHammerMinTargets + rand.Intn(LuckyLightningHammerMaxTargets-LuckyLightningHammerMinTargets+1)
	if hammerCount > len(candidates) {
		hammerCount = len(candidates)
	}

	// Fisher-Yates shuffle 選取
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	selectedTargets := candidates[:hammerCount]

	// 全服廣播：閃電錘觸發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyLightningHammer,
		Payload: ws.LuckyLightningHammerPayload{
			Event:       "hammer_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			HammerCount: hammerCount,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyLightningHammer, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⚡ %s 觸發閃電錘！瞬間錘擊 %d 個目標！30%% 機率直接擊破！",
			p.DisplayName, hammerCount),
		"color": "#FFD700",
	})
	g.broadcastAnnouncement(ann)

	// 執行錘擊（goroutine，帶視覺延遲）
	go g.runLightningHammerHits(p, selectedTargets)
}

// runLightningHammerHits 執行閃電錘擊（goroutine）
func (g *Game) runLightningHammerHits(p *player.Player, targetIIDs []string) {
	accumMult := 1.0
	hitCount := 0
	killCount := 0
	totalReward := 0

	for _, iid := range targetIIDs {
		time.Sleep(LuckyLightningHammerHitDelay)

		g.mu.Lock()
		t, ok := g.Targets[iid]
		if !ok || t.HP <= 0 {
			g.mu.Unlock()
			continue
		}

		// 計算 HP 傷害（-60%，最少 1）
		damage := int(float64(t.HP) * LuckyLightningHammerHPDamage)
		if damage < 1 {
			damage = 1
		}

		// 判斷是否直接擊破（30% 機率）
		killed := rand.Float64() < LuckyLightningHammerKillChance
		defID := t.Def.ID

		if killed {
			// 直接擊破：HP 歸零
			t.HP = 0
			killCount++
		} else {
			// 只造成傷害
			t.HP -= damage
			if t.HP < 1 {
				t.HP = 1
			}
		}
		g.mu.Unlock()

		hitCount++

		// 累積倍率（每次錘擊 ×1.2）
		accumMult *= LuckyLightningHammerMultPerHit

		// 直接擊破額外獎勵
		if killed {
			// 計算擊破獎勵（簡化：用 bet × 倍率 × killBonus）
			g.mu.RLock()
			var betCost int
			if pp, exists := g.Players[p.ID]; exists {
				betCost = pp.GetBetDef().BetCost
			}
			g.mu.RUnlock()

			killReward := int(float64(betCost) * LuckyLightningHammerKillBonus)
			totalReward += killReward

			log.Printf("[LuckyLightningHammer] target=%s KILLED! killReward=%d accumMult=%.2f",
				iid, killReward, accumMult)
		}

		log.Printf("[LuckyLightningHammer] hit #%d: target=%s killed=%v hp_damage=%d accum_mult=%.2f",
			hitCount, iid, killed, damage, accumMult)

		// 全服廣播：錘擊目標
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyLightningHammer,
			Payload: ws.LuckyLightningHammerPayload{
				Event:      "hammer_hit",
				PlayerID:   p.ID,
				PlayerName: p.DisplayName,
				InstanceID: iid,
				DefID:      defID,
				HPDamage:   damage,
				Killed:     killed,
				AccumMult:  accumMult,
			},
		})
	}

	// 結算
	g.doLightningHammerSettle(p, hitCount, killCount, accumMult, totalReward)
}

// doLightningHammerSettle 閃電錘結算
func (g *Game) doLightningHammerSettle(p *player.Player, hitCount, killCount int, finalMult float64, totalReward int) {
	log.Printf("[LuckyLightningHammer] settle: player=%s hits=%d kills=%d finalMult=%.2f totalReward=%d",
		p.ID, hitCount, killCount, finalMult, totalReward)

	// 全服廣播：錘擊結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyLightningHammer,
		Payload: ws.LuckyLightningHammerPayload{
			Event:       "hammer_end",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			HitCount:    hitCount,
			KillCount:   killCount,
			FinalMult:   finalMult,
			TotalReward: totalReward,
		},
	})

	// 全服廣播橫幅
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyLightningHammer,
		Payload: ws.LuckyLightningHammerPayload{
			Event:      "hammer_broadcast",
			PlayerName: p.DisplayName,
			HitCount:   hitCount,
			KillCount:  killCount,
			FinalMult:  finalMult,
		},
	})

	// 全服公告
	var annMsg string
	var annColor string
	if killCount >= 3 {
		annMsg = fmt.Sprintf("⚡ %s 的閃電錘大爆發！錘擊 %d 個目標，直接擊破 %d 個！累積倍率 ×%.1f！",
			p.DisplayName, hitCount, killCount, finalMult)
		annColor = "#FFD700" // 金色
	} else if killCount >= 1 {
		annMsg = fmt.Sprintf("⚡ 閃電錘結算！%s 錘擊 %d 個，擊破 %d 個，倍率 ×%.1f",
			p.DisplayName, hitCount, killCount, finalMult)
		annColor = "#FFA500" // 橙色
	} else {
		annMsg = fmt.Sprintf("⚡ 閃電錘結算！%s 錘擊 %d 個目標，倍率 ×%.1f",
			p.DisplayName, hitCount, finalMult)
		annColor = "#87CEEB" // 天藍
	}
	ann := g.Announce.Create(announce.EventLuckyLightningHammer, p.DisplayName, 0, map[string]string{
		"message": annMsg,
		"color":   annColor,
	})
	g.broadcastAnnouncement(ann)
}
