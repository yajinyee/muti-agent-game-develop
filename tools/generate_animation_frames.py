# -*- coding: utf-8 -*-
"""
動畫幀生成器
從單張角色圖生成多幀動畫 Spritesheet
支援：idle(4幀) / attack(3幀) / bigwin(4幀)
"""
import cv2
import numpy as np
from PIL import Image, ImageEnhance, ImageFilter
import os
import math

OUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"
SHEET_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\sheets"

FRAME_SIZE = 96  # 每幀尺寸

def img_to_pixel(img_pil, size=48, n_colors=16):
    """把 PIL 圖片轉成像素藝術風格"""
    # 縮小
    small = img_pil.resize((size, size), Image.LANCZOS)
    # 轉 numpy
    arr = np.array(small.convert("RGB"))
    # K-means 量化
    pixels = arr.reshape(-1, 3).astype(np.float32)
    criteria = (cv2.TERM_CRITERIA_EPS + cv2.TERM_CRITERIA_MAX_ITER, 20, 1.0)
    _, labels, centers = cv2.kmeans(pixels, n_colors, None, criteria, 10, cv2.KMEANS_RANDOM_CENTERS)
    centers = np.uint8(centers)
    quantized = centers[labels.flatten()].reshape(arr.shape)
    # 放大回 FRAME_SIZE
    result = cv2.resize(quantized, (FRAME_SIZE, FRAME_SIZE), interpolation=cv2.INTER_NEAREST)
    return Image.fromarray(result)

def remove_bg_and_pixelate(source_path, size=48, n_colors=16):
    """讀取圖片，去除白色背景，轉像素藝術"""
    img = Image.open(source_path).convert("RGBA")
    
    # 如果有透明通道，直接用
    r, g, b, a = img.split()
    
    # 去除白色背景（flood fill）
    from collections import deque
    w, h = img.size
    pixels = img.load()
    bg_mask = [[False]*h for _ in range(w)]
    queue = deque()
    
    for x in range(w):
        for y in [0, h-1]:
            r_val, g_val, b_val, a_val = pixels[x, y]
            if r_val > 220 and g_val > 220 and b_val > 220:
                queue.append((x, y))
                bg_mask[x][y] = True
    for y in range(h):
        for x in [0, w-1]:
            r_val, g_val, b_val, a_val = pixels[x, y]
            if r_val > 220 and g_val > 220 and b_val > 220 and not bg_mask[x][y]:
                queue.append((x, y))
                bg_mask[x][y] = True
    
    while queue:
        x, y = queue.popleft()
        for dx, dy in [(0,1),(0,-1),(1,0),(-1,0)]:
            nx, ny = x+dx, y+dy
            if 0 <= nx < w and 0 <= ny < h and not bg_mask[nx][ny]:
                r_val, g_val, b_val, a_val = pixels[nx, ny]
                if r_val > 215 and g_val > 215 and b_val > 215:
                    bg_mask[nx][ny] = True
                    queue.append((nx, ny))
    
    # 套用遮罩
    result = img.copy()
    result_pixels = result.load()
    for y in range(h):
        for x in range(w):
            if bg_mask[x][y]:
                result_pixels[x, y] = (0, 0, 0, 0)
    
    # 裁切非透明區域
    bbox = result.getbbox()
    if bbox:
        result = result.crop(bbox)
    
    # 縮放到目標尺寸
    result = result.resize((size, size), Image.LANCZOS)
    
    # 像素化 RGB 部分
    rgb = result.convert("RGB")
    arr = np.array(rgb)
    pixels_arr = arr.reshape(-1, 3).astype(np.float32)
    criteria = (cv2.TERM_CRITERIA_EPS + cv2.TERM_CRITERIA_MAX_ITER, 20, 1.0)
    _, labels, centers = cv2.kmeans(pixels_arr, n_colors, None, criteria, 10, cv2.KMEANS_RANDOM_CENTERS)
    centers = np.uint8(centers)
    quantized = centers[labels.flatten()].reshape(arr.shape)
    
    # 放大
    quantized_large = cv2.resize(quantized, (FRAME_SIZE, FRAME_SIZE), interpolation=cv2.INTER_NEAREST)
    
    # 重新套用透明通道
    alpha_arr = np.array(result.split()[3])
    alpha_large = cv2.resize(alpha_arr, (FRAME_SIZE, FRAME_SIZE), interpolation=cv2.INTER_NEAREST)
    _, alpha_binary = cv2.threshold(alpha_large, 128, 255, cv2.THRESH_BINARY)
    
    final = Image.fromarray(quantized_large).convert("RGBA")
    final_arr = np.array(final)
    final_arr[:, :, 3] = alpha_binary
    
    return Image.fromarray(final_arr)

