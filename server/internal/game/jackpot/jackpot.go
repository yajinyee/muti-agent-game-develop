// Package jackpot 實作 Progressive Jackpot 系統
// 四個等級：Mini / Minor / Major / Grand（DAY-095 升級）
// 每次攻擊抽取 0.5% 進入 Jackpot 池，達到門檻時觸發
package jackpot

import (
	"math/rand"
	"sync"
	"time"
)

// Level Jackpot 等級
type Level string

const (
	LevelMini  Level = "mini"  // 小獎：門檻 300x，觸發機率 1/300
	LevelMinor Level = "minor" // 次獎：門檻 1000x，觸發機率 1/800（DAY-095 新增）
	LevelMajor Level = "major" // 大獎：門檻 3000x，觸發機率 1/2000
	LevelGrand Level = "grand" // 超大獎：門檻 15000x，觸發機率 1/8000
)

// ContributionRate 每次攻擊抽取比例（0.5%）
const ContributionRate = 0.005

// JackpotWin 中獎記錄
type JackpotWin struct {
	Level    Level
	Amount   int
	WinnerID string
	WonAt    time.Time
}

// Pool 單一 Jackpot 池
type Pool struct {
	Level       Level
	Current     int // 當前累積金額（以 bet_cost 為單位）
	Threshold   int // 觸發門檻
	TriggerOdds int // 觸發機率分母（達到門檻後每次攻擊的觸發機率 = 1/TriggerOdds）
	BaseAmount  int // 基礎金額（重置後的起始值）
}

// Manager Jackpot 管理器
type Manager struct {
	mu    sync.RWMutex
	pools map[Level]*Pool
	rng   *rand.Rand
}

// NewManager 建立 Jackpot 管理器（四層）
func NewManager() *Manager {
	return &Manager{
		pools: map[Level]*Pool{
			LevelMini: {
				Level:       LevelMini,
				Current:     80,  // 起始 80x
				Threshold:   300, // 門檻 300x
				TriggerOdds: 300, // 1/300 機率
				BaseAmount:  80,
			},
			LevelMinor: {
				Level:       LevelMinor,
				Current:     200,  // 起始 200x
				Threshold:   1000, // 門檻 1000x
				TriggerOdds: 800,  // 1/800 機率
				BaseAmount:  200,
			},
			LevelMajor: {
				Level:       LevelMajor,
				Current:     500,  // 起始 500x
				Threshold:   3000, // 門檻 3000x
				TriggerOdds: 2000, // 1/2000 機率
				BaseAmount:  500,
			},
			LevelGrand: {
				Level:       LevelGrand,
				Current:     2000,  // 起始 2000x
				Threshold:   15000, // 門檻 15000x
				TriggerOdds: 8000,  // 1/8000 機率（非常稀有）
				BaseAmount:  2000,
			},
		},
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Contribute 每次攻擊時貢獻 Jackpot 池
// betCost: 本次攻擊的 bet_cost
// 回傳：是否有 Jackpot 觸發，觸發的等級和金額
func (m *Manager) Contribute(betCost int, playerID string) *JackpotWin {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 計算貢獻金額（0.5% of betCost，最少 4，確保四個池子都能分到至少 1）
	contribution := int(float64(betCost) * ContributionRate)
	if contribution < 4 {
		contribution = 4
	}

	// 四層分配：Mini 50%，Minor 25%，Major 15%，Grand 10%
	miniShare := int(float64(contribution) * 0.50)
	if miniShare < 1 {
		miniShare = 1
	}
	minorShare := int(float64(contribution) * 0.25)
	if minorShare < 1 {
		minorShare = 1
	}
	majorShare := int(float64(contribution) * 0.15)
	if majorShare < 1 {
		majorShare = 1
	}
	grandShare := contribution - miniShare - minorShare - majorShare
	if grandShare < 1 {
		grandShare = 1
	}

	m.pools[LevelMini].Current += miniShare
	m.pools[LevelMinor].Current += minorShare
	m.pools[LevelMajor].Current += majorShare
	m.pools[LevelGrand].Current += grandShare

	// 檢查觸發（從 Grand 開始，優先觸發大獎）
	for _, level := range []Level{LevelGrand, LevelMajor, LevelMinor, LevelMini} {
		pool := m.pools[level]
		if pool.Current >= pool.Threshold {
			// 達到門檻，以 1/TriggerOdds 機率觸發
			if m.rng.Intn(pool.TriggerOdds) == 0 {
				win := &JackpotWin{
					Level:    level,
					Amount:   pool.Current,
					WinnerID: playerID,
					WonAt:    time.Now(),
				}
				// 重置池子到基礎金額
				pool.Current = pool.BaseAmount
				return win
			}
		}
	}

	return nil
}

// GetSnapshot 取得當前 Jackpot 池快照（用於廣播）
func (m *Manager) GetSnapshot() map[Level]int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[Level]int{
		LevelMini:  m.pools[LevelMini].Current,
		LevelMinor: m.pools[LevelMinor].Current,
		LevelMajor: m.pools[LevelMajor].Current,
		LevelGrand: m.pools[LevelGrand].Current,
	}
}

// ForceWin 強制觸發指定等級的 Jackpot（測試用）
func (m *Manager) ForceWin(level Level, playerID string) *JackpotWin {
	m.mu.Lock()
	defer m.mu.Unlock()

	pool, ok := m.pools[level]
	if !ok {
		return nil
	}

	win := &JackpotWin{
		Level:    level,
		Amount:   pool.Current,
		WinnerID: playerID,
		WonAt:    time.Now(),
	}
	pool.Current = pool.BaseAmount
	return win
}

// PoolState Jackpot 池狀態快照（用於持久化）
type PoolState struct {
	Mini  int `json:"mini"`
	Minor int `json:"minor"` // DAY-095 新增
	Major int `json:"major"`
	Grand int `json:"grand"`
}

// SaveState 取得當前池狀態（用於持久化到 Redis）
func (m *Manager) SaveState() PoolState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return PoolState{
		Mini:  m.pools[LevelMini].Current,
		Minor: m.pools[LevelMinor].Current,
		Major: m.pools[LevelMajor].Current,
		Grand: m.pools[LevelGrand].Current,
	}
}

// LoadState 從持久化狀態恢復池金額（Server 重啟後呼叫）
// 只恢復大於基礎金額的值，防止異常數據
func (m *Manager) LoadState(state PoolState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if state.Mini > m.pools[LevelMini].BaseAmount {
		m.pools[LevelMini].Current = state.Mini
	}
	if state.Minor > m.pools[LevelMinor].BaseAmount {
		m.pools[LevelMinor].Current = state.Minor
	}
	if state.Major > m.pools[LevelMajor].BaseAmount {
		m.pools[LevelMajor].Current = state.Major
	}
	if state.Grand > m.pools[LevelGrand].BaseAmount {
		m.pools[LevelGrand].Current = state.Grand
	}
}

// GetLevelInfo 取得等級顯示資訊（名稱、顏色、圖示）
func GetLevelInfo(level Level) (name, color, icon string) {
	switch level {
	case LevelMini:
		return "MINI", "#C0C0C0", "🥈" // 銀色
	case LevelMinor:
		return "MINOR", "#FFD700", "🥇" // 金色
	case LevelMajor:
		return "MAJOR", "#FF6B35", "🔥" // 橙紅
	case LevelGrand:
		return "GRAND", "#FF0080", "👑" // 粉紅/紫
	default:
		return "JACKPOT", "#FFFFFF", "🎰"
	}
}
