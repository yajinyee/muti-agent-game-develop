"""
generate_targets_day312.py — DAY-312 T166-T170 精靈圖生成
T166 幸運星際門戶魚（深紫色 + 星際門戶光環 + 傳送粒子）
T167 幸運龍魂融合魚（龍紅色 + 龍魂光環 + 火焰紋路）
T168 幸運時空裂縫魚（深藍色 + 時空裂縫紋路 + 時鐘符號）
T169 幸運神聖審判魚（神聖橙金色 + 光柱 + 天秤符號）
T170 幸運宇宙大爆炸魚（深紅色 + 爆炸光芒 + 宇宙粒子）
"""
import os
import math
import random
from PIL import Image, ImageDraw

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64
os.makedirs(OUT_DIR, exist_ok=True)

rng = random.Random(312)

def fill_circle(draw, cx, cy, r, color):
    for y in range(max(0, cy-r), min(SIZE, cy+r+1)):
        for x in range(max(0, cx-r), min(SIZE, cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 <= r*r:
                draw.point((x, y), fill=color)

def fill_circle_shaded(img, cx, cy, r, base_color, light_color, dark_color):
    pixels = img.load()
    for y in range(max(0, cy-r), min(SIZE, cy+r+1)):
        for x in range(max(0, cx-r), min(SIZE, cx+r+1)):
            dx, dy = x-cx, y-cy
            if dx*dx + dy*dy <= r*r:
                if dx < -r//3 and dy < -r//3:
                    pixels[x, y] = light_color
                elif dx > r//3 and dy > r//3:
                    pixels[x, y] = dark_color
                else:
                    pixels[x, y] = base_color

def draw_ring(draw, cx, cy, r_outer, r_inner, color):
    for y in range(max(0, cy-r_outer), min(SIZE, cy+r_outer+1)):
        for x in range(max(0, cx-r_outer), min(SIZE, cx+r_outer+1)):
            d2 = (x-cx)**2 + (y-cy)**2
            if r_inner*r_inner <= d2 <= r_outer*r_outer:
                draw.point((x, y), fill=color)

def draw_ray(draw, cx, cy, angle_deg, length, width, color):
    angle = math.radians(angle_deg)
    for i in range(length):
        bx = cx + int(i * math.cos(angle))
        by = cy + int(i * math.sin(angle))
        for w in range(-width//2, width//2+1):
            px = bx + int(w * math.sin(angle))
            py = by - int(w * math.cos(angle))
            if 0 <= px < SIZE and 0 <= py < SIZE:
                draw.point((px, py), fill=color)

# ── T166 幸運星際門戶魚 ──────────────────────────────────────
def gen_t166():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2

    # 外層星際光暈（低透明度大橢圓）
    for r in range(30, 26, -1):
        alpha = int(40 + (30-r)*8)
        draw_ring(draw, cx, cy, r, r-1, (125, 28, 163, alpha))

    # 主體橢圓魚身（深紫色）
    fill_circle_shaded(img, cx, cy, 20,
        (125, 28, 163, 255),   # 深紫
        (186, 104, 200, 255),  # 亮紫
        (74, 0, 112, 255))     # 暗紫

    # 星際門戶光環（雙層）
    draw_ring(draw, cx, cy, 28, 26, (224, 64, 251, 180))
    draw_ring(draw, cx, cy, 24, 22, (186, 104, 200, 150))

    # 8 方向傳送粒子
    for i in range(8):
        angle = i * 45
        rad = math.radians(angle)
        for dist in range(22, 30):
            px = cx + int(dist * math.cos(rad))
            py = cy + int(dist * math.sin(rad))
            if 0 <= px < SIZE and 0 <= py < SIZE:
                draw.point((px, py), fill=(224, 64, 251, 200))

    # 中心星形符號（白色）
    for i in range(4):
        angle = i * 45
        draw_ray(draw, cx, cy, angle, 8, 1, (255, 255, 255, 220))

    # 眼睛（紫色發光眼）
    fill_circle(draw, cx-5, cy-3, 3, (224, 64, 251, 255))
    fill_circle(draw, cx+5, cy-3, 3, (224, 64, 251, 255))
    fill_circle(draw, cx-5, cy-3, 1, (255, 255, 255, 255))
    fill_circle(draw, cx+5, cy-3, 1, (255, 255, 255, 255))

    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 0:
                for dx2, dy2 in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x+dx2, y+dy2
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] == 0:
                        draw.point((nx, ny), fill=(74, 0, 112, 200))

    path = os.path.join(OUT_DIR, "T166_star_portal.png")
    img.save(path)
    pixels = sum(1 for y in range(SIZE) for x in range(SIZE) if img.getpixel((x,y))[3] > 0)
    print(f"T166 saved: {pixels} non-transparent pixels ({pixels/(SIZE*SIZE)*100:.1f}%)")

# ── T167 幸運龍魂融合魚 ──────────────────────────────────────
def gen_t167():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2

    # 外層龍魂光暈
    for r in range(30, 26, -1):
        alpha = int(35 + (30-r)*10)
        draw_ring(draw, cx, cy, r, r-1, (211, 47, 47, alpha))

    # 主體橢圓魚身（龍紅色）
    fill_circle_shaded(img, cx, cy, 20,
        (211, 47, 47, 255),   # 龍紅
        (239, 83, 80, 255),   # 亮紅
        (183, 28, 28, 255))   # 暗紅

    # 龍鱗紋路（菱形格）
    pixels = img.load()
    for y in range(cy-18, cy+18):
        for x in range(cx-18, cx+18):
            if 0 <= x < SIZE and 0 <= y < SIZE and img.getpixel((x,y))[3] > 0:
                if (abs(x-cx) + abs(y-cy)) % 6 == 0:
                    pixels[x, y] = (183, 28, 28, 255)

    # 8 方向龍魂火焰射線
    for i in range(8):
        angle = i * 45
        draw_ray(draw, cx, cy, angle, 12, 2, (255, 111, 0, 180))

    # 龍魂光環
    draw_ring(draw, cx, cy, 28, 26, (255, 111, 0, 160))
    draw_ring(draw, cx, cy, 25, 23, (211, 47, 47, 120))

    # 眼睛（火焰眼）
    fill_circle(draw, cx-5, cy-3, 3, (255, 111, 0, 255))
    fill_circle(draw, cx+5, cy-3, 3, (255, 111, 0, 255))
    fill_circle(draw, cx-5, cy-3, 1, (255, 255, 200, 255))
    fill_circle(draw, cx+5, cy-3, 1, (255, 255, 200, 255))

    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 0:
                for dx2, dy2 in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x+dx2, y+dy2
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] == 0:
                        draw.point((nx, ny), fill=(183, 28, 28, 200))

    path = os.path.join(OUT_DIR, "T167_dragon_soul.png")
    img.save(path)
    pixels = sum(1 for y in range(SIZE) for x in range(SIZE) if img.getpixel((x,y))[3] > 0)
    print(f"T167 saved: {pixels} non-transparent pixels ({pixels/(SIZE*SIZE)*100:.1f}%)")

# ── T168 幸運時空裂縫魚 ──────────────────────────────────────
def gen_t168():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2

    # 外層時空光暈
    for r in range(30, 26, -1):
        alpha = int(30 + (30-r)*8)
        draw_ring(draw, cx, cy, r, r-1, (21, 68, 191, alpha))

    # 主體橢圓魚身（深藍色）
    fill_circle_shaded(img, cx, cy, 20,
        (21, 68, 191, 255),   # 深藍
        (66, 133, 244, 255),  # 亮藍
        (13, 42, 120, 255))   # 暗藍

    # 時空裂縫紋路（Z字形）
    pixels = img.load()
    for i in range(-15, 15):
        x = cx + i
        y = cy + (i % 6) - 3
        if 0 <= x < SIZE and 0 <= y < SIZE and img.getpixel((x,y))[3] > 0:
            pixels[x, y] = (100, 181, 246, 255)

    # 時鐘符號（圓形 + 指針）
    draw_ring(draw, cx, cy, 12, 10, (100, 181, 246, 200))
    draw_ray(draw, cx, cy, -90, 8, 1, (255, 255, 255, 220))  # 12點方向
    draw_ray(draw, cx, cy, 0, 6, 1, (255, 255, 255, 220))    # 3點方向

    # 時空光環
    draw_ring(draw, cx, cy, 28, 26, (100, 181, 246, 160))

    # 眼睛（藍色）
    fill_circle(draw, cx-5, cy-3, 3, (100, 181, 246, 255))
    fill_circle(draw, cx+5, cy-3, 3, (100, 181, 246, 255))
    fill_circle(draw, cx-5, cy-3, 1, (255, 255, 255, 255))
    fill_circle(draw, cx+5, cy-3, 1, (255, 255, 255, 255))

    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 0:
                for dx2, dy2 in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x+dx2, y+dy2
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] == 0:
                        draw.point((nx, ny), fill=(13, 42, 120, 200))

    path = os.path.join(OUT_DIR, "T168_spacetime_rift.png")
    img.save(path)
    pixels = sum(1 for y in range(SIZE) for x in range(SIZE) if img.getpixel((x,y))[3] > 0)
    print(f"T168 saved: {pixels} non-transparent pixels ({pixels/(SIZE*SIZE)*100:.1f}%)")

