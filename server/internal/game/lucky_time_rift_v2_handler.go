// lucky_time_rift_v2_handler.go — 幸運時空裂縫魚系統（DAY-291）
// 業界依據：Fishing Fortune 2026「Time Freeze mechanic — all fish freeze in place for 8 seconds,
//          allowing players to target high-value fish without them escaping」
//          業界原創「時空裂縫 + 全場凍結 + 凍結期間傷害加倍 + 裂縫爆炸」機制
//
// 設計：
//   - 擊破 T249 後，觸發「時空裂縫」：全場所有目標凍結 8 秒（停止移動）
//   - 凍結期間，所有目標受到的傷害 ×2.0（凍結加成）
//   - 凍結結束時，場上所有目標 HP -30%（裂縫爆炸）
//   - 若凍結期間玩家擊破 ≥ 5 個目標 → 「時空完美」：全服 ×2.5 加成 6 秒
//   - 個人冷卻 22 秒；全服冷卻 38 秒
//
// 設計差異：
//   - 與幸運冰凍世界（T237，全場 HP -40% 一次性）不同，時空裂縫是「凍結 + 傷害加倍 + 爆炸」三段
//   - 「凍結期間傷害 ×2.0」讓玩家有「要趁凍結期間瘋狂打魚，每一發都值兩發」的動力
//   - 「凍結結束 HP -30%」讓玩家有「凍結結束後還有一波爆炸，雙重收益」的驚喜感
//   - 「擊破 ≥ 5 個觸發時空完美」讓玩家有「要趁 8 秒內打完 5 條魚」的緊迫感
//   - 「全服廣播凍結狀態」讓所有玩家看到「現在全場凍結，快去打魚」的社交感
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyTimeRiftV2PersonalCD = 22 * time.Second // 個人冷卻
	LuckyTimeRiftV2GlobalCD   = 38 * time.Second // 全服冷卻

	// 時空裂縫設計
	TimeRiftV2FreezeDuration = 8 * time.Second // 凍結持續時間
	TimeRiftV2DamageMult     = 2.0             // 凍結期間傷害倍率
	TimeRiftV2ExplosionDmg   = 0.30            // 凍結結束爆炸 HP -30%
	TimeRiftV2PerfectKills   = 5               // 時空完美門檻（擊破數）

	// 時空完美：全服加成
	TimeRiftV2PerfectMult     = 2.5             // 全服 ×2.5
	TimeRiftV2PerfectDuration = 6 * time.Second // 持續 6 秒
)

// timeRiftV2PerfectBoost 時空完美全服加成
type timeRiftV2PerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

// timeRiftV2Session 時空裂縫會話
type timeRiftV2Session struct {
	triggerPlayerID   string
	triggerPlayerName string
	killCount         int       // 凍結期間擊破數
	expiresAt         time.Time // 凍結結束時間
	settled           bool
}

// luckyTimeRiftV2Manager 幸運時空裂縫魚管理器
type luckyTimeRiftV2Manager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前凍結會話（全服同時只有一個）
	activeSession *timeRiftV2Session

	// 時空完美全服加成
	perfectBoost *timeRiftV2PerfectBoost

	// 是否正在凍結
	isFrozen bool
}

// newLuckyTimeRiftV2Manager 建立管理器
func newLuckyTimeRiftV2Manager() *luckyTimeRiftV2Manager {
	return &luckyTimeRiftV2Manager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyTimeRiftV2Fish 判斷是否為幸運時空裂縫魚
func isLuckyTimeRiftV2Fish(defID string) bool {
	return defID == "T249"
}

// getTimeRiftV2PerfectMult 取得時空完美全服倍率（供 handleKill 使用）
func (m *luckyTimeRiftV2Manager) getTimeRiftV2PerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	m.perfectBoost = nil
	return 1.0
}

// getTimeRiftV2DamageMult 取得凍結期間傷害倍率（供 handleKill 使用）
func (m *luckyTimeRiftV2Manager) getTimeRiftV2DamageMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isFrozen && m.activeSession != nil && time.Now().Before(m.activeSession.expiresAt) {
		return TimeRiftV2DamageMult
	}
	return 1.0
}

// isTimeRiftV2Active 判斷時空裂縫是否正在凍結
func (m *luckyTimeRiftV2Manager) isTimeRiftV2Active() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isFrozen && m.activeSession != nil && time.Now().Before(m.activeSession.expiresAt)
}

// notifyTimeRiftV2Kill 凍結期間擊破通知（增加擊破計數）
func (m *luckyTimeRiftV2Manager) notifyTimeRiftV2Kill() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession != nil && !m.activeSession.settled {
		m.activeSession.killCount++
	}
}

