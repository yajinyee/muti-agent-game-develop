"""
upgrade_basic_targets_day335.py — DAY-335 T001-T006 基礎目標物視覺特效升級
深度優先策略：讓最常出現的 6 個目標物有完整的視覺特效和手感

T001 草（雜草）：搖擺動畫 + 綠色粒子
T002 蟲G（綠蟲）：跳躍動畫 + 綠色光點
T003 蟲R（紅蟲）：快速移動 + 紅色尾跡
T004 蟲B（藍蟲）：Z字移動 + 藍色光點
T005 布丁：彈跳動畫 + 黃色光暈
T006 蘑菇：旋轉動畫 + 棕色粒子
"""
from PIL import Image, ImageDraw, ImageFilter
import os
import math
import random

OUTPUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"

def draw_pixel(draw, x, y, color, size=1):
    """畫一個像素點"""
    draw.rectangle([x, y, x+size-1, y+size-1], fill=color)

def generate_t001_grass_enhanced():
    """T001 草 — 更豐富的草叢，帶光澤感"""
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 草叢底部（深綠）
    for x in range(4, 28):
        draw_pixel(draw, x, 26, (30, 100, 30, 255))
        draw_pixel(draw, x, 27, (20, 80, 20, 255))
    
    # 草莖（多根，不同高度）
    stems = [
        (8, 10, (40, 160, 40, 255)),   # 左高莖
        (12, 14, (50, 180, 50, 255)),  # 左中莖
        (16, 8, (60, 200, 60, 255)),   # 中高莖（最高）
        (20, 13, (50, 175, 50, 255)),  # 右中莖
        (24, 11, (45, 165, 45, 255)),  # 右高莖
    ]
    for sx, sy, color in stems:
        for y in range(sy, 26):
            draw_pixel(draw, sx, y, color)
    
    # 草尖（亮綠色）
    tips = [(8, 9), (12, 13), (16, 7), (20, 12), (24, 10)]
    for tx, ty in tips:
        draw_pixel(draw, tx, ty, (120, 240, 80, 255))
        draw_pixel(draw, tx, ty-1, (150, 255, 100, 200))
    
    # 露珠（藍白色光點）
    dew_positions = [(10, 16), (18, 12), (22, 18)]
    for dx, dy in dew_positions:
        draw_pixel(draw, dx, dy, (180, 220, 255, 200))
        draw_pixel(draw, dx+1, dy, (220, 240, 255, 150))
    
    # 光澤高光
    draw_pixel(draw, 16, 8, (200, 255, 150, 220))
    
    return img

def generate_t002_bug_g_enhanced():
    """T002 綠蟲 — 可愛的綠色小蟲，帶觸角"""
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 身體（橢圓形，綠色）
    body_color = (60, 200, 60, 255)
    body_dark = (40, 150, 40, 255)
    body_light = (100, 240, 100, 255)
    
    # 主體
    for y in range(14, 22):
        for x in range(10, 22):
            dist = math.sqrt((x-16)**2 * 0.5 + (y-18)**2)
            if dist < 7:
                draw_pixel(draw, x, y, body_color)
    
    # 頭部
    for y in range(10, 17):
        for x in range(12, 20):
            dist = math.sqrt((x-16)**2 + (y-13)**2)
            if dist < 5:
                draw_pixel(draw, x, y, body_color)
    
    # 高光
    draw_pixel(draw, 14, 11, body_light)
    draw_pixel(draw, 15, 11, body_light)
    draw_pixel(draw, 14, 12, (150, 255, 150, 200))
    
    # 眼睛
    draw_pixel(draw, 14, 12, (0, 0, 0, 255))
    draw_pixel(draw, 18, 12, (0, 0, 0, 255))
    draw_pixel(draw, 14, 11, (255, 255, 255, 200))  # 眼白
    draw_pixel(draw, 18, 11, (255, 255, 255, 200))
    
    # 觸角
    draw_pixel(draw, 13, 9, body_dark)
    draw_pixel(draw, 12, 7, body_dark)
    draw_pixel(draw, 12, 6, (100, 220, 100, 255))  # 觸角尖
    draw_pixel(draw, 19, 9, body_dark)
    draw_pixel(draw, 20, 7, body_dark)
    draw_pixel(draw, 20, 6, (100, 220, 100, 255))
    
    # 腳（6隻）
    for i, (fx, fy) in enumerate([(10, 18), (10, 20), (10, 22), (22, 18), (22, 20), (22, 22)]):
        draw_pixel(draw, fx, fy, body_dark)
        if i < 3:
            draw_pixel(draw, fx-1, fy+1, body_dark)
        else:
            draw_pixel(draw, fx+1, fy+1, body_dark)
    
    return img

