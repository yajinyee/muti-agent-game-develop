// lucky_luck_totem_handler.go — 幸運幸運圖騰魚系統（DAY-275）
// 業界依據：Fish It Luck Totem（2026）「全場幸運加成」機制
//   進化版：「全服共享圖騰加成 + 觸發者個人雙倍加成」讓全服一起爽
//
// 設計：
//   擊破 T233 後，場上出現「幸運圖騰」（持續 15 秒）
//   圖騰期間：
//     全服所有玩家每次擊破任何目標 → ×1.3 全服加成
//     觸發玩家額外獲得 ×1.5 個人加成（疊加在全服加成上）
//   15 秒後圖騰消失，廣播結算（總擊破數/總獎勵）
//   個人冷卻 30 秒；全服冷卻 50 秒
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
	LuckyLuckTotemGlobalMult   = 1.3  // 全服加成倍率
	LuckyLuckTotemPersonalMult = 1.5  // 觸發者個人額外加成
	LuckyLuckTotemDuration     = 15   // 圖騰持續秒數
)

// luckTotemSession 幸運圖騰 session（全服共享）
type luckTotemSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	expiresAt         time.Time
	totalKills        int
	totalReward       int
	settled           bool
}

// luckyLuckTotemManager 幸運幸運圖騰魚系統管理器
type luckyLuckTotemManager struct {
	mu                sync.Mutex
	personalCooldowns map[string]time.Time
	globalCooldown    time.Time
	activeSession     *luckTotemSession
}

func newLuckyLuckTotemManager() *luckyLuckTotemManager {
	return &luckyLuckTotemManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyLuckTotemFish 判斷是否為幸運圖騰魚
func isLuckyLuckTotemFish(defID string) bool {
	return defID == "T233"
}

// isLuckTotemActive 判斷圖騰是否啟動中（供 handleKill 使用）
func (m *luckyLuckTotemManager) isLuckTotemActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil || m.activeSession.settled {
		return false
	}
	if time.Now().After(m.activeSession.expiresAt) {
		return false
	}
	return true
}

// getLuckTotemMult 取得圖騰倍率加成（供 handleKill 使用）
// 回傳 (globalMult, personalMult)
func (m *luckyLuckTotemManager) getLuckTotemMult(playerID string) (float64, float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil || m.activeSession.settled || time.Now().After(m.activeSession.expiresAt) {
		return 1.0, 1.0
	}
	globalMult := LuckyLuckTotemGlobalMult
	personalMult := 1.0
	if playerID == m.activeSession.triggerPlayerID {
		personalMult = LuckyLuckTotemPersonalMult
	}
	return globalMult, personalMult
}

// recordLuckTotemKill 記錄圖騰期間的擊破（供 handleKill 使用）
func (m *luckyLuckTotemManager) recordLuckTotemKill(reward int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil || m.activeSession.settled || time.Now().After(m.activeSession.expiresAt) {
		return
	}
	m.activeSession.totalKills++
	m.activeSession.totalReward += reward
}

// tryLuckyLuckTotemFish 擊破 T233 後觸發幸運圖騰
func (g *Game) tryLuckyLuckTotemFish(p *player.Player) {
	m := g.LuckyLuckTotem
	now := time.Now()

	m.mu.Lock()
	// 個人冷卻（30 秒）
	if cd, ok := m.personalCooldowns[p.ID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	// 全服冷卻（50 秒）
	if now.Before(m.globalCooldown) {
		m.mu.Unlock()
		return
	}
	// 已有進行中 session
	if m.activeSession != nil && !m.activeSession.settled && now.Before(m.activeSession.expiresAt) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(30 * time.Second)
	m.globalCooldown = now.Add(50 * time.Second)

	// 建立 session
	sess := &luckTotemSession{
		triggerPlayerID:   p.ID,
		triggerPlayerName: p.DisplayName,
		expiresAt:         now.Add(LuckyLuckTotemDuration * time.Second),
		totalKills:        0,
		totalReward:       0,
		settled:           false,
	}
	m.activeSession = sess
	m.mu.Unlock()

	log.Printf("[LuckTotem] player=%s triggered! global=×%.1f personal=×%.1f duration=%ds",
		p.ID, LuckyLuckTotemGlobalMult, LuckyLuckTotemPersonalMult, LuckyLuckTotemDuration)

	// 發送個人通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyLuckTotem,
		Payload: ws.LuckyLuckTotemPayload{
			Event:         "totem_start",
			PlayerID:      p.ID,
			PlayerName:    p.DisplayName,
			GlobalMult:    LuckyLuckTotemGlobalMult,
			PersonalMult:  LuckyLuckTotemPersonalMult,
			DurationSec:   LuckyLuckTotemDuration,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyLuckTotem,
		Payload: ws.LuckyLuckTotemPayload{
			Event:        "totem_broadcast",
			PlayerID:     p.ID,
			PlayerName:   p.DisplayName,
			GlobalMult:   LuckyLuckTotemGlobalMult,
			DurationSec:  LuckyLuckTotemDuration,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyLuckTotem, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🍀 %s 觸發幸運圖騰！全服 ×%.1f 加成 %d 秒！",
			p.DisplayName, LuckyLuckTotemGlobalMult, LuckyLuckTotemDuration),
	})

	// 啟動超時 goroutine
	go g.runLuckTotemTimeout(sess)
}

// notifyLuckTotemKill 圖騰期間擊破通知（供 handleKill 使用）
func (g *Game) notifyLuckTotemKill(p *player.Player, bonusReward int) {
	if bonusReward <= 0 {
		return
	}
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyLuckTotem,
		Payload: ws.LuckyLuckTotemPayload{
			Event:       "totem_kill",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			BonusReward: bonusReward,
		},
	})
}

// runLuckTotemTimeout 15 秒後結算圖騰
func (g *Game) runLuckTotemTimeout(sess *luckTotemSession) {
	remaining := time.Until(sess.expiresAt)
	if remaining > 0 {
		time.Sleep(remaining)
	}

	m := g.LuckyLuckTotem
	m.mu.Lock()
	if m.activeSession != sess || sess.settled {
		m.mu.Unlock()
		return
	}
	sess.settled = true
	totalKills := sess.totalKills
	totalReward := sess.totalReward
	m.mu.Unlock()

	log.Printf("[LuckTotem] expired! trigger=%s kills=%d reward=%d",
		sess.triggerPlayerID, totalKills, totalReward)

	// 全服廣播結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyLuckTotem,
		Payload: ws.LuckyLuckTotemPayload{
			Event:       "totem_end",
			PlayerID:    sess.triggerPlayerID,
			PlayerName:  sess.triggerPlayerName,
			TotalKills:  totalKills,
			TotalReward: totalReward,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyLuckTotem, sess.triggerPlayerName, totalReward, map[string]string{
		"message": fmt.Sprintf("🍀 幸運圖騰結束！全服共擊破 %d 個目標，總獎勵 +%d！",
			totalKills, totalReward),
		"color": "#00FF88",
	})
}
