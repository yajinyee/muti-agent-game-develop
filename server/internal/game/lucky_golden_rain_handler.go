// lucky_golden_rain_handler.go — T122 幸運黃金雨魚系統
// server-event-agent 負責維護
// 業界依據：Jackpot Fishing Jili「Golden Rain — coins shower from sky, each coin collected adds to jackpot」
// 設計：擊破 T122 後，觸發「黃金雨」：全場隨機生成 8-12 個「黃金幣」（虛擬目標）
// 玩家在 10 秒內點擊黃金幣可收集，每個黃金幣 = bet_cost × 3x 獎勵
// 收集 ≥ 8 個 → 「黃金豐收」：全服 ×2.0 加成 6 秒
// 個人冷卻 20 秒；全服冷卻 35 秒
package game

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyGoldenRainManager struct {
	mu              sync.Mutex
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
	activeSession   *goldenRainSession
	harvestBoost    *goldenHarvestBoost
}

type goldenRainSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	totalCoins        int // 生成的黃金幣總數
	collectedCoins    int // 已收集的黃金幣數
	totalReward       int
	expiresAt         time.Time
	settled           bool
}

type goldenHarvestBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyGoldenRainManager() *luckyGoldenRainManager {
	return &luckyGoldenRainManager{
		playerCooldowns: make(map[string]time.Time),
	}
}

func isLuckyGoldenRainFish(defID string) bool {
	return defID == "T122"
}

func (m *luckyGoldenRainManager) getGoldenHarvestMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.harvestBoost != nil && time.Now().Before(m.harvestBoost.expiresAt) {
		return m.harvestBoost.mult
	}
	return 1.0
}

func (g *Game) tryLuckyGoldenRain(playerID string, killerName string) {
	m := g.luckyGoldenRain
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

	// 生成 8-12 個黃金幣
	coinCount := 8 + rand.Intn(5)
	m.activeSession = &goldenRainSession{
		triggerPlayerID:   playerID,
		triggerPlayerName: killerName,
		totalCoins:        coinCount,
		expiresAt:         now.Add(10 * time.Second),
	}
	m.mu.Unlock()

	// 廣播觸發事件（含黃金幣位置）
	coinPositions := make([]protocol.GoldenCoinInfo, coinCount)
	for i := 0; i < coinCount; i++ {
		coinPositions[i] = protocol.GoldenCoinInfo{
			CoinID: i,
			X:      float64(100 + rand.Intn(1080)),
			Y:      float64(100 + rand.Intn(520)),
		}
	}

	g.hub.Broadcast(protocol.MsgLuckyGoldenRain, protocol.LuckyGoldenRainPayload{
		Event:         "trigger",
		TriggerID:     playerID,
		TriggerName:   killerName,
		TotalCoins:    coinCount,
		CoinPositions: coinPositions,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "🌧️💰 " + killerName + " 觸發黃金雨！" + string(rune('0'+coinCount)) + " 個黃金幣！快去收集！",
		Priority: "high",
		Color:    "#FFD700",
	})

	// 10 秒後結算
	go func() {
		time.Sleep(10 * time.Second)
		m.mu.Lock()
		s := m.activeSession
		if s == nil || s.settled || s.triggerPlayerID != playerID {
			m.mu.Unlock()
			return
		}
		s.settled = true
		collected := s.collectedCoins
		totalReward := s.totalReward
		isHarvest := collected >= 8
		m.mu.Unlock()

		if isHarvest {
			m.mu.Lock()
			m.harvestBoost = &goldenHarvestBoost{
				mult:      2.0,
				expiresAt: time.Now().Add(6 * time.Second),
			}
			m.mu.Unlock()

			g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
				Message:  "💰✨ 黃金豐收！" + killerName + " 收集 " + string(rune('0'+collected)) + " 個！全服 ×2.0 加成 6 秒！",
				Priority: "high",
				Color:    "#FFD700",
			})
			g.hub.Broadcast(protocol.MsgLuckyGoldenRain, protocol.LuckyGoldenRainPayload{
				Event:          "golden_harvest",
				TriggerID:      playerID,
				TriggerName:    killerName,
				CollectedCoins: collected,
				TotalReward:    totalReward,
			})

			go func() {
				time.Sleep(6 * time.Second)
				g.hub.Broadcast(protocol.MsgLuckyGoldenRain, protocol.LuckyGoldenRainPayload{
					Event:       "harvest_end",
					TriggerID:   playerID,
					TriggerName: killerName,
				})
			}()
		}

		g.hub.Broadcast(protocol.MsgLuckyGoldenRain, protocol.LuckyGoldenRainPayload{
			Event:          "settle",
			TriggerID:      playerID,
			TriggerName:    killerName,
			CollectedCoins: collected,
			TotalReward:    totalReward,
		})

		log.Printf("[GoldenRain] Player %s: collected=%d/%d, reward=%d, harvest=%v",
			playerID, collected, coinCount, totalReward, isHarvest)
	}()
}

// collectGoldenCoin 玩家點擊黃金幣時呼叫
func (g *Game) collectGoldenCoin(playerID string, coinID int) {
	m := g.luckyGoldenRain
	m.mu.Lock()
	s := m.activeSession
	if s == nil || s.settled || time.Now().After(s.expiresAt) {
		m.mu.Unlock()
		return
	}
	if s.collectedCoins >= s.totalCoins {
		m.mu.Unlock()
		return
	}

	g.mu.RLock()
	p, ok := g.players[playerID]
	var betCost int
	if ok {
		betCost = p.GetBetDef().BetCost
	}
	g.mu.RUnlock()

	reward := betCost * 3
	s.collectedCoins++
	s.totalReward += reward
	collected := s.collectedCoins
	totalReward := s.totalReward
	triggerID := s.triggerPlayerID
	triggerName := s.triggerPlayerName
	m.mu.Unlock()

	// 給收集者獎勵
	g.mu.Lock()
	if p2, ok2 := g.players[playerID]; ok2 {
		p2.AddCoins(reward)
		g.sendPlayerUpdate(playerID)
	}
	g.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyGoldenRain, protocol.LuckyGoldenRainPayload{
		Event:          "coin_collect",
		TriggerID:      triggerID,
		TriggerName:    triggerName,
		CollectorID:    playerID,
		CoinID:         coinID,
		CollectedCoins: collected,
		TotalReward:    totalReward,
	})
}
