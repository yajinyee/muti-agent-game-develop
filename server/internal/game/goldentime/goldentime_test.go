package goldentime

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
		t.Error("new manager should not be active")
	}
	if m.GetMultBoost() != 1.0 {
		t.Errorf("expected mult boost 1.0, got %.2f", m.GetMultBoost())
	}
}

func TestCanTrigger_Initial(t *testing.T) {
	m := New()
	if !m.CanTrigger() {
		t.Error("should be able to trigger initially")
	}
}

func TestStart_Silver(t *testing.T) {
	m := New()
	s := m.Start(TierSilver, TriggerRandom)
	if s == nil {
		t.Fatal("Start() returned nil")
	}
	if s.MultBoost != 1.5 {
		t.Errorf("expected mult boost 1.5, got %.2f", s.MultBoost)
	}
	if !m.IsActive() {
		t.Error("should be active after Start()")
	}
	if m.GetMultBoost() != 1.5 {
		t.Errorf("expected GetMultBoost() 1.5, got %.2f", m.GetMultBoost())
	}
}

func TestStart_Gold(t *testing.T) {
	m := New()
	s := m.Start(TierGold, TriggerBossKill)
	if s == nil {
		t.Fatal("Start() returned nil")
	}
	if s.MultBoost != 2.0 {
		t.Errorf("expected mult boost 2.0, got %.2f", s.MultBoost)
	}
}

func TestStart_Rainbow(t *testing.T) {
	m := New()
	s := m.Start(TierRainbow, TriggerRaidVictory)
	if s == nil {
		t.Fatal("Start() returned nil")
	}
	if s.MultBoost != 3.0 {
		t.Errorf("expected mult boost 3.0, got %.2f", s.MultBoost)
	}
}

func TestStart_AlreadyActive(t *testing.T) {
	m := New()
	m.Start(TierGold, TriggerBossKill)
	// 嘗試再次觸發，應該失敗
	s2 := m.Start(TierRainbow, TriggerRandom)
	if s2 != nil {
		t.Error("should not be able to start while active")
	}
}

func TestCanTrigger_WhileActive(t *testing.T) {
	m := New()
	m.Start(TierGold, TriggerBossKill)
	if m.CanTrigger() {
		t.Error("should not be able to trigger while active")
	}
}

func TestCheckExpiry_NotExpired(t *testing.T) {
	m := New()
	m.Start(TierSilver, TriggerRandom)
	if m.CheckExpiry() {
		t.Error("should not expire immediately")
	}
	if !m.IsActive() {
		t.Error("should still be active")
	}
}

func TestCheckExpiry_Expired(t *testing.T) {
	m := New()
	s := m.Start(TierSilver, TriggerRandom)
	if s == nil {
		t.Fatal("Start() returned nil")
	}
	// 手動設定過期
	m.mu.Lock()
	m.session.EndsAt = time.Now().Add(-1 * time.Second)
	m.mu.Unlock()

	if !m.CheckExpiry() {
		t.Error("should detect expiry")
	}
	if m.IsActive() {
		t.Error("should not be active after expiry")
	}
	if m.GetMultBoost() != 1.0 {
		t.Errorf("expected mult boost 1.0 after expiry, got %.2f", m.GetMultBoost())
	}
}

func TestGetSnapshot_Inactive(t *testing.T) {
	m := New()
	snap := m.GetSnapshot()
	if snap.IsActive {
		t.Error("snapshot should not be active")
	}
}

func TestGetSnapshot_Active(t *testing.T) {
	m := New()
	m.Start(TierGold, TriggerBossKill)
	snap := m.GetSnapshot()
	if !snap.IsActive {
		t.Error("snapshot should be active")
	}
	if snap.MultBoost != 2.0 {
		t.Errorf("expected mult boost 2.0, got %.2f", snap.MultBoost)
	}
	if snap.Tier != int(TierGold) {
		t.Errorf("expected tier %d, got %d", TierGold, snap.Tier)
	}
	if snap.TierName == "" {
		t.Error("tier name should not be empty")
	}
	if snap.SecondsLeft <= 0 {
		t.Error("seconds left should be > 0")
	}
	if snap.Color == "" {
		t.Error("color should not be empty")
	}
}

func TestSelectTier_BossKill(t *testing.T) {
	// 多次測試確保只回傳 Gold 或 Rainbow
	for i := 0; i < 100; i++ {
		tier := SelectTier(TriggerBossKill)
		if tier != TierGold && tier != TierRainbow {
			t.Errorf("boss kill should only trigger Gold or Rainbow, got %d", tier)
		}
	}
}

func TestSelectTier_RaidVictory(t *testing.T) {
	// Raid 勝利必定是 Rainbow
	for i := 0; i < 10; i++ {
		tier := SelectTier(TriggerRaidVictory)
		if tier != TierRainbow {
			t.Errorf("raid victory should always trigger Rainbow, got %d", tier)
		}
	}
}

func TestSelectTier_Random(t *testing.T) {
	// 隨機觸發可以是任何等級
	seen := map[Tier]bool{}
	for i := 0; i < 1000; i++ {
		tier := SelectTier(TriggerRandom)
		seen[tier] = true
	}
	// 1000 次應該能看到所有等級
	if !seen[TierSilver] {
		t.Error("random should sometimes trigger Silver")
	}
	if !seen[TierGold] {
		t.Error("random should sometimes trigger Gold")
	}
}

func TestTierDefs_Complete(t *testing.T) {
	tiers := []Tier{TierSilver, TierGold, TierRainbow}
	for _, tier := range tiers {
		def, ok := TierDefs[tier]
		if !ok {
			t.Errorf("missing tier def for %d", tier)
			continue
		}
		if def.Name == "" {
			t.Errorf("tier %d has empty name", tier)
		}
		if def.MultBoost <= 1.0 {
			t.Errorf("tier %d mult boost should be > 1.0, got %.2f", tier, def.MultBoost)
		}
		if def.Duration <= 0 {
			t.Errorf("tier %d duration should be > 0, got %d", tier, def.Duration)
		}
	}
}

func TestSecondsLeft(t *testing.T) {
	m := New()
	m.Start(TierGold, TriggerBossKill)
	snap := m.GetSnapshot()
	if snap.SecondsLeft > 45 || snap.SecondsLeft < 43 {
		t.Errorf("expected ~45 seconds left, got %d", snap.SecondsLeft)
	}
}
