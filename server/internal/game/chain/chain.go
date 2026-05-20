// Package chain 連鎖爆炸系統（DAY-088）
// 業界依據：Avalanche/Cascading Reels 是 2026 年最熱門的留存機制
// 擊破目標後有機率觸發連鎖，周圍目標同時爆炸，獎勵疊加，製造爽感
package chain

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

// ChainLevel 連鎖等級
type ChainLevel int

const (
	ChainNone   ChainLevel = 0 // 無連鎖
	ChainSmall  ChainLevel = 1 // 小連鎖（1個額外目標）
	ChainMedium ChainLevel = 2 // 中連鎖（2-3個額外目標）
	ChainBig    ChainLevel = 3 // 大連鎖（4-6個額外目標）
	ChainMega   ChainLevel = 4 // 超級連鎖（7-10個額外目標）
)

// ChainResult 連鎖結果
type ChainResult struct {
	Level       ChainLevel // 連鎖等級
	TargetIDs   []string   // 被連鎖擊破的目標 ID 列表
	BonusMult   float64    // 連鎖獎勵倍率加成（1.0 = 無加成）
	LevelName   string     // 等級名稱（用於 Client 顯示）
	LevelColor  string     // 顯示顏色（hex）
}

// ChainConfig 連鎖觸發設定
type ChainConfig struct {
	// 基礎觸發機率（依觸發目標倍率調整）
	// 倍率 2x: 5%, 5x: 10%, 10x: 15%, 20x+: 20%, 50x+: 30%
	BaseChance float64
	// 連鎖範圍（像素）
	Radius float64
	// 最大連鎖深度（防止無限連鎖）
	MaxDepth int
}

// DefaultConfig 預設連鎖設定
var DefaultConfig = ChainConfig{
	BaseChance: 0.05, // 基礎 5%
	Radius:     200.0,
	MaxDepth:   2, // 最多 2 層連鎖（防止 RTP 爆炸）
}

// TargetInfo 目標物資訊（用於連鎖計算）
type TargetInfo struct {
	ID         string
	X, Y       float64
	Multiplier float64
	DefID      string
}

// Manager 連鎖爆炸管理器
type Manager struct {
	mu  sync.Mutex
	rng *rand.Rand
	cfg ChainConfig
}

// New 建立連鎖爆炸管理器
func New(cfg ChainConfig) *Manager {
	return &Manager{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
		cfg: cfg,
	}
}

// NewDefault 使用預設設定建立管理器
func NewDefault() *Manager {
	return New(DefaultConfig)
}

// CalcChance 計算連鎖觸發機率（依觸發目標倍率）
func CalcChance(multiplier float64) float64 {
	switch {
	case multiplier >= 50:
		return 0.30 // 50x+ → 30%
	case multiplier >= 20:
		return 0.20 // 20x+ → 20%
	case multiplier >= 10:
		return 0.15 // 10x+ → 15%
	case multiplier >= 5:
		return 0.10 // 5x+ → 10%
	default:
		return 0.05 // 2x → 5%
	}
}

// TryChain 嘗試觸發連鎖爆炸
// triggerTarget: 觸發連鎖的目標
// allTargets: 場上所有目標（用於找周圍目標）
// depth: 當前連鎖深度（防止無限遞迴）
// 回傳 ChainResult（Level=ChainNone 表示未觸發）
func (m *Manager) TryChain(
	triggerTarget TargetInfo,
	allTargets []TargetInfo,
	depth int,
) ChainResult {
	if depth >= m.cfg.MaxDepth {
		return ChainResult{Level: ChainNone}
	}

	m.mu.Lock()
	chance := CalcChance(triggerTarget.Multiplier)
	roll := m.rng.Float64()
	m.mu.Unlock()

	if roll >= chance {
		return ChainResult{Level: ChainNone}
	}

	// 找周圍目標（排除觸發目標本身，排除 BOSS）
	nearby := m.findNearby(triggerTarget, allTargets)
	if len(nearby) == 0 {
		return ChainResult{Level: ChainNone}
	}

	// 依連鎖等級決定最多擊破幾個
	maxKills := m.calcMaxKills(triggerTarget.Multiplier)
	if maxKills > len(nearby) {
		maxKills = len(nearby)
	}

	// 隨機選取目標
	m.mu.Lock()
	m.rng.Shuffle(len(nearby), func(i, j int) {
		nearby[i], nearby[j] = nearby[j], nearby[i]
	})
	m.mu.Unlock()

	selected := nearby[:maxKills]
	ids := make([]string, len(selected))
	for i, t := range selected {
		ids[i] = t.ID
	}

	level := m.calcLevel(maxKills)
	bonusMult := m.calcBonusMult(level)
	levelName, levelColor := levelInfo(level)

	return ChainResult{
		Level:      level,
		TargetIDs:  ids,
		BonusMult:  bonusMult,
		LevelName:  levelName,
		LevelColor: levelColor,
	}
}

// findNearby 找觸發目標周圍的目標（排除 BOSS 和觸發目標本身）
func (m *Manager) findNearby(trigger TargetInfo, all []TargetInfo) []TargetInfo {
	var result []TargetInfo
	for _, t := range all {
		if t.ID == trigger.ID {
			continue
		}
		// 排除 BOSS（DefID 以 B 開頭）
		if len(t.DefID) > 0 && t.DefID[0] == 'B' {
			continue
		}
		dist := math.Sqrt(math.Pow(t.X-trigger.X, 2) + math.Pow(t.Y-trigger.Y, 2))
		if dist <= m.cfg.Radius {
			result = append(result, t)
		}
	}
	return result
}

// calcMaxKills 依倍率計算最多連鎖擊破數
func (m *Manager) calcMaxKills(multiplier float64) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	switch {
	case multiplier >= 50:
		return 7 + m.rng.Intn(4) // 7-10
	case multiplier >= 20:
		return 4 + m.rng.Intn(3) // 4-6
	case multiplier >= 10:
		return 2 + m.rng.Intn(2) // 2-3
	default:
		return 1 // 1
	}
}

// calcLevel 依擊破數計算連鎖等級
func calcLevel(count int) ChainLevel {
	switch {
	case count >= 7:
		return ChainMega
	case count >= 4:
		return ChainBig
	case count >= 2:
		return ChainMedium
	default:
		return ChainSmall
	}
}

func (m *Manager) calcLevel(count int) ChainLevel {
	return calcLevel(count)
}

// calcBonusMult 依連鎖等級計算獎勵倍率加成
func (m *Manager) calcBonusMult(level ChainLevel) float64 {
	switch level {
	case ChainMega:
		return 2.0 // 超級連鎖：獎勵 ×2.0
	case ChainBig:
		return 1.5 // 大連鎖：獎勵 ×1.5
	case ChainMedium:
		return 1.2 // 中連鎖：獎勵 ×1.2
	default:
		return 1.0 // 小連鎖：無額外倍率
	}
}

// levelInfo 回傳連鎖等級名稱和顏色
func levelInfo(level ChainLevel) (string, string) {
	switch level {
	case ChainMega:
		return "超級連鎖！", "#FF4500" // 橙紅
	case ChainBig:
		return "大連鎖！", "#FFD700" // 金色
	case ChainMedium:
		return "連鎖！", "#00BFFF" // 天藍
	default:
		return "小連鎖", "#FFFFFF" // 白色
	}
}
