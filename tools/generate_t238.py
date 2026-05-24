# -*- coding: utf-8 -*-
"""
T238 幸運連鎖稀有魚精靈圖生成（DAY-280）
連鎖稀有主題：橙紅+金+火橙+天藍，連鎖鎖鏈+稀有光環+層次感
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

def gen_T238():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    cx, cy = 32, 36

    # 主色：橙紅漸層
    BODY_MAIN = (255, 107, 53)   # #FF6B35 橙紅
    BODY_DARK = (200, 60, 20)    # 深橙紅
    BODY_LIGHT = (255, 160, 80)  # 淡橙
    ACCENT = (255, 215, 0)       # 金色
    CHAIN_COLOR = (0, 191, 255)  # 天藍（連鎖色）

    # 魚身（橢圓，橙紅漸層）
    for y in range(cy-11, cy+12):
        for x in range(cx-19, cx+20):
            dx = x - cx
            dy = y - cy
            if (dx/19)**2 + (dy/11)**2 <= 1.0:
                nx_ = dx/19
                ny_ = dy/11
                dot = -(nx_*(-0.7) + ny_*(-0.7))
                if dot > 0.25:
                    c = BODY_LIGHT
                elif dot < -0.1:
                    c = BODY_DARK
                else:
                    c = BODY_MAIN
                px(img, x, y, c + (255,))

    # 魚尾（三角形，橙紅）
    for y in range(cy-8, cy+9):
        for x in range(cx+14, cx+26):
            dx = x - (cx+14)
            dy = abs(y - cy)
            if dx > 0 and dy <= 8 - dx*0.6:
                px(img, x, y, BODY_DARK + (220,))

    # 魚鰭（上方，金色）
    for y in range(cy-15, cy-8):
        for x in range(cx-5, cx+6):
            dy = y - (cy-15)
            dx = abs(x - cx)
            if dy >= 0 and dx <= 5 - dy*0.4:
                px(img, x, y, ACCENT + (200,))

    # 輪廓
    for y in range(cy-13, cy+14):
        for x in range(cx-21, cx+22):
            dx = x - cx
            dy = y - cy
            if abs((dx/19)**2 + (dy/11)**2 - 1.0) < 0.15:
                px(img, x, y, (80, 30, 0, 255))

    # 眼睛（白+黑+高光）
    fill_circle_shaded(img, cx-8, cy-3, 4, (255, 255, 255))
    fill_circle_shaded(img, cx-8, cy-3, 2, (20, 20, 20))
    px(img, cx-9, cy-4, (255, 255, 255, 255))

    # 連鎖鎖鏈（魚身上方，天藍色鏈環）
    chain_y = cy - 20
    for i in range(5):
        link_x = cx - 16 + i * 8
        # 橢圓鏈環
        for angle_deg in range(0, 360, 10):
            angle = math.radians(angle_deg)
            lx = int(link_x + 4 * math.cos(angle))
            ly = int(chain_y + 2 * math.sin(angle))
            px(img, lx, ly, CHAIN_COLOR + (255,))
            px(img, lx, ly+1, CHAIN_COLOR + (200,))
        # 鏈環連接點
        if i < 4:
            px(img, link_x+4, chain_y, CHAIN_COLOR + (255,))
            px(img, link_x+5, chain_y, CHAIN_COLOR + (255,))

    # 稀有光環（金色，圍繞魚身）
    for angle_deg in range(0, 360, 15):
        angle = math.radians(angle_deg)
        gx = int(cx + 23 * math.cos(angle))
        gy = int(cy + 14 * math.sin(angle))
        px(img, gx, gy, ACCENT + (150,))
        px(img, gx+1, gy, ACCENT + (100,))

    # 層次感光點（5個，代表5層連鎖）
    layer_pts = [
        (cx-14, cy+18, (0, 191, 255)),    # 第1層：天藍
        (cx-7,  cy+20, (80, 255, 80)),    # 第2層：草綠
        (cx,    cy+21, (255, 215, 0)),    # 第3層：金
        (cx+7,  cy+20, (255, 107, 53)),   # 第4層：橙紅
        (cx+14, cy+18, (255, 69, 0)),     # 第5層：火橙
    ]
    for lx, ly, lc in layer_pts:
        fill_circle_shaded(img, lx, ly, 2, lc)
        px(img, lx, ly-1, (255, 255, 255, 180))  # 高光

    # 火焰尾跡（魚尾後方，橙紅→金）
    for i in range(6):
        fx = cx + 22 + i * 2
        fy = cy + (i % 3 - 1)
        alpha = max(0, 200 - i * 30)
        r = min(255, 255 - i * 10)
        g = min(255, 100 + i * 20)
        px(img, fx, fy, (r, g, 0, alpha))

    return img

def main():
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    fname = "T238_rare_chain.png"
    fpath = os.path.join(OUTPUT_DIR, fname)
    img = gen_T238()
    img.save(fpath)
    print(f"生成完成：{fpath}")
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    print(f"非透明像素：{non_transparent}/{SIZE*SIZE} ({non_transparent*100//(SIZE*SIZE)}%)")

if __name__ == "__main__":
    main()
