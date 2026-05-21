package awakenboss

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
	m.lastAt = time.Now().Add(-10 * time.Minute)

	triggered := false
	for i := 0; i < 5000; i++ {
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
		t.Error("expected at least one trigger in 5000 attempts")
	}
}

func TestShouldTrigger_ActiveBoss(t *testing.T) {
	m := New()
	m.lastAt = time.Now().Add(-10 * time.Minute)
	def := BossDefs[BossAwakenDragon]
	m.StartSession(uuid.New().String(), def)

	for i := 0; i < 100; i++ {
		ok, _ := m.ShouldTrigger()
		if ok {
			t.Error("should not trigger when boss is active")
		}
	}
}

func TestShouldTrigger_Cooldown(t *testing.T) {
	m := New()
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
	def := BossDefs[BossAwakenDragon]
	instanceID := uuid.New().String()

	s := m.StartSession(instanceID, def)
	if s == nil {
		t.Fatal("StartSession returned nil")
	}
	if s.InstanceID != instanceID {
		t.Errorf("expected instanceID %s, got %s", instanceID, s.InstanceID)
	}
	if !m.IsActive() {
		t.Error("manager should be active after StartSession")
	}
}

func TestRecordHit_Basic(t *testing.T) {
	m := New()
	def := BossDefs[BossAwakenDragon]
	instanceID := uuid.New().String()
	m.StartSession(instanceID, def)

	mult, reward, isPowerUp, ok := m.RecordHit(instanceID, "p1", "Player1", 100)
	if !ok {
		t.Fatal("RecordHit returned false")
	}
	if mult < def.MinMult {
		t.Errorf("multiplier %.0f below MinMult %.0f", mult, def.MinMult)
	}
	if reward <= 0 {
		t.Error("reward should be > 0")
	}
	// 第一次不應該是 Power Up（需要 5 次）
	if isPowerUp {
		t.Error("first hit should not be Power Up")
	}
}

func TestRecordHit_PowerUp(t *testing.T) {
	m := New()
	def := BossDefs[BossAwakenDragon]
	instanceID := uuid.New().String()
	m.StartSession(instanceID, def)

	// 打到第 PowerUpThreshold 次
	var lastIsPowerUp bool
	for i := 0; i < def.PowerUpThreshold; i++ {
		_, _, isPowerUp, ok := m.RecordHit(instanceID, "p1", "Player1", 100)
		if !ok {
			t.Fatalf("hit %d failed", i+1)
		}
		lastIsPowerUp = isPowerUp
	}

	if !lastIsPowerUp {
		t.Errorf("hit %d should be Power Up", def.PowerUpThreshold)
	}

	snap := m.GetSnapshot()
	if snap.PowerUpCount != 1 {
		t.Errorf("expected PowerUpCount=1, got %d", snap.PowerUpCount)
	}
	// Power Up 後進度應重置
	if snap.PowerUpProgress != 0.0 {
		t.Errorf("expected PowerUpProgress=0 after Power Up, got %.2f", snap.PowerUpProgress)
	}
}

func TestRecordHit_PowerUpMultiplier(t *testing.T) {
	m := New()
	def := BossDefs[BossAwakenDragon]
	instanceID := uuid.New().String()
	m.StartSession(instanceID, def)

	// 打到 Power Up
	var powerUpMult float64
	for i := 0; i < def.PowerUpThreshold; i++ {
		mult, _, isPowerUp, _ := m.RecordHit(instanceID, "p1", "Player1", 100)
		if isPowerUp {
			powerUpMult = mult
		}
	}

	// Power Up 倍率應該比基礎倍率高（MinMult × PowerUpMinMult）
	minExpected := def.MinMult * def.PowerUpMinMult
	if powerUpMult < minExpected {
		t.Errorf("Power Up mult %.0f should be >= %.0f", powerUpMult, minExpected)
	}
}

func TestRecordHit_WrongInstance(t *testing.T) {
	m := New()
	def := BossDefs[BossAwakenDragon]
	m.StartSession("correct-id", def)

	_, _, _, ok := m.RecordHit("wrong-id", "p1", "Player1", 100)
	if ok {
		t.Error("should not record hit for wrong instanceID")
	}
}

func TestRecordHit_Expired(t *testing.T) {
	m := New()
	def := &BossDef{
		ID:              BossAwakenDragon,
		Name:            "Test",
		MinMult:         90,
		MaxMult:         180,
		PowerUpMinMult:  6,
		PowerUpMaxMult:  10,
		PowerUpThreshold: 5,
		Duration:        0.001,
	}
	instanceID := uuid.New().String()
	m.StartSession(instanceID, def)
	time.Sleep(10 * time.Millisecond)

	_, _, _, ok := m.RecordHit(instanceID, "p1", "Player1", 100)
	if ok {
		t.Error("should not record hit for expired session")
	}
}

func TestCheckExpiry(t *testing.T) {
	m := New()
	def := &BossDef{
		ID:              BossAwakenDragon,
		Name:            "Test",
		MinMult:         90,
		MaxMult:         180,
		PowerUpMinMult:  6,
		PowerUpMaxMult:  10,
		PowerUpThreshold: 5,
		Duration:        0.001,
	}
	instanceID := uuid.New().String()
	m.StartSession(instanceID, def)
	time.Sleep(10 * time.Millisecond)

	expired := m.CheckExpiry()
	if expired == nil {
		t.Fatal("expected expired session")
	}
	if m.IsActive() {
		t.Error("manager should not be active after expiry")
	}
}

func TestGetSnapshot_Active(t *testing.T) {
	m := New()
	def := BossDefs[BossIcePhoenix]
	instanceID := uuid.New().String()
	m.StartSession(instanceID, def)

	snap := m.GetSnapshot()
	if !snap.Active {
		t.Error("snapshot should be active")
	}
	if snap.BossType != BossIcePhoenix {
		t.Errorf("expected BossIcePhoenix, got %s", snap.BossType)
	}
	if snap.PowerUpThreshold != def.PowerUpThreshold {
		t.Errorf("expected PowerUpThreshold=%d, got %d", def.PowerUpThreshold, snap.PowerUpThreshold)
	}
}

func TestPowerUpProgress(t *testing.T) {
	m := New()
	def := BossDefs[BossAwakenDragon]
	instanceID := uuid.New().String()
	m.StartSession(instanceID, def)

	// 打 2 次，進度應為 2/5 = 0.4
	m.RecordHit(instanceID, "p1", "Player1", 100)
	m.RecordHit(instanceID, "p1", "Player1", 100)

	snap := m.GetSnapshot()
	expected := 2.0 / float64(def.PowerUpThreshold)
	if snap.PowerUpProgress != expected {
		t.Errorf("expected PowerUpProgress=%.2f, got %.2f", expected, snap.PowerUpProgress)
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
		if def.PowerUpMinMult <= 1.0 {
			t.Errorf("%s: PowerUpMinMult should be > 1.0", bossType)
		}
		if def.PowerUpMaxMult <= def.PowerUpMinMult {
			t.Errorf("%s: PowerUpMaxMult should be > PowerUpMinMult", bossType)
		}
		if def.PowerUpThreshold <= 0 {
			t.Errorf("%s: PowerUpThreshold should be > 0", bossType)
		}
		if def.Duration <= 0 {
			t.Errorf("%s: Duration should be > 0", bossType)
		}
		if def.TriggerRate <= 0 || def.TriggerRate >= 1 {
			t.Errorf("%s: TriggerRate should be in (0, 1)", bossType)
		}
	}
}
