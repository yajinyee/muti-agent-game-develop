// lucky_time_bomb_handler.go — T162 幸運時間炸彈魚
// 業界依據：Time Bomb mechanic — 倒數計時炸彈，30 秒內擊破越多目標，爆炸威力越強
// 設計：擊破後觸發時間炸彈，30 秒倒數，每次擊破 +1 炸彈能量（最高 30）
//       倒數結束時爆炸：每點能量 HP -2%（最高 -60%），能量 ≥20 → 完美爆炸全服 ×6.0 加成 13 秒
//       個人冷卻 50 秒；全服冷卻 75 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyTimeBombManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	sessions   map[string]*timeBombSession
	perfectBoost *timeBombPerfectBoost
}

type timeBombSession struct {
	playerID  string
	energy    int
	expiresAt time.Time
}

type timeBombPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyTimeBombManager() *luckyTimeBombManager {
	return &luckyTimeBombManager{
		personalCD: make(map[string]time.Time),
		sessions:   make(map[string]*timeBombSession),
	}
}

func isLuckyTimeBombFish(defID string) bool {
	return defID == "T162"
}

func (m *luckyTimeBombManager) getTimeBombPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyTimeBombManager) addEnergy(playerID string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok || time.Now().After(sess.expiresAt) {
		return 0
	}
	sess.energy++
	if sess.energy > 30 {
		sess.energy = 30
	}
	return sess.energy
}

func (m *luckyTimeBombManager) tryLuckyTimeBombFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(50 * time.Second)
	m.globalCD = now.Add(75 * time.Second)
	m.sessions[p.ID] = &timeBombSession{
		playerID:  p.ID,
		energy:    0,
		expiresAt: now.Add(30 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyTimeBomb] Player %s triggered time bomb (30s countdown)", p.ID)

	g.hub.Broadcast(protocol.MsgLuckyTimeBomb, map[string]interface{}{
		"event":      "start",
		"player_id":  p.ID,
		"duration":   30,
		"max_energy": 30,
	})

	// 30 秒後爆炸
	go func() {
		time.Sleep(30 * time.Second)
		m.mu.Lock()
		sess, ok := m.sessions[p.ID]
		finalEnergy := 0
		if ok {
			finalEnergy = sess.energy
			delete(m.sessions, p.ID)
		}
		m.mu.Unlock()

		// 計算爆炸傷害
		dmgPct := float64(finalEnergy) * 0.02 // 每點能量 -2%，最高 -60%
		if dmgPct > 0.6 {
			dmgPct = 0.6
		}

		// 對全場目標造成傷害
		g.mu.Lock()
		hitCount := 0
		for _, t := range g.targets {
			if t.HP > 0 {
				dmg := int(float64(t.MaxHP) * dmgPct)
				t.HP -= dmg
				if t.HP < 0 {
					t.HP = 0
				}
				hitCount++
			}
		}
		g.mu.Unlock()

		// 完美爆炸判定
		if finalEnergy >= 20 {
			m.mu.Lock()
			m.perfectBoost = &timeBombPerfectBoost{
				mult:      6.0,
				expiresAt: time.Now().Add(13 * time.Second),
			}
			m.mu.Unlock()
			g.hub.Broadcast(protocol.MsgAnnounce, map[string]interface{}{
				"key":     "time_bomb_perfect",
				"message": fmt.Sprintf("💣 完美爆炸！%s 蓄積 %d 能量！全服 ×6.0 加成 13 秒！", p.ID, finalEnergy),
				"mult":    6.0,
				"duration": 13,
			})
		}

		g.hub.Broadcast(protocol.MsgLuckyTimeBomb, map[string]interface{}{
			"event":        "explode",
			"player_id":    p.ID,
			"final_energy": finalEnergy,
			"dmg_pct":      dmgPct,
			"hit_count":    hitCount,
			"perfect":      finalEnergy >= 20,
		})
	}()

	return true
}

// onKillDuringTimeBomb 擊破時增加炸彈能量
func (m *luckyTimeBombManager) onKillDuringTimeBomb(g *Game, p *Player) {
	energy := m.addEnergy(p.ID)
	if energy > 0 {
		g.hub.Broadcast(protocol.MsgLuckyTimeBomb, map[string]interface{}{
			"event":     "energy_update",
			"player_id": p.ID,
			"energy":    energy,
		})
	}
}
