# -*- coding: utf-8 -*-
"""
生成完整角色（頭+身體+手腳）
基於官方吉伊卡哇設計：大圓頭 + 小橢圓身體 + 小手小腳
畫布：48x48，放大到 96x96
"""
from PIL import Image, ImageDraw
import os
import math

OUTPUT_BASE = r"D:\Kiro\client\chiikawa-pixel\assets\sprites"

def new_img(w, h):
    return Image.new("RGBA", (w, h), (0, 0, 0, 0))

def px(img, x, y, c):
    if 0 <= x < img.width and 0 <= y < img.height:
        img.putpixel((x, y), c)

def fill_ellipse(img, cx, cy, rx, ry, color, outline=None):
    """填充橢圓"""
    draw = ImageDraw.Draw(img)
    if outline:
        draw.ellipse([cx-rx-1, cy-ry-1, cx+rx+1, cy+ry+1], fill=outline)
    draw.ellipse([cx-rx, cy-ry, cx+rx, cy+ry], fill=color)

def fill_circle(img, cx, cy, r, color, outline=None):
    fill_ellipse(img, cx, cy, r, r, color, outline)

def gen_chiikawa_full(state="idle"):
    """
    吉伊卡哇完整角色 48x48
    - 大圓頭（半徑 14）
    - 小橢圓身體（10x8）
    - 小手（圓形 3）
    - 小腳（橢圓 4x3）
    - 圓耳朵（半徑 5）
    """
    img = new_img(48, 48)
    draw = ImageDraw.Draw(img)

    # 官方顏色
    WHITE  = (255, 255, 247, 255)
    OUTLINE= (41,  42,  43,  255)
    BLUSH  = (239, 165, 201, 255)
    EAR_IN = (255, 200, 195, 255)
    PINK   = (255, 140, 190, 255)
    GLOW   = (255, 215, 235, 255)
    BLACK  = (20,  15,  10,  255)

    # 動作偏移
    dy = -3 if state == "attack" else (3 if state == "bigwin" else 0)
    arm_angle = -30 if state == "attack" else 0

    # === 腳（在身體下方）===
    # 左腳
    fill_ellipse(img, 18, 40+dy, 5, 4, WHITE, OUTLINE)
    # 右腳
    fill_ellipse(img, 30, 40+dy, 5, 4, WHITE, OUTLINE)

    # === 身體（橢圓）===
    fill_ellipse(img, 24, 33+dy, 10, 8, WHITE, OUTLINE)

    # === 手臂 ===
    if state == "attack":
        # 攻擊：右手舉起揮棒
        fill_circle(img, 36, 26+dy, 4, WHITE, OUTLINE)  # 右手（高）
        fill_circle(img, 12, 32+dy, 4, WHITE, OUTLINE)  # 左手
        # 討伐棒
        draw.line([(36, 26+dy), (44, 18+dy)], fill=PINK, width=3)
        fill_circle(img, 44, 17+dy, 3, GLOW)
    else:
        # 待機：手在兩側
        fill_circle(img, 13, 33+dy, 4, WHITE, OUTLINE)  # 左手
        fill_circle(img, 35, 33+dy, 4, WHITE, OUTLINE)  # 右手
        if state != "bigwin":
            # 持棒
            draw.line([(35, 33+dy), (41, 27+dy)], fill=PINK, width=2)
            fill_circle(img, 41, 26+dy, 2, GLOW)

    # === 耳朵（在頭部後面）===
    fill_circle(img, 13, 12+dy, 5, WHITE, OUTLINE)
    fill_circle(img, 35, 12+dy, 5, WHITE, OUTLINE)
    fill_circle(img, 13, 12+dy, 3, EAR_IN)
    fill_circle(img, 35, 12+dy, 3, EAR_IN)

    # === 大圓頭 ===
    fill_circle(img, 24, 18+dy, 14, WHITE, OUTLINE)

    # === 眼睛 ===
    if state == "bigwin":
        # 大獎：睜大眼睛
        fill_circle(img, 18, 16+dy, 3, BLACK)
        fill_circle(img, 30, 16+dy, 3, BLACK)
        px(img, 17, 15+dy, (255, 255, 255, 255))  # 高光
        px(img, 29, 15+dy, (255, 255, 255, 255))
    else:
        # 閉眼弧線
        for i in range(5):
            px(img, 15+i, 17+dy, OUTLINE)
            px(img, 27+i, 17+dy, OUTLINE)
        # 睫毛
        px(img, 15, 16+dy, OUTLINE)
        px(img, 19, 16+dy, OUTLINE)
        px(img, 27, 16+dy, OUTLINE)
        px(img, 31, 16+dy, OUTLINE)

    # === 腮紅（小圓點）===
    fill_circle(img, 13, 20+dy, 3, BLUSH)
    fill_circle(img, 35, 20+dy, 3, BLUSH)

    # === 嘴巴（小 V 形）===
    px(img, 22, 22+dy, OUTLINE)
    px(img, 23, 23+dy, OUTLINE)
    px(img, 24, 23+dy, OUTLINE)
    px(img, 25, 22+dy, OUTLINE)

    # 大獎：跳起星星
    if state == "bigwin":
        for sx, sy in [(6, 6), (42, 4), (4, 28)]:
            px(img, sx, sy+dy, (255, 230, 50, 230))
            for dx, dy2 in [(1,0),(-1,0),(0,1),(0,-1)]:
                nx, ny = sx+dx, sy+dy2+dy
                if 0 <= nx < 48 and 0 <= ny < 48:
                    px(img, nx, ny, (255, 230, 50, 150))

    return img

