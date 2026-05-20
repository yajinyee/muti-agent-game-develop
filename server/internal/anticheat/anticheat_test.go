package anticheat

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
}

func TestEnsureAndRemoveRecord(t *testing.T) {
	m := New()
	m.EnsureRecord("p1", "Player1")
	m.mu.RLock()
	_, ok := m.records["p1"]
	m.mu.RUnlock()
	if !ok {
		t.Fatal("record not created")
	}
	m.RemoveRecord("p1")
	m.mu.RLock()
	_, ok = m.records["p1"]
	m.mu.RUnlock()
	if ok {
		t.Fatal("record not removed")
	}
}

func TestRecordAttack_NormalRate(t *testing.T) {
	m := New()
	m.EnsureRecord("p1", "Player1")
	// 正常攻擊頻率（每秒 2 次）
	for i := 0; i < 5; i++ {
		alert := m.RecordAttack("p1", 100)
		if alert != nil {
			t.Errorf("unexpected alert for normal attack rate: %s", alert.Message)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func TestRecordAttack_BotRate(t *testing.T) {
	m := New()
	m.EnsureRecord("p1", "Player1")
	// 模擬 bot 攻擊（快速連續攻擊超過 80 次，超過 8次/秒 × 10秒 門檻）
	for i := 0; i < 85; i++ {
		m.RecordAttack("p1", 100)
	}
	// 應該觸發警告
	total, _, _ := m.GetAlertCount()
	if total == 0 {
		t.Error("expected bot attack alert, got none")
	}
}

func TestRecordReward_HighRTP(t *testing.T) {
	m := New()
	m.EnsureRecord("p1", "Player1")
	// 模擬 100 次攻擊，每次 bet=100
	for i := 0; i < 100; i++ {
		m.RecordAttack("p1", 100)
	}
	// 模擬超高 RTP（bet=10000, reward=30000 → RTP=300%）
	m.mu.Lock()
	m.records["p1"].TotalBet = 10000
	m.records["p1"].TotalReward = 30000
	m.records["p1"].AttackCount = 100
	m.mu.Unlock()
	alert := m.RecordReward("p1", 0)
	if alert == nil {
		t.Error("expected high RTP alert, got none")
	}
	if alert != nil && alert.Type != AlertHighRTP {
		t.Errorf("expected AlertHighRTP, got %s", alert.Type)
	}
}

func TestRecordReward_NormalRTP(t *testing.T) {
	m := New()
	m.EnsureRecord("p1", "Player1")
	// 正常 RTP（94%）
	m.mu.Lock()
	m.records["p1"].TotalBet = 10000
	m.records["p1"].TotalReward = 9400
	m.records["p1"].AttackCount = 100
	m.mu.Unlock()
	alert := m.RecordReward("p1", 0)
	if alert != nil {
		t.Errorf("unexpected alert for normal RTP: %s", alert.Message)
	}
}

func TestRecordCoins_Spike(t *testing.T) {
	m := New()
	m.EnsureRecord("p1", "Player1")
	// 模擬金幣暴增
	m.RecordCoins("p1", 1000)
	m.RecordCoins("p1", 60000) // 增加 59000，超過門檻 50000
	total, _, _ := m.GetAlertCount()
	if total == 0 {
		t.Error("expected coin spike alert, got none")
	}
}

func TestRecordBonus_Abuse(t *testing.T) {
	m := New()
	m.EnsureRecord("p1", "Player1")
	// 觸發 6 次 Bonus（超過門檻 5 次）
	for i := 0; i < 6; i++ {
		m.RecordBonus("p1")
	}
	total, _, _ := m.GetAlertCount()
	if total == 0 {
		t.Error("expected bonus abuse alert, got none")
	}
}

func TestRecordJackpot_Abuse(t *testing.T) {
	m := New()
	m.EnsureRecord("p1", "Player1")
	// 觸發 4 次 Jackpot（超過門檻 3 次）
	for i := 0; i < 4; i++ {
		m.RecordJackpot("p1")
	}
	total, critical, _ := m.GetAlertCount()
	if total == 0 {
		t.Error("expected jackpot abuse alert, got none")
	}
	if critical == 0 {
		t.Error("expected critical level for jackpot abuse")
	}
}

func TestAlertCooldown(t *testing.T) {
	m := New()
	m.EnsureRecord("p1", "Player1")
	// 第一次觸發
	for i := 0; i < 6; i++ {
		m.RecordBonus("p1")
	}
	count1, _, _ := m.GetAlertCount()
	// 立即再次觸發（應該被冷卻時間阻擋）
	for i := 0; i < 6; i++ {
		m.RecordBonus("p1")
	}
	count2, _, _ := m.GetAlertCount()
	if count2 != count1 {
		t.Errorf("cooldown not working: count1=%d, count2=%d", count1, count2)
	}
}

func TestGetAlerts(t *testing.T) {
	m := New()
	m.EnsureRecord("p1", "Player1")
	// 觸發多個警告
	for i := 0; i < 4; i++ {
		m.RecordJackpot("p1")
	}
	alerts := m.GetAlerts(10)
	if len(alerts) == 0 {
		t.Error("expected alerts, got none")
	}
	// 確認是倒序（最新的在前）
	if len(alerts) > 1 {
		if alerts[0].CreatedAt.Before(alerts[1].CreatedAt) {
			t.Error("alerts not in reverse chronological order")
		}
	}
}

func TestGetPlayerRTP(t *testing.T) {
	m := New()
	m.EnsureRecord("p1", "Player1")
	m.mu.Lock()
	m.records["p1"].TotalBet = 10000
	m.records["p1"].TotalReward = 9400
	m.mu.Unlock()
	rtp := m.GetPlayerRTP("p1")
	if rtp < 0.93 || rtp > 0.95 {
		t.Errorf("unexpected RTP: %.4f", rtp)
	}
}
