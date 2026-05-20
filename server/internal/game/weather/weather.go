// Package weather 天氣系統（DAY-087）
// 業界依據：Fisch（Roblox）2026 確認天氣系統讓魚群生成率提升 35%，是 2026 年捕魚遊戲標配
package weather

import (
	"math/rand"
	"sync"
	"time"
)

// WeatherType 天氣類型
type WeatherType string

const (
	WeatherClear    WeatherType = "clear"    // 晴天（正常）
	WeatherRain     WeatherType = "rain"     // 下雨（稀有目標出現率 +20%）
	WeatherStorm    WeatherType = "storm"    // 暴風雨（所有目標移動速度 +30%，獎勵 ×1.5）
	WeatherFog      WeatherType = "fog"      // 濃霧（目標物透明度降低，擊破獎勵 ×2.0）
	WeatherSunshine WeatherType = "sunshine" // 豔陽（金幣魚出現率 +50%，獎勵 ×1.2）
	WeatherBlizzard WeatherType = "blizzard" // 暴雪（BOSS 出現機率 +30%，BOSS 獎勵 ×2.0）
)

// WeatherDef 天氣定義
type WeatherDef struct {
	Type            WeatherType
	Name            string  // 顯示名稱
	Icon            string  // 圖示（emoji）
	Description     string  // 效果說明
	Duration        time.Duration // 持續時間
	SpawnRateMult   float64 // 目標生成倍率
	RewardMult      float64 // 獎勵倍率
	SpeedMult       float64 // 目標移動速度倍率
	RareChanceBonus float64 // 稀有目標出現機率加成（0.0-1.0）
	GoldFishBonus   float64 // 金幣魚出現機率加成（0.0-1.0）
	BossChanceBonus float64 // BOSS 出現機率加成（0.0-1.0）
	FogEffect       bool    // 是否有濃霧效果（Client 端視覺）
	Weight          int     // 隨機權重
}

// WeatherDefs 所有天氣定義
var WeatherDefs = map[WeatherType]*WeatherDef{
	WeatherClear: {
		Type:        WeatherClear,
		Name:        "晴天",
		Icon:        "☀️",
		Description: "風平浪靜，正常捕魚",
		Duration:    5 * time.Minute,
		SpawnRateMult:   1.0,
		RewardMult:      1.0,
		SpeedMult:       1.0,
		RareChanceBonus: 0.0,
		GoldFishBonus:   0.0,
		BossChanceBonus: 0.0,
		FogEffect:       false,
		Weight:          40,
	},
	WeatherRain: {
		Type:        WeatherRain,
		Name:        "下雨",
		Icon:        "🌧️",
		Description: "雨水帶來稀有魚群！稀有目標出現率 +20%",
		Duration:    4 * time.Minute,
		SpawnRateMult:   1.1,
		RewardMult:      1.0,
		SpeedMult:       1.0,
		RareChanceBonus: 0.20,
		GoldFishBonus:   0.0,
		BossChanceBonus: 0.0,
		FogEffect:       false,
		Weight:          25,
	},
	WeatherStorm: {
		Type:        WeatherStorm,
		Name:        "暴風雨",
		Icon:        "⛈️",
		Description: "狂風暴雨！目標移動加速 +30%，獎勵 ×1.5",
		Duration:    3 * time.Minute,
		SpawnRateMult:   1.2,
		RewardMult:      1.5,
		SpeedMult:       1.3,
		RareChanceBonus: 0.0,
		GoldFishBonus:   0.0,
		BossChanceBonus: 0.0,
		FogEffect:       false,
		Weight:          15,
	},
	WeatherFog: {
		Type:        WeatherFog,
		Name:        "濃霧",
		Icon:        "🌫️",
		Description: "神秘濃霧！目標若隱若現，擊破獎勵 ×2.0",
		Duration:    3 * time.Minute,
		SpawnRateMult:   0.9,
		RewardMult:      2.0,
		SpeedMult:       0.8,
		RareChanceBonus: 0.0,
		GoldFishBonus:   0.0,
		BossChanceBonus: 0.0,
		FogEffect:       true,
		Weight:          10,
	},
	WeatherSunshine: {
		Type:        WeatherSunshine,
		Name:        "豔陽高照",
		Icon:        "🌞",
		Description: "金光閃閃！金幣魚出現率 +50%，獎勵 ×1.2",
		Duration:    4 * time.Minute,
		SpawnRateMult:   1.0,
		RewardMult:      1.2,
		SpeedMult:       1.0,
		RareChanceBonus: 0.0,
		GoldFishBonus:   0.50,
		BossChanceBonus: 0.0,
		FogEffect:       false,
		Weight:          7,
	},
	WeatherBlizzard: {
		Type:        WeatherBlizzard,
		Name:        "暴雪",
		Icon:        "❄️",
		Description: "冰封海域！BOSS 出現機率 +30%，BOSS 獎勵 ×2.0",
		Duration:    3 * time.Minute,
		SpawnRateMult:   0.8,
		RewardMult:      1.0,
		SpeedMult:       0.9,
		RareChanceBonus: 0.0,
		GoldFishBonus:   0.0,
		BossChanceBonus: 0.30,
		FogEffect:       false,
		Weight:          3,
	},
}

