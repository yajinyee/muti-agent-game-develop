# -*- coding: utf-8 -*-
"""
T246 幸運多米諾魚精靈圖生成（DAY-288）
視覺設計：多米諾主題（深紅棕漸層魚身+多米諾骨牌紋路+連鎖光環+金色點數+橙棕魚尾）
"""
from PIL import Image
import math
import os

OUTPUT_PATH = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets\T246_domino.png"
SIZE = 64

def px(img, x, y, c):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), c)

def fill_ellipse_shaded(img, cx, cy, rx, ry, base_rgb):
    """帶陰影的橢圓：左上亮，右下暗"""
    r_v, g_v, b_v = base_rgb
    light = (min(255, r_v + 40), min(255, g_v + 30), min(255, b_v + 20), 255)
    mid   = (r_v, g_v, b_v, 255)
    dark  = (max(0, r_v - 50), max(0, g_v - 35), max(0, b_v - 20), 255)
    for y in range(cy - ry, cy + ry + 1):
        for x in range(cx - rx, cx + rx + 1):
            if (x - cx) ** 2 / rx ** 2 + (y - cy) ** 2 / ry ** 2 > 1.0:
                continue
            nx_ = (x - cx) / max(rx, 1)
            ny_ = (y - cy) / max(ry, 1)
            dot = -(nx_ * (-0.7) + ny_ * (-0.7))
            if dot > 0.3:
                c = light
            elif dot < -0.1:
                c = dark
            else:
                c = mid
            px(img, x, y, c)

