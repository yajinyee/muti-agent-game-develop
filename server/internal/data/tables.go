// Package data 定義遊戲所有靜態資料表（來自規格書）
package data

// TargetType 目標類型
type TargetType string

const (
	TargetTypeBasic   TargetType = "basic"
	TargetTypeSpecial TargetType = "special"
	TargetTypeBoss    TargetType = "boss"
	TargetTypeBonus   TargetType = "bonus"
)

// TargetDef 目標物定義（對應規格書 Target Table）
type TargetDef struct {
	ID              string
	Name            string
	Type            TargetType
	MultiplierMin   float64
	MultiplierMax   float64
	HP              int
	SpawnWeight     int
	Speed           float64 // pixels/sec
	Lifetime        float64 // seconds
	LaborGain       int
	DifficultyFactor float64
	SpecialBehavior string
}

// CharacterDef 角色定義
type CharacterDef struct {
	ID               string
	Name             string
	BetLevelMin      int
	BetLevelMax      int
	AttackColor      string
	KillModifier     float64
	FireRateModifier float64
	LaborModifier    float64
	VoiceText        string
}

// BetDef 投注等級定義
type BetDef struct {
	Level           int
	CharacterID     string
	BetCost         int
	AttackPower     int
	FireRate        float64 // shots/sec
	ProjectileSpeed float64 // pixels/sec
}

// BonusTargetDef Bonus Game 目標定義
type BonusTargetDef struct {
	ID           string
	Name         string
	ClickScore   int
	SpawnWeight  int
	SpecialEffect string
}

// ---- 靜態資料表 ----

