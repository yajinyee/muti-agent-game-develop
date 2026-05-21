# -*- coding: utf-8 -*-
"""
Pixel Art v5 - Faithful to original Chiikawa design
Reference: chiikawa.fandom.com
- Chiikawa: WHITE fur, round ears, small body
- Hachiware: WHITE with blue stripes, pointed ears
- Usagi: WHITE, long ears, red eyes
Key fixes:
- Correct white/cream color (not brown/grey)
- Smaller mouth, correct position
- Smaller blush dots
- Smaller body relative to head
"""
from PIL import Image
import os
import math

OUTPUT_BASE = r"D:\Kiro\client\chiikawa-pixel\assets\sprites"

def new_img(w, h):
    return Image.new("RGBA", (w, h), (0, 0, 0, 0))

def px(img, x, y, c):
    if 0 <= x < img.width and 0 <= y < img.height:
        img.putpixel((x, y), c)

def fill_circle(img, cx, cy, r, color):
    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            if (x-cx)**2 + (y-cy)**2 <= r**2:
                px(img, x, y, color)

def fill_circle_shaded(img, cx, cy, r, base, shadow_offset=(1, 1)):
    """帶陰影的圓形：左上亮，右下暗"""
    r_val, g_val, b_val = base[:3]
    light = (min(255, r_val+20), min(255, g_val+20), min(255, b_val+20), 255)
    mid   = (r_val, g_val, b_val, 255)
    dark  = (max(0, r_val-30), max(0, g_val-30), max(0, b_val-30), 255)

    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            if (x-cx)**2 + (y-cy)**2 > r**2:
                continue
            # 光源左上
            nx = (x - cx) / max(r, 1)
            ny = (y - cy) / max(r, 1)
            dot = -(nx * (-0.7) + ny * (-0.7))
            if dot > 0.2:
                c = light
            elif dot < -0.1:
                c = dark
            else:
                c = mid
            px(img, x, y, c)

def outline_circle(img, cx, cy, r, color):
    for y in range(cy - r - 1, cy + r + 2):
        for x in range(cx - r - 1, cx + r + 2):
            d = math.sqrt((x-cx)**2 + (y-cy)**2)
            if r + 0.1 <= d <= r + 1.4:
                px(img, x, y, color)

def draw_eye_simple(img, x, y, iris_color=(60, 40, 20, 255)):
    """簡潔的 2x2 眼睛，有高光"""
    BLACK = (15, 8, 3, 255)
    WHITE = (255, 255, 255, 255)
    # 眼白背景
    px(img, x,   y,   WHITE)
    px(img, x+1, y,   WHITE)
    px(img, x,   y+1, WHITE)
    px(img, x+1, y+1, WHITE)
    # 瞳孔
    px(img, x+1, y+1, BLACK)
    px(img, x,   y+1, iris_color)
    # 高光（左上）
    px(img, x, y, WHITE)

