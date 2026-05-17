# -*- coding: utf-8 -*-
"""
像素美術生成器 v3
基於 pixnote.net 的像素角色設計原則：
- 32x32 chibi 比例（頭大身小）
- 2x2 眼睛 + 1px 高光
- 圓形輪廓 = 可愛感
- 靠剪影就能辨識角色
"""
from PIL import Image, ImageDraw
import os

OUTPUT_BASE = r"D:\Kiro\client\chiikawa-pixel\assets\sprites"

def new_img(w, h):
    return Image.new("RGBA", (w, h), (0, 0, 0, 0))

def px(img, x, y, c):
    if 0 <= x < img.width and 0 <= y < img.height:
        img.putpixel((x, y), c)

def fill_circle(img, cx, cy, r, color):
    """填充圓形（像素精確）"""
    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            if (x - cx)**2 + (y - cy)**2 <= r**2:
                px(img, x, y, color)

def draw_outline_circle(img, cx, cy, r, color):
    """畫圓形輪廓"""
    for y in range(cy - r - 1, cy + r + 2):
        for x in range(cx - r - 1, cx + r + 2):
            dist = ((x - cx)**2 + (y - cy)**2) ** 0.5
            if r - 0.5 <= dist <= r + 1.2:
                px(img, x, y, color)

def draw_eye_2x2(img, x, y, pupil_color, highlight=True):
    """2x2 眼睛 + 高光（pixnote.net 建議）"""
    px(img, x,   y,   pupil_color)
    px(img, x+1, y,   pupil_color)
    px(img, x,   y+1, pupil_color)
    px(img, x+1, y+1, pupil_color)
    if highlight:
        px(img, x, y, (255, 255, 255, 255))  # 左上高光

def draw_smile(img, x, y, w, color):
    """像素微笑 - 向下弧線"""
    # 兩端高，中間低
    px(img, x,       y,   color)
    px(img, x+w-1,   y,   color)
    for i in range(1, w-1):
        px(img, x+i, y+1, color)
    # 嘴角加深
    px(img, x,     y+1, color)
    px(img, x+w-1, y+1, color)

def gen_chiikawa_v3(state="idle"):
    """
    吉伊卡哇 32x32
    特徵：圓頭、小耳朵、閉眼（努力感）、腮紅、粉紅劍氣
    """
    img = new_img(32, 32)

    # 顏色
    SKIN    = (255, 225, 200, 255)
    OUTLINE = (60,  30,  10,  255)
    EYE     = (40,  20,  10,  255)
    BLUSH   = (255, 160, 150, 255)
    EAR_IN  = (255, 190, 170, 255)
    PINK    = (255, 140, 190, 255)
    GLOW    = (255, 210, 230, 255)
    WHITE   = (255, 255, 255, 255)

    dy = -2 if state == "attack" else (2 if state == "bigwin" else 0)

    # 耳朵（小圓）
    fill_circle(img, 9,  8+dy, 4, SKIN)
    fill_circle(img, 23, 8+dy, 4, SKIN)
    fill_circle(img, 9,  8+dy, 2, EAR_IN)
    fill_circle(img, 23, 8+dy, 2, EAR_IN)
    draw_outline_circle(img, 9,  8+dy, 4, OUTLINE)
    draw_outline_circle(img, 23, 8+dy, 4, OUTLINE)

    # 大圓頭（chibi 比例：頭佔 60%）
    fill_circle(img, 16, 15+dy, 10, SKIN)
    draw_outline_circle(img, 16, 15+dy, 10, OUTLINE)

    # 眼睛
    if state == "bigwin":
        # 大獎：睜大眼睛 + 星星
        draw_eye_2x2(img, 10, 13+dy, EYE)
        draw_eye_2x2(img, 20, 13+dy, EYE)
        # 星星高光
        px(img, 10, 13+dy, WHITE)
        px(img, 20, 13+dy, WHITE)
    else:
        # 一般/攻擊：閉眼弧線（努力感）
        for i in range(4):
            px(img, 10+i, 14+dy, EYE)
        for i in range(4):
            px(img, 19+i, 14+dy, EYE)
        # 眼睫毛
        px(img, 10, 13+dy, EYE)
        px(img, 13, 13+dy, EYE)
        px(img, 19, 13+dy, EYE)
        px(img, 22, 13+dy, EYE)

    # 腮紅（2x2 圓點）
    fill_circle(img, 9,  17+dy, 2, BLUSH)
    fill_circle(img, 23, 17+dy, 2, BLUSH)

    # 嘴巴（小，在臉中央偏下）
    px(img, 14, 19+dy, OUTLINE)
    px(img, 15, 20+dy, OUTLINE)
    px(img, 16, 20+dy, OUTLINE)
    px(img, 17, 19+dy, OUTLINE)

    # 小身體
    fill_circle(img, 16, 26+dy, 5, SKIN)
    draw_outline_circle(img, 16, 26+dy, 5, OUTLINE)

    # 討伐棒
    if state == "attack":
        # 揮棒：斜向右上
        for i in range(7):
            px(img, 22+i, 10+dy-i, PINK)
            px(img, 22+i, 11+dy-i, PINK)
        # 劍氣光點
        fill_circle(img, 29, 4+dy, 2, GLOW)
    else:
        # 持棒
        for i in range(5):
            px(img, 22+i, 13+dy-i, PINK)
        fill_circle(img, 27, 9+dy, 2, GLOW)

    return img

