// lucky_infection_fish_handler.go — 幸運連鎖感染魚系統（DAY-219）
// 業界原創「病毒式蔓延」機制
//
// 設計：擊破 T177 後觸發「感染標記」：
//   - 場上隨機 2 個目標被「感染」（綠色標記）
//   - 感染目標每 2 秒向相鄰目標（300px 內）傳播感染（最多蔓延 3 層，最多 8 個感染目標）
//   - 感染目標被擊破：獎勵 ×2.0 倍率加成（乘法）
//   - 12 秒後所有感染目標同時「感染爆發」（75% 擊破機率，0.65x 倍率，全服共享）
//   - 個人冷卻 22 秒；全服廣播感染建立/蔓延/爆發
//
// 設計差異：
//   - 與彩虹稜鏡魚（染色固定目標）不同，感染魚是「動態蔓延」，讓玩家看到感染逐漸擴散的過程
//   - 「越多感染目標，爆發越強」讓玩家有「等待蔓延再打死感染目標」的策略決策
//   - 「感染蔓延」讓玩家有「感染在擴散，要趕快打還是等它蔓延更多？」的緊迫感
//   - 全服廣播感染蔓延讓所有玩家都看到感染進度，製造「全服一起等待爆發」的社交感
//   - 最多 8 個感染目標確保 RTP 平衡，不會無限蔓延
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyInfectionPersonalCD  = 22 * time.Second // 個人冷卻
	LuckyInfectionDuration    = 12 * time.Second // 感染持續時間
	LuckyInfectionSpreadCD    = 2 * time.Second  // 感染蔓延間隔
	LuckyInfectionSpreadRange = 300.0            // 感染蔓延範圍（px）
	LuckyInfectionMaxTargets  = 8                // 最大感染目標數
	LuckyInfectionMaxLayers   = 3                // 最大蔓延層數
	LuckyInfectionKillMult    = 2.0              // 感染目標擊破倍率加成
	LuckyInfectionBlastChance = 0.75             // 感染爆發擊破機率
	LuckyInfectionBlastMult   = 0.65             // 感染爆發倍率
)

// infectionEntry 感染目標記錄
type infectionEntry struct {
	instanceID string
	layer      int // 感染層數（0=初始，1/2/3=蔓延層）
	infectedAt time.Time
}

// luckyInfectionFishManager 幸運連鎖感染魚管理器
type luckyInfectionFishManager struct {
	mu sync.Mutex

	// 個人冷卻
	personalCooldowns map[string]time.Time

	// 當前感染狀態
	active          bool
	sessionID       string
	infectedTargets map[string]*infectionEntry // instanceID -> entry
	spreadLayer     int                        // 當前蔓延層數
	startedAt       time.Time
}

func newLuckyInfectionFishManager() *luckyInfectionFishManager {
	return &luckyInfectionFishManager{
		personalCooldowns: make(map[string]time.Time),
		infectedTargets:   make(map[string]*infectionEntry),
	}
}

// isLuckyInfectionFish 判斷是否為幸運連鎖感染魚
func isLuckyInfectionFish(defID string) bool {
	return defID == "T177"
}

// getLuckyInfectionKillMult 取得感染目標擊破倍率加成（供 handleKill 使用）
func (g *Game) getLuckyInfectionKillMult(instanceID string) float64 {
	mgr := g.LuckyInfectionFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if !mgr.active {
		return 1.0
	}
	if _, ok := mgr.infectedTargets[instanceID]; ok {
		return LuckyInfectionKillMult
	}
	return 1.0
}

// removeInfectionEntry 感染目標被擊破後移除（供 handleKill 使用）
func (g *Game) removeInfectionEntry(instanceID string) {
	mgr := g.LuckyInfectionFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	delete(mgr.infectedTargets, instanceID)
}

// isInfectedTarget 判斷是否為感染目標（供 handleKill 使用）
func (g *Game) isInfectedTarget(instanceID string) bool {
	mgr := g.LuckyInfectionFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if !mgr.active {
		return false
	}
	_, ok := mgr.infectedTargets[instanceID]
	return ok
}

