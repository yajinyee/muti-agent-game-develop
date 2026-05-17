# -*- coding: utf-8 -*-
"""
像素美術生成器 v2 - 更高品質的 16-bit 風格
使用更細緻的像素繪製技術
"""
from PIL import Image, ImageDraw
import os
import math

OUTPUT_BASE = r"D:\Kiro\client\chiikawa-pixel\assets\sprites"

def new_img(w, h):
    return Image.new("RGBA", (w, h), (0, 0, 0, 0))

def px(img, x, y, color):
    """安全設定像素"""
    if 0 <= x < img.width and 0 <= y < img.height:
        img.putpixel((x, y), color)

def draw_circle_pixels(img, cx, cy, r, color, outline=None):
    """畫像素圓形"""
    draw = ImageDraw.Draw(img)
    if outline:
        draw.ellipse([cx-r-1, cy-r-1, cx+r+1, cy+r+1], fill=outline)
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], fill=color)

def draw_rect_pixels(img, x, y, w, h, color, outline=None):
    draw = ImageDraw.Draw(img)
    if outline:
        draw.rectangle([x-1, y-1, x+w, y+h], fill=outline)
    draw.rectangle([x, y, x+w-1, y+h-1], fill=color)

# ---- 角色生成（更細緻版本）----

def gen_chiikawa(state="idle"):
    """吉伊卡哇 - 32x32 像素"""
    img = new_img(32, 32)
    draw = ImageDraw.Draw(img)

    # 顏色定義
    body = (255, 220, 200, 255)      # 膚色
    outline = (80, 40, 20, 255)      # 深棕輪廓
    eye = (40, 20, 10, 255)          # 眼睛
    blush = (255, 160, 160, 255)     # 腮紅
    ear = (255, 200, 180, 255)       # 耳朵
    weapon_pink = (255, 150, 200, 255)  # 粉紅討伐棒
    weapon_glow = (255, 200, 230, 255)  # 劍氣

    offset_y = -3 if state == "attack" else (3 if state == "bigwin" else 0)

    # 耳朵
    draw.ellipse([8, 4+offset_y, 14, 11+offset_y], fill=ear, outline=outline)
    draw.ellipse([18, 4+offset_y, 24, 11+offset_y], fill=ear, outline=outline)
    # 耳朵內側
    draw.ellipse([9, 5+offset_y, 13, 10+offset_y], fill=blush)
    draw.ellipse([19, 5+offset_y, 23, 10+offset_y], fill=blush)

    # 頭部（大圓）
    draw.ellipse([5, 7+offset_y, 27, 26+offset_y], fill=body, outline=outline)

    # 眼睛（閉著，努力感）
    if state == "bigwin":
        # 大獎：眼睛睜大
        draw.ellipse([9, 13+offset_y, 13, 17+offset_y], fill=eye)
        draw.ellipse([19, 13+offset_y, 23, 17+offset_y], fill=eye)
        draw.ellipse([10, 14+offset_y, 12, 16+offset_y], fill=(255,255,255,255))
    else:
        # 一般：閉眼弧線
        draw.arc([9, 13+offset_y, 13, 17+offset_y], 0, 180, fill=eye, width=2)
        draw.arc([19, 13+offset_y, 23, 17+offset_y], 0, 180, fill=eye, width=2)

    # 腮紅
    draw.ellipse([7, 17+offset_y, 11, 20+offset_y], fill=blush)
    draw.ellipse([21, 17+offset_y, 25, 20+offset_y], fill=blush)

    # 嘴巴
    draw.arc([13, 18+offset_y, 19, 22+offset_y], 0, 180, fill=outline, width=1)

    # 身體（小橢圓）
    draw.ellipse([9, 23+offset_y, 23, 31+offset_y], fill=body, outline=outline)

    # 討伐棒（攻擊狀態更明顯）
    if state == "attack":
        # 揮棒動作
        draw.line([(22, 10+offset_y), (30, 2+offset_y)], fill=weapon_pink, width=3)
        draw.ellipse([28, 0+offset_y, 32, 4+offset_y], fill=weapon_glow)
    else:
        # 持棒待機
        draw.line([(22, 14+offset_y), (28, 8+offset_y)], fill=weapon_pink, width=2)
        draw.ellipse([26, 6+offset_y, 30, 10+offset_y], fill=weapon_glow)

    return img

