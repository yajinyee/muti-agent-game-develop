// fevermode_test.go — 狂熱模式系統單元測試（DAY-133）
package fevermode

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() 回傳 nil")
	}
	if m.config.TriggerKills != 5 {
		t.Errorf("TriggerKills 應為 5，得到 %d", m.config.TriggerKills)
	}
}

func TestRecordKill_NoTrigger(t *testing.T) {
	m := New()
	// 擊破 4 次，不應觸發
	for i := 0; i < 4; i++ {
		triggered, extended, mult := m.RecordKill("p1")
		if triggered {
			t.Errorf("第 %d 次擊破不應觸發狂熱", i+1)
		}
		if extended {
			t.Errorf("第 %d 次擊破不應延長狂熱", i+1)
		}
		if mult != 1.0 {
			t.Errorf("未觸發時倍率應為 1.0，得到 %.1f", mult)
		}
	}
}

func TestRecordKill_Trigger(t *testing.T) {
	m := New()
	// 擊破 5 次，第 5 次應觸發
	for i := 0; i < 4; i++ {
		m.RecordKill("p1")
	}
	triggered, extended, mult := m.RecordKill("p1")
	if !triggered {
		t.Error("第 5 次擊破應觸發狂熱")
	}
	if extended {
		t.Error("觸發時不應同時延長")
	}
	if mult != m.config.MultBoost {
		t.Errorf("觸發後倍率應為 %.1f，得到 %.1f", m.config.MultBoost, mult)
	}
}

func TestRecordKill_ExtendDuringFever(t *testing.T) {
	m := New()
	// 觸發狂熱
	for i := 0; i < 5; i++ {
		m.RecordKill("p1")
	}
	// 狂熱中再擊破，應延長
	triggered, extended, mult := m.RecordKill("p1")
	if triggered {
		t.Error("狂熱中不應再次觸發")
	}
	if !extended {
		t.Error("狂熱中擊破應延長時間")
	}
	if mult != m.config.MultBoost {
		t.Errorf("狂熱中倍率應為 %.1f，得到 %.1f", m.config.MultBoost, mult)
	}
}

func TestGetMultBoost_Active(t *testing.T) {
	m := New()
	for i := 0; i < 5; i++ {
		m.RecordKill("p1")
	}
	boost := m.GetMultBoost("p1")
	if boost != m.config.MultBoost {
		t.Errorf("狂熱中倍率應為 %.1f，得到 %.1f", m.config.MultBoost, boost)
	}
}

func TestGetMultBoost_Idle(t *testing.T) {
	m := New()
	boost := m.GetMultBoost("p1")
	if boost != 1.0 {
		t.Errorf("未觸發時倍率應為 1.0，得到 %.1f", boost)
	}
}

func TestGetSnapshot_Idle(t *testing.T) {
	m := New()
	snap := m.GetSnapshot("p1")
	if snap.State != FeverStateIdle {
		t.Errorf("初始狀態應為 Idle，得到 %d", snap.State)
	}
	if snap.MultBoost != 1.0 {
		t.Errorf("初始倍率應為 1.0，得到 %.1f", snap.MultBoost)
	}
}

func TestGetSnapshot_Active(t *testing.T) {
	m := New()
	for i := 0; i < 5; i++ {
		m.RecordKill("p1")
	}
	snap := m.GetSnapshot("p1")
	if snap.State != FeverStateActive {
		t.Errorf("觸發後狀態應為 Active，得到 %d", snap.State)
	}
	if snap.SecondsLeft <= 0 {
		t.Error("狂熱中剩餘秒數應 > 0")
	}
	if snap.MultBoost != m.config.MultBoost {
		t.Errorf("狂熱中倍率應為 %.1f，得到 %.1f", m.config.MultBoost, snap.MultBoost)
	}
}

