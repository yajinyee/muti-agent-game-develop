# -*- coding: utf-8 -*-
"""
升級角色 idle 動畫：4幀 → 8幀
更流暢的呼吸感動畫，使用正弦波插值
Spritesheet 格式升級：
  - idle:   8幀（行0）
  - attack: 4幀（行1，補1幀）
  - bigwin: 4幀（行2）
  總計：8 cols × 3 rows = 24 幀
"""
import os
import json
import numpy as np
from PIL import Image, ImageEnhance

CHARS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"
SHEET_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\sheets"
FRAME_SIZE = 96

def gen_idle_8frames(base_img):
    """
    8幀 idle 動畫（正弦波呼吸感）
    幀0: 原始（0°）
    幀1: 輕微上移 1px（45°）
    幀2: 上移 2px + 放大 1.02x（90°，吸氣頂點）
    幀3: 上移 1px（135°）
    幀4: 原始（180°）
    幀5: 下移 1px（225°）
    幀6: 下移 1px + 縮小 0.99x（270°，呼氣底點）
    幀7: 下移 0.5px（315°）
    """
    import math
    frames = []
    
    for i in range(8):
        angle = i * math.pi / 4  # 0, 45, 90, 135, 180, 225, 270, 315 度
        # 正弦波：-1 到 1
        sin_val = math.sin(angle)
        
        # 位移：最大 ±2px
        offset_y = int(-sin_val * 2)  # 負值 = 上移
        
        # 縮放：1.0 到 1.02（吸氣時放大）
        scale = 1.0 + max(0, sin_val) * 0.02
        
        frame = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), (0, 0, 0, 0))
        
        if abs(scale - 1.0) > 0.001:
            new_size = int(FRAME_SIZE * scale)
            scaled = base_img.resize((new_size, new_size), Image.NEAREST)
            paste_x = (FRAME_SIZE - new_size) // 2
            paste_y = (FRAME_SIZE - new_size) // 2 + offset_y
            # 邊界保護
            paste_x = max(0, min(FRAME_SIZE - new_size, paste_x))
            paste_y = max(0, min(FRAME_SIZE - new_size, paste_y))
            frame.paste(scaled, (paste_x, paste_y))
        else:
            # 直接貼上，只做位移
            if offset_y != 0:
                temp = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), (0, 0, 0, 0))
                temp.paste(base_img, (0, offset_y))
                frame = temp
            else:
                frame = base_img.copy()
        
        frames.append(frame)
    
    return frames

def gen_attack_4frames(base_img):
    """4幀 attack（原3幀 + 補1幀）"""
    frames = []
    
    # 幀0：舉棒準備（向右上傾斜 -18度）
    f0 = base_img.rotate(-18, expand=False, fillcolor=(0,0,0,0))
    f0_canvas = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), (0,0,0,0))
    f0_canvas.paste(f0, (0, -3))
    f0_bright = ImageEnhance.Brightness(f0_canvas).enhance(1.12)
    frames.append(f0_bright)
    
    # 幀1：揮下衝擊（向左下傾斜 +12度 + 劍氣）
    f1 = base_img.rotate(12, expand=False, fillcolor=(0,0,0,0))
    f1_arr = np.array(f1)
    for i in range(8):
        x = FRAME_SIZE - 6 - i*3
        y = 6 + i*3
        if 0 <= x < FRAME_SIZE and 0 <= y < FRAME_SIZE:
            alpha = max(0, 220 - i*25)
            f1_arr[y, x] = [255, 150, 200, alpha]
            if x+1 < FRAME_SIZE:
                f1_arr[y, x+1] = [255, 180, 220, alpha//2]
    frames.append(Image.fromarray(f1_arr))
    
    # 幀2：收回（原始，輕微下移）
    f2 = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), (0,0,0,0))
    f2.paste(base_img, (0, 1))
    frames.append(f2)
    
    # 幀3：完全收回（原始）
    frames.append(base_img.copy())
    
    return frames

