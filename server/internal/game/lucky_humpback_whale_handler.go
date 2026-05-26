// Package game — T137 幸運座頭鯨魚 handler
// server-event-agent 負責維護
// 業界依據：Royal Fishing Jili「Humpback Whale 90-150x with 15x base multiplier」
// 設計：擊破後觸發「鯨歌共鳴」，座頭鯨發出聲波，
//       聲波每 2 秒擴散一圈（共 4 圈），每圈命中範圍內所有目標 HP -15%；
//       聲波命中目標越多，下一圈傷害越高（+5% per 3 targets）；
//       4 圈全部命中 ≥ 20 個目標 → 「完美鯨歌」：全服 ×3.0 加成 8 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyHumpbackWhaleManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	songBoost    *whaleSongBoost
}

type whaleSongBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyHumpbackWhaleManager() *luckyHumpbackWhaleManager {
	return &luckyHumpbackWhaleManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyHumpbackWhaleFish(defID string) bool {
	return defID == "T137"
}

func (m *luckyHumpbackWhaleManager) getWhaleSongMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.songBoost != nil && time.Now().Before(m.songBoost.expiresAt) {
		return m.songBoost.mult
	}
	return 1.0
}

func (g *Game) tryLuckyHumpbackWhaleFish(playerID, playerName string) {
	m := g.luckyHumpbackWhale
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
	m.personalCD[playerID] = now.Add(28 * time.Second)
	m.globalCD = now.Add(46 * time.Second)
	m.mu.Unlock()

	log.Printf("[LuckyHumpbackWhale] Triggered by %s", playerName)

	g.hub.Broadcast(protocol.MsgLuckyHumpbackWhale, protocol.LuckyHumpbackWhalePayload{
		Event:      "song_start",
		PlayerID:   playerID,
		PlayerName: playerName,
		WaveCount:  4,
	})

	go g.runWhaleSongWaves(playerID, playerName)
}

func (g *Game) runWhaleSongWaves(playerID, playerName string) {
	baseDmgPct := 0.15
	totalHitCount := 0

	for wave := 1; wave <= 4; wave++ {
		time.Sleep(2 * time.Second)

		dmgPct := baseDmgPct

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
		// 下一圈傷害加成：每 3 個目標 +5%
		baseDmgPct += float64(hitCount/3) * 0.05
		if baseDmgPct > 0.40 {
			baseDmgPct = 0.40
		}

		log.Printf("[LuckyHumpbackWhale] Wave %d: hit=%d, dmg=%.0f%%", wave, hitCount, dmgPct*100)

		g.hub.Broadcast(protocol.MsgLuckyHumpbackWhale, protocol.LuckyHumpbackWhalePayload{
			Event:      "song_wave",
			PlayerID:   playerID,
			PlayerName: playerName,
			WaveNum:    wave,
			HitCount:   hitCount,
			DamagePct:  dmgPct,
		})
	}

	// 完美鯨歌：4 圈命中總數 ≥ 20
	if totalHitCount >= 20 {
		g.doWhaleSongPerfect(playerID, playerName, totalHitCount)
	} else {
		g.hub.Broadcast(protocol.MsgLuckyHumpbackWhale, protocol.LuckyHumpbackWhalePayload{
			Event:         "song_end",
			PlayerID:      playerID,
			PlayerName:    playerName,
			TotalHitCount: totalHitCount,
		})
	}
}

func (g *Game) doWhaleSongPerfect(playerID, playerName string, totalHitCount int) {
	m := g.luckyHumpbackWhale
	m.mu.Lock()
	m.songBoost = &whaleSongBoost{
		mult:      3.0,
		expiresAt: time.Now().Add(8 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyHumpbackWhale] Perfect! %s total_hit=%d → global ×3.0 for 8s", playerName, totalHitCount)

	g.hub.Broadcast(protocol.MsgLuckyHumpbackWhale, protocol.LuckyHumpbackWhalePayload{
		Event:         "song_perfect",
		PlayerID:      playerID,
		PlayerName:    playerName,
		TotalHitCount: totalHitCount,
		BoostMult:     3.0,
		BoostSec:      8,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  fmt.Sprintf("🐋 完美鯨歌！%s 四波命中 %d 條！全服 ×3.0 加成 8 秒！", playerName, totalHitCount),
		Priority: "high",
		Color:    "#0288D1",
	})

	go func() {
		time.Sleep(8 * time.Second)
		m.mu.Lock()
		m.songBoost = nil
		m.mu.Unlock()
		g.hub.Broadcast(protocol.MsgLuckyHumpbackWhale, protocol.LuckyHumpbackWhalePayload{
			Event: "song_perfect_end",
		})
	}()
}