def gen_hachiware(state="idle"):
    """小八 - 32x32 像素"""
    img = new_img(32, 32)
    draw = ImageDraw.Draw(img)

    white = (240, 240, 240, 255)
    blue_stripe = (100, 150, 220, 255)
    outline = (40, 60, 100, 255)
    eye_blue = (60, 100, 200, 255)
    weapon_blue = (100, 160, 255, 255)
    weapon_glow = (180, 210, 255, 255)

    offset_y = -3 if state == "attack" else (3 if state == "bigwin" else 0)

    # 尖耳朵
    draw.polygon([(8, 12+offset_y), (11, 3+offset_y), (14, 12+offset_y)], fill=white, outline=outline)
    draw.polygon([(18, 12+offset_y), (21, 3+offset_y), (24, 12+offset_y)], fill=white, outline=outline)

    # 頭部
    draw.ellipse([5, 8+offset_y, 27, 27+offset_y], fill=white, outline=outline)

    # 藍色條紋（特徵）
    for i in range(3):
        x = 8 + i * 5
        draw.line([(x, 10+offset_y), (x, 25+offset_y)], fill=blue_stripe, width=2)

    # 眼睛（精神感）
    draw.ellipse([9, 13+offset_y, 14, 18+offset_y], fill=eye_blue, outline=outline)
    draw.ellipse([18, 13+offset_y, 23, 18+offset_y], fill=eye_blue, outline=outline)
    draw.ellipse([10, 14+offset_y, 13, 17+offset_y], fill=(255,255,255,255))
    draw.ellipse([19, 14+offset_y, 22, 17+offset_y], fill=(255,255,255,255))

    # 嘴巴（微笑）
    draw.arc([13, 19+offset_y, 19, 23+offset_y], 0, 180, fill=outline, width=1)

    # 身體
    draw.ellipse([9, 24+offset_y, 23, 31+offset_y], fill=white, outline=outline)

    # 討伐棒
    if state == "attack":
        draw.line([(22, 8+offset_y), (31, 1+offset_y)], fill=weapon_blue, width=3)
        draw.ellipse([29, -1+offset_y, 33, 3+offset_y], fill=weapon_glow)
    else:
        draw.line([(22, 12+offset_y), (29, 6+offset_y)], fill=weapon_blue, width=2)
        draw.ellipse([27, 4+offset_y, 31, 8+offset_y], fill=weapon_glow)

    return img

def gen_usagi(state="idle"):
    """烏薩奇 - 32x32 像素"""
    img = new_img(32, 32)
    draw = ImageDraw.Draw(img)

    white = (245, 245, 245, 255)
    pink_inner = (255, 180, 180, 255)
    outline = (80, 60, 80, 255)
    eye_red = (220, 60, 60, 255)
    weapon_yellow = (255, 220, 50, 255)
    weapon_glow = (255, 240, 150, 255)

    offset_y = -4 if state == "attack" else (4 if state == "bigwin" else 0)

    # 長耳朵（特徵）
    draw.ellipse([8, 0+offset_y, 14, 14+offset_y], fill=white, outline=outline)
    draw.ellipse([18, 0+offset_y, 24, 14+offset_y], fill=white, outline=outline)
    # 耳朵內側粉紅
    draw.ellipse([9, 1+offset_y, 13, 12+offset_y], fill=pink_inner)
    draw.ellipse([19, 1+offset_y, 23, 12+offset_y], fill=pink_inner)

    # 頭部
    draw.ellipse([5, 9+offset_y, 27, 27+offset_y], fill=white, outline=outline)

    # 眼睛（紅色，興奮感）
    draw.ellipse([9, 14+offset_y, 14, 19+offset_y], fill=eye_red, outline=outline)
    draw.ellipse([18, 14+offset_y, 23, 19+offset_y], fill=eye_red, outline=outline)
    draw.ellipse([10, 15+offset_y, 13, 18+offset_y], fill=(255,200,200,255))

    # 嘴巴（興奮）
    draw.arc([12, 20+offset_y, 20, 24+offset_y], 0, 180, fill=outline, width=2)

    # 身體
    draw.ellipse([9, 24+offset_y, 23, 31+offset_y], fill=white, outline=outline)

    # 討伐棒（旋轉感）
    if state == "attack":
        # 旋轉攻擊
        draw.line([(20, 8+offset_y), (31, -1+offset_y)], fill=weapon_yellow, width=3)
        draw.line([(22, 10+offset_y), (31, 3+offset_y)], fill=weapon_yellow, width=2)
        draw.ellipse([29, -3+offset_y, 33, 1+offset_y], fill=weapon_glow)
    else:
        draw.line([(22, 12+offset_y), (30, 5+offset_y)], fill=weapon_yellow, width=2)
        draw.ellipse([28, 3+offset_y, 32, 7+offset_y], fill=weapon_glow)

    return img

