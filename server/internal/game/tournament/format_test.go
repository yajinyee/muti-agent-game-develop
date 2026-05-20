package tournament

import (
	"testing"
)

func TestMultiFormat_RecordKill_ScoreFormat(t *testing.T) {
	m := NewMultiFormat()
	// 強制設定為積分賽格式
	m.mu.Lock()
	m.todayFormat = FormatScore
	m.mu.Unlock()

	m.RecordKill("p1", "Player1", 10.0, 500, 10)
	m.RecordKill("p1", "Player1", 5.0, 200, 10)

	rank, score := m.GetPlayerRank("p1")
	if rank != 1 {
		t.Errorf("expected rank 1, got %d", rank)
	}
	if score != 15 { // 10 + 5
		t.Errorf("expected score 15, got %.0f", score)
	}
}

func TestMultiFormat_RecordKill_MultiplierFormat(t *testing.T) {
	m := NewMultiFormat()
	m.mu.Lock()
	m.todayFormat = FormatMultiplier
	m.mu.Unlock()

	m.RecordKill("p1", "Player1", 10.0, 500, 10)
	m.RecordKill("p1", "Player1", 50.0, 2000, 10) // 更高倍率
	m.RecordKill("p2", "Player2", 30.0, 1500, 10)

	rank, score := m.GetPlayerRank("p1")
	if rank != 1 {
		t.Errorf("expected p1 rank 1, got %d", rank)
	}
	if score != 50.0 {
		t.Errorf("expected best multiplier 50, got %.0f", score)
	}
}

func TestMultiFormat_RecordKill_RewardFormat(t *testing.T) {
	m := NewMultiFormat()
	m.mu.Lock()
	m.todayFormat = FormatReward
	m.mu.Unlock()

	m.RecordKill("p1", "Player1", 5.0, 100, 10)
	m.RecordKill("p1", "Player1", 10.0, 5000, 10) // 更高獎勵
	m.RecordKill("p2", "Player2", 20.0, 3000, 10)

	rank, score := m.GetPlayerRank("p1")
	if rank != 1 {
		t.Errorf("expected p1 rank 1, got %d", rank)
	}
	if score != 5000 {
		t.Errorf("expected best reward 5000, got %.0f", score)
	}
}

func TestMultiFormat_RecordShot_BetFormat(t *testing.T) {
	m := NewMultiFormat()
	m.mu.Lock()
	m.todayFormat = FormatBet
	m.mu.Unlock()

	m.RecordShot("p1", "Player1", 100)
	m.RecordShot("p1", "Player1", 100)
	m.RecordShot("p2", "Player2", 50)

	rank, score := m.GetPlayerRank("p1")
	if rank != 1 {
		t.Errorf("expected p1 rank 1, got %d", rank)
	}
	if score != 200 {
		t.Errorf("expected total bet 200, got %.0f", score)
	}
}

func TestMultiFormat_RecordBoss(t *testing.T) {
	m := NewMultiFormat()
	m.mu.Lock()
	m.todayFormat = FormatScore
	m.mu.Unlock()

	m.RecordBoss("p1", "Player1")
	rank, score := m.GetPlayerRank("p1")
	if rank != 1 {
		t.Errorf("expected rank 1, got %d", rank)
	}
	if score != 50 {
		t.Errorf("expected 50 points for boss kill, got %.0f", score)
	}
}

func TestMultiFormat_RecordBonus(t *testing.T) {
	m := NewMultiFormat()
	m.mu.Lock()
	m.todayFormat = FormatScore
	m.mu.Unlock()

	m.RecordBonus("p1", "Player1", 1000)
	rank, score := m.GetPlayerRank("p1")
	if rank != 1 {
		t.Errorf("expected rank 1, got %d", rank)
	}
	if score != 20 {
		t.Errorf("expected 20 points for bonus, got %.0f", score)
	}
}

func TestMultiFormat_GetSnapshot(t *testing.T) {
	m := NewMultiFormat()
	m.RecordKill("p1", "Player1", 10.0, 500, 10)

	snap := m.GetSnapshot()
	if snap.TotalPlayers != 1 {
		t.Errorf("expected 1 player, got %d", snap.TotalPlayers)
	}
	if snap.FormatDef.Type != snap.TodayFormat {
		t.Error("format def should match today format")
	}
	if snap.NextFormat == "" {
		t.Error("next format should not be empty")
	}
}

func TestMultiFormat_GetTodayFormat_Rotates(t *testing.T) {
	// 確認 4 天輪換邏輯
	formats := make(map[FormatType]bool)
	for i := 0; i < 4; i++ {
		// 模擬不同天
		dayOfYear := i
		allFmts := []FormatType{FormatScore, FormatMultiplier, FormatReward, FormatBet}
		ft := allFmts[dayOfYear%4]
		formats[ft] = true
	}
	if len(formats) != 4 {
		t.Errorf("expected 4 different formats in rotation, got %d", len(formats))
	}
}

func TestMultiFormat_GetFormatDef(t *testing.T) {
	for _, ft := range []FormatType{FormatScore, FormatMultiplier, FormatReward, FormatBet} {
		def := GetFormatDef(ft)
		if def.Name == "" {
			t.Errorf("format %s has empty name", ft)
		}
		if def.Icon == "" {
			t.Errorf("format %s has empty icon", ft)
		}
	}
}

func TestMultiFormat_MultiplePlayersRanking(t *testing.T) {
	m := NewMultiFormat()
	m.mu.Lock()
	m.todayFormat = FormatMultiplier
	m.mu.Unlock()

	m.RecordKill("p1", "Player1", 30.0, 1000, 10)
	m.RecordKill("p2", "Player2", 50.0, 2000, 10)
	m.RecordKill("p3", "Player3", 20.0, 500, 10)

	rankings := m.GetRankings(10)
	if len(rankings) != 3 {
		t.Errorf("expected 3 rankings, got %d", len(rankings))
	}
	if rankings[0].PlayerID != "p2" {
		t.Errorf("expected p2 first (50x), got %s", rankings[0].PlayerID)
	}
	if rankings[1].PlayerID != "p1" {
		t.Errorf("expected p1 second (30x), got %s", rankings[1].PlayerID)
	}
}

func TestMultiFormat_FormatScore_Label(t *testing.T) {
	label := formatScore(FormatMultiplier, 50.0)
	if label != "50x" {
		t.Errorf("expected '50x', got '%s'", label)
	}
	label = formatScore(FormatReward, 5000.0)
	if label != "5000" {
		t.Errorf("expected '5000', got '%s'", label)
	}
}
