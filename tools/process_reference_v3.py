# -*- coding: utf-8 -*-
"""
處理 Perler Bead Pattern 參考圖 v3
修正：正確去除白色背景，讓透明區域真正透明
"""
from PIL import Image
import os
import math

REF_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\reference"
OUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"

# 官方顏色
OFFICIAL = {
    "chiikawa": {
        "body":    (255, 255, 247),
        "outline": (41,  42,  43),
        "blush":   (239, 165, 201),
    },
    "hachiware": {
        "body":    (255, 255, 247),
        "outline": (41,  42,  43),
        "stripe":  (51,  112, 192),
    },
    "usagi": {
        "body":    (255, 255, 247),
        "outline": (17,  17,  17),
        "eye":     (255, 91,  86),
        "ear_in":  (219, 104, 130),
    },
}

def color_dist(c1, c2):
    return math.sqrt(sum((a-b)**2 for a,b in zip(c1[:3], c2[:3])))

def find_bounds(img):
    """找非白色區域"""
    w, h = img.size
    pixels = img.load()
    rows, cols = [], []
    for y in range(h):
        for x in range(w):
            r, g, b = pixels[x, y][:3]
            if not (r > 210 and g > 210 and b > 210):
                rows.append(y)
                cols.append(x)
    if not rows:
        return 0, 0, w, h
    return min(cols), min(rows), max(cols)+1, max(rows)+1

def remove_white_bg(img, threshold=230):
    """
    去除白色背景，讓白色區域變透明
    保留角色的深色輪廓和彩色部分
    """
    result = img.convert("RGBA")
    pixels = result.load()
    w, h = result.size
    
    for y in range(h):
        for x in range(w):
            r, g, b, a = pixels[x, y]
            # 純白或接近白色 → 透明
            if r > threshold and g > threshold and b > threshold:
                pixels[x, y] = (0, 0, 0, 0)
            # 格線黃灰色 → 透明
            elif r > 180 and g > 180 and b < g - 25:
                pixels[x, y] = (0, 0, 0, 0)
    
    return result

def apply_official_colors(img, char_name):
    """套用官方顏色"""
    colors = OFFICIAL.get(char_name, {})
    result = img.copy()
    pixels = result.load()
    w, h = result.size
    
    for y in range(h):
        for x in range(w):
            r, g, b, a = pixels[x, y]
            if a < 30:
                continue
            
            # 深色 → 官方輪廓
            if r < 100 and g < 100 and b < 100:
                oc = colors.get("outline", (41, 42, 43))
                pixels[x, y] = (*oc, a)
            # 藍色 → 小八條紋
            elif char_name == "hachiware" and b > r + 40 and b > 100:
                oc = colors.get("stripe", (51, 112, 192))
                pixels[x, y] = (*oc, a)
            # 紅色 → 烏薩奇眼睛
            elif char_name == "usagi" and r > 180 and r > g + 60:
                oc = colors.get("eye", (255, 91, 86))
                pixels[x, y] = (*oc, a)
            # 粉紅 → 腮紅
            elif char_name == "chiikawa" and r > 200 and b > 150 and g < r - 20:
                oc = colors.get("blush", (239, 165, 201))
                pixels[x, y] = (*oc, a)
    
    return result

def process(ref_name, char_name, target_size=32):
    path = os.path.join(REF_DIR, ref_name)
    if not os.path.exists(path):
        print(f"  NOT FOUND: {ref_name}")
        return None
    
    img = Image.open(path).convert("RGB")
    
    # 找圖案區域
    x1, y1, x2, y2 = find_bounds(img)
    crop = img.crop((x1, y1, x2, y2))
    
    # 縮小到目標尺寸
    small = crop.resize((target_size, target_size), Image.NEAREST)
    
    # 去除白色背景（讓透明區域真正透明）
    no_bg = remove_white_bg(small, threshold=220)
    
    # 套用官方顏色
    corrected = apply_official_colors(no_bg, char_name)
    
    # 放大到 64x64
    final = corrected.resize((64, 64), Image.NEAREST)
    
    return final

def make_attack(idle):
    """攻擊：旋轉 + 右上角粉紅光點"""
    rotated = idle.rotate(-15, expand=False, fillcolor=(0, 0, 0, 0))
    result = rotated.copy()
    pixels = result.load()
    w, h = result.size
    # 右上角加幾個粉紅光點
    for sx, sy in [(w-8, 8), (w-12, 4), (w-4, 14)]:
        if 0 <= sx < w and 0 <= sy < h:
            pixels[sx, sy] = (255, 160, 210, 200)
            if sx-1 >= 0: pixels[sx-1, sy] = (255, 160, 210, 150)
            if sy-1 >= 0: pixels[sx, sy-1] = (255, 160, 210, 150)
    return result

def make_bigwin(idle):
    """大獎：向上位移 + 金色星星"""
    w, h = idle.size
    shifted = Image.new("RGBA", (w, h), (0, 0, 0, 0))
    shifted.paste(idle, (0, -5))
    result = shifted.copy()
    pixels = result.load()
    # 加星星
    for sx, sy in [(6, 6), (w-8, 4), (4, h//3)]:
        if 0 <= sx < w and 0 <= sy < h:
            pixels[sx, sy] = (255, 230, 50, 230)
            for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
                nx, ny = sx+dx, sy+dy
                if 0 <= nx < w and 0 <= ny < h:
                    pixels[nx, ny] = (255, 230, 50, 150)
    return result

def main():
    print("Processing v3 - correct transparency...")
    
    chars = [
        ("chiikawa_0.png", "chiikawa"),
        ("hachiware_0.png", "hachiware"),
        ("usagi_ref2_0.png", "usagi"),
    ]
    
    for ref_name, char_name in chars:
        print(f"\n[{char_name}]")
        idle = process(ref_name, char_name)
        if idle is None:
            continue
        
        idle.save(os.path.join(OUT_DIR, f"{char_name}_idle.png"))
        make_attack(idle).save(os.path.join(OUT_DIR, f"{char_name}_attack.png"))
        make_bigwin(idle).save(os.path.join(OUT_DIR, f"{char_name}_bigwin.png"))
        
        # 確認透明像素數量
        pixels = idle.load()
        transparent = sum(1 for y in range(64) for x in range(64) 
                         if idle.getpixel((x,y))[3] < 30)
        total = 64 * 64
        print(f"  Transparent: {transparent}/{total} ({transparent*100//total}%)")
        print(f"  Saved all 3 states")
    
    print("\nDone!")

if __name__ == "__main__":
    main()
