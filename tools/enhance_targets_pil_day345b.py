"""
DAY-345b: T161-T190 美術升級 + 損壞 PNG 修復
T161-T165 損壞需要重新生成
T161-T190 全部升級：飽和度+45%、對比度+35%、亮度+10%、三重銳化
光暈策略：
- T161-T170: 深紫/洋紅光暈（超高階神秘感）
- T171-T180: 金橙色光暈（Jackpot 感）
- T181-T190: 白銀色光暈（極致感）
"""
import os
import shutil
import math
from PIL import Image, ImageEnhance, ImageFilter, ImageDraw
import numpy as np

BASE = r'd:\Kiro\client\chiikawa-pixel\assets\sprites\targets'
BACKUP = r'd:\Kiro\tmp\targets_backup_day345b'
SIZE = 64

os.makedirs(BACKUP, exist_ok=True)

# ── 修復損壞的 T161-T165 ──────────────────────────────────────

def make_canvas():
    return Image.new('RGBA', (SIZE, SIZE), (0, 0, 0, 0))

def draw_circle_shaded(img, cx, cy, r, base_color, light_factor=1.3, dark_factor=0.6):
    br, bg, bb = base_color
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                x, y = cx+dx, cy+dy
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    t = ((-dx - dy) / (2*r) + 0.5)
                    t = max(0, min(1, t))
                    f = dark_factor + t * (light_factor - dark_factor)
                    r2 = int(min(255, br * f))
                    g2 = int(min(255, bg * f))
                    b2 = int(min(255, bb * f))
                    img.putpixel((x, y), (r2, g2, b2, 255))

def draw_star_rays(img, cx, cy, n, r_inner, r_outer, color):
    draw = ImageDraw.Draw(img)
    for i in range(n):
        angle = i * 2 * math.pi / n - math.pi/2
        x1 = int(cx + r_inner * math.cos(angle))
        y1 = int(cy + r_inner * math.sin(angle))
        x2 = int(cx + r_outer * math.cos(angle))
        y2 = int(cy + r_outer * math.sin(angle))
        draw.line([(x1, y1), (x2, y2)], fill=color, width=2)

def draw_ring(img, cx, cy, r, color, width=2):
    draw = ImageDraw.Draw(img)
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], outline=color, width=width)