# ── T169 幸運神聖審判魚 ──────────────────────────────────────
def gen_t169():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2

    # 外層神聖光暈
    for r in range(31, 26, -1):
        alpha = int(40 + (31-r)*10)
        draw_ring(draw, cx, cy, r, r-1, (245, 127, 23, alpha))

    # 主體大型圓形魚身（神聖橙金色）
    fill_circle_shaded(img, cx, cy, 22,
        (245, 127, 23, 255),   # 神聖橙
        (255, 183, 77, 255),   # 亮金
        (230, 81, 0, 255))     # 暗橙

    # 12 方向神聖光柱
    for i in range(12):
        angle = i * 30
        draw_ray(draw, cx, cy, angle, 14, 2, (255, 213, 79, 200))

    # 天秤符號（橫線 + 兩個圓）
    draw.line([(cx-10, cy), (cx+10, cy)], fill=(255, 255, 255, 220), width=2)
    draw.line([(cx, cy-8), (cx, cy)], fill=(255, 255, 255, 220), width=2)
    fill_circle(draw, cx-10, cy+4, 4, (255, 255, 255, 180))
    fill_circle(draw, cx+10, cy+4, 4, (255, 255, 255, 180))

    # 神聖光環（雙層）
    draw_ring(draw, cx, cy, 29, 27, (255, 213, 79, 180))
    draw_ring(draw, cx, cy, 26, 24, (245, 127, 23, 140))

    # 眼睛（金色）
    fill_circle(draw, cx-6, cy-5, 3, (255, 213, 79, 255))
    fill_circle(draw, cx+6, cy-5, 3, (255, 213, 79, 255))
    fill_circle(draw, cx-6, cy-5, 1, (255, 255, 255, 255))
    fill_circle(draw, cx+6, cy-5, 1, (255, 255, 255, 255))

    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 0:
                for dx2, dy2 in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x+dx2, y+dy2
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] == 0:
                        draw.point((nx, ny), fill=(230, 81, 0, 200))

    path = os.path.join(OUT_DIR, "T169_holy_judgment.png")
    img.save(path)
    pixels = sum(1 for y in range(SIZE) for x in range(SIZE) if img.getpixel((x,y))[3] > 0)
    print(f"T169 saved: {pixels} non-transparent pixels ({pixels/(SIZE*SIZE)*100:.1f}%)")

