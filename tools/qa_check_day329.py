#!/usr/bin/env python3
"""
qa_check_day329.py — DAY-329 QA 驗證腳本
驗證 T234-T238 五個新 Lucky 魚系統的完整性
"""
import os
import re

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

def read(path):
    try:
        with open(path, "r", encoding="utf-8") as f:
            return f.read()
    except:
        return ""

print("=== DAY-329 QA CHECK ===\n")

# ── Server Handler 檔案 ──────────────────────────────────────
print("【Server Handler 檔案】")
handlers = {
    "T234": "lucky_fever_boost_ultimate_handler.go",
    "T235": "lucky_rapid_riches_ultimate_handler.go",
    "T236": "lucky_ice_fishing_master_handler.go",
    "T237": "lucky_cosmic_miracle_handler.go",
    "T238": "lucky_genesis_ultimate_handler.go",
}
base = r"d:\Kiro\server\internal\game"
for tid, fname in handlers.items():
    path = os.path.join(base, fname)
    check(f"{tid} handler 存在", os.path.exists(path))
    if os.path.exists(path):
        content = read(path)
        check(f"{tid} 有 Manager struct", "Manager struct" in content or "Manager\n" in content or "Manager {" in content)
        check(f"{tid} 有 try 函數", f"try" in content)
        check(f"{tid} 有 get Mult 函數", "getMult" in content or "getCosmicMiracle" in content or "getGenesis" in content or "getFever" in content or "getRapid" in content or "getIce" in content)

# ── tables.go ────────────────────────────────────────────────
print("\n【tables.go 目標物定義】")
tables = read(r"d:\Kiro\server\internal\data\tables.go")
for tid in ["T234", "T235", "T236", "T237", "T238"]:
    check(f"{tid} 在 tables.go", f'"{tid}"' in tables)

# ── messages.go ──────────────────────────────────────────────
print("\n【messages.go 訊息常數】")
msgs = read(r"d:\Kiro\server\internal\protocol\messages.go")
for msg in ["MsgLuckyFeverBoostUltimate", "MsgLuckyRapidRichesUltimate", "MsgLuckyIceFishingMaster", "MsgLuckyCosmicMiracle", "MsgLuckyGenesisUltimate"]:
    check(f"{msg} 存在", msg in msgs)

# ── game.go ──────────────────────────────────────────────────
print("\n【game.go 整合】")
game = read(r"d:\Kiro\server\internal\game\game.go")
for field in ["luckyFeverBoostUltimate", "luckyRapidRichesUltimate", "luckyIceFishingMaster", "luckyCosmicMiracle", "luckyGenesisUltimate"]:
    check(f"{field} 欄位存在", field in game)
for init in ["newLuckyFeverBoostUltimateManager", "newLuckyRapidRichesUltimateManager", "newLuckyIceFishingMasterManager", "newLuckyCosmicMiracleManager", "newLuckyGenesisUltimateManager"]:
    check(f"{init} 初始化存在", init in game)
for case in ["isLuckyFeverBoostUltimateFish", "isLuckyRapidRichesUltimateFish", "isLuckyIceFishingMasterFish", "isLuckyCosmicMiracleFish", "isLuckyGenesisUltimateFish"]:
    check(f"{case} case 存在", case in game)

# ── Client Panel 檔案 ────────────────────────────────────────
print("\n【Client Lucky Panel 檔案】")
panels = {
    "T234": "LuckyFeverBoostUltimatePanel.gd",
    "T235": "LuckyRapidRichesUltimatePanel.gd",
    "T236": "LuckyIceFishingMasterPanel.gd",
    "T237": "LuckyCosmicMiraclePanel.gd",
    "T238": "LuckyGenesisUltimatePanel.gd",
}
ui_base = r"d:\Kiro\client\chiikawa-pixel\scripts\ui"
for tid, fname in panels.items():
    path = os.path.join(ui_base, fname)
    check(f"{tid} Panel 存在", os.path.exists(path))
    if os.path.exists(path):
        content = read(path)
        check(f"{tid} Panel 有 handle_event", "handle_event" in content)

# ── LuckyPanelRegistry.gd ────────────────────────────────────
print("\n【LuckyPanelRegistry.gd 映射】")
registry = read(r"d:\Kiro\client\chiikawa-pixel\scripts\ui\LuckyPanelRegistry.gd")
for sig in ["lucky_fever_boost_ultimate", "lucky_rapid_riches_ultimate", "lucky_ice_fishing_master", "lucky_cosmic_miracle", "lucky_genesis_ultimate"]:
    check(f"{sig} 在 Registry", sig in registry)

# ── GameManager.gd ───────────────────────────────────────────
print("\n【GameManager.gd 訊號】")
gm = read(r"d:\Kiro\client\chiikawa-pixel\scripts\game\GameManager.gd")
for sig in ["lucky_fever_boost_ultimate", "lucky_rapid_riches_ultimate", "lucky_ice_fishing_master", "lucky_cosmic_miracle", "lucky_genesis_ultimate"]:
    check(f"signal {sig} 存在", f"signal {sig}" in gm)
    check(f"emit {sig} 存在", f'emit_signal("{sig}"' in gm)

# ── TargetManager.gd ─────────────────────────────────────────
print("\n【TargetManager.gd 目標物映射】")
tm = read(r"d:\Kiro\client\chiikawa-pixel\scripts\game\TargetManager.gd")
for tid in ["T234", "T235", "T236", "T237", "T238"]:
    check(f"{tid} 精靈圖路徑存在", f'"{tid}"' in tm)

# ── 精靈圖檔案 ───────────────────────────────────────────────
print("\n【精靈圖檔案】")
sprites_base = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
sprites = {
    "T234": "T234_fever_boost_ultimate.png",
    "T235": "T235_rapid_riches_ultimate.png",
    "T236": "T236_ice_fishing_master.png",
    "T237": "T237_cosmic_miracle.png",
    "T238": "T238_genesis_ultimate.png",
}
for tid, fname in sprites.items():
    path = os.path.join(sprites_base, fname)
    check(f"{tid} 精靈圖存在", os.path.exists(path))
    if os.path.exists(path):
        size = os.path.getsize(path)
        check(f"{tid} 精靈圖大小 > 0", size > 0, f"size={size}")

# ── 結果 ─────────────────────────────────────────────────────
print(f"\n=== 結果：{PASS} 通過 / {FAIL} 失敗 ===")
if FAIL == 0:
    print("✅ DAY-329 QA 全部通過！")
else:
    print(f"❌ 有 {FAIL} 項需要修復")