# ---- 目標物生成（更細緻版本）----

def gen_target(target_id):
    """生成目標物像素圖"""
    targets = {
        "T001": gen_grass,
        "T002": lambda: gen_bug((80, 200, 80, 255)),
        "T003": lambda: gen_bug((220, 80, 80, 255)),
        "T004": lambda: gen_bug((80, 120, 220, 255)),
        "T005": gen_pudding,
        "T006": gen_mushroom,
        "T101": gen_mimic,
        "T102": gen_chest,
        "T103": gen_meteor,
        "T104": gen_gold_grass,
        "T105": gen_coin_fish,
        "B001": gen_boss,
    }
    fn = targets.get(target_id)
    if fn:
        return fn()
    return new_img(32, 32)

def gen_grass():
    img = new_img(20, 28)
    draw = ImageDraw.Draw(img)
    # 草莖
    draw.rectangle([8, 12, 11, 27], fill=(60, 140, 40, 255))
    # 葉子
    draw.polygon([(10, 12), (2, 4), (8, 10)], fill=(80, 180, 50, 255), outline=(40, 100, 20, 255))
    draw.polygon([(10, 10), (18, 3), (12, 9)], fill=(100, 200, 60, 255), outline=(40, 100, 20, 255))
    draw.polygon([(9, 8), (5, 0), (11, 6)], fill=(120, 220, 70, 255), outline=(40, 100, 20, 255))
    return img

def gen_bug(color):
    img = new_img(24, 16)
    draw = ImageDraw.Draw(img)
    r, g, b, a = color
    dark = (max(0,r-60), max(0,g-60), max(0,b-60), 255)
    # 身體
    draw.ellipse([4, 3, 19, 12], fill=color, outline=dark)
    # 頭
    draw.ellipse([0, 4, 8, 11], fill=color, outline=dark)
    # 眼睛
    draw.ellipse([1, 5, 4, 8], fill=(255,255,255,255))
    draw.ellipse([2, 6, 3, 7], fill=(20,20,20,255))
    # 觸角
    draw.line([(3, 4), (0, 0)], fill=dark, width=1)
    draw.line([(5, 4), (3, 0)], fill=dark, width=1)
    # 腳
    for i in range(3):
        y = 9 + i * 0
        draw.line([(6+i*4, 11), (4+i*4, 15)], fill=dark, width=1)
        draw.line([(6+i*4, 11), (8+i*4, 15)], fill=dark, width=1)
    return img

def gen_pudding():
    img = new_img(28, 32)
    draw = ImageDraw.Draw(img)
    yellow = (255, 220, 80, 255)
    dark_y = (180, 140, 20, 255)
    cream = (255, 240, 180, 255)
    # 布丁底部
    draw.ellipse([2, 18, 25, 30], fill=yellow, outline=dark_y)
    # 布丁主體
    draw.ellipse([4, 8, 23, 22], fill=yellow, outline=dark_y)
    # 頂部奶油
    draw.ellipse([8, 2, 19, 12], fill=cream, outline=dark_y)
    # 眼睛
    draw.ellipse([8, 12, 12, 16], fill=(40,20,10,255))
    draw.ellipse([16, 12, 20, 16], fill=(40,20,10,255))
    draw.ellipse([9, 13, 11, 15], fill=(255,255,255,255))
    draw.ellipse([17, 13, 19, 15], fill=(255,255,255,255))
    # 嘴巴
    draw.arc([11, 16, 17, 20], 0, 180, fill=(40,20,10,255), width=1)
    # 腳
    draw.rectangle([8, 28, 11, 31], fill=dark_y)
    draw.rectangle([16, 28, 19, 31], fill=dark_y)
    return img

def gen_mushroom():
    img = new_img(32, 32)
    draw = ImageDraw.Draw(img)
    red = (220, 60, 40, 255)
    dark_r = (140, 30, 20, 255)
    white = (240, 240, 240, 255)
    stem = (220, 200, 160, 255)
    # 莖
    draw.rectangle([11, 18, 20, 30], fill=stem, outline=(160,140,100,255))
    # 傘蓋
    draw.ellipse([2, 4, 29, 22], fill=red, outline=dark_r)
    # 白色斑點
    draw.ellipse([6, 8, 11, 13], fill=white)
    draw.ellipse([20, 7, 25, 12], fill=white)
    draw.ellipse([13, 5, 17, 9], fill=white)
    # 眼睛（在傘蓋下方）
    draw.ellipse([10, 17, 14, 21], fill=(40,20,10,255))
    draw.ellipse([18, 17, 22, 21], fill=(40,20,10,255))
    return img

