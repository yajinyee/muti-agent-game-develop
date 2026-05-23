// lucky_parasite_fish_handler.go — 幸運寄生魚系統（DAY-229）
// 業界原創「寄生附著+跳躍」機制
//
// 設計：擊破 T187 後觸發「寄生釋放」：
//   - 場上隨機 3 個目標被「寄生蟲附著」（綠色標記）
//   - 寄生目標每 2 秒自動損失 HP（-8%/次，最多 5 次）
//   - 寄生目標被擊破時，寄生蟲「跳躍」到最近的目標繼續寄生（最多跳躍 2 次）
//   - 玩家擊破寄生目標獲得 ×2.2 倍率加成（乘法）
//   - 個人冷卻 22 秒；全服廣播寄生附著/跳躍/消散
//
// 設計差異：
//   - 與感染魚（DAY-219，動態蔓延，時間驅動）不同，寄生魚是「附著+跳躍」，讓玩家有「要趁寄生蟲跳走前打死它」的緊迫感
//   - 與黑洞魚（DAY-221，重力傷害）不同，寄生魚是「個別目標 HP 損失」，更精準
//   - 「自動 HP 損失」讓玩家感受到「寄生蟲在幫我削血」的輔助感
//   - 「跳躍機制」讓寄生蟲像真正的寄生蟲一樣移動，視覺上更有趣
//   - 「最多跳躍 2 次」確保 RTP 平衡，不會無限跳躍
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyParasitePersonalCD  = 22 * time.Second // 個人冷卻
	LuckyParasiteTickInterval = 2 * time.Second  // 寄生 HP 損失間隔
	LuckyParasiteMaxTicks    = 5                 // 最多 HP 損失次數
	LuckyParasiteHPLoss      = 0.08              // 每次 HP 損失比例（8%）
	LuckyParasiteKillMult    = 2.2               // 擊破寄生目標倍率加成
	LuckyParasiteMaxJumps    = 2                 // 最多跳躍次數
	LuckyParasiteMaxTargets  = 3                 // 初始寄生目標數
	LuckyParasiteJumpRange   = 300.0             // 跳躍範圍（px）
)

// parasiteEntry 寄生蟲記錄
type parasiteEntry struct {
	instanceID string
	jumpLayer  int       // 已跳躍次數（0=初始，1=第一跳，2=第二跳）
	tickCount  int       // 已 tick 次數
	expiresAt  time.Time // 最終過期時間
}

// luckyParasiteFishManager 幸運寄生魚管理器
type luckyParasiteFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldown map[string]time.Time

	// 寄生目標（targetID → parasiteEntry）
	parasiteTargets map[string]*parasiteEntry

	// 當前 instance ID
	currentInstanceID string
}

func newLuckyParasiteFishManager() *luckyParasiteFishManager {
	return &luckyParasiteFishManager{
		personalCooldown: make(map[string]time.Time),
		parasiteTargets:  make(map[string]*parasiteEntry),
	}
}

// isLuckyParasiteFish 判斷是否為幸運寄生魚
func isLuckyParasiteFish(defID string) bool {
	return defID == "T187"
}

// isParasiteTarget 判斷目標是否被寄生（供 handleKill 使用）
func (g *Game) isParasiteTarget(targetID string) bool {
	mgr := g.LuckyParasiteFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	entry, ok := mgr.parasiteTargets[targetID]
	if !ok {
		return false
	}
	return time.Now().Before(entry.expiresAt)
}

// getLuckyParasiteKillMult 取得寄生目標擊破倍率（供 handleKill 使用）
func (g *Game) getLuckyParasiteKillMult(targetID string) float64 {
	if g.isParasiteTarget(targetID) {
		return LuckyParasiteKillMult
	}
	return 1.0
}

// notifyParasiteKill 寄生目標被玩家擊破時處理跳躍（供 handleKill 使用）
func (g *Game) notifyParasiteKill(p *player.Player, killedTargetID string, reward int) {
	mgr := g.LuckyParasiteFish
	mgr.mu.Lock()
	entry, ok := mgr.parasiteTargets[killedTargetID]
	if !ok {
		mgr.mu.Unlock()
		return
	}
	instanceID := entry.instanceID
	jumpLayer := entry.jumpLayer
	delete(mgr.parasiteTargets, killedTargetID)
	mgr.mu.Unlock()

	// 廣播寄生目標被擊破
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyParasiteFish,
		Payload: ws.LuckyParasiteFishPayload{
			Event:      "parasite_kill",
			InstanceID: instanceID,
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			TargetID:   killedTargetID,
			KillMult:   LuckyParasiteKillMult,
			KillReward: reward,
			JumpLayer:  jumpLayer,
		},
	})

	log.Printf("[LuckyParasite] player=%s killed parasite target=%s layer=%d reward=%d",
		p.ID, killedTargetID, jumpLayer, reward)

	// 嘗試跳躍到最近的目標
	if jumpLayer < LuckyParasiteMaxJumps {
		go g.doParasiteJump(instanceID, killedTargetID, jumpLayer+1)
	}
}

