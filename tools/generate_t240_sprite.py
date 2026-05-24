"""
generate_t240_sprite.py — T240 幸運星爆魚精靈圖生成（DAY-282）
設計：星爆主題魚（五角星魚身+爆炸光芒+天藍星爆點+金色光環+彩色星點）
輸出：64x64 PNG（NEAREST 插值，像素風格）
"""
from PIL import Image, ImageDraw
import math
import os

OUT_PATH = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets\T240_star_burst.png"
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

def draw_star(img, cx, cy, r_outer, r_inner, points, color, outline=None):
    """繪製多角星形"""
    coords = []
    for i in range(points * 2):
        angle = math.pi * i / points - math.pi / 2
        r = r_outer if i % 2 == 0 else r_inner
        x = int(cx + r * math.cos(angle))
        y = int(cy + r * math.sin(angle))
        coords.append((x, y))
    # 填充星形（掃描線法）
    if len(coords) < 3:
        return
    min_y = max(0, min(c[1] for c in coords))
    max_y = min(SIZE - 1, max(c[1] for c in coords))
    for y in range(min_y, max_y + 1):
        intersections = []
        n = len(coords)
        for i in range(n):
            x1, y1 = coords[i]
            x2, y2 = coords[(i + 1) % n]
            if (y1 <= y < y2) or (y2 <= y < y1):
                if y2 != y1:
                    x_int = x1 + (y - y1) * (x2 - x1) / (y2 - y1)
                    intersections.append(int(x_int))
        intersections.sort()
        for i in range(0, len(intersections) - 1, 2):
            for x in range(max(0, intersections[i]), min(SIZE, intersections[i + 1] + 1)):
                img.putpixel((x, y), color)
    # 輪廓
    if outline:
        for i in range(len(coords)):
            x1, y1 = coords[i]
            x2, y2 = coords[(i + 1) % len(coords)]
            steps = max(abs(x2 - x1), abs(y2 - y1)) + 1
            for s in range(steps):
                t = s / max(steps - 1, 1)
                x = int(x1 + t * (x2 - x1))
                y = int(y1 + t * (y2 - y1))
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    img.putpixel((x, y), outline)

def draw_ray(img, cx, cy, angle_deg, length, width, color):
    """繪製光芒射線"""
    angle = math.radians(angle_deg)
    for d in range(length):
        x = int(cx + d * math.cos(angle))
        y = int(cy + d * math.sin(angle))
        for w in range(-width // 2, width // 2 + 1):
            px_x = x + int(w * math.sin(angle))
            px_y = y - int(w * math.cos(angle))
            if 0 <= px_x < SIZE and 0 <= px_y < SIZE:
                img.putpixel((px_x, px_y), color)

# 顏色定義
OUTLINE      = (20,  20,  20,  255)   # 深黑輪廓
GOLD_LIGHT   = (255, 240, 100, 255)   # 金色亮
GOLD_MID     = (255, 215,   0, 255)   # 金色中（#FFD700）
GOLD_DARK    = (200, 160,   0, 255)   # 金色暗
SKY_BLUE     = (0,   191, 255, 255)   # 天藍（#00BFFF）
SKY_BLUE_L   = (100, 220, 255, 255)   # 天藍亮
SKY_BLUE_D   = (0,   140, 200, 255)   # 天藍暗
PINK         = (255, 105, 180, 255)   # 粉紅（#FF69B4）
LIME         = (127, 255,   0, 255)   # 草綠（#7FFF00）
WHITE        = (255, 255, 255, 255)   # 白色
ORANGE       = (255, 165,   0, 255)   # 橙色
TRANSPARENT  = (0,   0,   0,   0)

img = Image.new("RGBA", (SIZE, SIZE), TRANSPARENT)

# ── 1. 背景光暈（天藍色圓形光暈）────────────────────────────────────────────
for r in range(28, 20, -1):
    alpha = int(60 * (28 - r) / 8)
    for y in range(SIZE):
        for x in range(SIZE):
            if (x - 32) ** 2 + (y - 32) ** 2 <= r * r:
                cur = img.getpixel((x, y))
                if cur[3] == 0:
                    img.putpixel((x, y), (0, 191, 255, alpha))

# ── 2. 主體：五角星魚身（金色）──────────────────────────────────────────────
draw_star(img, 32, 32, 20, 9, 5, GOLD_MID, OUTLINE)

# 星形陰影（右下暗）
for y in range(SIZE):
    for x in range(SIZE):
        c = img.getpixel((x, y))
        if c == GOLD_MID:
            if x > 34 and y > 34:
                img.putpixel((x, y), GOLD_DARK)
            elif x < 30 and y < 30:
                img.putpixel((x, y), GOLD_LIGHT)

# ── 3. 爆炸光芒（8方向射線）────────────────────────────────────────────────
ray_colors = [GOLD_LIGHT, SKY_BLUE, PINK, LIME, GOLD_LIGHT, SKY_BLUE, ORANGE, WHITE]
for i, angle in enumerate(range(0, 360, 45)):
    draw_ray(img, 32, 32, angle, 28, 1, ray_colors[i % len(ray_colors)])

# ── 4. 星爆點（5個天藍小星）────────────────────────────────────────────────
burst_positions = [
    (14, 14), (50, 14), (14, 50), (50, 50), (32, 8)
]
for bx, by in burst_positions:
    draw_star(img, bx, by, 5, 2, 4, SKY_BLUE, OUTLINE)
    # 中心白點
    if 0 <= bx < SIZE and 0 <= by < SIZE:
        img.putpixel((bx, by), WHITE)

# ── 5. 金色光環（外圈）──────────────────────────────────────────────────────
for angle_deg in range(0, 360, 3):
    angle = math.radians(angle_deg)
    x = int(32 + 24 * math.cos(angle))
    y = int(32 + 24 * math.sin(angle))
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), GOLD_MID)

