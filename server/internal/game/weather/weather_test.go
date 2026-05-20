// weather_test.go — 天氣系統單元測試（DAY-087）
package weather

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := New()
	if m.GetCurrent() != WeatherClear {
		t.Errorf("expected initial weather to be clear, got %s", m.GetCurrent())
	}
}

func TestGetSnapshot(t *testing.T) {
	m := New()
	snap := m.GetSnapshot(false)
	if snap.Type != WeatherClear {
		t.Errorf("expected clear weather, got %s", snap.Type)
	}
	if snap.RewardMult != 1.0 {
		t.Errorf("expected reward mult 1.0, got %.2f", snap.RewardMult)
	}
	if snap.IsNew {
		t.Error("expected IsNew=false")
	}
}

func TestGetSnapshotIsNew(t *testing.T) {
	m := New()
	snap := m.GetSnapshot(true)
	if !snap.IsNew {
		t.Error("expected IsNew=true")
	}
}

func TestCheckAndRotate_NoChange(t *testing.T) {
	m := New()
	// 剛建立，不應該切換
	changed, _ := m.CheckAndRotate()
	if changed {
		t.Error("expected no change immediately after creation")
	}
}

func TestCheckAndRotate_Change(t *testing.T) {
	m := New()
	// 強制讓天氣過期
	m.startedAt = time.Now().Add(-10 * time.Minute)
	changed, snap := m.CheckAndRotate()
	if !changed {
		t.Error("expected weather to change after duration expired")
	}
	if snap.Type == WeatherClear {
		// 可能還是 clear，但機率很低（40/60 = 67%），重試一次
		m.startedAt = time.Now().Add(-10 * time.Minute)
		changed2, snap2 := m.CheckAndRotate()
		if !changed2 {
			t.Error("expected weather to change on second try")
		}
		_ = snap2
	}
	if snap.IsNew != true {
		t.Error("expected IsNew=true after rotation")
	}
}

func TestCheckAndRotate_NoDuplicate(t *testing.T) {
	m := New()
	// 多次切換，確認不會連續出現同一天氣
	prev := m.GetCurrent()
	for i := 0; i < 20; i++ {
		m.startedAt = time.Now().Add(-10 * time.Minute)
		changed, snap := m.CheckAndRotate()
		if changed && snap.Type == prev {
			t.Errorf("weather repeated: %s → %s", prev, snap.Type)
		}
		if changed {
			prev = snap.Type
		}
	}
}

func TestGetRewardMult(t *testing.T) {
	m := New()
	// 晴天獎勵倍率 = 1.0
	if m.GetRewardMult() != 1.0 {
		t.Errorf("expected 1.0, got %.2f", m.GetRewardMult())
	}
}

func TestGetSpeedMult(t *testing.T) {
	m := New()
	if m.GetSpeedMult() != 1.0 {
		t.Errorf("expected 1.0, got %.2f", m.GetSpeedMult())
	}
}

func TestWeatherDefs_AllPresent(t *testing.T) {
	expected := []WeatherType{
		WeatherClear, WeatherRain, WeatherStorm,
		WeatherFog, WeatherSunshine, WeatherBlizzard,
	}
	for _, wt := range expected {
		if _, ok := WeatherDefs[wt]; !ok {
			t.Errorf("missing weather def: %s", wt)
		}
	}
}

func TestWeatherDefs_ValidWeights(t *testing.T) {
	for wt, def := range WeatherDefs {
		if def.Weight <= 0 {
			t.Errorf("weather %s has invalid weight: %d", wt, def.Weight)
		}
		if def.RewardMult <= 0 {
			t.Errorf("weather %s has invalid reward mult: %.2f", wt, def.RewardMult)
		}
		if def.SpeedMult <= 0 {
			t.Errorf("weather %s has invalid speed mult: %.2f", wt, def.SpeedMult)
		}
	}
}

func TestPickNext_NeverSame(t *testing.T) {
	m := New()
	// 測試 pickNext 不會選到當前天氣
	for _, wt := range weatherOrder {
		m.current = wt
		for i := 0; i < 50; i++ {
			next := m.pickNext()
			if next == wt {
				t.Errorf("pickNext returned same weather: %s", wt)
			}
		}
	}
}

func TestRemainingSeconds(t *testing.T) {
	m := New()
	snap := m.GetSnapshot(false)
	// 剛建立，剩餘時間應接近 Duration
	def := WeatherDefs[WeatherClear]
	expected := int(def.Duration.Seconds())
	if snap.RemainingSeconds < expected-2 || snap.RemainingSeconds > expected {
		t.Errorf("expected remaining ~%d, got %d", expected, snap.RemainingSeconds)
	}
}

func TestGetBossChanceBonus_Blizzard(t *testing.T) {
	m := New()
	m.mu.Lock()
	m.current = WeatherBlizzard
	m.mu.Unlock()
	if m.GetBossChanceBonus() != 0.30 {
		t.Errorf("expected 0.30, got %.2f", m.GetBossChanceBonus())
	}
}

func TestGetGoldFishBonus_Sunshine(t *testing.T) {
	m := New()
	m.mu.Lock()
	m.current = WeatherSunshine
	m.mu.Unlock()
	if m.GetGoldFishBonus() != 0.50 {
		t.Errorf("expected 0.50, got %.2f", m.GetGoldFishBonus())
	}
}

func TestGetFogEffect(t *testing.T) {
	m := New()
	m.mu.Lock()
	m.current = WeatherFog
	m.mu.Unlock()
	if !m.GetFogEffect() {
		t.Error("expected fog effect for fog weather")
	}
	m.mu.Lock()
	m.current = WeatherClear
	m.mu.Unlock()
	if m.GetFogEffect() {
		t.Error("expected no fog effect for clear weather")
	}
}
