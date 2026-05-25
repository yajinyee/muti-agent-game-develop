// lucky_vampire_mult_handler.go — T120 幸運吸血鬼魚
// server-event-agent 負責維護
// 業界依據：Jili Games 2026「Vampire multiplier increases with each kill, chance to enter ×5 multiplier mode」
package game

import (
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type vampireSession struct {
	playerID    string
	playerName  string
	absorbCount int     // 已吸收次數
	currentMult float64 // 當前倍率（1.0 → 5.0）
	inMultMode  bool    // 是否在倍率模式
	multExpiry  time.Time
	expiresAt   time.Time
	settled     bool
}

type luckyVampireMultManager struct {
	mu                sync.Mutex
	personalCooldowns map[string]time.Time
	globalCooldown    time.Time
	activeSessions    map[string]*vampireSession // playerID -> session
}

func newLuckyVampireMultManager() *luckyVampireMultManager {
	return &luckyVampireMultManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*vampireSession),
	}
}

func isLuckyVampireMultFish(defID string) bool {
	return defID == "T120"
}

// getVampireMult 供 handleKill 使用，取得吸血鬼倍率加成
func (m *luckyVampireMultManager) getVampireMult(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	if !ok || sess.settled {
		return 1.0
	}
	if sess.inMultMode && time.Now().Before(sess.multExpiry) {
		return sess.currentMult
	}
	return 1.0
}

// isVampireActive 判斷玩家是否在吸血模式
func (m *luckyVampireMultManager) isVampireActive(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	return ok && !sess.settled && time.Now().Before(sess.expiresAt)
}

// notifyVampireKill 玩家擊破目標時吸收倍率
func (m *luckyVampireMultManager) notifyVampireKill(g *Game, playerID string) {
	m.mu.Lock()
	sess, ok := m.activeSessions[playerID]
	if !ok || sess.settled || time.Now().After(sess.expiresAt) {
		m.mu.Unlock()
		return
	}

	sess.absorbCount++
	// 每次吸收 +0.5x，最高 5.0x
	if sess.currentMult < 5.0 {
		sess.currentMult += 0.5
		if sess.currentMult > 5.0 {
			sess.currentMult = 5.0
		}
	}

	// 吸收 ≥ 8 次 → 進入倍率模式（10 秒）
	if sess.absorbCount >= 8 && !sess.inMultMode {
		sess.inMultMode = true
		sess.multExpiry = time.Now().Add(10 * time.Second)
		playerName := sess.playerName
		currentMult := sess.currentMult
		m.mu.Unlock()

		g.hub.Broadcast(protocol.MsgLuckyVampireMult, protocol.LuckyVampireMultPayload{
			Event:       "mult_mode",
			TriggerID:   playerID,
			TriggerName: playerName,
			AbsorbCount: sess.absorbCount,
			CurrentMult: currentMult,
			TimeLeft:    10.0,
		})

		// 10 秒後倍率模式結束
		go func() {
			time.Sleep(10 * time.Second)
			g.hub.Broadcast(protocol.MsgLuckyVampireMult, protocol.LuckyVampireMultPayload{
				Event:       "mult_end",
				TriggerID:   playerID,
				TriggerName: playerName,
			})
		}()
		return
	}

	absorbCount := sess.absorbCount
	currentMult := sess.currentMult
	playerName := sess.playerName
	m.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyVampireMult, protocol.LuckyVampireMultPayload{
		Event:       "absorb",
		TriggerID:   playerID,
		TriggerName: playerName,
		AbsorbCount: absorbCount,
		CurrentMult: currentMult,
	})
}

func (m *luckyVampireMultManager) tryLuckyVampireMult(g *Game, playerID, playerName string) {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCooldown) {
		m.mu.Unlock()
		return
	}
	if cd, ok := m.personalCooldowns[playerID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	// 個人冷卻 18 秒，全服冷卻 28 秒
	m.personalCooldowns[playerID] = now.Add(18 * time.Second)
	m.globalCooldown = now.Add(28 * time.Second)

	sess := &vampireSession{
		playerID:    playerID,
		playerName:  playerName,
		absorbCount: 0,
		currentMult: 1.0,
		inMultMode:  false,
		expiresAt:   now.Add(20 * time.Second), // 20 秒吸血模式
	}
	m.activeSessions[playerID] = sess
	m.mu.Unlock()

	// 廣播觸發
	g.hub.Broadcast(protocol.MsgLuckyVampireMult, protocol.LuckyVampireMultPayload{
		Event:       "trigger",
		TriggerID:   playerID,
		TriggerName: playerName,
		CurrentMult: 1.0,
	})

	// 20 秒後結算
	go func() {
		time.Sleep(20 * time.Second)

		m.mu.Lock()
		sess, ok := m.activeSessions[playerID]
		if !ok || sess.settled {
			m.mu.Unlock()
			return
		}
		sess.settled = true
		absorbCount := sess.absorbCount
		currentMult := sess.currentMult
		m.mu.Unlock()

		g.hub.Broadcast(protocol.MsgLuckyVampireMult, protocol.LuckyVampireMultPayload{
			Event:       "settle",
			TriggerID:   playerID,
			TriggerName: playerName,
			AbsorbCount: absorbCount,
			CurrentMult: currentMult,
		})
	}()
}
