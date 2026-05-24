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
