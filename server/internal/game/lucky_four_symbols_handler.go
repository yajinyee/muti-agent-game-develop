// lucky_four_symbols_handler.go — 幸運四象大獎魚系統（DAY-283）
// 業界依據：Jackpot Fishing by Jili「四層累進大獎」機制 + 中華文化「四象」主題
// 四象：青龍（×2.0）/ 白虎（×5.0）/ 朱雀（×12.0）/ 玄武（×30.0）
//
// 設計：
//   - Server 維護四象大獎池，每次任何玩家擊破任何目標貢獻 0.5% 到大獎池
//   - 擊破 T241 後，依機率抽取四象層級：青龍 60% / 白虎 25% / 朱雀 12% / 玄武 3%
//   - 大獎金額 = 當前大獎池 × 層級比例（青龍 10% / 白虎 25% / 朱雀 50% / 玄武 100%）
//   - 玄武大獎觸發後，大獎池重置為基礎值（30000）
//   - 全服廣播大獎結果；玄武大獎全服最高優先公告
//   - 個人冷卻 30 秒；全服冷卻 50 秒
//
// 設計差異：
//   - 與累積大獎池（T222，全服共同累積到爆發）不同，四象大獎是「個人觸發，依機率抽取層級」
//   - 「四象主題（青龍/白虎/朱雀/玄武）」讓大獎有文化感，比 Mini/Minor/Major/Grand 更有特色
//   - 「玄武 100% 大獎池」讓玩家有「要是觸發玄武就賺大了」的動力
//   - 「每次擊破貢獻 0.5%」讓大獎池持續增長
//   - 「玄武大獎全服最高優先公告」製造羨慕感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyFourSymbolsPersonalCD = 30 * time.Second // 個人冷卻
	LuckyFourSymbolsGlobalCD   = 50 * time.Second // 全服冷卻

	// 大獎池基礎值（重置後的初始值）
	FourSymbolsPoolBase = 30000

	// 每次擊破貢獻比例（0.5%）
	FourSymbolsContribRate = 0.005

	// 大獎觸發機率
	FourSymbolsQinglongChance = 0.60 // 青龍 60%
	FourSymbolsBaihuChance    = 0.25 // 白虎 25%
	FourSymbolsZhuqueChance   = 0.12 // 朱雀 12%
	FourSymbolsXuanwuChance   = 0.03 // 玄武 3%

	// 大獎金額比例（佔大獎池的比例）
	FourSymbolsQinglongPayout = 0.10 // 青龍取 10%
	FourSymbolsBaihuPayout    = 0.25 // 白虎取 25%
	FourSymbolsZhuquePayout   = 0.50 // 朱雀取 50%
	FourSymbolsXuanwuPayout   = 1.00 // 玄武取 100%（清空大獎池）

	// 大獎池廣播間隔（每 30 秒廣播一次大獎池狀態）
	FourSymbolsPoolBroadcastInterval = 30 * time.Second
)

// fourSymbolTier 四象層級定義
type fourSymbolTier struct {
	id      string
	name    string
	chance  float64
	payout  float64
	color   string
	emoji   string
}

var fourSymbolTiers = []fourSymbolTier{
	{id: "xuanwu",   name: "玄武", chance: FourSymbolsXuanwuChance,   payout: FourSymbolsXuanwuPayout,   color: "#000080", emoji: "🐢"},
	{id: "zhuque",   name: "朱雀", chance: FourSymbolsZhuqueChance,   payout: FourSymbolsZhuquePayout,   color: "#FF0000", emoji: "🦅"},
	{id: "baihu",    name: "白虎", chance: FourSymbolsBaihuChance,    payout: FourSymbolsBaihuPayout,    color: "#C0C0C0", emoji: "🐯"},
	{id: "qinglong", name: "青龍", chance: FourSymbolsQinglongChance, payout: FourSymbolsQinglongPayout, color: "#00AA00", emoji: "🐉"},
}

// luckyFourSymbolsManager 幸運四象大獎魚管理器
type luckyFourSymbolsManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 四象大獎池
	pool int
}

func newLuckyFourSymbolsManager() *luckyFourSymbolsManager {
	return &luckyFourSymbolsManager{
		personalCooldowns: make(map[string]time.Time),
		pool:              FourSymbolsPoolBase,
	}
}

