# -*- coding: utf-8 -*-
"""
分析規格提案圖中的角色設計
找出角色的位置、大小、比例
"""
from PIL import Image
import os

spec_dir = r"D:\Kiro\docs\規格提案"

for i in range(1, 10):
    img = Image.open(os.path.join(spec_dir, f"{i}.jpg")).convert("RGB")
    pixels = img.load()
    w, h = img.size
    
    # 找白色角色區域（吉伊卡哇是白色的）
    # 掃描找到白色密集區域
    white_clusters = []
    
    # 用 10x10 格子掃描
    for gy in range(0, h, 10):
        for gx in range(0, w, 10):
            white_count = 0
            for dy in range(10):
                for dx in range(10):
                    px_x, px_y = gx+dx, gy+dy
                    if px_x < w and px_y < h:
                        r, g, b = pixels[px_x, px_y]
                        if r > 230 and g > 230 and b > 230:
                            white_count += 1
            if white_count > 60:  # 60% 以上是白色
                white_clusters.append((gx, gy))
    
    if white_clusters:
        xs = [c[0] for c in white_clusters]
        ys = [c[1] for c in white_clusters]
        print(f"{i}.jpg: white clusters at x={min(xs)}-{max(xs)}, y={min(ys)}-{max(ys)}")
        
        # 裁切角色區域
        margin = 20
        x1 = max(0, min(xs) - margin)
        y1 = max(0, min(ys) - margin)
        x2 = min(w, max(xs) + 30 + margin)
        y2 = min(h, max(ys) + 30 + margin)
        
        crop = img.crop((x1, y1, x2, y2))
        crop.save(os.path.join(r"D:\Kiro\docs", f"spec_{i}_char.jpg"))
        print(f"  Saved spec_{i}_char.jpg ({crop.width}x{crop.height})")
    else:
        print(f"{i}.jpg: no white clusters found")
