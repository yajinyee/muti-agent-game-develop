// lucky_progressive_jackpot_handler.go — 幸運累積大獎池魚系統（DAY-262）
// 業界原創「全服累積大獎池+貢獻比例分配+大獎池爆發」機制
//
// 設計：全服所有玩家每次擊破任何目標時，自動累積「大獎池」（每次 +1%）：
//   - 擊破 T220 後，大獎池立即「爆發」：按貢獻比例分配給所有玩家
//   - 貢獻比例 = 玩家本局累積擊破次數 / 全服總擊破次數
//   - 大獎池最小值 = 100（確保有意義的爆發）
//   - 大獎池上限 = 10000（防止無限累積）
//   - 個人冷卻 60 秒；全服冷卻 90 秒
//
// 設計差異：
//   - 與全服充能（T214，合作達到目標數爆發）不同，累積大獎池是「持續累積」，
//     讓玩家有「每一槍都在累積大獎池」的動力
//   - 「貢獻比例分配」讓打得多的玩家獲得更多，公平且有激勵效果
//   - 「大獎池即時顯示」讓玩家看到「大獎池現在有多少」，製造「快要爆發了」的期待感
//   - 「擊破 T220 立即爆發」讓玩家有「要趕快找到 T220」的動機
//   - 「全服廣播大獎池金額」讓所有玩家都知道「現在大獎池有多少」，製造社交討論感
//   - 「貢獻排行榜廣播」讓玩家看到「誰貢獻最多」，製造競爭感
//   - 業界依據：Progressive Jackpot 是 2026 年捕魚機最熱門的留存機制（Fishing Fortune, Ocean King）
package game

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyProgressiveJackpotPersonalCD  = 60 * time.Second // 個人冷卻
	LuckyProgressiveJackpotGlobalCD    = 90 * time.Second // 全服冷卻
	LuckyProgressiveJackpotContribRate = 0.01             // 每次擊破貢獻 1% 到大獎池
	LuckyProgressiveJackpotMinPool     = 100              // 大獎池最小值
	LuckyProgressiveJackpotMaxPool     = 10000            // 大獎池上限
	LuckyProgressiveJackpotBroadcastInterval = 10 * time.Second // 大獎池廣播間隔
)

// jackpotContribution 玩家貢獻記錄
type jackpotContribution struct {
	playerID   string
	playerName string
	kills      int // 本局累積擊破次數
}

// luckyProgressiveJackpotManager 幸運累積大獎池魚管理器
type luckyProgressiveJackpotManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 大獎池當前金額
	pool int

	// 玩家貢獻記錄（playerID → contribution）
	contributions map[string]*jackpotContribution

	// 上次廣播時間
	lastBroadcastAt time.Time
}

func newLuckyProgressiveJackpotManager() *luckyProgressiveJackpotManager {
	return &luckyProgressiveJackpotManager{
		personalCooldowns: make(map[string]time.Time),
		pool:              LuckyProgressiveJackpotMinPool,
		contributions:     make(map[string]*jackpotContribution),
		lastBroadcastAt:   time.Now(),
	}
}

// isLuckyProgressiveJackpotFish 判斷是否為幸運累積大獎池魚
func isLuckyProgressiveJackpotFish(defID string) bool {
	return defID == "T220"
}

// accumulateJackpot 每次擊破時累積大獎池（由 handleKill 呼叫）
func (m *luckyProgressiveJackpotManager) accumulateJackpot(p *player.Player, killReward int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 累積大獎池（1% 的擊破獎勵）
	contrib := int(float64(killReward) * LuckyProgressiveJackpotContribRate)
	if contrib < 1 {
		contrib = 1
	}
	m.pool += contrib
	if m.pool > LuckyProgressiveJackpotMaxPool {
		m.pool = LuckyProgressiveJackpotMaxPool
	}

	// 記錄玩家貢獻
	if c, ok := m.contributions[p.ID]; ok {
		c.kills++
	} else {
		m.contributions[p.ID] = &jackpotContribution{
			playerID:   p.ID,
			playerName: p.DisplayName,
			kills:      1,
		}
	}
}

// getPool 取得當前大獎池金額（執行緒安全）
func (m *luckyProgressiveJackpotManager) getPool() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.pool
}

// shouldBroadcastPool 判斷是否需要廣播大獎池（每 10 秒一次）
func (m *luckyProgressiveJackpotManager) shouldBroadcastPool() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if time.Since(m.lastBroadcastAt) >= LuckyProgressiveJackpotBroadcastInterval {
		m.lastBroadcastAt = time.Now()
		return true
	}
	return false
}

