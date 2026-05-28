// lucky_genesis_epoch_handler.go — T200 幸運創世紀元魚（里程碑第 200 個目標）
// 設計：創世紀元，全場 HP 歸零，每個獎勵 ×25.0
//       觸發後全服 ×25.0 加成 50 秒（史上最高，超越 T195 的 ×22.0）
//       觸發率：0.010%（最稀有）；個人冷卻 180 秒；全服冷卻 240 秒
//       里程碑意義：第 200 個 Lucky 目標物，代表遊戲發展的重要里程碑
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyGenesisEpochManager struct {
	mu          sync.Mutex
	personalCD  map[string]time.Time
	globalCD    time.Time
	epochBoost  *genesisEpochBoost
}

type genesisEpochBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyGenesisEpochManager() *luckyGenesisEpochManager {
	return &luckyGenesisEpochManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyGenesisEpochFish(defID string) bool {
	return defID == "T200"
}

func (m *luckyGenesisEpochManager) getGenesisEpochMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.epochBoost != nil && time.Now().Before(m.epochBoost.expiresAt) {
		return m.epochBoost.mult
	}
	return 1.0
}

func (m *luckyGenesisEpochManager) tryLuckyGenesisEpochFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(180 * time.Second)
	m.globalCD = now.Add(240 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_genesis_epoch",
		Payload: map[string]interface{}{
			"event":        "genesis_epoch_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"milestone":    200,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🌌🎊 創世紀元！%s 開啟第 200 個 Lucky 目標！全場 HP 歸零！每個獎勵 ×25.0！全服 ×25.0！", p.GetDisplayName()), "critical", "#000000")
	log.Printf("[LuckyGenesisEpoch] %s 觸發創世紀元魚（里程碑第 200 個目標，最稀有）", p.GetDisplayName())

	go func() {
		time.Sleep(1000 * time.Millisecond)

		// 創世紀元：全場 HP 歸零，每個獎勵 ×25.0（史上最高單次清場倍率）
		hitCount := g.applyUltimateJudgment(p, 25.0)

		// 觸發全服 ×25.0 加成 50 秒（史上最高全服倍率，超越 T195 的 ×22.0）
		boostMult := 25.0
		boostSecs := 50
		m.mu.Lock()
		m.epochBoost = &genesisEpochBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_genesis_epoch",
			Payload: map[string]interface{}{
				"event":        "genesis_epoch_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"hit_count":    hitCount,
				"reward_mult":  25.0,
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
				"milestone":    200,
			},
		})
		g.sendAnnounce(fmt.Sprintf("🌌👑 創世紀元完成！%s 清場 %d 個！每個 ×25.0！全服 ×%.1f 加成 %d 秒！（里程碑第 200 個 Lucky 目標，史上最高）",
			p.GetDisplayName(), hitCount, boostMult, boostSecs), "critical", "#000000")
	}()
	return true
}
