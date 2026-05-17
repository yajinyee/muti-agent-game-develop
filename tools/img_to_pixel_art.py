# -*- coding: utf-8 -*-
"""
圖片轉像素藝術 - 使用 OpenCV + K-means
突破方法：把任何圖片轉成高品質像素藝術
比手工畫像素品質高很多
"""
import cv2
import numpy as np
from PIL import Image
import os

def img_to_pixel_art(input_path, output_path, pixel_size=8, n_colors=16, output_size=(64, 64)):
    """
    把圖片轉成像素藝術
    1. 縮小到目標像素尺寸
    2. K-means 顏色量化（限制顏色數）
    3. 放大回輸出尺寸（NEAREST 插值）
    """
    # 讀取圖片
    img = cv2.imread(input_path, cv2.IMREAD_UNCHANGED)
    if img is None:
        print(f"  Cannot read: {input_path}")
        return False
    
    # 處理透明通道
    has_alpha = img.shape[2] == 4 if len(img.shape) == 3 else False
    if has_alpha:
        alpha = img[:, :, 3]
        img_rgb = img[:, :, :3]
    else:
        alpha = None
        img_rgb = img
    
    # 縮小到像素藝術尺寸
    target_w = output_size[0] // pixel_size
    target_h = output_size[1] // pixel_size
    small = cv2.resize(img_rgb, (target_w, target_h), interpolation=cv2.INTER_AREA)
    
    # K-means 顏色量化
    pixels = small.reshape(-1, 3).astype(np.float32)
    criteria = (cv2.TERM_CRITERIA_EPS + cv2.TERM_CRITERIA_MAX_ITER, 20, 1.0)
    _, labels, centers = cv2.kmeans(pixels, n_colors, None, criteria, 10, cv2.KMEANS_RANDOM_CENTERS)
    centers = np.uint8(centers)
    quantized = centers[labels.flatten()].reshape(small.shape)
    
    # 放大回輸出尺寸（NEAREST 保持像素感）
    result = cv2.resize(quantized, output_size, interpolation=cv2.INTER_NEAREST)
    
    # 處理透明通道
    if has_alpha:
        alpha_small = cv2.resize(alpha, (target_w, target_h), interpolation=cv2.INTER_AREA)
        # 二值化透明通道
        _, alpha_binary = cv2.threshold(alpha_small, 128, 255, cv2.THRESH_BINARY)
        alpha_large = cv2.resize(alpha_binary, output_size, interpolation=cv2.INTER_NEAREST)
        result_rgba = cv2.cvtColor(result, cv2.COLOR_BGR2BGRA)
        result_rgba[:, :, 3] = alpha_large
        cv2.imwrite(output_path, result_rgba)
    else:
        cv2.imwrite(output_path, result)
    
    return True

def download_and_convert(url, save_name, output_path, pixel_size=6, n_colors=12):
    """下載圖片並轉成像素藝術"""
    import urllib.request
    
    headers = {"User-Agent": "Mozilla/5.0 Chrome/120.0.0.0"}
    temp_path = f"D:\\Kiro\\client\\chiikawa-pixel\\assets\\sprites\\reference\\temp_{save_name}"
    
    try:
        req = urllib.request.Request(url, headers=headers)
        with urllib.request.urlopen(req, timeout=15) as r:
            data = r.read()
        with open(temp_path, "wb") as f:
            f.write(data)
        print(f"  Downloaded: {len(data)} bytes")
        
        success = img_to_pixel_art(temp_path, output_path, pixel_size=pixel_size, n_colors=n_colors)
        if success:
            print(f"  Converted to pixel art: {output_path}")
        return success
    except Exception as e:
        print(f"  Error: {e}")
        return False

def convert_existing_reference(ref_path, output_path, pixel_size=6, n_colors=12):
    """把已有的參考圖轉成像素藝術"""
    # 先用 PIL 處理（去除白色背景）
    from PIL import Image as PILImage
    import math
    from collections import deque
    
    img = PILImage.open(ref_path).convert("RGB")
    w, h = img.size
    pixels = img.load()
    
    # 找圖案區域
    rows, cols = [], []
    for y in range(0, h, 3):
        for x in range(0, w, 3):
            r, g, b = pixels[x, y]
            if not (r > 210 and g > 210 and b > 210):
                rows.append(y)
                cols.append(x)
    
    if not rows:
        return False
    
    x1, y1 = min(cols), min(rows)
    x2, y2 = max(cols)+1, max(rows)+1
    crop = img.crop((x1, y1, x2, y2))
    
    # 儲存裁切後的圖
    temp_path = ref_path.replace(".png", "_crop.png").replace(".jpg", "_crop.jpg")
    crop.save(temp_path)
    
    # 轉像素藝術
    success = img_to_pixel_art(temp_path, output_path, pixel_size=pixel_size, n_colors=n_colors)
    
    # 清理暫存
    if os.path.exists(temp_path):
        os.remove(temp_path)
    
    return success

if __name__ == "__main__":
    ref_dir = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\reference"
    out_dir = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"
    
    from PIL import Image as PILImage
    
    print("Converting reference images to pixel art using OpenCV K-means...")
    
    conversions = [
        ("chiikawa_0.png", "chiikawa"),
        ("hachiware_0.png", "hachiware"),
        ("usagi_ref2_0.png", "usagi"),
    ]
    
    for ref_name, char_name in conversions:
        ref_path = os.path.join(ref_dir, ref_name)
        print(f"\n[{char_name}]")
        
        for state, pixel_sz in [("idle", 6), ("attack", 6), ("bigwin", 6)]:
            out_path = os.path.join(out_dir, f"{char_name}_{state}.png")
            success = convert_existing_reference(ref_path, out_path, pixel_size=pixel_sz, n_colors=14)
            if success:
                # 驗證輸出
                result = PILImage.open(out_path)
                print(f"  {state}: {result.width}x{result.height} OK")
            else:
                print(f"  {state}: FAILED")
    
    print("\nDone!")
