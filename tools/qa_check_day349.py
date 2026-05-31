#!/usr/bin/env python3
"""
qa_check_day349.py — DAY-349 QA 驗證腳本
成就系統 + 好友排行榜 + 同場排行榜
"""
import os
import sys
import subprocess

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SERVER_DIR = os.path.join(ROOT, "server")
CLIENT_DIR = os.path.join(ROOT, "client", "chiikawa-pixel")

passed = 0
failed = 0

def check(name, condition, detail=""):
    global passed, failed
    if condition:
        print(f"  ✅ {name}")
        passed += 1
    else:
        print(f"  ❌ {name}" + (f" — {detail}" if detail else ""))
        failed += 1

def file_exists(path):
    return os.path.isfile(path)

def file_contains(path, text):
    if not os.path.isfile(path):
        return False
    with open(path, "r", encoding="utf-8", errors="ignore") as f:
        return text in f.read()

print("=" * 60)
print("DAY-349 QA 驗證：成就系統 + 好友排行榜")
print("=" * 60)

# ── Server 編譯 ──────────────────────────────────────────────
print("\n[1] Server 編譯驗證")
result = subprocess.run(
    ["go", "build", "./..."],
    cwd=SERVER_DIR, capture_output=True, text=True
)
check("go build ./... 通過", result.returncode == 0, result.stderr[:200] if result.stderr else "")

result = subprocess.run(
    ["go", "vet", "./..."],
    cwd=SERVER_DIR, capture_output=True, text=True
)
check("go vet ./... 通過（零警告）", result.returncode == 0, result.stderr[:200] if result.stderr else "")

# ── Server 成就系統 ──────────────────────────────────────────
print("\n[2] Server 成就系統（achievement.go）")
ach_path = os.path.join(SERVER_DIR, "internal", "game", "achievement.go")
check("achievement.go 存在", file_exists(ach_path))
check("AchievementSystem 結構定義", file_contains(ach_path, "type AchievementSystem struct"))
check("25 個成就類型定義", file_contains(ach_path, "AchievTypeMegaRich"))
check("OnKill 方法", file_contains(ach_path, "func (a *AchievementSystem) OnKill"))
check("OnBossKill 方法", file_contains(ach_path, "func (a *AchievementSystem) OnBossKill"))
check("OnBonusComplete 方法", file_contains(ach_path, "func (a *AchievementSystem) OnBonusComplete"))
check("OnCombo 方法", file_contains(ach_path, "func (a *AchievementSystem) OnCombo"))
check("OnLuckyFish 方法", file_contains(ach_path, "func (a *AchievementSystem) OnLuckyFish"))
check("OnSeasonLevel 方法", file_contains(ach_path, "func (a *AchievementSystem) OnSeasonLevel"))
check("OnQuestComplete 方法", file_contains(ach_path, "func (a *AchievementSystem) OnQuestComplete"))
check("OnWeeklyComplete 方法", file_contains(ach_path, "func (a *AchievementSystem) OnWeeklyComplete"))
check("OnCoinsUpdate 方法", file_contains(ach_path, "func (a *AchievementSystem) OnCoinsUpdate"))
check("稀有度分類（legendary）", file_contains(ach_path, '"legendary"'))
check("成就獎勵定義（金幣）", file_contains(ach_path, "Reward      int"))

# ── Server 好友系統 ──────────────────────────────────────────
print("\n[3] Server 好友系統（friend.go）")
friend_path = os.path.join(SERVER_DIR, "internal", "game", "friend.go")
check("friend.go 存在", file_exists(friend_path))
check("FriendSystem 結構定義", file_contains(friend_path, "type FriendSystem struct"))
check("FriendEntry 結構定義", file_contains(friend_path, "type FriendEntry struct"))
check("UpdatePlayer 方法", file_contains(friend_path, "func (f *FriendSystem) UpdatePlayer"))
check("SetOffline 方法", file_contains(friend_path, "func (f *FriendSystem) SetOffline"))
check("GetRoomLeaderboard 方法", file_contains(friend_path, "func (f *FriendSystem) GetRoomLeaderboard"))
check("GetPlayerRank 方法", file_contains(friend_path, "func (f *FriendSystem) GetPlayerRank"))
check("GetOnlineCount 方法", file_contains(friend_path, "func (f *FriendSystem) GetOnlineCount"))

# ── Server 協定訊息 ──────────────────────────────────────────
print("\n[4] Server 協定訊息（messages.go）")
msg_path = os.path.join(SERVER_DIR, "internal", "protocol", "messages.go")
check("MsgAchievementUnlock 定義", file_contains(msg_path, "MsgAchievementUnlock"))
check("MsgAchievementList 定義", file_contains(msg_path, "MsgAchievementList"))
check("MsgRoomLeaderboard 定義", file_contains(msg_path, "MsgRoomLeaderboard"))
check("AchievementUnlockPayload 結構", file_contains(msg_path, "type AchievementUnlockPayload struct"))
check("AchievementListPayload 結構", file_contains(msg_path, "type AchievementListPayload struct"))
check("RoomLeaderboardPayload 結構", file_contains(msg_path, "type RoomLeaderboardPayload struct"))
check("RoomLeaderboardEntryPayload 結構", file_contains(msg_path, "type RoomLeaderboardEntryPayload struct"))
check("MsgAchievementListRequest 定義", file_contains(msg_path, "MsgAchievementListRequest"))
check("MsgRoomLeaderboardRequest 定義", file_contains(msg_path, "MsgRoomLeaderboardRequest"))

