// lucky_domino_chain_handler.go — T246 幸運骨牌連鎖魚
// 設計：Domino Chain Reaction 機制（全新機制）
//       骨牌效應：第一個目標倒下觸發連鎖，每個目標倒下觸發下一個
//       連鎖長度越長倍率越高（最高 20 個骨牌，每個 ×50.0）
//       完美連鎖（≥15 個）→ 全服 ×55.0 加成 110 秒
//       業界依據：Domino Chain Reaction 機制（2026 新趨勢）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type dominoChainBoost struct {
	mult      float64
	expiresAt time.Time
}

type luckyDominoChainManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	boost      *dominoChainBoost
}

func newLuckyDominoChainManager() *luckyDominoChainManager {
	return &luckyDominoChainManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyDominoChainFish(defID string) bool {
	return defID == "T246"
}

func (m *luckyDominoChainManager) getDominoChainMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.boost != nil && time.Now().Before(m.boost.expiresAt) {
		return m.boost.mult
	}
	return 1.0
}

func (m *luckyDominoChainManager) tryLuckyDominoChainFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(600 * time.Second)
	m.personalCD[p.ID] = now.Add(540 * time.Second)
	m.mu.Unlock()

	// 計算場上目標數量（最多 20 個骨牌）
	g.mu.Lock()
	dominoCount := len(g.targets)
	if dominoCount > 20 {
		dominoCount = 20
	}
	if dominoCount < 5 {
		dominoCount = 5
	}
	g.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyDominoChain,
		Payload: map[string]interface{}{
			"event":        "domino_chain_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"domino_count": dominoCount,
			"per_domino":   50.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("DOMINO CHAIN! %s started a Domino Chain! %d dominoes falling! Each x50.0!", p.GetDisplayName(), dominoCount), "critical", "#FF8C00")
	log.Printf("[LuckyDominoChain] %s triggered Domino Chain fish - %d dominoes", p.GetDisplayName(), dominoCount)

	go func() {
		betCost := float64(p.GetBetDef().BetCost)
		perDominoMult := 50.0
		totalReward := 0
		totalMult := 0.0

		// 骨牌依序倒下，每個間隔 200ms
		for domino := 1; domino <= dominoCount; domino++ {
			time.Sleep(200 * time.Millisecond)

			// 倍率隨連鎖長度遞增
			chainMult := perDominoMult * (1.0 + float64(domino-1)*0.1)
			dominoReward := int(chainMult * betCost)
			totalReward += dominoReward
			totalMult += chainMult

			g.mu.Lock()
			p.Coins += dominoReward
			g.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyDominoChain,
				Payload: map[string]interface{}{
					"event":         "domino_fall",
					"domino_no":     domino,
					"chain_mult":    chainMult,
					"domino_reward": dominoReward,
					"total_so_far":  totalMult,
				},
			})
		}

		// 完美連鎖獎勵（≥15 個骨牌）
		isPerfect := dominoCount >= 15
		if isPerfect {
			perfectBonus := perDominoMult * float64(dominoCount) * 0.3
			perfectReward := int(perfectBonus * betCost)
			totalReward += perfectReward
			totalMult += perfectBonus

			g.mu.Lock()
			p.Coins += perfectReward
			g.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: protocol.MsgLuckyDominoChain,
				Payload: map[string]interface{}{
					"event":         "perfect_chain",
					"perfect_bonus": perfectBonus,
					"perfect_reward": perfectReward,
				},
			})
		}

		// 全服 ×55.0 加成 110 秒（里程碑）
		globalBonus := 55.0
		globalDuration := 110

		m.mu.Lock()
		m.boost = &dominoChainBoost{
			mult:      globalBonus,
			expiresAt: time.Now().Add(time.Duration(globalDuration) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: protocol.MsgLuckyDominoChain,
			Payload: map[string]interface{}{
				"event":          "domino_chain_complete",
				"domino_count":   dominoCount,
				"is_perfect":     isPerfect,
				"total_mult":     totalMult,
				"total_reward":   totalReward,
				"global_bonus":   globalBonus,
				"global_seconds": globalDuration,
				"milestone":      "GLOBAL_55X",
			},
		})
		g.sendAnnounce(fmt.Sprintf("DOMINO CHAIN COMPLETE! %s: %d dominoes, total x%.1f! MILESTONE: GLOBAL x%.1f for %ds!", p.GetDisplayName(), dominoCount, totalMult, globalBonus, globalDuration), "critical", "#FF8C00")
		log.Printf("[LuckyDominoChain] MILESTONE! %s: dominoes=%d, total_mult=%.1f, global=x%.1f (NEW RECORD x55.0)", p.GetDisplayName(), dominoCount, totalMult, globalBonus)
	}()

	return true
}
