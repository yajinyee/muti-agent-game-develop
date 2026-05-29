"""
qa_check_day331.py — DAY-331 品質驗證腳本
驗證 T239-T243 五個新 Lucky 魚系統的完整性
"""
import os
import re

ROOT = r"d:\Kiro"
SERVER_GAME = os.path.join(ROOT, "server", "internal", "game")
SERVER_DATA = os.path.join(ROOT, "server", "internal", "data")
SERVER_PROTO = os.path.join(ROOT, "server", "internal", "protocol")
CLIENT_UI = os.path.join(ROOT, "client", "chiikawa-pixel", "scripts", "ui")
CLIENT_GAME = os.path.join(ROOT, "client", "chiikawa-pixel", "scripts", "game")
CLIENT_SPRITES = os.path.join(ROOT, "client", "chiikawa-pixel", "assets", "sprites", "targets")
CLIENT_SCENES = os.path.join(ROOT, "client", "chiikawa-pixel", "scenes")

passed = 0
failed = 0

def check(name, condition, detail=""):
    global passed, failed
    if condition:
        passed += 1
        print(f"  ✅ {name}")
    else:
        failed += 1
        print(f"  ❌ {name}" + (f" — {detail}" if detail else ""))

def file_contains(path, pattern):
    try:
        with open(path, "r", encoding="utf-8", errors="ignore") as f:
            return pattern in f.read()
    except:
        return False

def file_exists(path):
    return os.path.exists(path)

targets = [
    ("T239", "shark_spark", "幸運鯊魚閃電魚", 138888),
    ("T240", "winter_ice", "幸運冬季冰釣魚", 148888),
    ("T241", "atlantis_frenzy", "幸運大西洋狂潮魚", 158888),
    ("T242", "fishing_time_wheel", "幸運釣魚時間魚", 168888),
    ("T243", "ultimate_shark", "幸運終極鯊魚魚", 188888),
]

print("=" * 60)
print("DAY-331 QA 驗證")
print("=" * 60)

# ── Server Handler 驗證 ──────────────────────────────────────
print("\n[Server Handler]")
for tid, name, _, _ in targets:
    handler = os.path.join(SERVER_GAME, f"lucky_{name}_handler.go")
    check(f"{tid} handler 存在", file_exists(handler))
    check(f"{tid} isLucky{name.replace('_', ' ').title().replace(' ', '')}Fish", 
          file_contains(handler, f'return defID == "{tid}"'))

# ── Server tables.go 驗證 ────────────────────────────────────
print("\n[Server tables.go]")
tables_path = os.path.join(SERVER_DATA, "tables.go")
for tid, _, name, mult in targets:
    check(f"{tid} 在 tables.go", file_contains(tables_path, f'ID: "{tid}"'))
    check(f"{tid} 倍率 {mult}", file_contains(tables_path, str(mult)))

# ── Server messages.go 驗證 ─────────────────────────────────
print("\n[Server messages.go]")
msg_path = os.path.join(SERVER_PROTO, "messages.go")
for tid, name, _, _ in targets:
    msg_const = f"lucky_{name}"
    check(f"{tid} 訊息常數", file_contains(msg_path, msg_const))

# ── Server game.go 整合驗證 ──────────────────────────────────
print("\n[Server game.go]")
game_path = os.path.join(SERVER_GAME, "game.go")
for tid, name, _, _ in targets:
    manager_name = f"lucky{name.replace('_', ' ').title().replace(' ', '')}"
    check(f"{tid} manager 宣告", file_contains(game_path, f"lucky{name.replace('_', ' ').title().replace(' ', '')}"))
    check(f"{tid} handleKill 整合", file_contains(game_path, f"isLucky{name.replace('_', ' ').title().replace(' ', '')}Fish"))

# ── Client Panel 驗證 ────────────────────────────────────────
print("\n[Client Lucky Panel]")
panel_names = {
    "T239": "LuckySharkSparkPanel",
    "T240": "LuckyWinterIcePanel",
    "T241": "LuckyAtlantisFrenzyPanel",
    "T242": "LuckyFishingTimeWheelPanel",
    "T243": "LuckyUltimateSharkPanel",
}
for tid, panel_name in panel_names.items():
    panel_path = os.path.join(CLIENT_UI, f"{panel_name}.gd")
    check(f"{tid} {panel_name} 存在", file_exists(panel_path))
    check(f"{tid} handle_event 函數", file_contains(panel_path, "func handle_event"))

# ── Client GameManager 訊號驗證 ──────────────────────────────
print("\n[Client GameManager]")
gm_path = os.path.join(CLIENT_GAME, "GameManager.gd")
for tid, name, _, _ in targets:
    signal_name = f"lucky_{name}"
    check(f"{tid} signal 定義", file_contains(gm_path, f"signal {signal_name}"))
    check(f"{tid} emit_signal", file_contains(gm_path, f'emit_signal("{signal_name}"'))

# ── Client LuckyPanelRegistry 驗證 ──────────────────────────
print("\n[Client LuckyPanelRegistry]")
registry_path = os.path.join(CLIENT_UI, "LuckyPanelRegistry.gd")
for tid, name, _, _ in targets:
    check(f"{tid} Registry 映射", file_contains(registry_path, f"lucky_{name}"))

# ── Client TargetManager 驗證 ────────────────────────────────
print("\n[Client TargetManager]")
tm_path = os.path.join(CLIENT_GAME, "TargetManager.gd")
for tid, name, _, _ in targets:
    check(f"{tid} Sprite 路徑", file_contains(tm_path, f'"{tid}"'))
    check(f"{tid} 備用顏色", file_contains(tm_path, f'"{tid}":'))

# ── 精靈圖驗證 ───────────────────────────────────────────────
print("\n[精靈圖]")
for tid, name, _, _ in targets:
    png_path = os.path.join(CLIENT_SPRITES, f"{tid}_{name}.png")
    import_path = png_path + ".import"
    check(f"{tid} PNG 存在", file_exists(png_path))
    check(f"{tid} .import 存在", file_exists(import_path))
    if file_exists(png_path):
        size = os.path.getsize(png_path)
        check(f"{tid} PNG 大小 > 0", size > 0, f"size={size}")

# ── Main.tscn 驗證 ───────────────────────────────────────────
print("\n[Main.tscn]")
tscn_path = os.path.join(CLIENT_SCENES, "Main.tscn")
for tid, _, _, _ in targets:
    panel_name = panel_names[tid]
    check(f"{tid} {panel_name} 在 Main.tscn", file_contains(tscn_path, panel_name))
check("Main.tscn load_steps=148", file_contains(tscn_path, "load_steps=148"))

# ── 結果 ─────────────────────────────────────────────────────
print("\n" + "=" * 60)
total = passed + failed
print(f"結果：{passed}/{total} 通過")
if failed == 0:
    print("✅ 全部通過！DAY-331 品質驗證完成")
else:
    print(f"❌ {failed} 項失敗，需要修復")
print("=" * 60)