# ── T170 幸運宇宙大爆炸魚 ──────────────────────────────────────
def gen_t170():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2

    # 外層爆炸光暈（最大）
    for r in range(32, 26, -1):
        alpha = int(50 + (32-r)*12)
        draw_ring(draw, cx, cy, r, r-1, (183, 18, 18, alpha))

    # 主體大型圓形魚身（深紅色）
    fill_circle_shaded(img, cx, cy, 22,
        (183, 18, 18, 255),   # 深紅
        (229, 57, 53, 255),   # 亮紅
        (127, 0, 0, 255))     # 暗紅

    # 16 方向爆炸光芒（最多方向）
    for i in range(16):
        angle = i * 22.5
        length = 14 if i % 2 == 0 else 10
        draw_ray(draw, cx, cy, angle, length, 2, (255, 87, 34, 200))

    # 宇宙粒子（隨機散佈）
    rng2 = random.Random(170)
    for _ in range(20):
        px = rng2.randint(cx-28, cx+28)
        py = rng2.randint(cy-28, cy+28)
        if 0 <= px < SIZE and 0 <= py < SIZE:
            r2 = (px-cx)**2 + (py-cy)**2
            if 22*22 < r2 < 30*30:
                draw.point((px, py), fill=(255, 213, 79, 200))

    # 爆炸光環（三層）
    draw_ring(draw, cx, cy, 30, 28, (255, 87, 34, 180))
    draw_ring(draw, cx, cy, 27, 25, (183, 18, 18, 150))
    draw_ring(draw, cx, cy, 24, 22, (255, 213, 79, 120))

    # 中心爆炸核心（白色）
    fill_circle(draw, cx, cy, 5, (255, 255, 255, 200))

    # 眼睛（火紅眼）
    fill_circle(draw, cx-6, cy-5, 3, (255, 87, 34, 255))
    fill_circle(draw, cx+6, cy-5, 3, (255, 87, 34, 255))
    fill_circle(draw, cx-6, cy-5, 1, (255, 255, 255, 255))
    fill_circle(draw, cx+6, cy-5, 1, (255, 255, 255, 255))

    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 0:
                for dx2, dy2 in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x+dx2, y+dy2
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] == 0:
                        draw.point((nx, ny), fill=(127, 0, 0, 200))

    path = os.path.join(OUT_DIR, "T170_big_bang.png")
    img.save(path)
    pixels = sum(1 for y in range(SIZE) for x in range(SIZE) if img.getpixel((x,y))[3] > 0)
    print(f"T170 saved: {pixels} non-transparent pixels ({pixels/(SIZE*SIZE)*100:.1f}%)")

if __name__ == "__main__":
    print("=== DAY-312 T166-T170 精靈圖生成 ===")
    gen_t166()
    gen_t167()
    gen_t168()
    gen_t169()
    gen_t170()
    print("=== 全部完成 ===")
