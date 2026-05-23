// lucky_flag_fish_handler.go — 幸運奪旗魚系統（DAY-244）
// 業界原創「全服搶旗競爭」機制
//
// 設計：擊破 T202 後，場上隨機 1 個目標被「旗幟標記」（持續 15 秒）：
//   - 所有玩家射擊旗幟目標，每次命中累積「搶旗積分」（+1/命中，不消耗籌碼）
//   - 每 3 秒廣播即時排名，製造「全服競爭」的緊張感
//   - 15 秒後，積分最高的玩家「奪旗成功」→ 獲得 ×4.0 倍率加成 + 全服廣播英雄稱號
//   - 第 2 名 ×2.0，第 3 名 ×1.5（安慰獎）
//   - 若無人命中 → 旗幟目標自動爆炸（全服共享 ×1.2 倍率）
//   - 個人冷卻 25 秒；全服冷卻 40 秒
//
// 設計差異：
//   - 與幸運拍賣魚（T175，消耗籌碼競標）不同，奪旗魚是「射擊積分競爭」，不消耗籌碼，讓所有玩家都願意參與
//   - 「每 3 秒排名廣播」讓玩家即時看到自己的排名，製造「要趕快多打幾槍」的緊迫感
//   - 「第 2/3 名也有獎勵」讓玩家不會因為落後就放棄，保持全程參與度
//   - 「旗幟目標自動爆炸」確保即使沒人積極參與也有獎勵，不浪費機制
//   - ×4.0 倍率是目前全服競爭類最高倍率，讓玩家有「值得全力搶」的動機
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyFlagFishPersonalCD   = 25 * time.Second // 個人冷卻
	LuckyFlagFishGlobalCD     = 40 * time.Second // 全服冷卻
	LuckyFlagFishDuration     = 15 * time.Second // 搶旗持續時間
	LuckyFlagFishRankInterval = 3 * time.Second  // 排名廣播間隔
	LuckyFlagFishWinnerMult   = 4.0              // 第 1 名倍率
	LuckyFlagFishSecondMult   = 2.0              // 第 2 名倍率
	LuckyFlagFishThirdMult    = 1.5              // 第 3 名倍率
	LuckyFlagFishAutoBlastMult = 1.2             // 無人命中時自動爆炸倍率
)

// flagFishSession 搶旗 session
type flagFishSession struct {
	targetID    string
	targetDefID string
	targetX     float64
	targetY     float64
	startedAt   time.Time
	expiresAt   time.Time
	// 積分（playerID → score）
	scores map[string]int
	// 玩家名稱快取（playerID → displayName）
	names map[string]string
	mu    sync.Mutex
}

// luckyFlagFishManager 幸運奪旗魚管理器
type luckyFlagFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前搶旗 session（nil = 無進行中）
	activeSession *flagFishSession
}

