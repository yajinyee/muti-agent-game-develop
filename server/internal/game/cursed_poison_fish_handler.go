// cursed_poison_fish_handler.go — 詛咒毒魚系統（DAY-216）
// 業界原創「詛咒反轉」機制
//
// 設計：T174 詛咒毒魚出現後，場上隨機 3 個目標被「詛咒標記」（紫色）
//   - 詛咒目標被擊破：獎勵 ×2.5 倍率（高風險高報酬）
//   - 詛咒目標逃跑：觸發「詛咒懲罰」— 下一次擊破任何目標獎勵 ×0.5（持續 5 秒）
//   - 擊破 T174 本身：「解除詛咒」— 移除所有詛咒標記 + 解咒獎勵 10x betLevel
//   - 個人冷卻 18 秒；全服廣播詛咒標記/解除/懲罰
//
// 設計差異：
//   - 與彩虹稜鏡魚（染色 + 不同倍率）不同，詛咒毒魚是「高風險高報酬」
//     打到詛咒目標 ×2.5，但讓它跑掉就被懲罰 ×0.5
//   - 「詛咒懲罰」讓玩家有「一定要在它跑掉前打到」的緊迫感
//   - 「解除詛咒」讓玩家有「要不要先打毒魚解咒」的策略決策
//   - 全服廣播讓所有玩家都看到詛咒目標，製造「全服競爭搶打詛咒目標」的社交感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	CursedPoisonMaxTargets  = 3                    // 最多詛咒 3 個目標
	CursedPoisonKillMult    = 2.5                  // 詛咒目標被擊破的倍率加成
	CursedPoisonPenaltyMult = 0.5                  // 詛咒懲罰倍率（逃跑時）
	CursedPoisonPenaltySec  = 5                    // 詛咒懲罰持續秒數
	CursedPoisonCleanseBase = 10                   // 解咒基礎獎勵（× betLevel）
	CursedPoisonPersonalCD  = 18 * time.Second     // 個人冷卻
)

// cursedEntry 詛咒目標記錄
type cursedEntry struct {
	instanceID string
	defID      string
	killed     bool
}

// cursedPoisonFishManager 詛咒毒魚管理器
type cursedPoisonFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldown map[string]time.Time

	// 當前詛咒目標（instanceID → cursedEntry）
	cursedTargets map[string]*cursedEntry

	// 個人詛咒懲罰（playerID → penaltyUntil）
	penaltyUntil map[string]time.Time
}

func newCursedPoisonFishManager() *cursedPoisonFishManager {
	return &cursedPoisonFishManager{
		personalCooldown: make(map[string]time.Time),
		cursedTargets:    make(map[string]*cursedEntry),
		penaltyUntil:     make(map[string]time.Time),
	}
}

// isCursedPoisonFish 判斷是否為詛咒毒魚
func isCursedPoisonFish(defID string) bool {
	return defID == "T174"
}

// getCursedPoisonKillMult 取得詛咒目標擊破倍率加成（供 handleKill 使用）
// 若被擊破的目標是詛咒目標，回傳 ×2.5 乘法加成
func (g *Game) getCursedPoisonKillMult(instanceID string) float64 {
	mgr := g.CursedPoisonFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if entry, ok := mgr.cursedTargets[instanceID]; ok && !entry.killed {
		return CursedPoisonKillMult
	}
	return 1.0
}

// getCursedPoisonPenaltyMult 取得詛咒懲罰倍率（供 handleKill 使用）
// 若玩家正在受詛咒懲罰，回傳 ×0.5 乘法
func (g *Game) getCursedPoisonPenaltyMult(playerID string) float64 {
	mgr := g.CursedPoisonFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if until, ok := mgr.penaltyUntil[playerID]; ok && time.Now().Before(until) {
		return CursedPoisonPenaltyMult
	}
	return 1.0
}

// removeCursedEntry 詛咒目標被擊破後移除（供 handleKill 使用）
func (g *Game) removeCursedEntry(instanceID string) {
	mgr := g.CursedPoisonFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if entry, ok := mgr.cursedTargets[instanceID]; ok {
		entry.killed = true
		delete(mgr.cursedTargets, instanceID)
		log.Printf("[CursedPoisonFish] cursed target killed: instanceID=%s", instanceID)
	}
}

// isCursedTarget 判斷 instanceID 是否為詛咒目標
func (g *Game) isCursedTarget(instanceID string) bool {
	mgr := g.CursedPoisonFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	_, ok := mgr.cursedTargets[instanceID]
	return ok
}

