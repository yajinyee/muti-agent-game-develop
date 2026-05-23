// lucky_karma_cycle_handler.go — 幸運命運輪迴魚系統（DAY-264）
// 業界原創「業力累積+命運爆發」機制
//
// 設計：擊破 T222 後，觸發「命運輪迴」（持續 20 秒）：
//   - 玩家每次擊破任何目標，累積「業力值」（每次 +1，最多 10）
//   - 業力值達到 10 → 「命運爆發」：業力值 × ×1.5 倍率（個人，最高 ×15.0）
//   - 20 秒後未達到 10 → 「業力結算」：已累積業力值 × ×1.2 倍率（個人）
//   - 個人冷卻 30 秒；全服冷卻 50 秒
//
// 設計差異：
//   - 與元素融合（T221，確定性收集三種元素）不同，命運輪迴是「累積型爆發」，
//     讓玩家有「每一槍都在累積業力，業力滿了就爆發」的宿命感
//   - 「業力值 × ×1.5 最高 ×15.0」讓玩家有「要趁 20 秒內打滿 10 個目標」的緊迫感
//   - 「業力結算 ×1.2 × 業力值」確保即使沒打滿也有收益，降低挫敗感
//   - 「業力計數器即時顯示」讓玩家看到「業力現在有幾個」，製造「快滿了！」的期待感
//   - 「全服廣播命運爆發」讓所有玩家看到「有人業力滿了爆發了」，製造羨慕感
//   - 業界依據：2026 年最熱門「命運輪迴+業力累積」主題，讓玩家有宿命感
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
	LuckyKarmaCyclePersonalCD  = 30 * time.Second // 個人冷卻
	LuckyKarmaCycleGlobalCD    = 50 * time.Second // 全服冷卻
	LuckyKarmaCycleDuration    = 20 * time.Second // 命運輪迴持續時間
	LuckyKarmaCycleMaxKarma    = 10               // 最大業力值
	LuckyKarmaCycleBurstMult   = 1.5              // 業力爆發倍率係數（業力值 × 1.5）
	LuckyKarmaCycleSettleMult  = 1.2              // 業力結算倍率係數（業力值 × 1.2）
)

// karmaCycleSession 命運輪迴 session
type karmaCycleSession struct {
	playerID   string
	playerName string
	karma      int       // 當前業力值
	expiresAt  time.Time
}

// luckyKarmaCycleManager 幸運命運輪迴魚管理器
type luckyKarmaCycleManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 活躍 session（playerID → session）
	activeSessions map[string]*karmaCycleSession
}

func newLuckyKarmaCycleManager() *luckyKarmaCycleManager {
	return &luckyKarmaCycleManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*karmaCycleSession),
	}
}

// isLuckyKarmaCycleFish 判斷是否為幸運命運輪迴魚
func isLuckyKarmaCycleFish(defID string) bool {
	return defID == "T222"
}

// isKarmaCycleActive 判斷玩家是否在命運輪迴中（供 handleKill 使用）
func (m *luckyKarmaCycleManager) isKarmaCycleActive(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.activeSessions[playerID]; ok {
		if time.Now().Before(s.expiresAt) {
			return true
		}
		delete(m.activeSessions, playerID)
	}
	return false
}

// tryLuckyKarmaCycleFish 擊破 T222 後觸發命運輪迴
func (g *Game) tryLuckyKarmaCycleFish(p *player.Player) {
	m := g.LuckyKarmaCycle

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
	// 已有 session 時不重複觸發
	if _, ok := m.activeSessions[p.ID]; ok {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyKarmaCyclePersonalCD)
	m.globalCooldownUntil = now.Add(LuckyKarmaCycleGlobalCD)

	// 建立 session
	expiresAt := now.Add(LuckyKarmaCycleDuration)
	session := &karmaCycleSession{
		playerID:   p.ID,
		playerName: p.DisplayName,
		karma:      0,
		expiresAt:  expiresAt,
	}
	m.activeSessions[p.ID] = session
	m.mu.Unlock()

	log.Printf("[KarmaCycle] player=%s 觸發命運輪迴！持續 20 秒，業力目標 10",
		p.ID)

	// 個人通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyKarmaCycle,
		Payload: ws.LuckyKarmaCyclePayload{
			Event:       "karma_start",
			TriggerName: p.DisplayName,
			Duration:    int(LuckyKarmaCycleDuration.Seconds()),
			MaxKarma:    LuckyKarmaCycleMaxKarma,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyKarmaCycle,
		Payload: ws.LuckyKarmaCyclePayload{
			Event:       "karma_broadcast",
			TriggerName: p.DisplayName,
			Duration:    int(LuckyKarmaCycleDuration.Seconds()),
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyKarmaCycle, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("☯️ %s 觸發命運輪迴！20 秒內累積 10 業力可獲得 ×%.0f 大獎！",
			p.DisplayName, float64(LuckyKarmaCycleMaxKarma)*LuckyKarmaCycleBurstMult),
		"color": "#9B59B6",
	})

	// 啟動超時 goroutine
	go g.runKarmaCycleTimeout(p, expiresAt)
}

