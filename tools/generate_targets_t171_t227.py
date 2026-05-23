# -*- coding: utf-8 -*-
"""
目標物 T171-T227 批次生成 — Lucky 系列特殊目標物
DAY-109+ 補充：為 T171-T227 生成 Sprite（64x64 像素）
每個目標物用獨特顏色+圖案區分
"""
from PIL import Image, ImageDraw
import os, math

OUTPUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def new_img():
    return Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

def px(img, x, y, c):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), c)

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
                c = (min(255,r_v+45), min(255,g_v+45), min(255,b_v+45), 255)
            elif dot < -0.1:
                c = (max(0,r_v-50), max(0,g_v-50), max(0,b_v-50), 255)
            else:
                c = (r_v, g_v, b_v, 255)
            px(img, x, y, c)

def outline_circle(img, cx, cy, r, color):
    for y in range(max(0,cy-r-2), min(SIZE,cy+r+3)):
        for x in range(max(0,cx-r-2), min(SIZE,cx+r+3)):
            d = math.sqrt((x-cx)**2 + (y-cy)**2)
            if r+0.1 <= d <= r+1.5:
                px(img, x, y, color)

def draw_star(img, cx, cy, r_out, r_in, points, color):
    """畫星形"""
    draw = ImageDraw.Draw(img)
    pts = []
    for i in range(points * 2):
        angle = math.pi * i / points - math.pi / 2
        r = r_out if i % 2 == 0 else r_in
        pts.append((cx + r * math.cos(angle), cy + r * math.sin(angle)))
    draw.polygon(pts, fill=color)

def draw_diamond(img, cx, cy, r, color):
    draw = ImageDraw.Draw(img)
    pts = [(cx, cy-r), (cx+r, cy), (cx, cy+r), (cx-r, cy)]
    draw.polygon(pts, fill=color)

def draw_lightning(img, cx, cy, color):
    """閃電符號"""
    pts = [(cx+4, cy-14), (cx-2, cy-2), (cx+4, cy-2), (cx-4, cy+14), (cx+2, cy+2), (cx-4, cy+2)]
    draw = ImageDraw.Draw(img)
    draw.polygon(pts, fill=color)

def draw_spiral(img, cx, cy, r, color):
    """螺旋點"""
    for i in range(20):
        angle = i * 0.5
        rr = r * i / 20
        x = int(cx + rr * math.cos(angle))
        y = int(cy + rr * math.sin(angle))
        px(img, x, y, color)
        px(img, x+1, y, color)

def draw_cross(img, cx, cy, r, w, color):
    draw = ImageDraw.Draw(img)
    draw.rectangle([cx-w, cy-r, cx+w, cy+r], fill=color)
    draw.rectangle([cx-r, cy-w, cx+r, cy+w], fill=color)

def draw_eye(img, cx, cy, color):
    """眼睛符號"""
    fill_circle_shaded(img, cx, cy, 6, (255,255,255))
    fill_circle_shaded(img, cx, cy, 3, color)
    px(img, cx-1, cy-1, (255,255,255,200))

def add_glow_ring(img, cx, cy, r, color, alpha=80):
    """外圈光暈"""
    for y in range(max(0,cy-r-4), min(SIZE,cy+r+5)):
        for x in range(max(0,cx-r-4), min(SIZE,cx+r+5)):
            d = math.sqrt((x-cx)**2 + (y-cy)**2)
            if r+1 <= d <= r+4:
                fade = int(alpha * (1 - (d-r-1)/3))
                if fade > 0:
                    existing = img.getpixel((x,y))
                    if existing[3] == 0:
                        img.putpixel((x,y), (color[0], color[1], color[2], fade))

def save(img, name):
    path = os.path.join(OUTPUT_DIR, name)
    img.save(path)
    print(f"  saved: {name}")
