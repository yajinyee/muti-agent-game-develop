# -*- coding: utf-8 -*-
"""
T245 幸運宇宙脈衝魚精靈圖生成（DAY-287）
視覺設計：宇宙主題（深紫黑漸層魚身+同心圓脈衝波+宇宙光環+星點散落+金色脈衝核心）
"""
from PIL import Image, ImageDraw
import math
import os

OUTPUT_PATH = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets\T245_cosmic_pulse.png"
SIZE = 64

def new_img():
    return Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

def px(img, x, y, c):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), c)

def fill_circle_shaded(img, cx, cy, r, base_rgb, light_dir=(-0.7, -0.7)):
    """帶陰影的圓形：左上亮，右下暗"""
    r_v, g_v, b_v = base_rgb
    light = (min(255, r_v + 40), min(255, g_v + 40), min(255, b_v + 40), 255)
    mid   = (r_v, g_v, b_v, 255)
    dark  = (max(0, r_v - 50), max(0, g_v - 50), max(0, b_v - 50), 255)
    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            if (x - cx) ** 2 + (y - cy) ** 2 > r ** 2:
                continue
            nx_ = (x - cx) / max(r, 1)
            ny_ = (y - cy) / max(r, 1)
            dot = -(nx_ * light_dir[0] + ny_ * light_dir[1])
            if dot > 0.3:
                c = light
            elif dot < -0.1:
                c = dark
            else:
                c = mid
            px(img, x, y, c)

def draw_ring(img, cx, cy, r_inner, r_outer, color):
    """畫環形"""
    for y in range(cy - r_outer - 1, cy + r_outer + 2):
        for x in range(cx - r_outer - 1, cx + r_outer + 2):
            d2 = (x - cx) ** 2 + (y - cy) ** 2
            if r_inner ** 2 <= d2 <= r_outer ** 2:
                px(img, x, y, color)

