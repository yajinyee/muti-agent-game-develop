// lucky_chain_explosion_handler.go — T115 幸運連鎖爆炸魚系統
// server-event-agent 負責維護
// 業界依據：Classic Arcade Fishing「chain explosion — each kill triggers next explosion,
//           building a cascade of rewards」
// 設計：擊破 T115 後，觸發「連鎖爆炸模式」（持續 12 秒）；
//       模式期間，玩家每次擊破任何目標 → 觸發一次小爆炸（AOE r=120px，HP -30%）；
//       每次爆炸命中 ≥ 1 個目標 → 連鎖計數 +1，倍率 +0.5x（最高 ×8.0）；
//       連鎖計數 ≥ 6 → 「連鎖爆發」：全服 ×2.5 加成 6 秒；
//       個人冷卻 22 秒；全服冷卻 38 秒
package game

import (
	"log"
	"math"
	"time"

	"chiikawa-game/internal/protocol"
)

// luckyChainExplosionManager 管理連鎖爆炸系統
type luckyChainExplosionManager struct {
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
	// 活躍的連鎖爆炸會話（per player）
	activeSessions map[string]*chainExplosionSession
	// 連鎖爆發全服加成
	burstBoost *chainBurstBoost
}

type chainBurstBoost struct {
	mult      float64
	expiresAt time.Time
}

type chainExplosionSession struct {
	playerID    string
	playerName  string
	chainCount  int
	accumMult   float64
	totalReward int
	expiresAt   time.Time
	settled     bool
}

func newLuckyChainExplosionManager() *luckyChainExplosionManager {
	return &luckyChainExplosionManager{
		playerCooldowns: make(map[string]time.Time),
		activeSessions:  make(map[string]*chainExplosionSession),
	}
}

// isLuckyChainExplosionFish 判斷是否為連鎖爆炸魚
func isLuckyChainExplosionFish(defID string) bool {
	return defID == "T115"
}

// isChainExplosionActive 判斷玩家是否在連鎖爆炸模式
func (m *luckyChainExplosionManager) isChainExplosionActive(playerID string) bool {
	s, ok := m.activeSessions[playerID]
	if !ok {
		return false
	}
	return !s.settled && time.Now().Before(s.expiresAt)
}

// getChainBurstMult 取得連鎖爆發全服倍率（供 handleKill 使用）
func (m *luckyChainExplosionManager) getChainBurstMult() float64 {
	if m.burstBoost != nil && time.Now().Before(m.burstBoost.expiresAt) {
		return m.burstBoost.mult
	}
	return 1.0
}

