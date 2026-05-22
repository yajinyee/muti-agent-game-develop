// crocodile_hunter_handler.go — 巨型鱷魚獵食系統 handler（DAY-188）
// 業界依據：JILI Mega Fishing「giant crocodiles awaken to hunt fish on the fish farm
//  to accumulate big prizes」
// 設計：T146 巨型鱷魚出現後，每 2 秒主動「獵食」場上一個目標（優先高倍率普通目標），
// 獵食成功給全服玩家共享獎勵（0.4x 倍率），玩家擊破鱷魚本身獲得鱷魚倍率 + 累積獎池 50%
// 設計差異：與不死 BOSS（玩家打它）不同，鱷魚是「它主動打其他魚」，
// 製造「看著鱷魚在場上橫行霸道」的緊張感，玩家需要決策：讓它繼續獵食累積獎池，還是立刻擊破？
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	crocodileHunterCooldown    = 50 * time.Second // 全服冷卻
	crocodileHunterHuntInterval = 2 * time.Second  // 每次獵食間隔
	crocodileHunterMaxHunts    = 8                 // 最多獵食次數
	crocodileHunterKillChance  = 0.85              // 獵食成功機率
	crocodileHunterRewardMult  = 0.40              // 獵食獎勵倍率（全服共享）
	crocodileHunterPoolShare   = 0.50              // 玩家擊破時獲得獎池比例
)

// crocodileHunterHuntEntry 單次獵食記錄（DAY-188，避免與 crocodile_handler.go 的 crocodileHuntEntry 衝突）
type crocodileHunterHuntEntry struct {
	TargetDefID  string  `json:"target_def_id"`
	TargetName   string  `json:"target_name"`
	TargetMult   float64 `json:"target_mult"`
	Reward       int     `json:"reward"`
	HuntIndex    int     `json:"hunt_index"`
}

// crocodileHunterManager 巨型鱷魚獵食系統管理器（全服共享）
type crocodileHunterManager struct {
	mu           sync.Mutex
	isActive     bool
	instanceID   string
	huntCount    int       // 已獵食次數
	totalPool    int       // 累積獎池（全服共享）
	lastCooldown time.Time // 上次觸發時間
	stopHunt     chan struct{}
}

func newCrocodileHunterManager() *crocodileHunterManager {
	return &crocodileHunterManager{
		stopHunt: make(chan struct{}, 1),
	}
}

// isCrocodileHunter 判斷是否為巨型鱷魚
func isCrocodileHunter(defID string) bool {
	return defID == "T146"
}

// canTrigger 是否可以觸發（冷卻檢查）
func (m *crocodileHunterManager) canTrigger() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isActive {
		return false
	}
	return time.Since(m.lastCooldown) >= crocodileHunterCooldown
}

// startHunt 開始獵食模式
func (m *crocodileHunterManager) startHunt(instanceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isActive = true
	m.instanceID = instanceID
	m.huntCount = 0
	m.totalPool = 0
	m.lastCooldown = time.Now()
	// 重置 stop channel
	select {
	case <-m.stopHunt:
	default:
	}
}

// recordHunt 記錄一次獵食，回傳是否繼續
func (m *crocodileHunterManager) recordHunt(reward int) (huntIndex int, shouldContinue bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.huntCount++
	m.totalPool += reward
	huntIndex = m.huntCount
	shouldContinue = m.huntCount < crocodileHunterMaxHunts && m.isActive
	return
}

// getPool 取得當前累積獎池
func (m *crocodileHunterManager) getPool() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.totalPool
}

// stopAndGetPool 停止獵食，回傳累積獎池
func (m *crocodileHunterManager) stopAndGetPool() (pool int, huntCount int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.isActive {
		return 0, 0
	}
	m.isActive = false
	pool = m.totalPool
	huntCount = m.huntCount
	// 發送停止訊號
	select {
	case m.stopHunt <- struct{}{}:
	default:
	}
	return
}

// isHunting 是否正在獵食
func (m *crocodileHunterManager) isHunting() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isActive
}

