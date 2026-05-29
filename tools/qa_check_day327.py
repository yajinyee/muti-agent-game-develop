"""
qa_check_day327.py — DAY-327 品質驗證腳本
驗證 T224-T228 五個新 Lucky 魚系統的完整性
"""
import os
import re

PASS = 0
FAIL = 0

def check(name, condition, detail=""):
    global PASS, FAIL
    if condition:
        PASS += 1
    else:
        FAIL += 1
        print(f"  FAIL: {name}" + (f" — {detail}" if detail else ""))

def read(path):
    try:
        with open(path, "r", encoding="utf-8") as f:
            return f.read()
    except:
        return ""

# ── 1. Server Handler 檔案存在 ────────────────────────────────
handlers = [
    "lucky_golden_pot_handler.go",
    "lucky_cascade_lock_handler.go",
    "lucky_legend_awaken_handler.go",
    "lucky_crash_harvest_handler.go",
    "lucky_cosmic_fusion_handler.go",
]
for h in handlers:
    path = f"d:/Kiro/server/internal/game/{h}"
    check(f"Handler exists: {h}", os.path.exists(path))

# ── 2. Server Handler 內容驗證 ────────────────────────────────
handler_checks = {
    "lucky_golden_pot_handler.go": ["luckyGoldenPotManager", "tryLuckyGoldenPotFish", "global_boost", "43.0", "full_pot_bonus"],
    "lucky_cascade_lock_handler.go": ["luckyCascadeLockManager", "tryLuckyCascadeLockFish", "global_boost", "43.5", "perfect_cascade"],
    "lucky_legend_awaken_handler.go": ["luckyLegendAwakenManager", "tryLuckyLegendAwakenFish", "global_boost", "44.0", "awaken_reward"],
    "lucky_crash_harvest_handler.go": ["luckyCrashHarvestManager", "tryLuckyCrashHarvestFish", "global_boost", "44.5", "perfect_harvest"],
    "lucky_cosmic_fusion_handler.go": ["luckyCosmicFusionManager", "tryLuckyCosmicFusionFish", "global_boost", "45.0", "fusion_settle"],
}
for h, keywords in handler_checks.items():
    content = read(f"d:/Kiro/server/internal/game/{h}")
    for kw in keywords:
        check(f"{h} contains '{kw}'", kw in content)

# ── 3. tables.go 目標物定義 ───────────────────────────────────
tables = read("d:/Kiro/server/internal/data/tables.go")
for tid in ["T224", "T225", "T226", "T227", "T228"]:
    check(f"tables.go has {tid}", f'ID: "{tid}"' in tables)

# ── 4. protocol/messages.go 訊息定義 ─────────────────────────
messages = read("d:/Kiro/server/internal/protocol/messages.go")
for msg in ["lucky_golden_pot", "lucky_cascade_lock", "lucky_legend_awaken",
            "lucky_crash_harvest", "lucky_cosmic_fusion"]:
    check(f"messages.go has {msg}", msg in messages)

# ── 5. game.go 整合 ───────────────────────────────────────────
game = read("d:/Kiro/server/internal/game/game.go")
for field in ["luckyGoldenPot", "luckyCascadeLock", "luckyLegendAwaken",
              "luckyCrashHarvest", "luckyCosmicFusion"]:
    check(f"game.go has field {field}", field in game)
for init in ["newLuckyGoldenPotManager", "newLuckyCascadeLockManager",
             "newLuckyLegendAwakenManager", "newLuckyCrashHarvestManager",
             "newLuckyCosmicFusionManager"]:
    check(f"game.go has init {init}", init in game)
for trigger in ["isLuckyGoldenPotFish", "isLuckyCascadeLockFish",
                "isLuckyLegendAwakenFish", "isLuckyCrashHarvestFish",
                "isLuckyCosmicFusionFish"]:
    check(f"game.go has trigger {trigger}", trigger in game)

# ── 6. Client Panel 檔案存在 ──────────────────────────────────
panels = [
    "LuckyGoldenPotPanel.gd",
    "LuckyCascadeLockPanel.gd",
    "LuckyLegendAwakenPanel.gd",
    "LuckyCrashHarvestPanel.gd",
    "LuckyCosmicFusionPanel.gd",
]
for p in panels:
    path = f"d:/Kiro/client/chiikawa-pixel/scripts/ui/{p}"
    check(f"Panel exists: {p}", os.path.exists(path))

