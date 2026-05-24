# -*- coding: utf-8 -*-
"""
T234 幸運黃金颶風魚精靈圖生成
黃金颶風主題：螺旋颶風 + 金色魚身 + 旋風紋路 + 金色光點
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

def generate_t234():
    img = new_img()

    # ── 魚身（金色橢圓，帶陰影）──────────────────────────────────────────
    GOLD_BASE  = (220, 170, 0)    # 金色
    GOLD_DARK  = (150, 110, 0)    # 深金
    GOLD_LIGHT = (255, 230, 80)   # 淡金

    cx, cy = 32, 32
    for y in range(cy-11, cy+12):
        for x in range(cx-21, cx+22):
            if (x-cx)**2/21**2 + (y-cy)**2/11**2 <= 1.0:
                nx_ = (x-cx)/21.0
                ny_ = (y-cy)/11.0
                dot = -(nx_*(-0.7) + ny_*(-0.7))
                if dot > 0.3:
                    c = GOLD_LIGHT + (255,)
                elif dot < -0.1:
                    c = GOLD_DARK + (255,)
                else:
                    c = GOLD_BASE + (255,)
                px(img, x, y, c)

    # 魚尾（橙色）
    TAIL_COLOR = (200, 130, 0, 255)
    for y in range(cy-7, cy+8):
        tail_x = cx + 21 + abs(y - cy)
        for x in range(cx+21, min(tail_x+1, SIZE)):
            px(img, x, y, TAIL_COLOR)

    # 魚鰭（深金）
    FIN_COLOR = (180, 140, 0, 255)
    for i in range(7):
        fin_x = cx - 7 + i * 2
        fin_y = cy - 11 - (3 - abs(i - 3))
        fill_circle(img, fin_x, fin_y, 2, FIN_COLOR)

    # ── 螺旋颶風紋路（3條螺旋線，橙色）──────────────────────────────────
    SPIRAL_COLOR = (255, 140, 0, 200)   # 橙色螺旋
    SPIRAL_LIGHT = (255, 200, 50, 160)  # 淡橙螺旋

    # 三條螺旋線（從中心向外）
    for spiral_idx in range(3):
        base_angle = spiral_idx * (2 * math.pi / 3)
        for t_step in range(30):
            t = t_step / 30.0
            r_spiral = 3 + t * 14
            angle = base_angle + t * math.pi * 1.5
            sx = int(cx + r_spiral * math.cos(angle))
            sy = int(cy + r_spiral * math.sin(angle))
            # 只在魚身橢圓內繪製
            if (sx-cx)**2/21**2 + (sy-cy)**2/11**2 <= 0.9:
                fill_circle(img, sx, sy, 1, SPIRAL_COLOR)
                # 高光點
                if t_step % 5 == 0:
                    fill_circle(img, sx-1, sy-1, 1, SPIRAL_LIGHT)

    # ── 颶風中心（金色旋渦）──────────────────────────────────────────────
    VORTEX_GOLD = (255, 215, 0)
    VORTEX_CORE = (255, 255, 150)

    # 中心旋渦（小圓）
    fill_circle(img, cx, cy, 4, VORTEX_GOLD + (255,))
    fill_circle(img, cx, cy, 2, VORTEX_CORE + (255,))
    px(img, cx-1, cy-1, (255, 255, 255, 220))

    # ── 颶風光點（8個，圍繞魚身，螺旋排列）──────────────────────────────
    WIND_GOLD  = (255, 215, 0, 255)
    WIND_ORG   = (255, 165, 0, 200)
    WIND_LIGHT = (255, 240, 100, 180)

    for i in range(8):
        angle = math.pi * i / 4 + math.pi / 8
        r_wind = 18 + (i % 3) * 2
        wx = int(cx + r_wind * math.cos(angle))
        wy = int(cy + r_wind * math.sin(angle))
        fill_circle(img, wx, wy, 1, WIND_GOLD)
        # 外圈光點
        wx2 = int(cx + (r_wind + 4) * math.cos(angle + math.pi/8))
        wy2 = int(cy + (r_wind + 4) * math.sin(angle + math.pi/8))
        fill_circle(img, wx2, wy2, 1, WIND_ORG)

    # 額外小光點（12個，更外圈）
    for i in range(12):
        angle = math.pi * i / 6
        r_outer = 24 + (i % 4)
        ox = int(cx + r_outer * math.cos(angle))
        oy = int(cy + r_outer * math.sin(angle))
        fill_circle(img, ox, oy, 1, WIND_LIGHT)

    # ── 眼睛（橙色，帶高光）──────────────────────────────────────────────
    EYE_ORG = (255, 140, 0)
    fill_circle(img, cx-9, cy-2, 3, EYE_ORG + (255,))
    fill_circle(img, cx-9, cy-2, 2, (100, 50, 0, 255))
    px(img, cx-10, cy-3, (255, 255, 200, 220))

    # ── 輪廓（深金）──────────────────────────────────────────────────────
    OUTLINE = (100, 70, 0, 255)
    for y in range(cy-12, cy+13):
        for x in range(cx-22, cx+23):
            if (x-cx)**2/21**2 + (y-cy)**2/11**2 > 1.0:
                for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
                    nx2, ny2 = x+dx, y+dy
                    if (nx2-cx)**2/21**2 + (ny2-cy)**2/11**2 <= 1.0:
                        px(img, x, y, OUTLINE)
                        break

    # ── 金色光暈（外圈，半透明）──────────────────────────────────────────
    HALO_COLOR = (255, 200, 0, 50)
    for y in range(cy-17, cy+18):
        for x in range(cx-25, cx+26):
            d = math.sqrt((x-cx)**2/25**2 + (y-cy)**2/17**2)
            if 0.85 <= d <= 1.0:
                px(img, x, y, HALO_COLOR)

    # ── 颶風外圈弧線（橙色，4條弧）────────────────────────────────────────
    ARC_COLOR = (255, 165, 0, 150)
    for arc_idx in range(4):
        base_angle = arc_idx * (math.pi / 2)
        for t_step in range(15):
            t = t_step / 15.0
            angle = base_angle + t * math.pi / 2
            r_arc = 26 + arc_idx
            ax = int(cx + r_arc * math.cos(angle))
            ay = int(cy + r_arc * math.sin(angle))
            fill_circle(img, ax, ay, 1, ARC_COLOR)

    return img

if __name__ == "__main__":
    img = generate_t234()
    out_path = os.path.join(OUTPUT_DIR, "T234_golden_hurricane.png")
    img.save(out_path)
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    print(f"T234 精靈圖已生成：{out_path}")
    print(f"非透明像素：{non_transparent}/{SIZE*SIZE} ({non_transparent*100//(SIZE*SIZE)}%)")
