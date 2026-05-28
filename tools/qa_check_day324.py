"""
qa_check_day324.py — DAY-324 品質驗證腳本
驗證 T211-T215 五個新 Lucky 系統的完整性
"""
import os
import sys

PASS = 0
FAIL = 0

def check(name, condition, detail=""):
    global PASS, FAIL
    if condition:
        PASS += 1
        print(f"  ✅ {name}")
    else:
        FAIL += 1
        print(f"  ❌ {name}" + (f" — {detail}" if detail else ""))

print("=== DAY-324 QA 驗證 ===\n")

# ── Server Handler 檔案 ──────────────────────────────────────
print("【Server Handler 檔案】")
handlers = [
    "lucky_avalanche_handler.go",
    "lucky_crash_multiplier_handler.go",
    "lucky_multiplier_ladder_handler.go",
    "lucky_ice_fishing_wheel_handler.go",
    "lucky_global_avalanche_handler.go",
]
for h in handlers:
    path = f"d:/Kiro/server/internal/game/{h}"
    check(f"Handler 存在: {h}", os.path.exists(path))

# ── Server tables.go ─────────────────────────────────────────
print("\n【Server tables.go 目標物定義】")
tables_path = "d:/Kiro/server/internal/data/tables.go"
with open(tables_path, encoding="utf-8") as f:
    tables_content = f.read()

for tid in ["T211", "T212", "T213", "T214", "T215"]:
    check(f"tables.go 包含 {tid}", f'"{tid}"' in tables_content)

check("T211 倍率 17000", "17000" in tables_content)
check("T212 倍率 18000", "18000" in tables_content)
check("T213 倍率 19000", "19000" in tables_content)
check("T214 倍率 20000", "20000" in tables_content)
check("T215 倍率 22000", "22000" in tables_content)

# ── Server messages.go ───────────────────────────────────────
print("\n【Server messages.go 訊息類型】")
msg_path = "d:/Kiro/server/internal/protocol/messages.go"
with open(msg_path, encoding="utf-8") as f:
    msg_content = f.read()

for msg in ["lucky_avalanche", "lucky_crash_multiplier", "lucky_multiplier_ladder",
            "lucky_ice_fishing_wheel", "lucky_global_avalanche"]:
    check(f"messages.go 包含 {msg}", msg in msg_content)

# ── Server game.go ───────────────────────────────────────────
print("\n【Server game.go 整合】")
game_path = "d:/Kiro/server/internal/game/game.go"
with open(game_path, encoding="utf-8") as f:
    game_content = f.read()

for mgr in ["luckyAvalanche", "luckyCrashMultiplier", "luckyMultiplierLadder",
            "luckyIceFishingWheel", "luckyGlobalAvalanche"]:
    check(f"game.go 包含 manager: {mgr}", mgr in game_content)

for fn in ["newLuckyAvalancheManager", "newLuckyCrashMultiplierManager",
           "newLuckyMultiplierLadderManager", "newLuckyIceFishingWheelManager",
           "newLuckyGlobalAvalancheManager"]:
    check(f"game.go 初始化: {fn}", fn in game_content)

for fn in ["isLuckyAvalancheFish", "isLuckyCrashMultiplierFish",
           "isLuckyMultiplierLadderFish", "isLuckyIceFishingWheelFish",
           "isLuckyGlobalAvalancheFish"]:
    check(f"game.go 觸發: {fn}", fn in game_content)

check("game.go 包含 luckyMultiplierLadder.onKill", "luckyMultiplierLadder.onKill" in game_content)

# ── Client Panel 檔案 ────────────────────────────────────────
print("\n【Client Panel 檔案】")
panels = [
    "LuckyAvalanchePanel.gd",
    "LuckyCrashMultiplierPanel.gd",
    "LuckyMultiplierLadderPanel.gd",
    "LuckyIceFishingWheelPanel.gd",
    "LuckyGlobalAvalanchePanel.gd",
]
for p in panels:
    path = f"d:/Kiro/client/chiikawa-pixel/scripts/ui/{p}"
    check(f"Panel 存在: {p}", os.path.exists(path))

# ── Client GameManager.gd ────────────────────────────────────
print("\n【Client GameManager.gd 訊號】")
gm_path = "d:/Kiro/client/chiikawa-pixel/scripts/game/GameManager.gd"
with open(gm_path, encoding="utf-8") as f:
    gm_content = f.read()

for sig in ["lucky_avalanche", "lucky_crash_multiplier", "lucky_multiplier_ladder",
            "lucky_ice_fishing_wheel", "lucky_global_avalanche"]:
    check(f"GameManager 訊號: {sig}", f"signal {sig}" in gm_content)
    check(f"GameManager emit: {sig}", f'emit_signal("{sig}"' in gm_content)

# ── Client LuckyPanelRegistry.gd ────────────────────────────
print("\n【Client LuckyPanelRegistry.gd】")
reg_path = "d:/Kiro/client/chiikawa-pixel/scripts/ui/LuckyPanelRegistry.gd"
with open(reg_path, encoding="utf-8") as f:
    reg_content = f.read()

for sig in ["lucky_avalanche", "lucky_crash_multiplier", "lucky_multiplier_ladder",
            "lucky_ice_fishing_wheel", "lucky_global_avalanche"]:
    check(f"Registry 包含: {sig}", sig in reg_content)

# ── Client TargetManager.gd ──────────────────────────────────
print("\n【Client TargetManager.gd】")
tm_path = "d:/Kiro/client/chiikawa-pixel/scripts/game/TargetManager.gd"
with open(tm_path, encoding="utf-8") as f:
    tm_content = f.read()

for tid in ["T211", "T212", "T213", "T214", "T215"]:
    check(f"TargetManager Sprite: {tid}", f'"{tid}"' in tm_content)

# ── 精靈圖檔案 ───────────────────────────────────────────────
print("\n【精靈圖檔案】")
sprites = [
    "T211_avalanche.png",
    "T212_crash_multiplier.png",
    "T213_multiplier_ladder.png",
    "T214_ice_fishing_wheel.png",
    "T215_global_avalanche.png",
]
for s in sprites:
    path = f"d:/Kiro/client/chiikawa-pixel/assets/sprites/targets/{s}"
    check(f"精靈圖存在: {s}", os.path.exists(path))

# ── 精靈圖密度 ───────────────────────────────────────────────
print("\n【精靈圖密度（目標 > 35%）】")
try:
    from PIL import Image
    for s in sprites:
        path = f"d:/Kiro/client/chiikawa-pixel/assets/sprites/targets/{s}"
        if os.path.exists(path):
            img = Image.open(path).convert("RGBA")
            w, h = img.size
            non_transparent = sum(1 for y in range(h) for x in range(w) if img.getpixel((x, y))[3] > 0)
            density = non_transparent / (w * h) * 100
            check(f"{s} 密度 {density:.1f}% > 35%", density > 35.0)
except ImportError:
    print("  ⚠️ PIL 未安裝，跳過密度檢查")

# ── 結果 ─────────────────────────────────────────────────────
print(f"\n=== 結果：{PASS} 通過 / {FAIL} 失敗 ===")
if FAIL == 0:
    print("🎉 全部通過！DAY-324 品質驗證完成！")
else:
    print(f"⚠️ {FAIL} 項目需要修復")
    sys.exit(1)
