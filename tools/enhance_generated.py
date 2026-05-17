# -*- coding: utf-8 -*-
"""
用 OpenCV 後處理程式生成的角色，提升品質
方法：
1. 先用 generate_full_character.py 生成有正確顏色的角色
2. 用 OpenCV 做邊緣強化和顏色優化
3. 輸出更高品質的像素藝術
"""
import cv2
import numpy as np
from PIL import Image, ImageFilter, ImageEnhance
import os

OUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"

def enhance_pixel_art(img_path, output_path):
    """
    用 OpenCV 後處理像素藝術：
    1. 邊緣強化（讓輪廓更清晰）
    2. 顏色飽和度提升
    3. 對比度增強
    """
    # 用 PIL 讀取（保留透明通道）
    pil_img = Image.open(img_path).convert("RGBA")
    
    # 分離 RGB 和 Alpha
    r, g, b, a = pil_img.split()
    rgb = Image.merge("RGB", (r, g, b))
    
    # 1. 提升顏色飽和度
    enhancer = ImageEnhance.Color(rgb)
    rgb = enhancer.enhance(1.4)
    
    # 2. 提升對比度
    enhancer = ImageEnhance.Contrast(rgb)
    rgb = enhancer.enhance(1.2)
    
    # 3. 用 OpenCV 做邊緣強化
    cv_img = cv2.cvtColor(np.array(rgb), cv2.COLOR_RGB2BGR)
    
    # Unsharp mask（銳化）
    gaussian = cv2.GaussianBlur(cv_img, (3, 3), 0)
    sharpened = cv2.addWeighted(cv_img, 1.5, gaussian, -0.5, 0)
    
    # 轉回 PIL
    rgb_sharp = Image.fromarray(cv2.cvtColor(sharpened, cv2.COLOR_BGR2RGB))
    
    # 合併回 RGBA
    result = Image.merge("RGBA", (*rgb_sharp.split(), a))
    
    # 確保像素感（縮小再放大）
    small = result.resize((result.width // 2, result.height // 2), Image.NEAREST)
    final = small.resize((result.width, result.height), Image.NEAREST)
    
    final.save(output_path)
    return True

def process_all():
    """處理所有角色"""
    chars = ["chiikawa", "hachiware", "usagi"]
    states = ["idle", "attack", "bigwin"]
    
    # 先重新生成完整角色（有正確顏色）
    import subprocess
    print("Step 1: Regenerating full characters with correct colors...")
    result = subprocess.run(["py", "tools/generate_full_character.py"], 
                          capture_output=True, text=True, cwd=r"D:\Kiro")
    print(result.stdout[-200:] if result.stdout else "No output")
    
    print("\nStep 2: Enhancing with OpenCV...")
    for char in chars:
        for state in states:
            path = os.path.join(OUT_DIR, f"{char}_{state}.png")
            if os.path.exists(path):
                enhance_pixel_art(path, path)
                img = Image.open(path)
                # 統計顏色
                pixels = [img.getpixel((x,y)) for y in range(img.height) 
                         for x in range(img.width) if img.getpixel((x,y))[3] > 50]
                from collections import Counter
                top = Counter([(r//30*30, g//30*30, b//30*30) 
                               for r,g,b,a in pixels]).most_common(3)
                colors_str = " | ".join([f"#{r:02X}{g:02X}{b:02X}" for (r,g,b),_ in top])
                print(f"  {char}_{state}: {colors_str}")

if __name__ == "__main__":
    process_all()
    print("\nDone!")
