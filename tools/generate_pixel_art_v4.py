# -*- coding: utf-8 -*-
"""
像素美術生成器 v4
新增：3色陰影法 + Pillow Shading + 光源方向
來源：pixnote.net, hitpaw.com, sprite-ai.art
目標：從 Level 2 提升到 Level 3（55 -> 70分）
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

def get_px(img, x, y):
    if 0 <= x < img.width and 0 <= y < img.height:
        return img.getpixel((x, y))
    return (0, 0, 0, 0)

def shade_color(base_rgb, factor):
    """調整顏色亮度"""
    r, g, b = base_rgb[:3]
    return (
        max(0, min(255, int(r * factor))),
        max(0, min(255, int(g * factor))),
        max(0, min(255, int(b * factor))),
        255
    )

def fill_circle_shaded(img, cx, cy, r, base_color, light_dir=(-1, -1)):
    """
    填充帶陰影的圓形（Pillow Shading）
    光源方向：左上(-1,-1)
    3色：亮色(1.3x)、中間色(1.0x)、暗色(0.65x)
    """
    lx, ly = light_dir
    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            dist = math.sqrt((x - cx)**2 + (y - cy)**2)
            if dist > r:
                continue
            # 計算光照強度（點積）
            nx = (x - cx) / max(r, 1)
            ny = (y - cy) / max(r, 1)
            dot = -(nx * lx + ny * ly)  # 負號因為光源方向是指向光源
            # 邊緣加深（Pillow Shading）
            edge_factor = 1.0 - (dist / r) * 0.3
            # 合併光照和邊緣
            brightness = (dot * 0.4 + 0.6) * edge_factor
            brightness = max(0.55, min(1.35, brightness))
            color = shade_color(base_color, brightness)
            px(img, x, y, color)

def draw_outline_circle(img, cx, cy, r, color):
    """畫圓形輪廓"""
    for y in range(cy - r - 1, cy + r + 2):
        for x in range(cx - r - 1, cx + r + 2):
            dist = math.sqrt((x - cx)**2 + (y - cy)**2)
            if r + 0.2 <= dist <= r + 1.3:
                px(img, x, y, color)

def fill_circle(img, cx, cy, r, color):
    """純色填充圓形"""
    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            if (x - cx)**2 + (y - cy)**2 <= r**2:
                px(img, x, y, color)

def draw_eye_shaded(img, x, y, iris_color):
    """
    帶陰影的 3x3 眼睛
    外圈白色、虹膜、瞳孔、高光
    """
    WHITE = (255, 255, 255, 255)
    BLACK = (20, 10, 5, 255)
    HIGHLIGHT = (255, 255, 255, 255)

    # 白色眼白（3x3）
    for dy in range(3):
        for dx in range(3):
            px(img, x + dx, y + dy, WHITE)

    # 虹膜（2x2，中央偏下）
    iris_dark = shade_color(iris_color, 0.7)
    px(img, x + 1, y + 1, iris_color)
    px(img, x + 2, y + 1, iris_dark)
    px(img, x + 1, y + 2, iris_dark)
    px(img, x + 2, y + 2, iris_color)

    # 瞳孔（1x1）
    px(img, x + 1, y + 1, BLACK)

    # 高光（左上）
    px(img, x, y, HIGHLIGHT)

def draw_smile_v4(img, x, y, w, color):
    """改良版微笑：兩端高，中間低，有陰影"""
    dark = shade_color(color, 0.6)
    # 嘴角（高）
    px(img, x,     y,   color)
    px(img, x+w-1, y,   color)
    # 嘴唇（低）
    for i in range(1, w-1):
        px(img, x+i, y+1, color)
    # 嘴唇下陰影
    for i in range(1, w-1):
        px(img, x+i, y+2, dark)

def gen_chiikawa_v4(state="idle"):
    """
    吉伊卡哇 v4 - 32x32
    改進：3色陰影、帶陰影眼睛、改良嘴巴
    """
    img = new_img(32, 32)

    SKIN_BASE = (255, 220, 195)
    OUTLINE   = (55, 25, 8, 255)
    EYE_IRIS  = (80, 50, 30)
    BLUSH     = (255, 150, 140, 255)
    EAR_IN    = (255, 185, 165)
    PINK_BASE = (255, 130, 185)
    GLOW      = (255, 210, 235, 255)
    WHITE     = (255, 255, 255, 255)

    dy = -2 if state == "attack" else (2 if state == "bigwin" else 0)

    # 耳朵（帶陰影）
    fill_circle_shaded(img, 9,  8+dy, 4, EAR_IN)
    fill_circle_shaded(img, 23, 8+dy, 4, EAR_IN)
    draw_outline_circle(img, 9,  8+dy, 4, OUTLINE)
    draw_outline_circle(img, 23, 8+dy, 4, OUTLINE)

    # 大圓頭（帶陰影）
    fill_circle_shaded(img, 16, 15+dy, 10, SKIN_BASE)
    draw_outline_circle(img, 16, 15+dy, 10, OUTLINE)

    # 眼睛
    if state == "bigwin":
        draw_eye_shaded(img, 9, 12+dy, EYE_IRIS)
        draw_eye_shaded(img, 19, 12+dy, EYE_IRIS)
    else:
        # 閉眼弧線（努力感）
        for i in range(5):
            px(img, 9+i, 14+dy, OUTLINE)
            px(img, 19+i, 14+dy, OUTLINE)
        # 睫毛
        px(img, 9,  13+dy, OUTLINE)
        px(img, 13, 13+dy, OUTLINE)
        px(img, 19, 13+dy, OUTLINE)
        px(img, 23, 13+dy, OUTLINE)

    # 腮紅（半透明粉紅）
    for dy2 in range(-1, 2):
        for dx2 in range(-2, 3):
            if dx2*dx2 + dy2*dy2 <= 4:
                px(img, 9+dx2,  17+dy+dy2, BLUSH)
                px(img, 23+dx2, 17+dy+dy2, BLUSH)

    # 嘴巴（改良版）
    draw_smile_v4(img, 13, 18+dy, 6, OUTLINE)

    # 小身體（帶陰影）
    fill_circle_shaded(img, 16, 26+dy, 5, SKIN_BASE)
    draw_outline_circle(img, 16, 26+dy, 5, OUTLINE)

    # 討伐棒
    if state == "attack":
        for i in range(7):
            px(img, 22+i, 10+dy-i, PINK_BASE)
            px(img, 22+i, 11+dy-i, shade_color(PINK_BASE, 0.8))
        fill_circle(img, 29, 4+dy, 2, GLOW)
        # 劍氣殘影
        for i in range(3):
            px(img, 20+i, 14+dy-i, (*PINK_BASE[:3], 120))
    else:
        for i in range(5):
            px(img, 22+i, 13+dy-i, PINK_BASE)
        fill_circle(img, 27, 9+dy, 2, GLOW)

    return img

def gen_hachiware_v4(state="idle"):
    """小八 v4 - 32x32，帶陰影"""
    img = new_img(32, 32)

    WHITE_BASE = (240, 240, 240)
    OUTLINE    = (35, 45, 75, 255)
    STRIPE     = (90, 130, 210)
    EYE_IRIS   = (50, 90, 190)
    BLUE_BASE  = (110, 160, 250)
    GLOW       = (185, 210, 255, 255)
    EAR_IN     = (195, 215, 255)

    dy = -2 if state == "attack" else (2 if state == "bigwin" else 0)

    # 尖耳朵（三角形帶陰影）
    for i in range(7):
        for j in range(i+1):
            bri = 1.0 - j * 0.05
            c = shade_color(WHITE_BASE, bri)
            px(img, 8+j,  5+dy+i, c)
            px(img, 24-j, 5+dy+i, c)
    # 耳朵輪廓
    for i in range(7):
        px(img, 8,    5+dy+i, OUTLINE)
        px(img, 8+i,  5+dy+i, OUTLINE)
        px(img, 24,   5+dy+i, OUTLINE)
        px(img, 24-i, 5+dy+i, OUTLINE)
    # 耳朵內側
    for i in range(4):
        for j in range(i+1):
            px(img, 9+j,  7+dy+i, EAR_IN)
            px(img, 23-j, 7+dy+i, EAR_IN)

    # 大圓頭（帶陰影）
    fill_circle_shaded(img, 16, 16+dy, 10, WHITE_BASE)
    draw_outline_circle(img, 16, 16+dy, 10, OUTLINE)

    # 藍色條紋（帶陰影）
    for stripe_x in [10, 14, 18, 22]:
        for y in range(10+dy, 22+dy):
            if abs(y - (16+dy)) < 9:
                # 條紋也有陰影
                bri = 1.0 if stripe_x < 16 else 0.85
                px(img, stripe_x, y, shade_color(STRIPE, bri))

    # 眼睛（睜開）
    draw_eye_shaded(img, 10, 13+dy, EYE_IRIS)
    draw_eye_shaded(img, 20, 13+dy, EYE_IRIS)

    # 嘴巴
    draw_smile_v4(img, 13, 19+dy, 6, OUTLINE)

    # 小身體
    fill_circle_shaded(img, 16, 26+dy, 5, WHITE_BASE)
    draw_outline_circle(img, 16, 26+dy, 5, OUTLINE)

    # 討伐棒
    if state == "attack":
        for i in range(7):
            px(img, 22+i, 10+dy-i, BLUE_BASE)
            px(img, 22+i, 11+dy-i, shade_color(BLUE_BASE, 0.75))
        fill_circle(img, 29, 4+dy, 2, GLOW)
    else:
        for i in range(5):
            px(img, 22+i, 13+dy-i, BLUE_BASE)
        fill_circle(img, 27, 9+dy, 2, GLOW)

    return img

def gen_usagi_v4(state="idle"):
    """烏薩奇 v4 - 32x32，帶陰影"""
    img = new_img(32, 32)

    WHITE_BASE = (245, 245, 245)
    OUTLINE    = (65, 45, 65, 255)
    PINK_IN    = (255, 170, 170)
    EYE_IRIS   = (200, 40, 40)
    YELLOW_BASE= (255, 215, 35)
    GLOW       = (255, 245, 145, 255)

    dy = -3 if state == "attack" else (3 if state == "bigwin" else 0)

    # 長耳朵（帶陰影）
    for y in range(1+dy, 13+dy):
        for x in range(7, 14):
            if abs(x - 10) <= 3:
                bri = 1.1 if x < 10 else 0.85  # 左亮右暗
                px(img, x, y, shade_color(WHITE_BASE, bri))
        for x in range(19, 26):
            if abs(x - 22) <= 3:
                bri = 1.1 if x < 22 else 0.85
                px(img, x, y, shade_color(WHITE_BASE, bri))
    # 耳朵輪廓
    for y in range(1+dy, 13+dy):
        px(img, 7,  y, OUTLINE)
        px(img, 13, y, OUTLINE)
        px(img, 19, y, OUTLINE)
        px(img, 25, y, OUTLINE)
    # 耳朵內側（粉紅，帶陰影）
    for y in range(2+dy, 12+dy):
        for x in range(8, 13):
            if abs(x - 10) <= 2:
                bri = 1.0 if x <= 10 else 0.8
                px(img, x, y, shade_color(PINK_IN, bri))
        for x in range(20, 25):
            if abs(x - 22) <= 2:
                bri = 1.0 if x <= 22 else 0.8
                px(img, x, y, shade_color(PINK_IN, bri))

    # 大圓頭（帶陰影）
    fill_circle_shaded(img, 16, 17+dy, 10, WHITE_BASE)
    draw_outline_circle(img, 16, 17+dy, 10, OUTLINE)

    # 眼睛（紅色，興奮）
    draw_eye_shaded(img, 10, 14+dy, EYE_IRIS)
    draw_eye_shaded(img, 20, 14+dy, EYE_IRIS)

    # 嘴巴（興奮大笑）
    draw_smile_v4(img, 12, 20+dy, 8, OUTLINE)
    # 牙齒
    for tx in [14, 15, 17, 18]:
        px(img, tx, 21+dy, (245, 245, 245, 255))

    # 小身體
    fill_circle_shaded(img, 16, 27+dy, 4, WHITE_BASE)
    draw_outline_circle(img, 16, 27+dy, 4, OUTLINE)

    # 討伐棒（旋轉感）
    if state == "attack":
        for i in range(7):
            px(img, 21+i, 10+dy-i, YELLOW_BASE)
            px(img, 22+i, 10+dy-i, shade_color(YELLOW_BASE, 0.75))
        # 旋轉殘影（半透明）
        for i in range(4):
            c = (*YELLOW_BASE[:3], 100)
            px(img, 23+i, 12+dy-i, c)
        fill_circle(img, 29, 4+dy, 3, GLOW)
    else:
        for i in range(5):
            px(img, 22+i, 13+dy-i, YELLOW_BASE)
        fill_circle(img, 27, 9+dy, 2, GLOW)

    return img

def generate_characters_v4():
    chars_dir = os.path.join(OUTPUT_BASE, "characters")
    os.makedirs(chars_dir, exist_ok=True)

    print("[Characters v4 - 3-color shading + Pillow Shading]")
    fns = [
        ("chiikawa", gen_chiikawa_v4),
        ("hachiware", gen_hachiware_v4),
        ("usagi", gen_usagi_v4),
    ]
    for name, fn in fns:
        for state in ["idle", "attack", "bigwin"]:
            img = fn(state)
            img = img.resize((64, 64), Image.NEAREST)
            path = os.path.join(chars_dir, f"{name}_{state}.png")
            img.save(path)
            print(f"  OK {name}_{state}.png (64x64)")

if __name__ == "__main__":
    print("Pixel Art v4 - Shading upgrade...")
    generate_characters_v4()
    print("\nOK Done!")
