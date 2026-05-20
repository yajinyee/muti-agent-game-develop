package stats

import (
	"testing"
)

func TestNewPlayerStats(t *testing.T) {
	s := NewPlayerStats()
	if s.TotalSessions != 0 {
		t.Errorf("TotalSessions = %d, want 0", s.TotalSessions)
	}
	if s.FirstPlayAt.IsZero() {
		t.Error("FirstPlayAt should not be zero")
	}
}

func TestStartEndSession(t *testing.T) {
	s := NewPlayerStats()
	s.StartSession()
	if s.TotalSessions != 1 {
		t.Errorf("TotalSessions = %d, want 1", s.TotalSessions)
	}
	// 模擬一段時間
	s.EndSession()
	// TotalPlayTime 應該 >= 0
	if s.TotalPlayTime < 0 {
		t.Error("TotalPlayTime should be >= 0")
	}
}

func TestRecordKill(t *testing.T) {
	s := NewPlayerStats()
	s.RecordKill(5.0, 500)
	s.RecordKill(20.0, 2000)
	s.RecordKill(3.0, 300)

	if s.TotalKills != 3 {
		t.Errorf("TotalKills = %d, want 3", s.TotalKills)
	}
	if s.TotalReward != 2800 {
		t.Errorf("TotalReward = %d, want 2800", s.TotalReward)
	}
	if s.BestMultiplier != 20.0 {
		t.Errorf("BestMultiplier = %.1f, want 20.0", s.BestMultiplier)
	}
}

func TestRecordStreak(t *testing.T) {
	s := NewPlayerStats()
	s.RecordStreak(5)
	s.RecordStreak(12)
	s.RecordStreak(8)

	if s.BestStreak != 12 {
		t.Errorf("BestStreak = %d, want 12", s.BestStreak)
	}
}

func TestRecordJackpot(t *testing.T) {
	s := NewPlayerStats()
	s.RecordJackpot("mini", 300)
	s.RecordJackpot("minor", 1000)
	s.RecordJackpot("grand", 15000)

	if s.JackpotWins != 3 {
		t.Errorf("JackpotWins = %d, want 3", s.JackpotWins)
	}
	if s.JackpotMiniWins != 1 {
		t.Errorf("JackpotMiniWins = %d, want 1", s.JackpotMiniWins)
	}
	if s.JackpotGrandWins != 1 {
		t.Errorf("JackpotGrandWins = %d, want 1", s.JackpotGrandWins)
	}
	if s.TotalJackpotPayout != 16300 {
		t.Errorf("TotalJackpotPayout = %d, want 16300", s.TotalJackpotPayout)
	}
}

func TestGetHitRate(t *testing.T) {
	s := NewPlayerStats()
	// 0 shots
	if s.GetHitRate() != 0.0 {
		t.Error("HitRate should be 0 with no shots")
	}

	s.RecordKill(2.0, 200) // hit
	s.RecordMiss()         // miss
	s.RecordMiss()         // miss
	s.RecordKill(5.0, 500) // hit

	rate := s.GetHitRate()
	expected := 0.5 // 2 hits / 4 total
	if rate != expected {
		t.Errorf("HitRate = %.2f, want %.2f", rate, expected)
	}
}

func TestGetRTP(t *testing.T) {
	s := NewPlayerStats()
	// 0 bet
	if s.GetRTP() != 0.0 {
		t.Error("RTP should be 0 with no bet")
	}

	s.RecordShot(100)
	s.RecordShot(100)
	s.RecordKill(2.0, 180) // 180 reward for 200 bet = 90% RTP

	rtp := s.GetRTP()
	expected := 0.9
	if rtp != expected {
		t.Errorf("RTP = %.2f, want %.2f", rtp, expected)
	}
}

func TestSnapshot(t *testing.T) {
	s := NewPlayerStats()
	s.StartSession()
	s.RecordShot(100)
	s.RecordKill(10.0, 1000)
	s.RecordBonus(5000)
	s.RecordBossKill()
	s.RecordJackpot("mini", 300)
	s.RecordStreak(8)
	s.UpdateMaxCoins(50000)

	snap := s.Snapshot()

	if snap.TotalShots != 1 {
		t.Errorf("Snapshot TotalShots = %d, want 1", snap.TotalShots)
	}
	if snap.TotalKills != 1 {
		t.Errorf("Snapshot TotalKills = %d, want 1", snap.TotalKills)
	}
	if snap.BestMultiplier != 10.0 {
		t.Errorf("Snapshot BestMultiplier = %.1f, want 10.0", snap.BestMultiplier)
	}
	if snap.TotalBonuses != 1 {
		t.Errorf("Snapshot TotalBonuses = %d, want 1", snap.TotalBonuses)
	}
	if snap.TotalBossKills != 1 {
		t.Errorf("Snapshot TotalBossKills = %d, want 1", snap.TotalBossKills)
	}
	if snap.JackpotWins != 1 {
		t.Errorf("Snapshot JackpotWins = %d, want 1", snap.JackpotWins)
	}
	if snap.BestStreak != 8 {
		t.Errorf("Snapshot BestStreak = %d, want 8", snap.BestStreak)
	}
	if snap.MaxCoins != 50000 {
		t.Errorf("Snapshot MaxCoins = %d, want 50000", snap.MaxCoins)
	}
	// 當前 Session 時間應該 >= 0
	if snap.TotalPlayTimeSec < 0 {
		t.Error("TotalPlayTimeSec should be >= 0")
	}
}

func TestRecordBonus(t *testing.T) {
	s := NewPlayerStats()
	s.RecordBonus(1000)
	s.RecordBonus(5000)
	s.RecordBonus(2000)

	if s.TotalBonuses != 3 {
		t.Errorf("TotalBonuses = %d, want 3", s.TotalBonuses)
	}
	if s.BestBonusReward != 5000 {
		t.Errorf("BestBonusReward = %d, want 5000", s.BestBonusReward)
	}
}

func TestUpdateMaxCoins(t *testing.T) {
	s := NewPlayerStats()
	s.UpdateMaxCoins(1000)
	s.UpdateMaxCoins(5000)
	s.UpdateMaxCoins(3000) // 不應更新

	if s.MaxCoins != 5000 {
		t.Errorf("MaxCoins = %d, want 5000", s.MaxCoins)
	}
}
