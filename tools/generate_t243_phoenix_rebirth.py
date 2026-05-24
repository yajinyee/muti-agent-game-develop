# -*- coding: utf-8 -*-
"""
T243 幸運鳳凰涅槃魚精靈圖生成（DAY-285）
設計：鳳凰主題（火橙+金+深橙+白）
  - 橢圓魚身（火橙漸層）
  - 鳳凰羽毛紋路（金色弧線）
  - 涅槃火焰光環（深橙+金色）
  - 4方向火焰光芒
  - 鳳凰尾羽（橙紅漸層）
  - 中心涅槃核心（白色+金色）
"""
from PIL import Image, ImageDraw
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
    light = (min(255,r_v+40), min(255,g_v+30), min(255,b_v+10), 255)
    mid   = (r_v, g_v, b_v, 255)
    dark  = (max(0,r_v-50), max(0,g_v-40), max(0,b_v-20), 255)
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

def draw_arc_pixels(img, cx, cy, r, a_start, a_end, color, steps=40):
    """畫弧線（像素點）"""
    for i in range(steps+1):
        t = a_start + (a_end - a_start) * i / steps
        x = int(cx + r * math.cos(math.radians(t)))
        y = int(cy + r * math.sin(math.radians(t)))
        px(img, x, y, color)
        # 加粗
        px(img, x+1, y, color)
        px(img, x, y+1, color)

def draw_flame_ray(img, cx, cy, angle_deg, length, color_inner, color_outer):
    """畫火焰光芒"""
    angle = math.radians(angle_deg)
    for i in range(length):
        t = i / max(length-1, 1)
        x = int(cx + i * math.cos(angle))
        y = int(cy + i * math.sin(angle))
        # 顏色從內到外漸變
        r = int(color_inner[0] * (1-t) + color_outer[0] * t)
        g = int(color_inner[1] * (1-t) + color_outer[1] * t)
        b = int(color_inner[2] * (1-t) + color_outer[2] * t)
        a = int(255 * (1.0 - t * 0.6))
        px(img, x, y, (r, g, b, a))
        # 加粗（垂直方向）
        perp_x = int(-math.sin(angle))
        perp_y = int(math.cos(angle))
        if i < length * 0.6:
            px(img, x+perp_x, y+perp_y, (r, g, b, max(0, a-60)))