def gen_hachiware_full(state="idle"):
    """
    小八完整角色 48x48
    - 白色帶藍色條紋
    - 尖耳朵
    """
    img = new_img(48, 48)
    draw = ImageDraw.Draw(img)

    WHITE  = (255, 255, 247, 255)
    OUTLINE= (41,  42,  43,  255)
    STRIPE = (51,  112, 192, 255)
    EAR_IN = (190, 210, 250, 255)
    BLUE   = (100, 150, 240, 255)
    GLOW   = (180, 210, 255, 255)

    dy = -3 if state == "attack" else (3 if state == "bigwin" else 0)

    # 腳
    fill_ellipse(img, 18, 40+dy, 5, 4, WHITE, OUTLINE)
    fill_ellipse(img, 30, 40+dy, 5, 4, WHITE, OUTLINE)

    # 身體
    fill_ellipse(img, 24, 33+dy, 10, 8, WHITE, OUTLINE)
    # 條紋
    for sx in [18, 22, 26, 30]:
        for sy in range(26+dy, 40+dy):
            if (sy - (26+dy)) < 14:
                px(img, sx, sy, STRIPE)

    # 手臂
    if state == "attack":
        fill_circle(img, 36, 26+dy, 4, WHITE, OUTLINE)
        fill_circle(img, 12, 32+dy, 4, WHITE, OUTLINE)
        draw.line([(36, 26+dy), (44, 18+dy)], fill=BLUE, width=3)
        fill_circle(img, 44, 17+dy, 3, GLOW)
    else:
        fill_circle(img, 13, 33+dy, 4, WHITE, OUTLINE)
        fill_circle(img, 35, 33+dy, 4, WHITE, OUTLINE)
        draw.line([(35, 33+dy), (41, 27+dy)], fill=BLUE, width=2)
        fill_circle(img, 41, 26+dy, 2, GLOW)

    # 尖耳朵（三角形）
    draw.polygon([(10, 14+dy), (14, 4+dy), (18, 14+dy)], fill=WHITE, outline=OUTLINE)
    draw.polygon([(30, 14+dy), (34, 4+dy), (38, 14+dy)], fill=WHITE, outline=OUTLINE)
    draw.polygon([(11, 13+dy), (14, 6+dy), (17, 13+dy)], fill=EAR_IN)
    draw.polygon([(31, 13+dy), (34, 6+dy), (37, 13+dy)], fill=EAR_IN)

    # 大圓頭
    fill_circle(img, 24, 18+dy, 14, WHITE, OUTLINE)

    # 條紋（頭部）
    for sx in [16, 20, 24, 28, 32]:
        for sy in range(8+dy, 30+dy):
            dist = math.sqrt((sx-24)**2 + (sy-(18+dy))**2)
            if dist < 13:
                px(img, sx, sy, STRIPE)

    # 眼睛（睜開）
    fill_circle(img, 18, 16+dy, 3, (50, 90, 190, 255))
    fill_circle(img, 30, 16+dy, 3, (50, 90, 190, 255))
    px(img, 17, 15+dy, (255, 255, 255, 255))
    px(img, 29, 15+dy, (255, 255, 255, 255))

    # 嘴巴
    px(img, 22, 22+dy, OUTLINE)
    px(img, 23, 23+dy, OUTLINE)
    px(img, 24, 23+dy, OUTLINE)
    px(img, 25, 22+dy, OUTLINE)

    return img

