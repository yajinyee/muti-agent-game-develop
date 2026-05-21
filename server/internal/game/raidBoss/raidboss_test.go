// raidboss_test.go — Co-op Boss Raid 單元測試（DAY-115）
package raidboss

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := New()
	if m.GetState() != RaidStateIdle {
		t.Errorf("expected idle, got %s", m.GetState())
	}
}

func TestCanTrigger(t *testing.T) {
	m := New()
	if !m.CanTrigger("2026-05-21") {
		t.Error("should be able to trigger on fresh manager")
	}
}

func TestCanTrigger_SameDay(t *testing.T) {
	m := New()
	m.StartWarning()
	m.StartRaid("超強 BOSS", 10000, 50000, "2026-05-21")
	if m.CanTrigger("2026-05-21") {
		t.Error("should not trigger twice on same day")
	}
}

func TestCanTrigger_NextDay(t *testing.T) {
	m := New()
	m.StartWarning()
	m.StartRaid("超強 BOSS", 10000, 50000, "2026-05-21")
	m.Reset()
	if !m.CanTrigger("2026-05-22") {
		t.Error("should be able to trigger on next day")
	}
}

func TestStartWarning(t *testing.T) {
	m := New()
	raidID := m.StartWarning()
	if raidID == "" {
		t.Error("raidID should not be empty")
	}
	if m.GetState() != RaidStateWarning {
		t.Errorf("expected warning, got %s", m.GetState())
	}
}

func TestStartRaid(t *testing.T) {
	m := New()
	m.StartWarning()
	m.StartRaid("超強 BOSS", 10000, 50000, "2026-05-21")
	snap := m.GetSnapshot()
	if snap.State != RaidStateActive {
		t.Errorf("expected active, got %s", snap.State)
	}
	if snap.HP != 10000 {
		t.Errorf("expected HP=10000, got %d", snap.HP)
	}
	if snap.RewardPool != 50000 {
		t.Errorf("expected RewardPool=50000, got %d", snap.RewardPool)
	}
}

func TestRecordDamage_Kill(t *testing.T) {
	m := New()
	m.StartWarning()
	m.StartRaid("超強 BOSS", 1000, 10000, "2026-05-21")

	// 玩家 A 打 600 傷害
	hp, killed := m.RecordDamage("playerA", "玩家A", 600)
	if killed {
		t.Error("should not be killed yet")
	}
	if hp != 400 {
		t.Errorf("expected HP=400, got %d", hp)
	}

	// 玩家 B 打 400 傷害（擊殺）
	hp, killed = m.RecordDamage("playerB", "玩家B", 400)
	if !killed {
		t.Error("should be killed")
	}
	if hp != 0 {
		t.Errorf("expected HP=0, got %d", hp)
	}
	if m.GetState() != RaidStateResult {
		t.Errorf("expected result state, got %s", m.GetState())
	}
}

func TestRewardDistribution(t *testing.T) {
	m := New()
	m.StartWarning()
	m.StartRaid("超強 BOSS", 1000, 10000, "2026-05-21")

	// A 打 600，B 打 400
	m.RecordDamage("playerA", "玩家A", 600)
	m.RecordDamage("playerB", "玩家B", 400)

	rewardA, rankA, okA := m.GetContributorReward("playerA")
	rewardB, rankB, okB := m.GetContributorReward("playerB")

	if !okA || !okB {
		t.Error("both players should have rewards")
	}
	if rankA != 1 {
		t.Errorf("A should be rank 1, got %d", rankA)
	}
	if rankB != 2 {
		t.Errorf("B should be rank 2, got %d", rankB)
	}
	// A 應得 6000，B 應得 4000
	if rewardA != 6000 {
		t.Errorf("A reward expected 6000, got %d", rewardA)
	}
	if rewardB != 4000 {
		t.Errorf("B reward expected 4000, got %d", rewardB)
	}
	// 總和應等於獎勵池
	if rewardA+rewardB != 10000 {
		t.Errorf("total reward should be 10000, got %d", rewardA+rewardB)
	}
}

func TestCheckTimeout(t *testing.T) {
	m := New()
	m.StartWarning()
	m.StartRaid("超強 BOSS", 99999, 10000, "2026-05-21")

	// 手動設定 endsAt 為過去
	m.mu.Lock()
	m.endsAt = time.Now().Add(-1 * time.Second)
	m.mu.Unlock()

	timedOut := m.CheckTimeout()
	if !timedOut {
		t.Error("should have timed out")
	}
	if m.GetState() != RaidStateResult {
		t.Errorf("expected result state after timeout, got %s", m.GetState())
	}
}

func TestRecordDamage_NotActive(t *testing.T) {
	m := New()
	hp, killed := m.RecordDamage("playerA", "玩家A", 100)
	if killed {
		t.Error("should not kill when not active")
	}
	if hp != 0 {
		t.Errorf("expected HP=0 when not active, got %d", hp)
	}
}

func TestGetSnapshot_Contributors(t *testing.T) {
	m := New()
	m.StartWarning()
	m.StartRaid("超強 BOSS", 10000, 50000, "2026-05-21")

	m.RecordDamage("playerA", "玩家A", 300)
	m.RecordDamage("playerB", "玩家B", 700)

	snap := m.GetSnapshot()
	if len(snap.Contributors) != 2 {
		t.Errorf("expected 2 contributors, got %d", len(snap.Contributors))
	}
	// B 傷害更高，應排第一
	if snap.Contributors[0].PlayerID != "playerB" {
		t.Errorf("B should be first contributor")
	}
}

func TestReset(t *testing.T) {
	m := New()
	m.StartWarning()
	m.StartRaid("超強 BOSS", 1000, 10000, "2026-05-21")
	m.Reset()
	if m.GetState() != RaidStateIdle {
		t.Errorf("expected idle after reset, got %s", m.GetState())
	}
}

func TestIsActive(t *testing.T) {
	m := New()
	if m.IsActive() {
		t.Error("should not be active initially")
	}
	m.StartWarning()
	m.StartRaid("超強 BOSS", 1000, 10000, "2026-05-21")
	if !m.IsActive() {
		t.Error("should be active after start")
	}
}