func newLuckyFlagFishManager() *luckyFlagFishManager {
	return &luckyFlagFishManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyFlagFish 判斷是否為幸運奪旗魚
func isLuckyFlagFish(defID string) bool {
	return defID == "T202"
}

// isFlagTarget 判斷某個目標是否為當前旗幟目標
func (g *Game) isFlagTarget(instanceID string) bool {
	mgr := g.LuckyFlagFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if mgr.activeSession == nil {
		return false
	}
	if time.Now().After(mgr.activeSession.expiresAt) {
		return false
	}
	return mgr.activeSession.targetID == instanceID
}

// recordFlagHit 記錄玩家命中旗幟目標（供 handleAttack 使用）
func (g *Game) recordFlagHit(playerID, playerName string) {
	mgr := g.LuckyFlagFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if mgr.activeSession == nil || time.Now().After(mgr.activeSession.expiresAt) {
		return
	}

	sess := mgr.activeSession
	sess.mu.Lock()
	sess.scores[playerID]++
	sess.names[playerID] = playerName
	score := sess.scores[playerID]
	sess.mu.Unlock()

	log.Printf("[LuckyFlag] player=%s hit flag target, score=%d", playerID, score)
}

// getFlagWinnerMult 取得奪旗勝利倍率（供 handleKill 使用）
// 若玩家擊破旗幟目標，根據積分排名回傳對應倍率
func (g *Game) getFlagWinnerMult(playerID, instanceID string) float64 {
	mgr := g.LuckyFlagFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if mgr.activeSession == nil || mgr.activeSession.targetID != instanceID {
		return 1.0
	}

	sess := mgr.activeSession
	sess.mu.Lock()
	rank := g.calcFlagRankLocked(sess, playerID)
	sess.mu.Unlock()

	switch rank {
	case 1:
		return LuckyFlagFishWinnerMult
	case 2:
		return LuckyFlagFishSecondMult
	case 3:
		return LuckyFlagFishThirdMult
	default:
		return 1.0
	}
}

// calcFlagRankLocked 計算玩家在搶旗中的排名（需持有 sess.mu）
func (g *Game) calcFlagRankLocked(sess *flagFishSession, playerID string) int {
	myScore := sess.scores[playerID]
	if myScore == 0 {
		return 0
	}
	rank := 1
	for pid, score := range sess.scores {
		if pid != playerID && score > myScore {
			rank++
		}
	}
	return rank
}

// notifyFlagTargetKill 旗幟目標被擊破時呼叫（廣播奪旗結算）
func (g *Game) notifyFlagTargetKill(p *player.Player, instanceID string, reward int) {
	mgr := g.LuckyFlagFish
	mgr.mu.Lock()

	if mgr.activeSession == nil || mgr.activeSession.targetID != instanceID {
		mgr.mu.Unlock()
		return
	}

	sess := mgr.activeSession
	mgr.activeSession = nil
	mgr.mu.Unlock()

	// 計算排名
	sess.mu.Lock()
	type rankEntry struct {
		playerID   string
		playerName string
		score      int
		rank       int
	}
	var entries []rankEntry
	for pid, score := range sess.scores {
		entries = append(entries, rankEntry{
			playerID:   pid,
			playerName: sess.names[pid],
			score:      score,
		})
	}
	sess.mu.Unlock()

	// 按積分排序
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].score > entries[j].score
	})
	for i := range entries {
		entries[i].rank = i + 1
	}

	// 確認擊破者的排名
	killerRank := 0
	for _, e := range entries {
		if e.playerID == p.ID {
			killerRank = e.rank
			break
		}
	}

	var killerMult float64
	switch killerRank {
	case 1:
		killerMult = LuckyFlagFishWinnerMult
	case 2:
		killerMult = LuckyFlagFishSecondMult
	case 3:
		killerMult = LuckyFlagFishThirdMult
	default:
		killerMult = 1.0
	}

	log.Printf("[LuckyFlag] flag target killed by player=%s rank=%d mult=%.1f reward=%d",
		p.ID, killerRank, killerMult, reward)

	// 廣播奪旗結算
	var rankList []map[string]interface{}
	for _, e := range entries {
		var mult float64
		switch e.rank {
		case 1:
			mult = LuckyFlagFishWinnerMult
		case 2:
			mult = LuckyFlagFishSecondMult
		case 3:
			mult = LuckyFlagFishThirdMult
		default:
			mult = 1.0
		}
		rankList = append(rankList, map[string]interface{}{
			"player_id":   e.playerID,
			"player_name": e.playerName,
			"score":       e.score,
			"rank":        e.rank,
			"mult":        mult,
		})
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyFlagFish,
		Payload: ws.LuckyFlagFishPayload{
			Event:      "flag_captured",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			TargetID:   instanceID,
			Reward:     reward,
			KillerMult: killerMult,
			KillerRank: killerRank,
			RankList:   rankList,
		},
	})

	// 全服公告
	if killerRank == 1 {
		ann := g.Announce.Create(announce.EventLuckyFlagFish, p.DisplayName, reward, map[string]string{
			"message": fmt.Sprintf("🚩 %s 奪旗成功！×%.1f 倍率！獲得 %d 金幣！",
				p.DisplayName, killerMult, reward),
			"color": "#E74C3C",
		})
		g.broadcastAnnouncement(ann)
	}
}

// notifyFlagTargetGone 旗幟目標消失時呼叫（觸發自動爆炸）
func (g *Game) notifyFlagTargetGone(instanceID string) {
	mgr := g.LuckyFlagFish
	mgr.mu.Lock()

	if mgr.activeSession == nil || mgr.activeSession.targetID != instanceID {
		mgr.mu.Unlock()
		return
	}

	sess := mgr.activeSession
	mgr.activeSession = nil
	mgr.mu.Unlock()

	go g.doFlagAutoBlast(sess)
}

