// lucky_element_fusion_handler.go — 幸運元素融合魚系統（DAY-263）
// 業界原創「元素碎片收集+融合爆發」機制
//
// 設計：擊破 T221 後，Server 隨機將「火/水/風」三種元素碎片各 1 個分配給場上 3 個目標：
//   - 玩家擊破帶有元素碎片的目標，收集對應元素（個人）
//   - 集齊 3 種不同元素 → 「元素融合爆發」：×6.0 倍率（個人最高）
//   - 只集齊 2 種 → 「部分融合」：×2.5 倍率（個人）
//   - 只集齊 1 種 → 「元素殘留」：×1.3 倍率（個人）
//   - 元素碎片存活 25 秒；個人冷卻 35 秒；全服冷卻 55 秒
//
// 設計差異：
//   - 與寶藏獵人（T218，30%機率發現碎片）不同，元素融合是「確定性收集」，
//     讓玩家有「我知道那條魚有元素，要趕快打」的策略感
//   - 「三種元素各有主題色」讓玩家一眼看出「哪條魚有什麼元素」
//   - 「×6.0 全融合」是目前個人類最高倍率，製造「哇，集齊了！」的爽感
//   - 「部分融合也有獎勵」確保即使沒集齊也有收益，降低挫敗感
//   - 「全服廣播元素分配位置」讓所有玩家都看到「哪條魚有元素」，製造「全服一起搶」的競爭感
//   - 業界依據：Royal Fishing 的 Element Combo 系統，2026 年最熱門「元素組合+融合爆發」機制
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
	LuckyElementFusionPersonalCD = 35 * time.Second // 個人冷卻
	LuckyElementFusionGlobalCD   = 55 * time.Second // 全服冷卻
	LuckyElementFusionDuration   = 25 * time.Second // 元素碎片存活時間
	LuckyElementFusionFullMult   = 6.0              // 集齊 3 種元素倍率
	LuckyElementFusionPartialMult = 2.5             // 集齊 2 種元素倍率
	LuckyElementFusionSingleMult  = 1.3             // 集齊 1 種元素倍率
)

// elementType 元素類型
type elementType string

const (
	ElementFire  elementType = "fire"  // 火元素（紅色）
	ElementWater elementType = "water" // 水元素（藍色）
	ElementWind  elementType = "wind"  // 風元素（綠色）
)

// elementFragmentEntry 元素碎片記錄
type elementFragmentEntry struct {
	instanceID string      // 目標實例 ID
	defID      string      // 目標定義 ID
	element    elementType // 元素類型
	expiresAt  time.Time   // 過期時間
}

// elementFusionSession 玩家元素收集 session
type elementFusionSession struct {
	playerID   string
	playerName string
	collected  map[elementType]bool // 已收集的元素
	expiresAt  time.Time
}

// luckyElementFusionManager 幸運元素融合魚管理器
type luckyElementFusionManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 場上帶有元素碎片的目標（instanceID → entry）
	activeFragments map[string]*elementFragmentEntry

	// 玩家收集 session（playerID → session）
	activeSessions map[string]*elementFusionSession
}

func newLuckyElementFusionManager() *luckyElementFusionManager {
	return &luckyElementFusionManager{
		personalCooldowns: make(map[string]time.Time),
		activeFragments:   make(map[string]*elementFragmentEntry),
		activeSessions:    make(map[string]*elementFusionSession),
	}
}

// isLuckyElementFusionFish 判斷是否為幸運元素融合魚
func isLuckyElementFusionFish(defID string) bool {
	return defID == "T221"
}

// isElementFragmentTarget 判斷目標是否帶有元素碎片（供 handleKill 使用）
func (m *luckyElementFusionManager) isElementFragmentTarget(instanceID string) (elementType, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if entry, ok := m.activeFragments[instanceID]; ok {
		if time.Now().Before(entry.expiresAt) {
			return entry.element, true
		}
		delete(m.activeFragments, instanceID)
	}
	return "", false
}

// removeFragment 移除元素碎片
func (m *luckyElementFusionManager) removeFragment(instanceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.activeFragments, instanceID)
}

