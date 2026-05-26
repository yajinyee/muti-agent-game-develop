// Package game — T140 幸運品質魚 handler
// server-event-agent 負責維護
// 業界依據：Fishing Frenzy Chapter 3「Fish Quality tier system that raises the stakes
//           on every cast, adding more variation to each catch」
// 設計：擊破後觸發「品質鑑定」，隨機抽取品質等級（Common/Rare/Epic/Legendary）；
//       品質等級決定獎勵倍率：Common ×2.0 / Rare ×5.0 / Epic ×15.0 / Legendary ×50.0；
//       抽到 Legendary → 「傳說品質」：全服 ×5.0 加成 12 秒；
//       個人冷卻 20 秒；全服冷卻 35 秒
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyQualityFishManager struct {
	mu            sync.Mutex
	personalCD    map[string]time.Time
	globalCD      time.Time
	legendaryBoost *qualityLegendaryBoost
}

type qualityLegendaryBoost struct {
	mult      float64
	expiresAt time.Time
}

// 品質等級定義
type qualityTier struct {
	Name   string
	Mult   float64
	Weight int
}

var qualityTiers = []qualityTier{
	{"Common", 2.0, 50},
	{"Rare", 5.0, 30},
	{"Epic", 15.0, 15},
	{"Legendary", 50.0, 5},
}

func newLuckyQualityFishManager() *luckyQualityFishManager {
	return &luckyQualityFishManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyQualityFish(defID string) bool {
	return defID == "T140"
}

func (m *luckyQualityFishManager) getQualityLegendaryMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.legendaryBoost != nil && time.Now().Before(m.legendaryBoost.expiresAt) {
		return m.legendaryBoost.mult
	}
	return 1.0
}

func rollQualityTier() qualityTier {
	totalWeight := 0
	for _, t := range qualityTiers {
		totalWeight += t.Weight
	}
	r := rand.Intn(totalWeight)
	for _, t := range qualityTiers {
		r -= t.Weight
		if r < 0 {
			return t
		}
	}
	return qualityTiers[0]
}

func (g *Game) tryLuckyQualityFish(playerID, playerName string) {
	m := g.luckyQualityFish
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) {
		m.mu.Unlock()
		return
	}
	if cd, ok := m.personalCD[playerID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	m.personalCD[playerID] = now.Add(20 * time.Second)
	m.globalCD = now.Add(35 * time.Second)
	m.mu.Unlock()

	// 抽取品質等級
	tier := rollQualityTier()

	log.Printf("[LuckyQualityFish] Triggered by %s, tier=%s (×%.0f)", playerName, tier.Name, tier.Mult)

	// 計算獎勵
	g.mu.RLock()
	p, ok := g.players[playerID]
	g.mu.RUnlock()
	reward := 0
	if ok {
		betCost := p.GetBetDef().BetCost
		reward = int(float64(betCost) * tier.Mult * 3.0)
		g.mu.Lock()
		p.AddCoins(reward)
		g.mu.Unlock()
	}

	g.hub.Broadcast(protocol.MsgLuckyQualityFish, protocol.LuckyQualityFishPayload{
		Event:      "quality_result",
		PlayerID:   playerID,
		PlayerName: playerName,
		TierName:   tier.Name,
		TierMult:   tier.Mult,
		Reward:     reward,
	})

	// 傳說品質：全服加成
	if tier.Name == "Legendary" {
		g.doQualityLegendary(playerID, playerName, tier.Mult)
	}
}

func (g *Game) doQualityLegendary(playerID, playerName string, tierMult float64) {
	m := g.luckyQualityFish
	m.mu.Lock()
	m.legendaryBoost = &qualityLegendaryBoost{
		mult:      5.0,
		expiresAt: time.Now().Add(12 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyQualityFish] LEGENDARY! %s → global ×5.0 for 12s", playerName)

	g.hub.Broadcast(protocol.MsgLuckyQualityFish, protocol.LuckyQualityFishPayload{
		Event:      "legendary_boost",
		PlayerID:   playerID,
		PlayerName: playerName,
		TierName:   "Legendary",
		TierMult:   tierMult,
		BoostMult:  5.0,
		BoostSec:   12,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  fmt.Sprintf("✨ 傳說品質！%s 抽到 LEGENDARY！全服 ×5.0 加成 12 秒！", playerName),
		Priority: "critical",
		Color:    "#FFD700",
	})

	go func() {
		time.Sleep(12 * time.Second)
		m.mu.Lock()
		m.legendaryBoost = nil
		m.mu.Unlock()
		g.hub.Broadcast(protocol.MsgLuckyQualityFish, protocol.LuckyQualityFishPayload{
			Event: "legendary_end",
		})
	}()
}
