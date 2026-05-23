// lucky_time_rewind_handler.go — 幸運時光倒流魚系統（DAY-247）
// 業界原創「時光倒流+過去擊破重現」機制
//
// 設計：擊破 T205 後，Server 重播玩家「過去 10 秒內」擊破的最多 5 個目標：
//   - 每個重播目標以 ×1.6 倍率給予個人獎勵（不需要再次射擊，直接結算）
//   - 同時場上所有目標 HP 恢復到 60%（讓玩家有更多目標可打）
//   - 重播動畫：每個目標間隔 400ms 依序「閃現→爆炸」，製造「時光倒流」的視覺感
//   - 個人冷卻 25 秒；全服冷卻 40 秒
//
// 設計差異：
//   - 與時間凍結魚（DAY-212，全場靜止）不同，時光倒流是「重現過去擊破」
//     讓玩家有「剛才打的那些魚又回來了」的驚喜感
//   - 與鏡面時空魚（DAY-227，鏡像+時間）不同，時光倒流是「真正的歷史重播」
//     讓玩家有「我的每一槍都有意義」的成就感
//   - 「HP 恢復 60%」讓玩家在時光倒流後有更多目標可打，延長爽感
//   - 「間隔 400ms 依序爆炸」讓玩家看到「時光倒流」的視覺效果，不是一次全爆
//   - 全服廣播讓其他玩家看到「有人觸發了時光倒流」，製造羨慕感
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyTimeRewindPersonalCD   = 25 * time.Second // 個人冷卻
	LuckyTimeRewindGlobalCD     = 40 * time.Second // 全服冷卻
	LuckyTimeRewindHistoryWindow = 10 * time.Second // 歷史記錄視窗
	LuckyTimeRewindMaxReplay    = 5                 // 最多重播目標數
	LuckyTimeRewindReplayMult   = 1.6               // 重播倍率
	LuckyTimeRewindHPRestore    = 0.6               // HP 恢復比例
	LuckyTimeRewindReplayDelay  = 400 * time.Millisecond // 每個重播間隔
)

// killRecord 擊破記錄
type killRecord struct {
	instanceID string
	defID      string
	name       string
	mult       int    // 目標倍率
	killedAt   time.Time
}

// luckyTimeRewindManager 幸運時光倒流魚管理器
type luckyTimeRewindManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 玩家擊破歷史（playerID → []killRecord，最近 10 秒）
	killHistory map[string][]killRecord
}

func newLuckyTimeRewindManager() *luckyTimeRewindManager {
	return &luckyTimeRewindManager{
		personalCooldowns: make(map[string]time.Time),
		killHistory:       make(map[string][]killRecord),
	}
}

// isLuckyTimeRewindFish 判斷是否為幸運時光倒流魚
func isLuckyTimeRewindFish(defID string) bool {
	return defID == "T205"
}

// recordKillHistory 記錄玩家擊破歷史（由 handleKill 呼叫）
func (m *luckyTimeRewindManager) recordKillHistory(playerID, instanceID, defID, name string, mult int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-LuckyTimeRewindHistoryWindow)

	// 清理過期記錄
	history := m.killHistory[playerID]
	valid := history[:0]
	for _, r := range history {
		if r.killedAt.After(cutoff) {
			valid = append(valid, r)
		}
	}

	// 追加新記錄
	valid = append(valid, killRecord{
		instanceID: instanceID,
		defID:      defID,
		name:       name,
		mult:       mult,
		killedAt:   now,
	})
	m.killHistory[playerID] = valid
}

// getRecentKills 取得玩家最近 N 筆擊破記錄
func (m *luckyTimeRewindManager) getRecentKills(playerID string, maxCount int) []killRecord {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-LuckyTimeRewindHistoryWindow)

	history := m.killHistory[playerID]
	var valid []killRecord
	for _, r := range history {
		if r.killedAt.After(cutoff) {
			valid = append(valid, r)
		}
	}
	m.killHistory[playerID] = valid

	// 取最近的 maxCount 筆（從尾端取）
	if len(valid) > maxCount {
		return valid[len(valid)-maxCount:]
	}
	return valid
}

