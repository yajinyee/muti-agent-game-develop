// lucky_phantom_fish_handler.go — 幸運幽靈魚系統（DAY-245）
// 業界原創「幽靈殘影+死亡後復活攻擊」機制
//
// 設計：擊破 T203 後，玩家獲得「幽靈護盾」（12 秒）：
//   - 護盾期間，玩家每次擊破任何目標，目標留下「幽靈殘影」（持續 5 秒）
//   - 幽靈殘影可被再次擊破（50% 機率，×1.5 倍率，個人獎勵）
//   - 12 秒後「幽靈爆發」：所有場上幽靈殘影同時爆炸（100% 擊破，×2.0 倍率，個人獎勵）
//   - 個人冷卻 22 秒；全服冷卻 35 秒
//
// 設計差異：
//   - 與分身魚（T200，同時三方向射擊）不同，幽靈魚是「死亡後留下殘影可再次擊破」
//     讓玩家有「打死一條魚還能再賺一次」的爽感
//   - 「幽靈爆發」讓玩家有「等待殘影累積再一次爆發」的策略感
//   - 「50% 機率擊破殘影」讓玩家有「要不要賭一把」的刺激感
//   - 「×2.0 爆發倍率 > ×1.5 單次倍率」鼓勵玩家等待爆發而非逐一擊破
//   - 全服廣播讓其他玩家看到「有人觸發了幽靈護盾」，製造羨慕感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"

	"github.com/google/uuid"
)

const (
	LuckyPhantomFishPersonalCD   = 22 * time.Second // 個人冷卻
	LuckyPhantomFishGlobalCD     = 35 * time.Second // 全服冷卻
	LuckyPhantomFishDuration     = 12 * time.Second // 幽靈護盾持續時間
	LuckyPhantomFishGhostTTL     = 5 * time.Second  // 幽靈殘影存活時間
	LuckyPhantomFishGhostKillMult = 1.5             // 幽靈殘影擊破倍率
	LuckyPhantomFishBurstMult    = 2.0              // 幽靈爆發倍率
	LuckyPhantomFishGhostKillChance = 0.5           // 幽靈殘影擊破機率
)

// phantomGhost 幽靈殘影
type phantomGhost struct {
	ghostID     string
	origDefID   string
	x           float64
	y           float64
	expiresAt   time.Time
}

// phantomSession 幽靈護盾 session
type phantomSession struct {
	playerID  string
	expiresAt time.Time
	// 幽靈殘影列表（ghostID → ghost）
	ghosts map[string]*phantomGhost
	mu     sync.Mutex
}

// luckyPhantomFishManager 幸運幽靈魚管理器
type luckyPhantomFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前幽靈護盾 sessions（playerID → session）
	activeSessions map[string]*phantomSession
}

func newLuckyPhantomFishManager() *luckyPhantomFishManager {
	return &luckyPhantomFishManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*phantomSession),
	}
}

// isLuckyPhantomFish 判斷是否為幸運幽靈魚
func isLuckyPhantomFish(defID string) bool {
	return defID == "T203"
}

// isPhantomShieldActive 判斷玩家是否有幽靈護盾
func (m *luckyPhantomFishManager) isPhantomShieldActive(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	if !ok {
		return false
	}
	return time.Now().Before(sess.expiresAt)
}

