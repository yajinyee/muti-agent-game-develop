// Package game — T130 幸運崩潰魚 handler
// server-event-agent 負責維護
// 業界依據：Lucky Fish by AbraCadabra「crash mechanic — multiplier rises until crash」
// 設計：擊破後觸發崩潰倍率，玩家可隨時收割，崩潰前收割 ≥5.0x 觸發完美收割全服加成
package game

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyCrashFishManager struct {
	mu            sync.Mutex
	personalCD    map[string]time.Time // 個人冷卻
	globalCD      time.Time            // 全服冷卻
	activeSession *crashFishSession
	perfectBoost  *crashFishPerfectBoost
}

type crashFishPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type crashFishSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	currentMult       float64
	crashAt           time.Time
	settled           bool
	harvested         bool
	harvestMult       float64
}

func newLuckyCrashFishManager() *luckyCrashFishManager {
	return &luckyCrashFishManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyCrashFish(defID string) bool {
	return defID == "T130"
}

func (m *luckyCrashFishManager) getCrashPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (g *Game) tryLuckyCrashFish(playerID, playerName string) {
	m := g.luckyCrashFish
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) {
		m.mu.Unlock()
		return
	}
	if cd, ok := m.personalCD[playerID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	if m.activeSession != nil && !m.activeSession.settled {
		m.mu.Unlock()
		return
	}
	m.personalCD[playerID] = now.Add(20 * time.Second)
	m.globalCD = now.Add(35 * time.Second)

	// 隨機崩潰時間：5-12 秒
	crashDelay := time.Duration(5000+rand.Intn(7000)) * time.Millisecond
	session := &crashFishSession{
		triggerPlayerID:   playerID,
		triggerPlayerName: playerName,
		currentMult:       1.0,
		crashAt:           now.Add(crashDelay),
		settled:           false,
		harvested:         false,
	}
	m.activeSession = session
	m.mu.Unlock()

	log.Printf("[LuckyCrashFish] Triggered by %s, crash in %.1fs", playerName, crashDelay.Seconds())

	g.hub.Broadcast(protocol.MsgLuckyCrashFish, protocol.LuckyCrashFishPayload{
		Event:       "crash_start",
		PlayerID:    playerID,
		PlayerName:  playerName,
		CurrentMult: 1.0,
		CrashIn:     crashDelay.Seconds(),
	})

	go g.runCrashFishTick(playerID, playerName, session)
}

func (g *Game) runCrashFishTick(playerID, playerName string, session *crashFishSession) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C

		m := g.luckyCrashFish
		m.mu.Lock()

		if session.settled {
			m.mu.Unlock()
			return
		}

		now := time.Now()

		// 檢查是否已被收割
		if session.harvested {
			m.mu.Unlock()
			return
		}

		// 檢查是否崩潰
		if now.After(session.crashAt) {
			session.settled = true
			m.mu.Unlock()

			log.Printf("[LuckyCrashFish] CRASH! Player %s lost mult %.1fx", playerName, session.currentMult)
			g.hub.Broadcast(protocol.MsgLuckyCrashFish, protocol.LuckyCrashFishPayload{
				Event:       "crash",
				PlayerID:    playerID,
				PlayerName:  playerName,
				CurrentMult: session.currentMult,
			})
			return
		}

		// 倍率上升：每 0.5 秒 +0.3x，最高 10.0x
		session.currentMult += 0.3
		if session.currentMult > 10.0 {
			session.currentMult = 10.0
		}
		currentMult := session.currentMult
		timeLeft := session.crashAt.Sub(now).Seconds()
		m.mu.Unlock()

		g.hub.Broadcast(protocol.MsgLuckyCrashFish, protocol.LuckyCrashFishPayload{
			Event:       "mult_rise",
			PlayerID:    playerID,
			PlayerName:  playerName,
			CurrentMult: currentMult,
			TimeLeft:    timeLeft,
		})
	}
}

// handleCrashHarvest 玩家點擊收割
func (g *Game) handleCrashHarvest(playerID string) {
	m := g.luckyCrashFish
	m.mu.Lock()

	session := m.activeSession
	if session == nil || session.settled || session.harvested {
		m.mu.Unlock()
		return
	}
	if session.triggerPlayerID != playerID {
		m.mu.Unlock()
		return
	}

	session.harvested = true
	session.settled = true
	harvestMult := session.currentMult
	playerName := session.triggerPlayerName
	m.mu.Unlock()

	log.Printf("[LuckyCrashFish] Harvested by %s at %.1fx", playerName, harvestMult)

	// 計算獎勵
	g.mu.RLock()
	p, ok := g.players[playerID]
	g.mu.RUnlock()
	if !ok {
		return
	}
	betCost := p.GetBetDef().BetCost
	reward := int(float64(betCost) * harvestMult * 5.0)

	g.mu.Lock()
	p.AddCoins(reward)
	g.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyCrashFish, protocol.LuckyCrashFishPayload{
		Event:       "harvest",
		PlayerID:    playerID,
		PlayerName:  playerName,
		CurrentMult: harvestMult,
		Reward:      reward,
	})

	// 完美收割：≥ 5.0x
	if harvestMult >= 5.0 {
		g.doCrashPerfect(playerID, playerName, harvestMult)
	}
}

func (g *Game) doCrashPerfect(playerID, playerName string, harvestMult float64) {
	m := g.luckyCrashFish
	m.mu.Lock()
	m.perfectBoost = &crashFishPerfectBoost{
		mult:      2.0,
		expiresAt: time.Now().Add(5 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyCrashFish] Perfect harvest! %s ×%.1f → global ×2.0 for 5s", playerName, harvestMult)

	g.hub.Broadcast(protocol.MsgLuckyCrashFish, protocol.LuckyCrashFishPayload{
		Event:       "perfect_harvest",
		PlayerID:    playerID,
		PlayerName:  playerName,
		CurrentMult: harvestMult,
		BoostMult:   2.0,
		BoostSecs:   5,
	})

	time.AfterFunc(5*time.Second, func() {
		m.mu.Lock()
		m.perfectBoost = nil
		m.mu.Unlock()
		g.hub.Broadcast(protocol.MsgLuckyCrashFish, protocol.LuckyCrashFishPayload{
			Event:      "perfect_end",
			PlayerID:   playerID,
			PlayerName: playerName,
		})
	})
}
