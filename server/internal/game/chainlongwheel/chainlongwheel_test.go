// chainlongwheel_test.go — 千龍王強化輪盤系統單元測試（DAY-148）
package chainlongwheel

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.sessions == nil {
		t.Fatal("sessions map is nil")
	}
	if m.cooldown == nil {
		t.Fatal("cooldown map is nil")
	}
}

func TestCanTrigger_Initial(t *testing.T) {
	m := New()
	if !m.CanTrigger("player1") {
		t.Error("should be able to trigger initially")
	}
}

func TestCanTrigger_ActiveSession(t *testing.T) {
	m := New()
	m.StartSession("player1", 500.0, 5000)
	if m.CanTrigger("player1") {
		t.Error("should not trigger when session is active")
	}
}

func TestStartSession(t *testing.T) {
	m := New()
	s := m.StartSession("player1", 500.0, 5000)
	if s == nil {
		t.Fatal("StartSession returned nil")
	}
	if s.PlayerID != "player1" {
		t.Errorf("expected player1, got %s", s.PlayerID)
	}
	if s.TargetMult != 500.0 {
		t.Errorf("expected 500.0, got %f", s.TargetMult)
	}
	if s.BaseReward != 5000 {
		t.Errorf("expected 5000, got %d", s.BaseReward)
	}
	if s.IsStopped {
		t.Error("session should not be stopped initially")
	}
	// 驗證內外環結果在有效範圍內
	validInner := false
	for _, v := range InnerRing {
		if s.InnerResult == v {
			validInner = true
			break
		}
	}
	if !validInner {
		t.Errorf("InnerResult %f not in InnerRing", s.InnerResult)
	}
	validOuter := false
	for _, v := range OuterRing {
		if s.OuterResult == v {
			validOuter = true
			break
		}
	}
	if !validOuter {
		t.Errorf("OuterResult %f not in OuterRing", s.OuterResult)
	}
}

func TestHasActiveSession(t *testing.T) {
	m := New()
	if m.HasActiveSession("player1") {
		t.Error("should not have active session initially")
	}
	m.StartSession("player1", 500.0, 5000)
	if !m.HasActiveSession("player1") {
		t.Error("should have active session after start")
	}
}

func TestStopSession_Basic(t *testing.T) {
	m := New()
	m.StartSession("player1", 500.0, 5000)
	s := m.StopSession("player1")
	if s == nil {
		t.Fatal("StopSession returned nil")
	}
	if !s.IsStopped {
		t.Error("session should be stopped")
	}
	if s.FinalMultiplier() == 0 {
		t.Error("FinalMultiplier should not be 0 after stop")
	}
	// 驗證最終倍率 = 內環 × 外環
	expected := s.InnerResult * s.OuterResult
	if s.FinalMultiplier() != expected {
		t.Errorf("expected %f, got %f", expected, s.FinalMultiplier())
	}
}

func TestStopSession_NoSession(t *testing.T) {
	m := New()
	s := m.StopSession("player1")
	if s != nil {
		t.Error("StopSession should return nil when no session")
	}
}

func TestStopSession_AlreadyStopped(t *testing.T) {
	m := New()
	m.StartSession("player1", 500.0, 5000)
	m.StopSession("player1")
	s := m.StopSession("player1") // 第二次停止
	if s != nil {
		t.Error("second StopSession should return nil")
	}
}

func TestBonusReward(t *testing.T) {
	m := New()
	m.StartSession("player1", 500.0, 1000)
	s := m.StopSession("player1")
	if s == nil {
		t.Fatal("StopSession returned nil")
	}
	expected := int(float64(1000) * s.FinalMultiplier())
	if s.BonusReward() != expected {
		t.Errorf("expected %d, got %d", expected, s.BonusReward())
	}
}

func TestCooldown(t *testing.T) {
	m := New()
	m.StartSession("player1", 500.0, 5000)
	m.StopSession("player1")
	// 停止後應該在冷卻中
	if m.CanTrigger("player1") {
		t.Error("should be in cooldown after stop")
	}
	cd := m.GetCooldownLeft("player1")
	if cd <= 0 {
		t.Error("cooldown should be > 0")
	}
}

func TestGetCooldownLeft_NoSession(t *testing.T) {
	m := New()
	cd := m.GetCooldownLeft("player1")
	if cd != 0 {
		t.Errorf("expected 0, got %d", cd)
	}
}

