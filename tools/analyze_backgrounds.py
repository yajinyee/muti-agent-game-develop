# -*- coding: utf-8 -*-
"""分析背景品質，輸出縮圖"""
from PIL import Image
import os

BG_DIR = "D:/Kiro/client/chiikawa-pixel/assets/sprites/backgrounds"
OUT_DIR = "D:/Kiro/tools"

for f in ["sea_bg.png", "boss_bg.png", "bonus_bg.png"]:
    path = os.path.join(BG_DIR, f)
    img = Image.open(path).convert("RGB")
    
    # 縮圖（320x180）
    thumb = img.resize((320, 180), Image.LANCZOS)
    thumb.save(os.path.join(OUT_DIR, f"thumb_{f}"))
    
    # 分析顏色多樣性
    pixels = list(img.getdata())
    unique_colors = len(set(pixels))
    
    # 計算顏色方差（越高越豐富）
    import statistics
    r_vals = [p[0] for p in pixels[::100]]
    g_vals = [p[1] for p in pixels[::100]]
    b_vals = [p[2] for p in pixels[::100]]
    r_std = statistics.stdev(r_vals) if len(r_vals) > 1 else 0
    g_std = statistics.stdev(g_vals) if len(g_vals) > 1 else 0
    b_std = statistics.stdev(b_vals) if len(b_vals) > 1 else 0
    
    print(f"{f}:")
    print(f"  尺寸: {img.size}")
    print(f"  唯一顏色數: {unique_colors:,}")
    print(f"  顏色方差 R={r_std:.1f} G={g_std:.1f} B={b_std:.1f}")
    richness = (r_std + g_std + b_std) / 3
    level = "豐富" if richness > 30 else ("普通" if richness > 15 else "單調")
    print(f"  視覺豐富度: {richness:.1f} → {level}")
    print()
