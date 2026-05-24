# -*- coding: utf-8 -*-
"""
T233 幸運幸運圖騰魚精靈圖生成
翠綠幸運主題：四葉草 + 圖騰柱 + 翠綠光暈 + 金色幸運光點
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

def draw_clover_leaf(img, cx, cy, angle, r, color):
    """繪製一片四葉草葉子（圓形）"""
    lx = int(cx + r * math.cos(angle))
    ly = int(cy + r * math.sin(angle))
    fill_circle(img, lx, ly, r-1, color)

def generate_t233():
    img = new_img()

    # ── 魚身（翠綠橢圓，帶陰影）──────────────────────────────────────────
    GREEN_BASE = (0, 200, 100)    # 翠綠
    GREEN_DARK = (0, 130, 60)     # 深翠綠
    GREEN_LIGHT = (100, 255, 160) # 淡翠綠

    cx, cy = 32, 32
    for y in range(cy-11, cy+12):
        for x in range(cx-21, cx+22):
            if (x-cx)**2/21**2 + (y-cy)**2/11**2 <= 1.0:
                nx_ = (x-cx)/21.0
                ny_ = (y-cy)/11.0
                dot = -(nx_*(-0.7) + ny_*(-0.7))
                if dot > 0.3:
                    c = GREEN_LIGHT + (255,)
                elif dot < -0.1:
                    c = GREEN_DARK + (255,)
                else:
                    c = GREEN_BASE + (255,)
                px(img, x, y, c)

    # 魚尾
    TAIL_COLOR = (0, 160, 80, 255)
    for y in range(cy-7, cy+8):
        tail_x = cx + 21 + abs(y - cy)
        for x in range(cx+21, min(tail_x+1, SIZE)):
            px(img, x, y, TAIL_COLOR)

    # 魚鰭
    FIN_COLOR = (50, 220, 120, 255)
    for i in range(7):
        fin_x = cx - 7 + i * 2
        fin_y = cy - 11 - (3 - abs(i - 3))
        fill_circle(img, fin_x, fin_y, 2, FIN_COLOR)

    # ── 四葉草（中心，帶光澤）────────────────────────────────────────────
    CLOVER_GREEN = (0, 220, 110)   # 四葉草綠
    CLOVER_LIGHT = (80, 255, 160)  # 四葉草高光
    CLOVER_STEM  = (0, 150, 60)    # 莖

    # 四片葉子（上下左右）
    for angle in [0, math.pi/2, math.pi, 3*math.pi/2]:
        draw_clover_leaf(img, cx, cy, angle, 5, CLOVER_GREEN)

    # 葉子高光
    for angle in [0, math.pi/2, math.pi, 3*math.pi/2]:
        lx = int(cx + 5 * math.cos(angle))
        ly = int(cy + 5 * math.sin(angle))
        fill_circle(img, lx-1, ly-1, 1, CLOVER_LIGHT + (200,))

    # 中心圓（金色）
    GOLD = (255, 215, 0)
    fill_circle(img, cx, cy, 3, GOLD + (255,))
    fill_circle(img, cx-1, cy-1, 1, (255, 255, 200, 220))

    # ── 圖騰符文（圍繞四葉草，4個金色符文點）────────────────────────────
    RUNE_COLOR = (255, 215, 0, 200)
    for i in range(4):
        angle = math.pi * i / 2 + math.pi/4
        rx = int(cx + 12 * math.cos(angle))
        ry = int(cy + 12 * math.sin(angle))
        fill_circle(img, rx, ry, 2, RUNE_COLOR)
        # 符文光點
        rx2 = int(cx + 16 * math.cos(angle))
        ry2 = int(cy + 16 * math.sin(angle))
        fill_circle(img, rx2, ry2, 1, (255, 255, 100, 180))

    # ── 幸運光點（8個，圍繞魚身）────────────────────────────────────────
    LUCKY_COLOR = (255, 215, 0, 255)
    LUCKY_SMALL = (100, 255, 160, 200)
    for i in range(8):
        angle = math.pi * i / 4
        lx = int(cx + 20 * math.cos(angle))
        ly = int(cy + 20 * math.sin(angle))
        fill_circle(img, lx, ly, 1, LUCKY_COLOR)
        # 小光點
        lx2 = int(cx + 24 * math.cos(angle + math.pi/8))
        ly2 = int(cy + 24 * math.sin(angle + math.pi/8))
        fill_circle(img, lx2, ly2, 1, LUCKY_SMALL)

    # ── 眼睛（金色，帶高光）──────────────────────────────────────────────
    fill_circle(img, cx-9, cy-2, 3, GOLD + (255,))
    fill_circle(img, cx-9, cy-2, 2, (0, 100, 50, 255))
    px(img, cx-10, cy-3, (255, 255, 200, 220))

    # ── 輪廓（深翠綠）────────────────────────────────────────────────────
    OUTLINE = (0, 80, 40, 255)
    for y in range(cy-12, cy+13):
        for x in range(cx-22, cx+23):
            if (x-cx)**2/21**2 + (y-cy)**2/11**2 > 1.0:
                for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
                    nx2, ny2 = x+dx, y+dy
                    if (nx2-cx)**2/21**2 + (ny2-cy)**2/11**2 <= 1.0:
                        px(img, x, y, OUTLINE)
                        break

    # ── 翠綠光暈（外圈，半透明）──────────────────────────────────────────
    HALO_COLOR = (0, 200, 100, 55)
    for y in range(cy-17, cy+18):
        for x in range(cx-25, cx+26):
            d = math.sqrt((x-cx)**2/25**2 + (y-cy)**2/17**2)
            if 0.85 <= d <= 1.0:
                px(img, x, y, HALO_COLOR)

    return img

if __name__ == "__main__":
    img = generate_t233()
    out_path = os.path.join(OUTPUT_DIR, "T233_luck_totem.png")
    img.save(out_path)
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    print(f"T233 精靈圖已生成：{out_path}")
    print(f"非透明像素：{non_transparent}/{SIZE*SIZE} ({non_transparent*100//(SIZE*SIZE)}%)")
