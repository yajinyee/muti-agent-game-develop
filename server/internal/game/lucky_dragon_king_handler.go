// lucky_dragon_king_handler.go — 幸運龍王降臨魚系統（DAY-254）
// 業界原創「龍王降臨+龍息攻擊+龍王護盾+龍王爆發」機制
//
// 設計：擊破 T212 後，「龍王降臨」（持續 15 秒）：
//   - 每 2 秒「龍息攻擊」：隨機選 3 個目標，80% 擊破機率，×1.4 倍率（全服共享）
//   - 龍王降臨期間，觸發玩家獲得「龍王護盾」：下一次被扣費時免費（一次性保護）
//   - 15 秒後「龍王爆發」：場上所有目標 HP -60%，觸發玩家獲得 ×3.0 倍率加成（個人，5 秒）
//   - 個人冷卻 25 秒；全服冷卻 40 秒
//
// 設計差異：
//   - 與隕石雨（T211，隨機轟炸+連擊）不同，龍王降臨是「龍息定向攻擊」，每次選 3 個目標，更有「龍在選擇獵物」的感覺
//   - 「龍王護盾」讓觸發玩家有「我有龍王保護」的安心感，是遊戲中唯一的防禦型機制
//   - 「龍王爆發 HP -60%」讓玩家在爆發後有「全場魚都快死了，趕快打」的緊迫感
//   - 「×3.0 個人倍率加成 5 秒」讓觸發玩家在爆發後有「黃金 5 秒」的爆發感
//   - 全服廣播讓所有玩家看到「龍王降臨了」，製造「全服一起看龍王」的社交感
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
	LuckyDragonKingPersonalCD   = 25 * time.Second // 個人冷卻
	LuckyDragonKingGlobalCD     = 40 * time.Second // 全服冷卻
	LuckyDragonKingDuration     = 15 * time.Second // 龍王降臨持續時間
	LuckyDragonKingBreathInterval = 2 * time.Second // 龍息攻擊間隔
	LuckyDragonKingBreathTargets = 3                // 每次龍息攻擊目標數
	LuckyDragonKingBreathChance  = 0.80             // 龍息擊破機率
	LuckyDragonKingBreathMult    = 1.4              // 龍息倍率（全服共享）
	LuckyDragonKingBurstHPDrain  = 0.60             // 龍王爆發 HP 扣除比例
	LuckyDragonKingBurstMult     = 3.0              // 龍王爆發個人倍率加成
	LuckyDragonKingBurstDuration = 5 * time.Second  // 龍王爆發倍率持續時間
)

// dragonKingSession 龍王降臨會話
type dragonKingSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	expiresAt         time.Time
	breathCount       int
	mu                sync.Mutex
}

// luckyDragonKingManager 幸運龍王降臨魚管理器
type luckyDragonKingManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前活躍的龍王降臨會話（nil = 無）
	activeSession *dragonKingSession

	// 龍王爆發倍率加成（playerID → expiresAt）
	burstBoosts map[string]time.Time

	// 龍王護盾（playerID → 是否有護盾）
	dragonShields map[string]bool
}

func newLuckyDragonKingManager() *luckyDragonKingManager {
	return &luckyDragonKingManager{
		personalCooldowns: make(map[string]time.Time),
		burstBoosts:       make(map[string]time.Time),
		dragonShields:     make(map[string]bool),
	}
}

// isLuckyDragonKingFish 判斷是否為幸運龍王降臨魚
func isLuckyDragonKingFish(defID string) bool {
	return defID == "T212"
}

