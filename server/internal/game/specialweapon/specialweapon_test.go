package specialweapon

import (
	"testing"
)

func TestBuyWeapon_Success(t *testing.T) {
	m := New()
	ok, cost := m.BuyWeapon("p1", WeaponBomb, 1000)
	if !ok {
		t.Fatal("expected buy success")
	}
	if cost != 500 {
		t.Fatalf("expected cost=500, got %d", cost)
	}
	snap := m.GetSnapshot("p1")
	if snap.BombCharges != 1 {
		t.Fatalf("expected BombCharges=1, got %d", snap.BombCharges)
	}
}

func TestBuyWeapon_InsufficientCoins(t *testing.T) {
	m := New()
	ok, _ := m.BuyWeapon("p1", WeaponBomb, 100)
	if ok {
		t.Fatal("expected buy failure due to insufficient coins")
	}
}

func TestBuyWeapon_MaxCharges(t *testing.T) {
	m := New()
	for i := 0; i < 3; i++ {
		ok, _ := m.BuyWeapon("p1", WeaponLaser, 10000)
		if !ok {
			t.Fatalf("expected buy success on attempt %d", i+1)
		}
	}
	// 第 4 次應該失敗
	ok, _ := m.BuyWeapon("p1", WeaponLaser, 10000)
	if ok {
		t.Fatal("expected buy failure at max charges")
	}
}

func TestUseWeapon_Success(t *testing.T) {
	m := New()
	m.BuyWeapon("p1", WeaponFreeze, 1000)
	ok := m.UseWeapon("p1", WeaponFreeze)
	if !ok {
		t.Fatal("expected use success")
	}
	snap := m.GetSnapshot("p1")
	if snap.FreezeCharges != 0 {
		t.Fatalf("expected FreezeCharges=0, got %d", snap.FreezeCharges)
	}
}

func TestUseWeapon_NoCharges(t *testing.T) {
	m := New()
	ok := m.UseWeapon("p1", WeaponBomb)
	if ok {
		t.Fatal("expected use failure with no charges")
	}
}

func TestCalcBombTargets(t *testing.T) {
	targets := []TargetPos{
		{InstanceID: "t1", X: 100, Y: 100},
		{InstanceID: "t2", X: 400, Y: 100}, // 距離 300，超出半徑
		{InstanceID: "t3", X: 200, Y: 200}, // 距離 ~141，在半徑內
	}
	hit := CalcBombTargets(100, 100, targets)
	if len(hit) != 2 {
		t.Fatalf("expected 2 hits, got %d: %v", len(hit), hit)
	}
}

func TestCalcLaserTargets(t *testing.T) {
	targets := []TargetPos{
		{InstanceID: "t1", X: 100, Y: 300},
		{InstanceID: "t2", X: 200, Y: 350}, // Y 差 50，在範圍內
		{InstanceID: "t3", X: 300, Y: 500}, // Y 差 200，超出範圍
	}
	hit := CalcLaserTargets(300, targets)
	if len(hit) != 2 {
		t.Fatalf("expected 2 hits, got %d: %v", len(hit), hit)
	}
}

func TestCalcFreezeTargets(t *testing.T) {
	targets := []TargetPos{
		{InstanceID: "t1"},
		{InstanceID: "t2"},
		{InstanceID: "t3"},
	}
	hit := CalcFreezeTargets(targets)
	if len(hit) != 3 {
		t.Fatalf("expected 3 hits, got %d", len(hit))
	}
}

func TestRemovePlayer(t *testing.T) {
	m := New()
	m.BuyWeapon("p1", WeaponBomb, 1000)
	m.RemovePlayer("p1")
	snap := m.GetSnapshot("p1")
	if snap.BombCharges != 0 {
		t.Fatal("expected empty state after remove")
	}
}

func TestGetOrCreate_NewPlayer(t *testing.T) {
	m := New()
	s := m.GetOrCreate("newplayer")
	if s == nil {
		t.Fatal("expected non-nil state")
	}
	if s.BombCharges != 0 || s.LaserCharges != 0 || s.FreezeCharges != 0 {
		t.Fatal("expected zero charges for new player")
	}
}
