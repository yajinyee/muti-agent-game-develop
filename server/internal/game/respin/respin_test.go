package respin

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

func TestShouldTrigger_NoTrigger(t *testing.T) {
	m := New()
	// randFloat = 0.99，遠超觸發機率，不應觸發
	triggered, isChain, chainCount := m.ShouldTrigger("p1", 1, 0.99)
	if triggered {
		t.Error("should not trigger with high randFloat")
	}
	if isChain || chainCount != 0 {
		t.Error("should not be chain")
	}
}

func TestShouldTrigger_Trigger(t *testing.T) {
	m := New()
	// randFloat = 0.001，低於 BaseChance(0.04)，應觸發
	triggered, isChain, chainCount := m.ShouldTrigger("p1", 1, 0.001)
	if !triggered {
		t.Error("should trigger with low randFloat")
	}
	if isChain {
		t.Error("first trigger should not be chain")
	}
	if chainCount != 0 {
		t.Errorf("first trigger chainCount should be 0, got %d", chainCount)
	}
}

func TestShouldTrigger_HighBet(t *testing.T) {
	m := New()
	// LV10 觸發機率 8%，randFloat = 0.07 應觸發
	triggered, _, _ := m.ShouldTrigger("p1", 10, 0.07)
	if !triggered {
		t.Error("LV10 should trigger at 0.07")
	}
}

func TestShouldTrigger_Chain(t *testing.T) {
	m := New()
	// 第一次觸發
	triggered, isChain, chainCount := m.ShouldTrigger("p1", 10, 0.001)
	if !triggered || isChain || chainCount != 0 {
		t.Fatal("first trigger failed")
	}
	// 連鎖觸發（在視窗內，randFloat 很低）
	triggered2, isChain2, chainCount2 := m.ShouldTrigger("p1", 10, 0.001)
	if !triggered2 {
		t.Error("chain trigger should succeed")
	}
	if !isChain2 {
		t.Error("second trigger should be chain")
	}
	if chainCount2 != 1 {
		t.Errorf("chain count should be 1, got %d", chainCount2)
	}
}

func TestShouldTrigger_MaxChain(t *testing.T) {
	m := New()
	// 觸發到最大連鎖
	for i := 0; i < MaxChain; i++ {
		m.ShouldTrigger("p1", 10, 0.001)
	}
	// 再觸發應該失敗（達到最大連鎖）
	triggered, _, _ := m.ShouldTrigger("p1", 10, 0.001)
	if triggered {
		t.Error("should not trigger after max chain")
	}
}

func TestShouldTrigger_Cooldown(t *testing.T) {
	m := New()
	// 觸發後結束 session
	m.ShouldTrigger("p1", 10, 0.001)
	m.EndSession("p1")
	// 冷卻中不應觸發
	triggered, _, _ := m.ShouldTrigger("p1", 10, 0.001)
	if triggered {
		t.Error("should not trigger during cooldown")
	}
}

func TestShouldTrigger_ChainWindowExpired(t *testing.T) {
	m := New()
	// 觸發第一次
	m.ShouldTrigger("p1", 10, 0.001)
	// 手動讓 session 過期
	m.mu.Lock()
	if sess, ok := m.sessions["p1"]; ok {
		sess.LastRespinAt = time.Now().Add(-ChainWindow - time.Second)
	}
	m.mu.Unlock()
	// 連鎖視窗過期，不應連鎖（但可能觸發新的，因為沒有冷卻）
	// 這裡只確認 isChain = false
	_, isChain, _ := m.ShouldTrigger("p1", 10, 0.001)
	if isChain {
		t.Error("should not chain after window expired")
	}
}

func TestGetCurrentMult(t *testing.T) {
	tests := []struct {
		chainCount int
		expected   float64
	}{
		{0, 1.0},
		{1, 1.5},
		{2, 2.0},
		{3, 3.0},
		{4, 5.0},
		{5, 5.0}, // 超出範圍，取最後一個
	}
	for _, tt := range tests {
		sess := &Session{ChainCount: tt.chainCount}
		got := sess.GetCurrentMult()
		if got != tt.expected {
			t.Errorf("chainCount=%d: expected mult=%.1f, got=%.1f", tt.chainCount, tt.expected, got)
		}
	}
}

func TestRemovePlayer(t *testing.T) {
	m := New()
	m.ShouldTrigger("p1", 10, 0.001)
	m.RemovePlayer("p1")
	if m.GetSession("p1") != nil {
		t.Error("session should be removed")
	}
}

func TestGetSession_NoSession(t *testing.T) {
	m := New()
	if m.GetSession("nonexistent") != nil {
		t.Error("should return nil for nonexistent player")
	}
}
