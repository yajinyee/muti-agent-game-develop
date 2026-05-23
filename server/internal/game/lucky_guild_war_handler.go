// lucky_guild_war_handler.go — 幸運公會戰魚系統（DAY-257）
// 業界原創「全服分隊競爭→勝隊爆發」機制
//
// 設計：擊破 T215 後，全服玩家自動分成兩隊（紅隊/藍隊，依玩家 ID 奇偶分配）：
//   - 30 秒內競爭擊破數，每次擊破為己隊累積積分
//   - 每 5 秒廣播即時比分（讓全服看到「哪隊領先」）
//   - 30 秒後結算：勝隊全員 ×2.5 倍率加成（5 秒）；敗隊 ×1.2 安慰獎（全服共享）
//   - 平局：雙隊各獲得 ×1.8 倍率加成（5 秒）
//   - 個人冷卻 35 秒；全服冷卻 55 秒
//
// 設計差異：
//   - 與全服充能（T214，所有玩家合作）不同，公會戰是「玩家分隊競爭」，
//     讓玩家有「我要幫紅隊贏」的歸屬感和競爭動力
//   - 「依玩家 ID 奇偶分隊」讓分隊自動且公平，不需要玩家手動選擇
//   - 「每 5 秒比分廣播」讓玩家即時看到「哪隊領先」，製造「要趕快多打幾條」的緊迫感
//   - 「敗隊也有 ×1.2 安慰獎」確保即使輸了也有收益，降低挫敗感
//   - 「平局 ×1.8」鼓勵雙隊勢均力敵，製造「最後幾秒決勝負」的高潮感
//   - 觸發玩家獲得「戰爭發起者」稱號廣播，製造「是我開啟了這場公會戰」的成就感
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyGuildWarPersonalCD  = 35 * time.Second // 個人冷卻
	LuckyGuildWarGlobalCD    = 55 * time.Second // 全服冷卻
	LuckyGuildWarDuration    = 30 * time.Second // 公會戰時限
	LuckyGuildWarScoreTick   = 5 * time.Second  // 比分廣播間隔
	LuckyGuildWarWinMult     = 2.5              // 勝隊倍率加成
	LuckyGuildWarLoseMult    = 1.2              // 敗隊安慰獎倍率
	LuckyGuildWarDrawMult    = 1.8              // 平局倍率加成
	LuckyGuildWarBoostSec    = 5               // 倍率加成持續秒數
)

// guildWarTeam 公會戰隊伍
type guildWarTeam int

const (
	GuildWarTeamRed  guildWarTeam = 0 // 紅隊（玩家 ID 最後一位為偶數）
	GuildWarTeamBlue guildWarTeam = 1 // 藍隊（玩家 ID 最後一位為奇數）
)

// guildWarSession 公會戰會話
type guildWarSession struct {
	triggerPlayerID   string
	triggerPlayerName string
	expiresAt         time.Time
	redScore          int
	blueScore         int
	mu                sync.Mutex
}

// guildWarBoostEntry 公會戰倍率加成記錄
type guildWarBoostEntry struct {
	mult      float64
	expiresAt time.Time
}

// luckyGuildWarManager 幸運公會戰魚管理器
type luckyGuildWarManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前活躍的公會戰會話（nil = 無）
	activeSession *guildWarSession

	// 公會戰倍率加成（playerID → boostEntry）
	warBoosts map[string]guildWarBoostEntry
}

func newLuckyGuildWarManager() *luckyGuildWarManager {
	return &luckyGuildWarManager{
		personalCooldowns: make(map[string]time.Time),
		warBoosts:         make(map[string]guildWarBoostEntry),
	}
}

// isLuckyGuildWarFish 判斷是否為幸運公會戰魚
func isLuckyGuildWarFish(defID string) bool {
	return defID == "T215"
}

// getPlayerTeam 依玩家 ID 最後一個字元的奇偶決定隊伍
func getPlayerTeam(playerID string) guildWarTeam {
	if len(playerID) == 0 {
		return GuildWarTeamRed
	}
	lastChar := playerID[len(playerID)-1]
	if (lastChar-'0')%2 == 0 || lastChar == 'a' || lastChar == 'c' || lastChar == 'e' ||
		lastChar == 'A' || lastChar == 'C' || lastChar == 'E' {
		return GuildWarTeamRed
	}
	return GuildWarTeamBlue
}