def gen_usagi_full(state="idle"):
    """
    烏薩奇完整角色 48x48
    - 白色，長耳朵
    - 紅眼
    """
    img = new_img(48, 48)
    draw = ImageDraw.Draw(img)

    WHITE  = (255, 255, 247, 255)
    OUTLINE= (17,  17,  17,  255)
    PINK_IN= (255, 175, 175, 255)
    EYE_R  = (255, 91,  86,  255)
    YELLOW = (255, 215, 30,  255)
    GLOW   = (255, 245, 140, 255)

    dy = -3 if state == "attack" else (3 if state == "bigwin" else 0)

    # 腳
    fill_ellipse(img, 18, 40+dy, 5, 4, WHITE, OUTLINE)
    fill_ellipse(img, 30, 40+dy, 5, 4, WHITE, OUTLINE)

    # 身體
    fill_ellipse(img, 24, 33+dy, 10, 8, WHITE, OUTLINE)

    # 手臂
    if state == "attack":
        fill_circle(img, 36, 26+dy, 4, WHITE, OUTLINE)
        fill_circle(img, 12, 32+dy, 4, WHITE, OUTLINE)
        draw.line([(36, 26+dy), (44, 18+dy)], fill=YELLOW, width=3)
        draw.line([(37, 27+dy), (44, 20+dy)], fill=YELLOW, width=2)
        fill_circle(img, 44, 17+dy, 3, GLOW)
    else:
        fill_circle(img, 13, 33+dy, 4, WHITE, OUTLINE)
        fill_circle(img, 35, 33+dy, 4, WHITE, OUTLINE)
        draw.line([(35, 33+dy), (41, 27+dy)], fill=YELLOW, width=2)
        fill_circle(img, 41, 26+dy, 2, GLOW)

    # 長耳朵
    draw.ellipse([8, 0+dy, 18, 16+dy], fill=WHITE, outline=OUTLINE)
    draw.ellipse([30, 0+dy, 40, 16+dy], fill=WHITE, outline=OUTLINE)
    draw.ellipse([10, 1+dy, 16, 14+dy], fill=PINK_IN)
    draw.ellipse([32, 1+dy, 38, 14+dy], fill=PINK_IN)

    # 大圓頭
    fill_circle(img, 24, 20+dy, 14, WHITE, OUTLINE)

    # 眼睛（紅色）
    fill_circle(img, 18, 18+dy, 3, EYE_R)
    fill_circle(img, 30, 18+dy, 3, EYE_R)
    px(img, 17, 17+dy, (255, 200, 200, 255))
    px(img, 29, 17+dy, (255, 200, 200, 255))

    # 嘴巴（興奮）
    px(img, 21, 24+dy, OUTLINE)
    px(img, 22, 25+dy, OUTLINE)
    px(img, 23, 25+dy, OUTLINE)
    px(img, 24, 25+dy, OUTLINE)
    px(img, 25, 25+dy, OUTLINE)
    px(img, 26, 24+dy, OUTLINE)
    # 牙齒
    px(img, 22, 25+dy, (245, 245, 245, 255))
    px(img, 24, 25+dy, (245, 245, 245, 255))

    return img

def generate_all():
    chars_dir = os.path.join(OUTPUT_BASE, "characters")
    os.makedirs(chars_dir, exist_ok=True)

    print("Generating full characters (head + body + arms + legs)...")
    fns = [
        ("chiikawa", gen_chiikawa_full),
        ("hachiware", gen_hachiware_full),
        ("usagi", gen_usagi_full),
    ]

    for name, fn in fns:
        for state in ["idle", "attack", "bigwin"]:
            img = fn(state)
            # 放大 2x 到 96x96
            final = img.resize((96, 96), Image.NEAREST)
            path = os.path.join(chars_dir, f"{name}_{state}.png")
            final.save(path)
            # 統計非透明像素
            non_transparent = sum(1 for y in range(96) for x in range(96)
                                  if final.getpixel((x,y))[3] > 50)
            print(f"  OK {name}_{state}.png - {non_transparent} visible pixels")

if __name__ == "__main__":
    generate_all()
    print("\nDone!")
