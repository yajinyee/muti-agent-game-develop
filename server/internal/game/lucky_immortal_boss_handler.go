// lucky_immortal_boss_handler.go — 幸運永生 BOSS 魚系統（DAY-289）
// 業界依據：Royal Fishing Jili「Immortal Boss mechanic — Golden Toad and Ancient Crocodile bosses
//          appear randomly and award consecutive wins ranging from 50X to 150X until they leave the screen」
//          業界原創「永生 BOSS 降臨 + 5 條命連續復活 + 倍率遞增 + 最終爆發」機制
//
// 設計：
//   - 擊破 T247 後，召喚「永生 BOSS」（在場上隨機位置出現，持續 18 秒）
//   - 永生 BOSS 有 5 條命（每次被擊破後立即復活，HP 恢復 100%）
//   - 每次擊破永生 BOSS → 觸發玩家獲得 ×2.0 獎勵（第1次）
//   - 每次復活時，擊破倍率提升 +0.5x（第1次 ×2.0 → 第2次 ×2.5 → ... → 第5次 ×4.0）
//   - 若 5 條命全部在 18 秒內耗盡 → 「永生終結」：全服 ×3.5 加成 8 秒
//   - 全服廣播永生 BOSS 位置和每次復活
//   - 個人冷卻 32 秒；全服冷卻 50 秒
//
// 設計差異：
//   - 與多米諾（T246，場上 5 個目標依序連鎖倒下）不同，永生 BOSS 是「同一個目標反覆復活」
//   - 「每次復活倍率 +0.5x」讓玩家有「越打越值錢，要撐到第 5 次」的動力
//   - 「5 條命全部耗盡觸發永生終結」讓玩家有「要趁 18 秒內打完 5 次」的緊迫感
//   - 「全服 ×3.5 加成 8 秒」是最高全服加成之一，製造「全服一起爽」的社交感
//   - 「全服廣播永生 BOSS 位置」讓所有玩家看到「永生 BOSS 在哪裡」，製造全服搶打感
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

const (
	LuckyImmortalBossPersonalCD = 32 * time.Second // 個人冷卻
	LuckyImmortalBossGlobalCD   = 50 * time.Second // 全服冷卻

	// 永生 BOSS 設計
	ImmortalBossLives       = 5                // 永生 BOSS 條命數
	ImmortalBossBaseMult    = 2.0              // 第 1 次擊破倍率
	ImmortalBossMultStep    = 0.5              // 每次復活倍率增加
	ImmortalBossTimeout     = 18 * time.Second // 永生 BOSS 有效時間

	// 永生終結：全服加成
	ImmortalBossEndMult     = 3.5             // 全服 ×3.5
	ImmortalBossEndDuration = 8 * time.Second // 持續 8 秒
)

// immortalBossEndBoost 永生終結全服加成
type immortalBossEndBoost struct {
	mult      float64
	expiresAt time.Time
}

// immortalBossSession 永生 BOSS 會話
type immortalBossSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	instanceID        string  // 永生 BOSS 的 instanceID（每次復活更新）
	livesLeft         int     // 剩餘條命
	killCount         int     // 已擊破次數
	currentMult       float64 // 當前擊破倍率
	x, y              float64 // 位置
	expiresAt         time.Time
	settled           bool
}

// luckyImmortalBossManager 幸運永生 BOSS 魚管理器
type luckyImmortalBossManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前永生 BOSS 會話（全服同時只有一個）
	activeSession *immortalBossSession

	// 永生終結全服加成
	endBoost *immortalBossEndBoost
}

func newLuckyImmortalBossManager() *luckyImmortalBossManager {
	return &luckyImmortalBossManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyImmortalBossFish 判斷是否為幸運永生 BOSS 魚
func isLuckyImmortalBossFish(defID string) bool {
	return defID == "T247"
}

// getImmortalBossEndMult 取得永生終結全服加成倍率（供 handleKill 使用）
func (m *luckyImmortalBossManager) getImmortalBossEndMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.endBoost != nil && time.Now().Before(m.endBoost.expiresAt) {
		return m.endBoost.mult
	}
	return 1.0
}

