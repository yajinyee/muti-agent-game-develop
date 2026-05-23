// lucky_multiplier_stack_handler.go — 幸運倍率疊加魚系統（DAY-267）
// 業界依據：Fishing Fortune Multiplier Cascade 機制（2026 年最熱門）
//
// 設計：擊破 T225 後，觸發「倍率疊加模式」（持續 25 秒）：
//   - 玩家每次擊破任何目標，疊加倍率 +0.3x（從 1.0x 開始，最高 10.0x）
//   - 每次擊破都用「當前疊加倍率」計算獎勵（個人）
//   - 達到 10.0x 時觸發「倍率爆發」：最後一次擊破獲得 ×20.0 大獎（個人）
//   - 25 秒後未達到 10.0x → 「倍率結算」：用最終疊加倍率計算最後一次擊破獎勵
//   - 個人冷卻 32 秒；全服冷卻 50 秒
//
// 設計差異：
//   - 與連鎖爆炸（T224，空間擴散）不同，倍率疊加是「時間累積」，讓玩家有「越打越高，要趁 25 秒內打滿 10.0x」的緊迫感
//   - 「每次擊破 +0.3x」讓玩家有「每一槍都在累積倍率」的動力
//   - 「達到 10.0x 觸發 ×20.0 爆發」讓玩家有「要趁疊加期間打滿 30 個目標」的策略感
//   - 「倍率計數器即時顯示」讓玩家看到「現在疊加到幾倍了」，製造「快滿了！」的期待感
//   - 「全服廣播倍率爆發」讓所有玩家看到「有人疊加到 10.0x 爆發了」，製造羨慕感
//   - 業界依據：Fishing Fortune 的 Multiplier Cascade（2026 年最熱門），讓玩家有「每次擊破稀有魚都在疊加倍率，越打越高」的爽感
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
	LuckyMultiplierStackPersonalCD  = 32 * time.Second // 個人冷卻
	LuckyMultiplierStackGlobalCD    = 50 * time.Second // 全服冷卻
	LuckyMultiplierStackDuration    = 25 * time.Second // 疊加模式持續時間
	LuckyMultiplierStackStep        = 0.3              // 每次擊破疊加量
	LuckyMultiplierStackInitial     = 1.0              // 初始疊加倍率
	LuckyMultiplierStackMax         = 10.0             // 最大疊加倍率
	LuckyMultiplierStackBurstMult   = 20.0             // 爆發倍率（達到最大時）
)

// multiplierStackSession 倍率疊加 session
type multiplierStackSession struct {
	playerID   string
	playerName string
	expiresAt  time.Time
	stack      float64 // 當前疊加倍率
	killCount  int     // 本次疊加期間擊破數
	totalReward int    // 本次疊加期間總獎勵
	burst      bool    // 是否已爆發
}

// luckyMultiplierStackManager 幸運倍率疊加魚管理器
type luckyMultiplierStackManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 活躍 session（playerID → session）
	activeSessions map[string]*multiplierStackSession
}

func newLuckyMultiplierStackManager() *luckyMultiplierStackManager {
	return &luckyMultiplierStackManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*multiplierStackSession),
	}
}

// isLuckyMultiplierStackFish 判斷是否為幸運倍率疊加魚
func isLuckyMultiplierStackFish(defID string) bool {
	return defID == "T225"
}

// isMultiplierStackActive 判斷玩家是否在倍率疊加模式中（供 handleKill 使用）
func (m *luckyMultiplierStackManager) isMultiplierStackActive(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	if !ok {
		return false
	}
	if time.Now().After(sess.expiresAt) {
		delete(m.activeSessions, playerID)
		return false
	}
	return !sess.burst
}

// getMultiplierStackBonus 取得當前疊加倍率（供 handleKill 使用）
// 同時更新疊加值，回傳（倍率加成, 是否爆發）
func (m *luckyMultiplierStackManager) getMultiplierStackBonus(playerID string) (float64, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	if !ok {
		return 1.0, false
	}
	if time.Now().After(sess.expiresAt) {
		delete(m.activeSessions, playerID)
		return 1.0, false
	}
	if sess.burst {
		return 1.0, false
	}
	return sess.stack, false
}

// tryLuckyMultiplierStackFish 擊破 T225 後觸發倍率疊加
func (g *Game) tryLuckyMultiplierStackFish(p *player.Player) {
	m := g.LuckyMultiplierStack

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
	// 已有活躍 session 時不重複觸發
	if _, ok := m.activeSessions[p.ID]; ok {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyMultiplierStackPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyMultiplierStackGlobalCD)

	// 建立 session
	sess := &multiplierStackSession{
		playerID:   p.ID,
		playerName: p.DisplayName,
		expiresAt:  now.Add(LuckyMultiplierStackDuration),
		stack:      LuckyMultiplierStackInitial,
		killCount:  0,
		totalReward: 0,
		burst:      false,
	}
	m.activeSessions[p.ID] = sess
	m.mu.Unlock()

	log.Printf("[MultiplierStack] player=%s 觸發倍率疊加！持續 %.0f 秒",
		p.ID, LuckyMultiplierStackDuration.Seconds())

	// 個人訊息：觸發者
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyMultiplierStack,
		Payload: ws.LuckyMultiplierStackPayload{
			Event:        "stack_start",
			PlayerID:     p.ID,
			PlayerName:   p.DisplayName,
			CurrentStack: LuckyMultiplierStackInitial,
			MaxStack:     LuckyMultiplierStackMax,
			BurstMult:    LuckyMultiplierStackBurstMult,
			Duration:     LuckyMultiplierStackDuration.Seconds(),
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMultiplierStack,
		Payload: ws.LuckyMultiplierStackPayload{
			Event:      "stack_broadcast",
			PlayerName: p.DisplayName,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyMultiplierStack, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("📈 %s 觸發倍率疊加！每次擊破 +%.1fx，最高 ×%.0f！",
			p.DisplayName, LuckyMultiplierStackStep, LuckyMultiplierStackMax),
		"color": "#00FF88",
	})

	// 啟動超時 goroutine
	go g.runMultiplierStackTimeout(p, sess)
}

