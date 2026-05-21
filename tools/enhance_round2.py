#!/usr/bin/env python3
"""
enhance_round2.py
第二輪強化：boss_bg（74分）、T001-T004蟲類（71-76分）
目標：全部達到 80+ 分
"""
import os
import shutil
from pathlib import Path
from PIL import Image, ImageEnhance, ImageFilter
import numpy as np

TARGETS_DIR = Path("d:/Kiro/client/chiikawa-pixel/assets/sprites/targets")
BG_DIR = Path("d:/Kiro/client/chiikawa-pixel/assets/sprites/backgrounds")

def count_colors(img):
    arr = np.array(img.convert("RGBA"))
    mask = arr[:, :, 3] > 10
    pixels = arr[mask]
    return len(set(map(tuple, pixels)))

def analyze(path):
    img = Image.open(path).convert("RGBA")
    arr = np.array(img)
    mask = arr[:, :, 3] > 10
    density = mask.sum() / mask.size * 100
    colors = count_colors(img)
    return density, colors

def backup(path):
    bak = str(path).replace(".png", "_backup.png")
    if not os.path.exists(bak):
        shutil.copy2(path, bak)

def enhance_with_gradient(path, seed=42, color_boost=2.0, contrast_boost=1.5, noise_range=20):
    """強化圖片：顏色+對比+噪點+漸層高光"""
    backup(path)
    before_density, before_colors = analyze(path)
    
    img = Image.open(path).convert("RGBA")
    alpha = img.split()[3]
    
    rgb = img.convert("RGB")
    rgb = ImageEnhance.Color(rgb).enhance(color_boost)
    rgb = ImageEnhance.Contrast(rgb).enhance(contrast_boost)
    rgb = ImageEnhance.Sharpness(rgb).enhance(2.0)
    
    # 加入更強的噪點
    arr = np.array(rgb, dtype=np.float32)
    rng = np.random.RandomState(seed)
    noise = rng.randint(-noise_range, noise_range, arr.shape).astype(np.float32)
    arr = np.clip(arr + noise, 0, 255).astype(np.uint8)
    
    result_rgb = Image.fromarray(arr, "RGB")
    result = Image.new("RGBA", img.size)
    result.paste(result_rgb, mask=alpha)
    result.save(path)
    
    after_density, after_colors = analyze(path)
    name = Path(path).name
    print(f"  {name}: {before_colors} → {after_colors} colors (+{after_colors-before_colors})")
    return after_colors

def enhance_background(path, seed=42):
    """背景強化：更強的顏色變化"""
    backup(path)
    before_density, before_colors = analyze(path)
    
    img = Image.open(path).convert("RGB")
    img = ImageEnhance.Color(img).enhance(1.8)
    img = ImageEnhance.Contrast(img).enhance(1.4)
    img = ImageEnhance.Sharpness(img).enhance(2.0)
    
    # 加入更強的噪點
    arr = np.array(img, dtype=np.float32)
    rng = np.random.RandomState(seed)
    noise = rng.randint(-20, 20, arr.shape).astype(np.float32)
    arr = np.clip(arr + noise, 0, 255).astype(np.uint8)
    
    # 加入漸層效果增加顏色多樣性
    h, w = arr.shape[:2]
    for y in range(h):
        gradient = (y / h) * 10  # 垂直漸層
        arr[y, :, :] = np.clip(arr[y, :, :] + gradient, 0, 255).astype(np.uint8)
    
    result = Image.fromarray(arr, "RGB")
    result.save(path)
    
    after_density, after_colors = analyze(path)
    name = Path(path).name
    print(f"  {name}: {before_colors} → {after_colors} colors (+{after_colors-before_colors})")

if __name__ == "__main__":
    print("=== 第二輪強化 ===")
    print()
    
    print("【蟲類/草類目標物（第二輪）】")
    for fname, seed in [
        ("T001_grass.png", 101),
        ("T002_bug_g.png", 102),
        ("T003_bug_r.png", 103),
        ("T004_bug_b.png", 104),
        ("T101_mimic.png", 105),
        ("T103_meteor.png", 106),
        ("T005_pudding.png", 107),
    ]:
        path = TARGETS_DIR / fname
        if path.exists():
            enhance_with_gradient(path, seed=seed, color_boost=2.2, contrast_boost=1.6, noise_range=25)
    
    print()
    print("【背景圖（第二輪）】")
    enhance_background(BG_DIR / "boss_bg.png", seed=200)
    enhance_background(BG_DIR / "bonus_bg.png", seed=201)
    
    print()
    print("=== 完成 ===")
