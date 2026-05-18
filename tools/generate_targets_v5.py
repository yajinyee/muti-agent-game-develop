# -*- coding: utf-8 -*-
"""
目標物 v5 - 根本性密度提升
策略：
1. 草類目標物加入「土堆底座」，讓整體更飽滿
2. 蟲類目標物加入「陰影投影」，增加存在感
3. 所有目標物加入更豐富的細節（高光、紋理、輪廓）
4. 目標密度目標：>= 40%
"""
from PIL import Image, ImageDraw
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
    for y in range(max(0,cy-r), min(SIZE,cy+r+1)):
        for x in range(max(0,cx-r), min(SIZE,cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 <= r**2:
                px(img, x, y, color)

def fill_circle_shaded(img, cx, cy, r, base_rgb):
    r_v, g_v, b_v = base_rgb
    for y in range(max(0,cy-r), min(SIZE,cy+r+1)):
        for x in range(max(0,cx-r), min(SIZE,cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 > r**2:
                continue
            nx_ = (x-cx)/max(r,1)
            ny_ = (y-cy)/max(r,1)
            dot = -(nx_*(-0.7) + ny_*(-0.7))
            if dot > 0.25:
                c = (min(255,r_v+40), min(255,g_v+40), min(255,b_v+40), 255)
            elif dot < -0.1:
                c = (max(0,r_v-45), max(0,g_v-45), max(0,b_v-45), 255)
            else:
                c = (r_v, g_v, b_v, 255)
            px(img, x, y, c)

def outline_circle(img, cx, cy, r, color):
    for y in range(max(0,cy-r-2), min(SIZE,cy+r+3)):
        for x in range(max(0,cx-r-2), min(SIZE,cx+r+3)):
            d = math.sqrt((x-cx)**2 + (y-cy)**2)
            if r+0.2 <= d <= r+1.8:
                px(img, x, y, color)

def fill_rect(img, x1, y1, x2, y2, color):
    for y in range(max(0,y1), min(SIZE,y2)):
        for x in range(max(0,x1), min(SIZE,x2)):
            px(img, x, y, color)

def fill_rect_shaded(img, x1, y1, x2, y2, base_rgb):
    r_v, g_v, b_v = base_rgb
    w = max(x2-x1, 1); h = max(y2-y1, 1)
    for y in range(max(0,y1), min(SIZE,y2)):
        for x in range(max(0,x1), min(SIZE,x2)):
            nx_ = (x-x1)/w; ny_ = (y-y1)/h
            b = 1.0 - (nx_+ny_)*0.12
            px(img, x, y, (int(min(255,r_v*b)), int(min(255,g_v*b)), int(min(255,b_v*b)), 255))

def draw_soil_mound(img, cx, base_y, w, h, color=(120, 80, 40, 255)):
    """畫土堆底座（橢圓形），讓草類目標物更飽滿"""
    dark = (max(0,color[0]-30), max(0,color[1]-20), max(0,color[2]-10), 255)
    light = (min(255,color[0]+20), min(255,color[1]+15), min(255,color[2]+5), 255)
    for y in range(base_y - h, base_y + 1):
        t = (y - (base_y - h)) / max(h, 1)
        row_w = int(w * math.sqrt(max(0, 1 - (1-t)**2)))
        for x in range(cx - row_w, cx + row_w + 1):
            if 0 <= x < SIZE and 0 <= y < SIZE:
                nx_ = (x - cx) / max(row_w, 1)
                shade = light if nx_ < -0.3 else (dark if nx_ > 0.3 else color)
                px(img, x, y, shade)

def draw_leaf(img, base_x, base_y, tip_x, tip_y, width, color, highlight):
    """畫一片葉子（從底部到頂部，帶高光）"""
    steps = max(abs(tip_y - base_y), abs(tip_x - base_x), 1)
    for i in range(steps + 1):
        t = i / steps
        x = int(base_x + (tip_x - base_x) * t)
        y = int(base_y + (tip_y - base_y) * t)
        w = max(1, int(width * (1 - t * 0.7)))
        for dx in range(-w, w + 1):
            c = highlight if dx < 0 else color
            px(img, x + dx, y, c)

def draw_grass_blade(img, base_x, base_y, height, lean, color, highlight):
    """畫一根草葉（帶彎曲）"""
    for i in range(height):
        t = i / max(height, 1)
        x = base_x + int(lean * t * t)
        y = base_y - i
        w = max(1, int(2 * (1 - t * 0.6)))
        for dx in range(-w, w + 1):
            c = highlight if dx <= 0 else color
            px(img, x + dx, y, c)

# ── T001 像素雜草（帶土堆底座，密度目標 45%+）──────────────────────────────

def gen_T001_grass():
    img = new_img()
    OUTLINE  = (20, 60, 15, 255)
    STEM     = (50, 110, 30, 255)
    LEAF     = (70, 160, 45, 255)
    LEAF_L   = (90, 190, 55, 255)  # 高光
    LEAF_D   = (40, 90, 25, 255)   # 陰影

    # 土堆底座（讓草有根基感）
    draw_soil_mound(img, 32, 58, 22, 10, (110, 75, 35, 255))

    # 中央主莖（粗）
    for y in range(20, 52):
        w = 2 if y > 40 else 1
        for dx in range(-w, w+1):
            px(img, 32+dx, y, STEM)

    # 左側大葉片
    draw_grass_blade(img, 28, 50, 28, -6, LEAF, LEAF_L)
    draw_grass_blade(img, 24, 48, 22, -4, LEAF_D, LEAF, )

    # 右側大葉片
    draw_grass_blade(img, 36, 50, 26, 5, LEAF, LEAF_L)
    draw_grass_blade(img, 40, 47, 20, 4, LEAF_D, LEAF)

    # 中央頂部葉片（最高）
    draw_grass_blade(img, 32, 48, 32, 0, LEAF_L, (120, 220, 70, 255))

    # 小側葉（增加密度）
    draw_grass_blade(img, 26, 44, 14, -3, LEAF, LEAF_L)
    draw_grass_blade(img, 38, 44, 14, 3, LEAF, LEAF_L)
    draw_grass_blade(img, 30, 40, 10, -2, LEAF_D, LEAF)
    draw_grass_blade(img, 34, 40, 10, 2, LEAF_D, LEAF)

    # 輪廓加強
    outline_circle(img, 32, 35, 18, OUTLINE)

    return img

# ── T002/T003/T004 蟲類（帶陰影投影，密度目標 40%+）──────────────────────

def gen_bug(color_body, color_head, color_eye, color_leg, color_wing=None):
    """通用蟲類生成（帶陰影投影）"""
    img = new_img()
    OUTLINE = (20, 20, 20, 255)
    SHADOW  = (0, 0, 0, 60)  # 半透明陰影

    # 地面陰影（橢圓投影）
    for y in range(54, 60):
        t = (y - 54) / 6
        w = int(18 * (1 - t))
        for x in range(32 - w, 32 + w + 1):
            if 0 <= x < SIZE and 0 <= y < SIZE:
                existing = img.getpixel((x, y))
                img.putpixel((x, y), (0, 0, 0, int(50 * (1-t))))

    # 身體（大橢圓）
    fill_circle_shaded(img, 32, 36, 14, color_body)
    outline_circle(img, 32, 36, 14, OUTLINE)

    # 頭部
    fill_circle_shaded(img, 32, 20, 9, color_head)
    outline_circle(img, 32, 20, 9, OUTLINE)

    # 眼睛（2x2 白色 + 瞳孔）
    for ex, ey in [(28, 17), (36, 17)]:
        fill_circle(img, ex, ey, 3, (240, 240, 240, 255))
        fill_circle(img, ex+1, ey+1, 1, (20, 20, 20, 255))
        px(img, ex-1, ey-1, (255, 255, 255, 255))  # 高光

    # 觸角
    for i in range(8):
        px(img, 28 - i//2, 12 - i, OUTLINE)
        px(img, 36 + i//2, 12 - i, OUTLINE)
    fill_circle(img, 24, 5, 2, color_head)
    fill_circle(img, 40, 5, 2, color_head)

    # 腳（3對）
    for i, (lx, rx, y) in enumerate([(18,46,32),(16,48,38),(18,46,44)]):
        for dx in range(4):
            px(img, lx+dx, y+dx//2, color_leg)
            px(img, rx-dx, y+dx//2, color_leg)

    # 翅膀（若有）
    if color_wing:
        for side, cx in [(-1, 20), (1, 44)]:
            for y in range(26, 42):
                t = (y - 26) / 16
                w = int(8 * math.sin(t * math.pi))
                for x in range(cx - w, cx + w + 1):
                    if 0 <= x < SIZE:
                        px(img, x, y, (*color_wing[:3], 160))

    # 腹部條紋
    for i in range(3):
        y = 32 + i * 4
        for x in range(20, 44):
            if img.getpixel((x, y))[3] > 100:
                r, g, b, a = img.getpixel((x, y))
                px(img, x, y, (max(0,r-20), max(0,g-20), max(0,b-20), a))

    return img

def gen_T002_bug_g():
    return gen_bug((60,160,60), (40,120,40), (20,20,20,255), (30,100,30,255),
                   color_wing=(180,230,180))

def gen_T003_bug_r():
    return gen_bug((200,60,60), (160,40,40), (20,20,20,255), (120,30,30,255),
                   color_wing=(230,160,160))

def gen_T004_bug_b():
    return gen_bug((60,80,200), (40,60,160), (20,20,20,255), (30,50,140,255),
                   color_wing=(160,180,240))

# ── T101 擬態型怪物（基於 T001 草，但有詭異眼睛）──────────────────────────

def gen_T101_mimic():
    img = gen_T001_grass()  # 先畫草
    # 在草叢中加入詭異眼睛（讓玩家感覺「有東西藏在裡面」）
    # 左眼
    fill_circle(img, 27, 32, 4, (255, 255, 200, 255))
    fill_circle(img, 28, 33, 2, (200, 50, 50, 255))
    px(img, 27, 32, (255, 255, 255, 255))
    # 右眼
    fill_circle(img, 37, 32, 4, (255, 255, 200, 255))
    fill_circle(img, 38, 33, 2, (200, 50, 50, 255))
    px(img, 37, 32, (255, 255, 255, 255))
    # 詭異微笑
    for x in range(28, 37):
        px(img, x, 40, (20, 20, 20, 255))
    for x in [28, 36]:
        px(img, x, 39, (20, 20, 20, 255))
    return img

# ── T103 流星（帶完整光暈和尾焰）──────────────────────────────────────────

def gen_T103_meteor():
    img = new_img()
    # 尾焰（從右到左漸淡）
    for i in range(40):
        t = i / 40
        alpha = int(200 * (1 - t))
        r_c = int(255 * (1 - t * 0.3))
        g_c = int(180 * (1 - t * 0.5))
        b_c = int(80 * (1 - t * 0.8))
        x = 58 - i
        w = int(8 * (1 - t * 0.7))
        for dy in range(-w, w+1):
            y = 32 + dy
            if 0 <= x < SIZE and 0 <= y < SIZE:
                px(img, x, y, (r_c, g_c, b_c, alpha))

    # 外層光暈
    for r in range(16, 10, -1):
        alpha = int(80 * (16 - r) / 6)
        outline_circle(img, 20, 32, r, (255, 200, 100, alpha))

    # 主體（岩石核心）
    fill_circle_shaded(img, 20, 32, 12, (180, 120, 60))
    # 熔岩裂縫
    for x, y in [(18,28),(22,30),(16,34),(24,36),(20,32)]:
        fill_circle(img, x, y, 2, (255, 80, 20, 255))
    # 高光
    fill_circle(img, 15, 27, 3, (255, 240, 200, 200))
    # 輪廓
    outline_circle(img, 20, 32, 12, (80, 40, 10, 255))

    # 火花粒子
    sparks = [(8,20,255,180,50),(5,38,255,150,30),(12,44,255,200,80),
              (30,18,255,220,100),(35,42,255,160,40)]
    for sx, sy, sr, sg, sb in sparks:
        fill_circle(img, sx, sy, 2, (sr, sg, sb, 200))

    return img

# ── T104 金色雜草（帶土堆底座，金色光暈）──────────────────────────────────

def gen_T104_gold_grass():
    img = new_img()
    OUTLINE = (120, 80, 10, 255)
    STEM    = (180, 140, 30, 255)
    LEAF    = (220, 180, 40, 255)
    LEAF_L  = (255, 220, 80, 255)  # 高光
    LEAF_D  = (160, 120, 20, 255)  # 陰影

    # 金色土堆底座
    draw_soil_mound(img, 32, 58, 22, 10, (140, 100, 30, 255))

    # 金色光暈（底部）
    for r in range(20, 14, -1):
        alpha = int(40 * (20 - r) / 6)
        outline_circle(img, 32, 45, r, (255, 220, 50, alpha))

    # 主莖（粗，金色）
    for y in range(18, 52):
        w = 2 if y > 40 else 1
        for dx in range(-w, w+1):
            px(img, 32+dx, y, STEM)

    # 左側葉片（更寬更飽滿）
    draw_grass_blade(img, 27, 50, 30, -8, LEAF, LEAF_L)
    draw_grass_blade(img, 22, 47, 24, -5, LEAF_D, LEAF)
    draw_grass_blade(img, 18, 44, 18, -3, LEAF_D, LEAF)

    # 右側葉片
    draw_grass_blade(img, 37, 50, 30, 8, LEAF, LEAF_L)
    draw_grass_blade(img, 42, 47, 24, 5, LEAF_D, LEAF)
    draw_grass_blade(img, 46, 44, 18, 3, LEAF_D, LEAF)

    # 頂部穗狀花序（金色特徵）
    for i in range(5):
        x = 30 + i
        for y in range(8, 20):
            px(img, x, y, LEAF_L if y < 14 else LEAF)
        fill_circle(img, 30+i, 8, 2, (255, 240, 100, 255))

    # 中央頂葉
    draw_grass_blade(img, 32, 46, 34, 0, LEAF_L, (255, 240, 100, 255))

    # 小側葉
    draw_grass_blade(img, 26, 42, 14, -3, LEAF, LEAF_L)
    draw_grass_blade(img, 38, 42, 14, 3, LEAF, LEAF_L)

    # 輪廓
    outline_circle(img, 32, 33, 20, OUTLINE)

    return img

# ── 主程式 ──────────────────────────────────────────────────────────────────

def save_and_report(img, name):
    path = os.path.join(OUTPUT_DIR, name)
    img.save(path)
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 10)
    total = img.width * img.height
    pct = non_transparent / total * 100
    flag = "✅" if pct >= 30 else "⚠️"
    print(f"  {flag} {name}: {img.width}x{img.height}, {non_transparent}/{total} ({pct:.1f}%)")
    return pct

if __name__ == "__main__":
    print("=== 目標物 v5 生成（根本性密度提升）===\n")
    os.makedirs(OUTPUT_DIR, exist_ok=True)

    targets = [
        ("T001_grass.png",    gen_T001_grass),
        ("T002_bug_g.png",    gen_T002_bug_g),
        ("T003_bug_r.png",    gen_T003_bug_r),
        ("T004_bug_b.png",    gen_T004_bug_b),
        ("T101_mimic.png",    gen_T101_mimic),
        ("T103_meteor.png",   gen_T103_meteor),
        ("T104_gold_grass.png", gen_T104_gold_grass),
    ]

    total_pct = 0
    count = 0
    for name, gen_fn in targets:
        try:
            img = gen_fn()
            pct = save_and_report(img, name)
            total_pct += pct
            count += 1
        except Exception as e:
            print(f"  ❌ {name}: {e}")

    if count > 0:
        print(f"\n平均密度: {total_pct/count:.1f}%")
        print(f"目標: >= 40%")
        print(f"\n✅ 完成！輸出目錄: {OUTPUT_DIR}")
