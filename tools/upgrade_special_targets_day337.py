#!/usr/bin/env python3
"""
DAY-337: T101-T105 特殊目標物視覺升級
T101 擬態怪物 - 更有個性的偽裝外觀
T102 寶箱怪   - 更明顯的寶箱特徵
T103 流星     - 更有速度感的流星
T104 金草     - 更閃亮的金色草
T105 金幣魚   - 更有魚的特徵
"""
from PIL import Image, ImageDraw
import os
import math
import random

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def save_png(img, name):
    path = os.path.join(OUT_DIR, name)
    img.save(path)
    print(f"  ✅ 儲存 {name} ({img.size[0]}x{img.size[1]})")

def px(img, x, y, color):
    if 0 <= x < img.width and 0 <= y < img.height:
        img.putpixel((x, y), color)

def fill_circle(img, cx, cy, r, color):
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                px(img, cx+dx, cy+dy, color)

def fill_circle_shaded(img, cx, cy, r, light, mid, dark, outline=(0,0,0,255)):
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            dist = math.sqrt(dx*dx + dy*dy)
            if dist <= r:
                if dist >= r - 1:
                    c = outline
                elif dx < -r*0.2 and dy < -r*0.2:
                    c = light
                elif dx > r*0.2 and dy > r*0.2:
                    c = dark
                else:
                    c = mid
                px(img, cx+dx, cy+dy, c)

def fill_rect(img, x1, y1, x2, y2, color):
    for y in range(y1, y2+1):
        for x in range(x1, x2+1):
            px(img, x, y, color)

def draw_outline(img, x1, y1, x2, y2, color=(0,0,0,255)):
    for x in range(x1, x2+1):
        px(img, x, y1, color)
        px(img, x, y2, color)
    for y in range(y1, y2+1):
        px(img, x1, y, color)
        px(img, x2, y, color)

# ── T101 擬態怪物（Mimic）────────────────────────────────────
def gen_t101_mimic():
    """擬態怪物：看起來像普通草，但有詭異的眼睛"""
    img = Image.new("RGBA", (SIZE, SIZE), (0,0,0,0))
    
    # 草的主體（偽裝成 T001）
    GRASS_G = (34, 139, 34, 255)
    GRASS_L = (50, 180, 50, 255)
    GRASS_D = (20, 100, 20, 255)
    OUTLINE = (0, 0, 0, 255)
    
    # 畫幾根草莖（偽裝）
    stems = [(20, 55, 18, 30), (28, 55, 26, 22), (36, 55, 34, 28), (44, 55, 42, 25)]
    for x1, y1, x2, y2 in stems:
        for t in range(20):
            tt = t / 19.0
            x = int(x1 + (x2 - x1) * tt)
            y = int(y1 + (y2 - y1) * tt)
            c = GRASS_L if t < 5 else GRASS_G
            px(img, x, y, c)
            px(img, x+1, y, c)
    
    # 草葉（三角形）
    for i, (sx, sy, ex, ey) in enumerate(stems):
        # 左葉
        for t in range(12):
            tt = t / 11.0
            lx = int(sx - 8 * tt)
            ly = int(sy + (ey - sy) * 0.3 - 6 * tt)
            px(img, lx, ly, GRASS_G)
            px(img, lx, ly+1, GRASS_D)
        # 右葉
        for t in range(12):
            tt = t / 11.0
            rx = int(sx + 8 * tt)
            ry = int(sy + (ey - sy) * 0.3 - 6 * tt)
            px(img, rx, ry, GRASS_G)
            px(img, rx, ry+1, GRASS_D)
    
    # 地面
    fill_rect(img, 12, 54, 52, 57, (101, 67, 33, 255))
    fill_rect(img, 10, 55, 54, 56, (120, 80, 40, 255))
    draw_outline(img, 10, 54, 54, 57, OUTLINE)
    
    # 詭異的眼睛（隱藏在草叢中）
    EYE_W = (255, 255, 255, 255)
    EYE_R = (255, 0, 0, 255)
    EYE_B = (0, 0, 0, 255)
    
    # 左眼（在草叢中若隱若現）
    fill_circle(img, 22, 38, 4, EYE_W)
    fill_circle(img, 23, 38, 2, EYE_R)
    px(img, 23, 38, EYE_B)
    px(img, 22, 37, (255, 255, 255, 180))  # 高光
    
    # 右眼
    fill_circle(img, 40, 35, 4, EYE_W)
    fill_circle(img, 41, 35, 2, EYE_R)
    px(img, 41, 35, EYE_B)
    px(img, 40, 34, (255, 255, 255, 180))  # 高光
    
    # 詭異的嘴巴（鋸齒形）
    MOUTH = (200, 0, 0, 255)
    for i in range(6):
        x = 26 + i * 3
        y = 44 if i % 2 == 0 else 46
        px(img, x, y, MOUTH)
        px(img, x+1, y, MOUTH)
    
    # 金色輪廓（特殊目標標記）
    GOLD = (255, 215, 0, 255)
    for x in range(8, 56):
        for y in range(18, 60):
            if img.getpixel((x, y))[3] > 0:
                # 檢查是否是邊緣
                is_edge = False
                for dx, dy in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x+dx, y+dy
                    if 0 <= nx < SIZE and 0 <= ny < SIZE:
                        if img.getpixel((nx, ny))[3] == 0:
                            is_edge = True
                            break
                if is_edge:
                    px(img, x, y, GOLD)
    
    return img

