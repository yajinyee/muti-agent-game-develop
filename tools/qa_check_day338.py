#!/usr/bin/env python3
"""
qa_check_day338.py — DAY-338 品質驗證腳本
qa-playtest-agent 負責維護

DAY-338 改善項目：
1. 打擊感優化（Cannon.gd 本地預測命中音效）
2. 整合測試 recv_until 修復（能正確過濾非目標訊息）
3. Server 端對端整合測試通過率驗證

驗證項目：
- Cannon.gd 打擊感改善（本地預測 + Server 確認分離）
- integration_test_day334.py recv_until 修復
- Server 編譯狀態
- 核心 GDScript 檔案完整性
"""

import os
import sys
import subprocess
import re

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
CLIENT_SCRIPTS = os.path.join(ROOT, "client", "chiikawa-pixel", "scripts")
SERVER_DIR = os.path.join(ROOT, "server")
TOOLS_DIR = os.path.join(ROOT, "tools")

passed = 0
failed = 0

def ok(msg):
    global passed
    passed += 1
    print(f"  ✅ {msg}")

def fail(msg):
    global failed
    failed += 1
    print(f"  ❌ {msg}")

def check(condition, msg_ok, msg_fail):
    if condition:
        ok(msg_ok)
    else:
        fail(msg_fail)

def read_file(path):
    try:
        with open(path, "r", encoding="utf-8") as f:
            return f.read()
    except:
        return ""

print("\n" + "="*60)
print("DAY-338 品質驗證")
print("="*60)

# ── 1. Server 編譯狀態 ────────────────────────────────────────
print("\n[1] Server 編譯狀態")
try:
    result = subprocess.run(
        ["go", "build", "./..."],
        cwd=SERVER_DIR,
        capture_output=True, text=True, timeout=60
    )
    check(result.returncode == 0, "go build ./... 通過", f"go build 失敗: {result.stderr[:200]}")
except Exception as e:
    fail(f"go build 執行失敗: {e}")

try:
    result = subprocess.run(
        ["go", "vet", "./..."],
        cwd=SERVER_DIR,
        capture_output=True, text=True, timeout=60
    )
    check(result.returncode == 0, "go vet ./... 通過", f"go vet 失敗: {result.stderr[:200]}")
except Exception as e:
    fail(f"go vet 執行失敗: {e}")

# ── 2. Cannon.gd 打擊感優化驗證 ──────────────────────────────
print("\n[2] Cannon.gd 打擊感優化（DAY-338）")
cannon_path = os.path.join(CLIENT_SCRIPTS, "game", "Cannon.gd")
cannon = read_file(cannon_path)

check(cannon != "", "Cannon.gd 存在", "Cannon.gd 不存在")
check("_spawn_projectile_with_impact" in cannon,
      "新增 _spawn_projectile_with_impact 函數（本地預測）",
      "缺少 _spawn_projectile_with_impact 函數")
check("DAY-338" in cannon,
      "Cannon.gd 包含 DAY-338 標記",
      "Cannon.gd 缺少 DAY-338 標記")
check("target_id != \"\"" in cannon,
      "投射物到達時依 target_id 決定是否播放命中音效",
      "缺少 target_id 判斷邏輯")
check("is_kill" in cannon and "ScreenShake.add_trauma(0.35)" in cannon,
      "擊破時加強震動（0.35）",
      "缺少擊破加強震動邏輯")
# 確認不再在 _on_attack_result 中播放 HIT 音效（已移到本地預測）
check("AudioManager.play_sfx(AudioManager.SFX.HIT)" not in cannon.split("_on_attack_result")[1].split("func ")[0] if "_on_attack_result" in cannon else True,
      "_on_attack_result 不再重複播放 HIT 音效",
      "_on_attack_result 仍在重複播放 HIT 音效")

# ── 3. 整合測試修復驗證 ───────────────────────────────────────
print("\n[3] 整合測試 recv_until 修復")
test_path = os.path.join(TOOLS_DIR, "integration_test_day334.py")
test_content = read_file(test_path)

check(test_content != "", "integration_test_day334.py 存在", "integration_test_day334.py 不存在")
check("# 跳過其他類型的訊息，繼續等待" in test_content,
      "recv_until 已修復（跳過非目標訊息）",
      "recv_until 未修復")

# ── 4. HUD.gd 重構完整性驗證 ─────────────────────────────────
print("\n[4] HUD.gd 重構完整性（DAY-337）")
hud_path = os.path.join(CLIENT_SCRIPTS, "ui", "HUD.gd")
hud = read_file(hud_path)

