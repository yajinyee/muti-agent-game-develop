"""
generate_targets_day314.py — DAY-314 T176-T180 精靈圖生成
T176 多重宇宙魚（深紫色 + 三個宇宙環 + 星際粒子）
T177 時間迴圈魚（深藍色 + 時鐘符號 + 迴圈箭頭）
T178 命運之輪魚（火橙色 + 輪盤扇形 + 指針）
T179 神域降臨魚（神聖橙金色 + 5 道光柱 + 神聖光環）
T180 終焉之力魚（深紅色 + 16 道光芒 + 三層光環 + 骷髏符號）
"""
import os
import math
from PIL import Image, ImageDraw

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def fill_circle(draw, cx, cy, r, color):
    for y in range(max(0, cy-r), min(SIZE, cy+r+1)):
        for x in range(max(0, cx-r), min(SIZE, cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 <= r*r:
                draw.point((x, y), fill=color)

def draw_ring(draw, cx, cy, r_outer, r_inner, color):
    for y in range(max(0, cy-r_outer), min(SIZE, cy+r_outer+1)):
        for x in range(max(0, cx-r_outer), min(SIZE, cx+r_outer+1)):
            d2 = (x-cx)**2 + (y-cy)**2
            if r_inner*r_inner <= d2 <= r_outer*r_outer:
                draw.point((x, y), fill=color)

def draw_ray(draw, cx, cy, angle_deg, length, width, color):
    angle = math.radians(angle_deg)
    for i in range(length):
        bx = cx + int(i * math.cos(angle))
        by = cy + int(i * math.sin(angle))
        for dx in range(-width//2, width//2+1):
            for dy in range(-width//2, width//2+1):
                nx, ny = bx+dx, by+dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    draw.point((nx, ny), fill=color)

def make_t176_multiverse():
    """T176 多重宇宙魚：深紫色魚身 + 三個宇宙環 + 星際粒子"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32
    # 魚身（深紫橢圓）
    for y in range(cy-14, cy+15):
        for x in range(cx-20, cx+21):
            if ((x-cx)/20)**2 + ((y-cy)/14)**2 <= 1.0:
                shade = 1.0 - 0.3 * (((x-cx)/20)**2 + ((y-cy)/14)**2)
                r = int(100 * shade)
                g = int(20 * shade)
                b = int(160 * shade)
                draw.point((x, y), fill=(r, g, b, 255))
    # 三個宇宙環（不同半徑）
    for r_outer, r_inner, alpha in [(28, 25, 180), (22, 19, 140), (16, 13, 100)]:
        draw_ring(draw, cx, cy, r_outer, r_inner, (180, 100, 255, alpha))
    # 星際粒子（8個）
    import random
    rng = random.Random(314)
    for _ in range(12):
        angle = rng.uniform(0, 2*math.pi)
        dist = rng.uniform(10, 26)
        px = int(cx + dist * math.cos(angle))
        py = int(cy + dist * math.sin(angle))
        if 0 <= px < SIZE and 0 <= py < SIZE:
            draw.point((px, py), fill=(220, 180, 255, 255))
    # 眼睛
    draw.point((cx+6, cy-3), fill=(255, 255, 255, 255))
    draw.point((cx+7, cy-3), fill=(255, 255, 255, 255))
    draw.point((cx+6, cy-2), fill=(255, 255, 255, 255))
    draw.point((cx+7, cy-2), fill=(255, 255, 255, 255))
    draw.point((cx+6, cy-3), fill=(200, 150, 255, 255))
    return img

def make_t177_time_loop():
    """T177 時間迴圈魚：深藍色魚身 + 時鐘符號 + 迴圈箭頭"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32
    # 魚身（深藍橢圓）
    for y in range(cy-13, cy+14):
        for x in range(cx-19, cx+20):
            if ((x-cx)/19)**2 + ((y-cy)/13)**2 <= 1.0:
                shade = 1.0 - 0.3 * (((x-cx)/19)**2 + ((y-cy)/13)**2)
                r = int(20 * shade)
                g = int(80 * shade)
                b = int(180 * shade)
                draw.point((x, y), fill=(r, g, b, 255))
    # 時鐘圓環
    draw_ring(draw, cx, cy, 18, 15, (100, 180, 255, 200))
    # 時鐘指針（12點和3點）
    draw_ray(draw, cx, cy, -90, 12, 1, (200, 230, 255, 255))  # 12點
    draw_ray(draw, cx, cy, 0, 8, 1, (200, 230, 255, 255))     # 3點
    # 迴圈箭頭（外圈）
    for angle in range(0, 300, 15):
        a = math.radians(angle)
        px = int(cx + 22 * math.cos(a))
        py = int(cy + 22 * math.sin(a))
        if 0 <= px < SIZE and 0 <= py < SIZE:
            draw.point((px, py), fill=(150, 200, 255, 180))
    # 眼睛
    draw.point((cx+6, cy-3), fill=(255, 255, 255, 255))
    draw.point((cx+7, cy-3), fill=(255, 255, 255, 255))
    draw.point((cx+6, cy-2), fill=(255, 255, 255, 255))
    draw.point((cx+7, cy-2), fill=(255, 255, 255, 255))
    draw.point((cx+6, cy-3), fill=(100, 180, 255, 255))
    return img

def make_t178_fate_wheel():
    """T178 命運之輪魚：火橙色魚身 + 輪盤扇形 + 指針"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32
    # 魚身（火橙橢圓）
    for y in range(cy-13, cy+14):
        for x in range(cx-19, cx+20):
            if ((x-cx)/19)**2 + ((y-cy)/13)**2 <= 1.0:
                shade = 1.0 - 0.3 * (((x-cx)/19)**2 + ((y-cy)/13)**2)
                r = int(240 * shade)
                g = int(120 * shade)
                b = int(20 * shade)
                draw.point((x, y), fill=(r, g, b, 255))
    # 輪盤（8個扇形，交替顏色）
    colors = [(255, 200, 50, 180), (255, 100, 20, 180)]
    for i in range(8):
        angle_start = i * 45
        angle_end = (i + 1) * 45
        color = colors[i % 2]
        for angle in range(angle_start, angle_end, 2):
            a = math.radians(angle)
            for r in range(12, 20):
                px = int(cx + r * math.cos(a))
                py = int(cy + r * math.sin(a))
                if 0 <= px < SIZE and 0 <= py < SIZE:
                    draw.point((px, py), fill=color)
    # 中心圓
    fill_circle(draw, cx, cy, 5, (255, 220, 100, 255))
    # 指針
    draw_ray(draw, cx, cy, -60, 18, 1, (255, 255, 255, 255))
    # 眼睛
    draw.point((cx+6, cy-3), fill=(255, 255, 255, 255))
    draw.point((cx+7, cy-3), fill=(255, 255, 255, 255))
    draw.point((cx+6, cy-2), fill=(255, 255, 255, 255))
    draw.point((cx+7, cy-2), fill=(255, 255, 255, 255))
    draw.point((cx+6, cy-3), fill=(255, 150, 50, 255))
    return img

def make_t179_divine_realm():
    """T179 神域降臨魚：神聖橙金色魚身 + 5 道光柱 + 神聖光環"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32
    # 背景光暈
    for y in range(SIZE):
        for x in range(SIZE):
            d = math.sqrt((x-cx)**2 + (y-cy)**2)
            if d < 28:
                alpha = int(60 * (1 - d/28))
                draw.point((x, y), fill=(255, 200, 50, alpha))
    # 5 道神聖光柱（均勻分布）
    for i in range(5):
        angle = i * 72 - 90
        draw_ray(draw, cx, cy, angle, 28, 2, (255, 220, 100, 200))
    # 魚身（神聖橙金大型橢圓）
    for y in range(cy-15, cy+16):
        for x in range(cx-22, cx+23):
            if ((x-cx)/22)**2 + ((y-cy)/15)**2 <= 1.0:
                shade = 1.0 - 0.25 * (((x-cx)/22)**2 + ((y-cy)/15)**2)
                r = int(240 * shade)
                g = int(160 * shade)
                b = int(30 * shade)
                draw.point((x, y), fill=(r, g, b, 255))
    # 神聖光環
    draw_ring(draw, cx, cy, 24, 21, (255, 230, 120, 200))
    draw_ring(draw, cx, cy, 20, 18, (255, 200, 80, 160))
    # 眼睛
    draw.point((cx+7, cy-4), fill=(255, 255, 255, 255))
    draw.point((cx+8, cy-4), fill=(255, 255, 255, 255))
    draw.point((cx+7, cy-3), fill=(255, 255, 255, 255))
    draw.point((cx+8, cy-3), fill=(255, 255, 255, 255))
    draw.point((cx+7, cy-4), fill=(255, 180, 50, 255))
    return img

def make_t180_final_power():
    """T180 終焉之力魚：深紅色大型魚身 + 16 道光芒 + 三層光環 + 骷髏符號"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32
    # 背景暗紅光暈
    for y in range(SIZE):
        for x in range(SIZE):
            d = math.sqrt((x-cx)**2 + (y-cy)**2)
            if d < 30:
                alpha = int(80 * (1 - d/30))
                draw.point((x, y), fill=(180, 0, 0, alpha))
    # 16 道光芒
    for i in range(16):
        angle = i * 22.5
        draw_ray(draw, cx, cy, angle, 30, 1, (255, 50, 50, 180))
    # 魚身（深紅大型橢圓）
    for y in range(cy-16, cy+17):
        for x in range(cx-23, cx+24):
            if ((x-cx)/23)**2 + ((y-cy)/16)**2 <= 1.0:
                shade = 1.0 - 0.2 * (((x-cx)/23)**2 + ((y-cy)/16)**2)
                r = int(200 * shade)
                g = int(15 * shade)
                b = int(15 * shade)
                draw.point((x, y), fill=(r, g, b, 255))
    # 三層光環
    draw_ring(draw, cx, cy, 26, 23, (255, 80, 80, 200))
    draw_ring(draw, cx, cy, 22, 20, (200, 50, 50, 160))
    draw_ring(draw, cx, cy, 18, 16, (150, 30, 30, 120))
    # 骷髏符號（簡化版：兩個眼睛圓 + 鼻子 + 牙齒）
    fill_circle(draw, cx-4, cy-2, 3, (255, 200, 200, 220))
    fill_circle(draw, cx+4, cy-2, 3, (255, 200, 200, 220))
    draw.point((cx, cy+2), fill=(255, 200, 200, 200))
    for tx in range(cx-5, cx+6, 2):
        draw.point((tx, cy+5), fill=(255, 200, 200, 200))
    # 眼睛（魚眼）
    draw.point((cx+8, cy-5), fill=(255, 255, 255, 255))
    draw.point((cx+9, cy-5), fill=(255, 255, 255, 255))
    draw.point((cx+8, cy-4), fill=(255, 255, 255, 255))
    draw.point((cx+9, cy-4), fill=(255, 255, 255, 255))
    draw.point((cx+8, cy-5), fill=(255, 50, 50, 255))
    return img

def save_with_import(img, filename):
    path = os.path.join(OUT_DIR, filename)
    img.save(path)
    # 建立 .import 檔案
    import_path = path + ".import"
    import_content = f"""[remap]

importer="texture"
type="CompressedTexture2D"
uid="uid://auto_{filename.replace('.', '_')}"
path="res://.godot/imported/{filename}-auto.ctex"
metadata={{
"vram_texture": false
}}

[deps]

source_file="res://assets/sprites/targets/{filename}"
dest_files=["res://.godot/imported/{filename}-auto.ctex"]

[params]

compress/mode=0
compress/high_quality=false
compress/lossy_quality=0.7
compress/hdr_compression=1
compress/normal_map=0
compress/channel_pack=0
mipmaps/generate=false
mipmaps/limit=-1
roughness/mode=0
roughness/src_normal=""
process/fix_alpha_border=true
process/premult_alpha=false
process/normal_map_invert_y=false
process/hdr_as_srgb=false
process/hdr_clamp_exposure=false
process/size_limit=0
detect_3d/compress_to=1
svg/scale=1.0
editor/scale_with_editor_scale=false
editor/convert_colors_with_editor_theme=false
"""
    with open(import_path, 'w', encoding='utf-8') as f:
        f.write(import_content)
    # 計算非透明像素
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 10)
    pct = non_transparent / (SIZE * SIZE) * 100
    print(f"  {filename}: {non_transparent} 非透明像素 ({pct:.1f}%)")

if __name__ == "__main__":
    os.makedirs(OUT_DIR, exist_ok=True)
    print("DAY-314 T176-T180 精靈圖生成中...")
    
    sprites = [
        ("T176_multiverse.png", make_t176_multiverse()),
        ("T177_time_loop.png", make_t177_time_loop()),
        ("T178_fate_wheel.png", make_t178_fate_wheel()),
        ("T179_divine_realm.png", make_t179_divine_realm()),
        ("T180_final_power.png", make_t180_final_power()),
    ]
    
    for filename, img in sprites:
        save_with_import(img, filename)
    
    print("完成！")
