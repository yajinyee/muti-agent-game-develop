// Package analytics 單元測試
package analytics

import (
	"os"
	"testing"
	"time"
)

func TestTrackerInit(t *testing.T) {
	tmpDir := t.TempDir()
	tracker := newTracker("test-room", tmpDir)
	if tracker == nil {
		t.Fatal("newTracker should return non-nil tracker")
	}
	if tracker.roomID != "test-room" {
		t.Errorf("expected roomID=test-room, got %s", tracker.roomID)
	}
	tracker.Close()
}

func TestTrackPlayerJoinLeave(t *testing.T) {
	tracker := newTracker("test-room", t.TempDir())
	defer tracker.Close()

	// 玩家加入
	tracker.Track(EventPlayerJoin, "player-1", map[string]interface{}{})
	tracker.Track(EventPlayerJoin, "player-2", map[string]interface{}{})

	stats := tracker.GetRoomStats()
	if stats.TotalPlayers != 2 {
		t.Errorf("expected TotalPlayers=2, got %d", stats.TotalPlayers)
	}
	if stats.CurrentPlayers != 2 {
		t.Errorf("expected CurrentPlayers=2, got %d", stats.CurrentPlayers)
	}
	if stats.PeakPlayers != 2 {
		t.Errorf("expected PeakPlayers=2, got %d", stats.PeakPlayers)
	}

	// 玩家離開
	tracker.Track(EventPlayerLeave, "player-1", map[string]interface{}{})
	stats = tracker.GetRoomStats()
	if stats.CurrentPlayers != 1 {
		t.Errorf("expected CurrentPlayers=1 after leave, got %d", stats.CurrentPlayers)
	}
	// PeakPlayers 不應該下降
	if stats.PeakPlayers != 2 {
		t.Errorf("PeakPlayers should remain 2, got %d", stats.PeakPlayers)
	}
}

func TestTrackAttack(t *testing.T) {
	tracker := newTracker("test-room", t.TempDir())
	defer tracker.Close()

	tracker.Track(EventPlayerJoin, "player-1", map[string]interface{}{})

	// 攻擊 5 次
	for i := 0; i < 5; i++ {
		tracker.Track(EventAttack, "player-1", map[string]interface{}{
			"bet_level": 3,
			"bet_cost":  10,
			"is_hit":    true,
			"is_auto":   false,
		})
	}

	stats := tracker.GetRoomStats()
	if stats.TotalAttacks != 5 {
		t.Errorf("expected TotalAttacks=5, got %d", stats.TotalAttacks)
	}

	// 確認 session 統計
	sess := tracker.GetSessionStats("player-1")
	if sess == nil {
		t.Fatal("session should exist for player-1")
	}
	if sess.TotalAttacks != 5 {
		t.Errorf("expected session TotalAttacks=5, got %d", sess.TotalAttacks)
	}
	if sess.TotalBet != 50 {
		t.Errorf("expected session TotalBet=50, got %d", sess.TotalBet)
	}
	if sess.BetLevelDist[3] != 5 {
		t.Errorf("expected BetLevelDist[3]=5, got %d", sess.BetLevelDist[3])
	}
}

func TestTrackKillAndReward(t *testing.T) {
	tracker := newTracker("test-room", t.TempDir())
	defer tracker.Close()

	tracker.Track(EventPlayerJoin, "player-1", map[string]interface{}{})

	tracker.Track(EventKill, "player-1", map[string]interface{}{
		"def_id":      "T001",
		"target_type": "normal",
		"multiplier":  2.0,
		"reward":      20,
	})
	tracker.Track(EventReward, "player-1", map[string]interface{}{
		"source":     "target",
		"amount":     20,
		"multiplier": 2.0,
	})

	stats := tracker.GetRoomStats()
	if stats.TotalKills != 1 {
		t.Errorf("expected TotalKills=1, got %d", stats.TotalKills)
	}
	if stats.TotalReward != 20 {
		t.Errorf("expected TotalReward=20, got %d", stats.TotalReward)
	}

	sess := tracker.GetSessionStats("player-1")
	if sess.TotalKills != 1 {
		t.Errorf("expected session TotalKills=1, got %d", sess.TotalKills)
	}
	if sess.TotalReward != 20 {
		t.Errorf("expected session TotalReward=20, got %d", sess.TotalReward)
	}
	if sess.MaxSingleWin != 20 {
		t.Errorf("expected MaxSingleWin=20, got %d", sess.MaxSingleWin)
	}
	if sess.TargetKillDist["T001"] != 1 {
		t.Errorf("expected TargetKillDist[T001]=1, got %d", sess.TargetKillDist["T001"])
	}
}

