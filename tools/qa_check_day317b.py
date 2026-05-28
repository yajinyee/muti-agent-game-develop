#!/usr/bin/env python3
"""
QA Check DAY-317b — T191-T195 五個新 Lucky 魚系統
"""

import os
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
    except:
        return ""

print("=" * 60)
print("QA Check DAY-317b — T191-T195 五個新 Lucky 魚系統")
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

# ── 2. Server handler 檔案存在 ────────────────────────────────
print("\n[2] Server handler 檔案")
game_dir = os.path.join(ROOT, "server", "internal", "game")
handlers = {
    "T191": "lucky_pvp_battle_handler.go",
    "T192": "lucky_skill_chain_handler.go",
    "T193": "lucky_global_explosion_handler.go",
    "T194": "lucky_spacetime_fold_handler.go",
    "T195": "lucky_cosmic_end_handler.go",
}
for tid, fname in handlers.items():
    fpath = os.path.join(game_dir, fname)
    check(f"{tid} handler 存在 ({fname})", os.path.exists(fpath))

# ── 3. tables.go T191-T195 ────────────────────────────────────
print("\n[3] tables.go T191-T195")
tables = read_file(os.path.join(ROOT, "server", "internal", "data", "tables.go"))
for i in range(191, 196):
    check(f"T{i} 在 tables.go", f'"T{i}"' in tables)

# ── 4. protocol/messages.go 新訊息常數 ───────────────────────
print("\n[4] protocol/messages.go 新訊息常數")
messages = read_file(os.path.join(ROOT, "server", "internal", "protocol", "messages.go"))
msg_consts = [
    "MsgLuckyPvpBattle",
    "MsgLuckySkillChain",
    "MsgLuckyGlobalExplosion",
    "MsgLuckySpacetimeFold",
    "MsgLuckyCosmicEnd",
]
for c in msg_consts:
    check(f"{c} 存在", c in messages)

# ── 5. game.go 整合 ───────────────────────────────────────────
print("\n[5] game.go T191-T195 整合")
game = read_file(os.path.join(ROOT, "server", "internal", "game", "game.go"))
integrations = [
    "luckyPvpBattle",
    "luckySkillChain",
    "luckyGlobalExplosion",
    "luckySpacetimeFold",
    "luckyCosmicEnd",
]
for h in integrations:
    check(f"{h} 整合", h in game)

# ── 6. 精靈圖存在 ─────────────────────────────────────────────
print("\n[6] 精靈圖存在")
sprites_dir = os.path.join(ROOT, "client", "chiikawa-pixel", "assets", "sprites", "targets")
sprite_names = {
    191: "pvp_battle",
    192: "skill_chain",
    193: "global_explosion",
    194: "spacetime_fold",
    195: "cosmic_end",
}
for i, name in sprite_names.items():
    fname = f"T{i}_{name}.png"
    check(f"{fname} 存在", os.path.exists(os.path.join(sprites_dir, fname)))

# ── 7. GameManager.gd 新訊號 ─────────────────────────────────
print("\n[7] GameManager.gd 新訊號")
gm = read_file(os.path.join(ROOT, "client", "chiikawa-pixel", "scripts", "game", "GameManager.gd"))
signals = [
    "signal lucky_pvp_battle",
    "signal lucky_skill_chain",
    "signal lucky_global_explosion",
    "signal lucky_spacetime_fold",
    "signal lucky_cosmic_end",
]
for s in signals:
    check(f"{s}", s in gm)

# ── 8. TargetManager.gd T191-T195 ────────────────────────────
print("\n[8] TargetManager.gd T191-T195")
tm = read_file(os.path.join(ROOT, "client", "chiikawa-pixel", "scripts", "game", "TargetManager.gd"))
for i, name in sprite_names.items():
    check(f"T{i} Sprite 路徑", f'"T{i}": "res://assets/sprites/targets/T{i}_{name}.png"' in tm)
    check(f"T{i} 備用顏色", f'"T{i}": Color(' in tm)

# ── 9. LuckyPanelRegistry.gd 新 Panel 映射 ───────────────────
print("\n[9] LuckyPanelRegistry.gd 新 Panel 映射")
registry = read_file(os.path.join(ROOT, "client", "chiikawa-pixel", "scripts", "ui", "LuckyPanelRegistry.gd"))
panel_mappings = [
    '"lucky_pvp_battle"',
    '"lucky_skill_chain"',
    '"lucky_global_explosion"',
    '"lucky_spacetime_fold"',
    '"lucky_cosmic_end"',
]
for m in panel_mappings:
    check(f"{m} 映射", m in registry)

# ── 10. Client Panel 檔案存在 ─────────────────────────────────
print("\n[10] Client Panel 檔案")
ui_dir = os.path.join(ROOT, "client", "chiikawa-pixel", "scripts", "ui")
panels = [
    "LuckyPvpBattlePanel.gd",
    "LuckySkillChainPanel.gd",
    "LuckyGlobalExplosionPanel.gd",
    "LuckySpacetimeFoldPanel.gd",
    "LuckyCosmicEndPanel.gd",
]
for p in panels:
    check(f"{p} 存在", os.path.exists(os.path.join(ui_dir, p)))

# ── 結果 ──────────────────────────────────────────────────────
print("\n" + "=" * 60)
total = PASS + FAIL
print(f"結果：{PASS}/{total} 通過")
if FAIL == 0:
    print("🎉 全部通過！DAY-317b T191-T195 品質驗證完成")
else:
    print(f"⚠️  {FAIL} 項失敗，需要修正")
print("=" * 60)
sys.exit(0 if FAIL == 0 else 1)
