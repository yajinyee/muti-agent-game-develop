# -*- coding: utf-8 -*-
"""
特效 v2 - 更大更清晰的命中特效和投射物
命中特效：48x48（原本 24x24）
投射物：24x12（原本 12x8）
死亡粒子：48x48（原本 32x32）
"""
from PIL import Image
import os
import math

EFFECTS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\effects"
os.makedirs(EFFECTS_DIR, exist_ok=True)

def px(img, x, y, c):
    if 0 <= x < img.width and 0 <= y < img.height:
        img.putpixel((x, y), c)

def fill_circle(img, cx, cy, r, color):
    for y in range(max(0, cy-r), min(img.height, cy+r+1)):
        for x in range(max(0, cx-r), min(img.width, cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 <= r**2:
                px(img, x, y, color)

def gen_hit_effect(char_id):
    """
    命中特效 48x48
    - 中心爆炸圓
    - 放射狀光線（8方向）
    - 角色顏色
    """
    SIZE = 48
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    cx, cy = SIZE//2, SIZE//2

    # 顏色設定
    colors = {
        "chiikawa":  [(255, 160, 210), (255, 200, 230), (255, 100, 170)],
        "hachiware": [(100, 160, 255), (160, 200, 255), (60, 120, 220)],
        "usagi":     [(255, 230, 60),  (255, 245, 140), (220, 180, 20)],
    }
    CORE, LIGHT, DARK = colors.get(char_id, colors["chiikawa"])

    # 外圈光暈（半透明）
    for r in range(20, 24):
        for angle in range(0, 360, 3):
            rad = math.radians(angle)
            x = int(cx + r * math.cos(rad))
            y = int(cy + r * math.sin(rad))
            alpha = max(0, 100 - (r-20)*25)
            px(img, x, y, (*LIGHT, alpha))

    # 放射狀光線（8方向）
    for angle in range(0, 360, 45):
        rad = math.radians(angle)
        for r in range(8, 22):
            x = int(cx + r * math.cos(rad))
            y = int(cy + r * math.sin(rad))
            alpha = max(0, 220 - (r-8)*12)
            width = max(1, 3 - r//8)
            for dw in range(-width, width+1):
                px_x = int(cx + r * math.cos(rad) + dw * math.sin(rad))
                px_y = int(cy + r * math.sin(rad) - dw * math.cos(rad))
                px(img, px_x, px_y, (*CORE, alpha))

    # 中心爆炸圓
    fill_circle(img, cx, cy, 8, (*CORE, 255))
    fill_circle(img, cx, cy, 5, (*LIGHT, 255))
    fill_circle(img, cx, cy, 2, (255, 255, 255, 255))

    # 星形光芒（4方向，短）
    for angle in range(0, 360, 90):
        rad = math.radians(angle + 22)
        for r in range(9, 18):
            x = int(cx + r * math.cos(rad))
            y = int(cy + r * math.sin(rad))
            alpha = max(0, 200 - (r-9)*22)
            px(img, x, y, (255, 255, 255, alpha))

    return img


def gen_projectile(char_id):
    """
    投射物 32x16
    - 橢圓形彈丸
    - 尾焰
    - 角色顏色
    """
    W, H = 32, 16
    img = Image.new("RGBA", (W, H), (0, 0, 0, 0))
    cx, cy = W//2, H//2

    colors = {
        "chiikawa":  [(255, 140, 190), (255, 200, 230), (255, 80, 150)],
        "hachiware": [(80, 140, 240),  (160, 200, 255), (40, 100, 200)],
        "usagi":     [(255, 220, 40),  (255, 245, 140), (200, 160, 10)],
    }
    CORE, LIGHT, DARK = colors.get(char_id, colors["chiikawa"])

    # 尾焰（左側，漸淡）
    for x in range(0, 14):
        t = x / 14
        alpha = int(180 * t)
        r_tail = int(5 * t)
        for dy in range(-r_tail, r_tail+1):
            a = max(0, alpha - abs(dy)*20)
            px(img, x, cy+dy, (*DARK, a))

    # 彈丸主體（橢圓）
    for y in range(H):
        for x in range(10, W):
            if ((x-cx)/10)**2 + ((y-cy)/5)**2 <= 1.0:
                nx_ = (x-cx)/10
                ny_ = (y-cy)/5
                dot = -(nx_*(-0.7)+ny_*(-0.7))
                c = LIGHT if dot > 0.3 else (DARK if dot < -0.1 else CORE)
                px(img, x, y, (*c, 255))

    # 彈丸輪廓
    for y in range(H):
        for x in range(10, W):
            d = ((x-cx)/10)**2 + ((y-cy)/5)**2
            if 0.85 <= d <= 1.1:
                px(img, x, y, (*DARK, 255))

    # 前端高光
    px(img, W-4, cy-1, (255, 255, 255, 200))
    px(img, W-5, cy-2, (255, 255, 255, 150))

    return img


def gen_death_particles():
    """
    死亡粒子 Spritesheet 48x48
    - 金色/黃色粒子爆炸
    - 多個粒子散射
    """
    SIZE = 48
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    cx, cy = SIZE//2, SIZE//2

    import random
    rng = random.Random(42)

    # 粒子（8個，不同方向）
    COLORS = [
        (255, 220, 50, 255),   # 金色
        (255, 180, 30, 255),   # 橙金
        (255, 255, 100, 255),  # 亮黃
        (255, 140, 20, 255),   # 深橙
        (255, 255, 255, 200),  # 白色
    ]

    for i in range(8):
        angle = i * 45 + rng.randint(-15, 15)
        rad = math.radians(angle)
        dist = rng.randint(12, 20)
        px_x = int(cx + dist * math.cos(rad))
        px_y = int(cy + dist * math.sin(rad))
        color = COLORS[i % len(COLORS)]
        r_p = rng.randint(2, 5)
        fill_circle(img, px_x, px_y, r_p, color)
        # 高光
        px(img, px_x-1, px_y-1, (255, 255, 255, 180))

    # 中心爆炸
    fill_circle(img, cx, cy, 6, (255, 200, 50, 255))
    fill_circle(img, cx, cy, 3, (255, 255, 150, 255))
    px(img, cx, cy, (255, 255, 255, 255))

    # 小碎片（4個）
    for i in range(4):
        angle = i * 90 + 22
        rad = math.radians(angle)
        dist = rng.randint(6, 10)
        sx = int(cx + dist * math.cos(rad))
        sy = int(cy + dist * math.sin(rad))
        px(img, sx, sy, (255, 240, 80, 255))
        px(img, sx+1, sy, (255, 200, 40, 200))

    return img


def gen_warning_effect():
    """
    WARNING 特效 128x64
    - 紅色閃爍警告
    - 用於 BOSS 警告
    """
    W, H = 128, 64
    img = Image.new("RGBA", (W, H), (0, 0, 0, 0))

    # 背景（半透明紅）
    for y in range(H):
        for x in range(W):
            alpha = int(180 * (1 - abs(x - W//2) / (W//2)) * (1 - abs(y - H//2) / (H//2)))
            img.putpixel((x, y), (200, 20, 20, alpha))

    # WARNING 文字（像素字體模擬）
    # W
    pixels_W = [(2,4),(2,5),(2,6),(2,7),(2,8),(3,8),(4,7),(5,8),(6,8),(7,7),(7,6),(7,5),(7,4)]
    # A
    pixels_A = [(9,8),(10,6),(10,7),(10,8),(11,4),(11,5),(11,6),(11,7),(11,8),(12,6),(13,4),(13,5),(13,6),(13,7),(13,8)]
    # R
    pixels_R = [(15,4),(15,5),(15,6),(15,7),(15,8),(16,4),(16,6),(17,4),(17,5),(17,6),(17,7),(17,8)]
    # N
    pixels_N = [(19,4),(19,5),(19,6),(19,7),(19,8),(20,5),(21,6),(22,7),(23,4),(23,5),(23,6),(23,7),(23,8)]
    # I
    pixels_I = [(25,4),(25,5),(25,6),(25,7),(25,8)]
    # N
    pixels_N2 = [(27,4),(27,5),(27,6),(27,7),(27,8),(28,5),(29,6),(30,7),(31,4),(31,5),(31,6),(31,7),(31,8)]
    # G
    pixels_G = [(33,5),(33,6),(33,7),(34,4),(34,8),(35,4),(35,6),(35,7),(35,8),(36,4),(36,8),(37,5),(37,6),(37,7),(37,8)]

    all_pixels = pixels_W + pixels_A + pixels_R + pixels_N + pixels_I + pixels_N2 + pixels_G
    # 置中偏移
    offset_x = (W - 40) // 2
    offset_y = (H - 12) // 2

    for (x, y) in all_pixels:
        for scale in range(3):  # 3x 放大
            for dy in range(3):
                for dx in range(3):
                    px(img, offset_x + x*3 + dx, offset_y + y*3 + dy, (255, 255, 255, 255))

    return img


def main():
    print("=== 特效 v2 生成 ===\n")

    # 命中特效（3角色）
    for char_id in ["chiikawa", "hachiware", "usagi"]:
        img = gen_hit_effect(char_id)
        path = os.path.join(EFFECTS_DIR, f"hit_{char_id}.png")
        img.save(path)
        print(f"  ✅ hit_{char_id}.png: {img.size}")

    # 投射物（3角色）
    for char_id in ["chiikawa", "hachiware", "usagi"]:
        img = gen_projectile(char_id)
        path = os.path.join(EFFECTS_DIR, f"projectile_{char_id}.png")
        img.save(path)
        print(f"  ✅ projectile_{char_id}.png: {img.size}")

    # 死亡粒子
    img = gen_death_particles()
    path = os.path.join(EFFECTS_DIR, "death_particles.png")
    img.save(path)
    print(f"  ✅ death_particles.png: {img.size}")

    # WARNING 特效
    img = gen_warning_effect()
    path = os.path.join(EFFECTS_DIR, "warning.png")
    img.save(path)
    print(f"  ✅ warning.png: {img.size}")

    print(f"\n✅ 特效 v2 完成！")

if __name__ == "__main__":
    main()
