// lucky_time_bomb_fish_handler.go — 幸運時間炸彈魚系統（DAY-235）
// 業界原創「倒數計時+提前引爆+連鎖爆炸」機制
//
// 設計：擊破 T193 後，場上隨機 4 個目標被「時間炸彈標記」（倒數 8 秒）：
//   - 倒數結束時自動爆炸（80% 擊破機率，×1.6 倍率，個人獎勵）
//   - 玩家可以「提前引爆」（射擊命中炸彈目標）：
//     立即爆炸（100% 擊破）+ 引爆周圍 150px 內目標（60% 機率，×1.2 倍率）
//   - 提前引爆的目標獲得 ×2.0 倍率加成（比等待爆炸更高）
//   - 個人冷卻 20 秒
//
// 設計差異：
//   - 與鏈鎖爆炸魚（DAY-226，引爆標記+連鎖）不同，時間炸彈是「倒數計時」，
//     讓玩家有「要不要等倒數還是提前引爆」的策略決策
//   - 「提前引爆」讓玩家有「主動控制爆炸時機」的掌控感
//   - 「連鎖爆炸」讓提前引爆有額外獎勵，鼓勵積極射擊
//   - 「倒數 8 秒」讓玩家有緊迫感，不能無限等待
//   - 全服廣播炸彈標記讓所有玩家都看到炸彈位置，製造「全服競爭引爆」的社交感
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
	LuckyTimeBombPersonalCD    = 20 * time.Second // 個人冷卻
	LuckyTimeBombFuseTime      = 8 * time.Second  // 炸彈倒數時間
	LuckyTimeBombCount         = 4                // 初始炸彈數量
	LuckyTimeBombAutoChance    = 0.80             // 自動爆炸擊破機率
	LuckyTimeBombAutoMult      = 1.6              // 自動爆炸倍率
	LuckyTimeBombEarlyMult     = 2.0              // 提前引爆倍率（更高）
	LuckyTimeBombChainRadius   = 150.0            // 連鎖爆炸範圍（px）
	LuckyTimeBombChainChance   = 0.60             // 連鎖爆炸擊破機率
	LuckyTimeBombChainMult     = 1.2              // 連鎖爆炸倍率
)

// timeBombEntry 單個炸彈標記
type timeBombEntry struct {
	instanceID string
	expiresAt  time.Time
}

// luckyTimeBombFishManager 幸運時間炸彈魚管理器
type luckyTimeBombFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 炸彈標記（targetID → timeBombEntry）
	bombTargets map[string]*timeBombEntry

	// 當前 instanceID（用於區分不同批次的炸彈）
	currentInstanceID string
}

func newLuckyTimeBombFishManager() *luckyTimeBombFishManager {
	return &luckyTimeBombFishManager{
		personalCooldowns: make(map[string]time.Time),
		bombTargets:       make(map[string]*timeBombEntry),
	}
}

// isLuckyTimeBombFish 判斷是否為幸運時間炸彈魚
func isLuckyTimeBombFish(defID string) bool {
	return defID == "T193"
}

// isTimeBombTarget 判斷目標是否有炸彈標記（供 handleKill 使用）
func (g *Game) isTimeBombTarget(targetID string) bool {
	mgr := g.LuckyTimeBombFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	entry, ok := mgr.bombTargets[targetID]
	if !ok {
		return false
	}
	if time.Now().After(entry.expiresAt) {
		delete(mgr.bombTargets, targetID)
		return false
	}
	return true
}

// getLuckyTimeBombEarlyMult 取得提前引爆倍率（供 handleKill 使用）
func (g *Game) getLuckyTimeBombEarlyMult(targetID string) float64 {
	if g.isTimeBombTarget(targetID) {
		return LuckyTimeBombEarlyMult
	}
	return 1.0
}

// removeTimeBombEntry 移除炸彈標記（目標被擊破後呼叫）
func (g *Game) removeTimeBombEntry(targetID string) {
	mgr := g.LuckyTimeBombFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	delete(mgr.bombTargets, targetID)
}

