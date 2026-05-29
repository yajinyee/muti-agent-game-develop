#!/usr/bin/env python3
"""
generate_targets_day326.py — DAY-326 T221-T223 精靈圖生成
T221: 幸運骰子獎勵魚 (Dice Bonus, BGaming Shark & Spark Hold & Win 2026-05-25)
T222: 幸運雙Bonus魚 (Dual Bonus, BGaming Fishing Club 2 2026-04)
T223: 幸運Coin Respin魚 (Coin Respin, BGaming Shark & Spark Hold & Win 2026-05-28)
"""
import os
from PIL import Image, ImageDraw

OUTPUT_DIR = os.path.join(os.path.dirname(__file__), "..", "client", "chiikawa-pixel", "assets", "sprites", "targets")
SIZE = 64

def draw_outline(draw, size=64, color=(0, 0, 0, 255), width=1):
    """Draw 1px black outline around the sprite."""
    for x in range(size):
        for y in range(size):
            pass  # outline handled by border drawing

def create_base_fish(color_body, color_fin, color_eye=(255, 255, 255, 255)):
    """Create a base fish shape."""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # Body (ellipse)
    draw.ellipse([8, 16, 52, 48], fill=color_body, outline=(0, 0, 0, 255))
    # Tail
    draw.polygon([(8, 32), (0, 16), (0, 48)], fill=color_fin, outline=(0, 0, 0, 255))
    # Fin top
    draw.polygon([(20, 16), (35, 8), (40, 16)], fill=color_fin, outline=(0, 0, 0, 255))
    # Eye
    draw.ellipse([40, 22, 50, 32], fill=color_eye, outline=(0, 0, 0, 255))
    draw.ellipse([43, 25, 48, 30], fill=(0, 0, 0, 255))
    
    return img, draw

def draw_rays(draw, cx, cy, count, length, color, width=1):
    """Draw rays from center."""
    import math
    for i in range(count):
        angle = (2 * math.pi * i) / count
        x2 = int(cx + length * math.cos(angle))
        y2 = int(cy + length * math.sin(angle))
        draw.line([(cx, cy), (x2, y2)], fill=color, width=width)

def generate_t221_dice_bonus():
    """T221 幸運骰子獎勵魚 — 橙紅色魚身 + 骰子符號 + 光芒"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # Glow background
    for r in range(30, 0, -3):
        alpha = int(80 * (1 - r/30))
        draw.ellipse([32-r, 32-r, 32+r, 32+r], fill=(255, 100, 30, alpha))
    
    # 16 rays
    draw_rays(draw, 32, 32, 16, 28, (255, 180, 50, 200), 1)
    
    # Fish body (orange-red)
    draw.ellipse([8, 16, 52, 48], fill=(220, 80, 20, 255), outline=(0, 0, 0, 255))
    # Tail
    draw.polygon([(8, 32), (0, 16), (0, 48)], fill=(180, 60, 10, 255), outline=(0, 0, 0, 255))
    # Fin
    draw.polygon([(20, 16), (35, 8), (40, 16)], fill=(180, 60, 10, 255), outline=(0, 0, 0, 255))
    # Eye
    draw.ellipse([40, 22, 50, 32], fill=(255, 255, 255, 255), outline=(0, 0, 0, 255))
    draw.ellipse([43, 25, 48, 30], fill=(0, 0, 0, 255))
    
    # Dice symbol on body (white square with dots)
    draw.rectangle([18, 22, 34, 38], fill=(255, 255, 255, 230), outline=(0, 0, 0, 255))
    # Dice dots (showing 6)
    dot_color = (0, 0, 0, 255)
    draw.ellipse([20, 24, 23, 27], fill=dot_color)
    draw.ellipse([29, 24, 32, 27], fill=dot_color)
    draw.ellipse([20, 29, 23, 32], fill=dot_color)
    draw.ellipse([29, 29, 32, 32], fill=dot_color)
    draw.ellipse([20, 34, 23, 37], fill=dot_color)
    draw.ellipse([29, 34, 32, 37], fill=dot_color)
    
    # "x300" text area
    draw.rectangle([36, 42, 58, 52], fill=(255, 200, 0, 200), outline=(0, 0, 0, 255))
    
    return img

def generate_t222_dual_bonus():
    """T222 幸運雙Bonus魚 — 藍色魚身 + 雙選擇符號 + 光芒"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # Glow background
    for r in range(30, 0, -3):
        alpha = int(80 * (1 - r/30))
        draw.ellipse([32-r, 32-r, 32+r, 32+r], fill=(30, 100, 255, alpha))
    
    # 18 rays
    draw_rays(draw, 32, 32, 18, 28, (100, 200, 255, 200), 1)
    
    # Fish body (blue)
    draw.ellipse([8, 16, 52, 48], fill=(30, 100, 220, 255), outline=(0, 0, 0, 255))
    # Tail
    draw.polygon([(8, 32), (0, 16), (0, 48)], fill=(20, 70, 180, 255), outline=(0, 0, 0, 255))
    # Fin
    draw.polygon([(20, 16), (35, 8), (40, 16)], fill=(20, 70, 180, 255), outline=(0, 0, 0, 255))
    # Eye
    draw.ellipse([40, 22, 50, 32], fill=(255, 255, 255, 255), outline=(0, 0, 0, 255))
    draw.ellipse([43, 25, 48, 30], fill=(0, 0, 0, 255))
    
    # Dual choice symbol: two circles with A and B
    draw.ellipse([14, 22, 24, 32], fill=(255, 220, 0, 230), outline=(0, 0, 0, 255))
    draw.ellipse([26, 22, 36, 32], fill=(0, 220, 100, 230), outline=(0, 0, 0, 255))
    # A label
    draw.line([(16, 31), (19, 23)], fill=(0, 0, 0, 255), width=1)
    draw.line([(19, 23), (22, 31)], fill=(0, 0, 0, 255), width=1)
    draw.line([(17, 28), (21, 28)], fill=(0, 0, 0, 255), width=1)
    # B label
    draw.line([(28, 23), (28, 31)], fill=(0, 0, 0, 255), width=1)
    draw.arc([28, 23, 34, 27], 270, 90, fill=(0, 0, 0, 255), width=1)
    draw.arc([28, 27, 34, 31], 270, 90, fill=(0, 0, 0, 255), width=1)
    
    # "x42" text area
    draw.rectangle([36, 42, 58, 52], fill=(0, 200, 255, 200), outline=(0, 0, 0, 255))
    
    return img

