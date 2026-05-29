#!/usr/bin/env python3
"""
generate_targets_day329.py — DAY-329 T234-T238 精靈圖生成
T234 幸運FeverBoost升級魚：橙紅色魚身 + 火焰符號 + 特殊目標光環
T235 幸運快速暴富升級魚：金黃色魚身 + 閃電符號 + 速度線
T236 幸運冰釣大師魚：冰藍色魚身 + 雪花符號 + 釣竿圖示
T237 幸運宇宙奇蹟魚：深紫色大型魚身 + 8道光柱指示 + 多層光環
T238 幸運創世終極魚：純金色超大型魚身 + 12道光柱 + 里程碑光環
"""
import os
import math
from PIL import Image, ImageDraw

OUTPUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
os.makedirs(OUTPUT_DIR, exist_ok=True)

def draw_fish_body(draw, cx, cy, w, h, color, outline_color):
    """繪製基本魚身（橢圓）"""
    draw.ellipse([cx - w//2, cy - h//2, cx + w//2, cy + h//2],
                 fill=color, outline=outline_color, width=2)

def draw_tail(draw, cx, cy, size, color):
    """繪製魚尾"""
    points = [
        (cx + size//2, cy),
        (cx + size, cy - size//3),
        (cx + size, cy + size//3),
    ]
    draw.polygon(points, fill=color)

def draw_eye(draw, cx, cy, r, color=(255, 255, 255)):
    """繪製魚眼"""
    draw.ellipse([cx - r, cy - r, cx + r, cy + r], fill=color, outline=(0, 0, 0), width=1)
    draw.ellipse([cx - r//2, cy - r//2, cx + r//2, cy + r//2], fill=(0, 0, 0))

def draw_rays(draw, cx, cy, count, length, color, width=2):
    """繪製放射光芒"""
    for i in range(count):
        angle = (2 * math.pi * i) / count
        x2 = cx + int(length * math.cos(angle))
        y2 = cy + int(length * math.sin(angle))
        draw.line([(cx, cy), (x2, y2)], fill=color, width=width)

def draw_ring(draw, cx, cy, r, color, width=2):
    """繪製圓環"""
    draw.ellipse([cx - r, cy - r, cx + r, cy + r], outline=color, width=width)

def generate_t234():
    """T234 幸運FeverBoost升級魚：橙紅色魚身 + 火焰符號 + 特殊目標光環"""
    img = Image.new("RGBA", (64, 64), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32
    body_color = (255, 100, 0, 230)
    outline_color = (200, 50, 0, 255)
    draw_fish_body(draw, cx - 4, cy, 36, 22, body_color, outline_color)
    draw_tail(draw, cx - 4, cy, 16, (200, 60, 0, 200))
    draw_eye(draw, cx + 8, cy - 4, 4)
    # 火焰符號（三角形）
    flame_pts = [(cx, cy - 18), (cx - 8, cy - 6), (cx + 8, cy - 6)]
    draw.polygon(flame_pts, fill=(255, 200, 0, 200))
    flame_pts2 = [(cx, cy - 14), (cx - 5, cy - 6), (cx + 5, cy - 6)]
    draw.polygon(flame_pts2, fill=(255, 120, 0, 220))
    # 特殊目標光環
    draw_ring(draw, cx, cy, 28, (255, 150, 0, 180), 2)
    draw_ring(draw, cx, cy, 24, (255, 80, 0, 120), 1)
    # 光芒
    draw_rays(draw, cx, cy, 8, 30, (255, 180, 0, 100), 1)
    img.save(os.path.join(OUTPUT_DIR, "T234_fever_boost_ultimate.png"))
    print(f"T234 generated")

def generate_t235():
    """T235 幸運快速暴富升級魚：金黃色魚身 + 閃電符號 + 速度線"""
    img = Image.new("RGBA", (64, 64), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32
    body_color = (255, 215, 0, 230)
    outline_color = (200, 160, 0, 255)
    draw_fish_body(draw, cx - 4, cy, 36, 22, body_color, outline_color)
    draw_tail(draw, cx - 4, cy, 16, (200, 160, 0, 200))
    draw_eye(draw, cx + 8, cy - 4, 4)
    # 閃電符號
    lightning = [(cx + 2, cy - 16), (cx - 4, cy - 2), (cx + 2, cy - 2), (cx - 2, cy + 12), (cx + 6, cy - 4), (cx + 0, cy - 4)]
    draw.polygon(lightning, fill=(255, 255, 0, 220))
    # 速度線
    for i in range(3):
        y_off = cy - 6 + i * 6
        draw.line([(cx - 28, y_off), (cx - 16, y_off)], fill=(255, 200, 0, 150), width=2)
    # 光環
    draw_ring(draw, cx, cy, 28, (255, 220, 0, 180), 2)
    draw_rays(draw, cx, cy, 12, 30, (255, 200, 0, 100), 1)
    img.save(os.path.join(OUTPUT_DIR, "T235_rapid_riches_ultimate.png"))
    print(f"T235 generated")

def generate_t236():
    """T236 幸運冰釣大師魚：冰藍色魚身 + 雪花符號 + 釣竿圖示"""
    img = Image.new("RGBA", (64, 64), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32
    body_color = (0, 190, 255, 230)
    outline_color = (0, 120, 200, 255)
    draw_fish_body(draw, cx - 4, cy, 36, 22, body_color, outline_color)
    draw_tail(draw, cx - 4, cy, 16, (0, 140, 200, 200))
    draw_eye(draw, cx + 8, cy - 4, 4)
    # 雪花符號（6方向）
    draw_rays(draw, cx, cy - 16, 6, 8, (200, 240, 255, 220), 2)
    # 雪花小枝
    for i in range(6):
        angle = (2 * math.pi * i) / 6
        bx = cx + int(6 * math.cos(angle))
        by = (cy - 16) + int(6 * math.sin(angle))
        perp = angle + math.pi / 2
        draw.line([(bx - int(3 * math.cos(perp)), by - int(3 * math.sin(perp))),
                   (bx + int(3 * math.cos(perp)), by + int(3 * math.sin(perp)))],
                  fill=(200, 240, 255, 180), width=1)
    # 冰晶光環
    draw_ring(draw, cx, cy, 28, (0, 200, 255, 180), 2)
    draw_ring(draw, cx, cy, 22, (100, 220, 255, 120), 1)
    draw_rays(draw, cx, cy, 8, 30, (0, 180, 255, 80), 1)
    img.save(os.path.join(OUTPUT_DIR, "T236_ice_fishing_master.png"))
    print(f"T236 generated")

def generate_t237():
    """T237 幸運宇宙奇蹟魚：深紫色大型魚身 + 8道光柱指示 + 多層光環"""
    img = Image.new("RGBA", (64, 64), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32
    body_color = (148, 0, 211, 230)
    outline_color = (100, 0, 160, 255)
    draw_fish_body(draw, cx - 4, cy, 40, 26, body_color, outline_color)
    draw_tail(draw, cx - 4, cy, 18, (100, 0, 160, 200))
    draw_eye(draw, cx + 10, cy - 5, 5)
    # 8道光柱指示（小點）
    for i in range(8):
        angle = (2 * math.pi * i) / 8
        px = cx + int(20 * math.cos(angle))
        py = cy + int(20 * math.sin(angle))
        draw.ellipse([px - 3, py - 3, px + 3, py + 3], fill=(200, 100, 255, 200))
    # 多層光環
    draw_ring(draw, cx, cy, 30, (180, 0, 255, 180), 2)
    draw_ring(draw, cx, cy, 24, (140, 0, 200, 140), 2)
    draw_ring(draw, cx, cy, 18, (100, 0, 160, 100), 1)
    draw_rays(draw, cx, cy, 8, 32, (160, 0, 220, 100), 1)
    img.save(os.path.join(OUTPUT_DIR, "T237_cosmic_miracle.png"))
    print(f"T237 generated")

def generate_t238():
    """T238 幸運創世終極魚：純金色超大型魚身 + 12道光柱 + 里程碑光環"""
    img = Image.new("RGBA", (64, 64), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32
    body_color = (255, 215, 0, 240)
    outline_color = (180, 140, 0, 255)
    draw_fish_body(draw, cx - 4, cy, 44, 28, body_color, outline_color)
    draw_tail(draw, cx - 4, cy, 20, (180, 140, 0, 210))
    draw_eye(draw, cx + 12, cy - 6, 5)
    # 12道光柱指示（小點）
    for i in range(12):
        angle = (2 * math.pi * i) / 12
        px = cx + int(22 * math.cos(angle))
        py = cy + int(22 * math.sin(angle))
        draw.ellipse([px - 2, py - 2, px + 2, py + 2], fill=(255, 255, 100, 220))
    # 里程碑光環（三層）
    draw_ring(draw, cx, cy, 30, (255, 220, 0, 200), 3)
    draw_ring(draw, cx, cy, 25, (255, 180, 0, 160), 2)
    draw_ring(draw, cx, cy, 20, (255, 140, 0, 120), 1)
    # 16道光芒（里程碑特效）
    draw_rays(draw, cx, cy, 16, 32, (255, 200, 0, 120), 1)
    img.save(os.path.join(OUTPUT_DIR, "T238_genesis_ultimate.png"))
    print(f"T238 generated")

if __name__ == "__main__":
    generate_t234()
    generate_t235()
    generate_t236()
    generate_t237()
    generate_t238()
    print("DAY-329 T234-T238 sprites generated!")
