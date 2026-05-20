// dm_test.go — 私訊系統測試（DAY-103）
package dm

import (
	"testing"
)

func TestSend_Success_Online(t *testing.T) {
	m := New()
	delivered := false
	result := m.Send("p1", "Player1", "p2", "Hello!", func(msg *Message) bool {
		delivered = true
		return true
	})
	if !result.Success {
		t.Fatalf("expected success, got error: %s", result.ErrorMsg)
	}
	if !delivered {
		t.Error("expected message to be delivered")
	}
	// 不應有離線暫存
	if m.PendingCount("p2") != 0 {
		t.Error("expected no pending messages for online player")
	}
}

func TestSend_Success_Offline(t *testing.T) {
	m := New()
	result := m.Send("p1", "Player1", "p2", "Hello!", func(msg *Message) bool {
		return false // 接收者離線
	})
	if !result.Success {
		t.Fatalf("expected success, got error: %s", result.ErrorMsg)
	}
	if m.PendingCount("p2") != 1 {
		t.Errorf("expected 1 pending message, got %d", m.PendingCount("p2"))
	}
}

func TestSend_SelfMessage(t *testing.T) {
	m := New()
	result := m.Send("p1", "Player1", "p1", "Hello!", nil)
	if result.Success {
		t.Error("expected failure for self message")
	}
	if result.ErrorCode != "self_message" {
		t.Errorf("expected self_message error, got %s", result.ErrorCode)
	}
}

func TestSend_EmptyContent(t *testing.T) {
	m := New()
	result := m.Send("p1", "Player1", "p2", "", nil)
	if result.Success {
		t.Error("expected failure for empty content")
	}
	if result.ErrorCode != "empty_content" {
		t.Errorf("expected empty_content error, got %s", result.ErrorCode)
	}
}

func TestSend_TooLong(t *testing.T) {
	m := New()
	longMsg := ""
	for i := 0; i < 201; i++ {
		longMsg += "A"
	}
	result := m.Send("p1", "Player1", "p2", longMsg, nil)
	if result.Success {
		t.Error("expected failure for too long message")
	}
	if result.ErrorCode != "too_long" {
		t.Errorf("expected too_long error, got %s", result.ErrorCode)
	}
}

func TestGetPending_ClearsMessages(t *testing.T) {
	m := New()
	m.Send("p1", "Player1", "p2", "msg1", func(*Message) bool { return false })
	m.Send("p1", "Player1", "p2", "msg2", func(*Message) bool { return false })

	msgs := m.GetPending("p2")
	if len(msgs) != 2 {
		t.Errorf("expected 2 messages, got %d", len(msgs))
	}
	// 取出後應清除
	if m.PendingCount("p2") != 0 {
		t.Error("expected pending to be cleared after GetPending")
	}
}

func TestGetPending_Empty(t *testing.T) {
	m := New()
	msgs := m.GetPending("p2")
	if msgs != nil {
		t.Error("expected nil for player with no pending messages")
	}
}

func TestMaxStoredMessages(t *testing.T) {
	m := New()
	// 發送超過上限的訊息
	for i := 0; i < MaxStoredMessages+5; i++ {
		m.Send("p1", "Player1", "p2", "msg", func(*Message) bool { return false })
	}
	// 應該只保留最新的 MaxStoredMessages 則
	if m.PendingCount("p2") != MaxStoredMessages {
		t.Errorf("expected %d pending messages, got %d", MaxStoredMessages, m.PendingCount("p2"))
	}
}

func TestDailyCount(t *testing.T) {
	m := New()
	sent, remaining := m.GetDailyCount("p1")
	if sent != 0 || remaining != MaxDailyMessages {
		t.Errorf("expected 0 sent, %d remaining, got %d sent, %d remaining",
			MaxDailyMessages, sent, remaining)
	}

	m.Send("p1", "Player1", "p2", "msg", func(*Message) bool { return true })
	sent, remaining = m.GetDailyCount("p1")
	if sent != 1 || remaining != MaxDailyMessages-1 {
		t.Errorf("expected 1 sent, got %d", sent)
	}
}

func TestMessageID_Unique(t *testing.T) {
	m := New()
	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		result := m.Send("p1", "Player1", "p2", "msg", func(*Message) bool { return true })
		if ids[result.MessageID] {
			t.Errorf("duplicate message ID: %s", result.MessageID)
		}
		ids[result.MessageID] = true
	}
}