def gen_chiikawa_v5(state="idle"):
    """
    吉伊卡哇 v5 - 32x32
    白色毛皮，圓耳，小身體
    顏色：接近純白 #FFFCF5
    """
    img = new_img(32, 32)

    # 正確顏色（接近純白，略帶奶油色）
    WHITE  = (255, 252, 245)   # 主體白色
    OUTLINE= (50, 30, 15, 255) # 深棕輪廓
    BLUSH  = (255, 160, 155, 255) # 腮紅（小點）
    PINK   = (255, 140, 190, 255) # 討伐棒
    GLOW   = (255, 215, 235, 255) # 劍氣光暈
    EAR_IN = (255, 200, 195, 255) # 耳朵內側（淡粉）

    dy = -2 if state == "attack" else (1 if state == "bigwin" else 0)

    # 耳朵（小圓，帶內側粉紅）
    fill_circle_shaded(img, 9,  8+dy, 4, WHITE)
    fill_circle_shaded(img, 23, 8+dy, 4, WHITE)
    fill_circle(img, 9,  8+dy, 2, EAR_IN)
    fill_circle(img, 23, 8+dy, 2, EAR_IN)
    outline_circle(img, 9,  8+dy, 4, OUTLINE)
    outline_circle(img, 23, 8+dy, 4, OUTLINE)

    # 大圓頭（帶陰影）
    fill_circle_shaded(img, 16, 15+dy, 10, WHITE)
    outline_circle(img, 16, 15+dy, 10, OUTLINE)

    # 眼睛（閉著，努力感）
    if state == "bigwin":
        # 大獎：睜大眼睛
        draw_eye_simple(img, 10, 13+dy)
        draw_eye_simple(img, 19, 13+dy)
    else:
        # 閉眼：簡單弧線
        for i in range(4):
            px(img, 10+i, 14+dy, OUTLINE)
            px(img, 19+i, 14+dy, OUTLINE)
        # 眼睫毛（上方一點）
        px(img, 10, 13+dy, OUTLINE)
        px(img, 13, 13+dy, OUTLINE)
        px(img, 19, 13+dy, OUTLINE)
        px(img, 22, 13+dy, OUTLINE)

    # 腮紅（小圓點，2px）
    fill_circle(img, 9,  18+dy, 2, BLUSH)
    fill_circle(img, 23, 18+dy, 2, BLUSH)

    # 嘴巴（小，在臉中央，不是鬍子）
    # 只有 3 個像素，位置在眼睛下方 4px
    px(img, 15, 19+dy, OUTLINE)
    px(img, 16, 20+dy, OUTLINE)
    px(img, 17, 19+dy, OUTLINE)

    # 小身體（比頭小很多）
    fill_circle_shaded(img, 16, 26+dy, 4, WHITE)
    outline_circle(img, 16, 26+dy, 4, OUTLINE)

    # 討伐棒
    if state == "attack":
        for i in range(6):
            px(img, 22+i, 11+dy-i, PINK)
            px(img, 22+i, 12+dy-i, (200, 100, 150, 255))
        fill_circle(img, 28, 5+dy, 2, GLOW)
    else:
        for i in range(5):
            px(img, 22+i, 13+dy-i, PINK)
        fill_circle(img, 27, 9+dy, 2, GLOW)

    return img

def gen_hachiware_v5(state="idle"):
    """
    小八 v5 - 32x32
    白色帶藍色條紋，尖耳貓
    """
    img = new_img(32, 32)

    WHITE  = (248, 248, 248)
    OUTLINE= (30, 40, 70, 255)
    STRIPE = (80, 120, 200)
    BLUE   = (100, 150, 240, 255)
    GLOW   = (180, 210, 255, 255)
    EAR_IN = (190, 210, 250, 255)

    dy = -2 if state == "attack" else (1 if state == "bigwin" else 0)

    # 尖耳朵（三角形）
    for i in range(7):
        for j in range(i+1):
            px(img, 8+j,  5+dy+i, (*WHITE, 255))
            px(img, 24-j, 5+dy+i, (*WHITE, 255))
    for i in range(7):
        px(img, 8,    5+dy+i, OUTLINE)
        px(img, 8+i,  5+dy+i, OUTLINE)
        px(img, 24,   5+dy+i, OUTLINE)
        px(img, 24-i, 5+dy+i, OUTLINE)
    # 耳朵內側（藍色）
    for i in range(4):
        for j in range(i+1):
            px(img, 9+j,  7+dy+i, EAR_IN)
            px(img, 23-j, 7+dy+i, EAR_IN)

    # 大圓頭
    fill_circle_shaded(img, 16, 16+dy, 10, WHITE)
    outline_circle(img, 16, 16+dy, 10, OUTLINE)

    # 藍色條紋（特徵）
    for sx in [10, 14, 18, 22]:
        for y in range(10+dy, 22+dy):
            if (16+dy - y)**2 + (16 - sx)**2 < 100:
                px(img, sx, y, (*STRIPE, 255))

    # 眼睛（睜開）
    draw_eye_simple(img, 10, 13+dy, (50, 90, 190, 255))
    draw_eye_simple(img, 20, 13+dy, (50, 90, 190, 255))

    # 嘴巴（微笑）
    px(img, 14, 19+dy, OUTLINE)
    px(img, 15, 20+dy, OUTLINE)
    px(img, 16, 20+dy, OUTLINE)
    px(img, 17, 20+dy, OUTLINE)
    px(img, 18, 19+dy, OUTLINE)

    # 小身體
    fill_circle_shaded(img, 16, 26+dy, 4, WHITE)
    outline_circle(img, 16, 26+dy, 4, OUTLINE)

    # 討伐棒
    if state == "attack":
        for i in range(6):
            px(img, 22+i, 11+dy-i, BLUE)
            px(img, 22+i, 12+dy-i, (70, 110, 190, 255))
        fill_circle(img, 28, 5+dy, 2, GLOW)
    else:
        for i in range(5):
            px(img, 22+i, 13+dy-i, BLUE)
        fill_circle(img, 27, 9+dy, 2, GLOW)

    return img

