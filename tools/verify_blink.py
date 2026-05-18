# -*- coding: utf-8 -*-
"""驗證眨眼效果：比較幀 0（全開）和幀 5（全閉）的差異"""
import numpy as np
from PIL import Image

SHEET_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\sheets"
FRAME_SIZE = 96

for char in ["chiikawa", "hachiware", "usagi"]:
    sheet = Image.open(f"{SHEET_DIR}/{char}_animated.png").convert("RGBA")
    
    # 取幀 0（全開）和幀 5（全閉）
    frame0 = sheet.crop((0 * FRAME_SIZE, 0, 1 * FRAME_SIZE, FRAME_SIZE))
    frame5 = sheet.crop((5 * FRAME_SIZE, 0, 6 * FRAME_SIZE, FRAME_SIZE))
    
    arr0 = np.array(frame0)
    arr5 = np.array(frame5)
    
    # 計算差異
    diff = np.abs(arr0.astype(int) - arr5.astype(int))
    changed_pixels = np.sum(np.any(diff > 10, axis=2))
    
    print(f"{char}:")
    print(f"  幀0 vs 幀5 差異像素數: {changed_pixels}")
    
    # 找差異最大的區域（眼睛位置）
    diff_mask = np.any(diff > 10, axis=2)
    ys, xs = np.where(diff_mask)
    if len(ys) > 0:
        print(f"  差異區域: y={ys.min()}-{ys.max()}, x={xs.min()}-{xs.max()}")
        print(f"  ✅ 眨眼效果確認（差異在眼睛區域）" if ys.min() < 50 else "  ⚠️ 差異位置異常")
    else:
        print(f"  ❌ 無差異（眨眼效果未生效）")
    print()