// notifyCursedTargetEscape 詛咒目標逃跑時呼叫（由 gameLoop 觸發）
// 觸發詛咒懲罰：下一次擊破任何目標獎勵 ×0.5（持續 5 秒）
func (g *Game) notifyCursedTargetEscape(instanceID string) {
	mgr := g.CursedPoisonFish
	mgr.mu.Lock()

	entry, ok := mgr.cursedTargets[instanceID]
	if !ok {
		mgr.mu.Unlock()
		return
	}
	delete(mgr.cursedTargets, instanceID)
	mgr.mu.Unlock()

	log.Printf("[CursedPoisonFish] cursed target escaped: instanceID=%s defID=%s", instanceID, entry.defID)

	// 對所有玩家施加詛咒懲罰
	g.mu.RLock()
	playerIDs := make([]string, 0, len(g.Players))
	for id := range g.Players {
		playerIDs = append(playerIDs, id)
	}
	g.mu.RUnlock()

	penaltyUntil := time.Now().Add(time.Duration(CursedPoisonPenaltySec) * time.Second)
	mgr.mu.Lock()
	for _, pid := range playerIDs {
		mgr.penaltyUntil[pid] = penaltyUntil
	}
	mgr.mu.Unlock()

	// 廣播：詛咒懲罰
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCursedPoisonFish,
		Payload: ws.CursedPoisonFishPayload{
			Event:       "curse_escape",
			InstanceID:  instanceID,
			PenaltyMult: CursedPoisonPenaltyMult,
			PenaltySec:  CursedPoisonPenaltySec,
		},
	})

	log.Printf("[CursedPoisonFish] penalty applied to %d players for %ds", len(playerIDs), CursedPoisonPenaltySec)
}

// tryCleanseAllCurses 擊破 T174 本身時解除所有詛咒
func (g *Game) tryCleanseAllCurses(p *player.Player) {
	mgr := g.CursedPoisonFish
	mgr.mu.Lock()

	if len(mgr.cursedTargets) == 0 {
		mgr.mu.Unlock()
		return
	}

	// 收集詛咒目標列表
	killedCount := len(mgr.cursedTargets)
	// 清除所有詛咒
	for id := range mgr.cursedTargets {
		delete(mgr.cursedTargets, id)
	}
	// 清除玩家的詛咒懲罰
	delete(mgr.penaltyUntil, p.ID)
	mgr.mu.Unlock()

	// 解咒獎勵：10x betLevel
	betDef := data.GetBetDef(p.BetLevel)
	cleanseReward := CursedPoisonCleanseBase
	if betDef != nil {
		cleanseReward = CursedPoisonCleanseBase * betDef.BetCost
	}
	p.AddReward(cleanseReward)

	log.Printf("[CursedPoisonFish] cleanse by player=%s: %d curses removed, reward=%d",
		p.ID, killedCount, cleanseReward)

	// 廣播：解除詛咒
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCursedPoisonFish,
		Payload: ws.CursedPoisonFishPayload{
			Event:         "curse_cleanse",
			PlayerName:    p.DisplayName,
			KilledCount:   killedCount,
			CleanseReward: cleanseReward,
		},
	})

	// 發放解咒獎勵通知
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgReward,
		Payload: ws.RewardPayload{
			Source:     "curse_cleanse",
			Amount:     cleanseReward,
			Multiplier: float64(CursedPoisonCleanseBase),
			NewBalance: p.Coins,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventCursedPoisonFish, p.DisplayName, killedCount, map[string]string{
		"message": fmt.Sprintf("✨ %s 解除詛咒！%d 個詛咒目標解除！獲得解咒獎勵 %d 金幣！",
			p.DisplayName, killedCount, cleanseReward),
		"color": "#E8D5FF",
	})
	g.broadcastAnnouncement(ann)
}

