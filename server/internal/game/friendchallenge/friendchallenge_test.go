// friendchallenge_test.go — 好友挑戰系統測試（DAY-102）
package friendchallenge

import (
	"testing"
	"time"
)

func TestCreateChallenge_Success(t *testing.T) {
	m := New()
	c, errCode := m.CreateChallenge("p1", "Player1", "p2", "Player2")
	if errCode != "" {
		t.Fatalf("expected no error, got %s", errCode)
	}
	if c == nil {
		t.Fatal("expected challenge, got nil")
	}
	if c.Status != StatusPending {
		t.Errorf("expected pending, got %s", c.Status)
	}
	if c.Stake != ChallengeStake {
		t.Errorf("expected stake %d, got %d", ChallengeStake, c.Stake)
	}
}

func TestCreateChallenge_SelfChallenge(t *testing.T) {
	m := New()
	_, errCode := m.CreateChallenge("p1", "Player1", "p1", "Player1")
	if errCode != "self_challenge" {
		t.Errorf("expected self_challenge, got %s", errCode)
	}
}

func TestCreateChallenge_AlreadyInChallenge(t *testing.T) {
	m := New()
	c, _ := m.CreateChallenge("p1", "Player1", "p2", "Player2")
	m.AcceptChallenge(c.ID, "p2")

	_, errCode := m.CreateChallenge("p1", "Player1", "p3", "Player3")
	if errCode != "already_in_challenge" {
		t.Errorf("expected already_in_challenge, got %s", errCode)
	}
}

func TestAcceptChallenge_Success(t *testing.T) {
	m := New()
	c, _ := m.CreateChallenge("p1", "Player1", "p2", "Player2")
	accepted, errCode := m.AcceptChallenge(c.ID, "p2")
	if errCode != "" {
		t.Fatalf("expected no error, got %s", errCode)
	}
	if accepted.Status != StatusActive {
		t.Errorf("expected active, got %s", accepted.Status)
	}
	if accepted.EndAt.IsZero() {
		t.Error("expected EndAt to be set")
	}
}

func TestAcceptChallenge_WrongPlayer(t *testing.T) {
	m := New()
	c, _ := m.CreateChallenge("p1", "Player1", "p2", "Player2")
	_, errCode := m.AcceptChallenge(c.ID, "p3")
	if errCode != "not_your_challenge" {
		t.Errorf("expected not_your_challenge, got %s", errCode)
	}
}

func TestDeclineChallenge(t *testing.T) {
	m := New()
	c, _ := m.CreateChallenge("p1", "Player1", "p2", "Player2")
	ok := m.DeclineChallenge(c.ID, "p2")
	if !ok {
		t.Error("expected decline to succeed")
	}
	got := m.GetChallengeByID(c.ID)
	if got.Status != StatusDeclined {
		t.Errorf("expected declined, got %s", got.Status)
	}
}

func TestAddScore(t *testing.T) {
	m := New()
	c, _ := m.CreateChallenge("p1", "Player1", "p2", "Player2")
	m.AcceptChallenge(c.ID, "p2")

	m.AddScore("p1", 100)
	m.AddScore("p2", 200)
	m.AddScore("p1", 50)

	got := m.GetChallengeByID(c.ID)
	if got.ChallengerScore != 150 {
		t.Errorf("expected challenger score 150, got %d", got.ChallengerScore)
	}
	if got.ChallengedScore != 200 {
		t.Errorf("expected challenged score 200, got %d", got.ChallengedScore)
	}
}

func TestCheckAndFinish_ChallengerWins(t *testing.T) {
	m := New()
	c, _ := m.CreateChallenge("p1", "Player1", "p2", "Player2")
	accepted, _ := m.AcceptChallenge(c.ID, "p2")
	// 強制設定結束時間為過去
	m.mu.Lock()
	accepted.EndAt = time.Now().Add(-1 * time.Second)
	accepted.ChallengerScore = 500
	accepted.ChallengedScore = 300
	m.mu.Unlock()

	finished := m.CheckAndFinish()
	if len(finished) != 1 {
		t.Fatalf("expected 1 finished challenge, got %d", len(finished))
	}
	if finished[0].WinnerID != "p1" {
		t.Errorf("expected p1 to win, got %s", finished[0].WinnerID)
	}
	if finished[0].Prize != ChallengeStake*2 {
		t.Errorf("expected prize %d, got %d", ChallengeStake*2, finished[0].Prize)
	}
}

func TestCheckAndFinish_Draw(t *testing.T) {
	m := New()
	c, _ := m.CreateChallenge("p1", "Player1", "p2", "Player2")
	accepted, _ := m.AcceptChallenge(c.ID, "p2")
	m.mu.Lock()
	accepted.EndAt = time.Now().Add(-1 * time.Second)
	accepted.ChallengerScore = 300
	accepted.ChallengedScore = 300
	m.mu.Unlock()

	finished := m.CheckAndFinish()
	if len(finished) != 1 {
		t.Fatalf("expected 1 finished challenge, got %d", len(finished))
	}
	if finished[0].WinnerID != "" {
		t.Errorf("expected draw (no winner), got %s", finished[0].WinnerID)
	}
	if finished[0].Prize != ChallengeStake {
		t.Errorf("expected prize %d (refund), got %d", ChallengeStake, finished[0].Prize)
	}
}

func TestForceFinish_Disconnect(t *testing.T) {
	m := New()
	c, _ := m.CreateChallenge("p1", "Player1", "p2", "Player2")
	m.AcceptChallenge(c.ID, "p2")

	// p1 離線，p2 應該勝
	finished := m.ForceFinish("p1")
	if finished == nil {
		t.Fatal("expected finished challenge")
	}
	if finished.WinnerID != "p2" {
		t.Errorf("expected p2 to win, got %s", finished.WinnerID)
	}
}

func TestIsInChallenge(t *testing.T) {
	m := New()
	c, _ := m.CreateChallenge("p1", "Player1", "p2", "Player2")
	m.AcceptChallenge(c.ID, "p2")

	if !m.IsInChallenge("p1") {
		t.Error("expected p1 to be in challenge")
	}
	if !m.IsInChallenge("p2") {
		t.Error("expected p2 to be in challenge")
	}
	if m.IsInChallenge("p3") {
		t.Error("expected p3 to not be in challenge")
	}
}

func TestTimeRemaining(t *testing.T) {
	m := New()
	c, _ := m.CreateChallenge("p1", "Player1", "p2", "Player2")
	accepted, _ := m.AcceptChallenge(c.ID, "p2")

	remaining := accepted.TimeRemaining()
	if remaining <= 0 || remaining > int(ChallengeDuration.Seconds()) {
		t.Errorf("unexpected time remaining: %d", remaining)
	}
}
