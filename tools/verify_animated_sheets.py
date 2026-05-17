# -*- coding: utf-8 -*-
"""
驗證 animated sheets 格式是否符合 CharacterAnimator.gd 的期望
期望格式：
  - 384x288（4幀 × 3狀態 × 96px）
  - 行0：idle（4幀，4fps）
  - 行1：attack（3幀，8fps）
  - 行2：bigwin（4幀，6fps）
  - 每幀 96x96，RGBA
"""
from PIL import Image
import os
import json

SHEET_DIR = "D:/Kiro/client/chiikawa-pixel/assets/sprites/sheets"
FRAME_SIZE = 96
EXPECTED_W = 384  # 4 cols × 96
EXPECTED_H = 288  # 3 rows × 96

print("=== Animated Sheet 驗證 ===\n")

all_ok = True
for char in ["chiikawa", "hachiware", "usagi"]:
    sheet_path = os.path.join(SHEET_DIR, f"{char}_animated.png")
    meta_path  = os.path.join(SHEET_DIR, f"{char}_animated.json")

    issues = []

    # 檢查檔案存在
    if not os.path.exists(sheet_path):
        print(f"❌ {char}: sheet NOT FOUND")
        all_ok = False
        continue

    # 檢查尺寸
    sheet = Image.open(sheet_path).convert("RGBA")
    if sheet.width != EXPECTED_W or sheet.height != EXPECTED_H:
        issues.append(f"尺寸錯誤 {sheet.width}x{sheet.height}（期望 {EXPECTED_W}x{EXPECTED_H}）")

    # 檢查各幀有內容
    state_names = ["idle", "attack", "bigwin"]
    frame_counts = [4, 3, 4]
    for row, (state, n_frames) in enumerate(zip(state_names, frame_counts)):
        for col in range(n_frames):
            x = col * FRAME_SIZE
            y = row * FRAME_SIZE
            frame = sheet.crop((x, y, x + FRAME_SIZE, y + FRAME_SIZE))
            pixels_list = list(frame.getdata())
            non_t = sum(1 for px in pixels_list if px[3] > 10)
            if non_t < 100:
                issues.append(f"{state} 幀{col} 幾乎空白（{non_t} 像素）")

    # 檢查 metadata
    if not os.path.exists(meta_path):
        issues.append("metadata JSON 不存在")
    else:
        with open(meta_path) as f:
            meta = json.load(f)
        expected_keys = ["idle", "attack", "bigwin"]
        for k in expected_keys:
            if k not in meta.get("animations", {}):
                issues.append(f"metadata 缺少 {k}")

    if issues:
        print(f"⚠️  {char}:")
        for issue in issues:
            print(f"   - {issue}")
        all_ok = False
    else:
        print(f"✅ {char}: {sheet.width}x{sheet.height}, RGBA, 所有幀有內容")

print()
if all_ok:
    print("✅ 所有 animated sheets 格式正確，可供 Godot 使用")
else:
    print("⚠️  部分 sheets 有問題，需要修正")

print()
print("Godot 使用方式（CharacterAnimator.gd）：")
print("  CHAR_SHEETS = {")
print('    "chiikawa":  "res://assets/sprites/sheets/chiikawa_animated.png",')
print('    "hachiware": "res://assets/sprites/sheets/hachiware_animated.png",')
print('    "usagi":     "res://assets/sprites/sheets/usagi_animated.png",')
print("  }")
print("  FRAME_SIZE = 96, COLS = 4")
print("  idle: row=0, frames=4, fps=4.0")
print("  attack: row=1, frames=3, fps=8.0")
print("  bigwin: row=2, frames=4, fps=6.0")
