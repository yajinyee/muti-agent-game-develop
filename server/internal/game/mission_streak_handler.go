// mission_streak_handler.go — 每日任務連續完成獎勵（DAY-086）
// 玩家連續幾天完成所有任務，給遞增獎勵
// DAY-120：加入「寬限期」機制（Streaks with Mercy）
// 業界依據：nowg.net（2026-05-21）確認「Streaks with Mercy」是 2026 年最有效的留存機制
// actionnetwork.com 2026-05-09 確認連續登入+任務完成是留存率最高的機制
package game

import (
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/game/challenge"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// MissionStreakRecord 玩家任務連續完成記錄
type MissionStreakRecord struct {
	Streak        int       // 當前連續天數
	MaxStreak     int       // 歷史最高連續天數
	LastCompleted time.Time // 上次全部完成的日期
	TodayClaimed  bool      // 今天的連續獎勵是否已領取
	// 寬限期（DAY-120）
	MercyUsedAt   time.Time // 上次使用寬限期的時間（7天冷卻）
	MercyCount    int       // 本週已使用寬限次數（每7天重置）
}

// MissionStreakBonus 連續完成獎勵定義
var MissionStreakBonus = []struct {
	Days   int
	Reward int
	Label  string
}{
	{1, 500, "初次完成"},
	{2, 1000, "連續 2 天"},
	{3, 2000, "連續 3 天"},
	{5, 5000, "連續 5 天"},
	{7, 10000, "連續 7 天 🏆"},
	{14, 25000, "連續 14 天 👑"},
	{30, 100000, "連續 30 天 💎"},
}

// missionStreakMu 保護 missionStreaks map
var missionStreakMu sync.RWMutex
var missionStreaks = make(map[string]*MissionStreakRecord) // playerID → record

// getMissionStreakReward 計算連續天數對應的獎勵
func getMissionStreakReward(streak int) (int, string) {
	reward := 500
	label := "完成所有任務"
	for _, b := range MissionStreakBonus {
		if streak >= b.Days {
			reward = b.Reward
			label = b.Label
		}
	}
	return reward, label
}

// canUseMercy 判斷玩家是否可以使用寬限期（每 7 天最多 1 次）
func canUseMercy(rec *MissionStreakRecord) bool {
	if rec.MercyUsedAt.IsZero() {
		return true // 從未使用過
	}
	return time.Since(rec.MercyUsedAt) >= 7*24*time.Hour
}

// checkMissionStreakMercy 在玩家登入時檢查是否需要使用寬限期保護連續記錄（DAY-120）
// 如果玩家昨天沒有完成任務（中斷了），且有寬限期可用，自動保護連續記錄
func (g *Game) checkMissionStreakMercy(p *player.Player) {
	missionStreakMu.Lock()
	rec, ok := missionStreaks[p.ID]
	if !ok || rec.Streak == 0 {
		missionStreakMu.Unlock()
		return // 沒有連續記錄，不需要保護
	}

	// 計算時間差
	now := time.Now()
	loc := time.FixedZone("UTC+8", 8*60*60)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	if rec.LastCompleted.IsZero() {
		missionStreakMu.Unlock()
		return
	}

	lastDay := time.Date(rec.LastCompleted.Year(), rec.LastCompleted.Month(), rec.LastCompleted.Day(), 0, 0, 0, 0, loc)
	diff := today.Sub(lastDay)

	// 只有在「昨天沒完成（中斷 1 天）」且「今天還沒領取」時才觸發寬限期
	if diff <= 48*time.Hour && diff > 24*time.Hour && !rec.TodayClaimed {
		// 中斷了 1 天，檢查是否可以使用寬限期
		if canUseMercy(rec) && rec.Streak >= 3 { // 至少連續 3 天才值得保護
			rec.MercyUsedAt = now
			rec.MercyCount++
			streak := rec.Streak
			mercyLeft := 0 // 使用後本週剩餘 0 次
			missionStreakMu.Unlock()

			log.Printf("[MissionStreak] player=%s mercy used, streak=%d protected", p.ID, streak)

			// 通知玩家連續記錄被保護
			g.Hub.Send(p.ID, &ws.Message{
				Type: ws.MsgMissionMercyProtected,
				Payload: ws.MissionMercyProtectedPayload{
					Streak:    streak,
					MercyLeft: mercyLeft,
					Message:   "🛡️ 你的連續任務記錄被保護了！",
				},
			})
			return
		}
	}
	missionStreakMu.Unlock()
}

// notifyMissionAllComplete 全部任務完成後的連續獎勵處理（DAY-086）
func (g *Game) notifyMissionAllComplete(p *player.Player) {
	missionStreakMu.Lock()
	rec, ok := missionStreaks[p.ID]
	if !ok {
		rec = &MissionStreakRecord{}
		missionStreaks[p.ID] = rec
	}

	// 檢查今天是否已領取
	if rec.TodayClaimed {
		missionStreakMu.Unlock()
		return
	}

	// 計算連續天數
	now := time.Now()
	loc := time.FixedZone("UTC+8", 8*60*60)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	mercyUsed := false
	if rec.LastCompleted.IsZero() {
		// 第一次完成
		rec.Streak = 1
	} else {
		lastDay := time.Date(rec.LastCompleted.Year(), rec.LastCompleted.Month(), rec.LastCompleted.Day(), 0, 0, 0, 0, loc)
		diff := today.Sub(lastDay)
		if diff <= 24*time.Hour && diff >= 0 {
			// 同一天（不應該發生，但防禦性處理）
			missionStreakMu.Unlock()
			return
		} else if diff <= 48*time.Hour {
			// 昨天完成過，連續+1
			rec.Streak++
		} else if diff <= 72*time.Hour && canUseMercy(rec) && rec.Streak >= 3 {
			// 中斷了 2 天以內，且有寬限期可用（DAY-120）
			// 寬限期：連續記錄不重置，但不增加
			rec.MercyUsedAt = now
			rec.MercyCount++
			mercyUsed = true
			// Streak 保持不變（不增加，不重置）
		} else {
			// 中斷了，重置
			rec.Streak = 1
		}
	}

	rec.LastCompleted = now
	rec.TodayClaimed = true
	if rec.Streak > rec.MaxStreak {
		rec.MaxStreak = rec.Streak
	}

	streak := rec.Streak
	maxStreak := rec.MaxStreak
	mercyLeft := 0
	if canUseMercy(rec) {
		mercyLeft = 1
	}
	missionStreakMu.Unlock()

	// 計算獎勵（使用寬限期時獎勵減半）
	reward, label := getMissionStreakReward(streak)
	if mercyUsed {
		reward = reward / 2
		label = "🛡️ 寬限期保護（" + label + "）"
	}
	p.AddCoins(reward)

	log.Printf("[MissionStreak] player=%s streak=%d reward=%d (%s) mercy=%v",
		p.ID, streak, reward, label, mercyUsed)

	// 發送連續完成通知
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgMissionStreakBonus,
		Payload: ws.MissionStreakBonusPayload{
			Streak:     streak,
			MaxStreak:  maxStreak,
			Reward:     reward,
			Label:      label,
			NewBalance: p.GetCoins(),
			MercyUsed:  mercyUsed,
			MercyLeft:  mercyLeft,
		},
	}); err != nil {
		log.Printf("[MissionStreak] send error: %v", err)
	}

	// 連續 7 天解鎖隱藏挑戰（DAY-085 整合）
	if streak >= 7 {
		if def := g.Challenge.TryUnlock(p.ID, challenge.ChallengeMissionStreak7); def != nil {
			g.sendChallengeUnlocked(p, def)
		}
	}
	// 連續 30 天解鎖隱藏挑戰
	if streak >= 30 {
		if def := g.Challenge.TryUnlock(p.ID, challenge.ChallengeMissionStreak30); def != nil {
			g.sendChallengeUnlocked(p, def)
		}
	}
}
