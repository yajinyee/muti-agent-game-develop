// roulettecrab_test.go — 黃金輪盤螃蟹系統單元測試（DAY-167）
package roulettecrab

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

func TestCanTrigger_Initial(t *testing.T) {
	m := New()
	if !m.CanTrigger("p1") {
		t.Error("should be able to trigger initially")
	}
}

func TestCanTrigger_ActiveSession(t *testing.T) {
	m := New()
	m.StartSession("p1", 30.0, 100)
	if m.CanTrigger("p1") {
		t.Error("should not trigger when session is active")
	}
}

func TestStartSession(t *testing.T) {
	m := New()
	s := m.StartSession("p1", 30.0, 100)
	if s == nil {
		t.Fatal("StartSession returned nil")
	}
	if s.PlayerID != "p1" {
		t.Errorf("expected player_id=p1, got %s", s.PlayerID)
	}
	if s.TargetMult != 30.0 {
		t.Errorf("expected target_mult=30.0, got %.1f", s.TargetMult)
	}
	if s.BaseReward != 100 {
		t.Errorf("expected base_reward=100, got %d", s.BaseReward)
	}
	if s.IsStopped {
		t.Error("session should not be stopped initially")
	}
}

func TestStartSession_WheelResultInRange(t *testing.T) {
	m := New()
	for i := 0; i < 100; i++ {
		s := m.StartSession("p1", 30.0, 100)
		found := false
		for _, slot := range WheelSlots {
			if s.WheelResult == slot {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("wheel_result %.0f not in WheelSlots", s.WheelResult)
		}
		m.StopSession("p1")
	}
}

func TestStartSession_SlotIndexInRange(t *testing.T) {
	m := New()
	for i := 0; i < 100; i++ {
		s := m.StartSession("p1", 30.0, 100)
		if s.SlotIndex < 0 || s.SlotIndex >= len(WheelSlots) {
			t.Errorf("slot_index %d out of range [0, %d)", s.SlotIndex, len(WheelSlots))
		}
		m.StopSession("p1")
	}
}

func TestHasActiveSession(t *testing.T) {
	m := New()
	if m.HasActiveSession("p1") {
		t.Error("should not have active session initially")
	}
	m.StartSession("p1", 30.0, 100)
	if !m.HasActiveSession("p1") {
		t.Error("should have active session after start")
	}
}

func TestStopSession_Basic(t *testing.T) {
	m := New()
	m.StartSession("p1", 30.0, 100)
	s := m.StopSession("p1")
	if s == nil {
		t.Fatal("StopSession returned nil")
	}
	if !s.IsStopped {
		t.Error("session should be stopped")
	}
}

func TestStopSession_NoSession(t *testing.T) {
	m := New()
	s := m.StopSession("p1")
	if s != nil {
		t.Error("StopSession should return nil when no session")
	}
}

func TestStopSession_AlreadyStopped(t *testing.T) {
	m := New()
	m.StartSession("p1", 30.0, 100)
	m.StopSession("p1")
	s := m.StopSession("p1") // 第二次停止
	if s != nil {
		t.Error("second StopSession should return nil")
	}
}

func TestBonusReward(t *testing.T) {
	m := New()
	s := m.StartSession("p1", 30.0, 100)
	// 結果未停止時 BonusReward = 0
	if s.BonusReward() != 0 {
		t.Errorf("expected 0 before stop, got %d", s.BonusReward())
	}
	stopped := m.StopSession("p1")
	if stopped == nil {
		t.Fatal("StopSession returned nil")
	}
	// 停止後 BonusReward = BaseReward × WheelResult
	expected := int(float64(100) * stopped.WheelResult)
	if stopped.BonusReward() != expected {
		t.Errorf("expected bonus_reward=%d, got %d", expected, stopped.BonusReward())
	}
}

func TestCooldown(t *testing.T) {
	m := New()
	m.StartSession("p1", 30.0, 100)
	m.StopSession("p1")
	// 停止後應該在冷卻中
	if m.CanTrigger("p1") {
		t.Error("should be in cooldown after stop")
	}
}

func TestGetCooldownLeft_NoSession(t *testing.T) {
	m := New()
	left := m.GetCooldownLeft("p1")
	if left != 0 {
		t.Errorf("expected 0 cooldown initially, got %d", left)
	}
}

func TestRemovePlayer(t *testing.T) {
	m := New()
	m.StartSession("p1", 30.0, 100)
	m.RemovePlayer("p1")
	if m.HasActiveSession("p1") {
		t.Error("should not have session after remove")
	}
	if !m.CanTrigger("p1") {
		t.Error("should be able to trigger after remove (cooldown cleared)")
	}
}

func TestTickAutoStop(t *testing.T) {
	m := New()
	s := m.StartSession("p1", 30.0, 100)
	// 手動設定 StartedAt 為過去（超過 SpinDuration）
	m.mu.Lock()
	s.StartedAt = time.Now().Add(-time.Duration(SpinDuration*2) * time.Second)
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
	m.StartSession("p1", 30.0, 100)
	expired := m.TickAutoStop()
	if len(expired) != 0 {
		t.Errorf("expected 0 expired sessions, got %d", len(expired))
	}
}

func TestMultiplePlayers(t *testing.T) {
	m := New()
	m.StartSession("p1", 30.0, 100)
	m.StartSession("p2", 40.0, 200)

	if !m.HasActiveSession("p1") {
		t.Error("p1 should have active session")
	}
	if !m.HasActiveSession("p2") {
		t.Error("p2 should have active session")
	}

	m.StopSession("p1")
	if m.HasActiveSession("p1") {
		t.Error("p1 should not have active session after stop")
	}
	if !m.HasActiveSession("p2") {
		t.Error("p2 should still have active session")
	}
}

func TestWheelSlots_Count(t *testing.T) {
	if len(WheelSlots) != 8 {
		t.Errorf("expected 8 wheel slots, got %d", len(WheelSlots))
	}
}

func TestWheelWeights_Count(t *testing.T) {
	if len(WheelWeights) != len(WheelSlots) {
		t.Errorf("WheelWeights count %d != WheelSlots count %d", len(WheelWeights), len(WheelSlots))
	}
}

func TestWheelSlots_Ascending(t *testing.T) {
	for i := 1; i < len(WheelSlots); i++ {
		if WheelSlots[i] <= WheelSlots[i-1] {
			t.Errorf("WheelSlots[%d]=%.0f should be > WheelSlots[%d]=%.0f",
				i, WheelSlots[i], i-1, WheelSlots[i-1])
		}
	}
}

func TestWheelWeights_AllPositive(t *testing.T) {
	for i, w := range WheelWeights {
		if w <= 0 {
			t.Errorf("WheelWeights[%d]=%d should be > 0", i, w)
		}
	}
}

func TestGetSnapshot(t *testing.T) {
	m := New()
	s := m.StartSession("p1", 30.0, 100)
	stopped := m.StopSession("p1")
	snap := m.GetSnapshot(stopped)

	if snap.PlayerID != "p1" {
		t.Errorf("expected player_id=p1, got %s", snap.PlayerID)
	}
	if snap.WheelResult != s.WheelResult {
		t.Errorf("expected wheel_result=%.0f, got %.0f", s.WheelResult, snap.WheelResult)
	}
	if snap.SlotIndex != s.SlotIndex {
		t.Errorf("expected slot_index=%d, got %d", s.SlotIndex, snap.SlotIndex)
	}
	if !snap.IsStopped {
		t.Error("snapshot should show stopped")
	}
}

func TestResultIsPreDetermined(t *testing.T) {
	// 驗證結果在 StartSession 時就已決定（公平性保證）
	m := New()
	s := m.StartSession("p1", 30.0, 100)
	initialResult := s.WheelResult
	initialSlot := s.SlotIndex

	// 停止後結果不變
	stopped := m.StopSession("p1")
	if stopped.WheelResult != initialResult {
		t.Errorf("result changed after stop: %.0f -> %.0f", initialResult, stopped.WheelResult)
	}
	if stopped.SlotIndex != initialSlot {
		t.Errorf("slot changed after stop: %d -> %d", initialSlot, stopped.SlotIndex)
	}
}
