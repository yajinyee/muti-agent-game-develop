// lucky_rainbow_bridge_handler.go — 幸運彩虹橋魚系統（DAY-279）
// 業界原創「彩虹橋連接+跨目標連鎖傷害+彩虹爆發」機制
//
// 設計：擊破 T237 後，Server 在場上隨機選 3 個目標，用「彩虹橋」連接：
//   - 彩虹橋期間（12 秒），任何玩家擊破這 3 個目標中的任何一個
//     → 其他 2 個目標也獲得 HP -40%（連鎖傷害）
//   - 若 3 個目標都在 12 秒內被擊破 → 觸發「彩虹爆發」：全服 ×2.0 加成 6 秒
//   - 若 12 秒後未全部擊破 → 「彩虹消散」：剩餘目標 HP -60%（安慰獎）
//   - 全服廣播彩虹橋連接的目標和爆發結果
//   - 個人冷卻 25 秒；全服冷卻 40 秒
//
// 設計差異：
//   - 與連鎖爆炸（T224，從一點向外擴散）不同，彩虹橋是「連接特定目標」
//     讓玩家有「打這條魚，另外兩條也會受傷」的策略感
//   - 「3 個目標全部擊破觸發彩虹爆發」讓玩家有「要趁 12 秒內打完 3 條」的緊迫感
//   - 「全服 ×2.0 加成 6 秒」讓所有玩家都受益，製造「全服一起爽」的社交感
//   - 「彩虹消散 HP -60%」確保即使沒打完也有收益，降低挫敗感
//   - 「全服廣播彩虹橋目標」讓所有玩家看到「哪 3 條魚被連接了」，製造策略感
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
	LuckyRainbowBridgePersonalCD    = 25 * time.Second // 個人冷卻
	LuckyRainbowBridgeGlobalCD      = 40 * time.Second // 全服冷卻
	LuckyRainbowBridgeDuration      = 12 * time.Second // 彩虹橋持續時間
	LuckyRainbowBridgeTargetCount   = 3                // 連接目標數
	LuckyRainbowBridgeChainHPDamage = 0.40             // 連鎖傷害比例（-40%）
	LuckyRainbowBridgeBurstMult     = 2.0              // 彩虹爆發全服倍率
	LuckyRainbowBridgeBurstDuration = 6 * time.Second  // 彩虹爆發持續時間
	LuckyRainbowBridgeFadeHPDamage  = 0.60             // 彩虹消散傷害比例（-60%）
)

// rainbowBridgeSession 彩虹橋會話
type rainbowBridgeSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	targetIIDs        []string // 3 個連接目標的 instanceID
	targetNames       []string // 3 個連接目標的名稱
	killedIIDs        map[string]bool
	expiresAt         time.Time
	settled           bool
}

// rainbowBurstBoost 彩虹爆發加成
type rainbowBurstBoost struct {
	mult      float64
	expiresAt time.Time
}

// luckyRainbowBridgeManager 幸運彩虹橋魚管理器
type luckyRainbowBridgeManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前彩虹橋會話（同時只有一個）
	activeSession *rainbowBridgeSession

	// 彩虹爆發加成（全服）
	burstBoost *rainbowBurstBoost
}

func newLuckyRainbowBridgeManager() *luckyRainbowBridgeManager {
	return &luckyRainbowBridgeManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyRainbowBridgeFish 判斷是否為幸運彩虹橋魚
func isLuckyRainbowBridgeFish(defID string) bool {
	return defID == "T237"
}

// getRainbowBridgeBurstMult 取得彩虹爆發倍率（供 handleKill 使用）
func (m *luckyRainbowBridgeManager) getRainbowBridgeBurstMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.burstBoost != nil && time.Now().Before(m.burstBoost.expiresAt) {
		return m.burstBoost.mult
	}
	return 1.0
}

// isRainbowBridgeTarget 判斷是否為彩虹橋連接目標
func (m *luckyRainbowBridgeManager) isRainbowBridgeTarget(instanceID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil || m.activeSession.settled {
		return false
	}
	for _, iid := range m.activeSession.targetIIDs {
		if iid == instanceID {
			return true
		}
	}
	return false
}

