// Package room — 房間難度定義（DAY-091）
// 4 個難度等級：初級/中級/高級/VIP
// 業界依據：Ocean King 系列多難度房間是 2026 年捕魚機標配
package room

// Difficulty 房間難度等級
type Difficulty string

const (
	DifficultyBeginner     Difficulty = "beginner"     // 初級：bet 1-5
	DifficultyIntermediate Difficulty = "intermediate" // 中級：bet 5-20
	DifficultyAdvanced     Difficulty = "advanced"     // 高級：bet 20-50
	DifficultyVIP          Difficulty = "vip"          // VIP：bet 50-200
)

// DifficultyDef 難度定義
type DifficultyDef struct {
	ID              Difficulty
	Name            string  // 顯示名稱（中文）
	Icon            string  // 圖示 emoji
	Color           string  // 主題色（hex）
	MinBetLevel     int     // 最低 bet 等級
	MaxBetLevel     int     // 最高 bet 等級
	MinBetCost      int     // 最低 bet 金幣
	MaxBetCost      int     // 最高 bet 金幣
	MaxPlayers      int     // 最大玩家數
	TargetSpeedMult float64 // 目標物速度倍率（高難度目標更快）
	RewardMult      float64 // 獎勵倍率（高難度獎勵更高）
	JackpotMult     float64 // Jackpot 貢獻倍率（高難度累積更快）
	BossHPMult      float64 // BOSS HP 倍率（高難度 BOSS 更硬）
	SpawnRateMult   float64 // 目標物生成速率倍率（高難度更多目標）
	RTPTarget       float64 // 目標 RTP
	EntryFee        int     // 進場費（VIP 房間需要）
	Description     string  // 房間描述
}

// Difficulties 所有難度定義
var Difficulties = map[Difficulty]*DifficultyDef{
	DifficultyBeginner: {
		ID:              DifficultyBeginner,
		Name:            "初心者",
		Icon:            "🌱",
		Color:           "#4CAF50",
		MinBetLevel:     1,
		MaxBetLevel:     4,
		MinBetCost:      10,
		MaxBetCost:      100,
		MaxPlayers:      16,
		TargetSpeedMult: 0.8,  // 目標物較慢，新手友善
		RewardMult:      1.0,
		JackpotMult:     0.5,  // Jackpot 累積較慢
		BossHPMult:      0.7,  // BOSS 較弱
		SpawnRateMult:   0.9,  // 目標物較少
		RTPTarget:       0.92,
		EntryFee:        0,
		Description:     "適合新手，目標移動較慢，輕鬆享受遊戲",
	},
	DifficultyIntermediate: {
		ID:              DifficultyIntermediate,
		Name:            "一般",
		Icon:            "⚔️",
		Color:           "#2196F3",
		MinBetLevel:     3,
		MaxBetLevel:     7,
		MinBetCost:      50,
		MaxBetCost:      500,
		MaxPlayers:      12,
		TargetSpeedMult: 1.0,  // 標準速度
		RewardMult:      1.2,  // 獎勵略高
		JackpotMult:     1.0,
		BossHPMult:      1.0,
		SpawnRateMult:   1.0,
		RTPTarget:       0.94,
		EntryFee:        0,
		Description:     "標準難度，平衡的挑戰與獎勵",
	},
	DifficultyAdvanced: {
		ID:              DifficultyAdvanced,
		Name:            "高手",
		Icon:            "🔥",
		Color:           "#FF9800",
		MinBetLevel:     6,
		MaxBetLevel:     9,
		MinBetCost:      200,
		MaxBetCost:      2000,
		MaxPlayers:      8,
		TargetSpeedMult: 1.3,  // 目標物更快
		RewardMult:      1.5,  // 獎勵更高
		JackpotMult:     2.0,  // Jackpot 累積更快
		BossHPMult:      1.5,  // BOSS 更強
		SpawnRateMult:   1.2,  // 更多目標物
		RTPTarget:       0.95,
		EntryFee:        0,
		Description:     "高難度高獎勵，目標移動快速，Jackpot 累積更快",
	},
	DifficultyVIP: {
		ID:              DifficultyVIP,
		Name:            "VIP",
		Icon:            "👑",
		Color:           "#9C27B0",
		MinBetLevel:     8,
		MaxBetLevel:     10,
		MinBetCost:      500,
		MaxBetCost:      10000,
		MaxPlayers:      4,   // 限制人數，更私密
		TargetSpeedMult: 1.6, // 目標物非常快
		RewardMult:      2.0, // 獎勵翻倍
		JackpotMult:     5.0, // Jackpot 累積最快
		BossHPMult:      2.0, // BOSS 最強
		SpawnRateMult:   1.5, // 最多目標物
		RTPTarget:       0.96,
		EntryFee:        10000, // 需要 10000 金幣進場費
		Description:     "頂級玩家專屬，超高獎勵，Jackpot 累積最快，限 4 人",
	},
}

// GetDifficulty 取得難度定義
func GetDifficulty(d Difficulty) *DifficultyDef {
	if def, ok := Difficulties[d]; ok {
		return def
	}
	return Difficulties[DifficultyBeginner]
}

// GetDifficultyByBetLevel 根據 bet 等級推薦難度
func GetDifficultyByBetLevel(betLevel int) Difficulty {
	switch {
	case betLevel <= 4:
		return DifficultyBeginner
	case betLevel <= 7:
		return DifficultyIntermediate
	case betLevel <= 9:
		return DifficultyAdvanced
	default:
		return DifficultyVIP
	}
}

// AllDifficulties 按順序回傳所有難度（用於 UI 顯示）
func AllDifficulties() []*DifficultyDef {
	return []*DifficultyDef{
		Difficulties[DifficultyBeginner],
		Difficulties[DifficultyIntermediate],
		Difficulties[DifficultyAdvanced],
		Difficulties[DifficultyVIP],
	}
}
