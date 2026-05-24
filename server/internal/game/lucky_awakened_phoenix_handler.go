// lucky_awakened_phoenix_handler.go — T111 幸運覺醒鳳凰魚系統
// server-event-agent 負責維護
// 業界依據：Royal Fishing Jili「Rainbow Phoenix Power Up — 6x-10x boost for next attacks」
// 設計：擊破 T111 後，觸發「覺醒模式」：玩家下 5 次攻擊每次都有 Power Up 加成（6x-10x 隨機）
// 若 5 次全部命中目標（無空射）→「完美覺醒」：全服 ×2.0 加成 8 秒
// 個人冷卻 20 秒；全服冷卻 35 秒
package game

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

// awakenedPhoenixSession 覺醒鳳凰會話
type awakenedPhoenixSession struct {
	playerID    string
	playerName  string
	shotsLeft   int     // 剩餘 Power Up 次數（最多 5 次）
	hitCount    int     // 命中次數（非空射）
	totalReward int
	expiresAt   time.Time
	settled     bool
}

// awakenedPhoenixPerfectBoost 完美覺醒全服加成
type awakenedPhoenixPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

// luckyAwakenedPhoenixManager 管理覺醒鳳凰系統
type luckyAwakenedPhoenixManager struct {
	mu              sync.Mutex
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
	activeSessions  map[string]*awakenedPhoenixSession // playerID -> session
	perfectBoost    *awakenedPhoenixPerfectBoost
}

func newLuckyAwakenedPhoenixManager() *luckyAwakenedPhoenixManager {
	return &luckyAwakenedPhoenixManager{
		playerCooldowns: make(map[string]time.Time),
		activeSessions:  make(map[string]*awakenedPhoenixSession),
	}
}

// isLuckyAwakenedPhoenixFish 判斷是否為覺醒鳳凰魚
func isLuckyAwakenedPhoenixFish(defID string) bool {
	return defID == "T111"
}

// canTrigger 判斷是否可以觸發
func (m *luckyAwakenedPhoenixManager) canTrigger(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	if now.Before(m.globalCooldown) {
		return false
	}
	if cd, ok := m.playerCooldowns[playerID]; ok {
		if now.Before(cd) {
			return false
		}
	}
	return true
}

// getAwakenedPhoenixPerfectMult 取得完美覺醒全服倍率（供 handleKill 使用）
func (m *luckyAwakenedPhoenixManager) getAwakenedPhoenixPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost == nil {
		return 1.0
	}
	if time.Now().After(m.perfectBoost.expiresAt) {
		m.perfectBoost = nil
		return 1.0
	}
	return m.perfectBoost.mult
}

// isAwakenedPhoenixActive 判斷玩家是否在覺醒模式中
func (m *luckyAwakenedPhoenixManager) isAwakenedPhoenixActive(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	if !ok {
		return false
	}
	if time.Now().After(sess.expiresAt) || sess.shotsLeft <= 0 || sess.settled {
		return false
	}
	return true
}

// consumeAwakenedPhoenixShot 消耗一次 Power Up 機會，回傳本次加成倍率
// isHit: 是否命中目標（空射不算命中）
func (m *luckyAwakenedPhoenixManager) consumeAwakenedPhoenixShot(playerID string, isHit bool, betCost int) (float64, int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	if !ok || sess.settled || sess.shotsLeft <= 0 || time.Now().After(sess.expiresAt) {
		return 1.0, 0, false
	}

	// 隨機 Power Up 倍率 6x-10x
	powerUpMult := 6.0 + float64(rand.Intn(5)) // 6, 7, 8, 9, 10
	reward := int(float64(betCost) * powerUpMult)
	sess.totalReward += reward
	sess.shotsLeft--
	if isHit {
		sess.hitCount++
	}

	isDone := sess.shotsLeft <= 0 || time.Now().After(sess.expiresAt)
	return powerUpMult, reward, isDone
}

