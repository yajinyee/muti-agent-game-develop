package season

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
}

func TestGetOrCreate(t *testing.T) {
	m := New()
	data := m.GetOrCreate("player1")
	if data == nil {
		t.Fatal("GetOrCreate returned nil")
	}
	if data.PlayerID != "player1" {
		t.Errorf("expected player1, got %s", data.PlayerID)
	}
	if data.SeasonPoints != 0 {
		t.Errorf("expected 0 points, got %d", data.SeasonPoints)
	}
	if data.CurrentLevel != 0 {
		t.Errorf("expected level 0, got %d", data.CurrentLevel)
	}
}

func TestAddPoints_NoLevelUp(t *testing.T) {
	m := New()
	total, newLevels := m.AddPoints("player1", 50)
	if total != 50 {
		t.Errorf("expected 50, got %d", total)
	}
	if len(newLevels) != 0 {
		t.Errorf("expected no new levels, got %v", newLevels)
	}
}

func TestAddPoints_LevelUp(t *testing.T) {
	m := New()
	// 等級 1 需要 100 積分
	total, newLevels := m.AddPoints("player1", 100)
	if total != 100 {
		t.Errorf("expected 100, got %d", total)
	}
	if len(newLevels) != 1 || newLevels[0] != 1 {
		t.Errorf("expected [1], got %v", newLevels)
	}
}

func TestAddPoints_MultipleLevelUp(t *testing.T) {
	m := New()
	// 等級 1=100, 等級 2=200，一次加 250 積分應該解鎖等級 1 和 2
	total, newLevels := m.AddPoints("player1", 250)
	if total != 250 {
		t.Errorf("expected 250, got %d", total)
	}
	if len(newLevels) != 2 {
		t.Errorf("expected 2 new levels, got %v", newLevels)
	}
}

func TestClaimLevel_Success(t *testing.T) {
	m := New()
	m.AddPoints("player1", 100)
	result, ok := m.ClaimLevel("player1", 1)
	if !ok {
		t.Fatal("ClaimLevel should succeed")
	}
	if result.NewLevel != 1 {
		t.Errorf("expected level 1, got %d", result.NewLevel)
	}
	if result.CoinReward != 500 {
		t.Errorf("expected 500 coins, got %d", result.CoinReward)
	}
}

func TestClaimLevel_AlreadyClaimed(t *testing.T) {
	m := New()
	m.AddPoints("player1", 100)
	m.ClaimLevel("player1", 1)
	_, ok := m.ClaimLevel("player1", 1)
	if ok {
		t.Fatal("ClaimLevel should fail for already claimed level")
	}
}

func TestClaimLevel_InsufficientPoints(t *testing.T) {
	m := New()
	m.AddPoints("player1", 50)
	_, ok := m.ClaimLevel("player1", 1)
	if ok {
		t.Fatal("ClaimLevel should fail for insufficient points")
	}
}

func TestClaimLevel_SpecialReward_Skin(t *testing.T) {
	m := New()
	// 等級 5 需要 800 積分，有皮膚獎勵
	m.AddPoints("player1", 800)
	// 先領取 1-4 等級
	for i := 1; i <= 4; i++ {
		m.ClaimLevel("player1", i)
	}
	result, ok := m.ClaimLevel("player1", 5)
	if !ok {
		t.Fatal("ClaimLevel 5 should succeed")
	}
	if result.SpecialType != "skin" {
		t.Errorf("expected skin, got %s", result.SpecialType)
	}
	if result.SpecialID != "season_gold" {
		t.Errorf("expected season_gold, got %s", result.SpecialID)
	}
}

func TestClaimLevel_SpecialReward_Title(t *testing.T) {
	m := New()
	// 等級 10 需要 3300 積分，有稱號獎勵
	m.AddPoints("player1", 3300)
	// 先領取 1-9 等級
	for i := 1; i <= 9; i++ {
		m.ClaimLevel("player1", i)
	}
	result, ok := m.ClaimLevel("player1", 10)
	if !ok {
		t.Fatal("ClaimLevel 10 should succeed")
	}
	if result.SpecialType != "title" {
		t.Errorf("expected title, got %s", result.SpecialType)
	}
	if result.SpecialID != "season_legend" {
		t.Errorf("expected season_legend, got %s", result.SpecialID)
	}
}

func TestGetSnapshot_NewPlayer(t *testing.T) {
	m := New()
	snap := m.GetSnapshot("player1")
	if snap.SeasonPoints != 0 {
		t.Errorf("expected 0 points, got %d", snap.SeasonPoints)
	}
	if snap.NextLevel != 1 {
		t.Errorf("expected next level 1, got %d", snap.NextLevel)
	}
	if len(snap.Levels) != 10 {
		t.Errorf("expected 10 levels, got %d", len(snap.Levels))
	}
}

func TestGetSnapshot_WithProgress(t *testing.T) {
	m := New()
	m.AddPoints("player1", 150)
	snap := m.GetSnapshot("player1")
	if snap.SeasonPoints != 150 {
		t.Errorf("expected 150 points, got %d", snap.SeasonPoints)
	}
	// 等級 1 已解鎖（100積分）但未領取，NextLevel 應為 1（最低未領取等級）
	if snap.NextLevel != 1 {
		t.Errorf("expected next level 1 (unclaimed), got %d", snap.NextLevel)
	}
	// 等級 1 的 Unlocked 應為 true
	if !snap.Levels[0].Unlocked {
		t.Error("level 1 should be unlocked")
	}
	// 等級 2 的 Unlocked 應為 false（需要 200 積分）
	if snap.Levels[1].Unlocked {
		t.Error("level 2 should not be unlocked with 150 points")
	}
}

func TestSeasonLevels_Count(t *testing.T) {
	if len(SeasonLevels) != 10 {
		t.Errorf("expected 10 season levels, got %d", len(SeasonLevels))
	}
}

func TestSeasonLevels_Ascending(t *testing.T) {
	for i := 1; i < len(SeasonLevels); i++ {
		if SeasonLevels[i].PointsNeeded <= SeasonLevels[i-1].PointsNeeded {
			t.Errorf("level %d points (%d) should be > level %d points (%d)",
				SeasonLevels[i].Level, SeasonLevels[i].PointsNeeded,
				SeasonLevels[i-1].Level, SeasonLevels[i-1].PointsNeeded)
		}
	}
}
