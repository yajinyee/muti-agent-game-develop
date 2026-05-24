#!/usr/bin/env python3
"""分析精靈圖的實際品質"""
from PIL import Image
import os
import numpy as np

PY = "C:/Users/yajinyee0306/AppData/Local/Programs/Python/Python312/python.exe"
base = r"d:\Kiro\client\chiikawa-pixel\assets\sprites"

def analyze_sprite(path):
    img = Image.open(path).convert("RGBA")
    arr = np.array(img)
    
    # 計算非透明像素比例
    alpha = arr[:, :, 3]
    total = arr.shape[0] * arr.shape[1]
    non_transparent = np.sum(alpha > 10)
    ratio = non_transparent / total * 100
    
    # 計算顏色多樣性
    rgb = arr[alpha > 10, :3]
    if len(rgb) > 0:
        unique_colors = len(np.unique(rgb.reshape(-1, 3), axis=0))
    else:
        unique_colors = 0
    
    # 計算主要顏色
    if len(rgb) > 0:
        # 找最常見的顏色
        from collections import Counter
        color_counts = Counter(map(tuple, rgb.tolist()))
        top_colors = color_counts.most_common(3)
    else:
        top_colors = []
    
    return {
        "size": img.size,
        "non_transparent_pct": round(ratio, 1),
        "unique_colors": unique_colors,
        "top_colors": top_colors
    }

sprites = [
    ("T001 雜草", "targets/T001_grass.png"),
    ("T002 綠蟲", "targets/T002_bug_g.png"),
    ("T005 布丁", "targets/T005_pudding.png"),
    ("T101 擬態", "targets/T101_mimic.png"),
    ("T105 金幣魚", "targets/T105_coin_fish.png"),
    ("B001 BOSS", "targets/B001_boss.png"),
    ("吉伊卡哇", "characters/chiikawa_idle.png"),
    ("小八", "characters/hachiware_idle.png"),
    ("烏薩奇", "characters/usagi_idle.png"),
]

for name, rel_path in sprites:
    path = os.path.join(base, rel_path)
    if os.path.exists(path):
        info = analyze_sprite(path)
        print(f"\n{name} ({rel_path}):")
        print(f"  Size: {info['size']}")
        print(f"  Non-transparent: {info['non_transparent_pct']}%")
        print(f"  Unique colors: {info['unique_colors']}")
        if info['top_colors']:
            print(f"  Top colors: {info['top_colors'][:2]}")
    else:
        print(f"\n{name}: NOT FOUND")
