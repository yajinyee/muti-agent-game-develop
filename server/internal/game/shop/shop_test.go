package shop

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	items := m.GetItems()
	if len(items) == 0 {
		t.Error("expected items in shop, got 0")
	}
}

func TestGetItems_HasFlashSale(t *testing.T) {
	m := New()
	items := m.GetItems()
	flashCount := 0
	for _, item := range items {
		if item.IsFlashSale {
			flashCount++
		}
	}
	if flashCount != 2 {
		t.Errorf("expected 2 flash sale items, got %d", flashCount)
	}
}

func TestGetItems_FlashSaleDiscount(t *testing.T) {
	m := New()
	items := m.GetItems()
	for _, item := range items {
		if item.IsFlashSale {
			// 特賣價應該比原價低
			if item.Price >= item.OrigPrice {
				t.Errorf("flash sale item %s price %d should be less than orig price %d",
					item.ID, item.Price, item.OrigPrice)
			}
		}
	}
}

func TestBuyItem_Success(t *testing.T) {
	m := New()
	result := m.BuyItem("p1", "bomb_bundle", 10000)
	if !result.Success {
		t.Errorf("expected success, got failure: %s", result.Reason)
	}
	if result.Reward.BombCharge != 3 {
		t.Errorf("expected 3 bomb charges, got %d", result.Reward.BombCharge)
	}
}

func TestBuyItem_InsufficientCoins(t *testing.T) {
	m := New()
	result := m.BuyItem("p1", "bomb_bundle", 100) // 只有 100 金幣，不夠
	if result.Success {
		t.Error("expected failure due to insufficient coins")
	}
	if result.Reason != "insufficient_coins" {
		t.Errorf("expected reason 'insufficient_coins', got '%s'", result.Reason)
	}
}

func TestBuyItem_NotFound(t *testing.T) {
	m := New()
	result := m.BuyItem("p1", "nonexistent_item", 99999)
	if result.Success {
		t.Error("expected failure for nonexistent item")
	}
	if result.Reason != "item_not_found" {
		t.Errorf("expected reason 'item_not_found', got '%s'", result.Reason)
	}
}

func TestBuyItem_DailyLimit(t *testing.T) {
	m := New()
	// bomb_bundle 每日上限 3 次
	for i := 0; i < 3; i++ {
		result := m.BuyItem("p1", "bomb_bundle", 99999)
		if !result.Success {
			t.Errorf("purchase %d should succeed, got: %s", i+1, result.Reason)
		}
	}
	// 第 4 次應該失敗
	result := m.BuyItem("p1", "bomb_bundle", 99999)
	if result.Success {
		t.Error("expected failure due to daily limit")
	}
	if result.Reason != "daily_limit_reached" {
		t.Errorf("expected reason 'daily_limit_reached', got '%s'", result.Reason)
	}
}

func TestBuyItem_DifferentPlayers(t *testing.T) {
	m := New()
	// 不同玩家各自有獨立的每日上限
	for i := 0; i < 3; i++ {
		m.BuyItem("p1", "bomb_bundle", 99999)
	}
	// p2 應該還能購買
	result := m.BuyItem("p2", "bomb_bundle", 99999)
	if !result.Success {
		t.Errorf("p2 should be able to buy, got: %s", result.Reason)
	}
}

func TestGetPlayerDailyPurchases(t *testing.T) {
	m := New()
	m.BuyItem("p1", "bomb_bundle", 99999)
	m.BuyItem("p1", "bomb_bundle", 99999)
	m.BuyItem("p1", "laser_bundle", 99999)

	purchases := m.GetPlayerDailyPurchases("p1")
	if purchases["bomb_bundle"] != 2 {
		t.Errorf("expected 2 bomb_bundle purchases, got %d", purchases["bomb_bundle"])
	}
	if purchases["laser_bundle"] != 1 {
		t.Errorf("expected 1 laser_bundle purchase, got %d", purchases["laser_bundle"])
	}
}

func TestGetSnapshot(t *testing.T) {
	m := New()
	snap := m.GetSnapshot()
	if len(snap.Items) == 0 {
		t.Error("expected items in snapshot")
	}
	if snap.FlashSaleEndAt <= 0 {
		t.Error("expected flash sale end time")
	}
	if snap.SecondsLeft <= 0 {
		t.Error("expected positive seconds left")
	}
}

func TestGetFlashSaleEndAt(t *testing.T) {
	m := New()
	endAt := m.GetFlashSaleEndAt()
	if endAt <= 0 {
		t.Error("expected positive flash sale end time")
	}
}
