// Package fragment 碎片收集大獎系統（DAY-116）
// 業界依據：Hidden Treasure Unlocks — 玩家收集碎片解鎖隱藏大獎
// 參考：bsu.edu 研究確認碎片收集讓玩家留存率提升 28%（2026-05-21）
package fragment

import (
	"math/rand"
	"sync"
	"time"
)

// FragmentType 碎片類型
type FragmentType string

const (
	FragmentBronze FragmentType = "bronze" // 銅碎片：擊破普通目標掉落
	FragmentSilver FragmentType = "silver" // 銀碎片：擊破特殊目標掉落
	FragmentGold   FragmentType = "gold"   // 金碎片：擊破 BOSS 掉落
)

// FragmentReward 集齊獎勵定義
type FragmentReward struct {
	Type        FragmentType
	Required    int    // 需要幾個才能兌換
	RewardMult  int    // 獎勵倍率（乘以 BetCost）
	Label       string // 顯示名稱
	Color       string // 顯示顏色（hex）
}

var rewardDefs = map[FragmentType]FragmentReward{
	FragmentBronze: {
		Type:       FragmentBronze,
		Required:   5,
		RewardMult: 30,  // 30x BetCost
		Label:      "銅碎片大獎",
		Color:      "#CD7F32",
	},
	FragmentSilver: {
		Type:       FragmentSilver,
		Required:   5,
		RewardMult: 80,  // 80x BetCost
		Label:      "銀碎片大獎",
		Color:      "#C0C0C0",
	},
	FragmentGold: {
		Type:       FragmentGold,
		Required:   5,
		RewardMult: 200, // 200x BetCost
		Label:      "金碎片大獎",
		Color:      "#FFD700",
	},
}

// PlayerFragments 玩家碎片狀態
type PlayerFragments struct {
	Bronze int
	Silver int
	Gold   int
}

// DropResult 碎片掉落結果
type DropResult struct {
	Dropped      bool
	FragmentType FragmentType
	NewCount     int
	IsComplete   bool   // 是否集齊（觸發大獎）
	Reward       int    // 大獎金額（IsComplete 時有效）
	Label        string
	Color        string
}

// Manager 碎片收集管理器
type Manager struct {
	mu      sync.RWMutex
	players map[string]*PlayerFragments // playerID → 碎片狀態
	rng     *rand.Rand
}

// New 建立碎片管理器
func New() *Manager {
	return &Manager{
		players: make(map[string]*PlayerFragments),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// EnsurePlayer 確保玩家記錄存在
func (m *Manager) EnsurePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.players[playerID]; !ok {
		m.players[playerID] = &PlayerFragments{}
	}
}

// RemovePlayer 移除玩家記錄
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.players, playerID)
}

// TryDrop 嘗試掉落碎片（擊破目標後呼叫）
// defID: 目標物 ID（決定掉落類型）
// betCost: 投注金額（決定大獎金額）
// isBoss: 是否為 BOSS
func (m *Manager) TryDrop(playerID string, defID string, betCost int, isBoss bool) *DropResult {
	m.mu.Lock()
	defer m.mu.Unlock()

	pf, ok := m.players[playerID]
	if !ok {
		pf = &PlayerFragments{}
		m.players[playerID] = pf
	}

	// 決定掉落類型和機率
	var fragType FragmentType
	var dropChance float64

	if isBoss {
		// BOSS 擊破：50% 機率掉落金碎片
		fragType = FragmentGold
		dropChance = 0.50
	} else {
		switch defID {
		case "T103", "T104", "T105": // 特殊目標：30% 機率掉落銀碎片
			fragType = FragmentSilver
			dropChance = 0.30
		case "T101", "T102": // 稀有目標：20% 機率掉落銀碎片
			fragType = FragmentSilver
			dropChance = 0.20
		default: // 普通目標：8% 機率掉落銅碎片
			fragType = FragmentBronze
			dropChance = 0.08
		}
	}

	// 機率判定
	if m.rng.Float64() >= dropChance {
		return &DropResult{Dropped: false}
	}

	// 增加碎片
	var newCount int
	switch fragType {
	case FragmentBronze:
		pf.Bronze++
		newCount = pf.Bronze
	case FragmentSilver:
		pf.Silver++
		newCount = pf.Silver
	case FragmentGold:
		pf.Gold++
		newCount = pf.Gold
	}

	def := rewardDefs[fragType]
	result := &DropResult{
		Dropped:      true,
		FragmentType: fragType,
		NewCount:     newCount,
		Label:        def.Label,
		Color:        def.Color,
	}

	// 檢查是否集齊
	if newCount >= def.Required {
		// 重置碎片
		switch fragType {
		case FragmentBronze:
			pf.Bronze = 0
		case FragmentSilver:
			pf.Silver = 0
		case FragmentGold:
			pf.Gold = 0
		}
		result.IsComplete = true
		result.NewCount = 0
		result.Reward = betCost * def.RewardMult
	}

	return result
}

// GetSnapshot 取得玩家碎片快照
func (m *Manager) GetSnapshot(playerID string) PlayerFragments {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if pf, ok := m.players[playerID]; ok {
		return *pf
	}
	return PlayerFragments{}
}

// GetRewardDef 取得獎勵定義
func GetRewardDef(ft FragmentType) FragmentReward {
	return rewardDefs[ft]
}

// GetAllRewardDefs 取得所有獎勵定義
func GetAllRewardDefs() map[FragmentType]FragmentReward {
	return rewardDefs
}