// tryCrocodileHunterSpawn 鱷魚生成時觸發獵食模式（由 spawnTarget 呼叫）
func (g *Game) tryCrocodileHunterSpawn(instanceID string) {
	if g.CrocodileHunter == nil {
		return
	}
	if !g.CrocodileHunter.canTrigger() {
		return
	}

	g.CrocodileHunter.startHunt(instanceID)

	log.Printf("[CrocodileHunter] spawned: instance=%s, max_hunts=%d", instanceID, crocodileHunterMaxHunts)

	// 全服廣播：巨型鱷魚出現
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCrocodileHunter,
		Payload: ws.CrocodileHunterPayload{
			Phase:      "croc_appear",
			InstanceID: instanceID,
			MaxHunts:   crocodileHunterMaxHunts,
			Message:    "🐊 巨型鱷魚出現！牠將主動獵食場上的魚！擊破牠可獲得累積獎池！",
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventBossWarning, "巨型鱷魚", 0, map[string]string{
		"message": "🐊 巨型鱷魚降臨！牠正在獵食場上的魚，累積大獎！",
	})
	g.broadcastAnnouncement(ann)

	// 啟動獵食 goroutine
	go g.runCrocodileHunting(instanceID)
}

// runCrocodileHunting 執行鱷魚獵食循環（goroutine）
func (g *Game) runCrocodileHunting(instanceID string) {
	ticker := time.NewTicker(crocodileHunterHuntInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !g.CrocodileHunter.isHunting() {
				return
			}
			g.doCrocodileHunt(instanceID)

		case <-g.CrocodileHunter.stopHunt:
			log.Printf("[CrocodileHunter] hunt stopped by signal: instance=%s", instanceID)
			return

		case <-g.stopCh:
			return
		}
	}
}

// doCrocodileHunt 執行一次獵食
func (g *Game) doCrocodileHunt(instanceID string) {
	// 選擇獵食目標（優先高倍率普通目標，不選 BOSS 和特殊目標）
	g.mu.RLock()
	type targetCandidate struct {
		id   string
		mult float64
		name string
		x, y float64
	}
	var candidates []targetCandidate
	for id, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		if t.DefID == "B001" || t.DefID == "T144" || t.DefID == "T146" {
			continue // 不獵食 BOSS、龍龜、自己
		}
		candidates = append(candidates, targetCandidate{
			id:   id,
			mult: t.Multiplier,
			name: t.Def.Name, // 使用 TargetDef.Name 取得中文名稱
			x:    t.X,
			y:    t.Y,
		})
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		log.Printf("[CrocodileHunter] no targets to hunt, skipping")
		return
	}

	// 按倍率排序，優先選高倍率目標（前 30% 中隨機選）
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].mult > candidates[j].mult
	})
	topN := len(candidates) / 3
	if topN < 1 {
		topN = 1
	}
	chosen := candidates[rand.Intn(topN)]

	// 判斷獵食成功機率
	if rand.Float64() >= crocodileHunterKillChance {
		// 獵食失敗（目標逃脫）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgCrocodileHunter,
			Payload: ws.CrocodileHunterPayload{
				Phase:      "croc_miss",
				InstanceID: instanceID,
				TargetX:    chosen.x,
				TargetY:    chosen.y,
				Message:    fmt.Sprintf("🐊 鱷魚撲空了！%s 逃脫！", chosen.name),
			},
		})
		return
	}
	// 獵食成功：計算獎勵（全服共享）
	// 獎勵 = 目標倍率 × 全服平均 betLevel × 0.40
	avgBetLevel := g.getAverageBetLevel()
	reward := int(chosen.mult * float64(avgBetLevel) * crocodileHunterRewardMult)
	if reward < 1 {
		reward = 1
	}

	// 移除目標
	g.mu.Lock()
	t, exists := g.Targets[chosen.id]
	if !exists || t.HP <= 0 {
		g.mu.Unlock()
		return
	}
	t.HP = 0
	delete(g.Targets, chosen.id)
	g.mu.Unlock()

	// 廣播目標被鱷魚獵食（全服看到）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTargetKill,
		Payload: ws.TargetKillPayload{
			InstanceID: chosen.id,
			KillerID:   "crocodile_hunter",
			Reward:     0, // 鱷魚獵食不給個人獎勵，給全服共享
			Multiplier: chosen.mult,
		},
	})

	// 記錄獵食，更新獎池
	huntIndex, shouldContinue := g.CrocodileHunter.recordHunt(reward)
	currentPool := g.CrocodileHunter.getPool()

	log.Printf("[CrocodileHunter] hunt #%d: target=%s(%.0fx), reward=%d, pool=%d",
		huntIndex, chosen.name, chosen.mult, reward, currentPool)

	// 廣播獵食成功（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCrocodileHunter,
		Payload: ws.CrocodileHunterPayload{
			Phase:       "croc_hunt",
			InstanceID:  instanceID,
			HuntIndex:   huntIndex,
			TargetName:  chosen.name,
			TargetMult:  chosen.mult,
			TargetX:     chosen.x,
			TargetY:     chosen.y,
			HuntReward:  reward,
			TotalPool:   currentPool,
			MaxHunts:    crocodileHunterMaxHunts,
			Message:     fmt.Sprintf("🐊 鱷魚獵食了 %s（%.0fx）！獎池 +%d！", chosen.name, chosen.mult, reward),
		},
	})

	// 全服公告（≥4 次獵食時）
	if huntIndex >= 4 {
		ann := g.Announce.Create(announce.EventMegaWin, "巨型鱷魚", currentPool, map[string]string{
			"message": fmt.Sprintf("🐊 巨型鱷魚已獵食 %d 次！累積獎池 %d 金幣！快去擊破牠！", huntIndex, currentPool),
		})
		g.broadcastAnnouncement(ann)
	}

	// 達到最大獵食次數，鱷魚離開
	if !shouldContinue {
		g.onCrocodileHunterLeave(instanceID, "max_hunts")
	}
}

