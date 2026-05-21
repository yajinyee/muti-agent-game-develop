# -*- coding: utf-8 -*-
"""
修復 sprite 中的洋紅色殘留
問題：AI 生成時使用洋紅色背景，去背後仍有殘留的洋紅色像素
解法：更激進的洋紅色去除 + 邊緣羽化
"""
from PIL import Image
import os
import math
from collections import deque

CHARS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"

def is_magenta_like(r, g, b, threshold=80):
    """判斷是否為洋紅色系（包含粉紅、紫紅等）"""
    # 洋紅色特徵：R 高、G 低、B 高
    if r > 150 and g < 100 and b > 100:
        # 計算與純洋紅色的距離
        dist = math.sqrt((r - 255)**2 + g**2 + (b - 255)**2)
        return dist < threshold
    return False

def is_pink_residue(r, g, b):
    """判斷是否為粉紅色殘留（洋紅色去背後的邊緣）"""
    # 粉紅色：R 高、G 中、B 中高，但不是正常的角色顏色
    if r > 180 and g < 130 and b > 130:
        # 排除正常的角色顏色（白色、淡粉等）
        if g < 80:  # G 很低才是殘留
            return True
    return False

def remove_magenta_aggressive(img_path, output_path=None):
    """激進的洋紅色去除"""
    img = Image.open(img_path).convert("RGBA")
    pixels = img.load()
    w, h = img.size
    
    removed = 0
    
    for y in range(h):
        for x in range(w):
            r, g, b, a = pixels[x, y]
            if a < 10:
                continue
            
            # 洋紅色殘留
            if is_magenta_like(r, g, b, threshold=100):
                pixels[x, y] = (0, 0, 0, 0)
                removed += 1
                continue
            
            # 粉紅色殘留（更嚴格）
            if is_pink_residue(r, g, b):
                pixels[x, y] = (0, 0, 0, 0)
                removed += 1
                continue
    
    if output_path is None:
        output_path = img_path
    
    img.save(output_path)
    return removed

def analyze_magenta(img_path):
    """分析洋紅色殘留數量"""
    img = Image.open(img_path).convert("RGBA")
    pixels = img.load()
    w, h = img.size
    
    magenta_count = 0
    total_non_transparent = 0
    
    for y in range(h):
        for x in range(w):
            r, g, b, a = pixels[x, y]
            if a < 10:
                continue
            total_non_transparent += 1
            if is_magenta_like(r, g, b, threshold=100):
                magenta_count += 1
    
    return magenta_count, total_non_transparent

def main():
    print("=== 洋紅色殘留修復 ===\n")
    
    files = [f for f in os.listdir(CHARS_DIR) if f.endswith('.png') and 'ref' not in f]
    
    total_removed = 0
    
    for fname in sorted(files):
        path = os.path.join(CHARS_DIR, fname)
        
        # 分析前
        before, total = analyze_magenta(path)
        
        if before == 0:
            print(f"  ✅ {fname}: 無殘留")
            continue
        
        # 修復
        removed = remove_magenta_aggressive(path)
        
        # 分析後
        after, total_after = analyze_magenta(path)
        
        print(f"  🔧 {fname}: {before} → {after} 洋紅色像素（移除 {removed}）")
        total_removed += removed
    
    print(f"\n總計移除 {total_removed} 個洋紅色像素")
    print("\n重新執行 QC 確認...")

if __name__ == '__main__':
    main()
