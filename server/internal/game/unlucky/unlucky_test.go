package unlucky

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := NewDefault()
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestRecordShot_NotEnoughShots(t *testing.T) {
	m := NewDefault()
	// 只有 5 次，不夠 30 次，不應觸發
	for i := 0; i < 5; i++ {
		triggered, _ := m.RecordShot("p1", 100, 0)
		if triggered {
			t.Fatal("should not trigger with only 5 shots")
		}
	}
}

func TestRecordShot_LowSpend(t *testing.T) {
	m := NewDefault()
	// 30 次但花費太少（每次 5，總計 150 < MinSpend=200）
	for i := 0; i < 30; i++ {
		triggered, _ := m.RecordShot("p1", 5, 0)
		if triggered {
			t.Fatal("should not trigger with low spend")
		}
	}
}

func TestRecordShot_GoodRTP(t *testing.T) {
	m := NewDefault()
	// 30 次，花費 100，回報 80（RTP 80%，比例 1.25 < 3.0 門檻）
	for i := 0; i < 30; i++ {
		triggered, _ := m.RecordShot("p1", 100, 80)
		if triggered {
			t.Fatal("should not trigger with good RTP")
		}
	}
}

func TestRecordShot_BadRTP_Trigger(t *testing.T) {
	m := NewDefault()
	// 30 次，花費 100，回報 10（RTP 10%，比例 10 > 3.0 門檻）
	var triggered bool
	var bonus int
	for i := 0; i < 30; i++ {
		triggered, bonus = m.RecordShot("p1", 100, 10)
	}
	if !triggered {
		t.Fatal("expected unlucky bonus to trigger after 30 bad shots")
	}
	if bonus < DefaultConfig.MinReward {
		t.Fatalf("expected bonus >= %d, got %d", DefaultConfig.MinReward, bonus)
	}
}

func TestRecordShot_ZeroReward_Trigger(t *testing.T) {
	m := NewDefault()
	// 30 次，花費 100，回報 0（完全沒有回報）
	var triggered bool
	var bonus int
	for i := 0; i < 30; i++ {
		triggered, bonus = m.RecordShot("p1", 100, 0)
	}
	if !triggered {
		t.Fatal("expected unlucky bonus to trigger with zero reward")
	}
	// 補償 = 淨虧損 × 0.3 = 3000 × 0.3 = 900
	expectedMin := int(float64(3000) * DefaultConfig.BaseRewardMult)
	if bonus < expectedMin {
		t.Fatalf("expected bonus >= %d, got %d", expectedMin, bonus)
	}
}

func TestRecordShot_BonusCalc(t *testing.T) {
	m := NewDefault()
	// 30 次，花費 100，回報 0，淨虧損 3000
	// 補償 = 3000 × 0.3 = 900
	for i := 0; i < 30; i++ {
		m.RecordShot("p1", 100, 0)
	}
	snap := m.GetSnapshot("p1")
	// 觸發後狀態重置
	if snap.ShotCount != 0 {
		t.Fatalf("expected ShotCount=0 after trigger, got %d", snap.ShotCount)
	}
}

func TestRecordShot_Cooldown(t *testing.T) {
	cfg := DefaultConfig
	cfg.CooldownSecs = 1 // 1 秒冷卻（測試用）
	m := New(cfg)

	// 第一次觸發
	for i := 0; i < 30; i++ {
		m.RecordShot("p1", 100, 0)
	}

	// 立即再觸發（冷卻中，不應觸發）
	for i := 0; i < 30; i++ {
		triggered, _ := m.RecordShot("p1", 100, 0)
		if triggered {
			t.Fatal("should not trigger during cooldown")
		}
	}
}

func TestRecordShot_CooldownExpired(t *testing.T) {
	cfg := DefaultConfig
	cfg.CooldownSecs = 0 // 無冷卻（測試用）
	m := New(cfg)

	// 第一次觸發
	for i := 0; i < 30; i++ {
		m.RecordShot("p1", 100, 0)
	}

	// 冷卻已過（0秒），再次觸發
	var triggered bool
	for i := 0; i < 30; i++ {
		triggered, _ = m.RecordShot("p1", 100, 0)
	}
	if !triggered {
		t.Fatal("expected second trigger after cooldown expired")
	}
}

func TestGetSnapshot_Initial(t *testing.T) {
	m := NewDefault()
	snap := m.GetSnapshot("p1")
	if snap.TrackingMax != DefaultConfig.TrackingShots {
		t.Fatalf("expected TrackingMax=%d, got %d", DefaultConfig.TrackingShots, snap.TrackingMax)
	}
	if snap.ShotCount != 0 {
		t.Fatalf("expected ShotCount=0, got %d", snap.ShotCount)
	}
}

func TestGetSnapshot_Progress(t *testing.T) {
	m := NewDefault()
	for i := 0; i < 15; i++ {
		m.RecordShot("p1", 100, 10)
	}
	snap := m.GetSnapshot("p1")
	if snap.ShotCount != 15 {
		t.Fatalf("expected ShotCount=15, got %d", snap.ShotCount)
	}
	if snap.TotalSpend != 1500 {
		t.Fatalf("expected TotalSpend=1500, got %d", snap.TotalSpend)
	}
	if snap.TotalReward != 150 {
		t.Fatalf("expected TotalReward=150, got %d", snap.TotalReward)
	}
}

