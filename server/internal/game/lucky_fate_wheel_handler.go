// lucky_fate_wheel_handler.go — T178 幸運命運之輪魚
// 業界依據：「fate wheel mechanic」
// 設計：觸發命運之輪（8 個扇形，隨機停止），最高 ×50.0 單次獎勵
//       連續 3 次 ≥20x → 命運完美：全服 ×11.0 加成 24 秒
//       個人冷卻 68 秒；全服冷卻 108 秒
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

// 命運之輪 8 個扇形倍率
var fateWheelSlots = []float64{5, 10, 15, 20, 25, 30, 40, 50}
var fateWheelWeights = []int{30, 25, 18, 12, 7, 4, 3, 1}

type luckyFateWheelManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *fateWheelPerfectBoost
}

type fateWheelPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyFateWheelManager() *luckyFateWheelManager {
	return &luckyFateWheelManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyFateWheelFish(defID string) bool {
	return defID == "T178"
}

func (m *luckyFateWheelManager) getFateWheelPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func spinFateWheel() float64 {
	total := 0
	for _, w := range fateWheelWeights {
		total += w
	}
	r := rand.Intn(total)
	cum := 0
	for i, w := range fateWheelWeights {
		cum += w
		if r < cum {
			return fateWheelSlots[i]
		}
	}
	return fateWheelSlots[0]
}

func (m *luckyFateWheelManager) tryLuckyFateWheelFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) {
		m.mu.Unlock()
		return false
	}
	if cd, ok := m.personalCD[p.ID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return false
	}
	m.personalCD[p.ID] = now.Add(68 * time.Second)
	m.globalCD = now.Add(108 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_fate_wheel",
		Payload: map[string]interface{}{
			"event":        "fate_wheel_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"spins":        3,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🎡 命運之輪！%s 旋轉命運之輪！最高 ×50.0！", p.GetDisplayName()), "critical", "#F57F17")
	log.Printf("[LuckyFateWheel] %s 觸發命運之輪魚", p.GetDisplayName())

	go func() {
		highCount := 0
		totalReward := 0
		betCost := p.GetBetDef().BetCost
		for spin := 1; spin <= 3; spin++ {
			time.Sleep(time.Duration(spin) * 800 * time.Millisecond)
			mult := spinFateWheel()
			reward := int(mult * float64(betCost))
			totalReward += reward
			if mult >= 20 {
				highCount++
			}
			g.mu.Lock()
			p.AddCoins(reward)
			g.mu.Unlock()
			g.hub.Send(p.ID, protocol.MsgReward, protocol.RewardPayload{
				Source:     "fate_wheel",
				Amount:     reward,
				Multiplier: mult,
				NewBalance: p.Coins,
			})
			g.broadcast(protocol.Envelope{
				Type: "lucky_fate_wheel",
				Payload: map[string]interface{}{
					"event":        "wheel_spin",
					"spin_no":      spin,
					"mult":         mult,
					"reward":       reward,
					"trigger_name": p.GetDisplayName(),
				},
			})
		}
		if highCount >= 3 {
			m.mu.Lock()
			m.perfectBoost = &fateWheelPerfectBoost{
				mult:      11.0,
				expiresAt: time.Now().Add(24 * time.Second),
			}
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_fate_wheel",
				Payload: map[string]interface{}{
					"event":        "fate_wheel_perfect",
					"trigger_name": p.GetDisplayName(),
					"boost_mult":   11.0,
					"boost_secs":   24,
					"total_reward": totalReward,
				},
			})
			g.sendAnnounce(fmt.Sprintf("🎡✨ 命運完美！%s 三次 ≥20x！全服 ×11.0 加成 24 秒！", p.GetDisplayName()), "critical", "#FF8F00")
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_fate_wheel",
				Payload: map[string]interface{}{
					"event":        "fate_wheel_end",
					"trigger_name": p.GetDisplayName(),
					"total_reward": totalReward,
				},
			})
		}
	}()
	return true
}
