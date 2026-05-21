# -*- coding: utf-8 -*-
"""
目標物 v3 - 64x64 像素藝術，帶陰影和細節
所有目標物統一升級到 64x64
"""
from PIL import Image
import os
import math

OUTPUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def new_img():
    return Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

def px(img, x, y, c):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), c)

def fill_circle(img, cx, cy, r, color):
    for y in range(cy-r, cy+r+1):
        for x in range(cx-r, cx+r+1):
            if (x-cx)**2 + (y-cy)**2 <= r**2:
                px(img, x, y, color)

def fill_circle_shaded(img, cx, cy, r, base_rgb):
    """帶陰影的圓形：左上亮，右下暗"""
    r_v, g_v, b_v = base_rgb
    light = (min(255,r_v+35), min(255,g_v+35), min(255,b_v+35), 255)
    mid   = (r_v, g_v, b_v, 255)
    dark  = (max(0,r_v-40), max(0,g_v-40), max(0,b_v-40), 255)
    for y in range(cy-r, cy+r+1):
        for x in range(cx-r, cx+r+1):
            if (x-cx)**2 + (y-cy)**2 > r**2:
                continue
            nx_ = (x-cx)/max(r,1)
            ny_ = (y-cy)/max(r,1)
            dot = -(nx_*(-0.7) + ny_*(-0.7))
            if dot > 0.25:
                c = light
            elif dot < -0.1:
                c = dark
            else:
                c = mid
            px(img, x, y, c)

def outline_circle(img, cx, cy, r, color):
    for y in range(cy-r-1, cy+r+2):
        for x in range(cx-r-1, cx+r+2):
            d = math.sqrt((x-cx)**2 + (y-cy)**2)
            if r+0.1 <= d <= r+1.5:
                px(img, x, y, color)

def fill_rect_shaded(img, x1, y1, x2, y2, base_rgb):
    """帶陰影的矩形"""
    r_v, g_v, b_v = base_rgb
    w = x2 - x1
    h = y2 - y1
    for y in range(y1, y2):
        for x in range(x1, x2):
            # 左上亮，右下暗
            nx_ = (x - x1) / max(w, 1)
            ny_ = (y - y1) / max(h, 1)
            brightness = 1.0 - (nx_ + ny_) * 0.15
            r_c = int(min(255, r_v * brightness))
            g_c = int(min(255, g_v * brightness))
            b_c = int(min(255, b_v * brightness))
            px(img, x, y, (r_c, g_c, b_c, 255))


# ── 各目標物生成函數 ──────────────────────────────────────────────────────────

def gen_T001_grass():
    """像素雜草 2x - 三片葉子"""
    img = new_img()
    OUTLINE = (30, 80, 20, 255)
    STEM    = (60, 120, 40, 255)
    LEAF1   = (80, 180, 50, 255)
    LEAF2   = (100, 200, 60, 255)
    LEAF_L  = (50, 140, 35, 255)

    cx = 32
    # 草莖
    for y in range(38, 58):
        for x in range(cx-3, cx+4):
            px(img, x, y, STEM)
    # 輪廓
    for y in range(38, 58):
        px(img, cx-3, y, OUTLINE)
        px(img, cx+3, y, OUTLINE)

    # 左葉（斜向左上）
    for i in range(14):
        lx = cx - 2 - i*2
        ly = 38 - i*2
        for dx in range(-3, 4):
            c = LEAF1 if abs(dx) < 2 else LEAF_L
            px(img, lx+dx, ly, c)
        px(img, lx-3, ly, OUTLINE)
        px(img, lx+3, ly, OUTLINE)

    # 右葉（斜向右上）
    for i in range(14):
        lx = cx + 2 + i*2
        ly = 38 - i*2
        for dx in range(-3, 4):
            c = LEAF1 if abs(dx) < 2 else LEAF_L
            px(img, lx+dx, ly, c)
        px(img, lx-3, ly, OUTLINE)
        px(img, lx+3, ly, OUTLINE)

    # 中間葉（直向上）
    for i in range(18):
        ly = 38 - i*2
        for dx in range(-3, 4):
            c = LEAF2 if abs(dx) < 2 else LEAF1
            px(img, cx+dx, ly, c)
        px(img, cx-3, ly, OUTLINE)
        px(img, cx+3, ly, OUTLINE)

    # 頂部輪廓
    px(img, cx, 4, OUTLINE)
    for dx in range(-2, 3):
        px(img, cx+dx, 5, OUTLINE)

    return img


