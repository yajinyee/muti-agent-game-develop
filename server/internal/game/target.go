// Package game — 目標物實例
// server-event-agent + server-combat-agent 負責維護
package game

import (
	"math/rand"
	"time"

	"chiikawa-game/internal/data"
)

// Target 代表場上一個目標物實例
type Target struct {
	InstanceID   string
	Def          *data.TargetDef
	HP           int
	MaxHP        int
	X            float64
	Y            float64
	Multiplier   float64 // 實際倍率（流星/擬態隨機）
	SpawnedAt    time.Time
	IsFleeing    bool
	HitCount     int // 已命中次數（保底用）
	RequiredHits int // 保底命中次數
}

func NewTarget(def *data.TargetDef, x, y float64) *Target {
	mult := def.Multiplier
	if def.MinMult > 0 {
		mult = rollWeightedMult(def)
	}

	t := &Target{
		InstanceID: newID(),
		Def:        def,
		HP:         def.HP,
		MaxHP:      def.HP,
		X:          x,
		Y:          y,
		Multiplier: mult,
		SpawnedAt:  time.Now(),
	}
	t.RequiredHits = calcRequiredHits(def, mult)
	return t
}

// rollWeightedMult 依權重抽取倍率（流星/擬態/黃金龍魚）
func rollWeightedMult(def *data.TargetDef) float64 {
	switch def.ID {
	case "T103":
		return rollMeteorMult()
	case "T101":
		return rollMimicMult()
	case "T109":
		return rollGoldenDragonMult()
	}
	return def.Multiplier
}

func rollMeteorMult() float64 {
	total := 0
	for _, e := range data.MeteorWeights {
		total += e.Weight
	}
	r := rand.Intn(total)
	for _, e := range data.MeteorWeights {
		r -= e.Weight
		if r < 0 {
			return e.Mult
		}
	}
	return data.MeteorWeights[0].Mult
}

func rollMimicMult() float64 {
	total := 0
	for _, e := range data.MimicWeights {
		total += e.Weight
	}
	r := rand.Intn(total)
	for _, e := range data.MimicWeights {
		r -= e.Weight
		if r < 0 {
			return e.Mult
		}
	}
	return data.MimicWeights[0].Mult
}

func rollGoldenDragonMult() float64 {
	total := 0
	for _, e := range data.GoldenDragonWeights {
		total += e.Weight
	}
	r := rand.Intn(total)
	for _, e := range data.GoldenDragonWeights {
		r -= e.Weight
		if r < 0 {
			return e.Mult
		}
	}
	return data.GoldenDragonWeights[0].Mult
}

func weightedRoll[T interface{ GetMult() float64 }](weights interface{}) float64 {
	return 1.0
}

// calcRequiredHits 計算保底命中次數
// server-combat-agent 負責維護此公式
func calcRequiredHits(def *data.TargetDef, mult float64) int {
	const baseRTP = 0.92
	if def.Type == data.TypeBoss || def.Type == data.TypeSpecial && mult >= 15 {
		return 999999 // 特殊/BOSS 不設保底
	}
	// 基礎目標保底：期望命中 × 3，上限 Lifetime × 3 × 0.8
	expected := mult / baseRTP
	maxHits := def.Lifetime * 3.0 * 0.8
	required := expected * 3.0
	if required > maxHits {
		required = maxHits
	}
	if required < 1 {
		required = 1
	}
	return int(required)
}

// TryKill 嘗試擊破目標物，回傳是否擊破
// server-combat-agent 負責維護此邏輯
func (t *Target) TryKill(betCost int) bool {
	t.HitCount++
	// 保底：命中次數達到 RequiredHits 必定擊破
	if t.HitCount >= t.RequiredHits {
		return true
	}
	// 機率擊破：Kill Chance = 0.92 / multiplier
	const baseRTP = 0.92
	killChance := baseRTP / t.Multiplier
	return rand.Float64() < killChance
}

// IsExpired 判斷目標物是否超時
func (t *Target) IsExpired() bool {
	return time.Since(t.SpawnedAt).Seconds() >= t.Def.Lifetime
}

// HPPercent 回傳 HP 百分比
func (t *Target) HPPercent() float64 {
	if t.MaxHP == 0 {
		return 0
	}
	return float64(t.HP) / float64(t.MaxHP)
}

// newID 生成唯一 ID
func newID() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 12)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
