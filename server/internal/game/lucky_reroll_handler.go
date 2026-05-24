// lucky_reroll_handler.go — 幸運倍率重擲魚系統（DAY-271）
// 業界依據：GONE Fishing 2026-05-03 patch「4.25x multiplier reroll → 5x-10x max win formula」
//           2026 年最新 RNG 設計趨勢：高倍率時有機會再抽一次更高倍率
//
// 設計：擊破 T229 後，觸發「倍率重擲」：
//   - Server 為觸發玩家的下一次擊破「重擲倍率」（最多 3 次，取最高值）
//   - 每次重擲有 40% 機率提升倍率（×1.5 到 ×4.0 隨機）
//   - 最終用最高倍率計算獎勵（個人）
//   - 個人冷卻 20 秒；全服冷卻 35 秒
//
// 設計差異：
//   - 與輪盤（T227，固定扇區）不同，重擲是「動態累積最高值」，讓玩家有「再擲一次，說不定更高！」的期待感
//   - 「最多 3 次重擲取最高值」讓玩家有「每次重擲都可能更好」的動力
//   - 「40% 機率提升」讓重擲有不確定性，不是每次都能提升
//   - 「×1.5 到 ×4.0 隨機提升」讓每次重擲的提升幅度也有驚喜感
//   - 「全服廣播最終倍率」讓所有玩家看到「有人重擲到 ×4.0！」，製造羨慕感
//   - 「即時顯示每次重擲結果」讓玩家看到「第 1 擲 ×1.8，第 2 擲 ×3.2，第 3 擲 ×2.1 → 最終 ×3.2」的過程
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyRerollPersonalCD  = 20 * time.Second // 個人冷卻
	LuckyRerollGlobalCD    = 35 * time.Second // 全服冷卻
	LuckyRerollMaxRolls    = 3                // 最多重擲次數
	LuckyRerollUpChance    = 0.40             // 40% 機率提升
	LuckyRerollMinMult     = 1.5              // 最低提升倍率
	LuckyRerollMaxMult     = 4.0              // 最高提升倍率
	LuckyRerollSessionTTL  = 30 * time.Second // session 存活時間（等待下一次擊破）
)

// rerollSession 重擲 session（等待下一次擊破）
type rerollSession struct {
	playerID    string
	playerName  string
	rolls       []float64  // 每次重擲的倍率結果
	bestMult    float64    // 目前最高倍率
	expiresAt   time.Time
	used        bool       // 是否已用於一次擊破
}

// luckyRerollManager 幸運倍率重擲魚管理器
type luckyRerollManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 活躍 session（playerID → session）
	activeSessions map[string]*rerollSession
}

func newLuckyRerollManager() *luckyRerollManager {
	return &luckyRerollManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*rerollSession),
	}
}

// isLuckyRerollFish 判斷是否為幸運倍率重擲魚
func isLuckyRerollFish(defID string) bool {
	return defID == "T229"
}

// getRerollMult 取得重擲倍率（供 handleKill 使用）
// 回傳 (倍率, 是否有重擲加成)
func (m *luckyRerollManager) getRerollMult(playerID string) (float64, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.activeSessions[playerID]
	if !ok || session.used || time.Now().After(session.expiresAt) {
		if ok {
			delete(m.activeSessions, playerID)
		}
		return 1.0, false
	}

	// 標記為已使用
	session.used = true
	return session.bestMult, true
}

// consumeRerollSession 消耗 session（擊破後呼叫，取得最終倍率和 rolls 資訊）
func (m *luckyRerollManager) consumeRerollSession(playerID string) (*rerollSession, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.activeSessions[playerID]
	if !ok {
		return nil, false
	}
	delete(m.activeSessions, playerID)
	return session, true
}

