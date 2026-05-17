# -*- coding: utf-8 -*-
"""
修復 usagi bigwin 一致性問題 v2
目標：bigwin 的 bbox 高度和 idle 一致（78px），height diff <= 2px
策略：不上移，改用輕微縮放 + 金色色調 + 星星，保持與 idle 相同的 bbox 範圍
"""
from PIL import Image, ImageEnhance
import os
import numpy as np

CHARS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"

def fix_usagi_bigwin_v2():
    # 載入 idle
    idle_path = os.path.join(CHARS_DIR, "usagi_idle.png")
    img = Image.open(idle_path).convert("RGBA")
    idle_bbox = img.getbbox()
    print(f"usagi_idle: {img.size}, bbox={idle_bbox}")
    idle_h = idle_bbox[3] - idle_bbox[1]
    idle_w = idle_bbox[2] - idle_bbox[0]
    print(f"  content size: {idle_w}x{idle_h}px")

    SIZE = 96

    # bigwin 效果：
    # 1. 輕微放大 1.02x（幾乎不變，保持 bbox 一致）
    # 2. 金色色調
    # 3. 加星星光點（不影響 bbox）
    # 4. 上移 2px（在 idle bbox 範圍內）

    scale = 1.02
    new_size = int(SIZE * scale)

    # 放大
    scaled = img.resize((new_size, new_size), Image.NEAREST)

    # 建立畫布，置中並輕微上移（2px，不超出 idle 的 y 範圍）
    canvas = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    paste_x = (SIZE - new_size) // 2
    paste_y = (SIZE - new_size) // 2 - 2  # 只上移 2px
    canvas.paste(scaled, (paste_x, paste_y))

    # 金色色調（輕微）
    arr = np.array(canvas)
    mask = arr[:, :, 3] > 10
    arr[mask, 0] = np.clip(arr[mask, 0].astype(int) + 20, 0, 255)  # 加紅
    arr[mask, 1] = np.clip(arr[mask, 1].astype(int) + 12, 0, 255)  # 加綠
    arr[mask, 2] = np.clip(arr[mask, 2].astype(int) - 8, 0, 255)   # 減藍
    canvas = Image.fromarray(arr.copy()).copy()

    # 加星星光點（在 idle bbox 範圍內，不擴大 bbox）
    pixels = canvas.load()
    # idle bbox: (14, 12, 82, 90)，星星要在 y>=12, x 在 14-82 範圍內
    # 放在角色上半部的空白區域
    star_positions = [
        (20, 13), (72, 14),   # 上方兩顆（緊貼 idle bbox 頂部）
        (18, 30), (76, 28),   # 中上方兩顆
    ]
    STAR = (255, 235, 50, 255)
    STAR_DIM = (255, 235, 50, 160)
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
    arr_check = np.array(canvas)
    non_t = int((arr_check[:,:,3] > 10).sum())
    bbox = canvas.getbbox()
    print(f"usagi_bigwin (v2): {canvas.size}, {non_t}px, bbox={bbox}")
    if bbox:
        bw = bbox[2] - bbox[0]
        bh = bbox[3] - bbox[1]
        print(f"  content size: {bw}x{bh}px")
        h_diff = abs(bh - idle_h)
        w_diff = abs(bw - idle_w)
        status = "✅" if h_diff <= 2 and w_diff <= 4 else "⚠️ "
        print(f"  {status} vs idle: height diff={h_diff}px, width diff={w_diff}px")

if __name__ == "__main__":
    fix_usagi_bigwin_v2()
    print("Done!")
