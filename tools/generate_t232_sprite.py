# -*- coding: utf-8 -*-
"""
T232 幸運命運預言魚精靈圖生成
紫金命運主題：水晶球 + 星座符文 + 紫色光暈 + 金色預言光點
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

def outline_circle(img, cx, cy, r, color):
    for y in range(cy-r-1, cy+r+2):
        for x in range(cx-r-1, cx+r+2):
            d = math.sqrt((x-cx)**2 + (y-cy)**2)
            if r+0.1 <= d <= r+1.5:
                px(img, x, y, color)

def draw_star(img, cx, cy, r_outer, r_inner, points, color):
    """繪製星形"""
    for i in range(points * 2):
        angle = math.pi * i / points - math.pi / 2
        r = r_outer if i % 2 == 0 else r_inner
        x = int(cx + r * math.cos(angle))
        y = int(cy + r * math.sin(angle))
        fill_circle(img, x, y, 1, color)

def draw_rune(img, cx, cy, color):
    """繪製星座符文（六芒星輪廓）"""
    r = 8
    for i in range(6):
        angle = math.pi * i / 3
        x1 = int(cx + r * math.cos(angle))
        y1 = int(cy + r * math.sin(angle))
        x2 = int(cx + r * math.cos(angle + math.pi))
        y2 = int(cy + r * math.sin(angle + math.pi))
        # 畫線段
        steps = 12
        for s in range(steps + 1):
            t = s / steps
            lx = int(x1 + (x2 - x1) * t)
            ly = int(y1 + (y2 - y1) * t)
            px(img, lx, ly, color)

def generate_t232():
    img = new_img()

    # ── 魚身（紫色橢圓，帶陰影）──────────────────────────────────────────
    PURPLE_BASE = (155, 89, 182)   # #9B59B6 紫
    PURPLE_DARK = (100, 50, 130)   # 深紫
    PURPLE_LIGHT = (200, 150, 220) # 淡紫

    # 魚身主體（橫向橢圓）
    cx, cy = 32, 32
    for y in range(cy-12, cy+13):
        for x in range(cx-22, cx+23):
            # 橢圓方程
            if (x-cx)**2/22**2 + (y-cy)**2/12**2 <= 1.0:
                # 陰影計算
                nx_ = (x-cx)/22.0
                ny_ = (y-cy)/12.0
                dot = -(nx_*(-0.7) + ny_*(-0.7))
                if dot > 0.3:
                    c = PURPLE_LIGHT + (255,)
                elif dot < -0.1:
                    c = PURPLE_DARK + (255,)
                else:
                    c = PURPLE_BASE + (255,)
                px(img, x, y, c)

    # 魚尾（右側三角形）
    TAIL_COLOR = (130, 70, 160, 255)
    for y in range(cy-8, cy+9):
        tail_x = cx + 22 + abs(y - cy)
        for x in range(cx+22, min(tail_x+1, SIZE)):
            px(img, x, y, TAIL_COLOR)

    # 魚鰭（上方）
    FIN_COLOR = (180, 120, 200, 255)
    for i in range(8):
        fin_x = cx - 8 + i * 2
        fin_y = cy - 12 - (4 - abs(i - 4))
        fill_circle(img, fin_x, fin_y, 2, FIN_COLOR)

    # ── 水晶球（中心，帶光澤）────────────────────────────────────────────
    CRYSTAL_BASE = (180, 140, 220)  # 淡紫水晶
    CRYSTAL_GLOW = (220, 200, 255)  # 水晶高光
    CRYSTAL_DARK = (120, 80, 160)   # 水晶暗部

    fill_circle_shaded(img, cx, cy, 9, CRYSTAL_BASE)

    # 水晶球高光（左上角小圓）
    fill_circle(img, cx-3, cy-3, 2, CRYSTAL_GLOW + (200,))
    fill_circle(img, cx-2, cy-2, 1, (255, 255, 255, 220))

    # 水晶球輪廓
    outline_circle(img, cx, cy, 9, (80, 40, 120, 200))

    # ── 星座符文（水晶球內，金色）────────────────────────────────────────
    RUNE_COLOR = (255, 215, 0, 180)  # #FFD700 金色半透明
    draw_rune(img, cx, cy, RUNE_COLOR)

    # ── 預言光點（圍繞水晶球，6個金色光點）──────────────────────────────
    GLOW_COLOR = (255, 215, 0, 255)   # 金色
    GLOW_SMALL = (255, 180, 50, 200)  # 橙金
    for i in range(6):
        angle = math.pi * i / 3
        gx = int(cx + 14 * math.cos(angle))
        gy = int(cy + 14 * math.sin(angle))
        fill_circle(img, gx, gy, 2, GLOW_COLOR)
        # 小光點
        gx2 = int(cx + 18 * math.cos(angle + math.pi/6))
        gy2 = int(cy + 18 * math.sin(angle + math.pi/6))
        fill_circle(img, gx2, gy2, 1, GLOW_SMALL)

    # ── 星形光芒（4方向，金色）───────────────────────────────────────────
    STAR_COLOR = (255, 215, 0, 255)
    # 上下左右各一個小星
    for angle in [0, math.pi/2, math.pi, 3*math.pi/2]:
        sx = int(cx + 20 * math.cos(angle))
        sy = int(cy + 20 * math.sin(angle))
        draw_star(img, sx, sy, 3, 1, 4, STAR_COLOR)

    # ── 眼睛（金色，帶高光）──────────────────────────────────────────────
    EYE_COLOR = (255, 215, 0, 255)
    EYE_PUPIL = (80, 40, 120, 255)
    EYE_HIGHLIGHT = (255, 255, 255, 220)
    fill_circle(img, cx-10, cy-2, 3, EYE_COLOR)
    fill_circle(img, cx-10, cy-2, 2, EYE_PUPIL)
    px(img, cx-11, cy-3, EYE_HIGHLIGHT)

    # ── 輪廓（深紫）──────────────────────────────────────────────────────
    OUTLINE = (60, 20, 90, 255)
    # 魚身輪廓
    for y in range(cy-13, cy+14):
        for x in range(cx-23, cx+24):
            if (x-cx)**2/22**2 + (y-cy)**2/12**2 > 1.0:
                # 檢查鄰居是否在橢圓內
                for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
                    nx2, ny2 = x+dx, y+dy
                    if (nx2-cx)**2/22**2 + (ny2-cy)**2/12**2 <= 1.0:
                        px(img, x, y, OUTLINE)
                        break

    # ── 紫色光暈（外圈，半透明）──────────────────────────────────────────
    HALO_COLOR = (155, 89, 182, 60)
    for y in range(cy-18, cy+19):
        for x in range(cx-26, cx+27):
            d = math.sqrt((x-cx)**2/26**2 + (y-cy)**2/18**2)
            if 0.85 <= d <= 1.0:
                px(img, x, y, HALO_COLOR)

    return img

if __name__ == "__main__":
    img = generate_t232()
    out_path = os.path.join(OUTPUT_DIR, "T232_fortune_prophecy.png")
    img.save(out_path)
    # 統計非透明像素
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    print(f"T232 精靈圖已生成：{out_path}")
    print(f"非透明像素：{non_transparent}/{SIZE*SIZE} ({non_transparent*100//(SIZE*SIZE)}%)")
