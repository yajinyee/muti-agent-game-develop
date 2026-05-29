"""
qa_check_day332.py — DAY-332 QA 驗證腳本
驗證 T244-T248 五個新 Lucky 魚系統的完整性
"""
import os
import re

ROOT = r"d:\Kiro"
SERVER = os.path.join(ROOT, "server")
CLIENT = os.path.join(ROOT, "client", "chiikawa-pixel")
TOOLS = os.path.join(ROOT, "tools")

errors = []
passed = 0

def check(condition, msg):
    global passed
    if condition:
        passed += 1
    else:
        errors.append(f"❌ {msg}")

def file_contains(path, pattern):
    try:
        with open(path, encoding="utf-8") as f:
            return pattern in f.read()
    except:
        return False

def file_exists(path):
    return os.path.exists(path)

print("=== DAY-332 QA 驗證 ===\n")

# ── Server Handler 檔案 ──────────────────────────────────────
handlers = [
    ("T244", "lucky_wild_collector_handler.go"),
    ("T245", "lucky_lightning_eel_ultra_handler.go"),
    ("T246", "lucky_domino_chain_handler.go"),
    ("T247", "lucky_immortal_boss_ultra_handler.go"),
    ("T248", "lucky_quad_fusion_handler.go"),
]
for tid, fname in handlers:
    path = os.path.join(SERVER, "internal", "game", fname)
    check(file_exists(path), f"Handler 存在: {fname}")
    check(file_contains(path, f'"{tid}"'), f"Handler 包含 {tid} ID: {fname}")
    check(file_contains(path, "tryLucky"), f"Handler 包含 tryLucky 函數: {fname}")
    check(file_contains(path, "globalBonus"), f"Handler 包含 globalBonus: {fname}")

# ── tables.go ────────────────────────────────────────────────
tables_path = os.path.join(SERVER, "internal", "data", "tables.go")
for tid in ["T244", "T245", "T246", "T247", "T248"]:
    check(file_contains(tables_path, f'ID: "{tid}"'), f"tables.go 包含 {tid}")

# ── messages.go ──────────────────────────────────────────────
msg_path = os.path.join(SERVER, "internal", "protocol", "messages.go")
msg_consts = [
    "MsgLuckyWildCollector",
    "MsgLuckyLightningEelUltra",
    "MsgLuckyDominoChain",
    "MsgLuckyImmortalBossUltra",
    "MsgLuckyQuadFusion",
]
for const in msg_consts:
    check(file_contains(msg_path, const), f"messages.go 包含 {const}")

payload_types = [
    "LuckyWildCollectorPayload",
    "LuckyLightningEelUltraPayload",
    "LuckyDominoChainPayload",
    "LuckyImmortalBossUltraPayload",
    "LuckyQuadFusionPayload",
]
for pt in payload_types:
    check(file_contains(msg_path, pt), f"messages.go 包含 {pt}")

# ── game.go ──────────────────────────────────────────────────
game_path = os.path.join(SERVER, "internal", "game", "game.go")
game_fields = [
    "luckyWildCollector",
    "luckyLightningEelUltra",
    "luckyDominoChain",
    "luckyImmortalBossUltra",
    "luckyQuadFusion",
]
for field in game_fields:
    check(file_contains(game_path, field), f"game.go 包含 field: {field}")

# ── Client Panel 腳本 ────────────────────────────────────────
panels = [
    "LuckyWildCollectorPanel.gd",
    "LuckyLightningEelUltraPanel.gd",
    "LuckyDominoChainPanel.gd",
    "LuckyImmortalBossUltraPanel.gd",
    "LuckyQuadFusionPanel.gd",
]
for panel in panels:
    path = os.path.join(CLIENT, "scripts", "ui", panel)
    check(file_exists(path), f"Panel 存在: {panel}")
    check(file_contains(path, "handle_event"), f"Panel 包含 handle_event: {panel}")
    check(file_contains(path, "BaseLuckyPanel"), f"Panel 使用 BaseLuckyPanel: {panel}")

# ── GameManager.gd 訊號 ──────────────────────────────────────
gm_path = os.path.join(CLIENT, "scripts", "game", "GameManager.gd")
signals = [
    "lucky_wild_collector",
    "lucky_lightning_eel_ultra",
    "lucky_domino_chain",
    "lucky_immortal_boss_ultra",
    "lucky_quad_fusion",
]
for sig in signals:
    check(file_contains(gm_path, f"signal {sig}"), f"GameManager 包含 signal {sig}")
    check(file_contains(gm_path, f'emit_signal("{sig}"'), f"GameManager emit {sig}")

# ── TargetManager.gd ────────────────────────────────────────
tm_path = os.path.join(CLIENT, "scripts", "game", "TargetManager.gd")
for tid in ["T244", "T245", "T246", "T247", "T248"]:
    check(file_contains(tm_path, f'"{tid}"'), f"TargetManager 包含 {tid}")

# ── LuckyPanelRegistry.gd ───────────────────────────────────
reg_path = os.path.join(CLIENT, "scripts", "ui", "LuckyPanelRegistry.gd")
for sig in signals:
    check(file_contains(reg_path, sig), f"LuckyPanelRegistry 包含 {sig}")

# ── 精靈圖 PNG + .import ────────────────────────────────────
sprites = [
    ("T244_wild_collector.png", "T244_wild_collector.png.import"),
    ("T245_lightning_eel_ultra.png", "T245_lightning_eel_ultra.png.import"),
    ("T246_domino_chain.png", "T246_domino_chain.png.import"),
    ("T247_immortal_boss_ultra.png", "T247_immortal_boss_ultra.png.import"),
    ("T248_quad_fusion.png", "T248_quad_fusion.png.import"),
]
sprites_dir = os.path.join(CLIENT, "assets", "sprites", "targets")
for png, imp in sprites:
    check(file_exists(os.path.join(sprites_dir, png)), f"精靈圖存在: {png}")
    check(file_exists(os.path.join(sprites_dir, imp)), f".import 存在: {imp}")

# ── Main.tscn ───────────────────────────────────────────────
tscn_path = os.path.join(CLIENT, "scenes", "Main.tscn")
check(file_contains(tscn_path, "load_steps=153"), "Main.tscn load_steps=153")
for panel in ["LuckyWildCollectorPanel", "LuckyLightningEelUltraPanel",
              "LuckyDominoChainPanel", "LuckyImmortalBossUltraPanel", "LuckyQuadFusionPanel"]:
    check(file_contains(tscn_path, panel), f"Main.tscn 包含 {panel}")

# ── 結果 ────────────────────────────────────────────────────
total = passed + len(errors)
print(f"\n{'='*50}")
print(f"通過: {passed}/{total}")
if errors:
    print(f"\n失敗項目：")
    for e in errors:
        print(f"  {e}")
else:
    print("✅ 全部通過！")
