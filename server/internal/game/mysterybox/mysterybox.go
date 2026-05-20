// Package mysterybox 神秘寶箱系統（DAY-090）
// 業界依據：nerdbot.com 2026-05-02 確認「mystery rewards」是 2026 年 iGaming 最熱門留存機制
// 玩家擊破特定目標後有機率掉落神秘寶箱，開箱獲得隨機獎勵
package mysterybox

import (
	"math/rand"
	"sync"
	"time"
)

// BoxRarity 寶箱稀有度
type BoxRarity string

const (
	RarityCommon    BoxRarity = "common"    // 普通（灰色）
	RarityRare      BoxRarity = "rare"      // 稀有（藍色）
	RarityEpic      BoxRarity = "epic"      // 史詩（紫色）
	RarityLegendary BoxRarity = "legendary" // 傳說（金色）
)

// RewardType 獎勵類型
type RewardType string

const (
	RewardCoins         RewardType = "coins"          // 金幣
	RewardBombCharge    RewardType = "bomb_charge"     // 炸彈砲充能
	RewardLaserCharge   RewardType = "laser_charge"    // 雷射砲充能
	RewardFreezeCharge  RewardType = "freeze_charge"   // 冰凍砲充能
	RewardMultiplier    RewardType = "multiplier"      // 下一次攻擊倍率加成
	RewardJackpotTicket RewardType = "jackpot_ticket"  // Jackpot 抽獎券（增加 Jackpot 貢獻）
)

// BoxReward 寶箱獎勵定義
type BoxReward struct {
	Type     RewardType `json:"type"`
	Amount   int        `json:"amount"`   // 金幣數量 / 充能數量 / 倍率（×10 = 實際倍率）
	Label    string     `json:"label"`    // 顯示文字
	Icon     string     `json:"icon"`     // 圖示（emoji）
	Color    string     `json:"color"`    // 顏色（hex）
	Weight   int        `json:"-"`        // 抽獎權重（不傳給 Client）
}

// BoxDef 寶箱定義
type BoxDef struct {
	Rarity      BoxRarity   `json:"rarity"`
	Name        string      `json:"name"`
	Icon        string      `json:"icon"`
	Color       string      `json:"color"`
	GlowColor   string      `json:"glow_color"`
	DropChance  float64     `json:"-"` // 掉落機率（0.0-1.0）
	Rewards     []BoxReward `json:"-"` // 可能的獎勵池
}

