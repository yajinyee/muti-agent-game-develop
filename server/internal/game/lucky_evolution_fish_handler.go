// lucky_evolution_fish_handler.go — 幸運進化魚系統（DAY-218）
// 業界原創「三段進化」機制
//
// 設計：T176 幸運進化魚出現後，每次任何玩家命中它（不需要擊破），它會「進化」：
//   - 進化 1（命中 3 次）：HP -30%，倍率 ×1.5，全服廣播「進化！」
//   - 進化 2（命中 6 次）：HP -50%，倍率 ×2.5，全服廣播「二段進化！」
//   - 進化 3（命中 9 次）：HP -70%，倍率 ×4.0，全服廣播「終極進化！」
//   - 進化 3 後 3 秒：自動觸發「終極爆發」— 全場所有目標 HP -60%（保留 1）+ 全服 ×4.0 倍率加成 6 秒
//   - 玩家擊破進化魚本身：立即觸發「終極爆發」（不論進化階段）
//   - 全服冷卻 35 秒
//
// 設計差異：
//   - 與黃金累積魚（擊破累積）不同，進化魚是「命中累積」，讓玩家有「要不要現在打死它」的策略決策
//   - 「越打越強、越打越值錢」讓玩家有「等它進化再打死」的期待感
//   - 「終極爆發」讓玩家有「等待→爆發」的高潮設計
//   - 全服廣播進化階段讓所有玩家都看到進化進度，製造「全服一起等待爆發」的社交感
//   - 進化後倍率提升讓「等待進化」有實質獎勵，不只是視覺效果
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyEvolutionGlobalCD   = 35 * time.Second // 全服冷卻
	LuckyEvolutionBurstDelay = 3 * time.Second  // 終極進化後爆發延遲
	LuckyEvolutionBoostSec   = 6                // 終極爆發倍率加成持續秒數
	LuckyEvolutionBoostMult  = 4.0              // 終極爆發倍率加成（×4.0）
	LuckyEvolutionHPDrain    = 0.60             // 終極爆發 HP 削減比例
)

// evolutionStage 進化階段定義
type evolutionStage struct {
	hitThreshold int     // 觸發所需命中次數
	hpDrainPct   float64 // HP 削減比例
	multBoost    float64 // 倍率加成（乘法）
	name         string  // 階段名稱
	color        string  // 廣播顏色
}

var evolutionStages = []evolutionStage{
	{hitThreshold: 3, hpDrainPct: 0.30, multBoost: 1.5, name: "進化！", color: "#00FF88"},
	{hitThreshold: 6, hpDrainPct: 0.50, multBoost: 2.5, name: "二段進化！", color: "#00CCFF"},
	{hitThreshold: 9, hpDrainPct: 0.70, multBoost: 4.0, name: "終極進化！", color: "#FF00FF"},
}

// luckyEvolutionFishManager 幸運進化魚管理器
type luckyEvolutionFishManager struct {
	mu sync.Mutex

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前進化魚狀態
	active       bool
	instanceID   string
	hitCount     int     // 累積命中次數
	currentStage int     // 當前進化階段（0=未進化，1/2/3）
	currentMult  float64 // 當前倍率加成

	// 終極爆發倍率加成
	burstActive    bool
	burstUntil     time.Time
}

func newLuckyEvolutionFishManager() *luckyEvolutionFishManager {
	return &luckyEvolutionFishManager{}
}

// isLuckyEvolutionFish 判斷是否為幸運進化魚
func isLuckyEvolutionFish(defID string) bool {
	return defID == "T176"
}

// getLuckyEvolutionBurstBoost 取得終極爆發倍率加成（供 handleKill 使用）
func (g *Game) getLuckyEvolutionBurstBoost() float64 {
	mgr := g.LuckyEvolutionFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if mgr.burstActive && time.Now().Before(mgr.burstUntil) {
		return LuckyEvolutionBoostMult
	}
	return 1.0
}

// getLuckyEvolutionKillMult 取得進化魚當前倍率加成（供 handleKill 使用）
// 若被擊破的目標是進化魚本身，回傳當前進化倍率
func (g *Game) getLuckyEvolutionKillMult(instanceID string) float64 {
	mgr := g.LuckyEvolutionFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if mgr.active && mgr.instanceID == instanceID && mgr.currentMult > 1.0 {
		return mgr.currentMult
	}
	return 1.0
}

