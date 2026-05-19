// Package target 測試
package target

import (
	"testing"

	"digital-twin/server/internal/data"
)

// TestGetSpawnWeights 確認三段動態難度權重正確
func TestGetSpawnWeights(t *testing.T) {
	tests := []struct {
		betLevel     int
		wantBasic    float64
		wantSpecial  float64
		wantHigh     float64
	}{
		{1, 0.90, 0.09, 0.01},
		{3, 0.90, 0.09, 0.01},
		{4, 0.82, 0.15, 0.03},
		{7, 0.82, 0.15, 0.03},
		{8, 0.75, 0.20, 0.05},
		{10, 0.75, 0.20, 0.05},
	}

	for _, tc := range tests {
		w := GetSpawnWeights(tc.betLevel)
		if w.BasicRatio != tc.wantBasic {
			t.Errorf("BetLevel %d: BasicRatio = %.2f, want %.2f", tc.betLevel, w.BasicRatio, tc.wantBasic)
		}
		if w.SpecialRatio != tc.wantSpecial {
			t.Errorf("BetLevel %d: SpecialRatio = %.2f, want %.2f", tc.betLevel, w.SpecialRatio, tc.wantSpecial)
		}
		if w.HighRatio != tc.wantHigh {
			t.Errorf("BetLevel %d: HighRatio = %.2f, want %.2f", tc.betLevel, w.HighRatio, tc.wantHigh)
		}
		// 三個比例加總應接近 1.0
		total := w.BasicRatio + w.SpecialRatio + w.HighRatio
		if total < 0.99 || total > 1.01 {
			t.Errorf("BetLevel %d: total ratio = %.2f, want ~1.0", tc.betLevel, total)
		}
	}
}

// TestGetHighValuePool 確認高倍率 pool 只包含 T104 和 T105
func TestGetHighValuePool(t *testing.T) {
	pool := getHighValuePool()
	if len(pool) == 0 {
		t.Fatal("highValuePool is empty")
	}
	for _, d := range pool {
		if d.ID != "T104" && d.ID != "T105" {
			t.Errorf("unexpected target in highValuePool: %s", d.ID)
		}
	}
	// 確認 T104 和 T105 都在 pool 中
	ids := make(map[string]bool)
	for _, d := range pool {
		ids[d.ID] = true
	}
	if !ids["T104"] {
		t.Error("T104 (金色雜草) not in highValuePool")
	}
	if !ids["T105"] {
		t.Error("T105 (巨大金幣魚) not in highValuePool")
	}
}

// TestGetSpecialPool 確認一般特殊 pool 不包含 T104 和 T105
func TestGetSpecialPool(t *testing.T) {
	pool := getSpecialPool()
	for _, d := range pool {
		if d.ID == "T104" || d.ID == "T105" {
			t.Errorf("high-value target %s should not be in specialPool", d.ID)
		}
	}
}

// TestPickTargetDef_HighBetLevel 高投注等級應有機會選到高倍率目標
func TestPickTargetDef_HighBetLevel(t *testing.T) {
	s := NewSpawnSystem()
	highValueCount := 0
	totalRuns := 10000

	for i := 0; i < totalRuns; i++ {
		def := s.PickTargetDef(10, 0) // LV10，無補償
		if def.ID == "T104" || def.ID == "T105" {
			highValueCount++
		}
	}

	// LV10 的 HighRatio = 0.05，期望約 5% 的高倍率目標
	// 允許 ±2% 的統計誤差
	ratio := float64(highValueCount) / float64(totalRuns)
	if ratio < 0.03 || ratio > 0.08 {
		t.Errorf("LV10 high-value ratio = %.3f, expected ~0.05 (±0.02)", ratio)
	}
}

// TestPickTargetDef_LowBetLevel 低投注等級高倍率目標應很少
func TestPickTargetDef_LowBetLevel(t *testing.T) {
	s := NewSpawnSystem()
	highValueCount := 0
	totalRuns := 10000

	for i := 0; i < totalRuns; i++ {
		def := s.PickTargetDef(1, 0) // LV1，無補償
		if def.ID == "T104" || def.ID == "T105" {
			highValueCount++
		}
	}

	// LV1 的 HighRatio = 0.01，期望約 1% 的高倍率目標
	// 允許 ±1% 的統計誤差
	ratio := float64(highValueCount) / float64(totalRuns)
	if ratio > 0.03 {
		t.Errorf("LV1 high-value ratio = %.3f, expected ~0.01 (max 0.03)", ratio)
	}
}

// TestNewTarget_MeteorMultiplier 流星倍率應在 20-50 範圍內
func TestNewTarget_MeteorMultiplier(t *testing.T) {
	def := data.Targets["T103"]
	if def == nil {
		t.Fatal("T103 not found")
	}

	for i := 0; i < 1000; i++ {
		target := NewTarget("test-id", def, 100, 100)
		if target.Multiplier < 20 || target.Multiplier > 50 {
			t.Errorf("T103 multiplier = %.0f, expected 20-50", target.Multiplier)
		}
	}
}

// TestRequiredHits_SpecialTargetNoGuarantee 特殊目標不設保底
func TestRequiredHits_SpecialTargetNoGuarantee(t *testing.T) {
	def := data.Targets["T101"]
	if def == nil {
		t.Fatal("T101 not found")
	}
	target := NewTarget("test-id", def, 100, 100)
	required := target.RequiredHits(10)
	if required != 99999 {
		t.Errorf("Special target required hits = %d, expected 99999 (no guarantee)", required)
	}
}