// isLuckyFourSymbolsFish 判斷是否為幸運四象大獎魚
func isLuckyFourSymbolsFish(defID string) bool {
	return defID == "T241"
}

// contributeToFourSymbolsPool 每次擊破貢獻到大獎池（供 handleKill 使用）
func (m *luckyFourSymbolsManager) contributeToFourSymbolsPool(reward int) {
	if reward <= 0 {
		return
	}
	contrib := int(float64(reward) * FourSymbolsContribRate)
	if contrib < 1 {
		contrib = 1
	}
	m.mu.Lock()
	m.pool += contrib
	m.mu.Unlock()
}

// getPoolSize 取得大獎池大小（thread-safe）
func (m *luckyFourSymbolsManager) getPoolSize() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.pool
}

// rollFourSymbolTier 抽取四象層級
func rollFourSymbolTier() *fourSymbolTier {
	r := rand.Float64()
	cumulative := 0.0
	for i := range fourSymbolTiers {
		cumulative += fourSymbolTiers[i].chance
		if r < cumulative {
			return &fourSymbolTiers[i]
		}
	}
	return &fourSymbolTiers[len(fourSymbolTiers)-1] // fallback: 青龍
}

// tryLuckyFourSymbolsFish 擊破 T241 後觸發四象大獎（供 handleKill 使用）
func (g *Game) tryLuckyFourSymbolsFish(p *player.Player) {
	mgr := g.LuckyFourSymbols
	mgr.mu.Lock()

	now := time.Now()

	// 全服冷卻檢查
	if now.Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && now.Before(cd) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = now.Add(LuckyFourSymbolsPersonalCD)
	mgr.globalCooldownUntil = now.Add(LuckyFourSymbolsGlobalCD)

	// 抽取四象層級
	tier := rollFourSymbolTier()

	// 計算大獎金額
	reward := int(float64(mgr.pool) * tier.payout)

	// 玄武大獎：重置大獎池
	if tier.id == "xuanwu" {
		mgr.pool = FourSymbolsPoolBase
	}

	poolSize := mgr.pool
	mgr.mu.Unlock()

	// 給予玩家獎勵
	if reward > 0 {
		p.AddCoins(reward)
	}

	log.Printf("[FourSymbols] player=%s tier=%s(%s) reward=%d pool=%d",
		p.ID, tier.id, tier.name, reward, poolSize)

	// 全服廣播：四象大獎觸發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyFourSymbols,
		Payload: ws.LuckyFourSymbolsPayload{
			Event:      "symbol_trigger",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Symbol:     tier.id,
			SymbolName: tier.name,
			Reward:     reward,
			PoolSize:   poolSize,
		},
	})

	// 玄武大獎：額外全服廣播 + 最高優先公告
	if tier.id == "xuanwu" {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyFourSymbols,
			Payload: ws.LuckyFourSymbolsPayload{
				Event:      "symbol_xuanwu",
				PlayerName: p.DisplayName,
				Reward:     reward,
			},
		})
		ann := g.Announce.Create(announce.EventLuckyFourSymbols, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🐢🐢🐢 %s 觸發玄武大獎！獲得 %d 金幣！大獎池已重置！",
				p.DisplayName, reward),
			"color": "#000080",
		})
		g.broadcastAnnouncement(ann)
	} else {
		ann := g.Announce.Create(announce.EventLuckyFourSymbols, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("%s %s 觸發%s大獎！獲得 %d 金幣！",
				tier.emoji, p.DisplayName, tier.name, reward),
			"color": tier.color,
		})
		g.broadcastAnnouncement(ann)
	}
}

// startFourSymbolsPoolBroadcast 啟動四象大獎池定期廣播 goroutine
func (g *Game) startFourSymbolsPoolBroadcast() {
	go func() {
		ticker := time.NewTicker(FourSymbolsPoolBroadcastInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				poolSize := g.LuckyFourSymbols.getPoolSize()
				g.Hub.Broadcast(&ws.Message{
					Type: ws.MsgLuckyFourSymbols,
					Payload: ws.LuckyFourSymbolsPayload{
						Event:    "symbol_pool_update",
						PoolSize: poolSize,
					},
				})
			case <-g.stopCh:
				return
			}
		}
	}()
}