def gen_T002_bug(color_rgb, eye_color=(255,255,255)):
    """小蟲通用 - 橢圓身體+圓頭+觸角"""
    img = new_img()
    r_v, g_v, b_v = color_rgb
    OUTLINE = (max(0,r_v-60), max(0,g_v-60), max(0,b_v-60), 255)
    BODY    = (r_v, g_v, b_v, 255)
    LIGHT   = (min(255,r_v+40), min(255,g_v+40), min(255,b_v+40), 255)

    # 身體（橢圓）
    for y in range(30, 50):
        for x in range(12, 52):
            if ((x-32)/20)**2 + ((y-40)/10)**2 <= 1.0:
                nx_ = (x-32)/20
                ny_ = (y-40)/10
                dot = -(nx_*(-0.7)+ny_*(-0.7))
                c = LIGHT if dot > 0.3 else (BODY if dot > -0.1 else OUTLINE)
                px(img, x, y, c)

    # 身體輪廓
    for y in range(29, 51):
        for x in range(11, 53):
            if 0.9 <= ((x-32)/20)**2 + ((y-40)/10)**2 <= 1.15:
                px(img, x, y, OUTLINE)

    # 頭（圓形）
    fill_circle_shaded(img, 32, 22, 10, color_rgb)
    outline_circle(img, 32, 22, 10, OUTLINE)

    # 眼睛
    px(img, 27, 20, (255,255,255,255))
    px(img, 28, 20, (255,255,255,255))
    px(img, 27, 21, (255,255,255,255))
    px(img, 28, 21, (20,20,20,255))

    px(img, 36, 20, (255,255,255,255))
    px(img, 37, 20, (255,255,255,255))
    px(img, 36, 21, (255,255,255,255))
    px(img, 37, 21, (20,20,20,255))

    # 觸角
    for i in range(8):
        px(img, 28-i, 14-i, OUTLINE)
        px(img, 36+i, 14-i, OUTLINE)
    px(img, 20, 6, OUTLINE)
    px(img, 44, 6, OUTLINE)

    # 腳（3對）
    for i in range(3):
        y_leg = 35 + i*5
        for j in range(6):
            px(img, 12-j, y_leg+j//2, OUTLINE)
            px(img, 52+j, y_leg+j//2, OUTLINE)

    return img


def gen_T005_pudding():
    """會走路的布丁 8x"""
    img = new_img()
    YELLOW  = (255, 220, 60)
    CARAMEL = (200, 140, 30)
    CREAM   = (255, 245, 180)
    OUTLINE = (120, 70, 10, 255)
    SHINE   = (255, 255, 200, 255)

    # 布丁主體（梯形）
    for y in range(20, 52):
        width = 20 + (y-20)*14//32
        x1 = 32 - width//2
        x2 = 32 + width//2
        for x in range(x1, x2+1):
            nx_ = (x-32)/max(width//2,1)
            ny_ = (y-20)/32
            dot = -(nx_*(-0.7)+ny_*(-0.7))
            r_v = YELLOW[0] if dot > 0.2 else (CARAMEL[0] if dot < -0.1 else YELLOW[0])
            g_v = YELLOW[1] if dot > 0.2 else (CARAMEL[1] if dot < -0.1 else YELLOW[1])
            b_v = YELLOW[2] if dot > 0.2 else (CARAMEL[2] if dot < -0.1 else YELLOW[2])
            px(img, x, y, (r_v, g_v, b_v, 255))

    # 頂部奶油
    fill_circle_shaded(img, 32, 18, 10, CREAM)
    outline_circle(img, 32, 18, 10, OUTLINE)

    # 焦糖醬（頂部流下）
    for i in range(5):
        px(img, 26+i, 22+i, (*CARAMEL, 255))
        px(img, 38-i, 22+i, (*CARAMEL, 255))

    # 眼睛
    fill_circle(img, 28, 32, 3, (20, 20, 20, 255))
    fill_circle(img, 36, 32, 3, (20, 20, 20, 255))
    px(img, 27, 31, SHINE)
    px(img, 35, 31, SHINE)

    # 嘴巴
    for x in range(29, 36):
        px(img, x, 38, OUTLINE)
    px(img, 29, 37, OUTLINE)
    px(img, 35, 37, OUTLINE)

    # 輪廓
    for y in range(20, 52):
        width = 20 + (y-20)*14//32
        x1 = 32 - width//2
        x2 = 32 + width//2
        px(img, x1-1, y, OUTLINE)
        px(img, x2+1, y, OUTLINE)
    for x in range(22, 43):
        px(img, x, 52, OUTLINE)

    # 腳
    for i in range(4):
        px(img, 24+i*4, 54, OUTLINE)
        px(img, 24+i*4, 55, (*YELLOW, 255))
        px(img, 24+i*4, 56, OUTLINE)

    return img


def gen_T006_mushroom():
    """巨大蘑菇 10x"""
    img = new_img()
    RED    = (220, 50, 40)
    SPOT   = (255, 255, 255)
    STEM   = (240, 220, 180)
    OUTLINE= (80, 20, 10, 255)
    STEM_D = (180, 160, 120, 255)

    # 蘑菇傘（半圓）
    for y in range(8, 38):
        for x in range(6, 58):
            if ((x-32)/26)**2 + ((y-38)/30)**2 <= 1.0 and y <= 38:
                nx_ = (x-32)/26
                ny_ = (y-38)/30
                dot = -(nx_*(-0.7)+ny_*(-0.7))
                r_v = min(255, RED[0]+30) if dot > 0.3 else (max(0,RED[0]-30) if dot < -0.1 else RED[0])
                px(img, x, y, (r_v, RED[1], RED[2], 255))

    # 白色斑點
    for (sx, sy, sr) in [(22,18,4),(42,16,4),(32,10,3),(14,28,3),(50,26,3)]:
        fill_circle(img, sx, sy, sr, (*SPOT, 255))
        outline_circle(img, sx, sy, sr, OUTLINE)

    # 傘輪廓
    for y in range(7, 39):
        for x in range(5, 59):
            d = ((x-32)/26)**2 + ((y-38)/30)**2
            if 0.9 <= d <= 1.1 and y <= 38:
                px(img, x, y, OUTLINE)

    # 傘底邊
    for x in range(10, 54):
        px(img, x, 38, OUTLINE)
        px(img, x, 39, (*STEM, 255))

    # 莖
    for y in range(39, 58):
        for x in range(24, 40):
            nx_ = (x-32)/8
            c = STEM if abs(nx_) < 0.5 else STEM_D
            px(img, x, y, c)
    for y in range(39, 58):
        px(img, 23, y, OUTLINE)
        px(img, 40, y, OUTLINE)
    for x in range(23, 41):
        px(img, x, 58, OUTLINE)

    return img


def gen_T101_mimic():
    """擬態型怪物 15-30x - 看起來像普通雜草但有眼睛"""
    img = new_img()
    # 先畫雜草
    base = gen_T001_grass()
    img.paste(base, (0, 0))

    # 加上詭異的眼睛（偽裝破綻）
    EYE_W = (255, 255, 200, 255)
    EYE_P = (150, 0, 200, 255)
    OUTLINE = (50, 0, 80, 255)

    # 左眼（在葉子上）
    for (ex, ey) in [(20, 28), (44, 28)]:
        fill_circle(img, ex, ey, 5, EYE_W)
        fill_circle(img, ex, ey, 3, EYE_P)
        px(img, ex-1, ey-1, (255, 255, 255, 255))
        outline_circle(img, ex, ey, 5, OUTLINE)

    # 詭異的嘴（在莖上）
    for x in range(28, 37):
        px(img, x, 50, OUTLINE)
    for i in range(3):
        px(img, 29+i*3, 51, OUTLINE)

    return img


def gen_T102_chest():
    """寶箱怪 25x"""
    img = new_img()
    WOOD   = (160, 100, 40)
    WOOD_D = (100, 60, 20)
    GOLD   = (220, 180, 40)
    OUTLINE= (60, 30, 5, 255)
    LOCK   = (200, 160, 30, 255)
    EYE_W  = (255, 255, 255, 255)
    EYE_B  = (20, 20, 20, 255)

    # 箱體下半
    fill_rect_shaded(img, 8, 36, 56, 56, WOOD)
    # 箱蓋上半
    fill_rect_shaded(img, 8, 18, 56, 36, WOOD_D)

    # 金色邊框
    for x in range(8, 57):
        px(img, x, 18, (*GOLD, 255))
        px(img, x, 36, (*GOLD, 255))
        px(img, x, 56, (*GOLD, 255))
    for y in range(18, 57):
        px(img, 8, y, (*GOLD, 255))
        px(img, 56, y, (*GOLD, 255))

    # 金屬條紋
    for y in range(20, 56):
        px(img, 20, y, (*GOLD, 255))
        px(img, 44, y, (*GOLD, 255))

    # 鎖
    fill_circle(img, 32, 36, 5, LOCK)
    outline_circle(img, 32, 36, 5, OUTLINE)
    for y in range(38, 44):
        for x in range(29, 36):
            px(img, x, y, LOCK)
    for y in range(38, 44):
        px(img, 28, y, OUTLINE)
        px(img, 36, y, OUTLINE)
    px(img, 29, 44, OUTLINE)
    px(img, 35, 44, OUTLINE)

    # 眼睛（在箱蓋上）
    for (ex, ey) in [(22, 26), (42, 26)]:
        fill_circle(img, ex, ey, 4, EYE_W)
        fill_circle(img, ex, ey, 2, EYE_B)
        px(img, ex-1, ey-1, EYE_W)

    # 牙齒（箱蓋縫隙）
    for i in range(5):
        px(img, 16+i*8, 36, (255, 255, 255, 255))
        px(img, 16+i*8, 37, (255, 255, 255, 255))

    # 輪廓
    for x in range(7, 58):
        px(img, x, 17, OUTLINE)
        px(img, x, 57, OUTLINE)
    for y in range(17, 58):
        px(img, 7, y, OUTLINE)
        px(img, 57, y, OUTLINE)

    return img


def gen_T103_meteor():
    """流星 20-50x - 發光的星形"""
    img = new_img()
    CORE  = (255, 240, 100)
    GLOW1 = (255, 200, 50)
    GLOW2 = (255, 150, 20)
    TRAIL = (255, 100, 0)
    OUTLINE=(180, 80, 0, 255)

    # 流星尾巴（右上到左下）
    for i in range(20):
        x = 52 - i*2
        y = 12 + i*2
        alpha = max(0, 200 - i*10)
        for dx in range(-2, 3):
            a = max(0, alpha - abs(dx)*40)
            px(img, x+dx, y, (*TRAIL, a))

    # 星形核心
    fill_circle_shaded(img, 20, 44, 12, CORE)
    outline_circle(img, 20, 44, 12, OUTLINE)

    # 光暈
    for r in range(13, 17):
        for angle in range(0, 360, 5):
            rad = math.radians(angle)
            x = int(20 + r * math.cos(rad))
            y = int(44 + r * math.sin(rad))
            alpha = max(0, 120 - (r-13)*30)
            px(img, x, y, (*GLOW1, alpha))

    # 星芒（4方向）
    for i in range(6):
        px(img, 20, 44-13-i, (*GLOW2, max(0, 200-i*35)))
        px(img, 20, 44+13+i, (*GLOW2, max(0, 200-i*35)))
        px(img, 20-13-i, 44, (*GLOW2, max(0, 200-i*35)))
        px(img, 20+13+i, 44, (*GLOW2, max(0, 200-i*35)))

    return img


def gen_T104_gold_grass():
    """金色雜草 30x - 閃亮的金色"""
    img = new_img()
    GOLD   = (220, 180, 30)
    GOLD_L = (255, 230, 80)
    GOLD_D = (160, 120, 10)
    OUTLINE= (100, 60, 5, 255)
    SHINE  = (255, 255, 200, 255)

    cx = 32
    # 草莖（金色）
    for y in range(38, 58):
        for x in range(cx-3, cx+4):
            px(img, x, y, (*GOLD, 255))
    for y in range(38, 58):
        px(img, cx-3, y, OUTLINE)
        px(img, cx+3, y, OUTLINE)

    # 三片金色葉子
    for i in range(14):
        # 左葉
        lx = cx - 2 - i*2
        ly = 38 - i*2
        for dx in range(-3, 4):
            c = GOLD_L if abs(dx) < 2 else GOLD
            px(img, lx+dx, ly, (*c, 255))
        px(img, lx-3, ly, OUTLINE)
        px(img, lx+3, ly, OUTLINE)

        # 右葉
        lx = cx + 2 + i*2
        for dx in range(-3, 4):
            c = GOLD_L if abs(dx) < 2 else GOLD
            px(img, lx+dx, ly, (*c, 255))
        px(img, lx-3, ly, OUTLINE)
        px(img, lx+3, ly, OUTLINE)

    # 中間葉
    for i in range(18):
        ly = 38 - i*2
        for dx in range(-3, 4):
            c = GOLD_L if abs(dx) < 2 else GOLD
            px(img, cx+dx, ly, (*c, 255))
        px(img, cx-3, ly, OUTLINE)
        px(img, cx+3, ly, OUTLINE)

    # 閃光點
    for (sx, sy) in [(20, 20), (44, 16), (32, 8), (14, 32), (50, 28)]:
        px(img, sx, sy, SHINE)
        for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
            px(img, sx+dx, sy+dy, (*GOLD_L, 200))

    return img


def gen_T105_coin_fish():
    """巨大金幣魚 50x"""
    img = new_img()
    GOLD   = (220, 180, 40)
    GOLD_L = (255, 230, 80)
    GOLD_D = (160, 120, 10)
    OUTLINE= (80, 50, 5, 255)
    EYE_W  = (255, 255, 255, 255)
    EYE_B  = (20, 20, 20, 255)
    COIN   = (255, 200, 30)

    # 魚身（橢圓）
    for y in range(18, 50):
        for x in range(10, 54):
            if ((x-32)/22)**2 + ((y-34)/16)**2 <= 1.0:
                nx_ = (x-32)/22
                ny_ = (y-34)/16
                dot = -(nx_*(-0.7)+ny_*(-0.7))
                r_v = GOLD_L[0] if dot > 0.3 else (GOLD_D[0] if dot < -0.1 else GOLD[0])
                g_v = GOLD_L[1] if dot > 0.3 else (GOLD_D[1] if dot < -0.1 else GOLD[1])
                b_v = GOLD_L[2] if dot > 0.3 else (GOLD_D[2] if dot < -0.1 else GOLD[2])
                px(img, x, y, (r_v, g_v, b_v, 255))

    # 魚身輪廓
    for y in range(17, 51):
        for x in range(9, 55):
            d = ((x-32)/22)**2 + ((y-34)/16)**2
            if 0.9 <= d <= 1.1:
                px(img, x, y, OUTLINE)

    # 魚尾
    for i in range(12):
        for j in range(-i//2, i//2+1):
            px(img, 54+i, 34+j, (*GOLD, 255))
        px(img, 54+i, 34-i//2-1, OUTLINE)
        px(img, 54+i, 34+i//2+1, OUTLINE)

    # 魚鰭（上）
    for i in range(8):
        for j in range(i+1):
            px(img, 28+j, 18-i, (*GOLD_L, 255))
        px(img, 28+i, 18-i-1, OUTLINE)

    # 眼睛
    fill_circle(img, 20, 30, 5, EYE_W)
    fill_circle(img, 20, 30, 3, EYE_B)
    px(img, 19, 29, EYE_W)

    # 金幣符號（身上）
    for y in range(28, 40):
        for x in range(30, 42):
            if ((x-36)/6)**2 + ((y-34)/6)**2 <= 1.0:
                px(img, x, y, (*COIN, 255))
    for y in range(27, 41):
        for x in range(29, 43):
            d = ((x-36)/6)**2 + ((y-34)/6)**2
            if 0.9 <= d <= 1.15:
                px(img, x, y, OUTLINE)
    # ¥ 符號
    px(img, 36, 30, OUTLINE)
    px(img, 36, 31, OUTLINE)
    px(img, 36, 32, OUTLINE)
    for x in range(33, 40):
        px(img, x, 33, OUTLINE)
    for x in range(33, 40):
        px(img, x, 35, OUTLINE)

    return img


def gen_B001_boss():
    """那個孩子 BOSS 100-500x - 96x96"""
    SIZE_B = 96
    img = Image.new("RGBA", (SIZE_B, SIZE_B), (0, 0, 0, 0))

    def px_b(x, y, c):
        if 0 <= x < SIZE_B and 0 <= y < SIZE_B:
            img.putpixel((x, y), c)

    def fill_c(cx, cy, r, color):
        for y in range(cy-r, cy+r+1):
            for x in range(cx-r, cx+r+1):
                if (x-cx)**2 + (y-cy)**2 <= r**2:
                    px_b(x, y, color)

    def fill_cs(cx, cy, r, base_rgb):
        r_v, g_v, b_v = base_rgb
        for y in range(cy-r, cy+r+1):
            for x in range(cx-r, cx+r+1):
                if (x-cx)**2 + (y-cy)**2 > r**2:
                    continue
                nx_ = (x-cx)/max(r,1)
                ny_ = (y-cy)/max(r,1)
                dot = -(nx_*(-0.7)+ny_*(-0.7))
                if dot > 0.25:
                    c = (min(255,r_v+35), min(255,g_v+35), min(255,b_v+35), 255)
                elif dot < -0.1:
                    c = (max(0,r_v-40), max(0,g_v-40), max(0,b_v-40), 255)
                else:
                    c = (r_v, g_v, b_v, 255)
                px_b(x, y, c)

    def outline_c(cx, cy, r, color):
        for y in range(cy-r-1, cy+r+2):
            for x in range(cx-r-1, cx+r+2):
                d = math.sqrt((x-cx)**2 + (y-cy)**2)
                if r+0.1 <= d <= r+1.5:
                    px_b(x, y, color)

    WHITE   = (255, 252, 245)
    OUTLINE = (50, 30, 15, 255)
    BLUSH   = (255, 160, 155, 255)
    PINK    = (255, 140, 190, 255)
    EAR_IN  = (255, 200, 195, 255)
    RED_EYE = (200, 30, 30, 255)
    DARK_AURA = (80, 0, 80, 180)

    # 暗黑光環
    for r in range(44, 48):
        for angle in range(0, 360, 3):
            rad = math.radians(angle)
            x = int(48 + r * math.cos(rad))
            y = int(48 + r * math.sin(rad))
            alpha = max(0, 150 - (r-44)*35)
            px_b(x, y, (*DARK_AURA[:3], alpha))

    # 耳朵
    fill_cs(30, 22, 12, WHITE)
    fill_cs(66, 22, 12, WHITE)
    fill_c(30, 22, 6, EAR_IN[:3])
    fill_c(66, 22, 6, EAR_IN[:3])
    outline_c(30, 22, 12, OUTLINE)
    outline_c(66, 22, 12, OUTLINE)

    # 大圓頭
    fill_cs(48, 46, 28, WHITE)
    outline_c(48, 46, 28, OUTLINE)

    # 眼睛（紅色，BOSS 特徵）
    fill_c(38, 40, 5, (255,255,255,255))
    fill_c(38, 40, 3, RED_EYE[:3])
    px_b(37, 39, (255,255,255,255))
    fill_c(58, 40, 5, (255,255,255,255))
    fill_c(58, 40, 3, RED_EYE[:3])
    px_b(57, 39, (255,255,255,255))

    # 腮紅
    fill_c(32, 50, 5, BLUSH[:3])
    fill_c(64, 50, 5, BLUSH[:3])

    # 嘴巴（邪惡微笑）
    for x in range(40, 57):
        px_b(x, 58, OUTLINE)
    px_b(40, 57, OUTLINE)
    px_b(56, 57, OUTLINE)
    # 牙齒
    for i in range(3):
        for y in range(58, 62):
            px_b(42+i*5, y, (255,255,255,255))
            px_b(43+i*5, y, (255,255,255,255))

    # 身體
    fill_cs(48, 74, 12, WHITE)
    outline_c(48, 74, 12, OUTLINE)

    # 討伐棒（更大更威脅）
    for i in range(14):
        px_b(62+i, 28-i, (*PINK[:3], 255))
        px_b(62+i, 29-i, (200,100,150,255))
    fill_c(76, 14, 6, (255,215,235,255))
    outline_c(76, 14, 6, OUTLINE)

    return img


# ── 主程式 ────────────────────────────────────────────────────────────────────

def main():
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    print("=== 目標物 v3 生成（64x64）===\n")

    targets = [
        ("T001_grass",     gen_T001_grass()),
        ("T002_bug_g",     gen_T002_bug((60, 180, 60))),
        ("T003_bug_r",     gen_T002_bug((220, 60, 60))),
        ("T004_bug_b",     gen_T002_bug((60, 100, 220))),
        ("T005_pudding",   gen_T005_pudding()),
        ("T006_mushroom",  gen_T006_mushroom()),
        ("T101_mimic",     gen_T101_mimic()),
        ("T102_chest",     gen_T102_chest()),
        ("T103_meteor",    gen_T103_meteor()),
        ("T104_gold_grass",gen_T104_gold_grass()),
        ("T105_coin_fish", gen_T105_coin_fish()),
    ]

    for name, img in targets:
        # 放大到 64x64（已是 64x64，但確保一致）
        if img.size != (SIZE, SIZE):
            img = img.resize((SIZE, SIZE), Image.NEAREST)
        path = os.path.join(OUTPUT_DIR, f"{name}.png")
        img.save(path)
        non_t = sum(1 for px_v in img.getdata() if px_v[3] > 10)
        print(f"  ✅ {name}.png: {img.size}, {non_t}px")

    # BOSS 單獨處理（96x96）
    boss = gen_B001_boss()
    boss_path = os.path.join(OUTPUT_DIR, "B001_boss.png")
    boss.save(boss_path)
    non_t = sum(1 for px_v in boss.getdata() if px_v[3] > 10)
    print(f"  ✅ B001_boss.png: {boss.size}, {non_t}px")

    print(f"\n✅ 全部完成！輸出目錄: {OUTPUT_DIR}")

if __name__ == "__main__":
    main()
