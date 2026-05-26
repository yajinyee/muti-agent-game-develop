// Package game — T136 幸運龍怒蓄積魚 v2 handler
// server-event-agent 負責維護
// 業界依據：Royal Fishing Jili「Dragon Wrath system accumulates with every shot fired.
//           Once the wrath meter fills, players unleash a massive meteorite attack across
//           the centre screen, simultaneously targeting multiple fish」
// 設計：擊破後觸發「龍怒蓄積」30 秒，每次射擊 +1 怒氣（最高 30 點）；
//       30 秒後或怒氣滿 → 自動爆發「龍怒隕石雨」；
//       隕石數量 = 怒氣值（最少 5 顆，最多 30 顆）；
//       每顆隕石 HP -45%，AOE r=120px；
//       怒氣值 ≥ 20 → 「完美龍怒」：全服 ×3.5 加成 8 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyDragonWrathV2Manager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	activeSessions map[string]*dragonWrathV2Session
	perfectBoost   *dragonWrathV2PerfectBoost
}

type dragonWrathV2PerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type dragonWrathV2Session struct {
	playerID   string
	playerName string
	wrathValue int
	expiresAt  time.Time
	settled    bool
}

func newLuckyDragonWrathV2Manager() *luckyDragonWrathV2Manager {
	return &luckyDragonWrathV2Manager{
		personalCD:     make(map[string]time.Time),
		activeSessions: make(map[string]*dragonWrathV2Session),
	}
}

func isLuckyDragonWrathV2Fish(defID string) bool {
	return defID == "T136"
}

func (m *luckyDragonWrathV2Manager) getDragonWrathV2PerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyDragonWrathV2Manager) isDragonWrathV2Active(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.activeSessions[playerID]
	return ok && !s.settled && time.Now().Before(s.expiresAt)
}

func (m *luckyDragonWrathV2Manager) addWrathV2(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.activeSessions[playerID]
	if !ok || s.settled {
		return
	}
	s.wrathValue++
	if s.wrathValue > 30 {
		s.wrathValue = 30
	}
}

func (g *Game) tryLuckyDragonWrathV2Fish(playerID, playerName string) {
	m := g.luckyDragonWrathV2
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
	if s, ok := m.activeSessions[playerID]; ok && !s.settled {
		m.mu.Unlock()
		return
	}
	m.personalCD[playerID] = now.Add(30 * time.Second)
	m.globalCD = now.Add(50 * time.Second)
	session := &dragonWrathV2Session{
		playerID:   playerID,
		playerName: playerName,
		wrathValue: 0,
		expiresAt:  now.Add(30 * time.Second),
		settled:    false,
	}
	m.activeSessions[playerID] = session
	m.mu.Unlock()

	log.Printf("[LuckyDragonWrathV2] Triggered by %s", playerName)

	g.hub.Broadcast(protocol.MsgLuckyDragonWrathV2, protocol.LuckyDragonWrathV2Payload{
		Event:      "wrath_start",
		PlayerID:   playerID,
		PlayerName: playerName,
		Duration:   30.0,
		MaxWrath:   30,
	})

	go func() {
		time.Sleep(30 * time.Second)
		m.mu.Lock()
		s, ok := m.activeSessions[playerID]
		if !ok || s.settled {
			m.mu.Unlock()
			return
		}
		s.settled = true
		wrathValue := s.wrathValue
		m.mu.Unlock()
		g.doWrathV2Explosion(playerID, playerName, wrathValue)
	}()
}

func (g *Game) doWrathV2Explosion(playerID, playerName string, wrathValue int) {
	meteorCount := wrathValue
	if meteorCount < 5 {
		meteorCount = 5
	}

	g.mu.Lock()
	hitCount := 0
	for _, t := range g.targets {
		if t.HP <= 0 {
			continue
		}
		dmg := int(float64(t.MaxHP) * 0.45)
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

	log.Printf("[LuckyDragonWrathV2] Explosion! wrath=%d meteors=%d hit=%d", wrathValue, meteorCount, hitCount)

	g.hub.Broadcast(protocol.MsgLuckyDragonWrathV2, protocol.LuckyDragonWrathV2Payload{
		Event:       "wrath_explode",
		PlayerID:    playerID,
		PlayerName:  playerName,
		WrathValue:  wrathValue,
		MeteorCount: meteorCount,
		HitCount:    hitCount,
	})

	if wrathValue >= 20 {
		g.doDragonWrathV2Perfect(playerID, playerName, wrathValue)
	} else {
		g.hub.Broadcast(protocol.MsgLuckyDragonWrathV2, protocol.LuckyDragonWrathV2Payload{
			Event:      "wrath_end",
			PlayerID:   playerID,
			PlayerName: playerName,
			WrathValue: wrathValue,
		})
	}
}

func (g *Game) doDragonWrathV2Perfect(playerID, playerName string, wrathValue int) {
	m := g.luckyDragonWrathV2
	m.mu.Lock()
	m.perfectBoost = &dragonWrathV2PerfectBoost{
		mult:      3.5,
		expiresAt: time.Now().Add(8 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyDragonWrathV2] Perfect! %s wrath=%d → global ×3.5 for 8s", playerName, wrathValue)

	g.hub.Broadcast(protocol.MsgLuckyDragonWrathV2, protocol.LuckyDragonWrathV2Payload{
		Event:      "wrath_perfect",
		PlayerID:   playerID,
		PlayerName: playerName,
		WrathValue: wrathValue,
		BoostMult:  3.5,
		BoostSec:   8,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  fmt.Sprintf("🐉 完美龍怒！%s 蓄積 %d 怒氣！全服 ×3.5 加成 8 秒！", playerName, wrathValue),
		Priority: "high",
		Color:    "#FF4500",
	})

	go func() {
		time.Sleep(8 * time.Second)
		m.mu.Lock()
		m.perfectBoost = nil
		m.mu.Unlock()
		g.hub.Broadcast(protocol.MsgLuckyDragonWrathV2, protocol.LuckyDragonWrathV2Payload{
			Event: "wrath_perfect_end",
		})
	}()
}