def gen_usagi_v5(state="idle"):
    """
    烏薩奇 v5 - 32x32
    白色，長耳兔，紅眼
    """
    img = new_img(32, 32)

    WHITE  = (248, 248, 248)
    OUTLINE= (60, 40, 60, 255)
    PINK_IN= (255, 175, 175, 255)
    YELLOW = (255, 215, 30, 255)
    GLOW   = (255, 245, 140, 255)

    dy = -3 if state == "attack" else (1 if state == "bigwin" else 0)

    # 長耳朵（attack 時 dy=-3，耳朵起點要 clamp 到 0）
    ear_top = max(0, 0 + dy)
    for y in range(ear_top, 13+dy):
        for x in range(7, 14):
            if abs(x - 10) <= 3:
                px(img, x, y, (*WHITE, 255))
        for x in range(19, 26):
            if abs(x - 22) <= 3:
                px(img, x, y, (*WHITE, 255))
    # 耳朵輪廓
    for y in range(ear_top, 13+dy):
        px(img, 7,  y, OUTLINE)
        px(img, 13, y, OUTLINE)
        px(img, 19, y, OUTLINE)
        px(img, 25, y, OUTLINE)
    # 耳朵頂部
    for x in range(7, 14):
        px(img, x, ear_top, OUTLINE)
    for x in range(19, 26):
        px(img, x, ear_top, OUTLINE)
    # 耳朵內側（粉紅）
    for y in range(ear_top+1, 12+dy):
        for x in range(8, 13):
            if abs(x - 10) <= 2:
                px(img, x, y, PINK_IN)
        for x in range(20, 25):
            if abs(x - 22) <= 2:
                px(img, x, y, PINK_IN)

    # 大圓頭
    fill_circle_shaded(img, 16, 17+dy, 10, WHITE)
    outline_circle(img, 16, 17+dy, 10, OUTLINE)

    # 眼睛（紅色）
    draw_eye_simple(img, 10, 14+dy, (200, 40, 40, 255))
    draw_eye_simple(img, 20, 14+dy, (200, 40, 40, 255))

    # 嘴巴（興奮）
    px(img, 14, 20+dy, OUTLINE)
    px(img, 15, 21+dy, OUTLINE)
    px(img, 16, 21+dy, OUTLINE)
    px(img, 17, 21+dy, OUTLINE)
    px(img, 18, 20+dy, OUTLINE)
    # 牙齒
    px(img, 15, 21+dy, (245, 245, 245, 255))
    px(img, 17, 21+dy, (245, 245, 245, 255))

    # 小身體
    fill_circle_shaded(img, 16, 27+dy, 4, WHITE)
    outline_circle(img, 16, 27+dy, 4, OUTLINE)

    # 討伐棒（旋轉感）
    if state == "attack":
        for i in range(6):
            px(img, 21+i, 11+dy-i, YELLOW)
            px(img, 22+i, 11+dy-i, (200, 165, 20, 255))
        fill_circle(img, 28, 5+dy, 3, GLOW)
    else:
        for i in range(5):
            px(img, 22+i, 13+dy-i, YELLOW)
        fill_circle(img, 27, 9+dy, 2, GLOW)

    return img

def generate_all_v5():
    chars_dir = os.path.join(OUTPUT_BASE, "characters")
    os.makedirs(chars_dir, exist_ok=True)

    print("Characters v5 - White fur, correct proportions")
    fns = [
        ("chiikawa", gen_chiikawa_v5),
        ("hachiware", gen_hachiware_v5),
        ("usagi", gen_usagi_v5),
    ]
    for name, fn in fns:
        for state in ["idle", "attack", "bigwin"]:
            img = fn(state)
            # 2x 放大
            img = img.resize((64, 64), Image.NEAREST)
            path = os.path.join(chars_dir, f"{name}_{state}.png")
            img.save(path)
            print(f"  OK {name}_{state}.png")

if __name__ == "__main__":
    print("Pixel Art v5...")
    generate_all_v5()
    print("Done!")
