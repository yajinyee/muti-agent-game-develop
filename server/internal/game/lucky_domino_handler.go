// lucky_domino_handler.go — 幸運多米諾魚系統（DAY-288）
// 業界依據：Fishing Fortune 2026「multiplier cascade system — consecutive rare catches within 90s」
//          Royal Fishing Jili「chain reactions that jump between nearby fish」
//          業界原創「多米諾骨牌連鎖+場上目標依序倒下+最終爆發」機制
//
// 設計：
//   - 擊破 T246 後，場上隨機選 5 個目標設為「多米諾骨牌」
//   - 第 1 個骨牌立即 HP -80%（幾乎必死）；若被擊破 → 連鎖推倒第 2 個（HP -80%）→ 依序推倒到第 5 個
//   - 每推倒一個骨牌，觸發玩家獲得 ×1.5 累積倍率（最高 ×7.5）
//   - 若 5 個骨牌全部在 20 秒內被推倒 → 「多米諾完美」：全服 ×2.5 加成 7 秒
//   - 全服廣播骨牌位置和連鎖結果
//   - 個人冷卻 26 秒；全服冷卻 42 秒
//
// 設計差異：
//   - 與宇宙脈衝（T245，同心圓擴散波）不同，多米諾是「場上特定目標依序連鎖倒下」
//   - 「HP -80% 弱化」讓骨牌目標幾乎必死，玩家有「快去打第一個骨牌，後面會連鎖」的緊迫感
//   - 「每推倒一個 ×1.5 累積」讓玩家有「要把 5 個骨牌全部推倒才能拿到 ×7.5」的動力
//   - 「多米諾完美（全部推倒）→ 全服 ×2.5」製造「全服一起爽」的社交感
//   - 「全服廣播骨牌位置」讓所有玩家看到「哪 5 條魚是骨牌」，製造全服搶打感
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
	LuckyDominoPersonalCD = 26 * time.Second // 個人冷卻
	LuckyDominoGlobalCD   = 42 * time.Second // 全服冷卻

	// 多米諾骨牌設計
	DominoCount       = 5                      // 骨牌數量
	DominoHPDmg       = 0.80                   // 每個骨牌 HP -80%
	DominoAccumMult   = 1.5                    // 每推倒一個累積倍率
	DominoMaxMult     = 7.5                    // 最高累積倍率（5 個 × 1.5 = 7.5）
	DominoTimeout     = 20 * time.Second       // 骨牌有效時間

	// 多米諾完美：全服加成
	DominoPerfectMult     = 2.5                   // 全服 ×2.5
	DominoPerfectDuration = 7 * time.Second       // 持續 7 秒
)

// dominoPerfectBoost 多米諾完美全服加成
type dominoPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

// dominoTarget 多米諾骨牌目標
type dominoTarget struct {
	instanceID string
	defID      string
	name       string
	x, y       float64
	knocked    bool // 是否已被推倒
}

// dominoSession 多米諾骨牌會話
type dominoSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	targets           []*dominoTarget
	currentIdx        int       // 當前等待被推倒的骨牌序號（0-based）
	accumMult         float64
	expiresAt         time.Time
	settled           bool
}

// luckyDominoManager 幸運多米諾魚管理器
type luckyDominoManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前活躍的多米諾會話
	activeSession *dominoSession

	// 多米諾完美全服加成
	perfectBoost *dominoPerfectBoost
}

