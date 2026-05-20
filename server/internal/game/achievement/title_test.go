package achievement

import (
	"testing"
)

func TestNewTitleTracker(t *testing.T) {
	tracker := NewTitleTracker()
	if tracker == nil {
		t.Fatal("NewTitleTracker returned nil")
	}
	if !tracker.Unlocked[TitleNovice] {
		t.Error("Expected TitleNovice to be unlocked by default")
	}
	if tracker.ActiveTitle != TitleNovice {
		t.Errorf("Expected ActiveTitle to be TitleNovice, got %s", tracker.ActiveTitle)
	}
}

func TestOnAchievementUnlocked_KillMilestone(t *testing.T) {
	tracker := NewTitleTracker()

	// 解鎖 kill_5 成就 → 應解鎖 TitleHunter
	titleDef := tracker.OnAchievementUnlocked(AchKill5, 2)
	if titleDef == nil {
		t.Fatal("Expected TitleHunter to be unlocked")
	}
	if titleDef.ID != TitleHunter {
		t.Errorf("Expected TitleHunter, got %s", titleDef.ID)
	}
	if !tracker.Unlocked[TitleHunter] {
		t.Error("TitleHunter should be in Unlocked map")
	}
	if tracker.ActiveTitle != TitleHunter {
		t.Errorf("Expected ActiveTitle to be TitleHunter, got %s", tracker.ActiveTitle)
	}
}

func TestOnAchievementUnlocked_HigherPriorityWins(t *testing.T) {
	tracker := NewTitleTracker()

	// 先解鎖 Hunter（priority 10）
	tracker.OnAchievementUnlocked(AchKill5, 2)
	if tracker.ActiveTitle != TitleHunter {
		t.Errorf("Expected TitleHunter, got %s", tracker.ActiveTitle)
	}

	// 再解鎖 Legend（priority 50）→ 應覆蓋 Hunter
	tracker.OnAchievementUnlocked(AchKill100, 5)
	if tracker.ActiveTitle != TitleLegend {
		t.Errorf("Expected TitleLegend to override TitleHunter, got %s", tracker.ActiveTitle)
	}
}

func TestOnAchievementUnlocked_AllAround(t *testing.T) {
	tracker := NewTitleTracker()

	// 解鎖 8 個成就 → 應解鎖 TitleAllAround
	titleDef := tracker.OnAchievementUnlocked(AchFirstKill, 8)
	if titleDef == nil {
		t.Fatal("Expected TitleAllAround to be unlocked at 8 achievements")
	}
	if titleDef.ID != TitleAllAround {
		t.Errorf("Expected TitleAllAround, got %s", titleDef.ID)
	}
}

func TestOnAchievementUnlocked_NoDuplicate(t *testing.T) {
	tracker := NewTitleTracker()

	// 第一次解鎖
	first := tracker.OnAchievementUnlocked(AchKill5, 2)
	if first == nil {
		t.Fatal("Expected first unlock to succeed")
	}

	// 第二次解鎖同一個成就 → 應回傳 nil
	second := tracker.OnAchievementUnlocked(AchKill5, 3)
	if second != nil {
		t.Error("Expected second unlock of same achievement to return nil")
	}
}

func TestSetActiveTitle(t *testing.T) {
	tracker := NewTitleTracker()

	// 嘗試設定未解鎖的稱號 → 應失敗
	ok := tracker.SetActiveTitle(TitleLegend)
	if ok {
		t.Error("Expected SetActiveTitle to fail for unowned title")
	}

	// 解鎖 Hunter
	tracker.OnAchievementUnlocked(AchKill5, 2)
	// 解鎖 Legend
	tracker.OnAchievementUnlocked(AchKill100, 5)

	// 手動設定回 Hunter
	ok = tracker.SetActiveTitle(TitleHunter)
	if !ok {
		t.Error("Expected SetActiveTitle to succeed for owned title")
	}
	if tracker.ActiveTitle != TitleHunter {
		t.Errorf("Expected ActiveTitle to be TitleHunter, got %s", tracker.ActiveTitle)
	}
}

func TestGetActiveTitle(t *testing.T) {
	tracker := NewTitleTracker()
	def := tracker.GetActiveTitle()
	if def == nil {
		t.Fatal("GetActiveTitle returned nil")
	}
	if def.ID != TitleNovice {
		t.Errorf("Expected TitleNovice, got %s", def.ID)
	}
}

func TestGetUnlockedTitles(t *testing.T) {
	tracker := NewTitleTracker()
	titles := tracker.GetUnlockedTitles()
	if len(titles) != 1 {
		t.Errorf("Expected 1 unlocked title (Novice), got %d", len(titles))
	}

	tracker.OnAchievementUnlocked(AchKill5, 2)
	titles = tracker.GetUnlockedTitles()
	if len(titles) != 2 {
		t.Errorf("Expected 2 unlocked titles, got %d", len(titles))
	}
}

func TestAllTitleDefinitionsValid(t *testing.T) {
	for id, def := range TitleDefinitions {
		if def.ID != id {
			t.Errorf("Title %s has mismatched ID %s", id, def.ID)
		}
		if def.Name == "" {
			t.Errorf("Title %s has empty Name", id)
		}
		if def.Icon == "" {
			t.Errorf("Title %s has empty Icon", id)
		}
		if def.Color == "" {
			t.Errorf("Title %s has empty Color", id)
		}
	}
}

func TestBossSlayerTitle(t *testing.T) {
	tracker := NewTitleTracker()
	titleDef := tracker.OnAchievementUnlocked(AchKillBoss, 2)
	if titleDef == nil {
		t.Fatal("Expected TitleBossSlayer to be unlocked")
	}
	if titleDef.ID != TitleBossSlayer {
		t.Errorf("Expected TitleBossSlayer, got %s", titleDef.ID)
	}
}

func TestMillionaireTitle(t *testing.T) {
	tracker := NewTitleTracker()
	titleDef := tracker.OnAchievementUnlocked(AchCoins100k, 3)
	if titleDef == nil {
		t.Fatal("Expected TitleMillionaire to be unlocked")
	}
	if titleDef.ID != TitleMillionaire {
		t.Errorf("Expected TitleMillionaire, got %s", titleDef.ID)
	}
}
