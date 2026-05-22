# -*- coding: utf-8 -*-
"""
從 targets_sheet.json 自動生成 TargetManager.gd 的 SHEET_REGIONS
只包含主要目標物（不含 backup/swim 變體）
"""
import json, re

JSON_PATH = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\sheets\targets_sheet.json"
GD_PATH   = r"D:\Kiro\client\chiikawa-pixel\scripts\game\TargetManager.gd"

# 只保留主要目標物（不含 backup/swim/swim_backup）
VALID_IDS = [
    "T001","T002","T003","T004","T005","T006",
    "T101","T102","T103","T104","T105",
    "T106","T107","T108","T109","T110",
    "T111","T112","T113","T114","T115","T116","T117",
]
# T118-T126 用獨立 PNG，不加入 SHEET_REGIONS

with open(JSON_PATH, encoding="utf-8") as f:
    data = json.load(f)

sprites = data["sprites"]

# 建立 SHEET_REGIONS 字串
lines = []
for tid in VALID_IDS:
    # 找對應的 sprite key（可能是 T001_grass 等）
    matched = None
    for key, val in sprites.items():
        # 只匹配主要版本（不含 backup/swim）
        if key.startswith(tid + "_") and "backup" not in key and "swim" not in key:
            matched = val
            break
    if matched:
        lines.append(f'\t"{tid}": Rect2({matched["x"]}, {matched["y"]}, 64, 64),')
    else:
        print(f"  WARNING: {tid} not found in JSON")

# 加入 T118-T126 的說明
lines.append("\t# T118-T126 使用獨立 PNG（TARGET_SPRITES 備用路徑），不在 Spritesheet 中")

new_regions = "const SHEET_REGIONS = {\n" + "\n".join(lines) + "\n}"

# 讀取 GD 檔案
with open(GD_PATH, encoding="utf-8") as f:
    content = f.read()

# 替換 SHEET_REGIONS
pattern = r"const SHEET_REGIONS = \{[^}]+\}"
new_content = re.sub(pattern, new_regions, content, flags=re.DOTALL)

with open(GD_PATH, "w", encoding="utf-8") as f:
    f.write(new_content)

print("SHEET_REGIONS 已更新！")
for line in lines[:5]:
    print(" ", line)
print("  ...")
