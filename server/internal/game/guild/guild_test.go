package guild

import (
	"testing"
)

func TestCreateGuild(t *testing.T) {
	m := New()

	guildID, err := m.CreateGuild("p1", "Player1", "吉伊卡哇公會", "一起討伐！")
	if err != nil {
		t.Fatalf("CreateGuild failed: %v", err)
	}
	if guildID == "" {
		t.Fatal("expected non-empty guildID")
	}

	guild := m.GetGuild(guildID)
	if guild == nil {
		t.Fatal("guild not found")
	}
	if guild.Name != "吉伊卡哇公會" {
		t.Errorf("expected name '吉伊卡哇公會', got '%s'", guild.Name)
	}
	if guild.Level != 1 {
		t.Errorf("expected level 1, got %d", guild.Level)
	}
	if len(guild.Members) != 1 {
		t.Errorf("expected 1 member, got %d", len(guild.Members))
	}
	if guild.Members["p1"].Role != RoleLeader {
		t.Errorf("expected leader role, got %s", guild.Members["p1"].Role)
	}
	if len(guild.Tasks) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(guild.Tasks))
	}
}

func TestCreateGuild_DuplicateName(t *testing.T) {
	m := New()
	m.CreateGuild("p1", "Player1", "測試公會", "")
	_, err := m.CreateGuild("p2", "Player2", "測試公會", "")
	if err == nil {
		t.Fatal("expected error for duplicate guild name")
	}
}

func TestCreateGuild_AlreadyInGuild(t *testing.T) {
	m := New()
	m.CreateGuild("p1", "Player1", "公會A", "")
	_, err := m.CreateGuild("p1", "Player1", "公會B", "")
	if err == nil {
		t.Fatal("expected error when player already in guild")
	}
}

func TestJoinGuild(t *testing.T) {
	m := New()
	guildID, _ := m.CreateGuild("p1", "Player1", "測試公會", "")

	err := m.JoinGuild("p2", "Player2", guildID)
	if err != nil {
		t.Fatalf("JoinGuild failed: %v", err)
	}

	guild := m.GetGuild(guildID)
	if len(guild.Members) != 2 {
		t.Errorf("expected 2 members, got %d", len(guild.Members))
	}
	if guild.Members["p2"].Role != RoleMember {
		t.Errorf("expected member role, got %s", guild.Members["p2"].Role)
	}
}

func TestJoinGuild_AlreadyInGuild(t *testing.T) {
	m := New()
	guildID, _ := m.CreateGuild("p1", "Player1", "公會A", "")
	m.JoinGuild("p2", "Player2", guildID)

	guildID2, _ := m.CreateGuild("p3", "Player3", "公會B", "")
	err := m.JoinGuild("p2", "Player2", guildID2)
	if err == nil {
		t.Fatal("expected error when player already in guild")
	}
}

func TestLeaveGuild_Member(t *testing.T) {
	m := New()
	guildID, _ := m.CreateGuild("p1", "Player1", "測試公會", "")
	m.JoinGuild("p2", "Player2", guildID)

	_, err := m.LeaveGuild("p2")
	if err != nil {
		t.Fatalf("LeaveGuild failed: %v", err)
	}

	guild := m.GetGuild(guildID)
	if len(guild.Members) != 1 {
		t.Errorf("expected 1 member after leave, got %d", len(guild.Members))
	}
	if m.GetPlayerGuildID("p2") != "" {
		t.Error("expected player to have no guild after leaving")
	}
}

func TestLeaveGuild_LeaderTransfer(t *testing.T) {
	m := New()
	guildID, _ := m.CreateGuild("p1", "Player1", "測試公會", "")
	m.JoinGuild("p2", "Player2", guildID)
	m.PromoteMember("p1", "p2") // p2 升為副會長

	_, err := m.LeaveGuild("p1")
	if err != nil {
		t.Fatalf("LeaveGuild failed: %v", err)
	}

	guild := m.GetGuild(guildID)
	if guild == nil {
		t.Fatal("guild should still exist")
	}
	if guild.Members["p2"].Role != RoleLeader {
		t.Errorf("expected p2 to become leader, got %s", guild.Members["p2"].Role)
	}
}

func TestLeaveGuild_LastMember_Dissolve(t *testing.T) {
	m := New()
	guildID, _ := m.CreateGuild("p1", "Player1", "測試公會", "")

	_, err := m.LeaveGuild("p1")
	if err != nil {
		t.Fatalf("LeaveGuild failed: %v", err)
	}

	guild := m.GetGuild(guildID)
	if guild != nil {
		t.Error("guild should be dissolved when last member leaves")
	}
}

