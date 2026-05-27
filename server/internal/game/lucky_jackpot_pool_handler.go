// lucky_jackpot_pool_handler.go — Progressive Jackpot 累積獎池系統
// progressive-jackpot-agent 負責維護
// DAY-313：四層累積獎池（Mini/Minor/Major/Grand）
// 業界參考：Jili Jackpot Fishing 四層 Progressive Jackpot，RTP 97%，最高 888x
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

// JackpotTier 獎池層級
type JackpotTier string

const (
	JackpotMini  JackpotTier = "mini"  // Mini Jackpot：50x 起跳
	JackpotMinor JackpotTier = "minor" // Minor Jackpot：200x 起跳
	JackpotMajor JackpotTier = "major" // Major Jackpot：1000x 起跳
	JackpotGrand JackpotTier = "grand" // Grand Jackpot：5000x 起跳
)

// jackpotPool 累積獎池（內部使用）
type jackpotPool struct {
	mu sync.RWMutex

	// 四層獎池當前倍率（以 bet_cost 倍率計）
	MiniPool  float64 // 起始 50x，每次射擊 +0.01x
	MinorPool float64 // 起始 200x，每次射擊 +0.005x
	MajorPool float64 // 起始 1000x，每次射擊 +0.002x
	GrandPool float64 // 起始 5000x，每次射擊 +0.001x

	// 觸發機率（每次射擊）
	MiniTriggerRate  float64 // 0.005 = 0.5%
	MinorTriggerRate float64 // 0.001 = 0.1%
	MajorTriggerRate float64 // 0.0002 = 0.02%
	GrandTriggerRate float64 // 0.00005 = 0.005%

	// 上次觸發時間（防止連續觸發）
	lastMiniAt  time.Time
	lastMinorAt time.Time
	lastMajorAt time.Time
	lastGrandAt time.Time

	// 冷卻時間（秒）
	miniCooldown  float64
	minorCooldown float64
	majorCooldown float64
	grandCooldown float64
}

func newJackpotPool() *jackpotPool {
	return &jackpotPool{
		MiniPool:  50.0,
		MinorPool: 200.0,
		MajorPool: 1000.0,
		GrandPool: 5000.0,

		MiniTriggerRate:  0.005,
		MinorTriggerRate: 0.001,
		MajorTriggerRate: 0.0002,
		GrandTriggerRate: 0.00005,

		miniCooldown:  30.0,
		minorCooldown: 60.0,
		majorCooldown: 120.0,
		grandCooldown: 300.0,
	}
}

// addContribution 每次射擊時累積獎池
func (jp *jackpotPool) addContribution(betCost int) {
	jp.mu.Lock()
	defer jp.mu.Unlock()
	contribution := float64(betCost)
	jp.MiniPool += contribution * 0.01
	jp.MinorPool += contribution * 0.005
	jp.MajorPool += contribution * 0.002
	jp.GrandPool += contribution * 0.001
}

// tryTrigger 嘗試觸發 Jackpot（每次射擊時呼叫）
func (jp *jackpotPool) tryTrigger() JackpotTier {
	jp.mu.Lock()
	defer jp.mu.Unlock()

	now := time.Now()

	// 從高到低檢查（優先觸發高層級）
	if time.Since(jp.lastGrandAt).Seconds() >= jp.grandCooldown {
		if rand.Float64() < jp.GrandTriggerRate {
			jp.lastGrandAt = now
			return JackpotGrand
		}
	}
	if time.Since(jp.lastMajorAt).Seconds() >= jp.majorCooldown {
		if rand.Float64() < jp.MajorTriggerRate {
			jp.lastMajorAt = now
			return JackpotMajor
		}
	}
	if time.Since(jp.lastMinorAt).Seconds() >= jp.minorCooldown {
		if rand.Float64() < jp.MinorTriggerRate {
			jp.lastMinorAt = now
			return JackpotMinor
		}
	}
	if time.Since(jp.lastMiniAt).Seconds() >= jp.miniCooldown {
		if rand.Float64() < jp.MiniTriggerRate {
			jp.lastMiniAt = now
			return JackpotMini
		}
	}
	return ""
}

// payout 發放獎池並重置
func (jp *jackpotPool) payout(tier JackpotTier) float64 {
	jp.mu.Lock()
	defer jp.mu.Unlock()

	var amount float64
	switch tier {
	case JackpotMini:
		amount = jp.MiniPool
		jp.MiniPool = 50.0
	case JackpotMinor:
		amount = jp.MinorPool
		jp.MinorPool = 200.0
	case JackpotMajor:
		amount = jp.MajorPool
		jp.MajorPool = 1000.0
	case JackpotGrand:
		amount = jp.GrandPool
		jp.GrandPool = 5000.0
	}
	return amount
}

// getSnapshot 取得當前獎池快照
func (jp *jackpotPool) getSnapshot() map[string]float64 {
	jp.mu.RLock()
	defer jp.mu.RUnlock()
	return map[string]float64{
		"mini":  jp.MiniPool,
		"minor": jp.MinorPool,
		"major": jp.MajorPool,
		"grand": jp.GrandPool,
	}
}

// ── Lucky Jackpot Pool Manager ────────────────────────────────

type luckyJackpotPoolManager struct {
	pool             *jackpotPool
	lastBroadcastAt  time.Time
	broadcastInterval float64 // 秒
}

func newLuckyJackpotPoolManager() *luckyJackpotPoolManager {
	return &luckyJackpotPoolManager{
		pool:              newJackpotPool(),
		broadcastInterval: 5.0,
	}
}

