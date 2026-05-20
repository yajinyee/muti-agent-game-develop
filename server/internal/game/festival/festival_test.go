package festival

import (
	"testing"
	"time"
)

func TestNew_DetectsFestival(t *testing.T) {
	m := New()
	// 不管當前是否有節日，Manager 應該正常建立
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
	if m.progress == nil {
		t.Fatal("expected non-nil progress map")
	}
}

func TestGetDef_NoFestival(t *testing.T) {
	m := &Manager{
		progress: make(map[string]*PlayerFestivalProgress),
		current:  FestivalNone,
	}
	def := m.GetDef()
	if def != nil {
		t.Errorf("expected nil def for FestivalNone, got %v", def)
	}
}

func TestGetDef_DragonBoat(t *testing.T) {
	m := &Manager{
		progress: make(map[string]*PlayerFestivalProgress),
		current:  FestivalDragonBoat,
		startAt:  time.Now().Add(-1 * time.Hour),
		endAt:    time.Now().Add(24 * time.Hour),
	}
	def := m.GetDef()
	if def == nil {
		t.Fatal("expected non-nil def for FestivalDragonBoat")
	}
	if def.Type != FestivalDragonBoat {
		t.Errorf("expected FestivalDragonBoat, got %v", def.Type)
	}
	if def.JackpotMult != 1.5 {
		t.Errorf("expected JackpotMult=1.5, got %v", def.JackpotMult)
	}
}

func TestIsActive_True(t *testing.T) {
	m := &Manager{
		progress: make(map[string]*PlayerFestivalProgress),
		current:  FestivalMidAutumn,
		startAt:  time.Now().Add(-1 * time.Hour),
		endAt:    time.Now().Add(24 * time.Hour),
	}
	if !m.IsActive() {
		t.Error("expected IsActive=true")
	}
}

func TestIsActive_False_None(t *testing.T) {
	m := &Manager{
		progress: make(map[string]*PlayerFestivalProgress),
		current:  FestivalNone,
	}
	if m.IsActive() {
		t.Error("expected IsActive=false for FestivalNone")
	}
}

func TestIsActive_False_Expired(t *testing.T) {
	m := &Manager{
		progress: make(map[string]*PlayerFestivalProgress),
		current:  FestivalHalloween,
		startAt:  time.Now().Add(-48 * time.Hour),
		endAt:    time.Now().Add(-1 * time.Hour),
	}
	if m.IsActive() {
		t.Error("expected IsActive=false for expired festival")
	}
}

func TestGetRewardMult_Active(t *testing.T) {
	m := &Manager{
		progress: make(map[string]*PlayerFestivalProgress),
		current:  FestivalNewYear,
		startAt:  time.Now().Add(-1 * time.Hour),
		endAt:    time.Now().Add(24 * time.Hour),
	}
	mult := m.GetRewardMult()
	if mult != 1.5 {
		t.Errorf("expected RewardMult=1.5 for NewYear, got %v", mult)
	}
}

func TestGetRewardMult_Inactive(t *testing.T) {
	m := &Manager{
		progress: make(map[string]*PlayerFestivalProgress),
		current:  FestivalNone,
	}
	mult := m.GetRewardMult()
	if mult != 1.0 {
		t.Errorf("expected RewardMult=1.0 for no festival, got %v", mult)
	}
}

func TestRecordKill_DragonBoat(t *testing.T) {
	m := &Manager{
		progress: make(map[string]*PlayerFestivalProgress),
		current:  FestivalDragonBoat,
		startAt:  time.Now().Add(-1 * time.Hour),
		endAt:    time.Now().Add(24 * time.Hour),
	}
	playerID := "p1"
	// 擊破粽子，任務 dt_kill_10 需要 10 個
	for i := 0; i < 9; i++ {
		updated, completed := m.RecordKill(playerID, "zongzi")
		if !updated {
			t.Errorf("expected taskUpdated=true at kill %d", i+1)
		}
		if completed != "" {
			t.Errorf("expected no completion at kill %d, got %v", i+1, completed)
		}
	}
	// 第 10 個應該完成任務
	updated, completed := m.RecordKill(playerID, "zongzi")
	if !updated {
		t.Error("expected taskUpdated=true at kill 10")
	}
	if completed == "" {
		t.Error("expected task completion at kill 10")
	}
}

