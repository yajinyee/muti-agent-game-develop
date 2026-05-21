package flashchallenge

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.IsActive() {
		t.Error("new manager should not have active challenge")
	}
}

func TestShouldTrigger_NoChallenge(t *testing.T) {
	m := New()
	m.minIntervalSecs = 0 // 測試時不限制間隔
	// 第一次應該可以觸發（boss 類型必定觸發）
	if !m.ShouldTrigger("boss") {
		t.Error("boss trigger should always succeed when no active challenge")
	}
}

func TestShouldTrigger_ActiveChallenge(t *testing.T) {
	m := New()
	m.minIntervalSecs = 0
	m.StartChallenge()
	// 有進行中的挑戰時不應觸發
	if m.ShouldTrigger("boss") {
		t.Error("should not trigger when challenge is active")
	}
}

func TestShouldTrigger_Cooldown(t *testing.T) {
	m := New()
	m.minIntervalSecs = 300
	m.lastTriggerAt = time.Now() // 剛觸發過
	if m.ShouldTrigger("boss") {
		t.Error("should not trigger during cooldown")
	}
}

func TestStartChallenge(t *testing.T) {
	m := New()
	snap := m.StartChallenge()
	if snap == nil {
		t.Fatal("StartChallenge() returned nil")
	}
	if snap.State != StateActive {
		t.Errorf("expected state=active, got %s", snap.State)
	}
	if snap.Target <= 0 {
		t.Error("target should be > 0")
	}
	if snap.Duration <= 0 {
		t.Error("duration should be > 0")
	}
	if snap.BaseReward <= 0 {
		t.Error("base_reward should be > 0")
	}
}

func TestIsActive(t *testing.T) {
	m := New()
	if m.IsActive() {
		t.Error("should not be active before start")
	}
	m.StartChallenge()
	if !m.IsActive() {
		t.Error("should be active after start")
	}
}

func TestRecordKill_KillCount(t *testing.T) {
	m := New()
	// 強制使用 kill_count 類型
	m.current = &Challenge{
		Def: ChallengeDef{
			Type:        TypeKillCount,
			Target:      3,
			Duration:    90,
			BaseReward:  1000,
			BonusReward: 500,
		},
		State:    StateActive,
		StartedAt: time.Now(),
		EndsAt:   time.Now().Add(90 * time.Second),
		Progress: make(map[string]*PlayerProgress),
	}

	progress, completed, firstComplete := m.RecordKill("p1", "Player1", "T001", 2.0, 1)
	if progress != 1 {
		t.Errorf("expected progress=1, got %d", progress)
	}
	if completed {
		t.Error("should not be completed after 1 kill")
	}
	if firstComplete {
		t.Error("should not be first complete after 1 kill")
	}

	m.RecordKill("p1", "Player1", "T001", 2.0, 2)
	progress, completed, firstComplete = m.RecordKill("p1", "Player1", "T001", 2.0, 3)
	if progress != 3 {
		t.Errorf("expected progress=3, got %d", progress)
	}
	if !completed {
		t.Error("should be completed after 3 kills")
	}
	if !firstComplete {
		t.Error("should be first complete")
	}
}

func TestRecordKill_KillSpecific(t *testing.T) {
	m := New()
	m.current = &Challenge{
		Def: ChallengeDef{
			Type:        TypeKillSpecific,
			Target:      2,
			TargetDefID: "T105",
			Duration:    90,
			BaseReward:  5000,
			BonusReward: 2000,
		},
		State:    StateActive,
		StartedAt: time.Now(),
		EndsAt:   time.Now().Add(90 * time.Second),
		Progress: make(map[string]*PlayerProgress),
	}

	// 擊破非目標不計算
	progress, _, _ := m.RecordKill("p1", "Player1", "T001", 2.0, 1)
	if progress != 0 {
		t.Errorf("wrong target should not count, got progress=%d", progress)
	}

	// 擊破目標計算
	progress, _, _ = m.RecordKill("p1", "Player1", "T105", 10.0, 1)
	if progress != 1 {
		t.Errorf("expected progress=1, got %d", progress)
	}

	progress, completed, _ := m.RecordKill("p1", "Player1", "T105", 10.0, 1)
	if progress != 2 || !completed {
		t.Errorf("expected progress=2 completed=true, got progress=%d completed=%v", progress, completed)
	}
}

func TestRecordKill_HighMult(t *testing.T) {
	m := New()
	m.current = &Challenge{
		Def: ChallengeDef{
			Type:        TypeHighMult,
			Target:      2,
			Duration:    60,
			BaseReward:  6000,
			BonusReward: 3000,
		},
		State:    StateActive,
		StartedAt: time.Now(),
		EndsAt:   time.Now().Add(60 * time.Second),
		Progress: make(map[string]*PlayerProgress),
	}

	// 低倍率不計算
	progress, _, _ := m.RecordKill("p1", "Player1", "T001", 5.0, 1)
	if progress != 0 {
		t.Errorf("low mult should not count, got progress=%d", progress)
	}

	// 高倍率計算
	progress, _, _ = m.RecordKill("p1", "Player1", "T001", 10.0, 1)
	if progress != 1 {
		t.Errorf("expected progress=1, got %d", progress)
	}

	progress, completed, _ := m.RecordKill("p1", "Player1", "T001", 15.0, 1)
	if progress != 2 || !completed {
		t.Errorf("expected progress=2 completed=true, got progress=%d completed=%v", progress, completed)
	}
}

