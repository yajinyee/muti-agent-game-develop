package tournament

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tm := New()
	if tm == nil {
		t.Fatal("New() returned nil")
	}
	if len(tm.entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(tm.entries))
	}
}

func TestAddPoints_Kill(t *testing.T) {
	tm := New()
	pts := tm.AddPoints("p1", "Player1", PointKill, 5.0)
	if pts != 5 {
		t.Errorf("expected 5 points, got %d", pts)
	}
}

func TestAddPoints_KillMinimum(t *testing.T) {
	tm := New()
	// 低倍率目標至少 1 分
	pts := tm.AddPoints("p1", "Player1", PointKill, 0.5)
	if pts != 1 {
		t.Errorf("expected 1 point (minimum), got %d", pts)
	}
}

func TestAddPoints_Boss(t *testing.T) {
	tm := New()
	pts := tm.AddPoints("p1", "Player1", PointBoss, 0)
	if pts != 50 {
		t.Errorf("expected 50 points for boss kill, got %d", pts)
	}
}

func TestAddPoints_Bonus(t *testing.T) {
	tm := New()
	pts := tm.AddPoints("p1", "Player1", PointBonus, 0)
	if pts != 20 {
		t.Errorf("expected 20 points for bonus, got %d", pts)
	}
}

func TestAddPoints_Accumulate(t *testing.T) {
	tm := New()
	tm.AddPoints("p1", "Player1", PointKill, 10.0)
	tm.AddPoints("p1", "Player1", PointBoss, 0)
	pts := tm.AddPoints("p1", "Player1", PointBonus, 0)
	// 10 + 50 + 20 = 80
	if pts != 80 {
		t.Errorf("expected 80 accumulated points, got %d", pts)
	}
}

func TestGetRankings_Order(t *testing.T) {
	tm := New()
	tm.AddPoints("p1", "Alice", PointKill, 5.0)   // 5 pts
	tm.AddPoints("p2", "Bob", PointBoss, 0)        // 50 pts
	tm.AddPoints("p3", "Charlie", PointBonus, 0)   // 20 pts

	rankings := tm.GetRankings(3)
	if len(rankings) != 3 {
		t.Fatalf("expected 3 rankings, got %d", len(rankings))
	}
	if rankings[0].PlayerID != "p2" {
		t.Errorf("expected p2 (Bob) at rank 1, got %s", rankings[0].PlayerID)
	}
	if rankings[1].PlayerID != "p3" {
		t.Errorf("expected p3 (Charlie) at rank 2, got %s", rankings[1].PlayerID)
	}
	if rankings[2].PlayerID != "p1" {
		t.Errorf("expected p1 (Alice) at rank 3, got %s", rankings[2].PlayerID)
	}
}

func TestGetRankings_Prizes(t *testing.T) {
	tm := New()
	tm.AddPoints("p1", "Alice", PointBoss, 0)
	tm.AddPoints("p2", "Bob", PointKill, 5.0)
	tm.AddPoints("p3", "Charlie", PointBonus, 0)

	rankings := tm.GetRankings(3)
	// 第一名應有 50000 金幣獎勵
	if rankings[0].Prize != 50000 {
		t.Errorf("expected rank 1 prize 50000, got %d", rankings[0].Prize)
	}
	// 第二名應有 25000 金幣獎勵
	if rankings[1].Prize != 25000 {
		t.Errorf("expected rank 2 prize 25000, got %d", rankings[1].Prize)
	}
}

func TestGetPlayerRank(t *testing.T) {
	tm := New()
	tm.AddPoints("p1", "Alice", PointBoss, 0)   // 50 pts
	tm.AddPoints("p2", "Bob", PointKill, 10.0)  // 10 pts

	rank, pts := tm.GetPlayerRank("p2")
	if rank != 2 {
		t.Errorf("expected rank 2, got %d", rank)
	}
	if pts != 10 {
		t.Errorf("expected 10 points, got %d", pts)
	}
}

func TestGetPlayerRank_NotFound(t *testing.T) {
	tm := New()
	rank, pts := tm.GetPlayerRank("nonexistent")
	if rank != 0 {
		t.Errorf("expected rank 0 for nonexistent player, got %d", rank)
	}
	if pts != 0 {
		t.Errorf("expected 0 points for nonexistent player, got %d", pts)
	}
}

func TestGetWeekInfo(t *testing.T) {
	tm := New()
	start, end, left := tm.GetWeekInfo()
	if start.IsZero() {
		t.Error("week start should not be zero")
	}
	if end.IsZero() {
		t.Error("week end should not be zero")
	}
	if end.Before(start) {
		t.Error("week end should be after week start")
	}
	if left <= 0 {
		t.Error("seconds left should be positive for a new tournament")
	}
}

func TestCurrentWeekRange(t *testing.T) {
	// 測試週一
	loc := time.FixedZone("UTC+8", 8*3600)
	monday := time.Date(2026, 5, 18, 12, 0, 0, 0, loc) // 2026-05-18 是週一
	start, end := currentWeekRange(monday)

	// 開始應該是週一 00:00
	if start.Weekday() != time.Monday {
		t.Errorf("week start should be Monday, got %s", start.Weekday())
	}
	if start.Hour() != 0 || start.Minute() != 0 {
		t.Errorf("week start should be 00:00, got %02d:%02d", start.Hour(), start.Minute())
	}

	// 結束應該是週日 23:59:59
	if end.Weekday() != time.Sunday {
		t.Errorf("week end should be Sunday, got %s", end.Weekday())
	}
	if end.Hour() != 23 || end.Minute() != 59 {
		t.Errorf("week end should be 23:59, got %02d:%02d", end.Hour(), end.Minute())
	}
}

func TestGetSnapshot(t *testing.T) {
	tm := New()
	tm.AddPoints("p1", "Alice", PointBoss, 0)
	tm.AddPoints("p2", "Bob", PointKill, 5.0)

	snap := tm.GetSnapshot()
	if snap.TotalPlayers != 2 {
		t.Errorf("expected 2 total players, got %d", snap.TotalPlayers)
	}
	if len(snap.Rankings) != 2 {
		t.Errorf("expected 2 rankings, got %d", len(snap.Rankings))
	}
	if snap.SecondsLeft <= 0 {
		t.Error("seconds left should be positive")
	}
	if len(snap.Prizes) != 3 {
		t.Errorf("expected 3 prizes, got %d", len(snap.Prizes))
	}
}

func TestDisplayNameUpdate(t *testing.T) {
	tm := New()
	tm.AddPoints("p1", "OldName", PointKill, 5.0)
	tm.AddPoints("p1", "NewName", PointKill, 5.0)

	rankings := tm.GetRankings(1)
	if rankings[0].DisplayName != "NewName" {
		t.Errorf("expected display name 'NewName', got '%s'", rankings[0].DisplayName)
	}
}
