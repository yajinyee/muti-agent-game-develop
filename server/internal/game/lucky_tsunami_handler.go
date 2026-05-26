// Package game — T135 幸運海嘯魚 handler
// server-event-agent 負責維護
// 業界依據：Fishing Fortune 2026「multiplier cascade system — consecutive rare catches within 90s」
//           + 業界原創「海嘯波浪 — 三波衝擊，每波傷害遞增，最後一波全場清場」
// 設計：擊破後觸發「海嘯預警」2 秒，然後三波海嘯依序衝擊：
//       第一波：全場 HP -20%（2 秒後）
//       第二波：全場 HP -30%（4 秒後）
//       第三波：全場 HP -40%（6 秒後）
//       三波全部命中 ≥ 5 個目標 → 「完美海嘯」：全服 ×3.2 加成 8 秒
//       個人冷卻 30 秒；全服冷卻 48 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyTsunamiManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *tsunamiPerfectBoost
}

type tsunamiPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyTsunamiManager() *luckyTsunamiManager {
	return &luckyTsunamiManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyTsunamiFish(defID string) bool {
	return defID == "T135"
}

func (m *luckyTsunamiManager) getTsunamiPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (g *Game) tryLuckyTsunamiFish(playerID, playerName string) {
	m := g.luckyTsunami
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
	m.personalCD[playerID] = now.Add(30 * time.Second)
	m.globalCD = now.Add(48 * time.Second)
	m.mu.Unlock()

	log.Printf("[LuckyTsunami] Triggered by %s", playerName)

	g.hub.Broadcast(protocol.MsgLuckyTsunami, protocol.LuckyTsunamiPayload{
		Event:      "tsunami_warning",
		PlayerID:   playerID,
		PlayerName: playerName,
		WaveCount:  3,
	})

	// 三波海嘯
	go g.runTsunamiWaves(playerID, playerName)
}

func (g *Game) runTsunamiWaves(playerID, playerName string) {
	waveDelays := []time.Duration{2 * time.Second, 4 * time.Second, 6 * time.Second}
	waveDamages := []float64{0.20, 0.30, 0.40}
	totalHitCount := 0

	for i, delay := range waveDelays {
		time.Sleep(delay - func() time.Duration {
			if i == 0 {
				return 0
			}
			return waveDelays[i-1]
		}())

		dmgPct := waveDamages[i]
		waveNum := i + 1

		g.mu.Lock()
		hitCount := 0
		for _, t := range g.targets {
			if t.HP <= 0 {
				continue
			}
			dmg := int(float64(t.MaxHP) * dmgPct)
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

		totalHitCount += hitCount
		log.Printf("[LuckyTsunami] Wave %d: hit=%d, dmg=%.0f%%", waveNum, hitCount, dmgPct*100)

		g.hub.Broadcast(protocol.MsgLuckyTsunami, protocol.LuckyTsunamiPayload{
			Event:      "wave_hit",
			PlayerID:   playerID,
			PlayerName: playerName,
			WaveNum:    waveNum,
			HitCount:   hitCount,
			DamagePct:  dmgPct,
		})
	}

	// 結算：三波命中總數 ≥ 5 → 完美海嘯
	if totalHitCount >= 5 {
		g.doTsunamiPerfect(playerID, playerName, totalHitCount)
	} else {
		g.hub.Broadcast(protocol.MsgLuckyTsunami, protocol.LuckyTsunamiPayload{
			Event:         "tsunami_end",
			PlayerID:      playerID,
			PlayerName:    playerName,
			TotalHitCount: totalHitCount,
		})
	}
}

func (g *Game) doTsunamiPerfect(playerID, playerName string, totalHitCount int) {
	m := g.luckyTsunami
	m.mu.Lock()
	m.perfectBoost = &tsunamiPerfectBoost{
		mult:      3.2,
		expiresAt: time.Now().Add(8 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyTsunami] Perfect! %s total_hit=%d → global ×3.2 for 8s", playerName, totalHitCount)

	g.hub.Broadcast(protocol.MsgLuckyTsunami, protocol.LuckyTsunamiPayload{
		Event:         "tsunami_perfect",
		PlayerID:      playerID,
		PlayerName:    playerName,
		TotalHitCount: totalHitCount,
		BoostMult:     3.2,
		BoostSec:      8,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  fmt.Sprintf("🌊 完美海嘯！%s 三波命中 %d 條魚！全服 ×3.2 加成 8 秒！", playerName, totalHitCount),
		Priority: "high",
		Color:    "#0288D1",
	})

	go func() {
		time.Sleep(8 * time.Second)
		m.mu.Lock()
		m.perfectBoost = nil
		m.mu.Unlock()
		g.hub.Broadcast(protocol.MsgLuckyTsunami, protocol.LuckyTsunamiPayload{
			Event: "tsunami_perfect_end",
		})
	}()
}
