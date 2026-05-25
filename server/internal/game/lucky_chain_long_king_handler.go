// lucky_chain_long_king_handler.go — T116 幸運千龍王輪盤魚
// server-event-agent 負責維護
// 業界依據：Royal Fishing Jili「ChainLong King dual-ring roulette up to 1000x Mega Win」
package game

import (
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyChainLongKingManager struct {
	mu             sync.Mutex
	personalCooldowns map[string]time.Time // 個人冷卻
	globalCooldown    time.Time            // 全服冷卻
}

func newLuckyChainLongKingManager() *luckyChainLongKingManager {
	return &luckyChainLongKingManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

func isLuckyChainLongKingFish(defID string) bool {
	return defID == "T116"
}

func (m *luckyChainLongKingManager) tryLuckyChainLongKing(g *Game, playerID, playerName string) {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCooldown) {
		m.mu.Unlock()
		return
	}
	if cd, ok := m.personalCooldowns[playerID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	// 設定冷卻：個人 28 秒，全服 45 秒
	m.personalCooldowns[playerID] = now.Add(28 * time.Second)
	m.globalCooldown = now.Add(45 * time.Second)
	m.mu.Unlock()

	// 廣播觸發
	g.hub.Broadcast(protocol.MsgLuckyChainLongKing, protocol.LuckyChainLongKingPayload{
		Event:       "trigger",
		TriggerID:   playerID,
		TriggerName: playerName,
	})

	// 抽取輪盤結果
	go func() {
		time.Sleep(800 * time.Millisecond)

		innerMult := rollChainLongKingInner()
		outerMult := rollChainLongKingOuter()
		finalMult := innerMult * outerMult
		isMegaWin := finalMult >= 500

		// 廣播輪盤結果
		g.hub.Broadcast(protocol.MsgLuckyChainLongKing, protocol.LuckyChainLongKingPayload{
			Event:       "spin",
			TriggerID:   playerID,
			TriggerName: playerName,
			InnerMult:   innerMult,
			OuterMult:   outerMult,
			FinalMult:   finalMult,
			IsMegaWin:   isMegaWin,
		})

		time.Sleep(1200 * time.Millisecond)

		// 計算獎勵
		g.mu.Lock()
		p, ok := g.players[playerID]
		var reward int
		if ok {
			bet := p.GetBetDef()
			reward = int(float64(bet.BetCost) * finalMult)
			p.AddCoins(reward)
		}
		g.mu.Unlock()

		// 廣播結果
		g.hub.Broadcast(protocol.MsgLuckyChainLongKing, protocol.LuckyChainLongKingPayload{
			Event:       "result",
			TriggerID:   playerID,
			TriggerName: playerName,
			FinalMult:   finalMult,
			IsMegaWin:   isMegaWin,
			Reward:      reward,
		})

		if isMegaWin {
			time.Sleep(500 * time.Millisecond)
			g.hub.Broadcast(protocol.MsgLuckyChainLongKing, protocol.LuckyChainLongKingPayload{
				Event:       "mega_win",
				TriggerID:   playerID,
				TriggerName: playerName,
				FinalMult:   finalMult,
				Reward:      reward,
			})
		}

		if ok {
			g.sendPlayerUpdate(playerID)
		}
	}()
}

func rollChainLongKingInner() float64 {
	weights := []struct {
		Mult   float64
		Weight int
	}{
		{2, 40}, {5, 25}, {10, 15}, {15, 10}, {20, 7}, {25, 3},
	}
	total := 0
	for _, w := range weights {
		total += w.Weight
	}
	r := rand.Intn(total)
	cum := 0
	for _, w := range weights {
		cum += w.Weight
		if r < cum {
			return w.Mult
		}
	}
	return 2
}

func rollChainLongKingOuter() float64 {
	weights := []struct {
		Mult   float64
		Weight int
	}{
		{5, 40}, {10, 25}, {20, 15}, {30, 10}, {40, 7}, {50, 3},
	}
	total := 0
	for _, w := range weights {
		total += w.Weight
	}
	r := rand.Intn(total)
	cum := 0
	for _, w := range weights {
		cum += w.Weight
		if r < cum {
			return w.Mult
		}
	}
	return 5
}
