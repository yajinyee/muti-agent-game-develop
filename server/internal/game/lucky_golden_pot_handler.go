// lucky_golden_pot_handler.go — T224 幸運黃金鍋魚
// 設計：Gold Blitz™ Cash Collection 機制（Games Global 2026-05-28 最新）
//       收集金幣填滿黃金鍋（12格），每格有隨機金幣值（×5-×200）
//       填滿鍋子 → 額外 ×300 大獎 + 全服 ×43.0 加成 88 秒（新史上最高）
//       Enhanced Respin：每次新金幣落入重置 3 次機會，鎖定倍率不消失
//       觸發率：0.0004%（最稀有）；個人冷卻 330 秒；全服冷卻 390 秒
//       業界依據：Games Global「Fishin' Pots of Gold Gold Blitz Ultimate」（2026-05-28）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyGoldenPotManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
}

// goldenPotCoinEntry 黃金鍋金幣類型
type goldenPotCoinEntry struct {
	Name   string
	Mult   float64
	Weight int
}

var goldenPotCoinTable = []goldenPotCoinEntry{
	{Name: "Copper", Mult: 5.0, Weight: 40},
	{Name: "Silver", Mult: 20.0, Weight: 30},
	{Name: "Gold", Mult: 60.0, Weight: 18},
	{Name: "Platinum", Mult: 150.0, Weight: 9},
	{Name: "Diamond", Mult: 200.0, Weight: 3},
}

func newLuckyGoldenPotManager() *luckyGoldenPotManager {
	return &luckyGoldenPotManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyGoldenPotFish(defID string) bool {
	return defID == "T224"
}

func (m *luckyGoldenPotManager) tryLuckyGoldenPotFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(390 * time.Second)
	m.personalCD[p.ID] = now.Add(330 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_golden_pot",
		Payload: map[string]interface{}{
			"event":        "pot_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"grid_size":    12,
			"max_spins":    3,
		},
	})
	g.sendAnnounce(fmt.Sprintf("Gold Blitz! %s triggered Golden Pot! Cash Collection starts!", p.GetDisplayName()), "critical", "#FFD700")
	log.Printf("[LuckyGoldenPot] %s triggered Golden Pot fish", p.GetDisplayName())

	go func() {
		grid := make([]float64, 12) // 12 grid slots
		spinsLeft := 3
		totalMult := 0.0
		collectedCoins := []map[string]interface{}{}

		for spinsLeft > 0 {
			// 65% chance a coin lands this spin
			if rand.Float64() < 0.65 {
				emptySlots := []int{}
				for i, v := range grid {
					if v == 0 {
						emptySlots = append(emptySlots, i)
					}
				}
				if len(emptySlots) > 0 {
					slot := emptySlots[rand.Intn(len(emptySlots))]
					coin := rollGoldenPotCoin()
					grid[slot] = coin.Mult
					totalMult += coin.Mult
					spinsLeft = 3 // Enhanced Respin: reset on new coin

					collectedCoins = append(collectedCoins, map[string]interface{}{
						"slot":      slot,
						"coin_name": coin.Name,
						"coin_mult": coin.Mult,
					})

					g.broadcast(protocol.Envelope{
						Type: "lucky_golden_pot",
						Payload: map[string]interface{}{
							"event":      "coin_land",
							"slot":       slot,
							"coin_name":  coin.Name,
							"coin_mult":  coin.Mult,
							"total_mult": totalMult,
							"spins_left": spinsLeft,
						},
					})
				}
			}

			spinsLeft--
			if spinsLeft > 0 {
				g.broadcast(protocol.Envelope{
					Type: "lucky_golden_pot",
					Payload: map[string]interface{}{
						"event":      "spin_tick",
						"spins_left": spinsLeft,
					},
				})
			}
			time.Sleep(700 * time.Millisecond)
		}

		// Check full pot bonus (all 12 slots filled)
		fullPot := true
		for _, v := range grid {
			if v == 0 {
				fullPot = false
				break
			}
		}
		if fullPot {
			totalMult += 300.0
			g.broadcast(protocol.Envelope{
				Type: "lucky_golden_pot",
				Payload: map[string]interface{}{
					"event":      "full_pot_bonus",
					"bonus_mult": 300.0,
					"total_mult": totalMult,
				},
			})
			g.sendAnnounce("Full Pot! Extra x300.0 Gold Blitz bonus!", "critical", "#FF8C00")
		}

		betCost := float64(p.GetBetDef().BetCost)
		reward := int(totalMult * betCost)
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_golden_pot",
			Payload: map[string]interface{}{
				"event":           "pot_settle",
				"trigger_id":      p.ID,
				"trigger_name":    p.GetDisplayName(),
				"collected_coins": collectedCoins,
				"total_mult":      totalMult,
				"reward":          reward,
				"full_pot":        fullPot,
			},
		})

		// Global boost x43.0 for 88 seconds (new all-time high)
		g.broadcast(protocol.Envelope{
			Type: "lucky_golden_pot",
			Payload: map[string]interface{}{
				"event":        "global_boost",
				"global_mult":  43.0,
				"duration":     88,
				"trigger_name": p.GetDisplayName(),
				"total_mult":   totalMult,
			},
		})
		g.sendAnnounce(fmt.Sprintf("Gold Blitz settled! %s collected x%.0f! Global x43.0 for 88s! New all-time high!", p.GetDisplayName(), totalMult), "critical", "#FFD700")
		log.Printf("[LuckyGoldenPot] %s Golden Pot settled %.0fx, global x43.0 for 88s", p.GetDisplayName(), totalMult)
	}()

	return true
}

func rollGoldenPotCoin() goldenPotCoinEntry {
	total := 0
	for _, c := range goldenPotCoinTable {
		total += c.Weight
	}
	r := rand.Intn(total)
	for _, c := range goldenPotCoinTable {
		r -= c.Weight
		if r < 0 {
			return c
		}
	}
	return goldenPotCoinTable[0]
}
