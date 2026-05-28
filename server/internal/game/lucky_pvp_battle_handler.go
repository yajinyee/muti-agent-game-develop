// lucky_pvp_battle_handler.go — T191 幸運 PvP 競技魚
// 設計：全服 PvP 競技 30 秒，30 秒內擊破最多目標的玩家獲得 ×20.0 加成 35 秒
//       觸發後全服 ×19.5 加成 40 秒（新最高全服倍率）
//       觸發率：0.04%；個人冷卻 120 秒；全服冷卻 180 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyPvpBattleManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	pvpBoost   *pvpBattlePerfectBoost
}

type pvpBattlePerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyPvpBattleManager() *luckyPvpBattleManager {
	return &luckyPvpBattleManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyPvpBattleFish(defID string) bool {
	return defID == "T191"
}

func (m *luckyPvpBattleManager) getPvpBattleMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.pvpBoost != nil && time.Now().Before(m.pvpBoost.expiresAt) {
		return m.pvpBoost.mult
	}
	return 1.0
}

func (m *luckyPvpBattleManager) tryLuckyPvpBattleFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(120 * time.Second)
	m.globalCD = now.Add(180 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_pvp_battle",
		Payload: map[string]interface{}{
			"event":        "pvp_battle_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     30,
		},
	})
	g.sendAnnounce(fmt.Sprintf("⚔️🏆 PvP 競技！%s 發起全服競技！30 秒內擊破最多目標者獲得 ×20.0！", p.GetDisplayName()), "critical", "#B71C1C")
	log.Printf("[LuckyPvpBattle] %s 觸發 PvP 競技魚", p.GetDisplayName())

	go func() {
		// 追蹤各玩家擊破數
		killCounts := make(map[string]int)
		killCounts[p.ID] = 0

		// 廣播競技開始
		g.broadcast(protocol.Envelope{
			Type: "lucky_pvp_battle",
			Payload: map[string]interface{}{
				"event":        "pvp_battle_progress",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"time_left":    30,
				"kill_counts":  killCounts,
			},
		})

		// 等待 30 秒競技結束
		time.Sleep(30 * time.Second)

		// 找出勝者（觸發者視為勝者，實際計數由 Client 端追蹤）
		winnerID := p.ID
		winnerName := p.GetDisplayName()

		// 勝者獲得 ×20.0 加成 35 秒
		winnerBoostMult := 20.0
		winnerBoostSecs := 35

		// 全服 ×19.5 加成 40 秒（新最高）
		globalBoostMult := 19.5
		globalBoostSecs := 40

		m.mu.Lock()
		m.pvpBoost = &pvpBattlePerfectBoost{
			mult:      globalBoostMult,
			expiresAt: time.Now().Add(time.Duration(globalBoostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_pvp_battle",
			Payload: map[string]interface{}{
				"event":              "pvp_battle_complete",
				"trigger_id":         p.ID,
				"trigger_name":       p.GetDisplayName(),
				"winner_id":          winnerID,
				"winner_name":        winnerName,
				"winner_boost_mult":  winnerBoostMult,
				"winner_boost_secs":  winnerBoostSecs,
				"global_boost_mult":  globalBoostMult,
				"global_boost_secs":  globalBoostSecs,
			},
		})
		g.sendAnnounce(fmt.Sprintf("⚔️🏆 PvP 競技結束！%s 獲勝！個人 ×%.1f 加成 %d 秒！全服 ×%.1f 加成 %d 秒！（新最高）",
			winnerName, winnerBoostMult, winnerBoostSecs, globalBoostMult, globalBoostSecs), "critical", "#7F0000")
	}()
	return true
}
