#!/usr/bin/env python3
"""DAY-328 目標物精靈圖生成 — T229-T233"""
import os
import math
import random

try:
    from PIL import Image, ImageDraw
except ImportError:
    import subprocess
    subprocess.run(["pip", "install", "Pillow"], check=True)
    from PIL import Image, ImageDraw

OUTPUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
os.makedirs(OUTPUT_DIR, exist_ok=True)

TARGETS = [
    ("T229", "magnetic_attraction", (255, 102, 0),   "magnetic"),
    ("T230", "super_chain",         (0, 255, 255),   "chain"),
    ("T231", "holy_pillar",         (255, 255, 0),   "pillar"),
    ("T232", "time_stop",           (0, 200, 255),   "freeze"),
    ("T233", "cosmic_restart",      (255, 0, 255),   "restart"),
]

def draw_fish_body(draw, cx, cy, w, h, color, alpha=255):
    r, g, b = color
    draw.ellipse([cx-w//2, cy-h//2, cx+w//2, cy+h//2], fill=(r, g, b, alpha))
    # tail
    tail_pts = [(cx+w//2-4, cy), (cx+w//2+12, cy-10), (cx+w//2+12, cy+10)]
    draw.polygon(tail_pts, fill=(max(0,r-40), max(0,g-40), max(0,b-40), alpha))

def draw_rays(draw, cx, cy, n, length, color, alpha=200):
    r, g, b = color
    for i in range(n):
        angle = 2 * math.pi * i / n
        x2 = cx + int(length * math.cos(angle))
        y2 = cy + int(length * math.sin(angle))
        draw.line([(cx, cy), (x2, y2)], fill=(r, g, b, alpha), width=2)

def draw_ring(draw, cx, cy, radius, color, width=2, alpha=180):
    r, g, b = color
    draw.ellipse([cx-radius, cy-radius, cx+radius, cy+radius],
                 outline=(r, g, b, alpha), width=width)

def generate_target(tid, name, color, style):
    size = 64
    img = Image.new("RGBA", (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = size // 2, size // 2
    r, g, b = color

    if style == "magnetic":
        # 橙色魚身 + 磁力線 + 吸引符號
        draw_fish_body(draw, cx, cy, 36, 24, color)
        # 磁力弧線
        for i in range(4):
            angle_start = i * 90
            draw.arc([cx-20, cy-20, cx+20, cy+20], angle_start, angle_start+60,
                     fill=(255, 200, 0, 200), width=2)
        # 中心磁力點
        draw.ellipse([cx-4, cy-4, cx+4, cy+4], fill=(255, 255, 0, 255))
        draw_rays(draw, cx, cy, 8, 28, (255, 150, 0))

    elif style == "chain":
        # 青色魚身 + 連鎖符號
        draw_fish_body(draw, cx, cy, 36, 24, color)
        # 連鎖環
        for i in range(3):
            ox = cx - 8 + i * 8
            draw.ellipse([ox-5, cy-5, ox+5, cy+5], outline=(0, 255, 200, 220), width=2)
        draw_rays(draw, cx, cy, 12, 26, (0, 200, 255))
        draw_ring(draw, cx, cy, 20, (0, 255, 255), width=2)

    elif style == "pillar":
        # 金黃色魚身 + 光柱符號
        draw_fish_body(draw, cx, cy, 36, 24, color)
        # 12 道光柱指示
        for i in range(12):
            angle = 2 * math.pi * i / 12
            x1 = cx + int(16 * math.cos(angle))
            y1 = cy + int(16 * math.sin(angle))
            x2 = cx + int(28 * math.cos(angle))
            y2 = cy + int(28 * math.sin(angle))
            draw.line([(x1, y1), (x2, y2)], fill=(255, 255, 100, 200), width=2)
        draw_ring(draw, cx, cy, 14, (255, 255, 0), width=2)

    elif style == "freeze":
        # 冰藍色魚身 + 凍結符號
        draw_fish_body(draw, cx, cy, 36, 24, color)
        # 雪花符號
        for i in range(6):
            angle = math.pi * i / 3
            x2 = cx + int(18 * math.cos(angle))
            y2 = cy + int(18 * math.sin(angle))
            draw.line([(cx, cy), (x2, y2)], fill=(200, 240, 255, 220), width=2)
        # 凍結光環
        draw_ring(draw, cx, cy, 22, (0, 200, 255), width=2)
        draw_ring(draw, cx, cy, 28, (100, 220, 255), width=1)
        draw_rays(draw, cx, cy, 8, 30, (0, 180, 255))

    elif style == "restart":
        # 洋紅色大型魚身 + 重啟符號
        draw_fish_body(draw, cx, cy, 40, 28, color)
        # 重啟箭頭（圓形箭頭）
        draw.arc([cx-14, cy-14, cx+14, cy+14], 30, 330, fill=(255, 100, 255, 220), width=3)
        # 箭頭尖端
        draw.polygon([(cx+14, cy-2), (cx+10, cy-8), (cx+18, cy-8)], fill=(255, 0, 255, 220))
        # 多層光環
        for radius in [18, 24, 30]:
            draw_ring(draw, cx, cy, radius, (255, 0, 255), width=1)
        draw_rays(draw, cx, cy, 16, 30, (255, 0, 200))

    # 計算非透明像素比例
    pixels = list(img.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 10)
    ratio = non_transparent / len(pixels) * 100

    filename = f"{tid}_{name}.png"
    filepath = os.path.join(OUTPUT_DIR, filename)
    img.save(filepath)
    print(f"✅ {filename} — {ratio:.1f}% non-transparent")
    return ratio

if __name__ == "__main__":
    print("=== DAY-328 目標物精靈圖生成 ===")
    total_ratio = 0
    for tid, name, color, style in TARGETS:
        ratio = generate_target(tid, name, color, style)
        total_ratio += ratio
    avg = total_ratio / len(TARGETS)
    print(f"\n平均非透明像素比例：{avg:.1f}%")
    print("=== 生成完成 ===")
