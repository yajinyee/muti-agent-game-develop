// lucky_multiplier_ladder_handler.go — T213 幸運倍率梯魚
// 設計：Multiplier Ladder 機制（Relax Gaming「Cod of Thunder」2026）
//       觸發後 30 秒內每次擊破提升梯度（Lv.1→Lv.10），每級 +4.5 倍率
//       Lv.10（擊破 10 個）→ 全服 ×37.0 加成 74 秒（超越 T211 的 ×36.0）
//       觸發率：0.003%；個人冷卻 270 秒；全服冷卻 330 秒
//       業界依據：Relax Gaming「Cod of Thunder Dream Drop」Multiplier Ladder（2026）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyMultiplierLadderManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	active     *multiplierLadderState
}

type multiplierLadderState struct {
	playerID  string
	level     int
	kills     int
	expiresAt time.Time
}

func newLuckyMultiplierLadderManager() *luckyMultiplierLadderManager {
	return &luckyMultiplierLadderManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyMultiplierLadderFish(defID string) bool {
	return defID == "T213"
}

func (m *luckyMultiplierLadderManager) onKill(g *Game, playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.active == nil || time.Now().After(m.active.expiresAt) {
		m.active = nil
		return
	}
	if m.active.playerID != playerID {
		return
	}

	m.active.kills++
	if m.active.kills > 10 {
		m.active.kills = 10
	}
	newLevel := m.active.kills
	levelMult := float64(newLevel) * 4.5

	g.broadcast(protocol.Envelope{
		Type: "lucky_multiplier_ladder",
		Payload: map[string]interface{}{
			"event":      "level_up",
			"player_id":  playerID,
			"level":      newLevel,
			"level_mult": levelMult,
			"kills":      m.active.kills,
		},
	})

	if newLevel >= 10 {
		globalBoostMult := 37.0
		globalBoostSecs := 74
		m.active = nil

		g.broadcast(protocol.Envelope{
			Type: "lucky_multiplier_ladder",
			Payload: map[string]interface{}{
				"event":        "ladder_max",
				"player_id":    playerID,
				"global_mult":  globalBoostMult,
				"global_secs":  globalBoostSecs,
			},
		})
		g.sendAnnounce(fmt.Sprintf("🪜🌟 倍率梯頂端！全服 ×%.1f 加成 %d 秒！", globalBoostMult, globalBoostSecs), "critical", "#FFD700")
		log.Printf("[LuckyMultiplierLadder] 倍率梯頂端！全服 ×%.1f 加成 %d 秒（超越 T211 的 ×36.0）", globalBoostMult, globalBoostSecs)
	}
}

func (m *luckyMultiplierLadderManager) tryLuckyMultiplierLadderFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(270 * time.Second)
	m.globalCD = now.Add(330 * time.Second)

	m.active = &multiplierLadderState{
		playerID:  p.ID,
		level:     0,
		kills:     0,
		expiresAt: now.Add(30 * time.Second),
	}
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_multiplier_ladder",
		Payload: map[string]interface{}{
			"event":        "ladder_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     30,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🪜📈 倍率梯！%s 觸發倍率梯魚！30 秒內每次擊破提升等級，Lv.10 觸發全服 ×37.0！", p.GetDisplayName()), "critical", "#FFD700")
	log.Printf("[LuckyMultiplierLadder] %s 觸發倍率梯魚（30 秒，Lv.10 目標）", p.GetDisplayName())

	go func() {
		time.Sleep(30 * time.Second)
		m.mu.Lock()
		state := m.active
		m.active = nil
		m.mu.Unlock()

		if state != nil && state.kills < 10 {
			g.broadcast(protocol.Envelope{
				Type: "lucky_multiplier_ladder",
				Payload: map[string]interface{}{
					"event":  "ladder_timeout",
					"kills":  state.kills,
					"level":  state.kills,
				},
			})
			log.Printf("[LuckyMultiplierLadder] 倍率梯超時，最終等級 %d", state.kills)
		}
	}()
	return true
}
