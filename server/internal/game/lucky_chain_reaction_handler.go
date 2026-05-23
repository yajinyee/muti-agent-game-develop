// lucky_chain_reaction_handler.go — 幸運連鎖反應魚系統（DAY-241）
// 業界原創「多米諾骨牌效應」機制
//
// 設計：擊破 T199 後，場上隨機選 1 個目標作為「連鎖起點」：
//   - 起點目標被擊破後，自動引爆距離最近的目標（100% 擊破，×1.4 倍率）
//   - 被引爆的目標再引爆下一個最近目標（×1.3 倍率）
//   - 連鎖最多 8 層，每層倍率遞減 0.1（×1.4 → ×1.3 → ... → ×0.7）
//   - 連鎖期間全服廣播每一層的引爆，製造「多米諾骨牌」的視覺爽感
//   - 個人冷卻 25 秒
//
// 設計差異：
//   - 與鏈鎖爆炸魚（T184，空間爆炸，60% 機率）不同，連鎖反應是「100% 確定引爆，但倍率遞減」
//   - 讓玩家有「看著連鎖一層一層爆開」的期待感
//   - 「8 層連鎖」讓玩家有「要把魚排成一排才能最大化連鎖」的策略感
//   - 「倍率遞減」確保 RTP 平衡（期望值 = 1.4+1.3+...+0.7 = 8.4 層平均 1.05x）
//   - 全服廣播每一層讓所有玩家都看到連鎖進度，製造「全服一起數層數」的社交感
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
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyChainReactionMaxLayers  = 8                  // 最大連鎖層數
	LuckyChainReactionBaseMult   = 1.4                // 第一層倍率
	LuckyChainReactionMultDecay  = 0.1                // 每層倍率遞減
	LuckyChainReactionMinMult    = 0.7                // 最低倍率（第8層）
	LuckyChainReactionSearchRange = 350.0             // 連鎖搜尋範圍（px）
	LuckyChainReactionLayerDelay = 400 * time.Millisecond // 每層引爆間隔（視覺效果）
	LuckyChainReactionPersonalCD = 25 * time.Second   // 個人冷卻
)

// chainReactionEntry 連鎖起點記錄
type chainReactionEntry struct {
	instanceID string
	expiresAt  time.Time
}

// luckyChainReactionManager 幸運連鎖反應魚管理器
type luckyChainReactionManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldown map[string]time.Time

	// 連鎖起點目標（instanceID → entry）
	chainStarters map[string]*chainReactionEntry

	// 當前連鎖狀態（防止同時多條連鎖）
	chainActive bool
}

func newLuckyChainReactionManager() *luckyChainReactionManager {
	return &luckyChainReactionManager{
		personalCooldown: make(map[string]time.Time),
		chainStarters:    make(map[string]*chainReactionEntry),
	}
}

// isLuckyChainReactionFish 判斷是否為幸運連鎖反應魚
func isLuckyChainReactionFish(defID string) bool {
	return defID == "T199"
}

// isChainReactionStarter 判斷目標是否為連鎖起點
func (g *Game) isChainReactionStarter(instanceID string) bool {
	mgr := g.LuckyChainReaction
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	entry, ok := mgr.chainStarters[instanceID]
	if !ok {
		return false
	}
	if time.Now().After(entry.expiresAt) {
		delete(mgr.chainStarters, instanceID)
		return false
	}
	return true
}

// removeChainReactionStarter 移除連鎖起點標記
func (g *Game) removeChainReactionStarter(instanceID string) {
	mgr := g.LuckyChainReaction
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	delete(mgr.chainStarters, instanceID)
}

// getLuckyChainReactionStarterMult 取得連鎖起點擊破倍率（×1.4，供 handleKill 使用）
func (g *Game) getLuckyChainReactionStarterMult(instanceID string) float64 {
	if g.isChainReactionStarter(instanceID) {
		return LuckyChainReactionBaseMult
	}
	return 1.0
}

// tryLuckyChainReactionFish 擊破 T199 後觸發連鎖反應
func (g *Game) tryLuckyChainReactionFish(p *player.Player) {
	mgr := g.LuckyChainReaction
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldown[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}
	mgr.personalCooldown[p.ID] = time.Now().Add(LuckyChainReactionPersonalCD)
	mgr.mu.Unlock()

	// 選取連鎖起點（隨機選 1 個存活目標，排除 BOSS 和 T199 本身）
	g.mu.RLock()
	var candidates []*target.Target
	for _, t := range g.Targets {
		if t.HP > 0 && t.DefID != "B001" && !isLuckyChainReactionFish(t.DefID) {
			candidates = append(candidates, t)
		}
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		log.Printf("[LuckyChainReaction] no candidates for chain starter, player=%s", p.ID)
		return
	}

	// 隨機選起點
	starter := candidates[rand.Intn(len(candidates))]

	// 標記連鎖起點（存活 15 秒）
	mgr.mu.Lock()
	mgr.chainStarters[starter.InstanceID] = &chainReactionEntry{
		instanceID: starter.InstanceID,
		expiresAt:  time.Now().Add(15 * time.Second),
	}
	mgr.mu.Unlock()

	log.Printf("[LuckyChainReaction] chain starter set: instanceID=%s defID=%s player=%s",
		starter.InstanceID, starter.DefID, p.ID)

	// 全服廣播：連鎖起點標記
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyChainReaction,
		Payload: ws.LuckyChainReactionPayload{
			Event:      "chain_start",
			PlayerName: p.DisplayName,
			StarterID:  starter.InstanceID,
			StarterDef: starter.DefID,
			MaxLayers:  LuckyChainReactionMaxLayers,
			BaseMult:   LuckyChainReactionBaseMult,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyChainReaction, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🔗 %s 觸發連鎖反應！擊破標記目標可引發最多 %d 層連鎖爆炸！",
			p.DisplayName, LuckyChainReactionMaxLayers),
		"color": "#FF6B35",
	})
	g.broadcastAnnouncement(ann)
}

