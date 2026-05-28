// lucky_time_reversal_handler.go — T204 幸運時間逆流魚
// 設計：時間逆流，最近死亡的 10 個目標全部復活（HP 100%，獎勵 ×5.0）
//       全部擊破 → 完美逆流全服 ×29.0 加成 58 秒
//       觸發率：0.005%；個人冷卻 180 秒；全服冷卻 240 秒
//       業界依據：T199 神聖復活魚升級版 + 2026 時間逆流機制（最多 10 個目標）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/data"
	"chiikawa-game/internal/protocol"
)

type luckyTimeReversalManager struct {
	mu             sync.Mutex
	personalCD     map[string]time.Time
	globalCD       time.Time
	reversalBoost  *timeReversalBoost
	// 逆流目標：instanceID -> 倍率加成
	reversalTargets map[string]float64
}

type timeReversalBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyTimeReversalManager() *luckyTimeReversalManager {
	return &luckyTimeReversalManager{
		personalCD:      make(map[string]time.Time),
		reversalTargets: make(map[string]float64),
	}
}

func isLuckyTimeReversalFish(defID string) bool {
	return defID == "T204"
}

func (m *luckyTimeReversalManager) getTimeReversalMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.reversalBoost != nil && time.Now().Before(m.reversalBoost.expiresAt) {
		return m.reversalBoost.mult
	}
	return 1.0
}

func (m *luckyTimeReversalManager) getReversalTargetMult(instanceID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if mult, ok := m.reversalTargets[instanceID]; ok {
		return mult
	}
	return 1.0
}

func (m *luckyTimeReversalManager) onReversalTargetKilled(instanceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.reversalTargets, instanceID)
}

func (m *luckyTimeReversalManager) tryLuckyTimeReversalFish(g *Game, p *Player) bool {
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
		Type: "lucky_time_reversal",
		Payload: map[string]interface{}{
			"event":        "time_reversal_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("⏪🌀 時間逆流！%s 讓最近死亡的 10 個目標全部復活！獎勵 ×5.0！全服 ×29.0！", p.GetDisplayName()), "critical", "#8080FF")
	log.Printf("[LuckyTimeReversal] %s 觸發時間逆流魚（最多 10 個目標復活）", p.GetDisplayName())

	go func() {
		time.Sleep(800 * time.Millisecond)

		// 復活最多 10 個目標（使用基礎目標定義）
		reviveCount := 0
		reviveDefs := []string{"T001", "T002", "T003", "T004", "T005", "T006"}
		var revivedIDs []string

		g.mu.Lock()
		for reviveCount < 10 {
			defID := reviveDefs[reviveCount%len(reviveDefs)]
			var defPtr *data.TargetDef
			for i := range data.Targets {
				if data.Targets[i].ID == defID {
					defPtr = &data.Targets[i]
					break
				}
			}
			if defPtr == nil {
				reviveCount++
				continue
			}
			t := NewTarget(defPtr, SpawnX, spawnY())
			g.targets[t.InstanceID] = t
			revivedIDs = append(revivedIDs, t.InstanceID)
			g.hub.Broadcast(protocol.MsgTargetSpawn, g.targetSpawnPayload(t))
			reviveCount++
		}
		g.mu.Unlock()

		// 標記逆流目標（獎勵 ×5.0）
		m.mu.Lock()
		for _, id := range revivedIDs {
			m.reversalTargets[id] = 5.0
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_time_reversal",
			Payload: map[string]interface{}{
				"event":         "time_reversal_revived",
				"trigger_id":    p.ID,
				"trigger_name":  p.GetDisplayName(),
				"revive_count":  reviveCount,
				"revived_ids":   revivedIDs,
				"reward_mult":   5.0,
			},
		})

		// 45 秒後觸發全服加成
		time.Sleep(45 * time.Second)

		boostMult := 29.0
		boostSecs := 58
		m.mu.Lock()
		m.reversalTargets = make(map[string]float64) // 清除剩餘逆流目標
		m.reversalBoost = &timeReversalBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_time_reversal",
			Payload: map[string]interface{}{
				"event":        "time_reversal_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
			},
		})
		log.Printf("[LuckyTimeReversal] 時間逆流完成！復活 %d 個目標，全服 ×%.1f 加成 %d 秒", reviveCount, boostMult, boostSecs)
	}()
	return true
}
