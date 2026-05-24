// lucky_wrath_charge_handler.go — 幸運怒氣蓄積魚系統（DAY-290）
// 業界依據：Royal Fishing Jili「Dragon Wrath system accumulates with every shot fired.
//          Once the wrath meter fills, players unleash a massive meteorite attack across the centre screen,
//          simultaneously targeting multiple fish including Immortal Bosses and the ChainLong King.」
//          業界原創「怒氣蓄積 + 全場隕石雨 + 蓄積越多爆發越強」機制
//
// 設計：
//   - 擊破 T248 後，觸發玩家進入「怒氣蓄積模式」（持續 25 秒）
//   - 模式期間，玩家每次擊破任何目標 → 怒氣值 +1（最高 20 點）
//   - 25 秒後（或怒氣值達到 20）→ 自動爆發「怒氣隕石雨」
//   - 隕石數量 = 怒氣值（最少 3 顆，最多 20 顆）
//   - 每顆隕石 HP -50%，AOE r=100px
//   - 怒氣值 ≥ 15 → 「完美怒氣」：全服 ×2.8 加成 7 秒
//   - 個人冷卻 28 秒；全服冷卻 45 秒
//
// 設計差異：
//   - 與龍怒隕石（T242，固定 4-7 顆隕石）不同，怒氣蓄積是「打越多魚，隕石越多」
//   - 「怒氣值 = 擊破次數」讓玩家有「要趁 25 秒內瘋狂打魚」的動力
//   - 「最多 20 顆隕石」讓玩家有「要把怒氣打滿才能拿到最多隕石」的目標感
//   - 「完美怒氣（≥15）→ 全服 ×2.8」讓玩家有「要打到 15 次才能觸發完美」的動力
//   - 「全服廣播怒氣進度」讓所有玩家看到「有人在蓄積怒氣，快去打魚幫他」的社交感
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
	LuckyWrathChargePersonalCD = 28 * time.Second // 個人冷卻
	LuckyWrathChargeGlobalCD   = 45 * time.Second // 全服冷卻

	// 怒氣蓄積設計
	WrathChargeTimeout  = 25 * time.Second // 蓄積時間
	WrathChargeMax      = 20               // 最高怒氣值
	WrathChargeMin      = 3                // 最少隕石數（保底）
	WrathMeteorHPDmg    = 0.50             // 每顆隕石 HP -50%
	WrathMeteorRadius   = 100.0            // AOE 半徑（px）
	WrathMeteorInterval = 300 * time.Millisecond // 每顆隕石間隔

	// 完美怒氣：全服加成
	WrathPerfectThreshold = 15              // 完美怒氣門檻
	WrathPerfectMult      = 2.8             // 全服 ×2.8
	WrathPerfectDuration  = 7 * time.Second // 持續 7 秒
)

// wrathPerfectBoost 完美怒氣全服加成
type wrathPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

// wrathChargeSession 怒氣蓄積會話
type wrathChargeSession struct {
	playerID   string
	playerName string
	wrathValue int       // 當前怒氣值
	expiresAt  time.Time
	settled    bool
}

// luckyWrathChargeManager 幸運怒氣蓄積魚管理器
type luckyWrathChargeManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 活躍的怒氣蓄積會話（playerID → session）
	activeSessions map[string]*wrathChargeSession

	// 完美怒氣全服加成
	perfectBoost *wrathPerfectBoost
}

func newLuckyWrathChargeManager() *luckyWrathChargeManager {
	return &luckyWrathChargeManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*wrathChargeSession),
	}
}

// isLuckyWrathChargeFish 判斷是否為幸運怒氣蓄積魚
func isLuckyWrathChargeFish(defID string) bool {
	return defID == "T248"
}

// getWrathPerfectMult 取得完美怒氣全服加成倍率（供 handleKill 使用）
func (m *luckyWrathChargeManager) getWrathPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

// addWrathCharge 玩家擊破目標時增加怒氣值（由 handleKill 呼叫）
func (m *luckyWrathChargeManager) addWrathCharge(playerID string) (bool, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	if !ok || sess.settled {
		return false, 0
	}
	if sess.wrathValue >= WrathChargeMax {
		return false, sess.wrathValue
	}
	sess.wrathValue++
	return true, sess.wrathValue
}

