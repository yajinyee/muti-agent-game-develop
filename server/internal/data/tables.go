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