// notifyLuckyEvolutionFishSpawn 幸運進化魚生成時初始化
func (g *Game) notifyLuckyEvolutionFishSpawn(t *target.Target) {
	mgr := g.LuckyEvolutionFish
	mgr.mu.Lock()

	// 全服冷卻檢查
	if time.Now().Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}

	// 若已有進化魚在場，不重複觸發
	if mgr.active {
		mgr.mu.Unlock()
		return
	}

	mgr.active = true
	mgr.instanceID = t.InstanceID
	mgr.hitCount = 0
	mgr.currentStage = 0
	mgr.currentMult = 1.0
	mgr.globalCooldownUntil = time.Now().Add(LuckyEvolutionGlobalCD)
	mgr.mu.Unlock()

	log.Printf("[LuckyEvolutionFish] spawned: instanceID=%s", t.InstanceID)

	// 全服廣播：進化魚出現
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyEvolutionFish,
		Payload: ws.LuckyEvolutionFishPayload{
			Event:      "evolution_appear",
			InstanceID: t.InstanceID,
			Stage:      0,
			HitCount:   0,
			NextHit:    evolutionStages[0].hitThreshold,
			MultBoost:  1.0,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyEvolutionFish, "", 0, map[string]string{
		"message": "🌟 幸運進化魚出現！命中 3 次觸發進化！進化後倍率提升，終極進化後全場爆發！",
		"color":   "#00FF88",
	})
	g.broadcastAnnouncement(ann)
}

// notifyLuckyEvolutionFishHit 進化魚被命中時呼叫（由 handleAttack 呼叫）
func (g *Game) notifyLuckyEvolutionFishHit(instanceID string, p *player.Player) {
	mgr := g.LuckyEvolutionFish
	mgr.mu.Lock()

	if !mgr.active || mgr.instanceID != instanceID {
		mgr.mu.Unlock()
		return
	}

	mgr.hitCount++
	hitCount := mgr.hitCount
	currentStage := mgr.currentStage
	mgr.mu.Unlock()

	log.Printf("[LuckyEvolutionFish] hit by player=%s, hitCount=%d", p.ID, hitCount)

	// 檢查是否觸發進化
	for stageIdx, stage := range evolutionStages {
		if stageIdx <= currentStage-1 {
			continue // 已經過了這個階段
		}
		if hitCount >= stage.hitThreshold && stageIdx == currentStage {
			go g.triggerLuckyEvolution(instanceID, stageIdx+1, stage, p)
			return
		}
	}

	// 廣播命中進度
	mgr.mu.Lock()
	nextHit := 0
	if currentStage < len(evolutionStages) {
		nextHit = evolutionStages[currentStage].hitThreshold
	}
	mgr.mu.Unlock()

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyEvolutionFish,
		Payload: ws.LuckyEvolutionFishPayload{
			Event:      "evolution_hit",
			InstanceID: instanceID,
			PlayerName: p.DisplayName,
			HitCount:   hitCount,
			NextHit:    nextHit,
			Stage:      currentStage,
		},
	})
}

// triggerLuckyEvolution 觸發進化
func (g *Game) triggerLuckyEvolution(instanceID string, newStage int, stage evolutionStage, p *player.Player) {
	mgr := g.LuckyEvolutionFish
	mgr.mu.Lock()

	if !mgr.active || mgr.instanceID != instanceID {
		mgr.mu.Unlock()
		return
	}

	mgr.currentStage = newStage
	mgr.currentMult = stage.multBoost
	mgr.mu.Unlock()

	log.Printf("[LuckyEvolutionFish] evolved to stage %d: mult=%.1f", newStage, stage.multBoost)

	// 對進化魚本身削減 HP
	g.mu.Lock()
	t, ok := g.Targets[instanceID]
	if ok && t.HP > 0 {
		drain := int(float64(t.HP) * stage.hpDrainPct)
		t.HP -= drain
		if t.HP < 1 {
			t.HP = 1
		}
	}
	g.mu.Unlock()

	// 全服廣播：進化
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyEvolutionFish,
		Payload: ws.LuckyEvolutionFishPayload{
			Event:      "evolution_stage",
			InstanceID: instanceID,
			PlayerName: p.DisplayName,
			Stage:      newStage,
			StageName:  stage.name,
			MultBoost:  stage.multBoost,
			HitCount:   mgr.hitCount,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyEvolutionFish, p.DisplayName, newStage, map[string]string{
		"message": fmt.Sprintf("🌟 %s 觸發進化魚%s 倍率提升至 ×%.1f！", p.DisplayName, stage.name, stage.multBoost),
		"color":   stage.color,
	})
	g.broadcastAnnouncement(ann)

	// 若達到終極進化（第 3 段），3 秒後自動觸發終極爆發
	if newStage == 3 {
		go func() {
			time.Sleep(LuckyEvolutionBurstDelay)
			g.triggerLuckyEvolutionBurst(instanceID, false)
		}()
	}
}

