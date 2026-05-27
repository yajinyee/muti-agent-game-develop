// lucky_risk_level_handler.go — T184 幸運風險等級魚
// 業界依據：BGaming Fishing Club 2（5 risk levels, max x3000）
// 設計：擊破後選擇 5 個風險等級（低 ×5.0 / 中 ×20.0 / 高 ×100.0 / 極高 ×500.0 / 最高 ×3000.0）
//       最高等級（×3000.0）→ 全服 ×17.5 加成 36 秒
//       個人冷卻 100 秒；全服冷卻 155 秒
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

// 5 個風險等級定義
var riskLevels = []struct {
	Name       string
	Mult       float64
	Probability float64 // 觸發機率
	Color      string
}{
	{"低風險", 5.0, 0.40, "#4CAF50"},
	{"中風險", 20.0, 0.30, "#FF9800"},
	{"高風險", 100.0, 0.18, "#F44336"},
	{"極高風險", 500.0, 0.10, "#9C27B0"},
	{"最高風險", 3000.0, 0.02, "#FFD700"},
}

type luckyRiskLevelManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *riskLevelPerfectBoost
}

type riskLevelPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyRiskLevelManager() *luckyRiskLevelManager {
	return &luckyRiskLevelManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyRiskLevelFish(defID string) bool {
	return defID == "T184"
}

func (m *luckyRiskLevelManager) getRiskLevelMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyRiskLevelManager) rollRiskLevel() (string, float64, string) {
	r := rand.Float64()
	cumulative := 0.0
	for _, rl := range riskLevels {
		cumulative += rl.Probability
		if r < cumulative {
			return rl.Name, rl.Mult, rl.Color
		}
	}
	return riskLevels[0].Name, riskLevels[0].Mult, riskLevels[0].Color
}

func (m *luckyRiskLevelManager) tryLuckyRiskLevelFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(100 * time.Second)
	m.globalCD = now.Add(155 * time.Second)
	m.mu.Unlock()

	// 廣播開始，顯示 5 個等級選項
	g.broadcast(protocol.Envelope{
		Type: "lucky_risk_level",
		Payload: map[string]interface{}{
			"event":        "risk_level_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"levels": []map[string]interface{}{
				{"name": "低風險", "mult": 5.0, "prob": "40%", "color": "#4CAF50"},
				{"name": "中風險", "mult": 20.0, "prob": "30%", "color": "#FF9800"},
				{"name": "高風險", "mult": 100.0, "prob": "18%", "color": "#F44336"},
				{"name": "極高風險", "mult": 500.0, "prob": "10%", "color": "#9C27B0"},
				{"name": "最高風險", "mult": 3000.0, "prob": "2%", "color": "#FFD700"},
			},
		},
	})
	g.sendAnnounce(fmt.Sprintf("🎰⚡ 風險等級！%s 觸發風險選擇！最高 ×3000！", p.GetDisplayName()), "special", "#E65100")
	log.Printf("[LuckyRiskLevel] %s 觸發風險等級魚", p.GetDisplayName())

	go func() {
		// 1.5 秒後揭曉結果（模擬選擇動畫）
		time.Sleep(1500 * time.Millisecond)

		riskName, riskMult, riskColor := m.rollRiskLevel()

		// 套用風險倍率到全場目標
		hitCount := 0
		g.mu.Lock()
		for _, t := range g.targets {
			if t.HP > 0 && t.Def.Type != data.TypeBoss {
				damage := int(float64(t.HP) * 0.40)
				if damage < 1 {
					damage = 1
				}
				t.HP -= damage
				if t.HP <= 0 {
					t.HP = 0
				}
				hitCount++
			}
		}
		g.mu.Unlock()

		// 最高風險（×3000）→ 全服加成
		if riskMult >= 3000.0 {
			boostMult := 17.5
			boostSecs := 36
			m.mu.Lock()
			m.perfectBoost = &riskLevelPerfectBoost{
				mult:      boostMult,
				expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
			}
			m.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: "lucky_risk_level",
				Payload: map[string]interface{}{
					"event":        "risk_level_jackpot",
					"risk_name":    riskName,
					"risk_mult":    riskMult,
					"risk_color":   riskColor,
					"hit_count":    hitCount,
					"boost_mult":   boostMult,
					"boost_secs":   boostSecs,
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
				},
			})
			g.sendAnnounce(fmt.Sprintf("🎰🏆 最高風險！%s 抽中 ×3000！全服 ×%.1f 加成 %d 秒！", p.GetDisplayName(), boostMult, boostSecs), "critical", "#E65100")
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_risk_level",
				Payload: map[string]interface{}{
					"event":        "risk_level_result",
					"risk_name":    riskName,
					"risk_mult":    riskMult,
					"risk_color":   riskColor,
					"hit_count":    hitCount,
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
				},
			})
			g.sendAnnounce(fmt.Sprintf("🎰 風險結果：%s 抽中【%s】×%.1f！", p.GetDisplayName(), riskName, riskMult), "normal", riskColor)
		}
	}()
	return true
}
