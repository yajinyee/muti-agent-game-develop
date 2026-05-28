"""
generate_targets_day324.py — DAY-324 T211-T215 精靈圖生成
T211 幸運雪崩魚（冰藍色，Avalanche Cascade）
T212 幸運崩潰倍率魚（火橙紅色，Crash Multiplier）
T213 幸運倍率梯魚（金色，Multiplier Ladder）
T214 幸運冰釣輪盤魚（冰青色，Ice Fishing Wheel）
T215 幸運全服雪崩魚（天藍色，Global Avalanche 新史上最高）
"""
import os
import math
from PIL import Image, ImageDraw

SIZE = 64
OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
os.makedirs(OUT_DIR, exist_ok=True)

def fill_circle(draw, cx, cy, r, color):
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], fill=color)

def fill_circle_shaded(img, cx, cy, r, base_color):
    """帶陰影的圓形"""
    pixels = img.load()
    br, bg, bb = base_color[0], base_color[1], base_color[2]
    for y in range(max(0, cy-r), min(SIZE, cy+r+1)):
        for x in range(max(0, cx-r), min(SIZE, cx+r+1)):
            dx, dy = x - cx, y - cy
            if dx*dx + dy*dy <= r*r:
                dist = math.sqrt(dx*dx + dy*dy)
                # 光源左上
                lx, ly = -0.6, -0.6
                nx, ny = dx/max(dist,1), dy/max(dist,1)
                dot = max(0, -(nx*lx + ny*ly))
                light = 0.55 + 0.45 * dot
                edge = 1.0 - (dist/r) * 0.25
                factor = light * edge
                r2 = min(255, int(br * factor))
                g2 = min(255, int(bg * factor))
                b2 = min(255, int(bb * factor))
                pixels[x, y] = (r2, g2, b2, 255)

def draw_rays(draw, cx, cy, n_rays, length, color, width=2):
    """放射狀光芒"""
    for i in range(n_rays):
        angle = 2 * math.pi * i / n_rays
        ex = cx + math.cos(angle) * length
        ey = cy + math.sin(angle) * length
        draw.line([cx, cy, ex, ey], fill=color, width=width)

def draw_ring(draw, cx, cy, r, color, width=2):
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], outline=color, width=width)

def add_lucky_badge(draw, mult_text, badge_color):
    """右上角 Lucky 徽章"""
    draw.rectangle([42, 2, 62, 16], fill=(0, 0, 0, 200))
    draw.text((44, 3), mult_text, fill=badge_color)

# ── T211 幸運雪崩魚（冰藍色，Avalanche Cascade）────────────────
def gen_t211():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    pixels = img.load()

    # 冰藍色魚身（橢圓）
    fill_circle_shaded(img, 32, 32, 22, (0, 180, 240))

    # 8 波雪崩光芒（冰藍白色）
    draw_rays(draw, 32, 32, 8, 28, (180, 230, 255, 220), width=2)

    # 雪花符號（6 個方向）
    for i in range(6):
        angle = math.pi / 3 * i
        ex = 32 + math.cos(angle) * 18
        ey = 32 + math.sin(angle) * 18
        draw.line([32, 32, ex, ey], fill=(255, 255, 255, 200), width=1)

    # 三層冰環
    draw_ring(draw, 32, 32, 26, (0, 200, 255, 180), width=2)
    draw_ring(draw, 32, 32, 20, (100, 220, 255, 150), width=1)
    draw_ring(draw, 32, 32, 14, (200, 240, 255, 120), width=1)

    # 中心冰晶
    fill_circle(draw, 32, 32, 6, (220, 245, 255, 255))
    fill_circle(draw, 32, 32, 3, (255, 255, 255, 255))

    # 雪崩粒子（8 個散落冰晶）
    import random
    rng = random.Random(211)
    for _ in range(8):
        px = rng.randint(4, 60)
        py = rng.randint(4, 60)
        fill_circle(draw, px, py, 2, (180, 230, 255, 180))

    # Lucky 徽章
    draw.rectangle([38, 2, 62, 14], fill=(0, 0, 0, 200))
    draw.text((40, 3), "❄️36x", fill=(0, 200, 255, 255))

    img = img.resize((SIZE, SIZE), Image.NEAREST)
    img.save(os.path.join(OUT_DIR, "T211_avalanche.png"))
    print(f"T211 saved, density={sum(1 for y in range(SIZE) for x in range(SIZE) if img.getpixel((x,y))[3]>0)/(SIZE*SIZE)*100:.1f}%")

