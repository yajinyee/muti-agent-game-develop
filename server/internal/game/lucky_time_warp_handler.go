// lucky_time_warp_handler.go — T128 幸運時間扭曲魚系統
// server-event-agent 負責維護
// 業界依據：業界原創「時間扭曲 — 全場目標移動速度降低 70%，持續 10 秒，傷害 ×2.0，結束時全場爆炸 HP -20%」
// 設計：擊破 T128 後，觸發「時間扭曲」：
//   - 全場所有目標移動速度降低 70%（×0.3）
//   - 持續 10 秒，期間所有傷害 ×2.0
//   - 扭曲期間擊破 ≥ 6 個目標 → 「時間崩潰」：全服 ×2.5 加成 6 秒
//   - 扭曲結束時全場爆炸 HP -20%
//   - 個人冷卻 22 秒；全服冷卻 38 秒
package game

import (
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type timeWarpSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	killCount         int
	expiresAt         time.Time
	settled           bool
}

type luckyTimeWarpManager struct {
	mu              sync.Mutex
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
	activeSession   *timeWarpSession
	collapseBoost   *timeWarpCollapseBoost
}

type timeWarpCollapseBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyTimeWarpManager() *luckyTimeWarpManager {
	return &luckyTimeWarpManager{
		playerCooldowns: make(map[string]time.Time),
	}
}

func isLuckyTimeWarpFish(defID string) bool {
	return defID == "T128"
}

func (m *luckyTimeWarpManager) isTimeWarpActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.activeSession != nil && !m.activeSession.settled && time.Now().Before(m.activeSession.expiresAt)
}

func (m *luckyTimeWarpManager) getTimeWarpDamageMult() float64 {
	if m.isTimeWarpActive() {
		return 2.0
	}
	return 1.0
}

func (m *luckyTimeWarpManager) getCollapseBoostMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.collapseBoost != nil && time.Now().Before(m.collapseBoost.expiresAt) {
		return m.collapseBoost.mult
	}
	return 1.0
}

func (m *luckyTimeWarpManager) notifyWarpKill(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil || m.activeSession.settled {
		return
	}
	m.activeSession.killCount++
}

func (g *Game) tryLuckyTimeWarp(playerID string, killerName string) {
	m := g.luckyTimeWarp
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
	m.playerCooldowns[playerID] = now.Add(22 * time.Second)
	m.globalCooldown = now.Add(38 * time.Second)

	session := &timeWarpSession{
		triggerPlayerID:   playerID,
		triggerPlayerName: killerName,
		killCount:         0,
		expiresAt:         now.Add(10 * time.Second),
		settled:           false,
	}
	m.activeSession = session
	m.mu.Unlock()

	// 廣播時間扭曲開始（Client 端降低目標移動速度）
	g.hub.Broadcast(protocol.MsgLuckyTimeWarp, protocol.LuckyTimeWarpPayload{
		Event:       "warp_start",
		TriggerID:   playerID,
		TriggerName: killerName,
		Duration:    10.0,
		SpeedMult:   0.3,
		DamageMult:  2.0,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "⏰ " + killerName + " 觸發時間扭曲！全場慢速 10 秒！傷害 ×2.0！",
		Priority: "high",
		Color:    "#7B2FBE",
	})

	go g.runTimeWarpSession(session, playerID, killerName)
}

func (g *Game) runTimeWarpSession(session *timeWarpSession, playerID string, killerName string) {
	time.Sleep(10 * time.Second)

	g.luckyTimeWarp.mu.Lock()
	if session.settled {
		g.luckyTimeWarp.mu.Unlock()
		return
	}
	session.settled = true
	killCount := session.killCount
	g.luckyTimeWarp.mu.Unlock()

	// 扭曲結束：全場爆炸 HP -20%
	g.mu.Lock()
	for _, t := range g.targets {
		damage := t.MaxHP / 5
		t.HP -= damage
		if t.HP < 1 {
			t.HP = 1
		}
		g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
			InstanceID: t.InstanceID,
			HP:         t.HP,
			MaxHP:      t.MaxHP,
			X:          t.X,
			Y:          t.Y,
		})
	}
	g.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyTimeWarp, protocol.LuckyTimeWarpPayload{
		Event:       "warp_end",
		TriggerID:   playerID,
		TriggerName: killerName,
		KillCount:   killCount,
	})

	// 時間崩潰判定：扭曲期間擊破 ≥ 6 個
	if killCount >= 6 {
		g.luckyTimeWarp.mu.Lock()
		g.luckyTimeWarp.collapseBoost = &timeWarpCollapseBoost{
			mult:      2.5,
			expiresAt: time.Now().Add(6 * time.Second),
		}
		g.luckyTimeWarp.mu.Unlock()

		g.hub.Broadcast(protocol.MsgLuckyTimeWarp, protocol.LuckyTimeWarpPayload{
			Event:       "time_collapse",
			TriggerID:   playerID,
			TriggerName: killerName,
			KillCount:   killCount,
			BoostMult:   2.5,
			BoostSecs:   6,
		})
		g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
			Message:  "⏰💥 時間崩潰！" + killerName + " 扭曲期間擊破 " + string(rune('0'+killCount)) + " 條！全服 ×2.5 加成 6 秒！",
			Priority: "critical",
			Color:    "#7B2FBE",
		})

		go func() {
			time.Sleep(6 * time.Second)
			g.hub.Broadcast(protocol.MsgLuckyTimeWarp, protocol.LuckyTimeWarpPayload{
				Event: "collapse_end",
			})
		}()
	}

	log.Printf("[TimeWarp] Session ended. KillCount: %d", killCount)
}
