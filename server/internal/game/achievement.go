// Package game — 成就系統
// DAY-349：玩家成就解鎖、通知、進度追蹤
// achievement-agent 負責維護
package game

import (
	"sync"
	"time"
)

// AchievementType 成就類型
type AchievementType string

const (
	AchievTypeFirstKill      AchievementType = "first_kill"       // 第一次擊破目標
	AchievTypeFirstBoss      AchievementType = "first_boss"       // 第一次擊破 BOSS
	AchievTypeFirstBonus     AchievementType = "first_bonus"      // 第一次完成 Bonus
	AchievTypeCombo5         AchievementType = "combo_5"          // 達成 5 連擊
	AchievTypeCombo10        AchievementType = "combo_10"         // 達成 10 連擊
	AchievTypeCombo20        AchievementType = "combo_20"         // 達成 20 連擊
	AchievTypeCombo30        AchievementType = "combo_30"         // 達成 30 連擊
	AchievTypeKill100        AchievementType = "kill_100"         // 累積擊破 100 個
	AchievTypeKill500        AchievementType = "kill_500"         // 累積擊破 500 個
	AchievTypeKill1000       AchievementType = "kill_1000"        // 累積擊破 1000 個
	AchievTypeKill5000       AchievementType = "kill_5000"        // 累積擊破 5000 個
	AchievTypeMult50         AchievementType = "mult_50"          // 單次獲得 50x 以上
	AchievTypeMult100        AchievementType = "mult_100"         // 單次獲得 100x 以上
	AchievTypeMult500        AchievementType = "mult_500"         // 單次獲得 500x 以上
	AchievTypeMult1000       AchievementType = "mult_1000"        // 單次獲得 1000x 以上
	AchievTypeLucky10        AchievementType = "lucky_10"         // 觸發 10 次幸運魚
	AchievTypeLucky50        AchievementType = "lucky_50"         // 觸發 50 次幸運魚
	AchievTypeLucky100       AchievementType = "lucky_100"        // 觸發 100 次幸運魚
	AchievTypeSeasonLv5      AchievementType = "season_lv5"       // 賽季通行證達到 Lv.5
	AchievTypeSeasonLv10     AchievementType = "season_lv10"      // 賽季通行證達到 Lv.10（最高）
	AchievTypeQuestComplete5 AchievementType = "quest_complete_5" // 完成 5 個每日任務
	AchievTypeQuestComplete20 AchievementType = "quest_complete_20" // 完成 20 個每日任務
	AchievTypeWeeklyComplete AchievementType = "weekly_complete"  // 完成一次每週挑戰
	AchievTypeRichPlayer     AchievementType = "rich_player"      // 金幣達到 10000
	AchievTypeMegaRich       AchievementType = "mega_rich"        // 金幣達到 100000
)

// AchievementDef 成就定義
type AchievementDef struct {
	ID          AchievementType
	Name        string
	Description string
	Icon        string // emoji 圖示
	Rarity      string // common | rare | epic | legendary
	Reward      int    // 解鎖獎勵（金幣）
}