func TestKickMember(t *testing.T) {
	m := New()
	guildID, _ := m.CreateGuild("p1", "Player1", "測試公會", "")
	m.JoinGuild("p2", "Player2", guildID)

	err := m.KickMember("p1", "p2")
	if err != nil {
		t.Fatalf("KickMember failed: %v", err)
	}

	guild := m.GetGuild(guildID)
	if _, ok := guild.Members["p2"]; ok {
		t.Error("p2 should have been kicked")
	}
}

func TestKickMember_NoPermission(t *testing.T) {
	m := New()
	guildID, _ := m.CreateGuild("p1", "Player1", "測試公會", "")
	m.JoinGuild("p2", "Player2", guildID)
	m.JoinGuild("p3", "Player3", guildID)

	err := m.KickMember("p2", "p3") // 普通成員不能踢人
	if err == nil {
		t.Fatal("expected error when member tries to kick")
	}
}

func TestPromoteMember(t *testing.T) {
	m := New()
	guildID, _ := m.CreateGuild("p1", "Player1", "測試公會", "")
	m.JoinGuild("p2", "Player2", guildID)

	err := m.PromoteMember("p1", "p2")
	if err != nil {
		t.Fatalf("PromoteMember failed: %v", err)
	}

	guild := m.GetGuild(guildID)
	if guild.Members["p2"].Role != RoleOfficer {
		t.Errorf("expected officer role, got %s", guild.Members["p2"].Role)
	}
}

func TestUpdateTaskProgress(t *testing.T) {
	m := New()
	guildID, _ := m.CreateGuild("p1", "Player1", "測試公會", "")

	// 擊破 50 個目標（未完成）
	completed := m.UpdateTaskProgress("p1", TaskKillTargets, 50)
	if len(completed) != 0 {
		t.Errorf("expected 0 completed tasks, got %d", len(completed))
	}

	// 再擊破 50 個（完成）
	completed = m.UpdateTaskProgress("p1", TaskKillTargets, 50)
	if len(completed) != 1 {
		t.Errorf("expected 1 completed task, got %d", len(completed))
	}
	if completed[0].Type != TaskKillTargets {
		t.Errorf("expected kill_targets task, got %s", completed[0].Type)
	}

	guild := m.GetGuild(guildID)
	if !guild.Tasks[0].Completed {
		t.Error("task should be completed")
	}
}

func TestUpdateTaskProgress_NotInGuild(t *testing.T) {
	m := New()
	completed := m.UpdateTaskProgress("p_nobody", TaskKillTargets, 10)
	if completed != nil {
		t.Error("expected nil for player not in guild")
	}
}

func TestSetOnlineStatus(t *testing.T) {
	m := New()
	guildID, _ := m.CreateGuild("p1", "Player1", "測試公會", "")

	m.SetOnlineStatus("p1", false)
	guild := m.GetGuild(guildID)
	if guild.Members["p1"].IsOnline {
		t.Error("expected p1 to be offline")
	}

	m.SetOnlineStatus("p1", true)
	if !guild.Members["p1"].IsOnline {
		t.Error("expected p1 to be online")
	}
}

func TestGetAllGuilds(t *testing.T) {
	m := New()
	m.CreateGuild("p1", "Player1", "公會A", "")
	m.CreateGuild("p2", "Player2", "公會B", "")

	guilds := m.GetAllGuilds()
	if len(guilds) != 2 {
		t.Errorf("expected 2 guilds, got %d", len(guilds))
	}
}

func TestGuildLevelUp(t *testing.T) {
	m := New()
	guildID, _ := m.CreateGuild("p1", "Player1", "測試公會", "")

	// 擊破大量目標觸發公會升級
	// 每完成一個任務（100個目標）獲得 10 exp，需要 100 exp 升到 2 級
	for i := 0; i < 10; i++ {
		m.UpdateTaskProgress("p1", TaskKillTargets, 100)
		// 重置任務讓它可以再次完成
		guild := m.GetGuild(guildID)
		for _, task := range guild.Tasks {
			if task.Type == TaskKillTargets {
				task.Current = 0
				task.Completed = false
			}
		}
	}

	guild := m.GetGuild(guildID)
	if guild.Level < 2 {
		t.Errorf("expected guild level >= 2, got %d (exp: %d)", guild.Level, guild.Exp)
	}
}