def gen_idle_frames(base_img):
    """
    Idle 動畫：4幀呼吸感搖擺
    幀0: 原始（站立）
    幀1: 上移 2px + 輕微放大 1.03x（吸氣）
    幀2: 原始（站立）
    幀3: 下移 1px + 輕微縮小 0.98x（呼氣）
    """
    frames = []
    configs = [
        (0, 1.00),   # 幀0：原始
        (-2, 1.03),  # 幀1：上移+放大（吸氣）
        (0, 1.00),   # 幀2：原始
        (1, 0.98),   # 幀3：下移+縮小（呼氣）
    ]

    for offset, scale in configs:
        frame = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), (0, 0, 0, 0))
        if scale != 1.0:
            new_size = int(FRAME_SIZE * scale)
            scaled = base_img.resize((new_size, new_size), Image.NEAREST)
            # 置中貼上
            paste_x = (FRAME_SIZE - new_size) // 2
            paste_y = (FRAME_SIZE - new_size) // 2 + offset
            frame.paste(scaled, (paste_x, paste_y))
        else:
            frame.paste(base_img, (0, offset))
        frames.append(frame)

    return frames

def gen_attack_frames(base_img):
    """
    Attack 動畫：3幀揮棒（更明顯的動作感）
    幀0: 舉棒準備（向右上傾斜 -18度，上移 3px）
    幀1: 揮下衝擊（向左下傾斜 +12度，加劍氣光效）
    幀2: 收回（原始，輕微下移 1px）
    """
    frames = []

    # 幀0：舉棒準備
    f0 = base_img.rotate(-18, expand=False, fillcolor=(0,0,0,0))
    f0_canvas = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), (0,0,0,0))
    f0_canvas.paste(f0, (0, -3))
    f0_bright = ImageEnhance.Brightness(f0_canvas).enhance(1.12)
    frames.append(f0_bright)

    # 幀1：揮下衝擊
    f1 = base_img.rotate(12, expand=False, fillcolor=(0,0,0,0))
    f1_arr = np.array(f1)
    # 右上角加劍氣光點（粉紅/藍/黃依角色）
    for i in range(8):
        x = FRAME_SIZE - 6 - i*3
        y = 6 + i*3
        if 0 <= x < FRAME_SIZE and 0 <= y < FRAME_SIZE:
            alpha = max(0, 220 - i*25)
            f1_arr[y, x] = [255, 150, 200, alpha]
            if x+1 < FRAME_SIZE:
                f1_arr[y, x+1] = [255, 180, 220, alpha//2]
    frames.append(Image.fromarray(f1_arr))

    # 幀2：收回
    f2 = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), (0,0,0,0))
    f2.paste(base_img, (0, 1))
    frames.append(f2)

    return frames

