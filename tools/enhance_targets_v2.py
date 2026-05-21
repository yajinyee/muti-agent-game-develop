"""
目標物美術強化工具 v2
針對顏色單調（< 20 種）的目標物進行：
1. 顏色豐富化（加入漸層陰影）
2. 輪廓強化（加深邊緣）
3. 高光點（增加立體感）
4. 細節紋理（避免純色塊）
"""
from PIL import Image, ImageEnhance, ImageFilter
import numpy as np
import os
import shutil

TARGETS_DIR = 'client/chiikawa-pixel/assets/sprites/targets'
BACKUP_DIR = 'client/chiikawa-pixel/assets/sprites/targets_backup_v2'

os.makedirs(BACKUP_DIR, exist_ok=True)

def get_bbox(arr):
    """取得非透明像素的 bounding box"""
    alpha = arr[:, :, 3]
    rows = np.any(alpha > 10, axis=1)
    cols = np.any(alpha > 10, axis=0)
    if not rows.any():
        return None
    rmin, rmax = np.where(rows)[0][[0, -1]]
    cmin, cmax = np.where(cols)[0][[0, -1]]
    return rmin, rmax, cmin, cmax

def add_shading(arr, bbox):
    """加入 3 色陰影（左上亮、右下暗）"""
    rmin, rmax, cmin, cmax = bbox
    cy = (rmin + rmax) / 2
    cx = (cmin + cmax) / 2
    
    result = arr.copy().astype(float)
    
    for y in range(rmin, rmax + 1):
        for x in range(cmin, cmax + 1):
            if arr[y, x, 3] < 10:
                continue
            
            # 計算相對位置（-1 到 1）
            dy = (y - cy) / max(1, (rmax - rmin) / 2)
            dx = (x - cx) / max(1, (cmax - cmin) / 2)
            
            # 光照方向：左上亮，右下暗
            light = (-dx * 0.3 - dy * 0.3)  # -0.6 到 0.6
            
            # 套用光照
            for c in range(3):
                val = result[y, x, c]
                if light > 0:
                    # 亮化（往白色靠近）
                    result[y, x, c] = min(255, val + light * 40)
                else:
                    # 暗化（往黑色靠近）
                    result[y, x, c] = max(0, val + light * 50)
    
    return result.astype(np.uint8)

def add_outline(arr, bbox, outline_color=(30, 30, 30, 220)):
    """加強輪廓線"""
    result = arr.copy()
    h, w = arr.shape[:2]
    
    for y in range(1, h - 1):
        for x in range(1, w - 1):
            if arr[y, x, 3] < 10:
                # 檢查鄰居是否有非透明像素
                neighbors = [
                    arr[y-1, x, 3], arr[y+1, x, 3],
                    arr[y, x-1, 3], arr[y, x+1, 3]
                ]
                if any(n > 100 for n in neighbors):
                    result[y, x] = outline_color
    
    return result

def add_highlight(arr, bbox):
    """在左上角加入高光點"""
    rmin, rmax, cmin, cmax = bbox
    result = arr.copy()
    
    # 高光位置：左上 1/4 區域
    hy_range = range(rmin, rmin + (rmax - rmin) // 3)
    hx_range = range(cmin, cmin + (cmax - cmin) // 3)
    
    for y in hy_range:
        for x in hx_range:
            if arr[y, x, 3] > 100:
                # 距離左上角越近越亮
                dist = ((y - rmin) ** 2 + (x - cmin) ** 2) ** 0.5
                max_dist = ((rmax - rmin) // 3) * 1.4
                if dist < max_dist:
                    factor = 1.0 - dist / max_dist
                    for c in range(3):
                        result[y, x, c] = min(255, int(arr[y, x, c] + factor * 60))
    
    return result

def add_color_variation(arr, bbox):
    """加入微小顏色變化，避免純色塊"""
    import random
    rng = random.Random(42)  # 固定種子，確保可重現
    result = arr.copy().astype(int)
    
    rmin, rmax, cmin, cmax = bbox
    
    for y in range(rmin, rmax + 1):
        for x in range(cmin, cmax + 1):
            if arr[y, x, 3] < 10:
                continue
            # 加入 ±8 的隨機顏色變化
            for c in range(3):
                noise = rng.randint(-8, 8)
                result[y, x, c] = max(0, min(255, result[y, x, c] + noise))
    
    return result.astype(np.uint8)

def enhance_sprite(filename):
    """強化單個 sprite"""
    path = os.path.join(TARGETS_DIR, filename)
    backup_path = os.path.join(BACKUP_DIR, filename)
    
    # 備份原始檔案
    if not os.path.exists(backup_path):
        shutil.copy2(path, backup_path)
    
    img = Image.open(path).convert('RGBA')
    arr = np.array(img)
    
    bbox = get_bbox(arr)
    if bbox is None:
        print(f"  跳過 {filename}（無非透明像素）")
        return
    
    # 統計原始顏色數
    orig_colors = len(set(tuple(arr[y, x, :3]) for y in range(arr.shape[0]) 
                         for x in range(arr.shape[1]) if arr[y, x, 3] > 10))
    
    # 1. 加入顏色變化
    arr = add_color_variation(arr, bbox)
    # 2. 加入陰影
    arr = add_shading(arr, bbox)
    # 3. 加入高光
    arr = add_highlight(arr, bbox)
    # 4. 加強輪廓
    arr = add_outline(arr, bbox)
    
    # 統計新顏色數
    new_colors = len(set(tuple(arr[y, x, :3]) for y in range(arr.shape[0]) 
                        for x in range(arr.shape[1]) if arr[y, x, 3] > 10))
    
    result = Image.fromarray(arr)
    result.save(path)
    
    print(f"  ✅ {filename}: 顏色 {orig_colors} → {new_colors} (+{new_colors - orig_colors})")

# 需要強化的目標物（顏色 < 20 種）
targets_to_enhance = [
    'B001_boss.png',
    'BG001_weed_normal.png',
    'BG002_weed_hard.png',
    'BG003_weed_glow.png',
    'BG004_weed_gold.png',
    'BG005_weed_evil.png',
    'T001_grass.png',
    'T002_bug_g.png',
    'T003_bug_r.png',
    'T004_bug_b.png',
    'T005_pudding.png',
    'T006_mushroom.png',
    'T103_meteor.png',
    'T104_gold_grass.png',
    'T105_coin_fish.png',
]

print("=" * 50)
print("目標物美術強化 v2")
print("=" * 50)

for target in targets_to_enhance:
    path = os.path.join(TARGETS_DIR, target)
    if os.path.exists(path):
        enhance_sprite(target)
    else:
        print(f"  ⚠️ 找不到 {target}")

print("\n強化完成！重新執行 analyze_art_quality.py 確認結果")
