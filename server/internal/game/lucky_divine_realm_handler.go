// lucky_divine_realm_handler.go — T179 幸運神域降臨魚
// 業界依據：「divine realm mechanic」
// 設計：神域降臨 30 秒，每 6 秒神域波（全場 HP -35%）
//       5 波全部命中 ≥6 個目標 → 神域完美：全服 ×14.0 加成 30 秒
//       個人冷卻 80 秒；全服冷卻 125 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyDivineRealmManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	active       *divineRealmSession
	perfectBoost *divineRealmPerfectBoost
}

type divineRealmSession struct {
	triggerID   string
	triggerName string
	waveCount   int
	waveHits    []int
	expiresAt   time.Time
}

type divineRealmPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyDivineRealmManager() *luckyDivineRealmManager {
	return &luckyDivineRealmManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyDivineRealmFish(defID string) bool {
	return defID == "T179"
}

func (m *luckyDivineRealmManager) getDivineRealmPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyDivineRealmManager) tryLuckyDivineRealmFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(80 * time.Second)
	m.globalCD = now.Add(125 * time.Second)
	m.active = &divineRealmSession{
		triggerID:   p.ID,
		triggerName: p.GetDisplayName(),
		waveHits:    make([]int, 0, 5),
		expiresAt:   now.Add(32 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_divine_realm",
		Payload: map[string]interface{}{
			"event":        "divine_realm_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"waves":        5,
			"wave_interval": 6,
			"hp_damage":    0.35,
		},
	})
	g.sendAnnounce(fmt.Sprintf("✨ 神域降臨！%s 召喚神域！5 波神域光柱！每波 HP -35%%！", p.GetDisplayName()), "critical", "#F9A825")
	log.Printf("[LuckyDivineRealm] %s 觸發神域降臨魚", p.GetDisplayName())

	go func() {
		for wave := 1; wave <= 5; wave++ {
			time.Sleep(6 * time.Second)
			// 全場 AOE 傷害（使用全場範圍）
			hits := g.applyAOEDamage(640, 360, 9999, 0.35)
			m.mu.Lock()
			if m.active != nil {
				m.active.waveHits = append(m.active.waveHits, hits)
				m.active.waveCount = wave
			}
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_divine_realm",
				Payload: map[string]interface{}{
					"event":        "divine_wave",
					"wave_no":      wave,
					"hits":         hits,
					"trigger_name": p.GetDisplayName(),
				},
			})
		}
		m.mu.Lock()
		sess := m.active
		m.active = nil
		m.mu.Unlock()
		if sess == nil {
			return
		}
		perfectWaves := 0
		for _, h := range sess.waveHits {
			if h >= 6 {
				perfectWaves++
			}
		}
		if perfectWaves >= 5 {
			m.mu.Lock()
			m.perfectBoost = &divineRealmPerfectBoost{
				mult:      14.0,
				expiresAt: time.Now().Add(30 * time.Second),
			}
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_divine_realm",
				Payload: map[string]interface{}{
					"event":        "divine_realm_perfect",
					"trigger_name": sess.triggerName,
					"boost_mult":   14.0,
					"boost_secs":   30,
				},
			})
			g.sendAnnounce(fmt.Sprintf("✨🏆 神域完美！%s 全服 ×14.0 加成 30 秒！", sess.triggerName), "critical", "#FFD600")
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_divine_realm",
				Payload: map[string]interface{}{
					"event":        "divine_realm_end",
					"perfect_waves": perfectWaves,
					"trigger_name": sess.triggerName,
				},
			})
		}
	}()
	return true
}
