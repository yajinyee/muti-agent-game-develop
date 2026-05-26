// lucky_jackpot_fish_handler.go — T126 幸運進階 Jackpot 魚系統
// server-event-agent 負責維護
// 業界依據：Jackpot Fishing Jili「Progressive Jackpot — Grand/Major/Minor/Mini 四層獎池，每次下注貢獻 1%」
// 設計：擊破 T126 後，觸發「進階 Jackpot 抽獎」：
//   - 四層獎池：Grand（≥50000）/ Major（≥20000）/ Minor（≥5000）/ Mini（≥1000）
//   - 每次任何玩家下注，自動貢獻 0.5% 到獎池（Grand 0.1% / Major 0.15% / Minor 0.15% / Mini 0.1%）
//   - 擊破 T126 後，依機率抽取獎池層級：Mini 60% / Minor 25% / Major 12% / Grand 3%
//   - 中獎後該層獎池重置為最低值
//   - Grand Jackpot 觸發全服廣播 + 全服 ×3.0 加成 10 秒
//   - 個人冷卻 20 秒；全服冷卻 35 秒
package game

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type jackpotTier struct {
	Name    string
	Pool    int
	MinPool int
	Weight  int // 抽中機率權重
}

type luckyJackpotFishManager struct {
	mu              sync.Mutex
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
	tiers           [4]*jackpotTier // 0=Mini, 1=Minor, 2=Major, 3=Grand
	grandBoost      *jackpotGrandBoost
}

type jackpotGrandBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyJackpotFishManager() *luckyJackpotFishManager {
	return &luckyJackpotFishManager{
		playerCooldowns: make(map[string]time.Time),
		tiers: [4]*jackpotTier{
			{Name: "Mini", Pool: 1000, MinPool: 1000, Weight: 60},
			{Name: "Minor", Pool: 5000, MinPool: 5000, Weight: 25},
			{Name: "Major", Pool: 20000, MinPool: 20000, Weight: 12},
			{Name: "Grand", Pool: 50000, MinPool: 50000, Weight: 3},
		},
	}
}

func isLuckyJackpotFish(defID string) bool {
	return defID == "T126"
}

// ContributeBet 每次下注貢獻到獎池（在 handleAttack 中呼叫）
func (m *luckyJackpotFishManager) ContributeBet(betCost int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// 貢獻比例：Mini 0.1% / Minor 0.15% / Major 0.15% / Grand 0.1%
	contributions := []int{
		betCost / 1000,       // Mini: 0.1%
		betCost * 15 / 10000, // Minor: 0.15%
		betCost * 15 / 10000, // Major: 0.15%
		betCost / 1000,       // Grand: 0.1%
	}
	for i, c := range contributions {
		if c < 1 {
			c = 1
		}
		m.tiers[i].Pool += c
	}
}

func (m *luckyJackpotFishManager) getGrandBoostMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.grandBoost != nil && time.Now().Before(m.grandBoost.expiresAt) {
		return m.grandBoost.mult
	}
	return 1.0
}

func (m *luckyJackpotFishManager) rollTier() int {
	totalWeight := 0
	for _, t := range m.tiers {
		totalWeight += t.Weight
	}
	r := rand.Intn(totalWeight)
	cumulative := 0
	for i, t := range m.tiers {
		cumulative += t.Weight
		if r < cumulative {
			return i
		}
	}
	return 0
}

func (g *Game) tryLuckyJackpotFish(playerID string, killerName string) {
	m := g.luckyJackpotFish
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCooldown) {
		m.mu.Unlock()
		return
	}
	if cd, ok := m.playerCooldowns[playerID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	m.playerCooldowns[playerID] = now.Add(20 * time.Second)
	m.globalCooldown = now.Add(35 * time.Second)

	// 取得各層獎池快照
	tierPools := [4]int{}
	for i, t := range m.tiers {
		tierPools[i] = t.Pool
	}
	m.mu.Unlock()

	// 廣播觸發事件
	g.hub.Broadcast(protocol.MsgLuckyJackpotFish, protocol.LuckyJackpotFishPayload{
		Event:       "trigger",
		TriggerID:   playerID,
		TriggerName: killerName,
		MiniPool:    tierPools[0],
		MinorPool:   tierPools[1],
		MajorPool:   tierPools[2],
		GrandPool:   tierPools[3],
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "🏆 " + killerName + " 觸發進階 Jackpot！",
		Priority: "high",
		Color:    "#FFD700",
	})

	// 抽獎延遲（2.5 秒動畫）
	go func() {
		time.Sleep(2500 * time.Millisecond)

		m.mu.Lock()
		tierIdx := m.rollTier()
		tier := m.tiers[tierIdx]
		reward := tier.Pool
		tier.Pool = tier.MinPool // 重置獎池
		isGrand := tierIdx == 3
		m.mu.Unlock()

		g.mu.Lock()
		p, ok := g.players[playerID]
		if !ok {
			g.mu.Unlock()
			return
		}
		p.AddCoins(reward)
		g.sendPlayerUpdate(playerID)

		// Grand Jackpot：全服 ×3.0 加成 10 秒
		if isGrand {
			m.mu.Lock()
			m.grandBoost = &jackpotGrandBoost{
				mult:      3.0,
				expiresAt: time.Now().Add(10 * time.Second),
			}
			m.mu.Unlock()
		}
		g.mu.Unlock()

		g.hub.Broadcast(protocol.MsgLuckyJackpotFish, protocol.LuckyJackpotFishPayload{
			Event:       "jackpot_result",
			TriggerID:   playerID,
			TriggerName: killerName,
			TierName:    tier.Name,
			TierIdx:     tierIdx,
			Reward:      reward,
			MiniPool:    m.tiers[0].Pool,
			MinorPool:   m.tiers[1].Pool,
			MajorPool:   m.tiers[2].Pool,
			GrandPool:   m.tiers[3].Pool,
		})

		if isGrand {
			g.hub.Broadcast(protocol.MsgLuckyJackpotFish, protocol.LuckyJackpotFishPayload{
				Event:       "grand_boost",
				TriggerID:   playerID,
				TriggerName: killerName,
				BoostMult:   3.0,
				BoostSecs:   10,
			})
			g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
				Message:  "🏆✨ GRAND JACKPOT！" + killerName + " 獲得 " + formatReward(reward) + "！全服 ×3.0 加成 10 秒！",
				Priority: "critical",
				Color:    "#FFD700",
			})
			// 10 秒後結束加成
			go func() {
				time.Sleep(10 * time.Second)
				g.hub.Broadcast(protocol.MsgLuckyJackpotFish, protocol.LuckyJackpotFishPayload{
					Event: "grand_boost_end",
				})
			}()
		} else {
			g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
				Message:  "🏆 " + killerName + " 中 " + tier.Name + " Jackpot！獲得 " + formatReward(reward) + "！",
				Priority: "high",
				Color:    "#FFD700",
			})
		}

		log.Printf("[JackpotFish] Player %s: tier=%s, reward=%d", playerID, tier.Name, reward)
	}()
}

func formatReward(r int) string {
	if r >= 1000 {
		return string(rune('0'+r/1000)) + "K"
	}
	return string(rune('0' + r))
}
