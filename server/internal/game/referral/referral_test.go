package referral

import (
	"testing"
)

func TestGetOrCreateCode(t *testing.T) {
	m := NewManager()
	code1 := m.GetOrCreateCode("player1")
	if len(code1) != 6 {
		t.Errorf("expected 6-char code, got %q", code1)
	}
	// 再次取得應該是同一個碼
	code2 := m.GetOrCreateCode("player1")
	if code1 != code2 {
		t.Errorf("expected same code, got %q vs %q", code1, code2)
	}
}

func TestGetOrCreateCode_Unique(t *testing.T) {
	m := NewManager()
	code1 := m.GetOrCreateCode("player1")
	code2 := m.GetOrCreateCode("player2")
	if code1 == code2 {
		t.Error("expected different codes for different players")
	}
}

func TestUseCode_Success(t *testing.T) {
	m := NewManager()
	code := m.GetOrCreateCode("referrer")

	referrerID, err := m.UseCode("referee", code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if referrerID != "referrer" {
		t.Errorf("expected referrerID=referrer, got %s", referrerID)
	}

	// 確認推薦人統計更新
	info := m.GetInfo("referrer")
	if info.ReferralCount != 1 {
		t.Errorf("expected ReferralCount=1, got %d", info.ReferralCount)
	}
	if info.TotalReward != ReferrerReward {
		t.Errorf("expected TotalReward=%d, got %d", ReferrerReward, info.TotalReward)
	}

	// 確認被推薦人統計更新
	refereeInfo := m.GetInfo("referee")
	if refereeInfo.UsedCode != code {
		t.Errorf("expected UsedCode=%s, got %s", code, refereeInfo.UsedCode)
	}
	if refereeInfo.TotalReward != RefereeReward {
		t.Errorf("expected TotalReward=%d, got %d", RefereeReward, refereeInfo.TotalReward)
	}
}

func TestUseCode_AlreadyUsed(t *testing.T) {
	m := NewManager()
	code1 := m.GetOrCreateCode("referrer1")
	code2 := m.GetOrCreateCode("referrer2")

	m.UseCode("referee", code1)
	_, err := m.UseCode("referee", code2)
	if err == nil {
		t.Error("expected error when using code twice")
	}
}

func TestUseCode_InvalidCode(t *testing.T) {
	m := NewManager()
	_, err := m.UseCode("player", "INVALID")
	if err == nil {
		t.Error("expected error for invalid code")
	}
}

func TestUseCode_SelfReferral(t *testing.T) {
	m := NewManager()
	code := m.GetOrCreateCode("player1")
	_, err := m.UseCode("player1", code)
	if err == nil {
		t.Error("expected error for self-referral")
	}
}

func TestUseCode_MaxReferrals(t *testing.T) {
	m := NewManager()
	referrerCode := m.GetOrCreateCode("referrer")

	// 達到上限
	for i := 0; i < MaxReferrals; i++ {
		refereeID := "referee" + string(rune('A'+i))
		_, err := m.UseCode(refereeID, referrerCode)
		if err != nil {
			t.Fatalf("unexpected error at referral %d: %v", i, err)
		}
	}

	// 超過上限
	_, err := m.UseCode("refereeZ", referrerCode)
	if err == nil {
		t.Error("expected error when exceeding max referrals")
	}
}

func TestGetInfo_AutoCreateCode(t *testing.T) {
	m := NewManager()
	info := m.GetInfo("newplayer")
	if info.MyCode == "" {
		t.Error("expected auto-created code for new player")
	}
	if len(info.MyCode) != 6 {
		t.Errorf("expected 6-char code, got %q", info.MyCode)
	}
}

func TestGetStats(t *testing.T) {
	m := NewManager()
	m.GetOrCreateCode("p1")
	m.GetOrCreateCode("p2")
	code := m.GetOrCreateCode("p3")
	m.UseCode("p4", code)

	totalCodes, totalReferrals := m.GetStats()
	if totalCodes != 3 {
		t.Errorf("expected 3 codes, got %d", totalCodes)
	}
	if totalReferrals != 1 {
		t.Errorf("expected 1 referral, got %d", totalReferrals)
	}
}

func TestGetRecentRecords(t *testing.T) {
	m := NewManager()
	code := m.GetOrCreateCode("referrer")
	m.UseCode("referee1", code)

	records := m.GetRecentRecords()
	if len(records) != 1 {
		t.Errorf("expected 1 record, got %d", len(records))
	}
	if records[0].ReferrerID != "referrer" {
		t.Errorf("expected referrerID=referrer, got %s", records[0].ReferrerID)
	}
}
