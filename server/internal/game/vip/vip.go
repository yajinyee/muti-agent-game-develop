// Package vip VIP 等級系統（DAY-078）
// 累積消費金幣解鎖 VIP 等級，不重置
// 5 個等級：Bronze/Silver/Gold/Platinum/Diamond
// 每個等級提供：金幣返還加成、每日獎勵加成、特殊稱號
package vip

import (
	"sync"
	"time"
)

// VIPTier VIP 等級定義
type VIPTier struct {
	Level         int     `json:"level"`          // 1-5
	Name          string  `json:"name"`           // 等級名稱
	Icon          string  `json:"icon"`           // 圖示
	Color         string  `json:"color"`          // 顏色（hex）
	SpendRequired int     `json:"spend_required"` // 累積消費門檻（金幣）
	CashbackRate  float64 `json:"cashback_rate"`  // 金幣返還比例（0.01 = 1%）
	DailyBonusMult float64 `json:"daily_bonus_mult"` // 每日登入獎勵倍率
	TitleID       string  `json:"title_id"`       // 解鎖稱號 ID
	TitleName     string  `json:"title_name"`     // 解鎖稱號名稱
	WeeklyBonus   int     `json:"weekly_bonus"`   // 每週固定獎勵金幣
}

// VIPTiers VIP 等級定義（5 個等級）
var VIPTiers = []VIPTier{
	{
		Level:          1,
		Name:           "青銅會員",
		Icon:           "🥉",
		Color:          "#CD7F32",
		SpendRequired:  10000,
		CashbackRate:   0.01,
		DailyBonusMult: 1.1,
		TitleID:        "vip_bronze",
		TitleName:      "青銅會員",
		WeeklyBonus:    500,
	},
	{
		Level:          2,
		Name:           "白銀會員",
		Icon:           "🥈",
		Color:          "#C0C0C0",
		SpendRequired:  50000,
		CashbackRate:   0.02,
		DailyBonusMult: 1.2,
		TitleID:        "vip_silver",
		TitleName:      "白銀會員",
		WeeklyBonus:    1500,
	},
	{
		Level:          3,
		Name:           "黃金會員",
		Icon:           "🥇",
		Color:          "#FFD700",
		SpendRequired:  200000,
		CashbackRate:   0.03,
		DailyBonusMult: 1.5,
		TitleID:        "vip_gold",
		TitleName:      "黃金會員",
		WeeklyBonus:    5000,
	},
	{
		Level:          4,
		Name:           "白金會員",
		Icon:           "💎",
		Color:          "#E5E4E2",
		SpendRequired:  500000,
		CashbackRate:   0.05,
		DailyBonusMult: 2.0,
		TitleID:        "vip_platinum",
		TitleName:      "白金會員",
		WeeklyBonus:    15000,
	},
	{
		Level:          5,
		Name:           "鑽石會員",
		Icon:           "👑",
		Color:          "#B9F2FF",
		SpendRequired:  2000000,
		CashbackRate:   0.08,
		DailyBonusMult: 3.0,
		TitleID:        "vip_diamond",
		TitleName:      "鑽石會員",
		WeeklyBonus:    50000,
	},
}

// PlayerVIPData 玩家 VIP 資料
type PlayerVIPData struct {
	PlayerID       string    `json:"player_id"`
	TotalSpend     int       `json:"total_spend"`      // 累積消費金幣
	VIPLevel       int       `json:"vip_level"`        // 當前 VIP 等級（0=未達等級1）
	LastWeeklyAt   time.Time `json:"last_weekly_at"`   // 上次領取週獎勵時間
	LastUpdated    time.Time `json:"last_updated"`
}

// LevelUpResult VIP 升級結果
type LevelUpResult struct {
	NewLevel    int
	TierName    string
	TierIcon    string
	TierColor   string
	TitleID     string
	TitleName   string
	WeeklyBonus int
}

// WeeklyClaimResult 週獎勵領取結果
type WeeklyClaimResult struct {
	Coins   int
	VIPLevel int
	TierName string
}

// Manager VIP 等級管理器
type Manager struct {
	mu      sync.RWMutex
	players map[string]*PlayerVIPData // playerID → data
}

// New 建立新的 VIP 管理器
func New() *Manager {
	return &Manager{
		players: make(map[string]*PlayerVIPData),
	}
}

