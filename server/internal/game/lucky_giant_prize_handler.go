// lucky_giant_prize_handler.go — T154 幸運巨型獎勵魚
// 業界依據：Jili「Giant Prize Fish lets you easily win great prizes, with the chance for 5x multipliers」
// 設計：擊破後觸發 5 次隨機大獎（每次 ×5.0-×50.0），平均 ≥20x → 完美大獎全服 ×4.5 加成 10 秒
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyGiantPrizeManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *giantPrizePerfectBoost
}

type giantPrizePerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

var giantPrizeWeights = []struct {
	mult   float64
	weight int
}{
	{5.0, 40},
	{10.0, 25},
	{20.0, 15},
	{30.0, 10},
	{40.0, 7},
	{50.0, 3},
}

func newLuckyGiantPrizeManager() *luckyGiantPrizeManager {
	return &luckyGiantPrizeManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyGiantPrizeFish(defID string) bool {
	return defID == "T154"
}

func (m *luckyGiantPrizeManager) getGiantPrizePerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyGiantPrizeManager) tryLuckyGiantPrizeFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(24 * time.Second)
	m.globalCD = now.Add(40 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_giant_prize",
		Payload: map[string]interface{}{
			"event":        "giant_prize_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"prize_count":  5,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🎁 %s 觸發巨型獎勵魚！5 次隨機大獎！", p.GetDisplayName()), "high", "#FFD700")
	log.Printf("[LuckyGiantPrize] %s 觸發巨型獎勵魚", p.GetDisplayName())

	go func() {
		totalReward := 0
		betCost := p.GetBetDef().BetCost
		for i := 1; i <= 5; i++ {
			time.Sleep(800 * time.Millisecond)
			// 隨機選倍率
			total := 0
			for _, w := range giantPrizeWeights {
				total += w.weight
			}
			r := rand.Intn(total)
			cum := 0
			mult := giantPrizeWeights[0].mult
			for _, w := range giantPrizeWeights {
				cum += w.weight
				if r < cum {
					mult = w.mult
					break
				}
			}
			reward := int(float64(betCost) * mult)
			g.mu.Lock()
			p.AddCoins(reward)
			g.mu.Unlock()
			totalReward += reward

			g.broadcast(protocol.Envelope{
				Type: "lucky_giant_prize",
				Payload: map[string]interface{}{
					"event":      "prize_drop",
					"trigger_id": p.ID,
					"prize_no":   i,
					"prize_mult": mult,
					"reward":     reward,
				},
			})
		}

		// 完美大獎判定（平均倍率 ≥20x）
		avgMult := float64(totalReward) / float64(betCost) / 5.0
		if avgMult >= 20.0 {
			m.mu.Lock()
			m.perfectBoost = &giantPrizePerfectBoost{
				mult:      4.5,
				expiresAt: time.Now().Add(10 * time.Second),
			}
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_giant_prize",
				Payload: map[string]interface{}{
					"event":        "giant_prize_perfect",
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
					"total_reward": totalReward,
					"boost_mult":   4.5,
					"boost_secs":   10,
				},
			})
			g.sendAnnounce(fmt.Sprintf("🎁✨ 完美大獎！%s 獲得 %d！全服 ×4.5 加成 10 秒！", p.GetDisplayName(), totalReward), "high", "#FFD700")
			time.AfterFunc(10*time.Second, func() {
				m.mu.Lock()
				m.perfectBoost = nil
				m.mu.Unlock()
				g.broadcast(protocol.Envelope{
					Type: "lucky_giant_prize",
					Payload: map[string]interface{}{
						"event":      "giant_prize_perfect_end",
						"trigger_id": p.ID,
					},
				})
			})
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_giant_prize",
				Payload: map[string]interface{}{
					"event":        "giant_prize_end",
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
					"total_reward": totalReward,
				},
			})
		}
	}()
	return true
}
