package megaoctopus

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("NewManager returned nil")
	}
	if m.sessions == nil {
		t.Fatal("sessions map is nil")
	}
}

func TestStartSession(t *testing.T) {
	m := NewManager()
	session := m.StartSession("player1")
	if session == nil {
		t.Fatal("StartSession returned nil")
	}
	if session.PlayerID != "player1" {
		t.Errorf("PlayerID = %s, want player1", session.PlayerID)
	}
	if session.ResultIndex < 0 || session.ResultIndex >= len(WheelSlots) {
		t.Errorf("ResultIndex = %d, out of range [0, %d)", session.ResultIndex, len(WheelSlots))
	}
	if session.Stopped {
		t.Error("session should not be stopped initially")
	}
}

func TestHasActiveSession(t *testing.T) {
	m := NewManager()
	if m.HasActiveSession("player1") {
		t.Error("should not have active session before start")
	}
	m.StartSession("player1")
	if !m.HasActiveSession("player1") {
		t.Error("should have active session after start")
	}
}

func TestStopSession_Basic(t *testing.T) {
	m := NewManager()
	m.StartSession("player1")
	resultIndex, multiplier, ok := m.StopSession("player1")
	if !ok {
		t.Fatal("StopSession should succeed")
	}
	if resultIndex < 0 || resultIndex >= len(WheelSlots) {
		t.Errorf("resultIndex = %d, out of range", resultIndex)
	}
	if multiplier != WheelSlots[resultIndex].Multiplier {
		t.Errorf("multiplier = %d, want %d", multiplier, WheelSlots[resultIndex].Multiplier)
	}
	// session 應該被清除
	if m.HasActiveSession("player1") {
		t.Error("session should be removed after stop")
	}
}

func TestStopSession_NoSession(t *testing.T) {
	m := NewManager()
	_, _, ok := m.StopSession("player1")
	if ok {
		t.Error("StopSession should fail when no session exists")
	}
}

func TestStopSession_AlreadyStopped(t *testing.T) {
	m := NewManager()
	m.StartSession("player1")
	m.StopSession("player1")
	_, _, ok := m.StopSession("player1")
	if ok {
		t.Error("StopSession should fail when already stopped")
	}
}

func TestAutoStop(t *testing.T) {
	m := NewManager()
	// 建立一個過期的 session
	m.mu.Lock()
	m.sessions["player1"] = &Session{
		PlayerID:    "player1",
		ResultIndex: 0,
		StartedAt:   time.Now().Add(-10 * time.Second), // 10 秒前開始
		Stopped:     false,
	}
	m.mu.Unlock()

	expired := m.AutoStop()
	if len(expired) != 1 || expired[0] != "player1" {
		t.Errorf("AutoStop should return [player1], got %v", expired)
	}
}

func TestAutoStop_NotExpired(t *testing.T) {
	m := NewManager()
	m.StartSession("player1")
	expired := m.AutoStop()
	if len(expired) != 0 {
		t.Errorf("AutoStop should return empty for fresh session, got %v", expired)
	}
}

func TestRemovePlayer(t *testing.T) {
	m := NewManager()
	m.StartSession("player1")
	m.RemovePlayer("player1")
	if m.HasActiveSession("player1") {
		t.Error("session should be removed after RemovePlayer")
	}
}

func TestWheelSlots_Count(t *testing.T) {
	if len(WheelSlots) != 8 {
		t.Errorf("WheelSlots should have 8 slots, got %d", len(WheelSlots))
	}
}

func TestWheelSlots_Multipliers(t *testing.T) {
	expected := []int{50, 100, 150, 200, 300, 500, 750, 950}
	for i, slot := range WheelSlots {
		if slot.Multiplier != expected[i] {
			t.Errorf("WheelSlots[%d].Multiplier = %d, want %d", i, slot.Multiplier, expected[i])
		}
	}
}

func TestGetResultSlot(t *testing.T) {
	slot := GetResultSlot(0)
	if slot.Multiplier != 50 {
		t.Errorf("GetResultSlot(0).Multiplier = %d, want 50", slot.Multiplier)
	}
	slot = GetResultSlot(7)
	if slot.Multiplier != 950 {
		t.Errorf("GetResultSlot(7).Multiplier = %d, want 950", slot.Multiplier)
	}
}

func TestGetResultSlot_OutOfRange(t *testing.T) {
	slot := GetResultSlot(-1)
	if slot.Multiplier != 50 {
		t.Errorf("GetResultSlot(-1) should return first slot, got %d", slot.Multiplier)
	}
	slot = GetResultSlot(100)
	if slot.Multiplier != 50 {
		t.Errorf("GetResultSlot(100) should return first slot, got %d", slot.Multiplier)
	}
}

func TestMultiplePlayers(t *testing.T) {
	m := NewManager()
	m.StartSession("player1")
	m.StartSession("player2")
	if !m.HasActiveSession("player1") {
		t.Error("player1 should have active session")
	}
	if !m.HasActiveSession("player2") {
		t.Error("player2 should have active session")
	}
	m.StopSession("player1")
	if m.HasActiveSession("player1") {
		t.Error("player1 session should be removed")
	}
	if !m.HasActiveSession("player2") {
		t.Error("player2 session should still be active")
	}
}

func TestResultIsPreDetermined(t *testing.T) {
	// 確認結果在 StartSession 時就已決定，StopSession 只是讀取
	m := NewManager()
	session := m.StartSession("player1")
	expectedIndex := session.ResultIndex

	resultIndex, _, ok := m.StopSession("player1")
	if !ok {
		t.Fatal("StopSession should succeed")
	}
	if resultIndex != expectedIndex {
		t.Errorf("result should be pre-determined: got %d, want %d", resultIndex, expectedIndex)
	}
}
