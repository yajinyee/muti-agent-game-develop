// lucky_server_charge_handler.go — 幸運全服充能魚系統（DAY-256）
// 業界原創「全服共同充能→全服大爆發」機制
//
// 設計：擊破 T214 後，全服所有玩家共同累積「充能值」（每次任何玩家擊破任何目標 +1）：
//   - 充能值達到 20 時「全服大爆發」：全場所有目標 100% 擊破（×2.0 倍率，全服共享）
//   - 若 30 秒內未達到 20 → 「充能失敗」：已累積充能值 × 0.5 倍率（安慰獎，全服共享）
//   - 每次充能 +1 時廣播進度（讓全服玩家看到「還差幾個」）
//   - 個人冷卻 30 秒；全服冷卻 50 秒
//
// 設計差異：
//   - 與量子共鳴（T209，2個玩家1.5秒內協作）不同，全服充能是「所有玩家長期協作」，
//     讓玩家有「全服一起努力打魚，累積到20個就爆發」的集體感
//   - 「每次充能廣播進度」讓玩家即時看到「還差幾個」，製造「快了快了」的緊迫感
//   - 「充能失敗安慰獎」確保即使沒達到目標也有收益，降低挫敗感
//   - 「×2.0 全場 100% 擊破」是目前全服合力類最強的爆發效果
//   - 觸發玩家獲得「充能先鋒」稱號廣播，製造「是我開啟了這次充能」的成就感
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyServerChargePersonalCD  = 30 * time.Second // 個人冷卻
	LuckyServerChargeGlobalCD    = 50 * time.Second // 全服冷卻
	LuckyServerChargeDuration    = 30 * time.Second // 充能時限
	LuckyServerChargeTarget      = 20               // 充能目標值
	LuckyServerChargeBurstMult   = 2.0              // 大爆發倍率（全服共享）
	LuckyServerChargeFailMult    = 0.5              // 充能失敗安慰獎倍率係數
)

// serverChargeSession 全服充能會話
type serverChargeSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	expiresAt         time.Time
	chargeCount       int
	mu                sync.Mutex
}

// luckyServerChargeManager 幸運全服充能魚管理器
type luckyServerChargeManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前活躍的充能會話（nil = 無）
	activeSession *serverChargeSession
}

func newLuckyServerChargeManager() *luckyServerChargeManager {
	return &luckyServerChargeManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyServerChargeFish 判斷是否為幸運全服充能魚
func isLuckyServerChargeFish(defID string) bool {
	return defID == "T214"
}

// notifyServerChargeKill 任何玩家擊破任何目標時，若充能進行中則累積充能值
// 由 handleKill 呼叫（非 T214 目標）
func (g *Game) notifyServerChargeKill() {
	m := g.LuckyServerCharge
	m.mu.Lock()
	sess := m.activeSession
	if sess == nil {
		m.mu.Unlock()
		return
	}
	now := time.Now()
	if now.After(sess.expiresAt) {
		m.activeSession = nil
		m.mu.Unlock()
		return
	}
	m.mu.Unlock()

	sess.mu.Lock()
	sess.chargeCount++
	current := sess.chargeCount
	sess.mu.Unlock()

	log.Printf("[ServerCharge] 充能 +1！目前 %d/%d", current, LuckyServerChargeTarget)

	// 廣播充能進度
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyServerCharge,
		Payload: ws.LuckyServerChargePayload{
			Event:        "charge_progress",
			PlayerName:   sess.triggerPlayerName,
			ChargeCount:  current,
			ChargeTarget: LuckyServerChargeTarget,
		},
	})

	// 達到目標 → 觸發大爆發
	if current >= LuckyServerChargeTarget {
		go g.doServerChargeBurst(sess)
	}
}

// tryLuckyServerChargeFish 擊破 T214 後觸發全服充能
func (g *Game) tryLuckyServerChargeFish(p *player.Player) {
	m := g.LuckyServerCharge

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
	// 已有活躍充能
	if m.activeSession != nil && now.Before(m.activeSession.expiresAt) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyServerChargePersonalCD)
	m.globalCooldownUntil = now.Add(LuckyServerChargeGlobalCD)

	expiresAt := now.Add(LuckyServerChargeDuration)
	sess := &serverChargeSession{
		triggerPlayerID:   p.ID,
		triggerPlayerName: p.DisplayName,
		expiresAt:         expiresAt,
		chargeCount:       0,
	}
	m.activeSession = sess
	m.mu.Unlock()

	log.Printf("[ServerCharge] player=%s 觸發全服充能！目標 %d 次，時限 %ds",
		p.ID, LuckyServerChargeTarget, int(LuckyServerChargeDuration.Seconds()))

	// 個人訊息：充能先鋒
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyServerCharge,
		Payload: ws.LuckyServerChargePayload{
			Event:        "charge_start",
			PlayerID:     p.ID,
			PlayerName:   p.DisplayName,
			DurationSec:  int(LuckyServerChargeDuration.Seconds()),
			ChargeTarget: LuckyServerChargeTarget,
			BurstMult:    LuckyServerChargeBurstMult,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyServerCharge,
		Payload: ws.LuckyServerChargePayload{
			Event:        "charge_broadcast",
			PlayerName:   p.DisplayName,
			DurationSec:  int(LuckyServerChargeDuration.Seconds()),
			ChargeTarget: LuckyServerChargeTarget,
			BurstMult:    LuckyServerChargeBurstMult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyServerCharge, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⚡ %s 開啟全服充能！全服一起打魚累積 %d 次→大爆發！×%.1f 全服共享！",
			p.DisplayName, LuckyServerChargeTarget, LuckyServerChargeBurstMult),
		"color": "#FF8C00",
	})

	// 啟動超時 goroutine
	go g.runServerChargeTimeout(sess)
}