def generate_t161_combo_burst():
    """T161 連擊爆發 — 橙色爆炸"""
    img = make_canvas()
    draw = ImageDraw.Draw(img)
    # 爆炸中心
    draw_circle_shaded(img, 32, 32, 12, (255, 120, 0))
    # 爆炸碎片（8方向）
    for i in range(8):
        angle = i * math.pi / 4
        for r in range(14, 26):
            x = int(32 + r * math.cos(angle))
            y = int(32 + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                alpha = max(0, 255 - (r-14) * 20)
                img.putpixel((x, y), (255, 100 + r*3, 0, alpha))
    # 連擊數字 "x"
    draw.line([(26, 26), (38, 38)], fill=(255, 220, 0, 255), width=3)
    draw.line([(38, 26), (26, 38)], fill=(255, 220, 0, 255), width=3)
    draw_ring(img, 32, 32, 28, (255, 150, 0, 200))
    draw_star_rays(img, 32, 32, 12, 14, 28, (255, 180, 0, 180))
    return img

def generate_t162_time_bomb():
    """T162 時間炸彈 — 紅色炸彈"""
    img = make_canvas()
    draw = ImageDraw.Draw(img)
    # 炸彈主體
    draw_circle_shaded(img, 32, 36, 14, (180, 30, 30))
    # 導火線
    draw.line([(32, 22), (38, 14)], fill=(200, 150, 50, 255), width=3)
    # 火花
    for i in range(5):
        angle = -math.pi/4 + i * 0.3
        x = int(38 + 4 * math.cos(angle))
        y = int(14 + 4 * math.sin(angle))
        if 0 <= x < SIZE and 0 <= y < SIZE:
            img.putpixel((x, y), (255, 200, 0, 255))
    # 倒數符號
    draw.text((26, 30), "!", fill=(255, 255, 0, 255))
    draw_ring(img, 32, 36, 16, (255, 60, 0, 200))
    draw_star_rays(img, 32, 32, 8, 18, 30, (255, 80, 0, 180))
    return img

def generate_t163_elemental_fusion():
    """T163 元素融合 — 四色融合"""
    img = make_canvas()
    # 四象限顏色
    colors = [(255, 50, 50), (50, 150, 255), (50, 220, 50), (255, 200, 0)]
    for qi, (qx, qy) in enumerate([(0, 0), (32, 0), (0, 32), (32, 32)]):
        cr, cg, cb = colors[qi]
        for y in range(qy, qy+32):
            for x in range(qx, qx+32):
                dx, dy = x - 32, y - 32
                if dx*dx + dy*dy <= 28*28:
                    dist = math.sqrt(dx*dx + dy*dy)
                    alpha = int(200 * (1 - dist/28))
                    img.putpixel((x, y), (cr, cg, cb, alpha))
    # 中心融合球
    draw_circle_shaded(img, 32, 32, 8, (255, 255, 255))
    draw_ring(img, 32, 32, 28, (255, 255, 255, 180))
    draw_star_rays(img, 32, 32, 8, 10, 28, (255, 255, 255, 160))
    return img

def generate_t164_treasure_hunter():
    """T164 寶藏獵人 — 金色寶箱"""
    img = make_canvas()
    draw = ImageDraw.Draw(img)
    # 寶箱主體
    draw.rectangle([14, 28, 50, 52], fill=(180, 120, 20, 255), outline=(220, 160, 40, 255), width=2)
    # 寶箱蓋
    draw.rectangle([14, 18, 50, 30], fill=(200, 140, 30, 255), outline=(240, 180, 50, 255), width=2)
    # 鎖
    draw_circle_shaded(img, 32, 30, 4, (255, 200, 0))
    # 金幣溢出
    for i in range(6):
        angle = -math.pi/2 + (i - 2.5) * 0.4
        x = int(32 + 18 * math.cos(angle))
        y = int(20 + 10 * math.sin(angle))
        if 0 <= x < SIZE and 0 <= y < SIZE:
            draw_circle_shaded(img, x, y, 3, (255, 200, 0))
    draw_star_rays(img, 32, 32, 10, 16, 30, (255, 200, 0, 180))
    return img

def generate_t165_myth_awaken():
    """T165 神話覺醒 — 彩虹神話"""
    img = make_canvas()
    # 彩虹光環
    for r in range(28, 14, -2):
        hue = (28 - r) / 14.0
        hi = int(hue * 6)
        f = hue * 6 - hi
        if hi == 0: rc, gc, bc = 1, f, 0
        elif hi == 1: rc, gc, bc = 1-f, 1, 0
        elif hi == 2: rc, gc, bc = 0, 1, f
        elif hi == 3: rc, gc, bc = 0, 1-f, 1
        elif hi == 4: rc, gc, bc = f, 0, 1
        else: rc, gc, bc = 1, 0, 1-f
        draw_ring(img, 32, 32, r, (int(rc*255), int(gc*255), int(bc*255), 200))
    # 中心神話球
    draw_circle_shaded(img, 32, 32, 10, (200, 100, 255))
    draw_star_rays(img, 32, 32, 12, 12, 28, (255, 200, 255, 200))
    return img

# 修復損壞的 T161-T165
broken_generators = {
    'T161_combo_burst.png': generate_t161_combo_burst,
    'T162_time_bomb.png': generate_t162_time_bomb,
    'T163_elemental_fusion.png': generate_t163_elemental_fusion,
    'T164_treasure_hunter.png': generate_t164_treasure_hunter,
    'T165_myth_awaken.png': generate_t165_myth_awaken,
}

print("修復損壞的 T161-T165...")
for filename, gen_func in broken_generators.items():
    path = os.path.join(BASE, filename)
    img = gen_func()
    img.save(path)
    pixels = sum(1 for px in img.getdata() if px[3] > 50)
    print(f"  ✅ {filename}: {pixels} 非透明像素")

# ── 升級 T161-T190 ──────────────────────────────────────────

targets = []
for f in os.listdir(BASE):
    if not f.endswith('.png') or f.endswith('.import'):
        continue
    if not f.startswith('T'):
        continue
    try:
        num = int(f[1:4])
        if 161 <= num <= 190:
            targets.append(f)
    except ValueError:
        continue

targets.sort()
print(f"\n找到 {len(targets)} 個目標物（T161-T190）")

def add_ultra_tier_glow(img: Image.Image, target_num: int) -> Image.Image:
    """
    超高階光暈效果：
    T161-T170: 深紫/洋紅光暈（超高階神秘感）
    T171-T180: 金橙色光暈（Jackpot 感）
    T181-T190: 白銀色光暈（極致感）
    """
    arr = np.array(img, dtype=np.float32)
    alpha = arr[:, :, 3]
    mask = alpha > 50

    if not mask.any():
        return img

    if 161 <= target_num <= 170:
        # 深紫/洋紅：R+10%, B+15%, G-8%
        arr[:, :, 0] = np.where(mask, np.clip(arr[:, :, 0] * 1.10, 0, 255), arr[:, :, 0])
        arr[:, :, 1] = np.where(mask, np.clip(arr[:, :, 1] * 0.92, 0, 255), arr[:, :, 1])
        arr[:, :, 2] = np.where(mask, np.clip(arr[:, :, 2] * 1.15, 0, 255), arr[:, :, 2])
    elif 171 <= target_num <= 180:
        # 金橙色：R+15%, G+8%, B-12%
        arr[:, :, 0] = np.where(mask, np.clip(arr[:, :, 0] * 1.15, 0, 255), arr[:, :, 0])
        arr[:, :, 1] = np.where(mask, np.clip(arr[:, :, 1] * 1.08, 0, 255), arr[:, :, 1])
        arr[:, :, 2] = np.where(mask, np.clip(arr[:, :, 2] * 0.88, 0, 255), arr[:, :, 2])
    elif 181 <= target_num <= 190:
        # 白銀色：均勻提亮 +10%，輕微去飽和
        arr[:, :, 0] = np.where(mask, np.clip(arr[:, :, 0] * 1.10, 0, 255), arr[:, :, 0])
        arr[:, :, 1] = np.where(mask, np.clip(arr[:, :, 1] * 1.10, 0, 255), arr[:, :, 1])
        arr[:, :, 2] = np.where(mask, np.clip(arr[:, :, 2] * 1.10, 0, 255), arr[:, :, 2])

    return Image.fromarray(arr.astype(np.uint8), 'RGBA')

def enhance_ultra_tier_target(img: Image.Image, target_name: str) -> Image.Image:
    if img.mode != 'RGBA':
        img = img.convert('RGBA')

    r, g, b, a = img.split()
    rgb = Image.merge('RGB', (r, g, b))

    # 飽和度 +45%
    rgb = ImageEnhance.Color(rgb).enhance(1.45)
    # 對比度 +35%
    rgb = ImageEnhance.Contrast(rgb).enhance(1.35)
    # 亮度 +10%
    rgb = ImageEnhance.Brightness(rgb).enhance(1.10)
    # 三重銳化
    rgb = rgb.filter(ImageFilter.SHARPEN)
    rgb = rgb.filter(ImageFilter.SHARPEN)
    rgb = rgb.filter(ImageFilter.SHARPEN)

    r2, g2, b2 = rgb.split()
    result = Image.merge('RGBA', (r2, g2, b2, a))

    try:
        num = int(target_name[1:4])
        result = add_ultra_tier_glow(result, num)
    except ValueError:
        pass

    return result

enhanced_count = 0
for filename in targets:
    src_path = os.path.join(BASE, filename)
    backup_path = os.path.join(BACKUP, filename)

    shutil.copy2(src_path, backup_path)

    img = Image.open(src_path).convert('RGBA')
    original_pixels = sum(1 for px in img.getdata() if px[3] > 50)

    enhanced = enhance_ultra_tier_target(img, filename)
    enhanced_pixels = sum(1 for px in enhanced.getdata() if px[3] > 50)

    enhanced.save(src_path)
    enhanced_count += 1

    num = int(filename[1:4])
    tier = "深紫/洋紅" if num <= 170 else ("金橙" if num <= 180 else "白銀")
    print(f"  ✅ {filename} [{tier}光暈]: {original_pixels} → {enhanced_pixels} 像素")

print(f"\n完成！升級了 {enhanced_count} 個超高階 Lucky 目標物（T161-T190）")
print(f"備份位置：{BACKUP}")
print(f"策略：飽和度+45%、對比度+35%、亮度+10%、三重銳化")
print(f"T161-T170：深紫/洋紅光暈（超高階神秘感）")
print(f"T171-T180：金橙色光暈（Jackpot 感）")
print(f"T181-T190：白銀色光暈（極致感）")
