// Package raidboss — Co-op Boss Raid 系統（DAY-115）
// 全服玩家合作討伐超強 BOSS，依貢獻度分配獎勵池
package raidboss

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// RaidState 討伐狀態
type RaidState string

const (
	RaidStateIdle    RaidState = "idle"    // 無討伐
	RaidStateWarning RaidState = "warning" // 警告中（30 秒）
	RaidStateActive  RaidState = "active"  // 討伐進行中
	RaidStateResult  RaidState = "result"  // 結算中
)

// ContributorEntry 貢獻者記錄
type ContributorEntry struct {
	PlayerID    string
	DisplayName string
	Damage      int     // 累計傷害
	Reward      int     // 分配到的獎勵
	Rank        int     // 貢獻排名
}

// RaidSnapshot 討伐快照（用於廣播）
type RaidSnapshot struct {
	RaidID       string
	State        RaidState
	BossName     string
	HP           int
	MaxHP        int
	RewardPool   int
	Contributors []*ContributorEntry // 依傷害降序
	StartedAt    time.Time
	EndsAt       time.Time
	TimeLeft     float64 // 秒
}

// Manager Co-op Boss Raid 管理器
type Manager struct {
	mu sync.RWMutex

	raidID       string
	state        RaidState
	bossName     string
	hp           int
	maxHP        int
	rewardPool   int
	contributors map[string]*ContributorEntry
	startedAt    time.Time
	endsAt       time.Time
	lastRaidDate string // YYYY-MM-DD，防止同一天重複觸發
}

// RaidDuration 討伐持續時間（5 分鐘）
const RaidDuration = 5 * time.Minute

// New 建立新的 Raid 管理器
func New() *Manager {
	return &Manager{
		state:        RaidStateIdle,
		contributors: make(map[string]*ContributorEntry),
	}
}

// CanTrigger 是否可以觸發討伐（idle 狀態且今日未觸發）
func (m *Manager) CanTrigger(todayDate string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state == RaidStateIdle && m.lastRaidDate != todayDate
}

// StartWarning 進入警告狀態
func (m *Manager) StartWarning() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = RaidStateWarning
	m.raidID = uuid.New().String()
	return m.raidID
}

// StartRaid 開始討伐
// bossName: BOSS 名稱
// hp: BOSS 血量
// rewardPool: 獎勵池總金幣
func (m *Manager) StartRaid(bossName string, hp, rewardPool int, todayDate string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = RaidStateActive
	m.bossName = bossName
	m.hp = hp
	m.maxHP = hp
	m.rewardPool = rewardPool
	m.contributors = make(map[string]*ContributorEntry)
	m.startedAt = time.Now()
	m.endsAt = time.Now().Add(RaidDuration)
	m.lastRaidDate = todayDate
}

// RecordDamage 記錄玩家傷害貢獻，回傳 (newHP, isKilled)
func (m *Manager) RecordDamage(playerID, displayName string, damage int) (int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.state != RaidStateActive {
		return m.hp, false
	}

	// 記錄貢獻
	entry, ok := m.contributors[playerID]
	if !ok {
		entry = &ContributorEntry{
			PlayerID:    playerID,
			DisplayName: displayName,
		}
		m.contributors[playerID] = entry
	}
	entry.Damage += damage

	// 扣血
	m.hp -= damage
	if m.hp < 0 {
		m.hp = 0
	}

	if m.hp == 0 {
		m.state = RaidStateResult
		m.distributeRewards()
		return 0, true
	}
	return m.hp, false
}

// CheckTimeout 檢查是否超時，回傳是否超時
func (m *Manager) CheckTimeout() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.state != RaidStateActive {
		return false
	}
	if time.Now().After(m.endsAt) {
		m.state = RaidStateResult
		m.distributeRewards()
		return true
	}
	return false
}

// distributeRewards 依貢獻度分配獎勵（必須在 mu.Lock 下呼叫）
func (m *Manager) distributeRewards() {
	if len(m.contributors) == 0 {
		return
	}

	// 計算總傷害
	totalDamage := 0
	for _, e := range m.contributors {
		totalDamage += e.Damage
	}
	if totalDamage == 0 {
		return
	}

	// 依傷害比例分配獎勵
	distributed := 0
	entries := m.sortedContributors()
	for i, e := range entries {
		if i == len(entries)-1 {
			// 最後一名拿剩餘（避免浮點誤差）
			e.Reward = m.rewardPool - distributed
		} else {
			e.Reward = int(float64(m.rewardPool) * float64(e.Damage) / float64(totalDamage))
			distributed += e.Reward
		}
		e.Rank = i + 1
	}
}

// sortedContributors 依傷害降序排列（必須在 mu.Lock 下呼叫）
func (m *Manager) sortedContributors() []*ContributorEntry {
	entries := make([]*ContributorEntry, 0, len(m.contributors))
	for _, e := range m.contributors {
		entries = append(entries, e)
	}
	// 簡單插入排序（貢獻者通常不多）
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].Damage > entries[j-1].Damage; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
	return entries
}

// GetSnapshot 取得當前快照
func (m *Manager) GetSnapshot() *RaidSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	timeLeft := 0.0
	if m.state == RaidStateActive {
		timeLeft = time.Until(m.endsAt).Seconds()
		if timeLeft < 0 {
			timeLeft = 0
		}
	}

	contributors := m.sortedContributors()

	return &RaidSnapshot{
		RaidID:       m.raidID,
		State:        m.state,
		BossName:     m.bossName,
		HP:           m.hp,
		MaxHP:        m.maxHP,
		RewardPool:   m.rewardPool,
		Contributors: contributors,
		StartedAt:    m.startedAt,
		EndsAt:       m.endsAt,
		TimeLeft:     timeLeft,
	}
}

// GetState 取得當前狀態
func (m *Manager) GetState() RaidState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

// GetContributorReward 取得特定玩家的獎勵（結算後呼叫）
func (m *Manager) GetContributorReward(playerID string) (int, int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.contributors[playerID]
	if !ok {
		return 0, 0, false
	}
	return e.Reward, e.Rank, true
}

// Reset 重置為 idle 狀態
func (m *Manager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = RaidStateIdle
	m.hp = 0
	m.maxHP = 0
	m.contributors = make(map[string]*ContributorEntry)
}

// IsActive 是否正在進行討伐
func (m *Manager) IsActive() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state == RaidStateActive
}
