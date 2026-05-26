// Package game — T149 幸運創世魚 handler
// server-event-agent 負責維護
// 業界依據：Ultimate boss mechanic — 召喚創世神，全場目標直接擊破
// 設計：擊破後召喚「創世神」；
//       全場所有目標 HP 歸零（直接擊破）；
//       每個被擊破的目標獎勵 ×5.0；
//       觸發全服 ×6.0 加成 15 秒；
//       個人冷卻 50 秒；全服冷卻 80 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyGenesisManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *genesisPerfectBoost
}

type genesisPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyGenesisManager() *luckyGenesisManager {
	return &luckyGenesisManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyGenesisFish(defID string) bool {
	return defID == "T149"
}

func (m *luckyGenesisManager) getGenesisPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyGenesisManager) tryLuckyGenesisFish(g *Game, playerID, playerName string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	if cd, ok := m.personalCD[playerID]; ok && now.Before(cd) {
		return false
	}
	if now.Before(m.globalCD) {
		return false
	}

	m.personalCD[playerID] = now.Add(50 * time.Second)
	m.globalCD = now.Add(80 * time.Second)

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyGenesis,
		Payload: protocol.LuckyGenesisPayload{
			Event:      "genesis_descend",
			PlayerID:   playerID,
			PlayerName: playerName,
		},
	})

	log.Printf("[LuckyGenesis] %s 觸發創世神降臨", playerName)

	go m.runGenesisJudgment(g, playerID, playerName)
	return true
}

func (m *luckyGenesisManager) runGenesisJudgment(g *Game, playerID, playerName string) {
	time.Sleep(1 * time.Second)

	// 創世審判：全場所有目標 HP 歸零（直接擊破），每個獎勵 ×5.0
	g.mu.Lock()
	killCount := 0
	totalReward := 0
	for _, t := range g.targets {
		if t.Def.ID == "B001" {
			continue
		}
		t.HP = 0
		killCount++
		// 計算獎勵（×5.0 加成）
		totalReward += int(float64(t.Def.Multiplier) * 5.0)
	}
	g.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyGenesis,
		Payload: protocol.LuckyGenesisPayload{
			Event:       "genesis_judgment",
			PlayerID:    playerID,
			PlayerName:  playerName,
			KillCount:   killCount,
			TotalReward: totalReward,
			MultBoost:   5.0,
		},
	})

	// 觸發全服 ×6.0 加成
	m.doGenesisPerfect(g, playerID, playerName, killCount)
}

func (m *luckyGenesisManager) doGenesisPerfect(g *Game, playerID, playerName string, killCount int) {
	m.mu.Lock()
	m.perfectBoost = &genesisPerfectBoost{
		mult:      6.0,
		expiresAt: time.Now().Add(15 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyGenesis,
		Payload: protocol.LuckyGenesisPayload{
			Event:      "genesis_blessing",
			PlayerID:   playerID,
			PlayerName: playerName,
			KillCount:  killCount,
			BoostMult:  6.0,
			BoostSec:   15,
		},
	})

	g.sendAnnounce("🌟 創世神降臨！"+playerName+" 審判 "+fmt.Sprintf("%d", killCount)+" 個目標！全服 ×6.0 加成 15 秒！", "critical", "#FFD700")

	time.Sleep(15 * time.Second)
	m.mu.Lock()
	m.perfectBoost = nil
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: protocol.MsgLuckyGenesis,
		Payload: protocol.LuckyGenesisPayload{
			Event:      "genesis_blessing_end",
			PlayerID:   playerID,
			PlayerName: playerName,
		},
	})
}
