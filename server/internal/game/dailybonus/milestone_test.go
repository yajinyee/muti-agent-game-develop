// Package dailybonus — 登入里程碑獎勵系統測試（DAY-107）
package dailybonus

import (
	"testing"
)

func TestCheckMilestone_Day3(t *testing.T) {
	m := CheckMilestone(3)
	if m == nil {
		t.Fatal("expected milestone at day 3, got nil")
	}
	if m.Days != 3 {
		t.Errorf("expected days=3, got %d", m.Days)
	}
	if m.Name != "初心者" {
		t.Errorf("expected name=初心者, got %s", m.Name)
	}
	// 應有金幣 + 普通寶箱
	if len(m.Rewards) != 2 {
		t.Errorf("expected 2 rewards, got %d", len(m.Rewards))
	}
}

func TestCheckMilestone_Day7(t *testing.T) {
	m := CheckMilestone(7)
	if m == nil {
		t.Fatal("expected milestone at day 7, got nil")
	}
	if m.Icon != "⚔️" {
		t.Errorf("expected icon=⚔️, got %s", m.Icon)
	}
	// 應有稀有寶箱
	hasRare := false
	for _, r := range m.Rewards {
		if r.Type == MilestoneRewardMysteryBox && r.Rarity == "rare" {
			hasRare = true
		}
	}
	if !hasRare {
		t.Error("expected rare mystery box reward at day 7")
	}
}

func TestCheckMilestone_Day14_HasTitle(t *testing.T) {
	m := CheckMilestone(14)
	if m == nil {
		t.Fatal("expected milestone at day 14, got nil")
	}
	hasTitle := false
	for _, r := range m.Rewards {
		if r.Type == MilestoneRewardTitle {
			hasTitle = true
			if r.TitleID != "streak_veteran" {
				t.Errorf("expected title_id=streak_veteran, got %s", r.TitleID)
			}
		}
	}
	if !hasTitle {
		t.Error("expected title reward at day 14")
	}
}

func TestCheckMilestone_Day30_HasLegendary(t *testing.T) {
	m := CheckMilestone(30)
	if m == nil {
		t.Fatal("expected milestone at day 30, got nil")
	}
	hasLegendary := false
	for _, r := range m.Rewards {
		if r.Type == MilestoneRewardMysteryBox && r.Rarity == "legendary" {
			hasLegendary = true
		}
	}
	if !hasLegendary {
		t.Error("expected legendary mystery box at day 30")
	}
}

func TestCheckMilestone_NoMilestone(t *testing.T) {
	// 非里程碑天數應回傳 nil
	for _, day := range []int{1, 2, 4, 5, 6, 8, 10, 15, 20, 25} {
		m := CheckMilestone(day)
		if m != nil {
			t.Errorf("expected nil at day %d, got milestone %s", day, m.Name)
		}
	}
}

func TestGetAllMilestones_Count(t *testing.T) {
	all := GetAllMilestones()
	if len(all) != 6 {
		t.Errorf("expected 6 milestones, got %d", len(all))
	}
}

func TestGetNextMilestone_BeforeFirst(t *testing.T) {
	next := GetNextMilestone(1)
	if next == nil {
		t.Fatal("expected next milestone, got nil")
	}
	if next.Days != 3 {
		t.Errorf("expected next milestone at day 3, got %d", next.Days)
	}
}

func TestGetNextMilestone_AfterDay7(t *testing.T) {
	next := GetNextMilestone(7)
	if next == nil {
		t.Fatal("expected next milestone, got nil")
	}
	if next.Days != 14 {
		t.Errorf("expected next milestone at day 14, got %d", next.Days)
	}
}

func TestGetNextMilestone_AfterAll(t *testing.T) {
	next := GetNextMilestone(100)
	if next != nil {
		t.Errorf("expected nil after all milestones, got %s", next.Name)
	}
}

func TestMilestoneRewards_CoinsAmount(t *testing.T) {
	// 驗證金幣獎勵遞增
	prevCoins := 0
	for _, m := range GetAllMilestones() {
		for _, r := range m.Rewards {
			if r.Type == MilestoneRewardCoins {
				if r.Amount <= prevCoins {
					t.Errorf("milestone day %d coins %d should be > prev %d", m.Days, r.Amount, prevCoins)
				}
				prevCoins = r.Amount
			}
		}
	}
}
