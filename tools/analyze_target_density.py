# -*- coding: utf-8 -*-
"""分析目標物 sprite 的 bbox 利用率（更準確的密度評估）"""
from PIL import Image
import os

TARGET_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"

def analyze(path):
    img = Image.open(path).convert("RGBA")
    w, h = img.size
    pixels = img.load()

    # 找 bbox（非透明像素的邊界）
    min_x, min_y = w, h
    max_x, max_y = 0, 0
    non_transparent = 0

    for y in range(h):
        for x in range(w):
            if pixels[x, y][3] > 10:
                non_transparent += 1
                min_x = min(min_x, x)
                min_y = min(min_y, y)
                max_x = max(max_x, x)
                max_y = max(max_y, y)

    if non_transparent == 0:
        return 0, 0, 0

    bbox_w = max_x - min_x + 1
    bbox_h = max_y - min_y + 1
    bbox_area = bbox_w * bbox_h
    bbox_util = non_transparent / bbox_area * 100
    canvas_util = non_transparent / (w * h) * 100

    return canvas_util, bbox_util, (bbox_w, bbox_h)

print("=== 目標物密度分析（canvas vs bbox）===\n")
print(f"{'檔案':<30} {'畫布%':>8} {'bbox%':>8} {'bbox尺寸':>12}")
print("-" * 65)

for f in sorted(os.listdir(TARGET_DIR)):
    if not f.endswith('.png'):
        continue
    path = os.path.join(TARGET_DIR, f)
    canvas_pct, bbox_pct, bbox_size = analyze(path)
    canvas_flag = "✅" if canvas_pct >= 30 else "⚠️"
    bbox_flag = "✅" if bbox_pct >= 55 else "⚠️"
    print(f"{canvas_flag}{bbox_flag} {f:<28} {canvas_pct:>7.1f}% {bbox_pct:>7.1f}% {str(bbox_size):>12}")
