"""
generate_targets_day332.py — DAY-332 T244-T248 精靈圖生成
T244 幸運野生收集魚（黃金色 Wild Collector）
T245 幸運閃電鰻升級魚（青色 Lightning Eel Ultra）
T246 幸運骨牌連鎖魚（橙色 Domino Chain）
T247 幸運不死BOSS升級魚（深紅色 Immortal Boss Ultra）
T248 幸運四重終極融合魚（洋紅色 Quad Fusion）
"""
import os
import math
from PIL import Image, ImageDraw

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def fill_circle(draw, cx, cy, r, color):
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], fill=color)

def draw_rays(draw, cx, cy, count, r_inner, r_outer, color, width=2):
    for i in range(count):
        angle = math.radians(i * 360 / count)
        x1 = cx + r_inner * math.cos(angle)
        y1 = cy + r_inner * math.sin(angle)
        x2 = cx + r_outer * math.cos(angle)
        y2 = cy + r_outer * math.sin(angle)
        draw.line([x1, y1, x2, y2], fill=color, width=width)

def draw_ring(draw, cx, cy, r, color, width=2):
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], outline=color, width=width)

def make_fish_body(draw, cx, cy, rx, ry, color_main, color_dark, color_light):
    """畫橢圓魚身（帶陰影）"""
    draw.ellipse([cx-rx, cy-ry, cx+rx, cy+ry], fill=color_dark)
    draw.ellipse([cx-rx+2, cy-ry+2, cx+rx-2, cy+ry-2], fill=color_main)
    draw.ellipse([cx-rx+4, cy-ry+4, cx+rx//2, cy+ry//2], fill=color_light)

def gen_t244_wild_collector():
    """T244 幸運野生收集魚 — 黃金色 Wild Collector"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2
    # 黃金色魚身
    make_fish_body(draw, cx, cy, 22, 16, (255, 215, 0), (200, 160, 0), (255, 240, 100))
    # Wild 符號（W 字）
    draw.text((cx-8, cy-6), "W", fill=(255, 255, 255, 230))
    # 16 道金色光芒
    draw_rays(draw, cx, cy, 16, 18, 28, (255, 220, 50, 200), 2)
    # 三層光環
    draw_ring(draw, cx, cy, 26, (255, 200, 0, 180), 2)
    draw_ring(draw, cx, cy, 29, (255, 180, 0, 120), 1)
    draw_ring(draw, cx, cy, 31, (255, 160, 0, 80), 1)
    # 四個角落的 Wild 星星
    for dx, dy in [(-18, -18), (18, -18), (-18, 18), (18, 18)]:
        fill_circle(draw, cx+dx, cy+dy, 3, (255, 255, 100, 200))
    return img

def gen_t245_lightning_eel_ultra():
    """T245 幸運閃電鰻升級魚 — 青色 Lightning Eel Ultra"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2
    # 細長青色魚身（鰻魚形）
    draw.ellipse([cx-24, cy-10, cx+24, cy+10], fill=(0, 160, 180))
    draw.ellipse([cx-22, cy-8, cx+22, cy+8], fill=(0, 220, 240))
    draw.ellipse([cx-18, cy-5, cx+10, cy+5], fill=(100, 255, 255))
    # 8 條閃電（從中心向外）
    for i in range(8):
        angle = math.radians(i * 45)
        x1 = cx + 10 * math.cos(angle)
        y1 = cy + 10 * math.sin(angle)
        x2 = cx + 26 * math.cos(angle)
        y2 = cy + 26 * math.sin(angle)
        # 鋸齒閃電
        mx = (x1 + x2) / 2 + 3 * math.cos(angle + math.pi/2)
        my = (y1 + y2) / 2 + 3 * math.sin(angle + math.pi/2)
        draw.line([x1, y1, mx, my, x2, y2], fill=(0, 255, 255, 220), width=2)
    # 電弧紋路
    draw_ring(draw, cx, cy, 20, (0, 200, 255, 180), 2)
    draw_ring(draw, cx, cy, 24, (0, 180, 220, 120), 1)
    draw_ring(draw, cx, cy, 28, (0, 160, 200, 80), 1)
    # 跳躍點（3個）
    for i in range(3):
        angle = math.radians(i * 120)
        px = cx + 22 * math.cos(angle)
        py = cy + 22 * math.sin(angle)
        fill_circle(draw, int(px), int(py), 3, (255, 255, 0, 220))
    return img

def gen_t246_domino_chain():
    """T246 幸運骨牌連鎖魚 — 橙色 Domino Chain"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2
    # 橙色魚身
    make_fish_body(draw, cx, cy, 20, 15, (255, 140, 0), (200, 100, 0), (255, 180, 80))
    # 骨牌符號（兩個矩形）
    draw.rectangle([cx-8, cy-7, cx-2, cy+7], fill=(255, 255, 255, 200), outline=(0, 0, 0, 180))
    draw.rectangle([cx+2, cy-7, cx+8, cy+7], fill=(255, 255, 255, 200), outline=(0, 0, 0, 180))
    # 骨牌點數
    fill_circle(draw, cx-5, cy-3, 2, (0, 0, 0, 200))
    fill_circle(draw, cx-5, cy+3, 2, (0, 0, 0, 200))
    fill_circle(draw, cx+5, cy, 2, (0, 0, 0, 200))
    # 連鎖箭頭
    draw.line([cx-12, cy, cx-9, cy], fill=(255, 200, 0, 220), width=2)
    draw.line([cx+9, cy, cx+12, cy], fill=(255, 200, 0, 220), width=2)
    # 12 道光芒
    draw_rays(draw, cx, cy, 12, 18, 27, (255, 160, 0, 180), 2)
    # 光環
    draw_ring(draw, cx, cy, 25, (255, 140, 0, 160), 2)
    draw_ring(draw, cx, cy, 28, (255, 120, 0, 100), 1)
    draw_ring(draw, cx, cy, 30, (255, 100, 0, 60), 1)
    return img

def gen_t247_immortal_boss_ultra():
    """T247 幸運不死BOSS升級魚 — 深紅色 Immortal Boss Ultra"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2
    # 深紅色大型魚身
    make_fish_body(draw, cx, cy, 22, 17, (180, 0, 0), (120, 0, 0), (220, 50, 50))
    # 骷髏符號（簡化）
    fill_circle(draw, cx, cy-2, 7, (255, 255, 255, 200))
    fill_circle(draw, cx-3, cy+4, 3, (255, 255, 255, 180))
    fill_circle(draw, cx+3, cy+4, 3, (255, 255, 255, 180))
    # 眼睛（紅色）
    fill_circle(draw, cx-3, cy-3, 2, (255, 0, 0, 220))
    fill_circle(draw, cx+3, cy-3, 2, (255, 0, 0, 220))
    # 5 次復活光環（5層）
    for i, r in enumerate([20, 22, 24, 26, 28]):
        alpha = 200 - i * 30
        draw_ring(draw, cx, cy, r, (200, 0, 0, alpha), 1)
    # 16 道暗紅光芒
    draw_rays(draw, cx, cy, 16, 19, 29, (180, 0, 0, 160), 2)
    return img

def gen_t248_quad_fusion():
    """T248 幸運四重終極融合魚 — 洋紅色 Quad Fusion（里程碑）"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE//2, SIZE//2
    # 超大型洋紅色魚身
    make_fish_body(draw, cx, cy, 24, 18, (200, 0, 200), (140, 0, 140), (255, 100, 255))
    # 四象限符號（代表四個機制）
    # 左上：Wild（黃色）
    fill_circle(draw, cx-8, cy-8, 4, (255, 215, 0, 220))
    # 右上：Eel（青色）
    fill_circle(draw, cx+8, cy-8, 4, (0, 255, 255, 220))
    # 左下：Domino（橙色）
    fill_circle(draw, cx-8, cy+8, 4, (255, 140, 0, 220))
    # 右下：Boss（深紅）
    fill_circle(draw, cx+8, cy+8, 4, (200, 0, 0, 220))
    # 中心融合點
    fill_circle(draw, cx, cy, 5, (255, 255, 255, 240))
    # 四條融合線
    for dx, dy in [(-8, -8), (8, -8), (-8, 8), (8, 8)]:
        draw.line([cx, cy, cx+dx, cy+dy], fill=(255, 255, 255, 180), width=2)
    # 24 道洋紅光芒（里程碑）
    draw_rays(draw, cx, cy, 24, 20, 30, (255, 0, 255, 180), 2)
    # 六層光環（里程碑）
    for i, r in enumerate([22, 24, 26, 28, 30, 31]):
        alpha = 200 - i * 25
        draw_ring(draw, cx, cy, r, (220, 0, 220, alpha), 1)
    return img

def save_with_import(img, filename):
    path = os.path.join(OUT_DIR, filename)
    img.save(path)
    import_path = path + ".import"
    base = filename.replace(".png", "")
    content = f"""[remap]

importer="texture"
type="CompressedTexture2D"
uid="uid://b{abs(hash(base)) % 10**15}"
path="res://.godot/imported/{filename}-{abs(hash(base)) % 10**15}.ctex"
metadata={{
"vram_texture": false
}}

[deps]

source_file="res://assets/sprites/targets/{filename}"
dest_files=["res://.godot/imported/{filename}-{abs(hash(base)) % 10**15}.ctex"]

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
    with open(import_path, "w", encoding="utf-8") as f:
        f.write(content)
    print(f"  ✅ {filename} ({img.size[0]}x{img.size[1]})")

if __name__ == "__main__":
    print("=== DAY-332 T244-T248 精靈圖生成 ===")
    targets = [
        ("T244_wild_collector.png", gen_t244_wild_collector),
        ("T245_lightning_eel_ultra.png", gen_t245_lightning_eel_ultra),
        ("T246_domino_chain.png", gen_t246_domino_chain),
        ("T247_immortal_boss_ultra.png", gen_t247_immortal_boss_ultra),
        ("T248_quad_fusion.png", gen_t248_quad_fusion),
    ]
    for filename, gen_func in targets:
        img = gen_func()
        save_with_import(img, filename)
    print(f"\n✅ 全部完成！共生成 {len(targets)} 個精靈圖 + {len(targets)} 個 .import 檔案")