# ── 6. 中心高光（白色小圓）──────────────────────────────────────────────────
fill_circle(img, 29, 29, 3, WHITE)
fill_circle(img, 29, 29, 1, (255, 255, 255, 255))

# ── 7. 彩色星點散落（裝飾）──────────────────────────────────────────────────
star_dots = [
    (8,  32, PINK),
    (56, 32, LIME),
    (32, 56, ORANGE),
    (20, 48, SKY_BLUE_L),
    (44, 48, GOLD_LIGHT),
    (20, 16, LIME),
    (44, 16, PINK),
]
for sx, sy, sc in star_dots:
    if 0 <= sx < SIZE and 0 <= sy < SIZE:
        img.putpixel((sx, sy), sc)
    if 0 <= sx + 1 < SIZE and 0 <= sy < SIZE:
        img.putpixel((sx + 1, sy), sc)
    if 0 <= sx < SIZE and 0 <= sy + 1 < SIZE:
        img.putpixel((sx, sy + 1), sc)

# ── 8. 魚眼（天藍色）────────────────────────────────────────────────────────
fill_circle(img, 36, 28, 3, SKY_BLUE)
fill_circle(img, 36, 28, 2, SKY_BLUE_D)
img.putpixel((35, 27), WHITE)  # 高光

# ── 9. 魚尾（橙色扇形）──────────────────────────────────────────────────────
tail_pixels = [
    (52, 28), (53, 27), (54, 26), (55, 25),
    (52, 32), (53, 33), (54, 34), (55, 35),
    (52, 30), (53, 30), (54, 30),
]
for tx, ty in tail_pixels:
    if 0 <= tx < SIZE and 0 <= ty < SIZE:
        img.putpixel((tx, ty), ORANGE)

# 儲存
os.makedirs(os.path.dirname(OUT_PATH), exist_ok=True)
img.save(OUT_PATH)
print(f"T240 sprite saved: {OUT_PATH}")

# 驗證
loaded = Image.open(OUT_PATH)
pixels = loaded.load()
non_transparent = sum(1 for y in range(SIZE) for x in range(SIZE) if pixels[x, y][3] > 0)
total = SIZE * SIZE
print(f"Non-transparent pixels: {non_transparent}/{total} ({non_transparent*100//total}%)")
