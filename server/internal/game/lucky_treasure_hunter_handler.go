// lucky_treasure_hunter_handler.go — 幸運寶藏獵人魚系統（DAY-260）
// 業界原創「寶藏地圖碎片+挖掘+寶藏爆發」機制
//
// 設計：擊破 T218 後，Server 為觸發玩家啟動「寶藏獵人模式」（持續 20 秒）：
//   - 玩家每次擊破任何目標，有 30% 機率「發現碎片」（個人獎勵 ×1.8）
//   - 集齊 3 個碎片 → 「寶藏爆發」：×5.0 倍率大獎（個人）
//   - 20 秒後未集齊 → 「寶藏消失」：已收集碎片數 × ×1.2 安慰獎（個人）
//   - 個人冷卻 30 秒；全服冷卻 48 秒
//
// 設計差異：
//   - 與星座命運（T217，命運分配+標記目標）不同，寶藏獵人是「個人探索」，
//     讓玩家有「每一槍都可能發現寶藏」的期待感
//   - 「30% 機率發現碎片」讓玩家有「要趁 20 秒內多打幾條魚」的緊迫感
//   - 「集齊 3 個碎片 ×5.0 大獎」是目前個人類最高倍率，製造「哇，集齊了！」的爽感
//   - 「安慰獎 ×1.2 × 碎片數」確保即使沒集齊也有收益，降低挫敗感
//   - 全服廣播「有人觸發了寶藏獵人」讓其他玩家看到，製造羨慕感
//   - 全服廣播「有人集齊寶藏爆發了」讓所有玩家看到，製造「我也想觸發」的動機
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyTreasureHunterPersonalCD   = 30 * time.Second // 個人冷卻
	LuckyTreasureHunterGlobalCD     = 48 * time.Second // 全服冷卻
	LuckyTreasureHunterDuration     = 20 * time.Second // 寶藏獵人模式時限
	LuckyTreasureHunterFindChance   = 0.30             // 每次擊破發現碎片機率
	LuckyTreasureHunterFragMult     = 1.8              // 發現碎片倍率（個人）
	LuckyTreasureHunterBurstMult    = 5.0              // 寶藏爆發倍率（個人）
	LuckyTreasureHunterConsoleMult  = 1.2              // 安慰獎倍率（個人，每個碎片）
	LuckyTreasureHunterFragTarget   = 3                // 集齊碎片數
)

// treasureHunterSession 寶藏獵人會話
type treasureHunterSession struct {
	playerID   string
	playerName string
	expiresAt  time.Time
	fragments  int // 已收集碎片數
	mu         sync.Mutex
}

// luckyTreasureHunterManager 幸運寶藏獵人魚管理器
type luckyTreasureHunterManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前活躍的寶藏獵人會話（playerID → session）
	activeSessions map[string]*treasureHunterSession
}

func newLuckyTreasureHunterManager() *luckyTreasureHunterManager {
	return &luckyTreasureHunterManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*treasureHunterSession),
	}
}

// isLuckyTreasureHunterFish 判斷是否為幸運寶藏獵人魚
func isLuckyTreasureHunterFish(defID string) bool {
	return defID == "T218"
}

// isTreasureHunterActive 判斷玩家是否在寶藏獵人模式中
func (m *luckyTreasureHunterManager) isTreasureHunterActive(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	if !ok {
		return false
	}
	return time.Now().Before(sess.expiresAt)
}

// tryLuckyTreasureHunterFish 擊破 T218 後觸發寶藏獵人模式
func (g *Game) tryLuckyTreasureHunterFish(p *player.Player) {
	m := g.LuckyTreasureHunter

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
	// 已有活躍會話
	if sess, ok := m.activeSessions[p.ID]; ok && now.Before(sess.expiresAt) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyTreasureHunterPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyTreasureHunterGlobalCD)

	sess := &treasureHunterSession{
		playerID:   p.ID,
		playerName: p.DisplayName,
		expiresAt:  now.Add(LuckyTreasureHunterDuration),
		fragments:  0,
	}
	m.activeSessions[p.ID] = sess
	m.mu.Unlock()

	log.Printf("[TreasureHunter] player=%s 觸發寶藏獵人模式！時限 %ds，集齊 %d 個碎片→×%.1f 大獎！",
		p.ID, int(LuckyTreasureHunterDuration.Seconds()), LuckyTreasureHunterFragTarget, LuckyTreasureHunterBurstMult)

	// 個人訊息
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyTreasureHunter,
		Payload: ws.LuckyTreasureHunterPayload{
			Event:       "treasure_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyTreasureHunterDuration.Seconds()),
			FragTarget:  LuckyTreasureHunterFragTarget,
			FragMult:    LuckyTreasureHunterFragMult,
			BurstMult:   LuckyTreasureHunterBurstMult,
			FindChance:  LuckyTreasureHunterFindChance,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTreasureHunter,
		Payload: ws.LuckyTreasureHunterPayload{
			Event:      "treasure_broadcast",
			PlayerName: p.DisplayName,
			FragTarget: LuckyTreasureHunterFragTarget,
			BurstMult:  LuckyTreasureHunterBurstMult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyTreasureHunter, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🗺️ %s 觸發寶藏獵人！集齊 %d 個碎片→×%.1f 大獎！",
			p.DisplayName, LuckyTreasureHunterFragTarget, LuckyTreasureHunterBurstMult),
		"color": "#D4A017",
	})

	// 啟動超時 goroutine
	go g.runTreasureHunterTimeout(sess)
}

