#!/usr/bin/env python3
"""
修復精靈圖問題：
1. 移除洋紅色殘留 (220, 0, 220) 和類似顏色
2. 重新生成更好看的目標物精靈圖
"""
from PIL import Image, ImageDraw, ImageFilter
import numpy as np
import os
import math

base = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"

def remove_magenta(img_path):
    """移除洋紅色殘留"""
    img = Image.open(img_path).convert("RGBA")
    arr = np.array(img)
    
    # 找洋紅色像素 (R高, G低, B高)
    r, g, b, a = arr[:,:,0], arr[:,:,1], arr[:,:,2], arr[:,:,3]
    magenta_mask = (r > 150) & (g < 100) & (b > 150) & (a > 10)
    
    if np.sum(magenta_mask) > 0:
        arr[magenta_mask, 3] = 0  # 設為透明
        img = Image.fromarray(arr)
        img.save(img_path)
        return np.sum(magenta_mask)
    return 0

def draw_pixel_circle(draw, cx, cy, r, color, outline=None):
    """畫像素風格圓形"""
    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            dist = math.sqrt((x - cx)**2 + (y - cy)**2)
            if dist <= r:
                if outline and dist >= r - 1.5:
                    draw.point((x, y), fill=outline)
                else:
                    draw.point((x, y), fill=color)

def draw_pixel_eye(draw, cx, cy, size=3):
    """畫像素眼睛"""
    # 白色眼白
    for dy in range(-size, size+1):
        for dx in range(-size, size+1):
            if dx*dx + dy*dy <= size*size:
                draw.point((cx+dx, cy+dy), fill=(255, 255, 255, 255))
    # 黑色瞳孔
    s2 = max(1, size-1)
    for dy in range(-s2, s2+1):
        for dx in range(-s2, s2+1):
            if dx*dx + dy*dy <= s2*s2:
                draw.point((cx+dx, cy+dy), fill=(20, 20, 20, 255))
    # 高光
    draw.point((cx-1, cy-1), fill=(255, 255, 255, 255))

