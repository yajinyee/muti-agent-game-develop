// vampire_fish_handler.go — 吸血鬼魚累積倍率系統 handler（DAY-182）
// 業界依據：JILI 2026「The explicit multiplier of vampires increases the more you fight,
// and there is a chance that you can enter the multiplier mode, up to X5」
// 擊破 T140 後觸發「吸血鬼模式」：
//   1. 玩家進入吸血鬼模式，初始倍率 1.0x
//   2. 每擊破一個目標，倍率累積 +0.1x（最高 5.0x）
//   3. 持續 15 秒，時間到或達到 5.0x 後廣播結果
//   4. 全服廣播「吸血鬼模式激活」，讓其他玩家看到
// 設計差異：
//   - 與幸運星魚（固定 ×2，10秒）不同，吸血鬼魚是「累積型倍率」（越打越高），
//     製造「越打越爽」的正向反饋；玩家需要在 15 秒內盡量多打目標
//   - 與黃金鯊魚（全服固定 ×1.5）不同，吸血鬼魚是「個人累積」，讓玩家感受到「自己的努力有回報」
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// VampireFishDurationSec 吸血鬼模式持續時間（秒）
	VampireFishDurationSec = 15
	// VampireFishCooldownSec 個人冷卻時間（秒）
	VampireFishCooldownSec = 35
	// VampireFishInitMult 初始倍率
	VampireFishInitMult = 1.0
	// VampireFishMultStep 每次擊破倍率增加量
	VampireFishMultStep = 0.1
	// VampireFishMaxMult 最高倍率
	VampireFishMaxMult = 5.0
	// VampireFishAnnounceThreshold 全服公告倍率門檻（達到 3.0x 才公告）
	VampireFishAnnounceThreshold = 3.0
)

// vampireSession 吸血鬼模式 session（per-player）
type vampireSession struct {
	active    bool
	expiresAt time.Time
	currentMult float64
	killCount   int
}

// vampireFishManager 吸血鬼魚累積倍率管理器（per-player）
type vampireFishManager struct {
	mu       sync.Mutex
	sessions map[string]*vampireSession // playerID → session
	cooldown map[string]time.Time       // playerID → 冷卻結束時間
}

// newVampireFishManager 建立吸血鬼魚管理器
func newVampireFishManager() *vampireFishManager {
	return &vampireFishManager{
		sessions: make(map[string]*vampireSession),
		cooldown: make(map[string]time.Time),
	}
}

// isVampireFish 判斷是否為吸血鬼魚（T140）
func isVampireFish(defID string) bool {
	return defID == "T140"
}

// isOnCooldown 檢查玩家是否在冷卻中
func (m *vampireFishManager) isOnCooldown(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	cd, ok := m.cooldown[playerID]
	if !ok {
		return false
	}
	return time.Now().Before(cd)
}

// activate 激活吸血鬼模式
func (m *vampireFishManager) activate(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[playerID] = &vampireSession{
		active:      true,
		expiresAt:   time.Now().Add(time.Duration(VampireFishDurationSec) * time.Second),
		currentMult: VampireFishInitMult,
		killCount:   0,
	}
	m.cooldown[playerID] = time.Now().Add(time.Duration(VampireFishCooldownSec) * time.Second)
}

// recordKill 記錄擊破，累積倍率，回傳新倍率
func (m *vampireFishManager) recordKill(playerID string) (float64, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok || !sess.active {
		return 1.0, false
	}
	if time.Now().After(sess.expiresAt) {
		sess.active = false
		return 1.0, false
	}
	// 累積倍率
	sess.currentMult += VampireFishMultStep
	if sess.currentMult > VampireFishMaxMult {
		sess.currentMult = VampireFishMaxMult
	}
	sess.killCount++
	return sess.currentMult, true
}

// getVampireMult 取得當前吸血鬼倍率（供 handleKill 使用）
func (m *vampireFishManager) getVampireMult(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok || !sess.active {
		return 1.0
	}
	if time.Now().After(sess.expiresAt) {
		sess.active = false
		return 1.0
	}
	return sess.currentMult
}

// deactivate 結束吸血鬼模式，回傳最終倍率和擊破數
func (m *vampireFishManager) deactivate(playerID string) (float64, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok {
		return 1.0, 0
	}
	finalMult := sess.currentMult
	killCount := sess.killCount
	sess.active = false
	delete(m.sessions, playerID)
	return finalMult, killCount
}

// tryVampireFish 擊破 T140 後觸發吸血鬼模式（DAY-182）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryVampireFish(p *player.Player) {
	// 個人冷卻檢查
	if g.VampireFish.isOnCooldown(p.ID) {
		return
	}

	// 激活吸血鬼模式
	g.VampireFish.activate(p.ID)

	log.Printf("[VampireFish] player=%s activated vampire mode", p.ID)

	// 廣播吸血鬼模式激活（個人）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgVampireFish,
		Payload: ws.VampireFishPayload{
			Phase:       "vampire_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			CurrentMult: VampireFishInitMult,
			MaxMult:     VampireFishMaxMult,
			DurationSec: VampireFishDurationSec,
		},
	})

	// 全服廣播（讓其他玩家看到）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgVampireFish,
		Payload: ws.VampireFishPayload{
			Phase:      "vampire_broadcast",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
		},
	})

	// 等待吸血鬼模式結束
	time.Sleep(time.Duration(VampireFishDurationSec) * time.Second)

	// 結束吸血鬼模式
	finalMult, killCount := g.VampireFish.deactivate(p.ID)

	// 廣播吸血鬼模式結束（個人）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgVampireFish,
		Payload: ws.VampireFishPayload{
			Phase:       "vampire_end",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			CurrentMult: finalMult,
			KillCount:   killCount,
		},
	})

	// 全服公告（達到 3.0x 才公告）
	if finalMult >= VampireFishAnnounceThreshold {
		g.announceVampireFish(p.DisplayName, finalMult, killCount)
	}

	log.Printf("[VampireFish] player=%s ended: finalMult=%.1f kills=%d",
		p.ID, finalMult, killCount)
}

// getVampireMult 取得吸血鬼倍率（供 handleKill 使用）
func (g *Game) getVampireMult(playerID string) float64 {
	return g.VampireFish.getVampireMult(playerID)
}

// recordVampireKill 記錄吸血鬼模式擊破，累積倍率（供 handleKill 使用）
// 回傳新倍率（如果在吸血鬼模式中）
func (g *Game) recordVampireKill(p *player.Player) float64 {
	newMult, active := g.VampireFish.recordKill(p.ID)
	if !active {
		return 1.0
	}
	// 廣播倍率更新（個人）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgVampireFish,
		Payload: ws.VampireFishPayload{
			Phase:       "mult_update",
			PlayerID:    p.ID,
			CurrentMult: newMult,
		},
	})
	return newMult
}

// announceVampireFish 全服公告吸血鬼魚（DAY-182）
func (g *Game) announceVampireFish(playerName string, finalMult float64, killCount int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "vampire_fish",
			"message":    fmt.Sprintf("🧛 %s 吸血鬼模式達到 %.1fx！擊破 %d 個目標！", playerName, finalMult, killCount),
			"color":      "#8B0000", // 深紅色（吸血鬼感）
			"duration":   5.0,
			"priority":   3,
		},
	})
}
