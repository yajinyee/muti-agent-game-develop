"""
generate_t247_sprite.py — T247 幸運永生 BOSS 魚精靈圖生成（DAY-289）
視覺設計：永生主題（深紅 + 紅 + 火橙 + 金 + 白）
  - 圓形深紅漸層魚身（永生 BOSS 感）
  - 5 個金色生命之心（代表 5 條命）
  - 深紅光環 + 火焰光芒
  - 金色王冠符文（BOSS 感）
  - 紅色魚眼（永生之眼）
  - 火橙魚尾
"""
import os
from PIL import Image, ImageDraw

SIZE = 64
OUT_PATH = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets\T247_immortal_boss.png"

def px(img, x, y, color):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), color)

def fill_circle(img, cx, cy, r, color):
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                px(img, cx+dx, cy+dy, color)

def fill_circle_shaded(img, cx, cy, r, light, mid, dark):
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                # 左上亮，右下暗
                if dx < -r//3 and dy < -r//3:
                    c = light
                elif dx > r//3 and dy > r//3:
                    c = dark
                else:
                    c = mid
                px(img, cx+dx, cy+dy, c)

def draw_heart(img, cx, cy, size, color):
    """畫一個小愛心（生命之心）"""
    # 簡化版愛心：兩個小圓 + 三角形
    r = max(1, size // 2)
    # 左圓
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                px(img, cx - r//2 + dx, cy - r//2 + dy, color)
    # 右圓
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                px(img, cx + r//2 + dx, cy - r//2 + dy, color)
    # 下三角
    for i in range(r * 2):
        w = r * 2 - i
        for dx in range(-w, w+1):
            px(img, cx + dx, cy + i//2, color)

def draw_crown(img, cx, cy, color):
    """畫一個小王冠"""
    # 底部橫條
    for dx in range(-8, 9):
        px(img, cx+dx, cy+2, color)
        px(img, cx+dx, cy+3, color)
    # 三個尖頂
    for dy in range(6):
        px(img, cx, cy - dy, color)  # 中間尖
    for dy in range(4):
        px(img, cx - 5, cy - dy, color)  # 左尖
        px(img, cx + 5, cy - dy, color)  # 右尖
    # 王冠裝飾點
    px(img, cx, cy - 6, (255, 215, 0, 255))  # 中間頂點金色
    px(img, cx - 5, cy - 4, (255, 215, 0, 255))
    px(img, cx + 5, cy - 4, (255, 215, 0, 255))

img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

# 顏色定義
DEEP_RED   = (139, 0, 0, 255)      # #8B0000 深紅
RED        = (200, 0, 0, 255)      # 紅
RED_LIGHT  = (220, 50, 50, 255)    # 亮紅（高光）
FIRE       = (255, 69, 0, 255)     # #FF4500 火橙
GOLD       = (255, 215, 0, 255)    # #FFD700 金
GOLD_DARK  = (200, 160, 0, 255)    # 暗金
WHITE      = (255, 255, 255, 255)  # 白
BLACK      = (20, 0, 0, 255)       # 深黑輪廓

cx, cy = 32, 32

# 1. 外圈光環（深紅，半透明）
for dy in range(-22, 23):
    for dx in range(-22, 23):
        d = dx*dx + dy*dy
        if 380 <= d <= 484:  # r=19.5~22
            alpha = int(180 * (1 - (d - 380) / 104))
            px(img, cx+dx, cy+dy, (139, 0, 0, max(0, alpha)))

# 2. 主體魚身（圓形，深紅漸層）
fill_circle_shaded(img, cx, cy, 18, RED_LIGHT, RED, DEEP_RED)

# 3. 輪廓
for dy in range(-19, 20):
    for dx in range(-19, 20):
        d = dx*dx + dy*dy
        if 324 <= d <= 361:  # r=18~19
            px(img, cx+dx, cy+dy, BLACK)

# 4. 4方向火焰光芒
for i in range(12):
    # 上
    px(img, cx, cy - 20 - i, (FIRE[0], FIRE[1], FIRE[2], max(0, 200 - i*16)))
    # 下
    px(img, cx, cy + 20 + i, (FIRE[0], FIRE[1], FIRE[2], max(0, 200 - i*16)))
    # 左
    px(img, cx - 20 - i, cy, (FIRE[0], FIRE[1], FIRE[2], max(0, 200 - i*16)))
    # 右
    px(img, cx + 20 + i, cy, (FIRE[0], FIRE[1], FIRE[2], max(0, 200 - i*16)))

# 斜向光芒（較短）
for i in range(8):
    d = i
    px(img, cx - 14 - d, cy - 14 - d, (FIRE[0], FIRE[1], FIRE[2], max(0, 160 - i*20)))
    px(img, cx + 14 + d, cy - 14 - d, (FIRE[0], FIRE[1], FIRE[2], max(0, 160 - i*20)))
    px(img, cx - 14 - d, cy + 14 + d, (FIRE[0], FIRE[1], FIRE[2], max(0, 160 - i*20)))
    px(img, cx + 14 + d, cy + 14 + d, (FIRE[0], FIRE[1], FIRE[2], max(0, 160 - i*20)))

# 5. 王冠（頂部，金色）
draw_crown(img, cx, cy - 10, GOLD)

# 6. 5 個生命之心（排成弧形）
import math
heart_positions = []
for i in range(5):
    angle = math.pi + (i - 2) * 0.35  # 底部弧形排列
    hx = int(cx + 11 * math.cos(angle))
    hy = int(cy + 11 * math.sin(angle))
    heart_positions.append((hx, hy))

for i, (hx, hy) in enumerate(heart_positions):
    # 全部都是金色（代表全部存活）
    draw_heart(img, hx, hy, 3, GOLD)

# 7. 紅色魚眼（永生之眼）
fill_circle(img, cx - 5, cy - 3, 3, WHITE)
fill_circle(img, cx - 5, cy - 3, 2, (200, 0, 0, 255))
fill_circle(img, cx - 5, cy - 3, 1, (255, 50, 50, 255))
px(img, cx - 6, cy - 4, WHITE)  # 高光

# 8. 火橙魚尾（右側）
tail_points = [
    (cx + 16, cy - 6), (cx + 17, cy - 5), (cx + 18, cy - 4),
    (cx + 19, cy - 3), (cx + 20, cy - 2), (cx + 21, cy - 1),
    (cx + 22, cy),
    (cx + 21, cy + 1), (cx + 20, cy + 2), (cx + 19, cy + 3),
    (cx + 18, cy + 4), (cx + 17, cy + 5), (cx + 16, cy + 6),
    (cx + 17, cy - 3), (cx + 18, cy - 2), (cx + 19, cy - 1),
    (cx + 20, cy), (cx + 19, cy + 1), (cx + 18, cy + 2),
]
for (tx, ty) in tail_points:
    px(img, tx, ty, FIRE)
    px(img, tx + 1, ty, (FIRE[0], FIRE[1], FIRE[2], 180))

# 9. 中心永生核心（金色小圓）
fill_circle(img, cx + 3, cy + 2, 3, GOLD)
fill_circle(img, cx + 3, cy + 2, 2, WHITE)
px(img, cx + 3, cy + 2, GOLD)

# 10. 金色光點散落
import random
rng = random.Random(247)
for _ in range(12):
    gx = rng.randint(cx - 15, cx + 15)
    gy = rng.randint(cy - 15, cy + 15)
    d = (gx - cx)**2 + (gy - cy)**2
    if d > 100:  # 不在中心
        px(img, gx, gy, GOLD)

# 統計非透明像素
non_transparent = sum(1 for x in range(SIZE) for y in range(SIZE) if img.getpixel((x, y))[3] > 0)
total = SIZE * SIZE
print(f"T247 永生 BOSS 魚精靈圖：{non_transparent}/{total} = {non_transparent/total*100:.1f}% 非透明像素")

os.makedirs(os.path.dirname(OUT_PATH), exist_ok=True)
img.save(OUT_PATH)
print(f"已儲存：{OUT_PATH}")