// tryLuckyRainbowBridgeFish 擊破 T237 後觸發彩虹橋（供 handleKill 使用）
func (g *Game) tryLuckyRainbowBridgeFish(p *player.Player) {
	mgr := g.LuckyRainbowBridge
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
	// 已有活躍彩虹橋
	if mgr.activeSession != nil && !mgr.activeSession.settled {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = now.Add(LuckyRainbowBridgePersonalCD)
	mgr.globalCooldownUntil = now.Add(LuckyRainbowBridgeGlobalCD)
	mgr.mu.Unlock()

	// 選取 3 個目標
	g.mu.RLock()
	var candidates []string
	var candidateNames []string
	for iid, t := range g.Targets {
		if t.HP > 0 && !isLuckyRainbowBridgeFish(t.DefID) {
			candidates = append(candidates, iid)
			candidateNames = append(candidateNames, t.Def.Name)
		}
	}
	g.mu.RUnlock()

	if len(candidates) < 2 {
		log.Printf("[RainbowBridge] not enough targets (%d)", len(candidates))
		return
	}

	// 隨機選取（最多 3 個）
	count := LuckyRainbowBridgeTargetCount
	if count > len(candidates) {
		count = len(candidates)
	}
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
		candidateNames[i], candidateNames[j] = candidateNames[j], candidateNames[i]
	})
	selectedIIDs := candidates[:count]
	selectedNames := candidateNames[:count]

	session := &rainbowBridgeSession{
		triggerPlayerID:   p.ID,
		triggerPlayerName: p.DisplayName,
		targetIIDs:        selectedIIDs,
		targetNames:       selectedNames,
		killedIIDs:        make(map[string]bool),
		expiresAt:         time.Now().Add(LuckyRainbowBridgeDuration),
	}

	mgr.mu.Lock()
	mgr.activeSession = session
	mgr.mu.Unlock()

	log.Printf("[RainbowBridge] player=%s triggered rainbow bridge, targets=%v", p.ID, selectedIIDs)

	// 全服廣播：彩虹橋觸發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyRainbowBridge,
		Payload: ws.LuckyRainbowBridgePayload{
			Event:       "bridge_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			TargetIIDs:  selectedIIDs,
			TargetNames: selectedNames,
			Duration:    int(LuckyRainbowBridgeDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyRainbowBridge, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌈 %s 觸發彩虹橋！連接 %d 個目標！打一個，其他兩個 HP -40%%！12 秒內全打完觸發彩虹爆發！",
			p.DisplayName, count),
		"color": "#FF69B4",
	})
	g.broadcastAnnouncement(ann)

	// 啟動彩虹橋計時器
	go g.runRainbowBridgeTimer(p, session)
}

// notifyRainbowBridgeKill 彩虹橋目標被擊破時的連鎖處理（由 handleKill 呼叫）
func (g *Game) notifyRainbowBridgeKill(p *player.Player, killedIID string) {
	mgr := g.LuckyRainbowBridge
	mgr.mu.Lock()

	session := mgr.activeSession
	if session == nil || session.settled {
		mgr.mu.Unlock()
		return
	}

	// 記錄擊破
	session.killedIIDs[killedIID] = true
	killedCount := len(session.killedIIDs)
	totalCount := len(session.targetIIDs)

	// 找出其他未被擊破的目標
	var otherIIDs []string
	for _, iid := range session.targetIIDs {
		if !session.killedIIDs[iid] {
			otherIIDs = append(otherIIDs, iid)
		}
	}
	mgr.mu.Unlock()

	log.Printf("[RainbowBridge] target %s killed! killed=%d/%d others=%v",
		killedIID, killedCount, totalCount, otherIIDs)

	// 對其他目標造成連鎖傷害（HP -40%）
	for _, otherIID := range otherIIDs {
		g.mu.Lock()
		t, ok := g.Targets[otherIID]
		if ok && t.HP > 0 {
			damage := int(float64(t.HP) * LuckyRainbowBridgeChainHPDamage)
			if damage < 1 {
				damage = 1
			}
			t.HP -= damage
			if t.HP < 1 {
				t.HP = 1
			}
		}
		g.mu.Unlock()
	}

	// 廣播連鎖傷害
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyRainbowBridge,
		Payload: ws.LuckyRainbowBridgePayload{
			Event:       "bridge_chain",
			PlayerName:  p.DisplayName,
			KilledIID:   killedIID,
			OtherIIDs:   otherIIDs,
			KilledCount: killedCount,
			TotalCount:  totalCount,
		},
	})

	// 檢查是否全部擊破
	if killedCount >= totalCount {
		go g.doRainbowBridgeBurst(p, session)
	}
}

