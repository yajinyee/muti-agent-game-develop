#!/usr/bin/env python3
"""
QA Check DAY-317
- TargetManager T181-T190 Sprite 路徑補齊
- TargetManager T181-T190 備用顏色補齊
- Lucky Badge 圖示升級（T181+: 💫, T186+: 🌌）
- 視覺清晰度改善：高倍率字體大小、光暈大小
- Server BUILD OK + VET OK
"""

import os
import re
import subprocess
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
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

def read_file(path):
    try:
        with open(path, encoding="utf-8") as f:
            return f.read()
    except Exception as e:
        return ""

print("=" * 60)
print("QA Check DAY-317 — TargetManager T181-T190 + 視覺清晰度")
print("=" * 60)

# ── 1. Server 編譯 ────────────────────────────────────────────
print("\n[1] Server 編譯")
server_dir = os.path.join(ROOT, "server")
try:
    r = subprocess.run(["go", "build", "./..."], cwd=server_dir, capture_output=True, text=True, timeout=60)
    check("go build ./...", r.returncode == 0, r.stderr[:200] if r.returncode != 0 else "")
    r2 = subprocess.run(["go", "vet", "./..."], cwd=server_dir, capture_output=True, text=True, timeout=60)
    check("go vet ./...", r2.returncode == 0, r2.stderr[:200] if r2.returncode != 0 else "")
except Exception as e:
    check("go build/vet", False, str(e))

# ── 2. TargetManager T181-T190 Sprite 路徑 ───────────────────
print("\n[2] TargetManager T181-T190 Sprite 路徑")
tm_path = os.path.join(ROOT, "client", "chiikawa-pixel", "scripts", "game", "TargetManager.gd")
tm = read_file(tm_path)

for i in range(181, 191):
    tid = f"T{i}"
    names = {
        181: "mutation", 182: "arctic_storm", 183: "fisher_wild",
        184: "risk_level", 185: "cosmic_pulse", 186: "mirror_universe",
        187: "gravity_field", 188: "time_acceleration", 189: "nebula_vortex",
        190: "cosmic_judgment"
    }
    name = names.get(i, "unknown")
    expected = f'"{tid}": "res://assets/sprites/targets/{tid}_{name}.png"'
    check(f"{tid} Sprite 路徑", expected in tm)

# ── 3. TargetManager T181-T190 備用顏色 ──────────────────────
print("\n[3] TargetManager T181-T190 備用顏色")
for i in range(181, 191):
    tid = f"T{i}"
    check(f"{tid} 備用顏色", f'"{tid}": Color(' in tm)

# ── 4. Lucky Badge 圖示升級 ───────────────────────────────────
print("\n[4] Lucky Badge 圖示升級")
check("T186+ 宇宙圖示 🌌", '🌌' in tm, "T186+ 應有宇宙圖示")
check("T181+ 星圖示 💫", '💫' in tm, "T181+ 應有星圖示")
check("T171+ 老虎機 🎰", '🎰' in tm, "T171+ 應有老虎機圖示")
check("T106+ 閃光 ✨", '✨' in tm, "T106+ 應有閃光圖示")

# ── 5. 視覺清晰度改善 ─────────────────────────────────────────
print("\n[5] 視覺清晰度改善")
check("高倍率字體大小（font_size = 20）", "font_size = 20" in tm, "≥500x 應有更大字體")
check("高倍率字體大小（font_size = 17）", "font_size = 17" in tm, "≥100x 應有中等字體")
check("光暈大小依倍率調整（glow_size = 120）", "glow_size = 120.0" in tm, "≥1000x 應有更大光暈")
check("光暈大小依倍率調整（glow_size = 100）", "glow_size = 100.0" in tm, "≥500x 應有大光暈")
check("宇宙粉紅光暈（T186+）", "1.0, 0.0, 0.5" in tm, "超高倍率應有宇宙粉紅光暈")

# ── 6. Lucky Badge 光環顏色 T186+ ────────────────────────────
print("\n[6] Lucky Badge 光環顏色 T186+")
check("T186+ 宇宙粉紅光環", "1.0, 0.0, 0.5, 1.0" in tm, "T186+ 應有宇宙粉紅光環")
check("T181+ 最亮金光環", "0.95" in tm, "T181+ 應有最亮金光環")

# ── 7. 精靈圖檔案存在 ─────────────────────────────────────────
print("\n[7] 精靈圖檔案存在")
sprites_dir = os.path.join(ROOT, "client", "chiikawa-pixel", "assets", "sprites", "targets")
names_map = {
    181: "mutation", 182: "arctic_storm", 183: "fisher_wild",
    184: "risk_level", 185: "cosmic_pulse", 186: "mirror_universe",
    187: "gravity_field", 188: "time_acceleration", 189: "nebula_vortex",
    190: "cosmic_judgment"
}
for i in range(181, 191):
    fname = f"T{i}_{names_map[i]}.png"
    fpath = os.path.join(sprites_dir, fname)
    check(f"{fname} 存在", os.path.exists(fpath))

# ── 8. Server tables.go T181-T190 ────────────────────────────
print("\n[8] Server tables.go T181-T190")
tables_path = os.path.join(ROOT, "server", "internal", "data", "tables.go")
tables = read_file(tables_path)
for i in range(181, 191):
    check(f"T{i} 在 tables.go", f'"T{i}"' in tables)

# ── 9. Server game.go T181-T190 handler 整合 ─────────────────
print("\n[9] Server game.go T181-T190 handler 整合")
game_path = os.path.join(ROOT, "server", "internal", "game", "game.go")
game = read_file(game_path)
handlers = [
    "luckyMutation", "luckyArcticStorm", "luckyFisherWild",
    "luckyRiskLevel", "luckyCosmicPulse", "luckyMirrorUniverse",
    "luckyGravityField", "luckyTimeAcceleration", "luckyNebulaVortex",
    "luckyCosmicJudgment"
]
for h in handlers:
    check(f"{h} 整合", h in game)

# ── 10. GameManager.gd T181-T190 訊號 ────────────────────────
print("\n[10] GameManager.gd T181-T190 訊號")
gm_path = os.path.join(ROOT, "client", "chiikawa-pixel", "scripts", "game", "GameManager.gd")
gm = read_file(gm_path)
signals = [
    "lucky_mutation", "lucky_arctic_storm", "lucky_fisher_wild",
    "lucky_risk_level", "lucky_cosmic_pulse", "lucky_mirror_universe",
    "lucky_gravity_field", "lucky_time_acceleration", "lucky_nebula_vortex",
    "lucky_cosmic_judgment"
]
for s in signals:
    check(f"signal {s}", f"signal {s}" in gm)

# ── 結果 ──────────────────────────────────────────────────────
print("\n" + "=" * 60)
total = PASS + FAIL
print(f"結果：{PASS}/{total} 通過")
if FAIL == 0:
    print("🎉 全部通過！DAY-317 品質驗證完成")
else:
    print(f"⚠️  {FAIL} 項失敗，需要修正")
print("=" * 60)
sys.exit(0 if FAIL == 0 else 1)
