// lucky_mirror_fish_handler.go — 幸運鏡像魚系統（DAY-215）
// 業界原創「鏡像複製」機制
//
// 設計：擊破 T173 後觸發「鏡像複製」：
//   - 在場上隨機選最多 3 個目標，為每個目標建立「鏡像分身」
//   - 鏡像分身 HP = 原目標 HP × 50%，倍率 = 原目標倍率 × 1.5
//   - 鏡像分身持續 8 秒；擊破鏡像分身獲得 ×1.5 倍率加成（乘法）
//   - 8 秒後所有未被擊破的鏡像分身「鏡像爆炸」（60% 擊破機率，0.60x 倍率）
//   - 個人冷卻 20 秒；全服廣播鏡像建立/爆炸
//
// 設計差異：
//   - 與彩虹稜鏡魚（DAY-213，染色 5 個目標，顏色對應不同倍率）不同，
//     幸運鏡像魚是「複製分身」，讓玩家有「要先打分身還是本體」的策略選擇
//   - 分身 HP 只有 50%，更容易擊破，讓玩家有「快速連殺」的爽感
//   - 分身倍率 ×1.5，讓玩家有「打分身比打本體更划算」的感覺
//   - 8 秒後自動爆炸，製造「等待→爆發」的高潮感
//   - 全服廣播讓所有玩家都看到鏡像，製造「全服競爭搶打分身」的社交感
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
	LuckyMirrorMaxTargets   = 3                    // 最多複製 3 個目標
	LuckyMirrorHPRatio      = 0.50                 // 鏡像 HP = 原 HP × 50%
	LuckyMirrorMultRatio    = 1.5                  // 鏡像倍率 = 原倍率 × 1.5
	LuckyMirrorDuration     = 8 * time.Second      // 鏡像持續時間
	LuckyMirrorBlastChance  = 0.60                 // 鏡像爆炸擊破機率
	LuckyMirrorBlastMult    = 0.60                 // 鏡像爆炸倍率
	LuckyMirrorPersonalCD   = 20 * time.Second     // 個人冷卻
)

// mirrorEntry 鏡像分身記錄
type mirrorEntry struct {
	mirrorID      string
	originalID    string
	originalDefID string
	x, y          float64
	mirrorHP      int
	mirrorMult    float64
	killed        bool
}

// luckyMirrorFishManager 幸運鏡像魚管理器
type luckyMirrorFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldown map[string]time.Time

	// 當前活躍鏡像分身（mirrorID → mirrorEntry）
	activeMirrors map[string]*mirrorEntry
}

func newLuckyMirrorFishManager() *luckyMirrorFishManager {
	return &luckyMirrorFishManager{
		personalCooldown: make(map[string]time.Time),
		activeMirrors:    make(map[string]*mirrorEntry),
	}
}

// isLuckyMirrorFish 判斷是否為幸運鏡像魚
func isLuckyMirrorFish(defID string) bool {
	return defID == "T173"
}

// getLuckyMirrorMultiplier 取得鏡像分身倍率加成（供 handleKill 使用）
// 若被擊破的目標是鏡像分身，回傳 ×1.5 乘法加成
func (g *Game) getLuckyMirrorMultiplier(instanceID string) float64 {
	mgr := g.LuckyMirrorFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if entry, ok := mgr.activeMirrors[instanceID]; ok && !entry.killed {
		return LuckyMirrorMultRatio
	}
	return 1.0
}

// removeLuckyMirrorEntry 鏡像分身被擊破後移除（供 handleKill 使用）
func (g *Game) removeLuckyMirrorEntry(instanceID string) {
	mgr := g.LuckyMirrorFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if entry, ok := mgr.activeMirrors[instanceID]; ok {
		entry.killed = true
		// 廣播：鏡像分身被擊破
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyMirrorFish,
			Payload: ws.LuckyMirrorFishPayload{
				Event:    "mirror_kill",
				MirrorID: instanceID,
			},
		})
		log.Printf("[LuckyMirrorFish] mirror killed: mirrorID=%s", instanceID)
	}
}

// isLuckyMirrorEntry 判斷 instanceID 是否為鏡像分身
func (g *Game) isLuckyMirrorEntry(instanceID string) bool {
	mgr := g.LuckyMirrorFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	_, ok := mgr.activeMirrors[instanceID]
	return ok
}