// isDragonKingActive 判斷龍王降臨是否正在進行
func (m *luckyDragonKingManager) isDragonKingActive() bool {
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

// getLuckyDragonKingBurstMult 取得龍王爆發倍率加成（供 handleKill 使用）
func (m *luckyDragonKingManager) getLuckyDragonKingBurstMult(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if expiresAt, ok := m.burstBoosts[playerID]; ok {
		if time.Now().Before(expiresAt) {
			return LuckyDragonKingBurstMult
		}
		delete(m.burstBoosts, playerID)
	}
	return 1.0
}

// hasDragonShield 判斷玩家是否有龍王護盾
func (m *luckyDragonKingManager) hasDragonShield(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.dragonShields[playerID]
}

// consumeDragonShield 消耗龍王護盾（一次性）
func (m *luckyDragonKingManager) consumeDragonShield(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.dragonShields[playerID] {
		delete(m.dragonShields, playerID)
		return true
	}
	return false
}

// tryLuckyDragonKingFish 擊破 T212 後觸發龍王降臨
func (g *Game) tryLuckyDragonKingFish(p *player.Player) {
	m := g.LuckyDragonKing

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
	// 已有活躍龍王降臨
	if m.activeSession != nil && now.Before(m.activeSession.expiresAt) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyDragonKingPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyDragonKingGlobalCD)

	expiresAt := now.Add(LuckyDragonKingDuration)
	sess := &dragonKingSession{
		triggerPlayerID:   p.ID,
		triggerPlayerName: p.DisplayName,
		expiresAt:         expiresAt,
		breathCount:       0,
	}
	m.activeSession = sess

	// 給觸發玩家龍王護盾
	m.dragonShields[p.ID] = true
	m.mu.Unlock()

	log.Printf("[DragonKing] player=%s 觸發龍王降臨！持續 %ds，龍王護盾已啟動",
		p.ID, int(LuckyDragonKingDuration.Seconds()))

	// 個人訊息：龍王護盾
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyDragonKing,
		Payload: ws.LuckyDragonKingPayload{
			Event:       "dragon_king_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyDragonKingDuration.Seconds()),
			BreathMult:  LuckyDragonKingBreathMult,
			BurstMult:   LuckyDragonKingBurstMult,
			HasShield:   true,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyDragonKing,
		Payload: ws.LuckyDragonKingPayload{
			Event:      "dragon_king_broadcast",
			PlayerName: p.DisplayName,
			BreathMult: LuckyDragonKingBreathMult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyDragonKing, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🐉 %s 召喚龍王降臨！每 2 秒龍息攻擊 3 個目標！×%.1f 倍率！全服共享！",
			p.DisplayName, LuckyDragonKingBreathMult),
		"color": "#8B0000",
	})

	// 啟動龍王降臨 goroutine
	go g.runDragonKingDescent(sess, p)
}

// runDragonKingDescent 龍王降臨主循環 goroutine
func (g *Game) runDragonKingDescent(sess *dragonKingSession, triggerPlayer *player.Player) {
	ticker := time.NewTicker(LuckyDragonKingBreathInterval)
	defer ticker.Stop()

	maxBreaths := int(LuckyDragonKingDuration / LuckyDragonKingBreathInterval)

	for {
		select {
		case <-ticker.C:
			sess.mu.Lock()
			sess.breathCount++
			currentBreath := sess.breathCount
			sess.mu.Unlock()

			if currentBreath > maxBreaths {
				// 觸發龍王爆發
				g.doDragonKingBurst(sess, triggerPlayer)
				return
			}
			g.doDragonBreath(sess, currentBreath)

		case <-g.stopCh:
			return
		}
	}
}

