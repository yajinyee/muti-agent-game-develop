// filestore_test.go — FileStore 單元測試（DAY-098）
package store

import (
	"os"
	"testing"
	"time"
)

func TestFileStore_SaveAndLoad(t *testing.T) {
	// 使用暫存目錄
	tmpDir := t.TempDir()
	fs, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	// 儲存完整玩家狀態
	state := &FullPlayerState{
		PlayerID:       "player001",
		DisplayName:    "測試玩家",
		Coins:          50000,
		MaxCoins:       80000,
		BetLevel:       5,
		WeaponLevel:    2,
		KillCount:      123,
		TotalBet:       100000,
		TotalReward:    95000,
		LoginStreak:    7,
		MaxLoginStreak: 14,
		LastLoginDate:  "2026-05-20",
		EquippedSkin:   "season_gold",
		OwnedSkins:     []string{"default", "season_gold"},
		VIPTotalSpend:  200000,
		VIPLevel:       3,
		VIPLastWeeklyAt: time.Now().Add(-8 * 24 * time.Hour),
		SeasonPoints:   1500,
		SeasonLevel:    7,
		SeasonClaimed:  []int{1, 2, 3, 4, 5, 6, 7},
		CodexEntries: []CodexEntryState{
			{TargetID: "T001", Unlocked: true, KillCount: 50, MaxMultiplier: 2.0},
			{TargetID: "B001", Unlocked: true, KillCount: 3, MaxMultiplier: 100.0},
		},
		StatsTotalSessions:  10,
		StatsTotalPlayTime:  3600,
		StatsTotalShots:     5000,
		StatsTotalKills:     123,
		StatsTotalBet:       100000,
		StatsTotalReward:    95000,
		StatsBestMultiplier: 100.0,
		StatsBestStreak:     20,
		StatsJackpotWins:    2,
		StatsJackpotMini:    1,
		StatsJackpotGrand:   1,
	}

	if err := fs.SaveFull(state); err != nil {
		t.Fatalf("SaveFull: %v", err)
	}

	// 讀取並驗證
	loaded, err := fs.LoadFull("player001")
	if err != nil {
		t.Fatalf("LoadFull: %v", err)
	}
	if loaded == nil {
		t.Fatal("LoadFull returned nil")
	}

	if loaded.Coins != 50000 {
		t.Errorf("Coins: got %d, want 50000", loaded.Coins)
	}
	if loaded.VIPLevel != 3 {
		t.Errorf("VIPLevel: got %d, want 3", loaded.VIPLevel)
	}
	if loaded.SeasonPoints != 1500 {
		t.Errorf("SeasonPoints: got %d, want 1500", loaded.SeasonPoints)
	}
	if len(loaded.CodexEntries) != 2 {
		t.Errorf("CodexEntries: got %d, want 2", len(loaded.CodexEntries))
	}
	if loaded.StatsBestMultiplier != 100.0 {
		t.Errorf("StatsBestMultiplier: got %f, want 100.0", loaded.StatsBestMultiplier)
	}
	if loaded.StatsBestStreak != 20 {
		t.Errorf("StatsBestStreak: got %d, want 20", loaded.StatsBestStreak)
	}
	if len(loaded.SeasonClaimed) != 7 {
		t.Errorf("SeasonClaimed: got %d, want 7", len(loaded.SeasonClaimed))
	}
}

func TestFileStore_NewPlayer(t *testing.T) {
	tmpDir := t.TempDir()
	fs, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	// 新玩家應該回傳 nil（不是錯誤）
	loaded, err := fs.LoadFull("new_player_xyz")
	if err != nil {
		t.Fatalf("LoadFull new player: %v", err)
	}
	if loaded != nil {
		t.Error("Expected nil for new player, got non-nil")
	}
}

