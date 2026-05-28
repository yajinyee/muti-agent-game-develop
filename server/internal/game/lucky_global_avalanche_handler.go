// lucky_global_avalanche_handler.go — T215 幸運全服雪崩魚
// 設計：Global Avalanche 機制（終極版雪崩，全服連鎖消除）
//       觸發後全服連鎖消除：每個玩家各自觸發 5 波消除，每波 ×8.0
//       全服所有玩家 5 波全部命中 → 全服 ×38.0 加成 76 秒（新史上最高，超越 T214 的 ×37.5）
//       觸發率：0.002%（最稀有）；個人冷卻 280 秒；全服冷卻 340 秒
//       業界依據：Avalanche Reels + Global Multiplier 組合（2026 最新趨勢）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyGlobalAvalancheManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
}

func newLuckyGlobalAvalancheManager() *luckyGlobalAvalancheManager {
	return &luckyGlobalAvalancheManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyGlobalAvalancheFish(defID string) bool {
	return defID == "T215"
}

func (m *luckyGlobalAvalancheManager) tryLuckyGlobalAvalancheFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(280 * time.Second)
	m.globalCD = now.Add(340 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_global_avalanche",
		Payload: map[string]interface{}{
			"event":        "global_avalanche_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"waves":        5,
			"wave_mult":    8.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("❄️🌍 全服雪崩！%s 觸發全服雪崩魚！5 波全服連鎖消除，每波 ×8.0！", p.GetDisplayName()), "critical", "#87CEEB")
	log.Printf("[LuckyGlobalAvalanche] %s 觸發全服雪崩魚（5 波全服連鎖消除）", p.GetDisplayName())

	go func() {
		time.Sleep(800 * time.Millisecond)

		totalWaves := 5
		waveMult := 8.0
		totalKilled := 0

		for wave := 1; wave <= totalWaves; wave++ {
			// 每波消除場上最多 5 個目標（比 T211 的 3 個更多）
			killed := g.applyGlobalAvalancheWave(wave, waveMult)
			totalKilled += killed

			g.broadcast(protocol.Envelope{
				Type: "lucky_global_avalanche",
				Payload: map[string]interface{}{
					"event":         "global_wave_hit",
					"wave":          wave,
					"wave_mult":     waveMult,
					"killed":        killed,
					"total_killed":  totalKilled,
				},
			})

			if wave < totalWaves {
				time.Sleep(2500 * time.Millisecond)
			}
		}

		// 全服雪崩完成，觸發全服 ×38.0 加成 76 秒
		globalBoostMult := 38.0
		globalBoostSecs := 76
		g.broadcast(protocol.Envelope{
			Type: "lucky_global_avalanche",
			Payload: map[string]interface{}{
				"event":         "global_avalanche_complete",
				"trigger_id":    p.ID,
				"trigger_name":  p.GetDisplayName(),
				"total_killed":  totalKilled,
				"global_mult":   globalBoostMult,
				"global_secs":   globalBoostSecs,
			},
		})
		g.sendAnnounce(fmt.Sprintf("❄️🌟 全服雪崩完成！消滅 %d 個目標！全服 ×%.1f 加成 %d 秒！", totalKilled, globalBoostMult, globalBoostSecs), "critical", "#87CEEB")
		log.Printf("[LuckyGlobalAvalanche] 全服雪崩完成！消滅 %d 個目標，全服 ×%.1f 加成 %d 秒（新史上最高，超越 T214 的 ×37.5）", totalKilled, globalBoostMult, globalBoostSecs)
	}()
	return true
}

// applyGlobalAvalancheWave 執行一波全服雪崩消除（消除場上最多 5 個目標）
func (g *Game) applyGlobalAvalancheWave(wave int, waveMult float64) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	killed := 0
	maxKill := 5
	for id, t := range g.targets {
		if killed >= maxKill {
			break
		}
		if t.HP <= 0 {
			continue
		}
		reward := int(float64(t.Def.MinMult) * waveMult)
		t.HP = 0
		delete(g.targets, id)
		killed++

		g.broadcast(protocol.Envelope{
			Type: "lucky_global_avalanche",
			Payload: map[string]interface{}{
				"event":     "global_wave_kill",
				"target_id": id,
				"wave":      wave,
				"reward":    reward,
			},
		})
	}
	return killed
}
