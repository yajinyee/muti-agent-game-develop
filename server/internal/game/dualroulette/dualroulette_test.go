// dualroulette_test.go — 雙環輪盤系統單元測試（DAY-139）
package dualroulette

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := NewDefault()
	if m == nil {
		t.Fatal("NewDefault() returned nil")
	}
	if m.config.MinMultiplier != 30.0 {
		t.Errorf("expected MinMultiplier=30.0, got %.1f", m.config.MinMultiplier)
	}
	if m.config.TriggerChance != 0.15 {
		t.Errorf("expected TriggerChance=0.15, got %.2f", m.config.TriggerChance)
	}
}

func TestCanTrigger_MultTooLow(t *testing.T) {
	m := NewDefault()
	// 倍率低於門檻，不觸發
	for i := 0; i < 100; i++ {
		if m.CanTrigger("p1", 20.0) {
			t.Error("should not trigger when mult < MinMultiplier")
		}
	}
}

func TestCanTrigger_AlwaysTrigger(t *testing.T) {
	// 設定 100% 觸發機率
	cfg := DefaultConfig()
	cfg.TriggerChance = 1.0
	m := New(cfg)
	if !m.CanTrigger("p1", 50.0) {
		t.Error("should trigger with 100% chance and sufficient mult")
	}
}

func TestCanTrigger_NeverTrigger(t *testing.T) {
	// 設定 0% 觸發機率
	cfg := DefaultConfig()
	cfg.TriggerChance = 0.0
	m := New(cfg)
	for i := 0; i < 100; i++ {
		if m.CanTrigger("p1", 50.0) {
			t.Error("should never trigger with 0% chance")
		}
	}
}

func TestStartSession(t *testing.T) {
	m := NewDefault()
	s := m.StartSession("p1", 50.0, 1000)
	if s == nil {
		t.Fatal("StartSession returned nil")
	}
	if s.PlayerID != "p1" {
		t.Errorf("expected PlayerID=p1, got %s", s.PlayerID)
	}
	if s.TargetMult != 50.0 {
		t.Errorf("expected TargetMult=50.0, got %.1f", s.TargetMult)
	}
	if s.BaseReward != 1000 {
		t.Errorf("expected BaseReward=1000, got %d", s.BaseReward)
	}
	if s.IsStopped {
		t.Error("session should not be stopped immediately")
	}
	// 驗證內外環結果在合法範圍內
	validInner := false
	for _, v := range InnerRing {
		if s.InnerResult == v {
			validInner = true
			break
		}
	}
	if !validInner {
		t.Errorf("InnerResult %.1f not in InnerRing", s.InnerResult)
	}
	validOuter := false
	for _, v := range OuterRing {
		if s.OuterResult == v {
			validOuter = true
			break
		}
	}
	if !validOuter {
		t.Errorf("OuterResult %.1f not in OuterRing", s.OuterResult)
	}
}

func TestHasActiveSession(t *testing.T) {
	m := NewDefault()
	if m.HasActiveSession("p1") {
		t.Error("should not have active session before start")
	}
	m.StartSession("p1", 50.0, 1000)
	if !m.HasActiveSession("p1") {
		t.Error("should have active session after start")
	}
}

func TestStopSession_Basic(t *testing.T) {
	m := NewDefault()
	m.StartSession("p1", 50.0, 1000)
	s := m.StopSession("p1")
	if s == nil {
		t.Fatal("StopSession returned nil")
	}
	if !s.IsStopped {
		t.Error("session should be stopped")
	}
	// 驗證最終倍率 = 內環 × 外環
	expected := s.InnerResult * s.OuterResult
	if s.FinalMultiplier() != expected {
		t.Errorf("FinalMultiplier expected %.1f, got %.1f", expected, s.FinalMultiplier())
	}
	// 驗證獎勵 = 基礎獎勵 × 最終倍率
	expectedReward := int(float64(1000) * expected)
	if s.BonusReward() != expectedReward {
		t.Errorf("BonusReward expected %d, got %d", expectedReward, s.BonusReward())
	}
}

func TestStopSession_NoSession(t *testing.T) {
	m := NewDefault()
	s := m.StopSession("p1")
	if s != nil {
		t.Error("StopSession should return nil when no active session")
	}
}

func TestStopSession_AlreadyStopped(t *testing.T) {
	m := NewDefault()
	m.StartSession("p1", 50.0, 1000)
	m.StopSession("p1")
	// 再次停止應該回傳 nil
	s := m.StopSession("p1")
	if s != nil {
		t.Error("second StopSession should return nil")
	}
}