// GetOrCreate 取得或建立玩家 VIP 資料
func (m *Manager) GetOrCreate(playerID string) *PlayerVIPData {
	m.mu.Lock()
	defer m.mu.Unlock()

	if data, ok := m.players[playerID]; ok {
		return data
	}
	data := &PlayerVIPData{
		PlayerID:    playerID,
		TotalSpend:  0,
		VIPLevel:    0,
		LastUpdated: time.Now(),
	}
	m.players[playerID] = data
	return data
}

// AddSpend 增加累積消費，回傳新等級（若升級）和升級結果
func (m *Manager) AddSpend(playerID string, amount int) (newLevel int, levelUp *LevelUpResult) {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, ok := m.players[playerID]
	if !ok {
		data = &PlayerVIPData{
			PlayerID:    playerID,
			TotalSpend:  0,
			VIPLevel:    0,
			LastUpdated: time.Now(),
		}
		m.players[playerID] = data
	}

	data.TotalSpend += amount
	data.LastUpdated = time.Now()

	// 計算新等級
	newVIPLevel := 0
	for _, tier := range VIPTiers {
		if data.TotalSpend >= tier.SpendRequired {
			newVIPLevel = tier.Level
		}
	}

	// 檢查是否升級
	if newVIPLevel > data.VIPLevel {
		oldLevel := data.VIPLevel
		data.VIPLevel = newVIPLevel
		_ = oldLevel

		// 找到新等級定義
		var tierDef *VIPTier
		for i := range VIPTiers {
			if VIPTiers[i].Level == newVIPLevel {
				tierDef = &VIPTiers[i]
				break
			}
		}
		if tierDef != nil {
			return newVIPLevel, &LevelUpResult{
				NewLevel:    newVIPLevel,
				TierName:    tierDef.Name,
				TierIcon:    tierDef.Icon,
				TierColor:   tierDef.Color,
				TitleID:     tierDef.TitleID,
				TitleName:   tierDef.TitleName,
				WeeklyBonus: tierDef.WeeklyBonus,
			}
		}
	}

	return data.VIPLevel, nil
}

// GetCashback 計算金幣返還金額（依 VIP 等級）
func (m *Manager) GetCashback(playerID string, spendAmount int) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, ok := m.players[playerID]
	if !ok || data.VIPLevel == 0 {
		return 0
	}

	var tier *VIPTier
	for i := range VIPTiers {
		if VIPTiers[i].Level == data.VIPLevel {
			tier = &VIPTiers[i]
			break
		}
	}
	if tier == nil {
		return 0
	}

	cashback := int(float64(spendAmount) * tier.CashbackRate)
	return cashback
}

// GetDailyBonusMult 取得每日登入獎勵倍率（依 VIP 等級）
func (m *Manager) GetDailyBonusMult(playerID string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, ok := m.players[playerID]
	if !ok || data.VIPLevel == 0 {
		return 1.0
	}

	for _, tier := range VIPTiers {
		if tier.Level == data.VIPLevel {
			return tier.DailyBonusMult
		}
	}
	return 1.0
}

// ClaimWeeklyBonus 領取週獎勵，回傳獎勵金幣（0=不可領取）
func (m *Manager) ClaimWeeklyBonus(playerID string) *WeeklyClaimResult {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, ok := m.players[playerID]
	if !ok || data.VIPLevel == 0 {
		return nil
	}

	// 確認距離上次領取超過 7 天
	now := time.Now()
	if !data.LastWeeklyAt.IsZero() && now.Sub(data.LastWeeklyAt) < 7*24*time.Hour {
		return nil
	}

	// 找到等級定義
	var tier *VIPTier
	for i := range VIPTiers {
		if VIPTiers[i].Level == data.VIPLevel {
			tier = &VIPTiers[i]
			break
		}
	}
	if tier == nil {
		return nil
	}

	data.LastWeeklyAt = now
	return &WeeklyClaimResult{
		Coins:    tier.WeeklyBonus,
		VIPLevel: data.VIPLevel,
		TierName: tier.Name,
	}
}

