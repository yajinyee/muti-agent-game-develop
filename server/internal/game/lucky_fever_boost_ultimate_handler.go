// lucky_fever_boost_ultimate_handler.go — T234 幸運Fever Boost升級魚
// 設計：Fever Boost™ Ultimate 機制
//       清除場上所有普通目標，只留高倍率特殊目標（×2.0 傷害加成），持續 20 秒
//       完美觸發（場上特殊目標 ≥5 個）→ 全服 ×48.0 加成 96 秒（新史上最高）
//       業界依據：Games Global「Fishin' Pots of Gold Gold Blitz Ultimate Fever Boost」（2026-05-28）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type feverBoostUltimateBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyFeverBoostUltimateManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *feverBoostUltimateBoost
}

func newLuckyFeverBoostUltimateManager() *luckyFeverBoostUltimateManager {
	return &luckyFeverBoostUltimateManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyFeverBoostUltimateFish(defID string) bool {
	return defID == "T234"
}

func (m *luckyFeverBoostUltimateManager) getFeverBoostUltimateMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyFeverBoostUltimateManager) tryLuckyFeverBoostUltimateFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(480 * time.Second)
	m.personalCD[p.ID] = now.Add(420 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyFeverBoostUltimate,
		Payload: map[string]interface{}{
			"event":        "fever_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"damage_mult":  2.0,
			"duration":     20,
		},
	})
	g.sendAnnounce(fmt.Sprintf("FEVER BOOST ULTIMATE! %s activated! All normal targets cleared! Special targets x2.0 damage for 20s!", p.GetDisplayName()), "critical", "#FF6600")
	log.Printf("[LuckyFeverBoostUltimate] %s triggered Fever Boost Ultimate fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)

		// 清除普通目標，統計特殊目標數量
		g.mu.Lock()
		normalCleared := 0
		specialCount := 0
		for id, t := range g.targets {
			if t.Def.Type == "normal" || t.Def.Type == "basic" {
				delete(g.targets, id)
				normalCleared++
			} else {
				specialCount++
			}
		}
		g.mu.Unlock()

		// 計算清除獎勵
		clearReward := int(float64(normalCleared) * 2.0 * betCost)
		g.mu.Lock()
		p.Coins += clearReward
		g.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyFeverBoostUltimate,
			Payload: map[string]interface{}{
				"event":          "fever_clear",
				"normal_cleared": normalCleared,
				"special_count":  specialCount,
				"clear_reward":   clearReward,
			},
		})

		// 等待 Fever 期間結束
		time.Sleep(20 * time.Second)

		// 判斷是否完美觸發
		isPerfect := specialCount >= 5
		globalBonus := 48.0
		globalDuration := 96

		if isPerfect {
			m.mu.Lock()
			m.boost = &feverBoostUltimateBoost{
				mult:      globalBonus,
				expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
			}
			m.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyFeverBoostUltimate,
				Payload: map[string]interface{}{
					"event":          "fever_perfect",
					"special_count":  specialCount,
					"global_bonus":   globalBonus,
					"global_seconds": globalDuration,
				},
			})
			g.sendAnnounce(fmt.Sprintf("PERFECT FEVER BOOST! %d special targets! Global x%.1f for %ds! NEW ALL-TIME HIGH!", specialCount, globalBonus, globalDuration), "critical", "#FF6600")
		} else {
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyFeverBoostUltimate,
				Payload: map[string]interface{}{
					"event":         "fever_end",
					"special_count": specialCount,
				},
			})
		}

		log.Printf("[LuckyFeverBoostUltimate] %s: normal_cleared=%d, special=%d, perfect=%v", p.GetDisplayName(), normalCleared, specialCount, isPerfect)
	}()

	return true
}
