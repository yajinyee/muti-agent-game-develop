#!/usr/bin/env python3
"""
generate_targets_day313.py — DAY-313 Progressive Jackpot 系列精靈圖生成
T171-T175：Mini/Minor/Major/Grand Jackpot 魚 + Jackpot Trigger 魚
progressive-jackpot-agent 負責維護
"""
import os
from PIL import Image, ImageDraw
import math
import random

OUTPUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = (64, 64)

def draw_jackpot_fish(draw, color_main, color_accent, tier_symbol, tier_level):
    """繪製 Jackpot 系列魚的基礎形狀"""
    cx, cy = 32, 32
    
    # 外層光暈（依層級增大）
    halo_r = 28 + tier_level * 2
    for i in range(3):
        alpha = int(60 - i * 15)
        r = halo_r + i * 3
        draw.ellipse([cx-r, cy-r, cx+r, cy+r], 
                     fill=(*color_accent[:3], alpha))
    
    # 魚身（橢圓）
    body_w = 22 + tier_level * 2
    body_h = 16 + tier_level
    draw.ellipse([cx-body_w, cy-body_h, cx+body_w, cy+body_h], 
                 fill=(*color_main, 230))
    
    # 魚尾
    tail_pts = [
        (cx + body_w - 4, cy),
        (cx + body_w + 10, cy - 8),
        (cx + body_w + 10, cy + 8),
    ]
    draw.polygon(tail_pts, fill=(*color_accent, 200))
    
    # 魚眼
    draw.ellipse([cx - body_w + 6, cy - 4, cx - body_w + 12, cy + 2], 
                 fill=(255, 255, 255, 255))
    draw.ellipse([cx - body_w + 8, cy - 3, cx - body_w + 11, cy + 1], 
                 fill=(0, 0, 0, 255))
    
    # 層級光芒（依層級增加光芒數量）
    ray_count = 4 + tier_level * 2
    ray_len = 12 + tier_level * 3
    for i in range(ray_count):
        angle = (2 * math.pi * i / ray_count)
        x1 = cx + int(body_w * 0.6 * math.cos(angle))
        y1 = cy + int(body_h * 0.6 * math.sin(angle))
        x2 = cx + int((body_w + ray_len) * math.cos(angle))
        y2 = cy + int((body_h + ray_len) * math.sin(angle))
        draw.line([x1, y1, x2, y2], fill=(*color_accent, 180), width=2)
    
    return body_w, body_h

