"""
generate_targets_day323.py — DAY-323 T206-T210 精靈圖生成
T206 幸運 Fever Boost 魚（火橙色，Fever Boost™）
T207 幸運公會戰魚（金色，公會戰旗幟）
T208 幸運路徑魚（青藍色，路徑光軌）
T209 幸運連鎖電鰻魚（深紫色，紫粉電鰻）
T210 幸運終極奇蹟魚（純白色，新史上最高）
"""
import os
import math
import random
from PIL import Image, ImageDraw

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def fill_circle(draw, cx, cy, r, color):
    for y in range(max(0, cy - r), min(SIZE, cy + r + 1)):
        for x in range(max(0, cx - r), min(SIZE, cx + r + 1)):
            if (x - cx) ** 2 + (y - cy) ** 2 <= r * r:
                draw.point((x, y), fill=color)

def fill_circle_shaded(draw, cx, cy, r, base_color):
    br, bg, bb, ba = base_color
    for y in range(max(0, cy - r), min(SIZE, cy + r + 1)):
        for x in range(max(0, cx - r), min(SIZE, cx + r + 1)):
            dx, dy = x - cx, y - cy
            if dx * dx + dy * dy <= r * r:
                shade = 1.0 - 0.3 * (dx + dy) / (r * 2 + 1)
                shade = max(0.6, min(1.2, shade))
                draw.point((x, y), fill=(
                    min(255, int(br * shade)),
                    min(255, int(bg * shade)),
                    min(255, int(bb * shade)),
                    ba
                ))

def draw_rays(draw, cx, cy, n, r_inner, r_outer, color, width=2):
    for i in range(n):
        angle = 2 * math.pi * i / n
        x1 = int(cx + r_inner * math.cos(angle))
        y1 = int(cy + r_inner * math.sin(angle))
        x2 = int(cx + r_outer * math.cos(angle))
        y2 = int(cy + r_outer * math.sin(angle))
        draw.line([(x1, y1), (x2, y2)], fill=color, width=width)

def draw_ring(draw, cx, cy, r, color, width=2):
    draw.ellipse([cx - r, cy - r, cx + r, cy + r], outline=color, width=width)

def count_pixels(img):
    pixels = img.load()
    count = 0
    for y in range(SIZE):
        for x in range(SIZE):
            if pixels[x, y][3] > 0:
                count += 1
    return count

# ── T206 幸運 Fever Boost 魚 ──────────────────────────────────
def gen_t206():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 外層火焰光暈（低透明度）
    fill_circle(draw, cx, cy, 30, (255, 100, 0, 40))
    fill_circle(draw, cx, cy, 26, (255, 140, 0, 60))

    # 主體：火橙色橢圓魚身
    for y in range(SIZE):
        for x in range(SIZE):
            dx, dy = (x - cx) / 22, (y - cy) / 16
            if dx * dx + dy * dy <= 1.0:
                shade = 1.0 - 0.3 * (dx + dy) / 2
                shade = max(0.6, min(1.2, shade))
                r = min(255, int(255 * shade))
                g = min(255, int(100 * shade))
                b = 0
                draw.point((x, y), fill=(r, g, b, 255))

    # 輪廓
    draw.ellipse([cx - 22, cy - 16, cx + 22, cy + 16], outline=(200, 60, 0, 255), width=2)

    # Fever Boost™ 火焰光芒（12道）
    draw_rays(draw, cx, cy, 12, 18, 28, (255, 200, 0, 220), width=2)

    # 中心 "F" 符號（Fever）
    fill_circle(draw, cx, cy, 8, (255, 220, 0, 255))
    draw.text((cx - 4, cy - 6), "F", fill=(200, 50, 0, 255))

    # 火焰粒子（20個）
    rng = random.Random(206)
    for _ in range(20):
        px = rng.randint(cx - 25, cx + 25)
        py = rng.randint(cy - 20, cy + 20)
        pr = rng.randint(1, 3)
        fill_circle(draw, px, py, pr, (255, rng.randint(100, 200), 0, rng.randint(150, 220)))

    # 雙層外圈
    draw_ring(draw, cx, cy, 29, (255, 150, 0, 180), width=2)
    draw_ring(draw, cx, cy, 31, (255, 200, 0, 120), width=1)

    density = count_pixels(img) / (SIZE * SIZE) * 100
    print(f"T206 Fever Boost 魚：{density:.1f}%")
    return img

