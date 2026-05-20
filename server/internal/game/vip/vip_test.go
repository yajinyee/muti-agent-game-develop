package vip

import (
	"testing"
	"time"
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
		t.Errorf("expected player_id=player1, got %s", data.PlayerID)
	}
	if data.VIPLevel != 0 {
		t.Errorf("expected vip_level=0, got %d", data.VIPLevel)
	}
	if data.TotalSpend != 0 {
		t.Errorf("expected total_spend=0, got %d", data.TotalSpend)
	}

	// 再次取得應該是同一個
	data2 := m.GetOrCreate("player1")
	if data2.PlayerID != data.PlayerID {
		t.Error("GetOrCreate should return same data for same player")
	}
}

func TestAddSpend_NoLevelUp(t *testing.T) {
	m := New()
	level, result := m.AddSpend("player1", 5000)
	if level != 0 {
		t.Errorf("expected level=0, got %d", level)
	}
	if result != nil {
		t.Error("expected no level up result")
	}
}

func TestAddSpend_LevelUp_Bronze(t *testing.T) {
	m := New()
	level, result := m.AddSpend("player1", 10000)
	if level != 1 {
		t.Errorf("expected level=1 (Bronze), got %d", level)
	}
	if result == nil {
		t.Fatal("expected level up result")
	}
	if result.NewLevel != 1 {
		t.Errorf("expected new_level=1, got %d", result.NewLevel)
	}
	if result.TierName != "青銅會員" {
		t.Errorf("expected tier_name=青銅會員, got %s", result.TierName)
	}
	if result.TitleID != "vip_bronze" {
		t.Errorf("expected title_id=vip_bronze, got %s", result.TitleID)
	}
}

func TestAddSpend_LevelUp_Silver(t *testing.T) {
	m := New()
	m.AddSpend("player1", 50000)
	level, result := m.AddSpend("player1", 0)
	if level != 2 {
		t.Errorf("expected level=2 (Silver), got %d", level)
	}
	// 第二次 AddSpend 不應該再觸發升級（已是等級2）
	if result != nil {
		t.Error("should not level up again on same level")
	}
}

func TestAddSpend_MultiLevelUp(t *testing.T) {
	m := New()
	// 一次消費直接達到黃金等級
	level, result := m.AddSpend("player1", 200000)
	if level != 3 {
		t.Errorf("expected level=3 (Gold), got %d", level)
	}
	if result == nil {
		t.Fatal("expected level up result")
	}
	if result.NewLevel != 3 {
		t.Errorf("expected new_level=3, got %d", result.NewLevel)
	}
}

func TestGetCashback_NoVIP(t *testing.T) {
	m := New()
	cashback := m.GetCashback("player1", 1000)
	if cashback != 0 {
		t.Errorf("expected cashback=0 for non-VIP, got %d", cashback)
	}
}

func TestGetCashback_Bronze(t *testing.T) {
	m := New()
	m.AddSpend("player1", 10000)
	// Bronze = 1% cashback
	cashback := m.GetCashback("player1", 1000)
	if cashback != 10 {
		t.Errorf("expected cashback=10 (1%% of 1000), got %d", cashback)
	}
}

func TestGetCashback_Gold(t *testing.T) {
	m := New()
	m.AddSpend("player1", 200000)
	// Gold = 3% cashback
	cashback := m.GetCashback("player1", 1000)
	if cashback != 30 {
		t.Errorf("expected cashback=30 (3%% of 1000), got %d", cashback)
	}
}

func TestGetDailyBonusMult_NoVIP(t *testing.T) {
	m := New()
	mult := m.GetDailyBonusMult("player1")
	if mult != 1.0 {
		t.Errorf("expected mult=1.0 for non-VIP, got %f", mult)
	}
}

func TestGetDailyBonusMult_Bronze(t *testing.T) {
	m := New()
	m.AddSpend("player1", 10000)
	mult := m.GetDailyBonusMult("player1")
	if mult != 1.1 {
		t.Errorf("expected mult=1.1 for Bronze, got %f", mult)
	}
}

func TestClaimWeeklyBonus_NoVIP(t *testing.T) {
	m := New()
	result := m.ClaimWeeklyBonus("player1")
	if result != nil {
		t.Error("expected nil for non-VIP player")
	}
}

func TestClaimWeeklyBonus_Bronze(t *testing.T) {
	m := New()
	m.AddSpend("player1", 10000)
	result := m.ClaimWeeklyBonus("player1")
	if result == nil {
		t.Fatal("expected weekly bonus result")
	}
	if result.Coins != 500 {
		t.Errorf("expected coins=500 for Bronze, got %d", result.Coins)
	}
	if result.VIPLevel != 1 {
		t.Errorf("expected vip_level=1, got %d", result.VIPLevel)
	}
}