func newLuckyDominoManager() *luckyDominoManager {
	return &luckyDominoManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyDominoFish 判斷是否為幸運多米諾魚
func isLuckyDominoFish(defID string) bool {
	return defID == "T246"
}

// getDominoPerfectMult 取得多米諾完美全服加成倍率（供 handleKill 使用）
func (m *luckyDominoManager) getDominoPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

// isDominoTarget 判斷是否為多米諾骨牌目標，回傳（是否為骨牌, 骨牌序號, 累積倍率）
func (m *luckyDominoManager) isDominoTarget(instanceID string) (bool, int, float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil || m.activeSession.settled {
		return false, 0, 1.0
	}
	sess := m.activeSession
	if time.Now().After(sess.expiresAt) {
		return false, 0, 1.0
	}
	for i, t := range sess.targets {
		if t.instanceID == instanceID && !t.knocked && i == sess.currentIdx {
			return true, i, sess.accumMult
		}
	}
	return false, 0, 1.0
}

// tryLuckyDominoFish 擊破 T246 後觸發多米諾（供 handleKill 使用）
func (g *Game) tryLuckyDominoFish(p *player.Player) {
	mgr := g.LuckyDomino
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
	mgr.personalCooldowns[p.ID] = now.Add(LuckyDominoPersonalCD)
	mgr.globalCooldownUntil = now.Add(LuckyDominoGlobalCD)
	mgr.mu.Unlock()

	// 選取 5 個骨牌目標
	targets := g.selectDominoTargets(DominoCount)
	if len(targets) == 0 {
		log.Printf("[Domino] no targets available, skip")
		return
	}

	// 建立會話
	mgr.mu.Lock()
	mgr.activeSession = &dominoSession{
		triggerPlayerID:   p.ID,
		triggerPlayerName: p.DisplayName,
		targets:           targets,
		currentIdx:        0,
		accumMult:         1.0,
		expiresAt:         now.Add(DominoTimeout),
	}
	mgr.mu.Unlock()

	log.Printf("[Domino] player=%s targets=%d", p.ID, len(targets))

	// 對第 1 個骨牌施加 HP -80%
	g.applyDominoKnock(targets[0])

	// 廣播骨牌位置
	targetInfos := make([]ws.DominoTargetInfo, len(targets))
	for i, t := range targets {
		targetInfos[i] = ws.DominoTargetInfo{
			InstanceID: t.instanceID,
			Name:       t.name,
			X:          t.x,
			Y:          t.y,
			Idx:        i + 1,
		}
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyDomino,
		Payload: ws.LuckyDominoPayload{
			Event:      "domino_start",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Targets:    targetInfos,
			TotalCount: len(targets),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyDomino, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🀱🎯 %s 觸發多米諾！%d 個骨牌等待連鎖推倒！",
			p.DisplayName, len(targets)),
		"color": "#8B4513",
	})
	g.broadcastAnnouncement(ann)

	// 啟動超時 goroutine
	go g.runDominoTimeout(p)
}

// selectDominoTargets 選取場上 n 個存活目標作為骨牌
func (g *Game) selectDominoTargets(n int) []*dominoTarget {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var candidates []*dominoTarget
	for _, t := range g.Targets {
		if t == nil || t.HP <= 0 {
			continue
		}
		candidates = append(candidates, &dominoTarget{
			instanceID: t.InstanceID,
			defID:      t.DefID,
			name:       t.Name,
			x:          t.X,
			y:          t.Y,
		})
	}

	// 隨機選取 n 個（Fisher-Yates shuffle 前 n 個）
	for i := len(candidates) - 1; i > 0; i-- {
		j := int(time.Now().UnixNano()) % (i + 1)
		if j < 0 {
			j = -j
		}
		candidates[i], candidates[j] = candidates[j], candidates[i]
	}

	if len(candidates) > n {
		candidates = candidates[:n]
	}
	return candidates
}

// applyDominoKnock 對骨牌目標施加 HP -80%
func (g *Game) applyDominoKnock(dt *dominoTarget) {
	g.mu.Lock()
	defer g.mu.Unlock()

	t, ok := g.Targets[dt.instanceID]
	if !ok || t == nil || t.HP <= 0 {
		return
	}

	dmg := int(float64(t.HP) * DominoHPDmg)
	if dmg < 1 {
		dmg = 1
	}
	t.HP -= dmg
	if t.HP < 0 {
		t.HP = 0
	}

	// 廣播 HP 更新
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTargetUpdate,
		Payload: ws.TargetUpdatePayload{
			InstanceID: t.InstanceID,
			HP:         t.HP,
			MaxHP:      t.MaxHP,
		},
	})
}