def draw_domino_tile(img, x1, y1, w, h, dots, color_bg, color_dot, color_line):
    """畫多米諾骨牌格子"""
    # 背景
    for y in range(y1, y1 + h):
        for x in range(x1, x1 + w):
            px(img, x, y, color_bg)
    # 中線
    mid_y = y1 + h // 2
    for x in range(x1, x1 + w):
        px(img, x, mid_y, color_line)
    # 點數（上半）
    half_h = h // 2
    dot_positions = {
        1: [(w // 2, half_h // 2)],
        2: [(w // 4, half_h // 4), (3 * w // 4, 3 * half_h // 4)],
        3: [(w // 4, half_h // 4), (w // 2, half_h // 2), (3 * w // 4, 3 * half_h // 4)],
    }
    for dx, dy in dot_positions.get(dots, []):
        px(img, x1 + dx, y1 + dy, color_dot)
        px(img, x1 + dx + 1, y1 + dy, color_dot)
        px(img, x1 + dx, y1 + dy + 1, color_dot)

def gen_T246():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

    cx, cy = 32, 32

    # ── 1. 深紅棕漸層橢圓魚身 ────────────────────────────────────────────────
    fill_ellipse_shaded(img, cx, cy, 26, 18, (120, 40, 10))

    # ── 2. 多米諾骨牌紋路（魚身上的骨牌圖案）────────────────────────────────
    # 左半骨牌（3點）
    TILE_BG   = (180, 80, 20, 220)
    TILE_DOT  = (255, 215, 0, 255)   # 金色點數
    TILE_LINE = (80, 20, 5, 255)     # 深棕中線
    # 左骨牌
    for y in range(cy - 8, cy + 9):
        for x in range(cx - 18, cx - 2):
            if (x - cx) ** 2 / 26 ** 2 + (y - cy) ** 2 / 18 ** 2 <= 0.85:
                px(img, x, y, TILE_BG)
    # 右骨牌
    for y in range(cy - 8, cy + 9):
        for x in range(cx + 2, cx + 18):
            if (x - cx) ** 2 / 26 ** 2 + (y - cy) ** 2 / 18 ** 2 <= 0.85:
                px(img, x, y, TILE_BG)
    # 中線（骨牌分隔）
    for y in range(cy - 8, cy + 9):
        if (0) ** 2 / 26 ** 2 + (y - cy) ** 2 / 18 ** 2 <= 0.85:
            px(img, cx - 1, y, TILE_LINE)
            px(img, cx, y, TILE_LINE)
            px(img, cx + 1, y, TILE_LINE)

    # 左骨牌點數（3點）
    dot_positions_left = [(cx - 14, cy - 4), (cx - 10, cy), (cx - 6, cy + 4)]
    for dx, dy in dot_positions_left:
        px(img, dx, dy, TILE_DOT)
        px(img, dx + 1, dy, TILE_DOT)
        px(img, dx, dy + 1, TILE_DOT)
        px(img, dx + 1, dy + 1, TILE_DOT)

    # 右骨牌點數（5點）
    dot_positions_right = [
        (cx + 5, cy - 5), (cx + 13, cy - 5),
        (cx + 9, cy),
        (cx + 5, cy + 5), (cx + 13, cy + 5),
    ]
    for dx, dy in dot_positions_right:
        px(img, dx, dy, TILE_DOT)
        px(img, dx + 1, dy, TILE_DOT)
        px(img, dx, dy + 1, TILE_DOT)
        px(img, dx + 1, dy + 1, TILE_DOT)

    # ── 3. 連鎖光環（橙棕色，外圈）──────────────────────────────────────────
    RING_COLOR = (210, 105, 30, 160)
    for y in range(cy - 22, cy + 23):
        for x in range(cx - 30, cx + 31):
            d2 = (x - cx) ** 2 / 30 ** 2 + (y - cy) ** 2 / 22 ** 2
            if 0.88 <= d2 <= 1.0:
                px(img, x, y, RING_COLOR)

    # ── 4. 金色連鎖光點（4方向）──────────────────────────────────────────────
    GOLD = (255, 215, 0, 255)
    chain_pts = [
        (cx, cy - 20), (cx, cy + 20),
        (cx - 28, cy), (cx + 28, cy),
    ]
    for gx, gy in chain_pts:
        for i in range(3):
            alpha = 255 - i * 70
            px(img, gx, gy + i, (255, 215, 0, alpha))
            px(img, gx, gy - i, (255, 215, 0, alpha))

    # ── 5. 魚眼（金色，右側）────────────────────────────────────────────────
    px(img, cx + 14, cy - 3, (255, 220, 50, 255))
    px(img, cx + 15, cy - 3, (255, 220, 50, 255))
    px(img, cx + 14, cy - 2, (255, 220, 50, 255))
    px(img, cx + 14, cy - 4, (255, 255, 200, 200))  # 高光

    # ── 6. 橙棕魚尾（左側）──────────────────────────────────────────────────
    for i in range(5):
        alpha = 255 - i * 40
        c = (180, 70, 15, alpha)
        px(img, cx - 27 - i, cy - 4 + i, c)
        px(img, cx - 27 - i, cy - 3 + i, c)
        px(img, cx - 27 - i, cy - 2 + i, c)
        px(img, cx - 27 - i, cy + 3 - i, c)
        px(img, cx - 27 - i, cy + 4 - i, c)

    # ── 7. 輪廓（深棕黑）────────────────────────────────────────────────────
    OUTLINE = (30, 10, 0, 255)
    pixels_copy = [img.getpixel((x, y)) for y in range(SIZE) for x in range(SIZE)]
    for y in range(SIZE):
        for x in range(SIZE):
            if pixels_copy[y * SIZE + x][3] > 0:
                for dx, dy in [(-1, 0), (1, 0), (0, -1), (0, 1)]:
                    nx_, ny_ = x + dx, y + dy
                    if 0 <= nx_ < SIZE and 0 <= ny_ < SIZE:
                        if pixels_copy[ny_ * SIZE + nx_][3] == 0:
                            px(img, nx_, ny_, OUTLINE)

    return img

if __name__ == "__main__":
    img = gen_T246()
    os.makedirs(os.path.dirname(OUTPUT_PATH), exist_ok=True)
    img.save(OUTPUT_PATH)
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    total = SIZE * SIZE
    print(f"T246 多米諾魚精靈圖生成完成！")
    print(f"輸出路徑：{OUTPUT_PATH}")
    print(f"非透明像素：{non_transparent}/{total}（{non_transparent/total*100:.1f}%）")
