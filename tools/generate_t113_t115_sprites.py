"""
generate_t113_t115_sprites.py — T113/T114/T115 精靈圖生成
T113 幸運鑽頭魚雷魚：橙色機械魚身 + 鑽頭前端 + 魚雷尾焰
T114 幸運時間凍結魚：冰藍色魚身 + 冰晶紋路 + 凍結光環
T115 幸運連鎖爆炸魚：深紅色魚身 + 爆炸紋路 + 連鎖光環
"""
import os
from PIL import Image, ImageDraw

SIZE = 64
OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
os.makedirs(OUT_DIR, exist_ok=True)

def px(img, x, y, color):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), color)

def fill_ellipse_shaded(img, cx, cy, rx, ry, light, mid, dark):
    for y in range(cy - ry, cy + ry + 1):
        for x in range(cx - rx, cx + rx + 1):
            if ((x - cx) / rx) ** 2 + ((y - cy) / ry) ** 2 <= 1.0:
                if x < cx - rx * 0.3 and y < cy - ry * 0.3:
                    c = light
                elif x > cx + rx * 0.3 and y > cy + ry * 0.3:
                    c = dark
                else:
                    c = mid
                px(img, x, y, c)

def draw_circle(img, cx, cy, r, color):
    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            if (x - cx) ** 2 + (y - cy) ** 2 <= r * r:
                px(img, x, y, color)

def draw_ring(img, cx, cy, r_outer, r_inner, color):
    for y in range(cy - r_outer, cy + r_outer + 1):
        for x in range(cx - r_outer, cx + r_outer + 1):
            d2 = (x - cx) ** 2 + (y - cy) ** 2
            if r_inner * r_inner <= d2 <= r_outer * r_outer:
                px(img, x, y, color)