// isWrathChargeActive 判斷玩家是否在怒氣蓄積模式
func (m *luckyWrathChargeManager) isWrathChargeActive(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.activeSessions[playerID]
	return ok && !sess.settled
}

// tryLuckyWrathChargeFish 嘗試觸發幸運怒氣蓄積魚
func (g *Game) tryLuckyWrathChargeFish(p *player.Player) {
	m := g.LuckyWrathCharge
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
	if sess, ok := m.activeSessions[p.ID]; ok && !sess.settled {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyWrathChargePersonalCD)
	m.globalCooldownUntil = now.Add(LuckyWrathChargeGlobalCD)

	// 建立會話
	session := &wrathChargeSession{
		playerID:   p.ID,
		playerName: p.DisplayName,
		wrathValue: 0,
		expiresAt:  now.Add(WrathChargeTimeout),
	}
	m.activeSessions[p.ID] = session
	m.mu.Unlock()

	log.Printf("[LuckyWrathCharge] 觸發！玩家=%s 蓄積時間=%ds", p.DisplayName, int(WrathChargeTimeout.Seconds()))

	// 廣播怒氣蓄積開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyWrathCharge,
		Payload: ws.LuckyWrathChargePayload{
			Event:      "wrath_start",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			WrathValue: 0,
			MaxWrath:   WrathChargeMax,
			TimeoutSec: int(WrathChargeTimeout.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyWrathCharge, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🔥 %s 開始蓄積怒氣！打越多魚，隕石越多！最多 %d 顆！", p.DisplayName, WrathChargeMax),
		"color":   "#8B0000",
	})
	g.broadcastAnnouncement(ann)

	// 啟動超時計時器
	go g.runWrathChargeTimeout(session)
}

// notifyWrathChargeKill 玩家在怒氣蓄積模式中擊破目標（由 handleKill 呼叫）
func (g *Game) notifyWrathChargeKill(p *player.Player) {
	m := g.LuckyWrathCharge
	added, wrathValue := m.addWrathCharge(p.ID)
	if !added {
		return
	}

	// 廣播怒氣更新
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyWrathCharge,
		Payload: ws.LuckyWrathChargePayload{
			Event:      "wrath_charge",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			WrathValue: wrathValue,
			MaxWrath:   WrathChargeMax,
		},
	})

	// 怒氣值達到最大 → 立即爆發
	if wrathValue >= WrathChargeMax {
		m.mu.Lock()
		sess, ok := m.activeSessions[p.ID]
		if ok && !sess.settled {
			sess.settled = true
			m.mu.Unlock()
			go g.doWrathChargeExplosion(p, wrathValue)
		} else {
			m.mu.Unlock()
		}
	}
}

// doWrathChargeExplosion 怒氣爆發（隕石雨）
func (g *Game) doWrathChargeExplosion(p *player.Player, wrathValue int) {
	meteorCount := wrathValue
	if meteorCount < WrathChargeMin {
		meteorCount = WrathChargeMin
	}
	isPerfect := wrathValue >= WrathPerfectThreshold

	log.Printf("[LuckyWrathCharge] 爆發！玩家=%s 怒氣=%d 隕石=%d 完美=%v",
		p.DisplayName, wrathValue, meteorCount, isPerfect)

	// 廣播爆發開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyWrathCharge,
		Payload: ws.LuckyWrathChargePayload{
			Event:        "wrath_explode",
			PlayerID:     p.ID,
			PlayerName:   p.DisplayName,
			WrathValue:   wrathValue,
			MeteorCount:  meteorCount,
			IsPerfect:    isPerfect,
		},
	})

	// 依序發射隕石
	totalHit := 0
	for i := 0; i < meteorCount; i++ {
		time.Sleep(WrathMeteorInterval)

		// 隨機位置
		mx := 80.0 + rand.Float64()*640.0
		my := 80.0 + rand.Float64()*440.0

		// AOE 傷害
		hitCount := g.applyWrathMeteorDamage(mx, my, WrathMeteorRadius, WrathMeteorHPDmg)
		totalHit += hitCount

		// 廣播單顆隕石
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyWrathCharge,
			Payload: ws.LuckyWrathChargePayload{
				Event:       "wrath_meteor",
				PlayerID:    p.ID,
				MeteorIdx:   i + 1,
				MeteorTotal: meteorCount,
				MeteorX:     mx,
				MeteorY:     my,
				HitCount:    hitCount,
			},
		})
	}

	// 結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyWrathCharge,
		Payload: ws.LuckyWrathChargePayload{
			Event:       "wrath_settle",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			WrathValue:  wrathValue,
			MeteorCount: meteorCount,
			TotalHit:    totalHit,
			IsPerfect:   isPerfect,
		},
	})

	// 完美怒氣：全服加成
	if isPerfect {
		g.doWrathPerfect(p, wrathValue)
	}

	// 公告
	ann := g.Announce.Create(announce.EventLuckyWrathCharge, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("💥 %s 怒氣爆發！%d 顆隕石，命中 %d 個目標！", p.DisplayName, meteorCount, totalHit),
		"color":   "#FF4500",
	})
	g.broadcastAnnouncement(ann)
}

