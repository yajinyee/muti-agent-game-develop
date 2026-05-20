package codex

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	unlocked, total := m.GetStats()
	if total != 12 {
		t.Errorf("expected 12 total entries, got %d", total)
	}
	if unlocked != 0 {
		t.Errorf("expected 0 unlocked initially, got %d", unlocked)
	}
}

func TestRecordKill_FirstKill(t *testing.T) {
	m := NewManager()
	isNew, isComplete := m.RecordKill("T001", 2.0)
	if !isNew {
		t.Error("expected isNewUnlock=true on first kill")
	}
	if isComplete {
		t.Error("expected isComplete=false after first kill")
	}

	unlocked, _ := m.GetStats()
	if unlocked != 1 {
		t.Errorf("expected 1 unlocked, got %d", unlocked)
	}
}

func TestRecordKill_SecondKill_NotNew(t *testing.T) {
	m := NewManager()
	m.RecordKill("T001", 2.0)
	isNew, _ := m.RecordKill("T001", 2.0)
	if isNew {
		t.Error("expected isNewUnlock=false on second kill")
	}
}

func TestRecordKill_MaxMultiplier(t *testing.T) {
	m := NewManager()
	m.RecordKill("T103", 20.0)
	m.RecordKill("T103", 50.0)
	m.RecordKill("T103", 30.0)

	snap := m.GetSnapshot()
	for _, e := range snap {
		if e.TargetID == "T103" {
			if e.MaxMultiplier != 50.0 {
				t.Errorf("expected MaxMultiplier=50, got %f", e.MaxMultiplier)
			}
			if e.KillCount != 3 {
				t.Errorf("expected KillCount=3, got %d", e.KillCount)
			}
		}
	}
}

func TestRecordKill_Complete(t *testing.T) {
	m := NewManager()
	// 解鎖所有條目
	for _, id := range AllTargetIDs[:len(AllTargetIDs)-1] {
		m.RecordKill(id, 1.0)
	}
	// 最後一個
	_, isComplete := m.RecordKill(AllTargetIDs[len(AllTargetIDs)-1], 1.0)
	if !isComplete {
		t.Error("expected isComplete=true after unlocking all entries")
	}
}

func TestIsComplete(t *testing.T) {
	m := NewManager()
	if m.IsComplete() {
		t.Error("expected IsComplete=false initially")
	}
	for _, id := range AllTargetIDs {
		m.RecordKill(id, 1.0)
	}
	if !m.IsComplete() {
		t.Error("expected IsComplete=true after all kills")
	}
}

func TestRecordKill_UnknownTarget(t *testing.T) {
	m := NewManager()
	isNew, isComplete := m.RecordKill("UNKNOWN", 1.0)
	if isNew || isComplete {
		t.Error("expected no unlock for unknown target")
	}
}

func TestGetSnapshot_Order(t *testing.T) {
	m := NewManager()
	snap := m.GetSnapshot()
	if len(snap) != 12 {
		t.Errorf("expected 12 entries in snapshot, got %d", len(snap))
	}
	// 確認順序與 AllTargetIDs 一致
	for i, e := range snap {
		if e.TargetID != AllTargetIDs[i] {
			t.Errorf("expected entry[%d]=%s, got %s", i, AllTargetIDs[i], e.TargetID)
		}
	}
}

func TestLoadState(t *testing.T) {
	m := NewManager()
	saved := []*Entry{
		{TargetID: "T001", Unlocked: true, KillCount: 5, MaxMultiplier: 2.0},
		{TargetID: "B001", Unlocked: true, KillCount: 1, MaxMultiplier: 300.0},
	}
	m.LoadState(saved)

	unlocked, _ := m.GetStats()
	if unlocked != 2 {
		t.Errorf("expected 2 unlocked after LoadState, got %d", unlocked)
	}
}

func TestRarity(t *testing.T) {
	m := NewManager()
	snap := m.GetSnapshot()
	rarityMap := map[string]string{}
	for _, e := range snap {
		rarityMap[e.TargetID] = e.Rarity
	}

	tests := []struct {
		id     string
		rarity string
	}{
		{"T001", "common"},
		{"T101", "rare"},
		{"T103", "epic"},
		{"B001", "legendary"},
	}
	for _, tt := range tests {
		if rarityMap[tt.id] != tt.rarity {
			t.Errorf("expected %s rarity=%s, got %s", tt.id, tt.rarity, rarityMap[tt.id])
		}
	}
}
