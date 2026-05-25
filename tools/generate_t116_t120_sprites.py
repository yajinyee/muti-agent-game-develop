#!/usr/bin/env python3
"""
generate_t116_t120_sprites.py — DAY-295 T116-T120 精靈圖生成
T116 幸運千龍王輪盤魚 / T117 幸運龍力散彈魚 / T118 幸運火箭砲魚
T119 幸運深海漩渦魚 / T120 幸運吸血鬼魚
"""
import os
import math
from PIL import Image, ImageDraw

OUTPUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 32

def save_sprite(img, filename):
    path = os.path.join(OUTPUT_DIR, filename)
    img.save(path)
    # 計算非透明像素比例
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    total = len(pixels)
    print(f"  {filename}: {non_transparent/total*100:.1f}% 非透明像素")

def hex_to_rgba(hex_str, alpha=255):
    h = hex_str.lstrip('#')
    r, g, b = int(h[0:2], 16), int(h[2:4], 16), int(h[4:6], 16)
    return (r, g, b, alpha)

def draw_circle(draw, cx, cy, r, color):
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], fill=color)

def draw_ring(draw, cx, cy, r_outer, r_inner, color):
    """畫環形"""
    for y in range(cy - r_outer, cy + r_outer + 1):
        for x in range(cx - r_outer, cx + r_outer + 1):
            dx, dy = x - cx, y - cy
            dist = math.sqrt(dx*dx + dy*dy)
            if r_inner <= dist <= r_outer:
                draw.point((x, y), fill=color)

def draw_star(draw, cx, cy, r_outer, r_inner, points, color, rotation=0):
    """畫星形"""
    pts = []
    for i in range(points * 2):
        angle = math.pi * i / points + rotation
        r = r_outer if i % 2 == 0 else r_inner
        pts.append((cx + r * math.cos(angle), cy + r * math.sin(angle)))
    draw.polygon(pts, fill=color)

# ── T116 幸運千龍王輪盤魚 ─────────────────────────────────────
def gen_t116():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE // 2, SIZE // 2

    # 金色圓形魚身
    draw_circle(draw, cx, cy, 13, hex_to_rgba('#FFD700'))
    draw_circle(draw, cx, cy, 11, hex_to_rgba('#FFA500'))

    # 雙環輪盤紋路
    draw_ring(draw, cx, cy, 13, 11, hex_to_rgba('#FF8C00'))
    draw_ring(draw, cx, cy, 9, 7, hex_to_rgba('#FF6B00'))

    # 中心千龍王符文（金色十字）
    for dx in range(-2, 3):
        draw.point((cx + dx, cy), fill=hex_to_rgba('#FFFFFF'))
    for dy in range(-2, 3):
        draw.point((cx, cy + dy), fill=hex_to_rgba('#FFFFFF'))

    # 4 方向金色光芒
    for angle_deg in [0, 90, 180, 270]:
        angle = math.radians(angle_deg)
        for r in range(14, 16):
            x = int(cx + r * math.cos(angle))
            y = int(cy + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=hex_to_rgba('#FFD700'))

    # 金色眼睛
    draw.point((cx + 4, cy - 2), fill=hex_to_rgba('#FFFFFF'))

    # 金色魚尾
    for i in range(5):
        draw.point((cx - 12 + i, cy - 2 + i), fill=hex_to_rgba('#FF8C00'))
        draw.point((cx - 12 + i, cy + 2 - i), fill=hex_to_rgba('#FF8C00'))

    save_sprite(img, "T116_chain_long_king.png")

# ── T117 幸運龍力散彈魚 ───────────────────────────────────────
def gen_t117():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE // 2, SIZE // 2

    # 紫色橢圓魚身
    draw.ellipse([cx-12, cy-8, cx+12, cy+8], fill=hex_to_rgba('#7B2FBE'))
    draw.ellipse([cx-10, cy-6, cx+10, cy+6], fill=hex_to_rgba('#9B4FDE'))

    # 8 方向散彈光芒
    for i in range(8):
        angle = math.radians(i * 45)
        for r in range(12, 15):
            x = int(cx + r * math.cos(angle))
            y = int(cy + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=hex_to_rgba('#E040FB'))

    # 龍鱗紋路
    for i in range(-3, 4, 2):
        draw.point((cx + i, cy - 3), fill=hex_to_rgba('#C77DFF'))
        draw.point((cx + i, cy + 3), fill=hex_to_rgba('#C77DFF'))

    # 紫色眼睛
    draw.point((cx + 5, cy - 2), fill=hex_to_rgba('#FFFFFF'))

    # 紫色魚尾
    for i in range(4):
        draw.point((cx - 11 + i, cy - 3 + i), fill=hex_to_rgba('#6A1B9A'))
        draw.point((cx - 11 + i, cy + 3 - i), fill=hex_to_rgba('#6A1B9A'))

    # 中心散彈核心
    draw_circle(draw, cx, cy, 3, hex_to_rgba('#FFD700'))

    save_sprite(img, "T117_dragon_shotgun.png")