# ── T207 幸運公會戰魚 ─────────────────────────────────────────
def gen_t207():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 外層金色光暈
    fill_circle(draw, cx, cy, 30, (255, 215, 0, 40))
    fill_circle(draw, cx, cy, 26, (255, 215, 0, 60))

    # 主體：金色橢圓魚身
    for y in range(SIZE):
        for x in range(SIZE):
            dx, dy = (x - cx) / 22, (y - cy) / 16
            if dx * dx + dy * dy <= 1.0:
                shade = 1.0 - 0.3 * (dx + dy) / 2
                shade = max(0.6, min(1.2, shade))
                r = min(255, int(255 * shade))
                g = min(255, int(215 * shade))
                b = 0
                draw.point((x, y), fill=(r, g, b, 255))

    # 輪廓
    draw.ellipse([cx - 22, cy - 16, cx + 22, cy + 16], outline=(200, 160, 0, 255), width=2)

    # 公會戰旗幟（劍形）
    # 劍身
    draw.line([(cx, cy - 14), (cx, cy + 14)], fill=(200, 200, 200, 255), width=3)
    # 劍柄
    draw.line([(cx - 8, cy + 8), (cx + 8, cy + 8)], fill=(150, 100, 0, 255), width=3)
    # 劍尖
    draw.polygon([(cx, cy - 14), (cx - 3, cy - 8), (cx + 3, cy - 8)], fill=(220, 220, 220, 255))

    # 16道金色光芒
    draw_rays(draw, cx, cy, 16, 18, 28, (255, 215, 0, 200), width=2)

    # 三層光環
    draw_ring(draw, cx, cy, 27, (255, 215, 0, 180), width=2)
    draw_ring(draw, cx, cy, 29, (255, 180, 0, 120), width=1)
    draw_ring(draw, cx, cy, 31, (255, 215, 0, 80), width=1)

    density = count_pixels(img) / (SIZE * SIZE) * 100
    print(f"T207 公會戰魚：{density:.1f}%")
    return img

# ── T208 幸運路徑魚 ───────────────────────────────────────────
def gen_t208():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 外層青藍光暈
    fill_circle(draw, cx, cy, 30, (0, 255, 255, 40))
    fill_circle(draw, cx, cy, 26, (0, 200, 255, 60))

    # 主體：青藍色橢圓魚身
    for y in range(SIZE):
        for x in range(SIZE):
            dx, dy = (x - cx) / 22, (y - cy) / 16
            if dx * dx + dy * dy <= 1.0:
                shade = 1.0 - 0.3 * (dx + dy) / 2
                shade = max(0.6, min(1.2, shade))
                r = 0
                g = min(255, int(200 * shade))
                b = min(255, int(255 * shade))
                draw.point((x, y), fill=(r, g, b, 255))

    # 輪廓
    draw.ellipse([cx - 22, cy - 16, cx + 22, cy + 16], outline=(0, 150, 200, 255), width=2)

    # 路徑光軌（S形曲線）
    for i in range(20):
        t = i / 19.0
        px = int(cx - 18 + 36 * t)
        py = int(cy + 8 * math.sin(t * math.pi * 2))
        fill_circle(draw, px, py, 2, (0, 255, 255, 200))

    # 路徑節點（5個）
    for i in range(5):
        t = i / 4.0
        px = int(cx - 18 + 36 * t)
        py = int(cy + 8 * math.sin(t * math.pi * 2))
        fill_circle(draw, px, py, 3, (255, 255, 0, 255))

    # 20道青藍光芒
    draw_rays(draw, cx, cy, 20, 18, 28, (0, 255, 255, 180), width=1)

    # 雙層光環
    draw_ring(draw, cx, cy, 28, (0, 255, 255, 160), width=2)
    draw_ring(draw, cx, cy, 30, (0, 200, 255, 100), width=1)

    density = count_pixels(img) / (SIZE * SIZE) * 100
    print(f"T208 路徑魚：{density:.1f}%")
    return img