// GetSnapshot 取得玩家 VIP 快照
func (m *Manager) GetSnapshot(playerID string) VIPSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, ok := m.players[playerID]
	if !ok {
		return VIPSnapshot{
			PlayerID:      playerID,
			TotalSpend:    0,
			VIPLevel:      0,
			SpendToNext:   VIPTiers[0].SpendRequired,
			Progress:      0.0,
			CanClaimWeekly: false,
		}
	}

	// 計算下一等級所需消費
	spendToNext := 0
	progress := 1.0
	nextLevel := 0
	for _, tier := range VIPTiers {
		if data.TotalSpend < tier.SpendRequired {
			nextLevel = tier.Level
			spendToNext = tier.SpendRequired - data.TotalSpend
			// 計算進度
			prevRequired := 0
			if tier.Level > 1 {
				prevRequired = VIPTiers[tier.Level-2].SpendRequired
			}
			span := tier.SpendRequired - prevRequired
			earned := data.TotalSpend - prevRequired
			if span > 0 {
				progress = float64(earned) / float64(span)
				if progress > 1.0 {
					progress = 1.0
				}
				if progress < 0 {
					progress = 0
				}
			}
			break
		}
	}

	// 當前等級資訊
	var currentTier *VIPTier
	for i := range VIPTiers {
		if VIPTiers[i].Level == data.VIPLevel {
			currentTier = &VIPTiers[i]
			break
		}
	}

	// 是否可領取週獎勵
	canClaim := false
	if data.VIPLevel > 0 {
		canClaim = data.LastWeeklyAt.IsZero() || time.Since(data.LastWeeklyAt) >= 7*24*time.Hour
	}

	snap := VIPSnapshot{
		PlayerID:       playerID,
		TotalSpend:     data.TotalSpend,
		VIPLevel:       data.VIPLevel,
		NextLevel:      nextLevel,
		SpendToNext:    spendToNext,
		Progress:       progress,
		CanClaimWeekly: canClaim,
	}

	if currentTier != nil {
		snap.TierName        = currentTier.Name
		snap.TierIcon        = currentTier.Icon
		snap.TierColor       = currentTier.Color
		snap.CashbackRate    = currentTier.CashbackRate
		snap.DailyBonusMult  = currentTier.DailyBonusMult
		snap.WeeklyBonus     = currentTier.WeeklyBonus
	}

	return snap
}

// GetVIPLevel 取得玩家 VIP 等級（thread-safe）
func (m *Manager) GetVIPLevel(playerID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if data, ok := m.players[playerID]; ok {
		return data.VIPLevel
	}
	return 0
}

// GetTiers 取得所有 VIP 等級定義（用於 HTTP API）
func (m *Manager) GetTiers() []VIPTier {
	return VIPTiers
}

// VIPSnapshot 玩家 VIP 快照（用於 WebSocket 廣播）
type VIPSnapshot struct {
	PlayerID       string  `json:"player_id"`
	TotalSpend     int     `json:"total_spend"`
	VIPLevel       int     `json:"vip_level"`
	TierName       string  `json:"tier_name"`
	TierIcon       string  `json:"tier_icon"`
	TierColor      string  `json:"tier_color"`
	CashbackRate   float64 `json:"cashback_rate"`
	DailyBonusMult float64 `json:"daily_bonus_mult"`
	WeeklyBonus    int     `json:"weekly_bonus"`
	NextLevel      int     `json:"next_level"`      // 0 = 已滿級
	SpendToNext    int     `json:"spend_to_next"`   // 距離下一等級所需消費
	Progress       float64 `json:"progress"`        // 0.0-1.0 當前等級進度
	CanClaimWeekly bool    `json:"can_claim_weekly"` // 是否可領取週獎勵
}

// LoadState 從持久化資料恢復 VIP 狀態（DAY-098）
func (m *Manager) LoadState(playerID string, totalSpend int, vipLevel int, lastWeeklyAt time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.players[playerID] = &PlayerVIPData{
		PlayerID:    playerID,
		TotalSpend:  totalSpend,
		VIPLevel:    vipLevel,
		LastWeeklyAt: lastWeeklyAt,
		LastUpdated: time.Now(),
	}
}

// GetData 取得玩家 VIP 原始資料（用於持久化，DAY-098）
func (m *Manager) GetData(playerID string) (totalSpend int, vipLevel int, lastWeeklyAt time.Time) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if data, ok := m.players[playerID]; ok {
		return data.TotalSpend, data.VIPLevel, data.LastWeeklyAt
	}
	return 0, 0, time.Time{}
}
