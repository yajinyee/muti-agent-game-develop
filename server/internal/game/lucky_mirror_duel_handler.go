// lucky_mirror_duel_handler.go — 幸運鏡像對決魚系統（DAY-270）
// 業界依據：2026 年最熱門「PvP 鏡像對決」機制
//
// 設計：擊破 T228 後，觸發「鏡像對決」：
//   - Server 隨機選一個其他玩家作為「對手」
//   - 雙方進入 15 秒對決期
//   - 對決期間，雙方每次擊破目標，對手也獲得相同獎勵的 50%（鏡像分享）
//   - 15 秒後，擊破數多的玩家獲得「對決勝利」×2.0 加成（5 秒）
//   - 擊破數少的玩家獲得「對決失敗」×1.2 安慰獎（5 秒）
//   - 平局：雙方各獲得 ×1.5 加成（5 秒）
//   - 若無其他玩家，觸發「孤獨模式」：個人 ×1.5 加成 10 秒
//   - 個人冷卻 30 秒；全服冷卻 50 秒
//
// 設計差異：
//   - 與公會戰（T215，全服分隊）不同，鏡像對決是「1v1 個人對決」，讓玩家有「我要打贏對手」的直接競爭感
//   - 「鏡像分享 50%」讓雙方都有動力打魚，不是零和遊戲
//   - 「勝利 ×2.0」是個人競爭類最高倍率，製造「要趁 15 秒內打贏對手」的緊迫感
//   - 「平局 ×1.5」鼓勢均力敵，不讓任何一方感到被碾壓
//   - 「孤獨模式 ×1.5」確保單人遊戲也有收益
//   - 「全服廣播對決結果」讓所有玩家看到「誰贏了對決」，製造社交話題感
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
	LuckyMirrorDuelPersonalCD  = 30 * time.Second // 個人冷卻
	LuckyMirrorDuelGlobalCD    = 50 * time.Second // 全服冷卻
	LuckyMirrorDuelDuration    = 15 * time.Second // 對決持續時間
	LuckyMirrorDuelBoostDur    = 5 * time.Second  // 結算加成持續時間
	LuckyMirrorDuelSoloDur     = 10 * time.Second // 孤獨模式持續時間
	LuckyMirrorDuelShareRatio  = 0.5              // 鏡像分享比例
	LuckyMirrorDuelWinMult     = 2.0              // 勝利倍率
	LuckyMirrorDuelLoseMult    = 1.2              // 失敗安慰倍率
	LuckyMirrorDuelDrawMult    = 1.5              // 平局倍率
	LuckyMirrorDuelSoloMult    = 1.5              // 孤獨模式倍率
)

// mirrorDuelSession 對決 session
type mirrorDuelSession struct {
	player1ID   string
	player1Name string
	player2ID   string
	player2Name string
	expiresAt   time.Time
	score1      int // player1 擊破數
	score2      int // player2 擊破數
	settled     bool
}

// mirrorDuelBoost 對決結算加成
type mirrorDuelBoost struct {
	playerID  string
	mult      float64
	expiresAt time.Time
}

// luckyMirrorDuelManager 幸運鏡像對決魚管理器
type luckyMirrorDuelManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 活躍對決 session（sessionID → session）
	// sessionID = player1ID + ":" + player2ID
	activeSessions map[string]*mirrorDuelSession

	// 結算加成（playerID → boost）
	activeBoosts map[string]*mirrorDuelBoost
}

func newLuckyMirrorDuelManager() *luckyMirrorDuelManager {
	return &luckyMirrorDuelManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*mirrorDuelSession),
		activeBoosts:      make(map[string]*mirrorDuelBoost),
	}
}

// isLuckyMirrorDuelFish 判斷是否為幸運鏡像對決魚
func isLuckyMirrorDuelFish(defID string) bool {
	return defID == "T228"
}

// getMirrorDuelBoostMult 取得對決結算加成倍率（供 handleKill 使用）
func (m *luckyMirrorDuelManager) getMirrorDuelBoostMult(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	boost, ok := m.activeBoosts[playerID]
	if !ok {
		return 1.0
	}
	if time.Now().After(boost.expiresAt) {
		delete(m.activeBoosts, playerID)
		return 1.0
	}
	return boost.mult
}

