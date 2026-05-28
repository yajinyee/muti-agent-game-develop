// lucky_chain_eel_handler.go — T209 幸運連鎖電鰻魚
// 設計：Royal Fishing 紫粉色電鰻升級版（royal-fishing.co.uk）
//       連鎖電擊 8 條魚，每條 ×40.0，全部命中 → 完美連鎖
//       觸發後全服 ×34.0 加成 68 秒（超越 T208 的 ×33.0）
//       觸發率：0.002%（極稀有）；個人冷卻 280 秒；全服冷卻 340 秒
//       業界依據：Royal Fishing「Purple/Pink Lightning Eel chain reaction」升級版
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyChainEelManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	eelBoost   *chainEelBoost
}

type chainEelBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyChainEelManager() *luckyChainEelManager {
	return &luckyChainEelManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyChainEelFish(defID string) bool {
	return defID == "T209"
}

func (m *luckyChainEelManager) getChainEelMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.eelBoost != nil && time.Now().Before(m.eelBoost.expiresAt) {
		return m.eelBoost.mult
	}
	return 1.0
}

func (m *luckyChainEelManager) tryLuckyChainEelFish(g *Game, p *Player) bool {
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

	chainCount := 8
	rewardMult := 40.0

	g.broadcast(protocol.Envelope{
		Type: "lucky_chain_eel",
		Payload: map[string]interface{}{
			"event":        "eel_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"chain_count":  chainCount,
			"reward_mult":  rewardMult,
		},
	})
	g.sendAnnounce(fmt.Sprintf("⚡💜 連鎖電鰻！%s 觸發連鎖電鰻！8 條魚連鎖電擊，每條 ×40.0！", p.GetDisplayName()), "critical", "#CC00FF")
	log.Printf("[LuckyChainEel] %s 觸發連鎖電鰻魚（8 條連鎖，每條 ×40.0）", p.GetDisplayName())

	go func() {
		time.Sleep(800 * time.Millisecond)

		// 連鎖電擊 8 條魚，每條 ×40.0
		hitCount := g.applyUltimateJudgment(p, rewardMult)
		if hitCount > chainCount {
			hitCount = chainCount
		}

		isPerfect := hitCount >= chainCount

		// 觸發全服 ×34.0 加成 68 秒
		globalBoostMult := 34.0
		globalBoostSecs := 68
		m.mu.Lock()
		m.eelBoost = &chainEelBoost{
			mult:      globalBoostMult,
			expiresAt: time.Now().Add(time.Duration(globalBoostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_chain_eel",
			Payload: map[string]interface{}{
				"event":        "eel_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"hit_count":    hitCount,
				"reward_mult":  rewardMult,
				"is_perfect":   isPerfect,
				"global_mult":  globalBoostMult,
				"global_secs":  globalBoostSecs,
			},
		})

		if isPerfect {
			g.sendAnnounce(fmt.Sprintf("⚡💜 完美連鎖！%s 電擊 %d 條魚（×%.1f）！全服 ×%.1f 加成 %d 秒！", p.GetDisplayName(), hitCount, rewardMult, globalBoostMult, globalBoostSecs), "critical", "#CC00FF")
		} else {
			g.sendAnnounce(fmt.Sprintf("⚡💜 連鎖電鰻完成！%s 電擊 %d 條魚！全服 ×%.1f 加成 %d 秒！", p.GetDisplayName(), hitCount, globalBoostMult, globalBoostSecs), "critical", "#CC00FF")
		}
		log.Printf("[LuckyChainEel] 連鎖電鰻完成！命中 %d 條，全服 ×%.1f 加成 %d 秒", hitCount, globalBoostMult, globalBoostSecs)
	}()
	return true
}
