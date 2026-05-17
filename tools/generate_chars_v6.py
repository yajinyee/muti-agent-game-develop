# -*- coding: utf-8 -*-
"""
角色 v6 - 48x48 基礎尺寸，更豐富細節
改善點：
1. 更大的頭部（佔比更高，更 chibi）
2. 更清晰的眼睛（4x4 眼白 + 2x2 瞳孔 + 高光）
3. 更明顯的腮紅
4. 更完整的身體（手臂可見）
5. 更好的陰影（3色漸層，左上光源）
6. 輸出 96x96（2x 放大）
7. 統一調色板（skill-pixel-art-quality-2026）
"""
from PIL import Image
import os
import math

OUTPUT_BASE = r"D:\Kiro\client\chiikawa-pixel\assets\sprites"
SIZE = 48  # 基礎尺寸
OUT_SIZE = 96  # 輸出尺寸（2x 放大）

# ── 統一調色板（來自 skill-pixel-art-quality-2026.md）────────────────────────
CHIIKAWA_PALETTE = {
    "body":      (255, 252, 245),
    "shadow":    (220, 215, 205),
    "highlight": (255, 255, 255),
    "outline":   (45, 25, 10, 255),
    "blush":     (255, 155, 150, 255),
    "pink_rod":  (255, 130, 185, 255),
    "ear_in":    (255, 195, 190, 255),
    "glow":      (255, 210, 230, 255),
}

HACHIWARE_PALETTE = {
    "body":      (248, 248, 248),
    "shadow":    (210, 210, 215),
    "highlight": (255, 255, 255),
    "outline":   (25, 35, 65, 255),
    "stripe":    (75, 115, 195),
    "blue_rod":  (95, 145, 235, 255),
    "ear_in":    (185, 205, 248, 255),
    "glow":      (175, 205, 255, 255),
}

USAGI_PALETTE = {
    "body":      (248, 248, 248),
    "shadow":    (210, 210, 215),
    "highlight": (255, 255, 255),
    "outline":   (55, 35, 55, 255),
    "ear_pink":  (255, 170, 170, 255),
    "yellow_rod":(255, 210, 25, 255),
    "glow":      (255, 240, 135, 255),
}

def new_img():
    return Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

def px(img, x, y, c):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), c)