def gen_bigwin_4frames(base_img):
    """4幀 bigwin（跳起慶祝）"""
    frames = []
    configs = [
        (0,   1.00, False),
        (-10, 1.05, False),
        (-14, 1.08, True),
        (-5,  1.02, False),
    ]
    
    for offset, scale, add_stars in configs:
        frame = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), (0, 0, 0, 0))
        if scale != 1.0:
            new_size = int(FRAME_SIZE * scale)
            scaled = base_img.resize((new_size, new_size), Image.NEAREST)
            paste_x = (FRAME_SIZE - new_size) // 2
            paste_y = (FRAME_SIZE - new_size) // 2 + offset
            paste_x = max(0, min(FRAME_SIZE - new_size, paste_x))
            paste_y = max(0, min(FRAME_SIZE - new_size, paste_y))
            frame.paste(scaled, (paste_x, paste_y))
        else:
            frame.paste(base_img, (0, offset if offset >= 0 else 0))
        
        if add_stars:
            frame_arr = np.array(frame)
            star_positions = [
                (8, 8), (FRAME_SIZE-10, 8),
                (6, FRAME_SIZE//3), (FRAME_SIZE-8, FRAME_SIZE//3)
            ]
            for sx, sy in star_positions:
                if 0 <= sx < FRAME_SIZE and 0 <= sy < FRAME_SIZE:
                    frame_arr[sy, sx] = [255, 235, 50, 255]
                    for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
                        nx, ny = sx+dx, sy+dy
                        if 0 <= nx < FRAME_SIZE and 0 <= ny < FRAME_SIZE:
                            frame_arr[ny, nx] = [255, 235, 50, 180]
            frame = Image.fromarray(frame_arr)
        
        frames.append(frame)
    
    return frames

def build_spritesheet(char_name):
    """建立 8×3 Spritesheet"""
    idle_path = os.path.join(CHARS_DIR, f"{char_name}_idle.png")
    if not os.path.exists(idle_path):
        print(f"  SKIP: {idle_path} not found")
        return False
    
    base = Image.open(idle_path).convert("RGBA")
    if base.width != FRAME_SIZE:
        base = base.resize((FRAME_SIZE, FRAME_SIZE), Image.NEAREST)
    
    idle_frames   = gen_idle_8frames(base)   # 8幀
    attack_frames = gen_attack_4frames(base) # 4幀
    bigwin_frames = gen_bigwin_4frames(base) # 4幀
    
    # 建立 Spritesheet（8 cols × 3 rows）
    COLS = 8
    ROWS = 3
    sheet = Image.new("RGBA", (FRAME_SIZE * COLS, FRAME_SIZE * ROWS), (0, 0, 0, 0))
    
    all_frames = [idle_frames, attack_frames, bigwin_frames]
    for row, frames in enumerate(all_frames):
        for col, frame in enumerate(frames[:COLS]):
            sheet.paste(frame, (col * FRAME_SIZE, row * FRAME_SIZE))
    
    # 儲存
    sheet_path = os.path.join(SHEET_DIR, f"{char_name}_animated.png")
    sheet.save(sheet_path)
    
    # 更新 metadata
    meta = {
        "char": char_name,
        "frame_size": FRAME_SIZE,
        "cols": COLS,
        "rows": ROWS,
        "animations": {
            "idle":   {"row": 0, "frames": 8, "fps": 8.0},   # 升級：4→8幀，4→8fps
            "attack": {"row": 1, "frames": 3, "fps": 8.0},
            "bigwin": {"row": 2, "frames": 4, "fps": 6.0},
        }
    }
    meta_path = os.path.join(SHEET_DIR, f"{char_name}_animated.json")
    with open(meta_path, "w", encoding="utf-8") as f:
        json.dump(meta, f, indent=2)
    
    print(f"  ✅ {char_name}: {sheet.width}x{sheet.height} (8×3 frames)")
    return True

def main():
    print("=== 升級角色 idle 動畫：4幀 → 8幀 ===\n")
    for char in ["chiikawa", "hachiware", "usagi"]:
        print(f"[{char}]")
        build_spritesheet(char)
    print("\n✅ 完成！")

if __name__ == "__main__":
    main()
