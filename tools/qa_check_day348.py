#!/usr/bin/env python3
"""
QA 驗證腳本 — DAY-348
任務幣兌換商店 + 賽季排行榜
"""
import os
import re

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
PASS = 0
FAIL = 0

def check(name, condition, detail=""):
    global PASS, FAIL
    if condition:
        print(f"  ✅ {name}")
        PASS += 1
    else:
        print(f"  ❌ {name}" + (f" — {detail}" if detail else ""))
        FAIL += 1

def read(path):
    full = os.path.join(ROOT, path)
    if not os.path.exists(full):
        return ""
    with open(full, encoding="utf-8", errors="ignore") as f:
        return f.read()

print("=" * 60)
print("DAY-348 QA 驗證：任務幣兌換商店 + 賽季排行榜")
print("=" * 60)

# ── Server 端驗證 ──────────────────────────────────────────────
print("\n【Server — quest_shop.go】")
shop = read("server/internal/game/quest_shop.go")
check("QuestShop 結構定義", "type QuestShop struct" in shop)
check("ShopItem 結構定義", "type ShopItem struct" in shop)
check("ShopItemType 定義", "ShopItemType string" in shop)
check("BetBoost 類型", "ShopItemBetBoost" in shop)
check("CoinBonus 類型", "ShopItemCoinBonus" in shop)
check("XPBoost 類型", "ShopItemXPBoost" in shop)
check("LuckyCharm 類型", "ShopItemLuckyCharm" in shop)
check("AutoAmmo 類型", "ShopItemAutoAmmo" in shop)
check("9個商品定義", shop.count('"bet_boost_') + shop.count('"coin_') + shop.count('"xp_') + shop.count('"lucky_charm') + shop.count('"auto_ammo') >= 9)
check("Purchase 函數", "func (s *QuestShop) Purchase(" in shop)
check("ConsumeEffect 函數", "func (s *QuestShop) ConsumeEffect(" in shop)
check("HasActiveEffect 函數", "func (s *QuestShop) HasActiveEffect(" in shop)
check("GetActiveEffectsSummary 函數", "func (s *QuestShop) GetActiveEffectsSummary(" in shop)
check("NewQuestShop 函數", "func NewQuestShop()" in shop)

print("\n【Server — season_leaderboard.go】")
lb = read("server/internal/game/season_leaderboard.go")
check("SeasonLeaderboard 結構定義", "type SeasonLeaderboard struct" in lb)
check("LeaderboardEntry 結構定義", "type LeaderboardEntry struct" in lb)
check("UpdatePlayer 函數", "func (lb *SeasonLeaderboard) UpdatePlayer(" in lb)
check("GetTop 函數", "func (lb *SeasonLeaderboard) GetTop(" in lb)
check("GetPlayerRank 函數", "func (lb *SeasonLeaderboard) GetPlayerRank(" in lb)
check("GetSnapshot 函數", "func (lb *SeasonLeaderboard) GetSnapshot(" in lb)
check("Reset 函數", "func (lb *SeasonLeaderboard) Reset(" in lb)
check("排序邏輯（XP 降序）", "SeasonXP > players[j].SeasonXP" in lb)

print("\n【Server — protocol/messages.go】")
msgs = read("server/internal/protocol/messages.go")
check("MsgShopItems 定義", "MsgShopItems" in msgs)
check("MsgShopPurchaseResult 定義", "MsgShopPurchaseResult" in msgs)
check("MsgShopEffectUpdate 定義", "MsgShopEffectUpdate" in msgs)
check("MsgSeasonLeaderboard 定義", "MsgSeasonLeaderboard" in msgs)
check("MsgShopRequest 定義", "MsgShopRequest" in msgs)
check("MsgShopPurchase 定義", "MsgShopPurchase" in msgs)
check("ShopItemsPayload 結構", "type ShopItemsPayload struct" in msgs)
check("ShopPurchaseResultPayload 結構", "type ShopPurchaseResultPayload struct" in msgs)
check("SeasonLeaderboardPayload 結構", "type SeasonLeaderboardPayload struct" in msgs)
check("LeaderboardEntryPayload 結構", "type LeaderboardEntryPayload struct" in msgs)

