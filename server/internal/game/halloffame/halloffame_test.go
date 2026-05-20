package halloffame

import (
	"testing"
)

func TestHallOfFame_TryUpdate_NewRecord(t *testing.T) {
	m := New()
	isNew, old := m.TryUpdate("p1", "Player1", RecordBestStreak, 50, "50連擊", 5, 2)
	if !isNew {
		t.Error("expected new record")
	}
	if old != nil {
		t.Error("expected no old record")
	}
	e := m.GetRecord(RecordBestStreak)
	if e == nil || e.Value != 50 {
		t.Errorf("expected value=50, got %v", e)
	}
}

func TestHallOfFame_TryUpdate_BetterRecord(t *testing.T) {
	m := New()
	m.TryUpdate("p1", "Player1", RecordBestStreak, 50, "50連擊", 5, 2)
	isNew, old := m.TryUpdate("p2", "Player2", RecordBestStreak, 75, "75連擊", 7, 2)
	if !isNew {
		t.Error("expected new record")
	}
	if old == nil || old.PlayerID != "p1" {
		t.Error("expected old record from p1")
	}
	e := m.GetRecord(RecordBestStreak)
	if e.PlayerID != "p2" || e.Value != 75 {
		t.Errorf("expected p2 with 75, got %v", e)
	}
}

func TestHallOfFame_TryUpdate_WorseRecord(t *testing.T) {
	m := New()
	m.TryUpdate("p1", "Player1", RecordBestStreak, 75, "75連擊", 7, 2)
	isNew, _ := m.TryUpdate("p2", "Player2", RecordBestStreak, 50, "50連擊", 5, 2)
	if isNew {
		t.Error("should not update with worse record")
	}
	e := m.GetRecord(RecordBestStreak)
	if e.PlayerID != "p1" {
		t.Error("record should still belong to p1")
	}
}

func TestHallOfFame_GetAll(t *testing.T) {
	m := New()
	m.TryUpdate("p1", "Player1", RecordBestStreak, 50, "50連擊", 5, 2)
	m.TryUpdate("p2", "Player2", RecordBestMultiplier, 100.0, "100x", 10, 3)
	m.TryUpdate("p3", "Player3", RecordMaxCoins, 999999, "999999金幣", 10, 3)

	snap := m.GetAll()
	if len(snap.Records) != 3 {
		t.Errorf("expected 3 records, got %d", len(snap.Records))
	}
	if snap.UpdatedAt == 0 {
		t.Error("expected non-zero updated_at")
	}
}

func TestHallOfFame_IsRecordHolder(t *testing.T) {
	m := New()
	m.TryUpdate("p1", "Player1", RecordBestStreak, 50, "50連擊", 5, 2)

	if !m.IsRecordHolder("p1", RecordBestStreak) {
		t.Error("p1 should be record holder")
	}
	if m.IsRecordHolder("p2", RecordBestStreak) {
		t.Error("p2 should not be record holder")
	}
	if m.IsRecordHolder("p1", RecordMaxCoins) {
		t.Error("p1 should not hold max_coins record")
	}
}

func TestHallOfFame_GetPlayerRecords(t *testing.T) {
	m := New()
	m.TryUpdate("p1", "Player1", RecordBestStreak, 50, "50連擊", 5, 2)
	m.TryUpdate("p1", "Player1", RecordMaxCoins, 100000, "10萬金幣", 10, 3)
	m.TryUpdate("p2", "Player2", RecordBestMultiplier, 100.0, "100x", 10, 3)

	records := m.GetPlayerRecords("p1")
	if len(records) != 2 {
		t.Errorf("expected 2 records for p1, got %d", len(records))
	}
}

func TestHallOfFame_MultipleTypes(t *testing.T) {
	m := New()
	types := []RecordType{
		RecordBestStreak, RecordBestMultiplier, RecordBestBonusReward,
		RecordMostJackpots, RecordGrandJackpot, RecordBossKills,
		RecordMaxCoins, RecordBestRTP,
	}
	for i, rt := range types {
		m.TryUpdate("p1", "Player1", rt, float64(i+1)*10, "test", 5, 2)
	}
	snap := m.GetAll()
	if len(snap.Records) != len(types) {
		t.Errorf("expected %d records, got %d", len(types), len(snap.Records))
	}
}

func TestHallOfFame_RecordTypeLabel(t *testing.T) {
	if RecordTypeLabel(RecordBestStreak) == "" {
		t.Error("expected non-empty label")
	}
	if RecordTypeIcon(RecordGrandJackpot) == "" {
		t.Error("expected non-empty icon")
	}
}

func TestHallOfFame_ConcurrentUpdate(t *testing.T) {
	m := New()
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(val float64) {
			m.TryUpdate("p1", "Player1", RecordBestStreak, val, "test", 5, 2)
			done <- true
		}(float64(i * 10))
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	e := m.GetRecord(RecordBestStreak)
	if e == nil {
		t.Error("expected a record after concurrent updates")
	}
}

func TestHallOfFame_EqualValueNoUpdate(t *testing.T) {
	m := New()
	m.TryUpdate("p1", "Player1", RecordBestStreak, 50, "50連擊", 5, 2)
	isNew, _ := m.TryUpdate("p2", "Player2", RecordBestStreak, 50, "50連擊", 5, 2)
	if isNew {
		t.Error("equal value should not replace existing record")
	}
	e := m.GetRecord(RecordBestStreak)
	if e.PlayerID != "p1" {
		t.Error("original holder should remain")
	}
}
