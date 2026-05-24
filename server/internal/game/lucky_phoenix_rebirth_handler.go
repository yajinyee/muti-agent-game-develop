// lucky_phoenix_rebirth_handler.go — 幸運鳳凰涅槃魚系統（DAY-285）
// 業界依據：Royal Fishing Jili「Rainbow Phoenix Power Up + Awaken Boss」機制（2026）
// 業界原創「鳳凰涅槃+死亡重生+全場火焰洗禮+完全涅槃爆發」機制
//
// 設計：
//   - 擊破 T243 後，觸發「鳳凰涅槃」：
//     1. 全場火焰洗禮：場上所有目標 HP -25%
//     2. 涅槃重生：隨機選 3 個目標「涅槃重生」（HP 恢復 50%，但擊破倍率 ×4.0）
//     3. 若 3 個涅槃目標全部在 15 秒內被擊破 → 「鳳凰完全涅槃」：全服 ×3.0 加成 8 秒
//   - 全服廣播涅槃目標位置和結果
//   - 個人冷卻 28 秒；全服冷卻 45 秒
//
// 設計差異：
//   - 與龍怒隕石（T242，多點 AOE 墜落）不同，鳳凰涅槃是「先傷全場，再讓特定目標重生」
//   - 「涅槃目標 HP 恢復 50%」讓玩家有「這條魚又活了，但打死有 ×4.0！」的策略感
//   - 「3 個涅槃目標全部擊破觸發完全涅槃」讓玩家有「要趁 15 秒內打完 3 條涅槃魚」的緊迫感
//   - 「全服 ×3.0 加成 8 秒」是最高全服加成之一，製造「全服一起爽」的社交感
//   - 「全場火焰洗禮 HP -25%」讓所有魚都更容易打，製造「鳳凰降臨，全場魚都受傷」的爽感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyPhoenixRebirthPersonalCD = 28 * time.Second // 個人冷卻
	LuckyPhoenixRebirthGlobalCD   = 45 * time.Second // 全服冷卻

	// 鳳凰涅槃設計
	PhoenixFireBathHPDmg      = 0.25  // 全場火焰洗禮 HP -25%
	PhoenixRebirthHPRestore   = 0.50  // 涅槃目標 HP 恢復 50%
	PhoenixRebirthKillMult    = 4.0   // 涅槃目標擊破倍率 ×4.0
	PhoenixRebirthTargetCount = 3     // 涅槃目標數量
	PhoenixRebirthDuration    = 15 * time.Second // 涅槃持續時間

	// 鳳凰完全涅槃：全服加成
	PhoenixFullRebirthMult     = 3.0                   // 全服 ×3.0
	PhoenixFullRebirthDuration = 8 * time.Second       // 持續 8 秒
)

// phoenixRebirthTarget 涅槃重生目標
type phoenixRebirthTarget struct {
	instanceID string
	defID      string
	name       string
	x, y       float64
	killed     bool
}

// phoenixFullRebirthBoost 鳳凰完全涅槃全服加成
type phoenixFullRebirthBoost struct {
	mult      float64
	expiresAt time.Time
}

// phoenixRebirthSession 鳳凰涅槃會話
type phoenixRebirthSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	targets           []*phoenixRebirthTarget
	expiresAt         time.Time
	settled           bool
}

// luckyPhoenixRebirthManager 幸運鳳凰涅槃魚管理器
type luckyPhoenixRebirthManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前涅槃會話（同時只有一個）
	activeSession *phoenixRebirthSession

	// 鳳凰完全涅槃全服加成
	fullRebirthBoost *phoenixFullRebirthBoost
}

func newLuckyPhoenixRebirthManager() *luckyPhoenixRebirthManager {
	return &luckyPhoenixRebirthManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyPhoenixRebirthFish 判斷是否為幸運鳳凰涅槃魚
func isLuckyPhoenixRebirthFish(defID string) bool {
	return defID == "T243"
}

// getPhoenixFullRebirthMult 取得鳳凰完全涅槃全服加成倍率（供 handleKill 使用）
func (m *luckyPhoenixRebirthManager) getPhoenixFullRebirthMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.fullRebirthBoost != nil && time.Now().Before(m.fullRebirthBoost.expiresAt) {
		return m.fullRebirthBoost.mult
	}
	return 1.0
}

// isPhoenixRebirthTarget 判斷是否為涅槃重生目標（供 handleKill 使用）
// 回傳 (isRebirth bool, killMult float64)
func (m *luckyPhoenixRebirthManager) isPhoenixRebirthTarget(instanceID string) (bool, float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil || m.activeSession.settled {
		return false, 1.0
	}
	if time.Now().After(m.activeSession.expiresAt) {
		return false, 1.0
	}
	for _, t := range m.activeSession.targets {
		if t.instanceID == instanceID && !t.killed {
			return true, PhoenixRebirthKillMult
		}
	}
	return false, 1.0
}

