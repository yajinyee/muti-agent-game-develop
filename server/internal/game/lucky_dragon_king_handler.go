// lucky_dragon_king_handler.go — T196 幸運龍王輪盤魚
// 設計：雙環輪盤，內環 × 外環 = 最高 ×25.0
//       觸發後全服 ×23.0 加成 46 秒（超越 T195 的 ×22.0）
//       觸發率：0.018%；個人冷卻 145 秒；全服冷卻 210 秒
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyDragonKingManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	kingBoost  *dragonKingBoost
}

type dragonKingBoost struct {
	mult      float64
	expiresAt time.Time
}

// 龍王輪盤內環倍率權重
var dragonKingInnerWeights = []struct {
	Mult   float64
	Weight int
}{
	{2, 35},
	{5, 25},
	{8, 18},
	{12, 12},
	{18, 7},
	{25, 3},
}

// 龍王輪盤外環倍率權重
var dragonKingOuterWeights = []struct {
	Mult   float64
	Weight int
}{
	{3, 35},
	{6, 25},
	{10, 18},
	{15, 12},
	{20, 7},
	{25, 3},
}

func newLuckyDragonKingManager() *luckyDragonKingManager {
	return &luckyDragonKingManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyDragonKingFish(defID string) bool {
	return defID == "T196"
}

func (m *luckyDragonKingManager) getDragonKingMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.kingBoost != nil && time.Now().Before(m.kingBoost.expiresAt) {
		return m.kingBoost.mult
	}
	return 1.0
}

func weightedPickDragonKing(weights []struct {
	Mult   float64
	Weight int
}) float64 {
	total := 0
	for _, w := range weights {
		total += w.Weight
	}
	r := rand.Intn(total)
	cum := 0
	for _, w := range weights {
		cum += w.Weight
		if r < cum {
			return w.Mult
		}
	}
	return weights[0].Mult
}

func (m *luckyDragonKingManager) tryLuckyDragonKingFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(145 * time.Second)
	m.globalCD = now.Add(210 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_dragon_king",
		Payload: map[string]interface{}{
			"event":        "dragon_king_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("🐉👑 龍王輪盤！%s 召喚龍王！雙環輪盤啟動！最高 ×25.0！", p.GetDisplayName()), "critical", "#1A0A00")
	log.Printf("[LuckyDragonKing] %s 觸發龍王輪盤魚", p.GetDisplayName())

	go func() {
		time.Sleep(500 * time.Millisecond)

		// 雙環輪盤：內環 × 外環
		innerMult := weightedPickDragonKing(dragonKingInnerWeights)
		outerMult := weightedPickDragonKing(dragonKingOuterWeights)
		totalMult := innerMult * outerMult
		if totalMult > 25.0 {
			totalMult = 25.0
		}

		// 發放輪盤獎勵
		reward := int(float64(p.GetBetDef().BetCost) * totalMult)
		if reward < 1 {
			reward = 1
		}
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()
		g.sendPlayerUpdate(p.ID)

		// 觸發全服 ×23.0 加成 46 秒
		boostMult := 23.0
		boostSecs := 46
		m.mu.Lock()
		m.kingBoost = &dragonKingBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_dragon_king",
			Payload: map[string]interface{}{
				"event":        "dragon_king_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"inner_mult":   innerMult,
				"outer_mult":   outerMult,
				"total_mult":   totalMult,
				"reward":       reward,
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
			},
		})
		g.sendAnnounce(fmt.Sprintf("🐉🎰 龍王輪盤結果！%s 內環 ×%.0f × 外環 ×%.0f = ×%.1f！全服 ×%.1f 加成 %d 秒！",
			p.GetDisplayName(), innerMult, outerMult, totalMult, boostMult, boostSecs), "critical", "#2A1000")
	}()
	return true
}
