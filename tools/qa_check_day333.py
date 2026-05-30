#!/usr/bin/env python3
"""
qa_check_day333.py — DAY-333 QA 驗證腳本
驗證 T249-T253 五個新 Lucky 魚系統的完整性
"""
import os
import re

ROOT = r"d:\Kiro"
SERVER_GAME = os.path.join(ROOT, "server", "internal", "game")
SERVER_DATA = os.path.join(ROOT, "server", "internal", "data")
SERVER_PROTO = os.path.join(ROOT, "server", "internal", "protocol")
CLIENT_GAME = os.path.join(ROOT, "client", "chiikawa-pixel", "scripts", "game")
CLIENT_UI = os.path.join(ROOT, "client", "chiikawa-pixel", "scripts", "ui")
CLIENT_SCENES = os.path.join(ROOT, "client", "chiikawa-pixel", "scenes")
SPRITES = os.path.join(ROOT, "client", "chiikawa-pixel", "assets", "sprites", "targets")

pass_count = 0
fail_count = 0

def check(name, condition, detail=""):
    global pass_count, fail_count
    if condition:
        print(f"  ✅ {name}")
        pass_count += 1
    else:
        print(f"  ❌ {name}" + (f" — {detail}" if detail else ""))
        fail_count += 1

def file_contains(path, pattern):
    try:
        with open(path, "r", encoding="utf-8") as f:
            return pattern in f.read()
    except:
        return False

def file_exists(path):
    return os.path.exists(path)

print("=" * 60)
print("DAY-333 QA 驗證")
print("=" * 60)

# ── Server Handler 檔案 ──────────────────────────────────────
print("\n[1] Server Handler 檔案")
handlers = [
    ("lucky_electrical_frame_handler.go", "T249"),
    ("lucky_magnetic_respin_handler.go",  "T250"),
    ("lucky_fisherman_trail_handler.go",  "T251"),
    ("lucky_golden_gills_handler.go",     "T252"),
    ("lucky_penta_fusion_handler.go",     "T253"),
]
for filename, tid in handlers:
    path = os.path.join(SERVER_GAME, filename)
    check(f"{filename} 存在", file_exists(path))
    check(f"{filename} 包含 {tid}", file_contains(path, tid))

# ── Server tables.go ─────────────────────────────────────────
print("\n[2] Server tables.go")
tables_path = os.path.join(SERVER_DATA, "tables.go")
for tid in ["T249", "T250", "T251", "T252", "T253"]:
    check(f"tables.go 包含 {tid}", file_contains(tables_path, f'ID: "{tid}"'))

# ── Server protocol/messages.go ──────────────────────────────
print("\n[3] Server protocol/messages.go")
msg_path = os.path.join(SERVER_PROTO, "messages.go")
msg_consts = [
    "MsgLuckyElectricalFrame",
    "MsgLuckyMagneticRespin",
    "MsgLuckyFishermanTrail",
    "MsgLuckyGoldenGills",
    "MsgLuckyPentaFusion",
]
for const in msg_consts:
    check(f"messages.go 包含 {const}", file_contains(msg_path, const))

payload_types = [
    "LuckyElectricalFramePayload",
    "LuckyMagneticRespinPayload",
    "LuckyFishermanTrailPayload",
    "LuckyGoldenGillsPayload",
    "LuckyPentaFusionPayload",
]
for pt in payload_types:
    check(f"messages.go 包含 {pt}", file_contains(msg_path, pt))

# ── Server game.go ───────────────────────────────────────────
print("\n[4] Server game.go")
game_path = os.path.join(SERVER_GAME, "game.go")
game_checks = [
    "luckyElectricalFrame",
    "luckyMagneticRespin",
    "luckyFishermanTrail",
    "luckyGoldenGills",
    "luckyPentaFusion",
    "newLuckyElectricalFrameManager",
    "newLuckyMagneticRespinManager",
    "newLuckyFishermanTrailManager",
    "newLuckyGoldenGillsManager",
    "newLuckyPentaFusionManager",
    "isLuckyElectricalFrameFish",
    "isLuckyMagneticRespinFish",
    "isLuckyFishermanTrailFish",
    "isLuckyGoldenGillsFish",
    "isLuckyPentaFusionFish",
    "electricalFrameMult",
    "magneticRespinMult",
    "fishermanTrailMult",
    "goldenGillsMult",
    "pentaFusionMult",
]
for item in game_checks:
    check(f"game.go 包含 {item}", file_contains(game_path, item))