// markPhoenixRebirthTargetKilled 標記涅槃目標已被擊破
// 回傳 (allKilled bool) — 是否所有涅槃目標都被擊破
func (m *luckyPhoenixRebirthManager) markPhoenixRebirthTargetKilled(instanceID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil || m.activeSession.settled {
		return false
	}
	for _, t := range m.activeSession.targets {
		if t.instanceID == instanceID {
			t.killed = true
			break
		}
	}
	// 檢查是否全部擊破
	for _, t := range m.activeSession.targets {
		if !t.killed {
			return false
		}
	}
	return true
}

// tryLuckyPhoenixRebirthFish 擊破 T243 後觸發鳳凰涅槃（供 handleKill 使用）
func (g *Game) tryLuckyPhoenixRebirthFish(p *player.Player) {
	mgr := g.LuckyPhoenixRebirth
	mgr.mu.Lock()

	now := time.Now()

	// 全服冷卻檢查
	if now.Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && now.Before(cd) {
		mgr.mu.Unlock()
		return
	}
	// 已有活躍會話
	if mgr.activeSession != nil && !mgr.activeSession.settled && now.Before(mgr.activeSession.expiresAt) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = now.Add(LuckyPhoenixRebirthPersonalCD)
	mgr.globalCooldownUntil = now.Add(LuckyPhoenixRebirthGlobalCD)
	mgr.mu.Unlock()

	// Step 1：全場火焰洗禮 HP -25%
	fireBathCount := g.applyPhoenixFireBath(PhoenixFireBathHPDmg)

	// Step 2：選取涅槃目標（3 個）
	rebirthTargets := g.selectPhoenixRebirthTargets(PhoenixRebirthTargetCount)

	// 建立涅槃會話
	session := &phoenixRebirthSession{
		triggerPlayerID:   p.ID,
		triggerPlayerName: p.DisplayName,
		targets:           rebirthTargets,
		expiresAt:         now.Add(PhoenixRebirthDuration),
	}

	mgr.mu.Lock()
	mgr.activeSession = session
	mgr.mu.Unlock()

	log.Printf("[PhoenixRebirth] player=%s fireBath=%d rebirthTargets=%d",
		p.ID, fireBathCount, len(rebirthTargets))

	// 建立廣播用的目標資訊
	targetInfos := make([]ws.PhoenixRebirthTargetInfo, 0, len(rebirthTargets))
	for _, t := range rebirthTargets {
		targetInfos = append(targetInfos, ws.PhoenixRebirthTargetInfo{
			InstanceID: t.instanceID,
			Name:       t.name,
			X:          t.x,
			Y:          t.y,
		})
	}

	// 全服廣播：鳳凰涅槃開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyPhoenixRebirth,
		Payload: ws.LuckyPhoenixRebirthPayload{
			Event:         "rebirth_start",
			PlayerID:      p.ID,
			PlayerName:    p.DisplayName,
			FireBathCount: fireBathCount,
			RebirthTargets: targetInfos,
			Duration:      int(PhoenixRebirthDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyPhoenixRebirth, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🔥🦅 %s 召喚鳳凰涅槃！全場火焰洗禮，%d 條魚涅槃重生（×%.0f）！",
			p.DisplayName, len(rebirthTargets), PhoenixRebirthKillMult),
		"color": "#FF6B35",
	})
	g.broadcastAnnouncement(ann)

	// 啟動涅槃計時器
	go g.runPhoenixRebirthTimer(session)
}

// applyPhoenixFireBath 全場火焰洗禮，對所有目標造成 HP -25%，回傳命中目標數
func (g *Game) applyPhoenixFireBath(dmgPct float64) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	hitCount := 0
	for _, t := range g.Targets {
		if t == nil || t.HP <= 0 {
			continue
		}
		dmg := int(float64(t.HP) * dmgPct)
		if dmg < 1 {
			dmg = 1
		}
		t.HP -= dmg
		if t.HP < 0 {
			t.HP = 0
		}
		hitCount++

		// 廣播目標 HP 更新
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetUpdate,
			Payload: ws.TargetUpdatePayload{
				InstanceID: t.InstanceID,
				HP:         t.HP,
				MaxHP:      t.MaxHP,
			},
		})
	}
	return hitCount
}

// selectPhoenixRebirthTargets 選取涅槃重生目標，並恢復其 HP
func (g *Game) selectPhoenixRebirthTargets(count int) []*phoenixRebirthTarget {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 收集存活目標
	candidates := make([]string, 0)
	for iid, t := range g.Targets {
		if t != nil && t.HP > 0 && t.DefID != "T243" {
			candidates = append(candidates, iid)
		}
	}

	// 隨機選取
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	selected := count
	if selected > len(candidates) {
		selected = len(candidates)
	}

	result := make([]*phoenixRebirthTarget, 0, selected)
	for i := 0; i < selected; i++ {
		iid := candidates[i]
		t := g.Targets[iid]
		if t == nil {
			continue
		}

		// 恢復 HP 50%
		restore := int(float64(t.MaxHP) * PhoenixRebirthHPRestore)
		t.HP += restore
		if t.HP > t.MaxHP {
			t.HP = t.MaxHP
		}

		// 廣播 HP 恢復
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetUpdate,
			Payload: ws.TargetUpdatePayload{
				InstanceID: t.InstanceID,
				HP:         t.HP,
				MaxHP:      t.MaxHP,
				IsRebirth:  true,
			},
		})

		result = append(result, &phoenixRebirthTarget{
			instanceID: t.InstanceID,
			defID:      t.DefID,
			name:       t.Def.Name,
			x:          t.X,
			y:          t.Y,
		})
	}
	return result
}