# ── T212 幸運崩潰倍率魚（火橙紅色，Crash Multiplier）────────────
def gen_t212():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)

    # 火橙紅色魚身
    fill_circle_shaded(img, 32, 32, 22, (220, 60, 0))

    # 崩潰裂縫（Z 字形）
    draw.line([20, 20, 32, 32, 20, 44], fill=(255, 200, 0, 220), width=3)
    draw.line([44, 20, 32, 32, 44, 44], fill=(255, 200, 0, 220), width=3)

    # 上升箭頭（倍率上升感）
    draw.polygon([(32, 8), (26, 18), (38, 18)], fill=(255, 255, 0, 220))
    draw.line([32, 18, 32, 28], fill=(255, 255, 0, 200), width=2)

    # 爆炸光芒（12 道）
    draw_rays(draw, 32, 32, 12, 28, (255, 100, 0, 180), width=2)

    # 三層爆炸環
    draw_ring(draw, 32, 32, 26, (255, 80, 0, 180), width=2)
    draw_ring(draw, 32, 32, 20, (255, 150, 0, 150), width=1)

    # 中心爆炸
    fill_circle(draw, 32, 32, 7, (255, 200, 0, 255))
    fill_circle(draw, 32, 32, 4, (255, 255, 255, 255))

    # Lucky 徽章
    draw.rectangle([38, 2, 62, 14], fill=(0, 0, 0, 200))
    draw.text((40, 3), "💥36.5x", fill=(255, 100, 0, 255))

    img = img.resize((SIZE, SIZE), Image.NEAREST)
    img.save(os.path.join(OUT_DIR, "T212_crash_multiplier.png"))
    print(f"T212 saved, density={sum(1 for y in range(SIZE) for x in range(SIZE) if img.getpixel((x,y))[3]>0)/(SIZE*SIZE)*100:.1f}%")

