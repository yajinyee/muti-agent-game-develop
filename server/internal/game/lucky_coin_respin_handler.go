// lucky_coin_respin_handler.go — T223 幸運 Coin Respin 魚
// 設計：Coin Respin 機制（BGaming「Shark & Spark Hold & Win」，2026-05-28 最新）
//       Hold & Win 風格：空格盤面，金幣落下並鎖定，每次新金幣重置 3 次機會
//       金幣類型：Bronze ×10.0，Silver ×30.0，Gold ×80.0，Diamond ×200.0
//       填滿全盤（9格）→ 額外 ×500.0 大獎
//       全服 ×42.5 加成 86 秒（新史上最高，超越 T222 的 ×42.0）
//       觸發率：0.0005%（最稀有）；個人冷卻 325 秒；全服冷卻 385 秒
//       業界依據：BGaming「Shark & Spark Hold & Win」Coin Respin（2026-05-28）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyCoinRespinManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
}

// coinTypeEntry 金幣類型定義
type coinTypeEntry struct {
	Name   string
	Mult   float64
	Weight int
}

var coinTypeTable = []coinTypeEntry{
	{Name: "Bronze", Mult: 10.0, Weight: 50},
	{Name: "Silver", Mult: 30.0, Weight: 30},
	{Name: "Gold", Mult: 80.0, Weight: 15},
	{Name: "Diamond", Mult: 200.0, Weight: 5},
}

func newLuckyCoinRespinManager() *luckyCoinRespinManager {
	return &luckyCoinRespinManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyCoinRespinFish(defID string) bool {
	return defID == "T223"
}

func (m *luckyCoinRespinManager) tryLuckyCoinRespinFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(385 * time.Second)
	m.personalCD[p.ID] = now.Add(325 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_coin_respin",
		Payload: map[string]interface{}{
			"event":        "respin_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"grid_size":    9,
			"max_spins":    3,
		},
	})
	g.sendAnnounce(fmt.Sprintf("Coin Respin! %s triggered Coin Respin! Hold & Win starts!", p.GetDisplayName()), "critical", "#FFD700")
	log.Printf("[LuckyCoinRespin] %s triggered Coin Respin fish", p.GetDisplayName())

	go func() {
		grid := make([]float64, 9) // 9 grid slots, 0 = empty
		spinsLeft := 3
		totalMult := 0.0
		collectedCoins := []map[string]interface{}{}

		for spinsLeft > 0 {
			// 60% chance a coin lands this spin
			if rand.Float64() < 0.60 {
				emptySlots := []int{}
				for i, v := range grid {
					if v == 0 {
						emptySlots = append(emptySlots, i)
					}
				}
				if len(emptySlots) > 0 {
					slot := emptySlots[rand.Intn(len(emptySlots))]
					coin := rollCoinTypeEntry()
					grid[slot] = coin.Mult
					totalMult += coin.Mult
					spinsLeft = 3 // reset on new coin

					collectedCoins = append(collectedCoins, map[string]interface{}{
						"slot":      slot,
						"coin_name": coin.Name,
						"coin_mult": coin.Mult,
					})

					g.broadcast(protocol.Envelope{
						Type: "lucky_coin_respin",
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
					Type: "lucky_coin_respin",
					Payload: map[string]interface{}{
						"event":      "spin_tick",
						"spins_left": spinsLeft,
					},
				})
			}
			time.Sleep(800 * time.Millisecond)
		}

		// Check full board bonus
		fullBoard := true
		for _, v := range grid {
			if v == 0 {
				fullBoard = false
				break
			}
		}
		if fullBoard {
			totalMult += 500.0
			g.broadcast(protocol.Envelope{
				Type: "lucky_coin_respin",
				Payload: map[string]interface{}{
					"event":      "full_board_bonus",
					"bonus_mult": 500.0,
					"total_mult": totalMult,
				},
			})
			g.sendAnnounce("Full Board! Extra x500.0 bonus!", "critical", "#FF4500")
		}

		betCost := float64(p.GetBetDef().BetCost)
		reward := int(totalMult * betCost)
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_coin_respin",
			Payload: map[string]interface{}{
				"event":           "respin_settle",
				"trigger_id":      p.ID,
				"trigger_name":    p.GetDisplayName(),
				"collected_coins": collectedCoins,
				"total_mult":      totalMult,
				"reward":          reward,
				"full_board":      fullBoard,
			},
		})

		// Global boost x42.5 for 86 seconds (new all-time high)
		g.broadcast(protocol.Envelope{
			Type: "lucky_coin_respin",
			Payload: map[string]interface{}{
				"event":        "global_boost",
				"global_mult":  42.5,
				"duration":     86,
				"trigger_name": p.GetDisplayName(),
				"total_mult":   totalMult,
			},
		})
		g.sendAnnounce(fmt.Sprintf("Coin Respin settled! %s collected x%.0f! Global x42.5 for 86s! New all-time high!", p.GetDisplayName(), totalMult), "critical", "#FFD700")
		log.Printf("[LuckyCoinRespin] %s Coin Respin settled %.0fx, global x42.5 for 86s", p.GetDisplayName(), totalMult)
	}()

	return true
}

func rollCoinTypeEntry() coinTypeEntry {
	total := 0
	for _, c := range coinTypeTable {
		total += c.Weight
	}
	r := rand.Intn(total)
	for _, c := range coinTypeTable {
		r -= c.Weight
		if r < 0 {
			return c
		}
	}
	return coinTypeTable[0]
}
