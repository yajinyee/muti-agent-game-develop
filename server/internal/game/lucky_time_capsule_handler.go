// lucky_time_capsule_handler.go — 幸運時間膠囊魚系統（DAY-261）
// 業界原創「時間膠囊+預存獎勵+追加存入+膠囊開啟」機制
//
// 設計：擊破 T219 後，Server 為觸發玩家「封存」當前場上最高倍率目標的獎勵（×2.5 倍率）：
//   - 膠囊封存期間（15 秒），玩家每次擊破任何目標都會「追加存入」（×0.5 倍率，最多 5 次）
//   - 15 秒後「膠囊開啟」：一次性發放所有存入的獎勵（封存獎勵 + 追加獎勵）
//   - 個人冷卻 28 秒；全服冷卻 45 秒
//
// 設計差異：
//   - 與時光倒流（T205，重播過去擊破）不同，時間膠囊是「預存未來獎勵」，
//     讓玩家有「打開時間膠囊，看看裡面藏了什麼驚喜」的期待感
//   - 「封存最高倍率目標 ×2.5」讓玩家有「要趁膠囊期間找到高倍率目標」的動機
//   - 「追加存入 ×0.5 最多 5 次」讓玩家有「要趁 15 秒內多打幾條魚」的緊迫感
//   - 「15 秒後一次性開啟」製造「等待→開啟」的高潮設計，讓玩家有「期待感」
//   - 「追加存入計數器」讓玩家即時看到「膠囊裡存了幾個獎勵」，製造「快了快了」的緊迫感
//   - 全服廣播「有人觸發了時間膠囊」讓其他玩家看到，製造羨慕感
//   - 全服廣播「有人開啟了時間膠囊」讓所有玩家看到，製造「我也想觸發」的動機
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
	LuckyTimeCapsulePersonalCD  = 28 * time.Second // 個人冷卻
	LuckyTimeCapsuleGlobalCD    = 45 * time.Second // 全服冷卻
	LuckyTimeCapsuleDuration    = 15 * time.Second // 膠囊封存時限
	LuckyTimeCapsuleSealMult    = 2.5              // 封存最高倍率目標的倍率
	LuckyTimeCapsuleDepositMult = 0.5              // 追加存入倍率
	LuckyTimeCapsuleMaxDeposits = 5                // 最多追加存入次數
)

// timeCapsuleDeposit 追加存入記錄
type timeCapsuleDeposit struct {
	targetName string
	reward     int
}

// timeCapsuleSession 時間膠囊會話
type timeCapsuleSession struct {
	playerID    string
	playerName  string
	expiresAt   time.Time
	sealReward  int    // 封存獎勵（觸發時封存的最高倍率目標）
	sealTarget  string // 封存的目標名稱
	deposits    []timeCapsuleDeposit
	depositsMu  sync.Mutex
}

// luckyTimeCapsuleManager 幸運時間膠囊魚管理器
type luckyTimeCapsuleManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前活躍的時間膠囊會話（playerID → session）
	activeSessions map[string]*timeCapsuleSession
}

func newLuckyTimeCapsuleManager() *luckyTimeCapsuleManager {
	return &luckyTimeCapsuleManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*timeCapsuleSession),
	}
}

// isLuckyTimeCapsuleFish 判斷是否為幸運時間膠囊魚
func isLuckyTimeCapsuleFish(defID string) bool {
	return defID == "T219"
}

// isTimeCapsuleActive 判斷玩家是否在時間膠囊模式中
func (m *luckyTimeCapsuleManager) isTimeCapsuleActive(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	if !ok {
		return false
	}
	return time.Now().Before(sess.expiresAt)
}

// tryLuckyTimeCapsuleFish 擊破 T219 後觸發時間膠囊
func (g *Game) tryLuckyTimeCapsuleFish(p *player.Player) {
	m := g.LuckyTimeCapsule

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
	m.personalCooldowns[p.ID] = now.Add(LuckyTimeCapsulePersonalCD)
	m.globalCooldownUntil = now.Add(LuckyTimeCapsuleGlobalCD)
	m.mu.Unlock()

	// 找場上最高倍率目標（封存）
	avgBet := g.getAvgBetCost()
	sealTarget, sealMult, sealReward := g.findHighestMultTarget(avgBet, LuckyTimeCapsuleSealMult)

	m.mu.Lock()
	sess := &timeCapsuleSession{
		playerID:   p.ID,
		playerName: p.DisplayName,
		expiresAt:  now.Add(LuckyTimeCapsuleDuration),
		sealReward: sealReward,
		sealTarget: sealTarget,
		deposits:   make([]timeCapsuleDeposit, 0, LuckyTimeCapsuleMaxDeposits),
	}
	m.activeSessions[p.ID] = sess
	m.mu.Unlock()

	log.Printf("[TimeCapsule] player=%s 觸發時間膠囊！封存 %s（×%.1f）獎勵 %d，時限 %ds",
		p.ID, sealTarget, sealMult, sealReward, int(LuckyTimeCapsuleDuration.Seconds()))

	// 個人訊息
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyTimeCapsule,
		Payload: ws.LuckyTimeCapsulePayload{
			Event:       "capsule_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyTimeCapsuleDuration.Seconds()),
			SealTarget:  sealTarget,
			SealMult:    LuckyTimeCapsuleSealMult,
			SealReward:  sealReward,
			MaxDeposits: LuckyTimeCapsuleMaxDeposits,
			DepositMult: LuckyTimeCapsuleDepositMult,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTimeCapsule,
		Payload: ws.LuckyTimeCapsulePayload{
			Event:      "capsule_broadcast",
			PlayerName: p.DisplayName,
			SealTarget: sealTarget,
			SealMult:   LuckyTimeCapsuleSealMult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyTimeCapsule, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⏳ %s 觸發時間膠囊！封存 %s ×%.1f，15秒後開啟！",
			p.DisplayName, sealTarget, LuckyTimeCapsuleSealMult),
		"color": "#4A90D9",
	})

	// 啟動超時 goroutine
	go g.runTimeCapsuleTimeout(sess)
}

