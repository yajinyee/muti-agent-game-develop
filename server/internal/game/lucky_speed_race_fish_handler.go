// lucky_speed_race_fish_handler.go — 幸運競速賽魚系統（DAY-265）
// 業界原創「全服即時競速+排行榜爆發」機制
//
// 設計：擊破 T223 後，觸發「全服競速賽」（持續 30 秒）：
//   - 所有玩家競爭擊破數，每次擊破 +1 積分
//   - 每 5 秒廣播即時排行榜（前 3 名）
//   - 結算時：第 1 名 ×4.0、第 2 名 ×2.5、第 3 名 ×1.8、其他 ×1.2 安慰獎
//   - 個人冷卻 40 秒；全服冷卻 60 秒
//
// 設計差異：
//   - 與公會戰（T215，分隊競爭）不同，競速賽是「個人競速」，讓玩家有「我要打最多魚拿第一」的個人榮耀感
//   - 「即時排行榜每 5 秒廣播」讓玩家看到「我現在第幾名」，製造「要趕快多打幾條」的緊迫感
//   - 「第 1 名 ×4.0」是目前個人競爭類最高倍率，製造「拿第一超值」的動力
//   - 「所有人都有 ×1.2 安慰獎」確保即使沒進前三也有收益，降低挫敗感
//   - 「全服廣播最終排行榜」讓所有玩家看到「誰是冠軍」，製造社交話題感
//   - 業界依據：Fishing Frenzy Chapter 3（2026）Speed Race 機制，2026 年最熱門個人競速方向
package game

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckySpeedRaceFishPersonalCD = 40 * time.Second // 個人冷卻
	LuckySpeedRaceFishGlobalCD   = 60 * time.Second // 全服冷卻
	LuckySpeedRaceFishDuration   = 30 * time.Second // 競速賽時限
	LuckySpeedRaceFishScoreTick  = 5 * time.Second  // 排行榜廣播間隔
	LuckySpeedRaceFishRank1Mult  = 4.0              // 第 1 名倍率
	LuckySpeedRaceFishRank2Mult  = 2.5              // 第 2 名倍率
	LuckySpeedRaceFishRank3Mult  = 1.8              // 第 3 名倍率
	LuckySpeedRaceFishOtherMult  = 1.2              // 其他名次安慰獎倍率
	LuckySpeedRaceFishBoostSec   = 5                // 倍率加成持續秒數
)

// speedRaceEntry 競速賽玩家積分記錄
type speedRaceEntry struct {
	playerID   string
	playerName string
	score      int
}

// speedRaceSession 競速賽會話
type speedRaceSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	expiresAt         time.Time
	scores            map[string]*speedRaceEntry // playerID → entry
	mu                sync.Mutex
}

// speedRaceBoostEntry 競速賽倍率加成記錄
type speedRaceBoostEntry struct {
	mult      float64
	expiresAt time.Time
}

// luckySpeedRaceFishManager 幸運競速賽魚管理器
type luckySpeedRaceFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前活躍的競速賽會話（nil = 無）
	activeSession *speedRaceSession

	// 競速賽倍率加成（playerID → boostEntry）
	raceBoosts map[string]speedRaceBoostEntry
}

func newLuckySpeedRaceFishManager() *luckySpeedRaceFishManager {
	return &luckySpeedRaceFishManager{
		personalCooldowns: make(map[string]time.Time),
		raceBoosts:        make(map[string]speedRaceBoostEntry),
	}
}

// isLuckySpeedRaceFish 判斷是否為幸運競速賽魚
func isLuckySpeedRaceFish(defID string) bool {
	return defID == "T223"
}

// getLuckySpeedRaceBoostMult 取得競速賽倍率加成（供 handleKill 使用）
func (m *luckySpeedRaceFishManager) getLuckySpeedRaceBoostMult(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	entry, ok := m.raceBoosts[playerID]
	if !ok {
		return 1.0
	}
	if time.Now().After(entry.expiresAt) {
		delete(m.raceBoosts, playerID)
		return 1.0
	}
	return entry.mult
}

// isSpeedRaceActive 判斷競速賽是否進行中
func (m *luckySpeedRaceFishManager) isSpeedRaceActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil {
		return false
	}
	return time.Now().Before(m.activeSession.expiresAt)
}

