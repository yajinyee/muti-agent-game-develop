// rainbow_phoenix_handler.go — 彩虹鳳凰 Power Up 系統 handler（DAY-151）
// 業界依據：royal-fishing.co.uk 2026「Multicoloured phoenix (blue, pink, purple, orange) with magical aura.
// Awaken Boss with 30x basic multiplier. Power Up attack delivers 6x-10x boost for rewards up to 300 times bet.」
// 設計：T115 彩虹鳳凰擊破後觸發「Power Up 模式」
// 玩家在 8 秒內所有攻擊獲得隨機 6x-10x 倍率加成，最高 300x
// 全服廣播：讓其他玩家看到「有人觸發了彩虹鳳凰 Power Up」
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
	RainbowPhoenixDefID       = "T115"
	RainbowPhoenixDuration    = 8 * time.Second   // Power Up 持續時間
	RainbowPhoenixCooldown    = 90 * time.Second  // 每個玩家的冷卻時間
	RainbowPhoenixMinMult     = 6.0               // 最低 Power Up 倍率
	RainbowPhoenixMaxMult     = 10.0              // 最高 Power Up 倍率
	RainbowPhoenixMaxReward   = 300               // 最高獎勵倍率（300x betLevel）
)

// rainbowPhoenixSession Power Up session
type rainbowPhoenixSession struct {
	PlayerID    string
	StartAt     time.Time
	EndAt       time.Time
	PowerUpMult float64   // 本次 Power Up 倍率（6x-10x）
	KillCount   int
	TotalReward int
}

// rainbowPhoenixManager 管理所有玩家的 Power Up session
type rainbowPhoenixManager struct {
	mu        sync.Mutex
	sessions  map[string]*rainbowPhoenixSession // playerID → session
	cooldowns map[string]time.Time              // playerID → 冷卻結束時間
}

func newRainbowPhoenixManager() *rainbowPhoenixManager {
	return &rainbowPhoenixManager{
		sessions:  make(map[string]*rainbowPhoenixSession),
		cooldowns: make(map[string]time.Time),
	}
}

// CanTrigger 判斷玩家是否可以觸發 Power Up
func (m *rainbowPhoenixManager) CanTrigger(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, active := m.sessions[playerID]; active {
		return false
	}
	if cd, ok := m.cooldowns[playerID]; ok && time.Now().Before(cd) {
		return false
	}
	return true
}

// StartSession 開始 Power Up session，隨機決定倍率
func (m *rainbowPhoenixManager) StartSession(playerID string) *rainbowPhoenixSession {
	m.mu.Lock()
	defer m.mu.Unlock()
	// 隨機決定 Power Up 倍率（6x-10x，整數）
	mult := float64(6 + rand.Intn(5)) // 6, 7, 8, 9, 10
	sess := &rainbowPhoenixSession{
		PlayerID:    playerID,
		StartAt:     time.Now(),
		EndAt:       time.Now().Add(RainbowPhoenixDuration),
		PowerUpMult: mult,
	}
	m.sessions[playerID] = sess
	return sess
}

// GetActiveMult 取得玩家當前 Power Up 倍率（0 = 無 Power Up）
func (m *rainbowPhoenixManager) GetActiveMult(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok {
		return 0.0
	}
	if time.Now().After(sess.EndAt) {
		return 0.0
	}
	return sess.PowerUpMult
}

// RecordKill 記錄 Power Up 期間的擊破
func (m *rainbowPhoenixManager) RecordKill(playerID string, reward int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok {
		return
	}
	sess.KillCount++
	sess.TotalReward += reward
}

// CheckExpiry 檢查是否過期，過期則結束 session
func (m *rainbowPhoenixManager) CheckExpiry(playerID string) *rainbowPhoenixSession {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok {
		return nil
	}
	if time.Now().After(sess.EndAt) {
		delete(m.sessions, playerID)
		m.cooldowns[playerID] = time.Now().Add(RainbowPhoenixCooldown)
		return sess
	}
	return nil
}

// IsActive 判斷玩家是否在 Power Up 模式中
func (m *rainbowPhoenixManager) IsActive(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok {
		return false
	}
	return time.Now().Before(sess.EndAt)
}

// RemovePlayer 玩家離線時清理
func (m *rainbowPhoenixManager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, playerID)
}

// isRainbowPhoenix 判斷是否為彩虹鳳凰目標
func isRainbowPhoenix(defID string) bool {
	return defID == RainbowPhoenixDefID
}

