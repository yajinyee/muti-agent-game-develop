// Package dailybonus — 登入里程碑獎勵系統（DAY-107）
// 業界依據：ilogos.biz（2026）確認 gamified login streaks 讓留存率提升 75%
// 設計：連續登入達到里程碑時給予特殊獎勵（寶箱 + 金幣 + 稱號）
package dailybonus

// MilestoneRewardType 里程碑獎勵類型
type MilestoneRewardType string

const (
	MilestoneRewardCoins      MilestoneRewardType = "coins"       // 金幣獎勵
	MilestoneRewardMysteryBox MilestoneRewardType = "mystery_box" // 神秘寶箱
	MilestoneRewardTitle      MilestoneRewardType = "title"       // 特殊稱號解鎖
)

// MilestoneReward 里程碑獎勵項目
type MilestoneReward struct {
	Type     MilestoneRewardType `json:"type"`
	Amount   int                 `json:"amount"`   // 金幣數量 / 寶箱數量
	Rarity   string              `json:"rarity"`   // 寶箱稀有度（common/rare/epic/legendary）
	TitleID  string              `json:"title_id"` // 稱號 ID（type=title 時使用）
}

// Milestone 登入里程碑定義
type Milestone struct {
	Days        int               `json:"days"`         // 連續登入天數
	Name        string            `json:"name"`         // 里程碑名稱
	Description string            `json:"description"`  // 里程碑描述
	Icon        string            `json:"icon"`         // 圖示 emoji
	Color       string            `json:"color"`        // 顏色（hex）
	Rewards     []MilestoneReward `json:"rewards"`      // 獎勵列表
}

// milestones 里程碑定義表（按天數排序）
var milestones = []Milestone{
	{
		Days:        3,
		Name:        "初心者",
		Description: "連續登入 3 天",
		Icon:        "🌱",
		Color:       "#4CAF50",
		Rewards: []MilestoneReward{
			{Type: MilestoneRewardCoins, Amount: 3000},
			{Type: MilestoneRewardMysteryBox, Amount: 1, Rarity: "common"},
		},
	},
	{
		Days:        7,
		Name:        "一週勇士",
		Description: "連續登入 7 天",
		Icon:        "⚔️",
		Color:       "#2196F3",
		Rewards: []MilestoneReward{
			{Type: MilestoneRewardCoins, Amount: 8000},
			{Type: MilestoneRewardMysteryBox, Amount: 1, Rarity: "rare"},
		},
	},
	{
		Days:        14,
		Name:        "兩週老手",
		Description: "連續登入 14 天",
		Icon:        "🔥",
		Color:       "#FF9800",
		Rewards: []MilestoneReward{
			{Type: MilestoneRewardCoins, Amount: 20000},
			{Type: MilestoneRewardMysteryBox, Amount: 1, Rarity: "epic"},
			{Type: MilestoneRewardTitle, TitleID: "streak_veteran"},
		},
	},
	{
		Days:        30,
		Name:        "月度傳說",
		Description: "連續登入 30 天",
		Icon:        "👑",
		Color:       "#FFD700",
		Rewards: []MilestoneReward{
			{Type: MilestoneRewardCoins, Amount: 50000},
			{Type: MilestoneRewardMysteryBox, Amount: 2, Rarity: "legendary"},
			{Type: MilestoneRewardTitle, TitleID: "streak_legend"},
		},
	},
	{
		Days:        60,
		Name:        "兩月霸主",
		Description: "連續登入 60 天",
		Icon:        "💎",
		Color:       "#9C27B0",
		Rewards: []MilestoneReward{
			{Type: MilestoneRewardCoins, Amount: 120000},
			{Type: MilestoneRewardMysteryBox, Amount: 3, Rarity: "legendary"},
			{Type: MilestoneRewardTitle, TitleID: "streak_master"},
		},
	},
	{
		Days:        100,
		Name:        "百日神話",
		Description: "連續登入 100 天",
		Icon:        "🌟",
		Color:       "#FF5722",
		Rewards: []MilestoneReward{
			{Type: MilestoneRewardCoins, Amount: 300000},
			{Type: MilestoneRewardMysteryBox, Amount: 5, Rarity: "legendary"},
			{Type: MilestoneRewardTitle, TitleID: "streak_myth"},
		},
	},
}

// CheckMilestone 檢查是否達到里程碑
// 回傳達到的里程碑（nil = 未達到）
func CheckMilestone(newStreak int) *Milestone {
	for i := range milestones {
		if milestones[i].Days == newStreak {
			return &milestones[i]
		}
	}
	return nil
}

// GetAllMilestones 取得所有里程碑定義（用於 Client 顯示進度）
func GetAllMilestones() []Milestone {
	result := make([]Milestone, len(milestones))
	copy(result, milestones)
	return result
}

// GetNextMilestone 取得下一個里程碑（用於 Client 顯示「距離下一個里程碑還有 X 天」）
func GetNextMilestone(currentStreak int) *Milestone {
	for i := range milestones {
		if milestones[i].Days > currentStreak {
			return &milestones[i]
		}
	}
	return nil // 已達到所有里程碑
}
