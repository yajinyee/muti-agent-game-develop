// lucky_fisherman_trail_handler.go — T251 幸運漁夫路徑魚
// 設計：Bigger Bites 進階路徑機制（Reflex Gaming + Bragg 2026）
//       Fishermen 符號收集 + 路徑升級（10 個節點）
//       每個節點：魚升級 + 額外旋轉 + 倍率提升
//       路徑完成（≥8節點）→ 全服 ×57.5 加成 115 秒
//       業界依據：Reflex Gaming「Big Game Fishing Bigger Bites」進階路徑機制（2026-02）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type fishermanTrailBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyFishermanTrailManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *fishermanTrailBoost
}

func newLuckyFishermanTrailManager() *luckyFishermanTrailManager {
	return &luckyFishermanTrailManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyFishermanTrailFish(defID string) bool {
	return defID == "T251"
}

func (m *luckyFishermanTrailManager) getFishermanTrailMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyFishermanTrailManager) tryLuckyFishermanTrailFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(650 * time.Second)
	m.personalCD[p.ID] = now.Add(590 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyFishermanTrail,
		Payload: map[string]interface{}{
			"event":        "fisherman_trail_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"trail_nodes":  10,
			"max_mult":     500.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("FISHERMAN TRAIL! %s activated Bigger Bites system! 10 trail nodes, max x500!", p.GetDisplayName()), "critical", "#FF8C00")
	log.Printf("[LuckyFishermanTrail] %s triggered Fisherman Trail fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		totalReward := 0
		nodesReached := 0

		// 路徑節點：10 個，每個節點倍率遞增
		nodeMults := []float64{10.0, 20.0, 30.0, 50.0, 75.0, 100.0, 150.0, 200.0, 350.0, 500.0}
		for i, mult := range nodeMults {
			time.Sleep(400 * time.Millisecond)
			nodesReached++
			reward := int(mult * betCost)
			totalReward += reward
			g.mu.Lock()
			p.Coins += reward
			g.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyFishermanTrail,
				Payload: map[string]interface{}{
					"event":    "trail_node",
					"node_no":  i + 1,
					"mult":     mult,
					"reward":   reward,
					"upgrade":  fmt.Sprintf("Fish Upgrade Lv.%d", i+1),
				},
			})
		}

		// 路徑完成（≥8節點）→ 全服 ×57.5 加成 115 秒
		globalBonus := 57.5
		globalDuration := 115
		if nodesReached >= 8 {
			m.mu.Lock()
			m.boost = &fishermanTrailBoost{
				mult:      globalBonus,
				expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
			}
			m.mu.Unlock()
		}

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyFishermanTrail,
			Payload: map[string]interface{}{
				"event":          "fisherman_trail_complete",
				"nodes_reached":  nodesReached,
				"total_reward":   totalReward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
			},
		})
		g.sendAnnounce(fmt.Sprintf("FISHERMAN TRAIL COMPLETE! %s: %d nodes! GLOBAL x%.1f for %ds!", p.GetDisplayName(), nodesReached, globalBonus, globalDuration), "critical", "#FF8C00")
		log.Printf("[LuckyFishermanTrail] %s: nodes=%d, global=x%.1f", p.GetDisplayName(), nodesReached, globalBonus)
	}()

	return true
}