// notifyKarmaCycleKill 玩家在命運輪迴中擊破目標（由 handleKill 呼叫）
func (g *Game) notifyKarmaCycleKill(p *player.Player) {
	m := g.LuckyKarmaCycle

	m.mu.Lock()
	session, ok := m.activeSessions[p.ID]
	if !ok || time.Now().After(session.expiresAt) {
		if ok {
			delete(m.activeSessions, p.ID)
		}
		m.mu.Unlock()
		return
	}

	session.karma++
	karma := session.karma
	m.mu.Unlock()

	log.Printf("[KarmaCycle] player=%s 業力 +1 = %d/%d",
		p.ID, karma, LuckyKarmaCycleMaxKarma)

	// 個人通知：業力更新
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyKarmaCycle,
		Payload: ws.LuckyKarmaCyclePayload{
			Event:    "karma_update",
			PlayerID: p.ID,
			Karma:    karma,
			MaxKarma: LuckyKarmaCycleMaxKarma,
		},
	})

	// 業力達到 10 → 命運爆發
	if karma >= LuckyKarmaCycleMaxKarma {
		go g.doKarmaBurst(p, karma)
	}
}

// doKarmaBurst 命運爆發（業力滿）
func (g *Game) doKarmaBurst(p *player.Player, karma int) {
	m := g.LuckyKarmaCycle

	// 移除 session
	m.mu.Lock()
	delete(m.activeSessions, p.ID)
	m.mu.Unlock()

	// 計算獎勵：業力值 × 1.5 × avgBetCost
	mult := float64(karma) * LuckyKarmaCycleBurstMult
	reward := int(mult * float64(g.getAvgBetCost()))
	if reward < 1 {
		reward = 1
	}
	p.AddReward(reward)

	log.Printf("[KarmaCycle] player=%s 命運爆發！業力=%d，倍率=×%.1f，獎勵=%d",
		p.ID, karma, mult, reward)

	// 個人通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyKarmaCycle,
		Payload: ws.LuckyKarmaCyclePayload{
			Event:    "karma_burst",
			PlayerID: p.ID,
			Karma:    karma,
			Mult:     mult,
			Reward:   reward,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyKarmaCycle,
		Payload: ws.LuckyKarmaCyclePayload{
			Event:      "karma_burst_broadcast",
			PlayerName: p.DisplayName,
			Karma:      karma,
			Mult:       mult,
			Reward:     reward,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyKarmaCycle, p.DisplayName, reward, map[string]string{
		"message": fmt.Sprintf("☯️ %s 業力滿溢！命運爆發！×%.1f 大獎 +%d！",
			p.DisplayName, mult, reward),
		"color": "#FFD700",
	})
}

// runKarmaCycleTimeout 命運輪迴超時結算
func (g *Game) runKarmaCycleTimeout(p *player.Player, expiresAt time.Time) {
	remaining := time.Until(expiresAt)
	if remaining > 0 {
		time.Sleep(remaining)
	}

	m := g.LuckyKarmaCycle

	m.mu.Lock()
	session, ok := m.activeSessions[p.ID]
	if !ok {
		// 已被 doKarmaBurst 清除（業力滿了爆發）
		m.mu.Unlock()
		return
	}
	karma := session.karma
	delete(m.activeSessions, p.ID)
	m.mu.Unlock()

	if karma == 0 {
		// 沒有累積業力，直接結束
		_ = g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyKarmaCycle,
			Payload: ws.LuckyKarmaCyclePayload{
				Event:    "karma_expire",
				PlayerID: p.ID,
				Karma:    0,
			},
		})
		return
	}

	// 業力結算：業力值 × 1.2 × avgBetCost
	mult := float64(karma) * LuckyKarmaCycleSettleMult
	reward := int(mult * float64(g.getAvgBetCost()))
	if reward < 1 {
		reward = 1
	}

	// 確認玩家仍在線
	g.mu.RLock()
	pl, ok := g.Players[p.ID]
	g.mu.RUnlock()
	if ok {
		pl.AddReward(reward)
	}

	log.Printf("[KarmaCycle] player=%s 業力結算！業力=%d，倍率=×%.1f，獎勵=%d",
		p.ID, karma, mult, reward)

	// 個人通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyKarmaCycle,
		Payload: ws.LuckyKarmaCyclePayload{
			Event:    "karma_settle",
			PlayerID: p.ID,
			Karma:    karma,
			Mult:     mult,
			Reward:   reward,
		},
	})
}
