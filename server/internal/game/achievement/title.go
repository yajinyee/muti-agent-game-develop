// Package achievement — 稱號系統（DAY-068）
// 玩家達成特定成就後解鎖稱號，顯示在排行榜和玩家名稱旁
package achievement

// TitleID 稱號 ID
type TitleID string

const (
	TitleNovice       TitleID = "novice"        // 新手討伐者（預設）
	TitleHunter       TitleID = "hunter"        // 討伐獵人（擊破 5 個）
	TitleExpert       TitleID = "expert"        // 討伐達人（擊破 20 個）
	TitleMaster       TitleID = "master"        // 討伐高手（擊破 50 個）
	TitleLegend       TitleID = "legend"        // 討伐傳說（擊破 100 個）
	TitleBossSlayer   TitleID = "boss_slayer"   // BOSS 終結者（擊敗 BOSS）
	TitleBonusKing    TitleID = "bonus_king"    // 拔草之王（觸發 Bonus）
	TitleRichPlayer   TitleID = "rich_player"   // 小富翁（金幣 50k）
	TitleMillionaire  TitleID = "millionaire"   // 大富翁（金幣 100k）
	TitleMegaWinner   TitleID = "mega_winner"   // 超級大獎得主（50x+）
	TitleSpecialHunter TitleID = "special_hunter" // 特殊目標獵人
	TitleAllAround    TitleID = "all_around"    // 全能討伐者（解鎖 8 個以上成就）
)