// Targets 所有目標物（規格書 26.1）
// DifficultyFactor 修正說明（2026-05-12）：
//   正確公式：required_hits = ceil(multiplier / bet_cost × DifficultyFactor)
//   要讓保底在期望命中次數（= multiplier / RTP）的 1.5 倍觸發：
//   DifficultyFactor = 1.5 × bet_cost / RTP ≈ 1.5 × bet_cost / 0.94
//   對 bet_cost=10：DifficultyFactor ≈ 16
//   這樣 T001(2x,bet=10)：required = ceil(2/10×16) = ceil(3.2) = 4 次保底
//   期望命中 ≈ 1/0.47 = 2.1 次，保底 4 次，RTP ≈ 94% ✓
var Targets = map[string]*TargetDef{
	"T001": {ID: "T001", Name: "像素雜草", Type: TargetTypeBasic, MultiplierMin: 2, MultiplierMax: 2, HP: 3, SpawnWeight: 180, Speed: 0, Lifetime: 20, LaborGain: 1, DifficultyFactor: 16.0, SpecialBehavior: "static_sway"},
	"T002": {ID: "T002", Name: "綠色小蟲", Type: TargetTypeBasic, MultiplierMin: 3, MultiplierMax: 3, HP: 5, SpawnWeight: 160, Speed: 40, Lifetime: 18, LaborGain: 1, DifficultyFactor: 16.0, SpecialBehavior: "linear"},
	"T003": {ID: "T003", Name: "紅色小蟲", Type: TargetTypeBasic, MultiplierMin: 5, MultiplierMax: 5, HP: 8, SpawnWeight: 130, Speed: 55, Lifetime: 16, LaborGain: 1, DifficultyFactor: 16.0, SpecialBehavior: "jump"},
	"T004": {ID: "T004", Name: "藍色小蟲", Type: TargetTypeBasic, MultiplierMin: 6, MultiplierMax: 6, HP: 10, SpawnWeight: 110, Speed: 65, Lifetime: 15, LaborGain: 2, DifficultyFactor: 16.0, SpecialBehavior: "curve"},
	"T005": {ID: "T005", Name: "會走路的布丁", Type: TargetTypeBasic, MultiplierMin: 8, MultiplierMax: 8, HP: 16, SpawnWeight: 90, Speed: 35, Lifetime: 20, LaborGain: 2, DifficultyFactor: 16.0, SpecialBehavior: "sway"},
	"T006": {ID: "T006", Name: "巨大蘑菇", Type: TargetTypeBasic, MultiplierMin: 10, MultiplierMax: 10, HP: 22, SpawnWeight: 70, Speed: 25, Lifetime: 22, LaborGain: 3, DifficultyFactor: 16.0, SpecialBehavior: "sink"},
	"T101": {ID: "T101", Name: "擬態型怪物", Type: TargetTypeSpecial, MultiplierMin: 15, MultiplierMax: 30, HP: 35, SpawnWeight: 35, Speed: 50, Lifetime: 14, LaborGain: 5, DifficultyFactor: 16.0, SpecialBehavior: "mimic"},
	"T102": {ID: "T102", Name: "寶箱怪", Type: TargetTypeSpecial, MultiplierMin: 25, MultiplierMax: 25, HP: 55, SpawnWeight: 22, Speed: 70, Lifetime: 10, LaborGain: 6, DifficultyFactor: 16.0, SpecialBehavior: "flee"},
	"T103": {ID: "T103", Name: "流星", Type: TargetTypeSpecial, MultiplierMin: 20, MultiplierMax: 50, HP: 20, SpawnWeight: 18, Speed: 220, Lifetime: 4, LaborGain: 5, DifficultyFactor: 16.0, SpecialBehavior: "meteor"},
	"T104": {ID: "T104", Name: "金色雜草", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 30, HP: 45, SpawnWeight: 12, Speed: 0, Lifetime: 8, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "static"},
	"T105": {ID: "T105", Name: "巨大金幣魚", Type: TargetTypeSpecial, MultiplierMin: 50, MultiplierMax: 50, HP: 90, SpawnWeight: 8, Speed: 80, Lifetime: 8, LaborGain: 10, DifficultyFactor: 16.0, SpecialBehavior: "coin_rain"},
	// T106 鑽頭龍蝦（DAY-142）— 業界依據：Royal Fishing JILI 2026「Drill Bit Lobster (80X) — penetrating drill through multiple fish, self-detonates at end of trajectory」
	// 擊破後觸發穿透鑽頭，沿水平方向穿透所有目標，到達邊緣後爆炸，連帶擊破爆炸範圍內目標
	"T106": {ID: "T106", Name: "鑽頭龍蝦", Type: TargetTypeSpecial, MultiplierMin: 60, MultiplierMax: 80, HP: 60, SpawnWeight: 6, Speed: 45, Lifetime: 10, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "drill_lobster"},
	// T107 炸彈蟹（DAY-143）— 業界依據：royal-fishing.uk 2026「Worth 70x, this explosive crustacean triggers multiple large-scale detonations.
	// Each bomb creates expanding capture zones for massive multi-target eliminations.」
	// 擊破後觸發 3 波爆炸，每波爆炸半徑 150px，每波間隔 400ms，連帶擊破爆炸範圍內所有目標
	// T107 炸彈蟹（DAY-143）— 業界依據：royal-fishing.uk 2026「Worth 70x, this explosive crustacean triggers multiple large-scale detonations.
	// Each bomb creates expanding capture zones for massive multi-target eliminations.」
	// 擊破後觸發 3 波爆炸，每波爆炸半徑 150px，每波間隔 400ms，連帶擊破爆炸範圍內所有目標
	"T107": {ID: "T107", Name: "炸彈蟹", Type: TargetTypeSpecial, MultiplierMin: 60, MultiplierMax: 70, HP: 70, SpawnWeight: 5, Speed: 35, Lifetime: 12, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "bomb_crab"},
	// T108 巨型章魚（DAY-144）— 業界依據：JILI Mega Fishing「Mega Octopus Wheel – Defeat that giant octopus and enter
	// the bonus wheel round where you have a chance to win massive guaranteed prizes up to 950x.」
	// 擊破後觸發個人轉盤（8格：50x-950x），玩家點擊停止，結果預先決定（公平性保證）
	"T108": {ID: "T108", Name: "巨型章魚", Type: TargetTypeSpecial, MultiplierMin: 80, MultiplierMax: 120, HP: 120, SpawnWeight: 3, Speed: 30, Lifetime: 15, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "mega_octopus"},
	// T109 巨型鮟鱇魚（DAY-145）— 業界依據：jiligames.com 2026「Giant Anglerfish can shoot electricity to open treasure chests」
	// 擊破後觸發電擊，電流傳導到附近的寶箱目標（T102），強制開啟寶箱獲得額外獎勵
	"T109": {ID: "T109", Name: "巨型鮟鱇魚", Type: TargetTypeSpecial, MultiplierMin: 70, MultiplierMax: 90, HP: 90, SpawnWeight: 4, Speed: 25, Lifetime: 14, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "anglerfish_shock"},
	// T110 巨型鹹水鱷魚（DAY-146）— 業界依據：jiligames.com 2026「giant crocodiles awaken to hunt fish on the fish farm to accumulate big prizes!」
	// + megafishinggame.top「Giant Saltwater Crocodile」
	// 擊破後觸發「鱷魚獵魚」模式：鱷魚在 8 秒內自動獵殺場上的普通目標（T001-T006），累積獎勵給觸發玩家
	"T110": {ID: "T110", Name: "巨型鹹水鱷魚", Type: TargetTypeSpecial, MultiplierMin: 100, MultiplierMax: 150, HP: 150, SpawnWeight: 2, Speed: 20, Lifetime: 18, LaborGain: 18, DifficultyFactor: 16.0, SpecialBehavior: "crocodile_hunt"},
	// T111 夢幻巨型獎勵魚（DAY-147）— 業界依據：jiligames.com 2026「The dreamy Giant Prize Fish lets you easily win great prizes, with the chance for 5x multipliers」
	// 擊破後觸發「夢幻獎勵模式」：觸發玩家在 10 秒內所有擊破獎勵 ×5，讓玩家感受到「夢幻大獎」的爽感
	// 設計：低 HP（容易擊破）+ 中等倍率（40-60x）+ 觸發後 10 秒 5x 加成，是「容易觸發的短期爆發」機制
	"T111": {ID: "T111", Name: "夢幻巨型獎勵魚", Type: TargetTypeSpecial, MultiplierMin: 40, MultiplierMax: 60, HP: 80, SpawnWeight: 4, Speed: 35, Lifetime: 12, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "giant_prize_fish"},
	// T112 千龍王（DAY-148）— 業界依據：Royal Fishing JILI 2026「ChainLong King — capture this golden dragon to trigger
	// the dual-ring roulette. The ChainLong King itself can award up to 1000X mega wins.」
	// 擊破後觸發「千龍王強化輪盤」：內環（5x-50x）× 外環（2x-20x）= 最高 1000x
	// 設計：超高倍率（150-1000x）+ 超高 HP（300）+ 極低生成權重（1）= 終極稀有目標
	// 千龍王是全遊戲最高倍率目標，擊破後觸發專屬強化輪盤，最高 1000x 是業界最高水準
	"T112": {ID: "T112", Name: "千龍王", Type: TargetTypeSpecial, MultiplierMin: 150, MultiplierMax: 1000, HP: 300, SpawnWeight: 1, Speed: 15, Lifetime: 20, LaborGain: 30, DifficultyFactor: 16.0, SpecialBehavior: "chainlong_king"},
	// T113 黃金水母（DAY-149）— 業界依據：Ocean King 3 2026「Electric Jellyfish chain shocks across multiple targets.
	// Devastating against clustered schools.」— 擊破後觸發「全場電擊」，對畫面上所有目標發動電擊
	// 比閃電鰻（T103，200px 範圍跳躍 5 次）更強：全場範圍，最多 8 個目標，40% 擊破機率
	// 設計：中等倍率（60-80x）+ 中等 HP（80）+ 低生成權重（3）= 稀有但可遇到的強力目標
	"T113": {ID: "T113", Name: "黃金水母", Type: TargetTypeSpecial, MultiplierMin: 60, MultiplierMax: 80, HP: 80, SpawnWeight: 3, Speed: 30, Lifetime: 12, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "golden_jellyfish"},
	"B001": {ID: "B001", Name: "那個孩子", Type: TargetTypeBoss, MultiplierMin: 100, MultiplierMax: 500, HP: 3000, SpawnWeight: 0, Speed: 20, Lifetime: 60, LaborGain: 30, DifficultyFactor: 16.0, SpecialBehavior: "boss_phases"},
}

