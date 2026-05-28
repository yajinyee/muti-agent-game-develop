// lucky_ice_fishing_wheel_handler.go — T214 幸運冰釣輪盤魚
// 設計：Ice Fishing Wheel 機制（Evolution Gaming「Ice Fishing」2026）
//       觸發後旋轉冰釣輪盤（5 個扇形：×100/×500/×1000/×2000/×5000）
//       連續 3 次旋轉，倍率相乘，最高 ×5000 × ×5000 × ×5000（理論最高）
//       實際最高：×5000（單次），全服 ×37.5 加成 75 秒
//       觸發率：0.003%；個人冷卻 275 秒；全服冷卻 335 秒
//       業界依據：Evolution Gaming「Ice Fishing」最高 5000x（2026）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyIceFishingWheelManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
}

// 冰釣輪盤扇形（5 個）
var iceFishingWheelSegments = []struct {
	Mult   float64
	Weight int
	Label  string
}{
	{100, 40, "🐟 ×100"},
	{500, 25, "🐠 ×500"},
	{1000, 15, "🦈 ×1000"},
	{2000, 10, "🐋 ×2000"},
	{5000, 10, "❄️ ×5000"},
}

func newLuckyIceFishingWheelManager() *luckyIceFishingWheelManager {
	return &luckyIceFishingWheelManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyIceFishingWheelFish(defID string) bool {
	return defID == "T214"
}

func (m *luckyIceFishingWheelManager) spinWheel() (float64, string) {
	totalWeight := 0
	for _, s := range iceFishingWheelSegments {
		totalWeight += s.Weight
	}
	r := int(rand.Float64() * float64(totalWeight))
	cumulative := 0
	for _, s := range iceFishingWheelSegments {
		cumulative += s.Weight
		if r < cumulative {
			return s.Mult, s.Label
		}
	}
	return iceFishingWheelSegments[0].Mult, iceFishingWheelSegments[0].Label
}

func (m *luckyIceFishingWheelManager) tryLuckyIceFishingWheelFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(275 * time.Second)
	m.globalCD = now.Add(335 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_ice_fishing_wheel",
		Payload: map[string]interface{}{
			"event":        "wheel_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"spins":        3,
		},
	})
	g.sendAnnounce(fmt.Sprintf("❄️🎡 冰釣輪盤！%s 觸發冰釣輪盤魚！3 次旋轉，最高 ×5000！", p.GetDisplayName()), "critical", "#00CED1")
	log.Printf("[LuckyIceFishingWheel] %s 觸發冰釣輪盤魚（3 次旋轉）", p.GetDisplayName())

	go func() {
		time.Sleep(800 * time.Millisecond)

		totalMult := 0.0
		maxSingleMult := 0.0
		spinResults := make([]float64, 0, 3)

		for spin := 1; spin <= 3; spin++ {
			time.Sleep(1500 * time.Millisecond)

			mult, label := m.spinWheel()
			totalMult += mult
			spinResults = append(spinResults, mult)
			if mult > maxSingleMult {
				maxSingleMult = mult
			}

			g.broadcast(protocol.Envelope{
				Type: "lucky_ice_fishing_wheel",
				Payload: map[string]interface{}{
					"event":       "spin_result",
					"spin":        spin,
					"mult":        mult,
					"label":       label,
					"total_mult":  totalMult,
				},
			})
			log.Printf("[LuckyIceFishingWheel] 第 %d 次旋轉：%s（累計 ×%.0f）", spin, label, totalMult)
		}

		// 對場上最高倍率目標施加最高單次倍率獎勵
		g.applyIceFishingWheelReward(maxSingleMult)

		if maxSingleMult >= 2000 {
			globalBoostMult := 37.5
			globalBoostSecs := 75
			g.broadcast(protocol.Envelope{
				Type: "lucky_ice_fishing_wheel",
				Payload: map[string]interface{}{
					"event":        "wheel_jackpot",
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
					"max_mult":     maxSingleMult,
					"total_mult":   totalMult,
					"global_mult":  globalBoostMult,
					"global_secs":  globalBoostSecs,
				},
			})
			g.sendAnnounce(fmt.Sprintf("❄️🌟 冰釣大獎！最高 ×%.0f！全服 ×%.1f 加成 %d 秒！", maxSingleMult, globalBoostMult, globalBoostSecs), "critical", "#00CED1")
			log.Printf("[LuckyIceFishingWheel] 冰釣大獎！最高 ×%.0f，全服 ×%.1f 加成 %d 秒", maxSingleMult, globalBoostMult, globalBoostSecs)
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_ice_fishing_wheel",
				Payload: map[string]interface{}{
					"event":      "wheel_end",
					"max_mult":   maxSingleMult,
					"total_mult": totalMult,
				},
			})
		}
	}()
	return true
}

// applyIceFishingWheelReward 對場上最高倍率目標施加冰釣輪盤獎勵
func (g *Game) applyIceFishingWheelReward(mult float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	var bestTarget *Target
	var bestID string
	for id, t := range g.targets {
		if t.HP <= 0 {
			continue
		}
		if bestTarget == nil || t.Def.MinMult > bestTarget.Def.MinMult {
			bestTarget = t
			bestID = id
		}
	}

	if bestTarget != nil {
		reward := int(float64(bestTarget.Def.MinMult) * mult)
		bestTarget.HP = 0
		delete(g.targets, bestID)

		g.broadcast(protocol.Envelope{
			Type: "lucky_ice_fishing_wheel",
			Payload: map[string]interface{}{
				"event":     "wheel_kill",
				"target_id": bestID,
				"reward":    reward,
				"mult":      mult,
			},
		})
	}
}