# ── T213 幸運倍率梯魚（金色，Multiplier Ladder）────────────────
def gen_t213():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)

    # 金色魚身
    fill_circle_shaded(img, 32, 32, 22, (200, 160, 0))

    # 梯子符號（10 級）
    for i in range(10):
        y = 10 + i * 4
        width = 8 + i * 2
        draw.line([32 - width//2, y, 32 + width//2, y], fill=(255, 220, 0, 200), width=1)

    # 兩側梯柱
    draw.line([20, 10, 20, 50], fill=(255, 200, 0, 200), width=2)
    draw.line([44, 10, 44, 50], fill=(255, 200, 0, 200), width=2)

    # 16 道金色光芒
    draw_rays(draw, 32, 32, 16, 28, (255, 200, 0, 160), width=1)

    # 三層金環
    draw_ring(draw, 32, 32, 26, (255, 200, 0, 180), width=2)
    draw_ring(draw, 32, 32, 20, (255, 220, 100, 150), width=1)

    # 中心金星
    fill_circle(draw, 32, 32, 6, (255, 230, 0, 255))
    fill_circle(draw, 32, 32, 3, (255, 255, 200, 255))

    # Lucky 徽章
    draw.rectangle([38, 2, 62, 14], fill=(0, 0, 0, 200))
    draw.text((40, 3), "🪜37x", fill=(255, 200, 0, 255))

    img = img.resize((SIZE, SIZE), Image.NEAREST)
    img.save(os.path.join(OUT_DIR, "T213_multiplier_ladder.png"))
    print(f"T213 saved, density={sum(1 for y in range(SIZE) for x in range(SIZE) if img.getpixel((x,y))[3]>0)/(SIZE*SIZE)*100:.1f}%")

# ── T214 幸運冰釣輪盤魚（冰青色，Ice Fishing Wheel）────────────
def gen_t214():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)

    # 冰青色大型魚身
    fill_circle_shaded(img, 32, 32, 24, (0, 180, 190))

    # 輪盤（5 個扇形）
    colors = [
        (0, 200, 100, 200),   # 綠（×100）
        (0, 100, 255, 200),   # 藍（×500）
        (255, 150, 0, 200),   # 橙（×1000）
        (200, 0, 200, 200),   # 紫（×2000）
        (255, 255, 255, 200), # 白（×5000）
    ]
    for i, color in enumerate(colors):
        start_angle = i * 72 - 90
        end_angle = start_angle + 72
        draw.pieslice([14, 14, 50, 50], start=start_angle, end=end_angle, fill=color)

    # 輪盤邊框
    draw_ring(draw, 32, 32, 18, (255, 255, 255, 200), width=2)

    # 指針
    draw.polygon([(32, 14), (30, 22), (34, 22)], fill=(255, 255, 255, 255))

    # 外層冰環
    draw_ring(draw, 32, 32, 28, (0, 220, 230, 180), width=2)
    draw_ring(draw, 32, 32, 24, (100, 230, 240, 150), width=1)

    # 冰晶裝飾
    for i in range(6):
        angle = math.pi / 3 * i
        ex = 32 + math.cos(angle) * 30
        ey = 32 + math.sin(angle) * 30
        fill_circle(draw, int(ex), int(ey), 2, (200, 240, 255, 200))

    # Lucky 徽章
    draw.rectangle([36, 2, 62, 14], fill=(0, 0, 0, 200))
    draw.text((38, 3), "❄️37.5x", fill=(0, 200, 220, 255))

    img = img.resize((SIZE, SIZE), Image.NEAREST)
    img.save(os.path.join(OUT_DIR, "T214_ice_fishing_wheel.png"))
    print(f"T214 saved, density={sum(1 for y in range(SIZE) for x in range(SIZE) if img.getpixel((x,y))[3]>0)/(SIZE*SIZE)*100:.1f}%")

# ── T215 幸運全服雪崩魚（天藍色，Global Avalanche 新史上最高）──
def gen_t215():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)

    # 天藍色超大魚身
    fill_circle_shaded(img, 32, 32, 26, (100, 180, 240))

    # 40 道全服光芒（最多光芒）
    draw_rays(draw, 32, 32, 40, 30, (180, 220, 255, 160), width=1)

    # 全球符號（地球輪廓）
    draw_ring(draw, 32, 32, 20, (255, 255, 255, 200), width=2)
    # 赤道線
    draw.line([12, 32, 52, 32], fill=(255, 255, 255, 180), width=1)
    # 子午線
    draw.arc([18, 14, 46, 50], start=0, end=180, fill=(255, 255, 255, 150), width=1)
    draw.arc([18, 14, 46, 50], start=180, end=360, fill=(255, 255, 255, 150), width=1)

    # 五層光環（最多層）
    draw_ring(draw, 32, 32, 28, (135, 206, 250, 200), width=2)
    draw_ring(draw, 32, 32, 24, (150, 215, 255, 180), width=2)
    draw_ring(draw, 32, 32, 20, (180, 225, 255, 160), width=1)
    draw_ring(draw, 32, 32, 16, (200, 235, 255, 140), width=1)
    draw_ring(draw, 32, 32, 12, (220, 245, 255, 120), width=1)

    # 中心全服符號
    fill_circle(draw, 32, 32, 7, (255, 255, 255, 255))
    fill_circle(draw, 32, 32, 4, (135, 206, 250, 255))

    # 雪崩粒子（12 個散落）
    import random
    rng = random.Random(215)
    for _ in range(12):
        px = rng.randint(3, 61)
        py = rng.randint(3, 61)
        fill_circle(draw, px, py, 2, (200, 230, 255, 180))

    # Lucky 徽章（新史上最高標記）
    draw.rectangle([30, 2, 62, 14], fill=(0, 0, 0, 220))
    draw.text((32, 3), "🌍38x NEW", fill=(135, 206, 250, 255))

    img = img.resize((SIZE, SIZE), Image.NEAREST)
    img.save(os.path.join(OUT_DIR, "T215_global_avalanche.png"))
    print(f"T215 saved, density={sum(1 for y in range(SIZE) for x in range(SIZE) if img.getpixel((x,y))[3]>0)/(SIZE*SIZE)*100:.1f}%")

if __name__ == "__main__":
    print("=== DAY-324 T211-T215 精靈圖生成 ===")
    gen_t211()
    gen_t212()
    gen_t213()
    gen_t214()
    gen_t215()
    print("=== 全部完成 ===")