// doFlagAutoBlast 旗幟目標自動爆炸（無人命中時）
func (g *Game) doFlagAutoBlast(sess *flagFishSession) {
	sess.mu.Lock()
	hasParticipants := len(sess.scores) > 0
	sess.mu.Unlock()

	if hasParticipants {
		// 有人參與但目標消失，給所有參與者安慰獎
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyFlagFish,
			Payload: ws.LuckyFlagFishPayload{
				Event:    "flag_escaped",
				TargetID: sess.targetID,
			},
		})
		return
	}

	// 無人命中，自動爆炸（全服共享）
	g.mu.RLock()
	var avgBet int
	playerCount := len(g.Players)
	for _, p := range g.Players {
		betDef := data.GetBetDef(p.BetLevel)
		if betDef != nil {
			avgBet += betDef.BetCost
		}
	}
	g.mu.RUnlock()

	if playerCount > 0 {
		avgBet /= playerCount
	}
	if avgBet < 1 {
		avgBet = 1
	}

	blastReward := int(float64(avgBet) * LuckyFlagFishAutoBlastMult * 10)

	// 分配給所有玩家
	g.mu.RLock()
	for _, p := range g.Players {
		p.AddCoins(blastReward)
	}
	g.mu.RUnlock()

	log.Printf("[LuckyFlag] flag auto blast! reward=%d per player", blastReward)

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyFlagFish,
		Payload: ws.LuckyFlagFishPayload{
			Event:       "flag_auto_blast",
			TargetID:    sess.targetID,
			Reward:      blastReward,
			KillerMult:  LuckyFlagFishAutoBlastMult,
		},
	})

	ann := g.Announce.Create(announce.EventLuckyFlagFish, "全服", blastReward, map[string]string{
		"message": fmt.Sprintf("🚩 旗幟目標自動爆炸！全服每人獲得 %d 金幣！", blastReward),
		"color":   "#E67E22",
	})
	g.broadcastAnnouncement(ann)
}

// tryLuckyFlagFish 擊破 T202 後觸發搶旗
func (g *Game) tryLuckyFlagFish(p *player.Player) {
	mgr := g.LuckyFlagFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}

	// 全服冷卻檢查
	if time.Now().Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}

	// 已有進行中的搶旗
	if mgr.activeSession != nil && time.Now().Before(mgr.activeSession.expiresAt) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyFlagFishPersonalCD)
	mgr.globalCooldownUntil = time.Now().Add(LuckyFlagFishGlobalCD)
	mgr.mu.Unlock()

	// 選擇旗幟目標（場上隨機一個非 BOSS 目標）
	g.mu.RLock()
	var candidates []string
	for id, t := range g.Targets {
		if t.HP > 0 && t.DefID != "B001" && t.DefID != "T202" {
			candidates = append(candidates, id)
		}
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		log.Printf("[LuckyFlag] player=%s no valid targets for flag", p.ID)
		return
	}

	// 優先選高倍率目標（讓搶旗更有價值）
	targetID := candidates[rand.Intn(len(candidates))]
	g.mu.RLock()
	for _, id := range candidates {
		t := g.Targets[id]
		cur := g.Targets[targetID]
		if t != nil && cur != nil && t.Multiplier > cur.Multiplier {
			targetID = id
		}
	}
	target := g.Targets[targetID]
	var targetDefID string
	var targetX, targetY float64
	if target != nil {
		targetDefID = target.DefID
		targetX = target.X
		targetY = target.Y
	}
	g.mu.RUnlock()

	if target == nil {
		return
	}

	// 建立搶旗 session
	expiresAt := time.Now().Add(LuckyFlagFishDuration)
	sess := &flagFishSession{
		targetID:    targetID,
		targetDefID: targetDefID,
		targetX:     targetX,
		targetY:     targetY,
		startedAt:   time.Now(),
		expiresAt:   expiresAt,
		scores:      make(map[string]int),
		names:       make(map[string]string),
	}

	mgr.mu.Lock()
	mgr.activeSession = sess
	mgr.mu.Unlock()

	log.Printf("[LuckyFlag] player=%s triggered flag on target=%s (%s) for %v",
		p.ID, targetID, targetDefID, LuckyFlagFishDuration)

	// 全服廣播：搶旗開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyFlagFish,
		Payload: ws.LuckyFlagFishPayload{
			Event:       "flag_start",
			PlayerName:  p.DisplayName,
			TargetID:    targetID,
			TargetDefID: targetDefID,
			X:           targetX,
			Y:           targetY,
			DurationSec: int(LuckyFlagFishDuration.Seconds()),
			WinnerMult:  LuckyFlagFishWinnerMult,
			SecondMult:  LuckyFlagFishSecondMult,
			ThirdMult:   LuckyFlagFishThirdMult,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyFlagFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🚩 %s 觸發奪旗！%d 秒內搶旗積分最高者獲得 ×%.1f 倍率！",
			p.DisplayName, int(LuckyFlagFishDuration.Seconds()), LuckyFlagFishWinnerMult),
		"color": "#E74C3C",
	})
	g.broadcastAnnouncement(ann)

	// 啟動排名廣播 goroutine
	go g.runFlagRankBroadcast(sess)
}

