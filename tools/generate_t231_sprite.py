"""
generate_t231_sprite.py — T231 幸運共鳴波魚精靈圖生成（DAY-273）
主題：天藍共鳴波魚（同心圓波紋 + 魚身 + 波紋光環 + 天藍光點）
"""
from PIL import Image
import math
import os

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
OUT_FILE = os.path.join(OUT_DIR, "T231_resonance_wave.png")
SIZE = 64

def px(img, x, y, color):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), color)

def fill_circle(img, cx, cy, r, color):
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                px(img, cx+dx, cy+dy, color)

def draw_ring(img, cx, cy, r, color, width=1):
    """繪製圓環"""
    for angle_deg in range(0, 360, 2):
        angle = math.radians(angle_deg)
        for rr in range(r - width, r + 1):
            x = int(cx + rr * math.cos(angle))
            y = int(cy + rr * math.sin(angle))
            px(img, x, y, color)

def fill_ellipse(img, cx, cy, rx, ry, color):
    for dy in range(-ry, ry+1):
        for dx in range(-rx, rx+1):
            if (dx/rx)**2 + (dy/ry)**2 <= 1.0:
                px(img, cx+dx, cy+dy, color)

# 顏色定義
OUTLINE     = (10, 30, 50, 255)
BODY_BASE   = (200, 235, 255, 255)   # 淡天藍魚身
BODY_LIGHT  = (230, 248, 255, 255)   # 高光
BODY_DARK   = (120, 180, 220, 255)   # 陰影
FIN_COLOR   = (150, 210, 240, 255)   # 魚鰭
EYE_WHITE   = (255, 255, 255, 255)
EYE_PUPIL   = (0, 60, 120, 255)
EYE_SHINE   = (255, 255, 255, 255)
WAVE1       = (0, 191, 255, 200)     # 第 1 層波（天藍，半透明）
WAVE2       = (30, 144, 255, 150)    # 第 2 層波（道奇藍，半透明）
WAVE3       = (0, 100, 200, 100)     # 第 3 層波（深藍，半透明）
GOLD        = (255, 215, 0, 255)
SPARKLE     = (0, 255, 200, 255)     # 翠綠光點

img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

cx, cy = 32, 32

# ── 三層同心圓波紋（核心視覺）────────────────────────────────────────────────
# 第 3 層（最外圈，最淡）
draw_ring(img, cx, cy, 28, WAVE3, width=1)
# 第 2 層（中圈）
draw_ring(img, cx, cy, 20, WAVE2, width=2)
# 第 1 層（內圈，最亮）
draw_ring(img, cx, cy, 12, WAVE1, width=2)

# ── 魚身（橢圓形，帶陰影）────────────────────────────────────────────────────
for dy in range(-9, 10):
    for dx in range(-13, 14):
        if (dx/13)**2 + (dy/9)**2 <= 1.0:
            if dx < -3 and dy < -2:
                color = BODY_LIGHT
            elif dx > 3 and dy > 2:
                color = BODY_DARK
            else:
                color = BODY_BASE
            px(img, cx+dx, cy+dy, color)

# 魚身輪廓
for dy in range(-9, 10):
    for dx in range(-13, 14):
        dist = (dx/13)**2 + (dy/9)**2
        if 0.85 <= dist <= 1.0:
            px(img, cx+dx, cy+dy, OUTLINE)

# ── 魚尾 ──────────────────────────────────────────────────────────────────────
for dy in range(-6, 7):
    tail_len = int(7 * (1 - abs(dy) / 7.0))
    for dx in range(13, 13 + tail_len):
        px(img, cx+dx, cy+dy, FIN_COLOR)
for dy in [-6, 6]:
    for dx in range(13, 20):
        px(img, cx+dx, cy+dy, OUTLINE)

# ── 背鰭 ──────────────────────────────────────────────────────────────────────
for dx in range(-5, 4):
    fin_h = max(0, 4 - abs(dx+1))
    for dy in range(-9 - fin_h, -9):
        px(img, cx+dx, cy+dy, FIN_COLOR)

# ── 眼睛 ──────────────────────────────────────────────────────────────────────
eye_x, eye_y = cx - 7, cy - 2
fill_circle(img, eye_x, eye_y, 3, EYE_WHITE)
fill_circle(img, eye_x, eye_y, 2, EYE_PUPIL)
px(img, eye_x - 1, eye_y - 1, EYE_SHINE)
for angle_deg in range(0, 360, 45):
    angle = math.radians(angle_deg)
    ex = int(eye_x + 3 * math.cos(angle))
    ey = int(eye_y + 3 * math.sin(angle))
    px(img, ex, ey, OUTLINE)

# ── 波紋中心光點 ──────────────────────────────────────────────────────────────
fill_circle(img, cx, cy, 3, WAVE1)
fill_circle(img, cx, cy, 2, (255, 255, 255, 200))

# ── 翠綠光點（散落）──────────────────────────────────────────────────────────
sparkle_positions = [
    (cx-10, cy-8), (cx+8, cy-10), (cx-8, cy+8),
    (cx+10, cy+6), (cx-12, cy+2), (cx+4, cy-12),
    (cx+14, cy-4), (cx-4, cy+12),
]
for sx, sy in sparkle_positions:
    px(img, sx, sy, SPARKLE)
    px(img, sx+1, sy, SPARKLE)

# ── 波紋方向箭頭（表示擴散方向）──────────────────────────────────────────────
# 四個方向的小箭頭
for i, (ax, ay) in enumerate([(cx, cy-22), (cx, cy+22), (cx-22, cy), (cx+22, cy)]):
    angle = math.radians(i * 90)
    for t in range(4):
        x = int(ax + t * math.cos(angle + math.pi/2))
        y = int(ay + t * math.sin(angle + math.pi/2))
        px(img, x, y, WAVE1)
        x = int(ax - t * math.cos(angle + math.pi/2))
        y = int(ay - t * math.sin(angle + math.pi/2))
        px(img, x, y, WAVE1)

os.makedirs(OUT_DIR, exist_ok=True)
img.save(OUT_FILE)
print(f"✅ T231 精靈圖已生成：{OUT_FILE}")

non_transparent = sum(1 for x in range(SIZE) for y in range(SIZE) if img.getpixel((x, y))[3] > 0)
print(f"   非透明像素：{non_transparent} ({non_transparent/(SIZE*SIZE)*100:.1f}%)")