// tryLuckyPhantomFish 擊破 T203 後觸發幽靈護盾
func (g *Game) tryLuckyPhantomFish(p *player.Player) {
	// 計算個人 betCost
	betDef := data.GetBetDef(p.BetLevel)
	avgBet := betDef.BetCost
	if avgBet < 1 {
		avgBet = 1
	}

	m := g.LuckyPhantomFish
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
	m.personalCooldowns[p.ID] = now.Add(LuckyPhantomFishPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyPhantomFishGlobalCD)

	// 建立 session
	sess := &phantomSession{
		playerID:  p.ID,
		expiresAt: now.Add(LuckyPhantomFishDuration),
		ghosts:    make(map[string]*phantomGhost),
	}
	m.activeSessions[p.ID] = sess
	m.mu.Unlock()

	log.Printf("[PhantomFish] player=%s 幽靈護盾啟動（持續 %v）", p.ID, LuckyPhantomFishDuration)

	// 個人訊息：幽靈護盾啟動
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyPhantomFish,
		Payload: ws.LuckyPhantomFishPayload{
			Event:           "phantom_start",
			PlayerID:        p.ID,
			PlayerName:      p.DisplayName,
			DurationSec:     int(LuckyPhantomFishDuration.Seconds()),
			GhostKillMult:   LuckyPhantomFishGhostKillMult,
			BurstMult:       LuckyPhantomFishBurstMult,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyPhantomFish,
		Payload: ws.LuckyPhantomFishPayload{
			Event:      "phantom_broadcast",
			PlayerName: p.DisplayName,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyPhantomFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("👻 %s 觸發幽靈護盾！擊破目標留下殘影，12秒後幽靈爆發！", p.DisplayName),
		"color":   "#8E44AD",
	})

	// 啟動護盾計時 goroutine
	go g.runLuckyPhantomShield(p, sess)
}

// runLuckyPhantomShield 幽靈護盾計時 goroutine
func (g *Game) runLuckyPhantomShield(p *player.Player, sess *phantomSession) {
	timer := time.NewTimer(LuckyPhantomFishDuration)
	defer timer.Stop()

	select {
	case <-timer.C:
		g.doPhantomBurst(p, sess)
	case <-g.stopCh:
		return
	}
}

// createPhantomGhost 擊破目標後建立幽靈殘影（由 handleKill 呼叫）
func (g *Game) createPhantomGhost(p *player.Player, targetDefID string, x, y float64) {
	// 計算個人 betCost
	betDef := data.GetBetDef(p.BetLevel)
	avgBet := betDef.BetCost
	if avgBet < 1 {
		avgBet = 1
	}

	m := g.LuckyPhantomFish
	m.mu.Lock()
	sess, ok := m.activeSessions[p.ID]
	if !ok {
		m.mu.Unlock()
		return
	}
	if time.Now().After(sess.expiresAt) {
		m.mu.Unlock()
		return
	}

	ghostID := uuid.New().String()
	ghost := &phantomGhost{
		ghostID:   ghostID,
		origDefID: targetDefID,
		x:         x,
		y:         y,
		expiresAt: time.Now().Add(LuckyPhantomFishGhostTTL),
	}
	sess.mu.Lock()
	sess.ghosts[ghostID] = ghost
	sess.mu.Unlock()
	m.mu.Unlock()

	log.Printf("[PhantomFish] player=%s 幽靈殘影生成 ghostID=%s defID=%s", p.ID, ghostID, targetDefID)

	// 通知玩家幽靈殘影生成
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyPhantomFish,
		Payload: ws.LuckyPhantomFishPayload{
			Event:            "phantom_ghost_created",
			GhostID:          ghostID,
			OrigTargetDefID:  targetDefID,
			X:                x,
			Y:                y,
			GhostDurationSec: int(LuckyPhantomFishGhostTTL.Seconds()),
		},
	})

	// 啟動殘影過期清理 goroutine
	go func() {
		timer := time.NewTimer(LuckyPhantomFishGhostTTL)
		defer timer.Stop()
		select {
		case <-timer.C:
			m.mu.Lock()
			if s, ok := m.activeSessions[p.ID]; ok {
				s.mu.Lock()
				delete(s.ghosts, ghostID)
				s.mu.Unlock()
			}
			m.mu.Unlock()
		case <-g.stopCh:
		}
	}()
}