// tryLuckyTimeRewindFish 擊破 T205 後觸發時光倒流
func (g *Game) tryLuckyTimeRewindFish(p *player.Player) {
	m := g.LuckyTimeRewind
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
	m.personalCooldowns[p.ID] = now.Add(LuckyTimeRewindPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyTimeRewindGlobalCD)
	m.mu.Unlock()

	// 取得最近擊破記錄
	recentKills := m.getRecentKills(p.ID, LuckyTimeRewindMaxReplay)

	log.Printf("[TimeRewind] player=%s 時光倒流啟動，重播 %d 個目標", p.ID, len(recentKills))

	// 計算 betCost
	betDef := data.GetBetDef(p.BetLevel)
	avgBet := betDef.BetCost
	if avgBet < 1 {
		avgBet = 1
	}

	// 場上所有目標 HP 恢復到 60%
	restoredCount := g.doTimeRewindHPRestore()

	// 個人訊息：時光倒流啟動
	replayTargetIDs := make([]string, 0, len(recentKills))
	replayNames := make([]string, 0, len(recentKills))
	for _, r := range recentKills {
		replayTargetIDs = append(replayTargetIDs, r.instanceID)
		replayNames = append(replayNames, r.name)
	}

	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyTimeRewind,
		Payload: ws.LuckyTimeRewindPayload{
			Event:          "rewind_start",
			PlayerID:       p.ID,
			PlayerName:     p.DisplayName,
			ReplayCount:    len(recentKills),
			ReplayTargetIDs: replayTargetIDs,
			ReplayNames:    replayNames,
			ReplayMult:     LuckyTimeRewindReplayMult,
			RestoredCount:  restoredCount,
			HPRestorePct:   int(LuckyTimeRewindHPRestore * 100),
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTimeRewind,
		Payload: ws.LuckyTimeRewindPayload{
			Event:      "rewind_broadcast",
			PlayerName: p.DisplayName,
			ReplayCount: len(recentKills),
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyTimeRewind, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⏪ %s 觸發時光倒流！重播 %d 個目標，×%.1f 倍率！", p.DisplayName, len(recentKills), LuckyTimeRewindReplayMult),
		"color":   "#9B59B6",
	})

	// 啟動重播 goroutine
	go g.runTimeRewindReplay(p, recentKills, avgBet)
}

// doTimeRewindHPRestore 場上所有目標 HP 恢復到 60%
func (g *Game) doTimeRewindHPRestore() int {
	g.mu.Lock()
	defer g.mu.Unlock()

	restoredCount := 0
	for _, t := range g.Targets {
		if !t.IsAlive {
			continue
		}
		maxHP := t.MaxHP
		if maxHP <= 0 {
			continue
		}
		targetHP := int(float64(maxHP) * LuckyTimeRewindHPRestore)
		if targetHP < 1 {
			targetHP = 1
		}
		if t.HP < targetHP {
			t.HP = targetHP
			restoredCount++
		}
	}
	return restoredCount
}

// runTimeRewindReplay 時光倒流重播 goroutine
func (g *Game) runTimeRewindReplay(p *player.Player, kills []killRecord, avgBet int) {
	totalReward := 0

	for i, r := range kills {
		// 間隔 400ms
		if i > 0 {
			select {
			case <-time.After(LuckyTimeRewindReplayDelay):
			case <-g.stopCh:
				return
			}
		}

		// 計算獎勵
		reward := int(float64(avgBet) * float64(r.mult) * LuckyTimeRewindReplayMult / 10.0)
		if reward < int(float64(avgBet)*LuckyTimeRewindReplayMult) {
			reward = int(float64(avgBet) * LuckyTimeRewindReplayMult)
		}
		p.AddCoins(reward)
		totalReward += reward

		log.Printf("[TimeRewind] player=%s 重播 [%d/%d] defID=%s reward=%d",
			p.ID, i+1, len(kills), r.defID, reward)

		// 個人通知：每個重播目標
		_ = g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyTimeRewind,
			Payload: ws.LuckyTimeRewindPayload{
				Event:      "rewind_replay",
				TargetID:   r.instanceID,
				TargetName: r.name,
				ReplayMult: LuckyTimeRewindReplayMult,
				Reward:     reward,
				ReplayIdx:  i + 1,
				TotalCount: len(kills),
			},
		})
	}

	// 結束通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyTimeRewind,
		Payload: ws.LuckyTimeRewindPayload{
			Event:       "rewind_end",
			ReplayCount: len(kills),
			TotalReward: totalReward,
		},
	})

	if len(kills) >= 3 {
		g.Announce.Create(announce.EventLuckyTimeRewind, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("⏪ %s 時光倒流結束！重播 %d 個目標，獲得 %d 籌碼！",
				p.DisplayName, len(kills), totalReward),
			"color": "#8E44AD",
		})
	}
}