// tryLuckyAwakenedPhoenix 嘗試觸發覺醒鳳凰
func (g *Game) tryLuckyAwakenedPhoenix(playerID string, killerName string) {
	m := g.luckyAwakenedPhoenix
	if !m.canTrigger(playerID) {
		return
	}

	m.mu.Lock()
	now := time.Now()
	m.playerCooldowns[playerID] = now.Add(20 * time.Second)
	m.globalCooldown = now.Add(35 * time.Second)

	sess := &awakenedPhoenixSession{
		playerID:   playerID,
		playerName: killerName,
		shotsLeft:  5,
		hitCount:   0,
		expiresAt:  now.Add(30 * time.Second), // 30 秒內用完 5 次
	}
	m.activeSessions[playerID] = sess
	m.mu.Unlock()

	// 廣播觸發事件
	g.hub.Broadcast(protocol.MsgLuckyAwakenedPhoenix, protocol.LuckyAwakenedPhoenixPayload{
		Event:       "awaken_start",
		TriggerID:   playerID,
		TriggerName: killerName,
		ShotsLeft:   5,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "🔥 " + killerName + " 觸發覺醒鳳凰！下 5 次攻擊 Power Up 6x-10x！",
		Priority: "high",
		Color:    "#FF6B35",
	})

	log.Printf("[AwakenedPhoenix] Player %s triggered awakened phoenix mode", playerID)
}

// notifyAwakenedPhoenixShot 通知覺醒鳳凰射擊結果
func (g *Game) notifyAwakenedPhoenixShot(playerID string, killerName string, powerUpMult float64, reward int, isDone bool) {
	m := g.luckyAwakenedPhoenix
	m.mu.Lock()
	sess, ok := m.activeSessions[playerID]
	if !ok {
		m.mu.Unlock()
		return
	}
	shotsLeft := sess.shotsLeft
	hitCount := sess.hitCount
	totalReward := sess.totalReward
	m.mu.Unlock()

	// 廣播 Power Up 命中
	g.hub.Broadcast(protocol.MsgLuckyAwakenedPhoenix, protocol.LuckyAwakenedPhoenixPayload{
		Event:       "power_up",
		TriggerID:   playerID,
		TriggerName: killerName,
		PowerUpMult: powerUpMult,
		ShotsLeft:   shotsLeft,
		TotalReward: totalReward,
	})

	if isDone {
		g.settleAwakenedPhoenix(playerID, killerName, hitCount, totalReward)
	}
}

// settleAwakenedPhoenix 結算覺醒鳳凰
func (g *Game) settleAwakenedPhoenix(playerID string, killerName string, hitCount int, totalReward int) {
	m := g.luckyAwakenedPhoenix
	m.mu.Lock()
	sess, ok := m.activeSessions[playerID]
	if !ok || sess.settled {
		m.mu.Unlock()
		return
	}
	sess.settled = true
	m.mu.Unlock()

	// 完美覺醒：5 次全部命中
	isPerfect := hitCount >= 5
	if isPerfect {
		m.mu.Lock()
		m.perfectBoost = &awakenedPhoenixPerfectBoost{
			mult:      2.0,
			expiresAt: time.Now().Add(8 * time.Second),
		}
		m.mu.Unlock()

		g.hub.Broadcast(protocol.MsgLuckyAwakenedPhoenix, protocol.LuckyAwakenedPhoenixPayload{
			Event:       "perfect_awaken",
			TriggerID:   playerID,
			TriggerName: killerName,
			TotalReward: totalReward,
		})
		g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
			Message:  "🔥✨ 完美覺醒！" + killerName + " 全服 ×2.0 加成 8 秒！",
			Priority: "high",
			Color:    "#FFD700",
		})

		// 8 秒後廣播完美加成結束
		go func() {
			time.Sleep(8 * time.Second)
			g.hub.Broadcast(protocol.MsgLuckyAwakenedPhoenix, protocol.LuckyAwakenedPhoenixPayload{
				Event:       "perfect_end",
				TriggerID:   playerID,
				TriggerName: killerName,
			})
		}()
	}

	// 廣播結算
	g.hub.Broadcast(protocol.MsgLuckyAwakenedPhoenix, protocol.LuckyAwakenedPhoenixPayload{
		Event:       "awaken_end",
		TriggerID:   playerID,
		TriggerName: killerName,
		HitCount:    hitCount,
		TotalReward: totalReward,
	})

	log.Printf("[AwakenedPhoenix] Player %s settled: hits=%d, reward=%d, perfect=%v",
		playerID, hitCount, totalReward, isPerfect)
}
