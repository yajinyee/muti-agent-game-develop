"""
generate_targets_day318.py — DAY-318 T196-T200 精靈圖生成
T196 幸運龍王輪盤魚（3100x）：深橙金色魚身 + 雙環輪盤 + 龍紋
T197 幸運永恆循環魚（3200x）：深藍色魚身 + 無限符號 + 循環光環
T198 幸運混沌爆炸魚（3300x）：深紅色魚身 + 混沌爆炸光芒 + 裂縫紋路
T199 幸運神聖復活魚（3400x）：深金色魚身 + 神聖光芒 + 鳳凰羽翼
T200 幸運創世紀元魚（5000x）：純黑色大型魚身 + 25 道光芒 + 四層光環 + 里程碑符文
"""
import os
import math
from PIL import Image, ImageDraw

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def fill_circle(draw, cx, cy, r, color):
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], fill=color)

def fill_circle_shaded(draw, cx, cy, r, base_color, light_color, dark_color):
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                # 左上亮，右下暗
                if dx < -r//3 and dy < -r//3:
                    c = light_color
                elif dx > r//3 and dy > r//3:
                    c = dark_color
                else:
                    c = base_color
                draw.point([cx+dx, cy+dy], fill=c)

def draw_rays(draw, cx, cy, n_rays, r_inner, r_outer, color, width=1):
    for i in range(n_rays):
        angle = 2 * math.pi * i / n_rays
        x1 = cx + r_inner * math.cos(angle)
        y1 = cy + r_inner * math.sin(angle)
        x2 = cx + r_outer * math.cos(angle)
        y2 = cy + r_outer * math.sin(angle)
        draw.line([x1, y1, x2, y2], fill=color, width=width)

def draw_ring(draw, cx, cy, r, color, width=2):
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], outline=color, width=width)

def generate_t196():
    """T196 幸運龍王輪盤魚：深橙金色魚身 + 雙環輪盤 + 龍紋"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2

    # 外層光暈（深橙色）
    for r in range(30, 26, -1):
        alpha = int(60 * (r - 26) / 4)
        draw.ellipse([cx-r, cy-r, cx+r, cy+r], fill=(200, 80, 0, alpha))

    # 魚身（深橙金色，帶陰影）
    fill_circle_shaded(draw, cx, cy, 22,
        (200, 100, 0, 255),    # 基底深橙
        (255, 160, 40, 255),   # 高光橙金
        (140, 60, 0, 255))     # 陰影深棕

    # 雙環輪盤（外環）
    draw_ring(draw, cx, cy, 26, (255, 200, 0, 200), width=2)
    # 雙環輪盤（內環）
    draw_ring(draw, cx, cy, 18, (255, 220, 80, 180), width=1)

    # 輪盤分格線（8格）
    for i in range(8):
        angle = 2 * math.pi * i / 8
        x1 = cx + 18 * math.cos(angle)
        y1 = cy + 18 * math.sin(angle)
        x2 = cx + 26 * math.cos(angle)
        y2 = cy + 26 * math.sin(angle)
        draw.line([x1, y1, x2, y2], fill=(255, 200, 0, 150), width=1)

    # 龍紋（中心龍眼）
    fill_circle(draw, cx, cy, 5, (255, 50, 0, 255))
    fill_circle(draw, cx, cy, 3, (255, 200, 0, 255))
    fill_circle(draw, cx, cy, 1, (255, 255, 255, 255))

    # 輪廓
    draw_ring(draw, cx, cy, 22, (80, 30, 0, 255), width=2)

    # 放大到 128x128
    img = img.resize((128, 128), Image.NEAREST)
    path = os.path.join(OUT_DIR, "T196_dragon_king.png")
    img.save(path)
    print(f"T196 saved: {path}")
    return img

def generate_t197():
    """T197 幸運永恆循環魚：深藍色魚身 + 無限符號 + 循環光環"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2

    # 外層光暈（深藍色）
    for r in range(30, 26, -1):
        alpha = int(60 * (r - 26) / 4)
        draw.ellipse([cx-r, cy-r, cx+r, cy+r], fill=(0, 40, 150, alpha))

    # 魚身（深藍色，帶陰影）
    fill_circle_shaded(draw, cx, cy, 22,
        (0, 60, 180, 255),     # 基底深藍
        (40, 120, 255, 255),   # 高光藍
        (0, 30, 100, 255))     # 陰影深藍

    # 循環光環（3個同心圓）
    draw_ring(draw, cx, cy, 26, (100, 160, 255, 180), width=2)
    draw_ring(draw, cx, cy, 20, (80, 140, 255, 150), width=1)
    draw_ring(draw, cx, cy, 14, (60, 120, 255, 120), width=1)

    # 無限符號（∞）— 用兩個圓圈模擬
    # 左圓
    draw.ellipse([cx-12, cy-5, cx-2, cy+5], outline=(200, 220, 255, 220), width=2)
    # 右圓
    draw.ellipse([cx+2, cy-5, cx+12, cy+5], outline=(200, 220, 255, 220), width=2)

    # 中心高光
    fill_circle(draw, cx, cy, 3, (200, 220, 255, 255))
    fill_circle(draw, cx, cy, 1, (255, 255, 255, 255))

    # 輪廓
    draw_ring(draw, cx, cy, 22, (0, 20, 80, 255), width=2)

    # 放大到 128x128
    img = img.resize((128, 128), Image.NEAREST)
    path = os.path.join(OUT_DIR, "T197_eternal_cycle.png")
    img.save(path)
    print(f"T197 saved: {path}")
    return img

