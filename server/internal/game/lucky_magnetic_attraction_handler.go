// lucky_magnetic_attraction_handler.go — T229 幸運磁力吸引魚
// 設計：Magnetic Attraction 機制
//       磁力吸引全場所有目標到中心，每個目標獎勵 ×70.0
//       命中 ≥10 個目標 → 完美磁力，全服 ×45.5 加成 91 秒（超越 T228 的 ×45.0）
//       業界依據：Black Hole Fishing 引力機制升級版 + Magnetic Attraction 概念（2026）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type magneticAttractionBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyMagneticAttractionManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *magneticAttractionBoost
}

func newLuckyMagneticAttractionManager() *luckyMagneticAttractionManager {
	return &luckyMagneticAttractionManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyMagneticAttractionFish(defID string) bool {
	return defID == "T229"
}

func (m *luckyMagneticAttractionManager) getMagneticAttractionMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyMagneticAttractionManager) tryLuckyMagneticAttractionFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(420 * time.Second)
	m.personalCD[p.ID] = now.Add(360 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyMagneticAttraction,
		Payload: map[string]interface{}{
			"event":        "magnetic_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("MAGNETIC ATTRACTION! %s activated magnetic force! All targets pulled to center!", p.GetDisplayName()), "critical", "#FF6600")
	log.Printf("[LuckyMagneticAttraction] %s triggered Magnetic Attraction fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		time.Sleep(2 * time.Second)

		g.mu.Lock()
		targetCount := len(g.targets)
		for _, t := range g.targets {
			damage := int(float64(t.MaxHP) * 0.85)
			t.HP -= damage
			if t.HP <= 0 {
				t.HP = 0
			}
		}
		g.mu.Unlock()

		if targetCount == 0 {
			targetCount = 8
		}

		perTargetMult := 70.0
		totalMult := float64(targetCount) * perTargetMult
		reward := int(totalMult * betCost)
		if reward > 0 {
			g.mu.Lock()
			p.Coins += reward
			g.mu.Unlock()
		}

		isPerfect := targetCount >= 10
		globalBonus := 45.5
		globalDuration := 91

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyMagneticAttraction,
			Payload: map[string]interface{}{
				"event":          "magnetic_result",
				"hit_count":      targetCount,
				"total_mult":     totalMult,
				"reward":         reward,
				"is_perfect":     isPerfect,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
			},
		})

		if isPerfect {
			m.mu.Lock()
			m.boost = &magneticAttractionBoost{
				mult:      globalBonus,
				expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
			}
			m.mu.Unlock()
			g.sendAnnounce(fmt.Sprintf("PERFECT MAGNETIC ATTRACTION! %s pulled %d targets! Global x%.1f for %ds!", p.GetDisplayName(), targetCount, globalBonus, globalDuration), "critical", "#FF6600")
		}

		// 隨機觸發額外磁力波
		if rand.Float64() < 0.3 {
			time.Sleep(2 * time.Second)
			bonusMult := 35.0 + rand.Float64()*35.0
			bonusReward := int(bonusMult * betCost)
			g.mu.Lock()
			p.Coins += bonusReward
			g.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyMagneticAttraction,
				Payload: map[string]interface{}{
					"event":        "magnetic_bonus",
					"bonus_mult":   bonusMult,
					"bonus_reward": bonusReward,
				},
			})
		}

		log.Printf("[LuckyMagneticAttraction] %s: hit=%d, total_mult=%.1f, perfect=%v", p.GetDisplayName(), targetCount, totalMult, isPerfect)
	}()

	return true
}
