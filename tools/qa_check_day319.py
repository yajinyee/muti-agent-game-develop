"""
qa_check_day319.py — DAY-319 QA 驗證腳本
驗證 T201-T205 五個新 Lucky 魚系統的完整性
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

BASE = r"d:\Kiro"
SERVER_GAME = os.path.join(BASE, "server", "internal", "game")
SERVER_DATA = os.path.join(BASE, "server", "internal", "data")
SERVER_PROTO = os.path.join(BASE, "server", "internal", "protocol")
CLIENT_UI = os.path.join(BASE, "client", "chiikawa-pixel", "scripts", "ui")
CLIENT_GAME = os.path.join(BASE, "client", "chiikawa-pixel", "scripts", "game")
SPRITES = os.path.join(BASE, "client", "chiikawa-pixel", "assets", "sprites", "targets")

print("=== DAY-319 QA 驗證 ===\n")

# ── Server Handler 檔案 ──────────────────────────────────────
print("【Server Handler 檔案】")
handlers = [
    "lucky_energy_storm_handler.go",
    "lucky_crystal_resonance_handler.go",
    "lucky_fate_judgment_handler.go",
    "lucky_time_reversal_handler.go",
    "lucky_cosmic_singularity_handler.go",
]
for h in handlers:
    path = os.path.join(SERVER_GAME, h)
    check(f"Handler 存在: {h}", os.path.exists(path))

# ── Server Handler 內容驗證 ──────────────────────────────────
print("\n【Server Handler 內容】")
handler_checks = [
    ("lucky_energy_storm_handler.go", ["luckyEnergyStormManager", "tryLuckyEnergyStormFish", "getEnergyStormMult", "T201"]),
    ("lucky_crystal_resonance_handler.go", ["luckyCrystalResonanceManager", "tryLuckyCrystalResonanceFish", "getCrystalResonanceMult", "T202"]),
    ("lucky_fate_judgment_handler.go", ["luckyFateJudgmentManager", "tryLuckyFateJudgmentFish", "getFateJudgmentMult", "T203"]),
    ("lucky_time_reversal_handler.go", ["luckyTimeReversalManager", "tryLuckyTimeReversalFish", "getTimeReversalMult", "T204"]),
    ("lucky_cosmic_singularity_handler.go", ["luckyCosmicSingularityManager", "tryLuckyCosmicSingularityFish", "getCosmicSingularityMult", "T205"]),
]
for filename, keywords in handler_checks:
    path = os.path.join(SERVER_GAME, filename)
    if os.path.exists(path):
        content = open(path, encoding="utf-8").read()
        for kw in keywords:
            check(f"  {filename} 包含 {kw}", kw in content)

# ── tables.go 驗證 ──────────────────────────────────────────
print("\n【tables.go 目標物定義】")
tables_path = os.path.join(SERVER_DATA, "tables.go")
tables_content = open(tables_path, encoding="utf-8").read()
for tid in ["T201", "T202", "T203", "T204", "T205"]:
    check(f"tables.go 包含 {tid}", tid in tables_content)
check("T205 倍率 8888", "8888" in tables_content)
check("T205 HP 8888", "HP: 8888" in tables_content)

# ── messages.go 驗證 ──────────────────────────────────────────
print("\n【messages.go 訊息定義】")
msg_path = os.path.join(SERVER_PROTO, "messages.go")
msg_content = open(msg_path, encoding="utf-8").read()
for msg in ["MsgLuckyEnergyStorm", "MsgLuckyCrystalResonance", "MsgLuckyFateJudgment", "MsgLuckyTimeReversal", "MsgLuckyCosmicSingularity"]:
    check(f"messages.go 包含 {msg}", msg in msg_content)

# ── game.go 整合驗證 ──────────────────────────────────────────
print("\n【game.go 整合】")
game_path = os.path.join(SERVER_GAME, "game.go")
game_content = open(game_path, encoding="utf-8").read()
for field in ["luckyEnergyStorm", "luckyCrystalResonance", "luckyFateJudgment", "luckyTimeReversal", "luckyCosmicSingularity"]:
    check(f"game.go struct 欄位: {field}", field in game_content)
for init in ["newLuckyEnergyStormManager", "newLuckyCrystalResonanceManager", "newLuckyFateJudgmentManager", "newLuckyTimeReversalManager", "newLuckyCosmicSingularityManager"]:
    check(f"game.go 初始化: {init}", init in game_content)
for handler in ["tryLuckyEnergyStormFish", "tryLuckyCrystalResonanceFish", "tryLuckyFateJudgmentFish", "tryLuckyTimeReversalFish", "tryLuckyCosmicSingularityFish"]:
    check(f"game.go handleKill: {handler}", handler in game_content)
for getter in ["getEnergyStormMult", "getCrystalResonanceMult", "getFateJudgmentMult", "getTimeReversalMult", "getCosmicSingularityMult"]:
    check(f"game.go 全服加成: {getter}", getter in game_content)

# ── Client Panel 檔案 ──────────────────────────────────────────
print("\n【Client Lucky Panel 檔案】")
panels = [
    "LuckyEnergyStormPanel.gd",
    "LuckyCrystalResonancePanel.gd",
    "LuckyFateJudgmentPanel.gd",
    "LuckyTimeReversalPanel.gd",
    "LuckyCosmicSingularityPanel.gd",
]
for p in panels:
    path = os.path.join(CLIENT_UI, p)
    check(f"Panel 存在: {p}", os.path.exists(path))

# ── GameManager.gd 訊號 ──────────────────────────────────────
print("\n【GameManager.gd 訊號】")
gm_path = os.path.join(CLIENT_GAME, "GameManager.gd")
gm_content = open(gm_path, encoding="utf-8").read()
for sig in ["lucky_energy_storm", "lucky_crystal_resonance", "lucky_fate_judgment", "lucky_time_reversal", "lucky_cosmic_singularity"]:
    check(f"GameManager 訊號: {sig}", f"signal {sig}" in gm_content)
    check(f"GameManager emit: {sig}", f'emit_signal("{sig}"' in gm_content)

# ── TargetManager.gd 映射 ──────────────────────────────────────
print("\n【TargetManager.gd 映射】")
tm_path = os.path.join(CLIENT_GAME, "TargetManager.gd")
tm_content = open(tm_path, encoding="utf-8").read()
for tid in ["T201", "T202", "T203", "T204", "T205"]:
    check(f"TargetManager Sprite 映射: {tid}", tid in tm_content)

# ── LuckyPanelRegistry.gd 映射 ──────────────────────────────────
print("\n【LuckyPanelRegistry.gd 映射】")
reg_path = os.path.join(CLIENT_UI, "LuckyPanelRegistry.gd")
reg_content = open(reg_path, encoding="utf-8").read()
for sig in ["lucky_energy_storm", "lucky_crystal_resonance", "lucky_fate_judgment", "lucky_time_reversal", "lucky_cosmic_singularity"]:
    check(f"Registry 映射: {sig}", sig in reg_content)

# ── 精靈圖驗證 ──────────────────────────────────────────────
print("\n【精靈圖驗證】")
sprite_files = [
    "T201_energy_storm.png",
    "T202_crystal_resonance.png",
    "T203_fate_judgment.png",
    "T204_time_reversal.png",
    "T205_cosmic_singularity.png",
]
for sf in sprite_files:
    path = os.path.join(SPRITES, sf)
    check(f"精靈圖存在: {sf}", os.path.exists(path))
    if os.path.exists(path):
        size = os.path.getsize(path)
        check(f"精靈圖大小 > 1KB: {sf}", size > 1024, f"實際: {size} bytes")

# ── 結果 ──────────────────────────────────────────────────────
print(f"\n=== 結果：{PASS} 通過 / {FAIL} 失敗 ===")
if FAIL == 0:
    print("✅ 全部通過！DAY-319 QA 完成")
else:
    print(f"❌ {FAIL} 項失敗，需要修復")
    sys.exit(1)
