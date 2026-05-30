// lucky_electrical_frame_handler.go — T249 幸運電擊框架魚
// 設計：Catfish Hunters 電擊框架機制（Nolimit City 2026）
//       電擊框架系統：每次命中全局倍率翻倍（×1→×2→×4→...→×1024）
//       最多 10 次命中（×1024），每次命中獎勵 = 目標值 × 全局倍率
//       完美連鎖（≥8次）→ 全服 ×56.5 加成 113 秒
//       業界依據：Nolimit City「Catfish Hunters」電擊框架 + 全局倍率翻倍（2026-03）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type electricalFrameBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyElectricalFrameManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *electricalFrameBoost
}

func newLuckyElectricalFrameManager() *luckyElectricalFrameManager {
	return &luckyElectricalFrameManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyElectricalFrameFish(defID string) bool {
	return defID == "T249"
}

func (m *luckyElectricalFrameManager) getElectricalFrameMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyElectricalFrameManager) tryLuckyElectricalFrameFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(630 * time.Second)
	m.personalCD[p.ID] = now.Add(570 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyElectricalFrame,
		Payload: map[string]interface{}{
			"event":        "electrical_frame_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"max_hits":     10,
			"max_mult":     1024,
		},
	})
	g.sendAnnounce(fmt.Sprintf("ELECTRICAL FRAME! %s activated Catfish Hunters system! Global multiplier doubles each hit! Max x1024!", p.GetDisplayName()), "critical", "#00FFFF")
	log.Printf("[LuckyElectricalFrame] %s triggered Electrical Frame fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		totalReward := 0
		globalMult := 1.0
		hitCount := 0

		// 電擊框架：每次命中全局倍率翻倍，最多 10 次
		fishValues := []float64{5.0, 8.0, 12.0, 15.0, 20.0, 25.0, 30.0, 40.0, 50.0, 60.0}
		for i, fishVal := range fishValues {
			time.Sleep(300 * time.Millisecond)
			globalMult *= 2.0
			if globalMult > 1024.0 {
				globalMult = 1024.0
			}
			hitCount++
			reward := int(fishVal * globalMult * betCost * 0.1)
			totalReward += reward
			g.mu.Lock()
			p.Coins += reward
			g.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyElectricalFrame,
				Payload: map[string]interface{}{
					"event":       "frame_hit",
					"hit_no":      i + 1,
					"fish_value":  fishVal,
					"global_mult": globalMult,
					"reward":      reward,
				},
			})
		}

		// 完美連鎖（≥8次）→ 全服 ×56.5 加成 113 秒
		globalBonus := 56.5
		globalDuration := 113
		if hitCount >= 8 {
			m.mu.Lock()
			m.boost = &electricalFrameBoost{
				mult:      globalBonus,
				expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
			}
			m.mu.Unlock()
		}

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyElectricalFrame,
			Payload: map[string]interface{}{
				"event":          "electrical_frame_complete",
				"hit_count":      hitCount,
				"final_mult":     globalMult,
				"total_reward":   totalReward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
			},
		})
		g.sendAnnounce(fmt.Sprintf("ELECTRICAL FRAME COMPLETE! %s: %d hits, final x%.0f! GLOBAL x%.1f for %ds!", p.GetDisplayName(), hitCount, globalMult, globalBonus, globalDuration), "critical", "#00FFFF")
		log.Printf("[LuckyElectricalFrame] %s: hits=%d, final_mult=x%.0f, global=x%.1f", p.GetDisplayName(), hitCount, globalMult, globalBonus)
	}()

	return true
}
