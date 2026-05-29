"""
generate_targets_day325.py — DAY-325 T216-T220 精靈圖生成
業界依據：
  T216 漁網魚：BGaming Fishing Club 2 Fishing Net（2026-04）
  T217 TNT 爆炸魚：BGaming Fishing Club 2 TNT Bonus（2026-04）
  T218 擾動魚：Fisch Disturbance System（2026-01）
  T219 珍珠倍率魚：BGaming Shark & Spark Hold & Win（2026-05）
  T220 快速暴富魚：Reflex Gaming Big Game Fishing Rapid Riches（2026-05）
"""
import os
import math
import random
from PIL import Image, ImageDraw

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64
os.makedirs(OUT_DIR, exist_ok=True)

rng = random.Random(325)

def fill_circle(draw, cx, cy, r, color):
    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            if (x - cx) ** 2 + (y - cy) ** 2 <= r ** 2:
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    draw.point((x, y), fill=color)

def fill_ellipse_shaded(draw, cx, cy, rx, ry, base_color, light_color, dark_color):
    for y in range(cy - ry, cy + ry + 1):
        for x in range(cx - rx, cx + rx + 1):
            if (x - cx) ** 2 / rx ** 2 + (y - cy) ** 2 / ry ** 2 <= 1.0:
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    if x < cx - rx * 0.3 and y < cy - ry * 0.3:
                        draw.point((x, y), fill=light_color)
                    elif x > cx + rx * 0.3 and y > cy + ry * 0.3:
                        draw.point((x, y), fill=dark_color)
                    else:
                        draw.point((x, y), fill=base_color)

def draw_rays(draw, cx, cy, n_rays, r_inner, r_outer, color, width=1):
    for i in range(n_rays):
        angle = 2 * math.pi * i / n_rays
        x1 = int(cx + r_inner * math.cos(angle))
        y1 = int(cy + r_inner * math.sin(angle))
        x2 = int(cx + r_outer * math.cos(angle))
        y2 = int(cy + r_outer * math.sin(angle))
        draw.line([(x1, y1), (x2, y2)], fill=color, width=width)

def draw_ring(draw, cx, cy, r, color, width=2):
    for i in range(360):
        angle = math.radians(i)
        for w in range(width):
            x = int(cx + (r + w) * math.cos(angle))
            y = int(cy + (r + w) * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=color)

def count_pixels(img):
    pixels = img.load()
    count = 0
    for y in range(SIZE):
        for x in range(SIZE):
            if pixels[x, y][3] > 0:
                count += 1
    return count

# ── T216 幸運漁網魚 ──────────────────────────────────────────
# 深海藍魚身 + 漁網紋路 + 魚鉤符號 + 多層光環
def gen_T216():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 外層光暈（深海藍）
    fill_ellipse_shaded(draw, cx, cy, 28, 22, (20, 100, 200, 60), (40, 130, 220, 40), (10, 70, 160, 50))

    # 主體魚身（深海藍橢圓）
    fill_ellipse_shaded(draw, cx, cy, 22, 16,
        (30, 144, 255, 255),   # 深海藍
        (80, 180, 255, 255),   # 亮藍高光
        (15, 90, 180, 255))    # 深藍陰影

    # 漁網紋路（網格線）
    net_color = (0, 80, 180, 180)
    for i in range(-3, 4):
        # 斜線 /
        x1 = cx + i * 6 - 20
        y1 = cy - 16
        x2 = cx + i * 6 + 20
        y2 = cy + 16
        draw.line([(x1, y1), (x2, y2)], fill=net_color, width=1)
        # 斜線 \
        x1 = cx - i * 6 - 20
        y1 = cy - 16
        x2 = cx - i * 6 + 20
        y2 = cy + 16
        draw.line([(x1, y1), (x2, y2)], fill=net_color, width=1)

    # 魚鉤符號（右側）
    hook_color = (200, 220, 255, 255)
    draw.arc([(cx + 8, cy - 8), (cx + 18, cy + 2)], start=0, end=180, fill=hook_color, width=2)
    draw.line([(cx + 8, cy - 3), (cx + 8, cy + 8)], fill=hook_color, width=2)
    draw.arc([(cx + 5, cy + 5), (cx + 11, cy + 11)], start=180, end=360, fill=hook_color, width=2)

    # 光芒（12 道）
    draw_rays(draw, cx, cy, 12, 22, 30, (100, 180, 255, 200), 1)

    # 外圈光環
    draw_ring(draw, cx, cy, 26, (0, 150, 255, 150), 2)

    # 高光點
    fill_circle(draw, cx - 8, cy - 6, 3, (200, 230, 255, 220))

    density = count_pixels(img) / (SIZE * SIZE) * 100
    path = os.path.join(OUT_DIR, "T216_fishing_net.png")
    img.save(path)
    print(f"T216 漁網魚 → {path} ({density:.1f}%)")
    return density