// notifyPhoenixRebirthKill 涅槃目標被擊破時呼叫（由 handleKill 呼叫）
func (g *Game) notifyPhoenixRebirthKill(p *player.Player, instanceID string) {
	mgr := g.LuckyPhoenixRebirth
	allKilled := mgr.markPhoenixRebirthTargetKilled(instanceID)

	// 廣播涅槃目標被擊破
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyPhoenixRebirth,
		Payload: ws.LuckyPhoenixRebirthPayload{
			Event:      "rebirth_kill",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			InstanceID: instanceID,
			KillMult:   PhoenixRebirthKillMult,
		},
	})

	log.Printf("[PhoenixRebirth] rebirth target killed player=%s instanceID=%s allKilled=%v",
		p.ID, instanceID, allKilled)

	// 所有涅槃目標都被擊破 → 鳳凰完全涅槃
	if allKilled {
		mgr.mu.Lock()
		if mgr.activeSession != nil {
			mgr.activeSession.settled = true
		}
		mgr.mu.Unlock()
		g.doPhoenixFullRebirth(p)
	}
}

// doPhoenixFullRebirth 鳳凰完全涅槃：全服 ×3.0 加成 8 秒
func (g *Game) doPhoenixFullRebirth(p *player.Player) {
	mgr := g.LuckyPhoenixRebirth
	mgr.mu.Lock()
	mgr.fullRebirthBoost = &phoenixFullRebirthBoost{
		mult:      PhoenixFullRebirthMult,
		expiresAt: time.Now().Add(PhoenixFullRebirthDuration),
	}
	mgr.mu.Unlock()

	log.Printf("[PhoenixRebirth] FULL REBIRTH! player=%s global x%.1f for %v",
		p.ID, PhoenixFullRebirthMult, PhoenixFullRebirthDuration)

	// 全服廣播鳳凰完全涅槃
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyPhoenixRebirth,
		Payload: ws.LuckyPhoenixRebirthPayload{
			Event:        "rebirth_full",
			PlayerName:   p.DisplayName,
			FullRebirthMult: PhoenixFullRebirthMult,
			Duration:     int(PhoenixFullRebirthDuration.Seconds()),
		},
	})

	// 全服最高優先公告
	ann := g.Announce.Create(announce.EventLuckyPhoenixRebirth, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🔥🦅🔥 %s 鳳凰完全涅槃！全服 ×%.0f 加成 %d 秒！",
			p.DisplayName, PhoenixFullRebirthMult, int(PhoenixFullRebirthDuration.Seconds())),
		"color": "#FF4500",
	})
	g.broadcastAnnouncement(ann)

	// 8 秒後廣播完全涅槃結束
	go func() {
		time.Sleep(PhoenixFullRebirthDuration)
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyPhoenixRebirth,
			Payload: ws.LuckyPhoenixRebirthPayload{
				Event: "rebirth_full_end",
			},
		})
	}()
}

// runPhoenixRebirthTimer 涅槃計時器：15 秒後若未全部擊破，結算消散
func (g *Game) runPhoenixRebirthTimer(session *phoenixRebirthSession) {
	time.Sleep(PhoenixRebirthDuration)

	mgr := g.LuckyPhoenixRebirth
	mgr.mu.Lock()
	if mgr.activeSession != session || session.settled {
		mgr.mu.Unlock()
		return
	}
	session.settled = true
	mgr.mu.Unlock()

	// 統計未被擊破的涅槃目標數
	remaining := 0
	for _, t := range session.targets {
		if !t.killed {
			remaining++
		}
	}

	log.Printf("[PhoenixRebirth] timeout player=%s remaining=%d", session.triggerPlayerID, remaining)

	// 廣播涅槃消散
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyPhoenixRebirth,
		Payload: ws.LuckyPhoenixRebirthPayload{
			Event:     "rebirth_fade",
			PlayerID:  session.triggerPlayerID,
			Remaining: remaining,
		},
	})

	if remaining > 0 {
		ann := g.Announce.Create(announce.EventLuckyPhoenixRebirth, session.triggerPlayerName, 0, map[string]string{
			"message": fmt.Sprintf("🦅 鳳凰涅槃消散，尚有 %d 條涅槃魚未被擊破...", remaining),
			"color":   "#888888",
		})
		g.broadcastAnnouncement(ann)
	}
}

// end of lucky_phoenix_rebirth_handler.go
