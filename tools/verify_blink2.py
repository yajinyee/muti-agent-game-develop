# -*- coding: utf-8 -*-
"""更精確的眨眼驗證：比較相鄰幀的差異"""
import numpy as np
from PIL import Image

SHEET_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\sheets"
FRAME_SIZE = 96

for char in ["chiikawa"]:
    sheet = Image.open(f"{SHEET_DIR}/{char}_animated.png").convert("RGBA")
    
    print(f"=== {char} idle 幀分析 ===")
    frames = []
    for i in range(8):
        f = sheet.crop((i * FRAME_SIZE, 0, (i+1) * FRAME_SIZE, FRAME_SIZE))
        frames.append(np.array(f))
    
    # 分析每幀的非透明像素數
    for i, arr in enumerate(frames):
        non_trans = np.sum(arr[:,:,3] > 10)
        # 找眼睛區域（y=28-45）的深色像素
        eye_region = arr[28:45, 20:75]
        dark_pixels = np.sum(
            (eye_region[:,:,0] < 80) & 
            (eye_region[:,:,1] < 80) & 
            (eye_region[:,:,2] < 80) & 
            (eye_region[:,:,3] > 50)
        )
        print(f"  幀{i}: 非透明={non_trans}, 眼睛深色像素={dark_pixels}")
    
    print()
    print("  預期：幀4,5,6 的眼睛深色像素應該比幀0少（眼睛閉合）")
