// lucky_golden_hurricane_handler.go — 幸運黃金颶風魚系統（DAY-276）
// 業界依據：Royal Fishing Jili 2026「AOE 旋風掃場」機制進化版
//
// 設計：擊破 T234 後，觸發「黃金颶風」（螺旋掃場，持續 6 秒）：
//   - 颶風以螺旋路徑掃過整個場地，路徑上所有目標 HP -30%
//   - 颶風每掃過一個目標，觸發玩家獲得 ×1.5 倍率加成（累積，最高 ×8.0）
//   - 6 秒後颶風結算：廣播掃過目標數/累積倍率/總獎勵
//   - 個人冷卻 28 秒；全服冷卻 45 秒
//
// 設計差異：
//   - 與連鎖爆炸（T224，從一點向外擴散）不同，黃金颶風是「螺旋路徑掃場」，
//     讓玩家看到「颶風從中心螺旋向外掃過整個場地」的視覺爽感
//   - 「每掃過一個目標 ×1.5 累積」讓玩家有「颶風掃的目標越多，倍率越高」的期待感
//   - 「最高 ×8.0 累積倍率」讓玩家有「要是颶風掃過 5 個目標就賺大了」的動力
//   - 「HP -30% 弱化」讓颶風後的目標更容易擊破，製造「颶風過後趁機打」的策略感
//   - 「全服廣播颶風路徑」讓所有玩家看到「颶風在哪裡掃」，製造全服緊張感
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
	LuckyGoldenHurricanePersonalCD  = 28 * time.Second // 個人冷卻
	LuckyGoldenHurricaneGlobalCD    = 45 * time.Second // 全服冷卻
	LuckyGoldenHurricaneDuration    = 6 * time.Second  // 颶風持續時間
	LuckyGoldenHurricaneHPDamage    = 0.30             // 颶風 HP 傷害比例（-30%）
	LuckyGoldenHurricaneMultPerHit  = 1.5              // 每掃過一個目標的倍率加成（乘法）
	LuckyGoldenHurricaneMaxMult     = 8.0              // 最高累積倍率
	LuckyGoldenHurricaneSweepDelay  = 600 * time.Millisecond // 每次掃過間隔（6秒/最多10個目標）
)

// luckyGoldenHurricaneManager 幸運黃金颶風魚管理器
type luckyGoldenHurricaneManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 颶風狀態
	active      bool
	activeUntil time.Time
	instanceID  string

	// 觸發玩家累積倍率（playerID → accumMult）
	accumMults map[string]float64

	// 觸發玩家總獎勵（playerID → totalReward）
	totalRewards map[string]int
}

func newLuckyGoldenHurricaneManager() *luckyGoldenHurricaneManager {
	return &luckyGoldenHurricaneManager{
		personalCooldowns: make(map[string]time.Time),
		accumMults:        make(map[string]float64),
		totalRewards:      make(map[string]int),
	}
}

// isLuckyGoldenHurricaneFish 判斷是否為幸運黃金颶風魚
func isLuckyGoldenHurricaneFish(defID string) bool {
	return defID == "T234"
}

// getGoldenHurricaneMult 取得颶風累積倍率（供 handleKill 使用）
// 颶風期間，觸發玩家的每次擊破都套用累積倍率
func (g *Game) getGoldenHurricaneMult(playerID string) float64 {
	mgr := g.LuckyGoldenHurricane
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if !mgr.active || time.Now().After(mgr.activeUntil) {
		return 1.0
	}
	mult, ok := mgr.accumMults[playerID]
	if !ok || mult <= 1.0 {
		return 1.0
	}
	return mult
}

// recordGoldenHurricaneReward 記錄颶風期間獎勵（供 handleKill 使用）
func (g *Game) recordGoldenHurricaneReward(playerID string, reward int) {
	mgr := g.LuckyGoldenHurricane
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if !mgr.active || time.Now().After(mgr.activeUntil) {
		return
	}
	if _, ok := mgr.totalRewards[playerID]; ok {
		mgr.totalRewards[playerID] += reward
	}
}

