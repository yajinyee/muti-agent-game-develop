// Package dailybonus 每日登入獎勵系統（DAY-065）
// 業界標準：捕魚機遊戲的玩家留存率關鍵功能
// 設計：連續登入天數越多，獎勵越豐厚（最高 7 天循環）
package dailybonus

import (
	"time"
)

// 獎勵表（連續天數 → 金幣獎勵）
// 7 天為一個循環，第 7 天後重置到第 1 天（但保留 streak 計數）
var rewardTable = []int{
	500,  // Day 1
	800,  // Day 2
	1200, // Day 3
	1800, // Day 4
	2500, // Day 5
	3500, // Day 6
	5000, // Day 7（大獎）
}

// UTC+8 時區
var tz = time.FixedZone("UTC+8", 8*60*60)

// TodayDate 取得今日日期字串（UTC+8，格式 "2006-01-02"）
func TodayDate() string {
	return time.Now().In(tz).Format("2006-01-02")
}

// YesterdayDate 取得昨日日期字串（UTC+8）
func YesterdayDate() string {
	return time.Now().In(tz).AddDate(0, 0, -1).Format("2006-01-02")
}

// CheckAndCalc 檢查是否需要發放每日獎勵，計算獎勵金額
// 回傳：(reward, newStreak, isNewLogin)
//   - reward: 獎勵金幣數（0 = 今天已領過）
//   - newStreak: 更新後的連續天數
//   - isNewLogin: 是否是今天第一次登入
func CheckAndCalc(lastLoginDate string, currentStreak int) (reward int, newStreak int, isNewLogin bool) {
	today := TodayDate()
	yesterday := YesterdayDate()

	// 今天已領過
	if lastLoginDate == today {
		return 0, currentStreak, false
	}

	// 連續登入（昨天有登入）
	if lastLoginDate == yesterday {
		newStreak = currentStreak + 1
	} else {
		// 中斷連續（超過一天沒登入，或首次登入）
		newStreak = 1
	}

	// 計算獎勵（7 天循環）
	dayIdx := (newStreak - 1) % len(rewardTable)
	reward = rewardTable[dayIdx]

	return reward, newStreak, true
}

// GetRewardForDay 取得指定天數的獎勵金額（用於 Client 顯示獎勵預覽）
func GetRewardForDay(streak int) int {
	if streak <= 0 {
		streak = 1
	}
	dayIdx := (streak - 1) % len(rewardTable)
	return rewardTable[dayIdx]
}

// GetRewardTable 取得完整獎勵表（7 天）
func GetRewardTable() []int {
	result := make([]int, len(rewardTable))
	copy(result, rewardTable)
	return result
}
