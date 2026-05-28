#!/usr/bin/env python3
"""
generate_targets_day317.py — DAY-317 T191-T195 精靈圖生成器
生成 32x32 像素精靈圖，使用 Python Pillow

目標物：
- T191 PvP 競技魚：紅藍對抗色系，劍盾圖案
- T192 技能連鎖魚：彩虹漸層，連鎖符文
- T193 全服大爆炸魚：深紅爆炸，多層光環
- T194 時空折疊魚：藍紫折疊紋路，時空裂縫
- T195 宇宙終焉魚：黑色+金色，終焉符文，最大最亮
"""

import os
import math
from PIL import Image, ImageDraw

OUTPUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 32

def ensure_dir():
    os.makedirs(OUTPUT_DIR, exist_ok=True)

def save(img: Image.Image, filename: str):
    path = os.path.join(OUTPUT_DIR, filename)
    img.save(path)
    print(f"  ✓ 已生成: {filename}")

# ── 輔助函數 ──────────────────────────────────────────────────

def new_img():
    return Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

def draw_circle(draw, cx, cy, r, color, fill=True):
    x0, y0 = cx - r, cy - r
    x1, y1 = cx + r, cy + r
    if fill:
        draw.ellipse([x0, y0, x1, y1], fill=color)
    else:
        draw.ellipse([x0, y0, x1, y1], outline=color, width=1)

def draw_ring(draw, cx, cy, r_outer, r_inner, color):
    """繪製環形"""
    for y in range(SIZE):
        for x in range(SIZE):
            dx, dy = x - cx, y - cy
            dist = math.sqrt(dx*dx + dy*dy)
            if r_inner <= dist <= r_outer:
                draw.point((x, y), fill=color)

def lerp_color(c1, c2, t):
    """線性插值顏色"""
    return tuple(int(c1[i] + (c2[i] - c1[i]) * t) for i in range(len(c1)))

# ── T191 PvP 競技魚：紅藍對抗色系，劍盾圖案 ──────────────────

