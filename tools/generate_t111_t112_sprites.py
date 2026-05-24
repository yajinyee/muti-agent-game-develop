"""
generate_t111_t112_sprites.py — T111 覺醒鳳凰 + T112 全場震盪 精靈圖生成
DAY-293 新增幸運特殊魚
"""
from PIL import Image, ImageDraw
import math
import os

OUTPUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 48

def clamp(v, lo, hi):
    return max(lo, min(hi, v))

def blend(c1, c2, t):
    return tuple(int(c1[i] * (1 - t) + c2[i] * t) for i in range(3))

def set_pixel(img, x, y, color):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), color)

def draw_circle_outline(img, cx, cy, r, color, thickness=1):
    for angle in range(0, 360, 2):
        rad = math.radians(angle)
        for dr in range(thickness):
            px = int(cx + (r + dr) * math.cos(rad))
            py = int(cy + (r + dr) * math.sin(rad))
            set_pixel(img, px, py, color)

def draw_filled_circle(img, cx, cy, r, color):
    for dy in range(-r, r + 1):
        for dx in range(-r, r + 1):
            if dx * dx + dy * dy <= r * r:
                set_pixel(img, cx + dx, cy + dy, color)

def draw_ellipse_filled(img, cx, cy, rx, ry, color):
    for dy in range(-ry, ry + 1):
        for dx in range(-rx, rx + 1):
            if (dx / rx) ** 2 + (dy / ry) ** 2 <= 1.0:
                set_pixel(img, cx + dx, cy + dy, color)

