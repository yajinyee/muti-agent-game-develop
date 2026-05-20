package wheel

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestSlots_TotalWeight(t *testing.T) {
	if totalWeight != 100 {
		t.Errorf("expected totalWeight=100, got %d", totalWeight)
	}
}

func TestSlots_Count(t *testing.T) {
	if len(Slots) != 8 {
		t.Errorf("expected 8 slots, got %d", len(Slots))
	}
}

func TestSpin_ValidIndex(t *testing.T) {
	m := NewManager()
	for i := 0; i < 100; i++ {
		idx, slot := m.Spin()
		if idx < 0 || idx >= len(Slots) {
			t.Errorf("invalid slot index: %d", idx)
		}
		if slot.Multiplier <= 0 {
			t.Errorf("invalid multiplier: %f", slot.Multiplier)
		}
	}
}

func TestSpin_Distribution(t *testing.T) {
	m := NewManager()
	counts := make(map[int]int)
	n := 10000
	for i := 0; i < n; i++ {
		idx, _ := m.Spin()
		counts[idx]++
	}
	// 最高倍率（idx 7, weight=1）應該比最低倍率（idx 0, weight=30）少很多
	if counts[7] >= counts[0] {
		t.Errorf("expected rare slot (idx 7) to appear less than common (idx 0): %d vs %d", counts[7], counts[0])
	}
}

func TestShouldTrigger_KnownTarget(t *testing.T) {
	m := NewManager()
	triggered := 0
	n := 1000
	for i := 0; i < n; i++ {
		if m.ShouldTrigger("T103") {
			triggered++
		}
	}
	// T103 機率 15%，1000次應該在 50-250 之間
	if triggered < 50 || triggered > 250 {
		t.Errorf("T103 trigger count out of range: %d/1000", triggered)
	}
}

func TestShouldTrigger_UnknownTarget(t *testing.T) {
	m := NewManager()
	if m.ShouldTrigger("T001") {
		t.Error("T001 should never trigger wheel")
	}
}

func TestShouldTrigger_Boss(t *testing.T) {
	m := NewManager()
	triggered := 0
	n := 1000
	for i := 0; i < n; i++ {
		if m.ShouldTrigger("B001") {
			triggered++
		}
	}
	// B001 機率 50%，1000次應該在 350-650 之間
	if triggered < 350 || triggered > 650 {
		t.Errorf("B001 trigger count out of range: %d/1000", triggered)
	}
}

func TestSlots_MultiplierIncreasing(t *testing.T) {
	for i := 1; i < len(Slots); i++ {
		if Slots[i].Multiplier <= Slots[i-1].Multiplier {
			t.Errorf("Slots not sorted by multiplier at index %d", i)
		}
	}
}
