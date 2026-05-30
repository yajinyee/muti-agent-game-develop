#!/usr/bin/env python3
"""檢查基礎目標物的視覺狀態"""
from PIL import Image
import os

BASE = "client/chiikawa-pixel/assets/sprites/targets"
targets = [
    ("T001_grass", "草"),
    ("T002_bug_g", "綠蟲"),
    ("T003_bug_r", "紅蟲"),
    ("T004_bug_b", "藍蟲"),
    ("T005_pudding", "布丁"),
    ("T006_mushroom", "蘑菇"),
    ("T101_mimic", "擬態箱"),
    ("T102_chest", "寶箱"),
    ("T103_meteor", "隕石"),
    ("T104_gold_grass", "金草"),
    ("T105_coin_fish", "金幣魚"),
]

print("基礎目標物視覺狀態檢查")
print("=" * 60)
for fname, name in targets:
    path = os.path.join(BASE, fname + ".png")
    if os.path.exists(path):
        img = Image.open(path).convert("RGBA")
        w, h = img.size
        pixels = img.load()
        non_transparent = sum(1 for y in range(h) for x in range(w) if pixels[x, y][3] > 10)
        density = non_transparent / (w * h) * 100
        status = "✅" if density > 30 else "⚠️"
        print(f"{status} {fname} ({name}): {w}x{h}, 密度={density:.1f}%")
    else:
        print(f"❌ {fname} ({name}): 檔案不存在")
