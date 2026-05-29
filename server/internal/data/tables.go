// Package data — 遊戲資料表
// target-design-agent + balance-agent 負責維護
package data

// ── 角色定義 ─────────────────────────────────────────────────

type CharacterDef struct {
	ID              string
	Name            string
	MinBetLevel     int
	MaxBetLevel     int
	KillModifier    float64
	LaborModifier   float64
	FireRateModifier float64
}

var Characters = []CharacterDef{
	{ID: "chiikawa", Name: "Chiikawa", MinBetLevel: 1, MaxBetLevel: 3, KillModifier: 1.00, LaborModifier: 1.10, FireRateModifier: 1.00},
	{ID: "hachiware", Name: "Hachiware", MinBetLevel: 4, MaxBetLevel: 7, KillModifier: 1.00, LaborModifier: 1.00, FireRateModifier: 1.08},
	{ID: "usagi", Name: "Usagi", MinBetLevel: 8, MaxBetLevel: 10, KillModifier: 0.98, LaborModifier: 1.00, FireRateModifier: 1.20},
}

// ── 投注等級表 ────────────────────────────────────────────────

type BetLevel struct {
	Level           int
	CharacterID     string
	BetCost         int
	AttackPower     int
	FireRate        float64
	ProjectileSpeed float64
}

var BetLevels = []BetLevel{
	{1, "chiikawa", 1, 1, 2.0, 700},
	{2, "chiikawa", 2, 2, 2.0, 720},
	{3, "chiikawa", 3, 3, 2.1, 740},
	{4, "hachiware", 5, 5, 2.2, 780},
	{5, "hachiware", 10, 10, 2.3, 800},
	{6, "hachiware", 20, 20, 2.4, 820},
	{7, "hachiware", 30, 30, 2.5, 850},
	{8, "usagi", 50, 50, 2.7, 900},
	{9, "usagi", 80, 80, 2.9, 940},
	{10, "usagi", 100, 100, 3.0, 980},
}

func GetBetLevel(level int) BetLevel {
	if level < 1 {
		level = 1
	}
	if level > 10 {
		level = 10
	}
	return BetLevels[level-1]
}

func GetCharacter(id string) CharacterDef {
	for _, c := range Characters {
		if c.ID == id {
			return c
		}
	}
	return Characters[0]
}

// ── 目標物定義 ────────────────────────────────────────────────

type TargetType string

const (
	TypeBasic   TargetType = "basic"
	TypeSpecial TargetType = "special"
	TypeBoss    TargetType = "boss"
)

type Behavior string

const (
	BehaviorLinear Behavior = "linear"
	BehaviorSink   Behavior = "sink"
	BehaviorFlee   Behavior = "flee"
	BehaviorFast   Behavior = "fast"
)

type TargetDef struct {
	ID           string
	Name         string
	Type         TargetType
	Multiplier   float64 // 基礎倍率（流星用 MinMult/MaxMult）
	MinMult      float64 // 流星最小倍率
	MaxMult      float64 // 流星最大倍率
	HP           int
	SpawnWeight  int
	Speed        float64 // px/s
	Lifetime     float64 // 秒
	LaborGain    int
	Behavior     Behavior
	DiffFactor   float64 // 保底難度係數
}