def gen_hachiware_v3(state="idle"):
    """
    小八 32x32
    特徵：尖耳貓、藍白條紋、精神眼睛（睜開）、藍色劍氣
    """
    img = new_img(32, 32)

    WHITE   = (245, 245, 245, 255)
    OUTLINE = (40,  50,  80,  255)
    STRIPE  = (100, 140, 220, 255)
    EYE_B   = (60,  100, 200, 255)
    EYE_W   = (255, 255, 255, 255)
    BLUE    = (120, 170, 255, 255)
    GLOW    = (190, 215, 255, 255)
    EAR_IN  = (200, 220, 255, 255)

    dy = -2 if state == "attack" else (2 if state == "bigwin" else 0)

    # 尖耳朵（三角形）
    for i in range(6):
        for j in range(i+1):
            px(img, 8+j,  6+dy+i, WHITE)
            px(img, 24-j, 6+dy+i, WHITE)
    # 耳朵輪廓
    for i in range(6):
        px(img, 8,    6+dy+i, OUTLINE)
        px(img, 8+i,  6+dy+i, OUTLINE)
        px(img, 24,   6+dy+i, OUTLINE)
        px(img, 24-i, 6+dy+i, OUTLINE)
    # 耳朵內側
    for i in range(3):
        for j in range(i+1):
            px(img, 9+j,  8+dy+i, EAR_IN)
            px(img, 23-j, 8+dy+i, EAR_IN)

    # 大圓頭
    fill_circle(img, 16, 16+dy, 10, WHITE)
    draw_outline_circle(img, 16, 16+dy, 10, OUTLINE)

    # 藍色條紋（特徵）
    for stripe_x in [10, 14, 18, 22]:
        for y in range(10+dy, 22+dy):
            if abs(y - (16+dy)) < 9:
                px(img, stripe_x, y, STRIPE)

    # 眼睛（睜開，精神感）
    draw_eye_2x2(img, 10, 13+dy, EYE_B)
    draw_eye_2x2(img, 20, 13+dy, EYE_B)
    # 高光
    px(img, 10, 13+dy, EYE_W)
    px(img, 20, 13+dy, EYE_W)

    # 嘴巴（微笑）
    draw_smile(img, 13, 19+dy, 6, OUTLINE)

    # 小身體
    fill_circle(img, 16, 26+dy, 5, WHITE)
    draw_outline_circle(img, 16, 26+dy, 5, OUTLINE)

    # 討伐棒（藍色）
    if state == "attack":
        for i in range(7):
            px(img, 22+i, 10+dy-i, BLUE)
            px(img, 22+i, 11+dy-i, BLUE)
        fill_circle(img, 29, 4+dy, 2, GLOW)
    else:
        for i in range(5):
            px(img, 22+i, 13+dy-i, BLUE)
        fill_circle(img, 27, 9+dy, 2, GLOW)

    return img

