# -*- coding: utf-8 -*-
"""
處理 Perler Bead Pattern 參考圖 v4
關鍵修正：用 flood fill 從邊緣去除背景，保留角色內部的白色
"""
from PIL import Image
import os
import math
from collections import deque

REF_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\reference"
OUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"

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

def is_white_like(rgb, threshold=215):
    r, g, b = rgb[:3]
    return r > threshold and g > threshold and b > threshold

def is_grid_color(rgb):
    r, g, b = rgb[:3]
    return r > 180 and g > 180 and b < g - 20

def flood_fill_background(img, threshold=215):
    """
    從四個邊緣做 flood fill，把連通的白色區域標記為背景
    角色內部的白色不會被觸及
    """
    w, h = img.size
    pixels = img.load()
    visited = [[False] * h for _ in range(w)]
    bg_mask = [[False] * h for _ in range(w)]
    
    queue = deque()
    
    # 從四個邊緣開始
    for x in range(w):
        for y in [0, h-1]:
            if not visited[x][y] and is_white_like(pixels[x, y][:3], threshold):
                queue.append((x, y))
                visited[x][y] = True
    for y in range(h):
        for x in [0, w-1]:
            if not visited[x][y] and is_white_like(pixels[x, y][:3], threshold):
                queue.append((x, y))
                visited[x][y] = True
    
    # BFS flood fill
    while queue:
        x, y = queue.popleft()
        bg_mask[x][y] = True
        
        for dx, dy in [(0,1),(0,-1),(1,0),(-1,0)]:
            nx, ny = x+dx, y+dy
            if 0 <= nx < w and 0 <= ny < h and not visited[nx][ny]:
                rgb = pixels[nx, ny][:3]
                if is_white_like(rgb, threshold) or is_grid_color(rgb):
                    visited[nx][ny] = True
                    queue.append((nx, ny))
    
    return bg_mask

def process(ref_name, char_name, target_size=32):
    path = os.path.join(REF_DIR, ref_name)
    if not os.path.exists(path):
        print(f"  NOT FOUND: {ref_name}")
        return None
    
    img = Image.open(path).convert("RGB")
    w, h = img.size
    pixels = img.load()
    
    # 找非白色區域
    rows, cols = [], []
    for y in range(h):
        for x in range(w):
            r, g, b = pixels[x, y]
            if not (r > 210 and g > 210 and b > 210):
                rows.append(y)
                cols.append(x)
    
    if not rows:
        return None
    
    x1, y1 = min(cols), min(rows)
    x2, y2 = max(cols)+1, max(rows)+1
    crop = img.crop((x1, y1, x2, y2))
    
    # 縮小到目標尺寸
    small = crop.resize((target_size, target_size), Image.NEAREST)
    
    # 轉 RGBA
    result = small.convert("RGBA")
    
    # Flood fill 去除背景（從邊緣）
    bg_mask = flood_fill_background(small, threshold=210)
    
    result_pixels = result.load()
    small_pixels = small.load()
    
    for y in range(target_size):
        for x in range(target_size):
            if bg_mask[x][y]:
                result_pixels[x, y] = (0, 0, 0, 0)  # 透明
            else:
                r, g, b = small_pixels[x, y][:3]
                # 格線顏色也去除
                if is_grid_color((r, g, b)):
                    result_pixels[x, y] = (0, 0, 0, 0)
    
    # 套用官方顏色
    colors = OFFICIAL.get(char_name, {})
    for y in range(target_size):
        for x in range(target_size):
            r, g, b, a = result_pixels[x, y]
            if a < 30:
                continue
            
            # 深色 → 輪廓
            if r < 100 and g < 100 and b < 100:
                oc = colors.get("outline", (41, 42, 43))
                result_pixels[x, y] = (*oc, 255)
            # 白色 → 官方白
            elif r > 200 and g > 200 and b > 200:
                oc = colors.get("body", (255, 255, 247))
                result_pixels[x, y] = (*oc, 255)
            # 藍色 → 小八條紋
            elif char_name == "hachiware" and b > r + 40 and b > 100:
                oc = colors.get("stripe", (51, 112, 192))
                result_pixels[x, y] = (*oc, 255)
            # 紅色 → 烏薩奇眼睛
            elif char_name == "usagi" and r > 180 and r > g + 60:
                oc = colors.get("eye", (255, 91, 86))
                result_pixels[x, y] = (*oc, 255)
            # 粉紅 → 腮紅
            elif char_name == "chiikawa" and r > 180 and b > 130 and g < r - 20:
                oc = colors.get("blush", (239, 165, 201))
                result_pixels[x, y] = (*oc, 255)
    
    # 放大到 64x64
    final = result.resize((64, 64), Image.NEAREST)
    return final

def make_attack(idle):
    rotated = idle.rotate(-15, expand=False, fillcolor=(0, 0, 0, 0))
    result = rotated.copy()
    pixels = result.load()
    w, h = result.size
    for sx, sy in [(w-8, 8), (w-12, 4), (w-4, 14)]:
        if 0 <= sx < w and 0 <= sy < h:
            pixels[sx, sy] = (255, 160, 210, 200)
    return result

def make_bigwin(idle):
    w, h = idle.size
    shifted = Image.new("RGBA", (w, h), (0, 0, 0, 0))
    shifted.paste(idle, (0, -5))
    result = shifted.copy()
    pixels = result.load()
    for sx, sy in [(6, 6), (w-8, 4), (4, h//3)]:
        if 0 <= sx < w and 0 <= sy < h:
            pixels[sx, sy] = (255, 230, 50, 230)
            for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
                nx, ny = sx+dx, sy+dy
                if 0 <= nx < w and 0 <= ny < h:
                    pixels[nx, ny] = (255, 230, 50, 150)
    return result

def main():
    print("Processing v4 - flood fill background removal...")
    
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
        
        # 統計
        transparent = sum(1 for y in range(64) for x in range(64)
                         if idle.getpixel((x,y))[3] < 30)
        white_pixels = sum(1 for y in range(64) for x in range(64)
                          if idle.getpixel((x,y))[3] > 200 and
                          idle.getpixel((x,y))[0] > 200)
        print(f"  Transparent: {transparent}/4096 ({transparent*100//4096}%)")
        print(f"  White body: {white_pixels} pixels")
        
        idle.save(os.path.join(OUT_DIR, f"{char_name}_idle.png"))
        make_attack(idle).save(os.path.join(OUT_DIR, f"{char_name}_attack.png"))
        make_bigwin(idle).save(os.path.join(OUT_DIR, f"{char_name}_bigwin.png"))
        print(f"  Saved all 3 states")
    
    print("\nDone!")

if __name__ == "__main__":
    main()
