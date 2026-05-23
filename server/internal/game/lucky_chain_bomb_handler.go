// lucky_chain_bomb_handler.go — 幸運鏈鎖爆炸魚系統（DAY-226）
// 業界原創「連鎖爆炸」機制
//
// 設計：擊破 T184 後，場上隨機 3 個目標被「引爆標記」：
//   - 引爆標記目標被擊破後，立即引爆周圍 200px 內所有目標（60% 擊破機率，×1.5 倍率）
//   - 被引爆的目標如果也有引爆標記，繼續連鎖（最多 3 層連鎖）
//   - 連鎖爆炸獎勵給觸發者（擊破引爆標記目標的玩家）
//   - 個人冷卻 20 秒；全服廣播引爆標記/連鎖爆炸
//
// 設計差異：
//   - 與感染魚（DAY-219，動態蔓延，時間驅動）不同，鏈鎖爆炸是「空間爆炸」，即時觸發
//   - 與分裂魚（DAY-224，一魚分三）不同，鏈鎖爆炸是「引爆周圍目標」，空間策略感更強
//   - 「連鎖最多 3 層」讓玩家有「要把魚聚集在一起才能最大化連鎖」的策略感
//   - 「引爆標記」讓玩家有「要先打哪個引爆目標」的選擇感
//   - 全服廣播讓所有玩家看到連鎖爆炸的壯觀效果，製造「全服一起爽」的社交感
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyChainBombPersonalCD   = 20 * time.Second // 個人冷卻
	LuckyChainBombMarkDuration = 15 * time.Second // 引爆標記持續時間
	LuckyChainBombRadius       = 200.0            // 連鎖爆炸半徑（px）
	LuckyChainBombKillChance   = 0.60             // 連鎖爆炸擊破機率
	LuckyChainBombMult         = 1.5              // 連鎖爆炸倍率加成（乘法）
	LuckyChainBombMaxChain     = 3                // 最大連鎖層數
	LuckyChainBombInitialMarks = 3                // 初始引爆標記數量
)

// chainBombMark 引爆標記條目（避免與 chain_bomb_handler.go 的 chainBombEntry 衝突）
type chainBombMark struct {
	instanceID string
	expiresAt  time.Time
}

// luckyChainBombManager 幸運鏈鎖爆炸魚管理器
type luckyChainBombManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldown map[string]time.Time

	// 引爆標記（targetInstanceID → mark）
	markedTargets map[string]*chainBombMark

	// 當前 instance ID（用於廣播識別）
	currentInstanceID string
}

func newLuckyChainBombManager() *luckyChainBombManager {
	return &luckyChainBombManager{
		personalCooldown: make(map[string]time.Time),
		markedTargets:    make(map[string]*chainBombMark),
	}
}

// isLuckyChainBombFish 判斷是否為幸運鏈鎖爆炸魚
func isLuckyChainBombFish(defID string) bool {
	return defID == "T184"
}

// isChainBombMarked 判斷目標是否有引爆標記（供 handleKill 使用）
func (g *Game) isChainBombMarked(instanceID string) bool {
	mgr := g.LuckyChainBomb
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	mark, ok := mgr.markedTargets[instanceID]
	if !ok {
		return false
	}
	if time.Now().After(mark.expiresAt) {
		delete(mgr.markedTargets, instanceID)
		return false
	}
	return true
}

// removeChainBombMark 移除引爆標記（目標被擊破後）
func (g *Game) removeChainBombMark(instanceID string) {
	mgr := g.LuckyChainBomb
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	delete(mgr.markedTargets, instanceID)
}

// getLuckyChainBombMult 取得連鎖爆炸倍率（供 handleKill 使用）
// 如果目標有引爆標記，回傳 ×1.5 乘法加成
func (g *Game) getLuckyChainBombMult(instanceID string) float64 {
	if g.isChainBombMarked(instanceID) {
		return LuckyChainBombMult
	}
	return 1.0
}