# ── T102 寶箱怪（Chest）─────────────────────────────────────
def gen_t102_chest():
    """寶箱怪：明顯的寶箱外觀，有眼睛和牙齒"""
    img = Image.new("RGBA", (SIZE, SIZE), (0,0,0,0))
    
    CHEST_BROWN = (139, 90, 43, 255)
    CHEST_LIGHT = (180, 120, 60, 255)
    CHEST_DARK = (100, 60, 20, 255)
    GOLD = (255, 215, 0, 255)
    GOLD_D = (200, 160, 0, 255)
    OUTLINE = (0, 0, 0, 255)
    
    # 寶箱主體（下半部）
    fill_rect(img, 8, 36, 56, 56, CHEST_BROWN)
    fill_rect(img, 10, 38, 54, 54, CHEST_LIGHT)
    # 木紋
    for y in range(40, 54, 4):
        fill_rect(img, 10, y, 54, y, CHEST_DARK)
    
    # 寶箱蓋（上半部，微微打開）
    fill_rect(img, 8, 18, 56, 36, CHEST_BROWN)
    fill_rect(img, 10, 20, 54, 34, CHEST_LIGHT)
    # 蓋子木紋
    for y in range(22, 34, 4):
        fill_rect(img, 10, y, 54, y, CHEST_DARK)
    
    # 金色鎖扣
    fill_circle(img, 32, 36, 5, GOLD)
    fill_circle(img, 32, 36, 3, GOLD_D)
    px(img, 32, 35, (255, 255, 200, 255))  # 高光
    
    # 金色邊框
    for x in range(8, 57):
        px(img, x, 18, GOLD)
        px(img, x, 56, GOLD)
        px(img, x, 36, GOLD)
    for y in range(18, 57):
        px(img, 8, y, GOLD)
        px(img, 56, y, GOLD)
    
    # 眼睛（在蓋子上）
    EYE_W = (255, 255, 255, 255)
    EYE_Y = (255, 220, 0, 255)
    EYE_B = (0, 0, 0, 255)
    
    fill_circle(img, 22, 26, 5, EYE_W)
    fill_circle(img, 22, 26, 3, EYE_Y)
    px(img, 22, 26, EYE_B)
    px(img, 21, 25, (255, 255, 255, 200))
    
    fill_circle(img, 42, 26, 5, EYE_W)
    fill_circle(img, 42, 26, 3, EYE_Y)
    px(img, 42, 26, EYE_B)
    px(img, 41, 25, (255, 255, 255, 200))
    
    # 牙齒（在蓋子縫隙）
    TOOTH = (255, 255, 255, 255)
    for i in range(5):
        x = 16 + i * 8
        px(img, x, 35, TOOTH)
        px(img, x, 36, TOOTH)
        px(img, x+1, 35, TOOTH)
        px(img, x+1, 36, TOOTH)
    
    # 輪廓
    for x in range(8, 57):
        for y in range(18, 57):
            if img.getpixel((x, y))[3] > 0:
                for dx, dy in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x+dx, y+dy
                    if 0 <= nx < SIZE and 0 <= ny < SIZE:
                        if img.getpixel((nx, ny))[3] == 0:
                            px(img, x, y, OUTLINE)
                            break
    
    return img

