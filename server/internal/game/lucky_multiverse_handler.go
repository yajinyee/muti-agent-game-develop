// lucky_multiverse_handler.go — T176 幸運多重宇宙魚
// 業界依據：Fishing Frenzy Chapter 3「parallel dimension mechanic」
// 設計：擊破後開啟 3 個平行宇宙，每個宇宙獨立擊破計數（目標 5 個）
//       全部 3 個宇宙完成 → 多重宇宙完美：全服 ×13.0 加成 28 秒
//       個人冷卻 75 秒；全服冷卻 120 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyMultiverseManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	active     *multiverseSession
	perfectBoost *multiversePerfectBoost
}

type multiverseSession struct {
	triggerID   string
	triggerName string
	universes   [3]int // 每個宇宙的擊破計數
	target      int    // 每個宇宙目標擊破數
	expiresAt   time.Time
}

type multiversePerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyMultiverseManager() *luckyMultiverseManager {
	return &luckyMultiverseManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyMultiverseFish(defID string) bool {
	return defID == "T176"
}

func (m *luckyMultiverseManager) getMultiversePerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyMultiverseManager) tryLuckyMultiverseFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(75 * time.Second)
	m.globalCD = now.Add(120 * time.Second)
	m.active = &multiverseSession{
		triggerID:   p.ID,
		triggerName: p.GetDisplayName(),
		target:      5,
		expiresAt:   now.Add(35 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_multiverse",
		Payload: map[string]interface{}{
			"event":        "multiverse_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"target":       5,
			"universes":    3,
			"duration":     35,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🌌 多重宇宙！%s 開啟 3 個平行宇宙！每個宇宙擊破 5 個目標！", p.GetDisplayName()), "critical", "#7B1FA2")
	log.Printf("[LuckyMultiverse] %s 觸發多重宇宙魚", p.GetDisplayName())

	go func() {
		time.Sleep(35 * time.Second)
		m.mu.Lock()
		sess := m.active
		m.active = nil
		m.mu.Unlock()
		if sess == nil {
			return
		}
		completed := 0
		for _, cnt := range sess.universes {
			if cnt >= sess.target {
				completed++
			}
		}
		if completed >= 3 {
			m.mu.Lock()
			m.perfectBoost = &multiversePerfectBoost{
				mult:      13.0,
				expiresAt: time.Now().Add(28 * time.Second),
			}
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_multiverse",
				Payload: map[string]interface{}{
					"event":        "multiverse_perfect",
					"trigger_name": sess.triggerName,
					"boost_mult":   13.0,
					"boost_secs":   28,
				},
			})
			g.sendAnnounce(fmt.Sprintf("🌌✨ 多重宇宙完美！%s 全服 ×13.0 加成 28 秒！", sess.triggerName), "critical", "#E040FB")
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_multiverse",
				Payload: map[string]interface{}{
					"event":        "multiverse_end",
					"completed":    completed,
					"trigger_name": sess.triggerName,
				},
			})
		}
	}()
	return true
}

func (m *luckyMultiverseManager) onKill(g *Game, p *Player) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.active == nil || time.Now().After(m.active.expiresAt) {
		return
	}
	// 輪流分配到三個宇宙
	minIdx := 0
	for i := 1; i < 3; i++ {
		if m.active.universes[i] < m.active.universes[minIdx] {
			minIdx = i
		}
	}
	if m.active.universes[minIdx] < m.active.target {
		m.active.universes[minIdx]++
	}
	completed := 0
	for _, cnt := range m.active.universes {
		if cnt >= m.active.target {
			completed++
		}
	}
	g.broadcast(protocol.Envelope{
		Type: "lucky_multiverse",
		Payload: map[string]interface{}{
			"event":      "universe_progress",
			"universes":  m.active.universes,
			"completed":  completed,
			"target":     m.active.target,
		},
	})
}
