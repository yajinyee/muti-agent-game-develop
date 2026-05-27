"""
qa_check_day316.py — DAY-316 QA 驗證腳本
驗證 T186-T190 五個新 Lucky 魚系統的完整性
"""
import os
import sys

PASS = 0
FAIL = 0

def check(name, condition, detail=""):
    global PASS, FAIL
    if condition:
        print(f"  [PASS] {name}")
        PASS += 1
    else:
        print(f"  [FAIL] {name}" + (f" — {detail}" if detail else ""))
        FAIL += 1

def check_file(path, desc=""):
    check(desc or path, os.path.exists(path), f"File not found: {path}")

def check_content(path, keyword, desc=""):
    if not os.path.exists(path):
        check(desc or f"{path} contains '{keyword}'", False, "File not found")
        return
    with open(path, "r", encoding="utf-8", errors="ignore") as f:
        content = f.read()
    check(desc or f"{path} contains '{keyword}'", keyword in content)

ROOT = r"d:\Kiro"
SERVER = os.path.join(ROOT, "server")
CLIENT = os.path.join(ROOT, "client", "chiikawa-pixel")
SCRIPTS = os.path.join(CLIENT, "scripts")
TARGETS = os.path.join(CLIENT, "assets", "sprites", "targets")
TOOLS = os.path.join(ROOT, "tools")

print("=" * 60)
print("DAY-316 QA 驗證")
print("=" * 60)

# ── Server Handler 檔案 ──────────────────────────────────────
print("\n[1] Server Handler 檔案")
for fname in [
    "lucky_mirror_universe_handler.go",
    "lucky_gravity_field_handler.go",
    "lucky_time_acceleration_handler.go",
    "lucky_nebula_vortex_handler.go",
    "lucky_cosmic_judgment_handler.go",
]:
    check_file(os.path.join(SERVER, "internal", "game", fname), fname)

# ── Server tables.go ─────────────────────────────────────────
print("\n[2] Server tables.go 目標物定義")
tables_path = os.path.join(SERVER, "internal", "data", "tables.go")
for tid in ["T186", "T187", "T188", "T189", "T190"]:
    check_content(tables_path, f'ID: "{tid}"', f"tables.go 包含 {tid}")

# ── Server protocol/messages.go ──────────────────────────────
print("\n[3] Server protocol 訊息類型")
msg_path = os.path.join(SERVER, "internal", "protocol", "messages.go")
for msg in [
    "MsgLuckyMirrorUniverse",
    "MsgLuckyGravityField",
    "MsgLuckyTimeAcceleration",
    "MsgLuckyNebulaVortex",
    "MsgLuckyCosmicJudgment",
]:
    check_content(msg_path, msg, f"messages.go 包含 {msg}")

# ── Server game.go ───────────────────────────────────────────
print("\n[4] Server game.go 整合")
game_path = os.path.join(SERVER, "internal", "game", "game.go")
for keyword in [
    "luckyMirrorUniverse",
    "luckyGravityField",
    "luckyTimeAcceleration",
    "luckyNebulaVortex",
    "luckyCosmicJudgment",
    "newLuckyMirrorUniverseManager",
    "newLuckyGravityFieldManager",
    "newLuckyTimeAccelerationManager",
    "newLuckyNebulaVortexManager",
    "newLuckyCosmicJudgmentManager",
    "isLuckyMirrorUniverseFish",
    "isLuckyGravityFieldFish",
    "isLuckyTimeAccelerationFish",
    "isLuckyNebulaVortexFish",
    "isLuckyCosmicJudgmentFish",
    "getMirrorUniverseMult",
    "getGravityFieldMult",
    "getTimeAccelerationMult",
    "getNebulaVortexMult",
    "getCosmicJudgmentMult",
]:
    check_content(game_path, keyword, f"game.go 包含 {keyword}")

