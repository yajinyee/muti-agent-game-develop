// lucky_combo_burst_handler.go — T161 幸運連擊爆發魚
// 業界依據：Combo multiplier system — 連續擊破累積 Combo，Combo ×10 → 全服爆發
// 設計：擊破後觸發 20 秒連擊模式，每次擊破 Combo +1（最高 ×15.0 倍率加成）
//       Combo ≥10 → 完美連擊：全服 ×5.5 加成 12 秒
//       個人冷卻 45 秒；全服冷卻 70 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyComboBurstManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	sessions   map[string]*comboBurstSession
	perfectBoost *comboBurstPerfectBoost
}

type comboBurstSession struct {
	playerID  string
	combo     int
	expiresAt time.Time
}

type comboBurstPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyComboBurstManager() *luckyComboBurstManager {
	return &luckyComboBurstManager{
		personalCD: make(map[string]time.Time),
		sessions:   make(map[string]*comboBurstSession),
	}
}

func isLuckyComboBurstFish(defID string) bool {
	return defID == "T161"
}

func (m *luckyComboBurstManager) getComboBurstPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyComboBurstManager) getComboMult(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok || time.Now().After(sess.expiresAt) {
		return 1.0
	}
	// 每個 Combo 加 0.5x，最高 ×15.0
	mult := 1.0 + float64(sess.combo)*0.5
	if mult > 15.0 {
		mult = 15.0
	}
	return mult
}

func (m *luckyComboBurstManager) addCombo(playerID string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok || time.Now().After(sess.expiresAt) {
		return 0
	}
	sess.combo++
	return sess.combo
}

func (m *luckyComboBurstManager) tryLuckyComboBurstFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(45 * time.Second)
	m.globalCD = now.Add(70 * time.Second)
	// 建立連擊 session
	m.sessions[p.ID] = &comboBurstSession{
		playerID:  p.ID,
		combo:     0,
		expiresAt: now.Add(20 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyComboBurst] Player %s triggered combo burst mode (20s)", p.ID)

	// 廣播連擊模式開始
	g.hub.Broadcast(protocol.MsgLuckyComboBurst, map[string]interface{}{
		"event":      "start",
		"player_id":  p.ID,
		"duration":   20,
		"max_combo":  10,
		"max_mult":   15.0,
	})

	// 20 秒後結算
	go func() {
		time.Sleep(20 * time.Second)
		m.mu.Lock()
		sess, ok := m.sessions[p.ID]
		finalCombo := 0
		if ok {
			finalCombo = sess.combo
			delete(m.sessions, p.ID)
		}
		m.mu.Unlock()

		// 完美連擊判定
		if finalCombo >= 10 {
			m.mu.Lock()
			m.perfectBoost = &comboBurstPerfectBoost{
				mult:      5.5,
				expiresAt: time.Now().Add(12 * time.Second),
			}
			m.mu.Unlock()
			g.hub.Broadcast(protocol.MsgAnnounce, map[string]interface{}{
				"key":     "combo_burst_perfect",
				"message": fmt.Sprintf("🔥 完美連擊！%s 達成 %d 連擊！全服 ×5.5 加成 12 秒！", p.ID, finalCombo),
				"mult":    5.5,
				"duration": 12,
			})
		}

		g.hub.Broadcast(protocol.MsgLuckyComboBurst, map[string]interface{}{
			"event":       "end",
			"player_id":   p.ID,
			"final_combo": finalCombo,
			"perfect":     finalCombo >= 10,
		})
	}()

	return true
}

// onKillDuringComboBurst 擊破時觸發連擊加成
func (m *luckyComboBurstManager) onKillDuringComboBurst(g *Game, p *Player) float64 {
	m.mu.Lock()
	sess, ok := m.sessions[p.ID]
	if !ok || time.Now().After(sess.expiresAt) {
		m.mu.Unlock()
		return 1.0
	}
	sess.combo++
	combo := sess.combo
	mult := 1.0 + float64(combo)*0.5
	if mult > 15.0 {
		mult = 15.0
	}
	m.mu.Unlock()

	// 廣播連擊更新
	g.hub.Broadcast(protocol.MsgLuckyComboBurst, map[string]interface{}{
		"event":     "combo_update",
		"player_id": p.ID,
		"combo":     combo,
		"mult":      mult,
	})

	return mult
}
