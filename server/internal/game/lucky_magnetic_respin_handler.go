// lucky_magnetic_respin_handler.go — T250 幸運磁力連鎖魚
// 設計：Golden Gills 磁力連鎖 Respin 機制（Atomic Slot Lab 2026）
//       磁力連鎖 Respin：每次 Respin 磁力吸引相鄰目標（×75.0 旋轉倍率）
//       最多 8 次 Respin，每次吸引 1-3 個目標
//       完美連鎖（≥6次）→ 全服 ×57.0 加成 114 秒
//       業界依據：Atomic Slot Lab「Golden Gills」磁力連鎖 Respin + 75x 旋轉倍率（2026-02）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type magneticRespinBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyMagneticRespinManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *magneticRespinBoost
}

func newLuckyMagneticRespinManager() *luckyMagneticRespinManager {
	return &luckyMagneticRespinManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyMagneticRespinFish(defID string) bool {
	return defID == "T250"
}

func (m *luckyMagneticRespinManager) getMagneticRespinMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyMagneticRespinManager) tryLuckyMagneticRespinFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(640 * time.Second)
	m.personalCD[p.ID] = now.Add(580 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyMagneticRespin,
		Payload: map[string]interface{}{
			"event":        "magnetic_respin_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"max_respins":  8,
			"spin_mult":    75.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("MAGNETIC RESPIN! %s activated Golden Gills system! 8 respins with x75 multiplier!", p.GetDisplayName()), "critical", "#FFD700")
	log.Printf("[LuckyMagneticRespin] %s triggered Magnetic Respin fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		totalReward := 0
		respinCount := 0
		totalTargets := 0

		// 磁力連鎖 Respin：8 次，每次吸引 1-3 個目標
		for respin := 1; respin <= 8; respin++ {
			time.Sleep(350 * time.Millisecond)
			attracted := 1 + rand.Intn(3) // 1-3 個目標
			totalTargets += attracted
			spinMult := 75.0
			reward := int(spinMult * float64(attracted) * betCost)
			totalReward += reward
			respinCount++
			g.mu.Lock()
			p.Coins += reward
			g.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyMagneticRespin,
				Payload: map[string]interface{}{
					"event":     "respin",
					"respin_no": respin,
					"attracted": attracted,
					"spin_mult": spinMult,
					"reward":    reward,
				},
			})
		}

		// 完美連鎖（≥6次）→ 全服 ×57.0 加成 114 秒
		globalBonus := 57.0
		globalDuration := 114
		if respinCount >= 6 {
			m.mu.Lock()
			m.boost = &magneticRespinBoost{
				mult:      globalBonus,
				expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
			}
			m.mu.Unlock()
		}

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyMagneticRespin,
			Payload: map[string]interface{}{
				"event":          "magnetic_respin_complete",
				"respin_count":   respinCount,
				"total_targets":  totalTargets,
				"total_reward":   totalReward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
			},
		})
		g.sendAnnounce(fmt.Sprintf("MAGNETIC RESPIN COMPLETE! %s: %d respins, %d targets! GLOBAL x%.1f for %ds!", p.GetDisplayName(), respinCount, totalTargets, globalBonus, globalDuration), "critical", "#FFD700")
		log.Printf("[LuckyMagneticRespin] %s: respins=%d, targets=%d, global=x%.1f", p.GetDisplayName(), respinCount, totalTargets, globalBonus)
	}()

	return true
}
