package announce

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m.Count() != 0 {
		t.Errorf("Count = %d, want 0", m.Count())
	}
}

func TestCreate_JackpotWin(t *testing.T) {
	m := NewManager()
	ann := m.Create(EventJackpotWin, "TestPlayer", 5000, map[string]string{"level_name": "MINOR"})

	if ann.EventType != EventJackpotWin {
		t.Errorf("EventType = %s, want jackpot_win", ann.EventType)
	}
	if ann.PlayerName != "TestPlayer" {
		t.Errorf("PlayerName = %s, want TestPlayer", ann.PlayerName)
	}
	if ann.Amount != 5000 {
		t.Errorf("Amount = %d, want 5000", ann.Amount)
	}
	if ann.Priority != PriorityHigh {
		t.Errorf("Priority = %d, want %d", ann.Priority, PriorityHigh)
	}
	if ann.Duration <= 0 {
		t.Error("Duration should be > 0")
	}
	if ann.Icon == "" {
		t.Error("Icon should not be empty")
	}
	if ann.Color == "" {
		t.Error("Color should not be empty")
	}
	if ann.ID == "" {
		t.Error("ID should not be empty")
	}
}

func TestCreate_GrandJackpot(t *testing.T) {
	m := NewManager()
	ann := m.Create(EventGrandJackpot, "Winner", 15000, nil)

	if ann.Priority != PriorityCritical {
		t.Errorf("Grand Jackpot priority = %d, want %d", ann.Priority, PriorityCritical)
	}
	if ann.Duration < 6000 {
		t.Errorf("Grand Jackpot duration = %d, want >= 6000ms", ann.Duration)
	}
}

func TestCreate_BossKill(t *testing.T) {
	m := NewManager()
	ann := m.Create(EventBossKill, "Hero", 3000, map[string]string{"boss_name": "海龍王"})

	if ann.EventType != EventBossKill {
		t.Errorf("EventType = %s, want boss_kill", ann.EventType)
	}
	if ann.Priority != PriorityHigh {
		t.Errorf("BossKill priority = %d, want %d", ann.Priority, PriorityHigh)
	}
}

func TestCreate_PlayerJoin(t *testing.T) {
	m := NewManager()
	ann := m.Create(EventPlayerJoin, "NewPlayer", 0, nil)

	if ann.Priority != PriorityLow {
		t.Errorf("PlayerJoin priority = %d, want %d", ann.Priority, PriorityLow)
	}
}

func TestCreate_EmptyPlayerName(t *testing.T) {
	m := NewManager()
	ann := m.Create(EventBigWin, "", 1000, nil)

	// 空名稱應該用預設值
	if ann.Message == "" {
		t.Error("Message should not be empty even with empty player name")
	}
}

func TestGetRecent(t *testing.T) {
	m := NewManager()
	m.Create(EventPlayerJoin, "P1", 0, nil)
	m.Create(EventBigWin, "P2", 1000, nil)
	m.Create(EventBossKill, "P3", 3000, nil)

	recent := m.GetRecent(2)
	if len(recent) != 2 {
		t.Errorf("GetRecent(2) = %d items, want 2", len(recent))
	}
	// 最新的在前
	if recent[0].EventType != EventBossKill {
		t.Errorf("First item should be most recent (boss_kill), got %s", recent[0].EventType)
	}
}

func TestGetRecent_MoreThanAvailable(t *testing.T) {
	m := NewManager()
	m.Create(EventPlayerJoin, "P1", 0, nil)

	recent := m.GetRecent(10)
	if len(recent) != 1 {
		t.Errorf("GetRecent(10) with 1 item = %d, want 1", len(recent))
	}
}

func TestMaxSize(t *testing.T) {
	m := NewManager()
	// 建立超過 maxSize 的公告
	for i := 0; i < 60; i++ {
		m.Create(EventPlayerJoin, "Player", 0, nil)
	}

	if m.Count() > m.maxSize {
		t.Errorf("Count = %d, should not exceed maxSize %d", m.Count(), m.maxSize)
	}
}

func TestUniqueIDs(t *testing.T) {
	m := NewManager()
	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		ann := m.Create(EventPlayerJoin, "P", 0, nil)
		if ids[ann.ID] {
			t.Errorf("Duplicate ID: %s", ann.ID)
		}
		ids[ann.ID] = true
	}
}

func TestAllEventTypes(t *testing.T) {
	m := NewManager()
	events := []EventType{
		EventJackpotWin, EventBigWin, EventMegaWin, EventBossKill,
		EventStreakRecord, EventPlayerJoin, EventPlayerLeave,
		EventWeatherChange, EventEventStart, EventDailyReset,
		EventBossWarning, EventGrandJackpot,
	}
	for _, evt := range events {
		ann := m.Create(evt, "Player", 1000, nil)
		if ann.Title == "" {
			t.Errorf("Event %s has empty title", evt)
		}
		if ann.Icon == "" {
			t.Errorf("Event %s has empty icon", evt)
		}
	}
}