# ── T118 幸運火箭砲魚 ─────────────────────────────────────────
def gen_t118():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE // 2, SIZE // 2

    # 火紅橢圓魚身
    draw.ellipse([cx-12, cy-7, cx+12, cy+7], fill=hex_to_rgba('#FF3D00'))
    draw.ellipse([cx-10, cy-5, cx+10, cy+5], fill=hex_to_rgba('#FF6D00'))

    # 火箭砲紋路（橫向條紋）
    for i in range(-4, 5, 2):
        draw.line([(cx - 8, cy + i), (cx + 8, cy + i)], fill=hex_to_rgba('#FF8A65'), width=1)

    # 4 方向火焰光芒
    for angle_deg in [0, 90, 180, 270]:
        angle = math.radians(angle_deg)
        for r in range(13, 16):
            x = int(cx + r * math.cos(angle))
            y = int(cy + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=hex_to_rgba('#FF6D00'))

    # 火箭頭（右側尖端）
    for i in range(3):
        draw.point((cx + 12 - i, cy - i), fill=hex_to_rgba('#FFD700'))
        draw.point((cx + 12 - i, cy + i), fill=hex_to_rgba('#FFD700'))

    # 火焰尾（左側）
    for i in range(4):
        draw.point((cx - 12 + i, cy - 2 + i), fill=hex_to_rgba('#FF1744'))
        draw.point((cx - 12 + i, cy + 2 - i), fill=hex_to_rgba('#FF1744'))

    # 中心爆炸核心
    draw_circle(draw, cx, cy, 3, hex_to_rgba('#FFD700'))

    # 眼睛
    draw.point((cx + 4, cy - 2), fill=hex_to_rgba('#FFFFFF'))

    save_sprite(img, "T118_rocket_cannon.png")

# ── T119 幸運深海漩渦魚 ───────────────────────────────────────
def gen_t119():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE // 2, SIZE // 2

    # 深藍圓形魚身
    draw_circle(draw, cx, cy, 13, hex_to_rgba('#0D47A1'))
    draw_circle(draw, cx, cy, 11, hex_to_rgba('#1565C0'))

    # 漩渦紋路（螺旋線）
    for i in range(36):
        angle = math.radians(i * 10)
        r = 2 + i * 0.25
        if r > 10:
            break
        x = int(cx + r * math.cos(angle))
        y = int(cy + r * math.sin(angle))
        if 0 <= x < SIZE and 0 <= y < SIZE:
            draw.point((x, y), fill=hex_to_rgba('#00B4DC'))

    # 深海光環
    draw_ring(draw, cx, cy, 13, 11, hex_to_rgba('#00E5FF', 180))

    # 4 方向水流光芒
    for angle_deg in [45, 135, 225, 315]:
        angle = math.radians(angle_deg)
        for r in range(13, 16):
            x = int(cx + r * math.cos(angle))
            y = int(cy + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=hex_to_rgba('#00B4DC'))

    # 中心漩渦核心
    draw_circle(draw, cx, cy, 3, hex_to_rgba('#00E5FF'))

    # 眼睛
    draw.point((cx + 5, cy - 2), fill=hex_to_rgba('#FFFFFF'))

    # 魚尾
    for i in range(4):
        draw.point((cx - 12 + i, cy - 3 + i), fill=hex_to_rgba('#0D47A1'))
        draw.point((cx - 12 + i, cy + 3 - i), fill=hex_to_rgba('#0D47A1'))

    save_sprite(img, "T119_deep_whirlpool.png")

# ── T120 幸運吸血鬼魚 ─────────────────────────────────────────
def gen_t120():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE // 2, SIZE // 2

    # 深紫橢圓魚身
    draw.ellipse([cx-12, cy-8, cx+12, cy+8], fill=hex_to_rgba('#4A0072'))
    draw.ellipse([cx-10, cy-6, cx+10, cy+6], fill=hex_to_rgba('#6A1B9A'))

    # 吸血鬼牙齒（下方兩顆白牙）
    draw.point((cx - 2, cy + 5), fill=hex_to_rgba('#FFFFFF'))
    draw.point((cx + 2, cy + 5), fill=hex_to_rgba('#FFFFFF'))
    draw.point((cx - 2, cy + 6), fill=hex_to_rgba('#FFFFFF'))
    draw.point((cx + 2, cy + 6), fill=hex_to_rgba('#FFFFFF'))

    # 吸血鬼眼睛（紅色）
    draw.point((cx + 4, cy - 2), fill=hex_to_rgba('#FF1744'))
    draw.point((cx + 5, cy - 2), fill=hex_to_rgba('#FF1744'))

    # 倍率吸收光環
    draw_ring(draw, cx, cy, 13, 11, hex_to_rgba('#CE93D8', 160))

    # 4 方向紫色光芒
    for angle_deg in [0, 90, 180, 270]:
        angle = math.radians(angle_deg)
        for r in range(13, 16):
            x = int(cx + r * math.cos(angle))
            y = int(cy + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=hex_to_rgba('#AB47BC'))

    # 吸血鬼翅膀（上方兩側）
    for i in range(4):
        draw.point((cx - 8 + i, cy - 8 + i), fill=hex_to_rgba('#7B1FA2'))
        draw.point((cx + 8 - i, cy - 8 + i), fill=hex_to_rgba('#7B1FA2'))

    # 中心倍率核心
    draw_circle(draw, cx, cy, 3, hex_to_rgba('#E040FB'))

    # 魚尾
    for i in range(4):
        draw.point((cx - 11 + i, cy - 3 + i), fill=hex_to_rgba('#4A0072'))
        draw.point((cx - 11 + i, cy + 3 - i), fill=hex_to_rgba('#4A0072'))

    save_sprite(img, "T120_vampire_mult.png")

if __name__ == "__main__":
    print("=== DAY-295 T116-T120 精靈圖生成 ===")
    gen_t116()
    gen_t117()
    gen_t118()
    gen_t119()
    gen_t120()
    print("✅ 全部完成！")