// getLuckyGuildWarBoostMult 取得公會戰倍率加成（供 handleKill 使用）
func (m *luckyGuildWarManager) getLuckyGuildWarBoostMult(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	entry, ok := m.warBoosts[playerID]
	if !ok {
		return 1.0
	}
	if time.Now().After(entry.expiresAt) {
		delete(m.warBoosts, playerID)
		return 1.0
	}
	return entry.mult
}

// isGuildWarActive 判斷公會戰是否進行中
func (m *luckyGuildWarManager) isGuildWarActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeSession == nil {
		return false
	}
	return time.Now().Before(m.activeSession.expiresAt)
}

// notifyLuckyGuildWarKill 任何玩家擊破任何目標時，若公會戰進行中則累積積分
// 由 handleKill 呼叫（非 T215 目標）
func (g *Game) notifyLuckyGuildWarKill(p *player.Player) {
	m := g.LuckyGuildWar
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

	team := getPlayerTeam(p.ID)

	sess.mu.Lock()
	if team == GuildWarTeamRed {
		sess.redScore++
	} else {
		sess.blueScore++
	}
	red := sess.redScore
	blue := sess.blueScore
	sess.mu.Unlock()

	teamName := "🔴 紅隊"
	if team == GuildWarTeamBlue {
		teamName = "🔵 藍隊"
	}
	log.Printf("[GuildWar] %s +1！紅隊 %d vs 藍隊 %d", teamName, red, blue)
}

// tryLuckyGuildWarFish 擊破 T215 後觸發公會戰
func (g *Game) tryLuckyGuildWarFish(p *player.Player) {
	m := g.LuckyGuildWar

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
	// 已有活躍公會戰
	if m.activeSession != nil && now.Before(m.activeSession.expiresAt) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyGuildWarPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyGuildWarGlobalCD)

	expiresAt := now.Add(LuckyGuildWarDuration)
	sess := &guildWarSession{
		triggerPlayerID:   p.ID,
		triggerPlayerName: p.DisplayName,
		expiresAt:         expiresAt,
		redScore:          0,
		blueScore:         0,
	}
	m.activeSession = sess
	m.mu.Unlock()

	// 計算觸發玩家的隊伍
	triggerTeam := getPlayerTeam(p.ID)
	triggerTeamName := "🔴 紅隊"
	if triggerTeam == GuildWarTeamBlue {
		triggerTeamName = "🔵 藍隊"
	}

	log.Printf("[GuildWar] player=%s（%s）觸發公會戰！時限 %ds",
		p.ID, triggerTeamName, int(LuckyGuildWarDuration.Seconds()))

	// 個人訊息：戰爭發起者
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyGuildWar,
		Payload: ws.LuckyGuildWarPayload{
			Event:        "war_start",
			PlayerID:     p.ID,
			PlayerName:   p.DisplayName,
			PlayerTeam:   int(triggerTeam),
			DurationSec:  int(LuckyGuildWarDuration.Seconds()),
			WinMult:      LuckyGuildWarWinMult,
			LoseMult:     LuckyGuildWarLoseMult,
			DrawMult:     LuckyGuildWarDrawMult,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyGuildWar,
		Payload: ws.LuckyGuildWarPayload{
			Event:       "war_broadcast",
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyGuildWarDuration.Seconds()),
			WinMult:     LuckyGuildWarWinMult,
			LoseMult:    LuckyGuildWarLoseMult,
			DrawMult:    LuckyGuildWarDrawMult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyGuildWar, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⚔️ %s 發起公會戰！🔴紅隊 vs 🔵藍隊，30秒競爭擊破數！勝隊 ×%.1f 倍率加成！",
			p.DisplayName, LuckyGuildWarWinMult),
		"color": "#DC143C",
	})

	// 啟動比分廣播 + 結算 goroutine
	go g.runGuildWarSession(sess)
}

