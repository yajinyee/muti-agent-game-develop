// lucky_meteor_shower_handler.go — 幸運星際隕石魚系統（DAY-253）
// 業界原創「隕石雨+隨機轟炸+隕石連擊+最終隕石」機制
//
// 設計：擊破 T211 後，天空降下「隕石雨」（持續 8 秒）：
//   - 每 1 秒隨機轟炸場上 2 個目標（70% 擊破機率，×1.3 倍率，全服共享）
//   - 若連續 3 次都命中同一個目標 → 「隕石連擊」：×3.0 倍率（全服大獎）
//   - 8 秒後「最終隕石」：場上最高 HP 目標被 100% 擊破（×2.0 倍率，全服共享）
//   - 個人冷卻 20 秒；全服冷卻 30 秒
//
// 設計差異：
//   - 與黑洞爆炸魚（T207，吸收+能量爆炸）不同，隕石雨是「隨機轟炸」，讓玩家看到「隕石從天而降砸中魚」的視覺爽感
//   - 「隕石連擊」讓玩家有「要看哪條魚被連續砸中」的期待感，製造「哇，那條魚被砸了 3 次」的驚嘆感
//   - 「最終隕石」讓玩家有「等待→最後一擊」的高潮設計，確保最高 HP 目標被消滅
//   - 「全服共享獎勵」讓所有玩家都受益，製造「全服一起看隕石雨」的社交感
//   - 「70% 擊破機率」讓玩家有「這次會不會砸中」的刺激感，不是 100% 確定
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/ws"
)

const (
	LuckyMeteorPersonalCD   = 20 * time.Second // 個人冷卻
	LuckyMeteorGlobalCD     = 30 * time.Second // 全服冷卻
	LuckyMeteorDuration     = 8 * time.Second  // 隕石雨持續時間
	LuckyMeteorInterval     = 1 * time.Second  // 每次轟炸間隔
	LuckyMeteorBombCount    = 2                // 每次轟炸目標數
	LuckyMeteorKillChance   = 0.70             // 每次轟炸擊破機率
	LuckyMeteorBombMult     = 1.3              // 轟炸倍率（全服共享）
	LuckyMeteorComboMult    = 3.0              // 隕石連擊倍率（全服大獎）
	LuckyMeteorComboCount   = 3                // 連擊所需次數
	LuckyMeteorFinalMult    = 2.0              // 最終隕石倍率（全服共享）
)

// meteorShowerSession 隕石雨會話
type meteorShowerSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	expiresAt         time.Time
	// 連擊追蹤（instanceID → 連續命中次數）
	hitStreak map[string]int
	mu        sync.Mutex
}

// luckyMeteorShowerManager 幸運星際隕石魚管理器
type luckyMeteorShowerManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前活躍的隕石雨會話（nil = 無）
	activeSession *meteorShowerSession
}

