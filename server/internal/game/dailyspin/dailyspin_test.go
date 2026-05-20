// Package dailyspin — 每日轉盤測試（DAY-092）
package dailyspin

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("NewManager returned nil")
	}
	if len(m.normalSlots) != 8 {
		t.Errorf("expected 8 normal slots, got %d", len(m.normalSlots))
	}
	if len(m.superSlots) != 8 {
		t.Errorf("expected 8 super slots, got %d", len(m.superSlots))
	}
}

func TestCanSpinFirstTime(t *testing.T) {
	m := NewManager()
	if !m.CanSpin("player-1") {
		t.Error("new player should be able to spin")
	}
}

func TestSpinOnce(t *testing.T) {
	m := NewManager()
	result := m.Spin("player-1", 1)
	if result == nil {
		t.Fatal("first spin should succeed")
	}
	if result.SlotIndex < 0 || result.SlotIndex >= 8 {
		t.Errorf("slot index out of range: %d", result.SlotIndex)
	}
	if result.NextSpinAt <= 0 {
		t.Error("next_spin_at should be set")
	}
}

func TestCannotSpinTwiceToday(t *testing.T) {
	m := NewManager()
	m.Spin("player-1", 1)
	if m.CanSpin("player-1") {
		t.Error("should not be able to spin twice today")
	}
	result := m.Spin("player-1", 1)
	if result != nil {
		t.Error("second spin today should return nil")
	}
}

func TestSuperSpinRequires7Days(t *testing.T) {
	m := NewManager()
	// 連續 6 天不是超級轉盤
	if m.IsSuper("player-1") {
		t.Error("new player should not have super spin")
	}
	// 模擬 7 天連續登入
	result := m.Spin("player-1", 7)
	if result == nil {
		t.Fatal("spin should succeed")
	}
	if !result.IsSuper {
		t.Error("7-day streak should trigger super spin")
	}
}

func TestNormalSpinNotSuper(t *testing.T) {
	m := NewManager()
	result := m.Spin("player-1", 3)
	if result == nil {
		t.Fatal("spin should succeed")
	}
	if result.IsSuper {
		t.Error("3-day streak should not trigger super spin")
	}
}

func TestGetSnapshot(t *testing.T) {
	m := NewManager()
	snap := m.GetSnapshot("player-1")
	if snap["can_spin"] != true {
		t.Error("new player should be able to spin")
	}
	if snap["login_streak"] != 0 {
		t.Error("new player login streak should be 0")
	}
}

func TestGetSnapshotAfterSpin(t *testing.T) {
	m := NewManager()
	m.Spin("player-1", 5)
	snap := m.GetSnapshot("player-1")
	if snap["can_spin"] != false {
		t.Error("after spin, can_spin should be false")
	}
	if snap["login_streak"] != 5 {
		t.Errorf("login_streak should be 5, got %v", snap["login_streak"])
	}
	if snap["total_spins"] != 1 {
		t.Errorf("total_spins should be 1, got %v", snap["total_spins"])
	}
}

func TestSlotWeightsPositive(t *testing.T) {
	m := NewManager()
	for _, slot := range m.normalSlots {
		if slot.Weight <= 0 {
			t.Errorf("slot %d has non-positive weight: %d", slot.ID, slot.Weight)
		}
	}
	for _, slot := range m.superSlots {
		if slot.Weight <= 0 {
			t.Errorf("super slot %d has non-positive weight: %d", slot.ID, slot.Weight)
		}
	}
}

func TestSuperSlotsHigherReward(t *testing.T) {
	m := NewManager()
	// 超級轉盤的金幣獎勵應該比普通轉盤高
	normalMaxCoins := 0
	for _, slot := range m.normalSlots {
		if slot.Type == RewardCoins && slot.Amount > normalMaxCoins {
			normalMaxCoins = slot.Amount
		}
	}
	superMaxCoins := 0
	for _, slot := range m.superSlots {
		if slot.Type == RewardCoins && slot.Amount > superMaxCoins {
			superMaxCoins = slot.Amount
		}
	}
	if superMaxCoins <= normalMaxCoins {
		t.Errorf("super max coins (%d) should be > normal max coins (%d)", superMaxCoins, normalMaxCoins)
	}
}

func TestMultiplePlayersIndependent(t *testing.T) {
	m := NewManager()
	m.Spin("player-1", 1)
	// player-2 應該還能轉
	if !m.CanSpin("player-2") {
		t.Error("player-2 should be able to spin independently")
	}
}
