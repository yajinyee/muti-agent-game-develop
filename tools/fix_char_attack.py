# -*- coding: utf-8 -*-
"""
修復 chiikawa 和 hachiware 的 attack 幀一致性問題
策略：從 idle 幀做水平翻轉 + 光暈效果，確保 bbox 與 idle 一致
"""
from PIL import Image, ImageEnhance
import numpy as np
import os

CHARS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"
CELL_SIZE = 96

def fix_attack(char_name, glow_color=(255, 100, 180), glow_cx_offset=20, glow_cy_offset=-25):
    """
    從 idle 幀生成 attack 幀
    glow_color: 光暈顏色 (R, G, B)
    glow_cx_offset: 光暈中心相對角色中心的 X 偏移
    glow_cy_offset: 光暈中心相對角色中心的 Y 偏移
    """
    idle_path = os.path.join(CHARS_DIR, f"{char_name}_idle.png")
    img = Image.open(idle_path).convert("RGBA")
    idle_bbox = img.getbbox()
    idle_w = idle_bbox[2] - idle_bbox[0]
    idle_h = idle_bbox[3] - idle_bbox[1]
    idle_cx = (idle_bbox[0] + idle_bbox[2]) // 2
    idle_cy = (idle_bbox[1] + idle_bbox[3]) // 2
    print(f"{char_name}_idle: bbox={idle_bbox}, size={idle_w}x{idle_h}")

    # 1. 水平翻轉（bbox 不變，確保寬度一致）
    flipped = img.transpose(Image.FLIP_LEFT_RIGHT)

    # 2. 亮度提高（攻擊狀態更亮）
    flipped = ImageEnhance.Brightness(flipped).enhance(1.15)

    # 3. 加光暈（只在非透明像素上，不擴大 bbox）
    arr = np.array(flipped).copy()
    cx = idle_cx + glow_cx_offset
    cy = idle_cy + glow_cy_offset
    radius = 12
    gr, gg, gb = glow_color
    for y in range(max(0, cy - radius), min(CELL_SIZE, cy + radius)):
        for x in range(max(0, cx - radius), min(CELL_SIZE, cx + radius)):
            dist = ((x - cx) ** 2 + (y - cy) ** 2) ** 0.5
            if dist < radius and arr[y, x, 3] > 10:
                alpha = int(100 * (1 - dist / radius))
                arr[y, x, 0] = min(255, int(arr[y, x, 0]) + int(alpha * gr / 255))
                arr[y, x, 1] = max(0, int(arr[y, x, 1]) - alpha // 4)
                arr[y, x, 2] = min(255, int(arr[y, x, 2]) + int(alpha * gb / 255))

    result = Image.fromarray(arr.copy())

    # 儲存
    out_path = os.path.join(CHARS_DIR, f"{char_name}_attack.png")
    result.save(out_path)

    # 驗證
    new_bbox = result.getbbox()
    if new_bbox:
        nw = new_bbox[2] - new_bbox[0]
        nh = new_bbox[3] - new_bbox[1]
        h_diff = abs(nh - idle_h)
        w_diff = abs(nw - idle_w)
        status = "✅" if h_diff <= 2 and w_diff <= 4 else "⚠️ "
        print(f"{char_name}_attack: {nw}x{nh}px, {status} h_diff={h_diff}, w_diff={w_diff}")
    else:
        print(f"{char_name}_attack: 空白圖！")

    return result


if __name__ == "__main__":
    # chiikawa：粉紅色光暈（討伐棒）
    fix_attack("chiikawa", glow_color=(255, 120, 200), glow_cx_offset=18, glow_cy_offset=-20)

    # hachiware：藍色光暈（小八的攻擊感）
    fix_attack("hachiware", glow_color=(100, 180, 255), glow_cx_offset=18, glow_cy_offset=-20)

    print("\nDone! 執行 process_sprites.py --mode qc 確認結果")
