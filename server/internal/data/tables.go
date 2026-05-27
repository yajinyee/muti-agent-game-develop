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
