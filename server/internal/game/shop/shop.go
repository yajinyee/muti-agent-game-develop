// Package shop 商店系統（DAY-094）
// 提供限時特賣、道具購買、金幣包等消費管道
// 每日 UTC+8 00:00 重置特賣商品
package shop

import (
	"sync"
	"time"
)

// ItemType 商品類型
type ItemType string

const (
	ItemCoinPack      ItemType = "coin_pack"      // 金幣包（用真實貨幣購買，此處模擬）
	ItemSpecialItem   ItemType = "special_item"   // 特殊道具（用金幣購買）
	ItemFlashSale     ItemType = "flash_sale"     // 限時特賣（折扣商品）
)

// Item 商品定義
type Item struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        ItemType `json:"type"`
	Price       int      `json:"price"`        // 金幣價格（0=免費）
	OrigPrice   int      `json:"orig_price"`   // 原價（用於顯示折扣）
	Reward      ItemReward `json:"reward"`
	Stock       int      `json:"stock"`        // 庫存（-1=無限）
	LimitPerDay int      `json:"limit_per_day"` // 每日購買上限（0=無限）
	IsFlashSale bool     `json:"is_flash_sale"` // 是否為限時特賣
	FlashEndAt  int64    `json:"flash_end_at"`  // 特賣結束時間（Unix ms，0=不限時）
}

// ItemReward 商品獎勵
type ItemReward struct {
	Coins         int     `json:"coins"`          // 金幣
	BombCharge    int     `json:"bomb_charge"`    // 炸彈充能
	LaserCharge   int     `json:"laser_charge"`   // 雷射充能
	FreezeCharge  int     `json:"freeze_charge"`  // 冰凍充能
	AttackMult    float64 `json:"attack_mult"`    // 下次攻擊倍率加成
	SeasonPoints  int     `json:"season_points"`  // 賽季積分
}

// PurchaseRecord 購買記錄
type PurchaseRecord struct {
	PlayerID  string
	ItemID    string
	Count     int
	Date      string // "2026-05-20"
}

// DefaultItems 預設商品列表
var DefaultItems = []Item{
	// 金幣包（模擬，實際應接入支付系統）
	{
		ID:          "coin_pack_small",
		Name:        "💰 小金幣包",
		Description: "獲得 5000 金幣",
		Type:        ItemCoinPack,
		Price:       0, // 模擬免費（實際應有價格）
		OrigPrice:   0,
		Reward:      ItemReward{Coins: 5000},
		Stock:       -1,
		LimitPerDay: 1,
	},
	{
		ID:          "coin_pack_medium",
		Name:        "💰 中金幣包",
		Description: "獲得 15000 金幣",
		Type:        ItemCoinPack,
		Price:       0,
		OrigPrice:   0,
		Reward:      ItemReward{Coins: 15000},
		Stock:       -1,
		LimitPerDay: 1,
	},
	// 特殊道具（用金幣購買）
	{
		ID:          "bomb_bundle",
		Name:        "💣 炸彈套裝",
		Description: "獲得 3 發炸彈充能",
		Type:        ItemSpecialItem,
		Price:       2000,
		OrigPrice:   2000,
		Reward:      ItemReward{BombCharge: 3},
		Stock:       -1,
		LimitPerDay: 3,
	},
	{
		ID:          "laser_bundle",
		Name:        "⚡ 雷射套裝",
		Description: "獲得 3 發雷射充能",
		Type:        ItemSpecialItem,
		Price:       3000,
		OrigPrice:   3000,
		Reward:      ItemReward{LaserCharge: 3},
		Stock:       -1,
		LimitPerDay: 3,
	},
	{
		ID:          "freeze_bundle",
		Name:        "❄️ 冰凍套裝",
		Description: "獲得 3 發冰凍充能",
		Type:        ItemSpecialItem,
		Price:       1500,
		OrigPrice:   1500,
		Reward:      ItemReward{FreezeCharge: 3},
		Stock:       -1,
		LimitPerDay: 3,
	},
	{
		ID:          "attack_boost",
		Name:        "🔥 攻擊加成",
		Description: "下次攻擊倍率 ×2.0",
		Type:        ItemSpecialItem,
		Price:       5000,
		OrigPrice:   5000,
		Reward:      ItemReward{AttackMult: 2.0},
		Stock:       -1,
		LimitPerDay: 2,
	},
	{
		ID:          "season_boost",
		Name:        "⭐ 賽季加速",
		Description: "獲得 100 賽季積分",
		Type:        ItemSpecialItem,
		Price:       3000,
		OrigPrice:   3000,
		Reward:      ItemReward{SeasonPoints: 100},
		Stock:       -1,
		LimitPerDay: 5,
	},
}

// Manager 商店管理器
type Manager struct {
	mu        sync.RWMutex
	items     []Item                       // 商品列表
	purchases map[string]*PurchaseRecord   // "playerID:itemID:date" → record
	flashSaleEndAt time.Time               // 當前限時特賣結束時間
}

// New 建立新的商店管理器
func New() *Manager {
	m := &Manager{
		items:     make([]Item, len(DefaultItems)),
		purchases: make(map[string]*PurchaseRecord),
	}
	copy(m.items, DefaultItems)
	m.refreshFlashSale()
	return m
}

