# -*- coding: utf-8 -*-
"""
T236 幸運時間裂縫魚精靈圖生成
時間裂縫主題：紫色裂縫紋路 + 沙漏符文 + 紫色魚身 + 時間光點 + 裂縫光暈
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

def generate_t236():
    img = new_img()

    cx, cy = 32, 32

    # ── 魚身（深紫色橢圓，帶陰影）────────────────────────────────────────
    PURPLE_BASE  = (120, 60, 180)   # 深紫
    PURPLE_DARK  = (70, 30, 120)    # 更深紫
    PURPLE_LIGHT = (180, 120, 240)  # 淡紫

    for y in range(cy-11, cy+12):
        for x in range(cx-21, cx+22):
            if (x-cx)**2/21**2 + (y-cy)**2/11**2 <= 1.0:
                nx_ = (x-cx)/21.0
                ny_ = (y-cy)/11.0
                dot = -(nx_*(-0.7) + ny_*(-0.7))
                if dot > 0.3:
                    c = PURPLE_LIGHT + (255,)
                elif dot < -0.1:
                    c = PURPLE_DARK + (255,)
                else:
                    c = PURPLE_BASE + (255,)
                px(img, x, y, c)

    # 魚尾（深紫）
    TAIL_COLOR = (80, 30, 140, 255)
    for y in range(cy-7, cy+8):
        tail_x = cx + 21 + abs(y - cy)
        for x in range(cx+21, min(tail_x+1, SIZE)):
            px(img, x, y, TAIL_COLOR)

    # 魚鰭（淡紫）
    FIN_COLOR = (160, 100, 220, 255)
    for i in range(7):
        fin_x = cx - 7 + i * 2
        fin_y = cy - 11 - (3 - abs(i - 3))
        fill_circle(img, fin_x, fin_y, 2, FIN_COLOR)

    # ── 裂縫紋路（鋸齒形，金色）──────────────────────────────────────────
    RIFT_GOLD  = (255, 215, 0, 220)
    RIFT_WHITE = (240, 220, 255, 180)

    # 主裂縫（從左到右，鋸齒形）
    rift_points = [
        (cx-16, cy-3), (cx-10, cy+3),
        (cx-10, cy+3), (cx-4, cy-3),
        (cx-4, cy-3), (cx+2, cy+3),
        (cx+2, cy+3), (cx+8, cy-3),
    ]
    for i in range(0, len(rift_points), 2):
        x1, y1 = rift_points[i]
        x2, y2 = rift_points[i+1]
        steps = max(abs(x2-x1), abs(y2-y1))
        if steps == 0:
            continue
        for s in range(steps+1):
            lx = int(x1 + (x2-x1) * s / steps)
            ly = int(y1 + (y2-y1) * s / steps)
            fill_circle(img, lx, ly, 1, RIFT_GOLD)

    # 裂縫光芒（從裂縫中心向外輻射）
    rift_cx, rift_cy = cx - 4, cy
    for angle_deg in range(0, 360, 45):
        angle = math.radians(angle_deg)
        for r in range(3, 8):
            rx = int(rift_cx + r * math.cos(angle))
            ry = int(rift_cy + r * math.sin(angle))
            if r < 5:
                fill_circle(img, rx, ry, 1, RIFT_GOLD)
            else:
                fill_circle(img, rx, ry, 1, RIFT_WHITE)

    # ── 沙漏符文（右側，金色）────────────────────────────────────────────
    HOURGLASS_GOLD = (255, 200, 0, 255)
    HOURGLASS_DARK = (180, 130, 0, 200)

    # 沙漏上半部（三角形）
    hx, hy = cx + 10, cy
    for row in range(5):
        width = 5 - row
        for col in range(-width, width+1):
            px(img, hx + col, hy - 5 + row, HOURGLASS_GOLD)

    # 沙漏下半部（倒三角形）
    for row in range(5):
        width = row
        for col in range(-width, width+1):
            px(img, hx + col, hy + 1 + row, HOURGLASS_GOLD)

    # 沙漏中間（細腰）
    px(img, hx, hy, HOURGLASS_DARK)

    # 沙漏外框
    for row in range(11):
        if row == 0 or row == 10:
            for col in range(-5, 6):
                px(img, hx + col, hy - 5 + row, HOURGLASS_DARK)

    # ── 時間光點（8個，圍繞魚身）────────────────────────────────────────
    TIME_PURPLE = (200, 100, 255, 255)
    TIME_GOLD   = (255, 215, 0, 200)

    for i in range(8):
        angle = math.pi * i / 4 + math.pi / 8
        r_pt = 18 + (i % 3)
        tx = int(cx + r_pt * math.cos(angle))
        ty = int(cy + r_pt * math.sin(angle))
        color = TIME_GOLD if i % 2 == 0 else TIME_PURPLE
        fill_circle(img, tx, ty, 1, color)

    # ── 眼睛（金色，帶高光）──────────────────────────────────────────────
    EYE_GOLD = (255, 215, 0)
    fill_circle(img, cx-9, cy-2, 3, EYE_GOLD + (255,))
    fill_circle(img, cx-9, cy-2, 2, (60, 0, 100, 255))
    px(img, cx-10, cy-3, (255, 255, 200, 220))

    # ── 輪廓（深紫）──────────────────────────────────────────────────────
    OUTLINE = (50, 0, 80, 255)
    for y in range(cy-12, cy+13):
        for x in range(cx-22, cx+23):
            if (x-cx)**2/21**2 + (y-cy)**2/11**2 > 1.0:
                for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
                    nx2, ny2 = x+dx, y+dy
                    if (nx2-cx)**2/21**2 + (ny2-cy)**2/11**2 <= 1.0:
                        px(img, x, y, OUTLINE)
                        break

    # ── 裂縫光暈（外圈，半透明紫色）──────────────────────────────────────
    HALO_COLOR = (180, 100, 255, 40)
    for y in range(cy-17, cy+18):
        for x in range(cx-25, cx+26):
            d = math.sqrt((x-cx)**2/25**2 + (y-cy)**2/17**2)
            if 0.85 <= d <= 1.0:
                px(img, x, y, HALO_COLOR)

    # ── 裂縫中心光點（亮白紫）────────────────────────────────────────────
    fill_circle(img, rift_cx, rift_cy, 2, (255, 240, 255, 255))
    fill_circle(img, rift_cx, rift_cy, 1, (255, 255, 255, 255))

    return img

if __name__ == "__main__":
    img = generate_t236()
    out_path = os.path.join(OUTPUT_DIR, "T236_time_rift.png")
    img.save(out_path)
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    print(f"T236 精靈圖已生成：{out_path}")
    print(f"非透明像素：{non_transparent}/{SIZE*SIZE} ({non_transparent*100//(SIZE*SIZE)}%)")
