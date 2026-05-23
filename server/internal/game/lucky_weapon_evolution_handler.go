// lucky_weapon_evolution_handler.go — 幸運武器進化魚系統（DAY-252）
// 業界原創「武器進化+穿透+武器爆發」機制
//
// 設計：擊破 T210 後，玩家武器「進化」（持續 12 秒）：
//   - 進化等級 1→2：命中率 +30%，倍率 ×1.5（乘法）
//   - 進化期間再次擊破 T210 → 等級 2→3：穿透效果（子彈穿透第一個目標繼續飛行），倍率 ×2.5
//   - 進化結束時「武器爆發」：自動觸發 3 連射（×1.2 倍率，個人獎勵）
//   - 個人冷卻 18 秒；全服冷卻 25 秒
//
// 設計差異：
//   - 與充能魚（T183，累積能量爆發）不同，武器進化是「持續強化武器」，讓玩家有「這 12 秒我的武器更強」的掌控感
//   - 「等級 3 穿透」讓玩家有「要趁穿透期間打一排魚」的策略感
//   - 「武器爆發 3 連射」讓玩家有「進化結束前的最後爆發」的高潮設計
//   - 「再次擊破 T210 升級」讓玩家有「要趁進化期間找到另一條 T210」的動機
//   - 全服廣播讓其他玩家看到「有人武器進化了」，製造羨慕感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyWeaponEvoPersonalCD  = 18 * time.Second // 個人冷卻
	LuckyWeaponEvoGlobalCD    = 25 * time.Second // 全服冷卻
	LuckyWeaponEvoDuration    = 12 * time.Second // 進化持續時間
	LuckyWeaponEvoLv2Mult     = 1.5              // 等級 2 倍率加成（乘法）
	LuckyWeaponEvoLv3Mult     = 2.5              // 等級 3 倍率加成（乘法）
	LuckyWeaponEvoBurstMult   = 1.2              // 武器爆發倍率（個人）
	LuckyWeaponEvoBurstShots  = 3                // 武器爆發連射數
	LuckyWeaponEvoHitBonus    = 0.30             // 等級 2 命中率加成
)

// weaponEvoSession 武器進化會話
type weaponEvoSession struct {
	playerID    string
	playerName  string
	level       int       // 1=基礎, 2=進化, 3=穿透
	expiresAt   time.Time
	mu          sync.Mutex
}

// luckyWeaponEvoManager 幸運武器進化魚管理器
type luckyWeaponEvoManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 活躍的進化會話（playerID → session）
	activeSessions map[string]*weaponEvoSession
}

func newLuckyWeaponEvoManager() *luckyWeaponEvoManager {
	return &luckyWeaponEvoManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*weaponEvoSession),
	}
}

// isLuckyWeaponEvoFish 判斷是否為幸運武器進化魚
func isLuckyWeaponEvoFish(defID string) bool {
	return defID == "T210"
}

// isWeaponEvoActive 判斷玩家是否在武器進化狀態（供 handleAttack/handleKill 使用）
func (m *luckyWeaponEvoManager) isWeaponEvoActive(playerID string) (bool, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sess, ok := m.activeSessions[playerID]; ok {
		if time.Now().Before(sess.expiresAt) {
			return true, sess.level
		}
		// 已過期，清除
		delete(m.activeSessions, playerID)
	}
	return false, 0
}

// getLuckyWeaponEvoMult 取得武器進化倍率加成（供 handleKill 使用）
func (m *luckyWeaponEvoManager) getLuckyWeaponEvoMult(playerID string) float64 {
	active, level := m.isWeaponEvoActive(playerID)
	if !active {
		return 1.0
	}
	switch level {
	case 3:
		return LuckyWeaponEvoLv3Mult
	case 2:
		return LuckyWeaponEvoLv2Mult
	default:
		return 1.0
	}
}

// tryUpgradeWeaponEvo 嘗試升級武器進化等級（在進化期間再次擊破 T210）
func (m *luckyWeaponEvoManager) tryUpgradeWeaponEvo(playerID string) (bool, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sess, ok := m.activeSessions[playerID]; ok && time.Now().Before(sess.expiresAt) {
		sess.mu.Lock()
		defer sess.mu.Unlock()
		if sess.level < 3 {
			sess.level++
			return true, sess.level
		}
	}
	return false, 0
}

