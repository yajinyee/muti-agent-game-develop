// challenge_test.go - DAY-085 隱藏挑戰系統測試
package challenge

import (
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("NewManager returned nil")
	}
}

func TestInitPlayer(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")
	snap := m.GetSnapshot("p1")
	if len(snap) == 0 {
		t.Fatal("expected non-empty snapshot after init")
	}
}

func TestTryUnlock_Basic(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")

	def := m.TryUnlock("p1", ChallengeBossFirst)
	if def == nil {
		t.Fatal("expected unlock to succeed")
	}
	if def.ID != ChallengeBossFirst {
		t.Errorf("expected %s, got %s", ChallengeBossFirst, def.ID)
	}

	// 重複解鎖應回傳 nil
	def2 := m.TryUnlock("p1", ChallengeBossFirst)
	if def2 != nil {
		t.Fatal("expected nil on duplicate unlock")
	}
}

func TestClaimReward(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")
	m.TryUnlock("p1", ChallengeBossFirst)

	reward := m.ClaimReward("p1", ChallengeBossFirst)
	if reward <= 0 {
		t.Errorf("expected positive reward, got %d", reward)
	}

	// 重複領取應回傳 0
	reward2 := m.ClaimReward("p1", ChallengeBossFirst)
	if reward2 != 0 {
		t.Errorf("expected 0 on duplicate claim, got %d", reward2)
	}
}

func TestRecordKill_SpeedChallenge3s(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")

	// 快速擊破 3 個目標
	unlocked := m.RecordKill("p1", "T001", 2.0)
	unlocked = append(unlocked, m.RecordKill("p1", "T002", 3.0)...)
	unlocked = append(unlocked, m.RecordKill("p1", "T003", 5.0)...)

	found := false
	for _, u := range unlocked {
		if u.ID == ChallengeSpeed3 {
			found = true
		}
	}
	if !found {
		t.Error("expected ChallengeSpeed3 to be unlocked after 3 kills in 3s")
	}
}

func TestRecordKill_MultChallenge(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")

	unlocked := m.RecordKill("p1", "T103", 50.0)
	found := false
	for _, u := range unlocked {
		if u.ID == ChallengeMult50 {
			found = true
		}
	}
	if !found {
		t.Error("expected ChallengeMult50 to be unlocked")
	}
}

func TestRecordKill_Mult100(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")

	unlocked := m.RecordKill("p1", "T103", 100.0)
	found50, found100 := false, false
	for _, u := range unlocked {
		if u.ID == ChallengeMult50 {
			found50 = true
		}
		if u.ID == ChallengeMult100 {
			found100 = true
		}
	}
	if !found50 || !found100 {
		t.Errorf("expected both mult50 and mult100, got 50=%v 100=%v", found50, found100)
	}
}

func TestRecordCoins(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")

	unlocked := m.RecordCoins("p1", 10000)
	found := false
	for _, u := range unlocked {
		if u.ID == ChallengeRich10k {
			found = true
		}
	}
	if !found {
		t.Error("expected ChallengeRich10k to be unlocked at 10000 coins")
	}
}

func TestRecordCoins_50k(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")

	unlocked := m.RecordCoins("p1", 50000)
	found10k, found50k := false, false
	for _, u := range unlocked {
		if u.ID == ChallengeRich10k {
			found10k = true
		}
		if u.ID == ChallengeRich50k {
			found50k = true
		}
	}
	if !found10k || !found50k {
		t.Errorf("expected both 10k and 50k, got 10k=%v 50k=%v", found10k, found50k)
	}
}

func TestRecordStreak(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")

	unlocked := m.RecordStreak("p1", 5)
	found := false
	for _, u := range unlocked {
		if u.ID == ChallengeStreak5 {
			found = true
		}
	}
	if !found {
		t.Error("expected ChallengeStreak5 to be unlocked at streak=5")
	}
}

func TestRecordStreak_20(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")

	unlocked := m.RecordStreak("p1", 20)
	found5, found10, found20 := false, false, false
	for _, u := range unlocked {
		switch u.ID {
		case ChallengeStreak5:
			found5 = true
		case ChallengeStreak10:
			found10 = true
		case ChallengeStreak20:
			found20 = true
		}
	}
	if !found5 || !found10 || !found20 {
		t.Errorf("expected all streak challenges, got 5=%v 10=%v 20=%v", found5, found10, found20)
	}
}

func TestRecordAllTypes(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")

	allTypes := []string{"T001", "T002", "T003"}
	m.RecordKill("p1", "T001", 2.0)
	m.RecordKill("p1", "T002", 3.0)

	// 還沒全部擊破
	result := m.RecordAllTypes("p1", allTypes)
	if result != nil {
		t.Error("expected nil before all types killed")
	}

	m.RecordKill("p1", "T003", 5.0)
	result = m.RecordAllTypes("p1", allTypes)
	if result == nil {
		t.Error("expected ChallengeAllTypes to unlock after all types killed")
	}
}

func TestGetSnapshot_HiddenChallenge(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")

	snap := m.GetSnapshot("p1")
	for _, s := range snap {
		if s.ID == string(ChallengeStreak20) {
			if !s.IsHidden {
				t.Error("ChallengeStreak20 should be hidden before unlock")
			}
		}
	}
}

func TestGetSnapshot_UnlockedNotHidden(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")
	m.TryUnlock("p1", ChallengeStreak20)

	snap := m.GetSnapshot("p1")
	for _, s := range snap {
		if s.ID == string(ChallengeStreak20) {
			if s.IsHidden {
				t.Error("ChallengeStreak20 should NOT be hidden after unlock")
			}
			if !s.Unlocked {
				t.Error("ChallengeStreak20 should be marked as unlocked")
			}
		}
	}
}

func TestRemovePlayer(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")
	m.RecordKill("p1", "T001", 2.0)
	m.RemovePlayer("p1")

	// 移除後 session 應被清除，重新 init 後 KillTimestamps 應為空
	m.InitPlayer("p1")
	m.mu.RLock()
	sess := m.sessions["p1"]
	m.mu.RUnlock()
	if len(sess.KillTimestamps) != 0 {
		t.Error("expected empty KillTimestamps after RemovePlayer + InitPlayer")
	}
}

func TestKillTimestampCleanup(t *testing.T) {
	m := NewManager()
	m.InitPlayer("p1")

	// 模擬舊的時間戳（超過 10 秒）
	m.mu.Lock()
	m.sessions["p1"].KillTimestamps = []time.Time{
		time.Now().Add(-15 * time.Second),
		time.Now().Add(-12 * time.Second),
	}
	m.mu.Unlock()

	// 新的擊破應清理舊時間戳
	m.RecordKill("p1", "T001", 2.0)

	m.mu.RLock()
	count := len(m.sessions["p1"].KillTimestamps)
	m.mu.RUnlock()

	if count != 1 {
		t.Errorf("expected 1 timestamp after cleanup, got %d", count)
	}
}