// tryLuckyMirrorFish 擊破 T173 後觸發鏡像複製
func (g *Game) tryLuckyMirrorFish(p *player.Player) {
	mgr := g.LuckyMirrorFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldown[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}
	mgr.personalCooldown[p.ID] = time.Now().Add(LuckyMirrorPersonalCD)
	mgr.mu.Unlock()

	// 選取場上最多 3 個目標（排除 BOSS 和鏡像分身本身）
	g.mu.RLock()
	candidates := make([]*target.Target, 0, 10)
	for _, t := range g.Targets {
		if t.HP > 0 && t.DefID != "T173" && t.DefID != "B001" {
			// 排除已是鏡像分身的目標
			if !g.isLuckyMirrorEntry(t.InstanceID) {
				candidates = append(candidates, t)
			}
		}
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		log.Printf("[LuckyMirrorFish] no candidates for mirror")
		return
	}

	// 隨機打亂，取最多 3 個
	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	if len(candidates) > LuckyMirrorMaxTargets {
		candidates = candidates[:LuckyMirrorMaxTargets]
	}

	// 建立鏡像分身
	mirrors := make([]ws.LuckyMirrorFishInfo, 0, len(candidates))
	entries := make([]*mirrorEntry, 0, len(candidates))

	for _, orig := range candidates {
		mirrorID := "mirror_" + orig.InstanceID
		mirrorHP := int(float64(orig.HP) * LuckyMirrorHPRatio)
		if mirrorHP < 1 {
			mirrorHP = 1
		}
		mirrorMult := orig.Multiplier * LuckyMirrorMultRatio

		// 鏡像位置：在原目標附近偏移（±80px）
		offsetX := (rand.Float64()*2 - 1) * 80
		offsetY := (rand.Float64()*2 - 1) * 80
		mx := clampFloat(orig.X+offsetX, 50, 1230)
		my := clampFloat(orig.Y+offsetY, 50, 670)

		entry := &mirrorEntry{
			mirrorID:      mirrorID,
			originalID:    orig.InstanceID,
			originalDefID: orig.DefID,
			x:             mx,
			y:             my,
			mirrorHP:      mirrorHP,
			mirrorMult:    mirrorMult,
			killed:        false,
		}
		entries = append(entries, entry)

		mirrors = append(mirrors, ws.LuckyMirrorFishInfo{
			MirrorID:      mirrorID,
			OriginalID:    orig.InstanceID,
			OriginalDefID: orig.DefID,
			X:             mx,
			Y:             my,
			MirrorHP:      mirrorHP,
			MirrorMult:    mirrorMult,
		})
	}

	// 儲存鏡像分身
	mgr.mu.Lock()
	for _, entry := range entries {
		mgr.activeMirrors[entry.mirrorID] = entry
	}
	mgr.mu.Unlock()

	log.Printf("[LuckyMirrorFish] player=%s triggered mirror: %d mirrors created", p.ID, len(mirrors))

	// 全服廣播：鏡像複製開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMirrorFish,
		Payload: ws.LuckyMirrorFishPayload{
			Event:      "mirror_start",
			PlayerName: p.DisplayName,
			Mirrors:    mirrors,
			MultBoost:  LuckyMirrorMultRatio,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyMirrorFish, p.DisplayName, len(mirrors), map[string]string{
		"message": fmt.Sprintf("🪞 %s 觸發幸運鏡像魚！%d 個目標出現鏡像分身（×%.1f 倍率）！",
			p.DisplayName, len(mirrors), LuckyMirrorMultRatio),
		"color": "#00FFFF",
	})
	g.broadcastAnnouncement(ann)

	// 8 秒後觸發鏡像爆炸
	go g.runLuckyMirrorBlast(entries)
}

// runLuckyMirrorBlast 8 秒後對未被擊破的鏡像分身觸發爆炸
func (g *Game) runLuckyMirrorBlast(entries []*mirrorEntry) {
	time.Sleep(LuckyMirrorDuration)

	mgr := g.LuckyMirrorFish
	mgr.mu.Lock()

	// 找出未被擊破的鏡像分身
	surviving := make([]*mirrorEntry, 0)
	for _, entry := range entries {
		if !entry.killed {
			surviving = append(surviving, entry)
		}
	}

	// 清除所有鏡像分身記錄
	for _, entry := range entries {
		delete(mgr.activeMirrors, entry.mirrorID)
	}
	mgr.mu.Unlock()

	if len(surviving) == 0 {
		log.Printf("[LuckyMirrorFish] all mirrors killed before blast")
		return
	}

	log.Printf("[LuckyMirrorFish] blast: %d surviving mirrors", len(surviving))

	// 對每個存活的鏡像分身執行爆炸
	blastCount := 0
	totalReward := 0

	for _, entry := range surviving {
		if rand.Float64() < LuckyMirrorBlastChance {
			// 爆炸擊破：找場上最近的目標
			g.mu.RLock()
			var nearestTarget *target.Target
			minDist := 999999.0
			for _, t := range g.Targets {
				if t.HP > 0 && t.DefID != "B001" {
					dx := t.X - entry.x
					dy := t.Y - entry.y
					dist := dx*dx + dy*dy
					if dist < minDist {
						minDist = dist
						nearestTarget = t
					}
				}
			}
			g.mu.RUnlock()

			if nearestTarget != nil && minDist < 200*200 { // 200px 範圍內
				reward := int(float64(nearestTarget.Multiplier) * LuckyMirrorBlastMult * float64(g.getAvgBetCost()))
				if reward > 0 {
					g.distributeRewardToAll(reward)
					totalReward += reward
					blastCount++
				}
			}
		}
	}

	// 廣播：鏡像爆炸
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMirrorFish,
		Payload: ws.LuckyMirrorFishPayload{
			Event:       "mirror_blast",
			BlastCount:  blastCount,
			TotalReward: totalReward,
		},
	})

	// 廣播：鏡像結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMirrorFish,
		Payload: ws.LuckyMirrorFishPayload{
			Event:       "mirror_result",
			KilledCount: len(entries) - len(surviving),
			BlastCount:  blastCount,
			TotalReward: totalReward,
		},
	})

	// 全服公告（≥2 個爆炸時）
	if blastCount >= 2 {
		color := "#00FFFF"
		if blastCount >= 3 {
			color = "#00FF88"
		}
		ann := g.Announce.Create(announce.EventLuckyMirrorFish, "", blastCount, map[string]string{
			"message": fmt.Sprintf("🪞💥 鏡像爆炸！%d 個鏡像分身爆炸！全服共享獎勵！", blastCount),
			"color":   color,
		})
		g.broadcastAnnouncement(ann)
	}

	log.Printf("[LuckyMirrorFish] blast done: blastCount=%d totalReward=%d", blastCount, totalReward)
}

// getAvgBetCost 取得場上玩家的平均投注成本（供獎勵計算使用）
func (g *Game) getAvgBetCost() int {
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