# ── T103 流星（Meteor）──────────────────────────────────────
def gen_t103_meteor():
    """流星：有速度感的流星，帶火焰尾跡"""
    img = Image.new("RGBA", (SIZE, SIZE), (0,0,0,0))
    
    ROCK_G = (100, 100, 110, 255)
    ROCK_L = (140, 140, 150, 255)
    ROCK_D = (60, 60, 70, 255)
    FIRE_R = (255, 80, 0, 255)
    FIRE_O = (255, 160, 0, 255)
    FIRE_Y = (255, 220, 0, 255)
    GLOW = (255, 100, 0, 100)
    OUTLINE = (0, 0, 0, 255)
    
    # 火焰尾跡（從右上到左下）
    for i in range(30):
        t = i / 29.0
        x = int(52 - i * 1.2)
        y = int(12 + i * 0.8)
        r = max(1, int(8 - i * 0.2))
        alpha = int(255 * (1 - t * 0.8))
        if i < 10:
            c = (FIRE_Y[0], FIRE_Y[1], FIRE_Y[2], alpha)
        elif i < 20:
            c = (FIRE_O[0], FIRE_O[1], FIRE_O[2], alpha)
        else:
            c = (FIRE_R[0], FIRE_R[1], FIRE_R[2], alpha)
        fill_circle(img, x, y, r, c)
    
    # 流星主體（橢圓形岩石）
    cx, cy = 20, 44
    for dy in range(-12, 13):
        for dx in range(-16, 17):
            if (dx/16)**2 + (dy/12)**2 <= 1:
                dist_from_center = math.sqrt((dx/16)**2 + (dy/12)**2)
                if dist_from_center >= 0.85:
                    c = OUTLINE
                elif dx < -4 and dy < -4:
                    c = ROCK_L
                elif dx > 4 and dy > 4:
                    c = ROCK_D
                else:
                    c = ROCK_G
                px(img, cx+dx, cy+dy, c)
    
    # 岩石裂縫
    for i in range(5):
        x = cx - 8 + i * 3
        y = cy - 4 + (i % 3) * 2
        px(img, x, y, ROCK_D)
    
    # 發光效果（橘色光暈）
    for dy in range(-14, 15):
        for dx in range(-18, 19):
            if (dx/18)**2 + (dy/14)**2 <= 1 and (dx/16)**2 + (dy/12)**2 > 1:
                alpha = int(80 * (1 - ((dx/18)**2 + (dy/14)**2)))
                if alpha > 0:
                    existing = img.getpixel((cx+dx, cy+dy))
                    if existing[3] == 0:
                        px(img, cx+dx, cy+dy, (255, 100, 0, alpha))
    
    # 速度線
    for i in range(4):
        y = cy - 8 + i * 5
        for x in range(cx + 16, cx + 28):
            alpha = int(200 * (1 - (x - cx - 16) / 12))
            if alpha > 0:
                px(img, x, y, (255, 200, 100, alpha))
    
    return img

