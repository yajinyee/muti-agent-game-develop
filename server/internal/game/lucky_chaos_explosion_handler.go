// lucky_chaos_explosion_handler.go — T198 幸運混沌爆炸魚
// 設計：混沌爆炸，隨機 3-8 個目標同時爆炸，倍率疊加最高 ×30.0
//       觸發後全服 ×24.0 加成 48 秒（超越 T197 的 ×23.5）
//       觸發率：0.014%；個人冷卻 155 秒；全服冷卻 220 秒
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyChaosExplosionManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	chaosBoost *chaosExplosionBoost
}

type chaosExplosionBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyChaosExplosionManager() *luckyChaosExplosionManager {
	return &luckyChaosExplosionManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyChaosExplosionFish(defID string) bool {
	return defID == "T198"
}

func (m *luckyChaosExplosionManager) getChaosExplosionMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.chaosBoost != nil && time.Now().Before(m.chaosBoost.expiresAt) {
		return m.chaosBoost.mult
	}
	return 1.0
}

func (m *luckyChaosExplosionManager) tryLuckyChaosExplosionFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(155 * time.Second)
	m.globalCD = now.Add(220 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_chaos_explosion",
		Payload: map[string]interface{}{
			"event":        "chaos_explosion_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("💥🌪️ 混沌爆炸！%s 引發混沌！隨機目標同時爆炸！倍率疊加最高 ×30.0！", p.GetDisplayName()), "critical", "#1A0A00")
	log.Printf("[LuckyChaosExplosion] %s 觸發混沌爆炸魚", p.GetDisplayName())

	go func() {
		time.Sleep(600 * time.Millisecond)

		// 隨機選 3-8 個目標同時爆炸
		explodeCount := 3 + rand.Intn(6) // 3-8
		totalReward := 0
		totalMult := 0.0

		g.mu.Lock()
		// 收集存活目標，按倍率排序（高倍率優先）
		type targetInfo struct {
			id   string
			mult float64
		}
		var targets []targetInfo
		for id, t := range g.targets {
			if t.HP > 0 {
				targets = append(targets, targetInfo{id: id, mult: t.Multiplier})
			}
		}
		sort.Slice(targets, func(i, j int) bool {
			return targets[i].mult > targets[j].mult
		})
		if explodeCount > len(targets) {
			explodeCount = len(targets)
		}

		explodedIDs := make([]string, 0, explodeCount)
		for i := 0; i < explodeCount; i++ {
			t := g.targets[targets[i].id]
			if t == nil {
				continue
			}
			// 混沌爆炸：每個目標獎勵 ×3.0
			reward := int(float64(p.GetBetDef().BetCost) * t.Multiplier * 3.0)
			if reward < 1 {
				reward = 1
			}
			totalReward += reward
			totalMult += t.Multiplier * 3.0
			p.Coins += reward
			explodedIDs = append(explodedIDs, targets[i].id)
			delete(g.targets, targets[i].id)
			g.hub.Broadcast(protocol.MsgTargetKill, protocol.TargetKillPayload{
				InstanceID: targets[i].id,
				Reward:     reward,
			})
		}
		g.mu.Unlock()
		g.sendPlayerUpdate(p.ID)

		if totalMult > 30.0 {
			totalMult = 30.0
		}

		// 觸發全服 ×24.0 加成 48 秒
		boostMult := 24.0
		boostSecs := 48
		m.mu.Lock()
		m.chaosBoost = &chaosExplosionBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_chaos_explosion",
			Payload: map[string]interface{}{
				"event":          "chaos_explosion_complete",
				"trigger_id":     p.ID,
				"trigger_name":   p.GetDisplayName(),
				"explode_count":  explodeCount,
				"total_mult":     totalMult,
				"total_reward":   totalReward,
				"boost_mult":     boostMult,
				"boost_secs":     boostSecs,
			},
		})
		g.sendAnnounce(fmt.Sprintf("💥✨ 混沌爆炸完成！%s 爆炸 %d 個目標！倍率疊加 ×%.1f！全服 ×%.1f 加成 %d 秒！",
			p.GetDisplayName(), explodeCount, totalMult, boostMult, boostSecs), "critical", "#2A1500")
	}()
	return true
}
