#!/usr/bin/env python3
"""
qa_check_day313.py — DAY-313 品質驗證腳本
Progressive Jackpot 系統 + LuckyPanelRegistry 架構重構
"""
import os
import subprocess
import sys

WORKSPACE = r"d:\Kiro"
SERVER_DIR = os.path.join(WORKSPACE, "server")
CLIENT_SCRIPTS = os.path.join(WORKSPACE, "client", "chiikawa-pixel", "scripts")
AGENTS_DIR = os.path.join(WORKSPACE, "agents")
TOOLS_DIR = os.path.join(WORKSPACE, "tools")
SPRITES_DIR = os.path.join(WORKSPACE, "client", "chiikawa-pixel", "assets", "sprites", "targets")

results = []
passed = 0
failed = 0

def check(name, condition, detail=""):
    global passed, failed
    status = "✅" if condition else "❌"
    results.append(f"{status} {name}" + (f" — {detail}" if detail else ""))
    if condition:
        passed += 1
    else:
        failed += 1

def file_exists(path):
    return os.path.isfile(path)

def file_contains(path, text):
    if not os.path.isfile(path):
        return False
    with open(path, "r", encoding="utf-8", errors="ignore") as f:
        return text in f.read()

def dir_file_count(path, ext=".go"):
    if not os.path.isdir(path):
        return 0
    return sum(1 for f in os.listdir(path) if f.endswith(ext))

print("=== DAY-313 QA 驗證 ===\n")

# ── Server 編譯 ──────────────────────────────────────────────
print("【Server 編譯】")
result = subprocess.run(
    ["go", "build", "./..."],
    cwd=SERVER_DIR, capture_output=True, text=True
)
check("go build ./... 通過", result.returncode == 0, result.stderr[:100] if result.returncode != 0 else "")

result = subprocess.run(
    ["go", "vet", "./..."],
    cwd=SERVER_DIR, capture_output=True, text=True
)
check("go vet ./... 通過", result.returncode == 0, result.stderr[:100] if result.returncode != 0 else "")

# ── Server 新增目標物 ─────────────────────────────────────────
print("\n【Server 目標物 T171-T175】")
tables_path = os.path.join(SERVER_DIR, "internal", "data", "tables.go")
for tid in ["T171", "T172", "T173", "T174", "T175"]:
    check(f"tables.go 包含 {tid}", file_contains(tables_path, f'"{tid}"'))

# ── Server Progressive Jackpot Handler ───────────────────────
print("\n【Server Progressive Jackpot Handler】")
handler_path = os.path.join(SERVER_DIR, "internal", "game", "lucky_jackpot_pool_handler.go")
check("lucky_jackpot_pool_handler.go 存在", file_exists(handler_path))
check("jackpotPool 結構定義", file_contains(handler_path, "type jackpotPool struct"))
check("四層 Jackpot 常數", file_contains(handler_path, "JackpotGrand"))
check("addContribution 方法", file_contains(handler_path, "func (jp *jackpotPool) addContribution"))
check("tryTrigger 方法", file_contains(handler_path, "func (jp *jackpotPool) tryTrigger"))
check("payout 方法", file_contains(handler_path, "func (jp *jackpotPool) payout"))
check("onTargetKill 方法", file_contains(handler_path, "func (m *luckyJackpotPoolManager) onTargetKill"))
check("T175 隨機觸發邏輯", file_contains(handler_path, '"T175"'))
check("weightedPickIndex 函數", file_contains(handler_path, "func weightedPickIndex"))

# ── Server Protocol ───────────────────────────────────────────
print("\n【Server Protocol】")
protocol_path = os.path.join(SERVER_DIR, "internal", "protocol", "messages.go")
check("MsgLuckyJackpotPool 定義", file_contains(protocol_path, "MsgLuckyJackpotPool"))

