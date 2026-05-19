package jackpot

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	snap := m.GetSnapshot()

	if snap[LevelMini] != 100 {
		t.Errorf("Mini initial = %d, want 100", snap[LevelMini])
	}
	if snap[LevelMajor] != 500 {
		t.Errorf("Major initial = %d, want 500", snap[LevelMajor])
	}
	if snap[LevelGrand] != 2000 {
		t.Errorf("Grand initial = %d, want 2000", snap[LevelGrand])
	}
}

func TestContribute_IncreasesPool(t *testing.T) {
	m := NewManager()
	initialSnap := m.GetSnapshot()

	// 貢獻 1000 betCost，0.5% = 5，Mini 拿 60% = 3
	m.Contribute(1000, "player1")
	snap := m.GetSnapshot()

	if snap[LevelMini] <= initialSnap[LevelMini] {
		t.Errorf("Mini pool should increase after contribution, got %d (was %d)", snap[LevelMini], initialSnap[LevelMini])
	}
	if snap[LevelMajor] <= initialSnap[LevelMajor] {
		t.Errorf("Major pool should increase after contribution")
	}
	if snap[LevelGrand] <= initialSnap[LevelGrand] {
		t.Errorf("Grand pool should increase after contribution")
	}
}

func TestContribute_SmallBet(t *testing.T) {
	m := NewManager()
	initialSnap := m.GetSnapshot()

	// 小額 betCost = 10，0.5% = 0.05，最少貢獻 1
	m.Contribute(10, "player1")
	snap := m.GetSnapshot()

	// 至少有一個池子增加了
	increased := snap[LevelMini] > initialSnap[LevelMini] ||
		snap[LevelMajor] > initialSnap[LevelMajor] ||
		snap[LevelGrand] > initialSnap[LevelGrand]
	if !increased {
		t.Error("At least one pool should increase even with small bet")
	}
}

func TestForceWin_Mini(t *testing.T) {
	m := NewManager()
	initialAmount := m.GetSnapshot()[LevelMini]

	win := m.ForceWin(LevelMini, "player1")
	if win == nil {
		t.Fatal("ForceWin should return a win")
	}
	if win.Level != LevelMini {
		t.Errorf("Win level = %s, want mini", win.Level)
	}
	if win.Amount != initialAmount {
		t.Errorf("Win amount = %d, want %d", win.Amount, initialAmount)
	}
	if win.WinnerID != "player1" {
		t.Errorf("Winner = %s, want player1", win.WinnerID)
	}

	// 池子應該重置到基礎金額
	snap := m.GetSnapshot()
	if snap[LevelMini] != 100 {
		t.Errorf("Mini pool after win = %d, want 100 (base)", snap[LevelMini])
	}
}

func TestForceWin_Grand(t *testing.T) {
	m := NewManager()

	win := m.ForceWin(LevelGrand, "player2")
	if win == nil {
		t.Fatal("ForceWin Grand should return a win")
	}
	if win.Level != LevelGrand {
		t.Errorf("Win level = %s, want grand", win.Level)
	}

	// 池子重置
	snap := m.GetSnapshot()
	if snap[LevelGrand] != 2000 {
		t.Errorf("Grand pool after win = %d, want 2000 (base)", snap[LevelGrand])
	}
}

func TestForceWin_InvalidLevel(t *testing.T) {
	m := NewManager()
	win := m.ForceWin("invalid", "player1")
	if win != nil {
		t.Error("ForceWin with invalid level should return nil")
	}
}

func TestContribute_NoReturnNilNormally(t *testing.T) {
	m := NewManager()

	// 正常貢獻不應該觸發（機率極低）
	// 執行 100 次，期望大多數不觸發
	wins := 0
	for i := 0; i < 100; i++ {
		win := m.Contribute(100, "player1")
		if win != nil {
			wins++
		}
	}
	// 100 次中觸發超過 10 次是異常（機率 < 0.001%）
	if wins > 10 {
		t.Errorf("Too many jackpot wins in 100 contributions: %d", wins)
	}
}

// TestJackpotFrequency 驗證 Jackpot 觸發頻率合理性
// LV5 射擊速度 3 shots/sec，betCost=50
// Mini 應該平均每 5-15 分鐘觸發一次（900-2700 shots）
func TestJackpotFrequency(t *testing.T) {
	const shots = 100000
	const betCost = 50

	m := NewManager()
	miniWins := 0
	majorWins := 0
	grandWins := 0

	for i := 0; i < shots; i++ {
		win := m.Contribute(betCost, "player1")
		if win != nil {
			switch win.Level {
			case LevelMini:
				miniWins++
			case LevelMajor:
				majorWins++
			case LevelGrand:
				grandWins++
			}
		}
	}

	// 100,000 shots @ 3 shots/sec = 33,333 秒 ≈ 9.3 小時
	// Mini 期望觸發次數：合理範圍 5-100 次（平均每 1000-20000 shots 一次）
	t.Logf("Mini wins: %d (avg every %.0f shots)", miniWins, float64(shots)/float64(max(miniWins, 1)))
	t.Logf("Major wins: %d (avg every %.0f shots)", majorWins, float64(shots)/float64(max(majorWins, 1)))
	t.Logf("Grand wins: %d (avg every %.0f shots)", grandWins, float64(shots)/float64(max(grandWins, 1)))

	// Mini 不應該太頻繁（< 200 次 in 100k shots = 每 500 shots 一次）
	if miniWins > 200 {
		t.Errorf("Mini jackpot too frequent: %d wins in %d shots (avg every %.0f shots)",
			miniWins, shots, float64(shots)/float64(miniWins))
	}
	// Mini 不應該太稀少（> 0 次 in 100k shots）
	if miniWins == 0 {
		t.Error("Mini jackpot never triggered in 100k shots - too rare")
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