// getActiveDuelSession 取得玩家所在的對決 session
func (m *luckyMirrorDuelManager) getActiveDuelSession(playerID string) (*mirrorDuelSession, string) {
	for sessionID, s := range m.activeSessions {
		if (s.player1ID == playerID || s.player2ID == playerID) && !s.settled {
			if time.Now().Before(s.expiresAt) {
				return s, sessionID
			}
		}
	}
	return nil, ""
}

// tryLuckyMirrorDuelFish 擊破 T228 後觸發鏡像對決
func (g *Game) tryLuckyMirrorDuelFish(p *player.Player) {
	m := g.LuckyMirrorDuel

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
	m.personalCooldowns[p.ID] = now.Add(LuckyMirrorDuelPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyMirrorDuelGlobalCD)
	m.mu.Unlock()

	// 找一個其他玩家作為對手
	g.mu.RLock()
	var opponent *player.Player
	candidates := make([]*player.Player, 0)
	for _, other := range g.Players {
		if other.ID != p.ID {
			candidates = append(candidates, other)
		}
	}
	g.mu.RUnlock()

	if len(candidates) > 0 {
		rng := rand.New(rand.NewSource(now.UnixNano()))
		opponent = candidates[rng.Intn(len(candidates))]
	}

	if opponent == nil {
		// 孤獨模式：無其他玩家
		g.runMirrorDuelSoloMode(p, now)
		return
	}

	// 建立對決 session
	sessionID := p.ID + ":" + opponent.ID
	session := &mirrorDuelSession{
		player1ID:   p.ID,
		player1Name: p.DisplayName,
		player2ID:   opponent.ID,
		player2Name: opponent.DisplayName,
		expiresAt:   now.Add(LuckyMirrorDuelDuration),
	}

	m.mu.Lock()
	m.activeSessions[sessionID] = session
	m.mu.Unlock()

	log.Printf("[MirrorDuel] 對決開始！%s vs %s（session=%s）",
		p.DisplayName, opponent.DisplayName, sessionID)

	// 通知 player1（觸發者）
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyMirrorDuel,
		Payload: ws.LuckyMirrorDuelPayload{
			Event:       "duel_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			OpponentID:  opponent.ID,
			OpponentName: opponent.DisplayName,
			Duration:    LuckyMirrorDuelDuration.Seconds(),
			ShareRatio:  LuckyMirrorDuelShareRatio,
			IsChallenger: true,
		},
	})

	// 通知 player2（被挑戰者）
	_ = g.Hub.Send(opponent.ID, &ws.Message{
		Type: ws.MsgLuckyMirrorDuel,
		Payload: ws.LuckyMirrorDuelPayload{
			Event:       "duel_start",
			PlayerID:    opponent.ID,
			PlayerName:  opponent.DisplayName,
			OpponentID:  p.ID,
			OpponentName: p.DisplayName,
			Duration:    LuckyMirrorDuelDuration.Seconds(),
			ShareRatio:  LuckyMirrorDuelShareRatio,
			IsChallenger: false,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMirrorDuel,
		Payload: ws.LuckyMirrorDuelPayload{
			Event:       "duel_broadcast",
			PlayerName:  p.DisplayName,
			OpponentName: opponent.DisplayName,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyMirrorDuel, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🪞 %s 向 %s 發起鏡像對決！15 秒決勝負！",
			p.DisplayName, opponent.DisplayName),
		"color": "#9B59B6",
	})

	// 啟動對決計時器
	go g.runMirrorDuelTimer(session, sessionID)
}

// runMirrorDuelSoloMode 孤獨模式（無其他玩家）
func (g *Game) runMirrorDuelSoloMode(p *player.Player, now time.Time) {
	m := g.LuckyMirrorDuel

	boost := &mirrorDuelBoost{
		playerID:  p.ID,
		mult:      LuckyMirrorDuelSoloMult,
		expiresAt: now.Add(LuckyMirrorDuelSoloDur),
	}
	m.mu.Lock()
	m.activeBoosts[p.ID] = boost
	m.mu.Unlock()

	log.Printf("[MirrorDuel] player=%s 孤獨模式 ×%.1f（%v）",
		p.ID, LuckyMirrorDuelSoloMult, LuckyMirrorDuelSoloDur)

	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyMirrorDuel,
		Payload: ws.LuckyMirrorDuelPayload{
			Event:      "duel_solo",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			SoloMult:   LuckyMirrorDuelSoloMult,
			Duration:   LuckyMirrorDuelSoloDur.Seconds(),
		},
	})

	// 孤獨模式超時清理
	go func() {
		time.Sleep(LuckyMirrorDuelSoloDur)
		m.mu.Lock()
		if b, ok := m.activeBoosts[p.ID]; ok && b == boost {
			delete(m.activeBoosts, p.ID)
		}
		m.mu.Unlock()
	}()
}