check(hud != "", "HUD.gd 存在", "HUD.gd 不存在")
hud_lines = len(hud.splitlines())
check(hud_lines < 500, f"HUD.gd 行數 {hud_lines} < 500（重構後）", f"HUD.gd 行數 {hud_lines} 過多")
check("HUDLuckySignals" in hud, "HUD.gd 引用 HUDLuckySignals", "HUD.gd 缺少 HUDLuckySignals 引用")
check("_init_lucky_signals" in hud, "HUD.gd 包含 _init_lucky_signals", "HUD.gd 缺少 _init_lucky_signals")

# ── 5. HUDLuckySignals.gd 完整性 ─────────────────────────────
print("\n[5] HUDLuckySignals.gd 完整性")
lucky_signals_path = os.path.join(CLIENT_SCRIPTS, "ui", "HUDLuckySignals.gd")
lucky_signals = read_file(lucky_signals_path)

check(lucky_signals != "", "HUDLuckySignals.gd 存在", "HUDLuckySignals.gd 不存在")
lucky_lines = len(lucky_signals.splitlines())
check(lucky_lines > 1000, f"HUDLuckySignals.gd 行數 {lucky_lines} > 1000（包含所有訊號）",
      f"HUDLuckySignals.gd 行數 {lucky_lines} 不足")
check("connect_all_lucky_signals" in lucky_signals,
      "HUDLuckySignals.gd 包含 connect_all_lucky_signals",
      "HUDLuckySignals.gd 缺少 connect_all_lucky_signals")

# ── 6. GameManager.gd 訊號完整性 ─────────────────────────────
print("\n[6] GameManager.gd 訊號完整性")
gm_path = os.path.join(CLIENT_SCRIPTS, "game", "GameManager.gd")
gm = read_file(gm_path)

check(gm != "", "GameManager.gd 存在", "GameManager.gd 不存在")
# 確認核心訊號存在
core_signals = ["player_updated", "game_state_changed", "reward_received",
                "attack_result", "target_spawned", "target_updated", "target_killed",
                "boss_event", "bonus_event"]
for sig in core_signals:
    check(f"signal {sig}" in gm, f"GameManager 包含 signal {sig}", f"GameManager 缺少 signal {sig}")

# ── 7. TargetManager.gd 核心功能 ─────────────────────────────
print("\n[7] TargetManager.gd 核心功能")
tm_path = os.path.join(CLIENT_SCRIPTS, "game", "TargetManager.gd")
tm = read_file(tm_path)

check(tm != "", "TargetManager.gd 存在", "TargetManager.gd 不存在")
check("try_click_target" in tm, "TargetManager 包含 try_click_target", "TargetManager 缺少 try_click_target")
check("_flash_hit" in tm, "TargetManager 包含 _flash_hit（受擊閃白）", "TargetManager 缺少 _flash_hit")
check("HitEffect.spawn_kill" in tm, "TargetManager 呼叫 HitEffect.spawn_kill", "TargetManager 缺少 HitEffect.spawn_kill")
check("SHADER_HIT_FLASH" in tm, "TargetManager 使用 Hit Flash Shader", "TargetManager 缺少 Hit Flash Shader")

# ── 8. 核心 Shader 檔案存在 ──────────────────────────────────
print("\n[8] 核心 Shader 檔案")
shader_dir = os.path.join(ROOT, "client", "chiikawa-pixel", "assets", "shaders")
shaders = ["hit_flash.gdshader", "sprite_outline.gdshader", "tier_glow.gdshader"]
for shader in shaders:
    path = os.path.join(shader_dir, shader)
    check(os.path.exists(path), f"Shader 存在: {shader}", f"Shader 缺失: {shader}")

# ── 9. 音效檔案存在 ───────────────────────────────────────────
print("\n[9] 核心音效檔案")
sfx_dir = os.path.join(ROOT, "client", "chiikawa-pixel", "assets", "audio", "sfx")
sfx_files = ["attack_fire.wav", "attack_fire_hachiware.wav", "attack_fire_usagi.wav",
             "hit.wav", "kill.wav", "big_win.wav", "coin_drop.wav",
             "boss_warning.wav", "boss_enter.wav", "bonus_ready.wav",
             "bonus_game.wav", "weed_pull.wav"]
