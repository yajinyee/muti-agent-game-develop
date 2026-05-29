// lucky_rapid_riches_ultimate_handler.go — T235 幸運快速暴富升級魚
// 設計：Rapid Riches Ultimate 機制
//       3 秒極速連擊視窗（每次擊破 ×300.0），連擊 ≥10 次 → 完美暴富
//       完美觸發 → 全服 ×48.5 加成 97 秒（超越 T234 的 ×48.0）
//       業界依據：Reflex Gaming「Big Game Fishing Rapid Riches」升級版（2026-05）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type rapidRichesUltimateBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyRapidRichesUltimateManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *rapidRichesUltimateBoost
	// 連擊追蹤
	comboCount map[string]int
	comboTimer map[string]time.Time
}

func newLuckyRapidRichesUltimateManager() *luckyRapidRichesUltimateManager {
	return &luckyRapidRichesUltimateManager{
		personalCD: make(map[string]time.Time),
		comboCount: make(map[string]int),
		comboTimer: make(map[string]time.Time),
	}
}

func isLuckyRapidRichesUltimateFish(defID string) bool {
	return defID == "T235"
}

func (m *luckyRapidRichesUltimateManager) getRapidRichesUltimateMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyRapidRichesUltimateManager) tryLuckyRapidRichesUltimateFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(490 * time.Second)
	m.personalCD[p.ID] = now.Add(430 * time.Second)
	m.comboCount[p.ID] = 0
	m.comboTimer[p.ID] = now.Add(3 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyRapidRichesUltimate,
		Payload: map[string]interface{}{
			"event":        "rapid_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"per_mult":     300.0,
			"window_secs":  3,
		},
	})
	g.sendAnnounce(fmt.Sprintf("RAPID RICHES ULTIMATE! %s activated! 3-second combo window! Each kill x300.0!", p.GetDisplayName()), "critical", "#FFD700")
	log.Printf("[LuckyRapidRichesUltimate] %s triggered Rapid Riches Ultimate fish", p.GetDisplayName())

	go func() {
		betCost := float64(p.GetBetDef().BetCost)

		// 模擬 3 秒連擊視窗內的擊破
		time.Sleep(3 * time.Second)

		m.mu.Lock()
		comboHits := m.comboCount[p.ID]
		delete(m.comboCount, p.ID)
		delete(m.comboTimer, p.ID)
		m.mu.Unlock()

		// 計算獎勵
		perMult := 300.0
		totalMult := float64(comboHits) * perMult
		if totalMult < perMult {
			totalMult = perMult // 最少 1 次
			comboHits = 1
		}
		reward := int(totalMult * betCost)
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		isPerfect := comboHits >= 10
		globalBonus := 48.5
		globalDuration := 97

		if isPerfect {
			m.mu.Lock()
			m.boost = &rapidRichesUltimateBoost{
				mult:      globalBonus,
				expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
			}
			m.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyRapidRichesUltimate,
				Payload: map[string]interface{}{
					"event":          "rapid_perfect",
					"combo_hits":     comboHits,
					"per_mult":       perMult,
					"total_mult":     totalMult,
					"reward":         reward,
					"global_bonus":   globalBonus,
					"global_seconds": globalDuration,
				},
			})
			g.sendAnnounce(fmt.Sprintf("PERFECT RAPID RICHES! %s hit %d combos! Total x%.1f! Global x%.1f for %ds!", p.GetDisplayName(), comboHits, totalMult, globalBonus, globalDuration), "critical", "#FFD700")
		} else {
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyRapidRichesUltimate,
				Payload: map[string]interface{}{
					"event":      "rapid_end",
					"combo_hits": comboHits,
					"per_mult":   perMult,
					"total_mult": totalMult,
					"reward":     reward,
				},
			})
		}

		log.Printf("[LuckyRapidRichesUltimate] %s: combo=%d, total_mult=%.1f, perfect=%v", p.GetDisplayName(), comboHits, totalMult, isPerfect)
	}()

	return true
}