def generate_t003_bug_r_enhanced():
    """T003 紅蟲 — 快速移動的紅色蟲，帶速度感"""
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 速度線（尾跡）
    for i in range(3):
        alpha = 80 - i * 25
        for y in range(14, 20):
            draw_pixel(draw, 4 + i*2, y, (220, 60, 60, alpha))
    
    # 身體（流線型，紅色）
    body_color = (220, 60, 60, 255)
    body_light = (255, 120, 120, 255)
    
    # 主體（橢圓，稍微傾斜）
    for y in range(13, 21):
        for x in range(12, 26):
            dist = math.sqrt((x-19)**2 * 0.4 + (y-17)**2)
            if dist < 6:
                draw_pixel(draw, x, y, body_color)
    
    # 頭部（尖頭）
    for y in range(14, 19):
        for x in range(22, 28):
            dist = math.sqrt((x-24)**2 * 0.8 + (y-16)**2)
            if dist < 4:
                draw_pixel(draw, x, y, body_color)
    
    # 高光
    draw_pixel(draw, 20, 14, body_light)
    draw_pixel(draw, 21, 14, body_light)
    
    # 眼睛（紅色蟲的眼睛更兇）
    draw_pixel(draw, 24, 15, (0, 0, 0, 255))
    draw_pixel(draw, 24, 14, (255, 50, 50, 200))
    
    # 腳（快速移動姿態）
    for fy in [16, 18, 20]:
        draw_pixel(draw, 14, fy, (180, 40, 40, 255))
        draw_pixel(draw, 13, fy+1, (180, 40, 40, 255))
    
    return img