// notifyTimeBombKill 炸彈目標被玩家提前引爆（供 handleKill 使用）
func (g *Game) notifyTimeBombKill(p *player.Player, targetID string, targetX, targetY float64) {
	mgr := g.LuckyTimeBombFish
	mgr.mu.Lock()
	entry, ok := mgr.bombTargets[targetID]
	if !ok {
		mgr.mu.Unlock()
		return
	}
	instanceID := entry.instanceID
	delete(mgr.bombTargets, targetID)
	mgr.mu.Unlock()

	log.Printf("[LuckyTimeBomb] player=%s early detonated bomb target=%s", p.ID, targetID)

	// 廣播提前引爆
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTimeBombFish,
		Payload: ws.LuckyTimeBombFishPayload{
			Event:      "bomb_early_detonate",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			TargetID:   targetID,
			InstanceID: instanceID,
			X:          targetX,
			Y:          targetY,
			Mult:       LuckyTimeBombEarlyMult,
		},
	})

	// 連鎖爆炸：引爆周圍 150px 內目標
	go g.doTimeBombChain(p, targetID, targetX, targetY, instanceID)
}

// doTimeBombChain 連鎖爆炸（提前引爆後觸發）
func (g *Game) doTimeBombChain(p *player.Player, sourceID string, sourceX, sourceY float64, instanceID string) {
	g.mu.Lock()

	type chainResult struct {
		TargetID string
		Killed   bool
		Reward   int
		X        float64
		Y        float64
	}

	var results []chainResult
	totalReward := 0
	killedCount := 0

	for id, t := range g.Targets {
		if id == sourceID || t.HP <= 0 {
			continue
		}
		dx := t.X - sourceX
		dy := t.Y - sourceY
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist > LuckyTimeBombChainRadius {
			continue
		}

		// 60% 連鎖擊破機率
		killed := rand.Float64() < LuckyTimeBombChainChance
		reward := 0

		if killed {
			t.HP = 0
			reward = int(t.Multiplier * float64(1) * LuckyTimeBombChainMult)
			if reward < 1 {
				reward = 1
			}
			totalReward += reward
			killedCount++
			delete(g.Targets, id)
		}

		results = append(results, chainResult{
			TargetID: id,
			Killed:   killed,
			Reward:   reward,
			X:        t.X,
			Y:        t.Y,
		})
	}
	g.mu.Unlock()

	// 個人獎勵
	if totalReward > 0 {
		p.AddCoins(totalReward)
	}

	log.Printf("[LuckyTimeBomb] chain explosion: killed=%d totalReward=%d", killedCount, totalReward)

	// 廣播連鎖爆炸結果
	type chainResultPayload struct {
		TargetID string  `json:"target_id"`
		Killed   bool    `json:"killed"`
		Reward   int     `json:"reward"`
		X        float64 `json:"x"`
		Y        float64 `json:"y"`
	}
	payloadResults := make([]chainResultPayload, 0, len(results))
	for _, r := range results {
		payloadResults = append(payloadResults, chainResultPayload{
			TargetID: r.TargetID,
			Killed:   r.Killed,
			Reward:   r.Reward,
			X:        r.X,
			Y:        r.Y,
		})
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTimeBombFish,
		Payload: ws.LuckyTimeBombFishPayload{
			Event:        "bomb_chain_blast",
			PlayerID:     p.ID,
			InstanceID:   instanceID,
			SourceX:      sourceX,
			SourceY:      sourceY,
			KilledCount:  killedCount,
			TotalReward:  totalReward,
			ChainResults: payloadResults,
		},
	})
}