func TestClaimWeeklyBonus_CooldownNotExpired(t *testing.T) {
	m := New()
	m.AddSpend("player1", 10000)
	// 第一次領取
	m.ClaimWeeklyBonus("player1")
	// 立刻再領取應該失敗
	result := m.ClaimWeeklyBonus("player1")
	if result != nil {
		t.Error("should not be able to claim weekly bonus twice within 7 days")
	}
}

func TestGetSnapshot_NoVIP(t *testing.T) {
	m := New()
	snap := m.GetSnapshot("player1")
	if snap.VIPLevel != 0 {
		t.Errorf("expected vip_level=0, got %d", snap.VIPLevel)
	}
	if snap.SpendToNext != VIPTiers[0].SpendRequired {
		t.Errorf("expected spend_to_next=%d, got %d", VIPTiers[0].SpendRequired, snap.SpendToNext)
	}
	if snap.Progress != 0.0 {
		t.Errorf("expected progress=0.0, got %f", snap.Progress)
	}
}

func TestGetSnapshot_Bronze(t *testing.T) {
	m := New()
	m.AddSpend("player1", 10000)
	snap := m.GetSnapshot("player1")
	if snap.VIPLevel != 1 {
		t.Errorf("expected vip_level=1, got %d", snap.VIPLevel)
	}
	if snap.TierName != "青銅會員" {
		t.Errorf("expected tier_name=青銅會員, got %s", snap.TierName)
	}
	if snap.CashbackRate != 0.01 {
		t.Errorf("expected cashback_rate=0.01, got %f", snap.CashbackRate)
	}
	if snap.CanClaimWeekly != true {
		t.Error("expected can_claim_weekly=true for new VIP")
	}
}

func TestGetSnapshot_Progress(t *testing.T) {
	m := New()
	// 消費 25000（Bronze 門檻 10000，Silver 門檻 50000）
	// 進度 = (25000 - 10000) / (50000 - 10000) = 15000/40000 = 0.375
	m.AddSpend("player1", 25000)
	snap := m.GetSnapshot("player1")
	if snap.VIPLevel != 1 {
		t.Errorf("expected vip_level=1, got %d", snap.VIPLevel)
	}
	expectedProgress := float64(25000-10000) / float64(50000-10000)
	if snap.Progress < expectedProgress-0.01 || snap.Progress > expectedProgress+0.01 {
		t.Errorf("expected progress≈%f, got %f", expectedProgress, snap.Progress)
	}
}

func TestGetVIPLevel(t *testing.T) {
	m := New()
	if m.GetVIPLevel("player1") != 0 {
		t.Error("expected level=0 for new player")
	}
	m.AddSpend("player1", 10000)
	if m.GetVIPLevel("player1") != 1 {
		t.Error("expected level=1 after spending 10000")
	}
}

func TestVIPTiers_Consistency(t *testing.T) {
	// 確認等級定義一致性
	if len(VIPTiers) != 5 {
		t.Errorf("expected 5 VIP tiers, got %d", len(VIPTiers))
	}
	for i, tier := range VIPTiers {
		if tier.Level != i+1 {
			t.Errorf("tier[%d].Level should be %d, got %d", i, i+1, tier.Level)
		}
		if tier.SpendRequired <= 0 {
			t.Errorf("tier[%d].SpendRequired should be > 0", i)
		}
		if tier.CashbackRate <= 0 {
			t.Errorf("tier[%d].CashbackRate should be > 0", i)
		}
		if tier.DailyBonusMult < 1.0 {
			t.Errorf("tier[%d].DailyBonusMult should be >= 1.0", i)
		}
		if i > 0 && tier.SpendRequired <= VIPTiers[i-1].SpendRequired {
			t.Errorf("tier[%d].SpendRequired should be > tier[%d].SpendRequired", i, i-1)
		}
	}
}

func TestWeeklyBonusCooldown_Expired(t *testing.T) {
	m := New()
	m.AddSpend("player1", 10000)

	// 手動設定上次領取時間為 8 天前
	m.mu.Lock()
	m.players["player1"].LastWeeklyAt = time.Now().Add(-8 * 24 * time.Hour)
	m.mu.Unlock()

	result := m.ClaimWeeklyBonus("player1")
	if result == nil {
		t.Error("should be able to claim weekly bonus after 7 days")
	}
}
