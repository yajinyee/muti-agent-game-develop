// lucky_quantum_fish_handler.go — 幸運量子魚系統（DAY-228）
// 業界原創「量子疊加態」機制
//
// 設計：擊破 T186 後觸發「量子疊加」：
//   - 場上隨機 4 個目標進入「量子態」（同時疊加高倍率 ×3.0 和低倍率 ×0.8）
//   - 玩家「觀測」（射擊命中）量子態目標時，50% 機率坍縮為高倍率（×3.0），50% 機率坍縮為低倍率（×0.8）
//   - 量子態持續 10 秒；10 秒後所有未被觀測的量子態目標「量子爆炸」（70% 擊破機率，倍率隨機 ×1.0-×4.0）
//   - 個人冷卻 20 秒；全服廣播量子態建立/坍縮/爆炸
//
// 設計差異：
//   - 與彩虹稜鏡魚（DAY-213，染色固定倍率）不同，量子魚是「不確定性」，讓玩家有「要不要賭一把」的刺激感
//   - 與幸運鏡像魚（DAY-215，複製分身）不同，量子魚是「疊加態坍縮」，每次觀測結果不同
//   - 「50% 高倍率 / 50% 低倍率」讓玩家有「薛丁格的魚」的緊張感
//   - 「量子爆炸隨機倍率 ×1.0-×4.0」讓未被觀測的目標有驚喜感
//   - 視覺上量子態目標在高倍率（紫色）和低倍率（灰色）之間閃爍，製造「疊加態」的視覺效果
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
	LuckyQuantumPersonalCD  = 20 * time.Second // 個人冷卻
	LuckyQuantumDuration    = 10 * time.Second // 量子態持續時間
	LuckyQuantumHighMult    = 3.0              // 高倍率坍縮（×3.0）
	LuckyQuantumLowMult     = 0.8              // 低倍率坍縮（×0.8）
	LuckyQuantumBlastChance = 0.70             // 量子爆炸擊破機率
	LuckyQuantumMaxTargets  = 4               // 最多量子態目標數
)

// quantumEntry 量子態目標記錄
type quantumEntry struct {
	instanceID string
	expiresAt  time.Time
}

// luckyQuantumFishManager 幸運量子魚管理器
type luckyQuantumFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldown map[string]time.Time

	// 量子態目標（targetID → quantumEntry）
	quantumTargets map[string]*quantumEntry

	// 當前 instance ID
	currentInstanceID string
}

func newLuckyQuantumFishManager() *luckyQuantumFishManager {
	return &luckyQuantumFishManager{
		personalCooldown: make(map[string]time.Time),
		quantumTargets:   make(map[string]*quantumEntry),
	}
}

// isLuckyQuantumFish 判斷是否為幸運量子魚
func isLuckyQuantumFish(defID string) bool {
	return defID == "T186"
}

// isQuantumTarget 判斷目標是否處於量子態（供 handleKill 使用）
func (g *Game) isQuantumTarget(targetID string) bool {
	mgr := g.LuckyQuantumFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	entry, ok := mgr.quantumTargets[targetID]
	if !ok {
		return false
	}
	return time.Now().Before(entry.expiresAt)
}

// getLuckyQuantumCollapseMult 取得量子坍縮倍率（供 handleKill 使用）
// 50% 機率高倍率（×3.0），50% 機率低倍率（×0.8）
func (g *Game) getLuckyQuantumCollapseMult(targetID string) (float64, bool) {
	mgr := g.LuckyQuantumFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	entry, ok := mgr.quantumTargets[targetID]
	if !ok {
		return 1.0, false
	}
	if time.Now().After(entry.expiresAt) {
		return 1.0, false
	}
	// 量子坍縮：50% 高倍率，50% 低倍率
	if rand.Float64() < 0.5 {
		return LuckyQuantumHighMult, true
	}
	return LuckyQuantumLowMult, true
}

// removeQuantumEntry 移除量子態目標（被擊破後呼叫）
func (g *Game) removeQuantumEntry(targetID string) {
	mgr := g.LuckyQuantumFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	delete(mgr.quantumTargets, targetID)
}

// notifyQuantumKill 量子態目標被玩家擊破時廣播坍縮結果（供 handleKill 使用）
func (g *Game) notifyQuantumKill(p *player.Player, targetID string, collapseMult float64, reward int) {
	mgr := g.LuckyQuantumFish
	mgr.mu.Lock()
	instanceID := mgr.currentInstanceID
	mgr.mu.Unlock()

	isHigh := collapseMult >= LuckyQuantumHighMult

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyQuantumFish,
		Payload: ws.LuckyQuantumFishPayload{
			Event:          "quantum_collapse",
			InstanceID:     instanceID,
			PlayerID:       p.ID,
			PlayerName:     p.DisplayName,
			TargetID:       targetID,
			CollapseHigh:   isHigh,
			CollapseMult:   collapseMult,
			CollapseReward: reward,
		},
	})

	if isHigh {
		log.Printf("[LuckyQuantum] player=%s observed target=%s → HIGH collapse ×%.1f reward=%d",
			p.ID, targetID, collapseMult, reward)
	} else {
		log.Printf("[LuckyQuantum] player=%s observed target=%s → LOW collapse ×%.1f reward=%d",
			p.ID, targetID, collapseMult, reward)
	}
}

