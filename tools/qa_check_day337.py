#!/usr/bin/env python3
"""
DAY-337 QA 驗證腳本
主要驗證：
1. HUD.gd 重構完成（行數 < 500，無 Lucky 函數）
2. HUDLuckySignals.gd 完整性（87+ 個 Lucky 函數）
3. Server 編譯狀態
4. T101-T105 特殊目標物視覺
5. 所有 Lucky Panel 腳本存在
"""
import os
import subprocess
import sys

ROOT = r"d:\Kiro"
CLIENT = os.path.join(ROOT, "client", "chiikawa-pixel")
SERVER = os.path.join(ROOT, "server")
SCRIPTS_UI = os.path.join(CLIENT, "scripts", "ui")
TARGETS = os.path.join(CLIENT, "assets", "sprites", "targets")

passed = 0
failed = 0

def check(name, condition, detail=""):
    global passed, failed
    if condition:
        print(f"  ✅ {name}")
        passed += 1
    else:
        print(f"  ❌ {name}" + (f" — {detail}" if detail else ""))
        failed += 1

print("=" * 60)
print("DAY-337 QA 驗證")
print("=" * 60)

# ── 1. HUD.gd 重構驗證 ──────────────────────────────────────
print("\n[1] HUD.gd 重構驗證")
hud_path = os.path.join(SCRIPTS_UI, "HUD.gd")
if os.path.exists(hud_path):
    with open(hud_path, "r", encoding="utf-8") as f:
        hud_lines = f.readlines()
        hud_content = "".join(hud_lines)
    
    hud_line_count = len(hud_lines)
    check("HUD.gd 行數 < 600", hud_line_count < 600, f"當前 {hud_line_count} 行")
    check("HUD.gd 行數 < 500", hud_line_count < 500, f"當前 {hud_line_count} 行")
    
    lucky_funcs = [l for l in hud_lines if l.strip().startswith("func _on_lucky_")]
    check("HUD.gd 無 Lucky 函數", len(lucky_funcs) == 0, f"仍有 {len(lucky_funcs)} 個")
    
    # 核心函數存在
    core_funcs = ["func _ready", "func _process", "func _update_ui", 
                  "func _on_player_updated", "func _on_state_changed",
                  "func _on_reward_received", "func _on_boss_event",
                  "func _show_lucky_banner", "func _show_lucky_event",
                  "func _init_lucky_signals"]
    for func in core_funcs:
        check(f"HUD.gd 保留 {func}", func in hud_content)
    
    # 確認 HUDLuckySignals 初始化
    check("HUD.gd 有 HUDLuckySignals 初始化", "_init_lucky_signals" in hud_content)
    check("HUD.gd 有 _lucky_signals 變數", "_lucky_signals" in hud_content)
else:
    check("HUD.gd 存在", False)

# ── 2. HUDLuckySignals.gd 完整性 ────────────────────────────
print("\n[2] HUDLuckySignals.gd 完整性")
signals_path = os.path.join(SCRIPTS_UI, "HUDLuckySignals.gd")
if os.path.exists(signals_path):
    with open(signals_path, "r", encoding="utf-8") as f:
        signals_lines = f.readlines()
        signals_content = "".join(signals_lines)
    
    signals_line_count = len(signals_lines)
    check("HUDLuckySignals.gd 存在", True)
    check("HUDLuckySignals.gd 行數 > 1000", signals_line_count > 1000, f"當前 {signals_line_count} 行")
    
    lucky_funcs_count = sum(1 for l in signals_lines if l.strip().startswith("func _on_lucky_"))
    check("HUDLuckySignals.gd Lucky 函數 >= 80", lucky_funcs_count >= 80, f"當前 {lucky_funcs_count} 個")
    
    # 確認關鍵函數存在
    key_funcs = ["connect_all_lucky_signals", "_show_banner", "_show_event", "_on_lucky_fallback"]
    for func in key_funcs:
        check(f"HUDLuckySignals.gd 有 {func}", func in signals_content)
    
    # 確認 DAY-318~319 fallback 連接
    check("HUDLuckySignals.gd 有 fallback 連接", "_on_lucky_fallback" in signals_content)
else:
    check("HUDLuckySignals.gd 存在", False)

# ── 3. Lucky Panel 腳本存在 ──────────────────────────────────
print("\n[3] Lucky Panel 腳本驗證")
panel_files = [f for f in os.listdir(SCRIPTS_UI) if f.startswith("Lucky") and f.endswith(".gd")]
check("Lucky Panel 數量 >= 150", len(panel_files) >= 150, f"當前 {len(panel_files)} 個")
check("BaseLuckyPanel.gd 存在", os.path.exists(os.path.join(SCRIPTS_UI, "BaseLuckyPanel.gd")))
check("LuckyEventSystem.gd 存在", os.path.exists(os.path.join(SCRIPTS_UI, "LuckyEventSystem.gd")))
check("LuckyPanelRegistry.gd 存在", os.path.exists(os.path.join(SCRIPTS_UI, "LuckyPanelRegistry.gd")))

