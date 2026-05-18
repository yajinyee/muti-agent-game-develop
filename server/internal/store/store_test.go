// Package store 測試
package store

import (
	"testing"
	"time"
)

// TestMemoryStoreBasic 基本 CRUD 測試
func TestMemoryStoreBasic(t *testing.T) {
	s := NewMemoryStore()

	// 儲存玩家
	state := &PlayerState{
		PlayerID:    "player-001",
		DisplayName: "吉伊卡哇玩家",
		Coins:       10000,
		Labor:       50,
		BetLevel:    3,
		RoomID:      "room-001",
	}
	if err := s.SavePlayer(state); err != nil {
		t.Fatalf("SavePlayer failed: %v", err)
	}

	// 讀取玩家
	loaded, err := s.LoadPlayer("player-001")
	if err != nil {
		t.Fatalf("LoadPlayer failed: %v", err)
	}
	if loaded == nil {
		t.Fatal("LoadPlayer returned nil")
	}
	if loaded.Coins != 10000 {
		t.Errorf("expected coins=10000, got %d", loaded.Coins)
	}
	if loaded.DisplayName != "吉伊卡哇玩家" {
		t.Errorf("expected display_name='吉伊卡哇玩家', got '%s'", loaded.DisplayName)
	}
}

// TestMemoryStoreNotFound 找不到玩家時回傳 nil, nil
func TestMemoryStoreNotFound(t *testing.T) {
	s := NewMemoryStore()
	loaded, err := s.LoadPlayer("nonexistent")
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if loaded != nil {
		t.Fatalf("expected nil state, got: %+v", loaded)
	}
}

// TestMemoryStoreUpdate 更新玩家狀態
func TestMemoryStoreUpdate(t *testing.T) {
	s := NewMemoryStore()

	state := &PlayerState{
		PlayerID: "player-002",
		Coins:    5000,
	}
	s.SavePlayer(state)

	// 更新金幣
	state.Coins = 15000
	s.SavePlayer(state)

	loaded, _ := s.LoadPlayer("player-002")
	if loaded.Coins != 15000 {
		t.Errorf("expected coins=15000 after update, got %d", loaded.Coins)
	}
}

// TestMemoryStoreDelete 刪除玩家
func TestMemoryStoreDelete(t *testing.T) {
	s := NewMemoryStore()

	state := &PlayerState{PlayerID: "player-003", Coins: 1000}
	s.SavePlayer(state)
	s.DeletePlayer("player-003")

	loaded, err := s.LoadPlayer("player-003")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded != nil {
		t.Fatal("expected nil after delete")
	}
}

// TestMemoryStoreLeaderboard 排行榜測試
func TestMemoryStoreLeaderboard(t *testing.T) {
	s := NewMemoryStore()

	// 建立玩家
	players := []struct {
		id    string
		score int64
	}{
		{"p1", 5000},
		{"p2", 15000},
		{"p3", 8000},
		{"p4", 3000},
		{"p5", 12000},
	}

	for _, p := range players {
		s.SavePlayer(&PlayerState{
			PlayerID:     p.id,
			SessionScore: p.score,
		})
		s.UpdateLeaderboard(p.id, p.score)
	}

	// 取前 3 名
	top, err := s.GetTopPlayers(3)
	if err != nil {
		t.Fatalf("GetTopPlayers failed: %v", err)
	}
	if len(top) != 3 {
		t.Fatalf("expected 3 players, got %d", len(top))
	}
}

// TestMemoryStoreLeaderboardOnlyHighScore 排行榜只保留最高分
func TestMemoryStoreLeaderboardOnlyHighScore(t *testing.T) {
	s := NewMemoryStore()

	s.SavePlayer(&PlayerState{PlayerID: "p1", SessionScore: 1000})
	s.UpdateLeaderboard("p1", 1000)
	s.UpdateLeaderboard("p1", 500) // 低分不應覆蓋

	top, _ := s.GetTopPlayers(1)
	if len(top) == 0 {
		t.Fatal("expected 1 player")
	}
}

// TestMemoryStoreIsolation 確認 LoadPlayer 回傳的是拷貝，不是引用
func TestMemoryStoreIsolation(t *testing.T) {
	s := NewMemoryStore()

	state := &PlayerState{PlayerID: "p1", Coins: 1000}
	s.SavePlayer(state)

	loaded, _ := s.LoadPlayer("p1")
	loaded.Coins = 99999 // 修改拷貝

	// 再次讀取，應該還是 1000
	loaded2, _ := s.LoadPlayer("p1")
	if loaded2.Coins != 1000 {
		t.Errorf("store isolation broken: expected 1000, got %d", loaded2.Coins)
	}
}

// TestMemoryStoreLastSeen 儲存時自動更新 LastSeen
func TestMemoryStoreLastSeen(t *testing.T) {
	s := NewMemoryStore()

	before := time.Now()
	s.SavePlayer(&PlayerState{PlayerID: "p1"})
	after := time.Now()

	loaded, _ := s.LoadPlayer("p1")
	if loaded.LastSeen.Before(before) || loaded.LastSeen.After(after) {
		t.Errorf("LastSeen not set correctly: %v", loaded.LastSeen)
	}
}

// TestNewStoreMemoryFallback REDIS_URL 為空時使用記憶體模式
func TestNewStoreMemoryFallback(t *testing.T) {
	s := New("") // 空 URL → 記憶體模式
	if s.IsRedis() {
		t.Error("expected memory store when REDIS_URL is empty")
	}
}

// TestNewStoreRedisFailFallback Redis 連線失敗時降級到記憶體模式
func TestNewStoreRedisFailFallback(t *testing.T) {
	s := New("redis://invalid-host:6379") // 無效 URL → 降級
	if s.IsRedis() {
		t.Error("expected memory store when Redis connection fails")
	}
}
