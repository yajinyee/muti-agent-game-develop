# -*- coding: utf-8 -*-
"""
改善目標物：讓它們更大更清晰
重新生成所有目標物，尺寸統一為 48x48
"""
from PIL import Image, ImageDraw
import os
import math

OUTPUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
os.makedirs(OUTPUT_DIR, exist_ok=True)

def new_img(w, h):
    return Image.new("RGBA", (w, h), (0, 0, 0, 0))

def draw_pixel_char(size=48):
    """基礎像素角色模板"""
    img = new_img(size, size)
    draw = ImageDraw.Draw(img)
    return img, draw

def gen_grass(size=48):
    """像素雜草 T001"""
    img, draw = draw_pixel_char(size)
    GREEN = (60, 160, 50, 255)
    DARK_G = (30, 100, 20, 255)
    STEM = (80, 120, 40, 255)
    
    cx = size // 2
    # 草莖
    draw.rectangle([cx-3, size//2, cx+3, size-4], fill=STEM)
    # 左葉
    draw.polygon([(cx, size//2), (cx-size//3, size//4), (cx-2, size//2-4)], fill=GREEN, outline=DARK_G)
    # 右葉
    draw.polygon([(cx, size//2), (cx+size//3, size//4), (cx+2, size//2-4)], fill=GREEN, outline=DARK_G)
    # 頂葉
    draw.polygon([(cx, 4), (cx-8, size//3), (cx+8, size//3)], fill=(100, 200, 60, 255), outline=DARK_G)
    return img

def gen_bug(color, size=48):
    """小蟲"""
    img, draw = draw_pixel_char(size)
    r, g, b = color
    DARK = (max(0,r-60), max(0,g-60), max(0,b-60), 255)
    BODY = (*color, 255)
    
    cx, cy = size//2, size//2
    # 身體（橢圓）
    draw.ellipse([cx-size//3, cy-size//5, cx+size//3, cy+size//5], fill=BODY, outline=DARK)
    # 頭
    draw.ellipse([cx-size//4, cy-size//3, cx+size//4, cy-size//8], fill=BODY, outline=DARK)
    # 眼睛
    draw.ellipse([cx-10, cy-size//3+4, cx-4, cy-size//3+10], fill=(255,255,255,255))
    draw.ellipse([cx-9, cy-size//3+5, cx-5, cy-size//3+9], fill=(20,20,20,255))
    # 觸角
    draw.line([(cx-8, cy-size//3+2), (cx-14, 4)], fill=DARK, width=2)
    draw.line([(cx-2, cy-size//3+2), (cx+4, 4)], fill=DARK, width=2)
    # 腳
    for i in range(3):
        y_leg = cy + i*4
        draw.line([(cx-size//3, y_leg), (cx-size//3-8, y_leg+8)], fill=DARK, width=2)
        draw.line([(cx+size//3, y_leg), (cx+size//3+8, y_leg+8)], fill=DARK, width=2)
    return img

def gen_pudding(size=48):
    """布丁 T005"""
    img, draw = draw_pixel_char(size)
    YELLOW = (255, 220, 80, 255)
    DARK_Y = (180, 140, 20, 255)
    CREAM  = (255, 245, 180, 255)
    
    cx = size // 2
    # 底部
    draw.ellipse([4, size*2//3, size-5, size-4], fill=YELLOW, outline=DARK_Y)
    # 主體
    draw.ellipse([6, size//3, size-7, size*2//3+4], fill=YELLOW, outline=DARK_Y)
    # 頂部奶油
    draw.ellipse([10, 4, size-11, size//3+4], fill=CREAM, outline=DARK_Y)
    # 眼睛
    draw.ellipse([12, size//2-4, 18, size//2+2], fill=(40,20,10,255))
    draw.ellipse([size-18, size//2-4, size-12, size//2+2], fill=(40,20,10,255))
    draw.ellipse([13, size//2-3, 16, size//2+1], fill=(255,255,255,255))
    draw.ellipse([size-17, size//2-3, size-14, size//2+1], fill=(255,255,255,255))
    # 嘴巴
    draw.arc([cx-6, size//2+2, cx+6, size//2+10], 0, 180, fill=(40,20,10,255), width=2)
    return img

def gen_mushroom(size=48):
    """蘑菇 T006"""
    img, draw = draw_pixel_char(size)
    RED   = (220, 60, 40, 255)
    DARK_R= (140, 30, 20, 255)
    WHITE = (240, 240, 240, 255)
    STEM  = (220, 200, 160, 255)
    
    cx = size // 2
    # 莖
    draw.rectangle([cx-6, size//2, cx+6, size-4], fill=STEM, outline=(160,140,100,255))
    # 傘蓋
    draw.ellipse([4, 4, size-5, size//2+6], fill=RED, outline=DARK_R)
    # 白色斑點
    draw.ellipse([8, 10, 16, 18], fill=WHITE)
    draw.ellipse([size-16, 8, size-8, 16], fill=WHITE)
    draw.ellipse([cx-4, 6, cx+4, 14], fill=WHITE)
    # 眼睛
    draw.ellipse([12, size//2-2, 18, size//2+4], fill=(40,20,10,255))
    draw.ellipse([size-18, size//2-2, size-12, size//2+4], fill=(40,20,10,255))
    return img

def gen_mimic(size=48):
    """擬態型怪物 T101 - 看起來像寶石"""
    img, draw = draw_pixel_char(size)
    PURPLE = (160, 80, 220, 255)
    DARK_P = (80, 30, 120, 255)
    LIGHT  = (200, 140, 255, 255)
    
    cx, cy = size//2, size//2
    # 六角形寶石
    pts = []
    for i in range(6):
        angle = math.pi * i / 3 - math.pi/2
        r = size//2 - 4
        pts.append((cx + r*math.cos(angle), cy + r*math.sin(angle)))
    draw.polygon(pts, fill=PURPLE, outline=DARK_P)
    # 光澤
    draw.polygon([(cx, 6), (cx+12, cy-4), (cx, cy-2)], fill=LIGHT)
    # 隱藏眼睛
    draw.ellipse([cx-10, cy-4, cx-4, cy+2], fill=DARK_P)
    draw.ellipse([cx+4, cy-4, cx+10, cy+2], fill=DARK_P)
    draw.ellipse([cx-9, cy-3, cx-5, cy+1], fill=(255,200,255,255))
    draw.ellipse([cx+5, cy-3, cx+9, cy+1], fill=(255,200,255,255))
    return img

def gen_chest(size=48):
    """寶箱怪 T102"""
    img, draw = draw_pixel_char(size)
    GOLD  = (200, 160, 40, 255)
    DARK_G= (120, 80, 10, 255)
    BROWN = (140, 80, 30, 255)
    
    # 箱體
    draw.rectangle([4, size//2, size-5, size-4], fill=BROWN, outline=DARK_G)
    # 箱蓋
    draw.rectangle([4, 6, size-5, size//2+2], fill=GOLD, outline=DARK_G)
    # 金屬條
    draw.rectangle([4, size//2-2, size-5, size//2+2], fill=DARK_G)
    # 鎖
    draw.rectangle([size//2-5, size//2-6, size//2+5, size//2+4], fill=GOLD, outline=DARK_G)
    draw.ellipse([size//2-3, size//2-4, size//2+3, size//2+2], fill=DARK_G)
    # 眼睛
    draw.ellipse([8, size//2+6, 16, size//2+14], fill=(255,255,200,255), outline=DARK_G)
    draw.ellipse([size-16, size//2+6, size-8, size//2+14], fill=(255,255,200,255), outline=DARK_G)
    draw.ellipse([9, size//2+7, 15, size//2+13], fill=(40,20,10,255))
    draw.ellipse([size-15, size//2+7, size-9, size//2+13], fill=(40,20,10,255))
    # 牙齒
    for tx in range(6, size-6, 8):
        draw.rectangle([tx, size-8, tx+5, size-4], fill=(240,240,240,255))
    return img

def gen_meteor(size=48):
    """流星 T103"""
    img, draw = draw_pixel_char(size)
    YELLOW = (255, 230, 50, 255)
    WHITE  = (255, 255, 200, 255)
    ORANGE = (255, 150, 30, 255)
    
    # 尾巴
    for i in range(size*2//3):
        alpha = int(200 * (1 - i/(size*2//3)))
        half = max(1, 6 - i//6)
        draw.rectangle([i, size//2-half, i+3, size//2+half], fill=(255, 200, 50, alpha))
    # 星體
    draw.ellipse([size*2//3, 6, size-4, size-6], fill=YELLOW, outline=ORANGE)
    draw.ellipse([size*2//3+4, 10, size-8, size-10], fill=WHITE)
    # 光芒
    cx = size*5//6
    draw.line([(cx, 2), (cx, 6)], fill=YELLOW, width=2)
    draw.line([(cx, size-6), (cx, size-2)], fill=YELLOW, width=2)
    return img

def gen_gold_grass(size=48):
    """金色雜草 T104"""
    img, draw = draw_pixel_char(size)
    GOLD  = (255, 200, 30, 255)
    DARK_G= (180, 130, 10, 255)
    GLOW  = (255, 240, 150, 255)
    
    cx = size // 2
    draw.rectangle([cx-3, size//2, cx+3, size-4], fill=GOLD)
    draw.polygon([(cx, size//2), (cx-size//3, size//4), (cx-2, size//2-4)], fill=GOLD, outline=DARK_G)
    draw.polygon([(cx, size//2), (cx+size//3, size//4), (cx+2, size//2-4)], fill=GOLD, outline=DARK_G)
    draw.polygon([(cx, 4), (cx-8, size//3), (cx+8, size//3)], fill=GLOW, outline=DARK_G)
    # 閃光點
    draw.ellipse([6, 4, 12, 10], fill=(255,255,200,255))
    draw.ellipse([size-12, 2, size-6, 8], fill=(255,255,200,255))
    return img

def gen_coin_fish(size=48):
    """巨大金幣魚 T105"""
    img, draw = draw_pixel_char(size)
    GOLD  = (220, 180, 30, 255)
    DARK_G= (140, 100, 10, 255)
    LIGHT = (255, 230, 100, 255)
    
    # 魚尾
    draw.polygon([(size*2//3, size//2), (size-4, 6), (size-4, size-6)], fill=GOLD, outline=DARK_G)
    # 魚身
    draw.ellipse([4, 6, size*2//3+4, size-6], fill=GOLD, outline=DARK_G)
    # 金幣紋路
    draw.ellipse([10, 12, size//2, size-12], fill=LIGHT, outline=DARK_G)
    draw.ellipse([14, 16, size//2-4, size-16], fill=GOLD)
    # 眼睛
    draw.ellipse([6, size//2-8, 14, size//2], fill=(255,255,200,255), outline=DARK_G)
    draw.ellipse([7, size//2-7, 13, size//2-1], fill=(40,20,10,255))
    draw.ellipse([7, size//2-7, 9, size//2-5], fill=(255,255,255,255))
    return img

def gen_boss(size=96):
    """那個孩子 BOSS B001 - 更大更詳細"""
    img, draw = draw_pixel_char(size)
    BODY  = (220, 200, 180, 255)
    OUTLINE=(80, 50, 30, 255)
    RED   = (220, 50, 50, 255)
    WHITE = (255, 255, 255, 255)
    BLACK = (20, 10, 5, 255)
    BLUSH = (255, 150, 150, 255)
    
    cx, cy = size//2, size//2
    # 小圓耳朵
    draw.ellipse([10, 10, 28, 28], fill=BODY, outline=OUTLINE)
    draw.ellipse([size-28, 10, size-10, 28], fill=BODY, outline=OUTLINE)
    draw.ellipse([13, 13, 25, 25], fill=BLUSH)
    draw.ellipse([size-25, 13, size-13, 25], fill=BLUSH)
    # 大頭
    draw.ellipse([8, 12, size-9, size-10], fill=BODY, outline=OUTLINE)
    # 空洞眼睛
    draw.ellipse([18, 28, 36, 46], fill=WHITE, outline=OUTLINE)
    draw.ellipse([size-36, 28, size-18, 46], fill=WHITE, outline=OUTLINE)
    draw.ellipse([20, 30, 34, 44], fill=BLACK)
    draw.ellipse([size-34, 30, size-20, 44], fill=BLACK)
    draw.ellipse([20, 30, 24, 34], fill=WHITE)
    draw.ellipse([size-24, 30, size-20, 34], fill=WHITE)
    # 詭異微笑
    draw.arc([cx-18, cy+4, cx+18, cy+24], 0, 180, fill=OUTLINE, width=3)
    # 牙齒
    for tx in range(cx-14, cx+14, 8):
        draw.rectangle([tx, cy+12, tx+5, cy+20], fill=WHITE, outline=OUTLINE)
    # 身體
    draw.ellipse([20, size-30, size-20, size-4], fill=BODY, outline=OUTLINE)
    # 紅色警告邊框
    for t in range(4):
        draw.rectangle([t, t, size-1-t, size-1-t], outline=(*RED[:3], 255-t*40))
    return img

def main():
    print("Generating improved targets (48x48)...")
    
    targets = {
        "T001_grass":    (gen_grass, 48),
        "T002_bug_g":    (lambda s: gen_bug((80, 200, 80), s), 48),
        "T003_bug_r":    (lambda s: gen_bug((220, 80, 80), s), 48),
        "T004_bug_b":    (lambda s: gen_bug((80, 120, 220), s), 48),
        "T005_pudding":  (gen_pudding, 48),
        "T006_mushroom": (gen_mushroom, 48),
        "T101_mimic":    (gen_mimic, 48),
        "T102_chest":    (gen_chest, 48),
        "T103_meteor":   (gen_meteor, 48),
        "T104_gold_grass":(gen_gold_grass, 48),
        "T105_coin_fish":(gen_coin_fish, 48),
        "B001_boss":     (gen_boss, 96),
    }
    
    for name, (fn, size) in targets.items():
        img = fn(size)
        path = os.path.join(OUTPUT_DIR, f"{name}.png")
        img.save(path)
        print(f"  OK {name}.png ({size}x{size})")
    
    print("\nDone!")

if __name__ == "__main__":
    main()
