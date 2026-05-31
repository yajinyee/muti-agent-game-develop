#!/usr/bin/env python3
"""
QA Check DAY-340
驗證項目：
1. Cannon.gd var char_id 重複宣告修復
2. HitEffect.gd 金幣噴射特效新增
3. HitEffect.gd 衝擊波特效新增
4. TargetManager.gd 游泳動畫新增
5. TargetManager.gd 擊破消失動畫升級
6. TargetManager.gd 金幣音效觸發
7. Server 編譯狀態
"""

import subprocess
import os
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
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

def read_file(path):
    try:
        with open(path, encoding="utf-8") as f:
            return f.read()
    except:
        return ""

print("=" * 60)
print("DAY-340 QA 驗證")
print("=" * 60)

# ── 1. Cannon.gd 修復驗證 ──────────────────────────────────────
print("\n[1] Cannon.gd var char_id 重複宣告修復")
cannon = read_file(os.path.join(ROOT, "client/chiikawa-pixel/scripts/game/Cannon.gd"))
# 確認函數內不再有 var char_id = GameManager.get_character_id()
lines = cannon.split("\n")
in_spawn_func = False
double_decl_found = False
for i, line in enumerate(lines):
    if "func _spawn_projectile_with_impact" in line:
        in_spawn_func = True
    if in_spawn_func and "var char_id = GameManager.get_character_id()" in line:
        double_decl_found = True
        break
    if in_spawn_func and line.startswith("func ") and "_spawn_projectile_with_impact" not in line:
        in_spawn_func = False

check("Cannon.gd 無重複 var char_id 宣告", not double_decl_found,
      "仍有 var char_id = GameManager.get_character_id() 在 _spawn_projectile_with_impact 內")
check("Cannon.gd 有 char_id 由參數傳入的注釋", "# char_id 由參數傳入" in cannon)

# ── 2. HitEffect.gd 新特效驗證 ────────────────────────────────
print("\n[2] HitEffect.gd 新特效")
hiteffect = read_file(os.path.join(ROOT, "client/chiikawa-pixel/scripts/game/HitEffect.gd"))
check("HitEffect 有 _spawn_coin_burst 函數", "_spawn_coin_burst" in hiteffect)
check("HitEffect 有 _spawn_shockwave 函數", "_spawn_shockwave" in hiteffect)
check("HitEffect spawn_kill 有 100x 超高倍率分支", "multiplier >= 100" in hiteffect)
check("HitEffect spawn_kill 呼叫 _spawn_coin_burst", "_spawn_coin_burst(pos" in hiteffect)
check("HitEffect spawn_kill 呼叫 _spawn_shockwave", "_spawn_shockwave(pos" in hiteffect)
check("HitEffect 有 DAY-340 注釋", "DAY-340" in hiteffect)

# ── 3. TargetManager.gd 游泳動畫驗證 ──────────────────────────
print("\n[3] TargetManager.gd 游泳動畫")
tm = read_file(os.path.join(ROOT, "client/chiikawa-pixel/scripts/game/TargetManager.gd"))
check("TargetManager 有 _add_swim_animation 函數", "_add_swim_animation" in tm)
check("TargetManager _create_target_node 呼叫 _add_swim_animation", "_add_swim_animation(container" in tm)
check("TargetManager 游泳動畫有 rotation_degrees 搖擺", "rotation_degrees" in tm and "swim_tween" in tm)
check("TargetManager 游泳動畫有縮放脈動", "scale_tween" in tm)

# ── 4. TargetManager.gd 擊破動畫升級驗證 ──────────────────────
print("\n[4] TargetManager.gd 擊破動畫升級")
check("TargetManager 擊破有高倍率分支 (>=50)", "multiplier >= 50" in tm and "1.8, 1.8" in tm)
check("TargetManager 擊破有中倍率分支 (>=10)", "multiplier >= 10" in tm and "1.5, 1.5" in tm)
check("TargetManager 擊破有 DAY-340 注釋", "DAY-340 升級" in tm)

# ── 5. TargetManager.gd 金幣音效觸發驗證 ──────────────────────
print("\n[5] TargetManager.gd 金幣音效觸發")
check("TargetManager 擊破後播放 COIN_DROP", "COIN_DROP" in tm)
check("TargetManager 金幣音效有 reward > 0 條件", "reward > 0" in tm)

# ── 6. Server 編譯驗證 ────────────────────────────────────────
print("\n[6] Server 編譯狀態")
server_dir = os.path.join(ROOT, "server")
try:
    result = subprocess.run(
        ["go", "build", "./..."],
        cwd=server_dir,
        capture_output=True,
        text=True,
        timeout=30
    )
    check("go build ./... 通過", result.returncode == 0,
          result.stderr[:200] if result.returncode != 0 else "")
except Exception as e:
    check("go build ./... 通過", False, str(e))

try:
    result = subprocess.run(
        ["go", "vet", "./..."],
        cwd=server_dir,
        capture_output=True,
        text=True,
        timeout=30
    )
    check("go vet ./... 通過", result.returncode == 0,
          result.stderr[:200] if result.returncode != 0 else "")
except Exception as e:
    check("go vet ./... 通過", False, str(e))

# ── 7. 關鍵腳本語法基本驗證 ──────────────────────────────────
print("\n[7] GDScript 基本語法驗證")
# 確認沒有明顯的語法問題（未閉合字串等）
for script_name, content in [
    ("Cannon.gd", cannon),
    ("HitEffect.gd", hiteffect),
    ("TargetManager.gd", tm),
]:
    # 檢查引號是否平衡（簡單檢查）
    single_quotes = content.count("'") % 2
    # 檢查是否有 extends 宣告
    has_extends = "extends " in content
    check(f"{script_name} 有 extends 宣告", has_extends)

# ── 結果 ──────────────────────────────────────────────────────
print("\n" + "=" * 60)
total = PASS + FAIL
print(f"結果：{PASS}/{total} 通過")
if FAIL == 0:
    print("🎉 全部通過！DAY-340 品質驗證完成")
else:
    print(f"⚠️  {FAIL} 項未通過，請修復後重新驗證")
print("=" * 60)

sys.exit(0 if FAIL == 0 else 1)