// runGuildWarSession 公會戰主循環（比分廣播 + 結算）
func (g *Game) runGuildWarSession(sess *guildWarSession) {
	scoreTicker := time.NewTicker(LuckyGuildWarScoreTick)
	endTimer := time.NewTimer(LuckyGuildWarDuration)
	defer scoreTicker.Stop()
	defer endTimer.Stop()

	for {
		select {
		case <-scoreTicker.C:
			// 每 5 秒廣播即時比分
			sess.mu.Lock()
			red := sess.redScore
			blue := sess.blueScore
			sess.mu.Unlock()

			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgLuckyGuildWar,
				Payload: ws.LuckyGuildWarPayload{
					Event:     "war_score",
					RedScore:  red,
					BlueScore: blue,
				},
			})

		case <-endTimer.C:
			// 公會戰結束，結算
			m := g.LuckyGuildWar
			m.mu.Lock()
			if m.activeSession != sess {
				m.mu.Unlock()
				return
			}
			m.activeSession = nil
			m.mu.Unlock()

			sess.mu.Lock()
			red := sess.redScore
			blue := sess.blueScore
			sess.mu.Unlock()

			g.doGuildWarSettle(sess, red, blue)
			return

		case <-g.stopCh:
			return
		}
	}
}

// doGuildWarSettle 公會戰結算
func (g *Game) doGuildWarSettle(sess *guildWarSession, redScore, blueScore int) {
	m := g.LuckyGuildWar

	var winTeam guildWarTeam
	var winMult, loseMult float64
	var resultEvent string
	var resultMsg string

	if redScore > blueScore {
		// 紅隊勝
		winTeam = GuildWarTeamRed
		winMult = LuckyGuildWarWinMult
		loseMult = LuckyGuildWarLoseMult
		resultEvent = "war_result"
		resultMsg = fmt.Sprintf("⚔️ 公會戰結束！🔴紅隊勝利！%d vs %d！紅隊全員 ×%.1f 倍率加成 %ds！",
			redScore, blueScore, winMult, LuckyGuildWarBoostSec)
	} else if blueScore > redScore {
		// 藍隊勝
		winTeam = GuildWarTeamBlue
		winMult = LuckyGuildWarWinMult
		loseMult = LuckyGuildWarLoseMult
		resultEvent = "war_result"
		resultMsg = fmt.Sprintf("⚔️ 公會戰結束！🔵藍隊勝利！%d vs %d！藍隊全員 ×%.1f 倍率加成 %ds！",
			blueScore, redScore, winMult, LuckyGuildWarBoostSec)
	} else {
		// 平局
		winTeam = GuildWarTeamRed // 平局時兩隊都算勝
		winMult = LuckyGuildWarDrawMult
		loseMult = LuckyGuildWarDrawMult
		resultEvent = "war_draw"
		resultMsg = fmt.Sprintf("⚔️ 公會戰平局！%d vs %d！雙隊全員 ×%.1f 倍率加成 %ds！",
			redScore, blueScore, LuckyGuildWarDrawMult, LuckyGuildWarBoostSec)
	}

	// 套用倍率加成給所有玩家
	boostExpiry := time.Now().Add(time.Duration(LuckyGuildWarBoostSec) * time.Second)
	g.mu.RLock()
	players := make([]*player.Player, 0, len(g.Players))
	for _, p := range g.Players {
		players = append(players, p)
	}
	g.mu.RUnlock()

	m.mu.Lock()
	for _, p := range players {
		team := getPlayerTeam(p.ID)
		var mult float64
		if redScore == blueScore {
			// 平局：雙隊都用 drawMult
			mult = LuckyGuildWarDrawMult
		} else if team == winTeam {
			mult = winMult
		} else {
			mult = loseMult
		}
		m.warBoosts[p.ID] = guildWarBoostEntry{
			mult:      mult,
			expiresAt: boostExpiry,
		}
	}
	m.mu.Unlock()

	log.Printf("[GuildWar] 結算！紅隊 %d vs 藍隊 %d，勝隊 ×%.1f，敗隊 ×%.1f",
		redScore, blueScore, winMult, loseMult)

	// 全服廣播結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyGuildWar,
		Payload: ws.LuckyGuildWarPayload{
			Event:       resultEvent,
			PlayerName:  sess.triggerPlayerName,
			RedScore:    redScore,
			BlueScore:   blueScore,
			WinTeam:     int(winTeam),
			WinMult:     winMult,
			LoseMult:    loseMult,
			BoostSec:    LuckyGuildWarBoostSec,
		},
	})

	// 全服公告
	color := "#DC143C"
	if blueScore > redScore {
		color = "#1E90FF"
	} else if redScore == blueScore {
		color = "#9370DB"
	}
	g.Announce.Create(announce.EventLuckyGuildWar, sess.triggerPlayerName, 0, map[string]string{
		"message": resultMsg,
		"color":   color,
	})

	_ = winTeam // suppress unused warning if draw path sets winTeam=Red
}
