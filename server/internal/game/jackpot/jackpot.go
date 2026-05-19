// Package jackpot 實作 Progressive Jackpot 系統
// 三個等級：Mini / Major / Grand
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
	LevelMini  Level = "mini"  // 小獎：門檻 500x，觸發機率 1/200
	LevelMajor Level = "major" // 大獎：門檻 2000x，觸發機率 1/1000
	LevelGrand Level = "grand" // 超大獎：門檻 10000x，觸發機率 1/5000
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
	Current     int     // 當前累積金額（以 bet_cost 為單位）
	Threshold   int     // 觸發門檻
	TriggerOdds int     // 觸發機率分母（達到門檻後每次攻擊的觸發機率 = 1/TriggerOdds）
	BaseAmount  int     // 基礎金額（重置後的起始值）
}

// Manager Jackpot 管理器
type Manager struct {
	mu    sync.RWMutex
	pools map[Level]*Pool
	rng   *rand.Rand
}

// NewManager 建立 Jackpot 管理器
func NewManager() *Manager {
	return &Manager{
		pools: map[Level]*Pool{
			LevelMini: {
				Level:       LevelMini,
				Current:     100,   // 起始 100x（需要累積到 500x 才能觸發）
				Threshold:   500,   // 門檻 500x
				TriggerOdds: 500,   // 1/500 機率（達到門檻後，平均 500 次攻擊觸發一次）
				BaseAmount:  100,   // 重置後回到 100x
			},
			LevelMajor: {
				Level:       LevelMajor,
				Current:     500,   // 起始 500x
				Threshold:   2000,  // 門檻 2000x
				TriggerOdds: 2000,  // 1/2000 機率
				BaseAmount:  500,
			},
			LevelGrand: {
				Level:       LevelGrand,
				Current:     2000,  // 起始 2000x
				Threshold:   10000, // 門檻 10000x
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

	// 計算貢獻金額（0.5% of betCost，最少 3，確保三個池子都能分到至少 1）
	contribution := int(float64(betCost) * ContributionRate)
	if contribution < 3 {
		contribution = 3
	}

	// 依序從 Mini → Major → Grand 分配貢獻
	// Mini 拿 60%，Major 拿 30%，Grand 拿 10%（最少各 1）
	miniShare := int(float64(contribution) * 0.6)
	if miniShare < 1 {
		miniShare = 1
	}
	majorShare := int(float64(contribution) * 0.3)
	if majorShare < 1 {
		majorShare = 1
	}
	grandShare := contribution - miniShare - majorShare
	if grandShare < 1 {
		grandShare = 1
	}

	m.pools[LevelMini].Current += miniShare
	m.pools[LevelMajor].Current += majorShare
	m.pools[LevelGrand].Current += grandShare

	// 檢查觸發（從 Grand 開始，優先觸發大獎）
	for _, level := range []Level{LevelGrand, LevelMajor, LevelMini} {
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
