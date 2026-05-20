package event

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := New(30 * time.Minute)
	if m == nil {
		t.Fatal("New() returned nil")
	}
}

func TestGetCurrent_Initial(t *testing.T) {
	m := New(30 * time.Minute)
	ev := m.GetCurrent()
	if ev == nil {
		t.Fatal("GetCurrent() returned nil")
	}
	// 第一個活動應該是 EventGoldenHour（輪換順序第一個）
	if ev.Type != EventGoldenHour {
		t.Errorf("expected first event=%s, got %s", EventGoldenHour, ev.Type)
	}
}

func TestGetCurrent_IsActive(t *testing.T) {
	m := New(30 * time.Minute)
	ev := m.GetCurrent()
	if !ev.IsActive() {
		t.Error("first event should be active")
	}
}

func TestGetRewardMult_GoldenHour(t *testing.T) {
	m := New(30 * time.Minute)
	// 第一個活動是 GoldenHour，倍率應該是 1.5
	mult := m.GetRewardMult()
	if mult != 1.5 {
		t.Errorf("expected reward_mult=1.5 for GoldenHour, got %f", mult)
	}
}

func TestGetSpawnMult_GoldenHour(t *testing.T) {
	m := New(30 * time.Minute)
	// GoldenHour 不影響生成倍率
	mult := m.GetSpawnMult()
	if mult != 1.0 {
		t.Errorf("expected spawn_mult=1.0 for GoldenHour, got %f", mult)
	}
}

func TestGetKillChanceAdd_GoldenHour(t *testing.T) {
	m := New(30 * time.Minute)
	// GoldenHour 不影響擊破率
	add := m.GetKillChanceAdd()
	if add != 0.0 {
		t.Errorf("expected kill_chance_add=0.0 for GoldenHour, got %f", add)
	}
}

func TestTick_Advance(t *testing.T) {
	// 使用極短的 slot duration 測試切換
	m := New(10 * time.Millisecond)
	// 等待第一個活動過期
	time.Sleep(20 * time.Millisecond)
	changed := m.Tick()
	if !changed {
		t.Error("Tick() should return true when event changes")
	}
	// 第二個應該是 EventNone
	ev := m.GetCurrent()
	if ev.Type != EventNone {
		t.Errorf("expected second event=%s, got %s", EventNone, ev.Type)
	}
}

func TestTick_NoChange(t *testing.T) {
	m := New(30 * time.Minute)
	// 活動還沒過期，Tick 應該回傳 false
	changed := m.Tick()
	if changed {
		t.Error("Tick() should return false when event has not expired")
	}
}

func TestGetSnapshot(t *testing.T) {
	m := New(30 * time.Minute)
	snap := m.GetSnapshot()
	if snap.Type != string(EventGoldenHour) {
		t.Errorf("expected type=%s, got %s", EventGoldenHour, snap.Type)
	}
	if !snap.IsActive {
		t.Error("expected is_active=true")
	}
	if snap.TimeLeft <= 0 {
		t.Error("expected time_left > 0")
	}
	if snap.RewardMult != 1.5 {
		t.Errorf("expected reward_mult=1.5, got %f", snap.RewardMult)
	}
}

func TestEventRotation_Cycle(t *testing.T) {
	// 確認輪換順序正確
	if len(EventRotation) == 0 {
		t.Fatal("EventRotation should not be empty")
	}
	// 確認輪換中有 GoldenHour、FishFrenzy、LuckyMoment
	hasGolden := false
	hasFrenzy := false
	hasLucky := false
	for _, et := range EventRotation {
		switch et {
		case EventGoldenHour:
			hasGolden = true
		case EventFishFrenzy:
			hasFrenzy = true
		case EventLuckyMoment:
			hasLucky = true
		}
	}
	if !hasGolden {
		t.Error("EventRotation should contain GoldenHour")
	}
	if !hasFrenzy {
		t.Error("EventRotation should contain FishFrenzy")
	}
	if !hasLucky {
		t.Error("EventRotation should contain LuckyMoment")
	}
}

func TestEventDefs_Consistency(t *testing.T) {
	for et, def := range EventDefs {
		if def.Type != et {
			t.Errorf("EventDefs[%s].Type should be %s, got %s", et, et, def.Type)
		}
		if def.Name == "" {
			t.Errorf("EventDefs[%s].Name should not be empty", et)
		}
		if def.RewardMult < 1.0 {
			t.Errorf("EventDefs[%s].RewardMult should be >= 1.0", et)
		}
		if def.SpawnMult < 1.0 {
			t.Errorf("EventDefs[%s].SpawnMult should be >= 1.0", et)
		}
		if def.KillChanceAdd < 0 {
			t.Errorf("EventDefs[%s].KillChanceAdd should be >= 0", et)
		}
	}
}

func TestActiveEvent_TimeLeft(t *testing.T) {
	ev := &ActiveEvent{
		Type:    EventGoldenHour,
		StartAt: time.Now(),
		EndAt:   time.Now().Add(30 * time.Minute),
	}
	tl := ev.TimeLeft()
	if tl < 1700 || tl > 1800 {
		t.Errorf("expected time_left≈1800s, got %f", tl)
	}
}

func TestActiveEvent_IsActive_Expired(t *testing.T) {
	ev := &ActiveEvent{
		Type:    EventGoldenHour,
		StartAt: time.Now().Add(-2 * time.Hour),
		EndAt:   time.Now().Add(-1 * time.Hour),
	}
	if ev.IsActive() {
		t.Error("expired event should not be active")
	}
}

func TestNoneEvent_NoBonus(t *testing.T) {
	// 手動設定 EventNone 狀態
	m := New(10 * time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	m.Tick() // 切換到 EventNone

	mult := m.GetRewardMult()
	if mult != 1.0 {
		t.Errorf("expected reward_mult=1.0 for EventNone, got %f", mult)
	}
	spawnMult := m.GetSpawnMult()
	if spawnMult != 1.0 {
		t.Errorf("expected spawn_mult=1.0 for EventNone, got %f", spawnMult)
	}
	killAdd := m.GetKillChanceAdd()
	if killAdd != 0.0 {
		t.Errorf("expected kill_chance_add=0.0 for EventNone, got %f", killAdd)
	}
}
