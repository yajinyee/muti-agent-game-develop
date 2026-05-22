// dragon_king_handler.go — 深海龍王全服合力蓄力系統（DAY-208）
// 業界依據：Royal Fishing JILI「Dragon Wrath — accumulate wrath value through shooting,
//  then unleash devastating meteor strikes across the entire screen.
//  Includes both Immortal Boss and ChainLong King encounters.」
//
// 設計：擊破 T166 後觸發「龍王怒火蓄力模式」（12秒）：
//   1. 全服所有玩家的每次射擊都累積「龍怒值」（+1/shot）
//   2. 達到 20 點 → 觸發「龍怒隕石雨」（5 顆隕石，350px 半徑，80% 擊破機率，0.75x 倍率）
//   3. 12 秒內未達到 20 點 → 觸發「小型龍怒」（3 顆隕石，250px 半徑，65% 擊破機率，0.60x 倍率）
//   4. 全服冷卻 45 秒
//
// 設計差異（與現有 DragonWrath 個人系統的區別）：
//   - 個人 DragonWrath（DAY-154）：個人蓄力，個人使用，個人觸發
//   - 深海龍王（DAY-208）：全服合力蓄力，任何玩家射擊都累積，製造「大家一起打才能觸發龍怒」的社群感
//   - 「射擊越多越快觸發」讓玩家有「趕快多打幾槍」的緊迫感
//   - 全服廣播讓所有玩家看到蓄力進度條，製造「還差幾槍就爆發」的期待感
//   - 隕石雨全服共享獎勵，讓觸發者有「我幫全服觸發了龍怒」的英雄感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/ws"
)

// 深海龍王常數（DAY-208）
const (
	DragonKingCooldownSec      = 45    // 全服冷卻 45 秒
	DragonKingChargeDuration   = 12    // 蓄力持續時間 12 秒
	DragonKingChargeTarget     = 20    // 達到 20 點觸發龍怒隕石雨
	DragonKingMeteorCount      = 5     // 龍怒隕石雨：5 顆隕石
	DragonKingMeteorRadius     = 350.0 // 龍怒隕石雨：350px 半徑
	DragonKingMeteorKillChance = 0.80  // 龍怒隕石雨：80% 擊破機率
	DragonKingMeteorMult       = 0.75  // 龍怒隕石雨：0.75x 倍率
	DragonKingSmallCount       = 3     // 小型龍怒：3 顆隕石
	DragonKingSmallRadius      = 250.0 // 小型龍怒：250px 半徑
	DragonKingSmallKillChance  = 0.65  // 小型龍怒：65% 擊破機率
	DragonKingSmallMult        = 0.60  // 小型龍怒：0.60x 倍率
	DragonKingMeteorInterval   = 400   // 隕石間隔 400ms
)

// dragonKingManager 深海龍王管理器（全服共享蓄力狀態）
type dragonKingManager struct {
	mu          sync.Mutex
	isCharging  bool      // 是否正在蓄力中
	chargeEnd   time.Time // 蓄力結束時間
	cooldownAt  time.Time // 全服冷卻結束時間
	chargeCount int32     // 當前蓄力值（atomic 操作）
}

func newDragonKingManager() *dragonKingManager {
	return &dragonKingManager{}
}

// isDragonKingFish 判斷是否為深海龍王（T166，DAY-208）
func isDragonKingFish(defID string) bool {
	return defID == "T166"
}

// tryDragonKingCharge 擊破 T166 後觸發龍王怒火蓄力模式
func (g *Game) tryDragonKingCharge(t *target.Target) {
	mgr := g.DragonKing
	mgr.mu.Lock()

	// 全服冷卻檢查
	if time.Now().Before(mgr.cooldownAt) {
		mgr.mu.Unlock()
		return
	}
	if mgr.isCharging {
		mgr.mu.Unlock()
		return
	}
	mgr.isCharging = true
	mgr.chargeEnd = time.Now().Add(DragonKingChargeDuration * time.Second)
	mgr.cooldownAt = time.Now().Add(DragonKingCooldownSec * time.Second)
	atomic.StoreInt32(&mgr.chargeCount, 0)
	mgr.mu.Unlock()

	log.Printf("[DragonKing] charge mode triggered by killing T166 at (%.0f,%.0f)", t.X, t.Y)

	// 廣播蓄力開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgDragonKing,
		Payload: ws.DragonKingPayload{
			Event:       "charge_start",
			ChargeTarget: DragonKingChargeTarget,
			ChargeSec:   DragonKingChargeDuration,
			Current:     0,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventBigWin, "", 0, map[string]string{
		"message": "🐉 深海龍王覺醒！全服合力蓄力！射擊 20 次觸發龍怒隕石雨！",
		"color":   "#FF4500",
	})
	g.broadcastAnnouncement(ann)

	// 啟動蓄力計時 goroutine
	go g.runDragonKingChargeTimer()
}

// notifyDragonKingShot 玩家射擊時累積龍怒值（由 handleAttack 呼叫）
func (g *Game) notifyDragonKingShot() {
	mgr := g.DragonKing
	mgr.mu.Lock()
	if !mgr.isCharging || time.Now().After(mgr.chargeEnd) {
		mgr.mu.Unlock()
		return
	}
	mgr.mu.Unlock()

	// 原子累加
	newCount := atomic.AddInt32(&mgr.chargeCount, 1)

	// 廣播蓄力進度
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgDragonKing,
		Payload: ws.DragonKingPayload{
			Event:       "charge_progress",
			Current:     int(newCount),
			ChargeTarget: DragonKingChargeTarget,
		},
	})

	// 達到目標，觸發龍怒隕石雨
	if int(newCount) >= DragonKingChargeTarget {
		mgr.mu.Lock()
		if mgr.isCharging {
			mgr.isCharging = false
			mgr.mu.Unlock()
			log.Printf("[DragonKing] charge target reached (%d shots), triggering meteor rain!", newCount)
			go g.executeDragonKingMeteorRain(true)
		} else {
			mgr.mu.Unlock()
		}
	}
}

