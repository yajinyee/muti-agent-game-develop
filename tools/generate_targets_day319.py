"""
generate_targets_day319.py — DAY-319 T201-T205 精靈圖生成
T201 幸運能量風暴魚：青藍色魚身 + 5道閃電光芒 + 電弧紋路
T202 幸運水晶共鳴魚：水晶白藍色魚身 + 六角水晶結構 + 共鳴光環
T203 幸運命運審判魚：金色魚身 + 天秤符號 + 5個命運星
T204 幸運時間逆流魚：深藍紫色魚身 + 逆時針箭頭 + 時間漩渦
T205 幸運宇宙奇點魚：洋紅色大型魚身 + 30道光芒 + 五層光環 + 奇點符文（史上最高）
"""
from PIL import Image, ImageDraw
import math
import os
import random

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def fill_circle(draw, cx, cy, r, color):
    for y in range(max(0, cy-r), min(SIZE, cy+r+1)):
        for x in range(max(0, cx-r), min(SIZE, cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 <= r*r:
                draw.point((x, y), fill=color)

def fill_ellipse_shaded(draw, cx, cy, rx, ry, base_color, light_color, dark_color):
    for y in range(max(0, cy-ry), min(SIZE, cy+ry+1)):
        for x in range(max(0, cx-rx), min(SIZE, cx+rx+1)):
            if ((x-cx)/rx)**2 + ((y-cy)/ry)**2 <= 1.0:
                dx = (x - cx) / rx
                dy = (y - cy) / ry
                if dx < -0.2 and dy < -0.2:
                    color = light_color
                elif dx > 0.2 and dy > 0.2:
                    color = dark_color
                else:
                    color = base_color
                draw.point((x, y), fill=color)

def draw_ray(draw, cx, cy, angle_deg, length, color, width=1):
    angle = math.radians(angle_deg)
    ex = int(cx + math.cos(angle) * length)
    ey = int(cy + math.sin(angle) * length)
    draw.line([(cx, cy), (ex, ey)], fill=color, width=width)

def count_density(img):
    pixels = img.load()
    total = SIZE * SIZE
    non_transparent = sum(1 for y in range(SIZE) for x in range(SIZE) if pixels[x, y][3] > 10)
    return non_transparent / total * 100

# ── T201 幸運能量風暴魚 ──────────────────────────────────────
def gen_t201():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 外層電弧光暈
    for r in range(28, 31):
        for angle in range(0, 360, 3):
            a = math.radians(angle)
            x = int(cx + math.cos(a) * r)
            y = int(cy + math.sin(a) * r)
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=(0, 200, 255, 60))

    # 青藍色橢圓魚身（帶陰影）
    fill_ellipse_shaded(draw, cx, cy, 22, 16,
        (0, 180, 220, 255),
        (100, 240, 255, 255),
        (0, 100, 160, 255))

    # 5道閃電光芒（每72度一道）
    for i in range(5):
        angle = i * 72 - 90
        draw_ray(draw, cx, cy, angle, 26, (255, 255, 0, 255), 2)
        draw_ray(draw, cx, cy, angle + 5, 20, (200, 255, 255, 200), 1)

    # 電弧紋路（Z字形）
    for i in range(3):
        x_start = cx - 8 + i * 8
        draw.line([(x_start, cy-6), (x_start+4, cy), (x_start, cy+6)], fill=(255, 255, 100, 220), width=1)

    # 中心電核
    fill_circle(draw, cx, cy, 5, (255, 255, 255, 255))
    fill_circle(draw, cx, cy, 3, (0, 220, 255, 255))

    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 10:
                for dx, dy in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x+dx, y+dy
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] < 10:
                        draw.point((nx, ny), fill=(0, 0, 0, 200))

    path = os.path.join(OUT_DIR, "T201_energy_storm.png")
    img.save(path)
    density = count_density(img)
    print(f"T201 能量風暴魚: {density:.1f}% 非透明像素 → {path}")

