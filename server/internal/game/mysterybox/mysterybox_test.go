package mysterybox

import (
	"testing"
)

func TestTryDropBox_NormalTarget(t *testing.T) {
	m := New()
	// 跑 10000 次，確認有掉落（機率 8%+4%+1.5%+0.5% = 14%）
	dropped := 0
	for i := 0; i < 10000; i++ {
		if box := m.TryDropBox(false); box != nil {
			dropped++
		}
	}
	// 期望約 1400 次，允許 ±500 誤差
	if dropped < 900 || dropped > 1900 {
		t.Fatalf("unexpected drop count: %d (expected ~1400)", dropped)
	}
}

func TestTryDropBox_BossBoost(t *testing.T) {
	m := New()
	// BOSS 擊殺時傳說寶箱機率提升 10 倍（0.5% → 5%）
	legendaryCount := 0
	for i := 0; i < 10000; i++ {
		if box := m.TryDropBox(true); box != nil && box.Rarity == RarityLegendary {
			legendaryCount++
		}
	}
	// 期望約 500 次（5%），允許 ±200 誤差
	if legendaryCount < 300 || legendaryCount > 700 {
		t.Fatalf("unexpected legendary count with boss boost: %d (expected ~500)", legendaryCount)
	}
}

func TestOpenBox_CommonRewards(t *testing.T) {
	m := New()
	// 開 1000 個普通寶箱，確認獎勵都在定義範圍內
	validTypes := map[RewardType]bool{
		RewardCoins:        true,
		RewardFreezeCharge: true,
	}
	for i := 0; i < 1000; i++ {
		reward := m.OpenBox(RarityCommon)
		if reward == nil {
			t.Fatal("expected non-nil reward")
		}
		if !validTypes[reward.Type] {
			t.Fatalf("unexpected reward type for common box: %s", reward.Type)
		}
	}
}

func TestOpenBox_LegendaryRewards(t *testing.T) {
	m := New()
	// 開 1000 個傳說寶箱，確認獎勵都在定義範圍內
	validTypes := map[RewardType]bool{
		RewardCoins:         true,
		RewardBombCharge:    true,
		RewardLaserCharge:   true,
		RewardFreezeCharge:  true,
		RewardMultiplier:    true,
		RewardJackpotTicket: true,
	}
	for i := 0; i < 1000; i++ {
		reward := m.OpenBox(RarityLegendary)
		if reward == nil {
			t.Fatal("expected non-nil reward")
		}
		if !validTypes[reward.Type] {
			t.Fatalf("unexpected reward type for legendary box: %s", reward.Type)
		}
	}
}

func TestOpenBox_InvalidRarity(t *testing.T) {
	m := New()
	reward := m.OpenBox("invalid_rarity")
	if reward != nil {
		t.Fatal("expected nil reward for invalid rarity")
	}
}

func TestPendingMultiplier_SetAndConsume(t *testing.T) {
	m := New()
	m.SetPendingMultiplier("p1", 3.0)

	// 查詢（不消耗）
	mult := m.GetPendingMultiplier("p1")
	if mult != 3.0 {
		t.Fatalf("expected mult=3.0, got %.1f", mult)
	}

	// 消耗
	consumed := m.ConsumePendingMultiplier("p1")
	if consumed != 3.0 {
		t.Fatalf("expected consumed=3.0, got %.1f", consumed)
	}

	// 再次消耗應該是 1.0（已消耗）
	again := m.ConsumePendingMultiplier("p1")
	if again != 1.0 {
		t.Fatalf("expected 1.0 after consume, got %.1f", again)
	}
}

func TestPendingMultiplier_NoMultiplier(t *testing.T) {
	m := New()
	mult := m.ConsumePendingMultiplier("p1")
	if mult != 1.0 {
		t.Fatalf("expected 1.0 for no multiplier, got %.1f", mult)
	}
}

func TestRemovePlayer(t *testing.T) {
	m := New()
	m.SetPendingMultiplier("p1", 2.0)
	m.RemovePlayer("p1")
	mult := m.GetPendingMultiplier("p1")
	if mult != 1.0 {
		t.Fatal("expected 1.0 after remove")
	}
}

func TestGetBoxDef(t *testing.T) {
	def := GetBoxDef(RarityEpic)
	if def == nil {
		t.Fatal("expected non-nil box def")
	}
	if def.Rarity != RarityEpic {
		t.Fatalf("expected epic rarity, got %s", def.Rarity)
	}
}

func TestAllBoxesHaveRewards(t *testing.T) {
	for _, box := range AvailableBoxes {
		if len(box.Rewards) == 0 {
			t.Fatalf("box %s has no rewards", box.Rarity)
		}
		totalWeight := 0
		for _, r := range box.Rewards {
			if r.Weight <= 0 {
				t.Fatalf("box %s has reward with zero weight: %s", box.Rarity, r.Type)
			}
			totalWeight += r.Weight
		}
		if totalWeight <= 0 {
			t.Fatalf("box %s has zero total weight", box.Rarity)
		}
	}
}