// 成就定義表
var achievementDefs = map[AchievementType]*AchievementDef{
	AchievTypeFirstKill:      {AchievTypeFirstKill, "初次討伐", "第一次擊破目標物", "⚔️", "common", 50},
	AchievTypeFirstBoss:      {AchievTypeFirstBoss, "BOSS 終結者", "第一次擊破那個孩子", "👑", "rare", 500},
	AchievTypeFirstBonus:     {AchievTypeFirstBonus, "拔草達人", "第一次完成瘋狂拔草", "🌿", "rare", 300},
	AchievTypeCombo5:         {AchievTypeCombo5, "連擊新手", "達成 5 連擊", "🔥", "common", 100},
	AchievTypeCombo10:        {AchievTypeCombo10, "連擊高手", "達成 10 連擊", "💥", "rare", 300},
	AchievTypeCombo20:        {AchievTypeCombo20, "連擊大師", "達成 20 連擊", "⚡", "epic", 800},
	AchievTypeCombo30:        {AchievTypeCombo30, "連擊傳說", "達成 30 連擊", "🌟", "legendary", 2000},
	AchievTypeKill100:        {AchievTypeKill100, "百戰老兵", "累積擊破 100 個目標", "🎯", "common", 200},
	AchievTypeKill500:        {AchievTypeKill500, "五百討伐", "累積擊破 500 個目標", "⚔️", "rare", 500},
	AchievTypeKill1000:       {AchievTypeKill1000, "千人斬", "累積擊破 1000 個目標", "🗡️", "epic", 1500},
	AchievTypeKill5000:       {AchievTypeKill5000, "萬物終結", "累積擊破 5000 個目標", "💀", "legendary", 5000},
	AchievTypeMult50:         {AchievTypeMult50, "大獎初體驗", "單次獲得 50x 以上倍率", "💰", "common", 150},
	AchievTypeMult100:        {AchievTypeMult100, "百倍爆發", "單次獲得 100x 以上倍率", "💎", "rare", 400},
	AchievTypeMult500:        {AchievTypeMult500, "五百倍奇蹟", "單次獲得 500x 以上倍率", "🌈", "epic", 1000},
	AchievTypeMult1000:       {AchievTypeMult1000, "千倍傳說", "單次獲得 1000x 以上倍率", "🏆", "legendary", 3000},
	AchievTypeLucky10:        {AchievTypeLucky10, "幸運初探", "觸發 10 次幸運魚", "🍀", "common", 200},
	AchievTypeLucky50:        {AchievTypeLucky50, "幸運常客", "觸發 50 次幸運魚", "✨", "rare", 600},
	AchievTypeLucky100:       {AchievTypeLucky100, "幸運之神", "觸發 100 次幸運魚", "🌠", "epic", 1500},
	AchievTypeSeasonLv5:      {AchievTypeSeasonLv5, "賽季中堅", "賽季通行證達到 Lv.5", "🎖️", "rare", 500},
	AchievTypeSeasonLv10:     {AchievTypeSeasonLv10, "賽季傳說", "賽季通行證達到最高等級", "👑", "legendary", 2000},
	AchievTypeQuestComplete5: {AchievTypeQuestComplete5, "任務新手", "完成 5 個每日任務", "📋", "common", 200},
	AchievTypeQuestComplete20: {AchievTypeQuestComplete20, "任務達人", "完成 20 個每日任務", "📜", "rare", 600},
	AchievTypeWeeklyComplete: {AchievTypeWeeklyComplete, "週挑戰勇者", "完成一次每週挑戰", "🏅", "rare", 800},
	AchievTypeRichPlayer:     {AchievTypeRichPlayer, "小富翁", "金幣達到 10000", "💵", "rare", 500},
	AchievTypeMegaRich:       {AchievTypeMegaRich, "大富翁", "金幣達到 100000", "💰", "legendary", 3000},
}

// PlayerAchievement 玩家成就記錄
type PlayerAchievement struct {
	Type        AchievementType
	UnlockedAt  time.Time
	Notified    bool // 是否已通知玩家
}

// AchievementProgress 成就進度（用於有進度的成就）
type AchievementProgress struct {
	KillCount    int // 累積擊破數
	LuckyCount   int // 幸運魚觸發數
	QuestCount   int // 每日任務完成數
}

// AchievementSystem 成就系統
type AchievementSystem struct {
	mu          sync.RWMutex
	// playerID -> 已解鎖成就列表
	unlocked    map[string][]*PlayerAchievement
	// playerID -> 進度
	progress    map[string]*AchievementProgress
}

func newAchievementSystem() *AchievementSystem {
	return &AchievementSystem{
		unlocked: make(map[string][]*PlayerAchievement),
		progress: make(map[string]*AchievementProgress),
	}
}

