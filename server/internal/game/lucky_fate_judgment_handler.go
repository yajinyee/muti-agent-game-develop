// lucky_fate_judgment_handler.go — T203 幸運命運審判魚
// 設計：命運審判，隨機選 5 個目標，每個倍率 ×50-×500（加權隨機）
//       全部 5 個擊破 → 完美審判全服 ×28.0 加成 56 秒
//       觸發率：0.006%；個人冷卻 180 秒；全服冷卻 240 秒
//       業界依據：Fishing Fortune「Fate Judgment」升級版 + 2026 命運審判機制
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyFateJudgmentManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	judgmentBoost *fateJudgmentBoost
	// 命運目標：targetInstanceID -> 倍率
	fateTargets map[string]float64
}

type fateJudgmentBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyFateJudgmentManager() *luckyFateJudgmentManager {
	return &luckyFateJudgmentManager{
		personalCD:  make(map[string]time.Time),
		fateTargets: make(map[string]float64),
	}
}

func isLuckyFateJudgmentFish(defID string) bool {
	return defID == "T203"
}

func (m *luckyFateJudgmentManager) getFateJudgmentMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.judgmentBoost != nil && time.Now().Before(m.judgmentBoost.expiresAt) {
		return m.judgmentBoost.mult
	}
	return 1.0
}

func (m *luckyFateJudgmentManager) getFateTargetMult(instanceID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if mult, ok := m.fateTargets[instanceID]; ok {
		return mult
	}
	return 1.0
}

func (m *luckyFateJudgmentManager) onFateTargetKilled(instanceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.fateTargets, instanceID)
}

// 命運倍率加權隨機（×50-×500）
func rollFateMultiplier() float64 {
	r := rand.Float64()
	switch {
	case r < 0.40:
		return 50.0
	case r < 0.65:
		return 100.0
	case r < 0.82:
		return 200.0
	case r < 0.93:
		return 300.0
	case r < 0.98:
		return 400.0
	default:
		return 500.0
	}
}

func (m *luckyFateJudgmentManager) tryLuckyFateJudgmentFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(180 * time.Second)
	m.globalCD = now.Add(240 * time.Second)
	m.mu.Unlock()

	// 選取 5 個隨機目標，分配命運倍率
	g.mu.Lock()
	var selectedIDs []string
	var selectedMults []float64
	count := 0
	for id := range g.targets {
		if count >= 5 {
			break
		}
		mult := rollFateMultiplier()
		selectedIDs = append(selectedIDs, id)
		selectedMults = append(selectedMults, mult)
		count++
	}
	g.mu.Unlock()

	m.mu.Lock()
	for i, id := range selectedIDs {
		m.fateTargets[id] = selectedMults[i]
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_fate_judgment",
		Payload: map[string]interface{}{
			"event":        "fate_judgment_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"target_count": len(selectedIDs),
			"target_ids":   selectedIDs,
			"target_mults": selectedMults,
		},
	})
	g.sendAnnounce(fmt.Sprintf("⚖️🌟 命運審判！%s 標記 %d 個命運目標！最高 ×500！全服 ×28.0！", p.GetDisplayName(), len(selectedIDs)), "critical", "#FFD700")
	log.Printf("[LuckyFateJudgment] %s 觸發命運審判魚（%d 個命運目標）", p.GetDisplayName(), len(selectedIDs))

	go func() {
		// 30 秒後清除未擊破的命運目標，觸發全服加成
		time.Sleep(30 * time.Second)

		boostMult := 28.0
		boostSecs := 56
		m.mu.Lock()
		m.fateTargets = make(map[string]float64) // 清除剩餘命運目標
		m.judgmentBoost = &fateJudgmentBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_fate_judgment",
			Payload: map[string]interface{}{
				"event":        "fate_judgment_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
			},
		})
		log.Printf("[LuckyFateJudgment] 命運審判完成！全服 ×%.1f 加成 %d 秒", boostMult, boostSecs)
	}()
	return true
}