// tryLuckyTimeBombFish 擊破 T193 後觸發時間炸彈（供 handleKill 使用）
func (g *Game) tryLuckyTimeBombFish(p *player.Player) {
	mgr := g.LuckyTimeBombFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}

	// 設定個人冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyTimeBombPersonalCD)

	// 生成 instanceID
	instanceID := fmt.Sprintf("timebomb_%d", time.Now().UnixNano())
	mgr.currentInstanceID = instanceID
	mgr.mu.Unlock()

	// 選取隨機目標加上炸彈標記
	g.mu.Lock()
	var candidates []string
	for id, t := range g.Targets {
		if t.HP > 0 {
			candidates = append(candidates, id)
		}
	}

	// 隨機選取最多 LuckyTimeBombCount 個目標
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	if len(candidates) > LuckyTimeBombCount {
		candidates = candidates[:LuckyTimeBombCount]
	}

	type bombInfo struct {
		ID string
		X  float64
		Y  float64
	}
	var bombInfos []bombInfo
	expiresAt := time.Now().Add(LuckyTimeBombFuseTime)

	for _, id := range candidates {
		t := g.Targets[id]
		bombInfos = append(bombInfos, bombInfo{ID: id, X: t.X, Y: t.Y})
	}
	g.mu.Unlock()

	if len(bombInfos) == 0 {
		return
	}

	// 加入炸彈標記
	mgr.mu.Lock()
	for _, info := range bombInfos {
		mgr.bombTargets[info.ID] = &timeBombEntry{
			instanceID: instanceID,
			expiresAt:  expiresAt,
		}
	}
	mgr.mu.Unlock()

	log.Printf("[LuckyTimeBomb] player=%s placed %d bombs (instance=%s)", p.ID, len(bombInfos), instanceID)

	// 廣播炸彈標記
	type bombTarget struct {
		ID string  `json:"id"`
		X  float64 `json:"x"`
		Y  float64 `json:"y"`
	}
	targets := make([]bombTarget, 0, len(bombInfos))
	for _, info := range bombInfos {
		targets = append(targets, bombTarget{ID: info.ID, X: info.X, Y: info.Y})
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTimeBombFish,
		Payload: ws.LuckyTimeBombFishPayload{
			Event:       "bomb_placed",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			InstanceID:  instanceID,
			BombCount:   len(bombInfos),
			FuseSec:     int(LuckyTimeBombFuseTime.Seconds()),
			BombTargets: targets,
			EarlyMult:   LuckyTimeBombEarlyMult,
			AutoMult:    LuckyTimeBombAutoMult,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyTimeBombFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("💣 %s 放置了 %d 個時間炸彈！%d 秒後自動爆炸！",
			p.DisplayName, len(bombInfos), int(LuckyTimeBombFuseTime.Seconds())),
		"color": "#E74C3C",
	})
	g.broadcastAnnouncement(ann)

	// 啟動倒數計時 goroutine
	for _, info := range bombInfos {
		go g.runTimeBombFuse(p, info.ID, info.X, info.Y, instanceID, expiresAt)
	}
}

// runTimeBombFuse 單個炸彈倒數計時（goroutine）
func (g *Game) runTimeBombFuse(p *player.Player, targetID string, targetX, targetY float64, instanceID string, expiresAt time.Time) {
	// 每秒廣播倒數
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	remaining := int(LuckyTimeBombFuseTime.Seconds())

	for {
		select {
		case <-ticker.C:
			remaining--

			// 確認炸彈仍有效（可能已被提前引爆）
			mgr := g.LuckyTimeBombFish
			mgr.mu.Lock()
			entry, ok := mgr.bombTargets[targetID]
			if !ok || entry.instanceID != instanceID {
				mgr.mu.Unlock()
				return // 已被提前引爆或清除
			}
			mgr.mu.Unlock()

			if remaining <= 0 {
				// 倒數結束，自動爆炸
				g.doTimeBombAutoExplode(p, targetID, targetX, targetY, instanceID)
				return
			}

			// 廣播倒數更新（每秒）
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgLuckyTimeBombFish,
				Payload: ws.LuckyTimeBombFishPayload{
					Event:      "bomb_countdown",
					TargetID:   targetID,
					InstanceID: instanceID,
					Remaining:  remaining,
				},
			})
		}
	}
}

// doTimeBombAutoExplode 炸彈自動爆炸（倒數結束）
func (g *Game) doTimeBombAutoExplode(p *player.Player, targetID string, targetX, targetY float64, instanceID string) {
	// 移除炸彈標記
	mgr := g.LuckyTimeBombFish
	mgr.mu.Lock()
	entry, ok := mgr.bombTargets[targetID]
	if !ok || entry.instanceID != instanceID {
		mgr.mu.Unlock()
		return // 已被提前引爆
	}
	delete(mgr.bombTargets, targetID)
	mgr.mu.Unlock()

	// 嘗試擊破目標（80% 機率）
	killed := false
	reward := 0

	g.mu.Lock()
	t, exists := g.Targets[targetID]
	if exists && t.HP > 0 {
		if rand.Float64() < LuckyTimeBombAutoChance {
			killed = true
			reward = int(t.Multiplier * float64(1) * LuckyTimeBombAutoMult)
			if reward < 1 {
				reward = 1
			}
			t.HP = 0
			delete(g.Targets, targetID)
		}
	}
	g.mu.Unlock()

	// 個人獎勵
	if killed && reward > 0 {
		p.AddCoins(reward)
	}

	log.Printf("[LuckyTimeBomb] auto explode target=%s killed=%v reward=%d", targetID, killed, reward)

	// 廣播自動爆炸結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTimeBombFish,
		Payload: ws.LuckyTimeBombFishPayload{
			Event:      "bomb_auto_explode",
			PlayerID:   p.ID,
			TargetID:   targetID,
			InstanceID: instanceID,
			X:          targetX,
			Y:          targetY,
			Killed:     killed,
			Reward:     reward,
			Mult:       LuckyTimeBombAutoMult,
		},
	})
}
