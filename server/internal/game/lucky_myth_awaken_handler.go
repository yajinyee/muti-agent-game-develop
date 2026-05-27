// lucky_myth_awaken_handler.go — T165 幸運神話覺醒魚
// 業界依據：Myth mode — 全場目標倍率 ×3.0，持續 25 秒，最終爆發
// 設計：擊破後觸發神話覺醒模式，全場所有目標倍率 ×3.0，持續 25 秒
//       25 秒內擊破 ≥15 個目標 → 神話完美：全服 ×8.0 加成 20 秒（遊戲最長加成）
//       個人冷卻 65 秒；全服冷卻 100 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyMythAwakenManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	sessions     map[string]*mythAwakenSession
	perfectBoost *mythAwakenPerfectBoost
	// 全局神話模式（影響所有玩家）
	mythActive    bool
	mythExpiresAt time.Time
	mythMult      float64
}

type mythAwakenSession struct {
	playerID  string
	killCount int
	expiresAt time.Time
}

type mythAwakenPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyMythAwakenManager() *luckyMythAwakenManager {
	return &luckyMythAwakenManager{
		personalCD: make(map[string]time.Time),
		sessions:   make(map[string]*mythAwakenSession),
	}
}

func isLuckyMythAwakenFish(defID string) bool {
	return defID == "T165"
}

func (m *luckyMythAwakenManager) getMythAwakenPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyMythAwakenManager) getMythMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.mythActive && time.Now().Before(m.mythExpiresAt) {
		return m.mythMult
	}
	return 1.0
}

func (m *luckyMythAwakenManager) onKillDuringMyth(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok || time.Now().After(sess.expiresAt) {
		return
	}
	sess.killCount++
}

func (m *luckyMythAwakenManager) tryLuckyMythAwakenFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(65 * time.Second)
	m.globalCD = now.Add(100 * time.Second)
	// 啟動全局神話模式
	m.mythActive = true
	m.mythExpiresAt = now.Add(25 * time.Second)
	m.mythMult = 3.0
	m.sessions[p.ID] = &mythAwakenSession{
		playerID:  p.ID,
		killCount: 0,
		expiresAt: now.Add(25 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyMythAwaken] Player %s triggered myth awaken mode (25s, ×3.0 all targets)", p.ID)

	g.hub.Broadcast(protocol.MsgLuckyMythAwaken, map[string]interface{}{
		"event":     "start",
		"player_id": p.ID,
		"duration":  25,
		"myth_mult": 3.0,
	})

	// 25 秒後結算
	go func() {
		time.Sleep(25 * time.Second)
		m.mu.Lock()
		m.mythActive = false
		sess, ok := m.sessions[p.ID]
		finalKills := 0
		if ok {
			finalKills = sess.killCount
			delete(m.sessions, p.ID)
		}
		m.mu.Unlock()

		// 完美神話判定
		if finalKills >= 15 {
			m.mu.Lock()
			m.perfectBoost = &mythAwakenPerfectBoost{
				mult:      8.0,
				expiresAt: time.Now().Add(20 * time.Second),
			}
			m.mu.Unlock()
			g.hub.Broadcast(protocol.MsgAnnounce, map[string]interface{}{
				"key":     "myth_awaken_perfect",
				"message": fmt.Sprintf("🌟 神話覺醒！%s 擊破 %d 個目標！全服 ×8.0 加成 20 秒！", p.ID, finalKills),
				"mult":    8.0,
				"duration": 20,
			})
		}

		g.hub.Broadcast(protocol.MsgLuckyMythAwaken, map[string]interface{}{
			"event":       "end",
			"player_id":   p.ID,
			"final_kills": finalKills,
			"perfect":     finalKills >= 15,
		})
	}()

	return true
}