// runServerChargeTimeout 充能超時處理
func (g *Game) runServerChargeTimeout(sess *serverChargeSession) {
	timer := time.NewTimer(LuckyServerChargeDuration)
	defer timer.Stop()

	select {
	case <-timer.C:
		// 超時：檢查是否已爆發
		m := g.LuckyServerCharge
		m.mu.Lock()
		if m.activeSession != sess {
			// 已被爆發清除
			m.mu.Unlock()
			return
		}
		m.activeSession = nil
		m.mu.Unlock()

		sess.mu.Lock()
		finalCount := sess.chargeCount
		sess.mu.Unlock()

		// 充能失敗：安慰獎
		g.doServerChargeFail(sess, finalCount)

	case <-g.stopCh:
		return
	}
}

// doServerChargeBurst 全服大爆發（充能達到目標）
func (g *Game) doServerChargeBurst(sess *serverChargeSession) {
	// 清除活躍會話
	m := g.LuckyServerCharge
	m.mu.Lock()
	if m.activeSession != sess {
		m.mu.Unlock()
		return
	}
	m.activeSession = nil
	m.mu.Unlock()

	// 全場所有目標 100% 擊破
	g.mu.Lock()
	var toKill []*target.Target
	for _, t := range g.Targets {
		if t.IsAlive {
			toKill = append(toKill, t)
		}
	}
	for _, t := range toKill {
		t.IsAlive = false
		delete(g.Targets, t.InstanceID)
	}
	g.mu.Unlock()

	killCount := len(toKill)

	// 計算全服共享獎勵
	avgBet := g.getAvgBetCost()
	totalReward := 0
	for _, t := range toKill {
		reward := int(float64(avgBet) * t.Multiplier * LuckyServerChargeBurstMult)
		if reward < 1 {
			reward = 1
		}
		totalReward += reward
	}
	if totalReward > 0 {
		g.distributeRewardToAll(totalReward)
	}

	log.Printf("[ServerCharge] 全服大爆發！擊破 %d 個目標，全服獎勵 %d（充能 %d 次）",
		killCount, totalReward, LuckyServerChargeTarget)

	// 全服廣播大爆發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyServerCharge,
		Payload: ws.LuckyServerChargePayload{
			Event:        "charge_burst",
			PlayerName:   sess.triggerPlayerName,
			KillCount:    killCount,
			BurstMult:    LuckyServerChargeBurstMult,
			TotalReward:  totalReward,
			ChargeCount:  LuckyServerChargeTarget,
			ChargeTarget: LuckyServerChargeTarget,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyServerCharge, sess.triggerPlayerName, 0, map[string]string{
		"message": fmt.Sprintf("⚡ 全服充能成功！擊破 %d 個目標！全服獎勵 +%d！×%.1f 大爆發！",
			killCount, totalReward, LuckyServerChargeBurstMult),
		"color": "#FFD700",
	})
}

// doServerChargeFail 充能失敗（超時未達目標）
func (g *Game) doServerChargeFail(sess *serverChargeSession, finalCount int) {
	if finalCount == 0 {
		// 完全沒有充能，不給安慰獎
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyServerCharge,
			Payload: ws.LuckyServerChargePayload{
				Event:        "charge_fail",
				PlayerName:   sess.triggerPlayerName,
				ChargeCount:  0,
				ChargeTarget: LuckyServerChargeTarget,
				TotalReward:  0,
			},
		})
		return
	}

	// 安慰獎：已累積充能值 × 0.5 倍率
	avgBet := g.getAvgBetCost()
	consolationReward := int(float64(avgBet) * float64(finalCount) * LuckyServerChargeFailMult)
	if consolationReward < 1 {
		consolationReward = 1
	}
	g.distributeRewardToAll(consolationReward)

	log.Printf("[ServerCharge] 充能失敗！充能 %d/%d 次，安慰獎 %d",
		finalCount, LuckyServerChargeTarget, consolationReward)

	// 全服廣播失敗
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyServerCharge,
		Payload: ws.LuckyServerChargePayload{
			Event:        "charge_fail",
			PlayerName:   sess.triggerPlayerName,
			ChargeCount:  finalCount,
			ChargeTarget: LuckyServerChargeTarget,
			TotalReward:  consolationReward,
		},
	})

	// 全服公告（充能超過一半才公告）
	if finalCount >= LuckyServerChargeTarget/2 {
		g.Announce.Create(announce.EventLuckyServerCharge, sess.triggerPlayerName, 0, map[string]string{
			"message": fmt.Sprintf("⚡ 充能失敗！累積 %d/%d 次，安慰獎 +%d！下次加油！",
				finalCount, LuckyServerChargeTarget, consolationReward),
			"color": "#808080",
		})
	}
}