// getRainbowPhoenixMult 取得彩虹鳳凰 Power Up 倍率（供 handleKill 使用）
func (g *Game) getRainbowPhoenixMult(playerID string) float64 {
	if g.RainbowPhoenix == nil {
		return 1.0
	}
	mult := g.RainbowPhoenix.GetActiveMult(playerID)
	if mult <= 0 {
		return 1.0
	}
	return mult
}

// tryRainbowPhoenix 擊破彩虹鳳凰後觸發 Power Up 模式（由 handleKill 呼叫）
func (g *Game) tryRainbowPhoenix(p *player.Player, killedInstanceID string, killedX, killedY float64) {
	if !g.RainbowPhoenix.CanTrigger(p.ID) {
		return
	}

	sess := g.RainbowPhoenix.StartSession(p.ID)
	log.Printf("[RainbowPhoenix] player=%s triggered Power Up: mult=%.0fx, duration=8s", p.ID, sess.PowerUpMult)

	// 廣播 Power Up 開始（全服可見）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRainbowPhoenixActivate,
		Payload: ws.RainbowPhoenixActivatePayload{
			TriggerID:   killedInstanceID,
			TriggerX:    killedX,
			TriggerY:    killedY,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			PowerUpMult: sess.PowerUpMult,
			Duration:    int(RainbowPhoenixDuration.Seconds()),
			Message:     fmt.Sprintf("🌈 %s 觸發彩虹鳳凰！Power Up %.0fx！持續 8 秒！", p.DisplayName, sess.PowerUpMult),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventRainbowPhoenix, p.DisplayName, int(sess.PowerUpMult), nil)
	g.broadcastAnnouncement(ann)

	// 8 秒後結束 Power Up
	go func() {
		time.Sleep(RainbowPhoenixDuration)
		g.endRainbowPhoenixSession(p)
	}()
}

// notifyRainbowPhoenixKill 在 Power Up 期間擊破目標時呼叫（由 handleKill 呼叫）
// 回傳 Power Up 倍率（1.0 = 無 Power Up）
func (g *Game) notifyRainbowPhoenixKill(p *player.Player, baseReward int) float64 {
	if g.RainbowPhoenix == nil {
		return 1.0
	}
	mult := g.RainbowPhoenix.GetActiveMult(p.ID)
	if mult <= 0 {
		return 1.0
	}
	// 記錄擊破（Power Up 倍率已在 handleKill 中套用）
	g.RainbowPhoenix.RecordKill(p.ID, int(float64(baseReward)*mult))
	return mult
}

// endRainbowPhoenixSession 結束 Power Up session
func (g *Game) endRainbowPhoenixSession(p *player.Player) {
	sess := g.RainbowPhoenix.CheckExpiry(p.ID)
	if sess == nil {
		// 可能已被提前清理
		return
	}

	log.Printf("[RainbowPhoenix] player=%s Power Up ended: kills=%d reward=%d",
		p.ID, sess.KillCount, sess.TotalReward)

	// 廣播 Power Up 結束（全服可見）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRainbowPhoenixEnd,
		Payload: ws.RainbowPhoenixEndPayload{
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			PowerUpMult: sess.PowerUpMult,
			TotalKills:  sess.KillCount,
			TotalReward: sess.TotalReward,
			NewBalance:  p.GetCoins(),
			Message:     fmt.Sprintf("🌈 %s 的彩虹鳳凰 Power Up 結束！擊破 %d 個目標，獲得 %d 金幣！", p.DisplayName, sess.KillCount, sess.TotalReward),
		},
	})

	// 更新玩家狀態
	g.sendPlayerUpdate(p)

	// 全服公告：擊破 ≥3 個目標時廣播
	if sess.KillCount >= 3 {
		ann := g.Announce.Create(announce.EventRainbowPhoenixResult, p.DisplayName, sess.TotalReward, map[string]string{
			"kills": fmt.Sprintf("%d", sess.KillCount),
			"mult":  fmt.Sprintf("%.0f", sess.PowerUpMult),
		})
		g.broadcastAnnouncement(ann)
	}
}

// sendRainbowPhoenixStatus 登入時發送 Power Up 狀態
func (g *Game) sendRainbowPhoenixStatus(p *player.Player) {
	if g.RainbowPhoenix == nil {
		return
	}
	mult := g.RainbowPhoenix.GetActiveMult(p.ID)
	isActive := mult > 0
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgRainbowPhoenixStatus,
		Payload: ws.RainbowPhoenixStatusPayload{
			IsActive:    isActive,
			PowerUpMult: mult,
		},
	}); err != nil {
		log.Printf("[RainbowPhoenix] send status error: %v", err)
	}
}