# ── T217 幸運 TNT 爆炸魚 ──────────────────────────────────────
# 火橙紅魚身 + TNT 符號 + 爆炸光芒 + 裂縫紋路
def gen_T217():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 外層爆炸光暈
    fill_ellipse_shaded(draw, cx, cy, 30, 24, (255, 100, 0, 50), (255, 150, 50, 30), (200, 50, 0, 60))

    # 主體魚身（火橙紅橢圓）
    fill_ellipse_shaded(draw, cx, cy, 22, 16,
        (255, 69, 0, 255),     # 火橙紅
        (255, 140, 0, 255),    # 橙黃高光
        (180, 30, 0, 255))     # 深紅陰影

    # 爆炸裂縫紋路
    crack_color = (255, 200, 0, 200)
    # 中心向外的裂縫
    for angle_deg in [30, 90, 150, 210, 270, 330]:
        angle = math.radians(angle_deg)
        x1, y1 = cx, cy
        x2 = int(cx + 14 * math.cos(angle))
        y2 = int(cy + 14 * math.sin(angle))
        draw.line([(x1, y1), (x2, y2)], fill=crack_color, width=1)

    # TNT 文字符號（中心）
    tnt_color = (255, 255, 0, 255)
    # T
    draw.line([(cx - 8, cy - 4), (cx - 2, cy - 4)], fill=tnt_color, width=2)
    draw.line([(cx - 5, cy - 4), (cx - 5, cy + 4)], fill=tnt_color, width=2)
    # N
    draw.line([(cx - 1, cy - 4), (cx - 1, cy + 4)], fill=tnt_color, width=2)
    draw.line([(cx - 1, cy - 4), (cx + 3, cy + 4)], fill=tnt_color, width=2)
    draw.line([(cx + 3, cy - 4), (cx + 3, cy + 4)], fill=tnt_color, width=2)
    # T
    draw.line([(cx + 4, cy - 4), (cx + 10, cy - 4)], fill=tnt_color, width=2)
    draw.line([(cx + 7, cy - 4), (cx + 7, cy + 4)], fill=tnt_color, width=2)

    # 爆炸光芒（16 道）
    draw_rays(draw, cx, cy, 16, 22, 31, (255, 150, 0, 200), 1)

    # 外圈火焰環
    draw_ring(draw, cx, cy, 27, (255, 80, 0, 180), 2)
    draw_ring(draw, cx, cy, 29, (255, 200, 0, 120), 1)

    # 高光點
    fill_circle(draw, cx - 8, cy - 6, 3, (255, 220, 150, 220))

    density = count_pixels(img) / (SIZE * SIZE) * 100
    path = os.path.join(OUT_DIR, "T217_tnt_bonus.png")
    img.save(path)
    print(f"T217 TNT 爆炸魚 → {path} ({density:.1f}%)")
    return density

# ── T218 幸運擾動魚 ──────────────────────────────────────────
# 深青色魚身 + 波紋擾動紋路 + 活躍度指示器 + 漩渦光環
def gen_T218():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 外層擾動光暈
    fill_ellipse_shaded(draw, cx, cy, 29, 23, (0, 180, 200, 50), (0, 220, 230, 30), (0, 120, 150, 60))

    # 主體魚身（深青色橢圓）
    fill_ellipse_shaded(draw, cx, cy, 22, 16,
        (0, 206, 209, 255),    # 深青
        (0, 255, 240, 255),    # 亮青高光
        (0, 139, 139, 255))    # 深青陰影

    # 波紋擾動紋路（同心橢圓波）
    wave_color = (0, 255, 220, 150)
    for r in [8, 13, 18]:
        draw_ring(draw, cx, cy, r, wave_color, 1)

    # 活躍度指示器（右側 5 個小圓點）
    for i in range(5):
        dot_x = cx + 14
        dot_y = cy - 8 + i * 4
        alpha = 200 - i * 20
        fill_circle(draw, dot_x, dot_y, 2, (0, 255, 200, alpha))

    # 漩渦光芒（8 道）
    draw_rays(draw, cx, cy, 8, 22, 30, (0, 200, 220, 200), 1)

    # 外圈光環
    draw_ring(draw, cx, cy, 26, (0, 180, 200, 160), 2)

    # 高光點
    fill_circle(draw, cx - 8, cy - 6, 3, (150, 255, 250, 220))

    density = count_pixels(img) / (SIZE * SIZE) * 100
    path = os.path.join(OUT_DIR, "T218_disturbance.png")
    img.save(path)
    print(f"T218 擾動魚 → {path} ({density:.1f}%)")
    return density

