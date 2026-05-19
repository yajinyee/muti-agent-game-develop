// Package mission 每日任務系統單元測試（DAY-037）
package mission

import (
	"testing"
)

// TestNewManager 確認 Manager 正確初始化
func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("NewManager returned nil")
	}
	if m.progress == nil {
		t.Error("progress map should be initialized")
	}
	if m.resetAt.IsZero() {
		t.Error("resetAt should be set")
	}
}

// TestGetOrInitProgress 確認玩家進度正確初始化
func TestGetOrInitProgress(t *testing.T) {
	m := NewManager()
	progress := m.GetOrInitProgress("player-001")

	if len(progress) != len(DailyMissions) {
		t.Errorf("expected %d missions, got %d", len(DailyMissions), len(progress))
	}

	// 確認每個任務都有初始進度
	for _, mission := range DailyMissions {
		prog, ok := progress[mission.ID]
		if !ok {
			t.Errorf("mission %s not found in progress", mission.ID)
			continue
		}
		if prog.Current != 0 {
			t.Errorf("mission %s should start at 0, got %d", mission.ID, prog.Current)
		}
		if prog.Completed {
			t.Errorf("mission %s should not be completed initially", mission.ID)
		}
		if prog.Target != mission.Target {
			t.Errorf("mission %s target mismatch: expected %d, got %d",
				mission.ID, mission.Target, prog.Target)
		}
	}
}

// TestUpdateProgress_KillTargets 確認擊破目標任務進度更新
func TestUpdateProgress_KillTargets(t *testing.T) {
	m := NewManager()
	playerID := "player-002"

	// 初始化進度
	m.GetOrInitProgress(playerID)

	// 更新進度（擊破 5 個目標）
	completed := m.UpdateProgress(playerID, MissionKillTargets, 5)
	if len(completed) != 0 {
		t.Errorf("should not complete yet (5/10), got %d completed", len(completed))
	}

	// 確認進度
	progress := m.GetOrInitProgress(playerID)
	for _, mission := range DailyMissions {
		if mission.Type != MissionKillTargets {
			continue
		}
		prog := progress[mission.ID]
		if prog.Current != 5 {
			t.Errorf("expected current=5, got %d", prog.Current)
		}
	}

	// 再擊破 5 個，應該完成
	completed = m.UpdateProgress(playerID, MissionKillTargets, 5)
	if len(completed) == 0 {
		t.Error("should complete kill_targets mission (10/10)")
	}
	if completed[0].Type != MissionKillTargets {
		t.Errorf("expected MissionKillTargets, got %s", completed[0].Type)
	}
}

// TestUpdateProgress_CapAtTarget 確認進度不超過目標值
func TestUpdateProgress_CapAtTarget(t *testing.T) {
	m := NewManager()
	playerID := "player-003"
	m.GetOrInitProgress(playerID)

	// 一次性超過目標
	m.UpdateProgress(playerID, MissionKillTargets, 100)

	progress := m.GetOrInitProgress(playerID)
	for _, mission := range DailyMissions {
		if mission.Type != MissionKillTargets {
			continue
		}
		prog := progress[mission.ID]
		if prog.Current > prog.Target {
			t.Errorf("current (%d) should not exceed target (%d)", prog.Current, prog.Target)
		}
		if !prog.Completed {
			t.Error("mission should be completed")
		}
	}
}

// TestClaimReward 確認任務獎勵領取
func TestClaimReward(t *testing.T) {
	m := NewManager()
	playerID := "player-004"
	m.GetOrInitProgress(playerID)

	// 完成 kill_boss 任務
	m.UpdateProgress(playerID, MissionKillBoss, 1)

	// 找到 kill_boss 任務 ID
	var bossMissionID string
	for _, mission := range DailyMissions {
		if mission.Type == MissionKillBoss {
			bossMissionID = mission.ID
			break
		}
	}
	if bossMissionID == "" {
		t.Fatal("kill_boss mission not found")
	}

	// 領取獎勵
	reward := m.ClaimReward(playerID, bossMissionID)
	if reward <= 0 {
		t.Errorf("expected reward > 0, got %d", reward)
	}

	// 不能重複領取
	reward2 := m.ClaimReward(playerID, bossMissionID)
	if reward2 != 0 {
		t.Errorf("should not be able to claim reward twice, got %d", reward2)
	}
}