// runRainbowBridgeTimer 彩虹橋計時器（goroutine）
func (g *Game) runRainbowBridgeTimer(p *player.Player, session *rainbowBridgeSession) {
	select {
	case <-time.After(LuckyRainbowBridgeDuration):
	case <-g.stopCh:
		return
	}

	g.LuckyRainbowBridge.mu.Lock()
	if session.settled {
		g.LuckyRainbowBridge.mu.Unlock()
		return
	}
	session.settled = true
	killedCount := len(session.killedIIDs)
	totalCount := len(session.targetIIDs)
	g.LuckyRainbowBridge.mu.Unlock()

	if killedCount >= totalCount {
		// 已全部擊破（由 notifyRainbowBridgeKill 處理）
		return
	}

	// 彩虹消散：剩餘目標 HP -60%
	fadedCount := 0
	for _, iid := range session.targetIIDs {
		if session.killedIIDs[iid] {
			continue
		}
		g.mu.Lock()
		t, ok := g.Targets[iid]
		if ok && t.HP > 0 {
			damage := int(float64(t.HP) * LuckyRainbowBridgeFadeHPDamage)
			if damage < 1 {
				damage = 1
			}
			t.HP -= damage
			if t.HP < 1 {
				t.HP = 1
			}
			fadedCount++
		}
		g.mu.Unlock()
	}

	log.Printf("[RainbowBridge] FADE: killed=%d/%d faded=%d", killedCount, totalCount, fadedCount)

	// 廣播彩虹消散
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyRainbowBridge,
		Payload: ws.LuckyRainbowBridgePayload{
			Event:       "bridge_fade",
			PlayerName:  p.DisplayName,
			KilledCount: killedCount,
			TotalCount:  totalCount,
			FadedCount:  fadedCount,
		},
	})

	ann := g.Announce.Create(announce.EventLuckyRainbowBridge, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌈 彩虹橋消散！%s 擊破 %d/%d 個目標，剩餘目標 HP -60%%",
			p.DisplayName, killedCount, totalCount),
		"color": "#95A5A6",
	})
	g.broadcastAnnouncement(ann)
}

// doRainbowBridgeBurst 彩虹爆發（全部擊破時觸發）
func (g *Game) doRainbowBridgeBurst(p *player.Player, session *rainbowBridgeSession) {
	g.LuckyRainbowBridge.mu.Lock()
	if session.settled {
		g.LuckyRainbowBridge.mu.Unlock()
		return
	}
	session.settled = true

	// 設定全服爆發加成
	g.LuckyRainbowBridge.burstBoost = &rainbowBurstBoost{
		mult:      LuckyRainbowBridgeBurstMult,
		expiresAt: time.Now().Add(LuckyRainbowBridgeBurstDuration),
	}
	g.LuckyRainbowBridge.mu.Unlock()

	log.Printf("[RainbowBridge] BURST! player=%s all %d targets killed! global x%.1f for %ds",
		p.ID, len(session.targetIIDs), LuckyRainbowBridgeBurstMult, int(LuckyRainbowBridgeBurstDuration.Seconds()))

	// 全服廣播：彩虹爆發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyRainbowBridge,
		Payload: ws.LuckyRainbowBridgePayload{
			Event:        "bridge_burst",
			PlayerName:   p.DisplayName,
			BurstMult:    LuckyRainbowBridgeBurstMult,
			BurstSeconds: int(LuckyRainbowBridgeBurstDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyRainbowBridge, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌈 彩虹爆發！%s 擊破全部 %d 個目標！全服 ×%.1f 加成 %d 秒！",
			p.DisplayName, len(session.targetIIDs), LuckyRainbowBridgeBurstMult, int(LuckyRainbowBridgeBurstDuration.Seconds())),
		"color": "#FFD700",
	})
	g.broadcastAnnouncement(ann)

	// 爆發結束後清除加成
	go func() {
		select {
		case <-time.After(LuckyRainbowBridgeBurstDuration):
		case <-g.stopCh:
			return
		}
		g.LuckyRainbowBridge.mu.Lock()
		g.LuckyRainbowBridge.burstBoost = nil
		g.LuckyRainbowBridge.mu.Unlock()

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyRainbowBridge,
			Payload: ws.LuckyRainbowBridgePayload{
				Event: "bridge_burst_end",
			},
		})
	}()
}
