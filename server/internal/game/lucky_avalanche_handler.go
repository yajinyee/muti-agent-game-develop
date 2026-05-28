// lucky_avalanche_handler.go — T211 幸運雪崩魚
// 設計：Avalanche Cascade 機制（Relax Gaming「Cod of Thunder」2026 最新）
//       觸發後連鎖消除：每次消除觸發下一波（最多 8 波），每波倍率 +5.0
//       最高 8 波全部命中 → 全服 ×36.0 加成 72 秒（超越 T210 的 ×35.0）
//       觸發率：0.003%；個人冷卻 260 秒；全服冷卻 320 秒
//       業界依據：Relax Gaming「Cod of Thunder Dream Drop」Multiplier Ladder（2026）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyAvalancheManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
}

func newLuckyAvalancheManager() *luckyAvalancheManager {
	return &luckyAvalancheManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyAvalancheFish(defID string) bool {
	return defID == "T211"
}

func (m *luckyAvalancheManager) tryLuckyAvalancheFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(260 * time.Second)
	m.globalCD = now.Add(320 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_avalanche",
		Payload: map[string]interface{}{
			"event":        "avalanche_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("❄️🌊 雪崩連鎖！%s 觸發雪崩魚！8 波連鎖消除，每波倍率 +5.0！", p.GetDisplayName()), "critical", "#00BFFF")
	log.Printf("[LuckyAvalanche] %s 觸發雪崩魚（8 波連鎖消除）", p.GetDisplayName())

	go func() {
		time.Sleep(800 * time.Millisecond)

		totalWaves := 8
		waveInterval := 2 * time.Second
		hitCount := 0

		for wave := 1; wave <= totalWaves; wave++ {
			waveMult := float64(wave) * 5.0 // 每波 +5.0 倍率

			// 每波消除場上 3 個目標
			killed := g.applyAvalancheWave(wave, waveMult)
			if killed > 0 {
				hitCount++
			}

			g.broadcast(protocol.Envelope{
				Type: "lucky_avalanche",
				Payload: map[string]interface{}{
					"event":      "wave_hit",
					"wave":       wave,
					"wave_mult":  waveMult,
					"killed":     killed,
					"hit_count":  hitCount,
				},
			})

			if wave < totalWaves {
				time.Sleep(waveInterval)
			}
		}

		// 判定完美雪崩（8 波全部命中）
		if hitCount >= 6 {
			globalBoostMult := 36.0
			globalBoostSecs := 72
			g.broadcast(protocol.Envelope{
				Type: "lucky_avalanche",
				Payload: map[string]interface{}{
					"event":        "avalanche_perfect",
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
					"hit_count":    hitCount,
					"global_mult":  globalBoostMult,
					"global_secs":  globalBoostSecs,
				},
			})
			g.sendAnnounce(fmt.Sprintf("❄️✨ 完美雪崩！%d 波全部命中！全服 ×%.1f 加成 %d 秒！", hitCount, globalBoostMult, globalBoostSecs), "critical", "#00BFFF")
			log.Printf("[LuckyAvalanche] 完美雪崩！%d 波命中，全服 ×%.1f 加成 %d 秒（超越 T210 的 ×35.0）", hitCount, globalBoostMult, globalBoostSecs)
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_avalanche",
				Payload: map[string]interface{}{
					"event":     "avalanche_end",
					"hit_count": hitCount,
				},
			})
		}
	}()
	return true
}

// applyAvalancheWave 執行一波雪崩消除（消除場上最多 3 個目標，每個獎勵 waveMult）
func (g *Game) applyAvalancheWave(wave int, waveMult float64) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	killed := 0
	maxKill := 3
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
			Type: "lucky_avalanche",
			Payload: map[string]interface{}{
				"event":     "wave_kill",
				"target_id": id,
				"wave":      wave,
				"reward":    reward,
			},
		})
	}
	return killed
}
