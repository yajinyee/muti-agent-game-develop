// lucky_path_fish_handler.go — T208 幸運路徑魚
// 設計：Fish Road 機制（fishroad.eu）
//       路徑越遠倍率越高，每次擊破推進路徑（最高 20,000x）
//       觸發後全服 ×33.0 加成 66 秒（超越 T207 的 ×32.0）
//       觸發率：0.002%（極稀有）；個人冷卻 270 秒；全服冷卻 330 秒
//       業界依據：Fish Road「路徑越遠倍率越高，最高 20,000x」
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

// 路徑倍率表（每一步的倍率）
var pathMultipliers = []float64{
	1.5, 2.0, 3.0, 5.0, 8.0, 12.0, 20.0, 35.0, 60.0, 100.0,
	150.0, 250.0, 400.0, 700.0, 1200.0, 2000.0, 3500.0, 6000.0, 10000.0, 20000.0,
}

type luckyPathFishManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	pathBoost  *pathFishBoost
	pathStep   int     // 當前路徑步數（0-19）
	pathMult   float64 // 當前路徑倍率
	isActive   bool
	triggerID  string
}

type pathFishBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyPathFishManager() *luckyPathFishManager {
	return &luckyPathFishManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyPathFish(defID string) bool {
	return defID == "T208"
}

func (m *luckyPathFishManager) getPathFishMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.pathBoost != nil && time.Now().Before(m.pathBoost.expiresAt) {
		return m.pathBoost.mult
	}
	return 1.0
}

func (m *luckyPathFishManager) onKillDuringPath(g *Game, p *Player) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.isActive || m.triggerID != p.ID {
		return 1.0
	}
	if m.pathStep < len(pathMultipliers)-1 {
		m.pathStep++
	}
	m.pathMult = pathMultipliers[m.pathStep]

	g.broadcast(protocol.Envelope{
		Type: "lucky_path_fish",
		Payload: map[string]interface{}{
			"event":      "path_advance",
			"trigger_id": p.ID,
			"step":       m.pathStep,
			"mult":       m.pathMult,
			"max_steps":  len(pathMultipliers),
		},
	})
	return m.pathMult
}

func (m *luckyPathFishManager) tryLuckyPathFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(270 * time.Second)
	m.globalCD = now.Add(330 * time.Second)
	m.pathStep = 0
	m.pathMult = pathMultipliers[0]
	m.isActive = true
	m.triggerID = p.ID
	m.mu.Unlock()

	pathSecs := 40
	g.broadcast(protocol.Envelope{
		Type: "lucky_path_fish",
		Payload: map[string]interface{}{
			"event":        "path_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"path_secs":    pathSecs,
			"max_mult":     20000.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🛤️✨ 路徑魚！%s 開啟路徑！每次擊破推進路徑，最高 ×20,000！", p.GetDisplayName()), "critical", "#00FFFF")
	log.Printf("[LuckyPathFish] %s 觸發路徑魚（40 秒路徑，最高 ×20,000）", p.GetDisplayName())

	go func() {
		time.Sleep(time.Duration(pathSecs) * time.Second)

		m.mu.Lock()
		finalStep := m.pathStep
		finalMult := m.pathMult
		m.isActive = false
		m.mu.Unlock()

		// 觸發全服 ×33.0 加成 66 秒
		globalBoostMult := 33.0
		globalBoostSecs := 66
		m.mu.Lock()
		m.pathBoost = &pathFishBoost{
			mult:      globalBoostMult,
			expiresAt: time.Now().Add(time.Duration(globalBoostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_path_fish",
			Payload: map[string]interface{}{
				"event":        "path_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"final_step":   finalStep,
				"final_mult":   finalMult,
				"global_mult":  globalBoostMult,
				"global_secs":  globalBoostSecs,
			},
		})
		g.sendAnnounce(fmt.Sprintf("🛤️🌟 路徑完成！%s 到達第 %d 步（×%.1f）！全服 ×%.1f 加成 %d 秒！", p.GetDisplayName(), finalStep+1, finalMult, globalBoostMult, globalBoostSecs), "critical", "#00FFFF")
		log.Printf("[LuckyPathFish] 路徑完成！第 %d 步（×%.1f），全服 ×%.1f 加成 %d 秒", finalStep+1, finalMult, globalBoostMult, globalBoostSecs)
	}()
	return true
}