// runDragonKingChargeTimer 蓄力計時器（goroutine）
func (g *Game) runDragonKingChargeTimer() {
	time.Sleep(DragonKingChargeDuration * time.Second)

	mgr := g.DragonKing
	mgr.mu.Lock()
	if !mgr.isCharging {
		// 已被 notifyDragonKingShot 觸發，不需要再處理
		mgr.mu.Unlock()
		return
	}
	mgr.isCharging = false
	finalCount := int(atomic.LoadInt32(&mgr.chargeCount))
	mgr.mu.Unlock()

	log.Printf("[DragonKing] charge timer expired with %d shots, triggering small meteor", finalCount)
	go g.executeDragonKingMeteorRain(false)
}

// executeDragonKingMeteorRain 執行龍怒隕石雨
// isFull=true：龍怒隕石雨（5顆，350px，80%，0.75x）
// isFull=false：小型龍怒（3顆，250px，65%，0.60x）
func (g *Game) executeDragonKingMeteorRain(isFull bool) {
	meteorCount := DragonKingSmallCount
	radius := DragonKingSmallRadius
	killChance := DragonKingSmallKillChance
	mult := DragonKingSmallMult
	eventName := "small_meteor"

	if isFull {
		meteorCount = DragonKingMeteorCount
		radius = DragonKingMeteorRadius
		killChance = DragonKingMeteorKillChance
		mult = DragonKingMeteorMult
		eventName = "meteor_rain"
	}

	// 廣播隕石雨開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgDragonKing,
		Payload: ws.DragonKingPayload{
			Event:       eventName + "_start",
			MeteorCount: meteorCount,
			IsFull:      isFull,
		},
	})

	totalKills := 0
	totalReward := 0

	for i := 0; i < meteorCount; i++ {
		time.Sleep(DragonKingMeteorInterval * time.Millisecond)

		// 隨機選擇隕石落點（全場範圍）
		meteorX := 100.0 + rand.Float64()*1080.0
		meteorY := 100.0 + rand.Float64()*520.0

		kills, reward := g.doDragonKingMeteorBlast(meteorX, meteorY, radius, killChance, mult)
		totalKills += kills
		totalReward += reward

		// 廣播單顆隕石結果
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgDragonKing,
			Payload: ws.DragonKingPayload{
				Event:     "meteor_hit",
				MeteorIdx: i,
				MeteorX:   meteorX,
				MeteorY:   meteorY,
				KillCount: kills,
				Reward:    reward,
				IsFull:    isFull,
			},
		})
	}

	log.Printf("[DragonKing] %s complete: kills=%d reward=%d", eventName, totalKills, totalReward)

	// 廣播結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgDragonKing,
		Payload: ws.DragonKingPayload{
			Event:       eventName + "_result",
			TotalKills:  totalKills,
			TotalReward: totalReward,
			IsFull:      isFull,
		},
	})

	// 全服公告
	if isFull && totalKills >= 5 {
		color := "#FF4500"
		if totalKills >= 15 {
			color = "#FF0000"
		} else if totalKills >= 10 {
			color = "#FF6600"
		}
		msg := fmt.Sprintf("🐉 龍怒隕石雨！全服合力 %d 槍！擊破 %d 個目標！",
			DragonKingChargeTarget, totalKills)
		ann := g.Announce.Create(announce.EventMegaWin, "", totalReward, map[string]string{
			"message": msg,
			"color":   color,
		})
		g.broadcastAnnouncement(ann)
	} else if !isFull && totalKills >= 3 {
		msg := fmt.Sprintf("🐉 小型龍怒！擊破 %d 個目標！", totalKills)
		ann := g.Announce.Create(announce.EventBigWin, "", totalReward, map[string]string{
			"message": msg,
			"color":   "#FF8C00",
		})
		g.broadcastAnnouncement(ann)
	}
}

// doDragonKingMeteorBlast 執行單顆隕石爆炸
func (g *Game) doDragonKingMeteorBlast(cx, cy, radius, killChance, mult float64) (kills, reward int) {
	type candidate struct {
		instanceID string
		multiplier float64
	}

	g.mu.RLock()
	var candidates []candidate
	for _, t := range g.Targets {
		if !t.IsAlive || isDragonKingFish(t.DefID) {
			continue
		}
		dx := t.X - cx
		dy := t.Y - cy
		if dx*dx+dy*dy <= radius*radius {
			candidates = append(candidates, candidate{t.InstanceID, t.Multiplier})
		}
	}
	betLevel := 1
	for _, p := range g.Players {
		betLevel = p.BetLevel
		break
	}
	g.mu.RUnlock()

	for _, c := range candidates {
		if rand.Float64() >= killChance {
			continue
		}
		rewardAmt := int(float64(betLevel) * c.multiplier * mult)
		if rewardAmt < 1 {
			rewardAmt = 1
		}
		g.mu.Lock()
		if tgt, ok := g.Targets[c.instanceID]; ok && tgt.IsAlive {
			tgt.IsAlive = false
			tgt.HP = 0
			delete(g.Targets, c.instanceID)
			kills++
			reward += rewardAmt
			// 全服共享獎勵
			g.distributeRewardToAll(rewardAmt)
		}
		g.mu.Unlock()
	}
	return
}
