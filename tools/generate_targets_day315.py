"""
generate_targets_day315.py — DAY-315 T181-T185 精靈圖生成
業界依據：Fisch mutations / Arctic Mechanics / Big Bass Splash / BGaming Fishing Club 2 / Fishing Fortune
"""
import os
from PIL import Image, ImageDraw
import math
import random

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def fill_circle(draw, cx, cy, r, color):
    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            if (x - cx) ** 2 + (y - cy) ** 2 <= r ** 2:
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    draw.point((x, y), fill=color)

def fill_circle_shaded(draw, cx, cy, r, base_color, light_color, dark_color):
    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            if (x - cx) ** 2 + (y - cy) ** 2 <= r ** 2:
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    if x < cx and y < cy:
                        draw.point((x, y), fill=light_color)
                    elif x > cx and y > cy:
                        draw.point((x, y), fill=dark_color)
                    else:
                        draw.point((x, y), fill=base_color)

def draw_eye(draw, cx, cy):
    fill_circle(draw, cx, cy, 3, (255, 255, 255, 255))
    fill_circle(draw, cx, cy, 2, (30, 30, 30, 255))
    draw.point((cx - 1, cy - 1), fill=(255, 255, 255, 255))

def draw_rays(draw, cx, cy, n_rays, length, color, width=1):
    for i in range(n_rays):
        angle = 2 * math.pi * i / n_rays
        x2 = int(cx + length * math.cos(angle))
        y2 = int(cy + length * math.sin(angle))
        draw.line([(cx, cy), (x2, y2)], fill=color, width=width)

def count_pixels(img):
    pixels = img.load()
    count = 0
    for y in range(SIZE):
        for x in range(SIZE):
            if pixels[x, y][3] > 0:
                count += 1
    return count

# ── T181 幸運突變魚 ──────────────────────────────────────────
def gen_t181():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 深紫橢圓魚身（突變主題）
    BODY = (120, 30, 180, 255)
    BODY_L = (160, 60, 220, 255)
    BODY_D = (80, 10, 130, 255)
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - 32) / 22.0
            dy = (y - 32) / 16.0
            if dx * dx + dy * dy <= 1.0:
                if x < 32 and y < 32:
                    draw.point((x, y), fill=BODY_L)
                elif x > 32 and y > 32:
                    draw.point((x, y), fill=BODY_D)
                else:
                    draw.point((x, y), fill=BODY)
    
    # DNA 螺旋紋路（突變特徵）
    rng = random.Random(181)
    for i in range(20):
        angle = i * 0.6
        x = int(32 + 15 * math.cos(angle))
        y = int(32 + 8 * math.sin(angle * 2))
        if 0 <= x < SIZE and 0 <= y < SIZE:
            draw.point((x, y), fill=(200, 100, 255, 200))
    
    # 突變光環（多色）
    for i in range(16):
        angle = 2 * math.pi * i / 16
        r = 28
        x = int(32 + r * math.cos(angle))
        y = int(32 + r * math.sin(angle))
        hue = i / 16.0
        r_c = int(128 + 127 * math.sin(hue * 2 * math.pi))
        g_c = int(128 + 127 * math.sin(hue * 2 * math.pi + 2.094))
        b_c = int(128 + 127 * math.sin(hue * 2 * math.pi + 4.189))
        if 0 <= x < SIZE and 0 <= y < SIZE:
            draw.point((x, y), fill=(r_c, g_c, b_c, 180))
    
    # 眼睛
    draw_eye(draw, 40, 28)
    
    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 0:
                for dx2, dy2 in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x + dx2, y + dy2
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] == 0:
                        draw.point((nx, ny), fill=(50, 0, 80, 200))
    
    return img

# ── T182 幸運北極風暴魚 ──────────────────────────────────────
def gen_t182():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 冰藍橢圓魚身（北極主題）
    BODY = (0, 150, 220, 255)
    BODY_L = (100, 200, 255, 255)
    BODY_D = (0, 80, 160, 255)
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - 32) / 22.0
            dy = (y - 32) / 16.0
            if dx * dx + dy * dy <= 1.0:
                if x < 32 and y < 32:
                    draw.point((x, y), fill=BODY_L)
                elif x > 32 and y > 32:
                    draw.point((x, y), fill=BODY_D)
                else:
                    draw.point((x, y), fill=BODY)
    
    # 冰晶紋路（8 個冰晶）
    for i in range(8):
        angle = 2 * math.pi * i / 8
        r = 18
        x = int(32 + r * math.cos(angle))
        y = int(32 + r * math.sin(angle))
        fill_circle(draw, x, y, 2, (200, 240, 255, 200))
    
    # 8 波射線（北極風暴特徵）
    draw_rays(draw, 32, 32, 8, 28, (150, 220, 255, 160), 1)
    
    # 雪花符號（中心）
    for angle in [0, math.pi/3, 2*math.pi/3]:
        x1 = int(32 + 8 * math.cos(angle))
        y1 = int(32 + 8 * math.sin(angle))
        x2 = int(32 - 8 * math.cos(angle))
        y2 = int(32 - 8 * math.sin(angle))
        draw.line([(x1, y1), (x2, y2)], fill=(255, 255, 255, 220), width=1)
    
    # 眼睛
    draw_eye(draw, 40, 28)
    
    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 0:
                for dx2, dy2 in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x + dx2, y + dy2
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] == 0:
                        draw.point((nx, ny), fill=(0, 50, 120, 200))
    
    return img

