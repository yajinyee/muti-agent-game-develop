"""
DAY-347 QA 驗證腳本
驗證：T221-T253 美術升級 + 賽季通行證系統 + Server 編譯
"""
import os
import subprocess
import sys

SPRITES_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SERVER_DIR = r"d:\Kiro\server"
CLIENT_UI_DIR = r"d:\Kiro\client\chiikawa-pixel\scripts\ui"
CLIENT_GAME_DIR = r"d:\Kiro\client\chiikawa-pixel\scripts\game"
SERVER_GAME_DIR = r"d:\Kiro\server\internal\game"
SERVER_PROTO_DIR = r"d:\Kiro\server\internal\protocol"

passed = 0
failed = 0

def check(name, condition, detail=""):
    global passed, failed
    if condition:
        print(f"  ✅ {name}")
        passed += 1
    else:
        print(f"  ❌ {name}" + (f" ({detail})" if detail else ""))
        failed += 1

print("=== DAY-347 QA 驗證 ===\n")

# ── 1. T221-T253 美術升級驗證 ──────────────────────────────────
print("【1. T221-T253 美術升級】")
target_names = {
    221: "T221_dice_bonus", 222: "T222_dual_bonus", 223: "T223_coin_respin",
    224: "T224_golden_pot", 225: "T225_cascade_lock", 226: "T226_legend_awaken",
    227: "T227_crash_harvest", 228: "T228_cosmic_fusion", 229: "T229_magnetic_attraction",
    230: "T230_super_chain", 231: "T231_holy_pillar", 232: "T232_time_stop",
    233: "T233_cosmic_restart", 234: "T234_fever_boost_ultimate", 235: "T235_rapid_riches_ultimate",
    236: "T236_ice_fishing_master", 237: "T237_cosmic_miracle", 238: "T238_genesis_ultimate",
    239: "T239_shark_spark", 240: "T240_winter_ice", 241: "T241_atlantis_frenzy",
    242: "T242_fishing_time_wheel", 243: "T243_ultimate_shark", 244: "T244_wild_collector",
    245: "T245_lightning_eel_ultra", 246: "T246_domino_chain", 247: "T247_immortal_boss_ultra",
    248: "T248_quad_fusion", 249: "T249_electrical_frame", 250: "T250_magnetic_respin",
    251: "T251_fisherman_trail", 252: "T252_golden_gills", 253: "T253_penta_fusion",
}

for tid, name in target_names.items():
    png_path = os.path.join(SPRITES_DIR, f"{name}.png")
    import_path = os.path.join(SPRITES_DIR, f"{name}.png.import")
    check(f"T{tid} PNG 存在", os.path.exists(png_path))
    check(f"T{tid} .import 存在", os.path.exists(import_path))

# ── 2. 備份驗證 ──────────────────────────────────────────────
print("\n【2. 備份驗證】")
backup_dir = r"d:\Kiro\data\tmp\targets_backup_day347"
check("備份目錄存在", os.path.isdir(backup_dir))
if os.path.isdir(backup_dir):
    backup_count = len([f for f in os.listdir(backup_dir) if f.endswith(".png")])
    check(f"備份數量 = 33", backup_count == 33, f"實際: {backup_count}")

# ── 3. 賽季通行證 Server 端 ──────────────────────────────────
print("\n【3. 賽季通行證 Server 端】")
season_pass_file = os.path.join(SERVER_GAME_DIR, "season_pass.go")
check("season_pass.go 存在", os.path.exists(season_pass_file))

if os.path.exists(season_pass_file):
    with open(season_pass_file, "r", encoding="utf-8") as f:
        content = f.read()
    check("SeasonPassManager 定義", "SeasonPassManager" in content)
    check("SeasonPassTier 定義", "SeasonPassTier" in content)
    check("AddXP 函數", "func (m *SeasonPassManager) AddXP" in content)
    check("GetSnapshot 函數", "func (m *SeasonPassManager) GetSnapshot" in content)
    check("10個等級定義", "defaultSeasonTiers" in content)
    check("XP 常數定義", "XPPerKill" in content)