// tryHitPhantomGhost 嘗試擊破幽靈殘影（由 handleAttack 呼叫）
// 回傳是否命中殘影
func (g *Game) tryHitPhantomGhost(p *player.Player) bool {
	// 計算個人 betCost
	betDef := data.GetBetDef(p.BetLevel)
	avgBet := betDef.BetCost
	if avgBet < 1 {
		avgBet = 1
	}

	m := g.LuckyPhantomFish
	m.mu.Lock()
	sess, ok := m.activeSessions[p.ID]
	if !ok {
		m.mu.Unlock()
		return false
	}
	if time.Now().After(sess.expiresAt) {
		m.mu.Unlock()
		return false
	}

	// 找一個最近的幽靈殘影
	sess.mu.Lock()
	var target *phantomGhost
	for _, g2 := range sess.ghosts {
		if time.Now().Before(g2.expiresAt) {
			target = g2
			break
		}
	}
	if target == nil {
		sess.mu.Unlock()
		m.mu.Unlock()
		return false
	}
	ghostID := target.ghostID
	sess.mu.Unlock()
	m.mu.Unlock()

	// 50% 機率擊破
	if rand.Float64() >= LuckyPhantomFishGhostKillChance {
		return false
	}

	// 擊破殘影
	m.mu.Lock()
	if s, ok := m.activeSessions[p.ID]; ok {
		s.mu.Lock()
		delete(s.ghosts, ghostID)
		s.mu.Unlock()
	}
	m.mu.Unlock()

	reward := int(float64(avgBet) * LuckyPhantomFishGhostKillMult)
	p.AddCoins(reward)

	log.Printf("[PhantomFish] player=%s 擊破幽靈殘影 ghostID=%s reward=%d", p.ID, ghostID, reward)

	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyPhantomFish,
		Payload: ws.LuckyPhantomFishPayload{
			Event:    "phantom_ghost_killed",
			GhostID:  ghostID,
			Reward:   reward,
			KillMult: LuckyPhantomFishGhostKillMult,
		},
	})

	return true
}

// doPhantomBurst 幽靈爆發（護盾結束時觸發）
func (g *Game) doPhantomBurst(p *player.Player, sess *phantomSession) {
	// 計算個人 betCost
	betDef := data.GetBetDef(p.BetLevel)
	avgBet := betDef.BetCost
	if avgBet < 1 {
		avgBet = 1
	}

	m := g.LuckyPhantomFish

	// 收集所有存活的幽靈殘影
	sess.mu.Lock()
	now := time.Now()
	var aliveGhosts []*phantomGhost
	for _, ghost := range sess.ghosts {
		if now.Before(ghost.expiresAt) {
			aliveGhosts = append(aliveGhosts, ghost)
		}
	}
	// 清空殘影
	sess.ghosts = make(map[string]*phantomGhost)
	sess.mu.Unlock()

	// 移除 session
	m.mu.Lock()
	delete(m.activeSessions, p.ID)
	m.mu.Unlock()

	ghostCount := len(aliveGhosts)
	totalReward := 0

	if ghostCount > 0 {
		// 每個殘影 100% 爆炸，×2.0 倍率
		for range aliveGhosts {
			reward := int(float64(avgBet) * LuckyPhantomFishBurstMult)
			totalReward += reward
			p.AddCoins(reward)
		}
		log.Printf("[PhantomFish] player=%s 幽靈爆發 ghostCount=%d totalReward=%d", p.ID, ghostCount, totalReward)
	}

	// 通知玩家幽靈爆發結果
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyPhantomFish,
		Payload: ws.LuckyPhantomFishPayload{
			Event:       "phantom_burst",
			GhostCount:  ghostCount,
			TotalReward: totalReward,
			BurstMult:   LuckyPhantomFishBurstMult,
		},
	})

	// 護盾結束通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyPhantomFish,
		Payload: ws.LuckyPhantomFishPayload{
			Event: "phantom_end",
		},
	})

	if ghostCount >= 3 {
		g.Announce.Create(announce.EventLuckyPhantomFish, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("👻 %s 幽靈爆發！%d 個殘影同時爆炸！獲得 %d 籌碼！", p.DisplayName, ghostCount, totalReward),
			"color":   "#6C3483",
		})
	}
}