// tryLuckyGoldenHurricaneFish 擊破 T234 後觸發黃金颶風（供 handleKill 使用）
func (g *Game) tryLuckyGoldenHurricaneFish(p *player.Player) {
	mgr := g.LuckyGoldenHurricane
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
	// 已有颶風啟動中
	if mgr.active && time.Now().Before(mgr.activeUntil) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyGoldenHurricanePersonalCD)
	mgr.globalCooldownUntil = time.Now().Add(LuckyGoldenHurricaneGlobalCD)

	// 啟動颶風
	mgr.active = true
	mgr.activeUntil = time.Now().Add(LuckyGoldenHurricaneDuration)
	instanceID := fmt.Sprintf("hurricane_%d", time.Now().UnixNano())
	mgr.instanceID = instanceID

	// 初始化觸發玩家的累積倍率和總獎勵
	mgr.accumMults[p.ID] = 1.0
	mgr.totalRewards[p.ID] = 0
	mgr.mu.Unlock()

	log.Printf("[LuckyGoldenHurricane] player=%s activated golden hurricane for %v", p.ID, LuckyGoldenHurricaneDuration)

	// 取得場上目標數量
	g.mu.RLock()
	targetCount := 0
	for _, t := range g.Targets {
		if t.HP > 0 {
			targetCount++
		}
	}
	g.mu.RUnlock()

	// 全服廣播：颶風開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyGoldenHurricane,
		Payload: ws.LuckyGoldenHurricanePayload{
			Event:       "hurricane_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyGoldenHurricaneDuration.Seconds()),
			TargetCount: targetCount,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyGoldenHurricane, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌪️ %s 觸發黃金颶風！螺旋掃場 %d 秒，每掃過一個目標 ×%.1f 累積倍率！",
			p.DisplayName, int(LuckyGoldenHurricaneDuration.Seconds()), LuckyGoldenHurricaneMultPerHit),
		"color": "#FFD700",
	})
	g.broadcastAnnouncement(ann)

	// 啟動颶風螺旋掃場 goroutine
	go g.runGoldenHurricaneSweep(p, instanceID)
}

// runGoldenHurricaneSweep 颶風螺旋掃場（goroutine）
// 每隔 600ms 掃過一個隨機目標，最多掃 10 個（6秒/600ms = 10次）
func (g *Game) runGoldenHurricaneSweep(p *player.Player, instanceID string) {
	maxSweeps := int(LuckyGoldenHurricaneDuration / LuckyGoldenHurricaneSweepDelay)
	sweptCount := 0

	for i := 0; i < maxSweeps; i++ {
		time.Sleep(LuckyGoldenHurricaneSweepDelay)

		// 確認 instanceID 仍有效
		g.LuckyGoldenHurricane.mu.Lock()
		if g.LuckyGoldenHurricane.instanceID != instanceID {
			g.LuckyGoldenHurricane.mu.Unlock()
			break
		}
		g.LuckyGoldenHurricane.mu.Unlock()

		// 隨機選一個存活目標
		g.mu.Lock()
		var candidates []string
		for iid, t := range g.Targets {
			if t.HP > 0 && !isLuckyGoldenHurricaneFish(t.Def.ID) {
				candidates = append(candidates, iid)
			}
		}
		if len(candidates) == 0 {
			g.mu.Unlock()
			continue
		}

		// 隨機選一個目標
		targetIID := candidates[rand.Intn(len(candidates))]
		t, ok := g.Targets[targetIID]
		if !ok || t.HP <= 0 {
			g.mu.Unlock()
			continue
		}

		// 計算 HP 傷害（-30%，最少 1）
		damage := int(float64(t.HP) * LuckyGoldenHurricaneHPDamage)
		if damage < 1 {
			damage = 1
		}
		t.HP -= damage
		if t.HP < 1 {
			t.HP = 1
		}
		defID := t.Def.ID
		g.mu.Unlock()

		sweptCount++

		// 更新觸發玩家的累積倍率
		g.LuckyGoldenHurricane.mu.Lock()
		currentMult := g.LuckyGoldenHurricane.accumMults[p.ID]
		newMult := currentMult * LuckyGoldenHurricaneMultPerHit
		if newMult > LuckyGoldenHurricaneMaxMult {
			newMult = LuckyGoldenHurricaneMaxMult
		}
		g.LuckyGoldenHurricane.accumMults[p.ID] = newMult
		g.LuckyGoldenHurricane.mu.Unlock()

		log.Printf("[LuckyGoldenHurricane] sweep #%d: target=%s hp_damage=%d accum_mult=%.2f",
			sweptCount, targetIID, damage, newMult)

		// 全服廣播：颶風掃過目標
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyGoldenHurricane,
			Payload: ws.LuckyGoldenHurricanePayload{
				Event:      "hurricane_sweep",
				PlayerID:   p.ID,
				PlayerName: p.DisplayName,
				InstanceID: targetIID,
				DefID:      defID,
				HPDamage:   damage,
				AccumMult:  newMult,
			},
		})
	}

	// 颶風結束，結算
	g.doGoldenHurricaneSettle(p, instanceID, sweptCount)
}

