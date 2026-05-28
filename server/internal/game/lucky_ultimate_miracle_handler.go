// lucky_ultimate_miracle_handler.go — T210 幸運終極奇蹟魚
// 設計：終極機制，全場 HP 歸零（每個獎勵 ×50.0）
//       觸發後全服 ×35.0 加成 70 秒（新史上最高，超越 T205 的 ×30.0）
//       觸發率：0.001%（史上最稀有）；個人冷卻 300 秒；全服冷卻 360 秒
//       業界依據：終極奇蹟機制 + 2026 最高倍率設計（16888x 吉祥數字）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyUltimateMiracleManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	miracleBoost *ultimateMiracleBoost
}

type ultimateMiracleBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyUltimateMiracleManager() *luckyUltimateMiracleManager {
	return &luckyUltimateMiracleManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyUltimateMiracleFish(defID string) bool {
	return defID == "T210"
}

func (m *luckyUltimateMiracleManager) getUltimateMiracleMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.miracleBoost != nil && time.Now().Before(m.miracleBoost.expiresAt) {
		return m.miracleBoost.mult
	}
	return 1.0
}

func (m *luckyUltimateMiracleManager) tryLuckyUltimateMiracleFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(300 * time.Second)
	m.globalCD = now.Add(360 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_ultimate_miracle",
		Payload: map[string]interface{}{
			"event":        "miracle_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("🌟💫 終極奇蹟！%s 開啟終極奇蹟！全場 HP 歸零！每個獎勵 ×50.0！全服 ×35.0！", p.GetDisplayName()), "critical", "#FFFFFF")
	log.Printf("[LuckyUltimateMiracle] %s 觸發終極奇蹟魚（史上最高全服倍率 ×35.0）", p.GetDisplayName())

	go func() {
		time.Sleep(1500 * time.Millisecond)

		// 終極奇蹟：全場 HP 歸零，每個獎勵 ×50.0（史上最高單次清場倍率）
		hitCount := g.applyUltimateJudgment(p, 50.0)

		// 觸發全服 ×35.0 加成 70 秒（新史上最高全服倍率，超越 T205 的 ×30.0）
		boostMult := 35.0
		boostSecs := 70
		m.mu.Lock()
		m.miracleBoost = &ultimateMiracleBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_ultimate_miracle",
			Payload: map[string]interface{}{
				"event":        "miracle_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"hit_count":    hitCount,
				"reward_mult":  50.0,
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
			},
		})
		g.sendAnnounce(fmt.Sprintf("🌟💫 終極奇蹟完成！%s 清場 %d 個目標（×50.0）！全服 ×%.1f 加成 %d 秒（新史上最高）！", p.GetDisplayName(), hitCount, boostMult, boostSecs), "critical", "#FFFFFF")
		log.Printf("[LuckyUltimateMiracle] 終極奇蹟完成！命中 %d 個目標，全服 ×%.1f 加成 %d 秒（新史上最高）", hitCount, boostMult, boostSecs)
	}()
	return true
}
