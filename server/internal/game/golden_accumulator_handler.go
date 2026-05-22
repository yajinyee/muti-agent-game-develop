// golden_accumulator_handler.go — 黃金累積魚系統（DAY-214）
// 業界依據：Evolution Ice Fishing Live 2026「random multipliers ranging from 2x to 10x
//  to selected wheel segments, creating pathways to maximum 5000x payout」
// 業界原創「全服累積爆發」機制
//
// 設計：T172 黃金累積魚出現後，每次任何玩家擊破任何目標，累積槽 +1（最多 20 點）
//   - 全服廣播累積進度（每 5 點廣播一次）
//   - 累積槽滿（20 點）→ 自動觸發「黃金爆發」：
//       全場所有目標 HP -60%（立即）+ 全服 ×2.0 倍率加成 8 秒
//   - 玩家擊破黃金累積魚本身 → 「提前引爆」（不論累積多少）
//   - 全服冷卻 40 秒；黃金累積魚存活期間持續累積
//
// 設計差異：
//   - 與深海龍王（DAY-208，射擊累積）不同，黃金累積魚是「擊破累積」，
//     讓玩家有「打越多魚，累積越快」的正向回饋
//   - 「提前引爆」讓玩家有「要不要現在打黃金累積魚」的策略決策
//   - 全服 ×2.0 倍率加成讓所有玩家在爆發後 8 秒內都受益，製造「全服一起爽」的高潮感
//   - 全服廣播進度讓玩家感受到「還差幾個就爆發」的期待感
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
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	GoldenAccumulatorTarget   = 20          // 累積目標點數
	GoldenAccumulatorBoostMult = 2.0        // 黃金爆發倍率加成（乘法）
	GoldenAccumulatorBoostSec  = 8          // 黃金爆發持續秒數
	GoldenAccumulatorCooldown  = 40 * time.Second // 全服冷卻
	GoldenAccumulatorHPDrain   = 0.60       // 黃金爆發 HP 削減比例
)

// goldenAccumulatorManager 黃金累積魚管理器
type goldenAccumulatorManager struct {
	mu sync.Mutex

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前累積魚的 instanceID（空字串表示無活躍累積魚）
	activeInstanceID string

	// 累積點數（atomic）
	accumCount int64

	// 黃金爆發倍率加成（active 期間）
	boostActive    bool
	boostUntil     time.Time
}

func newGoldenAccumulatorManager() *goldenAccumulatorManager {
	return &goldenAccumulatorManager{}
}

// isGoldenAccumulatorFish 判斷是否為黃金累積魚
func isGoldenAccumulatorFish(defID string) bool {
	return defID == "T172"
}

// getGoldenAccumulatorBoost 取得黃金爆發倍率加成（供 handleKill 使用）
func (g *Game) getGoldenAccumulatorBoost() float64 {
	mgr := g.GoldenAccumulator
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if mgr.boostActive && time.Now().Before(mgr.boostUntil) {
		return GoldenAccumulatorBoostMult
	}
	if mgr.boostActive {
		mgr.boostActive = false
	}
	return 1.0
}

// notifyGoldenAccumulatorSpawn 黃金累積魚生成時呼叫（由 spawnTarget 觸發）
func (g *Game) notifyGoldenAccumulatorSpawn(t *target.Target) {
	mgr := g.GoldenAccumulator
	mgr.mu.Lock()

	// 全服冷卻檢查
	if time.Now().Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}
	// 已有活躍累積魚
	if mgr.activeInstanceID != "" {
		mgr.mu.Unlock()
		return
	}

	mgr.activeInstanceID = t.InstanceID
	atomic.StoreInt64(&mgr.accumCount, 0)
	mgr.mu.Unlock()

	log.Printf("[GoldenAccumulator] fish spawned: instanceID=%s", t.InstanceID)

	// 全服廣播：黃金累積魚出現
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGoldenAccumulator,
		Payload: ws.GoldenAccumulatorPayload{
			Event:       "accum_appear",
			InstanceID:  t.InstanceID,
			AccumCount:  0,
			AccumTarget: GoldenAccumulatorTarget,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventGoldenAccumulator, "黃金累積魚", 0, map[string]string{
		"message": fmt.Sprintf("🌟 黃金累積魚出現！全服合力擊破 %d 個目標觸發黃金爆發！", GoldenAccumulatorTarget),
		"color":   "#FFD700",
	})
	g.broadcastAnnouncement(ann)
}

// notifyGoldenAccumulatorKill 任何目標被擊破時呼叫（由 handleKill 觸發）
// 累積槽 +1，達到目標時觸發黃金爆發
func (g *Game) notifyGoldenAccumulatorKill() {
	mgr := g.GoldenAccumulator
	mgr.mu.Lock()
	if mgr.activeInstanceID == "" {
		mgr.mu.Unlock()
		return
	}
	mgr.mu.Unlock()

	// 累積 +1
	newCount := atomic.AddInt64(&mgr.accumCount, 1)

	// 每 5 點廣播一次進度
	if newCount%5 == 0 || newCount == int64(GoldenAccumulatorTarget) {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgGoldenAccumulator,
			Payload: ws.GoldenAccumulatorPayload{
				Event:       "accum_progress",
				AccumCount:  int(newCount),
				AccumTarget: GoldenAccumulatorTarget,
			},
		})
	}

	// 達到目標 → 觸發黃金爆發
	if newCount >= int64(GoldenAccumulatorTarget) {
		go g.triggerGoldenAccumulatorBurst(false, "")
	}
}

