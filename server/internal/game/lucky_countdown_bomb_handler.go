// lucky_countdown_bomb_handler.go — 幸運倒數炸彈魚系統（DAY-268）
// 業界原創「倒數充能+全服爆炸」機制
//
// 設計：擊破 T226 後，場上出現「倒數炸彈」（10 秒倒數）：
//   - 倒數期間，任何玩家每次擊破任何目標，炸彈充能 +1（最多 10 次）
//   - 10 秒後炸彈爆炸：充能數 × ×1.5 倍率（全服共享 AOE）
//   - 若充能達到 10 次，提前引爆：×3.0 倍率（全服大獎）
//   - 個人冷卻 28 秒；全服冷卻 45 秒
//
// 設計差異：
//   - 與倍率疊加（T225，個人累積）不同，倒數炸彈是「全服合力充能」，讓所有玩家有「快快快，還差幾個！」的緊迫感
//   - 「充能數 × ×1.5」讓玩家有「充能越多爆炸越強」的動力
//   - 「滿充能提前引爆 ×3.0」讓玩家有「要趁 10 秒內全服打滿 10 個目標」的策略感
//   - 「倒數計時器全服廣播」讓所有玩家看到「還剩幾秒」，製造緊迫感
//   - 「充能進度全服廣播」讓所有玩家看到「現在充能了幾個」，製造「快滿了！」的期待感
//   - 業界原創：結合「倒數計時+全服合力充能+爆炸獎勵」三個元素，製造全服社交緊迫感
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
	LuckyCountdownBombPersonalCD  = 28 * time.Second // 個人冷卻
	LuckyCountdownBombGlobalCD    = 45 * time.Second // 全服冷卻
	LuckyCountdownBombDuration    = 10 * time.Second // 倒數時間
	LuckyCountdownBombMaxCharge   = 10               // 最大充能次數
	LuckyCountdownBombNormalMult  = 1.5              // 普通爆炸倍率（每充能 ×1.5）
	LuckyCountdownBombBurstMult   = 3.0              // 滿充能爆炸倍率
)

// countdownBombSession 倒數炸彈 session
type countdownBombSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	expiresAt         time.Time
	chargeCount       int  // 當前充能次數
	exploded          bool // 是否已爆炸
}

// luckyCountdownBombManager 幸運倒數炸彈魚管理器
type luckyCountdownBombManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 活躍 session（全服只有一個）
	activeSession *countdownBombSession
}

func newLuckyCountdownBombManager() *luckyCountdownBombManager {
	return &luckyCountdownBombManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyCountdownBombFish 判斷是否為幸運倒數炸彈魚
func isLuckyCountdownBombFish(defID string) bool {
	return defID == "T226"
}

// isCountdownBombActive 判斷倒數炸彈是否活躍（供 handleKill 使用）
func (m *luckyCountdownBombManager) isCountdownBombActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil {
		return false
	}
	if time.Now().After(m.activeSession.expiresAt) {
		m.activeSession = nil
		return false
	}
	return !m.activeSession.exploded
}

// tryLuckyCountdownBombFish 擊破 T226 後觸發倒數炸彈
func (g *Game) tryLuckyCountdownBombFish(p *player.Player) {
	m := g.LuckyCountdownBomb

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
	// 已有活躍 session 時不重複觸發
	if m.activeSession != nil && !m.activeSession.exploded && now.Before(m.activeSession.expiresAt) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyCountdownBombPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyCountdownBombGlobalCD)

	// 建立 session
	sess := &countdownBombSession{
		triggerPlayerID:   p.ID,
		triggerPlayerName: p.DisplayName,
		expiresAt:         now.Add(LuckyCountdownBombDuration),
		chargeCount:       0,
		exploded:          false,
	}
	m.activeSession = sess
	m.mu.Unlock()

	log.Printf("[CountdownBomb] player=%s 觸發倒數炸彈！10 秒倒數開始", p.ID)

	// 全服廣播：炸彈出現
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyCountdownBomb,
		Payload: ws.LuckyCountdownBombPayload{
			Event:      "bomb_start",
			PlayerName: p.DisplayName,
			Countdown:  LuckyCountdownBombDuration.Seconds(),
			MaxCharge:  LuckyCountdownBombMaxCharge,
			NormalMult: LuckyCountdownBombNormalMult,
			BurstMult:  LuckyCountdownBombBurstMult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyCountdownBomb, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("💣 %s 觸發倒數炸彈！10 秒倒數！全服充能越多爆炸越強！",
			p.DisplayName),
		"color": "#FF4500",
	})

	// 啟動倒數 goroutine
	go g.runCountdownBombTimer(p, sess)
}

