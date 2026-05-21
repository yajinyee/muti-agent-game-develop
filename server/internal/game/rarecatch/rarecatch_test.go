package rarecatch

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.GetMultBoost("p1") != 1.0 {
		t.Error("new player should have mult boost 1.0")
	}
}

func TestIsRareTarget(t *testing.T) {
	rares := []string{"T101", "T102", "T103", "T104", "T105"}
	for _, id := range rares {
		if !IsRareTarget(id) {
			t.Errorf("%s should be rare", id)
		}
	}
	normals := []string{"T001", "T002", "T003", "T004", "T005", "T006", "B001"}
	for _, id := range normals {
		if IsRareTarget(id) {
			t.Errorf("%s should not be rare", id)
		}
	}
}

func TestRecordKill_FirstKill(t *testing.T) {
	m := New()
	count, mult, isLevelUp, shouldBroadcast := m.RecordKill("p1")
	if count != 1 {
		t.Errorf("expected count=1, got %d", count)
	}
	if mult != 2.0 {
		t.Errorf("expected mult=2.0, got %.1f", mult)
	}
	if !isLevelUp {
		t.Error("first kill should be level up")
	}
	if shouldBroadcast {
		t.Error("first kill should not broadcast (mult < 5.0)")
	}
}

func TestRecordKill_SecondKill(t *testing.T) {
	m := New()
	m.RecordKill("p1")
	count, mult, isLevelUp, shouldBroadcast := m.RecordKill("p1")
	if count != 2 {
		t.Errorf("expected count=2, got %d", count)
	}
	if mult != 3.0 {
		t.Errorf("expected mult=3.0, got %.1f", mult)
	}
	if !isLevelUp {
		t.Error("second kill should be level up")
	}
	if shouldBroadcast {
		t.Error("second kill should not broadcast (mult < 5.0)")
	}
}

func TestRecordKill_ThirdKill_Broadcast(t *testing.T) {
	m := New()
	m.RecordKill("p1")
	m.RecordKill("p1")
	count, mult, isLevelUp, shouldBroadcast := m.RecordKill("p1")
	if count != 3 {
		t.Errorf("expected count=3, got %d", count)
	}
	if mult != 5.0 {
		t.Errorf("expected mult=5.0, got %.1f", mult)
	}
	if !isLevelUp {
		t.Error("third kill should be level up")
	}
	if !shouldBroadcast {
		t.Error("third kill should broadcast (mult >= 5.0)")
	}
}

func TestRecordKill_MaxCap(t *testing.T) {
	m := New()
	for i := 0; i < 10; i++ {
		m.RecordKill("p1")
	}
	snap := m.GetSnapshot("p1")
	if snap.Count != MaxCascadeCount {
		t.Errorf("expected count capped at %d, got %d", MaxCascadeCount, snap.Count)
	}
	if snap.MultBoost != 15.0 {
		t.Errorf("expected max mult 15.0, got %.1f", snap.MultBoost)
	}
}

func TestGetSnapshot_Inactive(t *testing.T) {
	m := New()
	snap := m.GetSnapshot("p1")
	if snap.IsActive {
		t.Error("snapshot should not be active for new player")
	}
}

func TestGetSnapshot_Active(t *testing.T) {
	m := New()
	m.RecordKill("p1")
	snap := m.GetSnapshot("p1")
	if !snap.IsActive {
		t.Error("snapshot should be active after kill")
	}
	if snap.Count != 1 {
		t.Errorf("expected count=1, got %d", snap.Count)
	}
	if snap.MultBoost != 2.0 {
		t.Errorf("expected mult=2.0, got %.1f", snap.MultBoost)
	}
	if snap.SecondsLeft <= 0 || snap.SecondsLeft > 90 {
		t.Errorf("expected seconds_left in (0,90], got %d", snap.SecondsLeft)
	}
}

func TestCheckExpiry_NotExpired(t *testing.T) {
	m := New()
	m.RecordKill("p1")
	expired := m.CheckExpiry()
	if len(expired) != 0 {
		t.Errorf("expected no expired sessions, got %v", expired)
	}
}

func TestCheckExpiry_Expired(t *testing.T) {
	m := New()
	m.RecordKill("p1")
	// 手動設定過期
	m.mu.Lock()
	m.sessions["p1"].LastHitAt = time.Now().Add(-100 * time.Second)
	m.mu.Unlock()

	expired := m.CheckExpiry()
	if len(expired) != 1 || expired[0] != "p1" {
		t.Errorf("expected p1 to be expired, got %v", expired)
	}
	// 確認已清理
	if m.GetMultBoost("p1") != 1.0 {
		t.Error("expired player should have mult boost 1.0")
	}
}

func TestRemovePlayer(t *testing.T) {
	m := New()
	m.RecordKill("p1")
	m.RemovePlayer("p1")
	if m.GetMultBoost("p1") != 1.0 {
		t.Error("removed player should have mult boost 1.0")
	}
}

func TestMultiplePlayers(t *testing.T) {
	m := New()
	m.RecordKill("p1")
	m.RecordKill("p1")
	m.RecordKill("p2")

	snap1 := m.GetSnapshot("p1")
	snap2 := m.GetSnapshot("p2")

	if snap1.Count != 2 {
		t.Errorf("p1 expected count=2, got %d", snap1.Count)
	}
	if snap2.Count != 1 {
		t.Errorf("p2 expected count=1, got %d", snap2.Count)
	}
	if snap1.MultBoost != 3.0 {
		t.Errorf("p1 expected mult=3.0, got %.1f", snap1.MultBoost)
	}
	if snap2.MultBoost != 2.0 {
		t.Errorf("p2 expected mult=2.0, got %.1f", snap2.MultBoost)
	}
}

func TestCascadeLevels_Complete(t *testing.T) {
	if len(CascadeLevels) != MaxCascadeCount {
		t.Errorf("expected %d cascade levels, got %d", MaxCascadeCount, len(CascadeLevels))
	}
	for i, level := range CascadeLevels {
		if level.MultBoost <= 1.0 {
			t.Errorf("level %d mult boost should be > 1.0, got %.1f", i, level.MultBoost)
		}
		if level.Name == "" {
			t.Errorf("level %d should have a name", i)
		}
		if level.Icon == "" {
			t.Errorf("level %d should have an icon", i)
		}
	}
}

func TestSessionExpiry(t *testing.T) {
	s := &Session{Count: 3, LastHitAt: time.Now().Add(-100 * time.Second)}
	if !s.IsExpired() {
		t.Error("session should be expired")
	}
	if s.GetMultBoost() != 1.0 {
		t.Error("expired session should return mult 1.0")
	}
	if s.GetLevel() != nil {
		t.Error("expired session should return nil level")
	}
}