func TestTrackBossAndBonus(t *testing.T) {
	tracker := newTracker("test-room", t.TempDir())
	defer tracker.Close()

	tracker.Track(EventBossSpawn, "system", map[string]interface{}{
		"instance_id": "boss-001",
	})
	tracker.Track(EventBossKill, "system", map[string]interface{}{})
	tracker.Track(EventBonusStart, "system", map[string]interface{}{})
	tracker.Track(EventBonusStart, "system", map[string]interface{}{})

	stats := tracker.GetRoomStats()
	if stats.BossSpawnCount != 1 {
		t.Errorf("expected BossSpawnCount=1, got %d", stats.BossSpawnCount)
	}
	if stats.BossKillCount != 1 {
		t.Errorf("expected BossKillCount=1, got %d", stats.BossKillCount)
	}
	if stats.BonusCount != 2 {
		t.Errorf("expected BonusCount=2, got %d", stats.BonusCount)
	}
}

func TestRTPCalculation(t *testing.T) {
	tracker := newTracker("test-room", t.TempDir())
	defer tracker.Close()

	tracker.Track(EventPlayerJoin, "player-1", map[string]interface{}{})

	// 投入 100，獲得 95（RTP = 95%）
	for i := 0; i < 10; i++ {
		tracker.Track(EventAttack, "player-1", map[string]interface{}{
			"bet_level": 1,
			"bet_cost":  10,
		})
	}
	tracker.Track(EventReward, "player-1", map[string]interface{}{
		"amount": 95,
	})

	stats := tracker.GetRoomStats()
	if stats.TotalBet != 100 {
		t.Errorf("expected TotalBet=100, got %d", stats.TotalBet)
	}
	if stats.TotalReward != 95 {
		t.Errorf("expected TotalReward=95, got %d", stats.TotalReward)
	}
	expectedRTP := 0.95
	if stats.OverallRTP < expectedRTP-0.01 || stats.OverallRTP > expectedRTP+0.01 {
		t.Errorf("expected OverallRTP≈%.2f, got %.2f", expectedRTP, stats.OverallRTP)
	}
}

func TestJSONLLogOutput(t *testing.T) {
	tmpDir := t.TempDir()
	tracker := newTracker("test-room", tmpDir)

	tracker.Track(EventPlayerJoin, "player-1", map[string]interface{}{
		"room_id": "test-room",
	})
	tracker.Track(EventAttack, "player-1", map[string]interface{}{
		"bet_level": 1,
		"bet_cost":  5,
	})
	tracker.Close()

	// 確認日誌檔案存在且有內容
	today := time.Now().Format("2006-01-02")
	logPath := tmpDir + "/events-" + today + ".jsonl"
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("log file should exist: %v", err)
	}
	if info.Size() == 0 {
		t.Error("log file should not be empty")
	}
}

func TestNilTrackerSafe(t *testing.T) {
	// 確認 nil tracker 不會 panic
	var tracker *Tracker
	tracker.Track(EventAttack, "player-1", nil)
	stats := tracker.GetRoomStats()
	if stats.TotalAttacks != 0 {
		t.Error("nil tracker should return zero stats")
	}
	sess := tracker.GetSessionStats("player-1")
	if sess != nil {
		t.Error("nil tracker should return nil session")
	}
	tracker.Close() // 不應該 panic
}

func TestConcurrentTracking(t *testing.T) {
	// 並發安全測試
	tracker := newTracker("test-room", t.TempDir())
	defer tracker.Close()

	tracker.Track(EventPlayerJoin, "player-1", map[string]interface{}{})

	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				tracker.Track(EventAttack, "player-1", map[string]interface{}{
					"bet_level": 1,
					"bet_cost":  1,
				})
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}

	stats := tracker.GetRoomStats()
	if stats.TotalAttacks != 1000 {
		t.Errorf("expected TotalAttacks=1000 after concurrent writes, got %d", stats.TotalAttacks)
	}
}
