// Package combat 處理攻擊、命中、獎勵計算
package combat

import (
	"math"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/target"
)

// AttackRequest 攻擊請求
type AttackRequest struct {
	PlayerID    string
	TargetID    string // 目標 InstanceID，空字串表示自由攻擊
	BetLevel    int
	IsAuto      bool
	IsLock      bool
	ClickX      float64
	ClickY      float64
}

// AttackResult 攻擊結果
type AttackResult struct {
	AttackID    string
	PlayerID    string
	TargetID    string
	BetCost     int
	Damage      int
	IsHit       bool
	IsKill      bool
	Multiplier  float64
	Reward      int
	LaborGain   int
	CharacterID string
	// 特殊事件
	BossPhaseChanged bool
	BossPhase        int
}

// ProcessAttack 處理一次攻擊
func ProcessAttack(req AttackRequest, t *target.Target) *AttackResult {
	bet := data.GetBetDef(req.BetLevel)
	char := data.GetCharacterByBetLevel(req.BetLevel)

	result := &AttackResult{
		PlayerID:    req.PlayerID,
		TargetID:    req.TargetID,
		BetCost:     bet.BetCost,
		CharacterID: char.ID,
		IsHit:       true,
	}

	if t == nil || !t.IsAlive {
		result.IsHit = false
		return result
	}

	// 擊破判定（混合制）
	isKill, damage := t.TryKill(bet.BetCost, char.KillModifier)
	result.Damage = damage
	result.IsKill = isKill

	if isKill {
		result.Multiplier = t.Multiplier
		result.Reward = int(math.Round(float64(bet.BetCost) * t.Multiplier))
		result.LaborGain = int(float64(t.Def.LaborGain) * char.LaborModifier)

		// BOSS 獎勵依剩餘時間計算（規格書 28.3）
		if t.Def.Type == data.TargetTypeBoss {
			result.Reward = calcBossReward(bet.BetCost, t)
		}
	}

	// BOSS 階段變化
	if t.Def.Type == data.TargetTypeBoss {
		if t.UpdateBossPhase() {
			result.BossPhaseChanged = true
			result.BossPhase = t.Phase
		}
	}

	return result
}

// calcBossReward BOSS 獎勵計算（規格書 28.3 — 依真實剩餘時間）
func calcBossReward(betCost int, t *target.Target) int {
	elapsed := time.Since(t.SpawnedAt).Seconds()
	remaining := data.BossDuration - elapsed
	if remaining < 0 {
		remaining = 0
	}
	multiplier := CalcBossTimeMultiplier(remaining)
	return int(math.Round(float64(betCost) * multiplier))
}

// CalcBonusReward Bonus Game 獎勵計算（規格書 29.4）
// 倍率範圍：20-50x（Prototype 展示版，讓玩家有爽感）
// 正式版需要數值工程師依 RTP 分配精確調整
func CalcBonusReward(entryBetCost int, bonusScore int) (int, float64) {
	// 倍率：20 + score × 0.375，上限 50
	multiplier := math.Min(20+float64(bonusScore)*0.375, 50)
	reward := int(math.Round(float64(entryBetCost) * multiplier))
	return reward, multiplier
}

// CalcBossTimeMultiplier BOSS 依剩餘時間計算倍率（規格書 28.3）
func CalcBossTimeMultiplier(remainingSeconds float64) float64 {
	switch {
	case remainingSeconds <= 10:
		return 100
	case remainingSeconds <= 20:
		return 150
	case remainingSeconds <= 30:
		return 200
	case remainingSeconds <= 40:
		return 300
	case remainingSeconds <= 50:
		return 400
	default:
		return 500
	}
}
