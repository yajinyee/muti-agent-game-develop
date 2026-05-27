"""
qa_check_day314.py — DAY-314 QA 驗證腳本
驗證 T176-T180 五個新 Lucky 魚系統的完整性
"""
import os
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
    with open(path, 'r', encoding='utf-8', errors='ignore') as f:
        return text in f.read()

BASE = r"d:\Kiro"
SERVER = os.path.join(BASE, "server", "internal", "game")
PROTOCOL = os.path.join(BASE, "server", "internal", "protocol", "messages.go")
TABLES = os.path.join(BASE, "server", "internal", "data", "tables.go")
CLIENT_GAME = os.path.join(BASE, "client", "chiikawa-pixel", "scripts", "game")
CLIENT_UI = os.path.join(BASE, "client", "chiikawa-pixel", "scripts", "ui")
SPRITES = os.path.join(BASE, "client", "chiikawa-pixel", "assets", "sprites", "targets")

print("=" * 60)
print("DAY-314 QA 驗證")
print("=" * 60)

# ── Server Handler 檔案 ──────────────────────────────────────
print("\n[1] Server Handler 檔案")
for handler in ["lucky_multiverse_handler.go", "lucky_time_loop_handler.go",
                "lucky_fate_wheel_handler.go", "lucky_divine_realm_handler.go",
                "lucky_final_power_handler.go"]:
    check(handler, file_exists(os.path.join(SERVER, handler)))

# ── Protocol 訊息定義 ────────────────────────────────────────
print("\n[2] Protocol 訊息定義")
for msg in ["MsgLuckyMultiverse", "MsgLuckyTimeLoop", "MsgLuckyFateWheel",
            "MsgLuckyDivineRealm", "MsgLuckyFinalPower"]:
    check(msg, file_contains(PROTOCOL, msg))

# ── Tables 目標物定義 ────────────────────────────────────────
print("\n[3] Tables 目標物定義")
for tid in ["T176", "T177", "T178", "T179", "T180"]:
    check(f"{tid} 定義", file_contains(TABLES, f'ID: "{tid}"'))

# ── game.go 整合 ─────────────────────────────────────────────
print("\n[4] game.go 整合")
game_go = os.path.join(SERVER, "..", "..", "..", "internal", "game", "game.go")
game_go = os.path.join(BASE, "server", "internal", "game", "game.go")
for item in ["luckyMultiverse", "luckyTimeLoop", "luckyFateWheel",
             "luckyDivineRealm", "luckyFinalPower",
             "isLuckyMultiverseFish", "isLuckyTimeLoopFish",
             "isLuckyFateWheelFish", "isLuckyDivineRealmFish",
             "isLuckyFinalPowerFish"]:
    check(f"game.go: {item}", file_contains(game_go, item))

# ── Client Panel 腳本 ────────────────────────────────────────
print("\n[5] Client Lucky Panel 腳本")
for panel in ["LuckyMultiversePanel.gd", "LuckyTimeLoopPanel.gd",
              "LuckyFateWheelPanel.gd", "LuckyDivineRealmPanel.gd",
              "LuckyFinalPowerPanel.gd"]:
    check(panel, file_exists(os.path.join(CLIENT_UI, panel)))

# ── GameManager 訊號 ─────────────────────────────────────────
print("\n[6] GameManager 訊號")
gm = os.path.join(CLIENT_GAME, "GameManager.gd")
for sig in ["lucky_multiverse", "lucky_time_loop", "lucky_fate_wheel",
            "lucky_divine_realm", "lucky_final_power"]:
    check(f"signal {sig}", file_contains(gm, f"signal {sig}"))
    check(f"emit {sig}", file_contains(gm, f'emit_signal("{sig}"'))

# ── HUD 訊號連接 ─────────────────────────────────────────────
print("\n[7] HUD 訊號連接")
hud = os.path.join(CLIENT_UI, "HUD.gd")
for sig in ["lucky_multiverse", "lucky_time_loop", "lucky_fate_wheel",
            "lucky_divine_realm", "lucky_final_power"]:
    check(f"HUD connect {sig}", file_contains(hud, f"GameManager.{sig}.connect"))
    check(f"HUD handler {sig}", file_contains(hud, f"_on_{sig}"))

# ── TargetManager 映射 ───────────────────────────────────────
print("\n[8] TargetManager 映射")
tm = os.path.join(CLIENT_GAME, "TargetManager.gd")
for tid in ["T176", "T177", "T178", "T179", "T180"]:
    check(f"TargetManager {tid}", file_contains(tm, f'"{tid}"'))

# ── LuckyPanelRegistry 映射 ──────────────────────────────────
print("\n[9] LuckyPanelRegistry 映射")
reg = os.path.join(CLIENT_UI, "LuckyPanelRegistry.gd")
for sig in ["lucky_multiverse", "lucky_time_loop", "lucky_fate_wheel",
            "lucky_divine_realm", "lucky_final_power"]:
    check(f"Registry {sig}", file_contains(reg, sig))

# ── 精靈圖 ───────────────────────────────────────────────────
print("\n[10] 精靈圖")
for sprite in ["T176_multiverse.png", "T177_time_loop.png", "T178_fate_wheel.png",
               "T179_divine_realm.png", "T180_final_power.png"]:
    check(sprite, file_exists(os.path.join(SPRITES, sprite)))
    check(f"{sprite}.import", file_exists(os.path.join(SPRITES, sprite + ".import")))

# ── 結果 ─────────────────────────────────────────────────────
print("\n" + "=" * 60)
total = PASS + FAIL
print(f"結果：{PASS}/{total} 通過")
if FAIL == 0:
    print("✅ 全部通過！DAY-314 QA 完成")
else:
    print(f"❌ {FAIL} 項失敗，需要修復")
    sys.exit(1)
