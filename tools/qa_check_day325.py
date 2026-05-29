"""
qa_check_day325.py — DAY-325 QA 驗證腳本
驗證 T216-T220 五個新 Lucky 魚系統的完整性
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
    with open(path, "r", encoding="utf-8") as f:
        return text in f.read()

print("=== DAY-325 QA 驗證 ===\n")

# ── 1. Server Handler 檔案 ──────────────────────────────────
print("【1】Server Handler 檔案")
handlers = [
    "lucky_fishing_net_handler.go",
    "lucky_tnt_bonus_handler.go",
    "lucky_disturbance_handler.go",
    "lucky_pearl_multiplier_handler.go",
    "lucky_rapid_riches_handler.go",
]
for h in handlers:
    path = f"d:/Kiro/server/internal/game/{h}"
    check(f"Handler 存在：{h}", file_exists(path))

# ── 2. Server tables.go ──────────────────────────────────────
print("\n【2】Server tables.go")
tables = "d:/Kiro/server/internal/data/tables.go"
for tid in ["T216", "T217", "T218", "T219", "T220"]:
    check(f"tables.go 包含 {tid}", file_contains(tables, f'ID: "{tid}"'))

# ── 3. Server messages.go ────────────────────────────────────
print("\n【3】Server messages.go")
messages = "d:/Kiro/server/internal/protocol/messages.go"
for msg in ["lucky_fishing_net", "lucky_tnt_bonus", "lucky_disturbance", "lucky_pearl_multiplier", "lucky_rapid_riches"]:
    check(f"messages.go 包含 {msg}", file_contains(messages, msg))

# ── 4. Server game.go ────────────────────────────────────────
print("\n【4】Server game.go")
game_go = "d:/Kiro/server/internal/game/game.go"
for field in ["luckyFishingNet", "luckyTNTBonus", "luckyDisturbance", "luckyPearlMultiplier", "luckyRapidRiches"]:
    check(f"game.go 包含 {field}", file_contains(game_go, field))
for func_name in ["isLuckyFishingNetFish", "isLuckyTNTBonusFish", "isLuckyDisturbanceFish", "isLuckyPearlMultiplierFish", "isLuckyRapidRichesFish"]:
    check(f"game.go dispatch {func_name}", file_contains(game_go, func_name))

# ── 5. Server player.go ──────────────────────────────────────
print("\n【5】Server player.go")
player_go = "d:/Kiro/server/internal/game/player.go"
check("player.go 包含 RecentKills", file_contains(player_go, "RecentKills"))
check("player.go 包含 AddRecentKill", file_contains(player_go, "AddRecentKill"))

# ── 6. Server 編譯 ───────────────────────────────────────────
print("\n【6】Server 編譯")
result = subprocess.run(["go", "build", "./..."], cwd="d:/Kiro/server", capture_output=True, text=True)
check("go build 通過", result.returncode == 0, result.stderr[:200] if result.stderr else "")
result = subprocess.run(["go", "vet", "./..."], cwd="d:/Kiro/server", capture_output=True, text=True)
check("go vet 通過", result.returncode == 0, result.stderr[:200] if result.stderr else "")

# ── 7. Client Lucky Panel 檔案 ──────────────────────────────
print("\n【7】Client Lucky Panel 檔案")
panels = [
    "LuckyFishingNetPanel.gd",
    "LuckyTNTBonusPanel.gd",
    "LuckyDisturbancePanel.gd",
    "LuckyPearlMultiplierPanel.gd",
    "LuckyRapidRichesPanel.gd",
]
for p in panels:
    path = f"d:/Kiro/client/chiikawa-pixel/scripts/ui/{p}"
    check(f"Panel 存在：{p}", file_exists(path))
    check(f"Panel 繼承 BaseLuckyPanel：{p}", file_contains(path, "extends BaseLuckyPanel"))
    check(f"Panel 有 handle_event：{p}", file_contains(path, "func handle_event"))

# ── 8. GameManager.gd 訊號 ──────────────────────────────────
print("\n【8】GameManager.gd 訊號")
gm = "d:/Kiro/client/chiikawa-pixel/scripts/game/GameManager.gd"
for sig in ["lucky_fishing_net", "lucky_tnt_bonus", "lucky_disturbance", "lucky_pearl_multiplier", "lucky_rapid_riches"]:
    check(f"GameManager 訊號：{sig}", file_contains(gm, f"signal {sig}"))
    check(f"GameManager emit：{sig}", file_contains(gm, f'emit_signal("{sig}"'))

# ── 9. TargetManager.gd 映射 ────────────────────────────────
print("\n【9】TargetManager.gd 映射")
tm = "d:/Kiro/client/chiikawa-pixel/scripts/game/TargetManager.gd"
for tid in ["T216", "T217", "T218", "T219", "T220"]:
    check(f"TargetManager Sprite 映射：{tid}", file_contains(tm, f'"{tid}"'))

# ── 10. LuckyPanelRegistry.gd ───────────────────────────────
print("\n【10】LuckyPanelRegistry.gd")
registry = "d:/Kiro/client/chiikawa-pixel/scripts/ui/LuckyPanelRegistry.gd"
for sig in ["lucky_fishing_net", "lucky_tnt_bonus", "lucky_disturbance", "lucky_pearl_multiplier", "lucky_rapid_riches"]:
    check(f"Registry 包含：{sig}", file_contains(registry, sig))

# ── 11. 精靈圖 ──────────────────────────────────────────────
print("\n【11】精靈圖")
sprites = [
    "T216_fishing_net.png",
    "T217_tnt_bonus.png",
    "T218_disturbance.png",
    "T219_pearl_multiplier.png",
    "T220_rapid_riches.png",
]
for s in sprites:
    path = f"d:/Kiro/client/chiikawa-pixel/assets/sprites/targets/{s}"
    check(f"精靈圖存在：{s}", file_exists(path))
    if os.path.exists(path):
        size = os.path.getsize(path)
        check(f"精靈圖大小 > 1KB：{s}", size > 1024, f"實際大小：{size} bytes")

# ── 12. 業界研究記錄 ────────────────────────────────────────
print("\n【12】業界研究記錄")
check("T216 業界依據（BGaming Fishing Club 2）", file_contains(tables, "BGaming"))
check("T217 業界依據（TNT Bonus）", file_contains(tables, "TNT"))
check("T218 業界依據（Fisch Disturbance）", file_contains(tables, "Fisch"))
check("T219 業界依據（Shark & Spark）", file_contains(tables, "Shark"))
check("T220 業界依據（Rapid Riches）", file_contains(tables, "Rapid Riches"))

# ── 結果 ────────────────────────────────────────────────────
print(f"\n{'='*40}")
print(f"結果：{PASS} 通過 / {FAIL} 失敗 / {PASS + FAIL} 總計")
if FAIL == 0:
    print("🎉 全部通過！DAY-325 品質驗證完成！")
else:
    print(f"⚠️  {FAIL} 項需要修復")
    sys.exit(1)