def draw_ray(img, cx, cy, angle_deg, length, color, width=1):
    rad = math.radians(angle_deg)
    for i in range(length):
        px = int(cx + i * math.cos(rad))
        py = int(cy + i * math.sin(rad))
        for w in range(-width // 2, width // 2 + 1):
            set_pixel(img, px + w, py, color)
            set_pixel(img, px, py + w, color)

# ─────────────────────────────────────────────────────────────
# T111 覺醒鳳凰魚
# 設計：火橙漸層橢圓魚身 + 鳳凰羽翼光芒 + 覺醒光環 + 金色眼睛 + 火焰尾羽
# ─────────────────────────────────────────────────────────────
def generate_t111():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    cx, cy = SIZE // 2, SIZE // 2

    # 主體橢圓（火橙漸層）
    for dy in range(-14, 15):
        for dx in range(-18, 19):
            if (dx / 18) ** 2 + (dy / 14) ** 2 <= 1.0:
                dist = math.sqrt((dx / 18) ** 2 + (dy / 14) ** 2)
                # 火橙漸層：中心亮橙，邊緣深橙
                r = int(255 * (1.0 - dist * 0.3))
                g = int(clamp(107 - dist * 60, 30, 107))
                b = int(clamp(53 - dist * 40, 0, 53))
                alpha = 255
                set_pixel(img, cx + dx, cy + dy, (r, g, b, alpha))

    # 鳳凰羽翼光芒（8 方向，金色）
    GOLD = (255, 215, 0, 200)
    FIRE_ORANGE = (255, 107, 53, 180)
    for angle in range(0, 360, 45):
        rad = math.radians(angle)
        for i in range(16, 22):
            px = int(cx + i * math.cos(rad))
            py = int(cy + i * math.sin(rad))
            set_pixel(img, px, py, GOLD)
        # 羽翼尖端
        for i in range(22, 26):
            px = int(cx + i * math.cos(rad))
            py = int(cy + i * math.sin(rad))
            set_pixel(img, px, py, FIRE_ORANGE)

    # 覺醒光環（雙層圓環）
    RING1 = (255, 200, 50, 160)
    RING2 = (255, 140, 0, 120)
    draw_circle_outline(img, cx, cy, 17, RING1, 1)
    draw_circle_outline(img, cx, cy, 19, RING2, 1)

    # 鳳凰尾羽（右側，火橙→金色漸層）
    for i in range(5):
        tail_x = cx + 16 + i * 2
        tail_y = cy + i - 2
        r = int(255)
        g = int(107 + i * 20)
        b = int(53 + i * 10)
        draw_filled_circle(img, tail_x, tail_y, 2, (r, g, b, 220))

    # 金色眼睛
    EYE_GOLD = (255, 220, 0, 255)
    EYE_DARK = (80, 40, 0, 255)
    draw_filled_circle(img, cx - 5, cy - 2, 3, EYE_GOLD)
    draw_filled_circle(img, cx - 5, cy - 2, 1, EYE_DARK)
    # 眼睛高光
    set_pixel(img, cx - 6, cy - 3, (255, 255, 255, 255))

    # 覺醒核心（中心白色光點）
    CORE = (255, 255, 220, 240)
    draw_filled_circle(img, cx, cy, 3, CORE)

    # 火焰粒子散落
    SPARK = (255, 180, 50, 180)
    sparks = [(-20, -8), (20, -8), (-20, 8), (20, 8), (0, -22), (0, 22)]
    for sx, sy in sparks:
        draw_filled_circle(img, cx + sx, cy + sy, 1, SPARK)

    # 輪廓（深橙）
    OUTLINE = (180, 60, 0, 255)
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 0:
                for nx, ny in [(x-1,y),(x+1,y),(x,y-1),(x,y+1)]:
                    if 0 <= nx < SIZE and 0 <= ny < SIZE:
                        if img.getpixel((nx, ny))[3] == 0:
                            img.putpixel((nx, ny), OUTLINE)

    path = os.path.join(OUTPUT_DIR, "T111_awakened_phoenix.png")
    img.save(path)
    # 統計非透明像素
    pixels = [img.getpixel((x, y)) for y in range(SIZE) for x in range(SIZE)]
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    total = SIZE * SIZE
    print(f"T111 覺醒鳳凰魚: {path}")
    print(f"  非透明像素: {non_transparent}/{total} ({non_transparent/total*100:.1f}%)")
    return img

# ─────────────────────────────────────────────────────────────
# T112 全場震盪魚
# 設計：深橙漸層圓形魚身 + 震盪波紋（同心圓） + 爆炸光芒 + 震盪核心
# ─────────────────────────────────────────────────────────────
def generate_t112():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    cx, cy = SIZE // 2, SIZE // 2

    # 主體圓形（深橙漸層）
    for dy in range(-15, 16):
        for dx in range(-15, 16):
            if dx * dx + dy * dy <= 15 * 15:
                dist = math.sqrt(dx * dx + dy * dy) / 15
                # 深橙漸層：中心亮橙，邊緣深紅橙
                r = int(255 * (1.0 - dist * 0.2))
                g = int(clamp(69 - dist * 50, 10, 69))
                b = int(clamp(0, 0, 20))
                alpha = 255
                set_pixel(img, cx + dx, cy + dy, (r, g, b, alpha))

    # 震盪波紋（3 層同心圓，橙→金→白）
    WAVE1 = (255, 140, 0, 180)
    WAVE2 = (255, 200, 50, 150)
    WAVE3 = (255, 240, 120, 120)
    draw_circle_outline(img, cx, cy, 17, WAVE1, 1)
    draw_circle_outline(img, cx, cy, 20, WAVE2, 1)
    draw_circle_outline(img, cx, cy, 22, WAVE3, 1)

    # 爆炸光芒（4 方向，粗）
    BLAST = (255, 100, 0, 200)
    for angle in [0, 90, 180, 270]:
        rad = math.radians(angle)
        for i in range(16, 23):
            px = int(cx + i * math.cos(rad))
            py = int(cy + i * math.sin(rad))
            set_pixel(img, px, py, BLAST)
            set_pixel(img, px + 1, py, BLAST)
            set_pixel(img, px, py + 1, BLAST)

    # 斜向光芒（細）
    DIAG = (255, 160, 50, 160)
    for angle in [45, 135, 225, 315]:
        rad = math.radians(angle)
        for i in range(16, 21):
            px = int(cx + i * math.cos(rad))
            py = int(cy + i * math.sin(rad))
            set_pixel(img, px, py, DIAG)

    # 震盪裂縫紋路（X 形）
    CRACK = (200, 80, 0, 200)
    for i in range(-8, 9):
        set_pixel(img, cx + i, cy + i, CRACK)
        set_pixel(img, cx + i, cy - i, CRACK)

    # 眼睛（橙色）
    EYE_ORANGE = (255, 180, 0, 255)
    EYE_DARK = (100, 30, 0, 255)
    draw_filled_circle(img, cx - 5, cy - 2, 3, EYE_ORANGE)
    draw_filled_circle(img, cx - 5, cy - 2, 1, EYE_DARK)
    set_pixel(img, cx - 6, cy - 3, (255, 255, 255, 255))

    # 震盪核心（中心白色爆炸點）
    CORE_WHITE = (255, 255, 240, 255)
    CORE_ORANGE = (255, 140, 0, 220)
    draw_filled_circle(img, cx, cy, 4, CORE_ORANGE)
    draw_filled_circle(img, cx, cy, 2, CORE_WHITE)

    # 爆炸粒子
    PARTICLE = (255, 200, 100, 180)
    particles = [(-18, -10), (18, -10), (-18, 10), (18, 10),
                 (-10, -18), (10, -18), (-10, 18), (10, 18)]
    for px, py in particles:
        draw_filled_circle(img, cx + px, cy + py, 1, PARTICLE)

    # 輪廓（深紅橙）
    OUTLINE = (160, 40, 0, 255)
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 0:
                for nx, ny in [(x-1,y),(x+1,y),(x,y-1),(x,y+1)]:
                    if 0 <= nx < SIZE and 0 <= ny < SIZE:
                        if img.getpixel((nx, ny))[3] == 0:
                            img.putpixel((nx, ny), OUTLINE)

    path = os.path.join(OUTPUT_DIR, "T112_shockwave_bomb.png")
    img.save(path)
    pixels = [img.getpixel((x, y)) for y in range(SIZE) for x in range(SIZE)]
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    total = SIZE * SIZE
    print(f"T112 全場震盪魚: {path}")
    print(f"  非透明像素: {non_transparent}/{total} ({non_transparent/total*100:.1f}%)")
    return img

if __name__ == "__main__":
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    print("=== DAY-293 T111-T112 精靈圖生成 ===")
    generate_t111()
    generate_t112()
    print("=== 完成 ===")
