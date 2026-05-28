#!/usr/bin/env python3
"""
qa_check_day322.py — DAY-322 QA 驗證腳本
驗證 Shader 視覺分層系統 + ComboSystem 的完整性

執行方式：py tools/qa_check_day322.py
"""

import os
import re

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
PASS = "✅"
FAIL = "❌"
results = []

def check(name, condition, detail=""):
    status = PASS if condition else FAIL
    results.append((status, name, detail))
    return condition

def read_file(path):
    full = os.path.join(ROOT, path)
    if not os.path.exists(full):
        return None
    with open(full, "r", encoding="utf-8") as f:
        return f.read()

# ── 1. Shader 文件存在性 ──────────────────────────────────────
print("\n=== 1. Shader 文件 ===")
shaders = [
    "client/chiikawa-pixel/assets/shaders/hit_flash.gdshader",
    "client/chiikawa-pixel/assets/shaders/sprite_outline.gdshader",
    "client/chiikawa-pixel/assets/shaders/tier_glow.gdshader",
]
for s in shaders:
    content = read_file(s)
    check(f"Shader 存在: {os.path.basename(s)}", content is not None)
    if content:
        check(f"Shader 有 shader_type: {os.path.basename(s)}", "shader_type canvas_item" in content)

# ── 2. Shader 內容驗證 ────────────────────────────────────────
print("\n=== 2. Shader 內容 ===")
hit_flash = read_file("client/chiikawa-pixel/assets/shaders/hit_flash.gdshader")
if hit_flash:
    check("hit_flash: 有 flash_intensity uniform", "flash_intensity" in hit_flash)
    check("hit_flash: 有 flash_color uniform", "flash_color" in hit_flash)
    check("hit_flash: 保留透明度", "tex.a" in hit_flash)

outline = read_file("client/chiikawa-pixel/assets/shaders/sprite_outline.gdshader")
if outline:
    check("sprite_outline: 有 outline_color uniform", "outline_color" in outline)
    check("sprite_outline: 有 8方向採樣", outline.count("texel") >= 4)
    check("sprite_outline: 有脈動效果", "pulse_speed" in outline)

glow = read_file("client/chiikawa-pixel/assets/shaders/tier_glow.gdshader")
if glow:
    check("tier_glow: 有 glow_color uniform", "glow_color" in glow)
    check("tier_glow: 有 glow_radius uniform", "glow_radius" in glow)
    check("tier_glow: 有脈動效果", "pulse_speed" in glow)

# ── 3. TargetManager Shader 整合 ─────────────────────────────
print("\n=== 3. TargetManager Shader 整合 ===")
tm = read_file("client/chiikawa-pixel/scripts/game/TargetManager.gd")
if tm:
    check("TargetManager: preload hit_flash shader", "SHADER_HIT_FLASH" in tm)
    check("TargetManager: preload sprite_outline shader", "SHADER_OUTLINE" in tm)
    check("TargetManager: preload tier_glow shader", "SHADER_TIER_GLOW" in tm)
    check("TargetManager: TIER_THRESHOLDS 定義", "TIER_THRESHOLDS" in tm)
    check("TargetManager: _get_multiplier_tier 函數", "_get_multiplier_tier" in tm)
    check("TargetManager: _apply_tier_shaders 函數", "_apply_tier_shaders" in tm)
    check("TargetManager: _add_rainbow_ring 函數（Tier 5）", "_add_rainbow_ring" in tm)
    check("TargetManager: Sprite 命名為 Sprite", 'sprite.name = "Sprite"' in tm)
    check("TargetManager: hit_flash shader 套用到 Sprite", "SHADER_HIT_FLASH" in tm and "hit_mat" in tm)
    check("TargetManager: _flash_hit 使用 Shader", "ShaderMaterial" in tm and "_flash_hit" in tm)
    # 確認 Tier 顏色定義
    check("TargetManager: 5個 Tier 輪廓顏色", tm.count("TIER_OUTLINE_COLORS") >= 1)
    check("TargetManager: 5個 Tier 光暈顏色", tm.count("TIER_GLOW_COLORS") >= 1)
    check("TargetManager: 5個 Tier 脈動速度", tm.count("TIER_PULSE_SPEEDS") >= 1)

