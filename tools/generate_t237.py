# -*- coding: utf-8 -*-
"""
T237 幸運彩虹橋魚精靈圖生成（DAY-279）
彩虹橋主題：粉紅+金+天藍+草綠+橙，彩虹弧橋+連接光點+彩虹魚身
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

def outline_circle(img, cx, cy, r, color):
    for y in range(max(0, cy-r-2), min(SIZE, cy+r+3)):
        for x in range(max(0, cx-r-2), min(SIZE, cx+r+3)):
            d = math.sqrt((x-cx)**2 + (y-cy)**2)
            if r+0.1 <= d <= r+1.5:
                px(img, x, y, color)

def gen_T237():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    cx, cy = 32, 36

    # 彩虹顏色
    rainbow = [
        (255, 80, 150),   # 粉紅
        (255, 180, 0),    # 金黃
        (0, 200, 255),    # 天藍
        (80, 255, 80),    # 草綠
        (255, 140, 0),    # 橙
        (180, 80, 255),   # 紫
    ]

    # 魚身（橢圓，彩虹漸層）— 加大魚身
    for y in range(cy-12, cy+13):
        for x in range(cx-20, cx+21):
            dx = x - cx
            dy = y - cy
            if (dx/20)**2 + (dy/12)**2 <= 1.0:
                # 彩虹漸層：依 x 位置選顏色
                t = (dx + 20) / 40.0  # 0~1
                idx = int(t * (len(rainbow)-1))
                idx = min(idx, len(rainbow)-2)
                frac = t * (len(rainbow)-1) - idx
                r1, g1, b1 = rainbow[idx]
                r2, g2, b2 = rainbow[idx+1]
                r = int(r1 + (r2-r1)*frac)
                g = int(g1 + (g2-g1)*frac)
                b = int(b1 + (b2-b1)*frac)
                # 陰影
                nx_ = dx/20
                ny_ = dy/12
                dot = -(nx_*(-0.7) + ny_*(-0.7))
                if dot > 0.25:
                    r, g, b = min(255,r+40), min(255,g+40), min(255,b+40)
                elif dot < -0.1:
                    r, g, b = max(0,r-40), max(0,g-40), max(0,b-40)
                px(img, x, y, (r, g, b, 255))

    # 魚尾（三角形，彩虹）
    for y in range(cy-8, cy+9):
        for x in range(cx+14, cx+26):
            dx = x - (cx+14)
            dy = abs(y - cy)
            if dx > 0 and dy <= 8 - dx*0.6:
                t = dx / 12.0
                idx = int(t * (len(rainbow)-1))
                idx = min(idx, len(rainbow)-2)
                frac = t * (len(rainbow)-1) - idx
                r1, g1, b1 = rainbow[idx]
                r2, g2, b2 = rainbow[idx+1]
                r = int(r1 + (r2-r1)*frac)
                g = int(g1 + (g2-g1)*frac)
                b = int(b1 + (b2-b1)*frac)
                px(img, x, y, (r, g, b, 220))

    # 魚鰭（上方，粉紅）
    for y in range(cy-16, cy-8):
        for x in range(cx-6, cx+7):
            dy = y - (cy-16)
            dx = abs(x - cx)
            if dy >= 0 and dx <= 6 - dy*0.5:
                px(img, x, y, (255, 100, 180, 200))

    # 輪廓
    for y in range(cy-14, cy+15):
        for x in range(cx-22, cx+23):
            dx = x - cx
            dy = y - cy
            if abs((dx/20)**2 + (dy/12)**2 - 1.0) < 0.15:
                px(img, x, y, (80, 20, 60, 255))

    # 眼睛（白+黑+高光）
    fill_circle_shaded(img, cx-8, cy-3, 4, (255, 255, 255))
    fill_circle_shaded(img, cx-8, cy-3, 2, (20, 20, 20))
    px(img, cx-9, cy-4, (255, 255, 255, 255))  # 高光

    # 彩虹橋弧線（上方，三條彩虹弧）
    for arc_r, arc_color in [(20, (255,80,150,200)), (16, (0,200,255,200)), (12, (80,255,80,200))]:
        for angle_deg in range(0, 181, 2):
            angle = math.radians(angle_deg)
            ax = int(cx + arc_r * math.cos(math.pi - angle))
            ay = int(cy - 18 - arc_r * math.sin(angle) * 0.5)
            px(img, ax, ay, arc_color)
            px(img, ax, ay-1, arc_color)

    # 連接光點（彩虹橋兩端）
    for dot_x, dot_color in [(cx-18, (255,80,150,255)), (cx+18, (0,200,255,255))]:
        fill_circle_shaded(img, dot_x, cy-18, 3, dot_color[:3])
        outline_circle(img, dot_x, cy-18, 3, (255,255,255,180))

    # 彩虹光點（散落）
    light_pts = [
        (cx-12, cy-28, (255,80,150)),
        (cx,    cy-32, (255,200,0)),
        (cx+12, cy-28, (0,200,255)),
        (cx-6,  cy-30, (80,255,80)),
        (cx+6,  cy-30, (255,140,0)),
    ]
    for lx, ly, lc in light_pts:
        fill_circle_shaded(img, lx, ly, 2, lc)

    # 金色光暈（魚身周圍）
    for angle_deg in range(0, 360, 30):
        angle = math.radians(angle_deg)
        gx = int(cx + 22 * math.cos(angle))
        gy = int(cy + 12 * math.sin(angle))
        px(img, gx, gy, (255, 220, 0, 120))
        px(img, gx+1, gy, (255, 220, 0, 80))

    return img

def main():
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    fname = "T237_rainbow_bridge.png"
    fpath = os.path.join(OUTPUT_DIR, fname)
    img = gen_T237()
    img.save(fpath)
    print(f"生成完成：{fpath}")
    # 統計非透明像素
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    print(f"非透明像素：{non_transparent}/{SIZE*SIZE} ({non_transparent*100//(SIZE*SIZE)}%)")

if __name__ == "__main__":
    main()