def generate_t223_coin_respin():
    """T223 幸運Coin Respin魚 — 金色魚身 + 9格盤面 + 光芒"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # Glow background (gold)
    for r in range(30, 0, -3):
        alpha = int(100 * (1 - r/30))
        draw.ellipse([32-r, 32-r, 32+r, 32+r], fill=(255, 200, 0, alpha))
    
    # 20 rays (most impressive)
    draw_rays(draw, 32, 32, 20, 30, (255, 220, 50, 220), 1)
    
    # Fish body (bright gold)
    draw.ellipse([8, 14, 54, 50], fill=(255, 200, 0, 255), outline=(0, 0, 0, 255))
    # Tail
    draw.polygon([(8, 32), (0, 14), (0, 50)], fill=(220, 160, 0, 255), outline=(0, 0, 0, 255))
    # Fin
    draw.polygon([(20, 14), (36, 6), (42, 14)], fill=(220, 160, 0, 255), outline=(0, 0, 0, 255))
    # Eye
    draw.ellipse([42, 20, 52, 30], fill=(255, 255, 255, 255), outline=(0, 0, 0, 255))
    draw.ellipse([45, 23, 50, 28], fill=(0, 0, 0, 255))
    
    # 3x3 grid (Hold & Win board)
    grid_x, grid_y = 14, 20
    cell = 8
    for row in range(3):
        for col in range(3):
            x = grid_x + col * cell
            y = grid_y + row * cell
            # Alternate filled/empty cells
            if (row + col) % 2 == 0:
                draw.rectangle([x, y, x+cell-1, y+cell-1], fill=(255, 240, 100, 220), outline=(0, 0, 0, 255))
                # Coin dot
                draw.ellipse([x+2, y+2, x+cell-3, y+cell-3], fill=(255, 180, 0, 255))
            else:
                draw.rectangle([x, y, x+cell-1, y+cell-1], fill=(180, 140, 0, 150), outline=(0, 0, 0, 255))
    
    # "x42.5" text area
    draw.rectangle([36, 44, 62, 54], fill=(255, 220, 0, 220), outline=(0, 0, 0, 255))
    
    return img

def save_sprite(img, filename):
    path = os.path.join(OUTPUT_DIR, filename)
    img.save(path, "PNG")
    size_kb = os.path.getsize(path) / 1024
    print(f"  Saved: {filename} ({size_kb:.1f} KB)")

if __name__ == "__main__":
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    print("DAY-326 Target Sprite Generation")
    print("=" * 40)
    
    print("Generating T221 (Dice Bonus Fish)...")
    save_sprite(generate_t221_dice_bonus(), "T221_dice_bonus.png")
    
    print("Generating T222 (Dual Bonus Fish)...")
    save_sprite(generate_t222_dual_bonus(), "T222_dual_bonus.png")
    
    print("Generating T223 (Coin Respin Fish)...")
    save_sprite(generate_t223_coin_respin(), "T223_coin_respin.png")
    
    print("\nAll DAY-326 sprites generated successfully!")
    print("T221: Dice Bonus (BGaming Shark & Spark Hold & Win 2026-05-25)")
    print("T222: Dual Bonus (BGaming Fishing Club 2 2026-04)")
    print("T223: Coin Respin (BGaming Shark & Spark Hold & Win 2026-05-28)")