func newLuckyMeteorShowerManager() *luckyMeteorShowerManager {
	return &luckyMeteorShowerManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyMeteorShowerFish 判斷是否為幸運星際隕石魚
func isLuckyMeteorShowerFish(defID string) bool {
	return defID == "T211"
}

// isMeteorShowerActive 判斷隕石雨是否正在進行
func (m *luckyMeteorShowerManager) isMeteorShowerActive() bool {
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

// tryLuckyMeteorShowerFish 擊破 T211 後觸發隕石雨
func (g *Game) tryLuckyMeteorShowerFish(triggerPlayerID, triggerPlayerName string) {
	m := g.LuckyMeteorShower

	m.mu.Lock()
	now := time.Now()

	// 全服冷卻檢查
	if now.Before(m.globalCooldownUntil) {
		m.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if cd, ok := m.personalCooldowns[triggerPlayerID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	// 已有活躍隕石雨
	if m.activeSession != nil && now.Before(m.activeSession.expiresAt) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[triggerPlayerID] = now.Add(LuckyMeteorPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyMeteorGlobalCD)

	expiresAt := now.Add(LuckyMeteorDuration)
	sess := &meteorShowerSession{
		triggerPlayerID:   triggerPlayerID,
		triggerPlayerName: triggerPlayerName,
		expiresAt:         expiresAt,
		hitStreak:         make(map[string]int),
	}
	m.activeSession = sess
	m.mu.Unlock()

	log.Printf("[MeteorShower] player=%s 觸發隕石雨！持續 %ds", triggerPlayerID, int(LuckyMeteorDuration.Seconds()))

	// 全服廣播：隕石雨開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMeteorShower,
		Payload: ws.LuckyMeteorShowerPayload{
			Event:       "meteor_start",
			PlayerName:  triggerPlayerName,
			DurationSec: int(LuckyMeteorDuration.Seconds()),
			BombMult:    LuckyMeteorBombMult,
			FinalMult:   LuckyMeteorFinalMult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyMeteorShower, triggerPlayerName, 0, map[string]string{
		"message": fmt.Sprintf("☄️ %s 觸發星際隕石雨！每秒轟炸 %d 個目標！×%.1f 倍率！全服共享！",
			triggerPlayerName, LuckyMeteorBombCount, LuckyMeteorBombMult),
		"color": "#E74C3C",
	})

	// 啟動隕石雨 goroutine
	go g.runMeteorShower(sess)
}

// runMeteorShower 隕石雨主循環 goroutine
func (g *Game) runMeteorShower(sess *meteorShowerSession) {
	ticker := time.NewTicker(LuckyMeteorInterval)
	defer ticker.Stop()

	bombRound := 0
	maxRounds := int(LuckyMeteorDuration / LuckyMeteorInterval)

	for {
		select {
		case <-ticker.C:
			bombRound++
			if bombRound > maxRounds {
				// 觸發最終隕石
				g.doMeteorFinalStrike(sess)
				return
			}
			g.doMeteorBomb(sess, bombRound)

		case <-g.stopCh:
			return
		}
	}
}

// doMeteorBomb 執行一輪隕石轟炸
func (g *Game) doMeteorBomb(sess *meteorShowerSession, round int) {
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

	// Fisher-Yates 隨機打亂，取前 N 個
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(aliveTargets) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		aliveTargets[i], aliveTargets[j] = aliveTargets[j], aliveTargets[i]
	}

	bombCount := LuckyMeteorBombCount
	if len(aliveTargets) < bombCount {
		bombCount = len(aliveTargets)
	}

	for i := 0; i < bombCount; i++ {
		t := aliveTargets[i]

		// 70% 擊破機率
		if rng.Float64() > LuckyMeteorKillChance {
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
		reward := int(float64(avgBet) * t.Multiplier * LuckyMeteorBombMult)
		if reward < 1 {
			reward = 1
		}
		g.distributeRewardToAll(reward)

		log.Printf("[MeteorShower] 第 %d 輪轟炸命中 %s（%s）！全服獎勵 %d",
			round, t.InstanceID, t.DefID, reward)

		// 更新連擊計數
		sess.mu.Lock()
		sess.hitStreak[t.InstanceID]++
		streak := sess.hitStreak[t.InstanceID]
		sess.mu.Unlock()

		// 廣播轟炸命中
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyMeteorShower,
			Payload: ws.LuckyMeteorShowerPayload{
				Event:      "meteor_bomb",
				PlayerName: sess.triggerPlayerName,
				InstanceID: t.InstanceID,
				TargetName: t.DefID,
				Round:      round,
				Mult:       LuckyMeteorBombMult,
				Reward:     reward,
				X:          t.X,
				Y:          t.Y,
			},
		})

		// 隕石連擊判斷（連續 3 次命中同一目標）
		if streak >= LuckyMeteorComboCount {
			g.doMeteorCombo(sess, t, reward)
			// 重置連擊計數
			sess.mu.Lock()
			sess.hitStreak[t.InstanceID] = 0
			sess.mu.Unlock()
		}
	}
}

// doMeteorCombo 隕石連擊（連續 3 次命中同一目標）
func (g *Game) doMeteorCombo(sess *meteorShowerSession, t *target.Target, baseReward int) {
	// 連擊獎勵 = 基礎獎勵 × (comboMult / bombMult)，避免重複計算
	comboBonus := int(float64(baseReward) * (LuckyMeteorComboMult / LuckyMeteorBombMult))
	if comboBonus < 1 {
		comboBonus = 1
	}
	g.distributeRewardToAll(comboBonus)

	log.Printf("[MeteorShower] 隕石連擊！目標 %s 被連續命中 %d 次！全服額外獎勵 %d",
		t.InstanceID, LuckyMeteorComboCount, comboBonus)

	// 全服廣播連擊
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMeteorShower,
		Payload: ws.LuckyMeteorShowerPayload{
			Event:      "meteor_combo",
			PlayerName: sess.triggerPlayerName,
			InstanceID: t.InstanceID,
			TargetName: t.DefID,
			Mult:       LuckyMeteorComboMult,
			Reward:     comboBonus,
			X:          t.X,
			Y:          t.Y,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyMeteorShower, sess.triggerPlayerName, comboBonus, map[string]string{
		"message": fmt.Sprintf("☄️ 隕石連擊！%s 被連續命中 %d 次！全服獲得 %d 籌碼！",
			t.DefID, LuckyMeteorComboCount, comboBonus),
		"color": "#FF4500",
	})
}

// doMeteorFinalStrike 最終隕石（隕石雨結束時）
func (g *Game) doMeteorFinalStrike(sess *meteorShowerSession) {
	// 清除活躍會話
	m := g.LuckyMeteorShower
	m.mu.Lock()
	m.activeSession = nil
	m.mu.Unlock()

	// 找場上最高 HP 的目標
	g.mu.RLock()
	var finalTarget *target.Target
	maxHP := 0
	for _, t := range g.Targets {
		if t.IsAlive && t.HP > maxHP {
			maxHP = t.HP
			finalTarget = t
		}
	}
	g.mu.RUnlock()

	if finalTarget == nil {
		// 沒有目標，只廣播結束
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyMeteorShower,
			Payload: ws.LuckyMeteorShowerPayload{
				Event:      "meteor_end",
				PlayerName: sess.triggerPlayerName,
			},
		})
		return
	}

	// 100% 擊破最高 HP 目標
	g.mu.Lock()
	existing, ok := g.Targets[finalTarget.InstanceID]
	if !ok || !existing.IsAlive {
		g.mu.Unlock()
		// 目標已消失，廣播結束
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyMeteorShower,
			Payload: ws.LuckyMeteorShowerPayload{
				Event:      "meteor_end",
				PlayerName: sess.triggerPlayerName,
			},
		})
		return
	}
	existing.IsAlive = false
	delete(g.Targets, finalTarget.InstanceID)
	g.mu.Unlock()

	// 計算全服共享獎勵
	avgBet := g.getAvgBetCost()
	finalReward := int(float64(avgBet) * finalTarget.Multiplier * LuckyMeteorFinalMult)
	if finalReward < 1 {
		finalReward = 1
	}
	g.distributeRewardToAll(finalReward)

	log.Printf("[MeteorShower] 最終隕石！命中 %s（%s）！全服獎勵 %d",
		finalTarget.InstanceID, finalTarget.DefID, finalReward)

	// 全服廣播最終隕石
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMeteorShower,
		Payload: ws.LuckyMeteorShowerPayload{
			Event:       "meteor_final",
			PlayerName:  sess.triggerPlayerName,
			InstanceID:  finalTarget.InstanceID,
			TargetName:  finalTarget.DefID,
			Mult:        LuckyMeteorFinalMult,
			Reward:      finalReward,
			X:           finalTarget.X,
			Y:           finalTarget.Y,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyMeteorShower, sess.triggerPlayerName, finalReward, map[string]string{
		"message": fmt.Sprintf("☄️ 最終隕石！命中 %s！全服獲得 %d 籌碼！",
			finalTarget.DefID, finalReward),
		"color": "#C0392B",
	})
}

// end of lucky_meteor_shower_handler.go