# ── 4. Protocol 訊息 ──────────────────────────────────────────
print("\n【4. Protocol 訊息】")
messages_file = os.path.join(SERVER_PROTO_DIR, "messages.go")
if os.path.exists(messages_file):
    with open(messages_file, "r", encoding="utf-8") as f:
        content = f.read()
    check("MsgSeasonPassUpdate 定義", "MsgSeasonPassUpdate" in content)
    check("MsgSeasonPassLevelUp 定義", "MsgSeasonPassLevelUp" in content)
    check("SeasonPassUpdatePayload 定義", "SeasonPassUpdatePayload" in content)
    check("SeasonPassLevelUpPayload 定義", "SeasonPassLevelUpPayload" in content)

# ── 5. game.go 整合 ──────────────────────────────────────────
print("\n【5. game.go 整合】")
game_file = os.path.join(SERVER_GAME_DIR, "game.go")
if os.path.exists(game_file):
    with open(game_file, "r", encoding="utf-8") as f:
        content = f.read()
    check("seasonPass 欄位定義", "seasonPass *SeasonPassManager" in content)
    check("NewSeasonPassManager 初始化", "NewSeasonPassManager()" in content)
    check("addSeasonXP 函數", "func (g *Game) addSeasonXP" in content)
    check("kill 觸發 XP", "addSeasonXP(playerID, xpGain, xpSource)" in content)
    check("bonus 觸發 XP", "addSeasonXP(playerID, XPPerBonus" in content)
    check("combo 觸發 XP", "addSeasonXP(playerID, XPPerCombo5" in content)

# ── 6. Client 端 ──────────────────────────────────────────────
print("\n【6. Client 端】")
season_panel = os.path.join(CLIENT_UI_DIR, "SeasonPassPanel.gd")
check("SeasonPassPanel.gd 存在", os.path.exists(season_panel))

if os.path.exists(season_panel):
    with open(season_panel, "r", encoding="utf-8") as f:
        content = f.read()
    check("_build_ui 函數", "_build_ui" in content)
    check("_build_tier_grid 函數", "_build_tier_grid" in content)
    check("_on_season_pass_updated 函數", "_on_season_pass_updated" in content)
    check("_on_season_pass_level_up 函數", "_on_season_pass_level_up" in content)
    check("10個等級徽章", "TIER_BADGES" in content)

gm_file = os.path.join(CLIENT_GAME_DIR, "GameManager.gd")
if os.path.exists(gm_file):
    with open(gm_file, "r", encoding="utf-8") as f:
        content = f.read()
    check("season_pass_updated 訊號", "signal season_pass_updated" in content)
    check("season_pass_level_up 訊號", "signal season_pass_level_up" in content)
    check("season_pass_update 訊息處理", '"season_pass_update"' in content)
    check("season_pass_level_up 訊息處理", '"season_pass_level_up"' in content)

hud_file = os.path.join(CLIENT_UI_DIR, "HUD.gd")
if os.path.exists(hud_file):
    with open(hud_file, "r", encoding="utf-8") as f:
        content = f.read()
    check("_season_pass_panel 變數", "_season_pass_panel" in content)
    check("_create_season_button 函數", "_create_season_button" in content)
    check("SeasonPassPanel 載入", "SeasonPassPanel.gd" in content)

# ── 7. Server 編譯 ──────────────────────────────────────────
print("\n【7. Server 編譯】")
try:
    result = subprocess.run(
        ["go", "build", "./..."],
        cwd=SERVER_DIR,
        capture_output=True,
        text=True,
        timeout=60
    )
    check("go build 通過", result.returncode == 0, result.stderr[:200] if result.stderr else "")
    
    result2 = subprocess.run(
        ["go", "vet", "./..."],
        cwd=SERVER_DIR,
        capture_output=True,
        text=True,
        timeout=30
    )
    check("go vet 通過", result2.returncode == 0, result2.stderr[:200] if result2.stderr else "")
except Exception as e:
    check("Server 編譯", False, str(e))

# ── 結果 ──────────────────────────────────────────────────────
print(f"\n=== 結果：{passed}/{passed+failed} 通過 ===")
if failed == 0:
    print("🎉 全部通過！")
else:
    print(f"⚠️ {failed} 項失敗")

sys.exit(0 if failed == 0 else 1)
