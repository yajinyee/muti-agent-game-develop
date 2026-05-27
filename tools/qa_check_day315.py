"""
qa_check_day315.py — DAY-315 QA 驗證腳本
驗證 T181-T185 五個新 Lucky 魚系統的完整性
"""
import os
import subprocess
import sys

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

def file_exists(path):
    return os.path.exists(path)

def file_contains(path, text):
    if not os.path.exists(path):
        return False
    with open(path, "r", encoding="utf-8", errors="ignore") as f:
        return text in f.read()

BASE = r"d:\Kiro"
SERVER = os.path.join(BASE, "server")
CLIENT = os.path.join(BASE, "client", "chiikawa-pixel")
TOOLS = os.path.join(BASE, "tools")
TARGETS = os.path.join(CLIENT, "assets", "sprites", "targets")
SCRIPTS_UI = os.path.join(CLIENT, "scripts", "ui")
SCRIPTS_GAME = os.path.join(CLIENT, "scripts", "game")

print("=" * 60)
print("DAY-315 QA 驗證")
print("=" * 60)

# ── Server 編譯 ──────────────────────────────────────────────
print("\n[1] Server 編譯")
result = subprocess.run(["go", "build", "./..."], cwd=SERVER, capture_output=True, text=True)
check("go build ./...", result.returncode == 0, result.stderr[:200] if result.stderr else "")
result = subprocess.run(["go", "vet", "./..."], cwd=SERVER, capture_output=True, text=True)
check("go vet ./...", result.returncode == 0, result.stderr[:200] if result.stderr else "")

# ── Server Handler 檔案 ──────────────────────────────────────
print("\n[2] Server Handler 檔案")
handlers = [
    "lucky_mutation_handler.go",
    "lucky_arctic_storm_handler.go",
    "lucky_fisher_wild_handler.go",
    "lucky_risk_level_handler.go",
    "lucky_cosmic_pulse_handler.go",
]
for h in handlers:
    path = os.path.join(SERVER, "internal", "game", h)
    check(f"Handler: {h}", file_exists(path))

# ── Server tables.go 目標物定義 ──────────────────────────────
print("\n[3] Server tables.go 目標物定義")
tables_path = os.path.join(SERVER, "internal", "data", "tables.go")
for tid in ["T181", "T182", "T183", "T184", "T185"]:
    check(f"tables.go 包含 {tid}", file_contains(tables_path, f'"{tid}"'))

# ── Server protocol 訊息定義 ─────────────────────────────────
print("\n[4] Server protocol 訊息定義")
proto_path = os.path.join(SERVER, "internal", "protocol", "messages.go")
for msg in ["MsgLuckyMutation", "MsgLuckyArcticStorm", "MsgLuckyFisherWild", "MsgLuckyRiskLevel", "MsgLuckyCosmicPulse"]:
    check(f"protocol: {msg}", file_contains(proto_path, msg))

# ── Server game.go 整合 ──────────────────────────────────────
print("\n[5] Server game.go 整合")
game_path = os.path.join(SERVER, "internal", "game", "game.go")
for field in ["luckyMutation", "luckyArcticStorm", "luckyFisherWild", "luckyRiskLevel", "luckyCosmicPulse"]:
    check(f"game.go 欄位: {field}", file_contains(game_path, field))
for func_name in ["isLuckyMutationFish", "isLuckyArcticStormFish", "isLuckyFisherWildFish", "isLuckyRiskLevelFish", "isLuckyCosmicPulseFish"]:
    check(f"game.go 觸發: {func_name}", file_contains(game_path, func_name))

# ── Client 精靈圖 ────────────────────────────────────────────
print("\n[6] Client 精靈圖")
sprites = [
    "T181_mutation.png",
    "T182_arctic_storm.png",
    "T183_fisher_wild.png",
    "T184_risk_level.png",
    "T185_cosmic_pulse.png",
]
for s in sprites:
    path = os.path.join(TARGETS, s)
    check(f"精靈圖: {s}", file_exists(path))

# ── Client Lucky Panel 腳本 ──────────────────────────────────
print("\n[7] Client Lucky Panel 腳本")
panels = [
    "LuckyMutationPanel.gd",
    "LuckyArcticStormPanel.gd",
    "LuckyFisherWildPanel.gd",
    "LuckyRiskLevelPanel.gd",
    "LuckyCosmicPulsePanel.gd",
]
for p in panels:
    path = os.path.join(SCRIPTS_UI, p)
    check(f"Panel: {p}", file_exists(path))

# ── Client GameManager.gd 訊號 ───────────────────────────────
print("\n[8] Client GameManager.gd 訊號")
gm_path = os.path.join(SCRIPTS_GAME, "GameManager.gd")
for sig in ["lucky_mutation", "lucky_arctic_storm", "lucky_fisher_wild", "lucky_risk_level", "lucky_cosmic_pulse"]:
    check(f"GameManager 訊號: {sig}", file_contains(gm_path, f"signal {sig}"))
    check(f"GameManager 分發: {sig}", file_contains(gm_path, f'"{sig}"'))

# ── Client HUD.gd 訊號連接 ───────────────────────────────────
print("\n[9] Client HUD.gd 訊號連接")
hud_path = os.path.join(SCRIPTS_UI, "HUD.gd")
for sig in ["lucky_mutation", "lucky_arctic_storm", "lucky_fisher_wild", "lucky_risk_level", "lucky_cosmic_pulse"]:
    check(f"HUD 連接: {sig}", file_contains(hud_path, f"GameManager.{sig}.connect"))
    check(f"HUD 處理: _on_{sig}", file_contains(hud_path, f"func _on_{sig}"))

# ── Client LuckyPanelRegistry.gd ────────────────────────────
print("\n[10] Client LuckyPanelRegistry.gd")
registry_path = os.path.join(SCRIPTS_UI, "LuckyPanelRegistry.gd")
for panel in ["LuckyMutationPanel", "LuckyArcticStormPanel", "LuckyFisherWildPanel", "LuckyRiskLevelPanel", "LuckyCosmicPulsePanel"]:
    check(f"Registry: {panel}", file_contains(registry_path, panel))

# ── Client TargetManager.gd Lucky badge 範圍 ────────────────
print("\n[11] Client TargetManager.gd Lucky badge 範圍")
tm_path = os.path.join(SCRIPTS_GAME, "TargetManager.gd")
check("Lucky badge 範圍擴展到 T185", file_contains(tm_path, "tid_num <= 185"))
check("T181+ 最亮金色光環", file_contains(tm_path, "tid_num >= 181"))

# ── 結果統計 ─────────────────────────────────────────────────
print("\n" + "=" * 60)
total = PASS + FAIL
print(f"結果：{PASS}/{total} 通過")
if FAIL == 0:
    print("🎉 全部通過！DAY-315 QA 驗證完成！")
else:
    print(f"⚠️ {FAIL} 項失敗，請修復後重新驗證")
print("=" * 60)

sys.exit(0 if FAIL == 0 else 1)