// notifyLuckySpeedRaceKill 任何玩家擊破任何目標時，若競速賽進行中則累積積分
// 由 handleKill 呼叫（非 T223 目標）
func (g *Game) notifyLuckySpeedRaceKill(p *player.Player) {
	m := g.LuckySpeedRaceFish
	m.mu.Lock()
	sess := m.activeSession
	if sess == nil {
		m.mu.Unlock()
		return
	}
	now := time.Now()
	if now.After(sess.expiresAt) {
		m.activeSession = nil
		m.mu.Unlock()
		return
	}
	m.mu.Unlock()

	sess.mu.Lock()
	entry, ok := sess.scores[p.ID]
	if !ok {
		entry = &speedRaceEntry{
			playerID:   p.ID,
			playerName: p.DisplayName,
			score:      0,
		}
		sess.scores[p.ID] = entry
	}
	entry.score++
	score := entry.score
	sess.mu.Unlock()

	log.Printf("[SpeedRace] player=%s +1 積分（共 %d）", p.ID, score)
}

// tryLuckySpeedRaceFish 擊破 T223 後觸發競速賽
func (g *Game) tryLuckySpeedRaceFish(p *player.Player) {
	m := g.LuckySpeedRaceFish

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
	// 已有活躍競速賽
	if m.activeSession != nil && now.Before(m.activeSession.expiresAt) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckySpeedRaceFishPersonalCD)
	m.globalCooldownUntil = now.Add(LuckySpeedRaceFishGlobalCD)

	expiresAt := now.Add(LuckySpeedRaceFishDuration)
	sess := &speedRaceSession{
		triggerPlayerID:   p.ID,
		triggerPlayerName: p.DisplayName,
		expiresAt:         expiresAt,
		scores:            make(map[string]*speedRaceEntry),
	}
	// 觸發玩家預先加入積分表（確保在排行榜中）
	sess.scores[p.ID] = &speedRaceEntry{
		playerID:   p.ID,
		playerName: p.DisplayName,
		score:      0,
	}
	m.activeSession = sess
	m.mu.Unlock()

	log.Printf("[SpeedRace] player=%s 觸發競速賽！時限 %ds",
		p.ID, int(LuckySpeedRaceFishDuration.Seconds()))

	// 個人訊息：競速賽發起者
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckySpeedRaceFish,
		Payload: ws.LuckySpeedRaceFishPayload{
			Event:       "race_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckySpeedRaceFishDuration.Seconds()),
			Rank1Mult:   LuckySpeedRaceFishRank1Mult,
			Rank2Mult:   LuckySpeedRaceFishRank2Mult,
			Rank3Mult:   LuckySpeedRaceFishRank3Mult,
			OtherMult:   LuckySpeedRaceFishOtherMult,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckySpeedRaceFish,
		Payload: ws.LuckySpeedRaceFishPayload{
			Event:       "race_broadcast",
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckySpeedRaceFishDuration.Seconds()),
			Rank1Mult:   LuckySpeedRaceFishRank1Mult,
			Rank2Mult:   LuckySpeedRaceFishRank2Mult,
			Rank3Mult:   LuckySpeedRaceFishRank3Mult,
			OtherMult:   LuckySpeedRaceFishOtherMult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckySpeedRaceFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🏁 %s 發起競速賽！30秒內擊破最多目標！第1名 ×%.1f 倍率加成！",
			p.DisplayName, LuckySpeedRaceFishRank1Mult),
		"color": "#FF6B35",
	})

	// 啟動排行榜廣播 + 結算 goroutine
	go g.runSpeedRaceSession(sess)
}

// getSpeedRaceLeaderboard 取得競速賽排行榜（前 N 名）
func getSpeedRaceLeaderboard(sess *speedRaceSession, topN int) []speedRaceEntry {
	sess.mu.Lock()
	entries := make([]speedRaceEntry, 0, len(sess.scores))
	for _, e := range sess.scores {
		entries = append(entries, *e)
	}
	sess.mu.Unlock()

	// 依積分降序排列，積分相同依 playerID 字典序
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].score != entries[j].score {
			return entries[i].score > entries[j].score
		}
		return entries[i].playerID < entries[j].playerID
	})

	if topN > 0 && len(entries) > topN {
		return entries[:topN]
	}
	return entries
}