def generate_t171_jackpot_mini():
    """T171 幸運 Mini Jackpot 魚 — 綠色，4 道光芒"""
    img = Image.new("RGBA", SIZE, (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    color_main = (30, 180, 30)    # 綠色
    color_accent = (100, 255, 100)  # 亮綠
    
    # 外層大光暈
    cx, cy = 32, 32
    for r, a in [(30, 50), (26, 80), (22, 100)]:
        draw.ellipse([cx-r, cy-r, cx+r, cy+r], fill=(*color_accent, a))
    
    draw_jackpot_fish(draw, color_main, color_accent, "M", 0)
    
    # Mini 標誌：小星星
    for i in range(4):
        angle = math.pi / 4 + i * math.pi / 2
        sx = cx + int(8 * math.cos(angle))
        sy = cy + int(8 * math.sin(angle))
        draw.ellipse([sx-3, sy-3, sx+3, sy+3], fill=(255, 255, 100, 220))
    
    # 中心圓
    draw.ellipse([28, 28, 36, 36], fill=(255, 255, 255, 200))
    
    img.save(os.path.join(OUTPUT_DIR, "T171_jackpot_mini.png"))
    print(f"T171 生成完成（Mini Jackpot 綠色）")
    return img

def generate_t172_jackpot_minor():
    """T172 幸運 Minor Jackpot 魚 — 藍色，6 道光芒"""
    img = Image.new("RGBA", SIZE, (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    color_main = (30, 100, 220)    # 藍色
    color_accent = (100, 180, 255)  # 亮藍
    
    draw_jackpot_fish(draw, color_main, color_accent, "m", 1)
    
    # Minor 標誌：藍色菱形
    cx, cy = 32, 32
    diamond_pts = [(cx, cy-6), (cx+5, cy), (cx, cy+6), (cx-5, cy)]
    draw.polygon(diamond_pts, fill=(200, 230, 255, 220))
    
    img.save(os.path.join(OUTPUT_DIR, "T172_jackpot_minor.png"))
    print(f"T172 生成完成（Minor Jackpot 藍色）")
    return img

def generate_t173_jackpot_major():
    """T173 幸運 Major Jackpot 魚 — 橙色，8 道光芒，較大"""
    img = Image.new("RGBA", SIZE, (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    color_main = (220, 120, 20)    # 橙色
    color_accent = (255, 200, 50)  # 金橙
    
    draw_jackpot_fish(draw, color_main, color_accent, "M", 2)
    
    # Major 標誌：橙色六角星
    cx, cy = 32, 32
    for i in range(6):
        angle = i * math.pi / 3
        sx = cx + int(5 * math.cos(angle))
        sy = cy + int(5 * math.sin(angle))
        draw.ellipse([sx-2, sy-2, sx+2, sy+2], fill=(255, 220, 50, 230))
    draw.ellipse([30, 30, 34, 34], fill=(255, 255, 200, 255))
    
    img.save(os.path.join(OUTPUT_DIR, "T173_jackpot_major.png"))
    print(f"T173 生成完成（Major Jackpot 橙色）")
    return img

def generate_t174_jackpot_grand():
    """T174 幸運 Grand Jackpot 魚 — 金色，12 道光芒，最大，三層光環"""
    img = Image.new("RGBA", SIZE, (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    color_main = (200, 160, 20)    # 深金
    color_accent = (255, 220, 50)  # 亮金
    
    # 三層外圈光環
    cx, cy = 32, 32
    for r, alpha in [(30, 40), (26, 60), (22, 80)]:
        draw.ellipse([cx-r, cy-r, cx+r, cy+r], 
                     outline=(255, 220, 50, alpha), width=2)
    
    draw_jackpot_fish(draw, color_main, color_accent, "G", 3)
    
    # Grand 標誌：皇冠形狀（用三角形模擬）
    crown_pts = [
        (cx-6, cy+3), (cx-6, cy-2), (cx-3, cy-5),
        (cx, cy-2), (cx+3, cy-5), (cx+6, cy-2),
        (cx+6, cy+3)
    ]
    draw.polygon(crown_pts, fill=(255, 240, 100, 230))
    
    img.save(os.path.join(OUTPUT_DIR, "T174_jackpot_grand.png"))
    print(f"T174 生成完成（Grand Jackpot 金色，三層光環）")
    return img

def generate_t175_jackpot_trigger():
    """T175 幸運 Jackpot Trigger 魚 — 彩虹色，四色分區，隨機觸發"""
    img = Image.new("RGBA", SIZE, (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    cx, cy = 32, 32
    
    # 大背景光暈
    draw.ellipse([cx-30, cy-30, cx+30, cy+30], fill=(200, 200, 100, 60))
    draw.ellipse([cx-26, cy-26, cx+26, cy+26], fill=(200, 200, 100, 80))
    
    # 四色光暈（代表四層）
    colors = [
        (30, 180, 30, 80),   # 綠（Mini）
        (30, 100, 220, 80),  # 藍（Minor）
        (220, 120, 20, 80),  # 橙（Major）
        (200, 160, 20, 80),  # 金（Grand）
    ]
    for i, c in enumerate(colors):
        angle = i * math.pi / 2
        ox = int(8 * math.cos(angle))
        oy = int(8 * math.sin(angle))
        draw.ellipse([cx+ox-14, cy+oy-14, cx+ox+14, cy+oy+14], fill=c)
    
    # 主魚身（白金色）
    draw.ellipse([cx-18, cy-12, cx+18, cy+12], fill=(220, 210, 180, 230))
    
    # 魚尾
    draw.polygon([(cx+14, cy), (cx+24, cy-7), (cx+24, cy+7)], 
                 fill=(200, 190, 150, 210))
    
    # 魚眼
    draw.ellipse([cx-14, cy-3, cx-8, cy+3], fill=(255, 255, 255, 255))
    draw.ellipse([cx-12, cy-2, cx-9, cy+1], fill=(0, 0, 0, 255))
    
    # 四色光芒（每方向一色）
    ray_colors = [(30, 180, 30), (30, 100, 220), (220, 120, 20), (200, 160, 20)]
    for i, rc in enumerate(ray_colors):
        for j in range(3):
            angle = i * math.pi / 2 + j * math.pi / 6 - math.pi / 6
            x1 = cx + int(10 * math.cos(angle))
            y1 = cy + int(10 * math.sin(angle))
            x2 = cx + int(24 * math.cos(angle))
            y2 = cy + int(24 * math.sin(angle))
            draw.line([x1, y1, x2, y2], fill=(*rc, 180), width=2)
    
    # 中心問號（用圓點代替）
    draw.ellipse([29, 28, 35, 34], fill=(255, 255, 255, 230))
    
    img.save(os.path.join(OUTPUT_DIR, "T175_jackpot_trigger.png"))
    print(f"T175 生成完成（Jackpot Trigger 彩虹色）")
    return img

def calculate_density(img):
    """計算精靈圖的非透明像素密度"""
    pixels = img.getdata()
    non_transparent = sum(1 for p in pixels if p[3] > 30)
    total = SIZE[0] * SIZE[1]
    return non_transparent / total * 100

if __name__ == "__main__":
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    
    print("=== DAY-313 Progressive Jackpot 系列精靈圖生成 ===")
    
    imgs = [
        generate_t171_jackpot_mini(),
        generate_t172_jackpot_minor(),
        generate_t173_jackpot_major(),
        generate_t174_jackpot_grand(),
        generate_t175_jackpot_trigger(),
    ]
    
    names = ["T171", "T172", "T173", "T174", "T175"]
    print("\n=== 密度報告 ===")
    for name, img in zip(names, imgs):
        density = calculate_density(img)
        status = "✅" if density >= 35 else "⚠️"
        print(f"{status} {name}: {density:.1f}%")
    
    print("\n✅ DAY-313 精靈圖生成完成！")
