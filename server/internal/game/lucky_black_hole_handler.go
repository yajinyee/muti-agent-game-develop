// lucky_black_hole_handler.go — 幸運黑洞魚系統（DAY-221）
// 業界原創「重力黑洞」機制
// 設計：擊破 T179 後在場上建立「重力黑洞」（持續 10 秒）：
//   - 黑洞建立在場景中央附近（隨機偏移），半徑 350px
//   - 黑洞範圍內所有目標每 1 秒被「吸引」（HP -10%，模擬重力傷害）
//   - 10 秒後「奇點爆炸」：黑洞範圍內所有目標 85% 擊破機率（0.70x 倍率，全服共享）
//   - 個人冷卻 22 秒；全服冷卻 35 秒
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyBlackHolePersonalCD = 22 * time.Second // 個人冷卻
	LuckyBlackHoleGlobalCD   = 35 * time.Second // 全服冷卻

	BlackHoleActiveDuration  = 10 * time.Second       // 黑洞持續時間
	BlackHoleGravityInterval = 1 * time.Second        // 重力傷害間隔
	BlackHoleGravityRadius   = 350.0                  // 黑洞半徑（px）
	BlackHoleGravityDmg      = 0.10                   // 每秒 HP -10%
	BlackHoleSingularityKill = 0.85                   // 奇點爆炸擊破機率
	BlackHoleSingularityMult = 0.70                   // 奇點爆炸倍率
	BlackHoleKillMult        = 2.0                    // 黑洞範圍內擊破倍率加成
)

// blackHoleSession 黑洞會話
type blackHoleSession struct {
	playerID   string
	playerName string
	centerX    float64
	centerY    float64
	expiresAt  time.Time
}

// luckyBlackHoleManager 幸運黑洞魚管理器
type luckyBlackHoleManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前活躍黑洞
	activeSession *blackHoleSession
}

func newLuckyBlackHoleManager() *luckyBlackHoleManager {
	return &luckyBlackHoleManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyBlackHoleFish 判斷是否為幸運黑洞魚
func isLuckyBlackHoleFish(defID string) bool {
	return defID == "T179"
}

// getLuckyBlackHoleMultiplier 取得黑洞範圍內的倍率加成（供 handleKill 使用）
// 若目標在黑洞範圍內，回傳 BlackHoleKillMult（×2.0），否則回傳 1.0
func (g *Game) getLuckyBlackHoleMultiplier(targetX, targetY float64) float64 {
	mgr := g.LuckyBlackHole
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if mgr.activeSession == nil || time.Now().After(mgr.activeSession.expiresAt) {
		return 1.0
	}

	dx := targetX - mgr.activeSession.centerX
	dy := targetY - mgr.activeSession.centerY
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist <= BlackHoleGravityRadius {
		return BlackHoleKillMult
	}
	return 1.0
}

// tryLuckyBlackHoleFish 擊破 T179 後觸發黑洞（供 handleKill 使用）
func (g *Game) tryLuckyBlackHoleFish(p *player.Player) {
	mgr := g.LuckyBlackHole
	mgr.mu.Lock()

	now := time.Now()

	// 全服冷卻檢查
	if now.Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && now.Before(cd) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = now.Add(LuckyBlackHolePersonalCD)
	mgr.globalCooldownUntil = now.Add(LuckyBlackHoleGlobalCD)

	// 黑洞位置（場景中央附近）
	centerX := 400.0 + (rand.Float64()-0.5)*200.0
	centerY := 300.0 + (rand.Float64()-0.5)*150.0

	session := &blackHoleSession{
		playerID:   p.ID,
		playerName: p.DisplayName,
		centerX:    centerX,
		centerY:    centerY,
		expiresAt:  now.Add(BlackHoleActiveDuration),
	}
	mgr.activeSession = session
	mgr.mu.Unlock()

	log.Printf("[BlackHole] player=%s center=(%.0f,%.0f)", p.ID, centerX, centerY)

	// 全服廣播：黑洞建立
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyBlackHole,
		Payload: ws.LuckyBlackHolePayload{
			Event:       "singularity_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			X:           centerX,
			Y:           centerY,
			Radius:      BlackHoleGravityRadius,
			DurationSec: int(BlackHoleActiveDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyBlackHole, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌑 %s 召喚重力黑洞！範圍內目標正在被吸引！", p.DisplayName),
		"color":   "#1a1a2e",
	})
	g.broadcastAnnouncement(ann)

	// 啟動黑洞 goroutine
	go g.runBlackHoleGravity(p, session)
}

// runBlackHoleGravity 執行黑洞重力傷害（持續 10 秒後奇點爆炸）
func (g *Game) runBlackHoleGravity(p *player.Player, session *blackHoleSession) {
	ticker := time.NewTicker(BlackHoleGravityInterval)
	defer ticker.Stop()

	deadline := session.expiresAt

	for {
		select {
		case <-ticker.C:
			if time.Now().After(deadline) {
				// 奇點爆炸
				g.doBlackHoleSingularity(p, session)
				// 清除活躍黑洞
				g.LuckyBlackHole.mu.Lock()
				if g.LuckyBlackHole.activeSession == session {
					g.LuckyBlackHole.activeSession = nil
				}
				g.LuckyBlackHole.mu.Unlock()
				return
			}

			// 重力傷害
			g.applyBlackHoleGravityDamage(session.centerX, session.centerY)

		case <-g.stopCh:
			return
		}
	}
}

// applyBlackHoleGravityDamage 對黑洞範圍內目標造成重力傷害
func (g *Game) applyBlackHoleGravityDamage(cx, cy float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, t := range g.Targets {
		if t == nil || t.HP <= 0 {
			continue
		}
		dx := t.X - cx
		dy := t.Y - cy
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > BlackHoleGravityRadius {
			continue
		}

		dmg := int(float64(t.HP) * BlackHoleGravityDmg)
		if dmg < 1 {
			dmg = 1
		}
		t.HP -= dmg
		if t.HP < 0 {
			t.HP = 0
		}

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetUpdate,
			Payload: ws.TargetUpdatePayload{
				InstanceID: t.InstanceID,
				HP:         t.HP,
				MaxHP:      t.MaxHP,
			},
		})
	}
}