// notifyMirrorDuelKill 對決期間擊破目標時呼叫（由 handleKill 呼叫）
// 回傳鏡像分享獎勵（給對手）
func (g *Game) notifyMirrorDuelKill(p *player.Player, reward int) int {
	m := g.LuckyMirrorDuel

	m.mu.Lock()
	session, sessionID := m.getActiveDuelSession(p.ID)
	if session == nil {
		m.mu.Unlock()
		return 0
	}

	// 更新積分
	var opponentID string
	var opponentName string
	if session.player1ID == p.ID {
		session.score1++
		opponentID = session.player2ID
		opponentName = session.player2Name
	} else {
		session.score2++
		opponentID = session.player1ID
		opponentName = session.player1Name
	}
	score1 := session.score1
	score2 := session.score2
	_ = sessionID
	m.mu.Unlock()

	// 計算鏡像分享獎勵
	shareReward := int(float64(reward) * LuckyMirrorDuelShareRatio)
	if shareReward < 1 {
		shareReward = 1
	}

	// 通知擊破者（積分更新）
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyMirrorDuel,
		Payload: ws.LuckyMirrorDuelPayload{
			Event:        "duel_score",
			PlayerID:     p.ID,
			MyScore:      score1,
			OpponentScore: score2,
			ShareReward:  shareReward,
		},
	})

	// 通知對手（收到鏡像分享獎勵）
	_ = g.Hub.Send(opponentID, &ws.Message{
		Type: ws.MsgLuckyMirrorDuel,
		Payload: ws.LuckyMirrorDuelPayload{
			Event:        "duel_mirror_reward",
			PlayerID:     opponentID,
			OpponentName: p.DisplayName,
			ShareReward:  shareReward,
			MyScore:      score2,
			OpponentScore: score1,
		},
	})

	log.Printf("[MirrorDuel] %s 擊破目標，鏡像分享 %d 給 %s（比分 %d:%d）",
		p.DisplayName, shareReward, opponentName, score1, score2)

	return shareReward
}

// runMirrorDuelTimer 對決計時器 goroutine
func (g *Game) runMirrorDuelTimer(session *mirrorDuelSession, sessionID string) {
	timer := time.NewTimer(LuckyMirrorDuelDuration)
	defer timer.Stop()

	<-timer.C

	g.doMirrorDuelSettle(session, sessionID)
}

