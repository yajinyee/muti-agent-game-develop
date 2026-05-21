// Package bounty 全服目標懸賞系統（DAY-137）
// 業界依據：strivecloud.io 2026「social streaks + tiered rewards」
// 玩家可對高價值目標下懸賞，擊破者獲得懸賞金額，增加社交互動
package bounty

import (
	"fmt"
	"sync"
	"time"
)

// BountyConfig 懸賞設定
type BountyConfig struct {
	MaxActiveBounties int     // 同時最多懸賞數（預設 3）
	MinBountyAmount   int     // 最低懸賞金額（預設 100）
	MaxBountyAmount   int     // 最高懸賞金額（預設 5000）
	BountyDuration    float64 // 懸賞持續秒數（預設 60 秒）
	CooldownPerPlayer float64 // 每個玩家下懸賞的冷卻（預設 120 秒）
}

// DefaultConfig 預設設定
func DefaultConfig() BountyConfig {
	return BountyConfig{
		MaxActiveBounties: 3,
		MinBountyAmount:   100,
		MaxBountyAmount:   5000,
		BountyDuration:    60.0,
		CooldownPerPlayer: 120.0,
	}
}

// Bounty 單筆懸賞
type Bounty struct {
	ID               string
	TargetInstanceID string
	TargetDefID      string
	TargetName       string
	TargetMult       float64
	PosterID         string  // 下懸賞的玩家 ID
	PosterName       string  // 下懸賞的玩家名稱
	Amount           int     // 懸賞金額
	PostedAt         time.Time
	ExpiresAt        time.Time
	IsActive         bool
	KillerID         string // 擊破者 ID（空=未擊破）
	KillerName       string // 擊破者名稱
}

// Manager 懸賞管理器
type Manager struct {
	mu             sync.RWMutex
	config         BountyConfig
	bounties       map[string]*Bounty // bountyID -> Bounty
	playerCooldown map[string]time.Time // playerID -> 下次可下懸賞時間
	nextID         int
}

// New 建立懸賞管理器
func New(cfg BountyConfig) *Manager {
	return &Manager{
		config:         cfg,
		bounties:       make(map[string]*Bounty),
		playerCooldown: make(map[string]time.Time),
	}
}

// NewDefault 使用預設設定建立管理器
func NewDefault() *Manager {
	return New(DefaultConfig())
}

// CanPost 玩家是否可以下懸賞
func (m *Manager) CanPost(playerID string) (bool, int) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 檢查冷卻
	if cooldownEnd, ok := m.playerCooldown[playerID]; ok {
		if time.Now().Before(cooldownEnd) {
			remaining := int(time.Until(cooldownEnd).Seconds())
			return false, remaining
		}
	}

	// 檢查活躍懸賞數
	activeCount := 0
	for _, b := range m.bounties {
		if b.IsActive {
			activeCount++
		}
	}
	if activeCount >= m.config.MaxActiveBounties {
		return false, -1 // -1 表示懸賞已滿
	}

	return true, 0
}

// PostBounty 下懸賞
// 回傳：(bountyID, error_code)
// error_code: "" = 成功, "cooldown" = 冷卻中, "full" = 懸賞已滿, "invalid_amount" = 金額無效
func (m *Manager) PostBounty(playerID, playerName, instanceID, defID, targetName string, mult float64, amount int) (string, string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 金額驗證
	if amount < m.config.MinBountyAmount || amount > m.config.MaxBountyAmount {
		return "", "invalid_amount"
	}

	// 冷卻檢查
	if cooldownEnd, ok := m.playerCooldown[playerID]; ok {
		if time.Now().Before(cooldownEnd) {
			return "", "cooldown"
		}
	}

	// 活躍懸賞數檢查
	activeCount := 0
	for _, b := range m.bounties {
		if b.IsActive {
			activeCount++
		}
	}
	if activeCount >= m.config.MaxActiveBounties {
		return "", "full"
	}

	// 建立懸賞
	m.nextID++
	bountyID := fmt.Sprintf("bounty_%d", m.nextID)
	now := time.Now()
	b := &Bounty{
		ID:               bountyID,
		TargetInstanceID: instanceID,
		TargetDefID:      defID,
		TargetName:       targetName,
		TargetMult:       mult,
		PosterID:         playerID,
		PosterName:       playerName,
		Amount:           amount,
		PostedAt:         now,
		ExpiresAt:        now.Add(time.Duration(m.config.BountyDuration * float64(time.Second))),
		IsActive:         true,
	}
	m.bounties[bountyID] = b

	// 設定冷卻
	m.playerCooldown[playerID] = now.Add(time.Duration(m.config.CooldownPerPlayer * float64(time.Second)))

	return bountyID, ""
}