def draw_star_point(img, cx, cy, size, color):
    """畫星形光點"""
    for i in range(-size, size + 1):
        px(img, cx + i, cy, color)
        px(img, cx, cy + i, color)
    for i in range(-size // 2, size // 2 + 1):
        px(img, cx + i, cy + i, color)
        px(img, cx + i, cy - i, color)

def gen_T245():
    img = new_img()
    cx, cy = 32, 32

    # ── 1. 深紫黑漸層魚身（橢圓，橫向）──────────────────────────────────────
    # 外層：深紫黑
    for y in range(cy - 18, cy + 19):
        for x in range(cx - 26, cx + 27):
            # 橢圓方程
            if (x - cx) ** 2 / 26 ** 2 + (y - cy) ** 2 / 18 ** 2 <= 1.0:
                # 漸層：中心偏亮（深紫），邊緣偏暗（深黑）
                dist = math.sqrt((x - cx) ** 2 / 26 ** 2 + (y - cy) ** 2 / 18 ** 2)
                t = dist  # 0=中心, 1=邊緣
                r_v = int(30 + (1 - t) * 40)   # 30-70
                g_v = int(0 + (1 - t) * 10)    # 0-10
                b_v = int(60 + (1 - t) * 80)   # 60-140
                px(img, x, y, (r_v, g_v, b_v, 255))

    # ── 2. 同心圓脈衝波（3層，顏色由內到外：淡紫→紫→金）──────────────────
    # 第1層脈衝波（最內，淡紫）
    draw_ring(img, cx, cy, 8, 10, (200, 150, 255, 200))
    # 第2層脈衝波（中，紫）
    draw_ring(img, cx, cy, 14, 16, (150, 80, 220, 200))
    # 第3層脈衝波（最外，金色）
    draw_ring(img, cx, cy, 20, 22, (255, 215, 0, 180))

    # ── 3. 宇宙光環（外圈，淡紫白）──────────────────────────────────────────
    draw_ring(img, cx, cy, 24, 26, (220, 180, 255, 120))

    # ── 4. 金色脈衝核心（中心圓）────────────────────────────────────────────
    fill_circle_shaded(img, cx, cy, 6, (255, 200, 50))
    # 核心高光
    px(img, cx - 2, cy - 2, (255, 255, 200, 255))
    px(img, cx - 1, cy - 2, (255, 255, 200, 255))

    # ── 5. 4方向脈衝光芒（金色）─────────────────────────────────────────────
    # 上
    for i in range(1, 8):
        alpha = max(0, 255 - i * 30)
        px(img, cx, cy - 8 - i, (255, 215, 0, alpha))
    # 下
    for i in range(1, 8):
        alpha = max(0, 255 - i * 30)
        px(img, cx, cy + 8 + i, (255, 215, 0, alpha))
    # 左
    for i in range(1, 8):
        alpha = max(0, 255 - i * 30)
        px(img, cx - 8 - i, cy, (255, 215, 0, alpha))
    # 右
    for i in range(1, 8):
        alpha = max(0, 255 - i * 30)
        px(img, cx + 8 + i, cy, (255, 215, 0, alpha))

    # ── 6. 斜向脈衝光芒（淡紫）──────────────────────────────────────────────
    for i in range(1, 6):
        alpha = max(0, 200 - i * 35)
        c = (180, 120, 255, alpha)
        px(img, cx - 6 - i, cy - 6 - i, c)
        px(img, cx + 6 + i, cy - 6 - i, c)
        px(img, cx - 6 - i, cy + 6 + i, c)
        px(img, cx + 6 + i, cy + 6 + i, c)

    # ── 7. 星點散落（宇宙感）────────────────────────────────────────────────
    star_positions = [
        (10, 8), (52, 12), (8, 48), (54, 50),
        (18, 18), (46, 16), (16, 46), (48, 46),
        (6, 30), (58, 30), (30, 6), (30, 58),
    ]
    for sx, sy in star_positions:
        # 只在魚身外圍放星點
        dist = math.sqrt((sx - cx) ** 2 + (sy - cy) ** 2)
        if dist > 22:
            alpha = int(180 * (1 - (dist - 22) / 20))
            alpha = max(0, min(255, alpha))
            px(img, sx, sy, (200, 160, 255, alpha))
            if alpha > 100:
                px(img, sx + 1, sy, (255, 255, 255, alpha // 2))

    # ── 8. 魚眼（金色，右側）────────────────────────────────────────────────
    px(img, cx + 14, cy - 3, (255, 220, 50, 255))
    px(img, cx + 15, cy - 3, (255, 220, 50, 255))
    px(img, cx + 14, cy - 2, (255, 220, 50, 255))
    # 眼睛高光
    px(img, cx + 14, cy - 4, (255, 255, 200, 200))

    # ── 9. 魚尾（左側，紫色漸層）────────────────────────────────────────────
    for i in range(5):
        alpha = 255 - i * 40
        c = (120, 60, 200, alpha)
        px(img, cx - 27 - i, cy - 4 + i, c)
        px(img, cx - 27 - i, cy - 3 + i, c)
        px(img, cx - 27 - i, cy - 2 + i, c)
        px(img, cx - 27 - i, cy + 3 - i, c)
        px(img, cx - 27 - i, cy + 4 - i, c)

    # ── 10. 輪廓（深紫黑）────────────────────────────────────────────────────
    OUTLINE = (20, 0, 40, 255)
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 0:
                for dx, dy in [(-1, 0), (1, 0), (0, -1), (0, 1)]:
                    nx_, ny_ = x + dx, y + dy
                    if 0 <= nx_ < SIZE and 0 <= ny_ < SIZE:
                        if img.getpixel((nx_, ny_))[3] == 0:
                            px(img, nx_, ny_, OUTLINE)

    return img

if __name__ == "__main__":
    img = gen_T245()
    os.makedirs(os.path.dirname(OUTPUT_PATH), exist_ok=True)
    img.save(OUTPUT_PATH)
    # 統計非透明像素
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    total = SIZE * SIZE
    print(f"T245 宇宙脈衝魚精靈圖生成完成！")
    print(f"輸出路徑：{OUTPUT_PATH}")
    print(f"非透明像素：{non_transparent}/{total}（{non_transparent/total*100:.1f}%）")
