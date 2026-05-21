// Package bounty 懸賞系統單元測試
package bounty

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := NewDefault()
	if m == nil {
		t.Fatal("NewDefault() returned nil")
	}
}

func TestCanPost_Initial(t *testing.T) {
	m := NewDefault()
	ok, remaining := m.CanPost("p1")
	if !ok {
		t.Error("初始狀態應該可以下懸賞")
	}
	if remaining != 0 {
		t.Errorf("初始冷卻應該是 0，得到 %d", remaining)
	}
}

func TestPostBounty_Basic(t *testing.T) {
	m := NewDefault()
	id, errCode := m.PostBounty("p1", "玩家1", "inst1", "T103", "流星", 15.0, 500)
	if errCode != "" {
		t.Errorf("下懸賞應該成功，得到錯誤: %s", errCode)
	}
	if id == "" {
		t.Error("bountyID 不應該是空字串")
	}
}

func TestPostBounty_InvalidAmount_TooLow(t *testing.T) {
	m := NewDefault()
	_, errCode := m.PostBounty("p1", "玩家1", "inst1", "T103", "流星", 15.0, 50) // 低於 MinBountyAmount=100
	if errCode != "invalid_amount" {
		t.Errorf("金額太低應該回傳 invalid_amount，得到: %s", errCode)
	}
}

func TestPostBounty_InvalidAmount_TooHigh(t *testing.T) {
	m := NewDefault()
	_, errCode := m.PostBounty("p1", "玩家1", "inst1", "T103", "流星", 15.0, 99999) // 高於 MaxBountyAmount=5000
	if errCode != "invalid_amount" {
		t.Errorf("金額太高應該回傳 invalid_amount，得到: %s", errCode)
	}
}

func TestPostBounty_Cooldown(t *testing.T) {
	m := NewDefault()
	m.PostBounty("p1", "玩家1", "inst1", "T103", "流星", 15.0, 500)
	// 同一玩家再次下懸賞應該被冷卻
	_, errCode := m.PostBounty("p1", "玩家1", "inst2", "T104", "金草", 20.0, 500)
	if errCode != "cooldown" {
		t.Errorf("冷卻中應該回傳 cooldown，得到: %s", errCode)
	}
}

func TestPostBounty_Full(t *testing.T) {
	m := NewDefault()
	// 填滿 3 個懸賞（不同玩家）
	m.PostBounty("p1", "玩家1", "inst1", "T103", "流星", 15.0, 500)
	m.PostBounty("p2", "玩家2", "inst2", "T104", "金草", 20.0, 500)
	m.PostBounty("p3", "玩家3", "inst3", "T105", "金幣魚", 25.0, 500)
	// 第四個應該失敗
	_, errCode := m.PostBounty("p4", "玩家4", "inst4", "T101", "擬態怪", 30.0, 500)
	if errCode != "full" {
		t.Errorf("懸賞已滿應該回傳 full，得到: %s", errCode)
	}
}

func TestClaimBounty_Basic(t *testing.T) {
	m := NewDefault()
	m.PostBounty("p1", "玩家1", "inst1", "T103", "流星", 15.0, 500)
	total, claimed, isAny := m.ClaimBounty("inst1", "p2", "玩家2")
	if !isAny {
		t.Error("應該有懸賞可領取")
	}
	if total != 500 {
		t.Errorf("懸賞金額應該是 500，得到 %d", total)
	}
	if len(claimed) != 1 {
		t.Errorf("應該領取 1 筆懸賞，得到 %d", len(claimed))
	}
}

func TestClaimBounty_NoMatch(t *testing.T) {
	m := NewDefault()
	m.PostBounty("p1", "玩家1", "inst1", "T103", "流星", 15.0, 500)
	total, claimed, isAny := m.ClaimBounty("inst2", "p2", "玩家2") // 不同 instanceID
	if isAny {
		t.Error("不同 instanceID 不應該有懸賞")
	}
	if total != 0 || len(claimed) != 0 {
		t.Error("不同 instanceID 應該回傳 0 和空列表")
	}
}

func TestClaimBounty_MultipleBounties(t *testing.T) {
	m := NewDefault()
	// 兩個玩家對同一目標下懸賞
	m.PostBounty("p1", "玩家1", "inst1", "T103", "流星", 15.0, 500)
	m.PostBounty("p2", "玩家2", "inst1", "T103", "流星", 15.0, 300)
	total, claimed, isAny := m.ClaimBounty("inst1", "p3", "玩家3")
	if !isAny {
		t.Error("應該有懸賞可領取")
	}
	if total != 800 {
		t.Errorf("總懸賞金額應該是 800，得到 %d", total)
	}
	if len(claimed) != 2 {
		t.Errorf("應該領取 2 筆懸賞，得到 %d", len(claimed))
	}
}

func TestCheckExpiry(t *testing.T) {
	cfg := DefaultConfig()
	cfg.BountyDuration = 0.01 // 10ms，快速超時
	m := New(cfg)
	m.PostBounty("p1", "玩家1", "inst1", "T103", "流星", 15.0, 500)
	time.Sleep(20 * time.Millisecond)
	expired := m.CheckExpiry()
	if len(expired) != 1 {
		t.Errorf("應該有 1 筆過期懸賞，得到 %d", len(expired))
	}
}

func TestCancelBountyForTarget(t *testing.T) {
	m := NewDefault()
	m.PostBounty("p1", "玩家1", "inst1", "T103", "流星", 15.0, 500)
	cancelled := m.CancelBountyForTarget("inst1")
	if len(cancelled) != 1 {
		t.Errorf("應該取消 1 筆懸賞，得到 %d", len(cancelled))
	}
	// 確認懸賞已取消
	bounties := m.GetActiveBounties()
	if len(bounties) != 0 {
		t.Error("取消後不應該有活躍懸賞")
	}
}

func TestGetActiveBounties(t *testing.T) {
	m := NewDefault()
	m.PostBounty("p1", "玩家1", "inst1", "T103", "流星", 15.0, 500)
	m.PostBounty("p2", "玩家2", "inst2", "T104", "金草", 20.0, 300)
	bounties := m.GetActiveBounties()
	if len(bounties) != 2 {
		t.Errorf("應該有 2 筆活躍懸賞，得到 %d", len(bounties))
	}
}

func TestGetBountiesForTarget(t *testing.T) {
	m := NewDefault()
	m.PostBounty("p1", "玩家1", "inst1", "T103", "流星", 15.0, 500)
	m.PostBounty("p2", "玩家2", "inst1", "T103", "流星", 15.0, 300)
	m.PostBounty("p3", "玩家3", "inst2", "T104", "金草", 20.0, 200)
	bounties := m.GetBountiesForTarget("inst1")
	if len(bounties) != 2 {
		t.Errorf("inst1 應該有 2 筆懸賞，得到 %d", len(bounties))
	}
}

func TestGetPlayerCooldown(t *testing.T) {
	m := NewDefault()
	m.PostBounty("p1", "玩家1", "inst1", "T103", "流星", 15.0, 500)
	cooldown := m.GetPlayerCooldown("p1")
	if cooldown <= 0 {
		t.Error("下懸賞後應該有冷卻")
	}
}

func TestGetPlayerCooldown_NoPost(t *testing.T) {
	m := NewDefault()
	cooldown := m.GetPlayerCooldown("p1")
	if cooldown != 0 {
		t.Errorf("未下懸賞的玩家冷卻應該是 0，得到 %d", cooldown)
	}
}
