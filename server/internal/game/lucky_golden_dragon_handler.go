// lucky_golden_dragon_handler.go — T109 幸運黃金龍魚輪盤系統
// server-event-agent 負責維護
// 業界依據：Royal Fishing Jili「ChainLong King — dual-ring roulette, inner × outer = up to 350x」
// 設計：擊破 T109 後，觸發雙環輪盤（內環 × 外環 = 最終倍率，最高 350x）
// 個人冷卻 25 秒；全服冷卻 40 秒
package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"chiikawa-game/internal/protocol"
)

// 內環倍率（×基礎）
var innerRingWeights = []struct {
	Mult   float64
	Weight int
}{
	{2, 40},
	{3, 25},
	{5, 15},
	{7, 10},
	{10, 7},
	{14, 3},
}

// 外環倍率（×基礎）
var outerRingWeights = []struct {
	Mult   float64
	Weight int
}{
	{5, 35},
	{8, 25},
	{12, 18},
	{18, 12},
	{25, 7},
	{35, 3},
}

type luckyGoldenDragonManager struct {
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
}

func newLuckyGoldenDragonManager() *luckyGoldenDragonManager {
	return &luckyGoldenDragonManager{
		playerCooldowns: make(map[string]time.Time),
	}
}

func isLuckyGoldenDragonFish(defID string) bool {
	return defID == "T109"
}

func (m *luckyGoldenDragonManager) canTrigger(playerID string) bool {
	now := time.Now()
	if now.Before(m.globalCooldown) {
		return false
	}
	if cd, ok := m.playerCooldowns[playerID]; ok {
		if now.Before(cd) {
			return false
		}
	}
	return true
}

func rollInnerRing() float64 {
	total := 0
	for _, e := range innerRingWeights {
		total += e.Weight
	}
	r := rand.Intn(total)
	for _, e := range innerRingWeights {
		r -= e.Weight
		if r < 0 {
			return e.Mult
		}
	}
	return innerRingWeights[0].Mult
}

func rollOuterRing() float64 {
	total := 0
	for _, e := range outerRingWeights {
		total += e.Weight
	}
	r := rand.Intn(total)
	for _, e := range outerRingWeights {
		r -= e.Weight
		if r < 0 {
			return e.Mult
		}
	}
	return outerRingWeights[0].Mult
}

func (g *Game) tryLuckyGoldenDragon(playerID string, killerName string) {
	m := g.luckyGoldenDragon
	if !m.canTrigger(playerID) {
		return
	}

	now := time.Now()
	m.playerCooldowns[playerID] = now.Add(25 * time.Second)
	m.globalCooldown = now.Add(40 * time.Second)

	g.hub.Broadcast(protocol.MsgLuckyGoldenDragon, protocol.LuckyGoldenDragonPayload{
		Event:       "trigger",
		TriggerID:   playerID,
		TriggerName: killerName,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "🐉 " + killerName + " 觸發黃金龍魚輪盤！",
		Priority: "high",
		Color:    "#FFD700",
	})

	go g.runGoldenDragonSpin(playerID, killerName)
}

func (g *Game) runGoldenDragonSpin(playerID string, killerName string) {
	// 等待 1.5 秒讓 Client 顯示輪盤動畫
	time.Sleep(1500 * time.Millisecond)

	// 抽取輪盤結果
	innerMult := rollInnerRing()
	outerMult := rollOuterRing()
	finalMult := innerMult * outerMult

	// 廣播輪盤旋轉結果
	g.hub.Broadcast(protocol.MsgLuckyGoldenDragon, protocol.LuckyGoldenDragonPayload{
		Event:       "spin",
		TriggerID:   playerID,
		TriggerName: killerName,
		InnerMult:   innerMult,
		OuterMult:   outerMult,
		FinalMult:   finalMult,
	})

	// 等待 500ms 讓 Client 顯示結果
	time.Sleep(500 * time.Millisecond)

	g.mu.Lock()
	p, ok := g.players[playerID]
	reward := 0
	if ok {
		betCost := p.GetBetDef().BetCost
		reward = int(float64(betCost) * finalMult)
		p.AddCoins(reward)
		g.sendPlayerUpdate(playerID)
	}
	g.mu.Unlock()

	// 廣播結算
	g.hub.Broadcast(protocol.MsgLuckyGoldenDragon, protocol.LuckyGoldenDragonPayload{
		Event:       "result",
		TriggerID:   playerID,
		TriggerName: killerName,
		InnerMult:   innerMult,
		OuterMult:   outerMult,
		FinalMult:   finalMult,
		Reward:      reward,
	})

	// 高倍率全服廣播
	if finalMult >= 100 {
		g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
			Message:  "🐉🌟 " + killerName + " 黃金龍魚輪盤 ×" + formatMult(finalMult) + "！",
			Priority: "critical",
			Color:    "#FFD700",
		})
	}

	log.Printf("[GoldenDragon] Player %s: inner=%.0f, outer=%.0f, final=%.0f, reward=%d",
		playerID, innerMult, outerMult, finalMult, reward)
}

func formatMult(m float64) string {
	if m == float64(int(m)) {
		return fmt.Sprintf("%.0f", m)
	}
	return fmt.Sprintf("%.1f", m)
}