# ── T209 幸運連鎖電鰻魚 ───────────────────────────────────────
def gen_t209():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 外層深紫光暈
    fill_circle(draw, cx, cy, 30, (150, 0, 255, 40))
    fill_circle(draw, cx, cy, 26, (200, 0, 255, 60))

    # 主體：深紫色橢圓魚身（電鰻形狀，細長）
    for y in range(SIZE):
        for x in range(SIZE):
            dx, dy = (x - cx) / 24, (y - cy) / 12
            if dx * dx + dy * dy <= 1.0:
                shade = 1.0 - 0.3 * (dx + dy) / 2
                shade = max(0.6, min(1.2, shade))
                r = min(255, int(150 * shade))
                g = 0
                b = min(255, int(255 * shade))
                draw.point((x, y), fill=(r, g, b, 255))

    # 輪廓
    draw.ellipse([cx - 24, cy - 12, cx + 24, cy + 12], outline=(100, 0, 200, 255), width=2)

    # 連鎖電擊（8條閃電）
    rng = random.Random(209)
    for i in range(8):
        angle = 2 * math.pi * i / 8
        x1 = int(cx + 14 * math.cos(angle))
        y1 = int(cy + 10 * math.sin(angle))
        x2 = int(cx + 28 * math.cos(angle))
        y2 = int(cy + 20 * math.sin(angle))
        # 鋸齒閃電
        mid_x = (x1 + x2) // 2 + rng.randint(-4, 4)
        mid_y = (y1 + y2) // 2 + rng.randint(-4, 4)
        draw.line([(x1, y1), (mid_x, mid_y)], fill=(255, 100, 255, 220), width=2)
        draw.line([(mid_x, mid_y), (x2, y2)], fill=(255, 100, 255, 220), width=2)

    # 中心電弧
    fill_circle(draw, cx, cy, 6, (255, 200, 255, 255))
    fill_circle(draw, cx, cy, 4, (255, 255, 255, 255))

    # 電弧紋路
    for i in range(6):
        angle = 2 * math.pi * i / 6
        x1 = int(cx + 6 * math.cos(angle))
        y1 = int(cy + 6 * math.sin(angle))
        x2 = int(cx + 12 * math.cos(angle))
        y2 = int(cy + 12 * math.sin(angle))
        draw.line([(x1, y1), (x2, y2)], fill=(200, 0, 255, 200), width=1)

    # 雙層光環
    draw_ring(draw, cx, cy, 27, (200, 0, 255, 180), width=2)
    draw_ring(draw, cx, cy, 29, (255, 100, 255, 120), width=1)

    density = count_pixels(img) / (SIZE * SIZE) * 100
    print(f"T209 連鎖電鰻魚：{density:.1f}%")
    return img

# ── T210 幸運終極奇蹟魚 ───────────────────────────────────────
def gen_t210():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 外層純白光暈（最強）
    fill_circle(draw, cx, cy, 31, (255, 255, 255, 50))
    fill_circle(draw, cx, cy, 28, (255, 255, 255, 80))
    fill_circle(draw, cx, cy, 25, (255, 255, 255, 100))

    # 主體：純白大型魚身
    for y in range(SIZE):
        for x in range(SIZE):
            dx, dy = (x - cx) / 22, (y - cy) / 18
            if dx * dx + dy * dy <= 1.0:
                shade = 1.0 - 0.2 * (dx + dy) / 2
                shade = max(0.7, min(1.1, shade))
                v = min(255, int(255 * shade))
                draw.point((x, y), fill=(v, v, v, 255))

    # 輪廓（金色）
    draw.ellipse([cx - 22, cy - 18, cx + 22, cy + 18], outline=(255, 215, 0, 255), width=2)

    # 35道光芒（最多）
    draw_rays(draw, cx, cy, 35, 18, 30, (255, 255, 255, 200), width=1)
    draw_rays(draw, cx, cy, 12, 18, 30, (255, 215, 0, 220), width=2)

    # 五層光環（最多層）
    for r, alpha in [(22, 200), (24, 160), (26, 120), (28, 80), (30, 50)]:
        draw_ring(draw, cx, cy, r, (255, 255, 255, alpha), width=1)
    draw_ring(draw, cx, cy, 23, (255, 215, 0, 180), width=2)

    # 奇蹟符文（星形）
    for i in range(5):
        angle = 2 * math.pi * i / 5 - math.pi / 2
        px = int(cx + 10 * math.cos(angle))
        py = int(cy + 10 * math.sin(angle))
        fill_circle(draw, px, py, 2, (255, 215, 0, 255))

    # 中心高光
    fill_circle(draw, cx, cy, 5, (255, 255, 255, 255))
    fill_circle(draw, cx, cy, 3, (255, 215, 0, 255))

    # 散落光點（30個）
    rng = random.Random(210)
    for _ in range(30):
        px = rng.randint(cx - 28, cx + 28)
        py = rng.randint(cy - 28, cy + 28)
        pr = rng.randint(1, 2)
        fill_circle(draw, px, py, pr, (255, 255, rng.randint(200, 255), rng.randint(180, 255)))

    density = count_pixels(img) / (SIZE * SIZE) * 100
    print(f"T210 終極奇蹟魚：{density:.1f}%")
    return img

# ── 主程式 ────────────────────────────────────────────────────
if __name__ == "__main__":
    os.makedirs(OUT_DIR, exist_ok=True)

    targets = [
        ("T206_fever_boost", gen_t206),
        ("T207_guild_battle", gen_t207),
        ("T208_path_fish", gen_t208),
        ("T209_chain_eel", gen_t209),
        ("T210_ultimate_miracle", gen_t210),
    ]

    for name, gen_func in targets:
        img = gen_func()
        path = os.path.join(OUT_DIR, f"{name}.png")
        img.save(path)
        print(f"  → 已儲存：{path}")

    print("\n✅ DAY-323 T206-T210 精靈圖生成完成！")
