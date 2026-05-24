// lucky_fortune_prophecy_handler.go — 幸運命運預言魚系統（DAY-274）
// 業界依據：Lucky Fish by AbraCadabra（2026-05-16）crash mechanic + 倍率上升機制
//   進化版：「預言倍率門檻」讓玩家有「這次能達到預言嗎？」的期待感
//
// 設計：
//   擊破 T232 後，Server 預言「下一條被擊破的魚倍率門檻」（預言值 = 隨機 ×2.0 到 ×8.0）
//   玩家在 20 秒內擊破任何目標：
//     若實際倍率 ≥ 預言值 → 「預言成真」×3.0 加成（個人）
//     若實際倍率 < 預言值 → 「預言落空」×1.2 安慰獎（個人）
//   預言值越高，成真機率越低但更有挑戰感
//   個人冷卻 22 秒；全服冷卻 38 秒
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

// prophecySession 命運預言 session（個人）
type prophecySession struct {
	playerID    string
	playerName  string
	prophecyMult float64 // 預言倍率門檻（2.0 ~ 8.0）
	expiresAt   time.Time
	used        bool
}

// luckyFortuneProphecyManager 幸運命運預言魚系統管理器
type luckyFortuneProphecyManager struct {
	mu             sync.Mutex
	personalCooldowns map[string]time.Time // playerID → 冷卻到期時間
	globalCooldown    time.Time
	activeSessions    map[string]*prophecySession // playerID → session
}

func newLuckyFortuneProphecyManager() *luckyFortuneProphecyManager {
	return &luckyFortuneProphecyManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*prophecySession),
	}
}

// isLuckyFortuneProphecyFish 判斷是否為幸運命運預言魚
func isLuckyFortuneProphecyFish(defID string) bool {
	return defID == "T232"
}

// isFortuneProphecySessionActive 判斷玩家是否有進行中的預言 session（供 handleKill 使用）
func (m *luckyFortuneProphecyManager) isFortuneProphecySessionActive(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	if !ok {
		return false
	}
	if sess.used || time.Now().After(sess.expiresAt) {
		delete(m.activeSessions, playerID)
		return false
	}
	return true
}

// getFortuneProphecyMult 取得預言倍率加成（供 handleKill 使用）
// 回傳 (mult, consumed)：mult 是倍率加成，consumed 表示 session 已消耗
func (m *luckyFortuneProphecyManager) getFortuneProphecyMult(playerID string, actualMult float64) (float64, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	if !ok || sess.used || time.Now().After(sess.expiresAt) {
		delete(m.activeSessions, playerID)
		return 1.0, false
	}
	sess.used = true
	delete(m.activeSessions, playerID)
	if actualMult >= sess.prophecyMult {
		return 3.0, true // 預言成真 ×3.0
	}
	return 1.2, true // 預言落空 ×1.2 安慰獎
}

// rollProphecyMult 隨機抽取預言倍率門檻（2.0 ~ 8.0，分 6 個等級）
func rollProphecyMult() float64 {
	// 等級分布：2.0(30%) / 3.0(25%) / 4.0(20%) / 5.0(15%) / 6.0(7%) / 8.0(3%)
	r := rand.Float64()
	switch {
	case r < 0.30:
		return 2.0
	case r < 0.55:
		return 3.0
	case r < 0.75:
		return 4.0
	case r < 0.90:
		return 5.0
	case r < 0.97:
		return 6.0
	default:
		return 8.0
	}
}

// tryLuckyFortuneProphecyFish 擊破 T232 後觸發命運預言
func (g *Game) tryLuckyFortuneProphecyFish(p *player.Player) {
	m := g.LuckyFortuneProphecy
	now := time.Now()

	m.mu.Lock()
	// 個人冷卻檢查（22 秒）
	if cd, ok := m.personalCooldowns[p.ID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	// 全服冷卻檢查（38 秒）
	if now.Before(m.globalCooldown) {
		m.mu.Unlock()
		return
	}
	// 已有進行中 session
	if sess, ok := m.activeSessions[p.ID]; ok && !sess.used && now.Before(sess.expiresAt) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(22 * time.Second)
	m.globalCooldown = now.Add(38 * time.Second)

	// 抽取預言倍率門檻
	prophecyMult := rollProphecyMult()

	// 建立 session
	sess := &prophecySession{
		playerID:    p.ID,
		playerName:  p.DisplayName,
		prophecyMult: prophecyMult,
		expiresAt:   now.Add(20 * time.Second),
		used:        false,
	}
	m.activeSessions[p.ID] = sess
	m.mu.Unlock()

	log.Printf("[FortuneProphecy] player=%s prophecy=×%.1f expires=20s", p.ID, prophecyMult)

	// 發送個人預言通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyFortuneProphecy,
		Payload: ws.LuckyFortuneProphecyPayload{
			Event:        "prophecy_start",
			PlayerID:     p.ID,
			PlayerName:   p.DisplayName,
			PredictedMult: prophecyMult,
			ExpiresIn:    20,
		},
	})

	// 全服廣播（讓其他玩家知道有人觸發了命運預言）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyFortuneProphecy,
		Payload: ws.LuckyFortuneProphecyPayload{
			Event:        "prophecy_broadcast",
			PlayerID:     p.ID,
			PlayerName:   p.DisplayName,
			PredictedMult: prophecyMult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyFortuneProphecy, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🔮 %s 觸發命運預言！門檻 ×%.1f，20 秒內能成真嗎？", p.DisplayName, prophecyMult),
	})

	// 啟動超時 goroutine（20 秒後若未使用，發送超時通知）
	go g.runProphecyTimeout(p.ID, p.DisplayName, prophecyMult, sess.expiresAt)
}