// runSpeedRaceSession 競速賽主循環（排行榜廣播 + 結算）
func (g *Game) runSpeedRaceSession(sess *speedRaceSession) {
	scoreTicker := time.NewTicker(LuckySpeedRaceFishScoreTick)
	endTimer := time.NewTimer(LuckySpeedRaceFishDuration)
	defer scoreTicker.Stop()
	defer endTimer.Stop()

	for {
		select {
		case <-scoreTicker.C:
			// 每 5 秒廣播即時排行榜（前 3 名）
			top3 := getSpeedRaceLeaderboard(sess, 3)
			leaderboard := make([]ws.SpeedRaceLeaderboardEntry, 0, len(top3))
			for i, e := range top3 {
				leaderboard = append(leaderboard, ws.SpeedRaceLeaderboardEntry{
					Rank:       i + 1,
					PlayerName: e.playerName,
					Score:      e.score,
				})
			}

			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgLuckySpeedRaceFish,
				Payload: ws.LuckySpeedRaceFishPayload{
					Event:       "race_leaderboard",
					Leaderboard: leaderboard,
				},
			})

		case <-endTimer.C:
			// 競速賽結束，結算
			m := g.LuckySpeedRaceFish
			m.mu.Lock()
			if m.activeSession != sess {
				m.mu.Unlock()
				return
			}
			m.activeSession = nil
			m.mu.Unlock()

			g.doSpeedRaceSettle(sess)
			return

		case <-g.stopCh:
			return
		}
	}
}

// doSpeedRaceSettle 競速賽結算
func (g *Game) doSpeedRaceSettle(sess *speedRaceSession) {
	m := g.LuckySpeedRaceFish

	// 取得完整排行榜
	allEntries := getSpeedRaceLeaderboard(sess, 0)

	if len(allEntries) == 0 {
		log.Printf("[SpeedRace] 競速賽結束，無玩家參與")
		return
	}

	// 套用倍率加成
	boostExpiry := time.Now().Add(time.Duration(LuckySpeedRaceFishBoostSec) * time.Second)
	m.mu.Lock()
	for i, e := range allEntries {
		var mult float64
		switch i {
		case 0:
			mult = LuckySpeedRaceFishRank1Mult
		case 1:
			mult = LuckySpeedRaceFishRank2Mult
		case 2:
			mult = LuckySpeedRaceFishRank3Mult
		default:
			mult = LuckySpeedRaceFishOtherMult
		}
		m.raceBoosts[e.playerID] = speedRaceBoostEntry{
			mult:      mult,
			expiresAt: boostExpiry,
		}
	}
	m.mu.Unlock()

	// 建立前 3 名排行榜（用於廣播）
	top3 := allEntries
	if len(top3) > 3 {
		top3 = top3[:3]
	}
	leaderboard := make([]ws.SpeedRaceLeaderboardEntry, 0, len(top3))
	for i, e := range top3 {
		var mult float64
		switch i {
		case 0:
			mult = LuckySpeedRaceFishRank1Mult
		case 1:
			mult = LuckySpeedRaceFishRank2Mult
		case 2:
			mult = LuckySpeedRaceFishRank3Mult
		}
		leaderboard = append(leaderboard, ws.SpeedRaceLeaderboardEntry{
			Rank:       i + 1,
			PlayerName: e.playerName,
			Score:      e.score,
			Mult:       mult,
		})
	}

	winner := allEntries[0]
	log.Printf("[SpeedRace] 結算！冠軍 %s（%d 積分），共 %d 名玩家參與",
		winner.playerName, winner.score, len(allEntries))

	// 全服廣播結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckySpeedRaceFish,
		Payload: ws.LuckySpeedRaceFishPayload{
			Event:        "race_result",
			PlayerName:   sess.triggerPlayerName,
			WinnerName:   winner.playerName,
			WinnerScore:  winner.score,
			Leaderboard:  leaderboard,
			TotalPlayers: len(allEntries),
			Rank1Mult:    LuckySpeedRaceFishRank1Mult,
			BoostSec:     LuckySpeedRaceFishBoostSec,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckySpeedRaceFish, sess.triggerPlayerName, 0, map[string]string{
		"message": fmt.Sprintf("🏆 競速賽結束！冠軍 %s（%d 擊破）！×%.1f 倍率加成 %ds！",
			winner.playerName, winner.score, LuckySpeedRaceFishRank1Mult, LuckySpeedRaceFishBoostSec),
		"color": "#FFD700",
	})
}
