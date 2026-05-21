# -*- coding: utf-8 -*-
"""
修復 usagi bigwin 一致性問題 v3
目標：height diff=0px, width diff=0px
策略：完全不縮放不位移，只加金色色調 + 星星，確保 bbox 完全一致
"""
from PIL import Image
import os
import numpy as np

CHARS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"

def fix_usagi_bigwin_v3():
    idle_path = os.path.join(CHARS_DIR, "usagi_idle.png")
    img = Image.open(idle_path).convert("RGBA")
    idle_bbox = img.getbbox()
    idle_h = idle_bbox[3] - idle_bbox[1]
    idle_w = idle_bbox[2] - idle_bbox[0]
    print(f"usagi_idle: bbox={idle_bbox}, content={idle_w}x{idle_h}px")

    SIZE = 96

    # 直接從 idle 複製，不縮放不位移
    canvas = img.copy()
    arr = np.array(canvas)

    # 金色色調（只對非透明像素）
    mask = arr[:, :, 3] > 10
    arr[mask, 0] = np.clip(arr[mask, 0].astype(int) + 18, 0, 255)  # 加紅
    arr[mask, 1] = np.clip(arr[mask, 1].astype(int) + 10, 0, 255)  # 加綠
    arr[mask, 2] = np.clip(arr[mask, 2].astype(int) - 10, 0, 255)  # 減藍

    canvas = Image.fromarray(arr.copy()).copy()
    pixels = canvas.load()

    # 加星星光點（嚴格在 idle bbox 內部的空白區域）
    # idle bbox: (14, 12, 82, 90)
    # 找 idle 內部的透明區域放星星
    idle_arr = np.array(img)
    star_candidates = []
    for sy in range(idle_bbox[1], idle_bbox[1] + 20):  # 只在頂部 20px 找
        for sx in range(idle_bbox[0], idle_bbox[2]):
            if idle_arr[sy, sx, 3] < 10:  # 透明區域
                star_candidates.append((sx, sy))

    # 選幾個分散的位置放星星
    import random
    random.seed(42)
    if len(star_candidates) >= 4:
        chosen = random.sample(star_candidates, min(4, len(star_candidates)))
        STAR = (255, 235, 50, 255)
        STAR_DIM = (255, 235, 50, 150)
        for sx, sy in chosen:
            pixels[sx, sy] = STAR
            for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
                nx, ny = sx+dx, sy+dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    # 只在 idle bbox 內放星星光暈
                    if idle_bbox[0] <= nx < idle_bbox[2] and idle_bbox[1] <= ny < idle_bbox[3]:
                        if idle_arr[ny, nx, 3] < 10:  # 仍是透明區域
                            pixels[nx, ny] = STAR_DIM

    # 儲存
    out_path = os.path.join(CHARS_DIR, "usagi_bigwin.png")
    canvas.save(out_path)

    # 確認
    arr_check = np.array(canvas)
    non_t = int((arr_check[:,:,3] > 10).sum())
    bbox = canvas.getbbox()
    print(f"usagi_bigwin (v3): {non_t}px, bbox={bbox}")
    if bbox:
        bw = bbox[2] - bbox[0]
        bh = bbox[3] - bbox[1]
        h_diff = abs(bh - idle_h)
        w_diff = abs(bw - idle_w)
        status = "✅" if h_diff == 0 and w_diff == 0 else "⚠️ "
        print(f"  {status} vs idle: height diff={h_diff}px, width diff={w_diff}px")

if __name__ == "__main__":
    fix_usagi_bigwin_v3()
    print("Done!")