// notifyDominoKill 骨牌目標被擊破時呼叫（供 handleKill 使用）
func (g *Game) notifyDominoKill(p *player.Player, instanceID string) float64 {
	mgr := g.LuckyDomino
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if mgr.activeSession == nil || mgr.activeSession.settled {
		return 1.0
	}
	sess := mgr.activeSession
	if time.Now().After(sess.expiresAt) {
		return 1.0
	}

	// 確認是當前等待的骨牌
	if sess.currentIdx >= len(sess.targets) {
		return 1.0
	}
	current := sess.targets[sess.currentIdx]
	if current.instanceID != instanceID || current.knocked {
		return 1.0
	}

	// 標記推倒
	current.knocked = true
	newMult := sess.accumMult * DominoAccumMult
	if newMult > DominoMaxMult {
		newMult = DominoMaxMult
	}
	sess.accumMult = newMult
	knockedIdx := sess.currentIdx
	sess.currentIdx++

	log.Printf("[Domino] knocked %d/%d instanceID=%s accumMult=%.2f",
		knockedIdx+1, len(sess.targets), instanceID, newMult)

	// 廣播骨牌推倒
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyDomino,
		Payload: ws.LuckyDominoPayload{
			Event:      "domino_knock",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			KnockedIdx: knockedIdx + 1,
			TotalCount: len(sess.targets),
			AccumMult:  newMult,
		},
	})

	// 若還有下一個骨牌，施加 HP -80%
	if sess.currentIdx < len(sess.targets) {
		next := sess.targets[sess.currentIdx]
		go g.applyDominoKnock(next)

		// 廣播下一個骨牌提示
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyDomino,
			Payload: ws.LuckyDominoPayload{
				Event:      "domino_next",
				NextIdx:    sess.currentIdx + 1,
				TotalCount: len(sess.targets),
				NextTarget: ws.DominoTargetInfo{
					InstanceID: next.instanceID,
					Name:       next.name,
					X:          next.x,
					Y:          next.y,
					Idx:        sess.currentIdx + 1,
				},
			},
		})
	} else {
		// 全部推倒！多米諾完美
		sess.settled = true
		go g.doDominoPerfect(p, sess.accumMult)
	}

	return newMult
}

// doDominoPerfect 多米諾完美：全服 ×2.5 加成 7 秒
func (g *Game) doDominoPerfect(p *player.Player, accumMult float64) {
	mgr := g.LuckyDomino
	mgr.mu.Lock()
	mgr.perfectBoost = &dominoPerfectBoost{
		mult:      DominoPerfectMult,
		expiresAt: time.Now().Add(DominoPerfectDuration),
	}
	mgr.mu.Unlock()

	log.Printf("[Domino] PERFECT! player=%s accumMult=%.2f global x%.1f for %v",
		p.ID, accumMult, DominoPerfectMult, DominoPerfectDuration)

	// 全服廣播多米諾完美
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyDomino,
		Payload: ws.LuckyDominoPayload{
			Event:       "domino_perfect",
			PlayerName:  p.DisplayName,
			AccumMult:   accumMult,
			PerfectMult: DominoPerfectMult,
			Duration:    int(DominoPerfectDuration.Seconds()),
		},
	})

	// 全服最高優先公告
	ann := g.Announce.Create(announce.EventLuckyDomino, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🀱🎯🀱 %s 多米諾完美！累積 ×%.1f！全服 ×%.1f 加成 %d 秒！",
			p.DisplayName, accumMult, DominoPerfectMult, int(DominoPerfectDuration.Seconds())),
		"color": "#4A0000",
	})
	g.broadcastAnnouncement(ann)

	// 7 秒後廣播完美結束
	go func() {
		time.Sleep(DominoPerfectDuration)
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyDomino,
			Payload: ws.LuckyDominoPayload{
				Event: "domino_perfect_end",
			},
		})
	}()
}

// runDominoTimeout 多米諾超時處理（20 秒後若未全部推倒則結算）
func (g *Game) runDominoTimeout(p *player.Player) {
	time.Sleep(DominoTimeout)

	mgr := g.LuckyDomino
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if mgr.activeSession == nil || mgr.activeSession.settled {
		return
	}
	sess := mgr.activeSession
	sess.settled = true

	knockedCount := 0
	for _, t := range sess.targets {
		if t.knocked {
			knockedCount++
		}
	}

	log.Printf("[Domino] timeout player=%s knocked=%d/%d accumMult=%.2f",
		p.ID, knockedCount, len(sess.targets), sess.accumMult)

	// 廣播超時結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyDomino,
		Payload: ws.LuckyDominoPayload{
			Event:       "domino_end",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			KnockedIdx:  knockedCount,
			TotalCount:  len(sess.targets),
			AccumMult:   sess.accumMult,
			IsPerfect:   false,
		},
	})

	// 有推倒至少 2 個才公告
	if knockedCount >= 2 {
		ann := g.Announce.Create(announce.EventLuckyDomino, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🀱 %s 多米諾結算！推倒 %d/%d 個骨牌，累積 ×%.1f！",
				p.DisplayName, knockedCount, len(sess.targets), sess.accumMult),
			"color": "#8B4513",
		})
		g.broadcastAnnouncement(ann)
	}
}

// end of lucky_domino_handler.go