// notifyCrocodileHunterKill 玩家擊破鱷魚時呼叫（由 handleKill 呼叫）
func (g *Game) notifyCrocodileHunterKill(p *player.Player, instanceID string, baseMult float64) {
	if g.CrocodileHunter == nil {
		return
	}

	pool, huntCount := g.CrocodileHunter.stopAndGetPool()

	// 玩家獲得：鱷魚基礎獎勵 + 累積獎池 50%
	baseReward := int(baseMult * float64(p.BetLevel))
	poolBonus := int(float64(pool) * crocodileHunterPoolShare)
	totalReward := baseReward + poolBonus

	// 發放獎勵
	p.AddCoins(totalReward)

	log.Printf("[CrocodileHunter] killed by player=%s, base=%d, pool_bonus=%d, total=%d, hunts=%d",
		p.ID, baseReward, poolBonus, totalReward, huntCount)

	// 廣播鱷魚被擊破（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCrocodileHunter,
		Payload: ws.CrocodileHunterPayload{
			Phase:       "croc_killed",
			InstanceID:  instanceID,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			HuntCount:   huntCount,
			TotalPool:   pool,
			PoolBonus:   poolBonus,
			BaseReward:  baseReward,
			TotalReward: totalReward,
			NewBalance:  p.GetCoins(),
			Message:     fmt.Sprintf("🐊 %s 擊破了巨型鱷魚！獲得基礎獎勵 %d + 獎池 %d = 共 %d 金幣！", p.DisplayName, baseReward, poolBonus, totalReward),
		},
	})

	// 全服公告
	if totalReward >= 500 || huntCount >= 4 {
		ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("🐊 %s 擊破巨型鱷魚！獵食 %d 次後獲得 %d 金幣大獎！", p.DisplayName, huntCount, totalReward),
		})
		g.broadcastAnnouncement(ann)
	}

	// 動態牆
	if totalReward >= 300 {
		go g.notifyFeedMegaWin(p, baseMult, totalReward)
	}
}

// onCrocodileHunterLeave 鱷魚離開（達到最大獵食次數或超時）
func (g *Game) onCrocodileHunterLeave(instanceID string, reason string) {
	if g.CrocodileHunter == nil {
		return
	}

	pool, huntCount := g.CrocodileHunter.stopAndGetPool()
	if pool == 0 && huntCount == 0 {
		return // 已被玩家擊破，不重複廣播
	}

	log.Printf("[CrocodileHunter] left: reason=%s, hunts=%d, pool=%d", reason, huntCount, pool)

	// 廣播鱷魚離開（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCrocodileHunter,
		Payload: ws.CrocodileHunterPayload{
			Phase:      "croc_leave",
			InstanceID: instanceID,
			HuntCount:  huntCount,
			TotalPool:  pool,
			Message:    fmt.Sprintf("🐊 巨型鱷魚離去了！共獵食 %d 次，累積獎池 %d 金幣未被領取！", huntCount, pool),
		},
	})
}

// getAverageBetLevel 取得全服平均 betLevel（用於計算全服共享獎勵）
func (g *Game) getAverageBetLevel() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if len(g.Players) == 0 {
		return 3 // 預設值
	}
	total := 0
	for _, p := range g.Players {
		total += p.BetLevel
	}
	avg := total / len(g.Players)
	if avg < 1 {
		avg = 1
	}
	return avg
}