// isImmortalBossTarget 判斷是否為永生 BOSS 目標（返回是否為目標、當前倍率）
func (m *luckyImmortalBossManager) isImmortalBossTarget(instanceID string) (bool, float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil || m.activeSession.settled {
		return false, 0
	}
	if m.activeSession.instanceID == instanceID {
		return true, m.activeSession.currentMult
	}
	return false, 0
}

// tryLuckyImmortalBossFish 嘗試觸發幸運永生 BOSS 魚
func (g *Game) tryLuckyImmortalBossFish(p *player.Player) {
	m := g.LuckyImmortalBoss
	m.mu.Lock()

	now := time.Now()

	// 個人冷卻檢查
	if cd, ok := m.personalCooldowns[p.ID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	// 全服冷卻檢查
	if now.Before(m.globalCooldownUntil) {
		m.mu.Unlock()
		return
	}
	// 已有活躍會話
	if m.activeSession != nil && !m.activeSession.settled {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyImmortalBossPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyImmortalBossGlobalCD)

	// 隨機位置（場地 800x600，避開邊緣）
	x := 150.0 + rand.Float64()*500.0
	y := 100.0 + rand.Float64()*400.0

	// 建立會話（instanceID 先設為空，等 BOSS 實際出現時更新）
	session := &immortalBossSession{
		triggerPlayerID:   p.ID,
		triggerPlayerName: p.DisplayName,
		instanceID:        fmt.Sprintf("immortal_boss_%d", now.UnixNano()),
		livesLeft:         ImmortalBossLives,
		killCount:         0,
		currentMult:       ImmortalBossBaseMult,
		x:                 x,
		y:                 y,
		expiresAt:         now.Add(ImmortalBossTimeout),
	}
	m.activeSession = session
	m.mu.Unlock()

	log.Printf("[LuckyImmortalBoss] 觸發！玩家=%s 位置=(%.0f,%.0f) 條命=%d instanceID=%s",
		p.DisplayName, x, y, ImmortalBossLives, session.instanceID)

	// 廣播永生 BOSS 降臨
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyImmortalBoss,
		Payload: ws.LuckyImmortalBossPayload{
			Event:             "immortal_start",
			TriggerPlayerID:   p.ID,
			TriggerPlayerName: p.DisplayName,
			X:                 x,
			Y:                 y,
			LivesLeft:         ImmortalBossLives,
			CurrentMult:       ImmortalBossBaseMult,
			TimeoutSec:        int(ImmortalBossTimeout.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyImmortalBoss, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⚡ 永生 BOSS 降臨！%s 召喚了永生 BOSS！5 條命，打完全服 ×%.1f！", p.DisplayName, ImmortalBossEndMult),
		"color":   "#8B0000",
	})
	g.broadcastAnnouncement(ann)

	// 啟動超時計時器
	go g.runImmortalBossTimeout(session)
}

// notifyImmortalBossKill 永生 BOSS 被擊破通知（由 handleKill 呼叫）
// 返回本次擊破的倍率加成（>1.0 表示有加成）
func (g *Game) notifyImmortalBossKill(p *player.Player, instanceID string) float64 {
	m := g.LuckyImmortalBoss
	m.mu.Lock()

	if m.activeSession == nil || m.activeSession.settled {
		m.mu.Unlock()
		return 1.0
	}
	session := m.activeSession
	if session.instanceID != instanceID {
		m.mu.Unlock()
		return 1.0
	}

	killMult := session.currentMult
	session.killCount++
	session.livesLeft--

	log.Printf("[LuckyImmortalBoss] 擊破！玩家=%s 第%d次 倍率=×%.1f 剩餘條命=%d",
		p.DisplayName, session.killCount, killMult, session.livesLeft)

	if session.livesLeft <= 0 {
		// 5 條命全部耗盡 → 永生終結
		session.settled = true
		killCount := session.killCount
		m.mu.Unlock()
		g.doImmortalBossEnd(p, killCount)
		return killMult
	}

	// 復活：倍率提升 +0.5x，更新 instanceID
	session.currentMult += ImmortalBossMultStep
	session.instanceID = fmt.Sprintf("immortal_boss_%d_r%d", time.Now().UnixNano(), session.killCount)
	killCount := session.killCount
	newMult := session.currentMult
	livesLeft := session.livesLeft
	newInstanceID := session.instanceID
	x := session.x
	y := session.y
	m.mu.Unlock()

	// 廣播擊破通知
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyImmortalBoss,
		Payload: ws.LuckyImmortalBossPayload{
			Event:             "immortal_kill",
			TriggerPlayerID:   p.ID,
			TriggerPlayerName: p.DisplayName,
			KillCount:         killCount,
			KillMult:          killMult,
			NextMult:          newMult,
			LivesLeft:         livesLeft,
			X:                 x,
			Y:                 y,
		},
	})

	// 廣播復活通知（帶新 instanceID，讓 Client 知道要追蹤新的目標）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyImmortalBoss,
		Payload: ws.LuckyImmortalBossPayload{
			Event:       "immortal_revive",
			LivesLeft:   livesLeft,
			CurrentMult: newMult,
			X:           x,
			Y:           y,
		},
	})

	log.Printf("[LuckyImmortalBoss] 復活！新 instanceID=%s 倍率=×%.1f 剩餘條命=%d",
		newInstanceID, newMult, livesLeft)

	return killMult
}

