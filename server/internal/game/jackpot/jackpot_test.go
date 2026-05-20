package jackpot

import (
	"testing"
)

func TestNewManager_FourLevels(t *testing.T) {
	m := NewManager()
	snap := m.GetSnapshot()

	if snap[LevelMini] != 80 {
		t.Errorf("Mini initial = %d, want 80", snap[LevelMini])
	}
	if snap[LevelMinor] != 200 {
		t.Errorf("Minor initial = %d, want 200", snap[LevelMinor])
	}
	if snap[LevelMajor] != 500 {
		t.Errorf("Major initial = %d, want 500", snap[LevelMajor])
	}
	if snap[LevelGrand] != 2000 {
		t.Errorf("Grand initial = %d, want 2000", snap[LevelGrand])
	}
}

func TestContribute_IncreasesAllFourPools(t *testing.T) {
	m := NewManager()
	initialSnap := m.GetSnapshot()

	// 貢獻 1000 betCost，0.5% = 5，四層都應增加
	m.Contribute(1000, "player1")
	snap := m.GetSnapshot()

	if snap[LevelMini] <= initialSnap[LevelMini] {
		t.Errorf("Mini pool should increase after contribution, got %d (was %d)", snap[LevelMini], initialSnap[LevelMini])
	}
	if snap[LevelMinor] <= initialSnap[LevelMinor] {
		t.Errorf("Minor pool should increase after contribution")
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

	// 小額 betCost = 10，最少貢獻 4（每層至少 1）
	m.Contribute(10, "player1")
	snap := m.GetSnapshot()

	// 所有池子都應增加
	if snap[LevelMini] <= initialSnap[LevelMini] {
		t.Error("Mini pool should increase even with small bet")
	}
	if snap[LevelMinor] <= initialSnap[LevelMinor] {
		t.Error("Minor pool should increase even with small bet")
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
	if snap[LevelMini] != 80 {
		t.Errorf("Mini pool after win = %d, want 80 (base)", snap[LevelMini])
	}
}

func TestForceWin_Minor(t *testing.T) {
	m := NewManager()
	initialAmount := m.GetSnapshot()[LevelMinor]

	win := m.ForceWin(LevelMinor, "player2")
	if win == nil {
		t.Fatal("ForceWin Minor should return a win")
	}
	if win.Level != LevelMinor {
		t.Errorf("Win level = %s, want minor", win.Level)
	}
	if win.Amount != initialAmount {
		t.Errorf("Win amount = %d, want %d", win.Amount, initialAmount)
	}

	// 池子重置
	snap := m.GetSnapshot()
	if snap[LevelMinor] != 200 {
		t.Errorf("Minor pool after win = %d, want 200 (base)", snap[LevelMinor])
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
	wins := 0
	for i := 0; i < 100; i++ {
		win := m.Contribute(100, "player1")
		if win != nil {
			wins++
		}
	}
	// 100 次中觸發超過 10 次是異常
	if wins > 10 {
		t.Errorf("Too many jackpot wins in 100 contributions: %d", wins)
	}
}

func TestSaveLoadState_FourLevels(t *testing.T) {
	m := NewManager()

	// 貢獻一些金額
	for i := 0; i < 100; i++ {
		m.Contribute(500, "player1")
	}

	state := m.SaveState()
	if state.Minor == 0 {
		t.Error("Minor state should be saved")
	}

	// 建立新 manager 並載入狀態
	m2 := NewManager()
	m2.LoadState(state)
	snap := m2.GetSnapshot()

	if snap[LevelMinor] != state.Minor {
		t.Errorf("Minor after load = %d, want %d", snap[LevelMinor], state.Minor)
	}
	if snap[LevelMajor] != state.Major {
		t.Errorf("Major after load = %d, want %d", snap[LevelMajor], state.Major)
	}
}

func TestGetLevelInfo(t *testing.T) {
	name, color, icon := GetLevelInfo(LevelMini)
	if name != "MINI" {
		t.Errorf("Mini name = %s, want MINI", name)
	}
	if color == "" || icon == "" {
		t.Error("Mini color/icon should not be empty")
	}

	name, _, _ = GetLevelInfo(LevelMinor)
	if name != "MINOR" {
		t.Errorf("Minor name = %s, want MINOR", name)
	}

	name, _, _ = GetLevelInfo(LevelGrand)
	if name != "GRAND" {
		t.Errorf("Grand name = %s, want GRAND", name)
	}
}

// TestJackpotFrequency 驗證四層 Jackpot 觸發頻率合理性
func TestJackpotFrequency(t *testing.T) {
	const shots = 100000
	const betCost = 50

	m := NewManager()
	wins := map[Level]int{}

	for i := 0; i < shots; i++ {
		win := m.Contribute(betCost, "player1")
		if win != nil {
			wins[win.Level]++
		}
	}

	t.Logf("Mini wins: %d (avg every %.0f shots)", wins[LevelMini], float64(shots)/float64(maxInt(wins[LevelMini], 1)))
	t.Logf("Minor wins: %d (avg every %.0f shots)", wins[LevelMinor], float64(shots)/float64(maxInt(wins[LevelMinor], 1)))
	t.Logf("Major wins: %d (avg every %.0f shots)", wins[LevelMajor], float64(shots)/float64(maxInt(wins[LevelMajor], 1)))
	t.Logf("Grand wins: %d (avg every %.0f shots)", wins[LevelGrand], float64(shots)/float64(maxInt(wins[LevelGrand], 1)))

	// Mini 不應該太頻繁（< 300 次 in 100k shots）
	if wins[LevelMini] > 300 {
		t.Errorf("Mini jackpot too frequent: %d wins in %d shots", wins[LevelMini], shots)
	}
	// Mini 不應該太稀少（> 0 次 in 100k shots）
	if wins[LevelMini] == 0 {
		t.Error("Mini jackpot never triggered in 100k shots - too rare")
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
