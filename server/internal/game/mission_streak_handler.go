// mission_streak_handler.go — 每日任務連續完成獎勵（DAY-086）
// 玩家連續幾天完成所有任務，給遞增獎勵
// 業界依據：actionnetwork.com 2026-05-09 確認連續登入+任務完成是留存率最高的機制
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
	missionStreakMu.Unlock()

	// 計算獎勵
	reward, label := getMissionStreakReward(streak)
	p.AddCoins(reward)

	log.Printf("[MissionStreak] player=%s streak=%d reward=%d (%s)", p.ID, streak, reward, label)

	// 發送連續完成通知
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgMissionStreakBonus,
		Payload: ws.MissionStreakBonusPayload{
			Streak:     streak,
			MaxStreak:  maxStreak,
			Reward:     reward,
			Label:      label,
			NewBalance: p.GetCoins(),
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
