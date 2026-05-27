// lucky_treasure_hunter_handler.go — T164 幸運寶藏獵人魚
// 業界依據：Treasure Hunt mechanic — 隨機挖掘寶藏，每個 ×10-×100 隨機倍率
// 設計：擊破後觸發寶藏獵人模式，隨機標記場上 5 個目標為寶藏
//       每個寶藏擊破獎勵 ×10-×100 隨機倍率；30 秒內全部擊破 → 完美寶藏全服 ×7.0 加成 15 秒
//       個人冷卻 55 秒；全服冷卻 80 秒
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyTreasureHunterManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	sessions     map[string]*treasureHunterSession
	perfectBoost *treasureHunterPerfectBoost
}

type treasureHunterSession struct {
	playerID    string
	treasures   map[string]float64 // instanceID -> multiplier
	found       int
	expiresAt   time.Time
}

type treasureHunterPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyTreasureHunterManager() *luckyTreasureHunterManager {
	return &luckyTreasureHunterManager{
		personalCD: make(map[string]time.Time),
		sessions:   make(map[string]*treasureHunterSession),
	}
}

func isLuckyTreasureHunterFish(defID string) bool {
	return defID == "T164"
}

func (m *luckyTreasureHunterManager) getTreasureHunterPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyTreasureHunterManager) getTreasureMult(playerID, instanceID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok || time.Now().After(sess.expiresAt) {
		return 1.0
	}
	mult, isTreasure := sess.treasures[instanceID]
	if !isTreasure {
		return 1.0
	}
	return mult
}

func (m *luckyTreasureHunterManager) onTreasureKilled(g *Game, p *Player, instanceID string) {
	m.mu.Lock()
	sess, ok := m.sessions[p.ID]
	if !ok || time.Now().After(sess.expiresAt) {
		m.mu.Unlock()
		return
	}
	if _, isTreasure := sess.treasures[instanceID]; !isTreasure {
		m.mu.Unlock()
		return
	}
	delete(sess.treasures, instanceID)
	sess.found++
	found := sess.found
	remaining := len(sess.treasures)
	allFound := remaining == 0
	if allFound {
		delete(m.sessions, p.ID)
	}
	m.mu.Unlock()

	if allFound {
		m.mu.Lock()
		m.perfectBoost = &treasureHunterPerfectBoost{
			mult:      7.0,
			expiresAt: time.Now().Add(15 * time.Second),
		}
		m.mu.Unlock()
		g.hub.Broadcast(protocol.MsgAnnounce, map[string]interface{}{
			"key":     "treasure_hunter_perfect",
			"message": fmt.Sprintf("💎 完美寶藏！%s 找到全部 5 個寶藏！全服 ×7.0 加成 15 秒！", p.ID),
			"mult":    7.0,
			"duration": 15,
		})
	}

	g.hub.Broadcast(protocol.MsgLuckyTreasureHunter, map[string]interface{}{
		"event":     "treasure_found",
		"player_id": p.ID,
		"found":     found,
		"remaining": remaining,
		"perfect":   allFound,
	})
}

func (m *luckyTreasureHunterManager) tryLuckyTreasureHunterFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(55 * time.Second)
	m.globalCD = now.Add(80 * time.Second)
	m.mu.Unlock()

	// 標記 5 個隨機目標為寶藏
	g.mu.Lock()
	var targetIDs []string
	for id, t := range g.targets {
		if t.HP > 0 {
			targetIDs = append(targetIDs, id)
		}
	}
	g.mu.Unlock()

	treasures := make(map[string]float64)
	count := 5
	if len(targetIDs) < count {
		count = len(targetIDs)
	}
	// 隨機選取
	rand.Shuffle(len(targetIDs), func(i, j int) { targetIDs[i], targetIDs[j] = targetIDs[j], targetIDs[i] })
	treasureMultWeights := []float64{10, 20, 30, 50, 100}
	for i := 0; i < count; i++ {
		mult := treasureMultWeights[rand.Intn(len(treasureMultWeights))]
		treasures[targetIDs[i]] = mult
	}

	m.mu.Lock()
	m.sessions[p.ID] = &treasureHunterSession{
		playerID:  p.ID,
		treasures: treasures,
		found:     0,
		expiresAt: now.Add(30 * time.Second),
	}
	m.mu.Unlock()

	log.Printf("[LuckyTreasureHunter] Player %s triggered treasure hunt (%d treasures)", p.ID, count)

	// 廣播寶藏標記
	var treasureList []map[string]interface{}
	for id, mult := range treasures {
		treasureList = append(treasureList, map[string]interface{}{
			"instance_id": id,
			"mult":        mult,
		})
	}
	g.hub.Broadcast(protocol.MsgLuckyTreasureHunter, map[string]interface{}{
		"event":     "start",
		"player_id": p.ID,
		"treasures": treasureList,
		"duration":  30,
	})

	return true
}