# ── Client Panel 腳本 ────────────────────────────────────────
print("\n[5] Client Lucky Panel 腳本")
ui_dir = os.path.join(SCRIPTS, "ui")
for fname in [
    "LuckyMirrorUniversePanel.gd",
    "LuckyGravityFieldPanel.gd",
    "LuckyTimeAccelerationPanel.gd",
    "LuckyNebulaVortexPanel.gd",
    "LuckyCosmicJudgmentPanel.gd",
]:
    check_file(os.path.join(ui_dir, fname), fname)

# ── Client GameManager.gd ────────────────────────────────────
print("\n[6] Client GameManager.gd 訊號")
gm_path = os.path.join(SCRIPTS, "game", "GameManager.gd")
for sig in [
    "lucky_mirror_universe",
    "lucky_gravity_field",
    "lucky_time_acceleration",
    "lucky_nebula_vortex",
    "lucky_cosmic_judgment",
]:
    check_content(gm_path, f"signal {sig}", f"GameManager.gd 包含 signal {sig}")
    check_content(gm_path, f'"{sig}"', f"GameManager.gd 包含 emit {sig}")

# ── Client HUD.gd ────────────────────────────────────────────
print("\n[7] Client HUD.gd 事件處理")
hud_path = os.path.join(ui_dir, "HUD.gd")
for handler in [
    "_on_lucky_mirror_universe",
    "_on_lucky_gravity_field",
    "_on_lucky_time_acceleration",
    "_on_lucky_nebula_vortex",
    "_on_lucky_cosmic_judgment",
]:
    check_content(hud_path, handler, f"HUD.gd 包含 {handler}")

# ── Client LuckyPanelRegistry.gd ────────────────────────────
print("\n[8] Client LuckyPanelRegistry.gd 映射")
registry_path = os.path.join(ui_dir, "LuckyPanelRegistry.gd")
for mapping in [
    "LuckyMirrorUniversePanel",
    "LuckyGravityFieldPanel",
    "LuckyTimeAccelerationPanel",
    "LuckyNebulaVortexPanel",
    "LuckyCosmicJudgmentPanel",
]:
    check_content(registry_path, mapping, f"Registry 包含 {mapping}")

# ── Client TargetManager.gd ──────────────────────────────────
print("\n[9] Client TargetManager.gd Lucky badge 範圍")
tm_path = os.path.join(SCRIPTS, "game", "TargetManager.gd")
check_content(tm_path, "tid_num <= 190", "TargetManager Lucky badge 擴展到 T190")

# ── 精靈圖 ───────────────────────────────────────────────────
print("\n[10] 精靈圖")
for fname in [
    "T186_mirror_universe.png",
    "T187_gravity_field.png",
    "T188_time_acceleration.png",
    "T189_nebula_vortex.png",
    "T190_cosmic_judgment.png",
]:
    check_file(os.path.join(TARGETS, fname), fname)

# ── 精靈圖密度 ───────────────────────────────────────────────
print("\n[11] 精靈圖密度（目標 > 25%）")
try:
    from PIL import Image
    for fname in [
        "T186_mirror_universe.png",
        "T187_gravity_field.png",
        "T188_time_acceleration.png",
        "T189_nebula_vortex.png",
        "T190_cosmic_judgment.png",
    ]:
        path = os.path.join(TARGETS, fname)
        if os.path.exists(path):
            img = Image.open(path).convert("RGBA")
            w, h = img.size
            count = sum(1 for y in range(h) for x in range(w) if img.getpixel((x, y))[3] > 10)
            pct = count / (w * h) * 100
            check(f"{fname} 密度 {pct:.1f}%", pct >= 25.0, f"密度 {pct:.1f}% < 25%")
except ImportError:
    print("  [SKIP] PIL 未安裝，跳過密度檢查")

# ── 結果 ─────────────────────────────────────────────────────
print("\n" + "=" * 60)
total = PASS + FAIL
print(f"結果：{PASS}/{total} 通過")
if FAIL == 0:
    print("✅ 全部通過！DAY-316 QA 完成")
else:
    print(f"❌ {FAIL} 項失敗，需要修復")
    sys.exit(1)
