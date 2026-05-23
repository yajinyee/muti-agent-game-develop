// lucky_crystal_ball_fish_handler.go — 幸運水晶球魚系統（DAY-246）
// 業界原創「預測未來+命中率提升」機制
//
// 設計：擊破 T204 後，Server「預測」場上 3 個目標為「水晶預言目標」（持續 8 秒）：
//   - 玩家射擊水晶預言目標時，命中率提升至 100%（必中）
//   - 每次必中擊破獲得 ×2.5 倍率加成（個人獎勵）
//   - 8 秒後「水晶爆炸」：所有未擊破的水晶預言目標自動爆炸（×1.8 倍率，個人獎勵）
//   - 個人冷卻 20 秒；全服冷卻 30 秒
//
// 設計差異：
//   - 與預言魚（T201，指定 1 個目標×3.5）不同，水晶球魚是「3 個目標全部必中」
//     讓玩家有「趕快把這 3 條魚全打掉」的緊迫感
//   - 「必中」讓玩家感受到「這 8 秒我不會浪費任何一槍」的掌控感
//   - 「×2.5 倍率」比普通擊破高，讓玩家有「要集中打這 3 條」的動機
//   - 「水晶爆炸」確保即使沒打完也有獎勵，降低挫敗感
//   - 全服廣播讓其他玩家看到「有人觸發了水晶預言」，製造羨慕感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyCrystalBallPersonalCD  = 20 * time.Second // 個人冷卻
	LuckyCrystalBallGlobalCD    = 30 * time.Second // 全服冷卻
	LuckyCrystalBallDuration    = 8 * time.Second  // 水晶預言持續時間
	LuckyCrystalBallTargetCount = 3                // 預言目標數量
	LuckyCrystalBallHitMult     = 2.5              // 必中擊破倍率
	LuckyCrystalBallBlastMult   = 1.8              // 水晶爆炸倍率
)

// crystalBallEntry 水晶預言目標
type crystalBallEntry struct {
	instanceID string
	defID      string
	x          float64
	y          float64
	expiresAt  time.Time
}

// crystalBallSession 水晶預言 session
type crystalBallSession struct {
	playerID  string
	expiresAt time.Time
	// 水晶預言目標（instanceID → entry）
	targets map[string]*crystalBallEntry
	mu      sync.Mutex
	// 統計
	hitCount   int
	blastCount int
	totalReward int
}

// luckyCrystalBallFishManager 幸運水晶球魚管理器
type luckyCrystalBallFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前水晶預言 sessions（playerID → session）
	activeSessions map[string]*crystalBallSession
}

func newLuckyCrystalBallFishManager() *luckyCrystalBallFishManager {
	return &luckyCrystalBallFishManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*crystalBallSession),
	}
}

// isLuckyCrystalBallFish 判斷是否為幸運水晶球魚
func isLuckyCrystalBallFish(defID string) bool {
	return defID == "T204"
}

// isCrystalBallTarget 判斷目標是否為水晶預言目標（供 handleKill 使用）
// 回傳 (isCrystal bool, playerID string)
func (m *luckyCrystalBallFishManager) isCrystalBallTarget(instanceID string) (bool, string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for playerID, sess := range m.activeSessions {
		if time.Now().After(sess.expiresAt) {
			continue
		}
		sess.mu.Lock()
		_, ok := sess.targets[instanceID]
		sess.mu.Unlock()
		if ok {
			return true, playerID
		}
	}
	return false, ""
}

// removeCrystalBallTarget 移除水晶預言目標
func (m *luckyCrystalBallFishManager) removeCrystalBallTarget(playerID, instanceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sess, ok := m.activeSessions[playerID]; ok {
		sess.mu.Lock()
		delete(sess.targets, instanceID)
		sess.mu.Unlock()
	}
}

// tryLuckyCrystalBallFish 擊破 T204 後觸發水晶預言
func (g *Game) tryLuckyCrystalBallFish(p *player.Player) {
	m := g.LuckyCrystalBallFish
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
	m.personalCooldowns[p.ID] = now.Add(LuckyCrystalBallPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyCrystalBallGlobalCD)
	m.mu.Unlock()

	// 隨機選取場上目標
	g.mu.RLock()
	var candidates []*target.Target
	for _, t := range g.Targets {
		if t.DefID != "T204" && t.IsAlive {
			candidates = append(candidates, t)
		}
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		log.Printf("[CrystalBall] player=%s 場上無可用目標", p.ID)
		return
	}

	// 隨機打亂，取前 N 個
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	count := LuckyCrystalBallTargetCount
	if count > len(candidates) {
		count = len(candidates)
	}
	selected := candidates[:count]

	// 建立 session
	sess := &crystalBallSession{
		playerID:  p.ID,
		expiresAt: now.Add(LuckyCrystalBallDuration),
		targets:   make(map[string]*crystalBallEntry),
	}
	for _, t := range selected {
		sess.targets[t.InstanceID] = &crystalBallEntry{
			instanceID: t.InstanceID,
			defID:      t.DefID,
			x:          t.X,
			y:          t.Y,
			expiresAt:  now.Add(LuckyCrystalBallDuration),
		}
	}

	m.mu.Lock()
	m.activeSessions[p.ID] = sess
	m.mu.Unlock()

	// 收集目標 ID 列表
	targetIDs := make([]string, 0, count)
	for _, t := range selected {
		targetIDs = append(targetIDs, t.InstanceID)
	}

	log.Printf("[CrystalBall] player=%s 水晶預言啟動，%d 個目標", p.ID, count)

	// 個人訊息：水晶預言啟動
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyCrystalBallFish,
		Payload: ws.LuckyCrystalBallFishPayload{
			Event:       "crystal_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			TargetIDs:   targetIDs,
			DurationSec: int(LuckyCrystalBallDuration.Seconds()),
			HitMult:     LuckyCrystalBallHitMult,
			BlastMult:   LuckyCrystalBallBlastMult,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyCrystalBallFish,
		Payload: ws.LuckyCrystalBallFishPayload{
			Event:      "crystal_broadcast",
			PlayerName: p.DisplayName,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyCrystalBallFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🔮 %s 觸發水晶預言！%d 個目標必中，×%.1f 倍率！", p.DisplayName, count, LuckyCrystalBallHitMult),
		"color":   "#1ABC9C",
	})

	// 啟動計時 goroutine
	go g.runLuckyCrystalBall(p, sess)
}