# ── Server game.go 整合 ───────────────────────────────────────
print("\n【Server game.go 整合】")
game_path = os.path.join(SERVER_DIR, "internal", "game", "game.go")
check("luckyJackpotPool 欄位", file_contains(game_path, "luckyJackpotPool *luckyJackpotPoolManager"))
check("newLuckyJackpotPoolManager 初始化", file_contains(game_path, "newLuckyJackpotPoolManager()"))
check("onShot 呼叫", file_contains(game_path, "g.luckyJackpotPool.onShot"))
check("isLuckyJackpotPoolFish 觸發", file_contains(game_path, "isLuckyJackpotPoolFish"))

# ── Client Lucky Panel ────────────────────────────────────────
print("\n【Client Lucky Panel（T171-T175）】")
ui_dir = os.path.join(CLIENT_SCRIPTS, "ui")
for panel in ["LuckyJackpotMiniPanel", "LuckyJackpotMinorPanel", "LuckyJackpotMajorPanel",
              "LuckyJackpotGrandPanel", "LuckyJackpotTriggerPanel"]:
    check(f"{panel}.gd 存在", file_exists(os.path.join(ui_dir, f"{panel}.gd")))

# ── Client LuckyPanelRegistry ─────────────────────────────────
print("\n【Client LuckyPanelRegistry 架構重構】")
registry_path = os.path.join(ui_dir, "LuckyPanelRegistry.gd")
check("LuckyPanelRegistry.gd 存在", file_exists(registry_path))
check("SIGNAL_TO_PANEL 映射表", file_contains(registry_path, "SIGNAL_TO_PANEL"))
check("connect_all_signals 方法", file_contains(registry_path, "func connect_all_signals"))
check("T171-T175 映射", file_contains(registry_path, "lucky_jackpot_mini"))

# ── Client GameManager 訊號 ───────────────────────────────────
print("\n【Client GameManager 訊號】")
gm_path = os.path.join(CLIENT_SCRIPTS, "game", "GameManager.gd")
check("lucky_jackpot_pool 訊號", file_contains(gm_path, "signal lucky_jackpot_pool"))

# ── Client TargetManager 映射 ─────────────────────────────────
print("\n【Client TargetManager 映射】")
tm_path = os.path.join(CLIENT_SCRIPTS, "game", "TargetManager.gd")
for tid in ["T171", "T172", "T173", "T174", "T175"]:
    check(f"TargetManager 包含 {tid}", file_contains(tm_path, f'"{tid}"'))

# ── Client HUD 訊號連接 ───────────────────────────────────────
print("\n【Client HUD 訊號連接】")
hud_path = os.path.join(ui_dir, "HUD.gd")
check("HUD 連接 lucky_jackpot_pool", file_contains(hud_path, "lucky_jackpot_pool.connect"))
check("HUD _on_lucky_jackpot_pool 函數", file_contains(hud_path, "_on_lucky_jackpot_pool"))

# ── 美術精靈圖 ────────────────────────────────────────────────
print("\n【美術精靈圖 T171-T175】")
for tid, name in [("T171", "jackpot_mini"), ("T172", "jackpot_minor"),
                  ("T173", "jackpot_major"), ("T174", "jackpot_grand"),
                  ("T175", "jackpot_trigger")]:
    sprite_path = os.path.join(SPRITES_DIR, f"{tid}_{name}.png")
    check(f"{tid} 精靈圖存在", file_exists(sprite_path))
    if file_exists(sprite_path):
        size = os.path.getsize(sprite_path)
        check(f"{tid} 精靈圖大小 > 500 bytes", size > 500, f"{size} bytes")

# ── Agent 文件 ────────────────────────────────────────────────
print("\n【Agent 文件】")
check("progressive-jackpot-agent.md 存在", 
      file_exists(os.path.join(AGENTS_DIR, "progressive-jackpot-agent.md")))

# ── 總計 ─────────────────────────────────────────────────────
print("\n" + "="*50)
print(f"結果：{passed}/{passed+failed} 通過")
print("="*50)
for r in results:
    print(r)

if failed > 0:
    print(f"\n❌ {failed} 項未通過")
    sys.exit(1)
else:
    print(f"\n✅ 全部 {passed} 項通過！DAY-313 品質驗證完成！")
