#!/usr/bin/env python3
"""
sharpen_bugs.py
針對蟲類目標物進行銳化強化，提升清晰度分數
目標：T002-T005 從 74 → 80+
"""
from pathlib import Path
from PIL import Image, ImageEnhance, ImageFilter
import numpy as np

TARGETS_DIR = Path("d:/Kiro/client/chiikawa-pixel/assets/sprites/targets")

def analyze_sharpness(path):
    img = Image.open(path).convert("RGBA")
    gray = img.convert("L")
    import math
    from PIL import ImageFilter
    lap = gray.filter(ImageFilter.FIND_EDGES)
    lap_pixels = list(lap.getdata())
    sharpness = math.sqrt(sum(p*p for p in lap_pixels) / len(lap_pixels))
    return sharpness

def sharpen_sprite(path, seed=42):
    img = Image.open(path).convert("RGBA")
    alpha = img.split()[3]
    
    before_sharp = analyze_sharpness(path)
    
    rgb = img.convert("RGB")
    
    # 多次銳化
    for _ in range(3):
        rgb = rgb.filter(ImageFilter.SHARPEN)
    
    # Unsharp mask（更強的銳化）
    rgb = rgb.filter(ImageFilter.UnsharpMask(radius=1, percent=200, threshold=2))
    
    # 提升對比度讓邊緣更清晰
    rgb = ImageEnhance.Contrast(rgb).enhance(1.3)
    
    result = Image.new("RGBA", img.size)
    result.paste(rgb, mask=alpha)
    result.save(path)
    
    after_sharp = analyze_sharpness(path)
    name = Path(path).name
    print(f"  {name}: sharpness {before_sharp:.1f} → {after_sharp:.1f} (+{after_sharp-before_sharp:.1f})")

if __name__ == "__main__":
    print("=== 銳化蟲類目標物 ===")
    targets = [
        "T002_bug_g.png",
        "T003_bug_r.png", 
        "T004_bug_b.png",
        "T005_pudding.png",
        "T001_grass.png",
        "T101_mimic.png",
        "T103_meteor.png",
    ]
    for fname in targets:
        path = TARGETS_DIR / fname
        if path.exists():
            sharpen_sprite(path)
    print("=== 完成 ===")
