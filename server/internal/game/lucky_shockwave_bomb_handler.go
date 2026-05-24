// lucky_shockwave_bomb_handler.go — T112 幸運全場震盪魚系統
// server-event-agent 負責維護
// 業界依據：Classic Arcade Fishing「explosion rate upgrade — shockwave bomb hits all fish on screen」
// 設計：擊破 T112 後，觸發「全場震盪」：
//   - 全場所有目標 HP -35%（震盪傷害）
//   - 觸發玩家進入「震盪強化」模式 10 秒：攻擊力 ×2.0
//   - 若震盪傷害命中 ≥ 10 個目標 → 「超級震盪」：全服 ×1.8 加成 6 秒
// 個人冷卻 18 秒；全服冷卻 30 秒
package game

import (
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

// shockwavePowerBoost 震盪強化（個人攻擊力加成）
type shockwavePowerBoost struct {
	playerID  string
	mult      float64
	expiresAt time.Time
}

// shockwaveSuperBoost 超級震盪全服加成
type shockwaveSuperBoost struct {
	mult      float64
	expiresAt time.Time
}

// luckyShockwaveBombManager 管理全場震盪系統
type luckyShockwaveBombManager struct {
	mu              sync.Mutex
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
	powerBoosts     map[string]*shockwavePowerBoost // playerID -> boost
	superBoost      *shockwaveSuperBoost
}

func newLuckyShockwaveBombManager() *luckyShockwaveBombManager {
	return &luckyShockwaveBombManager{
		playerCooldowns: make(map[string]time.Time),
		powerBoosts:     make(map[string]*shockwavePowerBoost),
	}
}

// isLuckyShockwaveBombFish 判斷是否為全場震盪魚
func isLuckyShockwaveBombFish(defID string) bool {
	return defID == "T112"
}

// canTrigger 判斷是否可以觸發
func (m *luckyShockwaveBombManager) canTrigger(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	if now.Before(m.globalCooldown) {
		return false
	}
	if cd, ok := m.playerCooldowns[playerID]; ok {
		if now.Before(cd) {
			return false
		}
	}
	return true
}

// getShockwavePowerMult 取得震盪強化個人攻擊力倍率
func (m *luckyShockwaveBombManager) getShockwavePowerMult(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	boost, ok := m.powerBoosts[playerID]
	if !ok {
		return 1.0
	}
	if time.Now().After(boost.expiresAt) {
		delete(m.powerBoosts, playerID)
		return 1.0
	}
	return boost.mult
}

// getShockwaveSuperMult 取得超級震盪全服倍率（供 handleKill 使用）
func (m *luckyShockwaveBombManager) getShockwaveSuperMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.superBoost == nil {
		return 1.0
	}
	if time.Now().After(m.superBoost.expiresAt) {
		m.superBoost = nil
		return 1.0
	}
	return m.superBoost.mult
}

// tryLuckyShockwaveBomb 嘗試觸發全場震盪
func (g *Game) tryLuckyShockwaveBomb(playerID string, killerName string) {
	m := g.luckyShockwaveBomb
	if !m.canTrigger(playerID) {
		return
	}

	m.mu.Lock()
	now := time.Now()
	m.playerCooldowns[playerID] = now.Add(18 * time.Second)
	m.globalCooldown = now.Add(30 * time.Second)

	// 設定個人震盪強化（10 秒攻擊力 ×2.0）
	m.powerBoosts[playerID] = &shockwavePowerBoost{
		playerID:  playerID,
		mult:      2.0,
		expiresAt: now.Add(10 * time.Second),
	}
	m.mu.Unlock()

	// 廣播觸發事件
	g.hub.Broadcast(protocol.MsgLuckyShockwaveBomb, protocol.LuckyShockwaveBombPayload{
		Event:       "shockwave_start",
		TriggerID:   playerID,
		TriggerName: killerName,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "💥 " + killerName + " 觸發全場震盪！全場 HP -35%！攻擊力 ×2.0 持續 10 秒！",
		Priority: "high",
		Color:    "#FF4500",
	})

	// 在 goroutine 中執行震盪邏輯
	go g.runShockwaveBomb(playerID, killerName)
}

// runShockwaveBomb 執行全場震盪邏輯
func (g *Game) runShockwaveBomb(playerID string, killerName string) {
	// 等待 500ms 讓 Client 顯示震盪動畫
	time.Sleep(500 * time.Millisecond)

	g.mu.Lock()
	p, ok := g.players[playerID]
	if !ok {
		g.mu.Unlock()
		return
	}
	betCost := p.GetBetDef().BetCost

	// 全場震盪：所有目標 HP -35%
	hitCount := 0
	totalReward := 0
	for _, t := range g.targets {
		if t.Def.Type == "boss" {
			continue
		}
		damage := t.MaxHP * 35 / 100
		t.HP -= damage
		if t.HP < 1 {
			t.HP = 1
		}
		hitCount++

		// 廣播 HP 更新
		g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
			InstanceID: t.InstanceID,
			HP:         t.HP,
			MaxHP:      t.MaxHP,
			X:          t.X,
			Y:          t.Y,
		})
	}

	// 震盪獎勵：命中數 × betCost × 0.5
	totalReward = hitCount * betCost / 2
	if ok {
		p.AddCoins(totalReward)
		g.sendPlayerUpdate(playerID)
	}
	g.mu.Unlock()

	// 廣播震盪結果
	g.hub.Broadcast(protocol.MsgLuckyShockwaveBomb, protocol.LuckyShockwaveBombPayload{
		Event:       "shockwave_hit",
		TriggerID:   playerID,
		TriggerName: killerName,
		HitCount:    hitCount,
		TotalReward: totalReward,
	})

	// 超級震盪：命中 ≥ 10 個目標
	if hitCount >= 10 {
		g.doShockwaveSuper(playerID, killerName)
	}

	// 10 秒後震盪強化結束
	go func() {
		time.Sleep(10 * time.Second)
		g.hub.Broadcast(protocol.MsgLuckyShockwaveBomb, protocol.LuckyShockwaveBombPayload{
			Event:       "power_end",
			TriggerID:   playerID,
			TriggerName: killerName,
		})
		g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
			Message:  "💥 " + killerName + " 的震盪強化結束",
			Priority: "low",
			Color:    "#888888",
		})
	}()

	log.Printf("[ShockwaveBomb] Player %s: hitCount=%d, reward=%d", playerID, hitCount, totalReward)
}

// doShockwaveSuper 觸發超級震盪全服加成
func (g *Game) doShockwaveSuper(playerID string, killerName string) {
	m := g.luckyShockwaveBomb
	m.mu.Lock()
	m.superBoost = &shockwaveSuperBoost{
		mult:      1.8,
		expiresAt: time.Now().Add(6 * time.Second),
	}
	m.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyShockwaveBomb, protocol.LuckyShockwaveBombPayload{
		Event:       "super_shockwave",
		TriggerID:   playerID,
		TriggerName: killerName,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "💥🌊 超級震盪！" + killerName + " 全服 ×1.8 加成 6 秒！",
		Priority: "high",
		Color:    "#FF6B35",
	})

	go func() {
		time.Sleep(6 * time.Second)
		g.hub.Broadcast(protocol.MsgLuckyShockwaveBomb, protocol.LuckyShockwaveBombPayload{
			Event:       "super_end",
			TriggerID:   playerID,
			TriggerName: killerName,
		})
	}()
}
