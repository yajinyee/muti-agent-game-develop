# -*- coding: utf-8 -*-
"""
T249 幸運時空裂縫魚精靈圖生成（DAY-291）
視覺設計：時空主題（#00E5FF 青藍 + #0D47A1 深藍 + #7B2FBE 紫 + #FFD700 金 + #FFFFFF 白）
特徵：橢圓深藍漸層魚身+時空裂縫紋路+青藍光環+4方向時空光芒+金色裂縫核心+紫色魚尾+時空粒子
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

def fill_ellipse_shaded(img, cx, cy, rx, ry, base_rgb):
    """帶陰影的橢圓：左上亮，右下暗"""
    r_v, g_v, b_v = base_rgb
    light = (min(255,r_v+40), min(255,g_v+40), min(255,b_v+40), 255)
    mid   = (r_v, g_v, b_v, 255)
    dark  = (max(0,r_v-50), max(0,g_v-50), max(0,b_v-50), 255)
    for y in range(cy-ry, cy+ry+1):
        for x in range(cx-rx, cx+rx+1):
            if (x-cx)**2/max(rx**2,1) + (y-cy)**2/max(ry**2,1) > 1.0:
                continue
            nx_ = (x-cx)/max(rx,1)
            ny_ = (y-cy)/max(ry,1)
            dot = -(nx_*(-0.7) + ny_*(-0.7))
            if dot > 0.25:
                c = light
            elif dot < -0.1:
                c = dark
            else:
                c = mid
            px(img, x, y, c)

def draw_line(img, x0, y0, x1, y1, color, width=1):
    """Bresenham 直線"""
    dx = abs(x1-x0); dy = abs(y1-y0)
    sx = 1 if x0 < x1 else -1
    sy = 1 if y0 < y1 else -1
    err = dx - dy
    while True:
        for dw in range(-width//2, width//2+1):
            px(img, x0+dw, y0, color)
            px(img, x0, y0+dw, color)
        if x0 == x1 and y0 == y1:
            break
        e2 = 2*err
        if e2 > -dy:
            err -= dy; x0 += sx
        if e2 < dx:
            err += dx; y0 += sy

def generate_t249():
    img = new_img()

    # 顏色定義
    DEEP_BLUE   = (13, 71, 161)    # #0D47A1 深藍（魚身基底）
    CYAN        = (0, 229, 255)    # #00E5FF 青藍（光環/光芒）
    PURPLE      = (123, 47, 190)   # #7B2FBE 紫（魚尾/裂縫）
    GOLD        = (255, 215, 0)    # #FFD700 金（核心/裂縫邊緣）
    WHITE       = (255, 255, 255)  # 白（高光/粒子）
    LIGHT_CYAN  = (178, 245, 255)  # 淡青（高光）
    DARK_BLUE   = (5, 30, 80)      # 深藍黑（輪廓）
    LIGHT_PURPLE= (180, 130, 230)  # 淡紫（裂縫光）

    cx, cy = 32, 32

    # ── 1. 時空光環（最底層，大橢圓光環）──
    for r in range(28, 31):
        for angle in range(0, 360, 2):
            rad = math.radians(angle)
            ex = int(cx + r * 1.3 * math.cos(rad))
            ey = int(cy + r * 0.85 * math.sin(rad))
            alpha = 80 + int(60 * abs(math.sin(math.radians(angle * 2))))
            px(img, ex, ey, (*CYAN, alpha))

    # ── 2. 魚身（橢圓，深藍漸層）──
    fill_ellipse_shaded(img, cx, cy, 22, 16, DEEP_BLUE)

    # ── 3. 時空裂縫紋路（Z字形裂縫，貫穿魚身）──
    # 主裂縫（從左到右，Z字形）
    crack_points = [
        (12, 28), (18, 24), (22, 32), (28, 26), (34, 34), (40, 28), (46, 32)
    ]
    for i in range(len(crack_points)-1):
        x0, y0 = crack_points[i]
        x1, y1 = crack_points[i+1]
        draw_line(img, x0, y0, x1, y1, (*CYAN, 200), 1)
        draw_line(img, x0+1, y0, x1+1, y1, (*LIGHT_PURPLE, 120), 1)

    # 副裂縫（較短，斜向）
    sub_cracks = [
        [(20, 20), (24, 28)],
        [(36, 22), (32, 30)],
        [(28, 36), (34, 42)],
    ]
    for crack in sub_cracks:
        x0, y0 = crack[0]
        x1, y1 = crack[1]
        draw_line(img, x0, y0, x1, y1, (*PURPLE, 160), 1)

    # ── 4. 4方向時空光芒（十字形）──
    # 上下左右各一條光芒
    for length, angle_deg in [(14, 0), (14, 90), (10, 45), (10, 135)]:
        rad = math.radians(angle_deg)
        for d in range(1, length+1):
            ex = int(cx + d * math.cos(rad))
            ey = int(cy - d * math.sin(rad))
            alpha = max(20, 200 - d * 14)
            px(img, ex, ey, (*CYAN, alpha))
            # 對稱方向
            ex2 = int(cx - d * math.cos(rad))
            ey2 = int(cy + d * math.sin(rad))
            px(img, ex2, ey2, (*CYAN, alpha))

    # ── 5. 金色裂縫核心（中心圓）──
    fill_circle(img, cx, cy, 5, (*GOLD, 255))
    fill_circle(img, cx, cy, 3, (*WHITE, 255))
    # 核心光暈
    for r in range(6, 9):
        for angle in range(0, 360, 15):
            rad = math.radians(angle)
            ex = int(cx + r * math.cos(rad))
            ey = int(cy + r * math.sin(rad))
            px(img, ex, ey, (*GOLD, 120))

    # ── 6. 時空之眼（青藍眼睛）──
    eye_x, eye_y = 38, 27
    fill_circle(img, eye_x, eye_y, 4, (*WHITE, 255))
    fill_circle(img, eye_x, eye_y, 3, (*CYAN, 255))
    fill_circle(img, eye_x, eye_y, 1, (*DARK_BLUE, 255))
    px(img, eye_x-1, eye_y-1, (*WHITE, 200))  # 高光

    # ── 7. 紫色魚尾（扇形）──
    tail_cx, tail_cy = 10, 32
    for angle in range(-50, 51, 8):
        rad = math.radians(angle)
        for r in range(1, 12):
            tx = int(tail_cx - r * math.cos(rad))
            ty = int(tail_cy + r * math.sin(rad))
            alpha = max(60, 220 - r * 15)
            if r < 6:
                px(img, tx, ty, (*PURPLE, alpha))
            else:
                px(img, tx, ty, (*LIGHT_PURPLE, alpha))

    # ── 8. 時空粒子（散落的小點）──
    particles = [
        (15, 15, CYAN, 180), (50, 18, GOLD, 160), (48, 46, CYAN, 140),
        (18, 48, PURPLE, 150), (55, 32, WHITE, 200), (8, 32, GOLD, 170),
        (32, 10, CYAN, 160), (32, 54, PURPLE, 140), (22, 22, WHITE, 120),
        (44, 42, GOLD, 130), (14, 42, CYAN, 110), (50, 24, PURPLE, 120),
    ]
    for px_x, px_y, color, alpha in particles:
        px(img, px_x, px_y, (*color, alpha))
        # 十字形小粒子
        px(img, px_x+1, px_y, (*color, alpha//2))
        px(img, px_x-1, px_y, (*color, alpha//2))
        px(img, px_x, px_y+1, (*color, alpha//2))
        px(img, px_x, px_y-1, (*color, alpha//2))

    # ── 9. 輪廓（深藍黑）──
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 0:
                for dx, dy in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx_, ny_ = x+dx, y+dy
                    if 0 <= nx_ < SIZE and 0 <= ny_ < SIZE:
                        if img.getpixel((nx_, ny_))[3] == 0:
                            px(img, nx_, ny_, (*DARK_BLUE, 200))

    # 統計非透明像素
    non_transparent = sum(1 for y in range(SIZE) for x in range(SIZE) if img.getpixel((x,y))[3] > 0)
    total = SIZE * SIZE
    print(f"T249 非透明像素: {non_transparent}/{total} = {non_transparent/total*100:.1f}%")

    # 儲存
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    out_path = os.path.join(OUTPUT_DIR, "T249_time_rift_v2.png")
    img.save(out_path)
    print(f"已儲存: {out_path}")
    return out_path

if __name__ == "__main__":
    generate_t249()
    print("T249 精靈圖生成完成！")
