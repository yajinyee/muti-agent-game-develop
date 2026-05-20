package dailybonus

import (
	"testing"
)

func TestCheckAndCalc_FirstLogin(t *testing.T) {
	// 首次登入（lastLoginDate 為空）
	reward, streak, isNew := CheckAndCalc("", 0)
	if !isNew {
		t.Error("首次登入應該是 isNewLogin=true")
	}
	if streak != 1 {
		t.Errorf("首次登入 streak 應為 1，got %d", streak)
	}
	if reward != rewardTable[0] {
		t.Errorf("首次登入獎勵應為 %d，got %d", rewardTable[0], reward)
	}
}

func TestCheckAndCalc_AlreadyLoggedToday(t *testing.T) {
	// 今天已登入過
	today := TodayDate()
	reward, streak, isNew := CheckAndCalc(today, 3)
	if isNew {
		t.Error("今天已登入過，isNewLogin 應為 false")
	}
	if reward != 0 {
		t.Errorf("今天已登入過，reward 應為 0，got %d", reward)
	}
	if streak != 3 {
		t.Errorf("streak 不應改變，應為 3，got %d", streak)
	}
}

func TestCheckAndCalc_ConsecutiveDay(t *testing.T) {
	// 連續登入（昨天有登入）
	yesterday := YesterdayDate()
	reward, streak, isNew := CheckAndCalc(yesterday, 2)
	if !isNew {
		t.Error("連續登入應該是 isNewLogin=true")
	}
	if streak != 3 {
		t.Errorf("連續登入 streak 應為 3，got %d", streak)
	}
	if reward != rewardTable[2] {
		t.Errorf("Day 3 獎勵應為 %d，got %d", rewardTable[2], reward)
	}
}

func TestCheckAndCalc_BrokenStreak(t *testing.T) {
	// 中斷連續（超過一天沒登入）
	reward, streak, isNew := CheckAndCalc("2020-01-01", 10)
	if !isNew {
		t.Error("中斷後重新登入應該是 isNewLogin=true")
	}
	if streak != 1 {
		t.Errorf("中斷後 streak 應重置為 1，got %d", streak)
	}
	if reward != rewardTable[0] {
		t.Errorf("重置後獎勵應為 %d，got %d", rewardTable[0], reward)
	}
}

func TestCheckAndCalc_SevenDayCycle(t *testing.T) {
	// 第 7 天獎勵最高
	yesterday := YesterdayDate()
	reward, streak, _ := CheckAndCalc(yesterday, 6)
	if streak != 7 {
		t.Errorf("streak 應為 7，got %d", streak)
	}
	if reward != rewardTable[6] {
		t.Errorf("Day 7 獎勵應為 %d，got %d", rewardTable[6], reward)
	}
}

func TestCheckAndCalc_CycleReset(t *testing.T) {
	// 第 8 天應循環回 Day 1 的獎勵
	yesterday := YesterdayDate()
	reward, streak, _ := CheckAndCalc(yesterday, 7)
	if streak != 8 {
		t.Errorf("streak 應為 8，got %d", streak)
	}
	// Day 8 = 循環到 Day 1 的獎勵
	if reward != rewardTable[0] {
		t.Errorf("Day 8 應循環到 Day 1 獎勵 %d，got %d", rewardTable[0], reward)
	}
}

func TestGetRewardTable(t *testing.T) {
	table := GetRewardTable()
	if len(table) != 7 {
		t.Errorf("獎勵表應有 7 天，got %d", len(table))
	}
	// 確認獎勵遞增
	for i := 1; i < len(table); i++ {
		if table[i] <= table[i-1] {
			t.Errorf("Day %d 獎勵 %d 應大於 Day %d 獎勵 %d", i+1, table[i], i, table[i-1])
		}
	}
}

func TestGetRewardForDay(t *testing.T) {
	// Day 1
	if r := GetRewardForDay(1); r != rewardTable[0] {
		t.Errorf("Day 1 獎勵應為 %d，got %d", rewardTable[0], r)
	}
	// Day 7
	if r := GetRewardForDay(7); r != rewardTable[6] {
		t.Errorf("Day 7 獎勵應為 %d，got %d", rewardTable[6], r)
	}
	// Day 0（邊界）
	if r := GetRewardForDay(0); r != rewardTable[0] {
		t.Errorf("Day 0 應回傳 Day 1 獎勵 %d，got %d", rewardTable[0], r)
	}
}
