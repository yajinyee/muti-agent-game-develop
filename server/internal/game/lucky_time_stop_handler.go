// lucky_time_stop_handler.go — T232 幸運時間停止魚
// 設計：Time Stop 機制
//       全場凍結 15 秒（傷害 ×5.0），凍結結束全場 HP -70%
//       凍結期間擊破 ≥15 個 → 完美時間停止，全服 ×47.0 加成 94 秒（超越 T231 的 ×46.5）
//       業界依據：時間凍結機制終極升級版 + Time Stop 概念（2026）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type timeStopBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyTimeStopManager struct {
	globalCD     time.Time
	mu           sync.Mutex
	personalCD   map[string]time.Time
	boost        *timeStopBoost
	freezeActive bool
	freezeKills  map[string]int
}

func newLuckyTimeStopManager() *luckyTimeStopManager {
	return &luckyTimeStopManager{
		personalCD:  make(map[string]time.Time),
		freezeKills: make(map[string]int),
	}
}

func isLuckyTimeStopFish(defID string) bool {
	return defID == "T232"
}

func (m *luckyTimeStopManager) getTimeStopMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyTimeStopManager) notifyFreezeKill(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.freezeActive {
		m.freezeKills[playerID]++
	}
}

func (m *luckyTimeStopManager) tryLuckyTimeStopFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(450 * time.Second)
	m.personalCD[p.ID] = now.Add(390 * time.Second)
	m.freezeActive = true
	m.freezeKills[p.ID] = 0
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyTimeStop,
		Payload: map[string]interface{}{
			"event":          "freeze_start",
			"trigger_id":     p.ID,
			"trigger_name":   p.GetDisplayName(),
			"freeze_seconds": 15,
			"damage_mult":    5.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("TIME STOP! %s froze all time! 15 seconds of x5.0 damage!", p.GetDisplayName()), "critical", "#00CCFF")
	log.Printf("[LuckyTimeStop] %s triggered Time Stop fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)

		for sec := 15; sec > 0; sec-- {
			time.Sleep(1 * time.Second)
			m.mu.Lock()
			kills := m.freezeKills[p.ID]
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyTimeStop,
				Payload: map[string]interface{}{
					"event":      "freeze_tick",
					"time_left":  sec - 1,
					"kill_count": kills,
				},
			})
		}

		m.mu.Lock()
		m.freezeActive = false
		killCount := m.freezeKills[p.ID]
		m.mu.Unlock()

		g.mu.Lock()
		targetCount := 0
		for _, t := range g.targets {
			damage := int(float64(t.MaxHP) * 0.70)
			t.HP -= damage
			if t.HP <= 0 {
				t.HP = 0
			}
			targetCount++
		}
		g.mu.Unlock()

		endMult := float64(targetCount) * 35.0
		endReward := int(endMult * betCost)
		if endReward > 0 {
			g.mu.Lock()
			p.Coins += endReward
			g.mu.Unlock()
		}

		isPerfect := killCount >= 15
		globalBonus := 47.0
		globalDuration := 94

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyTimeStop,
			Payload: map[string]interface{}{
				"event":          "freeze_end",
				"kill_count":     killCount,
				"target_count":   targetCount,
				"end_mult":       endMult,
				"end_reward":     endReward,
				"is_perfect":     isPerfect,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
			},
		})

		if isPerfect {
			m.mu.Lock()
			m.boost = &timeStopBoost{
				mult:      globalBonus,
				expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
			}
			m.mu.Unlock()
			g.sendAnnounce(fmt.Sprintf("PERFECT TIME STOP! %s killed %d during freeze! Global x%.1f for %ds!", p.GetDisplayName(), killCount, globalBonus, globalDuration), "critical", "#00CCFF")
		}

		log.Printf("[LuckyTimeStop] %s: kills=%d, end_mult=%.1f, perfect=%v", p.GetDisplayName(), killCount, endMult, isPerfect)
	}()

	return true
}
