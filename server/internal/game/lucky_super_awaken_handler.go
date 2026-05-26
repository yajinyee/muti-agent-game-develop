// lucky_super_awaken_handler.go — T153 幸運超級覺醒魚
// 業界依據：Jili「Super Awakening Performance, bonus up to 3000x」
// 設計：擊破後超級覺醒，全場 HP 歸零（每個獎勵 ×4.0），觸發全服 ×7.0 加成 15 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckySuperAwakenManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *superAwakenPerfectBoost
}

type superAwakenPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckySuperAwakenManager() *luckySuperAwakenManager {
	return &luckySuperAwakenManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckySuperAwakenFish(defID string) bool {
	return defID == "T153"
}

func (m *luckySuperAwakenManager) getSuperAwakenPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckySuperAwakenManager) tryLuckySuperAwakenFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) {
		m.mu.Unlock()
		return false
	}
	if cd, ok := m.personalCD[p.ID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return false
	}
	m.personalCD[p.ID] = now.Add(35 * time.Second)
	m.globalCD = now.Add(55 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_super_awaken",
		Payload: map[string]interface{}{
			"event":        "super_awaken_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("⚡ %s 觸發超級覺醒！全場 HP 歸零！", p.GetDisplayName()), "critical", "#FF6F00")
	log.Printf("[LuckySuperAwaken] %s 觸發超級覺醒", p.GetDisplayName())

	// 全場 HP 歸零
	g.mu.Lock()
	hitCount := 0
	totalReward := 0
	for id, t := range g.targets {
		if t.HP > 0 {
			t.HP = 0
			delete(g.targets, id)
			reward := int(float64(p.GetBetDef().BetCost) * t.Def.Multiplier * 4.0)
			p.AddCoins(reward)
			totalReward += reward
			hitCount++
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgTargetKill,
				Payload: protocol.TargetKillPayload{
					InstanceID: t.InstanceID,
					DefID:      t.Def.ID,
					Multiplier: t.Def.Multiplier * 4.0,
					Reward:     reward,
					LaborGain:  0,
					KillerID:   p.ID,
				},
			})
		}
	}
	g.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_super_awaken",
		Payload: map[string]interface{}{
			"event":        "super_awaken_result",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"hit_count":    hitCount,
			"total_reward": totalReward,
		},
	})

	// 觸發全服 ×7.0 加成 15 秒
	m.mu.Lock()
	m.perfectBoost = &superAwakenPerfectBoost{
		mult:      7.0,
		expiresAt: time.Now().Add(15 * time.Second),
	}
	m.mu.Unlock()
	g.broadcast(protocol.Envelope{
		Type: "lucky_super_awaken",
		Payload: map[string]interface{}{
			"event":        "super_awaken_boost",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"boost_mult":   7.0,
			"boost_secs":   15,
		},
	})
	g.sendAnnounce(fmt.Sprintf("⚡✨ 超級覺醒！%s 全場審判 %d 條！全服 ×7.0 加成 15 秒！", p.GetDisplayName(), hitCount), "critical", "#FFD700")

	time.AfterFunc(15*time.Second, func() {
		m.mu.Lock()
		m.perfectBoost = nil
		m.mu.Unlock()
		g.broadcast(protocol.Envelope{
			Type: "lucky_super_awaken",
			Payload: map[string]interface{}{
				"event":      "super_awaken_boost_end",
				"trigger_id": p.ID,
			},
		})
	})
	return true
}
