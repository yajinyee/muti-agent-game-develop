package crystaldragon

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	snap := m.GetSnapshot()
	if snap.TotalCrystals != 0 {
		t.Errorf("expected 0 crystals, got %d", snap.TotalCrystals)
	}
	if snap.Goal != CrystalGoal {
		t.Errorf("expected goal=%d, got %d", CrystalGoal, snap.Goal)
	}
	if snap.CooldownSecs != 0 {
		t.Errorf("expected no cooldown, got %d", snap.CooldownSecs)
	}
}

func TestAddCrystals_Basic(t *testing.T) {
	m := New()
	total, triggered := m.AddCrystals("p1", "Player1", 5, 5)
	if total != 5 {
		t.Errorf("expected total=5, got %d", total)
	}
	if triggered {
		t.Error("should not trigger with only 5 crystals")
	}
}

func TestAddCrystals_MultipleContributors(t *testing.T) {
	m := New()
	m.AddCrystals("p1", "Player1", 5, 10)
	m.AddCrystals("p2", "Player2", 3, 8)
	snap := m.GetSnapshot()
	if snap.TotalCrystals != 18 {
		t.Errorf("expected 18 crystals, got %d", snap.TotalCrystals)
	}
}

func TestAddCrystals_SamePlayerAccumulates(t *testing.T) {
	m := New()
	m.AddCrystals("p1", "Player1", 5, 10)
	m.AddCrystals("p1", "Player1", 5, 10)
	snap := m.GetSnapshot()
	if snap.TotalCrystals != 20 {
		t.Errorf("expected 20 crystals, got %d", snap.TotalCrystals)
	}
}

func TestAddCrystals_TriggerOnGoal(t *testing.T) {
	m := New()
	// 加到剛好達到目標
	total, triggered := m.AddCrystals("p1", "Player1", 10, CrystalGoal)
	if total != CrystalGoal {
		t.Errorf("expected total=%d, got %d", CrystalGoal, total)
	}
	if !triggered {
		t.Error("should trigger when reaching goal")
	}
}

func TestAddCrystals_CapAtGoal(t *testing.T) {
	m := New()
	total, _ := m.AddCrystals("p1", "Player1", 10, CrystalGoal+100)
	if total != CrystalGoal {
		t.Errorf("expected total capped at %d, got %d", CrystalGoal, total)
	}
}

func TestTriggerHellDragon_Basic(t *testing.T) {
	m := New()
	m.AddCrystals("p1", "Player1", 5, 30)
	m.AddCrystals("p2", "Player2", 3, 20)
	// 手動設置達到目標
	m.mu.Lock()
	m.totalCrystals = CrystalGoal
	m.mu.Unlock()

	contributors := m.TriggerHellDragon()
	if len(contributors) != 2 {
		t.Errorf("expected 2 contributors, got %d", len(contributors))
	}
}

func TestTriggerHellDragon_ResetsState(t *testing.T) {
	m := New()
	m.mu.Lock()
	m.totalCrystals = CrystalGoal
	m.contributors["p1"] = &Contributor{PlayerID: "p1", PlayerName: "P1", Crystals: CrystalGoal, BetLevel: 5}
	m.mu.Unlock()

	m.TriggerHellDragon()

	snap := m.GetSnapshot()
	if snap.TotalCrystals != 0 {
		t.Errorf("expected crystals reset to 0, got %d", snap.TotalCrystals)
	}
	if snap.CooldownSecs <= 0 {
		t.Error("expected cooldown after trigger")
	}
}

func TestTriggerHellDragon_NotEnoughCrystals(t *testing.T) {
	m := New()
	m.AddCrystals("p1", "Player1", 5, 10)
	contributors := m.TriggerHellDragon()
	if contributors != nil {
		t.Error("should return nil when not enough crystals")
	}
}

func TestIsOnCooldown(t *testing.T) {
	m := New()
	if m.IsOnCooldown() {
		t.Error("should not be on cooldown initially")
	}

	// 觸發後應該在冷卻中
	m.mu.Lock()
	m.totalCrystals = CrystalGoal
	m.contributors["p1"] = &Contributor{PlayerID: "p1", Crystals: CrystalGoal, BetLevel: 5}
	m.mu.Unlock()
	m.TriggerHellDragon()

	if !m.IsOnCooldown() {
		t.Error("should be on cooldown after trigger")
	}
}

func TestAddCrystals_DuringCooldown(t *testing.T) {
	m := New()
	// 設置冷卻狀態
	m.mu.Lock()
	m.lastTriggerAt = time.Now()
	m.mu.Unlock()

	total, triggered := m.AddCrystals("p1", "Player1", 5, 10)
	if triggered {
		t.Error("should not trigger during cooldown")
	}
	if total != 0 {
		t.Errorf("should not add crystals during cooldown, got %d", total)
	}
}

func TestCheckDecay_NoDecayWhenEmpty(t *testing.T) {
	m := New()
	decayed := m.CheckDecay()
	if decayed {
		t.Error("should not decay when empty")
	}
}

func TestCheckDecay_DecayAfterInterval(t *testing.T) {
	m := New()
	m.AddCrystals("p1", "Player1", 5, 10)
	// 強制設置上次衰減時間為很久以前
	m.mu.Lock()
	m.lastDecayAt = time.Now().Add(-CrystalDecayInterval - time.Second)
	m.mu.Unlock()

	decayed := m.CheckDecay()
	if !decayed {
		t.Error("should decay after interval")
	}
	snap := m.GetSnapshot()
	if snap.TotalCrystals != 9 {
		t.Errorf("expected 9 crystals after decay, got %d", snap.TotalCrystals)
	}
}

func TestCalcReward_Basic(t *testing.T) {
	c := &Contributor{PlayerID: "p1", Crystals: 25, BetLevel: 10}
	reward := CalcReward(c, 50) // 50% 貢獻
	expected := int(0.5 * HellDragonBaseRewardMult * 10)
	if reward != expected {
		t.Errorf("expected reward=%d, got %d", expected, reward)
	}
}

func TestCalcReward_MinimumReward(t *testing.T) {
	c := &Contributor{PlayerID: "p1", Crystals: 1, BetLevel: 5}
	reward := CalcReward(c, 1000) // 極小貢獻
	if reward < c.BetLevel {
		t.Errorf("reward should be at least betLevel=%d, got %d", c.BetLevel, reward)
	}
}

func TestCalcReward_ZeroCrystals(t *testing.T) {
	c := &Contributor{PlayerID: "p1", Crystals: 0, BetLevel: 5}
	reward := CalcReward(c, 50)
	if reward != 0 {
		t.Errorf("expected 0 reward for 0 crystals, got %d", reward)
	}
}

func TestGetSnapshot_Progress(t *testing.T) {
	m := New()
	m.AddCrystals("p1", "Player1", 5, 25) // 50% progress
	snap := m.GetSnapshot()
	expected := 0.5
	if snap.Progress != expected {
		t.Errorf("expected progress=%.2f, got %.2f", expected, snap.Progress)
	}
}

func TestTotalTriggered(t *testing.T) {
	m := New()
	m.mu.Lock()
	m.totalCrystals = CrystalGoal
	m.contributors["p1"] = &Contributor{PlayerID: "p1", Crystals: CrystalGoal, BetLevel: 5}
	m.mu.Unlock()
	m.TriggerHellDragon()

	snap := m.GetSnapshot()
	if snap.TotalTriggered != 1 {
		t.Errorf("expected totalTriggered=1, got %d", snap.TotalTriggered)
	}
}