# ── Client GameManager.gd ────────────────────────────────────
print("\n[5] Client GameManager.gd")
gm_path = os.path.join(CLIENT_GAME, "GameManager.gd")
gm_signals = [
    "signal lucky_electrical_frame",
    "signal lucky_magnetic_respin",
    "signal lucky_fisherman_trail",
    "signal lucky_golden_gills",
    "signal lucky_penta_fusion",
    '"lucky_electrical_frame"',
    '"lucky_magnetic_respin"',
    '"lucky_fisherman_trail"',
    '"lucky_golden_gills"',
    '"lucky_penta_fusion"',
]
for sig in gm_signals:
    check(f"GameManager.gd 包含 {sig}", file_contains(gm_path, sig))

# ── Client TargetManager.gd ──────────────────────────────────
print("\n[6] Client TargetManager.gd")
tm_path = os.path.join(CLIENT_GAME, "TargetManager.gd")
for tid in ["T249", "T250", "T251", "T252", "T253"]:
    check(f"TargetManager.gd 包含 {tid} Sprite", file_contains(tm_path, f'"{tid}"'))

# ── Client Lucky Panel 腳本 ──────────────────────────────────
print("\n[7] Client Lucky Panel 腳本")
panels = [
    "LuckyElectricalFramePanel.gd",
    "LuckyMagneticRespinPanel.gd",
    "LuckyFishermanTrailPanel.gd",
    "LuckyGoldenGillsPanel.gd",
    "LuckyPentaFusionPanel.gd",
]
for panel in panels:
    path = os.path.join(CLIENT_UI, panel)
    check(f"{panel} 存在", file_exists(path))
    check(f"{panel} 包含 handle_event", file_contains(path, "handle_event"))

# ── Client LuckyPanelRegistry.gd ────────────────────────────
print("\n[8] Client LuckyPanelRegistry.gd")
reg_path = os.path.join(CLIENT_UI, "LuckyPanelRegistry.gd")
reg_checks = [
    "lucky_electrical_frame",
    "lucky_magnetic_respin",
    "lucky_fisherman_trail",
    "lucky_golden_gills",
    "lucky_penta_fusion",
    "LuckyElectricalFramePanel",
    "LuckyMagneticRespinPanel",
    "LuckyFishermanTrailPanel",
    "LuckyGoldenGillsPanel",
    "LuckyPentaFusionPanel",
]
for item in reg_checks:
    check(f"LuckyPanelRegistry.gd 包含 {item}", file_contains(reg_path, item))

# ── Client Main.tscn ─────────────────────────────────────────
print("\n[9] Client Main.tscn")
tscn_path = os.path.join(CLIENT_SCENES, "Main.tscn")
tscn_checks = [
    "load_steps=158",
    "LuckyElectricalFramePanel.gd",
    "LuckyMagneticRespinPanel.gd",
    "LuckyFishermanTrailPanel.gd",
    "LuckyGoldenGillsPanel.gd",
    "LuckyPentaFusionPanel.gd",
    'name="LuckyElectricalFramePanel"',
    'name="LuckyMagneticRespinPanel"',
    'name="LuckyFishermanTrailPanel"',
    'name="LuckyGoldenGillsPanel"',
    'name="LuckyPentaFusionPanel"',
]
for item in tscn_checks:
    check(f"Main.tscn 包含 {item}", file_contains(tscn_path, item))

# ── 美術精靈圖 ───────────────────────────────────────────────
print("\n[10] 美術精靈圖")
sprites = [
    "T249_electrical_frame.png",
    "T250_magnetic_respin.png",
    "T251_fisherman_trail.png",
    "T252_golden_gills.png",
    "T253_penta_fusion.png",
]
for sprite in sprites:
    path = os.path.join(SPRITES, sprite)
    check(f"{sprite} 存在", file_exists(path))
    import_path = path + ".import"
    check(f"{sprite}.import 存在", file_exists(import_path))

# ── 結果 ─────────────────────────────────────────────────────
print("\n" + "=" * 60)
total = pass_count + fail_count
print(f"結果：{pass_count}/{total} 通過")
if fail_count == 0:
    print("🎉 全部通過！DAY-333 T249-T253 驗證完成")
else:
    print(f"⚠️  {fail_count} 項失敗，需要修復")
print("=" * 60)