# ── T104 金草（Gold Grass）──────────────────────────────────
def gen_t104_gold_grass():
    """金草：閃亮的金色草，有光澤效果"""
    img = Image.new("RGBA", (SIZE, SIZE), (0,0,0,0))
    
    GOLD = (255, 215, 0, 255)
    GOLD_L = (255, 240, 100, 255)
    GOLD_D = (200, 160, 0, 255)
    GOLD_DARK = (160, 120, 0, 255)
    OUTLINE = (0, 0, 0, 255)
    SHINE = (255, 255, 200, 255)
    
    # 金草莖（5根，不同高度）
    stems = [
        (16, 55, 14, 20),
        (24, 55, 22, 14),
        (32, 55, 30, 10),
        (40, 55, 38, 16),
        (48, 55, 46, 22),
    ]
    
    for x1, y1, x2, y2 in stems:
        for t in range(25):
            tt = t / 24.0
            x = int(x1 + (x2 - x1) * tt)
            y = int(y1 + (y2 - y1) * tt)
            # 莖的顏色漸變
            if t < 5:
                c = GOLD_DARK
            elif t < 15:
                c = GOLD
            else:
                c = GOLD_L
            px(img, x, y, c)
            px(img, x+1, y, c)
            # 輪廓
            px(img, x-1, y, OUTLINE)
            px(img, x+2, y, OUTLINE)
    
    # 金草葉（三角形，帶光澤）
    for i, (sx, sy, ex, ey) in enumerate(stems):
        mid_x = (sx + ex) // 2
        mid_y = (sy + ey) // 2
        
        # 左葉
        for t in range(14):
            tt = t / 13.0
            lx = int(sx - 10 * tt)
            ly = int(mid_y - 8 * tt)
            c = GOLD_L if t < 4 else GOLD
            px(img, lx, ly, c)
            px(img, lx, ly+1, GOLD_D)
            # 光澤點
            if t == 2:
                px(img, lx, ly, SHINE)
        
        # 右葉
        for t in range(14):
            tt = t / 13.0
            rx = int(sx + 10 * tt)
            ry = int(mid_y - 8 * tt)
            c = GOLD_L if t < 4 else GOLD
            px(img, rx, ry, c)
            px(img, rx, ry+1, GOLD_D)
            if t == 2:
                px(img, rx, ry, SHINE)
    
    # 地面（金色土壤）
    fill_rect(img, 10, 54, 54, 57, GOLD_DARK)
    fill_rect(img, 12, 55, 52, 56, GOLD_D)
    draw_outline(img, 10, 54, 54, 57, OUTLINE)
    
    # 閃光效果（星形光點）
    for star_x, star_y in [(18, 18), (34, 12), (50, 20), (26, 30)]:
        px(img, star_x, star_y, SHINE)
        px(img, star_x-1, star_y, GOLD_L)
        px(img, star_x+1, star_y, GOLD_L)
        px(img, star_x, star_y-1, GOLD_L)
        px(img, star_x, star_y+1, GOLD_L)
    
    return img