func TestRecordKill_AlreadyCompleted(t *testing.T) {
	m := New()
	m.current = &Challenge{
		Def: ChallengeDef{
			Type:        TypeKillCount,
			Target:      1,
			Duration:    90,
			BaseReward:  1000,
			BonusReward: 500,
		},
		State:    StateActive,
		StartedAt: time.Now(),
		EndsAt:   time.Now().Add(90 * time.Second),
		Progress: make(map[string]*PlayerProgress),
	}

	m.RecordKill("p1", "Player1", "T001", 2.0, 1)
	// 再次擊破不應改變進度
	progress, completed, firstComplete := m.RecordKill("p1", "Player1", "T001", 2.0, 1)
	if progress != 1 {
		t.Errorf("expected progress=1, got %d", progress)
	}
	if !completed {
		t.Error("should still be completed")
	}
	if firstComplete {
		t.Error("should not be first complete again")
	}
}

func TestCheckExpiry(t *testing.T) {
	m := New()
	m.current = &Challenge{
		Def:      ChallengeDef{Type: TypeKillCount, Target: 10},
		State:    StateActive,
		StartedAt: time.Now().Add(-100 * time.Second),
		EndsAt:   time.Now().Add(-1 * time.Second), // 已超時
		Progress: make(map[string]*PlayerProgress),
	}

	expired := m.CheckExpiry()
	if !expired {
		t.Error("should detect expiry")
	}
	if m.current.State != StateFailed {
		t.Errorf("expected state=failed, got %s", m.current.State)
	}
}

func TestCalcReward_Completed(t *testing.T) {
	m := New()
	m.current = &Challenge{
		Def: ChallengeDef{
			Type:        TypeKillCount,
			Target:      10,
			BaseReward:  3000,
			BonusReward: 1500,
		},
		State:    StateActive,
		StartedAt: time.Now(),
		EndsAt:   time.Now().Add(90 * time.Second),
		Progress: make(map[string]*PlayerProgress),
	}

	reward := m.CalcReward(10, true)
	if reward != 4500 {
		t.Errorf("expected reward=4500, got %d", reward)
	}
}

func TestCalcReward_Partial(t *testing.T) {
	m := New()
	m.current = &Challenge{
		Def: ChallengeDef{
			Type:        TypeKillCount,
			Target:      10,
			BaseReward:  3000,
			BonusReward: 1500,
		},
		State:    StateActive,
		StartedAt: time.Now(),
		EndsAt:   time.Now().Add(90 * time.Second),
		Progress: make(map[string]*PlayerProgress),
	}

	// 完成 50%，安慰獎 = 3000 × 0.5 × 0.5 = 750
	reward := m.CalcReward(5, false)
	if reward != 750 {
		t.Errorf("expected reward=750, got %d", reward)
	}
}

func TestCalcReward_TooLow(t *testing.T) {
	m := New()
	m.current = &Challenge{
		Def: ChallengeDef{
			Type:        TypeKillCount,
			Target:      10,
			BaseReward:  3000,
			BonusReward: 1500,
		},
		State:    StateActive,
		StartedAt: time.Now(),
		EndsAt:   time.Now().Add(90 * time.Second),
		Progress: make(map[string]*PlayerProgress),
	}

	// 完成 5%（< 10%），無安慰獎
	reward := m.CalcReward(0, false)
	if reward != 0 {
		t.Errorf("expected reward=0, got %d", reward)
	}
}

func TestGetTopPlayers(t *testing.T) {
	m := New()
	m.current = &Challenge{
		Def: ChallengeDef{
			Type:        TypeKillCount,
			Target:      10,
			Duration:    90,
			BaseReward:  3000,
			BonusReward: 1500,
		},
		State:    StateActive,
		StartedAt: time.Now(),
		EndsAt:   time.Now().Add(90 * time.Second),
		Progress: map[string]*PlayerProgress{
			"p1": {PlayerID: "p1", PlayerName: "Alice", Progress: 8},
			"p2": {PlayerID: "p2", PlayerName: "Bob", Progress: 5},
			"p3": {PlayerID: "p3", PlayerName: "Carol", Progress: 10, Completed: true, CompletedAt: time.Now()},
		},
	}

	snap := m.GetSnapshot()
	if len(snap.TopPlayers) != 3 {
		t.Errorf("expected 3 top players, got %d", len(snap.TopPlayers))
	}
	// 完成者應排第一
	if !snap.TopPlayers[0].Completed {
		t.Error("completed player should be first")
	}
}

func TestMultiplePlayers(t *testing.T) {
	m := New()
	m.current = &Challenge{
		Def: ChallengeDef{
			Type:        TypeKillCount,
			Target:      3,
			Duration:    90,
			BaseReward:  1000,
			BonusReward: 500,
		},
		State:    StateActive,
		StartedAt: time.Now(),
		EndsAt:   time.Now().Add(90 * time.Second),
		Progress: make(map[string]*PlayerProgress),
	}

	m.RecordKill("p1", "Alice", "T001", 2.0, 1)
	m.RecordKill("p2", "Bob", "T001", 2.0, 1)
	m.RecordKill("p1", "Alice", "T001", 2.0, 2)
	m.RecordKill("p2", "Bob", "T001", 2.0, 2)

	p1Progress, _ := m.GetPlayerProgress("p1")
	p2Progress, _ := m.GetPlayerProgress("p2")

	if p1Progress != 2 {
		t.Errorf("p1 expected progress=2, got %d", p1Progress)
	}
	if p2Progress != 2 {
		t.Errorf("p2 expected progress=2, got %d", p2Progress)
	}
}