# ── Server game.go 整合 ──────────────────────────────────────
print("\n[5] Server game.go 整合")
game_path = os.path.join(SERVER_DIR, "internal", "game", "game.go")
check("achievementSystem 欄位", file_contains(game_path, "achievementSystem *AchievementSystem"))
check("friendSystem 欄位", file_contains(game_path, "friendSystem      *FriendSystem"))
check("newAchievementSystem() 初始化", file_contains(game_path, "newAchievementSystem()"))
check("newFriendSystem() 初始化", file_contains(game_path, "newFriendSystem()"))
check("notifyAchievements 函數", file_contains(game_path, "func (g *Game) notifyAchievements"))
check("handleAchievementListRequest 函數", file_contains(game_path, "func (g *Game) handleAchievementListRequest"))
check("handleRoomLeaderboardRequest 函數", file_contains(game_path, "func (g *Game) handleRoomLeaderboardRequest"))
check("updateFriendEntry 函數", file_contains(game_path, "func (g *Game) updateFriendEntry"))
check("擊破時觸發成就", file_contains(game_path, "achievementSystem.OnKill"))
check("Combo 時觸發成就", file_contains(game_path, "achievementSystem.OnCombo"))
check("Bonus 時觸發成就", file_contains(game_path, "achievementSystem.OnBonusComplete"))
check("BOSS 擊破時觸發成就", file_contains(game_path, "achievementSystem.OnBossKill"))
check("玩家離開時標記離線", file_contains(game_path, "friendSystem.SetOffline"))
check("訊息路由：achievement_list_request", file_contains(game_path, "MsgAchievementListRequest"))
check("訊息路由：room_leaderboard_request", file_contains(game_path, "MsgRoomLeaderboardRequest"))

# ── Client GameManager.gd ────────────────────────────────────
print("\n[6] Client GameManager.gd")
gm_path = os.path.join(CLIENT_DIR, "scripts", "game", "GameManager.gd")
check("achievement_unlock 訊號", file_contains(gm_path, "signal achievement_unlock"))
check("achievement_list 訊號", file_contains(gm_path, "signal achievement_list"))
check("room_leaderboard 訊號", file_contains(gm_path, "signal room_leaderboard"))
check("achievement_unlock 訊息處理", file_contains(gm_path, '"achievement_unlock"'))
check("achievement_list 訊息處理", file_contains(gm_path, '"achievement_list"'))
check("room_leaderboard 訊息處理", file_contains(gm_path, '"room_leaderboard"'))
check("request_achievement_list 函數", file_contains(gm_path, "func request_achievement_list"))
check("request_room_leaderboard 函數", file_contains(gm_path, "func request_room_leaderboard"))

# ── Client AchievementPanel.gd ───────────────────────────────
print("\n[7] Client AchievementPanel.gd")
ach_panel_path = os.path.join(CLIENT_DIR, "scripts", "ui", "AchievementPanel.gd")
check("AchievementPanel.gd 存在", file_exists(ach_panel_path))
check("extends CanvasLayer", file_contains(ach_panel_path, "extends CanvasLayer"))
check("layer = 24", file_contains(ach_panel_path, "layer = 24"))
check("成就解鎖通知彈出", file_contains(ach_panel_path, "_spawn_notify_popup"))
check("通知佇列機制", file_contains(ach_panel_path, "_notify_queue"))
check("稀有度顏色定義", file_contains(ach_panel_path, "RARITY_COLORS"))
check("成就列表更新", file_contains(ach_panel_path, "_on_achievement_list"))
check("show_panel 函數", file_contains(ach_panel_path, "func show_panel"))
check("右下角彈出動畫", file_contains(ach_panel_path, "1280 - NOTIFY_WIDTH"))

# ── Client RoomLeaderboardPanel.gd ──────────────────────────
print("\n[8] Client RoomLeaderboardPanel.gd")
room_lb_path = os.path.join(CLIENT_DIR, "scripts", "ui", "RoomLeaderboardPanel.gd")
check("RoomLeaderboardPanel.gd 存在", file_exists(room_lb_path))
check("extends CanvasLayer", file_contains(room_lb_path, "extends CanvasLayer"))
check("layer = 25", file_contains(room_lb_path, "layer = 25"))
check("排名獎牌顯示", file_contains(room_lb_path, "RANK_BADGES"))
check("在線人數顯示", file_contains(room_lb_path, "_online_label"))
check("我的排名顯示", file_contains(room_lb_path, "_my_rank_label"))
check("自己的行高亮", file_contains(room_lb_path, "is_me"))
check("show_panel 函數", file_contains(room_lb_path, "func show_panel"))
check("_on_room_leaderboard 函數", file_contains(room_lb_path, "_on_room_leaderboard"))

# ── Client HUD.gd 整合 ──────────────────────────────────────
print("\n[9] Client HUD.gd 整合")
hud_path = os.path.join(CLIENT_DIR, "scripts", "ui", "HUD.gd")
check("_create_achievement_button 函數", file_contains(hud_path, "_create_achievement_button"))
check("_create_room_leaderboard_button 函數", file_contains(hud_path, "_create_room_leaderboard_button"))
check("成就按鈕 🏅", file_contains(hud_path, '"🏅"'))
check("同場排行榜按鈕 👥", file_contains(hud_path, '"👥"'))
check("_ready 中呼叫 _create_achievement_button", file_contains(hud_path, "_create_achievement_button()"))
check("_ready 中呼叫 _create_room_leaderboard_button", file_contains(hud_path, "_create_room_leaderboard_button()"))

# ── 結果統計 ─────────────────────────────────────────────────
print("\n" + "=" * 60)
total = passed + failed
print(f"結果：{passed}/{total} 通過")
if failed == 0:
    print("🎉 全部通過！DAY-349 成就系統 + 好友排行榜 完成")
else:
    print(f"⚠️  {failed} 項未通過，請檢查上方錯誤")
print("=" * 60)

sys.exit(0 if failed == 0 else 1)