// tryLuckyElementFusionFish 擊破 T221 後觸發元素融合系統
func (g *Game) tryLuckyElementFusionFish(p *player.Player) {
	m := g.LuckyElementFusion

	m.mu.Lock()
	now := time.Now()

	// 全服冷卻檢查
	if now.Before(m.globalCooldownUntil) {
		m.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if cd, ok := m.personalCooldowns[p.ID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyElementFusionPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyElementFusionGlobalCD)
	m.mu.Unlock()

	// 取得場上目標（排除 T221 本身）
	g.mu.RLock()
	candidates := make([]struct {
		instanceID string
		defID      string
	}, 0, 8)
	for id, t := range g.Targets {
		if t.IsAlive && t.DefID != "T221" {
			candidates = append(candidates, struct {
				instanceID string
				defID      string
			}{id, t.DefID})
		}
	}
	g.mu.RUnlock()

	if len(candidates) < 1 {
		log.Printf("[ElementFusion] player=%s 場上目標不足，無法分配元素碎片", p.ID)
		return
	}

	// 隨機打亂候選目標
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	// 分配三種元素（最多 3 個目標）
	elements := []elementType{ElementFire, ElementWater, ElementWind}
	count := len(candidates)
	if count > 3 {
		count = 3
	}

	expiresAt := now.Add(LuckyElementFusionDuration)
	assignedFragments := make([]struct {
		instanceID string
		defID      string
		element    elementType
	}, 0, count)

	m.mu.Lock()
	for i := 0; i < count; i++ {
		entry := &elementFragmentEntry{
			instanceID: candidates[i].instanceID,
			defID:      candidates[i].defID,
			element:    elements[i],
			expiresAt:  expiresAt,
		}
		m.activeFragments[candidates[i].instanceID] = entry
		assignedFragments = append(assignedFragments, struct {
			instanceID string
			defID      string
			element    elementType
		}{candidates[i].instanceID, candidates[i].defID, elements[i]})
	}
	m.mu.Unlock()

	log.Printf("[ElementFusion] player=%s 觸發元素融合！分配 %d 個元素碎片，存活 25 秒",
		p.ID, count)

	// 個人通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyElementFusion,
		Payload: ws.LuckyElementFusionPayload{
			Event:       "fusion_start",
			TriggerName: p.DisplayName,
			FragmentCount: count,
			Duration:    int(LuckyElementFusionDuration.Seconds()),
		},
	})

	// 全服廣播（含元素碎片位置）
	fragInfos := make([]ws.ElementFragmentInfo, 0, len(assignedFragments))
	for _, f := range assignedFragments {
		fragInfos = append(fragInfos, ws.ElementFragmentInfo{
			InstanceID: f.instanceID,
			DefID:      f.defID,
			Element:    string(f.element),
		})
	}
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyElementFusion,
		Payload: ws.LuckyElementFusionPayload{
			Event:       "fusion_broadcast",
			TriggerName: p.DisplayName,
			Fragments:   fragInfos,
			Duration:    int(LuckyElementFusionDuration.Seconds()),
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyElementFusion, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🔥💧🌪️ %s 觸發元素融合！火/水/風三種元素碎片已分配到場上目標！集齊三種可獲得 ×%.0f 大獎！",
			p.DisplayName, LuckyElementFusionFullMult),
		"color": "#FF6B35",
	})

	// 啟動超時清理
	go g.runElementFusionTimeout(expiresAt, assignedFragments)
}

// notifyElementFragmentKill 玩家擊破帶有元素碎片的目標（由 handleKill 呼叫）
func (g *Game) notifyElementFragmentKill(p *player.Player, instanceID string, elem elementType, baseMult float64) {
	m := g.LuckyElementFusion

	// 移除碎片
	m.removeFragment(instanceID)

	// 取得或建立玩家 session
	m.mu.Lock()
	session, ok := m.activeSessions[p.ID]
	if !ok {
		session = &elementFusionSession{
			playerID:   p.ID,
			playerName: p.DisplayName,
			collected:  make(map[elementType]bool),
			expiresAt:  time.Now().Add(LuckyElementFusionDuration),
		}
		m.activeSessions[p.ID] = session
	}
	session.collected[elem] = true
	collectedCount := len(session.collected)
	collectedList := make([]string, 0, collectedCount)
	for e := range session.collected {
		collectedList = append(collectedList, string(e))
	}
	m.mu.Unlock()

	// 計算碎片獎勵（基礎倍率 × 0.5，作為收集獎勵）
	fragmentReward := int(baseMult * 0.5 * float64(g.getAvgBetCost()))
	if fragmentReward < 1 {
		fragmentReward = 1
	}
	p.AddReward(fragmentReward)

	log.Printf("[ElementFusion] player=%s 收集元素 %s，已收集 %d/3，碎片獎勵 %d",
		p.ID, elem, collectedCount, fragmentReward)

	// 個人通知：收集碎片
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyElementFusion,
		Payload: ws.LuckyElementFusionPayload{
			Event:          "fragment_collect",
			PlayerID:       p.ID,
			PlayerName:     p.DisplayName,
			Element:        string(elem),
			CollectedCount: collectedCount,
			CollectedList:  collectedList,
			FragmentReward: fragmentReward,
		},
	})

	// 判斷是否觸發融合
	if collectedCount >= 3 {
		go g.doElementFusionBurst(p, 3)
	}
}

