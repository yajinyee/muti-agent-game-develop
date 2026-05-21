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

func TestBuyWeapon_TornadoNotBuyable(t *testing.T) {
	// 龍捲風砲不能購買（DAY-134）
	m := New()
	ok, _ := m.BuyWeapon("p1", WeaponTornado, 99999)
	if ok {
		t.Fatal("expected tornado buy failure (not purchasable)")
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

func TestCalcTornadoTargets(t *testing.T) {
	// 龍捲風命中所有目標（DAY-134）
	targets := []TargetPos{
		{InstanceID: "t1", X: 100, Y: 100},
		{InstanceID: "t2", X: 500, Y: 300},
		{InstanceID: "t3", X: 800, Y: 500},
	}
	hit := CalcTornadoTargets(targets)
	if len(hit) != 3 {
		t.Fatalf("expected 3 hits (all targets), got %d", len(hit))
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
	if s.BombCharges != 0 || s.LaserCharges != 0 || s.FreezeCharges != 0 || s.TornadoCharges != 0 {
		t.Fatal("expected zero charges for new player")
	}
}

// ---- DAY-134 自動充能系統測試 ----

func TestRecordKill_BasicCharge(t *testing.T) {
	m := New()
	// 擊破 1 個普通目標（1x），充能進度 +1
	results := m.RecordKill("p1", 2.0)
	if len(results) == 0 {
		t.Fatal("expected charge results")
	}
	snap := m.GetSnapshot("p1")
	// 炸彈需要 20 次，進度應該是 1
	if snap.BombChargeProgress != 1 {
		t.Fatalf("expected BombChargeProgress=1, got %d", snap.BombChargeProgress)
	}
	// 龍捲風需要 50 次，進度應該是 1
	if snap.TornadoChargeProgress != 1 {
		t.Fatalf("expected TornadoChargeProgress=1, got %d", snap.TornadoChargeProgress)
	}
}

func TestRecordKill_HighMultiplierBonus(t *testing.T) {
	m := New()
	// 擊破 1 個高倍率目標（≥10x），充能進度 +2
	results := m.RecordKill("p1", 15.0)
	if len(results) == 0 {
		t.Fatal("expected charge results")
	}
	snap := m.GetSnapshot("p1")
	if snap.BombChargeProgress != 2 {
		t.Fatalf("expected BombChargeProgress=2 for 15x target, got %d", snap.BombChargeProgress)
	}
}

func TestRecordKill_VeryHighMultiplierBonus(t *testing.T) {
	m := New()
	// 擊破 1 個超高倍率目標（≥30x），充能進度 +3
	results := m.RecordKill("p1", 50.0)
	if len(results) == 0 {
		t.Fatal("expected charge results")
	}
	snap := m.GetSnapshot("p1")
	if snap.BombChargeProgress != 3 {
		t.Fatalf("expected BombChargeProgress=3 for 50x target, got %d", snap.BombChargeProgress)
	}
}

func TestRecordKill_FreezeChargeUnlock(t *testing.T) {
	m := New()
	// 冰凍砲需要 15 次，擊破 15 個普通目標應該充滿一發
	var unlocked bool
	for i := 0; i < 15; i++ {
		results := m.RecordKill("p1", 2.0)
		for _, r := range results {
			if r.WeaponType == WeaponFreeze && r.ChargeUnlocked {
				unlocked = true
			}
		}
	}
	if !unlocked {
		t.Fatal("expected freeze charge to unlock after 15 kills")
	}
	snap := m.GetSnapshot("p1")
	if snap.FreezeCharges != 1 {
		t.Fatalf("expected FreezeCharges=1, got %d", snap.FreezeCharges)
	}
}

func TestRecordKill_BombChargeUnlock(t *testing.T) {
	m := New()
	// 炸彈砲需要 20 次
	var unlocked bool
	for i := 0; i < 20; i++ {
		results := m.RecordKill("p1", 2.0)
		for _, r := range results {
			if r.WeaponType == WeaponBomb && r.ChargeUnlocked {
				unlocked = true
			}
		}
	}
	if !unlocked {
		t.Fatal("expected bomb charge to unlock after 20 kills")
	}
	snap := m.GetSnapshot("p1")
	if snap.BombCharges != 1 {
		t.Fatalf("expected BombCharges=1, got %d", snap.BombCharges)
	}
}

func TestRecordKill_TornadoChargeUnlock(t *testing.T) {
	m := New()
	// 龍捲風砲需要 50 次
	var unlocked bool
	for i := 0; i < 50; i++ {
		results := m.RecordKill("p1", 2.0)
		for _, r := range results {
			if r.WeaponType == WeaponTornado && r.ChargeUnlocked {
				unlocked = true
			}
		}
	}
	if !unlocked {
		t.Fatal("expected tornado charge to unlock after 50 kills")
	}
	snap := m.GetSnapshot("p1")
	if snap.TornadoCharges != 1 {
		t.Fatalf("expected TornadoCharges=1, got %d", snap.TornadoCharges)
	}
}

func TestRecordKill_MaxChargesNoProgress(t *testing.T) {
	m := New()
	// 先充滿炸彈砲（3 發）
	for i := 0; i < 3; i++ {
		m.AddCharge("p1", WeaponBomb)
	}
	// 再擊破目標，炸彈砲已滿，不應再累積進度
	m.RecordKill("p1", 2.0)
	snap := m.GetSnapshot("p1")
	if snap.BombChargeProgress != 0 {
		t.Fatalf("expected BombChargeProgress=0 when at max charges, got %d", snap.BombChargeProgress)
	}
}

func TestRecordKill_ProgressCarryOver(t *testing.T) {
	m := New()
	// 冰凍砲需要 15 次，擊破 17 次，進度應該是 2（17-15=2）
	for i := 0; i < 17; i++ {
		m.RecordKill("p1", 2.0)
	}
	snap := m.GetSnapshot("p1")
	if snap.FreezeCharges != 1 {
		t.Fatalf("expected FreezeCharges=1, got %d", snap.FreezeCharges)
	}
	if snap.FreezeChargeProgress != 2 {
		t.Fatalf("expected FreezeChargeProgress=2 (carry over), got %d", snap.FreezeChargeProgress)
	}
}

func TestRecordKill_MultiplePlayers(t *testing.T) {
	m := New()
	// 兩個玩家各自獨立充能
	for i := 0; i < 15; i++ {
		m.RecordKill("p1", 2.0)
	}
	for i := 0; i < 5; i++ {
		m.RecordKill("p2", 2.0)
	}
	snap1 := m.GetSnapshot("p1")
	snap2 := m.GetSnapshot("p2")
	if snap1.FreezeCharges != 1 {
		t.Fatalf("p1: expected FreezeCharges=1, got %d", snap1.FreezeCharges)
	}
	if snap2.FreezeCharges != 0 {
		t.Fatalf("p2: expected FreezeCharges=0, got %d", snap2.FreezeCharges)
	}
}