def generate_t198():
    """T198 幸運混沌爆炸魚：深紅色魚身 + 混沌爆炸光芒 + 裂縫紋路"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2

    # 外層光暈（深紅色）
    for r in range(30, 26, -1):
        alpha = int(80 * (r - 26) / 4)
        draw.ellipse([cx-r, cy-r, cx+r, cy+r], fill=(180, 20, 0, alpha))

    # 混沌爆炸光芒（12道，不規則長度）
    import random
    rng = random.Random(198)
    for i in range(12):
        angle = 2 * math.pi * i / 12 + rng.uniform(-0.1, 0.1)
        r_len = rng.randint(24, 30)
        x1 = cx + 20 * math.cos(angle)
        y1 = cy + 20 * math.sin(angle)
        x2 = cx + r_len * math.cos(angle)
        y2 = cy + r_len * math.sin(angle)
        alpha = rng.randint(150, 220)
        draw.line([x1, y1, x2, y2], fill=(255, 100, 0, alpha), width=2)

    # 魚身（深紅色，帶陰影）
    fill_circle_shaded(draw, cx, cy, 20,
        (180, 20, 0, 255),     # 基底深紅
        (255, 80, 40, 255),    # 高光橙紅
        (100, 10, 0, 255))     # 陰影深紅

    # 裂縫紋路（Z字形）
    draw.line([cx-8, cy-8, cx, cy, cx+8, cy-8], fill=(255, 200, 0, 200), width=1)
    draw.line([cx-8, cy+8, cx, cy, cx+8, cy+8], fill=(255, 200, 0, 200), width=1)

    # 中心爆炸圓
    fill_circle(draw, cx, cy, 6, (255, 150, 0, 255))
    fill_circle(draw, cx, cy, 3, (255, 255, 100, 255))
    fill_circle(draw, cx, cy, 1, (255, 255, 255, 255))

    # 輪廓
    draw_ring(draw, cx, cy, 20, (80, 10, 0, 255), width=2)

    # 放大到 128x128
    img = img.resize((128, 128), Image.NEAREST)
    path = os.path.join(OUT_DIR, "T198_chaos_explosion.png")
    img.save(path)
    print(f"T198 saved: {path}")
    return img

def generate_t199():
    """T199 幸運神聖復活魚：深金色魚身 + 神聖光芒 + 鳳凰羽翼"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2

    # 外層光暈（深金色）
    for r in range(30, 26, -1):
        alpha = int(70 * (r - 26) / 4)
        draw.ellipse([cx-r, cy-r, cx+r, cy+r], fill=(200, 180, 0, alpha))

    # 神聖光芒（16道）
    draw_rays(draw, cx, cy, 16, 20, 30, (255, 240, 100, 180), width=1)

    # 鳳凰羽翼（左右各一個弧形）
    draw.arc([cx-28, cy-12, cx-8, cy+12], start=30, end=150, fill=(255, 200, 50, 200), width=2)
    draw.arc([cx+8, cy-12, cx+28, cy+12], start=30, end=150, fill=(255, 200, 50, 200), width=2)

    # 魚身（深金色，帶陰影）
    fill_circle_shaded(draw, cx, cy, 20,
        (200, 160, 0, 255),    # 基底深金
        (255, 220, 60, 255),   # 高光亮金
        (120, 90, 0, 255))     # 陰影深棕金

    # 神聖光環
    draw_ring(draw, cx, cy, 24, (255, 240, 100, 200), width=2)

    # 復活符文（中心十字）
    draw.line([cx-6, cy, cx+6, cy], fill=(255, 255, 200, 220), width=2)
    draw.line([cx, cy-6, cx, cy+6], fill=(255, 255, 200, 220), width=2)

    # 中心高光
    fill_circle(draw, cx, cy, 4, (255, 240, 100, 255))
    fill_circle(draw, cx, cy, 2, (255, 255, 255, 255))

    # 輪廓
    draw_ring(draw, cx, cy, 20, (100, 80, 0, 255), width=2)

    # 放大到 128x128
    img = img.resize((128, 128), Image.NEAREST)
    path = os.path.join(OUT_DIR, "T199_divine_revival.png")
    img.save(path)
    print(f"T199 saved: {path}")
    return img