// ClaimBounty 擊破懸賞目標，領取懸賞
// 回傳：(totalBounty, claimedBounties, isAnyBounty)
func (m *Manager) ClaimBounty(instanceID, killerID, killerName string) (int, []*Bounty, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	totalAmount := 0
	claimed := make([]*Bounty, 0)

	for _, b := range m.bounties {
		if !b.IsActive || b.TargetInstanceID != instanceID {
			continue
		}
		b.IsActive = false
		b.KillerID = killerID
		b.KillerName = killerName
		totalAmount += b.Amount
		claimed = append(claimed, b)
	}

	return totalAmount, claimed, len(claimed) > 0
}

// CheckExpiry 檢查過期懸賞
// 回傳過期的懸賞列表
func (m *Manager) CheckExpiry() []*Bounty {
	m.mu.Lock()
	defer m.mu.Unlock()

	expired := make([]*Bounty, 0)
	now := time.Now()
	for _, b := range m.bounties {
		if b.IsActive && now.After(b.ExpiresAt) {
			b.IsActive = false
			expired = append(expired, b)
		}
	}
	return expired
}

// CancelBountyForTarget 目標消失時取消懸賞（退款）
// 回傳需要退款的懸賞列表
func (m *Manager) CancelBountyForTarget(instanceID string) []*Bounty {
	m.mu.Lock()
	defer m.mu.Unlock()

	cancelled := make([]*Bounty, 0)
	for _, b := range m.bounties {
		if b.IsActive && b.TargetInstanceID == instanceID {
			b.IsActive = false
			cancelled = append(cancelled, b)
		}
	}
	return cancelled
}

// GetActiveBounties 取得所有活躍懸賞
func (m *Manager) GetActiveBounties() []*BountySnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*BountySnapshot, 0)
	for _, b := range m.bounties {
		if b.IsActive {
			result = append(result, &BountySnapshot{
				ID:               b.ID,
				TargetInstanceID: b.TargetInstanceID,
				TargetDefID:      b.TargetDefID,
				TargetName:       b.TargetName,
				TargetMult:       b.TargetMult,
				PosterID:         b.PosterID,
				PosterName:       b.PosterName,
				Amount:           b.Amount,
				SecondsLeft:      time.Until(b.ExpiresAt).Seconds(),
			})
		}
	}
	return result
}

// GetBountiesForTarget 取得特定目標的所有活躍懸賞
func (m *Manager) GetBountiesForTarget(instanceID string) []*BountySnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*BountySnapshot, 0)
	for _, b := range m.bounties {
		if b.IsActive && b.TargetInstanceID == instanceID {
			result = append(result, &BountySnapshot{
				ID:               b.ID,
				TargetInstanceID: b.TargetInstanceID,
				TargetDefID:      b.TargetDefID,
				TargetName:       b.TargetName,
				TargetMult:       b.TargetMult,
				PosterID:         b.PosterID,
				PosterName:       b.PosterName,
				Amount:           b.Amount,
				SecondsLeft:      time.Until(b.ExpiresAt).Seconds(),
			})
		}
	}
	return result
}

// GetPlayerCooldown 取得玩家冷卻剩餘秒數
func (m *Manager) GetPlayerCooldown(playerID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if cooldownEnd, ok := m.playerCooldown[playerID]; ok {
		remaining := int(time.Until(cooldownEnd).Seconds())
		if remaining > 0 {
			return remaining
		}
	}
	return 0
}

// BountySnapshot 懸賞快照（用於廣播）
type BountySnapshot struct {
	ID               string
	TargetInstanceID string
	TargetDefID      string
	TargetName       string
	TargetMult       float64
	PosterID         string
	PosterName       string
	Amount           int
	SecondsLeft      float64
}
