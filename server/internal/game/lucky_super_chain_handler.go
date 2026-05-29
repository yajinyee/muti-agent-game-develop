// lucky_super_chain_handler.go — T230 幸運超級連鎖魚
// 設計：Super Chain 機制
//       每次擊破觸發 3 條連鎖（每條 ×80.0），連鎖 ≥5 次 → 超級連鎖爆發
//       超級連鎖爆發 → 全服 ×46.0 加成 92 秒（超越 T229 的 ×45.5）
//       業界依據：Royal Fishing 連鎖電擊升級版 + Super Chain 概念（2026）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type superChainBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckySuperChainManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *superChainBoost
}

func newLuckySuperChainManager() *luckySuperChainManager {
	return &luckySuperChainManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckySuperChainFish(defID string) bool {
	return defID == "T230"
}

func (m *luckySuperChainManager) getSuperChainMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckySuperChainManager) tryLuckySuperChainFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(430 * time.Second)
	m.personalCD[p.ID] = now.Add(370 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckySuperChain,
		Payload: map[string]interface{}{
			"event":        "chain_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"chain_count":  5,
			"per_mult":     80.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("SUPER CHAIN! %s triggered 5-chain reaction! Each chain x80.0!", p.GetDisplayName()), "critical", "#00FFFF")
	log.Printf("[LuckySuperChain] %s triggered Super Chain fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		totalMult := 0.0
		chainCount := 0

		for i := 1; i <= 5; i++ {
			delay := 1500
			if i > 3 {
				delay = 800
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)

			chainMult := 80.0
			totalMult += chainMult
			chainCount++
			reward := int(chainMult * betCost)
			g.mu.Lock()
			p.Coins += reward
			for _, t := range g.targets {
				damage := int(float64(t.MaxHP) * 0.40)
				t.HP -= damage
				if t.HP <= 0 {
					t.HP = 0
				}
			}
			g.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckySuperChain,
				Payload: map[string]interface{}{
					"event":       "chain_hit",
					"chain_index": i,
					"chain_mult":  chainMult,
					"total_mult":  totalMult,
					"reward":      reward,
					"is_bonus":    i > 3,
				},
			})
		}

		time.Sleep(1 * time.Second)
		isPerfect := chainCount >= 5
		globalBonus := 46.0
		globalDuration := 92

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckySuperChain,
			Payload: map[string]interface{}{
				"event":          "chain_result",
				"chain_count":    chainCount,
				"total_mult":     totalMult,
				"is_perfect":     isPerfect,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
			},
		})

		if isPerfect {
			m.mu.Lock()
			m.boost = &superChainBoost{
				mult:      globalBonus,
				expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
			}
			m.mu.Unlock()
			g.sendAnnounce(fmt.Sprintf("SUPER CHAIN BURST! %s chained %d times! Total x%.1f! Global x%.1f for %ds!", p.GetDisplayName(), chainCount, totalMult, globalBonus, globalDuration), "critical", "#00FFFF")
		}

		log.Printf("[LuckySuperChain] %s: chains=%d, total_mult=%.1f, perfect=%v", p.GetDisplayName(), chainCount, totalMult, isPerfect)
	}()

	return true
}
