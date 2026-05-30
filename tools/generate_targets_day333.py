#!/usr/bin/env python3
"""
generate_targets_day333.py — DAY-333 T249-T253 精靈圖生成
業界研究：Nolimit City Catfish Hunters 電擊框架 + Atomic Slot Lab Golden Gills 磁力連鎖
         Reflex Gaming Bigger Bites 漁夫路徑 + Penta Fusion 五重終極融合（2026）
"""
import os
import math
from PIL import Image, ImageDraw

OUTPUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = (64, 64)

def draw_glow_ring(draw, cx, cy, r, color, width=2):
    """畫光環"""
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], outline=color, width=width)

def draw_rays(draw, cx, cy, count, r_inner, r_outer, color, width=1):
    """畫放射光芒"""
    for i in range(count):
        angle = math.radians(i * 360 / count)
        x1 = cx + r_inner * math.cos(angle)
        y1 = cy + r_inner * math.sin(angle)
        x2 = cx + r_outer * math.cos(angle)
        y2 = cy + r_outer * math.sin(angle)
        draw.line([x1, y1, x2, y2], fill=color, width=width)

def generate_t249_electrical_frame():
    """T249 幸運電擊框架魚 — 青色魚身 + 電擊框架 + 全局倍率翻倍符號"""
    img = Image.new("RGBA", SIZE, (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 魚身（青色）
    draw.ellipse([14, 20, 50, 44], fill=(0, 200, 220, 255))
    draw.ellipse([16, 22, 48, 42], fill=(0, 230, 255, 255))

    # 電擊框架（矩形框）
    draw.rectangle([10, 16, 54, 48], outline=(0, 255, 255, 255), width=2)
    draw.rectangle([12, 18, 52, 46], outline=(100, 255, 255, 180), width=1)

    # 閃電符號（中央）
    pts = [(32, 20), (26, 32), (32, 30), (26, 44), (38, 30), (32, 32)]
    draw.polygon(pts, fill=(255, 255, 0, 255))

    # 全局倍率翻倍符號（×2）
    draw.text((44, 14), "×2", fill=(255, 255, 0, 255))

    # 光芒（10道）
    draw_rays(draw, cx, cy, 10, 26, 30, (0, 255, 255, 200), 1)

    # 光環
    draw_glow_ring(draw, cx, cy, 28, (0, 200, 255, 180), 1)
    draw_glow_ring(draw, cx, cy, 30, (0, 255, 255, 120), 1)

    return img

def generate_t250_magnetic_respin():
    """T250 幸運磁力連鎖魚 — 黃金色魚身 + 磁力弧線 + Respin符號"""
    img = Image.new("RGBA", SIZE, (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 魚身（黃金色）
    draw.ellipse([14, 20, 50, 44], fill=(180, 140, 0, 255))
    draw.ellipse([16, 22, 48, 42], fill=(255, 215, 0, 255))

    # 磁力弧線（U形磁鐵）
    draw.arc([20, 18, 44, 38], start=0, end=180, fill=(255, 100, 0, 255), width=3)
    draw.line([20, 28, 20, 38], fill=(255, 100, 0, 255), width=3)
    draw.line([44, 28, 44, 38], fill=(255, 100, 0, 255), width=3)
    # 磁極
    draw.rectangle([18, 36, 22, 40], fill=(255, 0, 0, 255))
    draw.rectangle([42, 36, 46, 40], fill=(0, 100, 255, 255))

    # 光芒（8道）
    draw_rays(draw, cx, cy, 8, 26, 30, (255, 215, 0, 200), 1)

    # 光環
    draw_glow_ring(draw, cx, cy, 28, (255, 180, 0, 180), 1)
    draw_glow_ring(draw, cx, cy, 30, (255, 215, 0, 120), 1)

    return img

def generate_t251_fisherman_trail():
    """T251 幸運漁夫路徑魚 — 橙色魚身 + 路徑節點 + 漁夫符號"""
    img = Image.new("RGBA", SIZE, (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 魚身（橙色）
    draw.ellipse([14, 20, 50, 44], fill=(180, 80, 0, 255))
    draw.ellipse([16, 22, 48, 42], fill=(255, 140, 0, 255))

    # 路徑（S形曲線節點）
    nodes = [(14, 38), (20, 28), (26, 36), (32, 26), (38, 34), (44, 24), (50, 32)]
    for i in range(len(nodes)-1):
        draw.line([nodes[i], nodes[i+1]], fill=(255, 220, 100, 200), width=2)
    for node in nodes:
        draw.ellipse([node[0]-3, node[1]-3, node[0]+3, node[1]+3], fill=(255, 255, 0, 255))

    # 漁夫符號（釣竿）
    draw.line([28, 16, 36, 24], fill=(200, 120, 0, 255), width=2)
    draw.line([36, 24, 36, 32], fill=(200, 120, 0, 200), width=1)

    # 光芒（12道）
    draw_rays(draw, cx, cy, 12, 26, 30, (255, 140, 0, 200), 1)

    # 光環
    draw_glow_ring(draw, cx, cy, 28, (255, 100, 0, 180), 1)
    draw_glow_ring(draw, cx, cy, 30, (255, 140, 0, 120), 1)

    return img

def generate_t252_golden_gills():
    """T252 幸運黃金鰓魚 — 黃金色大型魚身 + 4層Jackpot符號 + 鰓紋"""
    img = Image.new("RGBA", SIZE, (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 魚身（黃金色，較大）
    draw.ellipse([10, 18, 54, 46], fill=(160, 120, 0, 255))
    draw.ellipse([12, 20, 52, 44], fill=(255, 215, 0, 255))

    # 鰓紋（3條弧線）
    for i in range(3):
        x = 22 + i * 4
        draw.arc([x, 22, x+8, 42], start=270, end=90, fill=(200, 160, 0, 200), width=1)

    # Jackpot 4層符號（Mini/Minor/Major/Grand）
    jackpot_colors = [(180, 180, 180, 255), (0, 180, 255, 255), (180, 0, 255, 255), (255, 215, 0, 255)]
    for i, color in enumerate(jackpot_colors):
        r = 8 + i * 4
        draw.ellipse([cx-r, cy-r, cx+r, cy+r], outline=color, width=1)

    # 光芒（16道）
    draw_rays(draw, cx, cy, 16, 28, 31, (255, 215, 0, 200), 1)

    # 光環（三層）
    draw_glow_ring(draw, cx, cy, 29, (255, 180, 0, 180), 1)
    draw_glow_ring(draw, cx, cy, 31, (255, 215, 0, 120), 1)

    return img

def generate_t253_penta_fusion():
    """T253 幸運五重終極魚 — 熱粉紅超大型魚身 + 五象限符號 + 30道光芒 + 七層光環"""
    img = Image.new("RGBA", SIZE, (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 魚身（熱粉紅，超大型）
    draw.ellipse([8, 16, 56, 48], fill=(140, 0, 80, 255))
    draw.ellipse([10, 18, 54, 46], fill=(255, 105, 180, 255))

    # 五象限符號（五角星）
    pts = []
    for i in range(5):
        angle = math.radians(-90 + i * 72)
        pts.append((cx + 12 * math.cos(angle), cy + 12 * math.sin(angle)))
    # 五角星連線
    order = [0, 2, 4, 1, 3, 0]
    for i in range(len(order)-1):
        draw.line([pts[order[i]], pts[order[i+1]]], fill=(255, 255, 0, 255), width=2)

    # 五個相位符號（電擊/磁力/路徑/鰓/融合）
    phase_colors = [(0, 255, 255, 200), (255, 215, 0, 200), (255, 140, 0, 200), (255, 215, 0, 200), (255, 0, 255, 200)]
    for i, color in enumerate(phase_colors):
        angle = math.radians(-90 + i * 72)
        px = cx + 16 * math.cos(angle)
        py = cy + 16 * math.sin(angle)
        draw.ellipse([px-3, py-3, px+3, py+3], fill=color)

    # 光芒（30道）
    draw_rays(draw, cx, cy, 30, 26, 30, (255, 105, 180, 200), 1)

    # 光環（七層）
    for r, alpha in [(22, 200), (24, 180), (26, 160), (28, 140), (29, 120), (30, 100), (31, 80)]:
        draw_glow_ring(draw, cx, cy, r, (255, 105, 180, alpha), 1)

    return img

def save_with_import(img, filename):
    """儲存圖片並建立 .import 檔案"""
    path = os.path.join(OUTPUT_DIR, filename)
    img.save(path)
    print(f"  ✅ {filename} ({img.size[0]}x{img.size[1]})")

    # 建立 .import 檔案
    import_path = path + ".import"
    import_content = f"""[remap]

importer="texture"
type="CompressedTexture2D"
uid="uid://{filename.replace('.png','').replace('_','').lower()}"
path="res://.godot/imported/{filename}-{hash(filename) & 0xFFFFFFFF:08x}.ctex"
metadata={{
"vram_texture": false
}}

[deps]

source_file="res://assets/sprites/targets/{filename}"
dest_files=["res://.godot/imported/{filename}-{hash(filename) & 0xFFFFFFFF:08x}.ctex"]

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
        f.write(import_content)

def main():
    print("=== DAY-333 T249-T253 精靈圖生成 ===")
    print(f"輸出目錄：{OUTPUT_DIR}")

    targets = [
        ("T249_electrical_frame.png", generate_t249_electrical_frame),
        ("T250_magnetic_respin.png",  generate_t250_magnetic_respin),
        ("T251_fisherman_trail.png",  generate_t251_fisherman_trail),
        ("T252_golden_gills.png",     generate_t252_golden_gills),
        ("T253_penta_fusion.png",     generate_t253_penta_fusion),
    ]

    for filename, gen_func in targets:
        img = gen_func()
        save_with_import(img, filename)

    print("\n✅ 全部完成！T249-T253 精靈圖生成完畢")
    print("業界依據：")
    print("  T249: Nolimit City Catfish Hunters 電擊框架（2026-03）")
    print("  T250: Atomic Slot Lab Golden Gills 磁力連鎖（2026-02）")
    print("  T251: Reflex Gaming Bigger Bites 漁夫路徑（2026-02）")
    print("  T252: Atomic Slot Lab Golden Gills Jackpot（2026-02）")
    print("  T253: Penta Fusion 五重終極融合（2026 里程碑）")

if __name__ == "__main__":
    main()
