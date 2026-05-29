// lucky_ice_fishing_master_handler.go — T236 幸運冰釣大師魚
// 設計：Ice Fishing Master 機制
//       5 次旋轉（每次最高 ×8000），旋轉結果累積，最高單次 ≥3000 → 完美冰釣
//       完美觸發 → 全服 ×49.0 加成 98 秒（超越 T235 的 ×48.5）
//       業界依據：Evolution Gaming「Ice Fishing Live」最高 5000x 升級版（2026）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type iceFishingMasterBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyIceFishingMasterManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *iceFishingMasterBoost
}

// 冰釣大師輪盤倍率權重（5次旋轉，最高 8000x）
var iceFishingMasterWeights = []struct {
	Mult   float64
	Weight int
}{
	{50, 35},
	{100, 25},
	{300, 15},
	{500, 10},
	{1000, 7},
	{2000, 4},
	{5000, 3},
	{8000, 1},
}

func newLuckyIceFishingMasterManager() *luckyIceFishingMasterManager {
	return &luckyIceFishingMasterManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyIceFishingMasterFish(defID string) bool {
	return defID == "T236"
}

func (m *luckyIceFishingMasterManager) getIceFishingMasterMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func spinIceFishingMasterWheel() float64 {
	total := 0
	for _, w := range iceFishingMasterWeights {
		total += w.Weight
	}
	r := rand.Intn(total)
	for _, w := range iceFishingMasterWeights {
		r -= w.Weight
		if r < 0 {
			return w.Mult
		}
	}
	return iceFishingMasterWeights[0].Mult
}

func (m *luckyIceFishingMasterManager) tryLuckyIceFishingMasterFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(500 * time.Second)
	m.personalCD[p.ID] = now.Add(440 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyIceFishingMaster,
		Payload: map[string]interface{}{
			"event":        "ice_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"spin_count":   5,
			"max_mult":     8000,
		},
	})
	g.sendAnnounce(fmt.Sprintf("ICE FISHING MASTER! %s activated! 5 spins, max x8000!", p.GetDisplayName()), "critical", "#00BFFF")
	log.Printf("[LuckyIceFishingMaster] %s triggered Ice Fishing Master fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		totalReward := 0
		maxSpin := 0.0
		spins := make([]float64, 5)

		for i := 0; i < 5; i++ {
			time.Sleep(1200 * time.Millisecond)
			spinMult := spinIceFishingMasterWheel()
			spins[i] = spinMult
			spinReward := int(spinMult * betCost)
			totalReward += spinReward
			if spinMult > maxSpin {
				maxSpin = spinMult
			}

			g.mu.Lock()
			p.Coins += spinReward
			g.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyIceFishingMaster,
				Payload: map[string]interface{}{
					"event":        "ice_spin",
					"spin_no":      i + 1,
					"spin_mult":    spinMult,
					"spin_reward":  spinReward,
					"total_reward": totalReward,
				},
			})
		}

		isPerfect := maxSpin >= 3000
		globalBonus := 49.0
		globalDuration := 98

		if isPerfect {
			m.mu.Lock()
			m.boost = &iceFishingMasterBoost{
				mult:      globalBonus,
				expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
			}
			m.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyIceFishingMaster,
				Payload: map[string]interface{}{
					"event":          "ice_perfect",
					"max_spin":       maxSpin,
					"total_reward":   totalReward,
					"global_bonus":   globalBonus,
					"global_seconds": globalDuration,
				},
			})
			g.sendAnnounce(fmt.Sprintf("PERFECT ICE FISHING! %s max spin x%.0f! Total reward %d! Global x%.1f for %ds!", p.GetDisplayName(), maxSpin, totalReward, globalBonus, globalDuration), "critical", "#00BFFF")
		} else {
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyIceFishingMaster,
				Payload: map[string]interface{}{
					"event":        "ice_end",
					"max_spin":     maxSpin,
					"total_reward": totalReward,
				},
			})
		}

		log.Printf("[LuckyIceFishingMaster] %s: max_spin=%.0f, total=%d, perfect=%v", p.GetDisplayName(), maxSpin, totalReward, isPerfect)
	}()

	return true
}
