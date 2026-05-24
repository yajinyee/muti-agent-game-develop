# -*- coding: utf-8 -*-
"""
T235 幸運閃電錘魚精靈圖生成
閃電錘主題：閃電錘頭 + 電光魚身 + 閃電紋路 + 金色電光點
"""
from PIL import Image
import os
import math

OUTPUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def new_img():
    return Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

def px(img, x, y, c):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), c)

def fill_circle(img, cx, cy, r, color):
    for y in range(cy-r, cy+r+1):
        for x in range(cx-r, cx+r+1):
            if (x-cx)**2 + (y-cy)**2 <= r**2:
                px(img, x, y, color)

def fill_circle_shaded(img, cx, cy, r, base_rgb):
    r_v, g_v, b_v = base_rgb
    light = (min(255,r_v+40), min(255,g_v+40), min(255,b_v+40), 255)
    mid   = (r_v, g_v, b_v, 255)
    dark  = (max(0,r_v-45), max(0,g_v-45), max(0,b_v-45), 255)
    for y in range(cy-r, cy+r+1):
        for x in range(cx-r, cx+r+1):
            if (x-cx)**2 + (y-cy)**2 > r**2:
                continue
            nx_ = (x-cx)/max(r,1)
            ny_ = (y-cy)/max(r,1)
            dot = -(nx_*(-0.7) + ny_*(-0.7))
            if dot > 0.25:
                c = light
            elif dot < -0.1:
                c = dark
            else:
                c = mid
            px(img, x, y, c)