// MeteorMultiplierWeights 流星倍率權重（規格書 26.3）
var MeteorMultiplierWeights = []struct {
	Multiplier float64
	Weight     int
}{
	{20, 50},
	{30, 30},
	{40, 15},
	{50, 5},
}

// Characters 角色定義（規格書 5章）
var Characters = map[string]*CharacterDef{
	"chiikawa": {ID: "chiikawa", Name: "吉伊卡哇", BetLevelMin: 1, BetLevelMax: 3, AttackColor: "pink", KillModifier: 1.00, FireRateModifier: 1.00, LaborModifier: 1.10, VoiceText: "YaDa"},
	"hachiware": {ID: "hachiware", Name: "小八", BetLevelMin: 4, BetLevelMax: 7, AttackColor: "blue", KillModifier: 1.00, FireRateModifier: 1.08, LaborModifier: 1.00, VoiceText: "尖尖哇嘎乃"},
	"usagi": {ID: "usagi", Name: "烏薩奇", BetLevelMin: 8, BetLevelMax: 10, AttackColor: "yellow", KillModifier: 0.98, FireRateModifier: 1.20, LaborModifier: 0.95, VoiceText: "Yaha"},
}

// BetLevels 投注等級表（規格書 6章 & 25.5）
var BetLevels = []*BetDef{
	{Level: 1, CharacterID: "chiikawa", BetCost: 1, AttackPower: 1, FireRate: 2.0, ProjectileSpeed: 700},
	{Level: 2, CharacterID: "chiikawa", BetCost: 2, AttackPower: 2, FireRate: 2.0, ProjectileSpeed: 720},
	{Level: 3, CharacterID: "chiikawa", BetCost: 3, AttackPower: 3, FireRate: 2.1, ProjectileSpeed: 740},
	{Level: 4, CharacterID: "hachiware", BetCost: 5, AttackPower: 5, FireRate: 2.2, ProjectileSpeed: 780},
	{Level: 5, CharacterID: "hachiware", BetCost: 10, AttackPower: 10, FireRate: 2.3, ProjectileSpeed: 800},
	{Level: 6, CharacterID: "hachiware", BetCost: 20, AttackPower: 20, FireRate: 2.4, ProjectileSpeed: 820},
	{Level: 7, CharacterID: "hachiware", BetCost: 30, AttackPower: 30, FireRate: 2.5, ProjectileSpeed: 850},
	{Level: 8, CharacterID: "usagi", BetCost: 50, AttackPower: 50, FireRate: 2.7, ProjectileSpeed: 900},
	{Level: 9, CharacterID: "usagi", BetCost: 80, AttackPower: 80, FireRate: 2.9, ProjectileSpeed: 940},
	{Level: 10, CharacterID: "usagi", BetCost: 100, AttackPower: 100, FireRate: 3.0, ProjectileSpeed: 980},
}

