package activityfeed

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
	if len(m.GetRecent(10)) != 0 {
		t.Error("expected empty feed on init")
	}
}

func TestPush_SingleEvent(t *testing.T) {
	m := New()
	evt := NewAchievementEvent("p1", "Player1", "討伐傳說", "👑", "special")
	pushed := m.Push(evt)
	if pushed.ID == "" {
		t.Error("expected non-empty ID after push")
	}
	if pushed.Timestamp == 0 {
		t.Error("expected non-zero timestamp after push")
	}
	recent := m.GetRecent(10)
	if len(recent) != 1 {
		t.Errorf("expected 1 event, got %d", len(recent))
	}
}

func TestGetRecent_Order(t *testing.T) {
	m := New()
	for i := 0; i < 5; i++ {
		m.Push(NewBossKillEvent("p1", "Player1", i*100))
	}
	recent := m.GetRecent(3)
	if len(recent) != 3 {
		t.Errorf("expected 3 events, got %d", len(recent))
	}
	// 最新的應該在最前面
	if recent[0].Timestamp < recent[1].Timestamp {
		t.Error("expected newest event first")
	}
}

func TestGetRecent_LimitN(t *testing.T) {
	m := New()
	for i := 0; i < 20; i++ {
		m.Push(NewAchievementEvent("p1", "Player1", "test", "⭐", "normal"))
	}
	recent := m.GetRecent(5)
	if len(recent) != 5 {
		t.Errorf("expected 5 events, got %d", len(recent))
	}
}

func TestMaxSize_Overflow(t *testing.T) {
	m := New()
	// 推入超過 maxSize(50) 的事件
	for i := 0; i < 60; i++ {
		m.Push(NewAchievementEvent("p1", "Player1", "test", "⭐", "normal"))
	}
	recent := m.GetRecent(100)
	if len(recent) > 50 {
		t.Errorf("expected max 50 events, got %d", len(recent))
	}
}

func TestNewAchievementEvent_Rarity(t *testing.T) {
	evt := NewAchievementEvent("p1", "Player1", "討伐傳說", "👑", "special")
	if evt.Rarity != RarityEpic {
		t.Errorf("expected epic rarity for special achievement, got %s", evt.Rarity)
	}
	evt2 := NewAchievementEvent("p1", "Player1", "初次討伐", "⚔️", "normal")
	if evt2.Rarity != RarityCommon {
		t.Errorf("expected common rarity for normal achievement, got %s", evt2.Rarity)
	}
}

func TestNewJackpotEvent_Rarity(t *testing.T) {
	grand := NewJackpotEvent("p1", "Player1", "Grand", "👑", 15000)
	if grand.Rarity != RarityLegendary {
		t.Errorf("expected legendary for Grand Jackpot, got %s", grand.Rarity)
	}
	mini := NewJackpotEvent("p1", "Player1", "Mini", "🥈", 300)
	if mini.Rarity != RarityRare {
		t.Errorf("expected rare for Mini Jackpot, got %s", mini.Rarity)
	}
}

func TestNewMegaWinEvent_Rarity(t *testing.T) {
	evt := NewMegaWinEvent("p1", "Player1", 200.0, 10000)
	if evt.Rarity != RarityLegendary {
		t.Errorf("expected legendary for 200x, got %s", evt.Rarity)
	}
	evt2 := NewMegaWinEvent("p1", "Player1", 50.0, 2000)
	if evt2.Rarity != RarityRare {
		t.Errorf("expected rare for 50x, got %s", evt2.Rarity)
	}
}

func TestNewTitleEvent_Rarity(t *testing.T) {
	evt := NewTitleEvent("p1", "Player1", "百日神話", "🌟", 90)
	if evt.Rarity != RarityLegendary {
		t.Errorf("expected legendary for priority 90, got %s", evt.Rarity)
	}
}

func TestNewMilestoneEvent_Rarity(t *testing.T) {
	evt := NewMilestoneEvent("p1", "Player1", 100, "百日神話")
	if evt.Rarity != RarityLegendary {
		t.Errorf("expected legendary for 100-day milestone, got %s", evt.Rarity)
	}
	evt2 := NewMilestoneEvent("p1", "Player1", 3, "初心者")
	if evt2.Rarity != RarityUncommon {
		t.Errorf("expected uncommon for 3-day milestone, got %s", evt2.Rarity)
	}
}

func TestUniqueIDs(t *testing.T) {
	m := New()
	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		evt := m.Push(NewBossKillEvent("p1", "Player1", 1000))
		if ids[evt.ID] {
			t.Errorf("duplicate ID: %s", evt.ID)
		}
		ids[evt.ID] = true
	}
}