// tryLuckyWeaponEvoFish 擊破 T210 後觸發武器進化
func (g *Game) tryLuckyWeaponEvoFish(p *player.Player) {
	m := g.LuckyWeaponEvo

	// 先檢查是否已在進化狀態（升級）
	if upgraded, newLevel := m.tryUpgradeWeaponEvo(p.ID); upgraded {
		log.Printf("[WeaponEvo] player=%s 武器升級到等級 %d！", p.ID, newLevel)
		g.notifyWeaponEvoUpgrade(p, newLevel)
		return
	}

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
	m.personalCooldowns[p.ID] = now.Add(LuckyWeaponEvoPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyWeaponEvoGlobalCD)

	expiresAt := now.Add(LuckyWeaponEvoDuration)
	sess := &weaponEvoSession{
		playerID:   p.ID,
		playerName: p.DisplayName,
		level:      2, // 直接從等級 2 開始
		expiresAt:  expiresAt,
	}
	m.activeSessions[p.ID] = sess
	m.mu.Unlock()

	log.Printf("[WeaponEvo] player=%s 武器進化！等級 2，持續 %ds", p.ID, int(LuckyWeaponEvoDuration.Seconds()))

	// 個人訊息：進化啟動
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyWeaponEvo,
		Payload: ws.LuckyWeaponEvoPayload{
			Event:       "weapon_evo_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			Level:       2,
			DurationSec: int(LuckyWeaponEvoDuration.Seconds()),
			Mult:        LuckyWeaponEvoLv2Mult,
			HitBonus:    LuckyWeaponEvoHitBonus,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyWeaponEvo,
		Payload: ws.LuckyWeaponEvoPayload{
			Event:      "weapon_evo_broadcast",
			PlayerName: p.DisplayName,
			Level:      2,
			Mult:       LuckyWeaponEvoLv2Mult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyWeaponEvo, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⚔️ %s 武器進化！等級 2！命中率+30%%，倍率 ×%.1f！再擊破 T210 可升到等級 3！",
			p.DisplayName, LuckyWeaponEvoLv2Mult),
		"color": "#E67E22",
	})

	// 啟動進化計時 goroutine
	go g.runWeaponEvoDuration(p, expiresAt)
}

// notifyWeaponEvoUpgrade 武器升級到等級 3 的通知
func (g *Game) notifyWeaponEvoUpgrade(p *player.Player, newLevel int) {
	// 個人訊息：升級
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyWeaponEvo,
		Payload: ws.LuckyWeaponEvoPayload{
			Event:      "weapon_evo_upgrade",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Level:      newLevel,
			Mult:       LuckyWeaponEvoLv3Mult,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyWeaponEvo,
		Payload: ws.LuckyWeaponEvoPayload{
			Event:      "weapon_evo_broadcast",
			PlayerName: p.DisplayName,
			Level:      newLevel,
			Mult:       LuckyWeaponEvoLv3Mult,
		},
	})

	// 全服公告（等級 3 穿透）
	g.Announce.Create(announce.EventLuckyWeaponEvo, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⚔️ %s 武器升級到等級 3！穿透效果啟動！倍率 ×%.1f！",
			p.DisplayName, LuckyWeaponEvoLv3Mult),
		"color": "#E74C3C",
	})
}

// notifyWeaponEvoPierce 武器穿透效果：子彈穿透第一個目標後繼續飛行，命中第二個目標
// 由 handleAttack 在等級 3 時呼叫
func (g *Game) notifyWeaponEvoPierce(p *player.Player, firstTarget *target.Target) {
	// 找到第一個目標後方最近的目標（X 座標更小，即更靠近左側）
	g.mu.RLock()
	var pierceTarget *target.Target
	minDist := 300.0 // 穿透搜尋範圍 300px
	for _, t := range g.Targets {
		if !t.IsAlive || t.InstanceID == firstTarget.InstanceID {
			continue
		}
		// 穿透目標在第一個目標的後方（X 更小）
		if t.X >= firstTarget.X {
			continue
		}
		dist := firstTarget.X - t.X
		dy := t.Y - firstTarget.Y
		if dy < 0 {
			dy = -dy
		}
		if dy > 100 { // Y 偏差不超過 100px
			continue
		}
		if dist < minDist {
			minDist = dist
			pierceTarget = t
		}
	}
	g.mu.RUnlock()

	if pierceTarget == nil {
		return
	}

	// 穿透命中：60% 擊破機率，×0.8 倍率（個人獎勵）
	const pierceKillChance = 0.60
	const pierceMult = 0.8

	if rand.Float64() > pierceKillChance {
		return
	}

	// 消滅穿透目標
	g.mu.Lock()
	if t, ok := g.Targets[pierceTarget.InstanceID]; ok && t.IsAlive {
		t.IsAlive = false
		delete(g.Targets, pierceTarget.InstanceID)
	} else {
		g.mu.Unlock()
		return
	}
	g.mu.Unlock()

	// 計算獎勵
	reward := int(float64(data.GetBetDef(p.BetLevel).BetCost) * pierceTarget.Multiplier * pierceMult)
	if reward < 1 {
		reward = 1
	}
	p.AddCoins(reward)

	log.Printf("[WeaponEvo] player=%s 穿透命中 %s！獎勵 %d", p.ID, pierceTarget.InstanceID, reward)

	// 廣播穿透命中
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyWeaponEvo,
		Payload: ws.LuckyWeaponEvoPayload{
			Event:      "weapon_evo_pierce",
			PlayerName: p.DisplayName,
			InstanceID: pierceTarget.InstanceID,
			Mult:       pierceMult,
			Reward:     reward,
		},
	})
}