// doBlackHoleSingularity 奇點爆炸：黑洞範圍內所有目標 85% 擊破機率
func (g *Game) doBlackHoleSingularity(p *player.Player, session *blackHoleSession) {
	g.mu.Lock()

	hitEntries := make([]string, 0)
	totalReward := 0

	for id, t := range g.Targets {
		if t == nil || t.HP <= 0 {
			continue
		}
		dx := t.X - session.centerX
		dy := t.Y - session.centerY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > BlackHoleGravityRadius {
			continue
		}

		if rand.Float64() < BlackHoleSingularityKill {
			reward := int(float64(p.BetLevel) * t.Multiplier * BlackHoleSingularityMult)
			if reward < 1 {
				reward = 1
			}
			totalReward += reward
			p.Coins += reward
			if p.Coins > p.MaxCoins {
				p.MaxCoins = p.Coins
			}
			hitEntries = append(hitEntries, t.InstanceID)
			t.HP = 0
			delete(g.Targets, id)

			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgTargetKill,
				Payload: ws.TargetKillPayload{
					InstanceID: t.InstanceID,
					KillerID:   p.ID,
					Reward:     reward,
					Multiplier: t.Multiplier,
				},
			})
		}
	}
	g.mu.Unlock()

	log.Printf("[BlackHole] singularity player=%s killed=%d reward=%d",
		p.ID, len(hitEntries), totalReward)

	// 廣播奇點爆炸結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyBlackHole,
		Payload: ws.LuckyBlackHolePayload{
			Event:       "singularity_result",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			KilledCount: len(hitEntries),
			TotalReward: totalReward,
		},
	})

	if len(hitEntries) >= 3 {
		ann := g.Announce.Create(announce.EventLuckyBlackHole, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🌑 %s 黑洞奇點爆炸！擊破 %d 個目標，獲得 %d 金幣！",
				p.DisplayName, len(hitEntries), totalReward),
			"color": "#9B59B6",
		})
		g.broadcastAnnouncement(ann)
	}
}
