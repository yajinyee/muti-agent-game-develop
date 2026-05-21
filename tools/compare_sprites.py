# -*- coding: utf-8 -*-
"""比較 v5 和 v6 角色品質，輸出對比圖"""
from PIL import Image, ImageDraw, ImageFont
import os

CHARS_DIR = "D:/Kiro/client/chiikawa-pixel/assets/sprites/characters"
OUT_PATH  = "D:/Kiro/tools/sprite_comparison.png"

chars  = ["chiikawa", "hachiware", "usagi"]
states = ["idle", "attack", "bigwin"]
CELL   = 96
SCALE  = 2  # 放大顯示

# 建立對比圖（3角色 × 3狀態，每格放大 2x）
cols = len(states)
rows = len(chars)
W = cols * CELL * SCALE + 20
H = rows * CELL * SCALE + 20

canvas = Image.new("RGBA", (W, H), (40, 40, 40, 255))

for row, char in enumerate(chars):
    for col, state in enumerate(states):
        path = os.path.join(CHARS_DIR, f"{char}_{state}.png")
        if not os.path.exists(path):
            continue
        img = Image.open(path).convert("RGBA")
        # 放大 2x
        img = img.resize((CELL * SCALE, CELL * SCALE), Image.NEAREST)
        # 貼到畫布
        x = col * CELL * SCALE + 10
        y = row * CELL * SCALE + 10
        # 白色背景
        bg = Image.new("RGBA", (CELL * SCALE, CELL * SCALE), (60, 60, 60, 255))
        canvas.paste(bg, (x, y))
        canvas.paste(img, (x, y), img)

canvas.save(OUT_PATH)
print(f"對比圖儲存: {OUT_PATH}")
print(f"尺寸: {canvas.width}x{canvas.height}")

# 統計各角色的像素密度
print("\n像素密度分析：")
for char in chars:
    for state in states:
        path = os.path.join(CHARS_DIR, f"{char}_{state}.png")
        if not os.path.exists(path):
            continue
        img = Image.open(path).convert("RGBA")
        non_t = sum(1 for px in img.getdata() if px[3] > 10)
        bbox = img.getbbox()
        if bbox:
            w = bbox[2]-bbox[0]
            h = bbox[3]-bbox[1]
            density = non_t / (w*h) * 100 if w*h > 0 else 0
            print(f"  {char}_{state}: {w}x{h}px, {non_t}px, 密度={density:.0f}%")