func TestClaimTaskReward_Success(t *testing.T) {
	m := &Manager{
		progress: make(map[string]*PlayerFestivalProgress),
		current:  FestivalDragonBoat,
		startAt:  time.Now().Add(-1 * time.Hour),
		endAt:    time.Now().Add(24 * time.Hour),
	}
	playerID := "p2"
	// 手動設定進度完成
	p := &PlayerFestivalProgress{
		FestivalType: FestivalDragonBoat,
		TaskProgress: map[string]int{"dt_kill_10": 10},
		TaskDone:     make(map[string]bool),
	}
	m.progress[playerID] = p

	coins := m.ClaimTaskReward(playerID, "dt_kill_10")
	if coins != 5000 {
		t.Errorf("expected 5000 coins, got %d", coins)
	}
	// 重複領取應該回傳 0
	coins2 := m.ClaimTaskReward(playerID, "dt_kill_10")
	if coins2 != 0 {
		t.Errorf("expected 0 coins on duplicate claim, got %d", coins2)
	}
}

func TestClaimTitle_AllTasksDone(t *testing.T) {
	m := &Manager{
		progress: make(map[string]*PlayerFestivalProgress),
		current:  FestivalDragonBoat,
		startAt:  time.Now().Add(-1 * time.Hour),
		endAt:    time.Now().Add(24 * time.Hour),
	}
	playerID := "p3"
	def := FestivalDefs[FestivalDragonBoat]
	taskDone := make(map[string]bool)
	taskProgress := make(map[string]int)
	for _, task := range def.Tasks {
		taskDone[task.ID] = true
		taskProgress[task.ID] = task.Target
	}
	m.progress[playerID] = &PlayerFestivalProgress{
		FestivalType: FestivalDragonBoat,
		TaskProgress: taskProgress,
		TaskDone:     taskDone,
	}

	titleID, titleName, titleColor, ok := m.ClaimTitle(playerID)
	if !ok {
		t.Fatal("expected ClaimTitle to succeed")
	}
	if titleID != "festival_dragon_boat" {
		t.Errorf("expected titleID=festival_dragon_boat, got %v", titleID)
	}
	if titleName == "" || titleColor == "" {
		t.Error("expected non-empty titleName and titleColor")
	}
	// 重複領取應該失敗
	_, _, _, ok2 := m.ClaimTitle(playerID)
	if ok2 {
		t.Error("expected duplicate ClaimTitle to fail")
	}
}

func TestGetSnapshot_Active(t *testing.T) {
	m := &Manager{
		progress: make(map[string]*PlayerFestivalProgress),
		current:  FestivalHalloween,
		startAt:  time.Now().Add(-1 * time.Hour),
		endAt:    time.Now().Add(24 * time.Hour),
	}
	snap := m.GetSnapshot("p4")
	if snap.Type != string(FestivalHalloween) {
		t.Errorf("expected type=halloween, got %v", snap.Type)
	}
	if !snap.IsActive {
		t.Error("expected IsActive=true")
	}
	if len(snap.Tasks) == 0 {
		t.Error("expected non-empty tasks")
	}
	if len(snap.SpecialTargets) == 0 {
		t.Error("expected non-empty special targets")
	}
}

func TestRemovePlayer(t *testing.T) {
	m := &Manager{
		progress: make(map[string]*PlayerFestivalProgress),
		current:  FestivalMidAutumn,
		startAt:  time.Now().Add(-1 * time.Hour),
		endAt:    time.Now().Add(24 * time.Hour),
	}
	playerID := "p5"
	m.progress[playerID] = &PlayerFestivalProgress{
		FestivalType: FestivalMidAutumn,
		TaskProgress: make(map[string]int),
		TaskDone:     make(map[string]bool),
	}
	m.RemovePlayer(playerID)
	if _, ok := m.progress[playerID]; ok {
		t.Error("expected player to be removed")
	}
}