# ── T105 金幣魚（Coin Fish）─────────────────────────────────
def gen_t105_coin_fish():
    """金幣魚：有明顯魚形特徵，金色閃亮"""
    img = Image.new("RGBA", (SIZE, SIZE), (0,0,0,0))
    
    GOLD = (255, 215, 0, 255)
    GOLD_L = (255, 240, 120, 255)
    GOLD_D = (200, 160, 0, 255)
    GOLD_DARK = (160, 120, 0, 255)
    OUTLINE = (0, 0, 0, 255)
    EYE_W = (255, 255, 255, 255)
    EYE_B = (0, 0, 0, 255)
    COIN_Y = (255, 200, 0, 255)
    SHINE = (255, 255, 200, 255)
    
    # 魚身（橢圓形，橫向）
    cx, cy = 30, 32
    for dy in range(-14, 15):
        for dx in range(-22, 23):
            if (dx/22)**2 + (dy/14)**2 <= 1:
                dist = math.sqrt((dx/22)**2 + (dy/14)**2)
                if dist >= 0.88:
                    c = OUTLINE
                elif dx < -8 and dy < -4:
                    c = GOLD_L
                elif dx > 8 and dy > 4:
                    c = GOLD_D
                else:
                    c = GOLD
                px(img, cx+dx, cy+dy, c)
    
    # 魚鱗（弧形紋路）
    for i in range(4):
        for j in range(3):
            sx = cx - 15 + i * 8
            sy = cy - 6 + j * 6
            for t in range(8):
                angle = math.pi * t / 7
                x = int(sx + 4 * math.cos(angle))
                y = int(sy + 3 * math.sin(angle))
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    if img.getpixel((x, y))[3] > 0:
                        px(img, x, y, GOLD_D)
    
    # 魚尾（三角形）
    for dy in range(-10, 11):
        for dx in range(0, 16):
            if abs(dy) <= dx * 0.7:
                dist = math.sqrt(dx*dx + dy*dy)
                if dist >= 14:
                    c = OUTLINE
                elif dx < 8:
                    c = GOLD
                else:
                    c = GOLD_D
                px(img, cx + 20 + dx, cy + dy, c)
    
    # 魚鰭（上方）
    for i in range(10):
        x = cx - 10 + i * 2
        y = cy - 14 - i // 3
        px(img, x, y, GOLD_D)
        px(img, x, y-1, GOLD)
    
    # 眼睛
    fill_circle(img, cx - 14, cy - 2, 5, EYE_W)
    fill_circle(img, cx - 14, cy - 2, 3, COIN_Y)
    px(img, cx - 14, cy - 2, EYE_B)
    px(img, cx - 15, cy - 3, SHINE)
    
    # 嘴巴
    for i in range(4):
        px(img, cx - 22 + i, cy + 2, OUTLINE)
    
    # 金幣符號（在魚身上）
    COIN_C = (255, 180, 0, 255)
    fill_circle(img, cx + 2, cy, 7, COIN_C)
    fill_circle(img, cx + 2, cy, 5, COIN_Y)
    # ¥ 符號
    px(img, cx + 2, cy - 3, OUTLINE)
    px(img, cx + 2, cy - 2, OUTLINE)
    px(img, cx + 2, cy - 1, OUTLINE)
    px(img, cx + 2, cy, OUTLINE)
    px(img, cx + 2, cy + 1, OUTLINE)
    px(img, cx - 1, cy - 3, OUTLINE)
    px(img, cx + 5, cy - 3, OUTLINE)
    px(img, cx, cy - 2, OUTLINE)
    px(img, cx + 4, cy - 2, OUTLINE)
    px(img, cx, cy - 1, OUTLINE)
    px(img, cx + 4, cy - 1, OUTLINE)
    
    # 光澤
    px(img, cx - 5, cy - 8, SHINE)
    px(img, cx - 4, cy - 8, GOLD_L)
    px(img, cx - 5, cy - 7, GOLD_L)
    
    return img

# ── 主程式 ───────────────────────────────────────────────────
print("DAY-337: T101-T105 特殊目標物視覺升級")
print("=" * 50)

generators = [
    ("T101_mimic.png", gen_t101_mimic),
    ("T102_chest.png", gen_t102_chest),
    ("T103_meteor.png", gen_t103_meteor),
    ("T104_gold_grass.png", gen_t104_gold_grass),
    ("T105_coin_fish.png", gen_t105_coin_fish),
]

for filename, gen_func in generators:
    try:
        img = gen_func()
        save_png(img, filename)
        
        # 計算密度
        pixels = img.load()
        non_transparent = sum(1 for y in range(SIZE) for x in range(SIZE) if pixels[x,y][3] > 10)
        density = non_transparent / (SIZE * SIZE) * 100
        print(f"    密度: {density:.1f}%")
    except Exception as e:
        print(f"  ❌ {filename}: {e}")
        import traceback
        traceback.print_exc()

print("\n✅ T101-T105 視覺升級完成")
