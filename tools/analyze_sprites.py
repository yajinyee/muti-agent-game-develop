# -*- coding: utf-8 -*-
"""
分析 sprite 的顏色分布，找出品質問題
"""
from PIL import Image
import os
from collections import Counter

def analyze_sprite(path, name):
    img = Image.open(path).convert('RGBA')
    w, h = img.size
    pixels = [(x, y, img.getpixel((x, y))) for y in range(h) for x in range(w)]
    non_transparent = [(x, y, p) for x, y, p in pixels if p[3] > 10]
    
    total = w * h
    non_t = len(non_transparent)
    pct = non_t / total * 100
    
    print(f"\n[{name}] {w}x{h}, {non_t}/{total} ({pct:.0f}%)")
    
    # 顏色分布（不量化，直接看原始顏色）
    colors = Counter()
    for _, _, p in non_transparent:
        # 分類：白色、黑色、紅色、藍色、其他
        r, g, b, a = p
        if r > 200 and g > 200 and b > 200:
            colors['white/light'] += 1
        elif r < 50 and g < 50 and b < 50:
            colors['black/dark'] += 1
        elif r > 150 and g < 100 and b < 100:
            colors['red'] += 1
        elif r < 100 and g < 100 and b > 150:
            colors['blue'] += 1
        elif r > 150 and g > 150 and b < 100:
            colors['yellow/gold'] += 1
        elif r > 150 and g < 100 and b > 100:
            colors['pink/magenta'] += 1
        else:
            colors['other'] += 1
    
    for color, count in sorted(colors.items(), key=lambda x: -x[1]):
        pct_c = count / non_t * 100
        print(f"  {color}: {count} ({pct_c:.1f}%)")
    
    # 檢查是否有洋紅色殘留（去背不完整）
    magenta_count = sum(1 for _, _, p in non_transparent 
                       if p[0] > 200 and p[1] < 50 and p[2] > 200)
    if magenta_count > 0:
        print(f"  ⚠️  洋紅色殘留: {magenta_count} pixels")
    
    return non_t, total, pct

def main():
    char_dir = r'd:\Kiro\client\chiikawa-pixel\assets\sprites\characters'
    target_dir = r'd:\Kiro\client\chiikawa-pixel\assets\sprites\targets'
    
    print("=== 角色 Sprite 顏色分析 ===")
    for f in sorted(os.listdir(char_dir)):
        if not f.endswith('.png') or 'ref' in f:
            continue
        analyze_sprite(os.path.join(char_dir, f), f)
    
    print("\n=== 目標物 Sprite 顏色分析 ===")
    for f in sorted(os.listdir(target_dir)):
        if not f.endswith('.png'):
            continue
        analyze_sprite(os.path.join(target_dir, f), f)

if __name__ == '__main__':
    main()
