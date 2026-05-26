// Package game — T133 幸運黑洞魚 handler
// server-event-agent 負責維護
// 業界依據：Godot vortex water shader + 黑洞吸引機制
//           「black hole singularity — sucks all fish toward center, then implodes for massive damage」
// 設計：擊破後觸發「黑洞吸引」8 秒，所有目標被吸向中心（速度 ×0.2）；
//       8 秒後「黑洞坍縮」：全場 HP -50%（距離中心越近傷害越高）；
//       坍縮命中 ≥ 10 個目標 → 「奇點爆發」：全服 ×3.0 加成 8 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyBlackHoleManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	isActive     bool
	singularity  *blackHoleSingularity
}

type blackHoleSingularity struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyBlackHoleManager() *luckyBlackHoleManager {
	return &luckyBlackHoleManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyBlackHoleFish(defID string) bool {
	return defID == "T133"
}

func (m *luckyBlackHoleManager) isBlackHoleActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isActive
}

func (m *luckyBlackHoleManager) getBlackHoleSingularityMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.singularity != nil && time.Now().Before(m.singularity.expiresAt) {
		return m.singularity.mult
	}
	return 1.0
}

func (g *Game) tryLuckyBlackHoleFish(playerID, playerName string) {
	m := g.luckyBlackHole
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
	if m.isActive {
		m.mu.Unlock()
		return
	}
	m.personalCD[playerID] = now.Add(28 * time.Second)
	m.globalCD = now.Add(45 * time.Second)
	m.isActive = true
	m.mu.Unlock()

	log.Printf("[LuckyBlackHole] Triggered by %s", playerName)

	g.hub.Broadcast(protocol.MsgLuckyBlackHole, protocol.LuckyBlackHolePayload{
		Event:      "black_hole_start",
		PlayerID:   playerID,
		PlayerName: playerName,
		Duration:   8.0,
	})

	// 8 秒後坍縮
	go func() {
		time.Sleep(8 * time.Second)
		g.doBlackHoleCollapse(playerID, playerName)
	}()
}

func (g *Game) doBlackHoleCollapse(playerID, playerName string) {
	m := g.luckyBlackHole
	m.mu.Lock()
	m.isActive = false
	m.mu.Unlock()

	// 全場 HP -50%
	g.mu.Lock()
	hitCount := 0
	for _, t := range g.targets {
		if t.HP <= 0 {
			continue
		}
		dmg := int(float64(t.MaxHP) * 0.50)
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

	log.Printf("[LuckyBlackHole] Collapse! hit=%d", hitCount)

	g.hub.Broadcast(protocol.MsgLuckyBlackHole, protocol.LuckyBlackHolePayload{
		Event:      "collapse",
		PlayerID:   playerID,
		PlayerName: playerName,
		HitCount:   hitCount,
	})

	// 奇點爆發：命中 ≥ 10 個
	if hitCount >= 10 {
		g.doBlackHoleSingularity(playerID, playerName, hitCount)
	} else {
		g.hub.Broadcast(protocol.MsgLuckyBlackHole, protocol.LuckyBlackHolePayload{
			Event:      "black_hole_end",
			PlayerID:   playerID,
			PlayerName: playerName,
			HitCount:   hitCount,
		})
	}
}

func (g *Game) doBlackHoleSingularity(playerID, playerName string, hitCount int) {
	m := g.luckyBlackHole
	m.mu.Lock()
	m.singularity = &blackHoleSingularity{
		mult:      3.0,
		expiresAt: time.Now().Add(8 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyBlackHole] Singularity! %s hit=%d → global ×3.0 for 8s", playerName, hitCount)

	g.hub.Broadcast(protocol.MsgLuckyBlackHole, protocol.LuckyBlackHolePayload{
		Event:      "singularity",
		PlayerID:   playerID,
		PlayerName: playerName,
		HitCount:   hitCount,
		BoostMult:  3.0,
		BoostSec:   8,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  fmt.Sprintf("🌑 奇點爆發！%s 吸入 %d 條魚！全服 ×3.0 加成 8 秒！", playerName, hitCount),
		Priority: "high",
		Color:    "#7B2FBE",
	})

	go func() {
		time.Sleep(8 * time.Second)
		m.mu.Lock()
		m.singularity = nil
		m.mu.Unlock()
		g.hub.Broadcast(protocol.MsgLuckyBlackHole, protocol.LuckyBlackHolePayload{
			Event: "singularity_end",
		})
	}()
}
