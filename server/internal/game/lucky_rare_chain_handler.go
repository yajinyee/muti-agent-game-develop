// lucky_rare_chain_handler.go — 幸運連鎖稀有魚系統（DAY-280）
// 業界依據：Fishing Fortune 2026「Chain rare catches within 90-second windows」機制進化版
// 業界原創「稀有連鎖+倍率爬升+時間視窗」機制
//
// 設計：擊破 T238 後，觸發「稀有連鎖模式」（持續 20 秒）
//   - 模式期間，玩家每次擊破「稀有或以上目標」（倍率 ≥ 15x），連鎖計數 +1
//   - 倍率爬升：第1層×1.5 → 第2層×2.5 → 第3層×4.0 → 第4層×6.0 → 第5層×10.0
//   - 每次連鎖必須在 8 秒內完成（否則連鎖中斷，重置到第1層）
//   - 達到第 5 層（×10.0）觸發「連鎖爆發」：個人大獎 + 全服廣播
//   - 20 秒後模式結束，廣播最終連鎖層數和總獎勵
//   - 個人冷卻 22 秒；全服冷卻 38 秒
//
// 設計差異：
//   - 與倍率疊加（T225，每次擊破任何目標 +0.3x）不同，連鎖稀有是「只有稀有目標才能連鎖」
//     讓玩家有「要找稀有魚打，普通魚不算」的策略感
//   - 「8 秒連鎖視窗」讓玩家有「要趕快找下一條稀有魚」的緊迫感
//   - 「5 層倍率爬升（×1.5→×10.0）」讓玩家有「越打越高，要撐到第 5 層」的動力
//   - 「連鎖中斷重置到第1層」讓玩家有「不能讓連鎖斷掉」的壓力感
//   - 「第 5 層連鎖爆發全服廣播」讓所有玩家看到「有人達到 5 層連鎖了」，製造羨慕感
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
	LuckyRareChainPersonalCD  = 22 * time.Second // 個人冷卻
	LuckyRareChainGlobalCD    = 38 * time.Second // 全服冷卻
	LuckyRareChainDuration    = 20 * time.Second // 模式持續時間
	LuckyRareChainWindowSec   = 8 * time.Second  // 連鎖視窗（8 秒內必須擊破下一條稀有魚）
	LuckyRareChainRareMinMult = 15               // 稀有目標最低倍率門檻
	LuckyRareChainMaxLayer    = 5                // 最大連鎖層數
)

// 各層倍率
var rareChainLayerMults = []float64{0, 1.5, 2.5, 4.0, 6.0, 10.0} // index 0 未使用，1-5 對應各層

// rareChainSession 連鎖稀有模式會話
type rareChainSession struct {
	playerID   string
	playerName string
	layer      int       // 當前連鎖層數（0=未開始，1-5）
	expiresAt  time.Time // 模式結束時間
	chainUntil time.Time // 連鎖視窗截止時間（8 秒）
	totalReward int      // 累積獎勵
	burst      bool      // 是否已觸發爆發
	settled    bool      // 是否已結算
}

// luckyRareChainManager 幸運連鎖稀有魚管理器
type luckyRareChainManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 活躍會話（playerID → session）
	activeSessions map[string]*rareChainSession
}

func newLuckyRareChainManager() *luckyRareChainManager {
	return &luckyRareChainManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*rareChainSession),
	}
}

// isLuckyRareChainFish 判斷是否為幸運連鎖稀有魚
func isLuckyRareChainFish(defID string) bool {
	return defID == "T238"
}

// isRareTarget 判斷是否為稀有目標（倍率 ≥ 15x）
func isRareTarget(multMin int) bool {
	return multMin >= LuckyRareChainRareMinMult
}

// getRareChainMult 取得連鎖倍率（供 handleKill 使用）
// 回傳 (mult, isChainKill)
func (m *luckyRareChainManager) getRareChainMult(playerID string, targetMultMin int) (float64, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.activeSessions[playerID]
	if !ok || session.settled || session.burst {
		return 1.0, false
	}

	now := time.Now()

	// 模式已過期
	if now.After(session.expiresAt) {
		return 1.0, false
	}

	// 不是稀有目標，不計入連鎖
	if targetMultMin < LuckyRareChainRareMinMult {
		return 1.0, false
	}

	// 連鎖視窗已過期（連鎖中斷）
	if session.layer > 0 && now.After(session.chainUntil) {
		// 連鎖中斷，重置到第1層
		session.layer = 0
	}

	// 連鎖 +1
	session.layer++
	if session.layer > LuckyRareChainMaxLayer {
		session.layer = LuckyRareChainMaxLayer
	}

	// 更新連鎖視窗
	session.chainUntil = now.Add(LuckyRareChainWindowSec)

	mult := rareChainLayerMults[session.layer]
	return mult, true
}

// recordRareChainReward 記錄連鎖獎勵（供 handleKill 使用）
func (m *luckyRareChainManager) recordRareChainReward(playerID string, reward int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if session, ok := m.activeSessions[playerID]; ok && !session.settled {
		session.totalReward += reward
	}
}

// isRareChainBurst 判斷是否達到爆發條件（第5層）
func (m *luckyRareChainManager) isRareChainBurst(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if session, ok := m.activeSessions[playerID]; ok && !session.settled && !session.burst {
		return session.layer >= LuckyRareChainMaxLayer
	}
	return false
}

// markRareChainBurst 標記爆發已觸發
func (m *luckyRareChainManager) markRareChainBurst(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if session, ok := m.activeSessions[playerID]; ok {
		session.burst = true
	}
}