// tryLuckyRerollFish 擊破 T229 後觸發重擲
func (g *Game) tryLuckyRerollFish(p *player.Player) {
	m := g.LuckyReroll

	m.mu.Lock()
	now := time.Now()

	// 全服冷卻檢查
	if now.Before(m.globalCooldownUntil) {
		m.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if cd, ok := m.personalCooldowns[p.ID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyRerollPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyRerollGlobalCD)
	m.mu.Unlock()

	// 執行 3 次重擲
	rng := rand.New(rand.NewSource(now.UnixNano()))
	rolls := make([]float64, LuckyRerollMaxRolls)
	bestMult := 1.0

	for i := 0; i < LuckyRerollMaxRolls; i++ {
		if rng.Float64() < LuckyRerollUpChance {
			// 提升：隨機 ×1.5 到 ×4.0
			mult := LuckyRerollMinMult + rng.Float64()*(LuckyRerollMaxMult-LuckyRerollMinMult)
			// 四捨五入到 0.1
			mult = float64(int(mult*10+0.5)) / 10.0
			rolls[i] = mult
			if mult > bestMult {
				bestMult = mult
			}
		} else {
			// 未提升：保持 ×1.0
			rolls[i] = 1.0
		}
	}

	// 確保至少有一次提升（讓玩家不會完全失望）
	if bestMult == 1.0 {
		rolls[0] = LuckyRerollMinMult
		bestMult = LuckyRerollMinMult
	}

	log.Printf("[Reroll] player=%s 觸發重擲！rolls=%v bestMult=×%.1f",
		p.ID, rolls, bestMult)

	// 建立 session
	session := &rerollSession{
		playerID:   p.ID,
		playerName: p.DisplayName,
		rolls:      rolls,
		bestMult:   bestMult,
		expiresAt:  now.Add(LuckyRerollSessionTTL),
	}

	m.mu.Lock()
	m.activeSessions[p.ID] = session
	m.mu.Unlock()

	// 個人通知：重擲開始
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyReroll,
		Payload: ws.LuckyRerollPayload{
			Event:      "reroll_start",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Rolls:      rolls,
			BestMult:   bestMult,
			MaxRolls:   LuckyRerollMaxRolls,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyReroll,
		Payload: ws.LuckyRerollPayload{
			Event:      "reroll_broadcast",
			PlayerName: p.DisplayName,
			BestMult:   bestMult,
		},
	})

	// 全服公告（×3.0 以上才公告）
	if bestMult >= 3.0 {
		g.Announce.Create(announce.EventLuckyReroll, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🎲 %s 倍率重擲！最高 ×%.1f！下一擊必中！",
				p.DisplayName, bestMult),
			"color": "#FFD700",
		})
	} else {
		g.Announce.Create(announce.EventLuckyReroll, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🎲 %s 觸發倍率重擲！×%.1f 等待下一擊！",
				p.DisplayName, bestMult),
			"color": "#FF8C00",
		})
	}

	// session 超時清理
	go func() {
		time.Sleep(LuckyRerollSessionTTL)
		m.mu.Lock()
		if s, ok := m.activeSessions[p.ID]; ok && s == session {
			delete(m.activeSessions, p.ID)
			log.Printf("[Reroll] player=%s session 超時，未使用", p.ID)
		}
		m.mu.Unlock()

		// 通知超時
		_ = g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyReroll,
			Payload: ws.LuckyRerollPayload{
				Event:    "reroll_expire",
				PlayerID: p.ID,
			},
		})
	}()
}

// notifyRerollKill 重擲加成被使用時呼叫（由 handleKill 在套用倍率後呼叫）
func (g *Game) notifyRerollKill(p *player.Player, targetName string, reward int, bestMult float64, rolls []float64) {
	log.Printf("[Reroll] player=%s 使用重擲！×%.1f，目標=%s，獎勵=%d",
		p.ID, bestMult, targetName, reward)

	// 個人通知：重擲結算
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyReroll,
		Payload: ws.LuckyRerollPayload{
			Event:      "reroll_used",
			PlayerID:   p.ID,
			TargetName: targetName,
			Reward:     reward,
			BestMult:   bestMult,
			Rolls:      rolls,
		},
	})

	// 全服廣播（×3.0 以上）
	if bestMult >= 3.0 {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyReroll,
			Payload: ws.LuckyRerollPayload{
				Event:      "reroll_result_broadcast",
				PlayerName: p.DisplayName,
				BestMult:   bestMult,
				Reward:     reward,
			},
		})

		g.Announce.Create(announce.EventLuckyReroll, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🎲 %s 重擲命中！×%.1f 大獎 %d！",
				p.DisplayName, bestMult, reward),
			"color": "#FFD700",
		})
	}
}
