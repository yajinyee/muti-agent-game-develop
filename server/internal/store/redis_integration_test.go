// Redis 整合測試（需要真實 Redis 連線）
// 執行方式：REDIS_URL=redis://localhost:6379 go test ./internal/store/... -run TestRedis -v
// 沒有 Redis 時自動跳過
package store

import (
	"os"
	"testing"
	"time"
)

// getTestRedisURL 從環境變數取得測試用 Redis URL
func getTestRedisURL(t *testing.T) string {
	t.Helper()
	url := os.Getenv("REDIS_URL")
	if url == "" {
		t.Skip("REDIS_URL not set, skipping Redis integration tests")
	}
	return url
}

// TestRedisStoreBasic Redis 基本 CRUD
func TestRedisStoreBasic(t *testing.T) {
	url := getTestRedisURL(t)
	s, err := NewRedisStore(url)
	if err != nil {
		t.Fatalf("NewRedisStore failed: %v", err)
	}
	defer s.Close()

	playerID := "test-redis-" + time.Now().Format("150405")
	state := &PlayerState{
		PlayerID:    playerID,
		DisplayName: "Redis 測試玩家",
		Coins:       88888,
		Labor:       75,
		BetLevel:    5,
		RoomID:      "room-001",
	}

	// 儲存
	if err := s.SavePlayer(state); err != nil {
		t.Fatalf("SavePlayer failed: %v", err)
	}

	// 讀取
	loaded, err := s.LoadPlayer(playerID)
	if err != nil {
		t.Fatalf("LoadPlayer failed: %v", err)
	}
	if loaded == nil {
		t.Fatal("LoadPlayer returned nil")
	}
	if loaded.Coins != 88888 {
		t.Errorf("expected coins=88888, got %d", loaded.Coins)
	}
	if loaded.DisplayName != "Redis 測試玩家" {
		t.Errorf("expected display_name='Redis 測試玩家', got '%s'", loaded.DisplayName)
	}

	// 刪除（清理測試資料）
	if err := s.DeletePlayer(playerID); err != nil {
		t.Fatalf("DeletePlayer failed: %v", err)
	}

	// 確認已刪除
	loaded2, err := s.LoadPlayer(playerID)
	if err != nil {
		t.Fatalf("LoadPlayer after delete failed: %v", err)
	}
	if loaded2 != nil {
		t.Fatal("expected nil after delete")
	}
}

// TestRedisStoreLeaderboard Redis 排行榜
func TestRedisStoreLeaderboard(t *testing.T) {
	url := getTestRedisURL(t)
	s, err := NewRedisStore(url)
	if err != nil {
		t.Fatalf("NewRedisStore failed: %v", err)
	}
	defer s.Close()

	suffix := time.Now().Format("150405")
	players := []struct {
		id    string
		score int64
	}{
		{"redis-p1-" + suffix, 5000},
		{"redis-p2-" + suffix, 15000},
		{"redis-p3-" + suffix, 8000},
	}

	// 建立玩家並更新排行榜
	for _, p := range players {
		if err := s.SavePlayer(&PlayerState{
			PlayerID:     p.id,
			SessionScore: p.score,
		}); err != nil {
			t.Fatalf("SavePlayer failed: %v", err)
		}
		if err := s.UpdateLeaderboard(p.id, p.score); err != nil {
			t.Fatalf("UpdateLeaderboard failed: %v", err)
		}
	}

	// 取前 2 名
	top, err := s.GetTopPlayers(2)
	if err != nil {
		t.Fatalf("GetTopPlayers failed: %v", err)
	}
	if len(top) < 2 {
		t.Fatalf("expected at least 2 players, got %d", len(top))
	}

	// 清理
	for _, p := range players {
		s.DeletePlayer(p.id)
	}
}

// TestRedisStoreLeaderboardHighScoreOnly 排行榜只保留最高分
func TestRedisStoreLeaderboardHighScoreOnly(t *testing.T) {
	url := getTestRedisURL(t)
	s, err := NewRedisStore(url)
	if err != nil {
		t.Fatalf("NewRedisStore failed: %v", err)
	}
	defer s.Close()

	playerID := "redis-hiscore-" + time.Now().Format("150405")
	s.SavePlayer(&PlayerState{PlayerID: playerID, SessionScore: 10000})
	s.UpdateLeaderboard(playerID, 10000)
	s.UpdateLeaderboard(playerID, 5000) // 低分不應覆蓋

	top, err := s.GetTopPlayers(10)
	if err != nil {
		t.Fatalf("GetTopPlayers failed: %v", err)
	}

	for _, p := range top {
		if p.PlayerID == playerID {
			if p.SessionScore < 10000 {
				t.Errorf("expected score >= 10000, got %d", p.SessionScore)
			}
			break
		}
	}

	// 清理
	s.DeletePlayer(playerID)
}

// TestRedisStoreIsRedis 確認 IsRedis() 回傳 true
func TestRedisStoreIsRedis(t *testing.T) {
	url := getTestRedisURL(t)
	s, err := NewRedisStore(url)
	if err != nil {
		t.Fatalf("NewRedisStore failed: %v", err)
	}
	defer s.Close()

	if !s.IsRedis() {
		t.Error("expected IsRedis() = true for RedisStore")
	}
}