// notifyTreasureHunterKill 寶藏獵人模式中擊破目標（由 handleKill 呼叫）
func (g *Game) notifyTreasureHunterKill(p *player.Player, t *target.Target) {
	m := g.LuckyTreasureHunter

	m.mu.Lock()
	sess, ok := m.activeSessions[p.ID]
	if !ok || time.Now().After(sess.expiresAt) {
		m.mu.Unlock()
		return
	}
	m.mu.Unlock()

	// 30% 機率發現碎片
	if rand.Float64() >= LuckyTreasureHunterFindChance {
		return
	}

	sess.mu.Lock()
	sess.fragments++
	fragments := sess.fragments
	sess.mu.Unlock()

	avgBet := g.getAvgBetCost()
	fragReward := int(float64(avgBet) * t.Multiplier * LuckyTreasureHunterFragMult)
	if fragReward < 1 {
		fragReward = 1
	}
	p.AddReward(fragReward)

	log.Printf("[TreasureHunter] player=%s 發現碎片 %d/%d！獎勵 %d（×%.1f）",
		p.ID, fragments, LuckyTreasureHunterFragTarget, fragReward, LuckyTreasureHunterFragMult)

	// 個人通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyTreasureHunter,
		Payload: ws.LuckyTreasureHunterPayload{
			Event:      "treasure_fragment",
			PlayerName: p.DisplayName,
			Fragments:  fragments,
			FragTarget: LuckyTreasureHunterFragTarget,
			FragMult:   LuckyTreasureHunterFragMult,
			Reward:     fragReward,
			TargetName: t.Def.Name,
		},
	})

	// 集齊碎片 → 寶藏爆發
	if fragments >= LuckyTreasureHunterFragTarget {
		go g.doTreasureHunterBurst(sess, p)
	}
}

// doTreasureHunterBurst 寶藏爆發（集齊碎片）
func (g *Game) doTreasureHunterBurst(sess *treasureHunterSession, p *player.Player) {
	m := g.LuckyTreasureHunter

	m.mu.Lock()
	// 確認會話仍然有效
	if _, ok := m.activeSessions[p.ID]; !ok {
		m.mu.Unlock()
		return
	}
	delete(m.activeSessions, p.ID)
	m.mu.Unlock()

	avgBet := g.getAvgBetCost()
	burstReward := int(float64(avgBet) * float64(LuckyTreasureHunterFragTarget) * LuckyTreasureHunterBurstMult * 5)
	if burstReward < 1 {
		burstReward = 1
	}
	p.AddReward(burstReward)

	log.Printf("[TreasureHunter] player=%s 寶藏爆發！集齊 %d 個碎片！獎勵 %d（×%.1f）",
		p.ID, LuckyTreasureHunterFragTarget, burstReward, LuckyTreasureHunterBurstMult)

	// 個人通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyTreasureHunter,
		Payload: ws.LuckyTreasureHunterPayload{
			Event:      "treasure_burst",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Fragments:  LuckyTreasureHunterFragTarget,
			FragTarget: LuckyTreasureHunterFragTarget,
			BurstMult:  LuckyTreasureHunterBurstMult,
			Reward:     burstReward,
		},
	})

	// 全服廣播寶藏爆發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTreasureHunter,
		Payload: ws.LuckyTreasureHunterPayload{
			Event:      "treasure_burst_broadcast",
			PlayerName: p.DisplayName,
			BurstMult:  LuckyTreasureHunterBurstMult,
			Reward:     burstReward,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyTreasureHunter, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🏆 %s 集齊寶藏！×%.1f 大獎！獎勵 +%d！",
			p.DisplayName, LuckyTreasureHunterBurstMult, burstReward),
		"color": "#FFD700",
	})
}

// runTreasureHunterTimeout 寶藏獵人超時處理
func (g *Game) runTreasureHunterTimeout(sess *treasureHunterSession) {
	timer := time.NewTimer(LuckyTreasureHunterDuration)
	defer timer.Stop()

	select {
	case <-timer.C:
		m := g.LuckyTreasureHunter

		m.mu.Lock()
		// 確認會話仍然有效（可能已被 burst 清除）
		if _, ok := m.activeSessions[sess.playerID]; !ok {
			m.mu.Unlock()
			return
		}
		delete(m.activeSessions, sess.playerID)
		m.mu.Unlock()

		sess.mu.Lock()
		fragments := sess.fragments
		sess.mu.Unlock()

		log.Printf("[TreasureHunter] player=%s 寶藏消失！收集了 %d/%d 個碎片",
			sess.playerID, fragments, LuckyTreasureHunterFragTarget)

		// 安慰獎（已收集碎片數 × ×1.2）
		if fragments > 0 {
			avgBet := g.getAvgBetCost()
			consoleReward := int(float64(avgBet) * float64(fragments) * LuckyTreasureHunterConsoleMult * 3)
			if consoleReward < 1 {
				consoleReward = 1
			}

			// 找到玩家
			g.mu.RLock()
			p, ok := g.Players[sess.playerID]
			g.mu.RUnlock()
			if ok {
				p.AddReward(consoleReward)
			}

			_ = g.Hub.Send(sess.playerID, &ws.Message{
				Type: ws.MsgLuckyTreasureHunter,
				Payload: ws.LuckyTreasureHunterPayload{
					Event:      "treasure_timeout",
					PlayerName: sess.playerName,
					Fragments:  fragments,
					FragTarget: LuckyTreasureHunterFragTarget,
					Reward:     consoleReward,
				},
			})
		} else {
			_ = g.Hub.Send(sess.playerID, &ws.Message{
				Type: ws.MsgLuckyTreasureHunter,
				Payload: ws.LuckyTreasureHunterPayload{
					Event:      "treasure_timeout",
					PlayerName: sess.playerName,
					Fragments:  0,
					FragTarget: LuckyTreasureHunterFragTarget,
					Reward:     0,
				},
			})
		}

	case <-g.stopCh:
		return
	}
}