// notifyLuckyEvolutionFishKill 玩家擊破進化魚本身時立即觸發終極爆發
func (g *Game) notifyLuckyEvolutionFishKill(p *player.Player, instanceID string) {
	mgr := g.LuckyEvolutionFish
	mgr.mu.Lock()
	if !mgr.active || mgr.instanceID != instanceID {
		mgr.mu.Unlock()
		return
	}
	mgr.active = false
	mgr.mu.Unlock()

	log.Printf("[LuckyEvolutionFish] killed by player=%s, triggering burst", p.ID)
	go g.triggerLuckyEvolutionBurst(instanceID, true)
}

// triggerLuckyEvolutionBurst 終極爆發
// isKill=true 表示玩家擊破觸發，isKill=false 表示終極進化後自動觸發
func (g *Game) triggerLuckyEvolutionBurst(instanceID string, isKill bool) {
	mgr := g.LuckyEvolutionFish
	mgr.mu.Lock()

	// 若已不是 active（可能已被其他路徑清除），仍繼續爆發
	mgr.active = false
	mgr.burstActive = true
	mgr.burstUntil = time.Now().Add(time.Duration(LuckyEvolutionBoostSec) * time.Second)
	stage := mgr.currentStage
	mgr.mu.Unlock()

	log.Printf("[LuckyEvolutionFish] burst triggered: isKill=%v stage=%d", isKill, stage)

	// 全場所有目標 HP -60%（保留 1）
	g.mu.Lock()
	affectedCount := 0
	for _, t := range g.Targets {
		if t.HP > 0 && t.DefID != "B001" && t.InstanceID != instanceID {
			drain := int(float64(t.HP) * LuckyEvolutionHPDrain)
			t.HP -= drain
			if t.HP < 1 {
				t.HP = 1
			}
			affectedCount++
		}
	}
	g.mu.Unlock()

	eventName := "evolution_burst"
	if isKill {
		eventName = "evolution_kill_burst"
	}

	// 全服廣播：終極爆發開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyEvolutionFish,
		Payload: ws.LuckyEvolutionFishPayload{
			Event:         eventName,
			InstanceID:    instanceID,
			Stage:         stage,
			MultBoost:     LuckyEvolutionBoostMult,
			BoostSec:      LuckyEvolutionBoostSec,
			AffectedCount: affectedCount,
		},
	})

	// 全服公告
	burstMsg := fmt.Sprintf("💥 幸運進化魚終極爆發！全場 HP -60%%！全服 ×%.1f 倍率加成 %d 秒！",
		LuckyEvolutionBoostMult, LuckyEvolutionBoostSec)
	if isKill {
		burstMsg = fmt.Sprintf("💥 進化魚被擊破！提前引爆！全場 HP -60%%！全服 ×%.1f 倍率加成 %d 秒！",
			LuckyEvolutionBoostMult, LuckyEvolutionBoostSec)
	}
	ann := g.Announce.Create(announce.EventLuckyEvolutionFish, "", affectedCount, map[string]string{
		"message": burstMsg,
		"color":   "#FF00FF",
	})
	g.broadcastAnnouncement(ann)

	// 等待倍率加成結束
	time.Sleep(time.Duration(LuckyEvolutionBoostSec) * time.Second)

	mgr.mu.Lock()
	mgr.burstActive = false
	mgr.mu.Unlock()

	// 廣播：倍率加成結束
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyEvolutionFish,
		Payload: ws.LuckyEvolutionFishPayload{
			Event: "evolution_burst_end",
		},
	})

	log.Printf("[LuckyEvolutionFish] burst ended, affected=%d", affectedCount)
}

// notifyLuckyEvolutionFishLeave 進化魚逃跑時清除狀態
func (g *Game) notifyLuckyEvolutionFishLeave(instanceID string) {
	mgr := g.LuckyEvolutionFish
	mgr.mu.Lock()
	if !mgr.active || mgr.instanceID != instanceID {
		mgr.mu.Unlock()
		return
	}
	mgr.active = false
	stage := mgr.currentStage
	mgr.mu.Unlock()

	log.Printf("[LuckyEvolutionFish] escaped at stage=%d", stage)

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyEvolutionFish,
		Payload: ws.LuckyEvolutionFishPayload{
			Event:      "evolution_escape",
			InstanceID: instanceID,
			Stage:      stage,
		},
	})
}
