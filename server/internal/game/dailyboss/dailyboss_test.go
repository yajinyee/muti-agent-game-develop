package dailyboss

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	snap := m.GetSnapshot()
	if snap == nil {
		t.Fatal("snapshot should not be nil")
	}
	if snap.Status != BossStatusActive {
		t.Errorf("expected Active, got %s", snap.Status)
	}
	if snap.CurrentHP != snap.MaxHP {
		t.Errorf("HP should be full at start: %d/%d", snap.CurrentHP, snap.MaxHP)
	}
	if snap.MaxHP <= 0 {
		t.Errorf("MaxHP should be > 0, got %d", snap.MaxHP)
	}
}

func TestAddDamage_Normal(t *testing.T) {
	m := New()
	snap := m.GetSnapshot()
	initialHP := snap.MaxHP

	defeated, reward := m.AddDamage("player1", "玩家一", 100)
	if defeated {
		t.Error("should not be defeated with small damage")
	}
	if reward != 0 {
		t.Errorf("reward should be 0 before defeat, got %d", reward)
	}

	snap = m.GetSnapshot()
	if snap.CurrentHP != initialHP-100 {
		t.Errorf("expected HP %d, got %d", initialHP-100, snap.CurrentHP)
	}
}

func TestAddDamage_Defeat(t *testing.T) {
	m := New()
	snap := m.GetSnapshot()
	maxHP := snap.MaxHP

	// 一次打死
	defeated, reward := m.AddDamage("player1", "玩家一", maxHP)
	if !defeated {
		t.Error("should be defeated")
	}
	if reward <= 0 {
		t.Errorf("reward should be > 0, got %d", reward)
	}

	snap = m.GetSnapshot()
	if snap.Status != BossStatusDefeated {
		t.Errorf("expected Defeated, got %s", snap.Status)
	}
	if snap.CurrentHP != 0 {
		t.Errorf("HP should be 0, got %d", snap.CurrentHP)
	}
}

func TestAddDamage_MultipleContributors(t *testing.T) {
	m := New()
	snap := m.GetSnapshot()
	maxHP := snap.MaxHP

	// 三個玩家各打 1/3
	dmg := maxHP / 3
	m.AddDamage("player1", "玩家一", dmg)
	m.AddDamage("player2", "玩家二", dmg)
	defeated, _ := m.AddDamage("player3", "玩家三", maxHP) // 最後一擊

	if !defeated {
		t.Error("should be defeated")
	}

	// 確認貢獻記錄
	c1 := m.GetPlayerContribution("player1")
	c2 := m.GetPlayerContribution("player2")
	c3 := m.GetPlayerContribution("player3")

	if c1 == nil || c2 == nil || c3 == nil {
		t.Fatal("all contributions should be recorded")
	}

	// 確認獎勵分配（player3 貢獻最多，獎勵最多）
	if c3.Reward <= c1.Reward {
		t.Errorf("player3 should have more reward than player1: %d vs %d", c3.Reward, c1.Reward)
	}
}

func TestGetTopContributors(t *testing.T) {
	m := New()
	snap := m.GetSnapshot()
	maxHP := snap.MaxHP

	m.AddDamage("player1", "玩家一", maxHP/4)
	m.AddDamage("player2", "玩家二", maxHP/2)
	m.AddDamage("player3", "玩家三", maxHP/8)

	top := m.GetTopContributors(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 contributors, got %d", len(top))
	}
	if top[0].PlayerID != "player2" {
		t.Errorf("expected player2 first, got %s", top[0].PlayerID)
	}
	if top[1].PlayerID != "player1" {
		t.Errorf("expected player1 second, got %s", top[1].PlayerID)
	}
}

func TestGetHPPercent(t *testing.T) {
	m := New()
	snap := m.GetSnapshot()
	maxHP := snap.MaxHP

	pct := m.GetHPPercent()
	if pct != 1.0 {
		t.Errorf("expected 1.0, got %f", pct)
	}

	m.AddDamage("player1", "玩家一", maxHP/2)
	pct = m.GetHPPercent()
	if pct < 0.49 || pct > 0.51 {
		t.Errorf("expected ~0.5, got %f", pct)
	}
}

func TestDifficultyMod_ConsecutiveFails(t *testing.T) {
	m := New()
	m.consecutiveFails = 2 // 模擬連續 2 天未擊殺

	// 重新生成 BOSS
	m.spawnTodayBoss()
	snap := m.GetSnapshot()

	// 難度應該降低（HP 應該比基礎低）
	bossType := DailyBossTypes[time.Now().YearDay()%len(DailyBossTypes)]
	expectedHP := int(float64(bossType.BaseHP) * 0.6) // 1.0 - 2*0.2 = 0.6
	if snap.MaxHP != expectedHP {
		t.Errorf("expected HP %d (60%% of base), got %d", expectedHP, snap.MaxHP)
	}
}

func TestGetDateID(t *testing.T) {
	t1 := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	dateID := getDateID(t1)
	if dateID != "2026-05-20" {
		t.Errorf("expected 2026-05-20, got %s", dateID)
	}
}

func TestGetDayRange(t *testing.T) {
	t1 := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	start, end := getDayRange(t1)

	loc := time.FixedZone("UTC+8", 8*60*60)
	startLocal := start.In(loc)
	endLocal := end.In(loc)

	// start 應該是 2026-05-20 00:00:00 UTC+8
	if startLocal.Hour() != 0 || startLocal.Minute() != 0 {
		t.Errorf("start should be 00:00, got %02d:%02d", startLocal.Hour(), startLocal.Minute())
	}

	// end 應該是 2026-05-21 00:00:00 UTC+8
	if endLocal.Day() != 21 {
		t.Errorf("end should be day 21, got %d", endLocal.Day())
	}

	if !end.After(start) {
		t.Error("end should be after start")
	}
}

func TestAddDamage_AfterDefeated(t *testing.T) {
	m := New()
	snap := m.GetSnapshot()
	maxHP := snap.MaxHP

	// 擊殺
	m.AddDamage("player1", "玩家一", maxHP)

	// 再次攻擊，不應有效果
	defeated, reward := m.AddDamage("player2", "玩家二", 100)
	if defeated {
		t.Error("should not be defeated again")
	}
	if reward != 0 {
		t.Errorf("reward should be 0, got %d", reward)
	}
}

func TestFormatHP(t *testing.T) {
	boss := &DailyBoss{
		CurrentHP: 7500,
		MaxHP:     10000,
	}
	formatted := boss.FormatHP()
	if formatted != "7500/10000" {
		t.Errorf("expected '7500/10000', got '%s'", formatted)
	}
}
