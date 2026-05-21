package fragment

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
}

func TestEnsureAndRemovePlayer(t *testing.T) {
	m := New()
	m.EnsurePlayer("p1")
	snap := m.GetSnapshot("p1")
	if snap.Bronze != 0 || snap.Silver != 0 || snap.Gold != 0 {
		t.Error("new player should have 0 fragments")
	}
	m.RemovePlayer("p1")
	snap2 := m.GetSnapshot("p1")
	if snap2.Bronze != 0 {
		t.Error("removed player should return empty snapshot")
	}
}

func TestTryDrop_BossAlwaysDropsGold(t *testing.T) {
	m := New()
	// 強制多次嘗試，BOSS 50% 機率，100次應該至少有一次
	dropped := false
	for i := 0; i < 100; i++ {
		r := m.TryDrop("p1", "B001", 100, true)
		if r.Dropped {
			dropped = true
			if r.FragmentType != FragmentGold {
				t.Errorf("BOSS should drop gold, got %s", r.FragmentType)
			}
			break
		}
	}
	if !dropped {
		t.Error("BOSS should drop gold fragment within 100 tries")
	}
}

func TestTryDrop_NormalTargetDropsBronze(t *testing.T) {
	m := New()
	dropped := false
	for i := 0; i < 200; i++ {
		r := m.TryDrop("p1", "T001", 100, false)
		if r.Dropped {
			dropped = true
			if r.FragmentType != FragmentBronze {
				t.Errorf("normal target should drop bronze, got %s", r.FragmentType)
			}
			break
		}
	}
	if !dropped {
		t.Error("normal target should drop bronze fragment within 200 tries")
	}
}

func TestTryDrop_SpecialTargetDropsSilver(t *testing.T) {
	m := New()
	dropped := false
	for i := 0; i < 100; i++ {
		r := m.TryDrop("p1", "T103", 100, false)
		if r.Dropped {
			dropped = true
			if r.FragmentType != FragmentSilver {
				t.Errorf("special target should drop silver, got %s", r.FragmentType)
			}
			break
		}
	}
	if !dropped {
		t.Error("special target should drop silver fragment within 100 tries")
	}
}

func TestCollectBronzeComplete(t *testing.T) {
	m := New()
	m.EnsurePlayer("p1")

	// 手動填充到 4 個
	m.mu.Lock()
	m.players["p1"].Bronze = 4
	m.mu.Unlock()

	// 再掉落一個銅碎片（強制）
	var result *DropResult
	for i := 0; i < 500; i++ {
		r := m.TryDrop("p1", "T001", 100, false)
		if r.Dropped && r.FragmentType == FragmentBronze {
			result = r
			break
		}
	}
	if result == nil {
		t.Skip("could not force bronze drop in 500 tries")
	}
	if !result.IsComplete {
		t.Error("should be complete after 5 bronze fragments")
	}
	if result.Reward != 100*30 {
		t.Errorf("bronze reward should be 3000, got %d", result.Reward)
	}
	// 碎片應重置
	snap := m.GetSnapshot("p1")
	if snap.Bronze != 0 {
		t.Errorf("bronze should reset to 0 after complete, got %d", snap.Bronze)
	}
}

func TestCollectGoldComplete(t *testing.T) {
	m := New()
	m.EnsurePlayer("p1")

	// 手動填充到 4 個金碎片
	m.mu.Lock()
	m.players["p1"].Gold = 4
	m.mu.Unlock()

	// 再掉落一個金碎片（BOSS 50% 機率）
	var result *DropResult
	for i := 0; i < 100; i++ {
		r := m.TryDrop("p1", "B001", 500, true)
		if r.Dropped && r.FragmentType == FragmentGold {
			result = r
			break
		}
	}
	if result == nil {
		t.Skip("could not force gold drop in 100 tries")
	}
	if !result.IsComplete {
		t.Error("should be complete after 5 gold fragments")
	}
	if result.Reward != 500*200 {
		t.Errorf("gold reward should be 100000, got %d", result.Reward)
	}
}

func TestGetRewardDef(t *testing.T) {
	def := GetRewardDef(FragmentGold)
	if def.Required != 5 {
		t.Errorf("gold required should be 5, got %d", def.Required)
	}
	if def.RewardMult != 200 {
		t.Errorf("gold reward mult should be 200, got %d", def.RewardMult)
	}
}

func TestGetAllRewardDefs(t *testing.T) {
	defs := GetAllRewardDefs()
	if len(defs) != 3 {
		t.Errorf("should have 3 reward defs, got %d", len(defs))
	}
}

func TestGetSnapshot_Unknown(t *testing.T) {
	m := New()
	snap := m.GetSnapshot("unknown")
	if snap.Bronze != 0 || snap.Silver != 0 || snap.Gold != 0 {
		t.Error("unknown player should return empty snapshot")
	}
}

func TestCountIncrement(t *testing.T) {
	m := New()
	m.EnsurePlayer("p1")

	// 強制掉落 3 個銅碎片，確認計數遞增
	count := 0
	for i := 0; i < 1000 && count < 3; i++ {
		r := m.TryDrop("p1", "T001", 100, false)
		if r.Dropped && r.FragmentType == FragmentBronze && !r.IsComplete {
			count++
			if r.NewCount != count {
				t.Errorf("expected count %d, got %d", count, r.NewCount)
			}
		}
	}
}