// TitleDef 稱號定義
type TitleDef struct {
	ID          TitleID `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	Color       string  `json:"color"` // 顯示顏色（hex）
	Priority    int     `json:"priority"` // 優先級（越高越優先顯示）
}

// TitleDefinitions 所有稱號定義
var TitleDefinitions = map[TitleID]*TitleDef{
	TitleNovice: {
		ID:          TitleNovice,
		Name:        "新手討伐者",
		Description: "剛踏上討伐之路",
		Icon:        "🌱",
		Color:       "#AAAAAA",
		Priority:    0,
	},
	TitleHunter: {
		ID:          TitleHunter,
		Name:        "討伐獵人",
		Description: "累計擊破 5 個目標",
		Icon:        "🎯",
		Color:       "#88CC44",
		Priority:    10,
	},
	TitleExpert: {
		ID:          TitleExpert,
		Name:        "討伐達人",
		Description: "累計擊破 20 個目標",
		Icon:        "⚡",
		Color:       "#44AAFF",
		Priority:    20,
	},
	TitleMaster: {
		ID:          TitleMaster,
		Name:        "討伐高手",
		Description: "累計擊破 50 個目標",
		Icon:        "🏆",
		Color:       "#FF8800",
		Priority:    30,
	},
	TitleLegend: {
		ID:          TitleLegend,
		Name:        "討伐傳說",
		Description: "累計擊破 100 個目標",
		Icon:        "👑",
		Color:       "#FFD700",
		Priority:    50,
	},
	TitleBossSlayer: {
		ID:          TitleBossSlayer,
		Name:        "BOSS 終結者",
		Description: "首次擊敗 BOSS",
		Icon:        "🔥",
		Color:       "#FF4444",
		Priority:    40,
	},
	TitleBonusKing: {
		ID:          TitleBonusKing,
		Name:        "拔草之王",
		Description: "首次觸發 Bonus 遊戲",
		Icon:        "🌿",
		Color:       "#44FF88",
		Priority:    25,
	},
	TitleRichPlayer: {
		ID:          TitleRichPlayer,
		Name:        "小富翁",
		Description: "金幣達到 50,000",
		Icon:        "🪙",
		Color:       "#FFCC00",
		Priority:    35,
	},
	TitleMillionaire: {
		ID:          TitleMillionaire,
		Name:        "大富翁",
		Description: "金幣達到 100,000",
		Icon:        "💵",
		Color:       "#FFD700",
		Priority:    45,
	},
	TitleMegaWinner: {
		ID:          TitleMegaWinner,
		Name:        "超級大獎得主",
		Description: "獲得 50 倍以上超大獎",
		Icon:        "💎",
		Color:       "#CC44FF",
		Priority:    42,
	},
	TitleSpecialHunter: {
		ID:          TitleSpecialHunter,
		Name:        "特殊目標獵人",
		Description: "首次擊破特殊目標",
		Icon:        "✨",
		Color:       "#44FFFF",
		Priority:    22,
	},
	TitleAllAround: {
		ID:          TitleAllAround,
		Name:        "全能討伐者",
		Description: "解鎖 8 個以上成就",
		Icon:        "🌟",
		Color:       "#FF88FF",
		Priority:    55,
	},
}

// achievementToTitle 成就 → 稱號的對應關係
// 解鎖特定成就後，自動解鎖對應稱號
var achievementToTitle = map[AchievementID]TitleID{
	AchKill5:       TitleHunter,
	AchKill20:      TitleExpert,
	AchKill50:      TitleMaster,
	AchKill100:     TitleLegend,
	AchKillBoss:    TitleBossSlayer,
	AchBonus:       TitleBonusKing,
	AchCoins50k:    TitleRichPlayer,
	AchCoins100k:   TitleMillionaire,
	AchMegaWin:     TitleMegaWinner,
	AchKillSpecial: TitleSpecialHunter,
}

// TitleTracker 稱號追蹤器（每個玩家一個）
type TitleTracker struct {
	Unlocked    map[TitleID]bool // 已解鎖的稱號
	ActiveTitle TitleID          // 當前顯示的稱號（優先級最高的）
}

// NewTitleTracker 建立新的稱號追蹤器
func NewTitleTracker() *TitleTracker {
	return &TitleTracker{
		Unlocked:    map[TitleID]bool{TitleNovice: true},
		ActiveTitle: TitleNovice,
	}
}

// OnAchievementUnlocked 成就解鎖時呼叫，回傳是否解鎖了新稱號
func (t *TitleTracker) OnAchievementUnlocked(achID AchievementID, totalUnlocked int) *TitleDef {
	// 檢查是否有對應稱號
	if titleID, ok := achievementToTitle[achID]; ok {
		if !t.Unlocked[titleID] {
			t.Unlocked[titleID] = true
			t.recalcActiveTitle()
			return TitleDefinitions[titleID]
		}
	}

	// 檢查全能討伐者（8 個以上成就）
	if totalUnlocked >= 8 && !t.Unlocked[TitleAllAround] {
		t.Unlocked[TitleAllAround] = true
		t.recalcActiveTitle()
		return TitleDefinitions[TitleAllAround]
	}

	return nil
}

// recalcActiveTitle 重新計算當前顯示稱號（選優先級最高的）
func (t *TitleTracker) recalcActiveTitle() {
	bestPriority := -1
	bestTitle := TitleNovice
	for titleID := range t.Unlocked {
		def, ok := TitleDefinitions[titleID]
		if !ok {
			continue
		}
		if def.Priority > bestPriority {
			bestPriority = def.Priority
			bestTitle = titleID
		}
	}
	t.ActiveTitle = bestTitle
}

// GetActiveTitle 取得當前顯示稱號定義
func (t *TitleTracker) GetActiveTitle() *TitleDef {
	def, ok := TitleDefinitions[t.ActiveTitle]
	if !ok {
		return TitleDefinitions[TitleNovice]
	}
	return def
}

// GetUnlockedTitles 取得所有已解鎖稱號
func (t *TitleTracker) GetUnlockedTitles() []*TitleDef {
	result := make([]*TitleDef, 0, len(t.Unlocked))
	for titleID := range t.Unlocked {
		if def, ok := TitleDefinitions[titleID]; ok {
			result = append(result, def)
		}
	}
	return result
}

// SetActiveTitle 手動設定顯示稱號（玩家自選）
func (t *TitleTracker) SetActiveTitle(titleID TitleID) bool {
	if !t.Unlocked[titleID] {
		return false
	}
	t.ActiveTitle = titleID
	return true
}

// LoadState 從持久化資料恢復稱號狀態（DAY-100）
func (t *TitleTracker) LoadState(unlockedTitles []TitleID, activeTitle TitleID) {
	for _, id := range unlockedTitles {
		t.Unlocked[id] = true
	}
	if activeTitle != "" {
		if t.Unlocked[activeTitle] {
			t.ActiveTitle = activeTitle
		} else {
			t.recalcActiveTitle()
		}
	} else {
		t.recalcActiveTitle()
	}
}