// getOrCreateProgress 取得或建立玩家進度
func (a *AchievementSystem) getOrCreateProgress(playerID string) *AchievementProgress {
	if p, ok := a.progress[playerID]; ok {
		return p
	}
	p := &AchievementProgress{}
	a.progress[playerID] = p
	return p
}

// isUnlocked 確認成就是否已解鎖
func (a *AchievementSystem) isUnlocked(playerID string, t AchievementType) bool {
	for _, ach := range a.unlocked[playerID] {
		if ach.Type == t {
			return true
		}
	}
	return false
}

// unlock 解鎖成就，回傳是否為新解鎖
func (a *AchievementSystem) unlock(playerID string, t AchievementType) bool {
	if a.isUnlocked(playerID, t) {
		return false
	}
	a.unlocked[playerID] = append(a.unlocked[playerID], &PlayerAchievement{
		Type:       t,
		UnlockedAt: time.Now(),
		Notified:   false,
	})
	return true
}

// GetPendingNotifications 取得待通知的成就（並標記為已通知）
func (a *AchievementSystem) GetPendingNotifications(playerID string) []*AchievementDef {
	a.mu.Lock()
	defer a.mu.Unlock()
	var result []*AchievementDef
	for _, ach := range a.unlocked[playerID] {
		if !ach.Notified {
			ach.Notified = true
			if def, ok := achievementDefs[ach.Type]; ok {
				result = append(result, def)
			}
		}
	}
	return result
}

// GetAllUnlocked 取得玩家所有已解鎖成就
func (a *AchievementSystem) GetAllUnlocked(playerID string) []*AchievementDef {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var result []*AchievementDef
	for _, ach := range a.unlocked[playerID] {
		if def, ok := achievementDefs[ach.Type]; ok {
			result = append(result, def)
		}
	}
	return result
}

// GetProgress 取得玩家成就進度
func (a *AchievementSystem) GetProgress(playerID string) *AchievementProgress {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if p, ok := a.progress[playerID]; ok {
		return p
	}
	return &AchievementProgress{}
}

// OnKill 擊破目標時觸發
// 回傳新解鎖的成就列表
func (a *AchievementSystem) OnKill(playerID string, multiplier float64, isFirst bool) []*AchievementDef {
	a.mu.Lock()
	defer a.mu.Unlock()

	var newAchievs []*AchievementDef
	prog := a.getOrCreateProgress(playerID)
	prog.KillCount++

	// 首次擊破
	if isFirst && a.unlock(playerID, AchievTypeFirstKill) {
		newAchievs = append(newAchievs, achievementDefs[AchievTypeFirstKill])
	}

	// 累積擊破里程碑
	milestones := []struct {
		count int
		t     AchievementType
	}{
		{100, AchievTypeKill100},
		{500, AchievTypeKill500},
		{1000, AchievTypeKill1000},
		{5000, AchievTypeKill5000},
	}
	for _, m := range milestones {
		if prog.KillCount >= m.count && a.unlock(playerID, m.t) {
			newAchievs = append(newAchievs, achievementDefs[m.t])
		}
	}

	// 倍率里程碑
	multMilestones := []struct {
		mult float64
		t    AchievementType
	}{
		{50, AchievTypeMult50},
		{100, AchievTypeMult100},
		{500, AchievTypeMult500},
		{1000, AchievTypeMult1000},
	}
	for _, m := range multMilestones {
		if multiplier >= m.mult && a.unlock(playerID, m.t) {
			newAchievs = append(newAchievs, achievementDefs[m.t])
		}
	}

	return newAchievs
}

// OnBossKill BOSS 擊破時觸發
func (a *AchievementSystem) OnBossKill(playerID string) []*AchievementDef {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.unlock(playerID, AchievTypeFirstBoss) {
		return []*AchievementDef{achievementDefs[AchievTypeFirstBoss]}
	}
	return nil
}

