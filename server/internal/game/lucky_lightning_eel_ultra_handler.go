// lucky_lightning_eel_ultra_handler.go — T245 幸運閃電鰻升級魚
// 設計：Lightning Eel Ultra 機制（Royal Fishing 60x Lightning Eel 升級版）
//       閃電連鎖跳躍：8 條鰻魚依序觸發（每條 ×90.0），連鎖跳躍 3 次
//       完美連鎖（全部命中）→ 全服 ×54.5 加成 109 秒
//       業界依據：Royal Fishing Jili「60x Lightning Eel Chain Reaction」升級版（2026）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type lightningEelUltraBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyLightningEelUltraManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *lightningEelUltraBoost
}

func newLuckyLightningEelUltraManager() *luckyLightningEelUltraManager {
	return &luckyLightningEelUltraManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyLightningEelUltraFish(defID string) bool {
	return defID == "T245"
}

func (m *luckyLightningEelUltraManager) getLightningEelUltraMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyLightningEelUltraManager) tryLuckyLightningEelUltraFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(590 * time.Second)
	m.personalCD[p.ID] = now.Add(530 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyLightningEelUltra,
		Payload: map[string]interface{}{
			"event":        "lightning_eel_ultra_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"eel_count":    8,
			"per_eel":      90.0,
			"jump_count":   3,
		},
	})
	g.sendAnnounce(fmt.Sprintf("LIGHTNING EEL ULTRA! %s unleashed 8 Lightning Eels! Each x90.0! Chain jumps x3!", p.GetDisplayName()), "critical", "#00FFFF")
	log.Printf("[LuckyLightningEelUltra] %s triggered Lightning Eel Ultra fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		perEelMult := 90.0
		eelCount := 8
		jumpCount := 3
		totalReward := 0
		totalMult := 0.0

		// 8 條鰻魚依序觸發，每條跳躍 3 次
		for eel := 1; eel <= eelCount; eel++ {
			time.Sleep(250 * time.Millisecond)

			for jump := 1; jump <= jumpCount; jump++ {
				time.Sleep(150 * time.Millisecond)
				jumpMult := perEelMult * float64(jump) * 0.4 // 跳躍倍率遞增
				jumpReward := int(jumpMult * betCost)
				totalReward += jumpReward
				totalMult += jumpMult

				g.mu.Lock()
				p.Coins += jumpReward
				g.mu.Unlock()

				g.broadcast(protocol.Envelope{
					Type: protocol.MsgLuckyLightningEelUltra,
					Payload: map[string]interface{}{
						"event":       "eel_chain_jump",
						"eel_no":      eel,
						"jump_no":     jump,
						"jump_mult":   jumpMult,
						"jump_reward": jumpReward,
					},
				})
			}
		}

		// 完美連鎖獎勵（全部 8 條命中）
		perfectBonus := perEelMult * float64(eelCount) * 0.5
		perfectReward := int(perfectBonus * betCost)
		totalReward += perfectReward
		totalMult += perfectBonus

		g.mu.Lock()
		p.Coins += perfectReward
		g.mu.Unlock()

		// 全服 ×54.5 加成 109 秒
		globalBonus := 54.5
		globalDuration := 109

		m.mu.Lock()
		m.boost = &lightningEelUltraBoost{
			mult:      globalBonus,
			expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyLightningEelUltra,
			Payload: map[string]interface{}{
				"event":          "lightning_eel_ultra_complete",
				"eel_count":      eelCount,
				"total_mult":     totalMult,
				"total_reward":   totalReward,
				"perfect_bonus":  perfectBonus,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
			},
		})
		g.sendAnnounce(fmt.Sprintf("LIGHTNING EEL ULTRA COMPLETE! %s: %d eels, total x%.1f! GLOBAL x%.1f for %ds!", p.GetDisplayName(), eelCount, totalMult, globalBonus, globalDuration), "critical", "#00FFFF")
		log.Printf("[LuckyLightningEelUltra] %s: eels=%d, total_mult=%.1f, global=x%.1f", p.GetDisplayName(), eelCount, totalMult, globalBonus)
	}()

	return true
}