// notifyCrystalBallKill 玩家擊破水晶預言目標時觸發（由 handleKill 呼叫）
func (g *Game) notifyCrystalBallKill(p *player.Player, instanceID string, baseReward int) float64 {
	m := g.LuckyCrystalBallFish

	// 計算個人 betCost
	betDef := data.GetBetDef(p.BetLevel)
	avgBet := betDef.BetCost
	if avgBet < 1 {
		avgBet = 1
	}

	// 移除目標
	m.removeCrystalBallTarget(p.ID, instanceID)

	// 計算獎勵
	reward := int(float64(avgBet) * LuckyCrystalBallHitMult)
	p.AddCoins(reward)

	// 更新統計
	m.mu.Lock()
	if sess, ok := m.activeSessions[p.ID]; ok {
		sess.mu.Lock()
		sess.hitCount++
		sess.totalReward += reward
		sess.mu.Unlock()
	}
	m.mu.Unlock()

	log.Printf("[CrystalBall] player=%s 必中擊破 instanceID=%s reward=%d", p.ID, instanceID, reward)

	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyCrystalBallFish,
		Payload: ws.LuckyCrystalBallFishPayload{
			Event:    "crystal_hit",
			TargetID: instanceID,
			Reward:   reward,
			HitMult:  LuckyCrystalBallHitMult,
		},
	})

	return LuckyCrystalBallHitMult
}

// runLuckyCrystalBall 水晶預言計時 goroutine
func (g *Game) runLuckyCrystalBall(p *player.Player, sess *crystalBallSession) {
	timer := time.NewTimer(LuckyCrystalBallDuration)
	defer timer.Stop()

	select {
	case <-timer.C:
		g.doCrystalBallBlast(p, sess)
	case <-g.stopCh:
		return
	}
}

// doCrystalBallBlast 水晶爆炸（計時結束時觸發）
func (g *Game) doCrystalBallBlast(p *player.Player, sess *crystalBallSession) {
	m := g.LuckyCrystalBallFish

	// 計算個人 betCost
	betDef := data.GetBetDef(p.BetLevel)
	avgBet := betDef.BetCost
	if avgBet < 1 {
		avgBet = 1
	}

	// 收集剩餘存活的水晶預言目標
	sess.mu.Lock()
	var remaining []*crystalBallEntry
	for _, entry := range sess.targets {
		if time.Now().Before(entry.expiresAt) {
			remaining = append(remaining, entry)
		}
	}
	sess.targets = make(map[string]*crystalBallEntry)
	hitCount := sess.hitCount
	totalReward := sess.totalReward
	sess.mu.Unlock()

	// 移除 session
	m.mu.Lock()
	delete(m.activeSessions, p.ID)
	m.mu.Unlock()

	blastCount := 0
	for _, entry := range remaining {
		// 確認目標仍在場上
		g.mu.RLock()
		t, exists := g.Targets[entry.instanceID]
		g.mu.RUnlock()
		if !exists || !t.IsAlive {
			continue
		}

		// 水晶爆炸：移除目標
		g.mu.Lock()
		delete(g.Targets, entry.instanceID)
		g.mu.Unlock()

		reward := int(float64(avgBet) * LuckyCrystalBallBlastMult)
		p.AddCoins(reward)
		totalReward += reward
		blastCount++

		log.Printf("[CrystalBall] player=%s 水晶爆炸 instanceID=%s reward=%d", p.ID, entry.instanceID, reward)

		_ = g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyCrystalBallFish,
			Payload: ws.LuckyCrystalBallFishPayload{
				Event:     "crystal_blast",
				TargetID:  entry.instanceID,
				Reward:    reward,
				BlastMult: LuckyCrystalBallBlastMult,
			},
		})
	}

	// 結束通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyCrystalBallFish,
		Payload: ws.LuckyCrystalBallFishPayload{
			Event:       "crystal_end",
			HitCount:    hitCount,
			BlastCount:  blastCount,
			TotalReward: totalReward,
		},
	})

	if hitCount+blastCount >= 2 {
		g.Announce.Create(announce.EventLuckyCrystalBallFish, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("🔮 %s 水晶預言結束！命中 %d 個，爆炸 %d 個，獲得 %d 籌碼！",
				p.DisplayName, hitCount, blastCount, totalReward),
			"color": "#16A085",
		})
	}
}