// doElementFusionBurst 觸發元素融合爆發
func (g *Game) doElementFusionBurst(p *player.Player, collectedCount int) {
	m := g.LuckyElementFusion

	// 移除玩家 session
	m.mu.Lock()
	delete(m.activeSessions, p.ID)
	m.mu.Unlock()

	// 計算倍率
	var mult float64
	var eventType string
	var eventLabel string
	switch collectedCount {
	case 3:
		mult = LuckyElementFusionFullMult
		eventType = "fusion_burst"
		eventLabel = "🔥💧🌪️ 元素全融合！"
	case 2:
		mult = LuckyElementFusionPartialMult
		eventType = "fusion_partial"
		eventLabel = "⚡ 部分融合！"
	default:
		mult = LuckyElementFusionSingleMult
		eventType = "fusion_single"
		eventLabel = "✨ 元素殘留！"
	}

	// 計算獎勵
	reward := int(mult * float64(g.getAvgBetCost()))
	if reward < 1 {
		reward = 1
	}
	p.AddReward(reward)

	log.Printf("[ElementFusion] player=%s 元素融合爆發！收集 %d/3，倍率 ×%.1f，獎勵 %d",
		p.ID, collectedCount, mult, reward)

	// 個人通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyElementFusion,
		Payload: ws.LuckyElementFusionPayload{
			Event:          eventType,
			PlayerID:       p.ID,
			PlayerName:     p.DisplayName,
			CollectedCount: collectedCount,
			Mult:           mult,
			Reward:         reward,
		},
	})

	// 全融合時全服廣播
	if collectedCount == 3 {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyElementFusion,
			Payload: ws.LuckyElementFusionPayload{
				Event:      "fusion_burst_broadcast",
				PlayerName: p.DisplayName,
				Mult:       mult,
				Reward:     reward,
			},
		})
		// 全服公告
		g.Announce.Create(announce.EventLuckyElementFusion, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🔥💧🌪️ %s 集齊三種元素！元素全融合爆發！×%.0f 大獎 +%d！",
				p.DisplayName, mult, reward),
			"color": "#FFD700",
		})
	}

	_ = eventLabel // suppress unused warning
}

// runElementFusionTimeout 元素碎片超時清理
func (g *Game) runElementFusionTimeout(expiresAt time.Time, fragments []struct {
	instanceID string
	defID      string
	element    elementType
}) {
	remaining := time.Until(expiresAt)
	if remaining > 0 {
		time.Sleep(remaining)
	}

	m := g.LuckyElementFusion

	// 清理剩餘碎片
	m.mu.Lock()
	for _, f := range fragments {
		delete(m.activeFragments, f.instanceID)
	}

	// 對仍有 session 的玩家結算（部分融合或殘留）
	expiredSessions := make([]*elementFusionSession, 0)
	for _, s := range m.activeSessions {
		if time.Now().After(s.expiresAt) {
			expiredSessions = append(expiredSessions, s)
			delete(m.activeSessions, s.playerID)
		}
	}
	m.mu.Unlock()

	// 對超時 session 結算
	for _, s := range expiredSessions {
		count := len(s.collected)
		if count == 0 {
			continue
		}

		// 找到玩家
		g.mu.RLock()
		pl, ok := g.Players[s.playerID]
		g.mu.RUnlock()
		if !ok {
			continue
		}

		go g.doElementFusionBurst(pl, count)
	}

	// 廣播元素碎片消失
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyElementFusion,
		Payload: ws.LuckyElementFusionPayload{
			Event: "fusion_expire",
		},
	})
}
