"""
qa_check_day318.py — DAY-318 QA 驗證腳本
驗證 T196-T200 五個新 Lucky 魚系統的完整性
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

def check_file(path, desc=""):
    check(desc or path, os.path.exists(path), f"File not found: {path}")

def check_content(path, keyword, desc=""):
    if not os.path.exists(path):
        check(desc or keyword, False, f"File not found: {path}")
        return
    with open(path, encoding="utf-8", errors="ignore") as f:
        content = f.read()
    check(desc or keyword, keyword in content, f"Keyword '{keyword}' not in {path}")

print("=" * 60)
print("DAY-318 QA 驗證")
print("=" * 60)

# ── Server Handler 檔案 ──────────────────────────────────────
print("\n[1] Server Handler 檔案")
handlers = [
    ("lucky_dragon_king_handler.go", "T196 龍王輪盤"),
    ("lucky_eternal_cycle_handler.go", "T197 永恆循環"),
    ("lucky_chaos_explosion_handler.go", "T198 混沌爆炸"),
    ("lucky_divine_revival_handler.go", "T199 神聖復活"),
    ("lucky_genesis_epoch_handler.go", "T200 創世紀元"),
]
for fname, desc in handlers:
    check_file(f"server/internal/game/{fname}", f"Handler: {desc}")

# ── tables.go 目標物定義 ──────────────────────────────────────
print("\n[2] tables.go 目標物定義")
tables_path = "server/internal/data/tables.go"
for tid in ["T196", "T197", "T198", "T199", "T200"]:
    check_content(tables_path, f'"{tid}"', f"tables.go 包含 {tid}")

# ── messages.go 訊息類型 ──────────────────────────────────────
print("\n[3] messages.go 訊息類型")
msg_path = "server/internal/protocol/messages.go"
for msg in ["MsgLuckyDragonKing", "MsgLuckyEternalCycle", "MsgLuckyChaosExplosion",
            "MsgLuckyDivineRevival", "MsgLuckyGenesisEpoch"]:
    check_content(msg_path, msg, f"messages.go 包含 {msg}")

# ── game.go 整合 ──────────────────────────────────────────────
print("\n[4] game.go 整合")
game_path = "server/internal/game/game.go"
for field in ["luckyDragonKing", "luckyEternalCycle", "luckyChaosExplosion",
              "luckyDivineRevival", "luckyGenesisEpoch"]:
    check_content(game_path, field, f"game.go 包含 {field}")

for func_call in ["newLuckyDragonKingManager", "newLuckyEternalCycleManager",
                  "newLuckyChaosExplosionManager", "newLuckyDivineRevivalManager",
                  "newLuckyGenesisEpochManager"]:
    check_content(game_path, func_call, f"game.go 初始化 {func_call}")

for trigger in ["isLuckyDragonKingFish", "isLuckyEternalCycleFish",
                "isLuckyChaosExplosionFish", "isLuckyDivineRevivalFish",
                "isLuckyGenesisEpochFish"]:
    check_content(game_path, trigger, f"game.go 觸發 {trigger}")

# ── Client Panel 檔案 ─────────────────────────────────────────
print("\n[5] Client Panel 檔案")
panels = [
    ("LuckyDragonKingPanel.gd", "T196 龍王輪盤"),
    ("LuckyEternalCyclePanel.gd", "T197 永恆循環"),
    ("LuckyChaosExplosionPanel.gd", "T198 混沌爆炸"),
    ("LuckyDivineRevivalPanel.gd", "T199 神聖復活"),
    ("LuckyGenesisEpochPanel.gd", "T200 創世紀元"),
]
for fname, desc in panels:
    check_file(f"client/chiikawa-pixel/scripts/ui/{fname}", f"Panel: {desc}")

# ── GameManager.gd 訊號 ───────────────────────────────────────
print("\n[6] GameManager.gd 訊號")
gm_path = "client/chiikawa-pixel/scripts/game/GameManager.gd"
for sig in ["lucky_dragon_king", "lucky_eternal_cycle", "lucky_chaos_explosion",
            "lucky_divine_revival", "lucky_genesis_epoch"]:
    check_content(gm_path, f"signal {sig}", f"GameManager 訊號 {sig}")
    check_content(gm_path, f'"{sig}"', f"GameManager 處理 {sig}")

# ── LuckyPanelRegistry.gd 映射 ────────────────────────────────
print("\n[7] LuckyPanelRegistry.gd 映射")
reg_path = "client/chiikawa-pixel/scripts/ui/LuckyPanelRegistry.gd"
for sig in ["lucky_dragon_king", "lucky_eternal_cycle", "lucky_chaos_explosion",
            "lucky_divine_revival", "lucky_genesis_epoch"]:
    check_content(reg_path, sig, f"Registry 映射 {sig}")

# ── TargetManager.gd 更新 ─────────────────────────────────────
print("\n[8] TargetManager.gd 更新")
tm_path = "client/chiikawa-pixel/scripts/game/TargetManager.gd"
for tid, name in [("T196", "dragon_king"), ("T197", "eternal_cycle"),
                  ("T198", "chaos_explosion"), ("T199", "divine_revival"),
                  ("T200", "genesis_epoch")]:
    check_content(tm_path, f'"{tid}"', f"TargetManager Sprite {tid}")
    check_content(tm_path, name, f"TargetManager 包含 {name}")

check_content(tm_path, "tid_num <= 200", "TargetManager Lucky badge 範圍到 T200")

# ── 精靈圖檔案 ────────────────────────────────────────────────
print("\n[9] 精靈圖檔案")
sprites = [
    ("T196_dragon_king.png", "T196 龍王輪盤"),
    ("T197_eternal_cycle.png", "T197 永恆循環"),
    ("T198_chaos_explosion.png", "T198 混沌爆炸"),
    ("T199_divine_revival.png", "T199 神聖復活"),
    ("T200_genesis_epoch.png", "T200 創世紀元"),
]
for fname, desc in sprites:
    check_file(f"client/chiikawa-pixel/assets/sprites/targets/{fname}", f"Sprite: {desc}")

# ── 精靈圖尺寸驗證 ────────────────────────────────────────────
print("\n[10] 精靈圖尺寸驗證（應為 128x128）")
try:
    from PIL import Image
    for fname, desc in sprites:
        path = f"client/chiikawa-pixel/assets/sprites/targets/{fname}"
        if os.path.exists(path):
            img = Image.open(path)
            check(f"{desc} 尺寸 128x128", img.size == (128, 128), f"實際尺寸: {img.size}")
        else:
            check(f"{desc} 尺寸", False, "檔案不存在")
except ImportError:
    print("  ⚠️ PIL 未安裝，跳過尺寸驗證")

# ── 結果 ──────────────────────────────────────────────────────
print("\n" + "=" * 60)
total = PASS + FAIL
print(f"結果：{PASS}/{total} 通過")
if FAIL == 0:
    print("✅ 全部通過！DAY-318 QA 驗證完成")
else:
    print(f"❌ {FAIL} 項失敗，請修復後重新驗證")
    sys.exit(1)
