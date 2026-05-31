"""
DAY-345: 重新生成損壞的目標物 PNG
T141-T145, T156-T160 的 IDAT 資料損壞，需要重新生成
"""
import os
import numpy as np
from PIL import Image, ImageDraw

BASE = r'd:\Kiro\client\chiikawa-pixel\assets\sprites\targets'
SIZE = 64

def make_canvas():
    return Image.new('RGBA', (SIZE, SIZE), (0, 0, 0, 0))

def draw_circle_shaded(img, cx, cy, r, base_color, light_factor=1.3, dark_factor=0.6):
    """帶陰影的圓形"""
    draw = ImageDraw.Draw(img)
    br, bg, bb = base_color
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                x, y = cx+dx, cy+dy
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    # 左上亮，右下暗
                    t = ((-dx - dy) / (2*r) + 0.5)
                    t = max(0, min(1, t))
                    f = dark_factor + t * (light_factor - dark_factor)
                    r2 = int(min(255, br * f))
                    g2 = int(min(255, bg * f))
                    b2 = int(min(255, bb * f))
                    img.putpixel((x, y), (r2, g2, b2, 255))

def draw_star_rays(img, cx, cy, n, r_inner, r_outer, color):
    """放射狀光芒"""
    import math
    draw = ImageDraw.Draw(img)
    for i in range(n):
        angle = i * 2 * math.pi / n - math.pi/2
        x1 = int(cx + r_inner * math.cos(angle))
        y1 = int(cy + r_inner * math.sin(angle))
        x2 = int(cx + r_outer * math.cos(angle))
        y2 = int(cy + r_outer * math.sin(angle))
        draw.line([(x1, y1), (x2, y2)], fill=color, width=2)

def draw_ring(img, cx, cy, r, color, width=2):
    """圓環"""
    draw = ImageDraw.Draw(img)
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], outline=color, width=width)

def generate_t141_tornado():
    """T141 龍捲風 — 藍紫色螺旋"""
    img = make_canvas()
    import math
    # 螺旋線
    for t in range(0, 400):
        angle = t * 0.1
        r = t * 0.06
        x = int(32 + r * math.cos(angle))
        y = int(32 + r * math.sin(angle))
        if 0 <= x < SIZE and 0 <= y < SIZE:
            alpha = min(255, int(t * 0.6))
            img.putpixel((x, y), (120, 80, 220, alpha))
    # 中心漩渦
    draw_circle_shaded(img, 32, 32, 10, (100, 60, 200))
    # 光芒
    draw_star_rays(img, 32, 32, 8, 12, 28, (160, 120, 255, 200))
    # 外環
    draw_ring(img, 32, 32, 28, (180, 140, 255, 180))
    return img

def generate_t142_earthquake():
    """T142 地震 — 棕褐色裂縫"""
    img = make_canvas()
    draw = ImageDraw.Draw(img)
    # 地面
    draw_circle_shaded(img, 32, 38, 20, (139, 90, 43))
    # 裂縫
    draw.line([(20, 38), (32, 30), (44, 38)], fill=(60, 30, 10, 255), width=3)
    draw.line([(25, 45), (32, 35), (39, 45)], fill=(60, 30, 10, 255), width=2)
    # 震動波
    draw_ring(img, 32, 38, 22, (200, 150, 80, 180))
    draw_ring(img, 32, 38, 26, (200, 150, 80, 120))
    # 光芒
    draw_star_rays(img, 32, 32, 6, 14, 26, (220, 170, 80, 200))
    return img

def generate_t143_volcano():
    """T143 火山 — 橙紅色爆發"""
    img = make_canvas()
    draw = ImageDraw.Draw(img)
    # 火山主體（三角形）
    draw.polygon([(32, 8), (10, 52), (54, 52)], fill=(180, 60, 20, 255))
    draw.polygon([(32, 8), (10, 52), (54, 52)], outline=(220, 100, 40, 255), width=2)
    # 熔岩口
    draw_circle_shaded(img, 32, 14, 8, (255, 120, 0))
    # 熔岩噴發
    for i in range(6):
        import math
        angle = -math.pi/2 + (i - 2.5) * 0.3
        for r in range(5, 20):
            x = int(32 + r * math.cos(angle))
            y = int(14 - r * abs(math.sin(angle)) * 0.8)
            if 0 <= x < SIZE and 0 <= y < SIZE:
                alpha = max(0, 255 - r * 10)
                img.putpixel((x, y), (255, 80 + r*5, 0, alpha))
    # 光芒
    draw_star_rays(img, 32, 32, 8, 16, 30, (255, 140, 0, 180))
    return img