// notifyFortuneProphecyKill 玩家在預言期間擊破目標時，判斷預言結果並通知
// 回傳倍率加成（1.0 = 無加成，1.2 = 落空，3.0 = 成真）
func (g *Game) notifyFortuneProphecyKill(p *player.Player, targetName string, actualMult float64, baseReward int) float64 {
	mult, consumed := g.LuckyFortuneProphecy.getFortuneProphecyMult(p.ID, actualMult)
	if !consumed {
		return 1.0
	}

	isFulfilled := mult >= 3.0
	bonusReward := int(float64(baseReward) * (mult - 1.0))

	if isFulfilled {
		// 預言成真
		log.Printf("[FortuneProphecy] FULFILLED player=%s actual=×%.1f bonus=%d", p.ID, actualMult, bonusReward)
		_ = g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyFortuneProphecy,
			Payload: ws.LuckyFortuneProphecyPayload{
				Event:        "prophecy_fulfilled",
				PlayerID:     p.ID,
				PlayerName:   p.DisplayName,
				ActualMult:   actualMult,
				ResultMult:   mult,
				BonusReward:  bonusReward,
				TargetName:   targetName,
			},
		})
		// 全服廣播（預言成真）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyFortuneProphecy,
			Payload: ws.LuckyFortuneProphecyPayload{
				Event:        "prophecy_fulfilled_broadcast",
				PlayerID:     p.ID,
				PlayerName:   p.DisplayName,
				ActualMult:   actualMult,
				ResultMult:   mult,
				BonusReward:  bonusReward,
			},
		})
		// 全服公告（預言成真）
		g.Announce.Create(announce.EventLuckyFortuneProphecy, p.DisplayName, bonusReward, map[string]string{
			"message": fmt.Sprintf("🔮 %s 預言成真！×%.1f 命中！獎勵 +%d！", p.DisplayName, actualMult, bonusReward),
			"color":   "#FFD700",
		})
	} else {
		// 預言落空
		log.Printf("[FortuneProphecy] FAILED player=%s actual=×%.1f consolation=%d", p.ID, actualMult, bonusReward)
		_ = g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyFortuneProphecy,
			Payload: ws.LuckyFortuneProphecyPayload{
				Event:        "prophecy_failed",
				PlayerID:     p.ID,
				PlayerName:   p.DisplayName,
				ActualMult:   actualMult,
				ResultMult:   mult,
				BonusReward:  bonusReward,
				TargetName:   targetName,
			},
		})
	}

	return mult
}

// runProphecyTimeout 20 秒後若 session 未使用，發送超時通知
func (g *Game) runProphecyTimeout(playerID, playerName string, prophecyMult float64, expiresAt time.Time) {
	remaining := time.Until(expiresAt)
	if remaining > 0 {
		time.Sleep(remaining)
	}

	m := g.LuckyFortuneProphecy
	m.mu.Lock()
	sess, ok := m.activeSessions[playerID]
	if !ok || sess.used {
		m.mu.Unlock()
		return
	}
	delete(m.activeSessions, playerID)
	m.mu.Unlock()

	// 發送超時通知
	_ = g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgLuckyFortuneProphecy,
		Payload: ws.LuckyFortuneProphecyPayload{
			Event:        "prophecy_expire",
			PlayerID:     playerID,
			PlayerName:   playerName,
			PredictedMult: prophecyMult,
		},
	})
	log.Printf("[FortuneProphecy] expired player=%s prophecy=×%.1f", playerID, prophecyMult)
}
