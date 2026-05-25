// Package game — 目標物生成邏輯
// server-event-agent 負責維護
package game

import (
	"math/rand"

	"chiikawa-game/internal/data"
)

const (
	SpawnInterval  = 0.6  // 秒（從 0.8 加快到 0.6，增加目標物密度）
	MaxTargets     = 22   // 從 18 增加到 22，讓畫面更豐富
	MaxBossTargets = 8
	GameWidth      = 1280.0
	GameHeight     = 720.0
	SpawnX         = GameWidth + 50  // 從右側進入
)

// spawnY 隨機生成 Y 座標（避開頂部 UI 和底部 UI）
func spawnY() float64 {
	return 80 + rand.Float64()*(GameHeight-160)
}

// pickTargetDef 依 BetLevel 和權重選擇目標物定義
// server-event-agent 負責維護此邏輯
func pickTargetDef(betLevel int) *data.TargetDef {
	// 依 BetLevel 決定特殊目標出現機率
	var specialChance float64
	switch {
	case betLevel <= 3:
		specialChance = 0.10 // 10%
	case betLevel <= 7:
		specialChance = 0.18 // 18%
	default:
		specialChance = 0.25 // 25%
	}

	// 決定是基礎還是特殊
	isSpecial := rand.Float64() < specialChance

	// 收集候選目標
	var candidates []data.TargetDef
	for _, t := range data.Targets {
		if t.Type == data.TypeBoss {
			continue // BOSS 不在正常生成池
		}
		if isSpecial && t.Type == data.TypeSpecial {
			candidates = append(candidates, t)
		} else if !isSpecial && t.Type == data.TypeBasic {
			candidates = append(candidates, t)
		}
	}

	if len(candidates) == 0 {
		return &data.Targets[0] // 備用
	}

	// 加權隨機選擇
	total := 0
	for _, c := range candidates {
		total += c.SpawnWeight
	}
	r := rand.Intn(total)
	for _, c := range candidates {
		r -= c.SpawnWeight
		if r < 0 {
			def := c
			return &def
		}
	}
	return &candidates[0]
}