// refreshFlashSale 刷新限時特賣（每日 UTC+8 00:00 重置）
func (m *Manager) refreshFlashSale() {
	loc := time.FixedZone("UTC+8", 8*3600)
	now := time.Now().In(loc)
	// 特賣持續到今日 23:59:59
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)
	m.flashSaleEndAt = endOfDay

	// 隨機選 2 個特殊道具做限時特賣（7折）
	flashCount := 0
	for i := range m.items {
		if m.items[i].Type == ItemSpecialItem && flashCount < 2 {
			m.items[i].IsFlashSale = true
			m.items[i].OrigPrice = m.items[i].Price
			m.items[i].Price = int(float64(m.items[i].OrigPrice) * 0.7) // 7折
			m.items[i].FlashEndAt = endOfDay.UnixMilli()
			flashCount++
		}
	}
}

// GetItems 取得所有商品列表
func (m *Manager) GetItems() []Item {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.checkAndRefresh()
	result := make([]Item, len(m.items))
	copy(result, m.items)
	return result
}

// checkAndRefresh 檢查是否需要刷新特賣（必須在持有鎖的情況下呼叫）
func (m *Manager) checkAndRefresh() {
	if time.Now().After(m.flashSaleEndAt) {
		// 重置所有特賣標記
		for i := range m.items {
			if m.items[i].IsFlashSale {
				m.items[i].Price = m.items[i].OrigPrice
				m.items[i].IsFlashSale = false
				m.items[i].FlashEndAt = 0
			}
		}
		m.refreshFlashSale()
	}
}

// BuyResult 購買結果
type BuyResult struct {
	Success bool
	Reason  string
	Reward  ItemReward
	Item    Item
}

// BuyItem 購買商品
func (m *Manager) BuyItem(playerID string, itemID string, playerCoins int) BuyResult {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.checkAndRefresh()

	// 找到商品
	var item *Item
	for i := range m.items {
		if m.items[i].ID == itemID {
			item = &m.items[i]
			break
		}
	}
	if item == nil {
		return BuyResult{Success: false, Reason: "item_not_found"}
	}

	// 檢查庫存
	if item.Stock == 0 {
		return BuyResult{Success: false, Reason: "out_of_stock"}
	}

	// 檢查每日購買上限
	if item.LimitPerDay > 0 {
		loc := time.FixedZone("UTC+8", 8*3600)
		dateStr := time.Now().In(loc).Format("2006-01-02")
		key := playerID + ":" + itemID + ":" + dateStr
		if rec, ok := m.purchases[key]; ok && rec.Count >= item.LimitPerDay {
			return BuyResult{Success: false, Reason: "daily_limit_reached"}
		}
	}

	// 檢查金幣是否足夠
	if item.Price > 0 && playerCoins < item.Price {
		return BuyResult{Success: false, Reason: "insufficient_coins"}
	}

	// 扣除庫存
	if item.Stock > 0 {
		item.Stock--
	}

	// 記錄購買
	if item.LimitPerDay > 0 {
		loc := time.FixedZone("UTC+8", 8*3600)
		dateStr := time.Now().In(loc).Format("2006-01-02")
		key := playerID + ":" + itemID + ":" + dateStr
		if rec, ok := m.purchases[key]; ok {
			rec.Count++
		} else {
			m.purchases[key] = &PurchaseRecord{
				PlayerID: playerID,
				ItemID:   itemID,
				Count:    1,
				Date:     dateStr,
			}
		}
	}

	return BuyResult{
		Success: true,
		Reward:  item.Reward,
		Item:    *item,
	}
}

// GetPlayerDailyPurchases 取得玩家今日購買記錄（itemID → count）
func (m *Manager) GetPlayerDailyPurchases(playerID string) map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	loc := time.FixedZone("UTC+8", 8*3600)
	dateStr := time.Now().In(loc).Format("2006-01-02")

	result := make(map[string]int)
	for key, rec := range m.purchases {
		if rec.PlayerID == playerID && rec.Date == dateStr {
			_ = key
			result[rec.ItemID] = rec.Count
		}
	}
	return result
}

// GetFlashSaleEndAt 取得限時特賣結束時間（Unix ms）
func (m *Manager) GetFlashSaleEndAt() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.flashSaleEndAt.UnixMilli()
}

// Snapshot 商店快照
type Snapshot struct {
	Items          []Item         `json:"items"`
	FlashSaleEndAt int64          `json:"flash_sale_end_at"` // Unix ms
	SecondsLeft    int64          `json:"seconds_left"`      // 特賣剩餘秒數
}

// GetSnapshot 取得商店快照
func (m *Manager) GetSnapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.checkAndRefresh()

	items := make([]Item, len(m.items))
	copy(items, m.items)

	left := time.Until(m.flashSaleEndAt)
	if left < 0 {
		left = 0
	}

	return Snapshot{
		Items:          items,
		FlashSaleEndAt: m.flashSaleEndAt.UnixMilli(),
		SecondsLeft:    int64(left.Seconds()),
	}
}
