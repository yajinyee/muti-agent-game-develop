// chain_test.go — 連鎖爆炸系統單元測試（DAY-088）
package chain

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := NewDefault()
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestCalcChance(t *testing.T) {
	tests := []struct {
		mult     float64
		expected float64
	}{
		{2.0, 0.05},
		{5.0, 0.10},
		{10.0, 0.15},
		{20.0, 0.20},
		{50.0, 0.30},
		{100.0, 0.30},
	}
	for _, tc := range tests {
		got := CalcChance(tc.mult)
		if got != tc.expected {
			t.Errorf("CalcChance(%.0f) = %.2f, want %.2f", tc.mult, got, tc.expected)
		}
	}
}

func TestTryChain_NoNearby(t *testing.T) {
	m := NewDefault()
	trigger := TargetInfo{ID: "t1", X: 100, Y: 100, Multiplier: 50, DefID: "T001"}
	// 沒有其他目標
	result := m.TryChain(trigger, []TargetInfo{trigger}, 0)
	// 即使觸發機率高，沒有周圍目標就不會連鎖
	if result.Level != ChainNone {
		// 可能觸發但沒有目標，應該回傳 None
		if len(result.TargetIDs) > 0 {
			t.Error("expected no chain targets when no nearby targets")
		}
	}
}

func TestTryChain_MaxDepth(t *testing.T) {
	m := NewDefault()
	trigger := TargetInfo{ID: "t1", X: 100, Y: 100, Multiplier: 100, DefID: "T001"}
	nearby := TargetInfo{ID: "t2", X: 150, Y: 100, Multiplier: 5, DefID: "T002"}
	// 超過最大深度，不應觸發
	result := m.TryChain(trigger, []TargetInfo{trigger, nearby}, m.cfg.MaxDepth)
	if result.Level != ChainNone {
		t.Error("expected no chain at max depth")
	}
}

func TestTryChain_BossExcluded(t *testing.T) {
	m := NewDefault()
	trigger := TargetInfo{ID: "t1", X: 100, Y: 100, Multiplier: 100, DefID: "T001"}
	boss := TargetInfo{ID: "boss1", X: 110, Y: 100, Multiplier: 500, DefID: "B001"}
	// BOSS 不應被連鎖
	// 多次嘗試確認 BOSS 不在結果中
	for i := 0; i < 100; i++ {
		result := m.TryChain(trigger, []TargetInfo{trigger, boss}, 0)
		for _, id := range result.TargetIDs {
			if id == "boss1" {
				t.Error("BOSS should not be chain target")
			}
		}
	}
}

func TestTryChain_RadiusFilter(t *testing.T) {
	m := NewDefault()
	trigger := TargetInfo{ID: "t1", X: 100, Y: 100, Multiplier: 100, DefID: "T001"}
	farTarget := TargetInfo{ID: "t2", X: 500, Y: 100, Multiplier: 5, DefID: "T002"} // 400px 外
	// 超出範圍的目標不應被連鎖
	for i := 0; i < 100; i++ {
		result := m.TryChain(trigger, []TargetInfo{trigger, farTarget}, 0)
		for _, id := range result.TargetIDs {
			if id == "t2" {
				t.Error("far target should not be chain target")
			}
		}
	}
}

func TestCalcLevel(t *testing.T) {
	tests := []struct {
		count    int
		expected ChainLevel
	}{
		{1, ChainSmall},
		{2, ChainMedium},
		{3, ChainMedium},
		{4, ChainBig},
		{6, ChainBig},
		{7, ChainMega},
		{10, ChainMega},
	}
	for _, tc := range tests {
		got := calcLevel(tc.count)
		if got != tc.expected {
			t.Errorf("calcLevel(%d) = %d, want %d", tc.count, got, tc.expected)
		}
	}
}

func TestCalcBonusMult(t *testing.T) {
	m := NewDefault()
	tests := []struct {
		level    ChainLevel
		expected float64
	}{
		{ChainNone, 1.0},
		{ChainSmall, 1.0},
		{ChainMedium, 1.2},
		{ChainBig, 1.5},
		{ChainMega, 2.0},
	}
	for _, tc := range tests {
		got := m.calcBonusMult(tc.level)
		if got != tc.expected {
			t.Errorf("calcBonusMult(%d) = %.1f, want %.1f", tc.level, got, tc.expected)
		}
	}
}

func TestLevelInfo(t *testing.T) {
	levels := []ChainLevel{ChainSmall, ChainMedium, ChainBig, ChainMega}
	for _, level := range levels {
		name, color := levelInfo(level)
		if name == "" {
			t.Errorf("levelInfo(%d) returned empty name", level)
		}
		if color == "" {
			t.Errorf("levelInfo(%d) returned empty color", level)
		}
	}
}

func TestTryChain_WithNearby(t *testing.T) {
	m := NewDefault()
	trigger := TargetInfo{ID: "t1", X: 100, Y: 100, Multiplier: 100, DefID: "T001"}
	targets := []TargetInfo{
		trigger,
		{ID: "t2", X: 120, Y: 100, Multiplier: 5, DefID: "T002"},
		{ID: "t3", X: 140, Y: 100, Multiplier: 3, DefID: "T003"},
		{ID: "t4", X: 160, Y: 100, Multiplier: 2, DefID: "T004"},
	}
	// 100x 觸發機率 30%，多次嘗試應該至少觸發一次
	triggered := false
	for i := 0; i < 200; i++ {
		result := m.TryChain(trigger, targets, 0)
		if result.Level != ChainNone {
			triggered = true
			// 確認連鎖目標都在範圍內
			for _, id := range result.TargetIDs {
				if id == "t1" {
					t.Error("trigger target should not be in chain targets")
				}
			}
			break
		}
	}
	if !triggered {
		t.Error("expected chain to trigger at least once in 200 attempts with 30% chance")
	}
}

func TestTryChain_SelfExcluded(t *testing.T) {
	m := NewDefault()
	trigger := TargetInfo{ID: "t1", X: 100, Y: 100, Multiplier: 100, DefID: "T001"}
	nearby := TargetInfo{ID: "t2", X: 110, Y: 100, Multiplier: 5, DefID: "T002"}
	for i := 0; i < 100; i++ {
		result := m.TryChain(trigger, []TargetInfo{trigger, nearby}, 0)
		for _, id := range result.TargetIDs {
			if id == "t1" {
				t.Error("trigger target should not be in chain targets")
			}
		}
	}
}
