# -*- coding: utf-8 -*-
"""
處理下載的 Perler Bead Pattern 參考圖
1. 找到圖案區域
2. 裁切
3. 縮小到 32x32（像素化）
4. 顏色量化（限制到 8 色）
5. 放大到 64x64 輸出
"""
from PIL import Image, ImageFilter
import os

REF_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\reference"
OUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"

def find_pattern_bounds(img, white_threshold=220):
    """找到非白色的密集區域"""
    w, h = img.size
    pixels = img.load()
    
    def is_white(x, y):
        r, g, b = pixels[x, y][:3]
        return r > white_threshold and g > white_threshold and b > white_threshold
    
    # 找行密度
    row_density = []
    for y in range(h):
        count = sum(1 for x in range(w) if not is_white(x, y))
        row_density.append(count)
    
    # 找列密度
    col_density = []
    for x in range(w):
        count = sum(1 for y in range(h) if not is_white(x, y))
        col_density.append(count)
    
    max_row = max(row_density) if row_density else 1
    max_col = max(col_density) if col_density else 1
    
    dense_rows = [y for y, d in enumerate(row_density) if d > max_row * 0.2]
    dense_cols = [x for x, d in enumerate(col_density) if d > max_col * 0.2]
    
    if not dense_rows or not dense_cols:
        return 0, 0, w, h
    
    return min(dense_cols), min(dense_rows), max(dense_cols), max(dense_rows)

def pixelate_and_quantize(img, target_size=32, num_colors=12):
    """縮小到目標尺寸並量化顏色"""
    # 縮小（NEAREST 保持像素感）
    small = img.resize((target_size, target_size), Image.NEAREST)
    
    # 顏色量化（限制顏色數量，更像素藝術）
    # RGBA 圖片用 FASTOCTREE 方法
    quantized = small.quantize(colors=num_colors, method=Image.Quantize.FASTOCTREE)
    result = quantized.convert("RGBA")
    
    return result

def process_character(ref_name, char_name, states=None):
    """處理一個角色的參考圖"""
    if states is None:
        states = ["idle", "attack", "bigwin"]
    
    ref_path = os.path.join(REF_DIR, ref_name)
    if not os.path.exists(ref_path):
        print(f"  NOT FOUND: {ref_path}")
        return False
    
    img = Image.open(ref_path).convert("RGB")
    print(f"  Source: {img.width}x{img.height}")
    
    # 找圖案區域
    x1, y1, x2, y2 = find_pattern_bounds(img)
    print(f"  Pattern bounds: ({x1},{y1}) to ({x2},{y2})")
    
    # 裁切
    crop = img.crop((x1, y1, x2, y2))
    
    # 轉 RGBA
    crop_rgba = crop.convert("RGBA")
    
    # 像素化 + 量化
    pixelated = pixelate_and_quantize(crop_rgba, target_size=32, num_colors=10)
    
    # 放大到 64x64
    final = pixelated.resize((64, 64), Image.NEAREST)
    
    # 儲存所有狀態（暫時用同一張，之後可以分別處理）
    for state in states:
        out_path = os.path.join(OUT_DIR, f"{char_name}_{state}.png")
        final.save(out_path)
        print(f"  Saved: {char_name}_{state}.png")
    
    return True

def main():
    print("Processing reference images...")
    
    # 處理三個角色
    chars = [
        ("chiikawa_0.png", "chiikawa"),
        ("hachiware_0.png", "hachiware"),
        ("usagi_0.png", "usagi"),
    ]
    
    for ref_name, char_name in chars:
        print(f"\n[{char_name}]")
        process_character(ref_name, char_name)
    
    # 驗證輸出
    print("\n=== Output files ===")
    for f in sorted(os.listdir(OUT_DIR)):
        if f.endswith(".png"):
            path = os.path.join(OUT_DIR, f)
            img = Image.open(path)
            print(f"  {f}: {img.width}x{img.height}")

if __name__ == "__main__":
    main()
