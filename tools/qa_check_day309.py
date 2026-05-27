"""
qa_check_day309.py — DAY-309 QA 驗證腳本
驗證 T156-T160 五個新 Lucky 魚系統的完整性
"""
import os
import subprocess
import sys

PASS = "✅"
FAIL = "❌"
results = []

def check(name, condition, detail=""):
    status = PASS if condition else FAIL
    results.append((status, name, detail))
    print(f"{status} {name}" + (f" — {detail}" if detail else ""))

# ── Server 編譯 ───────────────────────────────────────────────
print("\n=== Server 編譯驗證 ===")
r = subprocess.run(["go", "build", "./..."], cwd=r"d:\Kiro\server", capture_output=True, text=True)
check("go build ./...", r.returncode == 0, r.stderr[:100] if r.stderr else "OK")

r = subprocess.run(["go", "vet", "./..."], cwd=r"d:\Kiro\server", capture_output=True, text=True)
check("go vet ./...", r.returncode == 0, r.stderr[:100] if r.stderr else "OK")

# ── Server Handler 檔案 ───────────────────────────────────────
print("\n=== Server Handler 檔案 ===")
handlers = [
    "lucky_ice_phoenix_handler.go",
    "lucky_dragon_fury_handler.go",
    "lucky_mult_cascade_handler.go",
    "lucky_awaken_boss_v2_handler.go",
    "lucky_ultimate_judgment_handler.go",
]
for h in handlers:
    path = os.path.join(r"d:\Kiro\server\internal\game", h)
    check(f"Handler: {h}", os.path.exists(path))

# ── tables.go T156-T160 ───────────────────────────────────────
print("\n=== tables.go 目標物定義 ===")
with open(r"d:\Kiro\server\internal\data\tables.go", encoding="utf-8") as f:
    tables_content = f.read()
for tid in ["T156", "T157", "T158", "T159", "T160"]:
    check(f"tables.go 包含 {tid}", f'"{tid}"' in tables_content)

# ── protocol/messages.go 訊息類型 ─────────────────────────────
print("\n=== protocol/messages.go 訊息類型 ===")
with open(r"d:\Kiro\server\internal\protocol\messages.go", encoding="utf-8") as f:
    proto_content = f.read()
for msg in ["MsgLuckyIcePhoenix", "MsgLuckyDragonFury", "MsgLuckyMultCascade", "MsgLuckyAwakenBossV2", "MsgLuckyUltimateJudgment"]:
    check(f"protocol 包含 {msg}", msg in proto_content)

# ── game.go 整合 ──────────────────────────────────────────────
print("\n=== game.go 整合 ===")
with open(r"d:\Kiro\server\internal\game\game.go", encoding="utf-8") as f:
    game_content = f.read()
for field in ["luckyIcePhoenix", "luckyDragonFury", "luckyMultCascade", "luckyAwakenBossV2", "luckyUltimateJudgment"]:
    check(f"game.go 包含 {field}", field in game_content)
check("game.go 包含 applyUltimateJudgment", "applyUltimateJudgment" in game_content)
check("game.go DAY-309 switch case", "isLuckyIcePhoenixFish" in game_content)
check("game.go DAY-309 notify", "notifyDragonFuryKill" in game_content)

# ── Client Panel 檔案 ─────────────────────────────────────────
print("\n=== Client Panel 檔案 ===")
panels = [
    "LuckyIcePhoenixPanel.gd",
    "LuckyDragonFuryPanel.gd",
    "LuckyMultCascadePanel.gd",
    "LuckyAwakenBossV2Panel.gd",
    "LuckyUltimateJudgmentPanel.gd",
]
for panel in panels:
    path = os.path.join(r"d:\Kiro\client\chiikawa-pixel\scripts\ui", panel)
    check(f"Panel: {panel}", os.path.exists(path))

# ── GameManager.gd 訊號 ───────────────────────────────────────
print("\n=== GameManager.gd 訊號 ===")
with open(r"d:\Kiro\client\chiikawa-pixel\scripts\game\GameManager.gd", encoding="utf-8") as f:
    gm_content = f.read()
for sig in ["lucky_ice_phoenix", "lucky_dragon_fury", "lucky_mult_cascade", "lucky_awaken_boss_v2", "lucky_ultimate_judgment"]:
    check(f"GameManager 訊號: {sig}", f"signal {sig}" in gm_content)
    check(f"GameManager emit: {sig}", f'emit_signal("{sig}"' in gm_content)

# ── HUD.gd 事件處理 ───────────────────────────────────────────
print("\n=== HUD.gd 事件處理 ===")
with open(r"d:\Kiro\client\chiikawa-pixel\scripts\ui\HUD.gd", encoding="utf-8") as f:
    hud_content = f.read()
for handler in ["_on_lucky_ice_phoenix", "_on_lucky_dragon_fury", "_on_lucky_mult_cascade", "_on_lucky_awaken_boss_v2", "_on_lucky_ultimate_judgment"]:
    check(f"HUD 事件處理: {handler}", handler in hud_content)

# ── TargetManager.gd 映射 ─────────────────────────────────────
print("\n=== TargetManager.gd 映射 ===")
with open(r"d:\Kiro\client\chiikawa-pixel\scripts\game\TargetManager.gd", encoding="utf-8") as f:
    tm_content = f.read()
for tid in ["T156", "T157", "T158", "T159", "T160"]:
    check(f"TargetManager 映射: {tid}", f'"{tid}"' in tm_content)

# ── 精靈圖 ────────────────────────────────────────────────────
print("\n=== 精靈圖 ===")
sprites = [
    "T156_ice_phoenix.png",
    "T157_dragon_fury.png",
    "T158_mult_cascade.png",
    "T159_awaken_boss_v2.png",
    "T160_ultimate_judgment.png",
]
for sprite in sprites:
    path = os.path.join(r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets", sprite)
    exists = os.path.exists(path)
    size = os.path.getsize(path) if exists else 0
    check(f"精靈圖: {sprite}", exists and size > 500, f"{size} bytes")

# ── 統計 ──────────────────────────────────────────────────────
print("\n=== 統計 ===")
passed = sum(1 for r in results if r[0] == PASS)
total = len(results)
print(f"\n結果：{passed}/{total} 通過")
if passed == total:
    print("🎉 DAY-309 QA 全部通過！")
else:
    print("⚠️ 有項目未通過，請檢查")
    for r in results:
        if r[0] == FAIL:
            print(f"  {r[0]} {r[1]}: {r[2]}")
sys.exit(0 if passed == total else 1)