# ── T183 幸運漁夫野生魚 ──────────────────────────────────────
def gen_t183():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 深藍橢圓魚身（漁夫主題）
    BODY = (20, 80, 180, 255)
    BODY_L = (60, 130, 220, 255)
    BODY_D = (10, 40, 120, 255)
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - 32) / 22.0
            dy = (y - 32) / 16.0
            if dx * dx + dy * dy <= 1.0:
                if x < 32 and y < 32:
                    draw.point((x, y), fill=BODY_L)
                elif x > 32 and y > 32:
                    draw.point((x, y), fill=BODY_D)
                else:
                    draw.point((x, y), fill=BODY)
    
    # Wild 標記（3 個金色星星）
    for i, (sx, sy) in enumerate([(20, 20), (44, 20), (32, 44)]):
        fill_circle(draw, sx, sy, 4, (255, 200, 0, 220))
        draw.text((sx - 3, sy - 4), "W", fill=(255, 255, 255, 255))
    
    # 釣魚線（漁夫特徵）
    draw.line([(32, 10), (32, 32)], fill=(200, 200, 200, 180), width=1)
    fill_circle(draw, 32, 10, 3, (255, 200, 0, 200))
    
    # 眼睛
    draw_eye(draw, 40, 28)
    
    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 0:
                for dx2, dy2 in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x + dx2, y + dy2
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] == 0:
                        draw.point((nx, ny), fill=(0, 20, 80, 200))
    
    return img

# ── T184 幸運風險等級魚 ──────────────────────────────────────
def gen_t184():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 火橙橢圓魚身（風險主題）
    BODY = (220, 80, 0, 255)
    BODY_L = (255, 140, 50, 255)
    BODY_D = (160, 40, 0, 255)
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - 32) / 22.0
            dy = (y - 32) / 16.0
            if dx * dx + dy * dy <= 1.0:
                if x < 32 and y < 32:
                    draw.point((x, y), fill=BODY_L)
                elif x > 32 and y > 32:
                    draw.point((x, y), fill=BODY_D)
                else:
                    draw.point((x, y), fill=BODY)
    
    # 5 個風險等級指示器（從綠到紅）
    level_colors = [
        (0, 200, 0, 200),    # 低風險 - 綠
        (255, 165, 0, 200),  # 中風險 - 橙
        (255, 50, 50, 200),  # 高風險 - 紅
        (150, 0, 200, 200),  # 極高 - 紫
        (255, 215, 0, 220),  # 最高 - 金
    ]
    for i, color in enumerate(level_colors):
        x = 12 + i * 10
        y = 50
        fill_circle(draw, x, y, 3, color)
    
    # 骰子符號（風險特徵）
    draw.rectangle([(24, 22), (40, 38)], outline=(255, 255, 255, 200), width=1)
    for dot_x, dot_y in [(27, 25), (37, 25), (32, 30), (27, 35), (37, 35)]:
        draw.point((dot_x, dot_y), fill=(255, 255, 255, 220))
    
    # 眼睛
    draw_eye(draw, 40, 28)
    
    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 0:
                for dx2, dy2 in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x + dx2, y + dy2
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] == 0:
                        draw.point((nx, ny), fill=(100, 20, 0, 200))
    
    return img

# ── T185 幸運宇宙脈衝魚 ──────────────────────────────────────
def gen_t185():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 深紫大型圓形魚身（宇宙主題）
    BODY = (60, 0, 120, 255)
    BODY_L = (120, 40, 200, 255)
    BODY_D = (30, 0, 80, 255)
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - 32) / 26.0
            dy = (y - 32) / 22.0
            if dx * dx + dy * dy <= 1.0:
                if x < 32 and y < 32:
                    draw.point((x, y), fill=BODY_L)
                elif x > 32 and y > 32:
                    draw.point((x, y), fill=BODY_D)
                else:
                    draw.point((x, y), fill=BODY)
    
    # 16 道宇宙脈衝光芒（超越 T180 的 12 道）
    draw_rays(draw, 32, 32, 16, 30, (180, 100, 255, 180), 1)
    
    # 三層光環
    for r, alpha in [(20, 120), (25, 80), (30, 50)]:
        for i in range(36):
            angle = 2 * math.pi * i / 36
            x = int(32 + r * math.cos(angle))
            y = int(32 + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=(200, 150, 255, alpha))
    
    # 宇宙粒子（散落）
    rng = random.Random(185)
    for _ in range(25):
        x = rng.randint(4, 59)
        y = rng.randint(4, 59)
        if img.getpixel((x, y))[3] > 0:
            draw.point((x, y), fill=(255, 255, 255, 180))
    
    # 中心爆炸圓
    fill_circle(draw, 32, 32, 8, (220, 180, 255, 220))
    fill_circle(draw, 32, 32, 4, (255, 255, 255, 240))
    
    # 眼睛
    draw_eye(draw, 40, 28)
    
    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 0:
                for dx2, dy2 in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x + dx2, y + dy2
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and img.getpixel((nx, ny))[3] == 0:
                        draw.point((nx, ny), fill=(30, 0, 60, 200))
    
    return img

def main():
    os.makedirs(OUT_DIR, exist_ok=True)
    
    targets = [
        ("T181_mutation", gen_t181),
        ("T182_arctic_storm", gen_t182),
        ("T183_fisher_wild", gen_t183),
        ("T184_risk_level", gen_t184),
        ("T185_cosmic_pulse", gen_t185),
    ]
    
    for name, gen_func in targets:
        img = gen_func()
        path = os.path.join(OUT_DIR, f"{name}.png")
        img.save(path)
        pixels = count_pixels(img)
        pct = pixels / (SIZE * SIZE) * 100
        print(f"  {name}: {pixels} 非透明像素 ({pct:.1f}%)")
    
    print(f"\n✅ DAY-315 T181-T185 精靈圖生成完成！")

if __name__ == "__main__":
    main()