// runWeaponEvoDuration 武器進化計時 goroutine
func (g *Game) runWeaponEvoDuration(p *player.Player, expiresAt time.Time) {
	remaining := time.Until(expiresAt)
	if remaining <= 0 {
		return
	}

	timer := time.NewTimer(remaining)
	defer timer.Stop()

	select {
	case <-timer.C:
		g.doWeaponEvoBurst(p)
	case <-g.stopCh:
		return
	}
}

// doWeaponEvoBurst 武器爆發（進化結束時自動 3 連射）
func (g *Game) doWeaponEvoBurst(p *player.Player) {
	m := g.LuckyWeaponEvo
	m.mu.Lock()
	delete(m.activeSessions, p.ID)
	m.mu.Unlock()

	// 找場上最高倍率的 3 個目標進行爆發
	g.mu.RLock()
	type targetMult struct {
		t    *target.Target
		mult float64
	}
	var candidates []targetMult
	for _, t := range g.Targets {
		if t.IsAlive {
			candidates = append(candidates, targetMult{t: t, mult: t.Multiplier})
		}
	}
	g.mu.RUnlock()

	// 按倍率排序（簡單選擇排序，最多 3 個）
	for i := 0; i < len(candidates)-1; i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].mult > candidates[i].mult {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	burstCount := LuckyWeaponEvoBurstShots
	if len(candidates) < burstCount {
		burstCount = len(candidates)
	}

	if burstCount == 0 {
		// 沒有目標，只發送結束通知
		_ = g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyWeaponEvo,
			Payload: ws.LuckyWeaponEvoPayload{
				Event:      "weapon_evo_end",
				PlayerName: p.DisplayName,
				BurstCount: 0,
			},
		})
		return
	}

	totalReward := 0
	var burstTargets []string

	for i := 0; i < burstCount; i++ {
		t := candidates[i].t

		// 消滅目標
		g.mu.Lock()
		if existing, ok := g.Targets[t.InstanceID]; ok && existing.IsAlive {
			existing.IsAlive = false
			delete(g.Targets, t.InstanceID)
		} else {
			g.mu.Unlock()
			continue
		}
		g.mu.Unlock()

		reward := int(float64(data.GetBetDef(p.BetLevel).BetCost) * t.Multiplier * LuckyWeaponEvoBurstMult)
		if reward < 1 {
			reward = 1
		}
		p.AddCoins(reward)
		totalReward += reward
		burstTargets = append(burstTargets, t.InstanceID)
	}

	log.Printf("[WeaponEvo] player=%s 武器爆發！%d 連射，總獎勵 %d", p.ID, len(burstTargets), totalReward)

	// 個人結束通知（含爆發結算）
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyWeaponEvo,
		Payload: ws.LuckyWeaponEvoPayload{
			Event:       "weapon_evo_end",
			PlayerName:  p.DisplayName,
			BurstCount:  len(burstTargets),
			TotalReward: totalReward,
			BurstMult:   LuckyWeaponEvoBurstMult,
		},
	})

	// 全服廣播爆發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyWeaponEvo,
		Payload: ws.LuckyWeaponEvoPayload{
			Event:       "weapon_evo_burst",
			PlayerName:  p.DisplayName,
			BurstCount:  len(burstTargets),
			TotalReward: totalReward,
			BurstMult:   LuckyWeaponEvoBurstMult,
		},
	})

	if len(burstTargets) >= 2 {
		g.Announce.Create(announce.EventLuckyWeaponEvo, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("⚔️ %s 武器爆發！%d 連射！獲得 %d 籌碼！",
				p.DisplayName, len(burstTargets), totalReward),
			"color": "#E74C3C",
		})
	}
}
