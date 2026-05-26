// lucky_coop_fish_handler.go — T127 幸運全服合作魚系統
// server-event-agent 負責維護
// 業界依據：業界原創「全服合作機制 — 所有玩家一起貢獻傷害，達到目標觸發全服大獎」
// 設計：擊破 T127 後，觸發「全服合作挑戰」：
//   - 全服所有玩家在 20 秒內共同擊破目標，累積「合作點數」
//   - 每擊破一個目標 +1 點（BOSS +10 點）
//   - 達到目標點數（依在線玩家數動態調整）→ 觸發「全服大爆發」：全服 ×4.0 加成 8 秒
//   - 每個玩家的貢獻比例決定個人獎勵分配
//   - 個人冷卻 25 秒；全服冷卻 45 秒
package game

import (
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type coopSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	targetPoints      int
	currentPoints     int
	contributions     map[string]int // playerID -> 貢獻點數
	expiresAt         time.Time
	settled           bool
}

type luckyCoopFishManager struct {
	mu              sync.Mutex
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
	activeSession   *coopSession
	coopBoost       *coopBoostInfo
}

type coopBoostInfo struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyCoopFishManager() *luckyCoopFishManager {
	return &luckyCoopFishManager{
		playerCooldowns: make(map[string]time.Time),
	}
}

func isLuckyCoopFish(defID string) bool {
	return defID == "T127"
}

func (m *luckyCoopFishManager) getCoopBoostMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.coopBoost != nil && time.Now().Before(m.coopBoost.expiresAt) {
		return m.coopBoost.mult
	}
	return 1.0
}

func (m *luckyCoopFishManager) isCoopActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.activeSession != nil && !m.activeSession.settled && time.Now().Before(m.activeSession.expiresAt)
}

func (m *luckyCoopFishManager) addCoopPoint(playerID string, points int) (int, int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil || m.activeSession.settled {
		return 0, 0, false
	}
	m.activeSession.contributions[playerID] += points
	m.activeSession.currentPoints += points
	current := m.activeSession.currentPoints
	target := m.activeSession.targetPoints
	reached := current >= target
	if reached {
		m.activeSession.settled = true
	}
	return current, target, reached
}

func (g *Game) tryLuckyCoopFish(playerID string, killerName string) {
	m := g.luckyCoopFish
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCooldown) {
		m.mu.Unlock()
		return
	}
	if cd, ok := m.playerCooldowns[playerID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	m.playerCooldowns[playerID] = now.Add(25 * time.Second)
	m.globalCooldown = now.Add(45 * time.Second)

	// 依在線玩家數動態調整目標點數
	playerCount := len(g.players)
	if playerCount < 1 {
		playerCount = 1
	}
	targetPoints := 5 + playerCount*3 // 1人=8點, 2人=11點, 4人=17點

	session := &coopSession{
		triggerPlayerID:   playerID,
		triggerPlayerName: killerName,
		targetPoints:      targetPoints,
		currentPoints:     0,
		contributions:     make(map[string]int),
		expiresAt:         now.Add(20 * time.Second),
		settled:           false,
	}
	m.activeSession = session
	m.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyCoopFish, protocol.LuckyCoopFishPayload{
		Event:        "coop_start",
		TriggerID:    playerID,
		TriggerName:  killerName,
		TargetPoints: targetPoints,
		CurrentPoints: 0,
		TimeLeft:     20.0,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "🤝 " + killerName + " 發起全服合作！目標 " + string(rune('0'+targetPoints)) + " 點！20 秒！",
		Priority: "high",
		Color:    "#00E5FF",
	})

	// 超時處理
	go func() {
		time.Sleep(20 * time.Second)
		m.mu.Lock()
		if m.activeSession == session && !session.settled {
			session.settled = true
			m.mu.Unlock()
			g.hub.Broadcast(protocol.MsgLuckyCoopFish, protocol.LuckyCoopFishPayload{
				Event:         "coop_timeout",
				TriggerID:     playerID,
				TriggerName:   killerName,
				CurrentPoints: session.currentPoints,
				TargetPoints:  session.targetPoints,
			})
			g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
				Message:  "🤝 合作挑戰時間到！達成 " + string(rune('0'+session.currentPoints)) + "/" + string(rune('0'+session.targetPoints)) + " 點",
				Priority: "normal",
				Color:    "#888888",
			})
		} else {
			m.mu.Unlock()
		}
	}()
}

func (g *Game) notifyCoopKill(playerID string, killerName string, isBoss bool) {
	m := g.luckyCoopFish
	points := 1
	if isBoss {
		points = 10
	}
	current, target, reached := m.addCoopPoint(playerID, points)
	if current == 0 {
		return
	}

	g.hub.Broadcast(protocol.MsgLuckyCoopFish, protocol.LuckyCoopFishPayload{
		Event:         "coop_progress",
		TriggerID:     playerID,
		TriggerName:   killerName,
		CurrentPoints: current,
		TargetPoints:  target,
		TimeLeft:      float64(time.Until(m.activeSession.expiresAt).Seconds()),
	})

	if reached {
		g.doCoopSuccess(playerID, killerName)
	}
}

func (g *Game) doCoopSuccess(playerID string, killerName string) {
	m := g.luckyCoopFish
	m.mu.Lock()
	session := m.activeSession
	m.coopBoost = &coopBoostInfo{
		mult:      4.0,
		expiresAt: time.Now().Add(8 * time.Second),
	}
	m.mu.Unlock()

	// 依貢獻比例分配獎勵
	g.mu.Lock()
	totalContrib := session.currentPoints
	if totalContrib < 1 {
		totalContrib = 1
	}
	for pid, contrib := range session.contributions {
		p, ok := g.players[pid]
		if !ok {
			continue
		}
		betCost := p.GetBetDef().BetCost
		ratio := float64(contrib) / float64(totalContrib)
		reward := int(float64(betCost) * 50.0 * ratio)
		if reward < betCost {
			reward = betCost
		}
		p.AddCoins(reward)
		g.sendPlayerUpdate(pid)
	}
	g.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyCoopFish, protocol.LuckyCoopFishPayload{
		Event:         "coop_success",
		TriggerID:     playerID,
		TriggerName:   killerName,
		CurrentPoints: session.currentPoints,
		TargetPoints:  session.targetPoints,
		BoostMult:     4.0,
		BoostSecs:     8,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "🤝✨ 全服合作成功！全服 ×4.0 加成 8 秒！",
		Priority: "critical",
		Color:    "#00E5FF",
	})

	go func() {
		time.Sleep(8 * time.Second)
		g.hub.Broadcast(protocol.MsgLuckyCoopFish, protocol.LuckyCoopFishPayload{
			Event: "coop_boost_end",
		})
	}()

	log.Printf("[CoopFish] Coop success! Points: %d/%d", session.currentPoints, session.targetPoints)
}
