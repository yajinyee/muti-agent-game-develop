"""
強化剩餘低分目標物：T101_mimic 和 boss_bg
"""
from PIL import Image, ImageEnhance, ImageFilter
import numpy as np
import os
import shutil

TARGETS_DIR = 'client/chiikawa-pixel/assets/sprites/targets'
BG_DIR = 'client/chiikawa-pixel/assets/sprites/backgrounds'

def get_bbox(arr):
    alpha = arr[:, :, 3]
    rows = np.any(alpha > 10, axis=1)
    cols = np.any(alpha > 10, axis=0)
    if not rows.any():
        return None
    rmin, rmax = np.where(rows)[0][[0, -1]]
    cmin, cmax = np.where(cols)[0][[0, -1]]
    return rmin, rmax, cmin, cmax

def enhance_t101():
    """強化 T101 擬態怪物 - 加入更豐富的偽裝紋理"""
    path = os.path.join(TARGETS_DIR, 'T101_mimic.png')
    backup = os.path.join(TARGETS_DIR, 'T101_mimic_backup.png')
    if not os.path.exists(backup):
        shutil.copy2(path, backup)
    
    img = Image.open(path).convert('RGBA')
    arr = np.array(img).astype(float)
    
    bbox = get_bbox(arr.astype(np.uint8))
    if bbox is None:
        return
    
    rmin, rmax, cmin, cmax = bbox
    cy = (rmin + rmax) / 2
    cx = (cmin + cmax) / 2
    
    import random
    rng = random.Random(101)
    
    for y in range(rmin, rmax + 1):
        for x in range(cmin, cmax + 1):
            if arr[y, x, 3] < 10:
                continue
            
            # 光照
            dy = (y - cy) / max(1, (rmax - rmin) / 2)
            dx = (x - cx) / max(1, (cmax - cmin) / 2)
            light = (-dx * 0.35 - dy * 0.35)
            
            # 顏色變化
            noise = rng.randint(-12, 12)
            
            for c in range(3):
                val = arr[y, x, c]
                if light > 0:
                    val = min(255, val + light * 45 + noise)
                else:
                    val = max(0, val + light * 55 + noise)
                arr[y, x, c] = val
    
    # 加強輪廓
    result = arr.astype(np.uint8)
    h, w = result.shape[:2]
    orig = result.copy()
    for y in range(1, h - 1):
        for x in range(1, w - 1):
            if orig[y, x, 3] < 10:
                neighbors = [orig[y-1,x,3], orig[y+1,x,3], orig[y,x-1,3], orig[y,x+1,3]]
                if any(n > 100 for n in neighbors):
                    result[y, x] = [20, 20, 20, 200]
    
    orig_colors = 13
    new_colors = len(set(tuple(result[y,x,:3]) for y in range(h) for x in range(w) if result[y,x,3] > 10))
    
    Image.fromarray(result).save(path)
    print(f"  ✅ T101_mimic.png: 顏色 {orig_colors} → {new_colors} (+{new_colors - orig_colors})")

def enhance_boss_bg():
    """強化 boss_bg - 加入更多細節和顏色變化"""
    path = os.path.join(BG_DIR, 'boss_bg.png')
    backup = os.path.join(BG_DIR, 'boss_bg_backup.png')
    if not os.path.exists(backup):
        shutil.copy2(path, backup)
    
    img = Image.open(path).convert('RGB')
    
    # 提升對比度和飽和度
    img = ImageEnhance.Contrast(img).enhance(1.3)
    img = ImageEnhance.Color(img).enhance(1.4)
    
    # 加入輕微銳化
    img = img.filter(ImageFilter.UnsharpMask(radius=1, percent=120, threshold=3))
    
    arr = np.array(img)
    import random
    rng = random.Random(999)
    
    # 加入細微紋理變化
    h, w = arr.shape[:2]
    for y in range(0, h, 4):
        for x in range(0, w, 4):
            noise = rng.randint(-5, 5)
            for dy in range(min(4, h-y)):
                for dx in range(min(4, w-x)):
                    for c in range(3):
                        arr[y+dy, x+dx, c] = max(0, min(255, int(arr[y+dy, x+dx, c]) + noise))
    
    orig_colors = 155
    new_colors = len(set(tuple(arr[y,x]) for y in range(0, h, 2) for x in range(0, w, 2)))
    
    Image.fromarray(arr).save(path)
    print(f"  ✅ boss_bg.png: 顏色 {orig_colors} → {new_colors}+ (採樣)")

print("強化剩餘低分目標物...")
enhance_t101()
enhance_boss_bg()
print("完成！")