def generate_t235():
    img = new_img()

    # ── 魚身（電藍色橢圓，帶陰影）────────────────────────────────────────
    BLUE_BASE  = (30, 144, 255)   # 道奇藍
    BLUE_DARK  = (0, 80, 180)     # 深藍
    BLUE_LIGHT = (135, 206, 250)  # 淡天藍

    cx, cy = 32, 32
    for y in range(cy-11, cy+12):
        for x in range(cx-21, cx+22):
            if (x-cx)**2/21**2 + (y-cy)**2/11**2 <= 1.0:
                nx_ = (x-cx)/21.0
                ny_ = (y-cy)/11.0
                dot = -(nx_*(-0.7) + ny_*(-0.7))
                if dot > 0.3:
                    c = BLUE_LIGHT + (255,)
                elif dot < -0.1:
                    c = BLUE_DARK + (255,)
                else:
                    c = BLUE_BASE + (255,)
                px(img, x, y, c)

    # 魚尾（深藍）
    TAIL_COLOR = (0, 100, 200, 255)
    for y in range(cy-7, cy+8):
        tail_x = cx + 21 + abs(y - cy)
        for x in range(cx+21, min(tail_x+1, SIZE)):
            px(img, x, y, TAIL_COLOR)

    # 魚鰭（電藍）
    FIN_COLOR = (100, 180, 255, 255)
    for i in range(7):
        fin_x = cx - 7 + i * 2
        fin_y = cy - 11 - (3 - abs(i - 3))
        fill_circle(img, fin_x, fin_y, 2, FIN_COLOR)

    # ── 閃電紋路（Z 字形，金色）──────────────────────────────────────────
    LIGHTNING_GOLD = (255, 215, 0, 230)
    LIGHTNING_WHITE = (255, 255, 200, 180)

    # 主閃電（Z 字形，從左到右）
    lightning_points = [
        (cx-14, cy-5), (cx-6, cy-5),   # 上橫
        (cx-6, cy-5), (cx-14, cy+5),   # 斜線
        (cx-14, cy+5), (cx-6, cy+5),   # 下橫
    ]
    for i in range(0, len(lightning_points), 2):
        x1, y1 = lightning_points[i]
        x2, y2 = lightning_points[i+1]
        steps = max(abs(x2-x1), abs(y2-y1))
        if steps == 0:
            continue
        for s in range(steps+1):
            lx = int(x1 + (x2-x1) * s / steps)
            ly = int(y1 + (y2-y1) * s / steps)
            fill_circle(img, lx, ly, 1, LIGHTNING_GOLD)

    # 第二條閃電（右側，較小）
    lightning2 = [
        (cx+2, cy-4), (cx+8, cy-4),
        (cx+8, cy-4), (cx+2, cy+4),
        (cx+2, cy+4), (cx+8, cy+4),
    ]
    for i in range(0, len(lightning2), 2):
        x1, y1 = lightning2[i]
        x2, y2 = lightning2[i+1]
        steps = max(abs(x2-x1), abs(y2-y1))
        if steps == 0:
            continue
        for s in range(steps+1):
            lx = int(x1 + (x2-x1) * s / steps)
            ly = int(y1 + (y2-y1) * s / steps)
            fill_circle(img, lx, ly, 1, LIGHTNING_WHITE)

    # ── 錘頭（金色矩形，左側）────────────────────────────────────────────
    HAMMER_GOLD = (255, 200, 0)
    HAMMER_DARK = (180, 130, 0)

    # 錘頭（左側，垂直矩形）
    for y in range(cy-8, cy+9):
        for x in range(cx-22, cx-14):
            if 0 <= x < SIZE and 0 <= y < SIZE:
                # 陰影
                if x < cx-20:
                    c = HAMMER_DARK + (255,)
                else:
                    c = HAMMER_GOLD + (255,)
                px(img, x, y, c)

    # 錘頭高光
    for y in range(cy-7, cy-3):
        px(img, cx-21, y, (255, 240, 100, 200))

    # ── 電光點（8個，圍繞魚身）────────────────────────────────────────────
    ELEC_GOLD  = (255, 215, 0, 255)
    ELEC_WHITE = (200, 230, 255, 200)

    for i in range(8):
        angle = math.pi * i / 4
        r_elec = 19 + (i % 3)
        ex = int(cx + r_elec * math.cos(angle))
        ey = int(cy + r_elec * math.sin(angle))
        fill_circle(img, ex, ey, 1, ELEC_GOLD)
        # 外圈電光
        ex2 = int(cx + (r_elec + 4) * math.cos(angle + math.pi/8))
        ey2 = int(cy + (r_elec + 4) * math.sin(angle + math.pi/8))
        fill_circle(img, ex2, ey2, 1, ELEC_WHITE)

    # ── 眼睛（金色，帶高光）──────────────────────────────────────────────
    EYE_GOLD = (255, 215, 0)
    fill_circle(img, cx-9, cy-2, 3, EYE_GOLD + (255,))
    fill_circle(img, cx-9, cy-2, 2, (0, 50, 150, 255))
    px(img, cx-10, cy-3, (255, 255, 200, 220))

    # ── 輪廓（深藍）──────────────────────────────────────────────────────
    OUTLINE = (0, 50, 120, 255)
    for y in range(cy-12, cy+13):
        for x in range(cx-22, cx+23):
            if (x-cx)**2/21**2 + (y-cy)**2/11**2 > 1.0:
                for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
                    nx2, ny2 = x+dx, y+dy
                    if (nx2-cx)**2/21**2 + (ny2-cy)**2/11**2 <= 1.0:
                        px(img, x, y, OUTLINE)
                        break

    # ── 電光光暈（外圈，半透明藍白）──────────────────────────────────────
    HALO_COLOR = (135, 206, 250, 45)
    for y in range(cy-17, cy+18):
        for x in range(cx-25, cx+26):
            d = math.sqrt((x-cx)**2/25**2 + (y-cy)**2/17**2)
            if 0.85 <= d <= 1.0:
                px(img, x, y, HALO_COLOR)

    return img

if __name__ == "__main__":
    img = generate_t235()
    out_path = os.path.join(OUTPUT_DIR, "T235_lightning_hammer.png")
    img.save(out_path)
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    print(f"T235 精靈圖已生成：{out_path}")
    print(f"非透明像素：{non_transparent}/{SIZE*SIZE} ({non_transparent*100//(SIZE*SIZE)}%)")