def gen_mimic():
    """擬態型怪物 - 初始看起來像寶石"""
    img = new_img(28, 28)
    draw = ImageDraw.Draw(img)
    purple = (160, 80, 220, 255)
    dark_p = (80, 30, 120, 255)
    light_p = (200, 140, 255, 255)
    # 寶石形狀
    draw.polygon([(14, 2), (24, 10), (24, 18), (14, 26), (4, 18), (4, 10)],
                 fill=purple, outline=dark_p)
    # 光澤
    draw.polygon([(14, 4), (20, 10), (14, 8)], fill=light_p)
    # 隱藏的眼睛（細看才看到）
    draw.ellipse([9, 12, 13, 16], fill=dark_p)
    draw.ellipse([15, 12, 19, 16], fill=dark_p)
    draw.ellipse([10, 13, 12, 15], fill=(255,200,255,255))
    draw.ellipse([16, 13, 18, 15], fill=(255,200,255,255))
    return img

def gen_chest():
    """寶箱怪"""
    img = new_img(32, 28)
    draw = ImageDraw.Draw(img)
    gold = (200, 160, 40, 255)
    dark_g = (120, 80, 10, 255)
    brown = (140, 80, 30, 255)
    # 箱體
    draw.rectangle([2, 12, 29, 26], fill=brown, outline=dark_g)
    # 箱蓋
    draw.rectangle([2, 4, 29, 14], fill=gold, outline=dark_g)
    # 金屬條
    draw.rectangle([2, 12, 29, 14], fill=dark_g)
    # 鎖
    draw.rectangle([13, 10, 18, 16], fill=gold, outline=dark_g)
    draw.ellipse([14, 11, 17, 14], fill=dark_g)
    # 眼睛（寶箱怪特徵）
    draw.ellipse([6, 16, 11, 21], fill=(255,255,200,255), outline=dark_g)
    draw.ellipse([21, 16, 26, 21], fill=(255,255,200,255), outline=dark_g)
    draw.ellipse([7, 17, 10, 20], fill=(40,20,10,255))
    draw.ellipse([22, 17, 25, 20], fill=(40,20,10,255))
    # 牙齒
    for i in range(4):
        draw.rectangle([5+i*6, 22, 7+i*6, 26], fill=(240,240,240,255))
    return img