def generate_t144_cosmic_ray():
    """T144 宇宙射線 — 青藍色光束"""
    img = make_canvas()
    import math
    # 中心能量球
    draw_circle_shaded(img, 32, 32, 12, (0, 180, 255))
    # 光束（4方向）
    for angle in [0, math.pi/2, math.pi, 3*math.pi/2]:
        for r in range(14, 30):
            x = int(32 + r * math.cos(angle))
            y = int(32 + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                alpha = max(0, 255 - (r-14) * 15)
                img.putpixel((x, y), (0, 200, 255, alpha))
                # 光束寬度
                for w in [-1, 1]:
                    xw = int(32 + r * math.cos(angle) + w * math.sin(angle))
                    yw = int(32 + r * math.sin(angle) - w * math.cos(angle))
                    if 0 <= xw < SIZE and 0 <= yw < SIZE:
                        img.putpixel((xw, yw), (0, 200, 255, alpha // 2))
    # 外環
    draw_ring(img, 32, 32, 28, (0, 220, 255, 200))
    draw_star_rays(img, 32, 32, 12, 14, 28, (100, 220, 255, 180))
    return img

def generate_t145_divine_dragon():
    """T145 神聖龍 — 金色龍形"""
    img = make_canvas()
    draw = ImageDraw.Draw(img)
    # 龍身（S形曲線）
    import math
    for t in range(0, 100):
        angle = t * 0.1
        x = int(32 + 15 * math.sin(angle))
        y = int(10 + t * 0.44)
        if 0 <= x < SIZE and 0 <= y < SIZE:
            img.putpixel((x, y), (255, 200, 0, 255))
            if x+1 < SIZE: img.putpixel((x+1, y), (255, 180, 0, 200))
            if x-1 >= 0: img.putpixel((x-1, y), (255, 180, 0, 200))
    # 龍頭
    draw_circle_shaded(img, 32, 12, 8, (255, 180, 0))
    # 眼睛
    img.putpixel((29, 10), (255, 50, 0, 255))
    img.putpixel((35, 10), (255, 50, 0, 255))
    # 光芒
    draw_star_rays(img, 32, 32, 10, 16, 30, (255, 220, 80, 200))
    draw_ring(img, 32, 32, 28, (255, 200, 0, 180))
    return img

def generate_t156_ice_phoenix():
    """T156 冰鳳凰 — 冰藍色鳳凰"""
    img = make_canvas()
    draw = ImageDraw.Draw(img)
    import math
    # 鳳凰身體
    draw_circle_shaded(img, 32, 28, 10, (100, 200, 255))
    # 翅膀（左右各一）
    for side in [-1, 1]:
        for i in range(15):
            angle = math.pi/2 + side * (0.3 + i * 0.1)
            r = 8 + i * 1.2
            x = int(32 + side * r * math.cos(angle - math.pi/2))
            y = int(28 - r * 0.5)
            if 0 <= x < SIZE and 0 <= y < SIZE:
                alpha = max(0, 255 - i * 15)
                img.putpixel((x, y), (150, 220, 255, alpha))
    # 尾羽
    for i in range(5):
        angle = math.pi/2 + (i - 2) * 0.25
        for r in range(10, 25):
            x = int(32 + r * math.cos(angle + math.pi/2))
            y = int(28 + r * math.sin(angle + math.pi/2))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                alpha = max(0, 200 - r * 6)
                img.putpixel((x, y), (180, 230, 255, alpha))
    # 冰晶光環
    draw_ring(img, 32, 28, 20, (200, 240, 255, 200))
    draw_star_rays(img, 32, 28, 8, 12, 26, (200, 240, 255, 180))
    return img

def generate_t157_dragon_fury():
    """T157 龍之怒 — 深紅色龍爪"""
    img = make_canvas()
    draw = ImageDraw.Draw(img)
    import math
    # 龍爪（3個爪子）
    for i in range(3):
        angle = -math.pi/2 + (i - 1) * 0.5
        for r in range(5, 25):
            x = int(32 + r * math.cos(angle))
            y = int(32 + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                alpha = max(0, 255 - r * 8)
                img.putpixel((x, y), (220, 30, 30, alpha))
                if r < 15:
                    for w in [-1, 1]:
                        xw = int(x + w * math.sin(angle))
                        yw = int(y - w * math.cos(angle))
                        if 0 <= xw < SIZE and 0 <= yw < SIZE:
                            img.putpixel((xw, yw), (200, 20, 20, alpha // 2))
    # 中心火焰
    draw_circle_shaded(img, 32, 32, 10, (200, 50, 0))
    # 光芒
    draw_star_rays(img, 32, 32, 8, 12, 28, (255, 80, 0, 200))
    draw_ring(img, 32, 32, 28, (255, 60, 0, 180))
    return img

def generate_t158_mult_cascade():
    """T158 倍率連鎖 — 金色連鎖符號"""
    img = make_canvas()
    draw = ImageDraw.Draw(img)
    import math
    # 連鎖圓環（3個）
    positions = [(20, 28), (32, 20), (44, 28)]
    for cx, cy in positions:
        draw_circle_shaded(img, cx, cy, 8, (255, 180, 0))
        draw_ring(img, cx, cy, 10, (255, 220, 80, 200))
    # 連接線
    draw.line([(28, 28), (24, 24)], fill=(255, 200, 0, 255), width=3)
    draw.line([(36, 24), (40, 28)], fill=(255, 200, 0, 255), width=3)
    # 倍率符號 ×
    draw.line([(28, 40), (36, 48)], fill=(255, 220, 0, 255), width=3)
    draw.line([(36, 40), (28, 48)], fill=(255, 220, 0, 255), width=3)
    # 光芒
    draw_star_rays(img, 32, 32, 10, 16, 30, (255, 200, 0, 180))
    draw_ring(img, 32, 32, 28, (255, 180, 0, 160))
    return img

def generate_t159_awaken_boss_v2():
    """T159 覺醒BOSS v2 — 深紫色覺醒形態"""
    img = make_canvas()
    draw = ImageDraw.Draw(img)
    import math
    # 覺醒光環（多層）
    for r, alpha in [(28, 80), (24, 120), (20, 160)]:
        draw_ring(img, 32, 32, r, (180, 0, 255, alpha), width=3)
    # BOSS 主體
    draw_circle_shaded(img, 32, 32, 14, (120, 0, 200))
    # 覺醒眼睛（3個）
    for i in range(3):
        angle = -math.pi/2 + (i - 1) * 0.8
        ex = int(32 + 8 * math.cos(angle))
        ey = int(32 + 8 * math.sin(angle))
        img.putpixel((ex, ey), (255, 0, 255, 255))
        img.putpixel((ex+1, ey), (255, 0, 255, 255))
    # 能量光芒
    draw_star_rays(img, 32, 32, 12, 16, 30, (200, 0, 255, 200))
    return img

def generate_t160_ultimate_judgment():
    """T160 終極審判 — 白金色神聖光柱"""
    img = make_canvas()
    draw = ImageDraw.Draw(img)
    import math
    # 神聖光柱（12道）
    for i in range(12):
        angle = i * 2 * math.pi / 12
        for r in range(14, 30):
            x = int(32 + r * math.cos(angle))
            y = int(32 + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                alpha = max(0, 200 - (r-14) * 12)
                img.putpixel((x, y), (255, 240, 180, alpha))
    # 中心神聖球
    draw_circle_shaded(img, 32, 32, 12, (255, 220, 100))
    # 十字光芒
    for angle in [0, math.pi/4, math.pi/2, 3*math.pi/4]:
        draw.line([
            (int(32 + 14*math.cos(angle)), int(32 + 14*math.sin(angle))),
            (int(32 + 30*math.cos(angle)), int(32 + 30*math.sin(angle)))
        ], fill=(255, 255, 200, 220), width=3)
    # 外環
    draw_ring(img, 32, 32, 28, (255, 240, 150, 200), width=2)
    draw_ring(img, 32, 32, 24, (255, 220, 100, 160), width=2)
    return img

# 生成並儲存
generators = {
    'T141_tornado.png': generate_t141_tornado,
    'T142_earthquake.png': generate_t142_earthquake,
    'T143_volcano.png': generate_t143_volcano,
    'T144_cosmic_ray.png': generate_t144_cosmic_ray,
    'T145_divine_dragon.png': generate_t145_divine_dragon,
    'T156_ice_phoenix.png': generate_t156_ice_phoenix,
    'T157_dragon_fury.png': generate_t157_dragon_fury,
    'T158_mult_cascade.png': generate_t158_mult_cascade,
    'T159_awaken_boss_v2.png': generate_t159_awaken_boss_v2,
    'T160_ultimate_judgment.png': generate_t160_ultimate_judgment,
}

print(f"重新生成 {len(generators)} 個損壞的目標物...")
for filename, gen_func in generators.items():
    path = os.path.join(BASE, filename)
    img = gen_func()
    img.save(path)
    pixels = sum(1 for px in img.getdata() if px[3] > 50)
    print(f"  ✅ {filename}: {pixels} 非透明像素")

print("\n完成！所有損壞的目標物已重新生成。")
