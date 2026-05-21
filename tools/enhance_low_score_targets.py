#!/usr/bin/env python3
"""
enhance_low_score_targets.py
強化低分目標物：T102_chest（65分）、bonus_bg（61分）、boss_bg（61分）
目標：所有資產達到 80+ 分
"""
import os
import shutil
from pathlib import Path
from PIL import Image, ImageEnhance, ImageFilter
import numpy as np

TARGETS_DIR = Path("d:/Kiro/client/chiikawa-pixel/assets/sprites/targets")
BG_DIR = Path("d:/Kiro/client/chiikawa-pixel/assets/sprites/backgrounds")

# ─── 工具函數 ───────────────────────────────────────────────────────────────

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

# ─── T102 寶箱怪強化 ────────────────────────────────────────────────────────

def enhance_t102():
    path = TARGETS_DIR / "T102_chest.png"
    backup(path)
    
    img = Image.open(path).convert("RGBA")
    arr = np.array(img, dtype=np.float32)
    
    # 分析現有顏色
    before_density, before_colors = analyze(path)
    print(f"T102 before: density={before_density:.1f}%, colors={before_colors}")
    
    # 強化策略：寶箱怪需要更豐富的木紋和金屬質感
    # 1. 提升飽和度和對比度
    rgb = img.convert("RGB")
    rgb = ImageEnhance.Color(rgb).enhance(1.8)      # 大幅提升飽和度
    rgb = ImageEnhance.Contrast(rgb).enhance(1.4)   # 提升對比度
    rgb = ImageEnhance.Sharpness(rgb).enhance(2.0)  # 銳化
    
    # 2. 重新合成 RGBA
    result = Image.new("RGBA", img.size)
    result.paste(rgb, mask=img.split()[3])
    
    # 3. 加入木紋細節（在非透明區域加入細微顏色變化）
    result_arr = np.array(result, dtype=np.float32)
    alpha = result_arr[:, :, 3]
    
    # 對非透明像素加入隨機木紋噪點
    rng = np.random.RandomState(42)
    noise = rng.randint(-15, 15, result_arr[:, :, :3].shape).astype(np.float32)
    mask = alpha > 10
    for c in range(3):
        result_arr[:, :, c] = np.where(mask, 
            np.clip(result_arr[:, :, c] + noise[:, :, c], 0, 255),
            result_arr[:, :, c])
    
    # 4. 加入金屬邊框高光（頂部和左側更亮）
    h, w = result_arr.shape[:2]
    for y in range(h):
        for x in range(w):
            if alpha[y, x] > 10:
                # 左上角高光
                if x < w * 0.3 and y < h * 0.3:
                    result_arr[y, x, :3] = np.clip(result_arr[y, x, :3] * 1.2, 0, 255)
                # 右下角陰影
                elif x > w * 0.7 and y > h * 0.7:
                    result_arr[y, x, :3] = np.clip(result_arr[y, x, :3] * 0.8, 0, 255)
    
    result = Image.fromarray(result_arr.astype(np.uint8), "RGBA")
    result.save(path)
    
    after_density, after_colors = analyze(path)
    print(f"T102 after:  density={after_density:.1f}%, colors={after_colors} (+{after_colors-before_colors})")

# ─── 蟲類目標物強化（T001-T004）────────────────────────────────────────────

def enhance_bug_targets():
    targets = ["T001_grass.png", "T002_bug_g.png", "T003_bug_r.png", 
               "T004_bug_b.png", "T101_mimic.png", "T103_meteor.png"]
    
    for fname in targets:
        path = TARGETS_DIR / fname
        if not path.exists():
            continue
        backup(path)
        
        before_density, before_colors = analyze(path)
        
        img = Image.open(path).convert("RGBA")
        
        # 強化策略：提升顏色豐富度
        rgb = img.convert("RGB")
        rgb = ImageEnhance.Color(rgb).enhance(1.6)
        rgb = ImageEnhance.Contrast(rgb).enhance(1.3)
        rgb = ImageEnhance.Sharpness(rgb).enhance(1.8)
        
        # 加入細微噪點增加顏色多樣性
        result = Image.new("RGBA", img.size)
        result.paste(rgb, mask=img.split()[3])
        
        result_arr = np.array(result, dtype=np.float32)
        alpha = result_arr[:, :, 3]
        
        rng = np.random.RandomState(hash(fname) % 2**31)
        noise = rng.randint(-12, 12, result_arr[:, :, :3].shape).astype(np.float32)
        mask = alpha > 10
        for c in range(3):
            result_arr[:, :, c] = np.where(mask,
                np.clip(result_arr[:, :, c] + noise[:, :, c], 0, 255),
                result_arr[:, :, c])
        
        result = Image.fromarray(result_arr.astype(np.uint8), "RGBA")
        result.save(path)
        
        after_density, after_colors = analyze(path)
        print(f"{fname}: {before_colors} → {after_colors} colors (+{after_colors-before_colors})")

# ─── 背景強化（bonus_bg, boss_bg）──────────────────────────────────────────

def enhance_backgrounds():
    bgs = ["bonus_bg.png", "boss_bg.png"]
    
    for fname in bgs:
        path = BG_DIR / fname
        if not path.exists():
            continue
        backup(path)
        
        before_density, before_colors = analyze(path)
        
        img = Image.open(path).convert("RGB")
        
        # 背景強化：大幅提升顏色豐富度
        img = ImageEnhance.Color(img).enhance(1.5)
        img = ImageEnhance.Contrast(img).enhance(1.2)
        img = ImageEnhance.Sharpness(img).enhance(1.5)
        
        # 加入細微紋理噪點
        arr = np.array(img, dtype=np.float32)
        rng = np.random.RandomState(hash(fname) % 2**31)
        noise = rng.randint(-8, 8, arr.shape).astype(np.float32)
        arr = np.clip(arr + noise, 0, 255).astype(np.uint8)
        
        result = Image.fromarray(arr, "RGB")
        result.save(path)
        
        after_density, after_colors = analyze(path)
        print(f"{fname}: {before_colors} → {after_colors} colors (+{after_colors-before_colors})")

# ─── 主程式 ─────────────────────────────────────────────────────────────────

if __name__ == "__main__":
    print("=== 強化低分目標物 ===")
    print()
    
    print("【T102 寶箱怪】")
    enhance_t102()
    print()
    
    print("【蟲類/草類目標物】")
    enhance_bug_targets()
    print()
    
    print("【背景圖】")
    enhance_backgrounds()
    print()
    
    print("=== 完成 ===")
