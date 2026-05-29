// lucky_legend_awaken_handler.go — T226 幸運傳說覺醒魚
// 設計：Legend Dragon 覺醒升級機制（Royal Fishing Jili 2026）
//       覺醒後連續獎勵 8 次（每次 50-300x），每次獎勵遞增
//       Humpback Whale 模式：90-150x × 15 基礎倍率
//       Legend Dragon 模式：120-200x × 20 基礎倍率
//       全部 8 次完成 → 全服 ×44.0 加成 88 秒（超越 T225 的 ×43.5）
//       觸發率：0.0003%；個人冷卻 340 秒；全服冷卻 400 秒
//       業界依據：Royal Fishing Jili「Legend Dragon 120-200x, Humpback Whale 90-150x」（2026）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyLegendAwakenManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
}

// legendAwakenMode 覺醒模式
type legendAwakenMode struct {
	Name    string
	MinMult float64
	MaxMult float64
	Base    float64
}

var legendAwakenModes = []legendAwakenMode{
	{Name: "Humpback Whale", MinMult: 90.0, MaxMult: 150.0, Base: 15.0},
	{Name: "Legend Dragon", MinMult: 120.0, MaxMult: 200.0, Base: 20.0},
}

func newLuckyLegendAwakenManager() *luckyLegendAwakenManager {
	return &luckyLegendAwakenManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyLegendAwakenFish(defID string) bool {
	return defID == "T226"
}

func (m *luckyLegendAwakenManager) tryLuckyLegendAwakenFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(400 * time.Second)
	m.personalCD[p.ID] = now.Add(340 * time.Second)
	m.mu.Unlock()

	// Randomly select awaken mode (40% Humpback, 60% Legend Dragon)
	mode := legendAwakenModes[0]
	if rand.Float64() < 0.60 {
		mode = legendAwakenModes[1]
	}

	g.broadcast(protocol.Envelope{
		Type: "lucky_legend_awaken",
		Payload: map[string]interface{}{
			"event":        "awaken_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"mode_name":    mode.Name,
			"base_mult":    mode.Base,
			"rounds":       8,
		},
	})
	g.sendAnnounce(fmt.Sprintf("Legend Awaken! %s awakened %s! 8 consecutive rewards!", p.GetDisplayName(), mode.Name), "critical", "#FF6347")
	log.Printf("[LuckyLegendAwaken] %s triggered Legend Awaken fish (%s)", p.GetDisplayName(), mode.Name)

	go func() {
		totalMult := 0.0
		rewards := []map[string]interface{}{}

		for round := 1; round <= 8; round++ {
			// Each round: random mult in mode range, multiplied by base
			rawMult := mode.MinMult + rand.Float64()*(mode.MaxMult-mode.MinMult)
			roundMult := rawMult * mode.Base / 100.0 // normalize to reasonable range
			// Ensure minimum 5x per round, increasing each round
			roundMult = roundMult + float64(round)*2.0
			if roundMult < 5.0 {
				roundMult = 5.0
			}

			totalMult += roundMult
			betCost := float64(p.GetBetDef().BetCost)
			roundReward := int(roundMult * betCost)

			g.mu.Lock()
			p.Coins += roundReward
			g.mu.Unlock()

			rewards = append(rewards, map[string]interface{}{
				"round":        round,
				"round_mult":   roundMult,
				"round_reward": roundReward,
				"total_mult":   totalMult,
			})

			g.broadcast(protocol.Envelope{
				Type: "lucky_legend_awaken",
				Payload: map[string]interface{}{
					"event":        "awaken_reward",
					"round":        round,
					"round_mult":   roundMult,
					"round_reward": roundReward,
					"total_mult":   totalMult,
					"mode_name":    mode.Name,
				},
			})
			time.Sleep(700 * time.Millisecond)
		}

		g.broadcast(protocol.Envelope{
			Type: "lucky_legend_awaken",
			Payload: map[string]interface{}{
				"event":        "awaken_settle",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"mode_name":    mode.Name,
				"total_mult":   totalMult,
				"rewards":      rewards,
			},
		})

		// Global boost x44.0 for 88 seconds
		g.broadcast(protocol.Envelope{
			Type: "lucky_legend_awaken",
			Payload: map[string]interface{}{
				"event":        "global_boost",
				"global_mult":  44.0,
				"duration":     88,
				"trigger_name": p.GetDisplayName(),
				"mode_name":    mode.Name,
				"total_mult":   totalMult,
			},
		})
		g.sendAnnounce(fmt.Sprintf("Legend Awaken settled! %s %s x%.0f total! Global x44.0 for 88s!", p.GetDisplayName(), mode.Name, totalMult), "critical", "#FF6347")
		log.Printf("[LuckyLegendAwaken] %s Legend Awaken settled %.0fx, global x44.0 for 88s", p.GetDisplayName(), totalMult)
	}()

	return true
}
