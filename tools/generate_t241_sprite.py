"""
generate_t241_sprite.py — T241 幸運四象大獎魚精靈圖生成（DAY-283）
設計：四象主題魚（四象符文+青龍綠/白虎銀/朱雀紅/玄武藍四色分區+金色光環+大獎光芒）
輸出：64x64 PNG（NEAREST 插值，像素風格）
"""
from PIL import Image
import math
import os

OUT_PATH = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets\T241_four_symbols.png"
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

# 顏色定義
OUTLINE      = (20,  20,  20,  255)   # 深黑輪廓
GOLD_LIGHT   = (255, 240, 100, 255)   # 金色亮
GOLD_MID     = (255, 215,   0, 255)   # 金色中（#FFD700）
GOLD_DARK    = (200, 160,   0, 255)   # 金色暗
QINGLONG_L   = (100, 220, 100, 255)   # 青龍亮綠
QINGLONG_M   = (0,   170,   0, 255)   # 青龍中綠
QINGLONG_D   = (0,   100,   0, 255)   # 青龍暗綠
BAIHU_L      = (240, 240, 240, 255)   # 白虎亮銀
BAIHU_M      = (192, 192, 192, 255)   # 白虎中銀
BAIHU_D      = (140, 140, 140, 255)   # 白虎暗銀
ZHUQUE_L     = (255, 100, 100, 255)   # 朱雀亮紅
ZHUQUE_M     = (220,   0,   0, 255)   # 朱雀中紅
ZHUQUE_D     = (150,   0,   0, 255)   # 朱雀暗紅
XUANWU_L     = (100, 100, 200, 255)   # 玄武亮藍
XUANWU_M     = (0,    0, 128, 255)    # 玄武中藍（#000080）
XUANWU_D     = (0,    0,  80, 255)    # 玄武暗藍
WHITE        = (255, 255, 255, 255)   # 白色
TRANSPARENT  = (0,   0,   0,   0)

img = Image.new("RGBA", (SIZE, SIZE), TRANSPARENT)

# ── 1. 主體：圓形魚身（四象四色分區）────────────────────────────────────────
cx, cy, r = 32, 32, 20
for y in range(max(0, cy - r), min(SIZE, cy + r + 1)):
    for x in range(max(0, cx - r), min(SIZE, cx + r + 1)):
        if (x - cx) ** 2 + (y - cy) ** 2 <= r * r:
            # 四象分區：左上青龍、右上白虎、左下朱雀、右下玄武
            if x <= cx and y <= cy:
                # 青龍（左上）
                if x < cx - 1 and y < cy - 1:
                    img.putpixel((x, y), QINGLONG_L)
                elif x > cx - 3 or y > cy - 3:
                    img.putpixel((x, y), QINGLONG_D)
                else:
                    img.putpixel((x, y), QINGLONG_M)
            elif x > cx and y <= cy:
                # 白虎（右上）
                if x > cx + 1 and y < cy - 1:
                    img.putpixel((x, y), BAIHU_L)
                elif x < cx + 3 or y > cy - 3:
                    img.putpixel((x, y), BAIHU_D)
                else:
                    img.putpixel((x, y), BAIHU_M)
            elif x <= cx and y > cy:
                # 朱雀（左下）
                if x < cx - 1 and y > cy + 1:
                    img.putpixel((x, y), ZHUQUE_L)
                elif x > cx - 3 or y < cy + 3:
                    img.putpixel((x, y), ZHUQUE_D)
                else:
                    img.putpixel((x, y), ZHUQUE_M)
            else:
                # 玄武（右下）
                if x > cx + 1 and y > cy + 1:
                    img.putpixel((x, y), XUANWU_L)
                elif x < cx + 3 or y < cy + 3:
                    img.putpixel((x, y), XUANWU_D)
                else:
                    img.putpixel((x, y), XUANWU_M)

# ── 2. 輪廓（深黑）──────────────────────────────────────────────────────────
for angle_deg in range(0, 360, 2):
    angle = math.radians(angle_deg)
    x = int(cx + r * math.cos(angle))
    y = int(cy + r * math.sin(angle))
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), OUTLINE)

# ── 3. 四象分隔線（金色十字）────────────────────────────────────────────────
for i in range(-r, r + 1):
    # 水平線
    if 0 <= cx + i < SIZE and 0 <= cy < SIZE:
        if (cx + i - cx) ** 2 + (cy - cy) ** 2 <= r * r:
            img.putpixel((cx + i, cy), GOLD_MID)
    # 垂直線
    if 0 <= cx < SIZE and 0 <= cy + i < SIZE:
        if (cx - cx) ** 2 + (cy + i - cy) ** 2 <= r * r:
            img.putpixel((cx, cy + i), GOLD_MID)