# ── T202 幸運水晶共鳴魚 ──────────────────────────────────────
def gen_t202():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 外層水晶光暈
    for r in range(27, 30):
        for angle in range(0, 360, 4):
            a = math.radians(angle)
            x = int(cx + math.cos(a) * r)
            y = int(cy + math.sin(a) * r)
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=(200, 200, 255, 80))

    # 水晶白藍色橢圓魚身
    fill_ellipse_shaded(draw, cx, cy, 21, 15,
        (180, 180, 255, 255),
        (240, 240, 255, 255),
        (100, 100, 200, 255))

    # 六角水晶結構
    for i in range(6):
        angle = i * 60
        draw_ray(draw, cx, cy, angle, 18, (220, 220, 255, 200), 2)
        # 六角頂點
        a = math.radians(angle)
        px = int(cx + math.cos(a) * 18)
        py = int(cy + math.sin(a) * 18)
        fill_circle(draw, px, py, 2, (255, 255, 255, 255))

    # 共鳴光環（三層）
    for r in [8, 13, 18]:
        for angle in range(0, 360, 6):
            a = math.radians(angle)
            x = int(cx + math.cos(a) * r)
            y = int(cy + math.sin(a) * r)
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=(255, 255, 255, 150))

    # 中心水晶核
    fill_circle(draw, cx, cy, 5, (255, 255, 255, 255))
    fill_circle(draw, cx, cy, 3, (200, 200, 255, 255))

    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 10:
                for dx, dy in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x+dx, y+dy
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] < 10:
                        draw.point((nx, ny), fill=(0, 0, 0, 200))

    path = os.path.join(OUT_DIR, "T202_crystal_resonance.png")
    img.save(path)
    density = count_density(img)
    print(f"T202 水晶共鳴魚: {density:.1f}% 非透明像素 → {path}")

# ── T203 幸運命運審判魚 ──────────────────────────────────────
def gen_t203():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 外層金色光暈
    for r in range(27, 30):
        for angle in range(0, 360, 3):
            a = math.radians(angle)
            x = int(cx + math.cos(a) * r)
            y = int(cy + math.sin(a) * r)
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=(255, 200, 0, 80))

    # 金色橢圓魚身
    fill_ellipse_shaded(draw, cx, cy, 21, 15,
        (220, 180, 0, 255),
        (255, 230, 50, 255),
        (160, 120, 0, 255))

    # 天秤符號（橫桿 + 兩個盤）
    draw.line([(cx-14, cy-2), (cx+14, cy-2)], fill=(255, 255, 200, 255), width=2)
    draw.line([(cx, cy-2), (cx, cy+8)], fill=(255, 255, 200, 255), width=2)
    fill_circle(draw, cx-12, cy+2, 4, (255, 220, 50, 200))
    fill_circle(draw, cx+12, cy+2, 4, (255, 220, 50, 200))

    # 5個命運星（五角星）
    for i in range(5):
        angle = i * 72 - 90
        a = math.radians(angle)
        sx = int(cx + math.cos(a) * 20)
        sy = int(cy + math.sin(a) * 20)
        fill_circle(draw, sx, sy, 3, (255, 255, 100, 255))

    # 12道光芒
    for i in range(12):
        draw_ray(draw, cx, cy, i * 30, 24, (255, 200, 0, 150), 1)

    # 中心金核
    fill_circle(draw, cx, cy, 4, (255, 255, 0, 255))

    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 10:
                for dx, dy in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x+dx, y+dy
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] < 10:
                        draw.point((nx, ny), fill=(0, 0, 0, 200))

    path = os.path.join(OUT_DIR, "T203_fate_judgment.png")
    img.save(path)
    density = count_density(img)
    print(f"T203 命運審判魚: {density:.1f}% 非透明像素 → {path}")