func TestRemovePlayer(t *testing.T) {
	m := New()
	m.StartSession("player1", 500.0, 5000)
	m.RemovePlayer("player1")
	if m.HasActiveSession("player1") {
		t.Error("session should be removed")
	}
}

func TestTickAutoStop(t *testing.T) {
	m := New()
	// 手動建立一個已過期的 session
	m.mu.Lock()
	m.sessions["player1"] = &Session{
		PlayerID:    "player1",
		TargetMult:  500.0,
		BaseReward:  5000,
		InnerResult: 10.0,
		OuterResult: 5.0,
		StartedAt:   time.Now().Add(-10 * time.Second), // 10 秒前開始，已超過 SpinDuration(4s)
	}
	m.mu.Unlock()

	expired := m.TickAutoStop()
	if len(expired) != 1 {
		t.Errorf("expected 1 expired session, got %d", len(expired))
	}
	if !expired[0].IsStopped {
		t.Error("expired session should be stopped")
	}
}

func TestTickAutoStop_NotExpired(t *testing.T) {
	m := New()
	m.StartSession("player1", 500.0, 5000) // 剛開始，未超時
	expired := m.TickAutoStop()
	if len(expired) != 0 {
		t.Errorf("expected 0 expired sessions, got %d", len(expired))
	}
}

func TestMultiplePlayers(t *testing.T) {
	m := New()
	m.StartSession("player1", 500.0, 5000)
	m.StartSession("player2", 300.0, 3000)

	if !m.HasActiveSession("player1") {
		t.Error("player1 should have active session")
	}
	if !m.HasActiveSession("player2") {
		t.Error("player2 should have active session")
	}

	s1 := m.StopSession("player1")
	if s1 == nil {
		t.Fatal("player1 StopSession returned nil")
	}
	// player2 仍然活躍
	if !m.HasActiveSession("player2") {
		t.Error("player2 should still have active session")
	}
}

func TestInnerRing_Count(t *testing.T) {
	if len(InnerRing) != 5 {
		t.Errorf("expected 5 inner ring slots, got %d", len(InnerRing))
	}
}

func TestOuterRing_Count(t *testing.T) {
	if len(OuterRing) != 6 {
		t.Errorf("expected 6 outer ring slots, got %d", len(OuterRing))
	}
}

func TestMaxCombined(t *testing.T) {
	// 最大倍率 = 最大內環 × 最大外環 = 50x × 20x = 1000x
	maxInner := InnerRing[len(InnerRing)-1]
	maxOuter := OuterRing[len(OuterRing)-1]
	maxCombined := maxInner * maxOuter
	if maxCombined != 1000.0 {
		t.Errorf("expected max combined 1000x, got %f", maxCombined)
	}
}

func TestWeights_Count(t *testing.T) {
	if len(InnerWeights) != len(InnerRing) {
		t.Errorf("InnerWeights count %d != InnerRing count %d", len(InnerWeights), len(InnerRing))
	}
	if len(OuterWeights) != len(OuterRing) {
		t.Errorf("OuterWeights count %d != OuterRing count %d", len(OuterWeights), len(OuterRing))
	}
}

func TestResultIsPreDetermined(t *testing.T) {
	// 驗證結果在 StartSession 時就已決定，StopSession 不改變結果
	m := New()
	s := m.StartSession("player1", 500.0, 5000)
	innerBefore := s.InnerResult
	outerBefore := s.OuterResult

	stopped := m.StopSession("player1")
	if stopped.InnerResult != innerBefore {
		t.Error("InnerResult should not change after stop")
	}
	if stopped.OuterResult != outerBefore {
		t.Error("OuterResult should not change after stop")
	}
}

func TestGetSnapshot(t *testing.T) {
	m := New()
	s := m.StartSession("player1", 500.0, 5000)
	stopped := m.StopSession("player1")
	snap := m.GetSnapshot(stopped)

	if snap.PlayerID != "player1" {
		t.Errorf("expected player1, got %s", snap.PlayerID)
	}
	if snap.TargetMult != 500.0 {
		t.Errorf("expected 500.0, got %f", snap.TargetMult)
	}
	if snap.BaseReward != 5000 {
		t.Errorf("expected 5000, got %d", snap.BaseReward)
	}
	if snap.Combined != s.InnerResult*s.OuterResult {
		t.Errorf("combined mismatch")
	}
	if !snap.IsStopped {
		t.Error("snapshot should show stopped")
	}
}
