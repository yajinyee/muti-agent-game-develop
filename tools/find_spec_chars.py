# -*- coding: utf-8 -*-
"""找規格提案圖中的角色位置"""
from PIL import Image
import os
from collections import Counter

spec_dir = r"D:\Kiro\docs\規格提案"
out_dir = r"D:\Kiro\docs"

for i in range(1, 10):
    img = Image.open(os.path.join(spec_dir, f"{i}.jpg")).convert("RGB")
    pixels = img.load()
    w, h = img.size
    
    # 找非白色區域
    colored = []
    for y in range(0, h, 3):
        for x in range(0, w, 3):
            r, g, b = pixels[x, y]
            if not (r > 200 and g > 200 and b > 200):
                colored.append((x, y))
    
    if not colored:
        print(f"{i}.jpg: all white")
        continue
    
    xs = [c[0] for c in colored]
    ys = [c[1] for c in colored]
    print(f"{i}.jpg: colored x={min(xs)}-{max(xs)}, y={min(ys)}-{max(ys)}, count={len(colored)}")
    
    # 找密集區域
    grid = Counter()
    for x, y in colored:
        grid[(x//80, y//80)] += 1
    
    top = grid.most_common(3)
    for (gx, gy), cnt in top:
        print(f"  Dense at ({gx*80},{gy*80}): {cnt} pixels")
    
    # 裁切最密集的區域
    if top:
        (gx, gy), _ = top[0]
        x1 = max(0, gx*80 - 40)
        y1 = max(0, gy*80 - 40)
        x2 = min(w, gx*80 + 120)
        y2 = min(h, gy*80 + 120)
        crop = img.crop((x1, y1, x2, y2))
        crop.save(os.path.join(out_dir, f"spec_{i}_dense.jpg"))