func TestCooldown(t *testing.T) {
	cfg := DefaultConfig()
	cfg.TriggerChance = 1.0
	cfg.CooldownSecs = 60
	m := New(cfg)

	m.StartSession("p1", 50.0, 1000)
	m.StopSession("p1")

	// 停止後應該在冷卻中
	if m.CanTrigger("p1", 50.0) {
		t.Error("should be in cooldown after stop")
	}
	if m.GetCooldownLeft("p1") <= 0 {
		t.Error("cooldown should be > 0 after stop")
	}
}

func TestGetCooldownLeft_NoSession(t *testing.T) {
	m := NewDefault()
	if m.GetCooldownLeft("p1") != 0 {
		t.Error("cooldown should be 0 for new player")
	}
}

func TestCanTrigger_ActiveSession(t *testing.T) {
	cfg := DefaultConfig()
	cfg.TriggerChance = 1.0
	m := New(cfg)
	m.StartSession("p1", 50.0, 1000)
	// 已有活躍 session，不應再觸發
	if m.CanTrigger("p1", 50.0) {
		t.Error("should not trigger when active session exists")
	}
}

func TestRemovePlayer(t *testing.T) {
	m := NewDefault()
	m.StartSession("p1", 50.0, 1000)
	m.RemovePlayer("p1")
	if m.HasActiveSession("p1") {
		t.Error("session should be removed after RemovePlayer")
	}
}

func TestGetSnapshot(t *testing.T) {
	m := NewDefault()
	s := m.StartSession("p1", 50.0, 1000)
	m.StopSession("p1")
	snap := m.GetSnapshot(s)
	if snap.PlayerID != "p1" {
		t.Errorf("expected PlayerID=p1, got %s", snap.PlayerID)
	}
	if snap.TargetMult != 50.0 {
		t.Errorf("expected TargetMult=50.0, got %.1f", snap.TargetMult)
	}
	if !snap.IsStopped {
		t.Error("snapshot should show IsStopped=true")
	}
	if snap.Combined != s.InnerResult*s.OuterResult {
		t.Errorf("snapshot Combined mismatch")
	}
}

func TestTickAutoStop(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SpinDuration = 0.01 // 10ms，方便測試
	m := New(cfg)
	m.StartSession("p1", 50.0, 1000)

	// 等待超時
	time.Sleep(20 * time.Millisecond)
	expired := m.TickAutoStop()
	if len(expired) != 1 {
		t.Errorf("expected 1 expired session, got %d", len(expired))
	}
	if expired[0].PlayerID != "p1" {
		t.Errorf("expected p1, got %s", expired[0].PlayerID)
	}
	if !expired[0].IsStopped {
		t.Error("auto-stopped session should be marked as stopped")
	}
}

func TestTickAutoStop_NotExpired(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SpinDuration = 60.0 // 60 秒，不會超時
	m := New(cfg)
	m.StartSession("p1", 50.0, 1000)
	expired := m.TickAutoStop()
	if len(expired) != 0 {
		t.Errorf("expected 0 expired sessions, got %d", len(expired))
	}
}

func TestMultiplePlayers(t *testing.T) {
	m := NewDefault()
	m.StartSession("p1", 50.0, 1000)
	m.StartSession("p2", 60.0, 2000)

	if !m.HasActiveSession("p1") {
		t.Error("p1 should have active session")
	}
	if !m.HasActiveSession("p2") {
		t.Error("p2 should have active session")
	}

	s1 := m.StopSession("p1")
	if s1 == nil {
		t.Fatal("p1 stop returned nil")
	}
	// p2 仍然活躍
	if !m.HasActiveSession("p2") {
		t.Error("p2 should still have active session after p1 stops")
	}
}

func TestFinalMultiplier_NotStopped(t *testing.T) {
	s := &Session{
		InnerResult: 5.0,
		OuterResult: 10.0,
		IsStopped:   false,
	}
	if s.FinalMultiplier() != 0 {
		t.Error("FinalMultiplier should be 0 when not stopped")
	}
}

func TestBonusReward_NotStopped(t *testing.T) {
	s := &Session{
		BaseReward:  1000,
		InnerResult: 5.0,
		OuterResult: 10.0,
		IsStopped:   false,
	}
	if s.BonusReward() != 0 {
		t.Error("BonusReward should be 0 when not stopped")
	}
}

func TestMaxCombined(t *testing.T) {
	// 最大組合：內環 10x × 外環 15x = 150x
	s := &Session{
		BaseReward:  100,
		InnerResult: 10.0,
		OuterResult: 15.0,
		IsStopped:   true,
	}
	if s.FinalMultiplier() != 150.0 {
		t.Errorf("max combined should be 150.0, got %.1f", s.FinalMultiplier())
	}
	if s.BonusReward() != 15000 {
		t.Errorf("max bonus reward should be 15000, got %d", s.BonusReward())
	}
}
