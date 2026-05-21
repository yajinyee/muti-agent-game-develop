package immortalboss

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.IsActive() {
		t.Error("new manager should not have active boss")
	}
}

func TestShouldTrigger_NoCooldown(t *testing.T) {
	m := New()
	// 強制設定 lastAt 為很久以前，確保冷卻已過
	m.lastAt = time.Now().Add(-10 * time.Minute)

	// 多次嘗試，至少有一次應該觸發（機率不為 0）
	triggered := false
	for i := 0; i < 1000; i++ {
		ok, def := m.ShouldTrigger()
		if ok {
			triggered = true
			if def == nil {
				t.Error("triggered but def is nil")
			}
			break
		}
	}
	if !triggered {
		t.Error("expected at least one trigger in 1000 attempts")
	}
}

func TestShouldTrigger_ActiveBoss(t *testing.T) {
	m := New()
	m.lastAt = time.Now().Add(-10 * time.Minute)

	// 啟動一個 session
	def := BossDefs[BossGoldenToad]
	m.StartSession(uuid.New().String(), def)

	// 有活躍 BOSS 時不應觸發
	for i := 0; i < 100; i++ {
		ok, _ := m.ShouldTrigger()
		if ok {
			t.Error("should not trigger when boss is active")
		}
	}
}

func TestShouldTrigger_Cooldown(t *testing.T) {
	m := New()
	// 剛觸發過（冷卻中）
	m.lastAt = time.Now()

	for i := 0; i < 100; i++ {
		ok, _ := m.ShouldTrigger()
		if ok {
			t.Error("should not trigger during cooldown")
		}
	}
}

func TestStartSession(t *testing.T) {
	m := New()
	def := BossDefs[BossGoldenToad]
	instanceID := uuid.New().String()

	s := m.StartSession(instanceID, def)
	if s == nil {
		t.Fatal("StartSession returned nil")
	}
	if s.InstanceID != instanceID {
		t.Errorf("expected instanceID %s, got %s", instanceID, s.InstanceID)
	}
	if s.Def != def {
		t.Error("session def mismatch")
	}
	if !m.IsActive() {
		t.Error("manager should be active after StartSession")
	}
}

func TestRecordHit(t *testing.T) {
	m := New()
	def := BossDefs[BossGoldenToad]
	instanceID := uuid.New().String()
	m.StartSession(instanceID, def)

	mult, reward, ok := m.RecordHit(instanceID, "p1", "Player1", 100)
	if !ok {
		t.Fatal("RecordHit returned false")
	}
	if mult < def.MinMult || mult > def.MaxMult {
		t.Errorf("multiplier %.0f out of range [%.0f, %.0f]", mult, def.MinMult, def.MaxMult)
	}
	if reward != int(mult)*100 {
		t.Errorf("reward %d != mult %.0f * betCost 100", reward, mult)
	}

	snap := m.GetSnapshot()
	if snap.HitCount != 1 {
		t.Errorf("expected HitCount=1, got %d", snap.HitCount)
	}
}

func TestRecordHit_WrongInstance(t *testing.T) {
	m := New()
	def := BossDefs[BossGoldenToad]
	m.StartSession("correct-id", def)

	_, _, ok := m.RecordHit("wrong-id", "p1", "Player1", 100)
	if ok {
		t.Error("should not record hit for wrong instanceID")
	}
}

func TestRecordHit_Expired(t *testing.T) {
	m := New()
	def := &BossDef{
		ID:       BossGoldenToad,
		Name:     "Test",
		MinMult:  50,
		MaxMult:  100,
		Duration: 0.001, // 幾乎立刻過期
	}
	instanceID := uuid.New().String()
	m.StartSession(instanceID, def)

	time.Sleep(10 * time.Millisecond)

	_, _, ok := m.RecordHit(instanceID, "p1", "Player1", 100)
	if ok {
		t.Error("should not record hit for expired session")
	}
}

func TestCheckExpiry(t *testing.T) {
	m := New()
	def := &BossDef{
		ID:       BossGoldenToad,
		Name:     "Test",
		MinMult:  50,
		MaxMult:  100,
		Duration: 0.001,
	}
	instanceID := uuid.New().String()
	m.StartSession(instanceID, def)

	time.Sleep(10 * time.Millisecond)

	expired := m.CheckExpiry()
	if expired == nil {
		t.Fatal("expected expired session")
	}
	if expired.InstanceID != instanceID {
		t.Errorf("expected instanceID %s, got %s", instanceID, expired.InstanceID)
	}
	if m.IsActive() {
		t.Error("manager should not be active after expiry")
	}
}

func TestCheckExpiry_NotExpired(t *testing.T) {
	m := New()
	def := BossDefs[BossAncientCrocodile]
	m.StartSession(uuid.New().String(), def)

	expired := m.CheckExpiry()
	if expired != nil {
		t.Error("should not expire a fresh session")
	}
	if !m.IsActive() {
		t.Error("manager should still be active")
	}
}

func TestGetSnapshot_Inactive(t *testing.T) {
	m := New()
	snap := m.GetSnapshot()
	if snap.Active {
		t.Error("snapshot should be inactive")
	}
}

func TestGetSnapshot_Active(t *testing.T) {
	m := New()
	def := BossDefs[BossGoldenToad]
	instanceID := uuid.New().String()
	m.StartSession(instanceID, def)

	snap := m.GetSnapshot()
	if !snap.Active {
		t.Error("snapshot should be active")
	}
	if snap.BossType != BossGoldenToad {
		t.Errorf("expected BossGoldenToad, got %s", snap.BossType)
	}
	if snap.RemainingSeconds <= 0 {
		t.Error("remaining seconds should be > 0")
	}
}

func TestMultipleHits(t *testing.T) {
	m := New()
	def := BossDefs[BossAncientCrocodile]
	instanceID := uuid.New().String()
	m.StartSession(instanceID, def)

	for i := 0; i < 5; i++ {
		_, _, ok := m.RecordHit(instanceID, "p1", "Player1", 200)
		if !ok {
			t.Fatalf("hit %d failed", i+1)
		}
	}

	snap := m.GetSnapshot()
	if snap.HitCount != 5 {
		t.Errorf("expected HitCount=5, got %d", snap.HitCount)
	}
	if snap.TotalReward <= 0 {
		t.Error("total reward should be > 0")
	}
}

func TestBossDefRanges(t *testing.T) {
	for bossType, def := range BossDefs {
		if def.MinMult <= 0 {
			t.Errorf("%s: MinMult should be > 0", bossType)
		}
		if def.MaxMult <= def.MinMult {
			t.Errorf("%s: MaxMult should be > MinMult", bossType)
		}
		if def.Duration <= 0 {
			t.Errorf("%s: Duration should be > 0", bossType)
		}
		if def.TriggerRate <= 0 || def.TriggerRate >= 1 {
			t.Errorf("%s: TriggerRate should be in (0, 1)", bossType)
		}
	}
}
