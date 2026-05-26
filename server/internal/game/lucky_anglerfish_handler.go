// Package game — T132 幸運巨型安康魚 handler
// server-event-agent 負責維護
// 業界依據：Jili Games「Giant Anglerfish can shoot electricity to open treasure chests,
//           giant crocodiles awaken to hunt fish on the fish farm to accumulate big prizes」
// 設計：擊破後觸發「深海誘餌」，安康魚的發光誘餌吸引場上所有目標向中心聚集 5 秒，
//       聚集期間傷害 ×1.8；5 秒後觸發「電擊爆炸」：全場 HP -30%；
//       電擊命中 ≥ 8 個目標 → 「完美誘捕」：全服 ×2.8 加成 7 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyAnglerFishManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	lureBoost    *anglerLureBoost
	perfectBoost *anglerPerfectBoost
}

type anglerLureBoost struct {
	dmgMult   float64
	expiresAt time.Time
}

type anglerPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyAnglerFishManager() *luckyAnglerFishManager {
	return &luckyAnglerFishManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyAnglerFish(defID string) bool {
	return defID == "T132"
}

func (m *luckyAnglerFishManager) getAnglerDamageMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.lureBoost != nil && time.Now().Before(m.lureBoost.expiresAt) {
		return m.lureBoost.dmgMult
	}
	return 1.0
}

func (m *luckyAnglerFishManager) getAnglerPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (g *Game) tryLuckyAnglerFish(playerID, playerName string) {
	m := g.luckyAnglerFish
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
	m.personalCD[playerID] = now.Add(24 * time.Second)
	m.globalCD = now.Add(40 * time.Second)

	// 啟動誘餌傷害加成 5 秒
	m.lureBoost = &anglerLureBoost{
		dmgMult:   1.8,
		expiresAt: now.Add(5 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyAnglerFish] Triggered by %s", playerName)

	g.hub.Broadcast(protocol.MsgLuckyAnglerFish, protocol.LuckyAnglerFishPayload{
		Event:      "lure_start",
		PlayerID:   playerID,
		PlayerName: playerName,
		LureSec:    5,
		DamageMult: 1.8,
	})

	// 5 秒後觸發電擊爆炸
	go func() {
		time.Sleep(5 * time.Second)
		g.doAnglerExplosion(playerID, playerName)
	}()
}

func (g *Game) doAnglerExplosion(playerID, playerName string) {
	m := g.luckyAnglerFish
	m.mu.Lock()
	m.lureBoost = nil
	m.mu.Unlock()

	// 全場 HP -30%
	g.mu.Lock()
	hitCount := 0
	for _, t := range g.targets {
		if t.HP <= 0 {
			continue
		}
		dmg := int(float64(t.MaxHP) * 0.30)
		t.HP -= dmg
		if t.HP < 0 {
			t.HP = 0
		}
		hitCount++
		g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
			InstanceID: t.InstanceID,
			HP:         t.HP,
			MaxHP:      t.MaxHP,
			X:          float64(t.X),
			Y:          float64(t.Y),
		})
	}
	g.mu.Unlock()

	log.Printf("[LuckyAnglerFish] Explosion! hit=%d", hitCount)

	g.hub.Broadcast(protocol.MsgLuckyAnglerFish, protocol.LuckyAnglerFishPayload{
		Event:      "explosion",
		PlayerID:   playerID,
		PlayerName: playerName,
		HitCount:   hitCount,
	})

	// 完美誘捕：命中 ≥ 8 個
	if hitCount >= 8 {
		g.doAnglerPerfect(playerID, playerName, hitCount)
	} else {
		g.hub.Broadcast(protocol.MsgLuckyAnglerFish, protocol.LuckyAnglerFishPayload{
			Event:      "lure_end",
			PlayerID:   playerID,
			PlayerName: playerName,
			HitCount:   hitCount,
		})
	}
}

func (g *Game) doAnglerPerfect(playerID, playerName string, hitCount int) {
	m := g.luckyAnglerFish
	m.mu.Lock()
	m.perfectBoost = &anglerPerfectBoost{
		mult:      2.8,
		expiresAt: time.Now().Add(7 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyAnglerFish] Perfect! %s hit=%d → global ×2.8 for 7s", playerName, hitCount)

	g.hub.Broadcast(protocol.MsgLuckyAnglerFish, protocol.LuckyAnglerFishPayload{
		Event:      "perfect",
		PlayerID:   playerID,
		PlayerName: playerName,
		HitCount:   hitCount,
		BoostMult:  2.8,
		BoostSec:   7,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  fmt.Sprintf("🎣 完美誘捕！%s 電擊 %d 條魚！全服 ×2.8 加成 7 秒！", playerName, hitCount),
		Priority: "high",
		Color:    "#00E5FF",
	})

	go func() {
		time.Sleep(7 * time.Second)
		m.mu.Lock()
		m.perfectBoost = nil
		m.mu.Unlock()
		g.hub.Broadcast(protocol.MsgLuckyAnglerFish, protocol.LuckyAnglerFishPayload{
			Event: "perfect_end",
		})
	}()
}
