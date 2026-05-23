// lucky_lightning_storm_handler.go — 幸運閃電風暴魚系統（DAY-258）
// 業界原創「閃電風暴+連鎖跳躍+超級閃電+全場電擊」機制
//
// 設計：擊破 T216 後，觸發「閃電風暴」（持續 12 秒）：
//   - 每 1.5 秒「閃電跳躍」：從場上隨機目標出發，連鎖跳躍到最近的 3 個目標
//     （每跳 ×1.3 倍率，全服共享）
//   - 若連鎖達到 5 跳以上 → 「超級閃電」：×3.0 倍率（全服大獎）
//   - 12 秒後「閃電爆炸」：場上所有目標 HP -40%（全服共享）
//   - 個人冷卻 20 秒；全服冷卻 32 秒
//
// 設計差異：
//   - 與閃電鰻（已有，單條連鎖）不同，閃電風暴是「多輪連鎖跳躍」，
//     讓玩家看到「閃電在魚群中不斷跳躍」的視覺爽感
//   - 「連鎖達到 5 跳觸發超級閃電」讓玩家有「要趁風暴期間多打魚讓閃電有更多目標跳」的策略感
//   - 「×3.0 超級閃電」是目前連鎖類最高倍率，製造「哇，超級閃電！」的驚嘆感
//   - 「閃電爆炸 HP -40%」讓玩家在風暴結束後有「全場魚都快死了，趕快打」的緊迫感
//   - 全服廣播每次跳躍讓所有玩家看到「閃電在哪裡跳」，製造「全服一起看閃電」的社交感
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyLightningStormPersonalCD  = 20 * time.Second  // 個人冷卻
	LuckyLightningStormGlobalCD    = 32 * time.Second  // 全服冷卻
	LuckyLightningStormDuration    = 12 * time.Second  // 風暴時限
	LuckyLightningStormJumpTick    = 1500 * time.Millisecond // 跳躍間隔
	LuckyLightningStormJumpMult    = 1.3               // 每跳倍率
	LuckyLightningStormSuperThresh = 5                 // 超級閃電觸發跳數
	LuckyLightningStormSuperMult   = 3.0               // 超級閃電倍率
	LuckyLightningStormBlastHPCut  = 0.4               // 閃電爆炸 HP 削減比例
	LuckyLightningStormMaxJumps    = 3                 // 每輪最多跳躍目標數
)

// lightningStormSession 閃電風暴會話
type lightningStormSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	expiresAt         time.Time
	totalJumps        int // 累計跳躍次數（用於超級閃電判斷）
	mu                sync.Mutex
}

// luckyLightningStormManager 幸運閃電風暴魚管理器
type luckyLightningStormManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前活躍的閃電風暴會話（nil = 無）
	activeSession *lightningStormSession
}

func newLuckyLightningStormManager() *luckyLightningStormManager {
	return &luckyLightningStormManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyLightningStormFish 判斷是否為幸運閃電風暴魚
func isLuckyLightningStormFish(defID string) bool {
	return defID == "T216"
}

// tryLuckyLightningStormFish 擊破 T216 後觸發閃電風暴
func (g *Game) tryLuckyLightningStormFish(p *player.Player) {
	m := g.LuckyLightningStorm

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
	// 已有活躍風暴
	if m.activeSession != nil && now.Before(m.activeSession.expiresAt) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyLightningStormPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyLightningStormGlobalCD)

	expiresAt := now.Add(LuckyLightningStormDuration)
	sess := &lightningStormSession{
		triggerPlayerID:   p.ID,
		triggerPlayerName: p.DisplayName,
		expiresAt:         expiresAt,
		totalJumps:        0,
	}
	m.activeSession = sess
	m.mu.Unlock()

	log.Printf("[LightningStorm] player=%s 觸發閃電風暴！時限 %ds",
		p.ID, int(LuckyLightningStormDuration.Seconds()))

	// 個人訊息
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyLightningStorm,
		Payload: ws.LuckyLightningStormPayload{
			Event:        "storm_start",
			PlayerID:     p.ID,
			PlayerName:   p.DisplayName,
			DurationSec:  int(LuckyLightningStormDuration.Seconds()),
			JumpMult:     LuckyLightningStormJumpMult,
			SuperMult:    LuckyLightningStormSuperMult,
			SuperThresh:  LuckyLightningStormSuperThresh,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyLightningStorm,
		Payload: ws.LuckyLightningStormPayload{
			Event:       "storm_broadcast",
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyLightningStormDuration.Seconds()),
			JumpMult:    LuckyLightningStormJumpMult,
			SuperMult:   LuckyLightningStormSuperMult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyLightningStorm, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⚡ %s 觸發閃電風暴！每 1.5 秒連鎖跳躍 ×%.1f！達到 %d 跳→超級閃電 ×%.1f！",
			p.DisplayName, LuckyLightningStormJumpMult, LuckyLightningStormSuperThresh, LuckyLightningStormSuperMult),
		"color": "#FFD700",
	})

	// 啟動風暴主循環
	go g.runLightningStorm(sess)
}