// tryLuckyRareChainFish 擊破 T238 後觸發連鎖稀有模式（供 handleKill 使用）
func (g *Game) tryLuckyRareChainFish(p *player.Player) {
	mgr := g.LuckyRareChain
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
	if _, ok := mgr.activeSessions[p.ID]; ok {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = now.Add(LuckyRareChainPersonalCD)
	mgr.globalCooldownUntil = now.Add(LuckyRareChainGlobalCD)

	session := &rareChainSession{
		playerID:   p.ID,
		playerName: p.DisplayName,
		layer:      0,
		expiresAt:  now.Add(LuckyRareChainDuration),
		chainUntil: now.Add(LuckyRareChainWindowSec),
	}
	mgr.activeSessions[p.ID] = session
	mgr.mu.Unlock()

	log.Printf("[RareChain] player=%s triggered rare chain mode, expires=%v", p.ID, session.expiresAt)

	// 個人訊息：連鎖稀有模式觸發
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyRareChain,
		Payload: ws.LuckyRareChainPayload{
			Event:      "chain_start",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Duration:   int(LuckyRareChainDuration.Seconds()),
			WindowSec:  int(LuckyRareChainWindowSec.Seconds()),
			MaxLayer:   LuckyRareChainMaxLayer,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyRareChain,
		Payload: ws.LuckyRareChainPayload{
			Event:      "chain_broadcast",
			PlayerName: p.DisplayName,
			Duration:   int(LuckyRareChainDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyRareChain, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🔗 %s 觸發連鎖稀有模式！20 秒內連鎖擊破稀有魚，倍率最高 ×10.0！",
			p.DisplayName),
		"color": "#FF6B35",
	})
	g.broadcastAnnouncement(ann)

	// 啟動模式計時器
	go g.runRareChainTimer(p, session)
}

// notifyRareChainKill 連鎖稀有模式中擊破稀有目標時通知（供 handleKill 使用）
func (g *Game) notifyRareChainKill(p *player.Player, layer int, mult float64, reward int) {
	mgr := g.LuckyRareChain
	mgr.mu.Lock()
	session, ok := mgr.activeSessions[p.ID]
	if !ok || session.settled {
		mgr.mu.Unlock()
		return
	}
	totalReward := session.totalReward
	mgr.mu.Unlock()

	log.Printf("[RareChain] player=%s chain kill! layer=%d mult=x%.1f reward=%d total=%d",
		p.ID, layer, mult, reward, totalReward)

	// 個人通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyRareChain,
		Payload: ws.LuckyRareChainPayload{
			Event:       "chain_kill",
			PlayerID:    p.ID,
			Layer:       layer,
			Mult:        mult,
			Reward:      reward,
			TotalReward: totalReward,
		},
	})

	// 第 5 層：觸發爆發
	if layer >= LuckyRareChainMaxLayer {
		go g.doRareChainBurst(p, session)
	}
}

// doRareChainBurst 連鎖爆發（達到第5層時觸發）
func (g *Game) doRareChainBurst(p *player.Player, session *rareChainSession) {
	mgr := g.LuckyRareChain
	mgr.mu.Lock()
	if session.burst || session.settled {
		mgr.mu.Unlock()
		return
	}
	session.burst = true
	totalReward := session.totalReward
	mgr.mu.Unlock()

	log.Printf("[RareChain] BURST! player=%s layer=%d totalReward=%d",
		p.ID, LuckyRareChainMaxLayer, totalReward)

	// 個人通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyRareChain,
		Payload: ws.LuckyRareChainPayload{
			Event:       "chain_burst",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			Layer:       LuckyRareChainMaxLayer,
			Mult:        rareChainLayerMults[LuckyRareChainMaxLayer],
			TotalReward: totalReward,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyRareChain,
		Payload: ws.LuckyRareChainPayload{
			Event:       "chain_burst_broadcast",
			PlayerName:  p.DisplayName,
			Layer:       LuckyRareChainMaxLayer,
			Mult:        rareChainLayerMults[LuckyRareChainMaxLayer],
			TotalReward: totalReward,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyRareChain, p.DisplayName, totalReward, map[string]string{
		"message": fmt.Sprintf("🔗 %s 達成 5 層連鎖爆發！×%.1f 大獎！總獎勵 %d！",
			p.DisplayName, rareChainLayerMults[LuckyRareChainMaxLayer], totalReward),
		"color": "#FFD700",
	})
	g.broadcastAnnouncement(ann)
}

// runRareChainTimer 連鎖稀有模式計時器（goroutine）
func (g *Game) runRareChainTimer(p *player.Player, session *rareChainSession) {
	select {
	case <-time.After(LuckyRareChainDuration):
	case <-g.stopCh:
		return
	}

	mgr := g.LuckyRareChain
	mgr.mu.Lock()
	if session.settled {
		mgr.mu.Unlock()
		return
	}
	session.settled = true
	finalLayer := session.layer
	totalReward := session.totalReward
	delete(mgr.activeSessions, p.ID)
	mgr.mu.Unlock()

	log.Printf("[RareChain] mode ended for player=%s finalLayer=%d totalReward=%d",
		p.ID, finalLayer, totalReward)

	// 個人通知：模式結束
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyRareChain,
		Payload: ws.LuckyRareChainPayload{
			Event:       "chain_end",
			PlayerID:    p.ID,
			Layer:       finalLayer,
			TotalReward: totalReward,
		},
	})
}