// tryLuckyChainBombFish 擊破 T184 後觸發引爆標記（供 handleKill 使用）
func (g *Game) tryLuckyChainBombFish(p *player.Player) {
	mgr := g.LuckyChainBomb
	mgr.mu.Lock()

	// 個人冷卻檢查
	if until, ok := mgr.personalCooldown[p.ID]; ok && time.Now().Before(until) {
		mgr.mu.Unlock()
		return
	}

	// 設定個人冷卻
	mgr.personalCooldown[p.ID] = time.Now().Add(LuckyChainBombPersonalCD)

	// 建立 instance ID
	instanceID := fmt.Sprintf("cbomb_%d", time.Now().UnixNano())
	mgr.currentInstanceID = instanceID
	mgr.mu.Unlock()

	// 取得場上目標（排除 BOSS）
	g.mu.RLock()
	candidates := make([]*target.Target, 0)
	for _, t := range g.Targets {
		if t.IsAlive && t.Def.Type != "boss" {
			candidates = append(candidates, t)
		}
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		return
	}

	// 隨機選取最多 3 個目標加上引爆標記
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	count := LuckyChainBombInitialMarks
	if len(candidates) < count {
		count = len(candidates)
	}

	type markedInfo struct {
		TargetID string  `json:"target_id"`
		X        float64 `json:"x"`
		Y        float64 `json:"y"`
	}
	markedList := make([]markedInfo, 0, count)

	expiresAt := time.Now().Add(LuckyChainBombMarkDuration)

	mgr.mu.Lock()
	for i := 0; i < count; i++ {
		t := candidates[i]
		mgr.markedTargets[t.InstanceID] = &chainBombMark{
			instanceID: instanceID,
			expiresAt:  expiresAt,
		}
		markedList = append(markedList, markedInfo{
			TargetID: t.InstanceID,
			X:        t.X,
			Y:        t.Y,
		})
	}
	mgr.mu.Unlock()

	log.Printf("[LuckyChainBomb] player=%s triggered, marked %d targets instance=%s",
		p.ID, count, instanceID)

	// 全服廣播引爆標記開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyChainBomb,
		Payload: ws.LuckyChainBombPayload{
			Event:       "chain_bomb_start",
			InstanceID:  instanceID,
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			Marked:      markedList,
			DurationSec: int(LuckyChainBombMarkDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyChainBombFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("💣 %s 觸發鏈鎖爆炸！%d 個目標被引爆標記！",
			p.DisplayName, count),
		"color": "#FF4500",
	})
	g.broadcastAnnouncement(ann)

	// 15 秒後清除過期標記
	go func() {
		time.Sleep(LuckyChainBombMarkDuration)
		mgr.mu.Lock()
		for _, info := range markedList {
			if mark, ok := mgr.markedTargets[info.TargetID]; ok && mark.instanceID == instanceID {
				delete(mgr.markedTargets, info.TargetID)
			}
		}
		mgr.mu.Unlock()

		// 廣播標記過期
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyChainBomb,
			Payload: ws.LuckyChainBombPayload{
				Event:      "chain_bomb_expire",
				InstanceID: instanceID,
			},
		})
	}()
}

