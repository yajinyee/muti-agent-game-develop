# -*- coding: utf-8 -*-
"""
處理 Perler Bead Pattern 參考圖 v2
改進：
1. 過濾格線顏色（黃色/灰色格線）
2. 用官方顏色校正（#FFFFF7 白色）
3. 更好的背景去除
"""
from PIL import Image
import os
import math

REF_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\reference"
OUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"

# 官方顏色（來源：color-hex.com Ichikawa Clan Palette）
OFFICIAL_COLORS = {
    "chiikawa": {
        "body":    (255, 255, 247),  # #FFFFF7 接近純白
        "outline": (41,  42,  43),   # #292A2B 深黑
        "blush":   (239, 165, 201),  # #EFA5C9 粉紅腮紅
        "weapon":  (255, 140, 190),  # 討伐棒粉紅
    },
    "hachiware": {
        "body":    (255, 255, 247),  # 白色
        "outline": (41,  42,  43),   # 深黑
        "stripe":  (51,  112, 192),  # #3370C0 藍色條紋
        "weapon":  (100, 150, 240),  # 藍色討伐棒
    },
    "usagi": {
        "body":    (255, 255, 247),  # 白色
        "outline": (17,  17,  17),   # #111111 近黑
        "eye":     (255, 91,  86),   # #FF5B56 紅眼
        "ear_in":  (219, 104, 130),  # #DB6882 耳朵內側
        "weapon":  (255, 215, 30),   # 黃色討伐棒
    },
}

# 格線顏色（需要過濾）
GRID_COLORS = [
    (220, 220, 180),  # 黃灰格線
    (200, 200, 160),
    (230, 230, 190),
    (210, 210, 170),
    (240, 240, 200),
]

def color_distance(c1, c2):
    return math.sqrt(sum((a-b)**2 for a,b in zip(c1[:3], c2[:3])))

def is_grid_color(rgb):
    """判斷是否為格線顏色"""
    r, g, b = rgb[:3]
    # 黃灰色格線特徵：R≈G > B，且都在 180-250 範圍
    if 170 < r < 255 and 170 < g < 255 and b < g - 20:
        return True
    for gc in GRID_COLORS:
        if color_distance(rgb, gc) < 25:
            return True
    return False

def find_pattern_bounds(img, white_threshold=215):
    """找到非白色非格線的密集區域"""
    w, h = img.size
    pixels = img.load()
    
    def is_background(x, y):
        rgb = pixels[x, y][:3]
        r, g, b = rgb
        # 純白背景
        if r > white_threshold and g > white_threshold and b > white_threshold:
            return True
        # 格線顏色
        if is_grid_color(rgb):
            return True
        return False
    
    row_density = [sum(1 for x in range(w) if not is_background(x, y)) for y in range(h)]
    col_density = [sum(1 for y in range(h) if not is_background(x, y)) for x in range(w)]
    
    max_r = max(row_density) if row_density else 1
    max_c = max(col_density) if col_density else 1
    
    dense_rows = [y for y, d in enumerate(row_density) if d > max_r * 0.15]
    dense_cols = [x for x, d in enumerate(col_density) if d > max_c * 0.15]
    
    if not dense_rows or not dense_cols:
        return 0, 0, w, h
    
    return min(dense_cols), min(dense_rows), max(dense_cols), max(dense_rows)

def color_correct(img, char_name):
    """
    顏色校正：把接近官方顏色的像素替換為精確的官方顏色
    """
    colors = OFFICIAL_COLORS.get(char_name, {})
    if not colors:
        return img
    
    result = img.copy()
    pixels = result.load()
    w, h = result.size
    
    for y in range(h):
        for x in range(w):
            r, g, b, a = pixels[x, y]
            if a < 30:
                continue
            
            rgb = (r, g, b)
            
            # 過濾格線顏色（變透明）
            if is_grid_color(rgb):
                pixels[x, y] = (0, 0, 0, 0)
                continue
            
            # 白色/接近白色 → 官方白色
            if r > 220 and g > 220 and b > 220:
                oc = colors.get("body", (255, 255, 247))
                pixels[x, y] = (*oc, a)
                continue
            
            # 深色 → 官方輪廓色
            if r < 80 and g < 80 and b < 80:
                oc = colors.get("outline", (41, 42, 43))
                pixels[x, y] = (*oc, a)
                continue
            
            # 藍色 → 條紋色（小八）
            if char_name == "hachiware" and b > r + 30 and b > g:
                oc = colors.get("stripe", (51, 112, 192))
                pixels[x, y] = (*oc, a)
                continue
            
            # 紅色 → 眼睛色（烏薩奇）
            if char_name == "usagi" and r > 180 and r > g + 50 and r > b + 50:
                oc = colors.get("eye", (255, 91, 86))
                pixels[x, y] = (*oc, a)
                continue
            
            # 粉紅色 → 腮紅/武器
            if r > 200 and b > 150 and g < r - 30:
                if char_name == "chiikawa":
                    oc = colors.get("blush", (239, 165, 201))
                    pixels[x, y] = (*oc, a)
                    continue
    
    return result