print("\n【Server — game.go 整合】")
game = read("server/internal/game/game.go")
check("questShop 欄位", "questShop         *QuestShop" in game)
check("seasonLeaderboard 欄位", "seasonLeaderboard *SeasonLeaderboard" in game)
check("NewQuestShop 初始化", "questShop:         NewQuestShop()" in game)
check("NewSeasonLeaderboard 初始化", "seasonLeaderboard: NewSeasonLeaderboard(" in game)
check("MsgShopRequest 處理", "MsgShopRequest" in game)
check("MsgShopPurchase 處理", "MsgShopPurchase" in game)
check("MsgSeasonLeaderboardRequest 處理", "MsgSeasonLeaderboardRequest" in game)
check("handleShopRequest 函數", "func (g *Game) handleShopRequest(" in game)
check("handleShopPurchase 函數", "func (g *Game) handleShopPurchase(" in game)
check("handleSeasonLeaderboardRequest 函數", "func (g *Game) handleSeasonLeaderboardRequest(" in game)

print("\n【Server — daily_quest.go SpendQuestCoins】")
dq = read("server/internal/game/daily_quest.go")
check("SpendQuestCoins 函數", "func (dqs *DailyQuestSystem) SpendQuestCoins(" in dq)

print("\n【Server — weekly_challenge.go SpendQuestCoins】")
wc = read("server/internal/game/weekly_challenge.go")
check("SpendQuestCoins 函數", "func (wcs *WeeklyChallengeSystem) SpendQuestCoins(" in wc)

# ── Client 端驗證 ──────────────────────────────────────────────
print("\n【Client — QuestShopPanel.gd】")
shop_panel = read("client/chiikawa-pixel/scripts/ui/QuestShopPanel.gd")
check("QuestShopPanel 存在", len(shop_panel) > 0)
check("extends CanvasLayer", "extends CanvasLayer" in shop_panel)
check("layer = 22", "layer = 22" in shop_panel)
check("show_panel 函數", "func show_panel()" in shop_panel)
check("_on_shop_items_received 函數", "_on_shop_items_received" in shop_panel)
check("_on_purchase_result 函數", "_on_purchase_result" in shop_panel)
check("_create_item_row 函數", "_create_item_row" in shop_panel)
check("shop_request 訊息發送", '"shop_request"' in shop_panel)
check("shop_purchase 訊息發送", '"shop_purchase"' in shop_panel)

print("\n【Client — SeasonLeaderboardPanel.gd】")
lb_panel = read("client/chiikawa-pixel/scripts/ui/SeasonLeaderboardPanel.gd")
check("SeasonLeaderboardPanel 存在", len(lb_panel) > 0)
check("extends CanvasLayer", "extends CanvasLayer" in lb_panel)
check("layer = 23", "layer = 23" in lb_panel)
check("show_panel 函數", "func show_panel()" in lb_panel)
check("_on_leaderboard_received 函數", "_on_leaderboard_received" in lb_panel)
check("_create_entry_row 函數", "_create_entry_row" in lb_panel)
check("排名顏色（金銀銅）", "RANK_COLORS" in lb_panel)
check("season_leaderboard_request 訊息", '"season_leaderboard_request"' in lb_panel)

print("\n【Client — GameManager.gd 訊號】")
gm = read("client/chiikawa-pixel/scripts/game/GameManager.gd")
check("shop_items_received 訊號", "signal shop_items_received" in gm)
check("shop_purchase_result 訊號", "signal shop_purchase_result" in gm)
check("shop_effect_update 訊號", "signal shop_effect_update" in gm)
check("season_leaderboard_received 訊號", "signal season_leaderboard_received" in gm)
check("shop_items 訊息處理", '"shop_items"' in gm)
check("shop_purchase_result 訊息處理", '"shop_purchase_result"' in gm)
check("season_leaderboard 訊息處理", '"season_leaderboard"' in gm)

print("\n【Client — HUD.gd 按鈕】")
hud = read("client/chiikawa-pixel/scripts/ui/HUD.gd")
check("_create_shop_button 函數", "_create_shop_button" in hud)
check("_create_leaderboard_button 函數", "_create_leaderboard_button" in hud)
check("商店按鈕 🛒", '"🛒"' in hud)
check("排行榜按鈕 📊", '"📊"' in hud)
check("_ready 中呼叫 _create_shop_button", "_create_shop_button()" in hud)
check("_ready 中呼叫 _create_leaderboard_button", "_create_leaderboard_button()" in hud)

# ── 總結 ──────────────────────────────────────────────────────
print("\n" + "=" * 60)
total = PASS + FAIL
print(f"結果：{PASS}/{total} 通過")
if FAIL == 0:
    print("🎉 全部通過！DAY-348 任務幣兌換商店 + 賽季排行榜完成")
else:
    print(f"⚠️  {FAIL} 項未通過，請檢查上方錯誤")
print("=" * 60)
