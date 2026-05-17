# -*- coding: utf-8 -*-
"""
預覽所有角色 Sprite，生成一張對比圖
讓我們能直接看到目前的美術品質
"""
from PIL import Image, ImageDraw
import os

CHARS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"
TARGETS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
OUT_PATH = r"D:\Kiro\docs\sprite_preview.png"

def make_preview():
    chars = ["chiikawa", "hachiware", "usagi"]
    states = ["idle", "attack", "bigwin"]
    
    cell = 80  # 每格大小
    padding = 10
    label_h = 20
    
    cols = len(states) + 1  # +1 for label
    rows = len(chars) + 1   # +1 for header
    
    w = cols * (cell + padding) + padding
    h = rows * (cell + padding + label_h) + padding
    
    preview = Image.new("RGB", (w, h), (30, 30, 50))
    draw = ImageDraw.Draw(preview)
    
    # Header labels
    for ci, state in enumerate(states):
        x = padding + (ci + 1) * (cell + padding) + cell // 2
        draw.text((x - 20, padding), state.upper(), fill=(200, 200, 200))
    
    # Character rows
    for ri, char in enumerate(chars):
        y_base = padding + (ri + 1) * (cell + padding + label_h)
        
        # Row label
        draw.text((padding + 5, y_base + cell // 2), char, fill=(200, 200, 200))
        
        for ci, state in enumerate(states):
            x = padding + (ci + 1) * (cell + padding)
            y = y_base
            
            path = os.path.join(CHARS_DIR, f"{char}_{state}.png")
            if os.path.exists(path):
                sprite = Image.open(path).convert("RGBA")
                # 縮放到 cell 大小
                sprite = sprite.resize((cell, cell), Image.NEAREST)
                # 貼上（處理透明度）
                bg = Image.new("RGB", (cell, cell), (50, 50, 80))
                bg.paste(sprite, (0, 0), sprite)
                preview.paste(bg, (x, y))
                
                # 邊框
                draw.rectangle([x-1, y-1, x+cell, y+cell], outline=(100, 100, 150))
            else:
                draw.rectangle([x, y, x+cell, y+cell], fill=(80, 30, 30))
                draw.text((x+5, y+cell//2), "MISSING", fill=(255, 100, 100))
    
    # 目標物預覽（下方）
    target_files = sorted([f for f in os.listdir(TARGETS_DIR) if f.endswith(".png")])
    t_y = h - padding - cell - label_h
    
    for ti, fname in enumerate(target_files[:8]):
        x = padding + ti * (cell//2 + padding//2)
        path = os.path.join(TARGETS_DIR, fname)
        sprite = Image.open(path).convert("RGBA")
        sprite = sprite.resize((cell//2, cell//2), Image.NEAREST)
        bg = Image.new("RGB", (cell//2, cell//2), (50, 50, 80))
        bg.paste(sprite, (0, 0), sprite)
        preview.paste(bg, (x, t_y))
        draw.text((x, t_y + cell//2 + 2), fname[:6], fill=(150, 150, 150))
    
    preview.save(OUT_PATH)
    print(f"Preview saved: {OUT_PATH}")
    print(f"Size: {preview.width}x{preview.height}")
    
    # 也輸出每個 sprite 的顏色統計
    print("\n=== Color Analysis ===")
    for char in chars:
        path = os.path.join(CHARS_DIR, f"{char}_idle.png")
        if os.path.exists(path):
            img = Image.open(path).convert("RGBA")
            # 非透明像素的主要顏色
            pixels = [img.getpixel((x, y)) for y in range(img.height) 
                     for x in range(img.width) if img.getpixel((x, y))[3] > 50]
            if pixels:
                from collections import Counter
                top_colors = Counter([(r//20*20, g//20*20, b//20*20) 
                                     for r,g,b,a in pixels]).most_common(3)
                print(f"  {char}: {[(f'#{r:02X}{g:02X}{b:02X}', cnt) for (r,g,b), cnt in top_colors]}")

if __name__ == "__main__":
    make_preview()
