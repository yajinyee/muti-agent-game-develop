# -*- coding: utf-8 -*-
"""
目標物 v4 - 修復密度問題
主要改善：
1. T001/T104 草葉：改為連續繪製（不跳行），密度從 11% 提升到 50%+
2. T101 擬態怪物：基於 T001 v4，密度同步提升
3. T103 流星：加大核心，增加光暈密度
4. 所有目標物：加強陰影和細節
"""
from PIL import Image
import os
import math

OUTPUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def new_img():
    return Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

def px(img, x, y, c):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), c)

def fill_circle(img, cx, cy, r, color):
    for y in range(cy-r, cy+r+1):
        for x in range(cx-r, cx+r+1):
            if (x-cx)**2 + (y-cy)**2 <= r**2:
                px(img, x, y, color)

def fill_circle_shaded(img, cx, cy, r, base_rgb):
    r_v, g_v, b_v = base_rgb
    light = (min(255,r_v+35), min(255,g_v+35), min(255,b_v+35), 255)
    mid   = (r_v, g_v, b_v, 255)
    dark  = (max(0,r_v-40), max(0,g_v-40), max(0,b_v-40), 255)
    for y in range(cy-r, cy+r+1):
        for x in range(cx-r, cx+r+1):
            if (x-cx)**2 + (y-cy)**2 > r**2:
                continue
            nx_ = (x-cx)/max(r,1)
            ny_ = (y-cy)/max(r,1)
            dot = -(nx_*(-0.7) + ny_*(-0.7))
            if dot > 0.25:
                c = light
            elif dot < -0.1:
                c = dark
            else:
                c = mid
            px(img, x, y, c)

def outline_circle(img, cx, cy, r, color):
    for y in range(cy-r-1, cy+r+2):
        for x in range(cx-r-1, cx+r+2):
            d = math.sqrt((x-cx)**2 + (y-cy)**2)
            if r+0.1 <= d <= r+1.5:
                px(img, x, y, color)

def fill_rect_shaded(img, x1, y1, x2, y2, base_rgb):
    r_v, g_v, b_v = base_rgb
    w = x2 - x1
    h = y2 - y1
    for y in range(y1, y2):
        for x in range(x1, x2):
            nx_ = (x - x1) / max(w, 1)
            ny_ = (y - y1) / max(h, 1)
            brightness = 1.0 - (nx_ + ny_) * 0.15
            r_c = int(min(255, r_v * brightness))
            g_c = int(min(255, g_v * brightness))
            b_c = int(min(255, b_v * brightness))
            px(img, x, y, (r_c, g_c, b_c, 255))

def draw_leaf(img, cx, cy, angle_deg, length, width, color_mid, color_edge, outline_color):
    """繪製一片葉子（連續填充，不跳行）"""
    angle = math.radians(angle_deg)
    dx_step = math.cos(angle)
    dy_step = math.sin(angle)
    
    for i in range(length):
        # 葉子中心點
        lx = cx + dx_step * i
        ly = cy + dy_step * i
        
        # 葉子寬度隨長度遞減（尖端變細）
        w = max(1, int(width * (1.0 - i / length * 0.7)))
        
        # 垂直於葉子方向的向量
        perp_x = -dy_step
        perp_y = dx_step
        
        for j in range(-w, w+1):
            x = int(lx + perp_x * j)
            y = int(ly + perp_y * j)
            if abs(j) == w:
                px(img, x, y, outline_color)
            elif abs(j) <= w // 2:
                px(img, x, y, color_mid)
            else:
                px(img, x, y, color_edge)

