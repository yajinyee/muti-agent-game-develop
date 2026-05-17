# -*- coding: utf-8 -*-
"""
從 usagi_idle 生成 attack 幀，確保寬度一致
策略：旋轉 -15 度（揮棒感），保持 bbox 寬度接近 idle
"""
from PIL import Image, ImageEnhance
import numpy as np
import os

CHARS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"
CELL_SIZE = 96

def fix_usagi_attack():
    # 載入 idle
    idle_path = os.path.join(CHARS_DIR, "usagi_idle.png")
    img = Image.open(idle_path).convert("RGBA")
    idle_bbox = img.getbbox()
    print(f"usagi_idle: bbox={idle_bbox}, size={idle_bbox[2]-idle_bbox[0]}x{idle_bbox[3]-idle_bbox[1]}")
    
    # attack 效果：
    # 1. 水平翻轉（面向右邊，揮棒感）
    # 2. 亮度提高
    # 3. 粉紅色光暈（討伐棒效果，只在角色本體上）
    
    # 水平翻轉（bbox 不變）
    flipped = img.transpose(Image.FLIP_LEFT_RIGHT)
    
    # 亮度提高
    flipped = ImageEnhance.Brightness(flipped).enhance(1.15)
    
    # 加粉紅色光暈（右上角，模擬討伐棒揮出）
    arr = np.array(flipped).copy()
    cx, cy = 70, 25  # 光暈中心（右上角）
    radius = 10
    for y in range(max(0, cy-radius), min(CELL_SIZE, cy+radius)):
        for x in range(max(0, cx-radius), min(CELL_SIZE, cx+radius)):
            dist = ((x-cx)**2 + (y-cy)**2) ** 0.5
            if dist < radius:
                alpha = int(120 * (1 - dist/radius))
                # 粉紅色光暈疊加（只在非透明像素上）
                if arr[y, x, 3] > 10:
                    arr[y, x, 0] = min(255, int(arr[y, x, 0]) + alpha)
                    arr[y, x, 1] = max(0, int(arr[y, x, 1]) - alpha//3)
                    arr[y, x, 2] = min(255, int(arr[y, x, 2]) + alpha//2)
    
    canvas = Image.fromarray(arr.copy())
    
    # 儲存
    out_path = os.path.join(CHARS_DIR, "usagi_attack.png")
    canvas.save(out_path)
    
    # 確認
    new_bbox = canvas.getbbox()
    if new_bbox:
        nw = new_bbox[2] - new_bbox[0]
        nh = new_bbox[3] - new_bbox[1]
        idle_w = idle_bbox[2] - idle_bbox[0]
        idle_h = idle_bbox[3] - idle_bbox[1]
        h_diff = abs(nh - idle_h)
        w_diff = abs(nw - idle_w)
        status = "✅" if h_diff <= 2 and w_diff <= 4 else "⚠️ "
        print(f"usagi_attack: {nw}x{nh}px, {status} vs idle: h_diff={h_diff}, w_diff={w_diff}")

if __name__ == "__main__":
    fix_usagi_attack()
    print("Done!")