def gen_meteor():
    """流星"""
    img = new_img(36, 20)
    draw = ImageDraw.Draw(img)
    yellow = (255, 230, 50, 255)
    white = (255, 255, 200, 255)
    orange = (255, 150, 30, 255)
    # 尾巴（漸層線條）
    for i in range(18):
        alpha = int(200 * (1 - i/18))
        y_center = 10
        half = max(1, 4 - i//5)
        draw.rectangle([i, y_center-half, i+2, y_center+half],
                       fill=(255, 200, 50, alpha))
    # 星體
    draw.ellipse([22, 4, 35, 16], fill=yellow, outline=orange)
    draw.ellipse([24, 6, 33, 14], fill=white)
    # 光芒
    draw.line([(28, 0), (28, 4)], fill=yellow, width=2)
    draw.line([(28, 16), (28, 20)], fill=yellow, width=2)
    draw.line([(22, 10), (18, 10)], fill=yellow, width=2)
    return img

def gen_gold_grass():
    """金色雜草"""
    img = new_img(20, 28)
    draw = ImageDraw.Draw(img)
    gold = (255, 200, 30, 255)
    dark_g = (180, 130, 10, 255)
    glow = (255, 240, 150, 255)
    # 草莖（金色）
    draw.rectangle([8, 12, 11, 27], fill=gold)
    # 葉子（金色，帶光暈）
    draw.polygon([(10, 12), (2, 4), (8, 10)], fill=gold, outline=dark_g)
    draw.polygon([(10, 10), (18, 3), (12, 9)], fill=gold, outline=dark_g)
    draw.polygon([(9, 8), (5, 0), (11, 6)], fill=glow, outline=dark_g)
    # 閃光點
    draw.ellipse([4, 2, 7, 5], fill=(255,255,200,255))
    draw.ellipse([14, 1, 17, 4], fill=(255,255,200,255))
    return img

def gen_coin_fish():
    """巨大金幣魚"""
    img = new_img(36, 28)
    draw = ImageDraw.Draw(img)
    gold = (220, 180, 30, 255)
    dark_g = (140, 100, 10, 255)
    light_g = (255, 230, 100, 255)
    # 魚尾
    draw.polygon([(28, 14), (36, 6), (36, 22)], fill=gold, outline=dark_g)
    # 魚身
    draw.ellipse([2, 4, 30, 24], fill=gold, outline=dark_g)
    # 金幣紋路
    draw.ellipse([8, 8, 22, 20], fill=light_g, outline=dark_g)
    draw.ellipse([11, 11, 19, 17], fill=gold)
    # 眼睛
    draw.ellipse([4, 8, 10, 14], fill=(255,255,200,255), outline=dark_g)
    draw.ellipse([5, 9, 9, 13], fill=(40,20,10,255))
    draw.ellipse([5, 9, 7, 11], fill=(255,255,255,255))
    # 鱗片
    for i in range(3):
        draw.arc([10+i*6, 6, 16+i*6, 14], 0, 180, fill=dark_g, width=1)
    return img

def gen_boss():
    """那個孩子 BOSS - 64x64"""
    img = new_img(64, 64)
    draw = ImageDraw.Draw(img)

    body = (220, 200, 180, 255)
    outline = (80, 50, 30, 255)
    eye_dark = (20, 10, 5, 255)
    red_glow = (220, 50, 50, 255)
    blush = (255, 150, 150, 255)

    # 小圓耳朵
    draw.ellipse([8, 8, 22, 22], fill=body, outline=outline)
    draw.ellipse([42, 8, 56, 22], fill=body, outline=outline)
    draw.ellipse([10, 10, 20, 20], fill=blush)
    draw.ellipse([44, 10, 54, 20], fill=blush)

    # 大頭
    draw.ellipse([6, 10, 58, 52], fill=body, outline=outline)

    # 空洞的眼睛（詭異感）
    draw.ellipse([14, 22, 26, 34], fill=(255,255,255,255), outline=outline)
    draw.ellipse([38, 22, 50, 34], fill=(255,255,255,255), outline=outline)
    draw.ellipse([16, 24, 24, 32], fill=eye_dark)
    draw.ellipse([40, 24, 48, 32], fill=eye_dark)
    # 眼睛高光
    draw.ellipse([16, 24, 19, 27], fill=(255,255,255,255))
    draw.ellipse([40, 24, 43, 27], fill=(255,255,255,255))

    # 詭異微笑
    draw.arc([20, 36, 44, 48], 0, 180, fill=outline, width=3)
    # 牙齒
    for i in range(4):
        draw.rectangle([22+i*5, 40, 24+i*5, 46], fill=(240,240,240,255), outline=outline)

    # 身體
    draw.ellipse([16, 48, 48, 63], fill=body, outline=outline)

    # 紅色警告邊框（BOSS 特徵）
    draw.rectangle([0, 0, 63, 63], outline=red_glow, width=3)

    return img

def generate_all():
    """生成所有角色和目標物"""
    chars_dir = os.path.join(OUTPUT_BASE, "characters")
    targets_dir = os.path.join(OUTPUT_BASE, "targets")
    os.makedirs(chars_dir, exist_ok=True)
    os.makedirs(targets_dir, exist_ok=True)

    print("[角色 Sprites v2]")
    chars = [("chiikawa", gen_chiikawa), ("hachiware", gen_hachiware), ("usagi", gen_usagi)]
    states = ["idle", "attack", "bigwin"]
    for char_name, gen_fn in chars:
        for state in states:
            img = gen_fn(state)
            # 放大 2x 讓細節更清晰
            img = img.resize((img.width * 2, img.height * 2), Image.NEAREST)
            path = os.path.join(chars_dir, f"{char_name}_{state}.png")
            img.save(path)
            print(f"  OK {char_name}_{state}.png ({img.width}x{img.height})")

    print("\n[目標物 Sprites v2]")
    target_ids = ["T001","T002","T003","T004","T005","T006",
                  "T101","T102","T103","T104","T105","B001"]
    target_names = {
        "T001":"grass","T002":"bug_g","T003":"bug_r","T004":"bug_b",
        "T005":"pudding","T006":"mushroom","T101":"mimic","T102":"chest",
        "T103":"meteor","T104":"gold_grass","T105":"coin_fish","B001":"boss"
    }
    for tid in target_ids:
        img = gen_target(tid)
        # 放大 2x
        img = img.resize((img.width * 2, img.height * 2), Image.NEAREST)
        name = target_names[tid]
        path = os.path.join(targets_dir, f"{tid}_{name}.png")
        img.save(path)
        print(f"  OK {tid}_{name}.png ({img.width}x{img.height})")

if __name__ == "__main__":
    print("生成像素美術 v2（更高品質）...")
    generate_all()
    print("\nOK 完成！")
