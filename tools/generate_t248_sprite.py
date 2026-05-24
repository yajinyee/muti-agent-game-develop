"""
generate_t248_sprite.py — T248 幸運怒氣蓄積魚精靈圖生成（DAY-290）
視覺設計：怒氣主題（深紅 + 火橙 + 橙紅 + 金 + 白）
  - 橢圓火橙漸層魚身（怒氣感）
  - 怒氣火焰光環（深紅+火橙）
  - 4方向火焰光芒
  - 怒氣計量條（右側，深紅→火橙→金）
  - 金色魚眼（怒氣之眼）
  - 深紅魚尾
"""
import os
from PIL import Image

SIZE = 64
OUT_PATH = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets\T248_wrath_charge.png"

def px(img, x, y, color):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), color)

def fill_ellipse_shaded(img, cx, cy, rx, ry, light, mid, dark):
    for dy in range(-ry, ry+1):
        for dx in range(-rx, rx+1):
            if (dx/rx)**2 + (dy/ry)**2 <= 1.0:
                if dx < -rx//3 and dy < -ry//3:
                    c = light
                elif dx > rx//3 and dy > ry//3:
                    c = dark
                else:
                    c = mid
                px(img, cx+dx, cy+dy, c)

def fill_circle(img, cx, cy, r, color):
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                px(img, cx+dx, cy+dy, color)

img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

# 顏色定義
DEEP_RED    = (139, 0, 0, 255)      # #8B0000 深紅
FIRE        = (255, 69, 0, 255)     # #FF4500 火橙
ORANGE_RED  = (255, 107, 53, 255)   # #FF6B35 橙紅
FIRE_LIGHT  = (255, 140, 60, 255)   # 亮火橙（高光）
GOLD        = (255, 215, 0, 255)    # #FFD700 金
WHITE       = (255, 255, 255, 255)  # 白
BLACK       = (20, 0, 0, 255)       # 深黑輪廓

cx, cy = 30, 32

# 1. 外圈怒氣光環（深紅，半透明）
for dy in range(-22, 23):
    for dx in range(-22, 23):
        d = dx*dx + dy*dy
        if 380 <= d <= 484:
            alpha = int(160 * (1 - (d - 380) / 104))
            px(img, cx+dx, cy+dy, (139, 0, 0, max(0, alpha)))

# 2. 主體魚身（橢圓，火橙漸層）
fill_ellipse_shaded(img, cx, cy, 16, 13, FIRE_LIGHT, FIRE, DEEP_RED)

# 3. 輪廓
for dy in range(-14, 15):
    for dx in range(-17, 18):
        if (dx/16)**2 + (dy/13)**2 > 1.0 and (dx/17)**2 + (dy/14)**2 <= 1.0:
            px(img, cx+dx, cy+dy, BLACK)

# 4. 4方向火焰光芒
for i in range(14):
    alpha = max(0, 220 - i*16)
    # 上
    px(img, cx, cy - 16 - i, (FIRE[0], FIRE[1], FIRE[2], alpha))
    # 下
    px(img, cx, cy + 16 + i, (FIRE[0], FIRE[1], FIRE[2], alpha))
    # 左
    px(img, cx - 16 - i, cy, (FIRE[0], FIRE[1], FIRE[2], alpha))
    # 右（較短，因為有魚尾）
    if i < 8:
        px(img, cx + 16 + i, cy, (FIRE[0], FIRE[1], FIRE[2], alpha))

# 斜向光芒
for i in range(8):
    alpha = max(0, 160 - i*20)
    px(img, cx - 12 - i, cy - 10 - i, (FIRE[0], FIRE[1], FIRE[2], alpha))
    px(img, cx + 12 + i, cy - 10 - i, (FIRE[0], FIRE[1], FIRE[2], alpha))
    px(img, cx - 12 - i, cy + 10 + i, (FIRE[0], FIRE[1], FIRE[2], alpha))
    px(img, cx + 12 + i, cy + 10 + i, (FIRE[0], FIRE[1], FIRE[2], alpha))

# 5. 怒氣計量條（右側，垂直）
bar_x = 56
for i in range(20):
    bar_y = 12 + i * 2
    # 背景
    px(img, bar_x, bar_y, (50, 0, 0, 200))
    px(img, bar_x + 1, bar_y, (50, 0, 0, 200))
    # 填充（漸層：深紅→火橙→金）
    if i < 7:
        c = DEEP_RED
    elif i < 14:
        c = FIRE
    else:
        c = GOLD
    px(img, bar_x, bar_y, c)
    px(img, bar_x + 1, bar_y, c)

# 計量條邊框
for i in range(20):
    bar_y = 12 + i * 2
    px(img, bar_x - 1, bar_y, BLACK)
    px(img, bar_x + 2, bar_y, BLACK)
px(img, bar_x - 1, 11, BLACK)
px(img, bar_x, 11, BLACK)
px(img, bar_x + 1, 11, BLACK)
px(img, bar_x + 2, 11, BLACK)
px(img, bar_x - 1, 51, BLACK)
px(img, bar_x, 51, BLACK)
px(img, bar_x + 1, 51, BLACK)
px(img, bar_x + 2, 51, BLACK)

# 6. 金色魚眼（怒氣之眼）
fill_circle(img, cx - 5, cy - 3, 3, WHITE)
fill_circle(img, cx - 5, cy - 3, 2, GOLD)
fill_circle(img, cx - 5, cy - 3, 1, (255, 100, 0, 255))
px(img, cx - 6, cy - 4, WHITE)  # 高光

# 7. 深紅魚尾（右側）
tail_points = [
    (cx + 14, cy - 5), (cx + 15, cy - 4), (cx + 16, cy - 3),
    (cx + 17, cy - 2), (cx + 18, cy - 1), (cx + 19, cy),
    (cx + 18, cy + 1), (cx + 17, cy + 2), (cx + 16, cy + 3),
    (cx + 15, cy + 4), (cx + 14, cy + 5),
    (cx + 16, cy - 2), (cx + 17, cy - 1), (cx + 18, cy),
    (cx + 17, cy + 1), (cx + 16, cy + 2),
]
for (tx, ty) in tail_points:
    px(img, tx, ty, DEEP_RED)
    px(img, tx + 1, ty, (DEEP_RED[0], DEEP_RED[1], DEEP_RED[2], 180))

# 8. 中心怒氣核心（金色小圓）
fill_circle(img, cx + 3, cy + 2, 3, GOLD)
fill_circle(img, cx + 3, cy + 2, 2, FIRE_LIGHT)
px(img, cx + 3, cy + 2, GOLD)

# 9. 火焰粒子散落
import random
rng = random.Random(248)
for _ in range(10):
    gx = rng.randint(cx - 14, cx + 14)
    gy = rng.randint(cy - 12, cy + 12)
    d = (gx - cx)**2 + (gy - cy)**2
    if d > 64:
        c = FIRE if rng.random() > 0.5 else GOLD
        px(img, gx, gy, c)

# 統計非透明像素
non_transparent = sum(1 for x in range(SIZE) for y in range(SIZE) if img.getpixel((x, y))[3] > 0)
total = SIZE * SIZE
print(f"T248 怒氣蓄積魚精靈圖：{non_transparent}/{total} = {non_transparent/total*100:.1f}% 非透明像素")

os.makedirs(os.path.dirname(OUT_PATH), exist_ok=True)
img.save(OUT_PATH)
print(f"已儲存：{OUT_PATH}")