def generate_t004_bug_b_enhanced():
    """T004 藍蟲 — Z字移動的藍色蟲，帶電光感"""
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 電光尾跡
    lightning_points = [(6, 12), (8, 16), (6, 20), (8, 24)]
    for i in range(len(lightning_points)-1):
        x1, y1 = lightning_points[i]
        x2, y2 = lightning_points[i+1]
        draw_pixel(draw, x1, y1, (100, 150, 255, 120))
        draw_pixel(draw, (x1+x2)//2, (y1+y2)//2, (80, 120, 220, 80))
    
    # 身體（圓形，藍色）
    body_color = (60, 120, 220, 255)
    body_light = (120, 180, 255, 255)
    body_dark = (40, 80, 180, 255)
    
    for y in range(11, 23):
        for x in range(11, 23):
            dist = math.sqrt((x-16)**2 + (y-17)**2)
            if dist < 6:
                draw_pixel(draw, x, y, body_color)
    
    # 高光（電光感）
    draw_pixel(draw, 13, 13, body_light)
    draw_pixel(draw, 14, 13, body_light)
    draw_pixel(draw, 13, 14, (180, 220, 255, 200))
    
    # 眼睛（發光藍眼）
    draw_pixel(draw, 14, 15, (200, 230, 255, 255))
    draw_pixel(draw, 18, 15, (200, 230, 255, 255))
    draw_pixel(draw, 14, 14, (255, 255, 255, 200))
    draw_pixel(draw, 18, 14, (255, 255, 255, 200))
    
    # 電光觸角
    draw_pixel(draw, 13, 10, body_dark)
    draw_pixel(draw, 12, 8, (150, 200, 255, 255))
    draw_pixel(draw, 11, 7, (200, 230, 255, 200))
    draw_pixel(draw, 19, 10, body_dark)
    draw_pixel(draw, 20, 8, (150, 200, 255, 255))
    draw_pixel(draw, 21, 7, (200, 230, 255, 200))
    
    return img

def generate_t005_pudding_enhanced():
    """T005 布丁 — 可愛的彈跳布丁，帶光澤"""
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 布丁底部（深黃）
    for y in range(22, 27):
        for x in range(8, 24):
            dist = math.sqrt((x-16)**2 * 0.3 + (y-24)**2)
            if dist < 5:
                draw_pixel(draw, x, y, (200, 140, 20, 255))
    
    # 布丁主體（黃色圓頂）
    body_color = (255, 200, 50, 255)
    body_light = (255, 240, 150, 255)
    body_dark = (220, 160, 30, 255)
    
    for y in range(10, 24):
        for x in range(8, 24):
            dist = math.sqrt((x-16)**2 * 0.7 + (y-18)**2)
            if dist < 8:
                draw_pixel(draw, x, y, body_color)
    
    # 頂部奶油
    for y in range(7, 13):
        for x in range(12, 20):
            dist = math.sqrt((x-16)**2 + (y-10)**2)
            if dist < 4:
                draw_pixel(draw, x, y, (255, 250, 220, 255))
    
    # 頂部草莓
    draw_pixel(draw, 15, 7, (220, 50, 50, 255))
    draw_pixel(draw, 16, 7, (220, 50, 50, 255))
    draw_pixel(draw, 15, 8, (200, 40, 40, 255))
    draw_pixel(draw, 16, 6, (100, 180, 50, 255))  # 草莓葉
    
    # 高光（光澤感）
    draw_pixel(draw, 12, 13, body_light)
    draw_pixel(draw, 13, 12, body_light)
    draw_pixel(draw, 12, 12, (255, 255, 200, 220))
    
    # 眼睛（可愛）
    draw_pixel(draw, 13, 17, (80, 40, 0, 255))
    draw_pixel(draw, 19, 17, (80, 40, 0, 255))
    draw_pixel(draw, 13, 16, (255, 255, 255, 200))
    draw_pixel(draw, 19, 16, (255, 255, 255, 200))
    
    # 嘴巴（微笑）
    draw_pixel(draw, 14, 20, (180, 100, 20, 255))
    draw_pixel(draw, 15, 21, (180, 100, 20, 255))
    draw_pixel(draw, 16, 21, (180, 100, 20, 255))
    draw_pixel(draw, 17, 20, (180, 100, 20, 255))
    
    return img

def generate_t006_mushroom_enhanced():
    """T006 蘑菇 — 可愛的蘑菇，帶斑點和光澤"""
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 蘑菇莖（白色）
    stem_color = (240, 230, 210, 255)
    stem_dark = (200, 190, 170, 255)
    for y in range(20, 28):
        for x in range(12, 20):
            draw_pixel(draw, x, y, stem_color)
    # 莖的陰影
    for y in range(20, 28):
        draw_pixel(draw, 12, y, stem_dark)
        draw_pixel(draw, 19, y, stem_dark)
    
    # 蘑菇傘（棕紅色）
    cap_color = (160, 80, 40, 255)
    cap_light = (200, 120, 70, 255)
    cap_dark = (120, 50, 20, 255)
    
    for y in range(8, 22):
        for x in range(6, 26):
            dist = math.sqrt((x-16)**2 * 0.5 + (y-16)**2)
            if dist < 9 and y < 21:
                draw_pixel(draw, x, y, cap_color)
    
    # 傘的高光
    for y in range(9, 14):
        for x in range(10, 16):
            dist = math.sqrt((x-13)**2 + (y-11)**2)
            if dist < 3:
                draw_pixel(draw, x, y, cap_light)
    
    # 白色斑點（蘑菇特徵）
    spots = [(11, 12), (18, 11), (14, 16), (20, 15), (9, 16)]
    for sx, sy in spots:
        draw_pixel(draw, sx, sy, (255, 255, 255, 220))
        draw_pixel(draw, sx+1, sy, (255, 255, 255, 180))
        draw_pixel(draw, sx, sy+1, (255, 255, 255, 180))
    
    # 眼睛
    draw_pixel(draw, 14, 17, (60, 30, 10, 255))
    draw_pixel(draw, 18, 17, (60, 30, 10, 255))
    draw_pixel(draw, 14, 16, (255, 255, 255, 200))
    draw_pixel(draw, 18, 16, (255, 255, 255, 200))
    
    # 嘴巴
    draw_pixel(draw, 15, 19, (100, 50, 20, 255))
    draw_pixel(draw, 16, 19, (100, 50, 20, 255))
    draw_pixel(draw, 17, 19, (100, 50, 20, 255))
    
    return img

def main():
    generators = {
        "T001_grass": generate_t001_grass_enhanced,
        "T002_bug_g": generate_t002_bug_g_enhanced,
        "T003_bug_r": generate_t003_bug_r_enhanced,
        "T004_bug_b": generate_t004_bug_b_enhanced,
        "T005_pudding": generate_t005_pudding_enhanced,
        "T006_mushroom": generate_t006_mushroom_enhanced,
    }
    
    print("=== DAY-335 T001-T006 基礎目標物視覺升級 ===")
    success = 0
    for name, gen_func in generators.items():
        try:
            img = gen_func()
            # 放大到 64x64（2x 像素放大，保持像素風格）
            img_large = img.resize((64, 64), Image.NEAREST)
            path = os.path.join(OUTPUT_DIR, f"{name}.png")
            img_large.save(path)
            print(f"  ✅ {name}.png ({img_large.size[0]}x{img_large.size[1]})")
            success += 1
        except Exception as e:
            print(f"  ❌ {name}: {e}")
    
    print(f"\n完成：{success}/{len(generators)} 個目標物升級")
    print("注意：這些是 32x32 設計放大到 64x64 的像素圖")
    print("在 Godot 中使用 TEXTURE_FILTER_NEAREST 確保像素風格")

if __name__ == "__main__":
    main()
