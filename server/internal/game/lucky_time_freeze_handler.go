// lucky_time_freeze_handler.go — T114 幸運時間凍結魚系統
// server-event-agent 負責維護
// 業界依據：Fishing Fortune 2026「Time Freeze mechanic — all fish freeze in place for 8 seconds,
//           allowing players to target high-value fish without them escaping」
// 設計：擊破 T114 後，全場所有目標凍結 8 秒（停止移動）；
//       凍結期間，所有目標受到的傷害 ×1.8（凍結加成）；
//       凍結結束時，場上所有目標 HP -25%（冰裂爆炸）；
//       若凍結期間玩家擊破 ≥ 4 個目標 → 「完美凍結」：全服 ×2.0 加成 5 秒；
//       個人冷卻 20 秒；全服冷卻 35 秒
package game

import (
	"log"
	"time"

	"chiikawa-game/internal/protocol"
)

// luckyTimeFreezeManager 管理時間凍結系統
type luckyTimeFreezeManager struct {
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time

	// 凍結狀態
	isFrozen      bool
	freezeExpires time.Time
	// 凍結期間擊破計數（per player）
	freezeKills map[string]int
	// 完美凍結全服加成
	perfectBoost *freezePerfectBoost
}

type freezePerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyTimeFreezeManager() *luckyTimeFreezeManager {
	return &luckyTimeFreezeManager{
		playerCooldowns: make(map[string]time.Time),
		freezeKills:     make(map[string]int),
	}
}

// isLuckyTimeFreezeFish 判斷是否為時間凍結魚
func isLuckyTimeFreezeFish(defID string) bool {
	return defID == "T114"
}

// isTimeFreezeActive 判斷是否凍結中
func (m *luckyTimeFreezeManager) isTimeFreezeActive() bool {
	return m.isFrozen && time.Now().Before(m.freezeExpires)
}

// getFreezeDamageMult 取得凍結期間傷害倍率
func (m *luckyTimeFreezeManager) getFreezeDamageMult() float64 {
	if m.isTimeFreezeActive() {
		return 1.8
	}
	return 1.0
}

// getFreezePerfectMult 取得完美凍結全服倍率（供 handleKill 使用）
func (m *luckyTimeFreezeManager) getFreezePerfectMult() float64 {
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

// notifyFreezeKill 凍結期間擊破計數
func (m *luckyTimeFreezeManager) notifyFreezeKill(playerID string) {
	if m.isTimeFreezeActive() {
		m.freezeKills[playerID]++
	}
}

func (m *luckyTimeFreezeManager) canTrigger(playerID string) bool {
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

// tryLuckyTimeFreeze 嘗試觸發時間凍結
func (g *Game) tryLuckyTimeFreeze(playerID string, killerName string) {
	m := g.luckyTimeFreeze
	if !m.canTrigger(playerID) {
		return
	}

	now := time.Now()
	m.playerCooldowns[playerID] = now.Add(20 * time.Second)
	m.globalCooldown = now.Add(35 * time.Second)

	// 設定凍結狀態
	m.isFrozen = true
	m.freezeExpires = now.Add(8 * time.Second)
	m.freezeKills = make(map[string]int)

	// 廣播凍結開始
	g.hub.Broadcast(protocol.MsgLuckyTimeFreeze, protocol.LuckyTimeFreezePayload{
		Event:       "freeze_start",
		TriggerID:   playerID,
		TriggerName: killerName,
		Duration:    8.0,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "❄️ " + killerName + " 觸發時間凍結！全場凍結 8 秒！",
		Priority: "high",
		Color:    "#00E5FF",
	})

	log.Printf("[TimeFreeze] Player %s triggered freeze", playerID)

	// 8 秒後凍結結束
	go g.runTimeFreezeEnd(playerID, killerName)
}

// runTimeFreezeEnd 凍結結束邏輯
func (g *Game) runTimeFreezeEnd(playerID string, killerName string) {
	time.Sleep(8 * time.Second)

	g.mu.Lock()
	m := g.luckyTimeFreeze
	m.isFrozen = false

	// 冰裂爆炸：全場 HP -25%
	for _, t := range g.targets {
		damage := int(float64(t.MaxHP) * 0.25)
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

	// 廣播凍結結束 + 冰裂爆炸
	g.hub.Broadcast(protocol.MsgLuckyTimeFreeze, protocol.LuckyTimeFreezePayload{
		Event:       "freeze_end",
		TriggerID:   playerID,
		TriggerName: killerName,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "❄️💥 冰裂爆炸！全場 HP -25%！",
		Priority: "normal",
		Color:    "#80D8FF",
	})

	// 完美凍結判定（擊破 ≥ 4 個）
	killCount := m.freezeKills[playerID]
	isPerfect := killCount >= 4

	if isPerfect {
		m.perfectBoost = &freezePerfectBoost{
			mult:      2.0,
			expiresAt: time.Now().Add(5 * time.Second),
		}
		g.hub.Broadcast(protocol.MsgLuckyTimeFreeze, protocol.LuckyTimeFreezePayload{
			Event:       "perfect_freeze",
			TriggerID:   playerID,
			TriggerName: killerName,
			KillCount:   killCount,
		})
		g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
			Message:  "❄️✨ 完美凍結！" + killerName + " 全服 ×2.0 加成 5 秒！",
			Priority: "high",
			Color:    "#00E5FF",
		})

		go func() {
			time.Sleep(5 * time.Second)
			g.mu.Lock()
			m.perfectBoost = nil
			g.mu.Unlock()
			g.hub.Broadcast(protocol.MsgLuckyTimeFreeze, protocol.LuckyTimeFreezePayload{
				Event:       "perfect_end",
				TriggerID:   playerID,
				TriggerName: killerName,
			})
		}()
	}
	g.mu.Unlock()

	log.Printf("[TimeFreeze] Freeze ended. Player %s kills=%d, perfect=%v", playerID, killCount, isPerfect)
}
