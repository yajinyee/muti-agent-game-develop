# -*- coding: utf-8 -*-
"""分析規格提案圖，了解角色設計"""
from PIL import Image
import os

spec_dir = r"D:\Kiro\docs\規格提案"

for i in range(1, 10):
    img = Image.open(os.path.join(spec_dir, f"{i}.jpg")).convert("RGB")
    print(f"\n=== {i}.jpg ({img.width}x{img.height}) ===")
    
    # 找白色/奶白色區域（可能是角色）
    pixels = img.load()
    white_regions = []
    for y in range(0, img.height, 10):
        for x in range(0, img.width, 10):
            r, g, b = pixels[x, y]
            if r > 220 and g > 220 and b > 220:
                white_regions.append((x, y))
    
    if white_regions:
        xs = [p[0] for p in white_regions]
        ys = [p[1] for p in white_regions]
        print(f"  White regions: x={min(xs)}-{max(xs)}, y={min(ys)}-{max(ys)}")
    
    # 找粉紅色（吉伊卡哇腮紅）
    pink_regions = []
    for y in range(0, img.height, 5):
        for x in range(0, img.width, 5):
            r, g, b = pixels[x, y]
            if r > 200 and b > 150 and g < r - 30:
                pink_regions.append((x, y))
    
    if pink_regions:
        xs = [p[0] for p in pink_regions]
        ys = [p[1] for p in pink_regions]
        print(f"  Pink regions: x={min(xs)}-{max(xs)}, y={min(ys)}-{max(ys)}")
    
    # 儲存縮圖
    thumb = img.resize((640, 360), Image.LANCZOS)
    thumb.save(os.path.join(r"D:\Kiro\docs", f"spec_{i}_thumb.jpg"))