// tryLuckyQuantumFish 擊破 T186 後觸發量子疊加（供 handleKill 使用）
func (g *Game) tryLuckyQuantumFish(p *player.Player) {
	mgr := g.LuckyQuantumFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if until, ok := mgr.personalCooldown[p.ID]; ok && time.Now().Before(until) {
		mgr.mu.Unlock()
		return
	}

	// 設定個人冷卻
	mgr.personalCooldown[p.ID] = time.Now().Add(LuckyQuantumPersonalCD)

	// 建立 instance ID
	instanceID := fmt.Sprintf("quantum_%d", time.Now().UnixNano())
	mgr.currentInstanceID = instanceID
	expiresAt := time.Now().Add(LuckyQuantumDuration)
	mgr.mu.Unlock()

	log.Printf("[LuckyQuantum] player=%s triggered quantum superposition instance=%s", p.ID, instanceID)

	// 選取最多 4 個存活目標進入量子態
	g.mu.Lock()
	candidates := make([]string, 0, 8)
	for id, t := range g.Targets {
		if t.IsAlive && t.Def.Type != "boss" && t.DefID != "T186" {
			candidates = append(candidates, id)
		}
	}
	// 隨機打亂候選目標
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	if len(candidates) > LuckyQuantumMaxTargets {
		candidates = candidates[:LuckyQuantumMaxTargets]
	}
	g.mu.Unlock()

	if len(candidates) == 0 {
		log.Printf("[LuckyQuantum] no candidates for quantum state")
		return
	}

	// 設定量子態
	mgr.mu.Lock()
	for _, id := range candidates {
		mgr.quantumTargets[id] = &quantumEntry{
			instanceID: instanceID,
			expiresAt:  expiresAt,
		}
	}
	mgr.mu.Unlock()

	log.Printf("[LuckyQuantum] set %d targets to quantum state", len(candidates))

	// 全服廣播量子疊加開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyQuantumFish,
		Payload: ws.LuckyQuantumFishPayload{
			Event:        "quantum_start",
			InstanceID:   instanceID,
			PlayerID:     p.ID,
			PlayerName:   p.DisplayName,
			QuantumCount: len(candidates),
			DurationSec:  int(LuckyQuantumDuration.Seconds()),
			HighMult:     LuckyQuantumHighMult,
			LowMult:      LuckyQuantumLowMult,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyQuantumFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⚛️ %s 觸發量子疊加！%d 個目標進入量子態，觀測即坍縮！",
			p.DisplayName, len(candidates)),
		"color": "#9B59B6",
	})
	g.broadcastAnnouncement(ann)

	// 10 秒後觸發量子爆炸
	go func() {
		time.Sleep(LuckyQuantumDuration)

		// 收集仍在量子態的目標（未被觀測）
		mgr.mu.Lock()
		remaining := make([]string, 0)
		for id, entry := range mgr.quantumTargets {
			if entry.instanceID == instanceID {
				remaining = append(remaining, id)
				delete(mgr.quantumTargets, id)
			}
		}
		mgr.mu.Unlock()

		if len(remaining) == 0 {
			log.Printf("[LuckyQuantum] all quantum targets observed before blast")
			return
		}

		log.Printf("[LuckyQuantum] quantum blast! %d unobserved targets", len(remaining))

		// 量子爆炸：70% 擊破機率，倍率隨機 ×1.0-×4.0
		type blastResult struct {
			TargetID string  `json:"target_id"`
			Killed   bool    `json:"killed"`
			Mult     float64 `json:"mult"`
			Reward   int     `json:"reward"`
		}

		results := make([]blastResult, 0, len(remaining))
		totalReward := 0
		blastCount := 0

		g.mu.Lock()
		for _, id := range remaining {
			t, ok := g.Targets[id]
			if !ok || !t.IsAlive {
				continue
			}
			// 70% 擊破機率
			if rand.Float64() < LuckyQuantumBlastChance {
				// 隨機倍率 ×1.0-×4.0
				blastMult := 1.0 + rand.Float64()*3.0
				reward := int(float64(t.Def.MultiplierMax) * blastMult * float64(g.getAvgBetCost()))
				if reward < 1 {
					reward = 1
				}
				t.IsAlive = false
				t.HP = 0
				totalReward += reward
				blastCount++
				results = append(results, blastResult{
					TargetID: id,
					Killed:   true,
					Mult:     blastMult,
					Reward:   reward,
				})
			} else {
				results = append(results, blastResult{
					TargetID: id,
					Killed:   false,
					Mult:     0,
					Reward:   0,
				})
			}
		}
		g.mu.Unlock()

		// 全服廣播量子爆炸結算
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyQuantumFish,
			Payload: ws.LuckyQuantumFishPayload{
				Event:        "quantum_blast",
				InstanceID:   instanceID,
				BlastResults: results,
				TotalReward:  totalReward,
				BlastCount:   blastCount,
			},
		})

		// 全服公告（≥2 個爆炸才公告）
		if blastCount >= 2 {
			color := "#9B59B6"
			if blastCount >= 3 {
				color = "#FF00FF"
			}
			ann2 := g.Announce.Create(announce.EventLuckyQuantumFish, p.DisplayName, 0, map[string]string{
				"message": fmt.Sprintf("⚛️ 量子爆炸！%d 個目標坍縮擊破！",
					blastCount),
				"color": color,
			})
			g.broadcastAnnouncement(ann2)
		}

		log.Printf("[LuckyQuantum] blast complete: %d/%d killed, totalReward=%d",
			blastCount, len(remaining), totalReward)
	}()
}