// runFlagRankBroadcast 每 3 秒廣播即時排名
func (g *Game) runFlagRankBroadcast(sess *flagFishSession) {
	ticker := time.NewTicker(LuckyFlagFishRankInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if time.Now().After(sess.expiresAt) {
				// 時間到，觸發超時結算
				g.doFlagTimeout(sess)
				return
			}

			// 廣播即時排名
			sess.mu.Lock()
			type rankEntry struct {
				playerID   string
				playerName string
				score      int
			}
			var entries []rankEntry
			for pid, score := range sess.scores {
				entries = append(entries, rankEntry{
					playerID:   pid,
					playerName: sess.names[pid],
					score:      score,
				})
			}
			sess.mu.Unlock()

			sort.Slice(entries, func(i, j int) bool {
				return entries[i].score > entries[j].score
			})

			var rankList []map[string]interface{}
			for i, e := range entries {
				rankList = append(rankList, map[string]interface{}{
					"player_id":   e.playerID,
					"player_name": e.playerName,
					"score":       e.score,
					"rank":        i + 1,
				})
			}

			remaining := int(time.Until(sess.expiresAt).Seconds())
			if remaining < 0 {
				remaining = 0
			}

			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgLuckyFlagFish,
				Payload: ws.LuckyFlagFishPayload{
					Event:       "flag_rank_update",
					TargetID:    sess.targetID,
					RankList:    rankList,
					RemainingSec: remaining,
				},
			})

		case <-g.stopCh:
			return
		}
	}
}

// doFlagTimeout 搶旗時間到，結算
func (g *Game) doFlagTimeout(sess *flagFishSession) {
	mgr := g.LuckyFlagFish
	mgr.mu.Lock()
	if mgr.activeSession == sess {
		mgr.activeSession = nil
	}
	mgr.mu.Unlock()

	sess.mu.Lock()
	hasParticipants := len(sess.scores) > 0
	sess.mu.Unlock()

	if !hasParticipants {
		// 無人參與，觸發自動爆炸
		go g.doFlagAutoBlast(sess)
		return
	}

	// 有人參與，廣播超時結算（旗幟目標仍存活，但搶旗結束）
	sess.mu.Lock()
	type rankEntry struct {
		playerID   string
		playerName string
		score      int
	}
	var entries []rankEntry
	for pid, score := range sess.scores {
		entries = append(entries, rankEntry{
			playerID:   pid,
			playerName: sess.names[pid],
			score:      score,
		})
	}
	sess.mu.Unlock()

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].score > entries[j].score
	})

	var rankList []map[string]interface{}
	for i, e := range entries {
		rank := i + 1
		var mult float64
		switch rank {
		case 1:
			mult = LuckyFlagFishWinnerMult
		case 2:
			mult = LuckyFlagFishSecondMult
		case 3:
			mult = LuckyFlagFishThirdMult
		default:
			mult = 1.0
		}
		rankList = append(rankList, map[string]interface{}{
			"player_id":   e.playerID,
			"player_name": e.playerName,
			"score":       e.score,
			"rank":        rank,
			"mult":        mult,
		})
	}

	log.Printf("[LuckyFlag] flag timeout! winner=%s score=%d",
		func() string {
			if len(entries) > 0 {
				return entries[0].playerName
			}
			return "none"
		}(),
		func() int {
			if len(entries) > 0 {
				return entries[0].score
			}
			return 0
		}())

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyFlagFish,
		Payload: ws.LuckyFlagFishPayload{
			Event:    "flag_timeout",
			TargetID: sess.targetID,
			RankList: rankList,
		},
	})

	// 全服公告（第 1 名）
	if len(entries) > 0 {
		ann := g.Announce.Create(announce.EventLuckyFlagFish, entries[0].playerName, 0, map[string]string{
			"message": fmt.Sprintf("🚩 搶旗結束！%s 積分最高（%d 分）！下次擊破旗幟目標獲得 ×%.1f 倍率！",
				entries[0].playerName, entries[0].score, LuckyFlagFishWinnerMult),
			"color": "#E74C3C",
		})
		g.broadcastAnnouncement(ann)
	}
}
