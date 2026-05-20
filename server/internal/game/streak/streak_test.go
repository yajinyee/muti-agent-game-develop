package streak

import (
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	snap := m.GetSnapshot()
	if snap.Current != 0 {
		t.Errorf("expected current=0, got %d", snap.Current)
	}
	if snap.MultBonus != 1.0 {
		t.Errorf("expected multBonus=1.0, got %f", snap.MultBonus)
	}
}

func TestRecordKill_Increments(t *testing.T) {
	m := NewManager()
	streak, mult, _ := m.RecordKill()
	if streak != 1 {
		t.Errorf("expected streak=1, got %d", streak)
	}
	if mult != 1.0 {
		t.Errorf("expected mult=1.0 at streak 1, got %f", mult)
	}
}

func TestRecordKill_Level3(t *testing.T) {
	m := NewManager()
	for i := 0; i < 2; i++ {
		m.RecordKill()
	}
	streak, mult, isNew := m.RecordKill() // 3rd kill
	if streak != 3 {
		t.Errorf("expected streak=3, got %d", streak)
	}
	if mult != 1.1 {
		t.Errorf("expected mult=1.1 at streak 3, got %f", mult)
	}
	if !isNew {
		t.Error("expected isNewLevel=true at streak 3")
	}
}

func TestRecordKill_Level20(t *testing.T) {
	m := NewManager()
	for i := 0; i < 19; i++ {
		m.RecordKill()
	}
	streak, mult, isNew := m.RecordKill() // 20th kill
	if streak != 20 {
		t.Errorf("expected streak=20, got %d", streak)
	}
	if mult != 2.0 {
		t.Errorf("expected mult=2.0 at streak 20, got %f", mult)
	}
	if !isNew {
		t.Error("expected isNewLevel=true at streak 20")
	}
}

func TestCheckTimeout_Resets(t *testing.T) {
	m := NewManager()
	m.RecordKill()
	m.RecordKill()
	// 手動設定 lastKillAt 到過去
	m.mu.Lock()
	m.lastKillAt = time.Now().Add(-5 * time.Second)
	m.mu.Unlock()

	reset := m.CheckTimeout()
	if !reset {
		t.Error("expected reset=true after timeout")
	}
	snap := m.GetSnapshot()
	if snap.Current != 0 {
		t.Errorf("expected current=0 after timeout, got %d", snap.Current)
	}
}

func TestCheckTimeout_NoReset(t *testing.T) {
	m := NewManager()
	m.RecordKill()
	reset := m.CheckTimeout()
	if reset {
		t.Error("expected reset=false when not timed out")
	}
}

func TestReset_Force(t *testing.T) {
	m := NewManager()
	for i := 0; i < 5; i++ {
		m.RecordKill()
	}
	m.Reset()
	snap := m.GetSnapshot()
	if snap.Current != 0 {
		t.Errorf("expected current=0 after Reset(), got %d", snap.Current)
	}
	// MaxStreak 應保留
	if snap.MaxStreak != 5 {
		t.Errorf("expected maxStreak=5 after Reset(), got %d", snap.MaxStreak)
	}
}

func TestGetSnapshot_MaxStreak(t *testing.T) {
	m := NewManager()
	for i := 0; i < 8; i++ {
		m.RecordKill()
	}
	// 超時重置
	m.mu.Lock()
	m.lastKillAt = time.Now().Add(-5 * time.Second)
	m.mu.Unlock()
	m.CheckTimeout()
	// 再打 3 個
	for i := 0; i < 3; i++ {
		m.RecordKill()
	}
	snap := m.GetSnapshot()
	if snap.MaxStreak != 8 {
		t.Errorf("expected maxStreak=8, got %d", snap.MaxStreak)
	}
	if snap.Current != 3 {
		t.Errorf("expected current=3, got %d", snap.Current)
	}
}

func TestLevels_Sorted(t *testing.T) {
	for i := 1; i < len(Levels); i++ {
		if Levels[i].MinStreak <= Levels[i-1].MinStreak {
			t.Errorf("Levels not sorted at index %d", i)
		}
		if Levels[i].MultBonus <= Levels[i-1].MultBonus {
			t.Errorf("MultBonus not increasing at index %d", i)
		}
	}
}

func TestRecordKill_AfterTimeout_ResetsFirst(t *testing.T) {
	m := NewManager()
	for i := 0; i < 5; i++ {
		m.RecordKill()
	}
	// 超時
	m.mu.Lock()
	m.lastKillAt = time.Now().Add(-5 * time.Second)
	m.mu.Unlock()
	// 下一次擊破應該從 1 開始
	streak, _, _ := m.RecordKill()
	if streak != 1 {
		t.Errorf("expected streak=1 after timeout+kill, got %d", streak)
	}
}
