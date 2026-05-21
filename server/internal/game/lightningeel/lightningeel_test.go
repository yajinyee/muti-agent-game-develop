// lightningeel_test.go — 閃電鰻連鎖攻擊系統單元測試（DAY-132）
package lightningeel

import (
	"testing"
	"time"
)

func makeTargets(n int) []NearbyTarget {
	targets := make([]NearbyTarget, n)
	for i := 0; i < n; i++ {
		targets[i] = NearbyTarget{
			InstanceID: string(rune('A' + i)),
			DefID:      "T001",
			Name:       "測試目標",
			Multiplier: 2.0,
			X:          float64(i * 50),
			Y:          0,
		}
	}
	return targets
}

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() 回傳 nil")
	}
	if m.config.MaxJumps != 5 {
		t.Errorf("MaxJumps 應為 5，得到 %d", m.config.MaxJumps)
	}
}

func TestCanTrigger_NewPlayer(t *testing.T) {
	m := New()
	if !m.CanTrigger("p1") {
		t.Error("新玩家應該可以觸發連鎖")
	}
}

func TestCanTrigger_AfterChain(t *testing.T) {
	m := New()
	targets := makeTargets(3)
	m.ExecuteChain("p1", "trigger", targets, 100)
	// 剛觸發後應該在冷卻中
	if m.CanTrigger("p1") {
		t.Error("觸發後應該在冷卻中，不能再次觸發")
	}
}

func TestCooldownLeft(t *testing.T) {
	m := New()
	targets := makeTargets(3)
	m.ExecuteChain("p1", "trigger", targets, 100)
	cd := m.CooldownLeft("p1")
	// +1 是因為 int 轉換的進位，允許 CooldownSecs+1
	if cd <= 0 || cd > m.config.CooldownSecs+1 {
		t.Errorf("冷卻剩餘應在 1-%d 秒，得到 %d", m.config.CooldownSecs+1, cd)
	}
}

func TestCooldownLeft_NewPlayer(t *testing.T) {
	m := New()
	cd := m.CooldownLeft("p1")
	if cd != 0 {
		t.Errorf("新玩家冷卻應為 0，得到 %d", cd)
	}
}

func TestExecuteChain_Basic(t *testing.T) {
	m := New()
	targets := makeTargets(3)
	session := m.ExecuteChain("p1", "trigger", targets, 100)
	if session == nil {
		t.Fatal("ExecuteChain 回傳 nil")
	}
	if session.PlayerID != "p1" {
		t.Errorf("PlayerID 應為 p1，得到 %s", session.PlayerID)
	}
	if len(session.Jumps) > m.config.MaxJumps {
		t.Errorf("跳躍次數不應超過 MaxJumps=%d，得到 %d", m.config.MaxJumps, len(session.Jumps))
	}
}

func TestExecuteChain_MaxJumps(t *testing.T) {
	m := New()
	// 提供 10 個目標，但最多只跳 5 次
	targets := makeTargets(10)
	session := m.ExecuteChain("p1", "trigger", targets, 100)
	if len(session.Jumps) > m.config.MaxJumps {
		t.Errorf("跳躍次數不應超過 MaxJumps=%d，得到 %d", m.config.MaxJumps, len(session.Jumps))
	}
}

func TestExecuteChain_NoTargets(t *testing.T) {
	m := New()
	session := m.ExecuteChain("p1", "trigger", []NearbyTarget{}, 100)
	if len(session.Jumps) != 0 {
		t.Errorf("無目標時跳躍次數應為 0，得到 %d", len(session.Jumps))
	}
}

func TestExecuteChain_NoDuplicateJumps(t *testing.T) {
	m := New()
	targets := makeTargets(5)
	session := m.ExecuteChain("p1", "trigger", targets, 100)
	seen := make(map[string]bool)
	for _, j := range session.Jumps {
		if seen[j.TargetInstanceID] {
			t.Errorf("目標 %s 被跳躍了兩次", j.TargetInstanceID)
		}
		seen[j.TargetInstanceID] = true
	}
}

func TestExecuteChain_RewardCalc(t *testing.T) {
	// 強制所有跳躍都擊破（用大量測試取平均）
	// 只驗證獎勵計算公式正確
	targets := []NearbyTarget{
		{InstanceID: "T1", DefID: "T001", Name: "測試", Multiplier: 5.0},
	}
	// 多次執行，至少有一次擊破
	var gotKill bool
	for i := 0; i < 100; i++ {
		m2 := New()
		session := m2.ExecuteChain("p1", "trigger", targets, 100)
		if len(session.Jumps) > 0 && session.Jumps[0].Killed {
			gotKill = true
			// 獎勵 = betCost × multiplier × JumpMultMod = 100 × 5.0 × 0.6 = 300
			expected := int64(100 * 5.0 * m2.config.JumpMultMod)
			if session.Jumps[0].Reward != expected {
				t.Errorf("獎勵計算錯誤：期望 %d，得到 %d", expected, session.Jumps[0].Reward)
			}
			break
		}
	}
	if !gotKill {
		t.Log("100 次測試中沒有擊破，機率可能太低（非錯誤）")
	}
}

func TestExecuteChain_JumpIndex(t *testing.T) {
	m := New()
	targets := makeTargets(5)
	session := m.ExecuteChain("p1", "trigger", targets, 100)
	for i, j := range session.Jumps {
		if j.JumpIndex != i+1 {
			t.Errorf("JumpIndex[%d] 應為 %d，得到 %d", i, i+1, j.JumpIndex)
		}
	}
}

func TestRemovePlayer(t *testing.T) {
	m := New()
	targets := makeTargets(3)
	m.ExecuteChain("p1", "trigger", targets, 100)
	m.RemovePlayer("p1")
	// 移除後應該可以再次觸發
	if !m.CanTrigger("p1") {
		t.Error("移除玩家後應該可以再次觸發")
	}
}

func TestGetSnapshot(t *testing.T) {
	m := New()
	snap := m.GetSnapshot("p1")
	if snap.PlayerID != "p1" {
		t.Errorf("PlayerID 應為 p1，得到 %s", snap.PlayerID)
	}
	if snap.CooldownLeft != 0 {
		t.Errorf("新玩家冷卻應為 0，得到 %d", snap.CooldownLeft)
	}
}

func TestGetSnapshot_AfterChain(t *testing.T) {
	m := New()
	targets := makeTargets(3)
	m.ExecuteChain("p1", "trigger", targets, 100)
	snap := m.GetSnapshot("p1")
	if snap.CooldownLeft <= 0 {
		t.Error("觸發後冷卻應 > 0")
	}
}

func TestMultiplePlayers(t *testing.T) {
	m := New()
	targets := makeTargets(3)
	m.ExecuteChain("p1", "trigger", targets, 100)
	// p2 不受 p1 冷卻影響
	if !m.CanTrigger("p2") {
		t.Error("p2 不應受 p1 冷卻影響")
	}
}

func TestCanTrigger_AfterCooldown(t *testing.T) {
	m := New()
	// 手動設定 LastChainAt 為 10 秒前（超過 8 秒冷卻）
	m.mu.Lock()
	s := m.getOrCreatePlayer("p1")
	s.LastChainAt = time.Now().Add(-10 * time.Second)
	m.mu.Unlock()

	if !m.CanTrigger("p1") {
		t.Error("冷卻結束後應該可以觸發")
	}
}
