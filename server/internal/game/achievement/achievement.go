// Package achievement 成就系統
package achievement

import "time"

// AchievementID 成就 ID
type AchievementID string

const (
	// 擊破類
	AchFirstKill    AchievementID = "first_kill"    // 首次擊破目標
	AchKill5        AchievementID = "kill_5"         // 擊破 5 個目標
	AchKill20       AchievementID = "kill_20"        // 擊破 20 個目標
	AchKill50       AchievementID = "kill_50"        // 擊破 50 個目標
	AchKill100      AchievementID = "kill_100"       // 擊破 100 個目標
	// 特殊目標類
	AchKillSpecial  AchievementID = "kill_special"   // 首次擊破特殊目標（T101-T105）
	AchKillBoss     AchievementID = "kill_boss"      // 首次擊敗 BOSS
	// 獎勵類
	AchBigWin       AchievementID = "big_win"        // 首次獲得大獎（≥20x）
	AchMegaWin      AchievementID = "mega_win"       // 首次獲得超大獎（≥50x）
	AchBonus        AchievementID = "first_bonus"    // 首次觸發 Bonus
	// 金幣類
	AchCoins50k     AchievementID = "coins_50k"      // 金幣達到 50,000
	AchCoins100k    AchievementID = "coins_100k"     // 金幣達到 100,000
)

// Achievement 成就定義
type Achievement struct {
	ID          AchievementID `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Icon        string        `json:"icon"` // emoji 圖示
	Type        string        `json:"type"` // 成就類型：normal/boss/bonus/special
}

// AchievementUnlock 成就解鎖記錄
type AchievementUnlock struct {
	ID          AchievementID `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Icon        string        `json:"icon"`
	Type        string        `json:"type"` // 成就類型（傳給 Client 用於 UI 顏色）
	UnlockedAt  time.Time     `json:"unlocked_at"`
}

// Definitions 所有成就定義
var Definitions = map[AchievementID]*Achievement{
	AchFirstKill: {
		ID:          AchFirstKill,
		Name:        "初次討伐",
		Description: "首次擊破目標",
		Icon:        "⚔️",
		Type:        "normal",
	},
	AchKill5: {
		ID:          AchKill5,
		Name:        "討伐新手",
		Description: "累計擊破 5 個目標",
		Icon:        "🌟",
		Type:        "normal",
	},
	AchKill20: {
		ID:          AchKill20,
		Name:        "討伐達人",
		Description: "累計擊破 20 個目標",
		Icon:        "💫",
		Type:        "normal",
	},
	AchKill50: {
		ID:          AchKill50,
		Name:        "討伐高手",
		Description: "累計擊破 50 個目標",
		Icon:        "🏆",
		Type:        "normal",
	},
	AchKill100: {
		ID:          AchKill100,
		Name:        "討伐傳說",
		Description: "累計擊破 100 個目標",
		Icon:        "👑",
		Type:        "special",
	},
	AchKillSpecial: {
		ID:          AchKillSpecial,
		Name:        "特殊目標獵人",
		Description: "首次擊破特殊目標",
		Icon:        "✨",
		Type:        "special",
	},
	AchKillBoss: {
		ID:          AchKillBoss,
		Name:        "BOSS 終結者",
		Description: "首次擊敗 BOSS",
		Icon:        "🔥",
		Type:        "boss",
	},
	AchBigWin: {
		ID:          AchBigWin,
		Name:        "大獎得主",
		Description: "首次獲得 20 倍以上大獎",
		Icon:        "💰",
		Type:        "normal",
	},
	AchMegaWin: {
		ID:          AchMegaWin,
		Name:        "超級大獎",
		Description: "首次獲得 50 倍以上超大獎",
		Icon:        "💎",
		Type:        "special",
	},
	AchBonus: {
		ID:          AchBonus,
		Name:        "瘋狂拔草",
		Description: "首次觸發 Bonus 遊戲",
		Icon:        "🌿",
		Type:        "bonus",
	},
	AchCoins50k: {
		ID:          AchCoins50k,
		Name:        "小富翁",
		Description: "金幣達到 50,000",
		Icon:        "🪙",
		Type:        "normal",
	},
	AchCoins100k: {
		ID:          AchCoins100k,
		Name:        "大富翁",
		Description: "金幣達到 100,000",
		Icon:        "💵",
		Type:        "special",
	},
}

// Tracker 成就追蹤器（每個玩家一個）
type Tracker struct {
	Unlocked map[AchievementID]time.Time // 已解鎖的成就及解鎖時間
}

// NewTracker 建立新的成就追蹤器
func NewTracker() *Tracker {
	return &Tracker{
		Unlocked: make(map[AchievementID]time.Time),
	}
}

// TryUnlock 嘗試解鎖成就，若成功回傳成就定義，否則回傳 nil
func (t *Tracker) TryUnlock(id AchievementID) *AchievementUnlock {
	if _, already := t.Unlocked[id]; already {
		return nil // 已解鎖，不重複觸發
	}
	def, ok := Definitions[id]
	if !ok {
		return nil
	}
	now := time.Now()
	t.Unlocked[id] = now
	return &AchievementUnlock{
		ID:          def.ID,
		Name:        def.Name,
		Description: def.Description,
		Icon:        def.Icon,
		Type:        def.Type,
		UnlockedAt:  now,
	}
}

// IsUnlocked 檢查成就是否已解鎖
func (t *Tracker) IsUnlocked(id AchievementID) bool {
	_, ok := t.Unlocked[id]
	return ok
}

// UnlockedList 取得所有已解鎖成就的清單
func (t *Tracker) UnlockedList() []AchievementUnlock {
	result := make([]AchievementUnlock, 0, len(t.Unlocked))
	for id, unlockedAt := range t.Unlocked {
		def := Definitions[id]
		if def == nil {
			continue
		}
		result = append(result, AchievementUnlock{
			ID:          def.ID,
			Name:        def.Name,
			Description: def.Description,
			Icon:        def.Icon,
			Type:        def.Type,
			UnlockedAt:  unlockedAt,
		})
	}
	return result
}