// tryLuckyInfectionFish 擊破 T177 後觸發感染（供 handleKill 使用）
func (g *Game) tryLuckyInfectionFish(p *player.Player) {
	mgr := g.LuckyInfectionFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok {
		if time.Now().Before(cd) {
			mgr.mu.Unlock()
			return
		}
	}

	// 若已有感染進行中，不重複觸發
	if mgr.active {
		mgr.mu.Unlock()
		return
	}

	// 設定個人冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyInfectionPersonalCD)

	// 初始化感染狀態
	mgr.active = true
	mgr.sessionID = fmt.Sprintf("inf_%d", time.Now().UnixNano())
	mgr.infectedTargets = make(map[string]*infectionEntry)
	mgr.spreadLayer = 0
	mgr.startedAt = time.Now()

	// 選取初始感染目標（隨機 2 個非 T177 目標）
	g.mu.RLock()
	var candidates []string
	for id, t := range g.Targets {
		if t.HP > 0 && !isLuckyInfectionFish(t.DefID) {
			candidates = append(candidates, id)
		}
	}
	g.mu.RUnlock()

	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	initialCount := 2
	if len(candidates) < initialCount {
		initialCount = len(candidates)
	}

	var initialTargets []ws.InfectionTargetInfo
	for i := 0; i < initialCount; i++ {
		iid := candidates[i]
		mgr.infectedTargets[iid] = &infectionEntry{
			instanceID: iid,
			layer:      0,
			infectedAt: time.Now(),
		}
		g.mu.RLock()
		t := g.Targets[iid]
		var x, y float64
		if t != nil {
			x, y = t.X, t.Y
		}
		g.mu.RUnlock()
		initialTargets = append(initialTargets, ws.InfectionTargetInfo{
			InstanceID: iid,
			Layer:      0,
			X:          x,
			Y:          y,
		})
	}

	sessionID := mgr.sessionID
	mgr.mu.Unlock()

	if len(initialTargets) == 0 {
		// 沒有可感染目標，重置
		mgr.mu.Lock()
		mgr.active = false
		mgr.mu.Unlock()
		return
	}

	// 全服廣播：感染開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyInfectionFish,
		Payload: ws.LuckyInfectionFishPayload{
			Event:           "infection_start",
			SessionID:       sessionID,
			TriggerPlayer:   p.DisplayName,
			InfectedTargets: initialTargets,
			TotalInfected:   len(initialTargets),
			MaxInfected:     LuckyInfectionMaxTargets,
			DurationSec:     int(LuckyInfectionDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyInfectionFish, p.DisplayName, len(initialTargets), map[string]string{
		"message": fmt.Sprintf("%s 觸發了感染魚！%d 個目標被感染，感染正在蔓延...", p.DisplayName, len(initialTargets)),
		"color":   "#00FF88",
	})
	g.broadcastAnnouncement(ann)

	log.Printf("[LuckyInfection] player=%s triggered infection, initial=%d targets", p.ID, len(initialTargets))

	// 啟動感染蔓延 goroutine
	go g.runLuckyInfectionSpread(sessionID)
}

// runLuckyInfectionSpread 感染蔓延 goroutine
func (g *Game) runLuckyInfectionSpread(sessionID string) {
	ticker := time.NewTicker(LuckyInfectionSpreadCD)
	defer ticker.Stop()

	endTimer := time.NewTimer(LuckyInfectionDuration)
	defer endTimer.Stop()

	for {
		select {
		case <-ticker.C:
			// 嘗試蔓延感染
			g.doInfectionSpread(sessionID)

		case <-endTimer.C:
			// 感染時間到，觸發感染爆發
			g.doLuckyInfectionBlast(sessionID)
			return

		case <-g.stopCh:
			return
		}
	}
}

// doInfectionSpread 執行感染蔓延
func (g *Game) doInfectionSpread(sessionID string) {
	mgr := g.LuckyInfectionFish
	mgr.mu.Lock()

	// 確認 session 仍有效
	if !mgr.active || mgr.sessionID != sessionID {
		mgr.mu.Unlock()
		return
	}

	// 已達最大感染數或最大層數，不再蔓延
	if len(mgr.infectedTargets) >= LuckyInfectionMaxTargets || mgr.spreadLayer >= LuckyInfectionMaxLayers {
		mgr.mu.Unlock()
		return
	}

	// 收集當前感染目標的位置
	type infectedPos struct {
		instanceID string
		x, y       float64
		layer      int
	}
	var infectedList []infectedPos

	g.mu.RLock()
	for iid, entry := range mgr.infectedTargets {
		t := g.Targets[iid]
		if t != nil && t.HP > 0 {
			infectedList = append(infectedList, infectedPos{
				instanceID: iid,
				x:          t.X,
				y:          t.Y,
				layer:      entry.layer,
			})
		}
	}
	g.mu.RUnlock()

	if len(infectedList) == 0 {
		mgr.mu.Unlock()
		return
	}

	// 找出可被感染的相鄰目標
	var newInfected []ws.InfectionTargetInfo
	g.mu.RLock()
	for _, inf := range infectedList {
		if inf.layer >= LuckyInfectionMaxLayers {
			continue
		}
		for iid, t := range g.Targets {
			if t.HP <= 0 || isLuckyInfectionFish(t.DefID) {
				continue
			}
			if _, alreadyInfected := mgr.infectedTargets[iid]; alreadyInfected {
				continue
			}
			// 計算距離
			dx := t.X - inf.x
			dy := t.Y - inf.y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= LuckyInfectionSpreadRange {
				// 感染！
				newLayer := inf.layer + 1
				mgr.infectedTargets[iid] = &infectionEntry{
					instanceID: iid,
					layer:      newLayer,
					infectedAt: time.Now(),
				}
				newInfected = append(newInfected, ws.InfectionTargetInfo{
					InstanceID: iid,
					Layer:      newLayer,
					X:          t.X,
					Y:          t.Y,
				})
				if len(mgr.infectedTargets) >= LuckyInfectionMaxTargets {
					break
				}
			}
		}
		if len(mgr.infectedTargets) >= LuckyInfectionMaxTargets {
			break
		}
	}
	g.mu.RUnlock()

	if len(newInfected) > 0 {
		mgr.spreadLayer++
	}
	totalInfected := len(mgr.infectedTargets)
	mgr.mu.Unlock()

	if len(newInfected) == 0 {
		return
	}

	// 全服廣播：感染蔓延
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyInfectionFish,
		Payload: ws.LuckyInfectionFishPayload{
			Event:           "infection_spread",
			SessionID:       sessionID,
			InfectedTargets: newInfected,
			TotalInfected:   totalInfected,
			MaxInfected:     LuckyInfectionMaxTargets,
		},
	})

	log.Printf("[LuckyInfection] spread: +%d new infected, total=%d", len(newInfected), totalInfected)
}

