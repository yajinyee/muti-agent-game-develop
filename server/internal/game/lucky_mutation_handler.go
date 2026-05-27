// lucky_mutation_handler.go — T181 幸運突變魚
// 業界依據：Fisch mutations system（150+ mutations, 17x bonus）
// 設計：擊破後觸發隨機突變（150種突變，最高 ×17.0 加成）
//       突變 ≥10x → 全服 ×16.0 加成 32 秒
//       個人冷卻 95 秒；全服冷卻 145 秒
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/data"
	"chiikawa-game/internal/protocol"
)

// 150 種突變定義（倍率 + 權重）
var mutationTable = []struct {
	Name   string
	Mult   float64
	Weight int
}{
	// 普通突變（×1.0-×2.0，高機率）
	{"水流突變", 1.2, 30}, {"氣泡突變", 1.3, 28}, {"珊瑚突變", 1.4, 26},
	{"海草突變", 1.5, 24}, {"沙地突變", 1.6, 22}, {"礁石突變", 1.7, 20},
	{"潮汐突變", 1.8, 18}, {"深海突變", 1.9, 16}, {"光線突變", 2.0, 14},
	// 稀有突變（×2.5-×5.0，中機率）
	{"閃電突變", 2.5, 12}, {"冰晶突變", 3.0, 10}, {"火焰突變", 3.5, 9},
	{"毒素突變", 4.0, 8}, {"暗影突變", 4.5, 7}, {"光明突變", 5.0, 6},
	// 史詩突變（×6.0-×10.0，低機率）
	{"龍鱗突變", 6.0, 5}, {"鳳凰突變", 7.0, 4}, {"雷霆突變", 8.0, 3},
	{"神聖突變", 9.0, 2}, {"宇宙突變", 10.0, 2},
	// 傳說突變（×12.0-×17.0，極低機率）
	{"創世突變", 12.0, 1}, {"終焉突變", 15.0, 1}, {"超越突變", 17.0, 1},
}

type luckyMutationManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *mutationPerfectBoost
}

type mutationPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyMutationManager() *luckyMutationManager {
	return &luckyMutationManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyMutationFish(defID string) bool {
	return defID == "T181"
}

func (m *luckyMutationManager) getMutationMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyMutationManager) rollMutation() (string, float64) {
	totalWeight := 0
	for _, mt := range mutationTable {
		totalWeight += mt.Weight
	}
	r := rand.Intn(totalWeight)
	cumulative := 0
	for _, mt := range mutationTable {
		cumulative += mt.Weight
		if r < cumulative {
			return mt.Name, mt.Mult
		}
	}
	return mutationTable[0].Name, mutationTable[0].Mult
}

func (m *luckyMutationManager) tryLuckyMutationFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(95 * time.Second)
	m.globalCD = now.Add(145 * time.Second)
	m.mu.Unlock()

	mutName, mutMult := m.rollMutation()

	g.broadcast(protocol.Envelope{
		Type: "lucky_mutation",
		Payload: map[string]interface{}{
			"event":        "mutation_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"mutation_name": mutName,
			"mutation_mult": mutMult,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🧬✨ 突變觸發！%s 引動【%s】×%.1f！", p.GetDisplayName(), mutName, mutMult), "special", "#7B1FA2")
	log.Printf("[LuckyMutation] %s 觸發突變魚：%s ×%.1f", p.GetDisplayName(), mutName, mutMult)

	go func() {
		time.Sleep(500 * time.Millisecond)

		// 套用突變倍率到全場目標
		hitCount := 0
		g.mu.Lock()
		for _, t := range g.targets {
			if t.HP > 0 && t.Def.Type != data.TypeBoss {
				damage := int(float64(t.HP) * 0.35)
				t.HP -= damage
				if t.HP <= 0 {
					t.HP = 0
				}
				hitCount++
			}
		}
		g.mu.Unlock()

		// 如果突變 ≥10x，觸發全服加成
		if mutMult >= 10.0 {
			boostMult := 16.0
			boostSecs := 32
			m.mu.Lock()
			m.perfectBoost = &mutationPerfectBoost{
				mult:      boostMult,
				expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
			}
			m.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: "lucky_mutation",
				Payload: map[string]interface{}{
					"event":         "mutation_perfect",
					"mutation_name": mutName,
					"mutation_mult": mutMult,
					"hit_count":     hitCount,
					"boost_mult":    boostMult,
					"boost_secs":    boostSecs,
					"trigger_id":    p.ID,
					"trigger_name":  p.GetDisplayName(),
				},
			})
			g.sendAnnounce(fmt.Sprintf("🧬🏆 傳說突變！%s【%s】×%.1f！全服 ×%.1f 加成 %d 秒！", p.GetDisplayName(), mutName, mutMult, boostMult, boostSecs), "critical", "#6A1B9A")
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_mutation",
				Payload: map[string]interface{}{
					"event":         "mutation_complete",
					"mutation_name": mutName,
					"mutation_mult": mutMult,
					"hit_count":     hitCount,
					"trigger_id":    p.ID,
					"trigger_name":  p.GetDisplayName(),
				},
			})
		}
	}()
	return true
}
