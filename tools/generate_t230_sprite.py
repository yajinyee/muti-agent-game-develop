"""
generate_t230_sprite.py — T230 幸運品質突變魚精靈圖生成（DAY-272）
主題：彩虹品質突變魚（五色品質光環 + 魚身 + 突變星芒 + 品質等級標記）
"""
from PIL import Image, ImageDraw
import math
import os

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
OUT_FILE = os.path.join(OUT_DIR, "T230_quality_mutation.png")
SIZE = 64

def px(img, x, y, color):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), color)

def fill_circle(img, cx, cy, r, color):
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                px(img, cx+dx, cy+dy, color)

def fill_ellipse(img, cx, cy, rx, ry, color):
    for dy in range(-ry, ry+1):
        for dx in range(-rx, rx+1):
            if (dx/rx)**2 + (dy/ry)**2 <= 1.0:
                px(img, cx+dx, cy+dy, color)

def draw_star(img, cx, cy, r, color, points=6):
    """繪製星形光芒"""
    for i in range(points):
        angle = math.pi * 2 * i / points
        for t in range(r+1):
            x = int(cx + t * math.cos(angle))
            y = int(cy + t * math.sin(angle))
            px(img, x, y, color)

def hsv_to_rgb(h, s, v):
    """HSV 轉 RGB（0-255）"""
    h = h % 1.0
    i = int(h * 6)
    f = h * 6 - i
    p = v * (1 - s)
    q = v * (1 - f * s)
    t = v * (1 - (1 - f) * s)
    if i == 0: r, g, b = v, t, p
    elif i == 1: r, g, b = q, v, p
    elif i == 2: r, g, b = p, v, t
    elif i == 3: r, g, b = p, q, v
    elif i == 4: r, g, b = t, p, v
    else: r, g, b = v, p, q
    return (int(r*255), int(g*255), int(b*255), 255)

# 顏色定義
OUTLINE     = (20, 20, 30, 255)
BODY_BASE   = (255, 255, 255, 255)   # 白色魚身（突變後閃光）
BODY_LIGHT  = (240, 240, 255, 255)   # 淡藍白高光
BODY_DARK   = (180, 180, 200, 255)   # 陰影
FIN_COLOR   = (200, 200, 220, 255)   # 魚鰭
EYE_WHITE   = (255, 255, 255, 255)
EYE_PUPIL   = (30, 30, 50, 255)
EYE_SHINE   = (255, 255, 255, 255)
GOLD        = (255, 215, 0, 255)
MYTHIC_PINK = (255, 105, 180, 255)

img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

# ── 彩虹光環（外圈）────────────────────────────────────────────────────────────
# 繪製彩虹色外圈（品質突變的核心視覺）
cx, cy = 32, 32
for angle_deg in range(0, 360, 3):
    angle = math.radians(angle_deg)
    hue = angle_deg / 360.0
    color = hsv_to_rgb(hue, 0.9, 1.0)
    # 外圈（r=28-30）
    for r in [27, 28, 29]:
        x = int(cx + r * math.cos(angle))
        y = int(cy + r * math.sin(angle))
        px(img, x, y, color)

