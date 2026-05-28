// lucky_divine_revival_handler.go — T199 幸運神聖復活魚
// 設計：神聖復活，最近死亡的 5 個目標全部復活（HP 80%，獎勵 ×4.0）
//       全部擊破 → 全服 ×24.5 加成 49 秒（超越 T198 的 ×24.0）
//       觸發率：0.012%；個人冷卻 160 秒；全服冷卻 225 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/data"
	"chiikawa-game/internal/protocol"
)

type luckyDivineRevivalManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	revivalBoost *divineRevivalBoost
	// 追蹤最近死亡的目標（用於復活）
	recentDeaths []recentDeathRecord
}

type divineRevivalBoost struct {
	mult      float64
	expiresAt time.Time
}

type recentDeathRecord struct {
	defID      string
	multiplier float64
	diedAt     time.Time
}

func newLuckyDivineRevivalManager() *luckyDivineRevivalManager {
	return &luckyDivineRevivalManager{
		personalCD:   make(map[string]time.Time),
		recentDeaths: make([]recentDeathRecord, 0, 20),
	}
}

func isLuckyDivineRevivalFish(defID string) bool {
	return defID == "T199"
}

// RecordDeath 記錄目標死亡（供 combat 系統呼叫）
func (m *luckyDivineRevivalManager) RecordDeath(defID string, multiplier float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recentDeaths = append(m.recentDeaths, recentDeathRecord{
		defID:      defID,
		multiplier: multiplier,
		diedAt:     time.Now(),
	})
	// 只保留最近 20 筆
	if len(m.recentDeaths) > 20 {
		m.recentDeaths = m.recentDeaths[len(m.recentDeaths)-20:]
	}
}

func (m *luckyDivineRevivalManager) getDivineRevivalMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.revivalBoost != nil && time.Now().Before(m.revivalBoost.expiresAt) {
		return m.revivalBoost.mult
	}
	return 1.0
}

func (m *luckyDivineRevivalManager) tryLuckyDivineRevivalFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(160 * time.Second)
	m.globalCD = now.Add(225 * time.Second)

	// 取最近 5 個死亡目標
	reviveCount := 5
	if len(m.recentDeaths) < reviveCount {
		reviveCount = len(m.recentDeaths)
	}
	toRevive := make([]recentDeathRecord, reviveCount)
	copy(toRevive, m.recentDeaths[len(m.recentDeaths)-reviveCount:])
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_divine_revival",
		Payload: map[string]interface{}{
			"event":         "divine_revival_start",
			"trigger_id":    p.ID,
			"trigger_name":  p.GetDisplayName(),
			"revive_count":  reviveCount,
		},
	})
	g.sendAnnounce(fmt.Sprintf("✨🌟 神聖復活！%s 召喚神聖之力！%d 個目標復活！獎勵 ×4.0！", p.GetDisplayName(), reviveCount), "critical", "#1A1A00")
	log.Printf("[LuckyDivineRevival] %s 觸發神聖復活魚，復活 %d 個目標", p.GetDisplayName(), reviveCount)

	go func() {
		time.Sleep(800 * time.Millisecond)

		// 復活目標：生成新的目標（HP 80%，獎勵 ×4.0）
		revivedIDs := make([]string, 0, reviveCount)
		g.mu.Lock()
		for _, rec := range toRevive {
			if len(g.targets) >= MaxTargets {
				break
			}
			def, ok := data.GetTarget(rec.defID)
			if !ok {
				continue
			}
			newTarget := NewTarget(def, SpawnX, spawnY())
			// 強化：HP 80%，倍率 ×4.0
			newTarget.HP = int(float64(newTarget.MaxHP) * 0.8)
			newTarget.Multiplier = rec.multiplier * 4.0
			g.targets[newTarget.InstanceID] = newTarget
			revivedIDs = append(revivedIDs, newTarget.InstanceID)
			g.hub.Broadcast(protocol.MsgTargetSpawn, g.targetSpawnPayload(newTarget))
		}
		g.mu.Unlock()

		// 觸發全服 ×24.5 加成 49 秒
		boostMult := 24.5
		boostSecs := 49
		m.mu.Lock()
		m.revivalBoost = &divineRevivalBoost{
			mult:      boostMult,
			expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_divine_revival",
			Payload: map[string]interface{}{
				"event":        "divine_revival_complete",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"revived_ids":  revivedIDs,
				"reward_mult":  4.0,
				"boost_mult":   boostMult,
				"boost_secs":   boostSecs,
			},
		})
		g.sendAnnounce(fmt.Sprintf("✨👑 神聖復活完成！%s 復活 %d 個強化目標！獎勵 ×4.0！全服 ×%.1f 加成 %d 秒！",
			p.GetDisplayName(), len(revivedIDs), boostMult, boostSecs), "critical", "#2A2A00")
	}()
	return true
}