func TestFileStore_Cache(t *testing.T) {
	tmpDir := t.TempDir()
	fs, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	state := &FullPlayerState{
		PlayerID: "cache_test",
		Coins:    12345,
	}
	if err := fs.SaveFull(state); err != nil {
		t.Fatalf("SaveFull: %v", err)
	}

	// 第一次讀取（從快取）
	loaded1, _ := fs.LoadFull("cache_test")
	// 第二次讀取（仍從快取）
	loaded2, _ := fs.LoadFull("cache_test")

	if loaded1.Coins != loaded2.Coins {
		t.Error("Cache inconsistency")
	}
}

func TestFileStore_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	fs, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	state := &FullPlayerState{PlayerID: "delete_test", Coins: 999}
	fs.SaveFull(state)

	if err := fs.DeletePlayer("delete_test"); err != nil {
		t.Fatalf("DeletePlayer: %v", err)
	}

	// 刪除後應該讀不到
	loaded, _ := fs.LoadFull("delete_test")
	if loaded != nil {
		t.Error("Expected nil after delete, got non-nil")
	}
}

func TestFileStore_AtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	fs, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	// 確認沒有 .tmp 殘留
	state := &FullPlayerState{PlayerID: "atomic_test", Coins: 777}
	if err := fs.SaveFull(state); err != nil {
		t.Fatalf("SaveFull: %v", err)
	}

	tmpPath := fs.playerPath("atomic_test") + ".tmp"
	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		t.Error("Temp file should not exist after successful save")
	}
}

func TestFileStore_SetGetJSON(t *testing.T) {
	tmpDir := t.TempDir()
	fs, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	type TestData struct {
		Value int    `json:"value"`
		Name  string `json:"name"`
	}

	original := TestData{Value: 42, Name: "jackpot_state"}
	if err := fs.SetJSON("jackpot_state", original, 0); err != nil {
		t.Fatalf("SetJSON: %v", err)
	}

	var loaded TestData
	if err := fs.GetJSON("jackpot_state", &loaded); err != nil {
		t.Fatalf("GetJSON: %v", err)
	}
	if loaded.Value != 42 || loaded.Name != "jackpot_state" {
		t.Errorf("GetJSON: got %+v, want %+v", loaded, original)
	}
}

func TestFileStore_Leaderboard(t *testing.T) {
	tmpDir := t.TempDir()
	fs, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	// 儲存幾個玩家
	for _, p := range []struct {
		id    string
		score int64
	}{
		{"p1", 1000},
		{"p2", 5000},
		{"p3", 3000},
	} {
		fs.SaveFull(&FullPlayerState{PlayerID: p.id, Coins: int(p.score)})
		fs.UpdateLeaderboard(p.id, p.score)
	}

	top, err := fs.GetTopPlayers(2)
	if err != nil {
		t.Fatalf("GetTopPlayers: %v", err)
	}
	if len(top) != 2 {
		t.Errorf("GetTopPlayers: got %d, want 2", len(top))
	}
	// 第一名應該是 p2（5000）
	if top[0].PlayerID != "p2" {
		t.Errorf("Top player: got %s, want p2", top[0].PlayerID)
	}
}

func TestFileStore_StoreInterface(t *testing.T) {
	tmpDir := t.TempDir()
	fs, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	// 確認 FileStore 實作 Store 介面
	var _ Store = fs

	// 測試 SavePlayer/LoadPlayer（向下相容）
	state := &PlayerState{
		PlayerID:    "compat_test",
		DisplayName: "相容測試",
		Coins:       9999,
		BetLevel:    3,
		LoginStreak: 5,
		EquippedSkin: "default",
		OwnedSkins:  []string{"default"},
	}
	if err := fs.SavePlayer(state); err != nil {
		t.Fatalf("SavePlayer: %v", err)
	}

	loaded, err := fs.LoadPlayer("compat_test")
	if err != nil {
		t.Fatalf("LoadPlayer: %v", err)
	}
	if loaded == nil {
		t.Fatal("LoadPlayer returned nil")
	}
	if loaded.Coins != 9999 {
		t.Errorf("Coins: got %d, want 9999", loaded.Coins)
	}
	if loaded.LoginStreak != 5 {
		t.Errorf("LoginStreak: got %d, want 5", loaded.LoginStreak)
	}
}