// AvailableBoxes 所有寶箱定義
var AvailableBoxes = []BoxDef{
	{
		Rarity:     RarityCommon,
		Name:       "普通寶箱",
		Icon:       "📦",
		Color:      "#A0A0A0",
		GlowColor:  "#C0C0C0",
		DropChance: 0.08, // 8% 掉落率
		Rewards: []BoxReward{
			{Type: RewardCoins, Amount: 200, Label: "+200 金幣", Icon: "🪙", Color: "#FFD700", Weight: 40},
			{Type: RewardCoins, Amount: 500, Label: "+500 金幣", Icon: "🪙", Color: "#FFD700", Weight: 30},
			{Type: RewardCoins, Amount: 1000, Label: "+1000 金幣", Icon: "🪙", Color: "#FFD700", Weight: 20},
			{Type: RewardFreezeCharge, Amount: 1, Label: "冰凍砲 ×1", Icon: "❄️", Color: "#87CEEB", Weight: 10},
		},
	},
	{
		Rarity:     RarityRare,
		Name:       "稀有寶箱",
		Icon:       "💎",
		Color:      "#4169E1",
		GlowColor:  "#6495ED",
		DropChance: 0.04, // 4% 掉落率
		Rewards: []BoxReward{
			{Type: RewardCoins, Amount: 1000, Label: "+1000 金幣", Icon: "🪙", Color: "#FFD700", Weight: 30},
			{Type: RewardCoins, Amount: 3000, Label: "+3000 金幣", Icon: "🪙", Color: "#FFD700", Weight: 25},
			{Type: RewardBombCharge, Amount: 1, Label: "炸彈砲 ×1", Icon: "💣", Color: "#FF6B35", Weight: 20},
			{Type: RewardLaserCharge, Amount: 1, Label: "雷射砲 ×1", Icon: "⚡", Color: "#00FFFF", Weight: 15},
			{Type: RewardMultiplier, Amount: 20, Label: "下次攻擊 ×2.0", Icon: "✨", Color: "#FFD700", Weight: 10},
		},
	},
	{
		Rarity:     RarityEpic,
		Name:       "史詩寶箱",
		Icon:       "🔮",
		Color:      "#9B59B6",
		GlowColor:  "#C39BD3",
		DropChance: 0.015, // 1.5% 掉落率
		Rewards: []BoxReward{
			{Type: RewardCoins, Amount: 5000, Label: "+5000 金幣", Icon: "🪙", Color: "#FFD700", Weight: 25},
			{Type: RewardCoins, Amount: 10000, Label: "+10000 金幣", Icon: "🪙", Color: "#FFD700", Weight: 20},
			{Type: RewardBombCharge, Amount: 2, Label: "炸彈砲 ×2", Icon: "💣", Color: "#FF6B35", Weight: 15},
			{Type: RewardLaserCharge, Amount: 2, Label: "雷射砲 ×2", Icon: "⚡", Color: "#00FFFF", Weight: 15},
			{Type: RewardMultiplier, Amount: 30, Label: "下次攻擊 ×3.0", Icon: "✨", Color: "#FFD700", Weight: 15},
			{Type: RewardJackpotTicket, Amount: 5, Label: "Jackpot 券 ×5", Icon: "🎰", Color: "#FF4500", Weight: 10},
		},
	},
	{
		Rarity:     RarityLegendary,
		Name:       "傳說寶箱",
		Icon:       "👑",
		Color:      "#FFD700",
		GlowColor:  "#FFA500",
		DropChance: 0.005, // 0.5% 掉落率（BOSS 擊殺時提升到 5%）
		Rewards: []BoxReward{
			{Type: RewardCoins, Amount: 20000, Label: "+20000 金幣", Icon: "🪙", Color: "#FFD700", Weight: 20},
			{Type: RewardCoins, Amount: 50000, Label: "+50000 金幣", Icon: "🪙", Color: "#FFD700", Weight: 10},
			{Type: RewardBombCharge, Amount: 3, Label: "炸彈砲 ×3", Icon: "💣", Color: "#FF6B35", Weight: 15},
			{Type: RewardLaserCharge, Amount: 3, Label: "雷射砲 ×3", Icon: "⚡", Color: "#00FFFF", Weight: 15},
			{Type: RewardFreezeCharge, Amount: 3, Label: "冰凍砲 ×3", Icon: "❄️", Color: "#87CEEB", Weight: 15},
			{Type: RewardMultiplier, Amount: 50, Label: "下次攻擊 ×5.0", Icon: "✨", Color: "#FFD700", Weight: 15},
			{Type: RewardJackpotTicket, Amount: 20, Label: "Jackpot 券 ×20", Icon: "🎰", Color: "#FF4500", Weight: 10},
		},
	},
}

// PendingMultiplier 待使用的攻擊倍率加成
type PendingMultiplier struct {
	Multiplier float64
	ExpiresAt  time.Time
}

// Manager 神秘寶箱管理器
type Manager struct {
	mu                sync.RWMutex
	pendingMultiplier map[string]*PendingMultiplier // playerID -> 待使用倍率
	rng               *rand.Rand
}