// notifyTimeCapsuleKill 時間膠囊模式中擊破目標（由 handleKill 呼叫）
func (g *Game) notifyTimeCapsuleKill(p *player.Player, t *target.Target) {
	m := g.LuckyTimeCapsule

	m.mu.Lock()
	sess, ok := m.activeSessions[p.ID]
	if !ok || time.Now().After(sess.expiresAt) {
		m.mu.Unlock()
		return
	}
	m.mu.Unlock()

	sess.depositsMu.Lock()
	if len(sess.deposits) >= LuckyTimeCapsuleMaxDeposits {
		sess.depositsMu.Unlock()
		return
	}

	avgBet := g.getAvgBetCost()
	depositReward := int(float64(avgBet) * t.Multiplier * LuckyTimeCapsuleDepositMult)
	if depositReward < 1 {
		depositReward = 1
	}

	sess.deposits = append(sess.deposits, timeCapsuleDeposit{
		targetName: t.Def.Name,
		reward:     depositReward,
	})
	depositCount := len(sess.deposits)
	sess.depositsMu.Unlock()

	log.Printf("[TimeCapsule] player=%s 追加存入 %d/%d！目標=%s 獎勵=%d（×%.1f）",
		p.ID, depositCount, LuckyTimeCapsuleMaxDeposits, t.Def.Name, depositReward, LuckyTimeCapsuleDepositMult)

	// 個人通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyTimeCapsule,
		Payload: ws.LuckyTimeCapsulePayload{
			Event:        "capsule_deposit",
			PlayerName:   p.DisplayName,
			DepositCount: depositCount,
			MaxDeposits:  LuckyTimeCapsuleMaxDeposits,
			DepositMult:  LuckyTimeCapsuleDepositMult,
			Reward:       depositReward,
			TargetName:   t.Def.Name,
		},
	})
}

// runTimeCapsuleTimeout 時間膠囊超時處理（15 秒後開啟膠囊）
func (g *Game) runTimeCapsuleTimeout(sess *timeCapsuleSession) {
	timer := time.NewTimer(LuckyTimeCapsuleDuration)
	defer timer.Stop()

	select {
	case <-timer.C:
		m := g.LuckyTimeCapsule

		m.mu.Lock()
		if _, ok := m.activeSessions[sess.playerID]; !ok {
			m.mu.Unlock()
			return
		}
		delete(m.activeSessions, sess.playerID)
		m.mu.Unlock()

		// 計算總獎勵
		sess.depositsMu.Lock()
		deposits := make([]timeCapsuleDeposit, len(sess.deposits))
		copy(deposits, sess.deposits)
		sess.depositsMu.Unlock()

		totalReward := sess.sealReward
		for _, d := range deposits {
			totalReward += d.reward
		}

		// 找到玩家並發放獎勵
		g.mu.RLock()
		p, ok := g.Players[sess.playerID]
		g.mu.RUnlock()
		if ok {
			p.AddReward(totalReward)
		}

		log.Printf("[TimeCapsule] player=%s 膠囊開啟！封存獎勵=%d，追加=%d次，總獎勵=%d",
			sess.playerID, sess.sealReward, len(deposits), totalReward)

		// 個人通知
		_ = g.Hub.Send(sess.playerID, &ws.Message{
			Type: ws.MsgLuckyTimeCapsule,
			Payload: ws.LuckyTimeCapsulePayload{
				Event:        "capsule_open",
				PlayerID:     sess.playerID,
				PlayerName:   sess.playerName,
				SealTarget:   sess.sealTarget,
				SealMult:     LuckyTimeCapsuleSealMult,
				SealReward:   sess.sealReward,
				DepositCount: len(deposits),
				MaxDeposits:  LuckyTimeCapsuleMaxDeposits,
				TotalReward:  totalReward,
			},
		})

		// 全服廣播膠囊開啟
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyTimeCapsule,
			Payload: ws.LuckyTimeCapsulePayload{
				Event:        "capsule_open_broadcast",
				PlayerName:   sess.playerName,
				SealTarget:   sess.sealTarget,
				DepositCount: len(deposits),
				TotalReward:  totalReward,
			},
		})

		// 全服公告（總獎勵豐厚時）
		if totalReward >= 50 {
			g.Announce.Create(announce.EventLuckyTimeCapsule, sess.playerName, 0, map[string]string{
				"message": fmt.Sprintf("⏳ %s 開啟時間膠囊！存入 %d 次，總獎勵 +%d！",
					sess.playerName, len(deposits), totalReward),
				"color": "#FFD700",
			})
		}

	case <-g.stopCh:
		return
	}
}

// findHighestMultTarget 找場上最高倍率目標（供時間膠囊封存使用）
// 回傳：目標名稱、倍率、計算後的獎勵
func (g *Game) findHighestMultTarget(avgBet int, sealMult float64) (string, float64, int) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var bestName string
	var bestMult float64
	for _, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		if t.Multiplier > bestMult {
			bestMult = t.Multiplier
			bestName = t.Def.Name
		}
	}

	if bestName == "" {
		bestName = "神秘目標"
		bestMult = 5.0
	}

	reward := int(float64(avgBet) * bestMult * sealMult * 3)
	if reward < 1 {
		reward = 1
	}
	return bestName, bestMult, reward
}