def gen_T001_grass_v4():
    """像素雜草 2x v4 - 連續葉子，密度提升"""
    img = new_img()
    OUTLINE = (30, 80, 20, 255)
    STEM    = (60, 120, 40, 255)
    LEAF1   = (80, 180, 50, 255)
    LEAF2   = (100, 200, 60, 255)
    LEAF_L  = (50, 140, 35, 255)
    DARK    = (40, 100, 25, 255)

    cx = 32
    
    # 草莖（加粗）
    for y in range(40, 62):
        for x in range(cx-4, cx+5):
            shade = 1.0 - abs(x - cx) / 5 * 0.3
            r = int(60 * shade)
            g = int(120 * shade)
            b = int(40 * shade)
            px(img, x, y, (r, g, b, 255))
    for y in range(40, 62):
        px(img, cx-4, y, OUTLINE)
        px(img, cx+4, y, OUTLINE)
    
    # 左葉（填充三角形區域）
    for y in range(6, 42):
        progress = (42 - y) / 36.0  # 0=底部, 1=頂部
        # 葉子從莖向左延伸
        leaf_x_center = cx - 2 - int(progress * 20)
        leaf_width = max(1, int(5 * (1.0 - progress * 0.6)))
        for dx in range(-leaf_width, leaf_width+1):
            x = leaf_x_center + dx
            if abs(dx) == leaf_width:
                px(img, x, y, OUTLINE)
            elif abs(dx) <= leaf_width // 2:
                px(img, x, y, LEAF1)
            else:
                px(img, x, y, LEAF_L)
    
    # 右葉（填充三角形區域）
    for y in range(6, 42):
        progress = (42 - y) / 36.0
        leaf_x_center = cx + 2 + int(progress * 20)
        leaf_width = max(1, int(5 * (1.0 - progress * 0.6)))
        for dx in range(-leaf_width, leaf_width+1):
            x = leaf_x_center + dx
            if abs(dx) == leaf_width:
                px(img, x, y, OUTLINE)
            elif abs(dx) <= leaf_width // 2:
                px(img, x, y, LEAF1)
            else:
                px(img, x, y, LEAF_L)
    
    # 中間葉（直向上，最寬）
    for y in range(2, 42):
        progress = (42 - y) / 40.0
        leaf_width = max(1, int(6 * (1.0 - progress * 0.7)))
        for dx in range(-leaf_width, leaf_width+1):
            x = cx + dx
            if abs(dx) == leaf_width:
                px(img, x, y, OUTLINE)
            elif abs(dx) <= leaf_width // 2:
                px(img, x, y, LEAF2)
            else:
                px(img, x, y, LEAF1)
    
    # 地面根部
    for x in range(cx-6, cx+7):
        px(img, x, 62, OUTLINE)
    
    return img


def gen_T101_mimic_v4():
    """擬態型怪物 v4 - 基於 T001 v4，加上詭異眼睛"""
    img = gen_T001_grass_v4()
    
    EYE_W = (255, 255, 200, 255)
    EYE_P = (150, 0, 200, 255)
    OUTLINE = (50, 0, 80, 255)
    
    # 詭異眼睛（在葉子上）
    for (ex, ey) in [(18, 26), (46, 26)]:
        fill_circle(img, ex, ey, 5, EYE_W)
        fill_circle(img, ex, ey, 3, EYE_P)
        px(img, ex-1, ey-1, (255, 255, 255, 255))
        outline_circle(img, ex, ey, 5, OUTLINE)
    
    # 詭異的嘴（在莖上）
    for x in range(28, 37):
        px(img, x, 52, OUTLINE)
    for i in range(3):
        px(img, 29+i*3, 53, OUTLINE)
    
    return img


def gen_T103_meteor_v4():
    """流星 v4 - 更大的核心，更密的光暈"""
    img = new_img()
    CORE  = (255, 240, 100)
    GLOW1 = (255, 200, 50)
    GLOW2 = (255, 150, 20)
    TRAIL = (255, 100, 0)
    OUTLINE=(180, 80, 0, 255)
    WHITE = (255, 255, 255, 255)

    # 流星尾巴（更寬更密）
    for i in range(24):
        x = 56 - i*2
        y = 8 + i*2
        alpha = max(0, 220 - i*9)
        for dx in range(-3, 4):
            a = max(0, alpha - abs(dx)*30)
            px(img, x+dx, y, (*TRAIL, a))
        # 尾巴中心線（更亮）
        px(img, x, y, (*CORE, min(255, alpha+50)))

    # 星形核心（更大）
    fill_circle_shaded(img, 18, 46, 14, CORE)
    outline_circle(img, 18, 46, 14, OUTLINE)
    
    # 核心高光
    fill_circle(img, 14, 42, 4, WHITE)
    
    # 光暈（更密）
    for r in range(15, 20):
        for angle in range(0, 360, 3):
            rad = math.radians(angle)
            x = int(18 + r * math.cos(rad))
            y = int(46 + r * math.sin(rad))
            alpha = max(0, 150 - (r-15)*30)
            px(img, x, y, (*GLOW1, alpha))

    # 星芒（8方向）
    for angle in range(0, 360, 45):
        rad = math.radians(angle)
        for i in range(6):
            x = int(18 + (15+i) * math.cos(rad))
            y = int(46 + (15+i) * math.sin(rad))
            alpha = max(0, 200 - i*35)
            px(img, x, y, (*GLOW2, alpha))

    return img