// tryCursedPoisonFish 詛咒毒魚出現時觸發詛咒標記
func (g *Game) tryCursedPoisonFish(p *player.Player) {
	mgr := g.CursedPoisonFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldown[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}
	mgr.personalCooldown[p.ID] = time.Now().Add(CursedPoisonPersonalCD)
	mgr.mu.Unlock()

	// 選取場上最多 3 個目標（排除 BOSS 和詛咒毒魚本身）
	g.mu.RLock()
	candidates := make([]*target.Target, 0, 10)
	for _, t := range g.Targets {
		if t.HP > 0 && t.DefID != "T174" && t.DefID != "B001" {
			if !g.isCursedTarget(t.InstanceID) {
				candidates = append(candidates, t)
			}
		}
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		log.Printf("[CursedPoisonFish] no candidates for curse")
		return
	}

	// 隨機打亂，取最多 3 個
	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	if len(candidates) > CursedPoisonMaxTargets {
		candidates = candidates[:CursedPoisonMaxTargets]
	}

	// 建立詛咒標記
	cursedInfos := make([]ws.CursedTargetInfo, 0, len(candidates))
	mgr.mu.Lock()
	for _, orig := range candidates {
		entry := &cursedEntry{
			instanceID: orig.InstanceID,
			defID:      orig.DefID,
			killed:     false,
		}
		mgr.cursedTargets[orig.InstanceID] = entry
		cursedInfos = append(cursedInfos, ws.CursedTargetInfo{
			InstanceID: orig.InstanceID,
			DefID:      orig.DefID,
			CurseMult:  CursedPoisonKillMult,
		})
	}
	mgr.mu.Unlock()

	log.Printf("[CursedPoisonFish] player=%s triggered curse: %d targets cursed", p.ID, len(cursedInfos))

	// 全服廣播：詛咒標記建立
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCursedPoisonFish,
		Payload: ws.CursedPoisonFishPayload{
			Event:         "curse_start",
			PlayerName:    p.DisplayName,
			CursedTargets: cursedInfos,
			CurseMult:     CursedPoisonKillMult,
			PenaltyMult:   CursedPoisonPenaltyMult,
			PenaltySec:    CursedPoisonPenaltySec,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventCursedPoisonFish, p.DisplayName, len(cursedInfos), map[string]string{
		"message": fmt.Sprintf("☠️ %s 觸發詛咒毒魚！%d 個目標被詛咒（×%.1f 倍率）！打到詛咒目標獲得高倍率，讓它跑掉受懲罰！",
			p.DisplayName, len(cursedInfos), CursedPoisonKillMult),
		"color": "#9B59B6",
	})
	g.broadcastAnnouncement(ann)
}

// spawnCursedPoisonMarks 詛咒毒魚生成時自動詛咒場上目標（不需要玩家擊破）
// 由 spawnTarget 呼叫
func (g *Game) spawnCursedPoisonMarks() {
	mgr := g.CursedPoisonFish

	// 選取場上最多 3 個目標（排除 BOSS 和詛咒毒魚本身）
	g.mu.RLock()
	candidates := make([]*target.Target, 0, 10)
	for _, t := range g.Targets {
		if t.HP > 0 && t.DefID != "T174" && t.DefID != "B001" {
			if !g.isCursedTarget(t.InstanceID) {
				candidates = append(candidates, t)
			}
		}
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		log.Printf("[CursedPoisonFish] no candidates for auto-curse on spawn")
		return
	}

	// 隨機打亂，取最多 3 個
	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	if len(candidates) > CursedPoisonMaxTargets {
		candidates = candidates[:CursedPoisonMaxTargets]
	}

	// 建立詛咒標記
	cursedInfos := make([]ws.CursedTargetInfo, 0, len(candidates))
	mgr.mu.Lock()
	for _, orig := range candidates {
		entry := &cursedEntry{
			instanceID: orig.InstanceID,
			defID:      orig.DefID,
			killed:     false,
		}
		mgr.cursedTargets[orig.InstanceID] = entry
		cursedInfos = append(cursedInfos, ws.CursedTargetInfo{
			InstanceID: orig.InstanceID,
			DefID:      orig.DefID,
			CurseMult:  CursedPoisonKillMult,
		})
	}
	mgr.mu.Unlock()

	log.Printf("[CursedPoisonFish] auto-curse on spawn: %d targets cursed", len(cursedInfos))

	// 全服廣播：詛咒標記建立
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCursedPoisonFish,
		Payload: ws.CursedPoisonFishPayload{
			Event:         "curse_start",
			CursedTargets: cursedInfos,
			CurseMult:     CursedPoisonKillMult,
			PenaltyMult:   CursedPoisonPenaltyMult,
			PenaltySec:    CursedPoisonPenaltySec,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventCursedPoisonFish, "", len(cursedInfos), map[string]string{
		"message": fmt.Sprintf("☠️ 詛咒毒魚出現！%d 個目標被詛咒（×%.1f 倍率）！打到詛咒目標獲得高倍率，讓它跑掉受懲罰！",
			len(cursedInfos), CursedPoisonKillMult),
		"color": "#9B59B6",
	})
	g.broadcastAnnouncement(ann)
}
