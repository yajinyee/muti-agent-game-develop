// lucky_ultimate_shark_handler.go — T243 幸運終極鯊魚魚
// 設計：Ultimate Shark 機制（里程碑：全服 ×53.0）
//       終極鯊魚清場：全場 HP 歸零（每個 ×180.0）+ 鯊魚咬合動畫
//       全服 ×53.0 加成 106 秒（新史上最高，超越 T238 的 ×50.0）
//       業界依據：Shark & Spark Hold & Win 終極版（2026-05-30）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type ultimateSharkBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyUltimateSharkManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *ultimateSharkBoost
}

func newLuckyUltimateSharkManager() *luckyUltimateSharkManager {
	return &luckyUltimateSharkManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyUltimateSharkFish(defID string) bool {
	return defID == "T243"
}

func (m *luckyUltimateSharkManager) getUltimateSharkMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyUltimateSharkManager) tryLuckyUltimateSharkFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(570 * time.Second)
	m.personalCD[p.ID] = now.Add(510 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyUltimateShark,
		Payload: map[string]interface{}{
			"event":         "ultimate_shark_start",
			"trigger_id":    p.ID,
			"trigger_name":  p.GetDisplayName(),
			"per_mult":      180.0,
			"bite_count":    14,
			"global_target": 53.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("ULTIMATE SHARK! %s triggered the Ultimate Shark! Every target x180.0! MILESTONE: Global x53.0!", p.GetDisplayName()), "critical", "#FF4500")
	log.Printf("[LuckyUltimateShark] %s triggered Ultimate Shark fish - MILESTONE x53.0", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		time.Sleep(2 * time.Second)

		// 全場清空
		g.mu.Lock()
		targetCount := len(g.targets)
		for _, t := range g.targets {
			t.HP = 0
		}
		g.mu.Unlock()

		if targetCount == 0 {
			targetCount = 16
		}

		perTargetMult := 180.0
		totalMult := float64(targetCount) * perTargetMult
		reward := int(totalMult * betCost)
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		// 鯊魚咬合動畫（14 次）
		for bite := 1; bite <= 14; bite++ {
			time.Sleep(200 * time.Millisecond)
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyUltimateShark,
				Payload: map[string]interface{}{
					"event":   "shark_bite",
					"bite_no": bite,
				},
			})
		}

		// 里程碑：全服 ×53.0（新史上最高）
		globalBonus := 53.0
		globalDuration := 106

		m.mu.Lock()
		m.boost = &ultimateSharkBoost{
			mult:      globalBonus,
			expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyUltimateShark,
			Payload: map[string]interface{}{
				"event":          "ultimate_shark_milestone",
				"target_count":   targetCount,
				"per_mult":       perTargetMult,
				"total_mult":     totalMult,
				"reward":         reward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
				"milestone":      "GLOBAL_53X",
			},
		})
		g.sendAnnounce(fmt.Sprintf("ULTIMATE SHARK MILESTONE! %s cleared %d targets! Total x%.1f! GLOBAL x%.1f for %ds! NEW RECORD!", p.GetDisplayName(), targetCount, totalMult, globalBonus, globalDuration), "critical", "#FF4500")
		log.Printf("[LuckyUltimateShark] MILESTONE! %s: targets=%d, total_mult=%.1f, global=x%.1f (NEW RECORD x53.0)", p.GetDisplayName(), targetCount, totalMult, globalBonus)
	}()

	return true
}
