// Package game — Daily Bonus（每日登入獎勵）handler（DAY-065）
// 業界標準：捕魚機遊戲的玩家留存率關鍵功能
// 設計：連續登入天數越多，獎勵越豐厚（7 天循環，最高 5000 金幣）
package game

import (
	"log"

	"digital-twin/server/internal/game/dailybonus"
	"digital-twin/server/internal/store"
	"digital-twin/server/internal/ws"
)

// checkAndSendDailyBonus 檢查並發放每日登入獎勵
// 在玩家加入後 200ms 非同步呼叫
// savedState 為 nil 時代表新玩家（首次登入）
func (g *Game) checkAndSendDailyBonus(playerID string, savedState *store.PlayerState) {
	g.mu.RLock()
	p, ok := g.Players[playerID]
	g.mu.RUnlock()
	if !ok {
		return
	}

	// 取得上次登入資訊
	lastLoginDate := p.LastLoginDate
	currentStreak := p.LoginStreak

	// 計算每日獎勵
	reward, newStreak, isNewLogin := dailybonus.CheckAndCalc(lastLoginDate, currentStreak)

	if !isNewLogin {
		// 今天已領過，不重複發放
		log.Printf("[DailyBonus] Player %s already claimed today (streak=%d)", playerID, currentStreak)
		return
	}

	// 更新玩家登入資訊
	g.mu.Lock()
	p.LoginStreak = newStreak
	p.LastLoginDate = dailybonus.TodayDate()
	if newStreak > p.MaxLoginStreak {
		p.MaxLoginStreak = newStreak
	}
	maxStreak := p.MaxLoginStreak
	g.mu.Unlock()

	// 發放獎勵
	p.AddReward(reward)
	log.Printf("[DailyBonus] Player %s day %d streak, reward=%d coins (new balance=%d)",
		playerID, newStreak, reward, p.Coins)

	// 通知玩家
	g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgDailyBonus,
		Payload: ws.DailyBonusPayload{
			Streak:      newStreak,
			Reward:      reward,
			NewBalance:  p.Coins,
			IsNewStreak: isNewLogin,
			MaxStreak:   maxStreak,
		},
	})

	// 同步更新玩家狀態到 Client
	g.sendPlayerUpdate(p)
}