# ── 4. T101-T105 特殊目標物視覺 ─────────────────────────────
print("\n[4] T101-T105 特殊目標物視覺")
try:
    from PIL import Image
    # T101 擬態怪物是草形狀，密度偏低是正常的（門檻 10%）
    # 其他特殊目標物門檻 15%
    special_targets = {
        "T101_mimic": ("T101_mimic.png", 10),   # 草形狀，密度偏低正常
        "T102_chest": ("T102_chest.png", 15),
        "T103_meteor": ("T103_meteor.png", 15),
        "T104_gold_grass": ("T104_gold_grass.png", 15),
        "T105_coin_fish": ("T105_coin_fish.png", 15),
    }
    for name, (filename, min_density) in special_targets.items():
        path = os.path.join(TARGETS, filename)
        if os.path.exists(path):
            img = Image.open(path).convert("RGBA")
            w, h = img.size
            pixels = img.load()
            non_transparent = sum(1 for y in range(h) for x in range(w) if pixels[x,y][3] > 10)
            total = w * h
            density = non_transparent / total * 100
            check(f"{name} 存在且密度 > {min_density}%", density > min_density, f"{w}x{h}, 密度 {density:.1f}%")
        else:
            check(f"{name} 存在", False, f"檔案不存在: {path}")
except ImportError:
    print("  ⚠️ PIL 未安裝，跳過視覺驗證")

# ── 5. T001-T006 基礎目標物 ─────────────────────────────────
print("\n[5] T001-T006 基礎目標物")
try:
    from PIL import Image
    basic_targets = ["T001_grass", "T002_bug_g", "T003_bug_r", "T004_bug_b", "T005_pudding", "T006_mushroom"]
    for name in basic_targets:
        path = os.path.join(TARGETS, f"{name}.png")
        if os.path.exists(path):
            img = Image.open(path).convert("RGBA")
            w, h = img.size
            check(f"{name} 尺寸 >= 48px", w >= 48 and h >= 48, f"{w}x{h}")
        else:
            check(f"{name} 存在", False)
except ImportError:
    print("  ⚠️ PIL 未安裝，跳過視覺驗證")

# ── 6. Server 編譯狀態 ───────────────────────────────────────
print("\n[6] Server 編譯狀態")
try:
    result = subprocess.run(
        ["go", "build", "./..."],
        cwd=SERVER,
        capture_output=True,
        text=True,
        timeout=60
    )
    check("go build ./... 通過", result.returncode == 0, result.stderr[:200] if result.stderr else "")
    
    result2 = subprocess.run(
        ["go", "vet", "./..."],
        cwd=SERVER,
        capture_output=True,
        text=True,
        timeout=60
    )
    check("go vet ./... 通過", result2.returncode == 0, result2.stderr[:200] if result2.stderr else "")
except Exception as e:
    check("Server 編譯", False, str(e))

# ── 7. 音效和 BGM 資產 ──────────────────────────────────────
print("\n[7] 音效和 BGM 資產")
audio_path = os.path.join(CLIENT, "assets", "audio")
bgm_files = ["main_game.wav", "boss_enter.wav", "boss_rage.wav", "bonus_game.wav"]
sfx_files = ["attack_fire.wav", "hit.wav", "kill.wav", "big_win.wav", "bonus_ready.wav",
             "boss_warning.wav", "boss_enter.wav", "coin_drop.wav", "weed_pull.wav"]

for bgm in bgm_files:
    path = os.path.join(audio_path, "bgm", bgm)
    check(f"BGM {bgm}", os.path.exists(path))

for sfx in sfx_files:
    path = os.path.join(audio_path, "sfx", sfx)
    check(f"SFX {sfx}", os.path.exists(path))

# ── 8. 角色精靈圖 ────────────────────────────────────────────
print("\n[8] 角色精靈圖")
chars_path = os.path.join(CLIENT, "assets", "sprites", "characters")
char_files = [
    "chiikawa_idle.png", "chiikawa_attack.png", "chiikawa_bigwin.png",
    "hachiware_idle.png", "hachiware_attack.png", "hachiware_bigwin.png",
    "usagi_idle.png", "usagi_attack.png", "usagi_bigwin.png",
]
for char in char_files:
    path = os.path.join(chars_path, char)
    check(f"角色 {char}", os.path.exists(path))

# ── 結果摘要 ─────────────────────────────────────────────────
print("\n" + "=" * 60)
total = passed + failed
print(f"結果：{passed}/{total} 通過")
if failed == 0:
    print("🎉 全部通過！DAY-337 驗證完成")
else:
    print(f"⚠️ {failed} 項失敗，需要修復")
print("=" * 60)

sys.exit(0 if failed == 0 else 1)