// doImmortalBossEnd 永生終結（5 條命全部耗盡）
func (g *Game) doImmortalBossEnd(p *player.Player, killCount int) {
	m := g.LuckyImmortalBoss
	m.mu.Lock()
	m.endBoost = &immortalBossEndBoost{
		mult:      ImmortalBossEndMult,
		expiresAt: time.Now().Add(ImmortalBossEndDuration),
	}
	m.mu.Unlock()

	log.Printf("[LuckyImmortalBoss] 永生終結！玩家=%s 全服×%.1f 持續%ds",
		p.DisplayName, ImmortalBossEndMult, int(ImmortalBossEndDuration.Seconds()))

	// 全服廣播永生終結
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyImmortalBoss,
		Payload: ws.LuckyImmortalBossPayload{
			Event:             "immortal_end",
			TriggerPlayerID:   p.ID,
			TriggerPlayerName: p.DisplayName,
			GlobalMult:        ImmortalBossEndMult,
			GlobalDurationSec: int(ImmortalBossEndDuration.Seconds()),
			KillCount:         killCount,
		},
	})

	// 全服公告（最高優先）
	ann := g.Announce.Create(announce.EventLuckyImmortalBoss, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("💀 永生終結！%s 擊敗了永生 BOSS %d 次！全服 ×%.1f 加成 %d 秒！",
			p.DisplayName, killCount, ImmortalBossEndMult, int(ImmortalBossEndDuration.Seconds())),
		"color": "#8B0000",
	})
	g.broadcastAnnouncement(ann)

	// 加成結束後廣播
	go func() {
		time.Sleep(ImmortalBossEndDuration)
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyImmortalBoss,
			Payload: ws.LuckyImmortalBossPayload{
				Event: "immortal_end_expire",
			},
		})
	}()
}

// runImmortalBossTimeout 永生 BOSS 超時計時器
func (g *Game) runImmortalBossTimeout(session *immortalBossSession) {
	time.Sleep(ImmortalBossTimeout)

	m := g.LuckyImmortalBoss
	m.mu.Lock()
	if m.activeSession != session || session.settled {
		m.mu.Unlock()
		return
	}
	session.settled = true
	killCount := session.killCount
	livesLeft := session.livesLeft
	m.mu.Unlock()

	log.Printf("[LuckyImmortalBoss] 超時！擊破次數=%d 剩餘條命=%d", killCount, livesLeft)

	// 廣播超時
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyImmortalBoss,
		Payload: ws.LuckyImmortalBossPayload{
			Event:     "immortal_timeout",
			KillCount: killCount,
			LivesLeft: livesLeft,
		},
	})

	ann := g.Announce.Create(announce.EventLuckyImmortalBoss, "", 0, map[string]string{
		"message": fmt.Sprintf("⏰ 永生 BOSS 消散了！共被擊破 %d 次，剩餘 %d 條命。", killCount, livesLeft),
		"color":   "#666666",
	})
	g.broadcastAnnouncement(ann)
}