# ── 4. 金色光環（外圈）──────────────────────────────────────────────────────
for angle_deg in range(0, 360, 3):
    angle = math.radians(angle_deg)
    x = int(cx + (r + 2) * math.cos(angle))
    y = int(cy + (r + 2) * math.sin(angle))
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), GOLD_MID)

# ── 5. 大獎光芒（4方向射線）────────────────────────────────────────────────
ray_colors = [QINGLONG_M, BAIHU_M, ZHUQUE_M, XUANWU_M]
for i, angle in enumerate([315, 45, 225, 135]):  # 左上/右上/左下/右下
    draw_ray(img, cx, cy, angle, 28, 2, ray_colors[i])

# ── 6. 四象符文（四個角落的小圓點）────────────────────────────────────────
symbol_dots = [
    (18, 18, QINGLONG_L),  # 青龍（左上）
    (46, 18, BAIHU_L),     # 白虎（右上）
    (18, 46, ZHUQUE_L),    # 朱雀（左下）
    (46, 46, XUANWU_L),    # 玄武（右下）
]
for sx, sy, sc in symbol_dots:
    fill_circle(img, sx, sy, 3, sc)
    if 0 <= sx < SIZE and 0 <= sy < SIZE:
        img.putpixel((sx, sy), WHITE)  # 中心高光

# ── 7. 中心金色圓（大獎核心）────────────────────────────────────────────────
fill_circle(img, cx, cy, 5, GOLD_MID)
fill_circle(img, cx, cy, 3, GOLD_LIGHT)
img.putpixel((cx - 1, cy - 1), WHITE)  # 高光

# ── 8. 魚眼（金色）──────────────────────────────────────────────────────────
fill_circle(img, 38, 26, 3, GOLD_MID)
fill_circle(img, 38, 26, 2, GOLD_DARK)
img.putpixel((37, 25), WHITE)  # 高光

# ── 9. 魚尾（金色扇形）──────────────────────────────────────────────────────
tail_pixels = [
    (52, 28), (53, 27), (54, 26), (55, 25),
    (52, 32), (53, 33), (54, 34), (55, 35),
    (52, 30), (53, 30), (54, 30),
]
for tx, ty in tail_pixels:
    if 0 <= tx < SIZE and 0 <= ty < SIZE:
        img.putpixel((tx, ty), GOLD_MID)

# ── 10. 額外裝飾：四象光點散落 ──────────────────────────────────────────────
extra_dots = [
    (8,  8,  QINGLONG_L), (9,  9,  QINGLONG_M),
    (55, 8,  BAIHU_L),    (54, 9,  BAIHU_M),
    (8,  55, ZHUQUE_L),   (9,  54, ZHUQUE_M),
    (55, 55, XUANWU_L),   (54, 54, XUANWU_M),
    (32, 4,  GOLD_LIGHT), (32, 5,  GOLD_MID),
    (4,  32, GOLD_LIGHT), (5,  32, GOLD_MID),
    (59, 32, GOLD_LIGHT), (58, 32, GOLD_MID),
    (32, 59, GOLD_LIGHT), (32, 58, GOLD_MID),
]
for dx, dy, dc in extra_dots:
    if 0 <= dx < SIZE and 0 <= dy < SIZE:
        img.putpixel((dx, dy), dc)

# ── 11. 四象邊框（外圈裝飾）────────────────────────────────────────────────
border_colors = [
    (QINGLONG_M, range(0, 16)),   # 左上角青龍
    (BAIHU_M,    range(48, 64)),  # 右上角白虎
    (ZHUQUE_M,   range(0, 16)),   # 左下角朱雀
    (XUANWU_M,   range(48, 64)),  # 右下角玄武
]
for y in range(0, 8):
    for x in range(0, 8):
        if 0 <= x < SIZE and 0 <= y < SIZE:
            img.putpixel((x, y), QINGLONG_D)
    for x in range(56, 64):
        if 0 <= x < SIZE and 0 <= y < SIZE:
            img.putpixel((x, y), BAIHU_D)
for y in range(56, 64):
    for x in range(0, 8):
        if 0 <= x < SIZE and 0 <= y < SIZE:
            img.putpixel((x, y), ZHUQUE_D)
    for x in range(56, 64):
        if 0 <= x < SIZE and 0 <= y < SIZE:
            img.putpixel((x, y), XUANWU_D)

# 儲存
os.makedirs(os.path.dirname(OUT_PATH), exist_ok=True)
img.save(OUT_PATH)
print(f"T241 sprite saved: {OUT_PATH}")

# 驗證
loaded = Image.open(OUT_PATH)
pixels = loaded.load()
non_transparent = sum(1 for y in range(SIZE) for x in range(SIZE) if pixels[x, y][3] > 0)
total = SIZE * SIZE
print(f"Non-transparent pixels: {non_transparent}/{total} ({non_transparent*100//total}%)")