def process_character(ref_name, char_name):
    ref_path = os.path.join(REF_DIR, ref_name)
    if not os.path.exists(ref_path):
        print(f"  NOT FOUND: {ref_path}")
        return False
    
    img = Image.open(ref_path).convert("RGBA")
    print(f"  Source: {img.width}x{img.height}")
    
    # 找圖案區域
    x1, y1, x2, y2 = find_pattern_bounds(img.convert("RGB"))
    print(f"  Bounds: ({x1},{y1}) to ({x2},{y2})")
    
    # 裁切
    crop = img.crop((x1, y1, x2, y2))
    
    # 縮小到 32x32
    small = crop.resize((32, 32), Image.NEAREST)
    
    # 顏色量化
    quantized = small.quantize(colors=12, method=Image.Quantize.FASTOCTREE).convert("RGBA")
    
    # 顏色校正
    corrected = color_correct(quantized, char_name)
    
    # 放大到 64x64
    final = corrected.resize((64, 64), Image.NEAREST)
    
    return final

def make_attack_state(idle_img):
    """攻擊狀態：旋轉 + 粉紅光暈"""
    rotated = idle_img.rotate(-12, expand=False, fillcolor=(0, 0, 0, 0))
    result = rotated.copy()
    pixels = result.load()
    w, h = result.size
    # 右上角粉紅光暈
    for y in range(h // 3):
        for x in range(w * 2 // 3, w):
            dist = math.sqrt((x - w)**2 + y**2)
            if dist < w // 3:
                r, g, b, a = pixels[x, y]
                if a < 30:
                    alpha = max(0, int(60 - dist * 2))
                    if alpha > 0:
                        pixels[x, y] = (255, 160, 210, alpha)
    return result

def make_bigwin_state(idle_img):
    """大獎狀態：跳起 + 金色光暈 + 星星"""
    w, h = idle_img.size
    shifted = Image.new("RGBA", (w, h), (0, 0, 0, 0))
    shifted.paste(idle_img, (0, -4))
    
    result = shifted.copy()
    pixels = result.load()
    
    # 金色色調
    for y in range(h):
        for x in range(w):
            r, g, b, a = pixels[x, y]
            if a > 50:
                pixels[x, y] = (min(255, r + 12), min(255, g + 8), max(0, b - 8), a)
    
    # 星星
    for sx, sy in [(w//4, h//6), (w*3//4, h//8), (w//8, h//3)]:
        for dy in range(-2, 3):
            for dx in range(-2, 3):
                if abs(dx) + abs(dy) <= 2:
                    nx, ny = sx + dx, sy + dy
                    if 0 <= nx < w and 0 <= ny < h:
                        pixels[nx, ny] = (255, 235, 80, 220)
    
    return result

def main():
    print("Processing reference images v2 (with color correction)...")
    
    chars = [
        ("chiikawa_0.png", "chiikawa"),
        ("hachiware_0.png", "hachiware"),
        ("usagi_ref2_0.png", "usagi"),
    ]
    
    for ref_name, char_name in chars:
        print(f"\n[{char_name}]")
        idle = process_character(ref_name, char_name)
        if idle is None:
            continue
        
        idle.save(os.path.join(OUT_DIR, f"{char_name}_idle.png"))
        print(f"  idle: saved")
        
        attack = make_attack_state(idle)
        attack.save(os.path.join(OUT_DIR, f"{char_name}_attack.png"))
        print(f"  attack: saved")
        
        bigwin = make_bigwin_state(idle)
        bigwin.save(os.path.join(OUT_DIR, f"{char_name}_bigwin.png"))
        print(f"  bigwin: saved")
    
    print("\nDone!")

if __name__ == "__main__":
    main()