def gen_usagi_v3(state="idle"):
    """
    烏薩奇 32x32
    特徵：長耳兔、紅眼、興奮感、黃色旋轉劍氣
    """
    img = new_img(32, 32)

    WHITE   = (248, 248, 248, 255)
    OUTLINE = (70,  50,  70,  255)
    PINK_IN = (255, 180, 180, 255)
    EYE_R   = (210, 50,  50,  255)
    EYE_W   = (255, 200, 200, 255)
    YELLOW  = (255, 220, 40,  255)
    GLOW    = (255, 245, 150, 255)

    dy = -3 if state == "attack" else (3 if state == "bigwin" else 0)

    # 長耳朵（橢圓形，高）
    fill_circle(img, 10, 6+dy, 3, WHITE)
    fill_circle(img, 22, 6+dy, 3, WHITE)
    for y in range(2+dy, 12+dy):
        for x in range(7, 14):
            if abs(x - 10) <= 3:
                px(img, x, y, WHITE)
        for x in range(19, 26):
            if abs(x - 22) <= 3:
                px(img, x, y, WHITE)
    # 耳朵輪廓
    for y in range(2+dy, 12+dy):
        px(img, 7,  y, OUTLINE)
        px(img, 13, y, OUTLINE)
        px(img, 19, y, OUTLINE)
        px(img, 25, y, OUTLINE)
    # 耳朵內側（粉紅）
    for y in range(3+dy, 11+dy):
        for x in range(8, 13):
            if abs(x - 10) <= 2:
                px(img, x, y, PINK_IN)
        for x in range(20, 25):
            if abs(x - 22) <= 2:
                px(img, x, y, PINK_IN)

    # 大圓頭
    fill_circle(img, 16, 17+dy, 10, WHITE)
    draw_outline_circle(img, 16, 17+dy, 10, OUTLINE)

    # 眼睛（紅色，興奮）
    draw_eye_2x2(img, 10, 14+dy, EYE_R)
    draw_eye_2x2(img, 20, 14+dy, EYE_R)
    px(img, 10, 14+dy, EYE_W)
    px(img, 20, 14+dy, EYE_W)

    # 嘴巴（興奮大笑）
    draw_smile(img, 12, 20+dy, 8, OUTLINE)
    # 牙齒
    px(img, 14, 21+dy, WHITE)
    px(img, 15, 21+dy, WHITE)
    px(img, 17, 21+dy, WHITE)
    px(img, 18, 21+dy, WHITE)

    # 小身體
    fill_circle(img, 16, 27+dy, 4, WHITE)
    draw_outline_circle(img, 16, 27+dy, 4, OUTLINE)

    # 討伐棒（黃色，旋轉感）
    if state == "attack":
        for i in range(7):
            px(img, 21+i, 10+dy-i, YELLOW)
            px(img, 22+i, 10+dy-i, YELLOW)
        # 旋轉殘影
        for i in range(4):
            px(img, 23+i, 12+dy-i, (255, 220, 40, 150))
        fill_circle(img, 29, 4+dy, 3, GLOW)
    else:
        for i in range(5):
            px(img, 22+i, 13+dy-i, YELLOW)
        fill_circle(img, 27, 9+dy, 2, GLOW)

    return img

def gen_background():
    """生成海底背景 1280x720"""
    img = Image.new("RGBA", (1280, 720), (0, 0, 0, 255))
    draw = ImageDraw.Draw(img)

    # 海底漸層
    for y in range(720):
        r = int(10 + (y/720) * 5)
        g = int(40 + (y/720) * 20)
        b = int(100 + (y/720) * 30)
        draw.line([(0, y), (1280, y)], fill=(r, g, b, 255))

    # 珊瑚
    import random
    random.seed(42)
    coral_colors = [(200, 80, 80, 255), (220, 120, 60, 255), (180, 60, 120, 255)]
    for _ in range(12):
        x = random.randint(50, 1200)
        h = random.randint(40, 90)
        c = random.choice(coral_colors)
        draw.rectangle([x, 720-h, x+6, 720], fill=c)
        draw.ellipse([x-5, 720-h-6, x+11, 720-h+6], fill=c)
        draw.ellipse([x-2, 720-h-10, x+8, 720-h], fill=c)

    # 海草
    import math
    for _ in range(25):
        x = random.randint(0, 1260)
        h = random.randint(30, 70)
        for i in range(h):
            ox = int(math.sin(i * 0.3) * 8)
            draw.rectangle([x+ox, 720-i-1, x+ox+4, 720-i], fill=(40, 120, 60, 255))

    # 氣泡
    for _ in range(40):
        bx = random.randint(10, 1260)
        by = random.randint(50, 650)
        br = random.randint(2, 6)
        draw.ellipse([bx-br, by-br, bx+br, by+br],
                     outline=(150, 200, 255, 180), width=1)

    # 沙底
    draw.rectangle([0, 700, 1280, 720], fill=(180, 160, 100, 255))
    for i in range(0, 1280, 20):
        draw.ellipse([i, 698, i+15, 706], fill=(200, 180, 120, 255))

    return img