// New 建立神秘寶箱管理器
func New() *Manager {
	return &Manager{
		pendingMultiplier: make(map[string]*PendingMultiplier),
		rng:               rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// TryDropBox 嘗試掉落寶箱（在擊破目標後呼叫）
// isBoss: 是否為 BOSS 擊殺（提升傳說寶箱機率）
// 回傳掉落的寶箱定義（nil = 未掉落）
func (m *Manager) TryDropBox(isBoss bool) *BoxDef {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range AvailableBoxes {
		box := &AvailableBoxes[i]
		chance := box.DropChance
		// BOSS 擊殺時，傳說寶箱機率提升 10 倍
		if isBoss && box.Rarity == RarityLegendary {
			chance *= 10.0
		}
		if m.rng.Float64() < chance {
			return box
		}
	}
	return nil
}

// OpenBox 開箱，回傳獲得的獎勵
func (m *Manager) OpenBox(rarity BoxRarity) *BoxReward {
	m.mu.Lock()
	defer m.mu.Unlock()

	var box *BoxDef
	for i := range AvailableBoxes {
		if AvailableBoxes[i].Rarity == rarity {
			box = &AvailableBoxes[i]
			break
		}
	}
	if box == nil || len(box.Rewards) == 0 {
		return nil
	}

	// 加權隨機選擇獎勵
	totalWeight := 0
	for _, r := range box.Rewards {
		totalWeight += r.Weight
	}
	roll := m.rng.Intn(totalWeight)
	cumulative := 0
	for _, r := range box.Rewards {
		cumulative += r.Weight
		if roll < cumulative {
			reward := r // 複製
			return &reward
		}
	}
	// fallback
	r := box.Rewards[0]
	return &r
}

// SetPendingMultiplier 設定待使用的攻擊倍率加成（開箱獲得 multiplier 時呼叫）
func (m *Manager) SetPendingMultiplier(playerID string, mult float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pendingMultiplier[playerID] = &PendingMultiplier{
		Multiplier: mult,
		ExpiresAt:  time.Now().Add(60 * time.Second), // 60 秒內有效
	}
}

// ConsumePendingMultiplier 消耗待使用的攻擊倍率加成（下次攻擊時呼叫）
// 回傳倍率（1.0 = 無加成）
func (m *Manager) ConsumePendingMultiplier(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	pm, ok := m.pendingMultiplier[playerID]
	if !ok {
		return 1.0
	}
	if time.Now().After(pm.ExpiresAt) {
		delete(m.pendingMultiplier, playerID)
		return 1.0
	}
	mult := pm.Multiplier
	delete(m.pendingMultiplier, playerID)
	return mult
}

// GetPendingMultiplier 查詢待使用的攻擊倍率加成（不消耗）
func (m *Manager) GetPendingMultiplier(playerID string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pm, ok := m.pendingMultiplier[playerID]
	if !ok {
		return 1.0
	}
	if time.Now().After(pm.ExpiresAt) {
		return 1.0
	}
	return pm.Multiplier
}

// RemovePlayer 移除玩家狀態
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.pendingMultiplier, playerID)
}

// GetBoxDef 取得寶箱定義（供 handler 使用）
func GetBoxDef(rarity BoxRarity) *BoxDef {
	for i := range AvailableBoxes {
		if AvailableBoxes[i].Rarity == rarity {
			return &AvailableBoxes[i]
		}
	}
	return nil
}

// ---- 背包管理 ----

// playerInventory 玩家背包（稀有度 -> 數量）
type playerInventory map[BoxRarity]int

// inventories 所有玩家背包
var inventories = struct {
	sync.RWMutex
	data map[string]playerInventory
}{data: make(map[string]playerInventory)}

// AddBox 加入寶箱到玩家背包
func (m *Manager) AddBox(playerID string, rarity BoxRarity) {
	inventories.Lock()
	defer inventories.Unlock()
	if _, ok := inventories.data[playerID]; !ok {
		inventories.data[playerID] = make(playerInventory)
	}
	inventories.data[playerID][rarity]++
}

// RemoveBox 從玩家背包移除一個寶箱
func (m *Manager) RemoveBox(playerID string, rarity BoxRarity) {
	inventories.Lock()
	defer inventories.Unlock()
	if inv, ok := inventories.data[playerID]; ok {
		if inv[rarity] > 0 {
			inv[rarity]--
		}
	}
}

// HasBox 確認玩家是否有指定稀有度的寶箱
func (m *Manager) HasBox(playerID string, rarity BoxRarity) bool {
	inventories.RLock()
	defer inventories.RUnlock()
	if inv, ok := inventories.data[playerID]; ok {
		return inv[rarity] > 0
	}
	return false
}

// GetBoxCount 取得玩家指定稀有度的寶箱數量
func (m *Manager) GetBoxCount(playerID string, rarity BoxRarity) int {
	inventories.RLock()
	defer inventories.RUnlock()
	if inv, ok := inventories.data[playerID]; ok {
		return inv[rarity]
	}
	return 0
}

// GetInventory 取得玩家完整背包快照
func (m *Manager) GetInventory(playerID string) map[BoxRarity]int {
	inventories.RLock()
	defer inventories.RUnlock()
	result := make(map[BoxRarity]int)
	if inv, ok := inventories.data[playerID]; ok {
		for k, v := range inv {
			result[k] = v
		}
	}
	return result
}