// BonusTargets Bonus Game 目標（規格書 29.3）
var BonusTargets = []*BonusTargetDef{
	{ID: "BG001", Name: "普通雜草", ClickScore: 1, SpawnWeight: 180, SpecialEffect: "none"},
	{ID: "BG002", Name: "硬雜草", ClickScore: 3, SpawnWeight: 80, SpecialEffect: "double_click"},
	{ID: "BG003", Name: "發光雜草", ClickScore: 8, SpawnWeight: 35, SpecialEffect: "multiplier_up"},
	{ID: "BG004", Name: "金色雜草", ClickScore: 20, SpawnWeight: 10, SpecialEffect: "coin_shower"},
	{ID: "BG005", Name: "搗亂怪草", ClickScore: -5, SpawnWeight: 20, SpecialEffect: "stun"},
}

// GetBetDef 取得投注等級定義
func GetBetDef(level int) *BetDef {
	if level < 1 || level > 10 {
		return BetLevels[0]
	}
	return BetLevels[level-1]
}

// GetCharacterByBetLevel 依投注等級取得角色
func GetCharacterByBetLevel(level int) *CharacterDef {
	bet := GetBetDef(level)
	return Characters[bet.CharacterID]
}

// BaseRTPFactor 基礎 RTP 係數（規格書 30章）
const BaseRTPFactor = 0.92

// LaborValueMax 勞動值上限
const LaborValueMax = 100

// SpawnInterval 目標生成間隔（秒）
const SpawnInterval = 0.8

// MaxTargetsOnScreen 畫面最大目標數
const MaxTargetsOnScreen = 18

// BossDuration BOSS 持續時間（秒）
const BossDuration = 60.0

// BonusDuration Bonus Game 持續時間（秒）
const BonusDuration = 15.0

// ---- 武器升級系統（DAY-067）----

// WeaponDef 武器定義
type WeaponDef struct {
	Level       int
	Name        string
	Icon        string    // 顯示圖示
	PowerMod    float64   // 攻擊力加成係數（1.0=無加成，1.25=+25%）
	ExtraCost   int       // 每次攻擊額外扣除的金幣（在 BetCost 之外）
	Color       string    // 投射物顏色（Client 端視覺）
	Description string
}

// Weapons 武器等級定義（DAY-067）
var Weapons = []*WeaponDef{
	{
		Level:       1,
		Name:        "標準砲",
		Icon:        "🔫",
		PowerMod:    1.00,
		ExtraCost:   0,
		Color:       "white",
		Description: "標準攻擊力，無額外費用",
	},
	{
		Level:       2,
		Name:        "強化砲",
		Icon:        "⚡",
		PowerMod:    1.25,
		ExtraCost:   50,  // 每次攻擊額外扣 50 金幣
		Color:       "cyan",
		Description: "攻擊力 +25%，每次攻擊額外消耗 50 金幣",
	},
	{
		Level:       3,
		Name:        "超級砲",
		Icon:        "🌟",
		PowerMod:    1.60,
		ExtraCost:   150, // 每次攻擊額外扣 150 金幣
		Color:       "gold",
		Description: "攻擊力 +60%，每次攻擊額外消耗 150 金幣",
	},
}

// GetWeaponDef 取得武器定義
func GetWeaponDef(level int) *WeaponDef {
	if level < 1 || level > len(Weapons) {
		return Weapons[0]
	}
	return Weapons[level-1]
}