for sfx in sfx_files:
    path = os.path.join(sfx_dir, sfx)
    check(os.path.exists(path), f"SFX 存在: {sfx}", f"SFX 缺失: {sfx}")

# ── 10. BGM 檔案存在 ──────────────────────────────────────────
print("\n[10] BGM 檔案")
bgm_dir = os.path.join(ROOT, "client", "chiikawa-pixel", "assets", "audio", "bgm")
bgm_files = ["main_game.wav", "boss_enter.wav", "boss_rage.wav", "bonus_game.wav"]
for bgm in bgm_files:
    path = os.path.join(bgm_dir, bgm)
    check(os.path.exists(path), f"BGM 存在: {bgm}", f"BGM 缺失: {bgm}")

# ── 11. 基礎目標物 Sprite 存在 ────────────────────────────────
print("\n[11] 基礎目標物 Sprite（T001-T006 + B001）")
sprite_dir = os.path.join(ROOT, "client", "chiikawa-pixel", "assets", "sprites", "targets")
basic_sprites = ["T001_grass.png", "T002_bug_g.png", "T003_bug_r.png",
                 "T004_bug_b.png", "T005_pudding.png", "T006_mushroom.png", "B001_boss.png"]
for sprite in basic_sprites:
    path = os.path.join(sprite_dir, sprite)
    check(os.path.exists(path), f"Sprite 存在: {sprite}", f"Sprite 缺失: {sprite}")

# ── 12. 特殊目標物 Sprite 存在（T101-T105）────────────────────
print("\n[12] 特殊目標物 Sprite（T101-T105，DAY-337 升級版）")
special_sprites = ["T101_mimic.png", "T102_chest.png", "T103_meteor.png",
                   "T104_gold_grass.png", "T105_coin_fish.png"]
for sprite in special_sprites:
    path = os.path.join(sprite_dir, sprite)
    check(os.path.exists(path), f"Sprite 存在: {sprite}", f"Sprite 缺失: {sprite}")

# ── 13. NetworkManager.gd 完整性 ─────────────────────────────
print("\n[13] NetworkManager.gd 完整性")
nm_path = os.path.join(CLIENT_SCRIPTS, "network", "NetworkManager.gd")
nm = read_file(nm_path)

check(nm != "", "NetworkManager.gd 存在", "NetworkManager.gd 不存在")
check("SERVER_URL_LOCAL = \"ws://localhost:7777/ws\"" in nm,
      "NetworkManager 連線到 Port 7777",
      "NetworkManager Port 設定錯誤")
check("send_attack" in nm, "NetworkManager 包含 send_attack", "NetworkManager 缺少 send_attack")
check("send_auto_toggle" in nm, "NetworkManager 包含 send_auto_toggle", "NetworkManager 缺少 send_auto_toggle")
check("send_bet_change" in nm, "NetworkManager 包含 send_bet_change", "NetworkManager 缺少 send_bet_change")

# ── 14. Main.tscn 場景結構 ────────────────────────────────────
print("\n[14] Main.tscn 場景結構")
tscn_path = os.path.join(ROOT, "client", "chiikawa-pixel", "scenes", "Main.tscn")
tscn = read_file(tscn_path)

check(tscn != "", "Main.tscn 存在", "Main.tscn 不存在")
check("TargetManager.gd" in tscn, "Main.tscn 包含 TargetManager", "Main.tscn 缺少 TargetManager")
check("Cannon.gd" in tscn, "Main.tscn 包含 Cannon", "Main.tscn 缺少 Cannon")
check("HUD.gd" in tscn, "Main.tscn 包含 HUD", "Main.tscn 缺少 HUD")
check("LuckyEventSystem.gd" in tscn, "Main.tscn 包含 LuckyEventSystem", "Main.tscn 缺少 LuckyEventSystem")

# ── 結果摘要 ──────────────────────────────────────────────────
total = passed + failed
print("\n" + "="*60)
print(f"DAY-338 QA 結果")
print("="*60)
print(f"  ✅ 通過: {passed}")
print(f"  ❌ 失敗: {failed}")
print(f"  總計: {total}")
rate = passed / total * 100 if total > 0 else 0
print(f"\n  通過率: {rate:.1f}%")

if failed == 0:
    print("\n  🎉 所有驗證通過！DAY-338 品質確認完成")
else:
    print(f"\n  ⚠️  有 {failed} 個驗證失敗，請檢查")

print("="*60 + "\n")
sys.exit(0 if failed == 0 else 1)
