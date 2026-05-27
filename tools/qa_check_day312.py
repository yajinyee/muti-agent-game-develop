"""
qa_check_day312.py — DAY-312 QA 驗證腳本
驗證 T166-T170 五個新 Lucky 魚系統的完整性
"""
import os
import sys

PASS = "✅"
FAIL = "❌"
results = []

def check(name, condition, detail=""):
    status = PASS if condition else FAIL
    results.append((status, name, detail))
    print(f"{status} {name}" + (f" — {detail}" if detail else ""))

# ── Server 檔案驗證 ──────────────────────────────────────────
server_game = r"d:\Kiro\server\internal\game"
server_proto = r"d:\Kiro\server\internal\protocol\messages.go"
server_tables = r"d:\Kiro\server\internal\data\tables.go"

# Handler 檔案
for tid, name in [("T166", "star_portal"), ("T167", "dragon_soul"),
                  ("T168", "spacetime_rift"), ("T169", "holy_judgment"), ("T170", "big_bang")]:
    fname = f"lucky_{name}_handler.go"
    path = os.path.join(server_game, fname)
    check(f"Server handler {fname}", os.path.exists(path))

# Protocol 訊息
with open(server_proto, encoding="utf-8") as f:
    proto_content = f.read()
for msg in ["lucky_star_portal", "lucky_dragon_soul", "lucky_spacetime_rift",
            "lucky_holy_judgment", "lucky_big_bang"]:
    check(f"Protocol msg {msg}", f'"{msg}"' in proto_content)

# Tables 目標物
with open(server_tables, encoding="utf-8") as f:
    tables_content = f.read()
for tid in ["T166", "T167", "T168", "T169", "T170"]:
    check(f"Tables target {tid}", f'ID: "{tid}"' in tables_content)

# game.go 整合
game_go = os.path.join(server_game, "game.go")
with open(game_go, encoding="utf-8") as f:
    game_content = f.read()
for name in ["luckyStarPortal", "luckyDragonSoul", "luckySpacetimeRift",
             "luckyHolyJudgment", "luckyBigBang"]:
    check(f"game.go manager {name}", name in game_content)
for func in ["tryLuckyStarPortalFish", "tryLuckyDragonSoulFish",
             "tryLuckySpacetimeRiftFish", "tryLuckyHolyJudgmentFish", "tryLuckyBigBangFish"]:
    check(f"game.go trigger {func}", func in game_content)

# ── Client 檔案驗證 ──────────────────────────────────────────
client_ui = r"d:\Kiro\client\chiikawa-pixel\scripts\ui"
client_game = r"d:\Kiro\client\chiikawa-pixel\scripts\game"

# Lucky Panel 腳本
for panel in ["LuckyStarPortalPanel", "LuckyDragonSoulPanel", "LuckySpacetimeRiftPanel",
              "LuckyHolyJudgmentPanel", "LuckyBigBangPanel"]:
    path = os.path.join(client_ui, f"{panel}.gd")
    check(f"Client panel {panel}.gd", os.path.exists(path))

# GameManager 訊號
gm_path = os.path.join(client_game, "GameManager.gd")
with open(gm_path, encoding="utf-8") as f:
    gm_content = f.read()
for sig in ["lucky_star_portal", "lucky_dragon_soul", "lucky_spacetime_rift",
            "lucky_holy_judgment", "lucky_big_bang"]:
    check(f"GameManager signal {sig}", f"signal {sig}" in gm_content)
    check(f"GameManager emit {sig}", f'emit_signal("{sig}"' in gm_content)

# HUD.gd 事件處理
hud_path = os.path.join(client_ui, "HUD.gd")
with open(hud_path, encoding="utf-8") as f:
    hud_content = f.read()
for handler in ["_on_lucky_star_portal", "_on_lucky_dragon_soul", "_on_lucky_spacetime_rift",
                "_on_lucky_holy_judgment", "_on_lucky_big_bang"]:
    check(f"HUD handler {handler}", handler in hud_content)

# TargetManager 映射
tm_path = os.path.join(client_game, "TargetManager.gd")
with open(tm_path, encoding="utf-8") as f:
    tm_content = f.read()
for tid, name in [("T166", "star_portal"), ("T167", "dragon_soul"),
                  ("T168", "spacetime_rift"), ("T169", "holy_judgment"), ("T170", "big_bang")]:
    check(f"TargetManager sprite {tid}", f'"{tid}"' in tm_content and name in tm_content)

# ── 精靈圖驗證 ──────────────────────────────────────────────
sprites_dir = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
for tid, name in [("T166", "star_portal"), ("T167", "dragon_soul"),
                  ("T168", "spacetime_rift"), ("T169", "holy_judgment"), ("T170", "big_bang")]:
    fname = f"{tid}_{name}.png"
    path = os.path.join(sprites_dir, fname)
    if os.path.exists(path):
        size = os.path.getsize(path)
        check(f"Sprite {fname}", size > 1000, f"{size} bytes")
    else:
        check(f"Sprite {fname}", False, "file not found")

# ── 結果統計 ──────────────────────────────────────────────────
total = len(results)
passed = sum(1 for r in results if r[0] == PASS)
failed = total - passed
print(f"\n{'='*50}")
print(f"DAY-312 QA 結果：{passed}/{total} 通過")
if failed > 0:
    print(f"失敗項目：")
    for r in results:
        if r[0] == FAIL:
            print(f"  {r[0]} {r[1]}")
    sys.exit(1)
else:
    print("全部通過！✅")