// removeParasiteEntry 移除寄生記錄（目標逃跑或過期時呼叫）
func (g *Game) removeParasiteEntry(targetID string) {
	mgr := g.LuckyParasiteFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	delete(mgr.parasiteTargets, targetID)
}

// doParasiteJump 寄生蟲跳躍到最近的目標
func (g *Game) doParasiteJump(instanceID, fromTargetID string, newJumpLayer int) {
	// 找到被擊破目標的位置
	g.mu.RLock()
	var fromX, fromY float64
	// 目標已被擊破，用場上最近的存活目標作為跳躍起點
	// 找場上最近的未寄生存活目標
	type candidate struct {
		id   string
		dist float64
	}
	candidates := make([]candidate, 0)
	mgr := g.LuckyParasiteFish

	for id, t := range g.Targets {
		if !t.IsAlive || t.Def.Type == "boss" {
			continue
		}
		mgr.mu.Lock()
		_, alreadyParasited := mgr.parasiteTargets[id]
		mgr.mu.Unlock()
		if alreadyParasited {
			continue
		}
		dx := t.X - fromX
		dy := t.Y - fromY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= LuckyParasiteJumpRange {
			candidates = append(candidates, candidate{id: id, dist: dist})
		}
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		log.Printf("[LuckyParasite] no jump target found for instance=%s layer=%d", instanceID, newJumpLayer)
		return
	}

	// 選最近的目標
	nearest := candidates[0]
	for _, c := range candidates[1:] {
		if c.dist < nearest.dist {
			nearest = c
		}
	}

	// 設定新的寄生目標
	expiresAt := time.Now().Add(LuckyParasiteTickInterval * time.Duration(LuckyParasiteMaxTicks+1))
	mgr.mu.Lock()
	mgr.parasiteTargets[nearest.id] = &parasiteEntry{
		instanceID: instanceID,
		jumpLayer:  newJumpLayer,
		tickCount:  0,
		expiresAt:  expiresAt,
	}
	mgr.mu.Unlock()

	log.Printf("[LuckyParasite] parasite jumped to target=%s layer=%d", nearest.id, newJumpLayer)

	// 廣播跳躍
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyParasiteFish,
		Payload: ws.LuckyParasiteFishPayload{
			Event:        "parasite_jump",
			InstanceID:   instanceID,
			FromTargetID: fromTargetID,
			ToTargetID:   nearest.id,
			JumpLayer:    newJumpLayer,
			KillMult:     LuckyParasiteKillMult,
		},
	})

	// 啟動新目標的 tick goroutine
	go g.runParasiteTick(instanceID, nearest.id, newJumpLayer)
}