// doGoldenHurricaneSettle 颶風結算
func (g *Game) doGoldenHurricaneSettle(p *player.Player, instanceID string, sweptCount int) {
	g.LuckyGoldenHurricane.mu.Lock()
	if g.LuckyGoldenHurricane.instanceID != instanceID {
		g.LuckyGoldenHurricane.mu.Unlock()
		return
	}
	g.LuckyGoldenHurricane.active = false
	finalMult := g.LuckyGoldenHurricane.accumMults[p.ID]
	totalReward := g.LuckyGoldenHurricane.totalRewards[p.ID]
	// 清理
	delete(g.LuckyGoldenHurricane.accumMults, p.ID)
	delete(g.LuckyGoldenHurricane.totalRewards, p.ID)
	g.LuckyGoldenHurricane.mu.Unlock()

	log.Printf("[LuckyGoldenHurricane] settle: player=%s swept=%d finalMult=%.2f totalReward=%d",
		p.ID, sweptCount, finalMult, totalReward)

	// 全服廣播：颶風結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyGoldenHurricane,
		Payload: ws.LuckyGoldenHurricanePayload{
			Event:       "hurricane_end",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			SweptCount:  sweptCount,
			FinalMult:   finalMult,
			TotalReward: totalReward,
		},
	})

	// 全服廣播橫幅
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyGoldenHurricane,
		Payload: ws.LuckyGoldenHurricanePayload{
			Event:      "hurricane_broadcast",
			PlayerName: p.DisplayName,
			SweptCount: sweptCount,
			FinalMult:  finalMult,
		},
	})

	// 全服公告
	var annMsg string
	var annColor string
	if finalMult >= 5.0 {
		annMsg = fmt.Sprintf("🌪️ %s 的黃金颶風結算！掃過 %d 個目標，累積倍率 ×%.1f！",
			p.DisplayName, sweptCount, finalMult)
		annColor = "#FFD700" // 金色
	} else {
		annMsg = fmt.Sprintf("🌪️ 黃金颶風結束！%s 掃過 %d 個目標，倍率 ×%.1f",
			p.DisplayName, sweptCount, finalMult)
		annColor = "#FFA500" // 橙色
	}
	ann := g.Announce.Create(announce.EventLuckyGoldenHurricane, p.DisplayName, 0, map[string]string{
		"message": annMsg,
		"color":   annColor,
	})
	g.broadcastAnnouncement(ann)
}