// TestClaimReward_NotCompleted 確認未完成任務無法領取
func TestClaimReward_NotCompleted(t *testing.T) {
	m := NewManager()
	playerID := "player-005"
	m.GetOrInitProgress(playerID)

	// 不完成任務，直接嘗試領取
	for _, mission := range DailyMissions {
		reward := m.ClaimReward(playerID, mission.ID)
		if reward != 0 {
			t.Errorf("should not be able to claim uncompleted mission %s, got %d", mission.ID, reward)
		}
	}
}

// TestGetPlayerMissions 確認取得任務列表
func TestGetPlayerMissions(t *testing.T) {
	m := NewManager()
	playerID := "player-006"

	statuses := m.GetPlayerMissions(playerID)
	if len(statuses) != len(DailyMissions) {
		t.Errorf("expected %d missions, got %d", len(DailyMissions), len(statuses))
	}

	// 確認每個任務都有正確的定義
	for _, s := range statuses {
		if s.Mission.ID == "" {
			t.Error("mission ID should not be empty")
		}
		if s.Mission.Name == "" {
			t.Error("mission Name should not be empty")
		}
		if s.Mission.Target <= 0 {
			t.Errorf("mission %s target should be > 0, got %d", s.Mission.ID, s.Mission.Target)
		}
		if s.Mission.Reward <= 0 {
			t.Errorf("mission %s reward should be > 0, got %d", s.Mission.ID, s.Mission.Reward)
		}
	}
}

// TestDailyMissionsDefinition 確認每日任務定義完整
func TestDailyMissionsDefinition(t *testing.T) {
	if len(DailyMissions) == 0 {
		t.Fatal("DailyMissions should not be empty")
	}

	// 確認每個任務 ID 唯一
	ids := make(map[string]bool)
	for _, m := range DailyMissions {
		if ids[m.ID] {
			t.Errorf("duplicate mission ID: %s", m.ID)
		}
		ids[m.ID] = true

		if m.ID == "" {
			t.Error("mission ID should not be empty")
		}
		if m.Name == "" {
			t.Errorf("mission %s name should not be empty", m.ID)
		}
		if m.Target <= 0 {
			t.Errorf("mission %s target should be > 0", m.ID)
		}
		if m.Reward <= 0 {
			t.Errorf("mission %s reward should be > 0", m.ID)
		}
	}
}

// TestUpdateProgress_Combo 確認連擊任務進度更新（DAY-038）
// 注意：game.go 每次 combo 傳入 amount=1（不是 comboCount），避免連擊串讓任務太快完成
func TestUpdateProgress_Combo(t *testing.T) {
	m := NewManager()
	playerID := "player-007"
	m.GetOrInitProgress(playerID)

	// 模擬 4 次 combo 事件（每次 +1，未完成，目標 5）
	for i := 0; i < 4; i++ {
		completed := m.UpdateProgress(playerID, MissionCombo, 1)
		if len(completed) != 0 {
			t.Errorf("should not complete yet (%d/5), got %d completed", i+1, len(completed))
		}
	}

	// 確認進度
	progress := m.GetOrInitProgress(playerID)
	for _, mission := range DailyMissions {
		if mission.Type != MissionCombo {
			continue
		}
		prog := progress[mission.ID]
		if prog.Current != 4 {
			t.Errorf("expected current=4, got %d", prog.Current)
		}
	}

	// 第 5 次 combo 事件（累積 5，達到目標，應完成）
	completed := m.UpdateProgress(playerID, MissionCombo, 1)
	if len(completed) == 0 {
		t.Error("should complete combo mission (5/5)")
	}
	if completed[0].Type != MissionCombo {
		t.Errorf("expected MissionCombo, got %s", completed[0].Type)
	}
}

// TestAllMissionTypesPresent 確認所有任務類型都有對應的 DailyMission（DAY-038）
func TestAllMissionTypesPresent(t *testing.T) {
	typeSet := make(map[MissionType]bool)
	for _, m := range DailyMissions {
		typeSet[m.Type] = true
	}

	requiredTypes := []MissionType{
		MissionKillTargets,
		MissionKillBoss,
		MissionPlayBonus,
		MissionEarnCoins,
		MissionKillHighMult,
		MissionCombo,
	}

	for _, mType := range requiredTypes {
		if !typeSet[mType] {
			t.Errorf("mission type %s not found in DailyMissions", mType)
		}
	}
}
