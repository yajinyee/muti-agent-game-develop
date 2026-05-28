#!/usr/bin/env python3
"""
qa_check_day321.py — DAY-321 品質驗證腳本
驗證 Main.tscn 完整性、Lucky Panel 存在性、Server 結構
"""
import os
import re

ROOT = r'd:\Kiro'
CLIENT = os.path.join(ROOT, 'client', 'chiikawa-pixel')
SERVER = os.path.join(ROOT, 'server', 'internal', 'game')
SCRIPTS_UI = os.path.join(CLIENT, 'scripts', 'ui')
MAIN_TSCN = os.path.join(CLIENT, 'scenes', 'Main.tscn')

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
print("DAY-321 QA 驗證")
print("=" * 60)

# ── 1. Main.tscn 存在 ────────────────────────────────────────
print("\n[1] Main.tscn 基本結構")
main_content = open(MAIN_TSCN, encoding='utf-8').read()
check("Main.tscn 存在", os.path.exists(MAIN_TSCN))
check("load_steps=109", "load_steps=109" in main_content)
check("LuckyEventSystem 節點存在", 'name="LuckyEventSystem"' in main_content)
check("LuckyPanelRegistry 節點存在", 'name="LuckyPanelRegistry"' in main_content)
check("BonusGame 節點存在", 'name="BonusGame"' in main_content)
check("BackgroundManager 節點存在", 'name="BackgroundManager"' in main_content)
check("TargetManager 節點存在", 'name="TargetManager"' in main_content)
check("Cannon 節點存在", 'name="Cannon"' in main_content)
check("HUD 節點存在", 'name="HUD"' in main_content)
check("Camera2D 節點存在", 'name="Camera2D"' in main_content)

# ── 2. Lucky Panel 節點 ──────────────────────────────────────
print("\n[2] Lucky Panel 節點（100 個）")
lucky_panels_in_tscn = re.findall(r'name="(Lucky\w+Panel)"', main_content)
check(f"Lucky Panel 節點數量 = 100", len(lucky_panels_in_tscn) == 100, f"實際: {len(lucky_panels_in_tscn)}")

# ── 3. Lucky Panel 腳本存在 ──────────────────────────────────
print("\n[3] Lucky Panel 腳本存在性")
panel_scripts = [f for f in os.listdir(SCRIPTS_UI) if f.startswith('Lucky') and f.endswith('.gd') 
                 and f not in ['LuckyEventSystem.gd', 'LuckyPanelRegistry.gd', 'BaseLuckyPanel.gd']]
check(f"Lucky Panel 腳本數量 = 100", len(panel_scripts) == 100, f"實際: {len(panel_scripts)}")

# 確認每個 tscn 中的 Panel 都有對應腳本
missing_scripts = []
for panel_name in lucky_panels_in_tscn:
    script_path = os.path.join(SCRIPTS_UI, panel_name + '.gd')
    if not os.path.exists(script_path):
        missing_scripts.append(panel_name)
check("所有 Panel 都有對應腳本", len(missing_scripts) == 0, f"缺少: {missing_scripts[:5]}")

# ── 4. Server 結構 ───────────────────────────────────────────
print("\n[4] Server 結構")
server_files = os.listdir(SERVER)
lucky_handlers = [f for f in server_files if f.startswith('lucky_') and f.endswith('_handler.go')]
check(f"Lucky Handler 數量 >= 90", len(lucky_handlers) >= 90, f"實際: {len(lucky_handlers)}")
check("game.go 存在", 'game.go' in server_files)
check("player.go 存在", 'player.go' in server_files)
check("spawn.go 存在", 'spawn.go' in server_files)
check("target.go 存在", 'target.go' in server_files)

# ── 5. Sprite 資產 ───────────────────────────────────────────
print("\n[5] Sprite 資產")
sprites_dir = os.path.join(CLIENT, 'assets', 'sprites', 'targets')
png_files = [f for f in os.listdir(sprites_dir) if f.endswith('.png')]
check(f"Target Sprite 數量 = 112", len(png_files) == 112, f"實際: {len(png_files)}")
check("B001_boss.png 存在", 'B001_boss.png' in png_files)
check("T001_grass.png 存在", 'T001_grass.png' in png_files)
check("T205_cosmic_singularity.png 存在", 'T205_cosmic_singularity.png' in png_files)

# ── 6. 核心腳本 ──────────────────────────────────────────────
print("\n[6] 核心腳本")
scripts_game = os.path.join(CLIENT, 'scripts', 'game')
core_scripts = ['GameManager.gd', 'NetworkManager.gd', 'AudioManager.gd', 
                'TargetManager.gd', 'Cannon.gd', 'HitEffect.gd', 
                'ScreenShake.gd', 'BackgroundManager.gd', 'BonusGame.gd', 'CharacterAnimator.gd']
scripts_network = os.path.join(CLIENT, 'scripts', 'network')
for script in core_scripts:
    path1 = os.path.join(scripts_game, script)
    path2 = os.path.join(scripts_network, script)
    exists = os.path.exists(path1) or os.path.exists(path2)
    check(f"{script} 存在", exists)

# ── 7. project.godot Autoload ────────────────────────────────
print("\n[7] project.godot Autoload")
project_godot = open(os.path.join(CLIENT, 'project.godot'), encoding='utf-8').read()
for autoload in ['GameManager', 'NetworkManager', 'AudioManager', 'ScreenShake', 'HitEffect']:
    check(f"Autoload: {autoload}", autoload in project_godot)

# ── 8. LuckyPanelRegistry 訊號映射完整性 ─────────────────────
print("\n[8] LuckyPanelRegistry 訊號映射")
registry_content = open(os.path.join(SCRIPTS_UI, 'LuckyPanelRegistry.gd'), encoding='utf-8').read()
signal_mappings = re.findall(r'"(lucky_\w+)":\s+"Lucky\w+Panel"', registry_content)
# 注意：Jackpot Panel 用特殊分發機制，不在 SIGNAL_TO_PANEL 中
check(f"訊號映射數量 >= 95", len(signal_mappings) >= 95, f"實際: {len(signal_mappings)}")

# ── 9. GameManager 訊號定義 ──────────────────────────────────
print("\n[9] GameManager 訊號定義")
gm_content = open(os.path.join(scripts_game, 'GameManager.gd'), encoding='utf-8').read()
gm_signals = re.findall(r'signal (lucky_\w+)', gm_content)
# GameManager 有 96 個 Lucky 訊號 + lucky_jackpot_pool（統一 Jackpot 訊號）
check(f"GameManager Lucky 訊號數量 >= 95", len(gm_signals) >= 95, f"實際: {len(gm_signals)}")

# ── 結果 ─────────────────────────────────────────────────────
print("\n" + "=" * 60)
total = passed + failed
print(f"結果：{passed}/{total} 通過")
if failed == 0:
    print("🎉 全部通過！")
else:
    print(f"⚠️  {failed} 項失敗，需要修復")
print("=" * 60)
