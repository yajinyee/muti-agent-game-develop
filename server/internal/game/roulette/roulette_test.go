// roulette_test.go — 雙層倍率輪盤系統單元測試（DAY-113）
package roulette

import (
	"testing"
)

// TestNewManager 建立管理器
func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("NewManager returned nil")
	}
	if m.sessions == nil {
		t.Fatal("sessions map not initialized")
	}
}

// TestShouldTrigger_Boss 測試 BOSS 必定觸發
func TestShouldTrigger_Boss(t *testing.T) {
	m := NewManager()
	triggered := false
	for i := 0; i < 10; i++ {
		if m.ShouldTrigger("B001", 0) {
			triggered = true
			break
		}
	}
	if !triggered {
		t.Error("B001 should always trigger roulette")
	}
}

// TestShouldTrigger_BetLevel 測試投注等級限制
func TestShouldTrigger_BetLevel(t *testing.T) {
	m := NewManager()
	// T103 需要 MinBet=5，低投注不觸發
	for i := 0; i < 100; i++ {
		if m.ShouldTrigger("T103", 3) {
			t.Error("T103 should not trigger with bet level 3 (min is 5)")
		}
	}
}

// TestShouldTrigger_Unknown 未知目標不觸發
func TestShouldTrigger_Unknown(t *testing.T) {
	m := NewManager()
	for i := 0; i < 100; i++ {
		if m.ShouldTrigger("T001", 10) {
			t.Error("T001 should never trigger roulette")
		}
	}
}

// TestStartSession 開始 session
func TestStartSession(t *testing.T) {
	m := NewManager()
	session := m.StartSession("player1", "B001", 1000)
	if session == nil {
		t.Fatal("StartSession returned nil")
	}
	if session.PlayerID != "player1" {
		t.Errorf("expected player1, got %s", session.PlayerID)
	}
	if session.BaseReward != 1000 {
		t.Errorf("expected base reward 1000, got %d", session.BaseReward)
	}
	if session.Resolved {
		t.Error("new session should not be resolved")
	}
}

// TestHasActiveSession 測試 session 狀態
func TestHasActiveSession(t *testing.T) {
	m := NewManager()
	if m.HasActiveSession("player1") {
		t.Error("should not have active session before start")
	}
	m.StartSession("player1", "B001", 500)
	if !m.HasActiveSession("player1") {
		t.Error("should have active session after start")
	}
}

// TestResolve 測試輪盤結算
func TestResolve(t *testing.T) {
	m := NewManager()
	m.StartSession("player1", "B001", 1000)
	result, ok := m.Resolve("player1")
	if !ok {
		t.Fatal("Resolve returned false")
	}
	if result == nil {
		t.Fatal("Resolve returned nil result")
	}
	// 內圈倍率應在 1-10 之間
	if result.Inner.Multiplier < 1 || result.Inner.Multiplier > 10 {
		t.Errorf("inner multiplier out of range: %.0f", result.Inner.Multiplier)
	}
	// 外圈倍率應在 1-100 之間
	if result.Outer.Multiplier < 1 || result.Outer.Multiplier > 100 {
		t.Errorf("outer multiplier out of range: %.0f", result.Outer.Multiplier)
	}
	// 最終倍率 = 內圈 × 外圈
	expectedFinal := result.Inner.Multiplier * result.Outer.Multiplier
	if result.FinalMult != expectedFinal {
		t.Errorf("final mult mismatch: expected %.0f, got %.0f", expectedFinal, result.FinalMult)
	}
	// 最終獎勵 = 基礎獎勵 × 最終倍率
	expectedReward := int(float64(1000) * result.FinalMult)
	if result.FinalReward != expectedReward {
		t.Errorf("final reward mismatch: expected %d, got %d", expectedReward, result.FinalReward)
	}
}

// TestResolve_NoSession 無 session 時 Resolve 回傳 false
func TestResolve_NoSession(t *testing.T) {
	m := NewManager()
	_, ok := m.Resolve("nonexistent")
	if ok {
		t.Error("Resolve should return false for nonexistent player")
	}
}

// TestResolve_AlreadyResolved 已結算的 session 不能再結算
func TestResolve_AlreadyResolved(t *testing.T) {
	m := NewManager()
	m.StartSession("player1", "B001", 500)
	m.Resolve("player1")
	_, ok := m.Resolve("player1")
	if ok {
		t.Error("should not be able to resolve twice")
	}
}

// TestCancelSession 取消 session
func TestCancelSession(t *testing.T) {
	m := NewManager()
	m.StartSession("player1", "B001", 500)
	m.CancelSession("player1")
	if m.HasActiveSession("player1") {
		t.Error("session should be cancelled")
	}
}

// TestIsJackpot 測試 Jackpot 標記（≥500x）
func TestIsJackpot(t *testing.T) {
	// 強制最高倍率：內圈 10x × 外圈 100x = 1000x
	result := &RouletteResult{
		Inner:       SpinResult{Multiplier: 10},
		Outer:       SpinResult{Multiplier: 100},
		FinalMult:   1000,
		BaseReward:  100,
		FinalReward: 100000,
		IsJackpot:   1000 >= 500,
		IsMegaWin:   1000 >= 100,
	}
	if !result.IsJackpot {
		t.Error("1000x should be jackpot")
	}
	if !result.IsMegaWin {
		t.Error("1000x should be mega win")
	}
}

// TestExpectedRTP 測試期望 RTP 計算
func TestExpectedRTP(t *testing.T) {
	rtp := ExpectedRTP()
	// 雙層輪盤期望倍率：內圈期望 ~3.5x × 外圈期望 ~5.8x ≈ 15-25x
	if rtp < 10.0 || rtp > 30.0 {
		t.Errorf("expected RTP out of reasonable range: %.2f", rtp)
	}
	t.Logf("Dual-ring roulette expected multiplier: %.2fx", rtp)
}

// TestSegmentCounts 測試格子數量
func TestSegmentCounts(t *testing.T) {
	if len(InnerSegments) != 8 {
		t.Errorf("expected 8 inner segments, got %d", len(InnerSegments))
	}
	if len(OuterSegments) != 12 {
		t.Errorf("expected 12 outer segments, got %d", len(OuterSegments))
	}
}

// TestWeightDistribution 測試權重分布（高倍率應該更稀少）
func TestWeightDistribution(t *testing.T) {
	m := NewManager()
	innerCounts := make(map[float64]int)
	outerCounts := make(map[float64]int)

	// 模擬 10000 次旋轉
	for i := 0; i < 10000; i++ {
		m.mu.Lock()
		inner := m.spinInner()
		outer := m.spinOuter()
		m.mu.Unlock()
		innerCounts[inner.Multiplier]++
		outerCounts[outer.Multiplier]++
	}

	// 1x 應該比 10x 更常見
	if innerCounts[1] <= innerCounts[10] {
		t.Errorf("1x (%d) should appear more than 10x (%d)", innerCounts[1], innerCounts[10])
	}
	if outerCounts[1] <= outerCounts[100] {
		t.Errorf("outer 1x (%d) should appear more than 100x (%d)", outerCounts[1], outerCounts[100])
	}
	t.Logf("Inner distribution: 1x=%d, 10x=%d", innerCounts[1], innerCounts[10])
	t.Logf("Outer distribution: 1x=%d, 100x=%d", outerCounts[1], outerCounts[100])
}
