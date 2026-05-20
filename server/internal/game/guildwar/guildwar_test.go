package guildwar

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	status, weekID, startAt, endAt := m.GetStatus()
	if status != WarStatusActive {
		t.Errorf("expected status Active, got %s", status)
	}
	if weekID == "" {
		t.Error("weekID should not be empty")
	}
	if startAt.IsZero() || endAt.IsZero() {
		t.Error("startAt/endAt should not be zero")
	}
	if !endAt.After(startAt) {
		t.Error("endAt should be after startAt")
	}
}

func TestEnsureGuildRegistered(t *testing.T) {
	m := New()
	m.EnsureGuildRegistered("guild-001", "測試公會", "⚔️", 5)

	score := m.GetGuildScore("guild-001")
	if score == nil {
		t.Fatal("guild should be registered")
	}
	if score.GuildName != "測試公會" {
		t.Errorf("expected name 測試公會, got %s", score.GuildName)
	}
	if score.MemberCount != 5 {
		t.Errorf("expected 5 members, got %d", score.MemberCount)
	}
}

func TestAddKillScore(t *testing.T) {
	m := New()
	m.EnsureGuildRegistered("guild-001", "A", "⚔️", 3)

	// 普通目標（2x）= 1 分
	m.AddKillScore("guild-001", 2)
	score := m.GetGuildScore("guild-001")
	if score.Score != 1 {
		t.Errorf("expected 1, got %d", score.Score)
	}

	// 高倍率（10x）= 2 分
	m.AddKillScore("guild-001", 10)
	score = m.GetGuildScore("guild-001")
	if score.Score != 3 {
		t.Errorf("expected 3, got %d", score.Score)
	}

	// 超高倍率（20x）= 3 分
	m.AddKillScore("guild-001", 20)
	score = m.GetGuildScore("guild-001")
	if score.Score != 6 {
		t.Errorf("expected 6, got %d", score.Score)
	}

	// 極高倍率（50x）= 5 分
	m.AddKillScore("guild-001", 50)
	score = m.GetGuildScore("guild-001")
	if score.Score != 11 {
		t.Errorf("expected 11, got %d", score.Score)
	}
}

func TestAddBossScore(t *testing.T) {
	m := New()
	m.EnsureGuildRegistered("guild-001", "A", "⚔️", 3)

	m.AddBossScore("guild-001")
	score := m.GetGuildScore("guild-001")
	if score.BossScore != 50 {
		t.Errorf("expected 50, got %d", score.BossScore)
	}
	if score.Score != 50 {
		t.Errorf("expected total 50, got %d", score.Score)
	}
}

func TestAddBonusScore(t *testing.T) {
	m := New()
	m.EnsureGuildRegistered("guild-001", "A", "⚔️", 3)

	m.AddBonusScore("guild-001")
	score := m.GetGuildScore("guild-001")
	if score.BonusScore != 20 {
		t.Errorf("expected 20, got %d", score.BonusScore)
	}
}

func TestGetRankings(t *testing.T) {
	m := New()
	m.EnsureGuildRegistered("guild-001", "A", "⚔️", 3)
	m.EnsureGuildRegistered("guild-002", "B", "🛡️", 5)
	m.EnsureGuildRegistered("guild-003", "C", "🔥", 2)

	m.AddBossScore("guild-002") // 50 分
	m.AddBossScore("guild-002") // 100 分
	m.AddKillScore("guild-001", 2) // 1 分
	m.AddBonusScore("guild-003") // 20 分

	rankings := m.GetRankings()
	if len(rankings) != 3 {
		t.Fatalf("expected 3 rankings, got %d", len(rankings))
	}
	if rankings[0].GuildID != "guild-002" {
		t.Errorf("expected guild-002 first, got %s", rankings[0].GuildID)
	}
	if rankings[1].GuildID != "guild-003" {
		t.Errorf("expected guild-003 second, got %s", rankings[1].GuildID)
	}
	if rankings[2].GuildID != "guild-001" {
		t.Errorf("expected guild-001 third, got %s", rankings[2].GuildID)
	}
}

func TestGetGuildRank(t *testing.T) {
	m := New()
	m.EnsureGuildRegistered("guild-001", "A", "⚔️", 3)
	m.EnsureGuildRegistered("guild-002", "B", "🛡️", 5)

	m.AddBossScore("guild-002") // 50 分

	rank := m.GetGuildRank("guild-002")
	if rank != 1 {
		t.Errorf("expected rank 1, got %d", rank)
	}

	rank = m.GetGuildRank("guild-001")
	if rank != 2 {
		t.Errorf("expected rank 2, got %d", rank)
	}

	rank = m.GetGuildRank("guild-999")
	if rank != 0 {
		t.Errorf("expected rank 0 for unknown guild, got %d", rank)
	}
}

func TestGetWeekID(t *testing.T) {
	// 2026-05-20 是第 21 週
	t1 := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	weekID := getWeekID(t1)
	if weekID != "2026-W21" {
		t.Errorf("expected 2026-W21, got %s", weekID)
	}
}

func TestGetWeekRange(t *testing.T) {
	// 2026-05-20（週三）
	t1 := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	start, end := getWeekRange(t1)

	// 週一應該是 2026-05-18
	loc := time.FixedZone("UTC+8", 8*60*60)
	startLocal := start.In(loc)
	if startLocal.Weekday() != time.Monday {
		t.Errorf("start should be Monday, got %s", startLocal.Weekday())
	}

	// 週日應該是 2026-05-24
	endLocal := end.In(loc)
	if endLocal.Weekday() != time.Sunday {
		t.Errorf("end should be Sunday, got %s", endLocal.Weekday())
	}

	// end 應該在 start 之後
	if !end.After(start) {
		t.Error("end should be after start")
	}
}

func TestParticipatingGuildCount(t *testing.T) {
	m := New()
	if m.GetParticipatingGuildCount() != 0 {
		t.Error("should start with 0 guilds")
	}

	m.EnsureGuildRegistered("guild-001", "A", "⚔️", 3)
	m.EnsureGuildRegistered("guild-002", "B", "🛡️", 5)

	if m.GetParticipatingGuildCount() != 2 {
		t.Errorf("expected 2, got %d", m.GetParticipatingGuildCount())
	}
}

func TestAddScoreForUnregisteredGuild(t *testing.T) {
	m := New()
	// 對未登記的公會加分，不應 panic
	m.AddKillScore("guild-999", 5)
	m.AddBossScore("guild-999")
	m.AddBonusScore("guild-999")

	score := m.GetGuildScore("guild-999")
	if score != nil {
		t.Error("unregistered guild should not have score")
	}
}