# ── T219 幸運珍珠倍率魚 ──────────────────────────────────────
# 金色大型魚身 + 珍珠符號 + 倍率數字 + 四層光環（里程碑）
def gen_T219():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 外層金色光暈（里程碑等級）
    fill_ellipse_shaded(draw, cx, cy, 31, 25, (255, 215, 0, 60), (255, 240, 100, 40), (200, 160, 0, 70))

    # 主體魚身（金色大型橢圓）
    fill_ellipse_shaded(draw, cx, cy, 24, 18,
        (255, 215, 0, 255),    # 金色
        (255, 255, 150, 255),  # 亮金高光
        (180, 140, 0, 255))    # 深金陰影

    # 珍珠符號（中心圓形）
    fill_circle(draw, cx, cy, 7, (255, 255, 240, 255))
    fill_circle(draw, cx - 2, cy - 2, 2, (255, 255, 255, 255))  # 高光

    # 珍珠光暈
    draw_ring(draw, cx, cy, 9, (255, 240, 100, 200), 2)

    # 倍率符號（×40 文字）
    mult_color = (255, 200, 0, 255)
    # × 符號
    draw.line([(cx - 18, cy - 8), (cx - 12, cy - 2)], fill=mult_color, width=2)
    draw.line([(cx - 12, cy - 8), (cx - 18, cy - 2)], fill=mult_color, width=2)
    # 4
    draw.line([(cx - 10, cy - 8), (cx - 10, cy - 4)], fill=mult_color, width=2)
    draw.line([(cx - 10, cy - 4), (cx - 6, cy - 4)], fill=mult_color, width=2)
    draw.line([(cx - 6, cy - 8), (cx - 6, cy - 2)], fill=mult_color, width=2)
    # 0
    draw.ellipse([(cx - 4, cy - 8), (cx + 2, cy - 2)], outline=mult_color, width=2)

    # 光芒（20 道，里程碑等級）
    draw_rays(draw, cx, cy, 20, 24, 31, (255, 220, 50, 200), 1)

    # 四層光環（里程碑）
    draw_ring(draw, cx, cy, 26, (255, 215, 0, 180), 2)
    draw_ring(draw, cx, cy, 28, (255, 240, 100, 130), 1)
    draw_ring(draw, cx, cy, 30, (255, 200, 0, 80), 1)

    # 高光點
    fill_circle(draw, cx - 9, cy - 7, 3, (255, 255, 200, 230))

    density = count_pixels(img) / (SIZE * SIZE) * 100
    path = os.path.join(OUT_DIR, "T219_pearl_multiplier.png")
    img.save(path)
    print(f"T219 珍珠倍率魚 → {path} ({density:.1f}%)")
    return density

# ── T220 幸運快速暴富魚 ──────────────────────────────────────
# 亮黃色超大型魚身 + 閃電符號 + 速度線 + 五層光環（新史上最高）
def gen_T220():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 外層超強光暈（新史上最高等級）
    fill_ellipse_shaded(draw, cx, cy, 32, 26, (255, 255, 0, 70), (255, 255, 150, 50), (200, 200, 0, 80))

    # 主體魚身（亮黃色超大型橢圓）
    fill_ellipse_shaded(draw, cx, cy, 25, 19,
        (255, 255, 0, 255),    # 亮黃
        (255, 255, 200, 255),  # 超亮高光
        (200, 200, 0, 255))    # 深黃陰影

    # 閃電符號（中心）
    lightning_color = (255, 140, 0, 255)
    # 大閃電
    points = [(cx - 4, cy - 10), (cx + 2, cy - 2), (cx - 2, cy - 2), (cx + 4, cy + 10)]
    for i in range(len(points) - 1):
        draw.line([points[i], points[i + 1]], fill=lightning_color, width=3)

    # 速度線（左側）
    speed_color = (255, 200, 0, 200)
    for i in range(4):
        y_offset = cy - 6 + i * 4
        x_start = cx - 22 + i * 2
        draw.line([(x_start, y_offset), (cx - 12, y_offset)], fill=speed_color, width=1)

    # 光芒（24 道，最高等級）
    draw_rays(draw, cx, cy, 24, 25, 32, (255, 240, 0, 200), 1)

    # 五層光環（新史上最高）
    draw_ring(draw, cx, cy, 27, (255, 255, 0, 200), 2)
    draw_ring(draw, cx, cy, 29, (255, 220, 0, 160), 2)
    draw_ring(draw, cx, cy, 31, (255, 200, 0, 120), 1)

    # 高光點
    fill_circle(draw, cx - 10, cy - 8, 4, (255, 255, 220, 240))
    fill_circle(draw, cx - 6, cy - 4, 2, (255, 255, 255, 200))

    density = count_pixels(img) / (SIZE * SIZE) * 100
    path = os.path.join(OUT_DIR, "T220_rapid_riches.png")
    img.save(path)
    print(f"T220 快速暴富魚 → {path} ({density:.1f}%)")
    return density

if __name__ == "__main__":
    print("=== DAY-325 T216-T220 精靈圖生成 ===")
    d216 = gen_T216()
    d217 = gen_T217()
    d218 = gen_T218()
    d219 = gen_T219()
    d220 = gen_T220()
    print(f"\n平均密度：{(d216 + d217 + d218 + d219 + d220) / 5:.1f}%")
    print("✅ 全部完成！")
