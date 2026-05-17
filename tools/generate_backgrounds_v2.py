# -*- coding: utf-8 -*-
"""
背景 v2 - 更豐富的像素藝術背景
1280x720，帶細節和層次感
"""
from PIL import Image, ImageDraw
import os
import math
import random

BG_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\backgrounds"
W, H = 1280, 720

def gen_sea_bg():
    """
    海底背景 v2
    - 深藍漸層（上深下淺）
    - 珊瑚礁（底部）
    - 氣泡（隨機分布）
    - 光線（從上方射入）
    - 海草（底部兩側）
    """
    img = Image.new("RGB", (W, H))
    draw = ImageDraw.Draw(img)
    rng = random.Random(42)

    # 海底漸層（上深下淺）
    for y in range(H):
        t = y / H
        r = int(5 + t * 20)
        g = int(30 + t * 60)
        b = int(80 + t * 60)
        draw.line([(0, y), (W, y)], fill=(r, g, b))

    # 光線（從上方射入，斜向）
    for i in range(8):
        x_start = rng.randint(100, W-100)
        width = rng.randint(20, 60)
        alpha_max = rng.randint(15, 35)
        for y in range(0, H//2):
            t = y / (H//2)
            x_offset = int(y * 0.3)
            alpha = int(alpha_max * (1 - t))
            for dx in range(-width//2, width//2):
                x = x_start + x_offset + dx
                if 0 <= x < W:
                    r, g, b = img.getpixel((x, y))
                    blend = alpha / 255
                    r = int(r + (180 - r) * blend * 0.3)
                    g = int(g + (220 - g) * blend * 0.3)
                    b = int(b + (255 - b) * blend * 0.3)
                    img.putpixel((x, y), (r, g, b))

    # 遠景珊瑚礁（底部，暗色）
    for i in range(12):
        x = rng.randint(50, W-50)
        h_coral = rng.randint(40, 120)
        w_coral = rng.randint(15, 40)
        color = (rng.randint(60,100), rng.randint(20,50), rng.randint(40,80))
        draw.ellipse([x-w_coral//2, H-h_coral, x+w_coral//2, H+10], fill=color)

    # 近景珊瑚礁（底部，亮色）
    coral_colors = [
        (220, 80, 60),   # 紅珊瑚
        (255, 140, 60),  # 橙珊瑚
        (180, 60, 160),  # 紫珊瑚
        (60, 180, 140),  # 綠珊瑚
    ]
    for i in range(18):
        x = rng.randint(0, W)
        h_coral = rng.randint(30, 90)
        w_coral = rng.randint(10, 30)
        color = coral_colors[rng.randint(0, len(coral_colors)-1)]
        # 珊瑚分支
        for branch in range(rng.randint(2, 5)):
            bx = x + rng.randint(-w_coral, w_coral)
            bh = rng.randint(h_coral//2, h_coral)
            bw = rng.randint(4, 12)
            draw.ellipse([bx-bw//2, H-bh, bx+bw//2, H+5], fill=color)
            # 頂部圓球
            draw.ellipse([bx-bw, H-bh-bw, bx+bw, H-bh+bw], fill=color)

    # 海草（底部兩側）
    for side in [0, 1]:
        x_base = rng.randint(0, 200) if side == 0 else rng.randint(W-200, W)
        for i in range(8):
            x = x_base + rng.randint(-80, 80)
            h_weed = rng.randint(60, 180)
            GREEN = (40, rng.randint(140, 200), 60)
            # 海草莖（波浪形）
            for y in range(H-h_weed, H):
                t = (y - (H-h_weed)) / h_weed
                wave = int(math.sin(t * math.pi * 3) * 8)
                for dx in range(-3, 4):
                    img.putpixel((min(W-1, max(0, x+wave+dx)), y), GREEN)

    # 氣泡
    for i in range(40):
        bx = rng.randint(0, W)
        by = rng.randint(0, H)
        br = rng.randint(2, 8)
        alpha = rng.randint(30, 80)
        for dy in range(-br, br+1):
            for dx in range(-br, br+1):
                if dx**2 + dy**2 <= br**2:
                    x, y = bx+dx, by+dy
                    if 0 <= x < W and 0 <= y < H:
                        r, g, b = img.getpixel((x, y))
                        blend = alpha / 255
                        r = int(r + (200-r)*blend)
                        g = int(g + (230-g)*blend)
                        b = int(b + (255-b)*blend)
                        img.putpixel((x, y), (r, g, b))
        # 氣泡高光
        hx, hy = bx-br//3, by-br//3
        if 0 <= hx < W and 0 <= hy < H:
            img.putpixel((hx, hy), (200, 230, 255))

    # 沙地（底部）
    for y in range(H-20, H):
        t = (y - (H-20)) / 20
        for x in range(W):
            r = int(180 + t*20 + rng.randint(-10, 10))
            g = int(160 + t*15 + rng.randint(-10, 10))
            b = int(100 + t*10 + rng.randint(-5, 5))
            img.putpixel((x, y), (min(255,r), min(255,g), min(255,b)))

    return img


def gen_boss_bg():
    """
    BOSS 背景 v2
    - 暗紅漸層
    - 閃電/裂縫效果
    - 暗黑光環
    - 警告紋路
    """
    img = Image.new("RGB", (W, H))
    draw = ImageDraw.Draw(img)
    rng = random.Random(123)

    # 暗紅漸層（中心亮，邊緣暗）
    for y in range(H):
        for x in range(W):
            cx = abs(x - W//2) / (W//2)
            cy = abs(y - H//2) / (H//2)
            dist = math.sqrt(cx**2 + cy**2)
            t = min(1.0, dist)
            r = int(80 - t*60)
            g = int(10 - t*8)
            b = int(10 - t*8)
            img.putpixel((x, y), (max(0,r), max(0,g), max(0,b)))

    # 警告條紋（斜向）
    for i in range(0, W+H, 80):
        for j in range(8):
            for y in range(H):
                x = i - y + j*2
                if 0 <= x < W:
                    r, g, b = img.getpixel((x, y))
                    r = min(255, r + 20)
                    img.putpixel((x, y), (r, g, b))

    # 裂縫效果
    for i in range(6):
        x = rng.randint(100, W-100)
        y = rng.randint(100, H-100)
        length = rng.randint(80, 200)
        angle = rng.uniform(0, math.pi*2)
        CRACK = (200, 50, 50)
        for j in range(length):
            angle += rng.uniform(-0.3, 0.3)
            cx = int(x + j * math.cos(angle))
            cy = int(y + j * math.sin(angle))
            if 0 <= cx < W and 0 <= cy < H:
                img.putpixel((cx, cy), CRACK)
                # 光暈
                for dx, dy in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = cx+dx, cy+dy
                    if 0 <= nx < W and 0 <= ny < H:
                        r, g, b = img.getpixel((nx, ny))
                        img.putpixel((nx, ny), (min(255,r+40), g, b))

    # 暗黑光環（中心）
    cx, cy = W//2, H//2
    for r in range(50, 300, 20):
        alpha = max(0, 60 - r//10)
        for angle in range(0, 360, 2):
            rad = math.radians(angle)
            x = int(cx + r * math.cos(rad))
            y = int(cy + r * math.sin(rad))
            if 0 <= x < W and 0 <= y < H:
                r_v, g_v, b_v = img.getpixel((x, y))
                blend = alpha / 255
                r_v = int(r_v + (180-r_v)*blend)
                img.putpixel((x, y), (r_v, g_v, b_v))

    # 底部地面（暗色石板）
    for y in range(H-60, H):
        t = (y - (H-60)) / 60
        for x in range(W):
            stone = rng.randint(-10, 10)
            r = int(40 + t*20 + stone)
            g = int(10 + t*5 + stone//2)
            b = int(10 + t*5 + stone//2)
            img.putpixel((x, y), (max(0,min(255,r)), max(0,min(255,g)), max(0,min(255,b))))

    # 石板縫隙
    for i in range(0, W, 120):
        for y in range(H-60, H):
            if 0 <= i < W:
                img.putpixel((i, y), (20, 5, 5))

    return img


def gen_bonus_bg():
    """
    Bonus 草地背景 v2
    - 明亮草地漸層
    - 草叢細節
    - 花朵
    - 天空（上方）
    - 雲朵
    """
    img = Image.new("RGB", (W, H))
    draw = ImageDraw.Draw(img)
    rng = random.Random(77)

    # 天空（上半部，藍色漸層）
    sky_h = H * 2 // 5
    for y in range(sky_h):
        t = y / sky_h
        r = int(135 + t*30)
        g = int(206 + t*20)
        b = int(235 - t*20)
        draw.line([(0, y), (W, y)], fill=(r, g, b))

    # 草地（下半部，綠色漸層）
    for y in range(sky_h, H):
        t = (y - sky_h) / (H - sky_h)
        r = int(60 + t*20)
        g = int(160 - t*30)
        b = int(40 - t*20)
        draw.line([(0, y), (W, y)], fill=(max(0,r), max(0,g), max(0,b)))

    # 雲朵
    for i in range(6):
        cx = rng.randint(100, W-100)
        cy = rng.randint(30, sky_h-60)
        for j in range(rng.randint(3, 6)):
            bx = cx + rng.randint(-60, 60)
            by = cy + rng.randint(-20, 20)
            br = rng.randint(25, 55)
            draw.ellipse([bx-br, by-br//2, bx+br, by+br//2], fill=(255, 255, 255))

    # 遠景樹木（天際線）
    for i in range(20):
        x = rng.randint(0, W)
        h_tree = rng.randint(40, 100)
        w_tree = rng.randint(20, 50)
        TREE = (30, rng.randint(100, 140), 30)
        draw.ellipse([x-w_tree//2, sky_h-h_tree, x+w_tree//2, sky_h+10], fill=TREE)

    # 草地細節（草叢）
    GRASS_COLORS = [
        (80, 180, 50),
        (60, 160, 40),
        (100, 200, 60),
        (70, 170, 45),
    ]
    for i in range(200):
        x = rng.randint(0, W)
        y = rng.randint(sky_h, H-20)
        h_g = rng.randint(8, 25)
        color = GRASS_COLORS[rng.randint(0, len(GRASS_COLORS)-1)]
        # 草葉
        for j in range(3):
            gx = x + rng.randint(-8, 8)
            draw.line([(gx, y), (gx + rng.randint(-5, 5), y-h_g)], fill=color, width=2)

    # 花朵
    FLOWER_COLORS = [
        (255, 100, 150),  # 粉紅
        (255, 220, 50),   # 黃色
        (200, 100, 255),  # 紫色
        (255, 150, 50),   # 橙色
        (255, 255, 255),  # 白色
    ]
    for i in range(60):
        x = rng.randint(0, W)
        y = rng.randint(sky_h+20, H-10)
        color = FLOWER_COLORS[rng.randint(0, len(FLOWER_COLORS)-1)]
        r_f = rng.randint(4, 10)
        # 花瓣
        for angle in range(0, 360, 60):
            rad = math.radians(angle)
            px_f = int(x + r_f * math.cos(rad))
            py_f = int(y + r_f * math.sin(rad))
            draw.ellipse([px_f-3, py_f-3, px_f+3, py_f+3], fill=color)
        # 花心
        draw.ellipse([x-3, y-3, x+3, y+3], fill=(255, 230, 50))

    # 地面（底部深色）
    for y in range(H-15, H):
        t = (y - (H-15)) / 15
        for x in range(W):
            r = int(50 + t*10)
            g = int(100 - t*20)
            b = int(30)
            img.putpixel((x, y), (r, g, b))

    return img


def main():
    os.makedirs(BG_DIR, exist_ok=True)
    print("=== 背景 v2 生成 ===\n")

    bgs = [
        ("sea_bg.png",   gen_sea_bg,   "海底背景"),
        ("boss_bg.png",  gen_boss_bg,  "BOSS 背景"),
        ("bonus_bg.png", gen_bonus_bg, "Bonus 草地背景"),
    ]

    for filename, fn, desc in bgs:
        print(f"[{desc}]")
        img = fn()
        path = os.path.join(BG_DIR, filename)
        img.save(path)
        # 統計顏色多樣性
        pixels = list(img.getdata())
        unique = len(set(pixels))
        print(f"  ✅ {filename}: {img.size}, {unique:,} 種顏色")
        print()

    print("✅ 背景 v2 生成完成！")

if __name__ == "__main__":
    main()