def generate_t200():
    """T200 幸運創世紀元魚（里程碑）：純黑色大型魚身 + 25 道光芒 + 四層光環 + 里程碑符文"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2

    # 外層光暈（純白色，最高階）
    for r in range(31, 26, -1):
        alpha = int(80 * (r - 26) / 5)
        draw.ellipse([cx-r, cy-r, cx+r, cy+r], fill=(255, 255, 255, alpha))

    # 25 道光芒（里程碑 200 → 25 道）
    draw_rays(draw, cx, cy, 25, 20, 31, (255, 255, 255, 200), width=1)

    # 四層光環
    draw_ring(draw, cx, cy, 30, (255, 255, 255, 180), width=2)
    draw_ring(draw, cx, cy, 26, (255, 220, 100, 160), width=1)
    draw_ring(draw, cx, cy, 22, (255, 255, 255, 140), width=1)
    draw_ring(draw, cx, cy, 18, (255, 200, 50, 120), width=1)

    # 魚身（純黑色，帶白色高光）
    fill_circle_shaded(draw, cx, cy, 16,
        (10, 10, 10, 255),     # 基底純黑
        (60, 60, 60, 255),     # 高光深灰
        (0, 0, 0, 255))        # 陰影純黑

    # 里程碑符文（200 的象徵：兩個圓 + 零）
    # 左圓（代表 2）
    draw.ellipse([cx-12, cy-5, cx-4, cy+5], outline=(255, 220, 0, 220), width=1)
    # 右圓（代表 0）
    draw.ellipse([cx+4, cy-5, cx+12, cy+5], outline=(255, 220, 0, 220), width=1)

    # 中心星形（最高階標誌）
    for i in range(8):
        angle = 2 * math.pi * i / 8
        x = cx + 6 * math.cos(angle)
        y = cy + 6 * math.sin(angle)
        draw.point([int(x), int(y)], fill=(255, 255, 255, 255))

    # 中心高光
    fill_circle(draw, cx, cy, 3, (255, 220, 0, 255))
    fill_circle(draw, cx, cy, 1, (255, 255, 255, 255))

    # 輪廓（白色，最高階）
    draw_ring(draw, cx, cy, 16, (255, 255, 255, 255), width=2)

    # 放大到 128x128
    img = img.resize((128, 128), Image.NEAREST)
    path = os.path.join(OUT_DIR, "T200_genesis_epoch.png")
    img.save(path)
    print(f"T200 saved: {path}")
    return img

if __name__ == "__main__":
    os.makedirs(OUT_DIR, exist_ok=True)
    generate_t196()
    generate_t197()
    generate_t198()
    generate_t199()
    generate_t200()
    print("DAY-318 T196-T200 精靈圖生成完成！")