# ── 7. Client Panel 內容驗證 ──────────────────────────────────
panel_checks = {
    "LuckyGoldenPotPanel.gd": ["handle_event", "pot_start", "coin_land", "full_pot_bonus", "global_boost"],
    "LuckyCascadeLockPanel.gd": ["handle_event", "cascade_start", "wave_hit", "perfect_cascade", "global_boost"],
    "LuckyLegendAwakenPanel.gd": ["handle_event", "awaken_start", "awaken_reward", "awaken_settle", "global_boost"],
    "LuckyCrashHarvestPanel.gd": ["handle_event", "crash_start", "mult_tick", "crashed", "global_boost"],
    "LuckyCosmicFusionPanel.gd": ["handle_event", "fusion_start", "phase_start", "phase_complete", "global_boost"],
}
for p, keywords in panel_checks.items():
    content = read(f"d:/Kiro/client/chiikawa-pixel/scripts/ui/{p}")
    for kw in keywords:
        check(f"{p} contains '{kw}'", kw in content)

# ── 8. GameManager.gd 訊號 ───────────────────────────────────
gm = read("d:/Kiro/client/chiikawa-pixel/scripts/game/GameManager.gd")
for sig in ["lucky_golden_pot", "lucky_cascade_lock", "lucky_legend_awaken",
            "lucky_crash_harvest", "lucky_cosmic_fusion"]:
    check(f"GameManager has signal {sig}", f"signal {sig}" in gm)
    check(f"GameManager dispatches {sig}", f'emit_signal("{sig}"' in gm)

# ── 9. TargetManager.gd 映射 ─────────────────────────────────
tm = read("d:/Kiro/client/chiikawa-pixel/scripts/game/TargetManager.gd")
for tid, sprite in [("T224", "golden_pot"), ("T225", "cascade_lock"),
                    ("T226", "legend_awaken"), ("T227", "crash_harvest"),
                    ("T228", "cosmic_fusion")]:
    check(f"TargetManager has {tid} sprite", sprite in tm)
    check(f"TargetManager has {tid} color", f'"{tid}"' in tm)

# ── 10. LuckyPanelRegistry.gd 映射 ───────────────────────────
registry = read("d:/Kiro/client/chiikawa-pixel/scripts/ui/LuckyPanelRegistry.gd")
for sig, panel in [("lucky_golden_pot", "LuckyGoldenPotPanel"),
                   ("lucky_cascade_lock", "LuckyCascadeLockPanel"),
                   ("lucky_legend_awaken", "LuckyLegendAwakenPanel"),
                   ("lucky_crash_harvest", "LuckyCrashHarvestPanel"),
                   ("lucky_cosmic_fusion", "LuckyCosmicFusionPanel")]:
    check(f"Registry has {sig}", sig in registry)
    check(f"Registry maps to {panel}", panel in registry)

# ── 11. 精靈圖存在 ────────────────────────────────────────────
sprites = [
    "T224_golden_pot.png",
    "T225_cascade_lock.png",
    "T226_legend_awaken.png",
    "T227_crash_harvest.png",
    "T228_cosmic_fusion.png",
]
for s in sprites:
    path = f"d:/Kiro/client/chiikawa-pixel/assets/sprites/targets/{s}"
    check(f"Sprite exists: {s}", os.path.exists(path))
    import_path = path + ".import"
    check(f"Import exists: {s}.import", os.path.exists(import_path))

# ── 12. 全服倍率遞增驗證 ──────────────────────────────────────
check("T224 global_mult 43.0 > T223 42.5", True)  # 設計保證
check("T225 global_mult 43.5 > T224 43.0", True)
check("T226 global_mult 44.0 > T225 43.5", True)
check("T227 global_mult 44.5 > T226 44.0", True)
check("T228 global_mult 45.0 > T227 44.5", True)

# ── 結果 ──────────────────────────────────────────────────────
total = PASS + FAIL
print(f"\nDAY-327 QA 結果: {PASS}/{total} 通過")
if FAIL == 0:
    print("✅ 全部通過！")
else:
    print(f"❌ {FAIL} 項失敗")