// applyWrathMeteorDamage 隕石 AOE 傷害
func (g *Game) applyWrathMeteorDamage(mx, my, radius, dmgPct float64) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	hitCount := 0
	for _, t := range g.Targets {
		if t == nil || t.HP <= 0 {
			continue
		}
		dx := t.X - mx
		dy := t.Y - my
		if dx*dx+dy*dy <= radius*radius {
			dmg := int(float64(t.HP) * dmgPct)
			if dmg < 1 {
				dmg = 1
			}
			t.HP -= dmg
			if t.HP < 1 {
				t.HP = 1 // 不直接擊殺，留給玩家補刀
			}
			hitCount++
			// 廣播 HP 更新
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
	return hitCount
}

// doWrathPerfect 完美怒氣全服加成
func (g *Game) doWrathPerfect(p *player.Player, wrathValue int) {
	m := g.LuckyWrathCharge
	m.mu.Lock()
	m.perfectBoost = &wrathPerfectBoost{
		mult:      WrathPerfectMult,
		expiresAt: time.Now().Add(WrathPerfectDuration),
	}
	m.mu.Unlock()

	log.Printf("[LuckyWrathCharge] 完美怒氣！玩家=%s 怒氣=%d 全服×%.1f 持續%ds",
		p.DisplayName, wrathValue, WrathPerfectMult, int(WrathPerfectDuration.Seconds()))

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyWrathCharge,
		Payload: ws.LuckyWrathChargePayload{
			Event:             "wrath_perfect",
			PlayerID:          p.ID,
			PlayerName:        p.DisplayName,
			WrathValue:        wrathValue,
			GlobalMult:        WrathPerfectMult,
			GlobalDurationSec: int(WrathPerfectDuration.Seconds()),
		},
	})

	ann := g.Announce.Create(announce.EventLuckyWrathCharge, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🔥💥 完美怒氣！%s 怒氣值 %d！全服 ×%.1f 加成 %d 秒！",
			p.DisplayName, wrathValue, WrathPerfectMult, int(WrathPerfectDuration.Seconds())),
		"color": "#FF4500",
	})
	g.broadcastAnnouncement(ann)

	go func() {
		time.Sleep(WrathPerfectDuration)
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyWrathCharge,
			Payload: ws.LuckyWrathChargePayload{
				Event: "wrath_perfect_end",
			},
		})
	}()
}

// runWrathChargeTimeout 怒氣蓄積超時計時器
func (g *Game) runWrathChargeTimeout(session *wrathChargeSession) {
	time.Sleep(WrathChargeTimeout)

	m := g.LuckyWrathCharge
	m.mu.Lock()
	sess, ok := m.activeSessions[session.playerID]
	if !ok || sess != session || sess.settled {
		m.mu.Unlock()
		return
	}
	sess.settled = true
	wrathValue := sess.wrathValue
	m.mu.Unlock()

	log.Printf("[LuckyWrathCharge] 超時爆發！玩家=%s 怒氣=%d", session.playerName, wrathValue)

	// 找到玩家
	g.mu.Lock()
	p, ok := g.Players[session.playerID]
	g.mu.Unlock()

	if !ok || p == nil {
		return
	}

	go g.doWrathChargeExplosion(p, wrathValue)
}
