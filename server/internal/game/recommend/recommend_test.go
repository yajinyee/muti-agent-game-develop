package recommend

import (
	"testing"
)

func TestRecommend_NewPlayer(t *testing.T) {
	e := New()
	b := PlayerBehavior{TotalShots: 5}
	recs := e.Analyze(b)
	if len(recs) == 0 {
		t.Error("expected at least one recommendation for new player")
	}
	if recs[0].Type != RecommendBetStay {
		t.Errorf("expected bet_stay for new player, got %s", recs[0].Type)
	}
}

func TestRecommend_LowCoins_BetDown(t *testing.T) {
	e := New()
	b := PlayerBehavior{
		TotalShots:   100,
		TotalKills:   30,
		TotalBet:     1000,
		TotalReward:  900,
		CurrentBetLv: 5,
		CurrentCoins: 50, // 很少金幣
	}
	recs := e.Analyze(b)
	found := false
	for _, r := range recs {
		if r.Type == RecommendBetDown {
			found = true
			if r.TargetBetLv >= 5 {
				t.Errorf("expected lower bet level, got %d", r.TargetBetLv)
			}
		}
	}
	if !found {
		t.Error("expected bet_down recommendation for low coins")
	}
}

func TestRecommend_HighRTP_BetUp(t *testing.T) {
	e := New()
	b := PlayerBehavior{
		TotalShots:   100,
		TotalKills:   50,
		TotalBet:     1000,
		TotalReward:  1500, // RTP = 150%
		CurrentBetLv: 3,
		CurrentCoins: 100000,
	}
	recs := e.Analyze(b)
	found := false
	for _, r := range recs {
		if r.Type == RecommendBetUp {
			found = true
			if r.TargetBetLv <= 3 {
				t.Errorf("expected higher bet level, got %d", r.TargetBetLv)
			}
		}
	}
	if !found {
		t.Error("expected bet_up recommendation for high RTP")
	}
}

func TestRecommend_LowRTP_BetDown(t *testing.T) {
	e := New()
	b := PlayerBehavior{
		TotalShots:   200,
		TotalKills:   20,
		TotalBet:     2000,
		TotalReward:  1000, // RTP = 50%
		CurrentBetLv: 7,
		CurrentCoins: 50000,
	}
	recs := e.Analyze(b)
	found := false
	for _, r := range recs {
		if r.Type == RecommendBetDown {
			found = true
		}
	}
	if !found {
		t.Error("expected bet_down recommendation for low RTP")
	}
}

func TestRecommend_LowHitRate_AutoMode(t *testing.T) {
	e := New()
	b := PlayerBehavior{
		TotalShots:   100,
		TotalKills:   2, // 命中率 2%
		TotalBet:     1000,
		TotalReward:  200,
		CurrentBetLv: 3,
		CurrentCoins: 10000,
	}
	recs := e.Analyze(b)
	found := false
	for _, r := range recs {
		if r.Type == RecommendAutoMode {
			found = true
		}
	}
	if !found {
		t.Error("expected auto_mode recommendation for low hit rate")
	}
}

func TestRecommend_HighStreak_LockTarget(t *testing.T) {
	e := New()
	b := PlayerBehavior{
		TotalShots:     100,
		TotalKills:     40,
		TotalBet:       1000,
		TotalReward:    1000,
		CurrentBetLv:   5,
		CurrentCoins:   20000,
		BestStreak:     15,
		BestMultiplier: 10.0, // 低倍率
	}
	recs := e.Analyze(b)
	found := false
	for _, r := range recs {
		if r.Type == RecommendLockTarget {
			found = true
		}
	}
	if !found {
		t.Error("expected lock_target recommendation for high streak but low multiplier")
	}
}

func TestRecommend_NoJackpot_JackpotRec(t *testing.T) {
	e := New()
	b := PlayerBehavior{
		TotalShots:   200,
		TotalKills:   80,
		TotalBet:     5000,
		TotalReward:  5000,
		CurrentBetLv: 6,
		CurrentCoins: 100000,
		JackpotWins:  0,
	}
	recs := e.Analyze(b)
	found := false
	for _, r := range recs {
		if r.Type == RecommendJackpot {
			found = true
		}
	}
	if !found {
		t.Error("expected jackpot recommendation for player with no jackpot wins")
	}
}

func TestRecommend_MaxThreeRecs(t *testing.T) {
	e := New()
	b := PlayerBehavior{
		TotalShots:   200,
		TotalKills:   80,
		TotalBet:     5000,
		TotalReward:  8000, // 高 RTP
		CurrentBetLv: 3,
		CurrentCoins: 200000,
		BestStreak:   20,
		JackpotWins:  0,
		VIPLevel:     4,
		LoginStreak:  10,
	}
	recs := e.Analyze(b)
	if len(recs) > 3 {
		t.Errorf("expected at most 3 recommendations, got %d", len(recs))
	}
}

func TestRecommend_Confidence_Range(t *testing.T) {
	e := New()
	b := PlayerBehavior{
		TotalShots:   100,
		TotalKills:   50,
		TotalBet:     1000,
		TotalReward:  1300,
		CurrentBetLv: 5,
		CurrentCoins: 50000,
	}
	recs := e.Analyze(b)
	for _, r := range recs {
		if r.Confidence < 0 || r.Confidence > 1 {
			t.Errorf("confidence out of range [0,1]: %f", r.Confidence)
		}
	}
}

func TestRecommend_BossKills_BossStrategy(t *testing.T) {
	e := New()
	b := PlayerBehavior{
		TotalShots:     100,
		TotalKills:     40,
		TotalBet:       1000,
		TotalReward:    1000,
		CurrentBetLv:   5,
		CurrentCoins:   20000,
		TotalBossKills: 5,
	}
	recs := e.Analyze(b)
	found := false
	for _, r := range recs {
		if r.Type == RecommendBoss {
			found = true
		}
	}
	if !found {
		t.Error("expected boss_focus recommendation for player with many boss kills")
	}
}