# ── T113 幸運鑽頭魚雷魚 ──────────────────────────────────────
def gen_t113():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    # 橙色機械魚身（橫向橢圓）
    ORANGE_L = (255, 160, 60, 255)
    ORANGE_M = (220, 110, 30, 255)
    ORANGE_D = (160, 70, 10, 255)
    fill_ellipse_shaded(img, 32, 32, 22, 14, ORANGE_L, ORANGE_M, ORANGE_D)
    # 鑽頭前端（右側三角形）
    DRILL = (180, 180, 200, 255)
    DRILL_L = (220, 220, 240, 255)
    for i in range(12):
        w = 12 - i
        for y in range(32 - w // 2, 32 + w // 2 + 1):
            px(img, 54 + i // 2, y, DRILL if i % 2 == 0 else DRILL_L)
    # 鑽頭螺旋紋
    SPIRAL = (100, 100, 120, 255)
    for i in range(0, 10, 2):
        px(img, 50 + i, 30 + i % 3, SPIRAL)
        px(img, 50 + i, 34 - i % 3, SPIRAL)
    # 尾焰（左側）
    FLAME1 = (255, 80, 0, 255)
    FLAME2 = (255, 200, 0, 255)
    for i in range(8):
        y_off = (i % 3) - 1
        px(img, 10 - i // 2, 32 + y_off, FLAME1 if i < 4 else FLAME2)
        px(img, 10 - i // 2, 31 + y_off, FLAME2 if i < 4 else FLAME1)
    # 機械紋路
    MECH = (140, 80, 20, 255)
    for x in range(20, 50, 6):
        px(img, x, 26, MECH)
        px(img, x, 38, MECH)
    # 眼睛（機械眼）
    draw_circle(img, 44, 29, 3, (40, 40, 60, 255))
    draw_circle(img, 44, 29, 2, (0, 200, 255, 255))
    px(img, 44, 28, (255, 255, 255, 255))
    # 光環
    draw_ring(img, 32, 32, 28, 25, (255, 140, 0, 120))
    img.save(os.path.join(OUT_DIR, "T113_drill_torpedo.png"))
    non_transparent = sum(1 for x in range(SIZE) for y in range(SIZE) if img.getpixel((x, y))[3] > 0)
    print(f"T113 saved: {non_transparent}/{SIZE*SIZE} = {non_transparent/(SIZE*SIZE)*100:.1f}% non-transparent")

# ── T114 幸運時間凍結魚 ──────────────────────────────────────
def gen_t114():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    # 冰藍色魚身
    ICE_L = (180, 240, 255, 255)
    ICE_M = (100, 190, 230, 255)
    ICE_D = (50, 130, 180, 255)
    fill_ellipse_shaded(img, 32, 32, 20, 13, ICE_L, ICE_M, ICE_D)
    # 冰晶紋路（六角形）
    CRYSTAL = (200, 240, 255, 255)
    CRYSTAL_D = (60, 160, 210, 255)
    # 中心六角冰晶
    cx, cy = 32, 32
    for angle_step in range(6):
        import math
        angle = math.radians(angle_step * 60)
        ex = int(cx + 10 * math.cos(angle))
        ey = int(cy + 10 * math.sin(angle))
        # 畫線
        steps = 10
        for s in range(steps + 1):
            lx = int(cx + (ex - cx) * s / steps)
            ly = int(cy + (ey - cy) * s / steps)
            px(img, lx, ly, CRYSTAL)
    # 冰晶尖端
    for angle_step in range(6):
        import math
        angle = math.radians(angle_step * 60 + 30)
        ex = int(cx + 14 * math.cos(angle))
        ey = int(cy + 14 * math.sin(angle))
        px(img, ex, ey, CRYSTAL_D)
        px(img, ex + 1, ey, CRYSTAL_D)
    # 凍結光環
    draw_ring(img, 32, 32, 26, 23, (150, 220, 255, 150))
    draw_ring(img, 32, 32, 30, 28, (100, 180, 230, 80))
    # 眼睛（冰藍眼）
    draw_circle(img, 44, 29, 3, (20, 80, 120, 255))
    draw_circle(img, 44, 29, 2, (180, 240, 255, 255))
    px(img, 44, 28, (255, 255, 255, 255))
    # 魚尾（冰藍）
    TAIL = (80, 160, 210, 255)
    for i in range(6):
        px(img, 12 - i, 28 + i // 2, TAIL)
        px(img, 12 - i, 36 - i // 2, TAIL)
    # 雪花粒子
    SNOW = (220, 245, 255, 200)
    for sx, sy in [(20, 20), (44, 18), (18, 44), (46, 44), (32, 16), (32, 48)]:
        px(img, sx, sy, SNOW)
        px(img, sx + 1, sy, SNOW)
        px(img, sx, sy + 1, SNOW)
    img.save(os.path.join(OUT_DIR, "T114_time_freeze.png"))
    non_transparent = sum(1 for x in range(SIZE) for y in range(SIZE) if img.getpixel((x, y))[3] > 0)
    print(f"T114 saved: {non_transparent}/{SIZE*SIZE} = {non_transparent/(SIZE*SIZE)*100:.1f}% non-transparent")

# ── T115 幸運連鎖爆炸魚 ──────────────────────────────────────
def gen_t115():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    # 深紅色魚身
    RED_L = (255, 100, 80, 255)
    RED_M = (200, 50, 40, 255)
    RED_D = (130, 20, 15, 255)
    fill_ellipse_shaded(img, 32, 32, 20, 13, RED_L, RED_M, RED_D)
    # 爆炸紋路（放射狀裂縫）
    CRACK = (255, 180, 0, 255)
    CRACK_D = (200, 100, 0, 255)
    import math
    cx, cy = 32, 32
    for i in range(8):
        angle = math.radians(i * 45)
        for r in range(4, 16):
            lx = int(cx + r * math.cos(angle))
            ly = int(cy + r * math.sin(angle))
            px(img, lx, ly, CRACK if r < 10 else CRACK_D)
    # 中心爆炸核心
    draw_circle(img, 32, 32, 4, (255, 220, 0, 255))
    draw_circle(img, 32, 32, 2, (255, 255, 200, 255))
    # 連鎖光環（雙環）
    draw_ring(img, 32, 32, 26, 23, (255, 100, 0, 130))
    draw_ring(img, 32, 32, 30, 28, (255, 60, 0, 80))
    # 眼睛（火焰眼）
    draw_circle(img, 44, 29, 3, (80, 10, 5, 255))
    draw_circle(img, 44, 29, 2, (255, 150, 0, 255))
    px(img, 44, 28, (255, 255, 200, 255))
    # 魚尾（深紅）
    TAIL = (160, 30, 20, 255)
    for i in range(6):
        px(img, 12 - i, 28 + i // 2, TAIL)
        px(img, 12 - i, 36 - i // 2, TAIL)
    # 爆炸粒子散落
    SPARK = (255, 200, 50, 200)
    for sx, sy in [(18, 18), (46, 18), (18, 46), (46, 46), (32, 14), (32, 50), (14, 32), (50, 32)]:
        px(img, sx, sy, SPARK)
        px(img, sx + 1, sy, SPARK)
    img.save(os.path.join(OUT_DIR, "T115_chain_explosion.png"))
    non_transparent = sum(1 for x in range(SIZE) for y in range(SIZE) if img.getpixel((x, y))[3] > 0)
    print(f"T115 saved: {non_transparent}/{SIZE*SIZE} = {non_transparent/(SIZE*SIZE)*100:.1f}% non-transparent")

if __name__ == "__main__":
    gen_t113()
    gen_t114()
    gen_t115()
    print("Done! T113-T115 sprites generated.")