// runLightningStorm 閃電風暴主循環
func (g *Game) runLightningStorm(sess *lightningStormSession) {
	jumpTicker := time.NewTicker(LuckyLightningStormJumpTick)
	endTimer := time.NewTimer(LuckyLightningStormDuration)
	defer jumpTicker.Stop()
	defer endTimer.Stop()

	jumpRound := 0

	for {
		select {
		case <-jumpTicker.C:
			jumpRound++
			g.doLightningJump(sess, jumpRound)

		case <-endTimer.C:
			// 風暴結束，觸發閃電爆炸
			m := g.LuckyLightningStorm
			m.mu.Lock()
			if m.activeSession != sess {
				m.mu.Unlock()
				return
			}
			m.activeSession = nil
			m.mu.Unlock()

			g.doLightningBlast(sess)
			return

		case <-g.stopCh:
			return
		}
	}
}

// doLightningJump 執行一輪閃電跳躍
func (g *Game) doLightningJump(sess *lightningStormSession, round int) {
	g.mu.Lock()
	// 收集所有存活目標
	var alive []*target.Target
	for _, t := range g.Targets {
		if t.IsAlive && !isLuckyLightningStormFish(t.DefID) {
			alive = append(alive, t)
		}
	}
	g.mu.Unlock()

	if len(alive) == 0 {
		return
	}

	// 隨機選起始目標
	startIdx := rand.Intn(len(alive))
	startTarget := alive[startIdx]

	// 從起始目標連鎖跳躍到最近的 N 個目標
	jumped := []*target.Target{startTarget}
	remaining := make([]*target.Target, 0, len(alive)-1)
	for i, t := range alive {
		if i != startIdx {
			remaining = append(remaining, t)
		}
	}

	maxJumps := LuckyLightningStormMaxJumps
	if maxJumps > len(remaining) {
		maxJumps = len(remaining)
	}

	// 貪心選最近的目標
	current := startTarget
	for i := 0; i < maxJumps; i++ {
		if len(remaining) == 0 {
			break
		}
		// 找最近的
		minDist := math.MaxFloat64
		minIdx := 0
		for j, t := range remaining {
			dx := t.X - current.X
			dy := t.Y - current.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < minDist {
				minDist = dist
				minIdx = j
			}
		}
		next := remaining[minIdx]
		jumped = append(jumped, next)
		remaining = append(remaining[:minIdx], remaining[minIdx+1:]...)
		current = next
	}

	// 計算本輪跳躍數（不含起始目標）
	thisJumps := len(jumped) - 1
	if thisJumps <= 0 {
		return
	}

	// 累計跳躍數
	sess.mu.Lock()
	sess.totalJumps += thisJumps
	totalJumps := sess.totalJumps
	sess.mu.Unlock()

	// 計算獎勵（全服共享）
	avgBet := g.getAvgBetCost()
	totalReward := 0
	var jumpedNames []string
	for _, t := range jumped[1:] { // 跳過起始目標（不給獎勵，只是起點）
		reward := int(float64(avgBet) * t.Multiplier * LuckyLightningStormJumpMult)
		if reward < 1 {
			reward = 1
		}
		totalReward += reward
		jumpedNames = append(jumpedNames, t.Def.Name)
	}
	if totalReward > 0 {
		g.distributeRewardToAll(totalReward)
	}

	log.Printf("[LightningStorm] 第 %d 輪跳躍！跳 %d 個目標，全服獎勵 %d（累計 %d 跳）",
		round, thisJumps, totalReward, totalJumps)

	// 廣播跳躍
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyLightningStorm,
		Payload: ws.LuckyLightningStormPayload{
			Event:       "storm_jump",
			PlayerName:  sess.triggerPlayerName,
			Round:       round,
			JumpCount:   thisJumps,
			TotalJumps:  totalJumps,
			JumpMult:    LuckyLightningStormJumpMult,
			TotalReward: totalReward,
			StartX:      startTarget.X,
			StartY:      startTarget.Y,
		},
	})

	// 超級閃電判斷（累計跳躍達到閾值）
	if totalJumps >= LuckyLightningStormSuperThresh && (totalJumps-thisJumps) < LuckyLightningStormSuperThresh {
		// 剛剛達到閾值，觸發超級閃電
		go g.doSuperLightning(sess)
	}
}

