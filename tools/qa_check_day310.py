"""
qa_check_day310.py — DAY-310 QA 驗證腳本
驗證 T161-T165 五個新 Lucky 魚系統的完整性
"""
import os, sys

PASS = "✅"
FAIL = "❌"
results = []

def check(name, condition, detail=""):
    status = PASS if condition else FAIL
    results.append((status, name, detail))
    return condition

# ── Server 端驗證 ──────────────────────────────────────────────
SERVER_GAME = r"d:\Kiro\server\internal\game"
SERVER_PROTOCOL = r"d:\Kiro\server\internal\protocol\messages.go"
SERVER_TABLES = r"d:\Kiro\server\internal\data\tables.go"
SERVER_GAME_GO = r"d:\Kiro\server\internal\game\game.go"

# Handler 檔案存在
for tid, name in [("T161", "combo_burst"), ("T162", "time_bomb"), ("T163", "elemental_fusion"),
                   ("T164", "treasure_hunter"), ("T165", "myth_awaken")]:
    path = os.path.join(SERVER_GAME, f"lucky_{name}_handler.go")
    check(f"Server handler: lucky_{name}_handler.go", os.path.exists(path))

# Protocol 訊息定義
with open(SERVER_PROTOCOL, encoding="utf-8") as f:
    proto_content = f.read()
for msg in ["lucky_combo_burst", "lucky_time_bomb", "lucky_elemental_fusion",
            "lucky_treasure_hunter", "lucky_myth_awaken"]:
    check(f"Protocol: {msg}", msg in proto_content)

# Tables 目標定義
with open(SERVER_TABLES, encoding="utf-8") as f:
    tables_content = f.read()
for tid in ["T161", "T162", "T163", "T164", "T165"]:
    check(f"Tables: {tid} 定義", f'ID: "{tid}"' in tables_content)

# game.go Manager 定義
with open(SERVER_GAME_GO, encoding="utf-8") as f:
    game_content = f.read()
for mgr in ["luckyComboBurst", "luckyTimeBomb", "luckyElementalFusion",
            "luckyTreasureHunter", "luckyMythAwaken"]:
    check(f"game.go: {mgr} manager", mgr in game_content)

# game.go 觸發邏輯
for fn in ["isLuckyComboBurstFish", "isLuckyTimeBombFish", "isLuckyElementalFusionFish",
           "isLuckyTreasureHunterFish", "isLuckyMythAwakenFish"]:
    check(f"game.go: {fn} 觸發", fn in game_content)

# ── Client 端驗證 ──────────────────────────────────────────────
CLIENT_UI = r"d:\Kiro\client\chiikawa-pixel\scripts\ui"
CLIENT_GAME = r"d:\Kiro\client\chiikawa-pixel\scripts\game"

# Panel 腳本存在
for name in ["LuckyComboBurstPanel", "LuckyTimeBombPanel", "LuckyElementalFusionPanel",
             "LuckyTreasureHunterPanel", "LuckyMythAwakenPanel"]:
    path = os.path.join(CLIENT_UI, f"{name}.gd")
    check(f"Client Panel: {name}.gd", os.path.exists(path))

# GameManager 訊號定義
gm_path = os.path.join(CLIENT_GAME, "GameManager.gd")
with open(gm_path, encoding="utf-8") as f:
    gm_content = f.read()
for sig in ["lucky_combo_burst", "lucky_time_bomb", "lucky_elemental_fusion",
            "lucky_treasure_hunter", "lucky_myth_awaken"]:
    check(f"GameManager: signal {sig}", f"signal {sig}" in gm_content)
    check(f"GameManager: emit {sig}", f'emit_signal("{sig}"' in gm_content)

# HUD 訊號連接
hud_path = os.path.join(CLIENT_UI, "HUD.gd")
with open(hud_path, encoding="utf-8") as f:
    hud_content = f.read()
for sig in ["lucky_combo_burst", "lucky_time_bomb", "lucky_elemental_fusion",
            "lucky_treasure_hunter", "lucky_myth_awaken"]:
    check(f"HUD: connect {sig}", f"GameManager.{sig}.connect" in hud_content)
    check(f"HUD: handler _on_{sig}", f"_on_{sig}" in hud_content)

# TargetManager 映射
tm_path = os.path.join(CLIENT_GAME, "TargetManager.gd")
with open(tm_path, encoding="utf-8") as f:
    tm_content = f.read()
for tid in ["T161", "T162", "T163", "T164", "T165"]:
    check(f"TargetManager: {tid} 映射", tid in tm_content)

# ── 美術資產驗證 ──────────────────────────────────────────────
SPRITES_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
for name in ["T161_combo_burst", "T162_time_bomb", "T163_elemental_fusion",
             "T164_treasure_hunter", "T165_myth_awaken"]:
    path = os.path.join(SPRITES_DIR, f"{name}.png")
    check(f"Sprite: {name}.png", os.path.exists(path) and os.path.getsize(path) > 100)

# ── Agent 文件驗證 ──────────────────────────────────────────────
check("Agent: combo-system-agent.md", os.path.exists(r"d:\Kiro\agents\combo-system-agent.md"))

# ── 結果輸出 ──────────────────────────────────────────────────
print("\n" + "="*60)
print("DAY-310 QA 驗證報告")
print("="*60)
passed = sum(1 for s, _, _ in results if s == PASS)
total = len(results)
for status, name, detail in results:
    print(f"{status} {name}" + (f" ({detail})" if detail else ""))
print("="*60)
print(f"結果：{passed}/{total} 通過")
if passed == total:
    print("🎉 全部通過！DAY-310 完成！")
else:
    print(f"⚠️ {total - passed} 項未通過，需要修復")
    sys.exit(1)
