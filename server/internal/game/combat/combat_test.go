package combat

import (
	"testing"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/target"
)

// TestCalcBossTimeMultiplier 測試 BOSS 時間倍率計算（規格書 28.3）
func TestCalcBossTimeMultiplier(t *testing.T) {
	tests := []struct {
		name      string
		remaining float64
		expected  float64
	}{
		{"最後 10 秒", 5.0, 100},
		{"10 秒整", 10.0, 100},
		{"11 秒", 11.0, 150},
		{"20 秒整", 20.0, 150},
		{"21 秒", 21.0, 200},
		{"30 秒整", 30.0, 200},
		{"31 秒", 31.0, 300},
		{"40 秒整", 40.0, 300},
		{"41 秒", 41.0, 400},
		{"50 秒整", 50.0, 400},
		{"51 秒（最高倍率）", 51.0, 500},
		{"60 秒（滿時間）", 60.0, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcBossTimeMultiplier(tt.remaining)
			if got != tt.expected {
				t.Errorf("CalcBossTimeMultiplier(%.1f) = %.0f, want %.0f", tt.remaining, got, tt.expected)
			}
		})
	}
}

// TestCalcBonusReward 測試 Bonus 獎勵計算（規格書 29.4）
func TestCalcBonusReward(t *testing.T) {
	tests := []struct {
		name        string
		entryBet    int
		bonusScore  int
		minReward   int
		maxReward   int
		minMult     float64
		maxMult     float64
	}{
		{"最低分（0分）", 10, 0, 200, 200, 20.0, 20.0},
		{"中等分（40分）", 10, 40, 350, 350, 35.0, 35.0},
		{"滿分（80分）", 10, 80, 500, 500, 50.0, 50.0},
		{"超過上限（100分）", 10, 100, 500, 500, 50.0, 50.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reward, mult := CalcBonusReward(tt.entryBet, tt.bonusScore)
			if reward < tt.minReward || reward > tt.maxReward {
				t.Errorf("CalcBonusReward(%d, %d) reward = %d, want [%d, %d]",
					tt.entryBet, tt.bonusScore, reward, tt.minReward, tt.maxReward)
			}
			if mult < tt.minMult || mult > tt.maxMult {
				t.Errorf("CalcBonusReward(%d, %d) mult = %.1f, want [%.1f, %.1f]",
					tt.entryBet, tt.bonusScore, mult, tt.minMult, tt.maxMult)
			}
		})
	}
}

// TestProcessAttack_NilTarget 測試攻擊 nil 目標
func TestProcessAttack_NilTarget(t *testing.T) {
	req := AttackRequest{
		PlayerID: "p1",
		BetLevel: 1,
	}
	result := ProcessAttack(req, nil)
	if result.IsHit {
		t.Error("攻擊 nil 目標應該 IsHit = false")
	}
	if result.IsKill {
		t.Error("攻擊 nil 目標應該 IsKill = false")
	}
}

// TestProcessAttack_DeadTarget 測試攻擊已死亡目標
func TestProcessAttack_DeadTarget(t *testing.T) {
	def, ok := data.Targets["T001"]
	if !ok {
		t.Skip("T001 定義不存在，跳過測試")
	}
	tgt := target.NewTarget("inst_001", def, 100, 100)
	tgt.IsAlive = false

	req := AttackRequest{
		PlayerID: "p1",
		BetLevel: 1,
	}
	result := ProcessAttack(req, tgt)
	if result.IsHit {
		t.Error("攻擊已死亡目標應該 IsHit = false")
	}
}

// TestProcessAttack_ValidTarget 測試攻擊有效目標
func TestProcessAttack_ValidTarget(t *testing.T) {
	def, ok := data.Targets["T001"]
	if !ok {
		t.Skip("T001 定義不存在，跳過測試")
	}
	tgt := target.NewTarget("inst_001", def, 100, 100)

	req := AttackRequest{
		PlayerID: "p1",
		BetLevel: 1,
	}
	result := ProcessAttack(req, tgt)

	// 攻擊有效目標應該 IsHit = true
	if !result.IsHit {
		t.Error("攻擊有效目標應該 IsHit = true")
	}
	// BetCost 應該 > 0
	if result.BetCost <= 0 {
		t.Errorf("BetCost 應該 > 0，got %d", result.BetCost)
	}
	// 如果擊殺，獎勵應該 > 0
	if result.IsKill && result.Reward <= 0 {
		t.Errorf("擊殺後 Reward 應該 > 0，got %d", result.Reward)
	}
}

// TestProcessAttack_BetLevels 測試不同投注等級的攻擊
func TestProcessAttack_BetLevels(t *testing.T) {
	def, ok := data.Targets["T001"]
	if !ok {
		t.Skip("T001 定義不存在，跳過測試")
	}

	for level := 1; level <= 10; level++ {
		tgt := target.NewTarget("inst_001", def, 100, 100)
		req := AttackRequest{
			PlayerID: "p1",
			BetLevel: level,
		}
		result := ProcessAttack(req, tgt)
		if !result.IsHit {
			t.Errorf("BetLevel %d: 攻擊有效目標應該 IsHit = true", level)
		}
		if result.BetCost <= 0 {
			t.Errorf("BetLevel %d: BetCost 應該 > 0，got %d", level, result.BetCost)
		}
	}
}
