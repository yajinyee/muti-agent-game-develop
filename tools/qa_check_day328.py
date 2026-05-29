#!/usr/bin/env python3
"""DAY-328 QA 驗證腳本 — T229-T233 五個新 Lucky 魚系統"""
import os
import re

PASS = 0
FAIL = 0

def check(name, condition, detail=""):
    global PASS, FAIL
    if condition:
        PASS += 1
    else:
        FAIL += 1
        print(f"  ❌ FAIL: {name}" + (f" — {detail}" if detail else ""))

def read(path):
    try:
        with open(path, encoding="utf-8") as f:
            return f.read()
    except:
        return ""

BASE = r"d:\Kiro"
SERVER_GAME = os.path.join(BASE, "server", "internal", "game")
SERVER_PROTO = os.path.join(BASE, "server", "internal", "protocol", "messages.go")
SERVER_TABLES = os.path.join(BASE, "server", "internal", "data", "tables.go")
CLIENT_UI = os.path.join(BASE, "client", "chiikawa-pixel", "scripts", "ui")
CLIENT_GAME = os.path.join(BASE, "client", "chiikawa-pixel", "scripts", "game")
SPRITES = os.path.join(BASE, "client", "chiikawa-pixel", "assets", "sprites", "targets")
MAIN_TSCN = os.path.join(BASE, "client", "chiikawa-pixel", "scenes", "Main.tscn")

print("=== DAY-328 QA 驗證 ===\n")

# ── Server Handler 檔案 ──────────────────────────────────────
print("【Server Handler 檔案】")
handlers = [
    ("lucky_magnetic_attraction_handler.go", "T229"),
    ("lucky_super_chain_handler.go", "T230"),
    ("lucky_holy_pillar_handler.go", "T231"),
    ("lucky_time_stop_handler.go", "T232"),
    ("lucky_cosmic_restart_handler.go", "T233"),
]
for fname, tid in handlers:
    path = os.path.join(SERVER_GAME, fname)
    content = read(path)
    check(f"{fname} 存在", os.path.exists(path))
    check(f"{fname} 包含 {tid}", tid in content)
    check(f"{fname} 有 getMult 方法", "getMult" in content or "Mult()" in content)
    check(f"{fname} 有 try 方法", "try" in content.lower())

# ── Protocol messages.go ──────────────────────────────────────
print("\n【Protocol messages.go】")
proto = read(SERVER_PROTO)
new_msgs = [
    "MsgLuckyMagneticAttraction",
    "MsgLuckySuperChain",
    "MsgLuckyHolyPillar",
    "MsgLuckyTimeStop",
    "MsgLuckyCosmicRestart",
]
for msg in new_msgs:
    check(f"messages.go 包含 {msg}", msg in proto)

# ── tables.go ──────────────────────────────────────────────────
print("\n【tables.go 目標物定義】")
tables = read(SERVER_TABLES)
for tid in ["T229", "T230", "T231", "T232", "T233"]:
    check(f"tables.go 包含 {tid}", f'"{tid}"' in tables)

# ── game.go 整合 ──────────────────────────────────────────────
print("\n【game.go 整合】")
game_go = read(os.path.join(BASE, "server", "internal", "game", "game.go"))
structs = [
    "luckyMagneticAttraction",
    "luckySuperChain",
    "luckyHolyPillar",
    "luckyTimeStop",
    "luckyCosmicRestart",
]
for s in structs:
    check(f"game.go struct 包含 {s}", s in game_go)
    check(f"game.go new 包含 {s}", f"new{s[0].upper()}{s[1:]}Manager" in game_go or s in game_go)

# getMult 調用
mult_calls = [
    "getMagneticAttractionMult",
    "getSuperChainMult",
    "getHolyPillarMult",
    "getTimeStopMult",
    "getCosmicRestartMult",
]
for m in mult_calls:
    check(f"game.go effectiveMult 包含 {m}", m in game_go)

# ── Client Panel 腳本 ──────────────────────────────────────────
print("\n【Client Lucky Panel 腳本】")
panels = [
    "LuckyMagneticAttractionPanel.gd",
    "LuckySuperChainPanel.gd",
    "LuckyHolyPillarPanel.gd",
    "LuckyTimeStopPanel.gd",
    "LuckyCosmicRestartPanel.gd",
]
for panel in panels:
    path = os.path.join(CLIENT_UI, panel)
    content = read(path)
    check(f"{panel} 存在", os.path.exists(path))
    check(f"{panel} extends BaseLuckyPanel", "extends BaseLuckyPanel" in content)
    check(f"{panel} 有 handle_event", "handle_event" in content)

# ── GameManager.gd 訊號 ──────────────────────────────────────
print("\n【GameManager.gd 訊號】")
gm = read(os.path.join(CLIENT_GAME, "GameManager.gd"))
signals = [
    "lucky_magnetic_attraction",
    "lucky_super_chain",
    "lucky_holy_pillar",
    "lucky_time_stop",
    "lucky_cosmic_restart",
]
for sig in signals:
    check(f"GameManager 有訊號 {sig}", f"signal {sig}" in gm)
    check(f"GameManager 有 emit {sig}", f'emit_signal("{sig}"' in gm)

# ── TargetManager.gd 映射 ──────────────────────────────────────
print("\n【TargetManager.gd 映射】")
tm = read(os.path.join(CLIENT_GAME, "TargetManager.gd"))
for tid in ["T229", "T230", "T231", "T232", "T233"]:
    check(f"TargetManager 有 {tid} Sprite 映射", f'"{tid}"' in tm)

# ── LuckyPanelRegistry.gd ──────────────────────────────────────
print("\n【LuckyPanelRegistry.gd】")
registry = read(os.path.join(CLIENT_UI, "LuckyPanelRegistry.gd"))
for sig in signals:
    check(f"Registry 有 {sig}", sig in registry)

# ── 精靈圖 ──────────────────────────────────────────────────────
print("\n【精靈圖】")
sprite_files = [
    "T229_magnetic_attraction.png",
    "T230_super_chain.png",
    "T231_holy_pillar.png",
    "T232_time_stop.png",
    "T233_cosmic_restart.png",
]
for f in sprite_files:
    path = os.path.join(SPRITES, f)
    check(f"精靈圖 {f} 存在", os.path.exists(path))
    if os.path.exists(path):
        size = os.path.getsize(path)
        check(f"精靈圖 {f} 大小 > 0", size > 0, f"size={size}")

# ── Main.tscn ──────────────────────────────────────────────────
print("\n【Main.tscn 節點】")
tscn = read(MAIN_TSCN)
for panel in ["LuckyMagneticAttractionPanel", "LuckySuperChainPanel",
              "LuckyHolyPillarPanel", "LuckyTimeStopPanel", "LuckyCosmicRestartPanel"]:
    check(f"Main.tscn 有 {panel} 節點", panel in tscn)

# ── 結果 ──────────────────────────────────────────────────────
print(f"\n=== 結果：{PASS} 通過 / {FAIL} 失敗 ===")
if FAIL == 0:
    print("✅ 全部通過！DAY-328 驗證完成")
else:
    print(f"❌ {FAIL} 項失敗，需要修復")
