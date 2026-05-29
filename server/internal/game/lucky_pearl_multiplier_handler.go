// lucky_pearl_multiplier_handler.go — T219 幸運珍珠倍率魚
// 設計：Pearl Multiplier 機制（BGaming Shark & Spark Hold & Win，2026-05 最新）
//       場上每個目標都有珍珠倍率（×1-×100），擊破時獎勵 = 基礎獎勵 × 珍珠倍率
//       全部珍珠收集 → 全服 ×40.0 加成 80 秒（新里程碑：全服 ×40.0，超越 T218 的 ×39.5）
//       觸發率：0.0012%（最稀有）；個人冷卻 300 秒；全服冷卻 360 秒
//       業界依據：BGaming「Shark & Spark Hold & Win」Pearl 倍率符號（2026-05）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyPearlMultiplierManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
}

func newLuckyPearlMultiplierManager() *luckyPearlMultiplierManager {
	return &luckyPearlMultiplierManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyPearlMultiplierFish(defID string) bool {
	return defID == "T219"
}

// 珍珠倍率權重（×1-×100）
var pearlMultWeights = []struct {
	Mult   float64
	Weight int
}{
	{1, 30},
	{2, 20},
	{5, 15},
	{10, 12},
	{20, 10},
	{30, 6},
	{50, 4},
	{80, 2},
	{100, 1},
}

func rollPearlMult() float64 {
	total := 0
	for _, w := range pearlMultWeights {
		total += w.Weight
	}
	r := rand.Intn(total)
	for _, w := range pearlMultWeights {
		r -= w.Weight
		if r < 0 {
			return w.Mult
		}
	}
	return 1.0
}

func (m *luckyPearlMultiplierManager) tryLuckyPearlMultiplierFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(360 * time.Second)
	m.personalCD[p.ID] = now.Add(300 * time.Second)
	m.mu.Unlock()

	// 為場上所有目標分配珍珠倍率
	pearlMults := g.assignPearlMultipliers()

	g.broadcast(protocol.Envelope{
		Type: "lucky_pearl_multiplier",
		Payload: map[string]interface{}{
			"event":        "pearl_assign",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"pearl_count":  len(pearlMults),
			"pearl_mults":  pearlMults,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🦪✨ 珍珠降臨！%s 觸發珍珠倍率魚！%d 個目標獲得珍珠倍率（最高 ×100）！", p.GetDisplayName(), len(pearlMults)), "critical", "#FFD700")
	log.Printf("[LuckyPearlMultiplier] %s 觸發珍珠倍率魚（%d 個目標獲得珍珠倍率）", p.GetDisplayName(), len(pearlMults))

	go func() {
		// 30 秒珍珠時間
		time.Sleep(30 * time.Second)

		// 統計收集到的珍珠
		collected := len(pearlMults)

		if collected >= 5 {
			globalBoostMult := 40.0
			globalBoostSecs := 80
			g.broadcast(protocol.Envelope{
				Type: "lucky_pearl_multiplier",
				Payload: map[string]interface{}{
					"event":             "pearl_perfect",
					"trigger_id":        p.ID,
					"trigger_name":      p.GetDisplayName(),
					"collected":         collected,
					"global_boost_mult": globalBoostMult,
					"global_boost_secs": globalBoostSecs,
				},
			})
			g.sendAnnounce(fmt.Sprintf("🦪🌟 珍珠完美收集！%d 個珍珠！全服 ×%.1f 加成 %d 秒！（新里程碑：全服 ×40.0）", collected, globalBoostMult, globalBoostSecs), "critical", "#FFD700")
			log.Printf("[LuckyPearlMultiplier] 珍珠完美收集！%d 個珍珠，全服 ×%.1f 加成 %d 秒（新里程碑：全服 ×40.0，超越 T218 的 ×39.5）", collected, globalBoostMult, globalBoostSecs)
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_pearl_multiplier",
				Payload: map[string]interface{}{
					"event":     "pearl_end",
					"collected": collected,
				},
			})
		}
	}()
	return true
}

// assignPearlMultipliers 為場上所有目標分配珍珠倍率
func (g *Game) assignPearlMultipliers() map[string]float64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	pearlMults := make(map[string]float64)
	for id, t := range g.targets {
		if t.HP <= 0 {
			continue
		}
		mult := rollPearlMult()
		pearlMults[id] = mult

		// 廣播珍珠倍率給 Client
		g.broadcast(protocol.Envelope{
			Type: "lucky_pearl_multiplier",
			Payload: map[string]interface{}{
				"event":       "pearl_on_target",
				"target_id":   id,
				"pearl_mult":  mult,
			},
		})
	}
	return pearlMults
}