func (m *luckyChainExplosionManager) canTrigger(playerID string) bool {
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

// tryLuckyChainExplosion 嘗試觸發連鎖爆炸
func (g *Game) tryLuckyChainExplosion(playerID string, killerName string) {
	m := g.luckyChainExplosion
	if !m.canTrigger(playerID) {
		return
	}

	now := time.Now()
	m.playerCooldowns[playerID] = now.Add(22 * time.Second)
	m.globalCooldown = now.Add(38 * time.Second)

	// 建立會話
	session := &chainExplosionSession{
		playerID:   playerID,
		playerName: killerName,
		chainCount: 0,
		accumMult:  1.0,
		expiresAt:  now.Add(12 * time.Second),
	}
	m.activeSessions[playerID] = session

	g.hub.Broadcast(protocol.MsgLuckyChainExplosion, protocol.LuckyChainExplosionPayload{
		Event:       "chain_start",
		TriggerID:   playerID,
		TriggerName: killerName,
		Duration:    12.0,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "💥 " + killerName + " 觸發連鎖爆炸！12 秒連鎖模式！",
		Priority: "high",
		Color:    "#FF4500",
	})

	// 12 秒後超時結算
	go g.runChainExplosionTimeout(playerID, killerName)
}

// notifyChainExplosionKill 玩家在連鎖爆炸模式中擊破目標
func (g *Game) notifyChainExplosionKill(playerID string, killerName string, killX float64, killY float64) {
	m := g.luckyChainExplosion
	s, ok := m.activeSessions[playerID]
	if !ok || s.settled || time.Now().After(s.expiresAt) {
		return
	}

	p, pok := g.players[playerID]
	if !pok {
		return
	}
	betCost := p.GetBetDef().BetCost

	// AOE 爆炸（r=120px）
	aoeHits := []string{}
	for id, t := range g.targets {
		dx := t.X - killX
		dy := t.Y - killY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= 120 {
			damage := int(float64(t.MaxHP) * 0.3)
			t.HP -= damage
			if t.HP < 1 {
				t.HP = 1
			}
			g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
				InstanceID: t.InstanceID,
				HP:         t.HP,
				MaxHP:      t.MaxHP,
				X:          t.X,
				Y:          t.Y,
			})
			aoeHits = append(aoeHits, id)
		}
	}

	if len(aoeHits) >= 1 {
		s.chainCount++
		s.accumMult += 0.5
		if s.accumMult > 8.0 {
			s.accumMult = 8.0
		}
		chainReward := int(float64(betCost) * 0.5)
		s.totalReward += chainReward
		p.AddCoins(chainReward)

		g.hub.Broadcast(protocol.MsgLuckyChainExplosion, protocol.LuckyChainExplosionPayload{
			Event:       "chain_explode",
			TriggerID:   playerID,
			TriggerName: killerName,
			ExplodeX:    killX,
			ExplodeY:    killY,
			HitTargets:  aoeHits,
			ChainCount:  s.chainCount,
			AccumMult:   s.accumMult,
			TotalReward: s.totalReward,
		})

		// 連鎖爆發（≥ 6 次）
		if s.chainCount >= 6 && m.burstBoost == nil {
			m.burstBoost = &chainBurstBoost{
				mult:      2.5,
				expiresAt: time.Now().Add(6 * time.Second),
			}
			g.hub.Broadcast(protocol.MsgLuckyChainExplosion, protocol.LuckyChainExplosionPayload{
				Event:       "chain_burst",
				TriggerID:   playerID,
				TriggerName: killerName,
				ChainCount:  s.chainCount,
				AccumMult:   2.5,
				TotalReward: s.totalReward,
			})
			g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
				Message:  "💥🔥 連鎖爆發！" + killerName + " 全服 ×2.5 加成 6 秒！",
				Priority: "high",
				Color:    "#FF6B35",
			})
			go func() {
				time.Sleep(6 * time.Second)
				g.mu.Lock()
				m.burstBoost = nil
				g.mu.Unlock()
				g.hub.Broadcast(protocol.MsgLuckyChainExplosion, protocol.LuckyChainExplosionPayload{
					Event:       "burst_end",
					TriggerID:   playerID,
					TriggerName: killerName,
				})
			}()
		}
	}
}

// runChainExplosionTimeout 連鎖爆炸超時結算
func (g *Game) runChainExplosionTimeout(playerID string, killerName string) {
	time.Sleep(12 * time.Second)

	g.mu.Lock()
	m := g.luckyChainExplosion
	s, ok := m.activeSessions[playerID]
	if !ok || s.settled {
		g.mu.Unlock()
		return
	}
	s.settled = true
	delete(m.activeSessions, playerID)

	g.hub.Broadcast(protocol.MsgLuckyChainExplosion, protocol.LuckyChainExplosionPayload{
		Event:       "chain_end",
		TriggerID:   playerID,
		TriggerName: killerName,
		ChainCount:  s.chainCount,
		AccumMult:   s.accumMult,
		TotalReward: s.totalReward,
	})
	g.mu.Unlock()

	log.Printf("[ChainExplosion] Player %s: chains=%d, mult=%.1f, reward=%d",
		playerID, s.chainCount, s.accumMult, s.totalReward)
}