def gen_bigwin_frames(base_img):
    """
    BigWin 動畫：4幀跳起慶祝
    幀0: 原始（準備跳）
    幀1: 上移 10px + 放大 1.05x（跳起）
    幀2: 上移 14px + 放大 1.08x（最高點）+ 金色星星
    幀3: 上移 5px（落下）
    """
    frames = []
    configs = [
        (0,   1.00, False),
        (-10, 1.05, False),
        (-14, 1.08, True),   # 最高點加星星
        (-5,  1.02, False),
    ]

    for offset, scale, add_stars in configs:
        frame = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), (0, 0, 0, 0))
        if scale != 1.0:
            new_size = int(FRAME_SIZE * scale)
            scaled = base_img.resize((new_size, new_size), Image.NEAREST)
            paste_x = (FRAME_SIZE - new_size) // 2
            paste_y = (FRAME_SIZE - new_size) // 2 + offset
            frame.paste(scaled, (paste_x, paste_y))
        else:
            frame.paste(base_img, (0, offset))

        if add_stars:
            frame_arr = np.array(frame)
            # 加金色星星（4個角落）
            star_positions = [
                (8, 8), (FRAME_SIZE-10, 8),
                (6, FRAME_SIZE//3), (FRAME_SIZE-8, FRAME_SIZE//3)
            ]
            for sx, sy in star_positions:
                sy_adj = max(0, min(FRAME_SIZE-1, sy + offset))
                if 0 <= sx < FRAME_SIZE and 0 <= sy_adj < FRAME_SIZE:
                    frame_arr[sy_adj, sx] = [255, 235, 50, 255]
                    for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
                        nx, ny = sx+dx, sy_adj+dy
                        if 0 <= nx < FRAME_SIZE and 0 <= ny < FRAME_SIZE:
                            frame_arr[ny, nx] = [255, 235, 50, 180]
            frame = Image.fromarray(frame_arr)

        frames.append(frame)

    return frames

def create_animation_spritesheet(char_name, source_path=None):
    """
    建立角色動畫 Spritesheet
    格式：每行一個動畫狀態，每列一幀
    idle:   4幀 → 行0
    attack: 3幀 → 行1（補1幀到4幀）
    bigwin: 4幀 → 行2
    """
    # 載入基礎圖
    if source_path and os.path.exists(source_path):
        print(f"  Using real image: {source_path}")
        base = remove_bg_and_pixelate(source_path)
    else:
        # 用現有的程式生成圖
        idle_path = os.path.join(OUT_DIR, f"{char_name}_idle.png")
        if not os.path.exists(idle_path):
            print(f"  SKIP: {idle_path} not found")
            return False
        base = Image.open(idle_path).convert("RGBA")
        if base.width != FRAME_SIZE:
            base = base.resize((FRAME_SIZE, FRAME_SIZE), Image.NEAREST)
        print(f"  Using generated image: {idle_path}")
    
    # 生成各狀態幀
    idle_frames   = gen_idle_frames(base)    # 4幀
    attack_frames = gen_attack_frames(base)  # 3幀 → 補到4幀
    attack_frames.append(base.copy())        # 補第4幀
    bigwin_frames = gen_bigwin_frames(base)  # 4幀
    
    # 建立 Spritesheet（4列 × 3行）
    cols = 4
    rows = 3
    sheet = Image.new("RGBA", (FRAME_SIZE * cols, FRAME_SIZE * rows), (0, 0, 0, 0))
    
    all_frames = [idle_frames, attack_frames, bigwin_frames]
    state_names = ["idle", "attack", "bigwin"]
    
    for row, (frames, state) in enumerate(zip(all_frames, state_names)):
        for col, frame in enumerate(frames[:cols]):
            sheet.paste(frame, (col * FRAME_SIZE, row * FRAME_SIZE))
        # 注意：不覆蓋 characters/ 下的 sprites，那些由 generate_pixel_art_v5.py 管理
    
    # 儲存 Spritesheet
    sheet_path = os.path.join(SHEET_DIR, f"{char_name}_animated.png")
    sheet.save(sheet_path)
    print(f"  Spritesheet: {sheet.width}x{sheet.height} -> {sheet_path}")
    
    # 儲存 metadata
    import json
    meta = {
        "char": char_name,
        "frame_size": FRAME_SIZE,
        "cols": cols,
        "rows": rows,
        "animations": {
            "idle":   {"row": 0, "frames": 4, "fps": 4},
            "attack": {"row": 1, "frames": 3, "fps": 8},
            "bigwin": {"row": 2, "frames": 4, "fps": 6},
        }
    }
    meta_path = os.path.join(SHEET_DIR, f"{char_name}_animated.json")
    with open(meta_path, "w") as f:
        json.dump(meta, f, indent=2)
    print(f"  Metadata: {meta_path}")
    
    return True

def main():
    ref_dir = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\reference"
    
    chars = [
        ("chiikawa",  os.path.join(ref_dir, "chiikawa_real.png")),
        ("hachiware", os.path.join(ref_dir, "hachiware_real.png")),
        ("usagi",     os.path.join(ref_dir, "usagi_real.png")),
    ]
    
    print("Generating animation spritesheets...")
    for char_name, real_path in chars:
        print(f"\n[{char_name}]")
        # 如果有真實圖片就用，否則用程式生成的
        source = real_path if os.path.exists(real_path) else None
        create_animation_spritesheet(char_name, source)
    
    print("\nDone!")

if __name__ == "__main__":
    main()