// tryLuckyProgressiveJackpotFish 擊破 T220 後觸發大獎池爆發
func (g *Game) tryLuckyProgressiveJackpotFish(p *player.Player) {
	m := g.LuckyProgressiveJackpot

	m.mu.Lock()
	now := time.Now()

	// 全服冷卻檢查
	if now.Before(m.globalCooldownUntil) {
		m.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if cd, ok := m.personalCooldowns[p.ID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyProgressiveJackpotPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyProgressiveJackpotGlobalCD)

	// 取得大獎池金額和貢獻記錄
	pool := m.pool
	if pool < LuckyProgressiveJackpotMinPool {
		pool = LuckyProgressiveJackpotMinPool
	}

	// 計算總擊破次數
	totalKills := 0
	for _, c := range m.contributions {
		totalKills += c.kills
	}

	// 複製貢獻記錄（用於分配）
	contribs := make([]*jackpotContribution, 0, len(m.contributions))
	for _, c := range m.contributions {
		contribs = append(contribs, &jackpotContribution{
			playerID:   c.playerID,
			playerName: c.playerName,
			kills:      c.kills,
		})
	}

	// 重置大獎池和貢獻記錄
	m.pool = LuckyProgressiveJackpotMinPool
	m.contributions = make(map[string]*jackpotContribution)
	m.mu.Unlock()

	log.Printf("[ProgressiveJackpot] player=%s 觸發大獎池爆發！池=%d，參與玩家=%d，總擊破=%d",
		p.ID, pool, len(contribs), totalKills)

	// 按貢獻比例分配獎勵
	type playerReward struct {
		playerID   string
		playerName string
		kills      int
		reward     int
		pct        float64
	}
	rewards := make([]playerReward, 0, len(contribs))

	// 確保觸發玩家在列表中
	triggerInList := false
	for _, c := range contribs {
		if c.playerID == p.ID {
			triggerInList = true
			break
		}
	}
	if !triggerInList {
		contribs = append(contribs, &jackpotContribution{
			playerID:   p.ID,
			playerName: p.DisplayName,
			kills:      1,
		})
		totalKills++
	}

	// 計算每個玩家的獎勵
	distributed := 0
	for _, c := range contribs {
		var pct float64
		if totalKills > 0 {
			pct = float64(c.kills) / float64(totalKills)
		} else {
			pct = 1.0 / float64(len(contribs))
		}
		reward := int(float64(pool) * pct)
		if reward < 1 {
			reward = 1
		}
		distributed += reward
		rewards = append(rewards, playerReward{
			playerID:   c.playerID,
			playerName: c.playerName,
			kills:      c.kills,
			reward:     reward,
			pct:        pct,
		})
	}

	// 剩餘給觸發玩家
	remainder := pool - distributed
	for i := range rewards {
		if rewards[i].playerID == p.ID {
			rewards[i].reward += remainder
			break
		}
	}

	// 按獎勵排序（高到低）
	sort.Slice(rewards, func(i, j int) bool {
		return rewards[i].reward > rewards[j].reward
	})

	// 發放獎勵給所有玩家
	g.mu.RLock()
	players := make(map[string]*player.Player, len(g.Players))
	for id, pl := range g.Players {
		players[id] = pl
	}
	g.mu.RUnlock()

	for _, r := range rewards {
		if pl, ok := players[r.playerID]; ok {
			pl.AddReward(r.reward)
		}
		// 個人通知
		_ = g.Hub.Send(r.playerID, &ws.Message{
			Type: ws.MsgLuckyProgressiveJackpot,
			Payload: ws.LuckyProgressiveJackpotPayload{
				Event:       "jackpot_burst",
				PlayerID:    r.playerID,
				PlayerName:  r.playerName,
				TriggerName: p.DisplayName,
				Pool:        pool,
				Kills:       r.kills,
				TotalKills:  totalKills,
				Pct:         r.pct,
				Reward:      r.reward,
			},
		})
	}

	// 全服廣播大獎池爆發
	topName := ""
	topReward := 0
	if len(rewards) > 0 {
		topName = rewards[0].playerName
		topReward = rewards[0].reward
	}
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyProgressiveJackpot,
		Payload: ws.LuckyProgressiveJackpotPayload{
			Event:       "jackpot_burst_broadcast",
			TriggerName: p.DisplayName,
			Pool:        pool,
			TopName:     topName,
			TopReward:   topReward,
			PlayerCount: len(rewards),
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyProgressiveJackpot, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("💰 %s 觸發大獎池爆發！池=%d，%d 位玩家按貢獻分配！最高獎勵 %s +%d！",
			p.DisplayName, pool, len(rewards), topName, topReward),
		"color": "#FFD700",
	})
}

// broadcastJackpotPool 廣播當前大獎池金額（由 game loop 定期呼叫）
func (g *Game) broadcastJackpotPool() {
	m := g.LuckyProgressiveJackpot
	if !m.shouldBroadcastPool() {
		return
	}
	pool := m.getPool()
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyProgressiveJackpot,
		Payload: ws.LuckyProgressiveJackpotPayload{
			Event: "jackpot_update",
			Pool:  pool,
		},
	})
}
