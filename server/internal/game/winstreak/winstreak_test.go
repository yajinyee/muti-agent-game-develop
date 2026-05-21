package winstreak

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
}

func TestRecordKill_Basic(t *testing.T) {
	m := New()
	streak, milestone, wasReset := m.RecordKill("p1")
	if streak != 1 {
		t.Errorf("expected streak=1, got %d", streak)
	}
	if milestone != nil {
		t.Error("no milestone should be reached at streak=1")
	}
	if wasReset {
		t.Error("should not reset on first kill")
	}
}

func TestRecordKill_BronzeMilestone(t *testing.T) {
	m := New()
	var lastMilestone *MilestoneDef
	for i := 0; i < 10; i++ {
		_, ms, _ := m.RecordKill("p1")
		if ms != nil {
			lastMilestone = ms
		}
	}
	if lastMilestone == nil {
		t.Fatal("expected bronze milestone at streak=10")
	}
	if lastMilestone.Level != MilestoneBronze {
		t.Errorf("expected MilestoneBronze, got %v", lastMilestone.Level)
	}
}

func TestRecordKill_SilverMilestone(t *testing.T) {
	m := New()
	var lastMilestone *MilestoneDef
	for i := 0; i < 25; i++ {
		_, ms, _ := m.RecordKill("p1")
		if ms != nil {
			lastMilestone = ms
		}
	}
	if lastMilestone == nil {
		t.Fatal("expected silver milestone at streak=25")
	}
	if lastMilestone.Level != MilestoneSilver {
		t.Errorf("expected MilestoneSilver, got %v", lastMilestone.Level)
	}
}

func TestRecordKill_NoDuplicateMilestone(t *testing.T) {
	m := New()
	milestoneCount := 0
	for i := 0; i < 15; i++ {
		_, ms, _ := m.RecordKill("p1")
		if ms != nil && ms.Level == MilestoneBronze {
			milestoneCount++
		}
	}
	if milestoneCount != 1 {
		t.Errorf("bronze milestone should only trigger once, got %d", milestoneCount)
	}
}

func TestRecordKill_Expiry(t *testing.T) {
	m := New()
	// 先打 5 次
	for i := 0; i < 5; i++ {
		m.RecordKill("p1")
	}

	// 強制設定 lastKillAt 為 31 秒前
	m.mu.Lock()
	m.players["p1"].LastKillAt = time.Now().Add(-31 * time.Second)
	m.mu.Unlock()

	// 再打一次，應該重置
	streak, _, wasReset := m.RecordKill("p1")
	if !wasReset {
		t.Error("should reset after timeout")
	}
	if streak != 1 {
		t.Errorf("expected streak=1 after reset, got %d", streak)
	}
}

func TestCheckExpiry(t *testing.T) {
	m := New()
	for i := 0; i < 5; i++ {
		m.RecordKill("p1")
	}

	// 強制過期
	m.mu.Lock()
	m.players["p1"].LastKillAt = time.Now().Add(-31 * time.Second)
	m.mu.Unlock()

	expired := m.CheckExpiry()
	if len(expired) != 1 || expired[0] != "p1" {
		t.Errorf("expected [p1] expired, got %v", expired)
	}

	snap := m.GetSnapshot("p1")
	if snap.Current != 0 {
		t.Errorf("expected Current=0 after expiry, got %d", snap.Current)
	}
}

func TestGetSnapshot(t *testing.T) {
	m := New()
	for i := 0; i < 5; i++ {
		m.RecordKill("p1")
	}

	snap := m.GetSnapshot("p1")
	if snap.Current != 5 {
		t.Errorf("expected Current=5, got %d", snap.Current)
	}
	if snap.NextMilestone != MilestoneBronze {
		t.Errorf("expected NextMilestone=Bronze, got %v", snap.NextMilestone)
	}
	if snap.ProgressToNext <= 0 {
		t.Error("ProgressToNext should be > 0")
	}
}

func TestGetProgressToNext(t *testing.T) {
	m := New()
	// 打 5 次，進度應為 5/10 = 0.5
	for i := 0; i < 5; i++ {
		m.RecordKill("p1")
	}
	snap := m.GetSnapshot("p1")
	expected := 0.5
	if snap.ProgressToNext != expected {
		t.Errorf("expected ProgressToNext=%.1f, got %.1f", expected, snap.ProgressToNext)
	}
}

func TestRemovePlayer(t *testing.T) {
	m := New()
	m.RecordKill("p1")
	m.RemovePlayer("p1")

	snap := m.GetSnapshot("p1")
	if snap.Current != 0 {
		t.Error("removed player should have empty snapshot")
	}
}

func TestMultiplePlayers(t *testing.T) {
	m := New()
	for i := 0; i < 5; i++ {
		m.RecordKill("p1")
	}
	for i := 0; i < 3; i++ {
		m.RecordKill("p2")
	}

	snap1 := m.GetSnapshot("p1")
	snap2 := m.GetSnapshot("p2")

	if snap1.Current != 5 {
		t.Errorf("p1 expected 5, got %d", snap1.Current)
	}
	if snap2.Current != 3 {
		t.Errorf("p2 expected 3, got %d", snap2.Current)
	}
}

func TestMilestoneReset_AfterExpiry(t *testing.T) {
	m := New()
	// 達成銅牌
	for i := 0; i < 10; i++ {
		m.RecordKill("p1")
	}

	// 強制過期
	m.mu.Lock()
	m.players["p1"].LastKillAt = time.Now().Add(-31 * time.Second)
	m.mu.Unlock()
	m.CheckExpiry()

	// 再次達成銅牌（里程碑應重置）
	var bronzeCount int
	for i := 0; i < 10; i++ {
		_, ms, _ := m.RecordKill("p1")
		if ms != nil && ms.Level == MilestoneBronze {
			bronzeCount++
		}
	}
	if bronzeCount != 1 {
		t.Errorf("bronze should trigger again after reset, got %d", bronzeCount)
	}
}