// doDragonBreath 執行一次龍息攻擊
func (g *Game) doDragonBreath(sess *dragonKingSession, breathNum int) {
	// 隨機選取場上目標
	g.mu.RLock()
	var aliveTargets []*target.Target
	for _, t := range g.Targets {
		if t.IsAlive {
			aliveTargets = append(aliveTargets, t)
		}
	}
	g.mu.RUnlock()

	if len(aliveTargets) == 0 {
		return
	}

	// Fisher-Yates 隨機打亂
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(aliveTargets) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		aliveTargets[i], aliveTargets[j] = aliveTargets[j], aliveTargets[i]
	}

	breathCount := LuckyDragonKingBreathTargets
	if len(aliveTargets) < breathCount {
		breathCount = len(aliveTargets)
	}

	totalReward := 0
	hitCount := 0

	for i := 0; i < breathCount; i++ {
		t := aliveTargets[i]

		// 80% 擊破機率
		if rng.Float64() > LuckyDragonKingBreathChance {
			continue
		}

		// 消滅目標
		g.mu.Lock()
		existing, ok := g.Targets[t.InstanceID]
		if !ok || !existing.IsAlive {
			g.mu.Unlock()
			continue
		}
		existing.IsAlive = false
		delete(g.Targets, t.InstanceID)
		g.mu.Unlock()

		// 計算全服共享獎勵
		avgBet := g.getAvgBetCost()
		reward := int(float64(avgBet) * t.Multiplier * LuckyDragonKingBreathMult)
		if reward < 1 {
			reward = 1
		}
		g.distributeRewardToAll(reward)
		totalReward += reward
		hitCount++
	}

	if hitCount == 0 {
		return
	}

	log.Printf("[DragonKing] 第 %d 次龍息！命中 %d 個目標，全服獎勵 %d",
		breathNum, hitCount, totalReward)

	// 全服廣播龍息攻擊
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyDragonKing,
		Payload: ws.LuckyDragonKingPayload{
			Event:       "dragon_breath",
			PlayerName:  sess.triggerPlayerName,
			BreathNum:   breathNum,
			HitCount:    hitCount,
			Mult:        LuckyDragonKingBreathMult,
			TotalReward: totalReward,
		},
	})
}

// doDragonKingBurst 龍王爆發（降臨結束時）
func (g *Game) doDragonKingBurst(sess *dragonKingSession, triggerPlayer *player.Player) {
	// 清除活躍會話
	m := g.LuckyDragonKing
	m.mu.Lock()
	m.activeSession = nil
	// 給觸發玩家龍王爆發倍率加成
	m.burstBoosts[triggerPlayer.ID] = time.Now().Add(LuckyDragonKingBurstDuration)
	m.mu.Unlock()

	// 場上所有目標 HP -60%
	g.mu.Lock()
	drainCount := 0
	for _, t := range g.Targets {
		if t.IsAlive {
			newHP := int(float64(t.HP) * (1.0 - LuckyDragonKingBurstHPDrain))
			if newHP < 1 {
				newHP = 1
			}
			t.HP = newHP
			drainCount++
		}
	}
	g.mu.Unlock()

	log.Printf("[DragonKing] 龍王爆發！%d 個目標 HP -60%%，player=%s 獲得 ×%.1f 倍率加成 %ds",
		drainCount, triggerPlayer.ID, LuckyDragonKingBurstMult, int(LuckyDragonKingBurstDuration.Seconds()))

	// 個人訊息：龍王爆發（含倍率加成）
	_ = g.Hub.Send(triggerPlayer.ID, &ws.Message{
		Type: ws.MsgLuckyDragonKing,
		Payload: ws.LuckyDragonKingPayload{
			Event:       "dragon_king_burst",
			PlayerID:    triggerPlayer.ID,
			PlayerName:  triggerPlayer.DisplayName,
			DrainCount:  drainCount,
			BurstMult:   LuckyDragonKingBurstMult,
			BurstSec:    int(LuckyDragonKingBurstDuration.Seconds()),
		},
	})

	// 全服廣播爆發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyDragonKing,
		Payload: ws.LuckyDragonKingPayload{
			Event:      "dragon_king_burst_broadcast",
			PlayerName: sess.triggerPlayerName,
			DrainCount: drainCount,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyDragonKing, sess.triggerPlayerName, 0, map[string]string{
		"message": fmt.Sprintf("🐉 龍王爆發！%d 個目標 HP -60%%！%s 獲得 ×%.1f 倍率加成 %d 秒！",
			drainCount, sess.triggerPlayerName, LuckyDragonKingBurstMult, int(LuckyDragonKingBurstDuration.Seconds())),
		"color": "#FF4500",
	})
}