// tryLuckyTimeRiftV2Fish 嘗試觸發時空裂縫
func (g *Game) tryLuckyTimeRiftV2Fish(p *player.Player) {
	m := g.LuckyTimeRiftV2
	m.mu.Lock()

	now := time.Now()

	// 個人冷卻檢查
	if cd, ok := m.personalCooldowns[p.ID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	// 全服冷卻檢查
	if now.Before(m.globalCooldownUntil) {
		m.mu.Unlock()
		return
	}
	// 已有活躍會話
	if m.activeSession != nil && !m.activeSession.settled {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyTimeRiftV2PersonalCD)
	m.globalCooldownUntil = now.Add(LuckyTimeRiftV2GlobalCD)

	// 建立會話
	session := &timeRiftV2Session{
		triggerPlayerID:   p.ID,
		triggerPlayerName: p.DisplayName,
		killCount:         0,
		expiresAt:         now.Add(TimeRiftV2FreezeDuration),
	}
	m.activeSession = session
	m.isFrozen = true
	m.mu.Unlock()

	log.Printf("[LuckyTimeRiftV2] %s 觸發時空裂縫！全場凍結 8 秒，傷害 ×2.0", p.DisplayName)

	// 廣播凍結開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTimeRiftV2,
		Payload: ws.LuckyTimeRiftV2Payload{
			Event:             "rift_start",
			TriggerPlayerID:   p.ID,
			TriggerPlayerName: p.DisplayName,
			FreezeDuration:    int(TimeRiftV2FreezeDuration.Seconds()),
			DamageMult:        TimeRiftV2DamageMult,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyTimeRiftV2, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⏸️ %s 觸發時空裂縫！全場凍結 8 秒！傷害 ×2.0！", p.DisplayName),
		"color":   "#00E5FF",
	})
	g.broadcastAnnouncement(ann)

	// 啟動凍結計時器
	go g.runTimeRiftV2Freeze(session)
}

// runTimeRiftV2Freeze 凍結計時器（goroutine）
func (g *Game) runTimeRiftV2Freeze(session *timeRiftV2Session) {
	time.Sleep(TimeRiftV2FreezeDuration)

	m := g.LuckyTimeRiftV2
	m.mu.Lock()
	if m.activeSession != session || session.settled {
		m.mu.Unlock()
		return
	}
	session.settled = true
	m.isFrozen = false
	killCount := session.killCount
	m.mu.Unlock()

	log.Printf("[LuckyTimeRiftV2] 凍結結束！%s 期間擊破 %d 個目標，執行裂縫爆炸", session.triggerPlayerName, killCount)

	// 裂縫爆炸：全場 HP -30%
	g.applyTimeRiftV2Explosion()

	// 廣播凍結結束
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTimeRiftV2,
		Payload: ws.LuckyTimeRiftV2Payload{
			Event:             "rift_end",
			TriggerPlayerID:   session.triggerPlayerID,
			TriggerPlayerName: session.triggerPlayerName,
			KillCount:         killCount,
			ExplosionDmg:      TimeRiftV2ExplosionDmg,
		},
	})

	// 判斷時空完美
	if killCount >= TimeRiftV2PerfectKills {
		g.doTimeRiftV2Perfect(session)
	}
}

// applyTimeRiftV2Explosion 裂縫爆炸：全場 HP -30%
func (g *Game) applyTimeRiftV2Explosion() {
	g.mu.Lock()

	type hpUpdate struct {
		instanceID string
		hp         int
		maxHP      int
	}
	var updates []hpUpdate

	for _, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		dmg := int(float64(t.MaxHP) * TimeRiftV2ExplosionDmg)
		if dmg < 1 {
			dmg = 1
		}
		t.HP -= dmg
		if t.HP < 0 {
			t.HP = 0
		}
		updates = append(updates, hpUpdate{
			instanceID: t.InstanceID,
			hp:         t.HP,
			maxHP:      t.MaxHP,
		})
	}
	g.mu.Unlock()

	// 廣播 HP 更新（在鎖外廣播）
	for _, u := range updates {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetUpdate,
			Payload: ws.TargetUpdatePayload{
				InstanceID: u.instanceID,
				HP:         u.hp,
				MaxHP:      u.maxHP,
			},
		})
	}
}

// doTimeRiftV2Perfect 時空完美：全服 ×2.5 加成 6 秒
func (g *Game) doTimeRiftV2Perfect(session *timeRiftV2Session) {
	m := g.LuckyTimeRiftV2
	m.mu.Lock()
	m.perfectBoost = &timeRiftV2PerfectBoost{
		mult:      TimeRiftV2PerfectMult,
		expiresAt: time.Now().Add(TimeRiftV2PerfectDuration),
	}
	m.mu.Unlock()

	log.Printf("[LuckyTimeRiftV2] 時空完美！%s 擊破 %d 個目標，全服 ×%.1f 加成 %v",
		session.triggerPlayerName, session.killCount, TimeRiftV2PerfectMult, TimeRiftV2PerfectDuration)

	// 廣播時空完美
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTimeRiftV2,
		Payload: ws.LuckyTimeRiftV2Payload{
			Event:             "rift_perfect",
			TriggerPlayerID:   session.triggerPlayerID,
			TriggerPlayerName: session.triggerPlayerName,
			KillCount:         session.killCount,
			PerfectMult:       TimeRiftV2PerfectMult,
			PerfectDuration:   int(TimeRiftV2PerfectDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyTimeRiftV2, session.triggerPlayerName, session.killCount, map[string]string{
		"message": fmt.Sprintf("⏸️ 時空完美！%s 凍結期間擊破 %d 個目標！全服 ×%.1f 加成 %d 秒！",
			session.triggerPlayerName, session.killCount, TimeRiftV2PerfectMult, int(TimeRiftV2PerfectDuration.Seconds())),
		"color": "#FFD700",
	})
	g.broadcastAnnouncement(ann)

	// 計時結束後廣播完美加成結束
	go func() {
		time.Sleep(TimeRiftV2PerfectDuration)
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyTimeRiftV2,
			Payload: ws.LuckyTimeRiftV2Payload{
				Event: "rift_perfect_end",
			},
		})
	}()
}
