# -*- coding: utf-8 -*-
"""
修復 usagi bigwin 一致性問題
從 usagi_idle（AI 生成，5304px）做變換生成 bigwin
確保大小與 idle 一致
"""
from PIL import Image, ImageEnhance
import os
import numpy as np

CHARS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"

def fix_usagi_bigwin():
    # 載入 idle（AI 生成，品質好）
    idle_path = os.path.join(CHARS_DIR, "usagi_idle.png")
    img = Image.open(idle_path).convert("RGBA")
    print(f"usagi_idle: {img.size}, bbox={img.getbbox()}")

    # bigwin 效果：
    # 1. 整體上移 8px（跳起感）
    # 2. 輕微放大 1.05x
    # 3. 金色色調
    # 4. 加星星光點

    SIZE = 96
    scale = 1.05
    new_size = int(SIZE * scale)

    # 放大
    scaled = img.resize((new_size, new_size), Image.NEAREST)

    # 建立畫布，置中並上移
    canvas = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    paste_x = (SIZE - new_size) // 2
    paste_y = (SIZE - new_size) // 2 - 8  # 上移 8px
    canvas.paste(scaled, (paste_x, paste_y))

    # 金色色調（輕微）
    arr = np.array(canvas)
    # 對非透明像素加金色調
    mask = arr[:, :, 3] > 10
    arr[mask, 0] = np.clip(arr[mask, 0].astype(int) + 15, 0, 255)  # 加紅
    arr[mask, 1] = np.clip(arr[mask, 1].astype(int) + 10, 0, 255)  # 加綠
    arr[mask, 2] = np.clip(arr[mask, 2].astype(int) - 5, 0, 255)   # 減藍
    canvas = Image.fromarray(arr).copy()  # copy() 確保可寫

    # 加星星光點（4個角落）
    pixels = canvas.load()
    star_positions = [(12, 12), (SIZE-14, 10), (10, SIZE//3), (SIZE-12, SIZE//3)]
    STAR = (255, 235, 50, 255)
    STAR_DIM = (255, 235, 50, 180)
    for sx, sy in star_positions:
        if 0 <= sx < SIZE and 0 <= sy < SIZE:
            pixels[sx, sy] = STAR
            for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
                nx, ny = sx+dx, sy+dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[nx, ny] = STAR_DIM

    # 儲存
    out_path = os.path.join(CHARS_DIR, "usagi_bigwin.png")
    canvas.save(out_path)

    # 確認品質
    non_t = sum(1 for px in canvas.getdata() if px[3] > 10)
    bbox = canvas.getbbox()
    print(f"usagi_bigwin (fixed): {canvas.size}, {non_t}px, bbox={bbox}")
    if bbox:
        print(f"  size: {bbox[2]-bbox[0]}x{bbox[3]-bbox[1]}px")

if __name__ == "__main__":
    fix_usagi_bigwin()
    print("Done!")
