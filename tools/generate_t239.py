# -*- coding: utf-8 -*-
"""
T239 幸運黃金突變魚精靈圖生成（DAY-281）
黃金突變主題：金色+橙金+深橙，黃金光環+突變星芒+金色魚身
"""
from PIL import Image, ImageDraw
import os, math

OUTPUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def px(img, x, y, c):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), c)

def fill_circle_shaded(img, cx, cy, r, base_rgb):
    r_v, g_v, b_v = base_rgb
    for y in range(max(0, cy-r), min(SIZE, cy+r+1)):
        for x in range(max(0, cx-r), min(SIZE, cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 > r**2:
                continue
            nx_ = (x-cx)/max(r,1)
            ny_ = (y-cy)/max(r,1)
            dot = -(nx_*(-0.7) + ny_*(-0.7))
            if dot > 0.25:
                c = (min(255,r_v+45), min(255,g_v+45), min(255,b_v+45), 255)
            elif dot < -0.1:
                c = (max(0,r_v-50), max(0,g_v-50), max(0,b_v-50), 255)
            else:
                c = (r_v, g_v, b_v, 255)
            px(img, x, y, c)

def draw_star(img, cx, cy, r_out, r_in, points, color):
    pts = []
    for i in range(points * 2):
        angle = math.pi * i / points - math.pi / 2
        r = r_out if i % 2 == 0 else r_in
        pts.append((int(cx + r * math.cos(angle)), int(cy + r * math.sin(angle))))
    draw = ImageDraw.Draw(img)
    draw.polygon(pts, fill=color)

def gen_T239():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    cx, cy = 32, 36

    GOLD_MAIN  = (255, 215, 0)    # #FFD700 金
    GOLD_LIGHT = (255, 240, 100)  # 淡金
    GOLD_DARK  = (200, 160, 0)    # 深金
    ORANGE_GOLD = (255, 165, 0)   # #FFA500 橙金
    DEEP_ORANGE = (255, 140, 0)   # #FF8C00 深橙
    OUTLINE    = (120, 80, 0, 255) # 深棕輪廓

    # 魚身（橢圓，金色漸層）
    for y in range(cy-11, cy+12):
        for x in range(cx-19, cx+20):
            dx = x - cx
            dy = y - cy
            if (dx/19)**2 + (dy/11)**2 <= 1.0:
                nx_ = dx/19
                ny_ = dy/11
                dot = -(nx_*(-0.7) + ny_*(-0.7))
                if dot > 0.25:
                    c = GOLD_LIGHT
                elif dot < -0.1:
                    c = GOLD_DARK
                else:
                    c = GOLD_MAIN
                px(img, x, y, c + (255,))

    # 魚尾（三角形，橙金）
    for y in range(cy-8, cy+9):
        for x in range(cx+14, cx+26):
            dx = x - (cx+14)
            dy = abs(y - cy)
            if dx > 0 and dy <= 8 - dx*0.6:
                px(img, x, y, ORANGE_GOLD + (220,))

    # 魚鰭（上方，深橙）
    for y in range(cy-15, cy-8):
        for x in range(cx-5, cx+6):
            dy = y - (cy-15)
            dx = abs(x - cx)
            if dy >= 0 and dx <= 5 - dy*0.4:
                px(img, x, y, DEEP_ORANGE + (200,))

    # 輪廓
    for y in range(cy-13, cy+14):
        for x in range(cx-21, cx+22):
            dx = x - cx
            dy = y - cy
            if abs((dx/19)**2 + (dy/11)**2 - 1.0) < 0.15:
                px(img, x, y, OUTLINE)

    # 眼睛（白+黑+高光）
    fill_circle_shaded(img, cx-8, cy-3, 4, (255, 255, 255))
    fill_circle_shaded(img, cx-8, cy-3, 2, (20, 20, 20))
    px(img, cx-9, cy-4, (255, 255, 255, 255))

    # 黃金光環（外圈，金色）
    for angle_deg in range(0, 360, 8):
        angle = math.radians(angle_deg)
        gx = int(cx + 24 * math.cos(angle))
        gy = int(cy + 15 * math.sin(angle))
        px(img, gx, gy, GOLD_MAIN + (200,))
        px(img, gx+1, gy, GOLD_MAIN + (150,))

    # 突變星芒（4個，金色六芒星）
    star_pts = [
        (cx-2, cy-22),   # 上方
        (cx+2, cy+22),   # 下方（魚身下）
        (cx-22, cy-2),   # 左方
        (cx+22, cy+2),   # 右方（魚尾前）
    ]
    for sx, sy in star_pts:
        draw_star(img, sx, sy, 5, 2, 6, GOLD_MAIN + (220,))

    # 金色光點（散落）
    light_pts = [
        (cx-10, cy-18, GOLD_LIGHT),
        (cx+8,  cy-20, ORANGE_GOLD),
        (cx-18, cy+8,  GOLD_MAIN),
        (cx+16, cy-8,  GOLD_LIGHT),
        (cx,    cy-25, DEEP_ORANGE),
    ]
    for lx, ly, lc in light_pts:
        fill_circle_shaded(img, lx, ly, 2, lc)
        px(img, lx, ly-1, (255, 255, 255, 180))

    # 突變紋路（魚身上的金色條紋）
    for i in range(3):
        stripe_x = cx - 8 + i * 8
        for dy in range(-8, 9):
            if abs(dy) < 8 - abs(stripe_x - cx) * 0.3:
                px(img, stripe_x, cy + dy, GOLD_LIGHT + (180,))

    return img

def main():
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    fname = "T239_gold_mutation.png"
    fpath = os.path.join(OUTPUT_DIR, fname)
    img = gen_T239()
    img.save(fpath)
    print(f"生成完成：{fpath}")
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    print(f"非透明像素：{non_transparent}/{SIZE*SIZE} ({non_transparent*100//(SIZE*SIZE)}%)")

if __name__ == "__main__":
    main()
