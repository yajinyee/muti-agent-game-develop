# -*- coding: utf-8 -*-
"""
修復低密度目標物：T106 鑽頭龍蝦（29.5%）和 T115 彩虹鳳凰（25.2%）
目標：密度 >= 35%
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
            nx_ = (x-cx)/max(r,1); ny_ = (y-cy)/max(r,1)
            dot = -(nx_*(-0.7)+ny_*(-0.7))
            if dot > 0.25:
                c = (min(255,r_v+40), min(255,g_v+40), min(255,b_v+40), 255)
            elif dot < -0.1:
                c = (max(0,r_v-45), max(0,g_v-45), max(0,b_v-45), 255)
            else:
                c = (r_v, g_v, b_v, 255)
            px(img, x, y, c)

def fill_rect(img, x1, y1, x2, y2, color):
    for y in range(max(0,y1), min(SIZE,y2)):
        for x in range(max(0,x1), min(SIZE,x2)):
            px(img, x, y, color)

def outline_all(img, color=(20,20,20,255)):
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

def count_pixels(img):
    count = 0
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x,y))[3] > 10:
                count += 1
    return count

# ---- T106 鑽頭龍蝦 v2（更飽滿）----
def gen_T106_v2():
    img = new_img()
    # 身體（更大的橙紅色橢圓）
    fill_circle_shaded(img, 30, 36, 18, (220, 80, 20))
    # 腹部（橙色）
    fill_circle_shaded(img, 30, 42, 14, (200, 100, 30))
    # 鑽頭（右側，更粗）
    for i in range(16):
        half = max(1, 8 - i//2)
        for y in range(36-half, 36+half+1):
            px(img, 48+i, y, (180,60,10,255))
    # 鑽頭尖端
    fill_circle_shaded(img, 58, 36, 5, (160, 160, 170))
    # 眼睛（大眼）
    fill_circle(img, 22, 28, 5, (255,255,255,255))
    fill_circle(img, 22, 28, 3, (0,150,255,255))
    fill_circle(img, 22, 28, 1, (0,0,0,255))
    px(img, 21, 27, (255,255,255,200))
    # 觸角（長）
    for i in range(12):
        px(img, 18-i, 22-i, (200,60,10,255))
        px(img, 20-i, 18-i, (200,60,10,255))
    # 大爪（左側，更大）
    fill_circle_shaded(img, 12, 40, 9, (200,70,15))
    fill_circle_shaded(img, 8, 32, 7, (200,70,15))
    # 小爪（右下）
    fill_circle_shaded(img, 16, 52, 6, (190,65,12))
    # 腳
    for i in range(4):
        px(img, 20+i*4, 54, (180,60,10,255))
        px(img, 20+i*4, 55, (160,50,8,255))
    # 鑽頭螺旋紋
    for i in range(4):
        px(img, 50+i*3, 32, (220,100,30,255))
        px(img, 51+i*3, 40, (220,100,30,255))
    outline_all(img, (80,20,0,255))
    return img

# ---- T115 彩虹鳳凰 v2（更飽滿）----
def gen_T115_v2():
    img = new_img()
    rainbow = [(255,0,0),(255,128,0),(255,255,0),(0,200,0),(0,100,255),(150,0,255)]
    # 翅膀（更大更飽滿）
    for i,c in enumerate(rainbow):
        # 左翼（更寬）
        for j in range(12):
            for k in range(3):
                px(img, 2+i*4+k, 16+j+i*2, c+(255,))
        # 右翼（更寬）
        for j in range(12):
            for k in range(3):
                px(img, 62-i*4-k, 16+j+i*2, c+(255,))
    # 翅膀填充（讓翅膀更實心）
    for y in range(16, 40):
        for x in range(4, 28):
            if img.getpixel((x,y))[3] == 0:
                # 找最近的彩虹色
                ci = min(5, (x-4)//4)
                c = rainbow[ci]
                px(img, x, y, c+(120,))
        for x in range(36, 60):
            if img.getpixel((x,y))[3] == 0:
                ci = min(5, (60-x)//4)
                c = rainbow[ci]
                px(img, x, y, c+(120,))
    # 身體（金色，更大）
    fill_circle_shaded(img, 32, 36, 14, (255, 200, 50))
    # 頭部（橙色）
    fill_circle_shaded(img, 32, 22, 10, (255, 150, 30))
    # 眼睛
    fill_circle(img, 27, 20, 4, (255,255,255,255))
    fill_circle(img, 27, 20, 2, (255,50,50,255))
    fill_circle(img, 27, 20, 1, (0,0,0,255))
    px(img, 26, 19, (255,255,255,200))
    # 鳳冠（彩虹，更高）
    for i,c in enumerate(rainbow[:5]):
        for j in range(3):
            px(img, 28+i*2, 10-i-j, c+(255,))
    # 尾羽（更長更寬）
    for i,c in enumerate(rainbow):
        for k in range(3):
            px(img, 26+i*2+k, 50+i*2, c+(255,))
            px(img, 26+i*2+k, 51+i*2, c+(200,))
    # 胸部羽毛
    for i,c in enumerate(rainbow[:4]):
        fill_circle(img, 32, 36+i*3, 3, c+(180,))
    outline_all(img, (150,80,0,255))
    return img

if __name__ == "__main__":
    # T106
    img106 = gen_T106_v2()
    path106 = os.path.join(OUTPUT_DIR, "T106_drill_lobster.png")
    img106.save(path106)
    d106 = count_pixels(img106)/(SIZE*SIZE)*100
    print(f"T106_drill_lobster.png: {d106:.1f}% 密度")

    # T115
    img115 = gen_T115_v2()
    path115 = os.path.join(OUTPUT_DIR, "T115_rainbow_phoenix.png")
    img115.save(path115)
    d115 = count_pixels(img115)/(SIZE*SIZE)*100
    print(f"T115_rainbow_phoenix.png: {d115:.1f}% 密度")

    print("修復完成！")
