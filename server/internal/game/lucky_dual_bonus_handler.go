// lucky_dual_bonus_handler.go — T222 Dual Bonus Fish
// Design: Dual Bonus mechanic (BGaming "Fishing Club 2", 2026-04)
//   After trigger, player chooses one of two bonuses:
//     Bonus A (Coin Collect): collect 5 coins, each x80.0, global x41.8 for 84s
//     Bonus B (Risk Wheel): spin risk wheel, max x500.0, global x42.0 for 85s
//   If no choice within 10s, auto-select Bonus A
//   Trigger rate: 0.0008%; personal CD 320s; global CD 380s
//   Industry ref: BGaming "Fishing Club 2" dual bonus games (2026-04)
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyDualBonusManager struct {
	globalCD      time.Time
	mu            sync.Mutex
	personalCD    map[string]time.Time
	pendingChoice map[string]*dualBonusChoice
}

type dualBonusChoice struct {
	PlayerID  string
	ChosenAt  time.Time
	BonusType string // "A" or "B"
	Resolved  bool
}

func newLuckyDualBonusManager() *luckyDualBonusManager {
	return &luckyDualBonusManager{
		personalCD:    make(map[string]time.Time),
		pendingChoice: make(map[string]*dualBonusChoice),
	}
}

func isLuckyDualBonusFish(defID string) bool {
	return defID == "T222"
}

func (m *luckyDualBonusManager) tryLuckyDualBonusFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(380 * time.Second)
	m.personalCD[p.ID] = now.Add(320 * time.Second)

	choice := &dualBonusChoice{
		PlayerID:  p.ID,
		ChosenAt:  now,
		BonusType: "A",
		Resolved:  false,
	}
	m.pendingChoice[p.ID] = choice
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_dual_bonus",
		Payload: map[string]interface{}{
			"event":        "choice_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"timeout":      10,
			"bonus_a":      "Coin Collect (5 coins x80.0)",
			"bonus_b":      "Risk Wheel (max x500.0)",
		},
	})
	g.sendAnnounce(fmt.Sprintf("Dual Bonus! %s triggered Dual Bonus Fish! Choose your bonus! 10s countdown!", p.GetDisplayName()), "critical", "#FFD700")
	log.Printf("[LuckyDualBonus] %s triggered Dual Bonus fish (10s choice)", p.GetDisplayName())

	go func() {
		time.Sleep(10 * time.Second)

		m.mu.Lock()
		c := m.pendingChoice[p.ID]
		bonusType := "A"
		if c != nil {
			bonusType = c.BonusType
		}
		delete(m.pendingChoice, p.ID)
		m.mu.Unlock()

		betCost := float64(p.GetBetDef().BetCost)

		switch bonusType {
		case "B":
			m.executeBonusB(g, p, betCost)
		default:
			m.executeBonusA(g, p, betCost)
		}
	}()

	return true
}

func (m *luckyDualBonusManager) executeBonusA(g *Game, p *Player, betCost float64) {
	// Bonus A: collect 5 coins, each x80.0
	totalMult := 0.0
	coins := []float64{}
	for i := 0; i < 5; i++ {
		coinMult := 80.0 * (1.0 + rand.Float64()*0.5) // 80-120x
		coins = append(coins, coinMult)
		totalMult += coinMult
		time.Sleep(600 * time.Millisecond)
		g.broadcast(protocol.Envelope{
			Type: "lucky_dual_bonus",
			Payload: map[string]interface{}{
				"event":      "coin_collect",
				"coin_index": i + 1,
				"coin_mult":  coinMult,
				"total_mult": totalMult,
			},
		})
	}

	reward := int(totalMult * betCost)
	g.mu.Lock()
	p.Coins += reward
	g.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_dual_bonus",
		Payload: map[string]interface{}{
			"event":        "bonus_a_settle",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"coins":        coins,
			"total_mult":   totalMult,
			"reward":       reward,
		},
	})

	// global boost via broadcast
	g.broadcast(protocol.Envelope{
		Type: "lucky_dual_bonus",
		Payload: map[string]interface{}{
			"event":        "global_boost",
			"global_mult":  41.8,
			"duration":     84,
			"trigger_name": p.GetDisplayName(),
			"bonus_type":   "A",
		},
	})
	g.sendAnnounce(fmt.Sprintf("Coin Collect! %s collected x%.0f! Global x41.8 for 84s!", p.GetDisplayName(), totalMult), "critical", "#FFD700")
	log.Printf("[LuckyDualBonus] %s Bonus A settled %.0fx, global x41.8 for 84s", p.GetDisplayName(), totalMult)
}

func (m *luckyDualBonusManager) executeBonusB(g *Game, p *Player, betCost float64) {
	// Risk Wheel: random result
	wheelResults := []float64{50.0, 100.0, 150.0, 200.0, 300.0, 500.0}
	weights := []int{30, 25, 20, 15, 8, 2}
	total := 0
	for _, w := range weights {
		total += w
	}
	r := rand.Intn(total)
	chosen := wheelResults[0]
	for i, w := range weights {
		r -= w
		if r < 0 {
			chosen = wheelResults[i]
			break
		}
	}

	reward := int(chosen * betCost)
	g.mu.Lock()
	p.Coins += reward
	g.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_dual_bonus",
		Payload: map[string]interface{}{
			"event":        "bonus_b_settle",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"wheel_mult":   chosen,
			"reward":       reward,
		},
	})

	// global boost via broadcast
	g.broadcast(protocol.Envelope{
		Type: "lucky_dual_bonus",
		Payload: map[string]interface{}{
			"event":        "global_boost",
			"global_mult":  42.0,
			"duration":     85,
			"trigger_name": p.GetDisplayName(),
			"bonus_type":   "B",
		},
	})
	g.sendAnnounce(fmt.Sprintf("Risk Wheel! %s spun x%.0f! Global x42.0 for 85s!", p.GetDisplayName(), chosen), "critical", "#FFD700")
	log.Printf("[LuckyDualBonus] %s Bonus B settled %.0fx, global x42.0 for 85s", p.GetDisplayName(), chosen)
}

// SetBonusChoice allows player to set their choice (called from game.go message handler)
func (m *luckyDualBonusManager) SetBonusChoice(playerID string, bonusType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c, ok := m.pendingChoice[playerID]; ok && !c.Resolved {
		if bonusType == "A" || bonusType == "B" {
			c.BonusType = bonusType
		}
	}
}