// notifyChainReactionKill 玩家擊破連鎖起點時觸發連鎖
// 由 handleKill 呼叫
func (g *Game) notifyChainReactionKill(p *player.Player, killedTarget *target.Target) {
	mgr := g.LuckyChainReaction
	mgr.mu.Lock()

	// 防止同時多條連鎖
	if mgr.chainActive {
		mgr.mu.Unlock()
		return
	}
	mgr.chainActive = true
	mgr.mu.Unlock()

	log.Printf("[LuckyChainReaction] chain triggered by player=%s at instanceID=%s",
		p.ID, killedTarget.InstanceID)

	// 啟動連鎖 goroutine
	go g.runChainReaction(p, killedTarget.X, killedTarget.Y, 1)
}

// runChainReaction 執行連鎖反應（遞迴，每層延遲 400ms）
func (g *Game) runChainReaction(p *player.Player, fromX, fromY float64, layer int) {
	if layer > LuckyChainReactionMaxLayers {
		// 連鎖結束
		mgr := g.LuckyChainReaction
		mgr.mu.Lock()
		mgr.chainActive = false
		mgr.mu.Unlock()

		log.Printf("[LuckyChainReaction] chain completed: player=%s layers=%d",
			p.ID, LuckyChainReactionMaxLayers)

		// 廣播：連鎖結束
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyChainReaction,
			Payload: ws.LuckyChainReactionPayload{
				Event:      "chain_complete",
				PlayerName: p.DisplayName,
				Layer:      LuckyChainReactionMaxLayers,
			},
		})
		return
	}

	// 等待視覺效果間隔
	time.Sleep(LuckyChainReactionLayerDelay)

	// 計算本層倍率
	mult := LuckyChainReactionBaseMult - float64(layer-1)*LuckyChainReactionMultDecay
	if mult < LuckyChainReactionMinMult {
		mult = LuckyChainReactionMinMult
	}

	// 找最近的存活目標（排除 BOSS）
	g.mu.Lock()
	var nearest *target.Target
	nearestDist := math.MaxFloat64
	for _, t := range g.Targets {
		if t.HP <= 0 || t.DefID == "B001" {
			continue
		}
		dx := t.X - fromX
		dy := t.Y - fromY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < nearestDist && dist <= LuckyChainReactionSearchRange {
			nearest = t
			nearestDist = dist
		}
	}

	if nearest == nil {
		// 範圍內無目標，連鎖中斷
		g.mu.Unlock()

		mgr := g.LuckyChainReaction
		mgr.mu.Lock()
		mgr.chainActive = false
		mgr.mu.Unlock()

		log.Printf("[LuckyChainReaction] chain broken at layer=%d (no target in range)", layer)

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyChainReaction,
			Payload: ws.LuckyChainReactionPayload{
				Event:      "chain_broken",
				PlayerName: p.DisplayName,
				Layer:      layer,
			},
		})
		return
	}

	// 擊破目標（100% 機率）
	def, ok := data.Targets[nearest.DefID]
	if !ok {
		g.mu.Unlock()
		mgr := g.LuckyChainReaction
		mgr.mu.Lock()
		mgr.chainActive = false
		mgr.mu.Unlock()
		return
	}

	nextX := nearest.X
	nextY := nearest.Y
	instanceID := nearest.InstanceID
	defID := nearest.DefID
	nearest.HP = 0

	// 計算獎勵
	betDef := data.GetBetDef(p.BetLevel)
	reward := 0
	if betDef != nil {
		avgMult := (def.MultiplierMin + def.MultiplierMax) / 2.0
		reward = int(float64(betDef.BetCost) * avgMult * mult)
		p.Coins += reward
	}

	// 從場上移除
	delete(g.Targets, instanceID)
	g.mu.Unlock()

	log.Printf("[LuckyChainReaction] layer=%d target=%s reward=%d mult=%.1f",
		layer, defID, reward, mult)

	// 廣播：本層連鎖引爆
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyChainReaction,
		Payload: ws.LuckyChainReactionPayload{
			Event:      "chain_explode",
			PlayerName: p.DisplayName,
			Layer:      layer,
			MaxLayers:  LuckyChainReactionMaxLayers,
			TargetID:   instanceID,
			TargetDef:  defID,
			Mult:       mult,
			Reward:     reward,
			FromX:      fromX,
			FromY:      fromY,
			ToX:        nextX,
			ToY:        nextY,
		},
	})

	// 廣播目標被擊破
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTargetKill,
		Payload: ws.TargetKillPayload{
			InstanceID: instanceID,
			DefID:      defID,
			KillerID:   p.ID,
			Reward:     reward,
		},
	})

	// 繼續下一層連鎖
	go g.runChainReaction(p, nextX, nextY, layer+1)
}
