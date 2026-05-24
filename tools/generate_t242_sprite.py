# -*- coding: utf-8 -*-
"""
generate_t242_sprite.py — T242 幸運龍怒隕石魚精靈圖生成（DAY-284）
設計：龍怒主題魚（火焰龍鱗魚身+隕石衝擊紋路+火橙光環+4方向火焰光芒+龍爪符文）
輸出：64x64 PNG（NEAREST 插值，像素風格）
"""
from PIL import Image
import math
import os

OUT_PATH = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets\T242_dragon_wrath.png"
SIZE = 64

def px(img, x, y, color):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), color)

def fill_circle(img, cx, cy, r, color):
    for y in range(max(0, cy - r), min(SIZE, cy + r + 1)):
        for x in range(max(0, cx - r), min(SIZE, cx + r + 1)):
            if (x - cx) ** 2 + (y - cy) ** 2 <= r * r:
                img.putpixel((x, y), color)

def fill_circle_shaded(img, cx, cy, r, light, mid, dark):
    for y in range(max(0, cy - r), min(SIZE, cy + r + 1)):
        for x in range(max(0, cx - r), min(SIZE, cx + r + 1)):
            if (x - cx) ** 2 + (y - cy) ** 2 <= r * r:
                if x < cx - 1 and y < cy - 1:
                    img.putpixel((x, y), light)
                elif x > cx + 1 and y > cy + 1:
                    img.putpixel((x, y), dark)
                else:
                    img.putpixel((x, y), mid)

def draw_ray(img, cx, cy, angle_deg, length, width, color):
    angle = math.radians(angle_deg)
    for d in range(length):
        x = int(cx + d * math.cos(angle))
        y = int(cy + d * math.sin(angle))
        for w in range(-width // 2, width // 2 + 1):
            px_x = x + int(w * math.sin(angle))
            px_y = y - int(w * math.cos(angle))
            if 0 <= px_x < SIZE and 0 <= px_y < SIZE:
                img.putpixel((px_x, px_y), color)

# 顏色定義（龍怒主題：火橙+深紅+金+黑）
OUTLINE      = (20,  10,   0, 255)   # 深黑輪廓（帶紅色調）
FIRE_LIGHT   = (255, 180,  50, 255)  # 火焰亮橙
FIRE_MID     = (255, 100,   0, 255)  # 火橙中（#FF6400）
FIRE_DARK    = (180,  40,   0, 255)  # 火橙暗
DEEP_RED_L   = (220,  50,  50, 255)  # 深紅亮
DEEP_RED_M   = (160,   0,   0, 255)  # 深紅中
DEEP_RED_D   = (100,   0,   0, 255)  # 深紅暗
GOLD_LIGHT   = (255, 240, 100, 255)  # 金色亮
GOLD_MID     = (255, 215,   0, 255)  # 金色中（#FFD700）
GOLD_DARK    = (200, 160,   0, 255)  # 金色暗
LAVA_L       = (255, 140,  20, 255)  # 熔岩亮
LAVA_M       = (220,  80,   0, 255)  # 熔岩中
LAVA_D       = (150,  40,   0, 255)  # 熔岩暗
WHITE        = (255, 255, 255, 255)  # 白色
TRANSPARENT  = (0,   0,   0,   0)

img = Image.new("RGBA", (SIZE, SIZE), TRANSPARENT)

# ── 1. 主體：圓形魚身（龍鱗火焰漸層）────────────────────────────────────────
cx, cy, r = 32, 32, 20
for y in range(max(0, cy - r), min(SIZE, cy + r + 1)):
    for x in range(max(0, cx - r), min(SIZE, cx + r + 1)):
        if (x - cx) ** 2 + (y - cy) ** 2 <= r * r:
            # 從中心向外漸層：中心金→中圈火橙→外圈深紅
            dist = math.sqrt((x - cx) ** 2 + (y - cy) ** 2)
            ratio = dist / r
            # 左上亮，右下暗（龍鱗光澤）
            nx_ = (x - cx) / max(r, 1)
            ny_ = (y - cy) / max(r, 1)
            dot = -(nx_ * (-0.7) + ny_ * (-0.7))

            if ratio < 0.35:
                # 中心：金色
                if dot > 0.2:
                    c = GOLD_LIGHT
                elif dot < -0.1:
                    c = GOLD_DARK
                else:
                    c = GOLD_MID
            elif ratio < 0.65:
                # 中圈：火橙
                if dot > 0.2:
                    c = FIRE_LIGHT
                elif dot < -0.1:
                    c = FIRE_DARK
                else:
                    c = FIRE_MID
            else:
                # 外圈：深紅
                if dot > 0.2:
                    c = DEEP_RED_L
                elif dot < -0.1:
                    c = DEEP_RED_D
                else:
                    c = DEEP_RED_M
            img.putpixel((x, y), c)

# ── 2. 龍鱗紋路（弧形鱗片）──────────────────────────────────────────────────
scale_positions = [
    (26, 26, 5), (38, 26, 5),
    (32, 32, 5),
    (26, 38, 5), (38, 38, 5),
]
for sx, sy, sr in scale_positions:
    if (sx - cx) ** 2 + (sy - cy) ** 2 <= (r - 2) ** 2:
        # 鱗片弧線（上半圓）
        for angle_deg in range(200, 340, 8):
            angle = math.radians(angle_deg)
            bx = int(sx + sr * math.cos(angle))
            by = int(sy + sr * math.sin(angle))
            if (bx - cx) ** 2 + (by - cy) ** 2 <= r * r:
                px(img, bx, by, DEEP_RED_D)

# ── 3. 輪廓（深黑帶紅色調）──────────────────────────────────────────────────
for angle_deg in range(0, 360, 2):
    angle = math.radians(angle_deg)
    x = int(cx + r * math.cos(angle))
    y = int(cy + r * math.sin(angle))
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), OUTLINE)