// doSuperLightning 超級閃電（累計跳躍達到閾值）
func (g *Game) doSuperLightning(sess *lightningStormSession) {
	// 選場上最高倍率目標
	g.mu.Lock()
	var bestTarget *target.Target
	for _, t := range g.Targets {
		if t.IsAlive && !isLuckyLightningStormFish(t.DefID) {
			if bestTarget == nil || t.Multiplier > bestTarget.Multiplier {
				bestTarget = t
			}
		}
	}
	g.mu.Unlock()

	avgBet := g.getAvgBetCost()
	superReward := 0
	targetName := "目標"
	if bestTarget != nil {
		superReward = int(float64(avgBet) * bestTarget.Multiplier * LuckyLightningStormSuperMult)
		if superReward < 1 {
			superReward = 1
		}
		targetName = bestTarget.Def.Name
		g.distributeRewardToAll(superReward)
	} else {
		// 無目標時給固定獎勵
		superReward = int(float64(avgBet) * LuckyLightningStormSuperMult * 5)
		if superReward < 1 {
			superReward = 1
		}
		g.distributeRewardToAll(superReward)
	}

	log.Printf("[LightningStorm] 超級閃電！目標=%s，全服獎勵 %d（×%.1f）",
		targetName, superReward, LuckyLightningStormSuperMult)

	// 全服廣播超級閃電
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyLightningStorm,
		Payload: ws.LuckyLightningStormPayload{
			Event:       "super_lightning",
			PlayerName:  sess.triggerPlayerName,
			SuperMult:   LuckyLightningStormSuperMult,
			TotalReward: superReward,
			TotalJumps:  LuckyLightningStormSuperThresh,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyLightningStorm, sess.triggerPlayerName, 0, map[string]string{
		"message": fmt.Sprintf("⚡ 超級閃電！累計 %d 跳！全服獎勵 +%d！×%.1f 大獎！",
			LuckyLightningStormSuperThresh, superReward, LuckyLightningStormSuperMult),
		"color": "#FFFFFF",
	})
}

// doLightningBlast 閃電爆炸（風暴結束）
func (g *Game) doLightningBlast(sess *lightningStormSession) {
	// 場上所有目標 HP -40%
	g.mu.Lock()
	var affected []*target.Target
	for _, t := range g.Targets {
		if t.IsAlive && !isLuckyLightningStormFish(t.DefID) {
			newHP := int(float64(t.HP) * (1.0 - LuckyLightningStormBlastHPCut))
			if newHP < 1 {
				newHP = 1
			}
			t.HP = newHP
			affected = append(affected, t)
		}
	}
	g.mu.Unlock()

	affectedCount := len(affected)

	sess.mu.Lock()
	totalJumps := sess.totalJumps
	sess.mu.Unlock()

	log.Printf("[LightningStorm] 閃電爆炸！影響 %d 個目標 HP-40%%，累計跳躍 %d 次",
		affectedCount, totalJumps)

	// 全服廣播閃電爆炸
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyLightningStorm,
		Payload: ws.LuckyLightningStormPayload{
			Event:         "storm_blast",
			PlayerName:    sess.triggerPlayerName,
			AffectedCount: affectedCount,
			TotalJumps:    totalJumps,
		},
	})

	// 全服公告（影響 ≥ 5 個才公告）
	if affectedCount >= 5 {
		g.Announce.Create(announce.EventLuckyLightningStorm, sess.triggerPlayerName, 0, map[string]string{
			"message": fmt.Sprintf("⚡ 閃電爆炸！%d 個目標 HP-40%%！累計跳躍 %d 次！",
				affectedCount, totalJumps),
			"color": "#87CEEB",
		})
	}
}