def make_grass(size=64):
    """重新生成像素雜草 T001"""
    img = Image.new("RGBA", (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    cx, cy = size//2, size//2
    
    # 草的莖（深綠色）
    stem_color = (34, 139, 34, 255)
    leaf_color = (50, 205, 50, 255)
    dark_leaf = (0, 100, 0, 255)
    
    # 主莖
    for y in range(cy+10, cy+22):
        draw.point((cx, y), fill=stem_color)
        draw.point((cx-1, y), fill=stem_color)
    
    # 左葉
    for i in range(12):
        x = cx - 2 - i//2
        y = cy + 8 - i
        draw.point((x, y), fill=leaf_color)
        draw.point((x-1, y), fill=dark_leaf)
    
    # 右葉
    for i in range(14):
        x = cx + 2 + i//2
        y = cy + 6 - i
        draw.point((x, y), fill=leaf_color)
        draw.point((x+1, y), fill=dark_leaf)
    
    # 中間葉（最高）
    for i in range(16):
        x = cx
        y = cy - 2 - i
        w = max(1, 3 - i//4)
        for dx in range(-w, w+1):
            draw.point((x+dx, y), fill=leaf_color)
    
    # 小眼睛（讓它有生命感）
    draw_pixel_eye(draw, cx-3, cy-2, 2)
    draw_pixel_eye(draw, cx+3, cy-2, 2)
    
    # 小嘴巴
    draw.point((cx-1, cy+2), fill=(20, 20, 20, 255))
    draw.point((cx, cy+3), fill=(20, 20, 20, 255))
    draw.point((cx+1, cy+2), fill=(20, 20, 20, 255))
    
    return img

def make_bug(size=64, color=(50, 200, 50)):
    """重新生成像素小蟲"""
    img = Image.new("RGBA", (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    cx, cy = size//2, size//2
    
    r, g, b = color
    body_color = (r, g, b, 255)
    dark_color = (max(0, r-60), max(0, g-60), max(0, b-60), 255)
    
    # 身體（橢圓）
    for dy in range(-8, 9):
        for dx in range(-12, 13):
            if (dx/12)**2 + (dy/8)**2 <= 1:
                dist = math.sqrt((dx/12)**2 + (dy/8)**2)
                if dist > 0.85:
                    draw.point((cx+dx, cy+dy), fill=dark_color)
                else:
                    draw.point((cx+dx, cy+dy), fill=body_color)
    
    # 頭部
    for dy in range(-6, 7):
        for dx in range(-6, 7):
            if dx*dx + dy*dy <= 36:
                draw.point((cx+dx+10, cy+dy), fill=body_color)
    
    # 眼睛
    draw_pixel_eye(draw, cx+12, cy-2, 2)
    draw_pixel_eye(draw, cx+8, cy-2, 2)
    
    # 觸角
    for i in range(5):
        draw.point((cx+14+i, cy-6-i), fill=dark_color)
        draw.point((cx+10+i, cy-6-i), fill=dark_color)
    
    # 腳
    for i in range(3):
        y_off = -4 + i*4
        draw.point((cx-8, cy+y_off), fill=dark_color)
        draw.point((cx-10, cy+y_off+1), fill=dark_color)
        draw.point((cx+2, cy+y_off), fill=dark_color)
        draw.point((cx+4, cy+y_off+1), fill=dark_color)
    
    return img

def make_pudding(size=64):
    """重新生成布丁"""
    img = Image.new("RGBA", (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    cx, cy = size//2, size//2 + 4
    
    # 布丁主體（黃色圓形）
    body_color = (255, 220, 50, 255)
    dark_body = (200, 160, 20, 255)
    caramel = (180, 100, 20, 255)
    
    # 主體
    for dy in range(-14, 15):
        for dx in range(-14, 15):
            dist = math.sqrt(dx*dx + dy*dy)
            if dist <= 14:
                if dist > 12:
                    draw.point((cx+dx, cy+dy), fill=dark_body)
                else:
                    draw.point((cx+dx, cy+dy), fill=body_color)
    
    # 焦糖頂部
    for dy in range(-16, -12):
        for dx in range(-8, 9):
            if dx*dx + (dy+14)**2 <= 25:
                draw.point((cx+dx, cy+dy), fill=caramel)
    
    # 眼睛
    draw_pixel_eye(draw, cx-4, cy-2, 3)
    draw_pixel_eye(draw, cx+4, cy-2, 3)
    
    # 笑臉
    for dx in range(-4, 5):
        if abs(dx) >= 2:
            draw.point((cx+dx, cy+5), fill=(20, 20, 20, 255))
    draw.point((cx-3, cy+6), fill=(20, 20, 20, 255))
    draw.point((cx+3, cy+6), fill=(20, 20, 20, 255))
    
    # 腳
    for i in range(4):
        draw.point((cx-6+i*4, cy+14), fill=dark_body)
        draw.point((cx-6+i*4, cy+15), fill=dark_body)
    
    return img

def make_mushroom(size=64):
    """重新生成蘑菇"""
    img = Image.new("RGBA", (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    cx, cy = size//2, size//2
    
    cap_color = (220, 50, 50, 255)
    cap_dark = (160, 20, 20, 255)
    stem_color = (240, 220, 180, 255)
    stem_dark = (200, 170, 130, 255)
    spot_color = (255, 255, 255, 200)
    
    # 莖
    for dy in range(4, 18):
        for dx in range(-7, 8):
            if abs(dx) <= 7 - abs(dy-11)//3:
                if abs(dx) >= 6:
                    draw.point((cx+dx, cy+dy), fill=stem_dark)
                else:
                    draw.point((cx+dx, cy+dy), fill=stem_color)
    
    # 傘蓋
    for dy in range(-18, 6):
        r = min(18, int(18 * math.sqrt(max(0, 1 - (dy/18)**2 * 0.3))))
        for dx in range(-r, r+1):
            dist = math.sqrt(dx*dx + (dy*1.2)**2)
            if dist <= 18:
                if dist > 16:
                    draw.point((cx+dx, cy+dy), fill=cap_dark)
                else:
                    draw.point((cx+dx, cy+dy), fill=cap_color)
    
    # 白色斑點
    for (sx, sy) in [(-6, -10), (4, -14), (-2, -6), (8, -8)]:
        for dy in range(-3, 4):
            for dx in range(-3, 4):
                if dx*dx + dy*dy <= 9:
                    draw.point((cx+sx+dx, cy+sy+dy), fill=spot_color)
    
    # 眼睛
    draw_pixel_eye(draw, cx-3, cy+2, 2)
    draw_pixel_eye(draw, cx+3, cy+2, 2)
    
    return img

# 修復洋紅色殘留
print("=== Removing magenta residuals ===")
for fname in os.listdir(base):
    if fname.endswith(".png") and not fname.endswith("_swim.png") and not fname.endswith("_backup.png"):
        path = os.path.join(base, fname)
        count = remove_magenta(path)
        if count > 0:
            print(f"  {fname}: removed {count} magenta pixels")

# 重新生成基礎目標物
print("\n=== Regenerating basic targets ===")

grass = make_grass()
grass.save(os.path.join(base, "T001_grass.png"))
print("  T001_grass.png regenerated")

bug_g = make_bug(color=(50, 200, 50))
bug_g.save(os.path.join(base, "T002_bug_g.png"))
print("  T002_bug_g.png regenerated")

bug_r = make_bug(color=(220, 60, 60))
bug_r.save(os.path.join(base, "T003_bug_r.png"))
print("  T003_bug_r.png regenerated")

bug_b = make_bug(color=(60, 100, 220))
bug_b.save(os.path.join(base, "T004_bug_b.png"))
print("  T004_bug_b.png regenerated")

pudding = make_pudding()
pudding.save(os.path.join(base, "T005_pudding.png"))
print("  T005_pudding.png regenerated")

mushroom = make_mushroom()
mushroom.save(os.path.join(base, "T006_mushroom.png"))
print("  T006_mushroom.png regenerated")

print("\nDone!")