# ── T204 幸運時間逆流魚 ──────────────────────────────────────
def gen_t204():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 外層時間漩渦光暈
    for r in range(27, 30):
        for angle in range(0, 360, 4):
            a = math.radians(angle)
            x = int(cx + math.cos(a) * r)
            y = int(cy + math.sin(a) * r)
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=(100, 100, 220, 80))

    # 深藍紫色橢圓魚身
    fill_ellipse_shaded(draw, cx, cy, 21, 15,
        (80, 80, 200, 255),
        (140, 140, 255, 255),
        (40, 40, 140, 255))

    # 逆時針箭頭（時間逆流）
    for i in range(3):
        angle = i * 120 + 30
        a = math.radians(angle)
        sx = int(cx + math.cos(a) * 14)
        sy = int(cy + math.sin(a) * 14)
        ex = int(cx + math.cos(a + 0.8) * 14)
        ey = int(cy + math.sin(a + 0.8) * 14)
        draw.line([(sx, sy), (ex, ey)], fill=(200, 200, 255, 220), width=2)

    # 時間漩渦（螺旋線）
    for i in range(0, 360, 15):
        r = 6 + i / 60
        a = math.radians(i)
        x = int(cx + math.cos(a) * r)
        y = int(cy + math.sin(a) * r)
        if 0 <= x < SIZE and 0 <= y < SIZE:
            draw.point((x, y), fill=(180, 180, 255, 200))

    # 時鐘符號（圓 + 指針）
    for angle in range(0, 360, 30):
        a = math.radians(angle)
        x = int(cx + math.cos(a) * 10)
        y = int(cy + math.sin(a) * 10)
        if 0 <= x < SIZE and 0 <= y < SIZE:
            draw.point((x, y), fill=(255, 255, 255, 200))
    draw.line([(cx, cy), (cx, cy-8)], fill=(255, 255, 255, 220), width=1)
    draw.line([(cx, cy), (cx+6, cy)], fill=(255, 255, 255, 220), width=1)

    # 中心時間核
    fill_circle(draw, cx, cy, 4, (200, 200, 255, 255))
    fill_circle(draw, cx, cy, 2, (255, 255, 255, 255))

    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 10:
                for dx, dy in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x+dx, y+dy
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] < 10:
                        draw.point((nx, ny), fill=(0, 0, 0, 200))

    path = os.path.join(OUT_DIR, "T204_time_reversal.png")
    img.save(path)
    density = count_density(img)
    print(f"T204 時間逆流魚: {density:.1f}% 非透明像素 → {path}")

# ── T205 幸運宇宙奇點魚（史上最高）──────────────────────────
def gen_t205():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 五層光環（最高階）
    for r, alpha in [(30, 40), (27, 60), (24, 80), (21, 100), (18, 120)]:
        for angle in range(0, 360, 2):
            a = math.radians(angle)
            x = int(cx + math.cos(a) * r)
            y = int(cy + math.sin(a) * r)
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=(255, 0, 255, alpha))

    # 洋紅色大型橢圓魚身
    fill_ellipse_shaded(draw, cx, cy, 23, 17,
        (200, 0, 200, 255),
        (255, 100, 255, 255),
        (140, 0, 140, 255))

    # 30道光芒（史上最多）
    for i in range(30):
        angle = i * 12
        length = 28 if i % 3 == 0 else 22
        color = (255, 255, 255, 200) if i % 3 == 0 else (255, 100, 255, 150)
        draw_ray(draw, cx, cy, angle, length, color, 1)

    # 奇點符文（∞ 無限符號）
    for t in range(0, 360, 5):
        a = math.radians(t)
        x = int(cx + 8 * math.cos(a) / (1 + math.sin(a)**2))
        y = int(cy + 8 * math.sin(a) * math.cos(a) / (1 + math.sin(a)**2))
        if 0 <= x < SIZE and 0 <= y < SIZE:
            draw.point((x, y), fill=(255, 255, 255, 220))

    # 中心奇點核（最亮）
    fill_circle(draw, cx, cy, 6, (255, 255, 255, 255))
    fill_circle(draw, cx, cy, 4, (255, 0, 255, 255))
    fill_circle(draw, cx, cy, 2, (255, 255, 255, 255))

    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 10:
                for dx, dy in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x+dx, y+dy
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] < 10:
                        draw.point((nx, ny), fill=(0, 0, 0, 200))

    path = os.path.join(OUT_DIR, "T205_cosmic_singularity.png")
    img.save(path)
    density = count_density(img)
    print(f"T205 宇宙奇點魚: {density:.1f}% 非透明像素 → {path}")

if __name__ == "__main__":
    os.makedirs(OUT_DIR, exist_ok=True)
    print("=== DAY-319 T201-T205 精靈圖生成 ===")
    gen_t201()
    gen_t202()
    gen_t203()
    gen_t204()
    gen_t205()
    print("=== 全部完成 ===")