// ChainBombBlastResult 連鎖爆炸單個目標結果
type ChainBombBlastResult struct {
	TargetID string  `json:"target_id"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Killed   bool    `json:"killed"`
	Reward   int     `json:"reward"`
	IsChain  bool    `json:"is_chain"`
}

// doChainBombExplosion 執行連鎖爆炸（供 notifyChainBombKill 使用）
// 當引爆標記目標被擊破時，引爆周圍 200px 內所有目標
func (g *Game) doChainBombExplosion(p *player.Player, triggerInstanceID string, triggerX, triggerY float64, chainLayer int) {
	if chainLayer > LuckyChainBombMaxChain {
		return
	}

	mgr := g.LuckyChainBomb
	mgr.mu.Lock()
	instanceID := mgr.currentInstanceID
	mgr.mu.Unlock()

	// 取得周圍目標
	g.mu.RLock()
	nearby := make([]*target.Target, 0)
	for _, t := range g.Targets {
		if t.InstanceID == triggerInstanceID || !t.IsAlive {
			continue
		}
		if t.Def.Type == "boss" {
			continue
		}
		dx := t.X - triggerX
		dy := t.Y - triggerY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= LuckyChainBombRadius {
			nearby = append(nearby, t)
		}
	}
	g.mu.RUnlock()

	if len(nearby) == 0 {
		return
	}

	results := make([]ChainBombBlastResult, 0, len(nearby))
	totalReward := 0

	type chainTarget struct {
		instanceID string
		x, y       float64
	}
	chainTargets := make([]chainTarget, 0)

	for _, t := range nearby {
		// 60% 擊破機率
		if rand.Float64() >= LuckyChainBombKillChance {
			results = append(results, ChainBombBlastResult{
				TargetID: t.InstanceID,
				X:        t.X,
				Y:        t.Y,
				Killed:   false,
			})
			continue
		}

		// 計算獎勵（×1.5 倍率）
		betCost := g.getAvgBetCost()
		mult := float64(t.Def.MultiplierMin+t.Def.MultiplierMax) / 2.0
		reward := int(float64(betCost) * mult * LuckyChainBombMult)
		totalReward += reward

		// 判斷是否有引爆標記（連鎖）
		isChain := g.isChainBombMarked(t.InstanceID)
		if isChain {
			chainTargets = append(chainTargets, chainTarget{
				instanceID: t.InstanceID,
				x:          t.X,
				y:          t.Y,
			})
			g.removeChainBombMark(t.InstanceID)
		}

		// 擊破目標
		g.mu.Lock()
		if tgt, ok := g.Targets[t.InstanceID]; ok && tgt.IsAlive {
			tgt.IsAlive = false
			tgt.HP = 0
			delete(g.Targets, t.InstanceID)
		}
		g.mu.Unlock()

		// 給玩家獎勵
		p.AddCoins(reward)

		results = append(results, ChainBombBlastResult{
			TargetID: t.InstanceID,
			X:        t.X,
			Y:        t.Y,
			Killed:   true,
			Reward:   reward,
			IsChain:  isChain,
		})

		log.Printf("[LuckyChainBomb] chain layer=%d killed target=%s reward=%d isChain=%v",
			chainLayer, t.InstanceID, reward, isChain)
	}

	// 廣播爆炸結果
	if len(results) > 0 {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyChainBomb,
			Payload: ws.LuckyChainBombPayload{
				Event:       "chain_bomb_blast",
				InstanceID:  instanceID,
				PlayerID:    p.ID,
				PlayerName:  p.DisplayName,
				TriggerID:   triggerInstanceID,
				TriggerX:    triggerX,
				TriggerY:    triggerY,
				ChainLayer:  chainLayer,
				BlastRadius: LuckyChainBombRadius,
				Results:     results,
				TotalReward: totalReward,
			},
		})
	}

	// 全服公告（連鎖層數 ≥ 2 時）
	if chainLayer >= 2 && totalReward > 0 {
		ann := g.Announce.Create(announce.EventLuckyChainBombFish, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("💥 %s 連鎖爆炸第 %d 層！獲得 🪙%d！",
				p.DisplayName, chainLayer, totalReward),
			"color": "#FF6B35",
		})
		g.broadcastAnnouncement(ann)
	}

	// 遞迴觸發連鎖（有引爆標記的目標繼續爆炸）
	for _, ct := range chainTargets {
		ctCopy := ct // 避免閉包問題
		go g.doChainBombExplosion(p, ctCopy.instanceID, ctCopy.x, ctCopy.y, chainLayer+1)
	}
}

// notifyChainBombKill 引爆標記目標被玩家擊破時觸發連鎖（供 handleKill 使用）
func (g *Game) notifyChainBombKill(p *player.Player, t *target.Target) {
	// 移除引爆標記
	g.removeChainBombMark(t.InstanceID)

	// 廣播引爆標記被擊破
	mgr := g.LuckyChainBomb
	mgr.mu.Lock()
	instanceID := mgr.currentInstanceID
	mgr.mu.Unlock()

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyChainBomb,
		Payload: ws.LuckyChainBombPayload{
			Event:      "chain_bomb_trigger",
			InstanceID: instanceID,
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			TriggerID:  t.InstanceID,
			TriggerX:   t.X,
			TriggerY:   t.Y,
			ChainLayer: 1,
		},
	})

	// 執行連鎖爆炸（第 1 層）
	go g.doChainBombExplosion(p, t.InstanceID, t.X, t.Y, 1)
}