# ── 品質等級標記（五個小圓點，代表五個品質等級）──────────────────────────────
quality_colors = [
    (170, 170, 170, 255),  # Normal 灰
    (74, 144, 217, 255),   # Rare 藍
    (155, 89, 182, 255),   # Epic 紫
    (255, 140, 0, 255),    # Legendary 橙
    (255, 105, 180, 255),  # Mythic 粉
]
for i, qc in enumerate(quality_colors):
    angle = math.radians(-90 + i * 72)  # 五個點均勻分布
    dot_x = int(cx + 22 * math.cos(angle))
    dot_y = int(cy + 22 * math.sin(angle))
    fill_circle(img, dot_x, dot_y, 3, qc)
    # 輪廓
    for ddx in [-4, 4, 0, 0]:
        for ddy in [0, 0, -4, 4]:
            if ddx != 0 or ddy != 0:
                px(img, dot_x+ddx//2, dot_y+ddy//2, OUTLINE)

# ── 魚身（橢圓形，帶陰影）────────────────────────────────────────────────────
# 主體橢圓（rx=14, ry=10）
for dy in range(-10, 11):
    for dx in range(-14, 15):
        if (dx/14)**2 + (dy/10)**2 <= 1.0:
            # 3色陰影
            if dx < -4 and dy < -2:
                color = BODY_LIGHT
            elif dx > 4 and dy > 2:
                color = BODY_DARK
            else:
                color = BODY_BASE
            px(img, cx+dx, cy+dy, color)

# 魚身輪廓
for dy in range(-10, 11):
    for dx in range(-14, 15):
        dist = (dx/14)**2 + (dy/10)**2
        if 0.85 <= dist <= 1.0:
            px(img, cx+dx, cy+dy, OUTLINE)

# ── 魚尾（三角形）────────────────────────────────────────────────────────────
for dy in range(-7, 8):
    tail_len = int(8 * (1 - abs(dy) / 8.0))
    for dx in range(14, 14 + tail_len):
        px(img, cx+dx, cy+dy, FIN_COLOR)
# 魚尾輪廓
for dy in [-7, 7]:
    for dx in range(14, 22):
        px(img, cx+dx, cy+dy, OUTLINE)
for dy in range(-7, 8):
    px(img, cx+14+int(8*(1-abs(dy)/8.0)), cy+dy, OUTLINE)

# ── 背鰭 ──────────────────────────────────────────────────────────────────────
for dx in range(-6, 5):
    fin_h = max(0, 5 - abs(dx+1))
    for dy in range(-10 - fin_h, -10):
        px(img, cx+dx, cy+dy, FIN_COLOR)
# 背鰭輪廓
for dx in range(-6, 5):
    fin_h = max(0, 5 - abs(dx+1))
    if fin_h > 0:
        px(img, cx+dx, cy-10-fin_h, OUTLINE)

# ── 眼睛 ──────────────────────────────────────────────────────────────────────
eye_x, eye_y = cx - 8, cy - 2
fill_circle(img, eye_x, eye_y, 4, EYE_WHITE)
fill_circle(img, eye_x, eye_y, 2, EYE_PUPIL)
px(img, eye_x - 1, eye_y - 1, EYE_SHINE)
# 眼睛輪廓
for angle_deg in range(0, 360, 30):
    angle = math.radians(angle_deg)
    ex = int(eye_x + 4 * math.cos(angle))
    ey = int(eye_y + 4 * math.sin(angle))
    px(img, ex, ey, OUTLINE)

# ── 突變星芒（中心）──────────────────────────────────────────────────────────
# 中心金色星形
draw_star(img, cx, cy, 8, GOLD, points=8)
# 中心白色核心
fill_circle(img, cx, cy, 3, (255, 255, 255, 255))
# 中心輪廓
fill_circle(img, cx, cy, 4, GOLD)
fill_circle(img, cx, cy, 3, (255, 255, 255, 255))

# ── 彩虹光點（散落在魚身上）──────────────────────────────────────────────────
sparkle_positions = [
    (cx-5, cy-5), (cx+5, cy-4), (cx-3, cy+4),
    (cx+7, cy+2), (cx-8, cy+1), (cx+3, cy-7),
]
for i, (sx, sy) in enumerate(sparkle_positions):
    hue = i / len(sparkle_positions)
    sc = hsv_to_rgb(hue, 0.8, 1.0)
    px(img, sx, sy, sc)
    px(img, sx+1, sy, sc)
    px(img, sx, sy+1, sc)

# ── 放大到 64x64（已是 64x64，確認輸出）────────────────────────────────────
os.makedirs(OUT_DIR, exist_ok=True)
img.save(OUT_FILE)
print(f"✅ T230 精靈圖已生成：{OUT_FILE}")

# 統計非透明像素
non_transparent = sum(1 for x in range(SIZE) for y in range(SIZE) if img.getpixel((x, y))[3] > 0)
print(f"   非透明像素：{non_transparent} ({non_transparent/(SIZE*SIZE)*100:.1f}%)")