// notifyCountdownBombKill 任何玩家擊破任何非 T226 目標時呼叫（充能）
func (g *Game) notifyCountdownBombKill(p *player.Player) {
	m := g.LuckyCountdownBomb

	m.mu.Lock()
	if m.activeSession == nil || m.activeSession.exploded || time.Now().After(m.activeSession.expiresAt) {
		m.mu.Unlock()
		return
	}

	m.activeSession.chargeCount++
	chargeCount := m.activeSession.chargeCount
	isFull := chargeCount >= LuckyCountdownBombMaxCharge
	if isFull {
		m.activeSession.exploded = true
	}
	m.mu.Unlock()

	// 廣播充能進度
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyCountdownBomb,
		Payload: ws.LuckyCountdownBombPayload{
			Event:       "bomb_charge",
			ChargeCount: chargeCount,
			MaxCharge:   LuckyCountdownBombMaxCharge,
			ChargerName: p.DisplayName,
		},
	})

	// 滿充能 → 提前引爆
	if isFull {
		log.Printf("[CountdownBomb] 滿充能！提前引爆！×%.1f", LuckyCountdownBombBurstMult)
		go g.doCountdownBombExplode(chargeCount, true)
	}
}

// runCountdownBombTimer 倒數計時 goroutine
func (g *Game) runCountdownBombTimer(p *player.Player, sess *countdownBombSession) {
	timer := time.NewTimer(LuckyCountdownBombDuration)
	defer timer.Stop()

	<-timer.C

	m := g.LuckyCountdownBomb
	m.mu.Lock()
	// 確認 session 仍然存在且未爆炸
	if m.activeSession == nil || m.activeSession.exploded {
		m.mu.Unlock()
		return
	}
	chargeCount := m.activeSession.chargeCount
	m.activeSession.exploded = true
	m.activeSession = nil
	m.mu.Unlock()

	log.Printf("[CountdownBomb] 倒數結束！充能 %d 次，普通爆炸", chargeCount)
	go g.doCountdownBombExplode(chargeCount, false)
}

// doCountdownBombExplode 執行炸彈爆炸
func (g *Game) doCountdownBombExplode(chargeCount int, isBurst bool) {
	// 計算爆炸倍率
	var mult float64
	if isBurst {
		mult = LuckyCountdownBombBurstMult
	} else {
		mult = float64(chargeCount) * LuckyCountdownBombNormalMult
		if mult < 1.0 {
			mult = 1.0
		}
	}

	// 對場上所有目標施加爆炸傷害並給予全服獎勵
	g.mu.RLock()
	targets := make([]targetSnapshot, 0, len(g.Targets))
	for _, t := range g.Targets {
		if t.HP > 0 && !isLuckyCountdownBombFish(t.DefID) {
			name := t.DefID
			if t.Def != nil {
				name = t.Def.Name
			}
			targets = append(targets, targetSnapshot{
				instanceID: t.InstanceID,
				defID:      t.DefID,
				name:       name,
				x:          t.X,
				y:          t.Y,
				hp:         t.HP,
				multiplier: t.Multiplier,
			})
		}
	}
	g.mu.RUnlock()

	// 計算全服總獎勵（基於平均 bet cost × 倍率 × 目標數）
	avgBet := g.getAvgBetCost()
	totalReward := int(float64(avgBet) * mult * float64(len(targets)))
	if totalReward < 1 && len(targets) > 0 {
		totalReward = len(targets)
	}

	// 對所有目標施加 HP -40% 傷害
	for _, t := range targets {
		g.applyExplosionDamage(t.instanceID, 0.40)
	}

	// 廣播爆炸結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyCountdownBomb,
		Payload: ws.LuckyCountdownBombPayload{
			Event:       "bomb_explode",
			ChargeCount: chargeCount,
			Mult:        mult,
			TotalReward: totalReward,
			IsBurst:     isBurst,
		},
	})

	// 全服公告
	if isBurst {
		g.Announce.Create(announce.EventLuckyCountdownBomb, "", totalReward, map[string]string{
			"message": fmt.Sprintf("💣 倒數炸彈滿充能爆炸！×%.1f！全服 AOE +%d！",
				mult, totalReward),
			"color": "#FFD700",
		})
	} else if chargeCount >= 5 {
		g.Announce.Create(announce.EventLuckyCountdownBomb, "", totalReward, map[string]string{
			"message": fmt.Sprintf("💣 倒數炸彈爆炸！充能 %d 次 ×%.1f！全服 AOE +%d！",
				chargeCount, mult, totalReward),
			"color": "#FF6B35",
		})
	}

	log.Printf("[CountdownBomb] 爆炸完成！充能 %d 次，倍率 ×%.1f，全服獎勵 %d，目標數 %d",
		chargeCount, mult, totalReward, len(targets))
}
