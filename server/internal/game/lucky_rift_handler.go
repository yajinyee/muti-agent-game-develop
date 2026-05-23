// lucky_rift_handler.go — 幸運時空裂縫魚系統（DAY-255）
// 業界原創「時空裂縫+傳送吸入+裂縫崩塌」機制
//
// 設計：擊破 T213 後，場景中央出現「時空裂縫」（持續 18 秒）：
//   - 每 3 秒「裂縫吸入」：吸入距離裂縫最近的目標，傳送到隨機位置（×1.6 倍率，全服共享）
//   - 最多吸入 5 個目標（達到上限後裂縫提前崩塌）
//   - 18 秒後「裂縫崩塌」：場上所有目標 HP -50%，全服 AOE 獎勵（×2.5 倍率，全服共享）
//   - 個人冷卻 22 秒；全服冷卻 35 秒
//
// 設計差異：
//   - 與黑洞爆炸魚（T207，吸收消滅+能量爆炸）不同，時空裂縫是「傳送而非消滅」，
//     目標被傳送到隨機位置後仍然存活，讓玩家有「魚突然出現在意想不到的地方」的驚喜感
//   - 「傳送後仍存活」讓玩家有「要趕快找到被傳送的魚繼續打」的追逐感
//   - 「裂縫崩塌 HP -50%」讓玩家在崩塌後有「全場魚都快死了，趕快打」的緊迫感
//   - 「×2.5 全服 AOE 獎勵」讓所有玩家都受益，製造「全服一起等崩塌」的社交感
//   - 全服廣播裂縫位置讓所有玩家都知道「裂縫在哪裡」，製造「全服一起看裂縫」的緊張感
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
	LuckyRiftPersonalCD    = 22 * time.Second // 個人冷卻
	LuckyRiftGlobalCD      = 35 * time.Second // 全服冷卻
	LuckyRiftDuration      = 18 * time.Second // 裂縫持續時間
	LuckyRiftSuckInterval  = 3 * time.Second  // 吸入間隔
	LuckyRiftMaxSuck       = 5                // 最多吸入目標數
	LuckyRiftSuckMult      = 1.6              // 吸入倍率（全服共享）
	LuckyRiftCollapseHPDrain = 0.50           // 崩塌 HP 扣除比例
	LuckyRiftCollapseMult  = 2.5              // 崩塌全服 AOE 倍率

	// 裂縫位置（場景中央）
	LuckyRiftCenterX = 500.0
	LuckyRiftCenterY = 300.0
)

// riftSession 時空裂縫會話
type riftSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	expiresAt         time.Time
	suckCount         int
	mu                sync.Mutex
}

// luckyRiftManager 幸運時空裂縫魚管理器
type luckyRiftManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前活躍的裂縫會話（nil = 無）
	activeSession *riftSession
}

func newLuckyRiftManager() *luckyRiftManager {
	return &luckyRiftManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyRiftFish 判斷是否為幸運時空裂縫魚
func isLuckyRiftFish(defID string) bool {
	return defID == "T213"
}

// isRiftActive 判斷時空裂縫是否正在進行
func (m *luckyRiftManager) isRiftActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil {
		return false
	}
	if time.Now().After(m.activeSession.expiresAt) {
		m.activeSession = nil
		return false
	}
	return true
}

// tryLuckyRiftFish 擊破 T213 後觸發時空裂縫
func (g *Game) tryLuckyRiftFish(p *player.Player) {
	m := g.LuckyRift

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
	// 已有活躍裂縫
	if m.activeSession != nil && now.Before(m.activeSession.expiresAt) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyRiftPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyRiftGlobalCD)

	expiresAt := now.Add(LuckyRiftDuration)
	sess := &riftSession{
		triggerPlayerID:   p.ID,
		triggerPlayerName: p.DisplayName,
		expiresAt:         expiresAt,
		suckCount:         0,
	}
	m.activeSession = sess
	m.mu.Unlock()

	log.Printf("[Rift] player=%s 觸發時空裂縫！持續 %ds，裂縫位置 (%.0f, %.0f)",
		p.ID, int(LuckyRiftDuration.Seconds()), LuckyRiftCenterX, LuckyRiftCenterY)

	// 個人訊息：裂縫啟動
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyRift,
		Payload: ws.LuckyRiftPayload{
			Event:       "rift_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyRiftDuration.Seconds()),
			SuckMult:    LuckyRiftSuckMult,
			CollapseMult: LuckyRiftCollapseMult,
			RiftX:       LuckyRiftCenterX,
			RiftY:       LuckyRiftCenterY,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyRift,
		Payload: ws.LuckyRiftPayload{
			Event:      "rift_broadcast",
			PlayerName: p.DisplayName,
			SuckMult:   LuckyRiftSuckMult,
			RiftX:      LuckyRiftCenterX,
			RiftY:      LuckyRiftCenterY,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyRift, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌀 %s 開啟時空裂縫！每 3 秒吸入最近目標傳送！×%.1f 倍率！全服共享！",
			p.DisplayName, LuckyRiftSuckMult),
		"color": "#6A0DAD",
	})

	// 啟動裂縫主循環 goroutine
	go g.runRiftSession(sess, p)
}

