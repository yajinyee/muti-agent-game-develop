# -*- coding: utf-8 -*-
"""
從 idle 參考圖生成 attack/bigwin 狀態
用圖像變換模擬不同動作
"""
from PIL import Image, ImageEnhance, ImageFilter
import os

REF_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\reference"
OUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"

def process_to_pixel(ref_path, size=32):
    """把參考圖處理成像素圖"""
    img = Image.open(ref_path).convert("RGB")
    
    # 找圖案區域
    pixels = img.load()
    row_d = [sum(1 for x in range(img.width) if not all(c > 220 for c in pixels[x,y][:3])) for y in range(img.height)]
    col_d = [sum(1 for y in range(img.height) if not all(c > 220 for c in pixels[x,y][:3])) for x in range(img.width)]
    
    max_r, max_c = max(row_d), max(col_d)
    dr = [y for y,d in enumerate(row_d) if d > max_r*0.2]
    dc = [x for x,d in enumerate(col_d) if d > max_c*0.2]
    
    if not dr or not dc:
        return None
    
    crop = img.crop((min(dc), min(dr), max(dc), max(dr))).convert("RGBA")
    small = crop.resize((size, size), Image.NEAREST)
    quantized = small.quantize(colors=10, method=Image.Quantize.FASTOCTREE).convert("RGBA")
    return quantized

def make_attack_state(idle_img):
    """
    攻擊狀態：
    - 整體向右傾斜（模擬揮棒）
    - 亮度提高（興奮感）
    - 加上粉紅色光暈在右上角（劍氣）
    """
    img = idle_img.copy()
    
    # 輕微旋轉（-15度，向右揮）
    rotated = img.rotate(-12, expand=False, fillcolor=(0,0,0,0))
    
    # 亮度提高
    enhancer = ImageEnhance.Brightness(rotated)
    bright = enhancer.enhance(1.15)
    
    # 在右上角加劍氣光暈
    result = bright.copy()
    pixels = result.load()
    w, h = result.size
    for y in range(h//4):
        for x in range(w*3//4, w):
            dist = ((x - w)**2 + y**2) ** 0.5
            if dist < w//4:
                r, g, b, a = pixels[x, y]
                if a > 0:
                    # 加粉紅色調
                    pixels[x, y] = (min(255, r+40), max(0, g-20), min(255, b+30), a)
                else:
                    # 透明區域加淡粉紅光暈
                    alpha = max(0, int(80 - dist * 3))
                    if alpha > 0:
                        pixels[x, y] = (255, 180, 220, alpha)
    
    return result

def make_bigwin_state(idle_img):
    """
    大獎狀態：
    - 整體向上位移（跳起）
    - 加金色光暈
    - 縮放略大（興奮膨脹感）
    """
    img = idle_img.copy()
    w, h = img.size
    
    # 向上位移 4px（跳起感）
    shifted = Image.new("RGBA", (w, h), (0, 0, 0, 0))
    shifted.paste(img, (0, -4))
    
    # 亮度提高
    enhancer = ImageEnhance.Brightness(shifted)
    bright = enhancer.enhance(1.2)
    
    # 加金色光暈（全身）
    result = bright.copy()
    pixels = result.load()
    for y in range(h):
        for x in range(w):
            r, g, b, a = pixels[x, y]
            if a > 50:
                # 加金色調
                pixels[x, y] = (min(255, r+15), min(255, g+10), max(0, b-10), a)
    
    # 在周圍加星星效果（幾個亮點）
    star_positions = [(w//4, h//6), (w*3//4, h//8), (w//8, h//3), (w*7//8, h//4)]
    for sx, sy in star_positions:
        for dy in range(-2, 3):
            for dx in range(-2, 3):
                if abs(dx) + abs(dy) <= 2:
                    nx, ny = sx+dx, sy+dy
                    if 0 <= nx < w and 0 <= ny < h:
                        pixels[nx, ny] = (255, 240, 100, 200)
    
    return result

def generate_all_states():
    chars = [
        ("chiikawa", "chiikawa_0.png"),
        ("hachiware", "hachiware_0.png"),
        ("usagi", "usagi_ref2_0.png"),
    ]
    
    for char_name, ref_file in chars:
        ref_path = os.path.join(REF_DIR, ref_file)
        if not os.path.exists(ref_path):
            print(f"SKIP {char_name}: {ref_file} not found")
            continue
        
        print(f"\n[{char_name}]")
        
        # 生成 idle
        idle = process_to_pixel(ref_path, size=32)
        if idle is None:
            print(f"  FAILED to process {ref_file}")
            continue
        
        idle_64 = idle.resize((64, 64), Image.NEAREST)
        idle_64.save(os.path.join(OUT_DIR, f"{char_name}_idle.png"))
        print(f"  idle: saved")
        
        # 生成 attack（基於 idle 變形）
        attack = make_attack_state(idle)
        attack_64 = attack.resize((64, 64), Image.NEAREST)
        attack_64.save(os.path.join(OUT_DIR, f"{char_name}_attack.png"))
        print(f"  attack: saved (rotated + pink glow)")
        
        # 生成 bigwin（基於 idle 變形）
        bigwin = make_bigwin_state(idle)
        bigwin_64 = bigwin.resize((64, 64), Image.NEAREST)
        bigwin_64.save(os.path.join(OUT_DIR, f"{char_name}_bigwin.png"))
        print(f"  bigwin: saved (shifted up + gold glow + stars)")

if __name__ == "__main__":
    print("Generating character states from reference images...")
    generate_all_states()
    print("\nDone!")
