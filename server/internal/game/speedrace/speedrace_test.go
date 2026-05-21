// Package speedrace 競速獵殺系統單元測試
package speedrace

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := NewDefault()
	if m == nil {
		t.Fatal("NewDefault() returned nil")
	}
	if !m.CanStart() {
		t.Error("新建管理器應該可以開始競速")
	}
}

func TestStartRace_Basic(t *testing.T) {
	m := NewDefault()
	sess := m.StartRace("inst1", "T103", "流星", 15.0)
	if sess == nil {
		t.Fatal("StartRace 應該成功")
	}
	if !sess.IsActive {
		t.Error("session 應該是 active")
	}
	if sess.TargetInstanceID != "inst1" {
		t.Errorf("TargetInstanceID 錯誤: %s", sess.TargetInstanceID)
	}
}

func TestStartRace_MultTooLow(t *testing.T) {
	m := NewDefault()
	sess := m.StartRace("inst1", "T001", "雜草", 2.0) // 低於 MinMultiplier=10
	if sess != nil {
		t.Error("倍率太低不應該開始競速")
	}
}

func TestStartRace_AlreadyActive(t *testing.T) {
	m := NewDefault()
	m.StartRace("inst1", "T103", "流星", 15.0)
	sess2 := m.StartRace("inst2", "T104", "金草", 20.0)
	if sess2 != nil {
		t.Error("已有進行中的競速，不應該開始新的")
	}
}

func TestRecordKill_NotRaceTarget(t *testing.T) {
	m := NewDefault()
	m.StartRace("inst1", "T103", "流星", 15.0)
	rank, mult, isRace := m.RecordKill("inst2", "p1", "玩家1") // 不同 instanceID
	if isRace {
		t.Error("不是競速目標，isRaceTarget 應該是 false")
	}
	if rank != 0 || mult != 1.0 {
		t.Errorf("非競速目標應回傳 rank=0, mult=1.0，得到 rank=%d, mult=%.1f", rank, mult)
	}
}

func TestRecordKill_FirstPlace(t *testing.T) {
	m := NewDefault()
	m.StartRace("inst1", "T103", "流星", 15.0)
	rank, mult, isRace := m.RecordKill("inst1", "p1", "玩家1")
	if !isRace {
		t.Error("應該是競速目標")
	}
	if rank != 1 {
		t.Errorf("第一個擊破應該是 rank=1，得到 %d", rank)
	}
	if mult != 3.0 {
		t.Errorf("第一名倍率應該是 3.0，得到 %.1f", mult)
	}
}

func TestRecordKill_FirstPlaceEndsRace(t *testing.T) {
	m := NewDefault()
	m.StartRace("inst1", "T103", "流星", 15.0)
	m.RecordKill("inst1", "p1", "玩家1")

	snap := m.GetSnapshot()
	if snap.IsActive {
		t.Error("第一名擊破後競速應該結束")
	}
}

func TestRecordKill_NoActiveRace(t *testing.T) {
	m := NewDefault()
	rank, mult, isRace := m.RecordKill("inst1", "p1", "玩家1")
	if isRace {
		t.Error("沒有競速時 isRaceTarget 應該是 false")
	}
	if rank != 0 || mult != 1.0 {
		t.Errorf("沒有競速應回傳 rank=0, mult=1.0，得到 rank=%d, mult=%.1f", rank, mult)
	}
}

func TestCheckExpiry(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Duration = 0.01 // 10ms，快速超時
	m := New(cfg)
	m.StartRace("inst1", "T103", "流星", 15.0)

	time.Sleep(20 * time.Millisecond)
	expired := m.CheckExpiry()
	if !expired {
		t.Error("應該已超時")
	}

	// 再次檢查不應該再次觸發
	expired2 := m.CheckExpiry()
	if expired2 {
		t.Error("超時後再次檢查不應該再次觸發")
	}
}

func TestCheckExpiry_NotExpired(t *testing.T) {
	m := NewDefault()
	m.StartRace("inst1", "T103", "流星", 15.0)
	expired := m.CheckExpiry()
	if expired {
		t.Error("剛開始不應該超時")
	}
}

func TestCancelRace(t *testing.T) {
	m := NewDefault()
	m.StartRace("inst1", "T103", "流星", 15.0)
	cancelled := m.CancelRace("inst1")
	if !cancelled {
		t.Error("取消競速應該成功")
	}
	snap := m.GetSnapshot()
	if snap.IsActive {
		t.Error("取消後競速應該結束")
	}
}

func TestCancelRace_WrongInstance(t *testing.T) {
	m := NewDefault()
	m.StartRace("inst1", "T103", "流星", 15.0)
	cancelled := m.CancelRace("inst2") // 不同 instanceID
	if cancelled {
		t.Error("不同 instanceID 不應該取消成功")
	}
}

func TestCooldown(t *testing.T) {
	cfg := DefaultConfig()
	cfg.CooldownSecs = 100.0 // 100 秒冷卻
	cfg.Duration = 0.01      // 快速超時
	m := New(cfg)
	m.StartRace("inst1", "T103", "流星", 15.0)
	time.Sleep(20 * time.Millisecond)
	m.CheckExpiry() // 觸發超時

	// 冷卻中不應該開始新競速
	if m.CanStart() {
		t.Error("冷卻中不應該可以開始新競速")
	}
	sess := m.StartRace("inst2", "T104", "金草", 20.0)
	if sess != nil {
		t.Error("冷卻中不應該開始新競速")
	}
}

func TestGetSnapshot_NoRace(t *testing.T) {
	m := NewDefault()
	snap := m.GetSnapshot()
	if snap.IsActive {
		t.Error("沒有競速時 IsActive 應該是 false")
	}
}

func TestGetSnapshot_Active(t *testing.T) {
	m := NewDefault()
	m.StartRace("inst1", "T103", "流星", 15.0)
	snap := m.GetSnapshot()
	if !snap.IsActive {
		t.Error("競速中 IsActive 應該是 true")
	}
	if snap.TargetInstanceID != "inst1" {
		t.Errorf("TargetInstanceID 錯誤: %s", snap.TargetInstanceID)
	}
	if snap.SecondsLeft <= 0 {
		t.Error("SecondsLeft 應該 > 0")
	}
	if snap.BonusMult != 3.0 {
		t.Errorf("BonusMult 應該是 3.0，得到 %.1f", snap.BonusMult)
	}
}

func TestIsRaceTarget(t *testing.T) {
	m := NewDefault()
	m.StartRace("inst1", "T103", "流星", 15.0)
	if !m.IsRaceTarget("inst1") {
		t.Error("inst1 應該是競速目標")
	}
	if m.IsRaceTarget("inst2") {
		t.Error("inst2 不應該是競速目標")
	}
}

func TestMultiplePlayers(t *testing.T) {
	m := NewDefault()
	m.StartRace("inst1", "T103", "流星", 15.0)

	// 第一名
	rank1, mult1, _ := m.RecordKill("inst1", "p1", "玩家1")
	if rank1 != 1 || mult1 != 3.0 {
		t.Errorf("第一名應該是 rank=1, mult=3.0，得到 rank=%d, mult=%.1f", rank1, mult1)
	}

	// 競速已結束，後續擊破不計入
	rank2, mult2, isRace2 := m.RecordKill("inst1", "p2", "玩家2")
	if isRace2 {
		t.Error("競速結束後不應該是競速目標")
	}
	if rank2 != 0 || mult2 != 1.0 {
		t.Errorf("競速結束後應回傳 rank=0, mult=1.0，得到 rank=%d, mult=%.1f", rank2, mult2)
	}
}
