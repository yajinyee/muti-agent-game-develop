# -*- coding: utf-8 -*-
"""
目標物 v6 - T106-T117 新特殊目標物生成
DAY-153 補充：為 DAY-106 到 DAY-153 新增的 12 個特殊目標物生成 Sprite
每個目標物 64x64 像素，帶陰影和細節
"""
from PIL import Image, ImageDraw
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
    for y in range(max(0,cy-r), min(SIZE,cy+r+1)):
        for x in range(max(0,cx-r), min(SIZE,cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 <= r**2:
                px(img, x, y, color)

def fill_circle_shaded(img, cx, cy, r, base_rgb):
    r_v, g_v, b_v = base_rgb
    for y in range(max(0,cy-r), min(SIZE,cy+r+1)):
        for x in range(max(0,cx-r), min(SIZE,cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 > r**2:
                continue
            nx_ = (x-cx)/max(r,1)
            ny_ = (y-cy)/max(r,1)
            dot = -(nx_*(-0.7) + ny_*(-0.7))
            if dot > 0.25:
                c = (min(255,r_v+40), min(255,g_v+40), min(255,b_v+40), 255)
            elif dot < -0.1:
                c = (max(0,r_v-45), max(0,g_v-45), max(0,b_v-45), 255)
            else:
                c = (r_v, g_v, b_v, 255)
            px(img, x, y, c)

def outline_circle(img, cx, cy, r, color):
    for y in range(max(0,cy-r-2), min(SIZE,cy+r+3)):
        for x in range(max(0,cx-r-2), min(SIZE,cx+r+3)):
            d = math.sqrt((x-cx)**2 + (y-cy)**2)
            if r+0.1 <= d <= r+1.5:
                px(img, x, y, color)

def fill_rect(img, x1, y1, x2, y2, color):
    for y in range(max(0,y1), min(SIZE,y2)):
        for x in range(max(0,x1), min(SIZE,x2)):
            px(img, x, y, color)

def fill_rect_shaded(img, x1, y1, x2, y2, base_rgb):
    r_v, g_v, b_v = base_rgb
    w = max(1, x2-x1); h = max(1, y2-y1)
    for y in range(max(0,y1), min(SIZE,y2)):
        for x in range(max(0,x1), min(SIZE,x2)):
            nx_ = (x-x1)/w*2-1; ny_ = (y-y1)/h*2-1
            dot = -(nx_*(-0.7)+ny_*(-0.7))
            if dot > 0.25:
                c = (min(255,r_v+35), min(255,g_v+35), min(255,b_v+35), 255)
            elif dot < -0.1:
                c = (max(0,r_v-40), max(0,g_v-40), max(0,b_v-40), 255)
            else:
                c = (r_v, g_v, b_v, 255)
            px(img, x, y, c)

def draw_eye(img, cx, cy, pupil_color=(20,20,20,255)):
    fill_circle(img, cx, cy, 3, (255,255,255,255))
    fill_circle(img, cx, cy, 2, pupil_color)
    px(img, cx-1, cy-1, (255,255,255,200))

def outline_all(img, color=(20,20,20,255)):
    """為所有非透明像素加輪廓"""
    orig = img.copy()
    w, h = orig.size
    for y in range(h):
        for x in range(w):
            if orig.getpixel((x,y))[3] > 10:
                continue
            for dx,dy in [(-1,0),(1,0),(0,-1),(0,1)]:
                nx,ny = x+dx, y+dy
                if 0<=nx<w and 0<=ny<h and orig.getpixel((nx,ny))[3]>10:
                    img.putpixel((x,y), color)
                    break

# ---- T106 鑽頭龍蝦（橙紅色，帶鑽頭）----
def gen_T106():
    img = new_img()
    # 身體（橙紅色橢圓）
    fill_circle_shaded(img, 32, 36, 16, (220, 80, 20))
    # 鑽頭（右側三角形）
    for i in range(14):
        y_off = 14 - i
        for x in range(46+i, 46+i+2):
            px(img, x, 36-y_off//2, (180,60,10,255))
            px(img, x, 36+y_off//2, (180,60,10,255))
    # 鑽頭尖端
    fill_circle_shaded(img, 56, 36, 5, (150, 150, 160))
    # 眼睛
    draw_eye(img, 24, 30)
    # 觸角
    for i in range(8):
        px(img, 20-i, 22-i, (200,60,10,255))
        px(img, 22-i, 20-i, (200,60,10,255))
    # 爪子
    fill_circle_shaded(img, 18, 44, 6, (200,70,15))
    fill_circle_shaded(img, 14, 38, 5, (200,70,15))
    outline_all(img, (80,20,0,255))
    return img

# ---- T107 炸彈蟹（橙色，帶骷髏炸彈）----
def gen_T107():
    img = new_img()
    # 身體（橙色圓形）
    fill_circle_shaded(img, 32, 36, 18, (220, 120, 20))
    # 炸彈（黑色圓形在身體上）
    fill_circle_shaded(img, 32, 30, 8, (30, 30, 30))
    # 骷髏眼睛
    fill_circle(img, 29, 28, 2, (255,255,255,255))
    fill_circle(img, 35, 28, 2, (255,255,255,255))
    fill_circle(img, 29, 28, 1, (0,0,0,255))
    fill_circle(img, 35, 28, 1, (0,0,0,255))
    # 炸彈導火線
    for i in range(5):
        px(img, 32, 22-i, (200,150,50,255))
    px(img, 32, 17, (255,200,0,255))
    # 蟹眼睛（在身體下方）
    draw_eye(img, 26, 40)
    draw_eye(img, 38, 40)
    # 蟹爪
    fill_circle_shaded(img, 14, 36, 7, (200,100,15))
    fill_circle_shaded(img, 50, 36, 7, (200,100,15))
    # 蟹腳
    for i in range(3):
        px(img, 18+i*2, 50+i, (180,90,10,255))
        px(img, 46-i*2, 50+i, (180,90,10,255))
    outline_all(img, (80,30,0,255))
    return img

# ---- T108 巨型章魚（紫色，多觸手）----
def gen_T108():
    img = new_img()
    # 頭部（紫色大圓）
    fill_circle_shaded(img, 32, 26, 20, (140, 60, 180))
    # 眼睛（大眼）
    fill_circle(img, 24, 22, 5, (255,255,255,255))
    fill_circle(img, 40, 22, 5, (255,255,255,255))
    fill_circle(img, 24, 22, 3, (255,100,0,255))
    fill_circle(img, 40, 22, 3, (255,100,0,255))
    fill_circle(img, 24, 22, 1, (0,0,0,255))
    fill_circle(img, 40, 22, 1, (0,0,0,255))
    # 觸手（8條，波浪形）
    tentacle_positions = [(12,44),(18,48),(24,50),(30,52),(34,52),(40,50),(46,48),(52,44)]
    for i,(tx,ty) in enumerate(tentacle_positions):
        wave = 2 if i%2==0 else -2
        for j in range(8):
            px(img, tx + (j-4)*1, ty + j + wave*(j%3), (120,50,160,255))
    # 嘴巴
    for x in range(27,38):
        px(img, x, 32, (80,20,120,255))
    outline_all(img, (50,0,80,255))
    return img

# ---- T109 巨型鮟鱇魚（深藍，帶發光誘餌）----
def gen_T109():
    img = new_img()
    # 身體（深藍色）
    fill_circle_shaded(img, 32, 36, 18, (20, 40, 120))
    # 大嘴（下方）
    fill_rect(img, 16, 44, 48, 52, (10,20,80,255))
    # 牙齒
    for i in range(5):
        fill_rect(img, 18+i*6, 44, 21+i*6, 48, (240,240,240,255))
    # 眼睛（大眼）
    fill_circle(img, 26, 28, 6, (255,255,200,255))
    fill_circle(img, 26, 28, 4, (0,200,100,255))
    fill_circle(img, 26, 28, 2, (0,0,0,255))
    # 發光誘餌（頭頂）
    for i in range(6):
        px(img, 32, 14-i, (100,80,60,255))
    fill_circle(img, 32, 8, 5, (255,220,50,255))
    fill_circle(img, 32, 8, 3, (255,255,150,255))
    # 光暈
    for r in range(7,10):
        outline_circle(img, 32, 8, r, (255,200,0,50))
    outline_all(img, (0,10,60,255))
    return img

# ---- T110 巨型鹹水鱷魚（深綠，大嘴）----
def gen_T110():
    img = new_img()
    # 身體（深綠色橢圓）
    fill_circle_shaded(img, 32, 36, 18, (30, 100, 30))
    # 大嘴（右側）
    fill_rect_shaded(img, 44, 28, 62, 44, (20, 80, 20))
    # 牙齒（上下）
    for i in range(4):
        fill_rect(img, 46+i*4, 28, 48+i*4, 32, (240,240,220,255))
        fill_rect(img, 46+i*4, 40, 48+i*4, 44, (240,240,220,255))
    # 眼睛（突出）
    fill_circle_shaded(img, 28, 24, 6, (60, 140, 40))
    fill_circle(img, 28, 24, 4, (255,200,0,255))
    fill_circle(img, 28, 24, 2, (0,0,0,255))
    # 鱗片紋理
    for y in range(30, 50, 4):
        for x in range(16, 44, 6):
            px(img, x, y, (20,80,20,255))
    # 尾巴
    fill_circle_shaded(img, 12, 40, 8, (25,90,25))
    fill_circle_shaded(img, 6, 36, 5, (20,80,20))
    outline_all(img, (0,40,0,255))
    return img

# ---- T111 夢幻巨型獎勵魚（粉紅，夢幻光暈）----
def gen_T111():
    img = new_img()
    # 光暈（粉紅漸層）
    for r in range(22, 18, -1):
        alpha = int(60 * (22-r)/4)
        outline_circle(img, 32, 32, r, (255,150,200,alpha))
    # 身體（粉紅色魚形）
    fill_circle_shaded(img, 30, 32, 16, (255, 140, 180))
    # 魚尾
    for i in range(10):
        px(img, 50+i, 28+i//2, (220,100,150,255))
        px(img, 50+i, 36-i//2, (220,100,150,255))
    # 眼睛（大眼，星形高光）
    fill_circle(img, 22, 28, 5, (255,255,255,255))
    fill_circle(img, 22, 28, 3, (255,100,150,255))
    fill_circle(img, 22, 28, 1, (0,0,0,255))
    px(img, 21, 27, (255,255,255,255))
    # 星星裝飾
    for sx,sy in [(38,18),(44,26),(40,40),(34,46)]:
        px(img, sx, sy, (255,220,50,255))
        px(img, sx-1, sy, (255,220,50,200))
        px(img, sx+1, sy, (255,220,50,200))
        px(img, sx, sy-1, (255,220,50,200))
        px(img, sx, sy+1, (255,220,50,200))
    # 嘴巴（微笑）
    for x in range(16,22):
        px(img, x, 34+(x-16)//3, (200,80,120,255))
    outline_all(img, (180,60,100,255))
    return img

# ---- T112 千龍王（金色，龍形）----
def gen_T112():
    img = new_img()
    # 龍身（金色S形）
    fill_circle_shaded(img, 32, 20, 14, (220, 170, 20))
    fill_circle_shaded(img, 28, 38, 12, (200, 150, 15))
    fill_circle_shaded(img, 36, 52, 8, (180, 130, 10))
    # 龍頭裝飾（角）
    for i in range(6):
        px(img, 24-i, 10-i, (255,200,50,255))
        px(img, 40+i, 10-i, (255,200,50,255))
    # 眼睛（紅色龍眼）
    fill_circle(img, 26, 16, 4, (255,50,50,255))
    fill_circle(img, 26, 16, 2, (200,0,0,255))
    px(img, 25, 15, (255,200,200,200))
    # 龍鱗（金色菱形）
    for y in range(18, 50, 5):
        for x in range(22, 44, 6):
            px(img, x, y, (255,200,30,255))
    # 火焰（嘴部）
    for i in range(5):
        px(img, 14-i, 18+i, (255,100,0,255))
        px(img, 13-i, 17+i, (255,200,0,200))
    # 金色光暈
    for r in range(16,20):
        outline_circle(img, 32, 20, r, (255,200,0,30))
    outline_all(img, (120,80,0,255))
    return img

# ---- T113 黃金水母（黃金，透明觸手）----
def gen_T113():
    img = new_img()
    # 傘部（金色半圓）
    for y in range(10, 36):
        for x in range(10, 54):
            if (x-32)**2 + (y-36)**2 <= 22**2 and y <= 36:
                nx_ = (x-32)/22; ny_ = (y-36)/22
                dot = -(nx_*(-0.7)+ny_*(-0.7))
                if dot > 0.25:
                    c = (255,220,50,255)
                elif dot < -0.1:
                    c = (180,140,0,255)
                else:
                    c = (220,180,20,255)
                px(img, x, y, c)
    # 傘邊緣（金色波浪）
    for i in range(8):
        cx = 14 + i*6
        fill_circle(img, cx, 36, 3, (255,200,30,255))
    # 觸手（半透明金色）
    tentacle_x = [16,20,24,28,32,36,40,44,48]
    for i,tx in enumerate(tentacle_x):
        length = 12 + (i%3)*4
        for j in range(length):
            wave = 2 if (i+j)%4<2 else -2
            alpha = max(80, 200-j*10)
            px(img, tx+wave, 38+j, (255,200,30,alpha))
    # 中心發光
    fill_circle(img, 32, 24, 6, (255,240,100,200))
    fill_circle(img, 32, 24, 3, (255,255,200,255))
    outline_all(img, (150,100,0,255))
    return img

# ---- T114 雷霆龍蝦（橙電，帶閃電）----
def gen_T114():
    img = new_img()
    # 身體（橙色）
    fill_circle_shaded(img, 32, 36, 17, (230, 100, 20))
    # 閃電紋路
    lightning_pts = [(32,20),(28,26),(34,28),(28,34),(36,36),(30,42),(36,44)]
    for i in range(len(lightning_pts)-1):
        x1,y1 = lightning_pts[i]; x2,y2 = lightning_pts[i+1]
        steps = max(abs(x2-x1),abs(y2-y1))
        for s in range(steps+1):
            t = s/max(steps,1)
            px(img, int(x1+t*(x2-x1)), int(y1+t*(y2-y1)), (255,255,100,255))
    # 眼睛（電藍色）
    fill_circle(img, 24, 30, 4, (100,200,255,255))
    fill_circle(img, 24, 30, 2, (0,100,255,255))
    # 觸角（帶電）
    for i in range(10):
        px(img, 20-i, 22-i, (255,200,50,255))
        px(img, 22-i, 20-i, (255,200,50,255))
    # 大爪（帶電光）
    fill_circle_shaded(img, 14, 38, 8, (200,80,15))
    fill_circle(img, 14, 38, 3, (255,255,100,200))
    fill_circle_shaded(img, 50, 38, 7, (200,80,15))
    # 電光粒子
    for ex,ey in [(10,20),(52,24),(8,44),(54,48),(32,10)]:
        px(img, ex, ey, (255,255,100,255))
        px(img, ex+1, ey, (255,255,100,200))
        px(img, ex, ey+1, (255,255,100,200))
    outline_all(img, (100,40,0,255))
    return img

# ---- T115 彩虹鳳凰（彩虹色，翅膀）----
def gen_T115():
    img = new_img()
    # 翅膀（左右展開，彩虹漸層）
    rainbow = [(255,0,0),(255,128,0),(255,255,0),(0,200,0),(0,100,255),(150,0,255)]
    for i,c in enumerate(rainbow):
        # 左翼
        for j in range(8):
            px(img, 4+i*3, 20+j+i*2, c+(255,))
            px(img, 5+i*3, 20+j+i*2, c+(200,))
        # 右翼
        for j in range(8):
            px(img, 60-i*3, 20+j+i*2, c+(255,))
            px(img, 59-i*3, 20+j+i*2, c+(200,))
    # 身體（金色）
    fill_circle_shaded(img, 32, 32, 12, (255, 200, 50))
    # 頭部（橙色）
    fill_circle_shaded(img, 32, 20, 8, (255, 150, 30))
    # 眼睛
    fill_circle(img, 28, 18, 3, (255,255,255,255))
    fill_circle(img, 28, 18, 2, (255,50,50,255))
    # 鳳冠（彩虹）
    for i,c in enumerate(rainbow[:4]):
        px(img, 30+i, 10-i, c+(255,))
        px(img, 30+i, 11-i, c+(200,))
    # 尾羽
    for i,c in enumerate(rainbow):
        px(img, 28+i*2, 46+i, c+(255,))
        px(img, 29+i*2, 47+i, c+(200,))
    outline_all(img, (150,80,0,255))
    return img

# ---- T116 吸血鬼（深紅，蝙蝠翅膀）----
def gen_T116():
    img = new_img()
    # 蝙蝠翅膀（左右）
    wing_pts_l = [(4,20),(8,14),(14,18),(18,24),(22,20),(26,26)]
    wing_pts_r = [(60,20),(56,14),(50,18),(46,24),(42,20),(38,26)]
    for pts in [wing_pts_l, wing_pts_r]:
        for i in range(len(pts)-1):
            x1,y1=pts[i]; x2,y2=pts[i+1]
            steps=max(abs(x2-x1),abs(y2-y1))
            for s in range(steps+1):
                t=s/max(steps,1)
                px(img,int(x1+t*(x2-x1)),int(y1+t*(y2-y1)),(80,0,0,255))
    # 翅膀填充
    for y in range(14,28):
        for x in range(6,28):
            if img.getpixel((x,y))[3]>0:
                fill_rect(img,x,y,x+1,y+1,(60,0,0,200))
        for x in range(36,58):
            if img.getpixel((x,y))[3]>0:
                fill_rect(img,x,y,x+1,y+1,(60,0,0,200))
    # 身體（深紅色）
    fill_circle_shaded(img, 32, 36, 14, (160, 20, 20))
    # 披風
    fill_rect_shaded(img, 22, 36, 42, 54, (100, 0, 0))
    # 臉（蒼白）
    fill_circle_shaded(img, 32, 28, 10, (220, 200, 200))
    # 眼睛（紅色）
    fill_circle(img, 27, 26, 3, (255,0,0,255))
    fill_circle(img, 37, 26, 3, (255,0,0,255))
    fill_circle(img, 27, 26, 1, (200,0,0,255))
    fill_circle(img, 37, 26, 1, (200,0,0,255))
    # 尖牙
    px(img, 30, 34, (255,255,255,255))
    px(img, 31, 35, (255,255,255,255))
    px(img, 34, 34, (255,255,255,255))
    px(img, 33, 35, (255,255,255,255))
    # 血月光暈
    for r in range(16,20):
        outline_circle(img, 32, 28, r, (200,0,0,30))
    outline_all(img, (40,0,0,255))
    return img

# ---- T117 水晶龍（紫水晶，龍形）----
def gen_T117():
    img = new_img()
    # 水晶光暈
    for r in range(24,20,-1):
        alpha = int(40*(24-r)/4)
        outline_circle(img, 32, 28, r, (180,100,255,alpha))
    # 龍身（紫水晶色）
    fill_circle_shaded(img, 32, 24, 16, (140, 80, 220))
    fill_circle_shaded(img, 28, 42, 10, (120, 60, 200))
    # 水晶鱗片（菱形）
    for y in range(16, 50, 5):
        for x in range(20, 46, 6):
            px(img, x, y, (200,150,255,255))
            px(img, x+1, y, (180,120,240,255))
    # 龍角（水晶）
    for i in range(7):
        px(img, 24-i, 12-i, (180,120,255,255))
        px(img, 40+i, 12-i, (180,120,255,255))
    # 眼睛（紫色發光）
    fill_circle(img, 26, 20, 4, (220,180,255,255))
    fill_circle(img, 26, 20, 2, (150,50,255,255))
    px(img, 25, 19, (255,255,255,200))
    # 水晶碎片裝飾
    for cx,cy in [(44,16),(48,28),(46,40),(20,44)]:
        fill_circle(img, cx, cy, 3, (200,150,255,255))
        px(img, cx, cy-3, (255,220,255,255))
    # 嘴（噴出水晶氣息）
    for i in range(6):
        px(img, 14-i, 22+i, (180,120,255,200))
        px(img, 13-i, 21+i, (220,180,255,150))
    outline_all(img, (60,0,120,255))
    return img

# ---- 主程式 ----
GENERATORS = {
    "T106": ("drill_lobster", gen_T106),
    "T107": ("bomb_crab",     gen_T107),
    "T108": ("mega_octopus",  gen_T108),
    "T109": ("anglerfish",    gen_T109),
    "T110": ("crocodile",     gen_T110),
    "T111": ("prize_fish",    gen_T111),
    "T112": ("chainlong",     gen_T112),
    "T113": ("jellyfish",     gen_T113),
    "T114": ("thunder_lobster", gen_T114),
    "T115": ("rainbow_phoenix", gen_T115),
    "T116": ("vampire",       gen_T116),
    "T117": ("crystal_dragon", gen_T117),
}

def count_pixels(img):
    count = 0
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x,y))[3] > 10:
                count += 1
    return count

if __name__ == "__main__":
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    print(f"生成 T106-T117 目標物 Sprite...")
    for tid, (name, gen_fn) in GENERATORS.items():
        img = gen_fn()
        filename = f"{tid}_{name}.png"
        path = os.path.join(OUTPUT_DIR, filename)
        img.save(path)
        density = count_pixels(img) / (SIZE*SIZE) * 100
        print(f"  {filename}: {density:.1f}% 密度")
    print(f"\n完成！共生成 {len(GENERATORS)} 個目標物 Sprite")
    print(f"輸出目錄：{OUTPUT_DIR}")
