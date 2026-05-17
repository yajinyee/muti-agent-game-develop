# -*- coding: utf-8 -*-
"""
結合最好的方法：
- idle: OpenCV K-means 從參考圖轉換（最高品質）
- attack: idle 旋轉 + 光效
- bigwin: idle 位移 + 金色光效
"""
from PIL import Image, ImageEnhance
import os
import math

OUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"

def make_attack(idle_img):
    """攻擊狀態：旋轉 + 右上角劍氣光效"""
    # 旋轉
    rotated = idle_img.rotate(-18, expand=False, fillcolor=(0, 0, 0, 0))
    
    # 加亮度
    enhancer = ImageEnhance.Brightness(rotated)
    bright = enhancer.enhance(1.1)
    
    result = bright.copy()
    pixels = result.load()
    w, h = result.size
    
    # 右上角加劍氣光點（粉紅/藍/黃依角色）
    for i in range(8):
        sx = w - 4 - i * 2
        sy = 4 + i * 2
        if 0 <= sx < w and 0 <= sy < h:
            r, g, b, a = pixels[sx, sy]
            if a < 30:
                pixels[sx, sy] = (255, 200, 230, max(0, 180 - i*20))
    
    return result

def make_bigwin(idle_img):
    """大獎狀態：跳起 + 金色光暈 + 星星"""
    w, h = idle_img.size
    
    # 向上位移
    shifted = Image.new("RGBA", (w, h), (0, 0, 0, 0))
    shifted.paste(idle_img, (0, -6))
    
    # 金色色調
    enhancer = ImageEnhance.Brightness(shifted)
    bright = enhancer.enhance(1.15)
    
    result = bright.copy()
    pixels = result.load()
    
    # 加金色色調
    for y in range(h):
        for x in range(w):
            r, g, b, a = pixels[x, y]
            if a > 50:
                pixels[x, y] = (min(255, r+10), min(255, g+8), max(0, b-5), a)
    
    # 星星
    for sx, sy in [(8, 8), (w-10, 6), (6, h//3), (w-8, h//3)]:
        if 0 <= sx < w and 0 <= sy < h:
            pixels[sx, sy] = (255, 235, 60, 240)
            for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
                nx, ny = sx+dx, sy+dy
                if 0 <= nx < w and 0 <= ny < h:
                    pixels[nx, ny] = (255, 235, 60, 160)
    
    return result

def process_all():
    chars = ["chiikawa", "hachiware", "usagi"]
    
    for char in chars:
        idle_path = os.path.join(OUT_DIR, f"{char}_idle.png")
        if not os.path.exists(idle_path):
            print(f"  SKIP {char}: idle not found")
            continue
        
        idle = Image.open(idle_path).convert("RGBA")
        print(f"[{char}] idle: {idle.width}x{idle.height}")
        
        # attack
        attack = make_attack(idle)
        attack.save(os.path.join(OUT_DIR, f"{char}_attack.png"))
        print(f"  attack: saved")
        
        # bigwin
        bigwin = make_bigwin(idle)
        bigwin.save(os.path.join(OUT_DIR, f"{char}_bigwin.png"))
        print(f"  bigwin: saved")

if __name__ == "__main__":
    print("Combining best methods...")
    process_all()
    print("Done!")