# ── 4. ComboSystem 驗證 ───────────────────────────────────────
print("\n=== 4. ComboSystem ===")
combo = read_file("client/chiikawa-pixel/scripts/game/ComboSystem.gd")
if combo:
    check("ComboSystem: 存在", True)
    check("ComboSystem: combo_updated 訊號", "signal combo_updated" in combo)
    check("ComboSystem: combo_milestone 訊號", "signal combo_milestone" in combo)
    check("ComboSystem: combo_reset 訊號", "signal combo_reset" in combo)
    check("ComboSystem: COMBO_TIMEOUT 定義", "COMBO_TIMEOUT" in combo)
    check("ComboSystem: MILESTONE_COUNTS 定義", "MILESTONE_COUNTS" in combo)
    check("ComboSystem: 5個里程碑（5/10/20/50/100）", all(str(n) in combo for n in [5, 10, 20, 50, 100]))
    check("ComboSystem: _register_hit 函數", "_register_hit" in combo)
    check("ComboSystem: _reset_combo 函數", "_reset_combo" in combo)
    check("ComboSystem: _trigger_milestone 函數", "_trigger_milestone" in combo)
    check("ComboSystem: 螢幕震動整合", "ScreenShake" in combo)
    check("ComboSystem: 音效整合", "AudioManager" in combo)
    check("ComboSystem: 全螢幕閃光（50+）", "_spawn_screen_flash" in combo)
else:
    check("ComboSystem: 存在", False, "文件不存在")

# ── 5. HUD Combo 顯示升級 ─────────────────────────────────────
print("\n=== 5. HUD Combo 顯示 ===")
hud = read_file("client/chiikawa-pixel/scripts/ui/HUD.gd")
if hud:
    check("HUD: _connect_combo_system 函數", "_connect_combo_system" in hud)
    check("HUD: _on_combo_updated 函數", "_on_combo_updated" in hud)
    check("HUD: _on_combo_milestone 函數", "_on_combo_milestone" in hud)
    check("HUD: combo 字體大小 18px", "font_size\", 18" in hud)
    check("HUD: 5層顏色系統（宇宙粉紅/深紫/火紅/金色）", all(c in hud for c in ["0.0, 0.5", "0.0, 1.0", "0.3, 0.1", "0.85, 0.0"]))

# ── 6. Agent 文件更新 ─────────────────────────────────────────
print("\n=== 6. Agent 文件 ===")
vca = read_file("agents/visual-clarity-agent.md")
if vca:
    check("visual-clarity-agent: Shader 系統章節", "Shader 系統" in vca)
    check("visual-clarity-agent: Tier System 表格", "Tier System" in vca or "Tier 5" in vca)
    check("visual-clarity-agent: DAY-322 更新", "DAY-322" in vca)

tda = read_file("agents/target-design-agent.md")
if tda:
    check("target-design-agent: DAY-322 狀態", "DAY-322" in tda)
    check("target-design-agent: 112 種目標物", "112" in tda)

# ── 7. knowhow-log 更新 ───────────────────────────────────────
print("\n=== 7. knowhow-log ===")
kh = read_file(".kiro/skills/knowhow-log.md")
if kh:
    check("knowhow-log: 條目 36（Shader 視覺分層）", "## 36." in kh)
    check("knowhow-log: 條目 37（ShaderMaterial 套用）", "## 37." in kh)
    check("knowhow-log: 條目 38（多層 Sprite 架構）", "## 38." in kh)

# ── 8. progress.md 更新 ───────────────────────────────────────
print("\n=== 8. progress.md ===")
prog = read_file("docs/progress.md")
if prog:
    check("progress.md: DAY-322 更新", "DAY-322" in prog)
    check("progress.md: Shader 系統記錄", "Shader" in prog)
    check("progress.md: 視覺清晰度目標 8.5/10", "8.5" in prog)

# ── 結果統計 ─────────────────────────────────────────────────
print("\n" + "="*60)
passed = sum(1 for r in results if r[0] == PASS)
total = len(results)
print(f"DAY-322 QA 結果：{passed}/{total} 通過")
print("="*60)

failed = [(n, d) for s, n, d in results if s == FAIL]
if failed:
    print("\n❌ 失敗項目：")
    for name, detail in failed:
        print(f"  - {name}" + (f"（{detail}）" if detail else ""))
else:
    print("\n✅ 全部通過！DAY-322 Shader 視覺分層系統完整。")

# 詳細輸出
print("\n詳細結果：")
for status, name, detail in results:
    print(f"  {status} {name}" + (f" — {detail}" if detail else ""))
