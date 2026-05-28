// lucky_guild_battle_handler.go — T207 幸運公會戰魚
// 設計：Fishing Frenzy Chapter 3 Guild Wars 機制（2026-05-27）
//       全服公會戰 45 秒，擊破最多目標的玩家獲得 ×35.0 加成
//       觸發後全服 ×32.0 加成 64 秒（超越 T206 的 ×31.0）
//       觸發率：0.002%（極稀有）；個人冷卻 260 秒；全服冷卻 320 秒
//       業界依據：Fishing Frenzy Chapter 3「Guild Wars + Boss Fish」（2026-05-27）
package game

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyGuildBattleManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	battleBoost *guildBattleBoost
	killCounts  map[string]int // playerID -> kill count during battle
	isActive   bool
}

type guildBattleBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyGuildBattleManager() *luckyGuildBattleManager {
	return &luckyGuildBattleManager{
		personalCD: make(map[string]time.Time),
		killCounts: make(map[string]int),
	}
}

func isLuckyGuildBattleFish(defID string) bool {
	return defID == "T207"
}

func (m *luckyGuildBattleManager) getGuildBattleMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.battleBoost != nil && time.Now().Before(m.battleBoost.expiresAt) {
		return m.battleBoost.mult
	}
	return 1.0
}

func (m *luckyGuildBattleManager) onKillDuringBattle(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isActive {
		m.killCounts[playerID]++
	}
}

func (m *luckyGuildBattleManager) tryLuckyGuildBattleFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(260 * time.Second)
	m.globalCD = now.Add(320 * time.Second)
	m.killCounts = make(map[string]int)
	m.isActive = true
	m.mu.Unlock()

	battleSecs := 45
	g.broadcast(protocol.Envelope{
		Type: "lucky_guild_battle",
		Payload: map[string]interface{}{
			"event":        "battle_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"battle_secs":  battleSecs,
		},
	})
	g.sendAnnounce(fmt.Sprintf("⚔️🏆 公會戰！%s 發動公會戰！45 秒內擊破最多目標者獲得 ×35.0！", p.GetDisplayName()), "critical", "#FFD700")
	log.Printf("[LuckyGuildBattle] %s 觸發公會戰魚（45 秒公會戰）", p.GetDisplayName())

	go func() {
		time.Sleep(time.Duration(battleSecs) * time.Second)

		m.mu.Lock()
		m.isActive = false
		// 找出擊破最多的玩家
		type playerKill struct {
			id    string
			kills int
		}
		var rankings []playerKill
		for pid, kills := range m.killCounts {
			rankings = append(rankings, playerKill{pid, kills})
		}
		sort.Slice(rankings, func(i, j int) bool {
			return rankings[i].kills > rankings[j].kills
		})
		m.mu.Unlock()

		winnerID := ""
		winnerKills := 0
		if len(rankings) > 0 {
			winnerID = rankings[0].id
			winnerKills = rankings[0].kills
		}

		// 勝利者獲得 ×35.0 個人加成
		winnerMult := 35.0
		// 觸發全服 ×32.0 加成 64 秒
		globalBoostMult := 32.0
		globalBoostSecs := 64
		m.mu.Lock()
		m.battleBoost = &guildBattleBoost{
			mult:      globalBoostMult,
			expiresAt: time.Now().Add(time.Duration(globalBoostSecs) * time.Second),
		}
		m.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_guild_battle",
			Payload: map[string]interface{}{
				"event":        "battle_complete",
				"winner_id":    winnerID,
				"winner_kills": winnerKills,
				"winner_mult":  winnerMult,
				"global_mult":  globalBoostMult,
				"global_secs":  globalBoostSecs,
				"rankings":     rankings,
			},
		})
		g.sendAnnounce(fmt.Sprintf("⚔️🏆 公會戰結束！勝者 %s（%d 擊破）獲得 ×%.1f！全服 ×%.1f 加成 %d 秒！", winnerID, winnerKills, winnerMult, globalBoostMult, globalBoostSecs), "critical", "#FFD700")
		log.Printf("[LuckyGuildBattle] 公會戰完成！勝者 %s（%d 擊破），全服 ×%.1f 加成 %d 秒", winnerID, winnerKills, globalBoostMult, globalBoostSecs)
	}()
	return true
}
