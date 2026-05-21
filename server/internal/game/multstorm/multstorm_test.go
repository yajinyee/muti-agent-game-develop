// Package multstorm 倍率風暴系統單元測試
package multstorm

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := NewDefault()
	if m == nil {
		t.Fatal("NewDefault() returned nil")
	}
	if m.IsActive() {
		t.Error("初始狀態不應該有活躍風暴")
	}
}

func TestGetMultBoost_NoStorm(t *testing.T) {
	m := NewDefault()
	boost := m.GetMultBoost()
	if boost != 1.0 {
		t.Errorf("無風暴時倍率加成應該是 1.0，得到 %.1f", boost)
	}
}

func TestForceStart_Basic(t *testing.T) {
	m := NewDefault()
	sess := m.ForceStart(0) // 閃電風暴
	if sess == nil {
		t.Fatal("ForceStart 應該成功")
	}
	if !sess.IsActive {
		t.Error("session 應該是 active")
	}
	if sess.Tier.MultBoost != 2.0 {
		t.Errorf("閃電風暴倍率應該是 2.0，得到 %.1f", sess.Tier.MultBoost)
	}
}

func TestForceStart_TierIndex(t *testing.T) {
	m := NewDefault()
	// 海嘯風暴（index=1）
	sess := m.ForceStart(1)
	if sess.Tier.MultBoost != 3.0 {
		t.Errorf("海嘯風暴倍率應該是 3.0，得到 %.1f", sess.Tier.MultBoost)
	}
	// 彩虹風暴（index=2）
	m2 := NewDefault()
	sess2 := m2.ForceStart(2)
	if sess2.Tier.MultBoost != 5.0 {
		t.Errorf("彩虹風暴倍率應該是 5.0，得到 %.1f", sess2.Tier.MultBoost)
	}
}

func TestForceStart_InvalidIndex(t *testing.T) {
	m := NewDefault()
	sess := m.ForceStart(99) // 無效 index，應該 fallback 到 0
	if sess == nil {
		t.Fatal("ForceStart 應該成功（fallback 到 index 0）")
	}
}

func TestGetMultBoost_Active(t *testing.T) {
	m := NewDefault()
	m.ForceStart(1) // 海嘯風暴 3.0x
	boost := m.GetMultBoost()
	if boost != 3.0 {
		t.Errorf("海嘯風暴倍率應該是 3.0，得到 %.1f", boost)
	}
}

func TestIsActive(t *testing.T) {
	m := NewDefault()
	if m.IsActive() {
		t.Error("初始狀態不應該是 active")
	}
	m.ForceStart(0)
	if !m.IsActive() {
		t.Error("ForceStart 後應該是 active")
	}
}

func TestCheckExpiry_NotExpired(t *testing.T) {
	m := NewDefault()
	m.ForceStart(0)
	expired := m.CheckExpiry()
	if expired {
		t.Error("剛開始不應該過期")
	}
}

func TestCheckExpiry_Expired(t *testing.T) {
	cfg := DefaultConfig()
	m := New(cfg)
	// 手動建立一個已過期的 session
	m.session = &StormSession{
		Tier:     StormTiers[0],
		StartAt:  time.Now().Add(-30 * time.Second),
		EndAt:    time.Now().Add(-1 * time.Second), // 已過期
		IsActive: true,
	}
	expired := m.CheckExpiry()
	if !expired {
		t.Error("應該已過期")
	}
	if m.IsActive() {
		t.Error("過期後不應該是 active")
	}
}

func TestCheckExpiry_NoSession(t *testing.T) {
	m := NewDefault()
	expired := m.CheckExpiry()
	if expired {
		t.Error("無 session 時不應該觸發過期")
	}
}

func TestGetSnapshot_NoStorm(t *testing.T) {
	m := NewDefault()
	snap := m.GetSnapshot()
	if snap.IsActive {
		t.Error("無風暴時 IsActive 應該是 false")
	}
}

func TestGetSnapshot_Active(t *testing.T) {
	m := NewDefault()
	m.ForceStart(2) // 彩虹風暴
	snap := m.GetSnapshot()
	if !snap.IsActive {
		t.Error("風暴中 IsActive 應該是 true")
	}
	if snap.MultBoost != 5.0 {
		t.Errorf("彩虹風暴倍率應該是 5.0，得到 %.1f", snap.MultBoost)
	}
	if snap.SecondsLeft <= 0 {
		t.Error("SecondsLeft 應該 > 0")
	}
	if snap.TierIcon != "🌈" {
		t.Errorf("彩虹風暴圖示應該是 🌈，得到 %s", snap.TierIcon)
	}
}

func TestCooldown(t *testing.T) {
	cfg := DefaultConfig()
	cfg.CooldownSecs = 100.0
	m := New(cfg)
	// 手動設定 lastEndAt
	m.lastEndAt = time.Now()
	// 嘗試觸發（應該被冷卻阻擋）
	sess := m.TryTrigger()
	if sess != nil {
		t.Error("冷卻中不應該觸發風暴")
	}
}

func TestStormTiers_Count(t *testing.T) {
	if len(StormTiers) != 3 {
		t.Errorf("應該有 3 個風暴等級，得到 %d", len(StormTiers))
	}
}

func TestStormTiers_MultBoosts(t *testing.T) {
	expected := []float64{2.0, 3.0, 5.0}
	for i, tier := range StormTiers {
		if tier.MultBoost != expected[i] {
			t.Errorf("StormTiers[%d].MultBoost 應該是 %.1f，得到 %.1f", i, expected[i], tier.MultBoost)
		}
	}
}
