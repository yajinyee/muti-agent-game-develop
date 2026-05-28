// lucky_fever_boost_handler.go — T206 幸運 Fever Boost 魚
// 設計：Fever Boost™ 機制（Games Global 2026-05-28 最新）
//       觸發後 30 秒內所有特效機率翻倍，全場目標倍率 ×2.0
//       觸發後全服 ×31.0 加成 62 秒（超越 T205 的 ×30.0）
//       觸發率：0.003%（最稀有）；個人冷卻 250 秒；全服冷卻 310 秒
//       業界依據：Games Global「Fishin' Pots of Gold」Fever Boost™（2026-05-28）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyFeverBoostManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	feverBoost *feverBoostState
}

type feverBoostState struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyFeverBoostManager() *luckyFeverBoostManager {
	return &luckyFeverBoostManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyFeverBoostFish(defID string) bool {
	return defID == "T206"
}

func (m *luckyFeverBoostManager) getFeverBoostMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.feverBoost != nil && time.Now().Before(m.feverBoost.expiresAt) {
		return m.feverBoost.mult
	}
	return 1.0
}

func (m *luckyFeverBoostManager) tryLuckyFeverBoostFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(250 * time.Second)
	m.globalCD = now.Add(310 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_fever_boost",
		Payload: map[string]interface{}{
			"event":        "fever_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("🔥⚡ Fever Boost！%s 觸發 Fever Boost！30 秒內所有特效機率翻倍！全場倍率 ×2.0！", p.GetDisplayName()), "critical", "#FF6600")
	log.Printf("[LuckyFeverBoost] %s 觸發 Fever Boost 魚（全場倍率 ×2.0，30 秒）", p.GetDisplayName())

	go func() {
		time.Sleep(800 * time.Millisecond)

		// Fever Boost：全場目標倍率 ×2.0，持續 30 秒
		boostMult := 2.0
		boostSecs := 30
		m.mu.Lock()
		m.feverBoost = &feverBoostState{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_fever_boost",
			Payload: map[string]interface{}{
				"event":        "fever_active",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
			},
		})

		// 等待 Fever Boost 結束
		time.Sleep(time.Duration(boostSecs) * time.Second)

		// 觸發全服 ×31.0 加成 62 秒（超越 T205 的 ×30.0）
		globalBoostMult := 31.0
		globalBoostSecs := 62
		m.mu.Lock()
		m.feverBoost = &feverBoostState{
			mult:      globalBoostMult,
			expiresAt: time.Now().Add(time.Duration(globalBoostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_fever_boost",
			Payload: map[string]interface{}{
				"event":        "fever_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"global_mult":  globalBoostMult,
				"global_secs":  globalBoostSecs,
			},
		})
		g.sendAnnounce(fmt.Sprintf("🔥🌟 Fever Boost 完成！全服 ×%.1f 加成 %d 秒！", globalBoostMult, globalBoostSecs), "critical", "#FF6600")
		log.Printf("[LuckyFeverBoost] Fever Boost 完成！全服 ×%.1f 加成 %d 秒（超越 T205 的 ×30.0）", globalBoostMult, globalBoostSecs)
	}()
	return true
}