def generate_t191():
    img = new_img()
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE // 2, SIZE // 2

    # 背景：左紅右藍漸層
    for x in range(SIZE):
        t = x / (SIZE - 1)
        r = int(200 * (1 - t) + 30 * t)
        g = int(20 * (1 - t) + 50 * t)
        b = int(20 * (1 - t) + 200 * t)
        for y in range(SIZE):
            dist = math.sqrt((x - cx)**2 + (y - cy)**2)
            if dist <= 14:
                alpha = 255 if dist <= 12 else int(255 * (14 - dist) / 2)
                img.putpixel((x, y), (r, g, b, alpha))

    # 中央分隔線（白色）
    for y in range(4, SIZE - 4):
        draw.point((cx, y), fill=(255, 255, 255, 200))

    # 左側：劍（紅色）
    # 劍身
    for y in range(6, 22):
        draw.point((cx - 4, y), fill=(255, 80, 80, 255))
        draw.point((cx - 5, y), fill=(255, 120, 120, 255))
    # 劍柄
    for x in range(cx - 8, cx - 1):
        draw.point((x, 20), fill=(200, 150, 50, 255))
    # 劍尖
    draw.point((cx - 4, 5), fill=(255, 255, 255, 255))
    draw.point((cx - 5, 6), fill=(255, 200, 200, 255))

    # 右側：盾（藍色）
    # 盾形
    for y in range(7, 23):
        w = min(4, 4 - abs(y - 15) // 3)
        for x in range(cx + 2, cx + 2 + w + 1):
            draw.point((x, y), fill=(80, 150, 255, 255))
    # 盾邊框
    draw.rectangle([cx + 2, 7, cx + 6, 22], outline=(200, 220, 255, 255), width=1)

    # 外圈光環（金色）
    draw_ring(draw, cx, cy, 15, 13, (255, 215, 0, 180))

    save(img, "T191_pvp_battle.png")

# ── T192 技能連鎖魚：彩虹漸層，連鎖符文 ──────────────────────

def generate_t192():
    img = new_img()
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE // 2, SIZE // 2

    # 背景：彩虹漸層圓形
    for y in range(SIZE):
        for x in range(SIZE):
            dx, dy = x - cx, y - cy
            dist = math.sqrt(dx*dx + dy*dy)
            if dist <= 14:
                angle = math.atan2(dy, dx)
                hue = (angle + math.pi) / (2 * math.pi)  # 0-1
                # 彩虹色
                h6 = hue * 6
                i = int(h6)
                f = h6 - i
                if i == 0:   r, g, b = 255, int(255*f), 0
                elif i == 1: r, g, b = int(255*(1-f)), 255, 0
                elif i == 2: r, g, b = 0, 255, int(255*f)
                elif i == 3: r, g, b = 0, int(255*(1-f)), 255
                elif i == 4: r, g, b = int(255*f), 0, 255
                else:        r, g, b = 255, 0, int(255*(1-f))
                alpha = 255 if dist <= 12 else int(255 * (14 - dist) / 2)
                img.putpixel((x, y), (r, g, b, alpha))

    # 連鎖符文：中央六芒星
    # 六個點
    for i in range(6):
        angle = i * math.pi / 3
        px = int(cx + 6 * math.cos(angle))
        py = int(cy + 6 * math.sin(angle))
        draw.ellipse([px-1, py-1, px+1, py+1], fill=(255, 255, 255, 255))

    # 連接線（白色）
    for i in range(6):
        a1 = i * math.pi / 3
        a2 = (i + 2) % 6 * math.pi / 3
        x1 = int(cx + 6 * math.cos(a1))
        y1 = int(cy + 6 * math.sin(a1))
        x2 = int(cx + 6 * math.cos(a2))
        y2 = int(cy + 6 * math.sin(a2))
        draw.line([x1, y1, x2, y2], fill=(255, 255, 255, 200), width=1)

    # 中心點（白色）
    draw.ellipse([cx-2, cy-2, cx+2, cy+2], fill=(255, 255, 255, 255))

    # 外圈（白色閃光）
    draw_ring(draw, cx, cy, 15, 13, (255, 255, 255, 150))

    save(img, "T192_skill_chain.png")

# ── T193 全服大爆炸魚：深紅爆炸，多層光環 ────────────────────

def generate_t193():
    img = new_img()
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE // 2, SIZE // 2

    # 背景：深紅爆炸放射狀
    for y in range(SIZE):
        for x in range(SIZE):
            dx, dy = x - cx, y - cy
            dist = math.sqrt(dx*dx + dy*dy)
            if dist <= 15:
                # 放射狀亮度
                angle = math.atan2(dy, dx)
                ray = abs(math.sin(angle * 8)) * 0.5 + 0.5
                intensity = (1 - dist / 15) * ray
                r = int(220 * intensity + 100)
                g = int(20 * intensity)
                b = int(10 * intensity)
                r = min(255, r)
                alpha = 255 if dist <= 13 else int(255 * (15 - dist) / 2)
                img.putpixel((x, y), (r, g, b, alpha))

    # 多層光環
    draw_ring(draw, cx, cy, 13, 11, (255, 100, 0, 200))
    draw_ring(draw, cx, cy, 10, 8, (255, 150, 0, 180))
    draw_ring(draw, cx, cy, 7, 5, (255, 200, 50, 160))

    # 中心爆炸核心（白色）
    draw.ellipse([cx-3, cy-3, cx+3, cy+3], fill=(255, 255, 255, 255))

    # 爆炸射線（8方向）
    for i in range(8):
        angle = i * math.pi / 4
        for r in range(4, 14):
            px = int(cx + r * math.cos(angle))
            py = int(cy + r * math.sin(angle))
            if 0 <= px < SIZE and 0 <= py < SIZE:
                alpha = int(255 * (1 - r / 14))
                draw.point((px, py), fill=(255, 200, 0, alpha))

    save(img, "T193_global_explosion.png")

# ── T194 時空折疊魚：藍紫折疊紋路，時空裂縫 ─────────────────

def generate_t194():
    img = new_img()
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE // 2, SIZE // 2

    # 背景：藍紫漸層
    for y in range(SIZE):
        for x in range(SIZE):
            dx, dy = x - cx, y - cy
            dist = math.sqrt(dx*dx + dy*dy)
            if dist <= 14:
                t = dist / 14
                r = int(80 * (1 - t) + 20 * t)
                g = int(20 * (1 - t) + 10 * t)
                b = int(200 * (1 - t) + 150 * t)
                alpha = 255 if dist <= 12 else int(255 * (14 - dist) / 2)
                img.putpixel((x, y), (r, g, b, alpha))

    # 折疊紋路（螺旋線）
    for i in range(100):
        t = i / 100
        angle = t * 4 * math.pi
        r = t * 10
        px = int(cx + r * math.cos(angle))
        py = int(cy + r * math.sin(angle))
        if 0 <= px < SIZE and 0 <= py < SIZE:
            draw.point((px, py), fill=(180, 100, 255, 220))

    # 時空裂縫（對角線裂縫）
    # 裂縫1
    for i in range(12):
        x = cx - 6 + i
        y = cy - 3 + i // 2
        if 0 <= x < SIZE and 0 <= y < SIZE:
            draw.point((x, y), fill=(255, 255, 255, 255))
            if y + 1 < SIZE:
                draw.point((x, y + 1), fill=(200, 150, 255, 180))

    # 裂縫2（反向）
    for i in range(10):
        x = cx + 2 - i
        y = cy + 2 + i // 2
        if 0 <= x < SIZE and 0 <= y < SIZE:
            draw.point((x, y), fill=(255, 255, 255, 220))

    # 外圈（紫色光環）
    draw_ring(draw, cx, cy, 15, 13, (150, 50, 255, 180))

    save(img, "T194_spacetime_fold.png")

# ── T195 宇宙終焉魚：黑色+金色，終焉符文，最大最亮 ──────────

def generate_t195():
    img = new_img()
    draw = ImageDraw.Draw(img)
    cx, cy = SIZE // 2, SIZE // 2

    # 背景：深黑色圓形
    for y in range(SIZE):
        for x in range(SIZE):
            dx, dy = x - cx, y - cy
            dist = math.sqrt(dx*dx + dy*dy)
            if dist <= 15:
                # 深黑到深紫漸層
                t = dist / 15
                r = int(10 * (1 - t) + 30 * t)
                g = int(5 * (1 - t) + 5 * t)
                b = int(20 * (1 - t) + 40 * t)
                alpha = 255 if dist <= 13 else int(255 * (15 - dist) / 2)
                img.putpixel((x, y), (r, g, b, alpha))

    # 終焉符文：金色八芒星
    for i in range(8):
        angle = i * math.pi / 4
        # 外點
        px_out = int(cx + 11 * math.cos(angle))
        py_out = int(cy + 11 * math.sin(angle))
        # 內點
        px_in = int(cx + 5 * math.cos(angle + math.pi / 8))
        py_in = int(cy + 5 * math.sin(angle + math.pi / 8))
        # 連線
        draw.line([cx, cy, px_out, py_out], fill=(255, 215, 0, 255), width=1)
        draw.ellipse([px_out-1, py_out-1, px_out+1, py_out+1], fill=(255, 255, 200, 255))

    # 中心：金色核心
    draw.ellipse([cx-4, cy-4, cx+4, cy+4], fill=(255, 215, 0, 255))
    draw.ellipse([cx-2, cy-2, cx+2, cy+2], fill=(255, 255, 255, 255))

    # 多層金色光環（最亮）
    draw_ring(draw, cx, cy, 15, 14, (255, 215, 0, 255))
    draw_ring(draw, cx, cy, 13, 12, (255, 200, 0, 220))
    draw_ring(draw, cx, cy, 11, 10, (255, 180, 0, 180))

    # 金色射線（12方向，最密集）
    for i in range(12):
        angle = i * math.pi / 6
        for r in range(5, 14):
            px = int(cx + r * math.cos(angle))
            py = int(cy + r * math.sin(angle))
            if 0 <= px < SIZE and 0 <= py < SIZE:
                alpha = int(255 * (1 - (r - 5) / 9))
                draw.point((px, py), fill=(255, 215, 0, alpha))

    save(img, "T195_cosmic_end.png")

# ── 主程式 ────────────────────────────────────────────────────

def main():
    ensure_dir()
    print("=== DAY-317 T191-T195 精靈圖生成 ===")
    print(f"輸出目錄: {OUTPUT_DIR}")
    print()

    print("生成 T191 PvP 競技魚（紅藍對抗色系，劍盾圖案）...")
    generate_t191()

    print("生成 T192 技能連鎖魚（彩虹漸層，連鎖符文）...")
    generate_t192()

    print("生成 T193 全服大爆炸魚（深紅爆炸，多層光環）...")
    generate_t193()

    print("生成 T194 時空折疊魚（藍紫折疊紋路，時空裂縫）...")
    generate_t194()

    print("生成 T195 宇宙終焉魚（黑色+金色，終焉符文，最大最亮）...")
    generate_t195()

    print()
    print("=== 全部完成！===")
    print("已生成 5 個精靈圖：T191-T195")

if __name__ == "__main__":
    main()
