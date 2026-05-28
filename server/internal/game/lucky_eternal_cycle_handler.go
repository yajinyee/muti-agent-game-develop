// lucky_eternal_cycle_handler.go — T197 幸運永恆循環魚
// 設計：永恆循環 10 波，每波獎勵遞增（×1.0 → ×2.0 → ... → ×10.0）
//       全部 10 波完成 → 全服 ×23.5 加成 47 秒（超越 T196 的 ×23.0）
//       觸發率：0.016%；個人冷卻 150 秒；全服冷卻 215 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyEternalCycleManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	cycleBoost *eternalCycleBoost
}

type eternalCycleBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyEternalCycleManager() *luckyEternalCycleManager {
	return &luckyEternalCycleManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyEternalCycleFish(defID string) bool {
	return defID == "T197"
}

func (m *luckyEternalCycleManager) getEternalCycleMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cycleBoost != nil && time.Now().Before(m.cycleBoost.expiresAt) {
		return m.cycleBoost.mult
	}
	return 1.0
}

func (m *luckyEternalCycleManager) tryLuckyEternalCycleFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(150 * time.Second)
	m.globalCD = now.Add(215 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_eternal_cycle",
		Payload: map[string]interface{}{
			"event":        "eternal_cycle_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"total_waves":  10,
		},
	})
	g.sendAnnounce(fmt.Sprintf("♾️🌀 永恆循環！%s 啟動永恆循環！10 波遞增獎勵！", p.GetDisplayName()), "critical", "#0A001A")
	log.Printf("[LuckyEternalCycle] %s 觸發永恆循環魚", p.GetDisplayName())

	go func() {
		totalReward := 0
		// 10 波遞增：每波倍率 = 波次 × 1.0（1x, 2x, 3x, ... 10x）
		for wave := 1; wave <= 10; wave++ {
			time.Sleep(400 * time.Millisecond)
			waveMult := float64(wave)
			waveReward := int(float64(p.GetBetDef().BetCost) * waveMult)
			if waveReward < 1 {
				waveReward = 1
			}
			totalReward += waveReward

			g.mu.Lock()
			p.Coins += waveReward
			g.mu.Unlock()
			g.sendPlayerUpdate(p.ID)

			g.broadcast(protocol.Envelope{
				Type: "lucky_eternal_cycle",
				Payload: map[string]interface{}{
					"event":       "eternal_cycle_wave",
					"trigger_id":  p.ID,
					"wave":        wave,
					"wave_mult":   waveMult,
					"wave_reward": waveReward,
				},
			})
		}

		// 全部 10 波完成 → 全服 ×23.5 加成 47 秒
		boostMult := 23.5
		boostSecs := 47
		m.mu.Lock()
		m.cycleBoost = &eternalCycleBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_eternal_cycle",
			Payload: map[string]interface{}{
				"event":        "eternal_cycle_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"total_reward": totalReward,
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
			},
		})
		g.sendAnnounce(fmt.Sprintf("♾️✨ 永恆循環完成！%s 10 波全部完成！總獎勵 %d！全服 ×%.1f 加成 %d 秒！",
			p.GetDisplayName(), totalReward, boostMult, boostSecs), "critical", "#1A0030")
	}()
	return true
}
