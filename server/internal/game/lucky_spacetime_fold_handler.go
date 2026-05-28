// lucky_spacetime_fold_handler.go — T194 幸運時空折疊魚
// 設計：時空折疊 20 秒，所有目標倍率 ×3.0，射擊速度 ×2.0
//       觸發後全服 ×21.0 加成 42 秒（新最高）
//       觸發率：0.025%；個人冷卻 135 秒；全服冷卻 195 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckySpacetimeFoldManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	foldBoost  *spacetimeFoldPerfectBoost
	// 折疊期間的倍率加成（目標倍率 ×3.0）
	foldActive    bool
	foldExpiresAt time.Time
}

type spacetimeFoldPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckySpacetimeFoldManager() *luckySpacetimeFoldManager {
	return &luckySpacetimeFoldManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckySpacetimeFoldFish(defID string) bool {
	return defID == "T194"
}

func (m *luckySpacetimeFoldManager) getSpacetimeFoldMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	// 折疊期間：目標倍率 ×3.0
	if m.foldActive && time.Now().Before(m.foldExpiresAt) {
		return 3.0
	}
	// 折疊後全服加成
	if m.foldBoost != nil && time.Now().Before(m.foldBoost.expiresAt) {
		return m.foldBoost.mult
	}
	return 1.0
}

func (m *luckySpacetimeFoldManager) tryLuckySpacetimeFoldFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(135 * time.Second)
	m.globalCD = now.Add(195 * time.Second)

	// 啟動折疊效果
	m.foldActive = true
	m.foldExpiresAt = now.Add(20 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_spacetime_fold",
		Payload: map[string]interface{}{
			"event":          "spacetime_fold_start",
			"trigger_id":     p.ID,
			"trigger_name":   p.GetDisplayName(),
			"duration":       20,
			"target_mult":    3.0,
			"fire_rate_mult": 2.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🌀⏰ 時空折疊！%s 折疊時空！20 秒內目標倍率 ×3.0！射擊速度 ×2.0！", p.GetDisplayName()), "critical", "#4A148C")
	log.Printf("[LuckySpacetimeFold] %s 觸發時空折疊魚", p.GetDisplayName())

	go func() {
		time.Sleep(20 * time.Second)

		m.mu.Lock()
		m.foldActive = false
		m.mu.Unlock()

		// 折疊結束後觸發全服 ×21.0 加成 42 秒（新最高）
		boostMult := 21.0
		boostSecs := 42
		m.mu.Lock()
		m.foldBoost = &spacetimeFoldPerfectBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_spacetime_fold",
			Payload: map[string]interface{}{
				"event":        "spacetime_fold_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
			},
		})
		g.sendAnnounce(fmt.Sprintf("🌀🏆 時空折疊結束！%s 全服 ×%.1f 加成 %d 秒！（新最高）",
			p.GetDisplayName(), boostMult, boostSecs), "critical", "#6A1B9A")
	}()
	return true
}