// doMirrorDuelSettle 對決結算
func (g *Game) doMirrorDuelSettle(session *mirrorDuelSession, sessionID string) {
	m := g.LuckyMirrorDuel

	m.mu.Lock()
	if session.settled {
		m.mu.Unlock()
		return
	}
	session.settled = true
	score1 := session.score1
	score2 := session.score2
	delete(m.activeSessions, sessionID)
	m.mu.Unlock()

	now := time.Now()

	var winner, loser string
	var winnerName, loserName string
	var winMult, loseMult float64
	var resultMsg string
	var resultColor string

	if score1 > score2 {
		// player1 勝利
		winner = session.player1ID
		winnerName = session.player1Name
		loser = session.player2ID
		loserName = session.player2Name
		winMult = LuckyMirrorDuelWinMult
		loseMult = LuckyMirrorDuelLoseMult
		resultMsg = fmt.Sprintf("🪞 %s 對決勝利！（%d:%d）×%.1f 加成！",
			winnerName, score1, score2, winMult)
		resultColor = "#FFD700"
	} else if score2 > score1 {
		// player2 勝利
		winner = session.player2ID
		winnerName = session.player2Name
		loser = session.player1ID
		loserName = session.player1Name
		winMult = LuckyMirrorDuelWinMult
		loseMult = LuckyMirrorDuelLoseMult
		resultMsg = fmt.Sprintf("🪞 %s 對決勝利！（%d:%d）×%.1f 加成！",
			winnerName, score2, score1, winMult)
		resultColor = "#FFD700"
	} else {
		// 平局
		winner = ""
		loser = ""
		winnerName = ""
		loserName = ""
		winMult = LuckyMirrorDuelDrawMult
		loseMult = LuckyMirrorDuelDrawMult
		resultMsg = fmt.Sprintf("🪞 %s vs %s 平局！雙方各獲得 ×%.1f 加成！",
			session.player1Name, session.player2Name, winMult)
		resultColor = "#9B59B6"
	}

	log.Printf("[MirrorDuel] 對決結算！%s(%d) vs %s(%d) → winner=%s",
		session.player1Name, score1, session.player2Name, score2, winnerName)

	// 設定加成
	boostExpiry := now.Add(LuckyMirrorDuelBoostDur)

	if winner != "" {
		// 有勝負
		m.mu.Lock()
		m.activeBoosts[winner] = &mirrorDuelBoost{
			playerID:  winner,
			mult:      winMult,
			expiresAt: boostExpiry,
		}
		m.activeBoosts[loser] = &mirrorDuelBoost{
			playerID:  loser,
			mult:      loseMult,
			expiresAt: boostExpiry,
		}
		m.mu.Unlock()

		// 通知勝者
		_ = g.Hub.Send(winner, &ws.Message{
			Type: ws.MsgLuckyMirrorDuel,
			Payload: ws.LuckyMirrorDuelPayload{
				Event:        "duel_result",
				PlayerID:     winner,
				IsWinner:     true,
				MyScore:      score1,
				OpponentScore: score2,
				ResultMult:   winMult,
				Duration:     LuckyMirrorDuelBoostDur.Seconds(),
			},
		})

		// 通知敗者
		_ = g.Hub.Send(loser, &ws.Message{
			Type: ws.MsgLuckyMirrorDuel,
			Payload: ws.LuckyMirrorDuelPayload{
				Event:        "duel_result",
				PlayerID:     loser,
				IsWinner:     false,
				MyScore:      score2,
				OpponentScore: score1,
				ResultMult:   loseMult,
				Duration:     LuckyMirrorDuelBoostDur.Seconds(),
			},
		})

		_ = loserName
	} else {
		// 平局
		m.mu.Lock()
		m.activeBoosts[session.player1ID] = &mirrorDuelBoost{
			playerID:  session.player1ID,
			mult:      winMult,
			expiresAt: boostExpiry,
		}
		m.activeBoosts[session.player2ID] = &mirrorDuelBoost{
			playerID:  session.player2ID,
			mult:      winMult,
			expiresAt: boostExpiry,
		}
		m.mu.Unlock()

		// 通知雙方平局
		for _, pid := range []string{session.player1ID, session.player2ID} {
			_ = g.Hub.Send(pid, &ws.Message{
				Type: ws.MsgLuckyMirrorDuel,
				Payload: ws.LuckyMirrorDuelPayload{
					Event:        "duel_draw",
					PlayerID:     pid,
					MyScore:      score1,
					OpponentScore: score2,
					ResultMult:   winMult,
					Duration:     LuckyMirrorDuelBoostDur.Seconds(),
				},
			})
		}
	}

	// 全服廣播結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMirrorDuel,
		Payload: ws.LuckyMirrorDuelPayload{
			Event:        "duel_settle_broadcast",
			PlayerName:   session.player1Name,
			OpponentName: session.player2Name,
			Score1:       score1,
			Score2:       score2,
			WinnerName:   winnerName,
			WinMult:      winMult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyMirrorDuel, winnerName, 0, map[string]string{
		"message": resultMsg,
		"color":   resultColor,
	})

	// 加成超時清理
	go func() {
		time.Sleep(LuckyMirrorDuelBoostDur)
		m.mu.Lock()
		for _, pid := range []string{session.player1ID, session.player2ID} {
			if b, ok := m.activeBoosts[pid]; ok && time.Now().After(b.expiresAt) {
				delete(m.activeBoosts, pid)
			}
		}
		m.mu.Unlock()
	}()
}