def fill_circle(img, cx, cy, r, color):
    for y in range(max(0, cy-r), min(SIZE, cy+r+1)):
        for x in range(max(0, cx-r), min(SIZE, cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 <= r**2:
                px(img, x, y, color)

def fill_circle_3shade(img, cx, cy, r, base_rgb):
    """3色陰影圓形：左上亮，右下暗"""
    rv, gv, bv = base_rgb
    LIGHT = (min(255,rv+30), min(255,gv+30), min(255,bv+30), 255)
    MID   = (rv, gv, bv, 255)
    DARK  = (max(0,rv-35), max(0,gv-35), max(0,bv-35), 255)
    for y in range(max(0,cy-r), min(SIZE,cy+r+1)):
        for x in range(max(0,cx-r), min(SIZE,cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 > r**2:
                continue
            nx_ = (x-cx)/max(r,1)
            ny_ = (y-cy)/max(r,1)
            dot = -(nx_*(-0.7) + ny_*(-0.7))
            c = LIGHT if dot > 0.2 else (DARK if dot < -0.1 else MID)
            px(img, x, y, c)

def outline_circle(img, cx, cy, r, color):
    for y in range(max(0,cy-r-1), min(SIZE,cy+r+2)):
        for x in range(max(0,cx-r-1), min(SIZE,cx+r+2)):
            d = math.sqrt((x-cx)**2 + (y-cy)**2)
            if r+0.1 <= d <= r+1.4:
                px(img, x, y, color)

def draw_eye_v6(img, x, y, iris_rgb=(60,40,20)):
    """4x4 眼睛：眼白+瞳孔+高光"""
    WHITE = (255, 255, 255, 255)
    BLACK = (15, 8, 3, 255)
    IRIS  = (*iris_rgb, 255)
    # 眼白（3x3）
    for dx in range(3):
        for dy in range(3):
            px(img, x+dx, y+dy, WHITE)
    # 瞳孔（2x2，右下）
    px(img, x+1, y+1, IRIS)
    px(img, x+2, y+1, BLACK)
    px(img, x+1, y+2, BLACK)
    px(img, x+2, y+2, BLACK)
    # 高光（左上）
    px(img, x, y, WHITE)
    # 眼框
    px(img, x, y+3, BLACK)
    px(img, x+1, y+3, BLACK)
    px(img, x+2, y+3, BLACK)
    px(img, x+3, y, BLACK)
    px(img, x+3, y+1, BLACK)
    px(img, x+3, y+2, BLACK)

def gen_chiikawa_v6(state="idle"):
    img = new_img()
    WHITE  = (255, 252, 245)
    OUTLINE= (45, 25, 10, 255)
    BLUSH  = (255, 155, 150, 255)
    PINK   = (255, 130, 185, 255)
    GLOW   = (255, 210, 230, 255)
    EAR_IN = (255, 195, 190, 255)

    # 動作偏移
    if state == "attack":
        dy = -3
    elif state == "bigwin":
        dy = 1
    else:
        dy = 0

    # 耳朵（圓形，帶內側）
    fill_circle_3shade(img, 13, 10+dy, 5, WHITE)
    fill_circle_3shade(img, 35, 10+dy, 5, WHITE)
    fill_circle(img, 13, 10+dy, 3, EAR_IN[:3])
    fill_circle(img, 35, 10+dy, 3, EAR_IN[:3])
    outline_circle(img, 13, 10+dy, 5, OUTLINE)
    outline_circle(img, 35, 10+dy, 5, OUTLINE)

    # 大圓頭（半徑 14，佔 48px 的 58%）
    fill_circle_3shade(img, 24, 22+dy, 14, WHITE)
    outline_circle(img, 24, 22+dy, 14, OUTLINE)

    # 眼睛
    if state == "bigwin":
        draw_eye_v6(img, 13, 18+dy)
        draw_eye_v6(img, 28, 18+dy)
    else:
        # 閉眼弧線（更自然）
        for i in range(5):
            px(img, 13+i, 20+dy, OUTLINE)
            px(img, 27+i, 20+dy, OUTLINE)
        px(img, 13, 19+dy, OUTLINE)
        px(img, 17, 19+dy, OUTLINE)
        px(img, 27, 19+dy, OUTLINE)
        px(img, 31, 19+dy, OUTLINE)

    # 腮紅（橢圓，更明顯）
    for dx in range(-3, 4):
        for dy2 in range(-1, 2):
            alpha = max(0, 200 - abs(dx)*40 - abs(dy2)*60)
            px(img, 13+dx, 25+dy+dy2, (*BLUSH[:3], alpha))
            px(img, 35+dx, 25+dy+dy2, (*BLUSH[:3], alpha))

    # 嘴巴（V形，3像素）
    px(img, 22, 27+dy, OUTLINE)
    px(img, 23, 28+dy, OUTLINE)
    px(img, 24, 27+dy, OUTLINE)

    # 身體（橢圓，比頭小）
    fill_circle_3shade(img, 24, 38+dy, 7, WHITE)
    outline_circle(img, 24, 38+dy, 7, OUTLINE)

    # 手臂（小圓）
    fill_circle_3shade(img, 14, 36+dy, 4, WHITE)
    fill_circle_3shade(img, 34, 36+dy, 4, WHITE)
    outline_circle(img, 14, 36+dy, 4, OUTLINE)
    outline_circle(img, 34, 36+dy, 4, OUTLINE)

    # 討伐棒
    if state == "attack":
        # 舉棒攻擊（斜向右上）
        for i in range(8):
            px(img, 34+i, 14+dy-i, PINK)
            px(img, 34+i, 15+dy-i, (195, 95, 145, 255))
        fill_circle(img, 42, 6+dy, 3, GLOW)
        outline_circle(img, 42, 6+dy, 3, OUTLINE)
    else:
        # 持棒（斜向右上，較小）
        for i in range(6):
            px(img, 33+i, 18+dy-i, PINK)
        fill_circle(img, 39, 13+dy, 2, GLOW)

    return img


def gen_hachiware_v6(state="idle"):
    img = new_img()
    WHITE  = (248, 248, 248)
    OUTLINE= (25, 35, 65, 255)
    STRIPE = (75, 115, 195)
    BLUE   = (95, 145, 235, 255)
    GLOW   = (175, 205, 255, 255)
    EAR_IN = (185, 205, 248, 255)

    if state == "attack":
        dy = -3
    elif state == "bigwin":
        dy = 1
    else:
        dy = 0

    # 尖耳朵（三角形）
    for i in range(9):
        for j in range(i+1):
            px(img, 10+j, 6+dy+i, (*WHITE, 255))
            px(img, 38-j, 6+dy+i, (*WHITE, 255))
    for i in range(9):
        px(img, 10, 6+dy+i, OUTLINE)
        px(img, 10+i, 6+dy+i, OUTLINE)
        px(img, 38, 6+dy+i, OUTLINE)
        px(img, 38-i, 6+dy+i, OUTLINE)
    # 耳朵內側
    for i in range(5):
        for j in range(i+1):
            px(img, 11+j, 8+dy+i, EAR_IN[:3])
            px(img, 37-j, 8+dy+i, EAR_IN[:3])

    # 大圓頭
    fill_circle_3shade(img, 24, 23+dy, 14, WHITE)
    outline_circle(img, 24, 23+dy, 14, OUTLINE)

    # 藍色條紋（特徵，3條）
    for sx in [16, 22, 28]:
        for y in range(14+dy, 32+dy):
            if (24-y)**2 + (24-sx)**2 < 196:
                px(img, sx, y, (*STRIPE, 255))

    # 眼睛（藍色虹膜）
    draw_eye_v6(img, 14, 19+dy, (45, 85, 185))
    draw_eye_v6(img, 29, 19+dy, (45, 85, 185))

    # 嘴巴（微笑）
    px(img, 21, 29+dy, OUTLINE)
    px(img, 22, 30+dy, OUTLINE)
    px(img, 23, 30+dy, OUTLINE)
    px(img, 24, 30+dy, OUTLINE)
    px(img, 25, 29+dy, OUTLINE)

    # 身體
    fill_circle_3shade(img, 24, 38+dy, 7, WHITE)
    outline_circle(img, 24, 38+dy, 7, OUTLINE)

    # 手臂
    fill_circle_3shade(img, 14, 36+dy, 4, WHITE)
    fill_circle_3shade(img, 34, 36+dy, 4, WHITE)
    outline_circle(img, 14, 36+dy, 4, OUTLINE)
    outline_circle(img, 34, 36+dy, 4, OUTLINE)

    # 討伐棒（藍色）
    if state == "attack":
        for i in range(8):
            px(img, 34+i, 14+dy-i, BLUE)
            px(img, 34+i, 15+dy-i, (65, 105, 185, 255))
        fill_circle(img, 42, 6+dy, 3, GLOW)
        outline_circle(img, 42, 6+dy, 3, OUTLINE)
    else:
        for i in range(6):
            px(img, 33+i, 18+dy-i, BLUE)
        fill_circle(img, 39, 13+dy, 2, GLOW)

    return img


def gen_usagi_v6(state="idle"):
    img = new_img()
    WHITE  = (248, 248, 248)
    OUTLINE= (55, 35, 55, 255)
    PINK_IN= (255, 170, 170, 255)
    YELLOW = (255, 210, 25, 255)
    GLOW   = (255, 240, 135, 255)

    if state == "attack":
        dy = -3
    elif state == "bigwin":
        dy = 1
    else:
        dy = 0

    # 長耳朵（兩條）
    ear_top = max(0, dy)
    for y in range(ear_top, 16+dy):
        for x in range(9, 16):
            if abs(x-12) <= 3:
                px(img, x, y, (*WHITE, 255))
        for x in range(33, 40):
            if abs(x-36) <= 3:
                px(img, x, y, (*WHITE, 255))
    # 耳朵輪廓
    for y in range(ear_top, 16+dy):
        px(img, 9,  y, OUTLINE)
        px(img, 15, y, OUTLINE)
        px(img, 33, y, OUTLINE)
        px(img, 39, y, OUTLINE)
    for x in range(9, 16):
        px(img, x, ear_top, OUTLINE)
    for x in range(33, 40):
        px(img, x, ear_top, OUTLINE)
    # 耳朵內側（粉紅）
    for y in range(ear_top+1, 15+dy):
        for x in range(10, 15):
            if abs(x-12) <= 2:
                px(img, x, y, PINK_IN[:3])
        for x in range(34, 39):
            if abs(x-36) <= 2:
                px(img, x, y, PINK_IN[:3])

    # 大圓頭
    fill_circle_3shade(img, 24, 24+dy, 14, WHITE)
    outline_circle(img, 24, 24+dy, 14, OUTLINE)

    # 眼睛（紅色虹膜）
    draw_eye_v6(img, 14, 20+dy, (195, 35, 35))
    draw_eye_v6(img, 29, 20+dy, (195, 35, 35))

    # 嘴巴（興奮，帶牙齒）
    px(img, 21, 30+dy, OUTLINE)
    px(img, 22, 31+dy, OUTLINE)
    px(img, 23, 31+dy, OUTLINE)
    px(img, 24, 31+dy, OUTLINE)
    px(img, 25, 30+dy, OUTLINE)
    # 牙齒
    px(img, 22, 31+dy, (245, 245, 245, 255))
    px(img, 24, 31+dy, (245, 245, 245, 255))

    # 身體
    fill_circle_3shade(img, 24, 39+dy, 7, WHITE)
    outline_circle(img, 24, 39+dy, 7, OUTLINE)

    # 手臂
    fill_circle_3shade(img, 14, 37+dy, 4, WHITE)
    fill_circle_3shade(img, 34, 37+dy, 4, WHITE)
    outline_circle(img, 14, 37+dy, 4, OUTLINE)
    outline_circle(img, 34, 37+dy, 4, OUTLINE)

    # 討伐棒（黃色，旋轉感）
    if state == "attack":
        for i in range(8):
            px(img, 34+i, 14+dy-i, YELLOW)
            px(img, 35+i, 14+dy-i, (195, 160, 15, 255))
        fill_circle(img, 42, 6+dy, 4, GLOW)
        outline_circle(img, 42, 6+dy, 4, OUTLINE)
    else:
        for i in range(6):
            px(img, 33+i, 18+dy-i, YELLOW)
        fill_circle(img, 39, 13+dy, 2, GLOW)

    return img


def generate_all_v6():
    chars_dir = os.path.join(OUTPUT_BASE, "characters")
    os.makedirs(chars_dir, exist_ok=True)

    print("=== 角色 v6（48x48 → 96x96）===")
    fns = [
        ("chiikawa",  gen_chiikawa_v6),
        ("hachiware", gen_hachiware_v6),
        ("usagi",     gen_usagi_v6),
    ]
    for name, fn in fns:
        for state in ["idle", "attack", "bigwin"]:
            img = fn(state)
            # 2x 放大（NEAREST 保持像素感）
            img = img.resize((OUT_SIZE, OUT_SIZE), Image.NEAREST)
            path = os.path.join(chars_dir, f"{name}_{state}.png")
            img.save(path)
            # 品質驗證（skill-pixel-art-quality-2026）
            non_t = sum(1 for px_v in img.getdata() if px_v[3] > 10)
            pct = non_t * 100 // (OUT_SIZE * OUT_SIZE)
            status = "✅" if pct >= 40 else "⚠️"
            print(f"  {status} {name}_{state}.png ({OUT_SIZE}x{OUT_SIZE}, {pct}% 非透明)")

if __name__ == "__main__":
    generate_all_v6()
    print("\nDone!")