// OnBonusComplete Bonus 完成時觸發
func (a *AchievementSystem) OnBonusComplete(playerID string) []*AchievementDef {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.unlock(playerID, AchievTypeFirstBonus) {
		return []*AchievementDef{achievementDefs[AchievTypeFirstBonus]}
	}
	return nil
}

// OnCombo 連擊達成時觸發
func (a *AchievementSystem) OnCombo(playerID string, comboCount int) []*AchievementDef {
	a.mu.Lock()
	defer a.mu.Unlock()
	var newAchievs []*AchievementDef
	milestones := []struct {
		count int
		t     AchievementType
	}{
		{5, AchievTypeCombo5},
		{10, AchievTypeCombo10},
		{20, AchievTypeCombo20},
		{30, AchievTypeCombo30},
	}
	for _, m := range milestones {
		if comboCount >= m.count && a.unlock(playerID, m.t) {
			newAchievs = append(newAchievs, achievementDefs[m.t])
		}
	}
	return newAchievs
}

// OnLuckyFish 幸運魚觸發時
func (a *AchievementSystem) OnLuckyFish(playerID string) []*AchievementDef {
	a.mu.Lock()
	defer a.mu.Unlock()
	prog := a.getOrCreateProgress(playerID)
	prog.LuckyCount++
	var newAchievs []*AchievementDef
	milestones := []struct {
		count int
		t     AchievementType
	}{
		{10, AchievTypeLucky10},
		{50, AchievTypeLucky50},
		{100, AchievTypeLucky100},
	}
	for _, m := range milestones {
		if prog.LuckyCount >= m.count && a.unlock(playerID, m.t) {
			newAchievs = append(newAchievs, achievementDefs[m.t])
		}
	}
	return newAchievs
}

// OnSeasonLevel 賽季等級提升時
func (a *AchievementSystem) OnSeasonLevel(playerID string, level int) []*AchievementDef {
	a.mu.Lock()
	defer a.mu.Unlock()
	var newAchievs []*AchievementDef
	if level >= 5 && a.unlock(playerID, AchievTypeSeasonLv5) {
		newAchievs = append(newAchievs, achievementDefs[AchievTypeSeasonLv5])
	}
	if level >= 10 && a.unlock(playerID, AchievTypeSeasonLv10) {
		newAchievs = append(newAchievs, achievementDefs[AchievTypeSeasonLv10])
	}
	return newAchievs
}

// OnQuestComplete 每日任務完成時
func (a *AchievementSystem) OnQuestComplete(playerID string) []*AchievementDef {
	a.mu.Lock()
	defer a.mu.Unlock()
	prog := a.getOrCreateProgress(playerID)
	prog.QuestCount++
	var newAchievs []*AchievementDef
	milestones := []struct {
		count int
		t     AchievementType
	}{
		{5, AchievTypeQuestComplete5},
		{20, AchievTypeQuestComplete20},
	}
	for _, m := range milestones {
		if prog.QuestCount >= m.count && a.unlock(playerID, m.t) {
			newAchievs = append(newAchievs, achievementDefs[m.t])
		}
	}
	return newAchievs
}

// OnWeeklyComplete 每週挑戰完成時
func (a *AchievementSystem) OnWeeklyComplete(playerID string) []*AchievementDef {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.unlock(playerID, AchievTypeWeeklyComplete) {
		return []*AchievementDef{achievementDefs[AchievTypeWeeklyComplete]}
	}
	return nil
}

// OnCoinsUpdate 金幣更新時
func (a *AchievementSystem) OnCoinsUpdate(playerID string, coins int) []*AchievementDef {
	a.mu.Lock()
	defer a.mu.Unlock()
	var newAchievs []*AchievementDef
	if coins >= 10000 && a.unlock(playerID, AchievTypeRichPlayer) {
		newAchievs = append(newAchievs, achievementDefs[AchievTypeRichPlayer])
	}
	if coins >= 100000 && a.unlock(playerID, AchievTypeMegaRich) {
		newAchievs = append(newAchievs, achievementDefs[AchievTypeMegaRich])
	}
	return newAchievs
}
