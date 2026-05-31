// Package game — 任務幣兌換商店（DAY-348）
// 靈感來源：BGaming Quests 2026-05-27 發布
// 設計：玩家用任務幣兌換 BET 加成、金幣、特殊道具
// 任務幣來源：每日任務 + 每週挑戰 + 賽季通行證
package game

import (
	"fmt"
	"sync"
	"time"
)

// ShopItemType 商店道具類型
type ShopItemType string

const (
	ShopItemBetBoost    ShopItemType = "bet_boost"    // BET 加成（下一局 BET ×N）
	ShopItemCoinBonus   ShopItemType = "coin_bonus"   // 直接獲得金幣
	ShopItemXPBoost     ShopItemType = "xp_boost"     // 賽季 XP 加成（30分鐘）
	ShopItemLuckyCharm  ShopItemType = "lucky_charm"  // 幸運符（提高 Lucky 魚出現率 5 分鐘）
	ShopItemAutoAmmo    ShopItemType = "auto_ammo"    // AUTO 彈藥補充（免費射擊 30 秒）
)

// ShopItem 商店道具定義
type ShopItem struct {
	ID          string       `json:"id"`
	Type        ShopItemType `json:"type"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Cost        int          `json:"cost"`        // 任務幣價格
	Value       int          `json:"value"`       // 道具數值（倍率/金幣數/秒數）
	Stock       int          `json:"stock"`       // -1 = 無限庫存
	Icon        string       `json:"icon"`
}

// PlayerShopState 玩家商店狀態
type PlayerShopState struct {
	PlayerID      string                 `json:"player_id"`
	ActiveEffects map[string]*ShopEffect `json:"active_effects"` // effectID -> effect
	PurchaseLog   []ShopPurchaseRecord   `json:"purchase_log"`
}

// ShopEffect 已啟用的道具效果
type ShopEffect struct {
	ItemID    string    `json:"item_id"`
	ItemType  ShopItemType `json:"item_type"`
	Value     int       `json:"value"`
	ExpiresAt time.Time `json:"expires_at"` // 零值表示永久（直到使用）
	Used      bool      `json:"used"`
}

// ShopPurchaseRecord 購買記錄
type ShopPurchaseRecord struct {
	ItemID    string    `json:"item_id"`
	ItemName  string    `json:"item_name"`
	Cost      int       `json:"cost"`
	PurchasedAt time.Time `json:"purchased_at"`
}

// QuestShop 任務幣兌換商店
type QuestShop struct {
	mu     sync.RWMutex
	items  []ShopItem
	states map[string]*PlayerShopState // playerID -> state
}

// 商店道具定義（固定商品，不隨時間變化）
var defaultShopItems = []ShopItem{
	// ── BET 加成 ──────────────────────────────────────────────
	{
		ID:          "bet_boost_2x",
		Type:        ShopItemBetBoost,
		Name:        "BET 雙倍加成",
		Description: "下一局 BET 效果 ×2（一次性）",
		Cost:        30,
		Value:       2,
		Stock:       -1,
		Icon:        "🎯",
	},
	{
		ID:          "bet_boost_3x",
		Type:        ShopItemBetBoost,
		Name:        "BET 三倍加成",
		Description: "下一局 BET 效果 ×3（一次性）",
		Cost:        80,
		Value:       3,
		Stock:       -1,
		Icon:        "🎯",
	},
	{
		ID:          "bet_boost_5x",
		Type:        ShopItemBetBoost,
		Name:        "BET 五倍加成",
		Description: "下一局 BET 效果 ×5（一次性）",
		Cost:        180,
		Value:       5,
		Stock:       -1,
		Icon:        "🎯",
	},
	// ── 金幣獎勵 ──────────────────────────────────────────────
	{
		ID:          "coin_500",
		Type:        ShopItemCoinBonus,
		Name:        "500 金幣包",
		Description: "立即獲得 500 金幣",
		Cost:        20,
		Value:       500,
		Stock:       -1,
		Icon:        "🪙",
	},
	{
		ID:          "coin_2000",
		Type:        ShopItemCoinBonus,
		Name:        "2000 金幣包",
		Description: "立即獲得 2000 金幣",
		Cost:        70,
		Value:       2000,
		Stock:       -1,
		Icon:        "🪙",
	},
	{
		ID:          "coin_5000",
		Type:        ShopItemCoinBonus,
		Name:        "5000 金幣包",
		Description: "立即獲得 5000 金幣",
		Cost:        150,
		Value:       5000,
		Stock:       -1,
		Icon:        "🪙",
	},
	// ── XP 加成 ──────────────────────────────────────────────
	{
		ID:          "xp_boost_30m",
		Type:        ShopItemXPBoost,
		Name:        "XP 加速（30分鐘）",
		Description: "30分鐘內賽季 XP 獲取 ×2",
		Cost:        50,
		Value:       1800, // 30分鐘 = 1800秒
		Stock:       -1,
		Icon:        "⭐",
	},
	// ── 幸運符 ──────────────────────────────────────────────
	{
		ID:          "lucky_charm_5m",
		Type:        ShopItemLuckyCharm,
		Name:        "幸運符（5分鐘）",
		Description: "5分鐘內 Lucky 魚出現率 +30%",
		Cost:        60,
		Value:       300, // 5分鐘 = 300秒
		Stock:       -1,
		Icon:        "🍀",
	},
	// ── AUTO 彈藥 ──────────────────────────────────────────────
	{
		ID:          "auto_ammo_30s",
		Type:        ShopItemAutoAmmo,
		Name:        "AUTO 彈藥（30秒）",
		Description: "30秒內 AUTO 射擊不消耗金幣",
		Cost:        40,
		Value:       30, // 30秒
		Stock:       -1,
		Icon:        "🔫",
	},
}

// NewQuestShop 建立任務幣兌換商店
func NewQuestShop() *QuestShop {
	return &QuestShop{
		items:  defaultShopItems,
		states: make(map[string]*PlayerShopState),
	}
}

// GetItems 取得所有商品列表
func (s *QuestShop) GetItems() []ShopItem {
	return s.items
}

// GetPlayerState 取得玩家商店狀態
func (s *QuestShop) GetPlayerState(playerID string) *PlayerShopState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if state, ok := s.states[playerID]; ok {
		return state
	}
	return &PlayerShopState{
		PlayerID:      playerID,
		ActiveEffects: make(map[string]*ShopEffect),
		PurchaseLog:   []ShopPurchaseRecord{},
	}
}

// Purchase 購買道具
// 回傳：(成功, 錯誤訊息, 購買的道具)
func (s *QuestShop) Purchase(playerID string, itemID string, questCoins int) (bool, string, *ShopItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// 找到道具
	var item *ShopItem
	for i := range s.items {
		if s.items[i].ID == itemID {
			item = &s.items[i]
			break
		}
	}
	if item == nil {
		return false, "道具不存在", nil
	}
	
	// 檢查任務幣是否足夠
	if questCoins < item.Cost {
		return false, fmt.Sprintf("任務幣不足（需要 %d，擁有 %d）", item.Cost, questCoins), nil
	}
	
	// 取得或建立玩家狀態
	state, ok := s.states[playerID]
	if !ok {
		state = &PlayerShopState{
			PlayerID:      playerID,
			ActiveEffects: make(map[string]*ShopEffect),
			PurchaseLog:   []ShopPurchaseRecord{},
		}
		s.states[playerID] = state
	}
	
	// 建立效果
	effect := &ShopEffect{
		ItemID:   item.ID,
		ItemType: item.Type,
		Value:    item.Value,
		Used:     false,
	}
	
	// 設定過期時間（時間限制型道具）
	switch item.Type {
	case ShopItemXPBoost:
		effect.ExpiresAt = time.Now().Add(time.Duration(item.Value) * time.Second)
	case ShopItemLuckyCharm:
		effect.ExpiresAt = time.Now().Add(time.Duration(item.Value) * time.Second)
	case ShopItemAutoAmmo:
		effect.ExpiresAt = time.Now().Add(time.Duration(item.Value) * time.Second)
	// BetBoost 和 CoinBonus 是一次性使用，不設過期時間
	}
	
	// 儲存效果（同類型效果疊加時間）
	effectKey := string(item.Type)
	if existing, ok := state.ActiveEffects[effectKey]; ok && !existing.Used {
		if !existing.ExpiresAt.IsZero() {
			// 時間疊加
			remaining := time.Until(existing.ExpiresAt)
			if remaining > 0 {
				effect.ExpiresAt = time.Now().Add(remaining + time.Duration(item.Value)*time.Second)
			}
		}
	}
	state.ActiveEffects[effectKey] = effect
	
	// 記錄購買
	state.PurchaseLog = append(state.PurchaseLog, ShopPurchaseRecord{
		ItemID:      item.ID,
		ItemName:    item.Name,
		Cost:        item.Cost,
		PurchasedAt: time.Now(),
	})
	// 只保留最近 20 筆記錄
	if len(state.PurchaseLog) > 20 {
		state.PurchaseLog = state.PurchaseLog[len(state.PurchaseLog)-20:]
	}
	
	return true, "", item
}

// ConsumeEffect 消耗一次性效果（BetBoost 使用後標記）
func (s *QuestShop) ConsumeEffect(playerID string, effectType ShopItemType) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	state, ok := s.states[playerID]
	if !ok {
		return 0, false
	}
	
	effectKey := string(effectType)
	effect, ok := state.ActiveEffects[effectKey]
	if !ok || effect.Used {
		return 0, false
	}
	
	// 檢查是否過期
	if !effect.ExpiresAt.IsZero() && time.Now().After(effect.ExpiresAt) {
		delete(state.ActiveEffects, effectKey)
		return 0, false
	}
	
	// 一次性效果標記為已使用
	if effect.ItemType == ShopItemBetBoost || effect.ItemType == ShopItemCoinBonus {
		effect.Used = true
	}
	
	return effect.Value, true
}

// HasActiveEffect 檢查是否有有效的效果
func (s *QuestShop) HasActiveEffect(playerID string, effectType ShopItemType) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	state, ok := s.states[playerID]
	if !ok {
		return 0, false
	}
	
	effectKey := string(effectType)
	effect, ok := state.ActiveEffects[effectKey]
	if !ok || effect.Used {
		return 0, false
	}
	
	// 檢查是否過期
	if !effect.ExpiresAt.IsZero() && time.Now().After(effect.ExpiresAt) {
		return 0, false
	}
	
	return effect.Value, true
}

// GetActiveEffectsSummary 取得玩家所有有效效果摘要
func (s *QuestShop) GetActiveEffectsSummary(playerID string) []map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	state, ok := s.states[playerID]
	if !ok {
		return []map[string]interface{}{}
	}
	
	result := []map[string]interface{}{}
	for _, effect := range state.ActiveEffects {
		if effect.Used {
			continue
		}
		if !effect.ExpiresAt.IsZero() && time.Now().After(effect.ExpiresAt) {
			continue
		}
		
		entry := map[string]interface{}{
			"item_id":   effect.ItemID,
			"item_type": string(effect.ItemType),
			"value":     effect.Value,
		}
		if !effect.ExpiresAt.IsZero() {
			entry["expires_in"] = int(time.Until(effect.ExpiresAt).Seconds())
		}
		result = append(result, entry)
	}
	return result
}

// CleanupExpiredEffects 清理過期效果（定期呼叫）
func (s *QuestShop) CleanupExpiredEffects() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for _, state := range s.states {
		for key, effect := range state.ActiveEffects {
			if effect.Used || (!effect.ExpiresAt.IsZero() && time.Now().After(effect.ExpiresAt)) {
				delete(state.ActiveEffects, key)
			}
		}
	}
}