// notifyGoldenAccumulatorFishKill 黃金累積魚本身被擊破時呼叫（提前引爆）
func (g *Game) notifyGoldenAccumulatorFishKill(p *player.Player) {
	mgr := g.GoldenAccumulator
	mgr.mu.Lock()
	if mgr.activeInstanceID == "" {
		mgr.mu.Unlock()
		return
	}
	mgr.activeInstanceID = ""
	mgr.mu.Unlock()

	log.Printf("[GoldenAccumulator] early detonation by player=%s", p.ID)

	// 全服廣播：提前引爆
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGoldenAccumulator,
		Payload: ws.GoldenAccumulatorPayload{
			Event:      "early_detonate",
			PlayerName: p.DisplayName,
			AccumCount: int(atomic.LoadInt64(&mgr.accumCount)),
		},
	})

	go g.triggerGoldenAccumulatorBurst(true, p.DisplayName)
}

// notifyGoldenAccumulatorLeave 黃金累積魚離開場上（超時）時呼叫
func (g *Game) notifyGoldenAccumulatorLeave(instanceID string) {
	mgr := g.GoldenAccumulator
	mgr.mu.Lock()
	if mgr.activeInstanceID != instanceID {
		mgr.mu.Unlock()
		return
	}
	mgr.activeInstanceID = ""
	mgr.mu.Unlock()

	log.Printf("[GoldenAccumulator] fish left without detonation")

	// 廣播：累積魚離開
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGoldenAccumulator,
		Payload: ws.GoldenAccumulatorPayload{
			Event:      "accum_escape",
			AccumCount: int(atomic.LoadInt64(&mgr.accumCount)),
		},
	})
}

// triggerGoldenAccumulatorBurst 觸發黃金爆發
// isEarly=true 表示提前引爆（玩家擊破累積魚）
func (g *Game) triggerGoldenAccumulatorBurst(isEarly bool, playerName string) {
	mgr := g.GoldenAccumulator

	mgr.mu.Lock()
	// 防止重複觸發
	if mgr.boostActive {
		mgr.mu.Unlock()
		return
	}
	// 清除活躍累積魚（若非提前引爆，這裡清除）
	if !isEarly {
		mgr.activeInstanceID = ""
	}
	// 設定全服冷卻
	mgr.globalCooldownUntil = time.Now().Add(GoldenAccumulatorCooldown)
	// 啟動倍率加成
	mgr.boostActive = true
	mgr.boostUntil = time.Now().Add(time.Duration(GoldenAccumulatorBoostSec) * time.Second)
	accumCount := int(atomic.LoadInt64(&mgr.accumCount))
	atomic.StoreInt64(&mgr.accumCount, 0)
	mgr.mu.Unlock()

	log.Printf("[GoldenAccumulator] burst triggered: isEarly=%v accumCount=%d", isEarly, accumCount)

	// 對場上所有目標 HP -60%
	affectedCount := 0
	g.mu.Lock()
	for _, t := range g.Targets {
		if t.HP > 0 && !isGoldenAccumulatorFish(t.DefID) {
			drain := int(float64(t.HP) * GoldenAccumulatorHPDrain)
			if drain < 1 {
				drain = 1
			}
			t.HP -= drain
			if t.HP < 1 {
				t.HP = 1 // 保留 1 HP，讓玩家有機會擊破
			}
			affectedCount++
		}
	}
	g.mu.Unlock()

	// 廣播：黃金爆發開始
	eventType := "burst_start"
	if isEarly {
		eventType = "early_burst_start"
	}
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGoldenAccumulator,
		Payload: ws.GoldenAccumulatorPayload{
			Event:         eventType,
			PlayerName:    playerName,
			AccumCount:    accumCount,
			AffectedCount: affectedCount,
			BoostMult:     GoldenAccumulatorBoostMult,
			BoostSec:      GoldenAccumulatorBoostSec,
		},
	})

	// 全服公告
	msg := fmt.Sprintf("🌟💥 黃金爆發！全場 %d 個目標 HP -60%%！全服 ×%.0f 倍率加成 %d 秒！",
		affectedCount, GoldenAccumulatorBoostMult, GoldenAccumulatorBoostSec)
	if isEarly {
		msg = fmt.Sprintf("🌟💥 %s 提前引爆黃金累積魚！全場 %d 個目標 HP -60%%！全服 ×%.0f 倍率加成 %d 秒！",
			playerName, affectedCount, GoldenAccumulatorBoostMult, GoldenAccumulatorBoostSec)
	}
	ann := g.Announce.Create(announce.EventGoldenAccumulator, playerName, affectedCount, map[string]string{
		"message": msg,
		"color":   "#FFD700",
	})
	g.broadcastAnnouncement(ann)

	// 等待倍率加成結束後廣播
	time.Sleep(time.Duration(GoldenAccumulatorBoostSec) * time.Second)

	mgr.mu.Lock()
	mgr.boostActive = false
	mgr.mu.Unlock()

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGoldenAccumulator,
		Payload: ws.GoldenAccumulatorPayload{
			Event: "burst_end",
		},
	})

	log.Printf("[GoldenAccumulator] burst ended")
}

// pickGoldenAccumulatorTargets 隨機選取場上目標（供爆炸使用）
func pickGoldenAccumulatorTargets(targets map[string]*target.Target, count int) []string {
	ids := make([]string, 0, len(targets))
	for id, t := range targets {
		if t.HP > 0 {
			ids = append(ids, id)
		}
	}
	rand.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
	if len(ids) > count {
		return ids[:count]
	}
	return ids
}