def gen_boss_bg():
    """BOSS 背景"""
    img = gen_background()
    draw = ImageDraw.Draw(img)
    # 暗化
    dark = Image.new("RGBA", (1280, 720), (0, 0, 0, 160))
    img = Image.alpha_composite(img, dark)
    draw = ImageDraw.Draw(img)
    # 紅色邊框
    for t in range(6):
        draw.rectangle([t, t, 1279-t, 719-t], outline=(200, 30, 30, 255-t*20))
    return img

def gen_bonus_bg():
    """Bonus 草地背景"""
    img = Image.new("RGBA", (1280, 720), (0, 0, 0, 255))
    draw = ImageDraw.Draw(img)
    # 天空
    for y in range(450):
        r = int(100 + (y/450) * 50)
        g = int(180 + (y/450) * 30)
        b = int(240 - (y/450) * 30)
        draw.line([(0, y), (1280, y)], fill=(r, g, b, 255))
    # 草地
    for y in range(450, 720):
        ratio = (y-450) / 270
        r = int(40 - ratio * 10)
        g = int(140 - ratio * 40)
        b = int(40 - ratio * 10)
        draw.line([(0, y), (1280, y)], fill=(r, g, b, 255))
    # 草叢
    import random, math
    random.seed(99)
    for _ in range(50):
        x = random.randint(0, 1260)
        h = random.randint(20, 50)
        for i in range(3):
            draw.polygon([
                (x+i*8, 450),
                (x+i*8+4, 450-h-random.randint(0,15)),
                (x+i*8+8, 450)
            ], fill=(60, 160, 50, 255))
    return img

def generate_all_v3():
    chars_dir = os.path.join(OUTPUT_BASE, "characters")
    targets_dir = os.path.join(OUTPUT_BASE, "targets")
    bg_dir = os.path.join(OUTPUT_BASE, "backgrounds")
    os.makedirs(chars_dir, exist_ok=True)
    os.makedirs(targets_dir, exist_ok=True)
    os.makedirs(bg_dir, exist_ok=True)

    print("[角色 v3 - chibi 比例 + 2x2 眼睛]")
    char_fns = [
        ("chiikawa", gen_chiikawa_v3),
        ("hachiware", gen_hachiware_v3),
        ("usagi", gen_usagi_v3),
    ]
    for name, fn in char_fns:
        for state in ["idle", "attack", "bigwin"]:
            img = fn(state)
            # 2x 放大保持像素感
            img = img.resize((64, 64), Image.NEAREST)
            path = os.path.join(chars_dir, f"{name}_{state}.png")
            img.save(path)
            print(f"  OK {name}_{state}.png")

    print("\n[背景 v3]")
    sea = gen_background()
    sea.save(os.path.join(bg_dir, "sea_bg.png"))
    print("  OK sea_bg.png")

    boss_bg = gen_boss_bg()
    boss_bg.save(os.path.join(bg_dir, "boss_bg.png"))
    print("  OK boss_bg.png")

    bonus_bg = gen_bonus_bg()
    bonus_bg.save(os.path.join(bg_dir, "bonus_bg.png"))
    print("  OK bonus_bg.png")

if __name__ == "__main__":
    print("像素美術 v3 生成中...")
    generate_all_v3()
    print("\nOK 完成！")