# ── 4. 火橙光環（外圈）──────────────────────────────────────────────────────
for angle_deg in range(0, 360, 3):
    angle = math.radians(angle_deg)
    x = int(cx + (r + 2) * math.cos(angle))
    y = int(cy + (r + 2) * math.sin(angle))
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), FIRE_MID)
# 外圈第二層（金色點綴）
for angle_deg in range(0, 360, 12):
    angle = math.radians(angle_deg)
    x = int(cx + (r + 3) * math.cos(angle))
    y = int(cy + (r + 3) * math.sin(angle))
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), GOLD_MID)

# ── 5. 4方向火焰光芒（龍怒光芒）────────────────────────────────────────────
# 上下左右四個方向
for angle in [270, 90, 180, 0]:  # 上/下/左/右
    draw_ray(img, cx, cy, angle, 26, 3, FIRE_MID)
    draw_ray(img, cx, cy, angle, 18, 1, GOLD_MID)  # 中心金色核心

# 斜向四個方向（較短）
for angle in [315, 45, 225, 135]:
    draw_ray(img, cx, cy, angle, 20, 2, FIRE_DARK)

# ── 6. 隕石衝擊紋路（放射狀裂縫）────────────────────────────────────────────
crack_angles = [30, 75, 120, 165, 210, 255, 300, 345]
for angle_deg in crack_angles:
    angle = math.radians(angle_deg)
    for d in range(8, 18):
        x = int(cx + d * math.cos(angle))
        y = int(cy + d * math.sin(angle))
        if (x - cx) ** 2 + (y - cy) ** 2 <= r * r:
            px(img, x, y, DEEP_RED_D)

# ── 7. 龍爪符文（四個角落）──────────────────────────────────────────────────
# 簡化龍爪：三條短線組成爪形
claw_positions = [
    (14, 14),  # 左上
    (50, 14),  # 右上
    (14, 50),  # 左下
    (50, 50),  # 右下
]
for clx, cly in claw_positions:
    # 爪形：中心點 + 三條短線
    px(img, clx, cly, GOLD_MID)
    # 三條爪線
    for i in range(1, 4):
        px(img, clx - i, cly - i, GOLD_DARK)  # 左上爪
        px(img, clx + i, cly - i, GOLD_DARK)  # 右上爪
        px(img, clx,     cly + i, GOLD_DARK)  # 下爪

# ── 8. 中心火焰核心 ──────────────────────────────────────────────────────────
fill_circle(img, cx, cy, 6, FIRE_MID)
fill_circle(img, cx, cy, 4, GOLD_MID)
fill_circle(img, cx, cy, 2, GOLD_LIGHT)
px(img, cx - 1, cy - 1, WHITE)  # 高光

# ── 9. 魚眼（火橙）──────────────────────────────────────────────────────────
fill_circle(img, 38, 26, 3, FIRE_MID)
fill_circle(img, 38, 26, 2, DEEP_RED_M)
px(img, 37, 25, WHITE)  # 高光

# ── 10. 魚尾（火焰扇形）─────────────────────────────────────────────────────
tail_pixels = [
    (52, 26), (53, 25), (54, 24), (55, 23), (56, 22),
    (52, 32), (53, 33), (54, 34), (55, 35), (56, 36),
    (52, 29), (53, 29), (54, 29), (55, 29),
]
for i, (tx, ty) in enumerate(tail_pixels):
    if 0 <= tx < SIZE and 0 <= ty < SIZE:
        c = FIRE_LIGHT if i % 3 == 0 else (FIRE_MID if i % 3 == 1 else FIRE_DARK)
        img.putpixel((tx, ty), c)

# ── 11. 火焰粒子散落 ─────────────────────────────────────────────────────────
fire_particles = [
    (6,  6,  FIRE_LIGHT), (7,  7,  FIRE_MID),
    (57, 6,  FIRE_LIGHT), (56, 7,  FIRE_MID),
    (6,  57, FIRE_LIGHT), (7,  56, FIRE_MID),
    (57, 57, FIRE_LIGHT), (56, 56, FIRE_MID),
    (32, 3,  GOLD_LIGHT), (32, 4,  GOLD_MID),
    (3,  32, GOLD_LIGHT), (4,  32, GOLD_MID),
    (60, 32, GOLD_LIGHT), (59, 32, GOLD_MID),
    (32, 60, GOLD_LIGHT), (32, 59, GOLD_MID),
    (20, 8,  FIRE_MID),   (44, 8,  FIRE_MID),
    (8,  20, FIRE_MID),   (8,  44, FIRE_MID),
    (56, 20, FIRE_MID),   (56, 44, FIRE_MID),
    (20, 56, FIRE_MID),   (44, 56, FIRE_MID),
]
for dx, dy, dc in fire_particles:
    if 0 <= dx < SIZE and 0 <= dy < SIZE:
        img.putpixel((dx, dy), dc)

# 儲存
os.makedirs(os.path.dirname(OUT_PATH), exist_ok=True)
img.save(OUT_PATH)
print(f"T242 sprite saved: {OUT_PATH}")

# 驗證
loaded = Image.open(OUT_PATH)
pixels = loaded.load()
non_transparent = sum(1 for y in range(SIZE) for x in range(SIZE) if pixels[x, y][3] > 0)
total = SIZE * SIZE
print(f"Non-transparent pixels: {non_transparent}/{total} ({non_transparent*100//total}%)")