// runRiftSession 時空裂縫主循環 goroutine
func (g *Game) runRiftSession(sess *riftSession, triggerPlayer *player.Player) {
	ticker := time.NewTicker(LuckyRiftSuckInterval)
	defer ticker.Stop()

	maxTicks := int(LuckyRiftDuration / LuckyRiftSuckInterval)

	for tick := 1; ; tick++ {
		select {
		case <-ticker.C:
			sess.mu.Lock()
			currentSuck := sess.suckCount
			sess.mu.Unlock()

			if tick > maxTicks || currentSuck >= LuckyRiftMaxSuck {
				// 觸發裂縫崩塌
				g.doRiftCollapse(sess, triggerPlayer)
				return
			}
			g.doRiftSuck(sess, tick)

		case <-g.stopCh:
			return
		}
	}
}

// doRiftSuck 執行一次裂縫吸入
func (g *Game) doRiftSuck(sess *riftSession, suckNum int) {
	// 找距離裂縫最近的存活目標
	g.mu.RLock()
	var closestTarget *target.Target
	closestDist := math.MaxFloat64

	for _, t := range g.Targets {
		if !t.IsAlive {
			continue
		}
		dx := t.X - LuckyRiftCenterX
		dy := t.Y - LuckyRiftCenterY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < closestDist {
			closestDist = dist
			closestTarget = t
		}
	}
	g.mu.RUnlock()

	if closestTarget == nil {
		return
	}

	// 傳送目標到隨機位置（不消滅，只移動）
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	newX := 100.0 + rng.Float64()*800.0 // 100~900 範圍
	newY := 50.0 + rng.Float64()*500.0  // 50~550 範圍

	g.mu.Lock()
	existing, ok := g.Targets[closestTarget.InstanceID]
	if !ok || !existing.IsAlive {
		g.mu.Unlock()
		return
	}
	oldX := existing.X
	oldY := existing.Y
	existing.X = newX
	existing.Y = newY
	g.mu.Unlock()

	// 更新吸入計數
	sess.mu.Lock()
	sess.suckCount++
	currentSuck := sess.suckCount
	sess.mu.Unlock()

	// 計算全服共享獎勵
	avgBet := g.getAvgBetCost()
	reward := int(float64(avgBet) * closestTarget.Multiplier * LuckyRiftSuckMult)
	if reward < 1 {
		reward = 1
	}
	g.distributeRewardToAll(reward)

	log.Printf("[Rift] 第 %d 次吸入！目標 %s (%.0f,%.0f) → (%.0f,%.0f)，全服獎勵 %d",
		suckNum, closestTarget.InstanceID, oldX, oldY, newX, newY, reward)

	// 全服廣播吸入事件
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyRift,
		Payload: ws.LuckyRiftPayload{
			Event:       "rift_suck",
			PlayerName:  sess.triggerPlayerName,
			SuckNum:     suckNum,
			SuckCount:   currentSuck,
			MaxSuck:     LuckyRiftMaxSuck,
			TargetName:  closestTarget.Def.Name,
			OldX:        oldX,
			OldY:        oldY,
			NewX:        newX,
			NewY:        newY,
			Mult:        LuckyRiftSuckMult,
			TotalReward: reward,
		},
	})
}

// doRiftCollapse 裂縫崩塌（持續時間結束或達到最大吸入數）
func (g *Game) doRiftCollapse(sess *riftSession, triggerPlayer *player.Player) {
	// 清除活躍會話
	m := g.LuckyRift
	m.mu.Lock()
	m.activeSession = nil
	m.mu.Unlock()

	// 場上所有目標 HP -50%
	g.mu.Lock()
	drainCount := 0
	totalCollapseReward := 0
	avgBet := 0
	for _, t := range g.Targets {
		if t.IsAlive {
			newHP := int(float64(t.HP) * (1.0 - LuckyRiftCollapseHPDrain))
			if newHP < 1 {
				newHP = 1
			}
			t.HP = newHP
			drainCount++
		}
	}
	g.mu.Unlock()

	// 計算全服 AOE 獎勵（基於場上目標數）
	avgBet = g.getAvgBetCost()
	if drainCount > 0 {
		totalCollapseReward = int(float64(avgBet) * LuckyRiftCollapseMult * float64(drainCount))
		if totalCollapseReward < 1 {
			totalCollapseReward = 1
		}
		g.distributeRewardToAll(totalCollapseReward)
	}

	sess.mu.Lock()
	finalSuckCount := sess.suckCount
	sess.mu.Unlock()

	log.Printf("[Rift] 裂縫崩塌！%d 個目標 HP -50%%，全服 AOE 獎勵 %d（吸入了 %d 個目標）",
		drainCount, totalCollapseReward, finalSuckCount)

	// 全服廣播崩塌
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyRift,
		Payload: ws.LuckyRiftPayload{
			Event:        "rift_collapse",
			PlayerName:   sess.triggerPlayerName,
			DrainCount:   drainCount,
			CollapseMult: LuckyRiftCollapseMult,
			TotalReward:  totalCollapseReward,
			SuckCount:    finalSuckCount,
		},
	})

	// 全服公告（drainCount >= 3 才公告）
	if drainCount >= 3 {
		g.Announce.Create(announce.EventLuckyRift, sess.triggerPlayerName, 0, map[string]string{
			"message": fmt.Sprintf("🌀 時空裂縫崩塌！%d 個目標 HP -50%%！全服 AOE ×%.1f 獎勵！",
				drainCount, LuckyRiftCollapseMult),
			"color": "#4B0082",
		})
	}
}