// onShot 每次射擊時呼叫（累積獎池 + 嘗試觸發）
func (m *luckyJackpotPoolManager) onShot(g *Game, playerID string, betCost int) {
	m.pool.addContribution(betCost)

	// 定期廣播獎池狀態
	if time.Since(m.lastBroadcastAt).Seconds() >= m.broadcastInterval {
		m.lastBroadcastAt = time.Now()
		m.broadcastPoolUpdate(g)
	}

	tier := m.pool.tryTrigger()
	if tier == "" {
		return
	}

	// 觸發 Jackpot！
	g.mu.Lock()
	p, ok := g.players[playerID]
	if !ok {
		g.mu.Unlock()
		return
	}
	amount := m.pool.payout(tier)
	reward := int(amount) * betCost
	p.Coins += reward
	displayName := p.GetDisplayName()
	g.mu.Unlock()

	log.Printf("[JackpotPool] %s 觸發 %s Jackpot！獎勵 %d 金幣", displayName, tier, reward)

	tierName := jackpotTierName(tier)
	g.broadcast(protocol.Envelope{
		Type: "lucky_jackpot_pool",
		Payload: map[string]interface{}{
			"event":         "jackpot_win",
			"tier":          string(tier),
			"tier_name":     tierName,
			"amount":        amount,
			"reward":        reward,
			"player_id":     playerID,
			"player_name":   displayName,
			"pool_snapshot": m.pool.getSnapshot(),
		},
	})
	g.sendAnnounce(
		fmt.Sprintf("🎰 %s 觸發 %s Jackpot！獎勵 %d 金幣！", displayName, tierName, reward),
		"jackpot",
		jackpotTierColor(tier),
	)
}

// onTargetKill T171-T175 擊破時觸發對應 Jackpot 層級
func (m *luckyJackpotPoolManager) onTargetKill(g *Game, p *Player, targetDefID string) {
	var tier JackpotTier
	switch targetDefID {
	case "T171":
		tier = JackpotMini
	case "T172":
		tier = JackpotMinor
	case "T173":
		tier = JackpotMajor
	case "T174":
		tier = JackpotGrand
	case "T175":
		// T175 隨機觸發一層（機率分佈：Mini 60%，Minor 30%，Major 8%，Grand 2%）
		tiers := []JackpotTier{JackpotMini, JackpotMinor, JackpotMajor, JackpotGrand}
		weights := []int{60, 30, 8, 2}
		tier = tiers[weightedPickIndex(weights)]
	default:
		return
	}

	amount := m.pool.payout(tier)
	betCost := p.GetBetDef().BetCost
	reward := int(amount) * betCost
	p.Coins += reward

	tierName := jackpotTierName(tier)
	log.Printf("[JackpotPool] %s 擊破 %s 觸發 %s Jackpot！獎勵 %d 金幣", p.GetDisplayName(), targetDefID, tier, reward)

	g.broadcast(protocol.Envelope{
		Type: "lucky_jackpot_pool",
		Payload: map[string]interface{}{
			"event":         "jackpot_win",
			"tier":          string(tier),
			"tier_name":     tierName,
			"amount":        amount,
			"reward":        reward,
			"player_id":     p.ID,
			"player_name":   p.GetDisplayName(),
			"target_id":     targetDefID,
			"pool_snapshot": m.pool.getSnapshot(),
		},
	})
	g.sendAnnounce(
		fmt.Sprintf("🎰🏆 %s 擊破 %s 觸發 %s Jackpot！獎勵 %d 金幣！",
			p.GetDisplayName(), targetDefID, tierName, reward),
		"jackpot",
		jackpotTierColor(tier),
	)

	m.broadcastPoolUpdate(g)
}

// broadcastPoolUpdate 廣播獎池當前狀態
func (m *luckyJackpotPoolManager) broadcastPoolUpdate(g *Game) {
	snapshot := m.pool.getSnapshot()
	g.broadcast(protocol.Envelope{
		Type: "lucky_jackpot_pool",
		Payload: map[string]interface{}{
			"event": "pool_update",
			"mini":  snapshot["mini"],
			"minor": snapshot["minor"],
			"major": snapshot["major"],
			"grand": snapshot["grand"],
		},
	})
}

// jackpotTierName 取得層級中文名稱
func jackpotTierName(tier JackpotTier) string {
	switch tier {
	case JackpotMini:
		return "Mini Jackpot"
	case JackpotMinor:
		return "Minor Jackpot"
	case JackpotMajor:
		return "Major Jackpot"
	case JackpotGrand:
		return "Grand Jackpot"
	}
	return "Jackpot"
}

// jackpotTierColor 取得層級顏色（用於公告）
func jackpotTierColor(tier JackpotTier) string {
	switch tier {
	case JackpotMini:
		return "#4CAF50" // 綠色
	case JackpotMinor:
		return "#2196F3" // 藍色
	case JackpotMajor:
		return "#FF9800" // 橙色
	case JackpotGrand:
		return "#FFD700" // 金色
	}
	return "#FFFFFF"
}

// weightedPickIndex 依權重隨機選擇索引
func weightedPickIndex(weights []int) int {
	total := 0
	for _, w := range weights {
		total += w
	}
	r := rand.Intn(total)
	for i, w := range weights {
		r -= w
		if r < 0 {
			return i
		}
	}
	return len(weights) - 1
}

// isLuckyJackpotPoolFish 判斷是否為 Progressive Jackpot 系列魚（T171-T175）
func isLuckyJackpotPoolFish(defID string) bool {
	switch defID {
	case "T171", "T172", "T173", "T174", "T175":
		return true
	}
	return false
}