// weatherOrder 天氣輪換順序（用於加權隨機）
var weatherOrder = []WeatherType{
	WeatherClear, WeatherRain, WeatherStorm, WeatherFog, WeatherSunshine, WeatherBlizzard,
}

// Snapshot 天氣快照（用於廣播）
type Snapshot struct {
	Type            WeatherType `json:"type"`
	Name            string      `json:"name"`
	Icon            string      `json:"icon"`
	Description     string      `json:"description"`
	RemainingSeconds int        `json:"remaining_seconds"`
	SpawnRateMult   float64     `json:"spawn_rate_mult"`
	RewardMult      float64     `json:"reward_mult"`
	SpeedMult       float64     `json:"speed_mult"`
	RareChanceBonus float64     `json:"rare_chance_bonus"`
	GoldFishBonus   float64     `json:"gold_fish_bonus"`
	BossChanceBonus float64     `json:"boss_chance_bonus"`
	FogEffect       bool        `json:"fog_effect"`
	IsNew           bool        `json:"is_new"` // 是否剛切換（用於 Client 端顯示通知）
}

// Manager 天氣管理器
type Manager struct {
	mu          sync.RWMutex
	current     WeatherType
	startedAt   time.Time
	rng         *rand.Rand
}

// New 建立天氣管理器（從晴天開始）
func New() *Manager {
	return &Manager{
		current:   WeatherClear,
		startedAt: time.Now(),
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GetSnapshot 取得當前天氣快照
func (m *Manager) GetSnapshot(isNew bool) Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	def := WeatherDefs[m.current]
	elapsed := time.Since(m.startedAt)
	remaining := def.Duration - elapsed
	if remaining < 0 {
		remaining = 0
	}

	return Snapshot{
		Type:             m.current,
		Name:             def.Name,
		Icon:             def.Icon,
		Description:      def.Description,
		RemainingSeconds: int(remaining.Seconds()),
		SpawnRateMult:    def.SpawnRateMult,
		RewardMult:       def.RewardMult,
		SpeedMult:        def.SpeedMult,
		RareChanceBonus:  def.RareChanceBonus,
		GoldFishBonus:    def.GoldFishBonus,
		BossChanceBonus:  def.BossChanceBonus,
		FogEffect:        def.FogEffect,
		IsNew:            isNew,
	}
}

// CheckAndRotate 檢查是否需要切換天氣，回傳 (changed bool, newSnap Snapshot)
func (m *Manager) CheckAndRotate() (bool, Snapshot) {
	m.mu.Lock()
	defer m.mu.Unlock()

	def := WeatherDefs[m.current]
	if time.Since(m.startedAt) < def.Duration {
		return false, Snapshot{}
	}

	// 加權隨機選擇下一個天氣（不重複選同一個）
	next := m.pickNext()
	m.current = next
	m.startedAt = time.Now()

	newDef := WeatherDefs[next]
	remaining := newDef.Duration

	snap := Snapshot{
		Type:             next,
		Name:             newDef.Name,
		Icon:             newDef.Icon,
		Description:      newDef.Description,
		RemainingSeconds: int(remaining.Seconds()),
		SpawnRateMult:    newDef.SpawnRateMult,
		RewardMult:       newDef.RewardMult,
		SpeedMult:        newDef.SpeedMult,
		RareChanceBonus:  newDef.RareChanceBonus,
		GoldFishBonus:    newDef.GoldFishBonus,
		BossChanceBonus:  newDef.BossChanceBonus,
		FogEffect:        newDef.FogEffect,
		IsNew:            true,
	}
	return true, snap
}

// pickNext 加權隨機選擇下一個天氣（不重複選當前天氣）
func (m *Manager) pickNext() WeatherType {
	totalWeight := 0
	for _, wt := range weatherOrder {
		if wt == m.current {
			continue
		}
		totalWeight += WeatherDefs[wt].Weight
	}

	r := m.rng.Intn(totalWeight)
	cumulative := 0
	for _, wt := range weatherOrder {
		if wt == m.current {
			continue
		}
		cumulative += WeatherDefs[wt].Weight
		if r < cumulative {
			return wt
		}
	}
	return WeatherClear
}

// GetRewardMult 取得當前天氣獎勵倍率
func (m *Manager) GetRewardMult() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return WeatherDefs[m.current].RewardMult
}

// GetSpeedMult 取得當前天氣速度倍率
func (m *Manager) GetSpeedMult() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return WeatherDefs[m.current].SpeedMult
}

// GetRareChanceBonus 取得稀有目標出現機率加成
func (m *Manager) GetRareChanceBonus() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return WeatherDefs[m.current].RareChanceBonus
}

// GetGoldFishBonus 取得金幣魚出現機率加成
func (m *Manager) GetGoldFishBonus() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return WeatherDefs[m.current].GoldFishBonus
}

// GetBossChanceBonus 取得 BOSS 出現機率加成
func (m *Manager) GetBossChanceBonus() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return WeatherDefs[m.current].BossChanceBonus
}

// GetFogEffect 取得是否有濃霧效果
func (m *Manager) GetFogEffect() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return WeatherDefs[m.current].FogEffect
}

// GetCurrent 取得當前天氣類型
func (m *Manager) GetCurrent() WeatherType {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.current
}