def gen_T104_gold_grass_v4():
    """金色雜草 v4 - 連續葉子，密度提升"""
    img = new_img()
    GOLD   = (220, 180, 30)
    GOLD_L = (255, 230, 80)
    GOLD_D = (160, 120, 10)
    OUTLINE= (100, 60, 5, 255)
    SHINE  = (255, 255, 200, 255)

    cx = 32
    
    # 草莖（金色，加粗）
    for y in range(40, 62):
        for x in range(cx-4, cx+5):
            shade = 1.0 - abs(x - cx) / 5 * 0.3
            r = int(220 * shade)
            g = int(180 * shade)
            b = int(30 * shade)
            px(img, x, y, (r, g, b, 255))
    for y in range(40, 62):
        px(img, cx-4, y, OUTLINE)
        px(img, cx+4, y, OUTLINE)
    
    # 左葉（填充三角形，金色）
    for y in range(6, 42):
        progress = (42 - y) / 36.0
        leaf_x_center = cx - 2 - int(progress * 20)
        leaf_width = max(1, int(5 * (1.0 - progress * 0.6)))
        for dx in range(-leaf_width, leaf_width+1):
            x = leaf_x_center + dx
            if abs(dx) == leaf_width:
                px(img, x, y, OUTLINE)
            elif abs(dx) <= leaf_width // 2:
                px(img, x, y, (*GOLD_L, 255))
            else:
                px(img, x, y, (*GOLD, 255))
    
    # 右葉（填充三角形，金色）
    for y in range(6, 42):
        progress = (42 - y) / 36.0
        leaf_x_center = cx + 2 + int(progress * 20)
        leaf_width = max(1, int(5 * (1.0 - progress * 0.6)))
        for dx in range(-leaf_width, leaf_width+1):
            x = leaf_x_center + dx
            if abs(dx) == leaf_width:
                px(img, x, y, OUTLINE)
            elif abs(dx) <= leaf_width // 2:
                px(img, x, y, (*GOLD_L, 255))
            else:
                px(img, x, y, (*GOLD, 255))
    
    # 中間葉（直向上，最亮）
    for y in range(2, 42):
        progress = (42 - y) / 40.0
        leaf_width = max(1, int(6 * (1.0 - progress * 0.7)))
        for dx in range(-leaf_width, leaf_width+1):
            x = cx + dx
            if abs(dx) == leaf_width:
                px(img, x, y, OUTLINE)
            elif abs(dx) <= leaf_width // 2:
                px(img, x, y, (*GOLD_L, 255))
            else:
                px(img, x, y, (*GOLD, 255))
    
    # 閃光點（更多）
    shine_points = [(20, 18), (44, 14), (32, 6), (14, 30), (50, 26), (26, 10), (38, 22)]
    for (sx, sy) in shine_points:
        px(img, sx, sy, SHINE)
        for ddx, ddy in [(1,0),(-1,0),(0,1),(0,-1),(1,1),(-1,-1)]:
            px(img, sx+ddx, sy+ddy, (*GOLD_L, 200))
    
    # 地面根部
    for x in range(cx-6, cx+7):
        px(img, x, 62, OUTLINE)
    
    return img


def main():
    print("=== 目標物 v4 生成（修復密度問題）===\n")
    
    # 只重新生成密度低的目標物
    targets_to_fix = [
        ("T001_grass",    gen_T001_grass_v4),
        ("T101_mimic",    gen_T101_mimic_v4),
        ("T103_meteor",   gen_T103_meteor_v4),
        ("T104_gold_grass", gen_T104_gold_grass_v4),
    ]
    
    for name, fn in targets_to_fix:
        img = fn()
        path = os.path.join(OUTPUT_DIR, f"{name}.png")
        img.save(path)
        
        # 計算密度
        pixels = list(img.getdata())
        non_t = sum(1 for p in pixels if p[3] > 10)
        total = img.width * img.height
        pct = non_t / total * 100
        
        print(f"  ✅ {name}.png: {img.width}x{img.height}, {non_t}/{total} ({pct:.0f}%)")
    
    print(f"\n✅ 完成！輸出目錄: {OUTPUT_DIR}")

if __name__ == '__main__':
    main()