// tryLuckyParasiteFish 擊破 T187 後觸發寄生釋放（供 handleKill 使用）
func (g *Game) tryLuckyParasiteFish(p *player.Player) {
	mgr := g.LuckyParasiteFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if until, ok := mgr.personalCooldown[p.ID]; ok && time.Now().Before(until) {
		mgr.mu.Unlock()
		return
	}

	// 設定個人冷卻
	mgr.personalCooldown[p.ID] = time.Now().Add(LuckyParasitePersonalCD)

	// 建立 instance ID
	instanceID := fmt.Sprintf("parasite_%d", time.Now().UnixNano())
	mgr.currentInstanceID = instanceID
	mgr.mu.Unlock()

	log.Printf("[LuckyParasite] player=%s triggered parasite release instance=%s", p.ID, instanceID)

	// 選取最多 3 個存活目標進行寄生
	g.mu.Lock()
	candidates := make([]string, 0, 8)
	for id, t := range g.Targets {
		if t.IsAlive && t.Def.Type != "boss" && t.DefID != "T187" {
			candidates = append(candidates, id)
		}
	}
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	if len(candidates) > LuckyParasiteMaxTargets {
		candidates = candidates[:LuckyParasiteMaxTargets]
	}
	g.mu.Unlock()

	if len(candidates) == 0 {
		log.Printf("[LuckyParasite] no candidates for parasite")
		return
	}

	// 設定寄生目標
	expiresAt := time.Now().Add(LuckyParasiteTickInterval * time.Duration(LuckyParasiteMaxTicks+1))
	mgr.mu.Lock()
	for _, id := range candidates {
		mgr.parasiteTargets[id] = &parasiteEntry{
			instanceID: instanceID,
			jumpLayer:  0,
			tickCount:  0,
			expiresAt:  expiresAt,
		}
	}
	mgr.mu.Unlock()

	log.Printf("[LuckyParasite] set %d targets to parasite state", len(candidates))

	// 全服廣播寄生釋放開始
	type targetInfo struct {
		ID string `json:"id"`
	}
	targetInfos := make([]targetInfo, len(candidates))
	for i, id := range candidates {
		targetInfos[i] = targetInfo{ID: id}
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyParasiteFish,
		Payload: ws.LuckyParasiteFishPayload{
			Event:         "parasite_start",
			InstanceID:    instanceID,
			PlayerID:      p.ID,
			PlayerName:    p.DisplayName,
			ParasiteCount: len(candidates),
			DurationSec:   int((LuckyParasiteTickInterval * time.Duration(LuckyParasiteMaxTicks)).Seconds()),
			KillMult:      LuckyParasiteKillMult,
			Targets:       targetInfos,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyParasiteFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🦠 %s 觸發寄生釋放！%d 個目標被寄生，擊破獲得 ×%.1f！",
			p.DisplayName, len(candidates), LuckyParasiteKillMult),
		"color": "#27AE60",
	})
	g.broadcastAnnouncement(ann)

	// 為每個寄生目標啟動 tick goroutine
	for _, id := range candidates {
		go g.runParasiteTick(instanceID, id, 0)
	}
}

// runParasiteTick 寄生 HP 損失 tick goroutine
func (g *Game) runParasiteTick(instanceID, targetID string, jumpLayer int) {
	for tick := 1; tick <= LuckyParasiteMaxTicks; tick++ {
		time.Sleep(LuckyParasiteTickInterval)

		// 確認目標仍在寄生狀態
		mgr := g.LuckyParasiteFish
		mgr.mu.Lock()
		entry, ok := mgr.parasiteTargets[targetID]
		if !ok || entry.instanceID != instanceID {
			mgr.mu.Unlock()
			return
		}
		entry.tickCount = tick
		mgr.mu.Unlock()

		// 對目標施加 HP 損失
		g.mu.Lock()
		t, exists := g.Targets[targetID]
		if !exists || !t.IsAlive {
			g.mu.Unlock()
			// 目標已死亡，清除寄生記錄
			mgr.mu.Lock()
			delete(mgr.parasiteTargets, targetID)
			mgr.mu.Unlock()
			return
		}
		hpLoss := int(float64(t.MaxHP) * LuckyParasiteHPLoss)
		if hpLoss < 1 {
			hpLoss = 1
		}
		t.HP -= hpLoss
		if t.HP < 1 {
			t.HP = 1
		}
		g.mu.Unlock()

		log.Printf("[LuckyParasite] tick=%d target=%s HP-=%d", tick, targetID, hpLoss)

		// 廣播 HP 損失
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyParasiteFish,
			Payload: ws.LuckyParasiteFishPayload{
				Event:      "parasite_tick",
				InstanceID: instanceID,
				TargetID:   targetID,
				HPLoss:     LuckyParasiteHPLoss,
				TickCount:  tick,
				JumpLayer:  jumpLayer,
			},
		})
	}

	// 最大 tick 達到，寄生消散
	mgr := g.LuckyParasiteFish
	mgr.mu.Lock()
	entry, ok := mgr.parasiteTargets[targetID]
	if ok && entry.instanceID == instanceID {
		delete(mgr.parasiteTargets, targetID)
	}
	mgr.mu.Unlock()

	log.Printf("[LuckyParasite] parasite expired on target=%s layer=%d", targetID, jumpLayer)

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyParasiteFish,
		Payload: ws.LuckyParasiteFishPayload{
			Event:      "parasite_end",
			InstanceID: instanceID,
			TargetID:   targetID,
			JumpLayer:  jumpLayer,
		},
	})
}
