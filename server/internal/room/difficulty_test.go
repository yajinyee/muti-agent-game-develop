// Package room — 房間難度系統測試（DAY-091）
package room

import (
	"testing"
)

func TestAllDifficultiesDefined(t *testing.T) {
	defs := AllDifficulties()
	if len(defs) != 4 {
		t.Errorf("expected 4 difficulties, got %d", len(defs))
	}
}

func TestDifficultyOrder(t *testing.T) {
	defs := AllDifficulties()
	expected := []Difficulty{DifficultyBeginner, DifficultyIntermediate, DifficultyAdvanced, DifficultyVIP}
	for i, def := range defs {
		if def.ID != expected[i] {
			t.Errorf("expected difficulty[%d] = %s, got %s", i, expected[i], def.ID)
		}
	}
}

func TestGetDifficulty(t *testing.T) {
	def := GetDifficulty(DifficultyVIP)
	if def.ID != DifficultyVIP {
		t.Errorf("expected VIP, got %s", def.ID)
	}
	if def.RewardMult != 2.0 {
		t.Errorf("VIP reward mult should be 2.0, got %.1f", def.RewardMult)
	}
	if def.JackpotMult != 5.0 {
		t.Errorf("VIP jackpot mult should be 5.0, got %.1f", def.JackpotMult)
	}
}

func TestGetDifficultyInvalid(t *testing.T) {
	// 無效難度應回傳初級
	def := GetDifficulty("invalid")
	if def.ID != DifficultyBeginner {
		t.Errorf("invalid difficulty should fallback to beginner, got %s", def.ID)
	}
}

func TestGetDifficultyByBetLevel(t *testing.T) {
	tests := []struct {
		betLevel int
		expected Difficulty
	}{
		{1, DifficultyBeginner},
		{4, DifficultyBeginner},
		{5, DifficultyIntermediate},
		{7, DifficultyIntermediate},
		{8, DifficultyAdvanced},
		{9, DifficultyAdvanced},
		{10, DifficultyVIP},
	}
	for _, tt := range tests {
		got := GetDifficultyByBetLevel(tt.betLevel)
		if got != tt.expected {
			t.Errorf("betLevel=%d: expected %s, got %s", tt.betLevel, tt.expected, got)
		}
	}
}

func TestDifficultyRewardMultIncreasing(t *testing.T) {
	// 難度越高，獎勵倍率越高
	defs := AllDifficulties()
	for i := 1; i < len(defs); i++ {
		if defs[i].RewardMult < defs[i-1].RewardMult {
			t.Errorf("difficulty[%d].RewardMult (%.1f) should be >= difficulty[%d].RewardMult (%.1f)",
				i, defs[i].RewardMult, i-1, defs[i-1].RewardMult)
		}
	}
}

func TestDifficultyJackpotMultIncreasing(t *testing.T) {
	// 難度越高，Jackpot 倍率越高
	defs := AllDifficulties()
	for i := 1; i < len(defs); i++ {
		if defs[i].JackpotMult < defs[i-1].JackpotMult {
			t.Errorf("difficulty[%d].JackpotMult (%.1f) should be >= difficulty[%d].JackpotMult (%.1f)",
				i, defs[i].JackpotMult, i-1, defs[i-1].JackpotMult)
		}
	}
}

func TestBeginnerNoEntryFee(t *testing.T) {
	def := GetDifficulty(DifficultyBeginner)
	if def.EntryFee != 0 {
		t.Errorf("beginner should have no entry fee, got %d", def.EntryFee)
	}
}

func TestVIPHasEntryFee(t *testing.T) {
	def := GetDifficulty(DifficultyVIP)
	if def.EntryFee <= 0 {
		t.Errorf("VIP should have entry fee > 0, got %d", def.EntryFee)
	}
}

func TestDifficultyMaxPlayersDecreasing(t *testing.T) {
	// 難度越高，最大玩家數越少（更私密）
	defs := AllDifficulties()
	for i := 1; i < len(defs); i++ {
		if defs[i].MaxPlayers > defs[i-1].MaxPlayers {
			t.Errorf("difficulty[%d].MaxPlayers (%d) should be <= difficulty[%d].MaxPlayers (%d)",
				i, defs[i].MaxPlayers, i-1, defs[i-1].MaxPlayers)
		}
	}
}

func TestDifficultyHasRequiredFields(t *testing.T) {
	for _, def := range AllDifficulties() {
		if def.Name == "" {
			t.Errorf("difficulty %s has empty name", def.ID)
		}
		if def.Icon == "" {
			t.Errorf("difficulty %s has empty icon", def.ID)
		}
		if def.Color == "" {
			t.Errorf("difficulty %s has empty color", def.ID)
		}
		if def.Description == "" {
			t.Errorf("difficulty %s has empty description", def.ID)
		}
		if def.MinBetCost <= 0 {
			t.Errorf("difficulty %s has invalid MinBetCost: %d", def.ID, def.MinBetCost)
		}
		if def.MaxBetCost <= def.MinBetCost {
			t.Errorf("difficulty %s MaxBetCost (%d) should be > MinBetCost (%d)",
				def.ID, def.MaxBetCost, def.MinBetCost)
		}
	}
}
