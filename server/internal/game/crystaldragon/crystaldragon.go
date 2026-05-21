// Package crystaldragon — 水晶龍收集大獎系統（DAY-153）
// 業界依據：jiligames.com JILI Flying Dragon 2026「collect crystals to get the grand prize!
// Kill the Underworld Dragon and win the prize!」
// 設計：T117 水晶龍擊破後掉落水晶碎片，全服玩家共同收集水晶，
// 達到目標數量後觸發「地獄龍大獎」，全服廣播並發放大獎
// 社交設計：全服合作收集，讓所有玩家都有參與感，增加留存
package crystaldragon

import (
	"sync"
	"time"
)

const (
	// CrystalGoal 全服水晶收集目標（達到後觸發地獄龍大獎）
	CrystalGoal = 50

	// CrystalPerKill 每次擊破水晶龍掉落的水晶數量
	CrystalPerKill = 5

	// CrystalDecayInterval 水晶衰減間隔（每 30 秒減少 1 個，防止永遠不觸發）
	CrystalDecayInterval = 30 * time.Second

	// CooldownDuration 地獄龍大獎觸發後的冷卻時間
	CooldownDuration = 120 * time.Second

	// HellDragonBaseRewardMult 地獄龍大獎基礎倍率（每個水晶貢獻者按比例分配）
	HellDragonBaseRewardMult = 200.0

	// MaxContributors 最多記錄的貢獻者數量
	MaxContributors = 20
)

// Contributor 水晶貢獻者記錄
type Contributor struct {
	PlayerID    string
	PlayerName  string
	Crystals    int     // 貢獻的水晶數量
	BetLevel    int     // 貢獻時的 betLevel（用於計算獎勵）
}

// Manager 水晶龍收集大獎管理器
type Manager struct {
	mu sync.Mutex

	// 全服水晶收集進度
	totalCrystals int
	contributors  map[string]*Contributor // playerID -> Contributor

	// 冷卻狀態
	lastTriggerAt time.Time
	lastDecayAt   time.Time

	// 統計
	totalTriggered int // 歷史觸發次數
}

// New 建立新的水晶龍管理器
func New() *Manager {
	return &Manager{
		contributors: make(map[string]*Contributor),
		lastDecayAt:  time.Now(),
	}
}

// AddCrystals 增加水晶（玩家擊破水晶龍後呼叫）
// 回傳：新的水晶總數、是否達到目標（觸發地獄龍大獎）
func (m *Manager) AddCrystals(playerID, playerName string, betLevel, amount int) (newTotal int, triggered bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 冷卻中不接受新水晶
	if !m.lastTriggerAt.IsZero() && time.Since(m.lastTriggerAt) < CooldownDuration {
		return m.totalCrystals, false
	}

	// 記錄貢獻者
	if c, ok := m.contributors[playerID]; ok {
		c.Crystals += amount
		c.BetLevel = betLevel // 更新為最新 betLevel
	} else {
		if len(m.contributors) < MaxContributors {
			m.contributors[playerID] = &Contributor{
				PlayerID:   playerID,
				PlayerName: playerName,
				Crystals:   amount,
				BetLevel:   betLevel,
			}
		}
	}

	m.totalCrystals += amount
	if m.totalCrystals > CrystalGoal {
		m.totalCrystals = CrystalGoal
	}

	newTotal = m.totalCrystals
	triggered = m.totalCrystals >= CrystalGoal
	return
}

// TriggerHellDragon 觸發地獄龍大獎（達到目標後呼叫）
// 回傳：貢獻者列表（用於發放獎勵）
func (m *Manager) TriggerHellDragon() []*Contributor {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.totalCrystals < CrystalGoal {
		return nil
	}

	// 收集貢獻者列表
	contributors := make([]*Contributor, 0, len(m.contributors))
	for _, c := range m.contributors {
		contributors = append(contributors, &Contributor{
			PlayerID:   c.PlayerID,
			PlayerName: c.PlayerName,
			Crystals:   c.Crystals,
			BetLevel:   c.BetLevel,
		})
	}

	// 重置狀態
	m.totalCrystals = 0
	m.contributors = make(map[string]*Contributor)
	m.lastTriggerAt = time.Now()
	m.totalTriggered++

	return contributors
}

// CheckDecay 水晶衰減檢查（每 30 秒減少 1 個，防止永遠不觸發）
// 回傳：是否有衰減發生
func (m *Manager) CheckDecay() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.totalCrystals <= 0 {
		return false
	}
	if time.Since(m.lastDecayAt) < CrystalDecayInterval {
		return false
	}

	m.totalCrystals--
	m.lastDecayAt = time.Now()
	return true
}

// GetSnapshot 取得當前狀態快照（thread-safe）
type Snapshot struct {
	TotalCrystals  int
	Goal           int
	Progress       float64 // 0.0 - 1.0
	CooldownSecs   int     // 冷卻剩餘秒數（0 = 可觸發）
	TotalTriggered int
}

func (m *Manager) GetSnapshot() Snapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	cooldown := 0
	if !m.lastTriggerAt.IsZero() {
		remaining := CooldownDuration - time.Since(m.lastTriggerAt)
		if remaining > 0 {
			cooldown = int(remaining.Seconds())
		}
	}

	progress := float64(m.totalCrystals) / float64(CrystalGoal)
	if progress > 1.0 {
		progress = 1.0
	}

	return Snapshot{
		TotalCrystals:  m.totalCrystals,
		Goal:           CrystalGoal,
		Progress:       progress,
		CooldownSecs:   cooldown,
		TotalTriggered: m.totalTriggered,
	}
}

// IsOnCooldown 是否在冷卻中
func (m *Manager) IsOnCooldown() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return !m.lastTriggerAt.IsZero() && time.Since(m.lastTriggerAt) < CooldownDuration
}

// CalcReward 計算貢獻者的地獄龍大獎獎勵
// 獎勵 = 貢獻比例 × HellDragonBaseRewardMult × betLevel
func CalcReward(contributor *Contributor, totalCrystals int) int {
	if totalCrystals <= 0 || contributor.Crystals <= 0 {
		return 0
	}
	ratio := float64(contributor.Crystals) / float64(totalCrystals)
	reward := int(ratio * HellDragonBaseRewardMult * float64(contributor.BetLevel))
	if reward < contributor.BetLevel {
		reward = contributor.BetLevel // 最少 1x betLevel
	}
	return reward
}
