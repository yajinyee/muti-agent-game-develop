// megacatch_test.go — 全服 Mega Catch 事件系統單元測試（DAY-140）
package megacatch

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := NewDefault()
	if m == nil {
		t.Fatal("NewDefault() returned nil")
	}
	if m.config.CooldownSecs != 120 {
		t.Errorf("expected CooldownSecs=120, got %d", m.config.CooldownSecs)
	}
}

func TestIsActive_Initial(t *testing.T) {
	m := NewDefault()
	if m.IsActive() {
		t.Error("should not be active initially")
	}
}

func TestGetRewardBoost_NoEvent(t *testing.T) {
	m := NewDefault()
	if m.GetRewardBoost() != 1.0 {
		t.Errorf("expected 1.0 when no event, got %.1f", m.GetRewardBoost())
	}
}

func TestGetSpawnBoost_NoEvent(t *testing.T) {
	m := NewDefault()
	if m.GetSpawnBoost() != 0.0 {
		t.Errorf("expected 0.0 when no event, got %.2f", m.GetSpawnBoost())
	}
}

func TestCanTrigger_Initial(t *testing.T) {
	m := NewDefault()
	if !m.CanTrigger() {
		t.Error("should be able to trigger initially")
	}
}

func TestForceStart_Basic(t *testing.T) {
	m := NewDefault()
	s := m.ForceStart(0)
	if s == nil {
		t.Fatal("ForceStart returned nil")
	}
	if s.Tier.Name != EventTiers[0].Name {
		t.Errorf("expected tier 0, got %s", s.Tier.Name)
	}
	if !m.IsActive() {
		t.Error("should be active after ForceStart")
	}
}

func TestForceStart_TierIndex(t *testing.T) {
	m := NewDefault()
	s := m.ForceStart(2)
	if s.Tier.Name != EventTiers[2].Name {
		t.Errorf("expected tier 2 (傳說豐收), got %s", s.Tier.Name)
	}
}

func TestForceStart_InvalidIndex(t *testing.T) {
	m := NewDefault()
	s := m.ForceStart(99)
	if s == nil {
		t.Fatal("ForceStart with invalid index should return tier 0")
	}
	if s.Tier.Name != EventTiers[0].Name {
		t.Errorf("expected tier 0 for invalid index, got %s", s.Tier.Name)
	}
}

func TestGetRewardBoost_Active(t *testing.T) {
	m := NewDefault()
	m.ForceStart(1) // 超級豐收 ×2.0
	boost := m.GetRewardBoost()
	if boost != 2.0 {
		t.Errorf("expected RewardBoost=2.0, got %.1f", boost)
	}
}

func TestGetSpawnBoost_Active(t *testing.T) {
	m := NewDefault()
	m.ForceStart(1) // 超級豐收 +35%
	boost := m.GetSpawnBoost()
	if boost != 0.35 {
		t.Errorf("expected SpawnBoost=0.35, got %.2f", boost)
	}
}

func TestCanTrigger_ActiveEvent(t *testing.T) {
	m := NewDefault()
	m.ForceStart(0)
	if m.CanTrigger() {
		t.Error("should not be able to trigger when event is active")
	}
}

func TestCheckExpiry_NotExpired(t *testing.T) {
	m := NewDefault()
	m.ForceStart(0)
	expired := m.CheckExpiry()
	if expired != nil {
		t.Error("should not expire immediately")
	}
}

func TestCheckExpiry_Expired(t *testing.T) {
	cfg := DefaultConfig()
	m := New(cfg)
	// 強制設定一個已過期的 session
	m.mu.Lock()
	m.session = &Session{
		Tier:      EventTiers[0],
		StartedAt: time.Now().Add(-20 * time.Second),
		EndsAt:    time.Now().Add(-1 * time.Second),
	}
	m.mu.Unlock()

	expired := m.CheckExpiry()
	if expired == nil {
		t.Fatal("should return expired session")
	}
	if m.IsActive() {
		t.Error("should not be active after expiry")
	}
}

func TestCheckExpiry_NoSession(t *testing.T) {
	m := NewDefault()
	expired := m.CheckExpiry()
	if expired != nil {
		t.Error("should return nil when no session")
	}
}

func TestGetSnapshot_NoEvent(t *testing.T) {
	m := NewDefault()
	snap := m.GetSnapshot()
	if snap.IsActive {
		t.Error("snapshot should show IsActive=false")
	}
}

func TestGetSnapshot_Active(t *testing.T) {
	m := NewDefault()
	m.ForceStart(2) // 傳說豐收
	snap := m.GetSnapshot()
	if !snap.IsActive {
		t.Error("snapshot should show IsActive=true")
	}
	if snap.TierName != EventTiers[2].Name {
		t.Errorf("expected tier 2 name, got %s", snap.TierName)
	}
	if snap.RewardBoost != 3.0 {
		t.Errorf("expected RewardBoost=3.0, got %.1f", snap.RewardBoost)
	}
	if snap.SecondsLeft <= 0 {
		t.Error("SecondsLeft should be > 0")
	}
}

func TestTryTriggerBossKill_AlwaysTrigger(t *testing.T) {
	cfg := DefaultConfig()
	cfg.BossKillChance = 1.0 // 100% 觸發
	m := New(cfg)
	s := m.TryTriggerBossKill()
	if s == nil {
		t.Fatal("should trigger with 100% chance")
	}
}

func TestTryTriggerBossKill_NeverTrigger(t *testing.T) {
	cfg := DefaultConfig()
	cfg.BossKillChance = 0.0 // 0% 觸發
	m := New(cfg)
	for i := 0; i < 100; i++ {
		s := m.TryTriggerBossKill()
		if s != nil {
			t.Error("should never trigger with 0% chance")
		}
	}
}

func TestTryTriggerBossKill_ActiveEvent(t *testing.T) {
	cfg := DefaultConfig()
	cfg.BossKillChance = 1.0
	m := New(cfg)
	m.ForceStart(0)
	s := m.TryTriggerBossKill()
	if s != nil {
		t.Error("should not trigger when event is active")
	}
}

func TestCooldown(t *testing.T) {
	cfg := DefaultConfig()
	cfg.BossKillChance = 1.0
	cfg.CooldownSecs = 60
	m := New(cfg)
	m.TryTriggerBossKill()

	// 強制讓 session 過期
	m.mu.Lock()
	m.session.EndsAt = time.Now().Add(-1 * time.Second)
	m.mu.Unlock()
	m.CheckExpiry()

	// 應該在冷卻中
	if m.CanTrigger() {
		t.Error("should be in cooldown after event ends")
	}
	if m.GetCooldownLeft() <= 0 {
		t.Error("cooldown should be > 0")
	}
}

func TestEventTiers_Count(t *testing.T) {
	if len(EventTiers) != 3 {
		t.Errorf("expected 3 tiers, got %d", len(EventTiers))
	}
}

func TestEventTiers_RewardBoosts(t *testing.T) {
	expected := []float64{1.5, 2.0, 3.0}
	for i, tier := range EventTiers {
		if tier.RewardBoost != expected[i] {
			t.Errorf("tier %d: expected RewardBoost=%.1f, got %.1f", i, expected[i], tier.RewardBoost)
		}
	}
}
