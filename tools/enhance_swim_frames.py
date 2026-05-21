"""
同步強化游泳動畫幀（_swim.png）
使用與靜態幀相同的強化邏輯
"""
from PIL import Image
import numpy as np
import os
import shutil

TARGETS_DIR = 'client/chiikawa-pixel/assets/sprites/targets'

def get_bbox(arr):
    alpha = arr[:, :, 3]
    rows = np.any(alpha > 10, axis=1)
    cols = np.any(alpha > 10, axis=0)
    if not rows.any():
        return None
    rmin, rmax = np.where(rows)[0][[0, -1]]
    cmin, cmax = np.where(cols)[0][[0, -1]]
    return rmin, rmax, cmin, cmax

def enhance_swim(filename):
    path = os.path.join(TARGETS_DIR, filename)
    if not os.path.exists(path):
        return
    
    backup = path.replace('.png', '_backup.png')
    if not os.path.exists(backup):
        shutil.copy2(path, backup)
    
    img = Image.open(path).convert('RGBA')
    arr = np.array(img)
    
    # swim 幀是 spritesheet（多幀橫排），需要分幀處理
    w, h = img.size
    
    # 估算幀數（通常是 4 幀）
    frame_w = h  # 假設幀是正方形
    n_frames = w // frame_w if frame_w > 0 else 1
    
    import random
    rng = random.Random(42)
    
    result = arr.copy().astype(float)
    
    for frame_idx in range(n_frames):
        x_start = frame_idx * frame_w
        x_end = x_start + frame_w
        frame = arr[:, x_start:x_end, :]
        
        bbox = get_bbox(frame)
        if bbox is None:
            continue
        
        rmin, rmax, cmin, cmax = bbox
        cy = (rmin + rmax) / 2
        cx = (cmin + cmax) / 2
        
        for y in range(rmin, rmax + 1):
            for x in range(cmin, cmax + 1):
                gx = x_start + x
                if arr[y, gx, 3] < 10:
                    continue
                
                dy = (y - cy) / max(1, (rmax - rmin) / 2)
                dx = (x - cx) / max(1, (cmax - cmin) / 2)
                light = (-dx * 0.3 - dy * 0.3)
                noise = rng.randint(-8, 8)
                
                for c in range(3):
                    val = result[y, gx, c]
                    if light > 0:
                        result[y, gx, c] = min(255, val + light * 40 + noise)
                    else:
                        result[y, gx, c] = max(0, val + light * 50 + noise)
    
    # 加強輪廓
    final = result.astype(np.uint8)
    orig = arr.copy()
    for y in range(1, h - 1):
        for x in range(1, w - 1):
            if orig[y, x, 3] < 10:
                neighbors = [orig[y-1,x,3], orig[y+1,x,3], orig[y,x-1,3], orig[y,x+1,3]]
                if any(n > 100 for n in neighbors):
                    final[y, x] = [20, 20, 20, 200]
    
    orig_colors = len(set(tuple(arr[y,x,:3]) for y in range(h) for x in range(w) if arr[y,x,3] > 10))
    new_colors = len(set(tuple(final[y,x,:3]) for y in range(h) for x in range(w) if final[y,x,3] > 10))
    
    Image.fromarray(final).save(path)
    print(f"  ✅ {filename}: {w}x{h} ({n_frames}幀), 顏色 {orig_colors} → {new_colors}")

swim_files = [f for f in os.listdir(TARGETS_DIR) if f.endswith('_swim.png')]

print(f"強化 {len(swim_files)} 個游泳動畫幀...")
for f in sorted(swim_files):
    enhance_swim(f)
print("完成！")