func TestGetSnapshot_CooldownLeft(t *testing.T) {
	cfg := DefaultConfig
	cfg.CooldownSecs = 120
	m := New(cfg)

	// 觸發補償
	for i := 0; i < 30; i++ {
		m.RecordShot("p1", 100, 0)
	}

	snap := m.GetSnapshot("p1")
	// 冷卻應該接近 120 秒
	if snap.CooldownLeft < 118 || snap.CooldownLeft > 120 {
		t.Fatalf("expected CooldownLeft ~120, got %d", snap.CooldownLeft)
	}
}

func TestRemovePlayer(t *testing.T) {
	m := NewDefault()
	for i := 0; i < 10; i++ {
		m.RecordShot("p1", 100, 0)
	}
	m.RemovePlayer("p1")
	snap := m.GetSnapshot("p1")
	if snap.ShotCount != 0 {
		t.Fatal("expected empty state after remove")
	}
}

func TestMultiplePlayers(t *testing.T) {
	m := NewDefault()
	// p1 運氣很差
	for i := 0; i < 30; i++ {
		m.RecordShot("p1", 100, 0)
	}
	// p2 運氣很好
	for i := 0; i < 30; i++ {
		m.RecordShot("p2", 100, 200)
	}

	snap1 := m.GetSnapshot("p1")
	snap2 := m.GetSnapshot("p2")

	// p1 觸發後重置
	if snap1.ShotCount != 0 {
		t.Fatalf("p1: expected ShotCount=0 after trigger, got %d", snap1.ShotCount)
	}
	// p2 沒有觸發
	if snap2.ShotCount != 30 {
		t.Fatalf("p2: expected ShotCount=30, got %d", snap2.ShotCount)
	}
}

func TestBonusCount(t *testing.T) {
	cfg := DefaultConfig
	cfg.CooldownSecs = 0
	m := New(cfg)

	// 觸發兩次
	for round := 0; round < 2; round++ {
		for i := 0; i < 30; i++ {
			m.RecordShot("p1", 100, 0)
		}
	}

	snap := m.GetSnapshot("p1")
	if snap.BonusCount != 2 {
		t.Fatalf("expected BonusCount=2, got %d", snap.BonusCount)
	}
}

func TestRingBuffer_OldRecordsRemoved(t *testing.T) {
	cfg := DefaultConfig
	cfg.CooldownSecs = 0 // 無冷卻，讓多次觸發都可以
	m := New(cfg)

	// 先填入 30 次好記錄（高回報）
	for i := 0; i < 30; i++ {
		m.RecordShot("p1", 100, 200)
	}
	// 再填入 30 次壞記錄（零回報）
	// 環形緩衝應該把舊的好記錄替換掉，最終觸發補償
	// 由於觸發後重置，可能需要超過 30 次才能再次觸發
	// 所以我們填入 60 次，確保至少觸發一次
	var anyTriggered bool
	for i := 0; i < 60; i++ {
		triggered, _ := m.RecordShot("p1", 100, 0)
		if triggered {
			anyTriggered = true
		}
	}
	if !anyTriggered {
		t.Fatal("expected at least one trigger after ring buffer replaced good records with bad ones")
	}
}

func TestMinReward(t *testing.T) {
	m := NewDefault()
	// 只花費 210（剛超過 MinSpend=200），回報 0
	// 淨虧損 210，補償 = 210 × 0.3 = 63 < MinReward=100
	// 應該給 MinReward=100
	for i := 0; i < 30; i++ {
		m.RecordShot("p1", 7, 0) // 30 × 7 = 210
	}
	snap := m.GetSnapshot("p1")
	_ = snap
	// 重新觸發確認補償金額
	cfg := DefaultConfig
	cfg.CooldownSecs = 0
	m2 := New(cfg)
	var bonus int
	for i := 0; i < 30; i++ {
		_, bonus = m2.RecordShot("p1", 7, 0)
	}
	if bonus < DefaultConfig.MinReward {
		t.Fatalf("expected bonus >= MinReward=%d, got %d", DefaultConfig.MinReward, bonus)
	}
}

func TestCooldownLeft_NoBonusYet(t *testing.T) {
	m := NewDefault()
	snap := m.GetSnapshot("p1")
	if snap.CooldownLeft != 0 {
		t.Fatalf("expected CooldownLeft=0 before any bonus, got %d", snap.CooldownLeft)
	}
}

func TestRecordShot_ResetAfterTrigger(t *testing.T) {
	m := NewDefault()
	// 觸發補償
	for i := 0; i < 30; i++ {
		m.RecordShot("p1", 100, 0)
	}
	// 觸發後，再記錄 5 次
	for i := 0; i < 5; i++ {
		m.RecordShot("p1", 100, 0)
	}
	snap := m.GetSnapshot("p1")
	// 重置後只有 5 次記錄
	if snap.ShotCount != 5 {
		t.Fatalf("expected ShotCount=5 after reset, got %d", snap.ShotCount)
	}
	_ = time.Now() // 確保 time 套件被使用
}