func TestGetSnapshot_KillProgress(t *testing.T) {
	m := New()
	m.RecordKill("p1")
	m.RecordKill("p1")
	snap := m.GetSnapshot("p1")
	if snap.KillProgress != 2 {
		t.Errorf("擊破進度應為 2，得到 %d", snap.KillProgress)
	}
}

func TestCheckExpiry(t *testing.T) {
	m := New()
	// 觸發狂熱
	for i := 0; i < 5; i++ {
		m.RecordKill("p1")
	}
	// 手動設定結束時間為過去
	m.mu.Lock()
	p := m.players["p1"]
	p.FeverEndAt = time.Now().Add(-1 * time.Second)
	m.mu.Unlock()

	expired := m.CheckExpiry("p1")
	if !expired {
		t.Error("應該偵測到狂熱結束")
	}
	snap := m.GetSnapshot("p1")
	if snap.State != FeverStateCooldown {
		t.Errorf("結束後應進入冷卻，得到 %d", snap.State)
	}
}

func TestTickExpiry(t *testing.T) {
	m := New()
	// 觸發兩個玩家的狂熱
	for i := 0; i < 5; i++ {
		m.RecordKill("p1")
		m.RecordKill("p2")
	}
	// 手動設定 p1 結束時間為過去
	m.mu.Lock()
	m.players["p1"].FeverEndAt = time.Now().Add(-1 * time.Second)
	m.mu.Unlock()

	expired := m.TickExpiry()
	if len(expired) != 1 || expired[0] != "p1" {
		t.Errorf("應該只有 p1 過期，得到 %v", expired)
	}
}

func TestRemovePlayer(t *testing.T) {
	m := New()
	for i := 0; i < 5; i++ {
		m.RecordKill("p1")
	}
	m.RemovePlayer("p1")
	snap := m.GetSnapshot("p1")
	if snap.State != FeverStateIdle {
		t.Error("移除後應回到 Idle 狀態")
	}
}

func TestMultiplePlayers(t *testing.T) {
	m := New()
	// p1 觸發狂熱
	for i := 0; i < 5; i++ {
		m.RecordKill("p1")
	}
	// p2 不受影響
	boost := m.GetMultBoost("p2")
	if boost != 1.0 {
		t.Errorf("p2 不應受 p1 影響，得到 %.1f", boost)
	}
}

func TestCooldownAfterFever(t *testing.T) {
	m := New()
	// 觸發狂熱
	for i := 0; i < 5; i++ {
		m.RecordKill("p1")
	}
	// 手動設定結束時間為過去
	m.mu.Lock()
	m.players["p1"].FeverEndAt = time.Now().Add(-1 * time.Second)
	m.mu.Unlock()
	m.CheckExpiry("p1")

	// 冷卻中擊破不應觸發
	triggered, _, _ := m.RecordKill("p1")
	if triggered {
		t.Error("冷卻中不應觸發狂熱")
	}
}

func TestTotalFevered(t *testing.T) {
	m := New()
	// 觸發狂熱
	for i := 0; i < 5; i++ {
		m.RecordKill("p1")
	}
	snap := m.GetSnapshot("p1")
	if snap.TotalFevered != 1 {
		t.Errorf("觸發次數應為 1，得到 %d", snap.TotalFevered)
	}
}

func TestWindowExpiry(t *testing.T) {
	m := New()
	// 手動加入過期的擊破記錄
	m.mu.Lock()
	p := m.getOrCreate("p1")
	// 加入 4 個超出時間窗口的記錄
	old := time.Now().Add(-10 * time.Second)
	for i := 0; i < 4; i++ {
		p.KillTimes = append(p.KillTimes, old)
	}
	m.mu.Unlock()

	// 再擊破 5 次（窗口內），應觸發
	for i := 0; i < 4; i++ {
		m.RecordKill("p1")
	}
	triggered, _, _ := m.RecordKill("p1")
	if !triggered {
		t.Error("窗口內 5 次擊破應觸發狂熱（舊記錄應被清除）")
	}
}
