// lucky_immortal_boss_ultra_handler.go — T247 幸運不死BOSS升級魚
// 設計：Immortal Boss Ultra 機制（Royal Fishing Immortal Boss 升級版）
//       不死 BOSS 連續獎勵：BOSS 被擊敗後立即復活（HP 50%），連續 5 次
//       每次復活獎勵遞增（×100→×150→×200→×250→×300），全服 ×55.5 加成 111 秒
//       業界依據：Royal Fishing Jili「Immortal Boss 50-150x consecutive wins」升級版（2026）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type immortalBossUltraBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyImmortalBossUltraManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *immortalBossUltraBoost
}

func newLuckyImmortalBossUltraManager() *luckyImmortalBossUltraManager {
	return &luckyImmortalBossUltraManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyImmortalBossUltraFish(defID string) bool {
	return defID == "T247"
}

func (m *luckyImmortalBossUltraManager) getImmortalBossUltraMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyImmortalBossUltraManager) tryLuckyImmortalBossUltraFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(610 * time.Second)
	m.personalCD[p.ID] = now.Add(550 * time.Second)
	m.mu.Unlock()

	// 不死 BOSS 連續獎勵倍率
	reviveMults := []float64{100.0, 150.0, 200.0, 250.0, 300.0}

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyImmortalBossUltra,
		Payload: map[string]interface{}{
			"event":        "immortal_boss_ultra_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"revive_count": len(reviveMults),
			"max_mult":     300.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("IMMORTAL BOSS ULTRA! %s summoned the Immortal Boss! 5 revivals! Up to x300.0!", p.GetDisplayName()), "critical", "#8B0000")
	log.Printf("[LuckyImmortalBossUltra] %s triggered Immortal Boss Ultra fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		totalReward := 0
		totalMult := 0.0

		// 5 次不死復活
		for revive, mult := range reviveMults {
			time.Sleep(1200 * time.Millisecond)

			reviveReward := int(mult * betCost)
			totalReward += reviveReward
			totalMult += mult

			g.mu.Lock()
			p.Coins += reviveReward
			g.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyImmortalBossUltra,
				Payload: map[string]interface{}{
					"event":         "boss_revive",
					"revive_no":     revive + 1,
					"revive_mult":   mult,
					"revive_reward": reviveReward,
					"hp_percent":    50.0,
				},
			})
		}

		// 全服 ×55.5 加成 111 秒
		globalBonus := 55.5
		globalDuration := 111

		m.mu.Lock()
		m.boost = &immortalBossUltraBoost{
			mult:      globalBonus,
			expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyImmortalBossUltra,
			Payload: map[string]interface{}{
				"event":          "immortal_boss_ultra_complete",
				"revive_count":   len(reviveMults),
				"total_mult":     totalMult,
				"total_reward":   totalReward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
			},
		})
		g.sendAnnounce(fmt.Sprintf("IMMORTAL BOSS ULTRA COMPLETE! %s: 5 revivals, total x%.1f! GLOBAL x%.1f for %ds!", p.GetDisplayName(), totalMult, globalBonus, globalDuration), "critical", "#8B0000")
		log.Printf("[LuckyImmortalBossUltra] %s: revivals=5, total_mult=%.1f, global=x%.1f", p.GetDisplayName(), totalMult, globalBonus)
	}()

	return true
}