def generate_t243():
    img = new_img()

    # ── 顏色定義 ──────────────────────────────────────────────────────────────
    FIRE_ORANGE  = (255, 107, 53)   # #FF6B35 火橙（魚身主色）
    GOLD         = (255, 215, 0)    # #FFD700 金（羽毛/光芒）
    DEEP_ORANGE  = (255, 69, 0)     # #FF4500 深橙（陰影/尾羽）
    LIGHT_ORANGE = (255, 200, 150)  # 淡橙（高光）
    WHITE        = (255, 255, 255)  # 白（涅槃核心）
    OUTLINE      = (120, 30, 0)     # 深棕（輪廓）
    PHOENIX_RED  = (220, 50, 20)    # 鳳凰紅（尾羽深色）

    cx, cy = 32, 32

    # ── 1. 涅槃火焰光環（最底層）──────────────────────────────────────────────
    # 外圈光環（深橙，半透明）
    for r in range(28, 31):
        for angle in range(0, 360, 3):
            x = int(cx + r * math.cos(math.radians(angle)))
            y = int(cy + r * math.sin(math.radians(angle)))
            alpha = 120 - (r - 28) * 30
            px(img, x, y, (DEEP_ORANGE[0], DEEP_ORANGE[1], DEEP_ORANGE[2], alpha))

    # 內圈光環（金色，半透明）
    for r in range(24, 27):
        for angle in range(0, 360, 4):
            x = int(cx + r * math.cos(math.radians(angle)))
            y = int(cy + r * math.sin(math.radians(angle)))
            alpha = 100 - (r - 24) * 20
            px(img, x, y, (GOLD[0], GOLD[1], GOLD[2], alpha))

    # ── 2. 4方向火焰光芒 ──────────────────────────────────────────────────────
    # 右（主方向）
    draw_flame_ray(img, cx+18, cy, 0, 12, GOLD, (255, 150, 0))
    # 左
    draw_flame_ray(img, cx-18, cy, 180, 10, GOLD, (255, 150, 0))
    # 上
    draw_flame_ray(img, cx, cy-16, 270, 10, GOLD, (255, 150, 0))
    # 下
    draw_flame_ray(img, cx, cy+16, 90, 10, GOLD, (255, 150, 0))
    # 斜向（較短）
    draw_flame_ray(img, cx+13, cy-13, 315, 8, FIRE_ORANGE, (255, 180, 50))
    draw_flame_ray(img, cx+13, cy+13, 45, 8, FIRE_ORANGE, (255, 180, 50))
    draw_flame_ray(img, cx-13, cy-13, 225, 7, FIRE_ORANGE, (255, 180, 50))
    draw_flame_ray(img, cx-13, cy+13, 135, 7, FIRE_ORANGE, (255, 180, 50))

    # ── 3. 鳳凰尾羽（左側，橙紅漸層）────────────────────────────────────────
    # 上尾羽
    for i in range(12):
        t = i / 11.0
        x = cx - 14 - i
        y = cy - 4 - int(i * 0.8)
        r = int(PHOENIX_RED[0] * (1-t) + DEEP_ORANGE[0] * t)
        g = int(PHOENIX_RED[1] * (1-t) + DEEP_ORANGE[1] * t)
        b = int(PHOENIX_RED[2] * (1-t) + DEEP_ORANGE[2] * t)
        a = int(255 * (1.0 - t * 0.4))
        px(img, x, y, (r, g, b, a))
        px(img, x, y+1, (r, g, b, max(0, a-40)))
    # 中尾羽
    for i in range(14):
        t = i / 13.0
        x = cx - 14 - i
        y = cy + int(i * 0.1)
        r = int(FIRE_ORANGE[0] * (1-t) + PHOENIX_RED[0] * t)
        g = int(FIRE_ORANGE[1] * (1-t) + PHOENIX_RED[1] * t)
        b = int(FIRE_ORANGE[2] * (1-t) + PHOENIX_RED[2] * t)
        a = int(255 * (1.0 - t * 0.3))
        px(img, x, y, (r, g, b, a))
        px(img, x, y-1, (r, g, b, max(0, a-40)))
    # 下尾羽
    for i in range(12):
        t = i / 11.0
        x = cx - 14 - i
        y = cy + 4 + int(i * 0.8)
        r = int(PHOENIX_RED[0] * (1-t) + DEEP_ORANGE[0] * t)
        g = int(PHOENIX_RED[1] * (1-t) + DEEP_ORANGE[1] * t)
        b = int(PHOENIX_RED[2] * (1-t) + DEEP_ORANGE[2] * t)
        a = int(255 * (1.0 - t * 0.4))
        px(img, x, y, (r, g, b, a))
        px(img, x, y-1, (r, g, b, max(0, a-40)))

    # ── 4. 魚身（橢圓，火橙漸層）────────────────────────────────────────────
    fill_ellipse_shaded(img, cx, cy, 18, 12, FIRE_ORANGE)

    # ── 5. 鳳凰羽毛紋路（金色弧線）──────────────────────────────────────────
    # 上弧
    draw_arc_pixels(img, cx, cy-4, 12, 200, 340, GOLD, 20)
    # 下弧
    draw_arc_pixels(img, cx, cy+4, 12, 20, 160, GOLD, 20)
    # 中弧（較小）
    draw_arc_pixels(img, cx-2, cy, 8, 210, 330, (255, 230, 100), 15)

    # ── 6. 魚眼（金色+黑色）──────────────────────────────────────────────────
    fill_circle(img, cx+12, cy-3, 3, GOLD)
    fill_circle(img, cx+12, cy-3, 2, (30, 20, 0, 255))
    # 眼睛高光
    px(img, cx+13, cy-4, (255, 255, 200, 255))

    # ── 7. 中心涅槃核心（白色+金色光點）──────────────────────────────────────
    fill_circle(img, cx, cy, 4, WHITE)
    fill_circle(img, cx, cy, 2, GOLD)
    # 核心光芒（4方向小點）
    for angle in [0, 90, 180, 270]:
        x = int(cx + 6 * math.cos(math.radians(angle)))
        y = int(cy + 6 * math.sin(math.radians(angle)))
        px(img, x, y, (255, 255, 200, 200))

    # ── 8. 輪廓（深棕）────────────────────────────────────────────────────────
    # 橢圓輪廓
    for angle in range(0, 360, 3):
        x = int(cx + 18 * math.cos(math.radians(angle)))
        y = int(cy + 12 * math.sin(math.radians(angle)))
        px(img, x, y, OUTLINE + (255,))

    # ── 9. 金色光點散落 ────────────────────────────────────────────────────────
    import random
    rng = random.Random(243)
    for _ in range(18):
        gx = rng.randint(4, 60)
        gy = rng.randint(4, 60)
        # 只在魚身外圍放光點
        dist = math.sqrt((gx-cx)**2 + (gy-cy)**2)
        if 20 < dist < 30:
            alpha = rng.randint(100, 200)
            px(img, gx, gy, (GOLD[0], GOLD[1], GOLD[2], alpha))

    return img

if __name__ == "__main__":
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    img = generate_t243()
    out_path = os.path.join(OUTPUT_DIR, "T243_phoenix_rebirth.png")
    img.save(out_path)
    print(f"T243 saved: {out_path}")

    # 統計非透明像素
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    total = SIZE * SIZE
    print(f"Non-transparent pixels: {non_transparent}/{total} ({non_transparent*100//total}%)")
