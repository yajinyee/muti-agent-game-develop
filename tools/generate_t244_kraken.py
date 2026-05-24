# -*- coding: utf-8 -*-
"""
T244 幸運深海克拉肯魚精靈圖生成（DAY-286）
設計：深海克拉肯主題（深藍+青藍+紫+白）
  - 圓形魚身（深藍漸層）
  - 8條觸手紋路（青藍弧線）
  - 深海光環（紫色+青藍）
  - 4方向觸手光芒
  - 克拉肯眼睛（黃色+黑色）
  - 中心深海核心（青藍+白色）
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
    light = (min(255,r_v+30), min(255,g_v+40), min(255,b_v+50), 255)
    mid   = (r_v, g_v, b_v, 255)
    dark  = (max(0,r_v-30), max(0,g_v-30), max(0,b_v-40), 255)
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

def draw_tentacle(img, cx, cy, angle_deg, length, color_base, color_tip):
    """畫觸手（彎曲弧線）"""
    angle = math.radians(angle_deg)
    curve = math.radians(25)  # 觸手彎曲角度
    for i in range(length):
        t = i / max(length-1, 1)
        cur_angle = angle + curve * t
        x = int(cx + i * math.cos(cur_angle))
        y = int(cy + i * math.sin(cur_angle))
        r = int(color_base[0] * (1-t) + color_tip[0] * t)
        g = int(color_base[1] * (1-t) + color_tip[1] * t)
        b = int(color_base[2] * (1-t) + color_tip[2] * t)
        a = int(255 * (1.0 - t * 0.5))
        px(img, x, y, (r, g, b, a))
        # 觸手加粗（靠近根部）
        if t < 0.4:
            perp_x = int(-math.sin(cur_angle))
            perp_y = int(math.cos(cur_angle))
            px(img, x+perp_x, y+perp_y, (r, g, b, max(0, a-80)))

def generate_t244():
    img = new_img()

    # ── 顏色定義 ──────────────────────────────────────────────────────────────
    DEEP_BLUE    = (10, 30, 80)     # #0A1E50 深藍（魚身主色）
    CYAN_BLUE    = (0, 180, 220)    # #00B4DC 青藍（觸手/光芒）
    PURPLE       = (100, 50, 180)   # #6432B4 紫（光環）
    LIGHT_CYAN   = (150, 230, 255)  # 淡青（高光）
    WHITE        = (255, 255, 255)  # 白（核心）
    YELLOW       = (255, 220, 0)    # 黃（眼睛）
    OUTLINE      = (0, 10, 40)      # 深藍黑（輪廓）
    TEAL         = (0, 150, 160)    # 青綠（觸手尖端）

    cx, cy = 32, 32

    # ── 1. 深海光環（最底層）──────────────────────────────────────────────────
    # 外圈光環（紫色，半透明）
    for r in range(28, 31):
        for angle in range(0, 360, 3):
            x = int(cx + r * math.cos(math.radians(angle)))
            y = int(cy + r * math.sin(math.radians(angle)))
            alpha = 100 - (r - 28) * 25
            px(img, x, y, (PURPLE[0], PURPLE[1], PURPLE[2], alpha))

    # 內圈光環（青藍，半透明）
    for r in range(23, 26):
        for angle in range(0, 360, 4):
            x = int(cx + r * math.cos(math.radians(angle)))
            y = int(cy + r * math.sin(math.radians(angle)))
            alpha = 90 - (r - 23) * 20
            px(img, x, y, (CYAN_BLUE[0], CYAN_BLUE[1], CYAN_BLUE[2], alpha))

    # ── 2. 8條觸手（從魚身向外延伸）──────────────────────────────────────────
    tentacle_angles = [0, 45, 90, 135, 180, 225, 270, 315]
    for i, angle in enumerate(tentacle_angles):
        start_x = int(cx + 14 * math.cos(math.radians(angle)))
        start_y = int(cy + 14 * math.sin(math.radians(angle)))
        length = 12 if i % 2 == 0 else 10  # 主觸手長，副觸手短
        draw_tentacle(img, start_x, start_y, angle, length, CYAN_BLUE, TEAL)

    # ── 3. 魚身（圓形，深藍漸層）────────────────────────────────────────────
    fill_circle_shaded(img, cx, cy, 14, DEEP_BLUE)

    # ── 4. 觸手紋路（青藍弧線，在魚身上）────────────────────────────────────
    # 4條主紋路
    for angle in [30, 120, 210, 300]:
        for r in range(6, 13):
            x = int(cx + r * math.cos(math.radians(angle)))
            y = int(cy + r * math.sin(math.radians(angle)))
            alpha = 180 - (r - 6) * 15
            px(img, x, y, (CYAN_BLUE[0], CYAN_BLUE[1], CYAN_BLUE[2], alpha))

    # ── 5. 克拉肯眼睛（2個，黃色+黑色）──────────────────────────────────────
    # 左眼
    fill_circle(img, cx-5, cy-3, 3, YELLOW)
    fill_circle(img, cx-5, cy-3, 2, (20, 10, 0, 255))
    px(img, cx-4, cy-4, (255, 255, 200, 255))  # 高光
    # 右眼
    fill_circle(img, cx+5, cy-3, 3, YELLOW)
    fill_circle(img, cx+5, cy-3, 2, (20, 10, 0, 255))
    px(img, cx+6, cy-4, (255, 255, 200, 255))  # 高光

    # ── 6. 嘴巴（青藍弧線）────────────────────────────────────────────────────
    for dx in range(-4, 5):
        y_offset = int(2 * math.sin(math.pi * (dx + 4) / 8))
        px(img, cx + dx, cy + 5 + y_offset, (CYAN_BLUE[0], CYAN_BLUE[1], CYAN_BLUE[2], 200))

    # ── 7. 中心深海核心（青藍+白色）──────────────────────────────────────────
    fill_circle(img, cx, cy, 3, CYAN_BLUE)
    fill_circle(img, cx, cy, 1, WHITE)

    # ── 8. 輪廓（深藍黑）──────────────────────────────────────────────────────
    for angle in range(0, 360, 3):
        x = int(cx + 14 * math.cos(math.radians(angle)))
        y = int(cy + 14 * math.sin(math.radians(angle)))
        px(img, x, y, OUTLINE + (255,))

    # ── 9. 青藍光點散落 ────────────────────────────────────────────────────────
    import random
    rng = random.Random(244)
    for _ in range(20):
        gx = rng.randint(4, 60)
        gy = rng.randint(4, 60)
        dist = math.sqrt((gx-cx)**2 + (gy-cy)**2)
        if 18 < dist < 30:
            alpha = rng.randint(80, 180)
            px(img, gx, gy, (CYAN_BLUE[0], CYAN_BLUE[1], CYAN_BLUE[2], alpha))

    return img

if __name__ == "__main__":
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    img = generate_t244()
    out_path = os.path.join(OUTPUT_DIR, "T244_kraken.png")
    img.save(out_path)
    print(f"T244 saved: {out_path}")

    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    total = SIZE * SIZE
    print(f"Non-transparent pixels: {non_transparent}/{total} ({non_transparent*100//total}%)")