// notifyMultiplierStackKill 玩家在疊加模式中擊破目標時呼叫
// 回傳疊加倍率加成（用於 handleKill 的獎勵計算）
func (g *Game) notifyMultiplierStackKill(p *player.Player, targetName string, baseReward int) float64 {
	m := g.LuckyMultiplierStack

	m.mu.Lock()
	sess, ok := m.activeSessions[p.ID]
	if !ok || sess.burst || time.Now().After(sess.expiresAt) {
		m.mu.Unlock()
		return 1.0
	}

	// 取得當前疊加倍率（用於本次獎勵）
	currentStack := sess.stack

	// 疊加 +0.3x
	sess.stack += LuckyMultiplierStackStep
	if sess.stack > LuckyMultiplierStackMax {
		sess.stack = LuckyMultiplierStackMax
	}
	newStack := sess.stack
	sess.killCount++

	// 計算本次獎勵（基礎獎勵 × 疊加倍率）
	stackReward := int(float64(baseReward) * currentStack)
	sess.totalReward += stackReward

	isBurst := newStack >= LuckyMultiplierStackMax
	if isBurst {
		sess.burst = true
	}
	m.mu.Unlock()

	// 發送疊加更新給玩家
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyMultiplierStack,
		Payload: ws.LuckyMultiplierStackPayload{
			Event:        "stack_update",
			PlayerID:     p.ID,
			CurrentStack: newStack,
			KillCount:    sess.killCount,
			TargetName:   targetName,
			Reward:       stackReward,
		},
	})

	// 達到最大疊加 → 觸發爆發
	if isBurst {
		go g.doMultiplierStackBurst(p, sess)
	}

	return currentStack
}

// doMultiplierStackBurst 倍率爆發（達到 10.0x 時觸發）
func (g *Game) doMultiplierStackBurst(p *player.Player, sess *multiplierStackSession) {
	// 計算爆發獎勵（最後一次擊破的 ×20.0 大獎）
	betDef := p.GetBetDef()
	betCost := 1
	if betDef != nil {
		betCost = betDef.BetCost
	}
	burstReward := int(float64(betCost) * LuckyMultiplierStackBurstMult * LuckyMultiplierStackMax)
	if burstReward < 1 {
		burstReward = 1
	}

	// 給予玩家爆發獎勵
	p.AddCoins(burstReward)

	totalReward := sess.totalReward + burstReward

	log.Printf("[MultiplierStack] player=%s 倍率爆發！疊加 %.1fx，爆發獎勵 %d，總獎勵 %d",
		p.ID, LuckyMultiplierStackMax, burstReward, totalReward)

	// 個人爆發訊息
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyMultiplierStack,
		Payload: ws.LuckyMultiplierStackPayload{
			Event:       "stack_burst",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			TotalStack:  LuckyMultiplierStackMax,
			BurstReward: burstReward,
			TotalReward: totalReward,
		},
	})

	// 全服廣播爆發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMultiplierStack,
		Payload: ws.LuckyMultiplierStackPayload{
			Event:       "stack_burst_broadcast",
			PlayerName:  p.DisplayName,
			TotalStack:  LuckyMultiplierStackMax,
			BurstReward: burstReward,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyMultiplierStack, p.DisplayName, burstReward, map[string]string{
		"message": fmt.Sprintf("📈 %s 倍率疊加達到 ×%.0f！爆發獎勵 +%d！",
			p.DisplayName, LuckyMultiplierStackMax, burstReward),
		"color": "#FFD700",
	})

	// 清除 session
	m := g.LuckyMultiplierStack
	m.mu.Lock()
	delete(m.activeSessions, p.ID)
	m.mu.Unlock()
}

// runMultiplierStackTimeout 超時後結算（25 秒後未達到 10.0x）
func (g *Game) runMultiplierStackTimeout(p *player.Player, sess *multiplierStackSession) {
	timer := time.NewTimer(LuckyMultiplierStackDuration)
	defer timer.Stop()

	<-timer.C

	m := g.LuckyMultiplierStack
	m.mu.Lock()
	// 確認 session 仍然存在且未爆發
	currentSess, ok := m.activeSessions[p.ID]
	if !ok || currentSess.burst {
		m.mu.Unlock()
		return
	}
	finalStack := currentSess.stack
	killCount := currentSess.killCount
	totalReward := currentSess.totalReward
	delete(m.activeSessions, p.ID)
	m.mu.Unlock()

	log.Printf("[MultiplierStack] player=%s 超時結算！最終疊加 %.1fx，擊破 %d 個，總獎勵 %d",
		p.ID, finalStack, killCount, totalReward)

	// 個人結算訊息
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyMultiplierStack,
		Payload: ws.LuckyMultiplierStackPayload{
			Event:       "stack_settle",
			PlayerID:    p.ID,
			FinalStack:  finalStack,
			KillCount:   killCount,
			TotalReward: totalReward,
		},
	})
}