var Targets = []TargetDef{
	// ── 基礎目標（2x-10x）────────────────────────────────────
	{ID: "T001", Name: "像素雜草", Type: TypeBasic, Multiplier: 2, HP: 3, SpawnWeight: 180, Speed: 0, Lifetime: 20, LaborGain: 1, Behavior: BehaviorSink, DiffFactor: 0.4},
	{ID: "T002", Name: "綠色小蟲", Type: TypeBasic, Multiplier: 3, HP: 5, SpawnWeight: 160, Speed: 40, Lifetime: 18, LaborGain: 1, Behavior: BehaviorLinear, DiffFactor: 0.4},
	{ID: "T003", Name: "紅色小蟲", Type: TypeBasic, Multiplier: 5, HP: 8, SpawnWeight: 130, Speed: 55, Lifetime: 16, LaborGain: 1, Behavior: BehaviorLinear, DiffFactor: 0.4},
	{ID: "T004", Name: "藍色小蟲", Type: TypeBasic, Multiplier: 6, HP: 10, SpawnWeight: 110, Speed: 65, Lifetime: 15, LaborGain: 2, Behavior: BehaviorLinear, DiffFactor: 0.4},
	{ID: "T005", Name: "會走路的布丁", Type: TypeBasic, Multiplier: 8, HP: 16, SpawnWeight: 90, Speed: 35, Lifetime: 20, LaborGain: 2, Behavior: BehaviorLinear, DiffFactor: 0.4},
	{ID: "T006", Name: "巨大蘑菇", Type: TypeBasic, Multiplier: 10, HP: 22, SpawnWeight: 70, Speed: 25, Lifetime: 22, LaborGain: 3, Behavior: BehaviorLinear, DiffFactor: 0.4},

	// ── 特殊目標（15x-50x）───────────────────────────────────
	{ID: "T101", Name: "擬態型怪物", Type: TypeSpecial, MinMult: 15, MaxMult: 30, HP: 35, SpawnWeight: 35, Speed: 50, Lifetime: 14, LaborGain: 5, Behavior: BehaviorLinear, DiffFactor: 0.7},
	{ID: "T102", Name: "寶箱怪", Type: TypeSpecial, Multiplier: 25, HP: 55, SpawnWeight: 22, Speed: 70, Lifetime: 10, LaborGain: 6, Behavior: BehaviorFlee, DiffFactor: 0.7},
	{ID: "T103", Name: "流星", Type: TypeSpecial, MinMult: 20, MaxMult: 50, HP: 20, SpawnWeight: 18, Speed: 220, Lifetime: 4, LaborGain: 5, Behavior: BehaviorFast, DiffFactor: 0.8},
	{ID: "T104", Name: "金色雜草", Type: TypeSpecial, Multiplier: 30, HP: 45, SpawnWeight: 12, Speed: 0, Lifetime: 8, LaborGain: 15, Behavior: BehaviorSink, DiffFactor: 0.7},
	{ID: "T105", Name: "巨大金幣魚", Type: TypeSpecial, Multiplier: 50, HP: 90, SpawnWeight: 8, Speed: 80, Lifetime: 8, LaborGain: 10, Behavior: BehaviorLinear, DiffFactor: 0.8},

	// ── 進階特殊目標（60x-120x）─────────────────────────────
	// T106 幸運連鎖閃電魚：擊破後觸發連鎖閃電，攻擊附近 3 條魚 HP -50%
	{ID: "T106", Name: "幸運連鎖閃電魚", Type: TypeSpecial, Multiplier: 60, HP: 80, SpawnWeight: 6, Speed: 90, Lifetime: 12, LaborGain: 12, Behavior: BehaviorLinear, DiffFactor: 0.9},
	// T107 幸運螃蟹魚雷：擊破後在場上隨機位置觸發 3 次 AOE 爆炸，每次 HP -40%
	{ID: "T107", Name: "幸運螃蟹魚雷", Type: TypeSpecial, Multiplier: 70, HP: 90, SpawnWeight: 5, Speed: 60, Lifetime: 14, LaborGain: 14, Behavior: BehaviorLinear, DiffFactor: 0.9},
	// T108 幸運渦旋海葵：擊破後全場所有目標被吸向中心，HP -30%，持續 5 秒
	{ID: "T108", Name: "幸運渦旋海葵", Type: TypeSpecial, Multiplier: 80, HP: 100, SpawnWeight: 4, Speed: 0, Lifetime: 16, LaborGain: 16, Behavior: BehaviorSink, DiffFactor: 0.9},
	// T109 幸運黃金龍魚：擊破後觸發雙環輪盤，最高 ×350 倍率
	{ID: "T109", Name: "幸運黃金龍魚", Type: TypeSpecial, MinMult: 80, MaxMult: 350, HP: 120, SpawnWeight: 3, Speed: 70, Lifetime: 14, LaborGain: 20, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T110 幸運雷霆龍蝦：擊破後 15 秒自動射擊模式，全場 AOE 傷害
	{ID: "T110", Name: "幸運雷霆龍蝦", Type: TypeSpecial, Multiplier: 100, HP: 140, SpawnWeight: 2, Speed: 50, Lifetime: 16, LaborGain: 25, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T111 幸運覺醒鳳凰魚：擊破後觸發覺醒模式，下 5 次攻擊每次 Power Up 6x-10x
	{ID: "T111", Name: "幸運覺醒鳳凰魚", Type: TypeSpecial, Multiplier: 90, HP: 110, SpawnWeight: 4, Speed: 75, Lifetime: 14, LaborGain: 18, Behavior: BehaviorLinear, DiffFactor: 0.95},
	// T112 幸運全場震盪魚：擊破後全場 HP -35%，觸發玩家 10 秒攻擊力 ×2.0
	{ID: "T112", Name: "幸運全場震盪魚", Type: TypeSpecial, Multiplier: 75, HP: 95, SpawnWeight: 5, Speed: 65, Lifetime: 13, LaborGain: 15, Behavior: BehaviorLinear, DiffFactor: 0.9},
	// T113 幸運鑽頭魚雷魚：擊破後發射鑽頭魚雷穿透最多 5 個目標（HP -60%），終點爆炸 AOE r=180px
	{ID: "T113", Name: "幸運鑽頭魚雷魚", Type: TypeSpecial, Multiplier: 85, HP: 105, SpawnWeight: 4, Speed: 80, Lifetime: 13, LaborGain: 17, Behavior: BehaviorLinear, DiffFactor: 0.95},
	// T114 幸運時間凍結魚：擊破後全場凍結 8 秒（傷害 ×1.8），凍結結束 HP -25%
	{ID: "T114", Name: "幸運時間凍結魚", Type: TypeSpecial, Multiplier: 95, HP: 115, SpawnWeight: 3, Speed: 55, Lifetime: 15, LaborGain: 19, Behavior: BehaviorLinear, DiffFactor: 0.95},
	// T115 幸運連鎖爆炸魚：擊破後 12 秒連鎖爆炸模式，每次擊破觸發 AOE r=120px HP -30%
	{ID: "T115", Name: "幸運連鎖爆炸魚", Type: TypeSpecial, Multiplier: 80, HP: 100, SpawnWeight: 4, Speed: 70, Lifetime: 12, LaborGain: 16, Behavior: BehaviorLinear, DiffFactor: 0.9},

	// ── DAY-295 新增特殊目標（120x-200x）────────────────────
	// T116 幸運千龍王輪盤魚：擊破後觸發千龍王輪盤，最高 1000x Mega Win
	{ID: "T116", Name: "幸運千龍王輪盤魚", Type: TypeSpecial, MinMult: 120, MaxMult: 1000, HP: 160, SpawnWeight: 2, Speed: 60, Lifetime: 16, LaborGain: 30, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T117 幸運龍力散彈魚：擊破後 8 方向散彈攻擊，每方向 HP -40%
	{ID: "T117", Name: "幸運龍力散彈魚", Type: TypeSpecial, Multiplier: 120, HP: 145, SpawnWeight: 3, Speed: 85, Lifetime: 13, LaborGain: 24, Behavior: BehaviorLinear, DiffFactor: 0.95},
	// T118 幸運火箭砲魚：擊破後召喚 3 枚火箭砲，每枚 AOE r=200px HP -50%
	{ID: "T118", Name: "幸運火箭砲魚", Type: TypeSpecial, Multiplier: 130, HP: 155, SpawnWeight: 3, Speed: 75, Lifetime: 14, LaborGain: 26, Behavior: BehaviorLinear, DiffFactor: 0.95},
	// T119 幸運深海漩渦魚：擊破後免費深海漩渦，全場 HP -50%，持續 6 秒
	{ID: "T119", Name: "幸運深海漩渦魚", Type: TypeSpecial, Multiplier: 150, HP: 175, SpawnWeight: 2, Speed: 45, Lifetime: 18, LaborGain: 30, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T120 幸運吸血鬼魚：擊破後進入吸血模式，每次擊破吸收倍率，最高 ×5 模式
	{ID: "T120", Name: "幸運吸血鬼魚", Type: TypeSpecial, Multiplier: 110, HP: 135, SpawnWeight: 3, Speed: 95, Lifetime: 12, LaborGain: 22, Behavior: BehaviorFlee, DiffFactor: 0.95},

	// ── DAY-296 新增特殊目標（110x-160x）────────────────────
	// T121 幸運鏡像魚：擊破後觸發鏡像模式，下 3 次攻擊自動複製一次
	{ID: "T121", Name: "幸運鏡像魚", Type: TypeSpecial, Multiplier: 110, HP: 130, SpawnWeight: 3, Speed: 70, Lifetime: 14, LaborGain: 22, Behavior: BehaviorLinear, DiffFactor: 0.95},
	// T122 幸運黃金雨魚：擊破後觸發黃金雨，全場生成 8-12 個黃金幣可收集
	{ID: "T122", Name: "幸運黃金雨魚", Type: TypeSpecial, Multiplier: 120, HP: 145, SpawnWeight: 3, Speed: 55, Lifetime: 16, LaborGain: 24, Behavior: BehaviorLinear, DiffFactor: 0.95},
	// T123 幸運冰凍炸彈魚：擊破後投擲冰凍炸彈，凍結 3 秒後爆炸 HP -60%
	{ID: "T123", Name: "幸運冰凍炸彈魚", Type: TypeSpecial, Multiplier: 130, HP: 155, SpawnWeight: 3, Speed: 65, Lifetime: 14, LaborGain: 26, Behavior: BehaviorLinear, DiffFactor: 0.95},
	// T124 幸運雷暴魚：擊破後觸發雷暴，10 秒內 6-7 道閃電隨機落下
	{ID: "T124", Name: "幸運雷暴魚", Type: TypeSpecial, Multiplier: 140, HP: 165, SpawnWeight: 2, Speed: 80, Lifetime: 13, LaborGain: 28, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T125 幸運大轉盤魚：擊破後觸發幸運大轉盤，8 格隨機獎勵
	{ID: "T125", Name: "幸運大轉盤魚", Type: TypeSpecial, Multiplier: 160, HP: 185, SpawnWeight: 2, Speed: 50, Lifetime: 18, LaborGain: 32, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-301 新增特殊目標（150x-200x）────────────────────
	// T126 幸運進階 Jackpot 魚：擊破後觸發四層 Jackpot 抽獎（Mini/Minor/Major/Grand）
	{ID: "T126", Name: "幸運進階Jackpot魚", Type: TypeSpecial, Multiplier: 150, HP: 180, SpawnWeight: 2, Speed: 55, Lifetime: 18, LaborGain: 30, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T127 幸運全服合作魚：擊破後觸發全服合作挑戰，20 秒內全服共同擊破目標點數
	{ID: "T127", Name: "幸運全服合作魚", Type: TypeSpecial, Multiplier: 140, HP: 170, SpawnWeight: 2, Speed: 60, Lifetime: 16, LaborGain: 28, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T128 幸運時間扭曲魚：擊破後全場目標速度 ×0.3，持續 10 秒，傷害 ×2.0
	{ID: "T128", Name: "幸運時間扭曲魚", Type: TypeSpecial, Multiplier: 130, HP: 160, SpawnWeight: 3, Speed: 70, Lifetime: 14, LaborGain: 26, Behavior: BehaviorLinear, DiffFactor: 0.95},
	// T129 幸運連鎖隕石魚：擊破後觸發連鎖隕石雨，5 顆隕石依序落下，每顆命中觸發連鎖擴大
	{ID: "T129", Name: "幸運連鎖隕石魚", Type: TypeSpecial, Multiplier: 160, HP: 190, SpawnWeight: 2, Speed: 65, Lifetime: 16, LaborGain: 32, Behavior: BehaviorLinear, DiffFactor: 1.0},

	// ── DAY-303 新增特殊目標（170x）────────────────────────────
	// T130 幸運崩潰魚：擊破後觸發崩潰倍率，倍率每 0.5 秒 +0.3x，玩家可隨時收割，崩潰前收割 ≥5.0x 觸發完美收割
	{ID: "T130", Name: "幸運崩潰魚", Type: TypeSpecial, Multiplier: 170, HP: 200, SpawnWeight: 2, Speed: 70, Lifetime: 16, LaborGain: 34, Behavior: BehaviorLinear, DiffFactor: 1.0},

	// ── DAY-304 新增特殊目標（180x-250x）────────────────────────
	// T131 幸運電鰻魚：擊破後持續放電 12 秒，每 1.5 秒電擊最近 3 條魚（HP -25%），連鎖加速，累積 ≥8 次 → 超級放電全服 ×2.5
	{ID: "T131", Name: "幸運電鰻魚", Type: TypeSpecial, Multiplier: 180, HP: 210, SpawnWeight: 2, Speed: 75, Lifetime: 16, LaborGain: 36, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T132 幸運巨型安康魚：擊破後誘餌吸引 5 秒（傷害 ×1.8），然後電擊爆炸全場 HP -30%，命中 ≥8 → 完美誘捕全服 ×2.8
	{ID: "T132", Name: "幸運巨型安康魚", Type: TypeSpecial, Multiplier: 190, HP: 220, SpawnWeight: 2, Speed: 60, Lifetime: 18, LaborGain: 38, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T133 幸運黑洞魚：擊破後黑洞吸引 8 秒（速度 ×0.2），然後坍縮全場 HP -50%，命中 ≥10 → 奇點爆發全服 ×3.0
	{ID: "T133", Name: "幸運黑洞魚", Type: TypeSpecial, Multiplier: 200, HP: 230, SpawnWeight: 1, Speed: 50, Lifetime: 20, LaborGain: 40, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T134 幸運賞金獵人魚：擊破後標記 3 個賞金目標（HP -20%），30 秒內全部擊破 → 完美賞金全服 ×3.5
	{ID: "T134", Name: "幸運賞金獵人魚", Type: TypeSpecial, Multiplier: 220, HP: 240, SpawnWeight: 1, Speed: 80, Lifetime: 16, LaborGain: 44, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T135 幸運海嘯魚：擊破後三波海嘯（HP -20%/-30%/-40%），三波命中 ≥5 → 完美海嘯全服 ×3.2
	{ID: "T135", Name: "幸運海嘯魚", Type: TypeSpecial, Multiplier: 250, HP: 260, SpawnWeight: 1, Speed: 65, Lifetime: 18, LaborGain: 50, Behavior: BehaviorLinear, DiffFactor: 1.0},

	// ── DAY-305 新增特殊目標（260x-350x）────────────────────────
	// T136 幸運龍怒蓄積魚 v2：擊破後 30 秒蓄積怒氣，每次射擊 +1（最高 30），爆發隕石雨，怒氣 ≥20 → 完美龍怒全服 ×3.5
	{ID: "T136", Name: "幸運龍怒蓄積魚", Type: TypeSpecial, Multiplier: 260, HP: 280, SpawnWeight: 1, Speed: 70, Lifetime: 18, LaborGain: 52, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T137 幸運座頭鯨魚：擊破後鯨歌共鳴 4 波（每波 HP -15%，命中越多下波傷害越高），命中 ≥20 → 完美鯨歌全服 ×3.0
	{ID: "T137", Name: "幸運座頭鯨魚", Type: TypeSpecial, Multiplier: 280, HP: 300, SpawnWeight: 1, Speed: 55, Lifetime: 20, LaborGain: 56, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T138 幸運傳說龍魚：擊破後傳說龍降臨 15 秒，每 3 秒噴火（HP -35%），4 次全部命中 ≥3 → 傳說龍怒全服 ×4.0
	{ID: "T138", Name: "幸運傳說龍魚", Type: TypeSpecial, Multiplier: 300, HP: 320, SpawnWeight: 1, Speed: 60, Lifetime: 20, LaborGain: 60, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T139 幸運公會戰魚：擊破後全服 30 秒積分挑戰，達成積分目標 → 公會勝利全服 ×4.5
	{ID: "T139", Name: "幸運公會戰魚", Type: TypeSpecial, Multiplier: 320, HP: 340, SpawnWeight: 1, Speed: 65, Lifetime: 18, LaborGain: 64, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T140 幸運品質魚：擊破後品質鑑定（Common/Rare/Epic/Legendary），Legendary → 傳說品質全服 ×5.0
	{ID: "T140", Name: "幸運品質魚", Type: TypeSpecial, Multiplier: 350, HP: 360, SpawnWeight: 1, Speed: 70, Lifetime: 16, LaborGain: 70, Behavior: BehaviorLinear, DiffFactor: 1.0},

	// ── DAY-306 新增特殊目標（360x-450x）────────────────────────
	// T141 幸運龍捲風魚：擊破後龍捲風橫掃 10 秒（每 2 秒 HP -40%），擊破 ≥8 → 完美龍捲風全服 ×3.8
	{ID: "T141", Name: "幸運龍捲風魚", Type: TypeSpecial, Multiplier: 360, HP: 380, SpawnWeight: 1, Speed: 75, Lifetime: 16, LaborGain: 72, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T142 幸運地震魚：擊破後三波地震（HP -25%/-35%/-45%），三波命中 ≥12 → 完美地震全服 ×4.0
	{ID: "T142", Name: "幸運地震魚", Type: TypeSpecial, Multiplier: 380, HP: 400, SpawnWeight: 1, Speed: 60, Lifetime: 18, LaborGain: 76, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T143 幸運火山魚：擊破後 10 顆熔岩彈（HP -35%），全部命中 → 完美火山全服 ×4.2
	{ID: "T143", Name: "幸運火山魚", Type: TypeSpecial, Multiplier: 400, HP: 420, SpawnWeight: 1, Speed: 65, Lifetime: 18, LaborGain: 80, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T144 幸運星際魚：擊破後 8 方向光束（HP -30%），命中 ≥16 → 完美星際全服 ×4.5
	{ID: "T144", Name: "幸運星際魚", Type: TypeSpecial, Multiplier: 420, HP: 440, SpawnWeight: 1, Speed: 70, Lifetime: 16, LaborGain: 84, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T145 幸運神龍魚：擊破後神龍降臨 20 秒（每 4 秒爪擊 HP -50%），5 次全部命中 ≥5 → 神龍完美全服 ×5.0
	{ID: "T145", Name: "幸運神龍魚", Type: TypeSpecial, Multiplier: 450, HP: 480, SpawnWeight: 1, Speed: 50, Lifetime: 22, LaborGain: 90, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-307 新增特殊目標（480x-600x）────────────────────────
	// T146 幸運量子魚：擊破後量子觀測（50% 機率 HP -60%），觀測 ≥10 → 量子坍縮全服 ×5.5
	{ID: "T146", Name: "幸運量子魚", Type: TypeSpecial, Multiplier: 480, HP: 500, SpawnWeight: 1, Speed: 55, Lifetime: 20, LaborGain: 96, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T147 幸運超新星魚：擊破後全場 HP -70%，5 秒倍率 ×3.0，命中 ≥8 → 超新星完美全服 ×5.5
	{ID: "T147", Name: "幸運超新星魚", Type: TypeSpecial, Multiplier: 500, HP: 520, SpawnWeight: 1, Speed: 60, Lifetime: 20, LaborGain: 100, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T148 幸運無限魚：擊破後 20 秒無限累積倍率（每次擊破 +1.0x），≥20x → 無限完美全服 ×6.0
	{ID: "T148", Name: "幸運無限魚", Type: TypeSpecial, Multiplier: 520, HP: 540, SpawnWeight: 1, Speed: 50, Lifetime: 22, LaborGain: 104, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T149 幸運創世魚：擊破後全場目標 HP 歸零（每個獎勵 ×5.0），觸發全服 ×6.0 加成 15 秒
	{ID: "T149", Name: "幸運創世魚", Type: TypeSpecial, Multiplier: 550, HP: 580, SpawnWeight: 1, Speed: 45, Lifetime: 24, LaborGain: 110, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T150 幸運重生魚：擊破後 15 秒重生之力（死亡目標復活 HP 50%，擊破獎勵 ×3.0），≥8 → 完美重生全服 ×6.5
	{ID: "T150", Name: "幸運重生魚", Type: TypeSpecial, Multiplier: 600, HP: 640, SpawnWeight: 1, Speed: 40, Lifetime: 26, LaborGain: 120, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-308 新增特殊目標（650x-750x）────────────────────────
	// T151 幸運覺醒鱷魚：擊破後覺醒鱷魚自動獵魚 20 秒（每次獵魚 ×3.0），獵魚 ≥8 → 完美覺醒全服 ×3.5 加成 9 秒
	{ID: "T151", Name: "幸運覺醒鱷魚", Type: TypeSpecial, Multiplier: 650, HP: 680, SpawnWeight: 1, Speed: 45, Lifetime: 26, LaborGain: 130, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T152 幸運吸血鬼升級魚：擊破後 25 秒吸血模式（每次擊破 +1.5x，最高 ×10.0），吸收 ≥10 → 完美吸血全服 ×4.0 加成 10 秒
	{ID: "T152", Name: "幸運吸血鬼升級魚", Type: TypeSpecial, Multiplier: 680, HP: 720, SpawnWeight: 1, Speed: 50, Lifetime: 24, LaborGain: 136, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T153 幸運超級覺醒魚：擊破後全場 HP 歸零（每個獎勵 ×4.0），觸發全服 ×7.0 加成 15 秒
	{ID: "T153", Name: "幸運超級覺醒魚", Type: TypeSpecial, Multiplier: 700, HP: 750, SpawnWeight: 1, Speed: 40, Lifetime: 28, LaborGain: 140, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T154 幸運巨型獎勵魚：擊破後 5 次隨機大獎（×5.0-×50.0），平均 ≥20x → 完美大獎全服 ×4.5 加成 10 秒
	{ID: "T154", Name: "幸運巨型獎勵魚", Type: TypeSpecial, Multiplier: 720, HP: 780, SpawnWeight: 1, Speed: 55, Lifetime: 22, LaborGain: 144, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T155 幸運不死 BOSS 魚：擊破後召喚不死 BOSS（5 條命，每次擊破倍率 +0.5x），18 秒內耗盡 5 條命 → 完美不死全服 ×5.0 加成 12 秒
	{ID: "T155", Name: "幸運不死 BOSS 魚", Type: TypeSpecial, Multiplier: 750, HP: 820, SpawnWeight: 1, Speed: 35, Lifetime: 30, LaborGain: 150, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-309 新增特殊目標（800x-1000x）────────────────────────
	// T156 幸運冰鳳凰魚：擊破後冰凍全場 10 秒（傷害 ×1.5），鳳凰重生爆炸（HP -60%），命中 ≥8 → 完美鳳凰全服 ×5.5 加成 12 秒
	{ID: "T156", Name: "幸運冰鳳凰魚", Type: TypeSpecial, Multiplier: 800, HP: 860, SpawnWeight: 1, Speed: 40, Lifetime: 28, LaborGain: 160, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T157 幸運龍怒能量魚：擊破後能量累積 15 秒（每次擊破 +10 能量），滿 100 → 龍怒全場（HP -80%），命中 ≥10 → 完美龍怒全服 ×6.0 加成 13 秒
	{ID: "T157", Name: "幸運龍怒能量魚", Type: TypeSpecial, Multiplier: 850, HP: 900, SpawnWeight: 1, Speed: 45, Lifetime: 26, LaborGain: 170, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T158 幸運倍率瀑布魚：擊破後 30 秒倍率瀑布（每次擊破 +0.5x，最高 ×20.0），30 秒內達到 ×15.0 → 完美瀑布全服 ×6.5 加成 14 秒
	{ID: "T158", Name: "幸運倍率瀑布魚", Type: TypeSpecial, Multiplier: 900, HP: 940, SpawnWeight: 1, Speed: 50, Lifetime: 24, LaborGain: 180, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T159 幸運覺醒 BOSS 魚 v2：擊破後 8 次 Power Up（每次 8x-15x 隨機），全部命中 → 完美覺醒全服 ×7.0 加成 15 秒
	{ID: "T159", Name: "幸運覺醒BOSS魚v2", Type: TypeSpecial, Multiplier: 950, HP: 980, SpawnWeight: 1, Speed: 40, Lifetime: 30, LaborGain: 190, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T160 幸運終極審判魚：擊破後全場目標 HP 歸零（每個獎勵 ×6.0），觸發全服 ×10.0 加成 20 秒（最高倍率機制）
	{ID: "T160", Name: "幸運終極審判魚", Type: TypeSpecial, Multiplier: 1000, HP: 1050, SpawnWeight: 1, Speed: 35, Lifetime: 32, LaborGain: 200, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-310 新增特殊目標（1050x-1250x）────────────────────────
	// T161 幸運連擊爆發魚：擊破後 20 秒連擊模式，每次擊破 Combo +1（倍率 +0.5x，最高 ×15.0），Combo ≥10 → 完美連擊全服 ×5.5 加成 12 秒
	{ID: "T161", Name: "幸運連擊爆發魚", Type: TypeSpecial, Multiplier: 1050, HP: 1100, SpawnWeight: 1, Speed: 45, Lifetime: 30, LaborGain: 210, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T162 幸運時間炸彈魚：擊破後 30 秒倒數炸彈，每次擊破 +1 能量（最高 30），爆炸時每點能量 HP -2%，能量 ≥20 → 完美爆炸全服 ×6.0 加成 13 秒
	{ID: "T162", Name: "幸運時間炸彈魚", Type: TypeSpecial, Multiplier: 1100, HP: 1150, SpawnWeight: 1, Speed: 40, Lifetime: 32, LaborGain: 220, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T163 幸運元素融合魚：擊破後 25 秒三元素融合（火/冰/雷），每元素全場 HP -25%，三元素全觸發 → 元素爆發全服 ×6.5 加成 14 秒
	{ID: "T163", Name: "幸運元素融合魚", Type: TypeSpecial, Multiplier: 1150, HP: 1200, SpawnWeight: 1, Speed: 35, Lifetime: 34, LaborGain: 230, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T164 幸運寶藏獵人魚：擊破後標記 5 個隨機寶藏（×10-×100），30 秒內全部擊破 → 完美寶藏全服 ×7.0 加成 15 秒
	{ID: "T164", Name: "幸運寶藏獵人魚", Type: TypeSpecial, Multiplier: 1200, HP: 1250, SpawnWeight: 1, Speed: 50, Lifetime: 30, LaborGain: 240, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T165 幸運神話覺醒魚：擊破後全場目標倍率 ×3.0，持續 25 秒，25 秒內擊破 ≥15 個 → 神話完美全服 ×8.0 加成 20 秒（最長加成）
	{ID: "T165", Name: "幸運神話覺醒魚", Type: TypeSpecial, Multiplier: 1250, HP: 1300, SpawnWeight: 1, Speed: 30, Lifetime: 36, LaborGain: 250, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-312 新增特殊目標（1300x-1500x）────────────────────────
	// T166 幸運星際門戶魚：擊破後開啟星際門戶，傳送 5 個目標到中央（HP -50%），全部傳送 → 完美門戶全服 ×5.5 加成 12 秒
	{ID: "T166", Name: "幸運星際門戶魚", Type: TypeSpecial, Multiplier: 1300, HP: 1360, SpawnWeight: 1, Speed: 40, Lifetime: 34, LaborGain: 260, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T167 幸運龍魂融合魚：擊破後龍魂融合 30 秒（每次擊破 +1 魂，最高 50 魂），50 魂 → 龍魂爆發全場 HP -90%，全服 ×9.0 加成 18 秒
	{ID: "T167", Name: "幸運龍魂融合魚", Type: TypeSpecial, Multiplier: 1350, HP: 1400, SpawnWeight: 1, Speed: 35, Lifetime: 36, LaborGain: 270, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T168 幸運時空裂縫魚：擊破後時空裂縫 20 秒（每 4 秒瞬間擊破 3 個目標，獎勵 ×4.0），裂縫期間擊破 ≥12 → 時空完美全服 ×7.5 加成 16 秒
	{ID: "T168", Name: "幸運時空裂縫魚", Type: TypeSpecial, Multiplier: 1400, HP: 1450, SpawnWeight: 1, Speed: 45, Lifetime: 32, LaborGain: 280, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T169 幸運神聖審判魚：擊破後神聖審判 25 秒（每 5 秒一波神聖光柱 HP -30%），5 波全部命中 ≥5 → 神聖完美全服 ×8.5 加成 18 秒
	{ID: "T169", Name: "幸運神聖審判魚", Type: TypeSpecial, Multiplier: 1450, HP: 1500, SpawnWeight: 1, Speed: 30, Lifetime: 38, LaborGain: 290, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T170 幸運宇宙大爆炸魚：擊破後全場 HP 歸零（每個獎勵 ×8.0），觸發全服 ×12.0 加成 25 秒（遊戲最高倍率機制）
	{ID: "T170", Name: "幸運宇宙大爆炸魚", Type: TypeSpecial, Multiplier: 1500, HP: 1560, SpawnWeight: 1, Speed: 25, Lifetime: 40, LaborGain: 300, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-313 新增特殊目標（Progressive Jackpot 系列）────────────────────────
	// T171 幸運 Mini Jackpot 魚：擊破後直接觸發 Mini Jackpot（50x 起跳累積獎池）
	{ID: "T171", Name: "幸運Mini Jackpot魚", Type: TypeSpecial, Multiplier: 50, HP: 60, SpawnWeight: 4, Speed: 85, Lifetime: 12, LaborGain: 10, Behavior: BehaviorFlee, DiffFactor: 0.85},
	// T172 幸運 Minor Jackpot 魚：擊破後直接觸發 Minor Jackpot（200x 起跳累積獎池）
	{ID: "T172", Name: "幸運Minor Jackpot魚", Type: TypeSpecial, Multiplier: 200, HP: 220, SpawnWeight: 2, Speed: 75, Lifetime: 14, LaborGain: 40, Behavior: BehaviorLinear, DiffFactor: 0.95},
	// T173 幸運 Major Jackpot 魚：擊破後直接觸發 Major Jackpot（1000x 起跳累積獎池）
	{ID: "T173", Name: "幸運Major Jackpot魚", Type: TypeSpecial, Multiplier: 1000, HP: 1050, SpawnWeight: 1, Speed: 60, Lifetime: 18, LaborGain: 200, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T174 幸運 Grand Jackpot 魚：擊破後直接觸發 Grand Jackpot（5000x 起跳累積獎池）
	{ID: "T174", Name: "幸運Grand Jackpot魚", Type: TypeSpecial, Multiplier: 5000, HP: 5200, SpawnWeight: 1, Speed: 40, Lifetime: 30, LaborGain: 1000, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T175 幸運 Jackpot Trigger 魚：擊破後隨機觸發四層之一（Mini 60%/Minor 30%/Major 8%/Grand 2%）
	{ID: "T175", Name: "幸運Jackpot觸發魚", Type: TypeSpecial, Multiplier: 200, HP: 240, SpawnWeight: 2, Speed: 70, Lifetime: 16, LaborGain: 40, Behavior: BehaviorLinear, DiffFactor: 1.0},

	// ── DAY-314 新增特殊目標（1550x-1800x）────────────────────────
	// T176 幸運多重宇宙魚：擊破後開啟 3 個平行宇宙，每個宇宙擊破 5 個目標，全部完成 → 全服 ×13.0 加成 28 秒
	{ID: "T176", Name: "幸運多重宇宙魚", Type: TypeSpecial, Multiplier: 1550, HP: 1620, SpawnWeight: 1, Speed: 30, Lifetime: 40, LaborGain: 310, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T177 幸運時間迴圈魚：擊破後 3 次時間迴圈（每次 15 秒，獎勵 ×1.5 遞增），全部完成 → 全服 ×10.0 加成 22 秒
	{ID: "T177", Name: "幸運時間迴圈魚", Type: TypeSpecial, Multiplier: 1600, HP: 1680, SpawnWeight: 1, Speed: 35, Lifetime: 50, LaborGain: 320, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T178 幸運命運之輪魚：擊破後觸發命運之輪（3 次旋轉，最高 ×50.0），連續 3 次 ≥20x → 全服 ×11.0 加成 24 秒
	{ID: "T178", Name: "幸運命運之輪魚", Type: TypeSpecial, Multiplier: 1650, HP: 1720, SpawnWeight: 1, Speed: 40, Lifetime: 36, LaborGain: 330, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T179 幸運神域降臨魚：擊破後神域降臨 30 秒（每 6 秒神域波 HP -35%），5 波全部命中 ≥6 → 全服 ×14.0 加成 30 秒
	{ID: "T179", Name: "幸運神域降臨魚", Type: TypeSpecial, Multiplier: 1700, HP: 1780, SpawnWeight: 1, Speed: 25, Lifetime: 42, LaborGain: 340, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T180 幸運終焉之力魚：擊破後全場 HP 歸零（每個獎勵 ×10.0），觸發全服 ×15.0 加成 30 秒（超越 T170 成為新最高倍率機制）
	{ID: "T180", Name: "幸運終焉之力魚", Type: TypeSpecial, Multiplier: 1800, HP: 1900, SpawnWeight: 1, Speed: 20, Lifetime: 45, LaborGain: 360, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-315 新增特殊目標（1850x-2100x）────────────────────────
	// T181 幸運突變魚：擊破後觸發隨機突變（150種突變，最高 ×17.0 加成），突變 ≥10x → 全服 ×16.0 加成 32 秒
	// 業界依據：Fisch mutations system（150+ mutations, 17x bonus）
	{ID: "T181", Name: "幸運突變魚", Type: TypeSpecial, Multiplier: 1850, HP: 1950, SpawnWeight: 1, Speed: 35, Lifetime: 42, LaborGain: 370, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T182 幸運北極風暴魚：擊破後快速節奏 500x 連擊（每 0.3 秒一波，共 8 波），全部命中 → 全服 ×16.5 加成 33 秒
	// 業界依據：Arctic Mechanics（500x multiplier, fast pace）
	{ID: "T182", Name: "幸運北極風暴魚", Type: TypeSpecial, Multiplier: 1900, HP: 2000, SpawnWeight: 1, Speed: 45, Lifetime: 38, LaborGain: 380, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T183 幸運漁夫野生魚：擊破後標記 3 個 Wild 目標（HP -50%，擊破獎勵 ×5.0），全部擊破 → 全服 ×17.0 加成 35 秒
	// 業界依據：Big Bass Splash 1000（Fisherman Wild + Fish Cash mechanic）
	{ID: "T183", Name: "幸運漁夫野生魚", Type: TypeSpecial, Multiplier: 1950, HP: 2050, SpawnWeight: 1, Speed: 30, Lifetime: 44, LaborGain: 390, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T184 幸運風險等級魚：擊破後選擇 5 個風險等級（低風險 ×5.0 / 中 ×20.0 / 高 ×100.0 / 極高 ×500.0 / 最高 ×3000.0），最高等級 → 全服 ×17.5 加成 36 秒
	// 業界依據：BGaming Fishing Club 2（5 risk levels, max x3000）
	{ID: "T184", Name: "幸運風險等級魚", Type: TypeSpecial, Multiplier: 2000, HP: 2100, SpawnWeight: 1, Speed: 25, Lifetime: 46, LaborGain: 400, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T185 幸運宇宙脈衝魚：擊破後宇宙脈衝波（全場 HP -45%，每個獎勵 ×12.0），觸發全服 ×16.0 加成 35 秒（新最高倍率機制）
	// 業界依據：Fishing Fortune multiplier cascade（2x→500x）升級版
	{ID: "T185", Name: "幸運宇宙脈衝魚", Type: TypeSpecial, Multiplier: 2100, HP: 2200, SpawnWeight: 1, Speed: 20, Lifetime: 48, LaborGain: 420, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-316 新增特殊目標（2150x-2500x）────────────────────────
	// T186 幸運鏡像宇宙魚：擊破後開啟鏡像宇宙 25 秒，複製場上最強 3 個目標（HP 50%，獎勵 ×2.0）
	// 業界依據：Royal Fishing「Mirror Fish」+ 量子糾纏概念
	{ID: "T186", Name: "幸運鏡像宇宙魚", Type: TypeSpecial, Multiplier: 2150, HP: 2260, SpawnWeight: 1, Speed: 30, Lifetime: 48, LaborGain: 430, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T187 幸運引力場魚：擊破後引力場 15 秒（目標速度 ×0.1），15 秒後引力爆炸（HP -55%，獎勵 ×9.0），命中 ≥12 → 全服 ×17.5 加成 37 秒
	// 業界依據：Black Hole Fishing（SDG Games, 2026）引力吸引機制
	{ID: "T187", Name: "幸運引力場魚", Type: TypeSpecial, Multiplier: 2200, HP: 2320, SpawnWeight: 1, Speed: 25, Lifetime: 50, LaborGain: 440, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T188 幸運時間加速魚：擊破後時間加速 30 秒（目標速度 ×0.15，射擊速度 ×3.0，獎勵 ×2.5），擊破 ≥20 → 全服 ×18.0 加成 38 秒（新最高）
	// 業界依據：Fishing Fortune「time warp」升級版 + 時間操控概念
	{ID: "T188", Name: "幸運時間加速魚", Type: TypeSpecial, Multiplier: 2300, HP: 2420, SpawnWeight: 1, Speed: 35, Lifetime: 46, LaborGain: 460, Behavior: BehaviorLinear, DiffFactor: 1.0},
	// T189 幸運星雲漩渦魚：擊破後星雲漩渦 20 秒（每秒全場 HP -8%，獎勵 ×1.5），累積命中 ≥20 → 全服 ×18.5 加成 39 秒
	// 業界依據：Fishing Carnival「vortex anemone」+ 星雲能量吸收概念
	{ID: "T189", Name: "幸運星雲漩渦魚", Type: TypeSpecial, Multiplier: 2400, HP: 2520, SpawnWeight: 1, Speed: 20, Lifetime: 52, LaborGain: 480, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T190 幸運宇宙審判魚：擊破後宇宙審判（全場 HP 歸零，每個獎勵 ×14.0），觸發全服 ×19.0 加成 40 秒（新最高全服倍率機制）
	// 業界依據：Fishing Fortune「ultimate judgment」+ 宇宙終極機制
	{ID: "T190", Name: "幸運宇宙審判魚", Type: TypeSpecial, Multiplier: 2500, HP: 2640, SpawnWeight: 1, Speed: 15, Lifetime: 55, LaborGain: 500, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-317 新增特殊目標（2600x-3000x）────────────────────────
	// T191 幸運 PvP 競技魚：全服 PvP 競技 30 秒，擊破最多目標者獲得 ×20.0 加成 35 秒，全服 ×19.5 加成 40 秒
	// 業界依據：PvP 競技機制 + 全服競爭設計
	{ID: "T191", Name: "幸運PvP競技魚", Type: TypeSpecial, Multiplier: 2600, HP: 2750, SpawnWeight: 1, Speed: 12, Lifetime: 58, LaborGain: 520, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T192 幸運技能連鎖魚：技能連鎖 25 秒，每次擊破提升技能等級（Lv.1-10），Lv.10 → 全服 ×20.0 加成 38 秒
	// 業界依據：技能連鎖升級機制 + 等級倍率設計
	{ID: "T192", Name: "幸運技能連鎖魚", Type: TypeSpecial, Multiplier: 2700, HP: 2860, SpawnWeight: 1, Speed: 10, Lifetime: 60, LaborGain: 540, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T193 幸運全服大爆炸魚：全服大爆炸，全場 HP 歸零，每個獎勵 ×15.0，觸發全服 ×20.5 加成 40 秒（新最高）
	// 業界依據：全場清空機制升級版 + 超高全服倍率
	{ID: "T193", Name: "幸運全服大爆炸魚", Type: TypeSpecial, Multiplier: 2800, HP: 2970, SpawnWeight: 1, Speed: 8, Lifetime: 62, LaborGain: 560, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T194 幸運時空折疊魚：時空折疊 20 秒，所有目標倍率 ×3.0，射擊速度 ×2.0，觸發全服 ×21.0 加成 42 秒（新最高）
	// 業界依據：時空操控機制 + 全場倍率加成設計
	{ID: "T194", Name: "幸運時空折疊魚", Type: TypeSpecial, Multiplier: 2900, HP: 3080, SpawnWeight: 1, Speed: 6, Lifetime: 65, LaborGain: 580, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T195 幸運宇宙終焉魚：宇宙終焉，全場 HP 歸零，每個獎勵 ×20.0，觸發全服 ×22.0 加成 45 秒（史上最高）
	// 業界依據：終極清場機制 + 史上最高全服倍率（超越 T190 的 ×19.0）
	{ID: "T195", Name: "幸運宇宙終焉魚", Type: TypeSpecial, Multiplier: 3000, HP: 3200, SpawnWeight: 1, Speed: 5, Lifetime: 68, LaborGain: 600, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-318 新增特殊目標（3100x-5000x）────────────────────────
	// T196 幸運龍王輪盤魚：雙環輪盤，內環 × 外環 = 最高 ×25.0，觸發全服 ×23.0 加成 46 秒
	// 業界依據：Royal Fishing「ChainLong King Wheel」雙環輪盤機制升級版
	{ID: "T196", Name: "幸運龍王輪盤魚", Type: TypeSpecial, Multiplier: 3100, HP: 3300, SpawnWeight: 1, Speed: 4, Lifetime: 70, LaborGain: 620, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T197 幸運永恆循環魚：永恆循環 10 波，每波獎勵遞增（×1.0 → ×10.0），全服 ×23.5 加成 47 秒
	// 業界依據：Fishing Fortune「time loop」升級版 + 永恆循環概念
	{ID: "T197", Name: "幸運永恆循環魚", Type: TypeSpecial, Multiplier: 3200, HP: 3420, SpawnWeight: 1, Speed: 3, Lifetime: 72, LaborGain: 640, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T198 幸運混沌爆炸魚：混沌爆炸，隨機 3-8 個目標同時爆炸，倍率疊加最高 ×30.0，全服 ×24.0 加成 48 秒
	// 業界依據：Fishing Carnival「Big Bang」升級版 + 混沌隨機機制
	{ID: "T198", Name: "幸運混沌爆炸魚", Type: TypeSpecial, Multiplier: 3300, HP: 3540, SpawnWeight: 1, Speed: 3, Lifetime: 74, LaborGain: 660, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T199 幸運神聖復活魚：最近死亡的 5 個目標全部復活（HP 80%，獎勵 ×4.0），全服 ×24.5 加成 49 秒
	// 業界依據：Phoenix rebirth 機制升級版 + 神聖復活概念
	{ID: "T199", Name: "幸運神聖復活魚", Type: TypeSpecial, Multiplier: 3400, HP: 3660, SpawnWeight: 1, Speed: 2, Lifetime: 76, LaborGain: 680, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T200 幸運創世紀元魚：里程碑第 200 個 Lucky 目標，全場 HP 歸零，每個獎勵 ×25.0，全服 ×25.0 加成 50 秒（史上最高）
	// 業界依據：終極里程碑機制 + 史上最高全服倍率（超越 T195 的 ×22.0）
	{ID: "T200", Name: "幸運創世紀元魚", Type: TypeSpecial, Multiplier: 5000, HP: 5000, SpawnWeight: 1, Speed: 1, Lifetime: 80, LaborGain: 1000, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-319 新增特殊目標（5200x-8888x）────────────────────────
	// T201 幸運能量風暴魚：5 波連鎖電擊（每波全場 HP -30%），全服 ×26.0 加成 52 秒
	// 業界依據：Royal Fishing「60x lightning eel chain reaction」升級版 + 2026 能量風暴機制
	{ID: "T201", Name: "幸運能量風暴魚", Type: TypeSpecial, Multiplier: 5200, HP: 5200, SpawnWeight: 1, Speed: 1, Lifetime: 82, LaborGain: 1040, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T202 幸運水晶共鳴魚：全場共鳴爆炸（每個獎勵 ×30.0），全服 ×27.0 加成 54 秒
	// 業界依據：Fishing Legend 2025「Crystal Resonance」全場共鳴爆炸機制
	{ID: "T202", Name: "幸運水晶共鳴魚", Type: TypeSpecial, Multiplier: 5500, HP: 5500, SpawnWeight: 1, Speed: 1, Lifetime: 84, LaborGain: 1100, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T203 幸運命運審判魚：隨機 5 個目標各 ×50-×500，全服 ×28.0 加成 56 秒
	// 業界依據：Fishing Fortune「Fate Judgment」升級版 + 2026 命運審判機制
	{ID: "T203", Name: "幸運命運審判魚", Type: TypeSpecial, Multiplier: 5800, HP: 5800, SpawnWeight: 1, Speed: 1, Lifetime: 86, LaborGain: 1160, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T204 幸運時間逆流魚：最近死亡的 10 個目標全部復活（HP 100%，獎勵 ×5.0），全服 ×29.0 加成 58 秒
	// 業界依據：T199 神聖復活魚升級版 + 2026 時間逆流機制（最多 10 個目標）
	{ID: "T204", Name: "幸運時間逆流魚", Type: TypeSpecial, Multiplier: 6000, HP: 6000, SpawnWeight: 1, Speed: 1, Lifetime: 88, LaborGain: 1200, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T205 幸運宇宙奇點魚：全場 HP 歸零（每個獎勵 ×30.0），全服 ×30.0 加成 60 秒（史上最高，超越 T200 的 ×25.0）
	// 業界依據：終極宇宙奇點機制 + 2026 最高倍率設計（8888x 吉祥數字）
	{ID: "T205", Name: "幸運宇宙奇點魚", Type: TypeSpecial, Multiplier: 8888, HP: 8888, SpawnWeight: 1, Speed: 1, Lifetime: 90, LaborGain: 1777, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-323 新增特殊目標（9500x-16888x）────────────────────────
	// T206 幸運 Fever Boost 魚：Fever Boost™ 機制（Games Global 2026-05-28），30 秒內所有特效機率翻倍，全場倍率 ×2.0，全服 ×31.0 加成 62 秒
	// 業界依據：Games Global「Fishin' Pots of Gold」Fever Boost™（2026-05-28 最新）
	{ID: "T206", Name: "幸運FeverBoost魚", Type: TypeSpecial, Multiplier: 9500, HP: 9500, SpawnWeight: 1, Speed: 1, Lifetime: 92, LaborGain: 1900, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T207 幸運公會戰魚：全服公會戰 45 秒，擊破最多目標的玩家獲得 ×35.0，全服 ×32.0 加成 64 秒
	// 業界依據：Fishing Frenzy Chapter 3「Guild Wars + Boss Fish」（2026-05-27）
	{ID: "T207", Name: "幸運公會戰魚", Type: TypeSpecial, Multiplier: 10000, HP: 10000, SpawnWeight: 1, Speed: 1, Lifetime: 94, LaborGain: 2000, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T208 幸運路徑魚：Fish Road 機制，路徑越遠倍率越高（最高 20,000x），全服 ×33.0 加成 66 秒
	// 業界依據：Fish Road「路徑越遠倍率越高，最高 20,000x」（fishroad.eu）
	{ID: "T208", Name: "幸運路徑魚", Type: TypeSpecial, Multiplier: 11000, HP: 11000, SpawnWeight: 1, Speed: 1, Lifetime: 96, LaborGain: 2200, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T209 幸運連鎖電鰻魚：Royal Fishing 紫粉色電鰻升級版，連鎖電擊 8 條魚，每條 ×40.0，全服 ×34.0 加成 68 秒
	// 業界依據：Royal Fishing「Purple/Pink Lightning Eel chain reaction」升級版（royal-fishing.co.uk）
	{ID: "T209", Name: "幸運連鎖電鰻魚", Type: TypeSpecial, Multiplier: 12000, HP: 12000, SpawnWeight: 1, Speed: 1, Lifetime: 98, LaborGain: 2400, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T210 幸運終極奇蹟魚：終極機制，全場 HP 歸零（每個獎勵 ×50.0），全服 ×35.0 加成 70 秒（新史上最高）
	// 業界依據：終極奇蹟機制 + 2026 最高倍率設計（16888x 吉祥數字，超越 T205 的 8888x）
	{ID: "T210", Name: "幸運終極奇蹟魚", Type: TypeSpecial, Multiplier: 16888, HP: 16888, SpawnWeight: 1, Speed: 1, Lifetime: 100, LaborGain: 3377, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-324 新增特殊目標（17000x-22000x）────────────────────────
	// T211 幸運雪崩魚：Avalanche Cascade 機制，8 波連鎖消除，每波倍率 +5.0，全部命中 → 全服 ×36.0 加成 72 秒
	// 業界依據：Relax Gaming「Cod of Thunder Dream Drop」Multiplier Ladder（2026）
	{ID: "T211", Name: "幸運雪崩魚", Type: TypeSpecial, Multiplier: 17000, HP: 17000, SpawnWeight: 1, Speed: 1, Lifetime: 102, LaborGain: 3400, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T212 幸運崩潰倍率魚：Crash Multiplier 機制，倍率持續上升直到崩潰，完美收割（≥40.0）→ 全服 ×36.5 加成 73 秒
	// 業界依據：cardsrealm.com「Hybrid Crash Game」趨勢（2026-05）
	{ID: "T212", Name: "幸運崩潰倍率魚", Type: TypeSpecial, Multiplier: 18000, HP: 18000, SpawnWeight: 1, Speed: 1, Lifetime: 104, LaborGain: 3600, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T213 幸運倍率梯魚：Multiplier Ladder 機制，每次擊破提升梯度（Lv.1-10），Lv.10 → 全服 ×37.0 加成 74 秒
	// 業界依據：Relax Gaming「Cod of Thunder Dream Drop」Multiplier Ladder（2026）
	{ID: "T213", Name: "幸運倍率梯魚", Type: TypeSpecial, Multiplier: 19000, HP: 19000, SpawnWeight: 1, Speed: 1, Lifetime: 106, LaborGain: 3800, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T214 幸運冰釣輪盤魚：Ice Fishing Wheel 機制，3 次旋轉（最高 ×5000），最高單次 ≥2000 → 全服 ×37.5 加成 75 秒
	// 業界依據：Evolution Gaming「Ice Fishing」最高 5000x（2026）
	{ID: "T214", Name: "幸運冰釣輪盤魚", Type: TypeSpecial, Multiplier: 20000, HP: 20000, SpawnWeight: 1, Speed: 1, Lifetime: 108, LaborGain: 4000, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T215 幸運全服雪崩魚：Global Avalanche 機制，5 波全服連鎖消除，每波 ×8.0，全服 ×38.0 加成 76 秒（新史上最高）
	// 業界依據：Avalanche Reels + Global Multiplier 組合（2026 最新趨勢）
	{ID: "T215", Name: "幸運全服雪崩魚", Type: TypeSpecial, Multiplier: 22000, HP: 22000, SpawnWeight: 1, Speed: 1, Lifetime: 110, LaborGain: 4400, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-325 新增特殊目標（23000x-28888x）────────────────────────
	// T216 幸運漁網魚：Fishing Net 機制，撒網捕獲全場所有目標，每個獎勵 ×60.0，全服 ×38.5 加成 77 秒
	// 業界依據：BGaming「Fishing Club 2」Fishing Net Bonus（×60 stake，2026-04）
	{ID: "T216", Name: "幸運漁網魚", Type: TypeSpecial, Multiplier: 23000, HP: 23000, SpawnWeight: 1, Speed: 1, Lifetime: 112, LaborGain: 4600, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T217 幸運 TNT 爆炸魚：TNT Bonus 機制，水下大爆炸（全場 HP -80%，每個 ×100.0），全服 ×39.0 加成 78 秒
	// 業界依據：BGaming「Fishing Club 2」TNT Bonus（×100 stake，2026-04）
	{ID: "T217", Name: "幸運TNT爆炸魚", Type: TypeSpecial, Multiplier: 24000, HP: 24000, SpawnWeight: 1, Speed: 1, Lifetime: 114, LaborGain: 4800, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T218 幸運擾動魚：Disturbance System，活躍度越高倍率越高（最高 ×50.0），全服 ×39.5 加成 79 秒
	// 業界依據：Fisch「Disturbance System」活躍度驅動稀有魚生成（2026-01）
	{ID: "T218", Name: "幸運擾動魚", Type: TypeSpecial, Multiplier: 25000, HP: 25000, SpawnWeight: 1, Speed: 1, Lifetime: 116, LaborGain: 5000, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T219 幸運珍珠倍率魚：Pearl Multiplier，場上每個目標都有珍珠倍率（×1-×100），全服 ×40.0 加成 80 秒（新里程碑）
	// 業界依據：BGaming「Shark & Spark Hold & Win」Pearl 倍率符號（2026-05）
	{ID: "T219", Name: "幸運珍珠倍率魚", Type: TypeSpecial, Multiplier: 26000, HP: 26000, SpawnWeight: 1, Speed: 1, Lifetime: 118, LaborGain: 5200, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T220 幸運快速暴富魚：Rapid Riches 機制，5 秒內快速連擊（每次 ×200.0），全服 ×41.0 加成 82 秒（新史上最高）
	// 業界依據：Reflex Gaming「Big Game Fishing Rapid Riches」快速獎勵機制（2026-05）
	{ID: "T220", Name: "幸運快速暴富魚", Type: TypeSpecial, Multiplier: 28888, HP: 28888, SpawnWeight: 1, Speed: 1, Lifetime: 120, LaborGain: 5776, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-326 新增特殊目標（30000x-35000x）────────────────────────
	// T221 幸運骰子獎勵魚：Dice Bonus 機制，擲骰3次（1-3點×50/4-5點×150/6點×300），全服×41.5加成83秒
	// 業界依據：BGaming「Shark & Spark Hold & Win」Dice Bonus（2026-05-25）
	{ID: "T221", Name: "幸運骰子獎勵魚", Type: TypeSpecial, Multiplier: 30000, HP: 30000, SpawnWeight: 1, Speed: 1, Lifetime: 122, LaborGain: 6000, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T222 幸運雙Bonus魚：Dual Bonus 機制，選擇金幣收集（5個×80，全服×41.8加成84秒）或風險輪盤（最高×500，全服×42.0加成85秒）
	// 業界依據：BGaming「Fishing Club 2」雙 Bonus 遊戲（2026-04）
	{ID: "T222", Name: "幸運雙Bonus魚", Type: TypeSpecial, Multiplier: 32000, HP: 32000, SpawnWeight: 1, Speed: 1, Lifetime: 124, LaborGain: 6400, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T223 幸運Coin Respin魚：Hold & Win 風格，9格盤面，Bronze×10/Silver×30/Gold×80/Diamond×200，填滿+×500
	// 全服×42.5加成86秒（新史上最高，超越T222的×42.0）
	// 業界依據：BGaming「Shark & Spark Hold & Win」Coin Respin（2026-05-28）
	{ID: "T223", Name: "幸運CoinRespin魚", Type: TypeSpecial, Multiplier: 35000, HP: 35000, SpawnWeight: 1, Speed: 1, Lifetime: 126, LaborGain: 7000, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-327 新增特殊目標（38000x-50000x）────────────────────────
	// T224 幸運黃金鍋魚：Gold Blitz™ Cash Collection，12格黃金鍋，Copper×5/Silver×20/Gold×60/Platinum×150/Diamond×200
	// 填滿鍋子 → 額外 ×300，全服 ×43.0 加成 88 秒（新史上最高，超越 T223 的 ×42.5）
	// 業界依據：Games Global「Fishin' Pots of Gold Gold Blitz Ultimate」（2026-05-28）
	{ID: "T224", Name: "幸運黃金鍋魚", Type: TypeSpecial, Multiplier: 38000, HP: 38000, SpawnWeight: 1, Speed: 1, Lifetime: 128, LaborGain: 7600, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T225 幸運瀑布鎖定魚：Cascading Wins + 鎖定倍率，8 波連鎖，Pearl 符號 ×2-×10 加成
	// 完美 8 波 → 額外 ×50，全服 ×43.5 加成 87 秒（超越 T224 的 ×43.0）
	// 業界依據：BGaming「Shark & Spark Hold & Win」Cascading Wins + Pearl Multipliers（2026-05-28）
	{ID: "T225", Name: "幸運瀑布鎖定魚", Type: TypeSpecial, Multiplier: 40000, HP: 40000, SpawnWeight: 1, Speed: 1, Lifetime: 130, LaborGain: 8000, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T226 幸運傳說覺醒魚：Legend Dragon 覺醒升級，8 次連續獎勵（Humpback 90-150x / Legend Dragon 120-200x）
	// 全部 8 次完成 → 全服 ×44.0 加成 88 秒（超越 T225 的 ×43.5）
	// 業界依據：Royal Fishing Jili「Legend Dragon 120-200x, Humpback Whale 90-150x」（2026）
	{ID: "T226", Name: "幸運傳說覺醒魚", Type: TypeSpecial, Multiplier: 43000, HP: 43000, SpawnWeight: 1, Speed: 1, Lifetime: 132, LaborGain: 8600, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T227 幸運崩潰收割魚：Crash Harvest 機制，倍率持續上升，玩家選擇收割時機
	// 完美收割（≥50x 未崩潰）→ 全服 ×44.5 加成 89 秒（超越 T226 的 ×44.0）
	// 業界依據：Lucky Fish AbraCadabra「Crash mechanic」（2026-05）
	{ID: "T227", Name: "幸運崩潰收割魚", Type: TypeSpecial, Multiplier: 46000, HP: 46000, SpawnWeight: 1, Speed: 1, Lifetime: 134, LaborGain: 9200, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T228 幸運宇宙大融合魚：終極融合機制，4 Phase（Coin Respin + Cascade Lock + Legend Awaken + 全場清空）
	// 全部完成 → 全服 ×45.0 加成 90 秒（新史上最高，超越 T227 的 ×44.5）
	// 業界依據：終極融合設計，整合 2026 年最新業界機制
	{ID: "T228", Name: "幸運宇宙大融合魚", Type: TypeSpecial, Multiplier: 50000, HP: 50000, SpawnWeight: 1, Speed: 1, Lifetime: 136, LaborGain: 10000, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-328 新增特殊目標（55000x-88888x）────────────────────────
	// T229 幸運磁力吸引魚：Magnetic Attraction 機制，磁力吸引全場目標到中心，每個 ×70.0
	// 命中 ≥10 個 → 完美磁力，全服 ×45.5 加成 91 秒（超越 T228 的 ×45.0）
	// 業界依據：Black Hole Fishing 引力機制升級版 + Magnetic Attraction 概念（2026）
	{ID: "T229", Name: "幸運磁力吸引魚", Type: TypeSpecial, Multiplier: 55000, HP: 55000, SpawnWeight: 1, Speed: 1, Lifetime: 138, LaborGain: 11000, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T230 幸運超級連鎖魚：Super Chain 機制，每次擊破觸發 3 條連鎖（每條 ×80.0），連鎖 ≥5 次 → 超級連鎖爆發
	// 全服 ×46.0 加成 92 秒（超越 T229 的 ×45.5）
	// 業界依據：Royal Fishing 連鎖電擊升級版 + Super Chain 概念（2026）
	{ID: "T230", Name: "幸運超級連鎖魚", Type: TypeSpecial, Multiplier: 60000, HP: 60000, SpawnWeight: 1, Speed: 1, Lifetime: 140, LaborGain: 12000, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T231 幸運神聖光柱魚：Holy Pillar 機制，12 道神聖光柱同時降下（每道 HP -50%），命中 ≥8 道 → 完美神聖
	// 全服 ×46.5 加成 93 秒（超越 T230 的 ×46.0）
	// 業界依據：神聖審判機制升級版 + Holy Pillar 概念（2026）
	{ID: "T231", Name: "幸運神聖光柱魚", Type: TypeSpecial, Multiplier: 65000, HP: 65000, SpawnWeight: 1, Speed: 1, Lifetime: 142, LaborGain: 13000, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T232 幸運時間停止魚：Time Stop 機制，全場凍結 15 秒（傷害 ×5.0），凍結結束全場 HP -70%
	// 凍結期間擊破 ≥15 個 → 完美時間停止，全服 ×47.0 加成 94 秒（超越 T231 的 ×46.5）
	// 業界依據：時間凍結機制終極升級版 + Time Stop 概念（2026）
	{ID: "T232", Name: "幸運時間停止魚", Type: TypeSpecial, Multiplier: 70000, HP: 70000, SpawnWeight: 1, Speed: 1, Lifetime: 144, LaborGain: 14000, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T233 幸運宇宙重啟魚：Cosmic Restart 機制，全場 HP 歸零（每個目標獎勵 ×100.0）
	// 觸發全服 ×47.5 加成 95 秒（新史上最高，超越 T228 的 ×45.0）
	// 業界依據：終極清場機制升級版 + Cosmic Restart 概念（2026）
	{ID: "T233", Name: "幸運宇宙重啟魚", Type: TypeSpecial, Multiplier: 88888, HP: 88888, SpawnWeight: 1, Speed: 1, Lifetime: 146, LaborGain: 17777, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── DAY-329 新增特殊目標（95000x-128888x）────────────────────────
	// T234 幸運Fever Boost升級魚：Fever Boost™ Ultimate，清除普通目標只留高倍率特殊目標（×2.0 傷害加成），持續 20 秒
	// 完美觸發（場上特殊目標 ≥5 個）→ 全服 ×48.0 加成 96 秒（新史上最高，超越 T233 的 ×47.5）
	// 業界依據：Games Global「Fishin' Pots of Gold Gold Blitz Ultimate Fever Boost」（2026-05-28）
	{ID: "T234", Name: "幸運FeverBoost升級魚", Type: TypeSpecial, Multiplier: 95000, HP: 95000, SpawnWeight: 1, Speed: 1, Lifetime: 148, LaborGain: 19000, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T235 幸運快速暴富升級魚：Rapid Riches Ultimate，3 秒極速連擊視窗（每次擊破 ×300.0），連擊 ≥10 次 → 完美暴富
	// 全服 ×48.5 加成 97 秒（超越 T234 的 ×48.0）
	// 業界依據：Reflex Gaming「Big Game Fishing Rapid Riches」升級版（2026-05）
	{ID: "T235", Name: "幸運快速暴富升級魚", Type: TypeSpecial, Multiplier: 100000, HP: 100000, SpawnWeight: 1, Speed: 1, Lifetime: 150, LaborGain: 20000, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T236 幸運冰釣大師魚：Ice Fishing Master，5 次旋轉（每次最高 ×8000），最高單次 ≥3000 → 完美冰釣
	// 全服 ×49.0 加成 98 秒（超越 T235 的 ×48.5）
	// 業界依據：Evolution Gaming「Ice Fishing Live」最高 5000x 升級版（2026）
	{ID: "T236", Name: "幸運冰釣大師魚", Type: TypeSpecial, Multiplier: 108888, HP: 108888, SpawnWeight: 1, Speed: 1, Lifetime: 152, LaborGain: 21777, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T237 幸運宇宙奇蹟魚：Cosmic Miracle，全場 HP 歸零（每個目標 ×120.0）+ 8 道宇宙光柱
	// 命中 ≥12 個 → 完美奇蹟，全服 ×49.5 加成 99 秒（超越 T236 的 ×49.0）
	// 業界依據：終極清場機制 + 神聖光柱機制融合升級版（2026）
	{ID: "T237", Name: "幸運宇宙奇蹟魚", Type: TypeSpecial, Multiplier: 118888, HP: 118888, SpawnWeight: 1, Speed: 1, Lifetime: 154, LaborGain: 23777, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T238 幸運創世終極魚：Genesis Ultimate（里程碑：全服 ×50.0）
	// 全場清空（每個目標 ×150.0）+ 12 道創世光柱，全服 ×50.0 加成 100 秒（史上第一個全服 ×50.0）
	// 業界依據：終極清場機制 + 創世機制融合終極版（2026）
	{ID: "T238", Name: "幸運創世終極魚", Type: TypeSpecial, Multiplier: 128888, HP: 128888, SpawnWeight: 1, Speed: 1, Lifetime: 156, LaborGain: 25777, Behavior: BehaviorSink, DiffFactor: 1.0},

	// DAY-331 新增：T239-T243（業界研究：BGaming Shark & Spark Hold & Win 2026-05-30 最新）
	// T239 幸運鯊魚閃電魚：Shark & Spark 機制（里程碑：全服 ×51.0）
	// 鯊魚閃電 + 珍珠倍率組合，閃電連鎖 6 條（每條 ×80.0），全服 ×51.0 加成 102 秒
	// 業界依據：BGaming Shark & Spark Hold & Win（2026-05-30 最新發布）
	{ID: "T239", Name: "幸運鯊魚閃電魚", Type: TypeSpecial, Multiplier: 138888, HP: 138888, SpawnWeight: 1, Speed: 1, Lifetime: 158, LaborGain: 27777, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T240 幸運冬季冰釣魚：Winter Ice Fishing 機制（全服 ×51.5）
	// 冰下魚群 + 53格輪盤，3 次旋轉（最高 ×500），全服 ×51.5 加成 103 秒
	// 業界依據：BGaming Winter Fishing Club（2026-01）+ Evolution Ice Fishing Live（2025-2026）
	{ID: "T240", Name: "幸運冬季冰釣魚", Type: TypeSpecial, Multiplier: 148888, HP: 148888, SpawnWeight: 1, Speed: 1, Lifetime: 160, LaborGain: 29777, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T241 幸運大西洋狂潮魚：Big Atlantis Frenzy 機制（全服 ×52.0）
	// 亞特蘭提斯爆炸 + 連鎖消除，7 波 Fish 符號（×5-×500），全服 ×52.0 加成 104 秒
	// 業界依據：BGaming Big Atlantis Frenzy（2025-2026）
	{ID: "T241", Name: "幸運大西洋狂潮魚", Type: TypeSpecial, Multiplier: 158888, HP: 158888, SpawnWeight: 1, Speed: 1, Lifetime: 162, LaborGain: 31777, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T242 幸運釣魚時間魚：Fishing Time Wheel 機制（全服 ×52.5）
	// 命運輪盤 + 倍率疊加，5 次旋轉（最高 ×10000），全服 ×52.5 加成 105 秒
	// 業界依據：BGaming Fishing Time（2026-04）
	{ID: "T242", Name: "幸運釣魚時間魚", Type: TypeSpecial, Multiplier: 168888, HP: 168888, SpawnWeight: 1, Speed: 1, Lifetime: 164, LaborGain: 33777, Behavior: BehaviorSink, DiffFactor: 1.0},
	// T243 幸運終極鯊魚魚：Ultimate Shark 機制（里程碑：全服 ×53.0）
	// 終極鯊魚清場（每個 ×180.0）+ 14 次鯊魚咬合，全服 ×53.0 加成 106 秒（新史上最高）
	// 業界依據：Shark & Spark Hold & Win 終極版（2026-05-30）
	{ID: "T243", Name: "幸運終極鯊魚魚", Type: TypeSpecial, Multiplier: 188888, HP: 188888, SpawnWeight: 1, Speed: 1, Lifetime: 166, LaborGain: 37777, Behavior: BehaviorSink, DiffFactor: 1.0},

	// ── BOSS ─────────────────────────────────────────────────
	{ID: "B001", Name: "那個孩子", Type: TypeBoss, MinMult: 100, MaxMult: 500, HP: 3000, SpawnWeight: 0, Speed: 0, Lifetime: 60, LaborGain: 30, Behavior: BehaviorLinear, DiffFactor: 1.5},
}

var targetMap map[string]*TargetDef

func init() {
	targetMap = make(map[string]*TargetDef)
	for i := range Targets {
		targetMap[Targets[i].ID] = &Targets[i]
	}
}

func GetTarget(id string) (*TargetDef, bool) {
	t, ok := targetMap[id]
	return t, ok
}

// 流星倍率權重
var MeteorWeights = []struct {
	Mult   float64
	Weight int
}{
	{20, 50},
	{30, 30},
	{40, 15},
	{50, 5},
}

// 擬態怪物倍率權重
var MimicWeights = []struct {
	Mult   float64
	Weight int
}{
	{15, 60},
	{20, 30},
	{30, 10},
}

// 黃金龍魚輪盤倍率權重（T109）
var GoldenDragonWeights = []struct {
	Mult   float64
	Weight int
}{
	{80, 40},
	{100, 25},
	{150, 15},
	{200, 10},
	{250, 6},
	{300, 3},
	{350, 1},
}

// 千龍王輪盤倍率權重（T116）— 內環 × 外環 = 最高 1000x
var ChainLongKingInnerWeights = []struct {
	Mult   float64
	Weight int
}{
	{2, 40},
	{5, 25},
	{10, 15},
	{15, 10},
	{20, 7},
	{25, 3},
}

var ChainLongKingOuterWeights = []struct {
	Mult   float64
	Weight int
}{
	{5, 40},
	{10, 25},
	{20, 15},
	{30, 10},
	{40, 7},
	{50, 3},
}

// ── Bonus 目標物 ──────────────────────────────────────────────

type BonusTargetDef struct {
	ID          string
	Name        string
	Score       int
	SpawnWeight int
	Special     string // "", "hard", "glow", "gold", "trouble"
}

var BonusTargets = []BonusTargetDef{
	{ID: "BG001", Name: "普通雜草", Score: 1, SpawnWeight: 180, Special: ""},
	{ID: "BG002", Name: "硬雜草", Score: 3, SpawnWeight: 80, Special: "hard"},
	{ID: "BG003", Name: "發光雜草", Score: 8, SpawnWeight: 35, Special: "glow"},
	{ID: "BG004", Name: "金色雜草", Score: 20, SpawnWeight: 10, Special: "gold"},
	{ID: "BG005", Name: "搗亂怪草", Score: -5, SpawnWeight: 20, Special: "trouble"},
}