// doLuckyInfectionBlast 感染爆發（12 秒後觸發）
func (g *Game) doLuckyInfectionBlast(sessionID string) {
	mgr := g.LuckyInfectionFish
	mgr.mu.Lock()

	if !mgr.active || mgr.sessionID != sessionID {
		mgr.mu.Unlock()
		return
	}

	// 收集所有感染目標
	var blastTargets []string
	for iid := range mgr.infectedTargets {
		blastTargets = append(blastTargets, iid)
	}

	// 重置狀態
	mgr.active = false
	mgr.infectedTargets = make(map[string]*infectionEntry)
	mgr.mu.Unlock()

	if len(blastTargets) == 0 {
		return
	}

	// 執行感染爆發
	killed := 0
	totalReward := 0
	var blastResults []ws.InfectionBlastResult

	// 取得平均 betCost 用於獎勵計算
	avgBet := g.getAvgBetCost()

	g.mu.Lock()
	for _, iid := range blastTargets {
		t := g.Targets[iid]
		if t == nil || t.HP <= 0 {
			continue
		}
		if rand.Float64() < LuckyInfectionBlastChance {
			// 擊破
			t.HP = 0
			killed++
			reward := int(t.Multiplier * float64(avgBet) * LuckyInfectionBlastMult)
			totalReward += reward
			blastResults = append(blastResults, ws.InfectionBlastResult{
				InstanceID: iid,
				Killed:     true,
				Reward:     reward,
			})
			delete(g.Targets, iid)
		} else {
			blastResults = append(blastResults, ws.InfectionBlastResult{
				InstanceID: iid,
				Killed:     false,
			})
		}
	}
	g.mu.Unlock()

	// 全服共享獎勵
	if totalReward > 0 {
		g.mu.Lock()
		for _, p := range g.Players {
			share := totalReward / len(g.Players)
			if share > 0 {
				p.AddCoins(share)
			}
		}
		g.mu.Unlock()
	}

	// 全服廣播：感染爆發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyInfectionFish,
		Payload: ws.LuckyInfectionFishPayload{
			Event:        "infection_blast",
			SessionID:    sessionID,
			BlastResults: blastResults,
			TotalKilled:  killed,
			TotalReward:  totalReward,
		},
	})

	// 全服公告（≥3 個擊破才公告）
	if killed >= 3 {
		color := "#00FF88"
		if killed >= 6 {
			color = "#00CCFF"
		}
		ann := g.Announce.Create(announce.EventLuckyInfectionFish, "", killed, map[string]string{
			"message": fmt.Sprintf("感染爆發！%d 個目標被感染消滅，全服共享 %d 金幣！", killed, totalReward),
			"color":   color,
		})
		g.broadcastAnnouncement(ann)
	}

	log.Printf("[LuckyInfection] blast: killed=%d, totalReward=%d", killed, totalReward)
}

// notifyLuckyInfectionFishKill 玩家擊破感染目標時的處理（供 handleKill 使用）
func (g *Game) notifyLuckyInfectionFishKill(instanceID string) {
	g.removeInfectionEntry(instanceID)

	mgr := g.LuckyInfectionFish
	mgr.mu.Lock()
	remaining := len(mgr.infectedTargets)
	sessionID := mgr.sessionID
	mgr.mu.Unlock()

	// 廣播感染目標被擊破
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyInfectionFish,
		Payload: ws.LuckyInfectionFishPayload{
			Event:         "infection_kill",
			SessionID:     sessionID,
			KilledTarget:  instanceID,
			TotalInfected: remaining,
		},
	})
}

// getAvgBetCostForInfection 取得平均投注成本（感染爆發獎勵計算用）
// 注意：此函數在 g.mu 外呼叫，需要自行加鎖
func (g *Game) getAvgBetCostForInfection() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if len(g.Players) == 0 {
		return 1
	}
	total := 0
	for _, p := range g.Players {
		betDef := data.GetBetDef(p.BetLevel)
		if betDef != nil {
			total += betDef.BetCost
		} else {
			total += 1
		}
	}
	return total / len(g.Players)
}
